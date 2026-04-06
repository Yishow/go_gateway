export type SourceRuleTagReviewDecisionAction = 'rename' | 'skip' | 'override';

export interface SourceRuleTagReviewDecision {
  source_rule_id: string;
  candidate_id: string;
  action: SourceRuleTagReviewDecisionAction;
  tag_key?: string;
  override_tag_id?: string;
  stale: boolean;
  stale_revision_id?: string;
  stale_at?: string;
  created_at: string;
  updated_at: string;
}

export interface UpsertSourceRuleTagReviewDecisionRequest {
  candidate_id: string;
  action: SourceRuleTagReviewDecisionAction;
  tag_key?: string;
  override_tag_id?: string;
}
