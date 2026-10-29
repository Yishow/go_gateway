import { expect, test, type Page } from '@playwright/test';
import {
  AUTHORING_VIEWPORTS, actionCalls, attachAuthoringEvidence, captureAuthoring, noAuthoringMutation,
  object, openAuthoring, selectDevice, testWritePanel, type AuthoringMock,
} from './recording-authoring.helpers';

const storageKey = 'studio.recording-test-write-recovery.v1';
async function submitTestWrite(page: Page, mock: AuthoringMock) {
  const panel = await testWritePanel(page);
  await expect(panel.getByTestId('group-test-write-preview')).toBeEnabled();
  await panel.getByTestId('group-test-write-preview').click();
  await expect(panel.getByTestId('group-test-write-confirm')).toBeEnabled();
  await panel.getByTestId('group-test-write-confirm').click();
  await expect.poll(() => actionCalls(mock, 'test-write-confirm').length).toBe(1);
  return panel;
}
function exactlyOneInsert(mock: AuthoringMock, repair = false) {
  expect(actionCalls(mock, 'test-write-preview')).toHaveLength(1);
  expect(actionCalls(mock, 'test-write-confirm')).toHaveLength(1);
  expect(mock.snapshot().insertCount).toBe(1);
  noAuthoringMutation(mock, ['schema-preview', 'test-write-preview', 'test-write-confirm', ...(repair ? ['cleanup-repair'] : [])]);
}

