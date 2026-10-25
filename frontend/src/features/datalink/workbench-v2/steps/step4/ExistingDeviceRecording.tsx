import type { FC } from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';

/** A saved legacy group remains editable; its presence never selects a different setup workflow. */
export const ExistingDeviceRecording: FC<{
  group: WriteGroup;
  disabled?: boolean;
  onEdit?: () => void;
  onOpenAdvanced?: (groupId: string) => void;
}> = ({ group, disabled = false, onEdit, onOpenAdvanced }) => {
  const { t } = useTranslation('workbench-v2');
  const applied = Boolean(group.applied_revision) && group.status !== 'disabled';
  return (
    <aside data-testid="device-recording-existing" className="space-y-2 rounded border border-slate-700 bg-slate-900/40 p-3 text-xs text-slate-300">
      <p className="font-semibold">{group.name}</p>
      <p>{t('step4.basic.existing.latest', { count: group.row_policy.interval_seconds })}</p>
      <p data-testid={applied ? 'device-recording-existing-applied' : 'device-recording-existing-draft'}>
        {t(applied ? 'step4.basic.existing.applied' : 'step4.basic.existing.draft')}
      </p>
      <div className="flex flex-wrap gap-3">
        {onEdit && <button type="button" data-testid="device-recording-edit-existing" disabled={disabled} onClick={onEdit} className="text-cyan-300 underline disabled:opacity-40">{t('step4.basic.existing.edit')}</button>}
        {applied && onOpenAdvanced && <button type="button" data-testid="device-recording-review-cutover" disabled={disabled} onClick={() => onOpenAdvanced(group.id)} className="text-cyan-300 underline disabled:opacity-40">{t('step4.basic.existing.review')}</button>}
      </div>
    </aside>
  );
};
