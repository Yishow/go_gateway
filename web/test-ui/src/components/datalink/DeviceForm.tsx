import { useState, useEffect } from 'react';
import type { Device, CreateDeviceRequest, UpdateDeviceRequest, ProtocolType } from '../../types/datalink';
import { protocolAPI } from '../../services/datalink';

interface DeviceFormProps {
  device?: Device;
  onSubmit: (data: CreateDeviceRequest | UpdateDeviceRequest) => Promise<void>;
  onCancel: () => void;
}

export default function DeviceForm({ device, onSubmit, onCancel }: DeviceFormProps) {
  const [name, setName] = useState(device?.name || '');
  const [description, setDescription] = useState(device?.description || '');
  const [protocol, setProtocol] = useState<ProtocolType>(device?.protocol || 'modbus_tcp');
  const [retryCount, setRetryCount] = useState(device?.retry_count || 3);
  const [retryDelay, setRetryDelay] = useState(device?.retry_delay_ms || 1000);
  
  // Dynamic Configuration
  const [config, setConfig] = useState<Record<string, any>>(
    device?.connection_config 
    ? (typeof device.connection_config === 'string' ? JSON.parse(device.connection_config) : device.connection_config)
    : {}
  );

  const [protocols, setProtocols] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchProtocols();
  }, []);

  const fetchProtocols = async () => {
    try {
      const list = await protocolAPI.list();
      setProtocols(list);
    } catch (err) {
      console.error('Failed to load protocols', err);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      if (device) {
        await onSubmit({
          name,
          description,
          retry_count: Number(retryCount),
          retry_delay_ms: Number(retryDelay),
          connection_config: config,
        } as UpdateDeviceRequest);
      } else {
        await onSubmit({
          name,
          description,
          protocol,
          retry_count: Number(retryCount),
          retry_delay_ms: Number(retryDelay),
          connection_config: config,
        } as CreateDeviceRequest);
      }
    } catch (err: any) {
      setError(err.message || 'Failed to save device');
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
              <label className="block text-sm font-medium text-slate-300">Host</label>
              <input
                type="text"
                value={config.host || ''}
                onChange={e => setConfig({...config, host: e.target.value})}
                placeholder="192.168.1.100"
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-300">Port</label>
              <input
                type="number"
                value={config.port || 502}
                onChange={e => setConfig({...config, port: Number(e.target.value)})}
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-300">Slave ID</label>
              <input
                type="number"
                value={config.slave_id || 1}
                onChange={e => setConfig({...config, slave_id: Number(e.target.value)})}
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                required
              />
            </div>
          </>
        );
      case 'modbus_rtu':
          return (
            <>
              <div>
                <label className="block text-sm font-medium text-slate-300">Port</label>
                <input
                  type="text"
                  value={config.port || ''}
                  onChange={e => setConfig({...config, port: e.target.value})}
                  placeholder="COM1 or /dev/ttyUSB0"
                  className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-slate-300">Baud Rate</label>
                <input
                  type="number"
                  value={config.baud_rate || 9600}
                  onChange={e => setConfig({...config, baud_rate: Number(e.target.value)})}
                  className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                  required
                />
              </div>
               <div>
                  <label className="block text-sm font-medium text-slate-300">Data Bits</label>
                  <input
                    type="number"
                    value={config.data_bits || 8}
                    onChange={e => setConfig({...config, data_bits: Number(e.target.value)})}
                    className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                    required
                  />
                </div>
                 <div>
                  <label className="block text-sm font-medium text-slate-300">Stop Bits</label>
                  <input
                    type="number"
                    value={config.stop_bits || 1}
                    onChange={e => setConfig({...config, stop_bits: Number(e.target.value)})}
                    className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                    required
                  />
                </div>
                 <div>
                  <label className="block text-sm font-medium text-slate-300">Parity</label>
                  <select
                     value={config.parity || 'N'}
                     onChange={e => setConfig({...config, parity: e.target.value})}
                     className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                  >
                      <option value="N">None</option>
                      <option value="E">Even</option>
                      <option value="O">Odd</option>
                  </select>
                </div>
              <div>
                <label className="block text-sm font-medium text-slate-300">Slave ID</label>
                <input
                  type="number"
                  value={config.slave_id || 1}
                  onChange={e => setConfig({...config, slave_id: Number(e.target.value)})}
                  className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                  required
                />
              </div>
            </>
          );
      case 'modbus_udp':
          // Similar to TCP
         return (
          <>
            <div>
              <label className="block text-sm font-medium text-slate-300">Host</label>
              <input
                type="text"
                value={config.host || ''}
                onChange={e => setConfig({...config, host: e.target.value})}
                placeholder="192.168.1.100"
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                required
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-slate-300">Port</label>
              <input
                type="number"
                value={config.port || 502}
                onChange={e => setConfig({...config, port: Number(e.target.value)})}
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                required
              />
            </div>
             <div>
                <label className="block text-sm font-medium text-slate-300">Slave ID</label>
                <input
                  type="number"
                  value={config.slave_id || 1}
                  onChange={e => setConfig({...config, slave_id: Number(e.target.value)})}
                  className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                  required
                />
              </div>
          </>
        );
      
      // Implement other protocols fields...
      default:
        // Generic JSON editor fallback
        return (
          <div className="col-span-2">
            <label className="block text-sm font-medium text-slate-300">Config (JSON)</label>
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
          <label className="block text-sm font-medium text-slate-300">Name</label>
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
          <label className="block text-sm font-medium text-slate-300">Description</label>
          <input
            type="text"
            value={description}
            onChange={e => setDescription(e.target.value)}
            className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
            placeholder="Optional device description"
          />
        </div>
      </div>

      <div className="border-t border-slate-700 pt-4">
          <h4 className="text-sm font-semibold text-slate-400 mb-3">Connection Settings</h4>
          <div className="grid grid-cols-2 gap-4">
               {!device && (
                   <div className="col-span-2">
                    <label className="block text-sm font-medium text-slate-300">Protocol</label>
                    <select
                        value={protocol}
                        onChange={e => {
                            setProtocol(e.target.value as ProtocolType);
                            setConfig({}); // Reset config on protocol change
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
                            <label className="block text-sm font-medium text-slate-300">Retry Count</label>
                            <input
                                type="number"
                                value={retryCount}
                                onChange={e => setRetryCount(Number(e.target.value))}
                                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                            />
                        </div>
                         <div>
                            <label className="block text-sm font-medium text-slate-300">Retry Delay (ms)</label>
                            <input
                                type="number"
                                value={retryDelay}
                                onChange={e => setRetryDelay(Number(e.target.value))}
                                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
                            />
                        </div>
                     </div>
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
          {loading ? 'Saving...' : 'Save Device'}
        </button>
      </div>
    </form>
  );
}
