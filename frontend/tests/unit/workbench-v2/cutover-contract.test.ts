import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { parseWriteGroupCutover, parseWriteGroupCutoverReview } from '@/utils/studioV2WriteGroupCutoverJson';
import type { WriteGroupCutoverRequest } from '@/types/studioV2WriteGroupCutover';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}));

const BOUNDARY = '2026-10-24T04:50:22.500123Z';
const member = { device_id: 'dev-1', point_id: 'pt-1', tag_id: 'tag-1' };
const view = (id: string, layout: string) => ({
  group_id: id, name: id, revision: `${id}-rev`, layout, interval_seconds: 15, connector_id: 'conn-1',
  table_schema: 'main', table_name: 'tbl', member_count: 1, schema_confirmed: true,
});
const review = (extra: Record<string, unknown> = {}) => ({
  kind: 'forward', boundary: BOUNDARY, digest: 'a'.repeat(64), applicable: true, issues: [],
  source: view('wide-1', 'managed_wide_v1'), targets: [view('long-1', 'managed_long_v1')],
  moved: [member], unselected: [], remaining: [], ...extra,
});
const request: WriteGroupCutoverRequest = {
  workspace_id: 'ws-1', expected_workspace_revision: 'w1', expected_source_revision: 'wide-1-rev',
  expected_connector_revision: 'crev-1', targets: [{ group_id: 'long-1', expected_revision: 'long-1-rev' }], boundary: BOUNDARY,
};

describe('cutover review contract', () => {
  it('keeps the reviewed instant exactly, with its microseconds', () => {
    expect(parseWriteGroupCutoverReview(review())?.boundary).toBe(BOUNDARY);
  });

  it('shows which sources move and which were not chosen', () => {
    const other = { device_id: 'dev-1', point_id: 'pt-2', tag_id: 'tag-2' };
    const parsed = parseWriteGroupCutoverReview(review({ applicable: false, issues: ['unselected_sources'], unselected: [other] }));
    expect(parsed).toMatchObject({ applicable: false, issues: ['unselected_sources'], moved: [member], unselected: [other] });
  });

  it('keeps the remaining wide revision a partial cutover leaves recording', () => {
    const other = { device_id: 'dev-1', point_id: 'pt-2', tag_id: 'tag-2' };
    const parsed = parseWriteGroupCutoverReview(review({ unselected: [other], remaining: [other] }));
    expect(parsed).toMatchObject({ applicable: true, unselected: [other], remaining: [other] });
    expect(parseWriteGroupCutoverReview(review({
      applicable: false, issues: ['remaining_draft_mismatch'], unselected: [other], remaining: [],
    }))).toMatchObject({ issues: ['remaining_draft_mismatch'] });
    expect(parseWriteGroupCutoverReview(review({
      applicable: false, issues: ['remaining_draft_omits_required_column', 'restored_group_not_ready'],
    }))).toMatchObject({ applicable: false, issues: ['remaining_draft_omits_required_column', 'restored_group_not_ready'] });
  });

  it.each([
    ['a digest that is not a SHA-256', { digest: 'abc' }],
    ['applicable with issues', { applicable: true, issues: ['boundary_not_future'] }],
    ['not applicable without a reason', { applicable: false, issues: [] }],
    ['an unknown issue', { applicable: false, issues: ['made_up_issue'] }],
    ['an unknown kind', { kind: 'sideways' }],
    ['an instant that is not a time', { boundary: 'tomorrow' }],
    ['a moved source without its tag', { moved: [{ device_id: 'dev-1', point_id: 'pt-1' }] }],
    ['no sources list', { moved: undefined }],
    ['no remaining list', { remaining: undefined }],
    ['a target that is not a group view', { targets: ['long-1'] }],
  ])('fails closed on %s', (_name, extra) => {
    expect(parseWriteGroupCutoverReview(review(extra))).toBeNull();
  });
});

