import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { readSchemaOperation, retainSchemaOperation } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/schemaOperationRecovery';
import { managedSchemaPreviewToken } from '../../fixtures/managedSchemaPreview';
import { GroupSchemaPanel } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { SchemaOperation } from '@/types/recordingPlan';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string) => key }),
}));

vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: {
    schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(),
  },
}));

const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const group: WriteGroup = {
  id: 'group-1', workspace_id: 'workspace-1', revision: 'group-rev-1', applied_revision: '',
  name: 'Line A', status: 'ready', members: [],
  destination: {
    connector_id: 'connector-1', connector_revision: 'connector-rev-1', database: 'recording.db',
    table_schema: 'main', table_name: 'gw_group_abc', storage_strategy: 'managed',
  },
  row_policy: { interval_seconds: 60, allowed_lateness_seconds: 0, record_key_column: 'record_id' },
  write_policy: { mode: 'append', dedupe_capability: 'receipt' }, migration: {},
  created_at: '2026-10-04T00:00:00Z', updated_at: '2026-10-04T00:00:00Z',
};

const token = managedSchemaPreviewToken;

const succeeded = {
  operation_id: 'operation-1', action: 'schema_apply', status: 'succeeded' as const, executed_statements: 2,
  verified_digest: 'verified-digest', created_at: '2026-10-04T00:00:00Z', updated_at: '2026-10-04T00:00:01Z',
};

function renderPanel(overrides: Partial<React.ComponentProps<typeof GroupSchemaPanel>> = {}) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <GroupSchemaPanel
        group={group} workspaceId="workspace-1" workspaceRevision="workspace-rev-1"
        destinationMatches dirty={false} readonly={false} onApplied={vi.fn()} {...overrides}
      />
    </QueryClientProvider>,
  );
}

beforeEach(() => {
  vi.clearAllMocks();
  window.sessionStorage.clear();
  retainSchemaOperation('workspace-1', 'group-1', undefined, readSchemaOperation('workspace-1', 'group-1') ?? '');
  retainSchemaOperation('workspace-1', 'group-2', undefined, readSchemaOperation('workspace-1', 'group-2') ?? '');
  api.schemaPreview.mockResolvedValue(token);
  api.schemaApply.mockResolvedValue(succeeded);
});

