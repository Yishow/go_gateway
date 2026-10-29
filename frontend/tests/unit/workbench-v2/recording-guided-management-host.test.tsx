import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { studioV2WorkspaceWriteGroupsAPI as api } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WriteGroupLifecycleAPI as lifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';
import { lifecycle } from '../../fixtures/writeGroupSuccessor';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  list: vi.fn(), delivery: vi.fn(), readiness: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(),
  ensureCircuitProfile: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(),
} }));
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: vi.fn() } }));
vi.mock('@/services/studioV2RuntimeContext', () => ({ studioV2RuntimeContextAPI: { getRuleEvidence: vi.fn() } }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel', () => ({ BasicRecordingPanel: (props: {
  reviewRequest?: { workspaceId: string; groupId: string };
}) => <div data-testid="selected-basic-scope" data-device="dev-1">{props.reviewRequest ? JSON.stringify(props.reviewRequest) : 'Device A'}</div> }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/Step4SupportPanels', () => ({ Step4SupportPanels: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/CommitSummary', () => ({ CommitSummary: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/ShareOutputSummary', () => ({ ShareOutputSummary: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/RecordingTestWriteRecoveryPanel', () => ({ RecordingTestWriteRecoveryPanel: () => null, RecoveryEntries: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/CutoverPanel', () => ({ CutoverPanel: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel', () => ({ GroupSchemaPanel: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({ useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }) }));
let groups: WriteGroup[];
const activate = vi.fn();
function managed(id: string, deviceId: string): WriteGroup {
  const base = savedGroup();
  return savedGroup({ id, name: id, layout_version: 'managed_wide_profile_v1', storage_mode: 'plain',
    recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: 'latest' },
    members: base.members.map(member => ({ ...member, device_id: deviceId, circuit_id: `circuit-${id}`, circuit_name: id, value_role: 'gauge' })),
    row_policy: { ...base.row_policy, interval_seconds: 15 }, destination: { ...base.destination, storage_strategy: 'managed' } });
}
beforeEach(() => {
  vi.clearAllMocks(); sessionStorage.clear();
  groups = [managed('managed-A', 'dev-1'), managed('managed-B', 'dev-2'), savedGroup({ id: 'custom-C', name: 'Custom C' })];
  vi.mocked(api.list).mockImplementation(async () => ({ workspace_id: 'ws-1', workspace_revision: 'ws-r1', groups: structuredClone(groups) }));
  vi.mocked(api.delivery).mockRejectedValue(new Error('delivery evidence unavailable'));
  vi.mocked(api.readiness).mockImplementation(async id => {
    const group = groups.find(value => value.id === id); if (!group) throw new Error('missing test scope');
    return { workspace_id: 'ws-1', workspace_revision: 'ws-r1', applied_revision: group.applied_revision, group_id: id, group_revision: group.revision, config_ready: true, schema_ready: true, ready: true, issues: [] };
  });
  vi.mocked(lifecycleAPI.lifecycle).mockImplementation(async id => {
    const group = groups.find(value => value.id === id); if (!group) throw new Error('missing test scope');
    const source = lifecycle();
    return { ...source, group_id: id, group_revision: group.revision, status: group.status, applied_revision: group.applied_revision,
      actions: source.actions.map(action => ({ ...action, allowed: group.status !== 'deleted' || action.action === 'recreate',
        reason_code: group.status === 'deleted' ? 'recording_removed' : '' })),
      history: { canonical: true, current_group_id: id, generations: [] } };
  });
});
afterEach(cleanup);
function show() {
  const state = savedGroupState(); state.devices.push({ ...state.devices[0], id: 'dev-2', name: 'Device B' });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><Step4Database state={state} dispatch={vi.fn()} workspaceId="ws-1" activateWorkspace={activate} /></QueryClientProvider>);
}
function noMutation() {
  for (const fn of [api.create, api.update, api.apply, api.disable, api.remove, api.ensureCircuitProfile, api.schemaPreview, api.schemaApply, activate]) expect(fn).not.toHaveBeenCalled();
}
it('labels Advanced as status, exposes exact saved-card management, and does not fetch a list/effective reader per card or activate', async () => {
  show(); const advanced = screen.getByTestId('step4-advanced-recording');
  expect(within(advanced).getByText('management.advancedTitle')).toBeInTheDocument();
  fireEvent.click(within(advanced).getByText('management.advancedTitle'));
  expect(within(advanced).getByTestId('recording-workspace-status-note')).toHaveTextContent('management.workspaceStatus');
  for (const id of ['managed-A', 'managed-B', 'custom-C']) {
    const card = await screen.findByTestId(`group-list-item-${id}`);
    await waitFor(() => expect(within(card).getByTestId('saved-recording-edit')).toBeEnabled());
    expect(within(card).getAllByTestId('saved-recording-actions')).toHaveLength(1);
    expect(within(card).queryByTestId('group-apply')).not.toBeInTheDocument();
  }
  expect(api.list).toHaveBeenCalledTimes(1); expect(studioV2RuntimeContextAPI.getRuleEvidence).not.toHaveBeenCalled();
  expect(screen.getByTestId('selected-basic-scope')).toHaveAttribute('data-device', 'dev-1'); noMutation();
});
it('routes the existing managed effective action to exact A review, not B or a workspace start', async () => {
  show(); fireEvent.click(screen.getByText('management.advancedTitle'));
  fireEvent.click(await screen.findByTestId('group-open-managed-A'));
  const action = await screen.findByTestId('group-apply');
  await waitFor(() => expect(action).toBeEnabled()); expect(action).toHaveTextContent('step4.group.lifecycle.review_start');
  fireEvent.click(action);
  expect(screen.getByTestId('selected-basic-scope')).toHaveTextContent('{"workspaceId":"ws-1","groupId":"managed-A"}');
  expect(screen.getByTestId('selected-basic-scope')).not.toHaveTextContent('managed-B'); noMutation();
});
it('opens custom saved settings via lifecycle Edit without converting identity, layout or destination', async () => {
  const before = structuredClone(groups[2]); show(); fireEvent.click(screen.getByText('management.advancedTitle'));
  const card = await screen.findByTestId('group-list-item-custom-C');
  const edit = await within(card).findByTestId('saved-recording-edit'); await waitFor(() => expect(edit).toBeEnabled()); fireEvent.click(edit);
  expect(await screen.findByTestId('group-name')).toHaveValue('Custom C');
  expect(screen.getByTestId('group-storage-custom')).toBeChecked(); expect(screen.getByTestId('group-table')).toHaveValue('readings');
  expect(screen.getByTestId('group-revisions')).toHaveTextContent(before.revision);
  fireEvent.click(screen.getByTestId('group-editor-close'));
  expect(await screen.findByTestId('group-list-item-custom-C')).toBeInTheDocument(); expect(groups[2]).toEqual(before); noMutation();
});
it('keeps removed custom history view-only without ordinary creation or conversion', async () => {
  groups = [savedGroup({ id: 'removed-C', status: 'deleted', name: 'Removed custom' })]; show(); fireEvent.click(screen.getByText('management.advancedTitle'));
  const removed = await screen.findByTestId('group-removed'); fireEvent.click(within(removed).getByText('step4.group.list.removed'));
  fireEvent.click(screen.getByTestId('group-open-removed-C'));
  expect(await screen.findByTestId('group-name')).toBeDisabled(); expect(screen.getByTestId('group-save')).toBeDisabled();
  expect(screen.getByTestId('group-storage-custom')).toBeChecked(); noMutation();
});
it.each(['foreign workspace', 'missing revision'] as const)('keeps card commands unavailable for a canonical list with %s', async invalid => {
  vi.mocked(api.list).mockImplementation(async () => ({ workspace_id: invalid === 'foreign workspace' ? 'other-workspace' : 'ws-1',
    workspace_revision: invalid === 'missing revision' ? '' : 'ws-r1', groups: structuredClone(groups) }));
  show(); const card = await screen.findByTestId('group-list-item-custom-C');
  await waitFor(() => expect(lifecycleAPI.lifecycle).toHaveBeenCalledWith('custom-C'));
  expect(within(card).getByTestId('group-disable')).toBeDisabled(); expect(within(card).getByTestId('group-delete')).toBeDisabled();
  expect(within(card).getByTestId('group-open-custom-C')).toBeEnabled(); noMutation();
});
