import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { basicRecordingScopeKey, saveBasicRecordingIntent } from '@/features/datalink/workbench-v2/state/basicRecording';
import { loadCircuitProfileLocalDraft } from '@/features/datalink/workbench-v2/state/circuitProfileBasic';
import { saveActiveCircuitId } from '@/features/datalink/workbench-v2/state/circuitProfileBatch';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { CircuitProfileBasicPanel } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel';
import type { RecordingStartOperation } from '@/types/studioV2RecordingStart';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { batchPanelState, httpFailure, longProfileGroup } from '../../fixtures/circuitBatchPanel';

const { query, measurements, ensure, update, start, operationQuery } = vi.hoisted(() => ({
  query: { data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] as WriteGroup[] }, isLoading: false, isError: false, refetch: vi.fn() },
  measurements: { data: [] as unknown[], isLoading: false, isError: false, refetch: vi.fn() },
  ensure: vi.fn(), update: vi.fn(), start: vi.fn(),
  operationQuery: { data: undefined as RecordingStartOperation | undefined, refetch: vi.fn() },
}));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
  useWriteGroupsQuery: () => query,
  useWriteGroupDeliveryQuery: () => ({ data: undefined, isError: false, isLoading: false }),
  useEnsureCircuitProfileMutation: () => ({ mutateAsync: ensure, isPending: false }),
  useUpdateWriteGroupMutation: () => ({ mutateAsync: update, isPending: false }),
}));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useRecordingStartMutation: () => ({ mutateAsync: start, isPending: false }),
  useRecordingStartOperationQuery: () => operationQuery,
}));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useStudioV2WorkspaceMeasurementsQuery: () => measurements,
  useCreateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn() }),
  useUpdateStudioV2MeasurementMutation: () => ({ mutateAsync: vi.fn() }),
}));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel', () => ({
  GroupSchemaPanel: () => <div data-testid="schema-panel" />,
}));

const state = batchPanelState();
let serverGroups: WriteGroup[];
let serverRevision: number;

type PutRequest = { expected_workspace_revision: string; group: object };

function panel(workspaceId = 'ws-1', over: Partial<React.ComponentProps<typeof CircuitProfileBasicPanel>> = {}, nextState = state) {
  return <CircuitProfileBasicPanel state={nextState} workspaceId={workspaceId} readonly={false} {...over} />;
}
function withConnectorRevision(revision: string) {
  return { ...state, db: { ...state.db, connector: { ...state.db.connector, identity_revision: revision } } };
}
/** The first PUT is held; later ones follow an in-memory server with the real CAS check. */
function holdFirstUpdate(options: { failSecondWith?: unknown } = {}) {
  let releaseFirst!: () => void;
  let calls = 0;
  update.mockImplementation(async ({ id, request }: { id: string; request: PutRequest }) => {
    calls += 1;
    if (calls === 1) await new Promise<void>((resolve) => { releaseFirst = resolve; });
    if (calls === 2 && options.failSecondWith) throw options.failSecondWith;
    if (request.expected_workspace_revision !== `wrev-${serverRevision}`) throw httpFailure(409, 'WRITE_GROUP_CIRCUIT_PROFILE_CONFLICT');
    const index = serverGroups.findIndex((group) => group.id === id);
    serverRevision += 1;
    serverGroups[index] = { ...serverGroups[index], ...(request.group as WriteGroup), id, revision: `rev-${id}-${serverRevision}` };
    return { workspace_revision: `wrev-${serverRevision}`, group: serverGroups[index] };
  });
  return { release: () => releaseFirst(), started: () => calls };
}
function confirmBatch(skip?: string) {
  fireEvent.click(screen.getByTestId('circuit-profile-open-batch'));
  if (skip) fireEvent.click(screen.getByTestId(`circuit-profile-batch-select-${skip}`));
  fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
  fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm'));
}
/** The pending-reload hook lets a test look at the new scope between the PUTs and the parent refresh. */
function holdReload() {
  let finish!: () => void;
  query.refetch.mockImplementation(() => new Promise((resolve) => { finish = () => resolve({ data: query.data }); }));
  return { finish: () => finish() };
}
function switchToWorkspaceTwo() {
  query.data = { workspace_id: 'ws-2', workspace_revision: 'wrev-1', groups: serverGroups.map((group) => ({ ...group })) };
}

