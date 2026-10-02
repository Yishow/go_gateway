import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WriteGroupMigrationAPI } from '@/services/studioV2WriteGroupMigration';
import type {
  WriteGroupMigrationCandidate,
  WriteGroupMigrationPreview,
} from '@/types/studioV2WriteGroupMigration';
import type { WriteGroup } from '@/types/studioV2WriteGroup';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

const zeroGoTime = '0001-01-01T00:00:00Z';

const candidateGroup: WriteGroupMigrationCandidate = {
  id: '',
  workspace_id: 'workspace-1',
  revision: '',
  applied_revision: '',
  name: 'raw temperatures',
  status: 'draft',
  members: [{
    device_id: 'device-1',
    point_id: 'point-1',
    tag_id: 'tag-1',
    source_revision: 'source-rev-1',
    mapping_revision: 'mapping-rev-1',
    target_column: 'temperature',
    required: true,
  }],
  destination: {
    connector_id: 'connector-1',
    connector_revision: 'connector-rev-1',
    database: 'gateway.db',
    table_schema: 'main',
    table_name: 'raw_values',
    storage_strategy: 'custom',
  },
  row_policy: {
    interval_seconds: 15,
    allowed_lateness_seconds: 0,
  },
  write_policy: {
    mode: 'append',
    dedupe_capability: 'limited',
  },
  migration: {
    source_kind: 'legacy_single_mapping',
    source_ids: ['legacy-A'],
    adapter_version: 'single-mapping-v1',
    review_result: 'needs_review',
  },
  created_at: zeroGoTime,
  updated_at: zeroGoTime,
};

const persistedGroup: WriteGroup = {
  id: 'group-1',
  workspace_id: 'workspace-1',
  revision: 'group-revision-1',
  applied_revision: '',
  name: 'raw temperatures',
  status: 'draft',
  members: [{
    device_id: 'device-1',
    point_id: 'point-1',
    tag_id: 'tag-1',
    source_revision: 'source-rev-1',
    mapping_revision: 'mapping-rev-1',
    target_column: 'temperature',
    required: true,
  }],
  destination: {
    connector_id: 'connector-1',
    connector_revision: 'connector-rev-1',
    database: 'gateway.db',
    table_schema: 'main',
    table_name: 'raw_values',
    storage_strategy: 'custom',
  },
  row_policy: {
    interval_seconds: 15,
    allowed_lateness_seconds: 0,
    incomplete_policy: 'skip_row',
  },
  write_policy: {
    mode: 'append',
    dedupe_capability: 'limited',
  },
  migration: {
    source_kind: 'legacy-single-mapping',
    source_ids: ['legacy-A'],
    source_revision: 'source-digest-1',
    adapter_version: 'single-mapping-v1',
    review_result: 'needs_review',
  },
  created_at: '2026-10-02T02:00:00Z',
  updated_at: '2026-10-02T02:00:00Z',
};

const previewResponse: WriteGroupMigrationPreview = {
  workspace_id: 'workspace-1',
  workspace_revision: '',
  connector_revision: 'connector-rev-1',
  adapter_version: 'single-mapping-v1',
  review_digest: 'review-digest-1',
  items: [{
    source_id: 'legacy-A',
    source_revision: 'source-rev-1',
    status: 'needs_review',
    differences: [{
      code: 'every-sample-to-periodic-snapshot',
      message: 'legacy every-sample writes do not prove periodic snapshot equivalence',
    }],
    issues: [],
    before_intent: {
      source_id: 'legacy-A',
      device_id: 'device-1',
      point_id: 'point-1',
      tag_id: 'tag-1',
      connector_id: 'connector-1',
      connector_revision: 'connector-rev-1',
      database: 'gateway.db',
      table_schema: 'main',
      table_name: 'raw_values',
      column_name: 'temperature',
      write_mode: 'insert',
      timestamp_column: 'observed_at',
      write_interval_seconds: 0,
      interval_source: 'legacy-every-sample',
      enabled: true,
    },
    candidate_group: candidateGroup,
  }, {
    source_id: 'legacy-B',
    source_revision: 'source-rev-2',
    status: 'blocked',
    differences: [],
    issues: [{
      code: 'multiple-enabled-source-mappings',
      message: 'source identity maps to multiple enabled points',
    }],
  }],
};

