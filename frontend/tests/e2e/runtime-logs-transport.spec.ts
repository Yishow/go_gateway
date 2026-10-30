import { expect, test } from '@playwright/test';
import { domainRecord, runtimeLogsSnapshot } from './fixtures/runtime-logs-state';
import { LOG_LANGUAGES, LOG_VIEWPORTS, applyDevice, captureLogs, disclose, fact, failStream, finishLogs,
  logRows, logWords, openLogs, readyLogs, sendFrame, sendRecord, snapshotCalls, streams } from './runtime-logs.helpers';

for (const viewport of LOG_VIEWPORTS) for (const language of LOG_LANGUAGES) {
  test.describe(`runtime logs transport ${language} ${viewport.width}`, () => {
    test.use({ viewport, timezoneId: 'UTC' });
    const words = logWords(language);

    test('RL-T01 logs_local_only403 stops automatic retries and never requests a proxy', async ({ page }, info) => {
      const journey = await openLogs(page, language, 'normal', mock => mock.failNext('snapshot', 403));
      await expect(page.getByRole('status')).toHaveText(words.denied);
      await expect(page.getByRole('alert')).toContainText(language === 'en' ? 'Automatic retries have stopped.' : '已停止自動重試。');
      await expect(page.getByRole('button', { name: words.pause, exact: true })).toBeDisabled();
      await page.clock.fastForward(31_000);
      expect(snapshotCalls(journey.mock)).toHaveLength(1); expect(snapshotCalls(journey.mock)[0].status).toBe(403); expect(await streams(page)).toEqual([]);
      expect(journey.mock.ledger.every(call => call.path === '/api/v1/system/logs')).toBe(true);
      await captureLogs(page, info, 'RL-T01', language); await finishLogs(info, journey);
    });

    test('RL-T02 A9007199254740993 late snapshot and cancelled stream cannot populate B', async ({ page }, info) => {
      const a = domainRecord('domain.source.rejected', { fields: { device_id: 'A' } });
      const b = domainRecord('domain.source.rejected', { sequence: '7', fields: { device_id: 'B' } });
      const journey = await openLogs(page, language, 'normal', mock => mock.setSnapshot(runtimeLogsSnapshot([a])));
      const first = await readyLogs(page, language);
      expect(new URL((await streams(page))[first].url).searchParams.get('after')).toBe('instance-a:9007199254740993');
      const late = journey.holdNext(call => call.path === '/api/v1/system/logs' && call.params.device_id === 'A');
      await applyDevice(page, language, 'A', 'RULE'); await late.reached;
      journey.mock.setSnapshot(runtimeLogsSnapshot([b]));
      await applyDevice(page, language, 'B'); const current = await readyLogs(page, language, 'B');
      expect(new URL((await streams(page))[current].url).searchParams.get('after')).toBe('instance-a:7');
      late.release(); await late.settled;
      await expect.poll(() => journey.failedRequests.some(url => new URL(url).searchParams.get('device_id') === 'A')).toBe(true);
      expect(await sendRecord(page, first, { ...a, sequence: '9007199254740994' })).toBe(false);
      expect((await streams(page))[first].aborted).toBe(true);
      expect(await sendRecord(page, current, { ...b, sequence: '8' })).toBe(true);
      await expect(logRows(page)).toHaveCount(2);
      for (const row of [logRows(page).first(), logRows(page).nth(1)]) {
        await disclose(row, language); await expect(fact(row, words.device)).toHaveText('B');
      }
      expect(snapshotCalls(journey.mock).map(call => call.params.device_id ?? '')).toEqual(['', 'A', 'B']);
      expect(snapshotCalls(journey.mock).slice(1).every(call => call.params.rule_id === 'RULE')).toBe(true);
      expect(new URL((await streams(page))[current].url).searchParams.get('rule_id')).toBe('RULE');
      expect((await streams(page)).every(stream => !stream.url.includes('9007199254740992'))).toBe(true);
      await captureLogs(page, info, 'RL-T02', language); await finishLogs(info, journey);
    });

    test('RL-T03 pause resume and clear preserve server records and exact consumed replay progress', async ({ page }, info) => {
      const journey = await openLogs(page, language); const first = await readyLogs(page, language);
      await page.getByRole('button', { name: words.pause, exact: true }).click();
      await expect(page.getByRole('status')).toHaveText(words.paused); expect((await streams(page))[first].aborted).toBe(true);
      await expect(logRows(page)).toHaveCount(1); await page.clock.fastForward(5000); expect(await streams(page)).toHaveLength(1);
      await page.getByRole('button', { name: words.resume, exact: true }).click(); const resumed = await readyLogs(page, language);
      expect(resumed).toBeGreaterThan(first); expect(snapshotCalls(journey.mock)).toHaveLength(1);
      const replay = domainRecord('domain.delivery.unknown', { sequence: '9007199254740994' });
      expect(await sendRecord(page, resumed, replay)).toBe(true);
      expect(await sendFrame(page, resumed, { type: 'heartbeat', metadata: runtimeLogsSnapshot([replay]), cursor: 'instance-a:9007199254740994' })).toBe(true);
      await expect(logRows(page)).toHaveCount(2); await page.getByRole('button', { name: words.clear, exact: true }).click();
      await expect(logRows(page)).toHaveCount(0); expect(snapshotCalls(journey.mock)).toHaveLength(1);
      await page.getByRole('button', { name: words.pause, exact: true }).click();
      await page.getByRole('button', { name: words.resume, exact: true }).click(); const afterClear = await readyLogs(page, language);
      expect(new URL((await streams(page))[afterClear].url).searchParams.get('after')).toBe('instance-a:9007199254740994');
      expect(await sendRecord(page, afterClear, replay)).toBe(true); // Old replay is transported, but not redisplayed.
      const next = { ...replay, sequence: '9007199254740995' };
      expect(await sendRecord(page, afterClear, next)).toBe(true); await expect(logRows(page)).toHaveCount(1);
      await expect(logRows(page).first()).toContainText(words.unknown);
      expect(journey.mock.ledger).toHaveLength(1);
      await captureLogs(page, info, 'RL-T03', language); await finishLogs(info, journey);
    });

    test('RL-T04 offline reconnect and retention gap remain unknown under actionable filtering', async ({ page }, info) => {
      const journey = await openLogs(page, language); const first = await readyLogs(page, language);
      await failStream(page, first); await expect(page.getByRole('status')).toHaveText(words.reconnecting);
      await expect(logRows(page)).toHaveCount(1); await expect(logRows(page).first()).toContainText(words.retrying);
      await page.clock.fastForward(1300); const recoveredTransport = await readyLogs(page, language);
      expect(recoveredTransport).toBeGreaterThan(first); expect(snapshotCalls(journey.mock)).toHaveLength(1);
      expect(await sendFrame(page, recoveredTransport, { type: 'gap', reason: 'retention', cursor: '' })).toBe(true);
      await expect(page.getByText(words.historyUnknown, { exact: false })).toBeVisible();
      await expect(page.getByRole('heading', { name: words.recovered, exact: true })).toHaveCount(0);
      // The next server-filtered snapshot also declares incomplete retention;
      // switching filters cannot fabricate recovery from the absence of a row.
      const retained = runtimeLogsSnapshot([domainRecord('domain.delivery.retrying')], { evicted: 2, capture_dropped: 1 });
      journey.mock.setSnapshot(retained);
      await page.getByLabel(words.actionable, { exact: true }).check();
      await page.getByRole('button', { name: words.apply, exact: true }).click();
      await expect.poll(() => snapshotCalls(journey.mock).length).toBe(2); await readyLogs(page, language, undefined, retained);
      await expect(page.getByText(words.historyUnknown, { exact: false })).toBeVisible();
      await expect(logRows(page)).toHaveCount(1); await expect(logRows(page).first()).toContainText(words.retrying);
      expect(snapshotCalls(journey.mock).at(-1)?.params.actionable).toBe('true');
      expect(new URL((await streams(page)).at(-1)!.url).searchParams.get('actionable')).toBe('true');
      await expect(page.getByRole('heading', { name: words.recovered, exact: true })).toHaveCount(0);
      await captureLogs(page, info, 'RL-T04', language); await finishLogs(info, journey);
    });
  });
}
