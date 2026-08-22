import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import DatalinkWorkbenchV2Page from '../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { studioV2WorkspaceAPI } from '../../../src/services/studioV2Workspace';
import { studioV2WorkspaceDevicesAPI } from '../../../src/services/studioV2WorkspaceDevices';
import { studioV2MappingsAPI } from '../../../src/services/studioV2Mappings';
import { studioV2RulesAPI } from '../../../src/services/studioV2Rules';
import { studioV2WorkspaceDatabaseAPI } from '../../../src/services/studioV2WorkspaceDatabase';

vi.mock('../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell', () => ({
  WorkbenchV2Shell: ({
    state,
    actions,
  }: {
    state: {
      devices: Array<{
        id: string;
        name: string;
        save_state: string;
        save_error?: string | null;
        runtime_apply_status?: string | null;
        availability_status?: string;
        running?: boolean;
      }>;
    };
    actions: { dispatch: (action: unknown) => void };
  }) => (
    <div data-testid="autosave-shell">
      <div data-testid="device-order">{state.devices.map((device) => device.id).join(',')}</div>
      {state.devices.map((device) => (
        <div key={device.id}>
          <div data-testid={`device-name-${device.id}`}>{device.name}</div>
          <div data-testid={`device-save-state-${device.id}`}>{device.save_state}</div>
          <div data-testid={`device-save-error-${device.id}`}>{device.save_error ?? ''}</div>
          <div data-testid={`device-runtime-apply-${device.id}`}>{device.runtime_apply_status ?? ''}</div>
          <div data-testid={`device-availability-${device.id}`}>{device.availability_status ?? ''}</div>
          <div data-testid={`device-running-${device.id}`}>{String(device.running ?? '')}</div>
        </div>
      ))}
      <button
        type="button"
        data-testid="seed-dev-01"
        onClick={() => actions.dispatch({
          type: 'addDevice',
          device: {
            id: 'dev-01',
            name: '',
            description: 'Modbus TCP PLC (Line A 主控)',
            protocol: 'modbus_tcp',
            config: {
              host: '192.168.1.100',
              port: 502,
              slave_id: 1,
              timeout: 5,
            },
            status: 'draft',
            test: null,
            persisted: false,
            save_state: 'idle',
            save_error: null,
            runtime_apply_status: null,
            runtime_apply_message: null,
            availability_status: 'available',
            availability_reason: null,
            running: false,
          },
        })}
      />
      <button
        type="button"
        data-testid="make-dev-01-valid"
        onClick={() => actions.dispatch({
          type: 'updateDevice',
          deviceId: 'dev-01',
          patch: { name: 'Line A Saved' },
        })}
      />
      <button
        type="button"
        data-testid="make-dev-01-invalid"
        onClick={() => actions.dispatch({
          type: 'updateDevice',
          deviceId: 'dev-01',
          patch: { name: '' },
        })}
      />
      <button
        type="button"
        data-testid="make-dev-A-fail"
        onClick={() => actions.dispatch({
          type: 'updateDevice',
          deviceId: 'dev-A',
          patch: { name: 'Line A Broken' },
        })}
      />
      <button
        type="button"
        data-testid="make-dev-B-valid"
        onClick={() => actions.dispatch({
          type: 'updateDevice',
          deviceId: 'dev-B',
          patch: { name: 'Line B Saved' },
        })}
      />
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
    updateAvailability: vi.fn(),
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

vi.mock('../../../src/services/studioV2Mappings', () => ({
  studioV2MappingsAPI: {
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

describe('DatalinkWorkbenchV2Page device autosave orchestration', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2RulesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(null);
    vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValue([]);
  });

  it('hydrates persisted devices in backend order on reload', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce({
      id: 'workspace-1',
      kind: 'single',
      status: 'ready',
      ordered_device_ids: ['dev-B', 'dev-A'],
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([
      {
        id: 'dev-B',
        name: 'Line B PLC',
        description: '',
        protocol: 'modbus_tcp',
        status: 'draft',
        connection_config: '{"host":"192.168.10.11","port":502,"slave_id":2,"timeout":5}',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
      {
        id: 'dev-A',
        name: 'Line A PLC',
        description: '',
        protocol: 'modbus_tcp',
        status: 'draft',
        connection_config: '{"host":"192.168.10.10","port":502,"slave_id":1,"timeout":5}',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
    ]);

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('device-order')).toHaveTextContent('dev-B,dev-A');
    });
  });

  it('hydrates persisted running truth on reload', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce({
      id: 'workspace-1',
      kind: 'single',
      status: 'ready',
      ordered_device_ids: ['dev-A'],
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([
      {
        id: 'dev-A',
        name: 'Line A PLC',
        description: '',
        protocol: 'modbus_tcp',
        status: 'active',
        connection_config: '{"host":"192.168.10.10","port":502,"slave_id":1,"timeout":5}',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        availability_status: 'available',
        running: true,
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
    ] as unknown as Awaited<ReturnType<typeof studioV2WorkspaceDevicesAPI.list>>);

    renderPage();

    await waitFor(() => {
      expect(screen.getByTestId('device-running-dev-A')).toHaveTextContent('true');
    });
  });

  it('saves one valid local device and marks it saved', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({
      id: 'workspace-1',
      kind: 'single',
      status: 'empty',
      ordered_device_ids: [],
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([]);
    vi.mocked(studioV2WorkspaceDevicesAPI.create).mockResolvedValueOnce({
      id: 'dev-01',
      name: 'Line A Saved',
      description: 'Modbus TCP PLC (Line A 主控)',
      protocol: 'modbus_tcp',
      status: 'draft',
      connection_config: '{"host":"192.168.1.100","port":502,"slave_id":1,"timeout":5}',
      last_test_at: null,
      last_test_success: null,
      last_test_error: '',
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
      runtime_apply_status: 'not_running',
    });

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
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce({
      id: 'workspace-1',
      kind: 'single',
      status: 'empty',
      ordered_device_ids: [],
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });
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
  it('marks session draft recovery state when a local device draft is not yet persisted', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce({
      id: 'workspace-1',
      kind: 'single',
      status: 'empty',
      ordered_device_ids: [],
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });
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
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({
      id: 'workspace-1',
      kind: 'single',
      status: 'ready',
      ordered_device_ids: ['dev-01'],
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([
      {
        id: 'dev-01',
        name: 'Line A PLC',
        description: '',
        protocol: 'modbus_tcp',
        status: 'active',
        connection_config: '{"host":"192.168.10.10","port":502,"slave_id":1,"timeout":5}',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
        running: true,
      },
    ]);
    vi.mocked(studioV2WorkspaceDevicesAPI.updateAvailability).mockResolvedValueOnce({
      id: 'dev-01',
      name: 'Line A PLC',
      description: '',
      protocol: 'modbus_tcp',
      status: 'active',
      connection_config: '{"host":"192.168.10.10","port":502,"slave_id":1,"timeout":5}',
      last_test_at: null,
      last_test_success: null,
      last_test_error: '',
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
      availability_status: 'unavailable',
      availability_reason: 'device form is invalid',
      running: false,
    } as unknown as Awaited<ReturnType<typeof studioV2WorkspaceDevicesAPI.updateAvailability>>);

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

  it('isolates save failure per device without blocking another valid save', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({
      id: 'workspace-1',
      kind: 'single',
      status: 'ready',
      ordered_device_ids: ['dev-A', 'dev-B'],
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([
      {
        id: 'dev-A',
        name: 'Line A PLC',
        description: '',
        protocol: 'modbus_tcp',
        status: 'draft',
        connection_config: '{"host":"192.168.10.10","port":502,"slave_id":1,"timeout":5}',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
      {
        id: 'dev-B',
        name: 'Line B PLC',
        description: '',
        protocol: 'modbus_tcp',
        status: 'draft',
        connection_config: '{"host":"192.168.10.11","port":502,"slave_id":2,"timeout":5}',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
    ]);
    vi.mocked(studioV2WorkspaceDevicesAPI.update).mockImplementation(async (deviceId, payload) => {
      if (deviceId === 'dev-A') {
        throw new Error('save failed');
      }

      return {
        id: deviceId,
        name: String(payload.name),
        description: '',
        protocol: 'modbus_tcp',
        status: 'draft',
        connection_config: '{"host":"192.168.10.11","port":502,"slave_id":2,"timeout":5}',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
        runtime_apply_status: 'not_running',
      };
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
      expect(screen.getByTestId('device-save-error-dev-A')).toHaveTextContent('save failed');
    });
  });

  it('keeps apply_failed visible instead of flattening it into saved', async () => {
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({
      id: 'workspace-1',
      kind: 'single',
      status: 'ready',
      ordered_device_ids: ['dev-A'],
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([
      {
        id: 'dev-A',
        name: 'Line A PLC',
        description: '',
        protocol: 'modbus_tcp',
        status: 'active',
        connection_config: '{"host":"192.168.10.10","port":502,"slave_id":1,"timeout":5}',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
    ]);
    vi.mocked(studioV2WorkspaceDevicesAPI.update).mockResolvedValueOnce({
      id: 'dev-A',
      name: 'Line A Broken',
      description: '',
      protocol: 'modbus_tcp',
      status: 'active',
      connection_config: '{"host":"192.168.10.20","port":502,"slave_id":1,"timeout":5}',
      last_test_at: null,
      last_test_success: null,
      last_test_error: '',
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
      runtime_apply_status: 'apply_failed',
      runtime_apply_message: 'runtime sync failed',
    });

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
    vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue({
      id: 'workspace-1',
      kind: 'single',
      status: 'ready',
      ordered_device_ids: ['dev-A'],
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
    });
    vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([
      {
        id: 'dev-A',
        name: 'Line A PLC',
        description: '',
        protocol: 'modbus_tcp',
        status: 'active',
        connection_config: '{"host":"192.168.10.10","port":502,"slave_id":1,"timeout":5}',
        last_test_at: null,
        last_test_success: null,
        last_test_error: '',
        created_at: '2026-05-30T00:00:00Z',
        updated_at: '2026-05-30T00:00:00Z',
      },
    ]);
    vi.mocked(studioV2WorkspaceDevicesAPI.update).mockResolvedValueOnce({
      id: 'dev-A',
      name: 'Line A Broken',
      description: '',
      protocol: 'modbus_tcp',
      status: 'active',
      connection_config: '{"host":"192.168.10.20","port":502,"slave_id":1,"timeout":5}',
      last_test_at: null,
      last_test_success: null,
      last_test_error: '',
      created_at: '2026-05-30T00:00:00Z',
      updated_at: '2026-05-30T00:00:00Z',
      runtime_apply_status: 'apply_failed',
      runtime_apply_message: 'runtime sync failed',
    });

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
