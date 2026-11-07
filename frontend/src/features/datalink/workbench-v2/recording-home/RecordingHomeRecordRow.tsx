import { Cpu, TriangleAlert } from 'lucide-react';
import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useWriteGroupLifecycleMutation } from '@/hooks/datalink/useStudioV2WriteGroups';
import { studioV2WriteGroupLifecycleAPI } from '@/services/studioV2WriteGroupLifecycle';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { WriteGroupLifecycle } from '@/types/studioV2WriteGroupLifecycle';
import type { RecordingReadonly } from '../state/recordingWorkflowView';
import { lifecycleMatches } from '../steps/step4/writeGroup/lifecycle/SavedRecordingActions';
import { SavedRecordingRemoveDialog } from '../steps/step4/writeGroup/lifecycle/SavedRecordingRemoveDialog';
import type { HomeDeviceSummary } from './recordingHomeCollections';
import type { recordingHomeResultSummary } from './recordingHomeResultSummary';
import { RecordingStatusBadge } from './RecordingHomeBadges';
import { RecordingHomeDestinationCell } from './RecordingHomeDestinationCell';
import { RecordingHomeResultFacts } from './RecordingHomeResultFacts';
import { RecordingHomeWriteStages } from './RecordingHomeWriteStages';
import { RecordingHomeRowActions } from './RecordingHomeRowActions';
import { recordingHomeLifecycleMatches } from './recordingHomeLifecycleMatches';
import type { RecordingHomeGuidanceItem } from './recordingHomeGuidance';
import type { RecordingHomeScope } from './recordingHomeProjection';

