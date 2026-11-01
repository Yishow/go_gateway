import { useEffect, useRef, useState, type RefObject } from 'react';
import { useTranslation } from 'react-i18next';
import type { StudioRecordingRepair } from '../../../../services/studioRecordingView';
import type { StudioLocatorEntryRead } from '../../../../types/studioV2LocatorEntry';
import { inspectSchemaOperation } from '../steps/step4/writeGroup/schemaOperationRecovery';
import { recordingTestWriteRecovery } from '../state/recordingTestWriteRecovery';
import { loadBasicRecordingIntent, subscribeBasicRecordingIntent } from '../state/basicRecording';
import { useRecordingNavigation } from './RecordingDraftNavigation';

interface Props {
  repair: StudioRecordingRepair; entry: StudioLocatorEntryRead; region: RefObject<HTMLElement | null>;
  canUseScope: () => boolean; onReviewStart: () => void; onAdvanced?: () => void;
}
/** Navigation only. Retained identity selects an existing control, never a request or permission. */
export function RecordingRepairFocus({ repair, entry, region, canUseScope, onReviewStart, onAdvanced }: Props) {
  const { t } = useTranslation('recording-editor'), { t: home } = useTranslation('recording-home'), { t: runtime } = useTranslation('runtime-dashboard');
  const navigation = useRecordingNavigation();
  const [, refresh] = useState(0);
  const ownerKey = `${repair.workspaceId}/saved-record/${repair.groupId}`;
  useEffect(() => {
    const update = () => refresh(value => value + 1);
    const unsubscribe = recordingTestWriteRecovery.subscribe(update), stop = subscribeBasicRecordingIntent(ownerKey, update);
    return () => { unsubscribe(); stop(); };
  }, [ownerKey]);
  const target = entry.selection.state === 'selected' ? entry.selection.target : undefined;
  const list = entry.entry.write_groups, group = list?.groups.length === 1 ? list.groups[0] : undefined;
  const exact = entry.entry.editor_hydrated && entry.workspace.id === repair.workspaceId && list?.workspace_id === repair.workspaceId &&
    list.workspace_revision === repair.workspaceRevision && entry.entry.workspace_revision === repair.workspaceRevision &&
    target?.step === 4 && target.workspaceId === repair.workspaceId && target.workspaceRevision === repair.workspaceRevision &&
    target.groupId === repair.groupId && target.groupRevision === repair.groupRevision && group?.workspace_id === repair.workspaceId &&
    group.id === repair.groupId && group.revision === repair.groupRevision && group.status !== 'deleted';
  const operation = repair.operation;
  const matchingOperation = () => {
    if (!operation || !group?.members.some(member => member.device_id === operation.deviceId &&
      (operation.circuitId === undefined || member.circuit_id === operation.circuitId))) return false;
    if (operation.owner === 'schema') {
      const held = inspectSchemaOperation(repair.workspaceId, repair.groupId);
      return held.state === 'retained' && held.operationId === operation.operationId;
    }
    if (operation.owner === 'test_write') return recordingTestWriteRecovery.read().references.some(reference =>
      reference.workspace_id === repair.workspaceId && reference.group_id === repair.groupId && reference.operation_id === operation.operationId);
    // A legacy device/circuit request cannot be reinterpreted as a saved-record request.
    const held = loadBasicRecordingIntent(ownerKey);
    return held?.owner?.kind === 'saved_record' && held.owner.groupId === repair.groupId && held.request.device_ids.includes(operation.deviceId) &&
      (operation.kind === 'recover_operation' || held.operation_id === operation.operationId);
  };
  const admitted = Boolean(exact && canUseScope() && (repair.section !== 'operation' || matchingOperation()));
  const consumed = useRef<string | undefined>(undefined), identity = JSON.stringify(repair);
  const live = useRef({ admitted, canUseScope, matchingOperation }); live.current = { admitted, canUseScope, matchingOperation };
  const selector = repair.section === 'sources' || repair.section === 'destination' ? `#recording-editor-${repair.section}`
    : operation?.owner === 'schema' ? '[data-testid="recording-preparation-lookup"]'
      : operation?.owner === 'test_write' ? '[data-testid^="group-test-write-recovery-query-"]'
        : operation?.owner === 'start' ? '[data-testid="recording-repair-start"]' : '[data-testid="group-editor"]';
  useEffect(() => {
    const root = region.current;
    if (!admitted || !root || consumed.current === identity) return;
    const focus = () => {
      if (!live.current.admitted || !live.current.canUseScope() || repair.section === 'operation' && !live.current.matchingOperation()) return false;
      const candidates = [...root.querySelectorAll<HTMLElement>(selector)];
      const node = operation?.owner === 'test_write' ? candidates.find(item => item.dataset.testid === `group-test-write-recovery-query-${operation.operationId}`) : candidates[0];
      if (!node) return false;
      for (let parent = node.parentElement; parent && parent !== root; parent = parent.parentElement) if (parent instanceof HTMLDetailsElement) parent.open = true;
      if (!node.matches('button,input,select,textarea,a,[tabindex]')) node.tabIndex = -1;
      consumed.current = identity; node.scrollIntoView?.({ block: 'center' }); node.focus({ preventScroll: true }); return true;
    };
    if (focus()) return;
    const observer = new MutationObserver(() => { if (focus()) observer.disconnect(); });
    observer.observe(root, { childList: true, subtree: true }); return () => observer.disconnect();
  }, [admitted, identity, operation?.owner, operation?.operationId, region, repair.section, selector]);
  if (!admitted) return <aside role="alert" data-testid="recording-repair-unavailable"><p>{operation ? t('start.recoveryUnavailable') : runtime('locatorEntry.stale')}</p>
    {onAdvanced && <button type="button" onClick={() => navigation.request(onAdvanced)}>{home('advanced')}</button>}</aside>;
  return operation?.owner === 'start' ? <button type="button" data-testid="recording-repair-start" onClick={() => {
    if (live.current.admitted && live.current.canUseScope() && live.current.matchingOperation()) navigation.request(() => { if (live.current.canUseScope() && live.current.matchingOperation()) onReviewStart(); });
  }}>{t('start.lookup')}</button> : null;
}
