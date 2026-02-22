import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import type { Device, CreateDeviceRequest, UpdateDeviceRequest, ProtocolType, ProtocolInfo } from '../../types/datalink';
import { protocolAPI, settingsAPI } from '../../services/datalink';
import { logger } from '../../utils/logger';

type DeviceConfig = Record<string, string | number | boolean | string[] | undefined>;

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
      return { port: 502, slave_id: 1, timeout: 5 };
    case 'modbus_rtu':
      return { baud_rate: 9600, data_bits: 8, stop_bits: 1, parity: 'none', slave_id: 1, timeout: 5 };
    case 'modbus_udp':
      return { port: 502, slave_id: 1, timeout: 5 };
    case 'fatek_fbs':
      return { mode: 'tcp', port: 500, station_no: 1, timeout: 5 };
    case 'mc_3e':
      return { port: 3000, network_no: 0, pc_no: 255, io_no: 1023, station_no: 0, timeout: 5, data_format: 'CDAB' };
    case 'mqtt':
      return { qos: 0, use_tls: false };
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

  /**
   * 依協議類型渲染連線設定欄位，與新增設備精靈的 DeviceConnectionStep 保持一致
   */
  const renderConfigFields = () => {
    /** 共用 input className */
    const inp = 'mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500';
    /** 共用 label className */
    const lbl = 'block text-sm font-medium text-slate-300';
    /** 數字欄位 onChange helper */
    const setNum = (key: string) => (e: React.ChangeEvent<HTMLInputElement>) => {
      const v = e.target.value;
      setConfig({ ...config, [key]: v === '' ? undefined : Number(v) });
    };

    switch (protocol) {
      // ─── Modbus TCP ───────────────────────────────────────────────────────────
      case 'modbus_tcp':
        return (
          <>
            <div className="col-span-2 grid grid-cols-2 gap-4">
              <div>
                <label htmlFor={fieldId('host')} className={lbl}>主機位址 (Host)</label>
                <input id={fieldId('host')} type="text" value={String(config.host ?? '')}
                  onChange={e => setConfig({ ...config, host: e.target.value })}
                  className={inp} placeholder="192.168.1.10" />
              </div>
              <div>
                <label htmlFor={fieldId('port')} className={lbl}>連接埠 (Port)</label>
                <input id={fieldId('port')} type="number" value={config.port ?? 502}
                  onChange={setNum('port')} className={inp} min={1} max={65535} />
              </div>
            </div>
            <div className="col-span-2 grid grid-cols-2 gap-4">
              <div>
                <label htmlFor={fieldId('slave-id')} className={lbl}>從站 ID (Slave ID)</label>
                <input id={fieldId('slave-id')} type="number" value={config.slave_id ?? 1}
                  onChange={setNum('slave_id')} className={inp} min={1} max={247} />
              </div>
              <div>
                <label htmlFor={fieldId('timeout')} className={lbl}>逾時秒數 (Timeout)</label>
                <input id={fieldId('timeout')} type="number" value={config.timeout ?? 5}
                  onChange={setNum('timeout')} className={inp} min={1} />
              </div>
            </div>
          </>
        );

      // ─── Modbus UDP ───────────────────────────────────────────────────────────
      case 'modbus_udp':
        return (
          <>
            <div className="col-span-2 grid grid-cols-2 gap-4">
              <div>
                <label htmlFor={fieldId('host')} className={lbl}>主機位址 (Host)</label>
                <input id={fieldId('host')} type="text" value={String(config.host ?? '')}
                  onChange={e => setConfig({ ...config, host: e.target.value })}
                  className={inp} placeholder="192.168.1.10" />
              </div>
              <div>
                <label htmlFor={fieldId('port')} className={lbl}>連接埠 (Port)</label>
                <input id={fieldId('port')} type="number" value={config.port ?? 502}
                  onChange={setNum('port')} className={inp} min={1} max={65535} />
              </div>
            </div>
            <div className="col-span-2 grid grid-cols-2 gap-4">
              <div>
                <label htmlFor={fieldId('slave-id')} className={lbl}>從站 ID (Slave ID)</label>
                <input id={fieldId('slave-id')} type="number" value={config.slave_id ?? 1}
                  onChange={setNum('slave_id')} className={inp} min={1} max={247} />
              </div>
              <div>
                <label htmlFor={fieldId('timeout')} className={lbl}>逾時秒數 (Timeout)</label>
                <input id={fieldId('timeout')} type="number" value={config.timeout ?? 5}
                  onChange={setNum('timeout')} className={inp} min={1} />
              </div>
            </div>
          </>
        );

      // ─── Modbus RTU ───────────────────────────────────────────────────────────
      case 'modbus_rtu':
        return (
          <>
            <div className="col-span-2">
              <label htmlFor={fieldId('serial-port')} className={lbl}>串列埠 (Serial Port)</label>
              <input id={fieldId('serial-port')} type="text" value={String(config.serial_port ?? '')}
                onChange={e => setConfig({ ...config, serial_port: e.target.value })}
                className={inp} placeholder="COM1 或 /dev/ttyUSB0" />
            </div>
            <div className="col-span-2 grid grid-cols-2 gap-4">
              <div>
                <label htmlFor={fieldId('baud-rate')} className={lbl}>波特率 (Baud Rate)</label>
                <select id={fieldId('baud-rate')} value={Number(config.baud_rate ?? 9600)}
                  onChange={e => setConfig({ ...config, baud_rate: parseInt(e.target.value) })}
                  className={inp}>
                  {[9600, 19200, 38400, 57600, 115200].map(br => (
                    <option key={br} value={br}>{br}</option>
                  ))}
                </select>
              </div>
              <div>
                <label htmlFor={fieldId('data-bits')} className={lbl}>資料位元 (Data Bits)</label>
                <select id={fieldId('data-bits')} value={Number(config.data_bits ?? 8)}
                  onChange={e => setConfig({ ...config, data_bits: parseInt(e.target.value) })}
                  className={inp}>
                  <option value={7}>7</option>
                  <option value={8}>8</option>
                </select>
              </div>
              <div>
                <label htmlFor={fieldId('stop-bits')} className={lbl}>停止位元 (Stop Bits)</label>
                <select id={fieldId('stop-bits')} value={Number(config.stop_bits ?? 1)}
                  onChange={e => setConfig({ ...config, stop_bits: parseInt(e.target.value) })}
                  className={inp}>
                  <option value={1}>1</option>
                  <option value={2}>2</option>
                </select>
              </div>
              <div>
                <label htmlFor={fieldId('parity')} className={lbl}>同位檢查 (Parity)</label>
                <select id={fieldId('parity')} value={String(config.parity ?? 'none')}
                  onChange={e => setConfig({ ...config, parity: e.target.value })}
                  className={inp}>
                  <option value="none">None（無）</option>
                  <option value="even">Even（偶）</option>
                  <option value="odd">Odd（奇）</option>
                </select>
              </div>
            </div>
            <div className="col-span-2 grid grid-cols-2 gap-4">
              <div>
                <label htmlFor={fieldId('slave-id')} className={lbl}>從站 ID (Slave ID)</label>
                <input id={fieldId('slave-id')} type="number" value={config.slave_id ?? 1}
                  onChange={setNum('slave_id')} className={inp} min={1} max={247} />
              </div>
              <div>
                <label htmlFor={fieldId('timeout')} className={lbl}>逾時秒數 (Timeout)</label>
                <input id={fieldId('timeout')} type="number" value={config.timeout ?? 5}
                  onChange={setNum('timeout')} className={inp} min={1} />
              </div>
            </div>
          </>
        );

      // ─── FATEK FBs ────────────────────────────────────────────────────────────
      case 'fatek_fbs': {
        const fatekMode = String(config.mode ?? 'tcp');
        return (
          <>
            <div className="col-span-2">
              <p className={lbl}>連線模式</p>
              <div className="mt-1 flex flex-wrap gap-4">
                {(['tcp', 'serial'] as const).map(m => (
                  <label key={m} htmlFor={fieldId(`mode-${m}`)} className="flex cursor-pointer items-center gap-2">
                    <input id={fieldId(`mode-${m}`)} type="radio" name={fieldId('mode')}
                      value={m} checked={fatekMode === m}
                      onChange={() => setConfig({ ...config, mode: m })} />
                    <span className="text-sm text-slate-300">{m === 'tcp' ? 'TCP/IP' : 'Serial (RS232/485)'}</span>
                  </label>
                ))}
              </div>
            </div>
            {fatekMode === 'serial' ? (
              <div className="col-span-2 grid grid-cols-2 gap-4">
                <div className="col-span-2">
                  <label htmlFor={fieldId('serial-port')} className={lbl}>串列埠 (Serial Port)</label>
                  <input id={fieldId('serial-port')} type="text" value={String(config.serial_port ?? '')}
                    onChange={e => setConfig({ ...config, serial_port: e.target.value })}
                    className={inp} placeholder="COM1 或 /dev/ttyUSB0" />
                </div>
                <div>
                  <label htmlFor={fieldId('baud-rate')} className={lbl}>波特率 (Baud Rate)</label>
                  <select id={fieldId('baud-rate')} value={Number(config.baud_rate ?? 9600)}
                    onChange={e => setConfig({ ...config, baud_rate: parseInt(e.target.value) })}
                    className={inp}>
                    {[9600, 19200, 38400, 57600, 115200].map(br => (
                      <option key={br} value={br}>{br}</option>
                    ))}
                  </select>
                </div>
              </div>
            ) : (
              <div className="col-span-2 grid grid-cols-2 gap-4">
                <div>
                  <label htmlFor={fieldId('host')} className={lbl}>主機位址 (Host)</label>
                  <input id={fieldId('host')} type="text" value={String(config.host ?? '')}
                    onChange={e => setConfig({ ...config, host: e.target.value })}
                    className={inp} placeholder="192.168.1.1" />
                </div>
                <div>
                  <label htmlFor={fieldId('port')} className={lbl}>連接埠 (Port)</label>
                  <input id={fieldId('port')} type="number" value={config.port ?? 500}
                    onChange={setNum('port')} className={inp} min={1} max={65535} />
                </div>
              </div>
            )}
            <div className="col-span-2 grid grid-cols-2 gap-4">
              <div>
                <label htmlFor={fieldId('station-no')} className={lbl}>站號 (Station No)</label>
                <input id={fieldId('station-no')} type="number" value={config.station_no ?? 1}
                  onChange={setNum('station_no')} className={inp} min={0} max={255} />
              </div>
              <div>
                <label htmlFor={fieldId('timeout')} className={lbl}>逾時秒數 (Timeout)</label>
                <input id={fieldId('timeout')} type="number" value={config.timeout ?? 5}
                  onChange={setNum('timeout')} className={inp} min={1} />
              </div>
            </div>
          </>
        );
      }

      // ─── Mitsubishi MC 3E ─────────────────────────────────────────────────────
      case 'mc_3e':
        return (
          <>
            <div className="col-span-2 grid grid-cols-2 gap-4">
              <div>
                <label htmlFor={fieldId('host')} className={lbl}>主機位址 (Host)</label>
                <input id={fieldId('host')} type="text" value={String(config.host ?? '')}
                  onChange={e => setConfig({ ...config, host: e.target.value })}
                  className={inp} placeholder="192.168.1.10" />
              </div>
              <div>
                <label htmlFor={fieldId('port')} className={lbl}>連接埠 (Port)</label>
                <input id={fieldId('port')} type="number" value={config.port ?? 3000}
                  onChange={setNum('port')} className={inp} min={1} max={65535} />
              </div>
            </div>
            <div className="col-span-2 grid grid-cols-2 gap-4">
              <div>
                <label htmlFor={fieldId('network-no')} className={lbl}>網路編號 (Network No)</label>
                <input id={fieldId('network-no')} type="number" value={config.network_no ?? 0}
                  onChange={setNum('network_no')} className={inp} min={0} />
              </div>
              <div>
                <label htmlFor={fieldId('pc-no')} className={lbl}>PC 編號 (PC No)</label>
                <input id={fieldId('pc-no')} type="number" value={config.pc_no ?? 255}
                  onChange={setNum('pc_no')} className={inp} min={0} max={255} />
              </div>
              <div>
                <label htmlFor={fieldId('io-no')} className={lbl}>I/O 編號 (IO No)</label>
                <input id={fieldId('io-no')} type="number" value={config.io_no ?? 1023}
                  onChange={setNum('io_no')} className={inp} min={0} />
              </div>
              <div>
                <label htmlFor={fieldId('station-no')} className={lbl}>站號 (Station No)</label>
                <input id={fieldId('station-no')} type="number" value={config.station_no ?? 0}
                  onChange={setNum('station_no')} className={inp} min={0} max={255} />
              </div>
            </div>
            <div className="col-span-2 grid grid-cols-2 gap-4">
              <div>
                <label htmlFor={fieldId('data-format')} className={lbl}>字節序格式 (Data Format)</label>
                <select id={fieldId('data-format')} value={String(config.data_format ?? 'CDAB')}
                  onChange={e => setConfig({ ...config, data_format: e.target.value })}
                  className={inp}>
                  <option value="CDAB">CDAB（三菱標準）</option>
                  <option value="ABCD">ABCD</option>
                  <option value="BADC">BADC</option>
                  <option value="DCBA">DCBA</option>
                </select>
              </div>
              <div>
                <label htmlFor={fieldId('timeout')} className={lbl}>逾時秒數 (Timeout)</label>
                <input id={fieldId('timeout')} type="number" value={config.timeout ?? 5}
                  onChange={setNum('timeout')} className={inp} min={1} />
              </div>
            </div>
          </>
        );

      // ─── MQTT ─────────────────────────────────────────────────────────────────
      case 'mqtt':
        return (
          <>
            <div className="col-span-2">
              <label htmlFor={fieldId('broker-url')} className={lbl}>Broker URL</label>
              <input id={fieldId('broker-url')} type="text" value={String(config.broker_url ?? '')}
                onChange={e => setConfig({ ...config, broker_url: e.target.value })}
                className={inp} placeholder="tcp://broker.hivemq.com:1883" />
            </div>
            <div className="col-span-2 grid grid-cols-2 gap-4">
              <div>
                <label htmlFor={fieldId('client-id')} className={lbl}>Client ID</label>
                <input id={fieldId('client-id')} type="text" value={String(config.client_id ?? '')}
                  onChange={e => setConfig({ ...config, client_id: e.target.value })}
                  className={inp} />
              </div>
              <div>
                <label htmlFor={fieldId('qos')} className={lbl}>QoS 等級</label>
                <select id={fieldId('qos')} value={config.qos !== undefined ? String(config.qos) : '0'}
                  onChange={e => setConfig({ ...config, qos: parseInt(e.target.value) })}
                  className={inp}>
                  <option value={0}>0 — 最多一次</option>
                  <option value={1}>1 — 至少一次</option>
                  <option value={2}>2 — 恰好一次</option>
                </select>
              </div>
            </div>
            <div className="col-span-2 grid grid-cols-2 gap-4">
              <div>
                <label htmlFor={fieldId('username')} className={lbl}>使用者名稱（選用）</label>
                <input id={fieldId('username')} type="text" value={String(config.username ?? '')}
                  onChange={e => setConfig({ ...config, username: e.target.value })}
                  className={inp} />
              </div>
              <div>
                <label htmlFor={fieldId('password')} className={lbl}>密碼（選用）</label>
                <input id={fieldId('password')} type="password" value={String(config.password ?? '')}
                  onChange={e => setConfig({ ...config, password: e.target.value })}
                  className={inp} />
              </div>
            </div>
            <div className="col-span-2">
              <label htmlFor={fieldId('topics')} className={lbl}>訂閱主題（逗號分隔）</label>
              <input id={fieldId('topics')} type="text"
                value={Array.isArray(config.topics) ? (config.topics as string[]).join(', ') : String(config.topics ?? '')}
                onChange={e => setConfig({ ...config, topics: e.target.value.split(',').map(s => s.trim()).filter(Boolean) })}
                className={inp} placeholder="sensors/#, factory/line1/+" />
            </div>
            <div className="col-span-2">
              <label htmlFor={fieldId('use-tls')} className="flex cursor-pointer items-center gap-2">
                <input id={fieldId('use-tls')} type="checkbox"
                  checked={config.use_tls === true || config.use_tls === 'true'}
                  onChange={e => setConfig({ ...config, use_tls: e.target.checked })} />
                <span className="text-sm text-slate-300">啟用 TLS/SSL</span>
              </label>
            </div>
          </>
        );

      default:
        return (
          <div className="col-span-2">
            <label htmlFor={fieldId('config-json')} className={lbl}>配置 JSON</label>
            <textarea
              id={fieldId('config-json')}
              value={JSON.stringify(config, null, 2)}
              onChange={e => {
                try { setConfig(JSON.parse(e.target.value)); } catch { /* ignore */ }
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
            placeholder="例：產線 PLC-01"
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
                    <h4 className="text-sm font-semibold text-slate-400 mb-3">進階設定</h4>
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
                            <label htmlFor={fieldId('retry-delay')} className="block text-sm font-medium text-slate-300">重試延遲 (ms)</label>
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
          取消
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
