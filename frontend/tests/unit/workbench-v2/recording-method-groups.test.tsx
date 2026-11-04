import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { RecordingMethodGroups } from '@/features/datalink/workbench-v2/recording-editor/RecordingMethodGroups';
import type { EditorDraft } from '@/features/datalink/workbench-v2/state/writeGroup/draft';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, options?: { count?: number }) => options?.count === undefined ? key : `${key}:${options.count}` }) }));

const member = (point_id: string, calculation: 'latest' | 'mean', interval: number) =>
  ({ device_id: 'D', point_id, tag_id: `t-${point_id}`, calculation, cadence: { interval, unit: 'second', alignment: 'utc_epoch' } });
const draft = { members: [member('a', 'latest', 15), member('b', 'latest', 15), member('c', 'mean', 60), member('d', 'latest', 15)] } as unknown as Pick<EditorDraft, 'members'>;

afterEach(cleanup);

describe('recording method groups', () => {
  it('groups members by method and interval with point counts, in member order', () => {
    render(<RecordingMethodGroups draft={draft} />);
    const rows = within(screen.getByTestId('recording-method-groups')).getAllByRole('listitem');
    expect(rows.map(row => row.textContent)).toEqual([
      'calculation.latest · 15 calculation.secondsimple.summary.points:3simple.goToMember',
      'calculation.mean · 60 calculation.secondsimple.summary.points:1simple.goToMember',
    ]);
  });

  it('moves focus to the first member of the group in the source table', () => {
    const table = document.createElement('div');
    table.innerHTML = '<div data-testid="group-member-calculation-c"><select aria-label="method c"></select></div>';
    document.body.append(table);
    const select = table.querySelector('select')!;
    select.scrollIntoView = vi.fn();
    render(<RecordingMethodGroups draft={draft} />);
    fireEvent.click(screen.getAllByRole('button', { name: 'simple.goToMember' })[1]);
    expect(select).toHaveFocus();
    table.remove();
  });

  it('states the empty case instead of an empty list', () => {
    render(<RecordingMethodGroups draft={{ members: [] } as unknown as Pick<EditorDraft, 'members'>} />);
    expect(screen.getByText('simple.methodEmpty')).toBeVisible();
    expect(screen.queryByRole('list')).not.toBeInTheDocument();
  });
});
