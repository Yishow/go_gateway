export type StudioV2WorkspaceReadinessSeverity = 'blocking' | 'warning';
export type StudioV2WorkspaceReadinessStep = 'Step 1' | 'Step 2' | 'Step 3' | 'Step 4';

export interface StudioV2WorkspaceReadinessIssue {
  code: string;
  severity: StudioV2WorkspaceReadinessSeverity;
  step: StudioV2WorkspaceReadinessStep;
  scope: string;
  message: string;
}

export interface StudioV2WorkspaceReadinessSummary {
  ready: boolean;
  blocking_count: number;
  warning_count: number;
  issues: StudioV2WorkspaceReadinessIssue[];
}
