import type { ReactNode } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, renderHook, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import type { EvidenceScope } from '@/types/runtimeRuleEvidence';
import type { RuntimeSourcePageRequest } from '@/types/studioV2RuntimeContext';
import { useRuntimeSourcePage } from '@/hooks/datalink/useRuntimeSourcePage';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
const scope: EvidenceScope = { workspace_id: 'W', workspace_revision: 'W1', device_id: 'A' };
const request: RuntimeSourcePageRequest = { kind: 'rules' };
const clients: QueryClient[] = [];
const unavailable = { ok: false, reason: 'unavailable', action: 'refresh' };
function response(observed = scope, selected = request) {
  const rules = selected.kind === 'rules' ? [{ id: observed.device_id === 'A' ? 'A6' : 'B1', device_id: observed.device_id, start_address: '40001', count: 1, data_type: 'float32', naming_prefix: 'Meter', enabled: true, revision_id: 'R6' }] : [];
  const mappings = selected.kind === 'mappings' ? [{ id: 'M', rule_id: selected.rule_id ?? 'A6', rule_revision: selected.rule_revision ?? 'R6', device_id: observed.device_id, point_id: 'P', address: '40001', enabled: true }] : [];
  return { data: { success: true, data: { workspace_id: observed.workspace_id, source_page: { contract_version: 1, ...observed, kind: selected.kind, query: selected.query?.trim() ?? '', limit: Math.min(selected.limit ?? 50, 100), list_revision: 'L1', rules, mappings, next_cursor: null } } } };
}
function deferred() {
  let resolve!: (value: ReturnType<typeof response>) => void;
  const promise = new Promise<ReturnType<typeof response>>((done) => { resolve = done; });
  return { promise, resolve };
}
function harness(observed: EvidenceScope | null = scope, selected: RuntimeSourcePageRequest | null = request) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  const wrapper = ({ children }: { children: ReactNode }) => <QueryClientProvider client={client}>{children}</QueryClientProvider>;
  return renderHook(({ observed, selected }: { observed: EvidenceScope | null; selected: RuntimeSourcePageRequest | null }) => useRuntimeSourcePage(observed, selected), { wrapper, initialProps: { observed, selected } });
}
beforeEach(() => vi.mocked(studioV2DatalinkApi.get).mockReset());
afterEach(() => { cleanup(); clients.splice(0).forEach((client) => client.clear()); });

describe('bounded source page query', () => {
  it('makes one page request without prefetching any row, group or detail', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response()); const view = harness();
    await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/runtime-context', { params: { ...scope, source_list: 'rules', q: '', limit: 50 }, signal: expect.any(AbortSignal) });
  });
  it('aborts delayed A6 and never replaces B1 after a device switch', async () => {
    const a = deferred(); const b = deferred(); vi.mocked(studioV2DatalinkApi.get).mockReturnValueOnce(a.promise).mockReturnValueOnce(b.promise);
    const view = harness(); await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1));
    const signalA = vi.mocked(studioV2DatalinkApi.get).mock.calls[0][1]?.signal; const observed = { ...scope, device_id: 'B' };
    view.rerender({ observed, selected: request }); expect(view.result.current.data).toBeUndefined();
    await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2)); expect(signalA?.aborted).toBe(true);
    await act(async () => { b.resolve(response(observed)); }); await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
    await act(async () => { a.resolve(response()); });
    const result = view.result.current.data; if (!result?.ok) throw new Error('expected B page');
    expect(result.value.rules.map((row) => row.id)).toEqual(['B1']);
  });
  it.each([
    { kind: 'mappings' as const }, { kind: 'rules' as const, query: 'Meter' },
    { kind: 'rules' as const, limit: 20 }, { kind: 'rules' as const, cursor: 'e30' },
    { kind: 'mappings' as const, rule_id: 'A6', rule_revision: 'R6' },
  ])('isolates changed kind/search/limit/cursor/filter %j', async (selected) => {
    const next = deferred(); vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response()).mockReturnValueOnce(next.promise);
    const view = harness(); await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
    view.rerender({ observed: scope, selected }); expect(view.result.current.data).toBeUndefined();
    await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2));
    await act(async () => { next.resolve(response(scope, selected)); }); await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2);
  });
  it.each(['rule_id', 'rule_revision'] as const)('isolates each parent filter field %s independently', async (field) => {
    const original: RuntimeSourcePageRequest = { kind: 'mappings', rule_id: 'A6', rule_revision: 'R6' }; const next = deferred();
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response(scope, original)).mockReturnValueOnce(next.promise);
    const view = harness(scope, original); await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
    const selected = { ...original, [field]: 'changed' }; view.rerender({ observed: scope, selected });
    expect(view.result.current.data).toBeUndefined(); await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2));
    await act(async () => { next.resolve(response(scope, selected)); }); await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
  });
  it.each(['workspace_id', 'workspace_revision'] as const)('isolates changed %s', async (field) => {
    const next = deferred(); vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response()).mockReturnValueOnce(next.promise);
    const view = harness(); await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
    const observed = { ...scope, [field]: 'changed' }; view.rerender({ observed, selected: request }); expect(view.result.current.data).toBeUndefined();
    await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(2));
    await act(async () => { next.resolve(response(observed)); }); await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
  });
  it.each(['scope', 'request'])('disables absent %s and never restores old page data', async (missing) => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(response()); const view = harness();
    await waitFor(() => expect(view.result.current.data?.ok).toBe(true));
    view.rerender({ observed: missing === 'scope' ? null : scope, selected: missing === 'request' ? null : request });
    expect(view.result.current.data).toBeUndefined();
    await act(async () => { await view.result.current.refetch(); }); await waitFor(() => expect(view.result.current.data).toEqual(unavailable));
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
  });
  it('returns a typed stale cursor result without automatically restarting at page one', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockRejectedValueOnce({ response: { data: { error: { code: 'stale-reference', message: '/private/secret' } } } });
    const view = harness(scope, { kind: 'rules', cursor: 'e30' });
    await waitFor(() => expect(view.result.current.data).toEqual({ ok: false, reason: 'stale-reference', action: 'refresh' }));
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1); expect(JSON.stringify(view.result.current.data)).not.toContain('secret');
  });
});
