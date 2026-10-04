import assert from 'node:assert/strict';
import { existsSync, mkdirSync, mkdtempSync, rmSync, statSync, symlinkSync } from 'node:fs';
import { join } from 'node:path';
import test, { after, before } from 'node:test';
import { sh } from './lib.mjs';
import {
  assertFreshDestination,
  createOwnedPostgresSchema,
  inspectPostgresSchema,
  preflightSQLiteDestination,
} from './destination-preflight.mjs';

let work;

before(() => { work = mkdtempSync('/tmp/gw-f-destination-preflight-'); });
after(() => { rmSync(work, { recursive: true, force: true }); });

test('missing SQLite preflight uses stat and leaves the destination absent', () => {
  const path = join(work, 'missing.db');
  const evidence = preflightSQLiteDestination({ path, work });
  assert.equal(evidence.exists, false);
  assert.equal(evidence.status, 'missing');
  assert.equal(evidence.stat, 'absent');
  assert.equal(existsSync(path), false);
});

test('existing SQLite preflight observes through a read-only URI and sees no recording table', () => {
  const path = join(work, 'existing.db');
  sh('sqlite3', [path, 'CREATE TABLE neighbor (id INTEGER);']);
  const before = statSync(path);
  const evidence = preflightSQLiteDestination({ path, work });
  const after = statSync(path);
  assert.equal(evidence.exists, true);
  assert.equal(evidence.status, 'existing');
  assert.equal(evidence.recording_table_present, false);
  assert.deepEqual(evidence.tables, ['neighbor']);
  assert.match(evidence.readonly_uri, /^file:.+\?mode=ro$/);
  assert.equal(after.size, before.size);
  assert.equal(after.mtimeMs, before.mtimeMs);
  assertFreshDestination(evidence);
});

test('fresh assertion rejects a destination already containing the recording table', () => {
  assert.throws(
    () => assertFreshDestination({ kind: 'sqlite', exists: true, recording_table_present: true }),
    /recording table/i,
  );
});

test('fresh assertion rejects unknown recording-table proof states', () => {
  for (const value of [undefined, null, 'false', 0]) {
    assert.throws(
      () => assertFreshDestination({ kind: 'sqlite', exists: true, recording_table_present: value }),
      /recording table proof/i,
      `proof state ${String(value)} must not pass as fresh`,
    );
  }
  assert.doesNotThrow(() => assertFreshDestination({ kind: 'sqlite', exists: true, recording_table_present: false }));
});

test('SQLite preflight rejects a path outside the owned run namespace', () => {
  assert.throws(
    () => preflightSQLiteDestination({ path: '/tmp/foreign-destination.db', work }),
    /owned/i,
  );
});

test('SQLite preflight rejects owned-looking directory and file aliases', () => {
  const foreign = mkdtempSync('/tmp/gw-f-preflight-neighbor-');
  const aliasWork = `${work}-alias`;
  try {
    const foreignDB = join(foreign, 'neighbor.db');
    sh('sqlite3', [foreignDB, 'CREATE TABLE neighbor (id INTEGER);']);
    symlinkSync(foreign, aliasWork);
    assert.throws(() => preflightSQLiteDestination({ path: join(aliasWork, 'neighbor.db'), work: aliasWork }), /alias|symlink/i);
    const nested = join(work, 'nested');
    mkdirSync(nested);
    symlinkSync(foreign, join(nested, 'alias'));
    assert.throws(() => preflightSQLiteDestination({ path: join(nested, 'alias', 'neighbor.db'), work }), /alias|symlink/i);
    symlinkSync(foreignDB, join(work, 'neighbor-alias.db'));
    assert.throws(() => preflightSQLiteDestination({ path: join(work, 'neighbor-alias.db'), work }), /alias|symlink/i);
    assert.equal(statSync(foreignDB).isFile(), true);
  } finally {
    rmSync(aliasWork, { force: true });
    rmSync(foreign, { recursive: true });
  }
});

function fakePsql({ tablePresent = false } = {}) {
  const calls = [];
  const psql = (sql) => {
    calls.push(sql);
    if (/^CREATE SCHEMA /.test(sql)) return '';
    if (sql.includes('FROM pg_namespace')) return 'gw_f_preflight\n';
    if (sql.includes('FROM pg_class')) return tablePresent ? 'readings\n' : '';
    throw new Error(`unexpected SQL: ${sql}`);
  };
  return { calls, psql };
}

test('PostgreSQL preflight uses SELECT-only observations and owned schema names', () => {
  const { calls, psql } = fakePsql();
  const evidence = inspectPostgresSchema({ psql, schema: 'gw_f_preflight', table: 'readings' });
  assert.equal(evidence.namespace_exists, true);
  assert.equal(evidence.recording_table_present, false);
  assert.ok(calls.length >= 2);
  assert.ok(calls.every((sql) => /^\s*SELECT\b/i.test(sql)));
  assertFreshDestination(evidence);
});

test('owned PostgreSQL setup creates only the empty schema, then rechecks by SELECT', () => {
  const { calls, psql } = fakePsql();
  const evidence = createOwnedPostgresSchema({ psql, schema: 'gw_f_preflight', table: 'readings' });
  assert.equal(evidence.recording_table_present, false);
  assert.match(calls[0], /^CREATE SCHEMA "gw_f_preflight"$/);
  assert.equal(calls.slice(1).every((sql) => /^\s*SELECT\b/i.test(sql)), true);
  assert.equal(calls.some((sql) => /CREATE TABLE|INSERT\s+INTO/i.test(sql)), false);
});

test('PostgreSQL preflight rejects foreign schemas and existing recording tables', () => {
  assert.throws(
    () => inspectPostgresSchema({ psql: () => '', schema: 'public', table: 'readings' }),
    /owned/i,
  );
  const { psql } = fakePsql({ tablePresent: true });
  assert.throws(
    () => assertFreshDestination(inspectPostgresSchema({ psql, schema: 'gw_f_preflight', table: 'readings' })),
    /recording table/i,
  );
});

test('owned schema is removed when observation fails after acknowledged creation', () => {
  const calls = [];
  assert.throws(() => createOwnedPostgresSchema({ schema: 'gw_f_failed_observer', psql: (sql) => {
    calls.push(sql);
    if (/^(CREATE|DROP) SCHEMA /.test(sql)) return '';
    throw new Error('observation failed');
  } }), /observation failed/i);
  assert.deepEqual(calls.filter((sql) => /^DROP /.test(sql)), ['DROP SCHEMA "gw_f_failed_observer"']);
});

test('schema validation and failed CREATE never drop an unowned existing namespace', () => {
  const calls = [];
  assert.throws(() => createOwnedPostgresSchema({ schema: 'gw_f_invalid_table', table: 'unsafe.table',
    psql: (sql) => calls.push(sql) }), /identifier/i);
  assert.deepEqual(calls, []);
  assert.throws(() => createOwnedPostgresSchema({ schema: 'gw_f_existing', psql: (sql) => {
    calls.push(sql);
    throw new Error('CREATE rejected: namespace exists');
  } }), /namespace exists/i);
  assert.equal(calls.length, 1);
  assert.match(calls[0], /^CREATE SCHEMA /);
});
