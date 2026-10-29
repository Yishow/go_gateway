import type { WriteGroup, WriteGroupStatus, WriteGroupMember, WriteGroupLayoutVersion } from '../types/studioV2WriteGroup';
import { boundedString, MAX_SAFE_JSON_ARRAY_LENGTH } from './safeJson';
import { isRecord, readPresentText, readRequiredText, readOptionalText } from './studioV2WriteGroupJsonPrimitives';
import { parseRecordingProfile, recordingProfileIntervalSeconds } from './studioV2RecordingProfile';
import { parseWriteGroupMember, parseWriteGroupDestination, parseWriteGroupRowPolicy, parseWriteGroupWritePolicy, parseWriteGroupMigration } from './studioV2WriteGroupDefinitionFields';
const WRITE_GROUP_STATUSES = new Set<WriteGroupStatus>(['draft', 'ready', 'running', 'disabled', 'deleted']);
const WRITE_GROUP_LAYOUT_VERSIONS = new Set<string>(['managed_wide_v1', 'managed_wide_profile_v1', 'managed_long_v1']);

/** Saved responses always require real persistent identity and table identity. */
export function parseStudioV2WriteGroup(value: unknown): WriteGroup | null { return parseDefinition(value, false); }

/** Only the server-sanitized candidate shape; never a persisted group response. */
export function parseWriteGroupSuccessorCandidate(value: unknown): WriteGroup | null {
  if (!isRecord(value) || value.id !== '' || value.revision !== '' || value.applied_revision !== '' || value.status !== 'draft' ||
    value.created_at !== '0001-01-01T00:00:00Z' || value.updated_at !== '0001-01-01T00:00:00Z' ||
    !isRecord(value.migration) || Object.keys(value.migration).length !== 0 || !isRecord(value.destination) ||
    value.destination.schema_revision !== undefined || value.destination.schema_digest !== undefined) return null;
  const group = parseDefinition(value, true);
  if (!group || group.destination.storage_strategy !== 'managed' ||
    (group.layout_version !== 'managed_long_v1' && group.destination.table_name !== '')) return null;
  return group;
}

function parseDefinition(value: unknown, candidate: boolean): WriteGroup | null {
  if (!isRecord(value)) return null;
  const id = candidate ? readPresentText(value, 'id') : readRequiredText(value, 'id');
  const workspaceId = readRequiredText(value, 'workspace_id');
  const revision = candidate ? readPresentText(value, 'revision') : readRequiredText(value, 'revision');
  const appliedRevision = readPresentText(value, 'applied_revision');
  const name = readRequiredText(value, 'name');
  const status = readRequiredText(value, 'status');
  const createdAt = readRequiredText(value, 'created_at');
  const updatedAt = readRequiredText(value, 'updated_at');
  const rawMembers = value.members;
  const basicDevice = value.basic_managed_device_id === undefined
    ? undefined : boundedString(value.basic_managed_device_id);
  if (value.basic_managed_device_id !== undefined && !basicDevice) return null;
  if (id === null || id === undefined || !workspaceId || revision === null || revision === undefined || appliedRevision === undefined || appliedRevision === null || !name ||
    !status || !WRITE_GROUP_STATUSES.has(status as WriteGroupStatus) || !createdAt || !updatedAt ||
    !Array.isArray(rawMembers) || rawMembers.length > MAX_SAFE_JSON_ARRAY_LENGTH) {
    return null;
  }

  const profile = parseRecordingProfile(value.recording_profile);
  const profileKind = profile?.kind;
  const members = rawMembers.map((member) => parseWriteGroupMember(member, profileKind));
  if (members.some((member) => member === null)) return null;
  const destination = parseWriteGroupDestination(value.destination, candidate);
  const rowPolicy = parseWriteGroupRowPolicy(value.row_policy, status as WriteGroupStatus);
  const writePolicy = parseWriteGroupWritePolicy(value.write_policy);
  const migration = parseWriteGroupMigration(value.migration);
  if (!destination || !rowPolicy || !writePolicy || !migration) return null;
  const layoutVersion = readOptionalText(value, 'layout_version');
  const storageMode = readOptionalText(value, 'storage_mode');
  // An unknown layout or mode, or an interval whose unit was misread, must never reach the editor.
  if (layoutVersion === null || (layoutVersion !== undefined && !WRITE_GROUP_LAYOUT_VERSIONS.has(layoutVersion)) ||
    storageMode === null || (storageMode !== undefined && storageMode !== 'plain' && storageMode !== 'timescale' && storageMode !== 'native_monthly') ||
    profile === null || (profile && recordingProfileIntervalSeconds(profile) !== rowPolicy.interval_seconds)) {
    return null;
  }
  // The long table serves profile groups with a mode; a wide group carries neither a mode nor a long layout.
  // The reduction-aware wide profile layout exists only as a managed profile table, never with a storage mode.
  if ((layoutVersion === 'managed_long_v1') !== (storageMode !== undefined) ||
    (layoutVersion === 'managed_long_v1' && !profile) ||
    (layoutVersion === 'managed_wide_profile_v1' && (!profile || destination.storage_strategy !== 'managed'))) {
    return null;
  }
  const parsedMembers = members as WriteGroupMember[];
  if (parsedMembers.some((member) => member.exact_recording_source !== undefined) &&
    (layoutVersion !== 'managed_long_v1' || destination.storage_strategy !== 'managed')) return null;
  if (profile?.kind === 'cumulative' && parsedMembers.some((member) => member.value_role === 'counter_raw')) return null;
  if (parsedMembers.some((member) => member.exact_recording_source !== undefined &&
    (profile?.kind === 'instantaneous' ? member.value_role !== 'counter_raw' :
      profile?.kind === 'cumulative' ? member.value_role !== undefined : true))) return null;
  if (layoutVersion === 'managed_long_v1' && parsedMembers.some((member) =>
    member.result_value_type === 'decimal' && member.source_value_type !== 'float32' &&
    member.source_value_type !== 'float64' && member.exact_recording_source === undefined &&
    ((profile?.kind === 'cumulative' && member.value_role === undefined) ||
      (profile?.kind === 'instantaneous' && member.value_role === 'counter_raw')))) return null;

  return {
    id,
    workspace_id: workspaceId,
    revision,
    applied_revision: appliedRevision,
    ...(basicDevice ? { basic_managed_device_id: basicDevice } : {}),
    name,
    status: status as WriteGroupStatus,
    members: members as WriteGroupMember[],
    destination,
    row_policy: rowPolicy,
    write_policy: writePolicy,
    migration,
    ...(layoutVersion ? { layout_version: layoutVersion as WriteGroupLayoutVersion } : {}),
    ...(profile ? { recording_profile: profile } : {}),
    ...(storageMode ? { storage_mode: storageMode } : {}),
    created_at: createdAt,
    updated_at: updatedAt,
  };
}

