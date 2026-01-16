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
