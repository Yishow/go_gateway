import * as React from 'react';
import { useTranslation } from 'react-i18next';

export interface WorkspaceScopeNoteProps {
  groupSummary: { state: 'loading' | 'error' | 'ready'; total: number; applied: number };
}

/**
 * What the workspace review below covers: the saved groups as they are saved. A saved group that is not applied is
 * a draft and records nothing; this review never changes an interval or a storage choice.
 */
export const WorkspaceScopeNote: React.FC<WorkspaceScopeNoteProps> = ({ groupSummary }) => {
  const { t } = useTranslation('workbench-v2');
  const draftCount = groupSummary.total - groupSummary.applied;
  return (
    <div data-testid="workspace-scope-note" data-total={groupSummary.total} data-applied={groupSummary.applied} className="space-y-1 text-xs text-slate-400">
      <p>{t('step4.workspace_scope.scope', { total: groupSummary.total, applied: groupSummary.applied })}</p>
      {draftCount > 0 && <p className="text-amber-300">{t('step4.workspace_scope.draft_note', { count: draftCount })}</p>}
      <p>{t('step4.workspace_scope.no_changes')}</p>
    </div>
  );
};
