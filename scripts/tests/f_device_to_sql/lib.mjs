// Shared pieces of the device-to-SQL acceptance harness. Everything here talks to
// the real gateway binary (with its embedded UI) through a real browser; the only
// stand-ins are the loopback Modbus simulators (devices) and the disposable
// destination database. Nothing inserts rows into the destination except the gateway.
import { chromium } from '../../../frontend/node_modules/playwright/index.mjs';
import { spawn, execFileSync } from 'node:child_process';
import { mkdirSync, writeFileSync, existsSync, rmSync, readFileSync, openSync, closeSync } from 'node:fs';
import { createHash } from 'node:crypto';
import { join, resolve } from 'node:path';
import os from 'node:os';

export const ROOT = resolve(new URL('../../..', import.meta.url).pathname);
// Another local service may already hold 3333 (it answers 404 for the UI), so the
// harness runs its gateway on its own port.
export const GW_PORT = process.env.GW_PORT ?? '3343';
export const BASE = `http://127.0.0.1:${GW_PORT}`;
function simulatorPort(name, fallback) {
  const raw = process.env[name] ?? String(fallback);
  const value = Number(raw);
  if (!/^\d+$/.test(raw) || !Number.isInteger(value) || value < 1024 || value > 65535) {
    throw new Error(`${name} must name an unprivileged TCP port`);
  }
  return value;
}
export const PORT_A = simulatorPort('F_SIM_PORT_A', 15020);
export const PORT_B = simulatorPort('F_SIM_PORT_B', 15021);

export function sh(cmd, args, options = {}) {
  return execFileSync(cmd, args, { encoding: 'utf8', ...options }).trim();
}

export function buildBinaries(work, { tags } = {}) {
  mkdirSync(work, { recursive: true });
  sh('make', ['sync-frontend-static'], { cwd: ROOT, stdio: 'pipe' });
  sh('go', ['build', ...(tags ? ['-tags', tags] : []), '-o', join(work, 'gw-ui'), './cmd/test_ui'], { cwd: ROOT });
  sh('go', ['build', '-o', join(work, 'f-sim'), './cmd/f_modbus_simulator'], { cwd: ROOT });
  const sha256 = (path) => createHash('sha256').update(readFileSync(path)).digest('hex');
  return { gateway_sha256: sha256(join(work, 'gw-ui')), simulator_sha256: sha256(join(work, 'f-sim')) };
}

export function writeRegisters(work, name, registers) {
  const file = join(work, `${name}.json`);
  writeFileSync(file, JSON.stringify(registers));
  return file;
}

export function startProcess(work, name, command, args, env = {}) {
  const log = join(work, `${name}.log`);
  const logFD = openSync(log, 'a');
  let child;
  try {
    child = spawn(command, args, {
      cwd: work, env: { ...process.env, ...env }, detached: true, stdio: ['ignore', logFD, logFD],
    });
  } finally {
    closeSync(logFD);
  }
  child.fHarnessName = name;
  child.fSpawnError = null;
  child.on('error', (error) => { child.fSpawnError = error.message.split('\n')[0]; });
  child.unref();
  return child;
}

const RUN_ID_PATTERN = /^[A-Za-z0-9][A-Za-z0-9_-]{0,80}$/;

export function makeRunNamespace(kind, requestedRunID) {
  if (!['sqlite', 'postgres'].includes(kind)) throw new Error(`invalid destination kind: ${kind}`);
  const run_id = requestedRunID ?? `${kind}-${Date.now()}`;
  if (!RUN_ID_PATTERN.test(run_id)) {
    throw new Error('F_RUN_ID must be a short alphanumeric run namespace using only - or _');
  }
  const work = `/tmp/gw-f-${run_id}`;
  if (resolve(work) !== work || !/^\/tmp\/gw-f-[A-Za-z0-9_-]+$/.test(work)) {
    throw new Error('F_RUN_ID resolved outside the owned /tmp/gw-f-* namespace');
  }
  return { run_id, work };
}

function childExited(child) {
  return child.exitCode !== null || child.signalCode !== null;
}

export function stopProcess(child) {
  if (!child?.pid) return;
  try {
    process.kill(-child.pid, 'SIGTERM');
  } catch (error) {
    if (error?.code !== 'ESRCH') throw error;
  }
}

function killProcess(child) {
  if (!child?.pid) return;
  try {
    process.kill(-child.pid, 'SIGKILL');
  } catch (error) {
    if (error?.code !== 'ESRCH') throw error;
  }
}

