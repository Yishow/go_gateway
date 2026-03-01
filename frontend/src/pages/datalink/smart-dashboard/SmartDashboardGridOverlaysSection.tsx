import { createPortal } from 'react-dom';
import type { RefObject } from 'react';
import type { QueryClient } from '@tanstack/react-query';
import type { TFunction } from 'i18next';
import { PointDetailPanel } from '../../../components/datalink/PointDetailPanel';
import { pointKeys } from '../../../hooks/datalink/keys';
import type { Point } from '../../../types/datalink';
import SmartDashboardTagPanel, { type SmartDashboardTagPanelProps } from './SmartDashboardTagPanel';
import type { GridContextMenuState } from './useSmartDashboardGridOverlays';

interface SmartDashboardGridOverlaysSectionProps {
  gridPopoverOpen: boolean;
  setGridPopoverOpen: (open: boolean) => void;
  gridPopoverPosition: { top: number; left: number } | null;
  gridPopoverAddress: string;
  gridPopoverPoint: Point | null;
  gridPopoverContentRef: RefObject<HTMLDivElement | null>;
  tagPanelProps: SmartDashboardTagPanelProps;
  queryClient: QueryClient;
  deletePoint: (id: string) => Promise<unknown>;
  selectedPoint: Point | null;
  setSelectedPoint: (point: Point | null) => void;
  showSuccess: (message: string) => void;
  showError: (message: string) => void;
  gridContextMenu: GridContextMenuState | null;
  gridContextMenuRef: RefObject<HTMLDivElement | null>;
  setGridContextMenu: (menu: GridContextMenuState | null) => void;
  t: TFunction;
}

export default function SmartDashboardGridOverlaysSection({
  gridPopoverOpen,
  setGridPopoverOpen,
  gridPopoverPosition,
  gridPopoverAddress,
  gridPopoverPoint,
  gridPopoverContentRef,
  tagPanelProps,
  queryClient,
  deletePoint,
  selectedPoint,
  setSelectedPoint,
  showSuccess,
  showError,
  gridContextMenu,
  gridContextMenuRef,
  setGridContextMenu,
  t,
}: SmartDashboardGridOverlaysSectionProps) {
  return (
    <>
      {gridPopoverOpen &&
        gridPopoverPosition &&
        createPortal(
          <div
            ref={gridPopoverContentRef}
            className="fixed z-[100] flex max-h-[85vh] w-[min(24rem,calc(100vw-1rem))] flex-col rounded-xl border border-slate-200 bg-white shadow-2xl dark:border-slate-700 dark:bg-slate-900"
            style={{
              top: Math.min(gridPopoverPosition.top, window.innerHeight - 400),
              left: Math.min(gridPopoverPosition.left, Math.max(0, window.innerWidth - 384)),
            }}
            role="dialog"
            aria-label={t('smartDashboard.gridCellPopover.ariaLabel')}
          >
            <div className="flex shrink-0 items-center justify-between border-b border-slate-200 px-3 py-2 dark:border-slate-700">
              <div>
                <span className="font-mono text-xs text-slate-500 dark:text-slate-400">
                  {gridPopoverAddress}
                </span>
                {gridPopoverPoint && (
                  <span className="ml-2 text-sm font-medium text-slate-800 dark:text-slate-200">
                    {gridPopoverPoint.name}
                  </span>
                )}
              </div>
              <button
                type="button"
                onClick={() => setGridPopoverOpen(false)}
                className="rounded p-1.5 text-slate-400 hover:bg-slate-100 hover:text-slate-600 dark:hover:bg-slate-700 dark:hover:text-slate-200"
                aria-label={t('common.close', '關閉')}
              >
                ×
              </button>
            </div>
            <div className="min-h-0 flex-1 overflow-y-auto p-3">
              {gridPopoverPoint ? (
                <PointDetailPanel
                  point={gridPopoverPoint}
                  onUpdate={(updatedPoint) => {
                    queryClient.invalidateQueries({ queryKey: pointKeys.lists() });
                    setSelectedPoint(updatedPoint);
                    setGridPopoverOpen(false);
                  }}
                  onDelete={async (id) => {
                    try {
                      await deletePoint(id);
                      showSuccess(t('smartDashboard.pointDeleted'));
                      setSelectedPoint(null);
                      setGridPopoverOpen(false);
                    } catch (error) {
                      showError(error instanceof Error ? error.message : t('smartDashboard.deletePointFailed'));
                    }
                  }}
                  onClose={() => setGridPopoverOpen(false)}
                />
              ) : (
                <div className="rounded-lg border border-white/10 bg-slate-800/40 p-3">
                  <SmartDashboardTagPanel {...tagPanelProps} />
                </div>
              )}
            </div>
          </div>,
          document.body,
        )}

      {gridContextMenu &&
        createPortal(
          <div
            ref={gridContextMenuRef}
            className="fixed z-[110] min-w-[160px] rounded-lg border border-slate-200 bg-white py-1 shadow-xl dark:border-slate-700 dark:bg-slate-800"
            style={{ left: gridContextMenu.x, top: gridContextMenu.y }}
            role="menu"
            aria-label={t('smartDashboard.gridContextMenu.ariaLabel')}
          >
            {gridContextMenu.pointsToDelete.length > 1 ? (
              <button
                type="button"
                role="menuitem"
                className="flex w-full items-center gap-2 px-3 py-2 text-left text-sm text-red-600 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/30"
                onClick={async () => {
                  const { pointsToDelete } = gridContextMenu;
                  if (
                    !window.confirm(
                      t('smartDashboard.gridContextMenu.deleteSelectedConfirm', {
                        count: pointsToDelete.length,
                      }),
                    )
                  ) {
                    return;
                  }
                  const deletedIds = new Set<string>();
                  for (const point of pointsToDelete) {
                    try {
                      await deletePoint(point.id);
                      deletedIds.add(point.id);
                    } catch (error) {
                      showError(error instanceof Error ? error.message : t('smartDashboard.deletePointFailed'));
                    }
                  }
                  if (deletedIds.size > 0) {
                    showSuccess(
                      t('smartDashboard.gridContextMenu.deleteSelectedDone', {
                        count: deletedIds.size,
                      }),
                    );
                    if (selectedPoint && deletedIds.has(selectedPoint.id)) setSelectedPoint(null);
                  }
                  setGridContextMenu(null);
                }}
              >
                {t('smartDashboard.gridContextMenu.deleteSelectedPoints', {
                  count: gridContextMenu.pointsToDelete.length,
                })}
              </button>
            ) : (
              <button
                type="button"
                role="menuitem"
                className="flex w-full items-center gap-2 px-3 py-2 text-left text-sm text-red-600 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/30"
                onClick={async () => {
                  const { point } = gridContextMenu;
                  if (!window.confirm(t('smartDashboard.gridContextMenu.deleteConfirm', { name: point.name }))) {
                    return;
                  }
                  try {
                    await deletePoint(point.id);
                    showSuccess(t('smartDashboard.pointDeleted'));
                    setGridContextMenu(null);
                    if (selectedPoint?.id === point.id) setSelectedPoint(null);
                  } catch (error) {
                    showError(error instanceof Error ? error.message : t('smartDashboard.deletePointFailed'));
                  }
                }}
              >
                {t('smartDashboard.gridContextMenu.deletePoint')}
              </button>
            )}
          </div>,
          document.body,
        )}
    </>
  );
}