export function RecordingStopConfirmation({ group, scope, returnFocus, onClose, onError }: {
  group: RecordingReadonly<WriteGroup>; scope: RecordingHomeScope; returnFocus: HTMLElement | null;
  onClose: () => void; onError: (kind: 'reloadHint' | 'stopFailed') => void;
}) {
  const disable = useWriteGroupLifecycleMutation('disable');
  const [open, setOpen] = useState(false), [busy, setBusy] = useState(false);
  const mounted = useRef(false), working = useRef(false);
  const intent = { workspaceId: scope.workspaceId, workspaceRevision: scope.workspaceRevision, groupId: group.id, groupRevision: group.revision };
  const permitted = (fact: WriteGroupLifecycle) => lifecycleMatches(fact, intent) && fact.applied_revision === group.applied_revision &&
    fact.connector_id === group.destination.connector_id && fact.connector_revision === group.destination.connector_revision &&
    fact.actions.some(action => action.action === 'stop' && action.allowed);
  useEffect(() => {
    mounted.current = true;
    const abort = new AbortController();
    void studioV2WriteGroupLifecycleAPI.lifecycle(group.id, abort.signal).then(fact => {
      if (!mounted.current || abort.signal.aborted) return;
      if (permitted(fact)) setOpen(true); else onError('reloadHint');
    }).catch(() => { if (mounted.current && !abort.signal.aborted) onError('reloadHint'); });
    return () => { mounted.current = false; abort.abort(); };
    // One captured scope per mounted confirmation; a changed revision remounts it.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);
  const confirm = async () => {
    if (working.current) return;
    working.current = true; setBusy(true);
    let started = false;
    try {
      const fact = await studioV2WriteGroupLifecycleAPI.lifecycle(group.id);
      if (!mounted.current) return;
      if (!permitted(fact)) { onError('reloadHint'); return; }
      started = true;
      await disable.mutateAsync({ id: group.id, request: { workspace_id: scope.workspaceId, expected_workspace_revision: scope.workspaceRevision,
        expected_group_revision: group.revision, expected_connector_revision: group.destination.connector_revision } });
      if (mounted.current) onClose();
    } catch {
      if (mounted.current) onError(started ? 'stopFailed' : 'reloadHint');
    } finally {
      working.current = false;
      if (mounted.current) setBusy(false);
    }
  };
  return <SavedRecordingRemoveDialog action="stop" open={open} busy={busy} objectLabel={group.name} returnFocus={returnFocus}
    onClose={onClose} onConfirm={() => { void confirm(); }} />;
}

export function RecordingHomeRecordRow({ group, devices, summary, attention, evidence, onOpen, onDelete, onRecheck }: {
  group: RecordingReadonly<WriteGroup>; devices: readonly HomeDeviceSummary[];
  summary?: ReturnType<typeof recordingHomeResultSummary>; attention: boolean; evidence?: RecordingHomeGuidanceItem; onOpen: () => void;
  onDelete?: (group: RecordingReadonly<WriteGroup>) => void;
  onRecheck: () => void;
}) {
  const { t } = useTranslation('recording-home');
  const [stopping, setStopping] = useState(false), [stopError, setStopError] = useState<'reloadHint' | 'stopFailed' | null>(null);
  const [refreshing, setRefreshing] = useState(false);
  const returnFocus = useRef<HTMLElement | null>(null), refreshObserved = useRef(false);
  const lifecycle = evidence?.evidence.lifecycle;
  const fact = lifecycle?.state === 'available' ? lifecycle.value : undefined;
  const scope = evidence?.action.scope;
  const matches = recordingHomeLifecycleMatches(group, evidence);
  const canStop = matches && fact?.actions.some(action => action.action === 'stop' && action.allowed);
  const canDelete = matches && fact?.actions.some(action => action.action === 'remove' && action.allowed);
  useEffect(() => {
    if (!refreshing) return;
    if (lifecycle?.state !== 'available') {
      refreshObserved.current = true;
      if (!lifecycle || lifecycle.reason !== 'loading') setRefreshing(false);
    } else if (refreshObserved.current) {
      if (matches) setStopError(null);
      setRefreshing(false);
    }
  }, [refreshing, lifecycle, matches]);
  const actions = [
    ...(group.status === 'running' ? [{ id: 'stop', label: t('rowEvidence.stop'), disabled: !canStop || stopping || Boolean(stopError), onSelect: () => {
      returnFocus.current = document.activeElement instanceof HTMLElement ? document.activeElement : null;
      setStopError(null); setStopping(true);
    } }] : []),
    ...(onDelete && (!matches || canDelete) ? [{ id: 'delete', label: t('rowEvidence.delete'), disabled: !matches || Boolean(stopError), onSelect: () => onDelete(group) }] : []),
  ];
  return <tr>
    <td><strong>{group.name}</strong></td>
    <td><RecordingStatusBadge status={group.status} />{attention && <p className="gw-observed-issue"><TriangleAlert aria-hidden="true" size={14} />{t('observedAttentionShort')}</p>}</td>
    <td data-label={t('recordSources')}><span className="gw-home-flow" data-flow="source"><Cpu aria-hidden="true" size={14} /><span>{[...new Set(group.members.map(member => member.device_id))].map(id => devices.find(device => device.id === id)?.name ?? id).join(' / ') || t('unknown')}</span></span><p>{t('recordPointCount', { count: group.members.length })}</p><p className="gw-source-state">{t(`result.current.${summary?.current === 'source_attention' ? 'source_attention' : 'unknown'}`)}</p></td>
    <td data-label={t('destinationFilter')}><RecordingHomeDestinationCell destination={group.destination} /></td>
    <td data-label={t('recordWrites')}>{summary ? <RecordingHomeResultFacts showSourceState={false} id={group.id} name={group.name} summary={summary} /> : <><RecordingHomeWriteStages name={group.name} /><p>{t('statusUnknown')}</p></>}</td>
    <td><div className="gw-record-row-controls"><button type="button" data-home-focus={`record-${group.id}`} onClick={onOpen}>{t('openRecord')}</button>
      <RecordingHomeRowActions label={t('rowEvidence.actionsFor', { name: group.name })} items={actions} hint={!matches ? t('rowEvidence.reloadHint') : actions.length === 0 ? t('rowEvidence.noActions') : undefined} /></div>
      {stopping && scope && <RecordingStopConfirmation key={JSON.stringify([scope.workspaceId, scope.workspaceRevision, group.id, group.revision])}
        group={group} scope={scope} returnFocus={returnFocus.current} onClose={() => setStopping(false)} onError={kind => { setStopping(false); setStopError(kind); }} />}
      {stopError && <div role="alert"><p>{t(`rowEvidence.${stopError}`)}</p><button type="button" disabled={refreshing} onClick={() => { refreshObserved.current = false; setRefreshing(true); onRecheck(); }}>{t('retry')}</button></div>}
    </td>
  </tr>;
}
