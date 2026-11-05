import { expect, test, type Locator, type Page } from '@playwright/test';

// Isolated, deterministic SAMPLE readings. This suite is not HSL or gateway integration evidence.
const layouts = ['grid', 'navigator', 'map', 'chapters', 'detail'] as const;
type Layout = typeof layouts[number];
async function open(page: Page, scenario = 'reference', light = true) {
  await page.goto(`/tests/e2e/fixtures/recording-home-memory-host.html?scenario=${scenario}&en${light ? '&light' : ''}`);
  await expect(page.getByTestId('fixture-sample-notice')).toContainText('MOCK SAMPLE');
  await expect(page.getByTestId('home-memory-values-status')).toHaveAttribute('data-status', 'live');
}
async function switchTo(page: Page, layout: Layout) {
  const button = page.getByTestId(`home-memory-layout-${layout}`);
  await button.click();
  await expect(button).toHaveAttribute('aria-pressed', 'true');
  await expect(page.getByTestId('home-memory-view')).toHaveAttribute('data-layout', layout);
}
async function locate(page: Page, address: string) {
  await page.getByRole('textbox').fill(address);
  await page.getByRole('button', { name: 'Go', exact: true }).click();
  await expect(page.getByTestId(`home-memory-cell-${address}`)).toHaveAttribute('tabindex', '0');
}
async function noOverflow(dialog: Locator) {
  // Clip descendants by each scroll container before comparing with the modal bounds.
  const overflow = await dialog.evaluate(root => {
    const bounds = root.getBoundingClientRect();
    return [...root.querySelectorAll<HTMLElement>('*')].filter(node => {
      const box = node.getBoundingClientRect();
      let { left, right } = box;
      for (let parent = node.parentElement; parent && parent !== root; parent = parent.parentElement) {
        if (getComputedStyle(parent).overflowX === 'visible') continue;
        const clip = parent.getBoundingClientRect();
        left = Math.max(left, clip.left); right = Math.min(right, clip.right);
      }
      return box.width > 0 && box.height > 0 && right > left && (left < bounds.left - 1 || right > bounds.right + 1);
    }).map(node => `${node.tagName}.${node.className}`);
  });
  expect(overflow).toEqual([]);
}
async function joinedGeometry(page: Page, start: number, width: number) {
  const first = page.getByTestId(`home-memory-cell-${start}`);
  await first.scrollIntoViewIfNeeded();
  const grid = page.locator('.gw-memory-grid');
  const columns = await grid.evaluate(node => getComputedStyle(node).gridTemplateColumns.split(' ').filter(Boolean).length);
  const offset = await first.getAttribute('data-offset');
  const firstOffset = await grid.locator('[data-offset]').first().getAttribute('data-offset');
  const column = (Number(offset) - Number(firstOffset)) % columns;
  const span = Math.min(width, columns - column);
  const shared = first.locator('.gw-memory-shared');
  await expect(shared).toHaveAttribute('data-shared-span', String(span));
  const end = page.getByTestId(`home-memory-cell-${start + span - 1}`);
  const [labelBox, firstBox, endBox] = await Promise.all([shared.boundingBox(), first.boundingBox(), end.boundingBox()]);
  expect(labelBox).not.toBeNull(); expect(firstBox).not.toBeNull(); expect(endBox).not.toBeNull();
  expect(Math.abs(firstBox!.y - endBox!.y)).toBeLessThan(1);
  // Assert the actual rendered label spans physical slots, not just its data attribute.
  expect(labelBox!.width).toBeGreaterThan(firstBox!.width * span - 12);
  expect(labelBox!.x + labelBox!.width).toBeLessThanOrEqual(endBox!.x + endBox!.width + 1);
  for (let part = 1; part < span; part++) {
    await expect(page.getByTestId(`home-memory-cell-${start + part}`)).toHaveAttribute('data-join-before', 'true');
  }
  if (span < width) {
    await expect(end).toContainText('next row');
    const continuation = page.getByTestId(`home-memory-cell-${start + span}`);
    await expect(continuation).toContainText(`cont. ${start}`);
    await expect(continuation).not.toHaveAttribute('data-join-before', 'true');
  }
}

