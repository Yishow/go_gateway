import { retainSchemaOperation } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery';
import { groupToDraft, draftToRequestGroup } from '@/features/datalink/workbench-v2/state/writeGroup/draft';
import { mixedDraftToRequestGroup } from '@/features/datalink/workbench-v2/recording-editor/recordingMixedDraft';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { GroupEditor, type GroupEditorProps } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';
import type { WriteGroupReadiness } from '@/types/studioV2WriteGroup';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({ useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }) }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  readiness: vi.fn(), delivery: vi.fn().mockResolvedValue(null), list: vi.fn(), get: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(), testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(),
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
      return { ...member, source_revision: `fresh-${original.source_revision}`, mapping_revision: `fresh-${original.mapping_revision}` };
    });
    vi.mocked(studioV2WorkspaceWriteGroupsAPI.readiness).mockResolvedValue({ ...stale(), workspace_revision: 'db-after', group_revision: 'rev-after', config_ready: true,
      issues: [{ code: 'schema-confirmation-required', severity: 'blocking', step: 'Step 4', scope: 'group-1', message: 'explicit schema confirmation is required' }] });
    return { workspace_revision: 'db-after', group: savedGroup({ ...request.group, members, revision: 'rev-after', destination: { ...persisted.destination, ...request.group.destination, schema_revision: '', schema_digest: '' } }) };
  });
});
afterEach(() => { retainSchemaOperation('ws-1', 'group-1', undefined, 'pending-DDL'); cleanup(); clients.splice(0).forEach(client => client.clear()); });

const capabilities = { version: 1 as const, result_shape: 'member-calculation-v1' as const, editor_supported: true, draft_supported: true, activation_supported: true, reason: 'supported',
  methods: ['latest', 'mean', 'delta'] as ('latest' | 'mean' | 'delta')[], cadence: { alignment: 'utc_epoch' as const, units: ['second', 'minute'] as ('second' | 'minute')[], minimum: 1, max_seconds: 600, max_minutes: 10, window_bounds_checked: true as const },
  schema_dialects: ['sqlite'] as 'sqlite'[], schema_layouts: ['managed_wide_profile_v1', 'managed_long_v1'] as ('managed_wide_profile_v1' | 'managed_long_v1')[], unsupported_combinations: [] };
function mixedGroup() {
  const base = savedGroup();
  return savedGroup({ record_id: '11111111-1111-4111-8111-111111111111', layout_version: 'managed_wide_profile_v1',
    recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' }, destination: { ...base.destination, storage_strategy: 'managed' },
    row_policy: { allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
    members: base.members.map(member => ({ ...member, calculation: 'latest', cadence: { interval: 15, unit: 'second', alignment: 'utc_epoch' }, source_value_type: 'float64' })) });
}
it('recovers a pinned mixed recording only by explicit current-source consent through its original update', async () => {
  const group = mixedGroup(), initial = JSON.stringify(group); const view = mount({ group, calculationCapabilities: capabilities });
  fireEvent.click(await screen.findByTestId('recording-source-refresh-review'));
  expect(screen.getByTestId('recording-source-refresh-confirmation')).toHaveTextContent('sourceRefresh.currentAtSave');
  fireEvent.click(screen.getByTestId('recording-source-refresh-cancel')); expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
  fireEvent.click(screen.getByTestId('recording-source-refresh-review')); fireEvent.click(screen.getByTestId('recording-source-refresh-confirm'));
  await waitFor(() => expect(view.props.onSaved).toHaveBeenCalledOnce());
  const [id, request] = vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mock.calls[0];
  expect(id).toBe(group.id); expect(request.expected_group_revision).toBe(group.revision);
  const expected = mixedDraftToRequestGroup(groupToDraft(group), 'ws-1', capabilities);
  expected.members = expected.members.map(member => { const copy = { ...member }; delete copy.source_revision; delete copy.mapping_revision; return copy; });
  expect(request.group).toEqual(expected);
  expect(view.props.onSaved).toHaveBeenCalledWith(expect.objectContaining({ id: group.id, record_id: group.record_id,
    destination: expect.objectContaining({ table_name: group.destination.table_name, schema_revision: '', schema_digest: '' }),
    members: expect.arrayContaining([expect.objectContaining({ source_revision: 'fresh-s0', mapping_revision: 'fresh-m0' })]),
  }), 'db-after', expect.anything());
  expect(studioV2WorkspaceWriteGroupsAPI.schemaApply).not.toHaveBeenCalled();
  expect(JSON.stringify(group)).toBe(initial); expect(studioV2WorkspaceWriteGroupsAPI.apply).not.toHaveBeenCalled(); expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
});
it('keeps exact source pins on an ordinary mixed Save rather than silently using the repair mode', async () => {
  const group = mixedGroup(); mount({ group, calculationCapabilities: capabilities });
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Only rename the record' } });
  await waitFor(() => expect(screen.getByTestId('group-save')).toBeEnabled()); fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledOnce());
  const request = vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mock.calls[0][1];
  expect(request.group.members.map(member => [member.source_revision, member.mapping_revision])).toEqual(group.members.map(member => [member.source_revision, member.mapping_revision]));
});
it('makes the original nonmixed profile save recovery reachable without changing its profile identity', async () => {
  const group = savedGroup({ recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: 'mean' }, layout_version: 'managed_wide_profile_v1' });
  mount({ group }); fireEvent.click(await screen.findByTestId('recording-source-refresh-review')); fireEvent.click(screen.getByTestId('recording-source-refresh-confirm'));
  await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledOnce());
  expect(vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mock.calls[0][1].group).toEqual(draftToRequestGroup(groupToDraft(group), 'ws-1'));
});
it('does not submit a mixed source review after the scope changes', async () => {
  const group = mixedGroup(); const view = mount({ group, calculationCapabilities: capabilities });
  fireEvent.click(await screen.findByTestId('recording-source-refresh-review'));
  view.update({ workspaceRevision: 'another-revision' });
  const confirm = screen.queryByTestId('recording-source-refresh-confirm'); if (confirm) fireEvent.click(confirm);
  expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled(); expect(group.members[0].calculation).toBe('latest');
});
it('keeps an unknown mixed save blocked rather than retrying or declaring the source repaired', async () => {
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mockRejectedValue(new Error('lost acknowledgement'));
  mount({ group: mixedGroup(), calculationCapabilities: capabilities });
  fireEvent.click(await screen.findByTestId('recording-source-refresh-review')); fireEvent.click(screen.getByTestId('recording-source-refresh-confirm'));
  await screen.findByTestId('group-save-unknown'); expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledOnce();
  expect(screen.queryByTestId('recording-source-refresh-review')).not.toBeInTheDocument(); expect(studioV2WorkspaceWriteGroupsAPI.apply).not.toHaveBeenCalled();
});
