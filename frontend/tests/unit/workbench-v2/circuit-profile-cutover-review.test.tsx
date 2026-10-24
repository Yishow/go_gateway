import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { CutoverPanel } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/CutoverPanel';
import { Step4Database } from '@/features/datalink/workbench-v2/steps/step4/Step4Database';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { quietDelivery, recordOf, reviewOf, target, wide } from '../../fixtures/cutoverFixtures';
import { savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({
  useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }),
}));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: {
  list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(),
  get: vi.fn(), schemaPreview: vi.fn(), schemaApply: vi.fn(), schemaOperation: vi.fn(), ensureCircuitProfile: vi.fn(),
  reviewCutover: vi.fn(), applyCutover: vi.fn(), getCutover: vi.fn(), listCutovers: vi.fn(),
} }));
vi.mock('@/services/studioV2WorkspaceMeasurements', () => ({ studioV2WorkspaceMeasurementsAPI: { list: vi.fn(async () => []) } }));
vi.mock('@/services/studioV2RecordingStart', () => ({ studioV2RecordingStartAPI: { start: vi.fn(), operation: vi.fn(), ensureBasic: vi.fn() } }));
const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);
const client = () => new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
function mount(groups: WriteGroup[] = [wide(), target()], host = false) {
  return render(<QueryClientProvider client={client()}>{host
    ? <Step4Database state={savedGroupState()} dispatch={vi.fn()} workspaceId="ws-1" />
    : <CutoverPanel groups={groups} workspaceId="ws-1" workspaceRevision="wrev-1" connectorRevision="crev-1" readonly={false} onOpenGroup={vi.fn()} />}
  </QueryClientProvider>);
}
async function review() {
  fireEvent.click(await screen.findByTestId('cutover-target-tgt-1'));
  fireEvent.change(screen.getByTestId('cutover-boundary'), { target: { value: '2099-01-01T00:00:00' } });
  fireEvent.click(screen.getByTestId('cutover-review-button'));
  await screen.findByTestId('cutover-review');
}
beforeEach(() => {
  vi.clearAllMocks();
  sessionStorage.clear();
  api.list.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'wrev-1', groups: [wide(), target()] });
  api.readiness.mockResolvedValue({ workspace_id: 'ws-1', workspace_revision: 'wrev-1', group_id: 'wide-1', group_revision: 'rev-wide', applied_revision: 'applied-wide', config_ready: true, schema_ready: true, ready: true, issues: [] });
  api.delivery.mockResolvedValue(quietDelivery as never);
  api.listCutovers.mockResolvedValue([]);
  api.reviewCutover.mockResolvedValue(reviewOf());
  api.applyCutover.mockResolvedValue(recordOf());
});

it('keeps source, mode, selection and boundary locked while the reviewed switch is actually being submitted', async () => {
  let finish!: (value: ReturnType<typeof recordOf>) => void;
  api.applyCutover.mockImplementation(() => new Promise((resolve) => { finish = resolve; }));
  mount(); await review();
  fireEvent.click(screen.getByTestId('cutover-confirm'));
  await waitFor(() => expect(api.applyCutover).toHaveBeenCalledTimes(1));
  try {
    for (const id of ['cutover-mode-forward', 'cutover-mode-rollback', 'cutover-source', 'cutover-target-tgt-1', 'cutover-boundary']) expect(screen.getByTestId(id), id).toBeDisabled();
  } finally { finish(recordOf()); }
});

it('does not discard an unknown frozen switch by allowing a new review before checking it', async () => {
  api.applyCutover.mockRejectedValue(Object.assign(new TypeError('Failed to fetch'), { outcome: 'unconfirmed' }));
  mount(); await review(); fireEvent.click(screen.getByTestId('cutover-confirm'));
  await screen.findByTestId('cutover-unknown');
  expect(screen.getByTestId('cutover-review-button')).toBeDisabled();
  expect(screen.getByTestId('cutover-mode-rollback')).toBeDisabled();
  fireEvent.click(screen.getByTestId('cutover-review-button'));
  expect(api.reviewCutover).toHaveBeenCalledTimes(1);
  expect(screen.getByTestId('cutover-unknown')).toBeInTheDocument();
});

it('uses the selected source group connector revision, even when the current connection controls show another destination', async () => {
  const old = wide();
  mount([wide({ destination: { ...old.destination, connector_id: 'conn-old', connector_revision: 'crev-old' } }), target()]);
  await review();
  expect(api.reviewCutover.mock.calls[0][1].expected_connector_revision).toBe('crev-old');
});

