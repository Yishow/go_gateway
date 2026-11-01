import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import type { GroupEditorProps } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import type { SavedRecordingStartReviewProps } from '@/features/datalink/workbench-v2/recording-editor/SavedRecordingStartReview';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { studioV2LocatorHydration } from '@/pages/datalink/workbench-v2/studioV2LocatorHydration';
import { readStudioLocatorEditorEntry } from '@/services/studioV2LocatorEntry';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { modbusShareAPI } from '@/services/datalink';
import { locatorRead, nativeDevice, nativeRule, nativeMapping } from './studio-locator-fixtures';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2LocatorEntry', () => ({ readStudioLocatorEditorEntry: vi.fn() }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn() } }));
vi.mock('@/services/datalink', async original => ({ ...await original<typeof import('@/services/datalink')>(), modbusShareAPI: { status: vi.fn() } }));
let editor: GroupEditorProps, review: SavedRecordingStartReviewProps;
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor', () => ({ GroupEditor: (props: GroupEditorProps) => { editor = props; return <button onClick={() => props.group && props.onReviewScopedStart?.(props.group)}>Review saved Start</button>; } }));
vi.mock('@/features/datalink/workbench-v2/recording-editor/SavedRecordingStartReview', () => ({ SavedRecordingStartReview: (props: SavedRecordingStartReviewProps) => { review = props; return <div data-testid="locator-start-review">{props.group.id}/{props.group.revision}<button onClick={props.onClose}>Close review</button></div>; } }));
function entry(w = 'W2', g = 'G3') {
 const read = locatorRead(4); const group = read.entry.write_groups!.groups[0];
 group.revision = g; group.destination.storage_strategy = 'managed'; group.recording_profile = { kind: 'mixed', result_shape: 'member-calculation-v1' }; group.members = [{ ...group.members[0], device_id: 'D', point_id: 'P-native', tag_id: 'T' }];
 read.entry.workspace_revision = w; read.entry.write_groups!.workspace_revision = w; read.entry.group_heads.G = g; read.entry.rule_heads.R = 'R2'; read.entry.revalidated_mapping_ids = ['M'];
 if (read.selection.state === 'selected') Object.assign(read.selection.target, { workspaceRevision: w, groupRevision: g });
 Object.assign(read.entry, { editor_hydrated: true, devices: [nativeDevice()], source_rules: [nativeRule()], mappings: [nativeMapping()], database_targets: [],
 hydration_scope: { complete_workspace: false, device_ids: ['D'], rule_ids: ['R'], mapping_ids: ['M'], group_ids: ['G'] },
 database_config: { id: group.destination.connector_id, identity_revision: group.destination.connector_revision, workspace_id: 'W', setup_revision: w, kind: 'sqlite', name: 'Saved', host: '', port: 0, database: '/d.db', username: '', schema: '', table: group.destination.table_name, write_mode: 'insert', write_interval_seconds: 15, timestamp_column: 'at', status: 'ready', created_at: 'now', updated_at: 'now' } });
 return read;
}
let current = entry(); const clients: QueryClient[] = [];
function mount() {
 const route = locatorRead(4); route.entry.write_groups!.groups[0].status = 'draft'; const client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); clients.push(client);
 const content = (read = route, suspended = false) => <QueryClientProvider client={client}><Step4Database state={studioV2LocatorHydration(route)!} workspaceId="W" dispatch={vi.fn()} locatorEntry={read} locatorOnly locatorSuspended={suspended} /></QueryClientProvider>;
 return { ...render(content()), content, route };
}
async function open() { const view = mount(); fireEvent.click(screen.getByTestId('group-open-G')); await screen.findByText('Review saved Start'); return view; }
beforeEach(() => { vi.clearAllMocks(); current = entry(); vi.mocked(readStudioLocatorEditorEntry).mockImplementation(async () => ({ ok: true, value: current }));
 vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockImplementation(async () => current.entry.write_groups!);
 vi.mocked(modbusShareAPI.status).mockResolvedValue({ workspace_id: 'W', workspace_revision: 'RUNTIME-9', settings_revision: 'SETTINGS-7', readiness_token: 'proof', readiness: true, enabled: false, port: 502, address: '', bind_state: 'disabled', mapping_count: 0 }); });
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });
it('uses exact hydrated locator editor scope and independent runtime proof without loading the complete workspace', async () => {
 await open(); expect(modbusShareAPI.status).not.toHaveBeenCalled(); fireEvent.click(screen.getByText('Review saved Start'));
 await screen.findByTestId('locator-start-review'); expect(review.readonly).toBe(false); expect(review.workspaceRevision).toBe('W2');
 expect(review.shareStatus?.workspace_revision).toBe('RUNTIME-9'); expect(review.state.devices.map(device => device.id)).toEqual(['D']);
 expect(studioV2WorkspaceWriteGroupsAPI.list).not.toHaveBeenCalled();
 await act(async () => { expect(await review.fetchLatest()).toEqual(current.entry.write_groups); });
 expect(readStudioLocatorEditorEntry).toHaveBeenLastCalledWith(current.locator, { workspaceRevision: 'W2', groupRevision: 'G3' }, expect.any(AbortSignal));
});
it.each(['foreign', 'stale', 'incomplete'] as const)('does not grant %s scoped read to the original Start', async kind => {
 await open(); fireEvent.click(screen.getByText('Review saved Start')); await screen.findByTestId('locator-start-review');
 current = entry(); if (kind === 'foreign') current.workspace.id = 'OTHER'; if (kind === 'stale') current.entry.group_heads.G = 'NEW'; if (kind === 'incomplete') current.entry.mappings = [];
 await act(async () => { expect(await review.fetchLatest()).toBeUndefined(); });
});
it.each(['suspended', 'route-switch', 'roundtrip'] as const)('invalidates captured Start reads after %s without borrowing a new scope', async kind => {
 const view = await open(); fireEvent.click(screen.getByText('Review saved Start')); await screen.findByTestId('locator-start-review'); const captured = review.fetchLatest;
 view.rerender(view.content(kind === 'suspended' ? view.route : locatorRead(4, 'OTHER'), kind === 'suspended'));
 if (kind === 'roundtrip') view.rerender(view.content(view.route));
 expect(review.readonly).toBe(true); await act(async () => { expect(await captured()).toBeUndefined(); });
});
it('uses current own-save readback rather than outer original URL revisions', async () => {
 await open(); current = entry('W3', 'G4'); await act(async () => { editor.onSaved(current.entry.write_groups!.groups[0], 'W3'); });
 await waitFor(() => expect(editor.group?.revision).toBe('G4')); fireEvent.click(screen.getByText('Review saved Start')); await screen.findByTestId('locator-start-review');
 expect(review.workspaceRevision).toBe('W3'); expect(review.group.revision).toBe('G4'); expect(review.readonly).toBe(false);
});

