import type { WriteGroup } from '../../../src/types/studioV2WriteGroup';
import type { WriteGroupLifecycleHistory } from '../../../src/types/studioV2WriteGroupLifecycle';
import { isRecord } from '../../../src/utils/studioV2LifecycleJsonPrimitives';
import { parseStudioV2WriteGroup } from '../../../src/utils/studioV2WriteGroupDefinitionJson';
import type { LifecycleMockCall, LifecycleMockRequest, LifecycleMockReply } from './recording-lifecycle-api';
import { fixtureLifecycle, fixtureReadiness, type LifecycleMockState } from './recording-lifecycle-state';
import { MockContractError } from './recording-lifecycle-successor';
import { quietDelivery } from '../../fixtures/cutoverFixtures';
import { AUTHORING_HASH as HASH, AUTHORING_TIME as TIME } from './recording-authoring-state';
import { initialGuidedGroups, recordingGuidedHostState, advanceGuidedHead, GUIDED_PATH as PATH, type GuidedScenario, type GuidedMockState } from './recording-guided-state';
import { GuidedRecordingOperations, guidedBody, guidedCAS, GUIDED_CAS_KEYS, rejectGuided as reject, type GuidedStartOutcome } from './recording-guided-operations';
import { guidedRuleEvidence } from './recording-guided-effective-evidence';

