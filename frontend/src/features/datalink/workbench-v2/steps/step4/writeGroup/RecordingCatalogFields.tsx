import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { ParameterHelp } from '../../../components/ParameterHelp';
import type { EditorDraft } from '../../../state/writeGroup/draft';

export interface RecordingCatalogView {
  status: 'loading' | 'ready' | 'error';
  defaultSchema: string;
  schemas: readonly { name: string; tables: readonly string[] }[];
  blockedTables?: readonly { schema: string; name: string }[];
  retry: () => void;
}

const FIELD = 'mt-1 w-full rounded border border-[var(--gw-control)] bg-[var(--gw-surface)] px-2 py-1 text-sm text-[var(--gw-text)]';
const manual = 'manual';
const listed = (value: string) => `listed:${value}`;
// Saved destinations normalize outer whitespace; do not offer a choice that would name a different table after saving.
const selectable = (value: string) => value === value.trim();

/** Catalog choices update the existing draft; free entry remains available for new or unlisted targets. */
export function RecordingCatalogFields({ draft, disabled, scopedPresentation, catalog, onChange }: {
  draft: EditorDraft; disabled: boolean; scopedPresentation: boolean; catalog: RecordingCatalogView;
  onChange: (patch: Partial<EditorDraft>) => void;
}) {
  const { t } = useTranslation('workbench-v2');
  const { t: text } = useTranslation('recording-editor');
  const [manualSchema, setManualSchema] = useState(false), [manualTable, setManualTable] = useState(false);
  const effectiveSchema = draft.table_schema.trim() || catalog.defaultSchema;
  const tables = catalog.schemas.find(schema => schema.name === effectiveSchema)?.tables ?? [];
  const owned = (name: string) => catalog.blockedTables?.some(table => table.schema === effectiveSchema && table.name === name) ?? false;
  const schemaEntry = manualSchema || Boolean(draft.table_schema && !catalog.schemas.some(schema => schema.name === draft.table_schema));
  const tableEntry = manualTable || Boolean(draft.table_name && !tables.includes(draft.table_name));
  const schemaChoice = schemaEntry ? manual : draft.table_schema ? listed(draft.table_schema) : '';
  const tableChoice = tableEntry ? manual : draft.table_name ? listed(draft.table_name) : '';
  return <>
    <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 text-xs text-[var(--gw-muted)]">
      {t('step4.group.editor.table')}
      <select value={tableChoice} disabled={disabled} className={FIELD} data-testid="group-table" onChange={event => {
        const value = event.target.value;
        if (value === manual) { setManualTable(true); return; }
        const selected = tables.find(table => listed(table) === value);
        if (value && (!selected || !selectable(selected) || owned(selected))) return;
        setManualTable(false); onChange({ table_name: selected ?? '' });
      }}>
        <option value="">{text(draft.storage_strategy === 'managed' ? 'newDestination.catalogAutoTable' : 'newDestination.catalogChooseTable')}</option>
        {tables.map(table => <option key={table} value={listed(table)} disabled={!selectable(table) || owned(table)}>{!selectable(table) ? `${JSON.stringify(table)} ${text('newDestination.catalogOuterSpace')}` : owned(table) ? `${table} ${text('newDestination.catalogOwned')}` : table}</option>)}
        <option value={manual}>{text('newDestination.catalogManualTable')}</option>
      </select>
      {tableEntry && <input type="text" aria-label={text('newDestination.catalogManualTable')} value={draft.table_name} disabled={disabled}
        onChange={event => onChange({ table_name: event.target.value })} className={FIELD} data-testid="group-table-manual" />}
      {tables.length === 0 && <span className="mt-1 block text-xs text-[var(--gw-muted)]">{text('newDestination.catalogNoTables')}</span>}
      {draft.storage_strategy === 'managed' && <span className="mt-1 block text-xs text-[var(--gw-muted)]">{t('step4.group.editor.table_managed_note')}</span>}
    </label><ParameterHelp compact={scopedPresentation} help="record.table" /></div>
    <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 text-xs text-[var(--gw-muted)]">
      {t('step4.group.editor.schema')}
      <select value={schemaChoice} disabled={disabled} className={FIELD} data-testid="group-schema" onChange={event => {
        const value = event.target.value;
        if (value === manual) { setManualSchema(true); return; }
        const selected = catalog.schemas.find(schema => listed(schema.name) === value);
        if (value && (!selected || !selectable(selected.name))) return;
        setManualSchema(false); onChange({ table_schema: selected?.name ?? '' });
      }}>
        <option value="">{text('newDestination.catalogDefaultSchema')}</option>
        {catalog.schemas.map(schema => <option key={schema.name} value={listed(schema.name)} disabled={!selectable(schema.name)}>{selectable(schema.name) ? schema.name : `${JSON.stringify(schema.name)} ${text('newDestination.catalogOuterSpace')}`}</option>)}
        <option value={manual}>{text('newDestination.catalogManualSchema')}</option>
      </select>
      {schemaEntry && <input type="text" aria-label={text('newDestination.catalogManualSchema')} value={draft.table_schema} disabled={disabled}
        onChange={event => onChange({ table_schema: event.target.value })} className={FIELD} data-testid="group-schema-manual" />}
    </label><ParameterHelp compact={scopedPresentation} help="record.schema" /></div>
  </>;
}
