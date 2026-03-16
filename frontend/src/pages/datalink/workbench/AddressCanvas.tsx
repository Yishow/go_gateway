import { useTranslation } from 'react-i18next';
import {
  formatSourceValue,
  type AddressCanvasItem,
  type SourceValueFormat,
  type SourceViewMode,
} from './sourceCanvasModel';

const statusClassName: Record<AddressCanvasItem['status'], string> = {
  conflict: 'border-rose-500/50 bg-rose-900/30 text-rose-100',
  gap: 'border-slate-800 bg-slate-950/40 text-slate-500',
  planned: 'border-sky-700/40 bg-sky-900/20 text-sky-100',
  used: 'border-emerald-700/40 bg-emerald-900/20 text-emerald-100',
};

function getPrimaryLabel(item: AddressCanvasItem, t: (key: string) => string) {
  switch (item.status) {
    case 'conflict':
      return t('workbench.source.canvas.conflictLabel');
    case 'gap':
      return t('workbench.source.canvas.gapLabel');
    case 'planned':
      return item.point?.name ?? t('workbench.source.canvas.plannedLabel');
    case 'used':
      return item.point?.name ?? t('workbench.source.canvas.usedLabel');
  }
}

function getPlanStateLabel(item: AddressCanvasItem, t: (key: string) => string) {
  switch (item.status) {
    case 'conflict':
      return t('workbench.source.canvas.conflictState');
    case 'gap':
      return t('workbench.source.canvas.gapState');
    case 'planned':
      return t('workbench.source.canvas.plannedState');
    case 'used':
      return t('workbench.source.canvas.usedLabel');
  }
}

function getOverlayValue(input: {
  item: AddressCanvasItem;
  t: (key: string) => string;
  valueFormat: SourceValueFormat;
  viewMode: SourceViewMode;
}) {
  const { item, t, valueFormat, viewMode } = input;

  switch (viewMode) {
    case 'plan':
      return getPlanStateLabel(item, t);
    case 'live':
      return formatSourceValue(item.liveValue, valueFormat);
    case 'link':
      return item.linkLabelKey ? t(item.linkLabelKey) : t('workbench.source.link.noState');
  }
}

export function AddressCanvas({
  items,
  onSelectAddress,
  selectedAddress,
  valueFormat,
  viewMode,
}: {
  items: AddressCanvasItem[];
  onSelectAddress?: (address: string) => void;
  selectedAddress?: string | null;
  valueFormat: SourceValueFormat;
  viewMode: SourceViewMode;
}) {
  const { t } = useTranslation();

  if (items.length === 0) {
    return (
      <div className="rounded-2xl border border-dashed border-slate-700 bg-slate-950/40 px-4 py-10 text-center text-sm text-slate-400">
        {t('workbench.source.canvas.empty')}
      </div>
    );
  }

  return (
    <div
      className="grid gap-3 sm:grid-cols-2 xl:grid-cols-5 2xl:grid-cols-6"
      data-testid="source-canvas"
      data-view-mode={viewMode}
    >
      {items.map((item) => (
        <button
          key={item.address}
          data-testid={`address-cell-${item.address}`}
          data-status={item.status}
          className={[
            'rounded-2xl border p-4 text-left transition',
            statusClassName[item.status],
            selectedAddress === item.address ? 'ring-2 ring-cyan-400/60' : '',
          ].join(' ')}
          onClick={() => onSelectAddress?.(item.address)}
          type="button"
        >
          <p className="font-mono text-xs uppercase tracking-[0.18em] text-slate-300">
            {item.address}
          </p>
          <div className="mt-3 space-y-1">
            <p className="text-sm font-semibold text-white">
              {getPrimaryLabel(item, t)}
            </p>
            <p className="text-xs text-slate-300">
              {getOverlayValue({ item, t, valueFormat, viewMode })}
            </p>
            {viewMode === 'plan' && item.mergeSpan > 1 ? (
              <p className="text-xs text-cyan-100">
                {t('workbench.source.canvas.mergeState', { cells: item.mergeSpan })}
              </p>
            ) : null}
            {viewMode === 'live' && item.liveTimestamp ? (
              <p className="text-xs text-slate-400">{item.liveTimestamp}</p>
            ) : null}
            {item.primaryRuleId ? (
              <p className="text-[11px] text-slate-400">{item.primaryRuleId}</p>
            ) : null}
          </div>
        </button>
      ))}
    </div>
  );
}
