import { act, cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { RecordingHomeMemoryView } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView';

beforeEach(async () => { await i18n.changeLanguage('en'); });
afterEach(cleanup);

it.each(['grid', 'navigator', 'map', 'chapters', 'detail'])('%s keeps explicit inspection expanded while a toolbar control receives focus', layout => {
  const onAdd = vi.fn(), onOpen = vi.fn();
  render(<RecordingHomeMemoryView protocol="modbus_tcp" status="ready" checking={false}
    inputs={[{ id: 'A', label: 'Pressure', start_address: '40001', count: 4, data_type: 'int16', skipped_addresses: [] }]}
    onAdd={onAdd} onOpen={onOpen} onBack={vi.fn()} onClose={vi.fn()} />);
  fireEvent.click(screen.getByTestId(`home-memory-layout-${layout}`));
  const first = screen.getByTestId('home-memory-cell-40001'), second = screen.getByTestId('home-memory-cell-40002');
  // A fallback tab stop must not expand the initial view before anyone inspects a cell.
  expect(first).not.toHaveAttribute('data-expanded', 'true');
  fireEvent.mouseMove(second);
  expect(second).not.toHaveAttribute('data-expanded', 'true');
  act(() => first.focus());
  expect(first).toHaveAttribute('data-expanded', 'true');
  act(() => within(screen.getByRole('group', { name: 'Address blocks' })).getByRole('button', { name: /Next/ }).focus());
  // Blur must not collapse the row and move the control between pointerdown and pointerup.
  expect(first).toHaveAttribute('data-expanded', 'true');
  act(() => second.focus());
  expect(first).not.toHaveAttribute('data-expanded', 'true');
  expect(second).toHaveAttribute('data-expanded', 'true');
  fireEvent.click(screen.getByTestId('home-memory-area-IR'));
  expect(screen.getByTestId('home-memory-cell-30002')).not.toHaveAttribute('data-expanded', 'true');
  expect(onAdd).not.toHaveBeenCalled();
  expect(onOpen).not.toHaveBeenCalled();
});
