import { useState, useEffect } from 'react';
import { mappingAPI } from '../../services/datalink';
import type { Mapping } from '../../types/datalink';
import MappingCard from '../../components/datalink/MappingCard';

export default function MappingsPage() {
  const [mappings, setMappings] = useState<Mapping[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

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

  const handleToggleStatus = async (id: string, currentEnabled: boolean) => {
    try {
      const updated = await mappingAPI.update(id, { enabled: !currentEnabled });
      setMappings(prev => prev.map(m => m.id === id ? updated : m));
    } catch (err: any) {
      alert(`Failed to update status: ${err.message}`);
    }
  };

  const handleEdit = (mapping: Mapping) => {
    // Navigate to canvas editor (To be implemented)
    alert(`Edit mapping ${mapping.id} (Canvas Editor Coming Soon)`);
  };

  const handleCreate = () => {
    // Navigate to canvas editor (To be implemented)
    alert('Create mapping (Canvas Editor Coming Soon)');
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold text-slate-100">Mappings</h2>
          <p className="text-slate-400">Manage data transformation pipelines</p>
        </div>
        <button
          onClick={handleCreate}
          className="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-lg font-medium transition-colors flex items-center space-x-2"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          <span>New Mapping</span>
        </button>
      </div>

      {loading ? (
        <div className="text-center py-12">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto"></div>
          <p className="mt-4 text-slate-400">Loading mappings...</p>
        </div>
      ) : error ? (
        <div className="bg-red-500/10 border border-red-500/20 rounded-lg p-4 text-red-400">
          Error: {error}
        </div>
      ) : mappings.length === 0 ? (
        <div className="text-center py-12 bg-slate-800/50 rounded-xl border border-slate-700 border-dashed">
          <p className="text-slate-400 mb-4">No mappings found</p>
          <button
            onClick={handleCreate}
            className="text-blue-400 hover:text-blue-300 font-medium"
          >
            Create your first mapping pipeline
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
    </div>
  );
}
