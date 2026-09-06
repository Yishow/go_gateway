import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, screen, waitFor } from '@testing-library/react';
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

  it('saves one valid local device and marks it saved', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(
      workspaceFixture({ status: 'empty', ordered_device_ids: [] }),
    );
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockResolvedValueOnce(
      deviceFixture({
        id: 'dev-01',
        name: 'Line A Saved',
        description: 'Modbus TCP PLC (Line A 主控)',
        connection_config: '{"host":"192.168.1.100","port":502,"slave_id":1,"timeout":5}',
        runtime_apply_status: 'not_running',
      }),
    );

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId('seed-dev-01'));
    fireEvent.click(screen.getByTestId('make-dev-01-valid'));

    await waitFor(() => {
      expect(studioV2WorkspaceDevicesAPI.create).toHaveBeenCalledWith(
        expect.objectContaining({
          id: 'dev-01',
          name: 'Line A Saved',
        }),
      );
    });

    await waitFor(() => {
      expect(screen.getByTestId('device-save-state-dev-01')).toHaveTextContent('saved');
      expect(screen.getByTestId('device-runtime-apply-dev-01')).toHaveTextContent('not_running');
    });
  });

  it('keeps invalid local values without overwriting the backend', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce(
      workspaceFixture({ status: 'empty', ordered_device_ids: [] }),
    );
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([]);

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId('seed-dev-01'));
    fireEvent.click(screen.getByTestId('make-dev-01-invalid'));

    await waitFor(() => {
      expect(screen.getByTestId('device-name-dev-01')).toHaveTextContent('');
      expect(screen.getByTestId('device-save-state-dev-01')).toHaveTextContent('draft-invalid');
    });

    expect(studioV2WorkspaceDevicesAPI.create).not.toHaveBeenCalled();
    expect(studioV2WorkspaceDevicesAPI.update).not.toHaveBeenCalled();
  });

  it('persists an MC protocol switch with backend-aligned connection fields', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(
      workspaceFixture({ ordered_device_ids: ['dev-01'] }),
    );
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([
      deviceFixture({ id: 'dev-01' }),
    ]);
    vi.mocked(studioV2WorkspaceDevicesAPI.update).mockResolvedValueOnce(
      deviceFixture({
        id: 'dev-01',
        protocol: 'mc_3e',
        connection_config: '{"host":"192.168.1.100","port":6000,"station_no":0,"network_no":0,"pc_no":255,"io_no":1023,"timeout":5}',
      }),
    );

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('device-save-state-dev-01')).toHaveTextContent('saved');
    });

    fireEvent.click(screen.getByTestId('change-dev-01-to-mc'));

    await waitFor(() => {
      expect(studioV2WorkspaceDevicesAPI.update).toHaveBeenCalledWith('dev-01', {
        name: 'Line A PLC',
        description: '',
        protocol: 'mc_3e',
        connection_config: {
          host: '192.168.1.100',
          port: 6000,
          station_no: 0,
          network_no: 0,
          pc_no: 255,
          io_no: 1023,
          timeout: 5,
        },
      });
    });
  });

  it('marks session draft recovery state when a local device draft is not yet persisted', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce(
      workspaceFixture({ status: 'empty', ordered_device_ids: [] }),
    );
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([]);

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTestId('seed-dev-01'));
    fireEvent.click(screen.getByTestId('make-dev-01-invalid'));

    await waitFor(() => {
      expect(window.sessionStorage.getItem('wbv2_unrecovered_draft')).toBe('1');
    });
  });

  it('keeps an invalid running device visible while marking it unavailable', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(
      workspaceFixture({ ordered_device_ids: ['dev-01'] }),
    );
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([
      deviceFixture({ id: 'dev-01', status: 'active', running: true }),
    ]);
    vi.mocked(studioV2WorkspaceDevicesAPI.updateAvailability).mockResolvedValueOnce(
      deviceFixture({
        id: 'dev-01',
        availability_status: 'unavailable',
        availability_reason: 'device form is invalid',
        running: false,
      }),
    );

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('device-order')).toHaveTextContent('dev-01');
      expect(screen.getByTestId('device-save-state-dev-01')).toHaveTextContent('saved');
      expect(screen.getByTestId('device-running-dev-01')).toHaveTextContent('true');
    });

    fireEvent.click(screen.getByTestId('make-dev-01-invalid'));

    await waitFor(() => {
      expect(studioV2WorkspaceDevicesAPI.updateAvailability).toHaveBeenCalledWith('dev-01', {
        availability_status: 'unavailable',
        availability_reason: 'device form is invalid',
      });
      expect(screen.getByTestId('device-order')).toHaveTextContent('dev-01');
      expect(screen.getByTestId('device-availability-dev-01')).toHaveTextContent('unavailable');
      expect(screen.getByTestId('device-running-dev-01')).toHaveTextContent('false');
    });
  });
});