beforeEach(() => {
  vi.clearAllMocks();
  update.mockReset();
  window.sessionStorage.clear();
  serverRevision = 1;
  serverGroups = [
    longProfileGroup({ id: 'grp-a', circuitId: 'cir-a', circuitName: 'Feeder A', interval: 15, points: [0, 1] }),
    longProfileGroup({ id: 'grp-b', circuitId: 'cir-b', circuitName: 'Feeder B', interval: 60, points: [2, 3] }),
    longProfileGroup({ id: 'grp-c', circuitId: 'cir-c', circuitName: 'Feeder C', interval: 15, points: [4, 5] }),
  ];
  query.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: serverGroups };
  query.refetch.mockImplementation(async () => ({ data: query.data }));
  measurements.data = [];
  ensure.mockImplementation(async ({ request }: { request: PutRequest }) =>
    ({ workspace_revision: request.expected_workspace_revision, group: serverGroups[1] }));
  start.mockResolvedValue({ operation_id: 'op-1', status: 'running', devices: [], groups: [] });
  operationQuery.data = undefined;
});

describe('A batch response belongs to the scope that launched it', () => {
  it('leaves the new scope draft, revision and busy state alone when the old workspace answers during a pending reload', async () => {
    saveActiveCircuitId('ws-2', 'dev-1', 'cir-b');
    const held = holdFirstUpdate();
    const view = render(panel());
    confirmBatch('cir-b');
    await waitFor(() => expect(held.started()).toBe(1));

    const reload = holdReload();
    switchToWorkspaceTwo();
    view.rerender(panel('ws-2'));
    await waitFor(() => expect(screen.getByTestId('circuit-profile-circuit-chooser')).toHaveValue('cir-b'));
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('60');
    expect(screen.queryByTestId('circuit-profile-batch-editor')).not.toBeInTheDocument();

    await act(async () => held.release());
    await waitFor(() => expect(update).toHaveBeenCalledTimes(2));
    await waitFor(() => expect(query.refetch).toHaveBeenCalled());

    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('60');
    expect(screen.getByTestId('circuit-profile-member-inst-tag-2')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-member-inst-tag-3')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-member-inst-tag-0')).not.toBeChecked();
    expect(screen.getByTestId('circuit-profile-member-inst-tag-1')).not.toBeChecked();
    expect(screen.getByTestId('circuit-profile-start')).toBeDisabled();

    await act(async () => reload.finish());
    await waitFor(() => expect(screen.getByTestId('circuit-profile-start')).toBeEnabled());

    fireEvent.click(screen.getByTestId('circuit-profile-start'));
    await waitFor(() => expect(ensure).toHaveBeenCalled());
    expect(ensure.mock.calls[0][0].request.expected_workspace_revision).toBe('wrev-1');
    expect(screen.queryByTestId('circuit-profile-error')).not.toBeInTheDocument();
  });

  it('keeps the old scope result in its own stored draft and shows no old error in the new scope', async () => {
    saveActiveCircuitId('ws-2', 'dev-1', 'cir-b');
    const held = holdFirstUpdate({ failSecondWith: httpFailure(409, 'WRITE_GROUP_CIRCUIT_PROFILE_CONFLICT') });
    const view = render(panel());
    confirmBatch('cir-b');
    await waitFor(() => expect(held.started()).toBe(1));

    switchToWorkspaceTwo();
    view.rerender(panel('ws-2'));
    await act(async () => held.release());
    await waitFor(() => expect(update).toHaveBeenCalledTimes(2));
    await act(async () => { await new Promise((resolve) => setTimeout(resolve, 0)); });

    expect(screen.queryByTestId('circuit-profile-error')).not.toBeInTheDocument();
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    const candidates = buildGroupCandidates(state).candidates.filter((candidate) => candidate.device_id === 'dev-1');
    const savedDraft = loadCircuitProfileLocalDraft('ws-1', 'dev-1', candidates, serverGroups, 'cir-a');
    expect(savedDraft.instantaneous_interval_input).toBe('30');
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('60');
  });

  it('closes the reviewed preview and never sends it when the destination changes before confirmation', () => {
    const view = render(panel());
    fireEvent.click(screen.getByTestId('circuit-profile-open-batch'));
    fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
    expect(screen.getByTestId('circuit-profile-batch-preview')).toBeInTheDocument();

    view.rerender(panel('ws-1', {}, withConnectorRevision('crev-2')));
    expect(screen.queryByTestId('circuit-profile-batch-editor')).not.toBeInTheDocument();
    expect(update).not.toHaveBeenCalled();
  });
});

