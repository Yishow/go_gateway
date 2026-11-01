import { useState } from 'react';
import { RecordingDraftNavigation, useRecordingDraftGuard, useRecordingReadUnavailable } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import { studioV2WorkspaceKeys } from '../../../src/hooks/datalink/keys';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
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
  WorkbenchV2Shell: function DraftShell() {
    const [text, setText] = useState(''); const unavailable = useRecordingReadUnavailable(); useRecordingDraftGuard({ version: text, dirty: Boolean(text) });
    return <input aria-label="record draft" data-unavailable={unavailable} value={text} onChange={event => setText(event.target.value)} />;
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

  it.each(['workspace', 'mappings'] as const)('retains unsaved editor fields across a failed background %s read', async kind => {
    const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(<QueryClientProvider client={queryClient}><RecordingDraftNavigation><DatalinkWorkbenchV2Page /></RecordingDraftNavigation></QueryClientProvider>);
    fireEvent.change(await screen.findByLabelText('record draft'), { target: { value: 'Unsaved page fields' } });
    const key = kind === 'workspace' ? studioV2WorkspaceKeys.bootstrap() : studioV2WorkspaceKeys.mappings();
    const read = kind === 'workspace' ? mocks.workspace.get : mocks.mappings.list;
    const saved = queryClient.getQueryData(key);
    read.mockRejectedValue(new Error('offline'));
    await act(async () => { await queryClient.refetchQueries({ queryKey: key, exact: true }); });
    await waitFor(() => expect(queryClient.getQueryState(key)?.status).toBe('error'));
    expect(screen.getByLabelText('record draft')).toHaveValue('Unsaved page fields');
    expect(screen.getByLabelText('record draft')).toHaveAttribute('data-unavailable', 'true');
    read.mockResolvedValue(saved);
    fireEvent.click(screen.getByTestId('recording-draft-refresh-failed').querySelector('button')!);
    await waitFor(() => expect(screen.queryByTestId('recording-draft-refresh-failed')).not.toBeInTheDocument());
    expect(screen.getByLabelText('record draft')).toHaveValue('Unsaved page fields');
    expect(screen.getByLabelText('record draft')).toHaveAttribute('data-unavailable', 'false');
  });
});
