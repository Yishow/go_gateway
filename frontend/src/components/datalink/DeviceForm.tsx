import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import type { Device, CreateDeviceRequest, UpdateDeviceRequest, ProtocolType, ProtocolInfo } from '../../types/datalink';
import { protocolAPI, settingsAPI } from '../../services/datalink';
import { logger } from '../../utils/logger';

type DeviceConfig = Record<string, string | number | undefined>;

interface DeviceFormProps {
  device?: Device;
  onSubmit: (data: CreateDeviceRequest | UpdateDeviceRequest) => Promise<void>;
  onCancel: () => void;
}

/**
 * 根據協議類型獲取默認配置值
 * 必須在組件外部定義，以便在 useState 初始化時使用
 */
function getDefaultConfigForProtocol(proto: ProtocolType): DeviceConfig {
  switch (proto) {
    case 'modbus_tcp':
      return {
        port: 502,
        slave_id: 1,
        timeout: 5,
      };
    case 'modbus_rtu':
      return {
        baud_rate: 9600,
        data_bits: 8,
        stop_bits: 1,
        parity: 'none',
        slave_id: 1,
        timeout: 5,
      };
    case 'modbus_udp':
      return {
        port: 502,
        slave_id: 1,
        timeout: 5,
      };
    case 'fatek_fbs':
      return {
        mode: 'tcp',
        station_no: 1,
        port: 500,
      };
    case 'mc_3e':
      return {
        port: 5000,
        network_no: 0,
        pc_no: 255,
        io_no: 1023,
        station_no: 0,
      };
    default:
      return {};
  }
}

