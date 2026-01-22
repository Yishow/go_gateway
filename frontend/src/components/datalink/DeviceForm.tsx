import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import type { Device, CreateDeviceRequest, UpdateDeviceRequest, ProtocolType } from '../../types/datalink';
import { protocolAPI, settingsAPI } from '../../services/datalink';

interface DeviceFormProps {
  device?: Device;
  onSubmit: (data: CreateDeviceRequest | UpdateDeviceRequest) => Promise<void>;
  onCancel: () => void;
}

/**
 * 根據協議類型獲取默認配置值
 * 必須在組件外部定義，以便在 useState 初始化時使用
 */
function getDefaultConfigForProtocol(proto: ProtocolType): Record<string, any> {
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
  const [name, setName] = useState(device?.name || '');
  const [description, setDescription] = useState(device?.description || '');
  const [protocol, setProtocol] = useState<ProtocolType>(device?.protocol || 'modbus_tcp');
  const [retryCount, setRetryCount] = useState(3);
  const [retryDelay, setRetryDelay] = useState(1000);
  
  // Dynamic Configuration
  const [config, setConfig] = useState<Record<string, any>>(() => {
    if (device?.connection_config) {
      return typeof device.connection_config === 'string' 
        ? JSON.parse(device.connection_config) 
        : device.connection_config;
    }
    // 新建設備時，根據協議設置默認值
    const initialProtocol = device?.protocol || 'modbus_tcp';
    return getDefaultConfigForProtocol(initialProtocol);
  });

  const [protocols, setProtocols] = useState<any[]>([]);
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
      console.error('Failed to load protocols', err);
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
      console.error('Failed to load settings', err);
    }
  };


  /**
   * 合併配置，確保必填欄位有默認值
   * 先設置默認值，然後用用戶輸入的值覆蓋（只要不是 undefined）
   */
  const prepareConfig = (): Record<string, any> => {
    const defaultConfig = getDefaultConfigForProtocol(protocol);
    const finalConfig: Record<string, any> = { ...defaultConfig };
    
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
    } catch (err: any) {
      // 提取錯誤訊息
      const errorMessage = err?.response?.data?.error?.message 
        || err?.response?.data?.message 
        || err?.message 
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
              <label className="block text-sm font-medium text-slate-300">{t('device.host')}</label>
              <input
                type="text"
                value={config.host ?? ''}
                onChange={e => setConfig({...config, host: e.target.value})}
                placeholder="192.168.1.100"
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-300">{t('device.port')}</label>
              <input
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
              <label className="block text-sm font-medium text-slate-300">{t('device.slaveId')}</label>
              <input
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
                <label className="block text-sm font-medium text-slate-300">{t('device.serialPort')}</label>
                <input
                  type="text"
                  value={config.serial_port ?? ''}
                  onChange={e => setConfig({...config, serial_port: e.target.value})}
                  placeholder="COM1 or /dev/ttyUSB0"
                  className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-300">{t('device.baudRate')}</label>
                <input
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
                  <label className="block text-sm font-medium text-slate-300">{t('device.dataBits')}</label>
                  <input
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
                  <label className="block text-sm font-medium text-slate-300">Stop Bits</label>
                  <input
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
                  <label className="block text-sm font-medium text-slate-300">{t('device.parity')}</label>
                  <select
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
                <label className="block text-sm font-medium text-slate-300">{t('device.slaveId')}</label>
                <input
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
              <label className="block text-sm font-medium text-slate-300">{t('device.host')}</label>
              <input
                type="text"
                value={config.host ?? ''}
                onChange={e => setConfig({...config, host: e.target.value})}
                placeholder="192.168.1.100"
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-300">{t('device.port')}</label>
              <input
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
                <label className="block text-sm font-medium text-slate-300">{t('device.slaveId')}</label>
                <input
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
              <label className="block text-sm font-medium text-slate-300">{t('device.mode')}</label>
              <select
                value={config.mode || 'tcp'}
                onChange={e => setConfig({...config, mode: e.target.value})}
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
              >
                <option value="tcp">{t('device.modeTcp')}</option>
                <option value="serial">{t('device.modeSerial')}</option>
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-300">{t('device.stationNo')}</label>
              <input
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
                  <label className="block text-sm font-medium text-slate-300">{t('device.host')}</label>
                  <input
                    type="text"
                    value={config.host ?? ''}
                    onChange={e => setConfig({...config, host: e.target.value})}
                    placeholder="192.168.1.100"
                    className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-300">{t('device.port')}</label>
                  <input
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
                  <label className="block text-sm font-medium text-slate-300">{t('device.serialPort')}</label>
                  <input
                    type="text"
                    value={config.serial_port ?? ''}
                    onChange={e => setConfig({...config, serial_port: e.target.value})}
                    placeholder={t('device.serialPortPlaceholder')}
                    className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-300">{t('device.baudRate')}</label>
                  <input
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
              <label className="block text-sm font-medium text-slate-300">{t('device.host')}</label>
              <input
                type="text"
                value={config.host ?? ''}
                onChange={e => setConfig({...config, host: e.target.value})}
                placeholder="192.168.1.100"
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-300">{t('device.port')}</label>
              <input
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
              <label className="block text-sm font-medium text-slate-300">Network No</label>
              <input
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
              <label className="block text-sm font-medium text-slate-300">{t('device.pcNo')}</label>
              <input
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
              <label className="block text-sm font-medium text-slate-300">I/O No</label>
              <input
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
              <label className="block text-sm font-medium text-slate-300">{t('device.stationNo')}</label>
              <input
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
            <label className="block text-sm font-medium text-slate-300">{t('device.configJson')}</label>
            <textarea
              value={JSON.stringify(config, null, 2)}
              onChange={e => {
                  try {
                      setConfig(JSON.parse(e.target.value));
                  } catch (er) {
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
          <label className="block text-sm font-medium text-slate-300">{t('device.name')}</label>
          <input
            type="text"
            value={name}
            onChange={e => setName(e.target.value)}
            className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
            required
            placeholder="e.g. Production PLC 01"
          />
        </div>
        <div className="col-span-2">
          <label className="block text-sm font-medium text-slate-300">{t('device.description')}</label>
          <input
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
                    <label className="block text-sm font-medium text-slate-300">{t('device.protocol')}</label>
                    <select
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
                            <label className="block text-sm font-medium text-slate-300">{t('device.retryCount')}</label>
                            <input
                                type="number"
                                value={retryCount}
                                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                                disabled
                                readOnly
                            />
                        </div>
                         <div>
                            <label className="block text-sm font-medium text-slate-300">Retry Delay (ms)</label>
                            <input
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
