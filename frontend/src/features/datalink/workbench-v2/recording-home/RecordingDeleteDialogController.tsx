import { useState, type ReactNode } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { recordingHomeSourcePrefix } from '@/hooks/datalink/useRecordingHomeSources';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { RecordingReadonly } from '../state/recordingWorkflowView';
import { buildSimpleRecordingTarget, buildStudioPresentationTarget } from '@/services/studioRecordingView';
import { DeviceDeleteDialog } from './DeviceDeleteDialog';
import { RecordingDeleteDialog } from './RecordingDeleteDialog';
import { RecordingStopConfirmation } from './RecordingHomeRecordRow';
import type { HomeDeviceSummary } from './recordingHomeCollections';
import type { RecordingHomeLocation } from './RecordingHome';

type SavedGroup = RecordingReadonly<WriteGroup>;
type Selection = { kind: 'recording' | 'stop'; group: SavedGroup; revision: string; returnFocus: HTMLElement | null } |
  { kind: 'device'; device: HomeDeviceSummary; revision: string };

/** Lives outside the revision-keyed home so a deletion ACK survives list refresh. */
export function RecordingDeleteDialogController({ workspaceId, workspaceRevision, groups, location, onNavigate, refresh, children }: {
  workspaceId: string; workspaceRevision: string; groups: readonly SavedGroup[]; location: RecordingHomeLocation;
  onNavigate?: (target: string) => void; refresh: () => Promise<boolean>;
  children: (actions: { onDelete: (group: SavedGroup) => void; onDeleteDevice: (device: HomeDeviceSummary) => void; notice?: ReactNode }) => ReactNode;
}) {
  const { t } = useTranslation('recording-home');
  const queryClient = useQueryClient();
  const [selection, setSelection] = useState<Selection | null>(null);
  const [notice, setNotice] = useState<string>();
  const [refreshFailed, setRefreshFailed] = useState(false);
  const refreshAfterUnknown = async () => {
    setRefreshFailed(false);
    try {
      const [fresh] = await Promise.all([refresh(), queryClient.refetchQueries({ queryKey: recordingHomeSourcePrefix({}, workspaceId), type: 'active' }, { throwOnError: true })]);
      if (!fresh) throw new Error('Refresh unavailable');
      setSelection(null); setNotice('deleteDialogs.refreshComplete');
    } catch { setRefreshFailed(true); }
  };
  const current = selection && (selection.kind === 'device' || selection.group.workspace_id === workspaceId) ? selection : null;
  return <>
    {children({
      onDelete: group => { setNotice(undefined); setRefreshFailed(false); setSelection({ kind: 'recording', group, revision: workspaceRevision,
        returnFocus: document.activeElement instanceof HTMLElement ? document.activeElement : null }); },
      onDeleteDevice: device => { setNotice(undefined); setRefreshFailed(false); setSelection({ kind: 'device', device, revision: workspaceRevision }); },
      notice: notice ? <p role="status" className="gw-notice">{t(notice)}</p> : undefined,
    })}
    {current?.kind === 'recording' && <RecordingDeleteDialog open group={current.group} workspaceRevision={current.revision} refreshFailed={refreshFailed}
      onClose={() => setSelection(null)} onDeleted={() => { setSelection(null); setNotice('deleteDialogs.recording.success'); }}
      onStop={() => setSelection({ ...current, kind: 'stop' })}
      onOpenRecording={() => {
        const target = current.group.status === 'deleted' ? buildStudioPresentationTarget('?step=output&target=database', location.hash, 'advanced') :
          buildSimpleRecordingTarget({ kind: 'edit', workspaceId, groupId: current.group.id, groupRevision: current.group.revision }, location.hash);
        setSelection(null); if (target) onNavigate?.(target);
      }} onRefresh={() => { void refreshAfterUnknown(); }} />}
    {current?.kind === 'stop' && <RecordingStopConfirmation group={current.group} scope={{ workspaceId, workspaceRevision: current.revision, epoch: 0 }}
      returnFocus={current.returnFocus} onClose={() => setSelection(null)} onError={() => { setSelection(null); setNotice('deleteDialogs.recording.stopUnavailable'); }} />}
    {current?.kind === 'device' && <DeviceDeleteDialog open device={current.device} workspaceId={workspaceId} recordings={groups} refreshFailed={refreshFailed} scopeValid={current.revision === workspaceRevision}
      onClose={() => setSelection(null)} onDeleted={() => { setSelection(null); setNotice('deleteDialogs.device.success'); }} onRefresh={() => { void refreshAfterUnknown(); }} />}
  </>;
}
