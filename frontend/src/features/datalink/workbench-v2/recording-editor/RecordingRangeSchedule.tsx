import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { RangeSourceOutcome } from '../../../../pages/datalink/workbench-v2/rangeSourceSubmission';
import { readRecordingRangeSchedule, type RecordingRangeScheduleRow } from './recordingRangeSchedule';
/** One scoped display read, independent from the submission and its completion. */
export function RecordingRangeSchedule({ workspaceId, deviceId, outcome }: { workspaceId: string; deviceId?: string; outcome?: RangeSourceOutcome }) {
  const { t } = useTranslation('recording-editor');
  const [refresh, setRefresh] = useState(0);
  const [observed, setObserved] = useState<{ outcome: RangeSourceOutcome; refresh: number; rows?: RecordingRangeScheduleRow[] }>();
  const canRead = Boolean(outcome?.savedPoints?.length && outcome.request.workspaceId === workspaceId && outcome.request.deviceId === deviceId);
  useEffect(() => {
    if (!canRead || !outcome || !deviceId) return;
    let current = true;
    void readRecordingRangeSchedule(deviceId, outcome.savedPoints!).then(rows => {
      if (current) setObserved({ outcome, refresh, rows });
    }, () => { if (current) setObserved({ outcome, refresh }); });
    return () => { current = false; };
  }, [canRead, deviceId, outcome, refresh, workspaceId]);
  const snapshot = canRead && observed && observed.outcome === outcome && observed.refresh === refresh ? observed : undefined;
  return <section data-testid="recording-range-schedule" aria-label={t('range.schedule')} className="space-y-2">
    <h3>{t('range.schedule')}</h3>
    {!canRead ? <p className="gw-muted">{t('range.scheduleUnassigned')}</p> : !snapshot ? <p role="status">{t('range.scheduleLoading')}</p> : snapshot.rows ? <>
      <p className="gw-muted">{t('range.scheduleObserved')}</p>
      <ul>{snapshot.rows.map(row => <li key={row.id}>{t('range.scheduleRow', { name: row.name, interval: row.intervalMs, count: row.points })} · {t(row.enabled ? 'range.scheduleEnabled' : 'range.scheduleDisabled')}</li>)}</ul>
    </> : <p role="status">{t('range.scheduleUnavailable')}</p>}
    {canRead && snapshot && <button type="button" onClick={() => setRefresh(value => value + 1)}>{t('range.scheduleRefresh')}</button>}
    <a href="/studio/v2" target="_blank" rel="noopener noreferrer">{t('range.globalScheduleSettings')}</a>
    <p className="gw-muted">{t('range.globalScheduleNote')}</p>
  </section>;
}
