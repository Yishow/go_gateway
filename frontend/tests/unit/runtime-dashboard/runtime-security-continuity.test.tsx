import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '../../../src/i18n/config';
import { RuntimeDashboardPage } from '../../../src/features/datalink/runtime-dashboard/RuntimeDashboardPage';
import { createRuntimeEntityHistory, readRuntimeEntityHistory } from '../../../src/features/datalink/runtime-dashboard/runtimeEntityHistory';
import { studioV2RuntimeContextAPI } from '../../../src/services/studioV2RuntimeContext';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import type { RuntimeSourceSelection } from '../../../src/types/studioV2RuntimeContext';
import { evidenceContext, evidenceScope, ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';

vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() } }));
const hostile = { password: 'fixture-password', database: '/private/fixture.sqlite', host: '[2001:db8::1]:5432',
  dsn: 'Server=localhost;Database=metrics', raw_error: 'driver exploded: SELECT fixture FROM credentials' };
const leaked = /fixture-password|\/private\/fixture|2001:db8|Server=localhost|driver exploded|SELECT fixture/;
const selection: RuntimeSourceSelection = { rule_id: 'R', rule_revision: 'R1', point_id: 'P', mapping_id: 'M' };
const clients: QueryClient[] = [];
function fixture() {
  const rule = { id: 'R', revision_id: 'R1', device_id: 'D', start_address: 'BFF', count: 1, data_type: 'int16', naming_prefix: '電力（A）', enabled: true };
  const mapping = { id: 'M', rule_id: 'R', rule_revision: 'R1', device_id: 'D', point_id: 'P', address: 'BFF', tag_id: 'T', tag_key: 'plant/power', display_name: '入口溫度', unit: 'm/s', enabled: true };
  const nativeConnector = { id: 'C1', name: '安全儲存（A）', kind: 'sqlite', schema: 'main', table: 'meter_readings', write_mode: 'insert', write_interval_seconds: 5, status: 'ready', ...hostile };
  const context = { workspace_id: 'W', workspace_revision: 'W1', default_device_id: 'D', devices: [{ device_id: 'D', name: 'Pump A', protocol: 'mc_3e', running: false, availability_status: 'available', ...hostile }],
    setup: { source_rules: [{ ...rule, ...hostile }], mappings: [{ ...mapping, ...hostile }], database_targets: [], database_config: nativeConnector }, ...hostile };
  const group = ruleEvidenceFixture(); const effective = group.groups[0].effective;
  if (!effective) throw new Error('Expected effective configuration');
  effective.destination.alias = '安全儲存 A';
  group.groups[0].delivery = { current: { binding: { workspace_id: 'W', workspace_revision: 'W1', entity_revision: 'A' }, intake: null, totals_scope: 'group', shared: false,
    stages: { collecting: 2, queued: 3, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 }, last_commit: null, sources: [], source_state: 'unavailable',
    next_due_at: '2026-10-27T10:00:15Z', next_retry_at: null, attempt_duration_ms: null, quota: null, storage: null }, backlog: [], counter_replay_quarantines: [], oldest_pending_scope: 'group-all-revisions', oldest_pending_seconds: 12, test_write: null };
  const detail = structuredClone(group);
  detail.source = { rule_id: 'R', saved_revision: 'R1', effective_revision: null, point_id: 'P', mapping_id: 'M', tag_id: 'T', state: 'saved-only', proof: 'unavailable', effective: null,
    saved: { binding: { workspace_id: 'W', workspace_revision: 'W1', entity_revision: 'R1' }, mapping_revision: 'M1', enabled: true, address: 'BFF', protocol_function: '03', read_type: 'int16', result_type: 'float64',
      polling: { group_id: 'poll-P', interval_ms: 1000, source: 'polling-group' }, transforms: [{ kind: 'scale', order: null, multiplier: 0.1, offset: 2, target_type: null }] } };
  return { context, nativeConnector, rule, mapping, group, detail };
}
function record(value: unknown): Record<string, unknown> {
  if (!value || typeof value !== 'object' || Array.isArray(value)) throw new Error('Expected scoped parameters');
  return Object.fromEntries(Object.entries(value));
}
function response(data: unknown) { return { data: { success: true, data } }; }
function transport(data: ReturnType<typeof fixture>) {
  vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url, config) => {
    expect(url).toBe('/studio-v2/workspace/runtime-context');
    if (!config?.params) return response(data.context);
    const params = record(config.params);
    expect(params).toMatchObject(evidenceScope);
    if (params.source_list) return response({ workspace_id: 'W', source_page: { contract_version: 1, ...evidenceScope, kind: params.source_list, query: '', limit: 50,
      list_revision: 'L1', next_cursor: null, rules: params.source_list === 'rules' ? [data.rule] : [], mappings: params.source_list === 'mappings' ? [data.mapping] : [] } });
    return response(evidenceContext(params.point_id ? data.detail : data.group));
  });
}
async function mount(data = fixture()) {
  transport(data);
  const context = await studioV2RuntimeContextAPI.get();
  const setup = context.setup;
  if (!setup) throw new Error('Expected projected setup fixture');
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  const selected = vi.fn<(value: RuntimeSourceSelection | null) => void>();
  const navigate = vi.fn();
  const view = render(<QueryClientProvider client={client}><RuntimeDashboardPage routeState="loading" selectedDeviceId="D" selectedDevice={context.devices[0]}
    devices={context.devices} setupContext={setup} sourceEvidenceScope={evidenceScope} snapshot={null} snapshotError={null} liveValues={{}} streamState="idle"
    onSelectDevice={vi.fn()} onRetrySnapshot={async () => undefined} onReconnectStream={vi.fn()} onSourceSelection={selected} navigateTo={navigate} /></QueryClientProvider>);
  return { ...view, context, selected, navigate };
}
function readOnly() {
  expect(studioV2DatalinkApi.post).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.put).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.delete).not.toHaveBeenCalled();
  expect(JSON.stringify(vi.mocked(studioV2DatalinkApi.get).mock.calls)).not.toMatch(/\/logs|\/connectors|fixture-password|private\/fixture/);
}
beforeEach(async () => { await i18n.changeLanguage('en'); vi.resetAllMocks(); });
afterEach(() => { cleanup(); clients.splice(0).forEach((client) => client.clear()); });

