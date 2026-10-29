import { useTranslation } from 'react-i18next';
import type { RuntimeDashboardRouteState } from './useRuntimeDashboardState';
import type { RuntimeStreamRecovery } from '../../../types/datalink';

interface LiveStateBannerProps {
  routeState: RuntimeDashboardRouteState;
  snapshotReceivedAt?: number | null;
  snapshotStale?: boolean;
  streamRecovery?: RuntimeStreamRecovery | null;
  onRetry?: () => Promise<unknown>;
  onReconnect?: () => void;
  navigateTo?: (target: string) => void;
}

/** Receipt metadata is browser observation, not server as-of or SQL evidence. */
export function LiveStateBanner({ routeState, snapshotReceivedAt, snapshotStale = false, streamRecovery, onRetry, onReconnect, navigateTo }: LiveStateBannerProps) {
  const { t } = useTranslation('runtime-dashboard');
  const live = routeState === 'live';
  const received = typeof snapshotReceivedAt === 'number' && Number.isFinite(snapshotReceivedAt) && snapshotReceivedAt > 0
    ? new Date(snapshotReceivedAt) : null;
  const receivedAt = received && Number.isFinite(received.getTime()) ? received.toISOString() : null;
  return (
    <div className={`rounded-2xl border px-4 py-3 text-sm ${live ? 'border-emerald-500/30 bg-emerald-500/10 text-emerald-100' : 'border-amber-500/30 bg-amber-500/10 text-amber-100'}`} data-testid="runtime-dashboard-live-state-banner">
      <p>{live ? t('banner.live', 'Live monitoring active') : routeState === 'degraded'
        ? t('banner.degraded', 'Monitoring degraded. Retained snapshots remain stale until a successful snapshot read and a ready stream.')
        : t('banner.waiting', 'Waiting for verified stream and snapshot readiness.')}</p>
      {live && <p className="mt-1 text-xs">{t('banner.liveScope', 'Stream and snapshot are ready; this does not prove every point is fresh, collector success or SQL delivery.')}</p>}
      <dl className="mt-2 text-xs">
        <dt>{t('banner.receivedAt', 'Snapshot received at (browser)')}</dt>
        <dd>{receivedAt ? <time data-testid="runtime-snapshot-received-at" dateTime={receivedAt}>{receivedAt}</time> : t('banner.receiptUnavailable', 'Snapshot receipt time unavailable')}</dd>
      </dl>
      {snapshotStale && <p className="mt-1 text-xs">{t('banner.staleSnapshot', 'Stale snapshot')}</p>}
      {streamRecovery?.requestId && <p className="mt-1 text-xs">{t('errors.requestId', 'Request ID')}: {streamRecovery.requestId}</p>}
      {!live && <div className="mt-3 flex flex-wrap gap-2">
        {onReconnect && <button type="button" onClick={onReconnect} className="rounded-lg border border-amber-200/30 bg-slate-950/30 px-3 py-1.5 text-xs font-medium focus-visible:outline focus-visible:outline-amber-200">
          {t('errors.reconnectStream', 'Reconnect live stream')}
        </button>}
        {onRetry && <button type="button" onClick={() => void onRetry()} className="rounded-lg border border-amber-200/30 bg-slate-950/30 px-3 py-1.5 text-xs font-medium focus-visible:outline focus-visible:outline-amber-200">
          {t('errors.retry', 'Retry snapshot')}
        </button>}
        {navigateTo && <button type="button" onClick={() => navigateTo('/studio/v2')} className="rounded-lg border border-amber-200/30 bg-slate-950/30 px-3 py-1.5 text-xs font-medium focus-visible:outline focus-visible:outline-amber-200">
          {t('header.unavailableAction', 'Review setup in Studio V2')}
        </button>}
      </div>}
    </div>
  );
}
