export type ProtocolType = 
  | 'modbus_tcp' 
  | 'modbus_rtu' 
  | 'fatek_tcp' 
  | 'fatek_serial' 
  | 'mc_tcp'
  | 'mc_serial';

export type ModbusOperation = 
  | 'read_coils'
  | 'read_discrete_inputs'
  | 'read_holding_registers'
  | 'read_input_registers'
  | 'write_single_coil'
  | 'write_single_register'
  | 'write_multiple_coils'
  | 'write_multiple_registers';

export type FatekOperation = 
  | 'read_status'
  | 'read_registers'
  | 'write_status'
  | 'write_registers';

export type MCOperation = 
  | 'batch_read_word'
  | 'batch_read_bit'
  | 'batch_write_word'
  | 'batch_write_bit';

export type ProtocolOperation = ModbusOperation | FatekOperation | MCOperation;