for (const viewport of AUTHORING_VIEWPORTS) {
  test.describe(`test-write recovery ${viewport.width}`, () => {
    test.use({ viewport, timezoneId: 'UTC' });

    test('AR-R01 lost ACK blocks replacement through reload, genuine browser Back and device changes', async ({ page }, info) => {
      const { mock } = await openAuthoring(page, 'test-write', { setup: value => value.loseNext('test-write-confirm') });
      let panel = await submitTestWrite(page, mock);
      await expect(panel.getByTestId('group-test-write-unconfirmed')).toBeVisible();
      await expect(panel.getByTestId('group-test-write-preview')).toBeDisabled();
      await captureAuthoring(page, info, 'test-write-lost-ack', panel.getByTestId('group-test-write-unconfirmed'));
      const original = structuredClone(actionCalls(mock, 'test-write-confirm')[0].body);
      await page.reload();
      const recovery = page.getByTestId('workspace-test-write-recovery');
      await expect(recovery).toContainText('O1');
      await page.getByTestId('fixture-leave').click();
      await expect(page.getByTestId('fixture-away')).toBeVisible();
      await page.goBack();
      await expect(recovery).toContainText('O1');
      await selectDevice(page, 'dev-2');
      await expect(recovery).toContainText('O1');
      await selectDevice(page, 'dev-1');
      panel = await testWritePanel(page);
      await expect(panel.getByTestId('group-test-write-preview')).toBeDisabled();
      await recovery.getByTestId('group-test-write-recovery-query-O1').click();
      await expect(recovery.getByTestId('recovery-O1-result')).toBeVisible();
      await expect(recovery.getByTestId('recovery-O1-write-phase')).toHaveText('Write not confirmed');
      await expect(recovery.getByTestId('recovery-O1-readback')).toHaveText('Readback not confirmed');
      expect(actionCalls(mock, 'test-write-confirm')[0].body).toEqual(original);
      exactlyOneInsert(mock);
      await captureAuthoring(page, info, 'test-write-recovered-after-back-device-reload', recovery);
      await attachAuthoringEvidence(info, mock);
    });

    test('AR-R02 server listing recovers O1 after local reference loss without exposing another workspace', async ({ page }, info) => {
      const { mock } = await openAuthoring(page, 'test-write', { setup: value => value.loseNext('test-write-confirm') });
      const panel = await submitTestWrite(page, mock);
      await expect(panel.getByTestId('group-test-write-unconfirmed')).toBeVisible();
      await page.evaluate(key => {
        localStorage.removeItem(key);
        localStorage.setItem(key, JSON.stringify([{ workspace_id: 'W2', group_id: 'FOREIGN-G', group_revision: 'FOREIGN-R',
          operation_id: 'FOREIGN-O', operation_kind: 'test_write', created_at: '2030-01-01T00:00:00Z' }]));
      }, storageKey);
      const priorQueries = actionCalls(mock, 'test-write-list').length;
      await page.reload();
      await expect.poll(() => actionCalls(mock, 'test-write-list').length).toBeGreaterThan(priorQueries);
      const recovery = page.getByTestId('workspace-test-write-recovery');
      await expect(recovery).toContainText('O1');
      await expect(page.getByText('FOREIGN-O', { exact: false })).toHaveCount(0);
      await recovery.getByTestId('group-test-write-recovery-query-O1').click();
      await expect(recovery.getByTestId('recovery-O1-result')).toBeVisible();
      expect(actionCalls(mock, 'operation-status').every(call => !call.url.includes('FOREIGN-O'))).toBe(true);
      exactlyOneInsert(mock);
      await captureAuthoring(page, info, 'server-ledger-recovery-local-reference-loss', recovery);
      await attachAuthoringEvidence(info, mock);
    });

    test('AR-R03 failed server read remains unknown and never authorizes replacement preview', async ({ page }, info) => {
      const { mock } = await openAuthoring(page, 'test-write', { setup: value => value.loseNext('test-write-confirm') });
      let panel = await submitTestWrite(page, mock);
      await expect(panel.getByTestId('group-test-write-unconfirmed')).toBeVisible();
      mock.rejectNext('test-write-list');
      mock.rejectNext('operation-status');
      await page.reload();
      const recovery = page.getByTestId('workspace-test-write-recovery');
      await expect(recovery).toContainText('not confirmed');
      panel = await testWritePanel(page);
      await expect(panel.getByTestId('group-test-write-preview')).toBeDisabled();
      await expect(recovery).toContainText('O1');
      exactlyOneInsert(mock);
      await captureAuthoring(page, info, 'test-write-ledger-unavailable-not-no-write', recovery);
      await attachAuthoringEvidence(info, mock, ['test-write-list', 'operation-status']);
    });

    for (const status of ['pending', 'running'] as const) {
      test(`AR-R04 ${status} remains separately queryable and blocks a new preview`, async ({ page }, info) => {
        const { mock } = await openAuthoring(page, 'test-write', { setup: value => value.setTestWriteOutcome(status) });
        const panel = await submitTestWrite(page, mock);
        await expect(panel.getByTestId('group-test-write-running')).toBeVisible();
        await expect(panel.getByTestId('group-test-write-preview')).toBeDisabled();
        await panel.getByTestId('group-test-write-check').click();
        await expect.poll(() => actionCalls(mock, 'operation-status').length).toBeGreaterThan(0);
        exactlyOneInsert(mock);
        await captureAuthoring(page, info, `test-write-${status}`, panel.getByTestId('group-test-write-result'));
        await attachAuthoringEvidence(info, mock);
      });
    }

    for (const cleanup of ['failed', 'unknown'] as const) {
      test(`AR-R05 verified write with ${cleanup} cleanup keeps phases separate and repairs the same O1 once`, async ({ page }, info) => {
        const { mock, holdNext } = await openAuthoring(page, 'test-write-sqlserver', { setup: value => value.setTestWriteOutcome('succeeded', 'written_verified', cleanup) });
        const panel = await submitTestWrite(page, mock);
        await expect(panel.getByTestId('group-test-write-write-phase')).toHaveText('Write committed');
        await expect(panel.getByTestId('group-test-write-readback')).toHaveText('Readback verified');
        await expect(panel.getByTestId('group-test-write-cleanup')).toHaveText(cleanup === 'failed'
          ? 'Cleanup failed; the test row may remain' : 'Unknown whether it was removed');
        await expect(panel.getByTestId('group-test-write-preview')).toBeDisabled();
        await captureAuthoring(page, info, `test-write-cleanup-${cleanup}`, panel.getByTestId('group-test-write-result'));
        const gate = holdNext('cleanup-repair');
        try {
          await panel.getByTestId('group-test-write-cleanup-repair').dblclick({ delay: 20 });
          await expect.poll(() => actionCalls(mock, 'cleanup-repair').length).toBe(1);
          const first = object(actionCalls(mock, 'test-write-confirm')[0].body);
          expect(actionCalls(mock, 'cleanup-repair')[0].body).toEqual({ ...first, action: 'cleanup_repair' });
          gate.release();
          await expect(panel.getByTestId('group-test-write-cleanup')).toHaveText('Test row removed');
          await expect(panel.getByTestId('group-test-write-cleanup-repair')).toHaveCount(0);
          expect(mock.snapshot().repairCount).toBe(1);
          exactlyOneInsert(mock, true);
        } finally { gate.release(); }
        await attachAuthoringEvidence(info, mock);
      });
    }

    test('AR-R06 after reload original-token loss gives read-only cleanup guidance, not invented repair authority', async ({ page }, info) => {
      const { mock } = await openAuthoring(page, 'test-write-sqlserver', { setup: value => value.setTestWriteOutcome('succeeded', 'written_verified', 'failed') });
      await submitTestWrite(page, mock);
      await page.reload();
      const recovery = page.getByTestId('workspace-test-write-recovery');
      await expect(recovery).toContainText('O1');
      await recovery.getByTestId('group-test-write-recovery-query-O1').click();
      await expect(recovery.getByTestId('recovery-O1-write-phase')).toHaveText('Write committed');
      await expect(recovery.getByTestId('recovery-O1-readback')).toHaveText('Readback verified');
      await expect(recovery.getByTestId('recovery-O1-manual-guidance')).toContainText('cannot authorize cleanup repair');
      const panel = await testWritePanel(page);
      await expect(panel.getByTestId('group-test-write-cleanup-repair')).toHaveCount(0);
      await expect(panel.getByTestId('group-test-write-preview')).toBeDisabled();
      expect(actionCalls(mock, 'cleanup-repair')).toHaveLength(0);
      exactlyOneInsert(mock);
      await captureAuthoring(page, info, 'cleanup-query-only-after-token-loss', recovery);
      await attachAuthoringEvidence(info, mock);
    });
  });
}
