import type {
  DbConnector,
  DbRowGroup,
  DbTarget,
  Device,
  Mapping,
  Point,
  SettingsConnector,
} from '../../state/types';
import type { DbColumn } from '../../state/dbSchemas';
import type { StudioV2DatabaseMetadata } from '../../../../../types/studioV2DatabaseMetadata';

/** SQLite 走檔案路徑，不需要密碼；其餘類型的連線憑證都綁在身分上。 */
export function databaseConnectorNeedsPassword(kind: DbConnector['kind']): boolean {
  return kind !== 'sqlite';
}

/**
 * 判斷選取的連接器是否換掉了連線身分。身分一改，既有密碼就屬於前一組連線，
 * 不能沿用。
 */
export function isDatabaseConnectorIdentityChange(
  connector: DbConnector,
  next: Pick<DbConnector, 'kind' | 'host' | 'port' | 'database' | 'username'>,
): boolean {
  return (
    next.kind !== connector.kind ||
    (next.host ?? '') !== (connector.host ?? '') ||
    next.port !== connector.port ||
    (next.database ?? '') !== (connector.database ?? '') ||
    (next.username ?? '') !== (connector.username ?? '')
  );
}

/** 查詢結果是否屬於目前選擇的已存連線、版本、資料庫與資料表；schema 未指定時沿用伺服器預設。 */
export function isMetadataForScope(metadata: StudioV2DatabaseMetadata, connector: DbConnector): boolean {
  return (
    metadata.connector_id === connector.connector_id &&
    metadata.connector_revision === connector.identity_revision &&
    metadata.database === connector.database &&
    (!connector.schema || metadata.schema === connector.schema) &&
    metadata.table === connector.table
  );
}

/** 只採用目前已存連線範圍內、已確認存在的資料表所查得的欄位。 */
export function metadataColumnsForScope(
  metadata: StudioV2DatabaseMetadata | undefined,
  connector: DbConnector,
): DbColumn[] {
  if (!metadata || metadata.inspection_status !== 'exists' || !isMetadataForScope(metadata, connector)) {
    return [];
  }
  return metadata.columns.map((column) => ({
    name: column.name, type: column.data_type, nullable: column.nullable, primary_key: column.primary_key,
  }));
}

/** 分組綁定到連線與資料表；改選另一個已存連線時，即使表名相同也屬於新的範圍。 */
export function isRowGroupScopeChange(connector: DbConnector, patch: Partial<DbConnector>): boolean {
  return (
    (patch.connector_id !== undefined && patch.connector_id !== connector.connector_id) ||
    (patch.kind !== undefined && patch.kind !== connector.kind) ||
    (patch.schema !== undefined && patch.schema !== connector.schema) ||
    (patch.table !== undefined && patch.table !== connector.table)
  );
}

export function rowGroupsForConnector(connector: DbConnector, rowGroups: DbRowGroup[]): DbRowGroup[] {
  return rowGroups.filter((group) => (
    group.table_name === connector.table &&
    (group.table_schema ?? '') === (connector.schema ?? '')
  ));
}

export function clearTargetRowGroupAssignments(targets: Record<string, DbTarget>): Record<string, DbTarget> {
  return Object.fromEntries(
    Object.entries(targets).map(([pointId, target]) => [
      pointId,
      target.row_group_id
        ? { ...target, row_group_id: undefined }
        : target,
    ]),
  );
}

export function syncTargetRowGroupMembership(
  rowGroups: DbRowGroup[],
  pointId: string,
  rowGroupID: string | undefined,
): DbRowGroup[] {
  return rowGroups.map((group) => {
    const memberPointIDs = group.member_point_ids.filter((memberPointId) => memberPointId !== pointId);
    if (group.id !== rowGroupID) {
      return { ...group, member_point_ids: memberPointIDs };
    }
    return { ...group, member_point_ids: [...memberPointIDs, pointId] };
  });
}

/** 唯讀狀態判斷 */
export function useStep4Readonly(phase: 'idle' | 'activating' | 'done'): boolean {
  return phase !== 'idle';
}

/** 產生 Schema 預覽簽名以避免重複重新整理 */
export interface SchemaPreviewScope {
  workspaceId?: string;
  workspaceRevision?: string;
  settingsRevision?: string;
  devices?: Array<Pick<Device, 'id' | 'persisted' | 'save_state'>>;
  points?: Array<Pick<Point, 'id' | 'device_id' | 'enabled' | 'skipped'>>;
  mappings?: Record<string, Pick<Mapping, 'mapping_id' | 'tag_id' | 'enabled' | 'persisted' | 'save_state'>>;
  rowGroups?: DbRowGroup[];
}

export function buildSchemaPreviewSignature(
  connector: DbConnector,
  targets: Record<string, DbTarget>,
  scope: SchemaPreviewScope = {},
): string {
  return JSON.stringify({
    workspace_id: scope.workspaceId ?? connector.workspace_id,
    workspace_revision: scope.workspaceRevision,
    settings_revision: scope.settingsRevision,
    connector_id: connector.connector_id,
    identity_revision: connector.identity_revision,
    kind: connector.kind,
    schema: connector.schema,
    table: connector.table,
    write_mode: connector.write_mode,
    timestamp_column: connector.timestamp_column,
    connector_save_state: connector.save_state,
    devices: (scope.devices ?? [])
      .map((device) => [device.id, device.persisted, device.save_state])
      .sort(),
    points: (scope.points ?? [])
      .map((point) => [point.id, point.device_id, point.enabled, point.skipped])
      .sort(),
    mappings: Object.entries(scope.mappings ?? {})
      .map(([pointId, mapping]) => [
        pointId,
        mapping.mapping_id,
        mapping.tag_id,
        mapping.enabled,
        mapping.persisted,
        mapping.save_state,
      ])
      .sort(),
    row_groups: (scope.rowGroups ?? [])
      .map((group) => [
        group.id,
        group.connector_id,
        group.table_schema,
        group.table_name,
        [...group.member_point_ids].sort(),
        [...(group.group_key_columns ?? [])].sort(),
        [...(group.unique_key_columns ?? [])].sort(),
      ])
      .sort(),
    targets: Object.entries(targets)
      .map(([pointId, target]) => [
        pointId,
        target.tag_id,
        target.column_name,
        target.enabled,
        target.row_group_id,
        target.persisted,
        target.save_state,
      ])
      .sort(),
  });
}

/** 依據選取的連線集產生 Connector patch */
export function createPoolConnectorPatch(
  poolConn: SettingsConnector,
): Partial<DbConnector> {
  return {
    connector_id: poolConn.id,
    identity_revision: poolConn.identity_revision,
    kind: poolConn.kind,
    name: poolConn.name,
    host: poolConn.host,
    port: poolConn.port,
    database: poolConn.database,
    username: poolConn.username,
    password: undefined,
    clear_password: false,
    password_required: false,
    schema: poolConn.schema,
    table: poolConn.table,
  };
}
