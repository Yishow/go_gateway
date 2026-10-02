import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroupReadiness } from '../../../../../../types/studioV2WriteGroup';

export interface GroupReadinessPanelProps {
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
      <span className={value ? 'text-emerald-300' : 'text-amber-300'}>{value ? t('step4.group.readiness.yes') : t('step4.group.readiness.no')}</span>
    </li>
  );
}

/** Shows the backend's persisted-readiness answer for the saved group; the editor never computes it. */
export const GroupReadinessPanel: React.FC<GroupReadinessPanelProps> = ({ readiness, loading, failed, dirty, onRetry }) => {
  const { t } = useTranslation('workbench-v2');
  if (loading) return <p className="text-xs text-slate-400" role="status" data-testid="group-readiness-loading">{t('step4.group.readiness.loading')}</p>;
  if (failed) {
    return (
      <div role="alert" className="text-xs text-amber-300" data-testid="group-readiness-failed">
        {t('step4.group.readiness.failed')}
        <button type="button" onClick={onRetry} className="ml-2 underline" data-testid="group-readiness-retry">{t('step4.group.retry')}</button>
      </div>
    );
  }
  if (!readiness) return <p className="text-xs text-slate-500" data-testid="group-readiness-none">{t('step4.group.readiness.none')}</p>;
  return (
    <div className="space-y-2 text-xs text-slate-300" data-testid="group-readiness">
      {dirty && <p className="text-amber-300" data-testid="group-readiness-stale">{t('step4.group.readiness.stale_draft')}</p>}
      <ul className="space-y-1">
        <Fact label={t('step4.group.readiness.config')} value={readiness.config_ready} />
        <Fact label={t('step4.group.readiness.schema')} value={readiness.schema_ready} />
        <Fact label={t('step4.group.readiness.ready')} value={readiness.ready} />
      </ul>
      {readiness.issues.length > 0 && (
        <ul className="space-y-1" data-testid="group-readiness-issues">
          {readiness.issues.map((issue) => (
            <li key={`${issue.code}-${issue.scope}`} className={issue.severity === 'blocking' ? 'text-red-300' : 'text-amber-300'}>
              <span className="font-semibold">{issue.severity === 'blocking' ? t('step4.group.readiness.blocking') : t('step4.group.readiness.warning')}</span>
              {' · '}{issue.message}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};
