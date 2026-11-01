import { useState } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { createMemoryRouter, RouterProvider, useLocation, useNavigate } from 'react-router-dom';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { WriteGroupSection } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection';
import type { GroupEditorProps } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { RecordingRouterNavigation } from '@/features/datalink/workbench-v2/recording-editor/RecordingRouterNavigation';
import { useRecordingDraftGuard } from '@/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import { buildSimpleRecordingTarget, parseStudioPresentation } from '@/services/studioRecordingView';
import { readStudioLocatorEditorEntry } from '@/services/studioV2LocatorEntry';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { savedGroupState } from '../../fixtures/writeGroupState';
import { locatorRead } from './studio-locator-fixtures';
import type { WriteGroupMutationResponse } from '@/types/studioV2WriteGroup';
import type { StudioLocatorEntryRead, StudioLocatorEntryResult } from '@/types/studioV2LocatorEntry';
const captured = vi.hoisted(() => ({ props: undefined as GroupEditorProps | undefined, lateInput: undefined as (() => void) | undefined, switchScope: undefined as ((id: string) => void) | undefined, nextAck: undefined as WriteGroupMutationResponse | undefined }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/SavedRecordingCard', () => ({ SavedRecordingCard: () => null }));
vi.mock('@/services/studioV2LocatorEntry', () => ({ readStudioLocatorEditorEntry: vi.fn() }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn() } }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor', () => ({ GroupEditor: function Owner(props: GroupEditorProps) {
  captured.props = props;
  const [name, setName] = useState(props.group?.name ?? 'Local new recording');
  const [dirty, setDirty] = useState(!props.group);
  captured.lateInput = () => { setName('Newer unsaved input'); setDirty(true); };
  useRecordingDraftGuard({ version: name, dirty });
  return <div data-testid="canonical-owner"><input aria-label="Recording name" value={name} disabled={props.readonly}
    onChange={event => { setName(event.target.value); setDirty(true); }} />
    <button onClick={() => { setDirty(false); const acknowledgement = captured.nextAck ?? { group: read.entry.write_groups!.groups[0], workspace_revision: 'W2' }; props.onSaved(acknowledgement.group, acknowledgement.workspace_revision, acknowledgement); }}>Accepted original Save</button>
    <button onClick={props.onReload}>Original Reload</button></div>;
} }));
const read = ((): StudioLocatorEntryRead => { const value = locatorRead(4); return { ...value, entry: { ...value.entry, editor_hydrated: true, database_config: null, database_targets: [] } }; })();
const clients: QueryClient[] = [];
function mount(url = '/studio/v2?view=recording&intent=new') {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); clients.push(client);
  function Entry() {
    const location = useLocation(), navigate = useNavigate();
    const [workspaceId, setWorkspaceId] = useState('W'); captured.switchScope = setWorkspaceId;
    const intent = parseStudioPresentation(location.search);
    const located = intent.locatorIntent.kind === 'locator';
    return <WriteGroupSection state={savedGroupState()} workspaceId={workspaceId} readonly={false} presentation="simple"
      initialSelection={located ? 'located' : 'new'} locatorOnly={located} locatorEntry={located ? read : undefined}
      onSavedLocation={group => { const target = buildSimpleRecordingTarget({ kind: 'edit', workspaceId: group.workspace_id, groupId: group.id, groupRevision: group.revision }); if (target) navigate(target, { replace: true }); }} />;
  }
  const router = createMemoryRouter([{ path: '/studio/v2', element: <QueryClientProvider client={client}><RecordingRouterNavigation><Entry /></RecordingRouterNavigation></QueryClientProvider> }], { initialEntries: [url] });
  return { router, ...render(<RouterProvider router={router} />) };
}
beforeEach(() => {
  vi.resetAllMocks(); captured.props = undefined; captured.lateInput = undefined; captured.nextAck = undefined; captured.switchScope = undefined;
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue(read.entry.write_groups!);
  vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue({ ok: true, value: read });
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });
describe('first accepted save becomes a reloadable canonical location', () => {
  it('replaces only after matching readback, then reloads the same group without another save', async () => {
    const view = mount(); await screen.findByTestId('canonical-owner');
    expect(captured.props?.group).toBeUndefined(); fireEvent.click(screen.getByText('Accepted original Save'));
    await waitFor(() => expect(view.router.state.location.search).toContain('entry=locator'));
    expect(view.router.state.historyAction).toBe('REPLACE');
    expect(view.router.state.location.search).toContain('group_id=G'); expect(view.router.state.location.search).toContain('group_revision=G3');
    expect(readStudioLocatorEditorEntry).toHaveBeenCalled();
    const savedURL = view.router.state.location.pathname + view.router.state.location.search;
    view.unmount(); mount(savedURL); await screen.findByTestId('canonical-owner');
    expect(captured.props?.group).toMatchObject({ id: 'G', workspace_id: 'W', revision: 'G3' });
    expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled(); expect(studioV2WorkspaceWriteGroupsAPI.update).not.toHaveBeenCalled();
  });
  it.each(['foreign', 'failed', 'wrong-success'] as const)('keeps the saved draft and new location when readback is %s; retry uses reads only', async mode => {
    vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue(mode === 'wrong-success'
      ? { ok: true, value: { ...read, selection: locatorRead(4, 'OTHER').selection } }
      : { ok: false, reason: mode === 'foreign' ? 'foreign-scope' : 'unavailable' });
    const view = mount(); await screen.findByTestId('canonical-owner'); fireEvent.click(screen.getByText('Accepted original Save'));
    await screen.findByTestId('studio-locator-open-error');
    expect(view.router.state.location.search).toBe('?view=recording&intent=new');
    expect(captured.props?.group?.id).toBe('G'); expect(screen.getByLabelText('Recording name')).toBeEnabled();
    vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue({ ok: true, value: read });
    fireEvent.click(screen.getByTestId('saved-location-retry'));
    await waitFor(() => expect(view.router.state.location.search).toContain('group_id=G'));
    expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
  });
  it('locks ordinary new input during readback, and does not erase a newer queued edit when the read arrives late', async () => {
    let finish!: (result: StudioLocatorEntryResult) => void;
    vi.mocked(readStudioLocatorEditorEntry).mockImplementation(() => new Promise(resolve => { finish = resolve; }));
    const view = mount(); await screen.findByTestId('canonical-owner'); fireEvent.click(screen.getByText('Accepted original Save'));
    await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledOnce());
    expect(screen.getByLabelText('Recording name')).toBeDisabled();
    await act(async () => { captured.lateInput?.(); });
    await act(async () => finish({ ok: true, value: read }));
    expect(view.router.state.location.search).toBe('?view=recording&intent=new');
    expect(screen.getByLabelText('Recording name')).toHaveValue('Newer unsaved input');
  });
  it('never navigates on an old readback after the host is left', async () => {
    let finish!: (result: StudioLocatorEntryResult) => void;
    vi.mocked(readStudioLocatorEditorEntry).mockImplementation(() => new Promise(resolve => { finish = resolve; }));
    const view = mount(); await screen.findByTestId('canonical-owner'); fireEvent.click(screen.getByText('Accepted original Save'));
    await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledOnce());
    const signal = vi.mocked(readStudioLocatorEditorEntry).mock.calls[0][2]!; view.unmount();
    expect(signal.aborted).toBe(true); await act(async () => finish({ ok: true, value: read }));
    expect(view.router.state.location.search).toBe('?view=recording&intent=new');
  });
});

