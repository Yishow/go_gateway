import * as React from 'react';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { QueryClient, QueryClientProvider, useQueryClient } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { BasicRecordingPanel } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel';
import type { MeasurementDefinition } from '@/types/measurement';
import type { RecordingStartOperation } from '@/types/studioV2RecordingStart';
import type { WriteGroup, WriteGroupProfileKind } from '@/types/studioV2WriteGroup';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import { confirmPreparedStart, guidedLifecycle, guidedPreview, guidedReadiness } from '../../fixtures/guidedRecordingPreparation';
import { parseRecordingStartOperation } from '@/utils/studioV2RecordingStartJson';
import type { WriteGroupMutationResponse } from '@/types/studioV2WriteGroup';
import { expandSources } from '../../fixtures/sourceReview';

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
const operationQuery = {
  data: undefined as RecordingStartOperation | undefined,
  isLoading: false,
  isError: false,
  refetch: vi.fn(),
};

const groupListeners = new Set<() => void>();
const publishGroups = () => groupListeners.forEach(listener => listener());
let workspaceVersion = 1;
function useFixtureGroups() {
  const client = useQueryClient();
  const canonical = groupsQuery.data;
  React.useLayoutEffect(() => { client.setQueryData(studioV2WorkspaceKeys.writeGroups(), canonical); }, [client, canonical]);
  const [, refresh] = React.useReducer(value => value + 1, 0);
  React.useEffect(() => { const notify = () => refresh(); groupListeners.add(notify); return () => { groupListeners.delete(notify); }; }, []);
  return groupsQuery;
}
function canonical(id: string) { const group = groupsQuery.data.groups.find(item => item.id === id); if (!group) throw new Error('Missing canonical fixture group'); return group; }
function saved(result: WriteGroupMutationResponse) {
  groupsQuery.data = { workspace_id: result.group.workspace_id, workspace_revision: result.workspace_revision, groups: [...groupsQuery.data.groups.filter(group => group.id !== result.group.id), structuredClone(result.group)] };
  publishGroups(); return result;
}
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', async importOriginal => ({
  ...await importOriginal<typeof import('@/hooks/datalink/useStudioV2WriteGroups')>(),
  useWriteGroupsQuery: useFixtureGroups,
  useWriteGroupDeliveryQuery: () => ({ data: undefined, isError: false }),
  useUpdateWriteGroupMutation: () => ({ mutateAsync: updateGroupAsync, isPending: false }),
  useEnsureCircuitProfileMutation: () => ({ mutateAsync: async (input: unknown) => saved(await ensureCircuitProfileAsync(input)), isPending: false }),
}));
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: vi.fn(async (id: string) => guidedLifecycle(canonical(id), groupsQuery.data.workspace_revision)) } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  readiness: vi.fn(async (id: string) => guidedReadiness(canonical(id), groupsQuery.data.workspace_revision)),
  schemaPreview: vi.fn(async (id: string, request: { workspace_id: string; expected_workspace_revision: string; expected_group_revision: string; expected_connector_revision: string }) => {
    expect(request).toEqual({ workspace_id: 'ws-1', expected_workspace_revision: groupsQuery.data.workspace_revision,
      expected_group_revision: canonical(id).revision, expected_connector_revision: canonical(id).destination.connector_revision });
    return guidedPreview(canonical(id), groupsQuery.data.workspace_revision);
  }),
  schemaApply: vi.fn(async (id: string, request: { token: string; operation_id: string; expected_workspace_revision: string; expected_group_revision: string; expected_connector_revision: string }) => {
    const group = canonical(id), preview = guidedPreview(group, groupsQuery.data.workspace_revision);
    expect(request).toMatchObject({ token: preview.token, operation_id: preview.operation_id, expected_workspace_revision: groupsQuery.data.workspace_revision,
      expected_group_revision: group.revision, expected_connector_revision: group.destination.connector_revision });
    saved({ workspace_revision: `wrev-${++workspaceVersion}`, group: { ...group, destination: { ...group.destination, schema_revision: group.revision, schema_digest: 'schema-proof' } } });
    return { operation_id: request.operation_id, status: 'succeeded', executed_statements: 1, created_at: '', updated_at: '' };
  }), schemaOperation: vi.fn(),
} }));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useEnsureBasicManagedMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useRecordingStartMutation: () => ({ mutateAsync: startAsync, isPending: false }),
  useRecordingStartOperationQuery: () => operationQuery,
}));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useStudioV2WorkspaceMeasurementsQuery: () => measurementsQuery,
  useCreateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useUpdateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

