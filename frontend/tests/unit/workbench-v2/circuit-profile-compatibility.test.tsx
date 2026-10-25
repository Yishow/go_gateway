import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import type { WorkbenchV2State } from '@/features/datalink/workbench-v2/state/types';
import { studioV2RecordingStartAPI } from '@/services/studioV2RecordingStart';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { ModbusShareStatus } from '@/types/modbusShare';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { batchPanelState, longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { savedGroup, TABLE_COLUMNS } from '../../fixtures/writeGroupState';
import { expandSources } from '../../fixtures/sourceReview';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string, options?: unknown) => (options && typeof options === 'object' && 'count' in options ? `${key}:${(options as { count: number }).count}` : key) }),
}));
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({
  useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }),
}));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: {
    list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(),
    testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), get: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(),
    schemaOperation: vi.fn(), ensureCircuitProfile: vi.fn(),
  },
}));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { list: vi.fn(async () => []) } }));
vi.mock('@/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), operation: vi.fn(), ensureBasic: vi.fn() } }));

const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const start = vi.mocked(studioV2RecordingStartAPI);
const quietDelivery = {
  group_id: 'x', intake: { state: 'not_running' }, stages: { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 },
  last_sql_committed_at: null, oldest_pending_seconds: 0, no_data_buckets: 0, skipped_buckets: 0, backlog: [],
  quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
};
/** Every call that would change a saved group, apply, disable, remove or start recording. Previews read only. */
const WRITES = [api.create, api.update, api.apply, api.disable, api.remove, api.ensureCircuitProfile, api.schemaApply, start.start, start.ensureBasic];
const expectNoWrites = () => WRITES.forEach((write) => expect(write).not.toHaveBeenCalled());

const power = (over: Partial<WriteGroup> = {}) => ({ ...longProfileGroup({ id: 'grp-i', circuitId: 'cir-a', circuitName: 'Feeder A', interval: 45, points: [0, 1] }), ...over });
const energy = (over: Partial<WriteGroup> = {}) => ({ ...longProfileGroup({ id: 'grp-c', circuitId: 'cir-a', circuitName: 'Feeder A', interval: 20, points: [0, 1], kind: 'cumulative' }), ...over });
/** The group the legacy Basic path owns: managed wide layout, one device, one row per interval. */
const legacyWide = (interval = 60, over: Partial<WriteGroup> = {}): WriteGroup => savedGroup({
  id: 'wide-1', name: 'PLC A raw', basic_managed_device_id: 'dev-1', layout_version: 'managed_wide_v1', status: 'ready', applied_revision: 'applied-w',
  destination: { ...savedGroup().destination, storage_strategy: 'managed', table_name: 'gw_basic_dev_1' },
  row_policy: { interval_seconds: interval, allowed_lateness_seconds: 0, incomplete_policy: 'skip_row' }, ...over,
});
const customGroup = (): WriteGroup => savedGroup({
  id: 'custom-1', name: 'Plant table', row_policy: { interval_seconds: 10, allowed_lateness_seconds: 0, incomplete_policy: 'partial', entity_key_column: 'entity', provenance_column: 'prov' },
  members: [
    { device_id: 'dev-1', point_id: 'pt-4', tag_id: 'tag-4', source_revision: 's4', mapping_revision: 'm4', target_column: 'temperature', required: true, entity_key: 'line-1' },
    { device_id: 'dev-2', point_id: 'pt-9', tag_id: 'tag-9', source_revision: 's9', mapping_revision: 'm9', target_column: 'pressure', required: false, entity_key: 'line-2' },
  ],
});

function renderHost(groups: WriteGroup[], over: { state?: WorkbenchV2State; shareStatus?: ModbusShareStatus } = {}) {
  api.list.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const element = () => (
    <QueryClientProvider client={client}>
      <Step4Database state={over.state ?? batchPanelState()} dispatch={vi.fn()} workspaceId="ws-1" shareStatus={over.shareStatus} />
    </QueryClientProvider>
  );
  const view = render(element());
  return { ...view, reload: () => { view.unmount(); return render(element()); } };
}

beforeEach(() => {
  vi.clearAllMocks();
  window.sessionStorage.clear();
  api.delivery.mockResolvedValue(quietDelivery as never);
  api.readiness.mockResolvedValue({ workspace_id: 'ws-1', group_id: 'x', group_revision: 'r', applied_revision: '', config_ready: true, schema_ready: true, ready: true, issues: [] } as never);
  api.schemaPreview.mockResolvedValue({ operation_id: 'op', token: 'tok', statements: [], tables: [], table_prefix: 'gw', no_change_reason: 'schema_compatible' } as never);
});

describe('A new Basic setup has one effective period: the one of the profile', () => {
  it('offers profile kinds with native periods and no legacy seconds, completeness or required controls', async () => {
    renderHost([]);
    expect(await screen.findByTestId('circuit-profile-basic-panel')).toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-interval-cumulative')).toBeInTheDocument();
    for (const id of ['basic-recording-panel', 'basic-recording-interval', 'basic-recording-completeness', 'basic-required-policy']) expect(screen.queryByTestId(id), id).not.toBeInTheDocument();
    expect(screen.queryByLabelText('step4.basic.interval')).not.toBeInTheDocument();
    expect(screen.queryByLabelText('step4.basic.completeness')).not.toBeInTheDocument();
  });

  it('keeps the connection card to connection identity: no table, write mode or period control', async () => {
    renderHost([]);
    const card = await screen.findByTestId('database-connection-card');
    expect(within(card).queryAllByLabelText(/table|write_mode|upsert|insert|interval|period/i)).toHaveLength(0);
    expect(card.textContent).not.toMatch(/upsert|insert|step4\.field_table|step4\.field_interval|write_mode/i);
  });

  it('shows the saved native period of each profile rather than a Basic default', async () => {
    renderHost([power(), energy()]);
    expect(await screen.findByTestId('circuit-profile-interval-instantaneous')).toHaveValue('45');
    expect(screen.getByTestId('circuit-profile-interval-cumulative')).toHaveValue('20');
    expectNoWrites();
  });
});

describe('Basic, Advanced and Basic again show one setting and write nothing', () => {
  it('returns to exactly the same profile, period, mode and sources after a visit to Advanced and after a reload', async () => {
    const view = renderHost([power(), energy()]);
    await screen.findByTestId('circuit-profile-basic-panel');
    const snapshot = () => ({
      name: (screen.getByTestId('circuit-profile-circuit-name') as HTMLInputElement).value,
      circuit: (screen.getByTestId('circuit-profile-circuit-chooser') as HTMLSelectElement).value,
      inst: (screen.getByTestId('circuit-profile-interval-instantaneous') as HTMLInputElement).value,
      cum: (screen.getByTestId('circuit-profile-interval-cumulative') as HTMLInputElement).value,
      toggles: [screen.getByTestId('circuit-profile-toggle-instantaneous'), screen.getByTestId('circuit-profile-toggle-cumulative')].map((el) => (el as HTMLInputElement).checked),
      sources: screen.getAllByTestId(/^circuit-profile-member-(inst|cum)-/).map((el) => `${el.getAttribute('data-testid')}:${(el as HTMLInputElement).checked}`),
      mode: (screen.getByTestId('circuit-profile-storage-plain') as HTMLInputElement).checked,
    });
    expandSources();
    const before = snapshot();

    fireEvent.click(screen.getByTestId('circuit-status-open-advanced-instantaneous'));
    expect(await screen.findByTestId('group-profile-interval')).toHaveTextContent('45');
    expect(screen.getByTestId('group-profile-storage')).toHaveTextContent('step4.group.profile.storage.plain');
    fireEvent.click(screen.getByTestId('group-editor-close'));
    expect(snapshot()).toEqual(before);

    view.reload();
    await screen.findByTestId('circuit-profile-basic-panel');
    expandSources();
    expect(snapshot()).toEqual(before);
    expectNoWrites();
  });
});

describe('Existing wide, custom and share-only setups keep their own editing and runtime', () => {
  it('leaves a legacy 60 second Basic group on the legacy panel, converts nothing and keeps its settings', async () => {
    renderHost([legacyWide()]);
    expect(await screen.findByTestId('basic-recording-panel')).toBeInTheDocument();
    expect(screen.queryByTestId('circuit-profile-basic-panel')).not.toBeInTheDocument();
    expect(screen.getByTestId('basic-recording-interval')).toHaveValue(60);
    expect(screen.getByTestId('basic-recording-completeness')).toBeInTheDocument();
    expect(screen.getByTestId('basic-required-policy')).toBeInTheDocument();
    expectNoWrites();
  });

  it('keeps a legacy period that the new instantaneous bounds would reject', async () => {
    renderHost([legacyWide(900)]);
    expect(await screen.findByTestId('basic-recording-panel')).toBeInTheDocument();
    fireEvent.click(await screen.findByTestId('basic-member-settings'));
    expect(await screen.findByTestId('group-interval')).toHaveValue(900);
    expect(screen.getByTestId('group-interval')).toBeEnabled();
    expect(screen.queryByTestId('group-profile-summary')).not.toBeInTheDocument();
    expect(screen.queryByTestId('group-issues')).not.toBeInTheDocument();
    expectNoWrites();
  });

  it('keeps the legacy panel and its 60 second settings when prepared long targets already exist beside it', async () => {
    renderHost([legacyWide(), power({ applied_revision: '', status: 'draft' })]);
    expect(await screen.findByTestId('basic-recording-panel')).toBeInTheDocument();
    expect(screen.queryByTestId('circuit-profile-basic-panel')).not.toBeInTheDocument();
    expect(screen.getByTestId('basic-recording-interval')).toHaveValue(60);
    expect(screen.getByTestId('legacy-migration-open')).toBeEnabled();
    expectNoWrites();
  });

  it('opens the legacy wide group in the original editor with seconds, completeness, required and column controls', async () => {
    renderHost([legacyWide()]);
    fireEvent.click(await screen.findByTestId('basic-member-settings'));
    await screen.findByTestId('group-editor');
    expect(screen.getByTestId('group-name')).toHaveValue('PLC A raw');
    expect(screen.getByTestId('group-interval')).toHaveValue(60);
    expect(screen.getByTestId('group-incomplete')).toHaveValue('skip_row');
    expect(screen.getAllByLabelText(/step4\.group\.members\.required_label/).length).toBeGreaterThan(0);
    expect(screen.getAllByTestId(/^group-member-column-/).length).toBeGreaterThan(0);
    expect(screen.getByTestId('group-storage-managed')).toBeEnabled();
    expectNoWrites();
  });

  it('does not turn a custom multi-entity cross-device group into a profile or touch it', async () => {
    renderHost([customGroup()]);
    expect(await screen.findByTestId('basic-recording-advanced-note')).toHaveTextContent('1');
    fireEvent.click(await screen.findByTestId('group-open-custom-1'));
    await screen.findByTestId('group-editor');
    expect(screen.getByTestId('group-interval')).toHaveValue(10);
    expect(screen.getByTestId('group-incomplete')).toHaveValue('partial');
    expect(screen.getByTestId('group-entity-column')).toHaveValue('entity');
    expect(screen.getByTestId('group-provenance-column')).toHaveValue('prov');
    expect(screen.getByTestId('group-storage-custom')).toBeChecked();
    expect(screen.queryByTestId('group-profile-summary')).not.toBeInTheDocument();
    expectNoWrites();
  });

  it('creates no database profile and asks for no probe when only Modbus Share is configured', async () => {
    const state = batchPanelState();
    const unsaved = { ...state, db: { ...state.db, connector: { ...state.db.connector, persisted: false, save_state: 'idle' as const, connector_id: undefined, identity_revision: undefined } } };
    renderHost([], { state: unsaved, shareStatus: { enabled: true, configured_enabled: true, port: 5020, address: '127.0.0.1', bind_state: 'pass', mapping_count: 3 } });
    expect(await screen.findByTestId('basic-recording-panel')).toBeInTheDocument();
    expect(screen.getByTestId('basic-recording-destination')).toHaveTextContent('step4.basic.share_only_destination');
    expect(screen.queryByTestId('circuit-profile-basic-panel')).not.toBeInTheDocument();
    await waitFor(() => expect(api.list).toHaveBeenCalled());
    expectNoWrites();
    expect(api.schemaPreview).not.toHaveBeenCalled();
  });
});

describe('Unchecking a profile in Basic changes the draft, never a running writer', () => {
  it('keeps the applied writer, says so, and leaves stopping to the explicit group lifecycle', async () => {
    renderHost([power({ applied_revision: 'applied-i', status: 'running' }), energy()]);
    fireEvent.click(await screen.findByTestId('circuit-profile-toggle-instantaneous'));
    const note = await screen.findByTestId('circuit-profile-draft-note-instantaneous');
    expect(note).toHaveTextContent('step4.basic.circuit.draft_not_stop');
    expectNoWrites();
    fireEvent.click(within(note).getByTestId('circuit-profile-draft-lifecycle-instantaneous'));
    expect(await screen.findByTestId('group-name')).toHaveValue('Feeder A (instantaneous)');
    expect(screen.getByTestId('group-disable')).toBeInTheDocument();
    expectNoWrites();
  });

  it('shows no such note for a profile that was never applied', async () => {
    renderHost([power(), energy()]);
    fireEvent.click(await screen.findByTestId('circuit-profile-toggle-instantaneous'));
    expect(screen.queryByTestId('circuit-profile-draft-note-instantaneous')).not.toBeInTheDocument();
  });
});
