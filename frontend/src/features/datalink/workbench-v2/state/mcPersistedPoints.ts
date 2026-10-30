import type {
  SourceRuleCandidateSnapshotView,
  SourceRuleTagCandidateView,
} from '../../../../types/sourceRuleCandidates';
import type {
  SourceRuleRecord,
  StudioV2WorkspaceDatabaseTargetRecord,
  StudioV2WorkspaceMappingRecord,
} from '../../../../types/datalink';
import type { DbTarget, Device, Mapping, Point, Rule, WorkbenchV2State } from './types';
import { workbenchV2Reducer } from './useWorkbenchV2State';
import { hydrateStudioV2DatabaseTarget } from './studioV2DatabaseAutosave';
import { hydrateStudioV2Mapping } from './studioV2MappingAutosave';
import { addressParser } from '../../../../utils/addressParser';
import { dataTypeWidth, projectPersistedMCPoints } from './sourceRule';

export type PersistedMCPointsStatus = 'checking' | 'compatible' | 'conflict' | 'unavailable';

export interface PersistedMCPointsScope {
  workspace_id: string;
  workspace_revision: string;
  revision_id: string;
}

export interface PersistedMCPointIdentity {
  point_id: string;
  address: string;
  tag_id?: string;
  mapping_id?: string;
}

export interface PersistedMCRuleState {
  status: PersistedMCPointsStatus;
  scope?: PersistedMCPointsScope;
  points: PersistedMCPointIdentity[];
  reason?: string;
}

export function persistedMCScopeKey(evidence: PersistedMCRuleState | undefined): string {
  const scope = evidence?.scope;
  return scope ? `${scope.workspace_id}:${scope.workspace_revision}:${scope.revision_id}` : '';
}

const MC_PROTOCOL = 'mc_3e';

function normalizeAddress(address: string): string {
  return address.trim().toUpperCase();
}

function isReadyStatus(status: PersistedMCPointsStatus): boolean {
  return status === 'compatible' || status === 'conflict';
}

function unavailable(scope: PersistedMCPointsScope | undefined, reason: string): PersistedMCRuleState {
  return { status: 'unavailable', ...(scope ? { scope } : {}), points: [], reason };
}

function expectedAddresses(rule: Pick<SourceRuleRecord, 'start_address' | 'count' | 'data_type' | 'skipped_addresses'>): string[] | null {
  if (!Number.isInteger(rule.count) || rule.count < 0) {
    return null;
  }
  if (!Array.isArray(rule.skipped_addresses) || rule.skipped_addresses.some((address) => typeof address !== 'string')) {
    return null;
  }

  const skipped = new Set(rule.skipped_addresses.map(normalizeAddress));
  const width = dataTypeWidth(rule.data_type);
  const addresses: string[] = [];
  try {
    for (let index = 0; index < rule.count; index += 1) {
      const address = addressParser.offset(rule.start_address, index * width, MC_PROTOCOL);
      if (!skipped.has(normalizeAddress(address))) {
        addresses.push(normalizeAddress(address));
      }
    }
  } catch {
    return null;
  }
  return addresses;
}

