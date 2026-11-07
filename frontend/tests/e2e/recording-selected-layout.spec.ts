import { expect, test } from '@playwright/test';

for (const width of [320, 1280]) test(`keeps the name help glyph close to its label at ${width}px`, async ({ page }) => {
  await page.setViewportSize({ width, height: 720 });
  await page.route('**/api/**', route => route.abort());
  await page.goto('/tests/e2e/fixtures/recording-selected-host.html');
  const name = page.locator('.recording-simple-name');
  const label = await name.locator('label span').first().boundingBox();
  const help = name.locator('[data-parameter-help="record.name"]');
  const hit = await help.boundingBox();
  const glyphLeft = await help.evaluate(element => { const range = document.createRange(); range.selectNodeContents(element); return range.getBoundingClientRect().left; });
  expect.soft(glyphLeft - (label!.x + label!.width)).toBeLessThanOrEqual(12);
  expect(hit!.width).toBeGreaterThanOrEqual(44); expect(hit!.height).toBeGreaterThanOrEqual(44);
  if (width === 1280) {
    const table = page.getByTestId('group-member-table');
    const hint = await table.locator('.recording-compact-members__name-hint').boundingBox();
    const info = await table.locator('.recording-compact-members__help').boundingBox();
    expect(hint!.y + hint!.height).toBeGreaterThanOrEqual(info!.y);
  }
});

test('sizes desktop columns from their contents and gives the remaining space to point names', async ({ page }) => {
  await page.setViewportSize({ width: 1280, height: 720 });
  await page.route('**/api/**', route => route.abort());
  const widths = async () => page.locator('th').evaluateAll(headers => Object.fromEntries(headers.map(header => [header.dataset.column!, header.getBoundingClientRect().width])));
  await page.goto('/tests/e2e/fixtures/recording-selected-host.html');
  const short = await widths();
  await page.goto('/tests/e2e/fixtures/recording-selected-host.html?longdevice');
  const long = await widths();
  expect(long.device).toBeGreaterThan(short.device + 30);
  expect(long.pointOnly).toBeLessThan(short.pointOnly - 30);
  expect(short.pointOnly).toBeGreaterThan(Math.max(short.device, short.addressOnly, short.wordOrder, short.dataType, short.actions));
});

