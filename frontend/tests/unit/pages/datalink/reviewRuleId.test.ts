import { describe, expect, it } from 'vitest';
import { resolvePersistedReviewRuleId } from '../../../../src/pages/datalink/workbench/reviewRuleId';

describe('resolvePersistedReviewRuleId', () => {
  it('returns the first persisted preferred rule id', () => {
    expect(
      resolvePersistedReviewRuleId(['persisted-rule-1'], [
        { id: 'persisted-rule-1' },
        { id: 'persisted-rule-2' },
      ]),
    ).toBe('persisted-rule-1');
  });

  it('falls back when the first preferred rule id is draft-only', () => {
    expect(
      resolvePersistedReviewRuleId(['rule-1', 'persisted-rule-2'], [
        { id: 'persisted-rule-1' },
        { id: 'persisted-rule-2' },
      ]),
    ).toBe('persisted-rule-2');
  });

  it('returns null when no preferred rule id is present', () => {
    expect(resolvePersistedReviewRuleId([], [{ id: 'persisted-rule-1' }])).toBeNull();
  });
});
