import { hydrateStudioV2Rule, toStudioV2RuleUpdateRequest } from '../../../src/features/datalink/workbench-v2/state/studioV2RuleAutosave';
import { compactSnapshot, snapshotDelivery, snapshotLifecycle, snapshotLocator, snapshotWorkspace } from './recording-home-snapshot-state';
export const SNAPSHOT_BASE = '/api/v1/datalink/studio-v2/workspace';
export const SNAPSHOT_SOURCE = '/api/v1/datalink/studio-v2/recording-home/sources';
export interface SnapshotRequest { method: string; url: string; body?: unknown }
export type SnapshotPhase = 'home' | 'editor' | 'source-modal';
export interface SnapshotCall extends SnapshotRequest {
  sequence: number; action: string; phase: SnapshotPhase; workspaceId: string; revision: string;
  requestedWorkspaceId?: string; mutation: boolean; unexpected: boolean; status: number; delivery: 'pending' | 'fulfilled' | 'finished' | 'cancelled' | 'not-delivered';
  response: unknown; groupId?: string; failure?: string;
}
export interface SnapshotReply { status: number; body: unknown; call: SnapshotCall; wait?: Promise<void> }
export interface SnapshotHold { entered: Promise<SnapshotCall>; release: () => void }
class Rejected extends Error {
  constructor(readonly code: string, readonly status = 404, readonly unexpected = false) { super(code); }
}
const reject = (code: string, status = 404, unexpected = false): never => { throw new Rejected(code, status, unexpected); };
const object = (value: unknown): value is Record<string, unknown> => typeof value === 'object' && value !== null && !Array.isArray(value);

/** Node-side simulator for the production /studio/v2 route. No test authority is injected into the app.
 * Each held response is cloned before waiting. Unknown methods/paths/query contracts fail the ledger.
 * Runtime facts remain unavailable and commands stay blocked, even when presentation is fresh.
 */
