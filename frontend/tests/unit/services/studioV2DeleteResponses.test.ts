import { beforeEach, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2WorkspaceDevicesAPI } from '@/services/studioV2WorkspaceDevices';

vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { delete: vi.fn() } }));

beforeEach(() => vi.resetAllMocks());

it('accepts the deletion receipt for the requested two-generation recording', async () => {
  vi.mocked(studioV2DatalinkApi.delete).mockResolvedValueOnce({ data: { success: true, data: {
    deleted: true, group_id: 'G2', lineage_id: 'L1', deleted_generations: 2,
  } } } as never);
  const request = { workspace_id: 'W1', expected_workspace_revision: 'W2', expected_group_revision: 'R2', expected_connector_revision: 'C1' };
  await expect(studioV2WorkspaceWriteGroupsAPI.remove('G2', request)).resolves.toEqual({
    deleted: true, group_id: 'G2', lineage_id: 'L1', deleted_generations: 2,
  });
  expect(studioV2DatalinkApi.delete).toHaveBeenCalledWith('/studio-v2/workspace/write-groups/G2', { data: request });
});

it('accepts an empty lineage only for one legacy generation', async () => {
  const request = { workspace_id: 'W1', expected_workspace_revision: 'W2', expected_group_revision: 'R2', expected_connector_revision: 'C1' };
  vi.mocked(studioV2DatalinkApi.delete).mockResolvedValueOnce({ data: { success: true, data: {
    deleted: true, group_id: 'G2', lineage_id: '', deleted_generations: 1,
  } } } as never);
  await expect(studioV2WorkspaceWriteGroupsAPI.remove('G2', request)).resolves.toEqual({ deleted: true, group_id: 'G2', lineage_id: '', deleted_generations: 1 });
  vi.mocked(studioV2DatalinkApi.delete).mockResolvedValueOnce({ data: { success: true, data: {
    deleted: true, group_id: 'G2', lineage_id: '', deleted_generations: 2,
  } } } as never);
  await expect(studioV2WorkspaceWriteGroupsAPI.remove('G2', request)).rejects.toMatchObject({ outcome: 'unconfirmed' });
});

it('preserves separate unconsumed samples and outbox states in the backlog conflict', async () => {
  vi.mocked(studioV2DatalinkApi.delete).mockRejectedValueOnce({ response: { status: 409, data: {
    success: false, error: { code: 'WRITE_GROUP_DELETE_BACKLOG', action: 'stop_and_wait', message: 'private',
      details: { samples_unconsumed: 3, outbox: { retrying: 1, unknown: 1, future_state: 2 } } },
  } } });
  const request = { workspace_id: 'W1', expected_workspace_revision: 'W2', expected_group_revision: 'R2', expected_connector_revision: 'C1' };
  await expect(studioV2WorkspaceWriteGroupsAPI.remove('G2', request)).rejects.toMatchObject({
    code: 'WRITE_GROUP_DELETE_BACKLOG', action: 'stop_and_wait', outcome: 'failed',
    details: { samples_unconsumed: 3, outbox: { retrying: 1, unknown: 1, future_state: 2 } },
  });
});

it('keeps an unfamiliar outbox state as a separate count', async () => {
  vi.mocked(studioV2DatalinkApi.delete).mockRejectedValueOnce({ response: { status: 409, data: {
    success: false, error: { code: 'WRITE_GROUP_DELETE_BACKLOG', action: 'stop_and_wait',
      details: { samples_unconsumed: 0, outbox: { 'NEW-STATE': 2 } } },
  } } });
  const request = { workspace_id: 'W1', expected_workspace_revision: 'W2', expected_group_revision: 'R2', expected_connector_revision: 'C1' };
  await expect(studioV2WorkspaceWriteGroupsAPI.remove('G2', request)).rejects.toMatchObject({
    outcome: 'failed', details: { samples_unconsumed: 0, outbox: { 'NEW-STATE': 2 } },
  });
});

it('classifies a stale revision as a refusal requiring refresh', async () => {
  vi.mocked(studioV2DatalinkApi.delete).mockRejectedValueOnce({ response: { status: 409, data: {
    success: false, error: { code: 'revision_mismatch', action: 'reload' },
  } } });
  const request = { workspace_id: 'W1', expected_workspace_revision: 'old', expected_group_revision: 'R2', expected_connector_revision: 'C1' };
  await expect(studioV2WorkspaceWriteGroupsAPI.remove('G2', request)).rejects.toMatchObject({
    outcome: 'failed', code: 'revision_mismatch', action: 'reload',
  });
});

it('treats an unfamiliar coded conflict as a refusal and offers refresh', async () => {
  vi.mocked(studioV2DatalinkApi.delete).mockRejectedValueOnce({ response: { status: 409, data: {
    success: false, error: { code: 'WRITE_GROUP_LIFECYCLE_BLOCKED', message: 'private details' },
  } } });
  const request = { workspace_id: 'W1', expected_workspace_revision: 'W2', expected_group_revision: 'R2', expected_connector_revision: 'C1' };
  const error = await studioV2WorkspaceWriteGroupsAPI.remove('G2', request).catch((value: unknown) => value);
  expect(error).toMatchObject({ outcome: 'failed', code: 'WRITE_GROUP_LIFECYCLE_BLOCKED', action: 'reload' });
  expect(JSON.stringify(error)).not.toContain('private details');
});

