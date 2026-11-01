import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { GroupLifecycleBar } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar';
import type { GroupEditorProps } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { RecordingDraftNavigation } from '@/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import { studioV2LocatorHydration } from '@/pages/datalink/workbench-v2/studioV2LocatorHydration';
import { readStudioLocatorEditorEntry } from '@/services/studioV2LocatorEntry';
import { studioV2WriteGroupLifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
import { studioV2WriteGroupSuccessorAPI } from '@/services/studioV2WriteGroupSuccessor';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { dbTargetAPI } from '@/services/datalink';
import type { StudioLocatorEntryRead } from '@/types/studioV2LocatorEntry';
import { recordingChangeKey, saveRecordingChange } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/recordingChangeRecovery';
import { fixture } from '../../fixtures/recordingChange';
import { locatorRead } from './studio-locator-fixtures';
const captured = vi.hoisted(() => ({ editor: undefined as GroupEditorProps | undefined }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2LocatorEntry', () => ({ readStudioLocatorEditorEntry: vi.fn() }));
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: vi.fn() } }));
vi.mock('@/services/studioV2WriteGroupSuccessor', () => ({ studioV2WriteGroupSuccessorAPI: {
  review: vi.fn(), reserve: vi.fn(), schemaPreview: vi.fn(), schemaConfirm: vi.fn(), activationReview: vi.fn(), activationReadiness: vi.fn(), confirmCutover: vi.fn(), confirmFirstStart: vi.fn(), operationStatus: vi.fn(),
} }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn(), get: vi.fn(), delivery: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn() } }));
vi.mock('@/services/datalink', async original => ({ ...await original<typeof import('@/services/datalink')>(), dbTargetAPI: { listConnectors: vi.fn() } }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor', () => ({ GroupEditor: (props: GroupEditorProps) => {
  captured.editor = props;
  return <div data-testid="canonical-editor"><input aria-label="Retained editor" defaultValue="same draft" />{props.group && <GroupLifecycleBar group={props.group}
    workspaceId={props.workspaceId} workspaceRevision={props.workspaceRevision} readonly={props.readonly} canApply={false}
    locatorScope={props.locatorScope} locatorGuard={props.locatorGuard} lifecycleHost={props.lifecycleHost} onChanged={props.onSaved} onConflict={props.onReload} />}</div>;
} }));
let f = fixture();
const clients: QueryClient[] = [];
function receipt(): StudioLocatorEntryRead {
  const read = locatorRead(4, f.group.id);
  read.locator.workspace_id = f.workspace; read.workspace.id = f.workspace;
  read.entry.workspace_revision = 'W1'; read.entry.group_heads = { [f.group.id]: f.group.revision };
  read.entry.write_groups = { workspace_id: f.workspace, workspace_revision: 'W1', groups: [f.group] };
  read.selection = { state: 'selected', target: { step: 4, workspaceId: f.workspace, workspaceRevision: 'W1', groupId: f.group.id, groupRevision: f.group.revision } };
  return { ...read, entry: { ...read.entry, editor_hydrated: true, database_config: null, database_targets: [] } };
}
function mount() {
  const route = receipt(), client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } }); clients.push(client);
  const content = (entry = route, suspended = false) => <QueryClientProvider client={client}><RecordingDraftNavigation><Step4Database
    state={studioV2LocatorHydration(route)!} dispatch={vi.fn()} workspaceId={f.workspace} presentation="simple" initialSelection="located"
    locatorEntry={entry} locatorOnly locatorSuspended={suspended} /></RecordingDraftNavigation></QueryClientProvider>;
  return { ...render(content()), content, route };
}
async function open() { const view = mount(); await screen.findByTestId('canonical-editor'); await waitFor(() => expect(screen.getByTestId('saved-recording-transition')).toBeEnabled()); return view; }
function noCommands() { for (const key of ['review', 'reserve', 'schemaPreview', 'schemaConfirm', 'activationReview', 'activationReadiness', 'confirmCutover', 'confirmFirstStart'] as const) expect(studioV2WriteGroupSuccessorAPI[key]).not.toHaveBeenCalled(); }
beforeEach(() => {
  vi.resetAllMocks(); captured.editor = undefined; sessionStorage.clear(); f = fixture();
  f.group.status = 'ready'; f.group.applied_revision = 'R1'; f.group.layout_version = 'managed_wide_profile_v1';
  f.group.recording_profile = { kind: 'instantaneous', interval: 15, unit: 'second' }; delete f.group.basic_managed_device_id;
  f.lifecycle.status = 'ready'; f.lifecycle.applied_revision = 'R1'; f.lifecycle.actions = [{ action: 'transition', allowed: true, reason_code: '', next_action: '', requires_review: true }];
  Object.assign(f.seed, { layout_version: f.group.layout_version, recording_profile: f.group.recording_profile });
  vi.mocked(readStudioLocatorEditorEntry).mockImplementation(async () => ({ ok: true, value: receipt() }));
  vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockImplementation(f.lifecycleRead);
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.get).mockImplementation(f.deps.getGroup);
  vi.mocked(dbTargetAPI.listConnectors).mockResolvedValue([]);
  vi.mocked(studioV2WriteGroupSuccessorAPI.review).mockImplementation(f.api.review);
  vi.mocked(studioV2WriteGroupSuccessorAPI.reserve).mockImplementation(f.api.reserve);
  vi.mocked(studioV2WriteGroupSuccessorAPI.schemaPreview).mockImplementation(f.api.schemaPreview);
  vi.mocked(studioV2WriteGroupSuccessorAPI.schemaConfirm).mockImplementation(f.api.schemaConfirm);
  vi.mocked(studioV2WriteGroupSuccessorAPI.activationReview).mockImplementation(f.api.activationReview);
  vi.mocked(studioV2WriteGroupSuccessorAPI.activationReadiness).mockImplementation(f.api.activationReadiness);
  vi.mocked(studioV2WriteGroupSuccessorAPI.confirmCutover).mockImplementation(f.api.confirmCutover);
  vi.mocked(studioV2WriteGroupSuccessorAPI.confirmFirstStart).mockImplementation(f.api.confirmFirstStart);
  vi.mocked(studioV2WriteGroupSuccessorAPI.operationStatus).mockImplementation(f.api.operationStatus);
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });
it('opens the original lifecycle dialog from exact saved scope, with one editor and on-demand destination metadata only', async () => {
  vi.mocked(dbTargetAPI.listConnectors).mockResolvedValue([{ id: 'C2', identity_revision: 'C2R', name: 'Other saved DB', kind: 'sqlite',
    connection_config: { password: 'PRIVATE_CREDENTIAL' }, status: 'ready', last_check_error: 'PRIVATE_DIAGNOSTIC', enabled: true, default_write_interval_seconds: 15, created_at: '', updated_at: '' }]);
  await open(); const input = screen.getByLabelText('Retained editor'); expect(dbTargetAPI.listConnectors).not.toHaveBeenCalled();
  fireEvent.click(screen.getByTestId('saved-recording-transition'));
  expect(await screen.findByRole('dialog')).toHaveTextContent('dialog.title');
  await waitFor(() => expect(dbTargetAPI.listConnectors).toHaveBeenCalledOnce());
  expect(studioV2WriteGroupLifecycleAPI.lifecycle).toHaveBeenLastCalledWith('G1');
  expect(await screen.findByRole('option', { name: 'Other saved DB' })).toBeInTheDocument();
  expect(document.body.textContent).not.toMatch(/PRIVATE_CREDENTIAL|PRIVATE_DIAGNOSTIC/);
  expect(screen.getAllByTestId('canonical-editor')).toHaveLength(1); expect(screen.getByLabelText('Retained editor')).toBe(input);
  expect(studioV2WorkspaceWriteGroupsAPI.list).not.toHaveBeenCalled(); noCommands();
});
it('rejects stale and foreign current availability before review or mutations', async () => {
  for (const kind of ['stale', 'foreign']) {
    const view = await open(); vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockResolvedValueOnce({ ...f.lifecycle, ...(kind === 'stale' ? { group_revision: 'NEW' } : { workspace_id: 'OTHER' }) });
    fireEvent.click(screen.getByTestId('saved-recording-transition')); expect(await screen.findByRole('alert')).toHaveTextContent('dialog.errors.stale');
    expect(screen.queryByLabelText('dialog.name')).not.toBeInTheDocument(); noCommands(); view.unmount();
  }
});
it('revokes a captured lifecycle action when the retained locator becomes suspended', async () => {
  const view = await open(); const host = captured.editor?.lifecycleHost; expect(host).toBeDefined();
  view.rerender(view.content(view.route, true)); act(() => host?.onAction({ ...f.intent, action: 'transition' }));
  expect(screen.queryByRole('dialog')).not.toBeInTheDocument(); noCommands();
});
it('does not revive a late action after an A to B to A locator roundtrip', async () => {
  const view = await open(); let finish!: (value: typeof f.lifecycle) => void;
  vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockImplementationOnce(() => new Promise(resolve => { finish = resolve; }));
  fireEvent.click(screen.getByTestId('saved-recording-transition'));
  const other = structuredClone(view.route); other.locator.group_id = 'OTHER';
  view.rerender(view.content(other)); view.rerender(view.content(view.route));
  await act(async () => { finish(f.lifecycle); });
  expect(screen.queryByLabelText('dialog.name')).not.toBeInTheDocument(); noCommands();
});
it('does not show or look up a different groups retained operation', async () => {
  const foreign = fixture(f.workspace); foreign.envelope.predecessor_group_id = 'OTHER';
  expect(saveRecordingChange({ version: 1, envelope: foreign.envelope, reviewDigest: foreign.review.review_digest,
    submitted: { action: 'reserve', request: { envelope: foreign.envelope, review_digest: foreign.review.review_digest } }, phase: 'impact', outcome: 'unknown' })).toBe(true);
  mount(); await screen.findByTestId('canonical-editor');
  expect(captured.editor?.lifecycleHost).toBeDefined(); expect(screen.queryByText('dialog.resume')).not.toBeInTheDocument();
  expect(screen.getByTestId('saved-recording-transition')).toBeDisabled();
  expect(studioV2WriteGroupSuccessorAPI.operationStatus).not.toHaveBeenCalled(); noCommands();
});
it('admits the fresh successor after a verified scheduled receipt without resuming its predecessors operation', async () => {
  const request = { envelope: f.activation.review.envelope, review_digest: f.activation.review.review_digest, readiness_digest: f.activation.readiness_digest };
  const result = await f.api.confirmFirstStart('G1', request);
  expect(saveRecordingChange({ version: 1, envelope: f.envelope, reviewDigest: f.review.review_digest,
    submitted: { action: 'confirmFirstStart', request }, preparation: f.preparation, phase: 'scheduled', outcome: 'recorded', result })).toBe(true);
  f.group = { ...f.preparation.group, status: 'ready', applied_revision: 'R2' };
  f.lifecycle = { ...f.lifecycle, group_id: 'G2', group_revision: 'R2', applied_revision: 'R2', history: { ...f.lifecycle.history, current_group_id: 'G2' } };
  vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockResolvedValue(f.lifecycle);
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.delivery).mockRejectedValue(new Error('Delivery evidence unavailable'));
  mount(); await screen.findByTestId('canonical-editor');
  await waitFor(() => expect(captured.editor?.readonly).toBe(false));
  expect(screen.queryByText('dialog.resume')).not.toBeInTheDocument();expect(screen.queryByTestId('recording-change-recovery')).not.toBeInTheDocument();
  expect(studioV2WriteGroupSuccessorAPI.operationStatus).not.toHaveBeenCalled();noCommands();
  await waitFor(() => expect(screen.getByTestId('saved-recording-transition')).toBeEnabled());
  fireEvent.click(screen.getByTestId('saved-recording-transition'));
  expect(await screen.findByLabelText('dialog.name')).toHaveValue(f.seed.name);
  expect(studioV2WriteGroupLifecycleAPI.lifecycle).toHaveBeenLastCalledWith('G2');
  expect(vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mock.calls.every(([id]) => id === 'G2')).toBe(true);
  expect(studioV2WriteGroupSuccessorAPI.operationStatus).not.toHaveBeenCalled();noCommands();
});
it('does not accept a scheduled string without a verified receipt as completed foreign recovery', async () => {
  const foreign = fixture(f.workspace); foreign.envelope.predecessor_group_id = 'OTHER';
  const invalid = { version: 1 as const, envelope: foreign.envelope, reviewDigest: foreign.review.review_digest,
    submitted: { action: 'reserve' as const, request: { envelope: foreign.envelope, review_digest: foreign.review.review_digest } }, phase: 'scheduled' as const, outcome: 'recorded' as const };
  expect(saveRecordingChange(invalid)).toBe(false);sessionStorage.setItem(recordingChangeKey(invalid), JSON.stringify(invalid));
  mount();await screen.findByTestId('canonical-editor');
  expect(captured.editor?.readonly).toBe(true);expect(screen.getByTestId('saved-recording-transition')).toBeDisabled();
  expect(screen.queryByText('dialog.resume')).not.toBeInTheDocument();expect(studioV2WriteGroupSuccessorAPI.operationStatus).not.toHaveBeenCalled();noCommands();
});
it('reopens matching unknown recovery and reads only the original operation without replay', async () => {
  expect(saveRecordingChange({ version: 1, envelope: f.envelope, reviewDigest: f.review.review_digest,
    submitted: { action: 'reserve', request: { envelope: f.envelope, review_digest: f.review.review_digest } }, phase: 'impact', outcome: 'unknown' })).toBe(true);
  mount(); await screen.findByTestId('canonical-editor'); fireEvent.click(await screen.findByText('dialog.resume'));
  fireEvent.click(await screen.findByTestId('recording-change-check_status'));
  await waitFor(() => expect(studioV2WriteGroupSuccessorAPI.operationStatus).toHaveBeenCalledWith('G1', { workspace_id: f.workspace, lineage_id: 'L', operation_id: 'ROOT' }));
  noCommands(); expect(studioV2WorkspaceWriteGroupsAPI.list).not.toHaveBeenCalled();
});

it('does not send schema confirmation after scope is revoked while its stage read is awaiting', async () => {
  const view = await open(); fireEvent.click(screen.getByTestId('saved-recording-transition'));
  await waitFor(() => expect(screen.getByTestId('recording-change-review')).toBeEnabled());
  fireEvent.click(screen.getByTestId('recording-change-review')); fireEvent.click(await screen.findByTestId('recording-change-reserve_preview'));
  const confirm = await screen.findByTestId('recording-change-confirm_schema'); await waitFor(() => expect(confirm).toBeEnabled());
  let finish!: (value: typeof f.lifecycle) => void;
  const next = await f.lifecycleRead('G1');
  vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockImplementationOnce(() => new Promise(resolve => { finish = resolve; }));
  fireEvent.click(confirm); await waitFor(() => expect(finish).toBeDefined());
  view.rerender(view.content(view.route, true));
  await act(async () => { finish(next); });
  expect(studioV2WriteGroupSuccessorAPI.schemaConfirm).not.toHaveBeenCalled();
  expect(studioV2WriteGroupSuccessorAPI.reserve).toHaveBeenCalledOnce();
});
