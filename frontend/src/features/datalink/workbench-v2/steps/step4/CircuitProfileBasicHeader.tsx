import * as React from 'react';
import { useTranslation } from 'react-i18next';

export const CircuitProfileBasicHeader: React.FC<{ prepareOnly: boolean }> = ({ prepareOnly }) => {
  const { t } = useTranslation('workbench-v2');
  const scope = prepareOnly ? 'step4.basic.migration' : 'step4.basic.circuit';
  return (
    <div>
      <h3 id="circuit-profile-basic-title" className="text-sm font-semibold text-slate-100">{t(`${scope}.title`)}</h3>
      <p className="mt-1 text-xs leading-5 text-slate-400">{t(`${scope}.subtitle`)}</p>
      {prepareOnly && <p role="note" data-testid="circuit-profile-prepare-only" className="mt-2 rounded border border-amber-700/60 bg-amber-950/30 p-2 text-xs text-amber-200">{t('step4.basic.migration.banner')}</p>}
    </div>
  );
};
