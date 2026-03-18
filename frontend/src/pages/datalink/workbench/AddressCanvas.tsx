import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import {
  formatSourceValue,
  type AddressCanvasItem,
  type SourceValueFormat,
  type SourceViewMode,
} from './sourceCanvasModel';

const LATTICE_COLUMNS = 16;

const statusClassName: Record<AddressCanvasItem['status'], string> = {
  conflict: 'border-rose-500/50 bg-rose-900/30 text-rose-100',
  gap: 'border-slate-800 bg-slate-950/40 text-slate-500',
  planned: 'border-sky-700/40 bg-sky-900/20 text-sky-100',
  used: 'border-emerald-700/40 bg-emerald-900/20 text-emerald-100',
};

type CanvasRowItem =
  | { kind: 'cell'; item: AddressCanvasItem; colSpan: number }
  | { kind: 'hidden'; item: AddressCanvasItem };

function buildRowLayout(
  rowItems: ReadonlyArray<AddressCanvasItem>,
): CanvasRowItem[] {
  const layout: CanvasRowItem[] = [];

  for (const item of rowItems) {
    if (item.mergeSpan > 1 && item.mergeOffset > 0) {
      layout.push({ kind: 'hidden', item });
      continue;
    }

    const remainingInRow =
      rowItems.length - rowItems.indexOf(item);
    const colSpan = Math.min(item.mergeSpan, remainingInRow);
    layout.push({ kind: 'cell', item, colSpan });
  }

  return layout;
}

function chunkItems(items: ReadonlyArray<AddressCanvasItem>, size: number) {
  const rows: AddressCanvasItem[][] = [];

  for (let index = 0; index < items.length; index += size) {
    rows.push(items.slice(index, index + size));
  }

  return rows;
}

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
  const rows = chunkItems(items, LATTICE_COLUMNS);

  const selectedLogicalAddresses = useMemo(() => {
    if (!selectedAddress) return new Set<string>();
    const rootIndex = items.findIndex((candidate) => candidate.address === selectedAddress);
    if (rootIndex === -1) return new Set<string>();
    const rootItem = items[rootIndex];
    if (rootItem.mergeSpan <= 1) return new Set([selectedAddress]);

    const addresses = new Set<string>();
    for (let offset = 0; offset < rootItem.mergeSpan && rootIndex + offset < items.length; offset++) {
      addresses.add(items[rootIndex + offset].address);
    }
    return addresses;
  }, [items, selectedAddress]);

  if (items.length === 0) {
    return (
      <div className="rounded-2xl border border-dashed border-slate-700 bg-slate-950/40 px-4 py-10 text-center text-sm text-slate-400">
        {t('workbench.source.canvas.empty')}
      </div>
    );
  }

  return (
    <div
      className="space-y-3"
      data-testid="source-canvas"
      data-view-mode={viewMode}
    >
      {rows.map((rowItems, rowIndex) => {
        const layout = buildRowLayout(rowItems);
        return (
          <div
            key={`source-canvas-row-${rowIndex}`}
            className="grid gap-0.5 rounded-2xl border border-slate-800 bg-slate-900/70 p-1"
            data-lattice-columns={String(LATTICE_COLUMNS)}
            data-row-end-address={rowItems.at(-1)?.address ?? ''}
            data-row-start-address={rowItems[0]?.address ?? ''}
            data-testid={`source-canvas-row-${rowIndex}`}
            style={{ gridTemplateColumns: `repeat(${LATTICE_COLUMNS}, minmax(0, 1fr))` }}
          >
            {layout.map((entry) => {
              if (entry.kind === 'hidden') {
                return (
                  <button
                    aria-hidden
                    aria-pressed={selectedLogicalAddresses.has(entry.item.address)}
                    key={entry.item.address}
                    data-merge-offset={String(entry.item.mergeOffset)}
                    data-merge-span={String(entry.item.mergeSpan)}
                    data-status={entry.item.status}
                    data-testid={`address-cell-${entry.item.address}`}
                    className="sr-only"
                    onClick={() => onSelectAddress?.(entry.item.address)}
                    type="button"
                    tabIndex={-1}
                  >
                    {entry.item.address}
                  </button>
                );
              }

              const { item, colSpan } = entry;
              return (
                <button
                  aria-pressed={selectedLogicalAddresses.has(item.address)}
                  key={item.address}
                  data-merge-offset={String(item.mergeOffset)}
                  data-merge-span={String(item.mergeSpan)}
                  data-status={item.status}
                  data-testid={`address-cell-${item.address}`}
                  className={[
                    'min-h-28 rounded-xl border p-3 text-left transition focus:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400/60',
                    statusClassName[item.status],
                    selectedLogicalAddresses.has(item.address) ? 'ring-2 ring-cyan-400/60' : '',
                  ].join(' ')}
                  onClick={() => onSelectAddress?.(item.address)}
                  style={colSpan > 1 ? { gridColumn: `span ${colSpan}` } : undefined}
                  type="button"
                >
                  <p className="font-mono text-[11px] uppercase tracking-[0.16em] text-slate-300">
                    {item.address}
                  </p>
                  <div className="mt-2 space-y-1">
                    <p className="text-xs font-semibold text-white">
                      {getPrimaryLabel(item, t)}
                    </p>
                    <p className="text-[11px] text-slate-300">
                      {getOverlayValue({ item, t, valueFormat, viewMode })}
                    </p>
                    {viewMode === 'plan' && colSpan > 1 ? (
                      <p className="text-[11px] text-cyan-100">
                        {t('workbench.source.canvas.mergeState', { cells: item.mergeSpan })}
                      </p>
                    ) : null}
                    {viewMode === 'live' && item.liveTimestamp ? (
                      <p className="text-[11px] text-slate-400">{item.liveTimestamp}</p>
                    ) : null}
                    {item.primaryRuleId ? (
                      <p className="text-[11px] text-slate-400">{item.primaryRuleId}</p>
                    ) : null}
                  </div>
                </button>
              );
            })}
          </div>
        );
      })}
    </div>
  );
}
