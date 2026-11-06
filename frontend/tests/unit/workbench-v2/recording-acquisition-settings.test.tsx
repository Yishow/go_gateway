import { useState } from 'react';
import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { RecordingDeviceSourceForm } from '@/features/datalink/workbench-v2/recording-editor/RecordingDeviceSourceForm';
import { parsedDeviceSourceBuffer } from '@/features/datalink/workbench-v2/recording-editor/deviceSourceBuffer';
import { hydrateStudioV2Device, toBackendConnectionConfig, toStudioV2DeviceUpdateRequest } from '@/features/datalink/workbench-v2/state/studioV2DeviceAutosave';
import { acquisitionInputFromSettings, acquisitionMatches, runtimeIntervalReadout } from '@/features/datalink/workbench-v2/state/acquisitionSettings';
import { projectRuntimeStatus } from '@/services/runtimeStatusProjection';
import type { DeviceSourceBuffer } from '@/pages/datalink/workbench-v2/deviceSourceSubmission';
import type { RuntimeIntervalReadout } from '@/types/acquisitionSettings';
import type { StudioV2WorkspaceDeviceRecord } from '@/services/studioV2WorkspaceDevices';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, values?: { title?: string }) => key === 'open' ? `About ${values?.title}` : key }) }));
afterEach(cleanup);

const mc = { host: 'plc.local', port: 6000, station: 0, network_no: 0, pc_no: 255, module_io_no: 1023, module_station_no: 0, timeout: 5 };
const fatek = { mode: 'tcp', host: 'plc.local', port: 500, station: 1, timeout: 5 };
const modbus = { host: 'plc.local', port: 502, slave_id: 1, timeout: 5 };
const saved = { interval_ms: 1000, auto_interval: false, stale_threshold_ms: null, frozen_threshold_ms: 0 };
function source(protocol: DeviceSourceBuffer['protocol'], config: DeviceSourceBuffer['config'], existing = false): DeviceSourceBuffer {
  return { name: 'Line PLC', description: '', protocol, config, ...(existing ? { acquisitionInput: acquisitionInputFromSettings(saved) } : {}) };
}
function Harness({ initial, runtime }: { initial: DeviceSourceBuffer; runtime?: RuntimeIntervalReadout }) {
  const [buffer, setBuffer] = useState(initial);
  const parsed = parsedDeviceSourceBuffer(buffer);
  const device = parsed ? { ...parsed, id: 'D', status: 'draft' as const, test: null } : undefined;
  return <><RecordingDeviceSourceForm buffer={buffer} lockedProtocol={false} onChange={setBuffer} onProtocol={setBuffer} runtimeInterval={runtime} />
    <output data-testid="payload">{device ? JSON.stringify(toStudioV2DeviceUpdateRequest(device)) : 'invalid'}</output></>;
}
const payload = () => { const text = screen.getByTestId('payload').textContent!; return text === 'invalid' ? undefined : JSON.parse(text); };
const declaration = () => toBackendAcquisition(payload()?.connection_config);
function toBackendAcquisition(config: Record<string, unknown> | undefined) { return config ? config.acquisition : 'invalid'; }
const options = () => within(screen.getByTestId('acquisition-profile')).getAllByRole('option').map(option => (option as HTMLOptionElement).value);
const change = (id: string, value: string) => fireEvent.change(screen.getByTestId(id), { target: { value } });