describe('GroupSchemaPanel', () => {
  it('previews the real server layout, then confirms once and reports the operation', async () => {
    const onApplied = vi.fn();
    renderPanel({ onApplied });
    fireEvent.click(screen.getByTestId('group-schema-preview'));
    expect(await screen.findByTestId('group-schema-preview-table')).toHaveTextContent('connector_default_only');
    expect(screen.getByTestId('group-schema-columns')).toHaveTextContent('record_id');
    expect(screen.queryByText('_gw_owner_efd49c6d03ab3796062ddecc')).not.toBeInTheDocument();
    fireEvent.click(screen.getByTestId('group-schema-apply'));
    await waitFor(() => expect(screen.getByTestId('group-schema-operation')).toHaveTextContent('succeeded'));
    expect(onApplied).toHaveBeenCalledTimes(1);
    expect(screen.getByTestId('group-schema-operation')).not.toHaveTextContent('step4.group.schema.operation_attention');
    expect(api.schemaPreview).toHaveBeenCalledWith('group-1', {
      workspace_id: 'workspace-1', expected_workspace_revision: 'workspace-rev-1',
      expected_group_revision: 'group-rev-1', expected_connector_revision: 'connector-rev-1',
    });
    expect(api.schemaApply).toHaveBeenCalledWith('group-1', {
      workspace_id: 'workspace-1', expected_workspace_revision: 'workspace-rev-1',
      expected_group_revision: 'group-rev-1', expected_connector_revision: 'connector-rev-1',
      token: 'token-1', operation_id: 'operation-1',
    });
  });

  it.each([
    ['unsaved', { group: undefined, dirty: true }, 'step4.group.schema.blocked_unsaved'],
    ['dirty', { dirty: true }, 'step4.group.schema.blocked_dirty'],
    ['destination', { destinationMatches: false }, 'step4.group.schema.blocked_destination'],
    ['readonly', { readonly: true }, 'step4.group.schema.blocked_readonly'],
  ])('blocks %s without calling the schema API', async (_name, props, message) => {
    renderPanel(props);
    expect(screen.getByTestId('group-schema-blocked')).toHaveTextContent(message);
    expect(screen.getByTestId('group-schema-preview')).toBeDisabled();
    expect(api.schemaPreview).not.toHaveBeenCalled();
  });

  it('keeps a failed operation on the same operation lookup path and never retries DDL', async () => {
    api.schemaApply.mockResolvedValue({ ...succeeded, status: 'failed', reason: 'permission_denied' });
    api.schemaOperation.mockResolvedValue({ ...succeeded, status: 'failed', reason: 'permission_denied' });
    renderPanel();
    fireEvent.click(screen.getByTestId('group-schema-preview'));
    fireEvent.click(await screen.findByTestId('group-schema-apply'));
    await waitFor(() => expect(screen.getByTestId('group-schema-operation')).toHaveTextContent('failed'));
    expect(screen.getByTestId('group-schema-operation')).toHaveTextContent('errors.RECORDING_SCHEMA_PERMISSION_DENIED');
    expect(api.schemaApply).toHaveBeenCalledTimes(1);
    fireEvent.click(screen.getByTestId('group-schema-check'));
    await waitFor(() => expect(api.schemaOperation).toHaveBeenCalledWith('operation-1'));
    expect(api.schemaApply).toHaveBeenCalledTimes(1);
  });

  it('shows an actionable safe error when same-operation lookup fails', async () => {
    api.schemaApply.mockResolvedValue({ ...succeeded, status: 'failed', reason: 'permission_denied' });
    api.schemaOperation.mockRejectedValue(new Error('raw backend exception with secret details'));
    renderPanel();
    fireEvent.click(screen.getByTestId('group-schema-preview'));
    fireEvent.click(await screen.findByTestId('group-schema-apply'));
    await waitFor(() => expect(screen.getByTestId('group-schema-operation')).toHaveTextContent('failed'));

    fireEvent.click(screen.getByTestId('group-schema-check'));
    await waitFor(() => expect(screen.getByTestId('group-schema-error')).toHaveTextContent('step4.group.schema.error_check_same_operation'));
    expect(screen.getByTestId('group-schema-error')).toHaveTextContent('errors.generic_failure');
    expect(screen.getByTestId('group-schema-error')).not.toHaveTextContent('raw backend exception');
    expect(api.schemaOperation).toHaveBeenCalledWith('operation-1');
    expect(api.schemaApply).toHaveBeenCalledTimes(1);
  });

  it('single-flights explicit confirmation while the server reply is pending', async () => {
    let release: (value: typeof succeeded) => void = () => undefined;
    api.schemaApply.mockReturnValue(new Promise((resolve) => { release = resolve; }));
    renderPanel();
    fireEvent.click(screen.getByTestId('group-schema-preview'));
    fireEvent.click(await screen.findByTestId('group-schema-apply'));
    fireEvent.click(screen.getByTestId('group-schema-apply'));
    await waitFor(() => expect(api.schemaApply).toHaveBeenCalledTimes(1));
    release(succeeded);
    await waitFor(() => expect(screen.getByTestId('group-schema-operation')).toHaveTextContent('succeeded'));
  });

  it('retains the terminal result when success invalidation refreshes the proof scope', async () => {
    const ScopeRefreshHarness = () => {
      const [workspaceRevision, setWorkspaceRevision] = React.useState('workspace-rev-1');
      return (
        <GroupSchemaPanel
          group={group} workspaceId="workspace-1" workspaceRevision={workspaceRevision}
          destinationMatches dirty={false} readonly={false}
          onApplied={() => setWorkspaceRevision('workspace-rev-2')}
        />
      );
    };
    render(
      <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })}>
        <ScopeRefreshHarness />
      </QueryClientProvider>,
    );
    fireEvent.click(screen.getByTestId('group-schema-preview'));
    fireEvent.click(await screen.findByTestId('group-schema-apply'));
    await waitFor(() => expect(screen.getByTestId('group-schema-operation')).toHaveTextContent('succeeded'));
    expect(screen.getByTestId('group-schema-operation')).toHaveTextContent('operation-1');
    expect(api.schemaApply).toHaveBeenCalledTimes(1);
  });

  it.each(['before_reply', 'after_reply'])('keeps a lost response queryable when same-group scope changes %s', async (timing) => {
    let rejectApply: (reason?: unknown) => void = () => undefined;
    api.schemaApply.mockReturnValue(new Promise((_resolve, reject) => { rejectApply = reject; }));
    api.schemaOperation.mockResolvedValue({ ...succeeded, status: 'unknown' });
    const ScopeRefreshHarness = () => {
      const [workspaceRevision, setWorkspaceRevision] = React.useState('workspace-rev-1');
      return (
        <>
          <GroupSchemaPanel
            group={group} workspaceId="workspace-1" workspaceRevision={workspaceRevision}
            destinationMatches dirty={false} readonly={false}
          />
          <button type="button" data-testid="scope-change" onClick={() => setWorkspaceRevision('workspace-rev-2')}>change scope</button>
        </>
      );
    };
    render(
      <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })}>
        <ScopeRefreshHarness />
      </QueryClientProvider>,
    );
    fireEvent.click(screen.getByTestId('group-schema-preview'));
    fireEvent.click(await screen.findByTestId('group-schema-apply'));
    if (timing === 'before_reply') fireEvent.click(screen.getByTestId('scope-change'));
    rejectApply(new Error('lost apply reply'));

    await waitFor(() => expect(screen.getByTestId('group-schema-check')).toBeInTheDocument());
    if (timing === 'after_reply') fireEvent.click(screen.getByTestId('scope-change'));
    await waitFor(() => expect(screen.getByTestId('group-schema-check')).toBeInTheDocument());
    expect(screen.queryByTestId('group-schema-preview-result')).not.toBeInTheDocument();
    fireEvent.click(screen.getByTestId('group-schema-check'));
    await waitFor(() => expect(api.schemaOperation).toHaveBeenCalledWith('operation-1'));
    expect(api.schemaApply).toHaveBeenCalledTimes(1);
  });

  it('drops preview and terminal evidence when switching to another group', async () => {
    const SwitchHarness = () => {
      const [selected, setSelected] = React.useState(group);
      return (
        <>
          <GroupSchemaPanel
            group={selected} workspaceId="workspace-1" workspaceRevision="workspace-rev-1"
            destinationMatches dirty={false} readonly={false}
          />
          <button type="button" data-testid="group-change" onClick={() => setSelected({ ...group, id: 'group-2' })}>change group</button>
        </>
      );
    };
    render(
      <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })}>
        <SwitchHarness />
      </QueryClientProvider>,
    );
    fireEvent.click(screen.getByTestId('group-schema-preview'));
    expect(await screen.findByTestId('group-schema-preview-result')).toBeInTheDocument();
    fireEvent.click(screen.getByTestId('group-schema-apply'));
    await waitFor(() => expect(screen.getByTestId('group-schema-operation')).toHaveTextContent('succeeded'));
    fireEvent.click(screen.getByTestId('group-change'));
    await waitFor(() => expect(screen.queryByTestId('group-schema-preview-result')).not.toBeInTheDocument());
    expect(screen.queryByTestId('group-schema-operation')).not.toBeInTheDocument();
    expect(screen.queryByTestId('group-schema-check')).not.toBeInTheDocument();
    expect(api.schemaApply).toHaveBeenCalledTimes(1);
    expect(api.schemaOperation).not.toHaveBeenCalled();
  });

  it('keeps a failed terminal operation queryable after a same-group scope refresh', async () => {
    const failed = { ...succeeded, status: 'failed' as const, reason: 'permission_denied' };
    api.schemaApply.mockResolvedValue(failed);
    api.schemaOperation.mockResolvedValue(failed);
    const ScopeRefreshHarness = () => {
      const [workspaceRevision, setWorkspaceRevision] = React.useState('workspace-rev-1');
      return (
        <>
          <GroupSchemaPanel
            group={group} workspaceId="workspace-1" workspaceRevision={workspaceRevision}
            destinationMatches dirty={false} readonly={false}
          />
          <button type="button" data-testid="scope-change" onClick={() => setWorkspaceRevision('workspace-rev-2')}>change scope</button>
        </>
      );
    };
    render(
      <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })}>
        <ScopeRefreshHarness />
      </QueryClientProvider>,
    );
    fireEvent.click(screen.getByTestId('group-schema-preview'));
    fireEvent.click(await screen.findByTestId('group-schema-apply'));
    await waitFor(() => expect(screen.getByTestId('group-schema-operation')).toHaveTextContent('failed'));
    fireEvent.click(screen.getByTestId('scope-change'));
    fireEvent.click(await screen.findByTestId('group-schema-check'));
    await waitFor(() => expect(api.schemaOperation).toHaveBeenCalledWith('operation-1'));
    expect(api.schemaApply).toHaveBeenCalledTimes(1);
  });

  it('drops a preview when the proof scope changes even if the group revision is unchanged', async () => {
    const ScopeHarness = () => {
      const [digest, setDigest] = React.useState('');
      const scopedGroup = { ...group, destination: { ...group.destination, schema_revision: digest ? 'schema-2' : undefined, schema_digest: digest || undefined } };
      return (
        <>
          <GroupSchemaPanel
            group={scopedGroup} workspaceId="workspace-1" workspaceRevision="workspace-rev-1"
            destinationMatches dirty={false} readonly={false}
          />
          <button type="button" data-testid="scope-change" onClick={() => setDigest('digest-2')}>change scope</button>
        </>
      );
    };
    render(
      <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } })}>
        <ScopeHarness />
      </QueryClientProvider>,
    );
    fireEvent.click(screen.getByTestId('group-schema-preview'));
    await screen.findByTestId('group-schema-preview-result');
    fireEvent.click(screen.getByTestId('scope-change'));
    await waitFor(() => expect(screen.queryByTestId('group-schema-preview-result')).not.toBeInTheDocument());
    expect(api.schemaApply).not.toHaveBeenCalled();
  });
});

