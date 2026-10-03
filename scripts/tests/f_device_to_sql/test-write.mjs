// ConfirmedTestWriteCleanup through the actual UI, SQL target and durable ledger.
import assert from 'node:assert/strict';
import { createHash } from 'node:crypto';
import { mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import {
  ROOT, api, buildBinaries, cleanWork, configureDestination, configureDevicesPointsAndTags,
  createAndApplyGroup, activateDevices, openBrowser, startProcess, waitForGateway,
  waitForPortFree, witness, writeRegisters, PORT_A, PORT_B, sh,
} from './lib.mjs';
import { control, eventually, query } from './quality-lib.mjs';
import { recoveryTargets } from './recovery-targets.mjs';
import { recoveryContext, saveApply } from './recovery-lib.mjs';
import { assertResult, confirm, operationSQL, permissionTargets, retained, safePreview, savedClone, screenshotOutcome, uiConfirm, uiPreview } from './test-write-lib.mjs';
import { cleanupFixture, publishAfterCleanup, safeRead } from './fixture-cleanup.mjs';

const kind = process.env.F_SQL_KIND ?? 'sqlite'; assert.ok(['sqlite', 'postgres'].includes(kind));
const runNumber = Date.now(); const runId = `gw-f-test-write-${kind}-${runNumber}`;
const work = `/tmp/${runId}`; const evidence = join(ROOT, 'docs/plans/studio-v2-write-groups/evidence-f');
const children = {}; const results = {};
let targets; let browser; let page; let problems = []; let e; let groups; let permissions; let gatewayStarted = false;
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
  targets = await recoveryTargets(kind, work, runNumber);
  for (const [name, port, registers] of [
    ['a', PORT_A, { 0: 215, 1: 1013, 2: 32, 3: 0, 4: 0, 5: 1, 6: 1 }],
    ['b', PORT_B, { 0: 187, 1: 777, 2: 0, 3: 0, 4: 0, 5: 2, 6: 0 }],
  ]) children[name] = startProcess(work, `sim-${name}`, join(work, 'f-sim'), ['-port', String(port), '-registers', writeRegisters(work, `regs-${name}`, registers)]);
  e = recoveryContext(work, children, targets, null, evidence, results);
  e.startGateway(); gatewayStarted = true; await waitForGateway();
  ({ browser, page, problems } = await openBrowser());
  e = recoveryContext(work, children, targets, page, evidence, results);
  if (kind === 'postgres') permissions = permissionTargets(targets, runNumber);
  await configureDevicesPointsAndTags(page); await configureDestination(page, permissions?.destination ?? targets.destination);
  await createAndApplyGroup(page, 'Line A', 'a'); await createAndApplyGroup(page, 'Line B', 'b');
  await activateDevices(page); await control('/pause'); await e.clock(0);
  const created = (await api('/studio-v2/workspace/write-groups')).groups;
  const a = await saveApply(created.find((group) => group.name === 'Line A'), { write_policy: { dedupe_capability: 'receipt' } });
  const b = await saveApply(created.find((group) => group.name === 'Line B'), { write_policy: { dedupe_capability: 'receipt' } });
  groups = { a, b }; await e.close(10);
  await e.feed(18, a, 'actual-neighbor-a'); await e.feed(18, b, 'actual-neighbor-b'); await e.close(20);
  await e.committed(a, 10); await e.committed(b, 10);
  const neighbors = targets.rows('a'); assert.equal(neighbors.length, 2);
  const preview = await uiPreview(page, a, () => targets.rows('a'));
  const op = await uiConfirm(page, a); assertResult(op, 'written_verified', 'cleaned');
  assert.deepEqual(targets.rows('a'), neighbors);
  e.check('ActualPreviewConfirmTypedReadbackOwnedCleanup', {
    preview: safePreview(preview), operation: op, neighbors, ui: await screenshotOutcome(page, evidence, kind, 'cleaned'),
    retained: await retained(a, preview, op, () => targets.rows('a')),
  });

  if (kind === 'sqlite') {
    const interrupted = await uiPreview(page, a, () => targets.rows('a'));
    const trigger = `CREATE TRIGGER f_test_cleanup_hold BEFORE DELETE ON readings WHEN OLD.line='${interrupted.owner_value}' BEGIN SELECT f_write_group_fixture_closure_hold('${a.id}'); END`;
    sh('sqlite3', [targets.database('a'), trigger]);
    await e.fault('closure_hold', a, true);
    const pending = confirm(a, interrupted).catch((error) => ({ interrupted: error.message }));
    const barrier = await e.reached('closure_hold', a);
    const beforeKill = targets.rows('a'); assert.equal(beforeKill.length, neighbors.length + 1);
    assert.equal(beforeKill.filter((row) => row.line === interrupted.owner_value).length, 1);
    const ownedRow = beforeKill.find((row) => row.line === interrupted.owner_value);
    assert.equal(ownedRow.pressure, '9007199254740993'); assert.equal(ownedRow.batch, '9007199254740993');
    const localSQL = operationSQL.replace('?', `'${interrupted.operation_id}'`);
    const before = query(join(work, 'gateway.db'), localSQL);
    assert.equal(JSON.parse(before[0].detail).phase, 'cleaning');
    await e.restart(20); await pending;
    assert.deepEqual(targets.rows('a'), beforeKill, 'SIGKILL rolls back only the uncommitted cleanup');
    sh('sqlite3', [targets.database('a'), 'DROP TRIGGER f_test_cleanup_hold']);
    await e.fault('closure_hold', a, false);
    const active = await confirm(a, interrupted); assert.equal(active.status, 202);
    assert.deepEqual(targets.rows('a'), beforeKill, 'a live operation lease prevents a new write');
    const adopted = await eventually('real three-minute operation lease expires and cleanup resumes', () => confirm(a, interrupted),
      (response) => response.status === 200, 220000);
    assertResult(adopted.body.data, 'written_verified', 'cleaned');
    assert.deepEqual(targets.rows('a'), neighbors);
    // Keep the actual UI preview. A confirm after adoption reads the retained
    // result and displays it without a second target effect.
    const uiReplay = await uiConfirm(page, a); assert.deepEqual(uiReplay, adopted.body.data);
    const saved = await retained(a, interrupted, adopted.body.data, () => targets.rows('a'));
    e.check('CleanupKillRestartSameOperationWithoutSecondWrite', {
      preview: safePreview(interrupted), fixture_ddl: trigger, barrier, ledger_before_kill: before,
      target_before_kill: beforeKill, in_progress_status: active.status, operation: adopted.body.data, retained: saved,
      ui: await screenshotOutcome(page, evidence, kind, 'restarted-cleaned'),
    });
    results.PermissionCases = { passed: null, reason: 'SQLite has no PostgreSQL SELECT/DELETE grants; actual permission cases run on owned PostgreSQL' };
  } else {
    for (const permission of ['no_select', 'no_delete']) {
      // Change grants on our run-owned non-superuser role, preserving the
      // saved connector identity and its production ownership guards.
      permissions.restrict(permission);
      const connector = { id: a.destination.connector_id, identity_revision: a.destination.connector_revision };
      const group = await savedClone(a, connector, `Permission ${permission}`);
      const limitedPreview = await uiPreview(page, group, () => targets.rows('a'));
      const limited = await uiConfirm(page, group);
      results[`ActualPermission_${permission}`] = { passed: false, stage: 'received_operation', preview: safePreview(limitedPreview), operation: limited };
      assertResult(limited, permission === 'no_select' ? 'written_unverified' : 'written_verified', 'failed');
      assert.equal(limited.cleanup_reason, 'cleanup-denied');
      const withOwned = targets.rows('a'); assert.equal(withOwned.length, neighbors.length + 1);
      assert.deepEqual(withOwned.filter((row) => row.line !== limitedPreview.owner_value), neighbors);
      const saved = await retained(group, limitedPreview, limited, () => targets.rows('a'));
      e.check(`ActualPermission_${permission}`, {
        group, role: permissions.role, revoked_privilege: permission, preview: safePreview(limitedPreview), operation: limited, owned_rows: withOwned, retained: saved,
        ui: await screenshotOutcome(page, evidence, kind, permission),
      });
      permissions.removeOwned(limitedPreview.owner_value); assert.deepEqual(targets.rows('a'), neighbors);
    }
  }
  assert.deepEqual(problems.filter((problem) => /pageerror/.test(problem)), []);
} catch (error) {
  results.failure = { passed: false, message: error.message, setup_cleanup: error.setupCleanup };
  await page?.screenshot({ path: join(evidence, `test-write-${kind}-failure.png`), fullPage: true }).catch(() => {});
  process.exitCode = 1; console.error(error);
} finally {
  const reportErrors = [];
  const snapshot = {
    binaryHash: safeRead('gateway binary hash', () => createHash('sha256').update(readFileSync(join(work, 'gw-ui'))).digest('hex'), null, reportErrors),
    actualRows: safeRead('target rows a', () => targets?.rows('a'), [], reportErrors),
  };
  const permissionsForCleanup = permissions ? {
    cleanup: async () => {
      await permissions.cleanup();
      return { removed: { role: permissions.role } };
    },
  } : undefined;
  await publishAfterCleanup({
    cleanup: () => cleanupFixture(browser, children, gatewayStarted ? waitForPortFree : undefined, permissionsForCleanup, targets),
    buildReport: async (fixtureCleanup) => witness({
      run_id: runId, binary_sha256: snapshot.binaryHash,
      commands: ['make sync-frontend-static', 'go build -tags f_write_group_fixture ./cmd/test_ui', 'go build ./cmd/f_modbus_simulator', `F_SQL_KIND=${kind} node scripts/tests/f_device_to_sql/test-write.mjs`],
      groups, results, actual_rows: snapshot.actualRows, target_queries: targets?.queries ?? null, local_operation_query: operationSQL,
      ui_problems: problems, report_errors: reportErrors, fixture_cleanup: fixtureCleanup,
      retained_evidence: { workdir: work, paths: [join(work, 'gateway.db'), join(work, '*.log')], external_sqlite_targets: kind === 'sqlite' ? [targets?.database('a'), targets?.database('b'), targets?.database('other')] : [], reason: 'diagnostic-evidence-retained' },
      passed: !process.exitCode && reportErrors.length === 0,
      coverage_scope: kind === 'sqlite' ? 'SQLite normal/restart/same-operation cases; permission counterparts require the separate PostgreSQL witness' : 'PostgreSQL normal/no-SELECT/no-DELETE cases; restart counterpart uses the separate SQLite witness',
      limits: ['disposable databases and loopback simulator only; no field acceptance'],
    }),
    writeReport: (report) => writeFileSync(join(evidence, `test-write-${kind}.json`), `${JSON.stringify(report)}\n`),
  });
}
