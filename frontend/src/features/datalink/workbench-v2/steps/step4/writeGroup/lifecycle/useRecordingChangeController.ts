import * as React from 'react';
import { useWriteGroupSuccessorActions, useCheckSuccessorOperation } from '../../../../../../../hooks/datalink/useStudioV2WriteGroupSuccessor';
import { studioV2WriteGroupLifecycleAPI } from '../../../../../../../services/studioV2WriteGroupLifecycle';
import { studioV2WorkspaceWriteGroupsAPI } from '../../../../../../../services/studioV2WorkspaceWriteGroups';
import type { studioV2WriteGroupSuccessorAPI } from '../../../../../../../services/studioV2WriteGroupSuccessor';
import type { WriteGroup } from '../../../../../../../types/studioV2WriteGroup';
import { RecordingChangeController } from './RecordingChangeController';
export interface UseRecordingChangeControllerOptions {
  workspaceId: string;
  readonly: boolean;
  onChanged?: (group: WriteGroup, workspaceRevision: string) => void;
}
export function useRecordingChangeController({ workspaceId, readonly, onChanged }: UseRecordingChangeControllerOptions) {
  const actions = useWriteGroupSuccessorActions();
  const check = useCheckSuccessorOperation();
  const api: typeof studioV2WriteGroupSuccessorAPI = {
    review: (groupId, request) => actions.review.mutateAsync({ groupId, request }),
    reserve: (groupId, request) => actions.reserve.mutateAsync({ groupId, request }),
    schemaPreview: (groupId, request) => actions.schemaPreview.mutateAsync({ groupId, request }),
    schemaConfirm: (groupId, request) => actions.schemaConfirm.mutateAsync({ groupId, request }),
    activationReview: (groupId, request) => actions.activationReview.mutateAsync({ groupId, request }),
    activationReadiness: (groupId, request) => actions.activationReadiness.mutateAsync({ groupId, request }),
    confirmCutover: (groupId, request) => actions.confirmCutover.mutateAsync({ groupId, request }),
    confirmFirstStart: (groupId, request) => actions.confirmFirstStart.mutateAsync({ groupId, request }),
    operationStatus: check,
  };
  // Workspace identity, not the selected mode, owns the controller and its recovery fence.
  const [controller, setController] = React.useState(() => new RecordingChangeController(workspaceId, {
    api, lifecycle: studioV2WriteGroupLifecycleAPI.lifecycle, getGroup: studioV2WorkspaceWriteGroupsAPI.get, onChanged,
  }));
  if (controller.workspaceId !== workspaceId) {
    controller.dispose();
    setController(new RecordingChangeController(workspaceId, {
      api, lifecycle: studioV2WriteGroupLifecycleAPI.lifecycle, getGroup: studioV2WorkspaceWriteGroupsAPI.get, onChanged,
    }));
  }
  React.useLayoutEffect(() => {
    controller.deps = { api, lifecycle: studioV2WriteGroupLifecycleAPI.lifecycle, getGroup: studioV2WorkspaceWriteGroupsAPI.get, onChanged };
    controller.readonly = readonly;
  });
  React.useEffect(() => { controller.activate(); return () => controller.dispose(); }, [controller]);
  const state = React.useSyncExternalStore(controller.subscribe, controller.snapshot, controller.snapshot);
  return { ...state, state, controller, requestAction: controller.requestAction, dismiss: controller.dismiss, unresolved: controller.unresolved };
}
export type RecordingChangeControllerView = ReturnType<typeof useRecordingChangeController>;
