import type { Device } from '../../../types/datalink';
import DeviceOnboardingWizard from '../../../components/datalink/wizard/DeviceOnboardingWizard';
import type { TFunction } from 'i18next';

interface SmartDashboardOverlaysProps {
  showSwitchConfirmDialog: boolean;
  setShowSwitchConfirmDialog: (open: boolean) => void;
  pendingSwitchDeviceId: string | null;
  applyDeviceSwitch: (deviceId: string) => Promise<void> | void;
  setSelectedAddresses: (addresses: string[]) => void;
  setPendingSwitchDeviceId: (deviceId: string | null) => void;
  deleteConfirmDevice: Device | null;
  setDeleteConfirmDevice: (device: Device | null) => void;
  deletingDeviceId: string | null;
  handleConfirmDeleteDevice: () => Promise<void> | void;
  isCreateDeviceModalOpen: boolean;
  closeCreateDeviceModal: () => void;
  setJustCreatedDeviceId: (deviceId: string | null) => void;
  justCreatedDeviceId: string | null;
  confirmSwitchToCreatedDevice: (switchImmediately: boolean) => void;
  t: TFunction;
}

export default function SmartDashboardOverlays({
  showSwitchConfirmDialog,
  setShowSwitchConfirmDialog,
  pendingSwitchDeviceId,
  applyDeviceSwitch,
  setSelectedAddresses,
  setPendingSwitchDeviceId,
  deleteConfirmDevice,
  setDeleteConfirmDevice,
  deletingDeviceId,
  handleConfirmDeleteDevice,
  isCreateDeviceModalOpen,
  closeCreateDeviceModal,
  setJustCreatedDeviceId,
  justCreatedDeviceId,
  confirmSwitchToCreatedDevice,
  t,
}: SmartDashboardOverlaysProps) {
  return (
    <>
      {showSwitchConfirmDialog && (
        <div className="fixed inset-0 z-[60] flex items-center justify-center bg-slate-950/70 p-4">
          <div className="w-full max-w-md rounded-2xl border border-amber-400/30 bg-slate-900 p-4 shadow-2xl">
            <h3 className="text-sm font-semibold text-amber-100">{t('smartDashboard.overlays.unsavedChangesTitle')}</h3>
            <p className="mt-2 text-xs text-slate-300">
              {t('smartDashboard.overlays.unsavedChangesDesc')}
            </p>
            <div className="mt-4 grid grid-cols-1 gap-2 sm:grid-cols-3">
              <button
                type="button"
                onClick={() => {
                  setShowSwitchConfirmDialog(false);
                  if (pendingSwitchDeviceId) void applyDeviceSwitch(pendingSwitchDeviceId);
                }}
                className="min-h-11 rounded-lg border border-blue-400/30 bg-blue-500/20 px-3 py-2 text-xs font-semibold text-blue-100 hover:bg-blue-500/30"
              >
                {t('smartDashboard.overlays.saveAndSwitch')}
              </button>
              <button
                type="button"
                onClick={() => {
                  setSelectedAddresses([]);
                  setShowSwitchConfirmDialog(false);
                  if (pendingSwitchDeviceId) void applyDeviceSwitch(pendingSwitchDeviceId);
                }}
                className="min-h-11 rounded-lg border border-amber-400/30 bg-amber-500/20 px-3 py-2 text-xs font-semibold text-amber-100 hover:bg-amber-500/30"
              >
                {t('smartDashboard.overlays.discardAndSwitch')}
              </button>
              <button
                type="button"
                autoFocus
                onClick={() => {
                  setPendingSwitchDeviceId(null);
                  setShowSwitchConfirmDialog(false);
                }}
                className="min-h-11 rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 text-xs font-semibold text-slate-200 hover:bg-slate-700"
              >
                {t('common.cancel')}
              </button>
            </div>
          </div>
        </div>
      )}

      {deleteConfirmDevice && (
        <div className="fixed inset-0 z-[70] flex items-center justify-center bg-slate-950/70 p-4">
          <div className="w-full max-w-md rounded-2xl border border-rose-400/30 bg-slate-900 p-4 shadow-2xl">
            <h3 className="text-sm font-semibold text-rose-100">{t('smartDashboard.overlays.deleteDeviceConfirmTitle')}</h3>
            <p className="mt-2 text-xs text-slate-300">
              {t('smartDashboard.overlays.deleteDeviceConfirmDesc', { name: deleteConfirmDevice.name })}
            </p>
            <p className="mt-1 text-[11px] text-slate-400">
              {t('smartDashboard.overlays.deleteDeviceConfirmHint')}
            </p>
            <div className="mt-4 grid grid-cols-2 gap-2">
              <button
                type="button"
                onClick={() => setDeleteConfirmDevice(null)}
                disabled={deletingDeviceId === deleteConfirmDevice.id}
                className="min-h-11 rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 text-xs font-semibold text-slate-200 hover:bg-slate-700 disabled:cursor-not-allowed disabled:opacity-60"
              >
                {t('common.cancel')}
              </button>
              <button
                type="button"
                onClick={() => void handleConfirmDeleteDevice()}
                disabled={deletingDeviceId === deleteConfirmDevice.id}
                className="min-h-11 rounded-lg border border-rose-400/40 bg-rose-500/20 px-3 py-2 text-xs font-semibold text-rose-100 hover:bg-rose-500/30 disabled:cursor-not-allowed disabled:opacity-60"
              >
                {deletingDeviceId === deleteConfirmDevice.id ? t('smartDashboard.overlays.deleting') : t('smartDashboard.overlays.confirmDelete')}
              </button>
            </div>
          </div>
        </div>
      )}

      {isCreateDeviceModalOpen && (
        <div className="fixed inset-0 z-[60] flex items-center justify-center overflow-y-auto bg-slate-950/70 p-2 sm:p-4">
          <div className="flex max-h-[95dvh] w-full max-w-5xl flex-col rounded-2xl border border-white/10 bg-slate-900/95 p-3 shadow-2xl sm:max-h-[92vh] sm:p-4">
            <div className="mb-2 flex shrink-0 items-center justify-between">
              <h3 className="text-sm font-semibold text-slate-100">{t('smartDashboard.overlays.createDeviceFlowTitle')}</h3>
              <button
                type="button"
                onClick={closeCreateDeviceModal}
                className="rounded-md border border-slate-700 px-2 py-1 text-xs text-slate-200 hover:bg-slate-800"
              >
                {t('common.close')}
              </button>
            </div>
            <div className="min-h-0 flex-1 overflow-y-auto overflow-x-hidden overscroll-contain">
              <DeviceOnboardingWizard
                embedded
                onClose={closeCreateDeviceModal}
                onActivated={(deviceId) => {
                  closeCreateDeviceModal();
                  setJustCreatedDeviceId(deviceId);
                }}
              />
            </div>
          </div>
        </div>
      )}

      {justCreatedDeviceId && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/70 p-4">
          <div className="w-full max-w-md rounded-2xl border border-blue-400/30 bg-slate-900 p-4 shadow-2xl">
            <h3 className="text-sm font-semibold text-blue-100">{t('smartDashboard.overlays.createDeviceDoneTitle')}</h3>
            <p className="mt-2 text-xs text-slate-300">
              {t('smartDashboard.overlays.createDeviceDoneDesc')}
            </p>
            <div className="mt-4 grid grid-cols-2 gap-2">
              <button
                type="button"
                onClick={() => confirmSwitchToCreatedDevice(false)}
                className="min-h-11 rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 text-xs font-semibold text-slate-200 hover:bg-slate-700"
              >
                {t('smartDashboard.overlays.switchLater')}
              </button>
              <button
                type="button"
                onClick={() => confirmSwitchToCreatedDevice(true)}
                className="min-h-11 rounded-lg bg-blue-600 px-3 py-2 text-xs font-semibold text-white hover:bg-blue-500"
              >
                {t('smartDashboard.overlays.switchNow')}
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
