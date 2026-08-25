import { act, fireEvent, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it } from 'vitest';
import { mockEventSources, mockRuntimeContext, mockRuntimeStatus, resetRuntimeDashboardFixture, renderRoute } from './runtime-dashboard-state.fixture';

describe('runtime dashboard route state', () => {
  beforeEach(resetRuntimeDashboardFixture);

  it('boots from workspace context without requiring device_id and selects the first available v2-ordered device', async () => {
    renderRoute('/studio/runtime');
    await waitFor(() => expect(mockRuntimeStatus).toHaveBeenCalledWith('device-A'));
    expect(screen.getByTestId('runtime-dashboard-selected-device')).toHaveTextContent('Mixer PLC');
    expect(screen.getByTestId('runtime-dashboard-device-switcher')).toHaveTextContent('Review the device setup in Studio V2 and retry.');
    expect(screen.getByTestId('runtime-dashboard-device-switcher')).not.toHaveTextContent('invalid Step 1 configuration');
    expect(screen.getByTestId('runtime-dashboard-location')).toHaveTextContent('');
  });

  it('passes configured setup context from runtime context to the dashboard', async () => {
    renderRoute('/studio/runtime?device_id=device-A');
    await waitFor(() => expect(mockRuntimeStatus).toHaveBeenCalledWith('device-A'));
    act(() => mockEventSources[0].open());
    const panel = await screen.findByTestId('runtime-dashboard-setup-context');
    expect(panel).toHaveTextContent('database-target-missing');
    expect(panel).toHaveTextContent('40001'); expect(panel).toHaveTextContent('line01.temp.inlet');
    expect(panel).toHaveTextContent('入口溫度'); expect(panel).toHaveTextContent('PostgreSQL Connector');
    expect(panel).toHaveTextContent('gateway_metrics');
  });

  it('updates query state when switching device', async () => {
    renderRoute('/studio/runtime?device_id=device-A'); await screen.findByTestId('runtime-dashboard-route');
    fireEvent.click(screen.getByRole('button', { name: 'Filler PLC' }));
    await waitFor(() => expect(screen.getByTestId('runtime-dashboard-location')).toHaveTextContent('?device_id=device-B'));
  });

  it('keeps device_id as an optional override when present', async () => {
    renderRoute('/studio/runtime?device_id=device-B');
    await waitFor(() => expect(mockRuntimeStatus).toHaveBeenCalledWith('device-B'));
    expect(screen.getByTestId('runtime-dashboard-selected-device')).toHaveTextContent('Filler PLC');
  });

  it('transitions from loading to live after snapshot and stream attach', async () => {
    renderRoute('/studio/runtime?device_id=device-A');
    expect(screen.getByTestId('runtime-dashboard-route-state')).toHaveTextContent('loading');
    await waitFor(() => expect(mockRuntimeStatus).toHaveBeenCalledWith('device-A'));
    expect(mockEventSources).toHaveLength(1); act(() => mockEventSources[0].open());
    await waitFor(() => expect(screen.getByTestId('runtime-dashboard-route-state')).toHaveTextContent('live'));
  });

  it('stays on runtime and shows an empty workspace state when no available devices exist', () => {
    mockRuntimeContext.devices.splice(0, mockRuntimeContext.devices.length, { device_id: 'device-B', name: 'Filler PLC', protocol: 'modbus_tcp', running: false, availability_status: 'unavailable', availability_reason: 'invalid Step 1 configuration' });
    mockRuntimeContext.default_device_id = null;
    renderRoute('/studio/runtime');
    expect(screen.getByTestId('runtime-dashboard-empty-workspace')).toBeInTheDocument();
    expect(screen.getByTestId('runtime-dashboard-device-switcher')).toHaveTextContent('Filler PLC');
    expect(mockRuntimeStatus).not.toHaveBeenCalled();
  });
});
