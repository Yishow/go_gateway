import type { DataType, DatabaseWriteMode, TransformStep } from './datalink';

export type SourceRuleCandidateSetStatus = 'ready' | 'blocked' | 'deferred';
export type SourceRuleMappingStatus =
  | 'draft'
  | 'validated'
  | 'active'
  | 'out_of_sync'
  | 'error';

export interface SourceRuleCandidateIdentity {
  source_rule_id: string;
  candidate_type: string;
  candidate_kind: string;
  derived_from_rule_address: string;
  target_binding_scope?: SourceRuleCandidateScopeField[];
}

export interface SourceRuleCandidateScopeField {
  key: string;
  value: string;
}

export interface SourceRuleTagCandidateView {
  id: string;
  identity: SourceRuleCandidateIdentity;
  proposed_signature: string;
  address: string;
  point_id: string;
  tag_id?: string;
  mapping_id?: string;
  tag_key: string;
  display_name: string;
  data_type: DataType;
  transform_pipeline?: TransformStep[];
  status: SourceRuleMappingStatus;
  last_applied_signature?: string;
  blocking_reason?: string;
}

export type SourceRuleOutputStatus = 'ready' | 'blocked' | 'out_of_sync';

export interface SourceRuleDatabaseOutputCandidateView {
  id: string;
  identity: SourceRuleCandidateIdentity;
  proposed_signature: string;
  address: string;
  point_id: string;
  tag_id?: string;
  mapping_id?: string;
  tag_key: string;
  display_name: string;
  data_type: DataType;
  status: SourceRuleOutputStatus;
  blocking_reason?: string;
  connector_id?: string;
  table_schema?: string;
  table_name?: string;
  column_name?: string;
  write_mode?: DatabaseWriteMode;
  timestamp_column?: string;
}

export interface SourceRuleCandidateSetView<TCandidate> {
  status: SourceRuleCandidateSetStatus;
  reason?: string;
  candidates: TCandidate[];
}

export interface SourceRuleCandidateSnapshotView {
  source_rule_id: string;
  revision_id: string;
  tags: SourceRuleCandidateSetView<SourceRuleTagCandidateView>;
  database_outputs: SourceRuleCandidateSetView<SourceRuleDatabaseOutputCandidateView>;
  local_modbus_outputs: SourceRuleCandidateSetView<Record<string, unknown>>;
}