function candidateIdentity(candidate: SourceRuleTagCandidateView): PersistedMCPointIdentity | null {
  if (!candidate || typeof candidate !== 'object') {
    return null;
  }
  const pointID = typeof candidate.point_id === 'string' ? candidate.point_id.trim() : '';
  const address = typeof candidate.address === 'string' ? normalizeAddress(candidate.address) : '';
  if (!pointID || !address) {
    return null;
  }
  return {
    point_id: pointID,
    address,
    ...(typeof candidate.tag_id === 'string' && candidate.tag_id ? { tag_id: candidate.tag_id } : {}),
    ...(typeof candidate.mapping_id === 'string' && candidate.mapping_id ? { mapping_id: candidate.mapping_id } : {}),
  };
}
export function buildPersistedMCRuleState(
  rule: SourceRuleRecord,
  scope: PersistedMCPointsScope | undefined,
  snapshot?: SourceRuleCandidateSnapshotView,
): PersistedMCRuleState {
  if (!scope || !rule.revision_id) {
    return unavailable(scope, 'scope_unavailable');
  }
  if (!snapshot || snapshot.source_rule_id !== rule.id || snapshot.workspace_id !== scope.workspace_id ||
    snapshot.workspace_revision !== scope.workspace_revision || snapshot.revision_id !== scope.revision_id) {
    return unavailable(scope, 'candidate_scope_mismatch');
  }
  if (!snapshot.tags || typeof snapshot.tags !== 'object' || snapshot.tags.status !== 'ready' || !Array.isArray(snapshot.tags.candidates)) {
    return unavailable(scope, 'candidate_set_incomplete');
  }

  const points: PersistedMCPointIdentity[] = [];
  const addresses = new Set<string>();
  const pointIDs = new Set<string>();
  for (const candidate of snapshot.tags.candidates) {
    const identity = candidateIdentity(candidate);
    if (!identity || addresses.has(identity.address) || pointIDs.has(identity.point_id)) {
      return unavailable(scope, 'candidate_set_incomplete');
    }
    addresses.add(identity.address);
    pointIDs.add(identity.point_id);
    points.push(identity);
  }

  const planned = expectedAddresses(rule);
  if (!planned) {
    return {
      status: 'conflict',
      scope,
      points,
      reason: 'stored_rule_invalid_under_current_planner',
    };
  }
  if (planned.length !== points.length) {
    return unavailable(scope, 'candidate_set_incomplete');
  }
  const plannedSet = new Set(planned);
  const compatible = plannedSet.size === addresses.size &&
    points.every((point) => plannedSet.has(point.address));
  return {
    status: compatible ? 'compatible' : 'conflict',
    scope,
    points,
    ...(compatible ? {} : { reason: 'persisted_address_set_conflict' }),
  };
}
export function isPersistedMCRule(rule: Pick<Rule, 'persisted' | 'device_id'>, devices: readonly Device[]): boolean {
  return rule.persisted === true && devices.some((device) => device.id === rule.device_id && device.protocol === MC_PROTOCOL);
}

export function withPersistedMCEvidence(
  rules: readonly Rule[],
  evidenceByRule: Readonly<Record<string, PersistedMCRuleState>>,
): Rule[] {
  return rules.map((rule) => {
    const evidence = evidenceByRule[rule.id];
    return evidence ? { ...rule, mc_persisted_points: evidence } : rule;
  });
}
export function markPersistedMCChecking(evidence: PersistedMCRuleState | undefined, revisionID?: string): PersistedMCRuleState | undefined {
  if (!evidence) return undefined;
  return {
    ...evidence,
    status: 'checking',
    reason: 'rule_revision_pending',
    scope: evidence.scope && revisionID ? { ...evidence.scope, revision_id: revisionID } : evidence.scope,
    points: [...evidence.points],
  };
}

export function mergePersistedMCPoints(
  currentPoints: readonly Point[],
  rules: readonly Rule[],
  evidenceByRule: Readonly<Record<string, PersistedMCRuleState>>,
  previousEvidenceByRule: Readonly<Record<string, PersistedMCRuleState | undefined>> = evidenceByRule,
): Point[] {
  const affectedRuleIDs = new Set(
    rules.filter((rule) => rule.persisted === true && rule.mc_persisted_points).map((rule) => rule.id),
  );
  const replacementByRule = new Map<string, Point[]>();
  for (const rule of rules) {
    const evidence = evidenceByRule[rule.id];
    if (rule.persisted === true && evidence && isReadyStatus(evidence.status)) {
      replacementByRule.set(rule.id, projectPersistedMCPoints({ ...rule, mc_persisted_points: evidence }));
    } else if (rule.persisted === true) {
      const previous = previousEvidenceByRule[rule.id];
      if (previous?.points.length) {
        replacementByRule.set(rule.id, projectPersistedMCPoints({ ...rule, mc_persisted_points: previous }));
      }
    }
  }

  const merged: Point[] = [];
  const inserted = new Set<string>();
  for (const point of currentPoints) {
    if (!affectedRuleIDs.has(point.rule_id)) {
      merged.push(point);
      continue;
    }
    if (!inserted.has(point.rule_id)) {
      merged.push(...(replacementByRule.get(point.rule_id) ?? []));
      inserted.add(point.rule_id);
    }
  }
  for (const [ruleID, points] of replacementByRule) {
    if (!inserted.has(ruleID)) {
      merged.push(...points);
    }
  }
  return merged;
}
function mappingKey(ruleID: string, address: string): string {
  return `${ruleID}::${normalizeAddress(address)}`;
}
function hasPendingSave(saveState?: string): boolean {
  return saveState === 'saving' || saveState === 'save-error' || saveState === 'draft-invalid';
}

