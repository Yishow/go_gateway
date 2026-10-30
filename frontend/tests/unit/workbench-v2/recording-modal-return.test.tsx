import { markReturnRegion } from '@/utils/returnFocus';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { GroupLifecycleBar } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupLifecycleBar';
import { savedGroup } from '../../fixtures/writeGroupState';
import { savedLifecycle } from '../../fixtures/writeGroupLifecycle';
const mocks = vi.hoisted(() => ({ lifecycle: vi.fn(), refetch: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn() }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/hooks/datalink/useStudioV2WriteGroupLifecycle', () => ({ useWriteGroupLifecycleQuery: mocks.lifecycle }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: mocks }));
let group = savedGroup();
beforeEach(() => {
  vi.clearAllMocks(); group = savedGroup();
  mocks.lifecycle.mockImplementation(() => ({ data: savedLifecycle(group), isError: false, refetch: mocks.refetch }));
  mocks.refetch.mockImplementation(async () => ({ data: savedLifecycle(group) }));
});
afterEach(cleanup);
function setup() {
  const client = new QueryClient({ defaultOptions: { mutations: { retry: false }, queries: { retry: false } } });
  const changed = vi.fn();
  const view = () => <QueryClientProvider client={client}><GroupLifecycleBar group={group} workspaceId={group.workspace_id} workspaceRevision="wrev-1" canApply readonly={false} onChanged={changed} onConflict={vi.fn()} /></QueryClientProvider>;
  const result = render(view()); return Object.assign(() => result.rerender(view()), { changed });
}
it.each(['revision', 'group', 'workspace'])('an open remove consent cannot be reused after %s changes, even after A→B→A', async kind => {
  const rerender = setup(); const original = group;
  fireEvent.click(screen.getByTestId('group-delete')); await screen.findByRole('alertdialog');
  group = { ...group, ...(kind === 'revision' ? { revision: 'rev-2' } : kind === 'group' ? { id: 'group-2' } : { workspace_id: 'ws-2' }) }; rerender();
  expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument();
  group = original; rerender(); expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument(); expect(mocks.remove).not.toHaveBeenCalled();
  fireEvent.click(screen.getByTestId('group-delete')); await screen.findByRole('alertdialog');
});
it('closing a confirmation with a hidden trigger returns to its current safe context rather than the hidden action', async () => {
  setup(); const trigger = screen.getByTestId('group-delete'); trigger.focus(); fireEvent.click(trigger); await screen.findByRole('alertdialog');
  trigger.hidden = true; fireEvent.click(screen.getByTestId('group-delete-cancel'));
  await waitFor(() => expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument()); expect(trigger).not.toHaveFocus();
  expect(screen.getByTestId('group-lifecycle-return')).toHaveFocus();
});

it('a late pre-confirmation read after A→B→A cannot reopen old consent', async () => {
  let reply!: (value: { data: ReturnType<typeof savedLifecycle> }) => void;
  mocks.refetch.mockReturnValueOnce(new Promise(resolve => { reply = resolve; }));
  const rerender = setup(); const original = group; fireEvent.click(screen.getByTestId('group-delete'));
  await waitFor(() => expect(mocks.refetch).toHaveBeenCalled());
  group = { ...original, revision: 'rev-2' }; rerender(); group = original; rerender();
  await act(async () => reply({ data: savedLifecycle(original) }));
  expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument(); expect(mocks.remove).not.toHaveBeenCalled();
});

it('an issued remove reply cannot settle a new A→B→A owner or refocus the reused trigger', async () => {
  let reply!: (value: { group: typeof group; workspace_revision: string }) => void;
  mocks.remove.mockReturnValueOnce(new Promise(resolve => { reply = resolve; }));
  const rerender = setup(); const original = group; fireEvent.click(screen.getByTestId('group-delete'));
  fireEvent.click(await screen.findByTestId('group-delete-confirm-button')); await waitFor(() => expect(mocks.remove).toHaveBeenCalledOnce());
  group = { ...original, revision: 'rev-2' }; rerender(); group = original; rerender();
  await act(async () => reply({ group: { ...original, status: 'deleted' }, workspace_revision: 'wrev-2' }));
  expect(rerender.changed).not.toHaveBeenCalled(); expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument();
  expect(screen.getByTestId('group-delete')).not.toHaveFocus();
  fireEvent.click(screen.getByTestId('group-delete')); await screen.findByRole('alertdialog');
  expect(mocks.remove).toHaveBeenCalledOnce();
});

it('closing Advanced while remove review is loading returns to the visible Step4 heading on cancel', async () => {
  let reply!: (value: { data: ReturnType<typeof savedLifecycle> }) => void;
  mocks.refetch.mockReturnValueOnce(new Promise(resolve => { reply = resolve; }));
  const owner = {}; const client = new QueryClient({ defaultOptions: { mutations: { retry: false } } });
  render(<QueryClientProvider client={client}><main ref={node => markReturnRegion(node, owner)} data-focus-return-region="ws-1/step4">
    <h2 tabIndex={-1} data-focus-return-heading>Step4 heading</h2><details open><summary>Advanced</summary>
      <GroupLifecycleBar group={group} workspaceId={group.workspace_id} workspaceRevision="wrev-1" canApply readonly={false} onChanged={vi.fn()} onConflict={vi.fn()} />
    </details></main></QueryClientProvider>);
  const trigger = screen.getByTestId('group-delete'); trigger.focus(); fireEvent.click(trigger);
  fireEvent.click(screen.getByText('Advanced')); expect(screen.getByText('Advanced').closest('details')).not.toHaveAttribute('open');
  await act(async () => reply({ data: savedLifecycle(group) })); await screen.findByRole('alertdialog');
  fireEvent.click(screen.getByTestId('group-delete-cancel'));
  await waitFor(() => expect(screen.getByRole('heading', { name: 'Step4 heading' })).toHaveFocus());
  expect(trigger).not.toHaveFocus(); expect(mocks.remove).not.toHaveBeenCalled();
});
