import { cleanup, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { RecordingHomeMemoryView } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView';
import { applyLiveEvent, indexMemoryValues, memoryValueKey, type MemoryValues } from '@/features/datalink/workbench-v2/recording-home/recordingHomeMemoryValues';
import type { Point } from '@/types/point';

const point = (address: string, data_type: Point['data_type'], last_value: unknown, extra: Partial<Point> = {}): Point => ({
  id: `p-${address}`, device_id: 'D', name: `TAG_${address}`, description: '', data_type, address, enabled: true, polling_group_id: '',
  last_value, last_read_at: '2026-11-04T22:05:51+08:00', last_error: '', error_count: 0, created_at: '', updated_at: '', ...extra,
});
const noop = () => undefined;

beforeEach(async () => { await i18n.changeLanguage('en'); });
afterEach(cleanup);

describe('memory view device values', () => {
  it('indexes readings by area, offset and type, skipping other devices and unparsable addresses', () => {
    const index = indexMemoryValues([
      point('D0', 'int16', '12'), point('D1', 'float32', 0.30000000000000004), point('D2', 'int16', '9', { device_id: 'other' }),
      point('??', 'int16', '1'), point('D3', 'int16', '', { last_error: 'timeout' }),
    ], 'D', 'mc_3e');
    expect(index.get(memoryValueKey('D', 0, 'int16'))?.value).toBe('12');
    expect(index.get(memoryValueKey('D', 0, 'float32'))).toBeUndefined();
    expect(index.get(memoryValueKey('D', 1, 'float32'))?.value).toBe('0.3');
    expect(index.has(memoryValueKey('D', 2, 'int16'))).toBe(false);
    expect(index.get(memoryValueKey('D', 3, 'int16'))).toMatchObject({ value: undefined, error: 'timeout' });
  });

  it('shows the device value instead of the point name, and marks missing or failed reads', () => {
    const values: MemoryValues = { status: 'live', updatedAt: Date.now(), byCell: indexMemoryValues([
      point('40001', 'int16', '1234'), point('40003', 'int16', null, { last_error: 'read failed' }),
    ], 'D', 'modbus_tcp') };
    render(<RecordingHomeMemoryView protocol="modbus_tcp" status="ready" checking={false} values={values} onOpen={vi.fn()} onAdd={noop} onBack={noop} onClose={noop}
      inputs={[{ id: 'A', label: 'A_', start_address: '40001', count: 4, data_type: 'int16', skipped_addresses: [] }]} />);
    const first = screen.getByTestId('home-memory-cell-40001');
    expect(first).toHaveTextContent('1234');
    expect(first).not.toHaveTextContent('A_');
    expect(first).toHaveAttribute('aria-label', expect.stringContaining('value 1234'));
    expect(screen.getByTestId('home-memory-cell-40002').querySelector('[data-reading]')).toHaveAttribute('data-reading', 'none');
    expect(screen.getByTestId('home-memory-cell-40003').querySelector('[data-reading]')).toHaveAttribute('data-reading', 'error');
    expect(screen.getByTestId('home-memory-values-status')).toHaveTextContent(/Live · every second · updated .* on-screen cells only/);
  });

  it('reports a failed refresh without hiding the grid', () => {
    const values: MemoryValues = { status: 'error', byCell: new Map() };
    render(<RecordingHomeMemoryView protocol="modbus_tcp" status="ready" checking={false} values={values} onOpen={vi.fn()} onAdd={noop} onBack={noop} onClose={noop}
      inputs={[{ id: 'A', label: 'A_', start_address: '40001', count: 1, data_type: 'int16', skipped_addresses: [] }]} />);
    expect(screen.getByTestId('home-memory-values-status')).toHaveAttribute('role', 'alert');
    expect(screen.getByTestId('home-memory-cell-40001')).toBeInTheDocument();
  });

  it('applies pushed readings: good replaces, bad keeps the last value as failed, stale is marked', () => {
    const base = { pointId: 'p', value: '1', readAt: 't0' };
    const event = { device_id: 'D', point_id: 'p', address: 'D0', raw_value: 7, transformed_value: 70, quality: 'good' as const, stale: false, timestamp: 't1' };
    expect(applyLiveEvent(base, event)).toMatchObject({ value: '70', error: undefined, readAt: 't1', stale: false });
    expect(applyLiveEvent(base, { ...event, quality: 'bad' })).toMatchObject({ value: '1', error: 'bad' });
    expect(applyLiveEvent(base, { ...event, stale: true })).toMatchObject({ value: '70', stale: true });
  });

  it('reports only on-screen value cells to the stream', () => {
    const observed: Element[] = [];
    let notify: (entries: IntersectionObserverEntry[], observer: IntersectionObserver) => void = () => undefined;
    vi.stubGlobal('IntersectionObserver', class { constructor(cb: typeof notify) { notify = cb; } observe(el: Element) { observed.push(el); } disconnect() {} });
    const setVisible = vi.fn();
    render(<RecordingHomeMemoryView protocol="modbus_tcp" status="ready" checking={false} values={{ status: 'live', byCell: new Map(), setVisible }} onOpen={vi.fn()} onAdd={noop} onBack={noop} onClose={noop}
      inputs={[{ id: 'A', label: 'A_', start_address: '40001', count: 2, data_type: 'int16', skipped_addresses: [] }]} />);
    expect(observed.map(el => el.getAttribute('data-value-key'))).toEqual([memoryValueKey('HR', 1, 'int16'), memoryValueKey('HR', 2, 'int16')]);
    notify([{ target: observed[1], isIntersecting: true } as unknown as IntersectionObserverEntry], {} as IntersectionObserver);
    expect(setVisible).toHaveBeenLastCalledWith([memoryValueKey('HR', 2, 'int16')]);
    vi.unstubAllGlobals();
  });
});
