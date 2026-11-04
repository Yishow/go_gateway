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

function Fact({ label, value, compact = false }: { label: string; value: boolean; compact?: boolean }) {
  const { t } = useTranslation('workbench-v2');
  if (compact) return <li className="recording-readiness-check" data-state={value ? 'pass' : 'fail'}>
    <span>{label}</span><span className="recording-readiness-check__value"><span aria-hidden="true">{value ? '✓' : '✗'}</span> {value ? t('step4.group.readiness.yes') : t('step4.group.readiness.no')}</span>
  </li>;
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
  const sourceStale = compact && readiness.issues.some(issue => issue.code === 'source-revision-stale' && issue.severity === 'blocking');
  const visibleIssues = sourceStale ? readiness.issues.filter(issue => issue.code !== 'schema-check-blocked') : readiness.issues;
  const facts = (
    <ul className={compact ? 'recording-readiness-checks' : 'space-y-1'} data-testid={compact ? 'recording-readiness-details' : undefined}>
        <Fact compact={compact} label={t('step4.group.readiness.config')} value={readiness.config_ready} />
        <Fact compact={compact} label={t('step4.group.readiness.schema')} value={readiness.schema_ready} />
        <Fact compact={compact} label={t('step4.group.readiness.ready')} value={readiness.ready} />
      </ul>
  );
  return (
    <div className="space-y-2 text-xs text-[var(--gw-text)]" data-testid="group-readiness">
      {dirty && <p className="text-[var(--gw-warning-fg)]" data-testid="group-readiness-stale">{t('step4.group.readiness.stale_draft')}</p>}
      {!compact && facts}
      {readiness.issues.length > 0 && (
        <ul className={compact ? 'recording-readiness-issues' : 'space-y-1'} data-testid="group-readiness-issues">
          {visibleIssues.map((issue) => {
            const severity = issue.severity === 'blocking' ? t('step4.group.readiness.blocking') : t('step4.group.readiness.warning');
            if (!compact) return <li key={`${issue.code}-${issue.scope}`} className={issue.severity === 'blocking' ? 'text-[var(--gw-danger-fg)]' : 'text-[var(--gw-warning-fg)]'}>
              <span className="font-semibold">{severity}</span>{' · '}{issue.message}
            </li>;
            const diagnostics = [issue, ...(sourceStale && issue.code === 'source-revision-stale' ? readiness.issues.filter(item => item.code === 'schema-check-blocked') : [])];
            // Plain-language instruction first; the server's code and message stay visible as a quiet diagnostic line.
            return <li key={`${issue.code}-${issue.scope}`} className="recording-readiness-issue" data-severity={issue.severity}>
              <span className="recording-readiness-issue__pill">{severity}</span>
              <span className="recording-readiness-issue__text">{editorText(issue.code === 'schema-confirmation-required' ? 'confirmation.schemaRequired' : issue.code === 'source-revision-stale' ? 'confirmation.sourceStale' : issue.code === 'schema-check-blocked' ? 'confirmation.configBlocked' : 'confirmation.unknownIssue')}</span>
              {diagnostics.map(item => <span key={`${item.code}-${item.scope}`} className="recording-readiness-issue__code">{t('step4.basic.diagnostics')}：{item.code} · {item.message}</span>)}
            </li>;
          })}
        </ul>
      )}
      {compact && facts}
    </div>
  );
};