it.each([400, 404, 503])('treats a typed HTTP %s deletion error as a failure requiring refresh', async status => {
  vi.mocked(studioV2DatalinkApi.delete).mockRejectedValueOnce({ response: { status, data: {
    success: false, error: { code: 'WRITE_GROUP_UNAVAILABLE', message: 'private details' },
  } } });
  const request = { workspace_id: 'W1', expected_workspace_revision: 'W2', expected_group_revision: 'R2', expected_connector_revision: 'C1' };
  const error = await studioV2WorkspaceWriteGroupsAPI.remove('G2', request).catch((value: unknown) => value);
  expect(error).toMatchObject({ outcome: 'failed', code: 'WRITE_GROUP_UNAVAILABLE', action: 'reload' });
  expect(JSON.stringify(error)).not.toContain('private details');
});

it.each([409, 503])('reports an HTTP %s error without a usable envelope as a failed request', async status => {
  vi.mocked(studioV2DatalinkApi.delete).mockRejectedValueOnce({ response: { status, data: {} } });
  const request = { workspace_id: 'W1', expected_workspace_revision: 'W2', expected_group_revision: 'R2', expected_connector_revision: 'C1' };
  await expect(studioV2WorkspaceWriteGroupsAPI.remove('G2', request)).rejects.toMatchObject({ outcome: 'failed', action: 'reload' });
});

it.each([
  ['WRITE_GROUP_DELETE_ACTIVE', 'stop_first'],
  ['WRITE_GROUP_DELETE_OPERATION_PENDING', 'open_recording'],
  ['WRITE_GROUP_DELETE_HAS_SUCCESSOR', 'open_recording'],
])('classifies %s as a refusal with its next action', async (code, action) => {
  vi.mocked(studioV2DatalinkApi.delete).mockRejectedValueOnce({ response: { status: 409, data: {
    success: false, error: { code, action, message: 'private' },
  } } });
  const request = { workspace_id: 'W1', expected_workspace_revision: 'W2', expected_group_revision: 'R2', expected_connector_revision: 'C1' };
  const error = await studioV2WorkspaceWriteGroupsAPI.remove('G2', request).catch((value: unknown) => value);
  expect(error).toMatchObject({ code, action, outcome: 'failed' });
  expect(JSON.stringify(error)).not.toContain('private');
});

it('exposes server-reported recordings for a referenced device', async () => {
  vi.mocked(studioV2DatalinkApi.delete).mockRejectedValueOnce({ response: { status: 409, data: {
    success: false, error: { code: 'DEVICE_REFERENCED_BY_RECORDINGS', action: 'remove_recordings_first',
      details: { recordings: [{ id: 'G1', name: 'First' }, { id: 'G2', name: 'Second' }], total: 2 } },
  } } });
  await expect(studioV2WorkspaceDevicesAPI.remove('device-1')).rejects.toMatchObject({
    code: 'DEVICE_REFERENCED_BY_RECORDINGS', outcome: 'failed',
    details: { recordings: [{ id: 'G1', name: 'First' }, { id: 'G2', name: 'Second' }], total: 2 },
  });
});

it('requires explicit device success and treats malformed replies or network loss as unconfirmed', async () => {
  vi.mocked(studioV2DatalinkApi.delete).mockResolvedValueOnce({ data: { success: true } } as never);
  await expect(studioV2WorkspaceDevicesAPI.remove('device-1')).resolves.toBeUndefined();
  vi.mocked(studioV2DatalinkApi.delete).mockResolvedValueOnce({ data: {} } as never);
  await expect(studioV2WorkspaceDevicesAPI.remove('device-1')).rejects.toMatchObject({ outcome: 'unconfirmed' });
  vi.mocked(studioV2DatalinkApi.delete).mockRejectedValueOnce(new Error('private network details'));
  const error = await studioV2WorkspaceDevicesAPI.remove('device-1').catch((value: unknown) => value);
  expect(error).toMatchObject({ outcome: 'unconfirmed' });
  expect(JSON.stringify(error)).not.toContain('private');
});

it('does not treat a lost or wrong-recording deletion receipt as success', async () => {
  const request = { workspace_id: 'W1', expected_workspace_revision: 'W2', expected_group_revision: 'R2', expected_connector_revision: 'C1' };
  vi.mocked(studioV2DatalinkApi.delete).mockResolvedValueOnce({ data: { success: true, data: {
    deleted: true, group_id: 'other', lineage_id: 'L1', deleted_generations: 2,
  } } } as never);
  await expect(studioV2WorkspaceWriteGroupsAPI.remove('G2', request)).rejects.toMatchObject({ outcome: 'unconfirmed' });
  vi.mocked(studioV2DatalinkApi.delete).mockRejectedValueOnce(new Error('private network details'));
  await expect(studioV2WorkspaceWriteGroupsAPI.remove('G2', request)).rejects.toMatchObject({ outcome: 'unconfirmed' });
});

it('preserves the dependent recordings supplied by a successor conflict', async () => {
  vi.mocked(studioV2DatalinkApi.delete).mockRejectedValueOnce({ response: { status: 409, data: {
    success: false, error: { code: 'WRITE_GROUP_DELETE_HAS_SUCCESSOR', action: 'open_recording',
      details: { referenced_by: [{ id: 'G3', name: 'Later recording' }], total: 2 } },
  } } });
  const request = { workspace_id: 'W1', expected_workspace_revision: 'W2', expected_group_revision: 'R2', expected_connector_revision: 'C1' };
  await expect(studioV2WorkspaceWriteGroupsAPI.remove('G2', request)).rejects.toMatchObject({
    code: 'WRITE_GROUP_DELETE_HAS_SUCCESSOR', details: { referenced_by: [{ id: 'G3', name: 'Later recording' }], total: 2 },
  });
});
