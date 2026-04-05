import { useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import {
  formatSourceValue,
  type AddressCanvasItem,
  type SourceValueFormat,
  type SourceViewMode,
} from './sourceCanvasModel';
import {
  clampLatticeColumns,
  LATTICE_COLUMNS_DEFAULT,
} from './sourceCanvasLatticeColumns';

/** 預設每列格數（與 {@link LATTICE_COLUMNS_DEFAULT} 相同，供外層對齊慣例）。 */
export const ADDRESS_CANVAS_LATTICE_COLUMNS = LATTICE_COLUMNS_DEFAULT;

/** 依狀態套用格位底色與邊框色；`used` 之語意主要由此視覺區分。 */
const statusClassName: Record<AddressCanvasItem['status'], string> = {
  conflict: 'border-rose-500/50 bg-rose-900/30 text-rose-100',
  gap: 'border-slate-800 bg-slate-950/40 text-slate-500',
  planned: 'border-sky-700/40 bg-sky-900/20 text-sky-100',
  unmanaged: 'border-amber-700/40 bg-amber-900/20 text-amber-100',
  used: 'border-emerald-700/40 bg-emerald-900/20 text-emerald-100',
};

type CanvasRowItem =
  | { kind: 'cell'; item: AddressCanvasItem; colSpan: number }
  | { kind: 'hidden'; item: AddressCanvasItem };

/**
 * 將同一列地址項目轉成可置於 CSS grid 的 cell／隱藏延續格配置。
 */
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

/**
 * 將平面地址清單切成每列固定欄數的列，供畫布網格使用。
 */
function chunkItems(items: ReadonlyArray<AddressCanvasItem>, size: number) {
  const rows: AddressCanvasItem[][] = [];

  for (let index = 0; index < items.length; index += size) {
    rows.push(items.slice(index, index + size));
  }

  return rows;
}

/**
 * 格位主標題：依狀態顯示衝突／空隙說明或 Tag／Point 名稱。
 */
function getPrimaryLabel(item: AddressCanvasItem, t: (key: string) => string) {
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

/** plan 檢視下覆蓋層的狀態短句（`used` 改在 getOverlayValue 顯示原始值）。 */
function getPlanStateLabel(item: AddressCanvasItem, t: (key: string) => string) {
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

/**
 * 依檢視模式決定格位覆蓋層文字；plan + `used` 時顯示格式化後的原始值（與底色語意不重複）。
 */
function getOverlayValue(input: {
  item: AddressCanvasItem;
  t: (key: string) => string;
  valueFormat: SourceValueFormat;
  viewMode: SourceViewMode;
}) {
  const { item, t, valueFormat, viewMode } = input;

  switch (viewMode) {
    case 'plan':
      if (item.status === 'used') {
        return formatSourceValue(item.liveValue, valueFormat);
      }
      return getPlanStateLabel(item, t);
    case 'live':
      return formatSourceValue(item.liveValue, valueFormat);
    case 'link':
      return item.linkLabelKey ? t(item.linkLabelKey) : t('workbench.source.link.noState');
  }
}

/**
 * 來源規則地址格狀畫布：以狀態色塊與多行文字呈現每個邏輯位址，支援 plan／live／link 三種覆蓋層。
 */
export function AddressCanvas({
  items,
  onSelectAddress,
  selectedAddress,
  valueFormat,
  viewMode,
  latticeColumns,
}: {
  items: AddressCanvasItem[];
  onSelectAddress?: (address: string) => void;
  selectedAddress?: string | null;
  valueFormat: SourceValueFormat;
  viewMode: SourceViewMode;
  /** 每列顯示幾個位址格；未傳則使用 {@link ADDRESS_CANVAS_LATTICE_COLUMNS}。 */
  latticeColumns?: number;
}) {
  const { t } = useTranslation();
  const columnCount = useMemo(
    () => clampLatticeColumns(latticeColumns ?? LATTICE_COLUMNS_DEFAULT),
    [latticeColumns],
  );
  const rows = chunkItems(items, columnCount);

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
            data-lattice-columns={String(columnCount)}
            data-row-end-address={rowItems.at(-1)?.address ?? ''}
            data-row-start-address={rowItems[0]?.address ?? ''}
            data-testid={`source-canvas-row-${rowIndex}`}
            style={{
              gridTemplateColumns: `repeat(${columnCount}, minmax(0, 1fr))`,
            }}
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
              const overlayText = getOverlayValue({ item, t, valueFormat, viewMode });
              const primaryText = getPrimaryLabel(item, t);
              const hoverTitle = [item.address, primaryText, overlayText]
                .filter((line) => line.length > 0)
                .join('\n');

              return (
                <button
                  aria-pressed={selectedLogicalAddresses.has(item.address)}
                  key={item.address}
                  data-merge-offset={String(item.mergeOffset)}
                  data-merge-span={String(item.mergeSpan)}
                  data-status={item.status}
                  data-testid={`address-cell-${item.address}`}
                  title={hoverTitle}
                  className={[
                    'flex min-h-28 min-w-0 flex-col items-stretch overflow-hidden rounded-xl border p-2 text-left transition',
                    'focus:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400/60',
                    statusClassName[item.status],
                    selectedLogicalAddresses.has(item.address) ? 'ring-2 ring-cyan-400/60' : '',
                  ].join(' ')}
                  onClick={() => onSelectAddress?.(item.address)}
                  style={colSpan > 1 ? { gridColumn: `span ${colSpan}` } : undefined}
                  type="button"
                >
                  <p className="break-all font-mono text-[10px] uppercase leading-snug tracking-[0.12em] text-slate-300/90 line-clamp-2">
                    {item.address}
                  </p>
                  <div className="mt-1.5 flex min-h-0 min-w-0 flex-1 flex-col gap-1 overflow-hidden">
                    <p className="min-w-0 break-words text-xs font-semibold leading-snug text-white line-clamp-2">
                      {primaryText}
                    </p>
                    <p className="min-w-0 break-all text-[11px] leading-snug text-slate-300 line-clamp-2">
                      {overlayText}
                    </p>
                    {viewMode === 'plan' && colSpan > 1 ? (
                      <p className="min-w-0 truncate text-[10px] leading-snug text-cyan-100/90">
                        {t('workbench.source.canvas.mergeState', { cells: item.mergeSpan })}
                      </p>
                    ) : null}
                    {viewMode === 'live' && item.liveTimestamp ? (
                      <p className="min-w-0 truncate text-[10px] text-slate-500">{item.liveTimestamp}</p>
                    ) : null}
                    {item.primaryRuleId ? (
                      <p className="min-w-0 truncate font-mono text-[10px] text-slate-500">{item.primaryRuleId}</p>
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
