import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { Point, Mapping } from '../../state/types';
import type { MappingPreviewResponse } from '../../../../../types/datalink';
import { runScale, castValue, formatFinal } from '../../state/transformPipeline';

function isPreviewScalar(value: unknown): value is number | boolean | string {
  return typeof value === 'number' || typeof value === 'boolean' || typeof value === 'string';
}

export interface PipelineStepsProps {
  point: Point;
  mapping: Mapping;
  rawValue: unknown;
  preview?: MappingPreviewResponse | null;
}

function toNumericValue(value: unknown): number {
  return typeof value === 'number' ? value : Number(value) || 0;
}

/**
 * 轉換管線預覽步驟元件
 * 
 * 落地設計決策：「轉換管線數值流：使用 raw seed 序列」
 * 渲染 4 步驟垂直管線，展示資料從 PLC 解碼到寫入 Tag 的完整運算軌跡與顏色語意。
 */
export const PipelineSteps: React.FC<PipelineStepsProps> = ({
  point,
  mapping,
  rawValue,
  preview = null,
}) => {
  const { t } = useTranslation('workbench-v2');
  const numericRawValue = toNumericValue(rawValue);

  const previewScaleValue = preview?.step_results.find((step) => step.step_type === 'scale')?.output_value;
  const previewCastValue = preview?.step_results.find((step) => step.step_type === 'cast')?.output_value;
  const scaled = typeof previewScaleValue === 'number'
    ? previewScaleValue
    : runScale(numericRawValue, mapping.scale, mapping.offset);
  const casted = isPreviewScalar(previewCastValue)
    ? previewCastValue
    : castValue(scaled, mapping.target_type);
  const finalSource = isPreviewScalar(preview?.final_value)
    ? preview.final_value
    : casted;
  const finalVal = formatFinal(finalSource, mapping.target_type);
  const decodeValue = preview?.raw_value ?? rawValue;

  return (
    <div className="relative pl-6 border-l border-slate-800 space-y-8 ml-3 py-1">
      {/* 步驟 1: 解碼 */}
      <div className="relative group" data-testid="step-decode">
        <span className="absolute -left-[31px] top-1.5 flex h-4 w-4 items-center justify-center rounded-full border border-slate-750 bg-slate-950 text-[10px] text-slate-500 font-mono">
          1
        </span>
        <div className="text-xs font-medium text-slate-400 font-mono">
          {t('step3.steps.decode', { defaultValue: 'decode (raw)' })}
        </div>
        <div className="mt-1 text-sm font-semibold text-slate-300 font-mono">
          {String(decodeValue)} <span className="text-xs text-slate-500 font-normal">({point.data_type})</span>
        </div>
      </div>

      {/* 步驟 2: 線性縮放 */}
      <div className="relative group" data-testid="step-scale">
        <span className="absolute -left-[31px] top-1.5 flex h-4 w-4 items-center justify-center rounded-full border border-blue-900 bg-slate-950 text-[10px] text-blue-400 font-mono">
          2
        </span>
        <div className="text-xs font-medium text-blue-400/80 font-mono">
          {t('step3.steps.scale', { defaultValue: 'scale (linear)' })}
        </div>
        <div className="mt-1 text-sm font-semibold text-blue-300 font-mono">
          {String(decodeValue)} × {mapping.scale} {mapping.offset >= 0 ? `+ ${mapping.offset}` : `- ${Math.abs(mapping.offset)}`} = {scaled.toFixed(2)}
        </div>
      </div>

      {/* 步驟 3: 型態轉換 */}
      <div className="relative group" data-testid="step-cast">
        <span className="absolute -left-[31px] top-1.5 flex h-4 w-4 items-center justify-center rounded-full border border-emerald-900 bg-slate-950 text-[10px] text-emerald-400 font-mono">
          3
        </span>
        <div className="text-xs font-medium text-emerald-400/80 font-mono">
          {t('step3.steps.cast', { defaultValue: 'cast to {{type}}', type: mapping.target_type })}
        </div>
        <div className="mt-1 text-sm font-semibold text-emerald-300 font-mono">
          {formatFinal(casted, mapping.target_type)}
        </div>
      </div>

      {/* 步驟 4: 輸出標籤 */}
      <div className="relative group" data-testid="step-final">
        <span className="absolute -left-[31px] top-1 flex h-4 w-4 items-center justify-center rounded-full border border-emerald-800 bg-slate-950 text-[10px] text-emerald-300 font-mono">
          4
        </span>
        <div className="text-xs font-medium text-emerald-300/80 font-mono">
          {t('step3.steps.final', { defaultValue: 'final → Tag' })}
        </div>
        <div className="mt-1 flex flex-col gap-0.5">
          <div className="text-2xl font-bold text-emerald-200 tracking-tight font-mono">
            {finalVal}
          </div>
          <div className="text-xs text-slate-500 font-mono truncate max-w-[200px]" title={`tag.${mapping.tag_key}`}>
            → tag.{mapping.tag_key || '(未設定)'}
          </div>
        </div>
      </div>
    </div>
  );
};
