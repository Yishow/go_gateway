import type { ReactNode } from 'react';
import type { WriteGroup, WriteGroupDraft } from '../../../../../types/studioV2WriteGroup';
import type { CandidateCounterStatus } from '../../state/circuitProfileBasic';
import type { RecordingOwnerObservationState } from '../../state/recordingOwnerObservation';
import type { RecordingScopeView, RenderRecordingScopeSummary } from './recordingScopeSummary';
import type { SchemaPreparationPresentation } from './writeGroup/schemaPreparationPresentation';

/** Status is an observation of the existing owner; callbacks retain all requests and authority there. */
export interface RecordingStartPresentation {
  readonly observation: RecordingOwnerObservationState;
  readonly busy: boolean;
  readonly prepare: () => Promise<void> | void;
  readonly start: () => Promise<void> | void;
  readonly check: () => Promise<void> | void;
}
export interface RecordingPreparationBridge {
  readonly online: boolean;
  /** Ephemeral parent consent lease; never operation or readiness authority. */
  readonly consentEpoch?: number;
  readonly newCommandsBlocked?: boolean;
  readonly blockedGroupIds?: ReadonlySet<string>;
  readonly otherOwnerBlockedGroupIds?: ReadonlySet<string>;
  readonly isNewOperationBlocked?: (groupId?: string) => boolean;
  /** Denial only for a known preview-only failure; all other owners still participate. */
  readonly isSchemaPreviewRetryBlocked?: (groupId: string) => boolean;
  readonly onPresentation: (groupId: string, report: SchemaPreparationPresentation | null) => void;
}
export interface RecordingPreparationCommandFence {
  readonly ownerKey: string;
  readonly scopeKey: string;
  readonly blocked: boolean;
  /** Local report alignment only, never a replacement for native command guards. */
  readonly handoff?: { readonly workspaceRevision: string; readonly consentKey: string; readonly state: () => 'settling' | 'current' | 'blocked' };
  /** Denial only. False cannot replace any native command guard or server proof. */
  readonly isBlocked: (groupId?: string) => boolean;
}
export interface GuidedRecordingPreparationProps {
  readonly scope: RecordingScopeView;
  readonly ownerKey: string;
  /** Native owner controls, used only to invalidate consent; not persisted or displayed as an operation. */
  readonly consentKey: string;
  readonly fenceKey?: string;
  readonly batchHandoff?: { readonly consentKey: string; readonly blocked: boolean };
  readonly onCommandFence?: (fence: RecordingPreparationCommandFence | null) => void;
  readonly groups: readonly WriteGroup[];
  readonly drafts: readonly WriteGroupDraft[];
  readonly dirtyGroupIds: ReadonlySet<string>;
  readonly counterStatuses?: readonly CandidateCounterStatus[];
  readonly canPrepare: boolean;
  readonly isNewCommandBlocked?: () => boolean;
  readonly readonly: boolean;
  readonly start: RecordingStartPresentation;
  readonly prepareOnly?: boolean;
  readonly summary?: RenderRecordingScopeSummary;
  /** Read-only result slot; only an already confirmed same-O1 continuation replaces a redundant start review. */
  readonly completedOperationId?: string;
  readonly renderStartResult?: (primary: boolean) => ReactNode;
  readonly renderPreparation: (bridge: RecordingPreparationBridge) => ReactNode;
  readonly onInspect?: () => void;
  readonly onRepair?: () => void;
  readonly testIdPrefix: 'circuit-profile' | 'basic-recording';
}
