// DeviceToSQLiteMixedRows / DeviceToPostgresMixedRows: real UI -> gateway -> SQL.
//   node scripts/tests/f_device_to_sql/run.mjs sqlite
//   POSTGRES_DSN='host=... port=... user=... password=... dbname=...' node scripts/tests/f_device_to_sql/run.mjs postgres
import { join } from 'node:path';
import { writeFileSync } from 'node:fs';
import {
  ROOT, api, buildBinaries, cleanupOwnedPostgresSchema, cleanWork, configureDestination, configureDevicesPointsAndTags, createAndApplyGroup,
  activateDevices, GW_PORT, makeRunNamespace, openBrowser, sh, startProcess, stopProcess, stopProcesses, waitForGateway, waitForPortFree, witness, writeRegisters, PORT_A, PORT_B,
  validateOwnedPostgresTarget,
} from './lib.mjs';
import { collectAndValidatePointObservations } from './mixed-identity.mjs';

const kind = process.argv[2] ?? 'sqlite';
if (!['sqlite', 'postgres'].includes(kind)) throw new Error(`invalid destination kind: ${kind}; expected sqlite or postgres`);
const { run_id, work } = makeRunNamespace(kind, process.env.F_RUN_ID);
const evidence = join(ROOT, 'docs/plans/studio-v2-write-groups/evidence-f');
const regsA = [215, 1013, 32, 0, 0, 1, 1, 0x3fc0, 0x0000, 0x405e, 0xdd2f, 0x1a9f, 0xbe77, 0x4131];
const regsB = [187, 777, 0, 0, 0, 2, 0, 0xc010, 0x0000, 0xc084, 0x7291, 0x6872, 0xb021, 0x4232];
const lines = [];
const log = (message) => { lines.push(message); console.log(message); };
const evidencePath = join(evidence, `device-to-${kind}-mixed.json`);
const commands = [
  'make sync-frontend-static',
  'go build ./cmd/test_ui',
  'go build ./cmd/f_modbus_simulator',
  `node scripts/tests/f_device_to_sql/run.mjs ${kind}`,
];
const typed_scope = {
  acquisition_types: ['int16', 'uint16', 'bool', 'uint64', 'float32', 'float64', 'string'],
  decimal: {
    acquisition: 'unsupported',
    acquisition_status: 'not-run',
    codec_sql_round_trip: 'separate-tests',
    separate_test_paths: [
      'internal/datalink/dbtarget/exact_value_codec_test.go',
      'internal/datalink/dbtarget/exact_value_codec_postgres_test.go',
    ],
  },
};

function psql(sql) {
  return sh('docker', ['exec', 'gw-wg-pg-test', 'psql', '-v', 'ON_ERROR_STOP=1', '-U', 'postgres', '-d', 'gwtest', '-At', '-F', '|', '-c', sql]);
}
function parseDsn(text) {
  return Object.fromEntries(text.trim().split(/\s+/).filter(Boolean).map((part) => part.split('=')));
}

