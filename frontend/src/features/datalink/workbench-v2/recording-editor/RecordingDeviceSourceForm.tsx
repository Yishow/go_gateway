import { useTranslation } from 'react-i18next';
import type { DeviceSourceBuffer } from '../../../../pages/datalink/workbench-v2/deviceSourceSubmission';
import { ConnectionConfigForm } from '../steps/step1/ConnectionConfigForm';
import { PROTOCOLS, getDefaultConfig } from '../state/protocols';
import { getProtocolSupport } from '../state/protocolSupport';

export function RecordingDeviceSourceForm({ buffer, lockedProtocol, onChange, onProtocol }: {
  buffer: DeviceSourceBuffer; lockedProtocol: boolean;
  onChange: (buffer: DeviceSourceBuffer) => void;
  onProtocol: (buffer: DeviceSourceBuffer) => void;
}) {
  const { t } = useTranslation('recording-editor');
  const { t: workbench } = useTranslation('workbench-v2');
  return <div className="space-y-4">
    <label className="block">{t('device.name')}<input data-testid="recording-source-name" value={buffer.name}
      onChange={event => onChange({ ...buffer, name: event.target.value })} /></label>
    <label className="block">{t('device.protocol')}<select value={buffer.protocol} disabled={lockedProtocol}
      onChange={event => {
        const protocol = PROTOCOLS.find(item => item.id === event.target.value)?.id;
        if (protocol && getProtocolSupport(protocol).available) onProtocol({ ...buffer, protocol, config: getDefaultConfig(protocol) });
      }}>{PROTOCOLS.map(protocol => <option key={protocol.id} value={protocol.id} disabled={!getProtocolSupport(protocol.id).available}>
        {workbench(protocol.nameKey)}{!getProtocolSupport(protocol.id).available ? ` — ${t('device.unsupported')}` : ''}
      </option>)}</select></label>
    <ConnectionConfigForm protocol={buffer.protocol} config={buffer.config} rawNumericText
      onChange={patch => onChange({ ...buffer, config: { ...buffer.config, ...patch } })} />
    <p className="gw-muted">{t('device.scopeNote')}</p>
  </div>;
}
