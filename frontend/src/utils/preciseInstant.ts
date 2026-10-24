/**
 * RFC 3339 instants with up to nanosecond precision. Date.parse keeps only milliseconds, which would make the
 * one-microsecond window a cutover can leave look empty, so whole seconds and the nanosecond part are kept apart.
 */
export interface PreciseInstant {
  /** Whole seconds since the Unix epoch. */
  seconds: number;
  /** 0 to 999999999. */
  nanos: number;
}

const RFC3339 = /^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})(?:\.(\d{1,9}))?(Z|[+-]\d{2}:\d{2})$/;

export function parsePreciseInstant(text: unknown): PreciseInstant | null {
  const match = typeof text === 'string' ? RFC3339.exec(text) : null;
  if (!match) return null;
  const whole = Date.parse(`${match[1]}${match[3]}`);
  return Number.isNaN(whole) ? null : { seconds: whole / 1000, nanos: Number((match[2] ?? '').padEnd(9, '0')) };
}

/** Negative, zero or positive like a comparator, exact to the nanosecond. */
export function comparePreciseInstants(a: PreciseInstant, b: PreciseInstant): number {
  return a.seconds - b.seconds || a.nanos - b.nanos;
}

/** `YYYY-MM-DD HH:MM:SS` in UTC, followed by the fraction only when there is one. */
export function formatPreciseUtc(text: string): string {
  const instant = parsePreciseInstant(text);
  if (!instant) return '';
  const whole = new Date(instant.seconds * 1000).toISOString().slice(0, 19).replace('T', ' ');
  const fraction = String(instant.nanos).padStart(9, '0').replace(/0+$/, '');
  return fraction ? `${whole}.${fraction}` : whole;
}
