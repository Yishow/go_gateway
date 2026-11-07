import { recordingHomeSourcesAPI } from '@/services/recordingHomeSources';
import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import { listRecordingTestWriteRecovery } from '@/services/recordingTestWriteRecovery';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import DatalinkWorkbenchV2Page from '@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WorkspaceDevicesAPI } from '@/services/studioV2WorkspaceDevices';
import { studioV2RulesAPI } from '@/services/studioV2Rules';
import { studioV2MappingsAPI } from '@/services/studioV2Mappings';
import { pointAPI } from '@/services/datalink';
import { useStudioV2AutosaveState } from '@/pages/datalink/workbench-v2/useStudioV2AutosaveState';
vi.mock('@/pages/datalink/workbench-v2/useStudioV2AutosaveState', () => ({ useStudioV2AutosaveState: vi.fn(() => { throw new Error('Home must not mount autosave owners'); }) }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2Workspace', async original => ({
  ...await original<typeof import('@/services/studioV2Workspace')>(), studioV2WorkspaceAPI: { get: vi.fn() },
}));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  list: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), schemaApply: vi.fn(), remove: vi.fn(),
} }));
vi.mock('@/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: { list: vi.fn(), create: vi.fn() } }));
vi.mock('@/services/studioV2Rules', () => ({ studioV2RulesAPI: { list: vi.fn(), create: vi.fn() } }));
vi.mock('@/services/studioV2Mappings', () => ({ studioV2MappingsAPI: { list: vi.fn(), create: vi.fn() } }));
vi.mock('@/services/datalink', () => ({ pointAPI: { list: vi.fn() } }));
vi.mock('@/services/studioV2RuntimeContext', () => ({ studioV2RuntimeContextAPI: { getRuleEvidence: vi.fn() } }));
vi.mock('@/services/recordingTestWriteRecovery', () => ({ listRecordingTestWriteRecovery: vi.fn() }));
vi.mock('@/services/recordingHomeSources', () => ({ recordingHomeSourcesAPI: { read: vi.fn() } }));
const emptySources = (workspace_id: string) => ({ workspace_id, state: 'available' as const, source_digest: 'a'.repeat(64), devices: [], candidates: [], excluded: [] });
const workspace = { id: 'W', kind: 'single' as const, status: 'empty' as const, ordered_device_ids: [], created_at: '', updated_at: '' };
const clients: QueryClient[] = [];
beforeEach(() => {
  vi.clearAllMocks();
  vi.mocked(recordingHomeSourcesAPI.read).mockReset().mockImplementation(async id => emptySources(id));
  vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence).mockReset().mockResolvedValue({ ok: false, reason: 'unavailable', action: 'refresh' });
  vi.mocked(listRecordingTestWriteRecovery).mockReset().mockResolvedValue({ workspace_id: 'W1', complete: true, operations: [] });
  vi.mocked(studioV2WorkspaceAPI.get).mockReset().mockResolvedValue(workspace);
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockReset().mockResolvedValue({ workspace_id: 'W', workspace_revision: '', groups: [] });
  vi.mocked(studioV2WorkspaceDevicesAPI.list).mockReset().mockResolvedValue([]); vi.mocked(studioV2RulesAPI.list).mockReset().mockResolvedValue([]);
  vi.mocked(studioV2MappingsAPI.list).mockReset().mockResolvedValue([]); vi.mocked(pointAPI.list).mockReset().mockResolvedValue([]);
});
afterEach(() => { cleanup(); clients.forEach(client => client.clear()); clients.length = 0; });
function mount() {
  const client = new QueryClient(); clients.push(client); const navigate = vi.fn();
  render(<QueryClientProvider client={client}><DatalinkWorkbenchV2Page presentation="recordings" locatorIntent={{ kind: 'absent' }}
    presentationLocation={{ search: '?view=recordings&section=database', hash: '#saved' }} navigateTo={navigate} /></QueryClientProvider>);
  return navigate;
}

import { group, delivery as deliveryFixture } from './recording-home-guidance-fixtures';
import { studioV2WriteGroupLifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: vi.fn() } }));