describe('remote runtime service/parser to actual page safety continuity', () => {
  it('keeps hostile native connector fields out of source/list/group/delivery/timing DOM, settings URLs and identity-only history', async () => {
    const data = fixture(); const original = structuredClone(data); const mounted = await mount(data);
    const open = await screen.findByRole('button', { name: 'Inspect source BFF' });
    await screen.findByRole('region', { name: 'Current revision delivery' });
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(4); // context, 2 bounded lists, 1 group read
    fireEvent.click(open);
    const detail = await screen.findByRole('region', { name: 'Source rule details' });
    await within(detail).findByText('1000 ms — polling group poll-P');
    expect(within(detail).getByText('mean / 15 seconds')).toBeInTheDocument();
    expect(screen.getAllByText('安全儲存 A').length).toBeGreaterThan(0);
    expect(screen.getByText('入口溫度')).toBeInTheDocument();
    expect(screen.getAllByText('2026-10-27T10:00:15Z').length).toBeGreaterThan(0);
    expect(mounted.container.innerHTML).not.toMatch(leaked);
    const urls = [...mounted.container.querySelectorAll('a')].map((link) => link.getAttribute('href'));
    expect(urls).toEqual(expect.arrayContaining([
      '/studio/v2?entry=locator&workspace_id=W&step=1&device_id=D',
      '/studio/v2?entry=locator&workspace_id=W&step=2&device_id=D&rule_id=R',
      '/studio/v2?entry=locator&workspace_id=W&step=3&device_id=D&mapping_id=M',
      '/studio/v2?entry=locator&workspace_id=W&step=4&device_id=D&group_id=G&group_revision=A',
    ]));
    expect(JSON.stringify(urls)).not.toMatch(/meter_|poll-P|入口|安全|fixture/);
    expect(mounted.selected).toHaveBeenCalledWith(selection);
    const history = createRuntimeEntityHistory(evidenceScope, selection);
    expect(history).toEqual({ runtimeEntityHistory: { version: 1, ...evidenceScope, selection } });
    expect(readRuntimeEntityHistory(history, evidenceScope)).toEqual(selection);
    expect(readRuntimeEntityHistory(history, { ...evidenceScope, workspace_revision: 'W2' })).toBeNull();
    expect(createRuntimeEntityHistory(evidenceScope, { ...selection, ...hostile })).toBeNull();
    expect(JSON.stringify(mounted.context)).not.toMatch(leaked);
    expect(data).toEqual(original); // native authorized setup input was neither stripped nor mutated
    expect(data.nativeConnector.password).toBe('fixture-password');
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(5); readOnly();
  });
  it.each(['[2001:db8::1]:5432', 'host=localhost dbname=metrics'])('rejects private source query %s before any outgoing URL', async (query) => {
    expect(await studioV2RuntimeContextAPI.getSourcePage(evidenceScope, { kind: 'rules', query })).toEqual({ ok: false, reason: 'unavailable', action: 'refresh' });
    expect(studioV2DatalinkApi.get).not.toHaveBeenCalled(); readOnly();
  });
  it.each(['source-list', 'mapping-list', 'source-detail', 'destination', 'delivery', 'timing'] as const)('rejects hostile %s before DOM, attributes or settings links can reflect it', async (surface) => {
    const data = fixture();
    const source = data.detail.source?.saved;
    const destination = data.group.groups[0].effective?.destination;
    const current = data.group.groups[0].delivery?.current;
    if (!source || !destination || !current) throw new Error('Expected source, destination and current delivery fixtures');
    if (surface === 'source-list') Object.assign(data.rule, hostile);
    if (surface === 'mapping-list') Object.assign(data.mapping, hostile);
    if (surface === 'source-detail') Object.assign(source, { driver: hostile.raw_error });
    if (surface === 'destination') Object.assign(destination, hostile);
    if (surface === 'delivery') Object.assign(current, hostile);
    if (surface === 'timing') Object.assign(current, { next_due_at: hostile.database });
    const mounted = await mount(data);
    if (surface === 'source-detail') fireEvent.click(await screen.findByRole('button', { name: 'Inspect source BFF' }));
    await screen.findAllByRole('alert');
    expect(mounted.container.innerHTML).not.toMatch(leaked);
    if (['destination', 'delivery', 'timing'].includes(surface)) expect(screen.queryByRole('region', { name: 'Current revision delivery' })).not.toBeInTheDocument();
    readOnly();
  });
});
