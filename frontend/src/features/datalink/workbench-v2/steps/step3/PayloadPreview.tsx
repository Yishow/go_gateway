import * as React from 'react';
import type { Point, Mapping } from '../../state/types';
import { buildPayload } from '../../state/transformPipeline';

export interface PayloadPreviewProps {
  point: Point;
  mapping: Mapping;
}

/**
 * API Payload 預覽元件
 * 
 * 落地設計決策：「轉換管線數值流：使用 raw seed 序列」
 * 呈現前端將傳送給 API 儲存的 Pipeline 完整 JSON 結構，不發送任何 API 請求。
 */
export const PayloadPreview: React.FC<PayloadPreviewProps> = ({ point, mapping }) => {
  const payload = buildPayload(mapping, point);
  const jsonStr = JSON.stringify(payload, null, 2);

  return (
    <div className="space-y-1.5" data-testid="payload-preview">
      <div className="text-xs font-semibold text-slate-500 font-mono uppercase tracking-wider">
        API payload preview
      </div>
      <pre className="text-[11px] font-mono bg-slate-950/80 p-3.5 rounded-lg border border-slate-900 text-slate-400 overflow-x-auto max-h-[220px] scrollbar-thin">
        <code>{jsonStr}</code>
      </pre>
    </div>
  );
};
