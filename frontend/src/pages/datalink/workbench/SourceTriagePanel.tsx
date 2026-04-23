import type { ReactNode } from 'react';

type SourceTriagePanelProps = {
  children: ReactNode;
  recoveryPanel: ReactNode;
};

export function SourceTriagePanel({
  children,
  recoveryPanel,
}: SourceTriagePanelProps) {
  return (
    <div
      className="flex min-h-0 flex-1 flex-col gap-4 overflow-hidden"
      data-testid="source-triage-workspace"
    >
      {recoveryPanel}
      {children}
    </div>
  );
}
