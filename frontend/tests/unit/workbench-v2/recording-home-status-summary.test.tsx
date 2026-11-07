import { createInstance } from 'i18next';
import { I18nextProvider, initReactI18next } from 'react-i18next';
import { cleanup, render, screen, within } from '@testing-library/react';
import { afterEach, expect, it } from 'vitest';
import { RecordingHomeStatusSummary } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeStatusSummary';
import enHome from '@/i18n/locales/en/recording-home.json';
import enStatus from '@/i18n/locales/en/recording-home-status.json';
import zhHome from '@/i18n/locales/zh-TW/recording-home.json';
import zhStatus from '@/i18n/locales/zh-TW/recording-home-status.json';
import type { RuntimeMetrics } from '@/types/datalink';

const i18n = createInstance();
await i18n.use(initReactI18next).init({ lng: 'en', fallbackLng: false, interpolation: { escapeValue: false },
  resources: { en: { 'recording-home': { ...enHome, ...enStatus } }, 'zh-TW': { 'recording-home': { ...zhHome, ...zhStatus } } } });

afterEach(cleanup);

function show(metrics?: RuntimeMetrics) {
  render(<I18nextProvider i18n={i18n}><RecordingHomeStatusSummary devices={[]} observations={new Map()} groups={[]} metrics={metrics} outdated={false} /></I18nextProvider>);
  return screen.getByTestId('runtime-write-counts');
}

it.each([
  { language: 'en', noWrites: 'No writes yet', unavailable: 'Success rate unavailable' },
  { language: 'zh-TW', noWrites: '尚無寫入', unavailable: '成功率無法取得' },
])('shows known zero writes in $language without reporting the rate as unavailable', async ({ language, noWrites, unavailable }) => {
  await i18n.changeLanguage(language);
  const card = show({ write_success_total: 0, write_error_total: 0 });
  expect(card.querySelector('.gw-stat-value > span')).toHaveTextContent(/^0$/);
  expect(within(card).getByText(noWrites)).toBeVisible();
  expect(within(card).queryByText(unavailable)).not.toBeInTheDocument();
});

it.each([
  { language: 'en', unavailable: 'Success rate unavailable', noWrites: 'No writes yet', counts: 'Successful — · Errors —' },
  { language: 'zh-TW', unavailable: '成功率無法取得', noWrites: '尚無寫入', counts: '成功 — · 錯誤 —' },
])('keeps missing metrics unavailable in $language', async ({ language, unavailable, noWrites, counts }) => {
  await i18n.changeLanguage(language);
  const card = show();
  expect(card.querySelector('.gw-stat-value > span')).toHaveTextContent('—');
  expect(within(card).getByText(unavailable)).toBeVisible();
  expect(within(card).queryByText(noWrites)).not.toBeInTheDocument();
  expect(card).toHaveTextContent(counts);
});

it.each([{ write_success_total: 0 }, { write_error_total: 0 }])('does not mistake partial metrics %j for no writes', async metrics => {
  await i18n.changeLanguage('en');
  const card = show(metrics);
  expect(within(card).getByText('Success rate unavailable')).toBeVisible();
  expect(within(card).queryByText('No writes yet')).not.toBeInTheDocument();
});

it('shows a zero success rate when writes failed', async () => {
  await i18n.changeLanguage('en');
  const card = show({ write_success_total: 0, write_error_total: 2 });
  expect(within(card).getByText('0% success')).toBeVisible();
  expect(within(card).queryByText('No writes yet')).not.toBeInTheDocument();
});

it('hides each distribution bar from assistive technology', async () => {
  await i18n.changeLanguage('en');
  show({ write_success_total: 10, write_error_total: 2 });
  const summary = screen.getByRole('region', { name: 'Home status summary' });
  const bars = summary.querySelectorAll('.gw-stat-bar');
  expect(bars).toHaveLength(3);
  for (const bar of bars) {
    expect(bar).toHaveAttribute('aria-hidden', 'true');
    expect(bar).not.toHaveAttribute('role');
    expect(bar).not.toHaveAttribute('aria-label');
  }
  expect(within(summary).queryAllByRole('img')).toHaveLength(0);
});

it('exposes counts in an individually named legend for each distribution', async () => {
  await i18n.changeLanguage('en');
  show({ write_success_total: 10, write_error_total: 2 });
  const summary = screen.getByRole('region', { name: 'Home status summary' });
  const devices = within(summary).getByRole('list', { name: 'Devices' });
  const recordings = within(summary).getByRole('list', { name: 'Recordings' });
  const writes = within(summary).getByRole('list', { name: 'Runtime cumulative writes' });
  expect(within(devices).getAllByRole('listitem')).toHaveLength(5);
  expect(within(devices).getByTestId('device-count-unavailable')).toHaveTextContent('Status unavailable 0');
  expect(within(recordings).getAllByRole('listitem')).toHaveLength(5);
  expect(within(recordings).getByTestId('record-count-draft')).toHaveTextContent('Draft 0');
  expect(within(writes).getAllByRole('listitem')).toHaveLength(2);
  expect(writes).toHaveTextContent('Successful 10');
  expect(writes).toHaveTextContent('Errors 2');
});
