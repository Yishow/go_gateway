import { type ProtocolType } from '../../../../types/datalink';

interface DeviceConnectionStepProps {
  protocol: ProtocolType;
  config: Record<string, any>;
  onChange: (config: Record<string, any>) => void;
  error?: string;
}

export default function DeviceConnectionStep({ protocol, config, onChange, error }: DeviceConnectionStepProps) {

  const handleChange = (key: string, value: any) => {
      onChange({ ...config, [key]: value });
  };

  const renderFields = () => {
    switch (protocol) {
      case 'modbus_tcp':
        return (
          <>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="label">Host IP Address</label>
                <input
                  type="text"
                  value={config.host || ''}
                  onChange={e => handleChange('host', e.target.value)}
                  className="input"
                  placeholder="192.168.1.10"
                />
              </div>
              <div>
                <label className="label">Port</label>
                <input
                  type="number"
                  value={config.port || 502}
                  onChange={e => handleChange('port', parseInt(e.target.value))}
                  className="input"
                />
              </div>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="label">Slave ID</label>
                <input
                  type="number"
                  value={config.slave_id || 1}
                  onChange={e => handleChange('slave_id', parseInt(e.target.value))}
                  className="input"
                />
              </div>
              <div>
                <label className="label">Timeout (ms)</label>
                <input
                  type="number"
                  value={config.timeout || 1000}
                  onChange={e => handleChange('timeout', parseInt(e.target.value))}
                  className="input"
                />
              </div>
            </div>
          </>
        );
      case 'modbus_rtu':
        return (
            <>
                <div>
                    <label className="label">Serial Port</label>
                    <input
                        type="text"
                        value={config.serial_port || ''}
                        onChange={e => handleChange('serial_port', e.target.value)}
                        className="input"
                        placeholder={window.navigator.userAgent.includes('Win') ? 'COM1' : '/dev/ttyUSB0'}
                    />
                </div>
                <div className="grid grid-cols-2 gap-4">
                    <div>
                        <label className="label">Baud Rate</label>
                        <select
                            value={config.baud_rate || 9600}
                            onChange={e => handleChange('baud_rate', parseInt(e.target.value))}
                            className="input"
                        >
                            {[9600, 19200, 38400, 57600, 115200].map(br => (
                                <option key={br} value={br}>{br}</option>
                            ))}
                        </select>
                    </div>
                    <div>
                        <label className="label">Data Bits</label>
                        <select
                            value={config.data_bits || 8}
                            onChange={e => handleChange('data_bits', parseInt(e.target.value))}
                            className="input"
                        >
                            <option value={8}>8</option>
                            <option value={7}>7</option>
                        </select>
                    </div>
                </div>
                <div className="grid grid-cols-2 gap-4">
                     <div>
                        <label className="label">Parity</label>
                        <select
                            value={config.parity || 'none'}
                            onChange={e => handleChange('parity', e.target.value)}
                            className="input"
                        >
                            <option value="none">None</option>
                            <option value="even">Even</option>
                            <option value="odd">Odd</option>
                        </select>
                    </div>
                    <div>
                        <label className="label">Stop Bits</label>
                        <select
                            value={config.stop_bits || 1}
                            onChange={e => handleChange('stop_bits', parseInt(e.target.value))}
                            className="input"
                        >
                            <option value={1}>1</option>
                            <option value={2}>2</option>
                        </select>
                    </div>
                </div>
                 <div className="grid grid-cols-2 gap-4">
                    <div>
                        <label className="label">Slave ID</label>
                        <input
                            type="number"
                            value={config.slave_id || 1}
                            onChange={e => handleChange('slave_id', parseInt(e.target.value))}
                            className="input"
                        />
                    </div>
                    <div>
                        <label className="label">Timeout (ms)</label>
                        <input
                            type="number"
                            value={config.timeout || 1000}
                            onChange={e => handleChange('timeout', parseInt(e.target.value))}
                            className="input"
                        />
                    </div>
                </div>
            </>
        );

       case 'fatek_fbs':
        return (
          <>
            <div className="mb-4">
               <label className="label">Connection Mode</label>
               <div className="flex space-x-4">
                  <label className="flex items-center space-x-2 cursor-pointer">
                      <input 
                        type="radio" 
                        name="fatek_mode" 
                        value="tcp" 
                        checked={config.mode === 'tcp'} 
                        onChange={() => handleChange('mode', 'tcp')}
                      />
                      <span className="text-slate-300">TCP/IP</span>
                  </label>
                  <label className="flex items-center space-x-2 cursor-pointer">
                      <input 
                        type="radio" 
                        name="fatek_mode" 
                        value="serial" 
                        checked={config.mode === 'serial'}
                        onChange={() => handleChange('mode', 'serial')}
                      />
                      <span className="text-slate-300">Serial (RS232/485)</span>
                  </label>
               </div>
            </div>

            {config.mode === 'tcp' ? (
                 <div className="grid grid-cols-2 gap-4 fade-in">
                    <div>
                        <label className="label">Host IP</label>
                        <input
                            type="text"
                            value={config.host || ''}
                            onChange={e => handleChange('host', e.target.value)}
                            className="input"
                            placeholder="192.168.1.1"
                        />
                    </div>
                    <div>
                        <label className="label">Port</label>
                        <input
                            type="number"
                            value={config.port || 500}
                            onChange={e => handleChange('port', parseInt(e.target.value))}
                            className="input"
                        />
                    </div>
                </div>
            ) : (
                 <div className="grid grid-cols-2 gap-4 fade-in">
                     <div>
                        <label className="label">Serial Port</label>
                        <input
                            type="text"
                            value={config.serial_port || ''}
                            onChange={e => handleChange('serial_port', e.target.value)}
                            className="input"
                            placeholder="COM1"
                        />
                    </div>
                     <div>
                        <label className="label">Baud Rate</label>
                        <select
                            value={config.baud_rate || 9600}
                            onChange={e => handleChange('baud_rate', parseInt(e.target.value))}
                            className="input"
                        >
                            {[9600, 19200, 38400, 57600, 115200].map(br => (
                                <option key={br} value={br}>{br}</option>
                            ))}
                        </select>
                    </div>
                </div>
            )}
            
            <div className="grid grid-cols-2 gap-4 mt-4">
                <div>
                    <label className="label">Station Number</label>
                    <input
                        type="number"
                        value={config.station_no || 1}
                        onChange={e => handleChange('station_no', parseInt(e.target.value))}
                        className="input"
                    />
                </div>
                 <div>
                    <label className="label">Timeout (ms)</label>
                    <input
                        type="number"
                        value={config.timeout || 1000}
                        onChange={e => handleChange('timeout', parseInt(e.target.value))}
                        className="input"
                    />
                </div>
            </div>
          </>
        );

      case 'mqtt':
          return (
              <>
                 <div>
                    <label className="label">Broker URL</label>
                    <input
                        type="text"
                        value={config.broker_url || ''}
                        onChange={e => handleChange('broker_url', e.target.value)}
                        className="input"
                        placeholder="tcp://broker.hivemq.com:1883"
                    />
                </div>
                 <div className="grid grid-cols-2 gap-4">
                    <div>
                        <label className="label">Client ID</label>
                        <input
                            type="text"
                            value={config.client_id || ''}
                            onChange={e => handleChange('client_id', e.target.value)}
                            className="input"
                        />
                    </div>
                     <div className="flex items-center pt-6">
                        <label className="flex items-center space-x-2 cursor-pointer">
                            <input 
                                type="checkbox"
                                checked={config.use_tls || false}
                                onChange={e => handleChange('use_tls', e.target.checked)}
                            />
                            <span className="text-slate-300">Enable TLS/SSL</span>
                        </label>
                    </div>
                </div>
                 <div className="grid grid-cols-2 gap-4">
                    <div>
                        <label className="label">Username (Optional)</label>
                        <input
                            type="text"
                            value={config.username || ''}
                            onChange={e => handleChange('username', e.target.value)}
                            className="input"
                        />
                    </div>
                    <div>
                        <label className="label">Password (Optional)</label>
                        <input
                            type="password"
                            value={config.password || ''}
                            onChange={e => handleChange('password', e.target.value)}
                            className="input"
                        />
                    </div>
                </div>
                 <div>
                    <label className="label">Topics (Comma separated)</label>
                    <input
                        type="text"
                        value={Array.isArray(config.topics) ? config.topics.join(', ') : (config.topics || '')}
                        onChange={e => handleChange('topics', e.target.value.split(',').map((s: string) => s.trim()))}
                        className="input"
                        placeholder="sensors/#, factory/line1/+"
                    />
                </div>
              </>
          );

      default:
        return <div className="text-slate-500 italic">Configuration not implemented for this protocol yet.</div>;
    }
  };

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-right-4 duration-300">
      <div>
        <h3 className="text-lg font-medium text-slate-200 mb-4">Connection Settings</h3>
        <style>{`
            .label { display: block; font-size: 0.875rem; font-weight: 500; color: #cbd5e1; margin-bottom: 0.25rem; }
            .input { width: 100%; background-color: #0f172a; border: 1px solid #334155; border-radius: 0.5rem; padding: 0.5rem 0.75rem; color: #e2e8f0; transition: all; }
            .input:focus { outline: none; border-color: #3b82f6; }
            .fade-in { animation: fadeIn 0.3s ease-in; }
            @keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
        `}</style>
        <div className="space-y-4">
             {renderFields()}
        </div>
      </div>
      {error && (
        <div className="bg-red-500/10 border border-red-500/20 rounded-lg p-3 text-red-400 text-sm">
            {error}
        </div>
      )}
    </div>
  );
}