it('revokes an in-flight saved-location read when the workspace changes', async () => {
  let finish!: (result: StudioLocatorEntryResult) => void;
  vi.mocked(readStudioLocatorEditorEntry).mockImplementation(() => new Promise(resolve => { finish = resolve; }));
  const view = mount(); await screen.findByTestId('canonical-owner'); fireEvent.click(screen.getByText('Accepted original Save'));
  await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledOnce());
  const signal = vi.mocked(readStudioLocatorEditorEntry).mock.calls[0][2]!;
  await act(async () => { captured.switchScope?.('OTHER'); }); expect(signal.aborted).toBe(true);
  await act(async () => finish({ ok: true, value: read }));
  expect(view.router.state.location.search).toBe('?view=recording&intent=new');
});
it('a newer accepted Save revokes the old read and cannot be replaced by its late location', async () => {
  let finishOld!: (result: StudioLocatorEntryResult) => void;
  vi.mocked(readStudioLocatorEditorEntry).mockImplementationOnce(() => new Promise(resolve => { finishOld = resolve; }));
  const view = mount(); await screen.findByTestId('canonical-owner'); fireEvent.click(screen.getByText('Accepted original Save'));
  await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledOnce());
  const oldSignal = vi.mocked(readStudioLocatorEditorEntry).mock.calls[0][2]!;
  const group = { ...read.entry.write_groups!.groups[0], revision: 'G4' };
  const next = { ...read, entry: { ...read.entry, workspace_revision: 'W3', group_heads: { G: 'G4' },
    write_groups: { workspace_id: 'W', workspace_revision: 'W3', groups: [group] } },
    selection: { state: 'selected' as const, target: { workspaceId: 'W', workspaceRevision: 'W3', step: 4 as const, groupId: 'G', groupRevision: 'G4' } } };
  captured.nextAck = { group, workspace_revision: 'W3' };
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue(next.entry.write_groups);
  vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue({ ok: true, value: next });
  fireEvent.click(screen.getByText('Accepted original Save'));
  await waitFor(() => expect(view.router.state.location.search).toContain('group_revision=G4'));
  expect(oldSignal.aborted).toBe(true);
  await act(async () => finishOld({ ok: true, value: read }));
  expect(view.router.state.location.search).toContain('group_revision=G4');
});