it('refreshes the original scoped editor after completion without granting the old review a new write lease', async () => {
 await open(); fireEvent.click(screen.getByText('Review saved Start')); await screen.findByTestId('locator-start-review');
 current = entry('W3', 'G3'); current.entry.write_groups!.groups[0].applied_revision = 'APPLIED';
 await act(async () => { review.onCompleted?.(); });
 await waitFor(() => expect(editor.workspaceRevision).toBe('W3'));
 expect(review.group.applied_revision).toBe('APPLIED'); expect(review.readonly).toBe(true);
});

it.each(['foreign', 'not-ready', 'missing-token', 'error'] as const)('denies fresh Start when original runtime barrier is %s', async kind => {
 if (kind === 'error') vi.mocked(modbusShareAPI.status).mockRejectedValueOnce(new Error('offline'));
 else vi.mocked(modbusShareAPI.status).mockResolvedValueOnce({ workspace_id: kind === 'foreign' ? 'OTHER' : 'W', workspace_revision: 'RUNTIME', settings_revision: 'SETTINGS', readiness_token: kind === 'missing-token' ? '' : 'token', readiness: kind !== 'not-ready', enabled: false, port: 502, address: '', bind_state: 'disabled', mapping_count: 0 });
 await open(); fireEvent.click(screen.getByText('Review saved Start'));
 await waitFor(() => expect(screen.getByRole('alert')).toHaveTextContent('start.changed'));
 expect(screen.queryByTestId('locator-start-review')).not.toBeInTheDocument();
 expect(studioV2WorkspaceWriteGroupsAPI.list).not.toHaveBeenCalled();
});
it('denies late scoped read after a route change and does not replace the retained editor', async () => {
 const view = await open(); fireEvent.click(screen.getByText('Review saved Start')); await screen.findByTestId('locator-start-review');
 let finish!: (value: Awaited<ReturnType<typeof readStudioLocatorEditorEntry>>) => void;
 vi.mocked(readStudioLocatorEditorEntry).mockImplementationOnce(() => new Promise(resolve => { finish = resolve; }));
 const pending = review.fetchLatest(); view.rerender(view.content(locatorRead(4, 'OTHER')));
 await act(async () => { finish({ ok: true, value: current }); expect(await pending).toBeUndefined(); });
 expect(editor.group?.id).toBe('G'); expect(review.readonly).toBe(true);
});
it('does not borrow changed source definitions from a same-head read', async () => {
 await open(); fireEvent.click(screen.getByText('Review saved Start')); await screen.findByTestId('locator-start-review');
 current = entry(); current.entry.source_rules![0].scale_multiplier = 100;
 await act(async () => { expect(await review.fetchLatest()).toBeUndefined(); });
});

it('closing a successfully refreshed scoped editor keeps its verified status on the same-route list, but Open still rechecks canonical heads', async () => {
 await open(); fireEvent.click(screen.getByText('Review saved Start')); await screen.findByTestId('locator-start-review');
 current = entry('W3', 'G3'); current.entry.write_groups!.groups[0].status = 'ready';
 await act(async () => { review.onCompleted?.(); }); await waitFor(() => expect(editor.workspaceRevision).toBe('W3'));
 fireEvent.click(screen.getByText('Close review')); act(() => editor.onClose());
 expect(screen.getByTestId('group-list-item-G')).toHaveTextContent('step4.group.status.ready');
 expect(screen.getByTestId('group-list-item-G')).not.toHaveTextContent('step4.group.status.draft');
 const calls = vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mock.calls.length;
 fireEvent.click(screen.getByTestId('group-open-G')); await screen.findByText('Review saved Start');
 expect(vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mock.calls.length).toBeGreaterThan(calls);
 expect(readStudioLocatorEditorEntry).toHaveBeenLastCalledWith(current.locator, { workspaceRevision: 'W3', groupRevision: 'G3' }, expect.any(AbortSignal));
});
it.each(['loader', 'suspended', 'roundtrip'] as const)('does not use the display receipt after %s owner invalidation', async kind => {
 const view = await open(); act(() => editor.onClose());
 expect(screen.getByTestId('group-list-item-G')).toHaveTextContent('step4.group.status.ready');
 view.rerender(view.content(kind === 'suspended' ? view.route : structuredClone(view.route), kind === 'suspended'));
 if (kind === 'roundtrip') view.rerender(view.content(view.route));
 expect(screen.getByTestId('group-list-item-G')).toHaveTextContent('step4.group.status.draft');
});
