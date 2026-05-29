import { useTranslation } from 'react-i18next';
import type { RuntimeDashboardRouteState } from './useRuntimeDashboardState';

interface LiveStateBannerProps {
  routeState: RuntimeDashboardRouteState;
}

export function LiveStateBanner({ routeState }: LiveStateBannerProps) {
  const { t } = useTranslation('runtime-dashboard');

  if (routeState === 'degraded') {
    return (
      <div
        className="rounded-2xl border border-amber-500/30 bg-amber-500/10 px-4 py-3 text-sm text-amber-100"
        data-testid="runtime-dashboard-live-state-banner"
      >
        {t(
          'banner.degraded',
          'Live stream degraded. Showing the last successful runtime snapshot until the stream recovers.',
        )}
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
