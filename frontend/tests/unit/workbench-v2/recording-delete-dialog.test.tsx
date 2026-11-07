import { act, cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { SavedRecordingRemoveDialog } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/lifecycle/SavedRecordingRemoveDialog';
import { savedGroup } from '../../fixtures/writeGroupState';
import { DeleteResponseError } from '@/services/studioV2WorkspaceDelete';
import { RecordingDeleteDialog } from '@/features/datalink/workbench-v2/recording-home/RecordingDeleteDialog';
import { GroupLifecycleBar } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WriteGroupLifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
import { savedLifecycle } from '../../fixtures/writeGroupLifecycle';
import { waitFor } from '@testing-library/react';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
afterEach(() => { cleanup(); vi.restoreAllMocks(); });

it('requires acknowledgement before deleting the recording', () => {
  const confirm = vi.fn();
  render(<SavedRecordingRemoveDialog open busy={false} group={savedGroup()} onClose={vi.fn()} onConfirm={confirm} />);
  expect(screen.getByTestId('group-delete-confirm-button')).toBeDisabled();
  expect(screen.getByText('deleteDialogs.recording.cannotUndo')).toBeInTheDocument();
  fireEvent.click(screen.getByRole('checkbox', { name: 'deleteDialogs.recording.acknowledge' }));
  fireEvent.click(screen.getByTestId('group-delete-confirm-button'));
  expect(confirm).toHaveBeenCalledOnce();
});

it('does not carry acknowledgement to another recording revision', () => {
  const group = savedGroup();
  const renderDialog = (revision: string) => <SavedRecordingRemoveDialog open busy={false} group={{ ...group, revision }} onClose={vi.fn()} onConfirm={vi.fn()} />;
  const view = render(renderDialog('R1'));
  fireEvent.click(screen.getByRole('checkbox', { name: 'deleteDialogs.recording.acknowledge' }));
  expect(screen.getByTestId('group-delete-confirm-button')).toBeEnabled();
  view.rerender(renderDialog('R2'));
  expect(screen.getByTestId('group-delete-confirm-button')).toBeDisabled();
});

it('explains a backlog without adding samples and outbox rows, and offers stop and refresh', () => {
  const stop = vi.fn(); const refresh = vi.fn();
  const error = new DeleteResponseError('failed', 'WRITE_GROUP_DELETE_BACKLOG', 'stop_and_wait', {
    samples_unconsumed: 3, outbox: { pending: 2, unknown: 1 },
  });
  render(<SavedRecordingRemoveDialog open busy={false} group={savedGroup()} error={error} onStop={stop} onRefresh={refresh} onClose={vi.fn()} onConfirm={vi.fn()} />);
  expect(screen.getByTestId('delete-backlog-samples')).toHaveTextContent('3');
  expect(screen.getByTestId('delete-backlog-pending')).toHaveTextContent('2');
  expect(screen.getByTestId('delete-backlog-unknown')).toHaveTextContent('1');
  fireEvent.click(screen.getByRole('button', { name: 'deleteDialogs.recording.stop' }));
  fireEvent.click(screen.getByRole('button', { name: 'deleteDialogs.recording.refresh' }));
  expect(stop).toHaveBeenCalledOnce(); expect(refresh).toHaveBeenCalledOnce();
});

it('names dependent recordings when a successor conflict includes them', () => {
  const error = new DeleteResponseError('failed', 'WRITE_GROUP_DELETE_HAS_SUCCESSOR', 'open_recording', {
    referenced_by: [{ id: 'G3', name: 'Later recording' }], total: 2,
  });
  render(<SavedRecordingRemoveDialog open busy={false} group={savedGroup()} error={error} onOpenRecording={vi.fn()} onRefresh={vi.fn()} onClose={vi.fn()} onConfirm={vi.fn()} />);
  expect(screen.getByText('Later recording')).toBeInTheDocument();
  expect(screen.getByTestId('delete-successor-total')).toHaveTextContent('2');
  expect(screen.getByRole('button', { name: 'deleteDialogs.recording.openRecording' })).toBeEnabled();
});

it('submits the saved revisions after a fresh authority read and reports a deletion receipt', async () => {
  const group = savedGroup();
  const lifecycle = vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockResolvedValue(savedLifecycle(group));
  const remove = vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'remove').mockResolvedValueOnce({ deleted: true, group_id: group.id, lineage_id: 'L1', deleted_generations: 2 });
  const onDeleted = vi.fn(); const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><RecordingDeleteDialog open group={group} workspaceRevision="wrev-1" onClose={vi.fn()} onDeleted={onDeleted}
    onStop={vi.fn()} onOpenRecording={vi.fn()} onRefresh={vi.fn()} /></QueryClientProvider>);
  fireEvent.click(screen.getByRole('checkbox', { name: 'deleteDialogs.recording.acknowledge' }));
  fireEvent.click(screen.getByTestId('group-delete-confirm-button'));
  await waitFor(() => expect(onDeleted).toHaveBeenCalledOnce());
  expect(lifecycle).toHaveBeenCalledWith(group.id);
  expect(remove).toHaveBeenCalledWith(group.id, { workspace_id: 'ws-1', expected_workspace_revision: 'wrev-1', expected_group_revision: 'rev-1', expected_connector_revision: 'crev-1' });
  lifecycle.mockRestore(); remove.mockRestore(); client.clear();
});

