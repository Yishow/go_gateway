import { cleanup, fireEvent, render, screen } from '@testing-library/react';
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
    // The ordinary reading is not replaced with a rule label; focused inspection may name its owner.
    expect(first.querySelector('.gw-memory-shared')).not.toHaveTextContent('A_');
    expect(first.querySelector('.gw-memory-focus-information')).toHaveTextContent('A_');
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
    expect(observed.map(el => el.getAttribute('data-value-key'))).toEqual(['HR:1:int16', 'HR:1:int16', 'HR:2:int16', 'HR:2:int16']);
    const second = screen.getByTestId('home-memory-cell-40002');
    const disclosure = second.querySelector('.gw-memory-focus-information')!;
    notify([{ target: second, isIntersecting: true }, { target: disclosure, isIntersecting: true }] as IntersectionObserverEntry[], {} as IntersectionObserver);
    expect(setVisible).toHaveBeenLastCalledWith([memoryValueKey('HR', 2, 'int16')]);
    vi.unstubAllGlobals();
  });
});

describe('typed physical register geometry', () => {
  const cell = (address: number) => screen.getByTestId(`home-memory-cell-${address}`);
  function mount(start: number, type: 'int16' | 'float32' | 'int64' | 'float64') {
    render(<RecordingHomeMemoryView protocol="modbus_tcp" status="ready" checking={false} onOpen={noop} onAdd={noop} onBack={noop} onClose={noop}
      inputs={[{ id: 'A', label: 'A_', start_address: `${start}`, count: 1, data_type: type, skipped_addresses: [] }]} />);
  }
  it.each([['int16', 1], ['float32', 2], ['int64', 4], ['float64', 4]] as const)('shares one %s label and value across %i physical registers', (type, width) => {
    mount(40001, type);
    expect(cell(40001).querySelector('[data-shared-span]')).toHaveAttribute('data-shared-span', `${width}`);
    for (let part = 0; part < width; part++) {
      expect(cell(40001 + part).tagName).toBe('BUTTON');
      expect(cell(40001 + part)).toHaveAccessibleName(expect.stringContaining(`register ${part + 1} of ${width}`));
      if (part) { expect(cell(40001 + part).querySelector('.gw-memory-type')).toBeNull(); expect(cell(40001 + part)).not.toHaveAttribute('data-value-key'); }
    }
  });
  it('splits a four-register value at the real row boundary with explicit continuation', () => {
    mount(40009, 'float64');
    expect(cell(40009).querySelector('[data-shared-span]')).toHaveAttribute('data-shared-span', '2');
    expect(cell(40010)).toHaveTextContent('next row');
    expect(cell(40011).querySelector('[data-shared-span]')).toHaveAttribute('data-shared-span', '2');
    expect(cell(40011)).toHaveTextContent('cont. 40009');
    expect(cell(40011)).toHaveTextContent('3/4');
  });
  it('keeps a page continuation addressable without inventing a new value key', () => {
    mount(40249, 'float64');
    // A capped 250-register window splits the value after its second register.
    fireEvent.change(screen.getByRole('textbox'), { target: { value: '40001' } });
    fireEvent.click(screen.getByRole('button', { name: 'Go' }));
    expect(cell(40249).querySelector('[data-shared-span]')).toHaveAttribute('data-shared-span', '2');
    fireEvent.click(screen.getByRole('button', { name: /Next/ }));
    expect(cell(40251)).toHaveTextContent('cont. 40249');
    expect(cell(40251)).toHaveTextContent('3/4');
    expect(cell(40251)).not.toHaveAttribute('data-value-key');
  });
});

it('uses the rendered five-column geometry for joins and keyboard movement', () => {
  const original = window.getComputedStyle;
  const style = vi.spyOn(window, 'getComputedStyle').mockImplementation((node, pseudo) => {
    const result = original(node, pseudo);
    if (node.classList.contains('gw-memory-grid')) Object.defineProperty(result, 'gridTemplateColumns', { value: '80px 80px 80px 80px 80px' });
    return result;
  });
  try {
    render(<RecordingHomeMemoryView protocol="modbus_tcp" status="ready" checking={false} onOpen={noop} onAdd={noop} onBack={noop} onClose={noop}
      inputs={[{ id: 'A', label: 'A_', start_address: '40005', count: 1, data_type: 'float64', skipped_addresses: [] }]} />);
    const first = screen.getByTestId('home-memory-cell-40005');
    expect(first.querySelector('[data-shared-span]')).toHaveAttribute('data-shared-span', '1');
    expect(screen.getByTestId('home-memory-cell-40006').querySelector('[data-shared-span]')).toHaveAttribute('data-shared-span', '3');
    first.focus(); fireEvent.keyDown(first, { key: 'ArrowDown' });
    expect(screen.getByTestId('home-memory-cell-40010')).toHaveFocus();
  } finally { style.mockRestore(); }
});