export function createRecordingHomeSnapshotAPI(options: { groupCount?: number } = {}) {
  const states = { A: snapshotWorkspace('A', options.groupCount), B: snapshotWorkspace('B', options.groupCount) };
  let active: 'A' | 'B' = 'A', phase: SnapshotPhase = 'home', unknownSources = 0;
  let approvedSave: { ruleId: string; prefix: string } | undefined;
  const ledger: SnapshotCall[] = [];
  const holds = new Map<string, { wait: Promise<void>; enter: (call: SnapshotCall) => void; release: () => void }[]>();
  const failures = new Map<string, number[]>();
  const activeHolds = new Set<() => void>();
  function actionOf(url: URL) {
    const path = url.pathname;
    if (path === SNAPSHOT_SOURCE) return 'sources';
    if (path === SNAPSHOT_BASE) return url.searchParams.get('entry') === 'locator' ? `locator:${url.searchParams.get('step')}` : 'bootstrap';
    if (path === `${SNAPSHOT_BASE}/write-groups`) return 'groups';
    const detail = path.match(/\/write-groups\/([^/]+)(?:\/(lifecycle|delivery|readiness|cutovers))?$/);
    if (detail) return detail[2] ?? 'group';
    if (path === `${SNAPSHOT_BASE}/source-rules` || path === `${SNAPSHOT_BASE}/devices` || path === `${SNAPSHOT_BASE}/mappings` || path === '/api/v1/datalink/points') return `raw:${path.split('/').at(-1)}`;
    if (path.startsWith(`${SNAPSHOT_BASE}/source-rules/`)) return 'save-rule';
    return path.slice('/api/v1/datalink/'.length);
  }
  function dispatch(request: SnapshotRequest, url: URL) {
    const state = states[active], { pathname: path, searchParams: params } = url;
    for (const key of ['workspace_id', 'locator_workspace_id']) {
      if (params.has(key) && params.get(key) !== active) reject('MOCK_FOREIGN_SCOPE', 409);
    }
    for (const key of ['expected_workspace_revision', 'workspace_revision', 'locator_workspace_revision']) {
      if (params.has(key) && params.get(key) !== state.revision) reject('MOCK_STALE_SCOPE', 409);
    }
    if (request.method === 'PUT' && path === `${SNAPSHOT_BASE}/source-rules/${approvedSave?.ruleId}`) {
      const body = request.body, saved = state.rules.find(rule => rule.id === approvedSave?.ruleId);
      if (!saved || !object(body)) return reject('MOCK_UNAPPROVED_RULE_BODY', 400, true);
      const expected = Object.fromEntries(Object.entries({ ...toStudioV2RuleUpdateRequest(hydrateStudioV2Rule(saved)),
        naming_prefix: approvedSave?.prefix }).filter(([, value]) => value !== undefined));
      if (Object.keys(body).length !== Object.keys(expected).length || Object.entries(expected).some(([key, value]) =>
        JSON.stringify(body[key]) !== JSON.stringify(value))) reject('MOCK_UNAPPROVED_RULE_BODY', 400, true);
      Object.assign(saved, body, { revision_id: `${saved.id}-r${++state.sourceVersion}` }); approvedSave = undefined;
      return { ...saved, runtime_apply_status: 'not_running' };
    }
    if (request.method !== 'GET') reject('MOCK_UNEXPECTED_WRITE', 404, true);
    if (path === SNAPSHOT_SOURCE) {
      if ([...params.keys()].length !== 1 || params.get('workspace_id') !== active) reject('MOCK_COMPACT_SCOPE', 400, true);
      if (unknownSources > 0) { unknownSources--; return { workspace_id: active, state: 'unknown', reason_code: 'source-changed' }; }
      return compactSnapshot(state);
    }
    if (path === SNAPSHOT_BASE) {
      if (!params.size) return state.workspace;
      if (params.get('entry') !== 'locator' || params.get('workspace_id') !== active) reject('MOCK_LOCATOR_QUERY', 400, true);
      const value = snapshotLocator(state, params);
      if (!value) reject('MOCK_LOCATOR_MISSING', 409);
      if (params.get('editor') === '1' && (params.get('expected_workspace_revision') !== state.revision ||
        params.get('expected_group_revision') !== state.groups.find(group => group.id === params.get('group_id'))?.revision)) reject('MOCK_STALE_LOCATOR', 409);
      return value;
    }
    if (path === `${SNAPSHOT_BASE}/write-groups`) return { workspace_id: active, workspace_revision: state.revision, groups: state.groups };
    if (path === `${SNAPSHOT_BASE}/database-operations/test-writes`) return { workspace_id: active, complete: true, operations: [] };
    if (path === `${SNAPSHOT_BASE}/runtime-context`) {
      if (!params.get('device_id') || !state.devices.some(device => device.id === params.get('device_id'))) reject('MOCK_RUNTIME_DEVICE', 400, true);
      return { workspace_id: active, devices: [], default_device_id: null, setup: null, rule_evidence: { contract_version: 1,
        workspace_id: active, workspace_revision: state.revision, device_id: params.get('device_id'), observed_at: null,
        freshness: 'unavailable', source: null, groups: [], locator: null } };
    }
    const detail = path.match(new RegExp(`^${SNAPSHOT_BASE}/write-groups/([^/]+)(?:/(lifecycle|delivery|readiness|cutovers))?$`));
    if (detail) {
      const group = state.groups.find(group => group.id === detail[1]);
      if (!group) return reject('MOCK_GROUP_MISSING', 404, true);
      if (params.has('expected_group_revision') && params.get('expected_group_revision') !== group.revision) reject('MOCK_STALE_GROUP', 409);
      if (detail[2] === 'lifecycle') return snapshotLifecycle(state, group);
      if (detail[2] === 'delivery') return snapshotDelivery(group);
      if (detail[2] === 'readiness') return { workspace_id: active, workspace_revision: state.revision, group_id: group.id,
        group_revision: group.revision, applied_revision: '', config_ready: true, schema_ready: false, ready: false,
        issues: [{ code: 'schema-confirmation-required', severity: 'blocking', step: 'Step 4', scope: group.id, message: 'Schema is not confirmed.' }] };
      if (detail[2] === 'cutovers') return [];
      return { workspace_revision: state.revision, group };
    }
    // Home must never silently fall back to the editor's complete raw transport.
    const raw = new Map<string, unknown>([[`${SNAPSHOT_BASE}/devices`, state.devices], [`${SNAPSHOT_BASE}/source-rules`, state.rules],
      [`${SNAPSHOT_BASE}/mappings`, state.mappings], ['/api/v1/datalink/points', state.points]]);
    if (raw.has(path)) {
      if (phase === 'home') reject('MOCK_RAW_HOME_TRANSPORT', 404, true);
      return raw.get(path);
    }
    if (path === '/api/v1/datalink/settings') return [{ key: 'theme', value: 'dark' }];
    if (path === '/api/v1/datalink/acquisition-anomalies/summary') return { items: state.devices.map(device => ({ device_id: device.id, active_count: 0 })), total_active: 0 };
    if (phase !== 'home' && path === '/api/v1/datalink/modbus-share/status') return { workspace_id: active, workspace_revision: state.revision,
      configured_enabled: false, enabled: false, running: false, port: 5020, address: '', bind_state: 'stopped', mapping_count: 0,
      hydration_state: 'pending', readiness: false };
    if (phase !== 'home' && path === '/api/v1/datalink/runtime/status') return { running: false, devices: [], groups: [] };
    if (phase === 'editor' && path === `/api/v1/datalink/points/${state.points[0].id}/format`) {
      const identities = { workspace_revision: state.revision, point_revision: 'point-r1', rule_revision: state.rules[0].revision_id, connection_revision: 'connection-r1' };
      if (params.get('device_id') !== state.devices[0].id) reject('MOCK_FORMAT_SCOPE', 409);
      return { workspace_id: active, device_id: state.devices[0].id, point_id: state.points[0].id, rule_id: state.rules[0].id,
        ...identities, workspace_revision_origin: 'persisted', editable: false, resolution: { ...identities, resolved: true,
          format_mode: 'legacy', configured_format: '', effective_format: '', effective_source: 'native', effective_policy: 'legacy_modbus',
          native_compatibility: true, applicable: false, format_version: 0 } };
    }
    if (phase === 'editor' && path === `${SNAPSHOT_BASE}/measurements`) return [];
    if (phase === 'editor' && path === `${SNAPSHOT_BASE}/database-config`) return null;
    if (phase === 'editor' && path === `${SNAPSHOT_BASE}/database-targets`) return [];
    reject('MOCK_UNEXPECTED_API', 404, true);
  }
  return {
    ledger,
    setPhase(next: SnapshotPhase) { phase = next; },
    snapshot: () => structuredClone(states[active]),
    switchWorkspace(id: 'A' | 'B', revision?: string) {
      active = id;
      if (revision) { states[id].revision = revision; states[id].devices[0].name = `Snapshot ${id} ${revision}`; states[id].sourceVersion++; }
    },
    allowRuleSave(ruleId: string, prefix: string) { approvedSave = { ruleId, prefix }; },
    unknownSourcesNext() { unknownSources++; },
    failNext(action: string, status = 503) { failures.set(action, [...failures.get(action) ?? [], status]); },
    holdNext(action = 'sources'): SnapshotHold {
      let release = () => {}, enter: (call: SnapshotCall) => void = () => {};
      const wait = new Promise<void>(resolve => { release = resolve; });
      const entered = new Promise<SnapshotCall>(resolve => { enter = resolve; });
      const releaseOnce = () => { release(); activeHolds.delete(releaseOnce); };
      activeHolds.add(releaseOnce);
      holds.set(action, [...holds.get(action) ?? [], { wait, enter, release: releaseOnce }]);
      return { entered, release: releaseOnce };
    },
    releaseAll() { for (const release of activeHolds) release(); holds.clear(); },
    handle(request: SnapshotRequest): SnapshotReply {
      const url = new URL(request.url, 'http://snapshot.fixture.invalid'), action = actionOf(url);
      const hold = holds.get(`${action}:${url.searchParams.get('workspace_id')}`)?.shift() ?? holds.get(action)?.shift();
      const failure = failures.get(action)?.shift();
      let status = 200, unexpected = false, body: unknown;
      try {
        if (failure) reject('MOCK_INJECTED_UNAVAILABLE', failure);
        body = { success: true, data: structuredClone(dispatch(request, url)) };
      } catch (error) {
        if (!(error instanceof Rejected)) throw error;
        status = error.status; unexpected = error.unexpected;
        body = { success: false, error: { code: error.code, message: 'RAW_FIXTURE_DIAGNOSTIC_MUST_NOT_REACH_DOM',
          retryable: true, action: 'reload', request_id: 'snapshot-fixture' } };
      }
      const call: SnapshotCall = { ...structuredClone(request), sequence: ledger.length + 1, action, phase, workspaceId: active,
        revision: states[active].revision, requestedWorkspaceId: url.searchParams.get('workspace_id') ?? url.searchParams.get('locator_workspace_id') ?? undefined, mutation: request.method !== 'GET' && request.method !== 'HEAD', unexpected, status,
        response: structuredClone(body), delivery: 'pending', groupId: url.pathname.match(/\/write-groups\/([^/]+)/)?.[1] };
      ledger.push(call); hold?.enter(call);
      return { status, body, call, wait: hold?.wait };
    },
  };
}
export type RecordingHomeSnapshotAPI = ReturnType<typeof createRecordingHomeSnapshotAPI>;
