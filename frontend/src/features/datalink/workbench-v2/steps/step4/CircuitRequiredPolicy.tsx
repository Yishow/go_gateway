import { ParameterHelp } from '../../components/ParameterHelp';
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
    <section data-testid="circuit-profile-required-policy" aria-label={t('step4.basic.required.title')} className="space-y-1 rounded border border-[var(--gw-control)] bg-[var(--gw-surface)] p-3 text-xs text-[var(--gw-text)]">
      <h4 className="text-xs font-semibold text-[var(--gw-text)]">{t('step4.basic.required.title')} <ParameterHelp help="member.required" /></h4>
      <p>{t('step4.basic.required.count', { count })}</p>
      <p className="text-[var(--gw-muted)]">{t(count > 0 ? 'step4.basic.required.wide_rule' : 'step4.basic.required.wide_optional')}</p>
      {onOpen && (
        <button type="button" data-testid="circuit-profile-required-open" onClick={onOpen} className="rounded border border-[var(--gw-control)] px-2 py-0.5 text-xs text-[var(--gw-text)] hover:bg-[var(--gw-raised)]">
          {t('step4.basic.required.open')}
        </button>
      )}
    </section>
  );
};
