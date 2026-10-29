import * as React from 'react';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { guidedLifecycle, guidedReadiness } from '../../fixtures/guidedRecordingPreparation';
import { QueryClient, QueryClientProvider, useQueryClient } from '@tanstack/react-query';
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

const clients = new Set<QueryClient>(), listeners = new Set<() => void>();
function publishGroups() { for (const client of clients) client.setQueryData(studioV2WorkspaceKeys.writeGroups(), groupsQuery.data); for (const listener of listeners) listener(); }
function useCanonicalGroups() {
  const client = useQueryClient(), [, refresh] = React.useReducer(value => value + 1, 0);
  React.useLayoutEffect(() => { clients.add(client); client.setQueryData(studioV2WorkspaceKeys.writeGroups(), groupsQuery.data); return () => { clients.delete(client); }; }, [client]);
  React.useEffect(() => { const listener = () => refresh(); listeners.add(listener); return () => { listeners.delete(listener); }; }, []);
  return groupsQuery;
}
function canonical(id: string) { const group = groupsQuery.data.groups.find(item => item.id === id); if (!group) throw new Error('Missing canonical group'); return group; }
async function nativeSave(run: () => Promise<unknown>) {
  const result = await run();
  // Existing native mutations invalidate/read the canonical list before resolving; model that publication, not a false guard.
  groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: currentWsRev, groups: [...serverGroups] }; publishGroups();
  return result;
}
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: vi.fn(async (id: string) => guidedLifecycle(canonical(id), groupsQuery.data.workspace_revision)) } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { readiness: vi.fn(async (id: string) => guidedReadiness(canonical(id), groupsQuery.data.workspace_revision)), schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn() } }));
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', async importOriginal => ({
  ...await importOriginal<typeof import('@/hooks/datalink/useStudioV2WriteGroups')>(),
  useWriteGroupsQuery: useCanonicalGroups,
  useWriteGroupDeliveryQuery: () => ({ data: undefined, isError: false }),
  useUpdateWriteGroupMutation: () => ({ mutateAsync: (input: unknown) => nativeSave(() => updateGroupAsync(input)), isPending: false }),
  useEnsureCircuitProfileMutation: () => ({ mutateAsync: (input: unknown) => nativeSave(() => ensureCircuitProfileAsync(input)), isPending: false }),
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
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, structuralSharing: false }, mutations: { retry: false } } });
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
  groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: currentWsRev, groups: [...serverGroups] };
  groupsQuery.refetch.mockImplementation(async () => {
    groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: currentWsRev, groups: [...serverGroups] };
    publishGroups(); return { data: groupsQuery.data };
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

describe('CircuitProfile multi-circuit chooser and scoped batch edits (Slice A)', () => {
  it('provides a circuit chooser, isolates per-circuit draft and retains circuit_id on rename', async () => {
    renderPanel();

    const chooser = screen.getByTestId('circuit-profile-circuit-chooser');
    expect(chooser).toBeInTheDocument();

    // Verify chooser lists all three saved circuits
    const options = screen.getAllByRole('option');
    const optionValues = options.map((opt) => (opt as HTMLOptionElement).value);
    expect(optionValues).toContain('cir-a');
    expect(optionValues).toContain('cir-b');
    expect(optionValues).toContain('cir-c');

    // Default active is Circuit A
    expect((chooser as HTMLSelectElement).value).toBe('cir-a');
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('Feeder A');
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('15');

    // Switch to Circuit B
    fireEvent.change(chooser, { target: { value: 'cir-b' } });
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('Feeder B');
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('60');

    // Modify Circuit B draft interval to 45
    fireEvent.change(screen.getByTestId('circuit-profile-interval-instantaneous'), { target: { value: '45' } });
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('45');

    // Switch back to Circuit A: Circuit A draft is NOT polluted
    fireEvent.change(chooser, { target: { value: 'cir-a' } });
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('Feeder A');
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('15');

    // Renaming Circuit A keeps its circuit_id intact
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Feeder A Renamed' } });
    expect((chooser as HTMLSelectElement).value).toBe('cir-a');
  });

  it('previews override on B, cancels without mutation, confirms updates A and B only, preserves C', async () => {
    const view = renderPanel();

    // Open batch editor
    const openBatchBtn = screen.getByTestId('circuit-profile-open-batch');
    fireEvent.click(openBatchBtn);

    // Batch editor view is open
    expect(screen.getByTestId('circuit-profile-batch-editor')).toBeInTheDocument();

    // Select Circuit A and B, unselect C
    const checkA = screen.getByTestId('circuit-profile-batch-select-cir-a') as HTMLInputElement;
    const checkB = screen.getByTestId('circuit-profile-batch-select-cir-b') as HTMLInputElement;
    const checkC = screen.getByTestId('circuit-profile-batch-select-cir-c') as HTMLInputElement;

    if (!checkA.checked) fireEvent.click(checkA);
    if (!checkB.checked) fireEvent.click(checkB);
    if (checkC.checked) fireEvent.click(checkC);

    expect(checkA.checked).toBe(true);
    expect(checkB.checked).toBe(true);
    expect(checkC.checked).toBe(false);

    // Set batch instantaneous interval to 30s
    fireEvent.change(screen.getByTestId('circuit-profile-batch-interval-inst'), { target: { value: '30' } });

    // Click preview
    fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));

    // Readonly preview lists old -> new and identifies B's override
    expect(screen.getByTestId('circuit-profile-batch-preview')).toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-batch-override-cir-b')).toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-batch-unselected-cir-c')).toBeInTheDocument();

    // Cancel preview: no mutation
    fireEvent.click(screen.getByTestId('circuit-profile-batch-cancel'));
    expect(screen.queryByTestId('circuit-profile-batch-preview')).not.toBeInTheDocument();
    expect(ensureCircuitProfileAsync).not.toHaveBeenCalled();
    expect(updateGroupAsync).not.toHaveBeenCalled();

    // Open preview again and Confirm
    fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
    fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm'));

    await waitFor(() => {
      // Canonical groups A and B updated via PUT, ensure called ZERO times
      expect(updateGroupAsync).toHaveBeenCalledTimes(2);
      expect(ensureCircuitProfileAsync).toHaveBeenCalledTimes(0);
    });

    // Check calls: only grp-a and grp-b were updated with same IDs
    const calls = updateGroupAsync.mock.calls;
    const calledGroupIds = calls.map((c) => c[0].id);
    expect(calledGroupIds).toContain('grp-a');
    expect(calledGroupIds).toContain('grp-b');
    expect(calledGroupIds).not.toContain('grp-c');

    // Sequential CAS: second call used the workspace revision returned by the first call
    expect(calls[0][0].request.expected_workspace_revision).toBe('wrev-1');
    expect(calls[1][0].request.expected_workspace_revision).toBe('wrev-2');

    // Verify canonical state in mock server: A and B are 30s, C is still 15s
    const updatedA = serverGroups.find((g) => g.id === 'grp-a');
    const updatedB = serverGroups.find((g) => g.id === 'grp-b');
    const preservedC = serverGroups.find((g) => g.id === 'grp-c');
    expect(updatedA?.recording_profile?.interval).toBe(30);
    expect(updatedB?.recording_profile?.interval).toBe(30);
    expect(preservedC).toEqual(
      expect.objectContaining({
        id: 'grp-c',
        revision: 'rev-grp-c-1',
        name: 'Feeder C (Instantaneous)',
        recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second' },
        members: expect.arrayContaining([
          expect.objectContaining({ point_id: 'pt-4', tag_id: 'tag-4', circuit_id: 'cir-c' }),
          expect.objectContaining({ point_id: 'pt-5', tag_id: 'tag-5', circuit_id: 'cir-c' }),
        ]),
      }),
    );

    // Verify full reload / remount: unmount, remount and verify Circuit C in UI
    view.unmount();
    renderPanel();
    const reloadedChooser = screen.getByTestId('circuit-profile-circuit-chooser');
    fireEvent.change(reloadedChooser, { target: { value: 'cir-c' } });
    expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('Feeder C');
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('15');

    // No start or apply was called
    expect(startAsync).not.toHaveBeenCalled();
  });
});
