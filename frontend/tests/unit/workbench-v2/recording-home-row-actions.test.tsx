import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { RecordingHomeCollection } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeCollection';
import { initialCollectionState } from '@/features/datalink/workbench-v2/recording-home/recordingHomeCollections';
import { recordingHomeGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import type { WriteGroupLifecycle, WriteGroupLifecycleAction } from '@/types/studioV2WriteGroupLifecycle';
import type { WriteGroup, WriteGroupStatus } from '@/types/studioV2WriteGroup';
import { studioV2WriteGroupLifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { fixture, group, groupScope } from './recording-home-guidance-fixtures';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
afterEach(() => { cleanup(); vi.restoreAllMocks(); });

const availability = (action: WriteGroupLifecycleAction['action'], allowed: boolean): WriteGroupLifecycleAction =>
  ({ action, allowed, reason_code: '', next_action: '', requires_review: false });

function renderRow({ stale = false, status = 'running', stop = true, remove = true }: { stale?: boolean; status?: WriteGroupStatus; stop?: boolean; remove?: boolean } = {}) {
  const item = group('R', status);
  const input = fixture(['D1'], true, [item]);
  const lifecycle = input.evidence.lifecycle?.R;
  if (!lifecycle || lifecycle.state !== 'available') throw new Error('fixture lifecycle missing');
  const guidance = recordingHomeGuidance({ ...input, evidence: { ...input.evidence, lifecycle: {
    R: { ...lifecycle, scope: groupScope(item), value: { ...lifecycle.value, group_revision: stale ? 'old' : item.revision,
      actions: [availability('stop', stop), availability('remove', remove)] } },
  } } });
  const state = initialCollectionState(), onDelete = vi.fn();
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const content = (record: WriteGroup) => <QueryClientProvider client={client}><RecordingHomeCollection groups={[record]} devices={[]} state={state} guidance={guidance} observedIds={[]} unverifiedObserved={false}
    sourcesAvailable restored={{ state, workspaceRevision: 'W2', scrollTop: 0, focusId: null }} onChange={vi.fn()} onOpen={vi.fn()}
    onOpenDevice={vi.fn()} onRemember={vi.fn()} onRecheck={vi.fn()} onDelete={onDelete} /></QueryClientProvider>;
  const view = render(content(item));
  const current: WriteGroupLifecycle = { workspace_id: 'W1', workspace_revision: 'W2', group_id: item.id, group_revision: item.revision,
    applied_revision: item.applied_revision, status: item.status, connector_id: item.destination.connector_id, connector_revision: item.destination.connector_revision,
    actions: [availability('stop', stop), availability('remove', remove)], history: { canonical: true, current_group_id: item.id, generations: [] } };
  return { row: screen.getByText('Recording R').closest('tr')!, onDelete, item, lifecycle: current,
    rerenderGroup: (record: WriteGroup) => view.rerender(content(record)) };
}

describe('recording row actions', () => {
  it('shows four unconfirmed stages when result evidence is unavailable', () => {
    const { row } = renderRow();
    const stages = [...row.querySelectorAll<HTMLElement>('[data-write-stage]')];
    expect(stages).toHaveLength(4);
    expect(stages.map(stage => stage.dataset.stageState)).toEqual(['unconfirmed', 'unconfirmed', 'unconfirmed', 'unconfirmed']);
    expect(stages[0]).not.toHaveTextContent('rowEvidence.draftNotRequested');
  });
  it('offers stop and delete as direct one-click buttons beside open, without a menu', () => {
    const { row } = renderRow();
    expect(within(row).queryByRole('button', { name: 'rowEvidence.actionsFor' })).toBeNull();
    expect(within(row).queryByRole('menu')).toBeNull();
    expect(within(row).getByRole('button', { name: 'rowEvidence.stopNamed' })).toBeEnabled();
    expect(within(row).getByRole('button', { name: 'rowEvidence.deleteNamed' })).toBeEnabled();
  });

  it('disables both actions with a reload hint when evidence has the wrong revision', () => {
    const { row, onDelete } = renderRow({ stale: true });
    expect(within(row).getByRole('button', { name: 'rowEvidence.stopNamed' })).toBeDisabled();
    const remove = within(row).getByRole('button', { name: 'rowEvidence.deleteNamed' });
    expect(remove).toBeDisabled();
    expect(within(row).getByText('rowEvidence.reloadHint')).toBeInTheDocument();
    fireEvent.click(remove);
    expect(onDelete).not.toHaveBeenCalled();
  });

  it('uses the server permissions and offers stop only for a running recording', () => {
    const { row, onDelete, item } = renderRow({ stop: false, remove: true });
    expect(within(row).getByRole('button', { name: 'rowEvidence.stopNamed' })).toBeDisabled();
    fireEvent.click(within(row).getByRole('button', { name: 'rowEvidence.deleteNamed' }));
    expect(onDelete).toHaveBeenCalledWith(item);
    cleanup();
    const draft = renderRow({ status: 'draft', remove: false });
    expect(within(draft.row).queryByRole('button', { name: 'rowEvidence.stopNamed' })).toBeNull();
    expect(within(draft.row).queryByRole('button', { name: 'rowEvidence.deleteNamed' })).toBeNull();
  });

  it('allows deleting a tombstone when matching server facts permit removal', () => {
    const { row, onDelete, item } = renderRow({ status: 'deleted' });
    fireEvent.click(within(row).getByRole('button', { name: 'rowEvidence.deleteNamed' }));
    expect(onDelete).toHaveBeenCalledExactlyOnceWith(item);
  });

  it('rechecks server authority, confirms impact, and submits stop with matching revisions', async () => {
    const { row, item, lifecycle } = renderRow();
    const read = vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockResolvedValue(lifecycle);
    const disable = vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'disable').mockResolvedValue({ group: { ...item, status: 'disabled' }, workspace_revision: 'W3' });
    fireEvent.click(within(row).getByRole('button', { name: 'rowEvidence.stopNamed' }));
    const confirm = await screen.findByTestId('group-stop-confirm-button');
    expect(disable).not.toHaveBeenCalled();
    fireEvent.click(confirm);
    await waitFor(() => expect(disable).toHaveBeenCalledWith('R', { workspace_id: 'W1', expected_workspace_revision: 'W2',
      expected_group_revision: 'R2', expected_connector_revision: 'DB2' }));
    expect(read).toHaveBeenCalledTimes(2);
  });

  it('returns focus to the stable recording button after cancelling a stop', async () => {
    const { row, lifecycle } = renderRow();
    vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockResolvedValue(lifecycle);
    fireEvent.click(within(row).getByRole('button', { name: 'rowEvidence.stopNamed' }));
    fireEvent.click(await screen.findByTestId('group-stop-cancel'));
    await waitFor(() => expect(within(row).getByRole('button', { name: 'openRecord' })).toHaveFocus());
  });

  it('returns focus to the recording when a successful stop removes the stop button', async () => {
    const { row, item, lifecycle, rerenderGroup } = renderRow();
    vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockResolvedValue(lifecycle);
    let resolve!: (value: Awaited<ReturnType<typeof studioV2WorkspaceWriteGroupsAPI.disable>>) => void;
    const disable = vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'disable').mockImplementation(() => new Promise(done => { resolve = done; }));
    fireEvent.click(within(row).getByRole('button', { name: 'rowEvidence.stopNamed' }));
    fireEvent.click(await screen.findByTestId('group-stop-confirm-button'));
    await waitFor(() => expect(disable).toHaveBeenCalledOnce());
    rerenderGroup({ ...item, status: 'disabled' });
    expect(within(row).queryByRole('button', { name: 'rowEvidence.stopNamed' })).not.toBeInTheDocument();
    resolve({ group: { ...item, status: 'disabled' }, workspace_revision: 'W3' });
    await waitFor(() => expect(within(row).getByRole('button', { name: 'openRecord' })).toHaveFocus());
  });

  it('blocks a stop if its authority changes after confirmation opened', async () => {
    const { row, lifecycle } = renderRow();
    vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockResolvedValueOnce(lifecycle).mockResolvedValueOnce({ ...lifecycle, group_revision: 'old' });
    const disable = vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'disable');
    fireEvent.click(within(row).getByRole('button', { name: 'rowEvidence.stopNamed' }));
    fireEvent.click(await screen.findByTestId('group-stop-confirm-button'));
    expect(await within(row).findByRole('alert')).toHaveTextContent('rowEvidence.reloadHint');
    expect(disable).not.toHaveBeenCalled();
    expect(within(row).getByRole('button', { name: 'rowEvidence.stopNamed' })).toBeDisabled();
    await waitFor(() => expect(within(row).getByRole('button', { name: 'openRecord' })).toHaveFocus());
  });

  it('returns focus to the recording if the stop authority cannot be loaded', async () => {
    const { row } = renderRow();
    vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockRejectedValue(new Error('offline'));
    fireEvent.click(within(row).getByRole('button', { name: 'rowEvidence.stopNamed' }));
    expect(await within(row).findByRole('alert')).toHaveTextContent('rowEvidence.reloadHint');
    await waitFor(() => expect(within(row).getByRole('button', { name: 'openRecord' })).toHaveFocus());
  });

  it('does not silently retry an uncertain stop or expose a raw error', async () => {
    const { row, lifecycle } = renderRow();
    vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockResolvedValue(lifecycle);
    vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'disable').mockRejectedValue(new Error('private backend error'));
    fireEvent.click(within(row).getByRole('button', { name: 'rowEvidence.stopNamed' }));
    fireEvent.click(await screen.findByTestId('group-stop-confirm-button'));
    expect(await within(row).findByRole('alert')).toHaveTextContent('rowEvidence.stopFailed');
    expect(row).not.toHaveTextContent('private backend error');
    expect(within(row).getByRole('button', { name: 'rowEvidence.stopNamed' })).toBeDisabled();
  });
});
