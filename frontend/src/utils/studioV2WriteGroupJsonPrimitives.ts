import { boundedString } from './safeJson';

export type JsonRecord = Record<string, unknown>;

export function isRecord(value: unknown): value is JsonRecord {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
}

export function readRequiredText(record: JsonRecord, key: string): string | null {
  return boundedString(record[key]) ?? null;
}

/** Reads a bounded string while allowing an explicit empty server value. */
export function readPresentText(record: JsonRecord, key: string): string | null | undefined {
  const value = record[key];
  if (value === undefined || value === null) return undefined;
  if (typeof value !== 'string' || value.length > 256) return null;
  return value;
}

export function readOptionalText(record: JsonRecord, key: string): string | undefined | null {
  const value = readPresentText(record, key);
  if (value === undefined) return undefined;
  if (value === null) return null;
  return boundedString(value) ?? null;
}

/**
 * Reads a tri-state flag: `undefined` when the key is absent or null (never
 * decided), the boolean when it is one, and `null` for anything else so the
 * caller fails closed.
 */
export function readOptionalBoolean(record: JsonRecord, key: string): boolean | undefined | null {
  const value = record[key];
  if (value === undefined || value === null) return undefined;
  return typeof value === 'boolean' ? value : null;
}

export function readRequiredInteger(record: JsonRecord, key: string, minimum = 0): number | null {
  const value = record[key];
  return typeof value === 'number' && Number.isInteger(value) && Number.isFinite(value) && value >= minimum
    ? value
    : null;
}

export function readOptionalInteger(record: JsonRecord, key: string, minimum = 0): number | undefined | null {
  const value = record[key];
  if (value === undefined || value === null) return undefined;
  return typeof value === 'number' && Number.isInteger(value) && Number.isFinite(value) && value >= minimum
    ? value
    : null;
}
