import { useTranslation } from 'react-i18next';
import type { AcquisitionPreviewGroup as Group } from '../../../../types/acquisitionPlanPreview';

/** Counts are rendered in their native units; never sum wire subsets or different cadence groups. */
export function AcquisitionPreviewGroup({ group }: { group: Group }) {
  const { t } = useTranslation('acquisition-preview');
  const id = group.polling_group_id;
  const actual = group.actual.status === 'available' && group.actual.initial_plan_revision === group.plan_revision &&
    group.actual.declared_plan_revision === group.declared_plan_revision ? group.actual : undefined;
  return <section data-testid={`acquisition-group-${id}`} className="space-y-2 border-t border-slate-700/60 pt-3">
    <h4>{t('group')} <span className="font-mono break-all">{id}</span></h4>
    <dl className="grid grid-cols-1 gap-2 sm:grid-cols-2">
      <div><dt>{t('cadence')}</dt><dd data-testid={`acquisition-cadence-${id}`} className="tabular-nums">{group.poll_interval_ms} ms</dd></div>
      <div><dt>{t('points')}</dt><dd data-testid={`acquisition-points-${id}`} className="tabular-nums">{group.acquisition_point_count}</dd></div>
      <div><dt>{t('planned')}</dt><dd data-testid={`acquisition-planned-${id}`} className="tabular-nums">{group.planned_requests}</dd></div>
    </dl>
    <p className="text-slate-400">{t(`basis.${group.basis}`)}</p>
    <h5>{t('wireUnitsTitle')}</h5>
    <ul>{group.unique_wire_units.map((unit, index) => <li key={`${unit.area}:${unit.function}:${unit.address_unit_bits}:${index}`}>
      {t(unit.address_unit_bits === 1 ? 'wireUnits.bit' : 'wireUnits.word', { area: unit.area, function: unit.function, count: unit.unique_address_units })}
    </li>)}</ul>
    <details data-testid={`acquisition-limits-${id}`}><summary>{t('limits')}</summary>
      {group.request_limits?.length ? <ul>{group.request_limits.map((limit, index) => <li key={index} className="space-y-1">
        <p>{t('limitRow', { kind: limit.kind, area: limit.area, function: limit.function, limit: limit.native_unit_limit })}</p>
        <p className="text-slate-400">{limit.frame_cost_known ? t('frameCost', { request: limit.request_frame_bytes, response: limit.response_frame_bytes }) : t('frameCostUnavailable')}</p>
      </li>)}</ul> : <p>{t('limitsUnavailable')}</p>}
    </details>
    {group.point_errors.length > 0 && <details><summary>{t('pointErrors')}</summary><ul>{group.point_errors.map((error, index) =>
      <li key={`${error.point_id}:${error.code}:${index}`}><span className="font-mono break-all">{error.point_id}</span> · <code className="break-all">{error.code}</code></li>)}</ul></details>}
    {group.fallback_reasons.length > 0 && <details><summary>{t('fallbackReasons')}</summary><ul>{group.fallback_reasons.map((reason, index) =>
      <li key={`${reason}:${index}`}><code className="break-all">{reason}</code></li>)}</ul></details>}
    <h5>{t('actualTitle')}</h5>
    {!actual ? <p className="text-slate-400">{t('actualUnavailable')}</p> : <>
      <p>{t('observedAt', { at: actual.observed_at })}</p>
      <dl className="grid grid-cols-1 gap-2 sm:grid-cols-2">
        <div><dt>{t('wireAttempts')}</dt><dd data-testid={`acquisition-wire-attempts-${id}`} className="tabular-nums">{actual.wire_attempts}</dd></div>
        <div><dt>{t('retry')}</dt><dd data-testid={`acquisition-retry-${id}`} className="tabular-nums">{actual.retry_wire_attempts}</dd></div>
        <div><dt>{t('fallback')}</dt><dd data-testid={`acquisition-fallback-${id}`} className="tabular-nums">{actual.fallback_wire_attempts}</dd></div>
        <div><dt>{t('latency')}</dt><dd className="tabular-nums">{actual.cycle_latency_ms} ms</dd></div>
      </dl>
      <p className="text-slate-400">{t('wireSubsetNote')}</p>
      {actual.logical_requests.status === 'observed'
        ? <p data-testid={`acquisition-logical-${id}`}>{t('logicalCounts', { succeeded: actual.logical_requests.succeeded,
          failed: actual.logical_requests.failed, canceled: actual.logical_requests.canceled })}</p>
        : <p>{t('logicalNotStarted')}</p>}
      <p>{t(actual.overrun ? 'overrun' : 'withinCadence')}</p>
      <p>{actual.last_good_age_ms === null ? t('lastGoodUnavailable') : t('lastGoodAge', { age: actual.last_good_age_ms })}</p>
    </>}
  </section>;
}
