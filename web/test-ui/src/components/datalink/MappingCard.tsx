import type { Mapping } from '../../types/datalink';

interface MappingCardProps {
  mapping: Mapping;
  onEdit: (mapping: Mapping) => void;
  onDelete: (id: string) => void;
  onToggleStatus: (id: string, enabled: boolean) => void;
}

export default function MappingCard({
  mapping,
  onEdit,
  onDelete,
  onToggleStatus,
}: MappingCardProps) {
  // Parsing pipeline for display
  const pipeline = typeof mapping.transform_pipeline === 'string' 
    ? JSON.parse(mapping.transform_pipeline) 
    : mapping.transform_pipeline;

  return (
    <div className="bg-slate-800 rounded-lg border border-slate-700 p-4 hover:border-purple-500/50 transition-colors">
      <div className="flex justify-between items-start mb-4">
        <div className="flex items-center space-x-2">
           <div className={`w-2 h-2 rounded-full ${mapping.enabled ? 'bg-emerald-400' : 'bg-slate-600'}`}></div>
           <span className="text-sm text-slate-400 font-mono text-xs">{mapping.id.substring(0, 8)}</span>
        </div>
        <div className="flex space-x-2">
            <button
                onClick={() => onEdit(mapping)}
                className="text-slate-400 hover:text-purple-400"
            >
                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                </svg>
            </button>
            <button
                onClick={() => onDelete(mapping.id)}
                className="text-slate-400 hover:text-red-400"
            >
                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
            </button>
        </div>
      </div>

      <div className="flex items-center justify-between mb-4">
          <div className="text-center">
              <div className="text-xs text-slate-500 mb-1">Point ID</div>
              <div className="text-sm font-mono text-slate-300 bg-slate-900 px-2 py-1 rounded">
                  {mapping.point_id.substring(0, 8)}...
              </div>
          </div>
          <div className="text-slate-600">
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 8l4 4m0 0l-4 4m4-4H3" />
            </svg>
          </div>
          <div className="text-center">
              <div className="text-xs text-slate-500 mb-1">Tag ID</div>
             <div className="text-sm font-mono text-slate-300 bg-slate-900 px-2 py-1 rounded">
                  {mapping.tag_id.substring(0, 8)}...
              </div>
          </div>
      </div>

      <div className="bg-slate-900/50 rounded p-2 mb-4">
         <div className="text-xs text-slate-500 mb-2">Pipeline Steps</div>
         <div className="flex flex-wrap gap-2">
             {Array.isArray(pipeline) && pipeline.map((step: any, idx: number) => (
                 <span key={idx} className="px-2 py-0.5 bg-purple-500/10 text-purple-400 text-xs rounded border border-purple-500/20">
                     {step.type}
                 </span>
             ))}
             {(!Array.isArray(pipeline) || pipeline.length === 0) && (
                 <span className="text-xs text-slate-600 italic">Direct Mapping</span>
             )}
         </div>
      </div>

      <div className="flex justify-end">
         <button
            onClick={() => onToggleStatus(mapping.id, mapping.enabled)}
            className={`text-xs font-medium px-2 py-1 rounded transition-colors ${
                mapping.enabled 
                ? 'text-amber-400 hover:bg-amber-400/10' 
                : 'text-emerald-400 hover:bg-emerald-400/10'
            }`}
        >
            {mapping.enabled ? 'Pause' : 'Resume'}
        </button>
      </div>
    </div>
  );
}
