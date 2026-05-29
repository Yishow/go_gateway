import { describe, it, expect } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useWorkbenchV2State } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';

describe('useWorkbenchV2State Hook', () => {
  it('should initialize with default state', () => {
    const { result } = renderHook(() => useWorkbenchV2State());
    expect(result.current.state.view).toBe('flow');
    expect(result.current.state.current).toBe(1);
    expect(result.current.state.sidebarCollapsed).toBe(false);
    expect(result.current.state.showSummaryRail).toBe(true);
    expect(result.current.state.devices[0].name).toBe('PLC-生產線-01');
    expect(result.current.state.rules[0].name).toBe('Holding Registers');
  });

  it('should toggle sidebar and showSummaryRail', () => {
    const { result } = renderHook(() => useWorkbenchV2State());
    
    act(() => {
      result.current.toggleSidebar();
    });
    expect(result.current.state.sidebarCollapsed).toBe(true);

    act(() => {
      result.current.toggleSummaryRail();
    });
    expect(result.current.state.showSummaryRail).toBe(false);
  });

  it('should set sidebar collapsed and show summary rail explicitly', () => {
    const { result } = renderHook(() => useWorkbenchV2State());

    act(() => {
      result.current.setSidebarCollapsed(true);
    });
    expect(result.current.state.sidebarCollapsed).toBe(true);

    act(() => {
      result.current.setShowSummaryRail(false);
    });
    expect(result.current.state.showSummaryRail).toBe(false);
  });

  it('should complete step and change current step', () => {
    const { result } = renderHook(() => useWorkbenchV2State());

    act(() => {
      result.current.completeStep(1);
    });
    expect(result.current.state.completed.has(1)).toBe(true);

    act(() => {
      result.current.setCurrent(2);
    });
    expect(result.current.state.current).toBe(2);
  });

  it('should switch view modes', () => {
    const { result } = renderHook(() => useWorkbenchV2State());

    act(() => {
      result.current.setView('settings');
    });
    expect(result.current.state.view).toBe('settings');
  });

  it('should reset flow state', () => {
    const { result } = renderHook(() => useWorkbenchV2State());

    act(() => {
      result.current.completeStep(1);
      result.current.setCurrent(2);
      result.current.setView('settings');
      result.current.resetFlow();
    });

    expect(result.current.state.view).toBe('flow');
    expect(result.current.state.current).toBe(1);
    expect(result.current.state.completed.has(1)).toBe(false);
  });
});
