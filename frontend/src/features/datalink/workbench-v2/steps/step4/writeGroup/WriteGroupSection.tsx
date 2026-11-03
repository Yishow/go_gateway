import { RecordingRepairFocus } from '../../../recording-editor/RecordingRepairFocus';
import type { StudioRecordingRepair } from '../../../../../../services/studioRecordingView';
import { LocatorRecordingLifecycleHost } from './LocatorRecordingLifecycleHost';
import { recordingCopyPrefill, type RecordingCopyPrefill } from '../../../recording-editor/recordingCopyDraft';
import { studioV2WorkspaceWriteGroupsAPI } from '../../../../../../services/studioV2WorkspaceWriteGroups';
import { LocatorSavedRecordingStartReview } from '../../../recording-editor/LocatorSavedRecordingStartReview';
import { locatorSavedRecordingStartScope } from '../../../recording-editor/locatorSavedRecordingStartScope';
import { useRecordingNavigation } from '../../../recording-editor/RecordingDraftNavigation';
import type { DeviceSourcePort } from '../../../../../../pages/datalink/workbench-v2/deviceSourceSubmission';
import { readStudioLocatorEditorEntry } from '../../../../../../services/studioV2LocatorEntry';
import { studioV2LocatorHydration } from '../../../../../../pages/datalink/workbench-v2/studioV2LocatorHydration';
import type { StudioLocatorEntryRead } from '../../../../../../types/studioV2LocatorEntry';
import * as React from 'react';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useWriteGroupsQuery } from '../../../../../../hooks/datalink/useStudioV2WriteGroups';
import type { WriteGroup, WriteGroupMutationResponse } from '../../../../../../types/studioV2WriteGroup';
import type { WorkbenchV2State } from '../../../state/types';
import { buildGroupCandidates } from '../../../state/writeGroup/candidates';
import type { WorkspaceReadinessStepNumber } from '../../../components/WorkspaceReadinessPanel';
import { CutoverPanel } from './CutoverPanel';
import { GroupEditor } from './GroupEditor';
import { SavedRecordingCard } from '../SavedRecordingCard';
import type { SavedRecordingActionHost } from './lifecycle/SavedRecordingActions';
import { captureReturnFocus, type ReturnFocusTicket } from '../../../../../../utils/returnFocus';

export interface WriteGroupSectionProps {
  repairIntent?: StudioRecordingRepair;
  onRepairAdvanced?: () => void;
  presentation?: 'guided' | 'simple';
  initialSelection?: 'new' | 'located' | 'copy';
  destinationContent?: React.ReactNode;
  onClose?: () => void;
  onSavedLocation?: (group: WriteGroup) => void;
  onCopyLocation?: (group: WriteGroup) => void;
  deviceSource?: DeviceSourcePort;
  locatorEntry?: StudioLocatorEntryRead;
  locatorOnly?: boolean;
  locatorSuspended?: boolean;
  onEditorMountedChange?: (mounted: boolean) => void;
  state: WorkbenchV2State;
  workspaceId?: string;
  readonly: boolean;
  lifecycleHost?: SavedRecordingActionHost;
  onNavigateStep?: (step: WorkspaceReadinessStepNumber) => void;
  openGroupRequest?: { id: string };
  /** The host's reviewed scoped start: a managed profile group's effective action goes there instead of applying on its own. */
  onReviewScopedStart?: (group: WriteGroup) => void;
}

type SavedLocationReadback = { acknowledgement: WriteGroupMutationResponse; lease: object };
type SavedLocationReady = { group: WriteGroup; readback: SavedLocationReadback };

type Selection = { kind: 'none' } | { kind: 'new'; copyPrefill?: RecordingCopyPrefill } | { kind: 'edit'; id: string; saved?: WriteGroup };

type Prerequisite = 'device' | 'tags' | 'destination' | null;
function editorFocusScope(read: StudioLocatorEntryRead | undefined): string | undefined {
  const target = read?.selection.state === 'selected' ? read.selection.target : undefined;
  return target?.step === 4 && target.groupId && target.groupRevision
    ? JSON.stringify([target.workspaceId, target.workspaceRevision, target.groupId, target.groupRevision]) : undefined;
}

/**
 * The only database output editor of Step 4. It reads the saved groups of the
 * workspace and says what is missing before anything can be created, keeping a
 * failed load, an empty list and groups for another destination apart.
 */
