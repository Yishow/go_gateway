import type { StudioV2LocatorReadScope } from '../types/studioV2LocatorEntry';
export type { StudioV2LocatorReadScope } from '../types/studioV2LocatorEntry';
import { captureStudioLocatorReadScope, studioLocatorReadParams, assertStudioLocatorWorkspaceBinding, throwStudioLocatorReadError } from './studioEntityLocator';
import type {
  APIResponse,
  StudioV2WorkspaceDatabaseConfigRecord,
  StudioV2WorkspaceDatabaseRowGroupRecord,
  StudioV2WorkspaceDatabaseTargetRecord,
} from '../types/datalink';
import type { StudioV2RuntimeAppliedRecord } from '../types/studioV2RuntimeApply';
import type { StudioV2DatabaseMetadata, StudioV2TableInspectionStatus } from '../types/studioV2DatabaseMetadata';
import type { StudioV2MySQLRecordingCapability } from '../types/studioV2MySQLRecordingCapability';
import type { StudioV2SQLServerRecordingCapability } from '../types/studioV2SQLServerRecordingCapability';
import { studioV2DatalinkApi } from './studioV2Workspace';

export interface StudioV2WorkspaceDatabaseConfigRequest {
  connector_id?: string;
  expected_connector_revision?: string;
  expected_setup_revision?: string;
  kind: StudioV2WorkspaceDatabaseConfigRecord['kind'];
  name: string;
  host: string;
  port: number;
  database: string;
  username: string;
  password?: string;
  clear_password?: boolean;
  schema: string;
  table: string;
  write_mode: StudioV2WorkspaceDatabaseConfigRecord['write_mode'];
  write_interval_seconds: number;
  timestamp_column: string;
  row_groups?: StudioV2WorkspaceDatabaseRowGroupRecord[];
}

export interface StudioV2WorkspaceDatabaseTargetRequest {
  column_name: string;
  enabled: boolean;
  row_group_id?: string;
  expected_setup_revision?: string;
}

export interface StudioV2WorkspaceSchemaGenerateResult {
  connector_id: string;
  dry_run: boolean;
  statements: string[];
  executed: number;
}

function capabilityText(record: Record<string, unknown>, key: string): string {
  const value = record[key];
  if (typeof value !== 'string' || value.length > 512) throw new Error('MySQL capability response is malformed');
  return value;
}

/** Parses the server projection and rejects any success bound to another saved connector revision. */
export function parseStudioV2MySQLRecordingCapability(
  value: unknown,
  expectedConnectorId: string,
  expectedConnectorRevision: string,
): StudioV2MySQLRecordingCapability {
  if (!value || typeof value !== 'object' || Array.isArray(value)) throw new Error('MySQL capability response is malformed');
  const record = value as Record<string, unknown>;
  const connectorId = capabilityText(record, 'connector_id');
  const connectorRevision = capabilityText(record, 'connector_revision');
  const database = capabilityText(record, 'database');
  const serverFamily = capabilityText(record, 'server_family');
  const serverVersion = capabilityText(record, 'server_version');
  const capabilityRevision = capabilityText(record, 'capability_revision');
  const metadataStatus = capabilityText(record, 'metadata_status');
  const booleanFields = [
    'connected', 'transactional_engine', 'supports_managed_long', 'supports_native_partition',
    'supports_group_test_write', 'supports_missing_repair',
  ] as const;
  if (connectorId !== expectedConnectorId || connectorRevision !== expectedConnectorRevision ||
    !connectorId || !connectorRevision || !capabilityRevision ||
    booleanFields.some((field) => typeof record[field] !== 'boolean') ||
    !Array.isArray(record.reasons) || record.reasons.length > 64 ||
    record.reasons.some((reason) => typeof reason !== 'string' || !reason || reason.length > 128)) {
    throw new Error('MySQL capability response does not match the saved connector');
  }
  return {
    connector_id: connectorId,
    connector_revision: connectorRevision,
    database,
    server_family: serverFamily,
    server_version: serverVersion,
    capability_revision: capabilityRevision,
    connected: record.connected as boolean,
    metadata_status: metadataStatus,
    transactional_engine: record.transactional_engine as boolean,
    supports_managed_long: record.supports_managed_long as boolean,
    supports_native_partition: record.supports_native_partition as boolean,
    supports_group_test_write: record.supports_group_test_write as boolean,
    supports_missing_repair: record.supports_missing_repair as boolean,
    reasons: record.reasons as string[],
  };
}

const SQLSERVER_PERMISSION_FIELDS = [
  'verified', 'database_view_definition', 'schema_view_definition',
  'create_table', 'alter_schema', 'select_schema', 'insert_schema', 'delete_schema', 'alter_any_dataspace',
] as const;

