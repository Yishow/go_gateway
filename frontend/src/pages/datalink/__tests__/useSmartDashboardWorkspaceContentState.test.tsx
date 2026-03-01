import { render, renderHook } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { useSmartDashboardWorkspaceContentState } from '../smart-dashboard/useSmartDashboardWorkspaceContentState';

describe('useSmartDashboardWorkspaceContentState', () => {
  it('derives planner bindings and clamps planner inputs', () => {
    const setPlanStartAddress = vi.fn();
    const setPlanCount = vi.fn();
    const getGridCenterAddress = vi.fn(() => '40001');

    const { result } = renderHook(() =>
      useSmartDashboardWorkspaceContentState({
        planStartAddress: '10008',
        setPlanStartAddress,
        setPlanCount,
        selectedDeviceProtocol: 'modbus_tcp',
        gridViewStartAddress: '',
        getGridCenterAddress,
        onGridViewShift: vi.fn(),
      }),
    );

    expect(result.current.modbusArea).toBe('1');
    expect(result.current.gridRange).toBe(100);
    expect(result.current.centerAddress).toBe('10008');

    result.current.handleModbusAreaChange('3');
    expect(setPlanStartAddress).toHaveBeenCalledWith('30001');

    result.current.handlePlanStartAddressChange('d10');
    expect(setPlanStartAddress).toHaveBeenCalledWith('D10');

    result.current.handlePlanCountChange('999');
    result.current.handlePlanCountChange('0');
    result.current.handlePlanCountChange('12.7');
    result.current.handlePlanCountChange('oops');

    expect(setPlanCount).toHaveBeenNthCalledWith(1, 200);
    expect(setPlanCount).toHaveBeenNthCalledWith(2, 1);
    expect(setPlanCount).toHaveBeenNthCalledWith(3, 12);
    expect(setPlanCount).toHaveBeenNthCalledWith(4, 1);
  });

  it('prioritizes grid view start address over planner and protocol defaults', () => {
    const getGridCenterAddress = vi.fn(() => 'D0');
    const { result, rerender } = renderHook(
      ({
        gridViewStartAddress,
        planStartAddress,
      }: {
        gridViewStartAddress: string;
        planStartAddress: string;
      }) =>
        useSmartDashboardWorkspaceContentState({
          planStartAddress,
          setPlanStartAddress: vi.fn(),
          setPlanCount: vi.fn(),
          selectedDeviceProtocol: 'fatek_fbs',
          gridViewStartAddress,
          getGridCenterAddress,
          onGridViewShift: vi.fn(),
        }),
      {
        initialProps: {
          gridViewStartAddress: 'D200',
          planStartAddress: 'D100',
        },
      },
    );

    expect(result.current.centerAddress).toBe('D200');

    rerender({ gridViewStartAddress: '', planStartAddress: 'D100' });
    expect(result.current.centerAddress).toBe('D100');

    rerender({ gridViewStartAddress: '', planStartAddress: '' });
    expect(result.current.centerAddress).toBe('D0');
    expect(getGridCenterAddress).toHaveBeenCalledWith('fatek_fbs');
  });

  it('throttles wheel shifts within 300ms and keeps preventDefault behavior', () => {
    const onGridViewShift = vi.fn();
    let now = 1000;
    const nowSpy = vi.spyOn(Date, 'now').mockImplementation(() => now);

    function Harness() {
      const state = useSmartDashboardWorkspaceContentState({
        planStartAddress: '40001',
        setPlanStartAddress: vi.fn(),
        setPlanCount: vi.fn(),
        selectedDeviceProtocol: 'modbus_tcp',
        gridViewStartAddress: '',
        getGridCenterAddress: () => '40001',
        onGridViewShift,
      });

      return <div data-testid="scroll" ref={state.gridScrollRef} />;
    }

    const { getByTestId } = render(<Harness />);
    const scroll = getByTestId('scroll');

    const first = new WheelEvent('wheel', { deltaY: 100, cancelable: true });
    const second = new WheelEvent('wheel', { deltaY: 120, cancelable: true });
    const third = new WheelEvent('wheel', { deltaY: -10, cancelable: true });

    expect(scroll.dispatchEvent(first)).toBe(false);
    now = 1100;
    expect(scroll.dispatchEvent(second)).toBe(false);
    now = 1401;
    expect(scroll.dispatchEvent(third)).toBe(false);
    expect(onGridViewShift).toHaveBeenNthCalledWith(1, 100);
    expect(onGridViewShift).toHaveBeenNthCalledWith(2, -100);

    nowSpy.mockRestore();
  });
});