it('automatically performs only the read-only preview and waits for explicit DDL consent', async () => {
  renderPanel({ autoPreview: true });
  await screen.findByTestId('group-schema-preview-result');
  expect(api.schemaPreview).toHaveBeenCalledTimes(1);
  expect(api.schemaApply).not.toHaveBeenCalled();
});
it('recovers an unknown DDL operation after remount without generating a new preview', async () => {
  api.schemaApply.mockRejectedValueOnce(new Error('lost reply'));
  api.schemaOperation.mockResolvedValue({ ...succeeded, status: 'unknown' });
  const first = renderPanel({ autoPreview: true });
  fireEvent.click(await screen.findByTestId('group-schema-apply'));
  await screen.findByTestId('group-schema-check');
  first.unmount();
  renderPanel({ autoPreview: true });
  expect(screen.getByTestId('group-schema-preview')).toBeDisabled();
  fireEvent.click(screen.getByTestId('group-schema-check'));
  await waitFor(() => expect(api.schemaOperation).toHaveBeenCalledWith('operation-1'));
  expect(api.schemaPreview).toHaveBeenCalledTimes(1);
  expect(api.schemaApply).toHaveBeenCalledTimes(1);
});

it('blocks a second already-open schema editor after another editor has an unknown apply', async () => {
  api.schemaApply.mockRejectedValueOnce(new Error('lost response'));
  renderPanel(); renderPanel();
  const previews = screen.getAllByTestId('group-schema-preview');
  fireEvent.click(previews[0]); fireEvent.click(previews[1]);
  await waitFor(() => expect(screen.getAllByTestId('group-schema-apply')).toHaveLength(2));
  fireEvent.click(screen.getAllByTestId('group-schema-apply')[0]);
  await screen.findByTestId('group-schema-check');
  fireEvent.click(screen.getAllByTestId('group-schema-apply')[1]);
  expect(api.schemaApply).toHaveBeenCalledTimes(1);
  await waitFor(() => expect(screen.getAllByTestId('group-schema-check')).toHaveLength(2));
});
it('requests a new read-only preview only after another explicit preparation request', async () => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const props = { group, workspaceId: 'workspace-1', workspaceRevision: 'workspace-rev-1', destinationMatches: true, dirty: false, readonly: false, autoPreview: true, previewRequest: 1 };
  const view = render(<QueryClientProvider client={client}><GroupSchemaPanel {...props} /></QueryClientProvider>);
  await screen.findByTestId('group-schema-preview-result');
  view.rerender(<QueryClientProvider client={client}><GroupSchemaPanel {...props} workspaceRevision="workspace-rev-2" /></QueryClientProvider>);
  expect(api.schemaPreview).toHaveBeenCalledTimes(1);
  view.rerender(<QueryClientProvider client={client}><GroupSchemaPanel {...props} workspaceRevision="workspace-rev-2" previewRequest={2} /></QueryClientProvider>);
  await waitFor(() => expect(api.schemaPreview).toHaveBeenCalledTimes(2));
  expect(api.schemaApply).not.toHaveBeenCalled();
});

