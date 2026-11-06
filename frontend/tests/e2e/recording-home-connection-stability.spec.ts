import { expect, test } from '@playwright/test';
import { locatorRead, nativeDevice, nativeRule } from '../unit/workbench-v2/studio-locator-fixtures';

for (const viewport of [{ width: 1440, height: 900 }, { width: 768, height: 1024 },
  { width: 390, height: 844 }, { width: 320, height: 568 }, { width: 900, height: 430 }, { width: 390, height: 844, light: true }]) {
  test(`home connection dialog loads without jumping at ${viewport.width}px ${'light' in viewport ? 'light' : 'dark'}`, async ({ page }, info) => {
    await page.setViewportSize({ width: viewport.width, height: viewport.height });
    let releaseDevice!: () => void;
    const deviceGate = new Promise<void>(resolve => { releaseDevice = resolve; });
    const requests: string[] = [];
    const workspace = { ...locatorRead().workspace, ordered_device_ids: ['D'] };
    const device = { ...nativeDevice(), connection_config: '{"host":"fixture.invalid","port":502,"slave_id":1,"timeout":1500}' };
    await page.route('**/api/**', async route => {
      const request = route.request(), url = new URL(request.url()), path = url.pathname;
      if (request.method() !== 'GET') { requests.push(`write:${request.method()}`); return route.abort(); }
      let data: unknown;
      if (path === '/api/v1/datalink/studio-v2/workspace' && url.searchParams.get('entry') === 'locator') {
        requests.push('locator'); await deviceGate;
        const read = locatorRead(); read.workspace = workspace; read.entry.devices = [device];
        data = { ...workspace, locator_entry: read.entry };
      } else if (path === '/api/v1/datalink/studio-v2/workspace') data = workspace;
      else if (path === '/api/v1/datalink/studio-v2/workspace/write-groups') data = { workspace_id: 'W', workspace_revision: 'W2', groups: [] };
      else if (path === '/api/v1/datalink/studio-v2/recording-home/sources') data = { workspace_id: 'W', state: 'available', source_digest: 'a'.repeat(64),
        devices: [device].map(({ id, name, protocol }) => ({ id, name, protocol, availability: 'unknown' })), candidates: [], excluded: [] };
      else if (path === '/api/v1/datalink/studio-v2/workspace/devices') data = [device];
      else if (path === '/api/v1/datalink/studio-v2/workspace/source-rules') data = [nativeRule()];
      else if (path === '/api/v1/datalink/studio-v2/workspace/mappings' || path === '/api/v1/datalink/points') data = [];
      else return route.abort();
      return route.fulfill({ json: { success: true, data } });
    });
    await page.goto(`/tests/e2e/fixtures/recording-home-source-host.html?home${'light' in viewport ? '&light' : ''}`);
    const trigger = page.locator('[data-home-focus="summary-connection-D"]');
    await expect(trigger).toBeVisible();
    await page.evaluate(() => {
      const seen: string[] = [];
      const values: string[] = [];
      new MutationObserver(() => {
        if (document.querySelector('.gw-home-source-modal [data-testid="studio-locator-notice"]')) seen.push('locator-notice');
      }).observe(document.body, { childList: true, subtree: true });
      (window as unknown as { transientConnectionViews: string[] }).transientConnectionViews = seen;
      (window as unknown as { connectionFirstValues: string[] }).connectionFirstValues = values;
      const sample = () => { const input = document.querySelector<HTMLInputElement>('[data-testid="recording-source-name"]');
        if (input?.getClientRects().length) values.push(JSON.stringify([input.value,
          document.querySelector<HTMLInputElement>('[data-testid="input-host"]')?.value,
          document.querySelector<HTMLInputElement>('[data-testid="input-port"]')?.value])); requestAnimationFrame(sample); };
      requestAnimationFrame(sample);
    });
    await trigger.click();
    const dialog = page.getByRole('dialog');
    await expect.poll(() => requests.includes('locator')).toBe(true);
    await expect(dialog.getByRole('status')).toBeVisible();
    const loading = await dialog.boundingBox();
    await page.screenshot({ path: info.outputPath(`connection-loading-${viewport.width}.png`) });
    await expect(dialog.getByRole('status')).toHaveText('正在讀取目前工作區…');
    // The loading state uses the same opaque settings frame as the ready form.
    await expect(dialog).not.toHaveCSS('background-color', 'rgba(0, 0, 0, 0)');
    await expect(dialog).not.toHaveCSS('border-top-color', 'rgba(0, 0, 0, 0)');
    await expect(dialog.locator('.gw-modal-header')).toHaveCSS('opacity', '1');
    await expect(dialog.getByRole('button', { name: '關閉', exact: true })).toHaveCount(1);
    await expect(dialog.getByRole('button', { name: '關閉', exact: true })).toBeFocused();
    await expect(dialog.getByTestId('studio-locator-notice')).toHaveCount(0);
    releaseDevice();
    const name = dialog.getByTestId('recording-source-name');
    await expect(name).toHaveValue('Device D');
    await page.evaluate(() => new Promise<void>(resolve => requestAnimationFrame(() => requestAnimationFrame(() => resolve()))));
    const values = await page.evaluate(() => (window as unknown as { connectionFirstValues: string[] }).connectionFirstValues);
    expect(values.length).toBeGreaterThan(0); expect(values.every(value => value === '["Device D","fixture.invalid","502"]')).toBe(true);
    await expect.poll(() => dialog.evaluate(element => element.contains(document.activeElement))).toBe(true);
    const ready = await dialog.boundingBox();
    await page.screenshot({ path: info.outputPath(`connection-ready-${viewport.width}.png`) });
    expect(loading).not.toBeNull(); expect(ready).not.toBeNull();
    expect(Math.abs(ready!.x - loading!.x)).toBeLessThanOrEqual(1);
    expect(Math.abs(ready!.y - loading!.y)).toBeLessThanOrEqual(1);
    expect(Math.abs(ready!.width - loading!.width)).toBeLessThanOrEqual(1);
    expect(Math.abs(ready!.height - loading!.height)).toBeLessThanOrEqual(1);
    expect(await page.evaluate(() => (window as unknown as { transientConnectionViews: string[] }).transientConnectionViews)).toEqual([]);
    await expect(dialog.getByTestId('recording-source-close')).toBeInViewport({ ratio: 1 });
    await name.fill('Edited device');
    await dialog.getByTestId('recording-source-close').click();
    await expect(dialog.getByTestId('recording-source-discard')).toBeVisible();
    await dialog.getByTestId('recording-source-discard').click();
    await expect(dialog).toHaveCount(0);
    await expect(trigger).toBeFocused();
    expect(requests.filter(value => value.startsWith('write:'))).toEqual([]);
  });
}
