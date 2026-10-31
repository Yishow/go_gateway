import type { WriteGroup, WriteGroupProfileKind } from '../../../../types/studioV2WriteGroup';
import { normalizeTypedEnvelope } from '../../../../utils/safeJson';
import { groupCircuitId, groupProfileKind } from './circuitProfileBasic';
import { existingProfileUpdateGroup, newProfileDraft, profileMatchesIntent } from './circuitProfileBatchIntent';

export type BatchProfileStatus = 'saved' | 'failed' | 'unknown' | 'not-attempted' | 'unchanged';
export type BatchTargetStatus = BatchProfileStatus | 'partial';

export interface BatchProfileResult {
  kind: WriteGroupProfileKind;
  status: BatchProfileStatus;
  /** The original rejection, kept so operator copy can use its typed code; never rendered raw. */
  cause?: unknown;
  group?: WriteGroup;
}

export type BatchProfileResults = Partial<Record<WriteGroupProfileKind, BatchProfileResult>>;

export interface BatchTargetExecutionResult {
  circuit_id: string;
  circuit_name: string;
  /** Derived from `profiles`; it never replaces them. */
  status: BatchTargetStatus;
  profiles: BatchProfileResults;
}

export const BATCH_PROFILE_ORDER: readonly WriteGroupProfileKind[] = ['instantaneous', 'cumulative'];

/**
 * Only an explicit application rejection proves the write was not committed: a validated `success:false`
 * envelope, or a 4xx that carries one of this backend's typed error codes. A lost reply, timeout, 5xx, a
 * bare status or a proxy page (HTML) says nothing about the server state.
 */
export function classifyBatchMutationError(error: unknown): 'failed' | 'unknown' {
  const record = typeof error === 'object' && error !== null ? error as Record<string, unknown> : {};
  if (record.outcome === 'failed') return 'failed';
  if (record.outcome === 'unconfirmed') return 'unknown';
  const response = typeof record.response === 'object' && record.response !== null ? record.response as Record<string, unknown> : {};
  const status = typeof response.status === 'number' ? response.status : typeof record.status === 'number' ? record.status : undefined;
  const rejected = status !== undefined && status >= 400 && status < 500 && status !== 408;
  return rejected && normalizeTypedEnvelope(error).code !== undefined ? 'failed' : 'unknown';
}

export function summarizeBatchTarget(profiles: BatchProfileResults): BatchTargetStatus {
  const statuses = BATCH_PROFILE_ORDER.flatMap((kind) => (profiles[kind] ? [profiles[kind].status] : []));
  if (statuses.length === 0 || statuses.every((status) => status === 'not-attempted')) return 'not-attempted';
  if (statuses.includes('unknown')) return 'unknown';
  if (statuses.every((status) => status === 'saved' || status === 'unchanged')) return statuses.includes('saved') ? 'saved' : 'unchanged';
  return statuses.includes('saved') ? 'partial' : 'failed';
}

export function groupProfileInterval(group: WriteGroup, kind: WriteGroupProfileKind): number | undefined {
  if (group.recording_profile?.kind === 'mixed') return undefined;
  const interval = group.recording_profile?.interval;
  if (interval !== undefined) return interval;
  if (group.row_policy.interval_seconds === undefined) return undefined;
  return (kind === 'cumulative' ? Math.round(group.row_policy.interval_seconds / 60) : group.row_policy.interval_seconds);
}

interface ReviewedTarget {
  circuit_id: string;
  circuit_name?: string;
  instantaneousGroup?: WriteGroup;
  cumulativeGroup?: WriteGroup;
  anyGroup?: WriteGroup;
}

/**
 * Upgrades an unknown profile to saved only when the reloaded server state carries the configuration that
 * was reviewed (the same intent the executor sent), not merely the same period: an existing writer must have
 * a new revision and match the reviewed update, a new writer must match the draft that would have been created.
 * Anything else stays unknown, so a changed source set, destination or policy is never reported as this save.
 */
export function reconcileBatchResults(
  targets: ReviewedTarget[],
  intervals: Record<WriteGroupProfileKind, number>,
  results: Record<string, BatchTargetExecutionResult>,
  liveGroups: WriteGroup[],
  workspaceId?: string,
): { results: Record<string, BatchTargetExecutionResult>; proven: WriteGroup[] } {
  const proven: WriteGroup[] = [];
  const next: Record<string, BatchTargetExecutionResult> = {};
  for (const [circuitId, result] of Object.entries(results)) {
    const target = targets.find((item) => item.circuit_id === circuitId);
    const profiles: BatchProfileResults = { ...result.profiles };
    for (const kind of BATCH_PROFILE_ORDER) {
      if (!target || profiles[kind]?.status !== 'unknown') continue;
      const reviewed = kind === 'instantaneous' ? target.instantaneousGroup : target.cumulativeGroup;
      let live: WriteGroup | undefined;
      if (reviewed) {
        const intent = existingProfileUpdateGroup(reviewed, kind, intervals[kind]);
        live = liveGroups.find((group) => group.id === reviewed.id && group.revision !== reviewed.revision && profileMatchesIntent(group, intent));
      } else if (target.anyGroup && workspaceId) {
        const intent = newProfileDraft(workspaceId, { circuit_id: circuitId, circuit_name: target.circuit_name ?? result.circuit_name, anyGroup: target.anyGroup }, kind, intervals[kind]);
        live = liveGroups.find((group) => groupCircuitId(group) === circuitId && groupProfileKind(group) === kind && profileMatchesIntent(group, intent));
      }
      if (live) {
        profiles[kind] = { kind, status: 'saved', group: live };
        proven.push(live);
      }
    }
    next[circuitId] = { ...result, profiles, status: summarizeBatchTarget(profiles) };
  }
  return { results: next, proven };
}
