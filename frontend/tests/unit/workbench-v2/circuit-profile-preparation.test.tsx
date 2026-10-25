import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { CircuitProfileBasicPanel } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel';
import { studioV2WorkspaceWriteGroupsAPI as groupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2RecordingStartAPI as startAPI } from '@/services/studioV2RecordingStart';
import type { SchemaOperation, SchemaPreviewToken } from '@/types/recordingPlan';
import type { WriteGroup, WriteGroupCreateRequest, WriteGroupSchemaApplyRequest, WriteGroupSchemaPreviewRequest } from '@/types/studioV2WriteGroup';
import type { RecordingStartOperation, RecordingStartRequest } from '@/types/studioV2RecordingStart';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import { expandSources } from '../../fixtures/sourceReview';

// Real React Query hooks and the real schema panel; only the service boundary
// is simulated. Server/SQL coverage remains a separate API integration gate.
const backend = vi.hoisted(() => ({ revision: 1, groups: [] as WriteGroup[], tableExists: false, loseConfirmReply: false,
  tokens: new Map<string, SchemaPreviewToken>(), operations: new Map<string, SchemaOperation>() }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  list: vi.fn(), ensureCircuitProfile: vi.fn(), update: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(),
} }));
vi.mock('@/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), operation: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: {
  list: vi.fn(async () => [{ id: 'measurement-0', workspace_id: 'ws-1', device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0',
    name: 'Water', equipment_id: 'equipment-water', definition_revision: 'rev-1', source_binding_revision: 'bind-1', series_epoch: 'epoch-1',
    quantity: 'volume', unit: 'm3', semantic_kind: 'counter', counter_policy: { allow_negative: false } }]),
} }));

const revision = () => `workspace-${backend.revision}`;
const now = '2026-10-24T00:00:00Z';
function currentGroup(id: string, request: WriteGroupSchemaPreviewRequest) {
  expect(request.expected_workspace_revision).toBe(revision());
  const group = backend.groups.find((item) => item.id === id);
  expect(group?.revision).toBe(request.expected_group_revision);
  if (!group) throw new Error('missing test group');
  return group;
}

beforeEach(() => {
  vi.clearAllMocks(); window.sessionStorage.clear();
  backend.revision = 1; backend.groups = []; backend.tableExists = false; backend.loseConfirmReply = false;
  backend.tokens.clear(); backend.operations.clear();
  vi.mocked(groupsAPI.list).mockImplementation(async () => ({ workspace_id: 'ws-1', workspace_revision: revision(), groups: structuredClone(backend.groups) }));
  vi.mocked(groupsAPI.ensureCircuitProfile).mockImplementation(async (request: WriteGroupCreateRequest) => {
    expect(request.expected_workspace_revision).toBe(revision());
    const group = savedGroup({ ...request.group, id: `group-${request.group.recording_profile?.kind}`, revision: 'group-rev-1', status: 'draft',
      members: request.group.members.map((member) => ({ ...member, source_revision: 'source-1', mapping_revision: 'mapping-1' })),
      destination: { ...request.group.destination, database: '/fixture.db', table_name: 'gw_long_fixture' } });
    backend.groups.push(group); backend.revision += 1;
    return { workspace_revision: revision(), group: structuredClone(group) };
  });
  vi.mocked(groupsAPI.schemaPreview).mockImplementation(async (id, request) => {
    const group = currentGroup(id, request);
    const operationId = `schema-${id}-${backend.tokens.size}`;
    const preview: SchemaPreviewToken = {
      token: operationId, operation_id: operationId, workspace_id: 'ws-1', workspace_revision: revision(), plan_id: id, plan_revision: group.revision,
      connector_id: 'conn-1', table_prefix: 'gw_long_fixture', statements: backend.tableExists ? [] : ['CREATE TABLE fixture (...)'], tables: [],
      digest: 'a'.repeat(64), created_at: now, expires_at: '2026-10-25T00:00:00Z', recreate_missing: true,
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
    return operation;
  });
});

async function prepare(both = false) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><CircuitProfileBasicPanel state={savedGroupState()} workspaceId="ws-1" readonly={false} /></QueryClientProvider>);
  await waitFor(() => expect(groupsAPI.list).toHaveBeenCalled());
  fireEvent.change(screen.getByTestId('circuit-profile-circuit-name'), { target: { value: '水量迴路' } });
  expandSources();
  fireEvent.click(screen.getByLabelText('line.pressure')); fireEvent.click(screen.getByLabelText('line.batch'));
  if (both) fireEvent.click(screen.getByTestId('circuit-profile-toggle-cumulative'));
  await waitFor(() => expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled());
  fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
  await screen.findByTestId('circuit-prepared-group-instantaneous');
  await waitFor(() => expect(groupsAPI.schemaPreview).toHaveBeenCalled());
}

it('confirms both profiles explicitly, refreshes the empty-DDL preview, then starts with the current revisions', async () => {
  await prepare(true);
  expect(groupsAPI.schemaApply).not.toHaveBeenCalled();
  expect(backend.tableExists).toBe(false);
  const inst = within(screen.getByTestId('circuit-prepared-group-instantaneous'));
  fireEvent.click(await inst.findByTestId('group-schema-apply'));
  await waitFor(() => expect(backend.revision).toBe(4));
  await waitFor(() => expect(vi.mocked(groupsAPI.schemaPreview).mock.calls.some(([id, req]) => id === 'group-cumulative' && req.expected_workspace_revision === 'workspace-4')).toBe(true));
  const cum = within(screen.getByTestId('circuit-prepared-group-cumulative'));
  await waitFor(() => expect(cum.getByTestId('group-schema-apply')).toBeEnabled());
  expect(cum.getByTestId('group-schema-preview-table')).toHaveTextContent('step4.group.schema.no_change');
  fireEvent.click(cum.getByTestId('group-schema-apply'));
  await waitFor(() => expect(backend.revision).toBe(5));
  await waitFor(() => expect(screen.getByTestId('circuit-profile-start')).toBeEnabled());
  fireEvent.click(screen.getByTestId('circuit-profile-start'));
  await waitFor(() => expect(startAPI.start).toHaveBeenCalledOnce());
  expect(vi.mocked(startAPI.start).mock.calls[0][0].groups).toHaveLength(2);
  expect(groupsAPI.schemaApply).toHaveBeenCalledTimes(2);
});

it('keeps a confirmation reply loss fenced to the same operation until readback succeeds', async () => {
  await prepare(); backend.loseConfirmReply = true;
  fireEvent.click(await screen.findByTestId('group-schema-apply'));
  const check = await screen.findByTestId('group-schema-check');
  expect(groupsAPI.schemaApply).toHaveBeenCalledOnce();
  fireEvent.click(check);
  await waitFor(() => expect(groupsAPI.schemaOperation).toHaveBeenCalledOnce());
  await waitFor(() => expect(queryRevisionShown()).toBe(revision()));
  expect(groupsAPI.schemaApply).toHaveBeenCalledOnce();
});

it('gives the two schema regions distinct accessible heading identities', async () => {
  await prepare(true);
  const headings = Array.from(document.querySelectorAll('[data-testid="group-schema-panel"] h4[id]'));
  expect(headings).toHaveLength(2);
  expect(new Set(headings.map((heading) => heading.id)).size).toBe(2);
});

function queryRevisionShown() {
  return vi.mocked(groupsAPI.schemaPreview).mock.lastCall?.[1].expected_workspace_revision;
}
