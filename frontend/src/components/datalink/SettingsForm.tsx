import { useState } from 'react';
import type { SystemSettings, UpdateSystemSettingsRequest } from '../../types/datalink';

interface SettingsFormProps {
  settings: SystemSettings;
  onSubmit: (data: UpdateSystemSettingsRequest) => Promise<void>;
}

export default function SettingsForm({ settings, onSubmit }: SettingsFormProps) {
  const [writePrecision, setWritePrecision] = useState(settings.write_precision);
  const [partitionInterval, setPartitionInterval] = useState(settings.partition_interval);
  const [batchSize, setBatchSize] = useState(settings.batch_size);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      await onSubmit({
        write_precision: writePrecision,
        partition_interval: partitionInterval,
        batch_size: Number(batchSize),
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <label className="block text-sm font-medium text-slate-300 mb-1">
            Write Timestamp Precision
          </label>
          <select
            value={writePrecision}
            onChange={(e) => setWritePrecision(e.target.value as SystemSettings['write_precision'])}
            className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
          >
            <option value="second">Seconds</option>
            <option value="millisecond">Milliseconds</option>
          </select>
          <p className="mt-1 text-xs text-slate-500">
            Precision of timestamps stored in the database.
          </p>
        </div>

        <div>
           <label className="block text-sm font-medium text-slate-300 mb-1">
            Partition Interval
          </label>
          <select
            value={partitionInterval}
            onChange={(e) => setPartitionInterval(e.target.value as SystemSettings['partition_interval'])}
             className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
          >
            <option value="daily">Daily</option>
            <option value="weekly">Weekly</option>
            <option value="monthly">Monthly</option>
          </select>
          <p className="mt-1 text-xs text-slate-500">
             Frequency of data table partitioning.
          </p>
        </div>

        <div>
          <label className="block text-sm font-medium text-slate-300 mb-1">
            Batch Size
          </label>
          <input
            type="number"
            value={batchSize}
            onChange={(e) => setBatchSize(Number(e.target.value))}
            min={1}
            max={10000}
            className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
          />
           <p className="mt-1 text-xs text-slate-500">
             Number of records to insert in a single transaction.
          </p>
        </div>
      </div>

      <div className="flex justify-end pt-4 border-t border-slate-700">
        <button
          type="submit"
          disabled={loading}
          className="px-6 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-lg font-medium transition-colors disabled:opacity-50"
        >
          {loading ? 'Saving...' : 'Save Changes'}
        </button>
      </div>
    </form>
  );
}
