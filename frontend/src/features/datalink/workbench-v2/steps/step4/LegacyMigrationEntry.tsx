import * as React from 'react';
import { useTranslation } from 'react-i18next';

/** The explicit way from a recording legacy wide group to preparing long targets; it changes nothing by itself. */
export const LegacyMigrationEntry: React.FC<{ disabled: boolean; onOpen: () => void }> = ({ disabled, onOpen }) => {
  const { t } = useTranslation('workbench-v2');
  return (
    <div data-testid="legacy-migration-entry" className="rounded border border-slate-800 p-3 text-xs text-slate-300">
      <p className="text-slate-400">{t('step4.basic.migration.hint')}</p>
      <button type="button" data-testid="legacy-migration-open" disabled={disabled} onClick={onOpen} className="mt-2 rounded border border-cyan-700 px-3 py-2 text-cyan-200 disabled:opacity-40">
        {t('step4.basic.migration.open')}
      </button>
    </div>
  );
};
