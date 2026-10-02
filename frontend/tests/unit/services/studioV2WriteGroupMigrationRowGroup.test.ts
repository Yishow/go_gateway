import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WriteGroupMigrationAPI } from '@/services/studioV2WriteGroupMigration';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type {
  WriteGroupMigrationCandidate,
  WriteGroupMigrationPreview,
  WriteGroupMigrationReviewResponse,
} from '@/types/studioV2WriteGroupMigration';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { WriteGroupDraft } from '@/types/studioV2WriteGroup';
import {
  parseWriteGroupListData,
} from '@/utils/studioV2WriteGroupJson';
import {
  parseWriteGroupMigrationPreviewData,
} from '@/utils/studioV2WriteGroupMigrationJson';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

const zeroGoTime = '0001-01-01T00:00:00Z';

const rowGroupIntent = {
  source_id: 'legacy-g',
  connector_id: 'connector-1',
  connector_revision: 'connector-rev-1',
  row_group: {
    id: 'legacy-g',
    connector_id: 'connector-1',
    table_schema: 'main',
    table_name: 'raw_values',
    member_point_ids: ['p-a', 'p-b'],
    group_key_columns: ['entity_id'],
    unique_key_columns: ['entity_id', 'observed_at'],
  },
  members: [{
    source_id: 'legacy-A',
    device_id: 'device-A',
    point_id: 'p-a',
    tag_id: 'tag-A',
    connector_id: 'connector-1',
    connector_revision: 'connector-rev-1',
    database: 'gateway.db',
    table_schema: 'main',
    table_name: 'raw_values',
    column_name: 'value',
    write_mode: 'insert',
    timestamp_column: 'observed_at',
    group_key: 'legacy-g:p-a',
    write_interval_seconds: 15,
    interval_source: 'legacy-target-override',
    enabled: true,
  }, {
    source_id: 'legacy-B',
    device_id: 'device-B',
    point_id: 'p-b',
    tag_id: 'tag-B',
    connector_id: 'connector-1',
    connector_revision: 'connector-rev-1',
    database: 'gateway.db',
    table_schema: 'main',
    table_name: 'raw_values',
    column_name: 'value',
    write_mode: 'insert',
    timestamp_column: 'observed_at',
    group_key: 'legacy-g:p-b',
    write_interval_seconds: 15,
    interval_source: 'legacy-target-override',
    enabled: true,
  }],
};