function waitForProcessExit(child, timeoutMs) {
  if (!child?.pid || childExited(child)) {
    return Promise.resolve({
      pid: child?.pid ?? null,
      state: 'exited',
      exit_code: child?.exitCode ?? null,
      signal: child?.signalCode ?? null,
    });
  }
  return new Promise((resolve) => {
    let timer;
    const finish = (state) => {
      clearTimeout(timer);
      child.removeListener('close', onClose);
      child.removeListener('error', onError);
      resolve({
        pid: child.pid,
        state,
        exit_code: child.exitCode,
        signal: child.signalCode,
      });
    };
    const onClose = () => finish('exited');
    const onError = (error) => finish(`error:${error.message.split('\n')[0]}`);
    child.once('close', onClose);
    child.once('error', onError);
    timer = setTimeout(() => finish('timeout'), timeoutMs);
  });
}

export async function stopProcesses(children, timeoutMs = 5000) {
  for (const child of children) {
    try { stopProcess(child); } catch (error) {
      child._fStopError = error.message.split('\n')[0];
    }
  }
  const termStatuses = await Promise.all(children.map((child) => waitForProcessExit(child, timeoutMs)));
  const statuses = await Promise.all(termStatuses.map(async (status, index) => {
    const child = children[index];
    if (status.state !== 'timeout') return { term_state: status.state, kill_sent: false, ...status };
    let killError;
    try { killProcess(child); } catch (error) { killError = error; }
    const afterKill = await waitForProcessExit(child, timeoutMs);
    return {
      term_state: 'timeout',
      kill_sent: !killError,
      ...(killError ? { kill_error: killError.message.split('\n')[0] } : {}),
      ...afterKill,
    };
  }));
  return statuses.map((status, index) => ({
    name: children[index]?.fHarnessName,
    ...status,
    ...(children[index]?.fSpawnError ? { spawn_error: children[index].fSpawnError } : {}),
    ...(children[index]?._fStopError ? { stop_error: children[index]._fStopError } : {}),
  }));
}

/** Fetch with one deadline covering connection, headers, and complete body consumption. */
export async function boundedFetch(url, options = {}, timeoutMs = 30000) {
  if (!Number.isFinite(timeoutMs) || timeoutMs <= 0 || timeoutMs > 2 ** 31 - 1) {
    throw new TypeError('timeoutMs must be a positive finite duration');
  }
  const timeoutSignal = AbortSignal.timeout(Math.max(1, Math.floor(timeoutMs)));
  const signal = options.signal ? AbortSignal.any([options.signal, timeoutSignal]) : timeoutSignal;
  const response = await fetch(url, { ...options, signal });
  const body = await response.arrayBuffer();
  return new Response(body, {
    status: response.status,
    statusText: response.statusText,
    headers: response.headers,
  });
}

/** A previous run's gateway can still be shutting down; starting on top of it would test the wrong process. */
export async function waitForPortFree(timeoutMs = 20000) {
  const deadline = Date.now() + timeoutMs;
  while (Date.now() < deadline) {
    try {
      await boundedFetch(`${BASE}/studio/v2`, {}, Math.min(3000, deadline - Date.now()));
    } catch (error) {
      if (error?.name !== 'TimeoutError' && error?.name !== 'AbortError') return;
    }
    await new Promise((r) => setTimeout(r, 300));
  }
  throw new Error(`port ${GW_PORT} is still in use`);
}

export async function waitForGateway(timeoutMs = 30000) {
  const deadline = Date.now() + timeoutMs;
  while (Date.now() < deadline) {
    try {
      const requestTimeout = Math.min(3000, deadline - Date.now());
      const api = await boundedFetch(`${BASE}/api/v1/datalink/studio-v2/workspace`, {}, requestTimeout);
      const ui = await boundedFetch(`${BASE}/studio/v2`, {}, requestTimeout);
      if (api.status < 500 && ui.status === 200) return;
    } catch { /* not up yet */ }
    await new Promise((r) => setTimeout(r, 300));
  }
  throw new Error('gateway did not start');
}

