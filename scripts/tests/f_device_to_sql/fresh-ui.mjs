// Fresh UI helpers for F runs. Mutations in this module are browser actions;
// fetch and SQL helpers are deliberately GET/SELECT-only observations.
import { chromium } from '../../../frontend/node_modules/playwright/index.mjs';
import { mkdirSync, existsSync, statSync, readFileSync, rmSync } from 'node:fs';
import { join, resolve } from 'node:path';
import { createHash } from 'node:crypto';
import { createConnection } from 'node:net';
import {
  ROOT, addDevice, addRule, boundedFetch, sh, startProcess, stopProcesses, writeRegisters,
} from './lib.mjs';
import { preflightSQLiteDestination } from './destination-preflight.mjs';

// Every runner records the binary hash it used; override these through the
// environment so a run never silently uses a stale build.
export const DEFAULT_GATEWAY_BINARY = process.env.F_GATEWAY_BINARY || '/private/tmp/gw-F-test-ui-production-final-v6';
export const DEFAULT_SIMULATOR_BINARY = process.env.F_SIMULATOR_BINARY || '/tmp/gw-f-early-sqlite-20261004/f-sim';
export const DEFAULT_EVIDENCE_DIR = join(ROOT, 'docs/plans/studio-v2-flow-completion/evidence-f');

export function monotonicNs() {
  return process.hrtime.bigint().toString();
}

export function makeFreshRun(runId) {
  if (!/^[A-Za-z0-9][A-Za-z0-9_-]{0,80}$/.test(runId)) throw new Error('invalid fresh run id');
  const work = `/tmp/gw-f-${runId}`;
  if (resolve(work) !== work) throw new Error('fresh work path escaped /tmp');
  try {
    mkdirSync(work);
  } catch (error) {
    if (error?.code === 'EEXIST' || existsSync(work)) throw new Error(`fresh run namespace already exists: ${work}`);
    throw error;
  }
  mkdirSync(DEFAULT_EVIDENCE_DIR, { recursive: true });
  return { run_id: runId, work, gateway_db: join(work, 'gateway.db'), destination_db: join(work, 'destination.db') };
}

export function assertFreshPortFree(port, timeoutMs = 3000) {
  if (!Number.isInteger(port) || port < 1024 || port > 65535) throw new Error('fresh port must be an unprivileged TCP port');
  return new Promise((resolvePort, rejectPort) => {
    const socket = createConnection({ host: '127.0.0.1', port });
    let settled = false;
    const finish = (error, value) => {
      if (settled) return;
      settled = true;
      clearTimeout(timer);
      socket.destroy();
      if (error) rejectPort(error);
      else resolvePort(value);
    };
    const timer = setTimeout(() => finish(new Error(`fresh port ${port} did not prove free before timeout`)), timeoutMs);
    socket.once('connect', () => finish(new Error(`fresh port ${port} is already occupied`)));
    socket.once('error', (error) => {
      if (error?.code === 'ECONNREFUSED') finish(null, { port, status: 'free' });
      else finish(new Error(`fresh port ${port} probe failed: ${error?.code ?? 'unknown'}`));
    });
  });
}

export function binaryIdentity(path) {
  return { path, sha256: createHash('sha256').update(readFileSync(path)).digest('hex') };
}

export function startFreshGateway({ work, binary = DEFAULT_GATEWAY_BINARY, port = 3380 }) {
  return startProcess(work, 'gateway', binary, [], {
    GATEWAY_DB_PATH: join(work, 'gateway.db'), HOST: '127.0.0.1', PORT: String(port),
  });
}

export function startFreshSimulator({ work, name = 'sim-a', binary = DEFAULT_SIMULATOR_BINARY, port, registers }) {
  const registerFile = writeRegisters(work, `${name}-registers`, registers);
  return startProcess(work, name, binary, ['-port', String(port), '-registers', registerFile]);
}

