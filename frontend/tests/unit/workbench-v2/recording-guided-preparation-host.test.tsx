import { onlineManager, QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { CircuitProfileBasicPanel } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel';
import { studioV2WorkspaceWriteGroupsAPI as groupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2RecordingStartAPI as startAPI } from '@/services/studioV2RecordingStart';
import type { SchemaOperation, SchemaPreviewToken } from '@/types/recordingPlan';
import type { WriteGroup, WriteGroupCreateRequest, WriteGroupSchemaApplyRequest, WriteGroupSchemaPreviewRequest } from '@/types/studioV2WriteGroup';
import type { RecordingStartOperation, RecordingStartRequest } from '@/types/studioV2RecordingStart';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import { expandSources } from '../../fixtures/sourceReview';
import { lifecycle } from '../../fixtures/writeGroupSuccessor';
import { recordingTestWriteRecovery as recovery } from '@/features/datalink/workbench-v2/state/recordingTestWriteRecovery';
import { retainSchemaOperation } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery';
import type { WorkbenchV2State } from '@/features/datalink/workbench-v2/state/types';
import type { WriteGroupTestWriteOperation } from '@/types/studioV2WriteGroupTestWrite';
import { studioV2WriteGroupLifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';

// Real React Query hooks and the real schema panel; only the service boundary
// is simulated. Server/SQL coverage remains a separate API integration gate.
const backend = vi.hoisted(() => ({ revision: 1, groups: [] as WriteGroup[], tableExists: false, loseConfirmReply: false,
  startOperations: new Map<string, RecordingStartOperation>(), tokens: new Map<string, SchemaPreviewToken>(), operations: new Map<string, SchemaOperation>() }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  list: vi.fn(), readiness: vi.fn(), ensureCircuitProfile: vi.fn(), update: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(), testWriteOperation: vi.fn(), testWritePreview: vi.fn(), testWrite: vi.fn(),
} }));
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: vi.fn() } }));
vi.mock('@/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), operation: vi.fn() } }));
vi.mock('@/services/recordingTestWriteRecovery', () => ({ listRecordingTestWriteRecovery: vi.fn(async (workspaceId: string) => ({ workspace_id: workspaceId, complete: true, operations: [] })) }));
vi.mock('@/services/studioV2WorkspaceDatabase', () => ({ studioV2WorkspaceDatabaseAPI: {
  getMySQLRecordingCapability: vi.fn(async () => ({ connector_id: 'conn-1', connector_revision: 'crev-1', database: 'fixture',
    server_family: 'mysql', server_version: '8.0', capability_revision: 'cap-1', connected: true, metadata_status: 'verified',
    transactional_engine: true, supports_managed_long: true, supports_native_partition: false, supports_group_test_write: true,
    supports_missing_repair: true, reasons: [] })),
} }));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: {
  list: vi.fn(async () => [{ id: 'measurement-0', workspace_id: 'ws-1', device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0',
    name: 'Water', equipment_id: 'equipment-water', definition_revision: 'rev-1', source_binding_revision: 'bind-1', series_epoch: 'epoch-1',
    quantity: 'volume', unit: 'm3', semantic_kind: 'counter', counter_policy: { allow_negative: false } }]),
} }));

const revision = () => `workspace-${backend.revision}`;
const now = new Date().toISOString();
const testWriteResult = (operationId = 'test-write-O1'): WriteGroupTestWriteOperation => ({ operation_id: operationId, action: 'test_write',
  status: 'succeeded', write_outcome: 'written_verified', cleanup_status: 'cleaned', created_at: now, updated_at: now, completed_at: now });
function currentGroup(id: string, request: WriteGroupSchemaPreviewRequest) {
  expect(request.expected_workspace_revision).toBe(revision());
  const group = backend.groups.find((item) => item.id === id);
  expect(group?.revision).toBe(request.expected_group_revision);
  if (!group) throw new Error('missing test group');
  return group;
}

