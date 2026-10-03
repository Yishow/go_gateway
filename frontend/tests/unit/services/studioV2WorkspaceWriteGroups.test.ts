import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { parseStudioV2WriteGroup } from '@/utils/studioV2WriteGroupJson';
import {
  studioV2WorkspaceWriteGroupsAPI,
} from '@/services/studioV2WorkspaceWriteGroups';
import type {
  WriteGroup,
  WriteGroupCreateRequest,
  WriteGroupDraft,
  WriteGroupListResponse,
  WriteGroupMutationResponse,
  WriteGroupReadiness,
} from '@/types/studioV2WriteGroup';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

const draft: WriteGroupDraft = {
  workspace_id: 'workspace-1',
  name: 'raw temperatures',
  members: [{
    device_id: 'device-1',
    point_id: 'point-1',
    tag_id: 'tag-1',
    target_column: 'temperature',
    required: true,
  }],
  destination: {
    connector_id: 'connector-1',
    connector_revision: 'connector-1-rev',
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
};

const draftWithoutServerFields = structuredClone(draft);

const draftWithServerFields = {
  ...draftWithoutServerFields,
  destination: {
    ...draftWithoutServerFields.destination,
    database: 'client-database-proof',
    schema_revision: 'client-schema-proof',
    schema_digest: 'client-schema-digest-proof',
  },
} as unknown as WriteGroupDraft;

const canonicalGroup: WriteGroup = {
  id: 'group-server-id',
  workspace_id: 'workspace-1',
  revision: 'group-rev-1',
  applied_revision: '',
  name: draft.name,
  status: 'draft',
  members: [{ ...draft.members[0], source_revision: 'source-1', mapping_revision: 'mapping-1' }],
  destination: { ...draft.destination, database: 'gateway.db' },
  row_policy: { ...draft.row_policy },
  write_policy: { ...draft.write_policy },
  migration: {},
  created_at: '2026-10-02T00:00:00Z',
  updated_at: '2026-10-02T00:00:00Z',
};

const mutationResponse: WriteGroupMutationResponse = {
  workspace_revision: 'workspace-rev-2',
  group: canonicalGroup,
};

const listResponse: WriteGroupListResponse = {
  workspace_id: 'workspace-1',
  workspace_revision: 'workspace-rev-2',
  groups: [canonicalGroup],
};

describe('persistent basic managed role', () => {
  it('retains the server role and rejects malformed role data', () => {
    const group = { ...canonicalGroup, destination: { ...canonicalGroup.destination, storage_strategy: 'managed' } };
    expect(parseStudioV2WriteGroup({ ...group, basic_managed_device_id: 'device-1' }))
      .toMatchObject({ basic_managed_device_id: 'device-1' });
    expect(parseStudioV2WriteGroup(group)).not.toHaveProperty('basic_managed_device_id');
    for (const role of ['', 1]) {
      expect(parseStudioV2WriteGroup({ ...group, basic_managed_device_id: role })).toBeNull();
    }
    // Explicit advanced edits can change layout while preserving the historical key.
    expect(parseStudioV2WriteGroup({ ...canonicalGroup, basic_managed_device_id: 'device-1' }))
      .toMatchObject({ basic_managed_device_id: 'device-1' });
  });
});

const readinessIssue = {
  code: 'schema-unverified',
  severity: 'blocking' as const,
  step: 'Step 4' as const,
  scope: 'group-server-id',
  message: 'destination schema has not been verified',
};

const readinessResponse: WriteGroupReadiness = {
  workspace_id: 'workspace-1',
  workspace_revision: 'workspace-rev-2',
  group_id: 'group-server-id',
  group_revision: 'group-rev-1',
  applied_revision: '',
  config_ready: true,
  schema_ready: false,
  ready: false,
  issues: [readinessIssue],
};

function envelope<T>(data: T) {
  return { data: { success: true, data } };
}

describe('studio V2 workspace write-group service', () => {
  beforeEach(() => {
    vi.mocked(studioV2DatalinkApi.get).mockReset();
    vi.mocked(studioV2DatalinkApi.post).mockReset();
    vi.mocked(studioV2DatalinkApi.put).mockReset();
    vi.mocked(studioV2DatalinkApi.delete).mockReset();
  });

  it('creates from a draft, accepts the server identity, and keeps optional measurement absent', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope(mutationResponse) as never);
    const request: WriteGroupCreateRequest = {
      workspace_id: 'workspace-1',
      expected_workspace_revision: '',
      expected_connector_revision: 'connector-1-rev',
      group: structuredClone(draftWithServerFields),
    };
    const before = structuredClone(request);

    const created = await studioV2WorkspaceWriteGroupsAPI.create(request);
    expect(created).toEqual(mutationResponse);
    expect(created.group.members[0]).toMatchObject({
      source_revision: 'source-1',
      mapping_revision: 'mapping-1',
    });

    expect(request).toEqual(before);
    expect(studioV2DatalinkApi.post).toHaveBeenCalledWith(
      '/studio-v2/workspace/write-groups',
      expect.objectContaining({
        workspace_id: 'workspace-1',
        expected_workspace_revision: '',
        expected_connector_revision: 'connector-1-rev',
        group: expect.objectContaining({ workspace_id: 'workspace-1' }),
      }),
    );
    const payload = vi.mocked(studioV2DatalinkApi.post).mock.calls[0]?.[1] as Record<string, unknown>;
    expect(payload).not.toHaveProperty('expected_group_revision');
    expect((payload.group as Record<string, unknown>).status).toBeUndefined();
    expect((payload.group as Record<string, unknown>).destination).not.toHaveProperty('database');
    expect((payload.group as Record<string, unknown>).destination).not.toHaveProperty('schema_revision');
    expect((payload.group as Record<string, unknown>).destination).not.toHaveProperty('schema_digest');
    expect((payload.group as Record<string, unknown>).members).toEqual([
      expect.not.objectContaining({ source_revision: expect.anything(), mapping_revision: expect.anything() }),
    ]);
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(envelope(mutationResponse) as never);
    expect((await studioV2WorkspaceWriteGroupsAPI.get('group-server-id')).group.members[0]?.measurement_id)
      .toBeUndefined();
  });

  it('sends workspace, group, and connector CAS fields on update with an encoded id', async () => {
    vi.mocked(studioV2DatalinkApi.put).mockResolvedValueOnce(envelope(mutationResponse) as never);
    const updateGroupWithServerFields = {
      ...structuredClone(draft),
      id: canonicalGroup.id,
      revision: canonicalGroup.revision,
      applied_revision: canonicalGroup.applied_revision,
      status: canonicalGroup.status,
      migration: structuredClone(canonicalGroup.migration),
      created_at: canonicalGroup.created_at,
      updated_at: canonicalGroup.updated_at,
      members: [{ ...draft.members[0], source_revision: 'source-1', mapping_revision: 'mapping-1' }],
      destination: {
        ...structuredClone(draft.destination),
        database: 'server-database-proof',
        schema_revision: 'server-schema-proof',
        schema_digest: 'server-schema-digest-proof',
      },
    };
    const request = {
      workspace_id: 'workspace-1',
      expected_workspace_revision: 'workspace-rev-1',
      expected_group_revision: 'group-rev-1',
      expected_connector_revision: 'connector-1-rev',
      group: updateGroupWithServerFields,
    };
    const before = structuredClone(request);

    await expect(studioV2WorkspaceWriteGroupsAPI.update('group/1', request)).resolves.toEqual(mutationResponse);
    expect(request).toEqual(before);
    expect(studioV2DatalinkApi.put).toHaveBeenCalledWith(
      '/studio-v2/workspace/write-groups/group%2F1',
      expect.objectContaining({
        workspace_id: 'workspace-1',
        expected_workspace_revision: 'workspace-rev-1',
        expected_group_revision: 'group-rev-1',
        expected_connector_revision: 'connector-1-rev',
      }),
    );
    const payload = vi.mocked(studioV2DatalinkApi.put).mock.calls[0]?.[1] as {
      group: Record<string, unknown>;
    };
    expect(payload.group).not.toHaveProperty('id');
    expect(payload.group).not.toHaveProperty('revision');
    expect(payload.group).not.toHaveProperty('applied_revision');
    expect(payload.group).not.toHaveProperty('status');
    expect(payload.group).not.toHaveProperty('migration');
    expect(payload.group).not.toHaveProperty('created_at');
    expect(payload.group).not.toHaveProperty('updated_at');
    expect(payload.group.destination).toEqual(expect.objectContaining({
      connector_id: 'connector-1',
      connector_revision: 'connector-1-rev',
      table_schema: 'main',
      table_name: 'raw_values',
      storage_strategy: 'custom',
    }));
    expect(payload.group.destination).not.toHaveProperty('database');
    expect(payload.group.destination).not.toHaveProperty('schema_revision');
    expect(payload.group.destination).not.toHaveProperty('schema_digest');
    expect(payload.group.members).toEqual([expect.objectContaining({
      source_revision: 'source-1',
      mapping_revision: 'mapping-1',
    })]);
  });

  it('reads config readiness without treating schema verification as complete', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(envelope(readinessResponse) as never);

    await expect(studioV2WorkspaceWriteGroupsAPI.readiness('group/with slash')).resolves.toEqual(readinessResponse);
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith(
      '/studio-v2/workspace/write-groups/group%2Fwith%20slash/readiness',
    );
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled();
    expect(studioV2DatalinkApi.put).not.toHaveBeenCalled();
    expect(studioV2DatalinkApi.delete).not.toHaveBeenCalled();
  });

  it('reads draft interval readiness with blocking issues and no fabricated ready state', async () => {
    const draftReadiness: WriteGroupReadiness = {
      ...readinessResponse,
      config_ready: false,
      schema_ready: false,
      ready: false,
      issues: [{
        code: 'interval-required',
        severity: 'blocking',
        step: 'Step 4',
        scope: 'group-server-id',
        message: 'write group interval must be positive',
      }],
    };
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(envelope(draftReadiness) as never);

    await expect(studioV2WorkspaceWriteGroupsAPI.readiness('group-server-id')).resolves.toEqual(draftReadiness);
  });

  it.each([
    {
      name: 'schema-ready response without a server digest',
      payload: { ...readinessResponse, schema_ready: true },
    },
    {
      name: 'ready response with a false schema gate',
      payload: { ...readinessResponse, ready: true },
    },
    {
      name: 'ready response with a blocking issue',
      payload: {
        ...readinessResponse,
        schema_ready: true,
        ready: true,
        schema_digest: 'server-digest-1',
      },
    },
    {
      name: 'schema-ready response with configuration not ready',
      payload: {
        ...readinessResponse,
        config_ready: false,
        schema_ready: true,
        ready: false,
        schema_digest: 'server-digest-2',
      },
    },
    {
      name: 'configuration-only response with a stale schema digest',
      payload: { ...readinessResponse, schema_digest: 'stale-server-digest' },
    },
  ])('rejects $name as malformed without exposing backend text', async ({ payload }) => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(envelope(payload) as never);

    const error = await studioV2WorkspaceWriteGroupsAPI.readiness('group-server-id').catch((value: unknown) => value);
    expect(error).toEqual(expect.objectContaining({
      name: 'WriteGroupResponseError',
      message: 'write-group readiness response was invalid',
    }));
    expect(error).not.toHaveProperty('detail');
  });

  it('rejects missing readiness versions and invalid issue shapes with a fixed error', async () => {
    const { workspace_revision: _workspaceRevision, ...missingVersion } = readinessResponse;
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(envelope(missingVersion) as never);
    await expect(studioV2WorkspaceWriteGroupsAPI.readiness('group-server-id')).rejects.toMatchObject({
      message: 'write-group readiness response was invalid',
    });

    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(envelope({
      ...readinessResponse,
      issues: [{ ...readinessIssue, severity: 'unsafe', message: 'postgres://secret' }],
    }) as never);
    const error = await studioV2WorkspaceWriteGroupsAPI.readiness('group-server-id').catch((value: unknown) => value);
    expect(error).toEqual(expect.objectContaining({
      message: 'write-group readiness response was invalid',
    }));
    expect(error).not.toHaveProperty('message', 'postgres://secret');
  });

  it('leaves HTTP readiness failures intact for UI error classification', async () => {
    const transportError = {
      response: {
        status: 503,
        data: { success: false, error: { code: 'WRITE_GROUP_UNAVAILABLE', request_id: 'req-readiness-1' } },
      },
    };
    vi.mocked(studioV2DatalinkApi.get).mockRejectedValueOnce(transportError);

    await expect(studioV2WorkspaceWriteGroupsAPI.readiness('group-server-id')).rejects.toBe(transportError);
  });

  it('returns canonical list data and rejects a nominal 200 without workspace revision', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(envelope(listResponse) as never);
    await expect(studioV2WorkspaceWriteGroupsAPI.list()).resolves.toEqual(listResponse);

    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(envelope({
      workspace_id: 'workspace-1',
      groups: [],
    }) as never);
    await expect(studioV2WorkspaceWriteGroupsAPI.list()).rejects.toMatchObject({
      message: 'write-group list response was invalid',
    });
  });

  it('rejects a nominal 200 with duplicate canonical group IDs', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(envelope({
      ...listResponse,
      groups: [canonicalGroup, { ...canonicalGroup }],
    }) as never);

    await expect(studioV2WorkspaceWriteGroupsAPI.list()).rejects.toMatchObject({
      name: 'WriteGroupResponseError',
      message: 'write-group list response was invalid',
    });
  });

  it('reloads a persisted draft with interval zero without inventing semantics', async () => {
    const draftResponse = {
      ...canonicalGroup,
      status: 'draft',
      members: [{ ...canonicalGroup.members[0] }],
      row_policy: { ...canonicalGroup.row_policy, interval_seconds: 0 },
    };
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(envelope({
      workspace_revision: 'workspace-rev-2',
      group: draftResponse,
    }) as never);

    const loaded = await studioV2WorkspaceWriteGroupsAPI.get('draft-0');
    expect(loaded).toMatchObject({ group: { status: 'draft', row_policy: { interval_seconds: 0 } } });
    expect(loaded.group.members[0]).not.toHaveProperty('measurement_id');
    expect(loaded.group.members[0]).toMatchObject({
      source_revision: 'source-1',
      mapping_revision: 'mapping-1',
    });
  });

  it.each(['source_revision', 'mapping_revision'] as const)(
    'rejects a canonical group missing persisted %s',
    async (missingField) => {
      const member = canonicalGroup.members[0]!;
      const missingMember = missingField === 'source_revision'
        ? (({ source_revision: _sourceRevision, ...rest }) => rest)(member)
        : (({ mapping_revision: _mappingRevision, ...rest }) => rest)(member);
      vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(envelope({
        workspace_revision: 'workspace-rev-2',
        group: { ...canonicalGroup, members: [missingMember] },
      }) as never);

      await expect(studioV2WorkspaceWriteGroupsAPI.get('missing-member-revision')).rejects.toMatchObject({
        message: 'write-group get response was invalid',
      });
    },
  );

  it('rejects malformed canonical groups without retaining backend diagnostics', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce(envelope({
      workspace_id: 'workspace-1',
      workspace_revision: 'workspace-rev-2',
      groups: [{ id: 'group-1', revision: 'group-rev-1', detail: 'postgres://secret' }],
    }) as never);

    const error = await studioV2WorkspaceWriteGroupsAPI.list().catch((value: unknown) => value);
    expect(error).toEqual(expect.objectContaining({
      message: 'write-group list response was invalid',
      name: 'WriteGroupResponseError',
    }));
    expect(error).not.toHaveProperty('detail');
  });

  it('preserves typed nominal errors and leaves rejected Axios errors intact for UI classification', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce({
      data: {
        success: false,
        error: {
          code: 'revision_mismatch',
          message: 'private backend details',
          retryable: true,
          action: 'reload',
          request_id: 'req-write-group-1',
        },
      },
    } as never);

    await expect(studioV2WorkspaceWriteGroupsAPI.create({
      workspace_id: 'workspace-1',
      expected_workspace_revision: 'stale',
      expected_connector_revision: 'connector-1-rev',
      group: structuredClone(draft),
    })).rejects.toMatchObject({
      name: 'WriteGroupResponseError',
      code: 'revision_mismatch',
      retryable: true,
      action: 'reload',
      request_id: 'req-write-group-1',
    });

    const axiosError = {
      response: {
        status: 409,
        data: { success: false, error: { code: 'revision_mismatch', request_id: 'req-write-group-2' } },
      },
    };
    vi.mocked(studioV2DatalinkApi.put).mockRejectedValueOnce(axiosError);
    await expect(studioV2WorkspaceWriteGroupsAPI.update('group-1', {
      workspace_id: 'workspace-1',
      expected_workspace_revision: 'workspace-rev-1',
      expected_group_revision: 'group-rev-1',
      expected_connector_revision: 'connector-1-rev',
      group: structuredClone(draft),
    })).rejects.toBe(axiosError);
  });

  it('sends delete CAS in axios config data and returns the canonical tombstone response', async () => {
    const deletedResponse: WriteGroupMutationResponse = {
      ...mutationResponse,
      group: { ...canonicalGroup, status: 'deleted' },
    };
    vi.mocked(studioV2DatalinkApi.delete).mockResolvedValueOnce(envelope(deletedResponse) as never);
    const request = {
      workspace_id: 'workspace-1',
      expected_workspace_revision: 'workspace-rev-2',
      expected_group_revision: 'group-rev-1',
      expected_connector_revision: 'connector-1-rev',
    };

    await expect(studioV2WorkspaceWriteGroupsAPI.remove('group/1', request)).resolves.toEqual(deletedResponse);
    expect(studioV2DatalinkApi.delete).toHaveBeenCalledWith(
      '/studio-v2/workspace/write-groups/group%2F1',
      { data: request },
    );
  });
});
