import { RecordingSourceRangeField, RecordingSourceRangeGrid } from './RecordingSourceRangeFields';
import { useTranslation } from 'react-i18next';
import type { DeviceSourceBuffer } from '../../../../pages/datalink/workbench-v2/deviceSourceSubmission';
import { ConnectionConfigForm } from '../steps/step1/ConnectionConfigForm';
import { PROTOCOLS, getDefaultConfig } from '../state/protocols';
import { getProtocolSupport } from '../state/protocolSupport';
import type { RuntimeIntervalReadout } from '../../../../types/acquisitionSettings';
import { RecordingAcquisitionSettings } from './RecordingAcquisitionSettings';

export function RecordingDeviceSourceForm({ buffer, lockedProtocol, onChange, onProtocol, onSubtype, existingDevice, runtimeInterval }: {
  buffer: DeviceSourceBuffer; lockedProtocol: boolean; existingDevice?: boolean; runtimeInterval?: RuntimeIntervalReadout;
  onChange: (buffer: DeviceSourceBuffer) => void;
  onProtocol: (buffer: DeviceSourceBuffer) => void;
  onSubtype?: (buffer: DeviceSourceBuffer) => void;
}) {
  const { t } = useTranslation('recording-editor');
  const { t: workbench } = useTranslation('workbench-v2');
  return <div className="recording-source-device-fields space-y-4"><RecordingSourceRangeGrid>
    <RecordingSourceRangeField label={t('device.name')} help="device.name"><input data-testid="recording-source-name" value={buffer.name}
      onChange={event => onChange({ ...buffer, name: event.target.value })} /></RecordingSourceRangeField>
    <RecordingSourceRangeField label={t('device.protocol')} help="device.protocol"><select value={buffer.protocol} disabled={lockedProtocol}
      onChange={event => {
        const protocol = PROTOCOLS.find(item => item.id === event.target.value)?.id;
        if (protocol && getProtocolSupport(protocol).available) onProtocol({ ...buffer, protocol, config: getDefaultConfig(protocol) });
      }}>{PROTOCOLS.map(protocol => <option key={protocol.id} value={protocol.id} disabled={!getProtocolSupport(protocol.id).available}>
        {workbench(protocol.nameKey)}{!getProtocolSupport(protocol.id).available ? ` — ${t('device.unsupported')}` : ''}
      </option>)}</select></RecordingSourceRangeField></RecordingSourceRangeGrid>
    <ConnectionConfigForm protocol={buffer.protocol} config={buffer.config} rawNumericText inlineAdvanced
      onChange={patch => {
        if (buffer.protocol === 'fatek_fbs' && patch.mode !== undefined && patch.mode !== (buffer.config.mode ?? 'tcp') && onSubtype) {
          onSubtype({ ...buffer, config: { mode: patch.mode } });
        } else onChange({ ...buffer, config: { ...buffer.config, ...patch } });
      }} />
    {buffer.protocol !== 'mqtt' && <RecordingAcquisitionSettings buffer={buffer} existingDevice={existingDevice} runtimeInterval={runtimeInterval} onChange={onChange} />}
    <p className="gw-muted">{t('device.scopeNote')}</p>
  </div>;
}
