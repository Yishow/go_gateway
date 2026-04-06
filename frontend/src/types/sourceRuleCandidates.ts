import type { DataType, TransformStep } from './datalink';

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
  target_binding_scope?: string[];
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

export interface SourceRuleCandidateSetView<TCandidate> {
  status: SourceRuleCandidateSetStatus;
  reason?: string;
  candidates: TCandidate[];
}

export interface SourceRuleCandidateSnapshotView {
  source_rule_id: string;
  revision_id: string;
  tags: SourceRuleCandidateSetView<SourceRuleTagCandidateView>;
  database_outputs: SourceRuleCandidateSetView<Record<string, unknown>>;
  local_modbus_outputs: SourceRuleCandidateSetView<Record<string, unknown>>;
}
