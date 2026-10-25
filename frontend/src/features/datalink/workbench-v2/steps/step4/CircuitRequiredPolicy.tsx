import * as React from 'react';
import { useTranslation } from 'react-i18next';

export interface CircuitRequiredPolicyProps {
  /** How many sources of the saved wide group are required, as saved. */
  count: number;
  /** Opens the saved group in Advanced, where the existing per source controls are; nothing is changed here. */
  onOpen?: () => void;
}

/**
 * Says how a whole wide row treats a missing source: an optional one is stored as NULL with its reason, a required
 * one still skips the row. It reads the saved state and never edits a member.
 */
export const CircuitRequiredPolicy: React.FC<CircuitRequiredPolicyProps> = ({ count, onOpen }) => {
  const { t } = useTranslation('workbench-v2');
  return (
    <section data-testid="circuit-profile-required-policy" aria-label={t('step4.basic.required.title')} className="space-y-1 rounded border border-slate-800/80 bg-slate-950/20 p-3 text-xs text-slate-300">
      <h4 className="text-xs font-semibold text-slate-200">{t('step4.basic.required.title')}</h4>
      <p>{t('step4.basic.required.count', { count })}</p>
      <p className="text-slate-400">{t(count > 0 ? 'step4.basic.required.wide_rule' : 'step4.basic.required.wide_optional')}</p>
      {onOpen && (
        <button type="button" data-testid="circuit-profile-required-open" onClick={onOpen} className="rounded border border-slate-700 px-2 py-0.5 text-[11px] text-slate-300 hover:bg-slate-800">
          {t('step4.basic.required.open')}
        </button>
      )}
    </section>
  );
};
