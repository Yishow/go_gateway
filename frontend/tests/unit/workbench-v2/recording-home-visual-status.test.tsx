import { useState } from 'react';
import { afterEach, describe, expect, it, vi } from 'vitest';
import { cleanup, render, screen, within } from '@testing-library/react';
import { RecordingHomeCollection } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeCollection';
import { initialCollectionState, type HomeDeviceSummary } from '@/features/datalink/workbench-v2/recording-home/recordingHomeCollections';
import { recordingHomeGuidance } from '@/features/datalink/workbench-v2/recording-home/recordingHomeGuidance';
import { DeviceAvailabilityBadge } from '@/features/datalink/workbench-v2/recording-home/RecordingHomeBadges';
import { recordingHomeResultSummary } from '@/features/datalink/workbench-v2/recording-home/recordingHomeResultSummary';
import { fixture, group } from './recording-home-guidance-fixtures';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, value?: Record<string, unknown>) => value && 'count' in value ? `${key}:${value.count}` : key }) }));
afterEach(cleanup);

const devices: HomeDeviceSummary[] = [
  { id: 'D1', name: 'Boiler PLC', protocol: 'modbus_tcp', availability: 'available' },
  { id: 'D2', name: 'Meter', protocol: 'mc_3e', availability: 'mystery' },
];
const groups = [group('A', 'draft'), group('B', 'draft'), { ...group('R', 'running'), name: 'Recording R' }, group('X', 'deleted')];
const input = fixture(['D1', 'D2'], true, groups);
const guidance = recordingHomeGuidance(input);
const delivered = guidance.recordingItems.find(item => item.id === 'R')!;
// R has a delivery last_sql_committed_at (historical) fact; A has its delivery removed so its SQL fact is unknown.
const summaries = {
  R: recordingHomeResultSummary(delivered.evidence),
  A: recordingHomeResultSummary({ ...guidance.recordingItems.find(item => item.id === 'A')!.evidence, delivery: { state: 'unavailable', reason: 'missing' } }),
};
function Host() {
  const [state, setState] = useState(initialCollectionState());
  return <RecordingHomeCollection groups={groups} devices={devices} state={state} guidance={guidance} summaries={summaries} observedIds={[]} unverifiedObserved={false}
    sourcesAvailable restored={{ state, workspaceRevision: 'W2', scrollTop: 0, focusId: null }}
    onChange={setState} onOpen={vi.fn()} onOpenDevice={vi.fn()} onRemember={vi.fn()} onRecheck={vi.fn()} />;
}
const row = (name: string) => within(screen.getByTestId('home-management-table')).getByText(name).closest('tr')!;
const icon = (node: Element | null) => node?.querySelector('svg[aria-hidden="true"]');

describe('recording rows read as source → recording → destination with status badges', () => {
  it('shows a saved-status badge with a decorative icon and keeps the readable text', () => {
    render(<Host />);
    const badge = row('Recording X').querySelector('[data-status-badge="deleted"]');
    expect(badge).toHaveTextContent('savedStatuses.deleted'); expect(icon(badge)).toBeTruthy();
    expect(row('Recording A').querySelector('[data-status-badge="draft"]')).toHaveTextContent('savedStatuses.draft');
  });
  it('marks the source and destination cells as a flow with icons, naming the real devices and table', () => {
    render(<Host />);
    const source = row('Recording R').querySelector('[data-flow="source"]');
    const destination = row('Recording R').querySelector('[data-flow="destination"]');
    expect(source).toHaveTextContent('Boiler PLC'); expect(icon(source)).toBeTruthy();
    expect(destination).toHaveTextContent('managedTable'); expect(icon(destination)).toBeTruthy();
  });
  it('SQL facts use a historical or unknown icon, never a success check', () => {
    render(<Host />);
    expect(row('Recording R').querySelector('[data-sql-fact="historical_commit"] svg[aria-hidden="true"]')).toBeTruthy();
    expect(row('Recording A').querySelector('[data-sql-fact="unknown"] svg[aria-hidden="true"]')).toBeTruthy();
    expect(screen.getByTestId('home-management-table').querySelector('.lucide-circle-check')).toBeNull();
  });
  it('saved running is drawn as a configuration switch, not live running or success', () => {
    render(<Host />);
    const badge = row('Recording R').querySelector('[data-status-badge="running"]')!;
    expect(badge.querySelector('.lucide-toggle-right')).toBeTruthy();
    expect(screen.getByTestId('home-management-table').querySelector('.lucide-play, .lucide-circle-check')).toBeNull();
  });
});



describe('untrusted availability strings', () => {
  it.each(['__proto__', 'constructor', 'toString', 'hasOwnProperty'])('%s renders as setting status unknown without crashing', value => {
    render(<DeviceAvailabilityBadge availability={value} />);
    const badge = document.querySelector('[data-availability-badge]');
    expect(badge).toHaveAttribute('data-availability-badge', 'unknown');
    expect(badge).toHaveTextContent('configAvailability.unknown');
    expect(icon(badge)).toBeTruthy();
  });
});
