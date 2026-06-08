import type {
  DataType,
  ProtocolType,
} from './datalink';
import type { StudioV2Availability } from './studioV2Availability';
import type { StudioV2WorkspaceReadinessSummary } from './studioV2WorkspaceReadiness';

export interface StudioV2RuntimeContextDevice {
  device_id: string;
  name: string;
  protocol: ProtocolType | string;
  running: boolean;
  availability_status: StudioV2Availability['availability_status'];
  availability_reason?: string | null;
  projection_alignment?: string;
  runtime_projection_version?: string;
  workspace_projection_version?: string;
  projection_message?: string;
}

export interface StudioV2RuntimeSetupSourceRule {
  id: string;
  device_id: string;
  start_address: string;
  count: number;
  data_type: DataType;
  naming_prefix: string;
  enabled: boolean;
  revision_id?: string;
  target_data_type?: DataType | null;
}

export interface StudioV2RuntimeSetupMapping {
  id?: string;
  rule_id: string;
  device_id: string;
  point_id: string;
  address: string;
  tag_id?: string;
  tag_key?: string;
  display_name?: string;
  unit?: string;
  target_type?: DataType;
  enabled: boolean;
  status?: string;
}

export interface StudioV2RuntimeSetupDatabaseConfig {
  id: string;
  workspace_id?: string;
  kind: 'sqlite' | 'postgres' | 'mysql' | 'sqlserver';
  name: string;
  host?: string;
  port?: number;
  database: string;
  username?: string;
  schema?: string;
  table: string;
  write_mode: 'insert' | 'upsert';
  write_interval_seconds: number;
  timestamp_column?: string;
  status?: string;
}

export interface StudioV2RuntimeSetupDatabaseTarget {
  id: string;
  tag_id: string;
  connector_id: string;
  table_schema: string;
  table_name: string;
  column_name: string;
  write_mode: 'insert' | 'upsert';
  timestamp_column?: string | null;
  group_key?: string | null;
  write_interval_seconds?: number | null;
  enabled: boolean;
}

export interface StudioV2RuntimeSetupContext {
  readiness_summary?: StudioV2WorkspaceReadinessSummary | null;
  source_rules: StudioV2RuntimeSetupSourceRule[];
  mappings: StudioV2RuntimeSetupMapping[];
  database_config?: StudioV2RuntimeSetupDatabaseConfig | null;
  database_targets: StudioV2RuntimeSetupDatabaseTarget[];
}

export interface StudioV2RuntimeContext {
  workspace_id: string;
  devices: StudioV2RuntimeContextDevice[];
  default_device_id: string | null;
  setup?: StudioV2RuntimeSetupContext | null;
}