describe('applied cutover contract', () => {
  const applied = (extra: Record<string, unknown> = {}) => ({
    id: 'cut-1', workspace_id: 'ws-1', kind: 'forward', source_group_id: 'wide-1', boundary: BOUNDARY,
    review_digest: 'a'.repeat(64), ended: [{ group_id: 'wide-1', revision: 'wide-1-rev' }],
    started: [{ group_id: 'long-1', revision: 'long-1-rev' }], moved: [member], created_at: '2026-10-24T04:49:00Z', ...extra,
  });

  it('keeps what ended and what started at the boundary', () => {
    expect(parseWriteGroupCutover(applied())).toMatchObject({ id: 'cut-1', boundary: BOUNDARY, ended: [{ group_id: 'wide-1' }], started: [{ group_id: 'long-1' }] });
  });

  it('keeps the cutover a rollback reverses', () => {
    expect(parseWriteGroupCutover(applied({ kind: 'rollback', prior_cutover_id: 'cut-0' }))).toMatchObject({ kind: 'rollback', prior_cutover_id: 'cut-0' });
  });

  it.each([
    ['an unknown kind', { kind: 'both' }],
    ['a rollback that names no cutover', { kind: 'rollback' }],
    ['a forward cutover that names a prior one', { prior_cutover_id: 'cut-0' }],
    ['an end without a revision', { ended: [{ group_id: 'wide-1' }] }],
  ])('fails closed on %s', (_name, extra) => {
    expect(parseWriteGroupCutover(applied(extra))).toBeNull();
  });
});

describe('cutover requests', () => {
  beforeEach(() => {
    vi.mocked(studioV2DatalinkApi.post).mockReset();
    vi.mocked(studioV2DatalinkApi.get).mockReset();
  });

  it('reads one recorded cutover and the list of a source', async () => {
    const record = {
      id: 'cut-1', workspace_id: 'ws-1', kind: 'forward', source_group_id: 'wide-1', boundary: BOUNDARY, review_digest: 'a'.repeat(64),
      ended: [{ group_id: 'wide-1', revision: 'r' }], started: [{ group_id: 'long-1', revision: 'r2' }], moved: [member], created_at: '2026-10-24T04:49:00Z',
    };
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success: true, data: record } } as never);
    expect((await studioV2WorkspaceWriteGroupsAPI.getCutover('wide-1', 'cut-1')).boundary).toBe(BOUNDARY);
    expect(vi.mocked(studioV2DatalinkApi.get).mock.calls[0][0]).toBe('/studio-v2/workspace/write-groups/wide-1/cutovers/cut-1');
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success: true, data: [record] } } as never);
    expect(await studioV2WorkspaceWriteGroupsAPI.listCutovers('wide-1')).toHaveLength(1);
    vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success: true, data: [{ ...record, kind: 'both' }] } } as never);
    await expect(studioV2WorkspaceWriteGroupsAPI.listCutovers('wide-1')).rejects.toBeTruthy();
  });

  it('reviews with exactly the reviewed values and sends the instant as text', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValue({ data: { success: true, data: review() } } as never);
    const result = await studioV2WorkspaceWriteGroupsAPI.reviewCutover('wide-1', { ...request, extra: 'ignored' } as WriteGroupCutoverRequest);
    expect(result.digest).toBe('a'.repeat(64));
    const [path, body] = vi.mocked(studioV2DatalinkApi.post).mock.calls[0];
    expect(path).toBe('/studio-v2/workspace/write-groups/wide-1/cutover-review');
    expect(body).toEqual({ ...request, kind: 'forward' });
    expect((body as WriteGroupCutoverRequest).boundary).toBe(BOUNDARY);
  });

  it('does not apply without the reviewed instant and digest, and sends nothing', async () => {
    await expect(studioV2WorkspaceWriteGroupsAPI.applyCutover('wide-1', request)).rejects.toThrow(/review/i);
    await expect(studioV2WorkspaceWriteGroupsAPI.applyCutover('wide-1', { ...request, boundary: undefined, review_digest: 'a'.repeat(64) })).rejects.toThrow(/review/i);
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled();
  });

  it('applies a reviewed cutover and returns its record', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValue({
      data: { success: true, data: {
        id: 'cut-1', workspace_id: 'ws-1', kind: 'forward', source_group_id: 'wide-1', boundary: BOUNDARY, review_digest: 'a'.repeat(64),
        ended: [{ group_id: 'wide-1', revision: 'r' }], started: [{ group_id: 'long-1', revision: 'r2' }], moved: [member], created_at: '2026-10-24T04:49:00Z',
      } },
    } as never);
    const result = await studioV2WorkspaceWriteGroupsAPI.applyCutover('wide-1', { ...request, review_digest: 'a'.repeat(64) });
    expect(result.id).toBe('cut-1');
    expect(vi.mocked(studioV2DatalinkApi.post).mock.calls[0][0]).toBe('/studio-v2/workspace/write-groups/wide-1/cutover');
  });
});
