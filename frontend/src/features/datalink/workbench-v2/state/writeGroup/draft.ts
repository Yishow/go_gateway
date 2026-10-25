import type {
  WriteGroup,
  WriteGroupDraft,
  WriteGroupLayoutVersion,
  WriteGroupRecordingProfile,
  WriteGroupResultValueType,
  WriteGroupStorageMode,
  WriteGroupStorageStrategy,
  WriteGroupValueRole,
} from '../../../../../types/studioV2WriteGroup';
import type { DataType } from '../../../../../types/datalink';
import { isValidRecordingProfile, recordingProfileIntervalSeconds } from '../../../../../utils/studioV2RecordingProfile';
import { findColumnConflicts } from './columns';

export type IncompletePolicy = 'skip_row' | 'partial';

/** One member as the editor holds it. `key` is the persisted point and tag IDs. */
export interface EditorMember {
  key: string;
  device_id: string;
  point_id: string;
  tag_id: string;
  /** Persisted by migrated groups but not currently exposed in this editor. */
  measurement_id?: string;
  entity_key: string;
  target_column: string;
  required: boolean;
  max_age_seconds?: number;
  /** Persisted circuit identity of a profile group; the editor carries it unchanged. */
  circuit_id?: string;
  circuit_name?: string;
  /** Persisted member role and value types of a profile group; carried unchanged, never decided here. */
  value_role?: WriteGroupValueRole;
  source_value_type?: DataType;
  result_value_type?: WriteGroupResultValueType;
}

export interface EditorDraft {
  name: string;
  connector_id: string;
  connector_revision: string;
  table_schema: string;
  table_name: string;
  storage_strategy: WriteGroupStorageStrategy;
  interval_seconds: number;
  allowed_lateness_seconds: number;
  incomplete_policy: IncompletePolicy;
  entity_key_column: string;
  provenance_column: string;
  dedupe_capability: string;
  /** Persisted recovery decision of managed long storage; the editor only carries it, never decides it. */
  recreate_missing?: boolean;
  /** Persisted row-group contract fields without editor controls. */
  group_key_columns?: string[];
  unique_key_columns?: string[];
  value_column?: string;
  quality_column?: string;
  /** Persisted managed layout bindings; the editor never invents these. */
  record_key_column?: string;
  bucket_start_column?: string;
  group_id_column?: string;
  device_id_column?: string;
  /** Persisted profile contract. When present, the profile alone decides `interval_seconds`. */
  layout_version?: WriteGroupLayoutVersion;
  recording_profile?: WriteGroupRecordingProfile;
  storage_mode?: WriteGroupStorageMode;
  members: EditorMember[];
}

/** Visible default for a new group's row interval; a configured value, not a recommendation. */
export const DEFAULT_INTERVAL_SECONDS = 60;

export function memberKey(pointId: string, tagId: string): string {
  return `${pointId}|${tagId}`;
}

export function newDraft(destination: { connector_id: string; connector_revision: string; table_schema: string; table_name: string }): EditorDraft {
  return {
    name: '', ...destination, storage_strategy: 'custom',
    interval_seconds: DEFAULT_INTERVAL_SECONDS, allowed_lateness_seconds: 0, incomplete_policy: 'skip_row',
    entity_key_column: '', provenance_column: '', dedupe_capability: '', members: [],
  };
}

