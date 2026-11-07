import { releaseSnapshotReply } from './helpers/recording-home-snapshot-ledger';
import { test, expect, HOME_URL, readyHome, rendered, settledRequests, sourceCount, observation, home, openEditor,
  returnHome, filterPageTwo, expectFiltersRestored, openSourceRange, assertPageDemand } from './helpers/recording-home-snapshot-browser';

// Real App/router/components served by Vite; API/test data are simulated. Not embedded/backend proof.
for (const count of [60, 1000]) test.describe(`${count} matching recordings, real routed frontend / mock API`, () => {
  test.use({ groupCount: count });
  test('editor return and genuine Back/Forward restore filters, page, scroll and focus with one compact GET per home visit', async ({ page, snapshot: api }, info) => {
    await page.goto(HOME_URL); await readyHome(page, api); expect(sourceCount(api)).toBe(1);
    await assertPageDemand(page, api, 0, true); await filterPageTwo(page, api, count);
    const trigger = page.locator('[data-home-focus="record-A-0034"]');
    await trigger.scrollIntoViewIfNeeded(); await trigger.focus();
    const scroll = await page.getByTestId('home-collection-scroll').evaluate(element => element.scrollTop);
    expect(scroll).toBeGreaterThan(0);
    await openEditor(page, api, 'A-0034');
    const editorURL = page.url(); expect(new URL(editorURL).hash).toBe('#snapshot-proof');
    let checkpoint = api.ledger.length, compactBefore = sourceCount(api);
    const held = api.holdNext(); await returnHome(page, api); const returning = await held.entered;
    expect(returning.workspaceId).toBe('A');
    await expect(observation(page)).toHaveAttribute('data-source-status', 'refreshing');
    await expectFiltersRestored(page, count); await expect(trigger).toBeFocused();
    await expect(trigger).toBeInViewport({ ratio: 1 });
    await expect.poll(async () => Math.abs(await page.getByTestId('home-collection-scroll').evaluate(element => element.scrollTop) - scroll)).toBeLessThanOrEqual(1);
    expect(sourceCount(api)).toBe(compactBefore + 1);
    await page.screenshot({ path: info.outputPath(`return-refreshing-${count}-MOCK.png`), fullPage: true });
    held.release(); await readyHome(page, api); await assertPageDemand(page, api, checkpoint, true);
    await expect(trigger).toBeInViewport({ ratio: 1 });
    expect(sourceCount(api)).toBe(compactBefore + 1);
    // Router history is exercised, not page.goto/reload or a standalone component mount.
    api.setPhase('editor'); await page.goBack(); await expect(page).toHaveURL(editorURL);
    await expect(page.getByTestId('group-editor')).toBeVisible(); await settledRequests(api);
    checkpoint = api.ledger.length; compactBefore = sourceCount(api);
    const forward = api.holdNext(); api.setPhase('home'); await page.goForward(); await forward.entered;
    await expect(observation(page)).toHaveAttribute('data-source-status', 'refreshing');
    await expectFiltersRestored(page, count); await expect(trigger).toBeFocused();
    await expect(trigger).toBeInViewport({ ratio: 1 });
    await expect.poll(async () => Math.abs(await page.getByTestId('home-collection-scroll').evaluate(element => element.scrollTop) - scroll)).toBeLessThanOrEqual(1);
    forward.release(); await readyHome(page, api); await assertPageDemand(page, api, checkpoint, true);
    await expect(trigger).toBeInViewport({ ratio: 1 });
    expect(sourceCount(api)).toBe(compactBefore + 1); expect(new URL(page.url()).hash).toBe('#snapshot-proof');
    expect(api.ledger.filter(call => call.mutation)).toEqual([]);
  });
});

test('ordinary source Close, connection Cancel and Escape do not revalidate the still-mounted home', async ({ page, snapshot: api }) => {
  await page.goto(HOME_URL); await readyHome(page, api); const before = sourceCount(api);
  await openSourceRange(page, api);
  await page.getByTestId('home-range-done').click(); await expect(page.getByRole('dialog')).toHaveCount(0);
  await expect(page.getByTestId('home-device-points-D-A')).toBeFocused();
  await rendered(page); await settledRequests(api); expect(sourceCount(api)).toBe(before);
  api.setPhase('source-modal'); await page.locator('[data-home-focus="summary-connection-D-A"]').click();
  await expect(page.getByTestId('recording-source-close')).toHaveText('Cancel');
  await expect(page.getByTestId('recording-source-close')).toBeEnabled(); await settledRequests(api);
  await page.getByTestId('recording-source-close').click(); await expect(page.getByRole('dialog')).toHaveCount(0);
  await expect(page.locator('[data-home-focus="summary-connection-D-A"]')).toBeFocused();
  await rendered(page); await settledRequests(api); expect(sourceCount(api)).toBe(before);
  await openSourceRange(page, api); await page.keyboard.press('Escape'); await expect(page.getByRole('dialog')).toHaveCount(0);
  await expect(page.getByTestId('home-device-points-D-A')).toBeFocused();
  api.setPhase('home'); await readyHome(page, api); expect(sourceCount(api)).toBe(before);
  expect(api.ledger.filter(call => call.mutation)).toEqual([]);
});

