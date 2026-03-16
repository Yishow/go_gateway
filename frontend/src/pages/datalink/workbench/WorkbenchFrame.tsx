import type { ReactNode } from 'react';
import { WorkbenchContextBar } from './WorkbenchContextBar';
import { WorkbenchStepRail } from './WorkbenchStepRail';
import { WorkbenchInspectorPanel } from './WorkbenchInspectorPanel';
import { WorkbenchBottomSummaryBar } from './WorkbenchBottomSummaryBar';

type WorkbenchFrameProps = {
  children: ReactNode;
};

/**
 * Five-region desktop shell for the Datalink Workbench.
 *
 * Layout (1920×1080 optimized):
 * ┌──────────────────────────────────────────────────────┐
 * │                    ContextBar                        │
 * ├──────┬──────────────────────────────┬────────────────┤
 * │ Step │                              │   Inspector    │
 * │ Rail │       PrimaryWorkArea        │     Panel      │
 * │      │                              │                │
 * ├──────┴──────────────────────────────┴────────────────┤
 * │                  BottomSummaryBar                    │
 * └──────────────────────────────────────────────────────┘
 */
export function WorkbenchFrame({ children }: WorkbenchFrameProps) {
  return (
    <div
      className="grid min-h-screen gap-2 bg-slate-950 p-3 text-slate-100"
      data-testid="workbench-frame"
      style={{
        gridTemplateRows: 'auto 1fr auto',
        gridTemplateColumns: '200px 1fr 280px',
        gridTemplateAreas: `
          "context  context   context"
          "rail     main      inspector"
          "summary  summary   summary"
        `,
      }}
    >
      {/* Row 1: ContextBar spans full width */}
      <div style={{ gridArea: 'context' }}>
        <WorkbenchContextBar />
      </div>

      {/* Row 2 col 1: StepRail */}
      <div style={{ gridArea: 'rail' }}>
        <WorkbenchStepRail />
      </div>

      {/* Row 2 col 2: PrimaryWorkArea */}
      <main
        className="overflow-auto rounded-2xl border border-slate-800 bg-slate-900/70 p-5"
        style={{ gridArea: 'main' }}
        data-testid="workbench-primary-work-area"
      >
        {children}
      </main>

      {/* Row 2 col 3: InspectorPanel */}
      <div style={{ gridArea: 'inspector' }}>
        <WorkbenchInspectorPanel />
      </div>

      {/* Row 3: BottomSummaryBar spans full width */}
      <div style={{ gridArea: 'summary' }}>
        <WorkbenchBottomSummaryBar />
      </div>
    </div>
  );
}
