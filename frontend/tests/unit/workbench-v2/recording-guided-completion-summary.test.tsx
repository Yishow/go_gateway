import type { ComponentProps } from 'react';
import { createInstance } from 'i18next';
import { I18nextProvider, initReactI18next } from 'react-i18next';
import { act, cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { RecordingCompletionSummary } from '@/features/datalink/workbench-v2/steps/step4/RecordingCompletionSummary';
import type { RecordingCompletionView } from '@/features/datalink/workbench-v2/steps/step4/recordingCompletion';
import en from '@/i18n/locales/en/workbench-v2-recording-guided.json';
import zh from '@/i18n/locales/zh-TW/workbench-v2-recording-guided.json';

const i18n = createInstance();
await i18n.use(initReactI18next).init({ lng: 'en', fallbackLng: 'en', interpolation: { escapeValue: false },
  resources: { en: { 'workbench-v2-recording-guided': en }, 'zh-TW': { 'workbench-v2-recording-guided': zh } } });
const available = (patch: Partial<Extract<RecordingCompletionView, { state: 'available' }>> = {}): RecordingCompletionView => ({
  state: 'available', operationId: 'operation-private-reference', status: 'partial', deviceIds: ['opaque-A', 'opaque-B'],
  groupIds: ['opaque-group-A', 'opaque-group-B'], activatedDeviceIds: ['opaque-A'], incompleteDeviceIds: ['opaque-B'],
  incompleteGroupIds: ['opaque-group-B'], handoffDeviceId: 'opaque-A', ...patch,
});
function show(view: RecordingCompletionView, props: Partial<ComponentProps<typeof RecordingCompletionSummary>> = {}) {
  const tree = (next: RecordingCompletionView) => <I18nextProvider i18n={i18n}><RecordingCompletionSummary {...props} view={next} /></I18nextProvider>;
  const rendered = render(tree(view)); return { ...rendered, replace: (next: RecordingCompletionView) => rendered.rerender(tree(next)) };
}
/** Capture the actual mounted React command so stale closures, not detached DOM no-ops, are exercised. */
function clickHandler(button: HTMLElement): () => void {
  const key = Object.keys(button).find(value => value.startsWith('__reactProps$'));
  if (!key) throw new Error('Missing mounted React props');
  return (button as unknown as Record<string, { onClick: () => void }>)[key].onClick;
}
beforeEach(async () => { await i18n.changeLanguage('en'); });
afterEach(cleanup);

it.each(['en', 'zh-TW'] as const)('shows confirmed A separately from incomplete B and groups with %s labels', async language => {
  await i18n.changeLanguage(language);
  const resource = language === 'en' ? en : zh;
  show(available(), { deviceNames: { 'opaque-A': 'Pump A', 'opaque-B': 'Pump B' }, groupNames: { 'opaque-group-B': 'Energy total' }, onContinue: vi.fn() });
  const region = screen.getByRole('region', { name: resource.completion.heading });
  expect(within(region).getByRole('list', { name: resource.completion.activated })).toHaveTextContent('Pump A');
  expect(within(region).getByRole('list', { name: resource.completion.activated })).not.toHaveTextContent('Pump B');
  expect(within(region).getByRole('list', { name: resource.completion.incompleteDevices })).toHaveTextContent('Pump B');
  expect(within(region).getByRole('list', { name: resource.completion.incompleteGroups })).toHaveTextContent('Energy total');
  expect(region).toHaveTextContent(resource.states.partial); expect(region).toHaveTextContent(resource.completion.deliverySeparate);
  expect(screen.getByTestId('basic-recording-continue-runtime')).toHaveTextContent(resource.completion.continue);
  expect(region.textContent).not.toMatch(/opaque-|operation-private-reference|completion\./);
});
it('uses full exact ID aliases only and safe numbered fallbacks, never a shortened or private alias', () => {
  show(available(), { deviceNames: { A: 'Wrong short alias', 'opaque-A': 'https://private.internal', 'opaque-B': '/private/db.sqlite' },
    groupNames: { B: 'Wrong group alias', 'opaque-group-B': 'password=secret' }, onContinue: vi.fn() });
  expect(screen.getByRole('list', { name: en.completion.activated })).toHaveTextContent('Device 1');
  expect(screen.getByRole('list', { name: en.completion.incompleteDevices })).toHaveTextContent('Device 2');
  expect(screen.getByRole('list', { name: en.completion.incompleteGroups })).toHaveTextContent('Recording group 2');
  expect(document.body.textContent).not.toMatch(/Wrong|private|opaque-|https:|password|secret/);
});
it.each(['pending', 'running', 'unknown', 'failed'] as const)('does not offer continuation for %s even with a claimed activated ID', status => {
  const onContinue = vi.fn(); show(available({ status }), { onContinue });
  expect(screen.queryByTestId('basic-recording-continue-runtime')).not.toBeInTheDocument(); expect(onContinue).not.toHaveBeenCalled();
});
it.each(['missing', 'scope_mismatch'] as const)('does not offer continuation for unavailable %s evidence', reason => {
  const onContinue = vi.fn(); show({ state: 'unavailable', reason }, { onContinue });
  expect(screen.getByRole('status')).toHaveTextContent(en.reasons[reason]);
  expect(screen.queryByTestId('basic-recording-continue-runtime')).not.toBeInTheDocument(); expect(onContinue).not.toHaveBeenCalled();
});
it('shows no activation confirmation and no Continue when activation is empty', () => {
  show(available({ activatedDeviceIds: [], incompleteDeviceIds: ['opaque-A', 'opaque-B'], handoffDeviceId: undefined }), { onContinue: vi.fn() });
  expect(screen.getByText(en.completion.noneActivated)).toBeVisible(); expect(screen.queryByTestId('basic-recording-continue-runtime')).not.toBeInTheDocument();
});
it.each([
  { handoffDeviceId: 'foreign' }, { activatedDeviceIds: ['opaque-A', 'opaque-B'] }, { deviceIds: ['opaque-B'] }, { handoffDeviceId: undefined },
])('never broadens singleton navigation for inconsistent handoff %j', patch => {
  show(available(patch), { onContinue: vi.fn() }); expect(screen.queryByTestId('basic-recording-continue-runtime')).not.toBeInTheDocument();
});
it('navigates only the confirmed singleton once, without callback on mount or focus', () => {
  const onContinue = vi.fn(); show(available(), { onContinue, primary: true });
  const button = screen.getByTestId('basic-recording-continue-runtime'); button.focus(); expect(onContinue).not.toHaveBeenCalled();
  expect(button).toHaveAttribute('data-primary', 'true'); expect(document.querySelectorAll('[data-primary="true"]')).toHaveLength(1);
  fireEvent.click(button); fireEvent.click(button); expect(onContinue).toHaveBeenCalledExactlyOnceWith('opaque-A');
});
it('does not claim primary styling unless explicitly selected by the owner', () => {
  show(available(), { onContinue: vi.fn() }); expect(screen.getByTestId('basic-recording-continue-runtime')).not.toHaveAttribute('data-primary');
});
it.each(['foreign', 'round_trip', 'unmount'] as const)('rejects a captured old navigation callback after %s', change => {
  const onContinue = vi.fn(); const first = available(); const mounted = show(first, { onContinue });
  const oldClick = clickHandler(screen.getByTestId('basic-recording-continue-runtime'));
  if (change === 'unmount') mounted.unmount();
  else {
    mounted.replace(available({ operationId: 'operation-B', deviceIds: ['other-device'], activatedDeviceIds: ['other-device'], handoffDeviceId: 'other-device' }));
    if (change === 'round_trip') mounted.replace(first);
  }
  act(() => oldClick()); expect(onContinue).not.toHaveBeenCalled();
});

it('keeps a read-only completion summary when no navigation callback is supplied', () => {
  show(available()); expect(screen.getByRole('region', { name: en.completion.heading })).toBeVisible();
  expect(screen.queryByTestId('basic-recording-continue-runtime')).not.toBeInTheDocument();
});
