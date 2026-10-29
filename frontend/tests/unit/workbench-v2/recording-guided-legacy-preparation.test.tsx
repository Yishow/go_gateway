import { onlineManager, QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { BasicRecordingPanel } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingPanel';
import { BasicSchemaPreparation } from '@/features/datalink/workbench-v2/steps/step4/BasicSchemaPreparation';
import { BasicRecordingPreparation } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingPreparation';
import { BasicMemberUpdate } from '@/features/datalink/workbench-v2/steps/step4/BasicMemberUpdate';
import type { RecordingPreparationCommandFence } from '@/features/datalink/workbench-v2/steps/step4/recordingPreparationPresentation';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { readSchemaOperation, retainSchemaOperation } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery';
import { basicManagedGroupDraft, basicRecordingScopeKey, loadBasicRecordingIntent, saveBasicRecordingIntent } from '@/features/datalink/workbench-v2/state/basicRecording';
import { buildBasicStartRequest } from '@/features/datalink/workbench-v2/steps/step4/basicStartRequest';
import { studioV2WorkspaceWriteGroupsAPI as groupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WriteGroupLifecycleAPI as lifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
import { studioV2RecordingStartAPI as startAPI } from '@/services/studioV2RecordingStart';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { RecordingStartOperation } from '@/types/studioV2RecordingStart';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import type { WriteGroupLifecycle } from '@/types/studioV2WriteGroupLifecycle';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import { managedSchemaPreviewToken } from '../../fixtures/managedSchemaPreview';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  list: vi.fn(), update: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(),
} }));
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: vi.fn() } }));
vi.mock('@/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { ensureBasic: vi.fn(), start: vi.fn(), operation: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { list: vi.fn(async () => []) } }));
let group: WriteGroup, revision: string, schemaReady: boolean, editable: boolean, foreignLifecycle: boolean;
function lifecycle(): WriteGroupLifecycle {
  return { workspace_id: 'ws-1', workspace_revision: foreignLifecycle ? 'foreign-revision' : revision,
    group_id: group.id, group_revision: group.revision, applied_revision: group.applied_revision, status: group.status,
    connector_id: group.destination.connector_id, connector_revision: group.destination.connector_revision,
    actions: [{ action: 'edit', allowed: editable, reason_code: editable ? '' : 'ownership_unavailable', next_action: 'reload', requires_review: false }],
    history: { canonical: true, current_group_id: group.id, generations: [] } };
}
beforeEach(() => {
  vi.restoreAllMocks(); vi.clearAllMocks(); retainSchemaOperation('ws-1', 'group-1', undefined, readSchemaOperation('ws-1', 'group-1')); sessionStorage.clear(); localStorage.clear();
  Object.defineProperty(navigator, 'onLine', { configurable: true, value: true }); onlineManager.setOnline(true);
  revision = 'workspace-1'; schemaReady = false; editable = true; foreignLifecycle = false;
  const base = savedGroup();
  group = savedGroup({ basic_managed_device_id: 'dev-1', layout_version: 'managed_wide_v1',
    members: [...base.members, { device_id: 'dev-1', point_id: 'pt-2', tag_id: 'tag-2', source_revision: 's2', mapping_revision: 'm2', target_column: 'batch', required: true }],
    destination: { ...base.destination, storage_strategy: 'managed' } });
  vi.mocked(groupsAPI.list).mockImplementation(async () => ({ workspace_id: 'ws-1', workspace_revision: revision, groups: [structuredClone(group)] }));
  vi.mocked(groupsAPI.delivery).mockRejectedValue(new Error('unavailable delivery'));
  vi.mocked(groupsAPI.readiness).mockImplementation(async () => ({ workspace_id: 'ws-1', workspace_revision: revision,
    group_id: group.id, group_revision: group.revision, applied_revision: group.applied_revision,
    config_ready: true, schema_ready: schemaReady, ready: schemaReady,
    issues: schemaReady ? [] : [{ code: 'schema-confirmation-required', severity: 'blocking', step: 'Step 4', scope: '', message: '' }] }));
  vi.mocked(lifecycleAPI.lifecycle).mockImplementation(async () => lifecycle());
  vi.mocked(groupsAPI.update).mockImplementation(async (id, request) => {
    expect(id).toBe(group.id); expect(request.expected_workspace_revision).toBe(revision);
    expect(request.expected_group_revision).toBe(group.revision); expect(request.expected_connector_revision).toBe('crev-1');
    group = { ...group, ...request.group, revision: 'rev-2', members: group.members,
      destination: { ...group.destination, ...request.group.destination } };
    revision = 'workspace-2'; schemaReady = false;
    return { workspace_revision: revision, group: structuredClone(group) };
  });
  vi.mocked(groupsAPI.schemaPreview).mockImplementation(async (id, request) => {
    expect(id).toBe(group.id); expect(request.expected_workspace_revision).toBe(revision); expect(request.expected_group_revision).toBe(group.revision);
    return { ...managedSchemaPreviewToken, workspace_id: 'ws-1', workspace_revision: revision, plan_id: group.id, plan_revision: group.revision };
  });
  vi.mocked(groupsAPI.schemaApply).mockImplementation(async (_id, request) => {
    expect(request.operation_id).toBe('operation-1'); expect(request.token).toBe('token-1');
    schemaReady = true; revision = 'workspace-3';
    group = { ...group, destination: { ...group.destination, schema_revision: group.revision, schema_digest: 'digest-1' } };
    return { operation_id: 'operation-1', status: 'succeeded', executed_statements: 1, created_at: '', updated_at: '' };
  });
  vi.mocked(startAPI.start).mockImplementation(async request => {
    expect(request.expected_workspace_revision).toBe(revision); expect(request.groups[0].expected_group_revision).toBe(group.revision);
    expect(request.groups[0].draft).toBeUndefined();
    return { operation_id: 'start-1', action: 'recording_start', status: 'succeeded', stage: 'complete', intent_digest: 'digest',
      workspace_id: 'ws-1', setup_revision: revision, device_ids: ['dev-1'], groups: [], devices: [], created_at: '', updated_at: '' };
  });
});
afterEach(() => { cleanup(); vi.restoreAllMocks(); Object.defineProperty(navigator, 'onLine', { configurable: true, value: true }); window.dispatchEvent(new Event('online')); onlineManager.setOnline(true); });
async function mount() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false, networkMode: 'always' } } });
  const result = render(<QueryClientProvider client={client}><BasicRecordingPanel state={savedGroupState()} workspaceId="ws-1" readonly={false} /></QueryClientProvider>);
  fireEvent.click(await screen.findByTestId('device-recording-edit-existing'));
  await screen.findByTestId('basic-recording-panel');
  return { ...result, client };
}
it('uses one primary and makes legacy navigation and native control edits mutation-free', async () => {
  await mount();
  fireEvent.change(screen.getByTestId('basic-recording-interval'), { target: { value: '30' } });
  fireEvent.change(screen.getByTestId('basic-recording-completeness'), { target: { value: 'partial' } });
  expect(document.querySelectorAll('button[data-primary="true"]')).toHaveLength(1);
  expect(groupsAPI.update).not.toHaveBeenCalled(); expect(groupsAPI.schemaPreview).not.toHaveBeenCalled();
  expect(groupsAPI.schemaApply).not.toHaveBeenCalled(); expect(startAPI.start).not.toHaveBeenCalled();
});
it('saves native legacy interval and completeness only on explicit prepare, then previews the canonical revisions separately', async () => {
  await mount(); const original = structuredClone(group);
  fireEvent.change(screen.getByTestId('basic-recording-interval'), { target: { value: '30' } });
  fireEvent.change(screen.getByTestId('basic-recording-completeness'), { target: { value: 'partial' } });
  await waitFor(() => expect(screen.getByTestId('basic-recording-prepare')).toBeEnabled());
  fireEvent.click(screen.getByTestId('basic-recording-prepare'));
  await waitFor(() => expect(groupsAPI.update).toHaveBeenCalledTimes(1));
  expect(groupsAPI.update).toHaveBeenCalledWith(group.id, expect.objectContaining({ group: basicManagedGroupDraft(original, 30, 'partial') }));
  expect(groupsAPI.schemaPreview).not.toHaveBeenCalled(); expect(groupsAPI.schemaApply).not.toHaveBeenCalled(); expect(startAPI.start).not.toHaveBeenCalled();
  await waitFor(() => expect(screen.getByTestId('group-schema-preview')).toBeEnabled()); fireEvent.click(screen.getByTestId('group-schema-preview'));
  await screen.findByTestId('group-schema-preview-result');
  expect(groupsAPI.schemaPreview).toHaveBeenCalledTimes(1); expect(groupsAPI.schemaApply).not.toHaveBeenCalled(); expect(startAPI.start).not.toHaveBeenCalled();
});

