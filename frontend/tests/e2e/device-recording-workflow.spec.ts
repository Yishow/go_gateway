import { expect, test } from '@playwright/test';
import fs from 'node:fs';
import type { WriteGroup } from '../../src/types/studioV2WriteGroup';
import type { RecordingStartRequest, RecordingStartOperation } from '../../src/types/studioV2RecordingStart';
import {
  api, catalog, committedReadback, groups, identifier, loadFixture, readSQL,
  recordingOperations, runtimeVersions,
} from './fixtures/device-recording-workflow';

const fixture = loadFixture();
test('one device recording journey: keyboard, purposes, real SQL and lost-reply recovery', async ({ page }, info) => {
  test.skip(!fixture, 'NOT_RUN: explicit owned normal gateway fixture is required');
  if (!fixture) throw new Error('Fixture required');
  test.setTimeout(300_000);
  page.setDefaultTimeout(15_000);
  page.setDefaultNavigationTimeout(15_000);
  const f = fixture;
  const ui = (id: string) => page.getByTestId(id);
  const tile = (id: string) => page.locator('[data-testid="basic-recording-device-tile"][data-device-id="' + id + '"]');
  const open = async () => {
    await page.goto(f.frontend + '/studio/v2?step=4');
    await expect(ui('circuit-profile-basic-panel')).toBeVisible();
  };
  const screenshot = async (name: string) => {
    await page.screenshot({ path: info.outputPath(name + '.png'), fullPage: true });
    await info.attach(name, { path: info.outputPath(name + '.png'), contentType: 'image/png' });
  };
  const attachJSON = async (name: string, value: unknown) => {
    const file = info.outputPath(name + '.json');
    fs.writeFileSync(file, JSON.stringify(value, null, 2));
    await info.attach(name, { path: file, contentType: 'application/json' });
  };
  const noOverflow = async () => {
    expect(await page.evaluate(() => document.documentElement.scrollWidth > innerWidth)).toBe(false);
  };
  const confirm = async (kind: string) => {
    const panel = ui('circuit-prepared-group-' + kind);
    const response = page.waitForResponse(r => r.url().endsWith('/schema-apply') && r.request().method() === 'POST');
    await panel.getByTestId('group-schema-apply').click();
    expect((await response).status()).toBe(200);
    await expect.poll(async () => (await groups(f)).filter(g => g.recording_profile?.kind === kind)
      .every(g => Boolean(g.destination.schema_digest))).toBe(true);
  };

  // Main seeded only Step 1–3 through formal APIs; this test refuses a used recording workspace.
  expect((await api<{ id: string }>(f, '')).id).toBe(f.workspace_id);
  const devices = await api<{ id: string; name: string }[]>(f, '/devices');
  expect(devices.map(d => d.id).sort()).toEqual([f.devices.a, f.devices.b].sort());
  expect(devices[0].name).toBe(devices[1].name);
  expect(await groups(f)).toHaveLength(0);
  expect(recordingOperations(f)).toHaveLength(0);
  expect(catalog(f)).toHaveLength(0);

  const mutations: string[] = [];
  page.on('request', request => {
    if (['POST', 'PUT', 'PATCH', 'DELETE'].includes(request.method()) && request.url().includes('/api/')) {
      mutations.push(request.method() + ' ' + new URL(request.url()).pathname);
    }
  });
  await page.addInitScript(() => localStorage.setItem('i18nextLng', 'zh-TW'));
  await open();
  for (const viewport of [{ width: 1024, height: 720 }, { width: 1600, height: 1000 }]) {
    await page.setViewportSize(viewport);
    await tile(f.devices.a).click();
    await expect(ui('circuit-profile-interval-instantaneous')).toHaveValue(viewport.width === 1024 ? '15' : '47');
    await ui('circuit-profile-interval-instantaneous').fill('47');
    await tile(f.devices.a).focus();
    await page.keyboard.press('Tab');
    await expect(tile(f.devices.b)).toBeFocused();
    await page.keyboard.press('Space');
    await expect(tile(f.devices.b)).toHaveAttribute('aria-pressed', 'true');
    await ui('circuit-profile-interval-instantaneous').fill('33');
    await tile(f.devices.b).focus();
    await page.keyboard.press('Shift+Tab');
    await expect(tile(f.devices.a)).toBeFocused();
    await page.keyboard.press('Enter');
    await expect(tile(f.devices.a)).toHaveAttribute('aria-pressed', 'true');
    await expect(ui('circuit-profile-interval-instantaneous')).toHaveValue('47');
    await page.reload();
    await expect(tile(f.devices.a)).toHaveAttribute('aria-pressed', 'true');
    await expect(ui('circuit-profile-interval-instantaneous')).toHaveValue('47');
    await tile(f.devices.b).click();
    await expect(ui('circuit-profile-interval-instantaneous')).toHaveValue('33');
    await tile(f.devices.a).click();
    await expect(page.locator('select[name="device"],select[aria-label*="device" i]')).toHaveCount(0);
    await expect(ui('circuit-profile-sources-summary')).toContainText('32');
    await expect(ui('circuit-profile-sources-toggle')).toHaveAttribute('aria-expanded', 'false');
    await noOverflow();
    await screenshot('drafts-' + viewport.width);
  }
  expect(mutations).toEqual([]);
  expect(await groups(f)).toHaveLength(0);

  await ui('circuit-profile-circuit-name').fill('電錶 A');
  await ui('circuit-profile-interval-instantaneous').fill('15');
  await ui('circuit-profile-purpose-wide').check();
  await ui('circuit-profile-sources-toggle').click();
  const list = ui('circuit-profile-member-list');
  expect(await list.evaluate(el => el.scrollHeight > el.clientHeight)).toBe(true);
  await list.hover();
  await page.mouse.wheel(0, 200);
  await expect.poll(() => list.evaluate(el => el.scrollTop)).toBeGreaterThan(0);
  const counter = f.mappings.find(m => m.device_id === f.devices.a && m.target_type === 'uint64');
  expect(counter).toBeDefined();
  const counterId = counter!.tag_id;
  await ui('circuit-profile-role-' + counterId).selectOption('counter_raw');
  await ui('circuit-profile-sources-toggle').click();
  await expect(ui('circuit-counter-confirmation')).toBeVisible();
  await expect(ui('circuit-profile-prepare')).toBeDisabled();
  await ui('circuit-confirm-counter-semantics').selectOption('counter');
  await ui('circuit-confirm-counter-unit').fill('kWh');
  await ui('circuit-confirm-counter-submit').click();
  await expect(ui('circuit-counter-confirmation')).toHaveCount(0);
  await ui('circuit-profile-toggle-cumulative').check();
  await expect(ui('circuit-profile-interval-cumulative')).toHaveValue('10');
  await ui('circuit-profile-interval-cumulative').fill('1');
  await ui('circuit-profile-sources-toggle').click();
  for (const mapping of f.mappings.filter(m => m.device_id === f.devices.a && m.tag_id !== counterId)) {
    const box = ui('circuit-profile-member-cum-' + mapping.tag_id);
    await box.evaluate(el => el.scrollIntoView({ block: 'center', inline: 'nearest' }));
    await box.uncheck();
  }
  await expect(ui('circuit-profile-member-cum-' + counterId)).toBeChecked();
  await ui('circuit-profile-sources-toggle').click();
  await expect(ui('circuit-profile-prepare')).toBeEnabled();
  await ui('circuit-profile-prepare').click();
  await expect(ui('group-schema-preview-result')).toHaveCount(2);
  expect(catalog(f)).toHaveLength(0); // preview must not issue DDL
  await confirm('instantaneous');
  await confirm('cumulative');

  let request: RecordingStartRequest | undefined;
  let committed: RecordingStartOperation | undefined;
  await page.route('**/studio-v2/workspace/recording-start', async route => {
    request = route.request().postDataJSON() as RecordingStartRequest;
    expect(request.request_id).toBeTruthy();
    expect(request.device_ids).toEqual([f.devices.a]);
    const response = await route.fetch();
    expect(response.status()).toBe(200);
    committed = (await response.json() as { data: RecordingStartOperation }).data;
    expect(committed.status).toBe('succeeded');
    await route.abort('failed'); // real backend commit, only its HTTP reply is lost
  }, { times: 1 });
  await ui('circuit-profile-start').click();
  await expect(ui('basic-recording-retry')).toBeVisible();
  await expect(ui('circuit-profile-start')).toBeDisabled();
  expect(request).toBeDefined();
  await expect.poll(() => committed?.status).toBe('succeeded');
  const aVersions = runtimeVersions(f);
  expect(aVersions).toHaveLength(2);
  const beforeRecovery = recordingOperations(f);
  expect(beforeRecovery).toHaveLength(1);
  expect(beforeRecovery[0].operation_id).toBe(committed!.operation_id);
  await tile(f.devices.b).click();
  await expect(ui('circuit-profile-interval-instantaneous')).toHaveValue('33');
  await expect(ui('basic-recording-retry')).toHaveCount(0);
  await page.reload();
  await expect(tile(f.devices.b)).toHaveAttribute('aria-pressed', 'true');
  await tile(f.devices.a).click();
  await expect(ui('basic-recording-retry')).toBeVisible();
  const resent = page.waitForRequest(r => r.url().endsWith('/recording-start') && r.method() === 'POST');
  await ui('basic-recording-retry').click();
  expect((await resent).postDataJSON()).toEqual(request);
  await page.waitForURL('**/studio/runtime**');
  expect(recordingOperations(f)).toEqual(beforeRecovery);
  expect(runtimeVersions(f)).toEqual(aVersions);
  console.info('Lost reply recovered: original request, one ledger operation, unchanged applied versions.');
  let saved = await groups(f);
  expect(saved).toHaveLength(2); // B was unprepared and did not block A or gain a recording group.
  expect(saved.every(g => g.members.every(m => m.device_id === f.devices.a))).toBe(true);

  await open();
  await tile(f.devices.b).click();
  await ui('circuit-profile-circuit-name').fill('電錶 B');
  await ui('circuit-profile-interval-instantaneous').fill('15');
  await ui('circuit-profile-purpose-long').check();
  await ui('circuit-profile-sources-toggle').click();
  const bMappings = f.mappings.filter(m => m.device_id === f.devices.b);
  expect(bMappings).toHaveLength(2);
  for (const mapping of bMappings) await ui('circuit-profile-member-inst-' + mapping.tag_id).check();
  await ui('circuit-profile-sources-toggle').click();
  await expect(ui('circuit-profile-prepare')).toBeEnabled();
  const bEnsureResponse = page.waitForResponse(r => r.url().endsWith('/circuit-profiles/ensure') && r.request().method() === 'POST');
  await ui('circuit-profile-prepare').click();
  const bResponse = await bEnsureResponse;
  const bIntent = bResponse.request().postDataJSON() as { group: WriteGroup };
  await attachJSON('B-ensure-intent', { request: bIntent, response: await bResponse.json() });
  expect(bIntent.group.members.map(m => m.device_id)).toEqual([f.devices.b, f.devices.b]);
  const aCircuitIds = saved.flatMap(g => g.members.map(m => m.circuit_id));
  expect(aCircuitIds).not.toContain(bIntent.group.members[0].circuit_id);
  expect(bResponse.status()).toBe(200);
  await expect(ui('group-schema-preview-result')).toHaveCount(1);
  await confirm('instantaneous');
  await ui('circuit-profile-start').click();
  await page.waitForURL('**/studio/runtime**');
  saved = await groups(f);
  expect(saved).toHaveLength(3);
  expect(recordingOperations(f)).toHaveLength(2);
  console.info('Both purposes started; waiting for real UTC windows and matching SQL receipts.');
  const groupFor = (device: string, kind: string): WriteGroup => {
    const matches = saved.filter(g => g.members.every(m => m.device_id === device) && g.recording_profile?.kind === kind);
    expect(matches).toHaveLength(1);
    expect(matches[0].applied_revision).toBeTruthy();
    return matches[0];
  };
  const a = groupFor(f.devices.a, 'instantaneous');
  const cumulative = groupFor(f.devices.a, 'cumulative');
  const b = groupFor(f.devices.b, 'instantaneous');
  expect(a.layout_version).toBe('managed_wide_profile_v1');
  expect(b.layout_version).toBe('managed_long_v1');
  expect(a.recording_profile?.sample_reduction).toBe('mean');
  expect(b.recording_profile?.sample_reduction).toBe('mean');
  expect(cumulative.members).toHaveLength(1);
  for (const group of [a, cumulative, b]) {
    await expect.poll(() => Boolean(committedReadback(f, group)), { timeout: 100_000 }).toBe(true);
  }
  const wide = committedReadback(f, a)!;
  for (const [type, value] of [['uint16', '20'], ['float32', '1.5'], ['bool', '1'], ['string', 'ON'], ['uint64', '9007199254740992']]) {
    const mapping = f.mappings.find(m => m.device_id === f.devices.a && m.target_type === type);
    expect(mapping).toBeDefined();
    const member = a.members.find(m => m.tag_id === mapping!.tag_id);
    expect(member).toBeDefined();
    const exact = readSQL(f.target_db, 'SELECT CAST(' + identifier(member!.target_column) + ' AS TEXT) AS value FROM ' +
      identifier(a.destination.table_name) + ' WHERE record_id=?', [wide.effect.record_id]);
    expect(exact[0].value).toBe(value);
  }
  const baseline = JSON.parse(String(committedReadback(f, cumulative)!.row.provenance)) as { counter: { status: string; delta?: string } }[];
  expect(baseline[0].counter.status).toBe('baseline');
  expect(baseline[0].counter.delta).toBeUndefined();
  const long = committedReadback(f, b)!;
  expect(long.row.group_revision).toBe(b.applied_revision);
  expect(long.row.value_decimal).toBe('20');
  await attachJSON('SQL-same-revision', { a: wide, cumulative: committedReadback(f, cumulative), b: long });

  await open();
  await tile(f.devices.a).click();
  const evidenceToggle = ui('circuit-status-sources-' + a.id + '-toggle');
  await expect(evidenceToggle).toHaveAttribute('aria-expanded', 'false');
  await evidenceToggle.click();
  const gauge = a.members.find(m => m.value_role === 'gauge')!;
  await expect(ui('circuit-status-aggregation-' + gauge.tag_id).first()).toContainText('未知');
  await expect(ui('circuit-status-valid-' + gauge.tag_id).first()).not.toHaveText('0');
  for (const viewport of [{ width: 1024, height: 720 }, { width: 1600, height: 1000 }]) {
    await page.setViewportSize(viewport);
    await ui('circuit-status-instantaneous').scrollIntoViewIfNeeded();
    await noOverflow();
    await screenshot('SQL-evidence-' + viewport.width);
  }
});
