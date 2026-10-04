export const FRESH_FIRST_SQL_MAX_MS = 300_000;

const POLLING_LIMITATION =
  'Local closure SELECT precedes target SELECT in each loop; observation delay includes up to one poll interval and actual SQL command durations.';

function parseMonotonicTimestamp(name, value) {
  if (value === undefined || value === null || value === '') {
    throw new TypeError(`${name} is required`);
  }

  let timestamp;
  try {
    timestamp = typeof value === 'bigint' ? value : BigInt(value);
  } catch (error) {
    throw new TypeError(`${name} must be a monotonic integer timestamp`, { cause: error });
  }
  if (timestamp < 0n) {
    throw new RangeError(`${name} must not be negative`);
  }
  return timestamp;
}

function milliseconds(nanoseconds) {
  return Number(nanoseconds) / 1e6;
}

export function assertFreshTimingWitness({ t_open, t_start, t_closed, t_sql, polls } = {}) {
  if (!Array.isArray(polls) || polls.length === 0) {
    throw new TypeError('polling observations are required');
  }

  const open = parseMonotonicTimestamp('t_open', t_open);
  const start = parseMonotonicTimestamp('t_start', t_start);
  const closed = parseMonotonicTimestamp('t_closed', t_closed);
  const sql = parseMonotonicTimestamp('t_sql', t_sql);

  if (open > start || start > closed || closed > sql) {
    throw new RangeError('monotonic observation timestamps must satisfy t_open <= t_start <= t_closed <= t_sql');
  }

  const endToEndNanoseconds = sql - open;
  if (endToEndNanoseconds > BigInt(FRESH_FIRST_SQL_MAX_MS) * 1_000_000n) {
    throw new RangeError('controlled first independent SQL must be at most 300 seconds');
  }

  return {
    t_open: open.toString(),
    t_start: start.toString(),
    t_closed: closed.toString(),
    t_sql: sql.toString(),
    setup_ms: milliseconds(start - open),
    system_wait_ms: milliseconds(closed - start),
    delivery_observed_ms: milliseconds(sql - closed),
    end_to_end_ms: milliseconds(endToEndNanoseconds),
    polls,
    polling_limitation: POLLING_LIMITATION,
  };
}
