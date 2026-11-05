import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { RecordingHomeMemoryView, type HomeMemoryBlock } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView';
import type { MemoryRuleInput } from '@/features/datalink/workbench-v2/recording-home/recordingHomeMemoryMap';

const layouts = ['Grid focus', 'Range navigation', 'Horizontal map', 'Range chapters', 'Grid with detail'];
const keys = ['grid', 'navigator', 'map', 'chapters', 'detail'];
const rule = (id: string, start: string, count = 1): MemoryRuleInput => ({ id, label: id, start_address: start, count, data_type: 'int16', skipped_addresses: [] });
beforeEach(async () => { await i18n.changeLanguage('en'); });
afterEach(() => { cleanup(); vi.unstubAllGlobals(); });
function view(inputs: MemoryRuleInput[], blocked?: HomeMemoryBlock) {
  const onAdd = vi.fn(), onOpen = vi.fn(), setVisible = vi.fn();
  render(<RecordingHomeMemoryView protocol="modbus_tcp" status="ready" checking={false} blocked={blocked}
    inputs={inputs} onAdd={onAdd} onOpen={onOpen} onBack={vi.fn()} onClose={vi.fn()}
    values={{ status: 'live', byCell: new Map([['IR:201:int16', { value: '42' }]]), setVisible }} />);
  return { onAdd, onOpen, setVisible };
}
function switchTo(index: number) {
  fireEvent.click(screen.getByRole('button', { name: layouts[index] }));
  expect(screen.getByTestId('home-memory-view')).toHaveAttribute('data-layout', keys[index]);
  expect(screen.getByRole('button', { name: layouts[index] })).toHaveAttribute('aria-pressed', 'true');
}

it('keeps the selected area, located window, focused address and live value across every layout', () => {
  const { onAdd, onOpen } = view([rule('A', '40001'), rule('B', '30201')]);
  expect(screen.getByRole('button', { name: layouts[0] })).toHaveAttribute('aria-pressed', 'true');
  fireEvent.click(screen.getByRole('tab', { name: /^3x/ }));
  fireEvent.change(screen.getByRole('textbox'), { target: { value: '30201' } });
  fireEvent.click(screen.getByRole('button', { name: 'Go' }));
  for (const index of [1, 2, 3, 4, 0]) {
    switchTo(index);
    expect(screen.getByRole('tab', { name: /^3x/ })).toHaveAttribute('aria-selected', 'true');
    expect(screen.getByRole('textbox')).toHaveValue('30201');
    expect(screen.getByTestId('home-memory-cell-30201')).toHaveAttribute('tabindex', '0');
    expect(screen.getByTestId('home-memory-cell-30201')).toHaveTextContent('42');
    expect(screen.queryByTestId('home-memory-cell-30001')).not.toBeInTheDocument();
  }
  expect(onAdd).not.toHaveBeenCalled(); expect(onOpen).not.toHaveBeenCalled();
});

it.each(layouts.map((layout, index) => ({ layout, index })))('$layout preserves empty, used and conflict actions and retains owner selection when switching', ({ index }) => {
  const { onAdd, onOpen } = view([rule('A', '40001', 2), rule('B', '40002')]);
  switchTo(index);
  fireEvent.click(screen.getByTestId('home-memory-cell-40003'));
  expect(onAdd).toHaveBeenCalledExactlyOnceWith('40003');
  fireEvent.click(screen.getByTestId('home-memory-cell-40001'));
  expect(onOpen).toHaveBeenCalledExactlyOnceWith('A'); onOpen.mockClear();
  fireEvent.click(screen.getByTestId('home-memory-cell-40002'));
  expect(onOpen).not.toHaveBeenCalled();
  switchTo((index + 1) % layouts.length);
  const owners = screen.getByTestId('home-memory-owners');
  expect(owners).toHaveTextContent('40002');
  fireEvent.click(within(owners).getByRole('button', { name: /^B/ }));
  expect(onOpen).toHaveBeenCalledExactlyOnceWith('B');
});

it.each([undefined, 'recovery', 'local', 'stale'] as const)('never enables uncertain cells across layouts with block %s', blocked => {
  const inputs: MemoryRuleInput[] = blocked ? [rule('A', '40001')] : [{ ...rule('S', '40001'), data_type: 'string' }];
  const { onAdd } = view(inputs, blocked);
  for (const index of [1, 2, 3, 4, 0]) {
    switchTo(index);
    const unknown = screen.getByTestId('home-memory-cell-40003');
    expect(unknown).toHaveAttribute('data-state', 'unconfirmed');
    expect(unknown.tagName).not.toBe('BUTTON'); fireEvent.click(unknown);
    expect(screen.getByRole('alert')).toBeInTheDocument();
  }
  expect(onAdd).not.toHaveBeenCalled();
});

it('reattaches column measurement and visible-value observers after each layout switch', () => {
  const resize = vi.fn(), intersect = vi.fn(), disconnectResize = vi.fn(), disconnectIntersect = vi.fn();
  vi.stubGlobal('ResizeObserver', class { observe = resize; disconnect = disconnectResize; });
  vi.stubGlobal('IntersectionObserver', class { observe = intersect; disconnect = disconnectIntersect; });
  view([rule('A', '40001')]);
  for (const index of [1, 2, 3, 4, 0]) {
    resize.mockClear(); intersect.mockClear(); disconnectResize.mockClear(); disconnectIntersect.mockClear();
    switchTo(index);
    expect(disconnectResize).toHaveBeenCalled(); expect(resize).toHaveBeenCalled();
    expect(disconnectIntersect).toHaveBeenCalled();
    expect(intersect).toHaveBeenCalledWith(screen.getByTestId('home-memory-cell-40001'));
  }
});

it('provides all layout controls in Traditional Chinese', async () => {
  await i18n.changeLanguage('zh-TW'); view([]);
  for (const label of ['格子專注', '範圍導航', '水平地圖', '範圍章節', '格子＋詳情']) {
    expect(screen.getByRole('button', { name: label })).toBeInTheDocument();
  }
});
