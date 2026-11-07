import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { act, cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { RecordingHomeMemoryView } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView';
import { useHomeMemoryValues } from '@/features/datalink/workbench-v2/recording-home/recordingHomeMemoryValues';
import type { MemoryRuleInput } from '@/features/datalink/workbench-v2/recording-home/recordingHomeMemoryMap';
import { pointAPI, runtimeAPI } from '@/services/datalink';
import type { RuntimeValueEvent } from '@/types/datalink';
import type { Point } from '@/types/point';

vi.mock('@/services/datalink', () => ({
  pointAPI: { list: vi.fn(), create: vi.fn(), update: vi.fn(), delete: vi.fn(), pollNow: vi.fn(), batchPoll: vi.fn() },
  runtimeAPI: { getStreamUrl: vi.fn((device: string, points: string[]) => `stream://${device}?points=${points.join(',')}`) },
}));

// Integration fixtures only: neither the point list nor events come from a live PLC/HSL server.
class MockEventSource {
  static instances: MockEventSource[] = [];
  listeners = new Map<string, (event: Event) => void>();
  onopen: (() => void) | null = null;
  onerror: (() => void) | null = null;
  close = vi.fn();
  constructor(readonly url: string) { MockEventSource.instances.push(this); }
  addEventListener(type: string, listener: (event: Event) => void) { this.listeners.set(type, listener); }
  removeEventListener(type: string) { this.listeners.delete(type); }
  emit(payload: RuntimeValueEvent) { this.listeners.get('value')?.(new MessageEvent('value', { data: JSON.stringify(payload) })); }
}

class MockIntersectionObserver {
  static instances: MockIntersectionObserver[] = [];
  targets: Element[] = [];
  disconnected = false;
  constructor(private notify: (entries: IntersectionObserverEntry[], observer: IntersectionObserver) => void) { MockIntersectionObserver.instances.push(this); }
  observe(target: Element) { this.targets.push(target); }
  disconnect() { this.disconnected = true; }
  showAll() {
    this.setVisibility(this.targets.map(target => ({ target, visible: true })));
  }
  setVisibility(entries: { target: Element; visible: boolean }[]) {
    this.notify(entries.map(({ target, visible }) => ({ target, isIntersecting: visible }) as IntersectionObserverEntry), this as unknown as IntersectionObserver);
  }
}

const point: Point = {
  id: 'persisted-01', device_id: 'device-01', name: 'SENSOR_01', description: '', address: '40001', data_type: 'int16',
  enabled: true, polling_group_id: '', last_value: 0, last_read_at: '2026-11-04T10:00:00Z', last_error: '',
  error_count: 0, created_at: '', updated_at: '',
};
const inputs: MemoryRuleInput[] = [{ id: 'range-01', label: 'Sensor', start_address: '40001', count: 1, data_type: 'int16', skipped_addresses: [] }];
const callbacks = { onOpen: vi.fn(), onAdd: vi.fn(), onBack: vi.fn(), onClose: vi.fn() };
const layouts = ['navigator', 'map', 'chapters', 'detail', 'grid'] as const;
let client: QueryClient;

// Match the real owner: the hook lives above one view instance, not inside each layout.
function HomeHarness({ rules = inputs }: { rules?: MemoryRuleInput[] }) {
  const values = useHomeMemoryValues('device-01', 'modbus_tcp', true);
  return <RecordingHomeMemoryView protocol="modbus_tcp" inputs={rules} status="ready" checking={false} values={values} {...callbacks} />;
}

beforeEach(async () => {
  await i18n.changeLanguage('en');
  vi.useFakeTimers();
  vi.clearAllMocks();
  MockEventSource.instances = [];
  MockIntersectionObserver.instances = [];
  vi.stubGlobal('EventSource', MockEventSource);
  vi.stubGlobal('IntersectionObserver', MockIntersectionObserver);
  vi.stubGlobal('ResizeObserver', class { observe() {} disconnect() {} });
  const getComputedStyle = window.getComputedStyle.bind(window);
  vi.spyOn(window, 'getComputedStyle').mockImplementation((node, pseudo) => {
    const style = getComputedStyle(node, pseudo);
    if (node.classList.contains('gw-memory-grid')) Object.defineProperty(style, 'gridTemplateColumns', { value: '80px '.repeat(10).trim() });
    return style;
  });
  vi.mocked(pointAPI.list).mockResolvedValue([point]);
  client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
});

afterEach(() => {
  cleanup();
  client.clear();
  vi.clearAllTimers();
  vi.useRealTimers();
  vi.restoreAllMocks();
  vi.unstubAllGlobals();
});

async function advance(ms: number) { await act(async () => { await vi.advanceTimersByTimeAsync(ms); }); }
async function showVisible() {
  act(() => { MockIntersectionObserver.instances.filter(observer => !observer.disconnected).forEach(observer => observer.showAll()); });
  await advance(300);
}
async function mount(rules?: MemoryRuleInput[]) {
  render(<QueryClientProvider client={client}><HomeHarness rules={rules} /></QueryClientProvider>);
  await advance(1);
  await showVisible();
  expect(MockEventSource.instances).toHaveLength(1);
  const stream = MockEventSource.instances[0];
  act(() => stream.onopen?.());
  return stream;
}
function reading() { return screen.getByTestId('home-memory-cell-40001').querySelector('[data-reading]'); }
function expectNoMutations() {
  for (const callback of Object.values(callbacks)) expect(callback).not.toHaveBeenCalled();
  for (const mutation of [pointAPI.create, pointAPI.update, pointAPI.delete, pointAPI.pollNow, pointAPI.batchPoll]) expect(mutation).not.toHaveBeenCalled();
  expect(client.getMutationCache().getAll()).toHaveLength(0);
}
function expectSubscription(stream: MockEventSource) {
  expect(MockEventSource.instances).toEqual([stream]);
  expect(stream.close).not.toHaveBeenCalled();
  expect(runtimeAPI.getStreamUrl).toHaveBeenCalledExactlyOnceWith('device-01', ['persisted-01']);
  expectNoMutations();
}
async function switchTo(layout: typeof layouts[number]) {
  fireEvent.click(screen.getByTestId(`home-memory-layout-${layout}`));
  expect(screen.getByTestId('home-memory-view')).toHaveAttribute('data-layout', layout);
  await showVisible();
}
function event(value: number, extra: Partial<RuntimeValueEvent> = {}): RuntimeValueEvent {
  return { device_id: 'device-01', point_id: 'persisted-01', address: '40001', raw_value: value, transformed_value: value,
    timestamp: '2026-11-04T10:00:01Z', quality: 'good', stale: false, ...extra };
}

async function showDetailOnly(detailFirst = false) {
  act(() => {
    for (const observer of MockIntersectionObserver.instances.filter(item => !item.disconnected)) {
      const entries = observer.targets.map(target => ({ target, visible: Boolean(target.closest('[data-memory-inspection]')) }));
      // Deliver duplicate-key visibility in both orders, as a browser observer may do.
      entries.sort((a, b) => detailFirst ? Number(b.visible) - Number(a.visible) : Number(a.visible) - Number(b.visible));
      observer.setVisibility(entries);
    }
  });
  await advance(300);
}

describe('memory layout and live hook integration continuity', () => {
  it.each(['map', 'detail'] as const)('keeps the %s detail subscribed when the grid value is offscreen', async layout => {
    const stream = await mount();
    await switchTo(layout);
    await showDetailOnly();
    expectSubscription(stream);
    act(() => stream.emit(event(58)));
    await advance(1000);
    expect(screen.getByTestId(layout === 'map' ? 'home-memory-summary' : 'home-memory-detail').querySelector('.gw-memory-detail-reading')).toHaveTextContent('58');
  });

  it.each(['map', 'detail'] as const)('does not drop the shared %s value key when detail-visible precedes grid-hidden', async layout => {
    const stream = await mount();
    await switchTo(layout);
    await showDetailOnly(true);
    expectSubscription(stream);
    act(() => stream.emit(event(61)));
    await advance(1000);
    expect(screen.getByTestId(layout === 'map' ? 'home-memory-summary' : 'home-memory-detail').querySelector('.gw-memory-detail-reading')).toHaveTextContent('61');
  });

  it.each(['map', 'detail'] as const)('%s observes the value-start key of a selected continuation even when its first register is outside the page', async layout => {
    vi.mocked(pointAPI.list).mockResolvedValue([{ ...point, address: '40249', data_type: 'int64' }]);
    const stream = await mount([{ ...inputs[0], start_address: '40249', data_type: 'int64' }]);
    await switchTo(layout);
    fireEvent.change(screen.getByRole('textbox'), { target: { value: '40001' } });
    fireEvent.click(screen.getByRole('button', { name: 'Go' }));
    fireEvent.click(screen.getByRole('button', { name: 'Next' }));
    expect(screen.queryByTestId('home-memory-cell-40249')).not.toBeInTheDocument();
    expect(screen.getByTestId('home-memory-cell-40251')).toHaveAttribute('data-part', 'rest');
    await showDetailOnly(true);
    expectSubscription(stream);
    const observedDetail = MockIntersectionObserver.instances.filter(observer => !observer.disconnected)
      .flatMap(observer => observer.targets).find(target => target.closest('[data-memory-inspection]'));
    expect(observedDetail).toHaveAttribute('data-value-key', 'HR:249:int64');
    act(() => stream.emit(event(84, { address: '40249' })));
    await advance(1000);
    expect(screen.getByTestId(layout === 'map' ? 'home-memory-summary' : 'home-memory-detail').querySelector('.gw-memory-detail-reading')).toHaveTextContent('84');
  });

  it('preserves zero, pushed values, stale and error state across all five layouts without resubscribing or mutating', async () => {
    const stream = await mount();
    const states = [
      { payload: undefined, state: 'value', text: '0' },
      { payload: event(37), state: 'value', text: '37' },
      { payload: event(37, { stale: true, quality: 'stale' }), state: 'stale', text: '37' },
      { payload: event(99, { quality: 'bad' }), state: 'error', text: '!' },
      { payload: event(42), state: 'value', text: '42' },
    ];
    for (const { payload, state, text } of states) {
      if (payload) { act(() => stream.emit(payload)); await advance(1000); }
      for (const layout of layouts) {
        expect(reading()).toHaveAttribute('data-reading', state);
        expect(reading()?.textContent).toBe(text);
        await switchTo(layout);
        expect(reading()).toHaveAttribute('data-reading', state);
        expect(reading()?.textContent).toBe(text);
        expect(screen.getByTestId('home-memory-values-status')).toHaveAttribute('data-status', 'live');
        expectSubscription(stream);
        expect(pointAPI.list).toHaveBeenCalledExactlyOnceWith({ device_id: 'device-01' });
      }
    }
  });

  it('flushes a pending reading after a layout change and keeps it when a later refresh fails', async () => {
    const stream = await mount();
    act(() => stream.emit(event(73)));
    await switchTo('navigator'); // 300 ms: the one-second SSE batch is still pending.
    expect(reading()?.textContent).toBe('0');
    expectSubscription(stream);
    await advance(700);
    expect(reading()?.textContent).toBe('73');
    vi.mocked(pointAPI.list).mockRejectedValue(new Error('fixture refresh unavailable'));
    await act(async () => { await client.refetchQueries({ queryKey: ['recording-home-memory-values', 'device-01'] }); });
    await advance(1);
    for (const layout of layouts) {
      await switchTo(layout);
      expect(reading()).toHaveAttribute('data-reading', 'value');
      expect(reading()?.textContent).toBe('73');
      expect(screen.getByTestId('home-memory-values-status')).toHaveAttribute('role', 'alert');
      expect(screen.getByTestId('home-memory-values-status')).toHaveAttribute('data-status', 'error');
      expectSubscription(stream);
    }
    expect(pointAPI.list).toHaveBeenCalledTimes(2); // Initial read and explicitly triggered failed refresh only.
    cleanup();
    expect(stream.close).toHaveBeenCalledTimes(1);
    expect(stream.listeners.size).toBe(0);
  });
});

it.each(['grid', 'navigator', 'chapters'] as const)('%s keeps a focused off-page continuation subscribed', async layout => {
  vi.mocked(pointAPI.list).mockResolvedValue([{ ...point, address: '40249', data_type: 'int64' }]);
  const stream = await mount([{ ...inputs[0], start_address: '40249', data_type: 'int64' }]);
  await switchTo(layout);
  fireEvent.change(screen.getByRole('textbox'), { target: { value: '40001' } });
  fireEvent.click(screen.getByRole('button', { name: 'Go' }));
  fireEvent.click(screen.getByRole('button', { name: 'Next' }));
  const continuation = screen.getByTestId('home-memory-cell-40251');
  act(() => continuation.focus());
  expect(screen.queryByTestId('home-memory-cell-40249')).not.toBeInTheDocument();
  const information = screen.getByTestId('home-memory-inspection')!;
  act(() => {
    for (const observer of MockIntersectionObserver.instances.filter(item => !item.disconnected)) {
      observer.setVisibility(observer.targets.map(target => ({ target, visible: target === information })));
    }
  });
  await advance(300);
  expect(information).toHaveAttribute('data-value-key', 'HR:249:int64');
  expectSubscription(stream);
  act(() => stream.emit(event(84, { address: '40249' })));
  await advance(1000);
  expect(screen.getByTestId('home-memory-inspection')).toHaveTextContent('84');
});
