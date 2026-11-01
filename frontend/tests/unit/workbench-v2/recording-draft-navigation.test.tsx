import { RecordingDraftReadState } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { GroupEditor } from '../../../src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { buildGroupCandidates } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
import { savedGroup, savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';
import { studioV2WorkspaceWriteGroupsAPI } from '../../../src/services/studioV2WorkspaceWriteGroups';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('../../../src/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({ useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }) }));
vi.mock('../../../src/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  readiness: vi.fn().mockResolvedValue(null), delivery: vi.fn().mockResolvedValue(null), list: vi.fn(), get: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(), testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(),
} }));
function mount(preflight?: Promise<{ workspace_revision: string; groups: ReturnType<typeof savedGroup>[] }>, existingClient?: QueryClient, scope?: { workspaceId: string; id: string }) {
  const state = savedGroupState(), group = savedGroup(scope ? { workspace_id: scope.workspaceId, id: scope.id } : {}), onClose = vi.fn(), onSaved = vi.fn();
  const client = existingClient ?? new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const tree = (workspaceId = scope?.workspaceId ?? "ws-1", unavailable = false) => <QueryClientProvider client={client}><RecordingDraftReadState unavailable={unavailable}><GroupEditor group={group} workspaceId={workspaceId} workspaceRevision="db-setup-9" connector={state.db.connector} candidates={buildGroupCandidates(state).candidates} excluded={[]} readonly={false} onSaved={onSaved} onClose={onClose} onReload={vi.fn()} fetchLatest={async () => preflight ?? ({ workspace_revision: 'db-setup-9', groups: [group] })} /></RecordingDraftReadState></QueryClientProvider>;
  const view = render(tree());
  return { ...view, client, onClose, onSaved, changeWorkspace: () => view.rerender(tree("ws-2")), setReadUnavailable: (value: boolean) => view.rerender(tree("ws-1", value)) };
}
beforeEach(() => vi.clearAllMocks());
it('Close retains the complete unsaved form until explicit discard and explains saved sources remain', async () => {
  const { onClose } = mount();
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Unsaved operator intent' } });
  fireEvent.click(screen.getByTestId('group-editor-close'));
  expect(onClose).not.toHaveBeenCalled();
  expect(await screen.findByRole('alertdialog')).toHaveTextContent('navigation.savedSources');
  fireEvent.click(screen.getByRole('button', { name: 'navigation.keep' }));
  expect(screen.getByTestId('group-name')).toHaveValue('Unsaved operator intent');
  fireEvent.click(screen.getByTestId('group-editor-close'));
  fireEvent.click(screen.getByRole('button', { name: 'navigation.discard' }));
  expect(onClose).toHaveBeenCalledOnce();
  expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
  expect(studioV2WorkspaceWriteGroupsAPI.remove).not.toHaveBeenCalled();
});
it('Escape offers the same guard; dismissing the confirmation retains the input and restores focus', async () => {
  const { onClose } = mount(); const name = screen.getByTestId('group-name');
  fireEvent.change(name, { target: { value: 'Retained' } }); name.focus();
  fireEvent.keyDown(name, { key: 'Escape' });
  expect(await screen.findByRole('alertdialog')).toBeInTheDocument();
  fireEvent.keyDown(screen.getByRole('alertdialog'), { key: 'Escape' });
  await waitFor(() => expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument());
  expect(name).toHaveValue('Retained'); expect(name).toHaveFocus(); expect(onClose).not.toHaveBeenCalled();
});
it('a failed save remains unsaved and warns on reload rather than claiming browser persistence', async () => {
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mockRejectedValue(new Error('offline'));
  const view = mount(); fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Not persisted' } });
  fireEvent.click(screen.getByTestId('group-save')); await screen.findByTestId('group-save-error');
  const unload = new Event('beforeunload', { cancelable: true }); window.dispatchEvent(unload);
  expect(unload.defaultPrevented).toBe(true);
  expect(screen.getByTestId('group-name')).toHaveValue('Not persisted');
  expect(screen.queryByTestId('group-saved-note')).not.toBeInTheDocument();
  view.unmount(); mount();
  expect(screen.getByTestId('group-name')).toHaveValue(savedGroup().name);
});

