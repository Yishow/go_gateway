import * as React from 'react';
import { SectionCard } from '../components';

/**
 * Step 1: 新增裝置 佔位元件 [Rollback Only]
 * 
 * @deprecated Replaced by `steps/step1/Step1Device.tsx`；保留以供 rollback
 * 落地設計決策：「Placeholder step and settings surfaces」
 * 不發起任何網路請求。
 */
export const Step1DevicePlaceholder: React.FC = () => {
  return (
    <div className="sweep-in space-y-4">
      <SectionCard title="Step 1 · 新增裝置" subtitle="Device & Protocol Settings">
        <div className="border border-blue-500/20 bg-blue-500/5 rounded-xl p-5 text-center">
          <p className="text-sm text-slate-300 font-medium">
            Step 1 內容即將上線
          </p>
          <p className="text-xs text-slate-500 mt-1">
            完整內容由後續 change 交付。
          </p>
        </div>
      </SectionCard>
    </div>
  );
};
export default Step1DevicePlaceholder;