export async function api(path, options = {}, timeoutMs = 30000) {
  const response = await boundedFetch(`${BASE}/api/v1/datalink${path}`, options, timeoutMs);
  let body;
  try {
    body = await response.json();
  } catch (error) {
    throw new Error(`API ${path} returned invalid JSON: ${error.message.split('\n')[0]}`);
  }
  if (!response.ok) {
    const code = body?.error?.code ?? `http_${response.status}`;
    throw new Error(`API ${path} failed (${response.status}, ${code})`);
  }
  if (!body || !Object.prototype.hasOwnProperty.call(body, 'data')) {
    throw new Error(`API ${path} response is missing data`);
  }
  return body.data;
}

export function validateOwnedPostgresTarget(dsn, image) {
  if (dsn?.host !== '127.0.0.1' || dsn?.port !== '55432' || dsn?.dbname !== 'gwtest' || dsn?.user !== 'postgres' || !dsn?.password) {
    throw new Error('POSTGRES_DSN must target the owned loopback PostgreSQL fixture at 127.0.0.1:55432/gwtest as postgres');
  }
  if (!/^postgres:16(?:-alpine)?$/.test(image ?? '')) {
    throw new Error(`owned PostgreSQL container must use postgres:16, got ${image ?? 'missing image'}`);
  }
}

export function cleanupOwnedPostgresSchema(resource, psql) {
  const schema = resource?.schema;
  const outcome = {
    resource: 'owned-postgres-schema',
    schema,
    passed: false,
    drop_sql: schema ? `DROP SCHEMA "${schema}" CASCADE` : null,
    verify_sql: schema ? `SELECT nspname FROM pg_namespace WHERE nspname = '${schema}'` : null,
  };
  if (!/^gw_f_[A-Za-z0-9_-]+$/.test(schema ?? '')) {
    outcome.error = 'invalid owned schema name';
    return outcome;
  }
  let dropError;
  try {
    psql(outcome.drop_sql);
    outcome.drop = { status: 'ok' };
  } catch (error) {
    dropError = error;
    outcome.drop = { status: 'error', error: error.message.split('\n')[0] };
  }
  let verifyResult;
  try {
    verifyResult = String(psql(outcome.verify_sql) ?? '').trim();
    outcome.verify = { status: 'ok', result: verifyResult };
  } catch (error) {
    outcome.verify = { status: 'error', error: error.message.split('\n')[0] };
  }
  if (dropError) outcome.error = 'DROP SCHEMA failed';
  else if (outcome.verify?.status !== 'ok') outcome.error = 'schema absence verification failed';
  else if (verifyResult !== '') outcome.error = `schema still exists: ${verifyResult}`;
  else outcome.passed = true;
  return outcome;
}

export async function openBrowser(width = 1440, height = 1400) {
  const browser = await chromium.launch();
  try {
    const page = await browser.newPage({ viewport: { width, height } });
    const problems = [];
    const mappingPreviewRequests = [];
    page.on('pageerror', (error) => problems.push(`pageerror: ${error.message}`));
    page.on('response', async (r) => {
      if (r.status() < 400 || !r.url().includes('/api/')) return;
      const path = new URL(r.url()).pathname;
      const responseText = await r.text().catch(() => '');
      if (path === '/api/v1/datalink/mappings/preview' && mappingPreviewRequests.length < 64) {
        let requestBody;
        try {
          requestBody = r.request().postDataJSON();
        } catch (error) {
          requestBody = { post_data_json_error: error instanceof Error ? error.message : String(error) };
        }
        let responseBody;
        try {
          responseBody = JSON.parse(responseText);
        } catch {
          responseBody = responseText.slice(0, 512);
        }
        mappingPreviewRequests.push({
          status: r.status(),
          method: r.request().method(),
          path,
          request: requestBody,
          response: responseBody,
        });
      }
      problems.push(`${r.status()} ${r.request().method()} ${r.url().replace(BASE, '')} ${responseText.slice(0, 160)}`);
    });
    await page.goto(`${BASE}/studio/v2`);
    await page.waitForTimeout(2000);
    return { browser, page, problems, mappingPreviewRequests };
  } catch (error) {
    await browser.close().catch(() => {});
    throw error;
  }
}

export async function addDevice(page, name, port) {
  await page.getByTestId('btn-add-device').click();
  await page.getByTestId('editor-input-name').fill(name);
  await page.getByTestId('input-host').fill('127.0.0.1');
  await page.getByTestId('input-port').fill(String(port));
  // The probe is persisted only for a saved device, so wait for the save first.
  await page.waitForFunction(() => document.querySelector('[data-testid=device-save-banner]')?.getAttribute('data-save-state') === 'saved', null, { timeout: 15000 });
  await page.getByTestId('run-test-button').click();
  await page.getByTestId('success-readiness-card').waitFor({ timeout: 15000 });
  await page.waitForTimeout(2500);
}

