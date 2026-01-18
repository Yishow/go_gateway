import type { Tag } from '../../types/datalink';

interface TagTableProps {
  tags: Tag[];
  onEdit: (tag: Tag) => void;
  onDelete: (id: string) => void;
  onActivate: (id: string) => void;
  onRetire: (id: string) => void;
}

export default function TagTable({
  tags,
  onEdit,
  onDelete,
  onActivate,
  onRetire,
}: TagTableProps) {
  
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20';
      case 'retired': return 'bg-slate-500/10 text-slate-400 border-slate-500/20';
      case 'draft': return 'bg-amber-500/10 text-amber-400 border-amber-500/20';
      default: return 'bg-slate-500/10 text-slate-400 border-slate-500/20';
    }
  };

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-left text-sm text-slate-400">
        <thead className="bg-slate-900/50 text-xs uppercase font-medium text-slate-500">
          <tr>
            <th className="px-6 py-4 font-semibold tracking-wider">Key / Name</th>
            <th className="px-6 py-4 font-semibold tracking-wider">Data Type</th>
            <th className="px-6 py-4 font-semibold tracking-wider">Unit</th>
            <th className="px-6 py-4 font-semibold tracking-wider">Labels</th>
            <th className="px-6 py-4 font-semibold tracking-wider">Status</th>
            <th className="px-6 py-4 text-right font-semibold tracking-wider">Actions</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-700/50">
          {tags.map((tag) => (
            <tr key={tag.id} className="hover:bg-slate-700/30 transition-colors group">
              <td className="px-6 py-4">
                <div className="font-mono font-medium text-slate-200">{tag.key}</div>
                <div className="text-xs text-slate-500 mt-0.5">{tag.name}</div>
              </td>
              <td className="px-6 py-4">
                 <span className="px-2 py-1 rounded bg-slate-700/50 font-mono text-xs border border-slate-600/50 text-slate-300">
                    {tag.data_type}
                 </span>
              </td>
              <td className="px-6 py-4 text-slate-300">
                  {tag.unit || '-'}
              </td>
              <td className="px-6 py-4">
                 <div className="flex flex-wrap gap-1">
                    {Object.entries((tag.labels as any) || {}).map(([Key, Value]) => (
                        <span key={Key} className="text-[10px] px-1.5 py-0.5 rounded bg-blue-500/10 text-blue-400 border border-blue-500/20">
                            {Key}:{String(Value)}
                        </span>
                    ))}
                    {!tag.labels && <span className="text-slate-600">-</span>}
                 </div>
              </td>
              <td className="px-6 py-4">
                <span className={`px-2.5 py-0.5 rounded-full text-xs font-medium border ${getStatusColor(tag.status)}`}>
                  {tag.status}
                </span>
              </td>
              <td className="px-6 py-4 text-right">
                <div className="flex items-center justify-end space-x-2 opacity-0 group-hover:opacity-100 transition-opacity">
                    {tag.status !== 'retired' && (
                        <button
                            onClick={() => onEdit(tag)}
                            className="p-1 text-slate-400 hover:text-blue-400 hover:bg-blue-500/10 rounded transition-colors"
                            title="Edit"
                        >
                            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                            </svg>
                        </button>
                    )}
                    
                    {tag.status === 'draft' && (
                        <button
                            onClick={() => onActivate(tag.id)}
                            className="p-1 text-slate-400 hover:text-emerald-400 hover:bg-emerald-500/10 rounded transition-colors"
                            title="Activate"
                        >
                            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                            </svg>
                        </button>
                    )}

                    {tag.status === 'active' && (
                        <button
                            onClick={() => onRetire(tag.id)}
                             className="p-1 text-slate-400 hover:text-amber-400 hover:bg-amber-500/10 rounded transition-colors"
                            title="Retire"
                        >
                             <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                            </svg>
                        </button>
                    )}

                    <button
                        onClick={() => onDelete(tag.id)}
                         className="p-1 text-slate-400 hover:text-red-400 hover:bg-red-500/10 rounded transition-colors"
                        title="Delete"
                    >
                         <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                    </button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
