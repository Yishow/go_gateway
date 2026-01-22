
import { type Point } from '../../types/datalink';

interface PointTableProps {
  points: Point[];
  onEdit: (point: Point) => void;
  onDelete: (id: string) => void;
  onPoll?: (id: string) => void;
}

export default function PointTable({ points, onEdit, onDelete, onPoll }: PointTableProps) {
  if (points.length === 0) {
    return (
      <div className="p-8 text-center text-slate-500">
        No points found matching your criteria.
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-left border-collapse">
        <thead>
          <tr className="bg-slate-900/50 text-slate-400 text-sm uppercase tracking-wider border-b border-slate-700">
            <th className="p-4 font-medium">Name</th>
            <th className="p-4 font-medium">Address</th>
            <th className="p-4 font-medium">Type</th>
            <th className="p-4 font-medium">Device ID</th>
             <th className="p-4 font-medium">Status</th>
            <th className="p-4 font-medium text-right">Actions</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-700/50">
          {points.map((point) => (
            <tr
              key={point.id}
              className="hover:bg-slate-800/50 transition-colors group"
            >
              <td className="p-4">
                <div className="font-medium text-slate-200">{point.name}</div>
                {point.description && (
                  <div className="text-xs text-slate-500 mt-0.5 max-w-xs truncate">
                    {point.description}
                  </div>
                )}
              </td>
              <td className="p-4">
                <code className="px-2 py-1 bg-slate-900 rounded text-xs font-mono text-purple-400 border border-slate-800">
                    {point.address}
                </code>
              </td>
              <td className="p-4 text-slate-300 text-sm">
                 <span className="px-2 py-1 bg-slate-800 rounded text-xs text-slate-400 border border-slate-700">
                    {point.data_type}
                 </span>
              </td>
               <td className="p-4 text-slate-400 text-sm font-mono text-xs">
                 {point.device_id.substring(0, 8)}...
              </td>
              <td className="p-4">
                <span
                  className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${
                    point.enabled
                      ? 'bg-green-500/10 text-green-400 border-green-500/20'
                      : 'bg-slate-700/50 text-slate-400 border-slate-600'
                  }`}
                >
                  {point.enabled ? 'Enabled' : 'Disabled'}
                </span>
                 {point.last_error && (
                    <div title={point.last_error} className="mt-1 text-xs text-red-400 flex items-center gap-1 cursor-help">
                         <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                           <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                         </svg>
                         Error
                    </div>
                )}
              </td>
              <td className="p-4 text-right space-x-2">
                {onPoll && (
                    <button
                        onClick={() => onPoll(point.id)}
                        className="text-blue-400 hover:text-blue-300 transition-colors p-1"
                        title="Poll Now"
                    >
                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                        </svg>
                    </button>
                )}
                <button
                  onClick={() => onEdit(point)}
                  className="text-slate-400 hover:text-white transition-colors p-1"
                  title="Edit Point"
                >
                  <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                  </svg>
                </button>
                <button
                  onClick={() => onDelete(point.id)}
                  className="text-slate-400 hover:text-red-400 transition-colors p-1"
                  title="Delete Point"
                >
                  <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
