import { RecordingDraftNavigation } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, expect, it, vi } from 'vitest';
import { act, fireEvent, render, screen } from '@testing-library/react';
import { WriteGroupSection } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string, options?: Record<string, unknown>) => (options && typeof options === 'object' && 'count' in options ? `${key}:${options.count}` : key) }),
}));

const metadata: { status: string; columns: typeof TABLE_COLUMNS } = { status: 'exists', columns: TABLE_COLUMNS };
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({
  useStep4TargetColumns: () => ({ columns: metadata.columns, status: metadata.status, refetch: vi.fn() }),
}));

vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: {
    list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(),
    testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), get: vi.fn(),
  },
}));

const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);

function list(groups: WriteGroup[], revision = 'wrev-1') {
  return { workspace_id: 'ws-1', workspace_revision: revision, groups };
}

function mount() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const tree = (openGroupRequest?: { id: string }) => <QueryClientProvider client={client}><RecordingDraftNavigation>
    <WriteGroupSection state={savedGroupState()} workspaceId="ws-1" readonly={false} deviceSource={{ submit: vi.fn(), inspect: vi.fn(), observations: () => [] }} openGroupRequest={openGroupRequest} />
  </RecordingDraftNavigation></QueryClientProvider>;
  const view = render(tree()); return { client, request: (id: string) => view.rerender(tree({ id })) };
}
beforeEach(() => { vi.clearAllMocks(); api.list.mockResolvedValue(list([savedGroup(), savedGroup({ id: 'group-2', name: 'Other group' })])); api.readiness.mockResolvedValue(null as never); api.delivery.mockResolvedValue(null as never); });
it('a group switch keeps the existing form until the same guard is explicitly discarded', async () => {
  const view = mount(); fireEvent.click(await screen.findByTestId('group-open-group-1'));
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Keep exact local intent' } });
  view.request('group-2'); expect(await screen.findByRole('alertdialog')).toBeInTheDocument();
  fireEvent.click(screen.getByText('navigation.keep')); expect(screen.getByTestId('group-name')).toHaveValue('Keep exact local intent');
  view.request('group-2'); fireEvent.click(await screen.findByText('navigation.discard'));
  expect(screen.getByTestId('group-name')).toHaveValue('Other group');
});
it('failed background list/autosave readback cannot unmount and lose the edited form', async () => {
  const view = mount(); fireEvent.click(await screen.findByTestId('group-open-group-1'));
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Not saved during failed readback' } });
  api.list.mockRejectedValue(new Error('offline'));
  await act(async () => { await view.client.invalidateQueries(); });
  expect(screen.getByTestId('group-name')).toHaveValue('Not saved during failed readback');
  expect(screen.queryByTestId('group-saved-note')).not.toBeInTheDocument();
});
it('a failed pre-save read cannot use cached revisions to submit an ordinary recording', async () => {
  mount(); fireEvent.click(await screen.findByTestId('group-open-group-1'));
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Never submit cached scope' } });
  api.list.mockRejectedValue(new Error('fresh read failed'));
  fireEvent.click(screen.getByTestId('group-save')); await screen.findByTestId('group-save-error');
  expect(api.update).not.toHaveBeenCalled();
  expect(screen.getByTestId('group-name')).toHaveValue('Never submit cached scope');
});
it('IME Escape retains the source picker pending selection; ordinary Escape cancels it', async () => {
  mount(); fireEvent.click(await screen.findByTestId('group-open-group-1'));
  fireEvent.click(screen.getByTestId('recording-points-open'));
  const checkbox = screen.getByRole('checkbox', { name: 'line.batch' }); fireEvent.click(checkbox);
  const search = screen.getByRole('searchbox');
  fireEvent.keyDown(search, { key: 'Escape', isComposing: true });
  expect(screen.getByRole('dialog')).toBeInTheDocument(); expect(checkbox).toBeChecked();
  fireEvent.keyDown(search, { key: 'Escape' });
  expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
});
