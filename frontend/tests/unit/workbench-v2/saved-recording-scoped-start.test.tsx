import { recordingTestWriteRecovery } from '../../../src/features/datalink/workbench-v2/state/recordingTestWriteRecovery';
import { safeBasicNextAction, safeBasicOperationMessage } from '../../../src/features/datalink/workbench-v2/state/basicRecordingMessages';
import type { ComponentProps } from 'react';
import { basicRecordingIntentKey, saveBasicRecordingIntent } from '../../../src/features/datalink/workbench-v2/state/basicRecording';
import type { WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { SavedRecordingStartReview } from '../../../src/features/datalink/workbench-v2/recording-editor/SavedRecordingStartReview';
import { batchPanelState, longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { studioV2RecordingStartAPI } from '../../../src/services/studioV2RecordingStart';
import { studioV2WorkspaceWriteGroupsAPI } from '../../../src/services/studioV2WorkspaceWriteGroups';
import type { RecordingStartOperation } from '../../../src/types/studioV2RecordingStart';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('../../../src/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), operation: vi.fn(), lookup: vi.fn(), ensureBasic: vi.fn() } }));
vi.mock('../../../src/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { readiness: vi.fn() } }));
const base = batchPanelState();
const state: WorkbenchV2State = { ...base, devices: [...base.devices, { ...base.devices[0], id: 'dev-2', name: 'PLC B' }, { ...base.devices[0], id: 'unselected-device', name: 'Do not start' }],
 rules: [{ ...base.rules[0], save_state: 'saved' as const }, { ...base.rules[0], id: 'rule-2', device_id: 'dev-2', save_state: 'saved' as const }],
 points: base.points.map((point, i) => i === 1 ? { ...point, device_id: 'dev-2', rule_id: 'rule-2' } : point),
 mappings: Object.fromEntries(Object.entries(base.mappings).map(([key, value], i) => [key, i === 1 ? { ...value, device_id: 'dev-2' } : value])) }; 
const saved = longProfileGroup({ id: 'independent-mixed', circuitId: 'c1', circuitName: 'Mixed', interval: 15, points: [0, 1] });
const group = { ...saved, members: saved.members.map((member, i) => i === 1 ? { ...member, device_id: 'dev-2' } : member), recording_profile: { kind: 'mixed' as const, result_shape: 'member-calculation-v1' as const } };
const api = vi.mocked(studioV2RecordingStartAPI);
const operation: RecordingStartOperation = { operation_id: 'start-only-this-record', action: 'recording_start', workspace_id: 'ws-1', setup_revision: 'w2', intent_digest: 'a'.repeat(64),
 status: 'succeeded', stage: 'complete', device_ids: ['dev-1', 'dev-2'], devices: [{ device_id: 'dev-1', activated: true }, { device_id: 'dev-2', activated: true }],
 groups: [{ group_id: group.id, group_revision: group.revision, applied_revision: group.revision, saved: true, ready: true, applied: true }], created_at: '2026-10-31T00:00:00Z', updated_at: '2026-10-31T00:00:01Z' };
const latest = { workspace_id: 'ws-1', workspace_revision: 'w1', groups: [group] };
function mount(fetchLatest = vi.fn(async () => latest), overrides: Partial<ComponentProps<typeof SavedRecordingStartReview>> = {}) {
 const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
 let currentState = state;
 const tree = () => <QueryClientProvider client={client}><SavedRecordingStartReview group={group} workspaceId="ws-1" workspaceRevision="w1" state={currentState}
  shareStatus={{ workspace_id: 'ws-1', enabled: false, port: 502, address: '', bind_state: 'disabled', mapping_count: 0, readiness_token: 'share-token', settings_revision: 'settings-1', workspace_revision: 'share-1' }} readonly={false} fetchLatest={fetchLatest} onClose={vi.fn()} onOperationActiveChange={vi.fn()} {...overrides} /></QueryClientProvider>;
 const view = render(tree());
 return { ...view, changeSource: () => { currentState = { ...state, rules: state.rules.map(rule => ({ ...rule, save_state: 'saving' as const })) }; view.rerender(tree()); } };
}
beforeEach(() => { vi.restoreAllMocks(); vi.clearAllMocks(); sessionStorage.clear(); api.start.mockResolvedValue(operation); api.operation.mockResolvedValue(operation);
 vi.mocked(studioV2WorkspaceWriteGroupsAPI.readiness).mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'w1', group_id: group.id, group_revision: group.revision, applied_revision: '', config_ready: true, schema_ready: true, ready: true, issues: [] }); });
it('starts only the saved independent record and its two selected sources once, never creates Basic, and waits for SQL evidence', async () => {
 mount(); const button = screen.getByTestId('saved-recording-start-confirm');
 fireEvent.click(button); fireEvent.click(button);
 await waitFor(() => expect(api.start).toHaveBeenCalledTimes(1));
 expect(api.start).toHaveBeenCalledWith(expect.objectContaining({ workspace_id: 'ws-1', expected_workspace_revision: 'w1', workspace_revision: 'share-1', settings_revision: 'settings-1', readiness_token: 'share-token', device_ids: ['dev-1', 'dev-2'], groups: [{ group_id: group.id, expected_group_revision: group.revision, expected_connector_revision: group.destination.connector_revision }] }));
 expect(api.ensureBasic).not.toHaveBeenCalled();
 expect(await screen.findByTestId('saved-recording-start-waiting')).toBeInTheDocument();
});
it('a source edit while fresh scope is loading prevents the unsent Start', async () => {
 let resolve!: (value: typeof latest) => void;
 const view = mount(vi.fn(() => new Promise(done => { resolve = done; })));
 fireEvent.click(screen.getByTestId('saved-recording-start-confirm')); view.changeSource();
 await act(async () => { resolve(latest); });
 expect(api.start).not.toHaveBeenCalled(); expect(api.ensureBasic).not.toHaveBeenCalled();
});
it('unknown acknowledgement exposes only original operation lookup without a replacement Start', async () => {
 api.start.mockResolvedValue({ ...operation, status: 'unknown', stage: 'activation' });
 api.operation.mockResolvedValue({ ...operation, status: 'unknown', stage: 'activation' });
 mount(); fireEvent.click(screen.getByTestId('saved-recording-start-confirm'));
 const check = await screen.findByTestId('saved-recording-start-check');
 await waitFor(() => expect(check).toBeEnabled()); fireEvent.click(check);
 await waitFor(() => expect(api.operation).toHaveBeenCalledWith(operation.operation_id));
 expect(api.start).toHaveBeenCalledTimes(1); expect(screen.queryByTestId('saved-recording-start-waiting')).not.toBeInTheDocument();
});
it('lost acknowledgement can explicitly look up its exact retained request without posting Start again', async () => {
 api.start.mockRejectedValueOnce(new TypeError('Failed to fetch'));
 api.lookup.mockResolvedValue(operation);
 mount(); fireEvent.click(screen.getByTestId('saved-recording-start-confirm'));
 await waitFor(() => expect(api.start).toHaveBeenCalledTimes(1));
 const lookup = await screen.findByTestId('saved-recording-start-lookup');
 await waitFor(() => expect(lookup).toBeEnabled()); fireEvent.click(lookup); fireEvent.click(lookup);
 await waitFor(() => expect(api.lookup).toHaveBeenCalledTimes(1));
 expect(api.lookup).toHaveBeenCalledWith(api.start.mock.calls[0][0]);
 expect(await screen.findByTestId('saved-recording-start-waiting')).toBeInTheDocument();
 expect(api.start).toHaveBeenCalledTimes(1);
});
it('a missing original request lookup remains unknown and never submits a replacement', async () => {
 api.start.mockRejectedValueOnce(new TypeError('Failed to fetch'));
 api.lookup.mockRejectedValueOnce({ code: 'RECORDING_START_NOT_FOUND', status: 404 });
 mount(); fireEvent.click(screen.getByTestId('saved-recording-start-confirm'));
 const lookup = await screen.findByTestId('saved-recording-start-lookup');
 await waitFor(() => expect(lookup).toBeEnabled()); fireEvent.click(lookup);
 await waitFor(() => expect(api.lookup).toHaveBeenCalledTimes(1));
 expect(screen.queryByTestId('saved-recording-start-waiting')).not.toBeInTheDocument();
 expect(api.start).toHaveBeenCalledTimes(1); expect(screen.queryByTestId('saved-recording-start-confirm')).not.toBeInTheDocument();
});

it.each(['revision', 'membership'] as const)('a fresh review can explicitly start a new saved %s only after the original request is proven terminal', async change => {
 const next = { ...group, revision: 'new-record-revision', members: change === 'membership' ? [group.members[0]] : group.members };
 const selectedDevices = change === 'membership' ? ['dev-1'] : ['dev-1', 'dev-2'];
 api.start.mockResolvedValue({ ...operation, operation_id: 'new-start-op', device_ids: selectedDevices, devices: selectedDevices.map(device_id => ({ device_id, activated: true })), groups: [{ ...operation.groups[0], group_revision: next.revision, applied_revision: next.revision }] });
 saveBasicRecordingIntent({ scope_key: `ws-1/saved-record/${group.id}`, owner: { kind: 'saved_record', groupId: group.id, deviceIds: ['dev-1', 'dev-2'] }, operation_id: operation.operation_id,
  request: { request_id: 'original-old-request', workspace_id: 'ws-1', expected_workspace_revision: 'older-db-head', device_ids: ['dev-1', 'dev-2'], groups: [{ group_id: group.id, expected_group_revision: group.revision, expected_connector_revision: group.destination.connector_revision }] } });
 vi.mocked(studioV2WorkspaceWriteGroupsAPI.readiness).mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'w1', group_id: group.id, group_revision: next.revision, applied_revision: '', config_ready: true, schema_ready: true, ready: true, issues: [] });
 const onCompleted = vi.fn();
 mount(vi.fn(async () => ({ ...latest, groups: [next] })), { group: next, onCompleted });
 await waitFor(() => expect(api.operation).toHaveBeenCalledWith(operation.operation_id));
 const confirm = await screen.findByTestId('saved-recording-start-confirm');
 await waitFor(() => expect(confirm).toBeEnabled());
 expect(screen.queryByTestId('saved-recording-start-waiting')).not.toBeInTheDocument();
 expect(onCompleted, 'historical terminal scope must not refresh and revoke this fresh review').not.toHaveBeenCalled();
 fireEvent.click(confirm);
 await waitFor(() => expect(api.start).toHaveBeenCalledTimes(1));
 expect(api.start.mock.calls[0][0]).toMatchObject({ device_ids: selectedDevices, groups: [{ group_id: group.id, expected_group_revision: next.revision }] });
 expect(api.start.mock.calls[0][0].request_id).not.toBe('original-old-request');
 await waitFor(() => expect(onCompleted).toHaveBeenCalledTimes(1));
});
it.each(['unknown', 'partial', 'retry_same_request'] as const)('a changed saved revision cannot authorize a new request while the original operation is %s', async status => {
 const next = { ...group, revision: 'new-record-revision' };
 saveBasicRecordingIntent({ scope_key: `ws-1/saved-record/${group.id}`, owner: { kind: 'saved_record', groupId: group.id, deviceIds: ['dev-1', 'dev-2'] }, operation_id: operation.operation_id,
  request: { request_id: 'original-old-request', workspace_id: 'ws-1', expected_workspace_revision: 'older-db-head', device_ids: ['dev-1', 'dev-2'], groups: [{ group_id: group.id, expected_group_revision: group.revision, expected_connector_revision: group.destination.connector_revision }] } });
 api.operation.mockResolvedValue({ ...operation, status: status === 'retry_same_request' ? 'failed' : status, next_action: status === 'retry_same_request' ? 'retry_same_request' : undefined, stage: 'activation' });
 mount(vi.fn(async () => ({ ...latest, groups: [next] })), { group: next });
 await waitFor(() => expect(api.operation).toHaveBeenCalledWith(operation.operation_id));
 expect(screen.queryByTestId('saved-recording-start-confirm')).not.toBeInTheDocument();
 expect(api.start).not.toHaveBeenCalled();
});

