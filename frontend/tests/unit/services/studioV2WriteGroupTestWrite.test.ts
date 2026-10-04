import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { parseRecordingConnectorCapability } from '@/utils/recordingPlanJson';
import {
  parseWriteGroupTestWriteOperation,
  parseWriteGroupTestWritePreview,
} from '@/utils/studioV2WriteGroupTestWriteJson';
import { BACKEND_ERROR_CODES } from '@/utils/backendErrorCodes';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}));

const OWNER = 'gw-test-6f1c2d3e';

function preview(overrides: Record<string, unknown> = {}) {
  return {
    token: 'tok-1', operation_id: 'op-6f1c2d3e', action: 'test_write', group_id: 'group-1', group_revision: 'rev-1',
    expires_at: '2026-10-02T12:10:00Z',
    target: { connector_id: 'connector-A', dialect: 'sqlite', database: '/data/line.db', schema: 'main', table: 'readings' },
    owner_column: 'entity', owner_value: OWNER, dedupe: 'receipt',
    values: [
      { column: 'temperature', type: 'float64', value: '1.5' },
      { column: 'entity', type: 'text', value: OWNER },
    ],
    cleanup: `remove only rows where entity = "${OWNER}"`,
    ...overrides,
  };
}

function operation(overrides: Record<string, unknown> = {}) {
  return {
    operation_id: 'op-1', action: 'test_write', status: 'succeeded', write_outcome: 'written_verified', cleanup_status: 'cleaned',
    created_at: '2026-10-02T12:00:00Z', updated_at: '2026-10-02T12:00:01Z', completed_at: '2026-10-02T12:00:01Z',
    ...overrides,
  };
}

