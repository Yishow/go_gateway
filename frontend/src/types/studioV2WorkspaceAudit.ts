export type StudioV2WorkspaceAuditEventType =
  | 'workspace_activation'
  | 'database_config_saved'
  | 'database_target_saved'
  | string;

export type StudioV2WorkspaceAuditResult =
  | 'success'
  | 'partial_success'
  | 'failure'
  | string;

export interface StudioV2WorkspaceAuditEntry {
  id: string;
  workspace_id: string;
  event_type: StudioV2WorkspaceAuditEventType;
  result: StudioV2WorkspaceAuditResult;
  scope: string;
  reference_id?: string;
  details?: string;
  occurred_at: string;
  created_at: string;
}

export interface StudioV2WorkspaceAuditHistory {
  entries: StudioV2WorkspaceAuditEntry[];
}
