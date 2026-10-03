import type { WriteGroup, WriteGroupDraft, WriteGroupStorageStrategy } from '../../../../../types/studioV2WriteGroup';
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
  /** Persisted row-group contract fields without editor controls. */
  group_key_columns?: string[];
  unique_key_columns?: string[];
  value_column?: string;
  quality_column?: string;
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
    members: group.members.map((member) => ({
      key: memberKey(member.point_id, member.tag_id),
      device_id: member.device_id, point_id: member.point_id, tag_id: member.tag_id,
      ...(member.measurement_id !== undefined ? { measurement_id: member.measurement_id } : {}),
      entity_key: member.entity_key ?? '', target_column: member.target_column, required: member.required,
      ...(member.max_age_seconds !== undefined ? { max_age_seconds: member.max_age_seconds } : {}),
    })),
    ...(group.row_policy.group_key_columns !== undefined ? { group_key_columns: [...group.row_policy.group_key_columns] } : {}),
    ...(group.row_policy.unique_key_columns !== undefined ? { unique_key_columns: [...group.row_policy.unique_key_columns] } : {}),
    ...(group.row_policy.value_column !== undefined ? { value_column: group.row_policy.value_column } : {}),
    ...(group.row_policy.quality_column !== undefined ? { quality_column: group.row_policy.quality_column } : {}),
  };
}

/** The editable payload. Source and mapping revisions are not sent: the server pins the current ones. */
export function draftToRequestGroup(draft: EditorDraft, workspaceId: string): WriteGroupDraft {
  return {
    workspace_id: workspaceId,
    name: draft.name.trim(),
    members: draft.members.map((member) => ({
      device_id: member.device_id, point_id: member.point_id, tag_id: member.tag_id,
      ...(member.measurement_id !== undefined ? { measurement_id: member.measurement_id } : {}),
      ...(member.entity_key.trim() ? { entity_key: member.entity_key.trim() } : {}),
      target_column: member.target_column.trim(), required: member.required,
      ...(member.max_age_seconds !== undefined ? { max_age_seconds: member.max_age_seconds } : {}),
    })),
    destination: {
      connector_id: draft.connector_id, connector_revision: draft.connector_revision,
      table_schema: draft.table_schema, table_name: draft.table_name.trim(), storage_strategy: draft.storage_strategy,
    },
    row_policy: {
      interval_seconds: draft.interval_seconds, allowed_lateness_seconds: draft.allowed_lateness_seconds,
      incomplete_policy: draft.incomplete_policy,
      ...(draft.entity_key_column.trim() ? { entity_key_column: draft.entity_key_column.trim() } : {}),
      ...(draft.group_key_columns !== undefined ? { group_key_columns: [...draft.group_key_columns] } : {}),
      ...(draft.unique_key_columns !== undefined ? { unique_key_columns: [...draft.unique_key_columns] } : {}),
      ...(draft.value_column !== undefined ? { value_column: draft.value_column } : {}),
      ...(draft.quality_column !== undefined ? { quality_column: draft.quality_column } : {}),
      ...(draft.provenance_column.trim() ? { provenance_column: draft.provenance_column.trim() } : {}),
    },
    write_policy: { ...(draft.dedupe_capability ? { dedupe_capability: draft.dedupe_capability } : {}) },
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
  | 'column-conflict' | 'interval-invalid' | 'lateness-invalid' | 'partial-needs-provenance' | 'entity-column-required';

export interface DraftIssue {
  code: DraftIssueCode;
  member_key?: string;
  column?: string;
}

export interface DraftValidationContext {
  destinationSaved: boolean;
}

/** Local checks only; the backend readiness result stays the authority for activation. */
export function validateDraft(draft: EditorDraft, context: DraftValidationContext): DraftIssue[] {
  const issues: DraftIssue[] = [];
  if (!draft.name.trim()) issues.push({ code: 'name-required' });
  if (!context.destinationSaved || !draft.connector_id || !draft.connector_revision) issues.push({ code: 'destination-unsaved' });
  if (!draft.table_name.trim()) issues.push({ code: 'table-required' });
  if (draft.members.length === 0) issues.push({ code: 'members-required' });
  for (const member of draft.members) {
    if (!member.target_column.trim()) issues.push({ code: 'column-required', member_key: member.key });
  }
  for (const conflict of findColumnConflicts(draft.members.map((member) => ({
    key: member.key, column: member.target_column.trim(), entity_key: member.entity_key,
  })))) {
    issues.push({ code: 'column-conflict', column: conflict.column });
  }
  if (!Number.isFinite(draft.interval_seconds) || draft.interval_seconds < 1) issues.push({ code: 'interval-invalid' });
  if (!Number.isFinite(draft.allowed_lateness_seconds) || draft.allowed_lateness_seconds < 0) issues.push({ code: 'lateness-invalid' });
  if (draft.incomplete_policy === 'partial' && !draft.provenance_column.trim()) issues.push({ code: 'partial-needs-provenance' });
  if (draft.members.some((member) => member.entity_key.trim()) && !draft.entity_key_column.trim()) {
    issues.push({ code: 'entity-column-required' });
  }
  return issues;
}
