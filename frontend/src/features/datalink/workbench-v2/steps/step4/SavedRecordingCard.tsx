import { useTranslation } from 'react-i18next';
import { useWriteGroupDeliveryQuery } from '../../../../../hooks/datalink/useStudioV2WriteGroups';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import { SavedRecordingManagement, type SavedRecordingManagementProps } from './SavedRecordingManagement';
import { safeRuntimeLabel } from '../../../../../utils/runtimeOperatorSafety';

export interface SavedRecordingCardProps {
  group: WriteGroup;
  current: boolean | null;
  highlighted?: boolean;
  onOpen: () => void;
  management?: Omit<SavedRecordingManagementProps, 'group'>;
}

/** One saved group, with its existing delivery read and a navigation-only opener. */
export function SavedRecordingCard({ group, current, onOpen, management, highlighted }: SavedRecordingCardProps) {
  const { t } = useTranslation('workbench-v2');
  const { t: guided } = useTranslation('workbench-v2-recording-guided');
  const { t: runtime } = useTranslation('runtime-dashboard');
  const delivery = useWriteGroupDeliveryQuery(group.id, true);
  return <li tabIndex={highlighted ? -1 : undefined} data-locator-selected={highlighted || undefined} className="flex flex-wrap items-center justify-between gap-2 rounded border border-[var(--gw-border)] p-3 text-xs text-[var(--gw-text)]" data-testid={`group-list-item-${group.id}`}>
    <div>
      <div className="text-sm font-medium text-[var(--gw-text)]">{safeRuntimeLabel(group.name) || guided('management.savedNameFallback')}</div>
      <div className="text-[var(--gw-muted)]">
        {safeRuntimeLabel(group.destination.table_name) || guided('management.tableUnavailable')} · {t(`step4.group.status.${group.status}`)} · {group.members.length} {t('step4.group.list.members')}
        {current === null && <span className="ml-2">{runtime('locatorEntry.destinationUnknown')}</span>}
        {current === false && <span className="ml-2 text-[var(--gw-warning-fg)]" data-testid={`group-list-mismatch-${group.id}`}>{t('step4.group.list.other_destination')}</span>}
      </div>
      {!delivery.isError && delivery.data?.group_id === group.id && <div className="text-[var(--gw-muted)]" data-testid={`group-list-delivery-${group.id}`}>
        {t('step4.group.list.delivery', { committed: delivery.data.stages.sql_committed, waiting: delivery.data.stages.queued + delivery.data.stages.retrying })}
      </div>}
    </div>
    <button type="button" onClick={onOpen} className="rounded border border-[var(--gw-control)] px-3 py-1 text-[var(--gw-text)]" data-testid={`group-open-${group.id}`}>
      {guided('management.openDetails')}
    </button>
    {management?.lifecycleHost && <div className="w-full space-y-2">
      {management.readonly && <p role="status" className="text-[var(--gw-warning-fg)]" data-testid={`saved-recording-management-unavailable-${group.id}`}>{guided('management.actionsUnavailable')}</p>}
      <SavedRecordingManagement {...management} group={group} />
    </div>}
  </li>;
}
