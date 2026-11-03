import { expect, test } from '@playwright/test';

for (const theme of ['light', 'dark']) for (const width of [390, 768, 1440]) {
  test(`${theme} editor uses semantic canvas and selected navigation at ${width}px`, async ({ page }, info) => {
    await page.setViewportSize({ width, height: 820 });
    await page.route('**/api/**', route => route.abort());
    await page.goto(`/tests/e2e/fixtures/recording-editor-real-host.html?theme=${theme}`);
    await expect(page.getByTestId('group-editor')).toBeVisible();
    const colors = await page.evaluate(() => {
      const token = (name: string) => {
        const probe = document.createElement('span'); probe.style.color = `var(--gw-${name})`; document.body.append(probe);
        const value = getComputedStyle(probe).color; probe.remove(); return value;
      };
      const canvas = getComputedStyle(document.querySelector('[data-workbench-v2]')!).backgroundColor;
      const panel = getComputedStyle(document.querySelector('.recording-editor-main')!).backgroundColor;
      const current = getComputedStyle(document.querySelector('.recording-editor-steps [aria-current="step"]')!);
      const back = getComputedStyle(document.querySelector('.gw-simple-recording-title button')!);
      const check = getComputedStyle(document.querySelector('.recording-editor-panel input[type="checkbox"]')!).accentColor;
      return { canvas, panel, check, back: back.backgroundColor, backBorder: back.borderColor, selected: current.backgroundColor, selectedText: current.color, ground: token('ground'), surface: token('surface'), control: token('control'), active: token('active'), accent: token('accent') };
    });
    expect(colors.canvas).toBe(colors.ground); expect(colors.panel).toBe(colors.surface); expect(colors.back).toBe(colors.surface); expect(colors.backBorder).toBe(colors.control);
    expect(colors.selected).toBe(colors.active); expect(colors.selectedText).toBe(colors.accent); expect(colors.check).toBe(colors.accent);
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
    await page.screenshot({ path: info.outputPath('editor.png') });
  });
}

test('editor typography, spacing and controls respond to their design tokens', async ({ page }) => {
  await page.setViewportSize({ width: 1440, height: 820 });
  await page.route('**/api/**', route => route.abort());
  await page.goto('/tests/e2e/fixtures/recording-editor-real-host.html');
  const styles = await page.evaluate(() => {
    const values = { '--gw-font-sans': 'monospace', '--gw-typography-title-size': '29px', '--gw-typography-section-size': '21px', '--gw-typography-body-size': '16px', '--gw-typography-small-size': '13px', '--gw-space-md': '18px', '--gw-space-lg': '26px', '--gw-control-height': '46px', '--gw-control-mobile': '48px', '--gw-radius-sm': '9px', '--gw-focus-width': '3px', '--gw-focus-offset': '3px' };
    for (const [name, value] of Object.entries(values)) document.documentElement.style.setProperty(name, value);
    const css = (selector: string) => getComputedStyle(document.querySelector(selector)!);
    return { titleFont: css('.gw-simple-recording-title h1').fontFamily, editorFont: css('.recording-editor-simple').fontFamily, title: css('.gw-simple-recording-title h1').fontSize, section: css('.recording-editor-panel header h4').fontSize,
      label: css('.recording-simple-name label').fontSize, summary: css('.recording-editor-summary dt').fontSize,
      sectionGap: css('.recording-editor-panel + .recording-editor-panel').marginTop, helpWidth: css('[data-parameter-help="record.name"]').width,
      fieldMinHeight: css('[data-testid="group-name"]').minHeight, fieldRadius: css('[data-testid="group-name"]').borderRadius };
  });
  expect(styles).toEqual({ titleFont: 'monospace', editorFont: 'monospace', title: '29px', section: '21px', label: '16px', summary: '13px', sectionGap: '26px', helpWidth: '48px', fieldMinHeight: '46px', fieldRadius: '9px' });
  const back = page.getByRole('button', { name: '返回總覽' }), name = page.getByTestId('group-name');
  await back.focus(); expect(await back.evaluate(element => getComputedStyle(element).outlineWidth)).toBe('3px');
  await name.focus(); expect(await name.evaluate(element => getComputedStyle(element).outlineWidth)).toBe('3px');
});

