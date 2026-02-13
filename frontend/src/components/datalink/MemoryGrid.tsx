
import { useMemo } from 'react';
import type { ProtocolType, Point, DataType } from '../../types/datalink';
import { addressParser } from '../../utils/addressParser';
import type { ConflictSeverity, OccupancyStatus } from '../../features/datalink/typedOccupancy';
import { resolveConflictSeverity } from '../../features/datalink/typedOccupancy';

export interface MemoryGridProps {
  deviceId: string;
  protocol: ProtocolType;
  /** Center address for the view (e.g. "D100") */
  centerAddress: string;
  /** Number of cells to show (default: 100) */
  range?: number;
  existingPoints: Point[];
  selectedAddresses: string[];
  onSelect: (addresses: string[]) => void;
  onCellClick: (address: string, point?: Point) => void;
  plannedAllocations?: PlannedAllocation[];
  linkedAddresses?: string[];
  showConflictsOnly?: boolean;
}

export interface PlannedAllocation {
  id: string;
  dataType: DataType;
  addresses: string[];
  label: string;
}

interface GridCellStart {
  address: string;
  point?: Point;
  status: OccupancyStatus;
  conflictSeverity: ConflictSeverity;
  isSelected: boolean;
  plan?: PlannedAllocation;
  planCellIndex?: number;
  isPlanGroupStart?: boolean;
  isPlanGroupEnd?: boolean;
}

