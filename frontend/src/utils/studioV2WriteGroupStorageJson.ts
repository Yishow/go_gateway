import type {
  WriteGroupStorage,
  WriteGroupStorageEvent,
  WriteGroupStorageEventKind,
} from '../types/studioV2WriteGroupDelivery';
import { boundedString } from './safeJson';
import { isRecord, type JsonRecord, readOptionalBoolean } from './studioV2WriteGroupJsonPrimitives';

const EVENT_KINDS = new Set<WriteGroupStorageEventKind>([
  'confirmed', 'recreated_by_confirmation', 'reconstructed', 'repair_blocked',
]);
/** The backend keeps at most this many diagnostics per group. */
const MAX_STORAGE_EVENTS = 20;

function wholeNumber(value: unknown, minimum: number): number | null {
  return typeof value === 'number' && Number.isSafeInteger(value) && value >= minimum ? value : null;
}

function isoTime(value: unknown): string | null {
  const text = boundedString(value);
  return text && !Number.isNaN(Date.parse(text)) ? text : null;
}

/** Optional value-free text: absent stays absent, anything but short text fails. */
function optionalText(record: JsonRecord, key: string): string | undefined | null {
  const value = record[key];
  if (value === undefined || value === null || value === '') return undefined;
  return boundedString(value) ?? null;
}

function parseEvent(value: unknown): WriteGroupStorageEvent | null {
  if (!isRecord(value)) return null;
  const kind = boundedString(value.kind);
  const generation = wholeNumber(value.generation, 0);
  const occurrences = wholeNumber(value.occurrences, 1);
  const occurredAt = isoTime(value.occurred_at);
  const groupId = optionalText(value, 'group_id');
  const code = optionalText(value, 'code');
  const reason = optionalText(value, 'reason');
  const operationId = optionalText(value, 'operation_id');
  if (!kind || !EVENT_KINDS.has(kind as WriteGroupStorageEventKind) || generation === null || occurrences === null ||
    !occurredAt || groupId === null || code === null || reason === null || operationId === null) return null;
  return {
    generation, kind: kind as WriteGroupStorageEventKind, occurrences, occurred_at: occurredAt,
    ...(groupId ? { group_id: groupId } : {}),
    ...(code ? { code } : {}),
    ...(reason ? { reason } : {}),
    ...(operationId ? { operation_id: operationId } : {}),
  };
}

/**
 * Parses the managed-storage truth of one long group. Absent (or null) means the
 * group has no managed long storage, which is not an error. A view that
 * contradicts itself is rejected, never shown: a rebuilt table must come with
 * the disclosure that earlier history may be gone.
 */
export function parseWriteGroupStorage(value: unknown): WriteGroupStorage | undefined | null {
  if (value === undefined || value === null) return undefined;
  if (!isRecord(value)) return null;
  const generation = wholeNumber(value.generation, 1);
  const reconstructions = wholeNumber(value.reconstructions, 0);
  const recreateMissing = readOptionalBoolean(value, 'recreate_missing');
  const repairReason = optionalText(value, 'repair_reason');
  const rebuiltAt = value.last_reconstructed_at === undefined ? undefined : isoTime(value.last_reconstructed_at);
  const rawEvents = value.events;
  if (generation === null || reconstructions === null || typeof value.history_may_be_lost !== 'boolean' ||
    recreateMissing === null || repairReason === null || rebuiltAt === null ||
    (value.repair_state !== undefined && value.repair_state !== '' && value.repair_state !== 'blocked') ||
    !Array.isArray(rawEvents) || rawEvents.length > MAX_STORAGE_EVENTS) return null;
  if ((reconstructions > 0 || value.history_may_be_lost) && generation < 2) return null;
  if (reconstructions > 0 && !value.history_may_be_lost) return null;
  const events = rawEvents.map(parseEvent);
  if (events.some((event) => event === null)) return null;
  return {
    generation, history_may_be_lost: value.history_may_be_lost, reconstructions,
    ...(rebuiltAt ? { last_reconstructed_at: rebuiltAt } : {}),
    ...(value.repair_state === 'blocked' ? { repair_state: 'blocked' as const } : {}),
    ...(repairReason ? { repair_reason: repairReason } : {}),
    ...(recreateMissing !== undefined ? { recreate_missing: recreateMissing } : {}),
    events: events as WriteGroupStorageEvent[],
  };
}