function makeProfileGroup(overrides: Partial<WriteGroup> & { profileKind?: WriteGroupProfileKind } = {}): WriteGroup {
  if (overrides.recording_profile?.kind === 'mixed') throw new Error('Legacy profile fixture cannot represent mixed');
  const profile: WriteGroupProfileKind = overrides.profileKind ?? overrides.recording_profile?.kind ?? 'instantaneous';
  const intervalVal = overrides.recording_profile?.interval ?? (profile === 'cumulative' ? 10 : 15);
  return savedGroup({
    id: profile === 'cumulative' ? 'grp-cum-1' : 'grp-inst-1',
    revision: profile === 'cumulative' ? 'rev-cum-1' : 'rev-inst-1',
    name: profile === 'cumulative' ? 'Main Circuit A (Cumulative)' : 'PLC A (Instantaneous)',
    status: 'ready',
    layout_version: 'managed_long_v1',
    recording_profile: { kind: profile, interval: intervalVal, unit: profile === 'cumulative' ? 'minute' : 'second' },
    storage_mode: 'plain',
    destination: { ...savedGroup().destination, storage_strategy: 'managed', table_name: 'gw_Safe_ws_1_long', database: '/recording.db' },
    members: overrides.members ?? [
      {
        device_id: 'dev-1',
        point_id: 'pt-0',
        tag_id: 'tag-0',
        source_revision: 's0',
        mapping_revision: 'm0',
        target_column: '',
        required: false,
        measurement_id: profile === 'cumulative' ? 'meas-0' : '',
        circuit_id: 'cir-main-1',
        circuit_name: 'Main Circuit A',
      },
    ],
    row_policy: {
      interval_seconds: profile === 'cumulative' ? intervalVal * 60 : intervalVal,
      allowed_lateness_seconds: 0,
      incomplete_policy: 'partial',
    },
    write_policy: { mode: 'append', dedupe_capability: 'receipt', recreate_missing: true },
    ...overrides,
  });
}

function makeConfirmedCounterMeasurement(pointId: string, tagId: string): MeasurementDefinition {
  return {
    id: `meas-${tagId}`,
    workspace_id: 'ws-1',
    definition_revision: 'mrev-1',
    name: `Energy ${tagId}`,
    equipment_id: 'eq-1',
    device_id: 'dev-1',
    point_id: pointId,
    tag_id: tagId,
    source_binding_revision: 's1',
    quantity: 'energy',
    semantic_kind: 'counter',
    unit: 'kWh',
    series_epoch: '1',
    counter_policy: { allow_negative: false },
    created_at: '2026-10-20T00:00:00Z',
    updated_at: '2026-10-20T00:00:00Z',
  };
}

function renderPanel() {
  // Hook fixture exposes this exact accepted object, as the real query observer does.
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, structuralSharing: false }, mutations: { retry: false } } });
  const view = render(
    <QueryClientProvider client={client}>
      <BasicRecordingPanel state={savedGroupState()} workspaceId="ws-1" readonly={false} />
    </QueryClientProvider>,
  );
  return { ...view, client };
}

