import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { focusWorkspaceReadiness } from './step4FocusTargets';

/**
 * Says what Advanced and the workspace are for. Neither is another table format or a competing recording setting:
 * Advanced edits the same saved profile, and the workspace is where readiness and activation are reviewed.
 */
export const CircuitProfileRoles: React.FC<{ onOpenAdvanced?: (groupId?: string) => void }> = ({ onOpenAdvanced }) => {
  const { t } = useTranslation('workbench-v2');
  return (
    <section data-testid="circuit-profile-roles" aria-label={t('step4.basic.roles.title')} className="space-y-1.5 rounded border border-[var(--gw-border)] bg-[var(--gw-surface)] p-3 text-xs text-[var(--gw-muted)]">
      <h4 className="text-xs font-semibold text-[var(--gw-text)]">{t('step4.basic.roles.title')}</h4>
      <p data-testid="circuit-roles-advanced">{t('step4.basic.roles.advanced')}</p>
      <p data-testid="circuit-roles-workspace">{t('step4.basic.roles.workspace')}</p>
      {onOpenAdvanced && (
        <button type="button" data-testid="circuit-profile-open-workspace" onClick={() => { onOpenAdvanced(); focusWorkspaceReadiness(); }}
          className="rounded border border-[var(--gw-control)] px-3 py-1.5 text-[var(--gw-text)] hover:bg-[var(--gw-surface)]">
          {t('step4.basic.roles.open_workspace')}
        </button>
      )}
    </section>
  );
};
