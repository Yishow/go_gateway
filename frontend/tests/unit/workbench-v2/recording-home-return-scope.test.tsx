import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, render, screen, waitFor } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { RecordingHomeRoute } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeRoute';
import { studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WorkspaceKeys } from '@/hooks/datalink/keys';
import { recordingHomeSourcesAPI } from '@/services/recordingHomeSources';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2Workspace', async original => ({ ...await original<typeof import('@/services/studioV2Workspace')>(),
  studioV2WorkspaceAPI: { get: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn() } }));
vi.mock('@/services/recordingHomeSources', () => ({ recordingHomeSourcesAPI: { read: vi.fn() } }));
vi.mock('@/hooks/datalink/useAcquisitionAnomalies', () => ({ useAcquisitionAnomalySummaryQuery: () => ({ isSuccess: false }) }));
let client: QueryClient;
afterEach(() => { cleanup(); client?.clear(); vi.restoreAllMocks(); });
it.each([1000, 6000])('revalidates home identity after a %i ms return before the list cache expires', async elapsed => {
  let now = Date.parse('2030-01-01T00:00:00Z'), workspaceId = 'A';
  vi.spyOn(Date, 'now').mockImplementation(() => now);
  vi.mocked(studioV2WorkspaceAPI.get).mockImplementation(async () => ({ id: workspaceId, kind: 'single', status: 'empty',
    ordered_device_ids: [], created_at: '', updated_at: '' }));
  vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mockImplementation(async () => ({ workspace_id: workspaceId, workspace_revision: `${workspaceId}1`, groups: [] }));
  vi.mocked(recordingHomeSourcesAPI.read).mockImplementation(async id => ({ workspace_id: id, state: 'available',
    source_digest: 'a'.repeat(64), devices: [], candidates: [], excluded: [] }));
  client = new QueryClient({ defaultOptions: { queries: { staleTime: 30_000, retry: false } } });
  const ui = (home: boolean) => <QueryClientProvider client={client}>{home && <RecordingHomeRoute location={{ search: '?view=recordings', hash: '' }} />}</QueryClientProvider>;
  const view = render(ui(true));
  await waitFor(() => expect(screen.getByTestId('home-source-observation')).toHaveAttribute('data-source-status', 'fresh'));
  view.rerender(ui(false));
  workspaceId = 'B'; now += elapsed;
  await act(async () => { view.rerender(ui(true)); });
  await waitFor(() => expect(recordingHomeSourcesAPI.read).toHaveBeenCalledWith('B', expect.anything(), expect.anything(), expect.anything()));
  await waitFor(() => expect(screen.getByTestId('home-source-observation')).toHaveAttribute('data-source-status', 'fresh'));
  expect(client.getQueryData(studioV2WorkspaceKeys.writeGroups())).toMatchObject({ workspace_id: 'B', workspace_revision: 'B1' });
  const reads = vi.mocked(studioV2WorkspaceWriteGroupsAPI.list).mock.calls.length;
  await act(async () => { view.rerender(ui(true)); });
  expect(studioV2WorkspaceWriteGroupsAPI.list).toHaveBeenCalledTimes(reads);
});