for (const width of [320, 390, 640, 641, 768, 1000, 1001, 1280, 1440]) {
  test(`keeps navigation and the selected-points table usable at ${width}px`, async ({ page }, testInfo) => {
    await page.setViewportSize({ width, height: 720 });
    await page.route('**/api/**', route => route.abort());
    await page.goto('/tests/e2e/fixtures/recording-selected-host.html');
    const heading = page.getByRole('heading', { name: '編輯記錄' });
    const nav = page.getByRole('navigation', { name: '本頁設定區段' });
    const table = page.getByTestId('group-member-table');
    const scroll = table.getByRole('region');
    await expect(heading).toBeVisible();
    await expect(table.getByTestId('group-member-search')).toHaveAttribute('placeholder', '名稱、設備或位址');
    await expect(page.getByTestId('group-member-scope')).toHaveClass(/recording-editor-point-count/);
    await expect(page.getByTestId('group-member-scope')).toContainText('記錄點位 · 32 點');
    await expect(page.getByTestId('group-member-scope')).not.toContainText('32 / 32');
    await expect(table.locator('[data-parameter-help="member.required"]')).toBeVisible();
    await expect(table.getByText('null', { exact: true })).toHaveCount(0);
    await expect(table.getByTestId('group-member-row-P0').getByRole('combobox', { name: '字序' })).toBeVisible();
    await expect(table.getByTestId('group-member-row-P1').getByText('不適用')).toBeVisible();
    await expect(table.getByTestId('group-member-row-P1').getByText('不適用')).toHaveAttribute('title', /int16.*16/);
    expect(await table.getByTestId('group-member-row-P0').getByRole('textbox').evaluate(element => parseFloat(getComputedStyle(element).borderWidth))).toBeGreaterThan(0);
    await expect(table.getByTestId('group-member-row-P0').locator('td[data-column="actions"] button').first()).toHaveAttribute('aria-pressed', 'true');
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
    const overflow = await scroll.evaluate(element => {
      const right = element.getBoundingClientRect().right;
      return { scrollWidth: element.scrollWidth, clientWidth: element.clientWidth, outside: [...element.querySelectorAll('*')].filter(child => child.getBoundingClientRect().right > right + 1 || child.scrollWidth > child.clientWidth + 1).slice(0, 8).map(child => ({ tag: child.tagName, text: child.textContent?.slice(0, 40), class: child.className, scroll: child.scrollWidth, client: child.clientWidth })) };
    });
    if (width <= 641) expect(overflow.scrollWidth, JSON.stringify(overflow)).toBeGreaterThan(overflow.clientWidth);
    else expect(overflow.scrollWidth, JSON.stringify(overflow)).toBeLessThanOrEqual(overflow.clientWidth);
    if (overflow.scrollWidth > overflow.clientWidth) await table.getByRole('columnheader', { name: '操作' }).scrollIntoViewIfNeeded();
    const regionBox = await scroll.boundingBox(), headerBox = await table.getByRole('columnheader', { name: '操作' }).boundingBox();
    expect(headerBox!.x).toBeGreaterThanOrEqual(regionBox!.x); expect(headerBox!.x + headerBox!.width).toBeLessThanOrEqual(regionBox!.x + regionBox!.width);
    if (overflow.scrollWidth > overflow.clientWidth) await scroll.evaluate(element => { element.scrollLeft = 0; });
    await page.screenshot({ path: testInfo.outputPath('top.png') });
    if (width === 1280) {
      await page.evaluate(() => window.scrollTo(0, 300));
      await page.screenshot({ path: testInfo.outputPath('scrolled.png') });
    }
    if (width <= 1000) {
      await page.evaluate(() => window.scrollTo(0, 300));
      const row = table.getByTestId('group-member-row-P0'), cells = row.locator('td');
      const deviceBox = await cells.nth(0).boundingBox(), pointBox = await cells.nth(1).boundingBox(), addressBox = await cells.nth(2).boundingBox();
      const orderBox = await cells.nth(3).boundingBox(), actionBox = await cells.nth(5).boundingBox();
      expect(deviceBox!.y).toBe(pointBox!.y); expect(addressBox!.y).toBe(pointBox!.y);
      expect(orderBox!.y).toBe(pointBox!.y);
      expect(actionBox!.y).toBe(orderBox!.y);
      const remove = row.getByRole('button', { name: /移除/ });
      const typeMetrics = await cells.nth(4).evaluate(element => ({ available: element.clientWidth, content: element.scrollWidth, font: getComputedStyle(element).fontSize, padding: getComputedStyle(element).padding }));
      expect(typeMetrics.content, JSON.stringify(typeMetrics)).toBeLessThanOrEqual(typeMetrics.available);
      const toggle = row.locator('button[aria-pressed]'), removeBox = await remove.boundingBox(), toggleBox = await toggle.boundingBox();
      expect(toggleBox!.y).toBe(removeBox!.y); expect(toggleBox!.x + toggleBox!.width).toBeLessThanOrEqual(actionBox!.x + actionBox!.width);
      await expect(toggle).toHaveAttribute('aria-pressed', 'true');
      const enabledColor = await toggle.evaluate(element => getComputedStyle(element).backgroundColor);
      await toggle.click(); await expect(toggle).toHaveAttribute('aria-pressed', 'false');
      expect(await toggle.evaluate(element => getComputedStyle(element).backgroundColor)).not.toBe(enabledColor);
      expect(await scroll.evaluate(element => element.scrollWidth > element.clientWidth)).toBe(width <= 641);
      if (width <= 390) {
        const order = row.getByRole('combobox', { name: '字序' }); await order.selectOption('CDAB');
        const selectBox = await order.boundingBox(), saveBox = await row.getByRole('button', { name: '儲存字序' }).boundingBox();
        await page.screenshot({ path: testInfo.outputPath('editing.png') });
        expect(selectBox!.width).toBeGreaterThanOrEqual(60);
        expect(saveBox!.x + saveBox!.width).toBeLessThanOrEqual(orderBox!.x + orderBox!.width);
      }
      await page.screenshot({ path: testInfo.outputPath('scrolled.png') });
      await remove.click(); await expect(page.getByTestId('group-member-scope')).toContainText('31 點');
    }
    await page.evaluate(() => window.scrollTo(0, 700));
    await expect.poll(async () => (await heading.boundingBox())?.y ?? -999).toBeGreaterThanOrEqual(0);
    await expect.poll(async () => (await nav.boundingBox())?.y ?? -999).toBeGreaterThanOrEqual(0);
    const headingBox = await heading.boundingBox(), navBox = await nav.boundingBox();
    expect(navBox!.y).toBeGreaterThanOrEqual(headingBox!.y + headingBox!.height - 1);
    if (width > 1350) {
      const summary = page.getByRole('complementary', { name: '這組記錄' });
      const summaryBox = await summary.boundingBox();
      expect(summaryBox!.y).toBeGreaterThanOrEqual(navBox!.y + navBox!.height - 1);
    }
    await nav.getByRole('button', { name: /2 記錄方式/ }).click();
    const section = page.getByTestId('recording-editor-method-note');
    const currentNavBox = await nav.boundingBox(), sectionBox = await section.boundingBox();
    expect(sectionBox!.y).toBeGreaterThanOrEqual(currentNavBox!.y + currentNavBox!.height - 1);
  });
}

