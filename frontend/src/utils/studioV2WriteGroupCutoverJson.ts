import type {
  CutoverGroupView,
  CutoverIssueCode,
  CutoverMemberRef,
  CutoverRevisionRef,
  WriteGroupCutover,
  WriteGroupCutoverKind,
  WriteGroupCutoverRequest,
  WriteGroupCutoverReview,
} from '../types/studioV2WriteGroupCutover';
import { boundedString, MAX_SAFE_JSON_ARRAY_LENGTH } from './safeJson';
import { isRecord, type JsonRecord } from './studioV2WriteGroupJsonPrimitives';

const KINDS = new Set<WriteGroupCutoverKind>(['forward', 'rollback']);
const ISSUES = new Set<CutoverIssueCode>([
  'source_unavailable', 'source_not_applied', 'source_has_unapplied_draft', 'source_already_ended', 'source_is_long',
  'target_not_long', 'target_already_active', 'target_not_ready', 'target_member_not_in_source', 'targets_overlap',
  'unselected_sources', 'boundary_not_future', 'boundary_before_active_revision', 'no_targets', 'cutover_not_rollbackable',
  'remaining_draft_mismatch', 'remaining_draft_not_ready', 'remaining_draft_changes_policy',
  'remaining_draft_omits_required_column', 'restored_group_not_ready',
]);
const DIGEST = /^[a-f0-9]{64}$/;

/** An RFC 3339 instant, kept as the text the server sent so no precision is lost. */
function instant(value: unknown): string | null {
  const text = boundedString(value);
  return text && !Number.isNaN(Date.parse(text)) ? text : null;
}

function whole(value: unknown): number | null {
  return typeof value === 'number' && Number.isSafeInteger(value) && value >= 0 ? value : null;
}

function list<T>(value: unknown, parse: (item: unknown) => T | null): T[] | null {
  if (!Array.isArray(value) || value.length > MAX_SAFE_JSON_ARRAY_LENGTH) return null;
  const out: T[] = [];
  for (const item of value) {
    const parsed = parse(item);
    if (parsed === null) return null;
    out.push(parsed);
  }
  return out;
}

function parseMember(value: unknown): CutoverMemberRef | null {
  if (!isRecord(value)) return null;
  const device = boundedString(value.device_id);
  const point = boundedString(value.point_id);
  const tag = boundedString(value.tag_id);
  return device && point && tag ? { device_id: device, point_id: point, tag_id: tag } : null;
}

function parseRevisionRef(value: unknown): CutoverRevisionRef | null {
  if (!isRecord(value)) return null;
  const group = boundedString(value.group_id);
  const revision = boundedString(value.revision);
  return group && revision ? { group_id: group, revision } : null;
}

function optionalText(record: JsonRecord, key: string): string | undefined | null {
  const value = record[key];
  if (value === undefined || value === null || value === '') return undefined;
  return boundedString(value) ?? null;
}

function parseGroupView(value: unknown): CutoverGroupView | null {
  if (!isRecord(value)) return null;
  const id = boundedString(value.group_id);
  const revision = boundedString(value.revision);
  const layout = boundedString(value.layout);
  const connector = boundedString(value.connector_id);
  const interval = whole(value.interval_seconds);
  const members = whole(value.member_count);
  const mode = optionalText(value, 'storage_mode');
  const profile = optionalText(value, 'profile_kind');
  if (!id || !revision || !layout || !connector || interval === null || members === null || mode === null || profile === null ||
    typeof value.name !== 'string' || typeof value.table_schema !== 'string' || typeof value.table_name !== 'string' ||
    typeof value.schema_confirmed !== 'boolean') return null;
  return {
    group_id: id, name: value.name, revision, layout, interval_seconds: interval, connector_id: connector,
    table_schema: value.table_schema, table_name: value.table_name, member_count: members, schema_confirmed: value.schema_confirmed,
    ...(mode ? { storage_mode: mode } : {}),
    ...(profile ? { profile_kind: profile } : {}),
  };
}

function parseKind(value: unknown): WriteGroupCutoverKind | null {
  const kind = boundedString(value);
  return kind && KINDS.has(kind as WriteGroupCutoverKind) ? kind as WriteGroupCutoverKind : null;
}

/**
 * Parses a cutover review. A review that contradicts itself (applicable with
 * issues, or not applicable with no reason) or names an issue this client does
 * not know is rejected, never shown as applicable.
 */
export function parseWriteGroupCutoverReview(value: unknown): WriteGroupCutoverReview | null {
  if (!isRecord(value)) return null;
  const kind = parseKind(value.kind);
  const boundary = instant(value.boundary);
  const digest = boundedString(value.digest);
  const issues = list(value.issues, (item) => {
    const code = boundedString(item);
    return code && ISSUES.has(code as CutoverIssueCode) ? code as CutoverIssueCode : null;
  });
  const source = parseGroupView(value.source);
  const targets = list(value.targets, parseGroupView);
  const moved = list(value.moved, parseMember);
  const unselected = list(value.unselected, parseMember);
  const remaining = list(value.remaining, parseMember);
  if (!kind || !boundary || !digest || !DIGEST.test(digest) || typeof value.applicable !== 'boolean' || !issues ||
    value.applicable !== (issues.length === 0) || !source || !targets || !moved || !unselected || !remaining) return null;
  return { kind, boundary, digest, applicable: value.applicable, issues, source, targets, moved, unselected, remaining };
}

export function parseWriteGroupCutoverList(value: unknown): WriteGroupCutover[] | null {
  return list(value, parseWriteGroupCutover);
}

/** Parses the immutable record of an applied cutover or rollback. */
export function parseWriteGroupCutover(value: unknown): WriteGroupCutover | null {
  if (!isRecord(value)) return null;
  const id = boundedString(value.id);
  const workspace = boundedString(value.workspace_id);
  const kind = parseKind(value.kind);
  const source = boundedString(value.source_group_id);
  const boundary = instant(value.boundary);
  const digest = boundedString(value.review_digest);
  const created = instant(value.created_at);
  const prior = optionalText(value, 'prior_cutover_id');
  const ended = list(value.ended, parseRevisionRef);
  const started = list(value.started, parseRevisionRef);
  const moved = list(value.moved, parseMember);
  if (!id || !workspace || !kind || !source || !boundary || !digest || !DIGEST.test(digest) || !created || prior === null ||
    !ended || !started || !moved || (kind === 'rollback') !== (prior !== undefined)) return null;
  return {
    id, workspace_id: workspace, kind, source_group_id: source, boundary, review_digest: digest,
    ended, started, moved, created_at: created, ...(prior ? { prior_cutover_id: prior } : {}),
  };
}

/** Only the fields the server reads; the instant stays the text it was reviewed as. */
export function cutoverBody(request: WriteGroupCutoverRequest): WriteGroupCutoverRequest {
  return {
    kind: request.kind ?? 'forward',
    workspace_id: request.workspace_id,
    expected_workspace_revision: request.expected_workspace_revision,
    expected_source_revision: request.expected_source_revision,
    expected_connector_revision: request.expected_connector_revision,
    ...(request.targets ? { targets: request.targets.map((target) => ({ ...target })) } : {}),
    ...(request.cutover_id ? { cutover_id: request.cutover_id } : {}),
    ...(request.boundary ? { boundary: request.boundary } : {}),
    ...(request.review_digest ? { review_digest: request.review_digest } : {}),
  };
}
