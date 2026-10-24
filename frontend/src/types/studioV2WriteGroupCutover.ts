export type WriteGroupCutoverKind = 'forward' | 'rollback';

/** One source, by its stable identity. */
export interface CutoverMemberRef {
  device_id: string;
  point_id: string;
  tag_id: string;
}

/** A long group and the revision the reviewer saw. */
export interface CutoverTargetRef {
  group_id: string;
  expected_revision: string;
}

/**
 * Request of a cutover review and its apply. `boundary` is the exact instant,
 * kept as the text the server returned so no precision is lost; apply needs the
 * reviewed `boundary` and `review_digest`.
 */
export interface WriteGroupCutoverRequest {
  kind?: WriteGroupCutoverKind;
  workspace_id: string;
  expected_workspace_revision: string;
  source_group_id?: string;
  expected_source_revision: string;
  expected_connector_revision: string;
  targets?: CutoverTargetRef[];
  /** The forward cutover a rollback reverses. */
  cutover_id?: string;
  boundary?: string;
  review_digest?: string;
}

/** Stable reasons a cutover cannot be applied. */
export type CutoverIssueCode =
  | 'source_unavailable' | 'source_not_applied' | 'source_has_unapplied_draft' | 'source_already_ended' | 'source_is_long'
  | 'target_not_long' | 'target_already_active' | 'target_not_ready' | 'target_member_not_in_source' | 'targets_overlap'
  | 'unselected_sources' | 'boundary_not_future' | 'boundary_before_active_revision' | 'no_targets' | 'cutover_not_rollbackable'
  | 'remaining_draft_mismatch' | 'remaining_draft_not_ready' | 'remaining_draft_changes_policy'
  | 'remaining_draft_omits_required_column' | 'restored_group_not_ready';

export interface CutoverGroupView {
  group_id: string;
  name: string;
  revision: string;
  layout: string;
  storage_mode?: string;
  profile_kind?: string;
  interval_seconds: number;
  connector_id: string;
  table_schema: string;
  table_name: string;
  member_count: number;
  schema_confirmed: boolean;
}

/** The read-only difference an operator confirms. */
export interface WriteGroupCutoverReview {
  kind: WriteGroupCutoverKind;
  boundary: string;
  digest: string;
  applicable: boolean;
  issues: CutoverIssueCode[];
  source: CutoverGroupView;
  targets: CutoverGroupView[];
  moved: CutoverMemberRef[];
  /** Sources of the applied group that no target takes; they are never moved silently. */
  unselected: CutoverMemberRef[];
  /** Sources of the saved draft that keeps recording in the wide group; it must name exactly the unselected ones. */
  remaining: CutoverMemberRef[];
}

export interface CutoverRevisionRef {
  group_id: string;
  revision: string;
}

/** The immutable record of an applied cutover. */
export interface WriteGroupCutover {
  id: string;
  workspace_id: string;
  kind: WriteGroupCutoverKind;
  source_group_id: string;
  boundary: string;
  prior_cutover_id?: string;
  review_digest: string;
  ended: CutoverRevisionRef[];
  started: CutoverRevisionRef[];
  moved: CutoverMemberRef[];
  created_at: string;
}
