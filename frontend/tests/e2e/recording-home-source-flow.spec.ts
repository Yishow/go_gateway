import { expect, test } from '@playwright/test';
import { locatorRead, nativeDevice, nativeRule } from '../unit/workbench-v2/studio-locator-fixtures';

for (const viewport of [{ width: 1440, height: 900 }, { width: 320, height: 568 }]) {
  test(`home keeps one loading view until its saved source editor is ready at ${viewport.width}px`, async ({ page }, info) => {
    await page.setViewportSize(viewport);
    let releaseDevice!: () => void, releaseList!: () => void, releaseRule!: () => void;
    const deviceGate = new Promise<void>(resolve => { releaseDevice = resolve; });
    const listGate = new Promise<void>(resolve => { releaseList = resolve; });
    const ruleGate = new Promise<void>(resolve => { releaseRule = resolve; });
    const requests: string[] = [];
    let opened = false;
    const workspace = { ...locatorRead().workspace, ordered_device_ids: ['D'] };
    const device = { ...nativeDevice(), connection_config: '{"host":"fixture.invalid","port":502,"slave_id":1,"timeout":1500}' };
    await page.route('**/api/**', async route => {
      const request = route.request(), url = new URL(request.url()), path = url.pathname;
      if (request.method() !== 'GET') { requests.push(`write:${request.method()}`); return route.abort(); }
      let data: unknown;
      if (path === '/api/v1/datalink/studio-v2/workspace' && url.searchParams.get('entry') === 'locator') {
        const step = url.searchParams.get('step') === '2' ? 2 : 1;
        requests.push(`locator:${step}`);
        await (step === 1 ? deviceGate : ruleGate);
        const read = locatorRead(step); read.workspace = workspace; read.entry.devices = [device];
        data = { ...workspace, locator_entry: read.entry };
      } else if (path === '/api/v1/datalink/studio-v2/workspace') data = workspace;
      else if (path === '/api/v1/datalink/studio-v2/workspace/write-groups') data = { workspace_id: 'W', workspace_revision: 'W2', groups: [] };
      else if (path === '/api/v1/datalink/studio-v2/workspace/devices') data = [device];
      else if (path === '/api/v1/datalink/studio-v2/workspace/source-rules') {
        if (opened) { requests.push('ranges'); await listGate; }
        data = [nativeRule()];
      } else if (path === '/api/v1/datalink/studio-v2/workspace/mappings' || path === '/api/v1/datalink/points') data = [];
      else return route.abort();
      return route.fulfill({ json: { success: true, data } });
    });
    await page.goto('/tests/e2e/fixtures/recording-home-source-host.html?home');
    const trigger = page.getByTestId('home-device-points-D');
    await expect(trigger).toBeVisible();
    await page.evaluate(() => {
      const seen: string[] = [];
      new MutationObserver(() => {
        const dialog = document.querySelector('.gw-home-source-modal');
        if (dialog?.querySelector('[data-testid="studio-locator-notice"]')) seen.push('locator-notice');
        if (dialog?.querySelector('.gw-home-range-list')) seen.push('range-list');
      }).observe(document.body, { childList: true, subtree: true });
      (window as unknown as { transientSourceViews: string[] }).transientSourceViews = seen;
    });
    opened = true;
    await trigger.click();
    const dialog = page.getByRole('dialog'), status = dialog.getByRole('status');
    await expect.poll(() => requests.includes('locator:1')).toBe(true);
    await expect(status).toHaveText('正在讀取目前工作區…');
    const first = await status.boundingBox();
    await expect(dialog.getByTestId('studio-locator-notice')).toHaveCount(0);
    await page.screenshot({ path: info.outputPath(`device-loading-${viewport.width}.png`) });
    releaseDevice();
    await expect.poll(() => requests.includes('ranges')).toBe(true);
    await expect(status).toHaveText('正在讀取目前工作區…');
    const second = await status.boundingBox();
    await expect(dialog.getByRole('region', { name: '讀取範圍' })).toHaveCount(0);
    await page.screenshot({ path: info.outputPath(`list-loading-${viewport.width}.png`) });
    releaseList();
    await expect.poll(() => requests.includes('locator:2')).toBe(true);
    await expect(status).toHaveText('正在讀取目前工作區…');
    const third = await status.boundingBox();
    await expect(dialog.getByTestId('studio-locator-notice')).toHaveCount(0);
    await page.screenshot({ path: info.outputPath(`rule-loading-${viewport.width}.png`) });
    expect(first).not.toBeNull(); expect(second).not.toBeNull(); expect(third).not.toBeNull();
    expect(Math.abs(second!.y - first!.y)).toBeLessThanOrEqual(1);
    expect(Math.abs(third!.y - first!.y)).toBeLessThanOrEqual(1);
    releaseRule();
    await expect(dialog.getByTestId('rule-prefix-input')).toHaveValue('Rule R');
    expect(await page.evaluate(() => (window as unknown as { transientSourceViews: string[] }).transientSourceViews)).toEqual([]);
    await page.screenshot({ path: info.outputPath(`ready-${viewport.width}.png`) });
    expect(requests.filter(value => value.startsWith('write:'))).toEqual([]);
  });
}