export const WriteGroupSection: React.FC<WriteGroupSectionProps> = ({ presentation = 'guided', initialSelection, destinationContent, onClose, onSavedLocation, onCopyLocation, deviceSource, state, workspaceId, readonly, onNavigateStep, openGroupRequest, onReviewScopedStart, lifecycleHost, locatorEntry, locatorOnly, locatorSuspended, onEditorMountedChange, repairIntent, onRepairAdvanced }) => {
  const { t } = useTranslation('workbench-v2');
  const { t: editorText } = useTranslation('recording-editor');
  const { request: requestNavigation, blocked: navigationBlocked } = useRecordingNavigation();
  const { t: runtime } = useTranslation('runtime-dashboard');
  const region = useRef<HTMLElement>(null);
  const [selection, setSelection] = useState<Selection>({ kind: 'none' });
  const returnAfterClose = useRef(false);
  const copyAfterClose = useRef<{ group: WriteGroup; entry?: StudioLocatorEntryRead; lease: object } | undefined>(undefined);
  const savedLocationReadback = useRef<SavedLocationReadback | undefined>(undefined);
  const savedLocationAfterClose = useRef<SavedLocationReady | undefined>(undefined);
  const [savedLocationReady, setSavedLocationReady] = useState<SavedLocationReady>();
  const consumedSimple = useRef<string | StudioLocatorEntryRead | undefined>(undefined);
  const groupsQuery = useWriteGroupsQuery(Boolean(workspaceId) && !locatorOnly);
  const [editorEntry, setEditorEntry] = useState<StudioLocatorEntryRead>();
  const [startReview, setStartReview] = useState<{ entry: StudioLocatorEntryRead; lease: object; focus: ReturnFocusTicket }>();
  const [startActive, setStartActive] = useState(false);
  const startLease = useRef({ route: locatorEntry, suspended: locatorSuspended, readonly, workspaceId });
  if (startLease.current.route !== locatorEntry || startLease.current.suspended !== locatorSuspended || startLease.current.readonly !== readonly || startLease.current.workspaceId !== workspaceId)
    startLease.current = { route: locatorEntry, suspended: locatorSuspended, readonly, workspaceId };
  const renderStartLease = startLease.current;
  const canUseStartScope = () => Boolean(startReview && startReview.lease === startLease.current && startReview.entry === editorEntry &&
    !readonly && !locatorSuspended && !openingRequest.current && selection.kind === 'edit');
  const [openError, setOpenError] = useState<'unavailable' | 'stale-reference' | 'foreign-scope' | 'copy-unsupported'>();
  const [opening, setOpening] = useState(false);
  const openingRequest = useRef<AbortController | undefined>(undefined);
  const openedLocator = useRef<StudioLocatorEntryRead | undefined>(undefined);
  useEffect(() => {
    openingRequest.current?.abort(); openingRequest.current = undefined; setOpening(false);
    return () => { openingRequest.current?.abort(); openingRequest.current = undefined; consumedSimple.current = undefined; savedLocationReadback.current = undefined; savedLocationAfterClose.current = undefined; copyAfterClose.current = undefined; };
  }, [locatorEntry, locatorSuspended, workspaceId]);
  // Display-only receipt: never changes the URL expectation, editor lease or fresh Open read.
  const listReceipt = useRef<{ route: StudioLocatorEntryRead; read: StudioLocatorEntryRead } | undefined>(undefined);
  if (locatorSuspended || listReceipt.current?.route !== locatorEntry) listReceipt.current = undefined;
  const observed = locatorEntry?.entry.write_groups ?? groupsQuery.data;
  const displayRead = locatorOnly && !locatorSuspended ? listReceipt.current?.read : undefined;
  const editorState = useMemo(() => editorEntry ? studioV2LocatorHydration(editorEntry) : null, [editorEntry]);
  // Retain only a matching display read for this simple form. It cannot grant
  // writes: a failed or foreign current read below is denial-only.
  const retainedSimpleRead = useRef<{ workspaceId: string; data: NonNullable<typeof observed> } | undefined>(undefined);
  if (retainedSimpleRead.current?.workspaceId !== workspaceId) retainedSimpleRead.current = undefined;
  if (presentation === 'simple' && !locatorEntry && workspaceId && observed?.workspace_id === workspaceId &&
      typeof observed.workspace_revision === 'string' && !groupsQuery.isError && !groupsQuery.isFetching)
    retainedSimpleRead.current = { workspaceId, data: observed };
  const editorObserved = editorEntry?.entry.write_groups ?? (presentation === 'simple' && !locatorEntry ? retainedSimpleRead.current?.data : observed);
  const open = async (group: WriteGroup) => {
    if (!locatorEntry) { requestNavigation(() => setSelection({ kind: 'edit', id: group.id, saved: group })); return; }
    const target = locatorEntry.selection.state === 'selected' ? locatorEntry.selection.target : undefined;
    if (!target?.groupRevision || target.groupId !== group.id || target.workspaceId !== workspaceId ||
      locatorSuspended || selection.kind !== 'none' || openingRequest.current) return;
    const request = new AbortController(); openingRequest.current = request; setOpening(true); setOpenError(undefined);
    try {
      let expected = { workspaceRevision: target.workspaceRevision, groupRevision: target.groupRevision };
      if (openedLocator.current === locatorEntry) {
        // Reopening the same selection rechecks canonical revisions. A source
        // saved in the previous editor may have advanced the workspace revision.
        const latest = await groupsQuery.refetch();
        if (request.signal.aborted || openingRequest.current !== request) return;
        if (latest.isError || !latest.data) { setOpenError('unavailable'); return; }
        if (latest.data.workspace_id !== target.workspaceId) { setOpenError('foreign-scope'); return; }
        const exact = latest.data.groups.find(value => value.id === target.groupId && value.status !== 'deleted');
        if (!exact || !exact.revision.trim()) { setOpenError('stale-reference'); return; }
        expected = { workspaceRevision: latest.data.workspace_revision, groupRevision: exact.revision };
      }
      const result = await readStudioLocatorEditorEntry(locatorEntry.locator, expected, request.signal);
      if (request.signal.aborted || openingRequest.current !== request) return;
      if (!result.ok) { setOpenError(result.reason); return; }
      const next = result.value;
      const selected = next.selection.state === 'selected' ? next.selection.target : undefined;
      const saved = next.entry.write_groups?.groups.find(value => value.id === group.id && value.status !== 'deleted');
      const nextState = studioV2LocatorHydration(next);
      if (!next.entry.editor_hydrated || !selected || selected.workspaceId !== target.workspaceId ||
        selected.workspaceRevision !== expected.workspaceRevision || selected.groupId !== target.groupId ||
        selected.groupRevision !== expected.groupRevision || !saved || !nextState) {
        setOpenError('stale-reference'); return;
      }
      const copyPrefill = initialSelection === 'copy' ? recordingCopyPrefill(saved, next.entry.write_groups?.calculation_capabilities, nextState.db.connector.kind) : undefined;
      if (initialSelection === 'copy' && !copyPrefill) { setOpenError('copy-unsupported'); return; }
      openedLocator.current = locatorEntry;
      if (locatorEntry) listReceipt.current = { route: locatorEntry, read: next };
      setEditorEntry(next); setSelection(copyPrefill ? { kind: 'new', copyPrefill } : { kind: 'edit', id: group.id, saved });
    } catch { if (!request.signal.aborted) setOpenError('unavailable'); }
    finally { if (openingRequest.current === request) { openingRequest.current = undefined; setOpening(false); } }
  };
  const reloadEditor = async (acknowledgement?: WriteGroupMutationResponse, exactEntity = false) => {
    if (!editorEntry || locatorSuspended || selection.kind !== 'edit' || openingRequest.current) return;
    const original = editorEntry.selection.state === 'selected' ? editorEntry.selection.target : undefined;
    if (!original?.groupId) return;
    if (acknowledgement && (acknowledgement.group.workspace_id !== original.workspaceId || acknowledgement.group.id !== original.groupId)) return;
    const handoff = acknowledgement && (presentation === 'simple' && initialSelection === 'located' || locatorOnly && presentation === 'guided') && onSavedLocation && startLease.current === renderStartLease
      ? { acknowledgement, lease: renderStartLease } : undefined;
    const request = new AbortController(); openingRequest.current = request; setOpening(true); setOpenError(undefined);
    try {
      const latest = exactEntity ? await studioV2WorkspaceWriteGroupsAPI.get(original.groupId).then(result => ({
        isError: false, data: { workspace_id: result.group.workspace_id, workspace_revision: result.workspace_revision, groups: [result.group] },
      })) : await groupsQuery.refetch();
      if (request.signal.aborted || openingRequest.current !== request) return;
      if (latest.isError || !latest.data) { setOpenError('unavailable'); return; }
      if (latest.data.workspace_id !== original.workspaceId) { setOpenError('foreign-scope'); return; }
      const exact = latest.data.groups.find(group => group.id === original.groupId && group.status !== 'deleted');
      if (!exact) { setOpenError('stale-reference'); return; }
      const result = await readStudioLocatorEditorEntry(editorEntry.locator,
        { workspaceRevision: latest.data.workspace_revision, groupRevision: exact.revision }, request.signal);
      if (request.signal.aborted || openingRequest.current !== request) return;
      if (!result.ok) { setOpenError(result.reason); return; }
      const target = result.value.selection.state === 'selected' ? result.value.selection.target : undefined;
      const group = result.value.entry.write_groups?.groups.find(value => value.id === original.groupId && value.status !== 'deleted');
      if (!target || target.workspaceId !== original.workspaceId || target.groupId !== original.groupId ||
        target.workspaceRevision !== latest.data.workspace_revision || target.groupRevision !== exact.revision ||
        !result.value.entry.editor_hydrated || !group || !studioV2LocatorHydration(result.value)) {
        setOpenError('stale-reference'); return;
      }
      if (handoff && (handoff.lease !== startLease.current || handoff.acknowledgement.workspace_revision !== target.workspaceRevision ||
          handoff.acknowledgement.group.revision !== group.revision || group.workspace_id !== original.workspaceId || result.value.entry.group_heads[group.id] !== group.revision)) {
        setOpenError('stale-reference'); return;
      }
      if (acknowledgement) deviceSource?.acknowledgeLocatorGroupSave?.({ previous: editorEntry, acknowledgement, current: result.value });
      if (locatorEntry) listReceipt.current = { route: locatorEntry, read: result.value };
      setEditorEntry(result.value); setSelection({ kind: 'edit', id: group.id, saved: group });
      if (handoff) { savedLocationReadback.current = handoff; setSavedLocationReady({ group, readback: handoff }); }
    } catch { if (!request.signal.aborted) setOpenError('unavailable'); }
    finally { if (openingRequest.current === request) { openingRequest.current = undefined; setOpening(false); } }
  };
  // This is a read-only URL handoff after an accepted Save, not another Save.
  const readSavedLocation = async (readback: SavedLocationReadback) => {
    if (readback !== savedLocationReadback.current || readback.lease !== startLease.current || openingRequest.current || readonly || locatorSuspended) return;
    const accepted = readback.acknowledgement.group;
    const request = new AbortController(); openingRequest.current = request; setOpening(true); setOpenError(undefined);
    const current = () => !request.signal.aborted && openingRequest.current === request && savedLocationReadback.current === readback && startLease.current === readback.lease;
    try {
      const latest = initialSelection === 'copy' ? await studioV2WorkspaceWriteGroupsAPI.get(accepted.id).then(result => ({
        isError: false, data: { workspace_id: result.group.workspace_id, workspace_revision: result.workspace_revision, groups: [result.group] },
      })) : await groupsQuery.refetch();
      if (!current()) return;
      if (latest.isError || !latest.data) { setOpenError('unavailable'); return; }
      if (latest.data.workspace_id !== accepted.workspace_id) { setOpenError('foreign-scope'); return; }
      const exact = latest.data.groups.find(group => group.id === accepted.id && group.workspace_id === accepted.workspace_id && group.status !== 'deleted');
      if (!exact || exact.revision !== accepted.revision) { setOpenError('stale-reference'); return; }
      const result = await readStudioLocatorEditorEntry({ step: 4, workspace_id: accepted.workspace_id, device_id: null,
        rule_id: null, mapping_id: null, group_id: accepted.id, group_revision: accepted.revision },
      { workspaceRevision: latest.data.workspace_revision, groupRevision: accepted.revision }, request.signal);
      if (!current()) return;
      if (!result.ok) { setOpenError(result.reason); return; }
      const target = result.value.selection.state === 'selected' ? result.value.selection.target : undefined;
      const group = result.value.entry.write_groups?.groups.find(value => value.id === accepted.id && value.workspace_id === accepted.workspace_id && value.revision === accepted.revision && value.status !== 'deleted');
      if (!target || target.workspaceId !== accepted.workspace_id || target.groupId !== accepted.id || target.groupRevision !== accepted.revision ||
          target.workspaceRevision !== latest.data.workspace_revision || !group || !result.value.entry.editor_hydrated || !studioV2LocatorHydration(result.value)) {
        setOpenError('stale-reference'); return;
      }
      setSavedLocationReady({ group, readback });
    } catch { if (current()) setOpenError('unavailable'); }
    finally { if (openingRequest.current === request) { openingRequest.current = undefined; setOpening(false); } }
  };
  React.useLayoutEffect(() => {
    if (!savedLocationReady) return;
    setSavedLocationReady(undefined);
    if (savedLocationReady.readback !== savedLocationReadback.current || savedLocationReady.readback.lease !== startLease.current) return;
    // A late picker/source callback can still introduce unsent intent. Never
    // replace that editor merely because an older read has finished.
    if (navigationBlocked()) { setOpenError('stale-reference'); return; }
    savedLocationAfterClose.current = savedLocationReady;
    setSelection({ kind: 'none' }); setEditorEntry(undefined);
  }, [savedLocationReady, navigationBlocked]);
  const located = locatorEntry?.selection.state === 'selected' ? locatorEntry.selection.target : undefined;
  useEffect(() => { if (located) region.current?.querySelector<HTMLElement>('[data-locator-selected="true"]')?.focus(); }, [located]);
  // Navigation intent only. Selection remains this section's single authority,
  // and a saved record still enters through the original fresh scoped read.
  const simpleOpen = useRef(open); simpleOpen.current = open;
  useEffect(() => {
    if (returnAfterClose.current || copyAfterClose.current || savedLocationAfterClose.current || presentation !== 'simple' || readonly || locatorSuspended || selection.kind !== 'none') return;
    if (initialSelection === 'new' && !locatorEntry && observed && observed.workspace_id === workspaceId &&
        typeof observed.workspace_revision === 'string' && !groupsQuery.isError && !groupsQuery.isFetching) {
      const intent = `new:${workspaceId}`;
      if (consumedSimple.current !== intent) { consumedSimple.current = intent; setSelection({ kind: 'new' }); }
    }
    if ((initialSelection === 'located' || initialSelection === 'copy') && locatorEntry && consumedSimple.current !== locatorEntry && located?.groupId) {
      const group = observed?.groups.find(value => value.id === located.groupId && value.status !== 'deleted');
      if (group) { consumedSimple.current = locatorEntry; void simpleOpen.current(group); }
    }
  }, [presentation, initialSelection, readonly, locatorSuspended, locatorEntry, located, observed, workspaceId, selection.kind, groupsQuery.isError, groupsQuery.isFetching]);
  const consumedOpen = useRef<typeof openGroupRequest>(undefined);
  useEffect(() => {
    if (openGroupRequest !== consumedOpen.current && openGroupRequest && groupsQuery.data?.groups.some(group => group.id === openGroupRequest.id)) {
      consumedOpen.current = openGroupRequest;
      requestNavigation(() => setSelection({ kind: 'edit', id: openGroupRequest.id }));
    }
  }, [openGroupRequest, groupsQuery.data, requestNavigation]);
  const { candidates, excluded } = useMemo(() => buildGroupCandidates(editorState ?? state), [editorState, state]);
  const connector = (editorState ?? state).db.connector;
  const destinationSaved = connector.persisted === true && connector.save_state === 'saved' &&
    Boolean(connector.connector_id && connector.identity_revision);

  const prerequisite: Prerequisite = state.devices.length === 0 ? 'device'
    : candidates.length === 0 ? 'tags'
      : !destinationSaved ? 'destination' : null;

  const groups = (presentation === 'simple' && !locatorEntry ? editorObserved : observed)?.groups ?? [];
  const live = (displayRead?.entry.write_groups?.groups ?? groups).filter((group) => group.status !== 'deleted' && (!located || group.id === located.groupId));
  const removed = groups.filter((group) => group.status === 'deleted' && !locatorEntry);
  // Right after a save the refetched list may not carry the group yet; the returned group bridges that gap.
  const selected = selection.kind === 'edit' ? editorEntry ? selection.saved ?? editorObserved?.groups.find(group => group.id === selection.id)
    : groups.find(group => group.id === selection.id) ?? selection.saved : undefined;
  const copySource = initialSelection === 'copy' ? editorEntry?.entry.write_groups?.groups.find(group => group.id === located?.groupId) : undefined;
  const deniedCopyIdentity = (group: WriteGroup) => initialSelection === 'copy' && (!copySource || group.workspace_id !== workspaceId ||
    group.id === copySource.id || !group.record_id || group.record_id === copySource.record_id);
  const copyIdentityDenied = Boolean(selected && deniedCopyIdentity(selected));
  const management = { workspaceId, workspaceRevision: observed?.workspace_revision ?? '', lifecycleHost,
    readonly: readonly || !workspaceId || observed?.workspace_id !== workspaceId || !observed?.workspace_revision || !locatorOnly && groupsQuery.isFetching,
    onReload: () => { void groupsQuery.refetch(); } };
  const matches = (group: WriteGroup) => !destinationSaved ? null : group.destination.connector_id === connector.connector_id &&
    group.destination.connector_revision === connector.identity_revision;

  const editing = selection.kind === 'new' || Boolean(selection.kind === 'edit' && selected);
  useEffect(() => {
    // The original editor first resolves its leave guard and unmounts; only then
    // may the router leave, without prompting for the same discarded draft twice.
    if (!editing && returnAfterClose.current) { returnAfterClose.current = false; onClose?.(); }
  }, [editing, onClose]);
  const retention = useRef(onEditorMountedChange); retention.current = onEditorMountedChange;
  useEffect(() => { retention.current?.(editing); return () => retention.current?.(false); }, [editing]);
  useEffect(() => {
    const pending = copyAfterClose.current;
    if (editing || !pending) return;
    copyAfterClose.current = undefined;
    if (pending.lease !== startLease.current) return;
    if (navigationBlocked()) { setEditorEntry(pending.entry); setSelection({ kind: 'edit', id: pending.group.id, saved: pending.group }); return; }
    onCopyLocation?.(pending.group);
  }, [editing, navigationBlocked, onCopyLocation]);
  useEffect(() => {
    const ready = savedLocationAfterClose.current;
    if (editing || !ready) return;
    savedLocationAfterClose.current = undefined;
    if (ready.readback !== savedLocationReadback.current || ready.readback.lease !== startLease.current) return;
    if (navigationBlocked()) { setSelection({ kind: 'edit', id: ready.group.id, saved: ready.group }); setOpenError('stale-reference'); return; }
    savedLocationReadback.current = undefined;
    onSavedLocation?.(ready.group);
  }, [editing, navigationBlocked, onSavedLocation]);
  const invalidSimpleScope = presentation === 'simple' && observed &&
    (observed.workspace_id !== workspaceId || typeof observed.workspace_revision !== 'string');
  if ((!editing || !editorObserved) && !locatorOnly && (groupsQuery.isError || invalidSimpleScope)) {
    return (
      <section role="alert" className="space-y-2 rounded-xl border border-[var(--gw-warning-border)] p-5 text-xs text-[var(--gw-warning-fg)]" data-testid="group-section-error">
        <p>{t('step4.group.section.load_failed')}</p>
        <button type="button" onClick={() => void groupsQuery.refetch()} className="underline" data-testid="group-section-retry">{t('step4.group.retry')}</button>
      </section>
    );
  }
  if (!editing && !locatorOnly && groupsQuery.isLoading) {
    return <section className="p-5 text-xs text-[var(--gw-muted)]" role="status" data-testid="group-section-loading">{t('step4.group.section.loading')}</section>;
  }

  const renderEditor = (editorLifecycleHost = lifecycleHost, lifecycleReadonly = false, isOriginalScope?: () => boolean) => editing && editorObserved ? (
        <GroupEditor scopedPresentation={Boolean(locatorOnly && presentation === 'guided')} presentation={presentation} scopeDiagnostics={collapseDiagnostics ? locatorDiagnostics : undefined} destinationContent={destinationContent} calculationCapabilities={editorObserved.calculation_capabilities} deviceSource={deviceSource} sourceDevices={(editorState ?? state).devices.filter(device => device.persisted)}
          copyPrefill={selection.kind === 'new' ? selection.copyPrefill : undefined}
          key={`${editorObserved.workspace_id}:${selection.kind === 'edit' ? selection.id : 'new'}`}
          group={selected} workspaceId={editorObserved.workspace_id} workspaceRevision={editorObserved.workspace_revision}
          locatorGuard={Boolean(locatorEntry || locatorSuspended)}
          saveReadbackPending={Boolean(opening && (editorEntry || savedLocationReadback.current)) || Boolean(isOriginalScope && !isOriginalScope())}
          locatorScope={editorEntry ? { workspaceId: editorEntry.workspace.id, workspaceRevision: editorEntry.entry.workspace_revision } : undefined}
          connector={connector} candidates={candidates} excluded={excluded} readonly={copyIdentityDenied || lifecycleReadonly || readonly || Boolean(locatorSuspended) || startActive || Boolean(opening && savedLocationReadback.current) || presentation === 'simple' && !locatorEntry &&
            (groupsQuery.isError || !observed || Boolean(invalidSimpleScope))} editorReadonly={selected?.status === 'deleted'}
          onClose={() => { openingRequest.current?.abort(); openingRequest.current = undefined; setOpening(false); savedLocationReadback.current = undefined; savedLocationAfterClose.current = undefined; returnAfterClose.current = true; setSelection({ kind: 'none' }); setEditorEntry(undefined); }}
          onCopy={presentation === 'simple' && initialSelection !== 'copy' && onCopyLocation && selected && recordingCopyPrefill(selected, editorObserved.calculation_capabilities, connector.kind) ? () => {
            if (startLease.current !== renderStartLease || readonly || locatorSuspended || openingRequest.current || copyAfterClose.current) return;
            copyAfterClose.current = { group: selected, entry: editorEntry, lease: renderStartLease };
            setSelection({ kind: 'none' }); setEditorEntry(undefined);
          } : undefined}
          onSaved={(group, savedRevision, ownSaveAcknowledgement) => {
            setSelection({ kind: 'edit', id: group.id, saved: group });
            if (editorEntry && initialSelection !== 'copy') { void reloadEditor(ownSaveAcknowledgement); return; }
            if (presentation !== 'simple' || initialSelection !== 'new' && initialSelection !== 'copy' || initialSelection !== 'copy' && locatorEntry || !onSavedLocation || startLease.current !== renderStartLease) return;
            openingRequest.current?.abort(); openingRequest.current = undefined; savedLocationReadback.current = undefined;
            if (deniedCopyIdentity(group)) { setOpenError('stale-reference'); return; }
            if (!ownSaveAcknowledgement || group.workspace_id !== workspaceId || ownSaveAcknowledgement.workspace_revision !== savedRevision ||
                ownSaveAcknowledgement.group.id !== group.id || ownSaveAcknowledgement.group.workspace_id !== workspaceId || ownSaveAcknowledgement.group.revision !== group.revision) return;
            const readback = { acknowledgement: ownSaveAcknowledgement, lease: renderStartLease };
            savedLocationReadback.current = readback; void readSavedLocation(readback);
          }}
          onReload={() => { if (savedLocationReadback.current) void readSavedLocation(savedLocationReadback.current); else if (editorEntry && initialSelection !== 'copy') void reloadEditor(); else if (initialSelection !== 'copy') void groupsQuery.refetch(); }}
          fetchLatest={async () => {
            // Ref fence also covers a cached Save callback before React renders
            // the pending flag. It never grants scope or changes picker state.
            if (openingRequest.current || copyIdentityDenied || isOriginalScope && !isOriginalScope()) throw new Error('Recording scope is unavailable.');
            if (initialSelection === 'copy' && editorEntry?.selection.state === 'selected') {
              const id = selected?.id ?? editorEntry.selection.target.groupId;
              if (!id) throw new Error('Recording scope is unavailable.');
              const latest = await studioV2WorkspaceWriteGroupsAPI.get(id);
              if (latest.group.id !== id || latest.group.workspace_id !== workspaceId || latest.group.status === 'deleted' ||
                  !selected && latest.group.revision !== editorEntry.selection.target.groupRevision) throw new Error('Recording scope is unavailable.');
              return { workspace_revision: latest.workspace_revision, groups: [latest.group] };
            }
            const latest = await groupsQuery.refetch();
            if (latest.isError || !latest.data || latest.data.workspace_id !== editorObserved.workspace_id ||
              selected && !latest.data.groups.some(group => group.id === selected.id && group.status !== 'deleted'))
              throw new Error('Recording scope is unavailable.');
            return latest.data;
          }}
          configuredDestinations={{ data: editorObserved, loading: groupsQuery.isFetching || opening, failed: groupsQuery.isError || Boolean(invalidSimpleScope),
            retry: () => { if (editorEntry) void reloadEditor(); else if (savedLocationReadback.current) void readSavedLocation(savedLocationReadback.current); else void groupsQuery.refetch(); }, open: locatorEntry ? undefined : group => { void open(group); } }}
          onReviewScopedStart={initialSelection === 'copy' ? undefined : locatorOnly ? group => {
            if (!editorEntry || locatorSuspended || readonly || openingRequest.current || startActive || openedLocator.current !== locatorEntry || isOriginalScope && !isOriginalScope()) return;
            const context = locatorSavedRecordingStartScope(editorEntry, workspaceId ?? '');
            if (context?.group.id !== group.id || context.group.revision !== group.revision) return;
            setStartReview({ entry: editorEntry, lease: startLease.current,
              focus: captureReturnFocus(startLease.current, document.activeElement instanceof HTMLElement ? document.activeElement : null) });
          } : onReviewScopedStart} lifecycleHost={editorLifecycleHost}
          repairContent={presentation === 'simple' && repairIntent && editorEntry && initialSelection === 'located' && <RecordingRepairFocus
          repair={repairIntent} entry={editorEntry} region={region} onAdvanced={onRepairAdvanced}
          canUseScope={() => startLease.current === renderStartLease && openedLocator.current === locatorEntry && !readonly && !lifecycleReadonly &&
            !locatorSuspended && !openingRequest.current && selection.kind === 'edit' && (!isOriginalScope || isOriginalScope()) &&
            (repairIntent.operation?.owner !== 'start' || Boolean(locatorSavedRecordingStartScope(editorEntry, workspaceId ?? '')))}
          onReviewStart={() => {
            if (startLease.current !== renderStartLease || readonly || lifecycleReadonly || locatorSuspended || openingRequest.current || startActive ||
              openedLocator.current !== locatorEntry || isOriginalScope && !isOriginalScope() || !locatorSavedRecordingStartScope(editorEntry, workspaceId ?? '')) return;
            setStartReview({ entry: editorEntry, lease: startLease.current,
              focus: captureReturnFocus(startLease.current, document.activeElement instanceof HTMLElement ? document.activeElement : null) });
          }} />}
        />
  ) : null;

  const draftBlocked = initialSelection !== 'copy' && located && selection.kind !== 'none' && (selection.kind !== 'edit' || selection.id !== located.groupId);
  const collapseDiagnostics = (presentation === 'simple' || locatorOnly && presentation === 'guided') && !readonly && !locatorSuspended && !opening && !openError && !draftBlocked && !invalidSimpleScope && !groupsQuery.isError &&
    selected && editorEntry?.selection.state === 'selected' && selected.id === located?.groupId && selected.workspace_id === workspaceId &&
    editorEntry.selection.target.groupId === selected.id && editorEntry.entry.group_heads[selected.id] === selected.revision;
  const collapseListDiagnostics = locatorOnly && presentation === 'guided' && !editing && !opening && !openError && !locatorSuspended && !readonly && !draftBlocked && !groupsQuery.isError && !invalidSimpleScope &&
    located?.step === 4 && located.workspaceId === workspaceId && observed?.workspace_revision === located.workspaceRevision &&
    observed.groups.some(group => group.workspace_id === workspaceId && group.id === located.groupId && group.revision === located.groupRevision && locatorEntry?.entry.group_heads[group.id] === group.revision);
  const locatorDiagnostics = <>
    {locatorOnly && <p role="status">{runtime('locatorEntry.unrequested')}</p>}
    {draftBlocked && <p role="alert">{runtime('locatorEntry.draftBlocked')}</p>}
    {locatorEntry?.selection.state === 'selected' && locatorEntry.selection.history && <p>{runtime('locatorEntry.historicalHint', { revision: locatorEntry.selection.history.revisionHint })}</p>}
    {located?.groupId && located.groupRevision && located.workspaceId === workspaceId &&
      <p>{runtime('locatorEntry.currentRevision')}: <span>{displayRead?.entry.group_heads[located.groupId] ?? located.groupRevision}</span></p>}
  </>;

  // A same-canonical completion readback may replace a DOM button, but never renews the old command lease.
  const focusReturn = startReview ? { ticket: startReview.focus, currentTarget: () => {
    const node = region.current, expected = editorFocusScope(startReview.entry);
    if (startReview.lease !== startLease.current || openingRequest.current || !node?.isConnected || !expected || node.dataset.focusScope !== expected) return null;
    const trigger = startReview.focus.trigger;
    if (trigger && node.contains(trigger)) return trigger;
    return trigger?.dataset.testid === 'recording-repair-start' ? node.querySelector<HTMLElement>('[data-testid="recording-repair-start"]') : null;
  } } : undefined;
  return (
    <section ref={region} className="space-y-4 rounded-xl border border-[var(--gw-border)] bg-[var(--gw-surface)] p-5" aria-labelledby={presentation === 'simple' ? undefined : 'group-section-title'} data-testid="group-section" data-editor-presentation={presentation} data-editor-active={Boolean(editing)}
      data-focus-scope={!opening && !openError && !locatorSuspended && selection.kind === 'edit' ? editorFocusScope(editorEntry) : undefined}>
      {presentation !== 'simple' && <div>
        <h3 id="group-section-title" className="text-sm font-semibold text-[var(--gw-text)]">{t('step4.group.section.title')}</h3>
        <p className="mt-1 text-xs text-[var(--gw-muted)]">{t('step4.group.section.subtitle')}</p>
      </div>}

      {editing && (groupsQuery.isError || invalidSimpleScope) && !locatorOnly && <p role="alert">{t('step4.group.section.load_failed')}</p>}
      {opening && <p role="status">{runtime('locatorEntry.loading')}</p>}
      {openError && <><p role="alert" data-testid="studio-locator-open-error">{openError === 'copy-unsupported' ? editorText('copy.unsupported') : runtime(`locatorEntry.${openError === 'stale-reference' ? 'stale' : openError === 'foreign-scope' ? 'foreign' : 'unavailable'}`)}</p>
        {savedLocationReadback.current && <button type="button" disabled={opening || readonly} data-testid="saved-location-retry"
          onClick={() => { if (savedLocationReadback.current) void readSavedLocation(savedLocationReadback.current); }}>{t('step4.group.retry')}</button>}</>}
      {!collapseDiagnostics && (collapseListDiagnostics ? <details data-testid="group-list-diagnostics" className="text-xs text-[var(--gw-muted)]"><summary>{editorText('versions.details')}</summary>{locatorDiagnostics}</details> : locatorDiagnostics)}
      {presentation !== 'simple' && !locatorOnly && prerequisite && (
        <div className="space-y-1 text-xs text-[var(--gw-warning-fg)]" role="status" data-testid={`group-prerequisite-${prerequisite}`}>
          <p>{t(`step4.group.prerequisite.${prerequisite}`)}</p>
          {prerequisite === 'device' && onNavigateStep && (
            <button type="button" onClick={() => onNavigateStep(1)} className="underline" data-testid="group-prerequisite-go">{t('step4.group.prerequisite.go_step1')}</button>
          )}
          {prerequisite === 'tags' && (
            <>
              {excluded.length > 0 && <p data-testid="group-prerequisite-excluded">{t('step4.group.prerequisite.excluded', { count: excluded.length })}</p>}
              {onNavigateStep && <button type="button" onClick={() => onNavigateStep(3)} className="underline" data-testid="group-prerequisite-go">{t('step4.group.prerequisite.go_step3')}</button>}
            </>
          )}
        </div>
      )}

      {!editing && (
        <>
          {live.length === 0 && removed.length === 0 ? (
            <p className="text-xs text-[var(--gw-muted)]" data-testid="group-empty">{t('step4.group.list.empty')}</p>
          ) : (
            <ul className="space-y-2" data-testid="group-list">
              {live.map((group) => <SavedRecordingCard key={group.id} group={group} highlighted={located?.groupId === group.id} current={matches(group)} onOpen={() => void open(group)}
                management={{ ...management, readonly: management.readonly || group.workspace_id !== workspaceId }} />)}
            </ul>
          )}
          {live.length > 0 && !live.some(matches) && destinationSaved && (
            <p className="text-xs text-[var(--gw-warning-fg)]" role="status" data-testid="group-none-match">{t('step4.group.list.none_match')}</p>
          )}
          {!locatorOnly && groupsQuery.data && observed && live.length > 0 && (
            <CutoverPanel groups={live} workspaceId={observed.workspace_id} workspaceRevision={observed.workspace_revision}
              connectorRevision={connector.identity_revision ?? ''} readonly={readonly || !destinationSaved} onOpenGroup={(id) => setSelection({ kind: 'edit', id })}
              candidates={candidates} />
          )}
          <button
            type="button" disabled={readonly || prerequisite !== null || !observed || Boolean(locatorOnly)}
            onClick={() => setSelection({ kind: 'new' })}
            className="rounded bg-[var(--gw-primary)] px-4 py-1.5 text-xs font-semibold text-[var(--gw-on-primary)] disabled:text-[var(--gw-muted)] disabled:bg-[var(--gw-ground)]" data-testid="group-create"
          >
            {t('step4.group.list.create')}
          </button>
          {removed.length > 0 && (
            <details data-testid="group-removed">
              <summary className="cursor-pointer text-xs text-[var(--gw-muted)]">{t('step4.group.list.removed', { count: removed.length })}</summary>
              <ul className="mt-2 space-y-2">
                {removed.map((group) => <SavedRecordingCard key={group.id} group={group} highlighted={located?.groupId === group.id} current={matches(group)} onOpen={() => void open(group)}
                management={{ ...management, readonly: management.readonly || group.workspace_id !== workspaceId }} />)}
              </ul>
            </details>
          )}
        </>
      )}

      {locatorOnly && editorEntry && initialSelection !== 'copy' ? <LocatorRecordingLifecycleHost entry={editorEntry} scopeLease={renderStartLease}
        readonly={readonly || Boolean(locatorSuspended) || startActive}
        canUseScope={() => startLease.current === renderStartLease && openedLocator.current === locatorEntry && !locatorSuspended && !openingRequest.current && selection.kind === 'edit'}
        onReload={() => { void reloadEditor(undefined, true); }}>{renderEditor}</LocatorRecordingLifecycleHost> : renderEditor()}

      {startReview && workspaceId && <LocatorSavedRecordingStartReview entry={editorEntry ?? startReview.entry} workspaceId={workspaceId}
        focusReturn={focusReturn}
        readonly={!canUseStartScope()} canUseScope={canUseStartScope} deviceSource={deviceSource}
        onClose={() => setStartReview(undefined)} onCompleted={() => { if (canUseStartScope()) void reloadEditor(); }} onOperationActiveChange={setStartActive} />}
    </section>
  );
};
