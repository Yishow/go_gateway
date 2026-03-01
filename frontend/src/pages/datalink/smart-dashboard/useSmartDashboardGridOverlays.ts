import { useCallback, useEffect, useLayoutEffect, useRef, useState, type MouseEvent as ReactMouseEvent } from 'react';
import type { Point } from '../../../types/datalink';

export interface GridContextMenuState {
  x: number;
  y: number;
  point: Point;
  shiftKey: boolean;
  pointsToDelete: Point[];
}

interface UseSmartDashboardGridOverlaysParams {
  allPoints: Point[];
  selectedAddresses: string[];
}

export function useSmartDashboardGridOverlays({
  allPoints,
  selectedAddresses,
}: UseSmartDashboardGridOverlaysParams) {
  const gridPopoverAnchorRef = useRef<HTMLElement | null>(null);
  const gridPopoverContentRef = useRef<HTMLDivElement | null>(null);
  const [gridPopoverOpen, setGridPopoverOpen] = useState(false);
  const [gridPopoverPoint, setGridPopoverPoint] = useState<Point | null>(null);
  const [gridPopoverAddress, setGridPopoverAddress] = useState('');
  const [gridPopoverPosition, setGridPopoverPosition] = useState<{ top: number; left: number } | null>(null);
  const gridContextMenuRef = useRef<HTMLDivElement | null>(null);
  const [gridContextMenu, setGridContextMenu] = useState<GridContextMenuState | null>(null);

  const handleCellClick = useCallback((address: string, point?: Point, event?: ReactMouseEvent) => {
    gridPopoverAnchorRef.current = (event?.currentTarget as HTMLElement) ?? null;
    setGridPopoverAddress(address);
    setGridPopoverPoint(point ?? null);
    setGridPopoverOpen(true);
  }, []);

  useLayoutEffect(() => {
    if (!gridPopoverOpen || !gridPopoverAnchorRef.current) {
      setGridPopoverPosition(null);
      return;
    }
    const rect = gridPopoverAnchorRef.current.getBoundingClientRect();
    setGridPopoverPosition({ top: rect.bottom + 4, left: rect.left });
  }, [gridPopoverAddress, gridPopoverOpen]);

  useEffect(() => {
    if (!gridPopoverOpen) return;
    const onMouseDown = (event: MouseEvent) => {
      const target = event.target as Node;
      if (
        gridPopoverContentRef.current?.contains(target) ||
        gridPopoverAnchorRef.current?.contains(target)
      ) {
        return;
      }
      setGridPopoverOpen(false);
    };
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') setGridPopoverOpen(false);
    };
    document.addEventListener('mousedown', onMouseDown);
    document.addEventListener('keydown', onKeyDown);
    return () => {
      document.removeEventListener('mousedown', onMouseDown);
      document.removeEventListener('keydown', onKeyDown);
    };
  }, [gridPopoverOpen]);

  const handleCellContextMenu = useCallback(
    (_address: string, point: Point, event: ReactMouseEvent) => {
      const shiftKey = event.shiftKey === true;
      const pointsToDelete =
        shiftKey && selectedAddresses.length > 0
          ? allPoints.filter((item) => selectedAddresses.includes(item.address))
          : [point];
      setGridContextMenu({
        x: event.clientX,
        y: event.clientY,
        point,
        shiftKey,
        pointsToDelete,
      });
    },
    [allPoints, selectedAddresses],
  );

  useEffect(() => {
    if (!gridContextMenu) return;
    const onMouseDown = (event: MouseEvent) => {
      const target = event.target as Node;
      if (gridContextMenuRef.current?.contains(target)) return;
      setGridContextMenu(null);
    };
    document.addEventListener('mousedown', onMouseDown);
    return () => document.removeEventListener('mousedown', onMouseDown);
  }, [gridContextMenu]);

  return {
    gridPopoverOpen,
    setGridPopoverOpen,
    gridPopoverPoint,
    gridPopoverAddress,
    gridPopoverPosition,
    gridPopoverContentRef,
    gridContextMenuRef,
    gridContextMenu,
    setGridContextMenu,
    handleCellClick,
    handleCellContextMenu,
  };
}
