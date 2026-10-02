import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WriteGroupMigrationAPI } from '@/services/studioV2WriteGroupMigration';
import type { WriteGroupMigrationReviewRequest } from '@/types/studioV2WriteGroupMigration';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

const recordingPlan = {
  id: 'legacy-plan-A',
  workspace_id: 'workspace-1',
  revision: '',
  applied_revision: 'applied-revision-1',
  name: 'Original every-sample plan',
  status: 'running',
  timezone: 'Asia/Taipei',
  members: null,
  streams: [{
    stream_id: 'stream-A',
    measurement_id: 'measurement-A',
    mode: 'legacy_future_mode',
    raw_policy: 'every_sample',
    destination_ids: ['destination-A'],
  }],
  destinations: null,
  retention: {
    raw_days: 12,
    summary_days: 30,
    events_days: 60,
    correction_horizon_hours: 8,
  },
  limits: {
    max_batch_size: 7,
    max_hold_seconds: 9,
    max_queue_bytes: 4096,
  },
  created_at: '2026-10-02T02:00:00Z',
  updated_at: '2026-10-02T02:00:00Z',
  legacy_raw_field: { mode: 'keep-this-for-repair' },
};

const beforeRecordingPlanIntent = {
  source_id: 'legacy-plan-A',
  plan: recordingPlan,
  sources: [{
    measurement_id: 'measurement-A',
    device_id: 'device-A',
    point_id: 'point-A',
    tag_id: 'tag-A',
    definition_revision: 'definition-1',
    source_binding_revision: 'binding-1',
    series_epoch: 'epoch-1',
    source_revision: 'source-revision-1',
    mapping_revision: 'mapping-revision-1',
    status: 'resolved',
  }],
};

const previewResponse = {
  workspace_id: 'workspace-1',
  workspace_revision: 'workspace-revision-1',
  connector_revision: 'connector-revision-1',
  adapter_version: 'recording-plan-v1',
  review_digest: 'recording-plan-digest-1',
  items: [{
    source_id: 'legacy-plan-A',
    source_revision: 'plan-source-revision-1',
    status: 'blocked',
    differences: [],
    issues: [{
      code: 'missing-snapshot-identity',
      message: 'the original plan has no canonical periodic snapshot identity',
    }],
    before_recording_plan_intent: beforeRecordingPlanIntent,
    repair_action: 'open_write_groups',
  }],
};

const reviewRequest: WriteGroupMigrationReviewRequest & Record<string, unknown> = {
  workspace_id: 'workspace-1',
  expected_workspace_revision: 'workspace-revision-1',
  expected_connector_revision: 'connector-revision-1',
  review_digest: 'recording-plan-digest-1',
  source_ids: ['legacy-plan-A'],
  confirm_snapshot_conversion: false,
  apply: true,
  group: { id: 'forged-group' },
  dsn: 'postgres://secret',
  server_readonly: false,
};

function envelope<T>(data: T) {
  return { data: { success: true, data } };
}