test('keeps a long device name and editable point usable on a narrow screen', async ({ page }) => {
  await page.setViewportSize({ width: 320, height: 720 });
  await page.route('**/api/**', route => route.abort());
  await page.goto('/tests/e2e/fixtures/recording-selected-host.html?longdevice');
  const row = page.getByTestId('group-member-row-P0');
  await expect(row.locator('td[data-column="device"]')).toHaveAttribute('title', 'MC Building North');
  const point = row.getByRole('textbox');
  expect((await point.boundingBox())!.width).toBeGreaterThanOrEqual(60);
  await point.focus(); await expect(point).toBeEditable();
  await point.fill('新名稱'); await point.press('Enter'); await expect(point).toHaveValue('新名稱');
  expect(await page.getByTestId('group-member-table').getByRole('region').evaluate(element => element.scrollWidth > element.clientWidth)).toBe(true);
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
});

test('keeps English per-point controls and actions usable together at 320px', async ({ page }) => {
  await page.setViewportSize({ width: 320, height: 720 });
  await page.route('**/api/**', route => route.abort());
  await page.goto('/tests/e2e/fixtures/recording-selected-host.html?lang=en&mixed');
  const row = page.getByTestId('group-member-row-P0'), order = row.getByRole('combobox', { name: 'Word order' });
  await order.selectOption('CDAB');
  const orderBox = await row.locator('td[data-column="wordOrder"]').boundingBox(), save = await row.getByRole('button', { name: 'Save word order' }).boundingBox();
  const metrics = await row.evaluate(element => ({ tracks: getComputedStyle(element.closest('table')!).gridTemplateColumns, cells: [...element.querySelectorAll('td')].map(cell => ({ name: cell.getAttribute('data-column') ?? cell.getAttribute('data-testid'), width: cell.getBoundingClientRect().width })) }));
  expect((await order.boundingBox())!.width, JSON.stringify(metrics)).toBeGreaterThanOrEqual(60);
  expect(save!.x + save!.width).toBeLessThanOrEqual(orderBox!.x + orderBox!.width);
  for (const button of await row.locator('td[data-column="actions"] button').all()) expect(await button.evaluate(element => element.scrollWidth <= element.clientWidth)).toBe(true);
  expect(await page.getByTestId('group-member-table').getByRole('region').evaluate(element => element.scrollWidth > element.clientWidth)).toBe(true);
});

