import { useId } from 'react';
import { useTranslation } from 'react-i18next';
import type { DeviceSourceBuffer } from '../../../../pages/datalink/workbench-v2/deviceSourceSubmission';
import type { AcquisitionField, DeviceAcquisitionInput, RuntimeIntervalReadout } from '../../../../types/acquisitionSettings';
import { ParameterHelp } from '../components/ParameterHelp';
import {
  ACQUISITION_INTERVAL_MAX_MS, ACQUISITION_INTERVAL_MIN_MS, acquisitionInputFromSettings, acquisitionIssues, acquisitionProfileOptions,
  randomPointsText, randomProfileLimit, readAcquisitionDeclaration, withAcquisitionProfile,
} from '../state/acquisitionSettings';
import { RecordingSourceRangeField, RecordingSourceRangeGrid } from './RecordingSourceRangeFields';
import './recordingAcquisitionSettings.css';

function RuntimeInterval({ readout }: { readout?: RuntimeIntervalReadout }) {
  const { t } = useTranslation('acquisition-settings');
  const state = readout ?? { status: 'unavailable' as const };
  return <div className="recording-acquisition-runtime" role="status" data-testid="acquisition-runtime" data-runtime-state={state.status}>
    <span className="recording-acquisition-runtime-title">{t('runtime.title')}</span>
    {state.status === 'ready' ? <span>{t('runtime.values', { configured: state.configured_ms.toLocaleString(), effective: state.effective_ms.toLocaleString() })}
      {state.effective_ms > state.configured_ms && <> · {t('runtime.stretched')}</>}</span>
      : <span>{t(`runtime.${state.status === 'not-running' ? 'notRunning' : state.status}`)}</span>}
  </div>;
}

/** Presentation over the device buffer only; validation lives in acquisitionSettings.ts. */
export function RecordingAcquisitionSettings({ buffer, existingDevice = false, runtimeInterval, onChange }: {
  buffer: DeviceSourceBuffer; existingDevice?: boolean; runtimeInterval?: RuntimeIntervalReadout; onChange: (buffer: DeviceSourceBuffer) => void;
}) {
  const { t } = useTranslation('acquisition-settings');
  const titleId = useId();
  const issues = acquisitionIssues(buffer), options = acquisitionProfileOptions(buffer.protocol);
  const input = buffer.acquisitionInput ?? (buffer.acquisition ? acquisitionInputFromSettings(buffer.acquisition) : undefined);
  const declared = readAcquisitionDeclaration(buffer.config), limit = randomProfileLimit(declared.profile);
  const error = (field: AcquisitionField) => issues[field] && t(`errors.${issues[field]}`, { limit, min: ACQUISITION_INTERVAL_MIN_MS.toLocaleString(), max: ACQUISITION_INTERVAL_MAX_MS.toLocaleString() });
  const setInput = (patch: Partial<DeviceAcquisitionInput>) => { if (input) onChange({ ...buffer, acquisitionInput: { ...input, ...patch } }); };
  const unrecognized = options && declared.profile !== 'undeclared' && !options.some(option => option === declared.profile);
  return <section className="recording-acquisition-settings space-y-4" aria-labelledby={titleId} data-testid="acquisition-settings">
    <h3 id={titleId} className="recording-acquisition-title">{t('title')}</h3>
    {input ? <>
      <RecordingSourceRangeGrid>
        <RecordingSourceRangeField label={t('fields.interval')} help="acquisition.interval" hint={t('hints.interval')} error={error('interval')} required>
          <input data-testid="acquisition-interval" inputMode="numeric" value={input.interval_ms} aria-invalid={Boolean(issues.interval)} onChange={event => setInput({ interval_ms: event.target.value })} />
        </RecordingSourceRangeField>
        <RecordingSourceRangeField label={t('fields.autoInterval')} help="acquisition.autoInterval" hint={t('hints.autoInterval')}>
          <span className="recording-acquisition-toggle"><input type="checkbox" data-testid="acquisition-auto" checked={input.auto_interval} onChange={event => setInput({ auto_interval: event.target.checked })} />
            {t(input.auto_interval ? 'autoOn' : 'autoOff')}</span>
        </RecordingSourceRangeField>
        <RecordingSourceRangeField label={t('fields.staleThreshold')} help="acquisition.staleThreshold" hint={t(input.stale_persisted ? 'hints.staleSaved' : 'hints.stale')} error={error('stale')}>
          <input data-testid="acquisition-stale" inputMode="numeric" value={input.stale_threshold_ms} placeholder={t('placeholders.stale')} aria-invalid={Boolean(issues.stale)} onChange={event => setInput({ stale_threshold_ms: event.target.value })} />
        </RecordingSourceRangeField>
        <RecordingSourceRangeField label={t('fields.frozenThreshold')} help="acquisition.frozenThreshold" hint={t('hints.frozen')} error={error('frozen')} required>
          <input data-testid="acquisition-frozen" inputMode="numeric" value={input.frozen_threshold_ms} aria-invalid={Boolean(issues.frozen)} onChange={event => setInput({ frozen_threshold_ms: event.target.value })} />
        </RecordingSourceRangeField>
      </RecordingSourceRangeGrid>
      <RuntimeInterval readout={runtimeInterval} />
    </> : <p className="gw-muted" data-testid={existingDevice ? 'acquisition-settings-unavailable' : 'acquisition-new-device-note'}>{t(existingDevice ? 'unavailable' : 'newDevice')}</p>}
    {options ? <RecordingSourceRangeGrid>
      <RecordingSourceRangeField label={t('fields.profile')} help="acquisition.profile" hint={t('hints.profile')} error={error('profile')}>
        <select data-testid="acquisition-profile" value={declared.profile} aria-invalid={Boolean(issues.profile)}
          onChange={event => { const next = randomProfileLimit(event.target.value); onChange({ ...buffer, config: withAcquisitionProfile(buffer.config, event.target.value), randomPointsInput: next ? String(next) : undefined }); }}>
          {options.map(option => <option key={option} value={option}>{t(`profiles.${option}`)}</option>)}
          {unrecognized && <option value={declared.profile}>{t('profiles.unrecognized', { profile: declared.profile })}</option>}
        </select>
      </RecordingSourceRangeField>
      {limit > 0 && <RecordingSourceRangeField label={t('fields.randomPoints')} help="acquisition.randomPoints" hint={t('hints.randomPoints', { limit })} error={error('randomPoints')}>
        <input data-testid="acquisition-random-points" inputMode="numeric" value={randomPointsText(buffer)} aria-invalid={Boolean(issues.randomPoints)} onChange={event => onChange({ ...buffer, randomPointsInput: event.target.value })} />
      </RecordingSourceRangeField>}
    </RecordingSourceRangeGrid> : <div className="recording-source-field recording-acquisition-note">
      <p className="gw-muted" data-testid="acquisition-modbus-note">{t('modbusNote')}</p>
      <ParameterHelp compact help="acquisition.contiguous" />
    </div>}
  </section>;
}
