import { describe, expect, it } from 'vitest';
import enWorkbench from '../../../src/i18n/locales/en/workbench-v2.json';
import enSQLServerWorkbench from '../../../src/i18n/locales/en/workbench-v2-circuit-sqlserver.json';
import zhWorkbench from '../../../src/i18n/locales/zh-TW/workbench-v2.json';
import zhSQLServerWorkbench from '../../../src/i18n/locales/zh-TW/workbench-v2-circuit-sqlserver.json';
import enRuntime from '../../../src/i18n/locales/en/runtime-dashboard.json';
import zhRuntime from '../../../src/i18n/locales/zh-TW/runtime-dashboard.json';
import enLogs from '../../../src/i18n/locales/en/runtime-logs.json';
import zhLogs from '../../../src/i18n/locales/zh-TW/runtime-logs.json';
import enMappingErrors from '../../../src/i18n/locales/en/mapping-errors.json';
import zhMappingErrors from '../../../src/i18n/locales/zh-TW/mapping-errors.json';
import i18n from '../../../src/i18n/config';
import { mergeTranslations } from '../../../src/i18n/mergeTranslations';
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
  it('registers complete en/zh-TW runtime-logs namespaces', () => {
    expect(getDeepKeys(enLogs).sort()).toEqual(getDeepKeys(zhLogs).sort());
    expect(i18n.getResourceBundle('en', 'runtime-logs')).toEqual(enLogs);
    expect(i18n.getResourceBundle('zh-TW', 'runtime-logs')).toEqual(zhLogs);
  });

  it('ensures workbench-v2 en and zh-TW have identical keys', () => {
    const enKeys = getDeepKeys(enWorkbench).sort();
    const zhKeys = getDeepKeys(zhWorkbench).sort();

    const missingInZh = enKeys.filter((k) => !zhKeys.includes(k));
    const missingInEn = zhKeys.filter((k) => !enKeys.includes(k));

    expect(missingInZh, 'Keys present in EN but missing in ZH-TW').toEqual([]);
    expect(missingInEn, 'Keys present in ZH-TW but missing in EN').toEqual([]);

    const enSqlKeys = getDeepKeys(enSQLServerWorkbench).sort();
    const zhSqlKeys = getDeepKeys(zhSQLServerWorkbench).sort();
    expect(enSqlKeys.filter((k) => !zhSqlKeys.includes(k)), 'SQL Server keys present in EN but missing in ZH-TW').toEqual([]);
    expect(zhSqlKeys.filter((k) => !enSqlKeys.includes(k)), 'SQL Server keys present in ZH-TW but missing in EN').toEqual([]);
  });

  it('ensures runtime-dashboard en and zh-TW have identical keys', () => {
    const enKeys = getDeepKeys(enRuntime).sort();
    const zhKeys = getDeepKeys(zhRuntime).sort();

    const missingInZh = enKeys.filter((k) => !zhKeys.includes(k));
    const missingInEn = zhKeys.filter((k) => !enKeys.includes(k));

    expect(missingInZh, 'Keys present in EN but missing in ZH-TW').toEqual([]);
    expect(missingInEn, 'Keys present in ZH-TW but missing in EN').toEqual([]);
  });

  it('ensures workbench-v2 merges SQL Server translations with distinct visible labels in en and zh-TW', () => {
    expect(i18n.getResource('en', 'workbench-v2', 'step4.basic.circuit.storage_native_monthly_sqlserver'))
      .toBe('Native monthly partitions (SQL Server)');
    expect(i18n.getResource('zh-TW', 'workbench-v2', 'step4.basic.circuit.storage_native_monthly_sqlserver'))
      .toBe('原生月分割區（SQL Server）');
    expect(i18n.getResource('en', 'workbench-v2', 'step4.basic.circuit.storage_native_monthly'))
      .toBe('Native monthly partitions (MySQL / MariaDB)');
    expect(i18n.getResource('zh-TW', 'workbench-v2', 'step4.basic.circuit.storage_native_monthly'))
      .toBe('原生月分割區（MySQL／MariaDB）');
    expect(i18n.getResource('en', 'workbench-v2', 'step4.group.test_write.blocked_sqlserver_capability'))
      .toBe('The saved SQL Server has not verified test-write support.');
    expect(i18n.getResource('zh-TW', 'workbench-v2', 'step4.group.test_write.blocked_sqlserver_capability'))
      .toBe('已儲存的 SQL Server 伺服器尚未驗證試寫能力。');
    expect(i18n.getResource('en', 'workbench-v2', 'step4.group.test_write.blocked_capability'))
      .toBe('The saved MySQL server has not verified test-write support.');
    expect(i18n.getResource('zh-TW', 'workbench-v2', 'step4.group.test_write.blocked_capability'))
      .toBe('已儲存的 MySQL 伺服器尚未驗證試寫能力。');
  });

  it('ensures configured mapping-errors en and zh-TW have identical keys', () => {
    expect(getDeepKeys(enMappingErrors).sort()).toEqual(getDeepKeys(zhMappingErrors).sort());
    expect(i18n.getResourceBundle('en', 'mapping-errors')).toEqual(enMappingErrors);
    expect(i18n.getResourceBundle('zh-TW', 'mapping-errors')).toEqual(zhMappingErrors);
  });

  it('ensures all required typed error codes are translated', () => {
    const requiredCodes = [...BACKEND_ERROR_CODES, 'generic_failure'];

    const enMergedWorkbench = mergeTranslations(enWorkbench, enSQLServerWorkbench) as Record<string, unknown>;
    const zhMergedWorkbench = mergeTranslations(zhWorkbench, zhSQLServerWorkbench) as Record<string, unknown>;
    const enErrors = enMergedWorkbench.errors as Record<string, string>;
    const zhErrors = zhMergedWorkbench.errors as Record<string, string>;

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
