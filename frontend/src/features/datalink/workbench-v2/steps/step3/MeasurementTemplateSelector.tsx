import * as React from 'react';
import { useState } from 'react';
import {
  useApplyMeasurementTemplateMutation,
  usePreviewMeasurementTemplateMutation,
  useStudioV2MeasurementTemplatesQuery,
} from '@/hooks/datalink/useStudioV2WorkspaceMeasurements';
import { SafeQueryBoundary } from '@/utils/SafeQueryBoundary';
import type { TemplateApplyPreview } from '@/types/measurement';

export interface MeasurementTemplateSelectorProps {
  deviceId?: string;
  onApplied?: () => void;
}

const MeasurementTemplateSelectorContent: React.FC<MeasurementTemplateSelectorProps> = ({
  deviceId,
  onApplied,
}) => {
  const { data: templates = [], isLoading } = useStudioV2MeasurementTemplatesQuery(true);
  const previewMutation = usePreviewMeasurementTemplateMutation();
  const applyMutation = useApplyMeasurementTemplateMutation();

  const [selectedTemplateId, setSelectedTemplateId] = useState<string>('');
  const [previewData, setPreviewData] = useState<TemplateApplyPreview | null>(null);

  const handlePreview = async () => {
    if (!selectedTemplateId || !deviceId) return;
    const res = await previewMutation.mutateAsync({
      template_id: selectedTemplateId,
      device_id: deviceId,
    });
    setPreviewData(res);
  };

  const handleApply = async () => {
    if (!previewData) return;
    await applyMutation.mutateAsync(previewData);
    setPreviewData(null);
    onApplied?.();
  };

  if (isLoading || templates.length === 0) {
    return null;
  }

  return (
    <div
      className="bg-slate-800/40 border border-slate-700/50 rounded-xl p-4 space-y-3"
      data-testid="measurement-template-selector"
    >
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-sm font-medium text-slate-200">
            推薦範本 (快速配置標準三相電表或感測器)
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            自動辨識電壓、電流、功率與累積電量單位與語意，無需手動逐項設定。
          </p>
        </div>
      </div>

      <div className="flex items-center gap-3">
        <select
          className="bg-slate-900 border border-slate-700 text-slate-200 text-xs rounded-lg px-3 py-2 focus:ring-1 focus:ring-sky-500"
          value={selectedTemplateId}
          onChange={(e) => {
            setSelectedTemplateId(e.target.value);
            setPreviewData(null);
          }}
          disabled={!deviceId}
          data-testid="template-select"
        >
          <option value="">-- 請選擇標準量測範本 --</option>
          {templates.map((tpl) => (
            <option key={tpl.id} value={tpl.id}>
              {tpl.name} ({tpl.version})
            </option>
          ))}
        </select>

        <button
          type="button"
          className="bg-slate-700 hover:bg-slate-600 disabled:opacity-50 text-slate-200 text-xs px-3 py-2 rounded-lg font-medium transition-colors"
          disabled={!selectedTemplateId || !deviceId || previewMutation.isPending}
          onClick={handlePreview}
          data-testid="preview-template-btn"
        >
          {previewMutation.isPending ? '預覽中…' : '預覽建議'}
        </button>

        {previewData && (
          <button
            type="button"
            className="bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white text-xs px-3 py-2 rounded-lg font-medium transition-colors"
            disabled={applyMutation.isPending}
            onClick={handleApply}
            data-testid="apply-template-btn"
          >
            {applyMutation.isPending ? '套用中…' : '確認套用範本'}
          </button>
        )}
      </div>

      {previewData && (
        <div
          className="bg-slate-900/60 border border-slate-700/40 rounded-lg p-3 text-xs space-y-2"
          data-testid="template-preview-summary"
        >
          <div className="text-slate-300 font-medium">
            範本預覽：{previewData.device_name} (共推薦 {previewData.proposed_items.length} 個標準量測項目)
          </div>
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 text-slate-400">
            {previewData.proposed_items.map((item) => (
              <div
                key={item.item_id}
                className="bg-slate-800/80 px-2 py-1.5 rounded border border-slate-700/50"
              >
                <div className="text-slate-200 font-mono text-[11px]">{item.name}</div>
                <div className="text-[10px] text-slate-400">
                  {item.semantic_kind} · {item.unit || '無單位'} ({item.data_type})
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export const MeasurementTemplateSelector: React.FC<MeasurementTemplateSelectorProps> = (props) => {
  return (
    <SafeQueryBoundary>
      <MeasurementTemplateSelectorContent {...props} />
    </SafeQueryBoundary>
  );
};
