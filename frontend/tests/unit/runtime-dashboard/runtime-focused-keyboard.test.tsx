import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '../../../src/i18n/config';
import { RuntimeDashboardPage } from '../../../src/features/datalink/runtime-dashboard/RuntimeDashboardPage';
import type { RuntimeDashboardState } from '../../../src/features/datalink/runtime-dashboard/useRuntimeDashboardState';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import { evidenceContext, evidenceScope, ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() } }));
const clients: QueryClient[] = [];
const receivedAt = Date.parse('2026-10-29T10:00:00Z');
function state(routeState: RuntimeDashboardState['routeState'] = 'live') {
  const device = { device_id: 'D', name: 'Meter', protocol: 'modbus_tcp', running: true, availability_status: 'available' as const };
  return { routeState, selectedDeviceId: 'D', selectedDevice: device, devices: [device], setupContext: { source_rules: [], mappings: [], database_targets: [] }, sourceEvidenceScope: evidenceScope,
    snapshot: { running: true, uptime_seconds: 10, collectors: [{ device_id: 'D', device_name: 'Meter', protocol: 'modbus_tcp', status: 'running' as const, availability_status: 'available' as const, running: true, points_total: 1, points_healthy: 1, points_stale: 0, points_error: 0, last_read_at: '2026-10-29T09:59:00Z', last_error: null, breaker_state: 'closed' }] },
    snapshotError: null, liveValues: {}, streamState: 'connected' as const, snapshotReceivedAt: receivedAt, snapshotStale: routeState === 'degraded',
    onSelectDevice: vi.fn(), onRetrySnapshot: vi.fn(async () => undefined), onReconnectStream: vi.fn() } satisfies RuntimeDashboardState;
}
function evidence(detail: boolean) {
  const value = ruleEvidenceFixture();
  value.groups[0].delivery = { current: { binding: { workspace_id: 'W', workspace_revision: 'W1', entity_revision: 'A' }, intake: null, totals_scope: 'group', shared: false,
    stages: { collecting: 2, queued: 3, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 }, last_commit: null, sources: [], source_state: 'unavailable', next_due_at: null, next_retry_at: null, attempt_duration_ms: null, quota: null, storage: null },
    backlog: [], counter_replay_quarantines: [], oldest_pending_scope: 'group-all-revisions', oldest_pending_seconds: null, test_write: null };
  if (detail) value.source = { rule_id: 'R', saved_revision: 'R1', point_id: 'P', mapping_id: 'M', tag_id: 'T', effective_revision: null, state: 'saved-only', proof: 'unavailable', effective: null,
    saved: { binding: { workspace_id: 'W', workspace_revision: 'W1', entity_revision: 'R1' }, mapping_revision: 'M1', enabled: true, address: '40001', protocol_function: '03', read_type: 'int16', result_type: 'float64', polling: null, transforms: [] } };
  return value;
}
function params(value: unknown): Record<string, unknown> {
  if (!value || typeof value !== 'object' || Array.isArray(value)) throw new Error('Expected scoped request');
  return Object.fromEntries(Object.entries(value));
}
function response(data: unknown) { return { data: { success: true, data } }; }
function mount(props = state()) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  const content = (next: typeof props) => <QueryClientProvider client={client}><RuntimeDashboardPage {...next} /></QueryClientProvider>;
  const view = render(content(props)); return { ...view, props, update(next: typeof props) { view.rerender(content(next)); } };
}
beforeEach(async () => {
  await i18n.changeLanguage('en'); vi.resetAllMocks();
  vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url, config) => {
    if (url === '/studio-v2/workspace/recording-plans') return response([]);
    expect(url).toBe('/studio-v2/workspace/runtime-context'); const p = params(config?.params); expect(p).toMatchObject(evidenceScope);
    if (p.source_list) return response({ workspace_id: 'W', source_page: { contract_version: 1, ...evidenceScope, kind: p.source_list, query: String(p.q ?? ''), limit: 50, list_revision: 'L1', next_cursor: null,
      rules: p.source_list === 'rules' ? [{ id: 'R', revision_id: 'R1', device_id: 'D', start_address: '40001', count: 1, data_type: 'int16', naming_prefix: 'Meter', enabled: true }] : [],
      mappings: p.source_list === 'mappings' ? [{ id: 'M', rule_id: 'R', rule_revision: 'R1', device_id: 'D', point_id: 'P', address: '40001', tag_id: 'T', display_name: 'Power', enabled: true }] : [] } });
    return response(evidenceContext(evidence(Boolean(p.point_id))));
  });
});
afterEach(() => { cleanup(); clients.splice(0).forEach((client) => client.clear()); });

