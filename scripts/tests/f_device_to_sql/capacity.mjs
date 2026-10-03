// CapacityPoisonAndConcurrentEdit on actual production wiring and owned SQLite.
import assert from 'node:assert/strict';
import { createHash } from 'node:crypto';
import { mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import {
  ROOT, api, buildBinaries, cleanWork, configureDestination, configureDevicesPointsAndTags,
  createAndApplyGroup, activateDevices, openBrowser, startProcess, waitForGateway,
  waitForPortFree, witness, writeRegisters, PORT_A, PORT_B,
} from './lib.mjs';
import { control } from './quality-lib.mjs';
import { recoveryTargets } from './recovery-targets.mjs';
import { recoveryContext, mutation, saveApply, localOutboxSQL, localSamplesSQL, localReceiptsSQL, localCheckpointSQL } from './recovery-lib.mjs';
import { runCapacityCases } from './capacity-cases.mjs';
import { runConcurrentEditCases } from './capacity-cas.mjs';
import { cleanupFixture, publishAfterCleanup, safeRead, safeReadAsync } from './fixture-cleanup.mjs';

const runNumber = Date.now(); const runId = `gw-f-capacity-${runNumber}`;
const work = `/tmp/${runId}`;
const evidence = join(ROOT, 'docs/plans/studio-v2-write-groups/evidence-f');
const children = {}; const results = {}; const gatewayEnv = {};
let targets; let browser; let page; let problems = []; let e; let groups; let capacityWitness; let casWitness; let gatewayStarted = false;
const kill = (name) => {
  const child = children[name];
  if (!child || child.exitCode !== null || child.signalCode !== null) return;
  try { process.kill(-child.pid, 'SIGTERM'); } catch { /* owned child stopped */ }
};
process.on('exit', () => { for (const name of Object.keys(children)) kill(name); });
try {
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
  let a = created.find((group) => group.name === 'Line A'); let b = created.find((group) => group.name === 'Line B');
  const connectorB = await mutation('/db-targets/connectors', 'POST', { name: 'Healthy capacity destination', kind: 'sqlite', connection_config: targets.bConfig, enabled: true });
  a = await saveApply(a, { write_policy: { dedupe_capability: 'receipt' } });
  b = await saveApply(b, {
    write_policy: { dedupe_capability: 'receipt' },
    destination: { ...b.destination, connector_id: connectorB.id, connector_revision: connectorB.identity_revision, database: targets.database('b') },
  });
  await e.close(10); groups = { a, b };
  capacityWitness = await runCapacityCases(e, groups, work, gatewayEnv);
  casWitness = await runConcurrentEditCases(e, groups);
  assert.deepEqual(problems.filter((problem) => /pageerror/.test(problem)), []);
} catch (error) {
  results.failure = { passed: false, message: error.message };
  await page?.screenshot({ path: join(evidence, 'capacity-sqlite-failure.png'), fullPage: true }).catch(() => {});
  process.exitCode = 1; console.error(error);
} finally {
  const reportErrors = [];
  const snapshot = {
    binaryHash: safeRead('gateway binary hash', () => createHash('sha256').update(readFileSync(join(work, 'gw-ui'))).digest('hex'), null, reportErrors),
    actualRows: {
      a: safeRead('target rows a', () => targets?.rows('a'), [], reportErrors),
      b: safeRead('target rows b', () => targets?.rows('b'), [], reportErrors),
    },
    samples: await safeReadAsync('local samples', () => e?.samples(), [], reportErrors),
    outbox: await safeReadAsync('local outbox', () => e?.outbox(), [], reportErrors),
    checkpoints: await safeReadAsync('local checkpoints', () => e?.checkpoints(), [], reportErrors),
    receipts: await safeReadAsync('local receipts', () => e?.receipts(), [], reportErrors),
  };
  await publishAfterCleanup({
    cleanup: () => cleanupFixture(browser, children, gatewayStarted ? waitForPortFree : undefined, undefined, targets),
    buildReport: async (fixtureCleanup) => witness({
      run_id: runId, binary_sha256: snapshot.binaryHash,
      commands: ['make sync-frontend-static', 'go build -tags f_write_group_fixture ./cmd/test_ui', 'go build ./cmd/f_modbus_simulator', 'node scripts/tests/f_device_to_sql/capacity.mjs'],
      groups, results, capacity: capacityWitness, concurrent_edits: casWitness,
      actual_rows: snapshot.actualRows, target_queries: targets?.queries ?? null,
      local_queries: [localSamplesSQL, localOutboxSQL, localCheckpointSQL, localReceiptsSQL],
      samples: snapshot.samples, outbox: snapshot.outbox, checkpoints: snapshot.checkpoints, receipts: snapshot.receipts, ui_problems: problems,
      report_errors: reportErrors, fixture_cleanup: fixtureCleanup,
      retained_evidence: { workdir: work, paths: [join(work, 'gateway.db'), join(work, '*.log')], external_sqlite_targets: [targets?.database('a'), targets?.database('b')], reason: 'diagnostic-evidence-retained' },
      passed: !process.exitCode && reportErrors.length === 0 && Boolean(capacityWitness && casWitness),
      limits: ['dispatcher overlap and explicit retry exhaustion require separate actual witness', 'no field/production DB acceptance'],
    }),
    writeReport: (report) => writeFileSync(join(evidence, 'capacity-sqlite.json'), `${JSON.stringify(report)}\n`),
  });
}
