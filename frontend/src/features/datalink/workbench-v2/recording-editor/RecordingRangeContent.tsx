import { ParameterHelp } from '../components/ParameterHelp';
import { RecordingRangeSchedule } from './RecordingRangeSchedule';
import { buildStudioEntityLocatorTarget } from '../../../../services/studioEntityLocator';
import { useTranslation } from 'react-i18next';
import { recordingRangeDataTypes } from './recordingRangeBuffer';
import type { RecordingRangeController } from './useRecordingRangeBuffer';
export function RecordingRangeContent({ controller: range, available, onBack, onSubmit, onInspect, onPrepare }: { controller: RecordingRangeController; available: boolean; onBack: () => void; onSubmit: () => void; onInspect: () => void; onPrepare: () => void }) {
  const { t } = useTranslation('recording-editor');
  const inspection = range.outcome?.inspection;
  const ruleId = range.outcome?.rule?.id ?? inspection?.ruleId;
  const reviewTarget = ruleId && range.outcome ? buildStudioEntityLocatorTarget({ step: 2,
    workspace_id: range.outcome.request.workspaceId, device_id: range.outcome.request.deviceId, rule_id: ruleId,
    mapping_id: null, group_id: null, group_revision: null }) : null;
  return <><div className="gw-modal-body space-y-4"><p data-testid="recording-range-device">{range.target?.name}</p>
    {available ? <p className="gw-muted" data-testid="recording-source-confirmed">{t('range.savedDeviceNote')}</p> : <p role="alert">{t('picker.unavailable')}</p>}
    {(['name', 'address', 'count'] as const).map(field => <div key={field} className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 block">{t(`range.${field}`)}<input disabled={!available} value={range.buffer[field]}
      inputMode={field === 'count' ? 'numeric' : undefined} onChange={event => range.change({ ...range.buffer, [field]: event.target.value })} /></label><ParameterHelp help={`range.${field}`} /></div>)}
    <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1 block">{t('range.dataType')}<select disabled={!available} value={range.buffer.dataType} onChange={event => { const dataType = recordingRangeDataTypes.find(type => type === event.target.value); if (dataType) range.change({ ...range.buffer, dataType }); }}>
      {recordingRangeDataTypes.map(type => <option key={type} value={type}>{type}</option>)}</select></label><ParameterHelp help="range.dataType" /></div>
    <details><summary>{t('range.transform')}</summary><div className="grid grid-cols-1 sm:grid-cols-2 gap-3">{(['scale', 'offset'] as const).map(field => <div key={field} className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1">{t(`range.${field}`)}<input disabled={!available} value={range.buffer[field]} inputMode="decimal" onChange={event => range.change({ ...range.buffer, [field]: event.target.value })} /></label><ParameterHelp help={`range.${field}`} /></div>)}
      <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1">{t('range.dataFormat')}<select disabled={!available} value={range.buffer.dataFormat} onChange={event => { const dataFormat = event.target.value; if (dataFormat === '' || dataFormat === 'ABCD' || dataFormat === 'BADC' || dataFormat === 'CDAB' || dataFormat === 'DCBA') range.change({ ...range.buffer, dataFormat }); }}>
        <option value="">{t('range.defaultFormat')}</option>{['ABCD', 'BADC', 'CDAB', 'DCBA'].map(format => <option key={format} value={format}>{format}</option>)}
      </select></label><ParameterHelp help="range.dataFormat" /></div>
      <div className="flex min-w-0 items-start gap-2"><label className="min-w-0 flex-1">{t('range.unit')}<input disabled={!available} value={range.buffer.unit} onChange={event => range.change({ ...range.buffer, unit: event.target.value })} /></label><ParameterHelp help="range.unit" /></div>
    </div><p className="gw-muted">{t('range.unitNote')}</p></details>
    {range.parsed && <p className="gw-muted">{t('range.prefix', { prefix: range.parsed.rule.naming_prefix })}</p>}
    {range.parsed && <p role="status">{t('range.span', { count: range.parsed.rule.count, span: range.parsed.span, end: range.parsed.end })}</p>}
    {range.dirty && !range.parsed && !range.preparing && !range.blocked && <p role="alert">{t('range.invalid')}</p>}
    {!range.preparing && range.device && !range.activationAllowed && !range.outcome && <p role="status">{t('range.offlineNote')}</p>}
    {range.preparing && <p role="status">{t('device.checkingScope')}</p>}
    {range.blocked && <p role="alert">{t('device.scopeBlocked')}</p>}
    {range.pending && <p role="status">{t('range.pending')}</p>}
    {range.outcome && range.outcome.status !== 'pending' && range.outcome.status !== 'confirmed' && <p role="alert" data-testid={`recording-range-${range.outcome.status}`}>{t(`range.${range.outcome.status}`)}</p>}
    {inspection?.status === 'current-state' && <p role="status" data-testid="recording-range-current-state">{t('range.inspectedUnknown')} · {inspection.currentRule?.start_address}</p>}
    {inspection?.status === 'unavailable' && <p role="alert" data-testid="recording-range-inspection-unavailable">{t(inspection.ruleId ? 'range.inspectFailed' : 'range.inspectUnavailable')}</p>}
    {reviewTarget && <a href={reviewTarget} target="_blank" rel="noopener noreferrer">{t('range.reviewSource')}</a>}
    <RecordingRangeSchedule workspaceId={range.workspaceId} deviceId={range.target?.id} outcome={range.outcome} />
    {range.outcome?.status === 'saved-inactive' && <p className="gw-muted">{t('range.prepareNote')}</p>}
    {range.outcome?.preparationError && <p role="alert">{t(`range.${range.outcome.preparationError}`)}</p>}
    {range.outcome && ['confirmed', 'saved-inactive'].includes(range.outcome.status) && range.dirty && <p role="status" data-testid="recording-range-newer">{t('range.newer')}</p>}
  </div><footer className="gw-modal-footer"><button type="button" data-testid="recording-range-back" onClick={onBack}>{t(range.outcome?.status === 'saved-inactive' ? 'range.keepOffline' : 'source.back')}</button>
    {range.locked && !range.pending && <button type="button" data-testid="recording-range-inspect" disabled={range.checking || !range.canInspect} onClick={onInspect}>{t('range.inspect')}</button>}
    {range.outcome?.status === 'saved-inactive' ? <button type="button" className="gw-primary" data-testid="recording-range-prepare" disabled={!available || !range.canPrepareSaved} onClick={onPrepare}>{t('range.prepare')}</button> : <button type="button" className="gw-primary" disabled={!available || !range.canSubmit} onClick={onSubmit}>{t('range.save')}</button>}
  </footer></>;
}
