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

export interface SchemaPreviewToken {
  token: string;
  workspace_id: string;
  plan_id: string;
  plan_revision: string;
  connector_id: string;
  table_prefix: string;
  statements: string[];
  expires_at: string;
  created_at: string;
}

export interface TestWriteResult {
  status: string;
  record_id: string;
  table: string;
  observed_at: string;
  delivered_at: string;
  message?: string;
}
