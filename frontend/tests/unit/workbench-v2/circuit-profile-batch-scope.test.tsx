import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { BasicRecordingPanel } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel';
import type { Mapping, Point, WorkbenchV2State } from '@/features/datalink/workbench-v2/state/types';
import type { MeasurementDefinition } from '@/types/measurement';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: Record<string, unknown>) =>
      options && 'count' in options ? `${key}:${options.count}` : key,
  }),
}));

const groupsQuery = {
  data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] as WriteGroup[] },
  isLoading: false,
  isError: false,
  refetch: vi.fn(),
};
const measurementsQuery = {
  data: [] as MeasurementDefinition[],
  isLoading: false,
  isError: false,
  refetch: vi.fn(),
};
const ensureCircuitProfileAsync = vi.fn();
const updateGroupAsync = vi.fn();
const startAsync = vi.fn();

vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
  useWriteGroupsQuery: () => groupsQuery,
  useWriteGroupDeliveryQuery: () => ({ data: undefined, isError: false }),
  useUpdateWriteGroupMutation: () => ({ mutateAsync: updateGroupAsync, isPending: false }),
  useEnsureCircuitProfileMutation: () => ({ mutateAsync: ensureCircuitProfileAsync, isPending: false }),
  useWriteGroupSchemaPreviewMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useWriteGroupSchemaApplyMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useEnsureBasicManagedMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useRecordingStartMutation: () => ({ mutateAsync: startAsync, isPending: false }),
  useRecordingStartOperationQuery: () => ({ data: undefined, isLoading: false, isError: false, refetch: vi.fn() }),
}));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useStudioV2WorkspaceMeasurementsQuery: () => measurementsQuery,
  useCreateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useUpdateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel', () => ({
  GroupSchemaPanel: (props: { group: WriteGroup; onApplied?: () => void }) => (
    <div data-testid={`circuit-schema-panel-${props.group.recording_profile?.kind ?? 'wide'}`}>
      <span>{props.group.id}</span>
      <button type="button" onClick={() => props.onApplied?.()}>Confirm</button>
    </div>
  ),
}));

function makeBatchState(): WorkbenchV2State {
  const base = savedGroupState();
  const points: Point[] = [];
  const mappings: Record<string, Mapping> = {};
  for (let i = 0; i < 6; i++) {
    const pid = `rule-1-p-${i}`;
    const name = `sensor.${i}`;
    points.push({
      id: pid, device_id: 'dev-1', rule_id: 'rule-1', rule_name: 'Rule 1', name, address: `4000${i + 1}`,
      data_type: 'float64', function: 'holding_register', width: 1, enabled: true, skipped: false,
      _rule_scale: 1, _rule_offset: 0,
    });
    mappings[pid] = {
      point_id: pid, tag_key: name, display_name: name, unit: '', target_type: 'float64', scale: 1, offset: 0,
      enabled: true, persisted: true, save_state: 'saved', persisted_point_id: `pt-${i}`, tag_id: `tag-${i}`, device_id: 'dev-1',
    };
  }
  return { ...base, points, mappings };
}

function makeCircuitGroup(circuitId: string, circuitName: string, intervalSeconds: number, id: string, pointIndices: [number, number]): WriteGroup {
  return savedGroup({
    id,
    revision: `rev-${id}-1`,
    name: `${circuitName} (Instantaneous)`,
    status: 'ready',
    layout_version: 'managed_long_v1',
    recording_profile: { kind: 'instantaneous', interval: intervalSeconds, unit: 'second' },
    storage_mode: 'plain',
    destination: { ...savedGroup().destination, storage_strategy: 'managed', table_name: 'gw_Safe_ws_1_long' },
    members: pointIndices.map((i) => ({
      device_id: 'dev-1',
      point_id: `pt-${i}`,
      tag_id: `tag-${i}`,
      source_revision: `s${i}`,
      mapping_revision: `m${i}`,
      target_column: '',
      required: false,
      circuit_id: circuitId,
      circuit_name: circuitName,
    })),
    row_policy: { interval_seconds: intervalSeconds, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
    write_policy: { mode: 'append', dedupe_capability: 'receipt', recreate_missing: true },
  });
}

function renderPanel(state = makeBatchState()) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <BasicRecordingPanel state={state} workspaceId="ws-1" readonly={false} />
    </QueryClientProvider>,
  );
}

let serverGroups: WriteGroup[] = [];
let currentWsRev = 'wrev-1';

