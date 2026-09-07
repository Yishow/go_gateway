import type { DbConnector, DbRowGroup, DbTarget, SettingsConnector } from '../../state/types';

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

export function isRowGroupScopeChange(connector: DbConnector, patch: Partial<DbConnector>): boolean {
  return (
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
export function buildSchemaPreviewSignature(
  connector: DbConnector,
  targets: Record<string, DbTarget>,
): string {
  return JSON.stringify({
    kind: connector.kind,
    host: connector.host,
    port: connector.port,
    database: connector.database,
    username: connector.username,
    password: connector.password ?? '',
    schema: connector.schema,
    table: connector.table,
    write_mode: connector.write_mode,
    timestamp_column: connector.timestamp_column,
    targets: Object.entries(targets)
      .map(([pointId, target]) => `${pointId}:${target.tag_id}:${target.column_name}:${target.enabled}`)
      .sort(),
  });
}

/** 依據選取的連線集產生 Connector patch */
export function createPoolConnectorPatch(
  connector: DbConnector,
  poolConn: SettingsConnector,
): Partial<DbConnector> {
  const identityChanged = isDatabaseConnectorIdentityChange(connector, poolConn);
  const poolPassword = (poolConn.password ?? '').trim();
  return {
    kind: poolConn.kind,
    name: poolConn.name,
    host: poolConn.host,
    port: poolConn.port,
    database: poolConn.database,
    username: poolConn.username,
    password: identityChanged ? poolPassword : (poolConn.password ?? connector.password),
    password_required: identityChanged
      && poolPassword === ''
      && databaseConnectorNeedsPassword(poolConn.kind),
    schema: poolConn.schema,
    table: poolConn.table,
  };
}
