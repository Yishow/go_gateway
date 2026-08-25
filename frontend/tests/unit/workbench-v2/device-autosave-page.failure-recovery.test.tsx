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

  it('isolates save failure per device without blocking another valid save', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(
      workspaceFixture({ ordered_device_ids: ['dev-A', 'dev-B'] }),
    );
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([
      deviceFixture(),
      deviceFixture({ id: 'dev-B', name: 'Line B PLC', connection_config: '{"host":"192.168.10.11","port":502,"slave_id":2,"timeout":5}' }),
    ]);
    vi.mocked(studioV2WorkspaceDevicesAPI.update).mockImplementation(async (deviceId, payload) => {
      if (deviceId === 'dev-A') {
        throw new Error('save failed');
      }

      return deviceFixture({
        id: deviceId,
        name: String(payload.name),
        connection_config: '{"host":"192.168.10.11","port":502,"slave_id":2,"timeout":5}',
        runtime_apply_status: 'not_running',
      });
    });

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('device-order')).toHaveTextContent('dev-A,dev-B');
    });

    fireEvent.click(screen.getByTestId('make-dev-A-fail'));
    fireEvent.click(screen.getByTestId('make-dev-B-valid'));

    await waitFor(() => {
      expect(studioV2WorkspaceDevicesAPI.update).toHaveBeenCalledWith(
        'dev-A',
        expect.objectContaining({ name: 'Line A Broken' }),
      );
      expect(studioV2WorkspaceDevicesAPI.update).toHaveBeenCalledWith(
        'dev-B',
        expect.objectContaining({ name: 'Line B Saved' }),
      );
    });

    await waitFor(() => {
      expect(screen.getByTestId('device-save-state-dev-A')).toHaveTextContent('save-error');
      expect(screen.getByTestId('device-save-state-dev-B')).toHaveTextContent('saved');
      expect(screen.getByTestId('device-save-error-dev-A')).toHaveTextContent('errors.autosave_failed');
    });
  });

  it('keeps apply_failed visible instead of flattening it into saved', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(workspaceFixture());
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([
      deviceFixture({ status: 'active' }),
    ]);
    vi.mocked(studioV2WorkspaceDevicesAPI.update).mockResolvedValueOnce(
      deviceFixture({
        name: 'Line A Broken',
        status: 'active',
        connection_config: '{"host":"192.168.10.20","port":502,"slave_id":1,"timeout":5}',
        runtime_apply_status: 'apply_failed',
        runtime_apply_message: 'runtime sync failed',
      }),
    );

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('device-order')).toHaveTextContent('dev-A');
    });

    fireEvent.click(screen.getByTestId('make-dev-A-fail'));

    await waitFor(() => {
      expect(screen.getByTestId('device-save-state-dev-A')).toHaveTextContent('save-error');
      expect(screen.getByTestId('device-save-error-dev-A')).toHaveTextContent('runtime sync failed');
      expect(screen.getByTestId('device-runtime-apply-dev-A')).toHaveTextContent('apply_failed');
    });
  });

  it('does not mark apply_failed as unrecovered draft loss', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(workspaceFixture());
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([
      deviceFixture({ status: 'active' }),
    ]);
    vi.mocked(studioV2WorkspaceDevicesAPI.update).mockResolvedValueOnce(
      deviceFixture({
        name: 'Line A Broken',
        status: 'active',
        connection_config: '{"host":"192.168.10.20","port":502,"slave_id":1,"timeout":5}',
        runtime_apply_status: 'apply_failed',
        runtime_apply_message: 'runtime sync failed',
      }),
    );

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('device-order')).toHaveTextContent('dev-A');
    });

    fireEvent.click(screen.getByTestId('make-dev-A-fail'));

    await waitFor(() => {
      expect(screen.getByTestId('device-runtime-apply-dev-A')).toHaveTextContent('apply_failed');
    });

    expect(window.sessionStorage.getItem('wbv2_unrecovered_draft')).toBeNull();
  });
});