it('does not call a status-only proxy rejection proof that the switch failed', async () => {
  api.applyCutover.mockRejectedValue(Object.assign(new Error('Request failed with status code 409'), { response: { status: 409, data: '<html>proxy response</html>' } }));
  mount(); await review(); fireEvent.click(screen.getByTestId('cutover-confirm'));
  expect(await screen.findByTestId('cutover-unknown')).toBeInTheDocument();
  expect(screen.queryByTestId('cutover-failed')).not.toBeInTheDocument();
});

it('shows the actual moved and remaining source names in the host review, rather than only their counts', async () => {
  mount(undefined, true); await review();
  expect(screen.getByTestId('cutover-moved')).toHaveTextContent('line.temperature');
  expect(screen.getByTestId('cutover-moved')).toHaveTextContent('line.pressure');
  expect(screen.getByTestId('cutover-unselected')).toHaveTextContent('line.batch');
});

it('retains an unknown reviewed request across remount and requires checking before another mutation', async () => {
  api.applyCutover.mockRejectedValue(Object.assign(new TypeError('Failed to fetch'), { outcome: 'unconfirmed' }));
  const first = mount(); await review(); fireEvent.click(screen.getByTestId('cutover-confirm'));
  await screen.findByTestId('cutover-unknown');
  const frozen = structuredClone(api.applyCutover.mock.calls[0]);
  first.unmount();
  mount();
  expect(await screen.findByTestId('cutover-unknown')).toBeInTheDocument();
  expect(screen.getByTestId('cutover-review-button')).toBeDisabled();
  expect(api.applyCutover.mock.calls).toEqual([frozen]);
  expect(api.reviewCutover).toHaveBeenCalledTimes(1);
});

it.each(['proxy rejected the request', '<h1>upstream error</h1>'])('keeps an unvalidated proxy body unknown: %s', async (data) => {
  api.applyCutover.mockRejectedValue(Object.assign(new Error('Request failed'), { response: { status: 409, data } }));
  mount(); await review(); fireEvent.click(screen.getByTestId('cutover-confirm'));
  expect(await screen.findByTestId('cutover-unknown')).toBeInTheDocument();
  expect(screen.queryByTestId('cutover-failed')).not.toBeInTheDocument();
});

it('keeps a pending intent in its own workspace when the host changes workspace', async () => {
  api.applyCutover.mockRejectedValue(Object.assign(new TypeError('Failed to fetch'), { outcome: 'unconfirmed' }));
  const view = mount(); await review(); fireEvent.click(screen.getByTestId('cutover-confirm'));
  await screen.findByTestId('cutover-unknown');
  const otherGroups = [wide({ workspace_id: 'ws-2' }), target({ workspace_id: 'ws-2' })];
  view.rerender(<QueryClientProvider client={client()}><CutoverPanel groups={otherGroups} workspaceId="ws-2" workspaceRevision="wrev-2" connectorRevision="crev-2" readonly={false} onOpenGroup={vi.fn()} /></QueryClientProvider>);
  expect(screen.queryByTestId('cutover-unknown')).not.toBeInTheDocument();
  expect(screen.getByTestId('cutover-boundary')).toBeEnabled();
  view.rerender(<QueryClientProvider client={client()}><CutoverPanel groups={[wide(), target()]} workspaceId="ws-1" workspaceRevision="wrev-1" connectorRevision="crev-1" readonly={false} onOpenGroup={vi.fn()} /></QueryClientProvider>);
  expect(await screen.findByTestId('cutover-unknown')).toBeInTheDocument();
  expect(api.applyCutover).toHaveBeenCalledTimes(1);
});

it('still offers the read-only check if a pending source no longer appears in the editable group list', async () => {
  api.applyCutover.mockRejectedValue(Object.assign(new TypeError('Failed to fetch'), { outcome: 'unconfirmed' }));
  const first = mount(); await review(); fireEvent.click(screen.getByTestId('cutover-confirm'));
  await screen.findByTestId('cutover-unknown'); first.unmount();
  mount([target()]);
  expect(await screen.findByTestId('cutover-unknown')).toBeInTheDocument();
  expect(screen.getByTestId('cutover-check')).toBeEnabled();
  fireEvent.click(screen.getByTestId('cutover-check'));
  await screen.findByTestId('cutover-check-none');
  expect(api.applyCutover).toHaveBeenCalledTimes(1);
});
