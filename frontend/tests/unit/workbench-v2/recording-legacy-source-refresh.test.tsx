import { retainSchemaOperation } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { GroupEditor, type GroupEditorProps } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { GroupReadinessPanel } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupReadinessPanel';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';
import type { WriteGroupReadiness } from '@/types/studioV2WriteGroup';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({ useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }) }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  readiness: vi.fn(), delivery: vi.fn().mockResolvedValue(null), list: vi.fn(), get: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(), testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(),
} }));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { list: vi.fn().mockResolvedValue([]) } }));
const clients: QueryClient[] = [];
function stale(): WriteGroupReadiness {
  return { workspace_id: 'ws-1', workspace_revision: 'db-setup-9', group_id: 'group-1', group_revision: 'rev-1', applied_revision: '', config_ready: false, schema_ready: false, ready: false,
    issues: [
      { code: 'source-revision-stale', severity: 'blocking', step: 'Step 4', scope: 'group-1', message: 'saved source or mapping revision is stale' },
      { code: 'schema-check-blocked', severity: 'blocking', step: 'Step 4', scope: 'group-1', message: 'destination schema check is blocked by configuration issues' },
    ] };
}
function mount(overrides: Partial<GroupEditorProps> = {}) {
  const state = savedGroupState();
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 }, mutations: { retry: false } } }); clients.push(client);
  const props: GroupEditorProps = { presentation: 'simple', group: savedGroup(), workspaceId: 'ws-1', workspaceRevision: 'db-setup-9', connector: state.db.connector,
    candidates: buildGroupCandidates(state).candidates, excluded: [], readonly: false, onSaved: vi.fn(), onClose: vi.fn(), onReload: vi.fn(),
    fetchLatest: vi.fn().mockResolvedValue({ workspace_revision: 'db-setup-9', groups: [savedGroup()] }), ...overrides };
  const tree = () => <QueryClientProvider client={client}><GroupEditor {...props} /></QueryClientProvider>;
  const view = render(tree());
  return { props, client, update(patch: Partial<GroupEditorProps>) { Object.assign(props, patch); view.rerender(tree()); } };
}
beforeEach(() => {
  vi.clearAllMocks(); vi.mocked(studioV2WorkspaceWriteGroupsAPI.readiness).mockResolvedValue(stale());
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mockImplementation(async (_id, request) => {
    const persisted = savedGroup();
    const members = request.group.members.map(member => {
      const original = persisted.members.find(value => value.device_id === member.device_id && value.point_id === member.point_id && value.tag_id === member.tag_id);
      if (!original) throw new Error('Fixture source identity is unavailable');
      return { ...member, source_revision: original.source_revision, mapping_revision: original.mapping_revision };
    });
    return { workspace_revision: 'db-after', group: savedGroup({ ...request.group, members, revision: 'rev-after', destination: persisted.destination }) };
  });
});
afterEach(() => { retainSchemaOperation('ws-1', 'group-1', undefined, 'pending-DDL'); cleanup(); clients.splice(0).forEach(client => client.clear()); });
it('explains one stale source cause while retaining its derived schema issue in diagnostics', () => {
  const value = stale();
  render(<GroupReadinessPanel compact readiness={value} loading={false} failed={false} dirty={false} onRetry={vi.fn()} />);
  const issues = screen.getByTestId('group-readiness-issues');
  expect(issues.children).toHaveLength(1);
  expect(within(issues).getByText('confirmation.sourceStale', { exact: false })).toBeVisible();
  const original = within(issues).getByText('schema-check-blocked: destination schema check is blocked by configuration issues');
  expect(original).not.toBeVisible(); fireEvent.click(original.closest('details')!.querySelector('summary')!); expect(original).toBeVisible();
  expect(value.issues).toHaveLength(2); expect(value.ready).toBe(false);
});
it('offers explicit source review on an unchanged saved legacy recording without requiring a fake name edit', async () => {
  const view = mount();
  fireEvent.click(await screen.findByTestId('recording-source-refresh-review'));
  const review = screen.getByTestId('recording-source-refresh-confirmation');
  expect(within(review).getByText('sourceRefresh.currentAtSave')).toBeVisible();
  expect(within(review).getAllByTestId('recording-source-refresh-member')).toHaveLength(2);
  expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
  fireEvent.click(within(review).getByTestId('recording-source-refresh-cancel'));
  expect(screen.queryByTestId('recording-source-refresh-confirmation')).not.toBeInTheDocument();
  fireEvent.click(screen.getByTestId('recording-source-refresh-review'));
  fireEvent.click(screen.getByTestId('recording-source-refresh-confirm'));
  await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledTimes(1));
  const [id, request] = vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mock.calls[0];
  expect(id).toBe('group-1'); expect(request).toMatchObject({ workspace_id: 'ws-1', expected_workspace_revision: 'db-setup-9', expected_group_revision: 'rev-1', expected_connector_revision: 'crev-1' });
  expect(request.group.name).toBe('Line A'); expect(request.group.members.map(m => [m.device_id, m.point_id, m.tag_id, m.target_column])).toEqual(savedGroup().members.map(m => [m.device_id, m.point_id, m.tag_id, m.target_column]));
  request.group.members.forEach(member => { expect(member).not.toHaveProperty('source_revision'); expect(member).not.toHaveProperty('mapping_revision'); });
  await waitFor(() => expect(view.props.onSaved).toHaveBeenCalledTimes(1)); expect(studioV2WorkspaceWriteGroupsAPI.apply).not.toHaveBeenCalled(); expect(studioV2WorkspaceWriteGroupsAPI.testWrite).not.toHaveBeenCalled();
});
it('retires the source review when the operator edits after opening it', async () => {
  mount(); fireEvent.click(await screen.findByTestId('recording-source-refresh-review'));
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Do not overwrite this edit' } });
  const confirm = screen.queryByTestId('recording-source-refresh-confirm'); if (confirm) fireEvent.click(confirm);
  expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled(); expect(screen.getByTestId('group-name')).toHaveValue('Do not overwrite this edit');
});
it('does not turn foreign readiness into permission to refresh saved source bindings', async () => {
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.readiness).mockResolvedValue({ ...stale(), workspace_revision: 'foreign' });
  mount(); await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.readiness).toHaveBeenCalled());
  expect(screen.queryByTestId('recording-source-refresh-review')).not.toBeInTheDocument(); expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
});
it('keeps a source refresh confirmation inert after its original editor becomes readonly', async () => {
  const view = mount(); fireEvent.click(await screen.findByTestId('recording-source-refresh-review')); view.update({ readonly: true });
  const confirm = screen.queryByTestId('recording-source-refresh-confirm'); if (confirm) fireEvent.click(confirm);
  expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
});

