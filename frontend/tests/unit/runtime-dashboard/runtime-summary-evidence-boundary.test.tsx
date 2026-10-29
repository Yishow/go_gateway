import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '../../../src/i18n/config';
import { RuntimeSummaryPanel } from '../../../src/features/datalink/runtime-dashboard/RuntimeSummaryPanel';
import { RuntimeDashboardPage, type RuntimeDashboardPageProps } from '../../../src/features/datalink/runtime-dashboard/RuntimeDashboardPage';
import { runtimeAPI } from '../../../src/services/datalink';
import { api } from '../../../src/services/datalinkClient';
import { projectRuntimeStatus } from '../../../src/services/runtimeStatusProjection';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import { evidenceContext, evidenceScope, ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() } }));
vi.mock('../../../src/services/datalinkClient', async (original) => ({
  ...(await original<typeof import('../../../src/services/datalinkClient')>()), api: { get: vi.fn() },
}));
const metrics = { collected_total: 100, write_success_total: 73, write_error_total: 5, mapping_error_total: 3, point_state_error_total: 2 };
const keys = Object.keys(metrics) as (keyof typeof metrics)[];
const clients: QueryClient[] = [];
function snapshot(observed: unknown = metrics) {
  return projectRuntimeStatus({ running: true, uptime_seconds: 65, metrics: observed,
    snapshot_state: { state: 'ready', empty: false, degraded: false, unavailable: false, stale: false },
    collectors: [{ device_id: 'D', device_name: 'Meter', protocol: 'modbus_tcp', availability_status: 'available', status: 'running', running: true, points_total: 1, points_healthy: 1, points_stale: 0, points_error: 0, last_read_at: '2026-10-29T10:00:00Z', last_error: null, breaker_state: 'closed' }] });
}
function groupEvidence() {
  const view = ruleEvidenceFixture();
  view.groups[0].delivery = { current: { binding: { workspace_id: 'W', workspace_revision: 'W1', entity_revision: 'A' }, intake: null, totals_scope: 'group', shared: false,
    stages: { collecting: 2, queued: 3, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 }, last_commit: null, sources: [], source_state: 'unavailable', next_due_at: null, next_retry_at: null, attempt_duration_ms: null, quota: null, storage: null },
    backlog: [], counter_replay_quarantines: [], oldest_pending_scope: 'group-all-revisions', oldest_pending_seconds: null, test_write: null };
  return view;
}
function pageProps(): RuntimeDashboardPageProps {
  const device = { device_id: 'D', name: 'Meter', protocol: 'modbus_tcp', availability_status: 'available' as const, running: true };
  return { routeState: 'live', selectedDeviceId: 'D', selectedDevice: device, devices: [device], sourceEvidenceScope: evidenceScope,
    snapshot: snapshot(), snapshotError: null, liveValues: {}, streamState: 'connected',
    setupContext: { source_rules: [], mappings: [], database_targets: [], database_config: { id: 'C', name: 'Saved connector', kind: 'sqlite', table: 'orphan_saved_draft', write_mode: 'insert', write_interval_seconds: 60 } },
    onSelectDevice: vi.fn(), onRetrySnapshot: vi.fn(async () => undefined), onReconnectStream: vi.fn() };
}
function mountPage() {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  return render(<QueryClientProvider client={client}><RuntimeDashboardPage {...pageProps()} /></QueryClientProvider>);
}
function record(value: unknown): Record<string, unknown> {
  if (!value || typeof value !== 'object' || Array.isArray(value)) throw new Error('Expected scoped parameters');
  return Object.fromEntries(Object.entries(value));
}
beforeEach(async () => {
  await i18n.changeLanguage('en'); vi.resetAllMocks();
  vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url, config) => {
    if (url === '/studio-v2/workspace/recording-plans') return { data: { success: true, data: [] } };
    expect(url).toBe('/studio-v2/workspace/runtime-context'); const params = record(config?.params); expect(params).toMatchObject(evidenceScope);
    return { data: { success: true, data: params.source_list
      ? { workspace_id: 'W', source_page: { contract_version: 1, ...evidenceScope, kind: params.source_list, query: '', limit: 50, list_revision: 'L1', rules: [], mappings: [], next_cursor: null } }
      : evidenceContext(groupEvidence()) } };
  });
});
afterEach(() => { cleanup(); clients.splice(0).forEach((client) => client.clear()); });

