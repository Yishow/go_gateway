// Immutable pre-capacity binary: use actual Modbus reads, never inserted samples.
import assert from 'node:assert/strict';
import { createHash } from 'node:crypto';
import { copyFileSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { ROOT, sh, startProcess, witness, writeRegisters } from './lib.mjs';

const prior = JSON.parse(readFileSync(join(ROOT, 'docs/plans/studio-v2-write-groups/evidence-f/quality-sqlite.json')));
const priorWork = join('/tmp', prior.run_id);
const binary = join(priorWork, 'gw-ui');
const hash = createHash('sha256').update(readFileSync(binary)).digest('hex');
assert.equal(hash, prior.binary_sha256);
const work = `/tmp/gw-f-capacity-red-${Date.now()}`;
mkdirSync(work); writeFileSync(join(work, '.f-write-group-fixture'), 'f_write_group_fixture');
copyFileSync(binary, join(work, 'gw-ui'));
sh('sqlite3', [join(priorWork, 'gateway.db'), `.backup '${join(work, 'gateway.db')}'`]);
const children = [];
const request = async (path, body = {}) => {
  const response = await fetch(`http://127.0.0.1:3379${path}`, {
    method: 'POST', headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body), signal: AbortSignal.timeout(5000),
  });
  assert.equal(response.ok, true); return response.json();
};
const group = prior.groups[0];
let result;
try {
  children.push(startProcess(work, 'sim-a', join(priorWork, 'f-sim'), ['-port', '15020', '-registers', writeRegisters(work, 'regs-a', { 0: 215, 1: 1013, 2: 32, 3: 0, 4: 0, 5: 1, 6: 1 })]));
  children.push(startProcess(work, 'gateway', join(work, 'gw-ui'), [], {
    GATEWAY_DB_PATH: join(work, 'gateway.db'), PORT: '3377', F_FIXTURE_PORT: '3379',
    F_FIXTURE_START_AT: '2026-01-01T00:01:18Z', F_FIXTURE_GROUP_QUOTA_BYTES: '1',
  }));
  const deadline = Date.now() + 20000;
  while (Date.now() < deadline) {
    try { if ((await fetch('http://127.0.0.1:3379/state', { signal: AbortSignal.timeout(1000) })).ok) break; } catch { /* owned startup */ }
    await new Promise((resolve) => setTimeout(resolve, 100));
  }
  await request('/pause');
  const attempts = [];
  for (const id of ['quota-one', 'quota-two']) {
    const captured = await request('/poll', { group_id: group.id, point_ids: [group.members[0].point_id], acquisition_id: id });
    assert.equal(captured.indices.length, 1);
    attempts.push((await request('/release', { indices: captured.indices })).results[0]);
  }
  result = { test: 'FixtureQuotaRefusesNewAck', attempts, expected_second_reason: 'quota-hard-limit', passed: !attempts[1].accepted && attempts[1].reason === 'quota-hard-limit' };
  console.log(`${result.passed ? 'PASS' : 'FAIL'} ${result.test}: ${JSON.stringify(attempts)}`);
  process.exitCode = result.passed ? 0 : 1;
} finally {
  for (const child of children) { try { process.kill(-child.pid, 'SIGKILL'); } catch { /* owned child stopped */ } }
  writeFileSync(join(ROOT, 'docs/plans/studio-v2-write-groups/evidence-f/capacity-quota-red.json'), JSON.stringify(witness({
    phase: 'pre-capacity immutable binary', baseline_run_id: prior.run_id,
    baseline_source_sha: prior.source_sha, binary_sha256: hash,
    command: 'node scripts/tests/f_device_to_sql/capacity-red.mjs', group_id: group.id, result,
  })) + '\n');
}
