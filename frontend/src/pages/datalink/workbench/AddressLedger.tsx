import { useTranslation } from 'react-i18next';
import type { AddressCanvasItem } from './sourceCanvasModel';

export function AddressLedger({ items }: { items: AddressCanvasItem[] }) {
  const { t } = useTranslation();

  if (items.length === 0) {
    return (
      <div className="rounded-2xl border border-dashed border-slate-700 bg-slate-950/40 px-4 py-10 text-center text-sm text-slate-400">
        {t('workbench.source.ledger.empty')}
      </div>
    );
  }

  return (
    <div className="overflow-hidden rounded-2xl border border-slate-800 bg-slate-950/40">
      <table className="min-w-full divide-y divide-slate-800 text-sm text-slate-200">
        <thead className="bg-slate-900/80 text-xs uppercase tracking-[0.16em] text-slate-400">
          <tr>
            <th className="px-4 py-3 text-left">{t('workbench.source.ledger.columns.address')}</th>
            <th className="px-4 py-3 text-left">{t('workbench.source.ledger.columns.name')}</th>
            <th className="px-4 py-3 text-left">{t('workbench.source.ledger.columns.status')}</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-800">
          {items.map((item) => (
            <tr key={item.address} data-testid={`address-ledger-row-${item.address}`}>
              <td className="px-4 py-3 font-mono text-xs text-slate-300">{item.address}</td>
              <td className="px-4 py-3">
                {item.status === 'conflict'
                  ? t('workbench.source.canvas.conflictLabel')
                  : item.point?.name ?? t('workbench.source.canvas.plannedLabel')}
              </td>
              <td className="px-4 py-3">
                {item.status === 'used'
                  ? t('workbench.source.canvas.usedLabel')
                  : item.status === 'conflict'
                    ? t('workbench.source.canvas.conflictState')
                    : t('workbench.source.canvas.plannedState')}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
