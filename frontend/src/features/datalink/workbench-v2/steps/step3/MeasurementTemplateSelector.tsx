import * as React from 'react';
import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  useApplyMeasurementTemplateMutation,
  usePreviewMeasurementTemplateMutation,
  useStudioV2MeasurementTemplatesQuery,
} from '@/hooks/datalink/useStudioV2WorkspaceMeasurements';
import { SafeQueryBoundary } from '@/utils/SafeQueryBoundary';
import type { TemplateApplyPreview } from '@/types/measurement';
import { getSafeErrorStateMessage } from '@/utils/typedErrors';

export interface MeasurementTemplateSelectorProps {
  deviceId?: string;
  onApplied?: () => void;
}

const MeasurementTemplateSelectorContent: React.FC<MeasurementTemplateSelectorProps> = ({
  deviceId,
  onApplied,
}) => {
  const { t } = useTranslation(['workbench-v2', 'studio-review']);
  const { data: templates = [], isLoading, isError, error: loadError, refetch } = useStudioV2MeasurementTemplatesQuery(true);
  const previewMutation = usePreviewMeasurementTemplateMutation();
  const applyMutation = useApplyMeasurementTemplateMutation();

  const [selectedTemplateId, setSelectedTemplateId] = useState<string>('');
  const [previewData, setPreviewData] = useState<TemplateApplyPreview | null>(null);
  const [error, setError] = useState<string | null>(null);
  const request = useRef(0);
  const busy = useRef(false);
  const mounted = useRef(false);
  useEffect(() => { mounted.current = true; return () => { mounted.current = false; }; }, []);

  const handlePreview = async () => {
    if (!selectedTemplateId || !deviceId || busy.current) return;
    const id = ++request.current;
    busy.current = true;
    setPreviewData(null);
    setError(null);
    try {
      const res = await previewMutation.mutateAsync({ template_id: selectedTemplateId, device_id: deviceId });
      if (mounted.current && request.current === id && res.device_id === deviceId && res.template_id === selectedTemplateId) setPreviewData(res);
    } catch (cause) {
      if (mounted.current && request.current === id) setError(getSafeErrorStateMessage(cause, t('template.preview_failed', { ns: 'studio-review' })));
    } finally { busy.current = false; }
  };

  const handleApply = async () => {
    if (!previewData || previewData.device_id !== deviceId || previewData.template_id !== selectedTemplateId || busy.current) return;
    const id = ++request.current;
    busy.current = true;
    setError(null);
    try {
      await applyMutation.mutateAsync({ ...previewData, confirmed: true });
      if (mounted.current && request.current === id) { setPreviewData(null); onApplied?.(); }
    } catch (cause) {
      if (mounted.current && request.current === id) {
        setPreviewData(null);
        setError(getSafeErrorStateMessage(cause, t('template.apply_failed', { ns: 'studio-review' })));
      }
    } finally { busy.current = false; }
  };

  if (!isLoading && !isError && templates.length === 0) {
    return null;
  }

  return (
    <details className="group" data-testid="measurement-advanced">
      <summary className="cursor-pointer text-xs font-medium text-slate-400 hover:text-slate-200 focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-500/40">
        {t('step3.advanced_measurement')}
      </summary>
      <div className="mt-2">
      {isLoading && <p role="status">{t('template.loading', { ns: 'studio-review' })}</p>}
      {isError && <div role="alert" data-testid="template-load-error" className="text-xs text-amber-300">
        {getSafeErrorStateMessage(loadError, t('template.load_failed', { ns: 'studio-review' }))}
        <button type="button" data-testid="template-load-retry" onClick={() => void refetch()}>{t('template.reload', { ns: 'studio-review' })}</button>
      </div>}
      {error && <p role="alert" data-testid="template-error" className="text-xs text-amber-300">{error}</p>}
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
            request.current++;
            setPreviewData(null);
            setError(null);
          }}
          disabled={!deviceId || isLoading || applyMutation.isPending}
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
          disabled={!selectedTemplateId || !deviceId || previewMutation.isPending || applyMutation.isPending || isError || isLoading}
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
      </div>
    </details>
  );
};

export const MeasurementTemplateSelector: React.FC<MeasurementTemplateSelectorProps> = (props) => {
  return (
    <SafeQueryBoundary>
      <MeasurementTemplateSelectorContent key={props.deviceId ?? 'none'} {...props} />
    </SafeQueryBoundary>
  );
};
