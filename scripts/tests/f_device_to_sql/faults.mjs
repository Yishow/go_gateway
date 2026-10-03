// Fault matrix on the production binary (SQLite destination), driven through the real UI setup:
//   OutageCrashAndUnknownCommit (destination offline + kill -9 + restart), silent device,
//   poison partition, and the confirmed test write from the group editor.
//   node scripts/tests/f_device_to_sql/faults.mjs
import { join } from 'node:path';
import { renameSync, writeFileSync } from 'node:fs';
import {
  ROOT, api, buildBinaries, cleanWork, configureDestination, configureDevicesPointsAndTags, createAndApplyGroup,
  activateDevices, GW_PORT, openBrowser, sh, startProcess, waitForGateway, waitForPortFree, witness, writeRegisters, PORT_A, PORT_B,
} from './lib.mjs';

const work = '/tmp/gw-f-faults';
const evidence = join(ROOT, 'docs/plans/studio-v2-write-groups/evidence-f');
const destination = join(work, 'destination.db');
const away = `${destination}.away`;
const lines = [];
const log = (m) => { lines.push(m); console.log(m); };
const results = {};
const check = (name, ok, detail) => { results[name] = { ok, detail }; log(`${ok ? 'PASS' : 'FAIL'} ${name}: ${detail}`); };
const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

await waitForPortFree();
cleanWork(work);
buildBinaries(work);
sh('sqlite3', [destination, 'CREATE TABLE readings (line TEXT, temperature INTEGER, pressure INTEGER, running INTEGER, batch INTEGER, prov TEXT);']);
sh('sqlite3', [destination, "INSERT INTO readings VALUES ('neighbor', 1, 2, 0, 7, '[]');"]);
const regsA = writeRegisters(work, 'regs-a', { 0: 215, 1: 1013, 2: 32, 3: 0, 4: 0, 5: 1, 6: 1 });
const regsB = writeRegisters(work, 'regs-b', { 0: 187, 1: 777, 2: 0, 3: 0, 4: 0, 5: 2, 6: 0 });
const children = {};
const kill = (name, signal = 'SIGTERM') => { try { process.kill(-children[name].pid, signal); } catch { /* gone */ } };
const startSim = (name, port, regs) => { children[name] = startProcess(work, name, join(work, 'f-sim'), ['-port', String(port), '-registers', regs]); };
const startGateway = (suffix = '') => { children.gateway = startProcess(work, `gateway${suffix}`, join(work, 'gw-ui'), [], { GATEWAY_DB_PATH: join(work, 'gateway.db'), PORT: GW_PORT }); };
process.on('exit', () => { for (const name of Object.keys(children)) kill(name); });

const q = (sql, path = destination) => sh('sqlite3', ['-separator', '|', path, sql]);
const bucketsOf = (line, path = destination) => q(`SELECT prov FROM readings WHERE line='${line}'`, path).split('\n').filter(Boolean)
  .map((p) => Math.floor(Date.parse(JSON.parse(p)[0].observed_at) / 10000) * 10000).sort((a, b) => a - b);
const delivery = async () => {
  const out = {};
  for (const group of (await api('/studio-v2/workspace/write-groups')).groups) out[group.name] = await api(`/studio-v2/workspace/write-groups/${group.id}/delivery`);
  return out;
};
const waitFor = async (what, predicate, ms = 90000) => {
  const deadline = Date.now() + ms;
  while (Date.now() < deadline) { if (await predicate()) return true; await sleep(1500); }
  log(`timed out waiting for ${what}`); return false;
};

