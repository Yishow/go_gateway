import type { SchemaOperationStatus } from '../../../src/types/recordingPlan';
import { isRecord } from '../../../src/utils/studioV2LifecycleJsonPrimitives';
import { quietDelivery } from '../../fixtures/cutoverFixtures';
import { fixtureLifecycle, fixtureReadiness, initialLifecycleGroup, FIXTURE_PATH, type LifecycleMockState, type LifecycleScenario } from './recording-lifecycle-state';
import { LifecycleSuccessorMock, MockContractError } from './recording-lifecycle-successor';

export interface LifecycleMockRequest { method: string; url: string; body?: unknown }
export interface LifecycleMockReply { status: number; body: unknown; transport?: 'abort'; delayMs?: number }
/** mutation counts attempted writes/token creation, including rejected calls; it is not a SQL/DDL success flag. */
export interface LifecycleMockCall extends LifecycleMockRequest { action: string; mutation: boolean; status: number }
const readOnlyActions = new Set(['successor-review', 'activation-review', 'activation-readiness', 'operation-status']);
type Effect = { transport?: 'abort'; delayMs?: number; reject?: boolean };

/** Node-side page.route simulator. Unlisted paths/methods and stale/foreign scopes fail closed.
 * The ledger describes MOCK requests, not real SQL/database or runtime evidence.
 * An ACK-loss effect is applied AFTER state processing; the next status read sees the recorded mock result.
 */