it('keeps unchanged preparation write-free and requires separate schema consent and final start consent', async () => {
  await mount();
  await waitFor(() => expect(screen.getByTestId('group-schema-preview')).toBeEnabled());
  fireEvent.click(screen.getByTestId('group-schema-preview'));
  await screen.findByTestId('group-schema-preview-result');
  expect(groupsAPI.update).not.toHaveBeenCalled(); expect(startAPI.ensureBasic).not.toHaveBeenCalled();
  expect(groupsAPI.schemaApply).not.toHaveBeenCalled(); expect(startAPI.start).not.toHaveBeenCalled();
  fireEvent.click(screen.getByTestId('group-schema-apply'));
  await waitFor(() => expect(screen.getByTestId('basic-recording-start')).toBeEnabled());
  expect(groupsAPI.schemaApply).toHaveBeenCalledTimes(1); expect(startAPI.start).not.toHaveBeenCalled();
  fireEvent.click(screen.getByTestId('basic-recording-start'));
  expect(screen.getByRole('dialog')).toBeVisible(); expect(startAPI.start).not.toHaveBeenCalled();
  fireEvent.click(screen.getByTestId('recording-start-cancel'));
  expect(screen.queryByRole('dialog')).not.toBeInTheDocument(); expect(startAPI.start).not.toHaveBeenCalled();
  fireEvent.click(screen.getByTestId('basic-recording-start'));
  fireEvent.click(screen.getByTestId('recording-start-confirm'));
  await waitFor(() => expect(startAPI.start).toHaveBeenCalledTimes(1));
  expect(groupsAPI.update).not.toHaveBeenCalled();
});
it.each(['interval', 'completeness'] as const)('invalidates open start consent when native %s changes without a saved revision change', async control => {
  schemaReady = true; await mount();
  await waitFor(() => expect(screen.getByTestId('basic-recording-start')).toBeEnabled());
  fireEvent.click(screen.getByTestId('basic-recording-start')); const oldConfirm = screen.getByTestId('recording-start-confirm');
  fireEvent.change(screen.getByTestId(`basic-recording-${control}`), { target: { value: control === 'interval' ? '20' : 'partial' } });
  expect(screen.queryByRole('dialog')).not.toBeInTheDocument(); fireEvent.click(oldConfirm);
  expect(startAPI.start).not.toHaveBeenCalled(); expect(groupsAPI.update).not.toHaveBeenCalled(); expect(groupsAPI.schemaApply).not.toHaveBeenCalled();
});
it.each(['denied', 'foreign'] as const)('cannot save dirty controls under a %s lifecycle edit proof', async kind => {
  editable = kind !== 'denied'; foreignLifecycle = kind === 'foreign';
  await mount(); fireEvent.change(screen.getByTestId('basic-recording-interval'), { target: { value: '30' } });
  await waitFor(() => expect(screen.getByTestId('basic-recording-prepare')).toBeDisabled());
  fireEvent.click(screen.getByTestId('basic-recording-prepare'));
  expect(groupsAPI.update).not.toHaveBeenCalled(); expect(groupsAPI.schemaPreview).not.toHaveBeenCalled(); expect(startAPI.start).not.toHaveBeenCalled();
});
it('waits for canonical saved revisions before allowing the separate preview command', async () => {
  await mount();
  let release!: (value: Awaited<ReturnType<typeof groupsAPI.list>>) => void;
  vi.mocked(groupsAPI.list).mockImplementationOnce(() => new Promise(resolve => { release = resolve; }));
  fireEvent.change(screen.getByTestId('basic-recording-interval'), { target: { value: '30' } });
  await waitFor(() => expect(screen.getByTestId('basic-recording-prepare')).toBeEnabled()); fireEvent.click(screen.getByTestId('basic-recording-prepare'));
  await waitFor(() => expect(groupsAPI.update).toHaveBeenCalledTimes(1));
  expect(screen.queryByTestId('group-schema-preview')).not.toBeInTheDocument(); expect(groupsAPI.schemaPreview).not.toHaveBeenCalled();
  await act(async () => { release({ workspace_id: 'ws-1', workspace_revision: revision, groups: [structuredClone(group)] }); });
  await waitFor(() => expect(screen.getByTestId('group-schema-preview')).toBeEnabled());
  fireEvent.click(screen.getByTestId('group-schema-preview'));
  await waitFor(() => expect(groupsAPI.schemaPreview).toHaveBeenCalledWith(group.id, expect.objectContaining({ expected_group_revision: 'rev-2', expected_workspace_revision: 'workspace-2' })));
});
it('keeps a retained unknown start on original O1 and never saves, previews, applies or resends it', async () => {
  schemaReady = true;
  const scope = basicRecordingScopeKey('ws-1', 'dev-1', 'managed-recording');
  const request = buildBasicStartRequest({ workspaceId: 'ws-1', deviceId: 'dev-1', connectorRevision: 'crev-1' }, group, revision, 'original-request');
  saveBasicRecordingIntent({ scope_key: scope, request, operation_id: 'O1' });
  vi.mocked(startAPI.operation).mockResolvedValue({ operation_id: 'O1', action: 'recording_start', status: 'unknown', stage: 'activation', intent_digest: 'digest',
    workspace_id: 'ws-1', setup_revision: revision, device_ids: ['dev-1'], groups: [], devices: [], created_at: '', updated_at: '' });
  await mount();
  await waitFor(() => expect(screen.getByTestId('basic-recording-check')).toBeEnabled()); fireEvent.click(screen.getByTestId('basic-recording-check'));
  await waitFor(() => expect(vi.mocked(startAPI.operation).mock.calls.length).toBeGreaterThanOrEqual(2));
  expect(vi.mocked(startAPI.operation).mock.calls.every(([id]) => id === 'O1')).toBe(true);
  expect(screen.queryByTestId('basic-recording-retry')).not.toBeInTheDocument();
  expect(groupsAPI.update).not.toHaveBeenCalled(); expect(groupsAPI.schemaPreview).not.toHaveBeenCalled(); expect(groupsAPI.schemaApply).not.toHaveBeenCalled(); expect(startAPI.start).not.toHaveBeenCalled();
});
it('retains a lost-ACK request without inventing an operation ID or resending start', async () => {
  const scope = basicRecordingScopeKey('ws-1', 'dev-1', 'managed-recording');
  saveBasicRecordingIntent({ scope_key: scope, request: buildBasicStartRequest({ workspaceId: 'ws-1', deviceId: 'dev-1', connectorRevision: 'crev-1' }, group, revision, 'original-request') });
  await mount();
  fireEvent.click(screen.getByTestId('recording-next-action'));
  expect(startAPI.operation).not.toHaveBeenCalled(); expect(startAPI.start).not.toHaveBeenCalled();
  expect(groupsAPI.schemaPreview).not.toHaveBeenCalled(); expect(groupsAPI.update).not.toHaveBeenCalled();
});
it('rechecks unreadable intent storage at final confirmation without sending start', async () => {
  schemaReady = true; await mount();
  await waitFor(() => expect(screen.getByTestId('basic-recording-start')).toBeEnabled()); fireEvent.click(screen.getByTestId('basic-recording-start'));
  vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => { throw new Error('private storage failure'); });
  fireEvent.click(screen.getByTestId('recording-start-confirm')); await act(async () => undefined);
  expect(startAPI.start).not.toHaveBeenCalled(); expect(groupsAPI.schemaApply).not.toHaveBeenCalled(); expect(document.body.textContent).not.toContain('private storage failure');
});
it('keeps offline local review available without network mutations', async () => {
  await mount();
  Object.defineProperty(navigator, 'onLine', { configurable: true, value: false }); fireEvent(window, new Event('offline'));
  fireEvent.change(screen.getByTestId('basic-recording-interval'), { target: { value: '30' } });
  await waitFor(() => expect(screen.getByTestId('basic-recording-prepare')).toBeDisabled()); fireEvent.click(screen.getByTestId('basic-recording-prepare'));
  expect(groupsAPI.update).not.toHaveBeenCalled(); expect(groupsAPI.schemaPreview).not.toHaveBeenCalled(); expect(startAPI.start).not.toHaveBeenCalled();
  expect(screen.getByTestId('group-schema-panel')).toBeInTheDocument();
});