/** Parses the server projection for SQL Server and rejects any success bound to another saved connector revision. */
export function parseStudioV2SQLServerRecordingCapability(
  value: unknown,
  expectedConnectorId: string,
  expectedConnectorRevision: string,
): StudioV2SQLServerRecordingCapability {
  if (!value || typeof value !== 'object' || Array.isArray(value)) throw new Error('SQL Server capability response is malformed');
  const record = value as Record<string, unknown>;
  const connectorId = capabilityText(record, 'connector_id');
  const connectorRevision = capabilityText(record, 'connector_revision');
  const database = capabilityText(record, 'database');
  const schema = capabilityText(record, 'schema');
  const serverVersion = capabilityText(record, 'server_version');
  const serverEdition = capabilityText(record, 'server_edition');
  const capabilityRevision = capabilityText(record, 'capability_revision');
  const metadataStatus = capabilityText(record, 'metadata_status');

  const engineEdition = record.engine_edition;
  const compatibilityLevel = record.compatibility_level;
  if (typeof engineEdition !== 'number' || !Number.isInteger(engineEdition) ||
      typeof compatibilityLevel !== 'number' || !Number.isInteger(compatibilityLevel)) {
    throw new Error('SQL Server capability response is malformed');
  }

  const booleanFields = [
    'connected', 'supports_managed_long', 'supports_native_partition',
    'supports_group_test_write', 'supports_missing_repair',
  ] as const;

  const permissions = record.permissions as Record<string, unknown> | undefined;
  if (!permissions || typeof permissions !== 'object' || Array.isArray(permissions)) {
    throw new Error('SQL Server capability response is malformed');
  }

  if (connectorId !== expectedConnectorId || connectorRevision !== expectedConnectorRevision ||
    !connectorId || !connectorRevision || !capabilityRevision ||
    booleanFields.some((field) => typeof record[field] !== 'boolean') ||
    SQLSERVER_PERMISSION_FIELDS.some((field) => typeof permissions[field] !== 'boolean') ||
    !Array.isArray(record.reasons) || record.reasons.length > 64 ||
    record.reasons.some((reason) => typeof reason !== 'string' || !reason || reason.length > 128)) {
    throw new Error('SQL Server capability response does not match the saved connector');
  }

  return {
    connector_id: connectorId,
    connector_revision: connectorRevision,
    database,
    schema,
    server_version: serverVersion,
    server_edition: serverEdition,
    engine_edition: engineEdition,
    compatibility_level: compatibilityLevel,
    connected: record.connected as boolean,
    metadata_status: metadataStatus,
    permissions: {
      verified: permissions.verified as boolean,
      database_view_definition: permissions.database_view_definition as boolean,
      schema_view_definition: permissions.schema_view_definition as boolean,
      create_table: permissions.create_table as boolean,
      alter_schema: permissions.alter_schema as boolean,
      select_schema: permissions.select_schema as boolean,
      insert_schema: permissions.insert_schema as boolean,
      delete_schema: permissions.delete_schema as boolean,
      alter_any_dataspace: permissions.alter_any_dataspace as boolean,
    },
    supports_managed_long: record.supports_managed_long as boolean,
    supports_native_partition: record.supports_native_partition as boolean,
    supports_group_test_write: record.supports_group_test_write as boolean,
    supports_missing_repair: record.supports_missing_repair as boolean,
    capability_revision: capabilityRevision,
    reasons: record.reasons as string[],
  };
}

const INSPECTION_STATUSES: ReadonlySet<string> = new Set(['exists', 'missing', 'forbidden', 'failed']);
const METADATA_SCOPE_FIELDS = ['workspace_id', 'connector_id', 'connector_revision', 'database', 'schema', 'table'] as const;

/** 解析實際欄位查詢結果；格式不符即丟錯，不以空欄位代替。非 exists 的結果一律不帶欄位。 */
export function parseStudioV2DatabaseMetadata(value: unknown): StudioV2DatabaseMetadata {
  if (!value || typeof value !== 'object') {
    throw new Error('Database metadata response is missing');
  }
  const record = value as Record<string, unknown>;
  const status = record.inspection_status;
  if (typeof status !== 'string' || !INSPECTION_STATUSES.has(status) || !Array.isArray(record.columns) ||
    METADATA_SCOPE_FIELDS.some((field) => typeof record[field] !== 'string')) {
    throw new Error('Database metadata response is malformed');
  }
  const columns = record.columns.map((column: unknown) => {
    const item = (column && typeof column === 'object' ? column : {}) as Record<string, unknown>;
    if (typeof item.name !== 'string' || typeof item.data_type !== 'string' ||
      typeof item.nullable !== 'boolean' || typeof item.primary_key !== 'boolean') {
      throw new Error('Database metadata column is malformed');
    }
    return {
      name: item.name, data_type: item.data_type, nullable: item.nullable, primary_key: item.primary_key,
      ...(typeof item.unique === 'boolean' ? { unique: item.unique } : {}),
    };
  });
  return {
    workspace_id: record.workspace_id as string,
    connector_id: record.connector_id as string,
    connector_revision: record.connector_revision as string,
    database: record.database as string,
    schema: record.schema as string,
    table: record.table as string,
    inspection_status: status as StudioV2TableInspectionStatus,
    ...(typeof record.reason === 'string' && record.reason ? { reason: record.reason } : {}),
    columns: status === 'exists' ? columns : [],
  };
}

