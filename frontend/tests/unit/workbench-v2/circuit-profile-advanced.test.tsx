import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import * as React from 'react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { WriteGroupSection } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection';
import { studioV2WorkspaceMeasurementsAPI } from '@/services/studioV2WorkspaceMeasurements';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { MeasurementDefinition } from '@/types/measurement';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { batchPanelState, longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string, options?: unknown) => (options && typeof options === 'object' && 'count' in options ? `${key}:${(options as { count: number }).count}` : key) }),
}));
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({
  useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }),
}));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: {
    list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(),
    testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), get: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(),
  },
}));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { list: vi.fn(async () => []) } }));
vi.mock('@/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), get: vi.fn() } }));

const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const measurementApi = vi.mocked(studioV2WorkspaceMeasurementsAPI);
const longGroup = (over: Partial<WriteGroup> = {}): WriteGroup => ({ ...longProfileGroup({ id: 'grp-a', circuitId: 'cir-a', circuitName: 'Feeder A', interval: 15, points: [0, 1] }), ...over });
const cumulativeGroup = (storage_mode: WriteGroup['storage_mode'] = 'plain'): WriteGroup => ({ ...longProfileGroup({ id: 'grp-c', circuitId: 'cir-a', circuitName: 'Feeder A', interval: 10, points: [0, 1], kind: 'cumulative' }), storage_mode });
const list = (groups: WriteGroup[]) => ({ workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups });
const emptyDelivery = {
  group_id: 'x', intake: { state: 'not_running' }, stages: { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 },
  last_sql_committed_at: null, oldest_pending_seconds: 0, no_data_buckets: 0, skipped_buckets: 0, backlog: [],
  quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
};
const client = () => new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });

function renderSection(groups: WriteGroup[], openId: string | undefined, props: Partial<React.ComponentProps<typeof WriteGroupSection>> = {}) {
  api.list.mockResolvedValue(list(groups));
  return render(
    <QueryClientProvider client={client()}>
      <WriteGroupSection state={batchPanelState()} workspaceId="ws-1" readonly={false} openGroupRequest={openId ? { id: openId } : undefined} {...props} />
    </QueryClientProvider>,
  );
}

beforeEach(() => {
  vi.clearAllMocks();
  measurementApi.list.mockResolvedValue([]);
  api.delivery.mockResolvedValue(emptyDelivery as never);
  api.readiness.mockResolvedValue({ workspace_id: 'ws-1', group_id: 'grp-a', group_revision: 'rev-grp-a-1', applied_revision: '', config_ready: true, schema_ready: true, ready: true, issues: [] } as never);
});