it('retains a verified recovered operation without an automatic replacement preview', async () => {
  retainSchemaOperation('workspace-1', group.id, 'operation-1');
  api.schemaOperation.mockResolvedValue(succeeded);
  renderPanel({ autoPreview: true });
  fireEvent.click(screen.getByTestId('group-schema-check'));
  await waitFor(() => expect(screen.getByTestId('group-schema-operation')).toHaveTextContent('succeeded'));
  expect(api.schemaPreview).not.toHaveBeenCalled();
  expect(api.schemaApply).not.toHaveBeenCalled();
});

it('permits a new read-only preview after a confirmed rolled-back failure while retaining its ledger result', async () => {
  api.schemaApply.mockResolvedValue({ ...succeeded, status: 'failed', reason: 'rolled_back' });
  renderPanel({ autoPreview: true });
  fireEvent.click(await screen.findByTestId('group-schema-apply'));
  await waitFor(() => expect(screen.getByTestId('group-schema-operation')).toHaveTextContent('failed'));
  expect(screen.getByTestId('group-schema-preview')).toBeEnabled();
  fireEvent.click(screen.getByTestId('group-schema-preview'));
  await waitFor(() => expect(api.schemaPreview).toHaveBeenCalledTimes(2));
  expect(api.schemaApply).toHaveBeenCalledTimes(1);
});

