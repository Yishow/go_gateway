import type { CutoverIssueCode } from '../../../../types/studioV2WriteGroupCutover';
import { normalizeTypedEnvelope } from '../../../../utils/safeJson';
import { classifyBatchMutationError } from './circuitProfileBatchOutcome';

/** Where an operator repairs a reason a cutover cannot be applied. */
export type CutoverRepair = 'source' | 'target' | 'boundary' | 'none';

export const CUTOVER_ISSUE_REPAIR: Record<CutoverIssueCode, CutoverRepair> = {
  source_unavailable: 'none', source_not_applied: 'source', source_has_unapplied_draft: 'source', source_already_ended: 'none', source_is_long: 'source',
  target_not_long: 'target', target_already_active: 'target', target_not_ready: 'target', target_member_not_in_source: 'target', targets_overlap: 'target',
  unselected_sources: 'source', boundary_not_future: 'boundary', boundary_before_active_revision: 'boundary', no_targets: 'target', cutover_not_rollbackable: 'none',
  remaining_draft_mismatch: 'source', remaining_draft_not_ready: 'source', remaining_draft_changes_policy: 'source',
  remaining_draft_omits_required_column: 'source', restored_group_not_ready: 'source',
};

const KNOWN_ERRORS = new Set(['revision_mismatch', 'cutover_not_applicable', 'cutover_required']);

/**
 * Classifies an apply error conservatively: only validated explicit application 4xx (except 408)
 * with a structured known envelope counts as a proven failure. Network loss, status-only, HTML,
 * bad response, timeout (408), 5xx, or unknown codes are classified as unknown.
 */
export function isCutoverUnknown(cause: unknown): boolean {
  return classifyBatchMutationError(cause) === 'unknown';
}

/** An operator sentence for a failed review or apply; the raw server text, SQL and DSNs are never shown. */
export function cutoverErrorMessage(cause: unknown, t: (key: string) => string): string {
  const code = normalizeTypedEnvelope(cause).code;
  return t(`step4.cutover.errors.${code && KNOWN_ERRORS.has(code) ? code : 'generic'}`);
}
