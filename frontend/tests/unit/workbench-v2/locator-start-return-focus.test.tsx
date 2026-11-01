import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import type { GroupEditorProps } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { studioV2LocatorHydration } from '@/pages/datalink/workbench-v2/studioV2LocatorHydration';
import { readStudioLocatorEditorEntry } from '@/services/studioV2LocatorEntry';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2RecordingStartAPI } from '@/services/studioV2RecordingStart';
import { modbusShareAPI } from '@/services/datalink';
import { saveBasicRecordingIntent } from '@/features/datalink/workbench-v2/state/basicRecording';
import type { RecordingStartOperation } from '@/types/studioV2RecordingStart';
import type { StudioLocatorEntryRead } from '@/types/studioV2LocatorEntry';
import { locatorRead, nativeDevice, nativeRule, nativeMapping } from './studio-locator-fixtures';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2LocatorEntry', () => ({ readStudioLocatorEditorEntry: vi.fn() }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn(), readiness: vi.fn(), delivery: vi.fn().mockResolvedValue(null) } }));
vi.mock('@/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), operation: vi.fn(), lookup: vi.fn(), ensureBasic: vi.fn() } }));
vi.mock('@/services/datalink', async original => ({ ...await original<typeof import('@/services/datalink')>(), modbusShareAPI: { status: vi.fn() } }));
// Only editor fields are replaced. WGS, locator host, Start review, completion/readback and Radix are real.
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor', () => ({ GroupEditor: (props: GroupEditorProps) => <div data-testid="group-editor"><input aria-label="record name" defaultValue={props.group?.name} />{props.repairContent}</div> }));
function entry(w = 'W2', g = 'G3'): StudioLocatorEntryRead {
 const read = locatorRead(4); const group = read.entry.write_groups!.groups[0];
 Object.assign(group, { revision: g, applied_revision: g, status: 'ready', recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' } });
 group.destination.storage_strategy = 'managed'; group.members = [{ ...group.members[0], device_id: 'D', point_id: 'P-native', tag_id: 'T' }];
 read.locator.group_revision = g; read.entry.workspace_revision = w; read.entry.write_groups!.workspace_revision = w;
 read.entry.group_heads.G = g; read.entry.rule_heads.R = 'R2'; read.entry.revalidated_mapping_ids = ['M'];
 if (read.selection.state === 'selected') Object.assign(read.selection.target, { workspaceRevision: w, groupRevision: g });
 return { ...read, entry: { ...read.entry, editor_hydrated: true, devices: [nativeDevice()], source_rules: [nativeRule()], mappings: [nativeMapping()], database_targets: [],
  hydration_scope: { complete_workspace: false, device_ids: ['D'], rule_ids: ['R'], mapping_ids: ['M'], group_ids: ['G'] },
  database_config: { id: group.destination.connector_id, identity_revision: group.destination.connector_revision, workspace_id: 'W', setup_revision: w,
   kind: 'sqlite', name: 'Saved', host: '', port: 0, database: '/d.db', username: '', schema: '', table: group.destination.table_name,
   write_mode: 'insert', write_interval_seconds: 15, timestamp_column: 'at', status: 'ready', created_at: 'now', updated_at: 'now' } } };
}
const api = vi.mocked(studioV2RecordingStartAPI), clients: QueryClient[] = [];
const operation: RecordingStartOperation = { operation_id: 'O1', action: 'recording_start', status: 'succeeded', stage: 'complete', intent_digest: 'a'.repeat(64),
 workspace_id: 'W', setup_revision: 'W2', device_ids: ['D'], devices: [{ device_id: 'D', activated: true }],
 groups: [{ group_id: 'G', group_revision: 'G3', applied_revision: 'G3', saved: true, ready: true, applied: true }], created_at: '2026-11-01T00:00:00Z', updated_at: '2026-11-01T00:00:00Z' };
let current: StudioLocatorEntryRead;
beforeEach(() => {
 vi.clearAllMocks(); sessionStorage.clear(); localStorage.clear(); current = entry();
 vi.mocked(readStudioLocatorEditorEntry).mockImplementation(async () => ({ ok: true, value: structuredClone(current) }));
 vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockImplementation(async () => current.entry.write_groups!);
 vi.mocked(modbusShareAPI.status).mockResolvedValue({ workspace_id: 'W', workspace_revision: 'SHARE', settings_revision: 'SETTINGS', readiness_token: 'proof', readiness: true, enabled: false, port: 502, address: '', bind_state: 'disabled', mapping_count: 0 });
 api.operation.mockResolvedValue(operation);
 saveBasicRecordingIntent({ scope_key: 'W/saved-record/G', operation_id: 'O1', owner: { kind: 'saved_record', groupId: 'G', deviceIds: ['D'] },
  request: { request_id: 'original-request', workspace_id: 'W', expected_workspace_revision: 'W2', device_ids: ['D'], groups: [{ group_id: 'G', expected_group_revision: 'G3', expected_connector_revision: current.entry.write_groups!.groups[0].destination.connector_revision }] } });
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); sessionStorage.clear(); localStorage.clear(); });
async function open() {
 const route = entry(); const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } }); clients.push(client);
 const content = (suspended = false) => <QueryClientProvider client={client}><Step4Database state={studioV2LocatorHydration(route)!} workspaceId="W" dispatch={vi.fn()}
  presentation="simple" initialSelection="located" locatorEntry={route} locatorOnly locatorSuspended={suspended}
  repairIntent={{ workspaceId: 'W', workspaceRevision: 'W2', groupId: 'G', groupRevision: 'G3', section: 'operation', operation: { kind: 'query_operation', owner: 'start', deviceId: 'D', operationId: 'O1' } }} /></QueryClientProvider>;
 const view = render(content()); const trigger = await screen.findByTestId('recording-repair-start'); trigger.focus(); fireEvent.click(trigger);
 await screen.findByRole('dialog'); await waitFor(() => expect(api.operation).toHaveBeenCalledWith('O1'));
 return { ...view, content, trigger };
}
it.each(['Escape', 'Close'])('returns from the real retained Start completion readback on %s without renewing its command lease', async dismissal => {
 let finish!: () => void;
 vi.mocked(readStudioLocatorEditorEntry).mockImplementationOnce(async () => ({ ok: true, value: entry() }))
  .mockImplementationOnce(() => new Promise(resolve => { finish = () => resolve({ ok: true, value: structuredClone(current) }); }));
 const view = await open();
 await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledTimes(2));
 await waitFor(() => expect(view.trigger.isConnected).toBe(false));
 await act(async () => finish());
 await waitFor(() => expect(screen.getByTestId('recording-repair-start')).toBeInTheDocument());
 const currentTrigger = screen.getByTestId('recording-repair-start'); expect(currentTrigger).not.toBe(view.trigger);
 expect(screen.queryByTestId('saved-recording-start-confirm')).toBeNull();
 if (dismissal === 'Escape') fireEvent.keyDown(document.activeElement ?? document, { key: 'Escape' });
 else fireEvent.click(screen.getByRole('button', { name: 'start.close' }));
 await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull()); await waitFor(() => expect(currentTrigger).toHaveFocus());
 expect(api.start).not.toHaveBeenCalled(); expect(api.ensureBasic).not.toHaveBeenCalled();
});
it.each(['revision', 'foreign', 'ABA'] as const)('does not borrow a replacement focus target after %s invalidates the real completion scope', async kind => {
 const read = vi.mocked(readStudioLocatorEditorEntry); let finish!: () => void;
 read.mockImplementationOnce(async () => ({ ok: true, value: entry() })).mockImplementationOnce(() => new Promise(resolve => { finish = () => resolve({ ok: true, value: current }); }));
 if (kind === 'revision') current = entry('W3', 'G4');
 const view = await open(); await waitFor(() => expect(read).toHaveBeenCalledTimes(2));
 if (kind === 'foreign') current = { ...entry(), workspace: { ...entry().workspace, id: 'OTHER' } };
 if (kind === 'ABA') { view.rerender(view.content(true)); view.rerender(view.content(false)); }
 await act(async () => finish());
 fireEvent.click(screen.getByRole('button', { name: 'start.close' }));
 await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull());
 expect(view.trigger).not.toHaveFocus(); const replacement = screen.queryByTestId('recording-repair-start'); if (replacement) expect(replacement).not.toHaveFocus();
 expect(api.start).not.toHaveBeenCalled(); expect(api.ensureBasic).not.toHaveBeenCalled();
});
it('does not wait for a pending completion readback or steal later user focus', async () => {
 let finish!: () => void;
 vi.mocked(readStudioLocatorEditorEntry).mockImplementationOnce(async () => ({ ok: true, value: entry() }))
  .mockImplementationOnce(() => new Promise(resolve => { finish = () => resolve({ ok: true, value: structuredClone(current) }); }));
 await open(); await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledTimes(2));
 fireEvent.click(screen.getByRole('button', { name: 'start.close' })); await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull());
 const next = screen.getByRole('textbox', { name: 'record name' }); next.focus();
 await act(async () => finish()); await screen.findByTestId('recording-repair-start');
 expect(next).toHaveFocus(); expect(api.start).not.toHaveBeenCalled(); expect(api.ensureBasic).not.toHaveBeenCalled();
});
