import type { RuntimeLogRecord, RuntimeLogSnapshot } from '../../../src/types/runtimeLogs';
import { projectRuntimeLogDomain } from '../../../src/services/runtimeLogDomainProjection';
import { parseEvidenceStage } from '../../../src/services/runtimeRuleEvidenceParser';
import { parseStudioEntityLocator } from '../../../src/services/studioEntityLocator';
import { safeRuntimeIdentity } from '../../../src/utils/runtimeOperatorSafety';
import { locatorRead } from '../../unit/workbench-v2/studio-locator-fixtures';
import { RUNTIME_LOGS_PATH, runtimeLogsScenarioSnapshot, type RuntimeLogsScenario } from './runtime-logs-state';
export type RuntimeLogsAction = 'snapshot' | 'locator' | 'delivery';
export interface RuntimeLogsRequest { method: string; url: string; body?: unknown }
export interface RuntimeLogsReply { status: number; body: unknown }
export interface RuntimeLogsCall { method: string; path: string; params: Record<string, string>; status: number; mutation: boolean }
const nativePath = '/api/v1/datalink/studio-v2/workspace';
class Rejected extends Error { constructor(readonly status: number, readonly unexpected = true) { super('Read-only fixture request rejected.'); } }
const reject = (status = 400): never => { throw new Rejected(status); };
function query(params: URLSearchParams, allowed: string[], required: string[] = []) {
  if ([...params.keys()].some(key => !allowed.includes(key) || params.getAll(key).length !== 1) || required.some(key => !params.get(key))) reject();
}
function validateLogQuery(params: URLSearchParams) {
  query(params, ['level', 'source', 'q', 'limit', 'workspace_id', 'device_id', 'rule_id', 'group_id', 'stage', 'actionable'], ['level']);
  if (!['debug', 'info', 'warn', 'error'].includes(params.get('level') ?? '') || params.has('limit') && !/^(?:[1-9]\d?|[1-4]\d\d|500)$/.test(params.get('limit') ?? '')) reject();
  for (const key of ['workspace_id', 'device_id', 'rule_id', 'group_id']) if (params.has(key) && !safeRuntimeIdentity(params.get(key))) reject();
  if (params.has('stage') && !parseEvidenceStage(params.get('stage'))) reject();
  if (params.has('actionable') && !['true', 'false'].includes(params.get('actionable') ?? '')) reject();
  if ((params.get('q')?.length ?? 0) > 256 || (params.get('source')?.length ?? 0) > 128) reject();
}
function matches(record: RuntimeLogRecord, params: URLSearchParams) {
  const levels = ['debug', 'info', 'warn', 'error'];
  if (levels.indexOf(record.level) < levels.indexOf(params.get('level') ?? 'debug')) return false;
  if (params.get('source') && record.source !== params.get('source')) return false;
  const scoped = ['workspace_id', 'device_id', 'rule_id', 'group_id', 'stage'].some(key => params.has(key)) || params.get('actionable') === 'true';
  if (scoped) {
    const view = projectRuntimeLogDomain(record);
    if (record.code === 'domain.unavailable' && record.source === 'runtime') {
      if (['workspace_id', 'device_id', 'rule_id', 'group_id'].some(key => params.has(key)) || params.has('stage') && params.get('stage') !== 'unavailable') return false;
    } else {
      if (view.kind !== 'domain') return false;
      for (const key of ['workspace_id', 'device_id', 'rule_id', 'group_id', 'stage']) if (params.has(key) && record.fields[key] !== params.get(key)) return false;
      if (params.get('actionable') === 'true' && view.action === 'none') return false;
    }
  }
  const literal = params.get('q')?.toLowerCase();
  return !literal || [record.message, record.code, ...Object.values(record.fields).filter(value => typeof value === 'string')]
    .some(value => String(value).toLowerCase().includes(literal));
}
function nativeLocator(url: URL, scenario: RuntimeLogsScenario) {
  const parsed = parseStudioEntityLocator(url.search);
  if (parsed.kind !== 'locator') return reject();
  const target = parsed.locator;
  if (target.workspace_id !== 'W') reject(403);
  if (target.device_id && !['D', 'A', 'B'].includes(target.device_id)) reject(404);
  if (target.group_id && !['G', 'G1', 'G2', 'B'].includes(target.group_id)) reject(404);
  if (target.rule_id && !['RULE', 'R'].includes(target.rule_id) || target.mapping_id && target.mapping_id !== 'M') reject(404);
  const read = locatorRead(target.step, target.group_id ?? target.device_id ?? 'D');
  read.workspace.ordered_device_ids = ['D', 'A', 'B'];
  if (read.entry.source_rules && target.rule_id) {
    read.entry.source_rules[0].id = target.rule_id;
    read.entry.rule_heads = { [target.rule_id]: 'R2' };
    read.entry.hydration_scope.rule_ids = [target.rule_id];
  }
  const group = read.entry.write_groups?.groups[0];
  if (group) {
    group.revision = scenario === 'renamed' ? 'R3' : 'R2';
    group.applied_revision = 'R1'; group.name = scenario === 'renamed' ? 'Renamed current recorder' : 'Shared recorder';
    read.entry.group_heads = { [group.id]: group.revision };
    if (scenario === 'deleted' && group.id === 'G1') group.status = 'deleted';
  }
  if (scenario === 'foreign') read.workspace.id = 'OTHER';
  return { ...read.workspace, locator_entry: read.entry };
}

