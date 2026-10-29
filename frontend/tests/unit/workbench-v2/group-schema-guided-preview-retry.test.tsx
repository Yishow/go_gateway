import type { ComponentProps } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { GroupSchemaPanel } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel';
import type { SchemaPreparationPresentation } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/schemaPreparationPresentation';
import { inspectSchemaOperation, readSchemaOperation, retainSchemaOperation } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { managedSchemaPreviewToken } from '../../fixtures/managedSchemaPreview';
import { savedGroup } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(),
} }));
const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const group = savedGroup({ id: 'group-1', workspace_id: 'workspace-1', revision: 'group-rev-1', destination: {
  connector_id: 'connector-1', connector_revision: 'connector-rev-1', database: 'mock-private', table_schema: 'main',
  table_name: 'readings', storage_strategy: 'managed', schema_revision: 'group-rev-1', schema_digest: 'a'.repeat(64),
} });
const failure = { response: { status: 503, data: { success: false, error: { code: 'RECORDING_SCHEMA_TARGET_UNAVAILABLE' } } } };
const clients: QueryClient[] = [];
let report: SchemaPreparationPresentation | null;
let retryBlocked: boolean, commandBlocked: boolean;
const observe = (value: SchemaPreparationPresentation | null) => { report = value; };
function show(patch: Partial<ComponentProps<typeof GroupSchemaPanel>> = {}) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } }); clients.push(client);
  const props: ComponentProps<typeof GroupSchemaPanel> = { group, workspaceId: 'workspace-1', workspaceRevision: 'W1', readonly: false,
    dirty: false, destinationMatches: true, guided: { online: true, suppressDefaultActions: true, onPresentation: observe,
      isNewOperationBlocked: () => commandBlocked, isPreviewRetryBlocked: () => retryBlocked }, ...patch };
  const tree = (value: typeof props) => <QueryClientProvider client={client}><GroupSchemaPanel {...value} /></QueryClientProvider>;
  const view = render(tree(props));
  return { ...view, rerenderPanel: (next: Partial<typeof props> = {}) => view.rerender(tree({ ...props, ...next })) };
}
async function failPreview() {
  const preview = report?.commands.preview; expect(preview).toBeTypeOf('function');
  await act(async () => { await preview?.(); });
  await screen.findByTestId('group-schema-error'); expect(api.schemaPreview).toHaveBeenCalledOnce();
}
function noDangerousCommand() { expect(api.schemaApply).not.toHaveBeenCalled(); expect(api.schemaOperation).not.toHaveBeenCalled(); }
beforeEach(() => {
  vi.restoreAllMocks(); vi.clearAllMocks(); sessionStorage.clear(); report = null; retryBlocked = false; commandBlocked = false;
  for (const id of ['group-1', 'group-2']) retainSchemaOperation('workspace-1', id, undefined, readSchemaOperation('workspace-1', id));
  api.schemaPreview.mockRejectedValueOnce(failure).mockResolvedValue(managedSchemaPreviewToken);
  api.schemaApply.mockRejectedValue(new Error('Acknowledgement lost'));
  api.schemaOperation.mockRejectedValue(new Error('Lookup unavailable'));
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); vi.restoreAllMocks(); });

