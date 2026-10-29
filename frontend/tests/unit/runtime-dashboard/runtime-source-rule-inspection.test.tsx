import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '../../../src/i18n/config';
import { MemoryRouter } from 'react-router-dom';
import { RuntimeDashboardPage } from '../../../src/features/datalink/runtime-dashboard/RuntimeDashboardPage';
import { useRuntimeDashboardState } from '../../../src/features/datalink/runtime-dashboard/useRuntimeDashboardState';
import { RuntimeSetupContextPanel } from '../../../src/features/datalink/runtime-dashboard/RuntimeSetupContextPanel';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import type { EvidenceScope, RuleEvidenceView } from '../../../src/types/runtimeRuleEvidence';
import type { StudioV2RuntimeContext, StudioV2RuntimeContextDevice, StudioV2RuntimeSetupContext } from '../../../src/types/studioV2RuntimeContext';
import { ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';

const bridge = vi.hoisted(() => ({ context: null as StudioV2RuntimeContext | null }));
vi.mock('../../../src/hooks/datalink/useStudioV2RuntimeContext', () => ({ useStudioV2RuntimeContextQuery: () => ({ data: bridge.context, error: null, isLoading: false, isError: false }) }));
vi.mock('../../../src/hooks/datalink/usePoints', () => ({ usePointsQuery: () => ({ data: [] }) }));
vi.mock('../../../src/features/datalink/runtime-dashboard/useRuntimeStatus', () => ({ useRuntimeStatus: () => ({ data: undefined, error: null, refetch: vi.fn() }) }));
vi.mock('../../../src/features/datalink/runtime-dashboard/useRuntimeStream', () => ({ useRuntimeDashboardStream: () => ({ liveValues: {}, connectionState: 'idle', reconnect: vi.fn() }) }));
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
const scope: EvidenceScope = { workspace_id: 'W', workspace_revision: 'W1', device_id: 'D' };
const device = (id = 'D'): StudioV2RuntimeContextDevice => ({ device_id: id, name: id, protocol: 'modbus_tcp', running: true, availability_status: 'available', runtime_projection_version: 'R1' });
function setup(): StudioV2RuntimeSetupContext {
  return {
    source_rules: ['E', 'D'].map((id) => ({ id: `rule-${id}`, revision_id: 'R2', device_id: id, start_address: id === 'D' ? '40001' : '49999', count: 1, data_type: 'int16', naming_prefix: '', enabled: true })),
    mappings: ['E', 'D'].map((id) => ({ id: `map-${id}`, rule_id: `rule-${id}`, point_id: 'P', device_id: id, address: id === 'D' ? '40001' : '49999', tag_id: 'T', display_name: `Tag ${id}`, enabled: true })),
    database_targets: [],
  };
}
function evidence(selectedScope = scope): RuleEvidenceView {
  const view = ruleEvidenceFixture(selectedScope);
  view.source = {
    rule_id: `rule-${selectedScope.device_id}`, point_id: 'P', mapping_id: `map-${selectedScope.device_id}`, tag_id: 'T',
    saved_revision: 'R2', effective_revision: null, state: 'saved-only', proof: 'unavailable', effective: null,
    saved: { binding: { workspace_id: selectedScope.workspace_id, workspace_revision: selectedScope.workspace_revision, entity_revision: 'R2' }, mapping_revision: 'M2', enabled: true,
      address: '40001', protocol_function: '03', read_type: 'int16', result_type: 'float64', polling: { group_id: 'poll-P', interval_ms: 1000, source: 'polling-group' },
      transforms: [{ kind: 'byte-order', order: 'CDAB', multiplier: null, offset: null, target_type: null }, { kind: 'scale', order: null, multiplier: 0.1, offset: 2, target_type: null }, { kind: 'cast', order: null, multiplier: null, offset: null, target_type: 'float64' }] },
  };
  return view;
}
const response = (value = evidence()) => ({ data: { success: true, data: { workspace_id: value.workspace_id, rule_evidence: value } } });
const detailRequest = vi.fn<typeof studioV2DatalinkApi.get>();
let pageContext = setup();
const clients: QueryClient[] = [];
function panel(selectedScope: EvidenceScope | null = scope, context = setup()) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  pageContext = context;
  const props = { setupContext: context, selectedDevice: device(selectedScope?.device_id), evidenceScope: selectedScope };
  const view = render(<QueryClientProvider client={client}><RuntimeSetupContextPanel {...props} /></QueryClientProvider>);
  return { ...view, changeContext(next: StudioV2RuntimeSetupContext) { pageContext = next; void client.invalidateQueries({ predicate: (query) => query.queryKey.includes('source-page') }); view.rerender(<QueryClientProvider client={client}><RuntimeSetupContextPanel {...props} setupContext={next} /></QueryClientProvider>); }, changeScope(next: EvidenceScope) { view.rerender(<QueryClientProvider client={client}><RuntimeSetupContextPanel {...props} selectedDevice={device(next.device_id)} evidenceScope={next} /></QueryClientProvider>); } };
}
beforeEach(async () => {
  await i18n.changeLanguage('en'); pageContext = setup(); detailRequest.mockReset(); vi.mocked(studioV2DatalinkApi.get).mockReset();
  vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url, config) => {
    const params = config?.params as Record<string, string>;
    if (params.rule_id && params.point_id) return detailRequest(url, config);
    if (!params.source_list) return response(ruleEvidenceFixture({ workspace_id: params.workspace_id, workspace_revision: params.workspace_revision, device_id: params.device_id }));
    const rules = pageContext.source_rules.filter((row) => row.device_id === params.device_id);
    const mappings = pageContext.mappings.filter((row) => row.device_id === params.device_id).map((row) => ({ ...row, rule_revision: rules.find((rule) => rule.id === row.rule_id)?.revision_id }));
    return { data: { success: true, data: { workspace_id: params.workspace_id, source_page: {
      contract_version: 1, workspace_id: params.workspace_id, workspace_revision: params.workspace_revision, device_id: params.device_id,
      kind: params.source_list, query: '', limit: 50, list_revision: 'L1', next_cursor: null,
      rules: params.source_list === 'rules' ? rules : [], mappings: params.source_list === 'mappings' ? mappings : [],
    } } } };
  });
});
afterEach(() => { cleanup(); clients.splice(0).forEach((client) => client.clear()); });