beforeEach(() => {
  vi.clearAllMocks(); window.sessionStorage.clear(); onlineManager.setOnline(true);
  // Explicit fixture reset, not product recovery: the real store intentionally remembers removed entries.
  recovery.read().references.forEach(ref => recovery.resolve(testWriteResult(ref.operation_id), ref.workspace_id));
  localStorage.clear();
  vi.mocked(groupsAPI.testWriteOperation).mockResolvedValue(testWriteResult());
  backend.revision = 1; backend.groups = []; backend.tableExists = false; backend.loseConfirmReply = false;
  backend.tokens.clear(); backend.operations.clear(); backend.startOperations.clear();
  vi.mocked(startAPI.operation).mockImplementation(async id => { const result = backend.startOperations.get(id); if (!result) throw new Error("unknown operation"); return structuredClone(result); });
  vi.mocked(groupsAPI.list).mockImplementation(async () => ({ workspace_id: 'ws-1', workspace_revision: revision(), groups: structuredClone(backend.groups) }));
  vi.mocked(groupsAPI.ensureCircuitProfile).mockImplementation(async (request: WriteGroupCreateRequest) => {
    expect(request.expected_workspace_revision).toBe(revision());
    const group = savedGroup({ ...request.group, id: `group-${request.group.recording_profile?.kind}`, revision: 'group-rev-1', status: 'draft',
      members: request.group.members.map((member) => ({ ...member, source_revision: 'source-1', mapping_revision: 'mapping-1' })),
      destination: { ...request.group.destination, database: '/fixture.db', table_name: 'gw_long_fixture' } });
    backend.groups.push(group); backend.revision += 1;
    return { workspace_revision: revision(), group: structuredClone(group) };
  });
  vi.mocked(groupsAPI.readiness).mockImplementation(async id => {
    const group = backend.groups.find(value => value.id === id);
    if (!group) throw new Error('missing readiness group');
    const ready = group.destination.schema_revision === group.revision;
    return { workspace_id: 'ws-1', workspace_revision: revision(), group_id: id, group_revision: group.revision,
      applied_revision: group.applied_revision, config_ready: true, schema_ready: ready, ready,
      issues: ready ? [] : [{ code: 'schema-confirmation-required', severity: 'blocking', step: 'Step 4', scope: id, message: 'Confirm schema' }] };
  });
  vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockImplementation(async id => {
    const group = backend.groups.find(value => value.id === id);
    if (!group) throw new Error('missing lifecycle group');
    return { ...lifecycle(), workspace_revision: revision(), group_id: id, group_revision: group.revision,
      applied_revision: group.applied_revision, status: group.status, actions: lifecycle().actions.map(action => ({ ...action, allowed: true, reason_code: '' })) };
  });
  vi.mocked(groupsAPI.schemaPreview).mockImplementation(async (id, request) => {
    const group = currentGroup(id, request);
    const operationId = `schema-${id}-${backend.tokens.size}`;
    const preview: SchemaPreviewToken = {
      token: operationId, operation_id: operationId, workspace_id: 'ws-1', workspace_revision: revision(), plan_id: id, plan_revision: group.revision,
      connector_id: 'conn-1', table_prefix: 'gw_long_fixture', statements: backend.tableExists ? [] : ['CREATE TABLE fixture (...)'], tables: [],
      digest: 'a'.repeat(64), created_at: now, expires_at: new Date(Date.now() + 600000).toISOString(), recreate_missing: true,
      ...(backend.tableExists ? { no_change_reason: 'schema_already_compatible' } : {}),
    };
    backend.tokens.set(operationId, preview);
    return preview;
  });
  vi.mocked(groupsAPI.schemaApply).mockImplementation(async (id: string, request: WriteGroupSchemaApplyRequest) => {
    const group = currentGroup(id, request);
    const token = backend.tokens.get(request.token);
    expect(token?.workspace_revision).toBe(revision());
    backend.tableExists = true;
    group.destination.schema_revision = group.revision;
    group.destination.schema_digest = 'a'.repeat(64);
    backend.revision += 1;
    const operation: SchemaOperation = { operation_id: request.operation_id, status: 'succeeded', executed_statements: token?.statements.length ?? 0,
      created_at: now, updated_at: now };
    backend.operations.set(request.operation_id, operation);
    if (backend.loseConfirmReply) { backend.loseConfirmReply = false; throw new TypeError('Failed to fetch'); }
    return operation;
  });
  vi.mocked(groupsAPI.schemaOperation).mockImplementation(async (id) => {
    const operation = backend.operations.get(id);
    if (!operation) throw new Error('missing fixture operation');
    return operation;
  });
  vi.mocked(startAPI.start).mockImplementation(async (request: RecordingStartRequest) => {
    expect(request.expected_workspace_revision).toBe(revision());
    for (const item of request.groups) {
      const group = backend.groups.find((saved) => saved.id === item.group_id);
      expect(group?.destination.schema_revision).toBe(item.expected_group_revision);
    }
    const operation: RecordingStartOperation = { operation_id: 'start-1', action: 'recording_start', status: 'succeeded', stage: 'complete',
      intent_digest: 'b'.repeat(64), workspace_id: 'ws-1', setup_revision: revision(), device_ids: ['dev-1'],
      groups: request.groups.map((group) => ({ group_id: group.group_id, group_revision: group.expected_group_revision, saved: true, applied: true, ready: true })),
      devices: [{ device_id: 'dev-1', activated: true }], created_at: now, updated_at: now };
    backend.startOperations.set(operation.operation_id, structuredClone(operation));
    return operation;
  });
});

