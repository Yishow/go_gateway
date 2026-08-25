import type { DataType } from './datalink';

export interface TypedAPIError {
  code: string;
  message: string;
  retryable: boolean;
  request_id: string;
  action?: string;
  workspace_revision?: string;
  settings_revision?: string;
  dirty_state?: string;
}

export interface ModbusShareRecoveryStatus {
  code: string;
  retryable: boolean;
  action: string;
  request_id?: string;
}

export interface ModbusShareOwnershipProof {
  verified: boolean;
  workspace_id: string;
  source_rule_id: string;
  source_rule_revision: string;
  basis: string;
}

export interface ModbusShareLifecycleRequest {
  port?: number;
  expected_settings_revision?: string;
}

export interface ModbusShareStatus {
  workspace_id?: string;
  enabled: boolean;
  configured_enabled?: boolean;
  running?: boolean;
  failed?: boolean;
  port: number;
  address: string;
  bind_address?: string;
  listener_state?: string;
  slave_id?: number;
  capacity_registers?: number;
  bind_state: 'pass' | 'fail' | 'disabled' | 'stopped';
  lifecycle_state?: string;
  mapping_count: number;
  settings_revision?: string;
  workspace_revision?: string;
  hydration_state?: 'ready' | 'pending' | 'failed' | string;
  readiness?: boolean;
  readiness_token?: string;
  dirty_state?: string;
  recovery?: ModbusShareRecoveryStatus;
  error?: TypedAPIError;
  canonical_plan?: ModbusShareCanonicalPlan;
  canonical_desired_mappings?: ModbusShareDesiredMapping[];
  /** Candidate/projection state is backend-owned; browser state may only add labels. */
  candidate_set_status?: string;
  candidate_statuses?: string[];
  outcome?: string;
  invalidated_count?: number;
}

/** Returns the persisted Share setting, failing closed when status is unknown. */
export function isModbusShareConfiguredEnabled(status: ModbusShareStatus | null | undefined): boolean {
  if (!status) {
    return false;
  }
  if (typeof status.configured_enabled === 'boolean') {
    return status.configured_enabled;
  }
  return status.enabled === true || status.bind_state !== 'disabled';
}

/** Returns only an explicit persisted setting, keeping transient runtime state separate. */
export function modbusShareConfiguredValue(status: ModbusShareStatus | null | undefined): boolean | undefined {
  if (!status) {
    return undefined;
  }
  if (typeof status.configured_enabled === 'boolean') {
    return status.configured_enabled;
  }
  return status.bind_state === 'disabled' ? false : undefined;
}

export interface ModbusShareCanonicalPlan {
  workspace_id: string;
  workspace_revision: string;
  settings_revision: string;
  desired_mappings: ModbusShareDesiredMapping[];
  signature: string;
}

export interface ModbusShareBootstrap {
  hydration_state: 'ready' | 'pending' | 'failed' | string;
  workspace_id?: string;
  workspace_revision?: string;
  settings_revision?: string;
  readiness: boolean;
  readiness_token?: string;
  status?: ModbusShareStatus;
  canonical_plan?: ModbusShareCanonicalPlan;
  canonical_desired_mappings?: ModbusShareDesiredMapping[];
  error?: TypedAPIError;
}

/** Normalizes the workspace bootstrap companion into the status seam used by the UI. */
export function modbusShareStatusFromBootstrap(
  bootstrap: ModbusShareBootstrap | null | undefined,
): ModbusShareStatus | null {
  if (!bootstrap?.status) {
    return null;
  }
  return {
    ...bootstrap.status,
    workspace_id: bootstrap.workspace_id ?? bootstrap.status?.workspace_id,
    workspace_revision: bootstrap.workspace_revision ?? bootstrap.status?.workspace_revision,
    settings_revision: bootstrap.settings_revision ?? bootstrap.status?.settings_revision,
    hydration_state: bootstrap.hydration_state,
    readiness: bootstrap.readiness,
    readiness_token: bootstrap.readiness_token ?? bootstrap.status?.readiness_token,
    canonical_plan: bootstrap.canonical_plan ?? bootstrap.status?.canonical_plan,
    canonical_desired_mappings: bootstrap.canonical_desired_mappings ?? bootstrap.status?.canonical_desired_mappings,
    error: bootstrap.error,
  };
}

export interface ModbusShareDesiredMapping {
  workspace_id: string;
  source_rule_id: string;
  source_rule_revision: string;
  tag_id: string;
  data_type: DataType;
  share_start_register: number;
  zero_based_register: number;
  span_registers: number;
  stride_registers: number;
  /** Persisted listener capacity returned by the canonical backend projection. */
  capacity_registers?: number;
  ownership_proof?: ModbusShareOwnershipProof;
  tag_key?: string;
  display_name?: string;
}

export interface ModbusShareReconcileRequest {
  workspace_id: string;
  expected_workspace_revision: string;
  expected_settings_revision: string;
  readiness_token: string;
  canonical_plan_signature?: string;
  desired_mappings: ModbusShareDesiredMapping[];
}

export interface ModbusShareReconcileDiagnostic {
  code: string;
  severity: string;
  message: string;
  retryable: boolean;
  action?: string;
  tag_id?: string;
  rule_id?: string;
  request_id?: string;
}

export interface ModbusShareReconcileOutcome {
  outcome: 'aligned' | 'applied' | 'disabled' | 'failed' | 'dirty_unknown' | 'invalidated_unknown' | string;
  new_workspace_revision: string;
  new_readiness_token?: string;
  settings_revision: string;
  applied_count: number;
  removed_count: number;
  invalidated_count: number;
  removed_spans: ModbusShareSpanRange[];
  invalidated_spans: ModbusShareSpanRange[];
  mappings: ModbusShareMirrorMapping[];
  diagnostics?: ModbusShareReconcileDiagnostic[];
  action?: string;
  retryable?: boolean;
  request_id?: string;
}

export interface ModbusShareSpanRange {
  tag_id: string;
  start: number;
  count: number;
  data_type: DataType;
}

export interface ModbusShareMirrorMapping {
  workspace_id?: string;
  source_rule_id?: string;
  source_rule_revision?: string;
  tag_id: string;
  register: number;
  share_start_register: number;
  zero_based_register: number;
  span_registers: number;
  stride_registers: number;
  capacity_registers: number;
  data_type: DataType;
  updated_at: string;
}

export interface ModbusShareMapping {
  tag_id: string;
  register: number;
  data_type: DataType;
  updated_at: string;
}