describe('Advanced shows the same canonical profile and offers no second period or completeness setting', () => {
  it('shows the profile kind, its native interval and unit, and the saved storage mode', async () => {
    renderSection([longGroup()], 'grp-a');
    const summary = await screen.findByTestId('group-profile-summary');
    expect(within(summary).getByTestId('group-profile-kind')).toHaveTextContent('step4.group.profile.kind.instantaneous');
    expect(within(summary).getByTestId('group-profile-interval')).toHaveTextContent('15');
    expect(within(summary).getByTestId('group-profile-interval')).toHaveTextContent('step4.group.profile.unit.second');
    expect(within(summary).getByTestId('group-profile-storage')).toHaveTextContent('step4.group.profile.storage.plain');
  });

  it('shows minutes and Timescale for a cumulative profile saved that way', async () => {
    renderSection([cumulativeGroup('timescale')], 'grp-c');
    const summary = await screen.findByTestId('group-profile-summary');
    expect(within(summary).getByTestId('group-profile-interval')).toHaveTextContent('10');
    expect(within(summary).getByTestId('group-profile-interval')).toHaveTextContent('step4.group.profile.unit.minute');
    expect(within(summary).getByTestId('group-profile-storage')).toHaveTextContent('step4.group.profile.storage.timescale');
  });

  it('offers the counter policy form only for a cumulative profile', async () => {
    const first = renderSection([cumulativeGroup()], 'grp-c');
    await screen.findByTestId('group-profile-summary');
    expect(screen.getByTestId('counter-policy')).toBeInTheDocument();
    first.unmount();
    renderSection([longGroup()], 'grp-a');
    await screen.findByTestId('group-profile-summary');
    expect(screen.queryByTestId('counter-policy')).not.toBeInTheDocument();
  });

  it('hands the persisted measurement reference of each member to the policy form', async () => {
    const definition = (id: string, modulus: string): MeasurementDefinition => ({
      id, workspace_id: 'ws-1', device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', equipment_id: 'eq-A', definition_revision: 'def-1', source_binding_revision: 'bind-1',
      series_epoch: 'epoch-1', name: id, quantity: 'energy', unit: 'kWh', semantic_kind: 'counter', numeric_encoding: 'uint64',
      counter_policy: { allow_negative: false, rollover_modulus: modulus, max_rate_per_second: '10' }, created_at: '2026-10-01T00:00:00Z', updated_at: '2026-10-01T00:00:00Z',
    });
    measurementApi.list.mockResolvedValue([definition('meas-other', '65536'), definition('meas-bound', '131072')]);
    const group = cumulativeGroup();
    renderSection([{ ...group, members: group.members.map((member, index) => (index === 0 ? { ...member, measurement_id: 'meas-bound' } : member)) }], 'grp-c');
    expect(await screen.findByTestId('counter-policy-modulus-tag-0')).toHaveValue('131072');
    expect(screen.getByTestId('counter-policy-missing-tag-1')).toBeInTheDocument();
  });

  it('has no seconds, completeness, dedupe or wide column controls that could compete with the profile', async () => {
    renderSection([longGroup()], 'grp-a');
    await screen.findByTestId('group-profile-summary');
    for (const id of ['group-interval', 'group-incomplete', 'group-dedupe', 'group-entity-column', 'group-provenance-column', 'group-column-proposal', 'group-member-metadata-note']) {
      expect(screen.queryByTestId(id), id).not.toBeInTheDocument();
    }
    expect(screen.queryAllByTestId(/^group-member-column-/)).toHaveLength(0);
    expect(screen.queryAllByLabelText(/step4\.group\.members\.required_label/)).toHaveLength(0);
    expect(screen.getByTestId('group-profile-members')).toHaveTextContent('sensor.0');
  });

  it('does not offer to change the layout, storage kind or shared table of a profile group', async () => {
    renderSection([longGroup()], 'grp-a');
    await screen.findByTestId('group-profile-summary');
    expect(screen.getByTestId('group-storage-managed')).toBeDisabled();
    expect(screen.getByTestId('group-storage-custom')).toBeDisabled();
    expect(screen.getByTestId('group-table')).toBeDisabled();
    expect(screen.getByTestId('group-schema')).toBeDisabled();
  });

  it('keeps the legacy wide group exactly as it was, with seconds, completeness and per member required', async () => {
    renderSection([savedGroup()], 'group-1', { state: savedGroupState() });
    await screen.findByTestId('group-editor');
    expect(screen.queryByTestId('group-profile-summary')).not.toBeInTheDocument();
    expect(screen.getByTestId('group-interval')).toBeEnabled();
    expect(screen.getByTestId('group-incomplete')).toBeEnabled();
    expect(screen.getAllByLabelText(/step4\.group\.members\.required_label/).length).toBeGreaterThan(0);
    expect(screen.getAllByTestId(/^group-member-column-/).length).toBeGreaterThan(0);
  });

  it('saves a rename through the same update with the unchanged profile and derived seconds', async () => {
    api.update.mockImplementation(async (id, request) => ({ workspace_revision: 'wrev-2', group: { ...longGroup(), ...(request.group as object), id, revision: 'rev-grp-a-2' } as WriteGroup }));
    renderSection([longGroup()], 'grp-a');
    await screen.findByTestId('group-profile-summary');
    fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Feeder A renamed' } });
    fireEvent.click(screen.getByTestId('group-save'));
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
    const [id, request] = api.update.mock.calls[0];
    expect(id).toBe('grp-a');
    expect(request).toMatchObject({
      expected_group_revision: 'rev-grp-a-1', expected_connector_revision: 'crev-1',
      group: {
        name: 'Feeder A renamed', layout_version: 'managed_long_v1', storage_mode: 'plain',
        recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' }, row_policy: { interval_seconds: 15 },
      },
    });
  });

  it('keeps every field locked while an operation is in flight', async () => {
    renderSection([longGroup()], 'grp-a', { readonly: true });
    await screen.findByTestId('group-profile-summary');
    expect(screen.getByTestId('group-name')).toBeDisabled();
    expect(screen.getByTestId('group-save')).toBeDisabled();
  });
});

