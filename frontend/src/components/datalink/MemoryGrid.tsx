import { useMemo } from "react";
import type { ProtocolType, Point, DataType } from "../../types/datalink";
import { addressParser } from "../../utils/addressParser";
import type { ConflictSeverity, OccupancyStatus } from "../../features/datalink/typedOccupancy";
import { getSpanByDataType, resolveConflictSeverity } from "../../features/datalink/typedOccupancy";

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
  /** 點擊格位時回呼；可傳入 event 供上層用於 Popover 定位 */
  onCellClick: (address: string, point?: Point, e?: React.MouseEvent) => void;
  /** 右鍵格位時回呼（僅有點位時觸發，供上層顯示刪除等選單） */
  onCellContextMenu?: (address: string, point: Point, e: React.MouseEvent) => void;
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

/** 顯示單元：單一格或合併多格（如 float32 兩格為一邏輯格），一個設定／一個 tag */
interface DisplayUnit {
  addresses: string[];
  span: 1 | 2 | 4;
  status: OccupancyStatus;
  conflictSeverity: ConflictSeverity;
  isSelected: boolean;
  plan?: PlannedAllocation;
  point?: Point;
  startCellIndex: number;
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
  onCellContextMenu,
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

    return addresses.map((addr) => {
      const point = existingPoints.find((p) => p.address === addr);
      const isSelected = selectedAddresses.includes(addr);
      const planned = plannedAddressMap.get(addr);
      const isLinked = linkedAddressSet.has(addr);

      const conflictSeverity = planned
        ? resolveConflictSeverity({
            hasUsedPoint: Boolean(point),
            hasLinkedAddress: isLinked,
          })
        : "none";

      let status: GridCellStart["status"] = "available";
      if (point) status = "used";
      if (isLinked) status = "linked";
      if (planned) status = conflictSeverity === "none" ? "planned" : "conflict";
      if (isSelected && status !== "conflict") status = "selected";

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

  /** 多格點位：僅「非起始」的位址集合，用於略過已合併的延續格 */
  const pointCoveredAddresses = useMemo(() => {
    const set = new Set<string>();
    existingPoints.forEach((point) => {
      const span = getSpanByDataType(point.data_type);
      if (span <= 1) return;
      const addrs = addressParser.expand(point.address, span, protocol);
      addrs.slice(1).forEach((addr) => set.add(addr));
    });
    return set;
  }, [existingPoints, protocol]);

  /** 將連續位址依「既有點位資料型別」與「規劃」合併為顯示單元：int16 一格、float32 兩格合一、int64/float64 四格合一；重載後仍依點位 data_type 正確顯示 */
  const displayUnits = useMemo(() => {
    const units: DisplayUnit[] = [];
    let i = 0;
    while (i < cells.length) {
      const cell = cells[i];
      if (pointCoveredAddresses.has(cell.address)) {
        i += 1;
        continue;
      }
      const point = cell.point;
      const pointSpan = point ? getSpanByDataType(point.data_type) : 1;
      const pointAddrs =
        point && pointSpan > 1 ? addressParser.expand(point.address, pointSpan, protocol) : [];
      const canMergeByPoint =
        point &&
        pointSpan > 1 &&
        point.address === cell.address &&
        pointAddrs.length === pointSpan &&
        pointAddrs.every((addr, j) => cells[i + j]?.address === addr);

      if (canMergeByPoint) {
        const s = pointSpan as 2 | 4;
        let status = cell.status;
        let conflictSeverity = cell.conflictSeverity;
        for (let j = 1; j < s; j++) {
          const c = cells[i + j];
          if (c?.status === "conflict") {
            status = "conflict";
            conflictSeverity = c.conflictSeverity;
          }
        }
        const isSelected = pointAddrs.some((a) => selectedAddresses.includes(a));
        units.push({
          addresses: pointAddrs,
          span: s,
          status,
          conflictSeverity,
          isSelected,
          plan: cell.plan,
          point,
          startCellIndex: i,
        });
        i += s;
        continue;
      }

      const planned = plannedAddressMap.get(cell.address);
      const plan = planned?.plan;
      const addrs = plan?.addresses ?? [];
      const span = addrs.length;
      const canMergeByPlan =
        planned?.isGroupStart &&
        span > 1 &&
        Array.from({ length: span }, (_, j) => cells[i + j]?.address === addrs[j]).every(Boolean);

      if (canMergeByPlan) {
        const s = span as 2 | 4;
        let status = cell.status;
        let conflictSeverity = cell.conflictSeverity;
        for (let j = 1; j < s; j++) {
          const c = cells[i + j];
          if (c?.status === "conflict") {
            status = "conflict";
            conflictSeverity = c.conflictSeverity;
          }
        }
        const isSelected = addrs.some((a) => selectedAddresses.includes(a));
        const pointAtStart = existingPoints.find((p) => p.address === addrs[0]);
        units.push({
          addresses: addrs,
          span: s,
          status,
          conflictSeverity,
          isSelected,
          plan: plan!,
          point: pointAtStart,
          startCellIndex: i,
        });
        i += s;
      } else {
        units.push({
          addresses: [cell.address],
          span: 1,
          status: cell.status,
          conflictSeverity: cell.conflictSeverity,
          isSelected: cell.isSelected,
          plan: cell.plan,
          point: cell.point,
          startCellIndex: i,
        });
        i += 1;
      }
    }
    return units;
  }, [cells, plannedAddressMap, pointCoveredAddresses, selectedAddresses, existingPoints]);

  const conflictVisibleIndexSet = useMemo(() => {
    const set = new Set<number>();
    cells.forEach((cell, index) => {
      if (cell.status !== "conflict") return;
      set.add(index);
      if (index > 0) set.add(index - 1);
      if (index < cells.length - 1) set.add(index + 1);
    });
    return set;
  }, [cells]);

  const visibleDisplayUnits = useMemo(() => {
    if (!showConflictsOnly) return displayUnits;
    return displayUnits.filter((unit) => {
      for (let k = 0; k < unit.span; k++) {
        if (conflictVisibleIndexSet.has(unit.startCellIndex + k)) return true;
      }
      return false;
    });
  }, [displayUnits, showConflictsOnly, conflictVisibleIndexSet]);

  const allAddresses = useMemo(() => cells.map((c) => c.address), [cells]);

  const handleUnitClick = (unit: DisplayUnit, e: React.MouseEvent) => {
    let newSelection: string[] = [];

    if (e.ctrlKey || e.metaKey) {
      const allInSelection = unit.addresses.every((a) => selectedAddresses.includes(a));
      if (allInSelection) {
        newSelection = selectedAddresses.filter((a) => !unit.addresses.includes(a));
      } else {
        newSelection = Array.from(new Set([...selectedAddresses, ...unit.addresses]));
      }
    } else if (e.shiftKey && selectedAddresses.length > 0) {
      const lastSelected = selectedAddresses[selectedAddresses.length - 1];
      const currentStart = allAddresses.indexOf(unit.addresses[0]);
      const lastIndex = allAddresses.indexOf(lastSelected);
      if (currentStart !== -1 && lastIndex !== -1) {
        const start = Math.min(currentStart, lastIndex);
        const end = Math.max(currentStart + unit.span - 1, lastIndex);
        const range = allAddresses.slice(start, end + 1);
        newSelection = Array.from(new Set([...selectedAddresses, ...range]));
      } else {
        newSelection = unit.addresses;
      }
    } else {
      newSelection = unit.addresses;
    }

    onSelect(newSelection);
    onCellClick(unit.addresses[0], unit.point, e);
  };

  return (
    <div
      data-testid="memory-grid"
      className="grid h-full min-h-0 flex-1 grid-cols-5 sm:grid-cols-8 md:grid-cols-10 gap-1.5 p-3 select-none"
      style={{ gridAutoRows: "minmax(2.5rem, 1fr)" }}
    >
      {visibleDisplayUnits.map((unit) => (
        <button
          type="button"
          key={unit.addresses[0]}
          data-testid="grid-cell"
          data-status={unit.status}
          data-conflict-severity={unit.conflictSeverity}
          data-span={unit.span}
          onClick={(e) => handleUnitClick(unit, e)}
          onContextMenu={(e) => {
            if (unit.point) {
              e.preventDefault();
              onCellContextMenu?.(unit.addresses[0], unit.point, e);
            }
          }}
          aria-label={
            unit.addresses.length > 1
              ? `${unit.plan?.label ?? unit.point?.name ?? "Block"} ${unit.addresses[0]}–${unit.addresses[unit.addresses.length - 1]}, ${unit.status}`
              : `Address ${unit.addresses[0]}, status ${unit.status}`
          }
          className={`
            relative min-w-0 min-h-[2.5rem] border rounded-md flex flex-col items-center justify-center cursor-pointer transition-colors duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 overflow-hidden
            ${
              unit.status === "used"
                ? "bg-green-100 border-green-500 text-green-800 dark:bg-green-900/30 dark:border-green-500/50 dark:text-green-300"
                : unit.status === "linked"
                  ? "bg-violet-100 border-violet-500 text-violet-800 dark:bg-violet-900/25 dark:border-violet-500/60 dark:text-violet-200"
                  : unit.status === "planned"
                    ? "bg-sky-100 border-sky-500 text-sky-800 dark:bg-sky-900/30 dark:border-sky-500 dark:text-sky-200"
                    : unit.status === "conflict"
                      ? "bg-rose-100 border-rose-500 text-rose-800 dark:bg-rose-900/30 dark:border-rose-500 dark:text-rose-200"
                      : unit.status === "selected"
                        ? "bg-blue-100 border-blue-500 text-blue-800 dark:bg-blue-900/30 dark:border-blue-500 dark:text-blue-200"
                        : "bg-white border-zinc-200 text-zinc-500 hover:border-zinc-400 dark:bg-zinc-800/50 dark:border-zinc-700 dark:text-zinc-400 dark:hover:border-zinc-500"
            }
            ${unit.isSelected ? "ring-2 ring-blue-400/70" : ""}
          `}
          style={{ gridColumn: `span ${unit.span}` }}
          title={
            unit.status === "conflict"
              ? `Conflict (${unit.conflictSeverity}): ${unit.addresses.join("–")}`
              : unit.point
                ? `Point: ${unit.point.name}`
                : unit.addresses.length > 1
                  ? `${unit.plan?.label ?? "Block"} ${unit.addresses[0]}–${unit.addresses[unit.addresses.length - 1]}`
                  : `Address ${unit.addresses[0]}`
          }
        >
          {unit.span === 1 ? (
            <>
              <span className="text-[10px] font-mono font-bold tracking-tight">{unit.addresses[0]}</span>
              {unit.plan && (
                <span className="mt-0.5 text-[9px] font-semibold tracking-wide opacity-80">{unit.plan.label}</span>
              )}
            </>
          ) : (
            <>
              <span className="text-[9px] font-mono font-bold tracking-tight leading-tight truncate max-w-full">
                {unit.addresses[0]}–{unit.addresses[unit.addresses.length - 1]}
              </span>
              <span className="mt-0.5 text-[8px] font-semibold tracking-wide opacity-90 leading-tight truncate max-w-full">
                {unit.plan?.label ?? unit.point?.name ?? "Block"} · {unit.plan?.dataType ?? unit.point?.data_type ?? ""}
              </span>
            </>
          )}
          {unit.point && <div className="absolute top-1 right-1 w-1.5 h-1.5 rounded-full bg-green-500 shadow-sm" />}
          {unit.status === "conflict" && (
            <div className="absolute inset-x-1 bottom-1 rounded bg-rose-500/20 text-[9px] text-rose-700 dark:text-rose-200 text-center">
              {unit.conflictSeverity === "hard" ? "硬衝突" : "軟衝突"}
            </div>
          )}
        </button>
      ))}

      {visibleDisplayUnits.length === 0 && (
        <div className="col-span-full py-8 text-center text-gray-400">
          {showConflictsOnly ? "目前沒有衝突格位" : "無位址資料或解析失敗"}
        </div>
      )}
    </div>
  );
}
