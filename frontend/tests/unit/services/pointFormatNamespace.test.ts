import { describe, expect, it } from 'vitest';
import i18n from '@/i18n/config';
import { pointFormatErrorKey, safePointFormatError } from '@/services/pointFormatErrors';

// Import the production registry, not an isolated resource object or a mocked translator.
describe('registered point format operator messages', () => {
  it.each(['en', 'zh-TW'])('resolves safe errors and pending state in %s', async (language) => {
    await i18n.changeLanguage(language);
    expect(i18n.hasResourceBundle(language, 'point-format')).toBe(true);
    for (const [status, code] of [
      [403, 'point_format_scope'], [409, 'point_format_conflict'],
      [422, 'point_format_invalid'], [503, 'point_format_unavailable'],
      [404, 'point_format_operation_not_found'],
    ] as const) {
      const error = safePointFormatError({ response: { status, data: { error: { code, message: 'SECRET' } } }, config: { password: 'SECRET' } }, 'original-request');
      const key = pointFormatErrorKey(error);
      const translated = i18n.t(key);
      expect(translated).not.toBe(key);
      expect(translated).not.toBe(`errors.${code}`);
      expect(translated).not.toContain('SECRET');
      expect(error.requestId).toBe('original-request');
      if (status === 503 || status === 404) expect(error.outcome).toBe('unknown');
    }
    const unknown = safePointFormatError(new Error('SECRET'));
    expect(i18n.t(pointFormatErrorKey(unknown))).not.toContain('SECRET');
    const pending = i18n.t('point-format:savedPending');
    expect(pending).not.toBe('savedPending');
    expect(pending).toMatch(language === 'zh-TW' ? /待重新套用設備設定/ : /still needs.*applied/i);
    expect(i18n.t('point-format:policy.legacy_mc3e')).not.toBe('policy.legacy_mc3e');
    expect(i18n.t('point-format:policy.legacy_fatek_components')).not.toBe('policy.legacy_fatek_components');
  });
});
