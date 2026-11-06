import { describe, expect, it } from 'vitest';
import { createRecordingHomeSnapshotAPI, SNAPSHOT_BASE, SNAPSHOT_SOURCE } from '../../e2e/fixtures/recording-home-snapshot-api';

describe('routed recording home mock API contract', () => {
  it('rejects unknown API reads and all unapproved writes instead of catch-all success', () => {
    const api = createRecordingHomeSnapshotAPI();
    expect(api.handle({ method: 'GET', url: `${SNAPSHOT_BASE}/not-supported` }).status).toBe(404);
    expect(api.handle({ method: 'POST', url: `${SNAPSHOT_BASE}/activate` }).status).toBe(404);
    expect(api.ledger).toHaveLength(2);
  });
  it('provides a scoped compact summary and actual singleton bootstrap, with no raw source configuration', () => {
    const api = createRecordingHomeSnapshotAPI();
    expect(api.handle({ method: 'GET', url: SNAPSHOT_BASE })).toMatchObject({ status: 200, body: { data: { id: 'A' } } });
    const reply = api.handle({ method: 'GET', url: `${SNAPSHOT_SOURCE}?workspace_id=A` });
    expect(reply).toMatchObject({ status: 200, body: { success: true, data: { workspace_id: 'A', state: 'available' } } });
    expect(JSON.stringify(reply)).not.toContain('connection_config');
  });
});

// Held replies must model a request-time response, not a read of state at release time.
describe('snapshot transport fences', () => {
  it('captures independent A / B / A2 response bodies before releasing old replies', async () => {
    const api = createRecordingHomeSnapshotAPI();
    const heldA = api.holdNext(), oldA = api.handle({ method: 'GET', url: `${SNAPSHOT_SOURCE}?workspace_id=A` });
    expect((await heldA.entered).workspaceId).toBe('A');
    api.switchWorkspace('B'); const heldB = api.holdNext();
    const oldB = api.handle({ method: 'GET', url: `${SNAPSHOT_SOURCE}?workspace_id=B` });
    api.switchWorkspace('A', 'A2');
    const current = api.handle({ method: 'GET', url: `${SNAPSHOT_SOURCE}?workspace_id=A` });
    heldB.release(); heldA.release(); await Promise.all([oldA.wait, oldB.wait]);
    expect(oldA.body).toMatchObject({ data: { devices: [expect.objectContaining({ name: 'Snapshot A D' }), expect.anything()] } });
    expect(oldB.body).toMatchObject({ data: { workspace_id: 'B' } });
    expect(current.body).toMatchObject({ data: { devices: [expect.objectContaining({ name: 'Snapshot A A2' }), expect.anything()] } });
    expect(oldA.call.response).toEqual(oldA.body);
  });
  it('rejects foreign/stale scope, duplicate compact params and raw home reads', () => {
    const api = createRecordingHomeSnapshotAPI();
    for (const url of [`${SNAPSHOT_SOURCE}?workspace_id=B`, `${SNAPSHOT_BASE}/runtime-context?workspace_id=A&workspace_revision=OLD&device_id=D-A`]) {
      expect(api.handle({ method: 'GET', url }).status).toBe(409);
    }
    expect(api.handle({ method: 'GET', url: `${SNAPSHOT_SOURCE}?workspace_id=A&workspace_id=A` }).status).toBe(400);
    expect(api.handle({ method: 'GET', url: `${SNAPSHOT_BASE}/devices` }).call.unexpected).toBe(true);
    api.setPhase('editor'); expect(api.handle({ method: 'GET', url: `${SNAPSHOT_BASE}/devices` }).status).toBe(200);
  });
  it('checks scoped locator-read parameters, not only ordinary workspace_id', () => {
    const api = createRecordingHomeSnapshotAPI();
    expect(api.handle({ method: 'GET', url: `${SNAPSHOT_BASE}/database-operations/test-writes?locator_workspace_id=B&locator_workspace_revision=B1` }).status).toBe(409);
    expect(api.handle({ method: 'GET', url: `${SNAPSHOT_BASE}/database-operations/test-writes?locator_workspace_id=A&locator_workspace_revision=OLD` }).status).toBe(409);
  });
  it('fails one read then supports independent read-only retry without advancing source state', () => {
    const api = createRecordingHomeSnapshotAPI(), before = api.snapshot();
    api.failNext('sources');
    expect(api.handle({ method: 'GET', url: `${SNAPSHOT_SOURCE}?workspace_id=A` }).status).toBe(503);
    expect(api.handle({ method: 'GET', url: `${SNAPSHOT_SOURCE}?workspace_id=A` }).status).toBe(200);
    expect(api.snapshot()).toEqual(before); expect(api.ledger.every(call => !call.mutation)).toBe(true);
  });
  it('a one-use prefix-save grant rejects unapproved geometry edits and permits the exact saved ACK', () => {
    const api = createRecordingHomeSnapshotAPI(), before = api.snapshot();
    api.setPhase('source-modal'); api.allowRuleSave('R-A', 'RENAMED_');
    const body = { device_id: 'D-A', start_address: '40001', count: 2, data_type: 'int16', enabled: false,
      naming_prefix: 'RENAMED_', skipped_addresses: [], scale_multiplier: 1, scale_offset: 0, data_format: '',
      share_enabled: false, share_start_register: null, share_stride: null, layout_mode: 'homogeneous' };
    expect(api.handle({ method: 'PUT', url: `${SNAPSHOT_BASE}/source-rules/R-A`, body: { ...body, scale_multiplier: 2 } }).status).toBe(400);
    const reply = api.handle({ method: 'PUT', url: `${SNAPSHOT_BASE}/source-rules/R-A`, body });
    expect(reply).toMatchObject({ status: 200, body: { data: { naming_prefix: 'RENAMED_', runtime_apply_status: 'not_running' } } });
    expect(api.snapshot().rules[0].revision_id).not.toBe(before.rules[0].revision_id);
    expect(api.handle({ method: 'PUT', url: `${SNAPSHOT_BASE}/source-rules/R-A`, body }).status).toBe(404);
  });
  it('releaseAll also releases an already-entered hold during failed-test teardown', async () => {
    const api = createRecordingHomeSnapshotAPI(), hold = api.holdNext();
    const reply = api.handle({ method: 'GET', url: `${SNAPSHOT_SOURCE}?workspace_id=A` }); await hold.entered;
    let released = false; void reply.wait?.then(() => { released = true; });
    api.releaseAll(); await Promise.resolve(); expect(released).toBe(true);
  });
});
