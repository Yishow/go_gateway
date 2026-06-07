import type { RuntimeStatus } from './datalink';
import type { RuntimeTruthState } from './runtimeTruth';

export type RuntimeFlowStage =
  | 'collector'
  | 'mapping'
  | 'runtime_projection'
  | 'database_delivery'
  | string;

export type RuntimeFlowStageStatus =
  | 'success'
  | 'failed'
  | 'unknown'
  | string;

export interface RuntimeFlowDiagnosticStage {
  stage: RuntimeFlowStage;
  status: RuntimeFlowStageStatus;
  observed_at?: string | null;
  reason?: string;
}

export interface RuntimeFlowDiagnostic {
  scope: string;
  device_id?: string;
  point_id?: string;
  tag_id?: string;
  last_success_at?: string | null;
  last_failure_at?: string | null;
  latest_successful_stage?: RuntimeFlowStage;
  failure_stage?: RuntimeFlowStage;
  failure_reason?: string;
  stages: RuntimeFlowDiagnosticStage[];
}

export type RuntimeStatusWithDiagnostics = RuntimeStatus & {
  diagnostics?: RuntimeFlowDiagnostic[];
  snapshot_state?: RuntimeTruthState;
};
