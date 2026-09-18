export type RecordingPlanStatus =
  | 'draft'
  | 'validating'
  | 'ready'
  | 'running'
  | 'partial'
  | 'blocked'
  | 'paused';

export type StreamMode =
  | 'raw_history'
  | 'window_summary'
  | 'usage_interval'
  | 'state_changes'
  | 'event_log'
  | 'batch_snapshot'
  | 'latest_only';

export type RawPolicy = 'every_sample' | 'on_change' | 'sampled';

export interface PlanMember {
  member_id: string;
  measurement_id: string;
  equipment_id: string;
  name: string;
}

export interface PlanStream {
  stream_id: string;
  measurement_id: string;
  equipment_id?: string;
  mode: StreamMode;
  raw_policy?: RawPolicy;
  on_change_deadband?: number;
  max_heartbeat_seconds?: number;
  summary_interval_seconds?: number;
  usage_interval_seconds?: number;
  batch_trigger_member_id?: string;
  batch_timeout_seconds?: number;
  destination_ids?: string[];
}

export interface PlanDestination {
  destination_id: string;
  connector_id: string;
  connector_revision?: string;
  table_prefix?: string;
  write_interval_seconds?: number;
  batch_size?: number;
}

export interface RetentionPolicy {
  raw_days: number;
  summary_days: number;
  events_days: number;
  correction_horizon_hours: number;
}

export interface PlanLimits {
  max_batch_size: number;
  max_hold_seconds: number;
  max_queue_bytes: number;
}

export interface RecordingPlan {
  id: string;
  workspace_id: string;
  revision: string;
  applied_revision?: string;
  name: string;
  status: RecordingPlanStatus;
  timezone: string;
  members: PlanMember[];
  streams: PlanStream[];
  destinations: PlanDestination[];
  retention: RetentionPolicy;
  limits: PlanLimits;
  created_at?: string;
  updated_at?: string;
}

export interface ConnectorCapability {
  kind: string;
  supported: boolean;
  supports_managed_schema: boolean;
  supports_transactions: boolean;
  supports_receipts: boolean;
  supports_test_writes: boolean;
  supported_modes: string[];
  notes?: string;
}

export type SchemaPreviewTableAction = 'create' | 'unchanged';

/** One managed table in a preview: created with the listed columns, or already compatible. */
export interface SchemaPreviewTable {
  name: string;
  action: SchemaPreviewTableAction;
  columns: string[];
}

/** Server-issued preview bound to the saved workspace, plan and connector revisions. */
export interface SchemaPreviewToken {
  token: string;
  operation_id: string;
  action?: string;
  workspace_id: string;
  workspace_revision: string;
  plan_id: string;
  plan_revision: string;
  connector_id: string;
  connector_revision?: string;
  dialect?: string;
  database?: string;
  schema?: string;
  table_prefix: string;
  statements: string[];
  tables: SchemaPreviewTable[];
  no_change_reason?: string;
  digest: string;
  expires_at: string;
  created_at: string;
}

export type SchemaOperationStatus = 'pending' | 'running' | 'succeeded' | 'partial' | 'failed' | 'unknown';

/** Durable state of one confirmed schema change, as reported by the backend. */
export interface SchemaOperation {
  operation_id: string;
  action?: string;
  status: SchemaOperationStatus;
  executed_statements: number;
  verified_digest?: string;
  reason?: string;
  next_action?: string;
  created_at: string;
  updated_at: string;
  completed_at?: string;
}

export type TestWriteResultStatus = 'written_verified' | 'written_unverified' | 'failed' | 'unknown';

export interface TestWriteResult {
  status: TestWriteResultStatus;
  record_id: string;
  table: string;
  observed_at: string;
  delivered_at: string;
  operation_id?: string;
  message?: string;
}
