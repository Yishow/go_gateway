import { useState, useEffect } from 'react';
import { tagAPI } from '../../services/datalink';
import type { Tag, CreateTagRequest, UpdateTagRequest } from '../../types/datalink';
import TagTable from '../../components/datalink/TagTable';
import TagForm from '../../components/datalink/TagForm';

export default function TagsPage() {
  const [tags, setTags] = useState<Tag[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Modal states
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingTag, setEditingTag] = useState<Tag | null>(null);

  // Search/Filter states
  const [searchKey, setSearchKey] = useState('');
  const [filterStatus, setFilterStatus] = useState<string>('');

  useEffect(() => {
    fetchTags();
  }, [searchKey, filterStatus]);

  const fetchTags = async () => {
    setLoading(true);
    try {
      // In a real app, debounce search and pass params
      const data = await tagAPI.list({
          key_prefix: searchKey || undefined,
          status: filterStatus || undefined,
      });
      setTags(data);
      setError(null);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch tags');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Are you sure you want to delete this tag?')) return;
    
    try {
      await tagAPI.delete(id);
      setTags(prev => prev.filter(t => t.id !== id));
    } catch (err: any) {
      alert(`Failed to delete tag: ${err.message}`);
    }
  };

  const handleActivate = async (id: string) => {
    try {
      const updated = await tagAPI.activate(id);
      setTags(prev => prev.map(t => t.id === id ? updated : t));
    } catch (err: any) {
      alert(`Failed to activate tag: ${err.message}`);
    }
  };

  const handleRetire = async (id: string) => {
    if (!window.confirm('Are you sure you want to retire this tag? retired tags cannot be used in new mappings.')) return;

    try {
      const updated = await tagAPI.retire(id);
      setTags(prev => prev.map(t => t.id === id ? updated : t));
    } catch (err: any) {
      alert(`Failed to retire tag: ${err.message}`);
    }
  };

  const handleEdit = (tag: Tag) => {
    setEditingTag(tag);
    setIsModalOpen(true);
  };

  const handleCreate = () => {
    setEditingTag(null);
    setIsModalOpen(true);
  };

  const handleSubmit = async (data: CreateTagRequest | UpdateTagRequest) => {
    try {
      if (editingTag) {
        const updated = await tagAPI.update(editingTag.id, data as UpdateTagRequest);
        setTags(prev => prev.map(t => t.id === editingTag.id ? updated : t));
      } else {
        const created = await tagAPI.create(data as CreateTagRequest);
        setTags(prev => [...prev, created]);
      }
      setIsModalOpen(false);
    } catch (err: any) {
      throw err;
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold text-slate-100">Tag Dictionary</h2>
          <p className="text-slate-400">Manage global standard tags and metadata</p>
        </div>
        <button
          onClick={handleCreate}
          className="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-lg font-medium transition-colors flex items-center space-x-2"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          <span>Create Tag</span>
        </button>
      </div>

      {/* Filters */}
      <div className="flex space-x-4 bg-slate-800 p-4 rounded-lg border border-slate-700">
          <div className="flex-1">
              <input
                type="text"
                placeholder="Search by tag key..."
                value={searchKey}
                onChange={e => setSearchKey(e.target.value)}
                className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
              />
          </div>
          <div>
              <select
                value={filterStatus}
                onChange={e => setFilterStatus(e.target.value)}
                className="bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
              >
                  <option value="">All Status</option>
                  <option value="draft">Draft</option>
                  <option value="active">Active</option>
                  <option value="retired">Retired</option>
              </select>
          </div>
      </div>

      {loading ? (
        <div className="text-center py-12">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500 mx-auto"></div>
          <p className="mt-4 text-slate-400">Loading tags...</p>
        </div>
      ) : error ? (
        <div className="bg-red-500/10 border border-red-500/20 rounded-lg p-4 text-red-400">
          Error: {error}
        </div>
      ) : tags.length === 0 ? (
        <div className="text-center py-12 bg-slate-800/50 rounded-xl border border-slate-700 border-dashed">
          <p className="text-slate-400 mb-4">No tags found</p>
          <button
            onClick={handleCreate}
            className="text-blue-400 hover:text-blue-300 font-medium"
          >
            Create your first tag
          </button>
        </div>
      ) : (
        <TagTable 
            tags={tags}
            onEdit={handleEdit}
            onDelete={handleDelete}
            onActivate={handleActivate}
            onRetire={handleRetire}
        />
      )}

      {/* Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50 backdrop-blur-sm">
          <div className="bg-slate-800 rounded-xl border border-slate-700 p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto shadow-2xl">
             <h3 className="text-xl font-bold text-slate-100 mb-6 pb-4 border-b border-slate-700">
                {editingTag ? 'Edit Tag' : 'New Tag'}
             </h3>
             
             <TagForm
                tag={editingTag || undefined}
                onSubmit={handleSubmit}
                onCancel={() => setIsModalOpen(false)}
             />
          </div>
        </div>
      )}
    </div>
  );
}
