import { existsSync, statSync } from 'node:fs';
import net from 'node:net';
import { sh, validateOwnedPostgresTarget } from './lib.mjs';

const IDENTIFIER = /^[A-Za-z_][A-Za-z0-9_]*$/;

export function quoteIdentifier(value, label = 'SQL identifier') {
  if (!IDENTIFIER.test(String(value ?? ''))) throw new Error(`${label} is not a safe identifier`);
  return `"${value}"`;
}

export function sqliteJSON(path, sql) {
  if (!existsSync(path)) return [];
  if (!statSync(path).isFile()) throw new Error('SQLite observation target is not a regular file');
  const raw = sh('sqlite3', ['-readonly', '-json', '-cmd', '.timeout 15000', `file:${path}?mode=ro`, sql]);
  return raw ? JSON.parse(raw) : [];
}

export function sqliteExec(path, sql) {
  if (!existsSync(path)) throw new Error('SQLite mutation target is missing');
  if (!statSync(path).isFile()) throw new Error('SQLite mutation target is not a regular file');
  return sh('sqlite3', ['-cmd', '.timeout 15000', path, sql]);
}

// The fresh witness may create only its own non-target neighbor fixture. Keep
// this separate from sqliteExec so a missing recording destination can never be
// provisioned by an observation or fault-control path.
export function sqliteOwnedNeighborExec(path, sql) {
  if (!/^\/tmp\/gw-f-[A-Za-z0-9_-]+\/neighbor\.db$/.test(path)) throw new Error('SQLite owned fixture path is outside the fresh neighbor namespace');
  if (existsSync(path) && !statSync(path).isFile()) throw new Error('SQLite owned fixture is not a regular file');
  return sh('sqlite3', ['-cmd', '.timeout 15000', path, sql]);
}

export function parseOwnedPostgresDSN(raw) {
  const values = Object.fromEntries(String(raw ?? '').trim().split(/\s+/).filter(Boolean).map((part) => {
    const separator = part.indexOf('=');
    return separator < 0 ? [part, ''] : [part.slice(0, separator), part.slice(separator + 1)];
  }));
  const image = sh('docker', ['inspect', '--format', '{{.Config.Image}}', 'gw-wg-pg-test']);
  validateOwnedPostgresTarget(values, image);
  return values;
}

export function postgresText({ database = 'gwtest' } = {}, sql) {
  return sh('docker', [
    'exec', '-i', 'gw-wg-pg-test', 'psql', '-v', 'ON_ERROR_STOP=1', '-U', 'postgres',
    '-d', database, '-At', '-c', sql,
  ]);
}

export function postgresJSON(connection, sql) {
  const statement = `SELECT COALESCE(json_agg(row_data), '[]'::json) FROM (${sql}) AS row_data`;
  const raw = postgresText(connection, statement);
  return raw ? JSON.parse(raw) : [];
}

export function postgresExec(connection, sql) {
  return postgresText(connection, sql);
}

export async function makeOwnedPostgresProxy({ listenPort = 55434, upstreamPort = 55432 } = {}) {
  const sockets = new Set();
  let server;
  const relay = (client) => {
    const upstream = net.connect({ host: '127.0.0.1', port: upstreamPort });
    sockets.add(client); sockets.add(upstream);
    const close = () => { sockets.delete(client); sockets.delete(upstream); };
    client.on('error', () => upstream.destroy());
    upstream.on('error', () => client.destroy());
    client.on('close', close); upstream.on('close', close);
    client.pipe(upstream); upstream.pipe(client);
  };
  const start = async () => {
    if (server?.listening) return;
    server = net.createServer(relay);
    await new Promise((resolve, reject) => {
      server.once('error', reject);
      server.listen(listenPort, '127.0.0.1', resolve);
    });
  };
  const stop = async () => {
    for (const socket of sockets) socket.destroy();
    sockets.clear();
    if (server?.listening) await new Promise((resolve) => server.close(resolve));
  };
  await start();
  return { start, stop, listenPort, upstreamPort };
}

export function readSQLiteTable(path, table) {
  const name = quoteIdentifier(table, 'SQLite table');
  return sqliteJSON(path, `SELECT * FROM ${name} ORDER BY bucket_start, record_id`);
}

export function readPostgresTable(connection, schema, table) {
  const namespace = quoteIdentifier(schema, 'PostgreSQL schema');
  const name = quoteIdentifier(table, 'PostgreSQL table');
  const columns = postgresJSON(connection, `SELECT column_name,data_type FROM information_schema.columns
    WHERE table_schema='${schema}' AND table_name='${table}' ORDER BY ordinal_position`);
  if (columns.length === 0) return [];
  const projection = columns.map((column) => {
    const identifier = quoteIdentifier(column.column_name, 'PostgreSQL column');
    const dataType = String(column.data_type ?? '').toLowerCase();
    return ['numeric', 'decimal', 'bigint'].includes(dataType)
      ? `${identifier}::text AS ${identifier}`
      : identifier;
  }).join(',');
  return postgresJSON(connection, `SELECT ${projection} FROM ${namespace}.${name} ORDER BY ${quoteIdentifier('bucket_start')}, ${quoteIdentifier('record_id')}`);
}

export function readSQLiteReceipt(path, effectKey) {
  if (!existsSync(path)) return [];
  if (!statSync(path).isFile()) throw new Error('SQLite observation target is not a regular file');
  const escaped = String(effectKey).replaceAll("'", "''");
  return sqliteJSON(path, `SELECT effect_key,payload_digest,committed_at FROM "gw_effect_receipts" WHERE effect_key = '${escaped}'`);
}

export function readPostgresReceipt(connection, schema, effectKey) {
  const escaped = String(effectKey).replaceAll("'", "''");
  return postgresJSON(connection, `SELECT effect_key,payload_digest,committed_at FROM ${quoteIdentifier(schema, 'PostgreSQL schema')}."gw_effect_receipts" WHERE effect_key = '${escaped}'`);
}
