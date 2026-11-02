import { RecordingDestinationMismatch, RecordingNewDestinationChooser, RecordingSavedDestination } from '../../../recording-editor/RecordingNewDestination';
import { ConfiguredDestinationChooser, type ConfiguredDestinationObservation } from '../../../recording-editor/ConfiguredDestinationChooser';
import { configuredDestinationPatch, configuredDestinationReason, configuredDestinationStillCurrent } from '../../../recording-editor/configuredRecordingDestinations';
import { RecordingSourceNavigation } from '../../../recording-editor/RecordingSourceNavigation';
import { useWriteGroupSaveActivity, writeGroupSaveOutcomeUnknown } from '../../../../../../hooks/datalink/useWriteGroupSaveActivity';
import { GroupEditorWithCapabilities, destinationSaved, type GroupEditorContentProps } from './GroupEditorCapabilities';
import { RecordingDraftBoundary, useRecordingDraftGuard, useRecordingReadUnavailable } from '../../../recording-editor/RecordingDraftNavigation';
import type { WriteGroupCalculationCapabilities } from '../../../../../../types/studioV2WriteGroupCalculation';
import { RecordingMixedEditor, RecordingMixedResults, recordingMixedEditorIssues } from '../../../recording-editor/RecordingMixedEditor';
import { createRecordingCopyDraft, type RecordingCopyPrefill } from '../../../recording-editor/recordingCopyDraft';
import { RecordingEditorSummary } from '../../../recording-editor/RecordingEditorSummary';
import { RecordingEditorLayout } from '../../../recording-editor/RecordingEditorLayout';
import { mixedDraftToRequestGroup, newMixedDraft } from '../../../recording-editor/recordingMixedDraft';
import { initializeMixedMembers } from '../../../recording-editor/recordingMemberCalculation';
import { RecordingMemberCalculationSource } from '../../../recording-editor/RecordingMemberCalculation';
import { RecordingPointOrder } from '../../../recording-editor/RecordingPointOrder';
import type { DeviceSourcePort } from '../../../../../../pages/datalink/workbench-v2/deviceSourceSubmission';
import { RecordingPointSelection } from '../../../recording-editor/RecordingPointSelection';
import { RecordingSelectedMembers } from '../../../recording-editor/RecordingSelectedMembers';
import type { RecordingPickerDiscovery } from '../../../recording-editor/useRecordingPickerDiscovery';
import { readRecordingHomeSources } from '../../../recording-home/recordingHomeSources';
import type { StudioV2LocatorReadScope } from '../../../../../../services/studioV2WorkspaceDatabase';
import * as React from 'react';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  useCreateWriteGroupMutation,
  useUpdateWriteGroupMutation,
  useWriteGroupReadinessQuery,
} from '../../../../../../hooks/datalink/useStudioV2WriteGroups';
import type { WriteGroup, WriteGroupStorageStrategy, WriteGroupMutationResponse } from '../../../../../../types/studioV2WriteGroup';
import { normalizeTypedEnvelope } from '../../../../../../utils/safeJson';
import { getSafeErrorMessage } from '../../../../../../utils/typedErrors';
import type { DbConnector } from '../../../state/types';
import type { ExcludedCandidate, GroupCandidate } from '../../../state/writeGroup/candidates';
import { suggestColumns } from '../../../state/writeGroup/columns';
import {
  draftIsDirty, draftToRequestGroup, groupToDraft, newDraft, validateDraft,
  type EditorDraft, type EditorMember,
} from '../../../state/writeGroup/draft';
import { useStep4TargetColumns } from '../useStep4TargetColumns';
import { GroupColumnProposal } from './GroupColumnProposal';
import { GroupDeliveryStrip } from './GroupDeliveryStrip';
import { GroupIdentityFields } from './GroupIdentityFields';
import { GroupLifecycleBar } from './GroupLifecycleBar';
import type { SavedRecordingActionHost } from './lifecycle/SavedRecordingActions';
import { GroupMemberTable } from './GroupMemberTable';
import { GroupPolicyFields } from './GroupPolicyFields';
import { GroupProfileMembers, GroupProfileSummary } from './GroupProfileSummary';
import { GroupReadinessPanel } from './GroupReadinessPanel';
import { GroupEditorPreparation } from '../../../recording-editor/GroupEditorPreparation';
import { GroupTestWritePanel } from './GroupTestWritePanel';