/** detail:0 is keyboard-origin activation in jsdom; real browser focus geometry remains separately gated. */
function keyboardActivate(button: HTMLElement, key = 'Enter') {
  button.focus(); expect(button).toHaveFocus(); fireEvent.keyDown(button, { key }); fireEvent.click(button, { detail: 0 }); fireEvent.keyUp(button, { key });
}
describe('localized focused-monitor keyboard controls', () => {
  it.each(['en', 'zh-TW'])('exposes reports and preserves source close focus in %s', async (locale) => {
    await i18n.changeLanguage(locale); mount();
    const reports = await screen.findByRole('button', { name: i18n.t('focused.showReports', { ns: 'runtime-dashboard' }) });
    keyboardActivate(reports, ' '); expect(reports).toHaveAttribute('aria-expanded', 'true'); expect(reports).toHaveFocus();
    expect(screen.getByTestId('history-reports-empty')).toBeVisible();
    const controlled = reports.getAttribute('aria-controls'); expect(controlled).toBeTruthy();
    expect(document.getElementById(controlled ?? '')).toContainElement(screen.getByTestId('history-reports-empty'));
    keyboardActivate(reports); expect(reports).toHaveAttribute('aria-expanded', 'false'); expect(reports).toHaveFocus();
    const open = await screen.findByRole('button', { name: locale === 'en' ? 'Inspect source 40001' : '檢視來源 40001' });
    keyboardActivate(open); await screen.findByText(i18n.t('sourceDetail.savedOnly', { ns: 'runtime-dashboard' }));
    const close = screen.getByRole('button', { name: i18n.t('sourceDetail.close', { ns: 'runtime-dashboard' }) });
    keyboardActivate(close); expect(open).toHaveFocus();
    expect(screen.queryByRole('region', { name: i18n.t('sourceDetail.title', { ns: 'runtime-dashboard' }) })).not.toBeInTheDocument();
  });
  it.each(['en', 'zh-TW'])('keeps typed reconnect and snapshot retry distinct in %s', async (locale) => {
    await i18n.changeLanguage(locale); const props = state('degraded');
    const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
    const content = (recovery: boolean) => <QueryClientProvider client={client}><RuntimeDashboardPage {...props} streamRecovery={recovery ? { code: 'runtime_stream_unavailable', retryable: true } : null} /></QueryClientProvider>;
    const view = render(content(true));
    const reconnect = screen.getByRole('button', { name: i18n.t('errors.reconnectStream', { ns: 'runtime-dashboard' }) });
    keyboardActivate(reconnect); expect(props.onReconnectStream).toHaveBeenCalledTimes(1); expect(props.onRetrySnapshot).not.toHaveBeenCalled();
    view.rerender(content(false));
    expect(screen.getByRole('button', { name: i18n.t('errors.reconnectStream', { ns: 'runtime-dashboard' }) })).toBeEnabled();
    keyboardActivate(screen.getByRole('button', { name: i18n.t('errors.retry', { ns: 'runtime-dashboard' }) }));
    expect(props.onRetrySnapshot).toHaveBeenCalledTimes(1); expect(props.onReconnectStream).toHaveBeenCalledTimes(1);
  });
});
