import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, within } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { BasicRecordingPanel } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel';
import type { Mapping, Point, WorkbenchV2State } from '@/features/datalink/workbench-v2/state/types';
import type { MeasurementDefinition } from '@/types/measurement';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { WriteGroupDelivery, WriteGroupDeliveryStages, WriteGroupSourceEvidenceItem, WriteGroupSourceRecord } from '@/types/studioV2WriteGroupDelivery';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
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
const measurementsQuery = { data: [] as MeasurementDefinition[], isLoading: false, isError: false, refetch: vi.fn() };
const deliveryQuery = { data: undefined as WriteGroupDelivery | undefined, isError: false };
/** A group's own answer when the test gives one; every other group gets `deliveryQuery.data`. */
const deliveryFor = new Map<string, WriteGroupDelivery>();
const ensureCircuitAsync = vi.fn();
const updateGroupAsync = vi.fn();
const startAsync = vi.fn();
const openAdvanced = vi.fn();

vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
  useWriteGroupsQuery: () => groupsQuery,
  useWriteGroupDeliveryQuery: (groupId?: string) => ({ data: deliveryFor.get(groupId ?? '') ?? deliveryQuery.data, isError: deliveryQuery.isError }),
  useUpdateWriteGroupMutation: () => ({ mutateAsync: updateGroupAsync, isPending: false }),
  useEnsureCircuitProfileMutation: () => ({ mutateAsync: ensureCircuitAsync, isPending: false }),
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
  GroupSchemaPanel: (props: { group: WriteGroup }) => <div data-testid="circuit-schema-panel">{props.group.id}</div>,
}));

/** `count` saved sources of one device; each has a display name except every fifth, which is known by its address only. */
function sourcesState(count: number): WorkbenchV2State {
  const base = savedGroupState();
  const points: Point[] = [];
  const mappings: Record<string, Mapping> = {};
  for (let i = 0; i < count; i += 1) {
    const id = `rule-1-p-${i}`;
    const named = i % 5 !== 0;
    points.push({
      id, device_id: 'dev-1', rule_id: 'rule-1', rule_name: 'Rule 1', name: `sensor.${i}`, address: `4${String(i + 1).padStart(4, '0')}`,
      data_type: 'float64', function: 'holding_register', width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0,
    });
    mappings[id] = {
      point_id: id, tag_key: `sensor.${i}`, display_name: named ? `Sensor ${i}` : '', unit: '', target_type: 'float64', scale: 1, offset: 0,
      enabled: true, persisted: true, save_state: 'saved', persisted_point_id: `pt-${i}`, tag_id: `tag-${i}`, device_id: 'dev-1',
    };
  }
  return { ...base, points, mappings };
}

/** A saved profile group of one circuit over the given point indexes; `ghost` adds a member no current source matches. */
function profileGroup(layout: 'managed_long_v1' | 'managed_wide_profile_v1', points: number[], over: Partial<WriteGroup> = {}, ghost = false): WriteGroup {
  const wide = layout === 'managed_wide_profile_v1';
  const member = (i: number | 'ghost') => ({
    device_id: 'dev-1', point_id: i === 'ghost' ? 'pt-lost' : `pt-${i}`, tag_id: i === 'ghost' ? 'tag-lost' : `tag-${i}`,
    source_revision: 's', mapping_revision: 'm', target_column: '', required: false, circuit_id: 'cir-main', circuit_name: 'Main',
    ...(wide ? { value_role: 'gauge' as const, source_value_type: 'float64' as const } : {}),
  });
  return savedGroup({
    id: 'grp-1', revision: 'rev-1', applied_revision: 'rev-1', name: 'Main (Instantaneous)', status: 'ready', layout_version: layout,
    recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: 'mean' },
    ...(wide ? {} : { storage_mode: 'plain' as const }),
    destination: { ...savedGroup().destination, storage_strategy: 'managed', table_name: wide ? '' : 'gw_long_fixture' },
    members: [...points.map(member), ...(ghost ? [member('ghost')] : [])],
    row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
    write_policy: wide ? { mode: 'append', dedupe_capability: 'receipt' } : { mode: 'append', dedupe_capability: 'receipt', recreate_missing: true },
    ...over,
  });
}

