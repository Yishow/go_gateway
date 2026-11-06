import { expect, type Locator, type Page } from '@playwright/test';

export async function tracks(page: Page) {
  return page.locator('.gw-memory-grid').evaluate(node => getComputedStyle(node).gridTemplateColumns.split(' ').filter(Boolean).length);
}
export async function box(locator: Locator) {
  const bounds = await locator.boundingBox();
  expect(bounds).not.toBeNull();
  return bounds!;
}
export async function noHorizontalOverflow(page: Page) {
  expect(await page.getByRole('dialog').evaluate(root => {
    const bounds = root.getBoundingClientRect();
    return [...root.querySelectorAll<HTMLElement>('*')].flatMap(node => {
      if (!node.checkVisibility()) return [];
      const rect = node.getBoundingClientRect(); let { left, right } = rect;
      for (let parent = node.parentElement; parent && parent !== root; parent = parent.parentElement) {
        if (getComputedStyle(parent).overflowX === 'visible') continue;
        const clip = parent.getBoundingClientRect(); left = Math.max(left, clip.left); right = Math.min(right, clip.right);
      }
      return rect.height > 0 && right > left && (left < bounds.left - 1 || right > bounds.right + 1) ? [node.className] : [];
    });
  })).toEqual([]);
}

/** Browser-resolved colors rendered into canvas, then alpha-composited over actual ancestor backgrounds.
 * No guessed hex/token values. Unsupported images/filters fail explicitly instead of claiming a measurement.
 * This supplements screenshots; it is not an antialiased-glyph or native range-thumb pixel audit.
 */
