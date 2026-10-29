import { useTranslation } from 'react-i18next';
import { useWriteGroupDeliveryQuery } from '../../../../../hooks/datalink/useStudioV2WriteGroups';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import { SavedRecordingManagement, type SavedRecordingManagementProps } from './SavedRecordingManagement';
import { safeRuntimeLabel } from '../../../../../utils/runtimeOperatorSafety';

export interface SavedRecordingCardProps {
  group: WriteGroup;
  current: boolean;
  onOpen: () => void;
  management?: Omit<SavedRecordingManagementProps, 'group'>;
}

/** One saved group, with its existing delivery read and a navigation-only opener. */
export function SavedRecordingCard({ group, current, onOpen, management }: SavedRecordingCardProps) {
  const { t } = useTranslation('workbench-v2');
  const { t: guided } = useTranslation('workbench-v2-recording-guided');
  const delivery = useWriteGroupDeliveryQuery(group.id, true);
  return <li className="flex flex-wrap items-center justify-between gap-2 rounded border border-slate-800 p-3 text-xs text-slate-300" data-testid={`group-list-item-${group.id}`}>
    <div>
      <div className="text-sm font-medium text-slate-100">{safeRuntimeLabel(group.name) || guided('management.savedNameFallback')}</div>
      <div className="text-slate-500">
        {safeRuntimeLabel(group.destination.table_name) || guided('management.tableUnavailable')} · {t(`step4.group.status.${group.status}`)} · {group.members.length} {t('step4.group.list.members')}
        {!current && <span className="ml-2 text-amber-300" data-testid={`group-list-mismatch-${group.id}`}>{t('step4.group.list.other_destination')}</span>}
      </div>
      {!delivery.isError && delivery.data?.group_id === group.id && <div className="text-slate-500" data-testid={`group-list-delivery-${group.id}`}>
        {t('step4.group.list.delivery', { committed: delivery.data.stages.sql_committed, waiting: delivery.data.stages.queued + delivery.data.stages.retrying })}
      </div>}
    </div>
    <button type="button" onClick={onOpen} className="rounded border border-slate-600 px-3 py-1 text-slate-100" data-testid={`group-open-${group.id}`}>
      {guided('management.openDetails')}
    </button>
    {management?.lifecycleHost && <div className="w-full space-y-2">
      {management.readonly && <p role="status" className="text-amber-300" data-testid={`saved-recording-management-unavailable-${group.id}`}>{guided('management.actionsUnavailable')}</p>}
      <SavedRecordingManagement {...management} group={group} />
    </div>}
  </li>;
}