describe('runtime-wide counters never replace scoped SQL delivery evidence', () => {
  it.each([metrics, { write_success_total: 73 }])('LegacyCounterSuccessNeverSQLReceipt: no computed delivery percentage or green progress for %j', (observed) => {
    const { container } = render(<RuntimeSummaryPanel snapshot={snapshot(observed)} />);
    expect(container).not.toHaveTextContent(/\d+(?:\.\d+)?%|交付成功率|資料庫交付與採集統計/);
    expect(container.querySelector('[style*="width"]')).not.toBeInTheDocument();
    expect(container).toHaveTextContent('73');
    expect(container).toHaveTextContent('These global counters are not SQL receipts, readback or cleanup proof.');
  });
  it.each([undefined, null, {}, { write_success_total: -1 }, { collected_total: '0' }, { mapping_error_total: Number.NaN }, { point_state_error_total: Number.POSITIVE_INFINITY }, { write_error_total: true }])('MissingCounterUnavailableNotZero: runtimeAPI rejects missing/invalid metric leaves %j', async (observed) => {
    vi.mocked(api.get).mockResolvedValueOnce({ data: { success: true, data: { ...snapshot(), metrics: observed } } });
    const value = await runtimeAPI.getStatus('D');
    for (const key of keys) expect(value.metrics?.[key]).toBeUndefined();
    render(<RuntimeSummaryPanel snapshot={value} />);
    expect(screen.getAllByText('Unavailable')).toHaveLength(5);
    expect(api.get).toHaveBeenCalledExactlyOnceWith('/runtime/status', { params: { device_id: 'D' } });
  });
  it('PartialMetricsPreserved: runtimeAPI retains valid zero and nonzero leaves beside unknown errors', async () => {
    vi.mocked(api.get).mockResolvedValueOnce({ data: { success: true, data: { ...snapshot(), metrics: { collected_total: 0, write_success_total: 73, write_error_total: -1, mapping_error_total: 3 } } } });
    const value = await runtimeAPI.getStatus('D');
    expect(value.metrics).toEqual({ collected_total: 0, write_success_total: 73, write_error_total: undefined, mapping_error_total: 3, point_state_error_total: undefined });
    const { container } = render(<RuntimeSummaryPanel snapshot={value} />);
    expect(within(container).getAllByText('Unavailable')).toHaveLength(2);
    for (const n of [0, 73, 3]) expect(within(container).getByText(String(n))).toBeInTheDocument();
  });
  it('SavedDraftTableNotAttachedToRuntimeCounter: saved setup stays separate from observed counters', async () => {
    mountPage(); await screen.findByRole('article', { name: 'Recording group G' });
    expect(screen.getByTestId('runtime-dashboard-setup-context')).toHaveTextContent('orphan_saved_draft');
    expect(screen.getByTestId('runtime-dashboard-summary-panel')).not.toHaveTextContent('orphan_saved_draft');
    expect(screen.getByTestId('runtime-dashboard-summary-panel').querySelector('[title="orphan_saved_draft"]')).not.toBeInTheDocument();
  });
  it('MissingCounterUnavailableNotZero: absent legacy metrics do not become a synthetic all-clear', () => {
    const { container } = render(<RuntimeSummaryPanel snapshot={{ ...snapshot(), metrics: undefined }} />);
    expect(container).not.toHaveTextContent(/0 異常|0 errors/);
    expect(within(container).getAllByText('Unavailable')).toHaveLength(5);
  });
  it.each(keys)('MissingCounterUnavailableNotZero: real read projection preserves absent %s', (key) => {
    const observed: Record<string, unknown> = { ...metrics }; delete observed[key];
    const { container } = render(<RuntimeSummaryPanel snapshot={snapshot(observed)} />);
    expect(within(container).getAllByText('Unavailable')).toHaveLength(1);
    expect(container).not.toHaveTextContent(/0 異常|10 異常/);
  });
  it('ExplicitObservedZeroPreserved: all five observed zeros remain zero without derived health or success', () => {
    const observed = Object.fromEntries(keys.map((key) => [key, 0]));
    const { container } = render(<RuntimeSummaryPanel snapshot={snapshot(observed)} />);
    expect(within(container).getAllByText('0')).toHaveLength(5);
    expect(container).not.toHaveTextContent(/Unavailable|0 異常|100\.0%/);
  });
  it.each(['en', 'zh-TW'])('GlobalCounterNotDeviceOrGroupTotal: localized runtime-wide scope in %s', async (locale) => {
    await i18n.changeLanguage(locale);
    const { container } = render(<RuntimeSummaryPanel snapshot={snapshot()} />);
    expect(container).toHaveTextContent(locale === 'en' ? 'Global runtime observations' : 'Runtime 全域觀測');
    expect(container).toHaveTextContent(locale === 'en'
      ? 'Counters cover the whole runtime, not the selected device or recording group.'
      : '計數涵蓋整個 runtime，不是所選設備或記錄群組的總數。');
    expect(container).toHaveTextContent(locale === 'en'
      ? 'These global counters are not SQL receipts, readback or cleanup proof.'
      : '這些全域計數不是 SQL receipt、讀回或清理完成的證據。');
    expect(container).toHaveTextContent('1m 5s');
    for (const value of Object.values(metrics)) expect(within(container).getByText(String(value))).toBeInTheDocument();
  });
  it('ExistingScopedGroupDeliveryStillVisible: global success does not relabel A or allocate commits; read budget unchanged', async () => {
    mountPage(); const group = await screen.findByRole('article', { name: 'Recording group G' });
    const delivery = within(group).getByRole('region', { name: 'Current revision delivery' });
    expect(group).toHaveTextContent('meter_old'); expect(group).toHaveTextContent('meter_draft'); expect(group).not.toHaveTextContent('orphan_saved_draft');
    expect(delivery).toHaveTextContent('2 accepted samples'); expect(delivery).toHaveTextContent('3 output rows');
    const committedLabel = within(delivery).getByText('SQL committed'); const fact = committedLabel.parentElement;
    if (!fact) throw new Error('Expected delivery fact container');
    expect(fact).toHaveTextContent('0 output rows'); expect(delivery).not.toHaveTextContent('73');
    await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(4));
    fireEvent.click(screen.getByRole('button', { name: 'Show historical reports' }));
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(4);
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.put).not.toHaveBeenCalled(); expect(studioV2DatalinkApi.delete).not.toHaveBeenCalled();
  });
});
