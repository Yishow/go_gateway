import type { EvidenceScope } from '../../../src/types/runtimeRuleEvidence';
import type { RuntimeSourcePageRequest } from '../../../src/types/studioV2RuntimeContext';
import type { StudioV2LocatorEntrySnapshot } from '../../../src/types/studioV2LocatorEntry';
import { parseStudioEntityLocator } from '../../../src/services/studioEntityLocator';
import { nativeDevice, nativeMapping, nativeRule } from '../../unit/workbench-v2/studio-locator-fixtures';
import { quietDelivery } from '../../fixtures/cutoverFixtures';
import { savedGroup } from '../../fixtures/writeGroupState';
import { RUNTIME_RECORDING_PATH as PATH, RUNTIME_RECORDING_TIME as TIME, RuntimeRecordingFixtureError, runtimeRecordingEvidence,
  type RuntimeRecordingAction, type RuntimeRecordingCall, type RuntimeRecordingReply, type RuntimeRecordingRequest, type RuntimeRecordingScenario } from './runtime-recording-state';
import { runtimeRecordingRows, runtimeRecordingSourceEvidence, runtimeRecordingSourcePage } from './runtime-recording-source';
export type { RuntimeRecordingAction, RuntimeRecordingCall, RuntimeRecordingReply, RuntimeRecordingRequest, RuntimeRecordingScenario } from './runtime-recording-state';
const workspacePath = `${PATH}/studio-v2/workspace`;
const hostile = { password: 'fixture-secret-password', host: 'db.private:5432', database: '/private/runtime-fixture.sqlite', raw_error: 'driver-secret SELECT secret FROM credentials Server=localhost;Database=metrics [::1]:502' };
function reject(status = 400, code = 'unavailable'): never { throw new RuntimeRecordingFixtureError(status, code); }
function query(params: URLSearchParams, allowed: string[], required: string[] = []) {
  if ([...params.keys()].some(key => !allowed.includes(key) || params.getAll(key).length !== 1) || required.some(key => !params.get(key))) reject();
}
function device(value: string | null): 'A' | 'B' { return value === 'A' || value === 'B' ? value : reject(403, 'foreign-scope'); }
function scope(params: URLSearchParams): EvidenceScope {
  if (params.get('workspace_id') !== 'W') reject(403, 'foreign-scope');
  if (params.get('workspace_revision') !== 'W1') reject(409, 'stale-reference');
  return { workspace_id: 'W', workspace_revision: 'W1', device_id: device(params.get('device_id')) };
}
const deviceView = (id: 'A' | 'B') => ({ device_id: id, name: `Device ${id}`, protocol: 'modbus_tcp', running: true,
  availability_status: 'available', projection_alignment: 'aligned', runtime_projection_version: 'R1', workspace_projection_version: 'R2' });

