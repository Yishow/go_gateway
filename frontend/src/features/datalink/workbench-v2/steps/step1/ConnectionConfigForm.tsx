import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { ProtocolId } from '../../state/types';
import { Field } from '../../components/Field';
import { Input } from '../../components/inputs';
import { ConnectionAdvancedFields, FatekModeField, SerialConnectionFields } from './ConnectionProtocolFields';

export interface ConnectionConfigFormProps {
  /** 協議 ID，決定渲染哪種配置佈局 */
  protocol: ProtocolId;
  /** 當前設備的連線配置資料 */
  config: Record<string, unknown>;
  /** 當設定值變更時觸發的 callback 函數 */
  onChange: (patch: Record<string, unknown>) => void;
  /** Source-local unsent buffers retain invalid numeric text until explicit submit. */
  rawNumericText?: boolean;
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
  rawNumericText = false,
}) => {
  const { t } = useTranslation('workbench-v2');

  // TCP 系包括 modbus_tcp, modbus_udp, fatek_fbs, mc_3e
  const isTcpLike =
    protocol === 'modbus_tcp' ||
    protocol === 'modbus_udp' ||
    (protocol === 'fatek_fbs' && (config.mode ?? 'tcp') === 'tcp') ||
    protocol === 'mc_3e';

  // 輔助更新屬性的 handler
  const handleChange = (key: string, value: string | number) => {
    onChange({ [key]: value });
  };

  if (isTcpLike) {
    const isStation = protocol === 'mc_3e' || protocol === 'fatek_fbs';
    // 只顯示設定裡真正存在的值；預設值由裝置載入正規化寫入，避免畫面與有效性判定不一致。
    const stationValue = (config.station_no as number) ?? (config.station as number) ?? (config.slave_id as number) ?? (config.module_station_no as number) ?? '';

    const handleStationChange = (val: string | number) => {
      const num = rawNumericText ? String(val) : val === '' ? '' : Number(val);
      if (isStation) {
        onChange({ station: num, station_no: num, slave_id: num });
      } else {
        onChange({ slave_id: num });
      }
    };

    return (
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        {protocol === 'fatek_fbs' && <FatekModeField config={config} onChange={onChange} />}
        <Field help="connection.host" label={t('step1.fields.host')} required>
          <Input
            type="text"
            value={(config.host as string) ?? ''}
            onChange={(e) => handleChange('host', e.target.value)}
            placeholder="192.168.1.100"
            data-testid="input-host"
          />
        </Field>

        <Field help="connection.port" label={t('step1.fields.port')} required>
          <Input
            type={rawNumericText ? "text" : "number"}
            inputMode={rawNumericText ? "decimal" : undefined}
            value={(config.port as number) ?? ''}
            onChange={(e) => handleChange('port', rawNumericText ? e.target.value : e.target.value === '' ? '' : Number(e.target.value))}
            placeholder={protocol === 'mc_3e' ? '6000' : protocol === 'fatek_fbs' ? '500' : '502'}
            data-testid="input-port"
          />
        </Field>

        <Field help={protocol === 'fatek_fbs' ? 'connection.fatekStation' : protocol === 'mc_3e' ? 'connection.mcStation' : protocol === 'modbus_udp' ? 'connection.udpSlave' : 'connection.modbusSlave'} label={isStation ? t('step1.fields.station', '站號') : t('step1.fields.slave_id', '從站 ID')} required>
          <Input
            type={rawNumericText ? "text" : "number"}
            inputMode={rawNumericText ? "decimal" : undefined}
            value={stationValue}
            onChange={(e) => handleStationChange(e.target.value)}
            placeholder={protocol === 'mc_3e' ? '0' : '1'}
            data-testid="input-slave-id"
          />
        </Field>

        <Field help={protocol === 'modbus_tcp' ? 'connection.timeoutModbus' : 'connection.timeoutOther'} label={t('step1.fields.timeout')} required hint={t('step1.fields.timeout_hint')}>
          <Input
            type={rawNumericText ? "text" : "number"}
            inputMode={rawNumericText ? "decimal" : undefined}
            value={(config.timeout as number) ?? ''}
            onChange={(e) => handleChange('timeout', rawNumericText ? e.target.value : e.target.value === '' ? '' : Number(e.target.value))}
            placeholder="5"
            data-testid="input-timeout"
          />
        </Field>
        <ConnectionAdvancedFields protocol={protocol} config={config} onChange={onChange} rawNumericText={rawNumericText} />
      </div>
    );
  }


  if (protocol === 'modbus_rtu' || protocol === 'fatek_fbs') {
    return <SerialConnectionFields protocol={protocol} config={config} onChange={onChange} rawNumericText={rawNumericText} />;
  }

  if (protocol === 'mqtt') {
    return (
      <div className="grid grid-cols-1 gap-4">
        <Field help="connection.mqttBroker" label={t('step1.fields.mqtt_broker')} required>
          <Input
            type="text"
            value={(config.broker_url as string) ?? (config.broker as string) ?? ''}
            onChange={(e) => handleChange('broker_url', e.target.value)}
            placeholder="mqtts://broker.local:8883"
            data-testid="input-mqtt-broker"
          />
        </Field>

        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <Field help="connection.mqttUsername" label={t('step1.fields.mqtt_username')}>
            <Input
              type="text"
              value={(config.username as string) ?? ''}
              onChange={(e) => handleChange('username', e.target.value)}
              placeholder="admin (選填)"
              data-testid="input-mqtt-username"
            />
          </Field>

          <Field help="connection.mqttClientId" label={t('step1.fields.mqtt_client_id')} required>
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
