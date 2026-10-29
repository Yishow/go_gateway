import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { cleanup, render, screen, waitFor, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '../../../src/i18n/config';
import { RuntimeDashboardPage, type RuntimeDashboardPageProps } from '../../../src/features/datalink/runtime-dashboard/RuntimeDashboardPage';
import { studioV2DatalinkApi } from '../../../src/services/studioV2Workspace';
import { evidenceContext, evidenceScope, ruleEvidenceFixture } from '../services/runtimeRuleEvidenceFixture';
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
const clients: QueryClient[] = [];
const baseProps: RuntimeDashboardPageProps = {
  routeState: 'loading', selectedDeviceId: 'D', selectedDevice: { device_id: 'D', name: 'Meter D', protocol: 'modbus_tcp', running: false, availability_status: 'available' },
  devices: [], sourceEvidenceScope: evidenceScope, setupContext: null, snapshot: null, snapshotError: null, liveValues: {}, streamState: 'connecting',
  onSelectDevice: vi.fn(), onRetrySnapshot: async () => undefined, onReconnectStream: vi.fn(),
};
function mount(props: Partial<RuntimeDashboardPageProps> = {}) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } }); clients.push(client);
  const content = (next: Partial<RuntimeDashboardPageProps>) => <QueryClientProvider client={client}><RuntimeDashboardPage {...baseProps} {...next} /></QueryClientProvider>;
  const result = render(content(props));
  return { ...result, update(next: Partial<RuntimeDashboardPageProps>) { result.rerender(content(next)); } };
}
beforeEach(async () => {
  await i18n.changeLanguage('en'); vi.mocked(studioV2DatalinkApi.get).mockReset();
  const view = ruleEvidenceFixture(); view.groups[0].applied_effective_at = null;
  vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url) => ({ data: { success: true, data: url === '/studio-v2/workspace/runtime-context' ? evidenceContext(view) : [] } }));
});
afterEach(() => { cleanup(); clients.splice(0).forEach((client) => client.clear()); });

describe('recording configuration on the existing runtime dashboard', () => {
  it('reads revision-bound recording evidence while collector snapshot is still loading', async () => {
    mount();
    await waitFor(() => expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/runtime-context', { params: evidenceScope, signal: expect.any(AbortSignal) }));
    expect(studioV2DatalinkApi.get).toHaveBeenCalledTimes(1);
  });
  it('shows effective A destination separately from saved C before the requested B boundary is known', async () => {
    mount();
    expect(await screen.findByText('meter_old')).toBeInTheDocument();
    expect(screen.getByText('meter_draft')).toBeInTheDocument();
    expect(screen.queryByText('meter_new')).not.toBeInTheDocument();
  });
  it.each(['loading', 'error', 'empty', 'live', 'degraded'] as const)('keeps recording evidence accessible in collector route state %s', async (routeState) => {
    mount({ routeState });
    const panel = screen.getByRole('region', { name: 'Recording rules' });
    await within(panel).findByText('meter_old');
    expect(panel).toBeVisible(); expect(panel).not.toHaveTextContent('SQL committed');
    const recordingCalls = vi.mocked(studioV2DatalinkApi.get).mock.calls.filter(([url]) => url === '/studio-v2/workspace/runtime-context');
    expect(recordingCalls).toHaveLength(1);
  });
  it('preserves the same evidence observation without extra group reads across loading/live/error transitions', async () => {
    const mounted = mount(); await screen.findByText('meter_old');
    mounted.update({ routeState: 'live' });
    await waitFor(() => expect(vi.mocked(studioV2DatalinkApi.get).mock.calls.some(([url]) => url === '/studio-v2/workspace/recording-plans')).toBe(true));
    expect(screen.getByRole('region', { name: 'Recording rules' })).toHaveTextContent('meter_old');
    mounted.update({ routeState: 'error' });
    expect(screen.getByRole('region', { name: 'Recording rules' })).toHaveTextContent('meter_old');
    expect(vi.mocked(studioV2DatalinkApi.get).mock.calls.filter(([url]) => url === '/studio-v2/workspace/runtime-context')).toHaveLength(1);
  });
  it.each(['empty-workspace', 'missing-device-context'] as const)('shows unavailable rather than a fallback device for %s', (routeState) => {
    mount({ routeState, selectedDeviceId: null, selectedDevice: null, sourceEvidenceScope: null });
    expect(screen.getByRole('region', { name: 'Recording rules' })).toHaveTextContent('Recording evidence unavailable until');
    expect(studioV2DatalinkApi.get).not.toHaveBeenCalled();
  });
});
