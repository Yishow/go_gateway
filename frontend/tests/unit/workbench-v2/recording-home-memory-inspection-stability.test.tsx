import { act, cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { RecordingHomeMemoryView } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView';

beforeEach(async () => { await i18n.changeLanguage('en'); });
afterEach(cleanup);

it.each(['grid', 'navigator', 'map', 'chapters', 'detail'])('%s keeps separate inspection visible while a toolbar control receives focus', layout => {
  const onAdd = vi.fn(), onOpen = vi.fn();
  render(<RecordingHomeMemoryView protocol="modbus_tcp" status="ready" checking={false}
    inputs={[{ id: 'A', label: 'Pressure', start_address: '40001', count: 4, data_type: 'int16', skipped_addresses: [] }]}
    onAdd={onAdd} onOpen={onOpen} onBack={vi.fn()} onClose={vi.fn()} />);
  fireEvent.click(screen.getByTestId(`home-memory-layout-${layout}`));
  const first = screen.getByTestId('home-memory-cell-40001'), second = screen.getByTestId('home-memory-cell-40002');
  const inspection = screen.getByTestId('home-memory-inspection');
  expect(inspection).not.toHaveTextContent('40001');
  fireEvent.mouseMove(second);
  expect(inspection).not.toHaveTextContent('40002');
  act(() => first.focus());
  expect(inspection).toHaveTextContent('40001');
  expect(first).not.toContainElement(inspection);
  act(() => within(screen.getByRole('group', { name: 'Address blocks' })).getByRole('button', { name: /Next/ }).focus());
  // Blur must not collapse the row and move the control between pointerdown and pointerup.
  expect(inspection).toHaveTextContent('40001');
  expect(first).not.toContainElement(inspection);
  act(() => second.focus());
  expect(inspection).not.toHaveTextContent('40001');
  expect(inspection).toHaveTextContent('40002');
  fireEvent.click(screen.getByTestId('home-memory-area-IR'));
  expect(inspection).not.toHaveTextContent('40002');
  expect(inspection).not.toHaveAttribute('data-value-key');
  expect(onAdd).not.toHaveBeenCalled();
  expect(onOpen).not.toHaveBeenCalled();
});
