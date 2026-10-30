import { useState } from 'react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import en from '@/i18n/locales/en/recording-home.json';
import zhTW from '@/i18n/locales/zh-TW/recording-home.json';
import { RecordingHome } from '@/features/datalink/workbench-v2/recording-home/RecordingHome';
import { RecordingHomeCollection } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeCollection';
import { initialCollectionState } from '@/features/datalink/workbench-v2/recording-home/recordingHomeCollections';
import { recordingHomeGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import { recordingHomeResultSummary } from '@/features/datalink/workbench-v2/recording-home/recordingHomeResultSummary';
import { fixture, group } from './recording-home-guidance-fixtures';
const actualLabels = vi.hoisted(() => ({ current: null as null | { details: string; detailsFor: string } }));
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, options?: { name?: string }) => actualLabels.current && key === 'result.details' ? actualLabels.current.details : actualLabels.current && key === 'result.detailsFor' ? actualLabels.current.detailsFor.replace('{{name}}', options?.name ?? '') : options?.name ? `${key}: ${options.name}` : key }) }));
afterEach(() => { cleanup(); actualLabels.current = null; });
const longName = '鍋爐室一號溫度與壓力資料記錄'.repeat(10);
const groups = [{ ...group('G'), name: longName }, { ...group('H'), name: 'Second recording' }];
const guidance = recordingHomeGuidance(fixture(['D1'], true, groups));
const summaries = Object.fromEntries(guidance.recordingItems.map(item => [item.id, recordingHomeResultSummary(item.evidence)]));
function Host({ restoredFocus = null, open = vi.fn(), remember = vi.fn() }: { restoredFocus?: string | null; open?: (id: string) => void; remember?: (top: number, focus: string | null, page: number) => void }) {
  const [state, setState] = useState(initialCollectionState());
  return <RecordingHome guidance={guidance} location={{ search: '?view=recordings', hash: '' }} onRetry={vi.fn()}
    collection={<RecordingHomeCollection groups={groups} devices={[{ id: 'D1', name: 'Device', protocol: 'modbus_tcp', availability: 'unknown' }]}
      state={state} guidance={guidance} summaries={summaries} observedIds={[]} unverifiedObserved={false} sourcesAvailable
      restored={{ state: initialCollectionState(), workspaceRevision: 'W2', scrollTop: 0, focusId: restoredFocus }} onChange={setState}
      onOpen={open} onOpenDevice={vi.fn()} onRemember={remember} onRecheck={vi.fn()} />} />;
}
describe('actual home semantic and keyboard/focus contracts (no geometry claims)', () => {
  it('uses the scoped dual-theme surface, real search labels and a keyboard-accessible named table region', () => {
    render(<Host />);
    expect(screen.getByRole('main')).toHaveClass('gw-recording-home');
    expect(screen.getByRole('search', { name: 'filters' })).toBeInTheDocument();
    const search = screen.getByRole('searchbox', { name: 'keyword' }); expect(search).toHaveAttribute('type', 'search');
    const region = screen.getByRole('region', { name: 'tableRegion' }); expect(region).toHaveAttribute('tabindex', '0');
    region.focus(); expect(region).toHaveFocus();
    expect(within(region).getByRole('table', { name: 'saved' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'previous' })).toBeDisabled();
  });
  it('keeps full long names and uniquely names each native evidence disclosure', () => {
    render(<Host />);
    expect(screen.getByRole('button', { name: longName })).toHaveTextContent(longName);
    const detail = screen.getByLabelText(`result.detailsFor: ${longName}`);
    expect(detail.tagName).toBe('SUMMARY'); expect(detail.closest('details')).not.toHaveAttribute('open');
    fireEvent.click(detail); expect(detail.closest('details')).toHaveAttribute('open');
    expect(screen.getByLabelText('result.detailsFor: Second recording')).toBeInTheDocument();
  });
  it('typing/filtering does not steal focus or submit any navigation and removed return targets safely focus the heading', () => {
    const open = vi.fn(); render(<Host open={open} restoredFocus="record-removed" />);
    expect(screen.getByRole('heading', { name: 'saved' })).toHaveFocus();
    const search = screen.getByRole('searchbox', { name: 'keyword' }); search.focus();
    fireEvent.change(search, { target: { value: 'nothing matches' } });
    expect(search).toHaveFocus(); expect(screen.getByText('noResults')).toBeInTheDocument();
    fireEvent.submit(screen.getByRole('search', { name: 'filters' })); expect(open).not.toHaveBeenCalled();
  });
  it('records and restores the actual visible record trigger, without turning focus into command permission', () => {
    const remember = vi.fn(), open = vi.fn(); render(<Host restoredFocus="record-G" remember={remember} open={open} />);
    const trigger = screen.getByRole('button', { name: longName }); expect(trigger).toHaveFocus();
    fireEvent.click(trigger); expect(remember).toHaveBeenCalledWith(0, 'record-G', 1); expect(open).toHaveBeenCalledWith('G');
  });
});
it('returns focus to the actual observed-issue trigger even when its table row is filtered out', () => {
  const state = { ...initialCollectionState(), keyword: 'no matching row' };
  const attention = { ...guidance, attentionItems: [{ ...guidance.recordingItems[0], reason: 'source_attention' as const }] };
  const remember = vi.fn();
  render(<RecordingHomeCollection groups={groups} devices={[]} state={state} guidance={attention} observedIds={['G']} unverifiedObserved={false}
    sourcesAvailable restored={{ state, workspaceRevision: 'W2', scrollTop: 0, focusId: 'attention-record-G' }}
    onChange={vi.fn()} onOpen={vi.fn()} onOpenDevice={vi.fn()} onRemember={remember} onRecheck={vi.fn()} />);
  const trigger = screen.getByRole('button', { name: longName }); expect(trigger).toHaveFocus();
  fireEvent.click(trigger); expect(remember).toHaveBeenCalledWith(0, 'attention-record-G', 1);
});
it('keeps all enabled native controls in the natural keyboard order without positive tabindex', () => {
  const view = render(<Host />);
  const controls = [...view.container.querySelectorAll<HTMLElement>('button:not(:disabled), input, select, summary, [tabindex]')];
  expect(controls.length).toBeGreaterThan(10);
  expect(controls.filter(node => Number(node.getAttribute('tabindex') ?? 0) > 0)).toEqual([]);
  for (const node of controls.filter(node => node.getAttribute('tabindex') !== '-1')) {
    node.focus(); expect(node).toHaveFocus();
  }
});

it.each([['en', en], ['zh-TW', zhTW]] as const)('%s disclosure accessible names contain their complete visible label and each record name', (_language, locale) => {
  actualLabels.current = locale.result;
  render(<Host />);
  const disclosures = screen.getAllByText(locale.result.details);
  expect(disclosures).toHaveLength(2);
  for (const group of groups) {
    const disclosure = within(screen.getByTestId(`home-result-${group.id}`)).getByText(locale.result.details);
    expect(disclosure).toHaveAccessibleName(expect.stringContaining(locale.result.details));
    expect(disclosure).toHaveAccessibleName(expect.stringContaining(group.name));
  }
  expect(disclosures[0].getAttribute('aria-label')).not.toBe(disclosures[1].getAttribute('aria-label'));
});
