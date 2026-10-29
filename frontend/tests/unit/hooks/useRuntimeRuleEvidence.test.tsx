import type { ReactNode } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, renderHook, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import type { EvidenceScope } from '@/types/runtimeRuleEvidence';
import { useRuntimeRuleEvidence } from '@/hooks/datalink/useRuntimeRuleEvidence';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { evidenceContext, evidenceScope, ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
const clients: QueryClient[] = [];
const unavailable = { ok: false, reason: 'unavailable', action: 'refresh' };
const response = (scope = evidenceScope, revision = 'A') => ({ data: { success: true, data: evidenceContext(ruleEvidenceFixture(scope, revision), scope.workspace_id) } });
function deferred() {
  let resolve!: (value: ReturnType<typeof response>) => void;
  const promise = new Promise<ReturnType<typeof response>>((done) => { resolve = done; });
  return { promise, resolve };
}
function harness(scope: EvidenceScope | null = evidenceScope) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  const wrapper = ({ children }: { children: ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  return renderHook(({ selected }: { selected: EvidenceScope | null }) => useRuntimeRuleEvidence(selected), { wrapper, initialProps: { selected: scope } });
}
beforeEach(() => { vi.mocked(studioV2DatalinkApi.get).mockReset(); });
afterEach(() => { cleanup(); clients.splice(0).forEach((client) => client.clear()); });

describe('shared scoped runtime rule evidence query', () => {
  it('reads server A then B without calculating a boundary or joining editable groups', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response()).mockResolvedValueOnce(response(evidenceScope, 'B'));
    const view = harness();
    await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
    const first = view.result.current.data;
    if (!first?.ok) throw new Error('missing A evidence');
    expect(first.value.groups[0]).toMatchObject({ saved_revision: 'C', applied_revision: 'B', effective_revision: 'A' });
    await act(async () => { await view.result.current.refetch(); });
    await waitFor(() => { const result = view.result.current.data; expect(result?.ok && result.value.groups[0].effective_revision).toBe('B'); });
    const second = view.result.current.data;
    if (!second?.ok) throw new Error('missing B evidence');
    expect(second.value.groups[0]).toMatchObject({ saved_revision: 'C', applied_revision: 'B', effective_revision: 'B' });
    expect(second.value.observed_at).toBe('2026-10-27T10:00:01Z');
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2);
    expect(vi.mocked(studioV2DatalinkApi.get).mock.calls.every(([url]) => url === '/studio-v2/workspace/runtime-context')).toBe(true);
  });
  it('aborts A and never replaces selected B with a late A response', async () => {
    const a = deferred(); const b = deferred();
    vi.mocked(studioV2DatalinkApi.get).mockReturnValueOnce(a.promise).mockReturnValueOnce(b.promise);
    const view = harness();
    await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1));
    const signalA = vi.mocked(studioV2DatalinkApi.get).mock.calls[0][1]?.signal;
    const scopeB = { ...evidenceScope, device_id: 'E' };
    view.rerender({ selected: scopeB });
    expect(view.result.current.data).toBeUndefined();
    await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2));
    expect(signalA?.aborted).toBe(true);
    await act(async () => { b.resolve(response(scopeB)); });
    await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
    await act(async () => { a.resolve(response()); });
    const result = view.result.current.data;
    if (!result?.ok) throw new Error('missing selected evidence');
    expect(result.value.device_id).toBe('E');
    expect(result.value.groups[0].effective?.members[0].device_id).toBe('E');
  });
  it.each([
    { ...evidenceScope, workspace_id: 'other' },
    { ...evidenceScope, workspace_revision: 'W2' },
    { ...evidenceScope, device_id: 'E' },
  ])('does not reuse completed data under a changed scope %s', async (selected) => {
    const next = deferred();
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response()).mockReturnValueOnce(next.promise);
    const view = harness();
    await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
    view.rerender({ selected });
    expect(view.result.current.data).toBeUndefined();
    await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2));
    await act(async () => { next.resolve(response(selected)); });
    await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
    const result = view.result.current.data;
    if (!result?.ok) throw new Error('missing scoped evidence');
    expect([result.value.workspace_id, result.value.workspace_revision, result.value.device_id]).toEqual([selected.workspace_id, selected.workspace_revision, selected.device_id]);
  });
  it('keeps missing scope disabled even after a prior successful selection', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response());
    const view = harness();
    await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
    view.rerender({ selected: null });
    expect(view.result.current.data).toBeUndefined();
    await act(async () => { await view.result.current.refetch(); });
    await waitFor(() => expect(view.result.current.data).toEqual(unavailable));
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
  });
  it('rejects invalid scope without transport or cached evidence', async () => {
    const view = harness({ ...evidenceScope, device_id: '/private/secret' });
    await waitFor(() => expect(view.result.current.data).toEqual(unavailable));
    expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
  });
  it('exposes safe typed stale-reference results without automatic retry or raw diagnostics', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockRejectedValueOnce({ response: { data: { error: { code: 'stale-reference', message: '/private/secret' } } } });
    const view = harness();
    await waitFor(() => expect(view.result.current.data).toEqual({ ok: false, reason: 'stale-reference', action: 'refresh' }));
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
    expect(JSON.stringify(view.result.current.data)).not.toContain('secret');
  });
});
