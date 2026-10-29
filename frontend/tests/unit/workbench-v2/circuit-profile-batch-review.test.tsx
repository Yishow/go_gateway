import { defaultCircuitProfileLocalDraft, loadCircuitProfileLocalDraft, saveCircuitProfileLocalDraft } from '@/features/datalink/workbench-v2/state/circuitProfileBasic';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
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


describe('Root review: batch boundaries', () => {
  it('rejects empty batch interval without default substitution or mutation', async () => {
    renderPanel();
    fireEvent.click(screen.getByTestId('circuit-profile-open-batch'));
    fireEvent.change(screen.getByTestId('circuit-profile-batch-interval-inst'), { target: { value: '' } });
    fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
    const confirm = screen.queryByTestId('circuit-profile-batch-confirm');
    if (confirm) fireEvent.click(confirm);
    await new Promise((resolve) => setTimeout(resolve, 0));
    expect(updateGroupAsync).not.toHaveBeenCalled();
    expect(ensureCircuitProfileAsync).not.toHaveBeenCalled();
    expect(screen.queryByTestId('circuit-profile-batch-preview')).not.toBeInTheDocument();
  });

  it('selects a genuinely new stable circuit instead of falling back to the first saved circuit', async () => {
    renderPanel();
    const chooser = screen.getByTestId('circuit-profile-circuit-chooser') as HTMLSelectElement;
    expect(chooser.value).toBe('cir-a');
    fireEvent.change(chooser, { target: { value: '__new__' } });
    await waitFor(() => {
      expect(chooser.value).not.toBe('cir-a');
      expect(screen.getByTestId('circuit-profile-circuit-name')).toHaveValue('');
    });
  });

  it('does not clone an unassigned legacy draft into another saved circuit', () => {
    const candidates = buildGroupCandidates(makeBatchState()).candidates;
    const legacy = { ...defaultCircuitProfileLocalDraft(candidates, serverGroups, 'cir-a'), circuit_id: '', circuit_name: '', instantaneous_interval_input: '99' };
    window.sessionStorage.setItem(`wbv2.circuit-profile-draft.${encodeURIComponent('ws-1/dev-1')}`, JSON.stringify(legacy));
    const loaded = loadCircuitProfileLocalDraft('ws-1', 'dev-1', candidates, serverGroups, 'cir-b');
    expect(loaded.instantaneous_interval_input).toBe('60');
    expect(loaded.instantaneous_keys).toEqual(defaultCircuitProfileLocalDraft(candidates, serverGroups, 'cir-b').instantaneous_keys);
  });

  it('refreshes the selected local draft after batch save rather than restoring the old override', async () => {
    const candidates = buildGroupCandidates(makeBatchState()).candidates;
    saveCircuitProfileLocalDraft('ws-1', 'dev-1', defaultCircuitProfileLocalDraft(candidates, serverGroups, 'cir-a'));
    renderPanel();
    fireEvent.click(screen.getByTestId('circuit-profile-open-batch'));
    fireEvent.click(screen.getByTestId('circuit-profile-batch-select-cir-c'));
    fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
    fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm'));
    await waitFor(() => expect(serverGroups.find((g) => g.id === 'grp-a')?.recording_profile?.interval).toBe(30));
    await waitFor(() => expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('30'));
  });

  it('does not claim a newly previewed cumulative profile was saved without creating it', async () => {
    ensureCircuitProfileAsync.mockImplementation(async ({ request }) => {
      expect(request.expected_workspace_revision).toBe(currentWsRev);
      const cid = request.group.members[0].circuit_id;
      expect(serverGroups.some((g) => g.members[0]?.circuit_id === cid && g.recording_profile?.kind === 'cumulative')).toBe(false);
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
  });
});
