import type { Device, SourceRuleRecord, StudioV2WorkspaceMappingRecord } from '../../../../types/datalink';
import type { StudioV2Workspace } from '../../../../types/studioV2Workspace';
import type { WriteGroup, WriteGroupListResponse, WriteGroupMember } from '../../../../types/studioV2WriteGroup';
import type { WorkbenchV2State } from './types';

/**
 * URL-independent viewing intent. IDs are canonical identities, never names.
 * Step chooses the visible panel (1–4); validated context IDs may accompany any
 * step so inspecting another panel does not discard the selected group. Rule +
 * group means independent same-device viewing context, NOT rule membership in
 * that group: WriteGroup members have no rule_id. Only an explicitly selected
 * mapping proves the device/point/tag membership link. None authorizes editing.
 */
export interface StudioSelectionRequest {
  readonly workspaceId: StudioV2Workspace['id'];
  readonly step: WorkbenchV2State['current'];
  readonly deviceId?: Device['id'];
  readonly ruleId?: SourceRuleRecord['id'];
  readonly mappingId?: StudioV2WorkspaceMappingRecord['id'];
  readonly groupId?: WriteGroup['id'];
  readonly groupRevisionHint?: WriteGroup['revision'];
}

type RuleIdentity = Readonly<Pick<SourceRuleRecord, 'id' | 'workspace_id' | 'device_id' | 'revision_id'>>;
type MappingIdentity = Readonly<Pick<StudioV2WorkspaceMappingRecord, 'id' | 'workspace_id' | 'device_id' | 'rule_id' | 'point_id' | 'tag_id'>>;
type GroupIdentity = Readonly<Pick<WriteGroup, 'id' | 'workspace_id' | 'revision' | 'status' | 'basic_managed_device_id'>> & {
  readonly members: readonly Readonly<Pick<WriteGroupMember, 'device_id' | 'point_id' | 'tag_id'>>[];
};

/**
 * Read-only provider projection, NOT a new backend DTO or revision authority.
 * `current` means ownership reads were revalidated, not merely cached or hydrated.
 * Workspace order proves device ownership; device records prove current existence.
 * Rule/group currentRevision comes from backend head revalidation, independently
 * of the returned record. Group workspaceRevision is the canonical list envelope;
 * workspaceRevision is the current database setup revision, NOT a device/mapping
 * revision. Mapping currentness must be revalidated by the provider (no revision
 * field exists on its canonical DTO). Missing/stale proof must not become current.
 */
export type StudioSelectionEvidence = {
  readonly state: 'current';
  readonly workspace: Readonly<Pick<StudioV2Workspace, 'id'>> & { readonly ordered_device_ids: readonly string[] };
  /** Backend newRecord starts at ""; that explicit current value is not missing proof. */
  readonly workspaceRevision: WriteGroupListResponse['workspace_revision'];
  readonly devices?: readonly Readonly<Pick<Device, 'id'>>[];
  readonly rules?: readonly { readonly value: RuleIdentity; readonly currentRevision: NonNullable<SourceRuleRecord['revision_id']> }[];
  readonly mappings?: readonly MappingIdentity[];
  readonly groups?: readonly {
    readonly value: GroupIdentity;
    readonly currentRevision: WriteGroup['revision'];
    readonly workspaceRevision: WriteGroupListResponse['workspace_revision'];
  }[];
} | { readonly state: 'loading' | 'error' | 'stale' | 'unavailable' };

/** Only read access; mounting an actual query/route receiver belongs to integration. */
export interface StudioSelectionEvidenceProvider {
  readonly readCurrent: (request: StudioSelectionRequest) => StudioSelectionEvidence | undefined;
}

export type StudioSelectionUnavailableReason = 'invalid_request' | 'scope_unavailable' |
  'evidence_unavailable' | 'evidence_loading' | 'evidence_error' | 'evidence_stale' |
  'target_unavailable' | 'revision_unavailable' | 'revision_conflict';

/** Selecting a saved target is not mutation permission; all later actions revalidate their guards. */
export type StudioSelectionResult = {
  readonly state: 'selected';
  readonly target: Omit<StudioSelectionRequest, 'groupRevisionHint'> & {
    readonly workspaceRevision: string;
    readonly ruleRevision?: string;
    readonly groupRevision?: string;
  };
  /** No historical/applied snapshot is supplied; a differing hint does not prove it exists. */
  readonly history?: { readonly groupId: string; readonly revisionHint: string; readonly relation: 'current' | 'different' };
} | { readonly state: 'unavailable'; readonly reason: StudioSelectionUnavailableReason };

function unavailable(reason: StudioSelectionUnavailableReason): StudioSelectionResult {
  return Object.freeze({ state: 'unavailable', reason });
}

function identity(value: unknown): value is string {
  return typeof value === 'string' && value.length > 0 && value.trim() === value;
}

function unique<T>(values: readonly T[] | undefined, matches: (value: T) => boolean): T | undefined {
  const found = values?.filter(matches);
  return found?.length === 1 ? found[0] : undefined;
}

/**
 * Revalidate a viewing request using only supplied current backend evidence.
 * No cache, fallback selection, session writes, URL parsing or mutation calls.
 * Results identify current saved records, not effective/running records. A device
 * remains viewing context: selecting a shared group does not make it device-owned.
 * Callers must retain drafts separately and revalidate CAS/readiness/lifecycle on
 * every later explicit action; neither this result nor its hint authorizes one.
 */