export async function waitForFreshGateway(base, timeoutMs = 30_000) {
  const deadline = Date.now() + timeoutMs;
  while (Date.now() < deadline) {
    try {
      const timeout = Math.min(3000, deadline - Date.now());
      const [apiResponse, uiResponse] = await Promise.all([
        boundedFetch(`${base}/api/v1/datalink/studio-v2/workspace`, {}, timeout),
        boundedFetch(`${base}/studio/v2`, {}, timeout),
      ]);
      if (uiResponse.status === 200 && apiResponse.status === 200) {
        const body = await apiResponse.json();
        if (body?.data?.id && body.data.readiness_summary != null) return body.data;
      }
    } catch { /* process is still starting */ }
    await new Promise((resolveWait) => setTimeout(resolveWait, 300));
  }
  throw new Error(`fresh gateway did not start: ${base}`);
}

export async function observeJSON(base, path) {
  const apiPath = path.startsWith('/api/v1/datalink/') ? path
    : path.startsWith('/studio-v2/') ? `/api/v1/datalink${path}` : '';
  if (!apiPath) throw new Error('observation path must be datalink API');
  const response = await boundedFetch(`${base}${apiPath}`, {}, 10_000);
  const body = await response.json();
  if (!response.ok) throw new Error(`observation ${apiPath} failed: ${response.status}`);
  return body?.data ?? body;
}

export async function waitForObservation(base, path, predicate, timeoutMs = 20_000) {
  const deadline = Date.now() + timeoutMs;
  let last;
  while (Date.now() < deadline) {
    last = await observeJSON(base, path);
    if (predicate(last)) return last;
    await new Promise((resolveWait) => setTimeout(resolveWait, 500));
  }
  throw new Error(`observation did not reach expected state: ${path} ${JSON.stringify(last).slice(0, 300)}`);
}

export async function openFreshUI({ port = 3380, width = 1440, height = 1400 }) {
  const base = `http://127.0.0.1:${port}`;
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage({ viewport: { width, height } });
  const apiTraffic = [];
  const pageErrors = [];
  page.on('pageerror', (error) => pageErrors.push(error.message));
  page.on('response', async (response) => {
    const url = new URL(response.url());
    if (!url.pathname.startsWith('/api/')) return;
    const request = response.request();
    const item = { method: request.method(), path: url.pathname, status: response.status() };
    if (apiTraffic.length < 512) {
      apiTraffic.push(item);
    }
  });
  const tOpen = monotonicNs();
  await page.goto(`${base}/studio/v2`);
  await page.getByTestId('workbench-v2-root').waitFor({ timeout: 20_000 });
  return { browser, page, base, t_open: tOpen, apiTraffic, pageErrors };
}

export async function setupFreshDevice(page, { name, port }) {
  await addDevice(page, name, String(port));
  await page.getByTestId('success-readiness-card').waitFor({ timeout: 20_000 });
  await page.getByTestId('btn-continue-step1').click();
  await page.getByTestId('step2-rule-container').waitFor({ timeout: 20_000 });
}

export async function setupFreshRules(page, rules) {
  for (const rule of rules) {
    await addRule(page, rule.name, String(rule.start), rule.count, rule.dataType, rule.prefix, rule.existingTab);
  }
  const continueButton = page.getByTestId('continue-step3-btn');
  try {
    await page.waitForFunction(() => !document.querySelector('[data-testid="continue-step3-btn"]')?.hasAttribute('disabled'), null, { timeout: 20_000 });
  } catch {
    const detail = await page.getByTestId('merged-point-table-container').innerText().catch(() => 'merged point table unavailable');
    throw new Error(`fresh Step 2 is not ready: ${detail.replace(/\s+/g, ' ').slice(-1200)}`);
  }
  await continueButton.click();
  await page.getByTestId('step3-mapping-container').waitFor({ timeout: 20_000 });
}

/** Add every fresh device before leaving Step 1, then create each device's
 * source rules at the same addresses. The UI remains the authority for both
 * mutations; callers only supply the disposable loopback ports. */
export async function setupFreshMultiDeviceSources(page, devices, rules) {
  for (const device of devices) await addDevice(page, device.name, String(device.port));
  await page.getByTestId('btn-continue-step1').click();
  await page.getByTestId('step2-rule-container').waitFor({ timeout: 20_000 });
  const lastTab = new Map();
  for (const rule of rules) {
    for (const device of devices) {
      const name = `${device.name} ${rule.name}`;
      await addRule(page, name, String(rule.start), rule.count, rule.dataType,
        `${device.prefix ?? device.name[0]}_${rule.prefix}`, lastTab.get(device.name));
      lastTab.set(device.name, name);
    }
  }
  const continueButton = page.getByTestId('continue-step3-btn');
  await page.waitForFunction(() => !document.querySelector('[data-testid="continue-step3-btn"]')?.hasAttribute('disabled'), null, { timeout: 30_000 });
  await continueButton.click();
  await page.getByTestId('step3-mapping-container').waitFor({ timeout: 20_000 });
}

