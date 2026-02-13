export const MOTION_TOKENS = {
  stageHandoffMs: 180,
  commitFeedbackMs: 220,
} as const;

export type IntentStage = 'idle' | 'grid' | 'commit';

export function resolveIntentMotionClass(stage: IntentStage, reducedMotion: boolean): string {
  if (stage === 'idle') return '';
  if (reducedMotion) {
    return stage === 'grid' ? 'ring-2 ring-sky-500/70 ring-offset-2 ring-offset-slate-900' : 'ring-2 ring-emerald-500/70';
  }
  return stage === 'grid'
    ? 'ring-2 ring-sky-500/70 ring-offset-2 ring-offset-slate-900 motion-safe:animate-pulse'
    : 'ring-2 ring-emerald-500/70 motion-safe:animate-pulse';
}

export function resolveScrollBehavior(reducedMotion: boolean): 'auto' | 'smooth' {
  return reducedMotion ? 'auto' : 'smooth';
}

export function buildMotionReadabilityGate(input: {
  stageHandoffMs: number;
  commitFeedbackMs: number;
  hasReducedMotionFallback: boolean;
  intentOnlyAnimations: boolean;
}): {
  pass: boolean;
  checklist: Array<{ id: string; pass: boolean }>;
} {
  const withinRange = (ms: number) => ms >= 150 && ms <= 300;
  const checklist = [
    { id: 'handoff_duration', pass: withinRange(input.stageHandoffMs) },
    { id: 'commit_duration', pass: withinRange(input.commitFeedbackMs) },
    { id: 'reduced_motion', pass: input.hasReducedMotionFallback },
    { id: 'intent_only', pass: input.intentOnlyAnimations },
  ];
  return {
    pass: checklist.every((item) => item.pass),
    checklist,
  };
}
