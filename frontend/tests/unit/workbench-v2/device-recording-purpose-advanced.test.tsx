import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { WriteGroupSection } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection';
import { studioV2WorkspaceMeasurementsAPI } from '@/services/studioV2WorkspaceMeasurements';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2RecordingStartAPI } from '@/services/studioV2RecordingStart';
import type { WriteGroup, WriteGroupDraft } from '@/types/studioV2WriteGroup';
import { batchPanelState } from '../../fixtures/circuitBatchPanel';
import { savedGroup, TABLE_COLUMNS } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string, options?: unknown) => (options && typeof options === 'object' && 'count' in options ? `${key}:${(options as { count: number }).count}` : key) }),
}));
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({
  useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }),
}));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: {
    list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), ensureCircuitProfile: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(),
    remove: vi.fn(), testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), get: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(),
  },
}));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { list: vi.fn(async () => []) } }));
vi.mock('@/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), get: vi.fn(), operation: vi.fn() } }));

const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const measurementApi = vi.mocked(studioV2WorkspaceMeasurementsAPI);
const startApi = vi.mocked(studioV2RecordingStartAPI);

/** A wide profile group of one circuit as the server saves it: an explicit mean, a role for every source, no storage mode or recreate policy. */
const wideGroup = (over: Partial<WriteGroup> = {}): WriteGroup => savedGroup({
  id: 'grp-w', revision: 'rev-grp-w-1', name: 'Feeder W (instantaneous)', status: 'ready', layout_version: 'managed_wide_profile_v1',
  recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: 'mean' },
  destination: { ...savedGroup().destination, storage_strategy: 'managed', table_name: 'gw_wide_fixture' },
  members: [0, 1].map((i) => ({
    device_id: 'dev-1', point_id: `pt-${i}`, tag_id: `tag-${i}`, source_revision: `s${i}`, mapping_revision: `m${i}`, target_column: `col_${i}`, required: false,
    circuit_id: 'cir-w', circuit_name: 'Feeder W', value_role: 'gauge' as const, source_value_type: 'float64' as const, result_value_type: 'decimal' as const,
  })),
  row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
  write_policy: { mode: 'append', dedupe_capability: 'receipt' },
  ...over,
});
const list = (groups: WriteGroup[]) => ({ workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups });
const emptyDelivery = {
  group_id: 'x', intake: { state: 'not_running' }, stages: { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 },
  last_sql_committed_at: null, oldest_pending_seconds: 0, no_data_buckets: 0, skipped_buckets: 0, backlog: [],
  quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
};
const client = () => new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });

beforeEach(() => {
  vi.clearAllMocks();
  window.sessionStorage.clear();
  measurementApi.list.mockResolvedValue([]);
  api.delivery.mockResolvedValue(emptyDelivery as never);
  api.readiness.mockResolvedValue({ workspace_id: 'ws-1', group_id: 'grp-w', group_revision: 'rev-grp-w-1', applied_revision: '', config_ready: true, schema_ready: true, ready: true, issues: [] } as never);
});

function renderHost(groups: WriteGroup[]) {
  api.list.mockResolvedValue(list(groups));
  return render(
    <QueryClientProvider client={client()}>
      <Step4Database state={batchPanelState()} dispatch={vi.fn()} workspaceId="ws-1" />
    </QueryClientProvider>,
  );
}

describe('Advanced opens the same saved wide profile group', () => {
  it('shows its layout, explicit mean, period, roles and saved-versus-applied state, with no competing control', async () => {
    renderHost([wideGroup()]);
    fireEvent.click(await screen.findByTestId('circuit-status-open-advanced-instantaneous'));
    expect(await screen.findByTestId('group-editor')).toBeInTheDocument();

    expect(screen.getByTestId('group-name')).toHaveValue('Feeder W (instantaneous)');
    const summary = within(screen.getByTestId('group-profile-summary'));
    expect(summary.getByTestId('group-profile-layout')).toHaveTextContent('step4.group.profile.layout.managed_wide_profile_v1');
    expect(summary.getByTestId('group-profile-reduction')).toHaveTextContent('step4.group.profile.reduction.mean');
    expect(summary.getByTestId('group-profile-interval')).toHaveTextContent('15');
    expect(summary.getByTestId('group-profile-interval')).toHaveTextContent('step4.group.profile.unit.second');
    expect(screen.getByTestId('group-profile-members')).toHaveTextContent('step4.group.profile.role.gauge');
    expect(screen.getByTestId('group-applied-revision')).toBeInTheDocument();
    for (const id of ['group-interval', 'group-incomplete', 'group-dedupe', 'group-entity-column', 'group-provenance-column', 'group-column-proposal']) {
      expect(screen.queryByTestId(id), id).not.toBeInTheDocument();
    }
    expect(screen.getByTestId('group-storage-managed')).toBeDisabled();
    expect(screen.getByTestId('group-table')).toBeDisabled();
  });

  it('returns to Basic with the same locked wide purpose and nothing created, updated, applied or started', async () => {
    renderHost([wideGroup()]);
    fireEvent.click(await screen.findByTestId('circuit-status-open-advanced-instantaneous'));
    await screen.findByTestId('group-editor');
    fireEvent.click(screen.getByTestId('group-editor-close'));

    expect(screen.getByTestId('circuit-profile-purpose-wide')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-purpose-long')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-reduction-instantaneous')).toHaveTextContent('step4.basic.circuit.reduction_mean:15');
    for (const call of [api.create, api.ensureCircuitProfile, api.update, api.apply, api.schemaApply, startApi.start]) {
      expect(call).not.toHaveBeenCalled();
    }
  });

  it('saves a rename with the layout, mean, roles and period unchanged', async () => {
    api.update.mockImplementation(async (id, request) => ({ workspace_revision: 'wrev-2', group: { ...wideGroup(), ...(request.group as object), id, revision: 'rev-grp-w-2' } as WriteGroup }));
    api.list.mockResolvedValue(list([wideGroup()]));
    render(
      <QueryClientProvider client={client()}>
        <WriteGroupSection state={batchPanelState()} workspaceId="ws-1" readonly={false} openGroupRequest={{ id: 'grp-w' }} />
      </QueryClientProvider>,
    );
    await screen.findByTestId('group-profile-summary');
    fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Feeder W renamed' } });
    fireEvent.click(screen.getByTestId('group-save'));
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));

    const [id, request] = api.update.mock.calls[0];
    const sent = request.group as WriteGroupDraft;
    expect(id).toBe('grp-w');
    expect(request.expected_group_revision).toBe('rev-grp-w-1');
    expect(sent).toMatchObject({
      name: 'Feeder W renamed', layout_version: 'managed_wide_profile_v1',
      recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: 'mean' }, row_policy: { interval_seconds: 15 },
    });
    expect(sent.members.map((member) => member.value_role)).toEqual(['gauge', 'gauge']);
    expect(sent.storage_mode).toBeUndefined();
    expect(sent.write_policy.recreate_missing).toBeUndefined();
  });
});
