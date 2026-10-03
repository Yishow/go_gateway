// Real disposable SQL destinations; only gateway writes sample/receipt rows.
import assert from 'node:assert/strict';
import { renameSync } from 'node:fs';
import { join } from 'node:path';
import net from 'node:net';
import { sh } from './lib.mjs';
import { query } from './quality-lib.mjs';

const tableSQLite = 'CREATE TABLE readings (line TEXT, temperature INTEGER, pressure INTEGER, running INTEGER, batch INTEGER, prov TEXT);';
const tablePostgres = 'CREATE TABLE readings (line TEXT, temperature BIGINT, pressure BIGINT, running BOOLEAN, batch BIGINT, prov JSONB);';
const receiptDDL = 'CREATE TABLE gw_effect_receipts (effect_key TEXT PRIMARY KEY, payload_digest TEXT NOT NULL, committed_at TEXT NOT NULL);';
const rowsSQL = 'SELECT line, temperature, CAST(pressure AS TEXT) AS pressure, running, CAST(batch AS TEXT) AS batch, CAST(prov AS TEXT) AS prov FROM readings';
const receiptsSQL = 'SELECT effect_key, payload_digest FROM gw_effect_receipts';

async function makeProxy(port) {
  const sockets = new Set();
  let server;
  const start = async () => {
    server = net.createServer((client) => {
      const upstream = net.connect({ host: '127.0.0.1', port: 55432 });
      sockets.add(client); sockets.add(upstream);
      client.on('error', () => upstream.destroy());
      upstream.on('error', () => client.destroy());
      client.on('close', () => sockets.delete(client));
      upstream.on('close', () => sockets.delete(upstream));
      client.pipe(upstream); upstream.pipe(client);
    });
    await new Promise((resolve, reject) => {
      server.once('error', reject);
      server.listen(port, '127.0.0.1', resolve);
    }).catch(async (error) => {
      for (const socket of sockets) socket.destroy();
      sockets.clear();
      if (server && !server.listening) server.close(() => {});
      throw error;
    });
  };
  const stop = async () => {
    for (const socket of sockets) socket.destroy();
    sockets.clear();
    if (server?.listening) await new Promise((resolve) => server.close(resolve));
  };
  await start();
  return { start, stop };
}

