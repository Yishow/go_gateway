import { execFileSync } from 'node:child_process';
import { createHash } from 'node:crypto';
import { expect, test, type Page, type TestInfo } from '@playwright/test';
import { assertContrast, box, contrastAudit, noHorizontalOverflow, phoneSizes, tracks } from './helpers/recording-memory-visual-audit';
import { longName, longNumber, unknownString } from './fixtures/recording-memory-visual-data';

// Vite isolated production component + deterministic mock inputs. NOT embedded delivery,
// production modal focus/history, gateway subscription, hardware or SQL acceptance.
const layouts = ['grid', 'navigator', 'map', 'chapters', 'detail'] as const;
type Layout = typeof layouts[number];
const references: Record<Layout, { id: string; sha256: string }> = {
  grid: { id: '51e970e9-3551-4cc5-9bf1-465239b04d0d', sha256: '625abc8a3cbbfaefd8197f3d46c79efb34b489de0d072c76ff1310c7aa772457' },
  navigator: { id: '86b48a15-0a46-4398-83f2-6b5e8d83da94', sha256: '9cdc5a9124bafdbd095037bccb1f86f5f21d8a04e484f8da5d3e65ad4a60596a' },
  map: { id: 'c91df0e8-3161-457e-9e8c-7c6146ddbf4b', sha256: '8bbf83722e31c1ef1d56c1ac448beebd9e4bdc81bfe673327f99d57d576f2baa' },
  chapters: { id: '943be82c-cc00-40cf-8873-bf81d82405f2', sha256: 'ac2c9e104cf6191c5b20c959c14aad75bd519365cdc45017457aa9d33774d47d' },
  detail: { id: '201f2c58-3543-4c36-b507-c703f07757d4', sha256: '5bce783eaaa1de2be897000a0a2a5b6a7030d253e03dfe537418cbe2d407482f' },
};
async function open(page: Page, scenario = 'reference', light = true) {
  // Abort any accidental API transport. The fixture itself needs no backend.
  await page.route('**/api/**', route => route.abort());
  await page.goto(`/tests/e2e/fixtures/recording-home-memory-host.html?scenario=${scenario}&en${light ? '&light' : ''}`);
  await expect(page.getByTestId('fixture-sample-notice')).toContainText('MOCK SAMPLE');
}
async function choose(page: Page, layout: Layout) {
  await page.getByTestId(`home-memory-layout-${layout}`).click();
  await expect(page.getByTestId('home-memory-view')).toHaveAttribute('data-layout', layout);
  await expect(page.locator('.gw-memory-grid')).toHaveCount(1);
}
async function locate(page: Page, address: string) {
  await page.getByRole('textbox').fill(address);
  await page.getByRole('button', { name: 'Go', exact: true }).click();
  await expect(page.getByTestId(`home-memory-cell-${address}`)).toHaveAttribute('tabindex', '0');
}
async function evidence(page: Page, info: TestInfo, name: string, data: object) {
  const root = execFileSync('git', ['rev-parse', '--show-toplevel'], { encoding: 'utf8' }).trim();
  const head = execFileSync('git', ['rev-parse', 'HEAD'], { encoding: 'utf8' }).trim();
  const diff = execFileSync('git', ['-C', root, 'diff', 'HEAD', '--', 'frontend/src', 'frontend/tests/e2e'], { encoding: 'utf8' });
  const status = execFileSync('git', ['-C', root, 'status', '--porcelain', '--', 'frontend/src', 'frontend/tests/e2e'], { encoding: 'utf8' });
  const activeLayout = await page.getByTestId('home-memory-view').getAttribute('data-layout') as Layout;
  const record = { reference: references[activeLayout], layout: activeLayout, productSha: head, dirty: Boolean(status.trim()), diffSha256: createHash('sha256').update(diff).digest('hex'),
    viewport: page.viewportSize(), fixture: 'Vite isolated recording-home-memory-host; deterministic SAMPLE values',
    visualReview: 'pending human reference comparison', intentionalDifference: 'Quiet mineral palette; existing safety controls retained; original sample values are not authority', ...data };
  await info.attach(`${name}.json`, { body: JSON.stringify(record, null, 2), contentType: 'application/json' });
  await info.attach(`${name}.png`, { body: await page.screenshot({ fullPage: true }), contentType: 'image/png' });
  expect(status.trim(), 'Final visual acceptance requires committed matching source and test inputs').toBe('');
}
async function structure(page: Page, layout: Layout, width: number) {
  const grid = await box(page.locator('.gw-memory-grid')), workspace = await box(page.locator('.gw-memory-workspace'));
  if (layout === 'grid' || layout === 'map' || layout === 'chapters') {
    expect(Math.abs(grid.width - workspace.width)).toBeLessThan(2);
    expect(Math.abs(grid.x - workspace.x)).toBeLessThan(2);
  }
  if (layout === 'navigator') {
    const nav = await box(page.getByTestId('home-memory-range-nav'));
    if (width >= 640) { expect(nav.x + nav.width).toBeLessThanOrEqual(grid.x); expect(Math.abs(nav.y - grid.y)).toBeLessThan(2); }
    else { expect(nav.y + nav.height).toBeLessThanOrEqual(grid.y); await expect(page.getByTestId('home-memory-range-nav').locator('details')).not.toHaveAttribute('open', ''); }
  }
  if (layout === 'chapters') {
    const band = await box(page.getByTestId('home-memory-chapters'));
    expect(band.y + band.height).toBeLessThanOrEqual(grid.y);
    expect(Math.abs(band.width - grid.width)).toBeLessThan(2);
    await expect(page.getByTestId('home-memory-chapter-F32').locator('summary')).toContainText('8 values');
    await expect(page.getByTestId('home-memory-chapter-F32').locator('summary')).toContainText('16 registers');
  }
  if (layout === 'map') {
    const strip = await box(page.getByTestId('home-memory-map')), summary = await box(page.getByTestId('home-memory-summary'));
    expect(strip.y + strip.height).toBeLessThanOrEqual(grid.y);
    expect(summary.y).toBeGreaterThanOrEqual(grid.y + grid.height);
    await expect(page.getByTestId('home-memory-detail')).toHaveCount(0);
    expect(await page.getByTestId('home-memory-summary').evaluate(node => getComputedStyle(node).display)).toBe('flex');
    // Compact wrapping facts stay proportional to their content, not a second full inspector.
    expect(summary.height).toBeLessThan(width >= 640 ? 240 : 400);
    await expect(page.locator('.gw-memory-map-strip > span')).toHaveCount(await page.locator('.gw-memory-grid > [data-offset]').count());
  }
  if (layout === 'detail') {
    const detail = await box(page.getByTestId('home-memory-detail'));
    if (width >= 640) { expect(detail.x).toBeGreaterThanOrEqual(grid.x + grid.width); expect(Math.abs(detail.y - grid.y)).toBeLessThan(2); }
    else expect(detail.y).toBeGreaterThanOrEqual(grid.y + grid.height);
    await expect(page.getByTestId('home-memory-detail').locator('[data-detail-group]')).toHaveCount(3);
  }
}

