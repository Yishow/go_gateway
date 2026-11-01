import { ParameterHelp } from '../../components/ParameterHelp';
import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { CircuitPurpose } from '../../state/circuitProfileBasic';
import type { CircuitLayoutBlockReason, CircuitLayoutOptions } from '../../state/circuitProfileCapabilities';
import { focusDestinationSettings } from './step4FocusTargets';

export interface CircuitPurposeFieldsProps {
  purpose: CircuitPurpose;
  /** A saved circuit keeps the layout it was created with; another layout needs the reviewed transition, not this choice. */
  locked: boolean;
  readonly: boolean;
  onChange: (purpose: CircuitPurpose) => void;
  layoutOptions?: CircuitLayoutOptions;
  blockedReason?: CircuitLayoutBlockReason;
}

const PURPOSES: CircuitPurpose[] = ['wide', 'long'];
const UNKNOWN_LAYOUTS: CircuitLayoutOptions = { wide: 'capability_unavailable', long: 'capability_unavailable' };

/**
 * The one choice that decides the table layout: a whole-device snapshot (one row per time window) or selected
 * measurements (one row per source per window). Both can record instantaneous and cumulative profiles.
 */
export const CircuitPurposeFields: React.FC<CircuitPurposeFieldsProps> = ({ purpose, locked, readonly, onChange, layoutOptions = UNKNOWN_LAYOUTS, blockedReason }) => {
  const { t } = useTranslation('workbench-v2');
  const reasonId = React.useId();
  return (
    <fieldset data-testid="circuit-profile-purpose" className="space-y-1.5 text-xs text-[var(--gw-text)]">
      <legend className="font-medium text-[var(--gw-text)]">{t('step4.basic.circuit.purpose_title')} <ParameterHelp help="record.layout" /></legend>
      <div className="grid gap-2 sm:grid-cols-2">
        {PURPOSES.map((item) => (
          <label key={item} className="block space-y-0.5 rounded border border-[var(--gw-control)] p-2">
            <span className="flex items-center gap-2 font-medium text-[var(--gw-text)]">
              <input type="radio" name="circuit-purpose" data-testid={`circuit-profile-purpose-${item}`} checked={purpose === item}
                disabled={readonly || locked || Boolean(layoutOptions[item])}
                aria-describedby={layoutOptions[item] ? `${reasonId}-${item}` : undefined}
                onChange={() => { if (!layoutOptions[item]) onChange(item); }} />
              {t(`step4.basic.circuit.purpose_${item}`)}
            </span>
            <span data-testid={`circuit-profile-purpose-hint-${item}`} className="block text-xs text-[var(--gw-muted)]">{t(`step4.basic.circuit.purpose_${item}_rows`)}</span>
            {layoutOptions[item] && <span id={`${reasonId}-${item}`} data-testid={`circuit-profile-purpose-reason-${item}`}
              className="block text-xs text-[var(--gw-warning-fg)]">{t(`step4.basic.circuit.${layoutOptions[item]}`)}</span>}
          </label>
        ))}
      </div>
      <p className="text-xs text-[var(--gw-muted)]">{t('step4.basic.circuit.layout_not_readiness')}</p>
      {blockedReason && <p role="status" data-testid="circuit-profile-layout-blocked" className="text-[var(--gw-warning-fg)]">
        {t(`step4.basic.circuit.${blockedReason}`)}{' '}
        <button type="button" onClick={focusDestinationSettings} className="underline">{t('step4.basic.blockers.repair_destination')}</button>
      </p>}
      {locked && <p data-testid="circuit-profile-purpose-locked" className="text-xs text-[var(--gw-muted)]">{t('step4.basic.circuit.purpose_locked')}</p>}
    </fieldset>
  );
};
