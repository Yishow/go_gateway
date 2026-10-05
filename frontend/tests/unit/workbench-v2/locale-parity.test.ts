import { describe, expect, it } from 'vitest';
import enWorkbench from '../../../src/i18n/locales/en/workbench-v2.json';
import zhWorkbench from '../../../src/i18n/locales/zh-TW/workbench-v2.json';
import enRuntime from '../../../src/i18n/locales/en/runtime-dashboard.json';
import zhRuntime from '../../../src/i18n/locales/zh-TW/runtime-dashboard.json';
import enMappingErrors from '../../../src/i18n/locales/en/mapping-errors.json';
import zhMappingErrors from '../../../src/i18n/locales/zh-TW/mapping-errors.json';
import i18n from '../../../src/i18n/config';
import { BACKEND_ERROR_CODES } from '../../../src/utils/typedErrors';

const MAPPING_ERROR_CODES = new Set<string>(['workspace_mapping_conflict', 'workspace_mapping_not_found', 'workspace_mapping_save_failed']);

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

  it('ensures configured mapping-errors en and zh-TW have identical keys', () => {
    expect(getDeepKeys(enMappingErrors).sort()).toEqual(getDeepKeys(zhMappingErrors).sort());
    expect(i18n.getResourceBundle('en', 'mapping-errors')).toEqual(enMappingErrors);
    expect(i18n.getResourceBundle('zh-TW', 'mapping-errors')).toEqual(zhMappingErrors);
  });

  it('ensures all required typed error codes are translated', () => {
    const requiredCodes = [...BACKEND_ERROR_CODES, 'generic_failure'];

    const enErrors = (enWorkbench as Record<string, unknown>).errors as Record<string, string>;
    const zhErrors = (zhWorkbench as Record<string, unknown>).errors as Record<string, string>;

    for (const code of requiredCodes) {
      const mappingError = MAPPING_ERROR_CODES.has(code);
      const namespace = mappingError ? 'mapping-errors' : 'workbench-v2';
      const key = mappingError ? code : `errors.${code}`;
      const expectedEn = mappingError ? (enMappingErrors as Record<string, string>)[code] : enErrors[code];
      const expectedZh = mappingError ? (zhMappingErrors as Record<string, string>)[code] : zhErrors[code];
      expect(expectedEn, `EN translation for ${code} in ${namespace}`).toBeTruthy();
      expect(expectedZh, `ZH-TW translation for ${code} in ${namespace}`).toBeTruthy();
      expect(i18n.getResource('en', namespace, key)).toBe(expectedEn);
      expect(i18n.getResource('zh-TW', namespace, key)).toBe(expectedZh);
      expect(i18n.getFixedT('en', namespace)(key)).toBe(expectedEn);
      expect(i18n.getFixedT('zh-TW', namespace)(key)).toBe(expectedZh);
    }
  });
});
