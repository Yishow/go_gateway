import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { GroupSchemaPanel } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel';
import type { SchemaPreparationPresentation } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/schemaPreparationPresentation';
import { readSchemaOperation, retainSchemaOperation } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery';
import { managedSchemaPreviewToken } from '../../fixtures/managedSchemaPreview';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { WriteGroup } from '@/types/studioV2WriteGroup';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(),
} }));
const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const group: WriteGroup = { id: 'group-1', workspace_id: 'workspace-1', revision: 'group-rev-1', applied_revision: '', name: 'Group', status: 'ready', members: [],
  destination: { connector_id: 'connector-1', connector_revision: 'connector-rev-1', database: 'private.db', table_schema: 'main', table_name: 'readings', storage_strategy: 'managed' },
  row_policy: { interval_seconds: 30, allowed_lateness_seconds: 0 }, write_policy: {}, migration: {}, created_at: '', updated_at: '' };
const succeeded = { operation_id: 'operation-1', action: 'schema_apply', status: 'succeeded' as const, executed_statements: 1,
  created_at: '2026-10-28T00:00:00Z', updated_at: '2026-10-28T00:00:01Z' };
let report: SchemaPreparationPresentation | null;
const reports: Array<SchemaPreparationPresentation | null> = [];
const onPresentation = (value: SchemaPreparationPresentation | null) => { report = value; reports.push(value); };
function mount(overrides: Partial<React.ComponentProps<typeof GroupSchemaPanel>> = {}) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const props = { group, workspaceId: 'workspace-1', workspaceRevision: 'workspace-rev-1', destinationMatches: true, dirty: false, readonly: false,
    guided: { online: true, suppressDefaultActions: true, onPresentation }, ...overrides };
  const tree = (next: typeof props) => <QueryClientProvider client={client}><details open><summary>Inspect schema</summary><GroupSchemaPanel {...next} /></details></QueryClientProvider>;
  const mounted = render(tree(props));
  return { ...mounted, rerenderPanel: (next: Partial<typeof props>) => mounted.rerender(tree({ ...props, ...next })) };
}
beforeEach(() => {
  vi.restoreAllMocks(); vi.clearAllMocks(); sessionStorage.clear(); report = null; reports.length = 0;
  for (const id of ['group-1', 'group-2']) retainSchemaOperation('workspace-1', id, undefined, readSchemaOperation('workspace-1', id));
  api.schemaPreview.mockResolvedValue(managedSchemaPreviewToken); api.schemaApply.mockResolvedValue(succeeded);
  api.schemaOperation.mockResolvedValue({ ...succeeded, status: 'unknown' });
});
afterEach(() => { cleanup(); vi.restoreAllMocks(); });

