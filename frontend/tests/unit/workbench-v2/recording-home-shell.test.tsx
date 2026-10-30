import { afterEach, describe, expect, it, vi } from 'vitest';
import { cleanup, fireEvent, render, screen } from '@testing-library/react';
import { RecordingHome } from '@/features/datalink/workbench-v2/recording-home/RecordingHome';
import { recordingHomeGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import { fixture, group } from './recording-home-guidance-fixtures';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
afterEach(cleanup);
describe('one readonly home renderer', () => {
  it.each([
    ['missing_devices', fixture()], ['missing_points', fixture(['D1'])], ['missing_recordings', fixture(['D1'], true)],
    ['draft', fixture(['D1'], true, [group('G', 'draft')])], ['recordings', fixture(['D1'], true, [group('G')])],
  ] as const)('renders %s through the same real resolver with navigation only', (reason, input) => {
    const navigate = vi.fn(); const retry = vi.fn();
    render(<RecordingHome guidance={recordingHomeGuidance(input)} location={{ search: '?view=recordings&section=database', hash: '#return' }} onNavigate={navigate} onRetry={retry} />);
    expect(screen.getByText(`reason.${reason}`)).toBeInTheDocument(); expect(navigate).not.toHaveBeenCalled();
    fireEvent.click(screen.getByText('openSettings')); expect(navigate).toHaveBeenCalledTimes(1);
    expect(navigate.mock.calls[0][0]).toMatch(/^\/studio\/v2\?/); expect(retry).not.toHaveBeenCalled();
  });
  it.each(['loading', 'error', 'unknown'] as const)('renders %s independently from empty, with safe recovery', state => {
    const base = fixture(); const input = { ...base, evidence: { ...base.evidence, workspace: { state, scope: base.evidence.scope } } };
    const retry = vi.fn(); render(<RecordingHome guidance={recordingHomeGuidance(input)} location={{ search: '?view=recordings', hash: '' }} onRetry={retry} />);
    expect(screen.getByText(state)).toBeInTheDocument(); expect(screen.queryByText('reason.missing_devices')).not.toBeInTheDocument();
    fireEvent.click(screen.getByText('retry')); expect(retry).toHaveBeenCalledOnce();
  });
  it('fault and draft remain named and reachable without mutating either', () => {
    const base = fixture(['D1'], true, [group('D', 'draft'), group('F')]);
    const delivery = base.evidence.delivery?.F;
    if (delivery?.state !== 'available') throw new Error('fixture');
    const input = { ...base, evidence: { ...base.evidence, delivery: { F: { ...delivery, value: { ...delivery.value, intake: { state: 'blocked' as const } } } } } };
    const navigate = vi.fn(); render(<RecordingHome guidance={recordingHomeGuidance(input)} location={{ search: '?view=recordings', hash: '' }} onNavigate={navigate} onRetry={vi.fn()} />);
    expect(screen.getByText('reason.delivery_attention')).toBeInTheDocument();
    fireEvent.click(screen.getByText('Recording D')); expect(navigate.mock.calls[0][0]).toContain('group_id=D');
    fireEvent.click(screen.getByText('Recording F')); expect(navigate.mock.calls[1][0]).toContain('group_id=F');
  });
  it('advanced return preserves query/hash and tools remain reachable', () => {
    const navigate = vi.fn(); render(<RecordingHome guidance={recordingHomeGuidance(fixture())} location={{ search: '?section=database&view=recordings', hash: '#saved' }} onNavigate={navigate} onRetry={vi.fn()} />);
    fireEvent.click(screen.getByText('advanced')); expect(navigate).toHaveBeenLastCalledWith('/studio/v2?section=database&view=advanced#saved');
    for (const [label, target] of [['runtime', '/studio/runtime'], ['logs', '/studio/logs'], ['test', '/test']]) {
      fireEvent.click(screen.getByText(label)); expect(navigate).toHaveBeenLastCalledWith(target);
    }
  });
  it('missing owner observations are unknown, not a fabricated fault requiring repair', () => {
    const input = fixture(['D1'], true, [group('G')]); const navigate = vi.fn();
    render(<RecordingHome guidance={recordingHomeGuidance({ ...input, owners: undefined })}
      location={{ search: '?view=recordings', hash: '' }} onNavigate={navigate} onRetry={vi.fn()} />);
    expect(screen.getByText('statusUnknown')).toBeInTheDocument(); expect(screen.queryByText('attention')).not.toBeInTheDocument();
    fireEvent.click(screen.getByText('Recording G')); expect(navigate.mock.calls[0][0]).toContain('group_id=G');
  });

});
