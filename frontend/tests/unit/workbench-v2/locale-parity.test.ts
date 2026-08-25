import { describe, expect, it } from 'vitest';
import enWorkbench from '../../../src/i18n/locales/en/workbench-v2.json';
import zhWorkbench from '../../../src/i18n/locales/zh-TW/workbench-v2.json';
import enRuntime from '../../../src/i18n/locales/en/runtime-dashboard.json';
import zhRuntime from '../../../src/i18n/locales/zh-TW/runtime-dashboard.json';
import { BACKEND_ERROR_CODES } from '../../../src/utils/typedErrors';

function getDeepKeys(obj: Record<string, unknown>, prefix = ''): string[] {
  return Object.keys(obj).flatMap((key) => {
    const val = obj[key];
    const fullKey = prefix ? `${prefix}.${key}` : key;
    if (typeof val === 'object' && val !== null && !Array.isArray(val)) {
      return getDeepKeys(val as Record<string, unknown>, fullKey);
    }
    return [fullKey];
  });
}

describe('i18n Locale Parity and Namespace Coverage', () => {
  it('ensures workbench-v2 en and zh-TW have identical keys', () => {
    const enKeys = getDeepKeys(enWorkbench).sort();
    const zhKeys = getDeepKeys(zhWorkbench).sort();

    const missingInZh = enKeys.filter((k) => !zhKeys.includes(k));
    const missingInEn = zhKeys.filter((k) => !enKeys.includes(k));

    expect(missingInZh, 'Keys present in EN but missing in ZH-TW').toEqual([]);
    expect(missingInEn, 'Keys present in ZH-TW but missing in EN').toEqual([]);
  });

  it('ensures runtime-dashboard en and zh-TW have identical keys', () => {
    const enKeys = getDeepKeys(enRuntime).sort();
    const zhKeys = getDeepKeys(zhRuntime).sort();

    const missingInZh = enKeys.filter((k) => !zhKeys.includes(k));
    const missingInEn = zhKeys.filter((k) => !enKeys.includes(k));

    expect(missingInZh, 'Keys present in EN but missing in ZH-TW').toEqual([]);
    expect(missingInEn, 'Keys present in ZH-TW but missing in EN').toEqual([]);
  });

  it('ensures all required typed error codes are translated', () => {
    const requiredCodes = [...BACKEND_ERROR_CODES, 'generic_failure'];

    const enErrors = (enWorkbench as Record<string, unknown>).errors as Record<string, string>;
    const zhErrors = (zhWorkbench as Record<string, unknown>).errors as Record<string, string>;

    for (const code of requiredCodes) {
      expect(enErrors[code], `EN translation for ${code}`).toBeTruthy();
      expect(zhErrors[code], `ZH-TW translation for ${code}`).toBeTruthy();
    }
  });
});
