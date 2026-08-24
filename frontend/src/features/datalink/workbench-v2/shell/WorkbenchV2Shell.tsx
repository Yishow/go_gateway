import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { TFunction } from 'i18next';
import type { WorkbenchV2State } from '../state/types';
import { useWorkbenchV2State } from '../state/useWorkbenchV2State';
import { TopBar } from './TopBar';
import { StepRail, STEPS } from './StepRail';
import { SummaryRail } from './SummaryRail';
import { TweaksPanel } from './TweaksPanel';
import { Icon, Button } from '../components';
import {
  buildRuntimeDashboardTarget,
  resolveRuntimeDashboardDevice,
} from './resolveRuntimeDashboardDevice';
import type { StudioV2ActivationResponse } from '../../../../types/studioV2Activation';
import type { StudioV2WorkspaceAuditEntry } from '../../../../types/studioV2WorkspaceAudit';
import type {
  StudioV2WorkspaceReadinessIssue,
  StudioV2WorkspaceReadinessSummary,
} from '../../../../types/studioV2WorkspaceReadiness';
import {
  readinessStepNumber,
  WorkspaceReadinessPanel,
  type WorkspaceReadinessStepNumber,
} from '../components/WorkspaceReadinessPanel';

// 導入真實的步驟元件與設定頁面
import { SettingsPage } from '../settings';

// 導入真實的步驟元件
import { Step1Device } from '../steps/step1';
import { Step2Rule } from '../steps/step2';
import { Step3Mapping } from '../steps/step3';
import { Step4Database } from '../steps/step4';
import { isStep2Ready } from '../state/sourceRule';

export interface WorkbenchV2ShellProps {
  state: WorkbenchV2State;
  actions: ReturnType<typeof useWorkbenchV2State>;
  navigateTo?: (target: string) => void;
  activateWorkspace?: () => Promise<StudioV2ActivationResponse>;
  workspaceReadiness?: StudioV2WorkspaceReadinessSummary | null;
  workspaceAuditHistory?: StudioV2WorkspaceAuditEntry[];
  workspaceAuditUnavailable?: boolean;
  runtimeReturnFocus?: WorkbenchV2RuntimeReturnFocus | null;
}

export interface WorkbenchV2RuntimeReturnFocus {
  step: WorkspaceReadinessStepNumber;
  issueCode?: string | null;
}

/**
 * Workbench V2 主 Shell 佈局元件
 * 
 * 落地設計決策：「Shell 五大區塊（Shell layout regions）」與「Keyboard shortcut」
 * 組合所有子佈局區塊，並管理鍵盤捷徑及步驟渲染切換。
 */