const candidateGroup: WriteGroupMigrationCandidate = {
  id: '',
  workspace_id: 'workspace-1',
  revision: '',
  applied_revision: '',
  name: 'raw values',
  status: 'draft',
  members: [{
    device_id: 'device-A',
    point_id: 'p-a',
    tag_id: 'tag-A',
    entity_key: 'legacy-g:p-a',
    source_revision: 'source-rev-A',
    mapping_revision: 'mapping-rev-A',
    target_column: 'value',
    required: true,
  }, {
    device_id: 'device-B',
    point_id: 'p-b',
    tag_id: 'tag-B',
    entity_key: 'legacy-g:p-b',
    source_revision: 'source-rev-B',
    mapping_revision: 'mapping-rev-B',
    target_column: 'value',
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
    group_key_columns: ['entity_id'],
    unique_key_columns: ['entity_id', 'observed_at'],
  },
  write_policy: { mode: 'append' },
  migration: {
    source_kind: 'legacy-row-group',
    source_ids: ['legacy-A', 'legacy-B'],
    source_revision: 'row-group-digest-1',
    adapter_version: 'row-group-v1',
    review_result: 'needs_review',
    legacy_row_group_id: 'legacy-g',
    target_mapping_points: { 'legacy-A': 'p-a', 'legacy-B': 'p-b' },
  },
  created_at: zeroGoTime,
  updated_at: zeroGoTime,
};

const rowGroupPreviewResponse: WriteGroupMigrationPreview = {
  workspace_id: 'workspace-1',
  workspace_revision: 'workspace-rev-1',
  connector_revision: 'connector-rev-1',
  adapter_version: 'row-group-v1',
  review_digest: 'row-group-digest-1',
  items: [{
    source_id: 'legacy-g',
    source_revision: 'row-group-source-rev-1',
    status: 'needs_review',
    differences: [{
      code: 'shared-column-to-entity-partition',
      message: 'shared target columns are partitioned by the persisted entity key',
    }],
    issues: [],
    before_row_group_intent: rowGroupIntent,
    candidate_group: candidateGroup,
  }, {
    source_id: 'legacy-blocked',
    source_revision: '',
    status: 'blocked',
    differences: [],
    issues: [{
      code: 'missing-unique-key',
      message: 'the row group has no persisted uniqueness metadata',
    }],
    before_row_group_intent: {
      ...rowGroupIntent,
      source_id: 'legacy-blocked',
      row_group: { ...rowGroupIntent.row_group, id: 'legacy-blocked', unique_key_columns: [] },
      members: [{ ...rowGroupIntent.members[0], device_id: '' }, rowGroupIntent.members[1]],
    },
  }],
};

const persistedGroup: WriteGroup = {
  ...candidateGroup,
  id: 'group-1',
  revision: 'group-revision-1',
  status: 'draft',
  members: candidateGroup.members.map((member) => ({
    ...member,
    source_revision: member.source_revision ?? 'source-revision',
    mapping_revision: member.mapping_revision ?? 'mapping-revision',
  })),
  created_at: '2026-10-02T02:00:00Z',
  updated_at: '2026-10-02T02:00:00Z',
};

const rowGroupReviewResponse: WriteGroupMigrationReviewResponse = {
  workspace_id: 'workspace-1',
  workspace_revision: 'workspace-rev-2',
  connector_revision: 'connector-rev-1',
  groups: [persistedGroup],
};

const reviewRequest = {
  workspace_id: 'workspace-1',
  expected_workspace_revision: 'workspace-rev-1',
  expected_connector_revision: 'connector-rev-1',
  review_digest: 'row-group-digest-1',
  source_ids: ['legacy-g'],
  confirm_snapshot_conversion: false,
  group: { id: 'forged-group', status: 'ready' },
  apply: true,
  dsn: 'postgres://secret',
  server_readonly: false,
};

function envelope<T>(data: T) {
  return { data: { success: true, data } };
}

describe('studio V2 row-group migration service', () => {
  beforeEach(() => {
    vi.mocked(studioV2DatalinkApi.post).mockReset();
    vi.mocked(studioV2DatalinkApi.put).mockReset();
  });

  it('preview sends only workspace and original row-group IDs without mutating input', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope(rowGroupPreviewResponse) as never);
    const request = {
      workspace_id: 'workspace-1',
      source_ids: ['legacy-g'],
      apply: true,
      dsn: 'postgres://secret',
    };
    const before = structuredClone(request);

    await expect(studioV2WriteGroupMigrationAPI.previewRowGroupMigration(request)).resolves.toEqual(rowGroupPreviewResponse);

    expect(request).toEqual(before);
    expect(studioV2DatalinkApi.post).toHaveBeenCalledWith(
      '/studio-v2/workspace/write-groups/migrations/row-groups/preview',
      { workspace_id: 'workspace-1', source_ids: ['legacy-g'] },
    );
  });

  it('review sends only the explicit confirmation envelope and preserves false', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope(rowGroupReviewResponse) as never);
    const request = structuredClone(reviewRequest);
    const before = structuredClone(request);

    await expect(studioV2WriteGroupMigrationAPI.reviewRowGroupMigration(request)).resolves.toEqual(rowGroupReviewResponse);

    expect(request).toEqual(before);
    expect(studioV2DatalinkApi.post).toHaveBeenCalledWith(
      '/studio-v2/workspace/write-groups/migrations/row-groups/review',
      {
        workspace_id: 'workspace-1',
        expected_workspace_revision: 'workspace-rev-1',
        expected_connector_revision: 'connector-rev-1',
        review_digest: 'row-group-digest-1',
        source_ids: ['legacy-g'],
        confirm_snapshot_conversion: false,
      },
    );
  });

  it('keeps safe typed metadata for row-group revision and validation failures', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce({
      data: {
        success: false,
        error: {
          code: 'revision_mismatch',
          message: 'postgres://secret backend exception',
          retryable: true,
          action: 'reload',
          request_id: 'req-row-group-review-409',
        },
      },
    } as never);

    const error = await studioV2WriteGroupMigrationAPI.reviewRowGroupMigration(reviewRequest)
      .catch((value: unknown) => value);
    expect(error).toEqual(expect.objectContaining({
      name: 'WriteGroupResponseError',
      code: 'revision_mismatch',
      action: 'reload',
      retryable: true,
      request_id: 'req-row-group-review-409',
      message: 'write-group row-group migration review response was invalid',
    }));
    expect(error).not.toHaveProperty('message', 'postgres://secret backend exception');
  });

  it('retains the complete shared-column row-group intent and candidate provenance', () => {
    const parsed = parseWriteGroupMigrationPreviewData(rowGroupPreviewResponse);

    expect(parsed).toEqual(rowGroupPreviewResponse);
    expect(parsed?.items[0]?.before_row_group_intent?.members.map((member) => member.group_key)).toEqual([
      'legacy-g:p-a',
      'legacy-g:p-b',
    ]);
    expect(parsed?.items[0]?.before_row_group_intent?.row_group).toEqual(rowGroupIntent.row_group);
    expect(parsed?.items[0]?.candidate_group?.members.map((member) => member.entity_key)).toEqual([
      'legacy-g:p-a',
      'legacy-g:p-b',
    ]);
    expect(parsed?.items[0]?.candidate_group?.row_policy.group_key_columns).toEqual(['entity_id']);
    expect(parsed?.items[0]?.candidate_group?.migration.target_mapping_points).toEqual({
      'legacy-A': 'p-a',
      'legacy-B': 'p-b',
    });
    expect(parsed?.items[0]?.candidate_group?.migration.source_ids).toEqual(['legacy-A', 'legacy-B']);
    expect(parsed?.items[0]?.candidate_group?.migration.source_ids).not.toContain('legacy-g');
    expect(parsed?.items[1]?.candidate_group).toBeUndefined();
    expect(parsed?.items[1]?.before_row_group_intent?.members[0]?.device_id).toBe('');
    expect(parsed?.items[1]?.before_row_group_intent?.row_group.unique_key_columns).toEqual([]);
  });

  it('requires complete source identity for a reviewable row-group intent', () => {
    const incomplete = structuredClone(rowGroupPreviewResponse) as unknown as Record<string, unknown>;
    const items = incomplete.items as Array<Record<string, unknown>>;
    const intent = items[0]?.before_row_group_intent as Record<string, unknown>;
    const members = intent.members as Array<Record<string, unknown>>;
    members[0].device_id = '';

    expect(parseWriteGroupMigrationPreviewData(incomplete)).toBeNull();
  });

  it('rejects a prototype key in provenance maps instead of rebuilding an unsafe object', () => {
    const unsafe = structuredClone(rowGroupPreviewResponse) as unknown as Record<string, unknown>;
    const items = unsafe.items as Array<Record<string, unknown>>;
    const candidate = items[0]?.candidate_group as Record<string, unknown>;
    const migration = candidate.migration as Record<string, unknown>;
    migration.target_mapping_points = JSON.parse('{"__proto__":"p-a"}') as unknown;

    expect(parseWriteGroupMigrationPreviewData(unsafe)).toBeNull();
  });

  it('rejects row-group arrays and provenance maps beyond safe bounds', () => {
    const oversizedMap = structuredClone(rowGroupPreviewResponse) as unknown as Record<string, unknown>;
    const mapItems = oversizedMap.items as Array<Record<string, unknown>>;
    const mapCandidate = mapItems[0]?.candidate_group as Record<string, unknown>;
    const mapMigration = mapCandidate.migration as Record<string, unknown>;
    mapMigration.target_mapping_points = Object.fromEntries(
      Array.from({ length: 65 }, (_, index) => [`legacy-${index}`, `point-${index}`]),
    );
    expect(parseWriteGroupMigrationPreviewData(oversizedMap)).toBeNull();

    const oversizedArray = structuredClone(rowGroupPreviewResponse) as unknown as Record<string, unknown>;
    const arrayItems = oversizedArray.items as Array<Record<string, unknown>>;
    const arrayIntent = arrayItems[0]?.before_row_group_intent as Record<string, unknown>;
    const arrayRowGroup = arrayIntent.row_group as Record<string, unknown>;
    arrayRowGroup.member_point_ids = Array.from({ length: 257 }, (_, index) => `point-${index}`);
    expect(parseWriteGroupMigrationPreviewData(oversizedArray)).toBeNull();
  });

  it('retains new fields when parsing the canonical list envelope', () => {
    const parsed = parseWriteGroupListData({
      workspace_id: rowGroupReviewResponse.workspace_id,
      workspace_revision: rowGroupReviewResponse.workspace_revision,
      groups: rowGroupReviewResponse.groups,
    });

    expect(parsed?.groups[0]?.members[0]?.entity_key).toBe('legacy-g:p-a');
    expect(parsed?.groups[0]?.row_policy.unique_key_columns).toEqual(['entity_id', 'observed_at']);
    expect(parsed?.groups[0]?.migration.legacy_row_group_id).toBe('legacy-g');
    expect(parsed?.groups[0]?.migration.target_mapping_points).toEqual({
      'legacy-A': 'p-a',
      'legacy-B': 'p-b',
    });
  });

  it('serializes row-group metadata on canonical update without sending server fields', async () => {
    vi.mocked(studioV2DatalinkApi.put).mockResolvedValueOnce(envelope({
      workspace_revision: 'workspace-rev-3',
      group: persistedGroup,
    }) as never);
    const group: WriteGroupDraft = {
      workspace_id: 'workspace-1',
      name: candidateGroup.name,
      members: candidateGroup.members,
      destination: {
        connector_id: candidateGroup.destination.connector_id,
        connector_revision: candidateGroup.destination.connector_revision,
        table_schema: candidateGroup.destination.table_schema,
        table_name: candidateGroup.destination.table_name,
        storage_strategy: candidateGroup.destination.storage_strategy,
      },
      row_policy: candidateGroup.row_policy,
      write_policy: candidateGroup.write_policy,
    };
    const request = {
      workspace_id: 'workspace-1',
      expected_workspace_revision: 'workspace-rev-2',
      expected_group_revision: 'group-revision-1',
      expected_connector_revision: 'connector-rev-1',
      group,
    };

    await expect(studioV2WorkspaceWriteGroupsAPI.update('group-1', request)).resolves.toMatchObject({
      group: persistedGroup,
    });

    const payload = vi.mocked(studioV2DatalinkApi.put).mock.calls[0]?.[1] as {
      group: {
        members: Array<Record<string, unknown>>;
        row_policy: Record<string, unknown>;
      };
    };
    expect(payload.group.members).toEqual(expect.arrayContaining([
      expect.objectContaining({ entity_key: 'legacy-g:p-a' }),
      expect.objectContaining({ entity_key: 'legacy-g:p-b' }),
    ]));
    expect(payload.group.row_policy).toEqual(expect.objectContaining({
      group_key_columns: ['entity_id'],
      unique_key_columns: ['entity_id', 'observed_at'],
    }));
    expect(payload.group).not.toHaveProperty('migration');
    expect(payload.group).not.toHaveProperty('status');
  });
});