for (const width of [320, 390, 1280, 1440]) {
  test(`keeps the English action buttons usable at ${width}px`, async ({ page }) => {
    await page.setViewportSize({ width, height: 720 });
    await page.route('**/api/**', route => route.abort());
    await page.goto('/tests/e2e/fixtures/recording-selected-host.html?lang=en');
    const scroll = page.getByTestId('group-member-table').getByRole('region');
    const row = page.getByTestId('group-member-row-P0'), cell = row.locator('td').last();
    const remove = cell.getByRole('button', { name: /Remove/ }), toggle = cell.locator('button[aria-pressed]');
    const overflow = await scroll.evaluate(element => {
      const right = element.getBoundingClientRect().right;
      return { scrollWidth: element.scrollWidth, clientWidth: element.clientWidth, outside: [...element.querySelectorAll('*')].filter(child => child.getBoundingClientRect().right > right + 1 || child.scrollWidth > child.clientWidth + 1).slice(0, 8).map(child => ({ tag: child.tagName, text: child.textContent?.slice(0, 40), class: child.className, scroll: child.scrollWidth, client: child.clientWidth })) };
    });
    if (width <= 640) expect(overflow.scrollWidth, JSON.stringify(overflow)).toBeGreaterThan(overflow.clientWidth);
    else expect(overflow.scrollWidth, JSON.stringify(overflow)).toBeLessThanOrEqual(overflow.clientWidth);
    const overhangingHeaders = await scroll.locator('th').evaluateAll(headers => headers.flatMap(header => {
      const button = header.querySelector('button'); return button && button.getBoundingClientRect().right > header.getBoundingClientRect().right + 1 ? [{ label: button.textContent, cell: header.getBoundingClientRect().width, button: button.getBoundingClientRect().width }] : [];
    }));
    expect(overhangingHeaders).toEqual([]);
    await cell.scrollIntoViewIfNeeded();
    const cellBox = await cell.boundingBox(), removeBox = await remove.boundingBox(), toggleBox = await toggle.boundingBox();
    const regionBox = await scroll.boundingBox();
    expect(toggleBox!.x).toBeGreaterThanOrEqual(regionBox!.x);
    expect(removeBox!.x + removeBox!.width).toBeLessThanOrEqual(regionBox!.x + regionBox!.width);
    expect(toggleBox!.y).toBe(removeBox!.y);
    expect(toggleBox!.x + toggleBox!.width).toBeLessThanOrEqual(cellBox!.x + cellBox!.width);
    const toggleSize = await toggle.evaluate(element => ({ scroll: element.scrollWidth, client: element.clientWidth, width: element.getBoundingClientRect().width, font: getComputedStyle(element).fontSize, padding: getComputedStyle(element).padding }));
    expect(toggleSize.scroll, JSON.stringify(toggleSize)).toBeLessThanOrEqual(toggleSize.client);
    expect(await remove.evaluate(element => element.scrollWidth <= element.clientWidth)).toBe(true);
    await toggle.click(); await expect(toggle).toHaveAttribute('aria-pressed', 'false');
    const optionalBox = await toggle.boundingBox();
    expect(optionalBox!.x + optionalBox!.width).toBeLessThanOrEqual(cellBox!.x + cellBox!.width);
    expect(await toggle.evaluate(element => element.scrollWidth <= element.clientWidth)).toBe(true);
  });
}

for (const width of [390, 1280]) {
  test(`marks the section currently reached by scrolling at ${width}px`, async ({ page }) => {
    await page.setViewportSize({ width, height: 720 });
    await page.route('**/api/**', route => route.abort());
    await page.goto('/tests/e2e/fixtures/recording-selected-host.html');
    const nav = page.getByRole('navigation', { name: '本頁設定區段' });
    const buttons = [nav.getByRole('button', { name: '1 來源' }), nav.getByRole('button', { name: '2 記錄方式' }), nav.getByRole('button', { name: '3 目的地與確認' })];
    await expect(buttons[0]).toHaveAttribute('aria-current', 'step');
    for (const [index, section] of ['method', 'destination', 'sources'].entries()) {
      await page.evaluate(name => document.getElementById(`recording-editor-${name}`)!.scrollIntoView({ block: 'start' }), section);
      const active = [1, 2, 0][index];
      await expect(buttons[active]).toHaveAttribute('aria-current', 'step');
      expect(await buttons[active].evaluate(element => getComputedStyle(element).backgroundColor)).not.toBe(await buttons[(active + 1) % 3].evaluate(element => getComputedStyle(element).backgroundColor));
    }
    await buttons[2].click(); await expect(buttons[2]).toHaveAttribute('aria-current', 'step');
  });
}

