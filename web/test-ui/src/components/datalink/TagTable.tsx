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
  return (
    <div className="bg-slate-800 rounded-lg border border-slate-700 overflow-hidden">
      <table className="min-w-full divide-y divide-slate-700">
        <thead className="bg-slate-900/50">
          <tr>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
              Tag Key
            </th>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
              Name
            </th>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
              Type / Unit
            </th>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
              Status
            </th>
            <th scope="col" className="relative px-6 py-3">
              <span className="sr-only">Actions</span>
            </th>
          </tr>
        </thead>
        <tbody className="bg-slate-800 divide-y divide-slate-700">
          {tags.map((tag) => (
            <tr key={tag.id} className="hover:bg-slate-700/50 transition-colors">
              <td className="px-6 py-4 whitespace-nowrap">
                <div className="text-sm font-mono font-medium text-blue-400">{tag.key}</div>
                <div className="text-xs text-slate-500">{tag.labels ? Object.entries(tag.labels).map(([k,v]) => `${k}=${v}`).join(', ') : ''}</div>
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <div className="text-sm text-slate-200">{tag.name}</div>
                <div className="text-xs text-slate-500">{tag.description}</div>
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <div className="text-sm text-slate-300">{tag.data_type}</div>
                {tag.unit && <div className="text-xs text-slate-500">{tag.unit}</div>}
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <span className={`px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full 
                  ${tag.status === 'active' ? 'bg-emerald-500/10 text-emerald-400' : 
                    tag.status === 'retired' ? 'bg-slate-600/10 text-slate-400' : 
                    'bg-amber-500/10 text-amber-400'}`}>
                  {tag.status}
                </span>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-2">
                {tag.status === 'draft' && (
                  <button onClick={() => onActivate(tag.id)} className="text-emerald-400 hover:text-emerald-300">Activate</button>
                )}
                {tag.status === 'active' && (
                  <button onClick={() => onRetire(tag.id)} className="text-amber-400 hover:text-amber-300">Retire</button>
                )}
                
                <button onClick={() => onEdit(tag)} className="text-blue-400 hover:text-blue-300">Edit</button>
                <button onClick={() => onDelete(tag.id)} className="text-red-400 hover:text-red-300">Delete</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
