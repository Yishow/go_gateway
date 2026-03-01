import { useCallback, useMemo } from 'react';
import type { Point } from '../../../types/datalink';

type PanelType = 'batch' | 'detail' | 'shortcuts' | null;

export interface UseSmartDashboardPanelsStateInput {
  panelType: PanelType;
  setPanelType: (panel: PanelType) => void;
  selectedPoint: Point | null;
  setSelectedAddresses: (addresses: string[]) => void;
  shortcutsTitle: string;
  pointDetailTitle: string;
  onBatchClose?: () => void;
}

export function useSmartDashboardPanelsState({
  panelType,
  setPanelType,
  selectedPoint,
  setSelectedAddresses,
  shortcutsTitle,
  pointDetailTitle,
  onBatchClose,
}: UseSmartDashboardPanelsStateInput) {
  const isBatchDialogOpen = panelType === 'batch';
  const isSlidePanelOpen = panelType === 'detail' || panelType === 'shortcuts';
  const showDetailPanel = panelType === 'detail' && Boolean(selectedPoint);
  const showShortcutsPanel = panelType === 'shortcuts';

  const slidePanelTitle = useMemo(
    () => (panelType === 'shortcuts' ? shortcutsTitle : pointDetailTitle),
    [panelType, pointDetailTitle, shortcutsTitle],
  );

  const handlePanelClose = useCallback(() => {
    setPanelType(null);
  }, [setPanelType]);

  const closeBatch = useCallback(() => {
    setPanelType(null);
    setSelectedAddresses([]);
    onBatchClose?.();
  }, [onBatchClose, setPanelType, setSelectedAddresses]);

  const handleBatchDialogOpenChange = useCallback(
    (open: boolean) => {
      if (open) return;
      setPanelType(null);
      onBatchClose?.();
    },
    [onBatchClose, setPanelType],
  );

  return {
    isBatchDialogOpen,
    isSlidePanelOpen,
    showDetailPanel,
    showShortcutsPanel,
    slidePanelTitle,
    closeBatch,
    handleBatchDialogOpenChange,
    handlePanelClose,
  };
}