it('read-only and offline reviews never dispatch Start', () => {
 const view = mount(undefined, { readonly: true });
 expect(screen.getByTestId('saved-recording-start-confirm')).toBeDisabled(); view.unmount();
 vi.spyOn(navigator, 'onLine', 'get').mockReturnValue(false);
 mount(); expect(screen.getByTestId('saved-recording-start-confirm')).toBeDisabled();
 expect(api.start).not.toHaveBeenCalled();
});
it('a failed original-request storage write blocks the unsent command rather than losing recovery', async () => {
 vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => { throw new Error('storage denied'); });
 mount(); fireEvent.click(screen.getByTestId('saved-recording-start-confirm'));
 await waitFor(() => expect(screen.getByRole('alert')).toBeInTheDocument());
 expect(api.start).not.toHaveBeenCalled();
});
it('a late original-request lookup cannot publish into another record review', async () => {
 let resolve!: (value: RecordingStartOperation) => void;
 api.start.mockRejectedValueOnce(new TypeError('Failed to fetch'));
 api.lookup.mockImplementation(() => new Promise(done => { resolve = done; }));
 const old = mount(); fireEvent.click(screen.getByTestId('saved-recording-start-confirm'));
 const lookup = await screen.findByTestId('saved-recording-start-lookup'); await waitFor(() => expect(lookup).toBeEnabled()); fireEvent.click(lookup);
 await waitFor(() => expect(api.lookup).toHaveBeenCalledTimes(1)); old.unmount();
 const other = { ...group, id: 'other-record' }; mount(vi.fn(async () => ({ ...latest, groups: [other] })), { group: other });
 await act(async () => { resolve(operation); });
 expect(screen.getByTestId('saved-recording-start-confirm')).toBeInTheDocument();
 expect(screen.queryByTestId('saved-recording-start-waiting')).not.toBeInTheDocument(); expect(api.start).toHaveBeenCalledTimes(1);
});
it('a proven terminal failure without same-request retry can be explicitly reviewed again with fresh proof', async () => {
 saveBasicRecordingIntent({ scope_key: `ws-1/saved-record/${group.id}`, owner: { kind: 'saved_record', groupId: group.id, deviceIds: ['dev-1', 'dev-2'] }, operation_id: operation.operation_id,
  request: { request_id: 'failed-old-request', workspace_id: 'ws-1', expected_workspace_revision: 'w1', device_ids: ['dev-1', 'dev-2'], groups: [{ group_id: group.id, expected_group_revision: group.revision, expected_connector_revision: group.destination.connector_revision }] } });
 api.operation.mockResolvedValue({ ...operation, status: 'failed', stage: 'readiness', next_action: 'revalidate' });
 mount(); await waitFor(() => expect(api.operation).toHaveBeenCalledWith(operation.operation_id));
 const confirm = await screen.findByTestId('saved-recording-start-confirm'); await waitFor(() => expect(confirm).toBeEnabled());
 fireEvent.click(confirm); await waitFor(() => expect(api.start).toHaveBeenCalledTimes(1));
 expect(api.start.mock.calls[0][0].request_id).not.toBe('failed-old-request');
});

