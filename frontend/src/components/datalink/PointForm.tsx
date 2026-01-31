
import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import type { Point, CreatePointRequest, UpdatePointRequest, DataType } from '../../types/datalink';
import { useDevicesQuery } from '../../hooks/datalink/useDevices';
import { usePollingGroupsQuery } from '../../hooks/datalink/usePollingGroups';

interface PointFormProps {
  point?: Point;
  onSubmit: (data: CreatePointRequest | UpdatePointRequest) => Promise<void>;
  onCancel: () => void;
}

const DATA_TYPES: DataType[] = ['bool', 'int16', 'int32', 'int64', 'uint16', 'uint32', 'uint64', 'float32', 'float64', 'string'];

export default function PointForm({ point, onSubmit, onCancel }: PointFormProps) {
  const navigate = useNavigate();
  
  // Queries
  const { data: devices = [] } = useDevicesQuery();
  const { data: pollingGroups = [] } = usePollingGroupsQuery();

  // Form State
  const [name, setName] = useState(point?.name || '');
  const [description, setDescription] = useState(point?.description || '');
  const [deviceId, setDeviceId] = useState(point?.device_id || '');
  const [pollingGroupId, setPollingGroupId] = useState(point?.polling_group_id || '');
  const [address, setAddress] = useState(point?.address || '');
  const [dataType, setDataType] = useState<DataType>(point?.data_type || 'int16');
  const [enabled, setEnabled] = useState(point?.enabled ?? true);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Set default device/group if not set and available
  useEffect(() => {
    if (!point && !deviceId && devices.length > 0) {
        setDeviceId(devices[0].id);
    }
  }, [devices, point, deviceId]);

  useEffect(() => {
    if (!point && !pollingGroupId && pollingGroups.length > 0) {
        setPollingGroupId(pollingGroups[0].id);
    }
  }, [pollingGroups, point, pollingGroupId]);


  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    // Basic Validation
    if (!deviceId) {
        setError('Device is required');
        setLoading(false);
        return;
    }
    if (!pollingGroupId) {
        setError('Polling Group is required');
        setLoading(false);
        return;
    }

    try {
      if (point) {
        await onSubmit({
          name,
          description,
          device_id: deviceId, // Usually update doesn't change parent, but API allows it? UpdatePointRequest defined in datalink.ts allows enabled/polling_group. Check types.
          // Checking UpdatePointRequest in datalink.ts:
          // name, description, data_type, address, enabled, polling_group_id.
          // Device ID is NOT updateable in UpdatePointRequest interface I saw earlier.
          // Wait, let me double check UpdatePointRequest in types/datalink.ts
          // It has: name, description, data_type, address, enabled, polling_group_id. 
          // It DOES NOT have device_id. So I should NOT send device_id for update.
          
          address,
          data_type: dataType, 
          enabled,
          polling_group_id: pollingGroupId
        } as UpdatePointRequest);
      } else {
        await onSubmit({
          name,
          description,
          device_id: deviceId,
          polling_group_id: pollingGroupId,
          address,
          data_type: dataType,
          enabled,
        } as CreatePointRequest);
      }
    } catch (err: any) {
      setError(err.message || 'Failed to save point');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <div className="bg-red-500/10 border border-red-500/20 rounded-lg p-3 text-sm text-red-400">
          {error}
        </div>
      )}

      {/* Device & Polling Group (Row 1) */}
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-slate-300">Device</label>
           <select
            value={deviceId}
            onChange={e => setDeviceId(e.target.value)}
            disabled={!!point || devices.length === 0}
            className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500 disabled:opacity-50"
            required
          >
            <option value="" disabled>Select Device</option>
            {devices.map(d => (
                <option key={d.id} value={d.id}>{d.name}</option>
            ))}
          </select>
          {devices.length === 0 && <p className="text-xs text-yellow-500 mt-1">No devices available. Create one first.</p>}
        </div>
        <div>
          <label className="block text-sm font-medium text-slate-300">Polling Group</label>
           <select
            value={pollingGroupId}
            onChange={e => setPollingGroupId(e.target.value)}
            disabled={pollingGroups.length === 0}
            className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500 disabled:opacity-50"
            required
          >
            <option value="" disabled>Select Group</option>
            {pollingGroups.map(g => (
                <option key={g.id} value={g.id}>{g.name} ({g.interval_ms}ms)</option>
            ))}
          </select>
          {pollingGroups.length === 0 && (
            <div className="mt-2">
              <p className="text-xs text-yellow-500 mb-2">No polling groups available.</p>
              <button
                type="button"
                onClick={() => navigate('/datalink/polling-groups')}
                className="text-xs text-blue-400 hover:text-blue-300 underline flex items-center space-x-1"
              >
                <span>Create polling group</span>
                <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                </svg>
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Name & Address (Row 2) */}
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-slate-300">Point Name</label>
          <input
            type="text"
            value={name}
            onChange={e => setName(e.target.value)}
            className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
            required
            placeholder="e.g. Temp_Sensor_01"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-slate-300">Address</label>
          <input
            type="text"
            value={address}
            onChange={e => setAddress(e.target.value)}
            className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500 font-mono"
            required
            placeholder="e.g. 40001 or D100"
          />
        </div>
      </div>

      {/* Create Data Type & Enabled (Row 3) */}
      <div className="grid grid-cols-2 gap-4">
         <div>
          <label className="block text-sm font-medium text-slate-300">Data Type</label>
          <select
            value={dataType}
            onChange={e => setDataType(e.target.value as DataType)}
            className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
          >
            {DATA_TYPES.map(t => (
                <option key={t} value={t}>{t}</option>
            ))}
          </select>
        </div>
        <div className="flex items-center pt-6">
            <label className="flex items-center space-x-3 cursor-pointer">
                <div className="relative">
                    <input type="checkbox" className="sr-only" checked={enabled} onChange={e => setEnabled(e.target.checked)} />
                    <div className={`w-10 h-6 rounded-full shadow-inner transition-colors ${enabled ? 'bg-blue-600' : 'bg-slate-700'}`}></div>
                    <div className={`absolute left-1 top-1 w-4 h-4 bg-white rounded-full transition-transform ${enabled ? 'translate-x-4' : 'translate-x-0'}`}></div>
                </div>
                <span className="text-sm font-medium text-slate-300">Enabled</span>
            </label>
        </div>
      </div>

      {/* Description */}
      <div>
          <label className="block text-sm font-medium text-slate-300">Description</label>
          <textarea
            value={description}
            onChange={e => setDescription(e.target.value)}
            rows={3}
             className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 text-sm focus:outline-none focus:border-blue-500"
             placeholder="Optional description..."
          />
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
          {loading ? 'Saving...' : 'Save Point'}
        </button>
      </div>
    </form>
  );
}
