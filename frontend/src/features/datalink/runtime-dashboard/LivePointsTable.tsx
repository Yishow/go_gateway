import { useTranslation } from 'react-i18next';
import type { RuntimeValueEvent } from '../../../types/datalink';

interface LivePointsTableProps {
  liveValues: Record<string, RuntimeValueEvent>;
}

function formatValue(value: unknown) {
  if (value === null || value === undefined) {
    return '-';
  }

  if (typeof value === 'object') {
    return JSON.stringify(value);
  }

  return String(value);
}

export function LivePointsTable({ liveValues }: LivePointsTableProps) {
  const { t } = useTranslation('runtime-dashboard');
  const rows = Object.values(liveValues).sort((left, right) =>
    left.address.localeCompare(right.address),
  );

  return (
    <section
      className="rounded-3xl border border-slate-800 bg-slate-900/70 p-6"
      data-testid="runtime-dashboard-live-points-table"
    >
      <div className="mb-4">
        <h2 className="text-lg font-semibold text-slate-50">
          {t('points.title', 'Live points')}
        </h2>
        <p className="mt-1 text-sm text-slate-400">
          {t('points.description', 'Live point rows update after the first runtime value event arrives.')}
        </p>
      </div>

      {rows.length === 0 ? (
        <div
          className="rounded-2xl border border-dashed border-slate-800 bg-slate-950/40 p-8 flex flex-col items-center justify-center space-y-4"
          data-testid="runtime-dashboard-live-points-placeholder"
        >
          <div className="flex items-center gap-2 text-cyan-400/90 text-sm font-medium">
            <span className="relative flex h-2.5 w-2.5">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-cyan-400 opacity-75"></span>
              <span className="relative inline-flex rounded-full h-2.5 w-2.5 bg-cyan-500"></span>
            </span>
            {t('points.placeholder', 'Waiting for the first live point update...')}
          </div>
          <div className="w-full max-w-md space-y-2.5 animate-pulse">
            <div className="h-3.5 bg-slate-900/80 rounded-lg w-full"></div>
            <div className="h-3.5 bg-slate-900/80 rounded-lg w-11/12 mx-auto"></div>
            <div className="h-3.5 bg-slate-900/80 rounded-lg w-10/12 mx-auto"></div>
          </div>
        </div>
      ) : (
        <div className="overflow-x-auto">
          <table className="min-w-full text-left text-sm text-slate-200">
            <thead className="text-xs uppercase tracking-[0.18em] text-slate-500">
              <tr>
                <th className="px-3 py-2">{t('points.address', 'Address')}</th>
                <th className="px-3 py-2">{t('points.rawValue', 'Raw value')}</th>
                <th className="px-3 py-2">{t('points.transformedValue', 'Transformed value')}</th>
                <th className="px-3 py-2">{t('points.quality', 'Quality')}</th>
                <th className="px-3 py-2">{t('points.timestamp', 'Timestamp')}</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((row) => (
                <tr key={row.point_id} className="border-t border-slate-800">
                  <td className="px-3 py-3 font-mono text-cyan-100">{row.address}</td>
                  <td className="px-3 py-3">{formatValue(row.raw_value)}</td>
                  <td className="px-3 py-3">{formatValue(row.transformed_value)}</td>
                  <td className="px-3 py-3">
                    {row.stale ? t('points.stale', 'stale') : row.quality}
                  </td>
                  <td className="px-3 py-3">{row.timestamp}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}