function sameMappingValue(left?: Mapping['local_value'], right?: Mapping['persisted_value']): boolean {
  return Boolean(
    left && right &&
    left.tag_key === right.tag_key &&
    left.display_name === right.display_name &&
    left.unit === right.unit &&
    left.target_type === right.target_type &&
    left.scale === right.scale &&
    left.offset === right.offset &&
    left.enabled === right.enabled,
  );
}

function hasLocalMappingDraft(mapping: Mapping | undefined): boolean {
  if (!mapping) return false;
  if (mapping.persisted !== true || hasPendingSave(mapping.save_state)) return true;
  return Boolean(mapping.local_value && mapping.persisted_value && !sameMappingValue(mapping.local_value, mapping.persisted_value));
}

function hasLocalDatabaseTargetDraft(target: DbTarget | undefined): boolean {
  return Boolean(target && (target.persisted !== true || hasPendingSave(target.save_state)));
}

function mergeLocalDraftPoints(
  currentPoints: readonly Point[],
  persistedPoints: readonly Point[],
  affectedRuleIDs: ReadonlySet<string>,
  mappings: Readonly<Record<string, Mapping>>,
  targets: Readonly<Record<string, DbTarget>>,
): Point[] {
  const localDraftPointIDs = new Set<string>();
  for (const mapping of Object.values(mappings)) {
    if (mapping.rule_id && affectedRuleIDs.has(mapping.rule_id) && hasLocalMappingDraft(mapping)) {
      localDraftPointIDs.add(mapping.point_id);
    }
  }
  for (const [pointID, target] of Object.entries(targets)) {
    const point = currentPoints.find((candidate) => candidate.id === pointID);
    if (point && affectedRuleIDs.has(point.rule_id) && hasLocalDatabaseTargetDraft(target)) {
      localDraftPointIDs.add(pointID);
    }
  }

  if (localDraftPointIDs.size === 0) return [...persistedPoints];

  const merged = [...persistedPoints];
  for (const point of currentPoints) {
    if (!localDraftPointIDs.has(point.id)) continue;
    const byID = merged.findIndex((candidate) => candidate.id === point.id);
    if (byID >= 0) {
      merged[byID] = point;
      continue;
    }
    const byAddress = merged.findIndex((candidate) => candidate.rule_id === point.rule_id && normalizeAddress(candidate.address) === normalizeAddress(point.address));
    if (byAddress >= 0) {
      merged[byAddress] = point;
      continue;
    }
    merged.push(point);
  }
  return merged;
}

