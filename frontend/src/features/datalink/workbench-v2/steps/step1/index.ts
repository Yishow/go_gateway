/**
 * Step 1 (Device Config) 模組匯出檔
 * 
 * 落地設計決策：「拆檔策略：以「資料區塊 / 互動行為」拆分而非單一巨檔」
 */

export { Step1Device } from './Step1Device';
export type { Step1DeviceProps } from './Step1Device';

export { DEVICE_COLORS, useDeviceColor } from '../../state/deviceColors';
export type { DeviceColor } from '../../state/deviceColors';
export { getStagesForProtocol, getDefaultConfig, PROTOCOLS } from '../../state/protocols';
export type { ConnectionConfigFormProps } from './ConnectionConfigForm';
export type { DeviceEditorProps } from './DeviceEditor';
export type { ConnectionTestPanelProps } from './ConnectionTestPanel';
export type { ReadinessStagesProps } from './ReadinessStages';
export type { DeviceTabRailProps } from './DeviceTabRail';