function renderPanel(state: WorkbenchV2State) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <BasicRecordingPanel state={state} workspaceId="ws-1" readonly={false} onOpenAdvanced={openAdvanced} />
    </QueryClientProvider>,
  );
}

const range = (n: number) => Array.from({ length: n }, (_, i) => i);

const STAGES: WriteGroupDeliveryStages = { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 1, skipped: 0 };

/** A verified record in the real wire form of the existing projection. */
function sourceRecord(over: Partial<WriteGroupSourceRecord> = {}): WriteGroupSourceRecord {
  return {
    available: true, group_revision: 'rev-1', connector_revision: 'crev-1', record_id: 'r0', effect_key: 'e0', payload_digest: 'd'.repeat(64),
    profile_kind: 'instantaneous', window_start: '2026-01-01T00:00:00Z', window_end: '2026-01-01T00:00:15Z', status: 'ok', raw_present: true,
    delivery: 'sql_committed', ...over,
  };
}

/** The per source aggregation evidence the projection will add; the field names and values are the real encoder's. */
interface AggregationEvidence {
  method: 'arithmetic-mean'; version: 'arithmetic-mean-v1'; basis: 'exact' | 'float64'; source_type: string; result_type: string;
  valid_count: number; rejected_count: number; coverage: 'complete' | 'partial' | 'unknown' | 'missing';
}

/** The scoped delivery response of one group's applied revision with the given per source evidence. */
function deliveryWith(sources: WriteGroupSourceEvidenceItem[], groupId = 'grp-1', revision = 'rev-1'): WriteGroupDelivery {
  return {
    group_id: groupId, intake: { state: 'active' }, stages: STAGES, last_sql_committed_at: null, oldest_pending_seconds: 0,
    no_data_buckets: 0, skipped_buckets: 0, backlog: [], quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
    revision_stages: { group_revision: revision, stages: STAGES },
    source_evidence: { state: 'available', group_revision: revision, connector_id: 'conn-1', connector_revision: 'crev-1', sources },
  };
}

/** The scoped delivery response of the applied revision with one source's latest record. */
function deliveryOf(latest: WriteGroupSourceRecord): WriteGroupDelivery {
  return deliveryWith([{ device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0', latest, last_sql_committed: null }]);
}
const memberRows = () => screen.queryAllByTestId(/^circuit-profile-member-inst-/);

beforeEach(() => {
  vi.clearAllMocks();
  window.sessionStorage.clear();
  groupsQuery.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] };
  groupsQuery.refetch.mockImplementation(async () => ({ data: groupsQuery.data }));
  measurementsQuery.data = [];
  deliveryQuery.data = undefined;
  deliveryFor.clear();
});

describe('compact source review', () => {
  it('collapses every source list behind an explicit expand, however short it is', () => {
    renderPanel(sourcesState(3));

    expect(screen.getByTestId('circuit-profile-sources-summary')).toHaveTextContent('step4.basic.sources.count:3');
    expect(screen.getByTestId('circuit-profile-sources-toggle')).toHaveAttribute('aria-expanded', 'false');
    expect(memberRows()).toHaveLength(0);

    fireEvent.click(screen.getByTestId('circuit-profile-sources-toggle'));
    expect(memberRows()).toHaveLength(3);
    for (const call of [ensureCircuitAsync, updateGroupAsync, startAsync]) expect(call).not.toHaveBeenCalled();
  });

  it('shows the count and the unavailable source of a 32 source group without 32 open rows', () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [profileGroup('managed_long_v1', range(31), {}, true)] };
    renderPanel(sourcesState(31));

    expect(screen.getByTestId('circuit-profile-sources-summary')).toHaveTextContent('step4.basic.sources.count:32');
    expect(screen.getByTestId('circuit-profile-sources-exceptions')).toHaveTextContent('step4.basic.sources.unavailable:1');
    expect(screen.getByTestId('circuit-profile-sources-toggle')).toHaveAttribute('aria-expanded', 'false');
    expect(memberRows()).toHaveLength(0);
  });

  it('keeps the repair of the unavailable source reachable while collapsed, without changing anything', () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [profileGroup('managed_long_v1', range(31), {}, true)] };
    renderPanel(sourcesState(31));

    fireEvent.click(screen.getByTestId('circuit-profile-sources-repair'));
    expect(openAdvanced).toHaveBeenCalledWith('grp-1');
    for (const call of [ensureCircuitAsync, updateGroupAsync, startAsync]) expect(call).not.toHaveBeenCalled();
  });

  it('reveals every scoped source inside an internally scrolling list when expanded, and collapses again', () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [profileGroup('managed_long_v1', range(31), {}, true)] };
    renderPanel(sourcesState(31));

    fireEvent.click(screen.getByTestId('circuit-profile-sources-toggle'));
    expect(screen.getByTestId('circuit-profile-sources-toggle')).toHaveAttribute('aria-expanded', 'true');
    expect(memberRows()).toHaveLength(31);
    expect(screen.getByTestId('circuit-profile-member-list').className).toMatch(/overflow-y-auto/);
    // A source without a saved name is known by its address; its technical key stays in diagnostics.
    expect(within(screen.getByTestId('circuit-profile-member-list')).getAllByText(/40001/).length).toBeGreaterThan(0);
    expect(screen.queryByText('tag-0')).toBeNull();
    expect(screen.getByTestId('circuit-profile-sources-exceptions')).toHaveTextContent('step4.basic.sources.unavailable:1');

    fireEvent.click(screen.getByTestId('circuit-profile-sources-toggle'));
    expect(memberRows()).toHaveLength(0);
    for (const call of [ensureCircuitAsync, updateGroupAsync, startAsync]) expect(call).not.toHaveBeenCalled();
  });
});

