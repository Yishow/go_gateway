export interface StudioV2MySQLRecordingCapability {
  connector_id: string;
  connector_revision: string;
  database: string;
  server_family: string;
  server_version: string;
  capability_revision: string;
  connected: boolean;
  metadata_status: string;
  transactional_engine: boolean;
  supports_managed_long: boolean;
  supports_native_partition: boolean;
  supports_group_test_write: boolean;
  supports_missing_repair: boolean;
  reasons: string[];
}
