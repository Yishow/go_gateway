import { act, cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { RecordingHomeMemoryView } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView';
import type { MemoryRuleInput } from '@/features/datalink/workbench-v2/recording-home/recordingHomeMemoryMap';
beforeEach(async () => { await i18n.changeLanguage('en'); });
afterEach(cleanup);
const rules: MemoryRuleInput[] = [{ id: 'F', label: 'Flow', start_address: '40001', count: 1, data_type: 'float32', skipped_addresses: [] }];
function view(inputs = rules) {
  const onAdd = vi.fn(), onOpen = vi.fn(), setVisible = vi.fn();
  render(<RecordingHomeMemoryView protocol="modbus_tcp" inputs={inputs} status="ready" checking={false}
    values={{ status: 'live', byCell: new Map([['HR:1:float32', { value: '12.5' }]]), setVisible }}
    onAdd={onAdd} onOpen={onOpen} onBack={vi.fn()} onClose={vi.fn()} />);
  return { onAdd, onOpen, setVisible };
}
it.each(['map', 'detail'])('%s keeps keyboard detail when scrolling brings a cell under a stationary pointer', layout => {
  view(); fireEvent.click(screen.getByTestId(`home-memory-layout-${layout}`));
  const first = screen.getByTestId('home-memory-cell-40001'), hovered = screen.getByTestId('home-memory-cell-40008');
  act(() => first.focus());
  const address = () => screen.getByTestId(layout === 'map' ? 'home-memory-summary' : 'home-memory-detail').querySelector('.gw-memory-detail-address');
  // Focus/arrow scrolling can emit enter without the user moving the mouse.
  fireEvent.mouseEnter(hovered);
  expect(address()).toHaveTextContent('40001');
  fireEvent.mouseMove(hovered);
  expect(address()).toHaveTextContent('40008');
  fireEvent.keyDown(first, { key: 'ArrowRight' });
  fireEvent.mouseEnter(screen.getByTestId('home-memory-cell-40014'));
  expect(address()).toHaveTextContent('40002');
  expect(screen.getByTestId('home-memory-cell-40002')).toHaveFocus();
  fireEvent.mouseMove(hovered);
  expect(address()).toHaveTextContent('40008');
});
it.each(['map', 'detail'])('%s previews hovered continuation without moving keyboard focus or invoking an action', layout => {
  const { onAdd, onOpen } = view();
  fireEvent.click(screen.getByTestId(`home-memory-layout-${layout}`));
  const first = screen.getByTestId('home-memory-cell-40001');
  act(() => first.focus());
  fireEvent.mouseMove(screen.getByTestId('home-memory-cell-40002'));
  const detail = screen.getByTestId(layout === 'map' ? 'home-memory-summary' : 'home-memory-detail');
  expect(detail.querySelector('.gw-memory-detail-address')).toHaveTextContent('40002');
  expect(detail.querySelector('.gw-memory-detail-reading')).toHaveTextContent('12.5');
  expect(first).toHaveFocus(); expect(first).toHaveAttribute('tabindex', '0');
  expect(onAdd).not.toHaveBeenCalled(); expect(onOpen).not.toHaveBeenCalled();
  fireEvent.click(screen.getByTestId('home-memory-detail-action'));
  expect(onOpen).toHaveBeenCalledExactlyOnceWith('F');
});
it('keyboard selection supersedes hover and detail navigation follows the displayed address', () => {
  view(); fireEvent.click(screen.getByTestId('home-memory-layout-detail'));
  const first = screen.getByTestId('home-memory-cell-40001'); act(() => first.focus());
  fireEvent.mouseMove(screen.getByTestId('home-memory-cell-40008'));
  fireEvent.keyDown(first, { key: 'ArrowRight' });
  expect(screen.getByTestId('home-memory-cell-40002')).toHaveFocus();
  expect(screen.getByTestId('home-memory-detail').querySelector('.gw-memory-detail-address')).toHaveTextContent('40002');
  fireEvent.mouseMove(screen.getByTestId('home-memory-cell-40008'));
  fireEvent.click(screen.getByTestId('home-memory-detail-next'));
  expect(screen.getByTestId('home-memory-detail').querySelector('.gw-memory-detail-address')).toHaveTextContent('40009');
});
it('hovered unknown cells remain non-actionable and changing area or page discards preview', () => {
  const { onAdd } = view([{ ...rules[0], data_type: 'string' }]);
  fireEvent.click(screen.getByTestId('home-memory-layout-detail'));
  fireEvent.mouseMove(screen.getByTestId('home-memory-cell-40008'));
  expect(screen.getByTestId('home-memory-detail').querySelector('.gw-memory-detail-address')).toHaveTextContent('40008');
  expect(screen.queryByTestId('home-memory-detail-action')).not.toBeInTheDocument();
  fireEvent.click(screen.getByTestId('home-memory-area-IR'));
  expect(screen.getByTestId('home-memory-detail').querySelector('.gw-memory-detail-address')).toHaveTextContent('30001');
  fireEvent.click(screen.getByTestId('home-memory-area-HR'));
  expect(screen.getByTestId('home-memory-detail').querySelector('.gw-memory-detail-address')).toHaveTextContent('40001');
  fireEvent.mouseMove(screen.getByTestId('home-memory-cell-40008'));
  fireEvent.click(screen.getByRole('button', { name: 'Next' }));
  expect(screen.getByTestId('home-memory-detail').querySelector('.gw-memory-detail-address')).toHaveTextContent('40021');
  expect(onAdd).not.toHaveBeenCalled();
});

it('reports a hovered detail reading even when its physical cell is offscreen', () => {
  let notify: (entries: IntersectionObserverEntry[]) => void = () => undefined;
  const observed: Element[] = [];
  vi.stubGlobal('IntersectionObserver', class {
    constructor(callback: typeof notify) { notify = callback; }
    observe(node: Element) { observed.push(node); }
    disconnect() { observed.length = 0; }
  });
  try {
    const { setVisible } = view([...rules, { ...rules[0], id: 'I', start_address: '40003', data_type: 'int16' }]);
    fireEvent.click(screen.getByTestId('home-memory-layout-detail'));
    fireEvent.mouseMove(screen.getByTestId('home-memory-cell-40003'));
    const reading = screen.getByTestId('home-memory-detail').querySelector('.gw-memory-detail-reading')!;
    expect(observed).toContain(reading);
    expect(reading).toHaveAttribute('data-value-key', 'HR:3:int16');
    act(() => notify([{ target: reading, isIntersecting: true } as IntersectionObserverEntry]));
    expect(setVisible).toHaveBeenLastCalledWith(['HR:3:int16']);
  } finally { vi.unstubAllGlobals(); }
});

it('hovering another cell does not steal focus from an open owner chooser', () => {
  view([...rules, { ...rules[0], id: 'Overlap', label: 'Overlap' }]);
  fireEvent.click(screen.getByTestId('home-memory-layout-detail'));
  fireEvent.click(screen.getByTestId('home-memory-cell-40001'));
  const cancel = screen.getByRole('button', { name: 'Cancel' });
  act(() => cancel.focus());
  fireEvent.mouseMove(screen.getByTestId('home-memory-cell-40008'));
  expect(cancel).toHaveFocus();
});