// 30 saved non-draft recordings: the first page (25) requests details; G025–G029 are never requested.
// G003's delivery read fails on this page; G004's lifecycle and delivery reads both fail.
it('unconfirmed list omits unrequested off-page items while a page read error stays visible through filters', async () => {
  const groups = Array.from({ length: 30 }, (_, index) => group(`G${String(index).padStart(3, '0')}`));
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ ...workspace, id: 'W1' });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'W1', workspace_revision: 'W2', groups });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.delivery).mockImplementation(async id => {
    const found = groups.find(item => item.id === id); if (!found || id === 'G003' || id === 'G004') throw new Error('private backend');
    return deliveryFixture(found);
  });
  vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockImplementation(async id => {
    const found = groups.find(item => item.id === id); if (!found || id === 'G004') throw new Error('private backend'); return {
      workspace_id: 'W1', workspace_revision: 'W2', group_id: id, group_revision: found.revision, applied_revision: found.applied_revision,
      status: found.status, connector_id: 'DB1', connector_revision: 'DB2', actions: [], history: { canonical: true, current_group_id: id, generations: [] },
    };
  });
  mount();
  const notice = await screen.findByRole('region', { name: 'unconfirmed.title' });
  await waitFor(() => expect(within(notice).getByText('Recording G003')).toBeInTheDocument());
  expect(within(within(notice).getByText('Recording G003').closest('li')!).getByText('unconfirmed.read.delivery')).toBeInTheDocument();
  await waitFor(() => expect(within(within(notice).getByText('Recording G004').closest('li')!).getByText('unconfirmed.read.lifecycle')).toBeInTheDocument());
  for (let index = 25; index < 30; index++) expect(notice).not.toHaveTextContent(`Recording G0${index}`);
  fireEvent.change(screen.getByLabelText('keyword'), { target: { value: 'Recording G010' } });
  await waitFor(() => expect(within(screen.getByTestId('home-management-table')).queryByText('Recording G003')).not.toBeInTheDocument());
  const filtered = screen.getByRole('region', { name: 'unconfirmed.title' });
  expect(within(filtered).getByText('Recording G003')).toBeInTheDocument();
  // Both reads actually failed in this visit; leaving the page must not turn them into "never requested".
  const failed = within(filtered).getByText('Recording G004').closest('li')!;
  expect(within(failed).getByText('unconfirmed.read.lifecycle')).toBeInTheDocument();
  expect(within(failed).getByText('unconfirmed.read.delivery')).toBeInTheDocument();
  for (let index = 25; index < 30; index++) expect(filtered).not.toHaveTextContent(`Recording G0${index}`);
  const requested = new Set([...vi.mocked(studioV2WorkspaceWriteGroupsAPI.delivery).mock.calls, ...vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mock.calls].map(([id]) => id));
  for (let index = 25; index < 30; index++) expect(requested.has(`G0${index}`)).toBe(false);
  expect(document.body.textContent).not.toContain('private backend');
  for (const action of ['create', 'update', 'apply', 'schemaApply', 'remove'] as const) expect(studioV2WorkspaceWriteGroupsAPI[action]).not.toHaveBeenCalled();
  expect(useStudioV2AutosaveState).not.toHaveBeenCalled();
});

function mockDetails(groups: ReturnType<typeof group>[], failing: readonly string[] = []) {
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({ ...workspace, id: 'W1' });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockResolvedValue({ workspace_id: 'W1', workspace_revision: 'W2', groups });
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.delivery).mockImplementation(async id => {
    const found = groups.find(item => item.id === id); if (!found || failing.includes(id)) throw new Error('private backend'); return deliveryFixture(found);
  });
  vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockImplementation(async id => {
    const found = groups.find(item => item.id === id); if (!found) throw new Error('fixture'); return {
      workspace_id: 'W1', workspace_revision: 'W2', group_id: id, group_revision: found.revision, applied_revision: found.applied_revision,
      status: found.status, connector_id: 'DB1', connector_revision: 'DB2', actions: [], history: { canonical: true, current_group_id: id, generations: [] },
    };
  });
}
const ids = (from: number, to: number) => Array.from({ length: to - from }, (_, index) => `G0${String(from + index).padStart(2, '0')}`);
it('a page of drafts keeps the draft guidance when only never-requested off-page recordings are unknown', async () => {
  mockDetails([...ids(0, 25).map(id => group(id, 'draft')), ...ids(25, 30).map(id => group(id))]);
  mount();
  await waitFor(() => expect(studioV2WorkspaceWriteGroupsAPI.delivery).toHaveBeenCalledTimes(25));
  await waitFor(() => expect(screen.getByTestId('recording-home-guidance')).toHaveAttribute('data-guidance-state', 'neutral'));
  expect(within(screen.getByTestId('recording-home-guidance')).getByText('reason.draft')).toBeInTheDocument();
  expect(screen.queryByText('reason.evidence_unconfirmed')).not.toBeInTheDocument();
  expect(screen.queryByRole('region', { name: 'unconfirmed.title' })).not.toBeInTheDocument();
  const requested = new Set(vi.mocked(studioV2WorkspaceWriteGroupsAPI.delivery).mock.calls.map(([id]) => id));
  for (const id of ids(25, 30)) expect(requested.has(id)).toBe(false);
});
it('an actual read error on the draft page still leads and is named', async () => {
  mockDetails([...ids(0, 24).map(id => group(id, 'draft')), group('G024'), ...ids(25, 30).map(id => group(id))], ['G024']);
  mount();
  const notice = await screen.findByRole('region', { name: 'unconfirmed.title' });
  await waitFor(() => expect(within(within(notice).getByText('Recording G024').closest('li')!).getByText('unconfirmed.read.delivery')).toBeInTheDocument());
  expect(within(screen.getByTestId('recording-home-guidance')).getByText('reason.evidence_unconfirmed')).toBeInTheDocument();
  for (const id of ids(25, 30)) expect(notice).not.toHaveTextContent(`Recording ${id}`);
});
