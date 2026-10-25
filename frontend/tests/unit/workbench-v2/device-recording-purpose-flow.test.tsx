import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { CircuitProfileBasicPanel } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileBasicPanel';
import { studioV2WorkspaceWriteGroupsAPI as groupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2RecordingStartAPI as startAPI } from '@/services/studioV2RecordingStart';
import type { SchemaOperation, SchemaPreviewToken } from '@/types/recordingPlan';
import type { WriteGroup, WriteGroupCreateRequest, WriteGroupSchemaApplyRequest, WriteGroupSchemaPreviewRequest } from '@/types/studioV2WriteGroup';
import type { RecordingStartOperation, RecordingStartRequest } from '@/types/studioV2RecordingStart';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';

// Real React Query hooks and the real schema panel; only the service boundary is simulated.
// The purpose decides the layout the group is saved with; the plain preview, the explicit
// confirmation and the start are the same for both.
const backend = vi.hoisted(() => ({ revision: 1, groups: [] as WriteGroup[], rejectPreview: undefined as unknown,
  tokens: new Map<string, SchemaPreviewToken>(), operations: new Map<string, SchemaOperation>() }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  list: vi.fn(), ensureCircuitProfile: vi.fn(), update: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(),
} }));
vi.mock('@/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), operation: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { list: vi.fn(async () => []) } }));

const revision = () => `workspace-${backend.revision}`;
const now = '2026-10-24T00:00:00Z';

beforeEach(() => {
  vi.clearAllMocks(); window.sessionStorage.clear();
  backend.revision = 1; backend.groups = []; backend.rejectPreview = undefined; backend.tokens.clear(); backend.operations.clear();
  vi.mocked(groupsAPI.list).mockImplementation(async () => ({ workspace_id: 'ws-1', workspace_revision: revision(), groups: structuredClone(backend.groups) }));
  vi.mocked(groupsAPI.ensureCircuitProfile).mockImplementation(async (request: WriteGroupCreateRequest) => {
    const group = savedGroup({ ...request.group, id: 'group-new', revision: 'group-rev-1', status: 'draft',
      members: request.group.members.map((member) => ({ ...member, source_revision: 'source-1', mapping_revision: 'mapping-1' })),
      destination: { ...request.group.destination, database: '/fixture.db', table_name: 'gw_fixture' } });
    backend.groups.push(group); backend.revision += 1;
    return { workspace_revision: revision(), group: structuredClone(group) };
  });
  vi.mocked(groupsAPI.schemaPreview).mockImplementation(async (id: string, request: WriteGroupSchemaPreviewRequest) => {
    if (backend.rejectPreview) throw backend.rejectPreview;
    const group = backend.groups.find((item) => item.id === id);
    if (!group) throw new Error('missing test group');
    const operationId = `schema-${id}-${backend.tokens.size}`;
    const preview: SchemaPreviewToken = {
      token: operationId, operation_id: operationId, workspace_id: 'ws-1', workspace_revision: request.expected_workspace_revision, plan_id: id,
      plan_revision: group.revision, connector_id: 'conn-1', table_prefix: 'gw_fixture', statements: ['CREATE TABLE fixture (...)'], tables: [],
      digest: 'a'.repeat(64), created_at: now, expires_at: '2026-10-25T00:00:00Z', recreate_missing: true,
    };
    backend.tokens.set(operationId, preview);
    return preview;
  });
  vi.mocked(groupsAPI.schemaApply).mockImplementation(async (id: string, request: WriteGroupSchemaApplyRequest) => {
    const group = backend.groups.find((item) => item.id === id);
    if (!group) throw new Error('missing test group');
    group.destination.schema_revision = group.revision; group.destination.schema_digest = 'a'.repeat(64); backend.revision += 1;
    const operation: SchemaOperation = { operation_id: request.operation_id, status: 'succeeded', executed_statements: 1, created_at: now, updated_at: now };
    backend.operations.set(request.operation_id, operation);
    return operation;
  });
  vi.mocked(startAPI.start).mockImplementation(async (request: RecordingStartRequest) => ({
    operation_id: 'start-1', action: 'recording_start', status: 'succeeded', stage: 'complete', intent_digest: 'b'.repeat(64), workspace_id: 'ws-1',
    setup_revision: revision(), device_ids: ['dev-1'], devices: [{ device_id: 'dev-1', activated: true }], created_at: now, updated_at: now,
    groups: request.groups.map((group) => ({ group_id: group.group_id, group_revision: group.expected_group_revision, saved: true, applied: true, ready: true })),
  } satisfies RecordingStartOperation));
});

async function prepare(purpose: 'wide' | 'long') {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><CircuitProfileBasicPanel state={savedGroupState()} workspaceId="ws-1" readonly={false} /></QueryClientProvider>);
  await waitFor(() => expect(groupsAPI.list).toHaveBeenCalled());
  if (purpose === 'wide') fireEvent.click(screen.getByTestId('circuit-profile-purpose-wide'));
  await waitFor(() => expect(screen.getByTestId('circuit-profile-prepare')).toBeEnabled());
  fireEvent.click(screen.getByTestId('circuit-profile-prepare'));
  await screen.findByTestId('circuit-prepared-group-instantaneous');
}

