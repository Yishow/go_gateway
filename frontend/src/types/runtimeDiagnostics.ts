import type { RuntimeStatus } from './datalink';
import type { RuntimeTruthState } from './runtimeTruth';

export type RuntimeFlowStage =
  | 'collector'
  | 'mapping'
  | 'runtime_projection'
  | 'database_delivery'
  | 'modbus_share_delivery'
  | string;

export type RuntimeFlowFailureCode =
  | 'runtime_snapshot_unavailable'
  | 'modbus_share_delivery'
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
  failure_code?: RuntimeFlowFailureCode;
  failure_reason?: string;
  request_id?: string;
  stages: RuntimeFlowDiagnosticStage[];
}

export interface DatabaseDeliveryDiagnostic {
  device_id?: string;
  point_id: string;
  tag_id: string;
  status: string;
  stages: string[];
  failed_stage?: string;
  error?: string;
  observed_at: string;
  last_success_at?: string | null;
  last_failure_at?: string | null;
  last_failure_reason?: string;
}

export interface ModbusShareDeliveryDiagnostic {
  device_id?: string;
  point_id: string;
  tag_id: string;
  status?: 'succeeded' | 'failed' | 'disabled' | string;
  stage?: string;
  error?: string;
  observed_at: string;
}

export type RuntimeStatusWithDiagnostics = RuntimeStatus & {
  diagnostics?: RuntimeFlowDiagnostic[];
  snapshot_state?: RuntimeTruthState;
};
