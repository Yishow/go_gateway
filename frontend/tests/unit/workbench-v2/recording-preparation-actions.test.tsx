import { recordingTestWriteRecovery } from '../../../src/features/datalink/workbench-v2/state/recordingTestWriteRecovery';
import { readSchemaOperation } from '../../../src/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { GroupEditor } from '../../../src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { buildGroupCandidates } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
import { batchPanelState, longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { managedSchemaPreviewToken } from '../../fixtures/managedSchemaPreview';
import { TABLE_COLUMNS } from '../../fixtures/writeGroupState';
import { studioV2WorkspaceWriteGroupsAPI } from '../../../src/services/studioV2WorkspaceWriteGroups';
import { studioV2RecordingStartAPI } from '../../../src/services/studioV2RecordingStart';
import type { SchemaOperation } from '../../../src/types/recordingPlan';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('../../../src/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({ useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }) }));
vi.mock('../../../src/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
 readiness: vi.fn(), delivery: vi.fn().mockResolvedValue(null), list: vi.fn(), get: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(), testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(),
} }));
vi.mock('../../../src/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), get: vi.fn() } }));
const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI), state = batchPanelState();
const group = longProfileGroup({ id: 'prepare-record', circuitId: 'c1', circuitName: 'My recording', interval: 15, points: [0, 1] });
const preview = { ...managedSchemaPreviewToken, workspace_id: 'ws-1', workspace_revision: 'w1', plan_id: group.id, plan_revision: group.revision,
 connector_id: 'conn-1', connector_revision: 'crev-1', expires_at: '2099-01-01T00:00:00Z', group_layout: { ...managedSchemaPreviewToken.group_layout!, indexes: [{ name: 'scope_time', columns: ['group_id', 'bucket_start'] }] } };
const terminal: SchemaOperation = { operation_id: preview.operation_id, status: 'succeeded', executed_statements: 2, created_at: '2026-10-31T00:00:00Z', updated_at: '2026-10-31T00:00:01Z' };
function mount() {
 const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
 const onClose = vi.fn();
 let revision = 'w1'; let readonly = false; let candidates = buildGroupCandidates(state).candidates;
 const tree = () => <QueryClientProvider client={client}><GroupEditor group={group} workspaceId="ws-1" workspaceRevision={revision} connector={state.db.connector} candidates={candidates} excluded={[]} readonly={readonly} onClose={onClose} onSaved={vi.fn()} onReload={vi.fn()} fetchLatest={async () => ({ workspace_revision: revision, groups: [group] })} /></QueryClientProvider>;
 const view = render(tree()); return { ...view, onClose, changeSource: () => { candidates = candidates.map((candidate, i) => i === 0 ? { ...candidate, address: '40099' } : candidate); view.rerender(tree()); }, resetSource: () => { candidates = buildGroupCandidates(state).candidates; view.rerender(tree()); }, changeRevision: () => { revision = 'w2'; view.rerender(tree()); }, setReadonly: (value: boolean) => { readonly = value; view.rerender(tree()); } };
}
beforeEach(() => { vi.restoreAllMocks(); vi.clearAllMocks(); sessionStorage.clear(); api.schemaPreview.mockResolvedValue(preview); api.schemaApply.mockResolvedValue(terminal); api.schemaOperation.mockResolvedValue(terminal);
 api.readiness.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'w1', group_id: group.id, group_revision: group.revision, applied_revision: '', config_ready: true, schema_ready: false, ready: false, issues: [] });
 api.update.mockResolvedValue({ workspace_revision: 'w2', group: { ...group, revision: 'new-own-revision', name: 'My changed recording' } }); });