export interface GroupEditorProps {
  presentation?: 'guided' | 'simple';
  copyPrefill?: RecordingCopyPrefill;
  onCopy?: () => void;
  /** One host-owned destination form, rendered here without adding a save owner. */
  destinationContent?: React.ReactNode;
  repairContent?: React.ReactNode;
  configuredDestinations?: ConfiguredDestinationObservation;
  calculationCapabilities?: WriteGroupCalculationCapabilities;
  deviceSource?: DeviceSourcePort;
  sourceDevices?: readonly { id: string; name: string }[];
  locatorScope?: StudioV2LocatorReadScope;
  /** Denial-only incoming navigation fence; never a read scope or write permission. */
  locatorGuard?: boolean;
  /** Denial-only fence until the preceding locator save has read back. */
  saveReadbackPending?: boolean;
  /** Undefined creates a new group. */
  group: WriteGroup | undefined;
  workspaceId: string;
  workspaceRevision: string;
  connector: DbConnector;
  candidates: GroupCandidate[];
  excluded: ExcludedCandidate[];
  readonly: boolean;
  /** Tombstone fields are immutable, but lifecycle review remains available. */
  editorReadonly?: boolean;
  lifecycleHost?: SavedRecordingActionHost;
  onClose: () => void;
  onSaved: (group: WriteGroup, workspaceRevision?: string, ownSaveAcknowledgement?: WriteGroupMutationResponse) => void;
  onReload: () => void;
  /**
   * Reads the current workspace revision and group just before a save. Other
   * setup saves (devices, rules, Tags) move the workspace revision without
   * touching this group, so a stale copy would be refused for no real conflict.
   */
  fetchLatest: () => Promise<{ workspace_revision: string; groups: WriteGroup[] } | undefined>;
  /** Sends a managed profile group's effective action to the host's reviewed scoped start. */
  onReviewScopedStart?: (group: WriteGroup) => void;
}

/**
 * One saved group: its name, destination table, row rules and members. The
 * same editor serves managed and custom storage; the backend stays the
 * authority for readiness, so Apply follows the server's answer for the saved
 * revision and never the local checks alone.
 */