export async function addRule(page, name, start, count, dataType, prefix, existingTab) {
  if (existingTab) await page.getByTestId('rule-tab-rail').getByText(existingTab, { exact: false }).first().click({ force: true, position: { x: 8, y: 8 } });
  await page.getByTestId('rule-add-btn').click();
  await page.waitForTimeout(600);
  await page.getByTestId('rule-name-input').fill(name);
  await page.getByTestId('rule-start-input').fill(start);
  await page.getByTestId('rule-count-input').fill(String(count));
  await page.getByTestId('rule-datatype-select').selectOption(dataType);
  await page.getByTestId('rule-prefix-input').fill(prefix);
  await page.waitForTimeout(2500);
}

const BASE_MIXED_SPECS = [
  ['temp', '40001', 'int16', 'TEMP_'],
  ['pressure', '40002', 'uint16', 'PRESS_'],
  ['batch', '40003', 'uint64', 'BATCH_'],
  ['running', '40007', 'bool', 'RUN_'],
];

// The extended fixture keeps the original four members and adds the decoder
// widths that the production Modbus path already supports. `text` uses the
// Rule UI's native string source/target so readiness and runtime exercise the
// same persisted type contract as the other members.
const EXTENDED_MIXED_SPECS = [
  ['f32', '40008', 'float32', 'F32_'],
  ['f64', '40010', 'float64', 'F64_'],
  ['text', '40014', 'string', 'TEXT_'],
];

/** Two devices with the same addresses and four default Tag types, optionally extended with float32/float64/text targets. */
export async function configureDevicesPointsAndTags(page, { extended = false } = {}) {
  await addDevice(page, 'Line A', PORT_A);
  await addDevice(page, 'Line B', PORT_B);
  await page.getByTestId('btn-continue-step1').click();
  await page.waitForTimeout(2000);
  const specs = [...BASE_MIXED_SPECS, ...(extended ? EXTENDED_MIXED_SPECS : [])];
  const last = {};
  for (const [name, start, dataType, prefix] of specs) {
    for (const device of ['Line A', 'Line B']) {
      await addRule(page, `${device} ${name}`, start, 1, dataType, `${device[5]}_${prefix}`, last[device]);
      last[device] = `${device} ${name}`;
    }
  }
  await page.getByTestId('continue-step3-btn').click();
  await page.waitForTimeout(2500);
  const tagTypes = { temp: 'int16', press: 'uint16', batch: 'uint64', run: 'bool', f32: 'float32', f64: 'float64', text: 'string' };
  for (const select of await page.locator('select[data-testid^=select-target-type]').all()) {
    const id = (await select.getAttribute('data-testid')).replace('select-target-type-', '');
    const tagKey = await page.getByTestId(`input-tag-key-${id}`).inputValue();
    const kind = Object.keys(tagTypes).find((k) => tagKey.includes(`.${k}.`));
    if (kind) await select.selectOption(tagTypes[kind]);
  }
  await page.waitForTimeout(6000);
  await page.getByTestId('btn-continue').click();
  await page.waitForTimeout(3000);
}

async function waitSaved(field, expected, page, refill) {
  for (let i = 0; i < 20; i++) {
    const cfg = await api('/studio-v2/workspace/database-config');
    if (Object.entries(expected).every(([k, v]) => cfg?.[k] === v)) return;
    await page.waitForTimeout(1000);
    if (i % 4 === 3) await refill();
  }
  throw new Error(`destination was not saved: ${field}`);
}