async function checkImpact() { await waitFor(() => expect(screen.getByTestId('recording-preparation-check')).toBeEnabled()); fireEvent.click(screen.getByTestId('recording-preparation-check')); }
function noStart() { expect(studioV2RecordingStartAPI.start).not.toHaveBeenCalled(); expect(api.apply).not.toHaveBeenCalled(); expect(api.testWrite).not.toHaveBeenCalled(); }
it('Save stays separate from explicit check, preview and starting', async () => {
 mount(); fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'My changed recording' } });
 expect(screen.getByTestId('recording-preparation-phase')).toHaveTextContent('preparation.stage.save');
 fireEvent.click(screen.getByTestId('group-save'));
 await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
 expect(await screen.findByTestId('group-saved-note')).toBeInTheDocument();
 expect(api.schemaPreview).not.toHaveBeenCalled(); expect(api.schemaApply).not.toHaveBeenCalled(); noStart();
});
it('explicit check exposes all real table impacts and requires separate schema consent without Start', async () => {
 mount(); expect(api.schemaPreview).not.toHaveBeenCalled();
 await checkImpact();
 expect(await screen.findByTestId('group-schema-preview-result')).toHaveTextContent('gw_effect_receipts');
 expect(screen.getByTestId('group-schema-preview-result')).toHaveTextContent('managed.db');
 expect(screen.getByTestId('group-schema-preview-result')).toHaveClass('bg-[var(--gw-surface)]', 'text-[var(--gw-text)]');
 expect(screen.getByTestId('recording-preparation-actions')).toHaveClass('bg-[var(--gw-surface)]', 'text-[var(--gw-text)]', 'border-[var(--gw-border)]');
 expect(screen.getByTestId('group-schema-preview-result')).toHaveTextContent('scope_time');
 expect(screen.getByTestId('recording-preparation-phase')).toHaveTextContent('preparation.stage.review');
 expect(api.schemaApply).not.toHaveBeenCalled(); noStart();
 const confirm = screen.getByTestId('recording-preparation-confirm');
 expect(confirm).toHaveClass('max-w-full', 'h-auto', 'min-h-9', 'whitespace-normal', 'break-words', 'py-2');
 fireEvent.click(confirm); fireEvent.click(confirm);
 await waitFor(() => expect(api.schemaApply).toHaveBeenCalledTimes(1)); noStart();
 expect(api.schemaApply).toHaveBeenCalledWith(group.id, expect.objectContaining({ operation_id: preview.operation_id, token: preview.token }));
 await waitFor(() => expect(readSchemaOperation('ws-1', group.id)).toBeUndefined());
 expect(screen.getByTestId('group-schema-operation')).toHaveTextContent('step4.group.schema.operation_status.succeeded');
});
it('selected source changes retire old consent even when the saved group revision and dirty flag stay unchanged', async () => {
 const view = mount(); await checkImpact();
 await screen.findByTestId('recording-preparation-confirm'); view.changeSource();
 expect(screen.queryByTestId('recording-preparation-confirm')).not.toBeInTheDocument();
 expect(screen.queryByTestId('group-schema-preview-result')).not.toBeInTheDocument();
 expect(api.schemaApply).not.toHaveBeenCalled(); noStart();
});

it('edit then revert cannot revive the old schema consent', async () => {
 mount(); await checkImpact(); await screen.findByTestId('recording-preparation-confirm');
 fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Changed intent' } });
 fireEvent.change(screen.getByTestId('group-name'), { target: { value: group.name } });
 expect(screen.queryByTestId('recording-preparation-confirm')).not.toBeInTheDocument();
 expect(screen.queryByTestId('group-schema-preview-result')).not.toBeInTheDocument();
 expect(api.schemaPreview).toHaveBeenCalledTimes(1); expect(api.schemaApply).not.toHaveBeenCalled(); noStart();
});
it.each(['workspace', 'source round trip'] as const)('a late preview after %s changes cannot publish confirmation or send a replacement command', async kind => {
 let resolve!: (value: typeof preview) => void; api.schemaPreview.mockImplementation(() => new Promise(done => { resolve = done; }));
 const view = mount(); await checkImpact(); await waitFor(() => expect(api.schemaPreview).toHaveBeenCalledTimes(1));
 if (kind === 'workspace') view.changeRevision(); else { view.changeSource(); view.resetSource(); }
 await act(async () => { resolve(preview); });
 expect(screen.queryByTestId('recording-preparation-confirm')).not.toBeInTheDocument();
 expect(screen.queryByTestId('group-schema-preview-result')).not.toBeInTheDocument();
 expect(api.schemaPreview).toHaveBeenCalledTimes(1); expect(api.schemaApply).not.toHaveBeenCalled(); noStart();
});
it('a lost schema acknowledgement offers only lookup of the same original operation even after editing', async () => {
 api.schemaApply.mockRejectedValue(new Error('Lost acknowledgement'));
 mount(); await checkImpact(); fireEvent.click(await screen.findByTestId('recording-preparation-confirm'));
 await screen.findByTestId('group-schema-error');
 fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Unsaved while checking' } });
 expect(screen.queryByTestId('group-save')).not.toBeInTheDocument();
 fireEvent.click(screen.getByTestId('recording-preparation-lookup'));
 await waitFor(() => expect(api.schemaOperation).toHaveBeenCalledWith(preview.operation_id));
 await waitFor(() => expect(screen.getByTestId('group-save')).toBeInTheDocument());
 expect(api.schemaPreview).toHaveBeenCalledTimes(1); expect(api.schemaApply).toHaveBeenCalledTimes(1);
 expect(api.update).not.toHaveBeenCalled(); noStart();
});
it('offline and readonly states preserve the form but deny fresh schema commands', async () => {
 const online = vi.spyOn(navigator, 'onLine', 'get').mockReturnValue(false);
 const view = mount(); await waitFor(() => expect(screen.getByTestId('recording-preparation-check')).toBeDisabled());
 expect(screen.getByTestId('recording-preparation-phase')).toHaveTextContent('preparation.stage.offline');
 online.mockReturnValue(true); act(() => window.dispatchEvent(new Event('online')));
 await checkImpact(); await screen.findByTestId('recording-preparation-confirm'); view.setReadonly(true);
 expect(screen.queryByTestId('recording-preparation-confirm')).not.toBeInTheDocument();
 expect(screen.getByTestId('group-save')).toBeDisabled(); expect(api.schemaApply).not.toHaveBeenCalled(); noStart();
});

