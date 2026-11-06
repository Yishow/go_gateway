import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { RecordingHomeMemoryView } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView';
import type { MemoryRuleInput } from '@/features/datalink/workbench-v2/recording-home/recordingHomeMemoryMap';
beforeEach(async () => { await i18n.changeLanguage('en'); });
afterEach(cleanup);
const rule: MemoryRuleInput = { id: 'flow', label: 'Flow', start_address: '40001', count: 2, data_type: 'float32', skipped_addresses: [] };
function mount(inputs = [rule]) {
  const onOpen = vi.fn(), onAdd = vi.fn();
  render(<RecordingHomeMemoryView protocol="modbus_tcp" inputs={inputs} status="ready" checking={false}
    onOpen={onOpen} onAdd={onAdd} onBack={vi.fn()} onClose={vi.fn()} />);
  fireEvent.click(screen.getByTestId('home-memory-layout-chapters'));
  return { onOpen, onAdd };
}
it('states value count separately from physical footprint in the chapter summary', () => {
  mount();
  const summary = screen.getByTestId('home-memory-chapter-flow').querySelector('summary');
  expect(summary).toHaveTextContent('2 values');
  expect(summary).toHaveTextContent('4 registers');
  expect(summary).toHaveTextContent('40001 – 40004');
  expect(summary).toHaveTextContent('float32');
});
it('does not invent a physical endpoint or footprint for unknown string width', () => {
  mount([{ ...rule, data_type: 'string' }]);
  const summary = screen.getByTestId('home-memory-chapter-flow').querySelector('summary');
  expect(summary).toHaveTextContent('2 values');
  expect(summary).toHaveTextContent('40001 – Unconfirmed');
  expect(summary).not.toHaveTextContent('4 registers');
});
it('keeps every chapter and its native disclosure plus last-range jump reachable', () => {
  const actions = mount(Array.from({ length: 40 }, (_, i) => ({ ...rule, id: `R${i}`, label: `Range ${i}`, start_address: String(40001 + i * 10) })));
  const list = screen.getByTestId('home-memory-chapters').querySelector('.gw-memory-chapter-list');
  expect(list?.querySelectorAll('details')).toHaveLength(40);
  const last = screen.getByTestId('home-memory-chapter-R39');
  expect(last.tagName).toBe('DETAILS');
  last.setAttribute('open', '');
  fireEvent.click(screen.getByTestId('home-memory-range-jump-R39'));
  expect(screen.getByTestId('home-memory-cell-40391')).toHaveAttribute('data-selected', 'true');
  expect(actions.onAdd).not.toHaveBeenCalled(); expect(actions.onOpen).not.toHaveBeenCalled();
});
