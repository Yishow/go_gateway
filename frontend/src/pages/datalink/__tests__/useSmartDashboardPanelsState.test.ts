import { renderHook } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { useSmartDashboardPanelsState } from '../smart-dashboard/useSmartDashboardPanelsState';
import type { Point } from '../../../types/datalink';

describe('useSmartDashboardPanelsState', () => {
  it('computes panel visibility and title by panel type', () => {
    const setPanelType = vi.fn();
    const setSelectedAddresses = vi.fn();
    const selectedPoint = { id: 'point-1' } as Point;

    const { result, rerender } = renderHook(
      ({ panelType }: { panelType: 'batch' | 'detail' | 'shortcuts' | null }) =>
        useSmartDashboardPanelsState({
          panelType,
          setPanelType,
          selectedPoint: panelType === 'detail' ? selectedPoint : null,
          setSelectedAddresses,
          shortcutsTitle: '快捷鍵',
          pointDetailTitle: '點位詳情',
          onBatchClose: vi.fn(),
        }),
      { initialProps: { panelType: 'batch' } as { panelType: 'batch' | 'detail' | 'shortcuts' | null } },
    );

    expect(result.current.isBatchDialogOpen).toBe(true);
    expect(result.current.isSlidePanelOpen).toBe(false);

    rerender({ panelType: 'shortcuts' });
    expect(result.current.isBatchDialogOpen).toBe(false);
    expect(result.current.isSlidePanelOpen).toBe(true);
    expect(result.current.slidePanelTitle).toBe('快捷鍵');
    expect(result.current.showShortcutsPanel).toBe(true);

    rerender({ panelType: 'detail' });
    expect(result.current.slidePanelTitle).toBe('點位詳情');
    expect(result.current.showDetailPanel).toBe(true);
  });

  it('closeBatch clears selection, resets panel and triggers batch-close callback', () => {
    const setPanelType = vi.fn();
    const setSelectedAddresses = vi.fn();
    const onBatchClose = vi.fn();

    const { result } = renderHook(() =>
      useSmartDashboardPanelsState({
        panelType: 'batch',
        setPanelType,
        selectedPoint: null,
        setSelectedAddresses,
        shortcutsTitle: '快捷鍵',
        pointDetailTitle: '點位詳情',
        onBatchClose,
      }),
    );

    result.current.closeBatch();

    expect(setPanelType).toHaveBeenCalledWith(null);
    expect(setSelectedAddresses).toHaveBeenCalledWith([]);
    expect(onBatchClose).toHaveBeenCalledTimes(1);
  });

  it('handles dialog and slide panel close actions consistently', () => {
    const setPanelType = vi.fn();
    const onBatchClose = vi.fn();

    const { result } = renderHook(() =>
      useSmartDashboardPanelsState({
        panelType: 'batch',
        setPanelType,
        selectedPoint: null,
        setSelectedAddresses: vi.fn(),
        shortcutsTitle: '快捷鍵',
        pointDetailTitle: '點位詳情',
        onBatchClose,
      }),
    );

    result.current.handleBatchDialogOpenChange(true);
    expect(setPanelType).not.toHaveBeenCalled();
    expect(onBatchClose).not.toHaveBeenCalled();

    result.current.handleBatchDialogOpenChange(false);
    expect(setPanelType).toHaveBeenCalledWith(null);
    expect(onBatchClose).toHaveBeenCalledTimes(1);

    result.current.handlePanelClose();
    expect(setPanelType).toHaveBeenCalledWith(null);
  });
});
