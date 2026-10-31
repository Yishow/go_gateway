import type { StudioV2WorkspaceReadinessIssue } from './studioV2WorkspaceReadiness';
import type { DataType } from './datalink';
import type { WriteGroupCalculationCapabilities } from './studioV2WriteGroupCalculation';
import type { CounterPolicy } from './measurement';

export type WriteGroupStatus = 'draft' | 'ready' | 'running' | 'disabled' | 'deleted';

export type WriteGroupStorageStrategy = 'managed' | 'custom';

/**
 * An absent layout always means `managed_wide_v1`. `managed_wide_profile_v1` is the
 * reduction-aware managed wide profile table; `managed_long_v1` is the fixed typed
 * table shared by the canonical profile groups of one storage scope.
 */
export type WriteGroupLayoutVersion = 'managed_wide_v1' | 'managed_wide_profile_v1' | 'managed_long_v1';

/** How a long layout is stored. Each non-default strategy is explicit and capability-gated. */
export type WriteGroupStorageMode = 'plain' | 'timescale' | 'native_monthly';

/** Legacy circuit-profile selector kinds; mixed uses the explicit recording-profile union below. */
export type WriteGroupProfileKind = 'instantaneous' | 'cumulative';

/** Each kind has one fixed unit: instantaneous in seconds, cumulative in minutes. */
export type WriteGroupIntervalUnit = 'second' | 'minute';

/** How a profile's samples are reduced; absent means `latest` for every legacy group. */
export type WriteGroupSampleReduction = 'latest' | 'mean';

/** Persisted role of one member's value inside a profile group; absent on groups saved before roles. */
export type WriteGroupValueRole = 'gauge' | 'counter_raw' | 'state';

/** Stored result type, independent of the source Tag type. */
export type WriteGroupResultValueType = 'bool' | 'text' | 'int64' | 'uint64' | 'decimal' | 'float64';

/** Backend-owned confirmation that an integer counter has an exact decimal source scale. */
export interface WriteGroupExactRecordingSource {
  measurement_id: string;
  definition_revision: string;
  source_binding_revision: string;
  series_epoch: string;
  original_data_type: DataType;
  scale: string;
  unit: string;
  semantic_kind: 'counter' | 'signed_counter';
  counter_policy?: CounterPolicy;
}

/** Recording profile of a canonical circuit group; the interval is in `unit`, never in seconds by default. */
export interface WriteGroupLegacyRecordingProfile {
  kind: WriteGroupProfileKind;
  interval: number;
  unit: WriteGroupIntervalUnit;
  /** Explicit reduction, carried unchanged; absent decodes as `latest` and is never rewritten as a choice. */
  sample_reduction?: WriteGroupSampleReduction;
}

/** Explicit per-member shape: no synthetic group interval or reduction. */
export interface WriteGroupMixedRecordingProfile {
  kind: 'mixed';
  result_shape: 'member-calculation-v1';
  interval?: never;
  unit?: never;
  sample_reduction?: never;
}
export type WriteGroupRecordingProfile = WriteGroupLegacyRecordingProfile | WriteGroupMixedRecordingProfile;
export type WriteGroupMemberCalculation = 'latest' | 'mean' | 'delta';
export interface WriteGroupMemberCadence {
  interval: number;
  unit: WriteGroupIntervalUnit;
  alignment: 'utc_epoch';
}

export interface WriteGroupMember {
  calculation?: WriteGroupMemberCalculation;
  cadence?: WriteGroupMemberCadence;
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
  /** Opaque workspace-scoped circuit identity; `circuit_name` is display text only. */
  circuit_id?: string;
  circuit_name?: string;
  /** Persisted role of this value; absent on groups saved before roles. */
  value_role?: WriteGroupValueRole;
  /** The Tag's own `schema.DataType` at save time; the Tag type itself is never changed. */
  source_value_type?: DataType;
  /** Stored result type of this member, independent of `source_value_type`. */
  result_value_type?: WriteGroupResultValueType;
  /** Backend-owned source proof for confirmed exact decimal integer counters. */
  exact_recording_source?: WriteGroupExactRecordingSource;
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
  /** Absent only on explicit mixed profiles; never default mixed to a group window. */
  interval_seconds?: number;
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
  due_at_column?: string;
  group_id_column?: string;
  device_id_column?: string;
}

export interface WriteGroupWritePolicy {
  mode?: string;
  dedupe_capability?: string;
  /**
   * Managed long storage only. `true` and `false` are decisions made at schema
   * confirmation; absent means none was ever made and never reads as `true`.
   */
  recreate_missing?: boolean;
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
  /** Canonical identity of an independent mixed recording. */
  record_id?: string;
  /** Server-owned persistent basic role; absent for pre-existing advanced groups. */
  basic_managed_device_id?: string;
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
  layout_version?: WriteGroupLayoutVersion;
  /** Legacy profiles derive the group interval; mixed records carry independent member cadences. */
  recording_profile?: WriteGroupRecordingProfile;
  /** Set by the server on long-layout groups; absent on wide groups. */
  storage_mode?: WriteGroupStorageMode;
  created_at: string;
  updated_at: string;
}

/** Editable fields accepted by create and update; server-owned identity is excluded. */
export type WriteGroupDraft = Omit<
  WriteGroup,
  'id' | 'revision' | 'applied_revision' | 'status' | 'migration' | 'created_at' | 'updated_at' |
  'members' | 'destination' | 'basic_managed_device_id'
> & {
  members: WriteGroupMemberDraft[];
  destination: WriteGroupDestinationDraft;
};

export interface WriteGroupListResponse {
  calculation_capabilities?: WriteGroupCalculationCapabilities;
  workspace_id: string;
  workspace_revision: string;
  groups: WriteGroup[];
}

export interface WriteGroupMutationResponse {
  calculation_capabilities?: WriteGroupCalculationCapabilities;
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
