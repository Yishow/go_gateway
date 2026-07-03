import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import type { DbConnector, Point, Rule, Mapping, DbTarget } from '../../state/types';

interface DestinationOverviewCardProps {
  connector: DbConnector;
  rules: Rule[];
  points: Point[];
  mappings: Record<string, Mapping>;
  targets: Record<string, DbTarget>;
  hasConflict: boolean;
}

export function DestinationOverviewCard({
  connector,
  rules,
  points,
  mappings,
  targets,
  hasConflict,
}: DestinationOverviewCardProps) {
  const { t } = useTranslation('workbench-v2');

  const enabledRuleCount = useMemo(
    () => rules.filter((rule) => rule.enabled).length,
    [rules],
  );
  const mappedPointCount = useMemo(
    () => points.filter((point) => point.enabled && mappings[point.id]).length,
    [points, mappings],
  );
  const enabledTargetCount = useMemo(
    () => Object.values(targets).filter((target) => target.enabled).length,
    [targets],
  );

  const destination =
    connector.kind === 'sqlite' || !connector.schema
      ? connector.table
      : `${connector.schema}.${connector.table}`;

  const items = [
    { label: t('step4.overview_destination', '目標表'), value: destination, mono: true },
    { label: t('step4.overview_kind', '資料庫'), value: connector.kind },
    { label: t('step4.overview_write_mode', '寫入策略'), value: connector.write_mode },
    {
      label: t('step4.overview_interval', '寫入間隔'),
      value: t('step4.overview_interval_value', '{{seconds}} 秒', {
        seconds: connector.write_interval_seconds,
      }),
    },
    {
      label: t('step4.overview_scope', '規則 / 點位'),
      value: t('step4.overview_scope_value', '{{rules}} 組 / {{points}} 點', {
        rules: enabledRuleCount,
        points: mappedPointCount,
      }),
    },
    {
      label: t('step4.overview_targets', '已啟用欄位'),
      value: String(enabledTargetCount),
    },
  ];

  return (
    <section
      data-testid="step4-destination-overview"
      className="rounded-2xl border border-cyan-500/20 bg-cyan-500/5 p-4"
    >
      <div className="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
        <div className="space-y-1">
          <div className="text-xs font-semibold uppercase tracking-[0.18em] text-cyan-200/80">
            {t('step4.overview_title', '單表寫入目的地')}
          </div>
          <div className="text-sm text-slate-200">
            {t(
              'step4.overview_subtitle',
              'Step 4 目前會把所有已啟用 point 配到同一張目標表；不同 rule 可以共表，但不能共用同一欄位。',
            )}
          </div>
        </div>

        <div className="shrink-0 rounded-full border px-3 py-1 text-xs font-medium text-cyan-100">
          {hasConflict
            ? t('step4.overview_conflict', '欄位衝突待處理')
            : t('step4.overview_ready', '欄位配置可檢查')}
        </div>
      </div>

      <div className="mt-4 grid gap-2 md:grid-cols-2 xl:grid-cols-3">
        {items.map((item) => (
          <div key={item.label} className="rounded-xl border border-slate-800/70 bg-slate-950/50 px-3 py-2">
            <div className="text-[11px] text-slate-500">{item.label}</div>
            <div className={`mt-1 text-sm font-semibold text-slate-100 ${item.mono ? 'font-mono' : ''}`}>
              {item.value}
            </div>
          </div>
        ))}
      </div>
    </section>
  );
}