it('names unavailable member positions and keeps source correction reachable instead of silently dropping members', async () => {
  mount({ candidates: [] }); await screen.findByTestId('recording-source-refresh');
  expect(screen.queryByTestId('recording-source-refresh-review')).not.toBeInTheDocument();
  expect(screen.getAllByText('sourceRefresh.missingMember')).toHaveLength(2);
  expect(screen.getByRole('link', { name: 'sourceRefresh.inspect' })).toHaveAttribute('href', '#recording-editor-sources');
  expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
});
it('rechecks source review after the original save owner waits for its current workspace read', async () => {
  let resolve!: (value: { workspace_revision: string; groups: ReturnType<typeof savedGroup>[] }) => void;
  const latest = new Promise<{ workspace_revision: string; groups: ReturnType<typeof savedGroup>[] }>(done => { resolve = done; });
  const view = mount({ fetchLatest: vi.fn().mockReturnValue(latest) });
  fireEvent.click(await screen.findByTestId('recording-source-refresh-review')); fireEvent.click(screen.getByTestId('recording-source-refresh-confirm'));
  await waitFor(() => expect(view.props.fetchLatest).toHaveBeenCalledOnce());
  view.client.setQueryData(studioV2WorkspaceKeys.writeGroupReadiness('group-1'), { ...stale(), workspace_revision: 'changed' });
  resolve({ workspace_revision: 'db-setup-9', groups: [savedGroup()] });
  await waitFor(() => expect(screen.queryByTestId('recording-source-refresh-confirm')).not.toBeInTheDocument());
  expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
});
it('does not repeat a source save when the original save acknowledgement is unknown', async () => {
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mockRejectedValue(new Error('connection lost'));
  mount(); fireEvent.click(await screen.findByTestId('recording-source-refresh-review')); fireEvent.click(screen.getByTestId('recording-source-refresh-confirm'));
  await screen.findByTestId('group-save-unknown');
  expect(screen.queryByTestId('recording-source-refresh-review')).not.toBeInTheDocument();
  expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledTimes(1); expect(studioV2WorkspaceWriteGroupsAPI.apply).not.toHaveBeenCalled();
});
it('does not expose source refresh for an invalid mixed recording', async () => {
  const group = savedGroup({ record_id: 'e94e7348-c90c-4d24-bba0-680c3d89859c', recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' }, layout_version: 'managed_wide_profile_v1' });
  const before = JSON.stringify(group); mount({ group, calculationCapabilities: {
    version: 1, result_shape: 'member-calculation-v1', editor_supported: true, draft_supported: true, activation_supported: false,
    reason: 'record_readiness_required', methods: ['latest', 'mean', 'delta'], cadence: { alignment: 'utc_epoch', units: ['second', 'minute'], minimum: 1, max_seconds: 600, max_minutes: 10, window_bounds_checked: true },
    schema_dialects: ['sqlite'], schema_layouts: ['managed_wide_profile_v1', 'managed_long_v1'], unsupported_combinations: [],
  } });
  await screen.findByTestId('recording-source-refresh');
  expect(screen.getByText('sourceRefresh.unavailable')).toBeVisible();
  expect(screen.queryByTestId('recording-source-refresh-review')).not.toBeInTheDocument();
  expect(JSON.stringify(group)).toBe(before); expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
});

it('rejects a confirmation whose viewed workspace changed while the original save is awaiting fresh metadata', async () => {
  let resolve!: (value: { workspace_revision: string; groups: ReturnType<typeof savedGroup>[] }) => void;
  const latest = new Promise<{ workspace_revision: string; groups: ReturnType<typeof savedGroup>[] }>(done => { resolve = done; });
  const view = mount({ fetchLatest: vi.fn().mockReturnValue(latest) });
  fireEvent.click(await screen.findByTestId('recording-source-refresh-review')); fireEvent.click(screen.getByTestId('recording-source-refresh-confirm'));
  await waitFor(() => expect(view.props.fetchLatest).toHaveBeenCalledOnce());
  view.update({ workspaceId: 'other-workspace' }); resolve({ workspace_revision: 'db-setup-9', groups: [savedGroup()] });
  await waitFor(() => expect(screen.queryByTestId('recording-source-refresh-confirm')).not.toBeInTheDocument());
  expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
});

it('does not start a source refresh while an original schema operation remains unresolved', async () => {
  retainSchemaOperation('ws-1', 'group-1', 'pending-DDL'); mount();
  await screen.findByTestId('recording-source-refresh');
  expect(screen.queryByTestId('recording-source-refresh-review')).not.toBeInTheDocument();
  expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
});