/** TEST ONLY strict HTTP simulator. Its responses are never backend/device/SQL proof. */
export function createRuntimeRecordingMock(initial: RuntimeRecordingScenario = 'normal') {
  let scenario = initial;
  const ledger: RuntimeRecordingCall[] = [];
  const failures = new Map<RuntimeRecordingAction, number>();
  function classify(url: URL): RuntimeRecordingAction {
    const path = url.pathname;
    if (path === `${PATH}/runtime/status`) return 'status';
    if (path === `${PATH}/points`) return 'points';
    if (path === `${workspacePath}/recording-plans`) return 'reports';
    if (path === `${workspacePath}/write-groups/G1/delivery`) return 'native-delivery';
    if (path === `${workspacePath}/measurements/templates`) return 'native-templates';
    if (path === workspacePath) return url.searchParams.get('entry') === 'locator' ? 'locator' : 'unknown';
    if (path !== `${workspacePath}/runtime-context`) return 'unknown';
    if (url.searchParams.get('source_list') === 'rules') return 'rules';
    if (url.searchParams.get('source_list') === 'mappings') return 'mappings';
    if (url.searchParams.has('point_id')) return 'detail';
    return url.searchParams.size ? 'evidence' : 'context';
  }
  function context() {
    const rows = [runtimeRecordingRows(scenario, 'A'), runtimeRecordingRows(scenario, 'B')];
    return { workspace_id: 'W', workspace_revision: 'W1', devices: [deviceView('A'), deviceView('B')], default_device_id: 'A',
      setup: { evidence_state: 'saved-only', readiness_summary: { ready: true, blocking_count: 0, warning_count: 0, issues: [] },
        source_rules: rows.flatMap(row => row.rules), mappings: rows.flatMap(row => row.mappings), database_targets: [],
        database_config: { id: 'C1', workspace_id: 'W', name: '安全資料表', kind: 'sqlite', table: 'saved_setup_table', schema: '',
          write_mode: 'insert', write_interval_seconds: 30, timestamp_column: 'observed_at', status: 'ready', ...(scenario === 'hostile' ? hostile : {}) } },
      ...(scenario === 'hostile' ? hostile : {}) };
  }
  function status(id: 'A' | 'B') {
    const count = runtimeRecordingRows(scenario, id).mappings.length;
    return { running: true, uptime_seconds: 600, metrics: { collected_total: 100, write_success_total: 73, write_error_total: 0, mapping_error_total: 0, point_state_error_total: 0 },
      snapshot_state: { state: 'ready', empty: false, degraded: false, unavailable: false, stale: false },
      collectors: [{ ...deviceView(id), device_name: `Device ${id}`, status: 'running', points_total: count, points_healthy: count, points_stale: 0, points_error: 0,
        last_read_at: '2026-10-29T09:59:00Z', last_error: null, breaker_state: 'closed' }], database_delivery: [], modbus_share_delivery: [], diagnostics: [],
      ...(scenario === 'hostile' ? hostile : {}) };
  }
  function points(id: 'A' | 'B') {
    return runtimeRecordingRows(scenario, id).mappings.map(row => ({ id: row.point_id, device_id: id, name: row.display_name,
      description: '', data_type: 'float32', data_format: 'CDAB', address: row.address, enabled: true, polling_group_id: `PG${id}`,
      last_value: null, last_read_at: TIME, last_error: '', error_count: 0, created_at: TIME, updated_at: TIME }));
  }
  function locator(url: URL) {
    if (scenario === 'foreign') reject(403, 'foreign-scope');
    const intent = parseStudioEntityLocator(url.search);
    if (intent.kind !== 'locator') reject();
    const target = intent.locator;
    if (target.workspace_id !== 'W') reject(403, 'foreign-scope');
    const id = target.device_id === null ? null : device(target.device_id);
    const rows = runtimeRecordingRows(scenario, id ?? 'A');
    const mapping = target.mapping_id ? rows.mappings.find(row => row.id === target.mapping_id) : undefined;
    if (target.mapping_id && !mapping) reject(409, 'stale-reference');
    const ruleId = target.rule_id ?? mapping?.rule_id;
    const rule = ruleId ? rows.rules.find(row => row.id === ruleId) : undefined;
    if (ruleId && !rule) reject(409, 'stale-reference');
    if (target.group_id && !['G1', 'G2'].includes(target.group_id)) reject(404);
    if (target.group_id && id === 'B' && !['shared', 'quota-group', 'quota-global'].includes(scenario)) reject(403, 'foreign-scope');
    const entry: StudioV2LocatorEntrySnapshot = { workspace_revision: 'W1',
      hydration_scope: { complete_workspace: false, device_ids: id ? [id] : [], rule_ids: rule ? [rule.id] : [], mapping_ids: mapping?.id ? [mapping.id] : [], group_ids: target.group_id ? [target.group_id] : [] },
      rule_heads: rule ? { [rule.id]: 'R2' } : {}, group_heads: target.group_id ? { [target.group_id]: 'B' } : {}, revalidated_mapping_ids: mapping?.id ? [mapping.id] : [] };
    if (id) entry.devices = [{ ...nativeDevice(id), name: `Device ${id}` }];
    if (rule) entry.source_rules = [{ ...nativeRule(rule.id, id ?? 'A'), start_address: rule.start_address, count: rule.count, naming_prefix: rule.naming_prefix,
      data_type: 'float32', target_data_type: 'float64', scale_multiplier: 0.1, scale_offset: 0, data_format: 'CDAB' }];
    if (mapping?.id) entry.mappings = [{ ...nativeMapping(mapping.id), device_id: id ?? 'A', rule_id: mapping.rule_id, point_id: mapping.point_id,
      tag_id: mapping.tag_id ?? 'TA1', address: mapping.address, tag_key: mapping.tag_key ?? '', display_name: mapping.display_name ?? '', unit: 'm/s' }];
    if (target.group_id) entry.write_groups = { workspace_id: 'W', workspace_revision: 'W1', groups: [savedGroup({ id: target.group_id,
      workspace_id: 'W', revision: 'B', applied_revision: 'A', name: 'Shared recorder', status: scenario === 'deleted' && target.group_id === 'G1' ? 'deleted' : 'ready',
      members: [{ device_id: 'A', point_id: 'PA1', tag_id: 'TA1', source_revision: 'R2', mapping_revision: 'MR2', target_column: 'value', required: true },
        ...(['shared', 'quota-group', 'quota-global'].includes(scenario) ? [{ device_id: 'B', point_id: 'PB1', tag_id: 'TB1', source_revision: 'R2', mapping_revision: 'MR2', target_column: 'other', required: true }] : [])],
      destination: { connector_id: 'C1', connector_revision: '7', database: 'fixture-native.sqlite', table_schema: '', table_name: 'meter_new', storage_strategy: 'custom' } })] };
    return { id: 'W', kind: 'single', status: 'ready', ordered_device_ids: ['A', 'B'], created_at: TIME, updated_at: TIME, locator_entry: entry };
  }
  function dispatch(action: RuntimeRecordingAction, url: URL): unknown {
    const params = url.searchParams;
    if (action === 'native-templates') { query(params, []); return []; }
    if (action === 'native-delivery') {
      query(params, []); const delivery = runtimeRecordingEvidence(scenario, 'A').groups[0]?.delivery; const current = delivery?.current;
      if (!delivery || !current) reject(503);
      return { ...structuredClone(quietDelivery), group_id: 'G1', stages: current.stages, intake: { state: current.intake?.state ?? 'active' },
        last_sql_committed_at: current.last_commit?.committed_at ?? null, backlog: delivery.backlog.map(item => ({ group_revision: item.binding.entity_revision, connector_id: item.destination.connector_id, connector_revision: item.destination.connector_revision, table_schema: item.destination.schema ?? '', table_name: item.destination.table, pending: item.pending, error_codes: item.reasons })) };
    }
    if (action === 'context') { query(params, []); return context(); }
    if (action === 'locator') return locator(url);
    if (action === 'status' || action === 'points' || action === 'reports') {
      query(params, ['device_id'], ['device_id']); const id = device(params.get('device_id'));
      return action === 'status' ? status(id) : action === 'points' ? points(id) : [];
    }
    const required = ['workspace_id', 'workspace_revision', 'device_id'];
    const allowed = action === 'rules' || action === 'mappings' ? [...required, 'source_list', 'q', 'limit', 'cursor', 'rule_id', 'rule_revision']
      : action === 'detail' ? [...required, 'rule_id', 'rule_revision', 'point_id', 'mapping_id'] : required;
    query(params, allowed, [...required, ...(action === 'detail' ? ['rule_id', 'rule_revision', 'point_id'] : [])]);
    const selected = scope(params);
    if (action === 'rules' || action === 'mappings') {
      if (params.get('source_list') !== action) reject();
      if (params.has('limit') && !/^\d+$/.test(params.get('limit') ?? '')) reject();
      const request: RuntimeSourcePageRequest = { kind: action, query: params.get('q') ?? '', ...(params.has('limit') ? { limit: Number(params.get('limit')) } : {}),
        ...(params.has('cursor') ? { cursor: params.get('cursor') ?? '' } : {}), ...(params.has('rule_id') ? { rule_id: params.get('rule_id') ?? '' } : {}),
        ...(params.has('rule_revision') ? { rule_revision: params.get('rule_revision') ?? '' } : {}) };
      return { workspace_id: 'W', source_page: runtimeRecordingSourcePage(scenario, selected, request) };
    }
    const evidence = action === 'detail' ? runtimeRecordingSourceEvidence(scenario, selected,
      { rule_id: params.get('rule_id') ?? '', rule_revision: params.get('rule_revision') ?? '', point_id: params.get('point_id') ?? '',
        ...(params.has('mapping_id') ? { mapping_id: params.get('mapping_id') ?? '' } : {}) }) : runtimeRecordingEvidence(scenario, device(selected.device_id));
    return { workspace_id: 'W', rule_evidence: evidence, ...(scenario === 'hostile' ? hostile : {}) };
  }
  function handle(request: RuntimeRecordingRequest): RuntimeRecordingReply {
    const url = new URL(request.url, 'http://runtime-fixture.invalid'); const method = request.method.toUpperCase(); const action = classify(url);
    let reply: RuntimeRecordingReply;
    try {
      if (method !== 'GET') reject(405);
      if (action === 'unknown') reject(404);
      if (request.body !== undefined && request.body !== null && request.body !== '') reject();
      const failure = failures.get(action); if (failure !== undefined) { failures.delete(action); reject(failure); }
      reply = { status: 200, body: { success: true, data: dispatch(action, url) } };
    } catch (error) {
      const status = error instanceof RuntimeRecordingFixtureError ? error.status : 500;
      reply = { status, body: { success: false, error: { code: error instanceof RuntimeRecordingFixtureError ? error.code : 'unavailable', message: 'Read-only fixture response unavailable.' } } };
    }
    ledger.push({ method, path: url.pathname, action, status: reply.status, mutation: method !== 'GET', params: Object.fromEntries(url.searchParams) });
    return reply;
  }
  return { ledger, handle, setScenario(next: RuntimeRecordingScenario) { scenario = next; },
    failNext(action: RuntimeRecordingAction, status = 503) { failures.set(action, status); }, get scenario() { return scenario; } };
}