export default function DeviceForm({ device, onSubmit, onCancel }: DeviceFormProps) {
  const { t } = useTranslation();
  const formIdPrefix = device?.id ? `device-form-${device.id}` : 'device-form-new';
  const fieldId = (name: string) => `${formIdPrefix}-${name}`;
  const [name, setName] = useState(device?.name || '');
  const [description, setDescription] = useState(device?.description || '');
  const [protocol, setProtocol] = useState<ProtocolType>(device?.protocol || 'modbus_tcp');
  const [retryCount, setRetryCount] = useState(3);
  const [retryDelay, setRetryDelay] = useState(1000);
  
  // Dynamic Configuration
  const [config, setConfig] = useState<DeviceConfig>(() => {
    if (device?.connection_config) {
      return typeof device.connection_config === 'string' 
        ? (JSON.parse(device.connection_config) as DeviceConfig)
        : device.connection_config;
    }
    // 新建設備時，根據協議設置默認值
    const initialProtocol = device?.protocol || 'modbus_tcp';
    return getDefaultConfigForProtocol(initialProtocol);
  });

  const [protocols, setProtocols] = useState<ProtocolInfo[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  useEffect(() => {
    fetchProtocols();
    fetchSettings();
  }, []);

  const fetchProtocols = async () => {
    try {
      const list = await protocolAPI.list();
      setProtocols(list);
    } catch (err) {
      logger.error('Failed to load protocols', err);
    }
  };

  const fetchSettings = async () => {
    try {
      const settings = await settingsAPI.get();
      if (settings.default_retry_count !== undefined) {
        setRetryCount(settings.default_retry_count);
      }
      if (settings.default_retry_delay !== undefined) {
        setRetryDelay(settings.default_retry_delay);
      }
    } catch (err) {
      logger.error('Failed to load settings', err);
    }
  };


  /**
   * 合併配置，確保必填欄位有默認值
   * 先設置默認值，然後用用戶輸入的值覆蓋（只要不是 undefined）
   */
  const prepareConfig = (): DeviceConfig => {
    const defaultConfig = getDefaultConfigForProtocol(protocol);
    const finalConfig: DeviceConfig = { ...defaultConfig };
    
    // 用用戶輸入的值覆蓋默認值（跳過 undefined，但保留 0、false、空字串等）
    for (const [key, value] of Object.entries(config)) {
      if (value !== undefined) {
        finalConfig[key] = value;
      }
    }
    
    return finalConfig;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      // 準備配置，確保必填欄位有值
      const finalConfig = prepareConfig();
      
      if (device) {
        await onSubmit({
          name,
          description,
          connection_config: finalConfig,
        } as UpdateDeviceRequest);
      } else {
        await onSubmit({
          name,
          description,
          protocol,
          connection_config: finalConfig,
        } as CreateDeviceRequest);
      }
    } catch (err: unknown) {
      const maybeError = err as {
        response?: { data?: { error?: { message?: string }; message?: string } };
        message?: string;
      };
      const errorMessage =
        maybeError.response?.data?.error?.message
        || maybeError.response?.data?.message
        || maybeError.message
        || t('device.failedToSave');
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  // 簡單的動態配置表單，未來可優化為根據 config_schema 生成
  const renderConfigFields = () => {
    switch (protocol) {
      case 'modbus_tcp':
        return (
          <>
            <div>
              <label htmlFor={fieldId('modbus-tcp-host')} className="block text-sm font-medium text-slate-300">{t('device.host')}</label>
              <input
                id={fieldId('modbus-tcp-host')}
                name="host"
                type="text"
                value={config.host ?? ''}
                onChange={e => setConfig({...config, host: e.target.value})}
                placeholder="192.168.1.100"
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                required
              />
            </div>
            <div>
              <label htmlFor={fieldId('modbus-tcp-port')} className="block text-sm font-medium text-slate-300">{t('device.port')}</label>
              <input
                id={fieldId('modbus-tcp-port')}
                name="port"
                type="number"
                value={config.port ?? 502}
                onChange={e => {
                  const value = e.target.value;
                  // 空值時設置為 undefined，否則轉換為數字
                  setConfig({...config, port: value === '' ? undefined : Number(value)});
                }}
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                required
                min={1}
                max={65535}
              />
            </div>
            <div>
              <label htmlFor={fieldId('modbus-tcp-slave-id')} className="block text-sm font-medium text-slate-300">{t('device.slaveId')}</label>
              <input
                id={fieldId('modbus-tcp-slave-id')}
                name="slave_id"
                type="number"
                value={config.slave_id ?? 1}
                onChange={e => {
                  const value = e.target.value;
                  setConfig({...config, slave_id: value === '' ? undefined : Number(value)});
                }}
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                required
                min={1}
                max={247}
              />
            </div>
          </>
        );
      case 'modbus_rtu':
           return (
             <>
               <div>
                <label htmlFor={fieldId('modbus-rtu-serial-port')} className="block text-sm font-medium text-slate-300">{t('device.serialPort')}</label>
                <input
                  id={fieldId('modbus-rtu-serial-port')}
                  name="serial_port"
                  type="text"
                  value={config.serial_port ?? ''}
                  onChange={e => setConfig({...config, serial_port: e.target.value})}
                  placeholder="COM1 or /dev/ttyUSB0"
                  className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                  required
                />
              </div>
              <div>
                <label htmlFor={fieldId('modbus-rtu-baud-rate')} className="block text-sm font-medium text-slate-300">{t('device.baudRate')}</label>
                <input
                  id={fieldId('modbus-rtu-baud-rate')}
                  name="baud_rate"
                  type="number"
                  value={config.baud_rate ?? 9600}
                  onChange={e => {
                    const value = e.target.value;
                    setConfig({...config, baud_rate: value === '' ? undefined : Number(value)});
                  }}
                  className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                  required
                />
              </div>
               <div>
                  <label htmlFor={fieldId('modbus-rtu-data-bits')} className="block text-sm font-medium text-slate-300">{t('device.dataBits')}</label>
                  <input
                    id={fieldId('modbus-rtu-data-bits')}
                    name="data_bits"
                    type="number"
                    value={config.data_bits ?? 8}
                    onChange={e => {
                      const value = e.target.value;
                      setConfig({...config, data_bits: value === '' ? undefined : Number(value)});
                    }}
                    className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                    required
                  />
                </div>
                 <div>
                  <label htmlFor={fieldId('modbus-rtu-stop-bits')} className="block text-sm font-medium text-slate-300">Stop Bits</label>
                  <input
                    id={fieldId('modbus-rtu-stop-bits')}
                    name="stop_bits"
                    type="number"
                    value={config.stop_bits ?? 1}
                    onChange={e => {
                      const value = e.target.value;
                      setConfig({...config, stop_bits: value === '' ? undefined : Number(value)});
                    }}
                    className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                    required
                  />
                </div>
                 <div>
                  <label htmlFor={fieldId('modbus-rtu-parity')} className="block text-sm font-medium text-slate-300">{t('device.parity')}</label>
                  <select
                     id={fieldId('modbus-rtu-parity')}
                     name="parity"
                     value={config.parity || 'none'}
                     onChange={e => setConfig({...config, parity: e.target.value})}
                     className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                  >
                      <option value="none">{t('device.parityNone')}</option>
                      <option value="even">{t('device.parityEven')}</option>
                      <option value="odd">{t('device.parityOdd')}</option>
                  </select>
                </div>
              <div>
                <label htmlFor={fieldId('modbus-rtu-slave-id')} className="block text-sm font-medium text-slate-300">{t('device.slaveId')}</label>
                <input
                  id={fieldId('modbus-rtu-slave-id')}
                  name="slave_id"
                  type="number"
                  value={config.slave_id ?? 1}
                  onChange={e => {
                    const value = e.target.value;
                    setConfig({...config, slave_id: value === '' ? undefined : Number(value)});
                  }}
                  className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                  required
                  min={1}
                  max={247}
                />
              </div>
            </>
          );
      case 'modbus_udp':
          // Similar to TCP
         return (
          <>
            <div>
              <label htmlFor={fieldId('modbus-udp-host')} className="block text-sm font-medium text-slate-300">{t('device.host')}</label>
              <input
                id={fieldId('modbus-udp-host')}
                name="host"
                type="text"
                value={config.host ?? ''}
                onChange={e => setConfig({...config, host: e.target.value})}
                placeholder="192.168.1.100"
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                required
              />
            </div>
            <div>
              <label htmlFor={fieldId('modbus-udp-port')} className="block text-sm font-medium text-slate-300">{t('device.port')}</label>
              <input
                id={fieldId('modbus-udp-port')}
                name="port"
                type="number"
                value={config.port ?? 502}
                onChange={e => {
                  const value = e.target.value;
                  setConfig({...config, port: value === '' ? undefined : Number(value)});
                }}
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                required
                min={1}
                max={65535}
              />
            </div>
             <div>
                <label htmlFor={fieldId('modbus-udp-slave-id')} className="block text-sm font-medium text-slate-300">{t('device.slaveId')}</label>
                <input
                  id={fieldId('modbus-udp-slave-id')}
                  name="slave_id"
                  type="number"
                  value={config.slave_id ?? 1}
                  onChange={e => {
                    const value = e.target.value;
                    setConfig({...config, slave_id: value === '' ? undefined : Number(value)});
                  }}
                  className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                  required
                  min={1}
                  max={247}
                />
              </div>
          </>
        );
      case 'fatek_fbs':
        return (
          <>
            <div>
              <label htmlFor={fieldId('fatek-mode')} className="block text-sm font-medium text-slate-300">{t('device.mode')}</label>
              <select
                id={fieldId('fatek-mode')}
                name="mode"
                value={config.mode || 'tcp'}
                onChange={e => setConfig({...config, mode: e.target.value})}
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
              >
                <option value="tcp">{t('device.modeTcp')}</option>
                <option value="serial">{t('device.modeSerial')}</option>
              </select>
            </div>
            <div>
              <label htmlFor={fieldId('fatek-station-no')} className="block text-sm font-medium text-slate-300">{t('device.stationNo')}</label>
              <input
                id={fieldId('fatek-station-no')}
                name="station_no"
                type="number"
                value={config.station_no ?? 1}
                onChange={e => {
                  const value = e.target.value;
                  setConfig({...config, station_no: value === '' ? undefined : Number(value)});
                }}
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                required
                min={0}
                max={255}
              />
            </div>
            {config.mode !== 'serial' ? (
              <>
                <div>
                  <label htmlFor={fieldId('fatek-host')} className="block text-sm font-medium text-slate-300">{t('device.host')}</label>
                  <input
                    id={fieldId('fatek-host')}
                    name="host"
                    type="text"
                    value={config.host ?? ''}
                    onChange={e => setConfig({...config, host: e.target.value})}
                    placeholder="192.168.1.100"
                    className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                    required
                  />
                </div>
                <div>
                  <label htmlFor={fieldId('fatek-port')} className="block text-sm font-medium text-slate-300">{t('device.port')}</label>
                  <input
                    id={fieldId('fatek-port')}
                    name="port"
                    type="number"
                    value={config.port ?? 500}
                    onChange={e => {
                      const value = e.target.value;
                      setConfig({...config, port: value === '' ? undefined : Number(value)});
                    }}
                    className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                    required
                    min={1}
                    max={65535}
                  />
                </div>
              </>
            ) : (
              <>
                <div>
                  <label htmlFor={fieldId('fatek-serial-port')} className="block text-sm font-medium text-slate-300">{t('device.serialPort')}</label>
                  <input
                    id={fieldId('fatek-serial-port')}
                    name="serial_port"
                    type="text"
                    value={config.serial_port ?? ''}
                    onChange={e => setConfig({...config, serial_port: e.target.value})}
                    placeholder={t('device.serialPortPlaceholder')}
                    className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                    required
                  />
                </div>
                <div>
                  <label htmlFor={fieldId('fatek-baud-rate')} className="block text-sm font-medium text-slate-300">{t('device.baudRate')}</label>
                  <input
                    id={fieldId('fatek-baud-rate')}
                    name="baud_rate"
                    type="number"
                    value={config.baud_rate ?? 9600}
                    onChange={e => {
                      const value = e.target.value;
                      setConfig({...config, baud_rate: value === '' ? undefined : Number(value)});
                    }}
                    className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                    required
                  />
                </div>
              </>
            )}
          </>
        );
      case 'mc_3e':
        return (
          <>
            <div>
              <label htmlFor={fieldId('mc3e-host')} className="block text-sm font-medium text-slate-300">{t('device.host')}</label>
              <input
                id={fieldId('mc3e-host')}
                name="host"
                type="text"
                value={config.host ?? ''}
                onChange={e => setConfig({...config, host: e.target.value})}
                placeholder="192.168.1.100"
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                required
              />
            </div>
            <div>
              <label htmlFor={fieldId('mc3e-port')} className="block text-sm font-medium text-slate-300">{t('device.port')}</label>
              <input
                id={fieldId('mc3e-port')}
                name="port"
                type="number"
                value={config.port ?? 5000}
                onChange={e => {
                  const value = e.target.value;
                  setConfig({...config, port: value === '' ? undefined : Number(value)});
                }}
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                required
                min={1}
                max={65535}
              />
            </div>
            <div>
              <label htmlFor={fieldId('mc3e-network-no')} className="block text-sm font-medium text-slate-300">Network No</label>
              <input
                id={fieldId('mc3e-network-no')}
                name="network_no"
                type="number"
                value={config.network_no ?? 0}
                onChange={e => {
                  const value = e.target.value;
                  setConfig({...config, network_no: value === '' ? undefined : Number(value)});
                }}
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                min={0}
              />
            </div>
            <div>
              <label htmlFor={fieldId('mc3e-pc-no')} className="block text-sm font-medium text-slate-300">{t('device.pcNo')}</label>
              <input
                id={fieldId('mc3e-pc-no')}
                name="pc_no"
                type="number"
                value={config.pc_no ?? 255}
                onChange={e => {
                  const value = e.target.value;
                  setConfig({...config, pc_no: value === '' ? undefined : Number(value)});
                }}
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                min={0}
                max={255}
              />
            </div>
            <div>
              <label htmlFor={fieldId('mc3e-io-no')} className="block text-sm font-medium text-slate-300">I/O No</label>
              <input
                id={fieldId('mc3e-io-no')}
                name="io_no"
                type="number"
                value={config.io_no ?? 1023}
                onChange={e => {
                  const value = e.target.value;
                  setConfig({...config, io_no: value === '' ? undefined : Number(value)});
                }}
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                min={0}
              />
            </div>
            <div>
              <label htmlFor={fieldId('mc3e-station-no')} className="block text-sm font-medium text-slate-300">{t('device.stationNo')}</label>
              <input
                id={fieldId('mc3e-station-no')}
                name="station_no"
                type="number"
                value={config.station_no ?? 0}
                onChange={e => {
                  const value = e.target.value;
                  setConfig({...config, station_no: value === '' ? undefined : Number(value)});
                }}
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                min={0}
                max={255}
              />
            </div>
          </>
        );
      
      // Implement other protocols fields...
      default:
        // Generic JSON editor fallback
        return (
          <div className="col-span-2">
            <label htmlFor={fieldId('config-json')} className="block text-sm font-medium text-slate-300">{t('device.configJson')}</label>
            <textarea
              id={fieldId('config-json')}
              name="config_json"
              value={JSON.stringify(config, null, 2)}
              onChange={e => {
                  try {
                      setConfig(JSON.parse(e.target.value));
                  } catch {
                      // ignore parse error while typing
                  }
              }}
              rows={5}
              className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 font-mono text-xs focus:outline-none focus:border-blue-500"
            />
          </div>
        );
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <div className="bg-red-500/10 border border-red-500/20 rounded-lg p-3 text-sm text-red-400">
          {error}
        </div>
      )}

      {/* Basic Info */}
      <div className="grid grid-cols-2 gap-4">
        <div className="col-span-2">
          <label htmlFor={fieldId('device-name')} className="block text-sm font-medium text-slate-300">{t('device.name')}</label>
          <input
            id={fieldId('device-name')}
            name="name"
            autoComplete="off"
            type="text"
            value={name}
            onChange={e => setName(e.target.value)}
            className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
            required
            placeholder="e.g. Production PLC 01"
          />
        </div>
        <div className="col-span-2">
          <label htmlFor={fieldId('device-description')} className="block text-sm font-medium text-slate-300">{t('device.description')}</label>
          <input
            id={fieldId('device-description')}
            name="description"
            type="text"
            value={description}
            onChange={e => setDescription(e.target.value)}
            className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
            placeholder={t('device.descriptionPlaceholder')}
          />
        </div>
      </div>

      <div className="border-t border-slate-700 pt-4">
          <h4 className="text-sm font-semibold text-slate-400 mb-3">{t('device.connectionSettings')}</h4>
          <div className="grid grid-cols-2 gap-4">
               {!device && (
                   <div className="col-span-2">
                    <label htmlFor={fieldId('device-protocol')} className="block text-sm font-medium text-slate-300">{t('device.protocol')}</label>
                    <select
                        id={fieldId('device-protocol')}
                        name="protocol"
                        value={protocol}
                        onChange={e => {
                            const newProtocol = e.target.value as ProtocolType;
                            setProtocol(newProtocol);
                            // 重置配置並設置默認值
                            setConfig(getDefaultConfigForProtocol(newProtocol));
                        }}
                        className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                    >
                        {protocols.map(p => (
                            <option key={p.type} value={p.type}>{p.name}</option>
                        ))}
                    </select>
                   </div>
               )}
               
               {renderConfigFields()}

                <div className="border-t border-slate-700 col-span-2 pt-4 mt-2">
                    <h4 className="text-sm font-semibold text-slate-400 mb-3">Advanced</h4>
                     <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label htmlFor={fieldId('retry-count')} className="block text-sm font-medium text-slate-300">{t('device.retryCount')}</label>
                            <input
                                id={fieldId('retry-count')}
                                name="retry_count"
                                type="number"
                                value={retryCount}
                                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                                disabled
                                readOnly
                            />
                        </div>
                         <div>
                            <label htmlFor={fieldId('retry-delay')} className="block text-sm font-medium text-slate-300">Retry Delay (ms)</label>
                            <input
                                id={fieldId('retry-delay')}
                                name="retry_delay"
                                type="number"
                                value={retryDelay}
                                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                                disabled
                                readOnly
                            />
                        </div>
                     </div>
                     <p className="mt-2 text-xs text-slate-500">
                        {t('device.retryValuesManaged')}
                     </p>
                </div>

          </div>
      </div>

      <div className="flex justify-end space-x-3 pt-6 border-t border-slate-700 mt-6">
        <button
          type="button"
          onClick={onCancel}
          className="px-4 py-2 text-slate-300 hover:text-white transition-colors"
          disabled={loading}
        >
          Cancel
        </button>
        <button
          type="submit"
          className="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          disabled={loading}
        >
          {loading ? t('common.saving') : t('device.saveDevice')}
        </button>
      </div>
    </form>
  );
}