describe('connection settings acquisition optimization', () => {
  it.each(['modbus_tcp', 'modbus_udp', 'modbus_rtu'] as const)('%s offers no random profile and explains contiguous aggregation', protocol => {
    render(<Harness initial={source(protocol, protocol === 'modbus_rtu' ? { serial_port: 'COM1', baud_rate: 9600, parity: 'none', slave_id: 1, timeout: 5 } : modbus)} />);
    expect(screen.queryByTestId('acquisition-profile')).toBeNull();
    expect(screen.queryByTestId('acquisition-random-points')).toBeNull();
    expect(screen.getByTestId('acquisition-modbus-note')).toHaveTextContent('modbusNote');
    expect(screen.getByRole('button', { name: 'About acquisition.contiguous.title' })).toBeInTheDocument();
    expect(declaration()).toBeUndefined();
  });

  it('offers exactly the MC 3E profiles and writes only connection_config.acquisition', () => {
    render(<Harness initial={source('mc_3e', mc)} />);
    expect(options()).toEqual(['undeclared', 'safe-v1', 'mitsubishi-fx5-cpu-3e-binary-v1', 'mitsubishi-q-iqr-3e-binary-v1']);
    expect(screen.getByTestId('acquisition-profile')).toHaveValue('undeclared');
    expect(declaration()).toBeUndefined();
    change('acquisition-profile', 'mitsubishi-q-iqr-3e-binary-v1');
    expect(screen.getByTestId('acquisition-random-points')).toHaveValue('192');
    expect(declaration()).toEqual({ version: 1, profile: 'mitsubishi-q-iqr-3e-binary-v1' });
    change('acquisition-random-points', '120');
    expect(declaration()).toEqual({ version: 1, profile: 'mitsubishi-q-iqr-3e-binary-v1', limits: { random_points: 120 } });
    expect(Object.keys(payload())).not.toContain('acquisition_interval_ms');
  });

  it('rejects a random limit above the profile limit with a visible field error', () => {
    render(<Harness initial={source('mc_3e', mc)} />);
    change('acquisition-profile', 'mitsubishi-q-iqr-3e-binary-v1');
    for (const value of ['300', '0', '1.5', 'abc']) {
      change('acquisition-random-points', value);
      expect(payload()).toBeUndefined();
      expect(screen.getByRole('alert')).toHaveTextContent('errors.randomRange');
    }
    change('acquisition-random-points', '');
    expect(declaration()).toEqual({ version: 1, profile: 'mitsubishi-q-iqr-3e-binary-v1' });
  });

  it('offers exactly the FATEK profiles with a 64 default limit that can be lowered', () => {
    render(<Harness initial={source('fatek_fbs', fatek)} />);
    expect(options()).toEqual(['undeclared', 'safe-v1', 'fatek-fbs-ascii-v1']);
    change('acquisition-profile', 'fatek-fbs-ascii-v1');
    expect(screen.getByTestId('acquisition-random-points')).toHaveValue('64');
    change('acquisition-random-points', '48');
    expect(declaration()).toEqual({ version: 1, profile: 'fatek-fbs-ascii-v1', limits: { random_points: 48 } });
    change('acquisition-random-points', '65');
    expect(payload()).toBeUndefined();
  });

  it('drops random limits for safe-v1 and removes the declaration when undeclared', () => {
    const declared = { ...mc, acquisition: { version: 1, profile: 'mitsubishi-q-iqr-3e-binary-v1', limits: { random_points: 100, request_bytes: 900 } } };
    render(<Harness initial={source('mc_3e', declared)} />);
    expect(screen.getByTestId('acquisition-random-points')).toHaveValue('100');
    change('acquisition-profile', 'safe-v1');
    expect(screen.queryByTestId('acquisition-random-points')).toBeNull();
    expect(declaration()).toEqual({ version: 1, profile: 'safe-v1', limits: { request_bytes: 900 } });
    change('acquisition-profile', 'undeclared');
    expect(payload().connection_config).not.toHaveProperty('acquisition');
  });

  it('blocks the FX5 CPU profile outside the own-station route the backend accepts', () => {
    render(<Harness initial={source('mc_3e', { ...mc, station: 3 })} />);
    change('acquisition-profile', 'mitsubishi-fx5-cpu-3e-binary-v1');
    expect(payload()).toBeUndefined();
    expect(screen.getByRole('alert')).toHaveTextContent('errors.fx5Route');
    cleanup();
    render(<Harness initial={source('mc_3e', mc)} />);
    change('acquisition-profile', 'mitsubishi-fx5-cpu-3e-binary-v1');
    expect(declaration()).toEqual({ version: 1, profile: 'mitsubishi-fx5-cpu-3e-binary-v1' });
  });

  it('shows device acquisition fields only for a saved device, each with contextual help', () => {
    render(<Harness initial={source('mc_3e', mc)} />);
    expect(screen.queryByTestId('acquisition-interval')).toBeNull();
    expect(screen.getByTestId('acquisition-new-device-note')).toBeInTheDocument();
    cleanup();
    render(<Harness initial={source('mc_3e', mc, true)} />);
    for (const key of ['interval', 'autoInterval', 'staleThreshold', 'frozenThreshold', 'profile', 'randomPoints']) {
      if (key === 'randomPoints') change('acquisition-profile', 'mitsubishi-q-iqr-3e-binary-v1');
      const help = screen.getByRole('button', { name: `About acquisition.${key}.title` });
      fireEvent.click(help); expect(screen.getByRole('note')).toHaveTextContent(`acquisition.${key}.how`);
      fireEvent.keyDown(window, { key: 'Escape' }); expect(help).toHaveFocus();
    }
  });

  it('validates interval and threshold ranges before saving', () => {
    render(<Harness initial={source('modbus_tcp', modbus, true)} />);
    expect(payload()).toMatchObject({ acquisition_interval_ms: 1000, acquisition_auto_interval: false, frozen_threshold_ms: 0 });
    expect(payload()).not.toHaveProperty('stale_threshold_ms');
    for (const [id, value, error] of [['acquisition-interval', '99', 'errors.intervalRange'], ['acquisition-interval', '3600001', 'errors.intervalRange'], ['acquisition-interval', '', 'errors.intervalRange'],
      ['acquisition-stale', '0', 'errors.staleInvalid'], ['acquisition-frozen', '-1', 'errors.frozenInvalid'], ['acquisition-frozen', '', 'errors.frozenInvalid']] as const) {
      change(id, value); expect(payload(), `${id}=${value}`).toBeUndefined();
      expect(screen.getByRole('alert')).toHaveTextContent(error);
      change(id, id === 'acquisition-interval' ? '1000' : id === 'acquisition-stale' ? '' : '0');
    }
    change('acquisition-interval', '100'); change('acquisition-stale', '2500'); change('acquisition-frozen', '60000');
    fireEvent.click(screen.getByTestId('acquisition-auto'));
    expect(payload()).toMatchObject({ acquisition_interval_ms: 100, acquisition_auto_interval: true, stale_threshold_ms: 2500, frozen_threshold_ms: 60000 });
    change('acquisition-interval', '3600000'); expect(payload()).toMatchObject({ acquisition_interval_ms: 3600000 });
  });

  it('does not pretend a saved stale threshold can be cleared back to automatic', () => {
    const initial = { ...source('modbus_tcp', modbus), acquisitionInput: acquisitionInputFromSettings({ ...saved, stale_threshold_ms: 4000 }) };
    render(<Harness initial={initial} />);
    expect(screen.getByTestId('acquisition-stale')).toHaveValue('4000');
    change('acquisition-stale', '');
    expect(payload()).toBeUndefined();
    expect(screen.getByRole('alert')).toHaveTextContent('errors.staleRequired');
  });

  it('shows runtime configured and effective intervals, and states when they are unavailable', () => {
    render(<Harness initial={source('modbus_tcp', modbus, true)} runtime={{ status: 'ready', configured_ms: 500, effective_ms: 800 }} />);
    expect(screen.getByTestId('acquisition-runtime')).toHaveTextContent('runtime.values');
    expect(screen.getByTestId('acquisition-runtime')).toHaveTextContent('runtime.stretched');
    cleanup();
    render(<Harness initial={source('modbus_tcp', modbus, true)} runtime={{ status: 'unavailable' }} />);
    expect(screen.getByTestId('acquisition-runtime')).toHaveTextContent('runtime.unavailable');
    cleanup();
    render(<Harness initial={source('modbus_tcp', modbus, true)} runtime={{ status: 'not-running' }} />);
    expect(screen.getByTestId('acquisition-runtime')).toHaveTextContent('runtime.notRunning');
  });
});