test('selected points follow the same typography, spacing and focus tokens', async ({ page }) => {
  await page.setViewportSize({ width: 1280, height: 820 });
  await page.route('**/api/**', route => route.abort());
  await page.goto('/tests/e2e/fixtures/recording-selected-host.html');
  await page.evaluate(() => {
    const values = { '--gw-typography-small-size': '13px', '--gw-space-2xs': '6px', '--gw-control-height': '46px', '--gw-control-mobile': '48px', '--gw-focus-width': '3px', '--gw-focus-offset': '3px' };
    for (const [name, value] of Object.entries(values)) document.documentElement.style.setProperty(name, value);
  });
  const region = page.getByTestId('group-member-table').getByRole('region'); await region.focus();
  const styles = await page.evaluate(() => {
    const css = (selector: string) => getComputedStyle(document.querySelector(selector)!);
    return { header: css('.recording-compact-members th').fontSize, hint: css('.recording-compact-members__name-hint').fontSize,
      searchGap: css('.recording-compact-members__search').gap, actionMinHeight: css('.recording-compact-members__actions button').minHeight,
      actionMinWidth: css('.recording-compact-members__actions button').minWidth, focusWidth: css('.recording-selected-table-scroll').outlineWidth, focusOffset: css('.recording-selected-table-scroll').outlineOffset };
  });
  expect(styles).toEqual({ header: '13px', hint: '13px', searchGap: '6px', actionMinHeight: '46px', actionMinWidth: '48px', focusWidth: '3px', focusOffset: '3px' });
});

for (const theme of ['light', 'dark']) for (const width of [390, 1440]) test(`${theme} selected points use semantic control and selection colors at ${width}px`, async ({ page }, info) => {
  await page.setViewportSize({ width, height: 820 });
  await page.route('**/api/**', route => route.abort());
  await page.goto(`/tests/e2e/fixtures/recording-selected-host.html?theme=${theme}`);
  const table = page.getByTestId('group-member-table'), row = table.getByTestId('group-member-row-P0');
  const required = row.locator('button[aria-pressed]'), input = row.getByRole('textbox');
  const colors = await row.evaluate(element => {
    const token = (name: string) => { const probe = document.createElement('span'); probe.style.color = `var(--gw-${name})`; document.body.append(probe); const color = getComputedStyle(probe).color; probe.remove(); return color; };
    const button = getComputedStyle(element.querySelector('button[aria-pressed]')!), field = getComputedStyle(element.querySelector('input')!);
    const placeholder = getComputedStyle(document.querySelector('.recording-compact-members__search input')!, '::placeholder').color;
    return { selected: button.backgroundColor, selectedText: button.color, field: field.backgroundColor, border: field.borderColor, placeholder, active: token('active'), accent: token('accent'), surface: token('surface'), control: token('control'), faint: token('faint') };
  });
  expect(colors.selected).toBe(colors.active); expect(colors.selectedText).toBe(colors.accent);
  expect(colors.field).toBe(colors.surface); expect(colors.border).toBe(colors.control); expect(colors.placeholder).toBe(colors.faint);
  expect((await required.boundingBox())!.height).toBeGreaterThanOrEqual(width < 640 ? 44 : 40);
  if (width < 640) expect(await table.locator('th[data-column="device"] button').evaluate(element => parseFloat(getComputedStyle(element).fontSize))).toBeGreaterThanOrEqual(14);
  await input.focus(); await expect(input).toBeEditable();
  await required.click(); await expect(required).toHaveAttribute('aria-pressed', 'false');
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth)).toBe(true);
  await table.scrollIntoViewIfNeeded(); await page.screenshot({ path: info.outputPath('points.png') });
});

test('theme colors can change without losing an unsaved recording name', async ({ page }) => {
  await page.setViewportSize({ width: 1280, height: 820 });
  await page.route('**/api/**', route => route.abort());
  await page.goto('/tests/e2e/fixtures/recording-editor-real-host.html?theme=dark');
  const name = page.getByTestId('group-name'); await name.fill('尚未儲存的名稱');
  await page.evaluate(() => { document.documentElement.classList.replace('dark', 'light'); });
  await expect(name).toHaveValue('尚未儲存的名稱');
  const colors = await page.evaluate(() => {
    const probe = document.createElement('span'); probe.style.backgroundColor = 'var(--gw-surface)'; document.body.append(probe);
    const expected = getComputedStyle(probe).backgroundColor; probe.remove();
    return { surface: getComputedStyle(document.querySelector('.recording-editor-main')!).backgroundColor, expected };
  });
  expect(colors.surface).toBe('rgb(255, 255, 255)'); expect(colors.surface).toBe(colors.expected);
});

test('reduced-motion preference removes editor control transitions', async ({ page }) => {
  await page.emulateMedia({ reducedMotion: 'reduce' });
  await page.route('**/api/**', route => route.abort());
  await page.goto('/tests/e2e/fixtures/recording-editor-real-host.html');
  const navigation = page.getByRole('navigation', { name: '本頁設定區段' }).getByRole('button').first();
  expect(await navigation.evaluate(element => getComputedStyle(element).transitionDuration)).toBe('0s');
});