async function prepare(both = false, connectorKind: 'sqlite' | 'mysql' = 'sqlite') {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const tree = (state: WorkbenchV2State) => <QueryClientProvider client={client}><CircuitProfileBasicPanel state={state} workspaceId="ws-1" readonly={false} /></QueryClientProvider>;
  const initial = savedGroupState(); initial.db.connector.kind = connectorKind;
  const view = render(tree(initial));
  await waitFor(() => expect(groupsAPI.list).toHaveBeenCalled());
  fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: '水量迴路' } });
  expandSources();
  fireEvent.click(screen.getByLabelText('line.pressure')); fireEvent.click(screen.getByLabelText('line.batch'));
  if (both) {
    fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
    fireEvent.click(screen.getByTestId('circuit-profile-member-cum-tag-0'));
  }
  await waitFor(() => expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled());
  fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
  await screen.findByTestId('circuit-prepared-group-instantaneous');
  expect(groupsAPI.schemaPreview).not.toHaveBeenCalled();
  const preview = await screen.findByTestId('group-schema-preview');
  await waitFor(() => expect(preview).toBeEnabled());
  fireEvent.click(preview);
  await waitFor(() => expect(groupsAPI.schemaPreview).toHaveBeenCalledOnce());
  return { client, rerenderState: (state: WorkbenchV2State) => view.rerender(tree(state)) };
}


afterEach(() => { cleanup(); vi.restoreAllMocks(); onlineManager.setOnline(true); });

it('requires a separate exact-scope confirmation before the existing start owner sends any start request', async () => {
  await prepare();
  fireEvent.click(await screen.findByTestId('group-schema-apply'));
  await waitFor(() => expect(backend.revision).toBe(3));
  const reviewStart = await screen.findByTestId('circuit-profile-start');
  await waitFor(() => expect(reviewStart).toBeEnabled());
  fireEvent.click(reviewStart);
  expect(screen.getByRole('dialog')).toBeVisible();
  expect(startAPI.start).not.toHaveBeenCalled();
  fireEvent.click(screen.getByTestId('recording-start-cancel'));
  expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
  expect(startAPI.start).not.toHaveBeenCalled();
});

it('mounts exactly one guided primary while keeping source controls editable', async () => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  render(<QueryClientProvider client={client}><CircuitProfileBasicPanel state={savedGroupState()} workspaceId="ws-1" readonly={false} /></QueryClientProvider>);
  await waitFor(() => expect(groupsAPI.list).toHaveBeenCalled());
  expect(document.querySelectorAll('button[data-primary="true"]')).toHaveLength(1);
  expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toBeEnabled();
  expect(groupsAPI.schemaPreview).not.toHaveBeenCalled();
  expect(groupsAPI.schemaApply).not.toHaveBeenCalled();
  expect(startAPI.start).not.toHaveBeenCalled();
});

