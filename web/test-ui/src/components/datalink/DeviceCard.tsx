import type { Device } from '../../types/datalink';

interface DeviceCardProps {
  device: Device;
  onEdit: (device: Device) => void;
  onDelete: (id: string) => void;
  onTestConnection: (id: string, protocol: string) => void;
  onToggleStatus: (id: string, currentStatus: string) => void;
}

export default function DeviceCard({
  device,
  onEdit,
  onDelete,
  onTestConnection,
  onToggleStatus,
}: DeviceCardProps) {
  const isOnline = device.status === 'active';
  const isError = device.status === 'error';
  
  // Status Color Logic
  const statusColor = isOnline 
    ? 'bg-emerald-500 text-emerald-400 border-emerald-500/20' 
    : isError 
      ? 'bg-red-500 text-red-400 border-red-500/20' 
      : 'bg-slate-600 text-slate-400 border-slate-600/20';
      
  const statusDot = isOnline
    ? 'bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.6)]'
    : isError
      ? 'bg-red-400'
      : 'bg-slate-400';

  return (
    <div className={`group bg-slate-800 rounded-xl border p-5 transition-all duration-300 hover:shadow-lg hover:shadow-blue-500/5 hover:-translate-y-1 relative overflow-hidden ${statusColor} hover:border-blue-500/50`}>
      {/* Status Line Top */}
      <div className={`absolute top-0 left-0 w-full h-1 ${isOnline ? 'bg-emerald-500' : 'bg-transparent'}`}></div>

      <div className="flex justify-between items-start mb-4">
        <div className="flex items-center space-x-3">
           <div className={`relative w-3 h-3 rounded-full ${statusDot} flex-shrink-0`}>
              {isOnline && <div className="absolute inset-0 rounded-full bg-emerald-400 animate-ping opacity-75"></div>}
           </div>
           <div>
               <h3 className="text-lg font-bold text-slate-100 group-hover:text-blue-400 transition-colors truncate max-w-[180px]">
                 {device.name}
               </h3>
               <div className="flex items-center space-x-2 text-xs text-slate-500 mt-0.5">
                   <span className="uppercase font-mono tracking-wide">{device.protocol}</span>
                   <span>•</span>
                   <span className="font-mono">{device.id.substring(0, 8)}</span>
               </div>
           </div>
        </div>
        
        {/* Actions Menu (Simplified as horizontal buttons for now) */}
        <div className="flex space-x-1 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
            <button
                onClick={() => onEdit(device)}
                className="p-1.5 text-slate-400 hover:text-blue-400 hover:bg-blue-500/10 rounded-lg transition-colors"
                title="Edit Device"
            >
                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                </svg>
            </button>
            <button
                onClick={() => onDelete(device.id)}
                className="p-1.5 text-slate-400 hover:text-red-400 hover:bg-red-500/10 rounded-lg transition-colors"
                title="Delete Device"
            >
                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
            </button>
        </div>
      </div>

      <p className="text-sm text-slate-400 mb-6 line-clamp-2 h-10">
          {device.description || <span className="text-slate-600 italic">No description provided</span>}
      </p>

      {/* Metrics / Info Grid */}
      <div className="grid grid-cols-2 gap-4 mb-6">
          <div className="bg-slate-900/50 p-2 rounded-lg border border-slate-700/50">
             <div className="text-[10px] text-slate-500 uppercase font-bold tracking-wider mb-1">Retry</div>
             <div className="text-sm font-mono text-slate-300">{device.retry_count} / {device.retry_delay_ms}ms</div>
          </div>
          <div className="bg-slate-900/50 p-2 rounded-lg border border-slate-700/50">
             <div className="text-[10px] text-slate-500 uppercase font-bold tracking-wider mb-1">Last seen</div>
             <div className="text-sm font-mono text-slate-300">
                {device.updated_at ? new Date(device.updated_at).toLocaleTimeString() : '-'}
             </div>
          </div>
      </div>

      <div className="flex gap-3 pt-4 border-t border-slate-700/50">
        <button
          onClick={() => onTestConnection(device.id, device.protocol)}
          className="flex-1 px-3 py-2 bg-slate-700 hover:bg-slate-600 text-slate-200 text-xs font-medium rounded-lg transition-colors flex items-center justify-center gap-2"
        >
          <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
          </svg>
          Test Conn
        </button>

        <button
          onClick={() => onToggleStatus(device.id, device.status)}
          className={`flex-1 px-3 py-2 text-xs font-medium rounded-lg transition-colors flex items-center justify-center gap-2 border ${
               isOnline 
               ? 'border-amber-500/30 text-amber-400 hover:bg-amber-500/10' 
               : 'border-emerald-500/30 text-emerald-400 hover:bg-emerald-500/10'
          }`}
        >
          {isOnline ? (
               <>
                 <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                 </svg>
                 Pause
               </>
          ) : (
               <>
                 <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                 </svg>
                 Activate
               </>
          )}
        </button>
      </div>
    </div>
  );
}
