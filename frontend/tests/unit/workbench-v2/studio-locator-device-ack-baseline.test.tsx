import { act, renderHook } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { useStudioV2LocatorHydration } from '@/pages/datalink/workbench-v2/studioV2LocatorHydration';
import { INITIAL_STATE, type useWorkbenchV2State } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';
import { toBackendConnectionConfig } from '@/features/datalink/workbench-v2/state/studioV2DeviceAutosave';
import type { Device, WorkbenchV2State } from '@/features/datalink/workbench-v2/state/types';
import { locatorRead, nativeDevice } from './studio-locator-fixtures';

function harness(fatek = false) {
  const entry = locatorRead();
  entry.entry.devices![0].connection_config = JSON.stringify({ host: 'device.example', port: 502, slave_id: 1, timeout: 5 });
  if (fatek) { entry.entry.devices![0].protocol = 'fatek_fbs'; entry.entry.devices![0].connection_config = JSON.stringify({ host: 'device.example', port: 2000, station: 0, mode: 'tcp', timeout: 5 }); }
  const ref = { current: structuredClone(INITIAL_STATE) as WorkbenchV2State };
  const actions: ReturnType<typeof useWorkbenchV2State> = { state: INITIAL_STATE, dispatch: vi.fn(), setView: vi.fn(), setCurrent: vi.fn(),
    completeStep: vi.fn(), toggleSidebar: vi.fn(), toggleSummaryRail: vi.fn(), setSidebarCollapsed: vi.fn(),
    setShowSummaryRail: vi.fn(), resetFlow: vi.fn(), selectRule: vi.fn() };
  const hook = renderHook(({ selected }) => useStudioV2LocatorHydration(actions, ref, false, selected), { initialProps: { selected: entry } });
  const submit = () => ({ ...structuredClone(ref.current.devices[0]), name: 'Saved D' });
  const latest = (value: Device) => { ref.current = { ...ref.current, devices: [{ ...value, save_state: 'saved', save_error: null,
    runtime_apply_status: 'not_running', runtime_apply_message: null, availability_status: 'available', availability_reason: null, running: false }] }; };
  const ack = (value: Device) => ({ ...nativeDevice(value.id), name: value.name, description: value.description, protocol: value.protocol,
    // Backend object-key order is intentionally not the submitted object's order.
    connection_config: JSON.stringify(Object.fromEntries(Object.entries(toBackendConnectionConfig(value)).reverse())), runtime_apply_status: 'not_running' as const });
  return { ...hook, ref, entry, submit, latest, ack };
}

describe('locator own-device ACK baseline boundary', () => {
  it.each([false, true])('advances only the existing accepted device after matching canonical ACK, FATEK alias=%s', fatek => {
    const h = harness(fatek), submitted = h.submit();
    const accept = h.result.current.captureDeviceSave(submitted, 'W');
    expect(accept).toBeTypeOf('function'); h.latest(submitted);
    act(() => accept!(h.ack(submitted)));
    h.rerender({ selected: locatorRead(1, 'B') });
    expect(h.result.current.ready).toBe(true); expect(h.ref.current.devices[0].id).toBe('B');
  });
  it('keeps a pending save blocked without an ACK even if local flags claim saved', () => {
    const h = harness(), submitted = h.submit();
    h.result.current.captureDeviceSave(submitted, 'W'); h.latest(submitted);
    h.rerender({ selected: locatorRead(1, 'B') });
    expect(h.result.current.blocked).toBe(true); expect(h.ref.current.devices[0].name).toBe('Saved D');
  });
  it.each(['foreign-ack', 'next-edit', 'other-rule', 'settings'] as const)('preserves dirty refusal for %s', reason => {
    const h = harness(), submitted = h.submit(), accept = h.result.current.captureDeviceSave(submitted, 'W');
    h.latest(reason === 'next-edit' ? { ...submitted, name: 'Newer D edit' } : submitted);
    if (reason === 'other-rule') h.ref.current = { ...h.ref.current, selectedRuleId: 'unsaved', rules: [{ id: 'unsaved', device_id: 'D', name: 'Unsaved range',
      start_address: '40001', count: 2, data_type: 'int16', naming_prefix: 'LOCAL_', enabled: false, scale_multiplier: 1, scale_offset: 0,
      data_format: '', skipped_addresses: [], share_enabled: false, share_start_register: null, share_stride: null }] };
    if (reason === 'settings') h.ref.current = { ...h.ref.current, settings: { ...h.ref.current.settings, modbus_share: { ...h.ref.current.settings.modbus_share, enabled: true } } };
    act(() => accept!({ ...h.ack(submitted), ...(reason === 'foreign-ack' ? { id: 'B' } : {}) }));
    h.rerender({ selected: locatorRead(1, 'B') });
    expect(h.result.current.blocked).toBe(true); expect(h.ref.current.devices[0].id).toBe('D');
  });
  it('rejects an old accepted-owner ACK across D to B to D ABA', () => {
    const h = harness(), submitted = h.submit(), accept = h.result.current.captureDeviceSave(submitted, 'W');
    h.rerender({ selected: locatorRead(1, 'B') }); expect(h.result.current.ready).toBe(true);
    h.rerender({ selected: structuredClone(h.entry) }); expect(h.result.current.ready).toBe(true);
    h.latest(submitted); act(() => accept!(h.ack(submitted)));
    h.rerender({ selected: locatorRead(1, 'B') }); expect(h.result.current.blocked).toBe(true);
  });
  it('rejects an ACK after a same-workspace foreign viewing intent even while the old accepted owner is retained', () => {
    const h = harness(), submitted = h.submit(), accept = h.result.current.captureDeviceSave(submitted, 'W');
    h.latest(submitted);
    h.rerender({ selected: locatorRead(1, 'B') }); expect(h.result.current.blocked).toBe(true);
    act(() => accept!(h.ack(submitted)));
    h.rerender({ selected: locatorRead(1, 'B') }); expect(h.result.current.blocked).toBe(true);
  });
  it('does not capture a newly created device or a foreign workspace', () => {
    const h = harness();
    expect(h.result.current.captureDeviceSave({ ...h.submit(), persisted: false }, 'W')).toBeUndefined();
    expect(h.result.current.captureDeviceSave(h.submit(), 'OTHER')).toBeUndefined();
    expect(h.result.current.captureDeviceSave({ ...h.submit(), id: 'NEW' }, 'W')).toBeUndefined();
  });
});
