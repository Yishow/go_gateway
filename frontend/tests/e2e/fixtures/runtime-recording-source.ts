/** TEST ONLY deterministic identities shared by mock source lists and selected details. */
export function runtimeRecordingMappingIdentity(index: number, deviceId: 'A' | 'B' = 'A', many = false) {
  const ruleNumber = deviceId === 'B' ? 1 : many || index <= 6 ? index : index - 6;
  const address = deviceId === 'B' ? 41001 : 40001 + (ruleNumber - 1) * 10 + (!many && index > 6 ? 1 : 0);
  return { rule_id: `${deviceId}${ruleNumber}`, rule_revision: 'R2', point_id: `P${deviceId}${index}`,
    mapping_id: `M${deviceId}${index}`, tag_id: `T${deviceId}${index}`, address: String(address) };
}

import type { EvidenceScope, RuleEvidenceView } from '../../../src/types/runtimeRuleEvidence';
import type { RuntimeSourcePage, RuntimeSourcePageRequest, RuntimeSourceSelection } from '../../../src/types/studioV2RuntimeContext';
import { normalizeRuntimeSourcePageRequest } from '../../../src/services/runtimeSourcePageParser';
import { RuntimeRecordingFixtureError, runtimeRecordingEvidence, type RuntimeRecordingScenario } from './runtime-recording-state';

export function runtimeRecordingRows(scenario: RuntimeRecordingScenario, deviceId: 'A' | 'B') {
  const many = scenario === 'many'; const ruleCount = deviceId === 'B' ? 1 : many ? 120 : 6;
  const mappingCount = deviceId === 'B' ? 1 : many ? 120 : 10;
  const rules: RuntimeSourcePage['rules'] = Array.from({ length: ruleCount }, (_, i) => ({
    id: `${deviceId}${i + 1}`, device_id: deviceId, revision_id: 'R2', start_address: String(deviceId === 'B' ? 41001 : 40001 + i * 10),
    count: !many && deviceId === 'A' && i < 4 ? 2 : 1, data_type: 'float32', naming_prefix: `Rule ${deviceId}${i + 1}`, enabled: true, target_data_type: 'float64',
  }));
  const mappings: RuntimeSourcePage['mappings'] = Array.from({ length: mappingCount }, (_, i) => {
    const identity = runtimeRecordingMappingIdentity(i + 1, deviceId, many);
    return { id: identity.mapping_id, rule_id: identity.rule_id, rule_revision: identity.rule_revision, device_id: deviceId,
      point_id: identity.point_id, address: identity.address, tag_id: identity.tag_id, tag_key: `line.${deviceId}.${i + 1}`,
      display_name: `Source ${deviceId}${i + 1}`, unit: 'm/s', target_type: 'float64', enabled: true, status: 'active' };
  });
  return { rules, mappings };
}

function cursorData(scope: EvidenceScope, request: RuntimeSourcePageRequest, offset: number) {
  return { ...scope, kind: request.kind, query: request.query, limit: request.limit, rule_id: request.rule_id ?? null, rule_revision: request.rule_revision ?? null, offset };
}
function encode(value: unknown) { return btoa(unescape(encodeURIComponent(JSON.stringify(value)))).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/g, ''); }
/** Bounded cursor fixture tied to exact scope/query/revision; no hidden detail fan-out. */
export function runtimeRecordingSourcePage(scenario: RuntimeRecordingScenario, scope: EvidenceScope, request: RuntimeSourcePageRequest): RuntimeSourcePage {
  const selected = normalizeRuntimeSourcePageRequest(request);
  if (!selected) throw new RuntimeRecordingFixtureError(400);
  const rows = runtimeRecordingRows(scenario, scope.device_id as 'A' | 'B');
  if (selected.rule_id && !rows.rules.some(rule => rule.id === selected.rule_id && rule.revision_id === selected.rule_revision)) throw new RuntimeRecordingFixtureError(409, 'stale-reference');
  const query = selected.query ?? ''; const limit = selected.limit ?? 50; let offset = 0;
  if (selected.cursor) {
    try {
      const decoded: unknown = JSON.parse(decodeURIComponent(escape(atob(selected.cursor.replace(/-/g, '+').replace(/_/g, '/')))));
      if (!decoded || typeof decoded !== 'object' || !('offset' in decoded) || typeof decoded.offset !== 'number' || !Number.isSafeInteger(decoded.offset) || decoded.offset < 0 || encode(decoded) !== selected.cursor ||
        JSON.stringify(decoded) !== JSON.stringify(cursorData(scope, selected, decoded.offset))) throw new Error('Wrong cursor');
      offset = decoded.offset;
    } catch { throw new RuntimeRecordingFixtureError(409, 'stale-reference'); }
  }
  const matches = (id: string, label: string, address: string) => !query || [id, label, address].some(text => text.toLowerCase().includes(query.toLowerCase()));
  const rules = rows.rules.filter(rule => matches(rule.id, rule.naming_prefix, rule.start_address));
  const mappings = rows.mappings.filter(row => (!selected.rule_id || row.rule_id === selected.rule_id) && matches(row.id ?? '', row.display_name ?? '', row.address));
  const count = selected.kind === 'rules' ? rules.length : mappings.length;
  if (offset > count) throw new RuntimeRecordingFixtureError(409, 'stale-reference');
  return { contract_version: 1, ...scope, kind: selected.kind, query, limit, list_revision: 'L1',
    rules: selected.kind === 'rules' ? rules.slice(offset, offset + limit) : [], mappings: selected.kind === 'mappings' ? mappings.slice(offset, offset + limit) : [],
    next_cursor: offset + limit < count ? encode(cursorData(scope, selected, offset + limit)) : null };
}

export function runtimeRecordingSourceEvidence(scenario: RuntimeRecordingScenario, scope: EvidenceScope, selection: RuntimeSourceSelection): RuleEvidenceView {
  const rows = runtimeRecordingRows(scenario, scope.device_id as 'A' | 'B');
  const row = rows.mappings.find(mapping => mapping.rule_id === selection.rule_id && mapping.rule_revision === selection.rule_revision && mapping.point_id === selection.point_id && (selection.mapping_id === undefined || mapping.id === selection.mapping_id));
  if (!row) throw new RuntimeRecordingFixtureError(409, 'stale-reference');
  const mapped = selection.mapping_id !== undefined;
  const view = runtimeRecordingEvidence(scenario, scope.device_id as 'A' | 'B');
  view.source = { rule_id: row.rule_id, point_id: row.point_id, mapping_id: mapped ? row.id ?? null : null, tag_id: mapped ? row.tag_id ?? null : null,
    saved_revision: 'R2', effective_revision: null, state: 'saved-only', proof: 'unavailable', effective: null,
    saved: { binding: { workspace_id: scope.workspace_id, workspace_revision: scope.workspace_revision, entity_revision: 'R2' }, mapping_revision: mapped ? 'MR2' : null,
      enabled: true, address: row.address, protocol_function: '03', read_type: 'float32', result_type: mapped ? 'float64' : null,
      polling: { group_id: `PG${scope.device_id}`, interval_ms: 1000, source: 'polling-group' },
      transforms: [{ kind: 'byte-order', order: 'CDAB', multiplier: null, offset: null, target_type: null },
        ...(mapped ? [{ kind: 'scale' as const, order: null, multiplier: 0.1, offset: 0, target_type: null }] : [])] } };
  return view;
}