describe('single selected-device source inspection', () => {
  it('filters source and mapping summaries before showing source details', async () => {
    panel();
    expect(screen.queryAllByText('49999')).toHaveLength(0);
    expect(screen.queryByText('Tag E')).not.toBeInTheDocument();
    expect(await screen.findByRole('button', { name: 'Inspect source 40001' })).toBeEnabled();
    expect(detailRequest).not.toHaveBeenCalled();
  });
  it('opens exact saved R2 scale and separates polling 1000 ms from recording mean / 15 seconds', async () => {
    detailRequest.mockResolvedValue(response());
    panel();
    const open = await screen.findByRole('button', { name: 'Inspect source 40001' });
    open.focus(); expect(open).toHaveFocus(); fireEvent.click(open);
    const detail = await screen.findByRole('region', { name: 'Source rule details' });
    await within(detail).findByText('Saved-only configuration');
    expect(within(detail).getByText('Effective source configuration unavailable')).toBeInTheDocument();
    expect(within(detail).getByText('1000 ms — polling group poll-P')).toBeInTheDocument();
    expect(within(detail).getByText('mean / 15 seconds')).toBeInTheDocument();
    expect(within(detail).getByText('CDAB')).toBeInTheDocument();
    expect(within(detail).getByText('× 0.1 + 2')).toBeInTheDocument();
    expect(within(detail).getByText('R2')).toBeInTheDocument();
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/runtime-context', expect.objectContaining({ params: { ...scope, rule_id: 'rule-D', rule_revision: 'R2', point_id: 'P', mapping_id: 'map-D' } }));
    const close = within(detail).getByRole('button', { name: 'Close source details' }); close.focus(); fireEvent.click(close);
    expect(screen.queryByRole('region', { name: 'Source rule details' })).not.toBeInTheDocument();
    expect(open).toHaveFocus();
  });
  it('does not read details when workspace or rule revision is unknown', async () => {
    const view = panel(null);
    expect(screen.queryByRole('button', { name: 'Inspect source 40001' })).not.toBeInTheDocument();
    view.unmount();
    const context = setup(); context.source_rules[1].revision_id = undefined;
    panel(scope, context);
    await screen.findAllByRole('alert');
    expect(screen.queryByRole('button', { name: 'Inspect source 40001' })).not.toBeInTheDocument();
    expect(detailRequest).not.toHaveBeenCalled();
  });
  it('uses an explicit point-only selection without silently selecting a mapping', async () => {
    const context = setup(); context.mappings[1].id = undefined;
    const value = evidence();
    if (!value.source?.saved) throw new Error('test source missing');
    value.source.mapping_id = null; value.source.tag_id = null;
    value.source.saved.mapping_revision = null; value.source.saved.result_type = null;
    value.source.saved.transforms = [];
    detailRequest.mockResolvedValue(response(value));
    panel(scope, context); fireEvent.click(await screen.findByRole('button', { name: 'Inspect source 40001' }));
    await screen.findByText('Saved-only configuration');
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/runtime-context', expect.objectContaining({ params: { ...scope, rule_id: 'rule-D', rule_revision: 'R2', point_id: 'P' } }));
    expect(screen.queryByText('mean / 15 seconds')).not.toBeInTheDocument();
  });
  it('does not use an unrelated applied member or saved draft as recording evidence', async () => {
    const value = evidence(); value.groups[0].effective!.members[0].point_id = 'other';
    detailRequest.mockResolvedValue(response(value));
    panel(); fireEvent.click(await screen.findByRole('button', { name: 'Inspect source 40001' }));
    await screen.findByText('Saved-only configuration');
    expect(screen.queryByText('mean / 15 seconds')).not.toBeInTheDocument();
    expect(screen.queryByText('mean / 60 seconds')).not.toBeInTheDocument();
  });
  it('shows unknown polling and unsupported transforms without inventing parameters', async () => {
    const value = evidence();
    if (!value.source?.saved) throw new Error('test source missing');
    value.source.saved.polling = null;
    value.source.saved.transforms = [{ kind: 'unsupported', order: null, multiplier: null, offset: null, target_type: null }];
    detailRequest.mockResolvedValue(response(value));
    panel(); fireEvent.click(await screen.findByRole('button', { name: 'Inspect source 40001' }));
    await screen.findByText('Unsupported transformation details');
    expect(screen.queryByText('1000 ms — polling group poll-P')).not.toBeInTheDocument();
    expect(screen.queryByText('× 0.1 + 2')).not.toBeInTheDocument();
    expect(screen.getByText('Effective source configuration unavailable')).toBeInTheDocument();
  });
  it.each(['unavailable', 'stale-reference', 'foreign-scope'])('renders safe %s failure and allows only a read-only refresh', async (code) => {
    detailRequest.mockRejectedValueOnce({ response: { data: { error: { code, message: 'password=unsafe /private/secret' } } } }).mockResolvedValueOnce(response());
    panel(); fireEvent.click(await screen.findByRole('button', { name: 'Inspect source 40001' }));
    const alert = await screen.findByRole('alert');
    expect(alert.textContent).not.toMatch(/password|private|secret/);
    fireEvent.click(within(alert).getByRole('button', { name: 'Refresh source evidence' }));
    await screen.findByText('Saved-only configuration');
    expect(detailRequest).toHaveBeenCalledTimes(2);
  });
  it('clears device A selection and ignores its late response after switching to B', async () => {
    let finishA!: (value: ReturnType<typeof response>) => void;
    const pendingA = new Promise<ReturnType<typeof response>>((resolve) => { finishA = resolve; });
    const nextScope = { ...scope, device_id: 'E' };
    const nextEvidence = evidence(nextScope); nextEvidence.source!.saved!.address = '49999';
    detailRequest.mockReturnValueOnce(pendingA).mockResolvedValueOnce(response(nextEvidence));
    const view = panel(); fireEvent.click(await screen.findByRole('button', { name: 'Inspect source 40001' }));
    await waitFor(() => expect(detailRequest).toHaveBeenCalledTimes(1));
    const signal = detailRequest.mock.calls[0][1]?.signal;
    view.changeScope(nextScope);
    expect(screen.queryByRole('region', { name: 'Source rule details' })).not.toBeInTheDocument();
    expect(signal?.aborted).toBe(true);
    fireEvent.click(await screen.findByRole('button', { name: 'Inspect source 49999' }));
    const detail = await screen.findByRole('region', { name: 'Source rule details' });
    await within(detail).findByText('rule-E');
    await act(async () => { finishA(response()); });
    expect(within(detail).queryByText('rule-D')).not.toBeInTheDocument();
    expect(within(detail).getByText('rule-E')).toBeInTheDocument();
    view.changeScope(scope);
    expect(screen.queryByRole('region', { name: 'Source rule details' })).not.toBeInTheDocument();
  });
  it('clears open detail after workspace revision changes without automatically fetching', async () => {
    detailRequest.mockResolvedValue(response());
    const view = panel(); fireEvent.click(await screen.findByRole('button', { name: 'Inspect source 40001' }));
    await screen.findByText('Saved-only configuration');
    view.changeScope({ ...scope, workspace_revision: 'W2' });
    expect(screen.queryByRole('region', { name: 'Source rule details' })).not.toBeInTheDocument();
    expect(detailRequest).toHaveBeenCalledTimes(1);
  });
  it('uses localized saved-only and unavailable evidence copy in Traditional Chinese', async () => {
    await i18n.changeLanguage('zh-TW');
    detailRequest.mockResolvedValue(response());
    panel(); fireEvent.click(await screen.findByRole('button', { name: '檢視來源 40001' }));
    await screen.findByText('僅已儲存設定');
    expect(screen.getByText('尚無生效來源設定證據')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '關閉來源詳情' })).toBeEnabled();
  });

  it('clears selected binding when it is removed from the observed summary', async () => {
    detailRequest.mockResolvedValue(response());
    const view = panel(); fireEvent.click(await screen.findByRole('button', { name: 'Inspect source 40001' }));
    await screen.findByText('Saved-only configuration');
    const changed = setup(); changed.mappings = [];
    view.changeContext(changed);
    await waitFor(() => expect(screen.queryByRole('region', { name: 'Source rule details' })).not.toBeInTheDocument());
    expect(detailRequest).toHaveBeenCalledTimes(1);
  });
  it('preserves missing-device legacy rendering without querying source details', () => {
    render(<RuntimeSetupContextPanel setupContext={setup()} selectedDevice={null} />);
    expect(screen.queryByText('Tag D')).not.toBeInTheDocument();
    expect(detailRequest).not.toHaveBeenCalled();
  });

  it('shows unavailable byte order even when a valid saved scale step exists', async () => {
    const value = evidence();
    value.source!.saved!.transforms = [{ kind: 'scale', order: null, multiplier: 0.1, offset: 2, target_type: null }];
    detailRequest.mockResolvedValue(response(value));
    panel(); fireEvent.click(await screen.findByRole('button', { name: 'Inspect source 40001' }));
    await screen.findByText('× 0.1 + 2');
    const detail = screen.getByRole('region', { name: 'Source rule details' });
    expect(within(detail).getByText('Byte order').parentElement).toHaveTextContent('Unavailable');
    expect(within(detail).queryByText('ABCD')).not.toBeInTheDocument();
  });

  it('passes observed workspace revision through the dashboard state and page without another authoring GET', async () => {
    bridge.context = { workspace_id: 'W', workspace_revision: 'W1', devices: [device()], default_device_id: 'D', setup: setup() };
    detailRequest.mockResolvedValue(response());
    const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
    function Dashboard() { return <RuntimeDashboardPage {...useRuntimeDashboardState()} />; }
    render(<QueryClientProvider client={client}><MemoryRouter><Dashboard /></MemoryRouter></QueryClientProvider>);
    fireEvent.click(await screen.findByRole('button', { name: 'Inspect source 40001' }));
    await screen.findByText('Saved-only configuration');
    expect(detailRequest).toHaveBeenCalledTimes(1);
    expect(detailRequest).toHaveBeenCalledWith('/studio-v2/workspace/runtime-context', expect.objectContaining({ params: { ...scope, rule_id: 'rule-D', rule_revision: 'R2', point_id: 'P', mapping_id: 'map-D' } }));
    const calls = vi.mocked(studioV2DatalinkApi.get).mock.calls;
    const groupCalls = calls.filter(([, config]) => {
      const params = config?.params as Record<string, string>;
      return !params.source_list && !params.rule_id && !params.point_id;
    });
    expect(groupCalls).toHaveLength(1);
    expect(groupCalls[0]).toEqual(['/studio-v2/workspace/runtime-context', { params: scope, signal: expect.any(AbortSignal) }]);
    expect(calls.every(([url]) => url === '/studio-v2/workspace/runtime-context')).toBe(true);
  });

});