export function groupToDraft(group: WriteGroup): EditorDraft {
  return {
    name: group.name,
    connector_id: group.destination.connector_id,
    connector_revision: group.destination.connector_revision,
    table_schema: group.destination.table_schema,
    table_name: group.destination.table_name,
    storage_strategy: group.destination.storage_strategy,
    interval_seconds: group.row_policy.interval_seconds,
    allowed_lateness_seconds: group.row_policy.allowed_lateness_seconds,
    incomplete_policy: group.row_policy.incomplete_policy === 'partial' ? 'partial' : 'skip_row',
    entity_key_column: group.row_policy.entity_key_column ?? '',
    provenance_column: group.row_policy.provenance_column ?? '',
    dedupe_capability: group.write_policy.dedupe_capability ?? '',
    ...(group.write_policy.recreate_missing !== undefined ? { recreate_missing: group.write_policy.recreate_missing } : {}),
    members: group.members.map((member) => ({
      key: memberKey(member.point_id, member.tag_id),
      device_id: member.device_id, point_id: member.point_id, tag_id: member.tag_id,
      ...(member.measurement_id !== undefined ? { measurement_id: member.measurement_id } : {}),
      entity_key: member.entity_key ?? '', target_column: member.target_column, required: member.required,
      ...(member.max_age_seconds !== undefined ? { max_age_seconds: member.max_age_seconds } : {}),
      ...(member.circuit_id !== undefined ? { circuit_id: member.circuit_id } : {}),
      ...(member.circuit_name !== undefined ? { circuit_name: member.circuit_name } : {}),
      ...(member.value_role !== undefined ? { value_role: member.value_role } : {}),
      ...(member.source_value_type !== undefined ? { source_value_type: member.source_value_type } : {}),
      ...(member.result_value_type !== undefined ? { result_value_type: member.result_value_type } : {}),
    })),
    ...(group.layout_version !== undefined ? { layout_version: group.layout_version } : {}),
    ...(group.recording_profile !== undefined ? { recording_profile: { ...group.recording_profile } } : {}),
    ...(group.storage_mode !== undefined ? { storage_mode: group.storage_mode } : {}),
    ...(group.row_policy.group_key_columns !== undefined ? { group_key_columns: [...group.row_policy.group_key_columns] } : {}),
    ...(group.row_policy.unique_key_columns !== undefined ? { unique_key_columns: [...group.row_policy.unique_key_columns] } : {}),
    ...(group.row_policy.value_column !== undefined ? { value_column: group.row_policy.value_column } : {}),
    ...(group.row_policy.quality_column !== undefined ? { quality_column: group.row_policy.quality_column } : {}),
    ...(group.row_policy.record_key_column !== undefined ? { record_key_column: group.row_policy.record_key_column } : {}),
    ...(group.row_policy.bucket_start_column !== undefined ? { bucket_start_column: group.row_policy.bucket_start_column } : {}),
    ...(group.row_policy.group_id_column !== undefined ? { group_id_column: group.row_policy.group_id_column } : {}),
    ...(group.row_policy.device_id_column !== undefined ? { device_id_column: group.row_policy.device_id_column } : {}),
  };
}

/** The editable payload. Source and mapping revisions are not sent: the server pins the current ones. */
export function draftToRequestGroup(draft: EditorDraft, workspaceId: string): WriteGroupDraft {
  // The profile is the only authority for the row interval; the editable seconds
  // field can never contradict it, and an invalid profile is blocked by validateDraft.
  const profileSeconds = draft.recording_profile ? recordingProfileIntervalSeconds(draft.recording_profile) : null;
  return {
    workspace_id: workspaceId,
    name: draft.name.trim(),
    members: draft.members.map((member) => ({
      device_id: member.device_id, point_id: member.point_id, tag_id: member.tag_id,
      ...(member.measurement_id !== undefined ? { measurement_id: member.measurement_id } : {}),
      ...(member.entity_key.trim() ? { entity_key: member.entity_key.trim() } : {}),
      target_column: member.target_column.trim(), required: member.required,
      ...(member.max_age_seconds !== undefined ? { max_age_seconds: member.max_age_seconds } : {}),
      ...(member.circuit_id !== undefined ? { circuit_id: member.circuit_id } : {}),
      ...(member.circuit_name !== undefined ? { circuit_name: member.circuit_name } : {}),
      ...(member.value_role !== undefined ? { value_role: member.value_role } : {}),
      ...(member.source_value_type !== undefined ? { source_value_type: member.source_value_type } : {}),
      ...(member.result_value_type !== undefined ? { result_value_type: member.result_value_type } : {}),
    })),
    destination: {
      connector_id: draft.connector_id, connector_revision: draft.connector_revision,
      table_schema: draft.table_schema, table_name: draft.table_name.trim(), storage_strategy: draft.storage_strategy,
    },
    row_policy: {
      interval_seconds: profileSeconds ?? draft.interval_seconds, allowed_lateness_seconds: draft.allowed_lateness_seconds,
      incomplete_policy: draft.incomplete_policy,
      ...(draft.entity_key_column.trim() ? { entity_key_column: draft.entity_key_column.trim() } : {}),
      ...(draft.group_key_columns !== undefined ? { group_key_columns: [...draft.group_key_columns] } : {}),
      ...(draft.unique_key_columns !== undefined ? { unique_key_columns: [...draft.unique_key_columns] } : {}),
      ...(draft.value_column !== undefined ? { value_column: draft.value_column } : {}),
      ...(draft.quality_column !== undefined ? { quality_column: draft.quality_column } : {}),
      ...(draft.provenance_column.trim() ? { provenance_column: draft.provenance_column.trim() } : {}),
      ...(draft.record_key_column?.trim() ? { record_key_column: draft.record_key_column.trim() } : {}),
      ...(draft.bucket_start_column?.trim() ? { bucket_start_column: draft.bucket_start_column.trim() } : {}),
      ...(draft.group_id_column?.trim() ? { group_id_column: draft.group_id_column.trim() } : {}),
      ...(draft.device_id_column?.trim() ? { device_id_column: draft.device_id_column.trim() } : {}),
    },
    write_policy: {
      ...(draft.dedupe_capability ? { dedupe_capability: draft.dedupe_capability } : {}),
      ...(draft.recreate_missing !== undefined ? { recreate_missing: draft.recreate_missing } : {}),
    },
    ...(draft.layout_version !== undefined ? { layout_version: draft.layout_version } : {}),
    ...(draft.recording_profile !== undefined ? { recording_profile: { ...draft.recording_profile } } : {}),
    ...(draft.storage_mode !== undefined ? { storage_mode: draft.storage_mode } : {}),
  };
}

