import type { ReactNode } from 'react';
import { RecordingSourceRangeField, RecordingSourceRangeGrid } from './RecordingSourceRangeFields';
import { RecordingRangeSchedule } from './RecordingRangeSchedulePanel';
import { buildStudioEntityLocatorTarget } from '../../../../services/studioEntityLocator';
import { useTranslation } from 'react-i18next';
import { recordingRangeDataTypes } from './recordingRangeBuffer';
import type { RecordingRangeController } from './useRecordingRangeBuffer';
export function RecordingRangeContent({ controller: range, available, onBack, onSubmit, onInspect, onPrepare, beforeForm, onReview, backLabel }: { controller: RecordingRangeController; available: boolean; onBack: () => void; onSubmit: () => void; onInspect: () => void; onPrepare: () => void; beforeForm?: ReactNode; backLabel?: string; onReview?: (ruleId: string) => void }) {
  const { t } = useTranslation('recording-editor');
  const { t: sourceText } = useTranslation('recording-home');
  const addressExample = range.device?.protocol.startsWith('modbus') ? '40001' : range.device && ['mc_3e', 'fatek_fbs'].includes(range.device.protocol) ? 'D0' : undefined;
  const inspection = range.outcome?.inspection;
  const ruleId = range.outcome?.rule?.id ?? inspection?.ruleId;
  const reviewTarget = ruleId && range.outcome ? buildStudioEntityLocatorTarget({ step: 2,
    workspace_id: range.outcome.request.workspaceId, device_id: range.outcome.request.deviceId, rule_id: ruleId,
    mapping_id: null, group_id: null, group_revision: null }) : null;
  return <><div className="gw-modal-body space-y-4">{beforeForm}<p data-testid="recording-range-device">{range.target?.name}</p>
    {available ? <p className="gw-muted" data-testid="recording-source-confirmed">{t('range.savedDeviceNote')}</p> : <p role="alert">{t('picker.unavailable')}</p>}
    <RecordingSourceRangeGrid>
      <RecordingSourceRangeField label={sourceText('source.fields.address')} help="range.address" hint={addressExample ? sourceText('source.fields.addressHint', { example: addressExample }) : sourceText('source.fields.addressUnknown')}>
        <input disabled={!available} value={range.buffer.address} placeholder={addressExample} onChange={event => range.change({ ...range.buffer, address: event.target.value })} />
      </RecordingSourceRangeField>
      <RecordingSourceRangeField label={t('range.count')} help="range.count" hint={sourceText('source.fields.countHint')}>
        <input disabled={!available} value={range.buffer.count} inputMode="numeric" placeholder="16" onChange={event => range.change({ ...range.buffer, count: event.target.value })} />
      </RecordingSourceRangeField>
      <RecordingSourceRangeField label={t('range.dataType')} help="range.dataType" hint={sourceText('source.fields.typeHint')}>
        <select disabled={!available} value={range.buffer.dataType} onChange={event => { const dataType = recordingRangeDataTypes.find(type => type === event.target.value); if (dataType) range.change({ ...range.buffer, dataType }); }}>
          {recordingRangeDataTypes.map(type => <option key={type} value={type}>{type}</option>)}
        </select>
      </RecordingSourceRangeField>
      <RecordingSourceRangeField label={sourceText('source.fields.prefix')} help="range.name" hint={sourceText('source.fields.prefixHint')}>
        <input disabled={!available} value={range.buffer.name} placeholder="TEMP_" onChange={event => range.change({ ...range.buffer, name: event.target.value })} />
      </RecordingSourceRangeField>
    </RecordingSourceRangeGrid>
    <details className="recording-source-option"><summary>{t('range.transform')}</summary><RecordingSourceRangeGrid>
      {(['scale', 'offset'] as const).map(field => <RecordingSourceRangeField key={field} label={t(`range.${field}`)} help={`range.${field}`}>
        <input disabled={!available} value={range.buffer[field]} inputMode="decimal" onChange={event => range.change({ ...range.buffer, [field]: event.target.value })} />
      </RecordingSourceRangeField>)}
      <RecordingSourceRangeField label={t('range.dataFormat')} help="range.dataFormat">
        <select disabled={!available} value={range.buffer.dataFormat} onChange={event => { const dataFormat = event.target.value; if (dataFormat === '' || dataFormat === 'ABCD' || dataFormat === 'BADC' || dataFormat === 'CDAB' || dataFormat === 'DCBA') range.change({ ...range.buffer, dataFormat }); }}>
          <option value="">{t('range.defaultFormat')}</option>{['ABCD', 'BADC', 'CDAB', 'DCBA'].map(format => <option key={format} value={format}>{format}</option>)}
        </select>
      </RecordingSourceRangeField>
      <RecordingSourceRangeField label={t('range.unit')} help="range.unit"><input disabled={!available} value={range.buffer.unit} onChange={event => range.change({ ...range.buffer, unit: event.target.value })} /></RecordingSourceRangeField>
    </RecordingSourceRangeGrid><p className="gw-muted">{t('range.unitNote')}</p></details>
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
    {reviewTarget && (onReview && ruleId ? <button type="button" onClick={() => onReview(ruleId)}>{t('range.reviewSource')}</button> : <a href={reviewTarget} target="_blank" rel="noopener noreferrer">{t('range.reviewSource')}</a>)}
    <details className="recording-source-option"><summary>{t('range.schedule')}</summary><RecordingRangeSchedule workspaceId={range.workspaceId} deviceId={range.target?.id} outcome={range.outcome} /></details>
    {range.outcome?.status === 'saved-inactive' && <p className="gw-muted">{t('range.prepareNote')}</p>}
    {range.outcome?.preparationError && <p role="alert">{t(`range.${range.outcome.preparationError}`)}</p>}
    {range.outcome && ['confirmed', 'saved-inactive'].includes(range.outcome.status) && range.dirty && <p role="status" data-testid="recording-range-newer">{t('range.newer')}</p>}
  </div><footer className="gw-modal-footer"><button type="button" data-testid="recording-range-back" onClick={onBack}>{backLabel ?? t(range.outcome?.status === 'saved-inactive' ? 'range.keepOffline' : 'source.back')}</button>
    {range.locked && !range.pending && <button type="button" data-testid="recording-range-inspect" disabled={range.checking || !range.canInspect} onClick={onInspect}>{t('range.inspect')}</button>}
    {range.outcome?.status === 'saved-inactive' ? <button type="button" className="gw-primary" data-testid="recording-range-prepare" disabled={!available || !range.canPrepareSaved} onClick={onPrepare}>{t('range.prepare')}</button> : <button type="button" className="gw-primary" disabled={!available || !range.canSubmit} onClick={onSubmit}>{t('range.save')}</button>}
  </footer></>;
}
