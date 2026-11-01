import { useState } from 'react';
import { fireEvent, render, screen, within } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { RecordingDeviceSourceForm } from '@/features/datalink/workbench-v2/recording-editor/RecordingDeviceSourceForm';
import { parsedDeviceSourceBuffer } from '@/features/datalink/workbench-v2/recording-editor/deviceSourceBuffer';
import { isStudioV2DeviceValid, toBackendConnectionConfig } from '@/features/datalink/workbench-v2/state/studioV2DeviceAutosave';
import type { DeviceSourceBuffer } from '@/pages/datalink/workbench-v2/deviceSourceSubmission';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));

function source(protocol: DeviceSourceBuffer['protocol'], config: DeviceSourceBuffer['config']): DeviceSourceBuffer {
  return { name: 'Line PLC', description: 'Preserved description', protocol, config };
}
function Harness({ initial }: { initial: DeviceSourceBuffer }) {
  const [buffer, setBuffer] = useState(initial);
  const parsed = parsedDeviceSourceBuffer(buffer);
  return <><RecordingDeviceSourceForm buffer={buffer} lockedProtocol={false} onChange={setBuffer} onProtocol={setBuffer} />
    <output data-testid="raw-config">{JSON.stringify(buffer.config)}</output>
    <output data-testid="submitted-config">{parsed ? JSON.stringify(toBackendConnectionConfig({ ...parsed, id: 'D', status: 'draft', test: null })) : 'invalid'}</output>
  </>;
}
function submitted() { return JSON.parse(screen.getByTestId('submitted-config').textContent!); }
function change(id: string, value: string) { fireEvent.change(screen.getByTestId(id), { target: { value } }); }
const rtu = { serial_port: 'COM7', baud_rate: 19200, slave_id: 7, parity: 'even', data_bits: 7,
  stop_bits: 2, timeout: 9, data_format: 'CDAB', retained_vendor_option: { enabled: true } };
const serial = { mode: 'serial', serial_port: 'COM8', baud_rate: 4800, station_no: 12,
  parity: 'odd', data_bits: 7, stop_bits: 2, timeout: 11, retained_vendor_option: 'keep' };

