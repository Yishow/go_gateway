import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroupStorageStrategy } from '../../../../../../types/studioV2WriteGroup';
import type { EditorDraft } from '../../../state/writeGroup/draft';

const FIELD = 'mt-1 w-full rounded border border-slate-700 bg-slate-950 px-2 py-1 text-sm text-slate-100';

export interface GroupIdentityFieldsProps {
  draft: EditorDraft;
  readonly: boolean;
  /** A profile group keeps the layout, storage kind and shared table it was created with. */
  locked: boolean;
  onChange: (patch: Partial<EditorDraft>) => void;
  onStorageStrategy: (strategy: WriteGroupStorageStrategy) => void;
}

/** Name, storage kind and destination table of a group. */
export const GroupIdentityFields: React.FC<GroupIdentityFieldsProps> = ({ draft, readonly, locked, onChange, onStorageStrategy }) => {
  const { t } = useTranslation('workbench-v2');
  return (
    <div className="grid gap-3 sm:grid-cols-2">
      <label className="text-xs text-slate-400">
        {t('step4.group.editor.name')}
        <input type="text" value={draft.name} disabled={readonly} onChange={(event) => onChange({ name: event.target.value })} className={FIELD} data-testid="group-name" />
      </label>
      <fieldset className="text-xs text-slate-400">
        <legend>{t('step4.group.editor.storage')}</legend>
        {(['custom', 'managed'] as WriteGroupStorageStrategy[]).map((strategy) => (
          <label key={strategy} className="mr-4 inline-flex items-center gap-1 text-slate-200">
            <input type="radio" name="group-storage" value={strategy} checked={draft.storage_strategy === strategy} disabled={readonly || locked} onChange={() => onStorageStrategy(strategy)} data-testid={`group-storage-${strategy}`} />
            {t(`step4.group.editor.storage_${strategy}`)}
          </label>
        ))}
        <p className="mt-1 text-slate-500" data-testid="group-storage-note">{t(`step4.group.editor.storage_${draft.storage_strategy}_note`)}</p>
      </fieldset>
      <label className="text-xs text-slate-400">
        {t('step4.group.editor.table')}
        <input type="text" value={draft.table_name} disabled={readonly || locked} onChange={(event) => onChange({ table_name: event.target.value })} className={FIELD} data-testid="group-table" />
        {draft.storage_strategy === 'managed' && <span className="mt-1 block text-[11px] text-slate-500">{t('step4.group.editor.table_managed_note')}</span>}
      </label>
      <label className="text-xs text-slate-400">
        {t('step4.group.editor.schema')}
        <input type="text" value={draft.table_schema} disabled={readonly || locked} onChange={(event) => onChange({ table_schema: event.target.value })} className={FIELD} data-testid="group-schema" />
      </label>
    </div>
  );
};
