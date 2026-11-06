import { useState, type ComponentProps } from 'react';
import { act, cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { RecordingHomeMemoryView } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView';
import type { MemoryRuleInput } from '@/features/datalink/workbench-v2/recording-home/recordingHomeMemoryMap';

type ViewProps = ComponentProps<typeof RecordingHomeMemoryView>;
const layouts = ['grid', 'navigator', 'map', 'chapters', 'detail'] as const;
const rule = (id: string, address: string, count = 1): MemoryRuleInput => ({
  id, label: id, start_address: address, count, data_type: 'int16', skipped_addresses: [],
});
// Two separate conflicts catch cancellation that incorrectly returns to the first cell/conflict.
const inputs = [rule('A', '40001', 9), rule('B', '40002'), rule('C', '40008')];
const callbacks = () => ({ onOpen: vi.fn(), onAdd: vi.fn(), onBack: vi.fn(), onClose: vi.fn(), onRetry: vi.fn() });
const base = { protocol: 'modbus_tcp', inputs, status: 'ready', checking: false } as const;

beforeEach(async () => { await i18n.changeLanguage('en'); });
afterEach(cleanup);

function mount(overrides: Partial<ViewProps> = {}) {
  const actions = callbacks();
  const props: ViewProps = { ...base, ...actions, ...overrides };
  const result = render(<RecordingHomeMemoryView {...props} />);
  return { ...actions, rerender: (next: Partial<ViewProps>) => result.rerender(<RecordingHomeMemoryView {...props} {...next} />) };
}
function switchTo(layout: typeof layouts[number]) {
  fireEvent.click(screen.getByTestId(`home-memory-layout-${layout}`));
  expect(screen.getByTestId('home-memory-view')).toHaveAttribute('data-layout', layout);
  expect(screen.getByTestId(`home-memory-layout-${layout}`)).toHaveAttribute('aria-pressed', 'true');
}
function locate(address: string) {
  fireEvent.change(screen.getByRole('textbox', { name: 'Go to address' }), { target: { value: address } });
  fireEvent.click(screen.getByRole('button', { name: 'Go' }));
}
function page(direction: 'Previous' | 'Next') {
  return screen.getByRole('button', { name: direction });
}
function addresses() {
  return screen.getAllByTestId(/^home-memory-cell-/).map(cell => cell.getAttribute('data-testid'));
}
function expectWindow(from: string, to: string) {
  expect(screen.getByRole('group', { name: 'Address blocks' })).toHaveTextContent(`Showing ${from} – ${to}`);
  expect(addresses()).toHaveLength(20);
  expect(addresses()[0]).toBe(`home-memory-cell-${from}`);
  expect(addresses().at(-1)).toBe(`home-memory-cell-${to}`);
}

function DismissalHarness({ actions }: { actions: ReturnType<typeof callbacks> }) {
  const [open, setOpen] = useState(true);
  if (!open) return <button type="button" onClick={() => setOpen(true)}>Reopen memory view</button>;
  return <RecordingHomeMemoryView {...base} {...actions}
    onBack={() => { actions.onBack(); setOpen(false); }}
    onClose={() => { actions.onClose(); setOpen(false); }} />;
}

describe.each(layouts)('%s memory recovery interactions', layout => {
  it('cancels each owner chooser back to the exact conflict cell without opening or adding', () => {
    const { onAdd, onOpen } = mount();
    switchTo(layout);
    for (const address of ['40008', '40002', '40008']) {
      const cell = screen.getByTestId(`home-memory-cell-${address}`);
      act(() => cell.focus());
      fireEvent.click(cell);
      const chooser = screen.getByTestId('home-memory-owners');
      expect(within(chooser).getByRole('heading')).toHaveTextContent(address);
      expect(within(chooser).getByRole('heading')).toHaveFocus();
      fireEvent.click(within(chooser).getByRole('button', { name: 'Cancel' }));
      expect(screen.queryByTestId('home-memory-owners')).not.toBeInTheDocument();
      expect(screen.getByTestId(`home-memory-cell-${address}`)).toHaveFocus();
      expect(cell).toHaveAttribute('tabindex', '0');
      expect(screen.getAllByTestId(/^home-memory-cell-/).filter(node => node.tabIndex === 0)).toEqual([cell]);
    }
    expect(onOpen).not.toHaveBeenCalled();
    expect(onAdd).not.toHaveBeenCalled();
  });

  it('pages both ways and rejects invalid addresses without replacing the current area or window', () => {
    const { onAdd, onOpen } = mount();
    switchTo(layout);
    expectWindow('40001', '40020');
    expect(page('Previous')).toBeDisabled();
    fireEvent.click(page('Next'));
    expectWindow('40021', '40040');
    fireEvent.click(page('Next'));
    expectWindow('40041', '40060');
    fireEvent.click(page('Previous'));
    expectWindow('40021', '40040');
    const savedWindow = addresses();
    for (const invalid of ['30001', 'not-an-address']) {
      locate(invalid);
      expect(screen.getByRole('textbox', { name: 'Go to address' })).toHaveAttribute('aria-invalid', 'true');
      expect(screen.getByRole('alert')).toHaveTextContent('That address is not in the 4x area.');
      expect(screen.getByRole('tab', { name: /^4x/ })).toHaveAttribute('aria-selected', 'true');
      expect(addresses()).toEqual(savedWindow);
    }
    locate('40028');
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    expect(screen.getByRole('textbox', { name: 'Go to address' })).toHaveAttribute('aria-invalid', 'false');
    expect(screen.getByTestId('home-memory-cell-40028')).toHaveAttribute('tabindex', '0');
    expectWindow('40021', '40040');
    fireEvent.click(screen.getByRole('tab', { name: /^3x/ }));
    expectWindow('30001', '30020');
    fireEvent.click(page('Next'));
    expectWindow('30021', '30040');
    fireEvent.click(screen.getByRole('tab', { name: /^4x/ }));
    expectWindow('40021', '40040');
    fireEvent.click(page('Previous'));
    expectWindow('40001', '40020');
    expect(page('Previous')).toBeDisabled();
    expect(onAdd).not.toHaveBeenCalled();
    expect(onOpen).not.toHaveBeenCalled();
  });

  it.each([
    ['recovery', 'A new read range save is still unconfirmed.'],
    ['local', 'Some range changes on this page have not been saved'],
    ['stale', 'The saved ranges changed since this page loaded'],
  ] as const)('blocks creation during %s while still opening a known owner', (blocked, notice) => {
    const { onAdd, onOpen } = mount({ blocked });
    switchTo(layout);
    expect(screen.getByRole('alert')).toHaveTextContent(notice);
    const unknown = screen.getByTestId('home-memory-cell-40010');
    expect(unknown).toHaveAttribute('data-state', 'unconfirmed');
    expect(unknown.tagName).not.toBe('BUTTON');
    act(() => unknown.focus());
    fireEvent.click(unknown);
    fireEvent.keyDown(unknown, { key: 'Enter' });
    fireEvent.keyDown(unknown, { key: ' ' });
    if (layout === 'map' || layout === 'detail') {
      expect(screen.getByTestId(layout === 'map' ? 'home-memory-summary' : 'home-memory-detail')).toHaveTextContent('40010');
      expect(screen.queryByTestId('home-memory-detail-action')).not.toBeInTheDocument();
    }
    fireEvent.click(screen.getByTestId('home-memory-cell-40001'));
    expect(onOpen).toHaveBeenCalledExactlyOnceWith('A');
    fireEvent.click(screen.getByTestId('home-memory-cell-40002'));
    fireEvent.click(within(screen.getByTestId('home-memory-owners')).getByRole('button', { name: /^B ·/ }));
    expect(onOpen.mock.calls).toEqual([['A'], ['B']]);
    expect(onAdd).not.toHaveBeenCalled();
  });

  it('retries a failed range list and keeps cells absent until the successful read arrives', () => {
    const { onRetry, onAdd, onOpen, rerender } = mount({ status: 'error' });
    switchTo(layout);
    expect(screen.getByRole('alert')).toHaveTextContent('Saved ranges could not be read');
    expect(screen.queryByTestId(/^home-memory-cell-/)).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Read ranges again' }));
    expect(onRetry).toHaveBeenCalledTimes(1);
    expect(onOpen).not.toHaveBeenCalled();
    expect(onAdd).not.toHaveBeenCalled();
    rerender({ status: 'loading' });
    expect(screen.getByRole('status')).toHaveTextContent('Reading the current workspace');
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    expect(screen.queryByTestId(/^home-memory-cell-/)).not.toBeInTheDocument();
    rerender({ status: 'ready' });
    expect(screen.getByTestId('home-memory-view')).toHaveAttribute('data-layout', layout);
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    expectWindow('40001', '40020');
    fireEvent.click(screen.getByTestId('home-memory-cell-40001'));
    expect(onOpen).toHaveBeenCalledExactlyOnceWith('A');
    expect(onAdd).not.toHaveBeenCalled();
  });

  it.each(['back', 'close'] as const)('%s calls only its callback and remount clears chooser and invalid address state', exit => {
    const actions = callbacks();
    render(<DismissalHarness actions={actions} />);
    switchTo(layout);
    locate('30001');
    fireEvent.click(screen.getByTestId('home-memory-cell-40008'));
    expect(screen.getByTestId('home-memory-owners')).toBeInTheDocument();
    expect(screen.getByRole('textbox', { name: 'Go to address' })).toHaveAttribute('aria-invalid', 'true');
    fireEvent.click(screen.getByTestId(`home-memory-${exit}`));
    expect(actions.onBack).toHaveBeenCalledTimes(exit === 'back' ? 1 : 0);
    expect(actions.onClose).toHaveBeenCalledTimes(exit === 'close' ? 1 : 0);
    expect(screen.queryByTestId('home-memory-view')).not.toBeInTheDocument();
    expect(screen.queryByTestId('home-memory-owners')).not.toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Reopen memory view' }));
    expect(screen.getByRole('heading', { name: 'Address areas' })).toHaveFocus();
    switchTo(layout);
    expect(screen.queryByTestId('home-memory-owners')).not.toBeInTheDocument();
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    expect(screen.getByRole('textbox', { name: 'Go to address' })).toHaveValue('');
    expect(screen.getByRole('textbox', { name: 'Go to address' })).toHaveAttribute('aria-invalid', 'false');
    expectWindow('40001', '40020');
    fireEvent.click(screen.getByTestId('home-memory-cell-40008'));
    fireEvent.click(within(screen.getByTestId('home-memory-owners')).getByRole('button', { name: 'Cancel' }));
    expect(screen.getByTestId('home-memory-cell-40008')).toHaveFocus();
    expect(actions.onOpen).not.toHaveBeenCalled();
    expect(actions.onAdd).not.toHaveBeenCalled();
  });
});
