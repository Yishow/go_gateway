import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { WriteGroupSection } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { WriteGroup, WriteGroupReadiness } from '@/types/studioV2WriteGroup';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string, options?: Record<string, unknown>) => (options && 'count' in options ? `${key}:${options.count}` : key) }),
}));
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({
  useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }),
}));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: {
    list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(),
    testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), get: vi.fn(),
  },
}));
const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);

const list = (groups: WriteGroup[], revision = 'wrev-1') => ({ workspace_id: 'ws-1', workspace_revision: revision, groups });
const readiness = (over: Partial<WriteGroupReadiness> = {}): WriteGroupReadiness => ({
  workspace_id: 'ws-1', workspace_revision: 'wrev-1', group_id: 'group-1', group_revision: 'rev-1', applied_revision: '',
  config_ready: true, schema_ready: true, ready: true, issues: [], ...over,
});
const typedError = (code: string, status = 409) => Object.assign(new Error(code), {
  response: { status, data: { success: false, error: { code, message: 'x', retryable: false, request_id: 'req-1', action: 'reload' } } },
});

function mount(props: Partial<React.ComponentProps<typeof WriteGroupSection>> = {}) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  return render(<QueryClientProvider client={client}><WriteGroupSection state={savedGroupState()} workspaceId="ws-1" readonly={false} {...props} /></QueryClientProvider>);
}

async function openSaved() {
  fireEvent.click(await screen.findByTestId('group-open-group-1'));
  await screen.findByTestId('group-editor');
}

beforeEach(() => {
  vi.clearAllMocks();
  api.list.mockResolvedValue(list([savedGroup()]));
  api.readiness.mockResolvedValue(readiness());
  api.delivery.mockResolvedValue({
    group_id: 'group-1', intake: { state: 'active' },
    stages: { collecting: 1, queued: 2, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 4, skipped: 0 },
    last_sql_committed_at: '2026-10-02T00:00:00Z', oldest_pending_seconds: 3, no_data_buckets: 0, skipped_buckets: 0, backlog: [],
    quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
  });
});

