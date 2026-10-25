import type {
  WriteGroupIntervalUnit,
  WriteGroupProfileKind,
  WriteGroupRecordingProfile,
  WriteGroupSampleReduction,
} from '../types/studioV2WriteGroup';

interface IntervalBounds {
  min: number;
  max: number;
  default: number;
  /** The one public unit of this kind; the backend reads the interval in the same unit. */
  unit: WriteGroupIntervalUnit;
  secondsPerUnit: number;
}

/** Instantaneous 1-300 seconds (default 15); cumulative 1-600 minutes (default 10). */
export const RECORDING_PROFILE_BOUNDS: Readonly<Record<WriteGroupProfileKind, IntervalBounds>> = {
  instantaneous: { min: 1, max: 300, default: 15, unit: 'second', secondsPerUnit: 1 },
  cumulative: { min: 1, max: 600, default: 10, unit: 'minute', secondsPerUnit: 60 },
};

/**
 * Default of a NEW profile draft only. A saved or submitted profile that omits
 * or zeroes its interval is invalid; it is never completed from these values.
 */
export function defaultRecordingProfile(kind: WriteGroupProfileKind): WriteGroupRecordingProfile {
  const bounds = RECORDING_PROFILE_BOUNDS[kind];
  return { kind, interval: bounds.default, unit: bounds.unit };
}

function isProfileKind(value: unknown): value is WriteGroupProfileKind {
  return value === 'instantaneous' || value === 'cumulative';
}

function isSampleReduction(value: unknown): value is WriteGroupSampleReduction {
  return value === 'latest' || value === 'mean';
}

/** True only for a whole-number interval inside its kind's bounds and in that kind's own unit. */
export function isValidRecordingProfile(value: unknown): value is WriteGroupRecordingProfile {
  if (typeof value !== 'object' || value === null || Array.isArray(value)) return false;
  const { kind, interval, unit, sample_reduction } = value as Record<string, unknown>;
  if (!isProfileKind(kind)) return false;
  // An unknown reduction and mean on cumulative totals are both rejected, never coerced to latest.
  if (sample_reduction !== undefined && !isSampleReduction(sample_reduction)) return false;
  if (kind === 'cumulative' && sample_reduction === 'mean') return false;
  const bounds = RECORDING_PROFILE_BOUNDS[kind];
  return unit === bounds.unit && typeof interval === 'number' && Number.isInteger(interval) &&
    interval >= bounds.min && interval <= bounds.max;
}

/** Execution seconds of a valid profile; null rather than a guess for anything else. */
export function recordingProfileIntervalSeconds(profile: WriteGroupRecordingProfile): number | null {
  return isValidRecordingProfile(profile)
    ? profile.interval * RECORDING_PROFILE_BOUNDS[profile.kind].secondsPerUnit
    : null;
}

/**
 * Absent stays undefined, a valid profile is copied, and anything else is rejected with null.
 * A missing reduction is copied as missing: it decodes as latest and is never rewritten as a choice.
 */
export function parseRecordingProfile(value: unknown): WriteGroupRecordingProfile | undefined | null {
  if (value === undefined || value === null) return undefined;
  if (!isValidRecordingProfile(value)) return null;
  return {
    kind: value.kind,
    interval: value.interval,
    unit: value.unit,
    ...(value.sample_reduction !== undefined ? { sample_reduction: value.sample_reduction } : {}),
  };
}
