import { onlineManager, QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, renderHook, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { useCircuitProfileStart } from '@/features/datalink/workbench-v2/steps/step4/useCircuitProfileStart';
import { saveBasicRecordingIntent } from '@/features/datalink/workbench-v2/state/basicRecording';
import { retainSchemaOperation } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery';
import { studioV2RecordingStartAPI as api } from '@/services/studioV2RecordingStart';
import type { RecordingStartOperation, RecordingStartRequest } from '@/types/studioV2RecordingStart';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), operation: vi.fn() } }));
const scope = 'W1/D1/circuit';
const request: RecordingStartRequest = { request_id: 'original-request', workspace_id: 'W1', expected_workspace_revision: 'old-head',
  device_ids: ['D1'], groups: [{ group_id: 'G1', expected_group_revision: 'old-group', expected_connector_revision: 'connector-1' }] };
const partial: RecordingStartOperation = { operation_id: 'O1', workspace_id: 'W1', action: 'recording_start', status: 'partial',
  next_action: 'retry_same_request', stage: 'activation', intent_digest: 'a'.repeat(64), setup_revision: 'new-head', device_ids: ['D1'],
  devices: [{ device_id: 'D1', activated: false }], groups: [{ group_id: 'G1', group_revision: 'new-group', saved: true, ready: true, applied: true }],
  created_at: '2026-10-28T00:00:00Z', updated_at: '2026-10-28T00:00:01Z' };
const clients: QueryClient[] = [];
function host() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false, networkMode: 'always' } } }); clients.push(client);
  const prepare = vi.fn();
  const result = renderHook(({ selected, readonly }: { selected: string; readonly: boolean }) => useCircuitProfileStart({
    scopeKey: selected, workspaceId: 'W1', deviceId: selected === scope ? 'D1' : 'D2', connectorRevision: 'connector-1',
    canStart: true, readonly, prepare,
  }), { initialProps: { selected: scope, readonly: false }, wrapper: ({ children }) => <QueryClientProvider client={client}>{children}</QueryClientProvider> });
  return { ...result, prepare };
}
function deferred() { let resolve!: (value: RecordingStartOperation) => void; const promise = new Promise<RecordingStartOperation>(done => { resolve = done; }); return { promise, resolve }; }
beforeEach(() => {
  vi.restoreAllMocks(); vi.clearAllMocks(); onlineManager.setOnline(true); sessionStorage.clear(); localStorage.clear();
  saveBasicRecordingIntent({ scope_key: scope, request, operation_id: 'O1' });
  vi.mocked(api.operation).mockResolvedValue(partial); vi.mocked(api.start).mockResolvedValue({ ...partial, status: 'succeeded', next_action: undefined });
});
afterEach(() => { cleanup(); vi.restoreAllMocks(); onlineManager.setOnline(true); clients.splice(0).forEach(client => client.clear()); retainSchemaOperation('W1', 'G1', undefined, 'schema-O1'); });

it('keeps immutable original revisions for an explicitly authorized retry without preparing new groups', async () => {
  const view = host(); await waitFor(() => expect(view.result.current.canRetrySameRequest).toBe(true));
  await act(async () => { await view.result.current.start(true); });
  expect(api.start).toHaveBeenCalledExactlyOnceWith(request); expect(view.prepare).not.toHaveBeenCalled();
});
it('does not revive a captured retry callback after owner A to B to A', async () => {
  const view = host(); await waitFor(() => expect(view.result.current.canRetrySameRequest).toBe(true));
  const stale = view.result.current.start; const before = vi.mocked(api.operation).mock.calls.length;
  view.rerender({ selected: 'W1/D2/circuit', readonly: false }); view.rerender({ selected: scope, readonly: false });
  await waitFor(() => expect(view.result.current.canRetrySameRequest).toBe(true));
  const afterRemount = vi.mocked(api.operation).mock.calls.length;
  await act(async () => { await stale(true); });
  expect(vi.mocked(api.operation).mock.calls.length).toBe(afterRemount); expect(afterRemount).toBeGreaterThanOrEqual(before);
  expect(api.start).not.toHaveBeenCalled();
});
it.each(['scope', 'readonly', 'replacement', 'schema'] as const)('rechecks %s after the fresh original-operation lookup before retry POST', async change => {
  const view = host(); await waitFor(() => expect(view.result.current.canRetrySameRequest).toBe(true));
  const pending = deferred(); vi.mocked(api.operation).mockReturnValueOnce(pending.promise);
  let running!: Promise<void>; act(() => { running = view.result.current.start(true); });
  await waitFor(() => expect(api.operation).toHaveBeenCalledTimes(2));
  if (change === 'scope') { view.rerender({ selected: 'W1/D2/circuit', readonly: false }); view.rerender({ selected: scope, readonly: false }); }
  if (change === 'readonly') view.rerender({ selected: scope, readonly: true });
  if (change === 'replacement') saveBasicRecordingIntent({ scope_key: scope, request: { ...request, request_id: 'R2' }, operation_id: 'O2' });
  if (change === 'schema') retainSchemaOperation('W1', 'G1', 'schema-O1');
  await act(async () => { pending.resolve(partial); await running; });
  expect(api.start).not.toHaveBeenCalled(); expect(view.prepare).not.toHaveBeenCalled();
});


it('queries the original O1 in readonly mode without permitting any mutation', async () => {
  const view = host(); await waitFor(() => expect(view.result.current.canRetrySameRequest).toBe(true));
  view.rerender({ selected: scope, readonly: true }); vi.mocked(api.operation).mockClear();
  await act(async () => { await view.result.current.check(); });
  expect(api.operation).toHaveBeenCalledExactlyOnceWith('O1');
  expect(api.start).not.toHaveBeenCalled(); expect(view.prepare).not.toHaveBeenCalled();
});
it('denies retry POST when the connection goes offline during the fresh original O1 lookup', async () => {
  const view = host(); await waitFor(() => expect(view.result.current.canRetrySameRequest).toBe(true));
  const pending = deferred(); vi.mocked(api.operation).mockReturnValueOnce(pending.promise);
  let running!: Promise<void>; act(() => { running = view.result.current.start(true); });
  await waitFor(() => expect(api.operation).toHaveBeenCalledTimes(2));
  vi.spyOn(navigator, 'onLine', 'get').mockReturnValue(false); fireEvent(window, new Event('offline'));
  await act(async () => { pending.resolve(partial); await running; });
  expect(api.start).not.toHaveBeenCalled(); expect(view.result.current.canRetrySameRequest).toBe(false);
});