describe('RevisionSafeGroupActivation: Apply', () => {
  it('follows the server readiness of the saved revision', async () => {
    mount();
    await openSaved();
    await waitFor(() => expect(screen.getByTestId('group-apply')).not.toBeDisabled());
    expect(screen.getByTestId('group-readiness')).toBeInTheDocument();
  });

  it('stays off when the server says the group is not ready, and lists its issues', async () => {
    api.readiness.mockResolvedValue(readiness({ ready: false, schema_ready: false, issues: [{ code: 'table-missing', severity: 'blocking', step: 'Step 4', scope: 'group-1', message: 'The table is missing' }] }));
    mount();
    await openSaved();
    expect(await screen.findByTestId('group-readiness-issues')).toHaveTextContent('The table is missing');
    expect(screen.getByTestId('group-apply')).toBeDisabled();
    expect(screen.getByTestId('group-apply-blocked')).toHaveTextContent('apply_blocked_not_ready');
  });

  it('stays off when readiness cannot be read; an unknown answer is not a yes', async () => {
    api.readiness.mockRejectedValue(new Error('down'));
    mount();
    await openSaved();
    expect(await screen.findByTestId('group-readiness-failed')).toBeInTheDocument();
    expect(screen.getByTestId('group-apply')).toBeDisabled();
    fireEvent.click(screen.getByTestId('group-readiness-retry'));
    await waitFor(() => expect(api.readiness.mock.calls.length).toBeGreaterThan(1));
  });

  it('does not offer Apply for readiness of an older revision, nor with unsaved edits', async () => {
    api.readiness.mockResolvedValue(readiness({ group_revision: 'rev-0' }));
    mount();
    await openSaved();
    await screen.findByTestId('group-readiness');
    expect(screen.getByTestId('group-apply')).toBeDisabled();
    api.readiness.mockResolvedValue(readiness());
  });

  it('turns Apply off as soon as the draft has unsaved changes', async () => {
    mount();
    await openSaved();
    await waitFor(() => expect(screen.getByTestId('group-apply')).not.toBeDisabled());
    fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Edited' } });
    expect(screen.getByTestId('group-apply')).toBeDisabled();
    expect(screen.getByTestId('group-apply-blocked')).toHaveTextContent('apply_blocked_dirty');
  });

  it('applies once per click with the revisions it was shown, and adopts what the server returns', async () => {
    let release: (value: unknown) => void = () => undefined;
    api.apply.mockReturnValue(new Promise((resolve) => { release = resolve; }) as never);
    mount();
    await openSaved();
    await waitFor(() => expect(screen.getByTestId('group-apply')).not.toBeDisabled());
    const apply = screen.getByTestId('group-apply');
    fireEvent.click(apply);
    fireEvent.click(apply);
    await waitFor(() => expect(api.apply).toHaveBeenCalledTimes(1));
    expect(api.apply).toHaveBeenCalledTimes(1);
    expect(api.apply).toHaveBeenCalledWith('group-1', {
      workspace_id: 'ws-1', expected_workspace_revision: 'wrev-1', expected_group_revision: 'rev-1', expected_connector_revision: 'crev-1',
    });
    api.list.mockResolvedValue(list([savedGroup({ applied_revision: 'rev-1', status: 'ready' })], 'wrev-2'));
    release({ workspace_revision: 'wrev-2', group: savedGroup({ applied_revision: 'rev-1' }) });
    await waitFor(() => expect(screen.getByTestId('group-applied-revision')).toHaveTextContent('rev-1'));
  });

  it('shows a stale apply as a conflict and keeps the editor usable', async () => {
    api.apply.mockRejectedValue(typedError('revision_mismatch'));
    mount();
    await openSaved();
    await waitFor(() => expect(screen.getByTestId('group-apply')).not.toBeDisabled());
    fireEvent.click(screen.getByTestId('group-apply'));
    expect(await screen.findByTestId('group-conflict')).toBeInTheDocument();
    expect(screen.getByTestId('group-lifecycle-error')).toBeInTheDocument();
  });

  it('is blocked for a read-only session even when the server says ready', async () => {
    mount({ readonly: true });
    await openSaved();
    await screen.findByTestId('group-readiness');
    expect(screen.getByTestId('group-apply')).toBeDisabled();
  });
});

describe('GroupLifecycleWithBacklog', () => {
  it('shows what stays behind before deleting, and deletes only after confirmation', async () => {
    api.remove.mockResolvedValue({ workspace_revision: 'wrev-2', group: savedGroup({ status: 'deleted' }) });
    mount();
    await openSaved();
    fireEvent.click(screen.getByTestId('group-delete'));
    expect(await screen.findByTestId('group-delete-impact')).toHaveTextContent('delete_impact:3');
    expect(api.remove).not.toHaveBeenCalled();
    fireEvent.click(screen.getByTestId('group-delete-confirm-button'));
    await waitFor(() => expect(api.remove).toHaveBeenCalledWith('group-1', expect.objectContaining({ expected_group_revision: 'rev-1' })));
  });

  it('says so when the backlog cannot be read, and still lets the operator cancel', async () => {
    api.delivery.mockRejectedValue(new Error('down'));
    mount();
    await openSaved();
    fireEvent.click(screen.getByTestId('group-delete'));
    expect(await screen.findByTestId('group-delete-impact')).toHaveTextContent('delete_impact_unknown');
    fireEvent.click(screen.getByTestId('group-delete-cancel'));
    expect(screen.queryByTestId('group-delete-confirm')).not.toBeInTheDocument();
    expect(api.remove).not.toHaveBeenCalled();
  });

  it('disables with the saved revisions and keeps the group listed', async () => {
    api.disable.mockResolvedValue({ workspace_revision: 'wrev-2', group: savedGroup({ status: 'disabled' }) });
    mount();
    await openSaved();
    fireEvent.click(screen.getByTestId('group-disable'));
    await waitFor(() => expect(api.disable).toHaveBeenCalledWith('group-1', expect.objectContaining({ expected_group_revision: 'rev-1', expected_connector_revision: 'crev-1' })));
  });
});

