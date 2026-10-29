import type { ReactNode } from 'react';
import type { RecordingOwnerObservationState } from '../../../state/recordingOwnerObservation';

/** Ephemeral owner identity, not a readiness proof or a mutation request. */
export interface SchemaPreparationScope {
  readonly workspaceId: string;
  readonly workspaceRevision: string;
  readonly groupId?: string;
  readonly groupRevision?: string;
  readonly connectorId?: string;
  readonly connectorRevision?: string;
  /** Monotonic for this mounted owner; A → B → A never revives an old command. */
  readonly epoch: number;
}
export type SchemaPreparationBlocker = 'missing_scope' | 'scope_mismatch' | 'unreadable' | 'offline' | 'readonly'
  | 'dirty' | 'destination' | 'deleted' | 'busy' | 'operation_unresolved' | 'preview_missing';
export interface SchemaPreparationPresentation {
  readonly scope: SchemaPreparationScope;
  readonly observation: { readonly owner: 'schema' } & RecordingOwnerObservationState;
  readonly phase: 'idle' | 'previewed' | 'unresolved' | 'resolved' | 'unavailable';
  readonly blocker?: SchemaPreparationBlocker;
  readonly preview: 'none' | 'ready';
  /** Opaque identity of only the current ready preview; not an unresolved-operation observation. */
  readonly previewOperationId?: string;
  /** A failed preview command in this exact mounted owner, not an unknown DDL operation. */
  readonly previewFailure?: true;
  /** Each invocation revalidates the current owner. These are commands, never authority. */
  readonly commands: {
    readonly preview?: () => Promise<void>;
    readonly retryPreview?: () => Promise<void>;
    /** Explicit operator consent for the exact preview still held by this owner. */
    readonly confirm?: () => Promise<void>;
    readonly check?: () => Promise<void>;
  };
}
export interface GuidedSchemaPresentation {
  readonly online: boolean;
  /** Parent-native consent generation; publication identity only, never mutation permission. */
  readonly consentEpoch?: number;
  readonly newOperationBlocked?: boolean;
  readonly isNewOperationBlocked?: () => boolean;
  /** Denial-only delegate for this preview failure; absence or a thrown error blocks retry. */
  readonly isPreviewRetryBlocked?: () => boolean;
  readonly suppressDefaultActions?: boolean;
  /** Use one stable callback per owner; null retires that callback's own report. */
  readonly onPresentation?: (presentation: SchemaPreparationPresentation | null) => void;
  readonly renderActions?: (presentation: SchemaPreparationPresentation) => ReactNode;
}