describe('acquisition settings adapters', () => {
  const record = (extra: Record<string, unknown>) => ({ id: 'D', name: 'P', description: '', protocol: 'modbus_tcp', status: 'draft', connection_config: '{}', last_test_at: null, last_test_success: null, last_test_error: '', created_at: '', updated_at: '', ...extra }) as unknown as StudioV2WorkspaceDeviceRecord;
  it('hydrates saved settings and fails closed on malformed values', () => {
    expect(hydrateStudioV2Device(record({ acquisition_interval_ms: 500, acquisition_auto_interval: true, stale_threshold_ms: 2500, frozen_threshold_ms: 0 })).acquisition)
      .toEqual({ interval_ms: 500, auto_interval: true, stale_threshold_ms: 2500, frozen_threshold_ms: 0 });
    expect(hydrateStudioV2Device(record({ acquisition_auto_interval: false, stale_threshold_ms: null, frozen_threshold_ms: 0 })).acquisition)
      .toEqual({ interval_ms: null, auto_interval: false, stale_threshold_ms: null, frozen_threshold_ms: 0 });
    expect(hydrateStudioV2Device(record({})).acquisition).toBeUndefined();
    expect(hydrateStudioV2Device(record({ acquisition_auto_interval: 'yes', frozen_threshold_ms: 0 })).acquisition).toBeUndefined();
    expect(hydrateStudioV2Device(record({ acquisition_auto_interval: false, frozen_threshold_ms: -1 })).acquisition).toBeUndefined();
  });
  it('confirms a save only when readback carries the submitted settings', () => {
    const settings = { interval_ms: 500, auto_interval: true, stale_threshold_ms: null, frozen_threshold_ms: 0 };
    expect(acquisitionMatches(record({ acquisition_interval_ms: 500, acquisition_auto_interval: true, stale_threshold_ms: null, frozen_threshold_ms: 0 }), settings)).toBe(true);
    expect(acquisitionMatches(record({ acquisition_interval_ms: 1000, acquisition_auto_interval: true, stale_threshold_ms: null, frozen_threshold_ms: 0 }), settings)).toBe(false);
    expect(acquisitionMatches(record({ acquisition_interval_ms: 500, acquisition_auto_interval: false, stale_threshold_ms: null, frozen_threshold_ms: 0 }), settings)).toBe(false);
    expect(acquisitionMatches(record({}), undefined)).toBe(true);
  });
  it('projects only finite positive runtime intervals and reads one device', () => {
    const status = projectRuntimeStatus({ collectors: [{ device_id: 'D', status: 'running', protocol: 'modbus_tcp', availability_status: 'available', running: true, configured_interval_ms: 500, effective_interval_ms: 800 },
      { device_id: 'E', configured_interval_ms: 'x', effective_interval_ms: -1 }] });
    expect(status.collectors[0]).toMatchObject({ configured_interval_ms: 500, effective_interval_ms: 800 });
    expect(status.collectors[1]).not.toHaveProperty('configured_interval_ms');
    expect(runtimeIntervalReadout('D', { data: status, isError: false, isPending: false })).toEqual({ status: 'ready', configured_ms: 500, effective_ms: 800 });
    expect(runtimeIntervalReadout('E', { data: status, isError: false, isPending: false })).toEqual({ status: 'not-running' });
    expect(runtimeIntervalReadout('D', { data: undefined, isError: true, isPending: false })).toEqual({ status: 'unavailable' });
    expect(runtimeIntervalReadout('D', { data: undefined, isError: false, isPending: true })).toEqual({ status: 'loading' });
  });
  it('keeps the connection config untouched for undeclared devices', () => {
    const parsed = parsedDeviceSourceBuffer(source('mc_3e', mc))!;
    expect(toBackendConnectionConfig({ ...parsed, id: 'D', status: 'draft', test: null })).not.toHaveProperty('acquisition');
  });
});
