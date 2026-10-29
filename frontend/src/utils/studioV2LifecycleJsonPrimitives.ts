import { MAX_SAFE_JSON_ARRAY_LENGTH, boundedString } from './safeJson';
import { isRecord } from './studioV2WriteGroupJsonPrimitives';
export { isRecord };
export function text(value: unknown, empty = false): string | null {
  if (empty && value === '') return '';
  const parsed = boundedString(value);
  return parsed && parsed === value ? parsed : null;
}
export function code(value: unknown, empty = false): string | null {
  const parsed = text(value, empty);
  return parsed !== null && (parsed === '' || /^[a-z][a-z0-9_-]{0,127}$/.test(parsed)) ? parsed : null;
}
export function digest(value: unknown): string | null {
  return typeof value === 'string' && /^[a-f0-9]{64}$/.test(value) ? value : null;
}
export function whole(value: unknown, min = 0): number | null {
  return typeof value === 'number' && Number.isSafeInteger(value) && value >= min ? value : null;
}
export function instant(value: unknown): string | null {
  const parsed = text(value);
  return parsed && /^\d{4}-\d\d-\d\dT\d\d:\d\d:\d\d(?:\.\d+)?(?:Z|[+-]\d\d:\d\d)$/.test(parsed) &&
    !Number.isNaN(Date.parse(parsed)) ? parsed : null;
}
export function list<T>(value: unknown, parse: (item: unknown) => T | null, nullable = false): T[] | null {
  if (value === null && nullable) return [];
  if (!Array.isArray(value) || value.length > MAX_SAFE_JSON_ARRAY_LENGTH) return null;
  const result: T[] = [];
  for (const item of value) { const parsed = parse(item); if (parsed === null) return null; result.push(parsed); }
  return result;
}
export function texts<K extends string>(value: Record<string, unknown>, keys: readonly K[]): Record<K, string> | null {
  const result = {} as Record<K, string>;
  for (const key of keys) { const parsed = text(value[key]); if (parsed === null) return null; result[key] = parsed; }
  return result;
}
export function optionalText(value: unknown): string | undefined | null {
  return value === undefined ? undefined : text(value);
}
export function optionalCode(value: unknown): string | undefined | null {
  return value === undefined ? undefined : code(value);
}