/** Strict request-time HTTP simulator. Unknown reads and all commands are evidence failures, not fallbacks. */
export function createRuntimeLogsMock(initial: RuntimeLogsScenario = 'normal') {
  let scenario = initial; let snapshot = runtimeLogsScenarioSnapshot(initial);
  const ledger: RuntimeLogsCall[] = [], unexpected: string[] = [];
  const failures = new Map<RuntimeLogsAction, number>();
  function handle(request: RuntimeLogsRequest): RuntimeLogsReply {
    const url = new URL(request.url, 'http://runtime-logs-fixture.invalid'); const method = request.method.toUpperCase();
    let reply: RuntimeLogsReply;
    try {
      if (method !== 'GET') reject(405);
      if (request.body !== undefined && request.body !== null && request.body !== '') reject();
      const action: RuntimeLogsAction = url.pathname === RUNTIME_LOGS_PATH ? 'snapshot' : url.pathname === nativePath ? 'locator'
        : /^\/api\/v1\/datalink\/studio-v2\/workspace\/write-groups\/(G|G1|G2|B)\/delivery$/.test(url.pathname) ? 'delivery' : reject(404);
      let data: unknown;
      if (action === 'snapshot') {
        validateLogQuery(url.searchParams);
        data = { ...snapshot, records: snapshot.records.filter(record => matches(record, url.searchParams)).slice(-Number(url.searchParams.get('limit') ?? 200)) };
      } else if (action === 'locator') data = nativeLocator(url, scenario);
      else {
        query(url.searchParams, []); // Existing native delivery(id) GET carries only the exact group path.
      }
      const failure = failures.get(action);
      if (failure !== undefined) { failures.delete(action); throw new Rejected(failure, false); }
      if (action === 'delivery' || action === 'locator' && scenario === 'unavailable') throw new Rejected(503, false);
      reply = { status: 200, body: structuredClone(action === 'snapshot' ? data : { success: true, data }) };
    } catch (error) {
      const known = error instanceof Rejected;
      if (!known || error.unexpected) unexpected.push(`${method} ${url.pathname}${url.search}`);
      reply = { status: known ? error.status : 500, body: { success: false, error: { code: known && error.status === 403 && url.pathname === RUNTIME_LOGS_PATH ? 'logs_local_only' : 'unavailable', message: 'Read-only fixture response unavailable.' } } };
    }
    ledger.push({ method, path: url.pathname, params: Object.fromEntries(url.searchParams), status: reply.status, mutation: method !== 'GET' });
    return reply;
  }
  return { ledger, unexpected, handle, setScenario(next: RuntimeLogsScenario) { scenario = next; snapshot = runtimeLogsScenarioSnapshot(next); },
    setSnapshot(next: RuntimeLogSnapshot) { snapshot = structuredClone(next); },
    failNext(action: RuntimeLogsAction, status = 503) { failures.set(action, status); } };
}
export type RuntimeLogsMock = ReturnType<typeof createRuntimeLogsMock>;
