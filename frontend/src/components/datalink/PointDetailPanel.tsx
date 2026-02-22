import { useState } from 'react';
import type { Point, DataType } from '../../types/datalink';
import { useToast } from '../../contexts/ToastContext';

export interface PointDetailPanelProps {
  point: Point;
  onUpdate: (updatedPoint: Point) => void;
  onDelete: (id: string) => void;
  onClose?: () => void;
}

export function PointDetailPanel({
  point,
  onUpdate,
  onDelete,
  onClose: _onClose
}: PointDetailPanelProps) {
  const [name, setName] = useState(point.name);
  const [description, setDescription] = useState(point.description || '');
  const [dataType, setDataType] = useState<DataType>(point.data_type);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const { showSuccess, showError } = useToast();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    try {
      const response = await fetch(`/api/v1/datalink/points/${point.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          device_id: point.device_id,
          name,
          description,
          data_type: dataType,
          polling_group_id: point.polling_group_id || null
        })
      });
      
      const result = await response.json();
      if (!response.ok) throw new Error(result.error?.message || '更新失敗');
      
      showSuccess('點位更新成功');
      onUpdate(result.data); // 注意後端 Helpers 會將資料包裝在 .data 中
    } catch (err) {
      showError(err instanceof Error ? err.message : '未知錯誤');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="space-y-6" data-testid="point-detail-panel">
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
            位址
          </label>
          <div className="px-3 py-2 bg-slate-100 dark:bg-slate-800 rounded-md font-mono text-sm text-slate-600 dark:text-slate-400 border border-slate-200 dark:border-slate-700">
            {point.address}
          </div>
        </div>

        <div>
          <label htmlFor="name" className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
            名稱
          </label>
          <input
            id="name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full px-3 py-2 bg-white dark:bg-slate-800 border border-slate-300 dark:border-slate-600 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 text-sm"
            required
          />
        </div>

        <div>
          <label htmlFor="desc" className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
            描述
          </label>
          <textarea
            id="desc"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            rows={2}
            className="w-full px-3 py-2 bg-white dark:bg-slate-800 border border-slate-300 dark:border-slate-600 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 text-sm"
          />
        </div>

        <div className="grid grid-cols-1 gap-4">
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
              資料型別
            </label>
            <select
              value={dataType}
              onChange={(e) => setDataType(e.target.value as DataType)}
              className="w-full px-3 py-2 bg-white dark:bg-slate-800 border border-slate-300 dark:border-slate-600 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 text-sm"
            >
              <option value="bool">Boolean</option>
              <option value="int16">Int16</option>
              <option value="uint16">Uint16</option>
              <option value="int32">Int32</option>
              <option value="float32">Float32</option>
            </select>
          </div>
        </div>

        <div className="flex gap-3 pt-4 border-t border-slate-200 dark:border-slate-700">
           <button
            type="button"
            onClick={() => {
              if (confirm('確定要刪除此點位嗎？')) onDelete(point.id);
            }}
            className="px-4 py-2 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/30 rounded-lg transition-colors font-medium text-sm"
           >
             刪除點位
           </button>
           <div className="flex-1" />
           <button
            type="submit"
            disabled={isSubmitting}
            className="px-6 py-2 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white rounded-lg transition-colors font-medium shadow-sm text-sm"
          >
            {isSubmitting ? '儲存中...' : '儲存變更'}
          </button>
        </div>
      </form>
    </div>
  );
}
