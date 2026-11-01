import { useTranslation } from 'react-i18next';
import { Field } from '../../components/Field';
import { Input, Select } from '../../components/inputs';
import { MODBUS_RTU_BAUD_RATES } from '../../state/protocolConnectionOptions';
import type { ConnectionConfigFormProps } from './ConnectionConfigForm';

type Props = Pick<ConnectionConfigFormProps, 'protocol' | 'config' | 'onChange' | 'rawNumericText'>;
export function FatekModeField({ config, onChange }: Pick<Props, 'config' | 'onChange'>) {
  const { t } = useTranslation('recording-editor');
  return <Field help="connection.fatekMode" label={t('connection.mode')} required><Select value={String(config.mode ?? 'tcp')}
    data-testid="select-fatek-mode" onChange={event => onChange(event.target.value === 'serial'
      ? { mode: 'serial', baud_rate: config.baud_rate ?? 9600, parity: config.parity ?? 'even',
        data_bits: config.data_bits ?? 7, stop_bits: config.stop_bits ?? 1 }
      : { mode: 'tcp' })}>
    <option value="tcp">TCP</option><option value="serial">{t('connection.serial')}</option>
  </Select></Field>;
}

export function SerialConnectionFields({ protocol, config, onChange, rawNumericText }: Props) {
  const { t } = useTranslation('workbench-v2');
  const { t: editor } = useTranslation('recording-editor');
  const fatek = protocol === 'fatek_fbs';
  const number = (value: string) => rawNumericText ? value : value === '' ? '' : Number(value);
  const parity = String(config.parity ?? (fatek ? 'even' : 'none')).toLowerCase();
  const displayedParity = ['even', 'e'].includes(parity) ? 'E' : ['odd', 'o'].includes(parity) ? 'O' : 'N';
  const baud = config.baud_rate ?? config.baud ?? 9600;
  return <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
    {fatek && <FatekModeField config={config} onChange={onChange} />}
    <Field help="connection.serialPort" label={t('step1.fields.serial_port')} required hint={editor('connection.serialHost')}>
      <Input value={String(config.serial_port ?? (!fatek ? config.port : '') ?? '')} data-testid="input-rtu-port"
        placeholder="/dev/ttyUSB0" onChange={event => onChange({ serial_port: event.target.value })} />
    </Field>
    <Field help={fatek ? 'connection.baudFatek' : 'connection.baudModbus'} label={t('step1.fields.baud_rate')} required>
      {fatek ? <Input type={rawNumericText ? 'text' : 'number'} inputMode="numeric" value={String(baud)}
        data-testid="input-fatek-baud" onChange={event => onChange({ baud_rate: number(event.target.value) })} /> :
        <Select value={String(baud)} data-testid="select-baud" onChange={event => onChange({ baud_rate: Number(event.target.value) })}>
          {!MODBUS_RTU_BAUD_RATES.includes(Number(baud)) && <option value={String(baud)} disabled>{String(baud)} — {editor('device.unsupported')}</option>}
          {MODBUS_RTU_BAUD_RATES.map(value => <option key={value} value={value}>{value}</option>)}
        </Select>}
    </Field>
    <Field help="connection.parity" label={t('step1.fields.parity')} required><Select value={displayedParity} data-testid="select-parity"
      onChange={event => onChange({ parity: event.target.value })}>
      <option value="N">{t('step1.parity.none')}</option><option value="E">{t('step1.parity.even')}</option><option value="O">{t('step1.parity.odd')}</option>
    </Select></Field>
    <Field help={fatek ? 'connection.fatekStation' : 'connection.modbusSlave'} label={t(fatek ? 'step1.fields.station' : 'step1.fields.slave_id')} required>
      <Input type={rawNumericText ? 'text' : 'number'} inputMode="numeric"
        value={String((fatek ? config.station_no ?? config.station ?? config.slave_id : config.slave_id) ?? '')}
        data-testid={fatek ? 'input-slave-id' : 'input-rtu-slave-id'} onChange={event => onChange(fatek
          ? { station_no: number(event.target.value), station: number(event.target.value), slave_id: number(event.target.value) }
          : { slave_id: number(event.target.value) })} />
    </Field>
    <Field help="connection.dataBits" label={editor('connection.dataBits')}><Select value={String(config.data_bits ?? (fatek ? 7 : 8))}
      data-testid="select-data-bits" onChange={event => onChange({ data_bits: Number(event.target.value) })}>
      <option value="7">7</option><option value="8">8</option>
    </Select></Field>
    <Field help="connection.stopBits" label={editor('connection.stopBits')}><Select value={String(config.stop_bits ?? 1)}
      data-testid="select-stop-bits" onChange={event => onChange({ stop_bits: Number(event.target.value) })}>
      <option value="1">1</option><option value="2">2</option>
    </Select></Field>
    <Field help={fatek ? 'connection.timeoutOther' : 'connection.timeoutModbus'} label={t('step1.fields.timeout')} hint={t('step1.fields.timeout_hint')}>
      <Input type={rawNumericText ? 'text' : 'number'} inputMode="numeric" value={String(config.timeout ?? 5)}
        data-testid="input-timeout" onChange={event => onChange({ timeout: number(event.target.value) })} />
    </Field>
    {!fatek && <ConnectionDataFormat config={config} onChange={onChange} />}
  </div>;
}

function ConnectionDataFormat({ config, onChange }: Pick<Props, 'config' | 'onChange'>) {
  const { t } = useTranslation('recording-editor');
  return <Field help="connection.dataFormat" label={t('connection.dataFormat')}><Select value={String(config.data_format ?? '')} data-testid="select-data-format"
    onChange={event => onChange({ data_format: event.target.value })}>
    <option value="">{t('connection.defaultFormat')}</option>
    {['ABCD', 'BADC', 'CDAB', 'DCBA'].map(format => <option key={format} value={format}>{format}</option>)}
  </Select></Field>;
}

export function ConnectionAdvancedFields({ protocol, config, onChange, rawNumericText }: Props) {
  const { t } = useTranslation('recording-editor');
  if (protocol === 'fatek_fbs') return null;
  return <details className="sm:col-span-2"><summary>{t('connection.advanced')}</summary>
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
      {protocol === 'mc_3e' && (['network_no', 'pc_no', 'io_no'] as const).map(key =>
        <Field key={key} help={key === 'network_no' ? 'connection.networkNo' : key === 'pc_no' ? 'connection.pcNo' : 'connection.ioNo'} label={t(`connection.${key}`)}><Input type={rawNumericText ? 'text' : 'number'} inputMode="numeric"
          value={String(config[key] ?? (key === 'io_no' ? config.module_io_no : undefined) ?? (key === 'pc_no' ? 255 : key === 'io_no' ? 1023 : 0))}
          data-testid={`input-${key.replaceAll('_', '-')}`} onChange={event => onChange({ [key]: rawNumericText ? event.target.value : event.target.value === '' ? '' : Number(event.target.value) })} />
        </Field>)}
      <ConnectionDataFormat config={config} onChange={onChange} />
    </div>
  </details>;
}
