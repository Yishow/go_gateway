import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { SectionCard } from '../../components/SectionCard';
import { PipelineSteps } from './PipelineSteps';
import { PayloadPreview } from './PayloadPreview';
import type { Point, Mapping } from '../../state/types';

export interface TransformPreviewProps {
  point: Point | null;
  mapping: Mapping | null;
  rawSeed: number | null;
}

/**
 * 點位轉換預覽主區塊元件
 * 
 * 落地設計決策：「拆檔策略：6 個元件 + 2 個 state module」
 * 組合 SectionCard、PipelineSteps 與 PayloadPreview。當無選取時顯示 Empty State。
 */
export const TransformPreview: React.FC<TransformPreviewProps> = ({ point, mapping, rawSeed }) => {
  const { t } = useTranslation('workbench-v2');

  if (!point || !mapping || rawSeed === null) {
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
        <PipelineSteps point={point} mapping={mapping} rawSeed={rawSeed} />
        <div className="border-t border-slate-900 pt-6">
          <PayloadPreview point={point} mapping={mapping} />
        </div>
      </div>
    </SectionCard>
  );
};
