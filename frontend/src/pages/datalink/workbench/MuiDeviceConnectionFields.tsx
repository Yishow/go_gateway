import Box from '@mui/material/Box';
import Checkbox from '@mui/material/Checkbox';
import FormControlLabel from '@mui/material/FormControlLabel';
import TextField from '@mui/material/TextField';
import { useTranslation } from 'react-i18next';
import type { ProtocolType } from '../../../types/datalink';
import { buildConnectionDataFormatOptions, normalizeConnectionDataFormatValue } from './workbenchDeviceDataFormat';
import type { DeviceConnectionConfig } from './workbenchDeviceFormModel';

type FieldErrorMap = Record<string, string>;

type Props = {
  protocol: ProtocolType;
  connectionConfig: DeviceConnectionConfig;
  fieldErrors: FieldErrorMap;
  onChangeValue: (key: string, value: DeviceConnectionConfig[string]) => void;
  onChangeNumber: (key: string, raw: string) => void;
  onChangeTopics: (raw: string) => void;
};

const row = { display: 'grid', gap: 2, gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' } } as const;

function str(config: DeviceConnectionConfig, key: string): string {
  const v = config[key];
  if (typeof v === 'string') return v;
  if (typeof v === 'number') return String(v);
  if (Array.isArray(v)) return v.join('\n');
  return '';
}

export function MuiDeviceConnectionFields({
  protocol,
  connectionConfig,
  fieldErrors,
  onChangeValue,
  onChangeNumber,
  onChangeTopics,
}: Props) {
  const { t } = useTranslation();
  const dataFormatOptions = buildConnectionDataFormatOptions(t);

  switch (protocol) {
    case 'modbus_tcp':
    case 'modbus_udp':
      return (
        <>
          <Box sx={row}>
            <TextField label={t('workbench.device.connection.host')} size="small" value={str(connectionConfig, 'host')} onChange={(e) => onChangeValue('host', e.target.value)} placeholder="192.168.1.10" error={Boolean(fieldErrors.host)} helperText={fieldErrors.host} />
            <TextField label={t('workbench.device.connection.port')} size="small" value={str(connectionConfig, 'port')} onChange={(e) => onChangeNumber('port', e.target.value)} placeholder="502" error={Boolean(fieldErrors.port)} helperText={fieldErrors.port} inputMode="numeric" />
          </Box>
          <Box sx={row}>
            <TextField label={t('workbench.device.connection.slaveId')} size="small" value={str(connectionConfig, 'slave_id')} onChange={(e) => onChangeNumber('slave_id', e.target.value)} placeholder="1" inputMode="numeric" />
            <TextField label={t('workbench.device.connection.timeout')} size="small" value={str(connectionConfig, 'timeout')} onChange={(e) => onChangeNumber('timeout', e.target.value)} placeholder="5" inputMode="numeric" />
          </Box>
          <TextField label={t('workbench.device.connection.dataFormat')} size="small" select SelectProps={{ native: true }} value={normalizeConnectionDataFormatValue(protocol, connectionConfig.data_format)} onChange={(e) => onChangeValue('data_format', e.target.value)}>
            {dataFormatOptions.map((o) => <option key={o.value} value={o.value}>{o.label}</option>)}
          </TextField>
        </>
      );

    case 'modbus_rtu':
      return (
        <>
          <TextField label={t('workbench.device.connection.serialPort')} size="small" fullWidth value={str(connectionConfig, 'serial_port')} onChange={(e) => onChangeValue('serial_port', e.target.value)} placeholder="/dev/ttyUSB0" error={Boolean(fieldErrors.serial_port)} helperText={fieldErrors.serial_port} />
          <Box sx={row}>
              <TextField label={t('workbench.device.connection.baudRate')} size="small" select SelectProps={{ native: true }} value={str(connectionConfig, 'baud_rate') || '9600'} onChange={(e) => onChangeValue('baud_rate', Number.parseInt(e.target.value, 10))}>
                {[9600, 19200, 38400, 57600, 115200].map((b) => <option key={b} value={String(b)}>{b}</option>)}
              </TextField>
              <TextField label={t('workbench.device.connection.parity')} size="small" select SelectProps={{ native: true }} value={str(connectionConfig, 'parity') || 'none'} onChange={(e) => onChangeValue('parity', e.target.value)}>
                {['none', 'even', 'odd'].map((p) => <option key={p} value={p}>{t(`device.parity${p.charAt(0).toUpperCase() + p.slice(1)}`)}</option>)}
              </TextField>
          </Box>
          <Box sx={row}>
            <TextField label={t('workbench.device.connection.slaveId')} size="small" value={str(connectionConfig, 'slave_id')} onChange={(e) => onChangeNumber('slave_id', e.target.value)} placeholder="1" inputMode="numeric" />
            <TextField label={t('workbench.device.connection.timeout')} size="small" value={str(connectionConfig, 'timeout')} onChange={(e) => onChangeNumber('timeout', e.target.value)} placeholder="5" inputMode="numeric" />
          </Box>
          <TextField label={t('workbench.device.connection.dataFormat')} size="small" select SelectProps={{ native: true }} value={normalizeConnectionDataFormatValue('modbus_rtu', connectionConfig.data_format)} onChange={(e) => onChangeValue('data_format', e.target.value)}>
            {dataFormatOptions.map((o) => <option key={o.value} value={o.value}>{o.label}</option>)}
          </TextField>
        </>
      );

    case 'fatek_fbs': {
      const mode = typeof connectionConfig.mode === 'string' ? connectionConfig.mode : 'tcp';
      return (
        <>
          <TextField label={t('workbench.device.connection.mode')} size="small" select SelectProps={{ native: true }} value={mode} onChange={(e) => onChangeValue('mode', e.target.value)}>
            <option value="tcp">{t('workbench.device.mode.tcp')}</option>
            <option value="serial">{t('workbench.device.mode.serial')}</option>
          </TextField>
          {mode === 'serial' ? (
            <Box sx={row}>
              <TextField label={t('workbench.device.connection.serialPort')} size="small" value={str(connectionConfig, 'serial_port')} onChange={(e) => onChangeValue('serial_port', e.target.value)} placeholder="/dev/ttyUSB0" error={Boolean(fieldErrors.serial_port)} helperText={fieldErrors.serial_port} />
               <TextField label={t('workbench.device.connection.baudRate')} size="small" select SelectProps={{ native: true }} value={str(connectionConfig, 'baud_rate') || '9600'} onChange={(e) => onChangeValue('baud_rate', Number.parseInt(e.target.value, 10))}>
                 {[9600, 19200, 38400, 57600, 115200].map((b) => <option key={b} value={String(b)}>{b}</option>)}
               </TextField>
            </Box>
          ) : (
            <Box sx={row}>
              <TextField label={t('workbench.device.connection.host')} size="small" value={str(connectionConfig, 'host')} onChange={(e) => onChangeValue('host', e.target.value)} placeholder="192.168.1.10" error={Boolean(fieldErrors.host)} helperText={fieldErrors.host} />
              <TextField label={t('workbench.device.connection.port')} size="small" value={str(connectionConfig, 'port')} onChange={(e) => onChangeNumber('port', e.target.value)} placeholder="500" error={Boolean(fieldErrors.port)} helperText={fieldErrors.port} inputMode="numeric" />
            </Box>
          )}
          <Box sx={row}>
            <TextField label={t('workbench.device.connection.stationNo')} size="small" value={str(connectionConfig, 'station_no')} onChange={(e) => onChangeNumber('station_no', e.target.value)} placeholder="1" error={Boolean(fieldErrors.station_no)} helperText={fieldErrors.station_no} inputMode="numeric" />
            <TextField label={t('workbench.device.connection.timeout')} size="small" value={str(connectionConfig, 'timeout')} onChange={(e) => onChangeNumber('timeout', e.target.value)} placeholder="5" inputMode="numeric" />
          </Box>
        </>
      );
    }

    case 'mc_3e':
      return (
        <>
          <Box sx={row}>
            <TextField label={t('workbench.device.connection.host')} size="small" value={str(connectionConfig, 'host')} onChange={(e) => onChangeValue('host', e.target.value)} placeholder="192.168.1.10" error={Boolean(fieldErrors.host)} helperText={fieldErrors.host} />
            <TextField label={t('workbench.device.connection.port')} size="small" value={str(connectionConfig, 'port')} onChange={(e) => onChangeNumber('port', e.target.value)} placeholder="5000" error={Boolean(fieldErrors.port)} helperText={fieldErrors.port} inputMode="numeric" />
          </Box>
          <Box sx={row}>
            <TextField label={t('workbench.device.connection.networkNo')} size="small" value={str(connectionConfig, 'network_no')} onChange={(e) => onChangeNumber('network_no', e.target.value)} placeholder="0" inputMode="numeric" />
            <TextField label={t('workbench.device.connection.stationNo')} size="small" value={str(connectionConfig, 'station_no')} onChange={(e) => onChangeNumber('station_no', e.target.value)} placeholder="0" inputMode="numeric" />
          </Box>
          <Box sx={row}>
            <TextField label={t('workbench.device.connection.timeout')} size="small" value={str(connectionConfig, 'timeout')} onChange={(e) => onChangeNumber('timeout', e.target.value)} placeholder="5" inputMode="numeric" />
             <TextField label={t('workbench.device.connection.dataFormat')} size="small" select SelectProps={{ native: true }} value={normalizeConnectionDataFormatValue('mc_3e', connectionConfig.data_format)} onChange={(e) => onChangeValue('data_format', e.target.value)}>
               {dataFormatOptions.map((o) => <option key={o.value} value={o.value}>{o.label}</option>)}
             </TextField>
          </Box>
        </>
      );

    case 'mqtt':
      return (
        <>
          <Box sx={row}>
            <TextField label={t('workbench.device.connection.brokerUrl')} size="small" value={str(connectionConfig, 'broker_url')} onChange={(e) => onChangeValue('broker_url', e.target.value)} placeholder="mqtt://broker.local:1883" error={Boolean(fieldErrors.broker_url)} helperText={fieldErrors.broker_url} />
            <TextField label={t('workbench.device.connection.clientId')} size="small" value={str(connectionConfig, 'client_id')} onChange={(e) => onChangeValue('client_id', e.target.value)} placeholder="datalink-workbench" error={Boolean(fieldErrors.client_id)} helperText={fieldErrors.client_id} />
          </Box>
          <Box sx={row}>
            <TextField label={t('workbench.device.connection.username')} size="small" value={str(connectionConfig, 'username')} onChange={(e) => onChangeValue('username', e.target.value)} />
            <TextField label={t('workbench.device.connection.password')} size="small" type="password" value={str(connectionConfig, 'password')} onChange={(e) => onChangeValue('password', e.target.value)} />
          </Box>
          <TextField label={t('workbench.device.connection.topics')} size="small" fullWidth multiline minRows={2} value={str(connectionConfig, 'topics')} onChange={(e) => onChangeTopics(e.target.value)} placeholder="factory/line1/temperature" error={Boolean(fieldErrors.topics)} helperText={fieldErrors.topics} />
          <Box sx={row}>
            <TextField label={t('workbench.device.connection.qos')} size="small" select SelectProps={{ native: true }} value={str(connectionConfig, 'qos') || '0'} onChange={(e) => onChangeValue('qos', Number.parseInt(e.target.value, 10))}>
              {[0, 1, 2].map((q) => <option key={q} value={String(q)}>{q}</option>)}
            </TextField>
            <FormControlLabel control={<Checkbox checked={Boolean(connectionConfig.use_tls)} onChange={(e) => onChangeValue('use_tls', e.target.checked)} size="small" />} label={t('workbench.device.connection.useTls')} />
          </Box>
        </>
      );
  }
}