test('source save ACK independently triggers exactly one held compact refresh; closing adds none', async ({ page, snapshot: api }, info) => {
  await page.goto(HOME_URL); await readyHome(page, api); const before = sourceCount(api);
  await openSourceRange(page, api); const originalRevision = api.snapshot().rules[0].revision_id;
  api.allowRuleSave('R-A', 'SAVED_PREFIX_');
  await page.getByTestId('rule-prefix-input').fill('SAVED_PREFIX_');
  const hold = api.holdNext(); await page.getByTestId('home-range-save').click(); await hold.entered;
  // Existing owner acknowledges first; held home observation cannot substitute for that ACK.
  await expect(page.getByTestId('home-range-saved')).toBeVisible();
  expect(api.snapshot().rules[0].revision_id).not.toBe(originalRevision);
  await expect(observation(page)).toHaveAttribute('data-source-status', 'refreshing');
  expect(sourceCount(api)).toBe(before + 1);
  const saveCalls = api.ledger.filter(call => call.mutation);
  expect(saveCalls).toHaveLength(1); expect(saveCalls[0]).toMatchObject({ method: 'PUT', action: 'save-rule', status: 200 });
  await page.getByTestId('home-range-done').click(); await expect(page.getByRole('dialog')).toHaveCount(0);
  await expect(page.getByTestId('home-device-points-D-A')).toBeFocused(); await rendered(page);
  expect(sourceCount(api)).toBe(before + 1); await expect(observation(page)).toHaveAttribute('data-source-status', 'refreshing');
  hold.release(); api.setPhase('home'); await readyHome(page, api); expect(sourceCount(api)).toBe(before + 1);
  await page.screenshot({ path: info.outputPath('source-save-fresh-MOCK.png'), fullPage: true });
});

for (const status of [409, 503]) test(`fresh cached presentation cannot authorize a source save when current locator proof returns ${status}`, async ({ page, snapshot: api }) => {
  await page.goto(HOME_URL); await readyHome(page, api); await openSourceRange(page, api);
  const original = api.snapshot(), before = sourceCount(api);
  await page.getByTestId('rule-prefix-input').fill('UNSENT_'); api.failNext('locator:2', status);
  const checkpoint = api.ledger.length;
  await page.getByTestId('home-range-save').click();
  await expect(page.getByRole('dialog').getByRole('alert')).toContainText('The saved range changed or could not be verified');
  await expect(page.getByTestId('rule-prefix-input')).toHaveValue('UNSENT_'); await settledRequests(api); await rendered(page);
  expect(api.ledger.slice(checkpoint).some(call => call.action === 'locator:2' && call.status === status)).toBe(true);
  expect(api.ledger.filter(call => call.mutation)).toEqual([]); expect(api.snapshot()).toEqual(original);
  expect(sourceCount(api)).toBe(before); await expect(page.locator('body')).not.toContainText('RAW_FIXTURE_DIAGNOSTIC');
  // Recover by reopening the existing owner; no readiness, activation, apply or test-write request is sent.
  await page.getByTestId('home-range-done').click(); await page.getByTestId('home-source-discard').click();
  await expect(page.getByRole('dialog')).toHaveCount(0); await openSourceRange(page, api);
  await expect(page.getByTestId('rule-prefix-input')).toHaveValue(original.rules[0].naming_prefix);
  await expect(page.getByTestId('home-range-save')).toBeDisabled();
  await page.getByTestId('home-range-done').click(); await expect(page.getByRole('dialog')).toHaveCount(0);
  await rendered(page); expect(api.ledger.filter(call => call.mutation)).toEqual([]);
});

test('failed return keeps observed filters and devices; two explicit read-only retries each add one compact GET', async ({ page, snapshot: api }) => {
  await page.goto(HOME_URL); await readyHome(page, api); await filterPageTwo(page, api, 60);
  await openEditor(page, api, 'A-0026'); const before = sourceCount(api);
  api.failNext('sources'); const held = api.holdNext(); await returnHome(page, api); await held.entered;
  await expect(observation(page)).toHaveAttribute('data-source-status', 'refreshing'); held.release();
  await expect(observation(page)).toHaveAttribute('data-source-status', 'error'); await expectFiltersRestored(page);
  await settledRequests(api); await rendered(page); expect(sourceCount(api)).toBe(before + 1);
  await expect(page.getByTestId('home-device-summary')).toContainText('Snapshot A D');
  await expect(page.locator('body')).not.toContainText('RAW_FIXTURE_DIAGNOSTIC');
  api.unknownSourcesNext();
  await home(page).getByRole('button', { name: 'Refresh', exact: true }).click();
  await expect(observation(page)).toHaveAttribute('data-source-status', 'stale'); await settledRequests(api); await rendered(page);
  expect(sourceCount(api)).toBe(before + 2); await expectFiltersRestored(page);
  await home(page).getByRole('button', { name: 'Refresh', exact: true }).click(); await readyHome(page, api);
  expect(sourceCount(api)).toBe(before + 3); await expectFiltersRestored(page);
  expect(api.ledger.filter(call => call.mutation)).toEqual([]);
});

