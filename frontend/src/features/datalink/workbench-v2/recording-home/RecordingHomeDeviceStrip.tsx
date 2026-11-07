import { useEffect, useState, type ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { ChevronDown, CircleCheck, CircleHelp, CirclePause, CircleX, Clock, ListTree, Plus, Search, Settings2, TriangleAlert } from 'lucide-react';
import type { DeviceConnectionObservation, DeviceConnectionState } from '../../../../types/deviceConnection';
import type { HomeDeviceSummary } from './recordingHomeCollections';
import { AcquisitionAnomalyButton, type AnomalySummaryStatus } from './AcquisitionAnomalyButton';
import { DeviceAvailabilityBadge } from './RecordingHomeBadges';

/** Collapsed summary size (DESIGN.md: device summary at most four, full management on the same page). */
export const DEVICE_SUMMARY_LIMIT = 4;
const icons = { collecting: CircleCheck, degraded: TriangleAlert, disconnected: CircleX, idle: CirclePause, unavailable: CircleHelp };
export function DeviceConnectionBadge({ state }: { state: DeviceConnectionState }) {
  const { t } = useTranslation('recording-home');
  const Icon = icons[state];
  return <span className="gw-connection-badge" data-connection-state={state}><Icon aria-hidden="true" size={14} />{t(`deviceStatus.states.${state}`)}</span>;
}
function relativeRead(lastRead: string | null | undefined, now: number, t: (key: string, options?: { count: number }) => string) {
  const then = lastRead ? Date.parse(lastRead) : NaN;
  if (!Number.isFinite(then)) return t('deviceStatus.readUnknown');
  const seconds = Math.max(0, Math.floor((now - then) / 1000));
  if (seconds < 60) return t('deviceStatus.secondsAgo', { count: seconds });
  if (seconds < 3600) return t('deviceStatus.minutesAgo', { count: Math.floor(seconds / 60) });
  return t('deviceStatus.hoursAgo', { count: Math.floor(seconds / 3600) });
}
type OpenKind = 'connection' | 'points';
/** Same-page device management: card body filters records; edit, view, anomalies and delete are one click away. */
export function RecordingHomeDeviceStrip({ devices, observations, selected, now, onSelect, onOpenDevice, onCreateDevice, renderDeviceActions, anomalySummary, onOpenAnomalies, onRemember, returnFocusId }: {
  devices: readonly HomeDeviceSummary[]; observations: ReadonlyMap<string, DeviceConnectionObservation>; selected: string; now: number;
  onSelect: (id: string) => void; onOpenDevice?: (id: string, kind: OpenKind) => void; onCreateDevice?: () => void;
  renderDeviceActions?: (device: HomeDeviceSummary) => ReactNode; anomalySummary?: AnomalySummaryStatus; onOpenAnomalies?: (device: HomeDeviceSummary) => void;
  onRemember?: (focusId: string) => void; returnFocusId?: string | null;
}) {
  const { t } = useTranslation('recording-home');
  const [expanded, setExpanded] = useState(false), [keyword, setKeyword] = useState('');
  const bounded = devices.length > DEVICE_SUMMARY_LIMIT;
  useEffect(() => { if (!bounded) { setExpanded(false); setKeyword(''); } }, [bounded]);
  const query = bounded ? keyword.trim().toLocaleLowerCase() : '';
  const matches = query ? devices.filter(device => `${device.name} ${device.protocol}`.toLocaleLowerCase().includes(query)) : devices;
  const open = bounded && (expanded || Boolean(query));
  const shown = open ? matches : matches.filter((device, index) => index < DEVICE_SUMMARY_LIMIT || device.id === selected ||
    returnFocusId === `device-${device.id}-connection` || returnFocusId === `device-${device.id}-points`);
  const openDevice = (event: React.MouseEvent<HTMLButtonElement>, id: string, kind: OpenKind) => {
    event.currentTarget.focus(); onRemember?.(`device-${id}-${kind}`); onOpenDevice?.(id, kind);
  };
  return <section className="gw-device-strip" aria-label={t('deviceStatus.strip')}>
    <div className="gw-section-head">
      <div className="gw-section-title"><h2>{t('deviceStatus.devices')}</h2><span className="gw-count">{devices.length}</span></div>
      {bounded && <div className="gw-section-tools">
        <label className="gw-search-field"><Search aria-hidden="true" size={16} /><span className="sr-only">{t('deviceKeyword')}</span>
          <input type="search" value={keyword} placeholder={t('deviceKeyword')} onChange={event => setKeyword(event.target.value)} /></label>
        <button type="button" className="gw-ghost" aria-expanded={open} onClick={() => { setExpanded(!open); if (open) setKeyword(''); }}>
          {t(open ? 'deviceStatus.showLess' : 'deviceStatus.showAll', { count: devices.length })}<ChevronDown aria-hidden="true" size={16} data-open={open || undefined} /></button>
      </div>}
    </div>
    {query && !matches.length && <p role="status" className="gw-empty-note">{t('noResults')}</p>}
    <div className="gw-device-grid">
      {shown.map(device => {
        const observation: DeviceConnectionObservation = observations.get(device.id) ?? { state: 'unavailable' };
        const collector = observation.collector;
        const health = collector && collector.points_total > 0 ? Math.min(100, Math.round(collector.points_healthy / collector.points_total * 100)) : 0;
        return <article className="gw-device-card" key={device.id} data-card-state={observation.state} aria-current={selected === device.id || undefined}>
          <button type="button" className="gw-device-card-main" data-testid={`device-status-${device.id}`} aria-pressed={selected === device.id}
            title={t('deviceStatus.filterHint')} onClick={() => onSelect(device.id)}>
            <span className="gw-device-card-head">
              <span className="gw-device-card-title"><strong>{device.name}</strong>
                <span>{t(`protocolNames.${device.protocol}`, { defaultValue: device.protocol })}</span></span>
              <DeviceConnectionBadge state={observation.state} />
            </span>
            {device.availability !== 'available' && <DeviceAvailabilityBadge availability={device.availability} />}
            <span className="gw-device-health">
              <span className="gw-device-health-line"><span>{collector ? t('deviceStatus.points', { healthy: collector.points_healthy, total: collector.points_total }) : t('deviceStatus.pointsUnknown')}</span>
                <span className="gw-device-card-read"><Clock aria-hidden="true" size={13} />{relativeRead(collector?.last_read_at, now, t)}</span></span>
              <span className="gw-device-meter" aria-hidden="true"><span style={{ width: `${health}%` }} /></span>
            </span>
          </button>
          {onOpenAnomalies && <div className="gw-device-card-anomalies">
            <AcquisitionAnomalyButton deviceId={device.id} status={anomalySummary} onOpen={() => onOpenAnomalies(device)} /></div>}
          <div className="gw-device-card-actions">
            {onOpenDevice && <>
              <button type="button" data-device-action="connection" data-home-focus={`device-${device.id}-connection`} aria-label={t('deviceStatus.editNamed', { name: device.name })} onClick={event => openDevice(event, device.id, 'connection')}>
                <Settings2 aria-hidden="true" size={15} /><span>{t('deviceConnection')}</span></button>
              <button type="button" data-device-action="points" data-home-focus={`device-${device.id}-points`} data-testid={`home-device-points-${device.id}`} aria-label={t('deviceStatus.viewNamed', { name: device.name })} onClick={event => openDevice(event, device.id, 'points')}>
                <ListTree aria-hidden="true" size={15} /><span>{t('devicePoints')}</span></button>
            </>}
            {renderDeviceActions?.(device)}
          </div>
        </article>;
      })}
      {onCreateDevice && <button type="button" className="gw-device-card-add" onClick={onCreateDevice}>
        <Plus aria-hidden="true" size={18} /><span>{t('createDevice')}</span>
      </button>}
    </div>
  </section>;
}