it.each(['failed', 'succeeded'] as const)('a delayed %s lookup cannot release a newer unknown operation', async (status) => {
  const terminal = { ...succeeded, status, reason: status === 'failed' ? 'rolled_back' : undefined };
  let finishOldLookup: (result: SchemaOperation) => void = () => undefined;
  api.schemaOperation
    .mockReturnValueOnce(new Promise(resolve => { finishOldLookup = resolve; }))
    .mockResolvedValueOnce(terminal);
  retainSchemaOperation('workspace-1', group.id, 'operation-1');
  const oldPanel = renderPanel();
  const resolvingPanel = renderPanel();
  fireEvent.click(within(oldPanel.container).getByTestId('group-schema-check'));
  await waitFor(() => expect(api.schemaOperation).toHaveBeenCalledTimes(1));
  fireEvent.click(within(resolvingPanel.container).getByTestId('group-schema-check'));
  await waitFor(() => expect(within(resolvingPanel.container).getByTestId('group-schema-operation')).toHaveTextContent(status));
  expect(readSchemaOperation('workspace-1', group.id)).toBeUndefined();

  api.schemaPreview
    .mockResolvedValueOnce({ ...token, operation_id: 'operation-2', token: 'token-2' })
    .mockResolvedValueOnce({ ...token, operation_id: 'operation-3', token: 'token-3' });
  const newerPanel = renderPanel({ autoPreview: true });
  const alreadyOpenPanel = renderPanel({ autoPreview: true });
  await within(newerPanel.container).findByTestId('group-schema-apply');
  await within(alreadyOpenPanel.container).findByTestId('group-schema-apply');
  api.schemaApply.mockRejectedValueOnce(new Error('lost newer response'));
  fireEvent.click(within(newerPanel.container).getByTestId('group-schema-apply'));
  await within(newerPanel.container).findByTestId('group-schema-check');
  expect(readSchemaOperation('workspace-1', group.id)).toBe('operation-2');

  await act(async () => { finishOldLookup(terminal); });
  await waitFor(() => expect(within(oldPanel.container).getByTestId('group-schema-operation')).toHaveTextContent(status));
  expect(readSchemaOperation('workspace-1', group.id)).toBe('operation-2');
  fireEvent.click(within(alreadyOpenPanel.container).getByTestId('group-schema-apply'));
  await within(alreadyOpenPanel.container).findByTestId('group-schema-check');
  expect(api.schemaApply).toHaveBeenCalledTimes(1);
  const freshPanel = renderPanel({ autoPreview: true });
  expect(within(freshPanel.container).getByTestId('group-schema-preview')).toBeDisabled();
  expect(api.schemaPreview).toHaveBeenCalledTimes(2);
});

