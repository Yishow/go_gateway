import { act, cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { I18nextProvider } from 'react-i18next';
import i18n from '@/i18n/config';
import { RecordingMethodGroups } from '@/features/datalink/workbench-v2/recording-editor/RecordingMethodGroups';
import { GroupEditorMethodSection } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditorMethodSection';
import type { EditorDraft } from '@/features/datalink/workbench-v2/state/writeGroup/draft';

const member = (point_id: string, calculation: 'latest' | 'mean', interval: number) =>
  ({ device_id: 'D', point_id, tag_id: `t-${point_id}`, calculation, cadence: { interval, unit: 'second', alignment: 'utc_epoch' } });
const draft = {
  layout_version: 'managed_long_v1',
  recording_profile: { kind: 'mixed', interval: 15, unit: 'second' },
  members: [member('a', 'latest', 15), member('b', 'latest', 15), member('c', 'mean', 60), member('d', 'latest', 15)],
} as unknown as EditorDraft;

beforeEach(async () => {
  await i18n.changeLanguage('zh-TW');
});
afterEach(cleanup);

describe('recording method groups', () => {
  it('groups members by method and interval with point counts in a single text span and renders localized button text in zh-TW and en', async () => {
    const { rerender } = render(
      <I18nextProvider i18n={i18n}>
        <GroupEditorMethodSection presentation="simple" draft={draft} saved={false} disabled={false} dialect="sqlite" onReplaceDraft={vi.fn()} />
      </I18nextProvider>,
    );
    const container = screen.getByTestId('recording-method-groups');
    const rows = within(container).getAllByRole('listitem');
    expect(rows).toHaveLength(2);
    // Point count must stay in the same value span as method and interval, not split into a trailing .recording-method-groups__count span.
    expect(container.querySelector('.recording-method-groups__count')).toBeNull();
    const valueSpansZh = Array.from(container.querySelectorAll('.recording-method-groups__value')).map(el => el.textContent);
    expect(valueSpansZh).toEqual([
      '最新值 · 15 秒 · 3 點',
      '時間窗平均值 · 60 秒 · 1 點',
    ]);
    expect(within(container).getAllByRole('button', { name: '到來源表格' })).toHaveLength(2);
    expect(within(container).queryByRole('button', { name: '前往點位' })).not.toBeInTheDocument();

    await act(async () => {
      await i18n.changeLanguage('en');
    });
    rerender(
      <I18nextProvider i18n={i18n}>
        <GroupEditorMethodSection presentation="simple" draft={draft} saved={false} disabled={false} dialect="sqlite" onReplaceDraft={vi.fn()} />
      </I18nextProvider>,
    );
    const valueSpansEn = Array.from(container.querySelectorAll('.recording-method-groups__value')).map(el => el.textContent);
    expect(valueSpansEn).toEqual([
      'Latest value · 15 Seconds · 3 points',
      'Average over a window · 60 Seconds · 1 points',
    ]);
    expect(within(container).getAllByRole('button', { name: 'Go to source table' })).toHaveLength(2);
  });

  it('scrolls to and focuses the first member calculation control of the selected group in the source table', () => {
    const table = document.createElement('div');
    table.innerHTML = [
      '<div data-testid="group-member-calculation-a"><select aria-label="method a"></select></div>',
      '<div data-testid="group-member-calculation-b"><select aria-label="method b"></select></div>',
      '<div data-testid="group-member-calculation-c"><select aria-label="method c"></select></div>',
    ].join('');
    document.body.append(table);
    const selectA = table.querySelector<HTMLSelectElement>('[aria-label="method a"]')!;
    const selectC = table.querySelector<HTMLSelectElement>('[aria-label="method c"]')!;
    selectA.scrollIntoView = vi.fn();
    selectC.scrollIntoView = vi.fn();

    render(
      <I18nextProvider i18n={i18n}>
        <RecordingMethodGroups draft={draft} />
      </I18nextProvider>,
    );
    const buttons = screen.getAllByRole('button', { name: '到來源表格' });

    fireEvent.click(buttons[0]);
    expect(selectA.scrollIntoView).toHaveBeenCalledWith({ block: 'center' });
    expect(selectA).toHaveFocus();

    fireEvent.click(buttons[1]);
    expect(selectC.scrollIntoView).toHaveBeenCalledWith({ block: 'center' });
    expect(selectC).toHaveFocus();
    table.remove();
  });

  it('states the empty case instead of an empty list', () => {
    render(
      <I18nextProvider i18n={i18n}>
        <RecordingMethodGroups draft={{ members: [] } as unknown as Pick<EditorDraft, 'members'>} />
      </I18nextProvider>,
    );
    expect(screen.getByText('先選擇點位，再於各列設定記錄數值與間隔。')).toBeVisible();
    expect(screen.queryByRole('list')).not.toBeInTheDocument();
  });
});
