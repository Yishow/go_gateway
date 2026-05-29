import { createContext, useContext } from 'react';
import type { Device } from './types';

/**
 * 6 色循環常數定義
 */
export const DEVICE_COLORS = ['blue', 'emerald', 'amber', 'fuchsia', 'cyan', 'rose'] as const;
export type DeviceColor = typeof DEVICE_COLORS[number];

export interface ColorTheme {
  name: DeviceColor;
  bg: string;
  border: string;
  text: string;
  solid: string;
}

/**
 * 裝置列表 Context，用於在子樹中共享裝置狀態，使色彩 Hook 能取得正確的裝置順序
 */
export const DeviceListContext = createContext<Device[]>([]);

/**
 * 依色彩名稱取得對應的 Tailwind 樣式類別集
 * 
 * @param color 色彩名稱
 * @returns 樣式物件
 */
export function getColorTheme(color: DeviceColor): ColorTheme {
  switch (color) {
    case 'emerald':
      return {
        name: 'emerald',
        bg: 'bg-emerald-500/10',
        border: 'border-emerald-500/30',
        text: 'text-emerald-400',
        solid: 'bg-emerald-500',
      };
    case 'amber':
      return {
        name: 'amber',
        bg: 'bg-amber-500/10',
        border: 'border-amber-500/30',
        text: 'text-amber-400',
        solid: 'bg-amber-500',
      };
    case 'fuchsia':
      return {
        name: 'fuchsia',
        bg: 'bg-fuchsia-500/10',
        border: 'border-fuchsia-500/30',
        text: 'text-fuchsia-400',
        solid: 'bg-fuchsia-500',
      };
    case 'cyan':
      return {
        name: 'cyan',
        bg: 'bg-cyan-500/10',
        border: 'border-cyan-500/30',
        text: 'text-cyan-400',
        solid: 'bg-cyan-500',
      };
    case 'rose':
      return {
        name: 'rose',
        bg: 'bg-rose-500/10',
        border: 'border-rose-500/30',
        text: 'text-rose-400',
        solid: 'bg-rose-500',
      };
    case 'blue':
    default:
      return {
        name: 'blue',
        bg: 'bg-blue-500/10',
        border: 'border-blue-500/30',
        text: 'text-blue-400',
        solid: 'bg-blue-500',
      };
  }
}

/**
 * 自訂 Hook，根據裝置 ID 在裝置清單中的索引計算並回傳其專屬的循環色彩主題
 * 
 * @param deviceId 裝置 ID
 * @returns 色彩主題物件
 */
export function useDeviceColor(deviceId: string): ColorTheme {
  const devices = useContext(DeviceListContext);
  const index = devices.findIndex((d) => d.id === deviceId);
  const colorIndex = (index === -1 ? 0 : index) % 6;
  const colorName = DEVICE_COLORS[colorIndex];
  return getColorTheme(colorName);
}