it.each([
  ['WRITE_GROUP_DELETE_ACTIVE', 'stop_first', 'active', 'stop'],
  ['WRITE_GROUP_DELETE_OPERATION_PENDING', 'open_recording', 'operationPending', 'openRecording'],
  ['WRITE_GROUP_DELETE_HAS_SUCCESSOR', 'open_recording', 'hasSuccessor', 'openRecording'],
])('keeps %s open with the matching recovery action', async (code, action, message, button) => {
  const group = savedGroup();
  vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockResolvedValue(savedLifecycle(group));
  vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'remove').mockRejectedValueOnce(new DeleteResponseError('failed', code, action));
  const onDeleted = vi.fn(); const onStop = vi.fn(); const onOpenRecording = vi.fn();
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><RecordingDeleteDialog open group={group} workspaceRevision="wrev-1" onClose={vi.fn()} onDeleted={onDeleted}
    onStop={onStop} onOpenRecording={onOpenRecording} onRefresh={vi.fn()} /></QueryClientProvider>);
  fireEvent.click(screen.getByRole('checkbox', { name: 'deleteDialogs.recording.acknowledge' }));
  fireEvent.click(screen.getByTestId('group-delete-confirm-button'));
  expect(await screen.findByText(`deleteDialogs.recording.${message}`)).toBeInTheDocument();
  fireEvent.click(screen.getByRole('button', { name: `deleteDialogs.recording.${button}` }));
  expect(button === 'stop' ? onStop : onOpenRecording).toHaveBeenCalledOnce();
  expect(onDeleted).not.toHaveBeenCalled(); client.clear();
});

it('shows an unconfirmed deletion and refresh after a network failure', async () => {
  const group = savedGroup();
  vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockResolvedValue(savedLifecycle(group));
  vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'remove').mockRejectedValueOnce(new Error('private network text'));
  const onRefresh = vi.fn(); const onDeleted = vi.fn();
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><RecordingDeleteDialog open group={group} workspaceRevision="wrev-1" onClose={vi.fn()} onDeleted={onDeleted}
    onStop={vi.fn()} onOpenRecording={vi.fn()} onRefresh={onRefresh} /></QueryClientProvider>);
  fireEvent.click(screen.getByRole('checkbox', { name: 'deleteDialogs.recording.acknowledge' }));
  fireEvent.click(screen.getByTestId('group-delete-confirm-button'));
  expect(await screen.findByText('deleteDialogs.recording.unknown')).toBeInTheDocument();
  expect(screen.queryByText('private network text')).not.toBeInTheDocument();
  fireEvent.click(screen.getByRole('button', { name: 'deleteDialogs.recording.refresh' }));
  expect(onRefresh).toHaveBeenCalledOnce(); expect(onDeleted).not.toHaveBeenCalled(); client.clear();
});

