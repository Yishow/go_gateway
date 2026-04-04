import type { ReactNode } from 'react';
import { WorkbenchBottomSummaryBar } from './WorkbenchBottomSummaryBar';
import { WorkbenchContextBar } from './WorkbenchContextBar';
import { WorkbenchInspectorPanel } from './WorkbenchInspectorPanel';
import { WB_SHELL_SURFACE } from './workbenchShellTokens';

type WorkbenchFrameProps = {
  children: ReactNode;
};

/**
 * Five-region desktop shell for the Datalink Workbench.
 *
 * 外層使用 `h-dvh` 鎖定視窗高度，讓中列 `1fr` 有明確剩餘空間；否則在 `min-h-screen` 下
 * Grid 會依內容伸長，`main` 會跟著被撐高而不在區塊內捲動。
 *
 * 主工作區與頂／右／底外殼共用 `WB_SHELL_SURFACE`（見 workbenchShellTokens），實心面板統一層次。
 * 主欄位預設 **不捲動**（`overflow-hidden`），由各步驟內部決定捲動容器（例如來源步驟僅記憶體格區捲動）。
 *
 * Layout（步驟導覽已併入 ContextBar 頂列）：
 * ┌──────────────────────────────────────────────────────┐
 * │              ContextBar（含水平步驟）                 │
 * ├──────────────────────────────┬─────────────────────┤
 * │       PrimaryWorkArea        │     Inspector       │
 * ├──────────────────────────────┴─────────────────────┤
 * │                  BottomSummaryBar                  │
 * └──────────────────────────────────────────────────────┘
 */
export function WorkbenchFrame({ children }: WorkbenchFrameProps) {
  return (
    <div
      className="grid h-dvh min-h-0 max-h-dvh gap-3 overflow-hidden bg-slate-950 p-3 text-slate-100"
      data-testid="workbench-frame"
      style={{
        gridTemplateRows: 'auto 1fr auto',
        gridTemplateColumns: '1fr 280px',
        gridTemplateAreas: `
          "context    context"
          "main       inspector"
          "summary    summary"
        `,
      }}
      >
        {/* Row 1: ContextBar spans full width */}
        <div style={{ gridArea: 'context' }}>
          <WorkbenchContextBar />
        </div>

      {/* Row 2 col 1: PrimaryWorkArea */}
      <main
        className={`flex h-full min-h-0 flex-col overflow-hidden p-6 ${WB_SHELL_SURFACE}`}
        style={{ gridArea: 'main' }}
        data-testid="workbench-primary-work-area"
      >
        {children}
      </main>

      {/* Row 2 col 3: InspectorPanel */}
      <div className="min-h-0" style={{ gridArea: 'inspector' }}>
        <WorkbenchInspectorPanel />
      </div>

      {/* Row 3: BottomSummaryBar spans full width */}
      <div style={{ gridArea: 'summary' }}>
        <WorkbenchBottomSummaryBar />
      </div>
    </div>
  );
}
