import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI as groupsApi } from '@/services/studioV2WorkspaceWriteGroups';
import { parseBoundedJson } from '@/utils/safeJson';
import { parseWriteGroupDeliveryData } from '@/utils/studioV2WriteGroupDeliveryJson';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}));

function source(id: string) {
  return {
    circuit_id: 'cir-1', device_id: 'device-1', point_id: `point-${id}`, tag_id: `tag-${id}`,
    latest: {
      available: true, delivery: 'sql_committed', group_revision: 'rev-1', connector_revision: 'crev-1',
      record_id: `record-${id}`, effect_key: `effect-${id}`, payload_digest: 'a'.repeat(64),
      profile_kind: 'instantaneous', window_start: '2026-10-25T04:03:45Z', window_end: '2026-10-25T04:04:00Z',
      status: 'ok', raw_present: true,
      aggregation: {
        method: 'arithmetic-mean', version: 'arithmetic-mean-v1', basis: 'exact',
        source_type: 'uint16', result_type: 'decimal', valid_count: 15, rejected_count: 0,
        coverage: 'unknown', reason: 'coverage-unverified',
      },
      committed_at: '2026-10-25T04:04:01.646429Z',
    },
    last_sql_committed: null,
  };
}

function meanEnvelope(extraDepth = false) {
  const first = source('a');
  if (extraDepth) {
    (first.latest.aggregation as Record<string, unknown>).extra = { note: 'deeper than the delivery contract allows' };
  }
  return {
    success: true,
    data: {
      attention: [], group_id: 'group-1', intake: { state: 'active' },
      stages: { collecting: 14, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 36, skipped: 0 },
      last_sql_committed_at: '2026-10-25T04:04:01.646429Z',
      last_sql_committed_effect: {
        group_revision: 'rev-1', connector_revision: 'crev-1', record_id: 'record-a', effect_key: 'effect-a',
        payload_digest: 'a'.repeat(64), committed_at: '2026-10-25T04:04:01.646429Z',
      },
      oldest_pending_seconds: 0, no_data_buckets: 0, skipped_buckets: 0,
      quota: { configured: false, state: 'unconfigured', scope: 'group', used_bytes: 0, max_bytes: 0, intake_refused: false },
      backlog: [],
      source_evidence: {
        state: 'available', group_revision: 'rev-1', connector_id: 'connector-1', connector_revision: 'crev-1',
        sources: [first, source('b')],
      },
    },
  };
}

describe('write group delivery envelope depth', () => {
  beforeEach(() => vi.mocked(studioV2DatalinkApi.get).mockReset());

  it('reads a mean source evidence envelope through the delivery service, not only the data parser', async () => {
    // The data parser itself already accepts the shape; only the service envelope is deeper.
    expect(parseWriteGroupDeliveryData(meanEnvelope().data)).not.toBeNull();
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue({ data: meanEnvelope() } as never);
    const delivery = await groupsApi.delivery('group-1');
    expect(delivery.source_evidence?.state).toBe('available');
    const latest = delivery.source_evidence?.sources[0].latest;
    expect(latest && 'aggregation' in latest ? latest.aggregation?.coverage : undefined).toBe('unknown');
  });

  it('still rejects an envelope deeper than the delivery contract', async () => {
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValue({ data: meanEnvelope(true) } as never);
    await expect(groupsApi.delivery('group-1')).rejects.toThrow();
  });

  it('keeps the default safe JSON depth at six for every other caller', () => {
    expect(parseBoundedJson({ a: { b: { c: { d: { e: { f: 1 } } } } } })).not.toBeNull();
    expect(parseBoundedJson({ a: { b: { c: { d: { e: { f: { g: 1 } } } } } } })).toBeNull();
  });
});