it('preserves standalone explicit preparation and its existing auto-preview, never implicit DDL', async () => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false, networkMode: 'always' } } });
  const ensure = vi.fn(async () => ({ workspace_revision: revision, group }));
  render(<QueryClientProvider client={client}><BasicSchemaPreparation group={group} workspaceId="ws-1" workspaceRevision={revision}
    disabled={false} readonly={false} destinationMatches ensureGroup={ensure} onApplied={() => undefined} /></QueryClientProvider>);
  expect(groupsAPI.schemaPreview).not.toHaveBeenCalled();
  fireEvent.click(screen.getByTestId('basic-recording-prepare'));
  await screen.findByTestId('group-schema-preview-result');
  expect(ensure).toHaveBeenCalledTimes(1); expect(groupsAPI.schemaPreview).toHaveBeenCalledTimes(1);
  expect(groupsAPI.schemaApply).not.toHaveBeenCalled(); expect(startAPI.start).not.toHaveBeenCalled();
});
it('rejects a ready response for an old workspace revision instead of offering start review', async () => {
  schemaReady = true;
  vi.mocked(groupsAPI.readiness).mockImplementation(async () => ({ workspace_id: 'ws-1', workspace_revision: 'old-workspace',
    group_id: group.id, group_revision: group.revision, applied_revision: group.applied_revision, config_ready: true, schema_ready: true, ready: true, issues: [] }));
  await mount(); await waitFor(() => expect(groupsAPI.readiness).toHaveBeenCalled());
  expect(screen.queryByTestId('basic-recording-start')).not.toBeInTheDocument();
  expect(screen.queryByTestId('recording-start-confirm')).not.toBeInTheDocument();
  expect(startAPI.start).not.toHaveBeenCalled(); expect(groupsAPI.schemaApply).not.toHaveBeenCalled();
});
it('shows a safe preparation failure without previewing, applying, starting or automatically retrying a failed save', async () => {
  vi.mocked(groupsAPI.update).mockRejectedValue(new Error('postgres://private.example/secret'));
  await mount(); fireEvent.change(screen.getByTestId('basic-recording-interval'), { target: { value: '30' } });
  await waitFor(() => expect(screen.getByTestId('basic-recording-prepare')).toBeEnabled()); fireEvent.click(screen.getByTestId('basic-recording-prepare'));
  expect(await screen.findByText('step4.basic.prepare_failed')).toBeInTheDocument();
  expect(groupsAPI.update).toHaveBeenCalledTimes(1); expect(document.body.textContent).not.toContain('private.example');
  expect(groupsAPI.schemaPreview).not.toHaveBeenCalled(); expect(groupsAPI.schemaApply).not.toHaveBeenCalled(); expect(startAPI.start).not.toHaveBeenCalled();
});

