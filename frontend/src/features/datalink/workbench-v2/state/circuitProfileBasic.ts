import type { MeasurementDefinition } from '../../../../types/measurement';
import type { RecordingStartOperation, RecordingStartRequest } from '../../../../types/studioV2RecordingStart';
import type {
  WriteGroup,
  WriteGroupDestinationDraft,
  WriteGroupDraft,
  WriteGroupLayoutVersion,
  WriteGroupMemberDraft,
  WriteGroupProfileKind,
  WriteGroupSampleReduction,
  WriteGroupStorageMode,
  WriteGroupValueRole,
} from '../../../../types/studioV2WriteGroup';
import { isValidRecordingProfile, recordingProfileIntervalSeconds } from '../../../../utils/studioV2RecordingProfile';
import { createBasicRecordingRequestId, recordingStartRevisionsMatch } from './basicRecording';
import type { GroupCandidate } from './writeGroup/candidates';
import { isConfirmedCounterDefinition } from './circuitCounterConfirmation';

/**
 * What the recording is for. A whole-device snapshot is the wide profile layout (one row per time window for the
 * device or circuit); selected measurements are the long layout (one row per source per window). The purpose only
 * picks the table layout: either one can record instantaneous and cumulative profiles.
 */
export type CircuitPurpose = 'wide' | 'long';

const PURPOSE_LAYOUT: Record<CircuitPurpose, WriteGroupLayoutVersion> = { wide: 'managed_wide_profile_v1', long: 'managed_long_v1' };

/** The purpose of a saved layout; any other layout (legacy wide, unknown) belongs to none of the two new paths. */
export function purposeOfLayout(layout: string | undefined): CircuitPurpose | undefined {
  if (layout === PURPOSE_LAYOUT.wide) return 'wide';
  if (layout === PURPOSE_LAYOUT.long) return 'long';
  return undefined;
}

const NUMERIC_TARGET_TYPES = ['int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64'];

/**
 * The role a source shows before the operator chooses: a number is averaged, anything else keeps its latest value.
 * It is never a counter: a counter is only what the operator selects and a saved definition confirms, whatever the
 * name or unit of the source says.
 */
export function defaultMemberRole(candidate: Pick<GroupCandidate, 'target_type'>): WriteGroupValueRole {
  return NUMERIC_TARGET_TYPES.includes(candidate.target_type) ? 'gauge' : 'state';
}

export function candidateKey(candidate: Pick<GroupCandidate, 'device_id' | 'point_id' | 'tag_id'>): string {
  return `${candidate.device_id}\0${candidate.point_id}\0${candidate.tag_id}`;
}

export interface CircuitProfileLocalDraft {
  circuit_id: string;
  circuit_name: string;
  selected_keys: string[];
  instantaneous_keys: string[];
  cumulative_keys: string[];
  /** The table layout chosen for a new circuit; a saved circuit keeps the layout it was created with. */
  purpose: CircuitPurpose;
  /** Roles the operator picked for instantaneous sources, by source key; an unpicked source shows its default role. */
  member_roles: Record<string, WriteGroupValueRole>;
  instantaneous_enabled: boolean;
  instantaneous_interval_input: string;
  cumulative_enabled: boolean;
  cumulative_interval_input: string;
  storage_mode: WriteGroupStorageMode;
  recreate_missing: boolean;
}

export interface ParsedProfileInterval {
  valid: boolean;
  intervalValue?: number;
  intervalSeconds?: number;
  errorCode?: 'profile_interval_invalid';
}

export interface CandidateCounterStatus {
  candidate: GroupCandidate;
  confirmed: boolean;
  measurement?: MeasurementDefinition;
  issueCode?:
    | 'counter-semantics-required'
    | 'counter-measurement-not-counter'
    | 'counter-measurement-foreign'
    | 'counter-measurement-unconfirmed'
    | 'counter-type-unsupported';
}

function draftStorageKey(workspaceId: string, deviceId: string, circuitId?: string): string {
  const scope = circuitId !== undefined && circuitId.trim()
    ? `${workspaceId.trim()}/${deviceId.trim()}/${circuitId.trim()}`
    : `${workspaceId.trim()}/${deviceId.trim()}`;
  return `wbv2.circuit-profile-draft.${encodeURIComponent(scope)}`;
}

function sessionStore(): Storage | undefined {
  if (typeof window === 'undefined') return undefined;
  try { return window.sessionStorage; } catch { return undefined; }
}

export function deriveCircuitId(circuitName: string, existingCircuitId = ''): string {
  if (existingCircuitId.trim()) return existingCircuitId.trim();
  return circuitName.trim() ? `cir-${createBasicRecordingRequestId()}` : '';
}

