import { render } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, vi } from 'vitest';
import DatalinkWorkbenchPage from '@/pages/datalink/workbench/DatalinkWorkbenchPage';
import type { Device, Mapping, Point, Tag } from '@/types/datalink';

const {
  mockDevices,
  mockPoints,
  mockTags,
  mockMappings,
  mockModbusShareAPI,
  mockDBTargetAPI,
} = vi.hoisted(() => ({
  mockDevices: [] as Device[],
  mockPoints: [] as Point[],
  mockTags: [] as Tag[],
  mockMappings: [] as Mapping[],
  mockModbusShareAPI: {
    status: vi.fn(), start: vi.fn(), stop: vi.fn(), listMappings: vi.fn(),
    upsertMapping: vi.fn(), deleteMapping: vi.fn(), writeTagValue: vi.fn(), sync: vi.fn(),
  },
  mockDBTargetAPI: {
    listConnectors: vi.fn(), getConnector: vi.fn(), createConnector: vi.fn(),
    updateConnector: vi.fn(), deleteConnector: vi.fn(), testConnector: vi.fn(),
    generateSchema: vi.fn(), dryRunMappings: vi.fn(), listTables: vi.fn(),
    validateConnector: vi.fn(), listMappings: vi.fn(), getMapping: vi.fn(),
    createMapping: vi.fn(), updateMapping: vi.fn(), deleteMapping: vi.fn(),
  },
}));

