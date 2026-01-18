import { useState, useEffect } from 'react';
import type { Tag, CreateTagRequest, UpdateTagRequest, DataType } from '../../types/datalink';
import { tagAPI } from '../../services/datalink';

interface TagFormProps {
  tag?: Tag;
  onSubmit: (data: CreateTagRequest | UpdateTagRequest) => Promise<void>;
  onCancel: () => void;
}

const DATA_TYPES: DataType[] = ['bool', 'int16', 'int32', 'int64', 'uint16', 'uint32', 'uint64', 'float32', 'float64', 'string'];

export default function TagForm({ tag, onSubmit, onCancel }: TagFormProps) {
  const [key, setKey] = useState(tag?.key || '');
  const [name, setName] = useState(tag?.name || '');
  const [description, setDescription] = useState(tag?.description || '');
  const [dataType, setDataType] = useState<DataType>(tag?.data_type || 'float64');
  const [unit, setUnit] = useState(tag?.unit || '');
  
  // Key validation state
  const [keyError, setKeyError] = useState<string | null>(null);
  const [isValidating, setIsValidating] = useState(false);

  // Labels (Simple key-value pairs for now)
  const [labelsInput, setLabelsInput] = useState(
    tag?.labels ? Object.entries(tag.labels).map(([k,v]) => `${k}=${v}`).join('\n') : ''
  );

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Validate key on change (debounce could be added here)
  useEffect(() => {
    if (!tag && key) {
        const timer = setTimeout(async () => {
            setIsValidating(true);
            try {
                const result = await tagAPI.validateKey(key);
                if (!result.valid) {
                    setKeyError(result.error || 'Invalid key format');
                } else if (result.exists) {
                    setKeyError('Tag key already exists');
                } else {
                    setKeyError(null);
                }
            } catch (err) {
                console.error('Validation failed', err);
            } finally {
                setIsValidating(false);
            }
        }, 500);
        return () => clearTimeout(timer);
    }
  }, [key, tag]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    // Parse labels
    const labels: Record<string, string> = {};
    labelsInput.split('\n').forEach(line => {
        const [k, v] = line.split('=').map(s => s.trim());
        if (k && v) labels[k] = v;
    });

    try {
      if (tag) {
        await onSubmit({
          name,
          description,
          unit,
          labels,
        } as UpdateTagRequest);
      } else {
        if (keyError) throw new Error("Please fix key errors");
        
        await onSubmit({
          key,
          name,
          description,
          data_type: dataType,
          unit,
          labels,
        } as CreateTagRequest);
      }
    } catch (err: any) {
      setError(err.message || 'Failed to save tag');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <div className="bg-red-500/10 border border-red-500/20 rounded-lg p-3 text-sm text-red-400">
          {error}
        </div>
      )}

      {/* Key (Read-only on edit) */}
      <div>
        <label className="block text-sm font-medium text-slate-300">Tag Key</label>
        <div className="relative">
            <input
                type="text"
                value={key}
                onChange={e => setKey(e.target.value)}
                disabled={!!tag}
                className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500 font-mono disabled:opacity-50"
                required
                placeholder="site.area.equipment.metric"
            />
            {isValidating && (
                <div className="absolute right-3 top-3">
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-500"></div>
                </div>
            )}
        </div>
        {!tag && (
            <p className={`text-xs mt-1 ${keyError ? 'text-red-400' : 'text-slate-500'}`}>
                {keyError || 'Format: segment1.segment2.segment3...'}
            </p>
        )}
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-slate-300">Name</label>
          <input
            type="text"
            value={name}
            onChange={e => setName(e.target.value)}
            className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
            required
            placeholder="Human readable name"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-slate-300">Data Type</label>
          <select
            value={dataType}
            onChange={e => setDataType(e.target.value as DataType)}
            disabled={!!tag}
            className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500 disabled:opacity-50"
          >
            {DATA_TYPES.map(t => (
                <option key={t} value={t}>{t}</option>
            ))}
          </select>
        </div>
      </div>

      <div>
          <label className="block text-sm font-medium text-slate-300">Description</label>
          <input
            type="text"
            value={description}
            onChange={e => setDescription(e.target.value)}
            className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
          />
      </div>

      <div className="grid grid-cols-2 gap-4">
         <div>
          <label className="block text-sm font-medium text-slate-300">Unit</label>
          <input
            type="text"
            value={unit}
            onChange={e => setUnit(e.target.value)}
            className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
            placeholder="e.g. °C, kW, rpm"
          />
        </div>
      </div>

      <div>
          <label className="block text-sm font-medium text-slate-300">Labels (key=value, one per line)</label>
          <textarea
            value={labelsInput}
            onChange={e => setLabelsInput(e.target.value)}
            rows={3}
             className="mt-1 w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 font-mono text-sm focus:outline-none focus:border-blue-500"
             placeholder="area=production&#10;type=sensor"
          />
      </div>

      <div className="flex justify-end space-x-3 pt-6 border-t border-slate-700 mt-6">
        <button
          type="button"
          onClick={onCancel}
          className="px-4 py-2 text-slate-300 hover:text-white transition-colors"
          disabled={loading}
        >
          Cancel
        </button>
        <button
          type="submit"
          className="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          disabled={loading || !!keyError}
        >
          {loading ? 'Saving...' : 'Save Tag'}
        </button>
      </div>
    </form>
  );
}
