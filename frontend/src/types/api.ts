export interface APIResponse<T = unknown> {
  data?: T;
  error?: string;
  status?: string;
}

export interface ConnectRequest {
  protocol: string;
  config: Record<string, unknown>;
}

export interface ConnectResponse {
  connection_id: string;
  status: string;
}

export interface ReadRequest {
  connection_id: string;
  operation: string;
  address: number;
  count: number;
  symbol?: string; // Fatek
  device?: string; // MC
  unit_id?: number; // Modbus (可選，覆蓋連線配置的站號)
  station?: number; // Fatek (可選，覆蓋連線配置的站號)
}

export interface ReadResponse {
  values: unknown[];
  count: number;
}

export interface WriteRequest {
  connection_id: string;
  operation: string;
  address: number;
  values: unknown[];
  symbol?: string; // Fatek
  device?: string; // MC
  unit_id?: number; // Modbus (可選，覆蓋連線配置的站號)
  station?: number; // Fatek (可選，覆蓋連線配置的站號)
}

export interface ConnectionState {
  id: string;
  protocol: string;
  config: Record<string, unknown>;
  connected: boolean;
  created_at: string;
}

export interface ConnectionListResponse {
  connections: ConnectionState[];
}

export interface BatchOperation {
  type: 'read' | 'write';
  read_request?: ReadRequest;
  write_request?: WriteRequest;
}

export interface BatchRequest {
  connection_id: string;
  operations: BatchOperation[];
}

export interface BatchResultItem {
  success: boolean;
  data?: unknown;
  error?: string;
}

export interface BatchResponse {
  results: BatchResultItem[];
}

export interface MonitorRequest {
  connection_id: string;
  items: ReadRequest[];
  interval: number;
}
