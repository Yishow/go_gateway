import { beforeEach, describe, expect, it, vi } from 'vitest';
import {
  parseStudioV2SQLServerRecordingCapability,
  studioV2WorkspaceDatabaseAPI,
} from '@/services/studioV2WorkspaceDatabase';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}));

const validEnvelope = {
  connector_id: 'sqlserver-conn-1',
  connector_revision: 'rev-alpha-1',
  database: 'gw_industrial',
  schema: 'dbo',
  server_version: '16.0.4165.4',
  server_edition: 'Developer Edition (64-bit)',
  engine_edition: 3,
  compatibility_level: 160,
  connected: true,
  metadata_status: 'verified',
  permissions: {
    verified: true,
    database_view_definition: true,
    schema_view_definition: true,
    create_table: true,
    alter_schema: true,
    select_schema: true,
    insert_schema: true,
    delete_schema: true,
    alter_any_dataspace: false,
  },
  supports_managed_long: false,
  supports_native_partition: false,
  supports_group_test_write: false,
  supports_missing_repair: false,
  capability_revision: 'cap-rev-sql-99',
  reasons: ['adapter_unverified', 'managed_long_unverified'],
};

describe('SQL Server recording capability parser & contract', () => {
  beforeEach(() => {
    vi.mocked(studioV2DatalinkApi.get).mockReset();
  });

  it('parses a valid SQL Server capability envelope successfully', () => {
    const parsed = parseStudioV2SQLServerRecordingCapability(
      validEnvelope,
      'sqlserver-conn-1',
      'rev-alpha-1',
    );
    expect(parsed.connector_id).toBe('sqlserver-conn-1');
    expect(parsed.connector_revision).toBe('rev-alpha-1');
    expect(parsed.database).toBe('gw_industrial');
    expect(parsed.schema).toBe('dbo');
    expect(parsed.server_version).toBe('16.0.4165.4');
    expect(parsed.server_edition).toBe('Developer Edition (64-bit)');
    expect(parsed.engine_edition).toBe(3);
    expect(parsed.compatibility_level).toBe(160);
    expect(parsed.connected).toBe(true);
    expect(parsed.metadata_status).toBe('verified');
    expect(parsed.permissions.schema_view_definition).toBe(true);
    expect(parsed.permissions.alter_any_dataspace).toBe(false);
    expect(parsed.supports_managed_long).toBe(false);
    expect(parsed.supports_native_partition).toBe(false);
    expect(parsed.supports_group_test_write).toBe(false);
    expect(parsed.supports_missing_repair).toBe(false);
    expect(parsed.reasons).toEqual(['adapter_unverified', 'managed_long_unverified']);
  });

  it('rejects when connector_id does not match expectedConnectorId', () => {
    expect(() =>
      parseStudioV2SQLServerRecordingCapability(validEnvelope, 'other-conn', 'rev-alpha-1'),
    ).toThrow('SQL Server capability response does not match the saved connector');
  });

  it('rejects when connector_revision does not match expectedConnectorRevision', () => {
    expect(() =>
      parseStudioV2SQLServerRecordingCapability(validEnvelope, 'sqlserver-conn-1', 'stale-rev'),
    ).toThrow('SQL Server capability response does not match the saved connector');
  });

  it('rejects malformed envelopes (null, primitive, or array)', () => {
    for (const bad of [null, undefined, 42, 'string', []]) {
      expect(() =>
        parseStudioV2SQLServerRecordingCapability(bad, 'sqlserver-conn-1', 'rev-alpha-1'),
      ).toThrow('SQL Server capability response is malformed');
    }
  });

  it('rejects non-integer engine_edition or compatibility_level', () => {
    expect(() =>
      parseStudioV2SQLServerRecordingCapability(
        { ...validEnvelope, engine_edition: '3' },
        'sqlserver-conn-1',
        'rev-alpha-1',
      ),
    ).toThrow('SQL Server capability response is malformed');

    expect(() =>
      parseStudioV2SQLServerRecordingCapability(
        { ...validEnvelope, compatibility_level: 160.5 },
        'sqlserver-conn-1',
        'rev-alpha-1',
      ),
    ).toThrow('SQL Server capability response is malformed');
  });

  it('rejects missing or non-boolean permission flags', () => {
    expect(() =>
      parseStudioV2SQLServerRecordingCapability(
        { ...validEnvelope, permissions: { ...validEnvelope.permissions, schema_view_definition: 'yes' } },
        'sqlserver-conn-1',
        'rev-alpha-1',
      ),
    ).toThrow('SQL Server capability response does not match the saved connector');

    expect(() =>
      parseStudioV2SQLServerRecordingCapability(
        { ...validEnvelope, permissions: undefined },
        'sqlserver-conn-1',
        'rev-alpha-1',
      ),
    ).toThrow('SQL Server capability response is malformed');
  });

  it('rejects invalid reasons array', () => {
    expect(() =>
      parseStudioV2SQLServerRecordingCapability(
        { ...validEnvelope, reasons: 'none' },
        'sqlserver-conn-1',
        'rev-alpha-1',
      ),
    ).toThrow('SQL Server capability response does not match the saved connector');

    expect(() =>
      parseStudioV2SQLServerRecordingCapability(
        { ...validEnvelope, reasons: [''] },
        'sqlserver-conn-1',
        'rev-alpha-1',
      ),
    ).toThrow('SQL Server capability response does not match the saved connector');
  });

  it('requests SQL Server capability through API with exact query params and parses response', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({
      data: { success: true, data: validEnvelope },
    } as never);

    const result = await studioV2WorkspaceDatabaseAPI.getSQLServerRecordingCapability(
      'sqlserver-conn-1',
      'rev-alpha-1',
    );

    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith(
      '/studio-v2/workspace/recording-plans/capabilities',
      {
        params: {
          kind: 'sqlserver',
          connector_id: 'sqlserver-conn-1',
          expected_connector_revision: 'rev-alpha-1',
        },
      },
    );
    expect(result.database).toBe('gw_industrial');
    expect(result.schema).toBe('dbo');
    expect(result.server_version).toBe('16.0.4165.4');
  });
});