it('leaving during preflight cannot submit a late mutation after the editor was discarded', async () => {
  let resolve!: (value: { workspace_revision: string; groups: ReturnType<typeof savedGroup>[] }) => void;
  const view = mount(new Promise(r => { resolve = r; }));
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Not submitted' } });
  fireEvent.click(screen.getByTestId('group-save'));
  fireEvent.click(screen.getByTestId('group-editor-close'));
  expect(await screen.findByRole('alertdialog')).toHaveTextContent('navigation.pending');
  fireEvent.click(screen.getByText('navigation.discard')); view.unmount();
  await act(async () => resolve({ workspace_revision: 'db-setup-9', groups: [savedGroup()] }));
  expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
});
it('leaving a submitted save does not call cancellation or let its late success reopen the editor', async () => {
  let resolve!: (value: { workspace_revision: string; group: ReturnType<typeof savedGroup> }) => void;
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mockReturnValue(new Promise(r => { resolve = r; }));
  const view = mount(); fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Submitted' } });
  fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledOnce());
  fireEvent.click(screen.getByTestId('group-editor-close')); fireEvent.click(screen.getByText('navigation.discard')); view.unmount();
  await act(async () => resolve({ workspace_revision: 'db-saved-10', group: savedGroup({ name: 'Submitted' }) }));
  expect(view.onSaved).not.toHaveBeenCalled(); expect(studioV2WorkspaceWriteGroupsAPI.remove).not.toHaveBeenCalled();
});

it('a late submitted save cannot publish into a newer viewing scope', async () => {
  let resolve!: (value: { workspace_revision: string; group: ReturnType<typeof savedGroup> }) => void;
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mockReturnValue(new Promise(r => { resolve = r; }));
  const view = mount(); fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Old workspace draft' } });
  fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledOnce());
  view.changeWorkspace();
  await act(async () => resolve({ workspace_revision: 'db-saved-10', group: savedGroup({ name: 'Old workspace response' }) }));
  await waitFor(() => expect(screen.getByTestId('group-save-fields')).not.toBeDisabled());
  expect(view.onSaved).not.toHaveBeenCalled();
  expect(screen.getByTestId('group-name')).not.toHaveValue('Old workspace response');
});

it('an uncertain submitted save cannot be blindly submitted again with a freshly rebased revision', async () => {
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mockRejectedValue(new Error('connection lost'));
  mount(); fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Uncertain submission' } });
  fireEvent.click(screen.getByTestId('group-save')); await screen.findByTestId('group-save-error');
  expect(screen.getByTestId('group-save')).toBeDisabled();
  expect(screen.getByTestId('group-save-review')).toBeEnabled();
  fireEvent.click(screen.getByTestId('group-save')); expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledOnce();
});
it('reopening the same saved record cannot submit while its earlier save remains in flight', async () => {
  let resolve!: (value: { workspace_revision: string; group: ReturnType<typeof savedGroup> }) => void;
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.update).mockReturnValue(new Promise(r => { resolve = r; }));
  const first = mount(); fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Submitted' } });
  fireEvent.click(screen.getByTestId('group-save')); await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledOnce());
  first.unmount(); mount(undefined, first.client);
  expect(screen.getByTestId('group-save-fields')).toBeDisabled();
  expect(screen.getByTestId('group-save')).toBeDisabled();
  await act(async () => resolve({ workspace_revision: 'db-saved-10', group: savedGroup({ name: 'Submitted' }) }));
  await waitFor(() => expect(screen.getByTestId('group-save-fields')).not.toBeDisabled());
  expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledOnce();
});

it('a failed page read retains fields but denies new commands until a fresh read succeeds', () => {
  const view = mount(); fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Preserved through refresh failure' } });
  view.setReadUnavailable(true);
  expect(screen.getByTestId('group-save-fields')).toBeDisabled(); expect(screen.getByTestId('group-save')).toBeDisabled();
  fireEvent.click(screen.getByTestId('group-save')); expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
  view.setReadUnavailable(false); expect(screen.getByTestId('group-name')).toHaveValue('Preserved through refresh failure');
  expect(screen.getByTestId('group-save-fields')).not.toBeDisabled();
});

