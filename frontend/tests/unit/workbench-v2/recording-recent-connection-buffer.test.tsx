import { act, renderHook } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { useRecordingDeviceBuffer } from '../../../src/features/datalink/workbench-v2/recording-editor/useRecordingDeviceBuffer';
import type { DeviceSourcePort } from '../../../src/pages/datalink/workbench-v2/deviceSourceSubmission';
const cached = { protocol: 'fatek_fbs' as const, config: { mode: 'serial', serial_port: '/dev/ttyS3', baud_rate: 19200, parity: 'even', data_bits: 7, stop_bits: 1, station_no: 4, timeout: 5 } };
function setup() {
  const recentConnections = { read: vi.fn(() => structuredClone(cached)), clear: vi.fn() };
  const port: DeviceSourcePort = { recentConnections, submit: vi.fn(), inspect: vi.fn(), observations: () => [] };
  const hook = renderHook(() => useRecordingDeviceBuffer('workspace-1', port));
  return { ...hook, port, recentConnections };
}
describe('new form recent connection initialization', () => {
  it('preselects last successful subtype only on a fresh form and keeps name new', async () => {
    const current = setup(); await act(async () => { await current.result.current.start(); });
    expect(current.result.current.buffer).toEqual({ name: '', description: '', ...cached });
    expect(current.result.current.reused).toBe(true);
    expect(current.result.current.dirty).toBe(false);
    expect(current.port.submit).not.toHaveBeenCalled();
  });
  it('does not apply a late cache update over typed values or an explicit protocol', async () => {
    const current = setup();
    act(() => current.result.current.change({ ...current.result.current.buffer, name: 'Typed first', config: { host: '192.0.2.90', port: 'raw' } }));
    await act(async () => { await current.result.current.start(); });
    expect(current.result.current.buffer).toMatchObject({ name: 'Typed first', protocol: 'modbus_tcp', config: { host: '192.0.2.90', port: 'raw' } });
    expect(current.recentConnections.read).not.toHaveBeenCalled();
  });
  it('clears all current-protocol subtypes, preserving name and selected mode, without mutation', async () => {
    const current = setup(); await act(async () => { await current.result.current.start(); });
    act(() => current.result.current.change({ ...current.result.current.buffer, name: 'Keep new name' }));
    act(() => current.result.current.clearRecent());
    expect(current.recentConnections.clear).toHaveBeenCalledWith('workspace-1', 'fatek_fbs');
    expect(current.result.current.buffer).toMatchObject({ name: 'Keep new name', protocol: 'fatek_fbs', config: { mode: 'serial', baud_rate: 9600, station_no: 1 } });
    expect(current.result.current.buffer.config).not.toHaveProperty('host');
    expect(current.result.current.buffer.config.serial_port ?? '').toBe('');
    expect(current.result.current.reused).toBe(false);
    expect(current.port.submit).not.toHaveBeenCalled();
  });
});


describe('late page-cache changes cannot rewrite an initialized form', () => {
  it('does not refresh an opened untouched form when another save changes recent settings', async () => {
    const current = setup(); await act(async () => { await current.result.current.start(); });
    current.recentConnections.read.mockImplementation(() => ({ ...cached, config: { ...cached.config, serial_port: '/dev/ttyS9' } }));
    await act(async () => { await current.result.current.start(); });
    expect(current.result.current.buffer.config.serial_port).toBe('/dev/ttyS3');
    expect(current.recentConnections.read).toHaveBeenCalledTimes(1);
  });
  it('does not overwrite typing while the existing prepare/readback is awaiting permission', async () => {
    let permit!: (value: boolean) => void;
    const recentConnections = { read: vi.fn(() => structuredClone(cached)), clear: vi.fn() };
    const port: DeviceSourcePort = { recentConnections, prepare: () => new Promise(resolve => { permit = resolve; }), submit: vi.fn(), inspect: vi.fn(), observations: () => [] };
    const current = renderHook(() => useRecordingDeviceBuffer('workspace-1', port));
    let started!: Promise<boolean>;
    act(() => { started = current.result.current.start(); });
    act(() => current.result.current.change({ ...current.result.current.buffer, name: 'Typed during read', protocol: 'modbus_udp', config: { host: 'typed.host', port: 'bad text' } }));
    await act(async () => { permit(true); expect(await started).toBe(true); });
    expect(current.result.current.buffer).toMatchObject({ name: 'Typed during read', protocol: 'modbus_udp', config: { host: 'typed.host', port: 'bad text' } });
    expect(recentConnections.read).not.toHaveBeenCalled();
  });
});
