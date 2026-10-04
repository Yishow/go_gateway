// Read-only SQLite observations shared by fresh destination runners.
import { sh } from './lib.mjs';

export function quoteSQLite(value) {
  return `'${String(value).replaceAll("'", "''")}'`;
}

export function readSQLiteRows(path, sql) {
  const raw = sh('sqlite3', ['-readonly', '-json', '-cmd', '.timeout 15000', `file:${path}?mode=ro`, sql]);
  return raw ? JSON.parse(raw) : [];
}

export async function waitForSQLiteRows({ path, sql, minimum = 3, timeoutMs = 270_000, pollMs = 2_000 } = {}) {
  const deadline = Date.now() + timeoutMs;
  const polls = [];
  let rows = [];
  let tSql;
  while (Date.now() < deadline) {
    const opened = process.hrtime.bigint().toString();
    rows = readSQLiteRows(path, sql);
    const completed = process.hrtime.bigint().toString();
    polls.push({ t_poll_open: opened, t_poll_complete: completed, row_count: rows.length, interval_ms: pollMs });
    if (rows.length > 0 && !tSql) tSql = completed;
    if (rows.length >= minimum) return { rows, polls, t_sql: tSql };
    await new Promise((resolve) => setTimeout(resolve, pollMs));
  }
  throw new Error(`SQLite observation did not reach ${minimum} rows; got ${rows.length}`);
}

export async function waitForClosedSQLiteBuckets({ path, groupId, groupRevision, timeoutMs = 270_000, pollMs = 2_000 } = {}) {
  const deadline = Date.now() + timeoutMs;
  const query = `SELECT group_id,group_revision,bucket_start,kind,reason,record_id,members
    FROM wg_delivery_buckets WHERE group_id = ${quoteSQLite(groupId)}
    AND group_revision = ${quoteSQLite(groupRevision)} AND kind = 'row'
    ORDER BY bucket_start, entity_key`;
  let rows = [];
  let tClosed;
  while (Date.now() < deadline) {
    rows = readSQLiteRows(path, query);
    if (rows.length > 0 && !tClosed) tClosed = process.hrtime.bigint().toString();
    if (rows.length >= 3) return { rows, query, t_closed: tClosed };
    await new Promise((resolve) => setTimeout(resolve, pollMs));
  }
  throw new Error(`gateway did not close three complete buckets; got ${rows.length}`);
}

export function readSQLiteDelivery({ path, groupId, groupRevision } = {}) {
  const scope = `group_id = ${quoteSQLite(groupId)} AND group_revision = ${quoteSQLite(groupRevision)}`;
  return {
    outbox: readSQLiteRows(path, `SELECT effect_key,record_id,group_id,group_revision,entity_key,bucket_start,table_name,state,payload_digest
      FROM wg_delivery_outbox WHERE ${scope} ORDER BY bucket_start,record_id`),
    receipts: readSQLiteRows(path, `SELECT effect_key,payload_digest,committed_at FROM wg_delivery_receipts
      WHERE effect_key IN (SELECT effect_key FROM wg_delivery_outbox WHERE ${scope})`),
    closedBuckets: readSQLiteRows(path, `SELECT group_revision,bucket_start,kind,record_id,members,entity_key
      FROM wg_delivery_buckets WHERE ${scope} AND kind = 'row' ORDER BY bucket_start,entity_key`),
  };
}
