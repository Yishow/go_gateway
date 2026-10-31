import { useTranslation } from 'react-i18next';
import { recordingRangeDataTypes } from './recordingRangeBuffer';
import type { RecordingRangeController } from './useRecordingRangeBuffer';
export function RecordingRangeContent({ controller: range, available, onBack, onSubmit, onInspect }: { controller: RecordingRangeController; available: boolean; onBack: () => void; onSubmit: () => void; onInspect: () => void }) {
  const { t } = useTranslation('recording-editor');
  return <><div className="gw-modal-body space-y-4"><p data-testid="recording-range-device">{range.target?.name}</p>
    {available ? <p className="gw-muted" data-testid="recording-source-confirmed">{t('range.savedDeviceNote')}</p> : <p role="alert">{t('picker.unavailable')}</p>}
    {(['name', 'address', 'count'] as const).map(field => <label key={field} className="block">{t(`range.${field}`)}<input disabled={!available} value={range.buffer[field]}
      inputMode={field === 'count' ? 'numeric' : undefined} onChange={event => range.change({ ...range.buffer, [field]: event.target.value })} /></label>)}
    <label className="block">{t('range.dataType')}<select disabled={!available} value={range.buffer.dataType} onChange={event => { const dataType = recordingRangeDataTypes.find(type => type === event.target.value); if (dataType) range.change({ ...range.buffer, dataType }); }}>
      {recordingRangeDataTypes.map(type => <option key={type} value={type}>{type}</option>)}</select></label>
    <details><summary>{t('range.transform')}</summary>{(['scale', 'offset'] as const).map(field => <label key={field}>{t(`range.${field}`)}<input disabled={!available} value={range.buffer[field]} inputMode="decimal" onChange={event => range.change({ ...range.buffer, [field]: event.target.value })} /></label>)}</details>
    {range.parsed && <p className="gw-muted">{t('range.prefix', { prefix: range.parsed.rule.naming_prefix })}</p>}
    {range.parsed && <p role="status">{t('range.span', { count: range.parsed.rule.count, span: range.parsed.span, end: range.parsed.end })}</p>}
    {range.dirty && !range.parsed && !range.preparing && !range.blocked && <p role="alert">{t('range.invalid')}</p>}
    {range.preparing && <p role="status">{t('device.checkingScope')}</p>}
    {range.blocked && <p role="alert">{t('device.scopeBlocked')}</p>}
    {range.pending && <p role="status">{t('range.pending')}</p>}
    {range.outcome && range.outcome.status !== 'pending' && range.outcome.status !== 'confirmed' && <p role="alert" data-testid={`recording-range-${range.outcome.status}`}>{t(`range.${range.outcome.status}`)}</p>}
    {range.outcome?.status === 'confirmed' && range.dirty && <p role="status" data-testid="recording-range-newer">{t('range.newer')}</p>}
  </div><footer className="gw-modal-footer"><button type="button" data-testid="recording-range-back" onClick={onBack}>{t('source.back')}</button>
    {range.locked && !range.pending && <button type="button" data-testid="recording-range-inspect" disabled={range.checking} onClick={onInspect}>{t('range.inspect')}</button>}
    <button type="button" className="gw-primary" disabled={!available || !range.canSubmit} onClick={onSubmit}>{t('range.save')}</button>
  </footer></>;
}
