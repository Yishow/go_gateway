import { describe, expect, it } from 'vitest';
import { buildGlobalTagGuardrail } from '@/features/datalink/globalTagGuardrails';

describe('globalTagGuardrails', () => {
  it('requires second confirmation when affected mappings exist', () => {
    const result = buildGlobalTagGuardrail(3);
    expect(result.requiresConfirmation).toBe(true);
    expect(result.warningMessage).toContain('3');
  });

  it('does not require confirmation when no mapping is affected', () => {
    const result = buildGlobalTagGuardrail(0);
    expect(result.requiresConfirmation).toBe(false);
    expect(result.warningMessage).toContain('不會影響既有映射');
  });
});
