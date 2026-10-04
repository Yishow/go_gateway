import { lstatSync, statSync } from 'node:fs';
import { join, relative, resolve } from 'node:path';
import { sh } from './lib.mjs';

const OWNED_WORK_PATTERN = /^\/tmp\/gw-f-[A-Za-z0-9_-]+$/;
const OWNED_SCHEMA_PATTERN = /^gw_f_[A-Za-z0-9_-]+$/;
const IDENTIFIER_PATTERN = /^[A-Za-z_][A-Za-z0-9_]*$/;

function assertIdentifier(value, label) {
  if (!IDENTIFIER_PATTERN.test(value ?? '')) throw new Error(`${label} must be a safe SQL identifier`);
  return value;
}

function assertOwnedSchema(schema) {
  if (!OWNED_SCHEMA_PATTERN.test(schema ?? '')) throw new Error('PostgreSQL schema must use the owned gw_f_* namespace');
  return schema;
}

function assertOwnedSQLitePath(path, work) {
  const ownedWork = resolve(work ?? '');
  const target = resolve(path ?? '');
  if (!OWNED_WORK_PATTERN.test(ownedWork) || !target.startsWith(`${ownedWork}/`)) {
    throw new Error('SQLite destination must stay below the owned /tmp/gw-f-* run namespace');
  }
  let current = ownedWork;
  for (const segment of ['', ...relative(ownedWork, target).split('/')]) {
    if (segment) current = join(current, segment);
    try {
      if (lstatSync(current).isSymbolicLink()) throw new Error('SQLite owned namespace must not contain a symlink alias');
    } catch (error) {
      if (error?.code !== 'ENOENT') throw error;
      break;
    }
  }
  return target;
}

function readOnlySQLiteURI(path) {
  return `file:${path}?mode=ro`;
}

function sqliteTables(path, sqliteBinary) {
  const sql = "SELECT name FROM sqlite_master WHERE type = 'table' ORDER BY name;";
  try {
    return sh(sqliteBinary, ['-readonly', readOnlySQLiteURI(path), sql])
      .split(/\r?\n/).map((value) => value.trim()).filter(Boolean);
  } catch (error) {
    throw new Error(`SQLite read-only observation failed: ${error.message.split('\n')[0]}`);
  }
}

/** Stat a missing file without opening it; inspect an existing file only through SQLite read-only URI. */
export function preflightSQLiteDestination({ path, work, table = 'readings', sqliteBinary = 'sqlite3' } = {}) {
  const target = assertOwnedSQLitePath(path, work);
  const recordingTable = assertIdentifier(table, 'recording table');
  let fileStat;
  try {
    fileStat = statSync(target);
  } catch (error) {
    if (error?.code !== 'ENOENT') throw new Error(`SQLite destination stat failed: ${error.message.split('\n')[0]}`);
    return {
      kind: 'sqlite', path: target, status: 'missing', exists: false, stat: 'absent',
      recording_table: recordingTable, recording_table_present: false,
    };
  }
  if (!fileStat.isFile()) throw new Error('SQLite destination must be a regular file');
  const tables = sqliteTables(target, sqliteBinary);
  return {
    kind: 'sqlite', path: target, status: 'existing', exists: true,
    stat: { size: fileStat.size, mtime_ms: fileStat.mtimeMs },
    readonly_uri: readOnlySQLiteURI(target), tables, recording_table: recordingTable,
    recording_table_present: tables.includes(recordingTable),
  };
}

function sqlLiteral(value) {
  return `'${value.replaceAll("'", "''")}'`;
}

/** Observe one owned PostgreSQL namespace and target table with SELECT statements only. */
export function inspectPostgresSchema({ psql, schema, table = 'readings' } = {}) {
  if (typeof psql !== 'function') throw new TypeError('psql observer is required');
  const ownedSchema = assertOwnedSchema(schema);
  const recordingTable = assertIdentifier(table, 'recording table');
  const namespaceRows = String(psql(`SELECT nspname FROM pg_namespace WHERE nspname = ${sqlLiteral(ownedSchema)}`) ?? '')
    .split(/\r?\n/).map((value) => value.trim()).filter(Boolean);
  const tableRows = String(psql(
    `SELECT c.relname FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace `
      + `WHERE n.nspname = ${sqlLiteral(ownedSchema)} AND c.relname = ${sqlLiteral(recordingTable)} `
      + "AND c.relkind IN ('r','p','v','m','f')",
  ) ?? '').split(/\r?\n/).map((value) => value.trim()).filter(Boolean);
  return {
    kind: 'postgres', schema: ownedSchema, status: 'existing', namespace_exists: namespaceRows.includes(ownedSchema),
    recording_table: recordingTable, recording_table_present: tableRows.includes(recordingTable),
    observation: 'select-only',
  };
}

/** Create only an owned empty schema, then prove the recording table is still absent by SELECT. */
export function createOwnedPostgresSchema({ psql, schema, table = 'readings' } = {}) {
  if (typeof psql !== 'function') throw new TypeError('psql observer is required');
  const ownedSchema = assertOwnedSchema(schema);
  assertIdentifier(table, 'recording table');
  psql(`CREATE SCHEMA "${ownedSchema}"`);
  try {
    return inspectPostgresSchema({ psql, schema: ownedSchema, table });
  } catch (error) {
    // Acknowledged CREATE proves ownership. Rejected/unknown CREATE never reaches this cleanup.
    try { psql(`DROP SCHEMA "${ownedSchema}"`); }
    catch (cleanupError) { throw new AggregateError([error, cleanupError], 'owned PostgreSQL preflight and cleanup failed'); }
    throw error;
  }
}

export function assertFreshDestination(evidence) {
  if (!evidence || !['sqlite', 'postgres'].includes(evidence.kind)) throw new TypeError('unknown destination preflight');
  if (evidence.kind === 'postgres' && evidence.namespace_exists !== true) {
    throw new Error('owned PostgreSQL namespace was not observed');
  }
  if (evidence.recording_table_present !== false) {
    throw new Error(`recording table proof is not confirmed absent: ${evidence.recording_table ?? 'unknown'}`);
  }
  return evidence;
}