describe('guided schema owner bridge', () => {
  it('mounts, opens and focuses schema details without any preview POST or DDL, even when standalone autoPreview is requested', async () => {
    mount({ autoPreview: true });
    fireEvent.click(screen.getByText('Inspect schema')); fireEvent.click(screen.getByText('Inspect schema'));
    screen.getByTestId('group-schema-panel').focus();
    await act(async () => undefined);
    expect(api.schemaPreview).not.toHaveBeenCalled(); expect(api.schemaApply).not.toHaveBeenCalled(); expect(api.schemaOperation).not.toHaveBeenCalled();
  });
  it('exposes a token-free scoped owner report instead of duplicate default action buttons', async () => {
    mount();
    await waitFor(() => expect(report).not.toBeNull());
    expect(report?.scope).toMatchObject({ workspaceId: 'workspace-1', groupId: 'group-1', workspaceRevision: 'workspace-rev-1', groupRevision: 'group-rev-1', connectorId: 'connector-1', connectorRevision: 'connector-rev-1' });
    expect(report?.phase).toBe('idle'); expect(report?.observation).toMatchObject({ owner: 'schema', state: 'observed', result: 'none' });
    expect(report?.commands.preview).toBeTypeOf('function'); expect(report?.commands.confirm).toBeUndefined();
    expect(screen.queryByTestId('group-schema-preview')).not.toBeInTheDocument();
    expect(JSON.stringify(report)).not.toMatch(/token|statements|private.db|CREATE TABLE/);
  });
  it('previews only on an explicit command and confirms through the same current owner handler once', async () => {
    mount(); const initial = report!;
    await act(async () => { await initial.commands.preview?.(); });
    expect(await screen.findByTestId('group-schema-preview-result')).toBeInTheDocument();
    expect(api.schemaPreview).toHaveBeenCalledTimes(1); expect(api.schemaApply).not.toHaveBeenCalled();
    await waitFor(() => expect(report?.commands.confirm).toBeTypeOf('function'));
    const confirm = report!.commands.confirm!;
    await act(async () => { await Promise.all([confirm(), confirm()]); });
    expect(api.schemaApply).toHaveBeenCalledTimes(1);
    expect(api.schemaApply).toHaveBeenCalledWith('group-1', expect.objectContaining({ token: 'token-1', operation_id: 'operation-1', expected_group_revision: 'group-rev-1' }));
    expect(report?.observation).toMatchObject({ state: 'observed', result: 'resolved', status: 'succeeded' });
    expect(JSON.stringify(report)).not.toMatch(/token-1|CREATE TABLE|statements/);
  });
  it('reports the opaque preview identity only for the current ready preview, never as an operation observation', async () => {
    const mounted = mount();
    expect(report?.previewOperationId).toBeUndefined();
    await act(async () => { await report?.commands.preview?.(); });
    expect(report?.preview).toBe('ready'); expect(report?.previewOperationId).toBe('operation-1');
    expect(report?.observation).toEqual({ owner: 'schema', state: 'observed', result: 'none' });
    expect(JSON.stringify(report)).not.toMatch(/token-1|statements|CREATE TABLE/);
    const firstReport = reports.length;
    mounted.rerenderPanel({ workspaceRevision: 'workspace-rev-2' });
    expect(reports.slice(firstReport).filter(value => value !== null).every(value => value.previewOperationId === undefined && value.preview === 'none')).toBe(true);
    expect(report?.commands.confirm).toBeUndefined();
  });
  it('does not revive an in-flight preview after an A → B → A scope round trip', async () => {
    let resolvePreview!: (value: typeof managedSchemaPreviewToken) => void;
    api.schemaPreview.mockImplementation(() => new Promise(resolve => { resolvePreview = resolve; }));
    const mounted = mount(); let pending!: Promise<void>;
    await act(async () => { pending = report!.commands.preview!(); });
    mounted.rerenderPanel({ workspaceRevision: 'workspace-rev-2' }); mounted.rerenderPanel({});
    await act(async () => { resolvePreview(managedSchemaPreviewToken); await pending; });
    expect(report?.previewOperationId).toBeUndefined(); expect(report?.preview).toBe('none');
    expect(report?.commands.confirm).toBeUndefined(); expect(screen.queryByTestId('group-schema-preview-result')).not.toBeInTheDocument();
    expect(api.schemaApply).not.toHaveBeenCalled();
  });
  it('checks retained O1 and never creates a replacement preview or DDL', async () => {
    retainSchemaOperation('workspace-1', group.id, 'O1'); api.schemaOperation.mockResolvedValue({ ...succeeded, operation_id: 'O1', status: 'unknown' });
    mount({ autoPreview: true });
    expect(report?.observation).toMatchObject({ state: 'unresolved', operationId: 'O1' });
    expect(report?.commands.preview).toBeUndefined(); expect(report?.commands.confirm).toBeUndefined();
    await act(async () => { await report?.commands.check?.(); });
    expect(api.schemaOperation).toHaveBeenCalledWith('O1');
    expect(api.schemaPreview).not.toHaveBeenCalled(); expect(api.schemaApply).not.toHaveBeenCalled();
  });
  it.each(['workspace', 'revision', 'connector', 'dirty', 'offline'] as const)('rejects a captured callback after %s changes', async (change) => {
    const mounted = mount(); const oldPreview = report!.commands.preview!;
    const next = change === 'workspace' ? { workspaceId: 'workspace-2', group: { ...group, workspace_id: 'workspace-2' } }
      : change === 'revision' ? { workspaceRevision: 'workspace-rev-2' }
        : change === 'connector' ? { group: { ...group, destination: { ...group.destination, connector_id: 'connector-2' } } }
          : change === 'dirty' ? { dirty: true } : { guided: { online: false, suppressDefaultActions: true, onPresentation } };
    mounted.rerenderPanel(next);
    await act(async () => { await oldPreview(); });
    expect(api.schemaPreview).not.toHaveBeenCalled(); expect(api.schemaApply).not.toHaveBeenCalled();
  });
  it('does not revive a captured A callback after A → B → A with identical IDs and revisions', async () => {
    const mounted = mount(); const first = report!; const oldPreview = first.commands.preview!;
    mounted.rerenderPanel({ workspaceId: 'workspace-2', group: { ...group, workspace_id: 'workspace-2' } });
    mounted.rerenderPanel({});
    expect(report!.scope.epoch).toBeGreaterThan(first.scope.epoch);
    await act(async () => { await oldPreview(); });
    expect(api.schemaPreview).not.toHaveBeenCalled();
    await act(async () => { await report?.commands.preview?.(); });
    expect(api.schemaPreview).toHaveBeenCalledTimes(1);
  });
  it('does not apply a captured preview confirmation after a scope round trip or unmount', async () => {
    const mounted = mount(); await act(async () => { await report?.commands.preview?.(); });
    const oldConfirm = report!.commands.confirm!;
    mounted.rerenderPanel({ workspaceRevision: 'workspace-rev-2' }); mounted.rerenderPanel({});
    await act(async () => { await oldConfirm(); }); expect(api.schemaApply).not.toHaveBeenCalled();
    const newPreview = report!.commands.preview!; mounted.unmount();
    await act(async () => { await newPreview(); }); expect(api.schemaPreview).toHaveBeenCalledTimes(1);
  });
  it('fails closed for unreadable storage, and rechecks storage on command invocation', async () => {
    const mounted = mount(); const oldPreview = report!.commands.preview!;
    vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => { throw new Error('private storage failure'); });
    await act(async () => { await oldPreview(); }); expect(api.schemaPreview).not.toHaveBeenCalled();
    mounted.rerenderPanel({ workspaceRevision: 'workspace-rev-2' });
    expect(report?.observation).toMatchObject({ state: 'unavailable', reason: 'unreadable' });
    expect(report?.commands.preview).toBeUndefined(); expect(report?.commands.confirm).toBeUndefined();
    expect(JSON.stringify(report)).not.toContain('private storage failure');
  });
  it('keeps offline inspection available while network commands remain absent', () => {
    retainSchemaOperation('workspace-1', group.id, 'O1');
    mount({ guided: { online: false, suppressDefaultActions: true, onPresentation } });
    expect(screen.getByTestId('group-schema-panel')).toBeInTheDocument();
    expect(report?.observation).toMatchObject({ state: 'unresolved', operationId: 'O1' });
    expect(report?.blocker).toBe('offline'); expect(Object.values(report!.commands).every(command => command === undefined)).toBe(true);
    expect(api.schemaOperation).not.toHaveBeenCalled(); expect(api.schemaPreview).not.toHaveBeenCalled();
  });

  it('does not let optional visible guided buttons bypass offline or unreadable-storage guards', () => {
    vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => { throw new Error('unreadable'); });
    mount({ guided: { online: true, onPresentation } });
    const previewButton = screen.getByTestId('group-schema-preview');
    expect(previewButton).toBeDisabled(); fireEvent.click(previewButton);
    expect(api.schemaPreview).not.toHaveBeenCalled();
  });
  it('never reports the previous workspace operation under a new workspace before scope cleanup', () => {
    retainSchemaOperation('workspace-1', group.id, 'O1');
    const mounted = mount();
    mounted.rerenderPanel({ workspaceId: 'workspace-2', group: { ...group, workspace_id: 'workspace-2' } });
    const newScopeReports = reports.filter(value => value?.scope.workspaceId === 'workspace-2');
    expect(newScopeReports.length).toBeGreaterThan(0);
    expect(newScopeReports.every(value => !(value && 'operationId' in value.observation && value.observation.operationId === 'O1'))).toBe(true);
    expect(report?.observation).toMatchObject({ state: 'observed', result: 'none' });
  });
  it('checks the current cross-editor O1 instead of sending a stale new-preview command', async () => {
    const mounted = mount(); const oldPreview = report!.commands.preview!;
    retainSchemaOperation('workspace-1', group.id, 'O1');
    await act(async () => { await oldPreview(); }); expect(api.schemaPreview).not.toHaveBeenCalled();
    mounted.rerenderPanel({});
    expect(report?.observation).toMatchObject({ state: 'unresolved', operationId: 'O1' });
    await act(async () => { await report?.commands.check?.(); });
    expect(api.schemaOperation).toHaveBeenCalledWith('O1');
  });

});
