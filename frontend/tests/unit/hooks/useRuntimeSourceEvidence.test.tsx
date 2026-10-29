import type { ReactNode } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, renderHook, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import type { EvidenceScope, RuleEvidenceView } from '@/types/runtimeRuleEvidence';
import type { RuntimeSourceSelection } from '@/types/studioV2RuntimeContext';
import { useRuntimeSourceEvidence } from '@/hooks/datalink/useRuntimeSourceEvidence';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { evidenceContext, evidenceScope, ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
const selection: RuntimeSourceSelection = { rule_id: 'R', rule_revision: 'R2', point_id: 'P', mapping_id: 'M' };
const clients: QueryClient[] = [];
const unavailable = { ok: false, reason: 'unavailable', action: 'refresh' };
function response(scope = evidenceScope, selected = selection) {
  const evidence: RuleEvidenceView = { ...ruleEvidenceFixture(scope), source: {
    rule_id: selected.rule_id, point_id: selected.point_id, mapping_id: selected.mapping_id ?? null,
    tag_id: selected.mapping_id ? 'T' : null, saved_revision: selected.rule_revision, effective_revision: null,
    state: 'saved-only', proof: 'unavailable', effective: null,
    saved: { binding: { workspace_id: scope.workspace_id, workspace_revision: scope.workspace_revision, entity_revision: selected.rule_revision },
      mapping_revision: selected.mapping_id ? 'MR2' : null, enabled: true, address: '40001', protocol_function: '03', read_type: 'float32',
      result_type: selected.mapping_id ? 'float64' : null, polling: { group_id: 'PG', interval_ms: 1000, source: 'polling-group' }, transforms: [] },
  } };
  return { data: { success: true, data: evidenceContext(evidence, scope.workspace_id) } };
}
function deferred() {
  let resolve!: (value: ReturnType<typeof response>) => void;
  const promise = new Promise<ReturnType<typeof response>>((done) => { resolve = done; });
  return { promise, resolve };
}
function harness(scope: EvidenceScope | null = evidenceScope, selected: RuntimeSourceSelection | null = selection) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  const wrapper = ({ children }: { children: ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  return renderHook(({ scope, selected }: { scope: EvidenceScope | null; selected: RuntimeSourceSelection | null }) => useRuntimeSourceEvidence(scope, selected), { wrapper, initialProps: { scope, selected } });
}
beforeEach(() => vi.mocked(studioV2DatalinkApi.get).mockReset());
afterEach(() => { cleanup(); clients.splice(0).forEach((client) => client.clear()); });

describe('scoped single-source runtime evidence query', () => {
  it('uses only the captured read scope and explicit source selection with an abort signal', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response()); const view = harness();
    await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/runtime-context', { params: { ...evidenceScope, ...selection }, signal: expect.any(AbortSignal) });
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
  });
  it('aborts A and never replaces selected B with a late A response', async () => {
    const a = deferred(); const b = deferred(); vi.mocked(studioV2DatalinkApi.get).mockReturnValueOnce(a.promise).mockReturnValueOnce(b.promise);
    const view = harness(); await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1));
    const signalA = vi.mocked(studioV2DatalinkApi.get).mock.calls[0][1]?.signal;
    const scopeB = { ...evidenceScope, device_id: 'B' }; const sourceB = { ...selection, rule_id: 'B1', point_id: 'PB', mapping_id: 'MB' };
    view.rerender({ scope: scopeB, selected: sourceB }); expect(view.result.current.data).toBeUndefined();
    await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2)); expect(signalA?.aborted).toBe(true);
    await act(async () => { b.resolve(response(scopeB, sourceB)); }); await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
    await act(async () => { a.resolve(response()); });
    const result = view.result.current.data; if (!result?.ok) throw new Error('expected B evidence');
    expect(result.value.device_id).toBe('B'); expect(result.value.source?.rule_id).toBe('B1');
  });
  it.each(['rule_id', 'rule_revision', 'point_id', 'mapping_id'] as const)('keys changed source %s separately with no prior placeholder', async (field) => {
    const next = deferred(); vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response()).mockReturnValueOnce(next.promise);
    const view = harness(); await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
    const selected = { ...selection, [field]: 'changed' }; view.rerender({ scope: evidenceScope, selected });
    expect(view.result.current.data).toBeUndefined(); await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2));
    await act(async () => { next.resolve(response(evidenceScope, selected)); }); await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
  });
  it.each(['workspace_id', 'workspace_revision', 'device_id'] as const)('keys changed scope %s separately', async (field) => {
    const next = deferred(); vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response()).mockReturnValueOnce(next.promise);
    const view = harness(); await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
    const scope = { ...evidenceScope, [field]: 'changed' }; view.rerender({ scope, selected: selection });
    expect(view.result.current.data).toBeUndefined(); await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2));
    await act(async () => { next.resolve(response(scope)); }); await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
  });
  it.each(['scope', 'selection'])('keeps absent %s disabled after a prior successful selection', async (missing) => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response()); const view = harness();
    await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
    view.rerender({ scope: missing === 'scope' ? null : evidenceScope, selected: missing === 'selection' ? null : selection });
    expect(view.result.current.data).toBeUndefined();
    await act(async () => { await view.result.current.refetch(); }); await waitFor(() => expect(view.result.current.data).toEqual(unavailable));
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
  });
  it('does not transport or retry invalid selection and typed failures', async () => {
    const view = harness(evidenceScope, { ...selection, rule_id: '/private/secret' });
    await waitFor(() => expect(view.result.current.data).toEqual(unavailable)); expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
    vi.mocked(studioV2DatalinkApi.get).mockRejectedValueOnce({ response: { data: { error: { code: 'stale-reference', message: '/private/secret' } } } });
    view.rerender({ scope: evidenceScope, selected: selection });
    await waitFor(() => expect(view.result.current.data).toEqual({ ok: false, reason: 'stale-reference', action: 'refresh' }));
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1); expect(JSON.stringify(view.result.current.data)).not.toContain('secret');
  });
});
