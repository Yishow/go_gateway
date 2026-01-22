
import { useMemo } from 'react';
import type { ProtocolType, Point } from '../../types/datalink';
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
}

interface GridCellStart {
  address: string;
  point?: Point;
  status: 'used' | 'available' | 'selected' | 'pending';
  isSelected: boolean;
}

export function MemoryGrid({
  // deviceId, // Unused for now in calculation
  protocol,
  centerAddress,
  range = 100,
  existingPoints,
  selectedAddresses,
  onSelect,
  onCellClick
}: MemoryGridProps) {
  
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
      
      let status: GridCellStart['status'] = 'available';
      if (point) status = 'used';
      if (isSelected) status = 'selected';
      
      return {
        address: addr,
        point,
        status,
        isSelected
      };
    });
  }, [centerAddress, range, protocol, existingPoints, selectedAddresses]);

  const handleCellClick = (cell: GridCellStart, e: React.MouseEvent) => {
    // Logical order:
    // 1. If Shift click -> Range select (TODO Phase 2 or simple implementation now)
    // 2. If Ctrl click -> Toggle select (TODO)
    // 3. Simple click -> Select single
    
    // For MVP test passing:
    onCellClick(cell.address, cell.point);
    onSelect([cell.address]);
  };

  return (
    <div 
      data-testid="memory-grid" 
      className="grid grid-cols-5 sm:grid-cols-8 md:grid-cols-10 gap-2 p-4 select-none"
    >
      {cells.map(cell => (
        <div
          key={cell.address}
          data-testid="grid-cell"
          onClick={(e) => handleCellClick(cell, e)}
          className={`
            relative aspect-[4/3] rounded-md border flex flex-col items-center justify-center cursor-pointer transition-all hover:scale-105 active:scale-95
            ${cell.status === 'used' 
              ? 'bg-green-100 border-green-500 text-green-800 dark:bg-green-900/30 dark:border-green-500/50 dark:text-green-300' 
              : cell.status === 'selected'
                ? 'bg-blue-100 border-blue-500 text-blue-800 dark:bg-blue-900/30 dark:border-blue-500 dark:text-blue-200'
                : 'bg-white border-zinc-200 text-zinc-500 hover:border-zinc-400 dark:bg-zinc-800/50 dark:border-zinc-700 dark:text-zinc-400 dark:hover:border-zinc-500'
            }
          `}
          title={cell.point ? `Point: ${cell.point.name}` : `Address: ${cell.address}`}
        >
          <span className="text-xs font-mono font-bold tracking-tight">
            {cell.address}
          </span>
          {cell.point && (
            <div className="absolute top-1 right-1 w-1.5 h-1.5 rounded-full bg-green-500 shadow-sm" />
          )}
        </div>
      ))}
      
      {cells.length === 0 && (
        <div className="col-span-full py-8 text-center text-gray-400">
          無位址資料或解析失敗
        </div>
      )}
    </div>
  );
}
