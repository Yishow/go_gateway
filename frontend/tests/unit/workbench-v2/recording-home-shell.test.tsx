import { afterEach, describe, expect, it, vi } from 'vitest';
import { cleanup, fireEvent, render, screen, within } from '@testing-library/react';
import { RecordingHome } from '@/features/datalink/workbench-v2/recording-home/RecordingHome';
import { parseStudioPresentation } from '@/services/studioRecordingView';
import { recordingHomeGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import { fixture, group } from './recording-home-guidance-fixtures';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, options?: { count?: number; name?: string; time?: string }) => key === 'deviceStatus.attention' ? `${options?.count} ${key}` : key === 'deviceStatus.resumeDraft' ? `${key} ${options?.name}` : key === 'deviceStatus.lastUpdated' ? `${key} ${options?.time}` : key }) }));
afterEach(cleanup);
describe('one readonly home renderer', () => {
  it('makes add recording the only primary button and names a resumable draft', () => {
    const guidance = recordingHomeGuidance(fixture(['D1'], true, [group('D', 'draft')]));
    const navigate = vi.fn();
    render(<RecordingHome guidance={guidance} location={{ search: '?view=recordings', hash: '' }} onNavigate={navigate} onRetry={vi.fn()} updatedAt={Date.parse('2026-11-07T12:00:00Z')} />);
    const main = screen.getByTestId('recording-home');
    const primary = main.querySelectorAll('button.gw-primary');
    expect(primary).toHaveLength(1);
    expect(primary[0]).toHaveTextContent('simple.create');
    expect(screen.getByTestId('home-draft-banner')).toHaveTextContent('Recording D');
    expect(screen.getByText(/^deviceStatus\.lastUpdated \d{2}:\d{2}:\d{2}$/)).toBeInTheDocument();
    expect(screen.getByLabelText('deviceStatus.freshnessHelp')).toBeInTheDocument();
    expect(screen.getByRole('tooltip', { hidden: true })).toHaveTextContent('observationBoundary');
    fireEvent.click(within(screen.getByTestId('home-draft-banner')).getByRole('button', { name: 'deviceStatus.resume' }));
    expect(navigate.mock.calls[0][0]).toContain('group_id=D');
    cleanup();
    render(<RecordingHome guidance={recordingHomeGuidance(fixture())} location={{ search: '', hash: '' }} onRetry={vi.fn()} />);
    expect(screen.queryByTestId('home-draft-banner')).not.toBeInTheDocument();
  });
  it('makes add device the only primary button in devices mode', () => {
    const createDevice = vi.fn();
    render(<RecordingHome guidance={recordingHomeGuidance(fixture())} mode="devices" onCreateDevice={createDevice}
      location={{ search: '?view=recordings', hash: '' }} onRetry={vi.fn()} />);
    const primary = screen.getByTestId('recording-home').querySelectorAll('button.gw-primary');
    expect(primary).toHaveLength(1); expect(primary[0]).toHaveTextContent('createDevice');
    expect(screen.queryByRole('button', { name: 'simple.create' })).not.toBeInTheDocument();
    fireEvent.click(primary[0]); expect(createDevice).toHaveBeenCalledOnce();
  });
  it.each([['deleted', false], ['running', true]] as const)('defaults the unconfirmed disclosure for %s to open=%s', (status, open) => {
    const input = fixture(['D1'], true, [group('G', status)]);
    const guidance = recordingHomeGuidance({ ...input, owners: undefined });
    render(<RecordingHome guidance={guidance} location={{ search: '', hash: '' }} onRetry={vi.fn()} />);
    const disclosure = screen.getByTestId('recording-home-unconfirmed');
    expect(disclosure).toHaveTextContent('1 deviceStatus.attention');
    expect(disclosure.hasAttribute('open')).toBe(open);
    fireEvent.click(within(disclosure).getByText('1 deviceStatus.attention'));
    expect(disclosure.hasAttribute('open')).toBe(!open);
  });
  it('opens the attention region for an observed delivery fault and names the actual fault', () => {
    const input = fixture(['D1'], true, [group('G')]);
    const delivery = input.evidence.delivery?.G;
    if (delivery?.state !== 'available') throw new Error('fixture');
    const guidance = recordingHomeGuidance({ ...input, evidence: { ...input.evidence, delivery: { G: { ...delivery,
      value: { ...delivery.value, intake: { state: 'blocked' } } } } } });
    render(<RecordingHome guidance={guidance} location={{ search: '', hash: '' }} onRetry={vi.fn()} />);
    const disclosure = screen.getByTestId('recording-home-unconfirmed');
    expect(disclosure).toHaveAttribute('open');
    expect(disclosure).toHaveTextContent('1 deviceStatus.attention');
    expect(within(disclosure).getByText('reason.delivery_attention')).toBeInTheDocument();
  });
  it('opens scoped logs for the selected record while the header tools remain global', () => {
    const navigate = vi.fn(); render(<RecordingHome guidance={recordingHomeGuidance(fixture(['D1'], true, [group('G')]))}
      location={{ search: '?view=recordings', hash: '#same-row' }} onNavigate={navigate} onRetry={vi.fn()} />);
    fireEvent.click(screen.getByTestId('recording-scoped-logs'));
    expect(navigate).toHaveBeenLastCalledWith('/studio/logs?entry=recording&workspace_id=W1&group_id=G&group_revision=G2#same-row');
    expect(navigate).toHaveBeenCalledOnce();
  });
  it.each(['sources', 'destination'] as const)('preserves the selected recording %s repair and return position', section => {
    const guidance = recordingHomeGuidance(fixture(['D1'], true, [group('G')]));
    const navigate = vi.fn();
    render(<RecordingHome guidance={{ ...guidance, primaryAction: { kind: 'navigate', target: 'recording',
      scope: { workspaceId: 'W1', workspaceRevision: 'W2', epoch: 2 }, groupId: 'G', groupRevision: 'G2', section, intent: null } }}
      location={{ search: '?view=recordings', hash: '#same-row' }} onNavigate={navigate} onRetry={vi.fn()} />);
    fireEvent.click(screen.getByText('openSettings'));
    const target = new URL(navigate.mock.calls[0][0], 'http://localhost');
    expect(target.hash).toBe('#same-row');
    expect(parseStudioPresentation(target.search)).toMatchObject({ repairIntent: {
      workspaceId: 'W1', workspaceRevision: 'W2', groupId: 'G', groupRevision: 'G2', section,
    } });
  });
  it('preserves the original operation repair instead of reducing it to an edit', () => {
    const value = group('G'), base = fixture(['D1'], true, [value]);
    if (base.owners?.G.state !== 'available') throw new Error('fixture');
    const guidance = recordingHomeGuidance({ ...base, owners: { G: { ...base.owners.G,
      value: base.owners.G.value.map(owner => owner.owner === 'start'
        ? { ...owner, state: 'unresolved' as const, reason: 'unknown' as const, operationId: 'O1' } : owner) } } });
    const navigate = vi.fn(); render(<RecordingHome guidance={guidance} location={{ search: '?view=recordings', hash: '#same-row' }} onNavigate={navigate} onRetry={vi.fn()} />);
    expect(screen.getByText('reason.operation_recovery')).toBeInTheDocument();
    fireEvent.click(screen.getByText('openSettings'));
    expect(parseStudioPresentation(new URL(navigate.mock.calls[0][0], 'http://localhost').search)).toMatchObject({
      repairIntent: { section: 'operation', operation: { kind: 'query_operation', owner: 'start', operationId: 'O1', deviceId: 'D1' } },
    });
  });
  it.each([
    ['missing_devices', fixture()], ['missing_points', fixture(['D1'])], ['missing_recordings', fixture(['D1'], true)],
    ['draft', fixture(['D1'], true, [group('G', 'draft')])], ['recordings', fixture(['D1'], true, [group('G')])],
  ] as const)('renders %s through the same real resolver with the explicit settings entry', (reason, input) => {
    const navigate = vi.fn(); const retry = vi.fn(); const configurePoints = vi.fn();
    render(<RecordingHome guidance={recordingHomeGuidance(input)} location={{ search: '?view=recordings&section=database', hash: '#return' }} onNavigate={navigate} onRetry={retry} onConfigurePoints={configurePoints} />);
    expect(screen.getByText(`reason.${reason}`)).toBeInTheDocument(); expect(navigate).not.toHaveBeenCalled(); expect(configurePoints).not.toHaveBeenCalled();
    fireEvent.click(screen.getByText('openSettings'));
    if (reason === 'missing_points') {
      expect(configurePoints).toHaveBeenCalledExactlyOnceWith('D1'); expect(navigate).not.toHaveBeenCalled();
    } else {
      expect(navigate).toHaveBeenCalledTimes(1); expect(navigate.mock.calls[0][0]).toMatch(/^\/studio\/v2\?/); expect(configurePoints).not.toHaveBeenCalled();
    }
    expect(retry).not.toHaveBeenCalled();
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
    expect(within(screen.getByTestId('recording-home-guidance')).getByText('reason.delivery_attention')).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Recording D' })); expect(navigate.mock.calls[0][0]).toContain('group_id=D');
    fireEvent.click(screen.getByRole('button', { name: 'Recording F' })); expect(navigate.mock.calls[1][0]).toContain('group_id=F');
  });
  it('advanced return preserves query/hash and tools remain reachable', () => {
    const navigate = vi.fn(); render(<RecordingHome guidance={recordingHomeGuidance(fixture())} location={{ search: '?section=database&view=recordings', hash: '#saved' }} onNavigate={navigate} onRetry={vi.fn()} />);
    fireEvent.click(screen.getByText('advanced')); expect(navigate).toHaveBeenLastCalledWith('/studio/v2?section=database&view=advanced#saved');
    for (const [label, target] of [['runtime', '/studio/runtime'], ['logs', '/studio/logs'], ['test', '/test']]) {
      fireEvent.click(within(screen.getByTestId('recording-product-header')).getByText(label)); expect(navigate).toHaveBeenLastCalledWith(target);
    }
  });
  it('missing owner observations are unknown, not a fabricated fault requiring repair', () => {
    const input = fixture(['D1'], true, [group('G')]); const navigate = vi.fn();
    render(<RecordingHome guidance={recordingHomeGuidance({ ...input, owners: undefined })}
      location={{ search: '?view=recordings', hash: '' }} onNavigate={navigate} onRetry={vi.fn()} />);
    expect(screen.getByText('statusUnknown')).toBeInTheDocument(); expect(screen.queryByText('attention')).not.toBeInTheDocument();
    // The unconfirmed list also names G; the saved-list entry is its button.
    fireEvent.click(screen.getByRole('button', { name: 'Recording G' })); expect(navigate.mock.calls[0][0]).toContain('group_id=G');
  });

});