it('an old scope finally cannot unlock the new scope save and does not borrow its result', async () => {
  type Result = { workspace_revision: string; group: ReturnType<typeof savedGroup> };
  let oldResolve!: (value: Result) => void, newResolve!: (value: Result) => void;
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.update)
    .mockImplementationOnce(() => new Promise(resolve => { oldResolve = resolve; }))
    .mockImplementationOnce(() => new Promise(resolve => { newResolve = resolve; }));
  const first = mount(); fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Old scope' } });
  fireEvent.click(screen.getByTestId('group-save')); await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledOnce());
  first.unmount(); const next = mount(undefined, first.client, { workspaceId: 'ws-2', id: 'group-2' });
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'New scope' } }); fireEvent.click(screen.getByTestId('group-save'));
  await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.update).toHaveBeenCalledTimes(2));
  await act(async () => oldResolve({ workspace_revision: 'old-saved', group: savedGroup({ name: 'Old saved' }) }));
  expect(screen.getByTestId('group-save-fields')).toBeDisabled(); expect(next.onSaved).not.toHaveBeenCalled();
  expect(screen.getByTestId('group-name')).toHaveValue('New scope');
  await act(async () => newResolve({ workspace_revision: 'new-saved', group: savedGroup({ workspace_id: 'ws-2', id: 'group-2', name: 'New scope' }) }));
  await waitFor(() => expect(next.onSaved).toHaveBeenCalledOnce()); expect(first.onSaved).not.toHaveBeenCalled();
});
it('a clean saved record closes normally without an unnecessary prompt or reload warning', () => {
  const view = mount(); fireEvent.click(screen.getByTestId('group-editor-close'));
  expect(view.onClose).toHaveBeenCalledOnce(); expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument();
  const unload = new Event('beforeunload', { cancelable: true }); window.dispatchEvent(unload); expect(unload.defaultPrevented).toBe(false);
});
it('destination save-state changes keep the same editor content rather than reseeding dirty fields', () => {
  const state = savedGroupState(); const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const tree = (persisted: boolean) => <QueryClientProvider client={client}><GroupEditor group={savedGroup()} workspaceId="ws-1" workspaceRevision="w1" connector={{ ...state.db.connector, kind: 'mysql', persisted }} candidates={buildGroupCandidates(state).candidates} excluded={[]} readonly={false} onClose={vi.fn()} onSaved={vi.fn()} onReload={vi.fn()} fetchLatest={vi.fn()} /></QueryClientProvider>;
  const view = render(tree(true));
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Unsent across destination authority change' } });
  view.rerender(tree(false));
  expect(screen.getByTestId('group-name')).toHaveValue('Unsent across destination authority change');
  expect(screen.getByTestId('group-save')).toBeDisabled();
});

import type { WriteGroupCalculationCapabilities } from '../../../src/types/studioV2WriteGroupCalculation';
const capabilities: WriteGroupCalculationCapabilities = { version: 1, result_shape: 'member-calculation-v1', editor_supported: true, draft_supported: true, activation_supported: false, reason: 'record_readiness_required', methods: ['latest', 'mean', 'delta'], cadence: { alignment: 'utc_epoch', units: ['second', 'minute'], minimum: 1, max_seconds: 600, max_minutes: 10, window_bounds_checked: true }, schema_dialects: ['sqlite'], schema_layouts: ['managed_wide_profile_v1', 'managed_long_v1'], unsupported_combinations: [] };
const base = savedGroup();
const group = savedGroup({ record_id: '11111111-1111-4111-8111-111111111111', layout_version: 'managed_wide_profile_v1', recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' }, destination: { ...base.destination, storage_strategy: 'managed' }, row_policy: { allowed_lateness_seconds: 0, incomplete_policy: 'partial' }, members: base.members.map((member, index) => ({ ...member, calculation: 'latest', cadence: { interval: index ? 60 : 15, unit: 'second', alignment: 'utc_epoch' }, value_role: 'gauge', source_value_type: 'float64' })) });

it('capability withdrawal retains the mixed draft but denies Save', () => {
 const state = savedGroupState(); const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
 const tree = (supported: boolean) => <QueryClientProvider client={client}><GroupEditor group={group} calculationCapabilities={{ ...capabilities, editor_supported: supported }} workspaceId="ws-1" workspaceRevision="w1" connector={state.db.connector} candidates={buildGroupCandidates(state).candidates} excluded={[]} readonly={false} onClose={vi.fn()} onSaved={vi.fn()} onReload={vi.fn()} fetchLatest={vi.fn()} /></QueryClientProvider>;
 const view = render(tree(true)); fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Keep mixed intent' } });
 view.rerender(tree(false)); expect(screen.getByTestId('group-name')).toHaveValue('Keep mixed intent');
 expect(screen.getByTestId('group-save')).toBeDisabled();
});