const reviewResponse = {
  workspace_id: 'workspace-1',
  workspace_revision: 'workspace-rev-2',
  connector_revision: 'connector-rev-1',
  groups: [persistedGroup],
};

const requestWithForgedFields = {
  workspace_id: 'workspace-1',
  source_ids: ['legacy-A', 'legacy-B'],
  group: { id: 'forged-group-id', status: 'ready' },
  expected_workspace_revision: 'forged-workspace-revision',
  apply: true,
  connector_dsn: 'postgres://secret',
};

const reviewRequestWithForgedFields = {
  workspace_id: 'workspace-1',
  expected_workspace_revision: 'workspace-rev-1',
  expected_connector_revision: 'connector-rev-1',
  review_digest: 'review-digest-1',
  source_ids: ['legacy-A'],
  confirm_snapshot_conversion: false,
  group: { id: 'forged-group-id', status: 'ready' },
  apply: true,
  dsn: 'postgres://secret',
  server_readonly: false,
};

function envelope<T>(data: T) {
  return { data: { success: true, data } };
}

describe('studio V2 write-group migration preview service', () => {
  beforeEach(() => {
    vi.mocked(studioV2DatalinkApi.post).mockReset();
  });

  it('sends only migration scope, preserves the request, and keeps review semantics explicit', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope(previewResponse) as never);
    const request = structuredClone(requestWithForgedFields);
    const before = structuredClone(request);

    const preview = await studioV2WriteGroupMigrationAPI.previewSingleMappingMigration(request);

    expect(preview).toEqual(previewResponse);
    expect(request).toEqual(before);
    expect(studioV2DatalinkApi.post).toHaveBeenCalledWith(
      '/studio-v2/workspace/write-groups/migrations/single-mappings/preview',
      { workspace_id: 'workspace-1', source_ids: ['legacy-A', 'legacy-B'] },
    );
    const payload = vi.mocked(studioV2DatalinkApi.post).mock.calls[0]?.[1] as Record<string, unknown>;
    expect(Object.keys(payload).sort()).toEqual(['source_ids', 'workspace_id']);
    expect(payload).not.toHaveProperty('group');
    expect(payload).not.toHaveProperty('apply');
    expect(preview.workspace_revision).toBe('');
    expect(preview.items[0]?.status).toBe('needs_review');
    expect(preview.items[0]?.differences[0]?.code).toBe('every-sample-to-periodic-snapshot');
    expect(preview.items[0]?.before_intent?.write_interval_seconds).toBe(0);
    expect(preview.items[0]?.candidate_group).toMatchObject({
      id: '',
      revision: '',
      applied_revision: '',
      status: 'draft',
      created_at: zeroGoTime,
      updated_at: zeroGoTime,
    });
    expect(preview.items[1]?.status).toBe('blocked');
    expect(preview.items[1]?.candidate_group).toBeUndefined();
  });

  it.each([
    {
      name: 'candidate claims ready status',
      data: {
        ...previewResponse,
        items: [{
          ...previewResponse.items[0]!,
          candidate_group: { ...candidateGroup, status: 'ready' },
        }],
      },
    },
    {
      name: 'item claims unsupported status',
      data: {
        ...previewResponse,
        items: [{ ...previewResponse.items[0]!, status: 'ready' }],
      },
    },
    {
      name: 'preview omits review digest',
      data: (({ review_digest: _reviewDigest, ...rest }) => rest)(previewResponse),
    },
  ])('rejects $name as malformed without exposing backend text', async ({ data }) => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope({
      ...data,
      detail: 'postgres://secret backend exception',
    }) as never);

    const error = await studioV2WriteGroupMigrationAPI.previewSingleMappingMigration(
      requestWithForgedFields,
    ).catch((value: unknown) => value);
    expect(error).toEqual(expect.objectContaining({
      name: 'WriteGroupResponseError',
      message: 'write-group migration preview response was invalid',
    }));
    expect(error).not.toHaveProperty('detail');
  });

  it('preserves typed failure metadata and rejected HTTP failures for UI classification', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce({
      data: {
        success: false,
        error: {
          code: 'revision_mismatch',
          message: 'private backend details',
          retryable: true,
          action: 'reload',
          request_id: 'req-migration-preview-1',
        },
      },
    } as never);

    await expect(studioV2WriteGroupMigrationAPI.previewSingleMappingMigration(
      requestWithForgedFields,
    )).rejects.toMatchObject({
      name: 'WriteGroupResponseError',
      code: 'revision_mismatch',
      retryable: true,
      action: 'reload',
      request_id: 'req-migration-preview-1',
    });

    const transportError = {
      response: {
        status: 409,
        data: { success: false, error: { code: 'revision_mismatch', request_id: 'req-migration-preview-2' } },
      },
    };
    vi.mocked(studioV2DatalinkApi.post).mockRejectedValueOnce(transportError);
    await expect(studioV2WriteGroupMigrationAPI.previewSingleMappingMigration(
      requestWithForgedFields,
    )).rejects.toBe(transportError);
  });

  it('accepts a blocked source with no unique revision and no candidate', async () => {
    const response: WriteGroupMigrationPreview = {
      ...previewResponse,
      items: [{ ...previewResponse.items[1]!, source_revision: '' }],
    };
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope(response) as never);
    await expect(studioV2WriteGroupMigrationAPI.previewSingleMappingMigration(
      requestWithForgedFields,
    )).resolves.toEqual(response);
  });
  it.each([
    ['before intent', { before_intent: undefined }], ['candidate', { candidate_group: undefined }],
    ['before intent and candidate', { before_intent: undefined, candidate_group: undefined }],
  ] as const)('rejects a needs-review item missing %s', async (_missing, fields) => {
    const incomplete = structuredClone(previewResponse) as unknown as { items: Array<Record<string, unknown>> };
    Object.assign(incomplete.items[0]!, fields);
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope(incomplete) as never);
    await expect(studioV2WriteGroupMigrationAPI.previewSingleMappingMigration(requestWithForgedFields)).rejects.toMatchObject({
      name: 'WriteGroupResponseError', message: 'write-group migration preview response was invalid',
    });
  });

  it('retains invalid legacy intent values in a blocked preview for repair', async () => {
    const response: WriteGroupMigrationPreview = {
      ...previewResponse,
      items: [{
        source_id: 'legacy-A',
        source_revision: 'source-rev-1',
        status: 'blocked',
        differences: [],
        issues: [{ code: 'negative-write-interval', message: 'review the saved interval' }],
        before_intent: {
          ...previewResponse.items[0]!.before_intent!,
          table_name: '',
          column_name: '',
          timestamp_column: ' observed_at ',
          group_key: '',
          write_interval_seconds: -1,
        },
      }],
    };
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope(response) as never);
    await expect(studioV2WriteGroupMigrationAPI.previewSingleMappingMigration(
      requestWithForgedFields,
    )).resolves.toEqual(response);
  });

  it.each([
    { ...previewResponse.items[0]!, source_revision: '' },
    { ...previewResponse.items[0]!, differences: [] },
    { ...previewResponse.items[0]!, issues: [{ code: 'ambiguous', message: 'repair source' }] },
    { ...previewResponse.items[1]!, candidate_group: candidateGroup },
  ])('rejects contradictory review status', async (item) => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope({
      ...previewResponse, items: [item],
    }) as never);
    await expect(studioV2WriteGroupMigrationAPI.previewSingleMappingMigration(
      requestWithForgedFields,
    )).rejects.toMatchObject({ name: 'WriteGroupResponseError' });
  });
});

