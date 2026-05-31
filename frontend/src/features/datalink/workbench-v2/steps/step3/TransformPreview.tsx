import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { SectionCard } from '../../components/SectionCard';
import { PipelineSteps } from './PipelineSteps';
import { PayloadPreview } from './PayloadPreview';
import { TargetTypeQuickActions } from './TargetTypeQuickActions';
import { useStep3LivePreview } from './useStep3LivePreview';
import type { Point, Mapping, TargetType } from '../../state/types';
import type { RuntimeStreamConnectionState } from '../../../../../types/datalink';

export interface TransformPreviewProps {
  point: Point | null;
  mapping: Mapping | null;
  rawValue: unknown | null;
  connectionState?: RuntimeStreamConnectionState;
  onApplyTargetTypeToAll?: () => void;
  onSelectTargetType?: (value: TargetType) => void;
}

/**
 * 點位轉換預覽主區塊元件
 * 
 * 落地設計決策：「拆檔策略：6 個元件 + 2 個 state module」
 * 組合 SectionCard、PipelineSteps 與 PayloadPreview。當無選取時顯示 Empty State。
 */
export const TransformPreview: React.FC<TransformPreviewProps> = ({
  point,
  mapping,
  rawValue,
  connectionState = 'disconnected',
  onApplyTargetTypeToAll,
  onSelectTargetType,
}) => {
  const { t } = useTranslation('workbench-v2');
  const preview = useStep3LivePreview(point, mapping, rawValue);

  if (!point || !mapping) {
    return (
      <SectionCard
        title={t('step3.preview.title', { defaultValue: 'Transform Preview' })}
        subtitle={t('step3.preview.emptySubtitle', { defaultValue: 'Select a row to preview' })}
      >
        <div className="flex flex-col items-center justify-center py-16 px-4 text-center border border-dashed border-slate-800 rounded-lg bg-slate-950/20" data-testid="preview-empty">
          <div className="text-2xl text-slate-700 mb-2 font-mono">⚡</div>
          <p className="text-xs text-slate-500 max-w-[200px] leading-relaxed">
            {t('step3.preview.emptyText', { defaultValue: 'Select any row on the left to preview its transform pipeline and API payload structure.' })}
          </p>
        </div>
      </SectionCard>
    );
  }

  return (
    <SectionCard
      title={t('step3.preview.title', { defaultValue: 'Transform Preview' })}
      subtitle={`${point.name} @ ${point.address}`}
      className="sticky top-6"
    >
      <div className="space-y-8" data-testid="preview-active">
        {mapping && onApplyTargetTypeToAll && onSelectTargetType ? (
          <div className="space-y-2">
            <div className="text-[11px] font-semibold uppercase tracking-wider text-slate-500">
              {t('step3.preview.quickTargetType', { defaultValue: '快速設定目標型態' })}
            </div>
            <TargetTypeQuickActions
              currentType={mapping.target_type}
              onApplyAll={onApplyTargetTypeToAll}
              onSelect={onSelectTargetType}
            />
          </div>
        ) : null}
        {rawValue === null || rawValue === undefined ? (
          <div
            className="rounded-lg border border-dashed border-slate-800 bg-slate-950/20 px-4 py-6 text-xs text-slate-400"
            data-testid="preview-waiting"
          >
            {connectionState === 'connecting'
              ? t('step3.preview.waitingConnecting', {
                  defaultValue: '正在連線裝置串流，等待第一筆即時值…',
                })
              : connectionState === 'error'
                ? t('step3.preview.waitingError', {
                    defaultValue: '尚未取得裝置即時值，請確認設備已啟動且串流正常。',
                  })
                : t('step3.preview.waiting', {
                    defaultValue: '等待裝置傳來即時值後再計算轉換預覽。',
                  })}
          </div>
        ) : preview.status === 'loading' ? (
          <div
            className="rounded-lg border border-dashed border-slate-800 bg-slate-950/20 px-4 py-6 text-xs text-slate-400"
            data-testid="preview-loading"
          >
            {t('step3.preview.loading', { defaultValue: '正在向後端計算轉換預覽…' })}
          </div>
        ) : preview.status === 'error' ? (
          <div
            className="rounded-lg border border-rose-900/60 bg-rose-950/20 px-4 py-6 text-xs text-rose-300"
            data-testid="preview-error"
          >
            {preview.error ?? t('step3.preview.errorFallback', { defaultValue: '轉換預覽失敗' })}
          </div>
        ) : (
          <PipelineSteps
            point={point}
            mapping={mapping}
            rawValue={rawValue}
            preview={preview.status === 'success' ? preview.data : null}
          />
        )}
        <div className="border-t border-slate-900 pt-6">
          <PayloadPreview point={point} mapping={mapping} />
        </div>
      </div>
    </SectionCard>
  );
};
