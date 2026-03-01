import { useMemo, type ComponentProps, type CSSProperties } from 'react';
import SmartDashboardSidebar from './SmartDashboardSidebar';
import SmartDashboardWorkspace from './SmartDashboardWorkspace';

export interface UseSmartDashboardWorkspaceStateInput {
  workspace: ComponentProps<typeof SmartDashboardWorkspace>;
  sidebar: ComponentProps<typeof SmartDashboardSidebar>;
  sidebarPanelMotion: {
    className: string;
    style: CSSProperties;
  };
}

export function useSmartDashboardWorkspaceState({
  workspace,
  sidebar,
  sidebarPanelMotion,
}: UseSmartDashboardWorkspaceStateInput) {
  return useMemo(
    () => ({
      workspaceProps: workspace,
      sidebarProps: sidebar,
      sidebarContainerClassName: sidebarPanelMotion.className,
      sidebarContainerStyle: sidebarPanelMotion.style,
    }),
    [sidebar, sidebarPanelMotion.className, sidebarPanelMotion.style, workspace],
  );
}
