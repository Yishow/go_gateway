import { useState, useEffect } from 'react';
import { deviceAPI } from '../../services/datalink';
import type { Device, CreateDeviceRequest, UpdateDeviceRequest } from '../../types/datalink';
import DeviceCard from '../../components/datalink/DeviceCard';
import DeviceForm from '../../components/datalink/DeviceForm';

export default function DevicesPage() {
  const [devices, setDevices] = useState<Device[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // Modal states
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingDevice, setEditingDevice] = useState<Device | null>(null);

  useEffect(() => {
    fetchDevices();
  }, []);

  const fetchDevices = async () => {
    setLoading(true);
    try {
      const data = await deviceAPI.list();
      setDevices(data);
      setError(null);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch devices');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Are you sure you want to delete this device?')) return;
    
    try {
      await deviceAPI.delete(id);
      setDevices(prev => prev.filter(d => d.id !== id));
    } catch (err: any) {
      alert(`Failed to delete device: ${err.message}`);
    }
  };

  const handleToggleStatus = async (id: string, currentStatus: string) => {
    try {
      if (currentStatus === 'active') {
        const updated = await deviceAPI.disable(id);
        setDevices(prev => prev.map(d => d.id === id ? updated : d));
      } else {
        const updated = await deviceAPI.activate(id);
        setDevices(prev => prev.map(d => d.id === id ? updated : d));
      }
    } catch (err: any) {
      alert(`Failed to update status: ${err.message}`);
    }
  };

  const handleTestConnection = async (id: string) => {
    try {
      const result = await deviceAPI.testConnection(id);
      if (result.success) {
        alert(`Connection Successful! Latency: ${result.latency_ms}ms`);
      } else {
        alert(`Connection Failed: ${result.error}`);
      }
    } catch (err: any) {
      alert(`Test failed: ${err.message}`);
    }
  };

  const handleEdit = (device: Device) => {
    setEditingDevice(device);
    setIsModalOpen(true);
  };

  const handleCreate = () => {
    setEditingDevice(null);
    setIsModalOpen(true);
  };

  const handleSubmit = async (data: CreateDeviceRequest | UpdateDeviceRequest) => {
    try {
      if (editingDevice) {
        const updated = await deviceAPI.update(editingDevice.id, data as UpdateDeviceRequest);
        setDevices(prev => prev.map(d => d.id === editingDevice.id ? updated : d));
      } else {
        const created = await deviceAPI.create(data as CreateDeviceRequest);
        setDevices(prev => [...prev, created]);
      }
      setIsModalOpen(false);
    } catch (err: any) {
      throw err; // 讓 Form 處理錯誤顯示
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold text-slate-100">Devices</h2>
          <p className="text-slate-400">Manage industrial devices and connections</p>
        </div>
        <button
          onClick={handleCreate}
          className="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-lg font-medium transition-colors flex items-center space-x-2"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          <span>Add Device</span>
        </button>
      </div>

      {loading ? (
        <div className="text-center py-12">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto"></div>
          <p className="mt-4 text-slate-400">Loading devices...</p>
        </div>
      ) : error ? (
        <div className="bg-red-500/10 border border-red-500/20 rounded-lg p-4 text-red-400">
          Error: {error}
        </div>
      ) : devices.length === 0 ? (
        <div className="text-center py-12 bg-slate-800/50 rounded-xl border border-slate-700 border-dashed">
          <p className="text-slate-400 mb-4">No devices found</p>
          <button
            onClick={handleCreate}
            className="text-blue-400 hover:text-blue-300 font-medium"
          >
            Create your first device
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {devices.map(device => (
            <DeviceCard
              key={device.id}
              device={device}
              onEdit={handleEdit}
              onDelete={handleDelete}
              onTest={handleTestConnection}
              onToggleStatus={handleToggleStatus}
            />
          ))}
        </div>
      )}

      {/* Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50 backdrop-blur-sm">
          <div className="bg-slate-800 rounded-xl border border-slate-700 p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto shadow-2xl">
             <h3 className="text-xl font-bold text-slate-100 mb-6 pb-4 border-b border-slate-700">
                {editingDevice ? 'Edit Device' : 'New Device'}
             </h3>
             
             <DeviceForm
                device={editingDevice || undefined}
                onSubmit={handleSubmit}
                onCancel={() => setIsModalOpen(false)}
             />
          </div>
        </div>
      )}
    </div>
  );
}
