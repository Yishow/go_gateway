import type { StudioV2WorkspaceReadinessIssue } from './studioV2WorkspaceReadiness';

export type StudioV2RuntimeApplyStatus = 'not_running' | 'applied' | 'apply_failed' | 'deferred';

export interface StudioV2RuntimeApply {
  runtime_apply_status: StudioV2RuntimeApplyStatus;
  runtime_apply_message?: string;
  runtime_apply_issues?: StudioV2WorkspaceReadinessIssue[];
}

export type StudioV2RuntimeAppliedRecord<T> = T & StudioV2RuntimeApply;
