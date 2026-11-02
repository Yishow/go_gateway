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
  return <details data-testid="group-revision-details" className="text-xs text-[var(--gw-muted)]">
    <summary className="cursor-pointer"><span>{text('versions.details')}</span> · {home('savedStatus')}: {status}</summary>
    {scopeDiagnostics && <section className="mt-2 space-y-1" data-testid="locator-version-details">{scopeDiagnostics}</section>}
    <section className="mt-2">{values}</section>
  </details>;
}
