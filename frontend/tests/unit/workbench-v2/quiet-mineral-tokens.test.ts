import { readFileSync } from 'node:fs';
import { describe, expect, it } from 'vitest';
const css = readFileSync('src/index.css', 'utf8');
const expected = {
  ground: ['#F7F8F5', '#191A19'], surface: ['#FFFFFF', '#202220'], raised: ['#FFFFFF', '#242624'], active: ['#E2F1EA', '#26473E'],
  text: ['#203332', '#EFEFEB'], muted: ['#586762', '#B8BEB9'], faint: ['#586762', '#B8BEB9'], border: ['#C8D1C9', '#586762'],
  control: ['#586762', '#929B94'], focus: ['#1D6D63', '#83C8B7'], accent: ['#1D6D63', '#83C8B7'], primary: ['#1D6D63', '#83C8B7'],
  'on-primary': ['#FFFFFF', '#191A19'], 'primary-hover': ['#175A52', '#BDE0D3'],
  'success-fg': ['#166534', '#BBF7D0'], 'success-bg': ['#F0FDF4', '#14532D'], 'success-border': ['#15803D', '#4ADE80'],
  'warning-fg': ['#854D0E', '#FEF08A'], 'warning-bg': ['#FEFCE8', '#713F12'], 'warning-border': ['#A16207', '#FACC15'],
  'danger-fg': ['#991B1B', '#FECACA'], 'danger-bg': ['#FEF2F2', '#7F1D1D'], 'danger-border': ['#B91C1C', '#F87171'],
  'unknown-fg': ['#394A43', '#EFEFEB'], 'unknown-bg': ['#E0E5DF', '#394A43'], 'unknown-border': ['#586762', '#929B94'],
};
function theme(mode: 'light' | 'dark') {
  const blocks = [...css.matchAll(/([^{}]+)\{([^{}]+)\}/g)].filter(([, selector]) => selector.includes(':root') || selector.trim() === `html.${mode}`);
  const values = Object.fromEntries(blocks.flatMap(([, , body]) => [...body.matchAll(/(--gw-[\w-]+)\s*:\s*([^;]+);/g)].map(([, key, value]) => [key, value.trim()])));
  const resolve = (key: string): string => values[key]?.startsWith('var(') ? resolve(values[key].slice(4, -1)) : values[key];
  return resolve;
}
function luminance(hex: string) {
  const rgb = hex.slice(1).match(/../g)!.map(value => parseInt(value, 16) / 255).map(value => value <= .04045 ? value / 12.92 : ((value + .055) / 1.055) ** 2.4);
  return rgb[0] * .2126 + rgb[1] * .7152 + rgb[2] * .0722;
}
const ratio = (a: string, b: string) => (Math.max(luminance(a), luminance(b)) + .05) / (Math.min(luminance(a), luminance(b)) + .05);
describe('DESIGN.md Quiet mineral semantic contract (not a browser render test)', () => {
  it.each(['light', 'dark'] as const)('%s resolves the exact agreed semantic palette and safe contrast pairs', mode => {
    const value = theme(mode), index = mode === 'light' ? 0 : 1;
    for (const [role, colors] of Object.entries(expected)) expect(value(`--gw-${role}`)?.toUpperCase(), role).toBe(colors[index]);
    for (const bg of ['ground', 'surface', 'raised']) for (const fg of ['text', 'muted', 'faint']) expect(ratio(value(`--gw-${fg}`), value(`--gw-${bg}`)), `${fg}/${bg}`).toBeGreaterThanOrEqual(4.5);
    for (const bg of ['ground', 'surface', 'raised']) for (const fg of ['focus', 'control']) expect(ratio(value(`--gw-${fg}`), value(`--gw-${bg}`)), `${fg}/${bg}`).toBeGreaterThanOrEqual(3);
    for (const kind of ['success', 'warning', 'danger', 'unknown']) expect(ratio(value(`--gw-${kind}-fg`), value(`--gw-${kind}-bg`)), kind).toBeGreaterThanOrEqual(4.5);
    expect(ratio(value('--gw-muted'), value('--gw-unknown-bg'))).toBeGreaterThanOrEqual(4.5);
    expect(ratio(value('--gw-on-primary'), value('--gw-primary'))).toBeGreaterThanOrEqual(4.5);
    expect(ratio(value('--gw-on-primary'), value('--gw-primary-hover'))).toBeGreaterThanOrEqual(4.5);
    for (const fg of ['accent', 'text', 'muted']) expect(ratio(value(`--gw-${fg}`), value('--gw-active')), `${fg}/active`).toBeGreaterThanOrEqual(4.5);
  });
  it('keeps fonts/CJK hierarchy and all new styles scoped, with mobile/reduced-motion rules', () => {
    expect(css).toContain('--gw-font-sans: Inter, "Noto Sans TC", "PingFang TC", "Microsoft JhengHei", system-ui, sans-serif');
    const home = readFileSync('src/features/datalink/workbench-v2/recording-home/recording-home.css', 'utf8');
    expect(home).not.toMatch(/#[0-9a-f]{3,8}\b|(?:slate|blue|navy)-/i);
    expect(home).toContain('prefers-reduced-motion: reduce'); expect(home).toContain('max-width: 639px');
    expect(home).toContain('var(--gw-control-mobile)'); expect(home).toContain(':focus-visible');
    expect(home).toContain('overflow-wrap: anywhere'); expect(home).not.toContain('75vw');
    expect(home).not.toMatch(/url\(|@import|@font-face/); // no new network-dependent font or asset
    expect(css.match(/fonts\.googleapis\.com/g)).toHaveLength(1); // pre-existing import only
    expect(home).toContain('forced-colors: active');
  });
});
