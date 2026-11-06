import { createHash } from 'node:crypto';
import { DatabaseSync } from 'node:sqlite';

export interface TableHash { rows: number; sha256: string; kind: 'configuration' | 'operational' }
export interface SQLiteSnapshot { reader: string; tables: Record<string, TableHash> }
// Everything not explicitly listed remains configuration, including all timestamps,
// source links, revisions, settings and future tables. Operational changes remain visible.
const operational = new Set(['timeseries', 'acquisition_anomalies', 'wg_runtime_versions',
  'wg_delivery_samples', 'wg_delivery_checkpoints', 'wg_delivery_buckets', 'wg_delivery_outbox',
  'wg_delivery_receipts', 'wg_delivery_counter_baselines', 'wg_delivery_sample_handoffs',
  'wg_delivery_mixed_progress', 'wg_delivery_mixed_outcomes', 'gw_delivery_journal',
  'gw_delivery_outbox', 'gw_delivery_receipts']);
const quote = (value: string) => `"${value.replaceAll('"', '""')}"`;
function serialize(value: unknown): string {
  return JSON.stringify(value, (_key, item: unknown) => typeof item === 'bigint' ? { integer: item.toString() } :
    item instanceof Uint8Array ? { blob: Buffer.from(item).toString('hex') } : item);
}
/** Read-only WAL-aware consistent transaction, not a hash of mutable SQLite files. */
export function snapshotSQLite(databasePath: string): SQLiteSnapshot {
  const db = new DatabaseSync(databasePath, { readOnly: true });
  try {
    db.exec('PRAGMA query_only=ON; BEGIN');
    const tables: SQLiteSnapshot['tables'] = {};
    const schemas = db.prepare("SELECT name, sql FROM sqlite_schema WHERE type='table' ORDER BY name").all();
    for (const schema of schemas) {
      const name = String(schema.name);
      const statement = db.prepare(`SELECT * FROM ${quote(name)}`);
      statement.setReadBigInts(true);
      const rows = statement.all().map(serialize).sort();
      const sha256 = createHash('sha256').update(serialize({ schema: schema.sql, rows })).digest('hex');
      tables[name] = { rows: rows.length, sha256, kind: operational.has(name) ? 'operational' : 'configuration' };
    }
    db.exec('ROLLBACK');
    return { reader: `node:sqlite ${process.version}; readOnly + query_only; one transaction`, tables };
  } finally { db.close(); }
}
export function changedTables(before: SQLiteSnapshot, after: SQLiteSnapshot) {
  return [...new Set([...Object.keys(before.tables), ...Object.keys(after.tables)])].sort().filter(name =>
    before.tables[name]?.sha256 !== after.tables[name]?.sha256).map(name => ({ name,
      kind: before.tables[name]?.kind ?? after.tables[name].kind, before: before.tables[name] ?? null, after: after.tables[name] ?? null }));
}
export function assertNoConfigurationChanges(before: SQLiteSnapshot, after: SQLiteSnapshot): void {
  const changed = changedTables(before, after).filter(table => table.kind === 'configuration');
  if (changed.length) throw new Error(`Configuration changed: ${changed.map(table => table.name).join(', ')}`);
}
export function assertEmptySources(snapshot: SQLiteSnapshot): void {
  for (const name of ['devices', 'source_rules', 'source_rule_links', 'points', 'tags', 'mappings', 'database_connectors', 'database_target_mappings', 'write_groups', 'recording_plans']) {
    if (snapshot.tables[name]?.rows !== 0) throw new Error(`Refusing existing or missing source table: ${name}`);
  }
}
export function assertSeedTableCounts(snapshot: SQLiteSnapshot): void {
  for (const [name, count] of Object.entries({ devices: 20, source_rules: 20, source_rule_links: 2000, points: 2000,
    tags: 0, mappings: 0, database_connectors: 0, database_target_mappings: 0, write_groups: 0, recording_plans: 0 })) {
    if (snapshot.tables[name]?.rows !== count) throw new Error(`Seed table count mismatch: ${name}`);
  }
}
