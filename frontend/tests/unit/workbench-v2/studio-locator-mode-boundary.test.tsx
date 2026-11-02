import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { MemoryRouter, useLocation, useNavigate } from 'react-router-dom';
import DatalinkWorkbenchV2Page from '@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2DatalinkApi, studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { buildStudioEntityLocatorTarget } from '@/services/studioEntityLocator';
import { parseStudioPresentation } from '@/services/studioRecordingView';
import i18n from '@/i18n/config';
import { locatorRead, nativeDevice, nativeRule } from './studio-locator-fixtures';

vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() }, studioV2WorkspaceAPI: { get: vi.fn() } }));
const clients: QueryClient[] = [];
let savedDeviceName = 'Device D';
let returnRead: 'ok' | 'error' | 'revision' = 'ok';
let deviceReads = 0;
let savedRule = nativeRule();
let foreignRuleHead = false;
function Host() {
  const location = useLocation(), navigate = useNavigate(), entry = parseStudioPresentation(location.search);
  const foreign = buildStudioEntityLocatorTarget(locatorRead(1, 'B').locator)!;
  return <><output data-testid="actual-url">{location.pathname + location.search}</output>
    <a href={foreign} onClick={event => { event.preventDefault(); navigate(foreign); }}>Select B</a>
    <button onClick={() => navigate(-1)}>History back</button><button onClick={() => navigate(1)}>History forward</button>
    <DatalinkWorkbenchV2Page navigateTo={navigate} locatorIntent={entry.locatorIntent} presentation={entry.presentation} presentationLocation={{ search: location.search, hash: location.hash }} /></>;
}
function mount(step: 1 | 2 = 1) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  render(<QueryClientProvider client={client}><MemoryRouter initialEntries={[buildStudioEntityLocatorTarget(locatorRead(step).locator)!]}><Host /></MemoryRouter></QueryClientProvider>);
}
beforeEach(async () => {
  vi.resetAllMocks(); savedDeviceName = 'Device D'; returnRead = 'ok'; deviceReads = 0; savedRule = nativeRule(); foreignRuleHead = false; await i18n.changeLanguage('en');
  vi.mocked(studioV2WorkspaceAPI.get).mockRejectedValue(new Error('Full bootstrap is allowed only after explicit handoff'));
  vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url, config) => {
    if (url !== '/studio-v2/workspace') throw new Error(`Unexpected read ${url}`);
    const params = config?.params as Record<string, string>;
    const read = locatorRead(params.step === '2' ? 2 : 1, params.device_id);
    if (params.device_id === 'D') {
      deviceReads += 1; read.entry.devices![0].name = savedDeviceName;
      if (deviceReads > 1 && returnRead === 'error') throw new Error('Read unavailable');
      if (deviceReads > 1 && returnRead === 'revision') read.entry.workspace_revision = 'W3';
    }
    if (params.step === '2') {
      const rule = { ...savedRule, ...(deviceReads > 1 && foreignRuleHead ? { revision_id: 'FOREIGN-REVISION' } : {}) };
      read.entry.source_rules = [rule]; read.entry.rule_heads = { R: rule.revision_id! };
    }
    read.entry.devices!.forEach(device => { device.connection_config = JSON.stringify({ host: 'device.example', port: 502, slave_id: 1, timeout: 5 }); });
    return { data: { success: true, data: { ...read.workspace, locator_entry: read.entry } } };
  });
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });
describe('precise locator versus full advanced navigation', () => {
  it('keeps a loaded rule editable without advertising unloaded steps or Continue to mapping', async () => {
    mount(2); expect(await screen.findByTestId('rule-start-input')).toBeVisible();
    expect(screen.queryByTestId('step-nav-button-3')).not.toBeInTheDocument();
    expect(screen.queryByText('STEP 02 / 04')).not.toBeInTheDocument();
    expect(screen.queryByTestId('continue-step3-btn')).not.toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Open full advanced setup' })).toBeEnabled();
    expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled();
    expect(studioV2DatalinkApi.put).not.toHaveBeenCalled();
  });
  it('opens the original complete advanced owner only after the explicit clean handoff', async () => {
    mount(); expect(await screen.findByTestId('editor-input-name')).toHaveValue('Device D');
    expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled();
    fireEvent.click(screen.getByRole('button', { name: 'Open full advanced setup' }));
    await waitFor(() => expect(screen.getByTestId('actual-url')).toHaveTextContent(/^\/studio\/v2\?view=advanced$/));
    await waitFor(() => expect(studioV2WorkspaceAPI.get).toHaveBeenCalled());
    expect(studioV2DatalinkApi.put).not.toHaveBeenCalled();
  });
  it('keeps an invalid local edit actionable and disables both exits without discarding it', async () => {
    mount(); fireEvent.change(await screen.findByTestId('editor-input-name'), { target: { value: '' } });
    await waitFor(() => expect(screen.getByTestId('device-save-banner')).toHaveAttribute('data-save-state', 'draft-invalid'));
    expect(screen.getByRole('button', { name: 'Open full advanced setup' })).toBeDisabled();
    expect(screen.getByRole('button', { name: 'Back to recordings' })).toBeDisabled();
    expect(screen.getByTestId('editor-input-name')).toBeEnabled();
    expect(screen.getByTestId('editor-input-name')).toHaveValue('');
    expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.put).not.toHaveBeenCalled();
  });
  it('replaces a refused B URL with the retained D owner and allows the same B anchor after D ACK', async () => {
    vi.mocked(studioV2DatalinkApi.put).mockImplementation(async (url, value) => {
      expect(url).toBe('/studio-v2/workspace/devices/D');
      const request = value as { name: string; connection_config: Record<string, unknown> };
      savedDeviceName = request.name;
      return { data: { success: true, data: { ...nativeDevice('D'), ...request, connection_config: JSON.stringify(request.connection_config), runtime_apply_status: 'not_running' } } };
    });
    mount(); fireEvent.change(await screen.findByTestId('editor-input-name'), { target: { value: '' } });
    await waitFor(() => expect(screen.getByTestId('device-save-banner')).toHaveAttribute('data-save-state', 'draft-invalid'));
    fireEvent.click(screen.getByRole('link', { name: 'Select B' }));
    await waitFor(() => expect(screen.getByTestId('actual-url')).toHaveTextContent(buildStudioEntityLocatorTarget(locatorRead().locator)!));
    await waitFor(() => expect(screen.getByTestId('editor-input-name')).toBeEnabled());
    expect(screen.getByTestId('editor-input-name').closest('fieldset')).not.toHaveAttribute('inert');
    expect(screen.getByTestId('editor-input-name')).toHaveValue('');
    fireEvent.change(screen.getByTestId('editor-input-name'), { target: { value: 'Device D repaired' } });
    await waitFor(() => expect(studioV2DatalinkApi.put).toHaveBeenCalledTimes(1));
    await waitFor(() => expect(screen.getByTestId('device-save-banner')).toHaveAttribute('data-save-state', 'saved'));
    await waitFor(() => expect(screen.getByRole('button', { name: 'Open full advanced setup' })).toBeEnabled());
    fireEvent.click(screen.getByRole('link', { name: 'Select B' }));
    await waitFor(() => expect(screen.getByTestId('editor-input-name')).toHaveValue('Device B'));
    expect(studioV2DatalinkApi.put).toHaveBeenCalledTimes(1);
    // Post-fix history control: the rejected B location was replaced, not pushed.
    fireEvent.click(screen.getByRole('button', { name: 'History back' }));
    await waitFor(() => expect(screen.getByTestId('editor-input-name')).toHaveValue('Device D repaired'));
    expect(screen.getByTestId('actual-url')).toHaveTextContent(buildStudioEntityLocatorTarget(locatorRead().locator)!);
    fireEvent.click(screen.getByRole('button', { name: 'History forward' }));
    await waitFor(() => expect(screen.getByTestId('editor-input-name')).toHaveValue('Device B'));
    expect(studioV2DatalinkApi.put).toHaveBeenCalledTimes(1);
    expect(deviceReads).toBe(3);
    expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled();
  });
  it.each([false, true])('uses only the acknowledged baseline rule revision on return (foreign=%s)', async foreign => {
    vi.mocked(studioV2DatalinkApi.put).mockImplementation(async (url, value) => {
      expect(url).toBe('/studio-v2/workspace/source-rules/R');
      savedRule = { ...savedRule, ...(value as Partial<typeof savedRule>), revision_id: `ACK-${vi.mocked(studioV2DatalinkApi.put).mock.calls.length}` };
      return { data: { success: true, data: { ...savedRule, runtime_apply_status: 'aligned' } } };
    });
    mount(2);
    fireEvent.change(await screen.findByTestId('rule-count-input'), { target: { value: '3' } });
    await waitFor(() => expect(screen.getByTestId('rule-save-status-R')).toHaveTextContent('Saved'));
    fireEvent.change(screen.getByTestId('rule-count-input'), { target: { value: '2' } });
    await waitFor(() => expect(studioV2DatalinkApi.put).toHaveBeenCalledTimes(2));
    await waitFor(() => expect(screen.getByTestId('rule-save-status-R')).toHaveTextContent('Saved'));
    fireEvent.change(screen.getByTestId('rule-name-input'), { target: { value: 'Unsaved local range label' } });
    foreignRuleHead = foreign;
    fireEvent.click(screen.getByRole('link', { name: 'Select B' }));
    await waitFor(() => expect(deviceReads).toBe(2));
    await waitFor(() => expect(screen.getByTestId('actual-url')).toHaveTextContent(buildStudioEntityLocatorTarget(locatorRead(2).locator)!));
    expect(screen.getByTestId('rule-name-input')).toHaveValue('Unsaved local range label');
    if (foreign) {
      await waitFor(() => expect(screen.getByTestId('rule-name-input')).toBeDisabled());
      expect(screen.getByRole('button', { name: 'Read settings again' })).toBeVisible();
    } else await waitFor(() => expect(screen.getByTestId('rule-name-input')).toBeEnabled());
    expect(screen.getByRole('button', { name: 'Open full advanced setup' })).toBeDisabled();
    expect(studioV2DatalinkApi.put).toHaveBeenCalledTimes(2);
    expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled();
  });
  it.each(['error', 'revision'] as const)('does not restore authority from a %s return read', async failure => {
    mount(); fireEvent.change(await screen.findByTestId('editor-input-name'), { target: { value: '' } });
    await waitFor(() => expect(screen.getByTestId('device-save-banner')).toHaveAttribute('data-save-state', 'draft-invalid'));
    returnRead = failure;
    fireEvent.click(screen.getByRole('link', { name: 'Select B' }));
    await waitFor(() => expect(deviceReads).toBe(2));
    await waitFor(() => expect(screen.getByRole('button', { name: 'Read settings again' })).toBeVisible());
    expect(screen.queryByRole('link', { name: 'Open device B' })).not.toBeInTheDocument();
    expect(screen.queryByTestId('studio-locator-notice')).not.toBeInTheDocument();
    expect(screen.getByTestId('actual-url')).toHaveTextContent(buildStudioEntityLocatorTarget(locatorRead().locator)!);
    expect(screen.getByTestId('editor-input-name')).toHaveValue('');
    expect(screen.getByRole('button', { name: 'Open full advanced setup' })).toBeDisabled();
    expect(screen.getByTestId('editor-input-name')).toBeDisabled();
    expect(screen.getByTestId('editor-input-name').closest('fieldset')).toHaveAttribute('inert');
    expect(studioV2DatalinkApi.put).not.toHaveBeenCalled();
    expect(studioV2WorkspaceAPI.get).not.toHaveBeenCalled();
    returnRead = 'ok';
    fireEvent.click(screen.getByRole('button', { name: 'Read settings again' }));
    await waitFor(() => expect(screen.getByTestId('editor-input-name')).toBeEnabled());
    expect(screen.getByTestId('editor-input-name')).toHaveValue('');
    expect(screen.getByRole('button', { name: 'Open full advanced setup' })).toBeDisabled();
    expect(studioV2DatalinkApi.put).not.toHaveBeenCalled();
  });
});
