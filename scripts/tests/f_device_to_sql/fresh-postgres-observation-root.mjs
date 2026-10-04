import assert from 'node:assert/strict';
import { sh, validateOwnedPostgresTarget } from './lib.mjs';

export function quotePGIdentifier(value) {
  assert.match(String(value), /^[A-Za-z_][A-Za-z0-9_]*$/, 'unsafe PostgreSQL identifier');
  return `"${value}"`;
}

export function ownedPGTarget(schema) {
  assert.match(schema, /^gw_f_[A-Za-z0-9_]+$/, 'PostgreSQL target must use a fresh owned schema');
  const dsn = Object.fromEntries(String(process.env.POSTGRES_DSN ?? '').split(/\s+/).filter(Boolean).map((part) => {
    const index = part.indexOf('=');
    return [part.slice(0, index), part.slice(index + 1)];
  }));
  const image = sh('docker', ['inspect', '--format', '{{.Config.Image}}', 'gw-wg-pg-test']);
  validateOwnedPostgresTarget(dsn, image);
  return { host: dsn.host, port: Number(dsn.port), database: dsn.dbname, username: dsn.user,
    password: dsn.password, schema, table: 'f_connector_placeholder' };
}

// Environment provisioning may CREATE/DROP only the caller's fresh schema.
// All subsequent recording observations use SELECT and never prepare target tables or values.
export function pgCommand(sql) {
  return sh('docker', ['exec', '-e', 'PGOPTIONS=-c statement_timeout=10000', '-i', 'gw-wg-pg-test',
    'psql', '-v', 'ON_ERROR_STOP=1', '-U', 'postgres', '-d', 'gwtest', '-At', '-c', sql], { timeout: 15_000 });
}

export function pgRows(sql) {
  assert.match(sql, /^SELECT\b/i, 'PostgreSQL row observer is SELECT-only');
  const raw = pgCommand(`SELECT COALESCE(json_agg(observation), '[]'::json) FROM (${sql}) observation`);
  return JSON.parse(raw);
}

export function pgRecordingRelations(schema) {
  quotePGIdentifier(schema);
  return pgRows(`SELECT c.relname AS name, c.relkind FROM pg_class c JOIN pg_namespace n ON n.oid=c.relnamespace
    WHERE n.nspname='${schema}' AND c.relkind IN ('r','p','v','m','f') ORDER BY c.relname`);
}

export function pgColumnMetadata(schema, table) {
  quotePGIdentifier(schema);
  quotePGIdentifier(table);
  return pgRows(`SELECT column_name AS name,data_type,numeric_precision,numeric_scale,is_nullable
    FROM information_schema.columns WHERE table_schema='${schema}' AND table_name='${table}' ORDER BY ordinal_position`);
}

export function pgRecordingRows(schema, table, columns) {
  const selection = columns.map(({ name }) => {
    const identifier = quotePGIdentifier(name);
    // Cast NUMERIC before JSON parsing so uint64 values never pass through a JS Number.
    const expression = name === 'bucket_start'
      ? `to_char(${identifier}::timestamptz AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')`
      : `${identifier}::text`;
    return `${expression} AS ${identifier}`;
  }).join(',');
  return pgRows(`SELECT ${selection} FROM ${quotePGIdentifier(schema)}.${quotePGIdentifier(table)} ORDER BY bucket_start,record_id`);
}

export function localDurableRows(path, sql) {
  assert.match(sql, /^SELECT\b/i, 'local durable observer is SELECT-only');
  const raw = sh('sqlite3', ['-readonly', '-json', '-cmd', '.timeout 10000', `file:${path}?mode=ro`, sql], { timeout: 15_000 });
  return raw ? JSON.parse(raw) : [];
}