it('keeps secondary member diff inspection available but blocks saving while the selected schema owner retains O1', async () => {
  group = { ...group, members: group.members.slice(0, 2).map((member, index) => ({ ...member, target_column: `v_${String(index).repeat(24)}` })) };
  retainSchemaOperation('ws-1', group.id, 'schema-O1');
  await mount();
  await waitFor(() => expect(screen.getByTestId('basic-recording-members-review')).toBeEnabled());
  fireEvent.click(screen.getByTestId('basic-recording-members-review'));
  expect(screen.getByTestId('basic-members-confirm')).toBeDisabled(); fireEvent.click(screen.getByTestId('basic-members-confirm'));
  expect(groupsAPI.update).not.toHaveBeenCalled(); expect(groupsAPI.schemaPreview).not.toHaveBeenCalled(); expect(startAPI.start).not.toHaveBeenCalled();
});

it('rechecks late schema O1 at the actual member confirmation instead of trusting the earlier enabled button', async () => {
  group = { ...group, members: group.members.slice(0, 2).map((member, index) => ({ ...member, target_column: `v_${String(index).repeat(24)}` })) }; await mount();
  await waitFor(() => expect(screen.getByTestId('basic-recording-members-review')).toBeEnabled());
  fireEvent.click(screen.getByTestId('basic-recording-members-review'));
  const confirm = screen.getByTestId('basic-members-confirm'); expect(confirm).toBeEnabled();
  retainSchemaOperation('ws-1', group.id, 'late-O1');
  fireEvent.click(confirm); await act(async () => undefined);
  expect(groupsAPI.update).not.toHaveBeenCalled(); expect(groupsAPI.schemaApply).not.toHaveBeenCalled(); expect(startAPI.start).not.toHaveBeenCalled();
});
it('rejects a retired owner fence after A → B → A at the real secondary save handler', async () => {
  group = { ...group, members: group.members.slice(0, 2).map((member, index) => ({ ...member, target_column: `v_${String(index).repeat(24)}` })) };
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false, networkMode: 'always' } } });
  const candidates = buildGroupCandidates(savedGroupState()).candidates;
  const ownerKey = basicRecordingScopeKey('ws-1', 'dev-1', 'managed-recording');
  let fence: RecordingPreparationCommandFence | null = null;
  const report = (value: RecordingPreparationCommandFence | null) => { fence = value; };
  let gate = () => false;
  const tree = (workspaceRevision: string) => <QueryClientProvider client={client}>
    <BasicRecordingPreparation workspaceId="ws-1" workspaceRevision={workspaceRevision} deviceId="dev-1" ownerKey={ownerKey}
      fenceKey={workspaceRevision} onCommandFence={report} group={group} candidates={candidates} intervalSeconds={10} incompletePolicy="skip_row"
      disabled={false} readonly={false} destinationMatches ensureGroup={async () => ({ workspace_revision: revision, group })}
      onPrepared={async () => undefined} start={async () => undefined} check={async () => undefined}
      live={{ unresolved: false, queryFailed: false, busy: false }} />
    <BasicMemberUpdate group={group} deviceId="dev-1" candidates={candidates} workspaceRevision={workspaceRevision}
      disabled={false} isBlocked={() => gate()} onUpdated={() => undefined} />
  </QueryClientProvider>;
  const view = render(tree('workspace-1'));
  await waitFor(() => expect(fence?.blocked).toBe(false)); const retired = fence!;
  view.rerender(tree('workspace-2')); view.rerender(tree('workspace-1'));
  expect(retired.isBlocked(group.id)).toBe(true);
  fireEvent.click(screen.getByTestId('basic-recording-members-review'));
  const confirm = screen.getByTestId('basic-members-confirm'); expect(confirm).toBeEnabled();
  gate = () => retired.isBlocked(group.id); fireEvent.click(confirm);
  await act(async () => undefined); expect(groupsAPI.update).not.toHaveBeenCalled();
});
it('fails closed when the optional secondary owner fence cannot be read', () => {
  group = { ...group, members: group.members.slice(0, 2).map((member, index) => ({ ...member, target_column: `v_${String(index).repeat(24)}` })) };
  const client = new QueryClient();
  expect(() => render(<QueryClientProvider client={client}><BasicMemberUpdate group={group} deviceId="dev-1"
    candidates={buildGroupCandidates(savedGroupState()).candidates} workspaceRevision={revision} disabled={false}
    isBlocked={() => { throw new Error('private fence failure'); }} onUpdated={() => undefined} /></QueryClientProvider>)).not.toThrow();
  fireEvent.click(screen.getByTestId('basic-recording-members-review'));
  expect(screen.getByTestId('basic-members-confirm')).toBeDisabled();
  expect(groupsAPI.update).not.toHaveBeenCalled(); expect(document.body.textContent).not.toContain('private fence failure');
});
it('keeps custom member repair navigation available while new mutations are fenced', () => {
  const onOpenAdvanced = vi.fn(), client = new QueryClient();
  render(<QueryClientProvider client={client}><BasicMemberUpdate group={group} deviceId="dev-1"
    candidates={buildGroupCandidates(savedGroupState()).candidates} workspaceRevision={revision} disabled={false}
    isBlocked={() => true} onUpdated={() => undefined} onOpenAdvanced={onOpenAdvanced} /></QueryClientProvider>);
  fireEvent.click(screen.getByTestId('basic-recording-members-review'));
  expect(onOpenAdvanced).toHaveBeenCalledTimes(1); expect(groupsAPI.update).not.toHaveBeenCalled();
});
it('allows an explicit original-request retry only after fresh matching partial proof, retaining historical request revisions', async () => {
  schemaReady = true;
  const scope = basicRecordingScopeKey('ws-1', 'dev-1', 'managed-recording');
  const request = buildBasicStartRequest({ workspaceId: 'ws-1', deviceId: 'dev-1', connectorRevision: 'crev-1' }, { ...group, revision: 'historical-rev' }, 'historical-workspace', 'original-request');
  saveBasicRecordingIntent({ scope_key: scope, request, operation_id: 'O1' });
  const partial = { operation_id: 'O1', action: 'recording_start' as const, status: 'partial' as const, next_action: 'retry_same_request', stage: 'activation' as const,
    intent_digest: 'digest', workspace_id: 'ws-1', setup_revision: revision, device_ids: ['dev-1'],
    groups: [{ group_id: group.id, group_revision: group.revision, saved: true, applied: true, ready: true }],
    devices: [{ device_id: 'dev-1', activated: false }], created_at: '', updated_at: '' };
  vi.mocked(startAPI.operation).mockResolvedValue(partial); vi.mocked(startAPI.start).mockResolvedValue(partial);
  await mount();
  const retry = await screen.findByTestId('basic-recording-retry'); expect(retry).toBeEnabled();
  expect(startAPI.start).not.toHaveBeenCalled(); const reads = vi.mocked(startAPI.operation).mock.calls.length;
  fireEvent.click(retry);
  await waitFor(() => expect(startAPI.start).toHaveBeenCalledTimes(1));
  expect(vi.mocked(startAPI.operation).mock.calls.length).toBeGreaterThan(reads);
  expect(startAPI.start).toHaveBeenCalledWith(request);
  expect(groupsAPI.update).not.toHaveBeenCalled(); expect(groupsAPI.schemaPreview).not.toHaveBeenCalled(); expect(groupsAPI.schemaApply).not.toHaveBeenCalled();
});

