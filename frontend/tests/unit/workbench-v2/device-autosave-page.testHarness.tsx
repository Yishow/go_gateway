import { vi } from 'vitest';
import type { StudioV2WorkspaceDeviceRecord } from '../../../src/services/studioV2WorkspaceDevices';
import type { StudioV2Workspace as WorkspaceType } from '../../../src/types/studioV2Workspace';

const studioV2ServiceMocks = vi.hoisted(() => ({
  workspace: { get: vi.fn() },
  devices: {
    list: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    updateAvailability: vi.fn(),
    remove: vi.fn(),
    updateOrder: vi.fn(),
  },
  rules: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() },
  mappings: { list: vi.fn(), create: vi.fn(), update: vi.fn(), remove: vi.fn() },
  database: {
    getConfig: vi.fn(),
    updateConfig: vi.fn(),
    listTargets: vi.fn(),
    upsertTarget: vi.fn(),
  },
}));

const studioV2WorkspaceAPI = studioV2ServiceMocks.workspace;
const studioV2WorkspaceDevicesAPI = studioV2ServiceMocks.devices;
const studioV2RulesAPI = studioV2ServiceMocks.rules;
const studioV2MappingsAPI = studioV2ServiceMocks.mappings;
const studioV2WorkspaceDatabaseAPI = studioV2ServiceMocks.database;

interface WorkbenchV2ShellMockProps<S> {
  state: S;
  actions: { dispatch: (action: unknown) => void };
}

export {
  studioV2WorkspaceAPI,
  studioV2WorkspaceDevicesAPI,
  studioV2RulesAPI,
  studioV2MappingsAPI,
  studioV2WorkspaceDatabaseAPI,
};

vi.mock('../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell', () => ({
  WorkbenchV2Shell: ({
    state,
    actions,
  }: WorkbenchV2ShellMockProps<{
    devices: Array<{
      id: string;
      name: string;
      save_state: string;
      save_error?: string | null;
      runtime_apply_status?: string | null;
      availability_status?: string;
      running?: boolean;
    }>;
  }>) => (
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
  studioV2WorkspaceAPI: studioV2ServiceMocks.workspace,
}));
vi.mock('../../../src/services/studioV2WorkspaceDevices', () => ({
  studioV2WorkspaceDevicesAPI: studioV2ServiceMocks.devices,
}));
vi.mock('../../../src/services/studioV2Rules', () => ({
  studioV2RulesAPI: studioV2ServiceMocks.rules,
}));
vi.mock('../../../src/services/studioV2Mappings', () => ({
  studioV2MappingsAPI: studioV2ServiceMocks.mappings,
}));
vi.mock('../../../src/services/studioV2WorkspaceDatabase', () => ({
  studioV2WorkspaceDatabaseAPI: studioV2ServiceMocks.database,
}));

export function workspaceFixture(overrides: Partial<WorkspaceType> = {}): WorkspaceType {
  return {
    id: 'workspace-1',
    kind: 'single',
    status: 'ready',
    ordered_device_ids: ['dev-A'],
    created_at: '2026-05-30T00:00:00Z',
    updated_at: '2026-05-30T00:00:00Z',
    ...overrides,
  };
}

export function deviceFixture(
  overrides: Partial<StudioV2WorkspaceDeviceRecord> = {},
): StudioV2WorkspaceDeviceRecord {
  return {
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
    ...overrides,
  };
}

export function resetDeviceAutosaveMocks() {
  vi.clearAllMocks();
  vi.mocked(studioV2WorkspaceDevicesAPI.list).mockResolvedValue([]);
  vi.mocked(studioV2RulesAPI.list).mockResolvedValue([]);
  vi.mocked(studioV2MappingsAPI.list).mockResolvedValue([]);
  vi.mocked(studioV2WorkspaceDatabaseAPI.getConfig).mockResolvedValue(null);
  vi.mocked(studioV2WorkspaceDatabaseAPI.listTargets).mockResolvedValue([]);
}
