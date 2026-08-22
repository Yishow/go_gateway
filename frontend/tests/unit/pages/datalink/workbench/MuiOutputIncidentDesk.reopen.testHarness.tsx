import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render } from '@testing-library/react';
import { useEffect } from 'react';
import { vi } from 'vitest';
import { MuiOutputIncidentDesk } from '@/pages/datalink/workbench/MuiOutputIncidentDesk';
import { WorkbenchProvider, useWorkbench } from '@/pages/datalink/workbench/WorkbenchProvider';
import type { Device, Mapping, Point, SourceRuleRecord, Tag } from '@/types/datalink';
import type { SourceRuleCandidateSnapshotView } from '@/types/sourceRuleCandidates';

const mocks = vi.hoisted(() => ({
  devices: [] as Device[], points: [] as Point[], tags: [] as Tag[], mappings: [] as Mapping[],
  sourceRules: [] as SourceRuleRecord[],
  candidateViews: {} as Record<string, SourceRuleCandidateSnapshotView>,
  candidateState: { isLoading: false, isError: false, error: null as Error | null, refetch: vi.fn() },
  modbusShareAPI: {
    status: vi.fn(), start: vi.fn(), stop: vi.fn(), listMappings: vi.fn(), upsertMapping: vi.fn(),
    deleteMapping: vi.fn(), writeTagValue: vi.fn(), sync: vi.fn(),
  },
  dbTargetAPI: {
    listConnectors: vi.fn(), getConnector: vi.fn(), createConnector: vi.fn(), updateConnector: vi.fn(),
    deleteConnector: vi.fn(), testConnector: vi.fn(), generateSchema: vi.fn(), dryRunMappings: vi.fn(),
    listTables: vi.fn(), validateConnector: vi.fn(), listMappings: vi.fn(), getMapping: vi.fn(),
    createMapping: vi.fn(), updateMapping: vi.fn(), deleteMapping: vi.fn(),
  },
}));

export { mocks };

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, params?: Record<string, unknown>) => {
      const translations: Record<string, string> = {
        'workbench.output.incident.summaryEyebrow': 'Incident context',
        'workbench.output.incident.metrics.linked': 'linked',
        'workbench.output.incident.metrics.scoped': 'scoped',
        'workbench.output.incident.metrics.attention': 'needs review',
        'workbench.output.incident.handoff.label': 'Next step',
        'workbench.output.incident.handoff.noRule': 'Return to Tag to restore the active source rule, then continue output review.',
        'workbench.output.incident.handoff.ready': 'The desk context is aligned with the shared workboard for this target.',
      };
      if (translations[key] !== undefined) return translations[key];
      if (params && typeof params.defaultValue === 'string') return params.defaultValue;
      return key;
    },
  }),
}));

