import { StrictMode } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { WriteGroupSection, type WriteGroupSectionProps } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection';
import type { GroupEditorProps } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { readStudioLocatorEditorEntry } from '@/services/studioV2LocatorEntry';
import { studioV2LocatorHydration } from '@/pages/datalink/workbench-v2/studioV2LocatorHydration';
import { locatorRead } from './studio-locator-fixtures';
import { savedGroupState } from '../../fixtures/writeGroupState';
import type { StudioLocatorEntryResult } from '@/types/studioV2LocatorEntry';
const captured = vi.hoisted(() => ({ props: undefined as GroupEditorProps | undefined }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor', () => ({ GroupEditor: (props: GroupEditorProps) => {
  captured.props = props; return <div data-testid="canonical-editor"><input aria-label="Local recording name" defaultValue="" />{props.destinationContent}<button onClick={props.onClose}>Close canonical editor</button></div>;
} }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/SavedRecordingCard', () => ({ SavedRecordingCard: ({ onOpen }: { onOpen: () => void }) => <button onClick={onOpen}>Open saved card</button> }));
vi.mock('@/services/studioV2LocatorEntry', () => ({ readStudioLocatorEditorEntry: vi.fn() }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn() } }));
const clients: QueryClient[] = [];
function mount(props: Partial<WriteGroupSectionProps> = {}, strict = false) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } }); clients.push(client);
  const state = savedGroupState(); state.devices = []; state.points = []; state.mappings = {};
  state.db.connector = { ...state.db.connector, persisted: false, save_state: 'idle', connector_id: undefined, identity_revision: undefined };
  const content = (next: Partial<WriteGroupSectionProps>) => <QueryClientProvider client={client}><WriteGroupSection state={state} workspaceId="ws-1" readonly={false} presentation="simple" initialSelection="new" {...next} /></QueryClientProvider>;
  return { client, ...render(strict ? <StrictMode>{content(props)}</StrictMode> : content(props)), content };
}
beforeEach(() => { vi.resetAllMocks(); captured.props = undefined; vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: '', groups: [] }); });
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });
describe('simple presentation delegates to the one canonical editor owner', () => {
  it('opens a local draft with no source or saved DB and injects a single destination owner', async () => {
    const onClose = vi.fn(); mount({ onClose, destinationContent: <div data-testid="original-db-owner" /> });
    expect(await screen.findAllByTestId('canonical-editor')).toHaveLength(1);
    expect(screen.getAllByTestId('original-db-owner')).toHaveLength(1);
    expect(captured.props).toMatchObject({ presentation: 'simple', group: undefined, workspaceId: 'ws-1', workspaceRevision: '', candidates: [] });
    expect(captured.props?.connector.persisted).toBe(false);
    expect(screen.queryByTestId('group-create')).not.toBeInTheDocument();
    expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
    fireEvent.click(screen.getByText('Close canonical editor')); expect(onClose).toHaveBeenCalledOnce();
  });
  it.each(['error', 'foreign', 'missing-revision'])('does not treat %s scope as a new empty workspace', async mode => {
    if (mode === 'error') vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockRejectedValue(new Error('private failure'));
    else vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: mode === 'foreign' ? 'other' : 'ws-1', ...(mode === 'foreign' ? { workspace_revision: '' } : {}), groups: [] } as never);
    mount(); await screen.findByTestId('group-section-error');
    expect(screen.queryByTestId('canonical-editor')).not.toBeInTheDocument();
    expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
  });
  it('auto-opens saved intent only after the existing fresh scoped read', async () => {
    const entry = locatorRead(4); let finish!: (result: StudioLocatorEntryResult) => void;
    vi.mocked(readStudioLocatorEditorEntry).mockImplementation(() => new Promise(resolve => { finish = resolve; }));
    mount({ state: studioV2LocatorHydration(entry)!, workspaceId: 'W', locatorEntry: entry, locatorOnly: true, initialSelection: 'located' });
    await waitFor(() => expect(readStudioLocatorEditorEntry).toHaveBeenCalledOnce());
    expect(readStudioLocatorEditorEntry).toHaveBeenCalledWith(entry.locator, { workspaceRevision: 'W2', groupRevision: 'G3' }, expect.any(AbortSignal));
    expect(screen.queryByTestId('canonical-editor')).not.toBeInTheDocument();
    await act(async () => finish({ ok: true, value: { ...entry, entry: { ...entry.entry, editor_hydrated: true, database_config: null, database_targets: [] } } }));
    expect(await screen.findAllByTestId('canonical-editor')).toHaveLength(1);
    expect(captured.props).toMatchObject({ group: { id: 'G', revision: 'G3' }, workspaceId: 'W', workspaceRevision: 'W2', locatorGuard: true });
    expect(studioV2WorkspaceWriteGroupsAPI.list).not.toHaveBeenCalled();
    expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
  });
  it('a failed saved read never falls back to a local new draft', async () => {
    const entry = locatorRead(4); vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue({ ok: false, reason: 'stale-reference' });
    mount({ state: studioV2LocatorHydration(entry)!, workspaceId: 'W', locatorEntry: entry, locatorOnly: true, initialSelection: 'located' });
    await screen.findByTestId('studio-locator-open-error');
    expect(screen.queryByTestId('canonical-editor')).not.toBeInTheDocument();
    expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
  });
});

