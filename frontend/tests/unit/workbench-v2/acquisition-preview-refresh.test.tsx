import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Step1Device } from '@/features/datalink/workbench-v2/steps/step1/Step1Device';
import { INITIAL_STATE } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import type { Device } from '@/features/datalink/workbench-v2/state/types';
import { act, fireEvent, render, renderHook, screen, waitFor } from '@testing-library/react';
import { afterEach, expect, it, vi } from 'vitest';
import { AcquisitionPlanPreviewPanel } from '@/features/datalink/workbench-v2/acquisition-preview/AcquisitionPlanPreviewPanel';
import { useAcquisitionPlanPreview } from '@/features/datalink/workbench-v2/acquisition-preview/useAcquisitionPlanPreview';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
afterEach(() => vi.resetAllMocks());
const scope = Object.freeze({ workspace_id: 'W', device_id: 'D', expected_workspace_revision: 'R' });
const groups = (workspace = 'W', revision = 'R2') => ({ data: { success: true, data: { workspace_id: workspace, workspace_revision: revision, groups: [] } } });
const preview = (device = 'D', revision = 'R2') => ({ data: { success: true, data: { workspace_id: 'W', device_id: device, workspace_revision: revision, status: 'unavailable', reason: 'installation_unavailable' } } });
const get = vi.mocked(studioV2DatalinkApi.get);
const stale = () => Promise.reject({ response: { status: 409 } });
async function refresh() { await screen.findByText('stale'); fireEvent.click(screen.getByRole('button', { name: 'refresh' })); }
it('refresh reads the canonical R2 head after schema advance, without changing the supplied R draft owner', async () => {
  // The caller's last acknowledged head remains R; schema-only R2 does not hydrate its draft owner.
  get.mockImplementation((url) => url.endsWith('/write-groups') ? Promise.resolve(groups()) : get.mock.calls.length === 1 ? stale() : Promise.resolve(preview()));
  render(<AcquisitionPlanPreviewPanel scope={scope} refreshWorkspaceHead />);
  await refresh();
  await waitFor(() => expect(get).toHaveBeenCalledTimes(3));
  expect(get.mock.calls[1][0]).toBe('/studio-v2/workspace/write-groups');
  expect(get.mock.calls[2][1]).toMatchObject({ params: { ...scope, expected_workspace_revision: 'R2' } });
  expect(await screen.findByText('unavailable')).toBeInTheDocument();
  expect(scope.expected_workspace_revision).toBe('R');
});
it.each(['OTHER', ''])('rejects foreign or malformed workspace head %s without preview retry', async workspace => {
  get.mockImplementation(url => url.endsWith('/write-groups') ? Promise.resolve(groups(workspace)) : stale());
  render(<AcquisitionPlanPreviewPanel scope={scope} refreshWorkspaceHead />); await refresh();
  await screen.findByText(workspace ? 'foreign' : 'unavailable'); expect(get).toHaveBeenCalledTimes(2);
});
it.each(['switch', 'suspend', 'unmount'] as const)('discards a late canonical read after %s', async action => {
  let resolve!: (value: ReturnType<typeof groups>) => void;
  const pending = new Promise<ReturnType<typeof groups>>(r => { resolve = r; });
  get.mockImplementation(url => url.endsWith('/write-groups') ? pending : stale());
  const view = render(<AcquisitionPlanPreviewPanel scope={scope} refreshWorkspaceHead />); await refresh();
  await waitFor(() => expect(get).toHaveBeenCalledTimes(2));
  if (action === 'unmount') view.unmount();
  else view.rerender(<AcquisitionPlanPreviewPanel scope={action === 'switch' ? { ...scope, device_id: 'D2' } : undefined} refreshWorkspaceHead />);
  const count = get.mock.calls.length;
  await act(async () => resolve(groups()));
  expect(get).toHaveBeenCalledTimes(count);
  expect((get.mock.calls[1][1] as { signal: AbortSignal }).signal.aborted).toBe(true);
});
it('keeps locator/default refresh on its exact supplied authority without global head lookup', async () => {
  get.mockImplementation(stale);
  render(<AcquisitionPlanPreviewPanel scope={scope} />); await refresh();
  await waitFor(() => expect(get).toHaveBeenCalledTimes(2));
  expect(get.mock.calls[1][0]).toBe('/studio-v2/workspace/acquisition-plan-preview');
  expect(get.mock.calls[1][1]).toMatchObject({ params: scope });
});