for (const width of [1440, 768, 390]) {
  for (const light of [true, false]) {
    test(`five layouts render typed SAMPLE readings at ${width}px ${light ? 'light' : 'dark'}`, async ({ page }, testInfo) => {
      await page.setViewportSize({ width, height: width === 390 ? 844 : 1000 });
      await open(page, 'reference', light);
      for (const layout of layouts) {
        await switchTo(page, layout);
        await expect(page.getByTestId('home-memory-cell-40001').locator('.gw-memory-reading')).toHaveText('0');
        await expect(page.getByTestId('home-memory-cell-40017').locator('.gw-memory-reading')).toHaveText('12.5');
        await expect(page.getByTestId('home-memory-cell-40032')).toHaveAttribute('data-state', 'used');
        await expect(page.getByTestId('home-memory-cell-40033')).toHaveAttribute('data-state', 'empty');
        await joinedGeometry(page, 40017, 2);
        await joinedGeometry(page, 40041, 4);
        await noOverflow(page.getByRole('dialog'));
        await page.getByTestId(`home-memory-layout-${layout}`).scrollIntoViewIfNeeded();
        await page.screenshot({ path: testInfo.outputPath(`memory-${layout}-${light ? 'light' : 'dark'}-${width}.png`), fullPage: true });
      }
    });
  }
}

test('switches preserve selected area, located window, selected address and reading', async ({ page }) => {
  await open(page);
  await page.getByTestId('home-memory-area-IR').click();
  await locate(page, '30201');
  for (const layout of [...layouts.slice(1), 'grid'] as Layout[]) {
    await switchTo(page, layout);
    await expect(page.getByTestId('home-memory-area-IR')).toHaveAttribute('aria-selected', 'true');
    await expect(page.getByRole('textbox')).toHaveValue('30201');
    await expect(page.getByTestId('home-memory-cell-30201')).toHaveAttribute('tabindex', '0');
    await expect(page.getByTestId('home-memory-cell-30201').locator('.gw-memory-reading')).toHaveText('65535');
    await expect(page.getByTestId('home-memory-cell-30001')).toHaveCount(0);
  }
});

test('navigator lists all ranges, chapters expand and jump with physical end addresses', async ({ page }) => {
  await open(page);
  await switchTo(page, 'navigator');
  const navigation = page.getByTestId('home-memory-range-nav');
  if (!(await navigation.locator('details').evaluate(node => node.hasAttribute('open')))) await navigation.locator('summary').click();
  await expect(navigation.getByTestId(/^home-memory-range-jump-/)).toHaveCount(6);
  await expect(navigation).toContainText('40032');
  for (const [id, address] of [['I16', '40001'], ['F32', '40017'], ['F64', '40041'], ['U16', '40045'], ['I32', '40047'], ['REMOTE', '40201']]) {
    await navigation.getByTestId(`home-memory-range-jump-${id}`).click();
    await expect(page.getByTestId(`home-memory-cell-${address}`)).toHaveAttribute('tabindex', '0');
  }
  await switchTo(page, 'chapters');
  const chapter = page.getByTestId('home-memory-chapter-F32');
  await expect(chapter).not.toHaveAttribute('open', '');
  await chapter.locator('summary').focus(); await page.keyboard.press('Enter');
  await expect(chapter).toHaveAttribute('open', '');
  await expect(chapter).toContainText('40017'); await expect(chapter).toContainText('40032');
  await chapter.getByTestId('home-memory-range-jump-F32').click();
  await expect(page.getByTestId('home-memory-cell-40017')).toHaveAttribute('tabindex', '0');
});