startSim('sim-a', PORT_A, regsA); startSim('sim-b', PORT_B, regsB); startGateway();
await waitForGateway();
const { browser, page } = await openBrowser();
try {
  await configureDevicesPointsAndTags(page);
  await configureDestination(page, { kind: 'sqlite', path: destination, table: 'readings' });
  await createAndApplyGroup(page, 'Line A', 'a', log);
  await createAndApplyGroup(page, 'Line B', 'b', log);
  // The confirmed test write is a pre-activation check, so run it before the collector starts.
  // --- Confirmed test write from the group editor (ConfirmedTestWriteCleanup) ---
  await page.getByTestId('group-list').locator('button').first().click();
  await page.getByTestId('group-editor').waitFor();
  await page.getByTestId('group-test-write-preview').click();
  await page.getByTestId('group-test-write-preview-card').waitFor({ timeout: 15000 });
  const neighborBefore = q("SELECT temperature, pressure, running, batch FROM readings WHERE line='neighbor'");
  await page.getByTestId('group-test-write-confirm').click();
  await page.getByTestId('group-test-write-outcome').waitFor({ timeout: 30000 });
  const outcome = await page.getByTestId('group-test-write-outcome').innerText();
  const cleanup = await page.getByTestId('group-test-write-cleanup').innerText();
  const leftover = q("SELECT COUNT(*) FROM readings WHERE line LIKE 'gw-test-%'");
  const neighborAfter = q("SELECT temperature, pressure, running, batch FROM readings WHERE line='neighbor'");
  await page.screenshot({ path: join(evidence, 'faults-test-write.png'), fullPage: true });
  check('ConfirmedTestWriteCleanup', /read back correctly|讀回/.test(outcome) && /removed|已移除|已清/.test(cleanup) && leftover === '0' && neighborBefore === neighborAfter,
    `outcome="${outcome}" cleanup="${cleanup}" leftover_test_rows=${leftover} neighbor_unchanged=${neighborBefore === neighborAfter}`);
  await page.getByTestId('group-editor-close').click();

  await activateDevices(page);
  await waitFor('two rows per line', () => bucketsOf('A').length >= 2 && bucketsOf('B').length >= 2);

  // --- Silent device: stop device B for ~35 s ---
  const beforeSilence = bucketsOf('B').length;
  kill('sim-b'); const silenceStart = Date.now();
  await sleep(35000);
  startSim('sim-b', PORT_B, regsB);
  await waitFor('B rows to resume', () => bucketsOf('B').length > beforeSilence + 1, 60000);
  const b = bucketsOf('B');
  const gaps = b.slice(1).map((t, i) => (t - b[i]) / 10000).filter((g) => g > 1);
  const view = await delivery();
  check('SilentDeviceNoInventedRows', gaps.length >= 1 && new Set(b).size === b.length && view['Line B'].skipped_buckets + view['Line B'].no_data_buckets > 0,
    `bucket gaps (in 10 s buckets)=${JSON.stringify(gaps)} skipped=${view['Line B'].skipped_buckets} no_data=${view['Line B'].no_data_buckets} duplicates=${b.length - new Set(b).size}`);

  // --- Destination offline, then kill -9 of the gateway during the outage, then restart and recover ---
  const aBefore = bucketsOf('A').length;
  const rowsBeforeOutage = Number(q('SELECT COUNT(*) FROM readings'));
  renameSync(destination, away);
  const outageStart = Date.now();
  await sleep(35000);
  const during = await delivery();
  const pending = (g) => g.stages.queued + g.stages.retrying + g.stages.blocked + g.stages.unknown;
  const rowsDuringOutage = Number(q('SELECT COUNT(*) FROM readings', away));
  check('OutageBacklogIsDurableAndNotWritten', pending(during['Line A']) >= 2 && rowsDuringOutage === rowsBeforeOutage,
    `A pending=${pending(during['Line A'])} B pending=${pending(during['Line B'])} destination rows before/during outage=${rowsBeforeOutage}/${rowsDuringOutage}`);
  const rowsAtKill = rowsDuringOutage;
  kill('gateway', 'SIGKILL');
  const killedAt = Date.now();
  await sleep(2000);
  startGateway('-restarted');
  await waitForGateway();
  renameSync(away, destination);
  const recovered = await waitFor('backlog drained after restart', async () => {
    try { const v = await delivery(); return pending(v['Line A']) === 0 && pending(v['Line B']) === 0 && Number(q('SELECT COUNT(*) FROM readings')) > rowsAtKill + 4; } catch { return false; }
  }, 150000);
  const aBuckets = bucketsOf('A');
  const outageBuckets = aBuckets.filter((t) => t >= outageStart - 10000 && t <= killedAt - 10000);
  const dupA = aBuckets.length - new Set(aBuckets).size;
  const dupB = bucketsOf('B').length - new Set(bucketsOf('B')).size;
  const expected = Math.floor((killedAt - 10000) / 10000) - Math.floor(outageStart / 10000);
  check('OutageCrashAndUnknownCommit', recovered && dupA === 0 && dupB === 0 && outageBuckets.length >= Math.max(1, expected - 1),
    `recovered=${recovered} duplicate_buckets A=${dupA} B=${dupB} outage_buckets_delivered=${outageBuckets.length}/expected~${expected} rows ${aBefore}->${aBuckets.length} (kill -9 during the outage)`);
  // A crash-restart must keep collecting, not only drain the backlog (this caught a mapping drift that made
  // every sample after a restart refuse with revision-mismatch).
  const afterRestart = Date.now();
  const resumed = await waitFor('new buckets after the restart', () => bucketsOf('A').some((t) => t > afterRestart) && bucketsOf('B').some((t) => t > afterRestart), 60000);
  check('RestartResumesCollection', resumed, `rows for buckets after the restart appeared=${resumed}`);
  results.NotCoveredByThisRun = { ok: null, detail: 'endpoint retarget of a backlog (service-level RevisionBoundBacklog test), Linux/Windows/ARM, real PLC' };

  // --- Poison partition: the destination rejects every row of line B ---
  q("CREATE TRIGGER poison BEFORE INSERT ON readings WHEN NEW.line='B' BEGIN SELECT RAISE(ABORT, 'poison'); END;");
  const aAtPoison = bucketsOf('A').length;
  const quarantined = await waitFor('line B quarantined while A continues', async () => {
    const v = await delivery(); return v['Line B'].stages.quarantined >= 1 && bucketsOf('A').length >= aAtPoison + 2;
  }, 120000);
  const vp = await delivery();
  check('PoisonPartitionIsolated', quarantined, `B quarantined=${vp['Line B'].stages.quarantined} blocked=${vp['Line B'].stages.blocked} while A kept advancing ${aAtPoison}->${bucketsOf('A').length}`);
  await page.reload(); await page.waitForTimeout(3000);
  await page.getByTestId('step-nav-button-4').click(); await page.waitForTimeout(3000);
  await page.screenshot({ path: join(evidence, 'faults-poison.png'), fullPage: true });
} finally {
  writeFileSync(join(evidence, 'faults-sqlite.json'), JSON.stringify(witness({ results, log: lines }), null, 2));
  await browser.close();
}
process.exitCode = Object.values(results).every((r) => r.ok !== false) ? 0 : 1;
