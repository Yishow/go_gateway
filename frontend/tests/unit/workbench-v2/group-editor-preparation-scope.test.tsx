import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { GroupEditor } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { managedSchemaPreviewToken } from '../../fixtures/managedSchemaPreview';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({
  useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }),
}));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  readiness: vi.fn().mockResolvedValue(null), delivery: vi.fn().mockResolvedValue(null), list: vi.fn(), get: vi.fn(),
  create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(), testWritePreview: vi.fn(),
  testWrite: vi.fn(), testWriteOperation: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(),
} }));

const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const state = savedGroupState();
const clients: QueryClient[] = [];
const groupA = savedGroup({ id: 'record-a', name: 'Saved A' });
function mount(group: WriteGroup | undefined, existingClient?: QueryClient) {
  const client = existingClient ?? new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  if (!existingClient) clients.push(client);
  const onClose = vi.fn(), onSaved = vi.fn(), onReload = vi.fn();
  const fetchLatest = vi.fn(async () => ({ workspace_revision: 'w1', groups: group ? [group] : [] }));
  const view = render(<QueryClientProvider client={client}><GroupEditor group={group} workspaceId="ws-1" workspaceRevision="w1"
    connector={state.db.connector} candidates={buildGroupCandidates(state).candidates} excluded={[]} readonly={false}
    onClose={onClose} onSaved={onSaved} onReload={onReload} fetchLatest={fetchLatest} /></QueryClientProvider>);
  return { ...view, client, onClose, onSaved, onReload, fetchLatest };
}
async function loseSaveA() {
  api.update.mockRejectedValueOnce(new Error('Acknowledgement lost'));
  const view = mount(groupA);
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Unknown A intent' } });
  fireEvent.click(screen.getByTestId('group-save'));
  await screen.findByTestId('group-save-error');
  await waitFor(() => expect(view.client.isMutating()).toBe(0));
  return view;
}
async function leaveUnknown(view: ReturnType<typeof mount>, review = false) {
  fireEvent.click(screen.getByTestId(review ? 'group-save-review' : 'group-editor-close'));
  expect(await screen.findByRole('alertdialog')).toHaveTextContent('navigation.pending');
  fireEvent.click(screen.getByRole('button', { name: 'navigation.discard' }));
  expect(view.onClose).toHaveBeenCalledOnce();
  view.unmount();
}
function noPreparationWrites() {
  expect(api.schemaApply).not.toHaveBeenCalled();
  expect(api.apply).not.toHaveBeenCalled();
  expect(api.testWrite).not.toHaveBeenCalled();
}
beforeEach(() => { vi.clearAllMocks(); sessionStorage.clear(); });
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });

it.each(['custom', 'managed'] as const)('A unknown does not block clean %s B from its explicit Check', async storage => {
  const first = await loseSaveA();
  await leaveUnknown(first);
  const groupB = savedGroup({ id: 'record-b', name: 'Saved B', destination: { ...savedGroup().destination, storage_strategy: storage } });
  api.schemaPreview.mockResolvedValue({ ...managedSchemaPreviewToken, workspace_id: 'ws-1', workspace_revision: 'w1',
    plan_id: 'record-b', plan_revision: 'rev-1', connector_id: 'conn-1', connector_revision: 'crev-1', expires_at: '2099-01-01T00:00:00Z' });
  mount(groupB, first.client);
  expect(await screen.findByTestId('group-save-unknown')).toHaveTextContent('navigation.unknownSave');
  const check = await screen.findByTestId('recording-preparation-check');
  await waitFor(() => expect(check).toBeEnabled());
  expect(api.schemaPreview).not.toHaveBeenCalled();
  const readinessReads = api.readiness.mock.calls.length;
  fireEvent.click(check);
  if (storage === 'managed') {
    await waitFor(() => expect(api.schemaPreview).toHaveBeenCalledExactlyOnceWith('record-b', {
      workspace_id: 'ws-1', expected_workspace_revision: 'w1', expected_group_revision: 'rev-1', expected_connector_revision: 'crev-1',
    }));
    expect(await screen.findByTestId('recording-preparation-confirm')).toBeEnabled();
  } else {
    await waitFor(() => expect(api.readiness).toHaveBeenCalledTimes(readinessReads + 1));
    expect(api.readiness).toHaveBeenLastCalledWith('record-b');
  }
  expect(api.update).toHaveBeenCalledOnce();
  expect(api.create).not.toHaveBeenCalled();
  noPreparationWrites();
});

