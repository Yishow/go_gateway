export type StudioV2TableInspectionStatus = 'exists' | 'missing' | 'forbidden' | 'failed';

export interface StudioV2DatabaseMetadataColumn {
  name: string;
  data_type: string;
  nullable: boolean;
  primary_key: boolean;
  unique?: boolean;
}

/** 已儲存工作區目標資料表的實際查詢結果；只有 exists 時 columns 才是查得的欄位。 */
export interface StudioV2DatabaseMetadata {
  workspace_id: string;
  connector_id: string;
  connector_revision: string;
  database: string;
  schema: string;
  table: string;
  inspection_status: StudioV2TableInspectionStatus;
  reason?: string;
  columns: StudioV2DatabaseMetadataColumn[];
}
