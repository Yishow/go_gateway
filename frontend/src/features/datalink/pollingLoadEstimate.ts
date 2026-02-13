import type { Point, PollingGroup } from '../../types/datalink';

export interface PollingLoadEstimate {
  baselineReadsPerSec: number;
  deltaReadsPerSec: number;
  projectedReadsPerSec: number;
  assumedIntervalMs: number;
}

const FALLBACK_INTERVAL_MS = 1000;

function roundTo(value: number, digits = 2): number {
  const factor = 10 ** digits;
  return Math.round(value * factor) / factor;
}

export function estimatePollingLoadDelta(
  points: Point[],
  pollingGroups: PollingGroup[],
  plannedPointCount: number
): PollingLoadEstimate {
  const enabledGroups = pollingGroups.filter((group) => group.enabled && group.interval_ms > 0);
  const minimumIntervalMs = enabledGroups.reduce(
    (min, group) => Math.min(min, group.interval_ms),
    Number.POSITIVE_INFINITY
  );
  const fallbackIntervalMs =
    minimumIntervalMs === Number.POSITIVE_INFINITY ? FALLBACK_INTERVAL_MS : minimumIntervalMs;
  const intervalByGroup = new Map(enabledGroups.map((group) => [group.id, group.interval_ms]));

  const baselineReadsPerSecRaw = points
    .filter((point) => point.enabled && intervalByGroup.has(point.polling_group_id))
    .reduce((sum, point) => {
      const intervalMs = intervalByGroup.get(point.polling_group_id) || FALLBACK_INTERVAL_MS;
      return sum + 1000 / intervalMs;
    }, 0);

  const normalizedPlannedCount = Math.max(0, plannedPointCount);
  const deltaReadsPerSecRaw = normalizedPlannedCount * (1000 / fallbackIntervalMs);
  const projectedReadsPerSecRaw = baselineReadsPerSecRaw + deltaReadsPerSecRaw;

  return {
    baselineReadsPerSec: roundTo(baselineReadsPerSecRaw),
    deltaReadsPerSec: roundTo(deltaReadsPerSecRaw),
    projectedReadsPerSec: roundTo(projectedReadsPerSecRaw),
    assumedIntervalMs: fallbackIntervalMs,
  };
}