it.each(['failed', 'succeeded'] as const)('a delayed %s apply response cannot release a newer unknown operation', async (status) => {
  const terminal = { ...succeeded, status, reason: status === 'failed' ? 'rolled_back' : undefined };
  let finishOldApply: (result: SchemaOperation) => void = () => undefined;
  api.schemaApply
    .mockReturnValueOnce(new Promise(resolve => { finishOldApply = resolve; }))
    .mockRejectedValueOnce(new Error('lost newer response'));
  api.schemaOperation.mockResolvedValueOnce(terminal);
  const oldPanel = renderPanel({ autoPreview: true });
  fireEvent.click(await within(oldPanel.container).findByTestId('group-schema-apply'));
  await waitFor(() => expect(api.schemaApply).toHaveBeenCalledTimes(1));

  const resolvingPanel = renderPanel();
  fireEvent.click(within(resolvingPanel.container).getByTestId('group-schema-check'));
  await waitFor(() => expect(within(resolvingPanel.container).getByTestId('group-schema-operation')).toHaveTextContent(status));
  expect(readSchemaOperation('workspace-1', group.id)).toBeUndefined();

  api.schemaPreview.mockResolvedValueOnce({ ...token, operation_id: 'operation-2', token: 'token-2' });
  const newerPanel = renderPanel({ autoPreview: true });
  fireEvent.click(await within(newerPanel.container).findByTestId('group-schema-apply'));
  await within(newerPanel.container).findByTestId('group-schema-check');
  expect(readSchemaOperation('workspace-1', group.id)).toBe('operation-2');
  await act(async () => { finishOldApply(terminal); });
  await waitFor(() => expect(within(oldPanel.container).getByTestId('group-schema-operation')).toHaveTextContent(status));
  expect(readSchemaOperation('workspace-1', group.id)).toBe('operation-2');
  const freshPanel = renderPanel({ autoPreview: true });
  expect(within(freshPanel.container).getByTestId('group-schema-preview')).toBeDisabled();
  expect(api.schemaPreview).toHaveBeenCalledTimes(2);
});