test('editor point discovery retains its full source reads without reintroducing raw transport on home', async ({ page, snapshot: api }) => {
  await page.goto(HOME_URL); await readyHome(page, api); const before = sourceCount(api);
  await openEditor(page, api, 'A-0001'); await settledRequests(api);
  const checkpoint = api.ledger.length;
  await page.getByTestId('recording-points-open').click();
  await expect.poll(() => api.ledger.slice(checkpoint).filter(call => call.action.startsWith('raw:')).length).toBe(8);
  await settledRequests(api); await expect(page.getByRole('dialog')).toBeVisible();
  expect(api.ledger.slice(checkpoint).filter(call => call.action.startsWith('raw:')).every(call => call.phase === 'editor')).toBe(true);
  expect(sourceCount(api)).toBe(before);
  await page.getByRole('dialog').getByRole('button', { name: 'Cancel', exact: true }).click();
  await expect(page.getByRole('dialog')).toHaveCount(0); await returnHome(page, api); await readyHome(page, api);
  expect(sourceCount(api)).toBe(before + 1);
  expect(api.ledger.filter(call => call.phase === 'home' && call.action.startsWith('raw:'))).toEqual([]);
  expect(api.ledger.filter(call => call.mutation)).toEqual([]);
});

for (const order of ['A-first', 'B-first'] as const) test(`A → B → A uses singleton bootstrap reads and fences old held replies (${order})`, async ({ page, snapshot: api }, info) => {
  await page.goto(HOME_URL); await readyHome(page, api);
  await openEditor(page, api, 'A-0001'); const oldA = api.holdNext('sources:A');
  await returnHome(page, api); const callA = await oldA.entered;
  await expect(observation(page)).toHaveAttribute('data-source-status', 'refreshing');
  await openEditor(page, api, 'A-0001');
  api.switchWorkspace('B'); const oldB = api.holdNext('sources:B');
  // Advance Date only past the production 5-second bootstrap staleTime. No cache/test setter.
  await page.clock.setFixedTime(new Date('2030-01-01T00:00:06Z'));
  const beforeB = api.ledger.length; await returnHome(page, api); const callB = await oldB.entered;
  await expect.poll(() => api.ledger.slice(beforeB).some(call => call.action === 'bootstrap' && call.workspaceId === 'B')).toBe(true);
  await expect(home(page)).toHaveAttribute('data-home-settling', 'false');
  await expect(page.locator('[data-home-focus="record-B-0001"]')).toBeVisible();
  await expect(page.getByTestId('home-device-summary')).not.toContainText('Snapshot A');
  await openEditor(page, api, 'B-0001'); api.switchWorkspace('A', 'A2');
  await page.clock.setFixedTime(new Date('2030-01-01T00:00:12Z'));
  const beforeA2 = api.ledger.length; await returnHome(page, api);
  await expect(observation(page)).toHaveAttribute('data-source-status', 'fresh');
  await expect(page.getByTestId('home-device-summary')).toContainText('Snapshot A A2');
  await settledRequests(api, [callA, callB]);
  const current = api.ledger.slice(beforeA2).filter(call => call.action === 'sources' && call.workspaceId === 'A' && call.revision === 'A2' && call.status === 200);
  expect(current).toHaveLength(1);
  const digest = (call: typeof callA) => (call.response as { data: { source_digest?: string } }).data.source_digest;
  expect(digest(current[0])).not.toBe(digest(callA));
  expect(digest(current[0])).toMatch(/^[a-f0-9]{64}$/);
  const assertCurrent = async () => {
    await expect(page.getByTestId('home-device-summary')).toContainText('Snapshot A A2');
    await expect(page.getByTestId('home-device-summary')).not.toContainText('Snapshot B');
    await expect(observation(page)).toHaveAttribute('data-source-status', 'fresh');
  };
  const replies = order === 'A-first' ? [[oldA, callA], [oldB, callB]] as const : [[oldB, callB], [oldA, callA]] as const;
  for (const [gate, call] of replies) {
    await releaseSnapshotReply(gate.release, async () => {
      await expect.poll(() => ['finished', 'cancelled', 'not-delivered'].includes(call.delivery)).toBe(true);
    }, () => rendered(page), assertCurrent);
  }
  await settledRequests(api); await rendered(page); await assertCurrent();
  // A browser cancellation is valid isolation evidence, never falsely labelled delivered-and-ignored.
  await info.attach('late-reply-delivery-classification', { contentType: 'application/json', body: Buffer.from(JSON.stringify({
    current: current[0], old: [callA, callB].map(call => ({ sequence: call.sequence, workspace: call.workspaceId,
      delivery: call.delivery, browserTransportCompleted: call.delivery === 'finished', failure: call.failure })),
    ignoredAbortUnitEvidenceRemainsSeparate: true }, null, 2)) });
  await page.screenshot({ path: info.outputPath(`scope-A2-${order}-MOCK.png`), fullPage: true });
  expect(api.ledger.filter(call => call.mutation)).toEqual([]);
});