export const studioV2WorkspaceDatabaseAPI = {
  async getConfig(): Promise<StudioV2WorkspaceDatabaseConfigRecord | null> {
    const res = await studioV2DatalinkApi.get<APIResponse<StudioV2WorkspaceDatabaseConfigRecord | null>>('/studio-v2/workspace/database-config');
    return res.data.data ?? null;
  },

  async updateConfig(request: StudioV2WorkspaceDatabaseConfigRequest): Promise<StudioV2RuntimeAppliedRecord<StudioV2WorkspaceDatabaseConfigRecord>> {
    const res = await studioV2DatalinkApi.put<APIResponse<StudioV2RuntimeAppliedRecord<StudioV2WorkspaceDatabaseConfigRecord>>>('/studio-v2/workspace/database-config', request);
    return res.data.data!;
  },

  async listTargets(): Promise<StudioV2WorkspaceDatabaseTargetRecord[]> {
    const res = await studioV2DatalinkApi.get<APIResponse<StudioV2WorkspaceDatabaseTargetRecord[]>>('/studio-v2/workspace/database-targets');
    return res.data.data ?? [];
  },

  async upsertTarget(pointId: string, request: StudioV2WorkspaceDatabaseTargetRequest): Promise<StudioV2RuntimeAppliedRecord<StudioV2WorkspaceDatabaseTargetRecord>> {
    const res = await studioV2DatalinkApi.put<APIResponse<StudioV2RuntimeAppliedRecord<StudioV2WorkspaceDatabaseTargetRecord>>>(`/studio-v2/workspace/database-targets/${pointId}`, request);
    return res.data.data!;
  },

  async getMySQLRecordingCapability(connectorId: string, expectedConnectorRevision: string, locatorScope?: StudioV2LocatorReadScope): Promise<StudioV2MySQLRecordingCapability> {
    const selected = captureStudioLocatorReadScope(locatorScope);
    const res = await studioV2DatalinkApi.get<APIResponse<unknown>>('/studio-v2/workspace/recording-plans/capabilities', {
      params: { kind: 'mysql', connector_id: connectorId, expected_connector_revision: expectedConnectorRevision, ...studioLocatorReadParams(selected) },
    }).catch(error => throwStudioLocatorReadError(error, selected));
    if (selected && res.data.success !== true) throw new Error('Scoped settings read unavailable.');
    assertStudioLocatorWorkspaceBinding(res.data.data, selected);
    return parseStudioV2MySQLRecordingCapability(res.data.data, connectorId, expectedConnectorRevision);
  },

  async getSQLServerRecordingCapability(connectorId: string, expectedConnectorRevision: string, locatorScope?: StudioV2LocatorReadScope): Promise<StudioV2SQLServerRecordingCapability> {
    const selected = captureStudioLocatorReadScope(locatorScope);
    const res = await studioV2DatalinkApi.get<APIResponse<unknown>>('/studio-v2/workspace/recording-plans/capabilities', {
      params: { kind: 'sqlserver', connector_id: connectorId, expected_connector_revision: expectedConnectorRevision, ...studioLocatorReadParams(selected) },
    }).catch(error => throwStudioLocatorReadError(error, selected));
    if (selected && res.data.success !== true) throw new Error('Scoped settings read unavailable.');
    assertStudioLocatorWorkspaceBinding(res.data.data, selected);
    return parseStudioV2SQLServerRecordingCapability(res.data.data, connectorId, expectedConnectorRevision);
  },

  async getMetadata(expectedConnectorRevision: string, groupScope?: { groupId: string; groupRevision: string }, locatorScope?: StudioV2LocatorReadScope): Promise<StudioV2DatabaseMetadata> {
    const selected = captureStudioLocatorReadScope(locatorScope);
    const res = await studioV2DatalinkApi.get<APIResponse<unknown>>('/studio-v2/workspace/database-metadata', {
      params: {
        expected_connector_revision: expectedConnectorRevision,
        ...studioLocatorReadParams(selected),
        ...(groupScope ? { group_id: groupScope.groupId, expected_group_revision: groupScope.groupRevision } : {}),
      },
    }).catch(error => throwStudioLocatorReadError(error, selected));
    if (selected && res.data.success !== true) throw new Error('Scoped settings read unavailable.');
    assertStudioLocatorWorkspaceBinding(res.data.data, selected);
    return parseStudioV2DatabaseMetadata(res.data.data);
  },

  async generateSchema(dryRun: boolean): Promise<StudioV2WorkspaceSchemaGenerateResult> {
    const res = await studioV2DatalinkApi.post<APIResponse<StudioV2WorkspaceSchemaGenerateResult>>('/studio-v2/workspace/database-schema/generate', { dry_run: dryRun });
    return res.data.data!;
  },
};