function retainedRetry() {
  schemaReady = true;
  const scope = basicRecordingScopeKey('ws-1', 'dev-1', 'managed-recording');
  const request = buildBasicStartRequest({ workspaceId: 'ws-1', deviceId: 'dev-1', connectorRevision: 'crev-1' }, group, revision, 'original-request');
  const partial: RecordingStartOperation = { operation_id: 'O1', action: 'recording_start', status: 'partial', next_action: 'retry_same_request', stage: 'activation',
    intent_digest: 'digest', workspace_id: 'ws-1', setup_revision: revision, device_ids: ['dev-1'],
    groups: [{ group_id: group.id, group_revision: group.revision, saved: true, applied: true, ready: true }],
    devices: [{ device_id: 'dev-1', activated: false }], created_at: '', updated_at: '' };
  saveBasicRecordingIntent({ scope_key: scope, request, operation_id: 'O1' });
  vi.mocked(startAPI.operation).mockReset().mockResolvedValue(partial);
  return { scope, request, partial };
}
it.each(['unknown', 'query_failed', 'foreign_operation', 'foreign_workspace', 'foreign_group', 'wrong_action'] as const)('does not resend the original request when fresh retry proof is %s', async kind => {
  const { scope, request, partial } = retainedRetry(); await mount();
  const retry = await screen.findByTestId('basic-recording-retry'); await waitFor(() => expect(retry).toBeEnabled());
  if (kind === 'query_failed') vi.mocked(startAPI.operation).mockRejectedValueOnce(new Error('private query failure'));
  else vi.mocked(startAPI.operation).mockResolvedValueOnce({ ...partial,
    ...(kind === 'unknown' ? { status: 'unknown' } : kind === 'foreign_operation' ? { operation_id: 'O2' }
      : kind === 'foreign_workspace' ? { workspace_id: 'foreign' } : kind === 'foreign_group' ? { groups: [{ ...partial.groups[0], group_id: 'foreign' }] } : { next_action: 'review_request' }) });
  const reads = vi.mocked(startAPI.operation).mock.calls.length; fireEvent.click(retry);
  await waitFor(() => expect(vi.mocked(startAPI.operation).mock.calls.length).toBeGreaterThan(reads)); await act(async () => undefined);
  expect(startAPI.start).not.toHaveBeenCalled(); expect(loadBasicRecordingIntent(scope)?.request).toEqual(request);
  expect(groupsAPI.update).not.toHaveBeenCalled(); expect(groupsAPI.schemaApply).not.toHaveBeenCalled(); expect(document.body.textContent).not.toContain('private query failure');
});
it('rechecks other-owner O1 after a fresh original-request recovery lookup', async () => {
  const { partial } = retainedRetry(); await mount(); const retry = await screen.findByTestId('basic-recording-retry');
  vi.mocked(startAPI.operation).mockImplementationOnce(async () => { retainSchemaOperation('ws-1', group.id, 'schema-O1'); return partial; });
  fireEvent.click(retry); await waitFor(() => expect(vi.mocked(startAPI.operation).mock.calls.length).toBeGreaterThanOrEqual(2));
  await act(async () => undefined); expect(startAPI.start).not.toHaveBeenCalled();
});
it.each(['round_trip', 'unmount', 'foreign_cache', 'readonly', 'offline'] as const)('rejects an old recovery continuation after %s while preserving the original request', async change => {
  const { scope, request, partial } = retainedRetry(); const mounted = await mount();
  const retry = await screen.findByTestId('basic-recording-retry'); let resolve!: (value: RecordingStartOperation) => void;
  vi.mocked(startAPI.operation).mockImplementationOnce(() => new Promise(accept => { resolve = accept; }));
  fireEvent.click(retry); await waitFor(() => expect(vi.mocked(startAPI.operation).mock.calls.length).toBeGreaterThanOrEqual(2));
  if (change === 'unmount') mounted.unmount();
  else if (change === 'foreign_cache') mounted.client.setQueryData(studioV2WorkspaceKeys.writeGroups(), { workspace_id: 'foreign', workspace_revision: revision, groups: [group] });
  else if (change === 'readonly') mounted.rerender(<QueryClientProvider client={mounted.client}><BasicRecordingPanel state={savedGroupState()} workspaceId="ws-1" readonly /></QueryClientProvider>);
  else if (change === 'offline') { Object.defineProperty(navigator, 'onLine', { configurable: true, value: false }); window.dispatchEvent(new Event('offline')); }
  else {
    await act(async () => { mounted.client.setQueryData(studioV2WorkspaceKeys.writeGroups(), { workspace_id: 'ws-1', workspace_revision: 'temporary-revision', groups: [group] }); });
    await act(async () => { mounted.client.setQueryData(studioV2WorkspaceKeys.writeGroups(), { workspace_id: 'ws-1', workspace_revision: revision, groups: [group] }); });
  }
  await act(async () => { resolve(partial); }); expect(startAPI.start).not.toHaveBeenCalled();
  expect(loadBasicRecordingIntent(scope)?.request).toEqual(request);
});

