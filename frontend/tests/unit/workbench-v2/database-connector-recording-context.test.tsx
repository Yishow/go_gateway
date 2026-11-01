import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2RecordingStartAPI } from '@/services/studioV2RecordingStart';
import { studioV2WriteGroupLifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
import { batchPanelState, longProfileGroup } from '../../fixtures/circuitBatchPanel';
import { guidedReadiness } from '../../fixtures/guidedRecordingPreparation';
import { savedLifecycle } from '../../fixtures/writeGroupLifecycle';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/studioV2WriteGroupLifecycle', () => ({ studioV2WriteGroupLifecycleAPI: { lifecycle: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { list: vi.fn(async () => []) } }));
vi.mock('@/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), operation: vi.fn(), ensureBasic: vi.fn() } }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), ensureCircuitProfile: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(),
  remove: vi.fn(), testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), get: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(),
} }));
const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const group = longProfileGroup({ id: 'group-a', circuitId: 'circuit-a', circuitName: 'Existing recording', interval: 15, points: [0, 1] });
beforeEach(() => {
  vi.clearAllMocks();
  sessionStorage.clear();
  api.list.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [group] });
  api.readiness.mockResolvedValue(guidedReadiness(group, 'wrev-1'));
  api.delivery.mockResolvedValue({ group_id: group.id, intake: { state: 'not_running' }, stages: { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 }, backlog: [], quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false } } as never);
  vi.mocked(studioV2WriteGroupLifecycleAPI.lifecycle).mockResolvedValue(savedLifecycle(group, 'wrev-1'));
});

it.each([false, true])('preserves the production database circuit draft readonly after save failure (Share configured: %s)', async (shareEnabled) => {
  let state = batchPanelState();
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const shareStatus = shareEnabled ? { workspace_id: 'ws-1', enabled: true, configured_enabled: true, port: 502, address: '127.0.0.1', bind_state: 'pass' as const, mapping_count: 1, readiness_token: 'share-token', settings_revision: 'settings-1', workspace_revision: 'share-wrev-1' } : undefined;
  const page = () => <QueryClientProvider client={client}><Step4Database state={state} dispatch={vi.fn()} workspaceId="ws-1" shareStatus={shareStatus} /></QueryClientProvider>;
  const view = render(page());
  await screen.findByTestId('circuit-profile-basic-panel');
  const interval = await screen.findByTestId('circuit-profile-interval-instantaneous');
  await waitFor(() => expect(interval).not.toBeDisabled());
  fireEvent.change(interval, { target: { value: '37' } });
  state = { ...state, db: { ...state.db, connector: { ...state.db.connector, name: 'Unsaved new name', save_state: 'save-error' } } };
  view.rerender(page());
  expect(screen.getByTestId('circuit-profile-basic-panel')).toBeVisible();
  expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toHaveValue('37');
  expect(screen.getByTestId('circuit-profile-interval-instantaneous')).toBeDisabled();
  expect(screen.queryByText('step4.basic.share_only_destination')).not.toBeInTheDocument();
  expect(api.ensureCircuitProfile).not.toHaveBeenCalled();
  expect(api.testWrite).not.toHaveBeenCalled();
  expect(studioV2RecordingStartAPI.start).not.toHaveBeenCalled();
});
