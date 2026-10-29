import * as Dialog from '@radix-ui/react-dialog';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { recordingNextAction, recordingWorkflowView, type RecordingNextActionInput, type RecordingNextIntent, type RecordingWorkflowScope } from '../../state/recordingWorkflowView';
import type { RecordingOwnerObservation } from '../../state/recordingOwnerObservation';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import { safeRuntimeLabel } from '../../../../../utils/runtimeOperatorSafety';
import { inspectSchemaOperation } from './writeGroup/schemaOperationRecovery';
import type { SchemaPreparationPresentation } from './writeGroup/schemaPreparationPresentation';
import { RecordingPreparationEvidence, type RecordingPreparationEvidenceView } from './RecordingPreparationEvidence';
import { RecordingNextAction } from './RecordingNextAction';
import { useRecordingNextActionLabels } from './recordingNextActionLabels';
import type { GuidedRecordingPreparationProps } from './recordingPreparationPresentation';
import { recordingStartObservation, useRecordingPreparationOwner } from './useRecordingPreparationOwner';
import { inspectBasicRecordingIntent } from '../../state/basicRecording';
import { RecordingPreparationTestWriteQuery, type RecordingPreparationTestWriteQueryReport } from './RecordingPreparationTestWriteQuery';
import type { RecordingTestWriteReference } from '../../state/recordingTestWriteRecovery';
import { recordingPreparationBlocked, recordingSchemaPreviewRetryBlocked, type RecordingPreparationFenceInput } from './recordingPreparationFence';

const rank = { query_operation: 0, recover_operation: 1, choose_sources: 2, resolve_semantics: 3, inspect: 4,
  save_prepare: 5, refresh_evidence: 6, review_schema: 7, confirm_schema: 8, preview_schema: 9, review_scoped_start: 10, monitor: 11 };
function EvidenceOwner({ group, scope, report }: { group: WriteGroup; scope: RecordingWorkflowScope;
  report: (id: string, evidence: RecordingPreparationEvidenceView | null) => void }) {
  const onEvidence = useCallback((value: RecordingPreparationEvidenceView | null) => report(group.id, value), [group.id, report]);
  return <RecordingPreparationEvidence group={group} scope={scope} onEvidence={onEvidence} />;
}
function TestWriteOwner({ reference, identity, report }: { reference: RecordingTestWriteReference; identity: string;
  report: (binding: string, query: RecordingPreparationTestWriteQueryReport | null) => void }) {
  const binding = JSON.stringify([identity, reference]);
  const onQuery = useCallback((value: RecordingPreparationTestWriteQueryReport | null) => report(binding, value), [binding, report]);
  return <RecordingPreparationTestWriteQuery reference={reference} ownerIdentity={binding} onQuery={onQuery} />;
}

