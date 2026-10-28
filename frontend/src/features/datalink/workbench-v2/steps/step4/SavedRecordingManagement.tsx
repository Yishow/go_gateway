import { useTranslation } from 'react-i18next';
import { useWriteGroupsQuery } from '@/hooks/datalink/useStudioV2WriteGroups';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { GroupLifecycleBar } from './writeGroup/GroupLifecycleBar';
import type { SavedRecordingActionHost } from './writeGroup/lifecycle/SavedRecordingActions';

export interface SavedRecordingManagementProps {
  group: WriteGroup;
  workspaceId?: string;
  workspaceRevision?: string;
  readonly?: boolean;
  lifecycleHost?: SavedRecordingActionHost;
}

/** Saved controls do not call ordinary ensure/create and never consume an edited Basic draft. */
export function SavedRecordingManagement({ group, workspaceId, workspaceRevision, readonly = false, lifecycleHost }: SavedRecordingManagementProps) {
  const groups = useWriteGroupsQuery(Boolean(workspaceId));
  return <GroupLifecycleBar group={group} workspaceId={workspaceId ?? ''} workspaceRevision={workspaceRevision ?? ''}
    canApply={false} hideApply readonly={readonly} lifecycleHost={lifecycleHost}
    onChanged={() => void groups.refetch()} onConflict={() => void groups.refetch()} />;
}

/** Tombstones are management-only; they never re-enter Basic's ordinary creation candidate list. */
export function RemovedRecordingManagement({ groups, deviceId, onOpenAdvanced, ...props }: Omit<SavedRecordingManagementProps, 'group'> & {
  groups: WriteGroup[]; deviceId: string; onOpenAdvanced?: (id: string) => void;
}) {
  const { t } = useTranslation('recording-lifecycle');
  const removed = groups.filter((group) => group.status === 'deleted' &&
    (group.basic_managed_device_id === deviceId || group.members.some((member) => member.device_id === deviceId)));
  if (!removed.length) return null;
  return <details className="rounded border border-slate-700 p-3 text-xs text-slate-300" data-testid="basic-removed-recordings">
    <summary className="cursor-pointer font-semibold">{t('removed.title', { count: removed.length })}</summary>
    <p className="my-2">{t('removed.note')}</p>
    {removed.map((group) => <section key={group.id} className="space-y-2 py-2" data-testid={`basic-removed-${group.id}`}>
      <h4>{group.name}</h4>
      {onOpenAdvanced && <button type="button" onClick={() => onOpenAdvanced(group.id)} className="underline" data-testid={`basic-removed-open-${group.id}`}>{t('removed.manage')}</button>}
      {props.lifecycleHost && <SavedRecordingManagement {...props} group={group} />}
    </section>)}
  </details>;
}
