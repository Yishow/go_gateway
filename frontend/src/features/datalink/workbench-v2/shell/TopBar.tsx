import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { Button, Icon } from '../components';

export interface TopBarProps {
  sidebarCollapsed: boolean;
  onToggleSidebar: () => void;
  view: 'flow' | 'settings';
  stepTitle?: string;
  scheduler: 'idle' | 'running';
}

/**
 * Workbench V2 頂列元件
 * 
 * 落地設計決策：「Shell 五大區塊（Shell layout regions）」
 * 整合側邊欄開關、麵包屑、Scheduler 狀態與動作按鈕。
 */
export const TopBar: React.FC<TopBarProps> = ({
  sidebarCollapsed,
  onToggleSidebar,
  view,
  stepTitle = '新增裝置',
  scheduler,
}) => {
  const { t } = useTranslation('workbench-v2');
  const isRunning = scheduler === 'running';

  return (
    <header className="wbv2-topbar border-b border-slate-800/80 bg-slate-950/60 backdrop-blur-md sticky top-0 z-30">
      <div className="px-4 sm:px-5 lg:px-7 py-3 flex items-center gap-3 lg:gap-4">
        {/* 收合按鈕 */}
        <button
          onClick={onToggleSidebar}
          className="grid place-items-center w-8 h-8 rounded-lg border border-slate-700/60 text-slate-400 hover:text-slate-100 hover:bg-slate-800/60 active:scale-[0.98] transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-300/60 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-950"
          title={sidebarCollapsed ? '展開側邊欄 (⌘B)' : '收合側邊欄 (⌘B)'}
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75" className="w-4 h-4">
            <rect x="3" y="4" width="18" height="16" rx="2" />
            <line x1="9" y1="4" x2="9" y2="20" />
            {!sidebarCollapsed && <line x1="6" y1="8" x2="6" y2="8.01" strokeWidth="2" />}
          </svg>
        </button>

        {/* 品牌標識 */}
        <div className="flex items-center gap-2.5">
          <div className="wbv2-brand-mark grid place-items-center w-9 h-9 rounded-lg border border-cyan-300/30">
            <Icon name="bolt" className="w-4 h-4 text-cyan-200" />
          </div>
          <div>
            <div className="text-sm font-semibold tracking-tight text-slate-100 leading-tight">Datalink Workbench</div>
            <div className="text-[10px] text-slate-500 leading-tight">go_gateway · v1 接入點位設定流程</div>
          </div>
        </div>

        {/* Breadcrumb */}
        <nav className="hidden md:flex items-center gap-1.5 ml-6 text-xs">
          <span className="text-slate-500">Datalink</span>
          <Icon name="chevron" className="w-3 h-3 text-slate-600" />
          <span className="text-slate-500">Workbench</span>
          <Icon name="chevron" className="w-3 h-3 text-slate-600" />
          <span className="text-slate-200 font-medium">
            {view === 'settings' ? '設定' : stepTitle}
          </span>
        </nav>

        {/* 右側按鈕與排程狀態 */}
        <div className="ml-auto flex items-center gap-2">
          <Button variant="ghost" size="sm" icon={<Icon name="save" className="w-3.5 h-3.5" />} className="hidden sm:inline-flex">
            儲存草稿
          </Button>
          <Button variant="secondary" size="sm" className="hidden md:inline-flex">
            取消
          </Button>
          <div className="hidden sm:block h-6 w-px bg-slate-700/60 mx-1" />
          <div className="wbv2-status-pill flex items-center gap-1.5 text-xs text-slate-400">
            <span
              className={`w-1.5 h-1.5 rounded-full ${
                isRunning 
                  ? 'bg-emerald-400 shadow-[0_0_8px_rgba(16,185,129,0.5)] animate-pulse' 
                  : 'bg-amber-400 shadow-[0_0_8px_rgba(245,158,11,0.5)]'
              }`}
            />
            <span>
              {isRunning 
                ? t('step4.scheduler_running', 'scheduler running') 
                : t('step4.scheduler_idle', 'scheduler idle')
              }
            </span>
          </div>
        </div>
      </div>
    </header>
  );
};