describe('the legacy wide source list', () => {
  /** A saved legacy wide group over `count` mapped sources plus one saved member no current source matches. */
  function legacyWide(count: number): WriteGroup {
    const member = (i: number | 'lost') => ({
      device_id: 'dev-1', point_id: i === 'lost' ? 'pt-lost' : `pt-${i}`, tag_id: i === 'lost' ? 'tag-lost' : `tag-${i}`,
      source_revision: 's', mapping_revision: 'm', target_column: `col_${i}`, required: true,
    });
    return savedGroup({
      id: 'legacy-wide', name: 'Legacy wide', status: 'ready', basic_managed_device_id: 'dev-1', layout_version: 'managed_wide_v1',
      destination: { ...savedGroup().destination, storage_strategy: 'managed', table_name: 'basic_readings' },
      members: [...range(count).map(member), member('lost')],
    });
  }
  const listed = () => within(screen.getByTestId('basic-recording-selected-members')).queryAllByRole('listitem');

  it('keeps the saved members behind an explicit expand, with the unavailable one summarized and repairable', () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [legacyWide(31)] };
    renderPanel(sourcesState(31));
    fireEvent.click(screen.getByTestId('device-recording-edit-existing'));

    expect(screen.getByTestId('basic-recording-sources-summary')).toHaveTextContent('step4.basic.sources.count:32');
    expect(screen.getByTestId('basic-recording-sources-exceptions')).toHaveTextContent('step4.basic.sources.unavailable:1');
    expect(screen.getByTestId('basic-recording-sources-toggle')).toHaveAttribute('aria-expanded', 'false');
    expect(listed()).toHaveLength(0);

    fireEvent.click(screen.getByTestId('basic-recording-sources-repair'));
    expect(openAdvanced).toHaveBeenCalledWith('legacy-wide');

    fireEvent.click(screen.getByTestId('basic-recording-sources-toggle'));
    expect(listed()).toHaveLength(32);
    for (const call of [ensureCircuitAsync, updateGroupAsync, startAsync]) expect(call).not.toHaveBeenCalled();
  });
});

