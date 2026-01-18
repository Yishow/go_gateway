import type { Device } from '../../types/datalink';

interface DeviceCardProps {
  device: Device;
  onEdit: (device: Device) => void;
  onDelete: (id: string) => void;
  onTest: (id: string) => void;
  onToggleStatus: (id: string, currentStatus: string) => void;
}

export default function DeviceCard({
  device,
  onEdit,
  onDelete,
  onTest,
  onToggleStatus,
}: DeviceCardProps) {
  const isOnline = device.status === 'active';
  const hasError = device.status === 'error';

  return (
    <div className="bg-slate-800 rounded-lg border border-slate-700 overflow-hidden hover:border-blue-500/50 transition-colors">
      {/* Header */}
      <div className="px-6 py-4 border-b border-slate-700 flex justify-between items-start">
        <div>
          <div className="flex items-center space-x-2">
            <h3 className="text-lg font-semibold text-slate-100">{device.name}</h3>
            <span
              className={`px-2 py-0.5 rounded text-xs font-medium uppercase ${
                isOnline
                  ? 'bg-emerald-500/10 text-emerald-400'
                  : hasError
                  ? 'bg-red-500/10 text-red-400'
                  : 'bg-slate-600/10 text-slate-400'
              }`}
            >
              {device.status}
            </span>
          </div>
          <p className="text-sm text-slate-400 mt-1">{device.description || 'No description'}</p>
        </div>
        <div className="flex space-x-2">
          <button
            onClick={() => onEdit(device)}
            className="p-1.5 text-slate-400 hover:text-blue-400 hover:bg-slate-700 rounded transition-colors"
            title="Edit"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
            </svg>
          </button>
          <button
            onClick={() => onDelete(device.id)}
            className="p-1.5 text-slate-400 hover:text-red-400 hover:bg-slate-700 rounded transition-colors"
            title="Delete"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          </button>
        </div>
      </div>

      {/* Body */}
      <div className="px-6 py-4 space-y-3">
        <div className="flex justify-between text-sm">
          <span className="text-slate-500">Protocol</span>
          <span className="text-slate-300 font-mono">{device.protocol}</span>
        </div>
        <div className="flex justify-between text-sm">
          <span className="text-slate-500">Retry Count</span>
          <span className="text-slate-300 font-mono">{device.retry_count}</span>
        </div>
        
        {device.last_error && (
            <div className="mt-2 p-2 bg-red-900/20 border border-red-900/50 rounded text-xs text-red-300 break-words">
                {device.last_error}
            </div>
        )}
      </div>

      {/* Footer */}
      <div className="px-6 py-3 bg-slate-900/50 border-t border-slate-700 flex justify-between items-center">
        <button
          onClick={() => onTest(device.id)}
          className="text-xs font-medium text-blue-400 hover:text-blue-300 transition-colors"
        >
          Test Connection
        </button>
        <button
            onClick={() => onToggleStatus(device.id, device.status)}
            className={`text-xs font-medium transition-colors ${
                isOnline 
                ? 'text-amber-400 hover:text-amber-300' 
                : 'text-emerald-400 hover:text-emerald-300'
            }`}
        >
            {isOnline ? 'Disable' : 'Activate'}
        </button>
      </div>
    </div>
  );
}
