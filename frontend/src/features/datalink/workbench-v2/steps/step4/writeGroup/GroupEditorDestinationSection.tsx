import type * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { StudioV2LocatorReadScope } from '../../../../../../services/studioV2WorkspaceDatabase';
import type { WriteGroupCalculationCapabilities } from '../../../../../../types/studioV2WriteGroupCalculation';
import type { WriteGroup, WriteGroupStorageStrategy } from '../../../../../../types/studioV2WriteGroup';
import { ConfiguredDestinationChooser, type ConfiguredDestinationObservation } from '../../../recording-editor/ConfiguredDestinationChooser';
import { configuredDestinationPatch, type ConfiguredDestinationContext } from '../../../recording-editor/configuredRecordingDestinations';
import { RecordingMixedEditor } from '../../../recording-editor/RecordingMixedEditor';
import { RecordingDestinationMismatch, RecordingNewDestinationChooser, RecordingSavedDestination, type NewRecordingDestination } from '../../../recording-editor/RecordingNewDestination';
import type { DbConnector } from '../../../state/types';
import type { GroupCandidate } from '../../../state/writeGroup/candidates';
import type { EditorDraft } from '../../../state/writeGroup/draft';
import type { Step4TargetMetadataStatus } from '../useStep4TargetColumns';
import { destinationLock } from './GroupEditorCapabilities';
import { createNewRecordingFromLegacy } from './GroupEditorMethodSection';
import { RecordingDestinationFields } from './RecordingDestinationFields';
import { RecordingDestinationLockNotice } from './RecordingDestinationLockNotice';

/** Configured-destination choice observed by the editor; the choice itself stays editor state. */
export interface GroupEditorConfiguredDestination {
  observation?: ConfiguredDestinationObservation;
  context: ConfiguredDestinationContext;
  choice?: WriteGroup;
  invalid: boolean;
  reason?: string;
  onChoose: (choice: WriteGroup | undefined) => void;
}

export interface GroupEditorDestinationSectionProps {
  presentation: 'guided' | 'simple';
  scopedPresentation: boolean;
  workspaceId: string;
  workspaceRevision: string;
  connector: DbConnector;
  group: WriteGroup | undefined;
  locatorScope?: StudioV2LocatorReadScope;
  /** The connector identity is persisted and saved. */
  saved: boolean;
  destinationMatches: boolean;
  draft: EditorDraft;
  readonly: boolean;
  /** A save is in flight. */
  disabled: boolean;
  saveUnknown: boolean;
  saveReadbackPending?: boolean;
  metadataStatus: Step4TargetMetadataStatus;
  onRetryMetadata: () => void;
  newDestination?: NewRecordingDestination;
  /** Host-owned destination form, rendered here without adding a save owner. */
  destinationContent?: React.ReactNode;
  configured: GroupEditorConfiguredDestination;
  ownedWideConflict: boolean;
  calculationCapabilities?: WriteGroupCalculationCapabilities;
  /** Evaluated at click time so an in-flight save or a just-saved group refuses a new connector. */
  canChooseNewDestination: () => boolean;
  onChange: (patch: Partial<EditorDraft>) => void;
  onStorageStrategy: (strategy: WriteGroupStorageStrategy) => void;
  onReplaceDraft: (next: EditorDraft) => void;
  candidates?: GroupCandidate[];
  /** Opens a new unsaved recording prefilled with the locked profile recording's points and connector. */
  onStartNew?: (next: EditorDraft) => void;
  /** Existing copy-as-new action, already behind the editor's navigation guard. */
  onCopy?: () => void;
  copyDisabled?: boolean;
}

/**
 * The "destination" section of the record editor. It only assembles the
 * existing destination views; locks come from `destinationLock` and the editor
 * remains the single draft and save owner.
 */
