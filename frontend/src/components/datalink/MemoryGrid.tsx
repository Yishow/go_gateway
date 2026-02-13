
import { useMemo } from 'react';
import type { ProtocolType, Point, DataType } from '../../types/datalink';
import { addressParser } from '../../utils/addressParser';

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
  status: 'used' | 'available' | 'selected' | 'planned' | 'conflict';
  isSelected: boolean;
  plan?: PlannedAllocation;
  planCellIndex?: number;
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
  showConflictsOnly = false,
}: MemoryGridProps) {
  const plannedAddressMap = useMemo(() => {
    const map = new Map<string, { plan: PlannedAllocation; index: number }>();
    plannedAllocations.forEach((plan) => {
      plan.addresses.forEach((address, index) => {
        map.set(address, { plan, index });
      });
    });
    return map;
  }, [plannedAllocations]);

  // Calculate cells
  const cells = useMemo(() => {
    // Determine start address. 
    // Ideally centerAddress should be in the middle, but for MVP let's start FROM centerAddress 
    // or slightly before if possible. 
    // addressParser.expand starts from given address.
    // For now, let's treat "centerAddress" as "startAddress" for simplicity matching the spec "Start D100, show 100"
    
    // Safety check
    if (!centerAddress) return [];

    const addresses = addressParser.expand(centerAddress, range, protocol);
    
    return addresses.map(addr => {
      const point = existingPoints.find(p => p.address === addr);
      const isSelected = selectedAddresses.includes(addr);
      const planned = plannedAddressMap.get(addr);
      
      let status: GridCellStart['status'] = 'available';
      if (point) status = 'used';
      if (planned) status = point ? 'conflict' : 'planned';
      if (isSelected && status !== 'conflict') status = 'selected';
      
      return {
        address: addr,
        point,
        status,
        isSelected,
        plan: planned?.plan,
        planCellIndex: planned?.index,
      };
    });
  }, [centerAddress, range, protocol, existingPoints, selectedAddresses, plannedAddressMap]);

  const visibleCells = useMemo(() => {
    if (!showConflictsOnly) return cells;
    return cells.filter((cell) => cell.status === 'conflict');
  }, [cells, showConflictsOnly]);

  const handleCellClick = (cell: GridCellStart, e: React.MouseEvent) => {
    // Determine new selection
    let newSelection: string[] = [];

    if (e.ctrlKey || e.metaKey) {
      // Toggle logic
      if (selectedAddresses.includes(cell.address)) {
        newSelection = selectedAddresses.filter(a => a !== cell.address);
      } else {
        newSelection = [...selectedAddresses, cell.address];
      }
    } else if (e.shiftKey && selectedAddresses.length > 0) {
      // Range logic
      // Ideally we need to find the index of the start and end address in the current grid
      const lastSelected = selectedAddresses[selectedAddresses.length - 1]; // Naive 'last' based on array
      // A better way is to find the index of the clicked cell and the index of the LAST clicked cell.
      // But we typically track 'anchor' in state. For this stateless component, let's use the last item in selectedAddresses as anchor if possible,
      // OR better, since we don't store order in parent likely, we find the range between MIN and MAX indices of "anchor" and "current".
      // Let's assume the user wants range from explicit Last Selected -> Current.
      
      const allAddresses = cells.map(c => c.address);
      const currentIndex = allAddresses.indexOf(cell.address);
      const lastIndex = allAddresses.indexOf(lastSelected);
      
      if (currentIndex !== -1 && lastIndex !== -1) {
        const start = Math.min(currentIndex, lastIndex);
        const end = Math.max(currentIndex, lastIndex);
        const range = allAddresses.slice(start, end + 1);
        // Union with existing selection to be safe, or just replace? Usual PC behavior is replace selection with range relative to anchor.
        // But for simplicity let's just ADD the range to existing (or Union). 
        // Actually standard is: Shift-click extends selection from anchor.
        // Let's go with: Union of current selection + new range.
        newSelection = Array.from(new Set([...selectedAddresses, ...range]));
      } else {
        newSelection = [cell.address];
      }
    } else {
      // Single select
      newSelection = [cell.address];
    }
    
    // Call onSelect first
    onSelect(newSelection);
    
    // Then notify click.
    // Important: SmartDashboard currently uses onCellClick to open panels.
    // We should pass the modifier keys or let SmartDashboard handle it.
    // But better: The parent can check if it wants to open panel.
    // We will just pass the event as is or let the parent deal with it?
    // The interface is `onCellClick: (address: string, point?: Point) => void;`
    // We can't change signature easily without breaking other things? 
    // Actually we can just call it for now, and fix SmartDashboard to ignore if not a point.
    onCellClick(cell.address, cell.point);
  };

  return (
    <div 
      data-testid="memory-grid" 
      className="grid grid-cols-5 sm:grid-cols-8 md:grid-cols-10 gap-2 p-4 select-none"
    >
      {visibleCells.map(cell => (
        <div
          key={cell.address}
          data-testid="grid-cell"
          onClick={(e) => handleCellClick(cell, e)}
          className={`
            relative aspect-[4/3] rounded-md border flex flex-col items-center justify-center cursor-pointer transition-all hover:scale-105 active:scale-95
            ${cell.status === 'used' 
              ? 'bg-green-100 border-green-500 text-green-800 dark:bg-green-900/30 dark:border-green-500/50 dark:text-green-300' 
              : cell.status === 'planned'
                ? 'bg-sky-100 border-sky-500 text-sky-800 dark:bg-sky-900/30 dark:border-sky-500 dark:text-sky-200'
                : cell.status === 'conflict'
                  ? 'bg-rose-100 border-rose-500 text-rose-800 dark:bg-rose-900/30 dark:border-rose-500 dark:text-rose-200'
              : cell.status === 'selected'
                ? 'bg-blue-100 border-blue-500 text-blue-800 dark:bg-blue-900/30 dark:border-blue-500 dark:text-blue-200'
                : 'bg-white border-zinc-200 text-zinc-500 hover:border-zinc-400 dark:bg-zinc-800/50 dark:border-zinc-700 dark:text-zinc-400 dark:hover:border-zinc-500'
            }
            ${cell.plan && cell.plan.addresses.length > 1 ? 'ring-1 ring-offset-1 ring-offset-transparent ring-sky-400/60' : ''}
            ${cell.isSelected ? 'ring-2 ring-blue-400/70' : ''}
          `}
          title={
            cell.status === 'conflict'
              ? `Conflict: ${cell.address} already used`
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
          {cell.point && (
            <div className="absolute top-1 right-1 w-1.5 h-1.5 rounded-full bg-green-500 shadow-sm" />
          )}
          {cell.status === 'conflict' && (
            <div className="absolute inset-x-1 bottom-1 rounded bg-rose-500/20 text-[9px] text-rose-700 dark:text-rose-200 text-center">
              衝突
            </div>
          )}
        </div>
      ))}
      
      {visibleCells.length === 0 && (
        <div className="col-span-full py-8 text-center text-gray-400">
          {showConflictsOnly ? '目前沒有衝突格位' : '無位址資料或解析失敗'}
        </div>
      )}
    </div>
  );
}
