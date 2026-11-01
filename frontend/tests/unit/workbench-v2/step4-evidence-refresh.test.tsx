import { useEffect, type ReactNode } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, renderHook, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import type { BasicRecordingPanelProps } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel';
import type { RecordingScopeSummaryProps } from '@/features/datalink/workbench-v2/steps/step4/recordingScopeModel';
import { useRuntimeRuleEvidence } from '@/hooks/datalink/useRuntimeRuleEvidence';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import type { WriteGroupListResponse } from '@/types/studioV2WriteGroup';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
let observedEvidence: RecordingScopeSummaryProps['evidence'] | undefined;
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel', () => ({ BasicRecordingPanel:
  ({ onViewedDeviceChange, renderScopeSummary }: BasicRecordingPanelProps) => {
    useEffect(() => { onViewedDeviceChange?.({ workspaceId: 'ws-1', deviceId: 'dev-1' }); }, [onViewedDeviceChange]);
    return renderScopeSummary?.({ workspaceId: 'ws-1', workspaceRevision: 'R', deviceId: 'dev-1', groupIds: ['group-1'],
      draft: { profiles: [], layout: 'wide', storageMode: 'plain', incompletePolicy: 'partial' } });
  } }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/RecordingScopeSummary', () => ({ RecordingScopeSummary:
  ({ evidence }: RecordingScopeSummaryProps) => { observedEvidence = evidence; return <><button disabled={evidence.isFetching} onClick={() => { void evidence.refetch(); }}>retry evidence</button>
    <output data-testid="evidence-state">{JSON.stringify({ data: evidence.data, failed: evidence.isError })}</output></>; } }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/Step4SupportPanels', () => ({ Step4SupportPanels: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/RecordingTestWriteRecoveryPanel', () => ({ RecordingTestWriteRecoveryPanel: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection', () => ({ WriteGroupSection: () => null }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn() } }));
vi.mock('@/services/studioV2RuntimeContext', () => ({ studioV2RuntimeContextAPI: { getRuleEvidence: vi.fn() } }));
const list = vi.mocked(studioV2WorkspaceWriteGroupsAPI.list);
const evidence = vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence);
const unavailable = { ok: false as const, reason: 'unavailable' as const, action: 'refresh' as const };
const groups = (revision = 'R', workspace = 'ws-1'): WriteGroupListResponse => ({ workspace_id: workspace, workspace_revision: revision, groups: [savedGroup({ workspace_id: workspace })] });
const clients: QueryClient[] = [];
function client() { const value = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(value); return value; }
function deferred() { let resolve!: (value: WriteGroupListResponse) => void; const promise = new Promise<WriteGroupListResponse>(r => { resolve = r; }); return { promise, resolve }; }
async function mount() {
  const cache = client(); cache.setQueryData(studioV2WorkspaceKeys.writeGroups(), groups());
  const state = savedGroupState(); const dispatch = vi.fn();
  const element = (dirty = false, workspaceId = 'ws-1') => <QueryClientProvider client={cache}><Step4Database state={dirty ? { ...state, devices: state.devices.map(device => ({ ...device, save_state: 'saving' as const })) } : state} dispatch={dispatch} workspaceId={workspaceId} /></QueryClientProvider>;
  const view = render(element());
  await waitFor(() => expect(evidence).toHaveBeenCalled());
  await waitFor(() => expect(screen.getByRole('button', { name: 'retry evidence' })).not.toBeDisabled());
  list.mockClear(); evidence.mockClear(); dispatch.mockClear();
  return { ...view, cache, state, dispatch, element };
}
beforeEach(() => { list.mockReset().mockResolvedValue(groups()); evidence.mockReset().mockResolvedValue(unavailable); });
afterEach(() => { cleanup(); clients.splice(0).forEach(value => value.clear()); vi.restoreAllMocks(); });
it('explicit Step4 retry first refreshes its groups authority and reads only R2 while preserving editable state', async () => {
  const view = await mount(); const before = JSON.stringify(view.state); const pending = deferred(); list.mockReturnValueOnce(pending.promise);
  fireEvent.click(screen.getByRole('button', { name: 'retry evidence' }));
  await waitFor(() => expect(list).toHaveBeenCalledTimes(1)); expect(evidence).not.toHaveBeenCalled();
  await act(async () => pending.resolve(groups('R2')));
  await waitFor(() => expect(evidence).toHaveBeenCalledTimes(1));
  expect(evidence.mock.calls[0][0]).toEqual({ workspace_id: 'ws-1', workspace_revision: 'R2', device_id: 'dev-1' });
  expect(view.dispatch).not.toHaveBeenCalled(); expect(JSON.stringify(view.state)).toBe(before);
});
it.each(['failed', 'foreign'] as const)('canonical %s refresh fails closed without an old-head retry', async kind => {
  await mount(); if (kind === 'failed') list.mockRejectedValueOnce(new Error('offline')); else list.mockResolvedValueOnce(groups('R2', 'OTHER'));
  fireEvent.click(screen.getByRole('button', { name: 'retry evidence' }));
  await waitFor(() => expect(list).toHaveBeenCalledTimes(1));
  await waitFor(() => expect(screen.getByRole('button', { name: 'retry evidence' })).not.toBeDisabled());
  expect(evidence).not.toHaveBeenCalled();
});
it.each(['dirty', 'workspace', 'unmount'] as const)('does not resume an old scope after pending refresh and %s', async change => {
  const view = await mount(); const pending = deferred(); list.mockReturnValueOnce(pending.promise);
  fireEvent.click(screen.getByRole('button', { name: 'retry evidence' }));
  await waitFor(() => expect(list).toHaveBeenCalledTimes(1));
  if (change === 'unmount') view.unmount(); else view.rerender(view.element(change === 'dirty', change === 'workspace' ? 'OTHER' : 'ws-1'));
  await act(async () => pending.resolve(groups('R2')));
  expect(evidence).not.toHaveBeenCalled();
});
it('the shared default caller still invalidates only runtimeContext on stale evidence', async () => {
  const cache = client(); const invalidate = vi.spyOn(cache, 'invalidateQueries'); evidence.mockResolvedValueOnce({ ok: false, reason: 'stale-reference', action: 'refresh' });
  const wrapper = ({ children }: { children: ReactNode }) => <QueryClientProvider client={cache}>{children}</QueryClientProvider>;
  renderHook(() => useRuntimeRuleEvidence({ workspace_id: 'ws-1', workspace_revision: 'R', device_id: 'dev-1' }), { wrapper });
  await waitFor(() => expect(invalidate).toHaveBeenCalledWith({ queryKey: studioV2WorkspaceKeys.runtimeContext(), exact: true }));
  expect(invalidate).not.toHaveBeenCalledWith({ queryKey: studioV2WorkspaceKeys.writeGroups(), exact: true });
});

it('Step4 stale evidence refreshes only its exact groups authority, without an invalidation loop', async () => {
  const invalidate = vi.spyOn(QueryClient.prototype, 'invalidateQueries');
  evidence.mockResolvedValueOnce({ ok: false, reason: 'stale-reference', action: 'refresh' });
  await mount();
  expect(invalidate).toHaveBeenCalledWith({ queryKey: studioV2WorkspaceKeys.writeGroups(), exact: true });
  expect(invalidate).not.toHaveBeenCalledWith({ queryKey: studioV2WorkspaceKeys.runtimeContext(), exact: true });
  expect(invalidate).toHaveBeenCalledTimes(1);
});
it('explicit same-head retry uses freshly confirmed R rather than making up a new head', async () => {
  await mount();
  fireEvent.click(screen.getByRole('button', { name: 'retry evidence' }));
  await waitFor(() => expect(list).toHaveBeenCalledTimes(1));
  await waitFor(() => expect(evidence).toHaveBeenCalledTimes(1));
  expect(evidence.mock.calls[0][0]).toEqual({ workspace_id: 'ws-1', workspace_revision: 'R', device_id: 'dev-1' });
});

it('a newer explicit refresh supersedes the older in-flight head receipt', async () => {
  const view = await mount(); const older = deferred();
  list.mockReturnValueOnce(older.promise).mockResolvedValueOnce(groups('R3'));
  act(() => { void observedEvidence?.refetch(); });
  await waitFor(() => expect(list).toHaveBeenCalledTimes(1));
  act(() => { void observedEvidence?.refetch(); });
  await waitFor(() => expect(list).toHaveBeenCalledTimes(2));
  await waitFor(() => expect(evidence).toHaveBeenCalledTimes(1));
  await act(async () => older.resolve(groups('R2')));
  expect(evidence).toHaveBeenCalledTimes(1);
  expect(evidence.mock.calls[0][0]).toEqual({ workspace_id: 'ws-1', workspace_revision: 'R3', device_id: 'dev-1' });
  expect(view.cache.getQueryData<WriteGroupListResponse>(studioV2WorkspaceKeys.writeGroups())?.workspace_revision).toBe('R3');
});
