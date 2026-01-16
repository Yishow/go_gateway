export interface APIResponse<T = any> {
  data?: T;
  error?: string;
  status?: string;
}

export interface ConnectRequest {
  protocol: string;
  config: Record<string, any>;
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
}

export interface ReadResponse {
  values: any[];
  count: number;
}

export interface WriteRequest {
  connection_id: string;
  operation: string;
  address: number;
  values: any[];
  symbol?: string; // Fatek
  device?: string; // MC
}

export interface ConnectionState {
  id: string;
  protocol: string;
  config: Record<string, any>;
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
  data?: any;
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
