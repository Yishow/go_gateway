import { beforeEach, describe, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WriteGroupSuccessorAPI as api } from '@/services/studioV2WriteGroupSuccessor';
import { activation, envelope, history, preparation } from '../../fixtures/writeGroupSuccessor';
import { DIGEST, recordOf } from '../../fixtures/cutoverFixtures';
import { managedSchemaPreviewToken } from '../../fixtures/managedSchemaPreview';
vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { post: vi.fn() } }));
const reply = (data: unknown) => vi.mocked(studioV2DatalinkApi.post).mockResolvedValueOnce({ data: { success: true, data } });
const reserved = () => ({ envelope: envelope(), review_digest: DIGEST });
const confirmed = () => ({ envelope: activation(), review_digest: DIGEST, readiness_digest: DIGEST });
const activationReview = () => ({ envelope: activation(), review_digest: DIGEST, history: history(), blockers: [] });
const schemaOperation = () => ({ operation_id: 'schema-op-1', status: 'succeeded', executed_statements: 1, verified_digest: DIGEST, created_at: '2026-10-28T00:00:00Z', updated_at: '2026-10-28T00:00:00Z' });
beforeEach(() => vi.resetAllMocks());
describe('successor stage routes and proof binding', () => {
  it('keeps review local route and reservation explicit', async () => {
    const review = { ...reserved(), history: history(), blockers: [], required_checks: [] }; reply(review);
    expect(await api.review('wide-1', envelope().request)).toEqual(review);
    expect(studioV2DatalinkApi.post).toHaveBeenLastCalledWith('/studio-v2/workspace/write-groups/wide-1/successor-review', envelope().request);
    reply(preparation()); expect(await api.reserve('wide-1', reserved())).toEqual(preparation());
    expect(studioV2DatalinkApi.post).toHaveBeenLastCalledWith('/studio-v2/workspace/write-groups/wide-1/successor', { ...reserved(), action: 'reserve' });
  });
  it('separates schema preview and exact token/digest confirmation', async () => {
    const stage = activation().request.stage;
    const preview = { ...managedSchemaPreviewToken, workspace_id: 'ws-1', workspace_revision: 'ws-r3', plan_id: 'new-2', plan_revision: 'g-r2', connector_id: 'conn-1', connector_revision: 'crev-1' };
    reply(preview); expect(await api.schemaPreview('wide-1', stage)).toEqual(preview);
    expect(studioV2DatalinkApi.post).toHaveBeenLastCalledWith('/studio-v2/workspace/write-groups/wide-1/successor', { ...stage, action: 'schema-preview' });
    const request = { ...stage, token: preview.token, digest: preview.digest, schema_operation_id: 'schema-op-1' };
    reply(schemaOperation()); expect(await api.schemaConfirm('wide-1', request)).toEqual(schemaOperation());
    expect(studioV2DatalinkApi.post).toHaveBeenLastCalledWith('/studio-v2/workspace/write-groups/wide-1/successor', { ...request, action: 'schema-confirm' });
  });
  it('requires verified digest for a purported schema success', async () => {
    reply({ ...schemaOperation(), verified_digest: '' });
    await expect(api.schemaConfirm('wide-1', { ...activation().request.stage, token: 't', digest: DIGEST, schema_operation_id: 'schema-op-1' })).rejects.toMatchObject({ outcome: 'unconfirmed' });
  });
  it('reviews activation then reads readiness without starting intake', async () => {
    reply(activationReview()); expect(await api.activationReview('wide-1', activation().request)).toEqual(activationReview());
    expect(studioV2DatalinkApi.post).toHaveBeenLastCalledWith('/studio-v2/workspace/write-groups/wide-1/successor', { ...activation().request, action: 'activation-review' });
    const readiness = { review: activationReview(), readiness: { workspace_id: 'ws-1', workspace_revision: 'ws-r3', group_id: 'new-2', group_revision: 'g-r2', applied_revision: '', config_ready: true, schema_ready: true, ready: true, schema_digest: DIGEST, issues: [] }, readiness_digest: DIGEST };
    reply(readiness); expect(await api.activationReadiness('wide-1', { envelope: activation(), review_digest: DIGEST })).toEqual(readiness);
  });
  it('rejects a changed activation envelope even if the returned digest text was copied', async () => {
    const changed = activationReview(); changed.envelope.request.stage.current.lineage_revision++;
    reply({ review: changed, readiness: { workspace_id: 'ws-1', workspace_revision: 'ws-r3', group_id: 'new-2', group_revision: 'g-r2', applied_revision: '', config_ready: true, schema_ready: true, ready: true, schema_digest: DIGEST, issues: [] }, readiness_digest: DIGEST });
    await expect(api.activationReadiness('wide-1', { envelope: activation(), review_digest: DIGEST })).rejects.toMatchObject({ outcome: 'unconfirmed' });
  });
  it('binds cutover to exact reviewed source and target revisions and digest', async () => {
    const cutover = recordOf({ started: [{ group_id: 'new-2', revision: 'g-r2' }] }); reply(cutover);
    expect(await api.confirmCutover('wide-1', confirmed())).toEqual(cutover);
    reply({ ...cutover, review_digest: 'b'.repeat(64) });
    await expect(api.confirmCutover('wide-1', confirmed())).rejects.toMatchObject({ outcome: 'unconfirmed' });
  });
  it('routes v2 first-start only to first-start-confirm, never legacy cutover', async () => {
    const confirmation = confirmed(); confirmation.envelope = { ...confirmation.envelope, version: 'successor-activation-review-v2', activation_kind: 'first_start', source_applied_revision: '' };
    confirmation.envelope.request.stage.envelope.source_applied_revision = '';
    const first = { operation_id: 'first-op', root_operation_id: 'root-op-1', workspace_revision: 'ws-r4', lineage_id: 'lineage-1', lineage_revision: 3, predecessor_group_id: 'wide-1', group_id: 'new-2', applied_revision: 'g-r2', boundary: activation().request.boundary, generation: 2, status: 'succeeded' };
    reply(first); expect(await api.confirmFirstStart('wide-1', confirmation)).toEqual(first);
    expect(studioV2DatalinkApi.post).toHaveBeenLastCalledWith('/studio-v2/workspace/write-groups/wide-1/successor', { ...confirmation, action: 'first-start-confirm' });
    for (const invalid of [{ ...first, applied_revision: 'wrong-revision' }, { ...first, generation: 99 }, { ...first, lineage_revision: 99 }]) {
      reply(invalid);
      await expect(api.confirmFirstStart('wide-1', confirmation)).rejects.toMatchObject({ outcome: 'unconfirmed' });
    }
    vi.mocked(studioV2DatalinkApi.post).mockClear();
    await expect(api.confirmCutover('wide-1', confirmation)).rejects.toMatchObject({ outcome: 'unconfirmed' });
    await expect(api.confirmFirstStart('wide-1', confirmed())).rejects.toMatchObject({ outcome: 'unconfirmed' });
    expect(studioV2DatalinkApi.post).not.toHaveBeenCalled();
  });
});
