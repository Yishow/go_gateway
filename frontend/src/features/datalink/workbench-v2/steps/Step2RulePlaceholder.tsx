import * as React from 'react';
import { SectionCard } from '../components';

/**
 * [ROLLBACK ONLY] Step 2: 接入規則 佔位元件
 * 
 * 此檔案僅保留供緊急 Rollback 使用，正式流程請使用 steps/step2/Step2Rule.tsx。
 * 落地設計決策：「Placeholder step and settings surfaces」
 */
export const Step2RulePlaceholder: React.FC = () => {
  return (
    <div className="sweep-in space-y-4">
      <SectionCard title="Step 2 · 接入規則" subtitle="Source Rules Settings">
        <div className="border border-blue-500/20 bg-blue-500/5 rounded-xl p-5 text-center">
          <p className="text-sm text-slate-300 font-medium">
            Step 2 內容即將上線
          </p>
          <p className="text-xs text-slate-500 mt-1">
            完整內容由後續 change 交付。
          </p>
        </div>
      </SectionCard>
    </div>
  );
};
export default Step2RulePlaceholder;
