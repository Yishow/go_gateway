import { render, screen } from '@testing-library/react';
import { expect, it, vi } from 'vitest';
import { SummaryRail } from '@/features/datalink/workbench-v2/shell/SummaryRail';
import { savedGroupState } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
it('keeps the connection summary from presenting legacy table and interval as effective group settings', () => {
  Object.defineProperty(window, 'innerWidth', { configurable: true, value: 1440 });
  const state = savedGroupState();
  state.db.connector = { ...state.db.connector, name: 'Owned SQLite', schema: 'old_schema', table: 'legacy_sensor_readings', write_interval_seconds: 5 };
  render(<SummaryRail state={state} />);
  const rail = screen.getByTestId('summary-rail');
  expect(rail).toHaveTextContent('Owned SQLite');
  expect(rail).toHaveTextContent('sqlite');
  expect(rail).not.toHaveTextContent('legacy_sensor_readings');
  expect(rail).not.toHaveTextContent('5s');
  expect(rail).toHaveTextContent('Basic');
  expect(rail).toHaveTextContent('Advanced');
});
