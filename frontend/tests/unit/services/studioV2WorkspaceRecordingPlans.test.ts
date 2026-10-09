import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2WorkspaceRecordingPlansAPI } from '@/services/studioV2WorkspaceRecordingPlans';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import {
  isRecordingOperationUnconfirmed,
  parseRecordingConnectorCapability,
  parseRecordingSchemaPreviewToken,
  RecordingPlanResponseError,
} from '@/utils/recordingPlanJson';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

describe('studio V2 recording plan service response validation', () => {
  beforeEach(() => {
    vi.mocked(studioV2DatalinkApi.get).mockReset();
    vi.mocked(studioV2DatalinkApi.post).mockReset();
    vi.mocked(studioV2DatalinkApi.put).mockReset();
    vi.mocked(studioV2DatalinkApi.delete).mockReset();
  });

  it('accepts only the low-level test-write statuses with bounded evidence', async () => {
    // The legacy plan-level testWrite service was removed: group test writes go
    // through studioV2WorkspaceWriteGroupsAPI, so these cases now cover the
    // shared result parser through the write-group path is out of scope here.
    // The parser itself is still exercised via parseRecordingTestWriteResult.
    const { parseRecordingTestWriteResult } = await import('@/utils/recordingPlanJson');
    for (const status of ['written_verified', 'written_unverified', 'failed'] as const) {
      expect(parseRecordingTestWriteResult({
        status,
        record_id: 'record-1',
        table: 'gw_record_samples',
        observed_at: '2026-09-07T10:00:00Z',
        delivered_at: '2026-09-07T10:00:00.012Z',
      })).toMatchObject({ status });
    }
  });

  it('rejects unfamiliar or incomplete nominal 200 test-write results as unconfirmed', async () => {
    const { parseRecordingTestWriteResult } = await import('@/utils/recordingPlanJson');
    expect(parseRecordingTestWriteResult({ status: 'success', record_id: 'record-1' })).toBeNull();
    expect(parseRecordingTestWriteResult({ operation_id: 'op-malformed-1', status: 'unknown' })).toBeNull();
    expect(parseRecordingTestWriteResult({
      status: 'written_verified',
      record_id: 'record-1',
      table: 'gw_record_samples',
      observed_at: '2026-09-07T10:00:00Z',
    })).toBeNull();
  });

  it('keeps typed error metadata from a nominal failure envelope without retaining raw text', async () => {
    const failure = {
      success: false,
      error: {
        code: 'RECORDING_TEST_WRITE_PLAN_UNRESOLVED',
        message: 'raw backend details',
        retryable: false,
        action: 'wait_for_supported_operation',
        request_id: 'req-recording-service',
      },
    };
    const { parseRecordingAPIData } = await import('@/utils/recordingPlanJson');
    try {
      parseRecordingAPIData(failure, 'recording test write', () => null);
      throw new Error('expected parseRecordingAPIData to throw');
    } catch (error) {
      expect(error).toMatchObject({
        code: 'RECORDING_TEST_WRITE_PLAN_UNRESOLVED',
        action: 'wait_for_supported_operation',
        retryable: false,
        request_id: 'req-recording-service',
        outcome: 'failed',
      });
    }
  });

  const preview = {
    token: 'tok-1',
    operation_id: 'op-1',
    workspace_id: 'ws-1',
    workspace_revision: 'setup-1',
    plan_id: 'plan-1',
    plan_revision: 'rev-1',
    connector_id: 'conn-1',
    table_prefix: 'gw_record_',
    statements: ['SELECT 1'],
    tables: [{ name: 'gw_record_samples', action: 'create', columns: ['record_id'] }],
    digest: 'a'.repeat(64),
    expires_at: '2026-09-07T04:00:00Z',
    created_at: '2026-09-07T03:50:00Z',
  };

  it('preserves a saved connector revision and rejects malformed supplied revisions', () => {
    expect(parseRecordingSchemaPreviewToken({ ...preview, connector_revision: 'identity-7' }))
      .toHaveProperty('connector_revision', 'identity-7');
    expect(parseRecordingSchemaPreviewToken({ ...preview, connector_revision: 7 })).toBeNull();
  });

  it('bounds the preview statement count before mapping', () => {
    expect(parseRecordingSchemaPreviewToken({
      ...preview,
      statements: Array.from({ length: 257 }, () => 'SELECT 1'),
    })).toBeNull();
  });

  it('bounds the total preview statement size before mapping', () => {
    expect(parseRecordingSchemaPreviewToken({
      ...preview,
      statements: Array.from({ length: 5 }, () => 'x'.repeat(16 * 1024)),
    })).toBeNull();
  });

  it('bounds the capability mode array before mapping', () => {
    expect(parseRecordingConnectorCapability({
      kind: 'sqlite',
      supported: true,
      supports_managed_schema: false,
      supports_transactions: true,
      supports_receipts: true,
      supports_test_writes: false,
      supported_modes: Array.from({ length: 257 }, () => 'custom_table'),
    })).toBeNull();
  });

  it('rejects an oversized capability list before returning UI flags', async () => {
    const capability = {
      kind: 'sqlite',
      supported: true,
      supports_managed_schema: false,
      supports_transactions: true,
      supports_receipts: true,
      supports_test_writes: false,
      supported_modes: ['custom_table'],
    };
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({
      data: {
        success: true,
        data: Array.from({ length: 257 }, () => capability),
      },
    } as never);

    await expect(studioV2WorkspaceRecordingPlansAPI.capabilities()).rejects.toMatchObject({
      outcome: 'unconfirmed',
    });
  });

  it.each([500, 502, 504])('keeps HTTP %s transport failures unconfirmed', (status) => {
    expect(isRecordingOperationUnconfirmed(new RecordingPlanResponseError('recording test write'))).toBe(true);
    expect(isRecordingOperationUnconfirmed({ response: { status } })).toBe(true);
  });

  it('keeps malformed HTTP status metadata unconfirmed', () => {
    expect(isRecordingOperationUnconfirmed({ response: { status: Number.NaN } })).toBe(true);
    expect(isRecordingOperationUnconfirmed({ response: { status: '502' } })).toBe(true);
  });

  it('keeps a fixed 501 unimplemented response definite and other client errors definite', () => {
    expect(isRecordingOperationUnconfirmed({
      response: {
        status: 501,
        data: {
          error: { code: 'RECORDING_TEST_WRITE_NOT_IMPLEMENTED' },
        },
      },
    })).toBe(false);
    expect(isRecordingOperationUnconfirmed({ response: { status: 409 } })).toBe(false);
  });
});