beforeEach(() => {
  vi.clearAllMocks();
  window.sessionStorage.clear();
  currentWsRev = 'wrev-1';
  serverGroups = [
    makeCircuitGroup('cir-a', 'Feeder A', 15, 'grp-a', [0, 1]),
    makeCircuitGroup('cir-b', 'Feeder B', 60, 'grp-b', [2, 3]),
    makeCircuitGroup('cir-c', 'Feeder C', 15, 'grp-c', [4, 5]),
  ];
  groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: currentWsRev, groups: serverGroups };
  groupsQuery.refetch.mockImplementation(async () => {
    groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: currentWsRev, groups: serverGroups };
    return { data: groupsQuery.data };
  });
  measurementsQuery.data = [];
  measurementsQuery.refetch.mockImplementation(async () => ({ data: measurementsQuery.data }));

  updateGroupAsync.mockImplementation(async ({ id, request }) => {
    if (request.expected_workspace_revision !== currentWsRev) {
      throw new Error(`stale_workspace_revision: expected ${request.expected_workspace_revision}, got ${currentWsRev}`);
    }
    const idx = serverGroups.findIndex((g) => g.id === id);
    if (idx === -1) throw new Error(`group_not_found: ${id}`);
    const existing = serverGroups[idx];
    if (existing.revision !== request.expected_group_revision) {
      throw new Error(`stale_group_revision: expected ${request.expected_group_revision}, got ${existing.revision}`);
    }
    const revNum = parseInt(currentWsRev.replace('wrev-', ''), 10) + 1;
    currentWsRev = `wrev-${revNum}`;
    const nextGroup: WriteGroup = {
      ...existing,
      ...request.group,
      id,
      revision: `rev-${id}-${revNum}`,
    };
    serverGroups[idx] = nextGroup;
    return { workspace_revision: currentWsRev, group: nextGroup };
  });
});



describe('Root review: batch scope and effective selection', () => {
  it('locks the main circuit chooser and start while a batch save is pending', async () => {
    updateGroupAsync.mockImplementation(() => new Promise(() => {}));
    renderPanel();
    fireEvent.click(screen.getByTestId('circuit-profile-open-batch'));
    fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
    fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm'));
    await waitFor(() => expect(updateGroupAsync).toHaveBeenCalledTimes(1));
    expect(screen.getByTestId('circuit-profile-circuit-chooser')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-start')).toBeDisabled();
  });

  it('restores the same unsaved new circuit identity and draft after reload', async () => {
    const panel = renderPanel();
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-chooser'), { target: { value: '__new__' } });
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'New circuit draft' } });
    const id = (screen.getByTestId('circuit-profile-circuit-chooser') as HTMLSelectElement).value;
    expect(id).not.toBe('cir-a');
    panel.unmount();
    renderPanel();
    await waitFor(() => {
      expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('New circuit draft');
      expect(screen.getByTestId('circuit-profile-circuit-chooser')).toHaveValue(id);
    });
  });

  it('gives both batch interval textboxes accessible names and links invalid state', () => {
    renderPanel();
    fireEvent.click(screen.getByTestId('circuit-profile-open-batch'));
    const inst = screen.getByTestId('circuit-profile-batch-interval-inst');
    const cum = screen.getByTestId('circuit-profile-batch-interval-cum');
    expect(inst).toHaveAccessibleName();
    expect(cum).toHaveAccessibleName();
    fireEvent.change(inst, { target: { value: '' } });
    expect(inst).toHaveAttribute('aria-invalid', 'true');
    const errorId = inst.getAttribute('aria-describedby');
    expect(errorId).toBeTruthy();
    expect(document.getElementById(errorId!)).toHaveAttribute('role', 'alert');
  });

  it('reflects a newly saved cumulative profile as selected with its actual sources', async () => {
    ensureCircuitProfileAsync.mockImplementation(async ({ request }) => {
      expect(request.expected_workspace_revision).toBe(currentWsRev);
      const cid = request.group.members[0].circuit_id;
      currentWsRev = `wrev-${Number(currentWsRev.slice(5)) + 1}`;
      const created = savedGroup({ ...request.group, id: `cum-${cid}`, revision: `rev-cum-${cid}`, status: 'draft' });
      serverGroups = [...serverGroups, created];
      return { workspace_revision: currentWsRev, group: created };
    });
    renderPanel();
    fireEvent.click(screen.getByTestId('circuit-profile-open-batch'));
    fireEvent.click(screen.getByLabelText('step4.basic.circuit.batch_inst_interval'));
    fireEvent.click(screen.getByLabelText('step4.basic.circuit.batch_cum_interval'));
    fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
    fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm'));
    await waitFor(() => expect(ensureCircuitProfileAsync).toHaveBeenCalledTimes(3));
    await waitFor(() => expect(screen.queryByTestId('circuit-profile-batch-editor')).not.toBeInTheDocument());
    expect(screen.getByTestId('circuit-profile-toggle-cumulative')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-member-cum-tag-0')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-member-cum-tag-1')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-member-cum-tag-2')).not.toBeChecked();
    expect(startAsync).not.toHaveBeenCalled();
  });
});
