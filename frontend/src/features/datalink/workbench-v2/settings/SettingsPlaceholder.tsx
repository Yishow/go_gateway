import * as React from 'react';
import { SectionCard } from '../components';

/**
 * 設定頁 佔位元件 [Rollback Only]
 * 
 * 落地設計決策：「Placeholder step and settings surfaces」
 */
export const SettingsPlaceholder: React.FC = () => {
  return (
    <div className="sweep-in space-y-4">
      <SectionCard title="設定" subtitle="System global settings">
        <div className="border border-slate-700 bg-slate-800/40 rounded-xl p-5 text-center">
          <p className="text-sm text-slate-300 font-medium">
            設定內容即將上線
          </p>
          <p className="text-xs text-slate-500 mt-1">
            完整內容由後續 change 交付。
          </p>
        </div>
      </SectionCard>
    </div>
  );
};
export default SettingsPlaceholder;