describe('The page sees one operation-active signal for batch and start', () => {
  it('stays active when the old batch settles while the new scope still has an unresolved start', async () => {
    const newScope = basicRecordingScopeKey('ws-1', 'dev-1', encodeURIComponent(JSON.stringify(['cir-a', 'conn-1', 'crev-2', ''])));
    saveBasicRecordingIntent({
      scope_key: newScope,
      request: { request_id: 'req-1', workspace_id: 'ws-1', expected_workspace_revision: 'wrev-1', device_ids: ['dev-1'], groups: [] },
    });
    const held = holdFirstUpdate();
    const active = vi.fn();
    const view = render(panel('ws-1', { onOperationActiveChange: active }));
    confirmBatch();
    await waitFor(() => expect(held.started()).toBe(1));
    expect(active).toHaveBeenLastCalledWith(true);

    view.rerender(panel('ws-1', { onOperationActiveChange: active }, withConnectorRevision('crev-2')));
    expect(await screen.findByTestId('basic-recording-retry')).toBeInTheDocument();
    await act(async () => held.release());
    await waitFor(() => expect(update).toHaveBeenCalledTimes(3));
    await act(async () => { await new Promise((resolve) => setTimeout(resolve, 0)); });

    expect(active).toHaveBeenLastCalledWith(true);
    expect(active.mock.calls.slice(-1)).toEqual([[true]]);
  });
});

describe('The parent reload is awaited and its outcome is reported on its own', () => {
  const unhandled: unknown[] = [];
  const track = (reason: unknown) => { unhandled.push(reason); };
  beforeEach(() => { unhandled.length = 0; process.on('unhandledRejection', track); });
  afterEach(() => { process.off('unhandledRejection', track); });

  it('shows a failed reload after committed saves without a mutation error, an unhandled rejection or a resend', async () => {
    const held = holdFirstUpdate();
    query.refetch.mockRejectedValue(new Error('GET failed dsn=postgres://user:secret@db/prod'));
    const view = render(panel());
    confirmBatch();
    await waitFor(() => expect(held.started()).toBe(1));
    await act(async () => held.release());

    expect(await screen.findByTestId('circuit-profile-batch-refresh-error')).toBeInTheDocument();
    await act(async () => { await new Promise((resolve) => setTimeout(resolve, 0)); });
    expect(unhandled).toEqual([]);
    expect(update).toHaveBeenCalledTimes(3);
    expect(view.container.textContent).not.toMatch(/secret|postgres:/);
    expect(screen.queryByTestId('circuit-profile-error')).not.toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-batch-status-cir-a-instantaneous')).toHaveTextContent('step4.basic.circuit.batch_status_saved');
  });

  it('checks an unknown outcome against the reloaded groups and syncs the draft only once the server proves it', async () => {
    update.mockRejectedValueOnce(new Error('socket hang up'));
    const view = render(panel());
    confirmBatch();
    await screen.findByTestId('circuit-profile-batch-unknown-note');
    expect(update).toHaveBeenCalledTimes(1);
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('15');

    fireEvent.click(screen.getByTestId('circuit-profile-batch-verify'));
    await waitFor(() => expect(query.refetch).toHaveBeenCalledTimes(1));
    view.rerender(panel());
    expect(screen.getByTestId('circuit-profile-batch-status-cir-a-instantaneous')).toHaveTextContent('step4.basic.circuit.batch_status_unknown');

    serverGroups[0] = {
      ...serverGroups[0], revision: 'rev-grp-a-9', recording_profile: { kind: 'instantaneous', interval: 30, unit: 'second' },
      row_policy: { ...serverGroups[0].row_policy, interval_seconds: 30 },
    };
    query.data = { ...query.data, groups: [...serverGroups] };
    fireEvent.click(screen.getByTestId('circuit-profile-batch-verify'));
    await waitFor(() => expect(query.refetch).toHaveBeenCalledTimes(2));
    view.rerender(panel());
    await waitFor(() => expect(screen.getByTestId('circuit-profile-batch-status-cir-a-instantaneous')).toHaveTextContent('step4.basic.circuit.batch_status_saved'));
    await waitFor(() => expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('30'));
    expect(update).toHaveBeenCalledTimes(1);
  });
});

