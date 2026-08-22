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
});
