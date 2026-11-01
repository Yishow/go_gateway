import * as React from 'react';
import { useTranslation } from 'react-i18next';

/** The explicit way from a recording legacy wide group to preparing long targets; it changes nothing by itself. */
export const LegacyMigrationEntry: React.FC<{ disabled: boolean; onOpen: () => void }> = ({ disabled, onOpen }) => {
  const { t } = useTranslation('workbench-v2');
  return (
    <div data-testid="legacy-migration-entry" className="rounded border border-[var(--gw-border)] p-3 text-xs text-[var(--gw-text)]">
      <p className="text-[var(--gw-muted)]">{t('step4.basic.migration.hint')}</p>
      <button type="button" data-testid="legacy-migration-open" disabled={disabled} onClick={onOpen} className="mt-2 rounded border border-[var(--gw-focus)] px-3 py-2 text-[var(--gw-accent)] disabled:text-[var(--gw-muted)] disabled:bg-[var(--gw-ground)]">
        {t('step4.basic.migration.open')}
      </button>
    </div>
  );
};
