import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { ProtocolId } from '../../state/types';
import { Field } from '../../components/Field';
import { Input, Select } from '../../components/inputs';

function displayParity(value: unknown): 'N' | 'E' | 'O' {
  switch (value) {
    case 'even':
    case 'E':
      return 'E';
    case 'odd':
    case 'O':
      return 'O';
    default:
      return 'N';
  }
}

export interface ConnectionConfigFormProps {
  /** 協議 ID，決定渲染哪種配置佈局 */
  protocol: ProtocolId;
  /** 當前設備的連線配置資料 */
  config: Record<string, unknown>;
  /** 當設定值變更時觸發的 callback 函數 */
  onChange: (patch: Record<string, unknown>) => void;
}

/**
 * 協議配置表單元件 (ConnectionConfigForm)
 * 
 * 落地需求 **Connection config form per protocol**。
 * 根據所選協議渲染專屬的參數配置欄位（TCP 系 / RTU 系 / MQTT），並進行型別正確的 onChange 回報。
 */
export const ConnectionConfigForm: React.FC<ConnectionConfigFormProps> = ({
  protocol,
  config,
  onChange,
}) => {
  const { t } = useTranslation('workbench-v2');

  // TCP 系包括 modbus_tcp, modbus_udp, fatek_fbs, mc_3e
  const isTcpLike =
    protocol === 'modbus_tcp' ||
    protocol === 'modbus_udp' ||
    protocol === 'fatek_fbs' ||
    protocol === 'mc_3e';

  // 輔助更新屬性的 handler
  const handleChange = (key: string, value: string | number) => {
    onChange({ [key]: value });
  };

  if (isTcpLike) {
    const isStation = protocol === 'mc_3e' || protocol === 'fatek_fbs';
    // 只顯示設定裡真正存在的值；預設值由裝置載入正規化寫入，避免畫面與有效性判定不一致。
    const stationValue = (config.station_no as number) ?? (config.station as number) ?? (config.slave_id as number) ?? '';

    const handleStationChange = (val: string | number) => {
      const num = val === '' ? '' : Number(val);
      if (isStation) {
        onChange({ station: num, station_no: num, slave_id: num });
      } else {
        onChange({ slave_id: num });
      }
    };

    return (
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <Field label={t('step1.fields.host')} required>
          <Input
            type="text"
            value={(config.host as string) ?? ''}
            onChange={(e) => handleChange('host', e.target.value)}
            placeholder="192.168.1.100"
            data-testid="input-host"
          />
        </Field>

        <Field label={t('step1.fields.port')} required>
          <Input
            type="number"
            value={(config.port as number) ?? ''}
            onChange={(e) => handleChange('port', e.target.value === '' ? '' : Number(e.target.value))}
            placeholder={protocol === 'mc_3e' ? '6000' : protocol === 'fatek_fbs' ? '500' : '502'}
            data-testid="input-port"
          />
        </Field>

        <Field label={isStation ? t('step1.fields.station', '站號') : t('step1.fields.slave_id', '從站 ID')} required>
          <Input
            type="number"
            value={stationValue}
            onChange={(e) => handleStationChange(e.target.value)}
            placeholder={protocol === 'mc_3e' ? '0' : '1'}
            data-testid="input-slave-id"
          />
        </Field>

        <Field label={t('step1.fields.timeout')} required hint={t('step1.fields.timeout_hint')}>
          <Input
            type="number"
            value={(config.timeout as number) ?? ''}
            onChange={(e) => handleChange('timeout', e.target.value === '' ? '' : Number(e.target.value))}
            placeholder="5"
            data-testid="input-timeout"
          />
        </Field>
      </div>
    );
  }


  if (protocol === 'modbus_rtu') {
    return (
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <Field label={t('step1.fields.serial_port')} required>
          <Input
            type="text"
            value={(config.serial_port as string) ?? (config.port as string) ?? ''}
            onChange={(e) => handleChange('serial_port', e.target.value)}
            placeholder="/dev/ttyUSB0"
            data-testid="input-rtu-port"
          />
        </Field>

        <Field label={t('step1.fields.baud_rate')} required>
          <Select
            value={(config.baud_rate as number) ?? (config.baud as number) ?? 9600}
            onChange={(e) => handleChange('baud_rate', Number(e.target.value))}
            data-testid="select-baud"
          >
            <option value={2400}>2400</option>
            <option value={4800}>4800</option>
            <option value={9600}>9600</option>
            <option value={19200}>19200</option>
            <option value={38400}>38400</option>
            <option value={57600}>57600</option>
            <option value={115200}>115200</option>
          </Select>
        </Field>

        <Field label={t('step1.fields.parity')} required>
          <Select
            value={displayParity(config.parity)}
            onChange={(e) => handleChange('parity', e.target.value)}
            data-testid="select-parity"
          >
            <option value="N">{t('step1.parity.none')}</option>
            <option value="E">{t('step1.parity.even')}</option>
            <option value="O">{t('step1.parity.odd')}</option>
          </Select>
        </Field>

        <Field label={t('step1.fields.slave_id')} required>
          <Input
            type="number"
            value={(config.slave_id as number) ?? ''}
            onChange={(e) => handleChange('slave_id', e.target.value === '' ? '' : Number(e.target.value))}
            placeholder="1"
            data-testid="input-rtu-slave-id"
          />
        </Field>
      </div>
    );
  }

  if (protocol === 'mqtt') {
    return (
      <div className="grid grid-cols-1 gap-4">
        <Field label={t('step1.fields.mqtt_broker')} required>
          <Input
            type="text"
            value={(config.broker_url as string) ?? (config.broker as string) ?? ''}
            onChange={(e) => handleChange('broker_url', e.target.value)}
            placeholder="mqtts://broker.local:8883"
            data-testid="input-mqtt-broker"
          />
        </Field>

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <Field label={t('step1.fields.mqtt_username')}>
            <Input
              type="text"
              value={(config.username as string) ?? ''}
              onChange={(e) => handleChange('username', e.target.value)}
              placeholder="admin (選填)"
              data-testid="input-mqtt-username"
            />
          </Field>

          <Field label={t('step1.fields.mqtt_client_id')} required>
            <Input
              type="text"
              value={(config.client_id as string) ?? ''}
              onChange={(e) => handleChange('client_id', e.target.value)}
              placeholder="gw-01"
              data-testid="input-mqtt-client-id"
            />
          </Field>
        </div>
      </div>
    );
  }

  return null;
};
