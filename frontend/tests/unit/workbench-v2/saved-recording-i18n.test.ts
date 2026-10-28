import { describe, expect, it } from 'vitest';
import i18n from '@/i18n/config';
import en from '@/i18n/locales/en/recording-lifecycle.json';
import zh from '@/i18n/locales/zh-TW/recording-lifecycle.json';
import enDialog from '@/i18n/locales/en/recording-lifecycle-dialog.json';
import zhDialog from '@/i18n/locales/zh-TW/recording-lifecycle-dialog.json';
function leaves(tree: Record<string, unknown>, prefix = ''): string[] {
  return Object.entries(tree).flatMap(([key, value]) => value && typeof value === 'object' && !Array.isArray(value)
    ? leaves(value as Record<string, unknown>, `${prefix}${key}.`) : [`${prefix}${key}`]);
}
describe('saved recording localized action and review copy', () => {
  it('has matching complete English and Traditional Chinese modules', () => {
    expect(leaves(en).sort()).toEqual(leaves(zh).sort());
    expect(leaves(enDialog).sort()).toEqual(leaves(zhDialog).sort());
  });
  it('merges action and dialog modules into the same runtime namespace', () => {
    for (const language of ['en', 'zh-TW']) {
      for (const key of [...leaves(en), ...leaves(enDialog)]) {
        expect(i18n.exists(key, { lng: language, ns: 'recording-lifecycle' }), `${language}:${key}`).toBe(true);
      }
    }
    expect(i18n.t('action.stop', { lng: 'zh-TW', ns: 'recording-lifecycle' })).toBe('停止新記錄');
    expect(i18n.t('remove.confirm', { lng: 'zh-TW', ns: 'recording-lifecycle' })).toBe('移除設定並保留歷史');
  });
});