async function readyForReview() {
  const view = await prepare();
  expect(groupsAPI.schemaApply).not.toHaveBeenCalled();
  fireEvent.click(await screen.findByTestId('group-schema-apply'));
  await waitFor(() => expect(backend.revision).toBe(3));
  await screen.findByTestId('circuit-profile-start');
  return view;
}

it('starts once only after explicit confirmation, and never calls a workspace-wide start', async () => {
  await readyForReview();
  fireEvent.click(screen.getByTestId('circuit-profile-start'));
  const confirm = screen.getByTestId('recording-start-confirm');
  expect(document.querySelectorAll('button[data-primary="true"]')).toHaveLength(1);
  fireEvent.click(confirm); fireEvent.click(confirm);
  await waitFor(() => expect(startAPI.start).toHaveBeenCalledOnce());
  expect(vi.mocked(startAPI.start).mock.calls[0][0].device_ids).toEqual(['dev-1']);
  expect(vi.mocked(startAPI.start).mock.calls[0][0].groups.map(group => group.group_id)).toEqual(['group-instantaneous']);
  expect(groupsAPI.schemaApply).toHaveBeenCalledOnce();
});

it('invalidates consent on native interval change and does not revive it after A to B to A controls', async () => {
  await readyForReview();
  fireEvent.click(screen.getByTestId('circuit-profile-start'));
  const interval = screen.getByTestId('circuit-profile-interval-instantaneous');
  const original = (interval as HTMLInputElement).value;
  fireEvent.change(interval, { target: { value: String(Number(original) + 1) } });
  expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
  fireEvent.change(interval, { target: { value: original } });
  expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
  expect(startAPI.start).not.toHaveBeenCalled();
  fireEvent.click(await screen.findByTestId('circuit-profile-start'));
  expect(screen.getByTestId('recording-start-confirm')).toBeEnabled();
  expect(startAPI.start).not.toHaveBeenCalled();
});

it('keeps offline scoped review and cancel local while disabling final network confirmation', async () => {
  await readyForReview();
  vi.spyOn(navigator, 'onLine', 'get').mockReturnValue(false);
  fireEvent(window, new Event('offline'));
  fireEvent.click(screen.getByTestId('circuit-profile-start'));
  expect(screen.getByRole('dialog')).toBeVisible();
  expect(screen.getByTestId('recording-start-confirm')).toBeDisabled();
  fireEvent.click(screen.getByTestId('recording-start-cancel'));
  expect(startAPI.start).not.toHaveBeenCalled();
  expect(groupsAPI.schemaPreview).toHaveBeenCalledOnce();
});

it('focuses the actual interval repair control without saving when an interval is invalid', async () => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  render(<QueryClientProvider client={client}><CircuitProfileBasicPanel state={savedGroupState()} workspaceId="ws-1" readonly={false} /></QueryClientProvider>);
  await waitFor(() => expect(groupsAPI.list).toHaveBeenCalled());
  const interval = screen.getByTestId('circuit-profile-interval-instantaneous');
  fireEvent.change(interval, { target: { value: '' } });
  fireEvent.click(document.querySelector('button[data-primary="true"]')!);
  expect(interval).toHaveFocus();
  expect(groupsAPI.ensureCircuitProfile).not.toHaveBeenCalled();
  expect(groupsAPI.schemaPreview).not.toHaveBeenCalled();
  expect(startAPI.start).not.toHaveBeenCalled();
});

it('retains a lost start acknowledgement without inventing an ID or resending a request', async () => {
  await readyForReview();
  vi.mocked(startAPI.start).mockRejectedValueOnce(new TypeError('Failed to fetch'));
  fireEvent.click(screen.getByTestId('circuit-profile-start'));
  fireEvent.click(screen.getByTestId('recording-start-confirm'));
  await waitFor(() => expect(startAPI.start).toHaveBeenCalledOnce());
  const original = JSON.stringify(vi.mocked(startAPI.start).mock.calls[0][0]);
  await waitFor(() => expect(screen.getByTestId('recording-next-action')).toHaveTextContent('next.actions.recover_operation'));
  fireEvent.click(screen.getByTestId('recording-next-action'));
  expect(startAPI.operation).not.toHaveBeenCalled();
  expect(startAPI.start).toHaveBeenCalledOnce();
  expect(JSON.stringify(vi.mocked(startAPI.start).mock.calls[0][0])).toBe(original);
  expect(screen.queryByTestId('circuit-profile-start')).not.toBeInTheDocument();
});

