import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroup } from '../../../../types/studioV2WriteGroup';
/** Saved configuration and locator readback evidence only; never a health or write acknowledgement. */
export function RecordingRevisionDetails({ group, simple, scopeDiagnostics }: { group: WriteGroup; simple: boolean; scopeDiagnostics?: ReactNode }) {
  const { t } = useTranslation('workbench-v2');
  const { t: text } = useTranslation('recording-editor');
  const { t: home } = useTranslation('recording-home');
  const status = simple ? home(`savedStatuses.${group.status}`) : group.status;
  const values = <dl className="grid grid-cols-2 gap-x-4 gap-y-1 text-xs text-[var(--gw-muted)] sm:grid-cols-4" data-testid="group-revisions">
    <dt>{simple ? home('savedStatus') : t('step4.group.editor.status')}</dt><dd className="text-[var(--gw-text)]">{status}</dd>
    <dt>{t('step4.group.editor.draft_revision')}</dt><dd className="font-mono text-[var(--gw-text)]">{group.revision}</dd>
    <dt>{t('step4.group.editor.applied_revision')}</dt><dd className="font-mono text-[var(--gw-text)]" data-testid="group-applied-revision">{group.applied_revision || t('step4.group.editor.not_applied')}</dd>
    <dt>{t('step4.group.editor.connector_revision')}</dt><dd className="font-mono text-[var(--gw-text)]">{group.destination.connector_revision}</dd>
  </dl>;
  if (!simple) return values;
  // Side rail: one row per value, IDs shortened with the full value kept in the tooltip and selectable text.
  const id = (value: string | undefined, empty: string) => value ? <span className="font-mono" title={value}>{value.length > 12 ? `${value.slice(0, 8)}…` : value}</span> : empty;
  return <section data-testid="group-revision-details" className="recording-revision-details" aria-label={text('versions.details')}>
    <h4>{text('versions.details')}</h4>
    {scopeDiagnostics && <div className="space-y-1" data-testid="locator-version-details">{scopeDiagnostics}</div>}
    <dl data-testid="group-revisions">
      <div><dt>{home('savedStatus')}</dt><dd>{status}</dd></div>
      <div><dt>{t('step4.group.editor.draft_revision')}</dt><dd>{id(group.revision, '—')}</dd></div>
      <div><dt>{t('step4.group.editor.applied_revision')}</dt><dd data-testid="group-applied-revision">{id(group.applied_revision, t('step4.group.editor.not_applied'))}</dd></div>
      <div><dt>{t('step4.group.editor.connector_revision')}</dt><dd>{id(group.destination.connector_revision, '—')}</dd></div>
    </dl>
  </section>;
}
