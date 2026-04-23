import type { ReactNode } from 'react';

type SourceRuleLayerPanelProps = {
  children: ReactNode;
  panelTestId?: string;
};

export function SourceRuleLayerPanel({
  children,
  panelTestId = 'source-rule-layer-panel',
}: SourceRuleLayerPanelProps) {
  return (
    <div className="flex min-h-0 min-w-0 flex-1 flex-col overflow-hidden" data-testid={panelTestId}>
      {children}
    </div>
  );
}
