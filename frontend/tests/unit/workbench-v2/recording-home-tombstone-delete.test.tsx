import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { RecordingHome } from '@/features/datalink/workbench-v2/recording-home/RecordingHome';
import { RecordingDeleteDialogController } from '@/features/datalink/workbench-v2/recording-home/RecordingDeleteDialogController';
import { recordingHomeGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import { studioV2WriteGroupLifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { savedLifecycle } from '../../fixtures/writeGroupLifecycle';
import { fixture, group } from './recording-home-guidance-fixtures';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
afterEach(() => { cleanup(); vi.restoreAllMocks(); });

function mountTombstone(stale = false) {
  const record = group('Test', 'deleted');
  const facts = savedLifecycle(record, 'W2');
  facts.actions = facts.actions.map(action => action.action === 'remove' ? { ...action, allowed: true } : action);
  const input = fixture(['D1'], true, [record]);
  const evidence = input.evidence.lifecycle?.Test;
  if (!evidence || evidence.state !== 'available') throw new Error('lifecycle fixture missing');
  const guidance = recordingHomeGuidance({ ...input, owners: undefined, evidence: { ...input.evidence, lifecycle: {
    Test: { ...evidence, value: { ...facts, group_revision: stale ? 'old' : record.revision } },
  } } });
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><RecordingDeleteDialogController workspaceId="W1" workspaceRevision="W2"
    groups={[record]} location={{ search: '', hash: '' }} refresh={async () => true}>
    {actions => <RecordingHome guidance={guidance} location={{ search: '', hash: '' }} onRetry={vi.fn()} onDelete={actions.onDelete} />}
  </RecordingDeleteDialogController></QueryClientProvider>);
  const disclosure = screen.getByTestId('recording-home-unconfirmed');
  return { record, facts, disclosure, client };
}

it('permanently deletes an unconfirmed tombstone after matching server permission and acknowledgement', async () => {
  const { record, facts, disclosure, client } = mountTombstone();
  vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockResolvedValue(facts);
  const remove = vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'remove').mockResolvedValue({ deleted: true, group_id: record.id, lineage_id: '', deleted_generations: 1 });
  fireEvent.click(within(disclosure).getByRole('button', { name: 'rowEvidence.deleteNamed' }));
  fireEvent.click(screen.getByRole('checkbox', { name: 'deleteDialogs.recording.acknowledge' }));
  fireEvent.click(screen.getByTestId('group-delete-confirm-button'));
  await waitFor(() => expect(remove).toHaveBeenCalledWith('Test', { workspace_id: 'W1', expected_workspace_revision: 'W2',
    expected_group_revision: 'Test2', expected_connector_revision: 'DB2' }));
  client.clear();
});

it('disables a tombstone delete when its lifecycle revision is stale', () => {
  const { disclosure, client } = mountTombstone(true);
  expect(within(disclosure).getByRole('button', { name: 'rowEvidence.deleteNamed' })).toBeDisabled();
  expect(within(disclosure).getByText('rowEvidence.reloadHint')).toBeInTheDocument();
  client.clear();
});

it('refuses to submit when tombstone authority changes after opening the dialog', async () => {
  const { facts, disclosure, client } = mountTombstone();
  vi.spyOn(studioV2WriteGroupLifecycleAPI, 'lifecycle').mockResolvedValue({ ...facts, connector_revision: 'changed' });
  const remove = vi.spyOn(studioV2WorkspaceWriteGroupsAPI, 'remove');
  fireEvent.click(within(disclosure).getByRole('button', { name: 'rowEvidence.deleteNamed' }));
  fireEvent.click(screen.getByRole('checkbox', { name: 'deleteDialogs.recording.acknowledge' }));
  fireEvent.click(screen.getByTestId('group-delete-confirm-button'));
  expect(await screen.findByText('deleteDialogs.recording.stale')).toBeInTheDocument();
  expect(remove).not.toHaveBeenCalled();
  client.clear();
});
