import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { parseWriteGroupDeliveryData } from '@/utils/studioV2WriteGroupDeliveryJson';
import type { WriteGroupDelivery } from '@/types/studioV2WriteGroupDelivery';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}));

function delivery(overrides: Partial<WriteGroupDelivery> = {}): WriteGroupDelivery {
  return {
    group_id: 'group-1',
    intake: { state: 'active' },
    stages: { collecting: 2, queued: 3, retrying: 1, blocked: 1, quarantined: 0, unknown: 1, sql_committed: 5, skipped: 0 },
    last_sql_committed_at: '2026-10-02T12:00:10Z',
    oldest_pending_seconds: 42.5,
    no_data_buckets: 1,
    skipped_buckets: 2,
    backlog: [{
      group_revision: 'rev-1', connector_id: 'connector-A', connector_revision: 'connector-1',
      table_schema: 'main', table_name: 'raw_values', pending: 4, error_codes: ['target-blocked'],
    }],
    quota: { configured: true, state: 'ok', scope: 'global', used_bytes: 10, max_bytes: 100, intake_refused: false },
    ...overrides,
  };
}

describe('write group delivery status', () => {
  beforeEach(() => vi.clearAllMocks());

  it('preserves actual revision-bound effect identity and rejects contradictory evidence', () => {
    const effect = {
      group_revision: 'rev-2', connector_revision: 'crev-2', record_id: 'record-2',
      effect_key: 'effect-2', payload_digest: 'a'.repeat(64), committed_at: '2026-10-02T12:00:10Z',
    };
    expect(parseWriteGroupDeliveryData({ ...delivery(), last_sql_committed_effect: effect }))
      .toMatchObject({ last_sql_committed_effect: effect });
    expect(parseWriteGroupDeliveryData({ ...delivery(), last_sql_committed_effect: null }))
      .toMatchObject({ last_sql_committed_effect: null });
    expect(parseWriteGroupDeliveryData({ ...delivery(), last_sql_committed_effect: { ...effect, payload_digest: 'bad' } })).toBeNull();
    expect(parseWriteGroupDeliveryData({ ...delivery(), last_sql_committed_effect: { ...effect, record_id: '' } })).toBeNull();
    expect(parseWriteGroupDeliveryData({ ...delivery(), last_sql_committed_effect: { ...effect, committed_at: 'bad' } })).toBeNull();
    expect(parseWriteGroupDeliveryData({
      ...delivery(), stages: { ...delivery().stages, sql_committed: 0 }, last_sql_committed_at: null,
      last_sql_committed_effect: effect,
    })).toBeNull();
  });

  it('keeps current revision stages separate from historical totals', () => {
    const current = { group_revision: 'rev-2', stages: {
      collecting: 1, queued: 1, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0,
    } };
    expect(parseWriteGroupDeliveryData({ ...delivery(), revision_stages: current }))
      .toMatchObject({ revision_stages: current });
    expect(parseWriteGroupDeliveryData({ ...delivery(), revision_stages: { ...current, group_revision: '' } })).toBeNull();
    expect(parseWriteGroupDeliveryData({ ...delivery(), revision_stages: { ...current, stages: { ...current.stages, queued: -1 } } })).toBeNull();
  });

  it('reads the delivery route for the encoded group id and keeps every stage', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success: true, data: delivery() } } as never);
    const result = await studioV2WorkspaceWriteGroupsAPI.delivery('group/with space');
    expect(studioV2DatalinkApi.get).toHaveBeenCalledWith('/studio-v2/workspace/write-groups/group%2Fwith%20space/delivery');
    expect(result.stages.sql_committed).toBe(5);
    expect(result.backlog[0].connector_revision).toBe('connector-1');
  });

  it('rejects a committed time without committed rows and committed rows without a time', () => {
    const buffered = delivery({
      stages: { collecting: 0, queued: 4, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 },
    });
    expect(parseWriteGroupDeliveryData(buffered)).toBeNull();
    expect(parseWriteGroupDeliveryData({ ...buffered, last_sql_committed_at: null })).not.toBeNull();
    expect(parseWriteGroupDeliveryData(delivery({ last_sql_committed_at: null }))).toBeNull();
  });

  it('rejects malformed counts, unknown states and contradictory quota', () => {
    expect(parseWriteGroupDeliveryData(delivery({ stages: { ...delivery().stages, queued: -1 } }))).toBeNull();
    expect(parseWriteGroupDeliveryData(delivery({ stages: { ...delivery().stages, queued: 1.5 } }))).toBeNull();
    expect(parseWriteGroupDeliveryData({ ...delivery(), intake: { state: 'running' } })).toBeNull();
    expect(parseWriteGroupDeliveryData({ ...delivery(), stages: { collecting: 1 } })).toBeNull();
    expect(parseWriteGroupDeliveryData(delivery({ oldest_pending_seconds: Number.NaN }))).toBeNull();
    expect(parseWriteGroupDeliveryData(delivery({
      quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: true },
    }))).toBeNull();
    expect(parseWriteGroupDeliveryData(delivery({
      quota: { configured: true, state: 'ok', scope: 'global', used_bytes: 1, max_bytes: 2, intake_refused: true },
    }))).toBeNull();
    expect(parseWriteGroupDeliveryData(delivery({
      backlog: [{ group_revision: 'rev-1', connector_id: 'c', connector_revision: 'r', table_schema: '', table_name: 't', pending: -3, error_codes: [] }],
    }))).toBeNull();
  });

  it('accepts an unconfigured quota, a hard limit with a loss-risk notice and an unknown group', () => {
    expect(parseWriteGroupDeliveryData(delivery({
      quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
    }))).not.toBeNull();
    const hard = parseWriteGroupDeliveryData(delivery({
      quota: {
        configured: true, state: 'hard_limit', scope: 'group', used_bytes: 95, max_bytes: 100,
        intake_refused: true, loss_risk_notice: 'new readings are not recorded',
      },
    }));
    expect(hard?.quota.loss_risk_notice).toBe('new readings are not recorded');
    expect(parseWriteGroupDeliveryData({
      ...delivery({ intake: { state: 'not_running' }, last_sql_committed_at: null, backlog: [] }),
      stages: { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 },
    })).not.toBeNull();
  });

  it('preserves recent bucket issues while reducing unknown causes to a safe fallback', () => {
    const parsed = parseWriteGroupDeliveryData({
      ...delivery(),
      recent_bucket_issues: [{
        group_revision: 'rev-1', bucket_start: '2026-10-03T00:00:00Z', kind: 'skipped',
        causes: ['missing', 'private-dsn'],
      }],
    });

    expect(parsed?.recent_bucket_issues).toEqual([{
      group_revision: 'rev-1', bucket_start: '2026-10-03T00:00:00Z', kind: 'skipped',
      causes: ['missing', 'unavailable'],
    }]);
    expect(parseWriteGroupDeliveryData({
      ...delivery(),
      recent_bucket_issues: Array.from({ length: 21 }, (_, index) => ({
        group_revision: 'rev-1', bucket_start: `2026-10-03T00:${String(index).padStart(2, '0')}:00Z`,
        kind: 'no_data', causes: ['no_data'],
      })),
    })).toBeNull();
  });

  it('turns a malformed envelope into a typed response error instead of empty data', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success: true, data: { group_id: 'x' } } } as never);
    await expect(studioV2WorkspaceWriteGroupsAPI.delivery('x')).rejects.toBeTruthy();
  });
});
