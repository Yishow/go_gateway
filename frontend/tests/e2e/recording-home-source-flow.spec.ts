import { expect, test, type Locator } from '@playwright/test';
import { locatorRead, nativeDevice, nativeRule } from '../unit/workbench-v2/studio-locator-fixtures';

async function expectActionsAtBottomRight(dialog: Locator, first: Locator, last: Locator) {
  const footer = dialog.locator('.gw-modal-footer');
  await expect(footer).toContainText(await first.innerText());
  const [frame, a, b] = await Promise.all([footer.boundingBox(), first.boundingBox(), last.boundingBox()]);
  expect(frame).not.toBeNull(); expect(a).not.toBeNull(); expect(b).not.toBeNull();
  expect(Math.abs(a!.y - b!.y)).toBeLessThanOrEqual(1);
  expect(b!.x - (a!.x + a!.width)).toBeGreaterThanOrEqual(0);
  expect(b!.x - (a!.x + a!.width)).toBeLessThanOrEqual(24);
  expect(frame!.x + frame!.width - (b!.x + b!.width)).toBeLessThanOrEqual(32);
}

for (const viewport of [{ width: 1440, height: 900 }, { width: 320, height: 568 }, { width: 390, height: 844, light: true }]) {
  test(`home keeps one loading view until its saved source editor is ready at ${viewport.width}px ${'light' in viewport ? 'light' : 'dark'}`, async ({ page }, info) => {
    await page.setViewportSize({ width: viewport.width, height: viewport.height });
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
    await page.goto(`/tests/e2e/fixtures/recording-home-source-host.html?home&strict${'light' in viewport ? '&light' : ''}`);
    const trigger = page.getByTestId('home-device-points-D');
    await expect(trigger).toBeVisible();
    await page.evaluate(() => {
      const seen: string[] = [];
      const values: string[] = [];
      const frames: string[] = [];
      new MutationObserver(() => {
        const dialog = document.querySelector('.gw-home-source-modal');
        if (dialog?.querySelector('[data-testid="studio-locator-notice"]')) seen.push('locator-notice');
        if (dialog?.querySelector('.gw-home-range-list')) seen.push('range-list');
        if (dialog?.textContent?.includes('目前無法確認此設備或讀取範圍')) seen.push('unavailable');
      }).observe(document.body, { childList: true, subtree: true });
      (window as unknown as { transientSourceViews: string[] }).transientSourceViews = seen;
      (window as unknown as { sourceFirstValues: string[] }).sourceFirstValues = values;
      (window as unknown as { sourceFrames: string[] }).sourceFrames = frames;
      const sample = () => { const input = document.querySelector<HTMLInputElement>('[data-testid="rule-prefix-input"]');
        const frame = document.querySelector<HTMLElement>('.gw-home-source-modal'), header = frame?.querySelector<HTMLElement>('.gw-modal-header');
        if (frame) { const box = frame.getBoundingClientRect(), style = getComputedStyle(frame);
          frames.push(JSON.stringify([Math.round(box.x), Math.round(box.y), Math.round(box.width), Math.round(box.height), style.backgroundColor, style.borderTopColor, header ? getComputedStyle(header).opacity : 'none'])); }
        if (input?.getClientRects().length) values.push(JSON.stringify([input.value,
          document.querySelector<HTMLInputElement>('[data-testid="rule-start-input"]')?.value,
          document.querySelector<HTMLInputElement>('[data-testid="rule-count-input"]')?.value])); requestAnimationFrame(sample); };
      requestAnimationFrame(sample);
    });
    opened = true;
    await trigger.click();
    const dialog = page.getByRole('dialog'), status = dialog.getByRole('status');
    await expect.poll(() => requests.includes('locator:1')).toBe(true);
    await expect(status).toHaveText('正在讀取目前工作區…');
    // One opaque settings frame from the first frame on: no transparent cue that later grows into the dialog.
    await expect(dialog).not.toHaveCSS('background-color', 'rgba(0, 0, 0, 0)');
    await expect(dialog).not.toHaveCSS('border-top-color', 'rgba(0, 0, 0, 0)');
    await expect(dialog.locator('.gw-modal-header')).toHaveCSS('opacity', '1');
    await expect(dialog.getByRole('heading', { name: /設定讀取點位/ })).toBeVisible();
    await expect(dialog.getByRole('button', { name: '關閉', exact: true })).toHaveCount(1);
    await expect(dialog.getByRole('button', { name: '關閉', exact: true })).toBeFocused();
    const first = await dialog.boundingBox();
    await expect(dialog.getByTestId('studio-locator-notice')).toHaveCount(0);
    await page.screenshot({ path: info.outputPath(`device-loading-${viewport.width}.png`) });
    releaseDevice();
    await expect.poll(() => requests.includes('ranges')).toBe(true);
    await expect(status).toHaveText('正在讀取目前工作區…');
    const second = await dialog.boundingBox();
    await expect(dialog.getByRole('region', { name: '讀取範圍' })).toHaveCount(0);
    await page.screenshot({ path: info.outputPath(`list-loading-${viewport.width}.png`) });
    releaseList();
    await expect.poll(() => requests.includes('locator:2')).toBe(true);
    await expect(status).toHaveText('正在讀取目前工作區…');
    const third = await dialog.boundingBox();
    await expect(dialog.getByTestId('studio-locator-notice')).toHaveCount(0);
    await page.screenshot({ path: info.outputPath(`rule-loading-${viewport.width}.png`) });
    expect(first).not.toBeNull(); expect(second).not.toBeNull(); expect(third).not.toBeNull();
    expect(Math.abs(second!.y - first!.y)).toBeLessThanOrEqual(1);
    expect(Math.abs(third!.y - first!.y)).toBeLessThanOrEqual(1);
    releaseRule();
    await expect(dialog.getByTestId('rule-prefix-input')).toHaveValue('Rule R');
    await page.evaluate(() => new Promise<void>(resolve => requestAnimationFrame(() => requestAnimationFrame(() => resolve()))));
    const values = await page.evaluate(() => (window as unknown as { sourceFirstValues: string[] }).sourceFirstValues);
    expect(values.length).toBeGreaterThan(0); expect(values.every(value => value === '["Rule R","40001","2"]')).toBe(true);
    const ready = await dialog.boundingBox();
    expect(Math.abs(ready!.y - first!.y)).toBeLessThanOrEqual(1); expect(Math.abs(ready!.height - first!.height)).toBeLessThanOrEqual(1);
    expect(Math.abs(ready!.width - first!.width)).toBeLessThanOrEqual(1);
    await expectActionsAtBottomRight(dialog, dialog.getByTestId('home-range-done'), dialog.getByTestId('home-range-save'));
    const frames = await page.evaluate(() => (window as unknown as { sourceFrames: string[] }).sourceFrames);
    expect(frames.length).toBeGreaterThan(0); expect([...new Set(frames)]).toEqual([frames[0]]);
    await expect.poll(() => dialog.evaluate(element => element.contains(document.activeElement))).toBe(true);
    expect(await page.evaluate(() => (window as unknown as { transientSourceViews: string[] }).transientSourceViews)).toEqual([]);
    await page.screenshot({ path: info.outputPath(`ready-${viewport.width}.png`) });
    // The address summary's function badge stays readable (WCAG AA 4.5:1) on the frame in both themes.
    const contrast = await dialog.getByTestId('function-code-badge').evaluate(badge => {
      // Let the browser convert any CSS color space (oklab/oklch) to sRGB and composite alpha on a 1px canvas.
      const context = document.createElement('canvas').getContext('2d', { willReadFrequently: true })!;
      const paint = (...colors: string[]) => { context.clearRect(0, 0, 1, 1); for (const color of colors) { context.fillStyle = color; context.fillRect(0, 0, 1, 1); } return [...context.getImageData(0, 0, 1, 1).data].slice(0, 3); };
      const frame = getComputedStyle(badge.closest('.gw-home-source-modal')!).backgroundColor, fill = getComputedStyle(badge).backgroundColor;
      const back = paint(frame, fill), text = paint(frame, fill, getComputedStyle(badge).color);
      const luminance = (rgb: number[]) => rgb.map(value => { const c = value / 255; return c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4; })
        .reduce((sum, value, index) => sum + value * [0.2126, 0.7152, 0.0722][index], 0);
      const [high, low] = [luminance(text), luminance(back)].sort((x, y) => y - x);
      return (high + 0.05) / (low + 0.05);
    });
    expect(contrast).toBeGreaterThanOrEqual(4.5);
    // The memory entry sits in the settings title bar; its back returns to this same setting.
    await expect(dialog.locator('.gw-modal-header').getByTestId('home-range-memory')).toBeVisible();
    await page.screenshot({ path: info.outputPath(`memory-entry-${viewport.width}.png`) });
    await dialog.getByTestId('home-range-memory').click();
    await expect(dialog.getByTestId('home-memory-view')).toBeVisible();
    await expect(dialog.getByTestId('home-memory-cell-40001')).toHaveAttribute('data-state', 'used');
    await expect(dialog.locator('.gw-modal-header').getByRole('button', { name: '關閉', exact: true })).toHaveCount(0);
    const memoryClose = dialog.locator('.gw-modal-footer').getByRole('button', { name: '關閉', exact: true });
    await expectActionsAtBottomRight(dialog, dialog.getByTestId('home-memory-back'), memoryClose);
    await dialog.getByTestId('home-memory-back').click();
    await expect(dialog.getByTestId('rule-prefix-input')).toHaveValue('Rule R');
    await dialog.getByTestId('home-range-memory').click();
    await expect(dialog.getByTestId('home-memory-view')).toBeVisible();
    await memoryClose.click();
    await expect(dialog).toBeHidden();
    expect(requests.filter(value => value.startsWith('write:'))).toEqual([]);
  });
}
