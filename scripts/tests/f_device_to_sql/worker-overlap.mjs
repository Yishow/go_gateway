// Two actual production pipelines contend for one owned durable store.
import assert from 'node:assert/strict';
import { createHash } from 'node:crypto';
import { mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import net from 'node:net';
import {
  ROOT, api, buildBinaries, cleanWork, configureDestination, configureDevicesPointsAndTags,
  createAndApplyGroup, activateDevices, openBrowser, startProcess, waitForGateway,
  waitForPortFree, witness, writeRegisters, PORT_A, PORT_B, sh,
} from './lib.mjs';
import { control, eventually, stamp } from './quality-lib.mjs';
import { recoveryTargets } from './recovery-targets.mjs';
import { cleanupFixture, publishAfterCleanup, safeRead, safeReadAsync } from './fixture-cleanup.mjs';
import { recoveryContext, mutation, saveApply, localOutboxSQL } from './recovery-lib.mjs';

const runNumber = Date.now(); const runId = `gw-f-workers-${runNumber}`;
const work = `/tmp/${runId}`; const evidence = join(ROOT, 'docs/plans/studio-v2-write-groups/evidence-f');
const secondaryPort = 3403; const secondaryFixturePort = 3405;
const children = {}; const results = {}; const gatewayEnv = {};
let targets; let browser; let page; let problems = []; let e; let groups; let gatewayStarted = false;
const kill = (name) => {
  const child = children[name];
  if (!child || child.exitCode !== null || child.signalCode !== null) return;
  try { process.kill(-child.pid, 'SIGTERM'); } catch { /* owned child stopped */ }
};
process.on('exit', () => { for (const name of Object.keys(children)) kill(name); });
const checkPortFree = async (port) => {
  const probe = net.createServer();
  try {
    await new Promise((resolve, reject) => { probe.once('error', reject); probe.listen(port, '127.0.0.1', resolve); });
  } finally {
    if (probe.listening) await new Promise((resolve, reject) => probe.close((error) => error ? reject(error) : resolve()));
  }
};
try {
  // Refuse occupied secondary ports before creating fixtures or starting children.
  for (const port of [secondaryPort, secondaryFixturePort]) await checkPortFree(port);
  await waitForPortFree(); cleanWork(work); mkdirSync(evidence, { recursive: true });
  writeFileSync(join(work, '.f-write-group-fixture'), 'f_write_group_fixture');
  buildBinaries(work, { tags: 'f_write_group_fixture' });
  targets = await recoveryTargets('sqlite', work, runNumber);
  for (const [name, port, registers] of [
    ['a', PORT_A, { 0: 215, 1: 1013, 2: 32, 3: 0, 4: 0, 5: 1, 6: 1 }],
    ['b', PORT_B, { 0: 187, 1: 777, 2: 0, 3: 0, 4: 0, 5: 2, 6: 0 }],
  ]) children[name] = startProcess(work, `sim-${name}`, join(work, 'f-sim'), ['-port', String(port), '-registers', writeRegisters(work, `regs-${name}`, registers)]);
  e = recoveryContext(work, children, targets, null, evidence, results, gatewayEnv);
  e.startGateway(); gatewayStarted = true; await waitForGateway();
  ({ browser, page, problems } = await openBrowser());
  e = recoveryContext(work, children, targets, page, evidence, results, gatewayEnv);
  await configureDevicesPointsAndTags(page); await configureDestination(page, targets.destination);
  await createAndApplyGroup(page, 'Line A', 'a'); await createAndApplyGroup(page, 'Line B', 'b');
  await activateDevices(page); await control('/pause'); await e.clock(0);
  const created = (await api('/studio-v2/workspace/write-groups')).groups;
  const a = await saveApply(created.find((group) => group.name === 'Line A'), { write_policy: { dedupe_capability: 'receipt' } });
  const connectorB = await mutation('/db-targets/connectors', 'POST', { name: 'Healthy overlap destination', kind: 'sqlite', connection_config: targets.bConfig, enabled: true });
  const b = await saveApply(created.find((group) => group.name === 'Line B'), {
    write_policy: { dedupe_capability: 'receipt' },
    destination: { ...created.find((group) => group.name === 'Line B').destination, connector_id: connectorB.id, connector_revision: connectorB.identity_revision, database: targets.database('b') },
  });
  groups = { a, b }; await e.close(10);
  await e.fault('target_commit_hold', a, true);
  await e.feed(18, a, 'exclusive-claim-original'); await e.close(20);
  const held = await e.reached('target_commit_hold', a); const primaryClaim = e.item(a, 10);
  assert.equal(primaryClaim.state, 'sending'); assert.equal(targets.rows('a').length, 1);
  assert.equal(e.receipts().some((row) => row.effect_key === primaryClaim.effect_key), false);

  // Separate process and worker incarnation, same real migrated gateway DB.
  children.secondary = startProcess(work, 'gateway-secondary', join(work, 'gw-ui'), [], {
    GATEWAY_DB_PATH: join(work, 'gateway.db'), PORT: String(secondaryPort),
    F_FIXTURE_PORT: String(secondaryFixturePort), F_FIXTURE_START_AT: stamp(20),
    F_FIXTURE_WORKER_NODE: 'secondary',
  });
  const secondaryState = await eventually('second actual worker is running', async () => {
    assert.equal(children.secondary.exitCode, null, 'owned secondary has not exited');
    assert.equal(children.secondary.signalCode, null, 'owned secondary has not been signaled');
    try {
      const owners = [secondaryPort, secondaryFixturePort].map((port) => sh('lsof', ['-nP', `-iTCP:${port}`, '-sTCP:LISTEN', '-t']).split(/\s+/).map(Number));
      if (!owners.every((pids) => pids.length === 1 && pids[0] === children.secondary.pid)) return undefined;
      const response = await fetch(`http://127.0.0.1:${secondaryPort}/api/v1/datalink/runtime/status`, { signal: AbortSignal.timeout(2000) });
      return { status: response.status, body: await response.json(), listener_pid: children.secondary.pid };
    } catch { return undefined; }
  }, (value) => value?.status === 200 && value.body.data.running);
  await e.feed(28, a, 'same-partition-held-successor'); await e.feed(28, b, 'other-worker-healthy-partition'); await e.close(30);
  await e.committed(b, 20);
  const healthyClaim = e.item(b, 20);
  assert.notEqual(healthyClaim.claim_owner, primaryClaim.claim_owner, 'the other production worker actually delivered the healthy partition');
  assert.equal(e.item(a, 10).claim_owner, primaryClaim.claim_owner); assert.equal(e.item(a, 10).state, 'sending');
  assert.equal(e.item(a, 20).state, 'pending'); assert.equal(targets.rows('a').length, 1);
  const exclusive = (await e.state()).faults.find((fault) => fault.group_id === a.id);
  assert.equal(exclusive.commits, 1); assert.equal(exclusive.reached, 1);
  await e.fault('target_commit_hold', a, false);
  await e.committed(a, 10); await e.committed(a, 20);
  assert.equal(targets.rows('a').length, 2); assert.equal(targets.rows('b').length, 1);
  e.check('ActualWorkerOverlapExclusiveClaimAndHealthyPartition', {
    processes: { primary_pid: children.gateway.pid, secondary_pid: children.secondary.pid }, secondary_runtime: secondaryState,
    primary_claim: primaryClaim, secondary_healthy_claim: healthyClaim, barrier: held, during_overlap: exclusive,
    confirmed: [e.proof(a, 10, 'a'), e.proof(a, 20, 'a'), e.proof(b, 20, 'b')],
  });
  kill('secondary');
  await eventually('owned second worker stopped', async () => {
    try { await fetch(`http://127.0.0.1:${secondaryPort}/studio/v2`, { signal: AbortSignal.timeout(1000) }); return false; } catch { return true; }
  }, Boolean);

  gatewayEnv.F_FIXTURE_MAX_RETRIES = '2'; await e.restart(40);
  await e.feed(48, a, 'capped-retry-failing-target'); await e.feed(48, b, 'capped-retry-healthy');
  await targets.disconnect(); await e.close(50); await e.committed(b, 40);
  const exhausted = await eventually('two actual transient attempts reach configured retry cap', () => e.item(a, 40),
    (row) => row?.state === 'blocked' && row.retry_count === 2, 40000);
  assert.equal(exhausted.last_error_code, 'target-unavailable');
  await e.feed(58, a, 'capped-retry-retained-successor'); await e.feed(58, b, 'capped-retry-healthy-successor'); await e.close(60);
  await e.committed(b, 50);
  assert.equal(e.item(a, 50).state, 'pending'); assert.equal(targets.rows('a').length, 2); assert.equal(targets.rows('b').length, 3);
  await targets.reconnect();
  e.check('ActualRetryExhaustionRetainsPayloadAndHealthyScope', {
    tagged_startup_max_retries: 2, exhausted, retained_successor: e.item(a, 50),
    healthy: [e.proof(b, 40, 'b'), e.proof(b, 50, 'b')], ui: await e.ui(a, 'retry-exhausted'),
  });
  assert.deepEqual(problems.filter((problem) => /pageerror/.test(problem)), []);
} catch (error) {
  results.failure = { passed: false, message: error.message }; process.exitCode = 1; console.error(error);
  await page?.screenshot({ path: join(evidence, 'worker-overlap-failure.png'), fullPage: true }).catch(() => {});
} finally {
  const reportErrors = [];
  const snapshot = {
    binaryHash: safeRead('gateway binary hash', () => createHash('sha256').update(readFileSync(join(work, 'gw-ui'))).digest('hex'), null, reportErrors),
    actualRows: {
      a: safeRead('target rows a', () => targets?.rows('a'), [], reportErrors),
      b: safeRead('target rows b', () => targets?.rows('b'), [], reportErrors),
    },
    outbox: await safeReadAsync('local outbox', () => e?.outbox(), [], reportErrors),
    receipts: await safeReadAsync('local receipts', () => e?.receipts(), [], reportErrors),
  };
  await publishAfterCleanup({
    cleanup: () => cleanupFixture(browser, children, gatewayStarted ? waitForPortFree : undefined, undefined, targets),
    buildReport: async (fixtureCleanup) => witness({
      run_id: runId, binary_sha256: snapshot.binaryHash, groups, results,
      commands: ['make sync-frontend-static', 'go build -tags f_write_group_fixture ./cmd/test_ui', 'node scripts/tests/f_device_to_sql/worker-overlap.mjs'],
      actual_rows: snapshot.actualRows, target_queries: targets?.queries ?? null, local_outbox_query: localOutboxSQL,
      outbox: snapshot.outbox, receipts: snapshot.receipts, ui_problems: problems, report_errors: reportErrors,
      fixture_cleanup: fixtureCleanup, retained_evidence: { workdir: work, paths: [join(work, 'gateway.db'), join(work, '*.log')], external_sqlite_targets: [targets?.database('a'), targets?.database('b')], reason: 'diagnostic-evidence-retained' },
      passed: !process.exitCode && reportErrors.length === 0,
      limits: ['test-only retry cap; production default remains unlimited', 'two loopback workers on an owned SQLite store, not field deployment'],
    }),
    writeReport: (report) => writeFileSync(join(evidence, 'worker-overlap-sqlite.json'), `${JSON.stringify(report)}\n`),
  });
}