it('a malformed original owner store is unavailable, never empty permission to replace a request', () => {
 sessionStorage.setItem(basicRecordingIntentKey(`ws-1/saved-record/${group.id}`), '{malformed');
 mount(); expect(screen.getByTestId('saved-recording-start-confirm')).toBeDisabled();
 expect(api.start).not.toHaveBeenCalled();
});
it('a competing retained original request appearing before dispatch cannot be overwritten', async () => {
 mount();
 const retained = { version: 1, scope_key: `ws-1/saved-record/${group.id}`, owner: { kind: 'saved_record', groupId: group.id, deviceIds: ['dev-1', 'dev-2'] },
  request: { request_id: 'competing-original', workspace_id: 'ws-1', expected_workspace_revision: 'w1', device_ids: ['dev-1', 'dev-2'], groups: [{ group_id: group.id, expected_group_revision: group.revision, expected_connector_revision: group.destination.connector_revision }] } };
 sessionStorage.setItem(basicRecordingIntentKey(retained.scope_key), JSON.stringify(retained));
 await act(async () => { fireEvent.click(screen.getByTestId('saved-recording-start-confirm')); });
 expect(api.start).not.toHaveBeenCalled();
 expect(JSON.parse(sessionStorage.getItem(basicRecordingIntentKey(retained.scope_key))!)).toEqual(retained);
});
it('an unreadable cleanup-owner store blocks fresh Start without requiring a test write to run', () => {
 vi.spyOn(recordingTestWriteRecovery, 'read').mockReturnValue({ references: [], unavailable: true, unreadable: true });
 mount(); expect(screen.getByTestId('saved-recording-start-confirm')).toBeDisabled();
 expect(api.start).not.toHaveBeenCalled();
});
it('does not borrow a foreign runtime barrier even when its token is populated', () => {
 mount(undefined, { shareStatus: { workspace_id: 'OTHER', workspace_revision: 'RUNTIME', settings_revision: 'SETTINGS', readiness_token: 'token', enabled: false, port: 502, address: '', bind_state: 'disabled', mapping_count: 0 } });
 expect(screen.getByTestId('saved-recording-start-confirm')).toBeDisabled(); expect(api.start).not.toHaveBeenCalled();
});

