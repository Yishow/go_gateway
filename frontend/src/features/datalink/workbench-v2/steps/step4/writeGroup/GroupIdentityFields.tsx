import { ParameterHelp } from '../../../components/ParameterHelp';
import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroupStorageStrategy } from '../../../../../../types/studioV2WriteGroup';
import type { EditorDraft } from '../../../state/writeGroup/draft';

const FIELD = 'mt-1 w-full rounded border border-[var(--gw-control)] bg-[var(--gw-surface)] px-2 py-1 text-sm text-[var(--gw-text)]';

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
      <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 text-xs text-[var(--gw-muted)]">
        {t('step4.group.editor.name')}
        <input type="text" value={draft.name} disabled={readonly} onChange={(event) => onChange({ name: event.target.value })} className={FIELD} data-testid="group-name" />
      </label><ParameterHelp help="record.name" /></div>
      <fieldset className="text-xs text-[var(--gw-muted)]">
        <legend>{t('step4.group.editor.storage')} <ParameterHelp help="record.strategy" /></legend>
        {(['custom', 'managed'] as WriteGroupStorageStrategy[]).map((strategy) => (
          <label key={strategy} className="mr-4 inline-flex items-center gap-1 text-[var(--gw-text)]">
            <input type="radio" name="group-storage" value={strategy} checked={draft.storage_strategy === strategy} disabled={readonly || locked} onChange={() => onStorageStrategy(strategy)} data-testid={`group-storage-${strategy}`} />
            {t(`step4.group.editor.storage_${strategy}`)}
          </label>
        ))}
        <p className="mt-1 text-[var(--gw-muted)]" data-testid="group-storage-note">{t(`step4.group.editor.storage_${draft.storage_strategy}_note`)}</p>
      </fieldset>
      <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 text-xs text-[var(--gw-muted)]">
        {t('step4.group.editor.table')}
        <input type="text" value={draft.table_name} disabled={readonly || locked} onChange={(event) => onChange({ table_name: event.target.value })} className={FIELD} data-testid="group-table" />
        {draft.storage_strategy === 'managed' && <span className="mt-1 block text-xs text-[var(--gw-muted)]">{t('step4.group.editor.table_managed_note')}</span>}
      </label><ParameterHelp help="record.table" /></div>
      <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 text-xs text-[var(--gw-muted)]">
        {t('step4.group.editor.schema')}
        <input type="text" value={draft.table_schema} disabled={readonly || locked} onChange={(event) => onChange({ table_schema: event.target.value })} className={FIELD} data-testid="group-schema" />
      </label><ParameterHelp help="record.schema" /></div>
    </div>
  );
};