it('queries the original schema operation after a lost apply reply, without previewing or applying again', async () => {
  await prepare(); backend.loseConfirmReply = true;
  fireEvent.click(await screen.findByTestId('group-schema-apply'));
  const check = await screen.findByTestId('basic-recording-check');
  await waitFor(() => expect(check).toBeEnabled());
  const original = vi.mocked(groupsAPI.schemaApply).mock.calls[0][1].operation_id;
  fireEvent.click(check);
  await waitFor(() => expect(groupsAPI.schemaOperation).toHaveBeenCalledWith(original));
  expect(groupsAPI.schemaPreview).toHaveBeenCalledOnce();
  expect(groupsAPI.schemaApply).toHaveBeenCalledOnce();
  expect(startAPI.start).not.toHaveBeenCalled();
});


it.each(['sources', 'profile', 'schema policy'] as const)('retires start consent after unsaved native %s changes without a revision change', async field => {
  await readyForReview();
  fireEvent.click(screen.getByTestId('circuit-profile-start'));
  const before = revision();
  const target = field === 'sources' ? 'circuit-profile-member-inst-tag-0'
    : field === 'profile' ? 'circuit-profile-toggle-cumulative' : 'circuit-profile-recreate-missing';
  fireEvent.click(screen.getByTestId(target));
  expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
  expect(revision()).toBe(before);
  expect(startAPI.start).not.toHaveBeenCalled();
  expect(groupsAPI.schemaApply).toHaveBeenCalledOnce();
});

it('retires consent when destination identity changes before any saved revision changes', async () => {
  const view = await readyForReview();
  fireEvent.click(screen.getByTestId('circuit-profile-start'));
  const state = savedGroupState();
  view.rerenderState({ ...state, db: { ...state.db, connector: { ...state.db.connector, connector_id: 'conn-B', identity_revision: 'crev-B' } } });
  expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
  expect(revision()).toBe('workspace-3');
  expect(startAPI.start).not.toHaveBeenCalled();
});

it('rechecks a freshly retained schema O1 before the captured final start command', async () => {
  await readyForReview();
  fireEvent.click(screen.getByTestId('circuit-profile-start'));
  const confirm = screen.getByTestId('recording-start-confirm');
  retainSchemaOperation('ws-1', 'group-instantaneous', 'late-schema-O1');
  fireEvent.click(confirm);
  expect(startAPI.start).not.toHaveBeenCalled();
  expect(groupsAPI.ensureCircuitProfile).toHaveBeenCalledOnce();
  expect(groupsAPI.schemaPreview).toHaveBeenCalledOnce();
  expect(groupsAPI.schemaApply).toHaveBeenCalledOnce();
  retainSchemaOperation('ws-1', 'group-instantaneous', undefined, 'late-schema-O1');
});

it('checks exact retained test-write O1 once and keeps its verified outcome visible after resolving the fence', async () => {
  await readyForReview();
  act(() => { recovery.retain({ workspace_id: 'ws-1', group_id: 'group-instantaneous', group_revision: 'older-original-revision',
    operation_id: 'test-write-O1', operation_kind: 'test_write', created_at: now }); });
  const check = await screen.findByTestId('basic-recording-check');
  await waitFor(() => expect(check).toBeEnabled());
  expect(groupsAPI.testWriteOperation).not.toHaveBeenCalled();
  fireEvent.click(check); fireEvent.click(check);
  await waitFor(() => expect(groupsAPI.testWriteOperation).toHaveBeenCalledExactlyOnceWith('test-write-O1'));
  expect(await screen.findByTestId('preparation-test-write-result')).toHaveTextContent('test-write-O1');
  expect(screen.getByTestId('preparation-test-write-write-phase')).toHaveTextContent('write_committed');
  expect(screen.getByTestId('preparation-test-write-readback')).toHaveTextContent('readback_verified');
  expect(screen.getByTestId('preparation-test-write-cleanup')).toHaveTextContent('cleanup_cleaned');
  expect(recovery.read().references).toEqual([]);
  await screen.findByTestId('circuit-profile-start');
  expect(screen.getByTestId('preparation-test-write-result')).toHaveTextContent('test-write-O1');
  expect(groupsAPI.testWritePreview).not.toHaveBeenCalled(); expect(groupsAPI.testWrite).not.toHaveBeenCalled();
  expect(startAPI.start).not.toHaveBeenCalled();
});