beforeEach(() => {
  vi.clearAllMocks(); workspaceVersion = 1;
  window.sessionStorage.clear();
  groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] };
  groupsQuery.refetch.mockImplementation(async () => ({ data: groupsQuery.data }));
  measurementsQuery.data = [];
  measurementsQuery.refetch.mockImplementation(async () => ({ data: measurementsQuery.data }));
  ensureCircuitProfileAsync.mockImplementation(async ({ request }) => ({
    workspace_revision: `wrev-${++workspaceVersion}`,
    group: makeProfileGroup({
      profileKind: request.group.recording_profile?.kind ?? 'instantaneous',
      recording_profile: request.group.recording_profile,
      storage_mode: request.group.storage_mode ?? 'plain',
      row_policy: request.group.row_policy,
      write_policy: request.group.write_policy,
      members: request.group.members,
    }),
  }));
  startAsync.mockImplementation(async request => {
    // Production parser rejects malformed succeeded groups; both selected profiles retain their own revisions.
    const parsed = parseRecordingStartOperation({ operation_id: 'op-circuit-1', action: 'recording_start', status: 'succeeded', stage: 'complete',
      intent_digest: 'b'.repeat(64), workspace_id: request.workspace_id, setup_revision: request.expected_workspace_revision, device_ids: request.device_ids,
      groups: request.groups.map((item: { group_id: string; expected_group_revision: string }) => ({ group_id: item.group_id, group_revision: item.expected_group_revision, applied_revision: item.expected_group_revision, saved: true, ready: true, applied: true })),
      devices: [{ device_id: 'dev-1', activated: true }], created_at: '2026-10-24T00:00:00Z', updated_at: '2026-10-24T00:00:01Z' });
    if (!parsed) throw new Error('Invalid operation fixture'); return parsed;
  });
});

