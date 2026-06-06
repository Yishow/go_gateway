import * as React from 'react';
import { Icon, type IconName } from '../components';

export interface Step {
  id: 1 | 2 | 3 | 4;
  key: string;
  title: string;
  subtitle: string;
  icon: IconName;
}

export const STEPS: Step[] = [
  { id: 1, key: 'device', title: '新增裝置', subtitle: 'Device + Protocol', icon: 'device' },
  { id: 2, key: 'rule', title: '接入規則', subtitle: 'Source Rules · 多條範圍', icon: 'rule' },
  { id: 3, key: 'mapping', title: '點位映射', subtitle: 'Point → Tag', icon: 'map' },
  { id: 4, key: 'database', title: '儲存資料庫', subtitle: 'DB Target · Commit', icon: 'db' },
];

export interface StepRailProps {
  view: 'flow' | 'settings';
  current: number;
  completed: Set<number>;
  collapsed: boolean;
  onJump: (stepId: 1 | 2 | 3 | 4) => void;
  onSwitchView: (view: 'flow' | 'settings') => void;
}

/**
 * Workbench V2 側邊步驟導航列元件
 * 
 * 落地設計決策：「Step rail navigation」與「Shell layout regions」
 * 限制未解鎖步驟的存取，並支援收合狀態下的 tooltip 提示。
 */
export const StepRail: React.FC<StepRailProps> = ({
  view,
  current,
  completed,
  collapsed,
  onJump,
  onSwitchView,
}) => {
  const focusClass =
    'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-300/60 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-950';

  return (
    <nav className="flex flex-col gap-1 w-full">
      {!collapsed && (
        <div className="px-1 mb-1 text-[10px] font-semibold uppercase tracking-[0.18em] text-slate-500">
          建置流程
        </div>
      )}
      
      {STEPS.map((s) => {
        const isCurrent = view === 'flow' && current === s.id;
        const isDone = completed.has(s.id);
        const isReachable = isDone || isCurrent || completed.has(s.id - 1) || s.id === 1;

        return (
          <button
            key={s.id}
            onClick={() => isReachable && (onSwitchView('flow'), onJump(s.id))}
            disabled={!isReachable}
            title={collapsed ? s.title : undefined}
            data-testid={`step-nav-button-${s.id}`}
            className={`relative flex items-center gap-3 rounded-xl border text-left transition-all duration-200 active:scale-[0.99] w-full ${focusClass} ${
              collapsed ? 'p-1.5 justify-center' : 'px-3 py-2.5'
            } ${
              isCurrent
                ? 'border-blue-500/40 bg-blue-500/10 shadow-lg shadow-blue-500/5'
                : isDone
                ? 'border-emerald-500/20 bg-emerald-500/[0.04] hover:bg-emerald-500/[0.08]'
                : !isReachable
                ? 'border-slate-800 bg-slate-900/30 opacity-60 cursor-not-allowed'
                : 'border-slate-800 bg-slate-900/30 hover:bg-slate-800/50'
            }`}
          >
            {/* 圖示外框 */}
            <span
              className={`grid place-items-center w-9 h-9 rounded-lg border flex-shrink-0 transition-colors ${
                isCurrent
                  ? 'border-blue-500/50 bg-blue-500/20 text-blue-200'
                  : isDone
                  ? 'border-emerald-500/40 bg-emerald-500/15 text-emerald-300'
                  : 'border-slate-700 bg-slate-800 text-slate-500'
              }`}
            >
              {isDone && !isCurrent ? (
                <Icon name="check" className="w-4 h-4" />
              ) : (
                <Icon name={s.icon} className="w-4 h-4" />
              )}
            </span>

            {/* 文字說明（非收合狀態） */}
            {!collapsed && (
              <>
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-1.5">
                    <span className="text-[10px] font-mono text-slate-500">0{s.id}</span>
                    <span
                      className={`text-sm font-semibold truncate ${
                        isCurrent ? 'text-blue-100' : 'text-slate-200'
                      }`}
                    >
                      {s.title}
                    </span>
                  </div>
                  <div className="text-[11px] text-slate-500 truncate">{s.subtitle}</div>
                </div>
                {isCurrent && <span className="w-1.5 h-1.5 rounded-full bg-blue-400 pulse-dot" />}
              </>
            )}
          </button>
        );
      })}

      <div className={`my-2 border-t border-slate-800 ${collapsed ? 'mx-1' : 'mx-1'}`} />

      {!collapsed && (
        <div className="px-1 mb-1 text-[10px] font-semibold uppercase tracking-[0.18em] text-slate-500">
          系統
        </div>
      )}

      {/* 設定按鈕 */}
      <button
        onClick={() => onSwitchView('settings')}
        title={collapsed ? '設定' : undefined}
        data-testid="settings-nav-button"
        className={`relative flex items-center gap-3 rounded-xl border text-left transition-all duration-200 active:scale-[0.99] w-full ${focusClass} ${
          collapsed ? 'p-1.5 justify-center' : 'px-3 py-2.5'
        } ${
          view === 'settings'
            ? 'border-blue-500/40 bg-blue-500/10 shadow-lg shadow-blue-500/5'
            : 'border-slate-800 bg-slate-900/30 hover:bg-slate-800/50'
        }`}
      >
        <span
          className={`grid place-items-center w-9 h-9 rounded-lg border flex-shrink-0 ${
            view === 'settings'
              ? 'border-blue-500/50 bg-blue-500/20 text-blue-200'
              : 'border-slate-700 bg-slate-800 text-slate-500'
          }`}
        >
          <Icon name="sliders" className="w-4 h-4" />
        </span>
        {!collapsed && (
          <div className="min-w-0 flex-1">
            <div className="text-sm font-semibold truncate text-slate-200">設定</div>
            <div className="text-[11px] text-slate-500 truncate">資料庫、排程、介面…</div>
          </div>
        )}
        {view === 'settings' && !collapsed && (
          <span className="w-1.5 h-1.5 rounded-full bg-blue-400 pulse-dot" />
        )}
      </button>
    </nav>
  );
};
