import type { MeasurementDefinition } from '../../../src/types/measurement';
import type { WriteGroup } from '../../../src/types/studioV2WriteGroup';
import type { WriteGroupTestWriteStatus, WriteGroupTestWriteOutcome, WriteGroupTestCleanupStatus } from '../../../src/types/studioV2WriteGroupTestWrite';
import { isRecord } from '../../../src/utils/studioV2LifecycleJsonPrimitives';
import { parseStudioV2WriteGroup } from '../../../src/utils/studioV2WriteGroupDefinitionJson';
import { quietDelivery } from '../../fixtures/cutoverFixtures';
import { fixtureLifecycle, fixtureReadiness, type LifecycleMockState } from './recording-lifecycle-state';
import { LifecycleSuccessorMock, MockContractError } from './recording-lifecycle-successor';
import type { LifecycleMockCall, LifecycleMockRequest, LifecycleMockReply } from './recording-lifecycle-api';
import { AuthoringRecoveryMock } from './recording-authoring-recovery';
import { authoringCapability, authoringKind, authoringMeasurements, initialAuthoringGroups, AUTHORING_PATH as PATH,
  AUTHORING_HASH as HASH, AUTHORING_TIME as TIME, type AuthoringScenario } from './recording-authoring-state';

type Effect = { transport?: 'abort'; delayMs?: number; reject?: boolean };
const reject = (code = 'MOCK_UNSUPPORTED_ROUTE', status = 404): never => { throw new MockContractError(code, status); };
/** Strict TEST-ONLY route simulator. A response is not backend/SQL/runtime proof. */
export function createRecordingAuthoringMock(scenario: AuthoringScenario = 'new-sqlite') {
  const groups = initialAuthoringGroups(scenario), frozenGroups = structuredClone(groups);
  const measurements = authoringMeasurements(scenario);
  let workspaceRevision = 'W1', schemaPreviewCount = 0;
  let capability = authoringCapability(scenario);
  const recovery = new AuthoringRecoveryMock(authoringKind(scenario));
  const effects = new Map<string, Effect>(), measurementErrors = new Set<string>();
  const ledger: LifecycleMockCall[] = [];
  const lifecycleStates = new Map<string, LifecycleMockState>();
  const successors = new Map<string, LifecycleSuccessorMock>();
  const group = (id: string): WriteGroup => groups.find(item => item.id === id) ?? reject('MOCK_GROUP_NOT_FOUND');
  function lifecycleState(selected: WriteGroup) {
    let state = lifecycleStates.get(selected.id);
    if (!state) {
      state = { scenario: 'applied-long', source: selected, workspaceRevision, lineageRevision: 1 };
      lifecycleStates.set(selected.id, state); successors.set(selected.id, new LifecycleSuccessorMock(state));
    }
    state.workspaceRevision = workspaceRevision;
    return state;
  }
  function actionOf(request: LifecycleMockRequest, url: URL) {
    const path = url.pathname;
    if (path.endsWith('/measurements')) return request.method === 'GET' ? `measurements:${url.searchParams.get('device_id') ?? 'all'}` : 'measurement-confirm';
    if (path.endsWith('/capabilities')) return 'capability';
    if (path.endsWith('/database-operations/test-writes')) return 'test-write-list';
    if (path.includes('/database-operations/')) return 'operation-status';
    if (path.endsWith('/test-write-preview')) return 'test-write-preview';
    if (path.endsWith('/test-write')) return isRecord(request.body) && request.body.action === 'cleanup_repair' ? 'cleanup-repair' : 'test-write-confirm';
    if (path.endsWith('/circuit-profiles/ensure')) return 'ensure';
    if (path.endsWith('/schema-preview')) return 'schema-preview';
    if (path.endsWith('/successor-review')) return 'successor-review';
    if (path.endsWith('/write-groups')) return request.method === 'GET' ? 'list' : 'forbidden-create';
    return path.split('/').at(-1) ?? 'unknown';
  }
  function schemaPreview(selected: WriteGroup, body: unknown) {
    if (!isRecord(body) || Object.keys(body).some(key => !['workspace_id', 'expected_workspace_revision', 'expected_group_revision', 'expected_connector_revision'].includes(key)) ||
      body.workspace_id !== 'ws-1' || body.expected_workspace_revision !== workspaceRevision || body.expected_group_revision !== selected.revision ||
      body.expected_connector_revision !== selected.destination.connector_revision) return reject('MOCK_REVISION_MISMATCH', 409);
    schemaPreviewCount++;
    return { token: `MOCK-SCHEMA-${schemaPreviewCount}`, operation_id: `SCHEMA-${schemaPreviewCount}`, action: 'schema_apply',
      workspace_id: 'ws-1', workspace_revision: workspaceRevision, plan_id: selected.id, plan_revision: selected.revision,
      connector_id: selected.destination.connector_id, connector_revision: selected.destination.connector_revision, table_prefix: '',
      statements: [], tables: [{ name: selected.destination.table_name, action: 'unchanged', columns: selected.members.map(member => member.target_column) }],
      group_layout: { table_name: selected.destination.table_name, owner_column: '__mock_owner', columns: selected.members.map(member => ({
        name: member.target_column, sql_type: 'REAL', nullable: true, primary_key: false })) }, digest: HASH,
      expires_at: '2099-01-01T00:00:00Z', created_at: TIME };
  }
  function ensure(body: unknown) {
    if (!isRecord(body) || Object.keys(body).some(key => !['workspace_id', 'expected_workspace_revision', 'expected_connector_revision', 'group'].includes(key)) ||
      body.workspace_id !== 'ws-1' || body.expected_workspace_revision !== workspaceRevision || body.expected_connector_revision !== 'crev-1' ||
      !isRecord(body.group) || !isRecord(body.group.destination) || !Array.isArray(body.group.members) || !isRecord(body.group.recording_profile)) return reject('MOCK_INVALID_ENSURE', 400);
    if (!scenario.startsWith('new-')) return reject('MOCK_SAVED_MEANING_MUTATION_FORBIDDEN', 409);
    const draft = body.group, destination = body.group.destination, members = body.group.members;
    const dynamic = ['mysql', 'sqlserver'].includes(authoringKind(scenario));
    if (destination.connector_id !== 'conn-1' || destination.connector_revision !== 'crev-1' ||
      dynamic && (draft.layout_version !== 'managed_long_v1' || capability.connector_id !== 'conn-1' || capability.connector_revision !== 'crev-1' ||
        !capability.connected || capability.metadata_status !== 'verified' || !capability.supports_managed_long)) return reject('MOCK_LAYOUT_UNSUPPORTED', 409);
    if (!members.length || members.some(member => !isRecord(member) || member.device_id !== 'dev-1' ||
      !['pt-0', 'pt-1', 'pt-2'].includes(String(member.point_id)) || member.tag_id !== `tag-${String(member.point_id).at(-1)}`)) return reject('MOCK_SOURCE_MISMATCH', 409);
    const profileKind = body.group.recording_profile.kind;
    if (members.some(member => !isRecord(member) || member.tag_id === 'tag-2' && member.value_role !== 'state' ||
      (profileKind === 'cumulative' || member.value_role === 'counter_raw') && !measurements.some(def =>
        def.id === member.measurement_id && def.device_id === member.device_id && def.point_id === member.point_id && def.tag_id === member.tag_id))) {
      return reject('MOCK_UNCONFIRMED_SOURCE_MEANING', 409);
    }
    if (draft.storage_mode && draft.storage_mode !== 'plain' && !(draft.layout_version === 'managed_long_v1' &&
      (draft.storage_mode === 'timescale' && authoringKind(scenario) === 'postgres' ||
        draft.storage_mode === 'native_monthly' && dynamic && capability.supports_native_partition))) return reject('MOCK_STORAGE_UNSUPPORTED', 409);
    const id = `G${groups.length + 1}`;
    const parsed = parseStudioV2WriteGroup({ ...draft, id, workspace_id: 'ws-1', revision: `R${groups.length + 1}`, applied_revision: '', status: 'draft',
      name: draft.name || 'Mock authored recording', created_at: TIME, updated_at: TIME, migration: {},
      destination: { ...destination, database: 'mock-only', table_name: `mock_recording_${id}` },
      members: members.map((member, index) => ({ ...(member as Record<string, unknown>), source_revision: `source-${index}`, mapping_revision: `mapping-${index}`,
        target_column: isRecord(member) && member.target_column ? member.target_column : `value_${index}` })),
    });
    if (!parsed) return reject('MOCK_INVALID_AUTHORED_DEFINITION', 400);
    groups.push(parsed); workspaceRevision = `W${groups.length + 1}`;
    return { workspace_revision: workspaceRevision, group: parsed };
  }
  function dispatch(request: LifecycleMockRequest, url: URL): unknown {
    const { pathname: path, searchParams: params } = url;
    if (!path.startsWith(`${PATH}/`)) return reject();
    if (request.method === 'GET') {
      if (path === `${PATH}/runtime-context`) {
        const keys = [...params.keys()];
        if (keys.length !== 3 || keys.some(key => !['workspace_id', 'workspace_revision', 'device_id'].includes(key))) return reject('MOCK_RUNTIME_SCOPE', 400);
        if (params.get('workspace_id') !== 'ws-1' || !['dev-1', 'dev-2'].includes(params.get('device_id') ?? '')) return reject('foreign-scope', 403);
        if (params.get('workspace_revision') !== workspaceRevision) return reject('stale-reference', 409);
        // Authoring fixtures prove no runtime facts. Keep a valid, explicitly unavailable read model, never clone saved groups.
        return { workspace_id: 'ws-1', devices: [], default_device_id: null, setup: null, rule_evidence: {
          contract_version: 1, workspace_id: 'ws-1', workspace_revision: workspaceRevision, device_id: params.get('device_id'),
          observed_at: null, freshness: 'unavailable', source: null, groups: [], locator: null,
        } };
      }
      if (path === `${PATH}/write-groups`) return { workspace_id: 'ws-1', workspace_revision: workspaceRevision, groups };
      if (path === `${PATH}/measurements`) {
        const device = params.get('device_id');
        if (device !== 'dev-1' && device !== 'dev-2') return reject('MOCK_MEASUREMENT_SCOPE', 400);
        if (measurementErrors.has(device)) return reject('MOCK_MEASUREMENT_UNAVAILABLE', 503);
        return measurements.filter(item => item.device_id === device);
      }
      if (path === `${PATH}/recording-plans/capabilities`) {
        if (!['crev-1', 'crev-2'].includes(params.get('expected_connector_revision') ?? '') || params.get('connector_id') !== 'conn-1' ||
          params.get('kind') !== authoringKind(scenario)) return reject('MOCK_CAPABILITY_SCOPE', 400);
        // Deliberately does not follow a new request revision: stale response is a required negative fixture.
        return capability;
      }
      if (path.startsWith(`${PATH}/database-operations/`)) {
        // Current production callers use the server's workspace context, not a query argument.
        // This fixture has only ws-1. A caller-supplied foreign scope is never accepted.
        if (params.has('workspace_id') && params.get('workspace_id') !== 'ws-1') return reject('MOCK_FOREIGN_WORKSPACE', 403);
        if (path === `${PATH}/database-operations/test-writes`) return recovery.list();
        if (path === `${PATH}/database-operations/O1`) return recovery.status('O1');
        return reject('MOCK_OPERATION_NOT_FOUND', 404);
      }
      if (path === `${PATH}/database-metadata`) {
        const selected = params.get('group_id') ? group(params.get('group_id')!) : groups[0];
        if (params.get('expected_connector_revision') !== 'crev-1' || selected && params.has('expected_group_revision') && params.get('expected_group_revision') !== selected.revision) return reject('MOCK_REVISION_MISMATCH', 409);
        return { workspace_id: 'ws-1', connector_id: 'conn-1', connector_revision: 'crev-1', database: 'mock-only', schema: 'main',
          table: selected?.destination.table_name ?? 'mock_unsaved', inspection_status: 'exists',
          columns: (selected?.members ?? []).map(member => ({ name: member.target_column, data_type: 'REAL', nullable: true, primary_key: false })) };
      }
      const match = path.slice(PATH.length).match(/^\/write-groups\/([^/]+)(?:\/(lifecycle|readiness|delivery|cutovers))?$/);
      if (!match) return reject();
      const selected = group(decodeURIComponent(match[1]));
      if (!match[2]) return { workspace_revision: workspaceRevision, group: selected };
      if (match[2] === 'lifecycle') return fixtureLifecycle(lifecycleState(selected), selected);
      if (match[2] === 'readiness') {
        const ready = Boolean(selected.destination.schema_digest);
        return { ...fixtureReadiness(lifecycleState(selected), selected), schema_ready: ready, ready, schema_digest: ready ? HASH : undefined,
          issues: ready ? [] : [{ code: 'schema-confirmation-required', severity: 'blocking', step: 'Step 4', scope: selected.id, message: 'Mock schema confirmation is still required.' }] };
      }
      if (match[2] === 'delivery') return { ...quietDelivery, group_id: selected.id, intake: { state: selected.applied_revision && selected.status === 'running' ? 'active' : 'not_running' } };
      return [];
    }
    if (request.method === 'POST' && path === `${PATH}/measurements`) {
      const body = request.body;
      if (!isRecord(body) || body.workspace_id !== 'ws-1' || body.device_id !== 'dev-1' || body.point_id !== 'pt-0' || body.tag_id !== 'tag-0' ||
        !body.equipment_id || !body.unit || !['counter', 'signed_counter'].includes(String(body.semantic_kind))) return reject('MOCK_MEASUREMENT_SCOPE', 400);
      const saved = { ...body, id: 'MEAS-0', definition_revision: 'DEF-1', source_binding_revision: 'BIND-1', series_epoch: 'EPOCH-1',
        created_at: TIME, updated_at: TIME } as unknown as MeasurementDefinition;
      if (measurements.some(item => item.id === saved.id)) return reject('MOCK_DUPLICATE_MEASUREMENT', 409);
      measurements.push(saved); return saved;
    }
    if (request.method === 'POST' && path === `${PATH}/write-groups/circuit-profiles/ensure`) return ensure(request.body);
    const action = path.slice(PATH.length).match(/^\/write-groups\/([^/]+)\/(schema-preview|test-write-preview|test-write|successor-review)$/);
    if (request.method !== 'POST' || !action) return reject();
    const selected = group(decodeURIComponent(action[1]));
    if (action[2] === 'schema-preview') return schemaPreview(selected, request.body);
    if (action[2] === 'test-write-preview') {
      if (request.body !== undefined && request.body !== null) return reject('MOCK_INVALID_REQUEST', 400);
      return recovery.makePreview(selected);
    }
    if (action[2] === 'test-write') return recovery.confirm(selected, workspaceRevision, request.body);
    lifecycleState(selected);
    return successors.get(selected.id)!.review(request.body);
  }
  return {
    ledger,
    snapshot: () => ({ scenario, groups: structuredClone(groups), frozenGroups: structuredClone(frozenGroups), workspaceRevision,
      measurements: structuredClone(measurements), capability: structuredClone(capability), schemaPreviewCount, ...recovery.snapshot() }),
    loseNext: (action: string) => effects.set(action, { ...effects.get(action), transport: 'abort' }),
    rejectNext: (action: string) => effects.set(action, { ...effects.get(action), reject: true }),
    delayNext: (action: string, delayMs: number) => effects.set(action, { ...effects.get(action), delayMs }),
    setMeasurementError: (device: string, failed: boolean) => { if (failed) measurementErrors.add(device); else measurementErrors.delete(device); },
    setCapability: (patch: Partial<ReturnType<typeof authoringCapability>>) => { capability = { ...capability, ...patch }; },
    setTestWriteOutcome: (status: WriteGroupTestWriteStatus, write?: WriteGroupTestWriteOutcome, cleanup?: WriteGroupTestCleanupStatus) => recovery.setOutcome(status, write, cleanup),
    handle(request: LifecycleMockRequest): LifecycleMockReply {
      const url = new URL(request.url, 'http://fixture.invalid'), action = actionOf(request, url), effect = effects.get(action);
      effects.delete(action);
      let reply: LifecycleMockReply;
      try {
        if (effect?.reject) returnFailure();
        reply = { status: 200, body: { success: true, data: structuredClone(dispatch(request, url)) } };
      } catch (error) {
        if (!(error instanceof MockContractError)) throw error;
        reply = { status: error.status, body: { success: false, error: { code: error.code, message: 'Mock request rejected safely',
          action: 'reload', retryable: false, request_id: 'mock-authoring' } } };
      }
      const mutation = request.method !== 'GET' && request.method !== 'HEAD' && !(request.method === 'POST' && action === 'successor-review');
      ledger.push({ ...structuredClone(request), action, mutation, status: reply.status });
      return { ...reply, ...(effect?.transport ? { transport: effect.transport } : {}), ...(effect?.delayMs ? { delayMs: effect.delayMs } : {}) };
    },
  };
}
function returnFailure(): never { throw new MockContractError('MOCK_CONTROLLED_UNAVAILABLE', 503); }
export type RecordingAuthoringMock = ReturnType<typeof createRecordingAuthoringMock>;
