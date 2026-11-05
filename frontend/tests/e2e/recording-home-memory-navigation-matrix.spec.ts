import { expect, test, type Page } from '@playwright/test';

// Isolated MOCK SAMPLE fixture. Reopen tests a fresh memory component, not the production modal shell.
const layouts = ['grid', 'navigator', 'map', 'chapters', 'detail'] as const;
async function locate(page: Page, address: string) {
  await page.getByRole('textbox', { name: 'Go to address', exact: true }).fill(address);
  await page.getByRole('button', { name: 'Go', exact: true }).click();
}

for (const width of [1440, 768, 390]) {
  for (const light of [true, false]) {
    for (const layout of layouts) {
      test(`${layout} cancel, paging, locate and fresh reopen at ${width}px ${light ? 'light' : 'dark'}`, async ({ page }) => {
        await page.setViewportSize({ width, height: width === 390 ? 844 : 1000 });
        const requests: string[] = [];
        await page.route('**/api/**', route => { requests.push(route.request().url()); return route.abort(); });
        await page.goto(`/tests/e2e/fixtures/recording-home-memory-host.html?scenario=conflict&en${light ? '&light' : ''}`);
        await expect(page.getByTestId('fixture-sample-notice')).toContainText('MOCK SAMPLE');
        await page.getByTestId(`home-memory-layout-${layout}`).click();
        const view = page.getByTestId('home-memory-view'), action = page.getByTestId('fixture-last-action');
        const navigation = page.getByRole('group', { name: 'Address blocks', exact: true });
        const previous = navigation.getByRole('button', { name: /Previous/ });
        const next = navigation.getByRole('button', { name: /Next/ });
        const grid = page.locator('.gw-memory-grid');
        await expect(view).toHaveAttribute('data-layout', layout);

        await test.step('cancel dismisses owner selection and returns keyboard focus to the exact cell', async () => {
          const conflict = page.getByTestId('home-memory-cell-40003');
          for (let repeat = 0; repeat < 2; repeat++) {
            await conflict.focus(); await page.keyboard.press('Enter');
            const owners = page.getByTestId('home-memory-owners');
            await expect(owners).toBeVisible();
            await expect(owners.getByRole('heading')).toBeFocused();
            const cancel = owners.getByRole('button', { name: 'Cancel', exact: true });
            await cancel.focus(); await page.getByTestId('home-memory-cell-40005').hover();
            await expect(cancel).toBeFocused();
            await cancel.click();
            await expect(owners).toHaveCount(0);
            await expect(conflict).toBeFocused();
            await expect(action).toHaveText('');
          }
          await page.keyboard.press('ArrowRight');
          await expect(page.getByTestId('home-memory-cell-40004')).toBeFocused();
        });

        await test.step('next and previous move one bounded window and restore the initial boundary', async () => {
          await expect(previous).toBeDisabled(); await expect(next).toBeEnabled();
          const initial = await grid.locator('[data-offset]').evaluateAll(nodes => nodes.map(node => Number((node as HTMLElement).dataset.offset)));
          expect(initial.length).toBeGreaterThan(0); expect(initial.length).toBeLessThanOrEqual(256);
          await next.click();
          await expect(grid.locator('[data-offset]').first()).toHaveAttribute('data-offset', String(initial.at(-1)! + 1));
          await expect(previous).toBeEnabled();
          await next.click();
          await expect(grid.locator('[data-offset]').first()).toHaveAttribute('data-offset', String(initial[0] + initial.length * 2));
          await previous.click(); await previous.click();
          await expect(grid.locator('[data-offset]').first()).toHaveAttribute('data-offset', String(initial[0]));
          await expect(previous).toBeDisabled();
          await locate(page, '499999');
          await expect(grid.locator('[data-offset]').last()).toHaveAttribute('data-offset', '99999');
          await expect(next).toBeDisabled(); await expect(previous).toBeEnabled();
          await locate(page, '40001');
          await expect(previous).toBeDisabled(); await expect(next).toBeEnabled();
          await expect(action).toHaveText('');
        });

        await test.step('invalid addresses cannot replace the current window or invoke any action', async () => {
          const input = page.getByRole('textbox', { name: 'Go to address', exact: true });
          for (const invalid of ['', 'D0', '30001', 'not-an-address', '-1', '999999999999999999999999']) {
            await locate(page, invalid);
            await expect(input).toHaveAttribute('aria-invalid', 'true');
            await expect(page.getByRole('alert')).toHaveText('That address is not in the 4x area.');
            await expect(page.getByTestId('home-memory-cell-40001')).toBeVisible();
            await expect(action).toHaveText('');
          }
          await input.fill('40021');
          await expect(input).toHaveAttribute('aria-invalid', 'false');
          await page.keyboard.press('Enter');
          await expect(page.getByTestId('home-memory-cell-40021')).toHaveAttribute('tabindex', '0');
          await expect(page.getByRole('alert')).toHaveCount(0);
          await expect(page.getByTestId('home-memory-cell-40001')).toHaveCount(0);
          await expect(action).toHaveText('');
        });

        await test.step('Back and Close each unmount and repeatedly reopen with fresh shared state', async () => {
          let instance = 1;
          for (const dismiss of ['back', 'close', 'back', 'close']) {
            await page.getByTestId('home-memory-area-IR').click();
            await locate(page, '30201');
            await page.getByTestId(`home-memory-layout-${layout}`).click();
            await page.getByTestId(`home-memory-${dismiss}`).click();
            await expect(page.getByRole('dialog')).toHaveCount(0);
            await expect(view).toHaveCount(0);
            await expect(action).toHaveText(dismiss);
            await expect(page.getByTestId('fixture-visible-keys')).toHaveText('[]');
            await page.getByTestId('fixture-reopen').click();
            await expect(page.getByRole('dialog')).toHaveAttribute('data-fixture-instance', String(++instance));
            await expect(view).toHaveAttribute('data-layout', 'grid');
            await expect(page.getByTestId('home-memory-layout-grid')).toHaveAttribute('aria-pressed', 'true');
            await expect(page.getByTestId('home-memory-area-HR')).toHaveAttribute('aria-selected', 'true');
            await expect(page.getByRole('textbox', { name: 'Go to address', exact: true })).toHaveValue('');
            await expect(page.getByTestId('home-memory-cell-40001')).toHaveAttribute('tabindex', '0');
            await expect(page.getByTestId('home-memory-owners')).toHaveCount(0);
            await expect(page.getByRole('heading', { name: 'Address areas', exact: true })).toBeFocused();
          }
        });
        expect(requests).toEqual([]);
      });
    }
  }
}
