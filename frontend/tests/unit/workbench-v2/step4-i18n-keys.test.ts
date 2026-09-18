import { describe, expect, it } from 'vitest';
import { readdirSync, readFileSync } from 'node:fs';
import { join } from 'node:path';
import i18n from '../../../src/i18n/config';

const CJK = /[\u4e00-\u9fff]/;

const STEP4_DIR = join(process.cwd(), 'src/features/datalink/workbench-v2/steps/step4');

/** Collect every literal step4.* translation key used via t('...') in the step4 sources. */
function usedStep4Keys(): string[] {
  const keys = new Set<string>();
  for (const file of readdirSync(STEP4_DIR).filter((name) => name.endsWith('.tsx') || name.endsWith('.ts'))) {
    const source = readFileSync(join(STEP4_DIR, file), 'utf8');
    for (const match of source.matchAll(/\bt\('([^']+)'/g)) {
      if (match[1].startsWith('step4.')) keys.add(match[1]);
    }
  }
  return [...keys];
}

describe('step4 keys resolve through the real i18n instance', () => {
  it.each(['zh-TW', 'en'])('resolves step4.conflict_tooltip in %s', async (lng) => {
    await i18n.changeLanguage(lng);
    const tooltip = i18n.t('step4.conflict_tooltip', { ns: 'workbench-v2' });
    expect(tooltip).not.toBe('step4.conflict_tooltip');
    expect(String(tooltip).length).toBeGreaterThan(0);
  });

  it('keeps the English UI free of leaked Chinese copy', async () => {
    await i18n.changeLanguage('en');
    for (const key of ['step4.conflict_tooltip', 'step4.target_remote_conflict']) {
      const text = i18n.t(key, { ns: 'workbench-v2' });
      expect(text, key).not.toBe(key);
      expect(String(text), key).not.toMatch(CJK);
    }
  });

  it('resolves every step4 key used by the step4 directory in both languages', async () => {
    const keys = usedStep4Keys();
    expect(keys.length).toBeGreaterThan(0);
    for (const lng of ['zh-TW', 'en'] as const) {
      await i18n.changeLanguage(lng);
      for (const key of keys) {
        expect(i18n.t(key, { ns: 'workbench-v2' }), `${lng}: ${key}`).not.toBe(key);
      }
    }
  });
});