export function mergePersistedMCState(
  state: WorkbenchV2State,
  evidenceByRule: Readonly<Record<string, PersistedMCRuleState>>,
  mappingRecords: readonly StudioV2WorkspaceMappingRecord[] | undefined,
  targetRecords: readonly StudioV2WorkspaceDatabaseTargetRecord[] | undefined,
): Pick<WorkbenchV2State, 'rules' | 'points' | 'mappings' | 'db'> {
  const previousEvidenceByRule = Object.fromEntries(state.rules.map((rule) => [rule.id, rule.mc_persisted_points]));
  const effectiveEvidenceByRule: Record<string, PersistedMCRuleState> = { ...evidenceByRule };
  for (const [ruleID, evidence] of Object.entries(evidenceByRule)) {
    const previous = previousEvidenceByRule[ruleID];
    if (previous?.points.length && evidence.points.length === 0 && (evidence.status === 'checking' || evidence.status === 'unavailable')) {
      effectiveEvidenceByRule[ruleID] = { ...evidence, points: previous.points.map((point) => ({ ...point })) };
    }
  }
  const rules = withPersistedMCEvidence(state.rules, effectiveEvidenceByRule).map((rule) =>
    effectiveEvidenceByRule[rule.id] || !rule.mc_persisted_points
      ? rule
      : { ...rule, mc_persisted_points: undefined },
  );
  const persistedPoints = mergePersistedMCPoints(state.points, rules, effectiveEvidenceByRule, previousEvidenceByRule);
  const affectedRuleIDs = new Set(Object.keys(effectiveEvidenceByRule));
  const points = mergeLocalDraftPoints(state.points, persistedPoints, affectedRuleIDs, state.mappings, state.db.targets);
  const currentByAddress = new Map<string, Mapping>();
  for (const mapping of Object.values(state.mappings)) {
    if (mapping.rule_id && affectedRuleIDs.has(mapping.rule_id)) {
      currentByAddress.set(mappingKey(mapping.rule_id, mapping.address ?? ''), mapping);
    }
  }

  const filteredMappings = Object.fromEntries(
    Object.entries(state.mappings).filter(([, mapping]) => !mapping.rule_id || !affectedRuleIDs.has(mapping.rule_id)),
  );
  const mappings = workbenchV2Reducer(
    { ...state, points, mappings: filteredMappings },
    { type: 'initMappingsForPoints', points },
  ).mappings;
  const recordByAddress = new Map((mappingRecords ?? []).map((record) => [mappingKey(record.rule_id, record.address), record]));
  for (const point of points) {
    const record = recordByAddress.get(mappingKey(point.rule_id, point.address));
    if (!record) continue;
    const current = affectedRuleIDs.has(point.rule_id)
      ? currentByAddress.get(mappingKey(point.rule_id, point.address))
      : mappings[point.id];
    mappings[point.id] = hydrateStudioV2Mapping(point, record, current);
  }
  for (const point of points) {
    if (!affectedRuleIDs.has(point.rule_id) || recordByAddress.has(mappingKey(point.rule_id, point.address))) continue;
    const current = currentByAddress.get(mappingKey(point.rule_id, point.address));
    if (hasLocalMappingDraft(current)) mappings[point.id] = current!;
  }

  const affectedPointIDs = new Set([
    ...state.points.filter((point) => affectedRuleIDs.has(point.rule_id)).map((point) => point.id),
    ...points.filter((point) => affectedRuleIDs.has(point.rule_id)).map((point) => point.id),
  ]);
  const targets = { ...state.db.targets };
  for (const pointID of affectedPointIDs) delete targets[pointID];
  const aliases = new Map<string, string>();
  for (const mapping of Object.values(mappings)) {
    aliases.set(mapping.point_id, mapping.point_id);
    if (mapping.persisted_point_id) aliases.set(mapping.persisted_point_id, mapping.point_id);
  }
  const hydratedTargetPointIDs = new Set<string>();
  for (const record of targetRecords ?? []) {
    const pointID = aliases.get(record.point_id);
    if (!pointID) continue;
    hydratedTargetPointIDs.add(pointID);
    const current = state.db.targets[pointID];
    targets[pointID] = current && hasPendingSave(current.save_state)
      ? current
      : hydrateStudioV2DatabaseTarget({ ...record, point_id: pointID }, current);
  }
  for (const point of points) {
    if (!affectedRuleIDs.has(point.rule_id) || hydratedTargetPointIDs.has(point.id)) continue;
    const current = state.db.targets[point.id];
    if (hasLocalDatabaseTargetDraft(current)) targets[point.id] = current!;
  }

  return { rules, points, mappings, db: { ...state.db, targets } };
}