/** Presentation only: private owner commands remain the only save/schema/start implementation. */
export function GuidedRecordingPreparation(props: GuidedRecordingPreparationProps) {
  const { t } = useTranslation('workbench-v2-recording-guided');
  const labels = useRecordingNextActionLabels();
  const { scope, groups, drafts, start } = props;
  const identity = JSON.stringify([props.ownerKey, props.consentKey, scope.workspaceRevision, groups.map(g => [g.id, g.revision,
    g.applied_revision, g.destination.connector_id, g.destination.connector_revision])]);
  const owner = useRecordingPreparationOwner(identity, scope.workspaceId, groups.map(group => group.id));
  const checkedReferences = useRef({ epoch: owner.epoch, references: new Map<string, RecordingTestWriteReference>() });
  if (checkedReferences.current.epoch !== owner.epoch) checkedReferences.current = { epoch: owner.epoch, references: new Map() };
  owner.references.forEach(reference => checkedReferences.current.references.set(reference.operation_id, reference));
  const queryIdentity = JSON.stringify([identity, owner.epoch]);
  const startObservation = recordingStartObservation(props.ownerKey, {
    unresolved: start.observation.state === 'unresolved',
    operationId: 'operationId' in start.observation ? start.observation.operationId : undefined,
    status: start.observation.state === 'observed' && start.observation.result === 'resolved' ? start.observation.status : undefined,
    busy: start.busy, queryFailed: start.observation.state === 'unresolved' && start.observation.reason === 'query_failed',
  });
  const [schemaReports, setSchemaReports] = useState<Record<string, SchemaPreparationPresentation>>({});
  const [evidence, setEvidence] = useState<Record<string, RecordingPreparationEvidenceView>>({});
  const fenceInput = useRef<RecordingPreparationFenceInput>({ ownerKey: props.ownerKey, workspaceId: scope.workspaceId,
    workspaceRevision: scope.workspaceRevision, groups, start: startObservation, schema: schemaReports });
  fenceInput.current = { ownerKey: props.ownerKey, workspaceId: scope.workspaceId,
    workspaceRevision: scope.workspaceRevision, groups, start: startObservation, schema: schemaReports };
  const [reviewEpoch, setReviewEpoch] = useState<number>();
  const reviewFocus = useRef<{ epoch: number; testId: string; dismissed: boolean } | undefined>(undefined);
  const [testWriteQueries, setTestWriteQueries] = useState<Record<string, RecordingPreparationTestWriteQueryReport>>({});
  const [busy, setBusy] = useState(false);
  const inFlight = useRef(false);
  const mounted = useRef(true);
  const region = useRef<HTMLElement>(null);
  useEffect(() => { mounted.current = true; return () => { mounted.current = false; }; }, []);
  const reportTestWriteQuery = useCallback((binding: string, query: RecordingPreparationTestWriteQueryReport | null) => {
    setTestWriteQueries(current => { if (current[binding] === query || !query && !current[binding]) return current;
      const next = { ...current }; if (query) next[binding] = query; else delete next[binding]; return next; });
  }, []);
  const externalFence = useRef(props.isNewCommandBlocked); externalFence.current = props.isNewCommandBlocked;
  const currentCommandBlocked = () => { try { return Boolean(externalFence.current?.()); } catch { return true; } };
  const blocked = currentCommandBlocked() || recordingPreparationBlocked(fenceInput.current);
  const blockedGroupIds = new Set(groups.filter(group => currentCommandBlocked() || recordingPreparationBlocked(fenceInput.current, group.id)).map(group => group.id));
  const otherOwnerBlockedGroupIds = new Set(groups.filter(group => recordingPreparationBlocked(fenceInput.current, undefined, false) ||
    owner.references.some(reference => reference.group_id !== group.id)).map(group => group.id));
  const generation = owner.epoch;
  // A batch may have just accepted a save while React is still publishing the new schema-owner reports.
  // This describes that handoff only; the batch must still re-run the unchanged exact-head denial below.
  const priorSchema = useRef<Record<string, SchemaPreparationPresentation>>({});
  const [publishedIdentity, setPublishedIdentity] = useState(identity);
  // Child owner effects have one publication pass; missing proof then fails closed instead of waiting forever.
  useEffect(() => { setPublishedIdentity(identity); }, [identity]);
  const publication = useRef(false); publication.current = publishedIdentity === identity;
  const handoffGuard = useRef({ blocked: props.batchHandoff?.blocked, recoveryAvailable: owner.recoveryAvailable });
  handoffGuard.current = { blocked: props.batchHandoff?.blocked, recoveryAvailable: owner.recoveryAvailable };
  const readHandoffState = (): 'settling' | 'current' | 'blocked' => {
    const input = fenceInput.current;
    const held = inspectBasicRecordingIntent(input.ownerKey);
    const resolved = input.start.state === 'observed' && input.start.result === 'resolved' ? input.start.operationId : undefined;
    if (!mounted.current || !navigator.onLine || handoffGuard.current.blocked || input.start.state !== 'observed' ||
      !handoffGuard.current.recoveryAvailable() || held.state === 'unavailable' ||
      held.state === 'retained' && (!resolved || held.operationId !== resolved)) return 'blocked';
    let settling = false;
    for (const group of input.groups) {
      const report = input.schema[group.id], previous = report ?? priorSchema.current[group.id];
      if (inspectSchemaOperation(input.workspaceId, group.id).state !== 'empty' || !previous ||
        previous.scope.workspaceId !== input.workspaceId || previous.scope.groupId !== group.id ||
        previous.scope.connectorId !== group.destination.connector_id || previous.scope.connectorRevision !== group.destination.connector_revision ||
        previous.observation.state !== 'observed') return 'blocked';
      if (!report || report.scope.workspaceRevision !== input.workspaceRevision || report.scope.groupRevision !== group.revision) settling = true;
    }
    return settling ? publication.current ? 'blocked' : 'settling' : 'current';
  };
  const handoffState = readHandoffState();
  Object.entries(schemaReports).forEach(([id, report]) => { priorSchema.current[id] = report; });
  const handoff = props.batchHandoff ? { workspaceRevision: scope.workspaceRevision, consentKey: props.batchHandoff.consentKey, state: readHandoffState } : undefined;
  const commandFence = useMemo(() => ({ ownerKey: props.ownerKey, scopeKey: props.fenceKey ?? props.ownerKey, blocked, handoff,
    isBlocked: (groupId?: string) => !mounted.current || !owner.isCurrent(generation) || currentCommandBlocked() || recordingPreparationBlocked(fenceInput.current, groupId) }),
    // The guarded callback reads current owner observations, never render-time permission.
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [props.ownerKey, props.fenceKey, blocked, generation, handoffState, props.batchHandoff?.consentKey, scope.workspaceRevision]);
  const notifyFence = props.onCommandFence;
  useEffect(() => { notifyFence?.(commandFence); }, [notifyFence, commandFence]);
  useEffect(() => () => notifyFence?.(null), [notifyFence]);
  const reportSchema = useCallback((id: string, value: SchemaPreparationPresentation | null) => setSchemaReports(current => {
    if (current[id] === value || (!value && !current[id])) return current;
    const next = { ...current }; if (value) next[id] = value; else delete next[id]; return next;
  }), []);
  const reportEvidence = useCallback((id: string, value: RecordingPreparationEvidenceView | null) => setEvidence(current => {
    if (current[id] === value || (!value && !current[id])) return current;
    const next = { ...current }; if (value) next[id] = value; else delete next[id]; return next;
  }), []);
  const groupScope = (group: WriteGroup): RecordingWorkflowScope => ({ workspaceId: scope.workspaceId,
    workspaceRevision: scope.workspaceRevision, deviceId: scope.deviceId, ...(scope.circuitId ? { circuitId: scope.circuitId } : {}),
    groupId: group.id, groupRevision: group.revision, appliedRevision: group.applied_revision });
  const currentSchema = (group: WriteGroup) => {
    const report = schemaReports[group.id];
    return report?.scope.workspaceId === scope.workspaceId && report.scope.workspaceRevision === scope.workspaceRevision &&
      report.scope.groupRevision === group.revision && report.scope.connectorId === group.destination.connector_id &&
      report.scope.connectorRevision === group.destination.connector_revision ? report : undefined;
  };
  const observations = (group?: WriteGroup): RecordingOwnerObservation[] => {
    const selected = group ? groupScope(group) : { workspaceId: scope.workspaceId, workspaceRevision: scope.workspaceRevision,
      deviceId: scope.deviceId, ...(scope.circuitId ? { circuitId: scope.circuitId } : {}) };
    const report = group ? currentSchema(group) : undefined;
    const stored = group ? inspectSchemaOperation(scope.workspaceId, group.id) : undefined;
    const schema = report?.observation ?? (stored?.state === 'retained'
      ? { state: 'unresolved' as const, reason: 'unknown' as const, operationId: stored.operationId }
      : group ? { state: 'unavailable' as const, reason: 'loading' as const }
        // There is no saved schema owner for this unsaved draft; this is not a server ledger query.
        : { state: 'observed' as const, result: 'none' as const });
    const reference = owner.references.find(value => value.group_id === group?.id);
    const testWrite = reference ? { state: 'unresolved' as const, reason: 'cleanup_pending' as const, operationId: reference.operation_id }
      : owner.recoveryUnavailable ? { state: 'unavailable' as const, reason: 'unreadable' as const } : { state: 'observed' as const, result: 'none' as const };
    return [{ ...schema, owner: 'schema', scope: selected }, { ...startObservation, owner: 'start', scope: selected },
      { ...testWrite, owner: 'test_write', scope: selected }];
  };
  const inputs: RecordingNextActionInput[] = [];
  for (const group of groups) {
    const selected = groupScope(group);
    const reports = observations(group);
    const proof = evidence[group.id];
    const currentProof = proof && JSON.stringify(proof.scope) === JSON.stringify(selected) ? proof : undefined;
    const nativeDraft = drafts.find(draft => draft.recording_profile?.kind === group.recording_profile?.kind);
    const schema = currentSchema(group);
    inputs.push({ kind: 'saved', scope: selected,
      workflow: recordingWorkflowView({ scope: selected, savedGroup: group, ownerObservations: reports,
        readiness: currentProof?.readiness, lifecycle: currentProof?.lifecycle }),
      ...(props.dirtyGroupIds.has(group.id) && nativeDraft ? { draft: nativeDraft } : {}),
      ownerObservations: reports, counterStatuses: props.counterStatuses,
      schemaReview: { scope: selected, stage: schema?.preview === 'ready' ? 'reviewed' : 'idle',
        ...(schema?.preview === 'ready' && schema.previewOperationId ? { operationId: schema.previewOperationId } : {}) },
      connection: owner.online ? 'online' : 'offline', nowMs: Date.now() });
  }
  for (const draft of drafts.filter(draft => !groups.some(group => group.recording_profile?.kind === draft.recording_profile?.kind))) {
    inputs.push({ kind: 'draft', scope: { workspaceId: scope.workspaceId, workspaceRevision: scope.workspaceRevision,
      deviceId: scope.deviceId, ...(scope.circuitId ? { circuitId: scope.circuitId } : {}) }, draft,
      ownerObservations: observations(), counterStatuses: props.counterStatuses, connection: owner.online ? 'online' : 'offline', nowMs: Date.now() });
  }
  const candidates = inputs.map(input => ({ input, model: recordingNextAction(input) }));
  candidates.sort((a, b) => rank[a.model.intent.kind] - rank[b.model.intent.kind] || Number(a.model.enabled) - Number(b.model.enabled));
  const selected = candidates[0];
  const allReady = candidates.length > 0 && candidates.every(({ model }) => model.intent.kind === 'review_scoped_start' && model.enabled);
  const selectedGroupId = selected?.input.kind === 'saved' ? selected.input.scope.groupId : undefined;
  const selectedGroup = groups.find(group => group.id === selectedGroupId);
  const schema = selectedGroup ? currentSchema(selectedGroup) : undefined;
  const open = reviewEpoch !== undefined && reviewEpoch === owner.epoch;
  const locked = props.readonly || busy || start.busy;
  const focus = (target: string) => {
    const control = region.current?.closest('[data-testid$="basic-panel"], [data-testid="basic-recording-panel"]')?.querySelector<HTMLElement>(target)
      ?? document.querySelector<HTMLElement>(target);
    control?.scrollIntoView?.({ block: 'center' });
    if (control && !control.matches('button,input,select,textarea,a,[tabindex]')) control.tabIndex = -1;
    control?.focus();
  };
  const inspect = () => { region.current?.focus(); props.onInspect?.(); };
  const repair = () => { if (props.onRepair) props.onRepair(); else focus('[data-testid="circuit-profile-interval-instantaneous"], [data-testid="basic-recording-interval"]'); };
  const run = async (command: (() => Promise<void> | void) | undefined, epoch: number, mutating = true) => {
    if (!command || !mounted.current || !owner.isCurrent(epoch) || !navigator.onLine || mutating && (current.current.props.readonly || !current.current.props.canPrepare) || inFlight.current ||
      mutating && !owner.recoveryAvailable()) return;
    if (mutating) {
      try { if (current.current.props.isNewCommandBlocked?.()) return; } catch { return; }
      const heldStart = inspectBasicRecordingIntent(props.ownerKey);
      const resolvedId = start.observation.state === 'observed' && start.observation.result === 'resolved' ? start.observation.operationId : undefined;
      if (heldStart.state === 'unavailable' || heldStart.state === 'retained' && (!resolvedId || heldStart.operationId !== resolvedId) ||
        groups.some(group => inspectSchemaOperation(scope.workspaceId, group.id).state !== 'empty')) return;
    }
    inFlight.current = true; setBusy(true);
    try { await command(); } finally { inFlight.current = false; if (mounted.current) setBusy(false); }
  };
  const current = useRef({ selected, schema, allReady, locked, owner, props });
  current.current = { selected, schema, allReady, locked, owner, props };
  const isSchemaPreviewRetryBlocked = (groupId: string) => {
    const live = current.current;
    return !mounted.current || !owner.isCurrent(generation) || !navigator.onLine || live.locked || !live.props.canPrepare ||
      !live.owner.recoveryAvailable() || currentCommandBlocked() || recordingSchemaPreviewRetryBlocked(fenceInput.current, groupId);
  };
  const dispatch = (intent: RecordingNextIntent) => {
    const live = current.current;
    if (!owner.isCurrent(owner.epoch) || !live.selected || (intent.kind !== 'inspect' &&
      JSON.stringify(live.selected.model.intent) !== JSON.stringify(intent))) return;
    if (intent.kind === 'inspect' || intent.kind === 'review_schema') { inspect(); return; }
    if (intent.kind === 'recover_operation') { focus('[data-testid="basic-recording-status"], [data-testid="workspace-test-write-recovery"], [data-testid="recording-guided-preparation"]'); return; }
    if (intent.kind === 'choose_sources') { focus('[data-testid="circuit-profile-sources-toggle"], [data-testid="circuit-profile-toggle-instantaneous"]'); return; }
    if (intent.kind === 'resolve_semantics') { focus('[data-testid^="circuit-counter-confirmation"] select, [data-testid^="circuit-counter-confirmation"] input, [data-testid="circuit-profile-measurements-error"] button'); return; }
    if (intent.kind === 'review_scoped_start') {
      if (!props.prepareOnly) {
        const opener = region.current?.querySelector<HTMLButtonElement>('button[data-primary="true"]');
        reviewFocus.current = opener?.dataset.testid ? { epoch: owner.epoch, testId: opener.dataset.testid, dismissed: false } : undefined;
        setReviewEpoch(owner.epoch);
      }
      return;
    }
    if (!live.selected.model.enabled || busy || start.busy) return;
    if (intent.kind === 'query_operation') {
      if (intent.operationOwner === 'test_write') {
        const reference = owner.references.find(value => value.operation_id === intent.operationId);
        const query = reference ? testWriteQueries[JSON.stringify([queryIdentity, reference])] : undefined;
        if (query && query.operationId === intent.operationId) void run(query.check, owner.epoch, false);
        return;
      }
      void run(intent.operationOwner === 'schema' ? schema?.commands.check : start.check, owner.epoch, false); return;
    }
    if (intent.kind === 'refresh_evidence') { void run(selectedGroup ? evidence[selectedGroup.id]?.refresh : undefined, owner.epoch, false); return; }
    if (!live.props.canPrepare || live.props.readonly) return;
    if (intent.kind === 'save_prepare') void run(start.prepare, owner.epoch);
    if (intent.kind === 'preview_schema') void run(schema?.commands.preview, owner.epoch);
    if (intent.kind === 'confirm_schema') void run(schema?.commands.confirm, owner.epoch);
  };
  const intent = selected?.model.intent.kind;
  const completionPrimary = Boolean(props.completedOperationId && 'operationId' in start.observation &&
    start.observation.operationId === props.completedOperationId && (intent === 'review_scoped_start' || intent === 'monitor'));
  const commandUnavailable = intent === 'preview_schema' ? !schema?.commands.preview : intent === 'confirm_schema' ? !schema?.commands.confirm : false;
  const command = intent === 'save_prepare' || intent === 'preview_schema' || intent === 'confirm_schema';
  const testId = intent === 'save_prepare' ? `${props.testIdPrefix}-prepare` : intent === 'review_scoped_start' ? `${props.testIdPrefix}-start`
    : intent === 'preview_schema' ? 'group-schema-preview' : intent === 'confirm_schema' ? 'group-schema-apply' : intent === 'query_operation' ? 'basic-recording-check' : 'recording-next-action';
  const focusReturn = reviewFocus.current;
  return <section ref={region} tabIndex={-1} data-testid="recording-guided-preparation" className="space-y-3">
    {groups.map(group => <EvidenceOwner key={`${scope.workspaceId}|${group.id}`} group={group} scope={groupScope(group)} report={reportEvidence} />)}
    {!open && props.summary?.(scope)}
    {!open && !completionPrimary && selected ? <RecordingNextAction input={selected.input} labels={labels}
      scopeLabel={safeRuntimeLabel(scope.draft.circuitName) || t('selectedDevice')} onIntent={dispatch} testId={testId}
      disabled={Boolean(intent === 'review_scoped_start' && props.prepareOnly) || Boolean(command && (!props.canPrepare || locked || commandUnavailable)) || intent === 'query_operation' && (busy || start.busy)}
      blockedReason={intent === 'review_scoped_start' && props.prepareOnly ? t('next.preparedOnly') : command && (!props.canPrepare || commandUnavailable) ? t('next.fixControls') : undefined}
      onRepair={repair} repairLabel={t('next.repair')} /> : !open && !completionPrimary && <>
      <p role="status">{t('next.fixControls')}</p><button type="button" data-primary="true" onClick={repair}>{labels.actions.inspect}</button>
    </>}
    {!open && props.renderStartResult?.(completionPrimary)}
    {props.renderPreparation({ online: owner.online, consentEpoch: generation, onPresentation: reportSchema, newCommandsBlocked: blocked,
      blockedGroupIds, otherOwnerBlockedGroupIds, isNewOperationBlocked: commandFence.isBlocked, isSchemaPreviewRetryBlocked })}
    {[...checkedReferences.current.references.values()].map(reference => <TestWriteOwner key={JSON.stringify([queryIdentity, reference])}
      reference={reference} identity={queryIdentity} report={reportTestWriteQuery} />)}
    <Dialog.Root open={open} onOpenChange={value => { if (!value) setReviewEpoch(undefined); }}>
      <Dialog.Portal><Dialog.Overlay className="fixed inset-0 z-40 bg-black/70" />
        <Dialog.Content onEscapeKeyDown={() => { if (focusReturn) focusReturn.dismissed = true; }}
          onCloseAutoFocus={event => {
            event.preventDefault();
            if (reviewFocus.current === focusReturn) reviewFocus.current = undefined;
            // The primary unmounts during review; restore only its current same-owner replacement, never an old DOM node.
            const target = region.current?.querySelector<HTMLButtonElement>('button[data-primary="true"]');
            if (focusReturn?.dismissed && mounted.current && current.current.owner.isCurrent(focusReturn.epoch) &&
              target?.isConnected && !target.disabled && target.dataset.testid === focusReturn.testId) target.focus();
          }} className="fixed left-1/2 top-1/2 z-50 max-h-[85vh] w-[min(94vw,48rem)] -translate-x-1/2 -translate-y-1/2 space-y-4 overflow-y-auto rounded-xl border border-slate-700 bg-slate-950 p-5 text-slate-100">
          <Dialog.Title>{t('next.startTitle')}</Dialog.Title><Dialog.Description>{t('next.startDescription')}</Dialog.Description>
          {props.summary?.(scope)}
          <p role="status">{t('next.deliverySeparate')}</p>
          <button type="button" data-primary="true" data-testid="recording-start-confirm"
            disabled={!allReady || !props.canPrepare || locked || !owner.online || props.prepareOnly}
            className="rounded bg-blue-600 px-3 py-2 disabled:opacity-50" onClick={() => {
              const consent = reviewEpoch;
              if (consent === undefined || !owner.isCurrent(consent) || !current.current.allReady || current.current.locked || !current.current.props.canPrepare) return;
              setReviewEpoch(undefined); void run(start.start, consent);
            }}>{t('next.confirmStart')}</button>
          <Dialog.Close asChild><button type="button" data-testid="recording-start-cancel" onClick={() => { if (focusReturn) focusReturn.dismissed = true; }}
            className="ml-2 rounded border border-slate-600 px-3 py-2">{t('next.cancel')}</button></Dialog.Close>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  </section>;
}
