import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { runtimeAPI } from '@/services/datalink';
import { api } from '@/services/datalinkClient';
import { studioV2RuntimeContextAPI } from '@/services/studioV2RuntimeContext';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { CollectorHealthPanel } from '@/features/datalink/runtime-dashboard/CollectorHealthPanel';
import { RuntimeDiagnosticsPanel } from '@/features/datalink/runtime-dashboard/RuntimeDiagnosticsPanel';
import { RuntimeSetupContextPanel } from '@/features/datalink/runtime-dashboard/RuntimeSetupContextPanel';
import type { RuntimeStatusWithDiagnostics } from '@/types/runtimeDiagnostics';
import { createRuntimeSetupFixture } from './runtimeSetupFixture';
vi.mock('@/services/datalinkClient', async (original) => ({ ...(await original<typeof import('@/services/datalinkClient')>()), api: { get: vi.fn() } }));
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, fallback?: unknown) => typeof fallback === 'string' ? fallback : key }) }));
const status = async (diagnostic: unknown) => {
  vi.mocked(api.get).mockResolvedValueOnce({ data: { success: true, data: { collectors: [], diagnostics: [diagnostic] } } });
  return await runtimeAPI.getStatus() as RuntimeStatusWithDiagnostics;
};
describe('real runtime service projection to operator UI seams', () => {
  it('does not turn successful device-scope delivery into a synthetic failure', async () => {
    const snapshot = await status({ scope: 'device:device-A', latest_successful_stage: 'database_delivery', stages: [{ stage: 'database_delivery', status: 'success' }] });
    expect(snapshot.diagnostics?.[0]).not.toHaveProperty('failure_stage');
    render(<RuntimeDiagnosticsPanel diagnostics={snapshot.diagnostics} selectedDeviceId="device-A" />);
    expect(screen.getByTestId('runtime-dashboard-diagnostics-panel')).toHaveTextContent('diagnostics.noErrorDetected');
    expect(screen.getByTestId('runtime-diagnostics-action')).toHaveClass('sr-only');
  });
  it('shows unavailable instead of warning for an unsupported collector status and does not claim running', async () => {
    vi.mocked(api.get).mockResolvedValueOnce({ data: { success: true, data: { collectors: [{ device_id: 'device-A', protocol: 'modbus_tcp', status: 'future-status', availability_status: 'available', running: true }] } } });
    const snapshot = await runtimeAPI.getStatus();
    expect(snapshot.collectors[0]).toMatchObject({ status: 'unavailable', running: false });
    render(<CollectorHealthPanel collector={snapshot.collectors[0]} />);
    expect(screen.getByTestId('runtime-dashboard-health-panel')).toHaveTextContent('unavailable');
  });
  it('retains allowlisted failure codes and the precise scope needed to choose actionable copy', async () => {
    const snapshot = await status({ scope: 'device:device-A', failure_stage: 'modbus_share_delivery', failure_code: 'modbus_share_delivery', stages: [] });
    render(<RuntimeDiagnosticsPanel diagnostics={snapshot.diagnostics} selectedDeviceId="device-A" />);
    expect(screen.getByTestId('runtime-diagnostics-action')).toHaveTextContent('Retry Modbus Share delivery');
  });
  it.each(['BFF', 'SB1A', 'W1A'])('keeps supported MC register %s and safe Unicode destination labels visible', async (address) => {
    const setup = createRuntimeSetupFixture(); setup.source_rules[0].start_address = address;
    Object.assign(setup.database_config!, { name: '泵浦 A (°C)', schema: '廠區', table: '電力資料（A）_Café', timestamp_column: '記錄時間' });
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success: true, data: { workspace_id: 'ws-1', devices: [], setup } } });
    const context = await studioV2RuntimeContextAPI.get();
    expect(context.setup?.database_config).toMatchObject({ schema: '廠區', table: '電力資料（A）_Café', timestamp_column: '記錄時間' });
    render(<RuntimeSetupContextPanel setupContext={context.setup!} selectedDevice={null} />);
    expect(screen.getByTestId('runtime-dashboard-setup-context')).toHaveTextContent(address);
    expect(screen.getByTestId('runtime-dashboard-setup-context')).toHaveTextContent('泵浦 A (°C)');
    expect(screen.getByTestId('runtime-dashboard-setup-context')).toHaveTextContent('電力資料（A）_Café');
  });
});
