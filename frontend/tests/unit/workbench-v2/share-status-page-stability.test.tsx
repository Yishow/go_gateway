import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchV2Page from '../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import type { ModbusShareStatus } from '../../../src/types/modbusShare';

const mocks = vi.hoisted(() => ({
  workspace: { get: vi.fn() },
  devices: { list: vi.fn(), create: vi.fn(), update: vi.fn(), updateAvailability: vi.fn(), remove: vi.fn(), updateOrder: vi.fn() },
  rules: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() },
  mappings: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() },
  database: { getConfig: vi.fn(), updateConfig: vi.fn(), listTargets: vi.fn(), upsertTarget: vi.fn() },
  modbusShare: { status: vi.fn() },
  shareStatuses: [] as Array<ModbusShareStatus | null | undefined>,
}));

vi.mock('../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell', () => ({
  WorkbenchV2Shell: ({ shareStatus, actions }: { shareStatus?: ModbusShareStatus | null; actions: { dispatch: (action: unknown) => void } }) => {
    mocks.shareStatuses.push(shareStatus);
    return (
      <button
        type="button"
        data-testid="unrelated-state-change"
        onClick={() => actions.dispatch({ type: 'SET_STATE', payload: { showSummaryRail: false } })}
      />
    );
  },
}));
vi.mock('../../../src/services/studioV2Workspace', () => ({ studioV2WorkspaceAPI: mocks.workspace }));
vi.mock('../../../src/services/studioV2WorkspaceDevices', () => ({ studioV2WorkspaceDevicesAPI: mocks.devices }));
vi.mock('../../../src/services/studioV2Rules', () => ({ studioV2RulesAPI: mocks.rules }));
vi.mock('../../../src/services/studioV2Mappings', () => ({ studioV2MappingsAPI: mocks.mappings }));
vi.mock('../../../src/services/studioV2WorkspaceDatabase', () => ({ studioV2WorkspaceDatabaseAPI: mocks.database }));
vi.mock('../../../src/services/modbusShare', () => ({ modbusShareAPI: mocks.modbusShare }));

const bootstrapStatus: ModbusShareStatus = {
  enabled: true,
  running: true,
  port: 5020,
  address: '127.0.0.1',
  bind_state: 'pass',
  mapping_count: 0,
};

describe('DatalinkWorkbenchV2Page Modbus Share status', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mocks.shareStatuses.splice(0);
    mocks.workspace.get.mockResolvedValue({
      id: 'workspace-1',
      kind: 'single',
      status: 'ready',
      ordered_device_ids: [],
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
      modbus_share: { hydration_state: 'ready', readiness: true, settings_revision: 'settings-1', status: bootstrapStatus },
    });
    mocks.modbusShare.status.mockResolvedValue({ ...bootstrapStatus, settings_revision: 'settings-1' });
    mocks.devices.list.mockResolvedValue([]);
    mocks.rules.list.mockResolvedValue([]);
    mocks.mappings.list.mockResolvedValue([]);
    mocks.database.getConfig.mockResolvedValue(null);
    mocks.database.listTargets.mockResolvedValue([]);
  });

  it('keeps the same share status object across unrelated re-renders', async () => {
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(
      <QueryClientProvider client={queryClient}>
        <DatalinkWorkbenchV2Page />
      </QueryClientProvider>,
    );
    await waitFor(() => expect(mocks.modbusShare.status).toHaveBeenCalled());
    const button = await screen.findByTestId('unrelated-state-change');
    await waitFor(() => expect(mocks.shareStatuses.at(-1)?.settings_revision).toBe('settings-1'));
    const before = mocks.shareStatuses.length;
    const statusBefore = mocks.shareStatuses.at(-1);

    fireEvent.click(button);

    await waitFor(() => expect(mocks.shareStatuses.length).toBeGreaterThan(before));
    expect(mocks.shareStatuses.at(-1)).toBe(statusBefore);
  });
});
