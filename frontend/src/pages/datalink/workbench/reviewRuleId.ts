import type { SourceRuleRecord } from '../../../types/datalink';

export function resolvePersistedReviewRuleId(
  preferredRuleIds: ReadonlyArray<string | null | undefined>,
  persistedRules: ReadonlyArray<Pick<SourceRuleRecord, 'id'>>,
): string | null {
  for (const ruleId of preferredRuleIds) {
    if (!ruleId) {
      continue;
    }

    if (persistedRules.some((rule) => rule.id === ruleId)) {
      return ruleId;
    }
  }

  return null;
}