describe('the source status list', () => {
  const item = (tag: string, latest: WriteGroupSourceEvidenceItem['latest']): WriteGroupSourceEvidenceItem => (
    { device_id: 'dev-1', point_id: `pt-${tag.replace('tag-', '')}`, tag_id: tag, latest, last_sql_committed: null }
  );
  /** 31 mapped sources and one saved member no current source matches; one record cannot be verified and one is bad. */
  function statusOf32(): WriteGroupSourceEvidenceItem[] {
    const sources = range(31).map((i) => {
      const record = i === 5
        ? { available: false as const, cause: 'payload_invalid' as const, delivery: 'queued' as const }
        : sourceRecord({ status: i === 6 ? 'bad' : 'ok', raw_present: i !== 6, record_id: `r${i}`, effect_key: `e${i}` });
      return item(`tag-${i}`, record);
    });
    return [...sources, item('tag-lost', sourceRecord({ record_id: 'rl', effect_key: 'el' }))];
  }
  const statusRows = () => screen.queryAllByTestId(/^circuit-status-source-tag-/);

  it('keeps the 32 status rows behind a count, with the exceptions and their repair in view', () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [profileGroup('managed_long_v1', range(31), {}, true)] };
    deliveryQuery.data = deliveryWith(statusOf32());
    renderPanel(sourcesState(31));

    expect(screen.getByTestId('circuit-status-sources-grp-1-summary')).toHaveTextContent('step4.basic.sources.count:32');
    expect(screen.getByTestId('circuit-status-sources-grp-1-exceptions')).toHaveTextContent('step4.basic.sources.needs_check:2');
    expect(screen.getByTestId('circuit-status-sources-grp-1-toggle')).toHaveAttribute('aria-expanded', 'false');
    expect(statusRows()).toHaveLength(0);

    fireEvent.click(screen.getByTestId('circuit-status-sources-grp-1-repair'));
    expect(openAdvanced).toHaveBeenCalledWith('grp-1');
    for (const call of [ensureCircuitAsync, updateGroupAsync, startAsync]) expect(call).not.toHaveBeenCalled();
  });

  it('lists every status row inside its own scrolling list when expanded, and never names a source by its technical key', () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [profileGroup('managed_long_v1', range(31), {}, true)] };
    deliveryQuery.data = deliveryWith(statusOf32());
    renderPanel(sourcesState(31));

    fireEvent.click(screen.getByTestId('circuit-status-sources-grp-1-toggle'));
    expect(statusRows()).toHaveLength(32);
    expect(screen.getByTestId('circuit-status-source-list-grp-1').className).toMatch(/overflow-y-auto/);
    expect(screen.getByTestId('circuit-status-source-tag-1')).toHaveTextContent('Sensor 1');
    // A source with neither a saved name nor an address is an unavailable saved member that is repaired in its group.
    const lost = screen.getByTestId('circuit-status-source-tag-lost');
    expect(lost).toHaveTextContent('step4.basic.saved_member_unavailable');
    expect(lost).not.toHaveTextContent('tag-lost');
    fireEvent.click(within(lost).getByTestId('circuit-status-source-repair-tag-lost'));
    expect(openAdvanced).toHaveBeenCalledWith('grp-1');
  });

  it('gives every group its own review: the toggles and the lists they control are distinct', () => {
    const cumulative = profileGroup('managed_long_v1', [0, 1], {
      id: 'grp-2', revision: 'rev-2', applied_revision: 'rev-2', name: 'Main (Cumulative)',
      recording_profile: { kind: 'cumulative', interval: 10, unit: 'minute' },
    });
    groupsQuery.data = { ...groupsQuery.data, groups: [profileGroup('managed_long_v1', [0, 1]), cumulative] };
    deliveryFor.set('grp-1', deliveryWith([item('tag-0', sourceRecord())]));
    deliveryFor.set('grp-2', deliveryWith([item('tag-0', sourceRecord({ profile_kind: 'cumulative', status: 'baseline', delta_known: false, group_revision: 'rev-2' }))], 'grp-2', 'rev-2'));
    renderPanel(sourcesState(3));

    const first = screen.getByTestId('circuit-status-sources-grp-1-toggle');
    const second = screen.getByTestId('circuit-status-sources-grp-2-toggle');
    expect(first.getAttribute('aria-controls')).not.toBe(second.getAttribute('aria-controls'));
    fireEvent.click(second);
    expect(second).toHaveAttribute('aria-expanded', 'true');
    expect(first).toHaveAttribute('aria-expanded', 'false');
    expect(document.getElementById(second.getAttribute('aria-controls') as string)).toBeInTheDocument();
  });
});