describe('editor protocol parameter parity through the existing submit adapter', () => {
  it('offers five usable protocols and explicitly disables MQTT without collecting credentials', () => {
    render(<Harness initial={source('modbus_tcp', { host: 'plc.local', port: 502, slave_id: 1, timeout: 5 })} />);
    const picker = screen.getByRole('combobox', { name: 'device.protocol' });
    expect(within(picker).getAllByRole('option').filter(option => !(option as HTMLOptionElement).disabled)).toHaveLength(5);
    expect(within(picker).getByRole('option', { name: /mqtt.*device.unsupported/ })).toBeDisabled();
    expect(document.querySelector('input[type="password"]')).toBeNull();
  });

  it('shows every RTU parameter and preserves hidden values when only the serial port changes', () => {
    render(<Harness initial={source('modbus_rtu', rtu)} />);
    expect(screen.getByTestId('select-data-bits')).toHaveValue('7');
    expect(screen.getByTestId('select-stop-bits')).toHaveValue('2');
    expect(screen.getByTestId('input-timeout')).toHaveValue('9');
    expect(screen.getByTestId('select-data-format')).toHaveValue('CDAB');
    expect(screen.getByTestId('select-parity')).toHaveValue('E');
    change('input-rtu-port', 'COM9');
    expect(submitted()).toEqual({ ...rtu, serial_port: 'COM9' });
  });

  it('offers exactly the current server RTU baud enum, not the old 2400/4800 options', () => {
    render(<Harness initial={source('modbus_rtu', rtu)} />);
    expect(within(screen.getByTestId('select-baud')).getAllByRole('option').map(option => (option as HTMLOptionElement).value))
      .toEqual(['9600', '19200', '38400', '57600', '115200']);
  });

  it.each([2400, 4800, 9600.5])('blocks an unsupported RTU baud %s without mutating the buffer', baud => {
    const buffer = source('modbus_rtu', { ...rtu, baud_rate: baud });
    expect(parsedDeviceSourceBuffer(buffer)).toBeUndefined();
    expect(buffer.config.baud_rate).toBe(baud);
  });

  it.each([9600, 19200, 38400, 57600, 115200])('accepts server-supported RTU baud %s', baud => {
    expect(parsedDeviceSourceBuffer(source('modbus_rtu', { ...rtu, baud_rate: String(baud) }))?.config.baud_rate).toBe(baud);
  });

  it('keeps invalid new timeout text and does not submit a fallback value', () => {
    render(<Harness initial={source('modbus_rtu', rtu)} />);
    change('input-timeout', 'later');
    expect(JSON.parse(screen.getByTestId('raw-config').textContent!).timeout).toBe('later');
    expect(screen.getByTestId('submitted-config')).toHaveTextContent('invalid');
  });

  it('loads FATEK serial without TCP substitution and round-trips its complete configuration', () => {
    render(<Harness initial={source('fatek_fbs', serial)} />);
    expect(screen.getByTestId('select-fatek-mode')).toHaveValue('serial');
    expect(screen.queryByTestId('input-host')).not.toBeInTheDocument();
    expect(screen.getByTestId('input-rtu-port')).toHaveValue('COM8');
    expect(screen.getByTestId('select-data-bits')).toHaveValue('7');
    expect(screen.getByTestId('select-stop-bits')).toHaveValue('2');
    change('input-slave-id', '13');
    expect(submitted()).toEqual({ ...serial, station_no: 13 });
  });

  it('lets the operator explicitly select the FATEK serial branch', () => {
    render(<Harness initial={source('fatek_fbs', { mode: 'tcp', host: 'fatek.local', port: 500, station_no: 1, timeout: 5 })} />);
    change('select-fatek-mode', 'serial');
    expect(screen.getByTestId('input-rtu-port')).toBeInTheDocument();
    expect(screen.queryByTestId('input-host')).not.toBeInTheDocument();
    expect(JSON.parse(screen.getByTestId('raw-config').textContent!).mode).toBe('serial');
  });

  it('can submit a new FATEK serial source after explicit mode selection and entering its port', () => {
    render(<Harness initial={source('fatek_fbs', { mode: 'tcp', host: 'fatek.local', port: 500, station_no: 1, timeout: 5 })} />);
    change('select-fatek-mode', 'serial');
    change('input-rtu-port', 'COM10');
    expect(screen.getByTestId('submitted-config')).not.toHaveTextContent('invalid');
    expect(submitted()).toMatchObject({ mode: 'serial', serial_port: 'COM10', baud_rate: 9600,
      parity: 'even', data_bits: 7, stop_bits: 1, station_no: 1, timeout: 5 });
  });

  it('does not impose the Modbus RTU baud enum on FATEK serial', () => {
    // FATEK schema baud_rate is an integer, with no RTU-specific enum.
    expect(parsedDeviceSourceBuffer(source('fatek_fbs', serial))?.config.baud_rate).toBe(4800);
  });

  it.each(['modbus_tcp', 'modbus_udp'] as const)('exposes and preserves %s connection data format', protocol => {
    render(<Harness initial={source(protocol, { host: 'plc.local', port: 502, slave_id: 4, timeout: 8, data_format: 'BADC', retained: 42 })} />);
    expect(screen.getByTestId('select-data-format')).toHaveValue('BADC');
    change('input-host', 'plc-new.local');
    expect(submitted()).toEqual({ host: 'plc-new.local', port: 502, slave_id: 4, timeout: 8, data_format: 'BADC', retained: 42 });
  });

  it('displays MC legacy aliases and updates one canonical field without losing other settings', () => {
    render(<Harness initial={source('mc_3e', { host: 'mc.local', port: 6000, station: 6, module_io_no: 1023,
      module_station_no: 6, network_no: 2, pc_no: 254, timeout: 12, data_format: 'DCBA', retained: 'untouched' })} />);
    expect(screen.getByTestId('input-network-no')).toHaveValue('2');
    expect(screen.getByTestId('input-pc-no')).toHaveValue('254');
    expect(screen.getByTestId('input-io-no')).toHaveValue('1023');
    expect(screen.getByTestId('select-data-format')).toHaveValue('DCBA');
    change('input-io-no', '1000');
    expect(submitted()).toEqual({ host: 'mc.local', port: 6000, station_no: 6, io_no: 1000,
      network_no: 2, pc_no: 254, timeout: 12, data_format: 'DCBA', retained: 'untouched' });
  });

  it('accepts an MC module_station_no-only legacy alias when editing another field', () => {
    render(<Harness initial={source('mc_3e', { host: 'mc.local', port: 6000, module_station_no: 6,
      module_io_no: 1023, network_no: 0, pc_no: 255, timeout: 5 })} />);
    expect(screen.getByTestId('input-slave-id')).toHaveValue('6');
    change('input-host', 'mc-new.local');
    expect(screen.getByTestId('submitted-config')).not.toHaveTextContent('invalid');
    expect(submitted()).toEqual({ host: 'mc-new.local', port: 6000, station_no: 6,
      io_no: 1023, network_no: 0, pc_no: 255, timeout: 5 });
  });

  it.each([-1, '', 'not-a-station'])('blocks invalid MC module-only station %s without altering raw aliases', station => {
    const buffer = source('mc_3e', { host: 'mc.local', port: 6000, module_station_no: station, timeout: 5 });
    expect(parsedDeviceSourceBuffer(buffer)).toBeUndefined();
    expect(buffer.config).toEqual({ host: 'mc.local', port: 6000, module_station_no: station, timeout: 5 });
  });

  it('accepts module-only MC station in the shared validator used by the original advanced form', () => {
    const device = { ...source('mc_3e', { host: 'mc.local', port: 6000, module_station_no: 6, timeout: 5 }),
      id: 'D', status: 'draft' as const, test: null };
    expect(isStudioV2DeviceValid(device)).toBe(true);
    expect(device.config).toEqual({ host: 'mc.local', port: 6000, module_station_no: 6, timeout: 5 });
  });

  it.each([-1, '', 'bad', '6'])('keeps the shared validator strict for invalid raw MC module station %s', station => {
    expect(isStudioV2DeviceValid({ ...source('mc_3e', { host: 'mc.local', port: 6000, module_station_no: station, timeout: 5 }),
      id: 'D', status: 'draft', test: null })).toBe(false);
  });

  it('uses canonical MC station over an invalid module alias in the shared validator', () => {
    expect(isStudioV2DeviceValid({ ...source('mc_3e', { host: 'mc.local', port: 6000, station_no: 3, module_station_no: -1, timeout: 5 }),
      id: 'D', status: 'draft', test: null })).toBe(true);
  });

  it('labels device format separately from range and point format', () => {
    render(<Harness initial={source('modbus_tcp', { host: 'plc.local', port: 502, slave_id: 1, timeout: 5 })} />);
    expect(screen.getByLabelText('connection.dataFormat')).toBe(screen.getByTestId('select-data-format'));
    expect(within(screen.getByTestId('select-data-format')).getByRole('option', { name: 'connection.defaultFormat' })).toHaveValue('');
  });

  it('preserves both existing FATEK mode configurations and unknown fields across explicit switches', () => {
    render(<Harness initial={source('fatek_fbs', { ...serial, host: 'fatek.local', port: 501 })} />);
    change('select-fatek-mode', 'tcp');
    expect(screen.getByTestId('input-host')).toHaveValue('fatek.local');
    expect(screen.getByTestId('input-port')).toHaveValue('501');
    change('select-fatek-mode', 'serial');
    expect(submitted()).toEqual({ ...serial, host: 'fatek.local', port: 501 });
  });

  it('gives canonical MC values priority over conflicting legacy aliases without changing untouched values', () => {
    render(<Harness initial={source('mc_3e', { host: 'mc.local', port: 5012, station_no: 3, station: 9,
      io_no: 800, module_io_no: 1023, network_no: 4, pc_no: 253, timeout: 6, data_format: 'CDAB' })} />);
    expect(screen.getByTestId('input-slave-id')).toHaveValue('3');
    expect(screen.getByTestId('input-io-no')).toHaveValue('800');
    change('input-host', 'mc-new.local');
    expect(submitted()).toEqual({ host: 'mc-new.local', port: 5012, station_no: 3,
      io_no: 800, network_no: 4, pc_no: 253, timeout: 6, data_format: 'CDAB' });
  });
});