it('A unknown does not poison a new B after its own successful Save acknowledgement', async () => {
  const first = await loseSaveA();
  await leaveUnknown(first);
  const groupB = savedGroup({ id: 'record-b', name: 'New B', members: [savedGroup().members[0]] });
  api.create.mockResolvedValueOnce({ workspace_revision: 'w2', group: groupB });
  const next = mount(undefined, first.client);
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'New B' } });
  fireEvent.click(screen.getByTestId('group-member-include-pt-0'));
  fireEvent.change(screen.getByTestId('group-member-column-pt-0'), { target: { value: 'temperature' } });
  fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(next.onSaved).toHaveBeenCalledOnce());
  const check = await screen.findByTestId('recording-preparation-check');
  await waitFor(() => expect(check).toBeEnabled());
  expect(screen.getByTestId('group-save-unknown')).toHaveTextContent('navigation.unknownSave');
  expect(screen.getByTestId('group-saved-note')).toBeInTheDocument();
  const reads = api.readiness.mock.calls.length;
  fireEvent.click(check);
  await waitFor(() => expect(api.readiness).toHaveBeenCalledTimes(reads + 1));
  expect(api.readiness).toHaveBeenLastCalledWith('record-b');
  expect(api.create).toHaveBeenCalledOnce();
  expect(api.update).toHaveBeenCalledOnce();
  noPreparationWrites();
});

it('the still-mounted unknown A retains its intent and cannot resubmit or prepare after editing again', async () => {
  const first = await loseSaveA();
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Further A intent' } });
  expect(screen.getByTestId('recording-preparation-phase')).toHaveTextContent('preparation.stage.saveUnknown');
  expect(screen.getByTestId('group-save')).toBeDisabled();
  expect(screen.queryByTestId('recording-preparation-check')).not.toBeInTheDocument();
  fireEvent.click(screen.getByTestId('group-save'));
  expect(api.update).toHaveBeenCalledOnce();
  expect(first.fetchLatest).toHaveBeenCalledOnce();
  expect(screen.getByTestId('group-name')).toHaveValue('Further A intent');
  noPreparationWrites();
});

it('review and reopen A restores canonical fields, never retries old intent, and preserves fresh group CAS', async () => {
  const first = await loseSaveA();
  await leaveUnknown(first, true);
  expect(first.onReload).toHaveBeenCalledOnce();
  const reopened = mount(groupA, first.client);
  expect(screen.getByTestId('group-name')).toHaveValue('Saved A');
  expect(screen.getByTestId('group-save-unknown')).toHaveTextContent('navigation.unknownSave');
  expect(api.update).toHaveBeenCalledOnce();
  expect(reopened.fetchLatest).not.toHaveBeenCalled();
  reopened.fetchLatest.mockResolvedValueOnce({ workspace_revision: 'w2', groups: [{ ...groupA, revision: 'rev-2' }] });
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Explicit new A intent' } });
  fireEvent.click(screen.getByTestId('group-save'));
  await screen.findByTestId('group-conflict');
  expect(api.update).toHaveBeenCalledOnce();
  expect(screen.getByTestId('group-name')).toHaveValue('Explicit new A intent');
  noPreparationWrites();
});

it('reopened A cannot submit while the matching original save is still pending', async () => {
  let resolve!: (value: { workspace_revision: string; group: WriteGroup }) => void;
  api.update.mockReturnValueOnce(new Promise(done => { resolve = done; }));
  const first = mount(groupA);
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Pending A' } });
  fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(api.update).toHaveBeenCalledOnce());
  await leaveUnknown(first);
  mount(groupA, first.client);
  expect(screen.getByTestId('group-save-fields')).toBeDisabled();
  expect(screen.getByTestId('group-save')).toBeDisabled();
  fireEvent.click(screen.getByTestId('group-save'));
  expect(api.update).toHaveBeenCalledOnce();
  await act(async () => resolve({ workspace_revision: 'w2', group: { ...groupA, revision: 'rev-2', name: 'Pending A' } }));
  await waitFor(() => expect(screen.getByTestId('group-save-fields')).not.toBeDisabled());
  expect(first.onSaved).not.toHaveBeenCalled();
  expect(api.update).toHaveBeenCalledOnce();
});
