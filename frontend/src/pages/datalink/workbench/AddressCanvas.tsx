import { useTranslation } from 'react-i18next';
import type { AddressCanvasItem } from './sourceCanvasModel';

const statusClassName: Record<AddressCanvasItem['status'], string> = {
  conflict: 'border-rose-500/50 bg-rose-900/30 text-rose-100',
  planned: 'border-sky-700/40 bg-sky-900/20 text-sky-100',
  used: 'border-emerald-700/40 bg-emerald-900/20 text-emerald-100',
};

export function AddressCanvas({ items }: { items: AddressCanvasItem[] }) {
  const { t } = useTranslation();

  if (items.length === 0) {
    return (
      <div className="rounded-2xl border border-dashed border-slate-700 bg-slate-950/40 px-4 py-10 text-center text-sm text-slate-400">
        {t('workbench.source.canvas.empty')}
      </div>
    );
  }

  return (
    <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-5">
      {items.map((item) => (
        <div
          key={item.address}
          data-testid={`address-cell-${item.address}`}
          data-status={item.status}
          className={`rounded-2xl border p-4 ${statusClassName[item.status]}`}
        >
          <p className="font-mono text-xs uppercase tracking-[0.18em] text-slate-300">
            {item.address}
          </p>
          <div className="mt-3 space-y-1">
            <p className="text-sm font-semibold text-white">
              {item.status === 'conflict'
                ? t('workbench.source.canvas.conflictLabel')
                : item.point?.name ?? t('workbench.source.canvas.plannedLabel')}
            </p>
            <p className="text-xs text-slate-300">
              {item.status === 'used'
                ? t('workbench.source.canvas.usedLabel')
                : item.status === 'conflict'
                  ? t('workbench.source.canvas.conflictState')
                  : t('workbench.source.canvas.plannedState')}
            </p>
            {item.point ? (
              <p className="text-xs text-cyan-100">
                {t('workbench.runtime.currentValue')}:{' '}
                {item.point.last_value === null || item.point.last_value === undefined
                  ? '—'
                  : String(item.point.last_value)}
              </p>
            ) : null}
          </div>
        </div>
      ))}
    </div>
  );
}