it('StrictMode abandons the first scope read and still opens exactly one fresh editor', async () => {
  const entry = locatorRead(4);
  vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue({ ok: true, value: { ...entry, entry: { ...entry.entry, editor_hydrated: true, database_config: null, database_targets: [] } } });
  mount({ state: studioV2LocatorHydration(entry)!, workspaceId: 'W', locatorEntry: entry, locatorOnly: true, initialSelection: 'located' }, true);
  expect(await screen.findAllByTestId('canonical-editor')).toHaveLength(1);
  expect(vi.mocked(readStudioLocatorEditorEntry).mock.calls[0][2]?.aborted).toBe(true);
  expect(captured.props?.group?.id).toBe('G');
});

it('collapses only fresh simple locator diagnostics while opening and guided diagnostics stay visible', async () => {
  const entry = locatorRead(4);let finish!: (result: StudioLocatorEntryResult) => void;
  vi.mocked(readStudioLocatorEditorEntry).mockImplementation(() => new Promise(resolve => { finish = resolve; }));
  const props: Partial<WriteGroupSectionProps> = { state: studioV2LocatorHydration(entry)!, workspaceId: 'W', locatorEntry: entry, locatorOnly: true, initialSelection: 'located' };
  const view = mount(props);
  expect(await screen.findByText('locatorEntry.loading')).toBeVisible();expect(screen.getByText('locatorEntry.unrequested')).toBeVisible();
  expect(screen.queryByTestId('locator-version-details')).not.toBeInTheDocument();
  await act(async () => finish({ ok: true, value: { ...entry, entry: { ...entry.entry, editor_hydrated: true, database_config: null, database_targets: [] } } }));
  await screen.findByTestId('canonical-editor');const details = screen.getByTestId('locator-version-details');
  expect(details.tagName).toBe('DETAILS');expect(details).not.toHaveAttribute('open');expect(details).toHaveTextContent('versions.details');
  for (const text of ['locatorEntry.unrequested', 'locatorEntry.historicalHint']) {
    expect(screen.getByText(text).closest('details')).toBe(details);expect(screen.getByText(text)).not.toBeVisible();
  }
  expect(screen.getByText(/locatorEntry.currentRevision/).closest('details')).toBe(details);
  view.rerender(view.content({ ...props, presentation: 'guided' }));
  expect(screen.queryByTestId('locator-version-details')).not.toBeInTheDocument();expect(screen.getByText('locatorEntry.unrequested')).toBeVisible();
  expect(screen.getByText('locatorEntry.historicalHint')).toBeVisible();
});
it.each(['stale-reference', 'foreign-scope', 'unavailable'] as const)('keeps %s locator diagnostics and denial outside any collapsed version details', async reason => {
  const entry = locatorRead(4);vi.mocked(readStudioLocatorEditorEntry).mockResolvedValue({ ok: false, reason });
  mount({ state: studioV2LocatorHydration(entry)!, workspaceId: 'W', locatorEntry: entry, locatorOnly: true, initialSelection: 'located' });
  const error = await screen.findByTestId('studio-locator-open-error');expect(error).toBeVisible();expect(error.closest('details')).toBeNull();
  expect(screen.getByText('locatorEntry.unrequested')).toBeVisible();expect(screen.queryByTestId('locator-version-details')).not.toBeInTheDocument();
});

it.each(['foreign', 'error'] as const)('keeps the same unsaved simple draft but denies writes after a %s refresh', async mode => {
  const { client, rerender, content } = mount();
  const input = await screen.findByLabelText('Local recording name');
  fireEvent.change(input, { target: { value: 'Unsaved local name' } });
  if (mode === 'foreign') vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'other', workspace_revision: 'X', groups: [] });
  else vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockRejectedValue(new Error('private failure'));
  await act(async () => { await client.refetchQueries(); });
  await waitFor(() => {
    expect(screen.getByLabelText('Local recording name')).toBe(input);
    expect(input).toHaveValue('Unsaved local name');
    expect(captured.props).toMatchObject({ workspaceId: 'ws-1', workspaceRevision: '', readonly: true });
  });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'W2', groups: [] });
  await act(async () => { captured.props?.onReload(); });
  await waitFor(() => expect(captured.props).toMatchObject({ workspaceId: 'ws-1', workspaceRevision: 'W2', readonly: false }));
  expect(screen.getByLabelText('Local recording name')).toBe(input); expect(input).toHaveValue('Unsaved local name');
  rerender(content({ workspaceId: 'different-scope' }));
  expect(screen.queryByTestId('canonical-editor')).not.toBeInTheDocument();
  expect(studioV2WorkspaceWriteGroupsAPI.create).not.toHaveBeenCalled();
});