describe('Basic opens Advanced on the same saved profile group through the real host', () => {
  function renderHost(groups: WriteGroup[]) {
    api.list.mockResolvedValue(list(groups));
    return render(
      <QueryClientProvider client={client()}>
        <Step4Database state={batchPanelState()} dispatch={vi.fn()} workspaceId="ws-1" />
      </QueryClientProvider>,
    );
  }

  it('opens the editor of the clicked profile group, not a generic list or another writer', async () => {
    renderHost([longGroup(), cumulativeGroup()]);
    fireEvent.click(await screen.findByTestId('circuit-status-open-advanced-cumulative'));
    expect(await screen.findByTestId('group-editor')).toBeInTheDocument();
    expect(screen.getByTestId('group-name')).toHaveValue('Feeder A (cumulative)');
    expect(screen.getByTestId('group-profile-kind')).toHaveTextContent('step4.group.profile.kind.cumulative');
    expect(screen.getByTestId('step4-advanced-recording')).toHaveAttribute('open');
  });

  it('opens the other profile of the same circuit when that one is clicked', async () => {
    renderHost([longGroup(), cumulativeGroup()]);
    fireEvent.click(await screen.findByTestId('circuit-status-open-advanced-instantaneous'));
    expect(await screen.findByTestId('group-editor')).toBeInTheDocument();
    expect(screen.getByTestId('group-name')).toHaveValue('Feeder A (instantaneous)');
  });

  it('opens a different group when asked again, even from the editor of the first', async () => {
    renderHost([longGroup(), cumulativeGroup()]);
    fireEvent.click(await screen.findByTestId('circuit-status-open-advanced-instantaneous'));
    await screen.findByTestId('group-editor');
    fireEvent.click(screen.getByTestId('group-editor-close'));
    fireEvent.click(screen.getByTestId('circuit-status-open-advanced-cumulative'));
    await waitFor(() => expect(screen.getByTestId('group-name')).toHaveValue('Feeder A (cumulative)'));
  });

  const typed = (code: string) => Object.assign(new Error('Request failed with status code 422'), { response: { status: 422, data: { error: { code } } } });

  it('explains a destination cause once in Basic and moves focus to the destination settings from there', async () => {
    api.schemaPreview.mockRejectedValue(typed('timescale_unavailable'));
    renderHost([longGroup(), cumulativeGroup()]);
    const card = await screen.findByTestId('circuit-blocker-timescale_unavailable-0');
    expect(screen.getAllByTestId(/^circuit-blocker-timescale_unavailable-/)).toHaveLength(1);
    expect(card).toHaveTextContent('Feeder A (instantaneous)');
    expect(card).toHaveTextContent('Feeder A (cumulative)');
    const toggle = screen.getByTestId('database-connection-card-toggle');
    expect(toggle).toHaveAttribute('aria-expanded', 'false');
    fireEvent.click(within(card).getByTestId('circuit-blocker-repair-destination'));
    expect(toggle).toHaveAttribute('aria-expanded', 'true');
    expect(toggle).toHaveFocus();
  });

  it('sends the operator to the counter confirmation of exactly the source that lacks one', async () => {
    renderHost([longGroup(), cumulativeGroup()]);
    fireEvent.click(await screen.findByTestId('circuit-status-open-advanced-cumulative'));
    await screen.findByTestId('group-editor');
    const target = await screen.findByTestId('circuit-counter-confirmation-source-tag-1');
    fireEvent.click(await screen.findByTestId('counter-policy-open-basic-tag-1'));
    expect(target.contains(document.activeElement)).toBe(true);
    expect(screen.getByTestId('circuit-counter-confirmation-source-tag-0').contains(document.activeElement)).toBe(false);
  });

  it('opens the exact blocked group in Advanced from its repair entry', async () => {
    api.schemaPreview.mockRejectedValue(typed('schema_repair_blocked'));
    renderHost([longGroup(), cumulativeGroup()]);
    fireEvent.click(await screen.findByTestId('circuit-blocker-repair-group-grp-c'));
    expect(await screen.findByTestId('group-editor')).toBeInTheDocument();
    expect(screen.getByTestId('group-name')).toHaveValue('Feeder A (cumulative)');
  });
});

describe('Basic says what Advanced and the workspace are for, and links to the existing readiness and activation', () => {
  const readiness = { ready: true, blocking_count: 0, warning_count: 0, issues: [] };
  function renderHost(groups: WriteGroup[]) {
    api.list.mockResolvedValue(list(groups));
    return render(
      <QueryClientProvider client={client()}>
        <Step4Database state={batchPanelState()} dispatch={vi.fn()} workspaceId="ws-1" workspaceReadiness={readiness} />
      </QueryClientProvider>,
    );
  }

  it('names Advanced as the editor of the same saved profile and the workspace as the scope of readiness and activation', async () => {
    renderHost([longGroup(), cumulativeGroup()]);
    const roles = await screen.findByTestId('circuit-profile-roles');
    expect(within(roles).getByTestId('circuit-roles-advanced')).toHaveTextContent('step4.basic.roles.advanced');
    expect(within(roles).getByTestId('circuit-roles-workspace')).toHaveTextContent('step4.basic.roles.workspace');
  });

  it('opens the existing workspace readiness and activation surface, not a second workflow', async () => {
    renderHost([longGroup(), cumulativeGroup()]);
    fireEvent.click(await screen.findByTestId('circuit-profile-open-workspace'));
    const region = await screen.findByTestId('step4-workspace-readiness');
    await waitFor(() => expect(region).toHaveFocus());
    expect(screen.getByTestId('step4-advanced-recording')).toHaveAttribute('open');
    expect(within(region).getByTestId('step4-readiness-panel')).toBeInTheDocument();
    expect(within(region).queryByTestId('group-interval')).not.toBeInTheDocument();
    expect(within(region).queryByTestId('group-storage-managed')).not.toBeInTheDocument();
  });

  it('counts saved and applied groups so an unapplied draft is not presented as recording', async () => {
    renderHost([longGroup({ applied_revision: 'applied-1' }), cumulativeGroup()]);
    const note = await screen.findByTestId('workspace-scope-note');
    await waitFor(() => expect(note).toHaveAttribute('data-total', '2'));
    expect(note).toHaveAttribute('data-applied', '1');
    expect(note).toHaveTextContent('step4.workspace_scope.draft_note');
  });
});

