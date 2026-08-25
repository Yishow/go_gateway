import { useTranslation } from 'react-i18next';
import type { RuntimeDashboardRouteState } from './useRuntimeDashboardState';
import type { RuntimeStreamRecovery } from '../../../types/datalink';

interface LiveStateBannerProps {
  routeState: RuntimeDashboardRouteState;
  streamRecovery?: RuntimeStreamRecovery | null;
  onRetry?: () => Promise<unknown>;
  onReconnect?: () => void;
  navigateTo?: (target: string) => void;
}

export function LiveStateBanner({ routeState, streamRecovery, onRetry, onReconnect, navigateTo }: LiveStateBannerProps) {
  const { t } = useTranslation('runtime-dashboard');

  if (routeState === 'degraded') {
    return (
      <div
        className="rounded-2xl border border-amber-500/30 bg-amber-500/10 px-4 py-3 text-sm text-amber-100"
        data-testid="runtime-dashboard-live-state-banner"
      >
        <div>
          {t(
          'banner.degraded',
          'Live stream degraded. Showing the last successful runtime snapshot until the stream recovers.',
          )}
          {streamRecovery?.requestId ? (
            <div className="mt-1 text-xs text-amber-100/80">
              {t('errors.requestId', 'Request ID')}: {streamRecovery.requestId}
            </div>
          ) : null}
          <div className="mt-3 flex flex-wrap gap-2">
            {streamRecovery?.code === 'runtime_stream_unavailable' && streamRecovery.retryable && onReconnect ? (
              <button
                type="button"
                onClick={onReconnect}
                className="rounded-lg border border-amber-200/30 bg-slate-950/30 px-3 py-1.5 text-xs font-medium text-amber-50"
              >
                {t('errors.reconnectStream', 'Reconnect live stream')}
              </button>
            ) : onRetry ? (
              <button
                type="button"
                onClick={() => void onRetry()}
                className="rounded-lg border border-amber-200/30 bg-slate-950/30 px-3 py-1.5 text-xs font-medium text-amber-50"
              >
                {t('errors.retry', 'Retry snapshot')}
              </button>
            ) : null}
            {navigateTo ? (
              <button
                type="button"
                onClick={() => navigateTo('/studio/v2')}
                className="rounded-lg border border-amber-200/30 bg-slate-950/30 px-3 py-1.5 text-xs font-medium text-amber-50"
              >
                {t('header.unavailableAction', 'Review setup in Studio V2')}
              </button>
            ) : null}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div
      className="rounded-2xl border border-emerald-500/30 bg-emerald-500/10 px-4 py-3 text-sm text-emerald-100"
      data-testid="runtime-dashboard-live-state-banner"
    >
      {t('banner.live', 'Live monitoring active')}
    </div>
  );
}
