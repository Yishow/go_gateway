import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '../../../src/i18n/config';
import { RuntimeSetupContextPanel } from '../../../src/features/datalink/runtime-dashboard/RuntimeSetupContextPanel';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import type { EvidenceScope, RuleEvidenceView } from '../../../src/types/runtimeRuleEvidence';
import type { StudioV2RuntimeSetupMapping, StudioV2RuntimeSetupSourceRule } from '../../../src/types/studioV2RuntimeContext';
import { ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';

vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
const scope: EvidenceScope = { workspace_id: 'W', workspace_revision: 'W1', device_id: 'A' };
const rules = (device = 'A'): StudioV2RuntimeSetupSourceRule[] => Array.from({ length: device === 'A' ? 6 : 1 }, (_, i) => ({ id: `${device}${i + 1}`, revision_id: `R${i + 1}`, device_id: device, start_address: String(40001 + i), count: 1, data_type: 'int16', naming_prefix: `Rule ${device}${i + 1}`, enabled: true }));
const mappings = (device = 'A'): StudioV2RuntimeSetupMapping[] => Array.from({ length: device === 'A' ? 10 : 1 }, (_, i) => ({ id: `M${device}${i + 1}`, rule_id: `${device}${Math.min(i + 1, 6)}`, rule_revision: `R${Math.min(i + 1, 6)}`, point_id: `P${i + 1}`, device_id: device, address: String(40001 + i), tag_id: `T${i + 1}`, display_name: `Tag ${device}${i + 1}`, enabled: true }));
function detail(device = 'A', number = 6): RuleEvidenceView {
  const view = ruleEvidenceFixture({ ...scope, device_id: device }); view.groups = [];
  view.source = { rule_id: `${device}${Math.min(number, 6)}`, saved_revision: `R${Math.min(number, 6)}`, point_id: `P${number}`, mapping_id: `M${device}${number}`, tag_id: `T${number}`, effective_revision: null, state: 'saved-only', proof: 'unavailable', effective: null,
    saved: { binding: { workspace_id: 'W', workspace_revision: 'W1', entity_revision: `R${Math.min(number, 6)}` }, mapping_revision: 'MR', enabled: true, address: String(40000 + number), protocol_function: '03', read_type: 'int16', result_type: 'int16', polling: null, transforms: [] } };
  return view;
}
function response<Data>(data: Data) { return { data: { success: true, data } }; }
function requestParams(value: unknown): Record<string, unknown> {
  if (typeof value !== 'object' || value === null || Array.isArray(value)) throw new Error('Expected scoped request parameters');
  return value as Record<string, unknown>;
}
function page(params: Record<string, unknown>) {
  const device = String(params.device_id); const kind = params.source_list as 'rules' | 'mappings'; const query = String(params.q ?? '');
  const ruleRows = rules(device).filter((row) => !query || JSON.stringify(row).includes(query));
  const mappingRows = mappings(device).filter((row) => (!params.rule_id || row.rule_id === params.rule_id) && (!query || JSON.stringify(row).includes(query)));
  const offset = params.cursor === 'cGFnZTI' ? 5 : 0;
  const count = kind === 'rules' ? ruleRows.length : mappingRows.length;
  return response({ workspace_id: 'W', source_page: { contract_version: 1, workspace_id: 'W', workspace_revision: params.workspace_revision, device_id: device, kind, query, limit: 50, list_revision: 'L1', rules: kind === 'rules' ? ruleRows.slice(offset, offset + 5) : [], mappings: kind === 'mappings' ? mappingRows.slice(offset, offset + 5) : [], next_cursor: offset + 5 < count ? 'cGFnZTI' : null } });
}
const clients: QueryClient[] = [];
function harness(selected = scope, legacyRows = false) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  const props = (current: EvidenceScope) => ({ evidenceScope: current, selectedDevice: { device_id: current.device_id, name: current.device_id, protocol: 'modbus_tcp', running: false, availability_status: 'available' as const }, setupContext: { source_rules: legacyRows ? rules(current.device_id) : [], mappings: legacyRows ? mappings(current.device_id) : [], database_targets: [] } });
  const view = render(<QueryClientProvider client={client}><RuntimeSetupContextPanel {...props(selected)} /></QueryClientProvider>);
  return { ...view, client, switchTo(next: EvidenceScope) { view.rerender(<QueryClientProvider client={client}><RuntimeSetupContextPanel {...props(next)} /></QueryClientProvider>); } };
}
const calls = () => vi.mocked(studioV2DatalinkApi.get).mock.calls;
beforeEach(async () => {
  await i18n.changeLanguage('en'); vi.mocked(studioV2DatalinkApi.get).mockReset();
  vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (_url, config) => {
    const params = requestParams(config?.params);
    return params.source_list ? page(params) : response({ workspace_id: 'W', rule_evidence: detail(String(params.device_id), Number(String(params.point_id).slice(1))) });
  });
});
afterEach(() => { cleanup(); clients.splice(0).forEach((client) => client.clear()); });

