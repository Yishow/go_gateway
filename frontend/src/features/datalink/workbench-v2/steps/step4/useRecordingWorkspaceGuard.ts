import { useCallback, useLayoutEffect, useRef } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from '../../../../../hooks/datalink/keys';
import type { WriteGroupListResponse } from '../../../../../types/studioV2WriteGroup';

/** Synchronous denial from the same canonical query; never fetches, writes, or grants authority. */
export function useRecordingWorkspaceGuard(workspaceId: string | undefined, accepted: WriteGroupListResponse | undefined) {
  const client = useQueryClient();
  const owner = useRef({ workspaceId, accepted, epoch: 0 });
  if (owner.current.workspaceId !== workspaceId || owner.current.accepted !== accepted) {
    owner.current = { workspaceId, accepted, epoch: owner.current.epoch + 1 };
  }
  const epoch = owner.current.epoch;
  const mounted = useRef(true);
  useLayoutEffect(() => { mounted.current = true; return () => { mounted.current = false; }; }, []);
  return useCallback(() => {
    const current = client.getQueryState<WriteGroupListResponse>(studioV2WorkspaceKeys.writeGroups());
    return !mounted.current || owner.current.epoch !== epoch || !workspaceId || !accepted ||
      !current || current.data !== accepted || current.status !== 'success' || current.fetchStatus !== 'idle' ||
      accepted.workspace_id !== workspaceId || !accepted.workspace_revision;
  }, [client, workspaceId, accepted, epoch]);
}