describe('RevisionSafeGroupActivationAndTestWrite: test write', () => {
  const preview = {
    token: 'tok-1', operation_id: 'op-1', action: 'test_write' as const, group_id: 'group-1', group_revision: 'rev-1', expires_at: '2026-10-02T12:10:00Z',
    target: { connector_id: 'conn-1', dialect: 'sqlite', database: '/d.db', schema: '', table: 'readings' },
    owner_column: 'entity', owner_value: 'gw-test-1', dedupe: '', cleanup: 'remove only the owned row',
    values: [{ column: 'temperature', type: 'float64', value: '1.5' }, { column: 'entity', type: 'text', value: 'gw-test-1' }],
  };
  const operation = (over: Record<string, unknown> = {}) => ({
    operation_id: 'op-1', action: 'test_write' as const, status: 'succeeded' as const, write_outcome: 'written_verified' as const,
    cleanup_status: 'cleaned' as const, created_at: 'a', updated_at: 'b', completed_at: 'c', ...over,
  });

  it('previews first, writes nothing until confirmed, then shows write and cleanup separately', async () => {
    api.testWritePreview.mockResolvedValue(preview);
    api.testWrite.mockResolvedValue(operation({ status: 'partial', write_outcome: 'written_unverified', reason: 'readback-denied', cleanup_status: 'failed', cleanup_reason: 'cleanup-denied' }));
    mount();
    await openSaved();
    fireEvent.click(screen.getByTestId('group-test-write-preview'));
    expect(await screen.findByTestId('group-test-write-preview-card')).toHaveTextContent('gw-test-1');
    expect(api.testWrite).not.toHaveBeenCalled();
    fireEvent.click(screen.getByTestId('group-test-write-confirm'));
    expect(await screen.findByTestId('group-test-write-outcome')).toHaveTextContent('outcome_written_unverified');
    expect(screen.getByTestId('group-test-write-outcome')).toHaveTextContent('readback-denied');
    expect(screen.getByTestId('group-test-write-cleanup')).toHaveTextContent('cleanup_failed');
    expect(api.testWrite).toHaveBeenCalledWith('group-1', { token: 'tok-1', operation_id: 'op-1' });
  });

  it('confirms once on a double click', async () => {
    api.testWritePreview.mockResolvedValue(preview);
    let release: (value: unknown) => void = () => undefined;
    api.testWrite.mockReturnValue(new Promise((resolve) => { release = resolve; }) as never);
    mount();
    await openSaved();
    fireEvent.click(screen.getByTestId('group-test-write-preview'));
    const confirm = await screen.findByTestId('group-test-write-confirm');
    fireEvent.click(confirm);
    fireEvent.click(confirm);
    await waitFor(() => expect(api.testWrite).toHaveBeenCalledTimes(1));
    expect(api.testWrite).toHaveBeenCalledTimes(1);
    release(operation());
    expect(await screen.findByTestId('group-test-write-outcome')).toHaveTextContent('outcome_written_verified');
  });

  it('after a lost reply it checks the operation and never writes again on its own', async () => {
    api.testWritePreview.mockResolvedValue(preview);
    api.testWrite.mockRejectedValue(new Error('Network Error'));
    api.testWriteOperation.mockResolvedValue(operation());
    mount();
    await openSaved();
    fireEvent.click(screen.getByTestId('group-test-write-preview'));
    fireEvent.click(await screen.findByTestId('group-test-write-confirm'));
    expect(await screen.findByTestId('group-test-write-unconfirmed')).toBeInTheDocument();
    expect(screen.getByTestId('group-test-write-confirm')).toBeDisabled();
    fireEvent.click(screen.getByTestId('group-test-write-check-unconfirmed'));
    expect(await screen.findByTestId('group-test-write-outcome')).toHaveTextContent('outcome_written_verified');
    expect(api.testWrite).toHaveBeenCalledTimes(1);
    expect(api.testWriteOperation).toHaveBeenCalledWith('op-1');
  });

  it('treats a typed refusal as a definite failure, not as an unconfirmed write', async () => {
    api.testWritePreview.mockResolvedValue(preview);
    api.testWrite.mockRejectedValue(typedError('WRITE_GROUP_TEST_WRITE_PREVIEW_STALE'));
    mount();
    await openSaved();
    fireEvent.click(screen.getByTestId('group-test-write-preview'));
    fireEvent.click(await screen.findByTestId('group-test-write-confirm'));
    expect(await screen.findByTestId('group-test-write-error')).toBeInTheDocument();
    expect(screen.queryByTestId('group-test-write-unconfirmed')).not.toBeInTheDocument();
  });

  it('is unavailable while the group has unsaved edits', async () => {
    mount();
    await openSaved();
    fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Edited' } });
    expect(screen.getByTestId('group-test-write-preview')).toBeDisabled();
    expect(screen.getByTestId('group-test-write-blocked')).toHaveTextContent('blocked_dirty');
  });

  it('shows a running operation as running, not as a result', async () => {
    api.testWritePreview.mockResolvedValue(preview);
    api.testWrite.mockResolvedValue({ operation_id: 'op-1', action: 'test_write', status: 'running', created_at: 'a', updated_at: 'b' });
    mount();
    await openSaved();
    fireEvent.click(screen.getByTestId('group-test-write-preview'));
    fireEvent.click(await screen.findByTestId('group-test-write-confirm'));
    expect(await screen.findByTestId('group-test-write-running')).toBeInTheDocument();
    expect(screen.queryByTestId('group-test-write-outcome')).not.toBeInTheDocument();
  });

  it('drops a reply that arrives after the group changed', async () => {
    api.list.mockResolvedValue(list([savedGroup(), savedGroup({ id: 'group-2', name: 'Other' })]));
    api.testWritePreview.mockResolvedValue(preview);
    let release: (value: unknown) => void = () => undefined;
    api.testWrite.mockReturnValue(new Promise((resolve) => { release = resolve; }) as never);
    mount();
    await openSaved();
    fireEvent.click(screen.getByTestId('group-test-write-preview'));
    fireEvent.click(await screen.findByTestId('group-test-write-confirm'));
    await waitFor(() => expect(api.testWrite).toHaveBeenCalledTimes(1));
    fireEvent.click(screen.getByTestId('group-editor-close'));
    fireEvent.click(await screen.findByTestId('group-open-group-2'));
    await screen.findByTestId('group-editor');
    release(operation());
    await new Promise((resolve) => setTimeout(resolve, 20));
    expect(screen.queryByTestId('group-test-write-result')).not.toBeInTheDocument();
  });

  it('drops a reply for an earlier revision of the same group after it was saved again', async () => {
    api.testWritePreview.mockResolvedValue(preview);
    let release: (value: unknown) => void = () => undefined;
    api.testWrite.mockReturnValue(new Promise((resolve) => { release = resolve; }) as never);
    api.update.mockResolvedValue({ workspace_revision: 'wrev-2', group: savedGroup({ revision: 'rev-2', name: 'Edited' }) });
    mount();
    await openSaved();
    fireEvent.click(screen.getByTestId('group-test-write-preview'));
    fireEvent.click(await screen.findByTestId('group-test-write-confirm'));
    await waitFor(() => expect(api.testWrite).toHaveBeenCalledTimes(1));
    fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Edited' } });
    api.list.mockResolvedValue(list([savedGroup({ revision: 'rev-2', name: 'Edited' })], 'wrev-2'));
    fireEvent.click(screen.getByTestId('group-save'));
    await waitFor(() => expect(api.update).toHaveBeenCalledTimes(1));
    await waitFor(() => expect(screen.queryByTestId('group-dirty-note')).not.toBeInTheDocument());
    release(operation());
    await new Promise((resolve) => setTimeout(resolve, 20));
    expect(screen.queryByTestId('group-test-write-result')).not.toBeInTheDocument();
    expect(screen.queryByTestId('group-test-write-preview-card')).not.toBeInTheDocument();
  });
});
