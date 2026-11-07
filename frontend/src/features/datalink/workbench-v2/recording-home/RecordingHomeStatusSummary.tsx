import { useTranslation } from 'react-i18next';
import type { RuntimeMetrics } from '../../../../types/datalink';
import type { DeviceConnectionObservation, DeviceConnectionState } from '../../../../types/deviceConnection';
import type { WriteGroup } from '../../../../types/studioV2WriteGroup';
import type { RecordingReadonly } from '../state/recordingWorkflowView';
import type { HomeDeviceSummary } from './recordingHomeCollections';
import { countRecordingHomeStatuses, recordingHomeStatuses } from './recordingHomeStatusCounts';

const states: readonly DeviceConnectionState[] = ['collecting', 'degraded', 'disconnected', 'idle', 'unavailable'];
export function RecordingHomeStatusSummary({ devices, observations, groups, metrics, outdated }: {
  devices: readonly HomeDeviceSummary[]; observations: ReadonlyMap<string, DeviceConnectionObservation>;
  groups: readonly RecordingReadonly<WriteGroup>[]; metrics?: RuntimeMetrics; outdated: boolean;
}) {
  const { t } = useTranslation('recording-home');
  const counts = countRecordingHomeStatuses(groups);
  return <section className="gw-status-summary" aria-label={t('deviceStatus.summary')} data-testid="home-status-summary">
    <div><strong>{t('deviceStatus.devices')}</strong> <span data-testid="device-count-total">{devices.length}</span>
      <ul>{states.map(state => <li key={state} data-testid={`device-count-${state}`}>{t(`deviceStatus.states.${state}`)} {devices.filter(device => (observations.get(device.id)?.state ?? 'unavailable') === state).length}</li>)}</ul>
    </div>
    <div><strong>{t('deviceStatus.recordings')}</strong> <span data-testid="record-count-total">{groups.length}</span>
      <ul>{recordingHomeStatuses.map(status => <li key={status} data-testid={`record-count-${status}`}>{t(`savedStatuses.${status}`)} {counts[status]}</li>)}</ul>
    </div>
    <div data-testid="runtime-write-counts"><strong>{t('deviceStatus.runtimeWrites')}</strong>
      <p>{t('deviceStatus.writeSuccess')} {metrics?.write_success_total ?? '—'} · {t('deviceStatus.writeError')} {metrics?.write_error_total ?? '—'}</p>
    </div>
    {outdated && <p role="status" className="gw-status-outdated">{t('deviceStatus.outdated')}</p>}
  </section>;
}
