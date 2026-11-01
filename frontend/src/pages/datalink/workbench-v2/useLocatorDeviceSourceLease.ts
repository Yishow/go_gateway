import { useEffect, useRef, useState } from 'react';
import type { StudioLocatorEntryRead } from '../../../types/studioV2LocatorEntry';
import type { StudioV2WorkspaceDeviceRecord } from '../../../services/studioV2WorkspaceDevices';
import { studioV2WorkspaceAPI } from '../../../services/studioV2Workspace';
import { studioV2WorkspaceDevicesAPI } from '../../../services/studioV2WorkspaceDevices';
import { studioV2WorkspaceWriteGroupsAPI } from '../../../services/studioV2WorkspaceWriteGroups';
import type { LocatorGroupSaveAcknowledgement } from './deviceSourceSubmission';

export interface LocatorDeviceSourceLease {
  allowed: boolean;
  acknowledgeLocatorGroupSave(proof: LocatorGroupSaveAcknowledgement): boolean;
  acquire(workspaceId: string): Promise<boolean>;
  readback(): Promise<readonly StudioV2WorkspaceDeviceRecord[] | undefined>;
}
interface Owner {
  origin?: StudioLocatorEntryRead;
  originRevision?: string;
  entry?: StudioLocatorEntryRead;
  workspaceId?: string;
  revision?: string;
  enabled: boolean;
}
function valid(owner: Owner): owner is Owner & { workspaceId: string; revision: string } {
  const selected = owner.entry?.selection;
  return Boolean(owner.enabled && owner.workspaceId?.trim() && owner.revision?.trim() &&
    owner.entry?.workspace.id === owner.workspaceId && owner.entry.locator.workspace_id === owner.workspaceId &&
    selected?.state === 'selected' && selected.target.workspaceId === owner.workspaceId && selected.target.workspaceRevision === owner.revision);
}
/** Local stable-read grant only. Neither an atomic source snapshot nor server CAS authority.
 * Uses database setup revision, never the unrelated Modbus Share revision.
 * No hydration dispatch occurs: an already-open recording draft is untouched. */
