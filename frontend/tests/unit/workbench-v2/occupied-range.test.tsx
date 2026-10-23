import { render, screen, cleanup } from '@testing-library/react';
import { afterEach, expect, it } from 'vitest';
import i18n from '@/i18n/config';
import { RangeSummary } from '@/features/datalink/workbench-v2/steps/step2/RangeSummary';
import type { ProtocolType } from '@/types/datalink';
afterEach(cleanup);
it.each([
  ['mc_3e', 'D0', 4, 'float32', 'D7'],
  ['mc_3e', 'D0', 4, 'float64', 'D15'],
  ['mc_3e', 'D0', 1, 'int64', 'D3'],
  ['mc_3e', 'D0', 1, 'uint64', 'D3'],
  ['mc_3e', 'D9', 2, 'uint16', 'D10'],
  ['mc_3e', 'XF', 2, 'bool', 'X10'],
  ['modbus_tcp', '40001', 4, 'float32', '40008'],
  ['modbus_tcp', '40001', 4, 'float64', '40016'],
  ['modbus_tcp', '40001', 1, 'int64', '40004'],
  ['modbus_tcp', '40001', 1, 'uint64', '40004'],
  ['modbus_tcp', '40009', 2, 'uint16', '40010'],
  ['modbus_tcp', '00009', 2, 'bool', '00010'],
  ['modbus_tcp', '40001', 1, 'uint16', '40001'],
] as const)('complete occupied range %s %s %s %s ends at %s', async (protocol, start, count, dataType, end) => {
  await i18n.changeLanguage('en');
  render(<RangeSummary protocol={protocol as ProtocolType} startAddress={start} count={count} dataType={dataType} />);
  expect(screen.getByTestId('range-summary')).toHaveTextContent(`${start} ~ ${end}`);
});
