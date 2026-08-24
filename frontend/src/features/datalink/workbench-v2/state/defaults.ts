import { getDefaultPlannerStartAddress } from '../../../../utils/addressParser';
import type { ProtocolType } from '../../../../types/datalink';
import type { Device, Rule } from './types';

export const DEFAULT_DEVICE: Device = {
  id: 'dev-01',
  name: 'PLC-生產線-01',
  description: 'Modbus TCP PLC (Line A 主控)',
  protocol: 'modbus_tcp',
  config: { host: '192.168.1.100', port: 502, slave_id: 1, timeout: 5 },
  status: 'draft',
  test: null,
  persisted: false,
  save_state: 'idle',
  save_error: null,
  runtime_apply_status: null,
  runtime_apply_message: null,
  availability_status: 'available',
  availability_reason: null,
  running: false,
};

/** Builds the first draft rule using the selected device protocol defaults. */
function createDefaultRule(
  deviceId: string = 'dev-01',
  protocol: ProtocolType = 'modbus_tcp',
): Rule {
  return {
    id: 'rule-01',
    device_id: deviceId,
    name: protocol === 'mc_3e' ? 'D Registers' : protocol === 'fatek_fbs' ? 'R Registers' : 'Holding Registers',
    start_address: getDefaultPlannerStartAddress(protocol),
    count: 8,
    data_type: 'int16',
    naming_prefix: 'SENSOR_',
    enabled: true,
    scale_multiplier: 0.1,
    scale_offset: 0,
    data_format: '',
    skipped_addresses: [],
    share_enabled: true,
    share_start_register: 40001,
    share_stride: null,
    persisted: false,
    save_state: 'idle',
    save_error: null,
  };
}

export const DEFAULT_RULE: Rule = createDefaultRule('dev-01', 'modbus_tcp');
