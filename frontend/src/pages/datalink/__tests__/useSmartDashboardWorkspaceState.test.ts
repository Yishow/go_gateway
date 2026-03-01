import { renderHook } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  useSmartDashboardWorkspaceState,
  type UseSmartDashboardWorkspaceStateInput,
} from '../smart-dashboard/useSmartDashboardWorkspaceState';

describe('useSmartDashboardWorkspaceState', () => {
  it('assembles workspace section props from workspace, sidebar and motion state', () => {
    const workspace = {
      selectedDevice: null,
    } as unknown as UseSmartDashboardWorkspaceStateInput['workspace'];
    const sidebar = {
      sidebarTab: 'plan',
    } as unknown as UseSmartDashboardWorkspaceStateInput['sidebar'];
    const sidebarPanelMotion = {
      className: 'motion-class',
      style: { transitionDuration: '1200ms' },
    };

    const { result } = renderHook(() =>
      useSmartDashboardWorkspaceState({
        workspace,
        sidebar,
        sidebarPanelMotion,
      }),
    );

    expect(result.current.workspaceProps).toBe(workspace);
    expect(result.current.sidebarProps).toBe(sidebar);
    expect(result.current.sidebarContainerClassName).toBe('motion-class');
    expect(result.current.sidebarContainerStyle).toBe(sidebarPanelMotion.style);
  });
});