export function GroupEditorDestinationSection(props: GroupEditorDestinationSectionProps) {
  const { presentation, scopedPresentation, workspaceId, workspaceRevision, connector, group, locatorScope, saved, destinationMatches, draft, readonly, disabled,
    saveUnknown, saveReadbackPending, metadataStatus, onRetryMetadata, newDestination, destinationContent, configured, ownedWideConflict, calculationCapabilities,
    canChooseNewDestination, onChange, onStorageStrategy, onReplaceDraft, candidates = [], onStartNew, onCopy, copyDisabled } = props;
  const { t } = useTranslation('workbench-v2');
  const { t: editorText } = useTranslation('recording-editor');
  const lock = destinationLock(draft, Boolean(group));
  const observation = configured.observation;
  const lockedProfileCopy = !onCopy && group && onStartNew
    && (draft.layout_version === 'managed_long_v1' || draft.layout_version === 'managed_wide_profile_v1')
    && (draft.recording_profile?.kind === 'instantaneous' || draft.recording_profile?.kind === 'cumulative')
    ? () => onStartNew(createNewRecordingFromLegacy(draft, candidates, calculationCapabilities, presentation, editorText('copy.name', { name: draft.name })))
    : undefined;
  return <>
    {presentation === 'simple' && <RecordingDestinationLockNotice draft={draft} saved={Boolean(group)} onCopy={onCopy ?? lockedProfileCopy} copyDisabled={copyDisabled} />}
    {presentation === 'simple' && draft.storage_strategy === 'managed' && metadataStatus !== 'exists' && <div className="recording-compact-members__metadata" data-testid="recording-destination-metadata-note"><p role="status">{editorText(`compact.destinationMetadata.${metadataStatus}`)}</p>{(metadataStatus === 'failed' || metadataStatus === 'forbidden') && <button type="button" className="gw-secondary rounded-md border border-[var(--gw-control)] bg-[var(--gw-surface)] px-3 py-1.5 text-sm text-[var(--gw-text)] hover:bg-[var(--gw-ground)]" onClick={onRetryMetadata}>{t('step4.group.retry')}</button>}</div>}
    {presentation === 'simple' && group ? <RecordingSavedDestination group={group} connector={connector} workspaceId={workspaceId} /> : newDestination && !newDestination.initialSetup ? !group && <RecordingNewDestinationChooser value={newDestination} selectedId={draft.connector_id}
      disabled={readonly || disabled || saveUnknown || Boolean(saveReadbackPending)} onSelect={id => {
        if (!canChooseNewDestination()) return;
        const next = newDestination.choose(id); if (!next) return;
        configured.onChoose(undefined); onChange({ connector_id: next.connector_id ?? '', connector_revision: next.identity_revision ?? '', table_schema: next.schema, table_name: '' });
      }} /> : destinationContent}
    {!destinationMatches && (!newDestination || newDestination.initialSetup) && <RecordingDestinationMismatch locatorUnknown={Boolean(locatorScope && !saved)} saved={saved} readonly={readonly}
      onUseCurrent={() => onChange({ connector_id: connector.connector_id ?? '', connector_revision: connector.identity_revision ?? '', table_schema: connector.schema ?? '',
        table_name: group?.layout_version === 'managed_long_v1' && draft.storage_strategy === 'managed' ? '' : connector.table ?? '' })} />}

    {presentation === 'simple' && <>
      <RecordingDestinationFields workspaceId={workspaceId} workspaceRevision={workspaceRevision} connector={connector} group={group} locatorScope={locatorScope} configuredDestinations={observation} scopedPresentation={scopedPresentation} draft={draft} readonly={readonly} {...lock} onChange={onChange} onStorageStrategy={onStorageStrategy} />

      <RecordingMixedEditor section="destination" draft={draft} saved={Boolean(group)} disabled={readonly || disabled} capabilities={calculationCapabilities} dialect={connector.kind} onChange={onReplaceDraft} />
    </>}
    {observation && <ConfiguredDestinationChooser observation={observation} context={configured.context} workspaceRevision={workspaceRevision}
      disabled={readonly || disabled} selected={configured.choice} onSelect={choice => { configured.onChoose(choice); onChange(configuredDestinationPatch(choice)); }}
      onReset={() => { configured.onChoose(undefined); onChange({ table_name: group?.destination.table_name ?? '' }); }} />}
    {configured.invalid && <p role="alert">{editorText(observation?.failed ? 'destinations.error' : observation?.loading ? 'destinations.loading' : configured.reason ? `destinations.reason.${configured.reason}` : 'destinations.stale')}</p>}{ownedWideConflict && <p role="alert">{editorText('destinations.reason.wideOwner')} {editorText('newDestination.catalogOwned')}</p>}
  </>;
}