export async function recoveryTargets(kind, work, runNumber) {
  if (kind === 'sqlite') {
    const a = join(work, 'a.db'); const b = join(work, 'b.db'); const other = join(work, 'elsewhere.db');
    for (const path of [a, b, other]) sh('sqlite3', [path, tableSQLite + receiptDDL]);
    let offline = false;
    return {
      destination: { kind, path: a, table: 'readings' },
      bConfig: { dsn: b }, otherConfig: { dsn: other },
      database: (name) => ({ a, b, other })[name], schema: () => 'main',
      queries: { rows: rowsSQL, receipts: receiptsSQL },
      rows: (name) => query(name === 'a' && offline ? `${a}.away` : ({ a, b, other })[name], rowsSQL),
      receipts: (name) => query(name === 'a' && offline ? `${a}.away` : ({ a, b, other })[name], receiptsSQL),
      disconnect: async () => { assert.equal(offline, false); renameSync(a, `${a}.away`); offline = true; },
      reconnect: async () => { if (offline) { renameSync(`${a}.away`, a); offline = false; } },
      cleanup: async () => ({
        retained: [a, `${a}.away`, b, other],
        reason: 'diagnostic-evidence-retained',
      }),
    };
  }
  assert.equal(kind, 'postgres');
  const dsn = Object.fromEntries((process.env.POSTGRES_DSN ?? '').trim().split(/\s+/).filter(Boolean).map((part) => part.split('=')));
  assert.equal(dsn.host, '127.0.0.1', 'only the owned loopback PostgreSQL fixture is allowed');
  assert.equal(dsn.port, '55432'); assert.equal(dsn.dbname, 'gwtest'); assert.equal(dsn.user, 'postgres');
  const ownedImage = sh('docker', ['inspect', '--format', '{{.Config.Image}}', 'gw-wg-pg-test']);
  assert.match(ownedImage, /^postgres:16(?:-alpine)?$/);
  const schemas = { a: `gw_f_recovery_a_${runNumber}`, b: `gw_f_recovery_b_${runNumber}`, other: `gw_f_recovery_a_${runNumber}` };
  const otherDB = `gw_f_recovery_other_${runNumber}`;
  const psql = (db, sql) => sh('docker', ['exec', 'gw-wg-pg-test', 'psql', '-v', 'ON_ERROR_STOP=1', '-U', 'postgres', '-d', db, '-At', '-c', sql]);
  const createdSchemas = [];
  const createdDatabases = [];
  let proxy;
  const cleanupPartial = async () => {
    const removedSchemas = []; const removedDatabases = []; const errors = [];
    let proxyStopped = !proxy;
    try { if (proxy) { await proxy.stop(); proxyStopped = true; } } catch (error) { errors.push(`proxy: ${error.message}`); }
    for (const { db, schema } of [...createdSchemas].reverse()) {
      try { psql(db, `DROP SCHEMA "${schema}" CASCADE`); removedSchemas.push({ db, schema }); }
      catch (error) { errors.push(`schema ${db}.${schema}: ${error.message}`); }
    }
    for (const db of [...createdDatabases].reverse()) {
      try { psql('gwtest', `DROP DATABASE "${db}" WITH (FORCE)`); removedDatabases.push(db); }
      catch (error) { errors.push(`database ${db}: ${error.message}`); }
    }
    return { removed: { schemas: removedSchemas, databases: removedDatabases }, proxy_stopped: proxyStopped, errors };
  };
  const create = (db, schema) => {
    psql(db, `CREATE SCHEMA "${schema}"`);
    createdSchemas.push({ db, schema });
    psql(db, `SET search_path TO "${schema}"; ${tablePostgres}${receiptDDL}`);
  };
  try {
    create('gwtest', schemas.a); create('gwtest', schemas.b);
    psql('gwtest', `CREATE DATABASE "${otherDB}"`);
    createdDatabases.push(otherDB);
    create(otherDB, schemas.other);
    proxy = await makeProxy(55433);
  } catch (error) {
    error.setupCleanup = await cleanupPartial();
    throw error;
  }
  const json = (name, sql) => {
    const db = name === 'other' ? otherDB : 'gwtest';
    const scoped = sql.replaceAll('readings', `"${schemas[name]}".readings`).replaceAll('gw_effect_receipts', `"${schemas[name]}".gw_effect_receipts`);
    return JSON.parse(psql(db, `SELECT COALESCE(json_agg(t),'[]'::json) FROM (${scoped}) AS t`));
  };
  const config = (port, database) => ({ host: '127.0.0.1', port, user: 'postgres', password: dsn.password, database, sslmode: 'disable' });
  return {
    destination: { kind, host: '127.0.0.1', port: 55433, database: 'gwtest', user: 'postgres', password: dsn.password, schema: schemas.a, table: 'readings' },
    bConfig: config(55432, 'gwtest'), otherConfig: config(55432, otherDB),
    database: (name) => name === 'other' ? otherDB : 'gwtest', schema: (name) => schemas[name],
    queries: { rows: rowsSQL, receipts: receiptsSQL, schema_scope: schemas },
    rows: (name) => json(name, rowsSQL), receipts: (name) => json(name, receiptsSQL),
    disconnect: proxy.stop, reconnect: proxy.start,
    cleanup: async () => {
      let cleanupError;
      try { await proxy.stop(); } catch (error) { cleanupError ??= error; }
      for (const name of ['a', 'b']) {
        try { psql('gwtest', `DROP SCHEMA "${schemas[name]}" CASCADE`); } catch (error) { cleanupError ??= error; }
      }
      try { psql('gwtest', `DROP DATABASE "${otherDB}" WITH (FORCE)`); } catch (error) { cleanupError ??= error; }
      if (cleanupError) throw cleanupError;
      return { removed: { schemas: [schemas.a, schemas.b, schemas.other], databases: [otherDB], proxy: '127.0.0.1:55433' } };
    },
  };
}
