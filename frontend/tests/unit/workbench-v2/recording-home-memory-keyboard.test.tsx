import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import i18n from '@/i18n/config';
import { RecordingHomeMemoryView } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView';

beforeEach(async () => { await i18n.changeLanguage('en'); });
afterEach(cleanup);
function view() {
  const onAdd = vi.fn(), onOpen = vi.fn();
  render(<RecordingHomeMemoryView protocol="mc_3e" status="ready" checking={false} onAdd={onAdd} onOpen={onOpen} onBack={vi.fn()}
    inputs={[{ id: 'I', label: 'TEMP_', start_address: 'D0', count: 2, data_type: 'int16', skipped_addresses: [], proven: ['D0', 'D1'] }]} />);
  return { onAdd, onOpen };
}

it('moves area selection and focus with arrow, Home and End keys on the tab list', () => {
  view();
  const d = screen.getByRole('tab', { name: /^D/ });
  expect(d).toHaveAttribute('aria-selected', 'true'); expect(d).toHaveAttribute('tabindex', '0');
  expect(screen.getByRole('tab', { name: /^W/ })).toHaveAttribute('tabindex', '-1');
  d.focus(); fireEvent.keyDown(d, { key: 'ArrowRight' });
  const w = screen.getByRole('tab', { name: /^W/ });
  expect(w).toHaveAttribute('aria-selected', 'true'); expect(w).toHaveFocus();
  fireEvent.keyDown(w, { key: 'End' });
  expect(screen.getByRole('tab', { name: /^Y/ })).toHaveFocus();
  fireEvent.keyDown(screen.getByRole('tab', { name: /^Y/ }), { key: 'ArrowRight' });
  expect(screen.getByRole('tab', { name: /^D/ })).toHaveFocus();
  fireEvent.keyDown(screen.getByRole('tab', { name: /^D/ }), { key: 'ArrowLeft' });
  expect(screen.getByRole('tab', { name: /^Y/ })).toHaveAttribute('aria-selected', 'true');
  fireEvent.keyDown(screen.getByRole('tab', { name: /^Y/ }), { key: 'Home' });
  expect(screen.getByRole('tab', { name: /^D/ })).toHaveFocus();
  expect(screen.getByRole('tabpanel')).toHaveAttribute('aria-labelledby', screen.getByRole('tab', { name: /^D/ }).id);
});

it('moves focus between cells with arrow keys and activates only the focused cell', () => {
  const { onAdd, onOpen } = view();
  const d0 = screen.getByTestId('home-memory-cell-D0');
  expect(d0).toHaveAttribute('tabindex', '0');
  d0.focus(); fireEvent.keyDown(d0, { key: 'ArrowRight' });
  expect(screen.getByTestId('home-memory-cell-D1')).toHaveFocus();
  fireEvent.keyDown(screen.getByTestId('home-memory-cell-D1'), { key: 'ArrowDown' });
  const d11 = screen.getByTestId('home-memory-cell-D11');
  expect(d11).toHaveFocus(); expect(d11).toHaveAttribute('tabindex', '0');
  fireEvent.click(d11);
  expect(onAdd).toHaveBeenCalledExactlyOnceWith('D11'); expect(onOpen).not.toHaveBeenCalled();
});
