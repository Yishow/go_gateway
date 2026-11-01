import * as React from 'react';
import { useTranslation } from 'react-i18next';

export const CircuitProfileBasicHeader: React.FC<{ prepareOnly: boolean }> = ({ prepareOnly }) => {
  const { t } = useTranslation('workbench-v2');
  const scope = prepareOnly ? 'step4.basic.migration' : 'step4.basic.circuit';
  return (
    <div>
      <h3 id="circuit-profile-basic-title" className="text-sm font-semibold text-[var(--gw-text)]">{t(`${scope}.title`)}</h3>
      <p className="mt-1 text-xs leading-5 text-[var(--gw-muted)]">{t(`${scope}.subtitle`)}</p>
      {prepareOnly && <p role="note" data-testid="circuit-profile-prepare-only" className="mt-2 rounded border border-[var(--gw-warning-border)] bg-[var(--gw-warning-bg)] p-2 text-xs text-[var(--gw-warning-fg)]">{t('step4.basic.migration.banner')}</p>}
    </div>
  );
};