export {
  mockDBTargetAPI,
  mockDevices,
  mockMappings,
  mockModbusShareAPI,
  mockPoints,
  mockTags,
};

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/datalink', () => ({ modbusShareAPI: mockModbusShareAPI, dbTargetAPI: mockDBTargetAPI }));
vi.mock('@/hooks/datalink/useDevices', () => ({
  useDevicesQuery: () => ({ data: mockDevices, isLoading: false }),
  useCreateDeviceMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useUpdateDeviceMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useTestConnectionMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useTestDraftConnectionMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useToggleDeviceStatusMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));
vi.mock('@/hooks/datalink/usePoints', () => ({
  usePointsQuery: (filters?: { device_id?: string }) => ({
    data: filters?.device_id ? mockPoints.filter((point) => point.device_id === filters.device_id) : [],
    isLoading: false,
  }),
  useCreatePointMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDeletePointMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));
vi.mock('@/hooks/datalink/useTags', () => ({
  useTagsQuery: () => ({ data: mockTags, isLoading: false }),
  useCreateTagMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDeleteTagMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));
vi.mock('@/hooks/datalink/useMappings', () => ({
  useMappingsQuery: () => ({ data: mockMappings, isLoading: false }),
  useCreateMappingMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDeleteMappingMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));

export function renderPage() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <DatalinkWorkbenchPage />
    </QueryClientProvider>,
  );
}

beforeEach(() => {
  vi.clearAllMocks();
  mockDevices.splice(0, mockDevices.length, {
    id: 'device-1', name: 'Mixer PLC', description: '', protocol: 'modbus_tcp', status: 'active',
    connection_config: '{}', last_test_at: null, last_test_success: null, last_test_error: '', created_at: '', updated_at: '',
  });
  mockPoints.splice(0, mockPoints.length,
    { id: 'point-1', device_id: 'device-1', name: 'Flow Sensor', description: '', data_type: 'int16', address: '40001', enabled: true, polling_group_id: '', last_value: 12, last_read_at: '', last_error: '', error_count: 0, created_at: '', updated_at: '' },
    { id: 'point-2', device_id: 'device-1', name: 'Pressure Sensor', description: '', data_type: 'int16', address: '40002', enabled: true, polling_group_id: '', last_value: 28, last_read_at: '', last_error: '', error_count: 0, created_at: '', updated_at: '' },
  );
  mockTags.splice(0, mockTags.length,
    { id: 'tag-1', key: 'TAG_40001', display_name: 'Flow Sensor', description: '', data_type: 'int16', unit: '', labels: null, status: 'draft', created_at: '', updated_at: '' },
    { id: 'tag-2', key: 'TAG_40002', display_name: 'Pressure Sensor', description: '', data_type: 'int16', unit: '', labels: null, status: 'draft', created_at: '', updated_at: '' },
  );
  mockMappings.splice(0, mockMappings.length,
    { id: 'mapping-1', point_id: 'point-1', tag_id: 'tag-1', enabled: true, transform_pipeline: '', created_at: '', updated_at: '' },
    { id: 'mapping-2', point_id: 'point-2', tag_id: 'tag-2', enabled: true, transform_pipeline: '', created_at: '', updated_at: '' },
  );

  mockModbusShareAPI.status.mockResolvedValue({ enabled: false, port: 5020, address: '', bind_state: 'fail', mapping_count: 0 });
  mockModbusShareAPI.listMappings.mockResolvedValue([]);
  mockModbusShareAPI.start.mockResolvedValue({ enabled: true, port: 5030, address: '127.0.0.1:5030', bind_state: 'pass', mapping_count: 0 });
  mockModbusShareAPI.stop.mockResolvedValue({ enabled: false, port: 0, address: '', bind_state: 'fail', mapping_count: 0 });
  mockModbusShareAPI.upsertMapping.mockResolvedValue({ tag_id: 'tag-1', register: 12, data_type: 'int16', updated_at: '' });
  mockModbusShareAPI.sync.mockResolvedValue({ updated: 1, skipped: 0, errors: [] });
  mockModbusShareAPI.writeTagValue.mockResolvedValue(undefined);

  mockDBTargetAPI.listConnectors.mockResolvedValue([{
    id: 'connector-1', name: 'Main SQLite', kind: 'sqlite', connection_config: { dsn: '/tmp/target.db' },
    status: 'ready', last_check_at: '', last_check_error: '', enabled: true, created_at: '', updated_at: '',
  }]);
  mockDBTargetAPI.listMappings.mockResolvedValue([]);
  mockDBTargetAPI.listTables.mockResolvedValue([{
    schema: 'main', name: 'sensor_values', columns: [
      { name: 'ts', data_type: 'datetime', nullable: false, primary_key: true },
      { name: 'value', data_type: 'real', nullable: false, primary_key: false },
    ],
  }]);
  mockDBTargetAPI.validateConnector.mockResolvedValue({ ready: true, issues: [] });
  mockDBTargetAPI.generateSchema.mockResolvedValue({ connector_id: 'connector-1', dry_run: false, statements: ['CREATE TABLE sensor_values (...)'], executed: 1 });
  mockDBTargetAPI.dryRunMappings.mockResolvedValue({ connector_id: 'connector-1', results: [{ candidate_id: 'db-mapping-1', status: 'ready', mapping_id: 'db-mapping-1', tag_id: 'tag-1' }] });
  mockDBTargetAPI.createConnector.mockResolvedValue({
    id: 'connector-1', name: 'Main SQLite', kind: 'sqlite', connection_config: { dsn: '/tmp/target.db' },
    status: 'ready', last_check_at: '', last_check_error: '', enabled: true, created_at: '', updated_at: '',
  });
  mockDBTargetAPI.updateConnector.mockResolvedValue({
    id: 'connector-1', name: 'Main SQLite', kind: 'sqlite', connection_config: { dsn: '/tmp/target.db' },
    status: 'ready', last_check_at: '', last_check_error: '', enabled: true, created_at: '', updated_at: '',
  });
  mockDBTargetAPI.createMapping.mockResolvedValue({
    id: 'db-mapping-1', tag_id: 'tag-1', connector_id: 'connector-1', table_schema: 'main', table_name: 'sensor_values',
    column_name: 'value', write_mode: 'insert', timestamp_column: null, enabled: true, created_at: '', updated_at: '',
  });
});
