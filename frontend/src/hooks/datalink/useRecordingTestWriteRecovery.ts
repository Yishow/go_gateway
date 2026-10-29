import type { StudioV2LocatorReadScope } from '../../types/studioV2LocatorEntry';
import { studioLocatorReadKey } from '../../services/studioEntityLocator';
import { useEffect, useRef, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { listRecordingTestWriteRecovery } from '../../services/recordingTestWriteRecovery';
import { recordingTestWriteRecovery } from '../../features/datalink/workbench-v2/state/recordingTestWriteRecovery';
/** Listing discovers locators, but its absence never resolves a previously known write. */
export function useRecordingTestWriteRecovery(workspaceId: string, locatorScope?: StudioV2LocatorReadScope) {
  const selectedLocatorScope = locatorScope === undefined ? undefined : { workspaceId: locatorScope.workspaceId, workspaceRevision: locatorScope.workspaceRevision };
  const query = useQuery({
    queryKey: ['studio-v2-workspace', 'test-write-recovery-list', workspaceId, ...studioLocatorReadKey(selectedLocatorScope)],
    queryFn: () => selectedLocatorScope ? listRecordingTestWriteRecovery(workspaceId, selectedLocatorScope) : listRecordingTestWriteRecovery(workspaceId), enabled: Boolean(workspaceId),
    retry: false, refetchOnMount: 'always', refetchOnWindowFocus: false,
  });
  const bindings = useRef(new Map<string, string>());
  const [accepted, setAccepted] = useState<unknown>(null);
  useEffect(() => {
    if (!query.data || query.isError || query.isFetching) return;
    const identities = query.data.operations.map((operation) => [operation.operation_id, JSON.stringify([
      operation.workspace_id, operation.workspace_revision, operation.group_id, operation.group_revision,
      operation.connector_id, operation.connector_revision, operation.target_fingerprint, operation.created_at,
    ])] as const);
    if (identities.some(([id, identity]) => bindings.current.has(id) && bindings.current.get(id) !== identity)) {
      setAccepted(null);
      return;
    }
    const existing = recordingTestWriteRecovery.read().references;
    if (query.data.workspace_id !== workspaceId || query.data.operations.some((operation) => existing.some((ref) =>
      ref.operation_id === operation.operation_id && (ref.workspace_id !== operation.workspace_id ||
        ref.group_id !== operation.group_id || ref.group_revision !== operation.group_revision)))) {
      setAccepted(null);
      return;
    }
    identities.forEach(([id, identity]) => bindings.current.set(id, identity));
    for (const operation of query.data.operations) recordingTestWriteRecovery.retain({
      workspace_id: workspaceId, group_id: operation.group_id, group_revision: operation.group_revision,
      operation_id: operation.operation_id, operation_kind: 'test_write', created_at: operation.created_at,
    });
    setAccepted(query.data);
  }, [query.data, query.isError, query.isFetching, workspaceId]);
  return { ...query, available: Boolean(query.data && accepted === query.data && !query.isError && !query.isFetching) };
}