export function groupProfileKind(group: WriteGroup | WriteGroupDraft | undefined): WriteGroupProfileKind | undefined {
  return group?.recording_profile?.kind;
}

export function groupCircuitId(group: WriteGroup | WriteGroupDraft | undefined): string {
  return (group?.members[0]?.circuit_id ?? '').trim();
}

export function groupCircuitName(group: WriteGroup | WriteGroupDraft | undefined): string {
  return (group?.members[0]?.circuit_name ?? '').trim();
}

export function defaultCircuitProfileLocalDraft(
  candidates: GroupCandidate[],
  savedLongGroups: WriteGroup[],
  targetCircuitId?: string, explicit = false,
): CircuitProfileLocalDraft {
  const selectedCircuitId = targetCircuitId !== undefined ? targetCircuitId : groupCircuitId(savedLongGroups[0]);
  const circuitGroups = savedLongGroups.filter((group) => groupCircuitId(group) === selectedCircuitId);
  const blank = explicit && circuitGroups.length === 0; // migration prepare: nothing is chosen for the operator
  const instGroup = circuitGroups.find((group) => groupProfileKind(group) === 'instantaneous');
  const cumGroup = circuitGroups.find((group) => groupProfileKind(group) === 'cumulative');
  const anyGroup = cumGroup ?? instGroup;
  const allCandidateKeys = candidates.map((candidate) => candidateKey(candidate));
  const instSaved = new Set((instGroup?.members ?? []).map((m) => `${m.device_id}\0${m.point_id}\0${m.tag_id}`));
  const cumSaved = new Set((cumGroup?.members ?? []).map((m) => `${m.device_id}\0${m.point_id}\0${m.tag_id}`));
  const instKeys = instGroup ? allCandidateKeys.filter((k) => instSaved.has(k)) : (savedLongGroups.length === 0 && !targetCircuitId ? allCandidateKeys : []);
  const cumKeys = cumGroup ? allCandidateKeys.filter((k) => cumSaved.has(k)) : (savedLongGroups.length === 0 && !targetCircuitId ? allCandidateKeys : []);
  const instSeconds = instGroup?.recording_profile?.interval ?? instGroup?.row_policy.interval_seconds ?? 15;
  const cumMinutes = cumGroup?.recording_profile?.interval ?? (cumGroup ? Math.max(1, Math.round(cumGroup.row_policy.interval_seconds / 60)) : 10);
  const circuitName = groupCircuitName(anyGroup);
  const circuitId = targetCircuitId !== undefined ? targetCircuitId : (groupCircuitId(anyGroup) || deriveCircuitId(circuitName));
  const savedRoles: Record<string, WriteGroupValueRole> = {};
  for (const member of instGroup?.members ?? []) {
    if (member.value_role) savedRoles[`${member.device_id}\0${member.point_id}\0${member.tag_id}`] = member.value_role;
  }
  return {
    circuit_id: circuitId, circuit_name: circuitName,
    selected_keys: blank ? [] : Array.from(new Set([...instKeys, ...cumKeys])),
    instantaneous_keys: blank ? [] : instKeys, cumulative_keys: blank ? [] : cumKeys,
    purpose: purposeOfLayout(anyGroup?.layout_version) ?? 'long', member_roles: savedRoles,
    instantaneous_enabled: !blank && (circuitGroups.length > 0 ? Boolean(instGroup) : true),
    instantaneous_interval_input: blank ? '' : String(instSeconds),
    cumulative_enabled: Boolean(cumGroup), cumulative_interval_input: blank ? '' : String(cumMinutes),
    storage_mode: anyGroup?.storage_mode ?? 'plain',
    recreate_missing: anyGroup?.write_policy.recreate_missing ?? true,
  };
}

