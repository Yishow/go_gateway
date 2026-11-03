import { GroupReadinessPanel } from '../steps/step4/writeGroup/GroupReadinessPanel';
import { useRef, type ReactNode } from 'react';
import { recordingTestWriteRecovery } from '../state/recordingTestWriteRecovery';
import { GroupSchemaPanel, type GroupSchemaPanelProps } from '../steps/step4/writeGroup/GroupSchemaPanel';
import { useRecordingPreparationOwner } from '../steps/step4/useRecordingPreparationOwner';
import type { WriteGroupReadiness } from '../../../../types/studioV2WriteGroup';
import { RecordingPreparationActions, type RecordingSavePresentation } from './RecordingPreparationActions';

export interface GroupEditorPreparationProps extends Omit<GroupSchemaPanelProps, 'guided' | 'autoPreview' | 'previewRequest'> {
  compact?: boolean; sourceRefresh?: ReactNode; managed: boolean; consentIdentity: string; save: RecordingSavePresentation; blocked: boolean;
  readiness: { data?: WriteGroupReadiness; loading: boolean; failed: boolean; updatedAt: number; check: () => Promise<unknown> };
}
/** Reuses the original save/schema/query owners; stores only a denial-only native consent generation. */
export function GroupEditorPreparation({ compact = false, managed, consentIdentity, save, blocked, readiness, sourceRefresh, ...schemaProps }: GroupEditorPreparationProps) {
  const { group, workspaceId, workspaceRevision, readonly } = schemaProps;
  const owner = useRecordingPreparationOwner(consentIdentity, workspaceId, group ? [group.id] : []);
  const lease = useRef({ epoch: owner.epoch, observedAt: 0 });
  if (lease.current.epoch !== owner.epoch) lease.current = { epoch: owner.epoch, observedAt: readiness.updatedAt };
  const fresh = !readiness.loading && !readiness.failed && readiness.updatedAt > lease.current.observedAt &&
    readiness.data?.workspace_id === workspaceId && readiness.data.workspace_revision === workspaceRevision && readiness.data.group_id === group?.id && readiness.data.group_revision === group?.revision;
  // Existing retained test-write cleanup may deny new schema work; no test write is requested or required here.
  const retainedTestWrite = owner.testWrite.state === 'unresolved';
  const newSchemaBlocked = () => blocked || retainedTestWrite || recordingTestWriteRecovery.read().references.some(
    reference => reference.workspace_id === workspaceId && reference.group_id === group?.id);
  const statusContent = compact && group ? <GroupReadinessPanel compact readiness={readiness.data} loading={readiness.loading && !save.dirty} failed={readiness.failed} dirty={save.dirty} onRetry={() => void readiness.check()} /> : undefined;
  const props = { statusContent, sourceRefresh, save, readonly, blocked: blocked || retainedTestWrite, pendingTestWrite: retainedTestWrite, online: owner.online,
    readinessReady: Boolean(fresh && readiness.data?.ready), readinessBusy: readiness.loading, checkReadiness: readiness.check };
  return managed && group ? <GroupSchemaPanel {...schemaProps} guided={{ online: owner.online, consentEpoch: owner.epoch,
    newOperationBlocked: newSchemaBlocked(), isNewOperationBlocked: newSchemaBlocked, isPreviewRetryBlocked: newSchemaBlocked,
    renderActions: schema => <RecordingPreparationActions {...props} schema={schema} /> }} /> : <RecordingPreparationActions {...props} />;
}
