import { useEffect, useRef, useState } from 'react';
import { useWriteGroupLifecycleMutation } from '@/hooks/datalink/useStudioV2WriteGroups';
import { studioV2WriteGroupLifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
import { DeleteResponseError } from '@/services/studioV2WorkspaceDelete';
import type { WriteGroupDeleteReceipt } from '@/types/studioV2Delete';
import { lifecycleMatches } from '../steps/step4/writeGroup/lifecycle/SavedRecordingActions';
import { SavedRecordingRemoveDialog } from '../steps/step4/writeGroup/lifecycle/SavedRecordingRemoveDialog';

interface DeleteGroup {
  readonly id: string;
  readonly name: string;
  readonly workspace_id: string;
  readonly revision: string;
  readonly applied_revision: string;
  readonly destination: { readonly connector_id: string; readonly connector_revision: string; readonly table_schema: string; readonly table_name: string };
}
export interface RecordingDeleteDialogProps {
  open: boolean;
  group: DeleteGroup;
  workspaceRevision: string;
  onClose: () => void;
  onDeleted: (receipt: WriteGroupDeleteReceipt) => void;
  onStop: () => void;
  onOpenRecording: () => void;
  onRefresh: () => void;
  refreshFailed?: boolean;
}

/** Rechecks server action facts before deleting the selected saved revision. */
export function RecordingDeleteDialog({ open, group, workspaceRevision, onClose, onDeleted, onStop, onOpenRecording, onRefresh, refreshFailed }: RecordingDeleteDialogProps) {
  const mutation = useWriteGroupLifecycleMutation('remove');
  const [failure, setFailure] = useState<{ owner: object; error: unknown }>();
  const [checking, setChecking] = useState(false);
  const inFlight = useRef(false);
  const mounted = useRef(false);
  const key = JSON.stringify([open, group.workspace_id, workspaceRevision, group.id, group.revision, group.applied_revision, group.destination.connector_id, group.destination.connector_revision]);
  const owner = useRef({ key });
  if (owner.current.key !== key) owner.current = { key };
  const error = failure?.owner === owner.current ? failure.error : null;
  useEffect(() => { mounted.current = true; return () => { mounted.current = false; }; }, []);
  const submit = async () => {
    if (inFlight.current || error) return;
    const selectedOwner = owner.current;
    inFlight.current = true; setChecking(true);
    try {
      const facts = await studioV2WriteGroupLifecycleAPI.lifecycle(group.id);
      if (!mounted.current || selectedOwner !== owner.current) return;
      if (!lifecycleMatches(facts, { workspaceId: group.workspace_id, workspaceRevision, groupId: group.id, groupRevision: group.revision }) ||
        facts.applied_revision !== group.applied_revision || facts.connector_id !== group.destination.connector_id ||
        facts.connector_revision !== group.destination.connector_revision || !facts.actions.some(action => action.action === 'remove' && action.allowed))
        throw new DeleteResponseError('failed', 'revision_mismatch', 'reload');
      const result = await mutation.mutateAsync({ id: group.id, request: {
        workspace_id: group.workspace_id, expected_workspace_revision: workspaceRevision,
        expected_group_revision: group.revision, expected_connector_revision: group.destination.connector_revision,
      } });
      if (!('deleted' in result)) throw new Error('Deletion receipt unavailable');
      if (mounted.current && selectedOwner === owner.current) onDeleted(result);
    } catch (cause) { if (mounted.current && selectedOwner === owner.current) setFailure({ owner: selectedOwner, error: cause }); }
    finally { inFlight.current = false; if (mounted.current) setChecking(false); }
  };
  return <SavedRecordingRemoveDialog open={open} busy={checking || mutation.isPending} group={group} objectLabel={group.name} error={error} refreshFailed={refreshFailed}
    onClose={onClose} onConfirm={() => { void submit(); }} onStop={onStop} onOpenRecording={onOpenRecording}
    onRefresh={onRefresh} />;
}
