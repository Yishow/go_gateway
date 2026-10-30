import { describe, expect, it } from 'vitest';
import i18n from '@/i18n/config';
describe('actual registered global theme namespace', () => {
  it.each(['en', 'zh-TW'])('%s resolves every persistence label without a missing-key fallback', language => {
    const t = i18n.getFixedT(language, 'theme');
    for (const key of ['saving', 'save_failed', 'retry', 'unsaved']) {
      expect(i18n.exists(key, { lng: language, ns: 'theme' })).toBe(true);
      expect(t(key)).not.toBe(key);
    }
  });
});
