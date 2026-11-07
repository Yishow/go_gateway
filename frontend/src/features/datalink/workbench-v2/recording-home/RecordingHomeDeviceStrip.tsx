import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { CircleCheck, CircleHelp, CirclePause, CircleX, TriangleAlert } from 'lucide-react';
import type { DeviceConnectionObservation, DeviceConnectionState } from '../../../../types/deviceConnection';
import type { HomeDeviceSummary } from './recordingHomeCollections';

const icons = { collecting: CircleCheck, degraded: TriangleAlert, disconnected: CircleX, idle: CirclePause, unavailable: CircleHelp };
export function DeviceConnectionBadge({ state }: { state: DeviceConnectionState }) {
  const { t } = useTranslation('recording-home');
  const Icon = icons[state];
  return <span className="gw-connection-badge" data-connection-state={state}><Icon aria-hidden="true" size={15} />{t(`deviceStatus.states.${state}`)}</span>;
}
function relativeRead(lastRead: string | null | undefined, now: number, t: (key: string, options?: { count: number }) => string) {
  const then = lastRead ? Date.parse(lastRead) : NaN;
  if (!Number.isFinite(then)) return t('deviceStatus.readUnknown');
  const seconds = Math.max(0, Math.floor((now - then) / 1000));
  if (seconds < 60) return t('deviceStatus.secondsAgo', { count: seconds });
  if (seconds < 3600) return t('deviceStatus.minutesAgo', { count: Math.floor(seconds / 60) });
  return t('deviceStatus.hoursAgo', { count: Math.floor(seconds / 3600) });
}
export function RecordingHomeDeviceStrip({ devices, observations, selected, now, onSelect, renderDeviceActions }: {
  devices: readonly HomeDeviceSummary[]; observations: ReadonlyMap<string, DeviceConnectionObservation>; selected: string; now: number;
  onSelect: (id: string) => void; renderDeviceActions?: (device: HomeDeviceSummary) => ReactNode;
}) {
  const { t } = useTranslation('recording-home');
  return <section className="gw-device-strip" aria-label={t('deviceStatus.strip')}>
    {devices.map(device => {
      const observation: DeviceConnectionObservation = observations.get(device.id) ?? { state: 'unavailable' };
      return <div className="gw-device-card" key={device.id}>
        <button type="button" data-testid={`device-status-${device.id}`} aria-pressed={selected === device.id} onClick={() => onSelect(device.id)}>
          <strong>{device.name}</strong><span>{t(`protocolNames.${device.protocol}`, { defaultValue: device.protocol })}</span>
          <DeviceConnectionBadge state={observation.state} />
          <span>{observation.collector ? t('deviceStatus.points', { healthy: observation.collector.points_healthy, total: observation.collector.points_total }) : t('deviceStatus.pointsUnknown')}</span>
          <span>{relativeRead(observation.collector?.last_read_at, now, t)}</span>
        </button>{renderDeviceActions?.(device)}
      </div>;
    })}
  </section>;
}
