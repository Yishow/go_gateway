import { useState, useEffect } from 'react';
import { mappingAPI } from '../../services/datalink';
import type { Mapping, CreateMappingRequest, UpdateMappingRequest } from '../../types/datalink';
import MappingCard from '../../components/datalink/MappingCard';
import MappingCanvas from '../../components/datalink/MappingCanvas';

export default function MappingsPage() {
  const [mappings, setMappings] = useState<Mapping[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Modal states for Canvas Editor
  const [isEditorOpen, setIsEditorOpen] = useState(false);
  const [editingMapping, setEditingMapping] = useState<Mapping | null>(null);

  useEffect(() => {
    fetchMappings();
  }, []);

  const fetchMappings = async () => {
    setLoading(true);
    try {
      const data = await mappingAPI.list();
      setMappings(data);
      setError(null);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch mappings');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Are you sure you want to delete this mapping?')) return;
    
    try {
      await mappingAPI.delete(id);
      setMappings(prev => prev.filter(m => m.id !== id));
    } catch (err: any) {
      alert(`Failed to delete mapping: ${err.message}`);
    }
  };

  const handleToggleStatus = async (id: string, enabled: boolean) => {
    try {
      // Create update request with partial data - assuming backend handles partial updates or we send full object
      const mapping = mappings.find(m => m.id === id);
      if (!mapping) return;

      const updated = await mappingAPI.update(id, { 
          enabled: !enabled 
      });

      setMappings(prev => prev.map(m => m.id === id ? updated : m));
    } catch (err: any) {
      alert(`Failed to update status: ${err.message}`);
    }
  };

  const handleEdit = (mapping: Mapping) => {
    setEditingMapping(mapping);
    setIsEditorOpen(true);
  };

  const handleCreate = () => {
    setEditingMapping(null);
    setIsEditorOpen(true);
  };

  const handleSave = async (data: CreateMappingRequest | UpdateMappingRequest) => {
    try {
      if (editingMapping) {
        const updated = await mappingAPI.update(editingMapping.id, data as UpdateMappingRequest);
        setMappings(prev => prev.map(m => m.id === editingMapping.id ? updated : m));
      } else {
        const created = await mappingAPI.create(data as CreateMappingRequest);
        setMappings(prev => [...prev, created]);
      }
      setIsEditorOpen(false);
    } catch (err: any) {
      alert(`Failed to save mapping: ${err.message}`);
    }
  };

  return (
    <div className="space-y-8">
      {/* Header with Glow */}
      <div className="flex justify-between items-center relative z-0">
        <div className="absolute -left-10 -top-10 w-32 h-32 bg-emerald-600/20 blur-3xl pointer-events-none"></div>
        <div>
          <h2 className="text-3xl font-bold text-slate-100 tracking-tight">Data Mappings</h2>
          <p className="text-slate-400 mt-1">Design data transformation pipelines from Source to Destination</p>
        </div>
        <button
          onClick={handleCreate}
          className="px-5 py-2.5 bg-blue-600 hover:bg-blue-500 text-white rounded-xl font-medium transition-all shadow-lg shadow-blue-500/20 hover:shadow-blue-500/40 flex items-center space-x-2"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          <span>New Mapping</span>
        </button>
      </div>

      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[1, 2, 3].map(i => (
                <div key={i} className="h-40 bg-slate-800/50 rounded-xl border border-slate-700/50 animate-pulse"></div>
            ))}
        </div>
      ) : error ? (
        <div className="bg-red-500/10 border border-red-500/20 rounded-xl p-6 text-red-400 flex items-center space-x-3">
          <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
           <span>{error}</span>
        </div>
      ) : mappings.length === 0 ? (
        <div className="text-center py-20 bg-slate-800/30 rounded-2xl border-2 border-slate-700/50 border-dashed group hover:border-emerald-500/30 transition-colors cursor-pointer" onClick={handleCreate}>
           <div className="w-16 h-16 bg-slate-700/50 rounded-full mx-auto flex items-center justify-center mb-4 group-hover:scale-110 transition-transform duration-300">
              <svg className="w-8 h-8 text-slate-500 group-hover:text-emerald-400 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
              </svg>
           </div>
           <h3 className="text-xl font-bold text-slate-200 mb-2">No mappings found</h3>
           <p className="text-slate-500 max-w-sm mx-auto mb-6">
             Create your first mapping pipeline to connect device data points to standardized tags.
           </p>
           <button
             className="text-emerald-400 hover:text-emerald-300 font-medium hover:underline"
           >
             Start mapping &rarr;
           </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {mappings.map(mapping => (
            <MappingCard
              key={mapping.id}
              mapping={mapping}
              onEdit={handleEdit}
              onDelete={handleDelete}
              onToggleStatus={handleToggleStatus}
            />
          ))}
        </div>
      )}

      {/* Editor Modal - Full Screen Overlay */}
      {isEditorOpen && (
        <div className="fixed inset-0 bg-slate-900 z-50 overflow-hidden flex flex-col animate-in fade-in duration-200">
            <MappingCanvas 
                initialMapping={editingMapping}
                onSave={handleSave}
                onCancel={() => setIsEditorOpen(false)}
            />
        </div>
      )}
    </div>
  );
}
