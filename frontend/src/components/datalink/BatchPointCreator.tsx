
import { useState, useMemo } from 'react';
import type { ProtocolType, DataType, Point } from '../../types/datalink';
import { addressParser } from '../../utils/addressParser';
import { useToast } from '../../contexts/ToastContext';
import { logger } from '../../utils/logger';

/** 批量建立點位 API 回應結構 */
interface BatchCreateResponse {
  created_count?: number;
  points?: Point[];
}

export interface BatchPointCreatorProps {
  deviceId: string;
  protocol: ProtocolType;
  preselectedAddresses?: string[];
  pollingGroups: { id: string; name: string }[];
  /** 從規劃「套用到網格」帶入的命名模板，例如 SRC-{index03} */
  initialTemplate?: string;
  onCreated: (points: Point[]) => void;
  onCancel: () => void;
}

interface GeneratedPoint {
  name: string;
  address: string;
}

export function BatchPointCreator({
  deviceId,
  protocol,
  preselectedAddresses = [],
  pollingGroups,
  initialTemplate,
  onCreated,
  onCancel
}: BatchPointCreatorProps) {
  const [template, setTemplate] = useState(initialTemplate ?? 'Pump_{index}');
  const [dataType, setDataType] = useState<DataType>('int16');
  const [pollingGroupId, setPollingGroupId] = useState(pollingGroups[0]?.id || '');
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Generate preview
  const previewPoints = useMemo<GeneratedPoint[]>(() => {
    return preselectedAddresses.map((addr, idx) => {
      let name = template;
      name = name.replace(/\{index03\}/g, (idx + 1).toString().padStart(3, '0'));
      name = name.replace(/\{index\}/g, idx.toString());
      name = name.replace(/\{address\}/g, addr);
      
      // Try to extract number if possible for MC/Fatek
      const parsed = addressParser.validate(addr, protocol);
      if (parsed.valid) {
          const p = addressParser.parse(addr, protocol);
          name = name.replace(/\{num\}/g, p.startNumber.toString());
      }

      return {
        name,
        address: addr
      };
    });
  }, [template, preselectedAddresses, protocol]);

  const { showSuccess, showError } = useToast();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (previewPoints.length === 0) return;
    
    setIsSubmitting(true);
    try {
      const payload = {
        device_id: deviceId,
        polling_group_id: pollingGroupId,
        data_type: dataType,
        enabled: true,
        points: previewPoints
      };
      
      const response = await fetch('/api/v1/datalink/points/batch', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });
      
      const result = await response.json() as { error?: { message?: string }; data?: BatchCreateResponse };
      if (!response.ok) throw new Error(result.error?.message || '批量建立失敗');

      const data = result.data ?? {};
      showSuccess(`成功建立 ${data.created_count ?? previewPoints.length} 個點位`);
      onCreated(data.points ?? []);
    } catch (err) {
      logger.error(err);
      showError('建立失敗: ' + (err instanceof Error ? err.message : '未知錯誤'));
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form 
      onSubmit={handleSubmit} 
      className="space-y-6" 
      data-testid="batch-point-creator"
    >
      <div className="space-y-4">
        {/* Template Input */}
        <div>
          <label htmlFor="template" className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
            命名模板
          </label>
          <input
            id="template"
            type="text"
            value={template}
            onChange={(e) => setTemplate(e.target.value)}
            className="w-full px-3 py-2 bg-white dark:bg-slate-800 border border-slate-300 dark:border-slate-600 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 font-mono text-sm"
            placeholder="例如: Pump_{index}"
            required
          />
          <p className="mt-1 text-xs text-slate-500">
            支援變數: {'{index}'} (序號), {'{index03}'} (三位數 1 起), {'{address}'} (位址), {'{num}'} (數值部分)
          </p>
        </div>

        {/* Data Type */}
        <div className="grid grid-cols-2 gap-4">
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
          
          <div>
            <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-1">
              輪詢群組
            </label>
            <select
              value={pollingGroupId}
              onChange={(e) => setPollingGroupId(e.target.value)}
              className="w-full px-3 py-2 bg-white dark:bg-slate-800 border border-slate-300 dark:border-slate-600 rounded-md shadow-sm focus:ring-2 focus:ring-blue-500 text-sm"
            >
              {pollingGroups.map(g => (
                <option key={g.id} value={g.id}>{g.name}</option>
              ))}
            </select>
          </div>
        </div>

        {/* Preview List */}
        <div>
          <label className="block text-sm font-medium text-slate-700 dark:text-slate-300 mb-2">
            預覽建立 ({previewPoints.length})
          </label>
          <div className="max-h-60 overflow-y-auto border border-slate-200 dark:border-slate-700 rounded-lg bg-slate-50 dark:bg-slate-900/50">
            <table className="w-full text-xs text-left">
              <thead className="sticky top-0 bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400">
                <tr>
                  <th className="px-3 py-2 font-semibold">名稱</th>
                  <th className="px-3 py-2 font-semibold">位址</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-200 dark:divide-slate-700">
                {previewPoints.map((p, i) => (
                  <tr key={i}>
                    <td className="px-3 py-2 font-medium text-slate-800 dark:text-slate-200">{p.name}</td>
                    <td className="px-3 py-2 font-mono text-slate-500 dark:text-slate-400">{p.address}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <div className="flex items-center gap-3">
        <button
          type="button"
          onClick={onCancel}
          className="flex-1 px-4 py-2 border border-slate-300 dark:border-slate-600 text-slate-700 dark:text-slate-300 rounded-lg hover:bg-slate-50 dark:hover:bg-slate-800 transition-colors font-medium"
        >
          取消
        </button>
        <button
          type="submit"
          disabled={isSubmitting || previewPoints.length === 0}
          className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white rounded-lg transition-colors font-medium shadow-sm"
        >
          {isSubmitting ? '建立中...' : '確認建立'}
        </button>
      </div>
    </form>
  );
}
