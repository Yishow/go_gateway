import { describe, it, expect } from 'vitest';
import { renderHook } from '@testing-library/react';
import { useDeviceColor, DeviceListContext, DEVICE_COLORS } from '../../../src/features/datalink/workbench-v2/state/deviceColors';
import type { Device } from '../../../src/features/datalink/workbench-v2/state/types';

describe('useDeviceColor Hook', () => {
  // 建立 Mock 裝置陣列
  const mockDevices = (count: number): Device[] =>
    Array.from({ length: count }, (_, i) => ({
      id: `dev-${i}`,
      name: `Device-${i}`,
      description: '',
      protocol: 'modbus_tcp',
      config: {},
      status: 'draft',
      test: null,
    }));

  it('在空裝置清單或找不到 ID 時，預設回傳第一個顏色 blue', () => {
    const { result } = renderHook(() => useDeviceColor('non-existent'), {
      wrapper: ({ children }) => (
        <DeviceListContext.Provider value={[]}>
          {children}
        </DeviceListContext.Provider>
      ),
    });
    expect(result.current.name).toBe('blue');
    expect(result.current.text).toBe('text-blue-400');
  });

  it('應正確根據裝置順序循環指派顏色', () => {
    const devices = mockDevices(8); // 8 個裝置

    // 測試前 6 個裝置 (0-5)
    for (let i = 0; i < 6; i++) {
      const { result } = renderHook(() => useDeviceColor(`dev-${i}`), {
        wrapper: ({ children }) => (
          <DeviceListContext.Provider value={devices}>
            {children}
          </DeviceListContext.Provider>
        ),
      });
      expect(result.current.name).toBe(DEVICE_COLORS[i]);
    }

    // 測試第 7 個裝置 (dev-6)，應 wrap-around 到 blue
    const { result: res6 } = renderHook(() => useDeviceColor('dev-6'), {
      wrapper: ({ children }) => (
        <DeviceListContext.Provider value={devices}>
          {children}
        </DeviceListContext.Provider>
      ),
    });
    expect(res6.current.name).toBe('blue');

    // 測試第 8 個裝置 (dev-7)，應 wrap-around 到 emerald
    const { result: res7 } = renderHook(() => useDeviceColor('dev-7'), {
      wrapper: ({ children }) => (
        <DeviceListContext.Provider value={devices}>
          {children}
        </DeviceListContext.Provider>
      ),
    });
    expect(res7.current.name).toBe('emerald');
  });
});
