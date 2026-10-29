import type { SchemaPreviewToken } from '../../../../../../../types/recordingPlan';
import type { WriteGroup } from '../../../../../../../types/studioV2WriteGroup';
import type { WriteGroupCutover } from '../../../../../../../types/studioV2WriteGroupCutover';
import type { WriteGroupLifecycle, SavedRecordingActionIntent } from '../../../../../../../types/studioV2WriteGroupLifecycle';
import type * as S from '../../../../../../../types/studioV2WriteGroupSuccessor';
export type { SavedRecordingActionIntent };
export interface RecordingDestinationChoice { id: string; name: string; revision: string }
export type RecordingChangePhase = 'choose' | 'impact' | 'reserved' | 'schema-plan' | 'schema-confirmed' | 'activation-ready' | 'scheduled';
export type RecordingChangeRequest =
  | { action: 'reserve'; request: S.WriteGroupSuccessorReservationRequest }
  | { action: 'schemaPreview'; request: S.WriteGroupSuccessorStageRequest }
  | { action: 'schemaConfirm'; request: S.WriteGroupSuccessorSchemaConfirmation }
  | { action: 'confirmCutover' | 'confirmFirstStart'; request: S.WriteGroupSuccessorActivationConfirmation };
/** Original envelope and exact submitted request are immutable across current CAS refreshes. */
export interface RecordingChangeRecovery {
  version: 1;
  attempt?: number;
  envelope: S.WriteGroupSuccessorEnvelope;
  reviewDigest: string;
  submitted: RecordingChangeRequest;
  phase: RecordingChangePhase;
  outcome: 'pending' | 'recorded' | 'unknown' | 'rejected';
  preparation?: S.WriteGroupSuccessorPreparation;
  preview?: Omit<SchemaPreviewToken, 'statements' | 'database'>;
  previewStage?: S.WriteGroupSuccessorStageRequest;
  stageRejected?: boolean;
  unreadable?: boolean;
  priorPreviewIds?: string[];
  priorPreviewStates?: Record<string, S.WriteGroupSuccessorOperationStatus['schema_operations'][number]['status']>;
  schemaConfirmationSubmitted?: boolean;
  terminalStatus?: S.WriteGroupSuccessorOperationStatus;
  result?: WriteGroupCutover | S.WriteGroupSuccessorFirstStartResult;
}
export interface UnreadableRecordingChangeScope { workspaceId: string; groupId: string; lineageId: string; operationId: string }
export interface RecordingChangeState {
  unreadableScope?: UnreadableRecordingChangeScope;
  readback?: S.WriteGroupSuccessorOperationStatus;
  isOpen: boolean;
  busy: boolean;
  phase: RecordingChangePhase;
  intent?: SavedRecordingActionIntent;
  lifecycle?: WriteGroupLifecycle;
  target?: WriteGroup;
  review?: S.WriteGroupSuccessorReview;
  recovery?: RecordingChangeRecovery;
  preview?: SchemaPreviewToken;
  previewStage?: S.WriteGroupSuccessorStageRequest;
  activation?: S.WriteGroupSuccessorActivationReadiness;
  activationReview?: S.WriteGroupSuccessorActivationReview;
  boundary: string;
  error?: 'unavailable' | 'stale' | 'failed' | 'unknown' | 'status_failed' | 'boundary' | 'storage' | 'unreadable';
  unknown: boolean;
  canRetry: boolean;
  canRepreview: boolean;
}
