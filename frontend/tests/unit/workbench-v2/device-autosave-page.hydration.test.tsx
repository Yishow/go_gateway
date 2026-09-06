import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { screen, waitFor } from '@testing-library/react';
import { render } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import {
  deviceFixture,
  resetDeviceAutosaveMocks,
  studioV2WorkspaceAPI,
  studioV2WorkspaceDevicesAPI,
  workspaceFixture,
} from './device-autosave-page.testHarness';
import DatalinkWorkbenchV2Page from '../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import {
  hydrateStudioV2Device,
  isStudioV2DeviceValid,
} from '../../../src/features/datalink/workbench-v2/state/studioV2DeviceAutosave';

function renderPage() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });
  render(
    <QueryClientProvider client={queryClient}>
      <DatalinkWorkbenchV2Page />
    </QueryClientProvider>,
  );
}

describe('DatalinkWorkbenchV2Page device autosave orchestration', () => {
  beforeEach(() => {
    resetDeviceAutosaveMocks();
  });

  it('hydrates persisted devices in backend order on reload', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce(
      workspaceFixture({ ordered_device_ids: ['dev-B', 'dev-A'] }),
    );
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([
      deviceFixture({ id: 'dev-B', name: 'Line B PLC', connection_config: '{"host":"192.168.10.11","port":502,"slave_id":2,"timeout":5}' }),
      deviceFixture(),
    ]);

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('device-order')).toHaveTextContent('dev-B,dev-A');
    });
  });

  it('hydrates persisted running truth on reload', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce(workspaceFixture());
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([
      deviceFixture({ status: 'active', availability_status: 'available', running: true }),
    ]);

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('device-running-dev-A')).toHaveTextContent('true');
    });
  });

  it('hydrates a successful persisted probe as tested on reload', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce(workspaceFixture());
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([
      deviceFixture({
        last_test_at: '2026-08-26T01:29:36+08:00',
        last_test_success: true,
        last_test_error: '',
      }),
    ]);

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('device-test-status-dev-A')).toHaveTextContent('success');
    });
  });
  it('backfills a missing station number from the protocol defaults on hydration', () => {
    const device = hydrateStudioV2Device(deviceFixture({
      id: 'dev-mc',
      protocol: 'mc_3e',
      connection_config: '{"host":"192.168.1.100","port":6000,"timeout":5}',
    }));

    expect(device.config.station_no).toBe(0);
    expect(isStudioV2DeviceValid(device)).toBe(true);
  });

  it('backfills a missing station number for FATEK devices on hydration', () => {
    const device = hydrateStudioV2Device(deviceFixture({
      id: 'dev-fatek',
      protocol: 'fatek_fbs',
      connection_config: '{"mode":"tcp","host":"192.168.1.100","port":500,"timeout":5}',
    }));

    expect(device.config.station_no).toBe(1);
    expect(isStudioV2DeviceValid(device)).toBe(true);
  });

  it('keeps a persisted station number instead of overwriting it with the default', () => {
    const device = hydrateStudioV2Device(deviceFixture({
      id: 'dev-mc',
      protocol: 'mc_3e',
      connection_config: '{"host":"192.168.1.100","port":6000,"station_no":7,"timeout":5}',
    }));

    expect(device.config.station_no).toBe(7);
  });

  it('does not autosave devices that load without a station number', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce(
      workspaceFixture({ ordered_device_ids: ['dev-mc'] }),
    );
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([
      deviceFixture({
        id: 'dev-mc',
        protocol: 'mc_3e',
        connection_config: '{"host":"192.168.1.100","port":6000,"timeout":5}',
      }),
    ]);

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('device-save-state-dev-mc')).toHaveTextContent('saved');
    });
    expect(studioV2WorkspaceDevicesAPI.update).not.toHaveBeenCalled();
    expect(studioV2WorkspaceDevicesAPI.create).not.toHaveBeenCalled();
  });
});