export const WorkbenchV2Shell: React.FC<WorkbenchV2ShellProps> = ({
  state,
  actions,
  navigateTo,
  activateWorkspace,
  workspaceReadiness,
  workspaceAuditHistory,
  workspaceAuditUnavailable,
  runtimeReturnFocus,
}) => {
  const { t } = useTranslation('workbench-v2');
  const {
    view,
    current,
    completed,
    sidebarCollapsed,
    showSummaryRail,
  } = state;
  const deviceProtocolMap = Object.fromEntries(state.devices.map((device) => [device.id, device.protocol]));
  const step2Ready = isStep2Ready(state.rules, deviceProtocolMap);

  const {
    setView,
    setCurrent,
    toggleSidebar,
    toggleSummaryRail,
    resetFlow,
    completeStep,
  } = actions;

  // 快捷鍵監聽 ⌘B / Ctrl+B
  React.useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      const isB = e.key === 'b' || e.key === 'B';
      const isMetaOrCtrl = e.metaKey || e.ctrlKey;

      if (isMetaOrCtrl && isB) {
        // 早退判斷：如果 focus 在 input, textarea, select 或 contenteditable 元素中，不觸發快捷鍵
        const target = e.target as HTMLElement;
        if (
          target &&
          (target.tagName === 'INPUT' ||
            target.tagName === 'TEXTAREA' ||
            target.tagName === 'SELECT' ||
            target.isContentEditable)
        ) {
          return;
        }

        e.preventDefault();
        toggleSidebar();

        // 寫入 localStorage 持久化
        try {
          if (typeof window !== 'undefined' && window.localStorage) {
            window.localStorage.setItem('wbv2_sidebar_collapsed', String(!sidebarCollapsed));
          }
        } catch {
          // 忽略
        }
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [sidebarCollapsed, toggleSidebar]);

  const stepMeta = STEPS.find((s) => s.id === current);
  const runtimeFocusStep = runtimeReturnFocus?.step;
  const runtimeFocusIssueCode = runtimeReturnFocus?.issueCode ?? null;
  const runtimeFocusedIssue = React.useMemo(() => {
    if (!runtimeFocusStep) {
      return null;
    }

    const issues = workspaceReadiness?.issues ?? [];
    const exactIssue = runtimeFocusIssueCode
      ? issues.find((issue) => issue.code === runtimeFocusIssueCode)
      : null;

    return exactIssue ?? issues.find((issue) => readinessStepNumber(issue.step) === runtimeFocusStep) ?? null;
  }, [runtimeFocusIssueCode, runtimeFocusStep, workspaceReadiness?.issues]);
  const runtimeFocusGuidance = runtimeFocusedIssue
    ? runtimeReturnGuidanceFor(runtimeFocusedIssue, t)
    : null;
  const runtimeFocusScope = runtimeFocusedIssue?.scope.trim();

  React.useEffect(() => {
    if (!runtimeFocusStep) {
      return;
    }

    setView('flow');
    setCurrent(runtimeFocusStep);
    document.getElementById('step-content')?.scrollTo?.({ top: 0, behavior: 'smooth' });
  }, [runtimeFocusStep, runtimeFocusIssueCode, setCurrent, setView]);

  const goBack = () => {
    if (current > 1) {
      setCurrent((current - 1) as 1 | 2 | 3 | 4);
      document.getElementById('step-content')?.scrollTo({ top: 0, behavior: 'smooth' });
    }
  };

  const handleReset = () => {
    resetFlow();
  };

  const handleContinueStep1 = () => {
    completeStep(1);
    setCurrent(2);
    document.getElementById('step-content')?.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handleContinueStep2 = () => {
    if (!step2Ready) {
      return;
    }
    completeStep(2);
    setCurrent(3);
    document.getElementById('step-content')?.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handleRuntimeDashboardHandoff = React.useCallback(() => {
    const target = buildRuntimeDashboardTarget(resolveRuntimeDashboardDevice(state));
    navigateTo?.(target);
  }, [navigateTo, state]);

  const handleNavigateStep = React.useCallback((step: WorkspaceReadinessStepNumber) => {
    setCurrent(step);
    document.getElementById('step-content')?.scrollTo({ top: 0, behavior: 'smooth' });
  }, [setCurrent]);

  // 渲染中央步驟內容
  const renderContent = () => {
    if (view === 'settings') {
      return (
        <SettingsPage
          state={state}
          dispatch={actions.dispatch}
        />
      );
    }

    switch (current) {
      case 1:
        return (
          <Step1Device
            state={state}
            dispatch={actions.dispatch}
            onContinue={handleContinueStep1}
          />
        );
      case 2:
        return (
          <Step2Rule
            state={state}
            dispatch={actions.dispatch}
            onContinue={handleContinueStep2}
          />
        );
      case 3:
        return (
          <Step3Mapping
            state={state}
            dispatch={actions.dispatch}
            onContinue={() => {
              completeStep(3);
              setCurrent(4);
              document.getElementById('step-content')?.scrollTo({ top: 0, behavior: 'smooth' });
            }}
            onBack={goBack}
          />
        );
      case 4:
        return (
          <Step4Database
            state={state}
            dispatch={actions.dispatch}
            onCommit={handleRuntimeDashboardHandoff}
            activateWorkspace={activateWorkspace}
            workspaceReadiness={workspaceReadiness}
            onNavigateStep={handleNavigateStep}
          />
        );
      default:
        return (
          <Step1Device
            state={state}
            dispatch={actions.dispatch}
            onContinue={handleContinueStep1}
          />
        );
    }
  };

  return (
    <div className="wbv2-shell min-h-screen flex flex-col" data-testid="workbench-v2-shell">
      {/* 1. 頂部列 */}
      <TopBar
        sidebarCollapsed={sidebarCollapsed}
        onToggleSidebar={toggleSidebar}
        view={view}
        stepTitle={stepMeta?.title}
        scheduler={state.devices.some((device) => device.running) ? 'running' : 'idle'}
      />

      {/* 主體區域 */}
      <div className="wbv2-workspace flex-1 flex gap-5 px-4 sm:px-5 lg:px-7 py-5">
        {/* 2. 左側步驟導航列 */}
        <aside
          className="min-w-0 overflow-hidden transition-all duration-200"
          style={{ flex: `0 0 ${sidebarCollapsed ? 64 : 232}px`, width: sidebarCollapsed ? 64 : 232 }}
          data-testid="sidebar-rail-aside"
          data-collapsed={sidebarCollapsed}
          aria-expanded={!sidebarCollapsed}
        >
          <div className="sticky top-[68px] space-y-4">
            <StepRail
              view={view}
              current={current}
              completed={completed}
              step2Ready={step2Ready}
              collapsed={sidebarCollapsed}
              onJump={setCurrent}
              onSwitchView={setView}
            />

            {!sidebarCollapsed && view === 'flow' && (
              <div className="wbv2-flow-note rounded-xl border border-slate-700/40 bg-slate-900/30 p-3 sweep-in">
                <div className="text-[10px] uppercase tracking-wider text-slate-500 mb-1.5">
                  本流程
                </div>
                <p className="text-[11px] leading-relaxed text-slate-400">
                  由「設備 → 接入規則 → 點位/映射 → 資料庫」四個階段組成，全部成功後才會啟動 collector 並寫入。
                </p>
              </div>
            )}
          </div>
        </aside>

        {/* 3. 中央主內容區 */}
        <main id="step-content" className="flex-1 min-w-0">
          {view === 'flow' ? (
            <>
              {/* 步驟 Progress Bar Header */}
              <div className="mb-4 flex items-center justify-between">
                <div>
                  <div className="flex items-center gap-2">
                    <span className="text-[11px] font-mono text-slate-500">
                      STEP 0{current} / 04
                    </span>
                    <span className="h-3 w-px bg-slate-700" />
                    <span className="text-[11px] text-slate-500">{stepMeta?.subtitle}</span>
                  </div>
                  <h2 className="text-xl font-semibold text-slate-100 mt-0.5">
                    {stepMeta?.title}
                  </h2>
                </div>

                <div className="hidden sm:flex items-center gap-1.5">
                  {STEPS.map((s) => (
                    <span
                      key={s.id}
                      className={`h-1 rounded-full transition-all ${
                        s.id < current
                          ? 'w-8 bg-emerald-500/80'
                          : s.id === current
                          ? 'w-12 bg-blue-500'
                          : 'w-8 bg-slate-700'
                      }`}
                    />
                  ))}
                </div>
              </div>

              <div className="mb-4">
                <WorkspaceReadinessPanel
                  summary={workspaceReadiness}
                  dataTestId="workbench-v2-readiness-banner"
                  compact
                  maxIssues={2}
                />
              </div>

              {runtimeFocusStep ? (
                <section
                  data-testid="workbench-v2-runtime-return-focus"
                  className="mb-4 rounded-xl border border-amber-400/40 bg-amber-400/10 px-4 py-3 text-sm text-amber-50"
                >
                  <div className="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
                    <div>
                      <p className="font-semibold">
                        {t('runtime_return.title', { step: `Step ${runtimeFocusStep}` })}
                      </p>
                      <p className="mt-1 text-xs leading-5 text-amber-100/80">
                        {runtimeFocusedIssue?.message ?? t('runtime_return.description')}
                      </p>
                      {runtimeFocusedIssue ? (
                        <p className="mt-1 font-mono text-[11px] text-amber-200/90">
                          {runtimeFocusedIssue.code}
                        </p>
                      ) : null}
                      {runtimeFocusGuidance ? (
                        <p
                          data-testid="workbench-v2-runtime-return-guidance"
                          className="mt-2 text-xs leading-5 text-amber-50"
                        >
                          {runtimeFocusGuidance}
                        </p>
                      ) : null}
                      {runtimeFocusScope ? (
                        <p
                          data-testid="workbench-v2-runtime-return-scope"
                          className="mt-1 font-mono text-[11px] text-amber-200/90"
                        >
                          {t('runtime_return.scope', { scope: runtimeFocusScope })}
                        </p>
                      ) : null}
                    </div>
                    <button
                      type="button"
                      onClick={() => handleNavigateStep(runtimeFocusStep)}
                      className="shrink-0 rounded-lg border border-amber-300/40 bg-amber-300/10 px-3 py-1.5 text-xs font-semibold text-amber-50 transition hover:bg-amber-300/20"
                    >
                      {t('runtime_return.action', { step: `Step ${runtimeFocusStep}` })}
                    </button>
                  </div>
                </section>
              ) : null}

              {/* 步驟元件 */}
              {renderContent()}

              {/* 底部導航欄 */}
              <div className="mt-6 flex items-center justify-between pt-4 border-t border-slate-800/80">
                <Button
                  onClick={goBack}
                  disabled={current === 1}
                  variant="ghost"
                  icon={<Icon name="chevron" className="w-4 h-4 rotate-180" />}
                >
                  上一步
                </Button>
                <div className="flex items-center gap-2 text-[11px] text-slate-500">
                  <span className="kbd">Ctrl</span>
                  <span>+</span>
                  <span className="kbd">B</span>
                  <span>收合側欄</span>
                </div>
                <div className="invisible">
                  <Button>placeholder</Button>
                </div>
              </div>
            </>
          ) : (
            // 設定元件
            renderContent()
          )}
        </main>

        {/* 4. 右側即時設定摘要欄 */}
        {showSummaryRail && view === 'flow' && (
          <aside className="hidden xl:block flex-shrink-0 w-[260px]">
            <div className="sticky top-[68px]">
              <div className="text-[10px] font-semibold uppercase tracking-wider text-slate-500 mb-2 px-1">
                即時設定摘要
              </div>
              <SummaryRail
                state={state}
                workspaceReadiness={workspaceReadiness}
                workspaceAuditHistory={workspaceAuditHistory}
                workspaceAuditUnavailable={workspaceAuditUnavailable}
                onOpenRuntime={handleRuntimeDashboardHandoff}
              />
            </div>
          </aside>
        )}
      </div>

      {/* 5. 浮動 Tweaks 面板 */}
      <TweaksPanel
        sidebarCollapsed={sidebarCollapsed}
        onToggleSidebar={toggleSidebar}
        showSummaryRail={showSummaryRail}
        onToggleSummaryRail={toggleSummaryRail}
        onResetFlow={handleReset}
      />
    </div>
  );
};

function runtimeReturnGuidanceFor(
  issue: StudioV2WorkspaceReadinessIssue,
  t: TFunction<'workbench-v2'>,
): string | null {
  switch (issue.code) {
    case 'point-missing':
      return t('runtime_return.point_missing_action');
    default:
      return null;
  }
}

export default WorkbenchV2Shell;
