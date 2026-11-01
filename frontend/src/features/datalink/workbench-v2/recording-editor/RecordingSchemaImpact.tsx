import { useTranslation } from 'react-i18next';
import type { SchemaPreviewToken } from '../../../../types/recordingPlan';
/** Exact service impact, including auxiliary tables. No inferred DDL or readiness. */
export function RecordingSchemaImpact({ preview, layout }: { preview: SchemaPreviewToken; layout?: string }) {
  const { t } = useTranslation('workbench-v2');
  const { t: editor } = useTranslation('recording-editor');
  const columns = preview.group_layout?.columns.filter(column => column.name !== preview.group_layout?.owner_column) ?? [];
  const addedColumns = new Set(preview.tables.filter(table => table.action === 'add_columns').flatMap(table => table.columns ?? []));
  return (
        <div className="max-h-72 space-y-3 overflow-auto rounded border border-[var(--gw-border)] bg-[var(--gw-surface)] p-3 text-xs text-[var(--gw-text)]" data-testid="group-schema-preview-result">
          <p role="status">{editor('preparation.impact.previewOnly')}</p>
          <dl className="grid gap-x-3 gap-y-1 sm:grid-cols-[auto_minmax(0,1fr)]" data-testid="recording-schema-impact-scope">
            <dt>{editor('preparation.impact.database')}</dt><dd className="break-words">{preview.database ?? '—'}</dd>
            <dt>{editor('preparation.impact.schema')}</dt><dd>{preview.schema || '—'}</dd>
            <dt>{editor('preparation.impact.layout')}</dt><dd>{layout || '—'}</dd>
            <dt>{editor('preparation.impact.workspaceRevision')}</dt><dd className="break-all">{preview.workspace_revision}</dd>
            <dt>{editor('preparation.impact.recordingRevision')}</dt><dd className="break-all">{preview.plan_revision}</dd>
            <dt>{editor('preparation.impact.connectorRevision')}</dt><dd className="break-all">{preview.connector_revision ?? '—'}</dd>
          </dl>
          <ul className="space-y-2" data-testid="recording-schema-impact-tables">{preview.tables.map((table, index) => <li key={`${table.name}:${index}`}>
            <strong className="break-words">{table.name}</strong> · {editor(`preparation.impact.action.${table.action}`)}
            <p className="break-words">{table.columns.join(', ') || editor('preparation.impact.noColumns')}</p>
          </li>)}</ul>
          {Boolean(preview.group_layout?.indexes?.length) && <div><p>{editor('preparation.impact.indexes')}</p>
            <ul>{preview.group_layout?.indexes?.map(index => <li key={index.name} className="break-words">{index.name}: {index.columns.join(', ')}</li>)}</ul>
          </div>}
          <div className="grid gap-1 sm:grid-cols-2" data-testid="group-schema-preview-table">
            <span className="text-[var(--gw-muted)]">{t('step4.group.schema.table')}</span>
            <span className="font-mono">{preview.group_layout?.table_name ?? preview.table_prefix}</span>
            {preview.no_change_reason && <><span className="text-[var(--gw-muted)]">{t('step4.group.schema.change')}</span><span>{t('step4.group.schema.no_change')}</span></>}
          </div>

          {columns.length > 0 && (
            <div className="overflow-x-auto" role="region" aria-label={t('step4.group.schema.columns_region')} tabIndex={0}>
              <table className="w-full min-w-[520px] border-collapse text-left" data-testid="group-schema-columns">
                <thead className="text-[var(--gw-muted)]">
                  <tr>
                    <th scope="col" className="p-2">{t('step4.group.schema.column_name')}</th>
                    <th scope="col" className="p-2">{t('step4.group.schema.sql_type')}</th>
                    <th scope="col" className="p-2">{t('step4.group.schema.nullable')}</th>
                    <th scope="col" className="p-2">{t('step4.group.schema.primary_key')}</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-[var(--gw-border)]">
                  {columns.map((column) => (
                    <tr key={column.name}>
                      <td className="p-2 font-mono">{column.name} {addedColumns.has(column.name) && <span className="font-semibold text-[var(--gw-text)]">{t('step4.group.schema.added_column')}</span>}</td>
                      <td className="p-2 font-mono">{column.sql_type}</td>
                      <td className="p-2">{column.nullable ? t('step4.group.schema.yes') : t('step4.group.schema.no')}</td>
                      <td className="p-2">{column.primary_key ? t('step4.group.schema.yes') : t('step4.group.schema.no')}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}

          <details>
            <summary className="cursor-pointer text-[var(--gw-muted)]">{t('step4.group.schema.advanced_sql')}</summary>
            <pre className="mt-2 max-w-full overflow-x-auto whitespace-pre-wrap text-[11px] text-[var(--gw-muted)]">{preview.statements.join('\n')}</pre>
          </details>
          <p className="text-[var(--gw-muted)]" role="status">{t('step4.group.schema.confirm_note')}</p>
        </div>
  );
}