export function MemoryGrid({
  // deviceId, // Unused for now in calculation
  protocol,
  centerAddress,
  range = 100,
  existingPoints,
  selectedAddresses,
  onSelect,
  onCellClick,
  plannedAllocations = [],
  linkedAddresses = [],
  showConflictsOnly = false,
}: MemoryGridProps) {
  const plannedAddressMap = useMemo(() => {
    const map = new Map<
      string,
      {
        plan: PlannedAllocation;
        index: number;
        isGroupStart: boolean;
        isGroupEnd: boolean;
      }
    >();
    plannedAllocations.forEach((plan) => {
      plan.addresses.forEach((address, index) => {
        map.set(address, {
          plan,
          index,
          isGroupStart: index === 0,
          isGroupEnd: index === plan.addresses.length - 1,
        });
      });
    });
    return map;
  }, [plannedAllocations]);
  const linkedAddressSet = useMemo(() => new Set(linkedAddresses), [linkedAddresses]);

  const cells = useMemo(() => {
    if (!centerAddress) return [];

    const addresses = addressParser.expand(centerAddress, range, protocol);

    return addresses.map(addr => {
      const point = existingPoints.find(p => p.address === addr);
      const isSelected = selectedAddresses.includes(addr);
      const planned = plannedAddressMap.get(addr);
      const isLinked = linkedAddressSet.has(addr);

      const conflictSeverity = planned
        ? resolveConflictSeverity({
            hasUsedPoint: Boolean(point),
            hasLinkedAddress: isLinked,
          })
        : 'none';

      let status: GridCellStart['status'] = 'available';
      if (point) status = 'used';
      if (isLinked) status = 'linked';
      if (planned) status = conflictSeverity === 'none' ? 'planned' : 'conflict';
      if (isSelected && status !== 'conflict') status = 'selected';

      return {
        address: addr,
        point,
        status,
        conflictSeverity,
        isSelected,
        plan: planned?.plan,
        planCellIndex: planned?.index,
        isPlanGroupStart: planned?.isGroupStart,
        isPlanGroupEnd: planned?.isGroupEnd,
      };
    });
  }, [centerAddress, range, protocol, existingPoints, selectedAddresses, plannedAddressMap, linkedAddressSet]);

  const visibleCells = useMemo(() => {
    if (!showConflictsOnly) return cells;
    const visibleIndexSet = new Set<number>();
    cells.forEach((cell, index) => {
      if (cell.status !== 'conflict') return;
      visibleIndexSet.add(index);
      if (index > 0) visibleIndexSet.add(index - 1);
      if (index < cells.length - 1) visibleIndexSet.add(index + 1);
    });
    return cells.filter((_, index) => visibleIndexSet.has(index));
  }, [cells, showConflictsOnly]);

  const handleCellClick = (cell: GridCellStart, e: React.MouseEvent) => {
    let newSelection: string[] = [];

    if (e.ctrlKey || e.metaKey) {
      if (selectedAddresses.includes(cell.address)) {
        newSelection = selectedAddresses.filter(a => a !== cell.address);
      } else {
        newSelection = [...selectedAddresses, cell.address];
      }
    } else if (e.shiftKey && selectedAddresses.length > 0) {
      const lastSelected = selectedAddresses[selectedAddresses.length - 1];
      const allAddresses = cells.map(c => c.address);
      const currentIndex = allAddresses.indexOf(cell.address);
      const lastIndex = allAddresses.indexOf(lastSelected);

      if (currentIndex !== -1 && lastIndex !== -1) {
        const start = Math.min(currentIndex, lastIndex);
        const end = Math.max(currentIndex, lastIndex);
        const range = allAddresses.slice(start, end + 1);
        newSelection = Array.from(new Set([...selectedAddresses, ...range]));
      } else {
        newSelection = [cell.address];
      }
    } else {
      newSelection = [cell.address];
    }

    onSelect(newSelection);

    onCellClick(cell.address, cell.point);
  };

  return (
    <div 
      data-testid="memory-grid" 
      className="grid grid-cols-5 sm:grid-cols-8 md:grid-cols-10 gap-2 p-4 select-none"
    >
      {visibleCells.map(cell => (
        <button
          type="button"
          key={cell.address}
          data-testid="grid-cell"
          data-status={cell.status}
          data-conflict-severity={cell.conflictSeverity}
          onClick={(e) => handleCellClick(cell, e)}
          aria-label={`Address ${cell.address}, status ${cell.status}`}
          className={`
            relative aspect-[4/3] border flex flex-col items-center justify-center cursor-pointer transition-colors duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500
            ${cell.status === 'used' 
              ? 'bg-green-100 border-green-500 text-green-800 dark:bg-green-900/30 dark:border-green-500/50 dark:text-green-300' 
              : cell.status === 'linked'
                ? 'bg-violet-100 border-violet-500 text-violet-800 dark:bg-violet-900/25 dark:border-violet-500/60 dark:text-violet-200'
              : cell.status === 'planned'
                ? 'bg-sky-100 border-sky-500 text-sky-800 dark:bg-sky-900/30 dark:border-sky-500 dark:text-sky-200'
                : cell.status === 'conflict'
                  ? 'bg-rose-100 border-rose-500 text-rose-800 dark:bg-rose-900/30 dark:border-rose-500 dark:text-rose-200'
              : cell.status === 'selected'
                ? 'bg-blue-100 border-blue-500 text-blue-800 dark:bg-blue-900/30 dark:border-blue-500 dark:text-blue-200'
                : 'bg-white border-zinc-200 text-zinc-500 hover:border-zinc-400 dark:bg-zinc-800/50 dark:border-zinc-700 dark:text-zinc-400 dark:hover:border-zinc-500'
            }
            ${!cell.plan || cell.plan.addresses.length <= 1 ? 'rounded-md' : ''}
            ${cell.plan && cell.plan.addresses.length > 1 && cell.isPlanGroupStart ? 'rounded-l-md rounded-r-none border-l-2' : ''}
            ${cell.plan && cell.plan.addresses.length > 1 && cell.isPlanGroupEnd ? 'rounded-r-md rounded-l-none border-r-2' : ''}
            ${cell.plan && cell.plan.addresses.length > 1 && !cell.isPlanGroupStart && !cell.isPlanGroupEnd ? 'rounded-none border-y-2' : ''}
            ${cell.isSelected ? 'ring-2 ring-blue-400/70' : ''}
          `}
          title={
            cell.status === 'conflict'
              ? `Conflict (${cell.conflictSeverity}): ${cell.address}`
              : cell.point
                ? `Point: ${cell.point.name}`
                : `Address: ${cell.address}`
          }
        >
          <span className="text-xs font-mono font-bold tracking-tight">
            {cell.address}
          </span>
          {cell.plan && (
            <span className="mt-0.5 text-[10px] font-semibold tracking-wide opacity-80">
              {cell.plan.label}
              {cell.plan.addresses.length > 1 ? ` ${cell.planCellIndex! + 1}/${cell.plan.addresses.length}` : ''}
            </span>
          )}
          {cell.plan && cell.plan.addresses.length > 1 && (
            <span className="absolute top-1 left-1 rounded bg-sky-500/25 px-1 text-[9px] font-bold text-sky-100">
              #{(cell.planCellIndex ?? 0) + 1}
            </span>
          )}
          {cell.point && (
            <div className="absolute top-1 right-1 w-1.5 h-1.5 rounded-full bg-green-500 shadow-sm" />
          )}
          {cell.status === 'conflict' && (
            <div className="absolute inset-x-1 bottom-1 rounded bg-rose-500/20 text-[9px] text-rose-700 dark:text-rose-200 text-center">
              {cell.conflictSeverity === 'hard' ? '硬衝突' : '軟衝突'}
            </div>
          )}
        </button>
      ))}
      
      {visibleCells.length === 0 && (
        <div className="col-span-full py-8 text-center text-gray-400">
          {showConflictsOnly ? '目前沒有衝突格位' : '無位址資料或解析失敗'}
        </div>
      )}
    </div>
  );
}
