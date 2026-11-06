import type { ProtocolType } from '../../../src/types/datalink';
import type { MemoryRuleInput } from '../../../src/features/datalink/workbench-v2/recording-home/recordingHomeMemoryMap';

// Deterministic display stress data, never device readings or evidence of protocol support.
export const longName = 'SAMPLE_BOILER_OUTLET_TEMPERATURE_SENSOR_NORTH_PRODUCTION_LINE_'.repeat(3);
export const longNumber = '-123456789012345678901234567890.125';
export const unknownString = 'SAMPLE_UNCONFIRMED_STRING_MUST_NOT_BE_SHOWN_AS_A_READING';
export const visualScenarios: Record<string, { protocol: ProtocolType; inputs: MemoryRuleInput[] }> = {
  long: { protocol: 'modbus_tcp', inputs: [
    { id: 'LONG', label: longName, start_address: '40001', count: 1, data_type: 'float64', skipped_addresses: [] },
    ...Array.from({ length: 32 }, (_, index) => ({ id: `LIST${index}`, label: `SAMPLE_RANGE_${index}_${longName}`,
      start_address: String(40101 + index * 4), count: 2, data_type: 'float32' as const, skipped_addresses: [] })),
  ] },
  dense256: { protocol: 'mc_3e', inputs: [{ id: 'W256', label: 'SAMPLE_256_PHYSICAL_SLOTS', start_address: 'W0',
    count: 256, data_type: 'int16', skipped_addresses: [], proven: Array.from({ length: 256 }, (_, index) => `W${index.toString(16).toUpperCase()}`) }] },
  crossing: { protocol: 'modbus_tcp', inputs: [
    // 250 is the canonical decimal window cap. 40249..40252 crosses that window.
    { id: 'CROSS', label: 'SAMPLE_FOUR_REGISTER_PAGE_CROSSING', start_address: '40249', count: 1, data_type: 'float64', skipped_addresses: [] },
    { id: 'ROW', label: 'SAMPLE_FOUR_REGISTER_ROW_CROSSING', start_address: '40008', count: 1, data_type: 'float64', skipped_addresses: [] },
  ] },
  unknownLong: { protocol: 'modbus_tcp', inputs: [{ id: 'STRING', label: longName, start_address: '40001', count: 1,
    data_type: 'string', skipped_addresses: [] }] },
};
