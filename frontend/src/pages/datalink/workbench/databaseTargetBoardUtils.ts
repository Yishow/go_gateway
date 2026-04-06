import type {
  DatabaseConnector,
  DatabaseConnectorKind,
  DatabaseTableColumn,
  DatabaseTableInfo,
} from '../../../types/datalink';

export type ConnectorDraft = {
  name: string;
  kind: DatabaseConnectorKind;
  enabled: boolean;
  sqliteDsn: string;
  postgresHost: string;
  postgresPort: string;
  postgresUser: string;
  postgresPassword: string;
  postgresDatabase: string;
  postgresSSLMode: string;
};

export function createEmptyDraft(): ConnectorDraft {
  return {
    name: '',
    kind: 'sqlite',
    enabled: true,
    sqliteDsn: '',
    postgresHost: '127.0.0.1',
    postgresPort: '5432',
    postgresUser: '',
    postgresPassword: '',
    postgresDatabase: '',
    postgresSSLMode: 'disable',
  };
}

export function getErrorMessage(error: unknown, fallback: string): string {
  return error instanceof Error ? error.message : fallback;
}

export function buildTableKey(table: DatabaseTableInfo): string {
  return `${table.schema}.${table.name}`;
}

export function parseTableKey(key: string): { schema: string; name: string } | null {
  const [schemaName, ...nameParts] = key.split('.');
  if (!schemaName || nameParts.length === 0) {
    return null;
  }
  return {
    schema: schemaName,
    name: nameParts.join('.'),
  };
}

export function defaultValueColumn(columns: DatabaseTableColumn[]): DatabaseTableColumn | null {
  return columns.find((column) => !column.primary_key) ?? columns[0] ?? null;
}

export function defaultTimestampColumn(columns: DatabaseTableColumn[]): DatabaseTableColumn | null {
  return (
    columns.find((column) =>
      ['ts', 'timestamp', 'created_at', 'updated_at'].includes(column.name.toLowerCase()) &&
      (column.primary_key || column.unique),
    ) ?? null
  );
}

export function applyConnectorDraft(connector: DatabaseConnector | null): ConnectorDraft {
  if (!connector) {
    return createEmptyDraft();
  }

  const config = connector.connection_config ?? {};

  return {
    name: connector.name,
    kind: connector.kind,
    enabled: connector.enabled,
    sqliteDsn: String(config.dsn ?? config.path ?? ''),
    postgresHost: String(config.host ?? '127.0.0.1'),
    postgresPort: String(config.port ?? '5432'),
    postgresUser: String(config.user ?? ''),
    postgresPassword: String(config.password ?? ''),
    postgresDatabase: String(config.database ?? config.dbname ?? ''),
    postgresSSLMode: String(config.sslmode ?? 'disable'),
  };
}

export function buildConnectionConfig(draft: ConnectorDraft): Record<string, unknown> {
  if (draft.kind === 'sqlite') {
    return {
      dsn: draft.sqliteDsn.trim(),
    };
  }

  return {
    host: draft.postgresHost.trim(),
    port: draft.postgresPort.trim(),
    user: draft.postgresUser.trim(),
    password: draft.postgresPassword,
    database: draft.postgresDatabase.trim(),
    sslmode: draft.postgresSSLMode.trim() || 'disable',
  };
}

export function statusBadgeClasses(status: DatabaseConnector['status']): string {
  switch (status) {
    case 'ready':
      return 'bg-emerald-500/10 text-emerald-200';
    case 'auth_failed':
      return 'bg-amber-500/10 text-amber-100';
    case 'unreachable':
      return 'bg-rose-500/10 text-rose-100';
    default:
      return 'bg-slate-800 text-slate-200';
  }
}
