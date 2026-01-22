/**
 * DevicesPage - 設備管理頁面
 *
 * 使用 TanStack Query 進行 Server State 管理。
 * 移除原有的 useState/useEffect 樣板代碼，改用 Query Hooks。
 */
import { useState } from 'react';
import type { Device, CreateDeviceRequest, UpdateDeviceRequest } from '../../types/datalink';
import DeviceCard from '../../components/datalink/DeviceCard';
import DeviceForm from '../../components/datalink/DeviceForm';
import {
  useDevicesQuery,
  useCreateDeviceMutation,
  useUpdateDeviceMutation,
  useDeleteDeviceMutation,
  useTestConnectionMutation,
  useToggleDeviceStatusMutation,
} from '../../hooks/datalink';

export default function DevicesPage() {
  // 使用 Query Hook 取代 useState + useEffect
  const { data: devices = [], isLoading, error } = useDevicesQuery();

  // Mutation Hooks
  const createMutation = useCreateDeviceMutation();
  const updateMutation = useUpdateDeviceMutation();
  const deleteMutation = useDeleteDeviceMutation();
  const testConnectionMutation = useTestConnectionMutation();
  const toggleStatusMutation = useToggleDeviceStatusMutation();

  // Modal 狀態（本地 UI 狀態，不需 Query 管理）
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingDevice, setEditingDevice] = useState<Device | null>(null);

  /**
   * 刪除設備
   */
  const handleDelete = async (id: string) => {
    if (!window.confirm('Are you sure you want to delete this device?')) return;

    try {
      await deleteMutation.mutateAsync(id);
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Unknown error';
      alert(`Failed to delete device: ${message}`);
    }
  };

  /**
   * 切換設備狀態（使用樂觀更新）
   */
  const handleToggleStatus = (id: string, currentStatus: string) => {
    toggleStatusMutation.mutate(
      { id, currentStatus },
      {
        onError: (err) => {
          alert(`Failed to update status: ${err.message}`);
        },
      }
    );
  };

  /**
   * 測試連線
   */
  const handleTestConnection = async (id: string) => {
    try {
      const result = await testConnectionMutation.mutateAsync(id);
      if (result.success) {
        alert(`Connection Successful! Latency: ${result.latency_ms}ms`);
      } else {
        alert(`Connection Failed: ${result.error}`);
      }
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : 'Unknown error';
      alert(`Test failed: ${message}`);
    }
  };

  /**
   * 開啟編輯 Modal
   */
  const handleEdit = (device: Device) => {
    setEditingDevice(device);
    setIsModalOpen(true);
  };

  /**
   * 開啟新增 Modal
   */
  const handleCreate = () => {
    setEditingDevice(null);
    setIsModalOpen(true);
  };

  /**
   * 提交表單（新增或更新）
   */
  const handleSubmit = async (data: CreateDeviceRequest | UpdateDeviceRequest) => {
    try {
      if (editingDevice) {
        await updateMutation.mutateAsync({
          id: editingDevice.id,
          data: data as UpdateDeviceRequest,
        });
      } else {
        await createMutation.mutateAsync(data as CreateDeviceRequest);
      }
      setIsModalOpen(false);
    } catch (err) {
      throw err; // Form will handle error display
    }
  };

  return (
    <div className="space-y-8">
      {/* Header with Glow */}
      <div className="flex justify-between items-center relative z-0">
        <div className="absolute -left-10 -top-10 w-32 h-32 bg-blue-600/20 blur-3xl pointer-events-none"></div>
        <div>
          <h2 className="text-3xl font-bold text-slate-100 tracking-tight">Devices</h2>
          <p className="text-slate-400 mt-1">Manage industrial devices and communication channels</p>
        </div>
        <button
          onClick={handleCreate}
          className="px-5 py-2.5 bg-blue-600 hover:bg-blue-500 text-white rounded-xl font-medium transition-all shadow-lg shadow-blue-500/20 hover:shadow-blue-500/40 flex items-center space-x-2"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          <span>Add Device</span>
        </button>
      </div>

      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {/* Skeleton Loaders */}
          {[1, 2, 3].map(i => (
            <div key={i} className="bg-slate-800/50 rounded-xl p-5 border border-slate-700/50 h-48 animate-pulse">
              <div className="flex space-x-3 mb-4">
                <div className="w-3 h-3 rounded-full bg-slate-700"></div>
                <div className="h-6 bg-slate-700 rounded w-1/2"></div>
              </div>
              <div className="h-4 bg-slate-700/50 rounded w-full mb-2"></div>
              <div className="h-4 bg-slate-700/50 rounded w-3/4"></div>
            </div>
          ))}
        </div>
      ) : error ? (
        <div className="bg-red-500/10 border border-red-500/20 rounded-xl p-6 text-red-400 flex items-center space-x-3">
          <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span>{error instanceof Error ? error.message : 'Failed to fetch devices'}</span>
        </div>
      ) : devices.length === 0 ? (
        <div
          className="text-center py-20 bg-slate-800/30 rounded-2xl border-2 border-slate-700/50 border-dashed group hover:border-blue-500/30 transition-colors cursor-pointer"
          onClick={handleCreate}
        >
          <div className="w-16 h-16 bg-slate-700/50 rounded-full mx-auto flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
            <svg className="w-8 h-8 text-slate-500 group-hover:text-blue-400 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
          </div>
          <h3 className="text-xl font-bold text-slate-200 mb-2">No devices configured</h3>
          <p className="text-slate-500 max-w-sm mx-auto mb-6">
            Get started by adding your first industrial device protocol connection.
          </p>
          <button className="text-blue-400 hover:text-blue-300 font-medium hover:underline">
            Create your first device &rarr;
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
              onTestConnection={handleTestConnection}
              onToggleStatus={handleToggleStatus}
            />
          ))}
        </div>
      )}

      {/* Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 bg-slate-900/80 backdrop-blur-sm flex items-center justify-center p-4 z-50">
          <div className="bg-slate-800 rounded-2xl border border-slate-700 p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto shadow-2xl ring-1 ring-white/10">
            <h3 className="text-xl font-bold text-slate-100 mb-6 pb-4 border-b border-slate-700 flex items-center justify-between">
              <span>{editingDevice ? 'Edit Device' : 'New Device'}</span>
              <button onClick={() => setIsModalOpen(false)} className="text-slate-400 hover:text-white transition-colors">
                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
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
