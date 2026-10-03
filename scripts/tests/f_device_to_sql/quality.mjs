// QualityBoundaryMatrix: actual UI configuration, actual Modbus acquisition,
// production typed runtime/journal/closure/sender and independent SQL queries.
import assert from 'node:assert/strict';
import { createHash } from 'node:crypto';
import { mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import {
  ROOT, api, buildBinaries, cleanWork, configureDestination, configureDevicesPointsAndTags,
  createAndApplyGroup, activateDevices, GW_PORT, openBrowser, startProcess,
  waitForGateway, waitForPortFree, witness, writeRegisters, PORT_A, PORT_B, sh, boundedFetch,
} from './lib.mjs';
import { control, eventually, FIXTURE_PORT, groupMutation, query, showDelivery, stamp } from './quality-lib.mjs';
import { cleanupFixture, publishAfterCleanup, safeRead, safeReadAsync } from './fixture-cleanup.mjs';

const runId = `gw-f-quality-${Date.now()}`;
const work = `/tmp/${runId}`;
const evidence = join(ROOT, 'docs/plans/studio-v2-write-groups/evidence-f');
const target = join(work, 'destination.db');
const gatewayDB = join(work, 'gateway.db');
const results = {};
const children = {};
let browser; let page; let problems = []; let groups; let gatewayStarted = false;
const commands = ['make sync-frontend-static', 'go build -tags f_write_group_fixture ./cmd/test_ui', 'go build ./cmd/f_modbus_simulator', 'node scripts/tests/f_device_to_sql/quality.mjs'];
const check = (name, detail) => { results[name] = { passed: true, ...detail }; console.log(`PASS ${name}`); };
const kill = (name) => {
  const child = children[name];
  if (!child || child.exitCode !== null || child.signalCode !== null) return;
  try { process.kill(-child.pid, 'SIGTERM'); } catch { /* owned child already stopped */ }
};
process.on('exit', () => { for (const name of Object.keys(children)) kill(name); });

const registersA = { 0: 321, 1: 1013, 2: 32, 3: 0, 4: 0, 5: 1, 6: 1 };
const registersB = { 0: 187, 1: 777, 2: 0, 3: 0, 4: 0, 5: 2, 6: 0 };
const delivery = (group) => api(`/studio-v2/workspace/write-groups/${group.id}/delivery`);
const targetSQL = 'SELECT rowid, line, temperature, pressure, running, CAST(batch AS TEXT) AS batch, prov FROM readings ORDER BY rowid';
const rows = () => query(target, targetSQL);
const setClock = (seconds) => control('/clock', { at: stamp(seconds) });
const poll = (group, acquisition, pointIDs) => control('/poll', {
  group_id: group.id, acquisition_id: acquisition, ...(pointIDs ? { point_ids: pointIDs } : {}),
});
const release = async (capture) => {
  const released = await control('/release', { indices: capture.indices });
  assert.equal(released.results.length, capture.indices.length);
  assert.equal(released.results.every((result) => result.accepted), true, JSON.stringify(released));
  return released;
};
const captureAt = async (seconds, group, identity, pointIDs) => {
  await setClock(seconds);
  const captured = await poll(group, identity, pointIDs);
  assert.equal(captured.indices.length, pointIDs?.length ?? group.members.length, 'every actual polled member reaches the typed capture sink');
  return captured;
};
const closeAt = async (seconds) => { await setClock(seconds); await control('/tick'); };
const bucketSQL = 'SELECT group_id, group_revision, bucket_start, kind, reason, members FROM wg_delivery_buckets ORDER BY bucket_start, group_id';

try {
  await waitForPortFree();
  cleanWork(work);
  writeFileSync(join(work, '.f-write-group-fixture'), 'f_write_group_fixture');
  mkdirSync(evidence, { recursive: true });
  buildBinaries(work, { tags: 'f_write_group_fixture' });
  // Empty, verified nullable fixture. The harness never inserts a sample row.
  sh('sqlite3', [target, 'CREATE TABLE readings (line TEXT, temperature INTEGER, pressure INTEGER, running INTEGER, batch INTEGER, prov TEXT);']);
  const fileA = writeRegisters(work, 'regs-a', registersA);
  const fileB = writeRegisters(work, 'regs-b', registersB);
  children.a = startProcess(work, 'sim-a', join(work, 'f-sim'), ['-port', String(PORT_A), '-registers', fileA]);
  children.b = startProcess(work, 'sim-b', join(work, 'f-sim'), ['-port', String(PORT_B), '-registers', fileB]);
  children.gateway = startProcess(work, 'gateway', join(work, 'gw-ui'), [], {
    GATEWAY_DB_PATH: gatewayDB, PORT: GW_PORT, F_FIXTURE_PORT: FIXTURE_PORT,
    F_FIXTURE_START_AT: '2025-12-31T23:59:50Z',
  });
  gatewayStarted = true;
  await waitForGateway();
  ({ browser, page, problems } = await openBrowser());
  await configureDevicesPointsAndTags(page);
  await configureDestination(page, { kind: 'sqlite', path: target, table: 'readings' });
  await createAndApplyGroup(page, 'Line A', 'a');
  await createAndApplyGroup(page, 'Line B', 'b');
  await activateDevices(page);
  await control('/pause');
  groups = (await api('/studio-v2/workspace/write-groups')).groups;
  let a = groups.find((group) => group.name === 'Line A');
  const b = groups.find((group) => group.name === 'Line B');
  assert.ok(a && b);
  assert.equal(rows().length, 0, 'capture is explicitly not durable ACK or SQL');

  const newer = await captureAt(8, a, 'order-eight');
  const delayed = await captureAt(8, a, 'late-eight');
  writeFileSync(fileA, JSON.stringify({ ...registersA, 0: 123 }));
  // The simulator reload is an acquisition barrier; the chosen values are
  // checked in the final SQL, never injected as typed envelopes.
  await new Promise((resolve) => setTimeout(resolve, 700));
  const older = await captureAt(3, a, 'order-three');
  await setClock(8);
  await release(newer);
  await release(older);
  const journalBefore = query(gatewayDB, 'SELECT COUNT(*) AS n FROM wg_delivery_samples')[0].n;
  await release(newer);
  assert.equal(query(gatewayDB, 'SELECT COUNT(*) AS n FROM wg_delivery_samples')[0].n, journalBefore);
  check('DuplicateNoExtraJournalEffect', { journal_before: journalBefore });
  const conflict = await captureAt(8, a, 'order-eight');
  const refused = await control('/release', { indices: conflict.indices });
  assert.ok(refused.results.some((result) => /conflict/.test(result.reason ?? '')), 'same acquisition identity with a different actual device value must conflict');
  check('ConflictingActualReadIdentity', { release: refused });

  const bTemp = b.members.find((member) => member.target_column === 'temperature').point_id;
  await release(await captureAt(8, b, 'missing-eight', [bTemp]));
  await closeAt(10);
  const first = await eventually('first exact SQL row', rows, (value) => value.length === 1);
  assert.equal(first[0].temperature, 321);
  assert.equal(String(first[0].batch), '9007199254740993');
  assert.equal(JSON.parse(first[0].prov).every((member) => member.observed_at === stamp(8).replace('.000Z', 'Z')), true);
  check('OutOfOrderNewestObservationWins', { rows: first });
  const late = await control('/release', { indices: delayed.indices });
  assert.equal(late.results.every((result) => result.reason === 'late-after-close'), true);
  assert.equal(rows().length, 1);
  check('LateNeverRewritesClosedSQL', { release: late });

  await release(await captureAt(10, a, 'boundary-ten'));
  kill('b');
  await new Promise((resolve) => setTimeout(resolve, 300));
  await release(await captureAt(18, b, 'bad-eighteen'));
  await closeAt(20);
  await eventually('next boundary SQL row', rows, (value) => value.length === 2);
  check('BoundaryTenBelongsToNextBucket', { rows: rows() });

  a = (await groupMutation(a, { members: a.members.map((member) => ({ ...member, max_age_seconds: 2 })) })).group;
  await closeAt(30); // Both groups entirely silent in [20,30).
  await release(await captureAt(33, a, 'stale-thirty-three'));
  await closeAt(40);
  assert.equal(rows().length, 2, 'stale and silent buckets never carry a previous value');
  const buckets = query(gatewayDB, bucketSQL);
  const causesFor = (group, start) => {
    const found = buckets.find((bucket) => bucket.group_id === group.id && bucket.group_revision === group.applied_revision && Date.parse(bucket.bucket_start) === Date.parse(stamp(start)));
    assert.ok(found, `durable bucket must exist for ${group.name} revision ${group.applied_revision} at ${stamp(start)}`);
    return found;
  };
  assert.equal(causesFor(b, 0).kind, 'skipped');
  assert.ok(JSON.parse(causesFor(b, 0).members).some((member) => member.status === 'missing'));
  assert.equal(causesFor(b, 10).kind, 'skipped');
  assert.equal(JSON.parse(causesFor(b, 10).members).every((member) => member.status === 'bad'), true);
  assert.equal(causesFor(b, 20).kind, 'no_data');
  assert.equal(causesFor(a, 30).kind, 'skipped');
  assert.equal(JSON.parse(causesFor(a, 30).members).every((member) => member.status === 'stale'), true);
  check('DefaultMissingBadStaleAndEntirelySilentWriteNoRows', { buckets });

  a = (await groupMutation(a, {
    row_policy: { ...a.row_policy, incomplete_policy: 'partial' },
    members: a.members.map((member) => ({ ...member, required: member.target_column === 'temperature', max_age_seconds: 10 })),
  })).group;
  await setClock(50);
  const aTemp = a.members.find((member) => member.target_column === 'temperature').point_id;
  await release(await captureAt(58, a, 'partial-fifty-eight', [aTemp]));
  await closeAt(60);
  const partialRows = await eventually('actual partial SQL row', rows, (value) => value.length === 3);
  const partial = partialRows.at(-1);
  assert.equal(partial.temperature, 123);
  for (const column of ['pressure', 'running', 'batch']) assert.equal(partial[column], null);
  assert.equal(JSON.parse(partial.prov).filter((member) => member.status === 'missing').length, 3);
  check('ExplicitPartialUsesActualNULLAndReason', { row: partial });

  const viewA = await delivery(a);
  const viewB = await delivery(b);
  assert.ok(viewA.recent_bucket_issues.some((issue) => issue.causes.includes('stale')));
  for (const reason of ['missing', 'bad', 'no_data']) assert.ok(viewB.recent_bucket_issues.some((issue) => issue.causes.includes(reason)));
  const uiA = await showDelivery(page, a, viewA, join(evidence, 'quality-line-a.png'));
  const uiB = await showDelivery(page, b, viewB, join(evidence, 'quality-line-b.png'));
  check('UIAndDurableSQLDescribeSameUnwrittenBuckets', { delivery: { a: viewA, b: viewB }, ui: { a: uiA, b: uiB } });
} catch (error) {
  results.failure = { passed: false, message: error.message };
  await page?.screenshot({ path: join(evidence, 'quality-failure.png'), fullPage: true }).catch(() => {});
  process.exitCode = 1;
  console.error(error);
} finally {
  const reportErrors = [];
  const snapshot = {
    state: await safeReadAsync('fixture state', async () => {
      const response = await boundedFetch(`http://127.0.0.1:${FIXTURE_PORT}/state`, {}, 3000);
      return response.json();
    }, null, reportErrors),
    binaryHash: safeRead('gateway binary hash', () => createHash('sha256').update(readFileSync(join(work, 'gw-ui'))).digest('hex'), null, reportErrors),
    rows: safeRead('target rows', rows, [], reportErrors),
    buckets: safeRead('gateway buckets', () => query(gatewayDB, bucketSQL), [], reportErrors),
  };
  await publishAfterCleanup({
    cleanup: () => cleanupFixture(browser, children, gatewayStarted ? waitForPortFree : undefined, undefined, undefined),
    buildReport: async (fixtureCleanup) => witness({
      run_id: runId, commands, binary_sha256: snapshot.binaryHash,
      fixture_clock: '2026-01-01T00:00:00Z', groups, results, captured: snapshot.state,
      target_queries: [targetSQL], rows: snapshot.rows,
      gateway_queries: [bucketSQL], buckets: snapshot.buckets, ui_problems: problems,
      report_errors: reportErrors, fixture_cleanup: fixtureCleanup,
      retained_evidence: { workdir: work, paths: [gatewayDB, target, join(work, '*.log')], reason: 'diagnostic-evidence-retained' },
      passed: !process.exitCode && reportErrors.length === 0 && Object.keys(results).length === 8,
    }),
    writeReport: (report) => writeFileSync(join(evidence, 'quality-sqlite.json'), `${JSON.stringify(report)}\n`),
  });
}
