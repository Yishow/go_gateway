import { useState } from 'react';
import type {
  PollingGroup,
  CreatePollingGroupRequest,
  UpdatePollingGroupRequest,
} from '../../types/datalink';

interface PollingGroupFormProps {
  pollingGroup?: PollingGroup;
  onSubmit: (data: CreatePollingGroupRequest | UpdatePollingGroupRequest) => Promise<void>;
  onCancel: () => void;
}

/**
 * PollingGroup 表單組件
 *
 * 用於建立或編輯輪詢群組，包含名稱、描述、輪詢間隔、優先級等欄位。
 */
export default function PollingGroupForm({
  pollingGroup,
  onSubmit,
  onCancel,
}: PollingGroupFormProps) {

  // Form State
  const [name, setName] = useState(pollingGroup?.name || '');
  const [description, setDescription] = useState(pollingGroup?.description || '');
  const [intervalMs, setIntervalMs] = useState(pollingGroup?.interval_ms || 1000);
  const [priority, setPriority] = useState(pollingGroup?.priority || 100);
  const [enabled, setEnabled] = useState(pollingGroup?.enabled ?? true);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  /**
   * 處理表單提交
   */
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    // 基本驗證
    if (!name.trim()) {
      setError('名稱不能為空');
      setLoading(false);
      return;
    }
    if (intervalMs < 100) {
      setError('輪詢間隔不能小於 100ms');
      setLoading(false);
      return;
    }

    try {
      if (pollingGroup) {
        await onSubmit({
          name,
          description,
          interval_ms: intervalMs,
          priority,
          enabled,
        } as UpdatePollingGroupRequest);
      } else {
        await onSubmit({
          name,
          description,
          interval_ms: intervalMs,
          priority,
          enabled,
        } as CreatePollingGroupRequest);
      }
    } catch (err: unknown) {
      const message =
        typeof err === 'object' && err !== null && 'message' in err
          ? String((err as { message?: string }).message ?? '儲存失敗')
          : '儲存失敗';
      setError(message);
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

      {/* Name */}
      <div>
        <label className="block text-sm font-medium text-slate-300 mb-1">
          名稱 <span className="text-red-400">*</span>
        </label>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
          required
          placeholder="e.g. High Frequency Group"
          disabled={loading}
        />
      </div>

      {/* Description */}
      <div>
        <label className="block text-sm font-medium text-slate-300 mb-1">描述</label>
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          rows={3}
          className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 text-sm focus:outline-none focus:border-blue-500"
          placeholder="可選的描述..."
          disabled={loading}
        />
      </div>

      {/* Interval & Priority (Row) */}
      <div className="grid grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-slate-300 mb-1">
            輪詢間隔 (ms) <span className="text-red-400">*</span>
          </label>
          <input
            type="number"
            value={intervalMs}
            onChange={(e) => setIntervalMs(parseInt(e.target.value) || 1000)}
            min={100}
            step={100}
            className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
            required
            disabled={loading}
          />
          <p className="text-xs text-slate-500 mt-1">最小值: 100ms</p>
        </div>
        <div>
          <label className="block text-sm font-medium text-slate-300 mb-1">優先級</label>
          <input
            type="number"
            value={priority}
            onChange={(e) => setPriority(parseInt(e.target.value) || 100)}
            min={1}
            max={1000}
            className="w-full bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-slate-200 focus:outline-none focus:border-blue-500"
            disabled={loading}
          />
          <p className="text-xs text-slate-500 mt-1">數值越大優先級越高 (1-1000)</p>
        </div>
      </div>

      {/* Enabled */}
      <div className="flex items-center pt-2">
        <label className="flex items-center space-x-3 cursor-pointer">
          <div className="relative">
            <input
              type="checkbox"
              className="sr-only"
              checked={enabled}
              onChange={(e) => setEnabled(e.target.checked)}
              disabled={loading}
            />
            <div
              className={`w-10 h-6 rounded-full shadow-inner transition-colors ${
                enabled ? 'bg-blue-600' : 'bg-slate-700'
              }`}
            ></div>
            <div
              className={`absolute left-1 top-1 w-4 h-4 bg-white rounded-full transition-transform ${
                enabled ? 'translate-x-4' : 'translate-x-0'
              }`}
            ></div>
          </div>
          <span className="text-sm font-medium text-slate-300">啟用</span>
        </label>
      </div>

      {/* Form Actions */}
      <div className="flex justify-end space-x-3 pt-6 border-t border-slate-700 mt-6">
        <button
          type="button"
          onClick={onCancel}
          className="px-4 py-2 text-slate-300 hover:text-white transition-colors"
          disabled={loading}
        >
          取消
        </button>
        <button
          type="submit"
          className="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          disabled={loading}
        >
          {loading ? '儲存中...' : pollingGroup ? '更新群組' : '建立群組'}
        </button>
      </div>
    </form>
  );
}