for (const width of [1440, 768, 390]) for (const light of [true, false]) for (const layout of layouts) {
  test(`reference ${layout} ${light ? 'light' : 'dark'} ${width}`, async ({ page }, info) => {
    await page.setViewportSize({ width, height: width === 390 ? 844 : 1000 });
    await open(page, 'reference', light); await choose(page, layout);
    await page.getByTestId('home-memory-view').evaluate(node => { node.scrollTop = 0; });
    await structure(page, layout, width); await noHorizontalOverflow(page);
    if (width >= 640) {
      expect(await page.getByTestId('home-memory-view').evaluate(node => {
        const viewport = node.getBoundingClientRect(), grid = node.querySelector('.gw-memory-grid')!;
        const count = getComputedStyle(grid).gridTemplateColumns.split(' ').filter(Boolean).length;
        return [...grid.children].slice(0, count * 2).every(cell => { const rect = cell.getBoundingClientRect(); return rect.top >= viewport.top && rect.bottom <= viewport.bottom; });
      }), 'Two complete initial rows remain visible').toBe(true);
    }
    const audit = await contrastAudit(page.getByRole('dialog'));
    const sizes = width === 390 ? await phoneSizes(page) : undefined;
    await evidence(page, info, `memory-${layout}-${light ? 'light' : 'dark'}-${width}`, { reference: references[layout], layout, theme: light ? 'light' : 'dark', audit, sizes });
    assertContrast(audit);
    await page.getByTestId('home-memory-cell-40001').focus(); await page.keyboard.press('ArrowRight');
    await expect(page.getByTestId('home-memory-cell-40002')).toBeFocused();
    const focus = await contrastAudit(page.getByRole('dialog'), true);
    await info.attach('focus-contrast.json', { body: JSON.stringify(focus), contentType: 'application/json' });
    assertContrast(focus);
  });
}

