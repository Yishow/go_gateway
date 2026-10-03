import assert from 'node:assert/strict';
import { spawn } from 'node:child_process';
import test from 'node:test';
import { cleanupFixture, publishAfterCleanup, safeRead } from './fixture-cleanup.mjs';

test('cleanup waits for the owned child to exit before reporting success', async () => {
  const child = spawn(process.execPath, ['-e', 'setInterval(() => {}, 1000)'], { detached: true, stdio: 'ignore' });
  const calls = [];
  const result = await cleanupFixture(
    { close: async () => calls.push('browser') }, { owned: child },
    async () => { assert.ok(child.exitCode !== null || child.signalCode !== null); calls.push('port'); },
    { cleanup: () => calls.push('role') }, { cleanup: async () => calls.push('target') },
  );
  assert.equal(result.passed, true);
  assert.deepEqual(calls, ['browser', 'port', 'role', 'target']);
});

test('cleanup tolerates setup failure before optional resources exist', async () => {
  const child = spawn(process.execPath, ['-e', 'setInterval(() => {}, 1000)'], { detached: true, stdio: 'ignore' });
  const calls = [];
  const result = await cleanupFixture(
    undefined, [child],
    async () => { assert.ok(child.exitCode !== null || child.signalCode !== null); calls.push('port'); }, undefined, undefined,
  );
  assert.equal(result.passed, true);
  assert.deepEqual(calls, ['port']);
  assert.equal(result.steps.browser_closed.passed, true);
  assert.equal(result.steps.owned_targets_removed.passed, true);
});

test('cleanup escalates an owned child that ignores SIGTERM', async () => {
  const priorExitCode = process.exitCode;
  const child = spawn(process.execPath, ['-e', 'process.on("SIGTERM", () => {}); console.log("ready"); setInterval(() => {}, 1000)'], { detached: true, stdio: ['ignore', 'pipe', 'ignore'] });
  try {
    await new Promise((resolve) => child.stdout.once('data', resolve));
    const result = await cleanupFixture(undefined, [child], undefined, undefined, undefined);
    assert.equal(result.passed, true);
    assert.equal(result.steps.owned_children_stopped.passed, true);
    const status = result.steps.owned_children_stopped.statuses[0];
    assert.equal(status.term_state, 'timeout');
    assert.equal(status.kill_sent, true);
    assert.equal(status.state, 'exited');
    assert.equal(status.signal, 'SIGKILL');
  } finally {
    try { if (child.exitCode === null && child.signalCode === null) process.kill(-child.pid, 'SIGKILL'); } catch { /* test-owned child */ }
    process.exitCode = priorExitCode;
  }
});

test('nonzero owned child exit remains a cleanup failure', async () => {
  const priorExitCode = process.exitCode;
  const child = spawn(process.execPath, ['-e', 'process.on("SIGTERM", () => process.exit(7)); process.stdout.write("ready"); setInterval(() => {}, 1000)'], { detached: true, stdio: ['ignore', 'pipe', 'ignore'] });
  try {
    await new Promise((resolve) => child.stdout.once('data', resolve));
    const result = await cleanupFixture(undefined, [child], undefined, undefined, undefined);
    const status = result.steps.owned_children_stopped.statuses[0];
    assert.equal(status.state, 'exited');
    assert.equal(status.exit_code, 7);
    assert.equal(result.passed, false);
    assert.equal(process.exitCode, 1);
  } finally {
    try { if (child.exitCode === null && child.signalCode === null) process.kill(-child.pid, 'SIGKILL'); } catch { /* test-owned child */ }
    process.exitCode = priorExitCode;
  }
});

test('report diagnostics stay bounded and publish runs after cleanup', async () => {
  const priorExitCode = process.exitCode;
  const events = [];
  const errors = [];
  try {
    assert.deepEqual(safeRead('target rows', () => { throw new Error('query failed'); }, [], errors), []);
    assert.deepEqual(errors, [{ label: 'target rows', reason: 'diagnostic-read-failed', message: 'query failed' }]);
    const report = await publishAfterCleanup({
      cleanup: async () => { events.push('cleanup'); return { passed: true }; },
      buildReport: async () => { events.push('build'); throw new Error('hash failed'); },
      writeReport: async (value) => { events.push('publish'); assert.equal(value.fixture_cleanup.passed, true); },
    });
    assert.deepEqual(events, ['cleanup', 'build', 'publish']);
    assert.equal(report.failure.stage, 'report');
    assert.equal(process.exitCode, 1);
  } finally { process.exitCode = priorExitCode; }
});

for (const failed of ['browser', 'role', 'target']) {
  test(`${failed} cleanup failure keeps failure explicit and runs later cleanup`, async () => {
    const priorExitCode = process.exitCode;
    const calls = [];
    const action = (name) => () => { calls.push(name); if (name === failed) throw new Error('fixture failure'); };
    try {
      const result = await cleanupFixture(
        { close: action('browser') }, {}, action('port'), { cleanup: action('role') }, { cleanup: action('target') },
      );
      assert.equal(result.passed, false);
      assert.equal(process.exitCode, 1);
      assert.deepEqual(calls, ['browser', 'port', 'role', 'target']);
      const key = { browser: 'browser_closed', role: 'permission_role_removed', target: 'owned_targets_removed' }[failed];
      assert.deepEqual(result.steps[key], { passed: false, reason: 'fixture-cleanup-failed' });
    } finally { process.exitCode = priorExitCode; }
  });
}

test('returned cleanup failure sets process exit status', async () => {
  const priorExitCode = process.exitCode;
  try {
    process.exitCode = undefined;
    const result = await cleanupFixture(
      undefined, undefined, undefined,
      { cleanup: async () => ({ passed: false, reason: 'role retained' }) }, undefined,
    );
    assert.equal(result.passed, false);
    assert.equal(process.exitCode, 1);
  } finally { process.exitCode = priorExitCode; }
});

test('failed cleanup cannot be published as a passing report', async () => {
  const priorExitCode = process.exitCode;
  const published = [];
  try {
    process.exitCode = undefined;
    const report = await publishAfterCleanup({
      cleanup: async () => ({ passed: false, steps: { owned_children_stopped: { passed: false } } }),
      buildReport: async (fixtureCleanup) => ({ passed: true, fixture_cleanup: fixtureCleanup }),
      writeReport: async (value) => published.push(value),
    });
    assert.equal(process.exitCode, 1);
    assert.equal(report.passed, false);
    assert.equal(published[0].passed, false);
    assert.equal(published[0].fixture_cleanup.passed, false);
  } finally { process.exitCode = priorExitCode; }
});
