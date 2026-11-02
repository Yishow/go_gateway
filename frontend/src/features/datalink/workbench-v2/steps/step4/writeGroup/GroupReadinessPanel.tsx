import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroupReadiness } from '../../../../../../types/studioV2WriteGroup';

export interface GroupReadinessPanelProps {
  compact?: boolean;
  readiness: WriteGroupReadiness | undefined;
  loading: boolean;
  failed: boolean;
  /** The editor holds changes the server has not seen, so its readiness describes older content. */
  dirty: boolean;
  onRetry: () => void;
}

function Fact({ label, value }: { label: string; value: boolean }) {
  const { t } = useTranslation('workbench-v2');
  return (
    <li className="flex items-center gap-2">
      <span aria-hidden="true">{value ? '✓' : '✗'}</span>
      <span>{label}</span>
      <span className={value ? 'text-[var(--gw-success-fg)]' : 'text-[var(--gw-warning-fg)]'}>{value ? t('step4.group.readiness.yes') : t('step4.group.readiness.no')}</span>
    </li>
  );
}

/** Shows the backend's persisted-readiness answer for the saved group; the editor never computes it. */
export const GroupReadinessPanel: React.FC<GroupReadinessPanelProps> = ({ readiness, loading, failed, dirty, onRetry, compact = false }) => {
  const { t } = useTranslation('workbench-v2');
  const { t: editorText } = useTranslation('recording-editor');
  if (loading) return <p className="text-xs text-[var(--gw-muted)]" role="status" data-testid="group-readiness-loading">{t('step4.group.readiness.loading')}</p>;
  if (failed) {
    return (
      <div role="alert" className="text-xs text-[var(--gw-warning-fg)]" data-testid="group-readiness-failed">
        {t('step4.group.readiness.failed')}
        <button type="button" onClick={onRetry} className="ml-2 underline" data-testid="group-readiness-retry">{t('step4.group.retry')}</button>
      </div>
    );
  }
  if (!readiness) return <p className="text-xs text-[var(--gw-muted)]" data-testid="group-readiness-none">{t('step4.group.readiness.none')}</p>;
  const facts = (
    <ul className="space-y-1">
        <Fact label={t('step4.group.readiness.config')} value={readiness.config_ready} />
        <Fact label={t('step4.group.readiness.schema')} value={readiness.schema_ready} />
        <Fact label={t('step4.group.readiness.ready')} value={readiness.ready} />
      </ul>
  );
  return (
    <div className="space-y-2 text-xs text-[var(--gw-text)]" data-testid="group-readiness">
      {dirty && <p className="text-[var(--gw-warning-fg)]" data-testid="group-readiness-stale">{t('step4.group.readiness.stale_draft')}</p>}
      {compact && !dirty && (readiness.ready || readiness.issues.length > 0) ? <details data-testid="recording-readiness-details"><summary>{editorText('workflow.readinessDetails')}</summary>{facts}</details> : facts}
      {readiness.issues.length > 0 && (
        <ul className="space-y-1" data-testid="group-readiness-issues">
          {readiness.issues.map((issue) => (
            <li key={`${issue.code}-${issue.scope}`} className={issue.severity === 'blocking' ? 'text-[var(--gw-danger-fg)]' : 'text-[var(--gw-warning-fg)]'}>
              <span className="font-semibold">{issue.severity === 'blocking' ? t('step4.group.readiness.blocking') : t('step4.group.readiness.warning')}</span>
              {' · '}{compact ? editorText(issue.code === 'schema-confirmation-required' ? 'confirmation.schemaRequired' : 'confirmation.unknownIssue') : issue.message}
              {compact && <details className="mt-1 text-[var(--gw-muted)]"><summary>{t('step4.basic.diagnostics')}</summary><p>{issue.code}: {issue.message}</p></details>}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};