for (const light of [true, false]) test(`long content, bounded full lists and keyboard disclosure ${light ? 'light' : 'dark'}`, async ({ page }, info) => {
  await page.setViewportSize({ width: 390, height: 844 }); await open(page, 'long', light);
  const first = page.getByTestId('home-memory-cell-40001');
  await first.focus();
  const full = first.locator('.gw-memory-focus-information');
  await expect(full).toBeVisible(); await expect(full).toContainText(longNumber); await expect(full).toContainText('float64');
  expect(await full.evaluate(node => node.scrollWidth <= node.clientWidth + 1)).toBe(true);
  await locate(page, '40001');
  await choose(page, 'navigator');
  await page.getByTestId('home-memory-range-nav').locator('summary').click();
  await expect(page.getByTestId('home-memory-range-jump-LONG')).toContainText(longName);
  const nav = page.getByTestId('home-memory-range-nav').locator('ul');
  expect(await nav.evaluate(node => node.scrollHeight > node.clientHeight)).toBe(true);
  await page.getByTestId('home-memory-range-jump-LIST31').focus(); await page.keyboard.press('Enter');
  await expect(page.getByTestId('home-memory-cell-40225')).toHaveAttribute('tabindex', '0');
  await choose(page, 'chapters');
  const chapter = page.getByTestId('home-memory-chapter-LIST31');
  expect(await page.locator('.gw-memory-chapter-list').evaluate(node => node.scrollHeight > node.clientHeight)).toBe(true);
  await chapter.locator('summary').focus(); await page.keyboard.press('Enter');
  await expect(chapter).toHaveAttribute('open', '');
  await chapter.getByTestId('home-memory-range-jump-LIST31').focus(); await page.keyboard.press('Enter');
  await expect(page.getByTestId('home-memory-cell-40225')).toHaveAttribute('tabindex', '0');
  await noHorizontalOverflow(page); await evidence(page, info, 'long-lists', { light, scenario: 'long' });
});

test('full 256-slot strip and exact keyboard slider endpoints/interior', async ({ page }, info) => {
  await open(page, 'dense256'); await choose(page, 'map');
  await expect(page.locator('.gw-memory-map-strip > span')).toHaveCount(256);
  const slider = page.getByTestId('home-memory-map-select');
  await slider.focus(); await page.keyboard.press('Home'); await expect(slider).toHaveAttribute('aria-valuetext', 'W0');
  for (let index = 0; index < 127; index++) await page.keyboard.press('ArrowRight');
  await expect(slider).toHaveAttribute('aria-valuetext', 'W7F');
  await expect(page.getByTestId('home-memory-summary').locator('.gw-memory-detail-address')).toHaveText('W7F');
  await page.keyboard.press('End'); await expect(slider).toHaveAttribute('aria-valuetext', 'WFF');
  await expect(page.locator('.gw-memory-map-ends')).toHaveText('W0WFF');
  await evidence(page, info, 'dense256', { scenario: 'dense256' });
});

test('four-register row/page continuation keeps every physical focus and original reading key', async ({ page }, info) => {
  await page.setViewportSize({ width: 390, height: 844 }); await open(page, 'crossing');
  // Choose a viewport whose measured tracks force 40008..40011 to cross a row.
  const count = await tracks(page); expect(Math.floor(7 / count)).not.toBe(Math.floor(10 / count));
  for (const address of [40008, 40009, 40010, 40011]) {
    await page.getByTestId(`home-memory-cell-${address}`).focus();
    await expect(page.getByTestId(`home-memory-cell-${address}`)).toBeFocused();
    await expect(page.getByTestId(`home-memory-cell-${address}`)).toHaveAttribute('aria-label', /float64/);
  }
  const boundary = 40001 + (Math.floor(7 / count) + 1) * count;
  await expect(page.getByTestId(`home-memory-cell-${boundary}`)).toContainText('cont. 40008');
  await expect(page.getByTestId('home-memory-cell-40250')).toContainText('next row');
  await page.getByRole('group', { name: 'Address blocks' }).getByRole('button', { name: 'Next', exact: true }).click();
  await expect(page.getByTestId('home-memory-cell-40251')).toContainText('cont. 40249');
  await page.getByTestId('home-memory-cell-40251').focus();
  await expect(page.getByTestId('home-memory-cell-40251').locator('.gw-memory-focus-information')).toHaveAttribute('data-value-key', 'HR:249:float64');
  await expect.poll(() => page.getByTestId('fixture-visible-keys').textContent()).toContain('HR:249:float64');
  await choose(page, 'map'); await page.getByTestId('home-memory-cell-40251').focus();
  await expect(page.getByTestId('home-memory-summary')).toContainText('40249');
  await expect(page.getByTestId('home-memory-summary').locator('[data-value-key]')).toHaveAttribute('data-value-key', 'HR:249:float64');
  await expect.poll(() => page.getByTestId('fixture-visible-keys').textContent()).toContain('HR:249:float64');
  await evidence(page, info, 'cross-page', { scenario: 'crossing', columns: count });
});

