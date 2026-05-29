import { expectTypeOf } from 'vitest';
import type { Settings, SettingsConnector, TimeseriesSettings, SchedulerSettings, ModbusShareSettings, GeneralSettings } from './types';

/**
 * @file types-settings.test-d.ts
 * @description 驗證 Settings 及其子結構型別契約。
 */

// 1. 斷言 SettingsConnector 結構
expectTypeOf<SettingsConnector['kind']>().toEqualTypeOf<'sqlite' | 'postgres' | 'mysql' | 'sqlserver'>();
expectTypeOf<SettingsConnector['status']>().toEqualTypeOf<'unknown' | 'testing' | 'ready' | 'unreachable' | 'auth_failed'>();
expectTypeOf<SettingsConnector['id']>().toEqualTypeOf<string>();

// 2. 斷言 TimeseriesSettings 結構
expectTypeOf<TimeseriesSettings['write_precision']>().toEqualTypeOf<'second' | 'millisecond'>();
expectTypeOf<TimeseriesSettings['partition_interval']>().toEqualTypeOf<'daily' | 'weekly' | 'monthly'>();

// 3. 斷言 SchedulerSettings 結構
expectTypeOf<SchedulerSettings['auto_start']>().toEqualTypeOf<boolean>();
expectTypeOf<SchedulerSettings['default_interval_ms']>().toEqualTypeOf<number>();

// 4. 斷言 ModbusShareSettings 結構
expectTypeOf<ModbusShareSettings['base_register']>().toEqualTypeOf<number>();
expectTypeOf<ModbusShareSettings['enabled']>().toEqualTypeOf<boolean>();

// 5. 斷言 GeneralSettings 結構
expectTypeOf<GeneralSettings['theme']>().toEqualTypeOf<'dark' | 'light' | 'auto'>();
expectTypeOf<GeneralSettings['locale']>().toEqualTypeOf<'zh-TW' | 'en'>();
expectTypeOf<GeneralSettings['api_version']>().toEqualTypeOf<'v1' | 'v2'>();

// 6. 斷言 Settings 包含子屬性
expectTypeOf<Settings['connectors']>().toEqualTypeOf<SettingsConnector[]>();
expectTypeOf<Settings['timeseries']>().toEqualTypeOf<TimeseriesSettings>();
expectTypeOf<Settings['scheduler']>().toEqualTypeOf<SchedulerSettings>();
expectTypeOf<Settings['modbus_share']>().toEqualTypeOf<ModbusShareSettings>();
expectTypeOf<Settings['general']>().toEqualTypeOf<GeneralSettings>();
