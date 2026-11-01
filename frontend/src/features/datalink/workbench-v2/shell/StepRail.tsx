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
  { id: 1, key: 'device', title: '連接設備', subtitle: '設備與通訊協定', icon: 'device' },
  { id: 2, key: 'rule', title: '選要讀的資料', subtitle: '設定一組或多組讀取範圍', icon: 'rule' },
  { id: 3, key: 'mapping', title: '命名與換算', subtitle: '名稱、單位與換算結果', icon: 'map' },
  { id: 4, key: 'database', title: '儲存並開始記錄', subtitle: '選目的地與確認開始', icon: 'db' },
];

export interface StepRailProps {
  view: 'flow' | 'settings';
  current: number;
  completed: Set<number>;
  step2Ready: boolean;
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
  step2Ready,
  collapsed,
  onJump,
  onSwitchView,
}) => {
  const focusClass =
    'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--gw-focus)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--gw-surface)]';

  return (
    <nav className="flex flex-col gap-1 w-full">
      {!collapsed && (
        <div className="px-1 mb-1 text-xs font-semibold text-[var(--gw-muted)]">
          建置流程
        </div>
      )}
      
      {STEPS.map((s) => {
        const isCurrent = view === 'flow' && current === s.id;
        const isDone = completed.has(s.id);
        const isStep2Blocked = s.id >= 3 && !step2Ready;
        const isReachable = !isStep2Blocked && (isDone || isCurrent || completed.has(s.id - 1) || s.id === 1);

        return (
          <button
            key={s.id}
            onClick={() => isReachable && (onSwitchView('flow'), onJump(s.id))}
            disabled={!isReachable}
            title={collapsed ? s.title : undefined}
            data-testid={`step-nav-button-${s.id}`}
            className={`relative flex items-center gap-3 rounded-xl border text-left transition-all duration-200 w-full ${focusClass} ${
              collapsed ? 'p-1.5 justify-center' : 'px-3 py-2.5'
            } ${
              isCurrent
                ? 'border-[var(--gw-focus)] bg-[var(--gw-active)]'
                : isDone
                ? 'border-[var(--gw-success-border)] bg-[var(--gw-success-bg)] hover:bg-[var(--gw-success-bg)]'
                : !isReachable
                ? 'border-[var(--gw-control)] bg-[var(--gw-surface)] disabled:text-[var(--gw-muted)] cursor-not-allowed'
                : 'border-[var(--gw-control)] bg-[var(--gw-surface)] hover:bg-[var(--gw-surface)]'
            }`}
          >
            {/* 圖示外框 */}
            <span
              className={`grid place-items-center w-9 h-9 rounded-lg border flex-shrink-0 transition-colors ${
                isCurrent
                  ? 'border-[var(--gw-focus)] bg-[var(--gw-active)] text-[var(--gw-accent)]'
                  : isDone
                  ? 'border-[var(--gw-success-border)] bg-[var(--gw-success-bg)] text-[var(--gw-success-fg)]'
                  : 'border-[var(--gw-border)] bg-[var(--gw-surface)] text-[var(--gw-muted)]'
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
                    <span className="text-xs font-mono text-[var(--gw-muted)]">0{s.id}</span>
                    <span
                      className={`text-sm font-semibold truncate ${
                        isCurrent ? 'text-[var(--gw-accent)]' : 'text-[var(--gw-text)]'
                      }`}
                    >
                      {s.title}
                    </span>
                  </div>
                  <div className="text-xs text-[var(--gw-muted)] truncate">{s.subtitle}</div>
                </div>
                {isCurrent && <span className="w-1.5 h-1.5 rounded-full bg-[var(--gw-active)]" />}
              </>
            )}
          </button>
        );
      })}

      <div className={`my-2 border-t border-[var(--gw-border)] ${collapsed ? 'mx-1' : 'mx-1'}`} />

      {!collapsed && (
        <div className="px-1 mb-1 text-xs font-semibold text-[var(--gw-muted)]">
          系統
        </div>
      )}

      {/* 設定按鈕 */}
      <button
        onClick={() => onSwitchView('settings')}
        title={collapsed ? '設定' : undefined}
        data-testid="settings-nav-button"
        className={`relative flex items-center gap-3 rounded-xl border text-left transition-all duration-200 w-full ${focusClass} ${
          collapsed ? 'p-1.5 justify-center' : 'px-3 py-2.5'
        } ${
          view === 'settings'
            ? 'border-[var(--gw-focus)] bg-[var(--gw-active)]'
            : 'border-[var(--gw-control)] bg-[var(--gw-surface)] hover:bg-[var(--gw-surface)]'
        }`}
      >
        <span
          className={`grid place-items-center w-9 h-9 rounded-lg border flex-shrink-0 ${
            view === 'settings'
              ? 'border-[var(--gw-focus)] bg-[var(--gw-active)] text-[var(--gw-accent)]'
              : 'border-[var(--gw-border)] bg-[var(--gw-surface)] text-[var(--gw-muted)]'
          }`}
        >
          <Icon name="sliders" className="w-4 h-4" />
        </span>
        {!collapsed && (
          <div className="min-w-0 flex-1">
            <div className="text-sm font-semibold truncate text-[var(--gw-text)]">設定</div>
            <div className="text-xs text-[var(--gw-muted)] truncate">資料庫、排程、介面…</div>
          </div>
        )}
        {view === 'settings' && !collapsed && (
          <span className="w-1.5 h-1.5 rounded-full bg-[var(--gw-active)]" />
        )}
      </button>
    </nav>
  );
};
