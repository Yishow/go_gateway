import type { Point, SourceRuleRecord } from '../../../types/datalink';
import type { SourceRuleCandidateSnapshotView } from '../../../types/sourceRuleCandidates';

export function resolveActiveRuleId(
  rules: SourceRuleRecord[],
  preferredIds: Array<string | null | undefined>,
) {
  for (const candidateId of preferredIds) {
    if (candidateId && rules.some((rule) => rule.id === candidateId)) {
      return candidateId;
    }
  }

  return rules[0]?.id ?? null;
}

export function scopePointsToActiveRule(
  points: Point[],
  activeRuleId: string | null,
  candidateView?: SourceRuleCandidateSnapshotView | null,
) {
  if (!activeRuleId || !candidateView) {
    return points;
  }

  const activePointIds = new Set(candidateView.tags.candidates.map((candidate) => candidate.point_id));
  return points.filter((point) => activePointIds.has(point.id));
}