test('explains required points and single-word decoding in accessible help', async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 720 });
  await page.route('**/api/**', route => route.abort());
  await page.goto('/tests/e2e/fixtures/recording-selected-host.html');
  const table = page.getByTestId('group-member-table');
  await table.locator('[data-parameter-help="member.required"]').click();
  const note = page.getByRole('note');
  await expect(note).toContainText('不保證一定寫入'); await expect(note).toContainText('儲存記錄設定');
  await expect(note).toContainText('NULL');
  await page.keyboard.press('Escape'); await expect(note).toHaveCount(0);
  await table.locator('[data-parameter-help="point.order"]').click();
  await expect(page.getByRole('note')).toContainText('int16');
  await expect(page.getByRole('note')).toContainText('一個 16 位元字');
});

test('edits a selected name, changes word-order intent and sorts without changing membership', async ({ page }) => {
  await page.setViewportSize({ width: 1280, height: 720 });
  await page.route('**/api/**', route => route.abort());
  await page.goto('/tests/e2e/fixtures/recording-selected-host.html');
  const table = page.getByTestId('group-member-table');
  const first = table.getByTestId('group-member-row-P0');
  const name = first.getByRole('textbox', { name: '點位名稱 TAG_0 (D0)' });
  await name.focus(); await expect(name).toBeEditable();
  await name.fill('主馬達'); await name.press('Enter');
  await expect(first.getByRole('textbox', { name: '點位名稱 主馬達' })).toHaveValue('主馬達');
  const order = first.getByRole('combobox', { name: '字序' });
  await order.selectOption('CDAB');
  await expect(first.getByRole('button', { name: '儲存字序' })).toBeVisible();
  const saveBox = await first.getByRole('button', { name: '儲存字序' }).boundingBox(), typeBox = await first.locator('td[data-column="dataType"]').boundingBox();
  expect(saveBox!.x + saveBox!.width).toBeLessThanOrEqual(typeBox!.x + 1);
  await expect(page.getByTestId('group-member-scope')).toContainText('32 點');
  await table.getByTestId('group-member-search').fill('TAG_1');
  await expect(page.getByTestId('group-member-scope')).toContainText('顯示 11 / 32 點');
  await table.getByTestId('group-member-search').clear();
  await expect(page.getByTestId('group-member-scope')).not.toContainText('顯示');
  const rowIds = () => table.locator('tbody tr').evaluateAll(rows => rows.map(row => row.getAttribute('data-testid')));
  expect((await rowIds()).slice(0, 3)).toEqual(['group-member-row-P0', 'group-member-row-P1', 'group-member-row-P2']);
  await table.getByRole('button', { name: '位址', exact: true }).click();
  expect((await rowIds()).slice(0, 3)).toEqual(['group-member-row-P31', 'group-member-row-P30', 'group-member-row-P29']);
  await expect(page.getByTestId('group-member-scope')).toContainText('32 點');
});