/** Fills the destination in Step 4, waits until the server holds it, then reloads so the UI holds its revision. */
export async function configureDestination(page, destination) {
  if (destination.kind === 'sqlite') {
    await page.getByRole('button', { name: /SQLite/ }).click();
    const fill = async () => {
      await page.getByLabel('Connection Name').fill('Line DB');
      await page.getByLabel('DB File Path').fill(destination.path);
      await page.getByLabel('Table Name').fill(destination.table);
    };
    await page.waitForTimeout(800);
    await fill();
    await waitSaved('sqlite', { database: destination.path, table: destination.table }, page, fill);
  } else {
    await page.getByRole('button', { name: /PostgreSQL/ }).click();
    const fill = async () => {
      await page.getByLabel('Connection Name').fill('Line DB');
      await page.getByLabel('Host IP Address').fill(destination.host);
      await page.getByLabel('Port').fill(String(destination.port));
      await page.getByLabel('Database Name').fill(destination.database);
      await page.getByLabel('Username').fill(destination.user);
      await page.locator('input[type=password]').first().fill(destination.password);
      await page.getByLabel('Schema').fill(destination.schema);
      await page.getByLabel('Table Name').fill(destination.table);
    };
    await page.waitForTimeout(800);
    await fill();
    await waitSaved('postgres', { host: destination.host, port: destination.port, table: destination.table, schema: destination.schema, database: destination.database, username: destination.user }, page, fill);
  }
  await page.waitForTimeout(2000);
  await page.reload();
  await page.waitForTimeout(3000);
  await page.getByTestId('step-nav-button-4').click();
  await page.waitForTimeout(2500);
}

async function showGroupList(page) {
  for (let attempt = 0; attempt < 4; attempt++) {
    if (await page.getByTestId('group-create').isVisible()) return;
    const close = page.getByTestId('group-editor-close');
    if (await close.isVisible()) await close.click();
    await page.waitForTimeout(1200);
  }
}

export async function createAndApplyGroup(page, name, letter, log = console.log, { extended = false } = {}) {
  await showGroupList(page);
  await page.getByTestId('group-create').click();
  await page.getByTestId('group-name').fill(name);
  const advanced = page.getByTestId('group-advanced');
  if ((await advanced.getAttribute('open')) === null) await advanced.locator('summary').click();
  await page.getByTestId('group-entity-column').fill('line');
  await page.getByTestId('group-provenance-column').fill('prov');
  await page.getByTestId('group-interval').fill('10');
  await page.getByTestId('group-member-search').fill(name);
  await page.getByTestId('group-bulk-include').click();
  await page.getByTestId('group-member-metadata-note').waitFor({ state: 'detached', timeout: 20000 });
  const columnFor = { temp: 'temperature', press: 'pressure', batch: 'batch', run: 'running' };
  if (extended) Object.assign(columnFor, { f32: 'reading_f32', f64: 'reading_f64', text: 'text_value' });
  for (const row of await page.locator('tr[data-testid^=group-member-row-]').all()) {
    const text = await row.innerText();
    const kind = Object.keys(columnFor).find((k) => text.includes(`.${letter}.${k}.`));
    if (!kind) continue;
    await row.locator('select').selectOption(columnFor[kind]);
    await row.locator('input[type=text]').fill(letter.toUpperCase());
  }
  await page.getByTestId('group-member-search').fill('');
  await page.getByTestId('group-save').click();
  await page.getByTestId('group-saved-note').waitFor({ timeout: 15000 });
  await page.getByTestId('group-readiness').waitFor({ timeout: 15000 });
  await page.waitForTimeout(1000);
  const readiness = (await page.getByTestId('group-readiness').innerText()).replace(/\s+/g, ' ');
  await page.getByTestId('group-apply').click();
  await page.waitForFunction(() => { const e = document.querySelector('[data-testid=group-applied-revision]'); return e && !/Not applied|尚未/.test(e.textContent); }, null, { timeout: 15000 });
  const applied = await page.getByTestId('group-applied-revision').innerText();
  log(`${name}: readiness=[${readiness}] applied_revision=${applied}`);
  await showGroupList(page);
  return { readiness, applied };
}

export async function activateDevices(page) {
  await page.getByRole('button', { name: /First Activate Devices/ }).click();
  await page.getByText('Activation confirmed').first().waitFor({ timeout: 20000 });
}

export function witness(extra) {
  return {
    recorded_at: new Date().toISOString(),
    source_sha: process.env.F_SOURCE_SHA ?? sh('git', ['rev-parse', 'HEAD'], { cwd: ROOT }),
    dirty_worktree: process.env.F_DIRTY_WORKTREE
      ? process.env.F_DIRTY_WORKTREE === 'true'
      : sh('git', ['status', '--porcelain'], { cwd: ROOT }) !== '',
    platform: `${os.platform()} ${os.arch()} ${os.release()}`,
    go: sh('go', ['version']),
    node: process.version,
    ...extra,
  };
}

export function cleanWork(work) {
  if (existsSync(work)) rmSync(work, { recursive: true, force: true });
  mkdirSync(work, { recursive: true });
}
