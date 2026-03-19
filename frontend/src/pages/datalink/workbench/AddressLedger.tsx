import { useTranslation } from 'react-i18next';
import {
  formatSourceValue,
  type AddressCanvasItem,
  type SourceValueFormat,
  type SourceViewMode,
} from './sourceCanvasModel';

function getStatusLabel(item: AddressCanvasItem, t: (key: string) => string) {
  switch (item.status) {
    case 'conflict':
      return t('workbench.source.canvas.conflictState');
    case 'gap':
      return t('workbench.source.canvas.gapState');
    case 'planned':
      return t('workbench.source.canvas.plannedState');
    case 'unmanaged':
      return t('workbench.source.canvas.unmanagedState');
    case 'used':
      return t('workbench.source.canvas.usedState');
  }
}

function getNameLabel(item: AddressCanvasItem, t: (key: string) => string) {
  switch (item.status) {
    case 'conflict':
      return t('workbench.source.canvas.conflictLabel');
    case 'gap':
      return t('workbench.source.canvas.gapLabel');
    case 'planned':
      return item.tagDisplayName ?? item.tagKey ?? t('workbench.source.canvas.plannedLabel');
    case 'unmanaged':
      return item.tagDisplayName ?? item.tagKey ?? item.point?.name ?? t('workbench.source.canvas.unmanagedLabel');
    case 'used':
      return item.tagDisplayName ?? item.tagKey ?? item.point?.name ?? t('workbench.source.canvas.usedLabel');
  }
}

function getValueLabel(input: {
  item: AddressCanvasItem;
  t: (key: string) => string;
  valueFormat: SourceValueFormat;
  viewMode: SourceViewMode;
}) {
  const { item, t, valueFormat, viewMode } = input;
  switch (viewMode) {
    case 'plan':
      return item.primaryRuleId ?? '—';
    case 'live':
      return formatSourceValue(item.liveValue, valueFormat);
    case 'link':
      return item.linkLabelKey ? t(item.linkLabelKey) : t('workbench.source.link.noState');
  }
}

export function AddressLedger({
  items,
  valueFormat,
  viewMode,
}: {
  items: AddressCanvasItem[];
  valueFormat: SourceValueFormat;
  viewMode: SourceViewMode;
}) {
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
            <th className="px-4 py-3 text-left">{t('workbench.runtime.currentValue')}</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-800">
          {items.map((item) => (
            <tr key={item.address} data-testid={`address-ledger-row-${item.address}`}>
              <td className="px-4 py-3 font-mono text-xs text-slate-300">{item.address}</td>
              <td className="px-4 py-3">
                {getNameLabel(item, t)}
              </td>
              <td className="px-4 py-3">{getStatusLabel(item, t)}</td>
              <td className="px-4 py-3 text-cyan-100">
                {getValueLabel({ item, t, valueFormat, viewMode })}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