for (const width of [320, 390, 1280]) {
  test(`keeps the optional entity key usable at ${width}px`, async ({ page }) => {
    await page.setViewportSize({ width, height: 720 });
    await page.route('**/api/**', route => route.abort());
    await page.goto('/tests/e2e/fixtures/recording-selected-host.html?entity');
    const scroll = page.getByTestId('group-member-table').getByRole('region');
    const cell = page.getByTestId('group-member-row-P0').locator('td').last();
    const input = cell.locator('input:not([type="checkbox"])');
    await input.scrollIntoViewIfNeeded();
    const cellBox = await cell.boundingBox(), inputBox = await input.boundingBox();
    const metrics = await scroll.evaluate(element => ({ scroll: element.scrollWidth, client: element.clientWidth, tracks: getComputedStyle(element.querySelector('table')!).gridTemplateColumns, cells: [...element.querySelectorAll('tbody tr:first-child td')].map(cell => ({ name: cell.getAttribute('data-column') ?? cell.getAttribute('data-testid'), width: cell.getBoundingClientRect().width, scroll: cell.scrollWidth, client: cell.clientWidth })) }));
    if (width <= 640) expect(metrics.scroll, JSON.stringify(metrics)).toBeGreaterThan(metrics.client);
    else expect(metrics.scroll, JSON.stringify(metrics)).toBeLessThanOrEqual(metrics.client);
    expect(inputBox!.width).toBeGreaterThanOrEqual(50);
    expect(inputBox!.x + inputBox!.width).toBeLessThanOrEqual(cellBox!.x + cellBox!.width);
    await input.fill('entity-1'); await expect(input).toHaveValue('entity-1');
  });
}

for (const width of [320, 390, 768, 1001, 1280, 1440]) {
  test(`keeps per-point calculation and interval controls usable at ${width}px`, async ({ page }) => {
    await page.setViewportSize({ width, height: 720 });
    await page.route('**/api/**', route => route.abort());
    await page.goto('/tests/e2e/fixtures/recording-selected-host.html?mixed');
    const scroll = page.getByTestId('group-member-table').getByRole('region');
    const row = page.getByTestId('group-member-row-P0');
    const cadence = row.getByTestId('group-member-cadence-P0');
    await cadence.scrollIntoViewIfNeeded();
    const cellBox = await cadence.boundingBox();
    const metrics = await scroll.evaluate(element => ({ scroll: element.scrollWidth, client: element.clientWidth, tracks: getComputedStyle(element.querySelector('table')!).gridTemplateColumns, cells: [...element.querySelectorAll('tbody tr:first-child td')].map(cell => ({ name: cell.getAttribute('data-column') ?? cell.getAttribute('data-testid'), width: cell.getBoundingClientRect().width, scroll: cell.scrollWidth, client: cell.clientWidth })) }));
    if (width <= 1024) expect(metrics.scroll, JSON.stringify(metrics)).toBeGreaterThan(metrics.client);
    else expect(metrics.scroll, JSON.stringify(metrics)).toBeLessThanOrEqual(metrics.client);
    for (const control of await cadence.locator('input, select').all()) {
      const box = await control.boundingBox();
      expect(box!.width).toBeGreaterThanOrEqual(40);
      expect(box!.x + box!.width, JSON.stringify({ control: box, cell: cellBox, metrics })).toBeLessThanOrEqual(cellBox!.x + cellBox!.width);
    }
    await expect(row.getByTestId('group-member-calculation-P0').getByRole('combobox')).toBeVisible();
    const interval = cadence.locator('input'); await interval.fill('30'); await expect(interval).toHaveValue('30');
    const order = row.getByRole('combobox', { name: '字序' }); await order.selectOption('CDAB');
    const orderCellBox = await row.locator('td[data-column="wordOrder"]').boundingBox(), selectBox = await order.boundingBox();
    const saveBox = await row.getByRole('button', { name: '儲存字序' }).boundingBox();
    const orderMetrics = await row.evaluate(element => ({ tracks: getComputedStyle(element.closest('table')!).gridTemplateColumns, cells: [...element.querySelectorAll('td')].map(cell => ({ name: cell.getAttribute('data-column') ?? cell.getAttribute('data-testid'), width: cell.getBoundingClientRect().width })) }));
    expect(selectBox!.width, JSON.stringify(orderMetrics)).toBeGreaterThanOrEqual(60);
    expect(saveBox!.x + saveBox!.width).toBeLessThanOrEqual(orderCellBox!.x + orderCellBox!.width);
    expect(await scroll.evaluate(element => element.scrollWidth > element.clientWidth)).toBe(width <= 1024);
  });
}