describe.each(['long', 'wide'] as const)('the %s purpose', (purpose) => {
  const layout = purpose === 'wide' ? 'managed_wide_profile_v1' : 'managed_long_v1';

  it('previews the plain schema, applies only after an explicit confirmation, then starts with the saved revision', async () => {
    await prepare(purpose);
    expect(vi.mocked(groupsAPI.ensureCircuitProfile).mock.calls[0][0].group.layout_version).toBe(layout);
    await waitFor(() => expect(groupsAPI.schemaPreview).toHaveBeenCalled());
    expect(groupsAPI.schemaApply).not.toHaveBeenCalled();
    expect(startAPI.start).not.toHaveBeenCalled();

    const prepared = within(screen.getByTestId('circuit-prepared-group-instantaneous'));
    fireEvent.click(await prepared.findByTestId('group-schema-apply'));
    await waitFor(() => expect(groupsAPI.schemaApply).toHaveBeenCalledTimes(1));
    await waitFor(() => expect(screen.getByTestId('circuit-profile-start')).toBeEnabled());
    fireEvent.click(screen.getByTestId('circuit-profile-start'));
    await waitFor(() => expect(startAPI.start).toHaveBeenCalledOnce());
    expect(vi.mocked(startAPI.start).mock.calls[0][0].groups.map((group) => group.group_id)).toEqual(['group-new']);
    expect(groupsAPI.update).not.toHaveBeenCalled();
  });

  it('refuses to confirm a schema preview once the settings changed after it', async () => {
    await prepare(purpose);
    await waitFor(() => expect(groupsAPI.schemaPreview).toHaveBeenCalled());

    fireEvent.change(screen.getByTestId('circuit-profile-interval-instantaneous'), { target: { value: '30' } });
    const prepared = within(screen.getByTestId('circuit-prepared-group-instantaneous'));
    await waitFor(() => expect(prepared.getByTestId('group-schema-blocked')).toHaveTextContent('step4.group.schema.blocked_dirty'));
    expect(prepared.getByTestId('group-schema-apply')).toBeDisabled();
    expect(groupsAPI.schemaApply).not.toHaveBeenCalled();
    expect(startAPI.start).not.toHaveBeenCalled();
  });

  it('shows an incompatible output as a destination blocker, with nothing confirmed or started', async () => {
    backend.rejectPreview = Object.assign(new Error('incompatible'), { code: 'RECORDING_SCHEMA_INCOMPATIBLE', status: 409 });
    await prepare(purpose);

    await waitFor(() => expect(screen.getByTestId('circuit-profile-blockers')).toBeInTheDocument());
    expect(screen.getByTestId('circuit-blocker-schema_incompatible-0')).toBeInTheDocument();
    expect(screen.queryByTestId('group-schema-apply')).toBeNull();
    expect(groupsAPI.schemaApply).not.toHaveBeenCalled();
    expect(startAPI.start).not.toHaveBeenCalled();
  });
});
