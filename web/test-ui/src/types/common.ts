export interface SelectOption {
  label: string;
  value: string;
}

export type ConnectionStatus = 'connected' | 'disconnected' | 'connecting' | 'error';
