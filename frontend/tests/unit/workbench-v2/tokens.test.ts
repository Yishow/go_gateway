import { describe, it, expect } from 'vitest';
import { COLORS, FONTS, RADII, ANIMATIONS } from '../../../src/features/datalink/workbench-v2/tokens';

describe('Workbench V2 Tokens', () => {
  it('should have correct core colors', () => {
    expect(COLORS.bg).toBe('#0b1220');
    expect(COLORS.success).toBe('#10b981');
    expect(COLORS.primary).toBe('#3b82f6');
  });

  it('should have correct font definitions', () => {
    expect(FONTS.sans).toContain('Inter');
    expect(FONTS.mono).toContain('JetBrains Mono');
  });

  it('should have correct radii and animations', () => {
    expect(RADII.full).toBe('9999px');
    expect(ANIMATIONS.pulseDot).toBe('pulseDot 1.6s ease-in-out infinite');
  });
});
