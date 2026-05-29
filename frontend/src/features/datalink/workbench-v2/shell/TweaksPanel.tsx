import * as React from 'react';
import { Toggle, Button, Icon } from '../components';

export interface TweaksPanelProps {
  sidebarCollapsed: boolean;
  onToggleSidebar: () => void;
  showSummaryRail: boolean;
  onToggleSummaryRail: () => void;
  onResetFlow: () => void;
}

/**
 * Workbench V2 Tweaks Panel (調試控制面板)
 * 
 * 落地設計決策：「Tweaks Panel：簡化原型版本」
 * 僅在 dev mode 或 localStorage.WBV2_TWEAKS === '1' 時載入渲染。
 */
export const TweaksPanel: React.FC<TweaksPanelProps> = ({
  sidebarCollapsed,
  onToggleSidebar,
  showSummaryRail,
  onToggleSummaryRail,
  onResetFlow,
}) => {
  const [isOpen, setIsOpen] = React.useState(false);
  const [shouldShow, setShouldShow] = React.useState(false);

  React.useEffect(() => {
    let show = false;

    // 1. 檢查是否為開發模式
    try {
      if (import.meta.env.DEV) {
        show = true;
      }
    } catch {
      // 忽略 env 錯誤
    }

    // 2. 檢查 localStorage 設定
    try {
      if (typeof window !== 'undefined' && window.localStorage) {
        if (window.localStorage.getItem('WBV2_TWEAKS') === '1') {
          show = true;
        }
      }
    } catch (e) {
      // eslint-disable-next-line no-console
      console.warn('localStorage is not available:', e);
    }

    setShouldShow(show);
  }, []);

  // 狀態持久化輔助函數
  const safeSetLocalStorage = (key: string, value: string) => {
    try {
      if (typeof window !== 'undefined' && window.localStorage) {
        window.localStorage.setItem(key, value);
      }
    } catch (e) {
      // eslint-disable-next-line no-console
      console.warn('localStorage is not available:', e);
    }
  };

  const handleToggleSidebar = () => {
    onToggleSidebar();
    safeSetLocalStorage('wbv2_sidebar_collapsed', String(!sidebarCollapsed));
  };

  const handleToggleSummary = () => {
    onToggleSummaryRail();
    safeSetLocalStorage('wbv2_show_summary_rail', String(!showSummaryRail));
  };

  const handleAddSecondDevice = () => {
    // TODO: 後續變更中實作 state 變更
    // eslint-disable-next-line no-console
    console.warn('pending follow-up change');
  };

  const handleAddSecondRule = () => {
    // TODO: 後續變更中實作 state 變更
    // eslint-disable-next-line no-console
    console.warn('pending follow-up change');
  };

  const handleJumpToStep4 = () => {
    // TODO: 後續變更中實作 state 變更
    // eslint-disable-next-line no-console
    console.warn('pending follow-up change');
  };

  if (!shouldShow) return null;

  return (
    <div className="fixed bottom-4 right-4 z-50">
      {/* 浮動齒輪按鈕 */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="w-10 h-10 rounded-full bg-slate-800 border border-slate-700 hover:bg-slate-700 text-slate-300 hover:text-white flex items-center justify-center shadow-lg transition-all focus:outline-none"
        title="開啟除錯面板"
        data-testid="tweaks-trigger"
      >
        <Icon name="sliders" className={`w-5 h-5 ${isOpen ? 'rotate-90' : ''} transition-transform duration-200`} />
      </button>

      {/* 面板主體 */}
      {isOpen && (
        <div
          className="absolute bottom-12 right-0 w-72 bg-slate-900 border border-slate-700/80 rounded-xl p-4 shadow-2xl sweep-in text-slate-200"
          data-testid="tweaks-panel-container"
        >
          <div className="flex items-center justify-between border-b border-slate-800 pb-2 mb-3">
            <span className="text-xs font-bold uppercase tracking-wider text-slate-400">
              🛠️ Workbench Tweaks
            </span>
            <button
              onClick={() => setIsOpen(false)}
              className="text-slate-500 hover:text-slate-300 text-xs"
            >
              關閉
            </button>
          </div>

          <div className="space-y-4">
            {/* Toggles */}
            <div className="space-y-2.5">
              <div className="flex items-center justify-between">
                <span className="text-xs font-medium">收合側邊欄</span>
                <Toggle checked={sidebarCollapsed} onChange={handleToggleSidebar} size="sm" />
              </div>
              <div className="flex items-center justify-between">
                <span className="text-xs font-medium">顯示右側摘要欄</span>
                <Toggle checked={showSummaryRail} onChange={handleToggleSummary} size="sm" />
              </div>
            </div>

            <div className="border-t border-slate-800/80 my-3" />

            {/* 示範按鈕 */}
            <div className="space-y-2">
              <span className="text-[10px] font-semibold text-slate-500 uppercase tracking-wider block">
                示範與除錯
              </span>
              <Button
                variant="danger"
                size="sm"
                className="w-full text-left justify-start"
                onClick={onResetFlow}
              >
                重置流程
              </Button>
              <Button
                variant="secondary"
                size="sm"
                className="w-full text-left justify-start text-slate-300"
                onClick={handleAddSecondDevice}
              >
                加入第二台設備
              </Button>
              <Button
                variant="secondary"
                size="sm"
                className="w-full text-left justify-start text-slate-300"
                onClick={handleAddSecondRule}
              >
                加入第二條規則
              </Button>
              <Button
                variant="secondary"
                size="sm"
                className="w-full text-left justify-start text-slate-300"
                onClick={handleJumpToStep4}
              >
                跳到 Step 4
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