type Effect = 'reject' | 'lose';
/** Strict TEST-ONLY page.route simulator: one canonical head; no SQL/device/runtime verification. */
export function createGuidedRecordingMock(scenario: GuidedScenario = 'new') {
  const state: GuidedMockState = { scenario, workspaceRevision: scenario === 'saved-different' ? 'W2' : 'W1', groups: initialGuidedGroups(scenario) };
  const host = recordingGuidedHostState(scenario), operations = new GuidedRecordingOperations(state);
  const effects = new Map<string, Effect>(), history = new Map<string, WriteGroupLifecycleHistory>();
  const ledger: LifecycleMockCall[] = [];
  const group = (id: string) => state.groups.find(item => item.id === id) ?? reject('MOCK_GROUP_NOT_FOUND');
  function lifecycleState(selected: WriteGroup): LifecycleMockState {
    return { scenario: scenario === 'blocked-wide' ? 'blocked-wide' : 'applied-long', source: selected, workspaceRevision: state.workspaceRevision, lineageRevision: 1 };
  }
  function lifecycle(selected: WriteGroup) {
    const result = fixtureLifecycle(lifecycleState(selected), selected);
    if (!history.has(selected.id)) history.set(selected.id, structuredClone(result.history));
    // Status/CAS follow the saved head; original destination/revision backlog never follows later edits.
    const original = history.get(selected.id)!;
    result.history = { ...structuredClone(original), generations: original.generations.map(item => ({ ...structuredClone(item),
      status: selected.status, group_revision: selected.revision })) };
    return result;
  }
  state.groups.forEach(lifecycle);
  function actionOf(request: LifecycleMockRequest, path: string) {
    if (path === `${PATH}/recording-start`) return 'start';
    if (path.startsWith(`${PATH}/recording-start/operations/`)) return 'start-operation';
    if (path === `${PATH}/database-operations/test-writes`) return 'test-write-list';
    if (path.startsWith(`${PATH}/database-operations/`)) return 'schema-operation';
    if (path.endsWith('/circuit-profiles/ensure') || path.includes('/write-groups/basic/')) return 'ensure';
    if (path === `${PATH}/write-groups`) return request.method === 'GET' ? 'list' : 'unsupported';
    if (path === `${PATH}/database-metadata`) return 'metadata';
    if (path.endsWith('/disable')) return 'stop';
    const match = path.match(/\/write-groups\/[^/]+$/);
    if (match) return request.method === 'DELETE' ? 'remove' : request.method === 'PUT' ? 'update' : 'get';
    return ['schema-preview', 'schema-apply', 'runtime-context', 'measurements', 'lifecycle', 'readiness', 'delivery', 'cutovers'].find(name => path.endsWith(`/${name}`)) ?? 'unsupported';
  }
  function readScope(params: URLSearchParams, allowed: string[] = []) {
    const keys = [...params.keys()];
    if (new Set(keys).size !== keys.length || keys.some(key => ![...allowed, 'workspace_id'].includes(key))) reject('MOCK_READ_SCOPE', 400);
    if (params.has('workspace_id') && params.get('workspace_id') !== 'ws-1') reject('MOCK_FOREIGN_WORKSPACE', 403);
  }
  function save(input: unknown, selected?: WriteGroup, basic = false) {
    const body = guidedBody(input, ['workspace_id', 'expected_workspace_revision', 'expected_connector_revision', 'group', ...(selected ? ['expected_group_revision'] : [])]);
    if (body.workspace_id !== 'ws-1' || body.expected_workspace_revision !== state.workspaceRevision || body.expected_connector_revision !== 'crev-1') reject('MOCK_REVISION_MISMATCH', 409);
    if (selected) {
      guidedCAS(state, selected, body);
      if (!lifecycle(selected).actions.some(item => item.action === 'edit' && item.allowed)) reject('MOCK_LIFECYCLE_BLOCKED', 409);
    } else if (!['new', 'long-sources', 'empty-device'].includes(scenario)) reject('MOCK_SAVED_CONVERSION_FORBIDDEN', 409);
    const draft = guidedBody(body.group, ['workspace_id', 'name', 'members', 'destination', 'row_policy', 'write_policy', 'layout_version', 'recording_profile', 'storage_mode']);
    const destination = guidedBody(draft.destination, ['connector_id', 'connector_revision', 'table_schema', 'table_name', 'storage_strategy']);
    if (draft.workspace_id !== 'ws-1' || destination.connector_id !== 'conn-1' || destination.connector_revision !== 'crev-1' ||
      destination.table_schema !== 'main' || !selected && destination.storage_strategy !== 'managed' ||
      draft.storage_mode && draft.storage_mode !== 'plain') reject('MOCK_DESTINATION_MISMATCH', 409);
    if (selected && (draft.layout_version !== selected.layout_version || destination.storage_strategy !== selected.destination.storage_strategy ||
      destination.table_name !== selected.destination.table_name)) reject('MOCK_SAVED_CONVERSION_FORBIDDEN', 409);
    if (!Array.isArray(draft.members) || !draft.members.length) return reject('MOCK_SOURCE_MISMATCH', 409);
    const sources = new Set<string>();
    for (const item of draft.members) {
      const member = guidedBody(item, ['device_id', 'point_id', 'tag_id', 'source_revision', 'mapping_revision', 'target_column', 'required',
        'entity_key', 'max_age_seconds', 'circuit_id', 'circuit_name', 'value_role', 'source_value_type', 'result_value_type', 'measurement_id']);
      const source = Object.values(host.mappings).find(mapping => mapping.device_id === member.device_id &&
        mapping.persisted_point_id === member.point_id && mapping.tag_id === member.tag_id);
      const key = JSON.stringify([member.device_id, member.point_id, member.tag_id]);
      if (!source || member.device_id !== 'dev-1' || sources.has(key) || member.measurement_id || member.value_role === 'counter_raw' ||
        member.source_value_type && member.source_value_type !== source.target_type || source.target_type === 'string' && member.value_role !== 'state' && !basic) reject('MOCK_SOURCE_MISMATCH', 409);
      sources.add(key);
    }
    if (!basic && (!isRecord(draft.recording_profile) || draft.recording_profile.kind !== 'instantaneous' || draft.recording_profile.unit !== 'second')) reject('MOCK_PROFILE_UNSUPPORTED', 409);
    const nextHead = `W${Number(state.workspaceRevision.slice(1)) + 1}`, id = selected?.id ?? `G${state.groups.length + 1}`;
    const parsed = parseStudioV2WriteGroup({ ...draft, id, workspace_id: 'ws-1', revision: `R${nextHead.slice(1)}`,
      applied_revision: selected?.applied_revision ?? '', status: 'draft', migration: selected?.migration ?? {}, created_at: selected?.created_at ?? TIME, updated_at: TIME,
      ...(selected?.basic_managed_device_id ? { basic_managed_device_id: selected.basic_managed_device_id }
        : !selected && basic ? { basic_managed_device_id: 'dev-1' } : {}),
      destination: { ...destination, database: '/mock-private/recordings.sqlite', table_name: destination.table_name || `mock_recording_${id}` },
      members: draft.members.map((item, index) => ({ ...(item as Record<string, unknown>), source_revision: `source-${index}`, mapping_revision: `mapping-${index}`,
        target_column: isRecord(item) && item.target_column ? item.target_column : `value_${index}` })),
    });
    if (!parsed) return reject('MOCK_INVALID_GROUP', 400);
    advanceGuidedHead(state);
    if (selected) state.groups.splice(state.groups.indexOf(selected), 1, parsed); else state.groups.push(parsed);
    return { workspace_revision: state.workspaceRevision, group: parsed };
  }
  function dispatch(request: LifecycleMockRequest, url: URL): unknown {
    const { pathname: path, searchParams: params } = url;
    if (!path.startsWith(`${PATH}/`)) return reject();
    if (request.method === 'GET') {
      if (request.body !== undefined) reject('MOCK_READ_BODY', 400);
      if (path === `${PATH}/runtime-context`) {
        readScope(params, ['workspace_revision', 'device_id']);
        const deviceId = params.get('device_id');
        if (params.get('workspace_id') !== 'ws-1' || deviceId === null || !host.devices.some(device => device.id === deviceId)) return reject('MOCK_RUNTIME_SCOPE', 403);
        if (params.get('workspace_revision') !== state.workspaceRevision) reject('MOCK_REVISION_MISMATCH', 409);
        return { workspace_id: 'ws-1', devices: [], default_device_id: null, setup: null, rule_evidence: guidedRuleEvidence(state, deviceId) };
      }
      if (path === `${PATH}/measurements`) {
        readScope(params, ['device_id']); if (!host.devices.some(device => device.id === params.get('device_id'))) reject('MOCK_DEVICE_SCOPE', 403); return [];
      }
      if (path === `${PATH}/database-metadata`) {
        readScope(params, ['expected_connector_revision', 'group_id', 'expected_group_revision']);
        const selected = params.has('group_id') ? group(params.get('group_id')!) : undefined;
        if (params.get('expected_connector_revision') !== 'crev-1' || selected && params.get('expected_group_revision') !== selected.revision ||
          !selected && params.has('expected_group_revision')) reject('MOCK_REVISION_MISMATCH', 409);
        return { workspace_id: 'ws-1', connector_id: 'conn-1', connector_revision: 'crev-1', database: '/mock-private/recordings.sqlite', schema: 'main',
          table: selected?.destination.table_name ?? 'mock_unsaved', inspection_status: selected?.destination.schema_digest ? 'exists' : 'missing',
          columns: (selected?.members ?? []).map(member => ({ name: member.target_column, data_type: 'REAL', nullable: true, primary_key: false })) };
      }
      readScope(params);
      if (path === `${PATH}/write-groups`) return { workspace_id: 'ws-1', workspace_revision: state.workspaceRevision, groups: state.groups };
      if (path === `${PATH}/database-operations/test-writes`) return { workspace_id: 'ws-1', complete: true, operations: [] };
      if (path.startsWith(`${PATH}/recording-start/operations/`)) return operations.operation(decodeURIComponent(path.slice(`${PATH}/recording-start/operations/`.length)));
      if (path.startsWith(`${PATH}/database-operations/`)) return operations.schemaOperation(decodeURIComponent(path.slice(`${PATH}/database-operations/`.length)));
      const match = path.slice(PATH.length).match(/^\/write-groups\/([^/]+)(?:\/(lifecycle|readiness|delivery|cutovers))?$/);
      if (!match) return reject();
      const selected = group(decodeURIComponent(match[1]));
      if (!match[2]) return { workspace_revision: state.workspaceRevision, group: selected };
      if (match[2] === 'lifecycle') return lifecycle(selected);
      if (match[2] === 'readiness') {
        const schemaReady = Boolean(selected.destination.schema_digest && selected.destination.schema_revision === selected.revision && !['disabled', 'deleted'].includes(selected.status));
        const blocked = scenario === 'blocked-wide', ready = schemaReady && !blocked;
        return { ...fixtureReadiness(lifecycleState(selected), selected), ready, config_ready: !blocked, schema_ready: schemaReady, schema_digest: schemaReady ? HASH : undefined,
          issues: ready ? [] : [{ code: blocked ? 'connector_identity_stale' : 'schema-confirmation-required', severity: 'blocking', step: 'Step 4', scope: selected.id,
            message: blocked ? 'Mock connector identity is stale.' : 'Mock schema confirmation is required.' }] };
      }
      if (match[2] === 'delivery') return { ...quietDelivery, group_id: selected.id, intake: { state: selected.status === 'running' ? 'active' : 'not_running' },
        stages: { ...quietDelivery.stages, queued: selected.applied_revision ? 1 : 0 } };
      return [];
    }
    if ([...params.keys()].length) reject('MOCK_WRITE_QUERY', 400);
    if (request.method === 'POST' && path === `${PATH}/recording-start`) return operations.start(request.body);
    if (request.method === 'POST' && path === `${PATH}/write-groups/circuit-profiles/ensure`) return save(request.body);
    if (request.method === 'POST' && path === `${PATH}/write-groups/basic/dev-1`) return save(request.body, undefined, true);
    const match = path.slice(PATH.length).match(/^\/write-groups\/([^/]+)(?:\/(schema-preview|schema-apply|disable))?$/);
    if (!match) return reject();
    const selected = group(decodeURIComponent(match[1]));
    if (request.method === 'PUT' && !match[2]) return save(request.body, selected, !selected.recording_profile);
    if (request.method === 'POST' && match[2] === 'schema-preview') return operations.preview(selected, request.body);
    if (request.method === 'POST' && match[2] === 'schema-apply') return operations.apply(selected, request.body);
    const remove = request.method === 'DELETE' && !match[2], stop = request.method === 'POST' && match[2] === 'disable';
    if (!remove && !stop) return reject();
    const body = guidedBody(request.body, GUIDED_CAS_KEYS); guidedCAS(state, selected, body);
    if (!lifecycle(selected).actions.some(item => item.action === (remove ? 'remove' : 'stop') && item.allowed)) reject('MOCK_LIFECYCLE_BLOCKED', 409);
    advanceGuidedHead(state); selected.revision = `R${state.workspaceRevision.slice(1)}`; selected.status = remove ? 'deleted' : 'disabled';
    return { workspace_revision: state.workspaceRevision, group: selected };
  }
  return {
    ledger,
    snapshot: () => ({ ...structuredClone(state), schemaPreviews: operations.previewCount, history: structuredClone([...history.values()]) }),
    rejectNext: (action: string) => { effects.set(action, 'reject'); },
    loseNext: (action: string) => { effects.set(action, 'lose'); },
    setStartOutcome: (outcome: GuidedStartOutcome) => operations.setStartOutcome(outcome),
    handle(request: LifecycleMockRequest): LifecycleMockReply {
      const url = new URL(request.url, 'http://fixture.invalid'), action = actionOf(request, url.pathname), effect = effects.get(action); effects.delete(action);
      let response: LifecycleMockReply;
      try {
        if (effect === 'reject') reject('MOCK_CONTROLLED_UNAVAILABLE', 503);
        response = { status: 200, body: { success: true, data: structuredClone(dispatch(request, url)) } };
      } catch (error) {
        if (!(error instanceof MockContractError)) throw error;
        response = { status: error.status, body: { success: false, error: { code: error.code, message: 'Test-only mock rejected request', action: 'reload', retryable: false, request_id: 'mock-request' } } };
      }
      // All attempted non-GETs count, including unknown/rejected writes and preview token creation.
      ledger.push({ ...structuredClone(request), action, mutation: request.method !== 'GET' && request.method !== 'HEAD', status: response.status });
      return { ...response, ...(effect === 'lose' && response.status === 200 ? { transport: 'abort' as const } : {}) };
    },
  };
}
export type GuidedRecordingMock = ReturnType<typeof createGuidedRecordingMock>;
