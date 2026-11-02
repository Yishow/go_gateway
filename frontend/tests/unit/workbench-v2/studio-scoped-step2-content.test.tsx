import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { MemoryRouter, useLocation, useNavigate } from 'react-router-dom';
import DatalinkWorkbenchV2Page from '@/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2DatalinkApi, studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import { buildStudioEntityLocatorTarget, parseStudioEntityLocator } from '@/services/studioEntityLocator';
import type { SourceRuleRecord } from '@/types/datalink';
import i18n from '@/i18n/config';
import { locatorRead, nativeRule } from './studio-locator-fixtures';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
  studioV2WorkspaceAPI: { get: vi.fn() },
}));
const clients: QueryClient[] = [];
let revision = 'W2';
let list: () => Promise<SourceRuleRecord[]>;
function Host() {
  const location = useLocation(), navigate = useNavigate();
  return <><button onClick={() => navigate(buildStudioEntityLocatorTarget(locatorRead(1, 'B').locator)!)}>Select device B</button>
    <DatalinkWorkbenchV2Page navigateTo={navigate} locatorIntent={parseStudioEntityLocator(location.search)} /></>;
}
function mount() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } });
  clients.push(client);
  return render(<QueryClientProvider client={client}><MemoryRouter initialEntries={[buildStudioEntityLocatorTarget(locatorRead().locator)!]}><Host /></MemoryRouter></QueryClientProvider>);
}
function noWrites() {
  expect(studioV2DatalinkApi.post).not.toHaveBeenCalled();
  expect(studioV2DatalinkApi.put).not.toHaveBeenCalled();
  expect(studioV2DatalinkApi.delete).not.toHaveBeenCalled();
}
async function continueD() {
  mount();
  expect(await screen.findByTestId('editor-input-name')).toHaveValue('Device D');
  fireEvent.click(screen.getByTestId('btn-continue-step1'));
}
beforeEach(async () => {
  vi.resetAllMocks(); revision = 'W2'; list = async () => [nativeRule('R1'), nativeRule('R2'), nativeRule('RB', 'B')];
  await i18n.changeLanguage('en');
  vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(locatorRead().workspace);
  vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url, config) => {
    if (url === '/studio-v2/workspace/source-rules') {
      expect(config?.params).toBeUndefined();
      return { data: { success: true, data: await list() } };
    }
    if (url !== '/studio-v2/workspace') throw new Error(`Unexpected read: ${url}`);
    const params = config?.params as Record<string, string> | undefined;
    expect(params).toBeDefined();
    expect(params?.entry).toBe('locator'); expect(params?.workspace_id).toBe('W');
    expect(['D', 'B']).toContain(params?.device_id);
    expect(['1', '2']).toContain(params?.step);
    const step = params!.step === '2' ? 2 : 1;
    const read = locatorRead(step, params!.device_id);
    read.entry.workspace_revision = revision;
    if (step === 1) expect(params?.rule_id).toBeUndefined();
    else {
      expect(['R1', 'R2', 'RB']).toContain(params?.rule_id);
      const rule = nativeRule(params!.rule_id, params!.device_id);
      read.entry.source_rules = [rule]; read.entry.rule_heads = { [rule.id]: rule.revision_id! };
      read.entry.hydration_scope.rule_ids = [rule.id];
    }
    return { data: { success: true, data: { ...read.workspace, locator_entry: read.entry } } };
  });
});
afterEach(() => { cleanup(); clients.splice(0).forEach(client => client.clear()); });

describe('scoped advanced Step2 through the real Page', () => {
  it('lists only the selected device ranges and opens an explicitly chosen full native rule', async () => {
    await continueD();
    expect(await screen.findByTestId('studio-scoped-rule-R1')).toBeVisible();
    expect(screen.getByTestId('studio-scoped-rule-R2')).toBeVisible();
    expect(screen.queryByTestId('studio-scoped-rule-RB')).not.toBeInTheDocument();
    expect(screen.queryByTestId('rule-editor')).not.toBeInTheDocument(); noWrites();
    fireEvent.click(screen.getByTestId('studio-scoped-rule-R2'));
    expect(await screen.findByTestId('rule-start-input')).toHaveValue('40001');
    expect(screen.getByTestId('rule-device-select')).toHaveValue('D');
    expect(screen.getByTestId('rule-tab-R2')).toBeVisible();
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace', expect.objectContaining({ params: {
      entry: 'locator', workspace_id: 'W', step: '2', device_id: 'D', rule_id: 'R2',
    } })); noWrites();
  });
  it('supports an empty initial workspace revision with the original advanced Add and original create queue', async () => {
    revision = ''; list = async () => [];
    vi.mocked(studioV2DatalinkApi.post).mockImplementation(async (url, value) => {
      expect(url).toBe('/studio-v2/workspace/source-rules');
      const request = value as { id: string; device_id: string };
      expect(request.device_id).toBe('D');
      return { data: { success: true, data: { ...nativeRule(request.id, 'D'), ...request, runtime_apply_status: 'not_running' } } };
    });
    await continueD();
    expect(await screen.findByTestId('rule-add-btn')).toBeVisible(); noWrites();
    fireEvent.click(screen.getByTestId('rule-add-btn'));
    expect(await screen.findByTestId('rule-editor')).toBeVisible();
    expect(screen.getByTestId('rule-device-select')).toHaveValue('D');
    await waitFor(() => expect(studioV2DatalinkApi.post).toHaveBeenCalledOnce());
    expect(studioV2DatalinkApi.put).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.delete).not.toHaveBeenCalled();
  });
  it('does not turn a failed range read into empty scope or an Add form', async () => {
    list = async () => { throw new Error('Read unavailable'); };
    await continueD();
    const panel = await screen.findByTestId('studio-scoped-step2-list');
    await waitFor(() => expect(panel).toHaveAttribute('role', 'alert'));
    expect(screen.queryByTestId('rule-add-btn')).not.toBeInTheDocument();
    expect(screen.queryByTestId('rule-editor')).not.toBeInTheDocument(); noWrites();
  });
  it('rejects a changed workspace ownership bracket without offering a foreign or empty editor', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce(locatorRead().workspace)
      .mockResolvedValue({ ...locatorRead().workspace, id: 'FOREIGN' });
    await continueD();
    const panel = await screen.findByTestId('studio-scoped-step2-list');
    await waitFor(() => expect(panel).toHaveAttribute('role', 'alert'));
    expect(screen.queryByTestId('studio-scoped-rule-R1')).not.toBeInTheDocument();
    expect(screen.queryByTestId('rule-add-btn')).not.toBeInTheDocument(); noWrites();
  });
  it('cannot publish a late D range list into the subsequently selected B form', async () => {
    let resolve!: (rules: SourceRuleRecord[]) => void;
    list = () => new Promise(done => { resolve = done; });
    await continueD();
    await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/source-rules'));
    fireEvent.click(screen.getByRole('button', { name: 'Select device B' }));
    await waitFor(() => expect(screen.getByTestId('editor-input-name')).toHaveValue('Device B'));
    await act(async () => { resolve([nativeRule('R1')]); });
    expect(screen.getByTestId('editor-input-name')).toHaveValue('Device B');
    expect(screen.queryByTestId('studio-scoped-rule-R1')).not.toBeInTheDocument();
    expect(screen.queryByTestId('rule-editor')).not.toBeInTheDocument(); noWrites();
  });
});
