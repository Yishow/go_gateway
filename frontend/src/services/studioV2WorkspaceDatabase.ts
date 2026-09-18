import type {
  APIResponse,
  StudioV2WorkspaceDatabaseConfigRecord,
  StudioV2WorkspaceDatabaseRowGroupRecord,
  StudioV2WorkspaceDatabaseTargetRecord,
} from '../types/datalink';
import type { StudioV2RuntimeAppliedRecord } from '../types/studioV2RuntimeApply';
import type { StudioV2DatabaseMetadata, StudioV2TableInspectionStatus } from '../types/studioV2DatabaseMetadata';
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

  async getMetadata(expectedConnectorRevision: string): Promise<StudioV2DatabaseMetadata> {
    const res = await studioV2DatalinkApi.get<APIResponse<unknown>>('/studio-v2/workspace/database-metadata', {
      params: { expected_connector_revision: expectedConnectorRevision },
    });
    return parseStudioV2DatabaseMetadata(res.data.data);
  },

  async generateSchema(dryRun: boolean): Promise<StudioV2WorkspaceSchemaGenerateResult> {
    const res = await studioV2DatalinkApi.post<APIResponse<StudioV2WorkspaceSchemaGenerateResult>>('/studio-v2/workspace/database-schema/generate', { dry_run: dryRun });
    return res.data.data!;
  },
};
