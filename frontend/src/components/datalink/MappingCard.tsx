import type { Mapping } from '../../types/datalink';

interface MappingCardProps {
  mapping: Mapping;
  onEdit: (mapping: Mapping) => void;
  onDelete: (id: string) => void;
  onToggleStatus: (id: string, enabled: boolean) => void;
}

interface MappingStep {
  type?: string;
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

  const isEnabled = mapping.enabled;

  return (
    <div className="group bg-slate-800 rounded-xl border border-slate-700 p-5 hover:border-violet-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-violet-500/5 hover:-translate-y-1 relative overflow-hidden">
      {/* Top Status Bar */}
      <div className={`absolute top-0 left-0 w-full h-1 ${isEnabled ? 'bg-violet-500' : 'bg-transparent'}`}></div>

      <div className="flex justify-between items-start mb-6">
        <div className="flex items-center space-x-3">
           <div className={`w-3 h-3 rounded-full ${isEnabled ? 'bg-violet-400 shadow-[0_0_8px_rgba(167,139,250,0.6)]' : 'bg-slate-600'} flex-shrink-0`}></div>
           <div>
               <h3 className="text-sm font-bold text-slate-100 font-mono tracking-wide">
                 {mapping.id.substring(0, 8)}
               </h3>
               <span className="text-[10px] uppercase font-bold text-slate-500 tracking-wider">Mapping ID</span>
           </div>
        </div>
        
        <div className="flex space-x-1 opacity-100 lg:opacity-0 group-hover:opacity-100 transition-opacity">
            <button
                onClick={() => onEdit(mapping)}
                className="p-1.5 text-slate-400 hover:text-violet-400 hover:bg-violet-500/10 rounded-lg transition-colors"
                title="Edit Mapping"
            >
                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                </svg>
            </button>
            <button
                onClick={() => onDelete(mapping.id)}
                className="p-1.5 text-slate-400 hover:text-red-400 hover:bg-red-500/10 rounded-lg transition-colors"
                title="Delete Mapping"
            >
                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
            </button>
        </div>
      </div>

      {/* Flow Visualization */}
      <div className="flex items-center justify-between mb-6 relative">
          {/* Connector Line */}
          <div className="absolute left-10 right-10 top-1/2 h-0.5 bg-slate-700 -z-10"></div>

          {/* Source Point */}
          <div className="flex flex-col items-center z-10">
              <div className="w-10 h-10 rounded-full bg-slate-700 flex items-center justify-center border-2 border-slate-600 mb-2 shadow-sm text-slate-300">
                  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                  </svg>
              </div>
              <div className="text-[10px] font-mono text-slate-400 bg-slate-900 px-1.5 py-0.5 rounded border border-slate-700 max-w-[80px] truncate">
                  {mapping.point_id.substring(0, 8)}...
              </div>
          </div>

          {/* Transform Icon */}
           <div className="bg-slate-800 p-1.5 rounded-full border border-slate-600 z-10">
                <svg className="w-4 h-4 text-violet-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 5l7 7-7 7M5 5l7 7-7 7" />
                </svg>
           </div>

           {/* Dest Tag */}
          <div className="flex flex-col items-center z-10">
              <div className="w-10 h-10 rounded-full bg-slate-700 flex items-center justify-center border-2 border-slate-600 mb-2 shadow-sm text-slate-300">
                  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                     <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                  </svg>
              </div>
              <div className="text-[10px] font-mono text-slate-400 bg-slate-900 px-1.5 py-0.5 rounded border border-slate-700 max-w-[80px] truncate">
                  {mapping.tag_id.substring(0, 8)}...
              </div>
          </div>
      </div>

      <div className="bg-slate-900/50 rounded-lg p-3 mb-4 border border-slate-700/50 min-h-[60px]">
         <div className="text-[10px] uppercase font-bold text-slate-500 mb-2 tracking-wider">Pipeline</div>
         <div className="flex flex-wrap gap-2">
             {Array.isArray(pipeline) && pipeline.length > 0 ? (
                 pipeline.map((step: MappingStep, idx: number) => (
                     <span key={idx} className="px-2 py-1 bg-violet-500/10 text-violet-300 text-[10px] font-medium rounded border border-violet-500/20 flex items-center">
                         {step.type}
                     </span>
                 ))
             ) : (
                 <span className="text-xs text-slate-500 italic flex items-center gap-1">
                    <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                    </svg>
                    Direct Pass-through
                 </span>
             )}
         </div>
      </div>

      <button
        onClick={() => onToggleStatus(mapping.id, mapping.enabled)}
        className={`w-full py-2 text-xs font-bold rounded-lg transition-all uppercase tracking-wider ${
            isEnabled 
            ? 'bg-slate-700 hover:bg-slate-600 text-slate-300' 
            : 'bg-emerald-500 hover:bg-emerald-600 text-white shadow-lg shadow-emerald-500/20'
        }`}
      >
        {isEnabled ? 'Pause Mapping' : 'Activate Mapping'}
      </button>
    </div>
  );
}