it('fences secondary new test writes across both selected profiles when a fresh sibling schema O1 is retained', async () => {
  await prepare(true, 'mysql');
  fireEvent.click(await screen.findByTestId('group-schema-apply'));
  const second = await screen.findByTestId('group-schema-preview');
  await waitFor(() => expect(second).toBeEnabled()); fireEvent.click(second);
  await waitFor(() => expect(groupsAPI.schemaPreview).toHaveBeenCalledTimes(2));
  fireEvent.click(await screen.findByTestId('group-schema-apply'));
  await screen.findByTestId('circuit-profile-start');
  const cumulativePreview = within(screen.getByTestId('circuit-prepared-group-cumulative')).getByTestId('group-test-write-preview');
  await waitFor(() => expect(cumulativePreview).toBeEnabled());
  retainSchemaOperation('ws-1', 'group-instantaneous', 'sibling-schema-O1');
  fireEvent.click(cumulativePreview); // Stored reference changed before a presentation rerender.
  expect(groupsAPI.testWritePreview).not.toHaveBeenCalled(); expect(groupsAPI.testWrite).not.toHaveBeenCalled();
  fireEvent(window, new Event('storage'));
  await waitFor(() => screen.getAllByTestId('group-test-write-preview').forEach(button => expect(button).toBeDisabled()));
  expect(startAPI.start).not.toHaveBeenCalled();
  expect(groupsAPI.schemaApply).toHaveBeenCalledTimes(2);
  retainSchemaOperation('ws-1', 'group-instantaneous', undefined, 'sibling-schema-O1');
});


it.each(['foreign', 'error', 'pending'] as const)('does not authorize a new draft save from %s canonical list evidence', async evidence => {
  if (evidence === 'foreign') vi.mocked(groupsAPI.list).mockResolvedValue({ workspace_id: 'other-workspace', workspace_revision: 'workspace-1', groups: [] });
  if (evidence === 'error') vi.mocked(groupsAPI.list).mockRejectedValue(new Error('list unavailable'));
  if (evidence === 'pending') vi.mocked(groupsAPI.list).mockReturnValue(new Promise(() => undefined));
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  render(<QueryClientProvider client={client}><CircuitProfileBasicPanel state={savedGroupState()} workspaceId="ws-1" workspaceRevision="workspace-1" readonly={false} /></QueryClientProvider>);
  await waitFor(() => expect(groupsAPI.list).toHaveBeenCalledOnce());
  fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: 'Water' } });
  expandSources(); fireEvent.click(screen.getByLabelText('line.pressure')); fireEvent.click(screen.getByLabelText('line.batch'));
  const prepare = await screen.findByTestId('circuit-profile-prepare');
  await waitFor(() => expect(prepare).toBeDisabled()); fireEvent.click(prepare);
  expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toBeEnabled();
  expect(groupsAPI.ensureCircuitProfile).not.toHaveBeenCalled(); expect(groupsAPI.schemaPreview).not.toHaveBeenCalled();
  expect(startAPI.start).not.toHaveBeenCalled();
});

