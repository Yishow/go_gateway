import type { WriteGroup } from '../../../../types/studioV2WriteGroup';
import type { WriteGroupCutover, WriteGroupCutoverKind, WriteGroupCutoverRequest, WriteGroupCutoverReview } from '../../../../types/studioV2WriteGroupCutover';
import { groupProfileKind } from './circuitProfileBasic';

export interface CutoverScope {
  workspaceId: string;
  workspaceRevision: string;
  connectorRevision: string;
  source: WriteGroup;
  /** Forward: the prepared long targets. Rollback: none. */
  targets: WriteGroup[];
  /** The boundary exactly as the operator entered it; the server answer, not this text, is what apply sends. */
  boundary: string;
  kind: WriteGroupCutoverKind;
  cutoverId?: string;
}

export interface CutoverIntentRecord {
  workspace_id: string;
  source_id: string;
  kind: WriteGroupCutoverKind;
  boundary: string;
  review_digest: string;
  expected_workspace_revision: string;
  expected_source_revision: string;
  expected_connector_revision: string;
  request: WriteGroupCutoverRequest;
  created_at: string;
}

export function buildCutoverRequest(scope: CutoverScope): WriteGroupCutoverRequest {
  const sourceConnectorRevision = scope.source.destination?.connector_revision?.trim() || scope.connectorRevision;
  return {
    kind: scope.kind,
    workspace_id: scope.workspaceId,
    expected_workspace_revision: scope.workspaceRevision,
    expected_source_revision: scope.source.revision,
    expected_connector_revision: sourceConnectorRevision,
    ...(scope.kind === 'forward' ? { targets: scope.targets.map((group) => ({ group_id: group.id, expected_revision: group.revision })) } : {}),
    ...(scope.cutoverId ? { cutover_id: scope.cutoverId } : {}),
    ...(scope.boundary ? { boundary: scope.boundary } : {}),
  };
}

/** Everything a review depends on. Any difference means the review no longer describes what would be applied. */
export function cutoverFingerprint(scope: CutoverScope): string {
  const request = buildCutoverRequest(scope);
  return JSON.stringify([
    request.kind, request.workspace_id, request.expected_workspace_revision, scope.source.id, request.expected_source_revision,
    scope.connectorRevision, request.expected_connector_revision, [...(request.targets ?? [])].sort((a, b) => a.group_id.localeCompare(b.group_id)).map((t) => [t.group_id, t.expected_revision]),
    request.cutover_id ?? '', request.boundary ?? '',
  ]);
}

/** The apply request: the reviewed text of the boundary and the digest of exactly that review, nothing recomputed. */
export function applyRequestFor(scope: CutoverScope, review: WriteGroupCutoverReview): WriteGroupCutoverRequest {
  return { ...buildCutoverRequest(scope), boundary: review.boundary, review_digest: review.digest };
}

/** A source that can be switched: an applied managed wide or custom group that is not itself a long profile. */
export function isCutoverSource(group: WriteGroup): boolean {
  return group.status !== 'deleted' && group.status !== 'disabled' && Boolean(group.applied_revision) && groupProfileKind(group) === undefined &&
    group.layout_version !== 'managed_long_v1';
}

/** A prepared long target: a saved long profile group that has not been applied yet. */
export function isCutoverTarget(group: WriteGroup): boolean {
  return group.status !== 'deleted' && groupProfileKind(group) !== undefined && !group.applied_revision;
}

/** Finds the recorded cutover that a lost answer would have produced; matches on what was reviewed, never on recency. */
export function findRecordedCutover(list: readonly WriteGroupCutover[], request: WriteGroupCutoverRequest, sourceId: string): WriteGroupCutover | undefined {
  return list.find((item) => item.source_group_id === sourceId && item.kind === (request.kind ?? 'forward') &&
    item.review_digest === request.review_digest && item.boundary === request.boundary &&
    (!request.workspace_id || !item.workspace_id || item.workspace_id === request.workspace_id));
}

export function cutoverIntentKey(workspaceId: string, sourceId?: string): string {
  return sourceId
    ? `wbv2.cutover.${encodeURIComponent(workspaceId)}.${encodeURIComponent(sourceId)}`
    : `wbv2.cutover.${encodeURIComponent(workspaceId)}`;
}

function storage(): Storage | undefined {
  if (typeof window === 'undefined') return undefined;
  try {
    return window.sessionStorage;
  } catch {
    return undefined;
  }
}

export function saveCutoverIntent(workspaceId: string, record: CutoverIntentRecord): void {
  try {
    const s = storage();
    if (!s) return;
    const json = JSON.stringify(record);
    s.setItem(cutoverIntentKey(workspaceId), json);
    if (record.source_id) {
      s.setItem(cutoverIntentKey(workspaceId, record.source_id), json);
    }
  } catch {
    // Best-effort session persistence
  }
}

export function loadCutoverIntent(workspaceId: string, sourceId?: string): CutoverIntentRecord | undefined {
  try {
    const s = storage();
    if (!s) return undefined;
    const raw = (sourceId ? s.getItem(cutoverIntentKey(workspaceId, sourceId)) : null) || s.getItem(cutoverIntentKey(workspaceId));
    if (!raw) return undefined;
    const parsed = JSON.parse(raw) as Partial<CutoverIntentRecord>;
    if (!parsed || typeof parsed !== 'object' || parsed.workspace_id !== workspaceId || !parsed.source_id || !parsed.request) {
      return undefined;
    }
    const request = parsed.request;
    if (typeof parsed.source_id !== 'string' || (sourceId !== undefined && parsed.source_id !== sourceId) ||
      typeof request !== 'object' || Array.isArray(request) || request.workspace_id !== workspaceId ||
      !['forward', 'rollback'].includes(parsed.kind ?? '') || (request.kind ?? 'forward') !== parsed.kind) return undefined;
    for (const key of ['boundary', 'review_digest', 'expected_workspace_revision', 'expected_source_revision', 'expected_connector_revision'] as const) {
      if (typeof parsed[key] !== 'string' || parsed[key] === '' || request[key] !== parsed[key]) return undefined;
    }
    return parsed as CutoverIntentRecord;
  } catch {
    return undefined;
  }
}

export function clearCutoverIntent(workspaceId: string, sourceId?: string): void {
  try {
    const s = storage();
    if (!s) return;
    s.removeItem(cutoverIntentKey(workspaceId));
    if (sourceId) {
      s.removeItem(cutoverIntentKey(workspaceId, sourceId));
    }
  } catch {
    // Best effort
  }
}

/** Forward cutovers of the source that have not been rolled back, newest first. */
export function rollbackCandidates(list: readonly WriteGroupCutover[]): WriteGroupCutover[] {
  const rolledBack = new Set(list.filter((item) => item.kind === 'rollback' && item.prior_cutover_id).map((item) => item.prior_cutover_id));
  return list.filter((item) => item.kind === 'forward' && !rolledBack.has(item.id)).sort((a, b) => b.created_at.localeCompare(a.created_at));
}

/** `datetime-local` text to the RFC 3339 instant the server reviews; second precision, in the browser's time zone. */
export function instantFromLocal(local: string): string {
  const date = new Date(local);
  return Number.isNaN(date.getTime()) ? '' : date.toISOString();
}

/** True when the server-reviewed boundary has not been reached yet, so the switch is scheduled and nothing has changed. */
export function boundaryIsScheduled(boundary: string, now: number = Date.now()): boolean {
  const at = Date.parse(boundary);
  return !Number.isNaN(at) && at > now;
}
