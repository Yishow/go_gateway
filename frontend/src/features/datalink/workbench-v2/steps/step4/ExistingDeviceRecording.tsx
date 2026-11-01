import { SavedRecordingManagement, type SavedRecordingManagementProps } from './SavedRecordingManagement';
import type { FC } from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';

/** A saved legacy group remains editable; its presence never selects a different setup workflow. */
export const ExistingDeviceRecording: FC<{
  group: WriteGroup;
  lifecycle?: Omit<SavedRecordingManagementProps, 'group'>;
  disabled?: boolean;
  onEdit?: () => void;
  onOpenAdvanced?: (groupId: string) => void;
}> = ({ group, disabled = false, onEdit, onOpenAdvanced, lifecycle }) => {
  const { t } = useTranslation('workbench-v2');
  const applied = Boolean(group.applied_revision) && group.status !== 'disabled';
  return (
    <aside data-testid="device-recording-existing" className="space-y-2 rounded border border-[var(--gw-border)] bg-[var(--gw-surface)] p-3 text-xs text-[var(--gw-text)]">
      <p className="font-semibold">{group.name}</p>
      <p>{t('step4.basic.existing.latest', { count: group.row_policy.interval_seconds })}</p>
      <p data-testid={applied ? 'device-recording-existing-applied' : 'device-recording-existing-draft'}>
        {t(applied ? 'step4.basic.existing.applied' : 'step4.basic.existing.draft')}
      </p>
      <div className="flex flex-wrap gap-3">
        {onEdit && <button type="button" data-testid="device-recording-edit-existing" disabled={disabled} onClick={onEdit} className="text-[var(--gw-accent)] underline disabled:text-[var(--gw-muted)] disabled:bg-[var(--gw-ground)]">{t('step4.basic.existing.edit')}</button>}
        {applied && onOpenAdvanced && <button type="button" data-testid="device-recording-review-cutover" disabled={disabled} onClick={() => onOpenAdvanced(group.id)} className="text-[var(--gw-accent)] underline disabled:text-[var(--gw-muted)] disabled:bg-[var(--gw-ground)]">{t('step4.basic.existing.review')}</button>}
      </div>
      {lifecycle?.lifecycleHost && <SavedRecordingManagement {...lifecycle} group={group} readonly={disabled || lifecycle.readonly} />}
    </aside>
  );
};
