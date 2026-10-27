export interface StudioV2SQLServerRecordingPermissions {
  verified: boolean;
  database_view_definition: boolean;
  schema_view_definition: boolean;
  create_table: boolean;
  alter_schema: boolean;
  select_schema: boolean;
  insert_schema: boolean;
  delete_schema: boolean;
  alter_any_dataspace: boolean;
}

export interface StudioV2SQLServerRecordingCapability {
  connector_id: string;
  connector_revision: string;
  database: string;
  schema: string;
  server_version: string;
  server_edition: string;
  engine_edition: number;
  compatibility_level: number;
  connected: boolean;
  metadata_status: string;
  permissions: StudioV2SQLServerRecordingPermissions;
  supports_managed_long: boolean;
  supports_native_partition: boolean;
  supports_group_test_write: boolean;
  supports_missing_repair: boolean;
  capability_revision: string;
  reasons: string[];
}
