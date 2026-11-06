import { useTranslation } from 'react-i18next';
import { activeAnomalyCount } from '../../../../services/acquisitionAnomalies';
import type { AcquisitionAnomalySummary } from '../../../../types/acquisitionAnomaly';
import './acquisitionAnomalies.css';

/** Home-level view of the anomaly summary read; `unavailable` covers failed or missing reads. */
export type AnomalySummaryStatus =
  | { state: 'checking' }
  | { state: 'unavailable' }
  | { state: 'ready'; summary: AcquisitionAnomalySummary };

/** Fails closed: an error hides older counts, because a stale count is not the current status. */
export function anomalySummaryStatus(query: { isError: boolean; data: AcquisitionAnomalySummary | undefined }): AnomalySummaryStatus {
  if (query.isError) return { state: 'unavailable' };
  return query.data ? { state: 'ready', summary: query.data } : { state: 'checking' };
}

type ButtonState = 'none' | 'active' | 'unknown' | 'checking';

function buttonState(status: AnomalySummaryStatus | undefined, deviceId: string): { state: ButtonState; count: number } {
  if (!status || status.state === 'unavailable') return { state: 'unknown', count: 0 };
  if (status.state === 'checking') return { state: 'checking', count: 0 };
  const count = activeAnomalyCount(status.summary, deviceId);
  if (count === null) return { state: 'unknown', count: 0 };
  return { state: count > 0 ? 'active' : 'none', count };
}

/** Sits next to "read points"; warning tokens and a text count only when active events exist. */
export function AcquisitionAnomalyButton({ deviceId, status, onOpen }: { deviceId: string; status: AnomalySummaryStatus | undefined; onOpen: () => void }) {
  const { t } = useTranslation('recording-home');
  const { state, count } = buttonState(status, deviceId);
  const label = state === 'active' ? t('acquisitionAnomalies.buttonActive', { count })
    : state === 'unknown' ? t('acquisitionAnomalies.buttonUnknown')
      : state === 'checking' ? t('acquisitionAnomalies.buttonChecking')
        : t('acquisitionAnomalies.button');
  return <button type="button" className="gw-anomaly-button" data-anomaly-state={state} data-testid={`home-device-anomalies-${deviceId}`}
    aria-busy={state === 'checking' || undefined} onClick={onOpen}>{label}</button>;
}
