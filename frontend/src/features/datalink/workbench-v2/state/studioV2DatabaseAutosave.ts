import type {
  DbConnector,
  DbTarget,
  Mapping,
} from './types';
import type {
  StudioV2WorkspaceDatabaseConfigRecord,
  StudioV2WorkspaceDatabaseTargetRecord,
} from '../../../../types/datalink';
import type {
  StudioV2WorkspaceDatabaseConfigRequest,
  StudioV2WorkspaceDatabaseTargetRequest,
} from '../../../../services/studioV2WorkspaceDatabase';
import type { DatabaseDeliveryOutcomeFields } from '../../../../types/databaseDelivery';

type StudioV2WorkspaceDatabaseConfigWithDelivery =
  StudioV2WorkspaceDatabaseConfigRecord & DatabaseDeliveryOutcomeFields;

export function hydrateStudioV2DatabaseConnector(
  record: StudioV2WorkspaceDatabaseConfigWithDelivery,
): DbConnector {
  return {
    kind: record.kind,
    name: record.name,
    host: record.host,
    port: record.port,
    database: record.database,
    username: record.username,
    schema: record.schema,
    table: record.table,
    write_mode: record.write_mode,
    write_interval_seconds: record.write_interval_seconds,
    timestamp_column: record.timestamp_column,
    password: undefined,
    status: record.status,
    last_schema_ensure_at: record.last_schema_ensure_at ?? null,
    last_schema_ensure_status: record.last_schema_ensure_status ?? '',
    last_schema_ensure_error: record.last_schema_ensure_error ?? '',
    last_write_at: record.last_write_at ?? null,
    last_write_status: record.last_write_status ?? '',
    last_write_error: record.last_write_error ?? '',
    last_flush_at: record.last_flush_at ?? null,
    last_flush_status: record.last_flush_status ?? '',
    last_flush_error: record.last_flush_error ?? '',
    connector_id: record.id,
    workspace_id: record.workspace_id,
    persisted: true,
    save_state: 'saved',
    save_error: null,
  };
}

export function hydrateStudioV2DatabaseTarget(
  record: StudioV2WorkspaceDatabaseTargetRecord,
  current?: DbTarget,
): DbTarget {
  return {
    ...current,
    tag_id: record.tag_id,
    column_name: record.column_name,
    enabled: record.enabled,
    point_id: record.point_id,
    row_id: record.id,
    workspace_id: record.workspace_id,
    persisted: true,
    save_state: 'saved',
    save_error: null,
  };
}

export function isStudioV2DatabaseConnectorValid(connector: DbConnector): boolean {
  if (!connector.name.trim() || !connector.database.trim() || !connector.table.trim()) {
    return false;
  }
  if (connector.write_interval_seconds <= 0) {
    return false;
  }
  if (connector.kind === 'postgres' && !connector.username.trim()) {
    return false;
  }
  return true;
}

export function isStudioV2DatabaseTargetValid(target: DbTarget, mapping: Mapping | undefined): boolean {
  return Boolean(
    target.column_name.trim() &&
    mapping?.tag_id &&
    mapping.persisted,
  );
}

export function toStudioV2DatabaseConfigRequest(
  connector: DbConnector,
): StudioV2WorkspaceDatabaseConfigRequest {
  const request: StudioV2WorkspaceDatabaseConfigRequest = {
    kind: connector.kind,
    name: connector.name,
    host: connector.host,
    port: connector.port,
    database: connector.database,
    username: connector.username,
    schema: connector.schema,
    table: connector.table,
    write_mode: connector.write_mode,
    write_interval_seconds: connector.write_interval_seconds,
    timestamp_column: connector.timestamp_column,
  };
  // 留空表示沿用既有密碼；只有實際輸入時才送出，避免後端被覆寫成空字串。
  if (connector.kind !== 'sqlite' && connector.password && connector.password.length > 0) {
    request.password = connector.password;
  }
  return request;
}

export function toStudioV2DatabaseTargetRequest(
  target: DbTarget,
): StudioV2WorkspaceDatabaseTargetRequest {
  return {
    column_name: target.column_name,
    enabled: target.enabled,
  };
}
