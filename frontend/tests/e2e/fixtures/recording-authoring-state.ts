import type { WorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/types';
import type { WriteGroup, WriteGroupValueRole } from '../../../src/types/studioV2WriteGroup';
import type { MeasurementDefinition } from '../../../src/types/measurement';
import type { StudioV2MySQLRecordingCapability } from '../../../src/types/studioV2MySQLRecordingCapability';
import type { StudioV2SQLServerRecordingCapability } from '../../../src/types/studioV2SQLServerRecordingCapability';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';

// Test identities only. No database, device, sampled values or production capability is represented.
export const authoringScenarios = ['new-sqlite', 'new-postgres', 'new-mysql', 'new-mariadb', 'new-sqlserver',
  'saved-gauge', 'saved-state', 'saved-counter', 'saved-cumulative', 'test-write', 'test-write-sqlserver', 'saved-unsupported-wide'] as const;
export type AuthoringScenario = typeof authoringScenarios[number];
export const AUTHORING_PATH = '/api/v1/datalink/studio-v2/workspace';
export const AUTHORING_TIME = '2030-01-01T00:00:00.000Z';
export const AUTHORING_HASH = 'a'.repeat(64);
export function authoringScenarioFrom(value: string | null): AuthoringScenario {
  return authoringScenarios.find(item => item === value) ?? 'new-sqlite';
}
export function authoringKind(scenario: AuthoringScenario): 'sqlite' | 'postgres' | 'mysql' | 'sqlserver' {
  return scenario === 'new-mysql' || scenario === 'new-mariadb' ? 'mysql' : (scenario === 'new-sqlserver' || scenario === 'test-write-sqlserver' || scenario === 'saved-unsupported-wide') ? 'sqlserver'
    : scenario === 'new-postgres' ? 'postgres' : 'sqlite';
}
export function recordingAuthoringHostState(scenario: AuthoringScenario): WorkbenchV2State {
  const state = savedGroupState();
  state.devices[0].name = 'Mock PLC A';
  state.devices.push({ ...state.devices[0], id: 'dev-2', name: 'Mock PLC B' });
  state.points[0] = { ...state.points[0], name: 'Machine state', data_type: 'int16' };
  state.mappings[state.points[0].id] = { ...state.mappings[state.points[0].id], display_name: 'Machine state', target_type: 'int16' };
  state.points.push({ ...state.points[0], id: 'b-point-0', device_id: 'dev-2', name: 'B state' });
  state.mappings['b-point-0'] = { ...state.mappings[state.points[0].id], point_id: 'b-point-0', device_id: 'dev-2',
    persisted_point_id: 'b-pt-0', tag_id: 'b-tag-0', display_name: 'B state' };
  state.db.connector = { ...state.db.connector, kind: authoringKind(scenario), database: 'mock-only', schema: 'main' };
  state.settings.connectors = [{ id: 'conn-1', identity_revision: 'crev-1', name: 'Mock authoring target',
    kind: authoringKind(scenario), host: '', port: 0, database: 'mock-only', username: '', schema: 'main', table: '',
    enabled: true, status: 'unknown', default_write_interval_seconds: 15 }];
  return state;
}
export function initialAuthoringGroups(scenario: AuthoringScenario): WriteGroup[] {
  if (scenario.startsWith('new-')) return [];
  const cumulative = scenario === 'saved-cumulative';
  const role: WriteGroupValueRole = scenario === 'saved-state' ? 'state'
    : scenario === 'saved-counter' || cumulative ? 'counter_raw' : 'gauge';
  return [savedGroup({ id: 'G1', revision: 'R1', applied_revision: 'APPLIED1', status: 'running', name: 'Mock saved recording',
    layout_version: scenario === 'saved-unsupported-wide' ? 'managed_wide_profile_v1' : 'managed_long_v1',
    ...(scenario === 'saved-unsupported-wide' ? {} : { storage_mode: 'plain' as const }),
    recording_profile: { kind: cumulative ? 'cumulative' : 'instantaneous', interval: cumulative ? 10 : 15,
      unit: cumulative ? 'minute' : 'second', sample_reduction: cumulative ? 'latest' : 'mean' },
    destination: { connector_id: 'conn-1', connector_revision: 'crev-1', database: 'mock-only', table_schema: 'main',
      table_name: 'mock_recording_G1', storage_strategy: 'managed', schema_revision: 'SCHEMA1', schema_digest: AUTHORING_HASH },
    members: (cumulative ? [0] : [0, 1]).map(index => ({ device_id: 'dev-1', point_id: `pt-${index}`, tag_id: `tag-${index}`,
      source_revision: `source-${index}`, mapping_revision: `mapping-${index}`, target_column: `value_${index}`, required: false,
      circuit_id: 'cir-saved', circuit_name: 'Saved circuit', ...(cumulative ? {} : { value_role: index ? 'gauge' as const : role }),
      source_value_type: index ? 'float64' : 'int16', result_value_type: index ? 'float64' : role === 'gauge' ? 'decimal' : 'int64',
      ...(role === 'counter_raw' && index === 0 ? { measurement_id: 'MEAS-0' } : {}) })),
    row_policy: { interval_seconds: cumulative ? 600 : 15, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
    write_policy: { mode: 'append', dedupe_capability: 'receipt', recreate_missing: true },
  })];
}
export function authoringMeasurements(scenario: AuthoringScenario): MeasurementDefinition[] {
  if (scenario !== 'saved-counter' && scenario !== 'saved-cumulative') return [];
  return [{ id: 'MEAS-0', workspace_id: 'ws-1', device_id: 'dev-1', point_id: 'pt-0', tag_id: 'tag-0',
    equipment_id: 'cir-saved', definition_revision: 'DEF-1', source_binding_revision: 'BIND-1', series_epoch: 'EPOCH-1',
    name: 'Saved counter', quantity: 'count', unit: 'count', semantic_kind: 'counter', numeric_encoding: 'int16',
    counter_policy: { allow_negative: false }, created_at: AUTHORING_TIME, updated_at: AUTHORING_TIME }];
}
export function authoringCapability(scenario: AuthoringScenario): StudioV2MySQLRecordingCapability | StudioV2SQLServerRecordingCapability {
  const common = { connector_id: 'conn-1', connector_revision: 'crev-1', database: 'mock-only', capability_revision: 'MOCK-CAP-1',
    connected: true, metadata_status: 'verified', supports_managed_long: true, supports_native_partition: true,
    supports_group_test_write: true, supports_missing_repair: true, reasons: [] };
  return (scenario === 'new-sqlserver' || scenario === 'test-write-sqlserver' || scenario === 'saved-unsupported-wide') ? { ...common, schema: 'main', server_version: '16.0.4165.4',
    server_edition: 'Mock Developer', engine_edition: 3, compatibility_level: 160,
    permissions: { verified: true, database_view_definition: true, schema_view_definition: true, create_table: true,
      alter_schema: true, select_schema: true, insert_schema: true, delete_schema: true, alter_any_dataspace: true } }
    : { ...common, server_family: scenario === 'new-mariadb' ? 'mariadb' : 'mysql', server_version: 'mock-version', transactional_engine: true };
}