it('rechecks the current canonical workspace head before a previously reviewed start confirmation', async () => {
  const view = await readyForReview();
  fireEvent.click(screen.getByTestId('circuit-profile-start'));
  act(() => { view.client.setQueryData(studioV2WorkspaceKeys.writeGroups(), {
    workspace_id: 'other-workspace', workspace_revision: revision(), groups: structuredClone(backend.groups),
  }); });
  const confirm = screen.queryByTestId('recording-start-confirm');
  // Deliberately invoke before React Query's observer notification renders new props.
  if (confirm) fireEvent.click(confirm);
  await act(async () => { await Promise.resolve(); });
  expect(startAPI.start).not.toHaveBeenCalled();
  expect(groupsAPI.ensureCircuitProfile).toHaveBeenCalledOnce(); expect(groupsAPI.schemaApply).toHaveBeenCalledOnce();
});


async function startWithRetryablePartial() {
  const view = await readyForReview();
  const native = vi.mocked(startAPI.start).getMockImplementation()!;
  vi.mocked(startAPI.start).mockImplementationOnce(async request => {
    const operation = await native(request);
    const partial = { ...operation, status: 'partial' as const, stage: 'activation' as const, next_action: 'retry_same_request',
      devices: operation.devices.map(device => ({ ...device, activated: false })) };
    backend.startOperations.set(partial.operation_id, partial); return partial;
  });
  fireEvent.click(screen.getByTestId('circuit-profile-start'));
  fireEvent.click(screen.getByTestId('recording-start-confirm'));
  return { view, retry: await screen.findByTestId('basic-recording-retry') };
}

it('keeps exact server-authorized same-request recovery secondary to querying the original operation', async () => {
  const { retry } = await startWithRetryablePartial();
  const original = structuredClone(vi.mocked(startAPI.start).mock.calls[0][0]);
  expect(screen.getByTestId('basic-recording-check')).toHaveAttribute('data-primary', 'true');
  expect(retry).not.toHaveAttribute('data-primary', 'true');
  const checks = vi.mocked(startAPI.operation).mock.calls.length;
  fireEvent.click(retry);
  await waitFor(() => expect(startAPI.start).toHaveBeenCalledTimes(2));
  expect(vi.mocked(startAPI.start).mock.calls[1][0]).toEqual(original);
  expect(vi.mocked(startAPI.operation).mock.calls.length).toBeGreaterThan(checks);
  expect(groupsAPI.ensureCircuitProfile).toHaveBeenCalledOnce(); expect(groupsAPI.schemaApply).toHaveBeenCalledOnce();
});

it('does not retry after the fresh lookup withdraws same-request authorization', async () => {
  const { retry } = await startWithRetryablePartial();
  const original = backend.startOperations.get('start-1')!;
  backend.startOperations.set('start-1', { ...original, status: 'unknown', next_action: undefined });
  const checks = vi.mocked(startAPI.operation).mock.calls.length;
  fireEvent.click(retry);
  await waitFor(() => expect(vi.mocked(startAPI.operation).mock.calls.length).toBeGreaterThan(checks));
  await waitFor(() => expect(screen.queryByTestId('basic-recording-retry')).not.toBeInTheDocument());
  expect(startAPI.start).toHaveBeenCalledOnce();
  expect(screen.getByTestId('basic-recording-check')).toHaveAttribute('data-primary', 'true');
});

it('rechecks the same canonical cache before a secondary new test-write preview dispatch', async () => {
  const view = await prepare(false, 'mysql');
  fireEvent.click(await screen.findByTestId('group-schema-apply'));
  await screen.findByTestId('circuit-profile-start');
  const preview = screen.getByTestId('group-test-write-preview');
  await waitFor(() => expect(preview).toBeEnabled());
  act(() => { view.client.setQueryData(studioV2WorkspaceKeys.writeGroups(), {
    workspace_id: 'other-workspace', workspace_revision: revision(), groups: structuredClone(backend.groups),
  }); });
  fireEvent.click(preview);
  await act(async () => { await Promise.resolve(); });
  expect(groupsAPI.testWritePreview).not.toHaveBeenCalled(); expect(groupsAPI.testWrite).not.toHaveBeenCalled();
  expect(startAPI.start).not.toHaveBeenCalled();
});