export function loadCircuitProfileLocalDraft(
  workspaceId: string, deviceId: string, candidates: GroupCandidate[], savedLongGroups: WriteGroup[], targetCircuitId?: string, explicit = false,
): CircuitProfileLocalDraft {
  const fallback = defaultCircuitProfileLocalDraft(candidates, savedLongGroups, targetCircuitId, explicit);
  if (!workspaceId.trim() || !deviceId.trim()) return fallback;
  try {
    const store = sessionStore();
    let raw = targetCircuitId ? store?.getItem(draftStorageKey(workspaceId, deviceId, targetCircuitId)) : undefined;
    if (!raw) {
      const legacyRaw = store?.getItem(draftStorageKey(workspaceId, deviceId));
      if (legacyRaw) {
        try {
          const parsed = JSON.parse(legacyRaw) as Partial<CircuitProfileLocalDraft>;
          const legacyCid = (parsed?.circuit_id ?? '').trim();
          const shouldAcceptLegacy = targetCircuitId !== undefined
            ? legacyCid === targetCircuitId.trim()
            : (savedLongGroups.length === 0 || (legacyCid !== '' && legacyCid === groupCircuitId(savedLongGroups[0])));
          if (shouldAcceptLegacy) {
            raw = legacyRaw;
          }
        } catch { /* ignore */ }
      }
    }
    if (!raw) return fallback;
    const parsed = JSON.parse(raw) as Partial<CircuitProfileLocalDraft>;
    if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) return fallback;
    const validKeys = new Set(candidates.map((candidate) => candidateKey(candidate)));
    const legacySelected = Array.isArray(parsed.selected_keys)
      ? parsed.selected_keys.filter((k): k is string => typeof k === 'string' && validKeys.has(k))
      : undefined;
    const instKeys = Array.isArray(parsed.instantaneous_keys)
      ? parsed.instantaneous_keys.filter((k): k is string => typeof k === 'string' && validKeys.has(k))
      : (legacySelected ?? fallback.instantaneous_keys);
    const cumKeys = Array.isArray(parsed.cumulative_keys)
      ? parsed.cumulative_keys.filter((k): k is string => typeof k === 'string' && validKeys.has(k))
      : (legacySelected ?? fallback.cumulative_keys);
    const circuitName = typeof parsed.circuit_name === 'string' ? parsed.circuit_name : fallback.circuit_name;
    const savedCircuit = savedLongGroups.find((group) => groupCircuitId(group) === (targetCircuitId ?? groupCircuitId(savedLongGroups[0])));
    const savedPurpose = purposeOfLayout(savedCircuit?.layout_version);
    const roles: Record<string, WriteGroupValueRole> = {};
    if (parsed.member_roles && typeof parsed.member_roles === 'object' && !Array.isArray(parsed.member_roles)) {
      for (const [key, role] of Object.entries(parsed.member_roles)) {
        if (validKeys.has(key) && (role === 'gauge' || role === 'state' || role === 'counter_raw')) roles[key] = role;
      }
    }
    const circuitId = typeof parsed.circuit_id === 'string' && parsed.circuit_id.trim()
      ? parsed.circuit_id.trim()
      : (targetCircuitId !== undefined ? targetCircuitId : deriveCircuitId(circuitName));
    return {
      circuit_id: circuitId, circuit_name: circuitName,
      selected_keys: Array.from(new Set([...instKeys, ...cumKeys])),
      instantaneous_keys: instKeys, cumulative_keys: cumKeys,
      // A saved circuit keeps its layout and saved roles; a new one restores what the operator had chosen.
      purpose: savedPurpose ?? (parsed.purpose === 'wide' ? 'wide' : 'long'),
      member_roles: savedPurpose ? { ...roles, ...fallback.member_roles } : roles,
      instantaneous_enabled: typeof parsed.instantaneous_enabled === 'boolean' ? parsed.instantaneous_enabled : fallback.instantaneous_enabled,
      instantaneous_interval_input: typeof parsed.instantaneous_interval_input === 'string' ? parsed.instantaneous_interval_input : fallback.instantaneous_interval_input,
      cumulative_enabled: typeof parsed.cumulative_enabled === 'boolean' ? parsed.cumulative_enabled : fallback.cumulative_enabled,
      cumulative_interval_input: typeof parsed.cumulative_interval_input === 'string' ? parsed.cumulative_interval_input : fallback.cumulative_interval_input,
      storage_mode: savedCircuit?.storage_mode ?? (parsed.storage_mode === 'plain' || parsed.storage_mode === 'timescale' || parsed.storage_mode === 'native_monthly'
        ? parsed.storage_mode : fallback.storage_mode),
      recreate_missing: typeof parsed.recreate_missing === 'boolean' ? parsed.recreate_missing : fallback.recreate_missing,
    };
  } catch {
    return fallback;
  }
}

export function saveCircuitProfileLocalDraft(workspaceId: string, deviceId: string, draft: CircuitProfileLocalDraft): void {
  if (!workspaceId.trim() || !deviceId.trim()) return;
  try {
    const store = sessionStore();
    const data = JSON.stringify(draft);
    if (draft.circuit_id) {
      store?.setItem(draftStorageKey(workspaceId, deviceId, draft.circuit_id), data);
    }
    store?.setItem(draftStorageKey(workspaceId, deviceId), data);
  } catch { /* ignore */ }
}

