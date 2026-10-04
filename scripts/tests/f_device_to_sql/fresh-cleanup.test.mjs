import test from 'node:test';
import assert from 'node:assert/strict';
import { execFileSync } from 'node:child_process';
import { chmodSync, existsSync, mkdirSync, mkdtempSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { isFreshCleanExit, removeOwnedWork } from './fresh-ui.mjs';

test('requires a zero exit code with no signal or forced kill', () => {
  assert.equal(isFreshCleanExit({ state: 'exited', exit_code: 0, signal: null, kill_sent: false }), true);
  for (const status of [
    { state: 'exited', exit_code: null, signal: null, kill_sent: false },
    { state: 'exited', exit_code: 1, signal: null, kill_sent: false },
    { state: 'exited', exit_code: 0, signal: 'SIGTERM', kill_sent: false },
    { state: 'exited', exit_code: 0, signal: null, kill_sent: true },
    { state: 'timeout', exit_code: 0, signal: null, kill_sent: false },
  ]) assert.equal(isFreshCleanExit(status), false, JSON.stringify(status));
});

test('removeOwnedWork reports removal and never throws, so evidence can still be written', () => {
  const root = mkdtempSync(join(tmpdir(), 'gw-f-owned-work-'));
  try {
    const work = join(root, 'work');
    mkdirSync(work);
    writeFileSync(join(work, 'gateway.db'), 'x');
    assert.deepEqual(removeOwnedWork(work), { removed: true });
    assert.equal(existsSync(work), false);
    assert.deepEqual(removeOwnedWork(work), { removed: true }, 'an already absent namespace is clean');

    const stuck = join(root, 'stuck');
    mkdirSync(stuck);
    writeFileSync(join(stuck, 'gateway.db'), 'x');
    chmodSync(root, 0o500); // removing entries of a read-only parent fails with EACCES
    let outcome;
    assert.doesNotThrow(() => { outcome = removeOwnedWork(stuck); });
    assert.equal(outcome.removed, false);
    assert.match(outcome.error, /EACCES|EPERM/);
    assert.equal(existsSync(stuck), true, 'a failed cleanup is never reported as removed');
  } finally {
    chmodSync(root, 0o700);
    rmSync(root, { recursive: true, force: true });
  }
});

test('default binary paths follow the F_GATEWAY_BINARY and F_SIMULATOR_BINARY overrides', () => {
  const probe = "import('./fresh-ui.mjs').then((m) => console.log(JSON.stringify([m.DEFAULT_GATEWAY_BINARY, m.DEFAULT_SIMULATOR_BINARY])))";
  const read = (env) => JSON.parse(execFileSync(process.execPath, ['-e', probe], {
    cwd: new URL('.', import.meta.url).pathname, env: { ...process.env, F_GATEWAY_BINARY: '', F_SIMULATOR_BINARY: '', ...env },
  }).toString());
  assert.deepEqual(read({ F_GATEWAY_BINARY: '/tmp/custom-gateway', F_SIMULATOR_BINARY: '/tmp/custom-sim' }), ['/tmp/custom-gateway', '/tmp/custom-sim']);
  const unset = read({});
  assert.match(unset[0], /gw-F-test-ui-production-final-v6$/, 'the unset default points at the final build');
});
