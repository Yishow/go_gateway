import type {
  WriteGroup,
  WriteGroupMemberDraft,
} from './studioV2WriteGroup';

export type WriteGroupMigrationPreviewStatus = 'needs_review' | 'blocked';

export interface WriteGroupMigrationFinding {
  code: string;
  message: string;
}

/** Persisted legacy intent, without connector credentials or DSNs. */
export interface WriteGroupMigrationIntent {
  source_id: string;
  device_id: string;
  point_id: string;
  tag_id: string;
  connector_id: string;
  connector_revision: string;
  database?: string;
  table_schema: string;
  table_name: string;
  column_name: string;
  write_mode: string;
  timestamp_column?: string;
  group_key?: string;
  write_interval_seconds?: number;
  interval_source?: string;
  enabled: boolean;
}

export interface WriteGroupMigrationRowGroup {
  id: string;
  connector_id: string;
  table_schema: string;
  table_name: string;
  member_point_ids: string[];
  group_key_columns: string[];
  unique_key_columns: string[];
}

/** Complete persisted row-group intent retained by the read-only preview. */
export interface WriteGroupMigrationRowGroupIntent {
  source_id: string;
  connector_id: string;
  connector_revision: string;
  row_group: WriteGroupMigrationRowGroup;
  members: WriteGroupMigrationIntent[];
}

/** Unsaved draft candidate; empty identity and Go zero timestamps are valid. */
export type WriteGroupMigrationCandidate = Omit<
  WriteGroup,
  'id' | 'revision' | 'applied_revision' | 'status' | 'members' | 'created_at' | 'updated_at'
> & {
  id: string;
  revision: string;
  applied_revision: string;
  status: 'draft';
  members: WriteGroupMemberDraft[];
  created_at: string;
  updated_at: string;
};

export interface WriteGroupMigrationPreviewItem {
  source_id: string;
  source_revision: string;
  status: WriteGroupMigrationPreviewStatus;
  differences: WriteGroupMigrationFinding[];
  issues: WriteGroupMigrationFinding[];
  before_intent?: WriteGroupMigrationIntent;
  before_row_group_intent?: WriteGroupMigrationRowGroupIntent;
  candidate_group?: WriteGroupMigrationCandidate;
}

export interface WriteGroupMigrationPreview {
  workspace_id: string;
  workspace_revision: string;
  connector_revision: string;
  adapter_version: string;
  review_digest: string;
  items: WriteGroupMigrationPreviewItem[];
}

export interface WriteGroupMigrationPreviewRequest {
  workspace_id: string;
  source_ids: string[];
}

/** Explicitly reviewed snapshot conversion request; server-owned group fields are excluded. */
export interface WriteGroupMigrationReviewRequest {
  workspace_id: string;
  expected_workspace_revision: string;
  expected_connector_revision: string;
  review_digest: string;
  source_ids: string[];
  confirm_snapshot_conversion: boolean;
}

/** Canonical persisted groups returned after a reviewed migration save. */
export interface WriteGroupMigrationReviewResponse {
  workspace_id: string;
  workspace_revision: string;
  connector_revision: string;
  groups: WriteGroup[];
}

/** JSON retained from an unsupported recording plan without applying defaults. */
export type RecordingPlanMigrationRawValue =
  | null
  | boolean
  | number
  | string
  | RecordingPlanMigrationRawValue[]
  | { [key: string]: RecordingPlanMigrationRawValue };

export type RecordingPlanMigrationRawPlan = {
  [key: string]: RecordingPlanMigrationRawValue;
};

export type RecordingPlanMigrationSourceStatus = 'resolved' | 'blocked';

/** Source-chain identity retained by a blocked recording-plan preview. */
export interface RecordingPlanMigrationSource {
  measurement_id: string;
  device_id: string;
  point_id: string;
  tag_id: string;
  definition_revision: string;
  source_binding_revision: string;
  series_epoch: string;
  source_revision?: string;
  mapping_revision?: string;
  status: RecordingPlanMigrationSourceStatus;
}

/** Complete persisted plan intent; the raw plan is never normalized or defaulted. */
export interface RecordingPlanMigrationIntent {
  source_id: string;
  plan: RecordingPlanMigrationRawPlan;
  sources: RecordingPlanMigrationSource[];
}

export type RecordingPlanMigrationRepairAction = 'open_write_groups';

export interface RecordingPlanMigrationPreviewItem extends Omit<
  WriteGroupMigrationPreviewItem,
  'before_intent' | 'before_row_group_intent' | 'candidate_group'
> {
  status: 'blocked';
  repair_action: RecordingPlanMigrationRepairAction;
  before_recording_plan_intent: RecordingPlanMigrationIntent;
}

export interface RecordingPlanMigrationPreview extends Omit<
  WriteGroupMigrationPreview,
  'adapter_version' | 'items'
> {
  adapter_version: 'recording-plan-v1';
  items: RecordingPlanMigrationPreviewItem[];
}

/** Compatibility name matching the backend result vocabulary. */
export type WriteGroupMigrationReviewResult = WriteGroupMigrationReviewResponse;
