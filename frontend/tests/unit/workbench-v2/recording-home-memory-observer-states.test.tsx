import { act, cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { RecordingHomeMemoryView } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView';

class Observer {
  static instances: Observer[] = [];
  disconnect = vi.fn();
  observe = vi.fn();
  constructor(readonly notify: (entries: IntersectionObserverEntry[]) => void) { Observer.instances.push(this); }
  show(target: Element) { this.notify([{ target, isIntersecting: true } as IntersectionObserverEntry]); }
}
beforeEach(async () => {
  await i18n.changeLanguage('en'); Observer.instances = [];
  vi.stubGlobal('IntersectionObserver', Observer);
});
afterEach(() => { cleanup(); vi.unstubAllGlobals(); });
const inputs = [{ id: 'A', label: 'A', start_address: '40001', count: 2, data_type: 'int16' as const, skipped_addresses: [] }];
function mount() {
  const setVisible = vi.fn();
  const props = { protocol: 'modbus_tcp' as const, inputs, checking: false, status: 'ready' as const,
    values: { status: 'live' as const, byCell: new Map(), setVisible },
    onOpen: vi.fn(), onAdd: vi.fn(), onBack: vi.fn(), onClose: vi.fn() };
  return { ...render(<RecordingHomeMemoryView {...props} />), props, setVisible };
}
it.each(['loading', 'error', 'checking'] as const)('clears visible subscription when ready becomes %s and ignores old queued observations', state => {
  const { props, setVisible, rerender } = mount();
  const first = screen.getByTestId('home-memory-cell-40001');
  const observer = Observer.instances.at(-1)!;
  act(() => observer.show(first));
  expect(setVisible).toHaveBeenLastCalledWith(['HR:1:int16']);
  rerender(<RecordingHomeMemoryView {...props} status={state === 'checking' ? 'ready' : state} checking={state === 'checking'} />);
  expect(screen.queryByTestId('home-memory-cell-40001')).not.toBeInTheDocument();
  expect(observer.disconnect).toHaveBeenCalledOnce();
  expect(setVisible).toHaveBeenLastCalledWith([]);
  act(() => observer.show(first));
  expect(setVisible).toHaveBeenLastCalledWith([]);
  rerender(<RecordingHomeMemoryView {...props} />);
  act(() => Observer.instances.at(-1)!.show(screen.getByTestId('home-memory-cell-40002')));
  expect(setVisible).toHaveBeenLastCalledWith(['HR:2:int16']);
});
it('does not transiently clear subscribed keys while hover rebinds a populated view', () => {
  const { setVisible } = mount();
  fireEvent.click(screen.getByTestId('home-memory-layout-detail'));
  act(() => Observer.instances.at(-1)!.show(screen.getByTestId('home-memory-cell-40001')));
  setVisible.mockClear();
  fireEvent.mouseEnter(screen.getByTestId('home-memory-cell-40002'));
  expect(setVisible).not.toHaveBeenCalledWith([]);
  const detail = screen.getByTestId('home-memory-detail').querySelector('[data-value-key]')!;
  act(() => Observer.instances.at(-1)!.show(detail));
  expect(setVisible).toHaveBeenLastCalledWith(['HR:2:int16']);
});
