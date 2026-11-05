import { expect, test, type Page } from '@playwright/test';

// MOCK SAMPLE presentation only. No gateway/HSL/device I/O or production recovery is exercised.
const layouts = ['grid', 'navigator', 'map', 'chapters', 'detail'] as const;
async function control(page: Page, label: string, value: string) {
  const controls = page.getByTestId('fixture-controls');
  await controls.locator('summary').click();
  await controls.getByLabel(label, { exact: true }).selectOption(value);
  await controls.locator('summary').click();
}

for (const width of [1440, 768, 390]) {
  for (const light of [true, false]) {
    for (const layout of layouts) {
      test(`${layout} mock recovery and reading states at ${width}px ${light ? 'light' : 'dark'}`, async ({ page }) => {
        await page.setViewportSize({ width, height: width === 390 ? 844 : 1000 });
        const requests: string[] = [];
        await page.route('**/api/**', route => { requests.push(route.request().url()); return route.abort(); });
        await page.goto(`/tests/e2e/fixtures/recording-home-memory-host.html?scenario=mc&en&controls${light ? '&light' : ''}`);
        await expect(page.getByTestId('fixture-sample-notice')).toContainText('MOCK SAMPLE');
        await page.getByTestId(`home-memory-layout-${layout}`).click();
        const view = page.getByTestId('home-memory-view'), action = page.getByTestId('fixture-last-action');
        await expect(view).toHaveAttribute('data-layout', layout);
        const used = page.getByTestId('home-memory-cell-D0'), free = page.getByTestId('home-memory-cell-D33');
        const detail = page.getByTestId('home-memory-detail');

        await test.step('all three source blocks prohibit creation but preserve known owners', async () => {
          for (const [block, notice] of [
            ['recovery', 'A new read range save is still unconfirmed.'],
            ['local', 'Some range changes on this page have not been saved'],
            ['stale', 'The saved ranges changed since this page loaded'],
          ]) {
            await control(page, 'Mock source blocking', block);
            await expect(page.getByRole('alert')).toContainText(notice);
            await expect(used).toHaveAttribute('data-state', 'used');
            await used.click();
            await expect(action).toHaveText('open:I');
            await expect(free).toHaveAttribute('data-state', 'unconfirmed');
            expect(await free.evaluate(node => node.tagName)).toBe('DIV');
            await free.focus(); await page.keyboard.press('Enter'); await free.click();
            await expect(action).toHaveText('open:I');
            if (layout === 'detail' || layout === 'map') {
              await expect(detail.locator('.gw-memory-detail-address')).toHaveText('D33');
              await expect(detail.getByTestId('home-memory-detail-action')).toHaveCount(0);
            }
            await expect(view).toHaveAttribute('data-layout', layout);
          }
          await control(page, 'Mock source blocking', 'none');
          await expect(free).toHaveAttribute('data-state', 'empty');
          await free.click();
          await expect(action).toHaveText('add:D33');
        });

        await test.step('pending and failed range reads show no cells, then retry restores the same layout', async () => {
          const visibleKeys = page.getByTestId('fixture-visible-keys');
          const key = await used.getAttribute('data-value-key');
          expect(key).not.toBeNull();
          const expectUsedObserved = async () => {
            await used.focus(); await used.scrollIntoViewIfNeeded();
            await expect(used).toHaveAttribute('data-selected', 'true');
            await expect.poll(async () => JSON.parse(await visibleKeys.textContent() ?? '[]') as string[]).toContain(key);
          };
          for (const state of ['loading', 'checking', 'error']) {
            // Each case starts populated, so an empty report cannot pass on a prior empty state.
            await expectUsedObserved();
            await control(page, 'Mock range read', state);
            await expect(page.getByTestId(/^home-memory-cell-/)).toHaveCount(0);
            // The host remains mounted; only the production view reports this cleared key list.
            await expect(visibleKeys).toHaveText('[]');
            if (state === 'checking') await expect(page.getByRole('status')).toHaveText('Checking saved MC points…');
            else if (state === 'error') await expect(page.getByRole('alert')).toContainText('Saved ranges could not be read');
            else await expect(page.getByRole('status')).toHaveText('Reading the current workspace…');
            await expect(view).toHaveAttribute('data-layout', layout);
            await expect(action).toHaveText('add:D33');
            if (state !== 'error') await control(page, 'Mock range read', 'ready');
          }
          await page.getByRole('button', { name: 'Read ranges again', exact: true }).click();
          await expect(page.getByTestId('fixture-retry-count')).toHaveText('1');
          await expect(used).toHaveAttribute('data-state', 'used');
          await expect(view).toHaveAttribute('data-layout', layout);
          await expect(action).toHaveText('add:D33');
          await expectUsedObserved();
        });

        await test.step('sample loading, zero, changed, stale and error values stay distinct', async () => {
          const modes = [
            { mode: 'loading', status: 'loading', state: 'pending', text: '…', freshness: 'Waiting for reading' },
            { mode: 'live', status: 'live', state: 'value', text: '0', freshness: 'Stored reading' },
            { mode: 'updated', status: 'live', state: 'value', text: '37', freshness: 'Stored reading' },
            { mode: 'stale', status: 'live', state: 'stale', text: '0', freshness: 'Stale reading' },
            { mode: 'cell-error', status: 'live', state: 'error', text: '!', freshness: 'Read failed' },
            { mode: 'error', status: 'error', state: 'stale', text: '0', freshness: 'Stale reading' },
            { mode: 'polling', status: 'polling', state: 'value', text: '0', freshness: 'Stored reading' },
          ];
          for (const mode of modes) {
            await control(page, 'Mock value presentation', mode.mode);
            await expect(page.getByTestId('home-memory-values-status')).toHaveAttribute('data-status', mode.status);
            await expect(used.locator('.gw-memory-reading')).toHaveAttribute('data-reading', mode.state);
            await expect(used.locator('.gw-memory-reading')).toHaveText(mode.text);
            await used.focus();
            if (layout === 'detail' || layout === 'map') {
              await expect(detail.locator('.gw-memory-detail-reading')).toHaveText(mode.text);
              await expect(detail).toContainText(mode.freshness);
            }
            await expect(used).toHaveAttribute('data-state', 'used');
            await expect(free).toHaveAttribute('data-state', 'empty');
            await expect(view).toHaveAttribute('data-layout', layout);
          }
          await expect(view).not.toContainText('MOCK_READ_FAILURE');
          await used.scrollIntoViewIfNeeded();
          const key = await used.getAttribute('data-value-key');
          expect(key).not.toBeNull();
          await expect.poll(async () => JSON.parse(await page.getByTestId('fixture-visible-keys').textContent() ?? '[]') as string[]).toContain(key);
          await used.click();
          await expect(action).toHaveText('open:I');
        });
        expect(requests).toEqual([]);
      });
    }
  }
}
