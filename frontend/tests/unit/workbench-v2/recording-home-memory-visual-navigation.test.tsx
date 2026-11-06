import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { RecordingHomeMemoryView } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView';
import type { MemoryRuleInput } from '@/features/datalink/workbench-v2/recording-home/recordingHomeMemoryMap';
beforeEach(async () => { await i18n.changeLanguage('en'); });
afterEach(cleanup);
const rules: MemoryRuleInput[] = [
  { id: 'A', label: 'Pressure', start_address: '40001', count: 2, data_type: 'int16', skipped_addresses: [] },
  { id: 'B', label: 'Remote flow', start_address: '40301', count: 2, data_type: 'float32', skipped_addresses: [] },
];
function mount(inputs = rules) {
  render(<RecordingHomeMemoryView protocol="modbus_tcp" inputs={inputs} status="ready" checking={false}
    onOpen={vi.fn()} onAdd={vi.fn()} onBack={vi.fn()} onClose={vi.fn()} />);
  fireEvent.click(screen.getByTestId('home-memory-layout-navigator'));
}
it('distinguishes an explicit selected owner from current-window membership', () => {
  mount();
  const a = screen.getByTestId('home-memory-range-jump-A');
  const b = screen.getByTestId('home-memory-range-jump-B');
  expect(a).not.toHaveAttribute('aria-current');
  expect(a.closest('li')).toHaveAttribute('data-in-window', 'true');
  expect(b.closest('li')).not.toHaveAttribute('data-in-window', 'true');
  expect(a.closest('li')).toHaveTextContent('In this window');
  fireEvent.focus(screen.getByTestId('home-memory-cell-40002'));
  expect(a).toHaveAttribute('aria-current', 'true');
  expect(a).toHaveTextContent('Selected range');
  fireEvent.mouseMove(screen.getByTestId('home-memory-cell-40003'));
  expect(a).toHaveAttribute('aria-current', 'true');
  fireEvent.click(b);
  expect(b).toHaveAttribute('aria-current', 'true');
  expect(a).not.toHaveAttribute('aria-current');
  expect(b.closest('li')).toHaveAttribute('data-in-window', 'true');
  expect(a.closest('li')).not.toHaveAttribute('data-in-window', 'true');
  expect(document.querySelectorAll('.gw-memory-grid')).toHaveLength(1);
});
it('does not choose an arbitrary owner of a conflict or unknown slot', () => {
  mount([rules[0], { ...rules[0], id: 'Overlap' }]);
  fireEvent.focus(screen.getByTestId('home-memory-cell-40001'));
  expect(screen.getByTestId('home-memory-range-nav').querySelector('[aria-current]')).toBeNull();
});
it('keeps every saved range in the bounded navigation list with exact type and address', () => {
  mount(Array.from({ length: 40 }, (_, i) => ({ ...rules[0], id: `R${i}`, label: `Range ${i}`, start_address: String(40001 + i * 10) })));
  const nav = screen.getByTestId('home-memory-range-nav');
  expect(within(nav).getAllByRole('button')).toHaveLength(40);
  expect(screen.getByTestId('home-memory-range-jump-R39')).toHaveTextContent('40391');
  expect(nav).toHaveTextContent('int16');
});
it('does not transfer explicit selection to another area with the same numeric offset', () => {
  mount([rules[0], { ...rules[0], id: 'IR', start_address: '30001' }]);
  fireEvent.focus(screen.getByTestId('home-memory-cell-40001'));
  expect(screen.getByTestId('home-memory-range-jump-A')).toHaveAttribute('aria-current', 'true');
  fireEvent.click(screen.getByTestId('home-memory-area-IR'));
  expect(screen.getByTestId('home-memory-range-jump-IR')).not.toHaveAttribute('aria-current');
  fireEvent.focus(screen.getByTestId('home-memory-cell-30001'));
  expect(screen.getByTestId('home-memory-range-jump-IR')).toHaveAttribute('aria-current', 'true');
});
