export interface BaseConfig {
  timeout?: number;
}

export interface TCPConfig extends BaseConfig {
  host: string;
  port: number;
}

export interface SerialConfig extends BaseConfig {
  port: string;
  baudRate: number;
  dataBits: number;
  stopBits: number;
  parity: 'N' | 'E' | 'O';
}

export interface ModbusConfig {
  unitID: number;
}

export interface FatekConfig {
  station: number;
}

export type ModbusTCPConfig = TCPConfig & ModbusConfig;
export type ModbusRTUConfig = SerialConfig & ModbusConfig;
export type FatekTCPConfig = TCPConfig & FatekConfig;
export type FatekSerialConfig = SerialConfig & FatekConfig;
export type MCTCPConfig = TCPConfig;
export type MCSerialConfig = SerialConfig;

export type ProtocolConfig = 
  | ModbusTCPConfig 
  | ModbusRTUConfig 
  | FatekTCPConfig 
  | FatekSerialConfig 
  | MCTCPConfig 
  | MCSerialConfig
  | Record<string, unknown>;