const processes = [];
let build;
let destination;
let query;
let browser;
let page;
let problems = [];
let mappingPreviewRequests = [];
let result;
let failure;
let cleanup = { owned_children: [] };
let postgresResource;
const signalChildren = () => { for (const child of processes) { try { stopProcess(child); } catch { /* best effort during process shutdown */ } } };
process.on('exit', signalChildren);
try {
  await waitForPortFree();
  cleanWork(work);
  build = buildBinaries(work);
  if (kind === 'sqlite') {
    const path = join(work, 'destination.db');
    sh('sqlite3', [path, 'CREATE TABLE readings (line TEXT, temperature INTEGER, pressure INTEGER, running INTEGER, batch INTEGER, reading_f32 REAL, reading_f64 REAL, text_value TEXT, prov TEXT);']);
    sh('sqlite3', [path, "INSERT INTO readings VALUES ('neighbor', 1, 2, 0, 7, NULL, NULL, NULL, '[]');"]);
    destination = { kind: 'sqlite', path, table: 'readings' };
    query = (sql) => sh('sqlite3', ['-cmd', '.timeout 15000', '-separator', '|', path, sql]);
  } else {
    const dsn = parseDsn(process.env.POSTGRES_DSN ?? '');
    if (!dsn.host) throw new Error('POSTGRES_DSN is required for the postgres run (blocked, not passed, without it)');
    const image = sh('docker', ['inspect', '--format', '{{.Config.Image}}', 'gw-wg-pg-test']);
    validateOwnedPostgresTarget(dsn, image);
    const schema = `gw_f_${Date.now()}`;
    psql(`CREATE SCHEMA "${schema}"`);
    postgresResource = { schema };
    psql(`CREATE TABLE "${schema}".readings (line TEXT, temperature BIGINT, pressure BIGINT, running BOOLEAN, batch BIGINT, reading_f32 DOUBLE PRECISION, reading_f64 DOUBLE PRECISION, text_value TEXT, prov JSONB); INSERT INTO "${schema}".readings VALUES ('neighbor', 1, 2, false, 7, NULL, NULL, NULL, '[]');`);
    destination = { kind: 'postgres', host: dsn.host, port: Number(dsn.port ?? 5432), database: dsn.dbname, user: dsn.user, password: dsn.password, schema, table: 'readings' };
    query = (sql) => psql(sql.replaceAll('readings', `"${schema}".readings`));
  }
  processes.push(startProcess(work, 'sim-a', join(work, 'f-sim'), ['-port', String(PORT_A), '-registers', writeRegisters(work, 'regs-a', Object.fromEntries(regsA.map((v, i) => [i, v])))]));
  processes.push(startProcess(work, 'sim-b', join(work, 'f-sim'), ['-port', String(PORT_B), '-registers', writeRegisters(work, 'regs-b', Object.fromEntries(regsB.map((v, i) => [i, v])))]));
  processes.push(startProcess(work, 'gateway', join(work, 'gw-ui'), [], { GATEWAY_DB_PATH: join(work, 'gateway.db'), PORT: GW_PORT }));
  await waitForGateway();
  ({ browser, page, problems, mappingPreviewRequests } = await openBrowser());
  await configureDevicesPointsAndTags(page, { extended: true });
  await configureDestination(page, destination);
  await createAndApplyGroup(page, 'Line A', 'a', log, { extended: true });
  await createAndApplyGroup(page, 'Line B', 'b', log, { extended: true });
  await activateDevices(page);
  log('activation confirmed through the UI');

  // The first rows land after the next 10 s bucket closes; wait for three of each line.
  const deadline = Date.now() + 90_000;
  let rows = [];
  while (Date.now() < deadline) {
    // The extended query is deliberately independent of the UI/runtime API.
    rows = query(`SELECT ${kind === 'sqlite' ? 'rowid' : 'ctid::text'} AS sql_row_id, line, temperature, pressure, running, batch, reading_f32, reading_f64, text_value, prov FROM readings WHERE line IN ('A','B') ORDER BY line, ${kind === 'sqlite' ? 'rowid' : 'prov::text'};`).split('\n').filter(Boolean);
    const rowLine = (row) => row.split('|', 3)[1];
    if (rows.filter((row) => rowLine(row) === 'A').length >= 3 && rows.filter((row) => rowLine(row) === 'B').length >= 3) break;
    await new Promise((r) => setTimeout(r, 2000));
  }
  await page.waitForTimeout(1500);
  await page.screenshot({ path: join(evidence, `device-to-${kind}-step4.png`), fullPage: true });

  const parsed = rows.map((row) => {
    const [sql_row_id, line, temperature, pressure, running, batch, reading_f32, reading_f64, text_value, ...prov] = row.split('|');
    return { sql_row_id, line, temperature, pressure, running, batch, reading_f32, reading_f64, text_value, prov: prov.join('|') };
  });
  const truthy = (value) => value === '1' || value === 't' || value === 'true';
  const failures = [];
  const expectRow = (row, want) => {
    const got = { temperature: row.temperature, pressure: row.pressure, running: truthy(row.running), batch: row.batch, reading_f32: row.reading_f32, reading_f64: row.reading_f64, text_value: row.text_value };
    if (JSON.stringify(got) !== JSON.stringify(want)) failures.push(`${row.line}: expected ${JSON.stringify(want)} got ${JSON.stringify(got)}`);
  };
  for (const row of parsed) expectRow(row, row.line === 'A'
    ? { temperature: '215', pressure: '1013', running: true, batch: '9007199254740993', reading_f32: '1.5', reading_f64: '123.456', text_value: 'A1' }
    : { temperature: '187', pressure: '777', running: false, batch: '2', reading_f32: '-2.25', reading_f64: '-654.321', text_value: 'B2' });
  const groups = (await api('/studio-v2/workspace/write-groups')).groups;
  const persistedMembers = groups.map((group) => ({
    id: group.id,
    name: group.name,
    workspace_id: group.workspace_id,
    revision: group.revision,
    applied_revision: group.applied_revision,
    members: group.members.map((member) => ({
      device_id: member.device_id,
      point_id: member.point_id,
      tag_id: member.tag_id,
      source_revision: member.source_revision,
      mapping_revision: member.mapping_revision,
      target_column: member.target_column,
    })),
  }));
  if (groups.length !== 2 || groups.some((group) => group.members.length !== 7)) failures.push('persisted groups do not contain two seven-member groups');
  const pointResult = await collectAndValidatePointObservations(
    groups,
    (deviceID) => api(`/points?device_id=${encodeURIComponent(deviceID)}`),
  );
  const pointObservations = pointResult.observations;
  failures.push(...pointResult.failures);
  for (const group of groups) {
    const keys = new Set(group.members.map((member) => JSON.stringify([member.device_id, member.point_id, member.tag_id])));
    for (const row of parsed.filter((candidate) => candidate.line === (group.name.endsWith('A') ? 'A' : 'B'))) {
      const entries = JSON.parse(row.prov);
      const actualKeys = new Set(entries.map((entry) => entry.member));
      if (actualKeys.size !== keys.size || [...keys].some((key) => !actualKeys.has(key))) failures.push(`${group.name}: SQL provenance member IDs do not match persisted group IDs`);
    }
  }
  for (const letter of ['A', 'B']) {
    const own = parsed.filter((r) => r.line === letter);
    if (own.length < 3) failures.push(`line ${letter}: only ${own.length} rows arrived`);
    const times = own.map((r) => JSON.parse(r.prov)[0]?.observed_at);
    if (new Set(times).size !== times.length) failures.push(`line ${letter}: duplicate bucket rows`);
    for (const r of own) {
      const entries = JSON.parse(r.prov);
      if (entries.length !== 7 || entries.some((e) => e.status !== 'ok' || e.quality !== 'good' || !/Z$/.test(e.observed_at))) failures.push(`line ${letter}: provenance is not seven good UTC samples`);
    }
  }
  const targetSQL = `SELECT ${kind === 'sqlite' ? 'rowid' : 'ctid::text'} AS sql_row_id, line, temperature, pressure, running, batch, reading_f32, reading_f64, text_value, prov FROM readings WHERE line IN ('A','B') ORDER BY line, ${kind === 'sqlite' ? 'rowid' : 'prov::text'};`;
  const neighborSQL = "SELECT line, temperature, pressure, batch, reading_f32, reading_f64, text_value FROM readings WHERE line = 'neighbor';";
  const neighbor = query(neighborSQL);
  if (neighbor !== 'neighbor|1|2|7|||') failures.push(`neighbor row changed: ${neighbor}`);
  const sampleSources = problems.filter((p) => !/source-rules|points\/poll/.test(p));
  if (sampleSources.length) failures.push(...sampleSources.map((p) => `browser/API problem: ${p}`));

  const delivery = [];
  for (const group of groups) {
    const view = await api(`/studio-v2/workspace/write-groups/${group.id}/delivery`);
    delivery.push({ group: group.name, applied_revision: group.applied_revision, intake: view.intake, stages: view.stages });
  }
  const sqlRows = parsed.map((row) => ({ ...row, provenance: JSON.parse(row.prov), prov: undefined }));
  result = { run_id, kind, fixture: 'extended-float32-float64-text', typed_scope, commands, build, ui_screenshot: `docs/plans/studio-v2-write-groups/evidence-f/device-to-${kind}-step4.png`, passed: failures.length === 0, failures, expected: { A: { temperature: '215', pressure: '1013', running: true, batch: '9007199254740993', reading_f32: '1.5', reading_f64: '123.456', text_value: 'A1' }, B: { temperature: '187', pressure: '777', running: false, batch: '2', reading_f32: '-2.25', reading_f64: '-654.321', text_value: 'B2' } }, persisted_members: persistedMembers, point_observations: pointObservations, sql: { target_query: targetSQL, neighbor_query: neighborSQL, rows: sqlRows, neighbor }, delivery, log: lines, ui_problems: problems, mapping_preview_requests: mappingPreviewRequests };
  log(result.passed ? `PASS: ${parsed.length} rows verified in ${kind}` : `FAIL: ${failures.join('; ')}`);
  process.exitCode = result.passed ? 0 : 1;
} catch (error) {
  failure = error;
  process.exitCode = 1;
  if (page) await page.screenshot({ path: join(work, 'failure.png'), fullPage: true }).catch(() => {});
  log(`ERROR: ${error?.message?.split('\n')[0] ?? String(error)}`);
} finally {
  cleanup = { owned_children: await stopProcesses(processes), retained_work: work };
  if (postgresResource) {
    cleanup.postgres = cleanupOwnedPostgresSchema(postgresResource, psql);
    if (!cleanup.postgres.passed) {
      const message = `owned PostgreSQL schema cleanup failed: ${cleanup.postgres.error ?? 'schema remains or verification failed'}`;
      if (result) {
        result.failures.push(message);
        result.passed = false;
      } else if (!failure) {
        failure = new Error(message);
      }
      process.exitCode = 1;
    }
  }
  const cleanupFailures = cleanup.owned_children.filter((status) => status.state !== 'exited'
    || (status.exit_code !== null && status.exit_code !== 0)
    || status.stop_error || status.spawn_error || status.kill_error);
  if (cleanupFailures.length) {
    const message = `owned child cleanup incomplete: ${JSON.stringify(cleanupFailures)}`;
    if (result) {
      result.failures.push(message);
      result.passed = false;
    } else if (!failure) {
      failure = new Error(message);
    }
    process.exitCode = 1;
  }
  try { await browser?.close(); } catch (error) {
    const message = `browser cleanup failed: ${error.message.split('\n')[0]}`;
    if (result) { result.failures.push(message); result.passed = false; }
    else if (!failure) failure = error;
    process.exitCode = 1;
  }
  process.removeListener('exit', signalChildren);
  const retained_artifacts = {
    work_directory: work,
    database_paths: [join(work, 'gateway.db'), ...(kind === 'sqlite' ? [join(work, 'destination.db')] : [])],
    log_paths: [join(work, 'gateway.log'), join(work, 'sim-a.log'), join(work, 'sim-b.log')],
    policy: 'owned database and process logs are retained for diagnostics; cleanup only stops children',
  };
  const errorMessage = failure?.message?.split('\n')[0];
  const payload = result ?? {
    run_id,
    kind,
    fixture: 'extended-float32-float64-text',
    typed_scope,
    commands,
    build: build ?? null,
    ui_screenshot: `docs/plans/studio-v2-write-groups/evidence-f/device-to-${kind}-step4.png`,
    passed: false,
    failures: [errorMessage ?? 'mixed acceptance setup failed'],
    log: lines,
    ui_problems: problems,
    mapping_preview_requests: mappingPreviewRequests,
  };
  payload.cleanup = cleanup;
  payload.retained_artifacts = retained_artifacts;
  writeFileSync(evidencePath, `${JSON.stringify(witness(payload))}\n`);
}

if (failure) throw failure;