export function useLocatorDeviceSourceLease(entry?: StudioLocatorEntryRead, sourceWorkspaceId?: string, enabled = true): LocatorDeviceSourceLease {
  const owner = useRef<Owner>({ origin: entry, originRevision: entry?.entry.workspace_revision, entry, workspaceId: sourceWorkspaceId, revision: entry?.entry.workspace_revision, enabled });
  if (owner.current.origin !== entry || owner.current.workspaceId !== sourceWorkspaceId || owner.current.enabled !== enabled || owner.current.originRevision !== entry?.entry.workspace_revision) {
    owner.current = { origin: entry, originRevision: entry?.entry.workspace_revision, entry, workspaceId: sourceWorkspaceId, revision: entry?.entry.workspace_revision, enabled };
  }
  const [grant, setGrant] = useState<Owner>();
  const granted = useRef<Owner | undefined>(undefined);
  // Retain only read-only recovery eligibility after a failed read. It never grants writes.
  const previouslyGranted = useRef<Owner | undefined>(undefined);
  const sequence = useRef(0), mounted = useRef(true);
  useEffect(() => { mounted.current = true; return () => { mounted.current = false; sequence.current += 1; }; }, []);
  const capture = owner.current;
  const allowed = valid(capture) && grant === capture && granted.current === capture;
  const isCurrent = (captured: Owner, request: number) => mounted.current && owner.current === captured && sequence.current === request && valid(captured);
  const revoke = (captured: Owner, request: number) => {
    if (!isCurrent(captured, request)) return;
    granted.current = undefined; setGrant(undefined);
  };
  const read = async (captured: Owner, request: number, acquiring: boolean) => {
    if (!valid(captured)) return;
    const beforeWorkspace = await studioV2WorkspaceAPI.get();
    if (!isCurrent(captured, request) || beforeWorkspace.id !== captured.workspaceId) return;
    const beforeGroups = await studioV2WorkspaceWriteGroupsAPI.list();
    if (!isCurrent(captured, request) || beforeGroups.workspace_id !== captured.workspaceId || !beforeGroups.workspace_revision.trim() ||
      acquiring && beforeGroups.workspace_revision !== captured.revision) return;
    const devices = await studioV2WorkspaceDevicesAPI.list();
    if (!isCurrent(captured, request)) return;
    const afterWorkspace = await studioV2WorkspaceAPI.get();
    if (!isCurrent(captured, request) || afterWorkspace.id !== captured.workspaceId) return;
    const afterGroups = await studioV2WorkspaceWriteGroupsAPI.list();
    if (!isCurrent(captured, request) || afterGroups.workspace_id !== captured.workspaceId || afterGroups.workspace_revision !== beforeGroups.workspace_revision ||
      JSON.stringify(beforeWorkspace.ordered_device_ids) !== JSON.stringify(afterWorkspace.ordered_device_ids)) return;
    const ids = devices.map(device => device.id);
    if (new Set(ids).size !== ids.length || ids.some(id => !id.trim()) ||
      new Set(afterWorkspace.ordered_device_ids).size !== ids.length || afterWorkspace.ordered_device_ids.length !== ids.length ||
      ids.some(id => !afterWorkspace.ordered_device_ids.includes(id))) return;
    return devices;
  };
  return {
    allowed,
    acknowledgeLocatorGroupSave: ({ previous, acknowledgement, current }) => {
      // The render's owner cannot borrow a later A -> B -> A lease. Only the
      // exact own mutation and its canonical editor readback advance this basis.
      const before = capture.entry?.selection, prior = previous.selection, next = current.selection;
      const group = acknowledgement.group, list = current.entry.write_groups;
      if (owner.current !== capture || !valid(capture) || before?.state !== 'selected' || prior.state !== 'selected' || next.state !== 'selected' ||
        before.target.step !== 4 || prior.target.step !== 4 || next.target.step !== 4 || !before.target.groupId ||
        prior.target.workspaceId !== capture.workspaceId || prior.target.workspaceRevision !== capture.revision || previous.entry.workspace_revision !== capture.revision ||
        prior.target.groupId !== before.target.groupId || prior.target.groupRevision !== before.target.groupRevision ||
        group.workspace_id !== capture.workspaceId || group.id !== before.target.groupId || group.status === 'deleted' || !group.revision.trim() ||
        !acknowledgement.workspace_revision.trim() || current.workspace.id !== capture.workspaceId || current.locator.workspace_id !== capture.workspaceId ||
        current.locator.group_id !== group.id || current.entry.editor_hydrated !== true || next.target.workspaceId !== capture.workspaceId ||
        next.target.workspaceRevision !== acknowledgement.workspace_revision || current.entry.workspace_revision !== acknowledgement.workspace_revision ||
        next.target.groupId !== group.id || next.target.groupRevision !== group.revision || current.entry.group_heads[group.id] !== group.revision ||
        list?.workspace_id !== capture.workspaceId || list.workspace_revision !== acknowledgement.workspace_revision || list.groups.length !== 1 ||
        list.groups[0].id !== group.id || list.groups[0].revision !== group.revision || list.groups[0].status === 'deleted') return false;
      capture.entry = current; capture.revision = acknowledgement.workspace_revision;
      // Existing pending source receipts retain read-only recovery eligibility.
      // Their next readback still has to pass both canonical brackets.
      sequence.current += 1; granted.current = undefined; setGrant(undefined);
      return true;
    },
    acquire: async workspaceId => {
      const captured = owner.current, request = ++sequence.current;
      if (!valid(captured) || workspaceId !== captured.workspaceId) { revoke(captured, request); return false; }
      try {
        const devices = await read(captured, request, true);
        if (!devices || !isCurrent(captured, request)) { revoke(captured, request); return false; }
        previouslyGranted.current = captured; granted.current = captured; setGrant(captured); return true;
      } catch { revoke(captured, request); return false; }
    },
    readback: async () => {
      const captured = owner.current, request = ++sequence.current;
      if (!valid(captured) || previouslyGranted.current !== captured) return;
      try {
        // An acknowledged source mutation may advance DB setup revision. Match
        // this read's two revisions, not the obsolete pre-submit locator basis.
        const devices = await read(captured, request, false);
        if (!devices || !isCurrent(captured, request)) { revoke(captured, request); return; }
        granted.current = captured; setGrant(captured);
        return devices;
      } catch { revoke(captured, request); return; }
    },
  };
}
