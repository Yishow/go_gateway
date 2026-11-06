// @vitest-environment node
import { afterEach, describe, expect, it } from 'vitest';
import { DatabaseSync } from 'node:sqlite';
import { mkdtempSync, rmSync } from 'node:fs';
import { join } from 'node:path';
import { tmpdir } from 'node:os';
import { assertEmptySources, assertNoConfigurationChanges, changedTables, snapshotSQLite } from '../../e2e/helpers/recording-home-sqlite';
const roots: string[] = [];
afterEach(() => { for (const root of roots.splice(0)) rmSync(root, { recursive: true, force: true }); });
function database() {
  const root = mkdtempSync(join(tmpdir(), 'home-sqlite-unit-')); roots.push(root);
  const path = join(root, 'datalink.db'); const db = new DatabaseSync(path);
  db.exec('PRAGMA journal_mode=WAL; CREATE TABLE devices(id TEXT, updated_at TEXT, payload BLOB); CREATE TABLE timeseries(value INTEGER)');
  return { path, db };
}
describe('read-only table snapshots', () => {
  it('sees committed WAL data without checkpointing, preserves timestamps and large integers, and ignores row insertion order', () => {
    const { path, db } = database();
    try {
      db.exec("INSERT INTO devices VALUES ('b','first',X'00FF'), ('a','first',NULL); INSERT INTO timeseries VALUES (9223372036854775807)");
      const before = snapshotSQLite(path);
      expect(before.tables.devices.rows).toBe(2);
      db.exec("DELETE FROM devices; INSERT INTO devices VALUES ('a','first',NULL), ('b','first',X'00FF')");
      expect(changedTables(before, snapshotSQLite(path))).toEqual([]);
      db.exec("UPDATE devices SET updated_at='changed' WHERE id='a'");
      const after = snapshotSQLite(path);
      expect(changedTables(before, after).map(item => item.name)).toEqual(['devices']);
      expect(() => assertNoConfigurationChanges(before, after)).toThrow('devices');
    } finally { db.close(); }
  });
  it('reports operational changes separately, never silently calls the whole DB unchanged', () => {
    const { path, db } = database();
    try { const before = snapshotSQLite(path); db.exec('INSERT INTO timeseries VALUES (1)'); const after = snapshotSQLite(path);
      expect(changedTables(before, after).map(item => [item.name, item.kind])).toEqual([['timeseries', 'operational']]);
      expect(() => assertNoConfigurationChanges(before, after)).not.toThrow();
    } finally { db.close(); }
  });
  it('fails closed on missing tables and missing database rather than creating it', () => {
    const { path, db } = database(); db.close();
    expect(() => assertEmptySources(snapshotSQLite(path))).toThrow();
    expect(() => snapshotSQLite(`${path}-missing`)).toThrow();
  });
});