export function parseStrictProfileIntervalInput(profile: WriteGroupProfileKind, rawInput: string): ParsedProfileInterval {
  const trimmed = rawInput.trim();
  if (!/^\d+$/.test(trimmed)) return { valid: false, errorCode: 'profile_interval_invalid' };
  const intervalValue = Number(trimmed);
  const unit = profile === 'instantaneous' ? 'second' : 'minute';
  if (!isValidRecordingProfile({ kind: profile, interval: intervalValue, unit })) {
    return { valid: false, errorCode: 'profile_interval_invalid' };
  }
  const intervalSeconds = recordingProfileIntervalSeconds({ kind: profile, interval: intervalValue, unit });
  return intervalSeconds === null
    ? { valid: false, errorCode: 'profile_interval_invalid' }
    : { valid: true, intervalValue, intervalSeconds };
}

export function evaluateCandidateCounterConfirmation(workspaceId: string, candidate: GroupCandidate, measurements: MeasurementDefinition[]): CandidateCounterStatus {
  if (!['int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64'].includes(candidate.target_type)) {
    return { candidate, confirmed: false, issueCode: 'counter-type-unsupported' };
  }
  const matching = measurements.filter((i) => i.workspace_id === workspaceId && i.device_id === candidate.device_id && i.point_id === candidate.point_id && i.tag_id === candidate.tag_id);
  if (matching.length === 0) return { candidate, confirmed: false, issueCode: 'counter-semantics-required' };
  const target = matching.find((item) => item.semantic_kind === 'counter' || item.semantic_kind === 'signed_counter') ?? matching[0];
  if (target.semantic_kind !== 'counter' && target.semantic_kind !== 'signed_counter') {
    return { candidate, confirmed: false, measurement: target, issueCode: 'counter-measurement-not-counter' };
  }
  if (!isConfirmedCounterDefinition(target)) {
    return { candidate, confirmed: false, measurement: target, issueCode: 'counter-measurement-unconfirmed' };
  }
  return { candidate, confirmed: true, measurement: target };
}

export function isManagedProfileGroupForDevice(group: WriteGroup, deviceId: string): boolean {
  const kind = groupProfileKind(group);
  return group.status !== 'deleted' && purposeOfLayout(group.layout_version) !== undefined &&
    group.destination.storage_strategy === 'managed' &&
    (kind === 'instantaneous' || kind === 'cumulative') &&
    group.members.length > 0 && group.members.every((member) => member.device_id === deviceId);
}

/** The one device a managed profile group records, or undefined for a custom or cross-device group. */
export function managedProfileDeviceId(group: WriteGroup): string | undefined {
  const deviceId = group.members[0]?.device_id;
  return deviceId && isManagedProfileGroupForDevice(group, deviceId) ? deviceId : undefined;
}

