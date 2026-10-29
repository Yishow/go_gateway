import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { Step4SupportPanels } from '@/features/datalink/workbench-v2/steps/step4/Step4SupportPanels';
import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import type { WriteGroupListResponse } from '@/types/studioV2WriteGroup';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, options?: { name?: string }) => options?.name ? `${key}: ${options.name}` : key }) }));
vi.mock('@/services/studioV2RuntimeContext', () => ({ studioV2RuntimeContextAPI: { getRuleEvidence: vi.fn() } }));
const groups = { data: { workspace_id: 'ws-1', workspace_revision: 'WR1', groups: [] } as WriteGroupListResponse | undefined, isLoading: false, isError: false, refetch: vi.fn() };
const mutate = vi.fn();
const measurements = { data: [], isLoading: false, isFetching: false, isError: false, refetch: vi.fn() };
vi.mock('@/hooks/datalink/useStudioV2WriteGroups', () => ({
  useWriteGroupsQuery: () => groups, useWriteGroupDeliveryQuery: () => ({ data: undefined, isError: false }),
  useUpdateWriteGroupMutation: () => ({ mutateAsync: mutate, isPending: false }),
  useEnsureCircuitProfileMutation: () => ({ mutateAsync: mutate, isPending: false }),
}));
vi.mock('@/hooks/datalink/useStudioV2RecordingStart', () => ({
  useEnsureBasicManagedMutation: () => ({ mutateAsync: mutate, isPending: false }),
  useRecordingStartMutation: () => ({ mutateAsync: mutate, isPending: false }),
  useRecordingStartOperationQuery: () => ({ data: undefined, isError: false }),
}));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceMeasurements', () => ({
  useStudioV2WorkspaceMeasurementsQuery: () => measurements,
  useCreateStudioV2MeasurementMutation: () => ({ mutateAsync: mutate, isPending: false }),
  useUpdateStudioV2MeasurementMutation: () => ({ mutateAsync: mutate, isPending: false }),
}));
vi.mock('@/features/datalink/workbench-v2/steps/step4/RecordingTestWriteRecoveryPanel', () => ({ RecordingTestWriteRecoveryPanel: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection', () => ({ WriteGroupSection: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/SavedRecordingManagement', () => ({ SavedRecordingManagement: () => null, RemovedRecordingManagement: () => null }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupSchemaPanel', () => ({ GroupSchemaPanel: () => null }));
const read = vi.mocked(studioV2RuntimeContextAPI.getRuleEvidence);
const clients: QueryClient[] = [];
function currentGroups() {
  if (!groups.data) throw new Error('Expected initialized test fixture');
  return groups.data;
}
function show(state = savedGroupState(), hasWorkspace = true) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  clients.push(client);
  const tree = () => <QueryClientProvider client={client}><Step4Database state={state} workspaceId={hasWorkspace ? 'ws-1' : undefined} dispatch={mutate} /></QueryClientProvider>;
  const view = render(tree());
  return { ...view, rerenderHost: () => view.rerender(tree()) };
}
beforeEach(() => {
  vi.clearAllMocks(); sessionStorage.clear();
  groups.data = { workspace_id: 'ws-1', workspace_revision: 'WR1', groups: [] };
  groups.isError = false; groups.isLoading = false;
  measurements.isError = false;
  read.mockResolvedValue({ ok: false, reason: 'unavailable', action: 'refresh' });
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });

describe('Step4 mounted guided summary bridge', () => {
  it('mounts a new-draft summary from actual Basic/Circuit controls and one readonly device query', async () => {
    show();
    expect(await screen.findByTestId('recording-scope-summary')).toBeInTheDocument();
    await waitFor(() => expect(read).toHaveBeenCalledTimes(1));
    expect(read.mock.calls[0][0]).toEqual({ workspace_id: 'ws-1', workspace_revision: 'WR1', device_id: 'dev-1' });
    fireEvent.change(screen.getByTestId('circuit-profile-interval-instantaneous'), { target: { value: '30' } });
    expect(screen.getByTestId('recording-scope-summary')).toHaveTextContent('30');
    const advanced = screen.getByTestId('step4-advanced-recording');
    fireEvent.click(within(advanced).getByText('management.advancedTitle'));
    fireEvent.click(screen.getByTestId('circuit-profile-purpose-wide'));
    expect(screen.getAllByTestId('recording-scope-summary')).toHaveLength(1);
    expect(read).toHaveBeenCalledTimes(1);
    expect(mutate).not.toHaveBeenCalled();
  });
  it('keeps the readonly provider mounted when legacy edit swaps the Circuit and Basic branches', async () => {
    currentGroups().groups = [savedGroup({ basic_managed_device_id: 'dev-1', layout_version: 'managed_wide_v1' })];
    show();
    await waitFor(() => expect(read).toHaveBeenCalledTimes(1));
    fireEvent.click(screen.getByTestId('device-recording-edit-existing'));
    expect(screen.getByTestId('basic-recording-panel')).toBeInTheDocument();
    expect(screen.getAllByTestId('recording-scope-summary')).toHaveLength(1);
    fireEvent.click(screen.getByTestId('device-recording-back-setup'));
    expect(screen.getByTestId('circuit-profile-basic-panel')).toBeInTheDocument();
    expect(screen.getAllByTestId('recording-scope-summary')).toHaveLength(1);
    expect(read).toHaveBeenCalledTimes(1);
    expect(mutate).not.toHaveBeenCalled();
  });
  it('keeps an empty selected device, reports no sources, and does not borrow the other device', async () => {
    const state = savedGroupState();
    state.devices.push({ ...state.devices[0], id: 'dev-empty', name: 'Empty device' });
    show(state);
    await waitFor(() => expect(read).toHaveBeenCalledTimes(1));
    fireEvent.click(screen.getByRole('button', { name: 'Empty device' }));
    await waitFor(() => expect(read).toHaveBeenCalledTimes(2));
    expect(read.mock.calls[1][0].device_id).toBe('dev-empty');
    expect(screen.getByTestId('recording-scope-summary')).toHaveTextContent('Empty device');
    const next = within(screen.getByTestId('recording-guided-preparation')).getByRole('button', { name: 'next.actions.inspect' });
    expect(next).toHaveAttribute('data-primary', 'true');
    fireEvent.click(next);
    expect(screen.queryByTestId('circuit-profile-prepare')).not.toBeInTheDocument();
    expect(read).toHaveBeenCalledTimes(2);
    expect(mutate).not.toHaveBeenCalled();
  });
  it('keeps ordinary non-counter controls usable without measurement evidence or any extra semantic mutation', async () => {
    measurements.isError = true;
    show();
    expect(await screen.findByTestId('recording-scope-summary')).toBeInTheDocument();
    expect(screen.queryByTestId('circuit-profile-measurements-error')).not.toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-prepare')).not.toBeDisabled();
    expect(screen.getByTestId('circuit-profile-toggle-instantaneous')).toBeChecked();
    expect(screen.getByTestId('circuit-profile-toggle-cumulative')).not.toBeChecked();
    expect(screen.getByTestId('circuit-profile-recreate-missing')).toBeChecked();
    expect(mutate).not.toHaveBeenCalled();
  });
  it.each(['foreign', 'error', 'loading'] as const)('does not query runtime from %s configuration evidence', async mode => {
    if (mode === 'foreign') currentGroups().workspace_id = 'other';
    if (mode === 'error') groups.isError = true;
    if (mode === 'loading') groups.isLoading = true;
    show();
    await waitFor(() => expect(screen.getByTestId('step4-basic-recording')).toBeInTheDocument());
    expect(read).not.toHaveBeenCalled();
    expect(mutate).not.toHaveBeenCalled();
  });
});

it('uses a safe connector alias in collapsed destination summary, without changing the intentional database editor', () => {
  const connector = { ...savedGroupState().db.connector, database: '/private/operator/archive.sqlite' };
  render(<Step4SupportPanels connector={connector} destinationReady onUpdateConnector={mutate} onKindChange={mutate} />);
  const card = screen.getByTestId('database-connection-card');
  expect(within(card).getByText('Line DB')).toBeVisible();
  expect(within(card).queryByText('/private/operator/archive.sqlite')).not.toBeInTheDocument();
  fireEvent.click(screen.getByTestId('database-connection-card-toggle'));
  expect(screen.getByDisplayValue('/private/operator/archive.sqlite')).toBeVisible();
  expect(mutate).not.toHaveBeenCalled();
});

it.each(['/private/operator/db', 'postgres://user:secret@db.example/private'])('does not derive a destination alias from unsafe name %s', name => {
  const connector = { ...savedGroupState().db.connector, name, database: '/private/operator/archive.sqlite' };
  render(<Step4SupportPanels connector={connector} destinationReady onUpdateConnector={mutate} onKindChange={mutate} />);
  const card = screen.getByTestId('database-connection-card');
  expect(within(card).queryByText(name)).not.toBeInTheDocument();
  expect(within(card).queryByText(connector.database)).not.toBeInTheDocument();
  expect(within(card).getAllByText('step4.connector_title').length).toBeGreaterThan(0);
  expect(mutate).not.toHaveBeenCalled();
});

it('does not expose a private-path alias in the legacy Basic destination summary', async () => {
  const state = savedGroupState();
  state.db.connector.name = '/private/operator/legacy.sqlite';
  currentGroups().groups = [savedGroup({ basic_managed_device_id: 'dev-1', layout_version: 'managed_wide_v1' })];
  show(state);
  fireEvent.click(screen.getByTestId('device-recording-edit-existing'));
  expect(screen.getByTestId('basic-recording-destination')).not.toHaveTextContent('/private/operator/legacy.sqlite');
  expect(mutate).not.toHaveBeenCalled();
});

it('keeps missing workspace and missing group data unavailable without throwing or querying runtime', () => {
  groups.data = undefined;
  show(savedGroupState(), false);
  expect(screen.getByTestId('recording-scope-summary')).toHaveTextContent('noScope');
  expect(read).not.toHaveBeenCalled();
  expect(mutate).not.toHaveBeenCalled();
});
it('withdraws readonly scope evidence if canonical group data disappears while Circuit remains mounted', async () => {
  const view = show();
  await waitFor(() => expect(read).toHaveBeenCalledTimes(1));
  groups.data = undefined;
  view.rerenderHost();
  expect(screen.getByTestId('circuit-profile-basic-panel')).toBeInTheDocument();
  expect(screen.getByTestId('recording-scope-summary')).toHaveTextContent('noScope');
  expect(read).toHaveBeenCalledTimes(1);
  expect(mutate).not.toHaveBeenCalled();
});
