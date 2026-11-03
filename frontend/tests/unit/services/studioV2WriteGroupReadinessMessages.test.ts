import { AxiosHeaders } from 'axios';
import { beforeEach, expect, it, vi } from 'vitest';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { WriteGroupReadiness } from '@/types/studioV2WriteGroup';

vi.mock('@/services/studioV2Workspace', () => ({ studioV2DatalinkApi: { get: vi.fn() } }));
const reductionMessage = "the recording reduction, member recording metadata or layout is not supported or cannot be verified; the saved group cannot start a writer";
function response(): WriteGroupReadiness {
  return { workspace_id: 'workspace-1', workspace_revision: 'workspace-rev-1', group_id: 'group-1', group_revision: 'group-rev-1', applied_revision: '',
    config_ready: false, schema_ready: false, ready: false, issues: [
      { code: 'source-revision-stale', severity: 'blocking', step: 'Step 4', scope: 'group-1', message: 'saved source or mapping revision is stale' },
      { code: 'schema-check-blocked', severity: 'blocking', step: 'Step 4', scope: 'group-1', message: 'destination schema check is blocked by configuration issues' },
      { code: 'recording-reduction-unavailable', severity: 'blocking', step: 'Step 4', scope: 'group-1', message: reductionMessage },
    ] };
}
function reply(data: WriteGroupReadiness) { vi.mocked(studioV2DatalinkApi.get).mockResolvedValueOnce({ data: { success: true, data }, status: 200, statusText: 'OK', headers: {}, config: { headers: new AxiosHeaders() } }); }
beforeEach(() => vi.clearAllMocks());

it('preserves the real 138-character readiness description and all three blocking issues', async () => {
  const data = response(); expect(reductionMessage).toHaveLength(138); reply(data);
  await expect(studioV2WorkspaceWriteGroupsAPI.readiness('group-1')).resolves.toEqual(data);
});
it('accepts a readiness description at the existing 256-character JSON boundary', async () => {
  const data = response(); data.issues[2].message = 'm'.repeat(256); reply(data);
  await expect(studioV2WorkspaceWriteGroupsAPI.readiness('group-1')).resolves.toEqual(data);
});
it('rejects a 257-character description with the original fixed safe error', async () => {
  const data = response(); data.issues[2].message = 'private-'.repeat(37).slice(0, 257); reply(data);
  const error: unknown = await studioV2WorkspaceWriteGroupsAPI.readiness('group-1').catch((value: unknown) => value);
  expect(error).toMatchObject({ message: 'write-group readiness response was invalid', outcome: 'unconfirmed' });
  expect(String(error)).not.toContain('private-');
});
it.each(['workspace_id', 'code'] as const)('keeps the 128-character identifier limit for %s', async field => {
  const data = response(); data.issues[2].message = 'short';
  if (field === 'workspace_id') data.workspace_id = 'i'.repeat(129); else data.issues[2].code = 'c'.repeat(129);
  reply(data); await expect(studioV2WorkspaceWriteGroupsAPI.readiness('group-1')).rejects.toMatchObject({ message: 'write-group readiness response was invalid' });
});
it('still accepts identifiers at 128 characters without widening their contract', async () => {
  const data = response(); data.workspace_id = 'i'.repeat(128); data.issues[2].code = 'c'.repeat(128); data.issues[2].message = 'short'; reply(data);
  await expect(studioV2WorkspaceWriteGroupsAPI.readiness('group-1')).resolves.toEqual(data);
});
