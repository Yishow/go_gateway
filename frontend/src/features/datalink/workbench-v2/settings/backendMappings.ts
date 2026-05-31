import type {
  CreateDatabaseConnectorRequest,
  DatabaseConnector,
  SettingItem,
  UpdateDatabaseConnectorRequest,
} from '../../../../types/datalink';
import type { Settings, SettingsConnector } from '../state/types';
import { DEFAULT_SETTINGS } from '../state/settingsDefaults';

type PersistableSettingEntry = {
  key: string;
  value: unknown;
};

function stringValue(value: unknown, fallback: string): string {
  return typeof value === 'string' ? value : fallback;
}

function numberValue(value: unknown, fallback: number): number {
  return typeof value === 'number' && Number.isFinite(value) ? value : fallback;
}

function booleanValue(value: unknown, fallback: boolean): boolean {
  return typeof value === 'boolean' ? value : fallback;
}

function connectorConfigValue(config: Record<string, unknown>, ...keys: string[]): string {
  for (const key of keys) {
    const value = config[key];
    if (typeof value === 'string' && value.trim() !== '') {
      return value;
    }
    if (typeof value === 'number') {
      return String(value);
    }
  }
  return '';
}

export function mapSettingItemsToWorkbenchSettings(items: SettingItem[]): Settings {
  const values = new Map(items.map((item) => [item.key, item.value]));

  return {
    ...DEFAULT_SETTINGS,
    connectors: DEFAULT_SETTINGS.connectors,
    timeseries: {
      write_precision: stringValue(
        values.get('write_precision'),
        DEFAULT_SETTINGS.timeseries.write_precision,
      ) as Settings['timeseries']['write_precision'],
      partition_interval: stringValue(
        values.get('partition_interval'),
        DEFAULT_SETTINGS.timeseries.partition_interval,
      ) as Settings['timeseries']['partition_interval'],
      batch_size: numberValue(values.get('batch_size'), DEFAULT_SETTINGS.timeseries.batch_size),
      retention_days: numberValue(
        values.get('retention_days'),
        DEFAULT_SETTINGS.timeseries.retention_days,
      ),
    },
    scheduler: {
      default_interval_ms: numberValue(
        values.get('default_interval_ms'),
        DEFAULT_SETTINGS.scheduler.default_interval_ms,
      ),
      default_retry_count: numberValue(
        values.get('default_retry_count'),
        DEFAULT_SETTINGS.scheduler.default_retry_count,
      ),
      default_retry_delay_ms: numberValue(
        values.get('default_retry_delay'),
        DEFAULT_SETTINGS.scheduler.default_retry_delay_ms,
      ),
      breaker_threshold: numberValue(
        values.get('breaker_threshold'),
        DEFAULT_SETTINGS.scheduler.breaker_threshold,
      ),
      auto_start: booleanValue(values.get('auto_start'), DEFAULT_SETTINGS.scheduler.auto_start),
    },
    modbus_share: {
      enabled: booleanValue(
        values.get('modbus_share_enabled'),
        DEFAULT_SETTINGS.modbus_share.enabled,
      ),
      bind_address: stringValue(
        values.get('modbus_share_bind_address'),
        DEFAULT_SETTINGS.modbus_share.bind_address,
      ),
      port: numberValue(values.get('modbus_share_port'), DEFAULT_SETTINGS.modbus_share.port),
      slave_id: numberValue(
        values.get('modbus_share_slave_id'),
        DEFAULT_SETTINGS.modbus_share.slave_id,
      ),
      base_register: numberValue(
        values.get('modbus_share_base_register'),
        DEFAULT_SETTINGS.modbus_share.base_register,
      ),
    },
    general: {
      theme: stringValue(values.get('theme'), DEFAULT_SETTINGS.general.theme) as Settings['general']['theme'],
      locale: stringValue(values.get('locale'), DEFAULT_SETTINGS.general.locale) as Settings['general']['locale'],
      addr_format: stringValue(
        values.get('addr_format'),
        DEFAULT_SETTINGS.general.addr_format,
      ) as Settings['general']['addr_format'],
      api_base: stringValue(values.get('api_base'), DEFAULT_SETTINGS.general.api_base),
      api_version: stringValue(
        values.get('api_version'),
        DEFAULT_SETTINGS.general.api_version,
      ) as Settings['general']['api_version'],
      timeout_seconds: numberValue(
        values.get('timeout_seconds'),
        DEFAULT_SETTINGS.general.timeout_seconds,
      ),
      log_level: stringValue(
        values.get('log_level'),
        DEFAULT_SETTINGS.general.log_level,
      ) as Settings['general']['log_level'],
      sse_heartbeat_seconds: numberValue(
        values.get('sse_heartbeat_seconds'),
        DEFAULT_SETTINGS.general.sse_heartbeat_seconds,
      ),
      enable_debug_panel: booleanValue(
        values.get('enable_debug_panel'),
        DEFAULT_SETTINGS.general.enable_debug_panel,
      ),
      enable_audit_log: booleanValue(
        values.get('enable_audit_log'),
        DEFAULT_SETTINGS.general.enable_audit_log,
      ),
    },
  };
}