function comparable(draft: EditorDraft): string {
  return JSON.stringify(draftToRequestGroup(draft, ''));
}

/** True only when the draft would change what the server stores. */
export function draftIsDirty(draft: EditorDraft, saved: WriteGroup): boolean {
  return comparable(draft) !== comparable(groupToDraft(saved));
}

export type DraftIssueCode =
  | 'name-required' | 'members-required' | 'destination-unsaved' | 'table-required' | 'column-required'
  | 'column-conflict' | 'interval-invalid' | 'profile-interval-invalid' | 'lateness-invalid' | 'partial-needs-provenance' | 'entity-column-required';

export interface DraftIssue {
  code: DraftIssueCode;
  member_key?: string;
  column?: string;
}

export interface DraftValidationContext {
  destinationSaved: boolean;
  storageStrategy?: WriteGroupStorageStrategy;
}

/** Local checks only; the backend readiness result stays the authority for activation. */
export function validateDraft(draft: EditorDraft, context: DraftValidationContext): DraftIssue[] {
  const issues: DraftIssue[] = [];
  const managed = (context.storageStrategy ?? draft.storage_strategy) === 'managed';
  if (!draft.name.trim()) issues.push({ code: 'name-required' });
  if (!context.destinationSaved || !draft.connector_id || !draft.connector_revision) issues.push({ code: 'destination-unsaved' });
  if (!managed && !draft.table_name.trim()) issues.push({ code: 'table-required' });
  if (draft.members.length === 0) issues.push({ code: 'members-required' });
  if (!managed) {
    for (const member of draft.members) {
      if (!member.target_column.trim()) issues.push({ code: 'column-required', member_key: member.key });
    }
    for (const conflict of findColumnConflicts(draft.members.map((member) => ({
      key: member.key, column: member.target_column.trim(), entity_key: member.entity_key,
    })))) {
      issues.push({ code: 'column-conflict', column: conflict.column });
    }
  }
  if (draft.recording_profile) {
    if (!isValidRecordingProfile(draft.recording_profile)) issues.push({ code: 'profile-interval-invalid' });
  } else if (!Number.isFinite(draft.interval_seconds) || draft.interval_seconds < 1) {
    issues.push({ code: 'interval-invalid' });
  }
  if (!Number.isFinite(draft.allowed_lateness_seconds) || draft.allowed_lateness_seconds < 0) issues.push({ code: 'lateness-invalid' });
  if (!managed && draft.incomplete_policy === 'partial' && !draft.provenance_column.trim()) issues.push({ code: 'partial-needs-provenance' });
  if (!managed && draft.members.some((member) => member.entity_key.trim()) && !draft.entity_key_column.trim()) {
    issues.push({ code: 'entity-column-required' });
  }
  return issues;
}
