import type { KeyboardShortcut } from '../../../hooks/useKeyboardShortcuts';
import type { Device, Point, PollingGroup } from '../../../types/datalink';
import type { DataType } from '../../../types/datalink';
import { BatchPointCreator } from '../../../components/datalink/BatchPointCreator';
import { PointDetailPanel } from '../../../components/datalink/PointDetailPanel';
import { SlidePanel } from '../../../components/datalink/SlidePanel';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '../../../components/ui/dialog';

type PanelType = 'batch' | 'detail' | 'shortcuts' | null;

interface SmartDashboardPanelsProps {
  panelType: PanelType;
  setPanelType: (panel: PanelType) => void;
  selectedDevice: Device | null;
  selectedPoint: Point | null;
  selectedAddresses: string[];
  setSelectedAddresses: (addresses: string[]) => void;
  pollingGroups: PollingGroup[];
  shortcuts: KeyboardShortcut[];
  batchCreateTitle: string;
  shortcutsTitle: string;
  pointDetailTitle: string;
  shortcutsHint: string;
  /** 從「套用到網格」帶入的命名模板，例如 SRC-{index03} */
  initialBatchTemplate?: string;
  /** 從「套用到網格」帶入的資料型別，與記憶體網格一致 */
  initialBatchDataType?: DataType;
  /** 批量建立關閉時回呼（用於清除 initial 狀態） */
  onBatchClose?: () => void;
  /** 批量建立成功後回呼，傳入已建立的點位（可依命名自動建立 Tag 並連結） */
  onBatchCreated?: (points: Point[]) => void;
}

export default function SmartDashboardPanels({
  panelType,
  setPanelType,
  selectedDevice,
  selectedPoint,
  selectedAddresses,
  setSelectedAddresses,
  pollingGroups,
  shortcuts,
  batchCreateTitle,
  shortcutsTitle,
  pointDetailTitle,
  shortcutsHint,
  initialBatchTemplate,
  initialBatchDataType,
  onBatchClose,
  onBatchCreated,
}: SmartDashboardPanelsProps) {
  const closeBatch = () => {
    setPanelType(null);
    setSelectedAddresses([]);
    onBatchClose?.();
  };

  return (
    <>
      {/* 批量建立：以 Modal 顯示，內容可與記憶體網格設定對應（套用到網格時帶入） */}
      <Dialog
        open={panelType === 'batch'}
        onOpenChange={(open) => {
          if (!open) {
            setPanelType(null);
            onBatchClose?.();
          }
        }}
      >
        <DialogContent
          className="max-w-lg max-h-[90vh] overflow-y-auto bg-white dark:bg-slate-900 border-slate-200 dark:border-slate-700"
          overlayClassName="bg-black/30"
        >
          <DialogHeader>
            <DialogTitle>{batchCreateTitle}</DialogTitle>
          </DialogHeader>
          {selectedDevice && (
            <BatchPointCreator
              deviceId={selectedDevice.id}
              protocol={selectedDevice.protocol}
              preselectedAddresses={selectedAddresses}
              pollingGroups={pollingGroups}
              initialTemplate={initialBatchTemplate}
              initialDataType={initialBatchDataType}
              onCreated={(points) => {
                setPanelType(null);
                onBatchClose?.();
                onBatchCreated?.(points);
                /* 不清空 selectedAddresses，讓 Tag 分頁可對應到剛建立的點位 */
              }}
              onCancel={closeBatch}
            />
          )}
        </DialogContent>
      </Dialog>

      {/* 點位詳情、快捷鍵：維持側滑面板 */}
      <SlidePanel
        isOpen={panelType === 'detail' || panelType === 'shortcuts'}
        title={panelType === 'shortcuts' ? shortcutsTitle : pointDetailTitle}
        onClose={() => setPanelType(null)}
      >
        {panelType === 'detail' && selectedPoint && (
          <PointDetailPanel
            point={selectedPoint}
            onUpdate={() => setPanelType(null)}
            onDelete={() => setPanelType(null)}
            onClose={() => setPanelType(null)}
          />
        )}

        {panelType === 'shortcuts' && (
          <div className="space-y-4 p-4">
            <p className="mb-4 text-sm text-slate-500 dark:text-slate-400">{shortcutsHint}</p>
            <div className="space-y-3">
              {shortcuts.map((shortcut, index) => (
                <div
                  key={index}
                  className="flex items-center justify-between rounded-lg bg-slate-100 px-3 py-2 dark:bg-slate-800/50"
                >
                  <span className="text-sm text-slate-700 dark:text-slate-300">{shortcut.description}</span>
                  <kbd className="rounded border border-slate-300 bg-slate-200 px-2 py-1 text-xs font-mono text-slate-700 dark:border-slate-600 dark:bg-slate-700 dark:text-slate-300">
                    {shortcut.ctrl && 'Ctrl+'}
                    {shortcut.alt && 'Alt+'}
                    {shortcut.shift && 'Shift+'}
                    {shortcut.key}
                  </kbd>
                </div>
              ))}
            </div>
          </div>
        )}
      </SlidePanel>
    </>
  );
}