it('offers explicit preview-only retry after 503 without treating a saved ready group as healthy or submitting DDL', async () => {
  const view = show(); await failPreview();
  expect(report?.previewFailure).toBe(true); expect(report?.observation).toMatchObject({ state: 'unavailable' });
  expect(report?.phase).toBe('unavailable'); expect(report?.commands.preview).toBeUndefined(); expect(report?.commands.confirm).toBeUndefined();
  expect(report?.commands.retryPreview).toBeTypeOf('function'); expect(inspectSchemaOperation('workspace-1', group.id)).toEqual({ state: 'empty' });
  const retry = screen.getByTestId('group-schema-preview'); expect(retry).toBeEnabled(); expect(retry).not.toHaveAttribute('data-primary', 'true');
  expect(screen.getAllByTestId('group-schema-preview')).toHaveLength(1);
  view.rerenderPanel(); await act(async () => undefined); expect(api.schemaPreview).toHaveBeenCalledOnce(); noDangerousCommand();
  fireEvent.click(retry); fireEvent.click(retry); await screen.findByTestId('group-schema-preview-result');
  expect(api.schemaPreview).toHaveBeenCalledTimes(2); expect(api.schemaPreview.mock.calls[1]).toEqual(api.schemaPreview.mock.calls[0]);
  await waitFor(() => expect(report?.preview).toBe('ready')); expect(report?.previewFailure).toBeUndefined(); noDangerousCommand();
});
it('uses the specialized denial only for explicit retry; global new-command denial remains effective', async () => {
  show(); await failPreview(); const retry = report?.commands.retryPreview; expect(retry).toBeTypeOf('function');
  commandBlocked = true; await act(async () => { await retry?.(); });
  expect(api.schemaPreview).toHaveBeenCalledTimes(2); await screen.findByTestId('group-schema-preview-result');
  const confirm = report?.commands.confirm; await act(async () => { await confirm?.(); });
  expect(api.schemaApply).not.toHaveBeenCalled();
});
it.each(['missing', 'true', 'throw'] as const)('denies preview retry when specialized denial is %s', async condition => {
  const view = show(); await failPreview(); const old = report?.commands.retryPreview;
  const deny = condition === 'missing' ? undefined : condition === 'throw' ? () => { throw new Error('unreadable sibling'); } : () => true;
  view.rerenderPanel({ guided: { online: true, suppressDefaultActions: true, onPresentation: observe, isPreviewRetryBlocked: deny } });
  expect(report?.commands.retryPreview).toBeUndefined(); await act(async () => { await old?.(); });
  expect(api.schemaPreview).toHaveBeenCalledOnce(); noDangerousCommand();
});
it.each(['readonly', 'dirty', 'destination', 'deleted', 'offline'] as const)('keeps native %s denial on a captured retry', async condition => {
  const view = show(); await failPreview(); const retry = report?.commands.retryPreview;
  const patch: Partial<ComponentProps<typeof GroupSchemaPanel>> = condition === 'readonly' ? { readonly: true } : condition === 'dirty' ? { dirty: true }
    : condition === 'destination' ? { destinationMatches: false } : condition === 'deleted' ? { group: { ...group, status: 'deleted' } }
      : { guided: { online: false, onPresentation: observe, suppressDefaultActions: true, isPreviewRetryBlocked: () => false } };
  view.rerenderPanel(patch); await act(async () => { await retry?.(); });
  expect(report?.commands.retryPreview).toBeUndefined(); expect(api.schemaPreview).toHaveBeenCalledOnce(); noDangerousCommand();
});
it.each(['retained', 'unreadable', 'other-owner'] as const)('freshly rejects %s evidence appearing after retry render', async condition => {
  show(); await failPreview(); const retry = report?.commands.retryPreview; expect(retry).toBeTypeOf('function');
  if (condition === 'retained') retainSchemaOperation('workspace-1', group.id, 'O1');
  if (condition === 'unreadable') vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => { throw new Error('storage unavailable'); });
  if (condition === 'other-owner') retryBlocked = true;
  await act(async () => { await retry?.(); }); expect(api.schemaPreview).toHaveBeenCalledOnce(); noDangerousCommand();
  if (condition === 'retained') expect(readSchemaOperation('workspace-1', group.id)).toBe('O1');
});
it('does not offer preview retry after an apply lost ACK; original O1 remains query-only', async () => {
  api.schemaPreview.mockReset().mockResolvedValue(managedSchemaPreviewToken); show();
  await act(async () => { await report?.commands.preview?.(); });
  await act(async () => { await report?.commands.confirm?.(); });
  expect(report?.observation).toMatchObject({ state: 'unresolved', operationId: 'operation-1' });
  expect(report?.previewFailure).toBeUndefined(); expect(report?.commands.retryPreview).toBeUndefined();
  expect(readSchemaOperation('workspace-1', group.id)).toBe('operation-1'); expect(api.schemaApply).toHaveBeenCalledOnce();
  await act(async () => { await report?.commands.check?.(); });
  expect(report?.commands.retryPreview).toBeUndefined(); expect(api.schemaOperation).toHaveBeenCalledExactlyOnceWith('operation-1');
  expect(api.schemaPreview).toHaveBeenCalledOnce(); expect(readSchemaOperation('workspace-1', group.id)).toBe('operation-1');
});
it('never adopts a late failed preview from A → B → A as current retry permission or error', async () => {
  let rejectPreview: ((cause: unknown) => void) | undefined;
  api.schemaPreview.mockReset().mockImplementationOnce(() => new Promise((_, reject) => { rejectPreview = reject; }));
  const view = show(); let pending: Promise<void> | undefined;
  await act(async () => { pending = report?.commands.preview?.(); });
  view.rerenderPanel({ group: { ...group, id: 'group-2' } }); view.rerenderPanel();
  await act(async () => { rejectPreview?.(failure); await pending; });
  expect(report?.previewFailure).toBeUndefined(); expect(report?.commands.retryPreview).toBeUndefined();
  expect(screen.queryByTestId('group-schema-error')).not.toBeInTheDocument(); expect(api.schemaPreview).toHaveBeenCalledOnce(); noDangerousCommand();
});
it('retires a captured retry across a scope round trip and unmount', async () => {
  const view = show(); await failPreview(); const old = report?.commands.retryPreview;
  view.rerenderPanel({ group: { ...group, id: 'group-2' } }); view.rerenderPanel();
  await act(async () => { await old?.(); }); view.unmount(); await act(async () => { await old?.(); });
  expect(api.schemaPreview).toHaveBeenCalledOnce(); noDangerousCommand();
});
it('does not borrow a new parent lease after native-consent A → B → A with unchanged schema scope', async () => {
  let generation = 0;
  const bridge = (captured: number) => ({ online: true, consentEpoch: captured, suppressDefaultActions: true, onPresentation: observe,
    isPreviewRetryBlocked: () => captured !== generation });
  const view = show({ guided: bridge(0) }); await failPreview();
  const original = report?.commands.retryPreview; expect(original).toBeTypeOf('function');
  const scope = report?.scope;
  generation = 1; view.rerenderPanel({ guided: bridge(1) });
  generation = 2; view.rerenderPanel({ guided: bridge(2) });
  expect(report?.scope).toMatchObject({ workspaceId: scope?.workspaceId, workspaceRevision: scope?.workspaceRevision,
    groupId: scope?.groupId, groupRevision: scope?.groupRevision }); // Saved schema identity did not change.
  const current = report?.commands.retryPreview;
  await act(async () => { await original?.(); });
  expect(api.schemaPreview).toHaveBeenCalledOnce(); noDangerousCommand();
  expect(current).toBeTypeOf('function');
  await act(async () => { await current?.(); }); await screen.findByTestId('group-schema-preview-result');
  expect(api.schemaPreview).toHaveBeenCalledTimes(2); noDangerousCommand();
});
it('can explicitly retry when equivalent parent callbacks are republished without a consent change', async () => {
  const bridge = () => ({ online: true, consentEpoch: 0, suppressDefaultActions: true, onPresentation: observe, isPreviewRetryBlocked: () => false });
  const view = show({ guided: bridge() }); await failPreview(); const retry = report?.commands.retryPreview;
  view.rerenderPanel({ guided: bridge() }); view.rerenderPanel({ guided: bridge() });
  await act(async () => { await retry?.(); }); await screen.findByTestId('group-schema-preview-result');
  expect(api.schemaPreview).toHaveBeenCalledTimes(2); noDangerousCommand();
});