function mappingList(data) {
  if (Array.isArray(data)) return data;
  return data?.mappings ?? data?.items ?? [];
}

async function waitForMappingFieldSaved(page, base, { pointId, fieldTestId, expected, apiPredicate }) {
  const savedData = await waitForObservation(base, '/studio-v2/workspace/mappings', apiPredicate);
  const field = page.getByTestId(fieldTestId);
  const saveState = page.getByTestId(`mapping-save-state-${pointId}`);
  const waitForSaved = async () => {
    await page.waitForFunction(({ fieldTestId: id, expectedValue, stateTestId }) => {
      const input = document.querySelector(`[data-testid="${id}"]`);
      const state = document.querySelector(`[data-testid="${stateTestId}"]`);
      return input?.value === expectedValue && /Saved|已儲存/.test(state?.textContent ?? '');
    }, { fieldTestId, expectedValue: expected, stateTestId: `mapping-save-state-${pointId}` }, { timeout: 20_000 });
    if (await field.inputValue() !== expected || !/Saved|已儲存/.test(await saveState.innerText())) {
      throw new Error(`mapping ${pointId} did not retain saved ${fieldTestId}`);
    }
  };
  await waitForSaved();
  // A mapping query invalidation can hydrate the row after the mutation response. Require
  // one short second observation so the next field is never edited over stale state.
  await page.waitForTimeout(150);
  await waitForSaved();
  return mappingList(savedData);
}

export async function persistFreshMappings(page, base, { tagPrefix = 'f', targetTypes = {}, targetTypeFor } = {}) {
  const rows = page.locator('[data-testid^="mapping-row-"]');
  await rows.first().waitFor({ timeout: 20_000 });
  const count = await rows.count();
  const tags = [];
  for (let index = 0; index < count; index += 1) {
    const row = rows.nth(index);
    const pointId = (await row.getAttribute('data-testid')).replace('mapping-row-', '');
    const address = (await row.locator('td').nth(1).innerText()).trim();
    const tag = `${tagPrefix}.${index + 1}`;
    await page.getByTestId(`input-tag-key-${pointId}`).fill(tag);
    await page.getByTestId(`input-tag-key-${pointId}`).press('Tab');
    await waitForMappingFieldSaved(page, base, {
      pointId,
      fieldTestId: `input-tag-key-${pointId}`,
      expected: tag,
      apiPredicate: (data) => {
        const saved = mappingList(data);
        return saved.some((item) => item.tag_key === tag && String(item.address) === address);
      },
    });
    await page.getByTestId(`input-display-name-${pointId}`).fill(tag);
    await page.getByTestId(`input-display-name-${pointId}`).press('Tab');
    await waitForMappingFieldSaved(page, base, {
      pointId,
      fieldTestId: `input-display-name-${pointId}`,
      expected: tag,
      apiPredicate: (data) => {
        const saved = mappingList(data);
        return saved.some((item) => item.tag_key === tag && String(item.address) === address && item.display_name === tag);
      },
    });
    const type = typeof targetTypeFor === 'function'
      ? targetTypeFor({ index, address, pointId })
      : targetTypes[index];
    let finalMappings;
    if (type) {
      const selector = page.getByTestId(`select-target-type-${pointId}`);
      await selector.selectOption(type);
      if (await selector.inputValue() !== type) throw new Error(`mapping target type control did not accept ${type}`);
      await selector.press('Tab');
      await waitForMappingFieldSaved(page, base, {
        pointId,
        fieldTestId: `select-target-type-${pointId}`,
        expected: type,
        apiPredicate: (data) => {
        const saved = mappingList(data);
          return saved.some((item) => item.tag_key === tag && String(item.address) === address && item.display_name === tag && item.target_type === type);
        },
      });
      finalMappings = await observeJSON(base, '/studio-v2/workspace/mappings').then(mappingList);
    }
    const persisted = (finalMappings ?? await observeJSON(base, '/studio-v2/workspace/mappings').then(mappingList))
      .find((item) => item.tag_key === tag && String(item.address) === address);
    if (!persisted) throw new Error(`mapping ${tag} did not expose persisted identity`);
    tags.push({
      ui_local_point_id: pointId,
      tag_key: tag,
      address,
      persisted_point_id: persisted.point_id,
      tag_id: persisted.tag_id,
      device_id: persisted.device_id,
      rule_id: persisted.rule_id,
    });
  }
  await waitForObservation(base, '/studio-v2/workspace/mappings', (data) => {
    const saved = mappingList(data);
    return tags.every((tag) => saved.some((item) => {
      return item.tag_key === tag.tag_key && String(item.address) === tag.address && item.point_id === tag.persisted_point_id;
    }));
  });
  await page.getByTestId('btn-continue').click();
  await page.getByTestId('step-nav-button-4').waitFor({ timeout: 20_000 });
  await page.waitForTimeout(800);
  return tags;
}

