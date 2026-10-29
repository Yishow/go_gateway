import * as React from 'react';
import { QueryClient, QueryClientProvider, useQueryClient } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { clickEnabled, guidedLifecycle, guidedReadiness } from '../../fixtures/guidedRecordingPreparation';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { basicRecordingScopeKey, saveBasicRecordingIntent } from '@/features/datalink/workbench-v2/state/basicRecording';
import { loadCircuitProfileLocalDraft, saveCircuitProfileLocalDraft } from '@/features/datalink/workbench-v2/state/circuitProfileBasic';
import { saveActiveCircuitId } from '@/features/datalink/workbench-v2/state/circuitProfileBatch';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { CircuitProfileBasicPanel } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel';
import type { RecordingStartOperation } from '@/types/studioV2RecordingStart';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { batchPanelState, httpFailure, longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { expandSources } from '../../fixtures/sourceReview';

const { query, measurements, ensure, update, start, operationQuery } = vi.hoisted(() => ({
  query: { data: { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [] as WriteGroup[] }, isLoading: false, isError: false, refetch: vi.fn() },
  measurements: { data: [] as unknown[], isLoading: false, isError: false, refetch: vi.fn() },
  ensure: vi.fn(), update: vi.fn(), start: vi.fn(),
  operationQuery: { data: undefined as RecordingStartOperation | undefined, refetch: vi.fn() },
}));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
let client: QueryClient;
const listeners = new Set<() => void>();
function publishGroups() { client.setQueryData(studioV2WorkspaceKeys.writeGroups(), query.data); for (const listener of listeners) listener(); }
function useCanonicalGroups() {
  const cache = useQueryClient(), [, refresh] = React.useReducer(value => value + 1, 0), accepted = query.data;
  React.useLayoutEffect(() => { cache.setQueryData(studioV2WorkspaceKeys.writeGroups(), accepted); }, [cache, accepted]);
  React.useEffect(() => { const listener = () => refresh(); listeners.add(listener); return () => { listeners.delete(listener); }; }, []);
  return query;
}
function canonical(id: string) { const group = query.data.groups.find(item => item.id === id); if (!group) throw new Error('Missing canonical group'); return group; }
async function nativeUpdate(input: { id: string; request: { workspace_id: string } }) {
  const result = await update(input);
  // Model native invalidation publication only for the still-selected workspace; never borrow old-head authority.
  if (input.request.workspace_id === query.data.workspace_id) {
    query.data = { workspace_id: input.request.workspace_id, workspace_revision: result.workspace_revision,
      groups: query.data.groups.map(group => group.id === input.id ? result.group : group) }; publishGroups();
  }
  return result;
}
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: vi.fn(async (id: string) => guidedLifecycle(canonical(id), query.data.workspace_revision)) } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { readiness: vi.fn(async (id: string) => guidedReadiness(canonical(id), query.data.workspace_revision)), schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn() } }));
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', async importOriginal => ({
  ...await importOriginal<typeof import('@/hooks/datalink/useStudioV2WriteGroups')>(),
  useWriteGroupsQuery: useCanonicalGroups,
  useWriteGroupDeliveryQuery: () => ({ data: undefined, isError: false, isLoading: false }),
  useEnsureCircuitProfileMutation: () => ({ mutateAsync: ensure, isPending: false }),
  useUpdateWriteGroupMutation: () => ({ mutateAsync: nativeUpdate, isPending: false }),
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

const state = batchPanelState();
let serverGroups: WriteGroup[];
let serverRevision: number;

type PutRequest = { expected_workspace_revision: string; group: object };