describe('studio V2 recording-plan migration service', () => {
  beforeEach(() => {
    vi.mocked(studioV2DatalinkApi.post).mockReset();
  });

  it('previews a blocked plan while preserving the complete raw intent', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope(previewResponse) as never);

    const request = {
      workspace_id: 'workspace-1',
      source_ids: ['legacy-plan-A'],
      apply: true,
      dsn: 'postgres://secret',
    };
    const before = structuredClone(request);

    await expect(studioV2WriteGroupMigrationAPI.previewRecordingPlanMigration(request)).resolves.toEqual(previewResponse);
    expect(request).toEqual(before);
    expect(studioV2DatalinkApi.post).toHaveBeenCalledWith(
      '/studio-v2/workspace/write-groups/migrations/recording-plans/preview',
      { workspace_id: 'workspace-1', source_ids: ['legacy-plan-A'] },
    );
  });

  it('sends only the reviewed fields and preserves an explicit false confirmation', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce({
      data: {
        success: false,
        error: {
          code: 'WRITE_GROUP_PLAN_MIGRATION_BLOCKED',
          action: 'open_write_groups',
          request_id: 'request-recording-plan-422',
          message: 'private backend detail',
        },
      },
    } as never);

    const request = structuredClone(reviewRequest);
    const before = structuredClone(request);
    await expect(studioV2WriteGroupMigrationAPI.reviewRecordingPlanMigration(request)).rejects.toMatchObject({
      name: 'WriteGroupResponseError',
      code: 'WRITE_GROUP_PLAN_MIGRATION_BLOCKED',
      action: 'open_write_groups',
      request_id: 'request-recording-plan-422',
    });
    expect(request).toEqual(before);
    expect(studioV2DatalinkApi.post).toHaveBeenCalledWith(
      '/studio-v2/workspace/write-groups/migrations/recording-plans/review',
      {
        workspace_id: 'workspace-1',
        expected_workspace_revision: 'workspace-revision-1',
        expected_connector_revision: 'connector-revision-1',
        review_digest: 'recording-plan-digest-1',
        source_ids: ['legacy-plan-A'],
        confirm_snapshot_conversion: false,
      },
    );
  });

  it('retains unresolved source fields as bounded empty strings', async () => {
    const response = structuredClone(previewResponse) as {
      items: Array<{ before_recording_plan_intent: { sources: Array<Record<string, unknown>> } }>;
    };
    const source = response.items[0]!.before_recording_plan_intent.sources[0]!;
    source.device_id = '';
    source.point_id = '';
    source.tag_id = '';
    source.definition_revision = '';
    source.source_binding_revision = '';
    source.series_epoch = '';
    source.status = 'blocked';
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope(response) as never);

    const parsed = await studioV2WriteGroupMigrationAPI.previewRecordingPlanMigration({
      workspace_id: 'workspace-1',
      source_ids: ['legacy-plan-A'],
    });
    expect(parsed.items[0]?.before_recording_plan_intent.sources[0]).toMatchObject({
      device_id: '',
      point_id: '',
      tag_id: '',
      definition_revision: '',
      source_binding_revision: '',
      series_epoch: '',
      status: 'blocked',
    });
  });

  it.each([
    {
      name: 'claims needs review',
      response: {
        ...previewResponse,
        items: [{ ...previewResponse.items[0], status: 'needs_review' }],
      },
    },
    {
      name: 'includes a candidate group',
      response: {
        ...previewResponse,
        items: [{ ...previewResponse.items[0], candidate_group: { id: 'forged' } }],
      },
    },
    {
      name: 'uses a different repair action',
      response: {
        ...previewResponse,
        items: [{ ...previewResponse.items[0], repair_action: 'apply' }],
      },
    },
    {
      name: 'does not bind the raw plan to the item source',
      response: {
        ...previewResponse,
        items: [{
          ...previewResponse.items[0],
          before_recording_plan_intent: {
            ...previewResponse.items[0].before_recording_plan_intent,
            source_id: 'other-plan',
          },
        }],
      },
    },
  ])('rejects a forged recording-plan preview that $name', async ({ response }) => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope(response) as never);
    await expect(studioV2WriteGroupMigrationAPI.previewRecordingPlanMigration({
      workspace_id: 'workspace-1',
      source_ids: ['legacy-plan-A'],
    })).rejects.toMatchObject({
      name: 'WriteGroupResponseError',
      message: 'write-group recording-plan migration preview response was invalid',
    });
  });

  it('rejects prototype keys in the preserved raw plan', async () => {
    const response = structuredClone(previewResponse) as Record<string, unknown>;
    const items = response.items as Array<Record<string, unknown>>;
    const intent = items[0]!.before_recording_plan_intent as Record<string, unknown>;
    const plan = intent.plan as Record<string, unknown>;
    plan.legacy_raw_field = JSON.parse('{"__proto__":"unsafe"}') as unknown;
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope(response) as never);

    await expect(studioV2WriteGroupMigrationAPI.previewRecordingPlanMigration({
      workspace_id: 'workspace-1',
      source_ids: ['legacy-plan-A'],
    })).rejects.toMatchObject({ name: 'WriteGroupResponseError' });
  });

  it('rejects a nominal recording-plan review success because this adapter has no candidate', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce(envelope({
      workspace_id: 'workspace-1',
      workspace_revision: 'workspace-revision-2',
      connector_revision: 'connector-revision-1',
      groups: [{ id: 'forged-canonical-group' }],
    }) as never);

    await expect(studioV2WriteGroupMigrationAPI.reviewRecordingPlanMigration(reviewRequest)).rejects.toMatchObject({
      name: 'WriteGroupResponseError',
      message: 'write-group recording-plan migration review response was invalid',
    });
  });
});