it('keeps failed retry_same_request on the original-operation query primary, not a new start review', async () => {
  const { partial } = retainedRetry(); vi.mocked(startAPI.operation).mockResolvedValue({ ...partial, status: 'failed' });
  await mount();
  await waitFor(() => expect(screen.getByTestId('basic-recording-retry')).toBeEnabled());
  expect(screen.getByTestId('basic-recording-check')).toHaveAttribute('data-primary', 'true');
  expect(screen.queryByTestId('basic-recording-start')).not.toBeInTheDocument();
  expect(startAPI.start).not.toHaveBeenCalled();
});

it.each(['start_confirmation', 'dirty_prepare'] as const)('rejects synchronous foreign canonical cache before %s', async command => {
  schemaReady = command === 'start_confirmation'; const mounted = await mount();
  let button: HTMLElement;
  if (command === 'start_confirmation') {
    await waitFor(() => expect(screen.getByTestId('basic-recording-start')).toBeEnabled());
    fireEvent.click(screen.getByTestId('basic-recording-start')); button = screen.getByTestId('recording-start-confirm');
  } else {
    fireEvent.change(screen.getByTestId('basic-recording-interval'), { target: { value: '30' } });
    await waitFor(() => expect(screen.getByTestId('basic-recording-prepare')).toBeEnabled()); button = screen.getByTestId('basic-recording-prepare');
  }
  mounted.client.setQueryData(studioV2WorkspaceKeys.writeGroups(), { workspace_id: 'foreign', workspace_revision: revision, groups: [group] });
  fireEvent.click(button); await act(async () => undefined);
  expect(startAPI.start).not.toHaveBeenCalled(); expect(groupsAPI.update).not.toHaveBeenCalled();
  expect(groupsAPI.schemaPreview).not.toHaveBeenCalled(); expect(groupsAPI.schemaApply).not.toHaveBeenCalled();
});