describe('studio V2 write-group migration review service', () => {
  beforeEach(() => {
    vi.mocked(studioV2DatalinkApi.post).mockReset();
  });

  it('sends only the reviewed migration fields, preserves false, and does not mutate input', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope(reviewResponse) as never);
    const request = structuredClone(reviewRequestWithForgedFields);
    const before = structuredClone(request);

    await expect(studioV2WriteGroupMigrationAPI.reviewSingleMappingMigration(request)).resolves.toEqual(reviewResponse);

    expect(request).toEqual(before);
    expect(studioV2DatalinkApi.post).toHaveBeenCalledWith(
      '/studio-v2/workspace/write-groups/migrations/single-mappings/review',
      {
        workspace_id: 'workspace-1',
        expected_workspace_revision: 'workspace-rev-1',
        expected_connector_revision: 'connector-rev-1',
        review_digest: 'review-digest-1',
        source_ids: ['legacy-A'],
        confirm_snapshot_conversion: false,
      },
    );
    const payload = vi.mocked(studioV2DatalinkApi.post).mock.calls[0]?.[1] as Record<string, unknown>;
    expect(Object.keys(payload).sort()).toEqual([
      'confirm_snapshot_conversion',
      'expected_connector_revision',
      'expected_workspace_revision',
      'review_digest',
      'source_ids',
      'workspace_id',
    ]);
    expect(payload).not.toHaveProperty('group');
    expect(payload).not.toHaveProperty('apply');
    expect(payload).not.toHaveProperty('dsn');
    expect(payload).not.toHaveProperty('server_readonly');
    expect(payload).toHaveProperty('confirm_snapshot_conversion', false);
  });

  it.each([
    {
      name: 'returns a group from another workspace',
      response: {
        ...reviewResponse,
        groups: [{ ...persistedGroup, workspace_id: 'workspace-foreign' }],
      },
    },
    {
      name: 'returns a group from another connector revision',
      response: {
        ...reviewResponse,
        groups: [{
          ...persistedGroup,
          destination: { ...persistedGroup.destination, connector_revision: 'connector-rev-stale' },
        }],
      },
    },
  ])('rejects $name without exposing backend details', async ({ response }) => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope({
      ...response,
      detail: 'postgres://secret backend exception',
    }) as never);

    const error = await studioV2WriteGroupMigrationAPI.reviewSingleMappingMigration(
      reviewRequestWithForgedFields,
    ).catch((value: unknown) => value);
    expect(error).toEqual(expect.objectContaining({
      name: 'WriteGroupResponseError',
      message: 'write-group migration review response was invalid',
    }));
    expect(error).not.toHaveProperty('detail');
  });

  it('rejects an unsaved candidate instead of treating it as canonical persisted state', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope({
      ...reviewResponse,
      groups: [{
        ...persistedGroup,
        id: '',
        revision: '',
        created_at: zeroGoTime,
        updated_at: zeroGoTime,
      }],
    }) as never);

    await expect(studioV2WriteGroupMigrationAPI.reviewSingleMappingMigration(
      reviewRequestWithForgedFields,
    )).rejects.toMatchObject({
      name: 'WriteGroupResponseError',
      message: 'write-group migration review response was invalid',
    });
  });

  it('accepts a persisted disabled group while leaving lifecycle decisions to the domain', async () => {
    const response = {
      ...reviewResponse,
      groups: [{
        ...persistedGroup,
        status: 'disabled' as const,
        row_policy: { ...persistedGroup.row_policy, interval_seconds: 0 },
      }],
    };
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope(response) as never);

    await expect(studioV2WriteGroupMigrationAPI.reviewSingleMappingMigration(
      reviewRequestWithForgedFields,
    )).resolves.toEqual(response);
  });

  it.each([
    { code: 'revision_mismatch', request_id: 'req-migration-review-409' },
    { code: 'validation', request_id: 'req-migration-review-422' },
  ])('keeps safe typed failure metadata for a rejected review ($code)', async ({ code, request_id }) => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce({
      data: {
        success: false,
        error: {
          code,
          message: 'postgres://secret backend exception',
          retryable: false,
          action: 'reload',
          request_id,
        },
      },
    } as never);

    const error = await studioV2WriteGroupMigrationAPI.reviewSingleMappingMigration(
      reviewRequestWithForgedFields,
    ).catch((value: unknown) => value);
    expect(error).toEqual(expect.objectContaining({
      name: 'WriteGroupResponseError',
      code,
      action: 'reload',
      request_id,
      retryable: false,
      message: 'write-group migration review response was invalid',
    }));
    expect(error).not.toHaveProperty('message', 'postgres://secret backend exception');
  });
});
