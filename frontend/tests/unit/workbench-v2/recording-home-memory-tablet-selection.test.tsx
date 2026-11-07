import { act, cleanup, render, screen, fireEvent } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { RecordingHomeMemoryView } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView';

beforeEach(async () => { await i18n.changeLanguage('zh-TW'); });
afterEach(cleanup);
it('keeps the complete selected description outside the address tile and updates it without opening settings', () => {
  const onOpen = vi.fn();
  render(<RecordingHomeMemoryView protocol="mc_3e" inputs={[{ id: 'range', label: '溫度長名稱', start_address: 'D1000', count: 2, data_type: 'int16', skipped_addresses: [], proven: ['D1000', 'D1001'] }]}
    status="ready" checking={false} values={{ status: 'live', byCell: new Map([['D:1000:int16', { value: '1234567890123456789' }]]) }}
    onOpen={onOpen} onAdd={vi.fn()} onBack={vi.fn()} onClose={vi.fn()} />);
  fireEvent.change(screen.getByRole('textbox'), { target: { value: 'D1000' } });
  fireEvent.submit(screen.getByRole('textbox').closest('form')!);
  const first = screen.getByTestId('home-memory-cell-D1000');
  act(() => first.focus());
  const inspection = screen.getByTestId('home-memory-inspection');
  expect(first).not.toContainElement(inspection);
  expect(first).not.toHaveTextContent('溫度長名稱');
  expect(inspection).toHaveTextContent('溫度長名稱');
  expect(inspection).toHaveTextContent('1234567890123456789');
  expect(inspection).toHaveAttribute('data-value-key', 'D:1000:int16');
  expect(first).toHaveAttribute('aria-label', expect.stringContaining('溫度長名稱'));
  act(() => screen.getByTestId('home-memory-cell-D1001').focus());
  expect(inspection).toHaveTextContent('D1001');
  expect(onOpen).not.toHaveBeenCalled();
  fireEvent.click(screen.getByTestId('home-memory-area-M'));
  expect(inspection).not.toHaveTextContent('D1001');
  expect(inspection).not.toHaveAttribute('data-value-key');
});
