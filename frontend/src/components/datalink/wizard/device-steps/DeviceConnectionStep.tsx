import { type ProtocolType } from '../../../../types/datalink';

interface DeviceConnectionStepProps {
  protocol: ProtocolType;
  config: Record<string, string | number | boolean | string[] | undefined>;
  onChange: (config: Record<string, string | number | boolean | string[] | undefined>) => void;
  error?: string;
}

/**
 * 連線設定步驟
 *
 * 依後端 schema_device_models.go 各協議 ConnectionConfig 結構定義欄位：
 * - modbus_tcp / modbus_udp: host, port, slave_id, timeout
 * - modbus_rtu: serial_port, baud_rate, data_bits, stop_bits, parity, slave_id, timeout
 * - fatek_fbs: mode, host/port (tcp) 或 serial_port/baud_rate (serial), station_no, timeout
 * - mc_3e: host, port, network_no, pc_no, io_no, station_no, timeout, data_format
 * - mqtt: broker_url, client_id, username, password, use_tls, topics, qos
 */
export default function DeviceConnectionStep({ protocol, config, onChange, error }: DeviceConnectionStepProps) {
  const fieldId = (name: string) => `wizard-${protocol}-${name}`;

  const handleChange = (key: string, value: string | number | boolean | string[] | undefined) => {
    onChange({ ...config, [key]: value });
  };

  const handleNumber = (key: string, raw: string) => {
    if (raw === '' || /^\d+$/.test(raw)) {
      handleChange(key, raw === '' ? undefined : parseInt(raw, 10));
    }
  };

  const renderFields = () => {
    switch (protocol) {
      // ─────────────────────────────────────────────────────────────────────────
      // Modbus TCP — host, port(502), slave_id, timeout
      // ─────────────────────────────────────────────────────────────────────────
      case 'modbus_tcp':
        return (
          <>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <label htmlFor={fieldId('host')} className="label">主機位址 (Host)</label>
                <input id={fieldId('host')} type="text" value={String(config.host ?? '')}
                  onChange={e => handleChange('host', e.target.value)}
                  className="input" placeholder="192.168.1.10" />
              </div>
              <div>
                <label htmlFor={fieldId('port')} className="label">連接埠 (Port)</label>
                <input id={fieldId('port')} type="text" inputMode="numeric"
                  value={config.port !== undefined ? String(config.port) : ''}
                  onChange={e => handleNumber('port', e.target.value)}
                  className="input" placeholder="502" />
              </div>
            </div>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <label htmlFor={fieldId('slave-id')} className="label">從站 ID (Slave ID)</label>
                <input id={fieldId('slave-id')} type="text" inputMode="numeric"
                  value={config.slave_id !== undefined ? String(config.slave_id) : ''}
                  onChange={e => handleNumber('slave_id', e.target.value)}
                  className="input" placeholder="1" />
              </div>
              <div>
                <label htmlFor={fieldId('timeout')} className="label">逾時秒數 (Timeout)</label>
                <input id={fieldId('timeout')} type="text" inputMode="numeric"
                  value={config.timeout !== undefined ? String(config.timeout) : ''}
                  onChange={e => handleNumber('timeout', e.target.value)}
                  className="input" placeholder="5" />
              </div>
            </div>
          </>
        );

      // ─────────────────────────────────────────────────────────────────────────
      // Modbus UDP — 與 TCP 相同欄位結構
      // ─────────────────────────────────────────────────────────────────────────
      case 'modbus_udp':
        return (
          <>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <label htmlFor={fieldId('host')} className="label">主機位址 (Host)</label>
                <input id={fieldId('host')} type="text" value={String(config.host ?? '')}
                  onChange={e => handleChange('host', e.target.value)}
                  className="input" placeholder="192.168.1.10" />
              </div>
              <div>
                <label htmlFor={fieldId('port')} className="label">連接埠 (Port)</label>
                <input id={fieldId('port')} type="text" inputMode="numeric"
                  value={config.port !== undefined ? String(config.port) : ''}
                  onChange={e => handleNumber('port', e.target.value)}
                  className="input" placeholder="502" />
              </div>
            </div>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <label htmlFor={fieldId('slave-id')} className="label">從站 ID (Slave ID)</label>
                <input id={fieldId('slave-id')} type="text" inputMode="numeric"
                  value={config.slave_id !== undefined ? String(config.slave_id) : ''}
                  onChange={e => handleNumber('slave_id', e.target.value)}
                  className="input" placeholder="1" />
              </div>
              <div>
                <label htmlFor={fieldId('timeout')} className="label">逾時秒數 (Timeout)</label>
                <input id={fieldId('timeout')} type="text" inputMode="numeric"
                  value={config.timeout !== undefined ? String(config.timeout) : ''}
                  onChange={e => handleNumber('timeout', e.target.value)}
                  className="input" placeholder="5" />
              </div>
            </div>
          </>
        );

      // ─────────────────────────────────────────────────────────────────────────
      // Modbus RTU — serial_port, baud_rate, data_bits, stop_bits, parity, slave_id, timeout
      // ─────────────────────────────────────────────────────────────────────────
      case 'modbus_rtu':
        return (
          <>
            <div>
              <label htmlFor={fieldId('serial-port')} className="label">串列埠 (Serial Port)</label>
              <input id={fieldId('serial-port')} type="text" value={String(config.serial_port ?? '')}
                onChange={e => handleChange('serial_port', e.target.value)}
                className="input" placeholder="COM1 或 /dev/ttyUSB0" />
            </div>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <label htmlFor={fieldId('baud-rate')} className="label">波特率 (Baud Rate)</label>
                <select id={fieldId('baud-rate')} value={Number(config.baud_rate ?? 9600)}
                  onChange={e => handleChange('baud_rate', parseInt(e.target.value))}
                  className="input">
                  {[9600, 19200, 38400, 57600, 115200].map(br => (
                    <option key={br} value={br}>{br}</option>
                  ))}
                </select>
              </div>
              <div>
                <label htmlFor={fieldId('data-bits')} className="label">資料位元 (Data Bits)</label>
                <select id={fieldId('data-bits')} value={Number(config.data_bits ?? 8)}
                  onChange={e => handleChange('data_bits', parseInt(e.target.value))}
                  className="input">
                  <option value={7}>7</option>
                  <option value={8}>8</option>
                </select>
              </div>
              <div>
                <label htmlFor={fieldId('stop-bits')} className="label">停止位元 (Stop Bits)</label>
                <select id={fieldId('stop-bits')} value={Number(config.stop_bits ?? 1)}
                  onChange={e => handleChange('stop_bits', parseInt(e.target.value))}
                  className="input">
                  <option value={1}>1</option>
                  <option value={2}>2</option>
                </select>
              </div>
              <div>
                <label htmlFor={fieldId('parity')} className="label">同位檢查 (Parity)</label>
                <select id={fieldId('parity')} value={String(config.parity ?? 'none')}
                  onChange={e => handleChange('parity', e.target.value)}
                  className="input">
                  <option value="none">None（無）</option>
                  <option value="even">Even（偶）</option>
                  <option value="odd">Odd（奇）</option>
                </select>
              </div>
            </div>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <label htmlFor={fieldId('slave-id')} className="label">從站 ID (Slave ID)</label>
                <input id={fieldId('slave-id')} type="text" inputMode="numeric"
                  value={config.slave_id !== undefined ? String(config.slave_id) : ''}
                  onChange={e => handleNumber('slave_id', e.target.value)}
                  className="input" placeholder="1" />
              </div>
              <div>
                <label htmlFor={fieldId('timeout')} className="label">逾時秒數 (Timeout)</label>
                <input id={fieldId('timeout')} type="text" inputMode="numeric"
                  value={config.timeout !== undefined ? String(config.timeout) : ''}
                  onChange={e => handleNumber('timeout', e.target.value)}
                  className="input" placeholder="5" />
              </div>
            </div>
          </>
        );

      // ─────────────────────────────────────────────────────────────────────────
      // FATEK FBs — mode(tcp|serial), host/port 或 serial_port/baud_rate, station_no, timeout
      // ─────────────────────────────────────────────────────────────────────────
      case 'fatek_fbs': {
        const mode = String(config.mode ?? 'tcp');
        return (
          <>
            <div>
              <p className="label">連線模式</p>
              <div className="flex flex-wrap gap-4">
                {(['tcp', 'serial'] as const).map(m => (
                  <label key={m} htmlFor={fieldId(`mode-${m}`)} className="flex cursor-pointer items-center gap-2">
                    <input id={fieldId(`mode-${m}`)} type="radio" name={fieldId('mode')}
                      value={m} checked={mode === m}
                      onChange={() => handleChange('mode', m)} />
                    <span className="text-slate-300">{m === 'tcp' ? 'TCP/IP' : 'Serial (RS232/485)'}</span>
                  </label>
                ))}
              </div>
            </div>
            {mode === 'serial' ? (
              <>
                <div>
                  <label htmlFor={fieldId('serial-port')} className="label">串列埠 (Serial Port)</label>
                  <input id={fieldId('serial-port')} type="text" value={String(config.serial_port ?? '')}
                    onChange={e => handleChange('serial_port', e.target.value)}
                    className="input" placeholder="COM1" />
                </div>
                <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                  <div>
                    <label htmlFor={fieldId('baud-rate')} className="label">波特率 (Baud Rate)</label>
                    <select id={fieldId('baud-rate')} value={Number(config.baud_rate ?? 9600)}
                      onChange={e => handleChange('baud_rate', parseInt(e.target.value))}
                      className="input">
                      {[9600, 19200, 38400, 57600, 115200].map(br => (
                        <option key={br} value={br}>{br}</option>
                      ))}
                    </select>
                  </div>
                </div>
              </>
            ) : (
              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                <div>
                  <label htmlFor={fieldId('host')} className="label">主機位址 (Host)</label>
                  <input id={fieldId('host')} type="text" value={String(config.host ?? '')}
                    onChange={e => handleChange('host', e.target.value)}
                    className="input" placeholder="192.168.1.1" />
                </div>
                <div>
                  <label htmlFor={fieldId('port')} className="label">連接埠 (Port)</label>
                  <input id={fieldId('port')} type="text" inputMode="numeric"
                    value={config.port !== undefined ? String(config.port) : ''}
                    onChange={e => handleNumber('port', e.target.value)}
                    className="input" placeholder="500" />
                </div>
              </div>
            )}
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <label htmlFor={fieldId('station-no')} className="label">站號 (Station No)</label>
                <input id={fieldId('station-no')} type="text" inputMode="numeric"
                  value={config.station_no !== undefined ? String(config.station_no) : ''}
                  onChange={e => handleNumber('station_no', e.target.value)}
                  className="input" placeholder="1" />
              </div>
              <div>
                <label htmlFor={fieldId('timeout')} className="label">逾時秒數 (Timeout)</label>
                <input id={fieldId('timeout')} type="text" inputMode="numeric"
                  value={config.timeout !== undefined ? String(config.timeout) : ''}
                  onChange={e => handleNumber('timeout', e.target.value)}
                  className="input" placeholder="5" />
              </div>
            </div>
          </>
        );
      }

      // ─────────────────────────────────────────────────────────────────────────
      // Mitsubishi MC 3E — host, port, network_no, pc_no, io_no, station_no, timeout, data_format
      // ─────────────────────────────────────────────────────────────────────────
      case 'mc_3e':
        return (
          <>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <label htmlFor={fieldId('host')} className="label">主機位址 (Host)</label>
                <input id={fieldId('host')} type="text" value={String(config.host ?? '')}
                  onChange={e => handleChange('host', e.target.value)}
                  className="input" placeholder="192.168.1.10" />
              </div>
              <div>
                <label htmlFor={fieldId('port')} className="label">連接埠 (Port)</label>
                <input id={fieldId('port')} type="text" inputMode="numeric"
                  value={config.port !== undefined ? String(config.port) : ''}
                  onChange={e => handleNumber('port', e.target.value)}
                  className="input" placeholder="3000" />
              </div>
            </div>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <label htmlFor={fieldId('network-no')} className="label">網路編號 (Network No)</label>
                <input id={fieldId('network-no')} type="text" inputMode="numeric"
                  value={config.network_no !== undefined ? String(config.network_no) : ''}
                  onChange={e => handleNumber('network_no', e.target.value)}
                  className="input" placeholder="0" />
              </div>
              <div>
                <label htmlFor={fieldId('pc-no')} className="label">PC 編號 (PC No)</label>
                <input id={fieldId('pc-no')} type="text" inputMode="numeric"
                  value={config.pc_no !== undefined ? String(config.pc_no) : ''}
                  onChange={e => handleNumber('pc_no', e.target.value)}
                  className="input" placeholder="255" />
              </div>
              <div>
                <label htmlFor={fieldId('io-no')} className="label">I/O 編號 (IO No)</label>
                <input id={fieldId('io-no')} type="text" inputMode="numeric"
                  value={config.io_no !== undefined ? String(config.io_no) : ''}
                  onChange={e => handleNumber('io_no', e.target.value)}
                  className="input" placeholder="1023" />
              </div>
              <div>
                <label htmlFor={fieldId('station-no')} className="label">站號 (Station No)</label>
                <input id={fieldId('station-no')} type="text" inputMode="numeric"
                  value={config.station_no !== undefined ? String(config.station_no) : ''}
                  onChange={e => handleNumber('station_no', e.target.value)}
                  className="input" placeholder="0" />
              </div>
            </div>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <label htmlFor={fieldId('data-format')} className="label">字節序格式 (Data Format)</label>
                <select id={fieldId('data-format')} value={String(config.data_format ?? 'CDAB')}
                  onChange={e => handleChange('data_format', e.target.value)}
                  className="input">
                  <option value="CDAB">CDAB（三菱標準）</option>
                  <option value="ABCD">ABCD</option>
                  <option value="BADC">BADC</option>
                  <option value="DCBA">DCBA</option>
                </select>
              </div>
              <div>
                <label htmlFor={fieldId('timeout')} className="label">逾時秒數 (Timeout)</label>
                <input id={fieldId('timeout')} type="text" inputMode="numeric"
                  value={config.timeout !== undefined ? String(config.timeout) : ''}
                  onChange={e => handleNumber('timeout', e.target.value)}
                  className="input" placeholder="5" />
              </div>
            </div>
          </>
        );

      // ─────────────────────────────────────────────────────────────────────────
      // MQTT — broker_url, client_id, qos, username, password, use_tls, topics
      // ─────────────────────────────────────────────────────────────────────────
      case 'mqtt':
        return (
          <>
            <div>
              <label htmlFor={fieldId('broker-url')} className="label">Broker URL</label>
              <input id={fieldId('broker-url')} type="text" value={String(config.broker_url ?? '')}
                onChange={e => handleChange('broker_url', e.target.value)}
                className="input" placeholder="tcp://broker.hivemq.com:1883" />
            </div>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <label htmlFor={fieldId('client-id')} className="label">Client ID</label>
                <input id={fieldId('client-id')} type="text" value={String(config.client_id ?? '')}
                  onChange={e => handleChange('client_id', e.target.value)}
                  className="input" />
              </div>
              <div>
                <label htmlFor={fieldId('qos')} className="label">QoS 等級</label>
                <select id={fieldId('qos')} value={config.qos !== undefined ? String(config.qos) : '0'}
                  onChange={e => handleChange('qos', parseInt(e.target.value))}
                  className="input">
                  <option value={0}>0 — 最多一次</option>
                  <option value={1}>1 — 至少一次</option>
                  <option value={2}>2 — 恰好一次</option>
                </select>
              </div>
            </div>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <label htmlFor={fieldId('username')} className="label">使用者名稱（選用）</label>
                <input id={fieldId('username')} type="text" value={String(config.username ?? '')}
                  onChange={e => handleChange('username', e.target.value)}
                  className="input" />
              </div>
              <div>
                <label htmlFor={fieldId('password')} className="label">密碼（選用）</label>
                <input id={fieldId('password')} type="password" value={String(config.password ?? '')}
                  onChange={e => handleChange('password', e.target.value)}
                  className="input" />
              </div>
            </div>
            <div>
              <label htmlFor={fieldId('topics')} className="label">訂閱主題（逗號分隔）</label>
              <input id={fieldId('topics')} type="text"
                value={Array.isArray(config.topics) ? config.topics.join(', ') : String(config.topics ?? '')}
                onChange={e => handleChange('topics', e.target.value.split(',').map((s: string) => s.trim()).filter(Boolean))}
                className="input" placeholder="sensors/#, factory/line1/+" />
            </div>
            <label htmlFor={fieldId('use-tls')} className="flex cursor-pointer items-center gap-2">
              <input id={fieldId('use-tls')} type="checkbox"
                checked={config.use_tls === true || config.use_tls === 'true'}
                onChange={e => handleChange('use_tls', e.target.checked)} />
              <span className="text-sm text-slate-300">啟用 TLS/SSL</span>
            </label>
          </>
        );

      default:
        return <p className="italic text-slate-500">此協議尚未支援設定介面。</p>;
    }
  };

  return (
    <div className="animate-in fade-in slide-in-from-right-4 space-y-6 duration-300">
      <div>
        <h3 className="mb-4 text-lg font-medium text-slate-200">連線設定</h3>
        <style>{`
          .label { display:block; font-size:.875rem; font-weight:500; color:#cbd5e1; margin-bottom:.25rem; }
          .input { width:100%; background-color:#0f172a; border:1px solid #334155; border-radius:.5rem; padding:.5rem .75rem; color:#e2e8f0; transition:border-color .2s ease; }
          .input:focus { outline:none; border-color:#3b82f6; }
          @media (prefers-reduced-motion:reduce) { .input { transition:none; } }
        `}</style>
        <div className="space-y-4">{renderFields()}</div>
      </div>
      {error && (
        <div className="rounded-lg border border-red-500/20 bg-red-500/10 p-3 text-sm text-red-400">
          {error}
        </div>
      )}
    </div>
  );
}
