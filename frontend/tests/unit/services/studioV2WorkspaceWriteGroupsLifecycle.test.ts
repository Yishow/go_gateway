import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import {
  studioV2WorkspaceWriteGroupsAPI,
} from '@/services/studioV2WorkspaceWriteGroups';
import type { WriteGroup, WriteGroupMutationResponse } from '@/types/studioV2WriteGroup';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

const disabledGroup: WriteGroup = {
  id: 'group-server-id',
  workspace_id: 'workspace-1',
  revision: 'group-rev-2',
  applied_revision: 'applied-rev-1',
  name: 'raw temperatures',
  status: 'disabled',
  members: [{
    device_id: 'device-1',
    point_id: 'point-1',
    tag_id: 'tag-1',
    source_revision: 'source-1',
    mapping_revision: 'mapping-1',
    target_column: 'temperature',
    required: true,
  }],
  destination: {
    connector_id: 'connector-1',
    connector_revision: 'connector-1-rev',
    database: 'gateway.db',
    table_schema: 'main',
    table_name: 'raw_values',
    storage_strategy: 'custom',
  },
  row_policy: {
    interval_seconds: 0,
    allowed_lateness_seconds: 0,
  },
  write_policy: {
    mode: 'append',
    dedupe_capability: 'limited',
  },
  migration: {},
  created_at: '2026-10-02T00:00:00Z',
  updated_at: '2026-10-02T00:01:00Z',
};

const disabledResponse: WriteGroupMutationResponse = {
  workspace_revision: 'workspace-rev-3',
  group: disabledGroup,
};

const disableRequestWithForgedFields = {
  workspace_id: 'workspace-1',
  expected_workspace_revision: 'workspace-rev-2',
  expected_group_revision: 'group-rev-1',
  expected_connector_revision: 'connector-1-rev',
  id: 'forged-group-id',
  revision: 'forged-group-revision',
  applied_revision: 'forged-applied-revision',
  status: 'running',
  group: { id: 'forged-group-id', status: 'running' },
};

function envelope<T>(data: T) {
  return { data: { success: true, data } };
}

describe('studio V2 workspace write-group disable service', () => {
  beforeEach(() => {
    vi.mocked(studioV2DatalinkApi.get).mockReset();
    vi.mocked(studioV2DatalinkApi.post).mockReset();
    vi.mocked(studioV2DatalinkApi.put).mockReset();
    vi.mocked(studioV2DatalinkApi.delete).mockReset();
  });

  it('disables an encoded group with exact CAS payload and reloads server state', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope(disabledResponse) as never);
    const request = structuredClone(disableRequestWithForgedFields);
    const before = structuredClone(request);

    await expect(studioV2WorkspaceWriteGroupsAPI.disable('group/with slash', request))
      .resolves.toEqual(disabledResponse);
    expect(request).toEqual(before);
    expect(studioV2DatalinkApi.post).toHaveBeenCalledWith(
      '/studio-v2/workspace/write-groups/group%2Fwith%20slash/disable',
      {
        workspace_id: 'workspace-1',
        expected_workspace_revision: 'workspace-rev-2',
        expected_group_revision: 'group-rev-1',
        expected_connector_revision: 'connector-1-rev',
      },
    );
    const payload = vi.mocked(studioV2DatalinkApi.post).mock.calls[0]?.[1] as Record<string, unknown>;
    expect(Object.keys(payload).sort()).toEqual([
      'expected_connector_revision',
      'expected_group_revision',
      'expected_workspace_revision',
      'workspace_id',
    ]);
    expect(payload).not.toHaveProperty('group');
    expect(payload).not.toHaveProperty('status');

    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(envelope(disabledResponse) as never);
    const reloaded = await studioV2WorkspaceWriteGroupsAPI.get('group/with slash');
    expect(reloaded.group).toMatchObject({
      id: 'group-server-id',
      revision: 'group-rev-2',
      applied_revision: 'applied-rev-1',
      status: 'disabled',
    });
    expect(reloaded.group.members).toEqual(disabledGroup.members);
  });

  it('rejects malformed disable success without retaining backend diagnostics', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope({
      ...disabledResponse,
      group: { ...disabledGroup, revision: undefined, detail: 'postgres://secret' },
    }) as never);

    const error = await studioV2WorkspaceWriteGroupsAPI.disable('group-server-id', disableRequestWithForgedFields)
      .catch((value: unknown) => value);
    expect(error).toEqual(expect.objectContaining({
      name: 'WriteGroupResponseError',
      message: 'write-group disable response was invalid',
    }));
    expect(error).not.toHaveProperty('detail');
  });

  it('preserves typed failure metadata and rejected HTTP errors for stale CAS classification', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce({
      data: {
        success: false,
        error: {
          code: 'revision_mismatch',
          message: 'private backend details',
          retryable: true,
          action: 'reload',
          request_id: 'req-disable-1',
        },
      },
    } as never);

    await expect(studioV2WorkspaceWriteGroupsAPI.disable('group-server-id', disableRequestWithForgedFields))
      .rejects.toMatchObject({
        name: 'WriteGroupResponseError',
        code: 'revision_mismatch',
        retryable: true,
        action: 'reload',
        request_id: 'req-disable-1',
      });

    const transportError = {
      response: {
        status: 409,
        data: { success: false, error: { code: 'revision_mismatch', request_id: 'req-disable-2' } },
      },
    };
    vi.mocked(studioV2DatalinkApi.post).mockRejectedValueOnce(transportError);
    await expect(studioV2WorkspaceWriteGroupsAPI.disable('group-server-id', disableRequestWithForgedFields))
      .rejects.toBe(transportError);
  });
});