describe('write group test write contract', () => {
  beforeEach(() => vi.clearAllMocks());

  it('previews through the encoded group route without sending a body', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce({ data: { success: true, data: preview() } } as never);
    const result = await studioV2WorkspaceWriteGroupsAPI.testWritePreview('group/with space');
    expect(studioV2DatalinkApi.post).toHaveBeenCalledWith('/studio-v2/workspace/write-groups/group%2Fwith%20space/test-write-preview');
    expect(result.owner_value).toBe(OWNER);
    expect(result.values.map((entry) => entry.column)).toEqual(['temperature', 'entity']);
  });

  it('accepts the production cleanup description with a complete owner UUID', async () => {
    const owner = 'gw-test-12345678-1234-1234-1234-123456789abc';
    const cleanup = `remove only rows where entity_id = "${owner}", and the receipt of this test when one is written`;
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce({ data: { success: true, data: preview({
      owner_column: 'entity_id', owner_value: owner, cleanup,
      values: [{ column: 'entity_id', type: 'text', value: owner }],
    }) } } as never);
    const result = await studioV2WorkspaceWriteGroupsAPI.testWritePreview('group-1');
    expect(result.cleanup).toBe(cleanup);
    expect(result.owner_value).toBe(owner);
  });

  it('keeps cleanup text bounded and identifier and owner validation intact', () => {
    for (const cleanup of ['', '  ', null, 'x'.repeat(257)]) {
      expect(parseWriteGroupTestWritePreview(preview({ cleanup }))).toBeNull();
    }
    expect(parseWriteGroupTestWritePreview(preview({ token: 'x'.repeat(129) }))).toBeNull();
    expect(parseWriteGroupTestWritePreview(preview({ owner_value: 'x'.repeat(129) }))).toBeNull();
  });

  it('confirms with only the token and operation, never a bare plan id', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce({ data: { success: true, data: operation() } } as never);
    const result = await studioV2WorkspaceWriteGroupsAPI.testWrite('group-1', { token: 'tok-1', operation_id: 'op-1' });
    expect(studioV2DatalinkApi.post).toHaveBeenCalledWith('/studio-v2/workspace/write-groups/group-1/test-write', { token: 'tok-1', operation_id: 'op-1' });
    expect(result.write_outcome).toBe('written_verified');
    expect(result.cleanup_status).toBe('cleaned');
  });

  it('reads a running operation without a result and a finished one from the shared endpoint', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({
      data: { success: true, data: { operation_id: 'op-1', action: 'test_write', status: 'running', created_at: 'a', updated_at: 'b' } },
    } as never);
    const running = await studioV2WorkspaceWriteGroupsAPI.testWriteOperation('op-1');
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/database-operations/op-1');
    expect(running.write_outcome).toBeUndefined();
    expect(running.cleanup_status).toBeUndefined();
  });

  it('treats a response that is not a recognizable operation as unconfirmed, not success', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce({ data: { success: true, data: { status: 'succeeded' } } } as never);
    await expect(studioV2WorkspaceWriteGroupsAPI.testWrite('group-1', { token: 't', operation_id: 'o' })).rejects.toThrow();
  });

  it('keeps write verification and cleanup as separate facts', () => {
    const verifiedButDirty = parseWriteGroupTestWriteOperation(operation({ cleanup_status: 'failed', cleanup_reason: 'cleanup-denied' }));
    expect(verifiedButDirty).toMatchObject({ write_outcome: 'written_verified', cleanup_status: 'failed', cleanup_reason: 'cleanup-denied' });
    const unverifiedClean = parseWriteGroupTestWriteOperation(operation({
      status: 'partial', write_outcome: 'written_unverified', reason: 'readback-denied', cleanup_status: 'cleaned',
    }));
    expect(unverifiedClean).toMatchObject({ write_outcome: 'written_unverified', reason: 'readback-denied', cleanup_status: 'cleaned' });
    expect(parseWriteGroupTestWriteOperation(operation({ status: 'unknown', write_outcome: 'unknown', cleanup_status: 'unknown' }))).not.toBeNull();
    expect(parseWriteGroupTestWriteOperation(operation({ status: 'failed', write_outcome: 'failed', cleanup_status: 'not_attempted' }))).not.toBeNull();
  });

  it('rejects contradictory operations instead of displaying them', () => {
    // A running operation cannot already carry a result or completion time.
    expect(parseWriteGroupTestWriteOperation(operation({ status: 'running' }))).toBeNull();
    expect(parseWriteGroupTestWriteOperation({
      operation_id: 'op-1', action: 'test_write', status: 'running', created_at: 'a', updated_at: 'b',
    })).not.toBeNull();
    // Status and outcome must agree.
    expect(parseWriteGroupTestWriteOperation(operation({ status: 'failed' }))).toBeNull();
    expect(parseWriteGroupTestWriteOperation(operation({ status: 'succeeded', write_outcome: 'written_unverified' }))).toBeNull();
    // A failed write wrote nothing to clean; a verified write was always followed by a cleanup attempt.
    expect(parseWriteGroupTestWriteOperation(operation({ status: 'failed', write_outcome: 'failed', cleanup_status: 'cleaned' }))).toBeNull();
    expect(parseWriteGroupTestWriteOperation(operation({ cleanup_status: 'not_attempted' }))).toBeNull();
    // Finished operations need both facts and unknown vocabulary is refused.
    expect(parseWriteGroupTestWriteOperation(operation({ cleanup_status: undefined }))).toBeNull();
    expect(parseWriteGroupTestWriteOperation(operation({ write_outcome: 'written' }))).toBeNull();
    expect(parseWriteGroupTestWriteOperation(operation({ cleanup_status: 'removed' }))).toBeNull();
    expect(parseWriteGroupTestWriteOperation(operation({ action: 'schema_apply' }))).toBeNull();
    expect(parseWriteGroupTestWriteOperation(operation({ status: 'done' }))).toBeNull();
  });

  it('rejects previews that cannot be confirmed safely', () => {
    expect(parseWriteGroupTestWritePreview(preview())).not.toBeNull();
    expect(parseWriteGroupTestWritePreview(preview({ action: 'schema_apply' }))).toBeNull();
    expect(parseWriteGroupTestWritePreview(preview({ token: '' }))).toBeNull();
    expect(parseWriteGroupTestWritePreview(preview({ expires_at: 'soon' }))).toBeNull();
    expect(parseWriteGroupTestWritePreview(preview({ values: [] }))).toBeNull();
    // The owner marker is the only thing cleanup may match: it must be among the listed values.
    expect(parseWriteGroupTestWritePreview(preview({ owner_value: 'gw-test-someone-else' }))).toBeNull();
    expect(parseWriteGroupTestWritePreview(preview({ owner_column: 'temperature' }))).toBeNull();
  });

  it('knows the typed error codes of the test write routes', () => {
    for (const code of [
      'WRITE_GROUP_TEST_WRITE_PREVIEW_NOT_FOUND', 'WRITE_GROUP_TEST_WRITE_PREVIEW_STALE', 'WRITE_GROUP_TEST_WRITE_PREVIEW_EXPIRED',
      'WRITE_GROUP_TEST_WRITE_TOKEN_KIND', 'WRITE_GROUP_TEST_WRITE_OPERATION_MISMATCH', 'WRITE_GROUP_TEST_WRITE_BUSY',
      'WRITE_GROUP_TEST_WRITE_UNSUPPORTED', 'WRITE_GROUP_TEST_WRITE_RESULT_UNKNOWN', 'WRITE_GROUP_TEST_WRITE_UNAVAILABLE',
      'RECORDING_TEST_WRITE_PLAN_UNRESOLVED',
    ]) {
      expect(BACKEND_ERROR_CODES as readonly string[], code).toContain(code);
    }
  });

  it('parses the group test write capability without enabling the legacy plan button', () => {
    const base = {
      kind: 'sqlite', supported: true, supports_managed_schema: true, supports_transactions: true, supports_receipts: true,
      supports_test_writes: false, supported_modes: ['managed_recording'],
    };
    expect(parseRecordingConnectorCapability({ ...base, supports_group_test_writes: true })).toMatchObject({
      supports_test_writes: false, supports_group_test_writes: true,
    });
    expect(parseRecordingConnectorCapability(base)?.supports_group_test_writes).toBeUndefined();
    expect(parseRecordingConnectorCapability({ ...base, supports_group_test_writes: 'yes' })).toBeNull();
  });
});
