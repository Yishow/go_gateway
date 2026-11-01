import { ParameterHelp } from '../../components/ParameterHelp';
import { useId } from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroupValueRole } from '../../../../../types/studioV2WriteGroup';

interface Props {
  testId: string;
  label: string;
  value: WriteGroupValueRole;
  numeric: boolean;
  readonly: boolean;
  saved: boolean;
  onChange: (role: WriteGroupValueRole) => void;
}

export function CircuitProfileRoleOptions({ testId, label, value, numeric, readonly, saved, onChange }: Props) {
  const { t } = useTranslation('workbench-v2');
  const descriptionId = useId();
  return (
    <div className="space-y-1">
      <select
        data-testid={testId}
        aria-label={`${label} ${t('step4.basic.circuit.role_label')}`}
        aria-describedby={saved ? descriptionId : undefined}
        value={value}
        disabled={readonly || saved}
        onChange={(event) => {
          if (readonly || saved) return;
          const role = event.target.value;
          if (role === 'state' || (numeric && (role === 'gauge' || role === 'counter_raw'))) onChange(role);
        }}
        className="rounded border border-[var(--gw-control)] bg-[var(--gw-surface)] px-1.5 py-0.5 text-xs text-[var(--gw-text)]"
      >
        <option value="state">{t('step4.basic.circuit.role_state')}</option>
        {numeric && <>
          <option value="gauge">{t('step4.basic.circuit.role_gauge')}</option>
          <option value="counter_raw">{t('step4.basic.circuit.role_counter_raw')}</option>
        </>}
      </select><ParameterHelp help="circuit.role" />
      {saved && <p id={descriptionId} data-testid={`${testId}-locked`} className="text-xs text-[var(--gw-muted)]">
        {t('step4.basic.circuit.saved_role_unsupported')}
      </p>}
    </div>
  );
}
