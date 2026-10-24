import { expect, test } from '@playwright/test';

for (const viewport of [{ width: 1600, height: 900 }, { width: 1024, height: 720 }]) {
  test(`mapping uses the same page layout as the other steps at ${viewport.width}px`, async ({ page }, testInfo) => {
    await page.setViewportSize(viewport);
    await page.route('**/api/**', route => route.fulfill({ json: { data: [] } }));
    await page.goto('/tests/e2e/fixtures/studio-mapping-scroll.html');
    await expect(page.getByTestId('mapping-table-wrapper').locator('tbody tr')).toHaveCount(48);

    const layout = () => page.evaluate(() => {
      const main = document.getElementById('step-content')!;
      const aside = document.querySelector('[data-testid=sidebar-rail-aside]')!;
      return {
        mainOverflow: getComputedStyle(main).overflowY,
        sidebarLeft: aside.getBoundingClientRect().left,
        sidebarWidth: aside.getBoundingClientRect().width,
        extraScrollers: [...document.querySelectorAll('#step-content, .wbv2-workspace > aside')]
          .filter(element => ['auto', 'scroll'].includes(getComputedStyle(element).overflowY)).length,
        horizontalOverflow: document.documentElement.scrollWidth > innerWidth,
      };
    });
    const mappingLayout = await layout();
    expect(mappingLayout.extraScrollers).toBe(0);
    expect(mappingLayout.horizontalOverflow).toBe(false);

    const tableStyle = (testId: string) => page.getByTestId(testId).evaluate(element => {
      const card = getComputedStyle(element);
      const scroller = getComputedStyle(element.querySelector('table')!.parentElement!);
      const heading = getComputedStyle(element.querySelector('thead')!);
      return {
        display: card.display, direction: card.flexDirection, gap: card.gap,
        padding: card.padding, background: card.backgroundColor,
        border: card.borderColor, radius: card.borderRadius,
        maxHeight: scroller.maxHeight, scrollerBorder: scroller.borderColor,
        scrollerRadius: scroller.borderRadius, scrollerBackground: scroller.backgroundColor,
        headingPosition: heading.position, headingBackground: heading.backgroundColor,
      };
    });
    const mappingTableStyle = await tableStyle('mapping-table-wrapper');
    expect(mappingTableStyle.maxHeight).toBe('520px');
    await page.getByTestId('step-nav-button-2').click();
    await expect(page.getByTestId('merged-point-table-container')).toBeAttached();
    const mergedPointTableStyle = await tableStyle('merged-point-table-container');
    expect(mergedPointTableStyle.maxHeight).toBe('300px');
    const { maxHeight: _m1, ...mappingCardStyle } = mappingTableStyle;
    const { maxHeight: _m2, ...mergedCardStyle } = mergedPointTableStyle;
    expect(mappingCardStyle).toEqual(mergedCardStyle);
    await page.getByTestId('step-nav-button-3').click();

    // Compare actual behavior with the shared Device page, instead of pinning a special Step 3 class.
    await page.getByTestId('step-nav-button-1').click();
    const deviceLayout = await layout();
    expect(mappingLayout).toEqual(deviceLayout);
    await page.getByTestId('step-nav-button-3').click();
    await expect(page.getByTestId('mapping-table-wrapper')).toBeVisible();

    // The page remains naturally scrollable and its sticky top bar stays visible.
    await page.mouse.move(viewport.width - 3, 400);
    await page.mouse.wheel(0, 250);
    await expect.poll(() => page.evaluate(() => scrollY)).toBeGreaterThan(0);
    await expect(page.locator('.wbv2-topbar')).toBeInViewport();
    await page.getByTestId('btn-continue').scrollIntoViewIfNeeded();
    await expect(page.getByTestId('btn-continue')).toBeInViewport();
    await page.evaluate(() => window.scrollTo(0, 0));
    await page.screenshot({ path: testInfo.outputPath('mapping-shared-layout.png') });
  });
}
