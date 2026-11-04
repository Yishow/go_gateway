import { ParameterHelp } from '../../../components/ParameterHelp';
import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroupStorageStrategy } from '../../../../../../types/studioV2WriteGroup';
import type { EditorDraft } from '../../../state/writeGroup/draft';
import { RecordingCatalogFields, type RecordingCatalogView } from './RecordingCatalogFields';

const FIELD = 'mt-1 w-full rounded border border-[var(--gw-control)] bg-[var(--gw-surface)] px-2 py-1 text-sm text-[var(--gw-text)]';

export interface GroupIdentityFieldsProps {
  scopedPresentation?: boolean;
  section?: 'all' | 'name' | 'destination';
  draft: EditorDraft;
  readonly: boolean;
  /** A profile group keeps the layout, storage kind and shared table it was created with. */
  locked: boolean;
  catalog?: RecordingCatalogView;
  onChange: (patch: Partial<EditorDraft>) => void;
  onStorageStrategy: (strategy: WriteGroupStorageStrategy) => void;
}

/** Name, storage kind and destination table of a group. */
export const GroupIdentityFields: React.FC<GroupIdentityFieldsProps> = ({ scopedPresentation = false, section = 'all', draft, readonly, locked, catalog, onChange, onStorageStrategy }) => {
  const { t } = useTranslation('workbench-v2');
  const { t: text } = useTranslation('recording-editor');
  return (
    <div className={section === 'name' ? 'grid gap-3' : section === 'destination' ? 'recording-destination-fields grid gap-3 sm:grid-cols-2' : 'grid gap-3 sm:grid-cols-2'}>
      {section !== 'destination' && <div className={section === 'name' ? 'recording-simple-name min-w-0' : 'flex min-w-0 items-start gap-2'}><label className="min-w-0 flex-1 text-xs text-[var(--gw-muted)]">
        <span>{section === 'name' ? text('simple.summary.name') : t('step4.group.editor.name')}</span>
        <input type="text" value={draft.name} disabled={readonly} onChange={(event) => onChange({ name: event.target.value })} className={FIELD} data-testid="group-name" />
      </label><ParameterHelp compact={scopedPresentation || section === 'name'} help="record.name" /></div>}
      {section !== 'name' && <>
      <fieldset className="text-xs text-[var(--gw-muted)]">
        <legend>{t('step4.group.editor.storage')} <ParameterHelp compact={scopedPresentation} help="record.strategy" /></legend>
        {(['custom', 'managed'] as WriteGroupStorageStrategy[]).map((strategy) => (
          <label key={strategy} className="mr-4 inline-flex items-center gap-1 text-[var(--gw-text)]">
            <input type="radio" name="group-storage" value={strategy} checked={draft.storage_strategy === strategy} disabled={readonly || locked} onChange={() => onStorageStrategy(strategy)} data-testid={`group-storage-${strategy}`} />
            {t(`step4.group.editor.storage_${strategy}`)}
          </label>
        ))}
        <p className="mt-1 text-[var(--gw-muted)]" data-testid="group-storage-note">{t(`step4.group.editor.storage_${draft.storage_strategy}_note`)}</p>
      </fieldset>
      {catalog?.status === 'ready' && !(scopedPresentation && locked && draft.storage_strategy === 'managed') ? <RecordingCatalogFields key={JSON.stringify([draft.connector_id, draft.connector_revision])} draft={draft} disabled={readonly || locked} scopedPresentation={scopedPresentation} catalog={catalog} onChange={onChange} /> : <><div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 text-xs text-[var(--gw-muted)]">
        {t('step4.group.editor.table')}
        {scopedPresentation && locked && draft.storage_strategy === 'managed' ? <details data-testid="group-managed-table"><summary>{text('compact.systemManaged')}</summary><code>{draft.table_name}</code></details> : <input type="text" value={draft.table_name} disabled={readonly || locked} onChange={(event) => onChange({ table_name: event.target.value })} className={FIELD} data-testid="group-table" />}
        {draft.storage_strategy === 'managed' && <span className="mt-1 block text-xs text-[var(--gw-muted)]">{t('step4.group.editor.table_managed_note')}</span>}
      </label><ParameterHelp compact={scopedPresentation} help="record.table" /></div>
      <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 text-xs text-[var(--gw-muted)]">
        {t('step4.group.editor.schema')}
        <input type="text" value={draft.table_schema} disabled={readonly || locked} onChange={(event) => onChange({ table_schema: event.target.value })} className={FIELD} data-testid="group-schema" />
      </label><ParameterHelp compact={scopedPresentation} help="record.schema" /></div></>}
      {catalog && catalog.status !== 'ready' && <div role={catalog.status === 'error' ? 'alert' : 'status'} className="text-xs text-[var(--gw-muted)]">{text(`newDestination.catalog${catalog.status === 'error' ? 'Error' : 'Loading'}`)}
        {catalog.status === 'error' && <button type="button" className="ml-2 underline" onClick={catalog.retry}>{text('newDestination.catalogRetry')}</button>}</div>}
      </>}
    </div>
  );
};