function field(page, english, chinese) {
  return page.getByLabel(new RegExp(`${english}|${chinese}`, 'i')).first();
}

export async function configureFreshSQLite(page, base, { path }) {
  await page.getByTestId('step-nav-button-4').click();
  await page.getByRole('button', { name: /^SQLite/i }).click();
  await page.waitForTimeout(800);
  const nameField = field(page, 'Connection Name', '連線名稱');
  const databaseField = field(page, 'DB File Path', '資料庫檔案路徑');
  const waitConfig = async (predicate, input, expected, label) => {
    const deadline = Date.now() + 20_000;
    let last;
    while (Date.now() < deadline) {
      last = await observeJSON(base, '/studio-v2/workspace/database-config');
      if (predicate(last)) break;
      await page.waitForTimeout(500);
    }
    if (!last || !predicate(last)) {
      const observed = last && {
        name: last.name,
        database: last.database,
        table: last.table,
        setup_revision: last.setup_revision,
      };
      throw new Error(`SQLite ${label} did not save requested value: ${JSON.stringify({ expected, observed })}`);
    }
    const saved = last;
    if (await input.inputValue() !== expected) throw new Error('SQLite connector response was overwritten in the UI');
    await page.waitForTimeout(150);
    if (await input.inputValue() !== expected) throw new Error('SQLite connector field did not remain saved');
    return saved;
  };
  await nameField.fill('F fresh SQLite');
  await nameField.press('Tab');
  await waitConfig((data) => data?.name === 'F fresh SQLite', nameField, 'F fresh SQLite', 'name');
  await databaseField.fill(path);
  await databaseField.press('Tab');
  await waitConfig((data) => data?.database === path, databaseField, path, 'database path');
  const savedConfig = await observeJSON(base, '/studio-v2/workspace/database-config');
  if (await page.getByTestId('input-table-name').count()) throw new Error('obsolete connector table setting is still present');
  await page.reload();
  await page.getByTestId('step-nav-button-4').click();
  await page.getByTestId('basic-recording-panel').waitFor({ timeout: 20_000 });
  return savedConfig;
}

export async function openFreshAdvanced(page) {
  const details = page.getByTestId('step4-advanced-recording');
  if ((await details.getAttribute('open')) === null) await details.locator('summary').click();
  await page.getByTestId('group-section').waitFor({ timeout: 20_000 });
}

export async function createFreshManagedGroup(page, { name = 'F managed group' } = {}) {
  await openFreshAdvanced(page);
  await page.getByTestId('group-create').click();
  await page.getByTestId('group-name').fill(name);
  await page.getByTestId('group-storage-managed').check();
  await page.getByTestId('group-bulk-include').click();
  await page.getByTestId('group-save').click();
  await page.getByTestId('group-saved-note').waitFor({ timeout: 20_000 });
  await page.getByTestId('group-readiness').waitFor({ timeout: 20_000 });
  const groups = await observeJSON(new URL(page.url()).origin, '/studio-v2/workspace/write-groups');
  const matches = groups.groups?.filter((group) => group.name === name) ?? [];
  if (matches.length !== 1) throw new Error(`expected one exact advanced group ${name}, found ${matches.length}`);
  return matches[0];
}

