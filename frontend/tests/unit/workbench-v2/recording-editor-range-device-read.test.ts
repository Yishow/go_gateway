import { beforeEach, expect, it, vi } from 'vitest';
import { readSourceRangeDevice } from '@/pages/datalink/workbench-v2/readSourceRangeDevice';
import { studioV2WorkspaceAPI } from '@/services/studioV2Workspace';
import type { StudioV2Workspace } from '@/types/studioV2Workspace';
import type { StudioV2WorkspaceDeviceRecord } from '@/services/studioV2WorkspaceDevices';
vi.mock('@/services/studioV2Workspace', () => ({ studioV2WorkspaceAPI: { get: vi.fn() } }));
const device = { id: 'D', protocol: 'modbus_tcp' } as StudioV2WorkspaceDeviceRecord;
const workspace = (revision?: string) => ({ id: 'W', ordered_device_ids: ['D'], modbus_share: revision ? { workspace_revision: revision } : undefined }) as StudioV2Workspace;
beforeEach(() => vi.mocked(studioV2WorkspaceAPI.get).mockReset());
it.each([undefined, ''])('does not authorize a range using a missing source-domain revision (%s)', async revision => {
 vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(workspace(revision));
 const read = vi.fn(async () => [device]);
 expect(await readSourceRangeDevice('W', 'D', read)).toBeUndefined();
 expect(read).not.toHaveBeenCalled();
});
it('rejects a source revision change during selected-device readback', async () => {
 vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValueOnce(workspace('R1')).mockResolvedValue(workspace('R2'));
 expect(await readSourceRangeDevice('W', 'D', async () => [device])).toBeUndefined();
});
it('authorizes only the exact unique device with stable workspace source scope', async () => {
 vi.mocked(studioV2WorkspaceAPI.get).mockResolvedValue(workspace('R1'));
 expect(await readSourceRangeDevice('W', 'D', async () => [device])).toEqual(device);
 expect(await readSourceRangeDevice('W', 'D', async () => [device, device])).toBeUndefined();
});
