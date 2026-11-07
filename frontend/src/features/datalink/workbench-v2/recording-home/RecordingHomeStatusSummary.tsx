import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { Activity, Cpu, Rows3, TriangleAlert } from 'lucide-react';
import type { RuntimeMetrics } from '../../../../types/datalink';
import type { DeviceConnectionObservation, DeviceConnectionState } from '../../../../types/deviceConnection';
import type { WriteGroup, WriteGroupStatus } from '../../../../types/studioV2WriteGroup';
import type { RecordingReadonly } from '../state/recordingWorkflowView';
import type { HomeDeviceSummary } from './recordingHomeCollections';
import { countRecordingHomeStatuses, recordingHomeStatuses } from './recordingHomeStatusCounts';

type Tone = 'success' | 'warning' | 'danger' | 'idle' | 'unknown' | 'accent';
interface Segment { key: string; label: string; value: number; tone: Tone; testId?: string }

const deviceStates: readonly DeviceConnectionState[] = ['collecting', 'degraded', 'disconnected', 'idle', 'unavailable'];
const deviceTones: Record<DeviceConnectionState, Tone> = { collecting: 'success', degraded: 'warning', disconnected: 'danger', idle: 'idle', unavailable: 'unknown' };
const recordTones: Record<WriteGroupStatus, Tone> = { draft: 'unknown', ready: 'accent', running: 'success', disabled: 'idle', deleted: 'danger' };
const formatCount = (value: number) => value.toLocaleString();

/** Proportional bar plus a legend; every count stays readable as text, the bar only adds shape. */
function Distribution({ label, segments }: { label: string; segments: readonly Segment[] }) {
  const total = segments.reduce((sum, segment) => sum + segment.value, 0);
  return <>
    <div className="gw-stat-bar" aria-hidden="true">
      {total > 0 && segments.filter(segment => segment.value > 0).map(segment =>
        <span key={segment.key} data-tone={segment.tone} style={{ flexGrow: segment.value }} />)}
    </div>
    <ul className="gw-stat-legend" aria-label={label}>{segments.map(segment => <li key={segment.key} data-testid={segment.testId} data-tone={segment.tone} data-zero={segment.value === 0 || undefined}>
      <i aria-hidden="true" />{segment.label} <b>{formatCount(segment.value)}</b>
    </li>)}</ul>
  </>;
}
function StatCard({ icon, label, children, testId }: { icon: ReactNode; label: string; children: ReactNode; testId?: string }) {
  return <div className="gw-stat-card" data-testid={testId}>
    <p className="gw-stat-head"><span className="gw-stat-icon" aria-hidden="true">{icon}</span><strong>{label}</strong></p>
    {children}
  </div>;
}

export function RecordingHomeStatusSummary({ devices, observations, groups, metrics, outdated }: {
  devices: readonly HomeDeviceSummary[]; observations: ReadonlyMap<string, DeviceConnectionObservation>;
  groups: readonly RecordingReadonly<WriteGroup>[]; metrics?: RuntimeMetrics; outdated: boolean;
}) {
  const { t } = useTranslation('recording-home');
  const counts = countRecordingHomeStatuses(groups);
  const success = metrics?.write_success_total, failed = metrics?.write_error_total;
  const writes = typeof success === 'number' && typeof failed === 'number' ? success + failed : undefined;
  const rate = writes ? Math.floor((success! / writes) * 1000) / 10 : undefined;
  return <section className="gw-status-summary" aria-label={t('deviceStatus.summary')} data-testid="home-status-summary">
    {outdated && <p role="status" className="gw-status-outdated"><TriangleAlert aria-hidden="true" size={16} />{t('deviceStatus.outdated')}</p>}
    <StatCard icon={<Cpu size={18} />} label={t('deviceStatus.devices')}>
      <p className="gw-stat-value"><span data-testid="device-count-total">{devices.length}</span></p>
      <Distribution label={t('deviceStatus.devices')} segments={deviceStates.map(state => ({ key: state, label: t(`deviceStatus.states.${state}`), tone: deviceTones[state],
        testId: `device-count-${state}`, value: devices.filter(device => (observations.get(device.id)?.state ?? 'unavailable') === state).length }))} />
    </StatCard>
    <StatCard icon={<Rows3 size={18} />} label={t('deviceStatus.recordings')}>
      <p className="gw-stat-value"><span data-testid="record-count-total">{groups.length}</span></p>
      <Distribution label={t('deviceStatus.recordings')} segments={recordingHomeStatuses.map(status => ({ key: status, label: t(`savedStatuses.${status}`),
        tone: recordTones[status], testId: `record-count-${status}`, value: counts[status] }))} />
    </StatCard>
    <StatCard icon={<Activity size={18} />} label={t('deviceStatus.runtimeWrites')} testId="runtime-write-counts">
      <p className="gw-stat-value"><span>{typeof success === 'number' ? formatCount(success) : '—'}</span>
        <small>{success === 0 && failed === 0 ? t('deviceStatus.noWritesYet') : rate === undefined ? t('deviceStatus.successRateUnknown') : t('deviceStatus.successRate', { rate })}</small></p>
      {writes === undefined ? <p className="gw-stat-empty">{t('deviceStatus.writeSuccess')} — · {t('deviceStatus.writeError')} —</p>
        : <Distribution label={t('deviceStatus.runtimeWrites')} segments={[
          { key: 'success', label: t('deviceStatus.writeSuccess'), tone: 'success', value: success! },
          { key: 'error', label: t('deviceStatus.writeError'), tone: 'danger', value: failed! }]} />}
    </StatCard>
  </section>;
}