it('a newer explicit refresh cancels the older head read and only its R3 receipt can publish', async () => {
  let resolve!: (value: ReturnType<typeof groups>) => void;
  const older = new Promise<ReturnType<typeof groups>>(r => { resolve = r; });
  get.mockRejectedValueOnce({ response: { status: 409 } }).mockReturnValueOnce(older)
    .mockResolvedValueOnce(groups('W', 'R3')).mockResolvedValueOnce(preview('D', 'R3'));
  const { result } = renderHook(() => useAcquisitionPlanPreview(scope, true));
  await waitFor(() => expect(result.current.loading).toBe(false));
  act(() => result.current.refresh());
  await waitFor(() => expect(get).toHaveBeenCalledTimes(2));
  act(() => result.current.refresh());
  await waitFor(() => expect(get).toHaveBeenCalledTimes(4));
  await waitFor(() => expect(result.current.loading).toBe(false));
  await act(async () => resolve(groups()));
  expect(get).toHaveBeenCalledTimes(4);
  expect(get.mock.calls[3][1]).toMatchObject({ params: { ...scope, expected_workspace_revision: 'R3' } });
  expect(result.current.result).toMatchObject({ ok: true, value: { workspace_revision: 'R3' } });
});
it('does not retry preview if the canonical head read fails', async () => {
  get.mockImplementation(stale);
  render(<AcquisitionPlanPreviewPanel scope={scope} refreshWorkspaceHead />); await refresh();
  await screen.findByText('unavailable'); expect(get).toHaveBeenCalledTimes(2);
});
it('accepts an explicitly empty canonical legacy head, not a missing head', async () => {
  get.mockRejectedValueOnce({ response: { status: 409 } }).mockResolvedValueOnce(groups('W', '')).mockResolvedValueOnce(preview('D', ''));
  render(<AcquisitionPlanPreviewPanel scope={scope} refreshWorkspaceHead />); await refresh();
  await screen.findByText('unavailable');
  expect(get.mock.calls[2][1]).toMatchObject({ params: { ...scope, expected_workspace_revision: '' } });
});

it('the ordinary saved-device surface refresh leaves dirty connector form state and dispatch untouched', async () => {
  const device: Device = { id: 'D', name: 'PLC', description: '', protocol: 'modbus_tcp', config: {}, status: 'draft', test: null,
    persisted: true, save_state: 'saved', availability_status: 'available' };
  const state = { ...INITIAL_STATE, devices: [device], db: { ...INITIAL_STATE.db,
    connector: { ...INITIAL_STATE.db.connector, workspace_id: 'W', setup_revision: 'R', database: 'unsaved target', save_state: 'draft-invalid' as const },
  } };
  const before = JSON.stringify(state); const dispatch = vi.fn();
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  get.mockRejectedValueOnce({ response: { status: 409 } }).mockResolvedValueOnce(groups()).mockResolvedValueOnce(preview());
  render(<QueryClientProvider client={client}><Step1Device state={state} dispatch={dispatch} onContinue={vi.fn()}
    previewWorkspace={{ workspace_id: 'W', expected_workspace_revision: 'R' }} /></QueryClientProvider>);
  await refresh(); await screen.findByText('unavailable');
  expect(get.mock.calls[2][1]).toMatchObject({ params: { ...scope, expected_workspace_revision: 'R2' } });
  expect(dispatch).not.toHaveBeenCalled(); expect(JSON.stringify(state)).toBe(before);
  expect(client.getQueryCache().getAll()).toHaveLength(0);
});
