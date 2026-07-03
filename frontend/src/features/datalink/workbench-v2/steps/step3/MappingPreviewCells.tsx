import * as React from 'react';
import type { Point, Mapping } from '../../state/types';
import { castValue, formatFinal, runScale } from '../../state/transformPipeline';
import { useStep3LivePreview } from './useStep3LivePreview';

export interface MappingPreviewCellsProps {
  point: Point;
  mapping: Mapping;
  rawValue?: unknown;
}

function isPreviewScalar(value: unknown): value is number | boolean | string {
  return typeof value === 'number' || typeof value === 'boolean' || typeof value === 'string';
}

function toNumericValue(value: unknown): number {
  return typeof value === 'number' ? value : Number(value) || 0;
}

function formatPreviewValue(value: unknown, fallback = '--'): string {
  if (value === null || value === undefined) {
    return fallback;
  }
  if (typeof value === 'boolean') {
    return value ? 'true' : 'false';
  }
  if (typeof value === 'number') {
    return value.toFixed(2);
  }
  return String(value);
}

export const MappingPreviewCells: React.FC<MappingPreviewCellsProps> = ({
  point,
  mapping,
  rawValue,
}) => {
  const preview = useStep3LivePreview(
    point,
    mapping,
    rawValue === undefined ? null : rawValue,
  );

  if (rawValue === null || rawValue === undefined) {
    return (
      <>
        <td className="p-2 text-xs font-mono text-slate-500" data-testid={`preview-scale-${point.id}`}>--</td>
        <td className="p-2 text-xs font-mono text-slate-500" data-testid={`preview-cast-${point.id}`}>--</td>
        <td className="p-2 text-xs font-mono text-slate-500" data-testid={`preview-final-${point.id}`}>--</td>
      </>
    );
  }

  const numericRawValue = toNumericValue(rawValue);
  const previewScaleValue = preview.data?.step_results.find((step) => step.step_type === 'scale')?.output_value;
  const previewCastValue = preview.data?.step_results.find((step) => step.step_type === 'cast')?.output_value;
  const scaled = typeof previewScaleValue === 'number'
    ? previewScaleValue
    : runScale(numericRawValue, mapping.scale, mapping.offset);
  const casted = isPreviewScalar(previewCastValue)
    ? previewCastValue
    : castValue(scaled, mapping.target_type);
  const finalSource = isPreviewScalar(preview.data?.final_value)
    ? preview.data?.final_value
    : casted;

  return (
    <>
      <td className="p-2 text-xs font-mono text-blue-300" data-testid={`preview-scale-${point.id}`}>
        {formatPreviewValue(scaled)}
      </td>
      <td className="p-2 text-xs font-mono text-emerald-300" data-testid={`preview-cast-${point.id}`}>
        {formatFinal(casted, mapping.target_type)}
      </td>
      <td className="p-2 text-xs font-mono text-emerald-200" data-testid={`preview-final-${point.id}`}>
        {formatFinal(finalSource, mapping.target_type)}
      </td>
    </>
  );
};
