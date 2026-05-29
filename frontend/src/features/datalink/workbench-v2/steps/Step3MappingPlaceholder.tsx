import * as React from 'react';
import { SectionCard } from '../components';

/**
 * Step 3: 點位映射 佔位元件
 * 
 * ROLLBACK ONLY: This file is preserved for easy rollback support only. Do not import this file in production.
 * 落地設計決策：「Placeholder step and settings surfaces」
 */
export const Step3MappingPlaceholder: React.FC = () => {
  return (
    <div className="sweep-in space-y-4">
      <SectionCard title="Step 3 · 點位映射" subtitle="Point mapping to Tag">
        <div className="border border-blue-500/20 bg-blue-500/5 rounded-xl p-5 text-center">
          <p className="text-sm text-slate-300 font-medium">
            Step 3 內容即將上線
          </p>
          <p className="text-xs text-slate-500 mt-1">
            完整內容由後續 change 交付。
          </p>
        </div>
      </SectionCard>
    </div>
  );
};
export default Step3MappingPlaceholder;