export async function contrastAudit(root: Locator, focusOnly = false) {
  return root.evaluate((rootNode, onlyFocus) => {
    type RGBA = [number, number, number, number];
    const canvas = document.createElement('canvas'); canvas.width = canvas.height = 1;
    const ctx = canvas.getContext('2d', { willReadFrequently: true })!;
    const rgba = (color: string): RGBA => {
      if (!CSS.supports('color', color)) throw new Error(`Unsupported computed color: ${color}`);
      ctx.clearRect(0, 0, 1, 1); ctx.fillStyle = color; ctx.fillRect(0, 0, 1, 1);
      const pixel = ctx.getImageData(0, 0, 1, 1).data;
      return [pixel[0], pixel[1], pixel[2], pixel[3] / 255];
    };
    const over = (fg: RGBA, bg: RGBA): RGBA => {
      const alpha = fg[3] + bg[3] * (1 - fg[3]);
      return [0, 1, 2].map(index => (fg[index] * fg[3] + bg[index] * bg[3] * (1 - fg[3])) / (alpha || 1)).concat(alpha) as RGBA;
    };
    const background = (node: Element | null): RGBA => {
      if (!node) return [255, 255, 255, 1];
      const style = getComputedStyle(node);
      if (style.backgroundImage !== 'none' || style.filter !== 'none' || Number(style.opacity) !== 1) {
        throw new Error(`Unsupported painted backdrop at ${node.tagName}.${node.className}`);
      }
      return over(rgba(style.backgroundColor), background(node.parentElement));
    };
    const luminance = (color: RGBA) => color.slice(0, 3).map(value => {
      const channel = value / 255; return channel <= .04045 ? channel / 12.92 : ((channel + .055) / 1.055) ** 2.4;
    }).reduce((sum, value, index) => sum + value * [.2126, .7152, .0722][index], 0);
    const ratio = (a: RGBA, b: RGBA) => (Math.max(luminance(a), luminance(b)) + .05) / (Math.min(luminance(a), luminance(b)) + .05);
    const samples: { kind: string; text: string; ratio: number; minimum: number }[] = [];
    const unsupported: string[] = [];
    const visible = (node: HTMLElement) => {
      if (!node.checkVisibility()) return false;
      const rect = node.getBoundingClientRect();
      if (!rect.width || !rect.height) return false;
      let top = Math.max(0, rect.top), bottom = Math.min(innerHeight, rect.bottom);
      let left = Math.max(0, rect.left), right = Math.min(innerWidth, rect.right);
      for (let parent = node.parentElement; parent; parent = parent.parentElement) {
        const style = getComputedStyle(parent), clip = parent.getBoundingClientRect();
        if (style.overflowX !== 'visible') { left = Math.max(left, clip.left); right = Math.min(right, clip.right); }
        if (style.overflowY !== 'visible') { top = Math.max(top, clip.top); bottom = Math.min(bottom, clip.bottom); }
      }
      return right > left && bottom > top;
    };
    const nodes = onlyFocus ? [document.activeElement as HTMLElement] : [...rootNode.querySelectorAll<HTMLElement>('*')];
    for (const node of nodes) {
      if (!node || !rootNode.contains(node) || node.closest(':disabled') || !visible(node)) continue;
      const style = getComputedStyle(node), text = [...node.childNodes].filter(child => child.nodeType === Node.TEXT_NODE).map(child => child.textContent).join('').trim();
      const textSample = !onlyFocus && Boolean(text) && style.clipPath === 'none';
      const controlSample = !onlyFocus && node.matches('input:not([type=range]), select');
      // A decorative hatch swatch has no contrast sample; do not classify it as unmeasurable text.
      if (!onlyFocus && !textSample && !controlSample) continue;
      try {
        const bg = background(node);
        if (textSample) samples.push({ kind: 'text', text: text.slice(0, 90), ratio: ratio(over(rgba(style.color), bg), bg), minimum: 4.5 });
        if (onlyFocus) {
          if (parseFloat(style.outlineWidth) < 2 || style.outlineStyle === 'none') throw new Error('Focused control has no visible 2px outline');
          const adjacent = parseFloat(style.outlineOffset) < 0 ? bg : background(node.parentElement);
          samples.push({ kind: 'focus', text: node.getAttribute('aria-label') ?? node.tagName, ratio: ratio(over(rgba(style.outlineColor), adjacent), adjacent), minimum: 3 });
        } else if (controlSample) {
          samples.push({ kind: 'control-boundary', text: node.getAttribute('aria-label') ?? node.tagName,
            ratio: ratio(over(rgba(style.borderTopColor), bg), background(node.parentElement)), minimum: 3 });
        }
      } catch (error) { unsupported.push(String(error)); }
    }
    return { samples, unsupported };
  }, focusOnly);
}
export function assertContrast(audit: Awaited<ReturnType<typeof contrastAudit>>) {
  expect(audit.unsupported, 'Unmeasurable backgrounds need explicit pixel/manual evidence, never a claimed pass').toEqual([]);
  expect(audit.samples.length).toBeGreaterThan(0);
  expect(audit.samples.filter(sample => sample.ratio < sample.minimum), 'Measured adjacent-background contrast').toEqual([]);
}
export async function phoneSizes(page: Page) {
  const measurements = await page.getByRole('dialog').evaluate(root => {
    const targets = [...root.querySelectorAll<HTMLElement>('button, input, select, summary')].filter(node => node.checkVisibility());
    const labels = [...root.querySelectorAll<HTMLElement>('.gw-memory-tabs > button, .gw-memory-layouts > button, .gw-memory-toolbar label, .gw-memory-toolbar button, .gw-memory-nav p, .gw-memory-jumps > button, .gw-memory-address, .gw-memory-type, .gw-memory-reading, .gw-memory-range-nav strong, .gw-memory-chapter-summary, .gw-memory-summary-field, .gw-memory-detail-facts')].filter(node => node.checkVisibility());
    return { targets: targets.map(node => ({ name: node.getAttribute('data-testid') ?? node.tagName, height: node.getBoundingClientRect().height, width: node.getBoundingClientRect().width })),
      labels: labels.map(node => ({ name: node.className, size: parseFloat(getComputedStyle(node).fontSize) })) };
  });
  expect(measurements.targets.filter(item => item.height < 43.5 || item.width < 43.5)).toEqual([]);
  expect(measurements.labels.filter(item => item.size < 14)).toEqual([]);
  return measurements;
}