function panel(workspaceId = 'ws-1', over: Partial<React.ComponentProps<typeof CircuitProfileBasicPanel>> = {}, nextState = state) {
  return <QueryClientProvider client={client}><CircuitProfileBasicPanel state={nextState} workspaceId={workspaceId} readonly={false} {...over} /></QueryClientProvider>;
}
function withConnectorRevision(revision: string) {
  return { ...state, db: { ...state.db, connector: { ...state.db.connector, identity_revision: revision } } };
}
/** One explicit PUT is held; all targets follow an in-memory server with the real CAS check. */
function holdUpdate(options: { failSecondWith?: unknown; at?: number } = {}) {
  let releaseFirst!: () => void;
  let calls = 0;
  update.mockImplementation(async ({ id, request }: { id: string; request: PutRequest }) => {
    calls += 1;
    if (calls === (options.at ?? 1)) await new Promise<void>((resolve) => { releaseFirst = resolve; });
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
  query.data = { workspace_id: 'ws-2', workspace_revision: 'wrev-1', groups: serverGroups.map((group) => ({ ...group, workspace_id: 'ws-2' })) };
}

beforeEach(() => {
  vi.clearAllMocks();
  client = new QueryClient({ defaultOptions: { queries: { retry: false, structuralSharing: false }, mutations: { retry: false } } });
  update.mockReset();
  window.sessionStorage.clear();
  serverRevision = 1;
  serverGroups = [
    longProfileGroup({ id: 'grp-a', circuitId: 'cir-a', circuitName: 'Feeder A', interval: 15, points: [0, 1] }),
    longProfileGroup({ id: 'grp-b', circuitId: 'cir-b', circuitName: 'Feeder B', interval: 60, points: [2, 3] }),
    longProfileGroup({ id: 'grp-c', circuitId: 'cir-c', circuitName: 'Feeder C', interval: 15, points: [4, 5] }),
  ];
  query.data = { workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [...serverGroups] };
  query.refetch.mockImplementation(async () => ({ data: query.data }));
  measurements.data = [];
  ensure.mockImplementation(async ({ request }: { request: PutRequest }) =>
    ({ workspace_revision: request.expected_workspace_revision, group: serverGroups[1] }));
  start.mockRejectedValue(new Error('No start is authorized by a batch'));
  operationQuery.data = undefined;
});

describe('A batch response belongs to the scope that launched it', () => {
  it('leaves the new scope draft, revision and busy state alone when the old workspace answers during a pending reload', async () => {
    saveActiveCircuitId('ws-2', 'dev-1', 'cir-b');
    const held = holdUpdate();
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
    // Scoped command lease retired: retain accepted first save, abort every not-yet-dispatched target.
    await waitFor(() => expect(query.refetch).toHaveBeenCalled());
    expect(update).toHaveBeenCalledTimes(1);
    await waitFor(() => expect(query.refetch).toHaveBeenCalled());

    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('60');
    expandSources();
    expect(screen.getByTestId('circuit-profile-member-inst-tag-2')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-member-inst-tag-3')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-member-inst-tag-0')).not.toBeChecked();
    expect(screen.getByTestId('circuit-profile-member-inst-tag-1')).not.toBeChecked();
    expect(screen.getByTestId('recording-guided-preparation').querySelector('button[data-primary]')).toBeDisabled();

    await act(async () => reload.finish());
    await waitFor(() => expect(screen.getByTestId('group-schema-preview')).toBeEnabled());
    expect(query.data.workspace_revision).toBe('wrev-1');
    expect(client.getQueryData(studioV2WorkspaceKeys.writeGroups())).toBe(query.data);
    expect(serverGroups[0].recording_profile?.interval).toBe(30);
    expect(serverGroups[2].recording_profile?.interval).toBe(15);
    // New explicit save is scoped to B's current workspace/head, never an old batch continuation.
    update.mockImplementationOnce(async ({ id, request }) => ({ workspace_revision: 'wrev-2', group: { ...canonical(id), ...request.group, revision: 'ws2-group-2' } }));
    fireEvent.change(screen.getByTestId('circuit-profile-interval-instantaneous'), { target: { value: '45' } });
    await clickEnabled('circuit-profile-prepare');
    await waitFor(() => expect(update).toHaveBeenCalledTimes(2));
    expect(update.mock.calls[1][0].request).toMatchObject({ workspace_id: 'ws-2', expected_workspace_revision: 'wrev-1' });
    expect(ensure).not.toHaveBeenCalled(); expect(start).not.toHaveBeenCalled();
    expect(screen.queryByTestId('circuit-profile-error')).not.toBeInTheDocument();
  });

  it('keeps the old scope result in its own stored draft and shows no old error in the new scope', async () => {
    saveActiveCircuitId('ws-2', 'dev-1', 'cir-b');
    const held = holdUpdate({ failSecondWith: httpFailure(409, 'WRITE_GROUP_CIRCUIT_PROFILE_CONFLICT') });
    const view = render(panel());
    confirmBatch('cir-b');
    await waitFor(() => expect(held.started()).toBe(1));

    switchToWorkspaceTwo();
    view.rerender(panel('ws-2'));
    await act(async () => held.release());
    // Scoped command lease retired: retain accepted first save, abort every not-yet-dispatched target.
    await waitFor(() => expect(query.refetch).toHaveBeenCalled());
    expect(update).toHaveBeenCalledTimes(1);
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
    const held = holdUpdate();
    const active = vi.fn();
    const view = render(panel('ws-1', { onOperationActiveChange: active }));
    confirmBatch();
    await waitFor(() => expect(held.started()).toBe(1));
    expect(active).toHaveBeenLastCalledWith(true);

    view.rerender(panel('ws-1', { onOperationActiveChange: active }, withConnectorRevision('crev-2')));
    // Guide design owner-local recovery: no O1 means manual guidance, never a retry POST.
    await waitFor(() => expect(screen.getByTestId('recording-next-action')).toHaveTextContent('next.actions.recover_operation'));
    expect(screen.queryByTestId('basic-recording-retry')).not.toBeInTheDocument();
    await act(async () => held.release());
    await waitFor(() => expect(query.refetch).toHaveBeenCalled());
    expect(update).toHaveBeenCalledTimes(1);
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

  it('keeps the reviewed remaining targets after delayed own first success with a healthy canonical reload', async () => {
    const held = holdUpdate();
    render(panel()); confirmBatch();
    await waitFor(() => expect(held.started()).toBe(1));
    expect(update.mock.calls[0][0].request.expected_workspace_revision).toBe('wrev-1');
    await act(async () => held.release());
    // No native controls changed: accepted own progress is not new operator consent or an external head advance.
    await waitFor(() => expect(update).toHaveBeenCalledTimes(3));
    expect(update.mock.calls.map(([input]) => input.request.expected_workspace_revision)).toEqual(['wrev-1', 'wrev-2', 'wrev-3']);
    expect(serverGroups.map(group => group.recording_profile?.interval)).toEqual([30, 30, 30]);
    await waitFor(() => expect(screen.queryByTestId('circuit-profile-batch-editor')).not.toBeInTheDocument());
    expect(ensure).not.toHaveBeenCalled(); expect(start).not.toHaveBeenCalled();
  });

  it.each([['changed', ['21'], '21', false], ['ABA', ['21', '15'], '15', false], ['edit then navigate', ['21'], '21', true]] as const)('retires the reviewed batch after captured native interval %s without overwriting the latest edit', async (_case, values, expected, navigate) => {
    const held = holdUpdate(); const view = render(panel());
    const interval = screen.getByTestId('circuit-profile-interval-instantaneous');
    const key = Object.keys(interval).find(value => value.startsWith('__reactProps$'));
    if (!key) throw new Error('Missing mounted native input props');
    const captured = (interval as unknown as Record<string, { onChange: (event: { target: { value: string } }) => void }>)[key].onChange;
    confirmBatch(); await waitFor(() => expect(held.started()).toBe(1));
    const original = structuredClone(update.mock.calls[0][0]);
    expect(interval).toBeDisabled();
    for (const value of values) await act(async () => captured({ target: { value } }));
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue(expected);
    if (navigate) {
      saveActiveCircuitId('ws-2', 'dev-1', 'cir-b'); switchToWorkspaceTwo(); view.rerender(panel('ws-2'));
      await waitFor(() => expect(screen.getByTestId('circuit-profile-circuit-chooser')).toHaveValue('cir-b'));
      expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('60');
    }
    await act(async () => held.release());
    await waitFor(() => expect(query.refetch).toHaveBeenCalled());
    expect(update).toHaveBeenCalledTimes(1); expect(update.mock.calls[0][0]).toEqual(original);
    expect(original.request.group.recording_profile.interval).toBe(30);
    expect(serverGroups.map(group => group.recording_profile?.interval)).toEqual([30, 60, 15]);
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue(navigate ? '60' : expected);
    const candidates = buildGroupCandidates(state).candidates.filter(candidate => candidate.device_id === 'dev-1');
    expect(loadCircuitProfileLocalDraft('ws-1', 'dev-1', candidates, serverGroups, 'cir-a').instantaneous_interval_input).toBe(expected);
    expect(ensure).not.toHaveBeenCalled(); expect(start).not.toHaveBeenCalled();
  });

  it.each([['edit', ['21'], '21', false], ['ABA to old 60', ['21', '60'], '60', true], ['pure NAV', [], '30', true]] as const)('settles an accepted non-launch target after %s without losing actual edits or suppressing untouched progress', async (_case, values, expected, seed) => {
    const candidates = buildGroupCandidates(state).candidates.filter(candidate => candidate.device_id === 'dev-1');
    if (seed) saveCircuitProfileLocalDraft('ws-1', 'dev-1', loadCircuitProfileLocalDraft('ws-1', 'dev-1', candidates, serverGroups, 'cir-b'));
    const held = holdUpdate({ at: 3 }); render(panel());
    const chooser = screen.getByTestId('circuit-profile-circuit-chooser');
    const chooserKey = Object.keys(chooser).find(value => value.startsWith('__reactProps$'));
    if (!chooserKey) throw new Error('Missing mounted chooser props');
    const choose = (chooser as unknown as Record<string, { onChange: (event: { target: { value: string } }) => void }>)[chooserKey].onChange;
    confirmBatch(); await waitFor(() => expect(held.started()).toBe(3));
    expect(serverGroups.map(group => group.recording_profile?.interval)).toEqual([30, 30, 15]);
    const original = structuredClone(update.mock.calls.map(([input]) => input));
    await act(async () => choose({ target: { value: 'cir-b' } }));
    await waitFor(() => expect(screen.getByTestId('circuit-profile-circuit-chooser')).toHaveValue('cir-b'));
    const interval = screen.getByTestId('circuit-profile-interval-instantaneous');
    expect(interval).toHaveValue(seed ? '60' : '30');
    const inputKey = Object.keys(interval).find(value => value.startsWith('__reactProps$'));
    if (!inputKey) throw new Error('Missing mounted native input props');
    const edit = (interval as unknown as Record<string, { onChange: (event: { target: { value: string } }) => void }>)[inputKey].onChange;
    for (const value of values) await act(async () => edit({ target: { value } }));
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue(values.at(-1) ?? '60');
    await act(async () => held.release()); await waitFor(() => expect(query.refetch).toHaveBeenCalled());
    expect(update).toHaveBeenCalledTimes(3); expect(update.mock.calls.map(([input]) => input)).toEqual(original);
    expect(original.map(input => input.request.expected_workspace_revision)).toEqual(['wrev-1', 'wrev-2', 'wrev-3']);
    expect(serverGroups.map(group => group.recording_profile?.interval)).toEqual([30, 30, 30]);
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue(expected);
    expect(loadCircuitProfileLocalDraft('ws-1', 'dev-1', candidates, serverGroups, 'cir-b').instantaneous_interval_input).toBe(expected);
    expect(ensure).not.toHaveBeenCalled(); expect(start).not.toHaveBeenCalled();
  });

  it('never overwrites a remounted native owner edit when the old unmounted first save settles', async () => {
    const held = holdUpdate(); const old = render(panel()); confirmBatch();
    await waitFor(() => expect(held.started()).toBe(1));
    const original = structuredClone(update.mock.calls[0][0]);
    old.unmount(); render(panel());
    const interval = screen.getByTestId('circuit-profile-interval-instantaneous');
    await waitFor(() => expect(interval).toBeEnabled());
    fireEvent.change(interval, { target: { value: '21' } });
    expect(interval).toHaveValue('21');
    await act(async () => held.release()); await waitFor(() => expect(query.refetch).toHaveBeenCalled());
    expect(update).toHaveBeenCalledTimes(1); expect(update.mock.calls[0][0]).toEqual(original);
    expect(serverGroups.map(group => group.recording_profile?.interval)).toEqual([30, 60, 15]);
    expect(client.getQueryData(studioV2WorkspaceKeys.writeGroups())).toMatchObject({ workspace_revision: 'wrev-2' });
    expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('21');
    const candidates = buildGroupCandidates(state).candidates.filter(candidate => candidate.device_id === 'dev-1');
    expect(loadCircuitProfileLocalDraft('ws-1', 'dev-1', candidates, serverGroups, 'cir-a').instantaneous_interval_input).toBe('21');
    expect(ensure).not.toHaveBeenCalled(); expect(start).not.toHaveBeenCalled();
  });

  it('shows a failed reload after committed saves without a mutation error, an unhandled rejection or a resend', async () => {
    const held = holdUpdate();
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
    const held = holdUpdate();
    render(panel());
    expandSources();
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
    expandSources();
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