const failedReadiness: RecordingStartOperation = { ...operation, status: 'failed', stage: 'readiness', reason: 'invalid_scope', next_action: 'review_selection',
 groups: operation.groups.map(progress => ({ ...progress, ready: false, applied: false, applied_revision: undefined })),
 devices: operation.devices.map(progress => ({ ...progress, activated: false })) };
it('shows a failed Start ACK as a failure with its safe reason, never as an older completed start', async () => {
 api.start.mockResolvedValue(failedReadiness); api.operation.mockResolvedValue(failedReadiness);
 const onCompleted = vi.fn(); mount(undefined, { onCompleted });
 fireEvent.click(screen.getByTestId('saved-recording-start-confirm'));
 expect(await screen.findByText('step4.basic.operation.failed')).toBeInTheDocument();
 expect(screen.getByText('step4.basic.reason.invalid_scope')).toBeInTheDocument();
 expect(screen.getByText('step4.basic.action.review_selection')).toBeInTheDocument();
 expect(screen.queryByText('start.previousTerminal')).not.toBeInTheDocument();
 expect(screen.queryByTestId('saved-recording-start-waiting')).not.toBeInTheDocument();
 const check = screen.getByTestId('saved-recording-start-check');
 await waitFor(() => expect(check).toBeEnabled()); fireEvent.click(check);
 await waitFor(() => expect(api.operation).toHaveBeenCalledWith(operation.operation_id));
 expect(api.start).toHaveBeenCalledTimes(1); expect(onCompleted).not.toHaveBeenCalled();
 expect(JSON.parse(sessionStorage.getItem(basicRecordingIntentKey(`ws-1/saved-record/${group.id}`))!).operation_id).toBe(operation.operation_id);
});
it('reloads a retained failed Start as failure without replaying or calling it historical success', async () => {
 saveBasicRecordingIntent({ scope_key: `ws-1/saved-record/${group.id}`, owner: { kind: 'saved_record', groupId: group.id, deviceIds: ['dev-1', 'dev-2'] }, operation_id: operation.operation_id,
  request: { request_id: 'failed-retained-request', workspace_id: 'ws-1', expected_workspace_revision: 'w1', device_ids: ['dev-1', 'dev-2'], groups: [{ group_id: group.id, expected_group_revision: group.revision, expected_connector_revision: group.destination.connector_revision }] } });
 api.operation.mockResolvedValue(failedReadiness);
 const onCompleted = vi.fn(); mount(undefined, { onCompleted });
 expect(await screen.findByText('step4.basic.operation.failed')).toBeInTheDocument();
 expect(screen.getByText('step4.basic.reason.invalid_scope')).toBeInTheDocument();
 expect(screen.queryByText('start.previousTerminal')).not.toBeInTheDocument();
 expect(screen.queryByTestId('saved-recording-start-waiting')).not.toBeInTheDocument();
 expect(screen.getByTestId('saved-recording-start-confirm')).toBeEnabled();
 expect(api.start).not.toHaveBeenCalled(); expect(onCompleted).not.toHaveBeenCalled();
});
it('renders a safe fallback for an unknown failed Start reason without exposing raw diagnostics', async () => {
 const failed = { ...failedReadiness, reason: 'private-endpoint-and-credential', next_action: 'raw-private-action' };
 api.start.mockResolvedValue(failed); api.operation.mockResolvedValue(failed);
 mount(); fireEvent.click(screen.getByTestId('saved-recording-start-confirm'));
 expect(await screen.findByText('step4.basic.operation.failed')).toBeInTheDocument();
 expect(screen.getByText('step4.basic.reason.unconfirmed')).toBeInTheDocument();
 expect(screen.queryByText('start.previousTerminal')).not.toBeInTheDocument();
 expect(document.body.textContent).not.toContain('private-endpoint-and-credential');
 expect(document.body.textContent).not.toContain('raw-private-action');
 expect(api.start).toHaveBeenCalledTimes(1);
});
it.each(['toString', 'constructor', '__proto__'])('does not translate inherited failed Start reason or action %s', key => {
 const failed = { ...failedReadiness, reason: key, next_action: key };
 const translate = vi.fn((message: string) => message);
 expect(safeBasicOperationMessage(failed, translate)).toBe('step4.basic.reason.unconfirmed');
 expect(safeBasicNextAction(failed, translate)).toBeUndefined();
 expect(translate).toHaveBeenCalledExactlyOnceWith('step4.basic.reason.unconfirmed');
});