export function createRecordingLifecycleMock(scenario: LifecycleScenario = 'applied-wide') {
  const state: LifecycleMockState = { scenario, source: initialLifecycleGroup(scenario), workspaceRevision: 'W1', lineageRevision: 1 };
  const successor = new LifecycleSuccessorMock(state);
  const ledger: LifecycleMockCall[] = [];
  const effects = new Map<string, Effect>();
  let staleLifecycle = false;
  let ordinaryPreviewCount = 0;
  const reject = (code: string, status = 409): never => { throw new MockContractError(code, status); };
  const group = (id: string) => id === 'G1' ? state.source : id === 'G2' && state.target ? state.target : reject('WRITE_GROUP_NOT_FOUND', 404);
  function actionOf(request: LifecycleMockRequest, path: string) {
    if (path.endsWith('/successor-review')) return 'successor-review';
    if (path.endsWith('/successor')) return isRecord(request.body) && typeof request.body.action === 'string' ? request.body.action : 'invalid-successor';
    if (path.endsWith('/schema-preview')) return 'ordinary-schema-preview';
    if (request.method === 'DELETE') return 'remove';
    if (path === `${FIXTURE_PATH}/write-groups`) return 'list';
    return path.split('/').at(-1) === 'G1' || path.split('/').at(-1) === 'G2' ? 'get' : path.split('/').at(-1) ?? 'unknown';
  }
  function dispatch(request: LifecycleMockRequest, url: URL): unknown {
    const { pathname: path, searchParams: params } = url;
    if (!path.startsWith(`${FIXTURE_PATH}/`)) return reject('MOCK_UNSUPPORTED_ROUTE', 404);
    if (request.method === 'GET') {
      if (path === `${FIXTURE_PATH}/runtime-context`) {
        const keys = ['workspace_id', 'workspace_revision', 'device_id'];
        if (request.body !== undefined || [...params.keys()].some(key => !keys.includes(key)) ||
          keys.some(key => params.getAll(key).length !== 1)) return reject('MOCK_INVALID_REQUEST', 400);
        if (params.get('workspace_id') !== 'ws-1' || params.get('device_id') !== 'dev-1') return reject('foreign-scope');
        if (params.get('workspace_revision') !== state.workspaceRevision) return reject('stale-reference');
        // This lifecycle simulator owns no verified runtime evidence; saved groups are not effective proof.
        return { workspace_id: 'ws-1', devices: [], default_device_id: null, setup: null,
          rule_evidence: { contract_version: 1, workspace_id: 'ws-1', workspace_revision: state.workspaceRevision,
            device_id: 'dev-1', observed_at: null, freshness: 'unavailable', source: null, groups: [], locator: null } };
      }
      if (path === `${FIXTURE_PATH}/database-operations/test-writes`) {
        // The workspace comes from this mock server session; the real list service sends no query/body.
        if (request.body !== undefined || params.size !== 0) return reject('MOCK_INVALID_REQUEST', 400);
        return { workspace_id: 'ws-1', complete: true, operations: [] };
      }
      if (path === `${FIXTURE_PATH}/write-groups`) return { workspace_id: 'ws-1', workspace_revision: state.workspaceRevision, groups: [state.source, ...(state.target ? [state.target] : [])] };
      if (path === `${FIXTURE_PATH}/measurements` && (!params.has('device_id') || params.get('device_id') === 'dev-1')) return [];
      if (path === `${FIXTURE_PATH}/database-metadata`) {
        const id = params.get('group_id');
        const selected = id ? group(id) : state.source;
        if (params.get('expected_connector_revision') !== selected.destination.connector_revision ||
          id && params.get('expected_group_revision') !== selected.revision) return reject('revision_mismatch');
        return { workspace_id: 'ws-1', connector_id: selected.destination.connector_id, connector_revision: selected.destination.connector_revision,
          database: selected.destination.database, schema: '', table: selected.destination.table_name, inspection_status: 'exists',
          columns: selected.members.map(member => ({ name: member.target_column || 'value', data_type: 'REAL', nullable: true, primary_key: false })) };
      }
      const match = path.slice(FIXTURE_PATH.length).match(/^\/write-groups\/([^/]+)(?:\/(lifecycle|readiness|delivery|cutovers))?$/);
      if (!match) return reject('MOCK_UNSUPPORTED_ROUTE', 404);
      const selected = group(decodeURIComponent(match[1]));
      if (!match[2]) return { workspace_revision: state.workspaceRevision, group: selected };
      if (match[2] === 'lifecycle') {
        if (staleLifecycle && selected.id === 'G1') { staleLifecycle = false; selected.revision = 'R1-STALE'; }
        return fixtureLifecycle(state, selected);
      }
      if (match[2] === 'readiness') return fixtureReadiness(state, selected);
      if (match[2] === 'delivery') return { ...quietDelivery, group_id: selected.id, intake: { state: selected.status === 'running' ? 'active' : 'not_running' } };
      return [];
    }
    if (request.method === 'POST' && path === `${FIXTURE_PATH}/write-groups/G1/successor-review`) return successor.review(request.body);
    if (request.method === 'POST' && path === `${FIXTURE_PATH}/write-groups/G1/successor`) {
      if (!isRecord(request.body) || typeof request.body.action !== 'string') return reject('MOCK_INVALID_REQUEST', 400);
      return successor.handle(request.body.action, request.body);
    }
    const ordinary = path.slice(FIXTURE_PATH.length).match(/^\/write-groups\/(G1|G2)\/schema-preview$/);
    if (request.method === 'POST' && ordinary) {
      // Explicit ordinary preview and standalone autoPreview create tokens; record both as mutations.
      const body = request.body, source = group(ordinary[1]);
      if (source.id === 'G2' && state.status?.cutover.status !== 'succeeded' && state.status?.first_start.status !== 'succeeded') return reject('WRITE_GROUP_LIFECYCLE_BLOCKED');
      if (!isRecord(body) || Object.keys(body).some(key => !['workspace_id', 'expected_workspace_revision', 'expected_group_revision', 'expected_connector_revision'].includes(key))) return reject('WRITE_GROUP_INVALID_REQUEST', 400);
      if (body.workspace_id !== 'ws-1' || body.expected_workspace_revision !== state.workspaceRevision || body.expected_group_revision !== source.revision ||
        body.expected_connector_revision !== source.destination.connector_revision || source.status === 'deleted') return reject('revision_mismatch');
      ordinaryPreviewCount++;
      return { token: `MOCK-ORDINARY-${ordinaryPreviewCount}`, operation_id: `ORDINARY-SCHEMA-${ordinaryPreviewCount}`, action: 'schema_apply',
        workspace_id: 'ws-1', workspace_revision: state.workspaceRevision, plan_id: source.id, plan_revision: source.revision,
        connector_id: source.destination.connector_id, connector_revision: source.destination.connector_revision, table_prefix: '', statements: [],
        tables: [{ name: source.destination.table_name, action: 'unchanged', columns: source.members.map(member => member.target_column) }],
        group_layout: { table_name: source.destination.table_name, owner_column: '__owner', columns: source.members.map(member => ({ name: member.target_column, sql_type: 'REAL', nullable: true, primary_key: false })) },
        digest: 'a'.repeat(64), expires_at: '2099-01-01T00:00:00Z', created_at: '2026-10-01T00:00:00Z' };
    }
    const remove = request.method === 'DELETE' && path === `${FIXTURE_PATH}/write-groups/G1`;
    const stop = request.method === 'POST' && path === `${FIXTURE_PATH}/write-groups/G1/disable`;
    if (remove || stop) {
      const body = request.body;
      if (!isRecord(body) || Object.keys(body).some(key => !['workspace_id', 'expected_workspace_revision', 'expected_group_revision', 'expected_connector_revision'].includes(key))) return reject('WRITE_GROUP_INVALID_REQUEST', 400);
      if (!isRecord(body) || body.workspace_id !== 'ws-1' || body.expected_workspace_revision !== state.workspaceRevision ||
        body.expected_group_revision !== state.source.revision || body.expected_connector_revision !== state.source.destination.connector_revision ||
        !fixtureLifecycle(state, state.source).actions.some(item => item.action === (remove ? 'remove' : 'stop') && item.allowed)) return reject('revision_mismatch');
      state.source.status = remove ? 'deleted' : 'disabled';
      state.source.revision += remove ? '-REMOVED' : '-STOPPED';
      state.workspaceRevision += remove ? '-REMOVED' : '-STOPPED';
      return { workspace_revision: state.workspaceRevision, group: state.source };
    }
    // No ordinary apply/ensure/start/schema-apply, arbitrary writes, or real backend passthrough.
    return reject('MOCK_UNSUPPORTED_ROUTE', 404);
  }
  return {
    ledger,
    snapshot: () => ({ ...structuredClone(state), ordinaryPreviewCount }),
    loseNext: (action: string) => effects.set(action, { ...effects.get(action), transport: 'abort' }),
    delayNext: (action: string, delayMs: number) => effects.set(action, { ...effects.get(action), delayMs }),
    rejectNext: (action: string) => effects.set(action, { reject: true }),
    staleNextLifecycle: () => { staleLifecycle = true; },
    setSchemaOperationStatus: (status: SchemaOperationStatus) => successor.setSchemaOperationStatus(status),
    setActivationReady: (ready: boolean) => successor.setActivationReady(ready),
    handle(request: LifecycleMockRequest): LifecycleMockReply {
      const url = new URL(request.url, 'http://fixture.invalid');
      const action = actionOf(request, url.pathname), effect = effects.get(action);
      effects.delete(action);
      let response: LifecycleMockReply;
      try {
        if (effect?.reject) throw new MockContractError('MOCK_CONTROLLED_UNAVAILABLE', 503);
        response = { status: 200, body: { success: true, data: structuredClone(dispatch(request, url)) } };
      } catch (error) {
        if (!(error instanceof MockContractError)) throw error;
        response = { status: error.status, body: { success: false, error: { code: error.code, message: 'Test-only mock rejected request', action: 'reload', retryable: false, request_id: 'mock-request' } } };
      }
      // Count forbidden write attempts too: a fail-closed 404 must not make a zero-mutation assertion pass.
      const readOnly = request.method === 'GET' || request.method === 'HEAD' || request.method === 'POST' &&
        (url.pathname === `${FIXTURE_PATH}/write-groups/G1/successor-review` && action === 'successor-review' ||
          url.pathname === pathForSuccessor() && readOnlyActions.has(action));
      ledger.push({ ...structuredClone(request), action, mutation: !readOnly, status: response.status });
      return { ...response, ...(effect?.transport ? { transport: effect.transport } : {}), ...(effect?.delayMs ? { delayMs: effect.delayMs } : {}) };
    },
  };
}
function pathForSuccessor() { return `${FIXTURE_PATH}/write-groups/G1/successor`; }
export type RecordingLifecycleMock = ReturnType<typeof createRecordingLifecycleMock>;
