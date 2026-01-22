
import { useState, useMemo } from 'react';
import type { Point, CreatePointRequest, UpdatePointRequest } from '../../types/datalink';
import PointTable from '../../components/datalink/PointTable';
import PointForm from '../../components/datalink/PointForm';
import ConfirmDialog from '../../components/datalink/ConfirmDialog';
import { useToast } from '../../contexts/ToastContext';
import {
  usePointsQuery,
  useCreatePointMutation,
  useUpdatePointMutation,
  useDeletePointMutation,
  usePollPointMutation,
} from '../../hooks/datalink/usePoints';
import { useDevicesQuery } from '../../hooks/datalink/useDevices';
import { usePollingGroupsQuery } from '../../hooks/datalink/usePollingGroups';

export default function PointsPage() {
  // Filters
  const [selectedDeviceId, setSelectedDeviceId] = useState('');
  const [selectedPollingGroupId, setSelectedPollingGroupId] = useState('');
  const [filterEnabled, setFilterEnabled] = useState<string>(''); // 'true', 'false', or ''

  const { showSuccess, showError } = useToast();

  const filters = useMemo(() => ({
      device_id: selectedDeviceId || undefined,
      polling_group_id: selectedPollingGroupId || undefined,
      enabled: filterEnabled === 'true' ? true : filterEnabled === 'false' ? false : undefined,
  }), [selectedDeviceId, selectedPollingGroupId, filterEnabled]);

  // Queries
  const { data: points = [], isLoading, error } = usePointsQuery(filters);
  const { data: devices = [] } = useDevicesQuery();
  const { data: pollingGroups = [] } = usePollingGroupsQuery();

  // Mutations
  const createMutation = useCreatePointMutation();
  const updateMutation = useUpdatePointMutation();
  const deleteMutation = useDeletePointMutation();
  const pollMutation = usePollPointMutation();

  // Modal State
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingPoint, setEditingPoint] = useState<Point | null>(null);

  // Confirm Dialog
  const [confirmDialog, setConfirmDialog] = useState<{
    isOpen: boolean;
    title: string;
    message: string;
    onConfirm: () => void;
    variant?: 'danger' | 'warning' | 'default';
  }>({
    isOpen: false,
    title: '',
    message: '',
    onConfirm: () => {},
  });

  const showConfirm = (title: string, message: string, onConfirm: () => void, variant: 'danger' | 'warning' | 'default' = 'default') => {
    setConfirmDialog({ isOpen: true, title, message, onConfirm, variant });
  };

  const closeConfirm = () => {
    setConfirmDialog(prev => ({ ...prev, isOpen: false }));
  };

  // Handlers
  const handleDelete = (id: string) => {
    showConfirm(
      'Delete Point',
      'Are you sure you want to delete this point? This action cannot be undone.',
      async () => {
        try {
          await deleteMutation.mutateAsync(id);
          showSuccess('Point deleted successfully');
          closeConfirm();
        } catch (err: unknown) {
          const message = err instanceof Error ? err.message : 'Unknown error';
          showError(`Failed to delete point: ${message}`);
          closeConfirm();
        }
      },
      'danger'
    );
  };

  const handlePoll = async (id: string) => {
    try {
        await pollMutation.mutateAsync(id);
        showSuccess('Poll request sent successfully');
    } catch (err: unknown) {
        const message = err instanceof Error ? err.message : 'Unknown error';
        showError(`Failed to poll point: ${message}`);
    }
  };

  const handleEdit = (point: Point) => {
    setEditingPoint(point);
    setIsModalOpen(true);
  };

  const handleCreate = () => {
    setEditingPoint(null);
    setIsModalOpen(true);
  };

  const handleSubmit = async (data: CreatePointRequest | UpdatePointRequest) => {
    try {
      if (editingPoint) {
        await updateMutation.mutateAsync({
          id: editingPoint.id,
          data: data as UpdatePointRequest,
        });
        showSuccess('Point updated successfully');
      } else {
        await createMutation.mutateAsync(data as CreatePointRequest);
        showSuccess('Point created successfully');
      }
      setIsModalOpen(false);
    } catch (err) {
      throw err;
    }
  };

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex justify-between items-center relative z-0">
        <div className="absolute -left-10 -top-10 w-32 h-32 bg-blue-600/20 blur-3xl pointer-events-none"></div>
        <div>
          <h2 className="text-3xl font-bold text-slate-100 tracking-tight">Data Points</h2>
          <p className="text-slate-400 mt-1">Manage data collection points for your devices</p>
        </div>
        <button
          onClick={handleCreate}
          className="px-5 py-2.5 bg-blue-600 hover:bg-blue-500 text-white rounded-xl font-medium transition-all shadow-lg shadow-blue-500/20 hover:shadow-blue-500/40 flex items-center space-x-2"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          <span>Create Point</span>
        </button>
      </div>

       {/* Filters */}
      <div className="flex flex-wrap gap-4 bg-slate-800/80 backdrop-blur-sm p-4 rounded-xl border border-slate-700 shadow-sm">
        <div className="min-w-[200px] flex-1">
             <select
                value={selectedDeviceId}
                onChange={e => setSelectedDeviceId(e.target.value)}
                className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
              >
                <option value="">All Devices</option>
                {devices.map(d => (
                    <option key={d.id} value={d.id}>{d.name}</option>
                ))}
              </select>
        </div>
        <div className="min-w-[200px] flex-1">
             <select
                value={selectedPollingGroupId}
                onChange={e => setSelectedPollingGroupId(e.target.value)}
                className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
              >
                <option value="">All Polling Groups</option>
                {pollingGroups.map(g => (
                    <option key={g.id} value={g.id}>{g.name}</option>
                ))}
              </select>
        </div>
        <div className="min-w-[150px]">
             <select
                value={filterEnabled}
                onChange={e => setFilterEnabled(e.target.value)}
                className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
              >
                <option value="">All Status</option>
                <option value="true">Enabled</option>
                <option value="false">Disabled</option>
              </select>
        </div>
      </div>

      {/* Content */}
      {isLoading ? (
        <div className="space-y-4">
          {[1, 2, 3].map(i => (
            <div key={i} className="h-16 bg-slate-800/50 rounded-lg animate-pulse border border-slate-700/50"></div>
          ))}
        </div>
      ) : error ? (
         <div className="bg-red-500/10 border border-red-500/20 rounded-xl p-6 text-red-400">
           {error instanceof Error ? error.message : 'Failed to fetch points'}
         </div>
      ) : (
        <div className="bg-slate-800 rounded-xl border border-slate-700 overflow-hidden shadow-sm">
            <PointTable
                points={points}
                onEdit={handleEdit}
                onDelete={handleDelete}
                onPoll={handlePoll}
            />
        </div>
      )}

      {/* Confirm Dialog */}
      <ConfirmDialog
        isOpen={confirmDialog.isOpen}
        title={confirmDialog.title}
        message={confirmDialog.message}
        variant={confirmDialog.variant || 'default'}
        onConfirm={confirmDialog.onConfirm}
        onCancel={closeConfirm}
      />

      {/* Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 bg-slate-900/80 backdrop-blur-sm flex items-center justify-center p-4 z-50">
          <div className="bg-slate-800 rounded-2xl border border-slate-700 p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto shadow-2xl ring-1 ring-white/10">
            <h3 className="text-xl font-bold text-slate-100 mb-6 pb-4 border-b border-slate-700 flex items-center justify-between">
              <span>{editingPoint ? 'Edit Point' : 'New Point'}</span>
              <button onClick={() => setIsModalOpen(false)} className="text-slate-400 hover:text-white transition-colors">
                <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </h3>

            <PointForm
              point={editingPoint || undefined}
              onSubmit={handleSubmit}
              onCancel={() => setIsModalOpen(false)}
            />
          </div>
        </div>
      )}
    </div>
  );
}
