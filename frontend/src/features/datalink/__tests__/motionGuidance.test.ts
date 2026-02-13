import { describe, expect, it } from 'vitest';
import {
  MOTION_TOKENS,
  buildMotionReadabilityGate,
  resolveIntentMotionClass,
  resolveScrollBehavior,
} from '../motionGuidance';

describe('motionGuidance', () => {
  it('uses intent-only stage motion classes with reduced-motion fallback', () => {
    expect(resolveIntentMotionClass('grid', false)).toContain('animate-pulse');
    expect(resolveIntentMotionClass('grid', true)).not.toContain('animate-pulse');
    expect(resolveIntentMotionClass('commit', true)).toContain('ring-2');
  });

  it('switches scroll behavior for reduced-motion mode', () => {
    expect(resolveScrollBehavior(false)).toBe('smooth');
    expect(resolveScrollBehavior(true)).toBe('auto');
  });

  it('validates motion readability checklist and timing gate', () => {
    const gate = buildMotionReadabilityGate({
      stageHandoffMs: MOTION_TOKENS.stageHandoffMs,
      commitFeedbackMs: MOTION_TOKENS.commitFeedbackMs,
      hasReducedMotionFallback: true,
      intentOnlyAnimations: true,
    });
    expect(gate.pass).toBe(true);
    expect(gate.checklist).toHaveLength(4);
  });
});