export function resolveStudioSelection(
  currentWorkspaceId: string,
  request: StudioSelectionRequest,
  provider?: StudioSelectionEvidenceProvider,
): StudioSelectionResult {
  const optional = ['deviceId', 'ruleId', 'mappingId', 'groupId', 'groupRevisionHint'] as const;
  if (!request || !identity(request.workspaceId) || ![1, 2, 3, 4].includes(request.step) ||
      optional.some(key => request[key] !== undefined && !identity(request[key])) ||
      (request.groupRevisionHint !== undefined && request.groupId === undefined)) return unavailable('invalid_request');
  if (!identity(currentWorkspaceId) || request.workspaceId !== currentWorkspaceId) return unavailable('scope_unavailable');
  if (!provider) return unavailable('evidence_unavailable');

  // Whitelist and freeze before invoking caller code; never forward extra fields.
  const { workspaceId, step, deviceId, ruleId, mappingId, groupId, groupRevisionHint } = request;
  const intent = Object.freeze({ workspaceId, step,
    ...(deviceId !== undefined ? { deviceId } : {}), ...(ruleId !== undefined ? { ruleId } : {}),
    ...(mappingId !== undefined ? { mappingId } : {}), ...(groupId !== undefined ? { groupId } : {}),
    ...(groupRevisionHint !== undefined ? { groupRevisionHint } : {}) });
  let evidence: StudioSelectionEvidence | undefined;
  try {
    evidence = provider.readCurrent(intent);
  } catch {
    return unavailable('evidence_error');
  }
  if (!evidence) return unavailable('evidence_unavailable');
  switch (evidence.state) {
    case 'current': break;
    case 'loading': case 'error': case 'stale': case 'unavailable': return unavailable(`evidence_${evidence.state}`);
    default: return unavailable('evidence_unavailable');
  }
  if (evidence.workspace.id !== workspaceId) return unavailable('scope_unavailable');
  // service_database_setup.go:newRecord leaves DatabaseSetupRevision empty until
  // the first setup change; missing (undefined) remains different from explicit "".
  if (typeof evidence.workspaceRevision !== 'string') return unavailable('revision_unavailable');
  const ownsDevice = (id: string) => identity(id) &&
    evidence.workspace.ordered_device_ids.filter(owned => owned === id).length === 1 &&
    !!unique(evidence.devices, device => device.id === id);

  const mapping = mappingId === undefined ? undefined : unique(evidence.mappings, item => item.id === mappingId);
  if (mappingId !== undefined && (!mapping || mapping.workspace_id !== workspaceId ||
      !identity(mapping.rule_id) || !identity(mapping.device_id) || !identity(mapping.point_id) || !identity(mapping.tag_id) ||
      (ruleId !== undefined && mapping.rule_id !== ruleId) || (deviceId !== undefined && mapping.device_id !== deviceId))) {
    return unavailable('target_unavailable');
  }
  const selectedRuleId = ruleId ?? mapping?.rule_id;
  const rule = selectedRuleId === undefined ? undefined : unique(evidence.rules, item => item.value.id === selectedRuleId);
  if (selectedRuleId !== undefined && (!rule || rule.value.workspace_id !== workspaceId ||
      !identity(rule.value.device_id) || (deviceId !== undefined && rule.value.device_id !== deviceId) ||
      (mapping !== undefined && rule.value.device_id !== mapping.device_id))) return unavailable('target_unavailable');
  if (rule && (!identity(rule.value.revision_id) || !identity(rule.currentRevision))) return unavailable('revision_unavailable');
  if (rule && rule.value.revision_id !== rule.currentRevision) return unavailable('revision_conflict');
  const selectedDeviceId = deviceId ?? rule?.value.device_id;
  if (selectedDeviceId !== undefined && !ownsDevice(selectedDeviceId)) return unavailable('target_unavailable');

  const group = groupId === undefined ? undefined : unique(evidence.groups, item => item.value.id === groupId);
  if (groupId !== undefined && (!group || group.value.workspace_id !== workspaceId ||
      !['draft', 'ready', 'running', 'disabled'].includes(group.value.status))) {
    return unavailable('target_unavailable');
  }
  if (group) {
    if (!identity(group.value.revision) || !identity(group.currentRevision)) return unavailable('revision_unavailable');
    if (group.value.revision !== group.currentRevision || group.workspaceRevision !== evidence.workspaceRevision) {
      return unavailable('revision_conflict');
    }
    const owner = group.value.basic_managed_device_id;
    const members = group.value.members;
    // Unselected orphan member references must not prevent opening a canonical
    // group for repair. Only the explicitly selected chain needs live ownership.
    if ((owner !== undefined && selectedDeviceId !== undefined && selectedDeviceId !== owner) ||
        (selectedDeviceId !== undefined && owner !== selectedDeviceId && !members.some(member => member.device_id === selectedDeviceId)) ||
        (mapping !== undefined && !members.some(member => member.device_id === mapping.device_id &&
          member.point_id === mapping.point_id && member.tag_id === mapping.tag_id))) return unavailable('target_unavailable');
  }

  const target = Object.freeze({ workspaceId, step, workspaceRevision: evidence.workspaceRevision,
    ...(selectedDeviceId !== undefined ? { deviceId: selectedDeviceId } : {}),
    ...(rule ? { ruleId: rule.value.id, ruleRevision: rule.value.revision_id } : {}),
    ...(mapping ? { mappingId: mapping.id } : {}),
    ...(group ? { groupId: group.value.id, groupRevision: group.value.revision } : {}) });
  return Object.freeze({ state: 'selected', target,
    ...(group && groupRevisionHint !== undefined ? { history: Object.freeze({ groupId: group.value.id,
      revisionHint: groupRevisionHint, relation: groupRevisionHint === group.value.revision ? 'current' as const : 'different' as const }) } : {}) });
}
