import assert from 'node:assert/strict';
import { createServer } from 'node:http';
import { mkdirSync, existsSync, rmSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { spawn } from 'node:child_process';
import test, { after, before } from 'node:test';
import { MIXED_POINT_EXPECTATIONS, collectAndValidatePointObservations } from './mixed-identity.mjs';

const ROOT = new URL('../../..', import.meta.url).pathname;

let server;
let port;
let harness;

before(async () => {
  server = createServer((request, response) => {
    if (request.url === '/api/v1/datalink/non2xx') {
      response.writeHead(500, { 'content-type': 'application/json' });
      response.end(JSON.stringify({ error: { code: 'fixture_failure' } }));
      return;
    }
    if (request.url === '/api/v1/datalink/missing-data') {
      response.writeHead(200, { 'content-type': 'application/json' });
      response.end(JSON.stringify({ success: true }));
      return;
    }
    if (request.url === '/hang' || request.url === '/studio/v2') {
      response.writeHead(200, { 'content-type': 'text/plain' });
      response.write('partial');
      return;
    }
    response.writeHead(404); response.end('{}');
  });
  await new Promise((resolve) => server.listen(0, '127.0.0.1', resolve));
  port = server.address().port;
  process.env.GW_PORT = String(port);
  harness = await import(`./lib.mjs?mixed-regression=${Date.now()}`);
});

after(async () => new Promise((resolve, reject) => {
  server.closeAllConnections();
  server.close((error) => error ? reject(error) : resolve());
}));

test('invalid destination kind fails before touching its work path', async () => {
  const work = '/tmp/gw-f-..';
  mkdirSync(work, { recursive: true });
  const sentinel = join(work, 'sentinel');
  writeFileSync(sentinel, 'preserve');
  try {
    const child = spawn(process.execPath, ['scripts/tests/f_device_to_sql/run.mjs', '..'], { cwd: ROOT, stdio: ['ignore', 'pipe', 'pipe'] });
    let output = '';
    child.stdout.on('data', (chunk) => { output += chunk; });
    child.stderr.on('data', (chunk) => { output += chunk; });
    const exitCode = await new Promise((resolve) => child.once('close', resolve));
    assert.equal(exitCode, 1);
    assert.match(output, /invalid destination kind/);
    assert.equal(existsSync(sentinel), true);
  } finally {
    rmSync(work, { recursive: true, force: true });
  }
});

test('F_RUN_ID accepts only a safe run namespace and derives work below /tmp/gw-f-*', () => {
  const namespace = harness.makeRunNamespace('sqlite', 'sqlite-regression_1');
  assert.equal(namespace.run_id, 'sqlite-regression_1');
  assert.equal(namespace.work, '/tmp/gw-f-sqlite-regression_1');
  for (const value of ['../escape', '/tmp/custom', 'run with spaces', '']) {
    assert.throws(() => harness.makeRunNamespace('sqlite', value), /F_RUN_ID/);
  }
});

test('startProcess treats shell characters in log paths literally and stop waits for child exit', async () => {
  const work = `/tmp/f-mixed-process-${process.pid}`;
  mkdirSync(work, { recursive: true });
  const name = 'child-$(touch injected-marker)';
  try {
    const child = harness.startProcess(work, name, process.execPath, ['-e', 'setTimeout(() => {}, 1000)']);
    const statuses = await harness.stopProcesses([child], 2000);
    assert.equal(statuses[0].state, 'exited');
    assert.equal(existsSync(join(work, 'injected-marker')), false);
    assert.equal(existsSync(join(work, `${name}.log`)), true);
  } finally {
    rmSync(work, { recursive: true, force: true });
  }
});

test('stopProcesses preserves exit status for an already exited child', async () => {
  const work = `/tmp/f-mixed-process-exited-${process.pid}`;
  mkdirSync(work, { recursive: true });
  let child;
  try {
    child = harness.startProcess(work, 'exited-nonzero', process.execPath, ['-e', 'process.exit(17)']);
    await new Promise((resolve, reject) => {
      child.once('close', resolve);
      child.once('error', reject);
    });
    const statuses = await harness.stopProcesses([child], 2000);
    assert.equal(statuses[0].state, 'exited');
    assert.equal(statuses[0].exit_code, 17);
    assert.equal(statuses[0].signal, null);
  } finally {
    rmSync(work, { recursive: true, force: true });
  }
});

test('stopProcesses escalates a child that ignores SIGTERM and waits for close', async () => {
  const work = `/tmp/f-mixed-process-ignore-term-${process.pid}`;
  mkdirSync(work, { recursive: true });
  let child;
  try {
    child = harness.startProcess(work, 'ignore-term', process.execPath, ['-e', "process.on('SIGTERM', () => {}); setInterval(() => {}, 1000);"]);
    await new Promise((resolve) => setTimeout(resolve, 100));
    const statuses = await harness.stopProcesses([child], 50);
    assert.equal(statuses[0].state, 'exited');
    assert.equal(statuses[0].term_state, 'timeout');
    assert.equal(statuses[0].kill_sent, true);
    assert.equal(statuses[0].signal, 'SIGKILL');
  } finally {
    if (child?.pid) {
      try { process.kill(-child.pid, 'SIGKILL'); } catch (error) { if (error?.code !== 'ESRCH') throw error; }
    }
    rmSync(work, { recursive: true, force: true });
  }
});

test('boundedFetch rejects a response whose body never completes', async () => {
  await assert.rejects(
    harness.boundedFetch(`http://127.0.0.1:${port}/hang`, {}, 40),
    /Timeout|aborted|abort/i,
  );
});

test('boundedFetch preserves a caller abort signal', async () => {
  const controller = new AbortController();
  const request = harness.boundedFetch(`http://127.0.0.1:${port}/hang`, { signal: controller.signal }, 5000);
  setTimeout(() => controller.abort(new Error('caller cancelled')), 20);
  await assert.rejects(request, /abort|cancel/i);
});

test('waitForPortFree does not treat an occupied but stalled port as free', async () => {
  await assert.rejects(harness.waitForPortFree(60), /still in use/);
});

test('api rejects non-2xx and missing-data responses', async () => {
  await assert.rejects(harness.api('/non2xx'), /API \/non2xx failed \(500/);
  await assert.rejects(harness.api('/missing-data'), /response is missing data/);
});

test('PostgreSQL target validation rejects an arbitrary endpoint or image', () => {
  const owned = { host: '127.0.0.1', port: '55432', dbname: 'gwtest', user: 'postgres', password: 'gwtest' };
  assert.doesNotThrow(() => harness.validateOwnedPostgresTarget(owned, 'postgres:16-alpine'));
  assert.throws(() => harness.validateOwnedPostgresTarget({ ...owned, host: '10.0.0.4' }, 'postgres:16-alpine'), /owned loopback/);
  assert.throws(() => harness.validateOwnedPostgresTarget(owned, 'postgres:15'), /postgres:16/);
});

test('PostgreSQL schema cleanup must drop and verify absence before passing', () => {
  const calls = [];
  const passed = harness.cleanupOwnedPostgresSchema({ schema: 'gw_f_mixed_test' }, (sql) => { calls.push(sql); return ''; });
  assert.equal(passed.passed, true);
  assert.deepEqual(calls, [
    'DROP SCHEMA "gw_f_mixed_test" CASCADE',
    "SELECT nspname FROM pg_namespace WHERE nspname = 'gw_f_mixed_test'",
  ]);
  const failed = harness.cleanupOwnedPostgresSchema({ schema: 'gw_f_mixed_test' }, (sql) => {
    if (sql.startsWith('DROP')) throw new Error('drop refused');
    return '';
  });
  assert.equal(failed.passed, false);
});

function groupsFixture() {
  return Object.entries(MIXED_POINT_EXPECTATIONS).map(([line, expected]) => ({
    name: `Line ${line}`,
    members: expected.map((point, index) => ({ device_id: `dev-${line}`, point_id: `${line}-${index}` })),
  }));
}

function pointsFixture(line) {
  return MIXED_POINT_EXPECTATIONS[line].map((point, index) => ({
    id: `${line}-${index}`,
    device_id: `dev-${line}`,
    ...point,
  }));
}

test('point identity validation rejects API failure, empty, missing, wrong device, type, and address', async () => {
  const groups = groupsFixture();
  const valid = await collectAndValidatePointObservations(groups, async (deviceID) => pointsFixture(deviceID.endsWith('A') ? 'A' : 'B'));
  assert.deepEqual(valid.failures, []);
  const cases = [
    ['500', async () => { throw new Error('API /points failed (500)'); }, /point API failed/],
    ['empty', async () => [], /expected point identity/],
    ['missing', async (deviceID) => pointsFixture(deviceID.endsWith('A') ? 'A' : 'B').slice(1), /point .* is missing/],
    ['wrong-device', async (deviceID) => pointsFixture(deviceID.endsWith('A') ? 'A' : 'B').map((point, index) => index === 0 ? { ...point, device_id: 'other-device' } : point), /belongs to/],
    ['wrong-type', async (deviceID) => pointsFixture(deviceID.endsWith('A') ? 'A' : 'B').map((point, index) => index === 0 ? { ...point, data_type: 'float64' } : point), /unexpected address\/type/],
    ['wrong-address', async (deviceID) => pointsFixture(deviceID.endsWith('A') ? 'A' : 'B').map((point, index) => index === 0 ? { ...point, address: '49999' } : point), /unexpected address\/type/],
  ];
  for (const [name, reader, message] of cases) {
    const actual = await collectAndValidatePointObservations(groups, reader);
    assert.match(actual.failures.join('; '), message, name);
  }
});
