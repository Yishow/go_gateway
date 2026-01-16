/**
 * 監控項目介面
 */
export interface MonitorItem {
  id: string;
  operation: string;
  address: number;
  count: number;
  symbol?: string; // Fatek
  device?: string; // MC Protocol
  label?: string; // 自定義標籤
}

/**
 * Profile 配置介面
 * 包含完整的連線配置資訊
 */
export interface Profile {
  /** Profile 唯一識別碼 */
  id: string;
  /** Profile 顯示名稱 */
  name: string;
  /** 通訊協定 (例如: modbus_tcp, fatek_serial) */
  protocol: string;
  /** 連線模式 (tcp, udp, serial) */
  connectionMode: string;
  /** 連線配置參數 */
  config: Record<string, any>;
  /** 監控配置 */
  monitorConfig?: {
    items: MonitorItem[];
    interval: number;
  };
  /** 建立時間戳記 */
  createdAt: number;
  /** 最後更新時間戳記 */
  updatedAt: number;
}

/**
 * Profile 管理相關的常數
 */
export const PROFILE_STORAGE_KEY = 'test_ui_profiles';
export const DEFAULT_PROFILE_NAME = '預設配置';
