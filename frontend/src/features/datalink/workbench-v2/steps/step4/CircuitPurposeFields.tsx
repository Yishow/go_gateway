import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { CircuitPurpose } from '../../state/circuitProfileBasic';

export interface CircuitPurposeFieldsProps {
  purpose: CircuitPurpose;
  /** A saved circuit keeps the layout it was created with; another layout needs the reviewed transition, not this choice. */
  locked: boolean;
  readonly: boolean;
  onChange: (purpose: CircuitPurpose) => void;
}

const PURPOSES: CircuitPurpose[] = ['wide', 'long'];

/**
 * The one choice that decides the table layout: a whole-device snapshot (one row per time window) or selected
 * measurements (one row per source per window). Both can record instantaneous and cumulative profiles.
 */
export const CircuitPurposeFields: React.FC<CircuitPurposeFieldsProps> = ({ purpose, locked, readonly, onChange }) => {
  const { t } = useTranslation('workbench-v2');
  return (
    <fieldset data-testid="circuit-profile-purpose" className="space-y-1.5 text-xs text-slate-300">
      <legend className="font-medium text-slate-200">{t('step4.basic.circuit.purpose_title')}</legend>
      <div className="grid gap-2 sm:grid-cols-2">
        {PURPOSES.map((item) => (
          <label key={item} className="block space-y-0.5 rounded border border-slate-800 p-2">
            <span className="flex items-center gap-2 font-medium text-slate-200">
              <input type="radio" name="circuit-purpose" data-testid={`circuit-profile-purpose-${item}`} checked={purpose === item}
                disabled={readonly || locked} onChange={() => onChange(item)} />
              {t(`step4.basic.circuit.purpose_${item}`)}
            </span>
            <span data-testid={`circuit-profile-purpose-hint-${item}`} className="block text-[11px] text-slate-400">{t(`step4.basic.circuit.purpose_${item}_rows`)}</span>
          </label>
        ))}
      </div>
      {locked && <p data-testid="circuit-profile-purpose-locked" className="text-[11px] text-slate-400">{t('step4.basic.circuit.purpose_locked')}</p>}
    </fieldset>
  );
};