test('map selection and detail previous/next resolve continuation values and actions', async ({ page }) => {
  await open(page); await switchTo(page, 'map');
  const slider = page.getByTestId('home-memory-map-select');
  await slider.focus(); await page.keyboard.press('Home');
  await page.keyboard.press('ArrowRight');
  await expect(slider).toHaveAttribute('aria-valuetext', '40002');
  await expect(page.locator('.gw-memory-detail-address')).toHaveText('40002');
  await locate(page, '40017');
  await switchTo(page, 'detail');
  await expect(page.locator('.gw-memory-detail-reading')).toHaveText('12.5');
  await page.getByTestId('home-memory-detail-next').click();
  await expect(page.locator('.gw-memory-detail-address')).toHaveText('40018');
  await expect(page.locator('.gw-memory-detail-reading')).toHaveText('12.5');
  await page.getByTestId('home-memory-detail-previous').click();
  await expect(page.locator('.gw-memory-detail-address')).toHaveText('40017');
  await page.getByTestId('home-memory-detail-action').click();
  await expect(page.getByTestId('fixture-last-action')).toHaveText('open:F32');
  await locate(page, '40033');
  await page.getByTestId('home-memory-detail-action').click();
  await expect(page.getByTestId('fixture-last-action')).toHaveText('add:40033');
});

for (const layout of layouts) {
  test(`${layout} preserves used, empty, conflict and unconfirmed semantics`, async ({ page }) => {
    await open(page, 'conflict'); await switchTo(page, layout);
    await page.getByTestId('home-memory-cell-40001').click();
    await expect(page.getByTestId('fixture-last-action')).toHaveText('open:A');
    await page.getByTestId('home-memory-cell-40005').click();
    await expect(page.getByTestId('fixture-last-action')).toHaveText('add:40005');
    await page.getByTestId('home-memory-cell-40003').click();
    const owners = page.getByTestId('home-memory-owners');
    await expect(owners).toBeVisible();
    await owners.getByRole('button', { name: /RIGHT_/ }).click();
    await expect(page.getByTestId('fixture-last-action')).toHaveText('open:B');
    await open(page, 'unconfirmed'); await switchTo(page, layout);
    await locate(page, 'D100');
    const unknown = page.getByTestId('home-memory-cell-D101');
    await expect(unknown).toHaveAttribute('data-state', 'unconfirmed');
    expect(await unknown.evaluate(node => node.tagName)).not.toBe('BUTTON');
    await unknown.click();
    await expect(page.getByTestId('fixture-last-action')).toHaveText('');
    if (layout === 'detail') await expect(page.getByTestId('home-memory-detail-action')).toHaveCount(0);
  });
}

test('layout buttons, physical cells and area tabs remain keyboard operable', async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 }); await open(page);
  const detail = page.getByTestId('home-memory-layout-detail');
  await detail.focus(); await page.keyboard.press('Enter');
  await expect(detail).toHaveAttribute('aria-pressed', 'true');
  await page.getByTestId('home-memory-cell-40001').focus();
  await page.keyboard.press('ArrowRight');
  await expect(page.getByTestId('home-memory-cell-40002')).toBeFocused();
  const columns = await page.locator('.gw-memory-grid').evaluate(node => getComputedStyle(node).gridTemplateColumns.split(' ').filter(Boolean).length);
  await page.keyboard.press('ArrowDown');
  await expect(page.getByTestId(`home-memory-cell-${40002 + columns}`)).toBeFocused();
  await page.keyboard.press('ArrowUp'); await page.keyboard.press('Enter');
  await expect(page.getByTestId('fixture-last-action')).toHaveText('open:I16');
  await page.getByTestId('home-memory-area-HR').focus(); await page.keyboard.press('ArrowRight');
  await expect(page.getByRole('tab', { selected: true })).toBeFocused();
  await noOverflow(page.getByRole('dialog'));
});