describe('missing values by layout', () => {
  it('tells a wide group which sources are required, how a gap is handled, and reaches the existing member controls', () => {
    const wide = profileGroup('managed_wide_profile_v1', [0, 1]);
    wide.members[0] = { ...wide.members[0], required: true };
    groupsQuery.data = { ...groupsQuery.data, groups: [wide] };
    renderPanel(sourcesState(3));

    const policy = screen.getByTestId('circuit-profile-required-policy');
    expect(policy).toHaveTextContent('step4.basic.required.count:1');
    expect(policy).toHaveTextContent('step4.basic.required.wide_rule');
    fireEvent.click(screen.getByTestId('circuit-profile-required-open'));
    expect(openAdvanced).toHaveBeenCalledWith('grp-1');
    // Looking at the rule never changes a member.
    expect(updateGroupAsync).not.toHaveBeenCalled();
    expect(wide.members.map((member) => member.required)).toEqual([true, false]);
  });

  it('shows a long group no whole-row rule: each source carries its own missing result', () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [profileGroup('managed_long_v1', [0, 1])] };
    deliveryQuery.data = deliveryOf(sourceRecord({ status: 'missing', reason: 'no-samples', raw_present: false, delivery: 'queued' }));
    renderPanel(sourcesState(3));

    expect(screen.queryByTestId('circuit-profile-required-policy')).toBeNull();
    expect(screen.queryByTestId('basic-recording-completeness')).toBeNull();
    expandSources('circuit-status-sources-grp-1');
    expect(screen.getByTestId('circuit-status-record-latest-tag-0')).toHaveTextContent('step4.basic.status.status_value.missing');
    expect(screen.getByTestId('circuit-status-record-latest-tag-0')).toHaveTextContent('step4.basic.status.reason.no-samples');
  });
});

describe('evidence stays honest', () => {
  it('shows a mean as counts, method and unknown coverage: never complete and never a percentage', () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [profileGroup('managed_long_v1', [0, 1])] };
    const mean: WriteGroupSourceRecord & { aggregation: AggregationEvidence } = {
      ...sourceRecord(), aggregation: {
        method: 'arithmetic-mean', version: 'arithmetic-mean-v1', basis: 'exact', source_type: 'uint64', result_type: 'decimal',
        valid_count: 2, rejected_count: 0, coverage: 'unknown',
      },
    };
    deliveryQuery.data = deliveryOf(mean);
    renderPanel(sourcesState(3));

    expandSources('circuit-status-sources-grp-1');
    const view = screen.getByTestId('circuit-status-aggregation-tag-0');
    expect(view).toHaveTextContent('step4.basic.status.method.arithmetic-mean');
    expect(view).toHaveTextContent('step4.basic.status.coverage.unknown');
    expect(screen.getByTestId('circuit-status-valid-tag-0')).toHaveTextContent('2');
    expect(screen.getByTestId('circuit-status-rejected-tag-0')).toHaveTextContent('0');
    expect(view.textContent).not.toMatch(/complete|%/);
  });

  it('shows a counter baseline as a first reading with an unknown difference, not zero and not a failure', () => {
    groupsQuery.data = { ...groupsQuery.data, groups: [profileGroup('managed_long_v1', [0, 1], {
      recording_profile: { kind: 'cumulative', interval: 10, unit: 'minute' },
    })] };
    deliveryQuery.data = deliveryOf(sourceRecord({ profile_kind: 'cumulative', status: 'baseline', reason: 'first-reading', delta_known: false }));
    renderPanel(sourcesState(3));

    expandSources('circuit-status-sources-grp-1');
    const record = screen.getByTestId('circuit-status-record-latest-tag-0');
    expect(record).toHaveTextContent('step4.basic.status.status_value.baseline');
    expect(screen.getByTestId('circuit-status-delta-tag-0')).toHaveTextContent('step4.basic.status.difference_unknown');
    expect(screen.getByTestId('circuit-status-delta-tag-0')).not.toHaveTextContent(/^0$|zero/);
    expect(screen.queryByRole('alert')).toBeNull();
  });
});

describe('recording windows', () => {
  it('explains that acquisition feeds the independent recording windows and is not a report', () => {
    renderPanel(sourcesState(3));
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));

    const explainer = screen.getByTestId('circuit-profile-interval-explainer');
    expect(explainer).toHaveTextContent('step4.basic.circuit.interval_explainer');
    expect(explainer).toHaveTextContent('step4.basic.circuit.interval_not_report');
  });

  it('keeps exactly one period control per profile in Basic', () => {
    renderPanel(sourcesState(3));
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));

    expect(screen.getAllByTestId(/^circuit-profile-interval-(instantaneous|cumulative)$/)).toHaveLength(2);
    expect(screen.queryByTestId('basic-recording-interval')).toBeNull();
  });
});
