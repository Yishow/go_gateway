import type { SchemaPreviewToken } from '@/types/recordingPlan';

/** Sanitized copy of the embedded UI managed schema-preview response shape. */
export const managedSchemaPreviewToken: SchemaPreviewToken = {
  token: 'token-1',
  operation_id: 'operation-1',
  action: 'schema_apply',
  workspace_id: 'workspace-1',
  workspace_revision: 'workspace-rev-1',
  plan_id: 'group-1',
  plan_revision: 'group-rev-1',
  connector_id: 'connector-1',
  connector_revision: 'connector-rev-1',
  dialect: 'sqlite',
  database: 'managed.db',
  schema: 'main',
  table_prefix: '',
  statements: [
    'CREATE TABLE "connector_default_only" ("record_id" TEXT NOT NULL PRIMARY KEY, "group_id" TEXT NOT NULL, "device_id" TEXT NOT NULL, "bucket_start" TEXT NOT NULL, "provenance" TEXT NOT NULL, "_gw_owner_efd49c6d03ab3796062ddecc" TEXT NOT NULL DEFAULT \'gateway\', "v_633ec80d4f7f79b3687e8b31" INTEGER, "v_75e054788d18608a52b2aba9" INTEGER);',
    'CREATE TABLE "gw_effect_receipts" ("effect_key" TEXT PRIMARY KEY, "payload_digest" TEXT NOT NULL, "committed_at" TEXT NOT NULL, "_gw_effect_receipts_v1" TEXT NOT NULL DEFAULT \'gateway\');',
  ],
  tables: [
    {
      name: 'connector_default_only',
      action: 'create',
      columns: [
        'record_id', 'group_id', 'device_id', 'bucket_start', 'provenance',
        '_gw_owner_efd49c6d03ab3796062ddecc', 'v_633ec80d4f7f79b3687e8b31', 'v_75e054788d18608a52b2aba9',
      ],
    },
    {
      name: 'gw_effect_receipts',
      action: 'create',
      columns: ['effect_key', 'payload_digest', 'committed_at', '_gw_effect_receipts_v1'],
    },
  ],
  group_layout: {
    table_name: 'connector_default_only',
    columns: [
      { name: 'record_id', sql_type: 'TEXT', nullable: false, primary_key: true },
      { name: 'group_id', sql_type: 'TEXT', nullable: false, primary_key: false },
      { name: 'device_id', sql_type: 'TEXT', nullable: false, primary_key: false },
      { name: 'bucket_start', sql_type: 'TEXT', nullable: false, primary_key: false },
      { name: 'provenance', sql_type: 'TEXT', nullable: false, primary_key: false },
      { name: '_gw_owner_efd49c6d03ab3796062ddecc', sql_type: 'TEXT', nullable: false, primary_key: false },
      { name: 'v_633ec80d4f7f79b3687e8b31', sql_type: 'INTEGER', nullable: true, primary_key: false },
      { name: 'v_75e054788d18608a52b2aba9', sql_type: 'INTEGER', nullable: true, primary_key: false },
    ],
    owner_column: '_gw_owner_efd49c6d03ab3796062ddecc',
  },
  source_digest: 'bea2ef8e1e0cd485efbca09ee6c5da700a45d167605719275201dda3402204c4',
  digest: 'ed03197208a29d9345109b258a7bb47a84561b9ba4a65ac87ef52ecb16740616',
  expires_at: '2026-10-04T00:01:00Z',
  created_at: '2026-10-04T00:00:00Z',
};

export const managedSchemaPreviewResponse = { success: true, data: managedSchemaPreviewToken };