describe('bounded runtime source lists', () => {
  it('makes all six rules and ten mappings inspectable without per-row fan-out', async () => {
    harness();
    const ruleList = await screen.findByRole('region', { name: 'Source rules' });
    const mappingList = screen.getByRole('region', { name: 'Source mappings' });
    await within(ruleList).findByRole('button', { name: 'Show mappings for A1' });
    await within(mappingList).findByRole('button', { name: 'Inspect source 40001' });
    expect(calls()).toHaveLength(2);
    for (let i = 1; i <= 5; i++) {
      expect(within(ruleList).getByRole('button', { name: `Show mappings for A${i}` })).toBeInTheDocument();
      expect(within(mappingList).getByRole('button', { name: `Inspect source ${40000 + i}` })).toBeInTheDocument();
    }
    fireEvent.click(within(ruleList).getByRole('button', { name: 'Next source rules page' }));
    await within(ruleList).findByRole('button', { name: 'Show mappings for A6' });
    expect(calls()).toHaveLength(3);
    fireEvent.click(within(mappingList).getByRole('button', { name: 'Next source mappings page' }));
    await within(mappingList).findByRole('button', { name: 'Inspect source 40010' });
    expect(calls()).toHaveLength(4);
    for (let i = 6; i <= 10; i++) expect(within(mappingList).getByRole('button', { name: `Inspect source ${40000 + i}` })).toBeInTheDocument();
    fireEvent.click(within(mappingList).getByRole('button', { name: 'Inspect source 40006' }));
    await screen.findByText('Saved-only configuration');
    expect(calls()).toHaveLength(5);
    expect(calls()[4][1]?.params).toEqual({ ...scope, rule_id: 'A6', rule_revision: 'R6', point_id: 'P6', mapping_id: 'MA6' });
    expect(screen.queryByText(/^\+\d+ more$/)).not.toBeInTheDocument();
  });
  it('submits bounded search without per-keystroke requests, then filters mappings by exact parent revision', async () => {
    harness();
    const list = await screen.findByRole('region', { name: 'Source rules' });
    await within(list).findByRole('button', { name: 'Show mappings for A1' });
    const input = within(list).getByRole('textbox', { name: 'Search source rules' });
    input.focus(); fireEvent.change(input, { target: { value: 'A6' } });
    expect(input).toHaveFocus(); expect(calls()).toHaveLength(2);
    fireEvent.submit(input.closest('form')!);
    await within(list).findByRole('button', { name: 'Show mappings for A6' });
    expect(calls()).toHaveLength(3);
    expect(calls()[2][1]?.params).toMatchObject({ source_list: 'rules', q: 'A6', limit: 50 });
    expect(calls()[2][1]?.params).not.toHaveProperty('cursor');
    fireEvent.click(within(list).getByRole('button', { name: 'Show mappings for A6' }));
    await screen.findByText('Mappings for A6');
    await screen.findByRole('button', { name: 'Inspect source 40010' });
    expect(calls()).toHaveLength(4);
    expect(calls()[3][1]?.params).toMatchObject({ source_list: 'mappings', rule_id: 'A6', rule_revision: 'R6' });
    expect(screen.queryByRole('region', { name: 'Source rule details' })).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Show all mappings' }));
    await screen.findByRole('button', { name: 'Inspect source 40001' });
    expect(calls()).toHaveLength(5);
  });
  it('supports explicit previous and next pages and resets cursor on search', async () => {
    harness();
    const list = await screen.findByRole('region', { name: 'Source rules' });
    await within(list).findByRole('button', { name: 'Show mappings for A1' });
    expect(within(list).getByRole('button', { name: 'Previous source rules page' })).toBeDisabled();
    fireEvent.click(within(list).getByRole('button', { name: 'Next source rules page' }));
    await within(list).findByRole('button', { name: 'Show mappings for A6' });
    expect(within(list).getByRole('button', { name: 'Next source rules page' })).toBeDisabled();
    fireEvent.click(within(list).getByRole('button', { name: 'Previous source rules page' }));
    await within(list).findByRole('button', { name: 'Show mappings for A1' });
    expect(calls()).toHaveLength(4);
    fireEvent.click(within(list).getByRole('button', { name: 'Next source rules page' }));
    await within(list).findByRole('button', { name: 'Show mappings for A6' });
    const search = within(list).getByRole('textbox'); fireEvent.change(search, { target: { value: 'A2' } }); fireEvent.submit(search.closest('form')!);
    await within(list).findByRole('button', { name: 'Show mappings for A2' });
    expect(calls().at(-1)?.[1]?.params).not.toHaveProperty('cursor');
    expect(within(list).getByText('Page 1')).toBeInTheDocument();
  });
  it('keeps selected source details while the same mapping page revalidates', async () => {
    const view = harness();
    const mappingList = await screen.findByRole('region', { name: 'Source mappings' });
    fireEvent.click(await within(mappingList).findByRole('button', { name: 'Inspect source 40001' }));
    const details = await screen.findByRole('region', { name: 'Source rule details' });
    await within(details).findByText('Saved-only configuration');

    let resolveMappings!: (result: ReturnType<typeof page>) => void;
    let mappingRequest: Record<string, unknown> | null = null;
    const pendingMappings = new Promise<ReturnType<typeof page>>((resolve) => { resolveMappings = resolve; });
    vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (_url, config) => {
      const params = requestParams(config?.params);
      if (params.source_list === 'mappings') {
        mappingRequest = params;
        return pendingMappings;
      }
      return params.source_list ? page(params) : response({ workspace_id: 'W', rule_evidence: detail(String(params.device_id), Number(String(params.point_id).slice(1))) });
    });
    let invalidation: Promise<void> | undefined;
    act(() => {
      invalidation = view.client.invalidateQueries({ predicate: (query) => query.queryKey.includes('source-page') });
    });
    await within(mappingList).findByRole('status');
    expect(screen.getByRole('region', { name: 'Source rule details' })).toBeInTheDocument();

    await act(async () => {
      resolveMappings(page(mappingRequest ?? { ...scope, source_list: 'mappings' }));
      await invalidation;
    });
    await within(mappingList).findByRole('button', { name: 'Inspect source 40001' });
    expect(screen.getByRole('region', { name: 'Source rule details' })).toBeInTheDocument();
  });
  it('clears selected source details after a settled page omits the selected mapping', async () => {
    const view = harness();
    const mappingList = await screen.findByRole('region', { name: 'Source mappings' });
    fireEvent.click(await within(mappingList).findByRole('button', { name: 'Inspect source 40001' }));
    const details = await screen.findByRole('region', { name: 'Source rule details' });
    await within(details).findByText('Saved-only configuration');

    let resolveMappings!: (result: ReturnType<typeof page>) => void;
    let mappingRequest: Record<string, unknown> | null = null;
    const pendingMappings = new Promise<ReturnType<typeof page>>((resolve) => { resolveMappings = resolve; });
    vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (_url, config) => {
      const params = requestParams(config?.params);
      if (params.source_list === 'mappings') {
        mappingRequest = params;
        return pendingMappings;
      }
      return params.source_list ? page(params) : response({ workspace_id: 'W', rule_evidence: detail(String(params.device_id), Number(String(params.point_id).slice(1))) });
    });
    let invalidation: Promise<void> | undefined;
    act(() => {
      invalidation = view.client.invalidateQueries({ predicate: (query) => query.queryKey.includes('source-page') });
    });
    await within(mappingList).findByRole('status');
    expect(screen.getByRole('region', { name: 'Source rule details' })).toBeInTheDocument();

    await act(async () => {
      const removedPage = page(mappingRequest ?? { ...scope, source_list: 'mappings' });
      removedPage.data.data.source_page.mappings = removedPage.data.data.source_page.mappings.filter((row) => row.point_id !== 'P1');
      resolveMappings(removedPage);
      await invalidation;
    });
    await waitFor(() => expect(screen.queryByRole('region', { name: 'Source rule details' })).not.toBeInTheDocument());
  });
  it('clears A6 detail on device switch and never replaces B1 with late A6', async () => {
    let finishA!: (result: ReturnType<typeof response>) => void;
    const pending = new Promise<ReturnType<typeof response>>((resolve) => { finishA = resolve; });
    vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (_url, config) => {
      const params = requestParams(config?.params);
      if (params.source_list) return page(params);
      return params.device_id === 'A' ? pending : response({ workspace_id: 'W', rule_evidence: detail('B', 1) });
    });
    const view = harness();
    const list = await screen.findByRole('region', { name: 'Source mappings' });
    await within(list).findByRole('button', { name: 'Inspect source 40001' });
    fireEvent.click(within(list).getByRole('button', { name: 'Next source mappings page' }));
    fireEvent.click(await within(list).findByRole('button', { name: 'Inspect source 40006' }));
    await waitFor(() => expect(calls()).toHaveLength(4));
    const signal = calls()[3][1]?.signal;
    view.switchTo({ ...scope, device_id: 'B' });
    expect(screen.queryByRole('region', { name: 'Source rule details' })).not.toBeInTheDocument();
    expect(signal?.aborted).toBe(true);
    await screen.findByRole('button', { name: 'Show mappings for B1' });
    fireEvent.click(await screen.findByRole('button', { name: 'Inspect source 40001' }));
    const region = await screen.findByRole('region', { name: 'Source rule details' });
    await within(region).findByText('B1');
    await act(async () => { finishA(response({ workspace_id: 'W', rule_evidence: detail() })); });
    expect(within(region).queryByText('A6')).not.toBeInTheDocument();
    expect(within(region).getByText('B1')).toBeInTheDocument();
    expect(calls()).toHaveLength(7);
    const bPages = calls().filter(([, config]) => { const params = requestParams(config?.params); return params.device_id === 'B' && params.source_list; });
    expect(bPages).toHaveLength(2);
    for (const [, config] of bPages) { expect(config?.params).not.toHaveProperty('cursor'); expect(config?.params).not.toHaveProperty('rule_id'); }
  });
  it('ignores late prior search pages and makes no source-detail calls', async () => {
    let finishOld!: (result: ReturnType<typeof response>) => void;
    const pending = new Promise<ReturnType<typeof response>>((resolve) => { finishOld = resolve; });
    vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (_url, config) => {
      const params = requestParams(config?.params);
      return params.q === 'A1' ? pending : page(params);
    });
    harness();
    const list = await screen.findByRole('region', { name: 'Source rules' });
    await within(list).findByRole('button', { name: 'Show mappings for A1' });
    const input = within(list).getByRole('textbox');
    fireEvent.change(input, { target: { value: 'A1' } }); fireEvent.submit(input.closest('form')!);
    await waitFor(() => expect(calls()).toHaveLength(3));
    const oldSignal = calls()[2][1]?.signal;
    fireEvent.change(input, { target: { value: 'A6' } }); fireEvent.submit(input.closest('form')!);
    await within(list).findByRole('button', { name: 'Show mappings for A6' });
    expect(oldSignal?.aborted).toBe(true);
    await act(async () => { finishOld(page({ ...scope, source_list: 'rules', q: 'A1' })); });
    expect(within(list).queryByRole('button', { name: 'Show mappings for A1' })).not.toBeInTheDocument();
    expect(calls()).toHaveLength(4);
    expect(calls().every(([, config]) => Boolean(requestParams(config?.params).source_list))).toBe(true);
  });
  it('shows stale page failure without silently restarting or loading a legacy fallback', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (_url, config) => {
      const params = requestParams(config?.params);
      if (params.cursor) throw { response: { data: { error: { code: 'stale-reference', message: '/private/secret' } } } };
      return page(params);
    });
    harness();
    const list = await screen.findByRole('region', { name: 'Source rules' });
    await within(list).findByRole('button', { name: 'Show mappings for A1' });
    fireEvent.click(within(list).getByRole('button', { name: 'Next source rules page' }));
    const alert = await within(list).findByRole('alert');
    expect(alert).toHaveTextContent('Source configuration changed');
    expect(alert).not.toHaveTextContent('private');
    expect(within(list).queryByRole('button', { name: 'Show mappings for A1' })).not.toBeInTheDocument();
    expect(calls()).toHaveLength(3);
    fireEvent.click(within(alert).getByRole('button', { name: 'Restart this list' }));
    await within(list).findByRole('button', { name: 'Show mappings for A1' });
    expect(calls()).toHaveLength(4);
  });
  it('localizes paging and search controls in Traditional Chinese', async () => {
    await i18n.changeLanguage('zh-TW'); harness();
    await screen.findByRole('button', { name: '檢視 A1 的映射' });
    expect(screen.getByRole('textbox', { name: '搜尋來源規則' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '下一頁來源規則' })).toBeEnabled();
  });

  it('does not populate unavailable page evidence from legacy saved arrays', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockRejectedValue({ response: { data: { error: { code: 'unavailable', message: 'password=secret' } } } });
    harness(scope, true);
    await screen.findAllByRole('alert');
    expect(screen.queryByRole('button', { name: 'Show mappings for A1' })).not.toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Inspect source 40001' })).not.toBeInTheDocument();
    expect(screen.queryByText('Tag A1')).not.toBeInTheDocument();
    expect(document.body.textContent).not.toContain('password=secret');
    expect(calls()).toHaveLength(2);
  });
  it('resets cursor and search on a new workspace revision and does not open a detail automatically', async () => {
    const view = harness();
    const list = await screen.findByRole('region', { name: 'Source rules' });
    await within(list).findByRole('button', { name: 'Show mappings for A1' });
    fireEvent.click(within(list).getByRole('button', { name: 'Next source rules page' }));
    await within(list).findByRole('button', { name: 'Show mappings for A6' });
    view.switchTo({ ...scope, workspace_revision: 'W2' });
    await screen.findByRole('button', { name: 'Show mappings for A1' });
    const refreshed = calls().filter(([, config]) => requestParams(config?.params).workspace_revision === 'W2');
    expect(refreshed).toHaveLength(2);
    for (const [, config] of refreshed) expect(config?.params).not.toHaveProperty('cursor');
    expect(screen.queryByRole('region', { name: 'Source rule details' })).not.toBeInTheDocument();
  });

  it('keeps safe rule counts and Tag keys visible alongside display names', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (_url, config) => {
      const result = page(requestParams(config?.params));
      const body = result.data.data.source_page;
      if (body.kind === 'rules') body.rules[0].count = 8;
      else body.mappings[0].tag_key = 'line01.temp.inlet';
      return result;
    });
    harness();
    const list = await screen.findByRole('region', { name: 'Source rules' });
    await within(list).findByRole('button', { name: 'Show mappings for A1' });
    expect(list).toHaveTextContent('8');
    expect(screen.getByRole('region', { name: 'Source mappings' })).toHaveTextContent('line01.temp.inlet');
    expect(screen.getByRole('region', { name: 'Source mappings' })).toHaveTextContent('Tag A1');
  });

});
