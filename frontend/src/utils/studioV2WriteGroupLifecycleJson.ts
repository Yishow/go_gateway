import type { WriteGroupStatus } from '../types/studioV2WriteGroup';
import type { WriteGroupLifecycle, WriteGroupLifecycleAction, WriteGroupLifecycleDestination, WriteGroupLifecycleGeneration, WriteGroupLifecycleHistory, WriteGroupLifecyclePending, WriteGroupLifecycleRevision, WriteGroupSuccessorSeed } from '../types/studioV2WriteGroupLifecycle';
import { parseWriteGroupSuccessorCandidate } from './studioV2WriteGroupDefinitionJson';
import { isRecord, text, texts, code, digest, whole, instant, list, optionalText } from './studioV2LifecycleJsonPrimitives';

function status(value: unknown): WriteGroupStatus | null {
  return value === 'draft' || value === 'ready' || value === 'running' || value === 'disabled' || value === 'deleted' ? value : null;
}
function action(value: unknown): WriteGroupLifecycleAction | null {
  if (!isRecord(value)) return null;
  const name = value.action;
  const reason = code(value.reason_code, true), next = code(value.next_action, true);
  if ((name !== 'edit' && name !== 'stop' && name !== 'remove' && name !== 'recreate' && name !== 'transition') ||
    typeof value.allowed !== 'boolean' || typeof value.requires_review !== 'boolean' || reason === null || next === null ||
    (!value.allowed && !reason) || (value.allowed && reason !== '') ||
    ((name === 'recreate' || name === 'transition') && !value.requires_review)) return null;
  return { action: name, allowed: value.allowed, reason_code: reason, next_action: next, requires_review: value.requires_review };
}
function destination(value: unknown): WriteGroupLifecycleDestination | null {
  return isRecord(value) ? texts(value, ['connector_id', 'connector_revision', 'destination_scope']) : null;
}
function groupKey(value: unknown) {
  // LifecycleHandoff.MarshalJSON publishes explicitly scoped, lowercase identities.
  return isRecord(value) ? texts(value, ['workspace_id', 'group_id', 'group_revision']) : null;
}
export function parseWriteGroupLifecyclePending(value: unknown): WriteGroupLifecyclePending | null {
  if (!isRecord(value) || !isRecord(value.outbox)) return null;
  const samples = whole(value.samples);
  const keys = ['outstanding', 'total', 'queued', 'retrying', 'blocked', 'quarantined', 'unknown', 'sql_committed', 'skipped'] as const;
  const outbox = {} as WriteGroupLifecyclePending['outbox'];
  for (const key of keys) { const n = whole(value.outbox[key]); if (n === null) return null; outbox[key] = n; }
  if (samples === null || keys.some(key => outbox[key] > outbox.total) ||
    outbox.outstanding !== outbox.queued + outbox.retrying + outbox.blocked + outbox.quarantined + outbox.unknown ||
    outbox.total !== outbox.outstanding + outbox.sql_committed + outbox.skipped) return null;
  const handoffs = list(value.handoffs, item => {
    if (!isRecord(item)) return null;
    const origin = groupKey(item.origin), source = groupKey(item.source), target = groupKey(item.target);
    const count = whole(item.samples), state = code(item.state);
    return origin && source && target && count !== null && state && ['pending', 'accepted', 'finalized', 'quarantined'].includes(state)
      ? { origin, source, target, state, samples: count } : null;
  }, true);
  return handoffs ? { samples, outbox, handoffs } : null;
}
function revision(value: unknown): WriteGroupLifecycleRevision | null {
  if (!isRecord(value)) return null;
  const rev = text(value.group_revision), effective = instant(value.effective_at);
  const ends = value.ends_at === undefined ? undefined : instant(value.ends_at);
  const dest = destination(value.destination), pending = parseWriteGroupLifecyclePending(value.pending);
  const sources = list(value.sources, item => isRecord(item) ? texts(item, ['device_id', 'point_id', 'tag_id']) : null);
  if (!rev || !effective || ends === null || (ends && Date.parse(ends) < Date.parse(effective)) || !dest || !pending || !sources) return null;
  return { group_revision: rev, effective_at: effective, ...(ends ? { ends_at: ends } : {}), sources, destination: dest, pending };
}
function generation(value: unknown): WriteGroupLifecycleGeneration | null {
  if (!isRecord(value)) return null;
  const ids = texts(value, ['group_id', 'group_revision']);
  const applied = text(value.applied_revision, true), state = status(value.status), dest = destination(value.saved_destination);
  const predecessor = optionalText(value.predecessor_group_id), number = value.generation === undefined ? undefined : whole(value.generation, 1);
  const revisions = list(value.revisions, revision);
  if (!ids || applied === null || !state || !dest || predecessor === null || number === null || !revisions ||
    typeof value.current !== 'boolean' || typeof value.reserved !== 'boolean' ||
    (value.reserved && (value.current || state !== 'draft' || applied !== ''))) return null;
  return { ...ids, applied_revision: applied, status: state, saved_destination: dest, revisions,
    current: value.current, reserved: value.reserved, ...(predecessor ? { predecessor_group_id: predecessor } : {}),
    ...(number !== undefined ? { generation: number } : {}) };
}
export function parseWriteGroupLifecycleHistory(value: unknown): WriteGroupLifecycleHistory | null {
  if (!isRecord(value) || typeof value.canonical !== 'boolean') return null;
  const current = text(value.current_group_id), lineage = optionalText(value.lineage_id), reserved = optionalText(value.reserved_group_id);
  const rev = value.lineage_revision === undefined ? undefined : whole(value.lineage_revision, 1);
  const generations = list(value.generations, generation);
  if (!current || lineage === null || reserved === null || rev === null || !generations || generations.length === 0 ||
    (value.canonical ? !lineage || rev === undefined : lineage !== undefined || rev !== undefined || reserved !== undefined) ||
    new Set(generations.map(g => g.group_id)).size !== generations.length ||
    generations.filter(g => g.current && g.group_id === current).length !== 1 ||
    generations.some(g => g.current !== (g.group_id === current) || g.reserved !== (g.group_id === reserved) ||
      (value.canonical ? g.generation === undefined : g.generation !== undefined || g.predecessor_group_id !== undefined)) ||
    (reserved && !generations.some(g => g.group_id === reserved))) return null;
  return { canonical: value.canonical, current_group_id: current, generations,
    ...(lineage ? { lineage_id: lineage } : {}), ...(rev !== undefined ? { lineage_revision: rev } : {}), ...(reserved ? { reserved_group_id: reserved } : {}) };
}
function seed(value: unknown): WriteGroupSuccessorSeed | null {
  if (!isRecord(value) || (value.basis_kind !== 'applied' && value.basis_kind !== 'saved')) return null;
  const revision = text(value.basis_revision), hash = digest(value.basis_digest), group = parseWriteGroupSuccessorCandidate(value.group);
  return revision && hash && group ? { basis_kind: value.basis_kind, basis_revision: revision, basis_digest: hash, group } : null;
}
export function parseWriteGroupLifecycle(value: unknown): WriteGroupLifecycle | null {
  if (!isRecord(value)) return null;
  const ids = texts(value, ['workspace_id', 'workspace_revision', 'group_id', 'group_revision', 'connector_id', 'connector_revision']);
  const applied = text(value.applied_revision, true), state = status(value.status), history = parseWriteGroupLifecycleHistory(value.history);
  const actions = list(value.actions, action), successor = value.successor_seed === undefined ? undefined : seed(value.successor_seed);
  if (!ids || applied === null || !state || !history || !actions || successor === null || actions.length !== 5 ||
    new Set(actions.map(a => a.action)).size !== 5 || !history.generations.some(g => g.group_id === ids.group_id && g.group_revision === ids.group_revision && g.status === state && g.applied_revision === applied) ||
    (successor && (successor.group.workspace_id !== ids.workspace_id || !history.canonical ||
      successor.basis_revision !== (successor.basis_kind === 'applied' ? applied : ids.group_revision))) ||
    (actions.some(a => a.allowed && (a.action === 'recreate' || a.action === 'transition')) && !successor)) return null;
  return { ...ids, applied_revision: applied, status: state, actions, history, ...(successor ? { successor_seed: successor } : {}) };
}