export function buildPersistableSettingEntries(settings: Settings): PersistableSettingEntry[] {
  return [
    { key: 'write_precision', value: settings.timeseries.write_precision },
    { key: 'partition_interval', value: settings.timeseries.partition_interval },
    { key: 'batch_size', value: settings.timeseries.batch_size },
    { key: 'retention_days', value: settings.timeseries.retention_days },
    { key: 'default_interval_ms', value: settings.scheduler.default_interval_ms },
    { key: 'default_retry_count', value: settings.scheduler.default_retry_count },
    { key: 'default_retry_delay', value: settings.scheduler.default_retry_delay_ms },
    { key: 'breaker_threshold', value: settings.scheduler.breaker_threshold },
    { key: 'auto_start', value: settings.scheduler.auto_start },
    { key: 'modbus_share_enabled', value: settings.modbus_share.enabled },
    { key: 'modbus_share_bind_address', value: settings.modbus_share.bind_address },
    { key: 'modbus_share_port', value: settings.modbus_share.port },
    { key: 'modbus_share_slave_id', value: settings.modbus_share.slave_id },
    { key: 'modbus_share_base_register', value: settings.modbus_share.base_register },
    { key: 'theme', value: settings.general.theme },
    { key: 'locale', value: settings.general.locale },
    { key: 'addr_format', value: settings.general.addr_format },
    { key: 'api_base', value: settings.general.api_base },
    { key: 'api_version', value: settings.general.api_version },
    { key: 'timeout_seconds', value: settings.general.timeout_seconds },
    { key: 'log_level', value: settings.general.log_level },
    { key: 'sse_heartbeat_seconds', value: settings.general.sse_heartbeat_seconds },
    { key: 'enable_debug_panel', value: settings.general.enable_debug_panel },
    { key: 'enable_audit_log', value: settings.general.enable_audit_log },
  ];
}

export function mapDatabaseConnectorToSettingsConnector(
  connector: DatabaseConnector,
  passwordFallback = '',
): SettingsConnector {
  const config = connector.connection_config ?? {};
  const rawPort = config.port;
  const port =
    typeof rawPort === 'number'
      ? rawPort
      : typeof rawPort === 'string'
        ? parseInt(rawPort, 10) || 0
        : 0;

  return {
    id: connector.id,
    name: connector.name,
    kind: connector.kind,
    host: connectorConfigValue(config, 'host'),
    port,
    database: connectorConfigValue(config, 'database', 'dbname', 'path', 'dsn'),
    username: connectorConfigValue(config, 'user', 'username'),
    password: connectorConfigValue(config, 'password') || passwordFallback,
    schema: connectorConfigValue(config, 'schema'),
    table: connectorConfigValue(config, 'table'),
    enabled: connector.enabled,
    status:
      connector.status === 'error'
        ? 'error'
        : connector.status,
    last_check_at: connector.last_check_at ?? undefined,
    last_check_error: connector.last_check_error || undefined,
    default_write_interval_seconds: connector.default_write_interval_seconds,
  };
}

function buildConnectorConfig(connector: SettingsConnector): Record<string, unknown> {
  if (connector.kind === 'sqlite') {
    return {
      path: connector.database,
      table: connector.table,
    };
  }

  const config: Record<string, unknown> = {
    host: connector.host,
    port: connector.port,
    database: connector.database,
    user: connector.username,
    table: connector.table,
  };

  if (connector.password) {
    config.password = connector.password;
  }

  if (connector.schema) {
    config.schema = connector.schema;
  }

  return config;
}

export function buildCreateConnectorRequest(
  connector: SettingsConnector,
): CreateDatabaseConnectorRequest {
  return {
    name: connector.name,
    kind: connector.kind,
    enabled: connector.enabled,
    default_write_interval_seconds: connector.default_write_interval_seconds,
    connection_config: buildConnectorConfig(connector),
  };
}

export function buildUpdateConnectorRequest(
  connector: SettingsConnector,
): UpdateDatabaseConnectorRequest {
  return {
    name: connector.name,
    kind: connector.kind,
    enabled: connector.enabled,
    default_write_interval_seconds: connector.default_write_interval_seconds,
    connection_config: buildConnectorConfig(connector),
  };
}