describe('CircuitProfile Independent Sources per Profile (Slice A)', () => {
  it('allows selecting power (gauge) for instantaneous and energy (counter) for cumulative on one circuit', async () => {
    // tag-0 is line.temperature (float64, gauge)
    // tag-1 is line.pressure (float64, gauge) -> we will provide counter measurement for tag-1 to simulate energy
    // tag-2 is line.batch (string)
    measurementsQuery.data = [makeConfirmedCounterMeasurement('pt-1', 'tag-1')];

    const host = renderPanel();
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Main Feeder 1' } });
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));

    expandSources();
    // tag-0: select for instantaneous ONLY (uncheck cumulative)
    // tag-1: select for cumulative ONLY (uncheck instantaneous)
    // tag-2: uncheck both
    const instTag0 = screen.getByTestId('circuit-profile-member-inst-tag-0');
    const cumTag0 = screen.getByTestId('circuit-profile-member-cum-tag-0');
    const instTag1 = screen.getByTestId('circuit-profile-member-inst-tag-1');
    const cumTag1 = screen.getByTestId('circuit-profile-member-cum-tag-1');
    const instTag2 = screen.getByTestId('circuit-profile-member-inst-tag-2');
    const cumTag2 = screen.getByTestId('circuit-profile-member-cum-tag-2');

    expect(instTag0).toBeInTheDocument();
    expect(cumTag0).toBeInTheDocument();

    // Configure independent selections
    if ((cumTag0 as HTMLInputElement).checked) fireEvent.click(cumTag0);
    if (!(instTag0 as HTMLInputElement).checked) fireEvent.click(instTag0);

    if ((instTag1 as HTMLInputElement).checked) fireEvent.click(instTag1);
    if (!(cumTag1 as HTMLInputElement).checked) fireEvent.click(cumTag1);

    if ((instTag2 as HTMLInputElement).checked) fireEvent.click(instTag2);
    if ((cumTag2 as HTMLInputElement).checked) fireEvent.click(cumTag2);

    // Because tag-0 (gauge) is NOT in cumulative, it does not require counter confirmation!
    // And tag-1 (counter) is in cumulative and has confirmed measurement, so unconfirmed is 0!
    expect(screen.queryByTestId('circuit-counter-confirmation')).not.toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled();

    // Prepare creates two groups with independent member sets
    fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
    await waitFor(() => expect(ensureCircuitProfileAsync).toHaveBeenCalledTimes(2));

    const firstGroup = ensureCircuitProfileAsync.mock.calls[0][0].request.group;
    const secondGroup = ensureCircuitProfileAsync.mock.calls[1][0].request.group;

    const instCall = firstGroup.recording_profile.kind === 'instantaneous' ? firstGroup : secondGroup;
    const cumCall = firstGroup.recording_profile.kind === 'cumulative' ? firstGroup : secondGroup;

    expect(instCall.members).toHaveLength(1);
    expect(instCall.members[0].tag_id).toBe('tag-0');
    expect(instCall.members[0].circuit_name).toBe('Main Feeder 1');

    expect(cumCall.members).toHaveLength(1);
    expect(cumCall.members[0].tag_id).toBe('tag-1');
    expect(cumCall.members[0].circuit_name).toBe('Main Feeder 1');
    expect(cumCall.members[0].circuit_id).toBe(instCall.members[0].circuit_id);

    await waitFor(() => expect(host.client.getQueryData(studioV2WorkspaceKeys.writeGroups())).toBe(groupsQuery.data));
    // Can start both groups only after separate schema confirmations and explicit start review.
    await confirmPreparedStart('circuit-profile');
    await waitFor(() => expect(startAsync).toHaveBeenCalledTimes(1));
    expect(startAsync.mock.calls[0][0].groups.map((item: { group_id: string }) => item.group_id).sort()).toEqual(['grp-cum-1', 'grp-inst-1']);
  });

  it('preserves independent instantaneous and cumulative selections across reload', async () => {
    measurementsQuery.data = [makeConfirmedCounterMeasurement('pt-1', 'tag-1')];

    const first = renderPanel();
    fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Main Feeder 1' } });
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));

    expandSources();
    const cumTag0 = screen.getByTestId('circuit-profile-member-cum-tag-0');
    const instTag1 = screen.getByTestId('circuit-profile-member-inst-tag-1');
    if ((cumTag0 as HTMLInputElement).checked) fireEvent.click(cumTag0);
    if ((instTag1 as HTMLInputElement).checked) fireEvent.click(instTag1);
    fireEvent.click(screen.getByTestId('circuit-profile-member-cum-tag-1'));

    first.unmount();

    renderPanel();
    expandSources();
    expect((screen.getByTestId('circuit-profile-member-inst-tag-0') as HTMLInputElement).checked).toBe(true);
    expect((screen.getByTestId('circuit-profile-member-cum-tag-0') as HTMLInputElement).checked).toBe(false);
    expect((screen.getByTestId('circuit-profile-member-inst-tag-1') as HTMLInputElement).checked).toBe(false);
    expect((screen.getByTestId('circuit-profile-member-cum-tag-1') as HTMLInputElement).checked).toBe(true);
  });

  it('hydrates legacy instantaneous selection without inventing cumulative intent from shared selected_keys', async () => {
    window.sessionStorage.setItem(
      'wbv2.circuit-profile-draft.ws-1%2Fdev-1',
      JSON.stringify({
        circuit_id: 'cir-legacy-1',
        circuit_name: 'Legacy Feeder',
        selected_keys: ['dev-1\0pt-0\0tag-0'],
        instantaneous_enabled: true,
        cumulative_enabled: true,
      }),
    );
    renderPanel();
    expandSources();
    expect((screen.getByTestId('circuit-profile-member-inst-tag-0') as HTMLInputElement).checked).toBe(true);
    expect((screen.getByTestId('circuit-profile-member-cum-tag-0') as HTMLInputElement).checked).toBe(false);
    expect(screen.getByTestId('circuit-profile-cumulative-sources-error')).toBeVisible();
    expect((screen.getByTestId('circuit-profile-member-inst-tag-1') as HTMLInputElement).checked).toBe(false);
    expect((screen.getByTestId('circuit-profile-member-cum-tag-1') as HTMLInputElement).checked).toBe(false);
  });
});
