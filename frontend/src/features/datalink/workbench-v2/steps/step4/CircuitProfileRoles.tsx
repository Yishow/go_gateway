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
    <section data-testid="circuit-profile-roles" aria-label={t('step4.basic.roles.title')} className="space-y-1.5 rounded border border-slate-800/80 bg-slate-950/20 p-3 text-xs text-slate-400">
      <h4 className="text-xs font-semibold text-slate-200">{t('step4.basic.roles.title')}</h4>
      <p data-testid="circuit-roles-advanced">{t('step4.basic.roles.advanced')}</p>
      <p data-testid="circuit-roles-workspace">{t('step4.basic.roles.workspace')}</p>
      {onOpenAdvanced && (
        <button type="button" data-testid="circuit-profile-open-workspace" onClick={() => { onOpenAdvanced(); focusWorkspaceReadiness(); }}
          className="rounded border border-slate-700 px-3 py-1.5 text-slate-300 hover:bg-slate-800">
          {t('step4.basic.roles.open_workspace')}
        </button>
      )}
    </section>
  );
};