export const GroupEditor: React.FC<GroupEditorProps> = props => <RecordingDraftBoundary><RecordingSourceNavigation workspaceId={props.workspaceId} source={props.deviceSource}><GroupEditorWithCapabilities {...props} Content={GroupEditorContent} /></RecordingSourceNavigation></RecordingDraftBoundary>;
const GroupEditorContent: React.FC<GroupEditorContentProps> = ({
  presentation = 'guided', newDestination, copyPrefill, onCopy, destinationContent, repairContent, configuredDestinations, calculationCapabilities, deviceSource, sourceDevices = [], group: propGroup, workspaceId, workspaceRevision, connector, candidates: suppliedCandidates, excluded, readonly: operatorReadonly, editorReadonly = false, lifecycleHost, onClose, onSaved, onReload, fetchLatest, onReviewScopedStart, mysqlCapability, sqlserverCapability, locatorScope, locatorGuard, saveReadbackPending,
}) => {
  const { t } = useTranslation('workbench-v2');
  const { t: editorText } = useTranslation('recording-editor');
  const readUnavailable = useRecordingReadUnavailable();
  const observedPickerScope = JSON.stringify([workspaceId, workspaceRevision, propGroup?.id ?? 'new', propGroup?.revision]);
  const ownPickerAcknowledgement = useRef<WriteGroupMutationResponse | undefined>(undefined);
  const pickerIdentity = useRef({ observed: observedPickerScope, scope: observedPickerScope, workspaceRevision });
  if (pickerIdentity.current.observed !== observedPickerScope) {
    const acknowledgement = ownPickerAcknowledgement.current;
    const ownReadback = locatorScope && acknowledgement && acknowledgement.group.workspace_id === workspaceId &&
      acknowledgement.group.id === propGroup?.id && acknowledgement.group.revision === propGroup?.revision &&
      (workspaceRevision === acknowledgement.workspace_revision || workspaceRevision === pickerIdentity.current.workspaceRevision);
    if (!ownReadback) ownPickerAcknowledgement.current = undefined;
    // Only our acknowledged revision handoff keeps this local form session.
    // Unrelated scope changes retain the existing invalidation behavior.
    pickerIdentity.current = { observed: observedPickerScope, scope: ownReadback ? pickerIdentity.current.scope : observedPickerScope, workspaceRevision };
  }
  const pickerScope = pickerIdentity.current.scope;
  const [observedSources, setObservedSources] = useState<{ scope: string; value: RecordingPickerDiscovery }>();
  const candidates = useMemo(() => observedSources?.scope === pickerScope ? [...observedSources.value.candidates] : suppliedCandidates, [observedSources, pickerScope, suppliedCandidates]);
  const readonly = operatorReadonly || readUnavailable || editorReadonly || propGroup?.status === 'deleted';
  const renderCalculation = (member: EditorMember) => <RecordingMemberCalculationSource member={member} disabled={readonly || saveBusy} compact={presentation === 'simple'}
    workspaceRevision={workspaceRevision} locatorScope={locatorScope} context={{ workspaceId, recordId: draft.record_id, capabilities: calculationCapabilities, layout: draft.layout_version ?? '', dialect: connector.kind }}
    onChange={next => set({ members: draft.members.map(value => value === member ? next : value) })} />;
  const renderPointOrder = (member: EditorMember) => <>{deviceSource?.pointOrder && <RecordingPointOrder port={deviceSource.pointOrder}
    workspaceId={workspaceId} deviceId={member.device_id} pointId={member.point_id} disabled={readonly || saveBusy} />}{mixed && presentation !== 'simple' && renderCalculation(member)}</>;
  const saved = destinationSaved(connector);
  const isMySQL = String(connector.kind).toLowerCase() === 'mysql';
  const mysqlConnectorId = connector.connector_id ?? '';
  const mysqlConnectorRevision = connector.identity_revision ?? '';
  const mysqlTestWriteAvailable = !isMySQL || Boolean(saved && mysqlCapability &&
    mysqlCapability.connector_id === mysqlConnectorId && mysqlCapability.connector_revision === mysqlConnectorRevision &&
    mysqlCapability.connected && mysqlCapability.metadata_status === 'verified' && mysqlCapability.supports_group_test_write);
  const isSQLServer = String(connector.kind).toLowerCase() === 'sqlserver';
  const sqlserverConnectorId = connector.connector_id ?? '';
  const sqlserverConnectorRevision = connector.identity_revision ?? '';
  const sqlserverTestWriteAvailable = !isSQLServer || Boolean(saved && sqlserverCapability &&
    sqlserverCapability.connector_id === sqlserverConnectorId && sqlserverCapability.connector_revision === sqlserverConnectorRevision &&
    sqlserverCapability.connected && sqlserverCapability.metadata_status === 'verified' && sqlserverCapability.supports_group_test_write);
  // Keep the just-returned canonical group visible while the list query is
  // refetching. This prevents stale list data from making server-assigned
  // managed table/columns look dirty again.
  const [editorGroup, setEditorGroup] = useState<WriteGroup | undefined>(propGroup);
  const previousPropGroup = useRef(propGroup);
  const group = editorGroup;
  const groupRef = useRef(group); groupRef.current = group;
  const [draft, setDraft] = useState<EditorDraft>(() => {
    const destination = { connector_id: connector.connector_id ?? '', connector_revision: connector.identity_revision ?? '', table_schema: connector.schema ?? '', table_name: connector.table ?? '' };
    // Explicit simple creation is local per-point intent, not a capability or permission to save/start.
    return group ? groupToDraft(group) : copyPrefill ? createRecordingCopyDraft(copyPrefill, destination, crypto.randomUUID(), editorText('copy.name', { name: copyPrefill.sourceName })) : presentation === 'simple' ? newMixedDraft({ ...destination, table_name: '' }, crypto.randomUUID()) : newDraft(destination);
  });
  const [configuredChoice, setConfiguredChoice] = useState<WriteGroup>();
  const configuredContext = { workspaceId, connector, draft, groupId: group?.id };
  const configuredInvalid = Boolean(configuredChoice && (configuredDestinationReason(configuredChoice, configuredContext) || configuredDestinations?.failed || configuredDestinations?.loading || !configuredDestinationStillCurrent(configuredChoice, configuredDestinations?.data)));
  const configuredReason = configuredChoice && configuredDestinationReason(configuredChoice, configuredContext);
  const [saveError, setSaveError] = useState<unknown>(null);
  const [conflict, setConflict] = useState(false);
  const createMutation = useCreateWriteGroupMutation();
  const updateMutation = useUpdateWriteGroupMutation();
  const savingRef = useRef(false), saveOperation = useRef<object | undefined>(undefined);
  const saveActivity = useWriteGroupSaveActivity(workspaceId, group?.id);
  const [saveUnknown, setSaveUnknown] = useState(false);
  const [ownSaveBusy, setSaveBusy] = useState(false);
  const saveBusy = ownSaveBusy || saveActivity.pending;
  const customDedupeRef = useRef(draft.dedupe_capability);
  const mountedRef = useRef(false);
  useEffect(() => {
    mountedRef.current = true;
    return () => { mountedRef.current = false; };
  }, []);

  // Legacy profiles keep their canonical group policy; mixed members own their individual cadence.
  const profile = (draft.layout_version === 'managed_long_v1' || draft.layout_version === 'managed_wide_profile_v1') && draft.recording_profile !== undefined;
  const mixed = draft.recording_profile?.kind === 'mixed';
  const columnsAssignedOnSave = presentation === 'simple' && !group && draft.storage_strategy === 'managed';
  const dirty = group ? draftIsDirty(draft, group) : true;
  const saveScopeKey = JSON.stringify([workspaceId, propGroup?.id, connector.connector_id, connector.identity_revision, readonly]);
  const saveScope = useRef({ key: saveScopeKey });
  if (saveScope.current.key !== saveScopeKey) saveScope.current = { key: saveScopeKey };
  const navigation = useRecordingDraftGuard({ version: JSON.stringify([workspaceId, propGroup?.id, draft]), dirty, pending: saveBusy || saveUnknown,
    discard: () => { saveScope.current = { key: saveScopeKey }; } });
  const closeEditor = () => navigation.request(onClose);
  // Adopt a new server revision only while there is nothing unsaved to lose.
  const revision = group?.revision;
  const dirtyRef = useRef(dirty);
  dirtyRef.current = dirty;
  useEffect(() => {
    if (previousPropGroup.current !== propGroup && !dirtyRef.current) {
      previousPropGroup.current = propGroup;
      setEditorGroup(propGroup);
    }
  }, [propGroup]);
  useEffect(() => {
    if (group && !dirtyRef.current) {
      const next = groupToDraft(group);
      customDedupeRef.current = next.dedupe_capability;
      setDraft(next);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [revision]);

  const destinationMatches = draft.connector_id === connector.connector_id &&
    draft.connector_revision === connector.identity_revision;
  const scopedConnector = useMemo<DbConnector>(() => ({
    ...connector, table: draft.table_name, schema: draft.table_schema,
  }), [connector, draft.table_name, draft.table_schema]);
  const tableSaved = group
    ? draft.table_name === group.destination.table_name && draft.table_schema === group.destination.table_schema
    : draft.table_name === connector.table && draft.table_schema === connector.schema;
  const metadata = useStep4TargetColumns(destinationMatches && tableSaved && (draft.storage_strategy === 'custom' || Boolean(group)) ? scopedConnector : { ...scopedConnector, persisted: false },
    group ? { groupId: group.id, groupRevision: group.revision } : undefined, locatorScope);
  // A new managed group has no target table to inspect. Present an empty
  // server-owned column view so the editor never asks for manual assignments.
  const editorMetadata = draft.storage_strategy === 'managed' && !group
    ? { ...metadata, status: 'exists' as const, columns: [] }
    : metadata;
  const dialect = connector.kind === 'postgres' ? 'postgres' : connector.kind === 'sqlite' ? 'sqlite' : 'portable';
  const confirmed = useMemo(
    () => Object.fromEntries(draft.members.filter((m) => m.target_column).map((m) => [m.key, m.target_column])),
    [draft.members],
  );
  const suggestions = useMemo(
    () => (editorMetadata.status === 'exists' ? suggestColumns(candidates, editorMetadata.columns, confirmed, dialect, draft.storage_strategy === 'managed') : { assignments: {}, unmatched: candidates.map((c) => c.key) }),
    [candidates, editorMetadata.columns, editorMetadata.status, confirmed, dialect, draft.storage_strategy],
  );
  const sourceOwnerRef = useRef(deviceSource); sourceOwnerRef.current = deviceSource;
  const unsettledSources = draft.members.filter(member => deviceSource?.range?.sourceWriteState?.(member.device_id, member.point_id));
  const sourceWritePending = () => draft.members.some(member => sourceOwnerRef.current?.range?.sourceWriteState?.(member.device_id, member.point_id));
  const issues = mixed ? [] : validateDraft(draft, { destinationSaved: saved && destinationMatches, storageStrategy: draft.storage_strategy });
  const mixedIssues = mixed ? recordingMixedEditorIssues(draft, calculationCapabilities, connector.kind, saved && destinationMatches) : [];
  const invalid = issues.length > 0 || mixedIssues.length > 0 || configuredInvalid || Boolean(!group && newDestination && !newDestination.current());
  const readiness = useWriteGroupReadinessQuery(group?.id, Boolean(group) && !dirty);
  const set = (patch: Partial<EditorDraft>) => { if (!savingRef.current) setDraft((current) => ({ ...current, ...patch })); };
  const setMembers = (members: EditorMember[]) => set({ members: mixed && (calculationCapabilities || presentation === 'simple') ? initializeMixedMembers(members, draft.members, candidates, calculationCapabilities,
    presentation === 'simple' ? { interval: 15, unit: 'second', alignment: 'utc_epoch' } : undefined) : members });
  const changeStorageStrategy = (storage_strategy: WriteGroupStorageStrategy) => setDraft((current) => {
    if (savingRef.current || storage_strategy === current.storage_strategy) return current;
    if (storage_strategy === 'managed') {
      customDedupeRef.current = current.dedupe_capability;
      return { ...current, storage_strategy, dedupe_capability: 'receipt' };
    }
    return { ...current, storage_strategy, dedupe_capability: customDedupeRef.current };
  });

  // A locator read grants no enduring command lease. Suspension or owner changes cannot be undone by returning to A.
  const locatorLeaseKey = locatorScope || locatorGuard ? JSON.stringify([locatorScope?.workspaceId, locatorScope?.workspaceRevision, workspaceId,
    workspaceRevision, group?.id, group?.revision, readonly]) : null;
  const locatorLease = useRef({ key: locatorLeaseKey, blocked: readonly });
  if (locatorLease.current.key !== locatorLeaseKey) locatorLease.current = { key: locatorLeaseKey, blocked: readonly };
  const renderLocatorLease = locatorLease.current;

  const handleSave = async () => {
    if (savingRef.current || saveActivity.isPending() || saveUnknown || saveReadbackPending || invalid || !group && newDestination && !newDestination.current() || sourceWritePending() || readonly || locatorLease.current !== renderLocatorLease || renderLocatorLease.key !== null &&
      (!mountedRef.current || locatorLease.current.blocked)) return;
    savingRef.current = true;
    const operation = {}; saveOperation.current = operation;
    setSaveBusy(true);
    const requestLease = renderLocatorLease;
    const requestScope = saveScope.current;
    setSaveError(null);
    let submitted = false;
    try {
      const body = mixed && calculationCapabilities ? mixedDraftToRequestGroup(draft, workspaceId, calculationCapabilities) : draftToRequestGroup(draft, workspaceId);
      const latest = await fetchLatest();
      if (!mountedRef.current || saveScope.current !== requestScope || locatorLease.current !== requestLease || requestLease.key !== null && locatorLease.current.blocked) return;
      if (sourceWritePending() || !group && newDestination && !newDestination.current()) return;
      if (configuredChoice && (!configuredDestinationStillCurrent(configuredChoice, latest) || configuredDestinationReason(configuredChoice, configuredContext))) throw new Error(editorText('destinations.stale'));
      const freshRevision = latest?.workspace_revision ?? workspaceRevision;
      if (group) {
        const current = latest?.groups.find((entry) => entry.id === group.id);
        if (current && current.revision !== group.revision) {
          // This group itself changed elsewhere: that is a real conflict.
          setConflict(true);
          return;
        }
      }
      submitted = true;
      const result = group
        ? await updateMutation.mutateAsync({ id: group.id, request: {
          workspace_id: workspaceId, expected_workspace_revision: freshRevision, expected_group_revision: group.revision,
          expected_connector_revision: draft.connector_revision, group: body,
        } })
        : await createMutation.mutateAsync({
          workspace_id: workspaceId, expected_workspace_revision: freshRevision,
          expected_connector_revision: draft.connector_revision, group: body,
      });
      if (mountedRef.current && saveScope.current === requestScope) {
        const acknowledgement = locatorLease.current === requestLease && !requestLease.blocked &&
          result.group.workspace_id === workspaceId && (!group || result.group.id === group.id) ? result : undefined;
        ownPickerAcknowledgement.current = acknowledgement;
        // Adopt every server-assigned managed column and binding before the
        // parent refetches. Server proof stays outside the editable draft.
        setConfiguredChoice(undefined);
        setDraft(groupToDraft(result.group));
        groupRef.current = result.group; setEditorGroup(result.group);
        setConflict(false);
        onSaved(result.group, result.workspace_revision, acknowledgement);
      }
    } catch (cause) {
      if (!mountedRef.current || saveScope.current !== requestScope) return;
      setSaveError(cause);
      if (submitted && writeGroupSaveOutcomeUnknown(cause)) setSaveUnknown(true);
      if (normalizeTypedEnvelope(cause).code === 'revision_mismatch') setConflict(true);
    } finally {
      if (saveOperation.current === operation) {
        savingRef.current = false; saveOperation.current = undefined;
        if (mountedRef.current) setSaveBusy(false);
      }
    }
  };

  const saving = saveBusy || createMutation.isPending || updateMutation.isPending;
  const canApply = Boolean(group) && !dirty && !invalid && (!mixed || calculationCapabilities?.activation_supported === true) && readiness.data?.ready === true &&
    readiness.data.group_revision === group?.revision;
  const applyBlocked = !group ? t('step4.group.apply_blocked_unsaved')
    : dirty ? t('step4.group.apply_blocked_dirty')
      : readiness.isLoading ? t('step4.group.apply_blocked_loading')
        : readiness.isError ? t('step4.group.apply_blocked_unknown')
          : readiness.data && !readiness.data.ready ? t('step4.group.apply_blocked_not_ready') : undefined;
  const mysqlSchemaReady = !isMySQL || Boolean(group && readiness.data?.group_revision === group.revision &&
    readiness.data.schema_ready && !readiness.isFetching && !readiness.isError);
  const sqlserverSchemaReady = !isSQLServer || Boolean(group && readiness.data?.group_revision === group.revision &&
    readiness.data.schema_ready && !readiness.isFetching && !readiness.isError);
  const testWriteEnabled = Boolean(group) && !readonly && !saveReadbackPending && !dirty && destinationMatches && group?.status !== 'deleted' &&
    mysqlTestWriteAvailable && mysqlSchemaReady && sqlserverTestWriteAvailable && sqlserverSchemaReady;
  const testBlocked = !group ? t('step4.group.test_write.blocked_unsaved')
    : dirty ? t('step4.group.test_write.blocked_dirty')
      : !destinationMatches ? t('step4.group.test_write.blocked_destination')
        : isMySQL && !mysqlTestWriteAvailable ? t('step4.group.test_write.blocked_capability')
          : isMySQL && readiness.isError ? t('step4.group.test_write.blocked_readiness_unknown')
            : isMySQL && (readiness.isLoading || readiness.isFetching) ? t('step4.group.test_write.blocked_readiness_loading')
              : isMySQL && !mysqlSchemaReady ? t('step4.group.test_write.blocked_schema_not_ready')
                : isSQLServer && !sqlserverTestWriteAvailable ? t('step4.group.test_write.blocked_sqlserver_capability')
                  : isSQLServer && readiness.isError ? t('step4.group.test_write.blocked_readiness_unknown')
                    : isSQLServer && (readiness.isLoading || readiness.isFetching) ? t('step4.group.test_write.blocked_readiness_loading')
                      : isSQLServer && !sqlserverSchemaReady ? t('step4.group.test_write.blocked_schema_not_ready') : undefined;

  const revisionDetails = group && (
    <dl className="grid grid-cols-2 gap-x-4 gap-y-1 text-xs text-[var(--gw-muted)] sm:grid-cols-4" data-testid="group-revisions">
      <dt>{t('step4.group.editor.status')}</dt><dd className="text-[var(--gw-text)]">{group.status}</dd>
      <dt>{t('step4.group.editor.draft_revision')}</dt><dd className="font-mono text-[var(--gw-text)]">{group.revision}</dd>
      <dt>{t('step4.group.editor.applied_revision')}</dt><dd className="font-mono text-[var(--gw-text)]" data-testid="group-applied-revision">{group.applied_revision || t('step4.group.editor.not_applied')}</dd>
      <dt>{t('step4.group.editor.connector_revision')}</dt><dd className="font-mono text-[var(--gw-text)]">{group.destination.connector_revision}</dd>
    </dl>
  );
  return (
    <div className={presentation === 'simple' ? 'recording-editor-simple min-w-0 space-y-5' : 'min-w-0 space-y-5 rounded-xl border border-[var(--gw-border)] p-5'} data-testid="group-editor"
      onKeyDown={event => { if (event.key === 'Escape' && !event.nativeEvent.isComposing && !event.defaultPrevented && !(event.target as HTMLElement).closest('[role="dialog"], [role="alertdialog"], [role="tooltip"], [data-state="open"]')) { event.preventDefault(); closeEditor(); } }}>
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-semibold text-[var(--gw-text)]">{group ? group.name : presentation === 'simple' ? editorText('simple.create') : t('step4.group.editor.new_title')}</h3>
        {presentation === 'simple' && group && onCopy && <button type="button" className="gw-secondary" disabled={readonly || saveBusy || saveUnknown || saveReadbackPending} onClick={() => navigation.request(onCopy)} data-testid="recording-copy">{editorText('copy.action')}</button>}
        {presentation !== 'simple' && <button type="button" onClick={closeEditor} className="text-xs text-[var(--gw-muted)] underline" data-testid="group-editor-close">{t('step4.group.editor.back')}</button>}
      </div>

      {!group && copyPrefill && <aside role="status" className="space-y-2 text-sm" data-testid="recording-copy-intent"><p>{editorText('copy.intent')}</p><p>{editorText('copy.concurrent')}</p>{copyPrefill.layout === 'managed_long_v1' && <p>{editorText('copy.longIsolation')}</p>}</aside>}
      {group && (presentation === 'simple' ? <details data-testid="group-revision-details">
        <summary><span>{editorText('versions.details')}</span> · {t('step4.group.editor.status')}: {group.status}</summary>
        {revisionDetails}
      </details> : revisionDetails)}

      <fieldset disabled={saveBusy || readUnavailable} className="min-w-0 space-y-5" data-testid="group-save-fields">
      <RecordingEditorLayout presentation={presentation} perPoint={mixed} memberCount={draft.members.length}
        summary={presentation === 'simple' && <RecordingEditorSummary draft={draft} destinationName={saved && destinationMatches ? connector.name || connector.kind : undefined} saved={Boolean(group)} dirty={dirty} />}
        identity={<GroupIdentityFields section={presentation === 'simple' ? 'name' : 'all'} draft={draft} readonly={readonly} locked={profile || draft.layout_version === 'managed_long_v1'} onChange={set} onStorageStrategy={changeStorageStrategy} />}
        method={<>
      <RecordingMixedEditor section={presentation === 'simple' ? 'method' : 'all'} draft={draft} saved={Boolean(group)} disabled={readonly || saveBusy} capabilities={calculationCapabilities} dialect={connector.kind} onChange={next => { if (!savingRef.current) setDraft(next); }} />
      {profile && !mixed && <GroupProfileSummary draft={draft} />}
        </>}
        destination={<>
      {presentation === 'simple' && group ? <RecordingSavedDestination group={group} connector={connector} workspaceId={workspaceId} /> : newDestination && !newDestination.initialSetup ? !group && <RecordingNewDestinationChooser value={newDestination} selectedId={draft.connector_id}
        disabled={readonly || saveBusy || saveUnknown || Boolean(saveReadbackPending)} onSelect={id => {
          if (groupRef.current || readonly || savingRef.current || saveUnknown || saveReadbackPending) return;
          const next = newDestination.choose(id); if (!next) return;
          setConfiguredChoice(undefined); set({ connector_id: next.connector_id ?? '', connector_revision: next.identity_revision ?? '', table_schema: next.schema, table_name: '' });
        }} /> : destinationContent}
      {!destinationMatches && (!newDestination || newDestination.initialSetup) && <RecordingDestinationMismatch locatorUnknown={Boolean(locatorScope && !saved)} saved={saved} readonly={readonly}
        onUseCurrent={() => set({ connector_id: connector.connector_id ?? '', connector_revision: connector.identity_revision ?? '', table_schema: connector.schema ?? '',
          table_name: group?.layout_version === 'managed_long_v1' && draft.storage_strategy === 'managed' ? '' : connector.table ?? '' })} />}

          {presentation === 'simple' && <>
      <GroupIdentityFields section="destination" draft={draft} readonly={readonly} locked={profile || draft.layout_version === 'managed_long_v1'} onChange={set} onStorageStrategy={changeStorageStrategy} />

      <RecordingMixedEditor section="destination" draft={draft} saved={Boolean(group)} disabled={readonly || saveBusy} capabilities={calculationCapabilities} dialect={connector.kind} onChange={next => { if (!savingRef.current) setDraft(next); }} />
          </>}
      {configuredDestinations && <ConfiguredDestinationChooser observation={configuredDestinations} context={configuredContext} workspaceRevision={workspaceRevision}
        disabled={readonly || saveBusy} selected={configuredChoice} onSelect={choice => { setConfiguredChoice(choice); set(configuredDestinationPatch(choice)); }}
        onReset={() => { setConfiguredChoice(undefined); set({ table_name: group?.destination.table_name ?? '' }); }} />}
      {configuredInvalid && <p role="alert">{editorText(configuredDestinations?.failed ? 'destinations.error' : configuredDestinations?.loading ? 'destinations.loading' : configuredReason ? `destinations.reason.${configuredReason}` : 'destinations.stale')}</p>}

        </>}
        policy={<GroupPolicyFields expanded={presentation === 'simple'} draft={draft} readonly={readonly} profile={profile} onChange={set} onDedupe={(capability) => { customDedupeRef.current = capability; set({ dedupe_capability: capability }); }} />}
        sources={profile && !mixed ? <GroupProfileMembers renderSource={renderPointOrder} draft={draft} candidates={candidates} locatorScope={locatorScope} workspaceId={workspaceId} readonly={readonly} /> : (
        <>
        {deviceSource ? <RecordingPointSelection key={pickerScope} scope={pickerScope} workspaceId={workspaceId} candidates={candidates}
          members={draft.members} devices={sourceDevices} disabled={readonly || saveBusy} deviceSource={deviceSource} onChange={setMembers}
          onDiscovery={setObservedSources} discover={locatorScope ? async () => {
            const read = await readRecordingHomeSources({ workspaceId, workspaceRevision, epoch: 0 });
            if (read.state !== 'available') throw new Error('Source observation unavailable');
            return { candidates: read.value.candidates, devices: read.value.devices };
          } : undefined}
          renderMembers={actions => <RecordingSelectedMembers presentation={presentation} columnsAssignedOnSave={columnsAssignedOnSave} renderCalculation={presentation === 'simple' && mixed ? renderCalculation : undefined} renderSource={renderPointOrder} {...actions} managed={draft.storage_strategy === 'managed'} dialect={dialect}
            candidates={candidates} members={draft.members} columns={editorMetadata.columns} metadataStatus={editorMetadata.status}
            suggestions={suggestions} issues={issues} showEntityKey={draft.entity_key_column.trim() !== ''} disabled={readonly || saveBusy}
            onChange={setMembers} labels={{ remove: member => editorText('picker.removeMember', { point: member.point_id }),
              unavailable: editorText('picker.unavailable'), undo: editorText('picker.undo') }} />} /> :
        <GroupMemberTable presentation={presentation} columnsAssignedOnSave={columnsAssignedOnSave} renderCalculation={presentation === 'simple' && mixed ? renderCalculation : undefined} renderSource={renderPointOrder} managed={draft.storage_strategy === 'managed'} dialect={dialect}
          candidates={candidates} excluded={excluded} members={draft.members} columns={editorMetadata.columns} metadataStatus={editorMetadata.status}
          suggestions={suggestions} issues={issues} showEntityKey={draft.entity_key_column.trim() !== ''} disabled={readonly || saveBusy}
          onChange={setMembers} />}
        {!columnsAssignedOnSave && <GroupColumnProposal
          dialect={dialect}
          managed={draft.storage_strategy === 'managed'}
          existingColumns={editorMetadata.columns.map((column) => column.name)}
          unmatched={candidates.filter((candidate) => {
            const member = draft.members.find((entry) => entry.key === candidate.key);
            if (!member) return false;
            const column = metadata.columns.find((entry) => entry.name === member.target_column);
            return editorMetadata.status === 'exists' && (!member.target_column || !column);
          })}
        />}
        {editorMetadata.status === 'failed' && (
          <button type="button" onClick={metadata.refetch} className="text-xs underline" data-testid="group-metadata-retry">{t('step4.group.retry')}</button>
        )}

        </>
      )} />

      </fieldset>
      {unsettledSources.length > 0 && <div role="alert" data-testid="group-source-unsettled">
        <p>{editorText('range.selectedUnsettled')}</p><ul>{unsettledSources.map(member => {
          const candidate = candidates.find(row => row.device_id === member.device_id && row.point_id === member.point_id);
          return <li key={member.key}>{candidate ? `${candidate.device_name} · ${candidate.label} · ${candidate.address}` : `${sourceDevices.find(device => device.id === member.device_id)?.name ?? member.device_id} · ${member.point_id}`}</li>;
        })}</ul>
      </div>}
      {issues.length > 0 && (
        <ul role="alert" className="list-disc pl-5 text-xs text-[var(--gw-danger-fg)]" data-testid="group-issues">
          {issues.map((issue, index) => <li key={`${issue.code}-${issue.member_key ?? issue.column ?? index}`}>{t(`step4.group.issue.${issue.code.replace(/-/g, '_')}`, { column: issue.column ?? '' })}</li>)}
        </ul>
      )}
      {mixedIssues.length > 0 && <ul role="alert" data-testid="mixed-group-issues">{mixedIssues.map((code, index) => <li key={`${code}-${index}`}>{editorText(`calculation.issue.${code}`, { defaultValue: editorText('calculation.unavailable') })}</li>)}</ul>}
      {conflict && (
        <div role="alert" className="text-xs text-[var(--gw-warning-fg)]" data-testid="group-conflict">
          {t('step4.group.editor.conflict')}
          <button type="button" onClick={() => { setConflict(false); onReload(); }} className="ml-2 underline" data-testid="group-conflict-reload">{t('step4.group.editor.reload')}</button>
        </div>
      )}
      {saveError !== null && !conflict && (
        <p role="alert" className="text-xs text-[var(--gw-danger-fg)]" data-testid="group-save-error">{getSafeErrorMessage(saveError, t).message}</p>
      )}

      {(saveUnknown || saveActivity.priorUnknown) && <div role="status" data-testid="group-save-unknown">
        <p>{editorText('navigation.unknownSave')}</p>
        <button type="button" className="rounded-md border border-[var(--gw-control)] bg-[var(--gw-surface)] px-3 py-2 text-sm text-[var(--gw-text)] hover:bg-[var(--gw-active)] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--gw-focus)]" data-testid="group-save-review" onClick={() => navigation.request(() => { onReload(); onClose(); })}>{editorText('navigation.reviewSaved')}</button>
      </div>}
      {repairContent}
      <GroupEditorPreparation managed={draft.storage_strategy === 'managed' || group?.destination.storage_strategy === 'managed'}
        group={group} workspaceId={workspaceId} workspaceRevision={workspaceRevision} destinationMatches={destinationMatches} dirty={dirty} readonly={readonly}
        consentIdentity={JSON.stringify([workspaceId, workspaceRevision, group?.id, group?.revision, draft, connector.connector_id, connector.identity_revision,
          candidates.filter(candidate => draft.members.some(member => member.point_id === candidate.point_id && member.tag_id === candidate.tag_id)), readonly, unsettledSources])}
        save={{ saved: Boolean(group), dirty, busy: saving, unknown: saveUnknown,
          disabled: readonly || saving || saveUnknown || saveReadbackPending || invalid || unsettledSources.length > 0 || (Boolean(group) && !dirty), invoke: handleSave }}
        blocked={readonly || saving || saveUnknown || saveReadbackPending || invalid || unsettledSources.length > 0}
        readiness={{ data: readiness.data, loading: readiness.isFetching, failed: readiness.isError, updatedAt: readiness.dataUpdatedAt, check: readiness.refetch }}
        onApplied={() => { metadata.refetch(); void readiness.refetch(); onReload(); }} />

      {group && (
        <>
          <section aria-label={t('step4.group.readiness.title')} className="space-y-2">
            <h4 className="text-sm font-semibold text-[var(--gw-text)]">{t('step4.group.readiness.title')}</h4>
            <GroupReadinessPanel readiness={readiness.data} loading={readiness.isLoading && !dirty} failed={readiness.isError} dirty={dirty} onRetry={() => void readiness.refetch()} />
          </section>
          <GroupDeliveryStrip group={group} readonly={readonly} />
          {mixed && <RecordingMixedResults group={group} candidates={candidates} />}
          <GroupLifecycleBar
            locatorScope={locatorScope} locatorGuard={locatorGuard} confirmStop={presentation === 'simple'}
            group={group} workspaceId={workspaceId} workspaceRevision={workspaceRevision}
            canApply={canApply} applyBlockedReason={applyBlocked} readonly={operatorReadonly || readUnavailable} lifecycleHost={lifecycleHost}
            onChanged={onSaved} onConflict={() => setConflict(true)} onReviewScopedStart={onReviewScopedStart}
          />
          <GroupTestWritePanel
            locatorScope={locatorScope}
            workspaceId={workspaceId}
            groupId={group.id}
            scopeKey={`${group.id}|${group.revision}`}
            enabled={testWriteEnabled} newOperationBlocked={Boolean(saveReadbackPending)}
            disabledReason={readonly ? t('step4.group.test_write.blocked_readonly') : testBlocked}
          />
        </>
      )}
    </div>
  );
};
