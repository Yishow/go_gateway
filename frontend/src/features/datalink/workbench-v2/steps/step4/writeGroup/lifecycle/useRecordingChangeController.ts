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
  /** Denial-only caller lease, checked before every transport; absent preserves ordinary callers. */
  canUseScope?: () => boolean;
  onChanged?: (group: WriteGroup, workspaceRevision: string) => void;
}
export function useRecordingChangeController({ workspaceId, readonly, canUseScope, onChanged }: UseRecordingChangeControllerOptions) {
  const actions = useWriteGroupSuccessorActions();
  const check = useCheckSuccessorOperation();
  const scoped = <T,>(send: () => Promise<T>): Promise<T> => canUseScope?.() === false
    ? Promise.reject(new Error('Recording scope is unavailable.')) : send();
  const lifecycle = canUseScope ? (id: string) => scoped(() => studioV2WriteGroupLifecycleAPI.lifecycle(id)) : studioV2WriteGroupLifecycleAPI.lifecycle;
  const getGroup = canUseScope ? (id: string) => scoped(() => studioV2WorkspaceWriteGroupsAPI.get(id)) : studioV2WorkspaceWriteGroupsAPI.get;
  const api: typeof studioV2WriteGroupSuccessorAPI = {
    review: (groupId, request) => scoped(() => actions.review.mutateAsync({ groupId, request })),
    reserve: (groupId, request) => scoped(() => actions.reserve.mutateAsync({ groupId, request })),
    schemaPreview: (groupId, request) => scoped(() => actions.schemaPreview.mutateAsync({ groupId, request })),
    schemaConfirm: (groupId, request) => scoped(() => actions.schemaConfirm.mutateAsync({ groupId, request })),
    activationReview: (groupId, request) => scoped(() => actions.activationReview.mutateAsync({ groupId, request })),
    activationReadiness: (groupId, request) => scoped(() => actions.activationReadiness.mutateAsync({ groupId, request })),
    confirmCutover: (groupId, request) => scoped(() => actions.confirmCutover.mutateAsync({ groupId, request })),
    confirmFirstStart: (groupId, request) => scoped(() => actions.confirmFirstStart.mutateAsync({ groupId, request })),
    operationStatus: (groupId, request) => scoped(() => check(groupId, request)),
  };
  // Workspace identity, not the selected mode, owns the controller and its recovery fence.
  const [controller, setController] = React.useState(() => new RecordingChangeController(workspaceId, {
    api, lifecycle, getGroup, onChanged,
  }));
  if (controller.workspaceId !== workspaceId) {
    controller.dispose();
    setController(new RecordingChangeController(workspaceId, {
      api, lifecycle, getGroup, onChanged,
    }));
  }
  React.useLayoutEffect(() => {
    controller.deps = { api, lifecycle, getGroup, onChanged };
    controller.readonly = readonly;
  });
  React.useEffect(() => { controller.activate(); return () => controller.dispose(); }, [controller]);
  const state = React.useSyncExternalStore(controller.subscribe, controller.snapshot, controller.snapshot);
  return { ...state, state, controller, requestAction: controller.requestAction, dismiss: controller.dismiss, unresolved: controller.unresolved };
}
export type RecordingChangeControllerView = ReturnType<typeof useRecordingChangeController>;
