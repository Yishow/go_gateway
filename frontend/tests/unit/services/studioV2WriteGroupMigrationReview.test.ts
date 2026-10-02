import { describe, expect, it } from 'vitest';
import { parseWriteGroupMigrationReviewData } from '@/utils/studioV2WriteGroupMigrationJson';

const canonicalGroup = {
  id: 'group-1',
  workspace_id: 'workspace-1',
  revision: 'group-revision-1',
  applied_revision: '',
  name: 'raw values',
  status: 'draft',
  members: [{
    device_id: 'device-1', point_id: 'point-1', tag_id: 'tag-1',
    source_revision: 'source-revision-1', mapping_revision: 'mapping-revision-1',
    target_column: 'value', required: true,
  }],
  destination: {
    connector_id: 'connector-1', connector_revision: 'connector-revision-1',
    database: 'gateway.db', table_schema: 'main', table_name: 'raw_values', storage_strategy: 'custom',
  },
  row_policy: { interval_seconds: 15, allowed_lateness_seconds: 0 },
  write_policy: {},
  migration: { source_kind: 'legacy-single-mapping', source_ids: ['legacy-1'] },
  created_at: '2026-10-02T00:00:00Z',
  updated_at: '2026-10-02T00:00:00Z',
};

const reviewData = (groups: unknown[]) => ({
  workspace_id: 'workspace-1',
  workspace_revision: 'workspace-revision-1',
  connector_revision: 'connector-revision-1',
  groups,
});

describe('studio V2 migration review canonical response parser', () => {
  it('rejects a successful review with no persisted groups', () => {
    expect(parseWriteGroupMigrationReviewData(reviewData([]))).toBeNull();
  });

  it('rejects duplicate persisted group IDs', () => {
    expect(parseWriteGroupMigrationReviewData(reviewData([
      canonicalGroup,
      { ...canonicalGroup, revision: 'group-revision-2' },
    ]))).toBeNull();
  });
});