describe('A batch edits the period and leaves the operator source choices alone', () => {
  function uncheckOneSourceOfCircuitA() {
    fireEvent.click(screen.getByTestId('circuit-profile-member-inst-tag-1'));
    expect(screen.getByTestId('circuit-profile-member-inst-tag-1')).not.toBeChecked();
  }
  const sentGroup = () => (update.mock.calls[0][0] as { request: { group: object } }).request.group as WriteGroup;

  it('keeps a source the operator unchecked after the period was saved', async () => {
    const held = holdFirstUpdate();
    render(panel());
    uncheckOneSourceOfCircuitA();
    confirmBatch();
    await waitFor(() => expect(held.started()).toBe(1));
    await act(async () => held.release());
    await waitFor(() => expect(screen.queryByTestId('circuit-profile-batch-editor')).not.toBeInTheDocument());

    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('30');
    expect(screen.getByTestId('circuit-profile-member-inst-tag-0')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-member-inst-tag-1')).not.toBeChecked();
  });

  it('keeps that choice when an unknown save is later proven by the reloaded group', async () => {
    update.mockRejectedValueOnce(new Error('socket hang up'));
    const view = render(panel());
    uncheckOneSourceOfCircuitA();
    confirmBatch();
    await screen.findByTestId('circuit-profile-batch-unknown-note');

    // The server did commit exactly what was sent.
    serverGroups[0] = { ...sentGroup(), id: 'grp-a', revision: 'rev-grp-a-9' };
    query.data = { ...query.data, groups: [...serverGroups] };
    fireEvent.click(screen.getByTestId('circuit-profile-batch-verify'));
    await waitFor(() => expect(query.refetch).toHaveBeenCalled());
    view.rerender(panel());
    await waitFor(() => expect(screen.getByTestId('circuit-profile-batch-status-cir-a-instantaneous')).toHaveTextContent('step4.basic.circuit.batch_status_saved'));
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('30');
    expect(screen.getByTestId('circuit-profile-member-inst-tag-1')).not.toBeChecked();
    expect(update).toHaveBeenCalledTimes(1);
  });

  it('does not report an unknown save as saved when the reloaded group lost a source', async () => {
    update.mockRejectedValueOnce(new Error('socket hang up'));
    const view = render(panel());
    confirmBatch();
    await screen.findByTestId('circuit-profile-batch-unknown-note');

    const sent = sentGroup();
    serverGroups[0] = { ...sent, id: 'grp-a', revision: 'rev-grp-a-9', members: sent.members.slice(0, 1) };
    query.data = { ...query.data, groups: [...serverGroups] };
    fireEvent.click(screen.getByTestId('circuit-profile-batch-verify'));
    await waitFor(() => expect(query.refetch).toHaveBeenCalled());
    view.rerender(panel());
    expect(screen.getByTestId('circuit-profile-batch-status-cir-a-instantaneous')).toHaveTextContent('step4.basic.circuit.batch_status_unknown');
    expect(screen.getByTestId('circuit-profile-batch-unknown-note')).toBeInTheDocument();
  });
});
