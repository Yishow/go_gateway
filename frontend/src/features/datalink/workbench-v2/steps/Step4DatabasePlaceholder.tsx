import * as React from 'react';
import { SectionCard } from '../components';

/**
 * Step 4: 儲存資料庫 佔位元件
 * 
 * 落地設計決策：「Placeholder step and settings surfaces」
 * [Rollback Only] 此元件目前已無直接調用，僅保留作為緊急回滾備用。
 */
export const Step4DatabasePlaceholder: React.FC = () => {
  return (
    <div className="sweep-in space-y-4">
      <SectionCard title="Step 4 · 儲存資料庫" subtitle="Database target settings">
        <div className="border border-blue-500/20 bg-blue-500/5 rounded-xl p-5 text-center">
          <p className="text-sm text-slate-300 font-medium">
            Step 4 內容即將上線
          </p>
          <p className="text-xs text-slate-500 mt-1">
            完整內容由後續 change 交付。
          </p>
        </div>
      </SectionCard>
    </div>
  );
};
export default Step4DatabasePlaceholder;
