import { expect, test, type Page } from '@playwright/test';

const fields = [
  ['record.table', 'group-table'],
  ['record.schema', 'group-schema'],
  ['policy.interval', 'group-interval'],
  ['policy.lateness', 'group-lateness'],
  ['policy.incomplete', 'group-incomplete'],
  ['policy.dedupe', 'group-dedupe'],
  ['policy.entityColumn', 'group-entity-column'],
  ['policy.provenanceColumn', 'group-provenance-column'],
] as const;

type FieldGeometry = {
  label: { top: number; bottom: number; left: number; right: number };
  controls: Array<{ top: number; bottom: number; left: number; right: number }>;
  help: { top: number; bottom: number; left: number; right: number };
  controlEdgeProbes: Array<{ ok: boolean; y: number; hit: string | null; hitTestId: string | null; hitHelp: string | null }>;
};

async function fieldGeometry(page: Page, helpKey: string, controlTestId: string): Promise<FieldGeometry> {
  return page.locator(`[data-parameter-help="${helpKey}"]`).evaluate((button, testId) => {
    const label = button.previousElementSibling;
    if (!(label instanceof HTMLElement)) throw new Error(`missing label for ${helpKey}`);
    const textNode = [...label.childNodes].find(node => node.nodeType === Node.TEXT_NODE && node.textContent?.trim());
    if (!textNode) throw new Error(`missing label text for ${helpKey}`);
    const range = document.createRange();
    range.selectNodeContents(textNode);
    const labelBox = range.getBoundingClientRect();
    const controls = [...label.children]
      .filter(element => ['INPUT', 'SELECT', 'DETAILS'].includes(element.tagName))
      .map(element => element.getBoundingClientRect())
      .map(rect => ({ top: rect.top, bottom: rect.bottom, left: rect.left, right: rect.right }));
    const expected = document.querySelector(`[data-testid="${testId}"]`);
    if (!expected || !controls.some(control => Math.abs(control.top - expected.getBoundingClientRect().top) < 1)) {
      throw new Error(`missing first control for ${helpKey}`);
    }
    const helpBox = button.getBoundingClientRect();
    const helpColumnX = (helpBox.left + helpBox.right) / 2;
    const controlElements = [...label.children].filter(element => ['INPUT', 'SELECT', 'DETAILS'].includes(element.tagName));
    const controlEdgeProbes = controlElements.flatMap(element => {
      const box = element.getBoundingClientRect();
      return [box.top + 2, box.bottom - 2].map(y => {
        const hit = document.elementFromPoint(helpColumnX, y);
        return {
          ok: hit === element || hit?.closest(`[data-testid="${testId}"]`) === element,
          y,
          hit: hit?.tagName ?? null,
          hitTestId: hit?.getAttribute('data-testid') ?? null,
          hitHelp: hit?.getAttribute('data-parameter-help') ?? null,
        };
      });
    });
    return {
      label: { top: labelBox.top, bottom: labelBox.bottom, left: labelBox.left, right: labelBox.right },
      controls,
      help: { top: helpBox.top, bottom: helpBox.bottom, left: helpBox.left, right: helpBox.right },
      controlEdgeProbes,
    };
  }, controlTestId);
}

for (const [width, fixture] of [[390, 'recording-editor-real-host.html'], [1440, 'recording-selected-host.html']] as const) {
  test(`keeps destination labels and controls compact at ${width}px`, async ({ page }, testInfo) => {
    await page.setViewportSize({ width, height: 900 });
    await page.route('**/api/**', route => route.abort());
    await page.goto(`/tests/e2e/fixtures/${fixture}`);
    await page.getByTestId('group-advanced').scrollIntoViewIfNeeded();

    for (const [helpKey, controlTestId] of fields) {
      await page.locator(`[data-parameter-help="${helpKey}"]`).scrollIntoViewIfNeeded();
      await page.locator(`[data-testid="${controlTestId}"]`).scrollIntoViewIfNeeded();
      const geometry = await fieldGeometry(page, helpKey, controlTestId);
      const firstControl = geometry.controls[0];
      const labelControlGap = firstControl.top - geometry.label.bottom;
      const labelHelpOffset = Math.abs((geometry.label.top + geometry.label.bottom) / 2 - (geometry.help.top + geometry.help.bottom) / 2);
      expect(labelControlGap, `${helpKey}: ${JSON.stringify(geometry)}`).toBeGreaterThanOrEqual(0);
      expect(labelControlGap, `${helpKey}: ${JSON.stringify(geometry)}`).toBeLessThanOrEqual(12);
      expect(labelHelpOffset, `${helpKey}: ${JSON.stringify(geometry)}`).toBeLessThanOrEqual(5);
      expect(geometry.help.right - geometry.help.left, helpKey).toBeGreaterThanOrEqual(44);
      expect(geometry.help.bottom - geometry.help.top, helpKey).toBeGreaterThanOrEqual(44);
      expect(geometry.controlEdgeProbes.every(probe => probe.ok), `${helpKey}: help hit area intercepts control edge: ${JSON.stringify(geometry.controlEdgeProbes)}`).toBe(true);
      for (let index = 1; index < geometry.controls.length; index += 1) {
        expect(geometry.controls[index].top, `${helpKey}: controls overlap`).toBeGreaterThanOrEqual(geometry.controls[index - 1].bottom);
      }
    }
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
    await page.screenshot({ path: testInfo.outputPath(`destination-${width}.png`), fullPage: true });
  });
}
