import type {
  APIResponse,
  StudioV2WorkspaceDatabaseConfigRecord,
  StudioV2WorkspaceDatabaseRowGroupRecord,
  StudioV2WorkspaceDatabaseTargetRecord,
} from '../types/datalink';
import type { StudioV2RuntimeAppliedRecord } from '../types/studioV2RuntimeApply';
import { studioV2DatalinkApi } from './studioV2Workspace';

export interface StudioV2WorkspaceDatabaseConfigRequest {
  kind: StudioV2WorkspaceDatabaseConfigRecord['kind'];
  name: string;
  host: string;
  port: number;
  database: string;
  username: string;
  password?: string;
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
}

export interface StudioV2WorkspaceSchemaGenerateResult {
  connector_id: string;
  dry_run: boolean;
  statements: string[];
  executed: number;
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

  async generateSchema(dryRun: boolean): Promise<StudioV2WorkspaceSchemaGenerateResult> {
    const res = await studioV2DatalinkApi.post<APIResponse<StudioV2WorkspaceSchemaGenerateResult>>('/studio-v2/workspace/database-schema/generate', { dry_run: dryRun });
    return res.data.data!;
  },
};
