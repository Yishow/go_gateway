import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchV2Page from '../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2WorkspaceAPI } from '../../../src/services/studioV2Workspace';
import { studioV2WorkspaceDevicesAPI } from '../../../src/services/studioV2WorkspaceDevices';
import { studioV2RulesAPI } from '../../../src/services/studioV2Rules';
import { studioV2WorkspaceDatabaseAPI } from '../../../src/services/studioV2WorkspaceDatabase';

vi.mock('../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell', () => ({
  WorkbenchV2Shell: ({ state }: { state: { devices: unknown[]; rules: unknown[] } }) => (
    <div
      data-testid="boot-shell"
      data-device-count={state.devices.length}
      data-rule-count={state.rules.length}
    >
      boot-shell
    </div>
  ),
}));

vi.mock('../../../src/services/studioV2Workspace', () => ({
  studioV2WorkspaceAPI: {
    get: vi.fn(),
  },
}));

vi.mock('../../../src/services/studioV2WorkspaceDevices', () => ({
  studioV2WorkspaceDevicesAPI: {
    list: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    remove: vi.fn(),
    updateOrder: vi.fn(),
  },
}));

vi.mock('../../../src/services/studioV2Rules', () => ({
  studioV2RulesAPI: {
    list: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    remove: vi.fn(),
  },
}));

vi.mock('../../../src/services/studioV2WorkspaceDatabase', () => ({
  studioV2WorkspaceDatabaseAPI: {
    getConfig: vi.fn(),
    updateConfig: vi.fn(),
    listTargets: vi.fn(),
    upsertTarget: vi.fn(),
  },
}));

function renderPage() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

  render(
    <QueryClientProvider client={queryClient}>
      <DatalinkWorkbenchV2Page />
    </QueryClientProvider>,
  );
}

describe('DatalinkWorkbenchV2Page workspace bootstrap', () => {
  it('boots from the workspace API on first open', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce({
      id: 'workspace-1',
      kind: 'single',
      status: 'empty',
      ordered_device_ids: [],
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(null);
    vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValue([]);

    renderPage();

    expect(screen.getByTestId('workbench-v2-bootstrap-loading')).toBeInTheDocument();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-root')).toBeInTheDocument();
    });

    expect(studioV2WorkspaceAPI.get).toHaveBeenCalledTimes(1);
    expect(screen.getByTestId('workbench-v2-root')).toHaveAttribute('data-workspace-id', 'workspace-1');
    expect(screen.getByTestId('workbench-v2-root')).toHaveAttribute('data-workspace-status', 'empty');
    expect(screen.getByTestId('boot-shell')).toBeInTheDocument();
  });

  it('clears seeded local draft data when the workspace has no devices or rules', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce({
      id: 'workspace-empty',
      kind: 'single',
      status: 'empty',
      ordered_device_ids: [],
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(null);
    vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValue([]);

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('boot-shell')).toHaveAttribute('data-device-count', '0');
    });

    expect(screen.getByTestId('boot-shell')).toHaveAttribute('data-rule-count', '0');
    expect(screen.getByTestId('workbench-v2-root')).toHaveAttribute('data-workspace-id', 'workspace-empty');
  });

  it('shows an actionable bootstrap error when the workspace API fails', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockRejectedValueOnce(new Error('db offline'));
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(null);
    vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValue([]);

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('workbench-v2-bootstrap-error')).toBeInTheDocument();
    });

    expect(screen.getByText('無法載入 Studio V2 工作區。請檢查後端資料庫連線後重試。')).toBeInTheDocument();
  });
});