it('keeps an accepted two-target batch moving when its own preceding target advances the canonical head', async () => {
  backend.groups = ['A', 'B'].map((suffix, index) => savedGroup({ id: `batch-${suffix}`, name: suffix,
    layout_version: 'managed_long_v1', storage_mode: 'plain', status: 'draft',
    members: [{ ...savedGroup().members[index], circuit_id: `circuit-${suffix}`, circuit_name: suffix,
      required: false, target_column: '', value_role: 'gauge' }],
    recording_profile: { kind: 'instantaneous', interval: 15, unit: 'second', sample_reduction: 'mean' },
    destination: { ...savedGroup().destination, storage_strategy: 'managed', table_name: `table_${suffix}`, schema_revision: 'rev-1' },
    row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
    write_policy: { mode: 'append', recreate_missing: true },
  }));
  vi.mocked(groupsAPI.update).mockImplementation(async (id, request) => {
    const group = backend.groups.find(value => value.id === id);
    expect(request.expected_workspace_revision).toBe(revision()); expect(request.expected_group_revision).toBe(group?.revision);
    if (!group) throw new Error('missing batch fixture target');
    Object.assign(group, request.group, { revision: `${group.revision}-next` }); backend.revision += 1;
    return { workspace_revision: revision(), group: structuredClone(group) };
  });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><CircuitProfileBasicPanel state={savedGroupState()} workspaceId="ws-1" readonly={false} /></QueryClientProvider>);
  const chooser = await screen.findByTestId('circuit-profile-circuit-chooser');
  await waitFor(() => expect(chooser.querySelector('option[value="circuit-B"]')).not.toBeNull());
  fireEvent.change(chooser, { target: { value: 'circuit-B' } });
  await waitFor(() => expect(screen.getByTestId('circuit-profile-start')).toBeEnabled());
  const listReadsBeforeBatch = vi.mocked(groupsAPI.list).mock.calls.length;
  fireEvent.click(screen.getByTestId('circuit-profile-open-batch'));
  fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
  expect(groupsAPI.update).not.toHaveBeenCalled();
  fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm'));
  await waitFor(() => expect(groupsAPI.update).toHaveBeenCalledTimes(2));
  await waitFor(() => expect(screen.queryByTestId('circuit-profile-batch-editor')).not.toBeInTheDocument());
  expect(groupsAPI.list).toHaveBeenCalledTimes(listReadsBeforeBatch + 3); // Two native invalidations and the existing batch host reload.
  expect(vi.mocked(groupsAPI.update).mock.calls.map(([id]) => id)).toEqual(['batch-A', 'batch-B']);
  expect(backend.groups.map(group => group.row_policy.interval_seconds)).toEqual([30, 30]);
  expect(vi.mocked(groupsAPI.update).mock.calls[1][1].expected_workspace_revision).toBe('workspace-2');
  expect(groupsAPI.ensureCircuitProfile).not.toHaveBeenCalled(); expect(startAPI.start).not.toHaveBeenCalled();
});


it.each(['foreign', 'A to B to A'] as const)('denies original retry after %s canonical cache changes during its fresh GET', async transition => {
  const { retry, view } = await startWithRetryablePartial();
  let resolve!: (value: RecordingStartOperation) => void;
  const pending = new Promise<RecordingStartOperation>(done => { resolve = done; });
  vi.mocked(startAPI.operation).mockReturnValueOnce(pending);
  const checks = vi.mocked(startAPI.operation).mock.calls.length;
  fireEvent.click(retry);
  await waitFor(() => expect(vi.mocked(startAPI.operation).mock.calls.length).toBeGreaterThan(checks));
  const saved = view.client.getQueryData(studioV2WorkspaceKeys.writeGroups());
  act(() => { view.client.setQueryData(studioV2WorkspaceKeys.writeGroups(), { workspace_id: 'other-workspace', workspace_revision: revision(), groups: structuredClone(backend.groups) }); });
  view.rerenderState(savedGroupState());
  if (transition === 'A to B to A') {
    act(() => { view.client.setQueryData(studioV2WorkspaceKeys.writeGroups(), saved); });
    view.rerenderState(savedGroupState());
  }
  await act(async () => { resolve(backend.startOperations.get('start-1')!); await pending; });
  expect(startAPI.start).toHaveBeenCalledOnce();
  expect(groupsAPI.ensureCircuitProfile).toHaveBeenCalledOnce();
});
