import type { RuleEvidenceView } from './runtimeRuleEvidence';
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
  data_type: DataType | '';
  naming_prefix: string;
  enabled: boolean;
  revision_id?: string;
  target_data_type?: DataType | '' | null;
}

export interface StudioV2RuntimeSetupMapping {
  /** Required on bounded source-page rows; absent legacy rows remain compatible. */
  rule_revision?: string;
  id?: string;
  rule_id: string;
  device_id: string;
  point_id: string;
  address: string;
  tag_id?: string;
  tag_key?: string;
  display_name?: string;
  unit?: string;
  target_type?: DataType | '';
  enabled: boolean;
  status?: string;
}

export interface StudioV2RuntimeSetupDatabaseConfig {
  id: string;
  workspace_id?: string;
  kind: 'sqlite' | 'postgres' | 'mysql' | 'sqlserver' | 'unavailable';
  name: string;
  schema?: string;
  table: string;
  write_mode: 'insert' | 'upsert' | 'unavailable';
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
  write_mode: 'insert' | 'upsert' | 'unavailable';
  timestamp_column?: string | null;
  group_key?: string | null;
  write_interval_seconds?: number | null;
  enabled: boolean;
}

export interface StudioV2RuntimeSetupContext {
  evidence_state?: 'saved-only' | 'unavailable';
  readiness_summary?: StudioV2WorkspaceReadinessSummary | null;
  source_rules: StudioV2RuntimeSetupSourceRule[];
  mappings: StudioV2RuntimeSetupMapping[];
  database_config?: StudioV2RuntimeSetupDatabaseConfig | null;
  database_targets: StudioV2RuntimeSetupDatabaseTarget[];
}

/** Explicit read selection, not an editing locator or mutation permission. */
export interface RuntimeSourceSelection {
  rule_id: string;
  rule_revision: string;
  point_id: string;
  mapping_id?: string;
}

export interface StudioV2RuntimeContext {
  /** Safe revision of the same observed workspace; absent legacy evidence stays unknown. */
  workspace_revision?: string;
  /** Only supplied for an explicit scoped read; the shared service validates it before use. */
  rule_evidence?: RuleEvidenceView | null;
  workspace_id: string;
  devices: StudioV2RuntimeContextDevice[];
  default_device_id: string | null;
  setup?: StudioV2RuntimeSetupContext | null;
}

/** Read-only page selection; no locator or mutation capability. */
export interface RuntimeSourcePageRequest {
  kind: 'rules' | 'mappings';
  query?: string;
  limit?: number;
  cursor?: string;
  rule_id?: string;
  rule_revision?: string;
}

/** Bounded safe saved summaries, not another source-facts or execution authority. */
export interface RuntimeSourcePage {
  contract_version: 1;
  workspace_id: string;
  workspace_revision: string;
  device_id: string;
  kind: 'rules' | 'mappings';
  query: string;
  limit: number;
  list_revision: string;
  rules: (StudioV2RuntimeSetupSourceRule & { revision_id: string })[];
  mappings: (StudioV2RuntimeSetupMapping & { rule_revision: string })[];
  next_cursor: string | null;
}
