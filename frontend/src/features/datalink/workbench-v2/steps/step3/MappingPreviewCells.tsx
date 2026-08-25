import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { Point, Mapping } from '../../state/types';
import type { RuntimeStreamConnectionState } from '../../../../../types/datalink';
import { formatFinal } from '../../state/transformPipeline';
import { useStep3LivePreview } from './useStep3LivePreview';

export interface MappingPreviewCellsProps {
  point: Point;
  mapping: Mapping;
  rawValue?: unknown;
  connectionState?: RuntimeStreamConnectionState;
  workspaceId?: string;
}

function isPreviewScalar(value: unknown): value is number | boolean | string {
  return typeof value === 'number' || typeof value === 'boolean' || typeof value === 'string';
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
  connectionState = 'connected',
  workspaceId,
}) => {
  const { t } = useTranslation('workbench-v2');
  const preview = useStep3LivePreview(
    point,
    mapping,
    rawValue === undefined ? null : rawValue,
    connectionState,
    workspaceId,
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

  const previewScaleValue = preview.data?.step_results.find((step) => step.step_type === 'scale')?.output_value;
  const previewCastValue = preview.data?.step_results.find((step) => step.step_type === 'cast')?.output_value;
  const hasServerPreview = preview.data !== null &&
    ['live', 'stale', 'reconnecting', 'degraded'].includes(preview.status);
  const scaled = hasServerPreview && isPreviewScalar(previewScaleValue) ? previewScaleValue : null;
  const casted = hasServerPreview && isPreviewScalar(previewCastValue) ? previewCastValue : null;
  const finalSource = hasServerPreview && isPreviewScalar(preview.data?.final_value)
    ? preview.data?.final_value
    : null;
  const stateLabel = t(`step3.previewStates.${preview.status}`);
  const safeError = preview.errorCode
    ? t(`errors.${preview.errorCode}`, { defaultValue: t('errors.preview_unavailable') })
    : t('errors.preview_unavailable');

  return (
    <>
      <td className="p-2 text-xs font-mono text-blue-300" data-testid={`preview-scale-${point.id}`} aria-label={stateLabel}>
        {formatPreviewValue(scaled)}
      </td>
      <td className="p-2 text-xs font-mono text-emerald-300" data-testid={`preview-cast-${point.id}`}>
        {casted === null ? '--' : formatFinal(casted, mapping.target_type)}
      </td>
      <td className="p-2 text-xs font-mono text-emerald-200" data-testid={`preview-final-${point.id}`}>
        {finalSource === null ? '--' : formatFinal(finalSource, mapping.target_type)}
        {preview.error && (
          <div className="mt-1 space-y-1 font-sans text-[10px] text-rose-300" data-testid={`preview-error-${point.id}`} role="alert">
            <div>{safeError}</div>
            {preview.requestId && (
              <div data-testid={`preview-request-id-${point.id}`}>
                {t('errors.request_id')}: {preview.requestId}
              </div>
            )}
            {preview.error && (
              <button
                type="button"
                className="text-rose-200 underline underline-offset-2 hover:text-white"
                onClick={preview.retry}
                data-testid={`preview-retry-${point.id}`}
              >
                {t('errors.retry')}
              </button>
            )}
          </div>
        )}
        <span className="sr-only" data-testid={`preview-state-${point.id}`} data-state={preview.status} role="status" aria-live="polite">
          {stateLabel}
        </span>
      </td>
    </>
  );
};
