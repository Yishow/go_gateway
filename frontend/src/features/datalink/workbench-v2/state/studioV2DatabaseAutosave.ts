import type {
  DbConnector,
  DbRowGroup,
  DbTarget,
  Mapping,
} from './types';
import type {
  StudioV2WorkspaceDatabaseConfigRecord,
  StudioV2WorkspaceDatabaseRowGroupRecord,
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
  current?: DbConnector,
): DbConnector {
  const sameIdentity = Boolean(record.identity_revision) && current?.connector_id === record.id &&
    current?.identity_revision === record.identity_revision;
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
    password: sameIdentity ? current?.password : undefined,
    clear_password: undefined,
    // password_required 是純前端旗標（後端不回傳）：hydration 若把它丟掉，
    // 「換了連線身分但還沒重新輸入密碼」的狀態就會被洗掉，舊憑證又會被沿用。
    password_required: current?.password_required,
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
    identity_revision: record.identity_revision,
    setup_revision: record.setup_revision,
    workspace_id: record.workspace_id,
    persisted: true,
    save_state: 'saved',
    save_error: null,
  };
}

const TARGET_EDIT_FIELDS = ['tag_id', 'column_name', 'enabled', 'row_group_id'] as const;
const CONNECTOR_EDIT_FIELDS = [
  'kind', 'name', 'host', 'port', 'database', 'username', 'password', 'clear_password', 'password_required',
  'schema', 'table', 'write_mode', 'write_interval_seconds', 'timestamp_column',
] as const;

function sameFields<T>(left: T, right: T, fields: readonly (keyof T)[]): boolean {
  return fields.every((field) => left[field] === right[field]);
}

/** 同一列的使用者可編輯欄位是否仍與已送出的內容相同。 */
export function isSameStudioV2DatabaseTargetEdit(current: DbTarget | undefined, submitted: DbTarget): boolean {
  return Boolean(current) && sameFields(current as DbTarget, submitted, TARGET_EDIT_FIELDS);
}

/**
 * 回覆抵達時若同一列已有較新的本地編輯（佇列會再送一次），只採用伺服器指派的列身分，
 * 保留本地值並維持儲存中；沒有較新編輯時才完整套用回覆。
 */
export function mergeSavedStudioV2DatabaseTarget(
  current: DbTarget | undefined,
  submitted: DbTarget,
  record: StudioV2WorkspaceDatabaseTargetRecord,
): DbTarget {
  const hydrated = hydrateStudioV2DatabaseTarget(record, submitted);
  if (!current || isSameStudioV2DatabaseTargetEdit(current, submitted)) return hydrated;
  return {
    ...current, point_id: hydrated.point_id, row_id: hydrated.row_id, workspace_id: hydrated.workspace_id,
    persisted: true, save_state: 'saving', save_error: null,
  };
}

/**
 * 連線回覆同理：本地已有較新編輯時只更新伺服器指派的連線身分與版本。儲存期間改選了
 * 另一個已存連線時，回覆描述的是舊連線，只採用工作區層級的設定版本，保留本地選擇。
 */
export function mergeSavedStudioV2DatabaseConnector(
  current: DbConnector,
  submitted: DbConnector,
  record: StudioV2WorkspaceDatabaseConfigWithDelivery,
): DbConnector {
  const hydrated = hydrateStudioV2DatabaseConnector(record, submitted);
  if ((current.connector_id ?? '') !== (submitted.connector_id ?? '')) {
    return { ...current, setup_revision: hydrated.setup_revision, save_state: 'saving', save_error: null };
  }
  if (sameFields(current, submitted, CONNECTOR_EDIT_FIELDS)) return hydrated;
  return {
    ...current, connector_id: hydrated.connector_id, identity_revision: hydrated.identity_revision,
    setup_revision: hydrated.setup_revision, workspace_id: hydrated.workspace_id, persisted: true,
    save_state: 'saving', save_error: null,
  };
}

/** 分組在送出後被本地改動時保留本地分組，由排隊中的下一筆連線儲存送出。 */
export function mergeSavedStudioV2DatabaseRowGroups(
  current: DbRowGroup[] | undefined,
  submitted: DbRowGroup[] | undefined,
  records: StudioV2WorkspaceDatabaseRowGroupRecord[] | undefined,
): DbRowGroup[] {
  if (JSON.stringify(current ?? []) !== JSON.stringify(submitted ?? [])) return current ?? [];
  return hydrateStudioV2DatabaseRowGroups(records);
}

/**
 * 排隊中的連線儲存改用執行當下的伺服器版本；只有快照明確改選另一個已存連線時，
 * 才沿用快照裡該連線的身分版本。
 */
export function withLatestStudioV2DatabaseRevisions(snapshot: DbConnector, latest: DbConnector): DbConnector {
  const sameConnector = !snapshot.connector_id || snapshot.connector_id === latest.connector_id;
  return {
    ...snapshot,
    setup_revision: latest.setup_revision,
    connector_id: sameConnector ? latest.connector_id : snapshot.connector_id,
    identity_revision: sameConnector ? latest.identity_revision : snapshot.identity_revision,
  };
}

export function hydrateStudioV2DatabaseRowGroups(
  records: StudioV2WorkspaceDatabaseRowGroupRecord[] | undefined,
): DbRowGroup[] {
  return (records ?? []).map((record) => ({
    id: record.id,
    connector_id: record.connector_id,
    table_schema: record.table_schema,
    table_name: record.table_name,
    member_point_ids: [...record.member_point_ids],
    group_key_columns: record.group_key_columns ? [...record.group_key_columns] : undefined,
    unique_key_columns: record.unique_key_columns ? [...record.unique_key_columns] : undefined,
  }));
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
    row_group_id: record.row_group_id,
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
  // 連線身分換掉後，既有密碼屬於前一組連線；未重新輸入就存檔會沿用舊憑證，
  // 使用者只會看到一個看似正確的表單配上 access denied。
  if (connector.password_required && !(connector.password ?? '').length && !connector.clear_password) {
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

export function resolveStudioV2DatabaseTargetPointID(
  pointId: string,
  mapping: Mapping | undefined,
): string | undefined {
  if (!mapping?.persisted) {
    return undefined;
  }
  return mapping.persisted_point_id || pointId;
}

export function toStudioV2DatabaseConfigRequest(
  connector: DbConnector,
  rowGroups: DbRowGroup[] = [],
): StudioV2WorkspaceDatabaseConfigRequest {
  const request: StudioV2WorkspaceDatabaseConfigRequest = {
    connector_id: connector.connector_id,
    expected_connector_revision: connector.identity_revision,
    expected_setup_revision: connector.setup_revision,
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
    row_groups: rowGroups,
  };
  if (connector.clear_password) request.clear_password = true;
  // 留空表示沿用既有密碼；只有實際輸入時才送出，避免後端被覆寫成空字串。
  if (connector.kind !== 'sqlite' && connector.password && connector.password.length > 0) {
    request.password = connector.password;
  }
  return request;
}

export function toStudioV2DatabaseTargetRequest(
  target: DbTarget,
  setupRevision?: string,
): StudioV2WorkspaceDatabaseTargetRequest {
  return {
    column_name: target.column_name,
    enabled: target.enabled,
    row_group_id: target.row_group_id,
    expected_setup_revision: setupRevision,
  };
}
