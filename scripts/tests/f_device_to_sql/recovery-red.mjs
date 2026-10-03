// Use the preserved pre-F2.2 binary, rather than reverting anyone's source.
import assert from 'node:assert/strict';
import { createHash } from 'node:crypto';
import { copyFileSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { ROOT, sh, startProcess, witness } from './lib.mjs';

const priorPath = join(ROOT, 'docs/plans/studio-v2-write-groups/evidence-f/quality-sqlite.json');
const prior = JSON.parse(readFileSync(priorPath, 'utf8'));
const priorWork = join('/tmp', prior.run_id);
const binary = join(priorWork, 'gw-ui');
const hash = createHash('sha256').update(readFileSync(binary)).digest('hex');
assert.equal(hash, prior.binary_sha256);
const work = `/tmp/gw-f-recovery-red-${Date.now()}`;
mkdirSync(work);
writeFileSync(join(work, '.f-write-group-fixture'), 'f_write_group_fixture');
copyFileSync(binary, join(work, 'gw-ui'));
sh('sqlite3', [join(priorWork, 'gateway.db'), `.backup '${join(work, 'gateway.db')}'`]);
const group = prior.groups[0];
const ports = { gateway: 3367, control: 3371 };
const child = startProcess(work, 'gateway', join(work, 'gw-ui'), [], {
  GATEWAY_DB_PATH: join(work, 'gateway.db'), PORT: String(ports.gateway),
  F_FIXTURE_PORT: String(ports.control), F_FIXTURE_START_AT: '2026-01-01T00:00:00Z',
});
let result;
try {
  const deadline = Date.now() + 20000;
  while (Date.now() < deadline) {
    try {
      if ((await fetch(`http://127.0.0.1:${ports.control}/state`)).ok) break;
    } catch { /* wait for the owned process */ }
    await new Promise((resolve) => setTimeout(resolve, 100));
  }
  const response = await fetch(`http://127.0.0.1:${ports.control}/fault`, {
    method: 'POST', headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ kind: 'target_commit_response_lost', group_id: group.id, enabled: true }),
  });
  result = { test: 'ProductionFaultControlAvailable', expected_status: 200, actual_status: response.status, body: await response.text(), group_id: group.id, passed: response.status === 200 };
  console.log(`${result.passed ? 'PASS' : 'FAIL'} ${result.test}: expected=200 actual=${response.status}`);
  process.exitCode = result.passed ? 0 : 1;
} finally {
  try { process.kill(-child.pid, 'SIGKILL'); } catch { /* owned child already stopped */ }
  writeFileSync(join(ROOT, 'docs/plans/studio-v2-write-groups/evidence-f/recovery-control-red.json'), JSON.stringify(witness({
    phase: 'pre-F2.2 immutable binary', baseline_run_id: prior.run_id, binary_sha256: hash,
    baseline_source_sha: prior.source_sha, baseline_dirty_worktree: prior.dirty_worktree,
    command: 'node scripts/tests/f_device_to_sql/recovery-red.mjs', result,
  })) + '\n');
}