/** Explicit Basic preparation owns the canonical group; DDL still requires preview and confirmation. */
export async function prepareFreshBasicGroup(page, base, deviceId) {
  if (deviceId) await page.getByTestId('basic-recording-device').selectOption(deviceId);
  const tStart = monotonicNs();
  await page.getByTestId('basic-recording-prepare').click();
  await page.getByTestId('group-schema-panel').waitFor({ timeout: 30_000 });
  const devices = await observeJSON(base, '/studio-v2/workspace/devices');
  const device = (Array.isArray(devices) ? devices : devices.devices ?? []).find((item) => item.id === deviceId) ??
    (Array.isArray(devices) ? devices : devices.devices ?? [])[0];
  const groups = await observeJSON(base, '/studio-v2/workspace/write-groups');
  const basic = (groups.groups ?? []).filter((group) => group.basic_managed_device_id === device?.id);
  if (basic.length !== 1) throw new Error(`expected one Basic canonical group, found ${basic.length}`);
  return { t_start: tStart, first_status_text: 'Basic group prepared; DDL awaits explicit preview and confirmation', group: basic[0] };
}

export async function previewAndApplyFreshSchema(page, { beforeConfirm } = {}) {
  await page.getByTestId('group-schema-preview').click();
  await page.getByTestId('group-schema-preview-result').waitFor({ timeout: 20_000 });
  const previewText = await page.getByTestId('group-schema-preview-result').innerText();
  if (beforeConfirm) await beforeConfirm();
  await page.getByTestId('group-schema-apply').click();
  await page.getByTestId('group-schema-operation').waitFor({ timeout: 30_000 });
  await page.waitForFunction(() => /succeeded|成功/.test(document.querySelector('[data-testid="group-schema-operation"]')?.textContent ?? ''), null, { timeout: 30_000 });
  return { preview_text: previewText, operation_text: await page.getByTestId('group-schema-operation').innerText() };
}

export async function startFreshBasic(page) {
  const tStart = monotonicNs();
  await page.getByTestId('basic-recording-start').click();
  await page.getByTestId('basic-recording-status').waitFor({ timeout: 30_000 });
  await page.waitForFunction(() => /succeeded|成功|partial|部分|failed|失敗/.test(document.querySelector('[data-testid="basic-recording-status"]')?.textContent ?? ''), null, { timeout: 120_000 });
  return { t_start: tStart, status_text: await page.getByTestId('basic-recording-status').innerText() };
}

export async function ensureFreshSidebarCollapsed(page) {
  const rail = page.getByTestId('sidebar-rail-aside');
  if ((await rail.getAttribute('data-collapsed')) === 'true') return;
  for (const shortcut of ['Control+b', 'Meta+b']) {
    await page.keyboard.press(shortcut);
    await page.waitForTimeout(100);
    if ((await rail.getAttribute('data-collapsed')) === 'true') return;
  }
  throw new Error('fresh UI sidebar did not collapse with keyboard shortcut');
}

export function freshSQLiteObservation(path, sql) {
  if (!existsSync(path)) throw new Error(`SQLite destination is still missing: ${path}`);
  statSync(path);
  return sh('sqlite3', ['-readonly', '-cmd', '.timeout 15000', '-separator', '|', `file:${path}?mode=ro`, sql]);
}

export function preflightFreshSQLite(path, work) {
  return preflightSQLiteDestination({ path, work });
}

export async function closeFreshUI(browser) {
  await browser?.close();
}

/** A process cleanup witness is valid only after a normal, unforced exit. */
export function isFreshCleanExit(status) {
  return status?.state === 'exited' && status.exit_code === 0 && !status.signal &&
    status.kill_sent !== true && !status.spawn_error && !status.stop_error;
}

/**
 * Removes a run-owned work directory without throwing, so the caller can still
 * write its evidence. A directory that remains is reported as not removed.
 */
export function removeOwnedWork(work) {
  let error;
  try { rmSync(work, { recursive: true, force: true }); } catch (cause) { error = String(cause?.message ?? cause); }
  const removed = !existsSync(work);
  return removed ? { removed } : { removed, error: error ?? 'work directory still exists' };
}

export { stopProcesses };
