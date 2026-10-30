import { expect, test } from '@playwright/test';
import { domainRecord, legacyRecord, runtimeLogsSnapshot, RUNTIME_LOGS_TIME } from './fixtures/runtime-logs-state';
import { LOG_LANGUAGES, LOG_VIEWPORTS, captureLogs, disclose, fact, finishLogs, locatorCalls,
  logRows, logWords, openLogs, readyLogs, sendFrame, sendRecord, sendRecords, snapshotCalls } from './runtime-logs.helpers';

// Formal browser preparation only. Mock rows cannot establish native SQL producer
// correctness or writer nonblocking behavior; those remain the current Go proofs.
for (const viewport of LOG_VIEWPORTS) for (const language of LOG_LANGUAGES) {
  test.describe(`runtime logs evidence ${language} ${viewport.width}`, () => {
    test.use({ viewport, timezoneId: 'UTC' });
    const words = logWords(language);

    test('RL-E01 G/R1 retry becomes matching receipt E without relabeling history', async ({ page }, info) => {
      const retry = domainRecord('domain.delivery.retrying', { fields: { reference_id: 'E' } });
      const journey = await openLogs(page, language, 'normal', mock => mock.setSnapshot(runtimeLogsSnapshot([retry])));
      const stream = await readyLogs(page, language);
      await expect(logRows(page)).toHaveCount(1); await expect(logRows(page).first()).toContainText(words.retrying);
      const recovered = domainRecord('domain.delivery.recovered', { sequence: '9007199254740994', fields: { reference_id: 'E', effect_key: 'E' } });
      expect(await sendRecord(page, stream, recovered)).toBe(true);
      await expect(logRows(page)).toHaveCount(2); await expect(logRows(page).nth(1)).toContainText(words.recovered);
      for (const row of [logRows(page).nth(0), logRows(page).nth(1)]) {
        await disclose(row, language); await expect(fact(row, words.revision)).toHaveText('R1');
        await expect(fact(row, words.reference)).toHaveText('E');
      }
      const receipt = logRows(page).nth(1).getByRole('region', { name: words.receipt, exact: true });
      await expect(fact(receipt, words.effect)).toHaveText('E'); await expect(fact(receipt, words.receiptRevision)).toHaveText('R1');
      expect(snapshotCalls(journey.mock)).toHaveLength(1); expect(locatorCalls(journey.mock)).toHaveLength(0);
      expect(journey.mock.ledger).toHaveLength(1); // Expanding every row does not hydrate settings.
      await captureLogs(page, info, 'RL-E01', language); await finishLogs(info, journey);
    });

    test('RL-E02 accepted queued uncertain and mismatched R0 receipt never become R1 commitment', async ({ page }, info) => {
      // A queued E has no commitment event. Intake is separate, uncertainty is
      // explicit, and a hostile R0-vs-R1 claimed receipt must fail closed in UI.
      const records = [domainRecord('domain.group.applied', { fields: { reference_id: 'E' } }),
        domainRecord('domain.delivery.unknown', { sequence: '9007199254740994', fields: { reference_id: 'E' } }),
        domainRecord('domain.delivery.first-confirmed', { sequence: '9007199254740995', fields: { reference_id: 'E', effect_key: 'E', group_revision: 'R0' } })];
      const journey = await openLogs(page, language, 'normal', mock => mock.setSnapshot(runtimeLogsSnapshot(records)));
      await readyLogs(page, language); await expect(logRows(page)).toHaveCount(3);
      await expect(logRows(page).nth(1)).toContainText(words.unknown); await expect(logRows(page).nth(2)).toContainText(words.unavailable);
      for (const row of [logRows(page).first(), logRows(page).nth(1), logRows(page).nth(2)]) await disclose(row, language);
      await expect(page.getByRole('heading', { name: words.firstCommit, exact: true })).toHaveCount(0);
      await expect(page.getByRole('region', { name: words.receipt, exact: true })).toHaveCount(0);
      await expect(logRows(page).nth(2).getByRole('link')).toHaveCount(0);
      expect(journey.mock.ledger).toHaveLength(1);
      await captureLogs(page, info, 'RL-E02', language); await finishLogs(info, journey);
    });

    test('RL-E03 capture and file loss preserve own-effect receipt and expose incomplete history', async ({ page }, info) => {
      const committed = domainRecord('domain.delivery.first-confirmed', { fields: { reference_id: 'E', effect_key: 'E' } });
      const journey = await openLogs(page, language, 'normal', mock => mock.setSnapshot(runtimeLogsSnapshot([committed], {
        capture_dropped: 7, file_dropped: 2, disk_errors: 1, sink_health: 'degraded',
      })));
      await readyLogs(page, language, undefined, runtimeLogsSnapshot([committed], { capture_dropped: 7, file_dropped: 2, disk_errors: 1, sink_health: 'degraded' }));
      await expect(page.getByText(words.historyUnknown, { exact: false })).toBeVisible();
      await expect(page.getByText(words.gap, { exact: false })).toBeVisible();
      await disclose(logRows(page).first(), language);
      await expect(fact(logRows(page).first(), words.effect)).toHaveText('E');
      await expect(logRows(page).first()).toContainText(words.firstCommit);
      expect(locatorCalls(journey.mock)).toHaveLength(0);
      await captureLogs(page, info, 'RL-E03', language); await finishLogs(info, journey);
    });

    test('RL-E04 hostile raw candidates and unknown codes never leak or execute markup', async ({ page }, info) => {
      const hostile = [domainRecord('domain.delivery.blocked', { message: '<img src=x onerror="alert(1)">password_fixture', fields: { raw_error: 'password_fixture' } }),
        { ...domainRecord('domain.delivery.blocked', { sequence: '9007199254740994' }), code: 'domain.password_fixture' },
        { ...legacyRecord('runtime.started', { sequence: '9007199254740995', message: '<script>alert(1)</script>private.example' }), code: 'future.password_fixture' }];
      const journey = await openLogs(page, language, 'normal', mock => mock.setSnapshot(runtimeLogsSnapshot(hostile)));
      const dialogs: string[] = []; page.on('dialog', async dialog => { dialogs.push(dialog.message()); await dialog.dismiss(); });
      await readyLogs(page, language); await expect(logRows(page)).toHaveCount(3);
      await expect(page.locator('tbody')).not.toContainText(/password_fixture|private\.example|raw_error|<script|<img/);
      await expect(page.locator('tbody script, tbody img, tbody iframe')).toHaveCount(0);
      await expect(page.locator('tbody a')).toHaveCount(0); expect(dialogs).toEqual([]);
      // A wire-level nested payload is rejected by the real transport parser,
      // rather than cast into a legal scalar RuntimeLogRecord by the fixture.
      const stream = await readyLogs(page, language);
      const nested = JSON.stringify({ ...domainRecord(), sequence: '9007199254740996', fields: { credential: { password: 'password_fixture' } } });
      await page.evaluate(({ stream, frame }) => {
        const fixture = window as Window & { __runtimeLogsStream: { send(index: number, frame: string): boolean } };
        fixture.__runtimeLogsStream.send(stream, frame);
      }, { stream, frame: `event: log\nid: instance-a:9007199254740996\ndata: ${nested}\n\n` });
      await expect(page.getByRole('status')).toHaveText(language === 'en' ? 'Failed' : '載入失敗');
      await expect(page.locator('body')).not.toContainText(/password_fixture|private\.example/);
      await page.getByRole('button', { name: language === 'en' ? 'Retry connection' : '重試連線', exact: true }).click();
      const retried = await readyLogs(page, language);
      expect(await sendRecord(page, retried, domainRecord('domain.delivery.retrying', {
        sequence: '9007199254740997', message: 'oversized_secret'.repeat(600),
      }))).toBe(true);
      await expect(page.getByRole('status')).toHaveText(language === 'en' ? 'Failed' : '載入失敗');
      await expect(page.locator('body')).not.toContainText(/password_fixture|private\.example|oversized_secret/);
      expect(snapshotCalls(journey.mock)).toHaveLength(2);
      await captureLogs(page, info, 'RL-E04', language); await finishLogs(info, journey);
    });

    test('RL-E05 G/R1 occurrences3 and G/R2 first fault remain independent after R1 recovery', async ({ page }, info) => {
      const records = [domainRecord('domain.delivery.retrying', { fields: { first_seen_at: '2026-10-29T09:00:00Z', last_seen_at: RUNTIME_LOGS_TIME, occurrences: '3' } }),
        domainRecord('domain.delivery.retrying', { sequence: '9007199254740994', fields: { entity_revision: 'R2', group_revision: 'R2', first_seen_at: RUNTIME_LOGS_TIME, last_seen_at: RUNTIME_LOGS_TIME, occurrences: '1' } })];
      const journey = await openLogs(page, language, 'revisions', mock => mock.setSnapshot(runtimeLogsSnapshot(records)));
      const stream = await readyLogs(page, language);
      for (const [index, revision, count] of [[0, 'R1', '3'], [1, 'R2', '1']] as const) {
        const row = logRows(page).nth(index); await disclose(row, language);
        await expect(fact(row, words.revision)).toHaveText(revision); await expect(fact(row, words.occurrences)).toHaveText(count);
      }
      expect(await sendRecord(page, stream, domainRecord('domain.delivery.recovered', { sequence: '9007199254740995' }))).toBe(true);
      await expect(logRows(page)).toHaveCount(3); await expect(fact(logRows(page).nth(1), words.revision)).toHaveText('R2');
      await expect(logRows(page).nth(1)).toContainText(words.retrying); await expect(fact(logRows(page).first(), words.occurrences)).toHaveText('3');
      await captureLogs(page, info, 'RL-E05', language); await finishLogs(info, journey);
    });

    test('RL-E06 browser record and byte bounds remove old rows without inventing capture counts', async ({ page }, info) => {
      const journey = await openLogs(page, language, 'empty'); const stream = await readyLogs(page, language);
      const records = Array.from({ length: 1005 }, (_, index) => domainRecord('domain.delivery.retrying', {
        sequence: String(index + 1), fields: { group_id: `G${index + 1}` },
      }));
      // Real encoded frames and consumer append loop; no buffer or React state injection.
      expect(await sendRecords(page, stream, records)).toBe(1005);
      await expect(logRows(page)).toHaveCount(1000);
      await expect(page.getByText(words.truncated, { exact: false })).toBeVisible();
      await expect(logRows(page).first()).toContainText('G6');
      await page.getByRole('button', { name: words.clear, exact: true }).click();
      const large = Array.from({ length: 400 }, (_, index) => legacyRecord('runtime.started', {
        sequence: String(2000 + index), message: 'x'.repeat(7000),
      }));
      expect(await sendRecords(page, stream, large)).toBe(400);
      await expect(page.locator('tbody tr').last()).toHaveAttribute('data-log-cursor', 'instance-a:2399');
      const retained = await logRows(page).count();
      expect(retained).toBeGreaterThan(0); expect(retained).toBeLessThan(400);
      expect(retained * 7000).toBeLessThanOrEqual(2 * 1024 * 1024); expect(snapshotCalls(journey.mock)).toHaveLength(1);
      await captureLogs(page, info, 'RL-E06', language); await finishLogs(info, journey);
    });

    test('RL-E07 process I2 reset does not replay historical E as a new SQL commitment', async ({ page }, info) => {
      const committed = domainRecord('domain.delivery.first-confirmed', { instance_id: 'I1', sequence: '9', fields: { reference_id: 'E', effect_key: 'E' } });
      const journey = await openLogs(page, language, 'normal', mock => mock.setSnapshot(runtimeLogsSnapshot([committed], { instance_id: 'I1', latest_cursor: 'I1:9', oldest_cursor: 'I1:9' })));
      const stream = await readyLogs(page, language, undefined, runtimeLogsSnapshot([committed], { instance_id: 'I1', latest_cursor: 'I1:9', oldest_cursor: 'I1:9' }));
      await expect(logRows(page)).toHaveCount(1);
      expect(await sendFrame(page, stream, { type: 'reset', reason: 'process_restart', cursor: '' })).toBe(true);
      expect(await sendFrame(page, stream, { type: 'handshake', metadata: runtimeLogsSnapshot([], { instance_id: 'I2', oldest_cursor: '', latest_cursor: 'I2:0' }), cursor: 'I2:0' })).toBe(true);
      await expect(logRows(page)).toHaveCount(0); await expect(page.getByText(words.reset, { exact: false })).toBeVisible();
      await expect(page.getByText(words.historyUnknown, { exact: false })).toBeVisible();
      await expect(page.getByRole('heading', { name: words.firstCommit, exact: true })).toHaveCount(0);
      await captureLogs(page, info, 'RL-E07', language); await finishLogs(info, journey);
    });
  });
}