it('states that an unfamiliar conflict refused deletion and offers refresh', () => {
  const refresh = vi.fn();
  render(<SavedRecordingRemoveDialog open busy={false} group={savedGroup()}
    error={new DeleteResponseError('failed', 'WRITE_GROUP_LIFECYCLE_BLOCKED', 'reload')}
    onRefresh={refresh} onClose={vi.fn()} onConfirm={vi.fn()} />);
  expect(screen.getByText('deleteDialogs.recording.refused')).toBeInTheDocument();
  expect(screen.queryByText('deleteDialogs.recording.unknown')).not.toBeInTheDocument();
  fireEvent.click(screen.getByRole('button', { name: 'deleteDialogs.recording.refresh' }));
  expect(refresh).toHaveBeenCalledOnce();
});

it('does not send deletion when the selected revision changes while authority is loading', async () => {
  const group = savedGroup();
  let release!: (value: ReturnType<typeof savedLifecycle>) => void;
  vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockImplementationOnce(() => new Promise(resolve => { release = resolve; }));
  const remove = vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'remove');
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const dialog = (revision: string) => <QueryClientProvider client={client}><RecordingDeleteDialog open group={{ ...group, revision }} workspaceRevision="wrev-1" onClose={vi.fn()} onDeleted={vi.fn()}
    onStop={vi.fn()} onOpenRecording={vi.fn()} onRefresh={vi.fn()} /></QueryClientProvider>;
  const view = render(dialog('rev-1'));
  fireEvent.click(screen.getByRole('checkbox', { name: 'deleteDialogs.recording.acknowledge' }));
  fireEvent.click(screen.getByTestId('group-delete-confirm-button'));
  view.rerender(dialog('rev-2'));
  await act(async () => { release(savedLifecycle(group)); });
  expect(remove).not.toHaveBeenCalled(); client.clear();
});

it('blocks deletion if refreshed lifecycle facts have a different connector revision', async () => {
  const group = savedGroup();
  vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockResolvedValue({ ...savedLifecycle(group), connector_revision: 'changed' });
  const remove = vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'remove');
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><RecordingDeleteDialog open group={group} workspaceRevision="wrev-1" onClose={vi.fn()} onDeleted={vi.fn()}
    onStop={vi.fn()} onOpenRecording={vi.fn()} onRefresh={vi.fn()} /></QueryClientProvider>);
  fireEvent.click(screen.getByRole('checkbox', { name: 'deleteDialogs.recording.acknowledge' }));
  fireEvent.click(screen.getByTestId('group-delete-confirm-button'));
  expect(await screen.findByText('deleteDialogs.recording.stale')).toBeInTheDocument();
  expect(remove).not.toHaveBeenCalled(); client.clear();
});

it('editor refuses to send deletion when refreshed connector facts no longer match', async () => {
  const group = savedGroup(); const facts = savedLifecycle(group);
  const lifecycle = vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockResolvedValue(facts);
  const remove = vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'remove');
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><GroupLifecycleBar group={group} workspaceId="ws-1" workspaceRevision="wrev-1"
    readonly={false} canApply onChanged={vi.fn()} onConflict={vi.fn()} /></QueryClientProvider>);
  await waitFor(() => expect(screen.getByTestId('group-delete')).toBeEnabled());
  fireEvent.click(screen.getByTestId('group-delete'));
  fireEvent.click(await screen.findByRole('checkbox', { name: 'deleteDialogs.recording.acknowledge' }));
  lifecycle.mockResolvedValue({ ...facts, connector_revision: 'changed' });
  fireEvent.click(screen.getByTestId('group-delete-confirm-button'));
  expect(await screen.findByText('reason.stale_scope')).toBeInTheDocument();
  expect(remove).not.toHaveBeenCalled(); client.clear();
});

it('keeps the editor delete action disabled when its initial connector facts are stale', async () => {
  const group = savedGroup();
  const lifecycle = vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockResolvedValue({ ...savedLifecycle(group), connector_revision: 'changed' });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><GroupLifecycleBar group={group} workspaceId="ws-1" workspaceRevision="wrev-1"
    readonly={false} canApply onChanged={vi.fn()} onConflict={vi.fn()} /></QueryClientProvider>);
  await waitFor(() => expect(lifecycle).toHaveBeenCalledOnce());
  await waitFor(() => expect(client.isFetching()).toBe(0));
  expect(screen.getByTestId('saved-recording-refresh')).toBeEnabled();
  expect(screen.getByTestId('group-delete')).toBeDisabled(); client.clear();
});