vi.mock('@/services/datalink', () => ({ modbusShareAPI: mocks.modbusShareAPI, dbTargetAPI: mocks.dbTargetAPI }));
vi.mock('@/hooks/datalink/useDevices', () => ({ useDevicesQuery: () => ({ data: mocks.devices, isLoading: false }) }));
vi.mock('@/hooks/datalink/usePoints', () => ({
  usePointsQuery: (filters?: { device_id?: string }) => ({
    data: filters?.device_id ? mocks.points.filter((point) => point.device_id === filters.device_id) : [],
    isLoading: false,
  }),
  useDeletePointMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));
vi.mock('@/hooks/datalink/useTags', () => ({
  useTagsQuery: () => ({ data: mocks.tags, isLoading: false, refetch: vi.fn().mockResolvedValue({ data: mocks.tags }) }),
  useCreateTagMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDeleteTagMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));
vi.mock('@/hooks/datalink/useMappings', () => ({
  useMappingsQuery: () => ({ data: mocks.mappings, isLoading: false }),
  useCreateMappingMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
  useDeleteMappingMutation: () => ({ mutateAsync: vi.fn(), isPending: false }),
}));
vi.mock('@/hooks/datalink/useSourceRules', () => ({
  useSourceRulesQuery: (filters?: { device_id?: string }) => ({
    data: filters?.device_id ? mocks.sourceRules.filter((rule) => rule.device_id === filters.device_id) : mocks.sourceRules,
    isLoading: false, isSuccess: true, isRefetching: false,
    refetch: vi.fn().mockResolvedValue({ data: mocks.sourceRules }),
  }),
}));
vi.mock('@/hooks/datalink/useSourceRuleCandidates', () => ({
  useSourceRuleCandidatesQuery: (ruleId?: string | null) => ({
    data: ruleId ? mocks.candidateViews[ruleId] ?? null : null,
    isLoading: mocks.candidateState.isLoading, isFetching: false,
    isError: mocks.candidateState.isError, error: mocks.candidateState.error, refetch: mocks.candidateState.refetch,
  }),
}));

function Bootstrap({ focusedRuleId }: { focusedRuleId?: string }) {
  const { setFocusedRuleId, setSelectedDeviceId } = useWorkbench();
  useEffect(() => {
    setSelectedDeviceId('device-1');
    if (focusedRuleId) setFocusedRuleId(focusedRuleId);
  }, [focusedRuleId, setFocusedRuleId, setSelectedDeviceId]);
  return null;
}

export function renderDesk(focusedRuleId = 'rule-1') {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  return render(
    <QueryClientProvider client={queryClient}>
      <WorkbenchProvider>
        <Bootstrap focusedRuleId={focusedRuleId} />
        <MuiOutputIncidentDesk />
      </WorkbenchProvider>
    </QueryClientProvider>,
  );
}

const readyDevice: Device = {
  id: 'device-1', name: 'Mixer PLC', description: '', protocol: 'modbus_tcp', status: 'active',
  connection_config: '{}', last_test_at: null, last_test_success: null, last_test_error: '',
  created_at: '', updated_at: '',
};
const readyPoints: Point[] = [
  { id: 'point-1', device_id: 'device-1', name: 'Flow Sensor', description: '', data_type: 'int16',
    address: '40001', enabled: true, polling_group_id: '', last_value: 12, last_read_at: '', last_error: '',
    error_count: 0, created_at: '', updated_at: '' },
  { id: 'point-2', device_id: 'device-1', name: 'Pressure Sensor', description: '', data_type: 'int16',
    address: '40002', enabled: true, polling_group_id: '', last_value: 28, last_read_at: '', last_error: '',
    error_count: 0, created_at: '', updated_at: '' },
];
const readyTags: Tag[] = [
  { id: 'tag-1', key: 'TAG_40001', display_name: 'Flow Sensor', description: '', data_type: 'int16', unit: '',
    labels: null, status: 'draft', created_at: '', updated_at: '' },
  { id: 'tag-2', key: 'TAG_40002', display_name: 'Pressure Sensor', description: '', data_type: 'int16', unit: '',
    labels: null, status: 'draft', created_at: '', updated_at: '' },
];
const readyMappings: Mapping[] = [
  { id: 'mapping-1', point_id: 'point-1', tag_id: 'tag-1', enabled: true, transform_pipeline: '',
    created_at: '', updated_at: '' },
  { id: 'mapping-2', point_id: 'point-2', tag_id: 'tag-2', enabled: true, transform_pipeline: '',
    created_at: '', updated_at: '' },
];
const readySourceRule: SourceRuleRecord = {
  id: 'rule-1', device_id: 'device-1', start_address: '40001', count: 2, data_type: 'int16',
  naming_prefix: 'SRC', enabled: true, locked: false, origin: 'manual', skipped_addresses: [],
  revision_id: 'rev-1', created_at: '', updated_at: '',
};
const readyCandidateView: SourceRuleCandidateSnapshotView = {
  source_rule_id: 'rule-1', revision_id: 'rev-1', tags: { status: 'ready', candidates: [] },
  local_modbus_outputs: {
    status: 'ready',
    candidates: [{
      id: 'modbus-candidate-1',
      identity: { source_rule_id: 'rule-1', candidate_type: 'local_modbus_outputs', candidate_kind: 'local_modbus_output', derived_from_rule_address: '40001' },
      proposed_signature: 'modbus-sig-1', address: '40001', point_id: 'point-1', tag_id: 'tag-1',
      tag_key: 'TAG_40001', display_name: 'Flow Sensor', data_type: 'int16', register: 12,
      register_count: 1, status: 'ready',
    }],
  },
  database_outputs: {
    status: 'ready',
    candidates: [{
      id: 'database-candidate-1',
      identity: { source_rule_id: 'rule-1', candidate_type: 'database_outputs', candidate_kind: 'database_output', derived_from_rule_address: '40002' },
      proposed_signature: 'db-sig-1', address: '40002', point_id: 'point-2', tag_id: 'tag-2',
      tag_key: 'TAG_40002', display_name: 'Pressure Sensor', data_type: 'int16', status: 'ready',
      connector_id: 'connector-1', table_schema: 'main', table_name: 'sensor_values', column_name: 'value',
      write_mode: 'insert',
    }],
  },
};

const readyConnector = {
  id: 'connector-1', name: 'Main SQLite', kind: 'sqlite', connection_config: { dsn: '/tmp/target.db' },
  status: 'ready', last_check_at: '', last_check_error: '', enabled: true, created_at: '', updated_at: '',
};
const readyTable = {
  schema: 'main', name: 'sensor_values',
  columns: [
    { name: 'ts', data_type: 'datetime', nullable: false, primary_key: true, unique: true },
    { name: 'value', data_type: 'real', nullable: false, primary_key: false, unique: false },
  ],
};
const readyDatabaseMapping = {
  id: 'db-mapping-1', tag_id: 'tag-2', connector_id: 'connector-1', table_schema: 'main',
  table_name: 'sensor_values', column_name: 'value', write_mode: 'insert', timestamp_column: null,
  enabled: true, created_at: '', updated_at: '',
};

export function seedReadyOutputFlow() {
  mocks.devices.splice(0, mocks.devices.length, { ...readyDevice });
  mocks.points.splice(0, mocks.points.length, ...readyPoints.map((point) => ({ ...point })));
  mocks.tags.splice(0, mocks.tags.length, ...readyTags.map((tag) => ({ ...tag })));
  mocks.mappings.splice(0, mocks.mappings.length, ...readyMappings.map((mapping) => ({ ...mapping })));
  mocks.sourceRules.splice(0, mocks.sourceRules.length, { ...readySourceRule });
  mocks.candidateViews['rule-1'] = readyCandidateView;
}

export function resetIncidentMocks() {
  vi.clearAllMocks();
  mocks.devices.splice(0); mocks.points.splice(0); mocks.tags.splice(0); mocks.mappings.splice(0);
  mocks.sourceRules.splice(0); Object.keys(mocks.candidateViews).forEach((key) => delete mocks.candidateViews[key]);
  mocks.candidateState.isLoading = false; mocks.candidateState.isError = false; mocks.candidateState.error = null;
  mocks.candidateState.refetch.mockResolvedValue({ data: null });
  mocks.modbusShareAPI.status.mockResolvedValue({ enabled: false, port: 5020, address: '', bind_state: 'fail', mapping_count: 1 });
  mocks.modbusShareAPI.listMappings.mockResolvedValue([{ tag_id: 'tag-1', register: 12, data_type: 'int16', updated_at: '' }]);
  mocks.modbusShareAPI.start.mockResolvedValue({ enabled: true, port: 5020, address: '127.0.0.1:5020', bind_state: 'pass', mapping_count: 1 });
  mocks.modbusShareAPI.stop.mockResolvedValue({ enabled: false, port: 0, address: '', bind_state: 'fail', mapping_count: 1 });
  mocks.modbusShareAPI.upsertMapping.mockResolvedValue({ tag_id: 'tag-1', register: 12, data_type: 'int16', updated_at: '' });
  mocks.modbusShareAPI.deleteMapping.mockResolvedValue(undefined);
  mocks.modbusShareAPI.writeTagValue.mockResolvedValue(undefined);
  mocks.modbusShareAPI.sync.mockResolvedValue({ updated: 1, skipped: 0, errors: [] });
  mocks.dbTargetAPI.listConnectors.mockResolvedValue([readyConnector]);
  mocks.dbTargetAPI.getConnector.mockResolvedValue(readyConnector);
  mocks.dbTargetAPI.listMappings.mockResolvedValue([]);
  mocks.dbTargetAPI.listTables.mockResolvedValue([readyTable]);
  mocks.dbTargetAPI.validateConnector.mockResolvedValue({ ready: true, issues: [] });
  mocks.dbTargetAPI.generateSchema.mockResolvedValue({ connector_id: 'connector-1', dry_run: false, statements: ['CREATE TABLE sensor_values (...)'], executed: 1 });
  mocks.dbTargetAPI.dryRunMappings.mockResolvedValue({ connector_id: 'connector-1', results: [{ candidate_id: 'database-candidate-1', status: 'ready', mapping_id: 'db-mapping-1', tag_id: 'tag-2' }] });
  mocks.dbTargetAPI.createConnector.mockResolvedValue(readyConnector);
  mocks.dbTargetAPI.updateConnector.mockResolvedValue(readyConnector);
  mocks.dbTargetAPI.deleteConnector.mockResolvedValue(undefined);
  mocks.dbTargetAPI.testConnector.mockResolvedValue({ success: true, latency_ms: 42, message: 'ok' });
  mocks.dbTargetAPI.createMapping.mockResolvedValue(readyDatabaseMapping);
  mocks.dbTargetAPI.updateMapping.mockResolvedValue(undefined);
  mocks.dbTargetAPI.deleteMapping.mockResolvedValue(undefined);
}
