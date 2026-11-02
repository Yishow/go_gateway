import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { MemoryRouter, useLocation, useNavigate } from 'react-router-dom';
import DatalinkWorkbenchV2Page from '@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2DatalinkApi, studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { buildStudioEntityLocatorTarget, parseStudioEntityLocator } from '@/services/studioEntityLocator';
import i18n from '@/i18n/config';
import { locatorRead, nativeDevice } from './studio-locator-fixtures';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
  studioV2WorkspaceAPI: { get: vi.fn() },
}));
const clients: QueryClient[] = [];
let completeConnection = false;
function Host() {
  const location = useLocation(), navigate = useNavigate();
  return <DatalinkWorkbenchV2Page navigateTo={navigate} locatorIntent={parseStudioEntityLocator(location.search)} />;
}
function mount() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } });
  clients.push(client);
  return render(<QueryClientProvider client={client}><MemoryRouter initialEntries={[buildStudioEntityLocatorTarget(locatorRead().locator)!]}><Host /></MemoryRouter></QueryClientProvider>);
}
beforeEach(async () => {
  vi.resetAllMocks(); completeConnection = false;
  await i18n.changeLanguage('en');
  vi.mocked(studioV2WorkspaceAPI.get).mockRejectedValue(new Error('Unscoped bootstrap is forbidden'));
  vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url, config) => {
    if (url !== '/studio-v2/workspace') throw { response: { status: 503 } };
    const params = config?.params as Record<string, unknown> | undefined;
    const read = locatorRead(1, params?.device_id === 'B' ? 'B' : 'D');
    if (completeConnection) read.entry.devices!.forEach(device => { device.connection_config = JSON.stringify({ host: 'device.example', port: 502, slave_id: 1, timeout: 5 }); });
    return { data: { success: true, data: { ...read.workspace, locator_entry: read.entry } } };
  });
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); sessionStorage.removeItem('wbv2.fixture-recovery'); });

function expectNoMutation() {
  expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled();
  expect(studioV2DatalinkApi.post).not.toHaveBeenCalled();
  expect(studioV2DatalinkApi.put).not.toHaveBeenCalled();
  expect(studioV2DatalinkApi.delete).not.toHaveBeenCalled();
}
describe('scoped device source return through the real page and shell', () => {
  it('returns from an unread Settings to the same selected device form instead of looping its device list', async () => {
    mount();
    expect(await screen.findByTestId('editor-input-name')).toHaveValue('Device D');
    fireEvent.click(screen.getByTestId('settings-nav-button'));
    expect(await screen.findByTestId('studio-locator-notice')).toHaveTextContent('have not been read');
    fireEvent.click(screen.getByRole('link', { name: 'Open device D' }));
    await waitFor(() => expect(screen.getByTestId('editor-input-name')).toBeVisible());
    expect(screen.getByTestId('editor-input-name')).toHaveValue('Device D');
    expectNoMutation();
  });
  it('keeps the original invalid device edit visible and editable when another locator is refused', async () => {
    sessionStorage.setItem('wbv2.fixture-recovery', 'retained');
    mount();
    const name = await screen.findByTestId('editor-input-name');
    fireEvent.change(name, { target: { value: '' } });
    await waitFor(() => expect(screen.getByTestId('device-save-banner')).toHaveAttribute('data-save-state', 'draft-invalid'));
    fireEvent.click(screen.getByTestId('settings-nav-button'));
    expect(await screen.findByTestId('studio-locator-notice')).toHaveTextContent('have not been read');
    fireEvent.click(screen.getByRole('link', { name: 'Open device B' }));
    expect(await screen.findByText('Your existing edits are kept. Resolve them before opening a different settings selection.')).toBeVisible();
    expect(screen.getByTestId('editor-input-name')).toBeVisible();
    expect(screen.getByTestId('editor-input-name')).toHaveValue('');
    expect(screen.getByTestId('editor-input-name')).toBeEnabled();
    expect(sessionStorage.getItem('wbv2.fixture-recovery')).toBe('retained');
    expectNoMutation();
  });
  it('saves only the original D edit after a blocked foreign locator and explicit return to D', async () => {
    completeConnection = true;
    vi.mocked(studioV2DatalinkApi.put).mockImplementation(async (url, value) => {
      expect(url).toBe('/studio-v2/workspace/devices/D');
      const request = value as { name: string; connection_config: Record<string, unknown> };
      expect(request.name).toBe('Device D repaired');
      return { data: { success: true, data: { ...nativeDevice('D'), ...request, connection_config: JSON.stringify(Object.fromEntries(Object.entries(request.connection_config).reverse())), runtime_apply_status: 'not_running' } } };
    });
    mount();
    fireEvent.change(await screen.findByTestId('editor-input-name'), { target: { value: '' } });
    await waitFor(() => expect(screen.getByTestId('device-save-banner')).toHaveAttribute('data-save-state', 'draft-invalid'));
    fireEvent.click(screen.getByTestId('settings-nav-button'));
    fireEvent.click(await screen.findByRole('link', { name: 'Open device B' }));
    expect(await screen.findByText('Your existing edits are kept. Resolve them before opening a different settings selection.')).toBeVisible();
    expect(screen.getByTestId('editor-input-name')).toHaveValue('');
    fireEvent.click(screen.getByRole('link', { name: 'Open device D' }));
    await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenLastCalledWith('/studio-v2/workspace', expect.objectContaining({ params: expect.objectContaining({ device_id: 'D' }) })));
    await waitFor(() => expect(screen.getByTestId('editor-input-name')).toBeVisible());
    expect(screen.getByTestId('editor-input-name')).toHaveValue('');
    fireEvent.change(screen.getByTestId('editor-input-name'), { target: { value: 'Device D repaired' } });
    await waitFor(() => expect(screen.getByTestId('device-save-banner')).toHaveAttribute('data-save-state', 'saved'));
    expect(screen.getByTestId('editor-input-name')).toHaveValue('Device D repaired');
    expect(studioV2DatalinkApi.put).toHaveBeenCalledExactlyOnceWith('/studio-v2/workspace/devices/D', expect.objectContaining({ name: 'Device D repaired' }));
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled();
    expect(studioV2DatalinkApi.delete).not.toHaveBeenCalled();
    expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled();
    // Once the retained edit is genuinely saved, the original refusal's advice
    // must be actionable: selecting B again opens B without another mutation.
    fireEvent.click(screen.getByTestId('settings-nav-button'));
    fireEvent.click((await screen.findAllByRole('link', { name: 'Open device B' }))[0]);
    await waitFor(() => expect(screen.getByTestId('editor-input-name')).toHaveValue('Device B'));
    expect(screen.getByTestId('editor-input-name')).toBeVisible();
    expect(screen.queryByText('Your existing edits are kept. Resolve them before opening a different settings selection.')).not.toBeInTheDocument();
    expect(studioV2DatalinkApi.put).toHaveBeenCalledTimes(1);
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled();
    expect(studioV2DatalinkApi.delete).not.toHaveBeenCalled();
  });
});