it('a known preview-only failure can be explicitly retried after native source A to B to A without borrowing old consent', async () => {
 api.schemaPreview.mockRejectedValueOnce(new Error('Preview unavailable'));
 const view = mount(); await checkImpact(); await screen.findByTestId('group-schema-error');
 view.changeSource(); view.resetSource();
 expect(api.schemaPreview).toHaveBeenCalledTimes(1); expect(api.schemaApply).not.toHaveBeenCalled();
 await checkImpact(); await screen.findByTestId('recording-preparation-confirm');
 expect(api.schemaPreview).toHaveBeenCalledTimes(2); expect(api.schemaApply).not.toHaveBeenCalled(); noStart();
});
it('native source changes invalidate prior readiness until a fresh matching check, even with unchanged saved revision', async () => {
 api.readiness.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'w1', group_id: group.id, group_revision: group.revision, applied_revision: '', config_ready: true, schema_ready: true, ready: true, issues: [] });
 const view = mount(); await waitFor(() => expect(screen.getByTestId('recording-preparation-phase')).toHaveTextContent('preparation.stage.ready'));
 view.changeSource(); expect(screen.getByTestId('recording-preparation-phase')).not.toHaveTextContent('preparation.stage.ready');
 expect(api.schemaPreview).not.toHaveBeenCalled(); expect(api.schemaApply).not.toHaveBeenCalled(); noStart();
});

it('Close during a submitted schema change uses the same pending guard and never cancels its original operation', async () => {
 let resolve!: (value: SchemaOperation) => void; api.schemaApply.mockImplementation(() => new Promise(done => { resolve = done; }));
 const view = mount(); await checkImpact(); fireEvent.click(await screen.findByTestId('recording-preparation-confirm'));
 await waitFor(() => expect(api.schemaApply).toHaveBeenCalledTimes(1));
 fireEvent.click(screen.getByTestId('group-editor-close'));
 expect(await screen.findByRole('alertdialog')).toHaveTextContent('navigation.pending');
 expect(view.onClose).not.toHaveBeenCalled(); fireEvent.click(screen.getByText('navigation.keep'));
 expect(readSchemaOperation('ws-1', group.id)).toBe(preview.operation_id);
 fireEvent.click(screen.getByTestId('group-editor-close')); fireEvent.click(await screen.findByText('navigation.discard'));
 expect(view.onClose).toHaveBeenCalledTimes(1); expect(readSchemaOperation('ws-1', group.id)).toBe(preview.operation_id);
 await act(async () => { resolve(terminal); }); await waitFor(() => expect(readSchemaOperation('ws-1', group.id)).toBeUndefined());
 expect(api.schemaPreview).toHaveBeenCalledTimes(1); expect(api.schemaApply).toHaveBeenCalledTimes(1); noStart();
});

it('a preview-only in-flight read is not presented as a submitted schema change when leaving', async () => {
 let resolve!: (value: typeof preview) => void; api.schemaPreview.mockImplementation(() => new Promise(done => { resolve = done; }));
 const view = mount(); await checkImpact(); await waitFor(() => expect(api.schemaPreview).toHaveBeenCalledTimes(1));
 fireEvent.click(screen.getByTestId('group-editor-close'));
 expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument(); expect(view.onClose).toHaveBeenCalledTimes(1);
 view.unmount(); await act(async () => { resolve(preview); });
 expect(readSchemaOperation('ws-1', group.id)).toBeUndefined(); expect(api.schemaApply).not.toHaveBeenCalled(); noStart();
});

it.each(['preview', 'confirm'] as const)('rechecks a newly retained test-write reference before %s without waiting for subscription render', async stage => {
 mount();
 if (stage === 'confirm') { await checkImpact(); await screen.findByTestId('recording-preparation-confirm'); }
 else await waitFor(() => expect(screen.getByTestId('recording-preparation-check')).toBeEnabled());
 // Another owner retains cleanup after this render; the command must inspect the original live store.
 const read = vi.spyOn(recordingTestWriteRecovery, 'read').mockReturnValue({ references: [{ workspace_id: 'ws-1', group_id: group.id, group_revision: group.revision,
  operation_id: 'retained-cleanup', operation_kind: 'test_write', created_at: '2026-10-31T00:00:00Z' }], unavailable: true, unreadable: true });
 await act(async () => { fireEvent.click(screen.getByTestId(stage === 'confirm' ? 'recording-preparation-confirm' : 'recording-preparation-check')); });
 expect(api.schemaPreview).toHaveBeenCalledTimes(stage === 'confirm' ? 1 : 0);
 expect(api.schemaApply).not.toHaveBeenCalled(); expect(read).toHaveBeenCalled(); noStart();
});
it('an unreadable empty optional test-write store does not introduce a mandatory test-write prerequisite', async () => {
 vi.spyOn(recordingTestWriteRecovery, 'read').mockReturnValue({ references: [], unavailable: true, unreadable: true });
 mount(); await checkImpact();
 expect(await screen.findByTestId('recording-preparation-confirm')).toBeEnabled();
 expect(api.schemaPreview).toHaveBeenCalledTimes(1); noStart();
});
