import { describe, expect, it } from 'vitest';
import { parseRecordingSchemaPreviewToken } from '@/utils/recordingPlanJson';

const serverToken = {
  token: 'tok-1', operation_id: 'op-1', action: 'schema_apply', workspace_id: 'ws-1', workspace_revision: 'setup-9',
  plan_id: 'plan-a', plan_revision: 'rev-2', connector_id: 'db-1', connector_revision: 'identity-7', dialect: 'sqlite',
  database: '/data/line-a.db', schema: 'main', table_prefix: 'gw_record_', statements: [],
  tables: [{ name: 'gw_record_samples', action: 'unchanged' }], no_change_reason: 'schema_already_compatible',
  digest: 'a'.repeat(64), expires_at: '2026-09-16T01:00:00Z', created_at: '2026-09-16T00:50:00Z',
};

describe('parseRecordingSchemaPreviewToken', () => {
  it('keeps the protected preview scope and table summary', () => {
    expect(parseRecordingSchemaPreviewToken(serverToken)).toEqual(expect.objectContaining({
      operation_id: 'op-1', workspace_revision: 'setup-9', digest: 'a'.repeat(64), dialect: 'sqlite', schema: 'main',
      database: '/data/line-a.db', no_change_reason: 'schema_already_compatible',
      tables: [{ name: 'gw_record_samples', action: 'unchanged', columns: [] }],
    }));
  });

  it('keeps the columns of each table to create', () => {
    const parsed = parseRecordingSchemaPreviewToken({ ...serverToken, no_change_reason: undefined, statements: ['CREATE TABLE IF NOT EXISTS gw_record_samples (id TEXT)'],
      tables: [{ name: 'gw_record_samples', action: 'create', columns: ['record_id', 'observed_at'] }] });
    expect(parsed?.tables).toEqual([{ name: 'gw_record_samples', action: 'create', columns: ['record_id', 'observed_at'] }]);
    expect(parsed?.no_change_reason).toBeUndefined();
  });

  it('rejects a preview without its operation identity, workspace revision or digest', () => {
    expect(parseRecordingSchemaPreviewToken({ ...serverToken, operation_id: undefined })).toBeNull();
    expect(parseRecordingSchemaPreviewToken({ ...serverToken, workspace_revision: '' })).toBeNull();
    expect(parseRecordingSchemaPreviewToken({ ...serverToken, digest: '' })).toBeNull();
  });

  it('rejects unknown table actions and malformed table summaries', () => {
    expect(parseRecordingSchemaPreviewToken({ ...serverToken, tables: [{ name: 'gw_record_samples', action: 'drop' }] })).toBeNull();
    expect(parseRecordingSchemaPreviewToken({ ...serverToken, tables: [{ name: '', action: 'create' }] })).toBeNull();
    expect(parseRecordingSchemaPreviewToken({ ...serverToken, tables: 'gw_record_samples' })).toBeNull();
  });
});