for (const scenario of ['conflict', 'unknownLong']) test(`${scenario} keeps safety ahead of decorative typing`, async ({ page }, info) => {
  await page.setViewportSize({ width: 390, height: 844 }); await open(page, scenario); await choose(page, 'map');
  if (scenario === 'conflict') {
    await page.getByTestId('home-memory-cell-40003').click();
    await expect(page.getByTestId('home-memory-owners')).toContainText('RIGHT_');
    await expect(page.getByTestId('fixture-last-action')).toHaveText('');
  } else {
    // A saved start proves ownership, not string width or a usable reading.
    const start = page.getByTestId('home-memory-cell-40001');
    await expect(start).toHaveAttribute('data-state', 'used');
    await expect(start).toHaveAttribute('aria-label', /span after this address cannot be confirmed/);
    await expect(page.getByRole('dialog')).not.toContainText(unknownString);
    await expect(page.getByRole('dialog').locator('[data-value-key="HR:1:string"]')).toHaveCount(0);
    await start.click();
    await expect(page.getByTestId('fixture-last-action')).toHaveText('open:STRING');
    const remainder = page.getByTestId('home-memory-cell-40002');
    await expect(remainder).toHaveAttribute('data-state', 'unconfirmed');
    await remainder.click();
    await expect(page.getByTestId('fixture-last-action')).toHaveText('open:STRING');
    await expect(page.getByTestId('home-memory-summary').getByTestId('home-memory-detail-action')).toHaveCount(0);
  }
  await noHorizontalOverflow(page); await evidence(page, info, scenario, { scenario });
});

test('low-height safe-area sizing keeps all controls reachable', async ({ page }, info) => {
  await page.setViewportSize({ width: 768, height: 420 }); await open(page); await choose(page, 'chapters');
  const dialog = await box(page.getByRole('dialog'));
  expect(dialog.width).toBeCloseTo(752, 0); expect(dialog.height).toBeCloseTo(404, 0);
  for (const id of ['home-memory-cell-40001', 'home-memory-cell-40201', 'home-memory-close', 'home-memory-back']) {
    await page.getByTestId(id).scrollIntoViewIfNeeded(); await expect(page.getByTestId(id)).toBeInViewport({ ratio: 1 });
  }
  await noHorizontalOverflow(page); await evidence(page, info, 'low-height', { scenario: 'reference' });
});

test.describe('genuine touch context', () => {
  test.use({ hasTouch: true, viewport: { width: 390, height: 844 } });
  test('disclosure, exact slider and layout selection respond to taps', async ({ page }, info) => {
    await open(page, 'dense256');
    await page.getByTestId('home-memory-layout-navigator').tap();
    const disclosure = page.getByTestId('home-memory-range-nav').locator('details');
    await disclosure.locator('summary').tap(); await expect(disclosure).toHaveAttribute('open', '');
    await page.getByTestId('home-memory-range-jump-W256').tap();
    await page.getByTestId('home-memory-layout-map').tap();
    const slider = page.getByTestId('home-memory-map-select'); await slider.scrollIntoViewIfNeeded();
    const bounds = await box(slider);
    await page.touchscreen.tap(bounds.x + 1, bounds.y + bounds.height / 2);
    await expect(slider).toHaveAttribute('aria-valuetext', 'W0');
    await page.touchscreen.tap(bounds.x + bounds.width - 1, bounds.y + bounds.height / 2);
    await expect(slider).toHaveAttribute('aria-valuetext', 'WFF');
    await expect(page.getByTestId('fixture-last-action')).toHaveText('');
    await evidence(page, info, 'touch-slider', { scenario: 'dense256', hasTouch: true });
  });
});


test('navigator selection is separate from window membership and resets between equal-offset areas', async ({ page }, info) => {
  await open(page); await choose(page, 'navigator');
  await page.getByTestId('home-memory-cell-40001').focus();
  await expect(page.getByTestId('home-memory-range-jump-I16')).toHaveAttribute('aria-current', 'true');
  await expect(page.getByTestId('home-memory-range-jump-F32')).not.toHaveAttribute('aria-current', 'true');
  await expect(page.getByTestId('home-memory-range-jump-F32').locator('..')).toHaveAttribute('data-in-window', 'true');
  await page.getByTestId('home-memory-area-CS').click();
  await expect(page.getByTestId('home-memory-range-jump-BOOL')).not.toHaveAttribute('aria-current', 'true');
  await expect(page.getByTestId('home-memory-range-jump-BOOL').locator('..')).toHaveAttribute('data-in-window', 'true');
  await evidence(page, info, 'selection-membership', { scenario: 'reference' });
});
