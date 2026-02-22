import type { TFunction } from 'i18next';
import type { FlowState, FlowSegment } from '../../../features/flow/stateMachine';

/** 流程區段診斷品質對應的 Tailwind 背景色 */
const QUALITY_DOT_CLASS: Record<string, string> = {
  good: 'bg-emerald-400',
  warning: 'bg-yellow-400',
  bad: 'bg-red-400',
  unknown: 'bg-slate-500',
} as const;

const EMPTY_VALUE = '-';

interface SmartDashboardFlowStatusSectionProps {
  flowSegments: readonly FlowSegment[];
  flowState: FlowState;
  statusStyle: Record<string, string>;
  hasError: boolean;
  t: TFunction;
}

/**
 * 來源資料流程狀態區塊：顯示流程標題、整體狀態與各區段（Source / Grid / Tag / Sink）診斷。
 * 空值以 i18n 佔位符顯示，品質點具備無障礙標籤。
 */
export default function SmartDashboardFlowStatusSection({
  flowSegments,
  flowState,
  statusStyle,
  hasError,
  t,
}: SmartDashboardFlowStatusSectionProps) {
  return (
    <section className="mb-6 rounded-2xl border border-white/10 bg-slate-900/50 p-4" aria-label={t('smartDashboard.flowTitle')}>
      <div className="mb-3 flex flex-wrap items-center justify-between gap-2">
        <h3 className="text-sm font-semibold text-slate-100">{t('smartDashboard.flowTitle')}</h3>
        <span
          className={`rounded-lg border px-2 py-1 text-xs ${statusStyle[flowState.status] ?? ''}`}
          aria-label={t('smartDashboard.flowStatus.' + flowState.status)}
        >
          {t(`smartDashboard.flowStatus.${flowState.status}`)}
        </span>
      </div>
      {hasError && (
        <p className="mb-3 rounded-lg border border-red-500/40 bg-red-500/10 px-3 py-2 text-xs text-red-200" role="alert">
          {t('smartDashboard.flowErrorHint')}
        </p>
      )}
      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-4">
        {flowSegments.map((segment) => {
          const diag = flowState.diagnostics[segment];
          const hasValue = diag.latestValue !== EMPTY_VALUE && diag.latestValue.trim() !== '';
          const hasTime = diag.timestamp !== EMPTY_VALUE && diag.timestamp.trim() !== '';
          const qualityKey = diag.quality in QUALITY_DOT_CLASS ? diag.quality : 'unknown';
          const qualityLabel = t(`smartDashboard.flowSegmentQuality.${qualityKey}`);
          const segmentTitle = t(`smartDashboard.flowSegments.${segment}.title`);
          const segmentSubtitle = t(`smartDashboard.flowSegments.${segment}.subtitle`);
          return (
            <article
              key={segment}
              className="flex flex-col rounded-xl border border-white/10 bg-slate-800/40 p-3"
              aria-labelledby={`flow-segment-${segment}-title`}
            >
              <div className="flex items-center justify-between gap-2">
                <p id={`flow-segment-${segment}-title`} className="text-xs font-medium uppercase tracking-wider text-slate-400">
                  {segmentTitle}
                </p>
                <span
                  className={`h-2 w-2 shrink-0 rounded-full ${QUALITY_DOT_CLASS[qualityKey] ?? QUALITY_DOT_CLASS.unknown}`}
                  aria-label={t('smartDashboard.flowSegmentQuality.label', { quality: qualityLabel })}
                  title={qualityLabel}
                />
              </div>
              <p className="mt-0.5 text-[11px] text-slate-500">{segmentSubtitle}</p>
              <p className="mt-2 min-h-[1.25rem] truncate font-mono text-sm text-slate-100">
                {hasValue ? diag.latestValue : t('smartDashboard.flowSegmentNotSet')}
              </p>
              {hasTime && (
                <p className="mt-0.5 truncate text-[11px] text-slate-400">{diag.timestamp}</p>
              )}
              {diag.error && (
                <p className="mt-1 truncate text-[11px] text-red-300" role="status">
                  {diag.error}
                </p>
              )}
            </article>
          );
        })}
      </div>
    </section>
  );
}
