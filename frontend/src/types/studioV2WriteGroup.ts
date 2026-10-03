import type { StudioV2WorkspaceReadinessIssue } from './studioV2WorkspaceReadiness';

export type WriteGroupStatus = 'draft' | 'ready' | 'running' | 'disabled' | 'deleted';

export type WriteGroupStorageStrategy = 'managed' | 'custom';

export interface WriteGroupMember {
  device_id: string;
  point_id: string;
  tag_id: string;
  entity_key?: string;
  source_revision: string;
  mapping_revision: string;
  measurement_id?: string;
  target_column: string;
  required: boolean;
  max_age_seconds?: number;
}

export interface WriteGroupDestination {
  connector_id: string;
  connector_revision: string;
  database: string;
  table_schema: string;
  table_name: string;
  storage_strategy: WriteGroupStorageStrategy;
  schema_revision?: string;
  schema_digest?: string;
}

export type WriteGroupMemberDraft = Omit<WriteGroupMember, 'source_revision' | 'mapping_revision'> & {
  source_revision?: string;
  mapping_revision?: string;
};

export type WriteGroupDestinationDraft = Omit<WriteGroupDestination, 'database' | 'schema_revision' | 'schema_digest'>;

export interface WriteGroupRowPolicy {
  interval_seconds: number;
  allowed_lateness_seconds: number;
  incomplete_policy?: string;
  entity_key_column?: string;
  group_key_columns?: string[];
  unique_key_columns?: string[];
  value_column?: string;
  quality_column?: string;
  provenance_column?: string;
  /** Server-owned managed layout bindings. They are retained when editing a saved group. */
  record_key_column?: string;
  bucket_start_column?: string;
  group_id_column?: string;
  device_id_column?: string;
}

export interface WriteGroupWritePolicy {
  mode?: string;
  dedupe_capability?: string;
}

export interface WriteGroupMigration {
  source_kind?: string;
  source_ids?: string[];
  source_revision?: string;
  adapter_version?: string;
  review_result?: string;
  legacy_row_group_id?: string;
  target_mapping_points?: Record<string, string>;
}

/** Canonical persisted authority returned by the workspace write-group API. */
export interface WriteGroup {
  id: string;
  workspace_id: string;
  revision: string;
  applied_revision: string;
  name: string;
  status: WriteGroupStatus;
  members: WriteGroupMember[];
  destination: WriteGroupDestination;
  row_policy: WriteGroupRowPolicy;
  write_policy: WriteGroupWritePolicy;
  migration: WriteGroupMigration;
  created_at: string;
  updated_at: string;
}

/** Editable fields accepted by create and update; server-owned identity is excluded. */
export type WriteGroupDraft = Omit<
  WriteGroup,
  'id' | 'revision' | 'applied_revision' | 'status' | 'migration' | 'created_at' | 'updated_at' |
  'members' | 'destination'
> & {
  members: WriteGroupMemberDraft[];
  destination: WriteGroupDestinationDraft;
};

export interface WriteGroupListResponse {
  workspace_id: string;
  workspace_revision: string;
  groups: WriteGroup[];
}

export interface WriteGroupMutationResponse {
  workspace_revision: string;
  group: WriteGroup;
}

export interface WriteGroupCreateRequest {
  workspace_id: string;
  /** Empty is valid for the first save of a newly initialized workspace. */
  expected_workspace_revision: string;
  expected_connector_revision: string;
  group: WriteGroupDraft;
}

export interface WriteGroupUpdateRequest {
  workspace_id: string;
  expected_workspace_revision: string;
  expected_group_revision: string;
  expected_connector_revision: string;
  group: WriteGroupDraft;
}

export interface WriteGroupDeleteRequest {
  workspace_id: string;
  expected_workspace_revision: string;
  expected_group_revision: string;
  expected_connector_revision: string;
}

/** Revision-only scope accepted by canonical group schema preview. */
export interface WriteGroupSchemaPreviewRequest {
  workspace_id: string;
  expected_workspace_revision: string;
  expected_group_revision: string;
  expected_connector_revision: string;
}

/** Explicit confirmation of the exact preview operation. */
export interface WriteGroupSchemaApplyRequest extends WriteGroupSchemaPreviewRequest {
  token: string;
  operation_id: string;
}

/** Read-only local configuration and verified destination schema gate. */
export interface WriteGroupReadiness {
  workspace_id: string;
  workspace_revision: string;
  group_id: string;
  group_revision: string;
  applied_revision: string;
  config_ready: boolean;
  schema_ready: boolean;
  ready: boolean;
  schema_digest?: string;
  issues: StudioV2WorkspaceReadinessIssue[];
}