export function buildCircuitProfileGroupDraft(params: {
  workspaceId: string;
  deviceName: string;
  profile: WriteGroupProfileKind;
  intervalValue: number;
  intervalSeconds: number;
  storageMode: WriteGroupStorageMode;
  recreateMissing: boolean;
  circuitId: string;
  circuitName: string;
  destination: WriteGroupDestinationDraft;
  selectedCandidates: GroupCandidate[];
  measurementsByPointKey: Map<string, MeasurementDefinition>;
  purpose: CircuitPurpose;
  /** `mean` for a new instantaneous profile; a saved profile passes the reduction it already has (`latest` = none). */
  reduction?: WriteGroupSampleReduction;
  memberRoles?: Record<string, WriteGroupValueRole>;
}): WriteGroupDraft {
  const { workspaceId, deviceName, profile, intervalValue, intervalSeconds, storageMode, recreateMissing, circuitId, circuitName, destination, selectedCandidates, measurementsByPointKey, purpose } = params;
  const labelBase = circuitName.trim() || deviceName.trim() || 'Circuit';
  const suffix = profile === 'cumulative' ? 'Cumulative' : 'Instantaneous';
  // Only an instantaneous profile averages, and only when it says so: every source then states its role.
  const mean = profile === 'instantaneous' && (params.reduction ?? 'mean') === 'mean';
  const wide = purpose === 'wide';
  const members: WriteGroupMemberDraft[] = selectedCandidates.map((candidate) => {
    const key = candidateKey(candidate);
    const role = mean ? (params.memberRoles?.[key] ?? defaultMemberRole(candidate)) : undefined;
    const counted = profile === 'cumulative' || role === 'counter_raw';
    return {
      device_id: candidate.device_id, point_id: candidate.point_id, tag_id: candidate.tag_id,
      measurement_id: counted ? (measurementsByPointKey.get(key)?.id ?? '') : '',
      target_column: '', required: false,
      ...(role ? { value_role: role } : {}),
      ...(circuitId ? { circuit_id: circuitId } : {}),
      ...(circuitName.trim() ? { circuit_name: circuitName.trim() } : {}),
    };
  });
  return {
    workspace_id: workspaceId,
    name: `${labelBase} (${suffix})`,
    layout_version: PURPOSE_LAYOUT[purpose],
    recording_profile: { kind: profile, interval: intervalValue, unit: profile === 'instantaneous' ? 'second' : 'minute', ...(mean ? { sample_reduction: 'mean' as const } : {}) },
    // The wide profile table has neither a storage mode nor a recreate policy to choose.
    ...(wide ? {} : { storage_mode: storageMode }),
    members,
    destination: {
      connector_id: destination.connector_id, connector_revision: destination.connector_revision,
      table_schema: destination.table_schema, table_name: '', storage_strategy: 'managed',
    },
    row_policy: { interval_seconds: intervalSeconds, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
    write_policy: wide ? { mode: 'append', dedupe_capability: 'receipt' } : { mode: 'append', dedupe_capability: 'receipt', recreate_missing: recreateMissing },
  };
}

/** An omitted reduction is the stored form of `latest`; an explicit mean is a different configuration. */
function effectiveSampleReduction(group: WriteGroup | WriteGroupDraft): WriteGroupSampleReduction {
  return group.recording_profile?.sample_reduction ?? 'latest';
}

export function canonicalProfileGroupMatchesDraft(group: WriteGroup, draft: WriteGroupDraft, connectorKind?: string): boolean {
  const schema = (value: string) => connectorKind === 'sqlite' && !value ? 'main' : value;
  // Only the two known profile layouts can match, and only the very layout the draft carries.
  if (groupProfileKind(group) !== groupProfileKind(draft) || purposeOfLayout(group.layout_version) === undefined || group.layout_version !== draft.layout_version ||
    effectiveSampleReduction(group) !== effectiveSampleReduction(draft) ||
    (group.storage_mode ?? 'plain') !== (draft.storage_mode ?? 'plain') ||
    (group.write_policy.recreate_missing ?? false) !== (draft.write_policy.recreate_missing ?? false) ||
    groupCircuitId(group) !== groupCircuitId(draft) || groupCircuitName(group) !== groupCircuitName(draft) ||
    group.destination.connector_id !== draft.destination.connector_id || group.destination.connector_revision !== draft.destination.connector_revision ||
    schema(group.destination.table_schema) !== schema(draft.destination.table_schema) || group.row_policy.interval_seconds !== draft.row_policy.interval_seconds ||
    group.members.length !== draft.members.length) return false;
  const draftMembersByKey = new Map(
    draft.members.map((m) => [`${m.device_id}\0${m.point_id}\0${m.tag_id}`, m]),
  );
  return group.members.every((member) => {
    const expected = draftMembersByKey.get(`${member.device_id}\0${member.point_id}\0${member.tag_id}`);
    return expected !== undefined && (member.measurement_id ?? '') === (expected.measurement_id ?? '') &&
      (member.value_role ?? '') === (expected.value_role ?? '');
  });
}

export function sameCircuitRecordingStartRequest(
  stored: RecordingStartRequest | undefined,
  workspaceId: string,
  workspaceRevision: string,
  settingsRevision: string | undefined,
  deviceId: string,
  targetGroups: WriteGroup[],
  operation?: RecordingStartOperation,
): boolean {
  if (!stored || stored.workspace_id !== workspaceId || stored.device_ids.length !== 1 || stored.device_ids[0] !== deviceId) return false;
  const ownsRevision = Boolean(operation?.setup_revision && workspaceRevision && operation.setup_revision === workspaceRevision) &&
    targetGroups.every((g) => operation?.groups.some((p) => p.group_id === g.id && p.group_revision === g.revision && (!p.applied_revision || p.applied_revision === g.applied_revision)));
  if (stored.expected_workspace_revision !== workspaceRevision &&
    !recordingStartRevisionsMatch(stored, settingsRevision, workspaceRevision) && !ownsRevision) return false;
  if (stored.groups.length !== targetGroups.length) return false;
  return targetGroups.every((target, idx) => {
    const item = stored.groups[idx];
    return Boolean(item && item.group_id === target.id && (item.expected_group_revision === target.revision || ownsRevision));
  });
}
