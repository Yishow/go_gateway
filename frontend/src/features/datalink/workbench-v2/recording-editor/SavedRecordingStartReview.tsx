import { inspectBasicRecordingIntent } from '../state/basicRecording';
import { sameRecordingDeviceIds } from '../state/recordingStartOwner';
import { safeBasicNextAction, safeBasicOperationMessage } from '../state/basicRecordingMessages';
import type { PointOrderObservation } from './useRecordingPointOrder';
import { buildGroupCandidates } from '../state/writeGroup/candidates';
import { SavedRecordingStartEvidence } from './SavedRecordingStartEvidence';
import { useEffect, useRef, useState, useSyncExternalStore } from 'react';
import * as Dialog from '@radix-ui/react-dialog';
import { useTranslation } from 'react-i18next';
import { Button } from '../../../../components/ui/button';
import type { WriteGroup, WriteGroupListResponse } from '../../../../types/studioV2WriteGroup';
import type { WorkbenchV2State } from '../state/types';
import type { ModbusShareStatus } from '../../../../types/modbusShare';
import type { DeviceSourcePort } from '../../../../pages/datalink/workbench-v2/deviceSourceSubmission';
import { studioV2WorkspaceWriteGroupsAPI } from '../../../../services/studioV2WorkspaceWriteGroups';
import { useCircuitProfileStart } from '../steps/step4/useCircuitProfileStart';
import { inspectSchemaOperation } from '../steps/step4/writeGroup/schemaOperationRecovery';
import { recordingTestWriteRecovery } from '../state/recordingTestWriteRecovery';
import { useRecordingDraftGuard, useRecordingNavigation, useRecordingReadUnavailable } from './RecordingDraftNavigation';
import { savedRecordingStartDevices, savedRecordingStartSourceIdentity, savedRecordingStartSourcesSaved } from './savedRecordingStartScope';
import { captureReturnFocus, restoreReturnFocus, type ReturnFocusTicket } from '../../../../utils/returnFocus';
/** DOM-only return after a same-canonical caller readback; never a command permission. */
export interface SavedRecordingStartFocusReturn { ticket: ReturnFocusTicket; currentTarget: () => HTMLElement | null }
export interface SavedRecordingStartReviewProps {
  group: WriteGroup; workspaceId: string; workspaceRevision: string; state: WorkbenchV2State; readonly: boolean;
  fetchLatest: () => Promise<WriteGroupListResponse | undefined>;
  onClose: () => void; onOperationActiveChange?: (active: boolean) => void;
  shareStatus?: ModbusShareStatus | null; deviceSource?: DeviceSourcePort;
  /** Denial-only live caller lease; never expands the saved-record scope. */
  canUseScope?: () => boolean;
  focusReturn?: SavedRecordingStartFocusReturn;
  onCompleted?: () => void;
}
const emptyPointObservations: readonly PointOrderObservation[] = [];
const noSubscribe = () => () => undefined;
const buttonStyle = 'max-w-full h-auto min-h-9 whitespace-normal break-words border border-[var(--gw-control)] bg-[var(--gw-surface)] px-3 py-2 text-[var(--gw-text)] hover:bg-[var(--gw-active)] disabled:bg-[var(--gw-surface)] disabled:text-[var(--gw-muted)] focus-visible:outline-solid focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--gw-focus)]';
/** Explicit saved-record review, using the original retained request and scoped-start operation owner. */
export function SavedRecordingStartReview({ group, workspaceId, workspaceRevision, state, readonly, fetchLatest, onClose, onOperationActiveChange, shareStatus, deviceSource, canUseScope, onCompleted, focusReturn }: SavedRecordingStartReviewProps) {
  const { t } = useTranslation('recording-editor');
  const { t: startT } = useTranslation('workbench-v2');
  const navigation = useRecordingNavigation(), readUnavailable = useRecordingReadUnavailable();
  const devices = savedRecordingStartDevices(group);
  const candidates = buildGroupCandidates(state).candidates;
  const pointObservations = useSyncExternalStore(deviceSource?.pointOrder?.subscribe ?? noSubscribe, () => deviceSource?.pointOrder?.observations?.() ?? emptyPointObservations);
  const selectedPointObservations = pointObservations.filter(({ scope }) => scope.workspace_id === workspaceId && group.members.some(member => member.device_id === scope.device_id && member.point_id === scope.point_id));
  const scope = `${workspaceId}/saved-record/${group.id}`;
  const identity = JSON.stringify([workspaceId, workspaceRevision, group, savedRecordingStartSourceIdentity(group, state), selectedPointObservations, readonly, readUnavailable, shareStatus?.settings_revision, shareStatus?.workspace_revision]);
  const current = useRef(identity); current.current = identity;
  const accepted = useRef(identity);
  const focusAvailable = canUseScope?.() !== false;
  const focusOwner = useRef({ identity, available: focusAvailable });
  if (focusOwner.current.identity !== identity || focusOwner.current.available !== focusAvailable) focusOwner.current = { identity, available: focusAvailable };
  const returnFocus = useRef<ReturnFocusTicket | null>(null), closing = useRef(false);
  const closeTarget = useRef<HTMLElement | null>(null);
  const mounted = useRef(true);
  useEffect(() => { mounted.current = true; return () => { mounted.current = false; if (!closing.current) focusOwner.current = { ...focusOwner.current }; }; }, []);
  const [error, setError] = useState(false);
  const recoveryUnavailable = () => inspectBasicRecordingIntent(scope).state === 'unavailable' || recordingTestWriteRecovery.read().unavailable;
  const denied = () => canUseScope?.() === false || recoveryUnavailable() || group.workspace_id !== workspaceId || readonly || readUnavailable || !navigator.onLine || shareStatus?.workspace_id !== workspaceId || !shareStatus?.readiness_token || !shareStatus.settings_revision || !shareStatus.workspace_revision || !savedRecordingStartSourcesSaved(group, state) ||
    group.members.some(member => deviceSource?.range?.sourceWriteState?.(member.device_id, member.point_id)) ||
    deviceSource?.pointOrder?.observations?.().some(({ scope: pointScope, entry }) => pointScope.workspace_id === workspaceId &&
      group.members.some(member => member.device_id === pointScope.device_id && member.point_id === pointScope.point_id) &&
      (Boolean(entry.intent) && !entry.saved && !entry.recovered || ['saving', 'unknown', 'inspecting', 'recovering', 'recovery-unknown'].includes(entry.status))) ||
    inspectSchemaOperation(workspaceId, group.id).state !== 'empty' ||
    recordingTestWriteRecovery.read().references.some(ref => ref.workspace_id === workspaceId && ref.group_id === group.id);
  const valid = current.current === accepted.current && !denied();
  const flow = useCircuitProfileStart({ scopeKey: scope, canStart: valid, readonly: readonly || readUnavailable, workspaceId,
    deviceId: devices[0] ?? '', savedRecord: { kind: 'saved_record', groupId: group.id, deviceIds: devices },
    connectorRevision: group.destination.connector_revision,
    readinessToken: shareStatus?.readiness_token, settingsRevision: shareStatus?.settings_revision, shareWorkspaceRevision: shareStatus?.workspace_revision,
    retryBlocked: denied, onOperationActiveChange,
    prepare: async () => {
      const lease = identity;
      const stillCurrent = () => mounted.current && current.current === lease && accepted.current === lease && !denied();
      if (!stillCurrent()) { if (mounted.current) setError(true); return undefined; }
      setError(false);
      try {
        const latest = await fetchLatest();
        const matches = latest?.groups.filter(candidate => candidate.id === group.id) ?? [];
        const saved = matches.length === 1 ? matches[0] : undefined;
        if (!stillCurrent() || latest?.workspace_id !== workspaceId || latest.workspace_revision !== workspaceRevision || !saved || saved.workspace_id !== workspaceId ||
          saved.revision !== group.revision || saved.destination.connector_revision !== group.destination.connector_revision ||
          JSON.stringify(savedRecordingStartDevices(saved)) !== JSON.stringify(devices)) { if (mounted.current) setError(true); return undefined; }
        const ready = await studioV2WorkspaceWriteGroupsAPI.readiness(group.id);
        if (!stillCurrent() || ready.workspace_id !== workspaceId || ready.workspace_revision !== workspaceRevision || ready.group_id !== group.id || ready.group_revision !== group.revision || !ready.ready) {
          if (mounted.current) setError(true); return undefined;
        }
        return { groups: [saved], workspaceRevision };
      } catch { if (mounted.current) setError(true); return undefined; }
    },
  });
  useRecordingDraftGuard({ version: JSON.stringify([scope, flow.operationId, flow.status]), dirty: false, pending: flow.busy || flow.unresolved });
  const close = () => navigation.request(() => {
    closeTarget.current = focusReturn?.currentTarget() ?? null;
    if (canUseScope?.() === false) focusOwner.current = { ...focusOwner.current, available: false };
    closing.current = true; mounted.current = false; onClose();
  });
  const originalTerminal = flow.completion.state === 'available' && ['succeeded', 'failed'].includes(flow.completion.status) && !flow.unresolved;
  const matchesSaved = flow.completion.state === 'available' && sameRecordingDeviceIds(flow.completion.deviceIds, devices) &&
    flow.operation?.groups.some(progress => progress.group_id === group.id && progress.group_revision === group.revision);
  const terminalFailure = originalTerminal && flow.status === 'failed';
  const previousTerminal = originalTerminal && !terminalFailure && !matchesSaved;
  const canReviewAgain = previousTerminal || terminalFailure;
  const failureReason = terminalFailure ? safeBasicOperationMessage(flow.operation, startT) : undefined;
  const failureNextAction = terminalFailure ? safeBasicNextAction(flow.operation, startT) : undefined;
  const started = flow.completion.state === 'available' && flow.completion.status === 'succeeded' && matchesSaved;
  const completed = useRef<string | undefined>(undefined);
  useEffect(() => {
    if (!started || !flow.operationId || completed.current === flow.operationId) return;
    completed.current = flow.operationId; onCompleted?.();
  }, [started, flow.operationId, onCompleted]);
  return <Dialog.Root open onOpenChange={open => { if (!open) close(); }}><Dialog.Portal>
    <Dialog.Overlay className="gw-modal-overlay" />
    <Dialog.Content className="gw-modal bg-[var(--gw-surface)] text-[var(--gw-text)]" data-size="confirm" data-testid="saved-recording-start-review"
      onOpenAutoFocus={() => { returnFocus.current = captureReturnFocus(focusOwner.current, document.activeElement instanceof HTMLElement ? document.activeElement : null); }}
      onCloseAutoFocus={event => {
        event.preventDefault(); if (!closing.current) return;
        if (focusReturn) {
          const target = closeTarget.current;
          if (target && focusReturn.currentTarget() === target) restoreReturnFocus(focusReturn.ticket, focusReturn.ticket.owner, target);
        } else if (returnFocus.current?.owner === focusOwner.current) restoreReturnFocus(returnFocus.current, focusOwner.current);
      }}
      onEscapeKeyDown={event => { if (event.isComposing) event.preventDefault(); }}>
      <header className="gw-modal-header"><Dialog.Title className="text-lg font-semibold">{t('start.title')}</Dialog.Title>
      <Dialog.Description className="mt-2 text-sm text-[var(--gw-muted)]">{t('start.scopeNote')}</Dialog.Description></header>
      <div className="gw-modal-body space-y-3">
        <p className="font-semibold break-words">{group.name}</p>
        <ul>{devices.map(id => <li key={id}>{state.devices.find(device => device.id === id)?.name || id}<span className="ml-2 break-all text-xs text-[var(--gw-muted)]">{id}</span></li>)}</ul>
        <ul className="max-h-40 overflow-auto text-sm">{group.members.map(member => <li key={`${member.device_id}/${member.point_id}/${member.tag_id}`}>
          {state.devices.find(device => device.id === member.device_id)?.name || member.device_id} / {candidates.find(candidate => candidate.device_id === member.device_id && candidate.point_id === member.point_id && candidate.tag_id === member.tag_id)?.label || member.point_id}
        </li>)}</ul>
        <p className="text-sm text-[var(--gw-muted)]">{t('start.noSql')}</p>
        {(error || !valid) && !recoveryUnavailable() && !flow.recoveryUnavailable && (!flow.status || canReviewAgain) && <p role="alert">{t('start.changed')}</p>}
        {flow.error && <p role="alert">{flow.error}</p>}
        {(recoveryUnavailable() || flow.recoveryUnavailable) && <p role="alert">{t('start.recoveryUnavailable')}</p>}
        {flow.queryFailed && <p role="alert">{t('start.lookupFailed')}</p>}
        {started && flow.operation ? <SavedRecordingStartEvidence group={group} workspaceRevision={workspaceRevision} operation={flow.operation} />
          : terminalFailure ? <div role="alert"><p>{startT('step4.basic.operation.failed')}</p>
            {failureReason && <p>{failureReason}</p>}{failureNextAction && <p>{failureNextAction}</p>}</div>
          : flow.status && <p role="status">{t(previousTerminal ? 'start.previousTerminal' : flow.unresolved ? 'start.unknown' : 'start.incomplete')}</p>}
      </div>
      <footer className="gw-modal-footer">
        {(!flow.status || canReviewAgain) && <Button className={`${buttonStyle} gw-primary`} disabled={!valid || flow.busy} data-testid="saved-recording-start-confirm" onClick={() => void flow.start()}>{t('start.confirm')}</Button>}
        {flow.canLookupRequest && <Button className={buttonStyle} disabled={flow.busy} data-testid="saved-recording-start-lookup" onClick={() => void flow.check()}>{t('start.lookup')}</Button>}
        {flow.operationId && <Button className={buttonStyle} disabled={flow.busy} data-testid="saved-recording-start-check" onClick={() => void flow.check()}>{t('start.check')}</Button>}
        {flow.canRetrySameRequest && <Button className={buttonStyle} disabled={flow.busy} onClick={() => void flow.start(true)}>{t('start.retry')}</Button>}
        <Button className={buttonStyle} onClick={close}>{t('start.close')}</Button>
      </footer>
    </Dialog.Content>
  </Dialog.Portal></Dialog.Root>;
}
