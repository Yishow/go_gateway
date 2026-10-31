import { studioV2WorkspaceAPI } from '../../../services/studioV2Workspace';
import type { StudioV2WorkspaceDeviceRecord } from '../../../services/studioV2WorkspaceDevices';
/** Exact-device read in the existing owner's scope, not mutation CAS. */
export async function readSourceRangeDevice(workspaceId: string, deviceId: string, readback: () => Promise<readonly StudioV2WorkspaceDeviceRecord[] | undefined>) {
  const before = await studioV2WorkspaceAPI.get();
  const revision = before.modbus_share?.workspace_revision;
  if (!revision?.trim() || before.id !== workspaceId || !before.ordered_device_ids.includes(deviceId)) return;
  const rows = await readback();
  const after = await studioV2WorkspaceAPI.get();
  if (after.modbus_share?.workspace_revision !== revision || after.id !== workspaceId || JSON.stringify(after.ordered_device_ids) !== JSON.stringify(before.ordered_device_ids)) return;
  const exact = rows?.filter(row => row.id === deviceId);
  return exact?.length === 1 && exact[0].runtime_apply_status !== 'apply_failed' ? exact[0] : undefined;
}
