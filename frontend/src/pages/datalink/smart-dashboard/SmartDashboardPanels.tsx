import type { KeyboardShortcut } from '../../../hooks/useKeyboardShortcuts';
import type { Device, Point, PollingGroup } from '../../../types/datalink';
import { BatchPointCreator } from '../../../components/datalink/BatchPointCreator';
import { PointDetailPanel } from '../../../components/datalink/PointDetailPanel';
import { SlidePanel } from '../../../components/datalink/SlidePanel';

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
  /** 批量建立面板關閉時回呼（用於清除 initialBatchTemplate） */
  onBatchClose?: () => void;
  /** 批量建立成功後回呼（例如切換至 Tag 分頁） */
  onBatchCreated?: () => void;
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
  onBatchClose,
  onBatchCreated,
}: SmartDashboardPanelsProps) {
  return (
    <SlidePanel
      isOpen={panelType !== null}
      title={
        panelType === 'batch'
          ? batchCreateTitle
          : panelType === 'shortcuts'
            ? shortcutsTitle
            : pointDetailTitle
      }
      onClose={() => setPanelType(null)}
    >
      {panelType === 'batch' && selectedDevice && (
        <BatchPointCreator
          deviceId={selectedDevice.id}
          protocol={selectedDevice.protocol}
          preselectedAddresses={selectedAddresses}
          pollingGroups={pollingGroups}
          initialTemplate={initialBatchTemplate}
          onCreated={() => {
            setPanelType(null);
            setSelectedAddresses([]);
            onBatchClose?.();
            onBatchCreated?.();
          }}
          onCancel={() => {
            setPanelType(null);
            onBatchClose?.();
          }}
        />
      )}

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
  );
}
