import type { CreateDeviceRequest, Device, UpdateDeviceRequest } from '../../../types/datalink';
import type { TFunction } from 'i18next';
import DeviceForm from '../../../components/datalink/DeviceForm';

type DeviceStatusFilter = 'all' | 'active' | 'disabled' | 'draft';
type PanelType = 'batch' | 'detail' | 'shortcuts' | null;
type DashboardTab = 'overview' | 'devices' | 'settings';

interface SmartDashboardWorkflowModalProps {
  modalIntent: string | null;
  modalIntentLabel: string;
  closeWorkflowModal: () => void;
  deviceSearchQuery: string;
  setDeviceSearchQuery: (value: string) => void;
  deviceStatusFilter: DeviceStatusFilter;
  setDeviceStatusFilter: (value: DeviceStatusFilter) => void;
  handleCreateDevice: () => void;
  deviceSummary: {
    total: number;
    active: number;
    disabled: number;
    draft: number;
  };
  filteredDevices: Device[];
  openDeviceSetupModal: (deviceId?: string | null) => void;
  handleTestDeviceConnection: (deviceId: string) => Promise<void> | void;
  testingDeviceId: string | null;
  handleToggleDeviceStatusDirect: (deviceId: string, closeModalOnSuccess?: boolean) => Promise<void> | void;
  activatingDeviceId: string | null;
  requestDeleteDevice: (deviceId: string) => void;
  deletingDeviceId: string | null;
  requestDeviceSwitch: (deviceId: string) => void;
  /** 點選設備卡片時：切換到該設備（draft 會先啟用再切換），保持連接並關閉 modal */
  handleSelectDevice: (device: Device) => void;
  isSwitchingDevice: boolean;
  selectedDeviceId: string | null;
  editingDeviceInModal: Device | null;
  setEditingDeviceInModal: (device: Device | null) => void;
  handleSubmitDeviceSetup: (data: CreateDeviceRequest | UpdateDeviceRequest) => Promise<void>;
  setPanelType: (panel: PanelType) => void;
  setActiveTab: (tab: DashboardTab) => void;
  goToLocalModbusWorkbench: () => void;
  t: TFunction;
}

export default function SmartDashboardWorkflowModal({
  modalIntent,
  modalIntentLabel,
  closeWorkflowModal,
  deviceSearchQuery,
  setDeviceSearchQuery,
  deviceStatusFilter,
  setDeviceStatusFilter,
  handleCreateDevice,
  deviceSummary,
  filteredDevices,
  openDeviceSetupModal,
  handleTestDeviceConnection,
  testingDeviceId,
  handleToggleDeviceStatusDirect,
  activatingDeviceId,
  requestDeleteDevice,
  deletingDeviceId,
  requestDeviceSwitch: _requestDeviceSwitch,
  handleSelectDevice,
  isSwitchingDevice: _isSwitchingDevice,
  selectedDeviceId,
  editingDeviceInModal,
  setEditingDeviceInModal,
  handleSubmitDeviceSetup,
  setPanelType: _setPanelType,
  setActiveTab: _setActiveTab,
  goToLocalModbusWorkbench: _goToLocalModbusWorkbench,
  t,
}: SmartDashboardWorkflowModalProps) {
  if (!modalIntent) return null;

  const deviceStatusLabel = (status: Device['status']) => t(`device.status.${status}`);

  // 收斂為設備入口：非設備 intent 顯示簡化訊息
  if (modalIntent !== 'devices') {
    return (
      <div className="fixed inset-0 z-50 flex items-center justify-center overflow-y-auto bg-slate-950/70 p-2 sm:p-4">
        <div className="my-auto flex w-full max-w-lg flex-col rounded-2xl border border-white/10 bg-gradient-to-br from-slate-900/95 via-slate-900 to-slate-950 p-4 shadow-2xl">
          <div className="mb-4 flex shrink-0 items-center justify-between gap-2">
            <div className="min-w-0">
              <p className="text-[10px] uppercase tracking-wider text-indigo-200/80 sm:text-xs">{t('smartDashboard.workflowModal.title')}</p>
              <h3 className="truncate text-base font-semibold text-slate-100">{modalIntentLabel}</h3>
            </div>
            <button
              type="button"
              onClick={closeWorkflowModal}
              className="rounded-md border border-slate-700 px-2 py-1 text-xs text-slate-200 hover:bg-slate-800"
            >
              {t('common.close')}
            </button>
          </div>
          <div className="space-y-3">
            <div className="rounded-xl border border-amber-500/20 bg-amber-500/5 p-3 text-sm text-amber-100">
              <p className="font-semibold mb-1">{t('smartDashboard.workflowModal.legacyIntentNotice', { defaultValue: '此功能已整合至主工作流程' })}</p>
              <p className="text-xs text-amber-200/80">
                {modalIntent === 'points' && t('smartDashboard.workflowModal.legacyIntentHint.points', { defaultValue: '點位管理請使用右側「規劃」分頁進行設定。' })}
                {modalIntent === 'mappings' && t('smartDashboard.workflowModal.legacyIntentHint.mappings', { defaultValue: '映射設定請使用右側「Tag」與「Modbus」分頁。' })}
                {modalIntent === 'tags' && t('smartDashboard.workflowModal.legacyIntentHint.tags', { defaultValue: 'Tag 設定請使用右側「Tag」分頁。' })}
                {modalIntent === 'polling-groups' && t('smartDashboard.workflowModal.legacyIntentHint.pollingGroups', { defaultValue: '輪詢群組設定請使用設備設定頁面。' })}
                {modalIntent === 'settings' && t('smartDashboard.workflowModal.legacyIntentHint.settings', { defaultValue: '系統設定請使用頂部選單。' })}
                {modalIntent === 'wizard' && t('smartDashboard.workflowModal.legacyIntentHint.wizard', { defaultValue: '請先選擇或建立設備，然後使用主工作流程進行設定。' })}
              </p>
            </div>
            <button
              type="button"
              onClick={() => {
                closeWorkflowModal();
                // 重定向到設備選擇
                const newParams = new URLSearchParams(window.location.search);
                newParams.set('modal', 'devices');
                window.history.replaceState({}, '', `${window.location.pathname}?${newParams.toString()}`);
              }}
              className="w-full rounded-lg border border-blue-500/40 bg-blue-500/20 px-3 py-2 text-xs font-semibold text-blue-100 hover:bg-blue-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
            >
              {t('smartDashboard.workflowModal.goToDevices', { defaultValue: '前往設備管理中心' })}
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center overflow-y-auto bg-slate-950/70 p-2 sm:p-4">
      <div className="my-auto flex w-full max-w-5xl flex-col rounded-2xl border border-white/10 bg-gradient-to-br from-slate-900/95 via-slate-900 to-slate-950 p-3 shadow-2xl sm:my-0 sm:max-h-[95dvh] sm:p-4">
        <div className="mb-3 flex shrink-0 items-center justify-between gap-2">
          <div className="min-w-0">
            <p className="text-[10px] uppercase tracking-wider text-indigo-200/80 sm:text-xs">{t('smartDashboard.workflowModal.deviceCenterTitle', { defaultValue: '設備管理中心' })}</p>
            <h3 className="truncate text-base font-semibold text-slate-100">{t('nav.devices')}</h3>
          </div>
          <button
            type="button"
            onClick={closeWorkflowModal}
            className="rounded-md border border-slate-700 px-2 py-1 text-xs text-slate-200 hover:bg-slate-800"
          >
            {t('common.close')}
          </button>
        </div>
        {modalIntent === 'devices' && (
          <div className="min-h-0 flex-1 overflow-y-auto">
            <div className="grid min-h-0 grid-cols-1 gap-4 xl:grid-cols-[minmax(0,1.15fr)_minmax(0,0.85fr)]">
            <section className="flex min-h-0 flex-col rounded-2xl border border-white/10 bg-slate-900/60 p-3 sm:p-4">
              <div className="flex flex-col gap-2 sm:flex-row sm:flex-wrap sm:items-center">
                <input
                  value={deviceSearchQuery}
                  onChange={(e) => setDeviceSearchQuery(e.target.value)}
                  placeholder={t('smartDashboard.workflowModal.searchPlaceholder')}
                  className="min-h-11 min-w-0 flex-1 rounded-lg border border-slate-700 bg-slate-800/80 px-3 py-2 text-sm text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500 sm:min-w-[180px]"
                />
                <select
                  value={deviceStatusFilter}
                  onChange={(e) => setDeviceStatusFilter(e.target.value as DeviceStatusFilter)}
                  className="min-h-11 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-3 py-2 text-sm text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500 sm:w-auto sm:min-w-[120px]"
                >
                  <option value="all">{t('smartDashboard.workflowModal.statusAll')}</option>
                  <option value="active">{deviceStatusLabel('active')}</option>
                  <option value="disabled">{deviceStatusLabel('disabled')}</option>
                  <option value="draft">{deviceStatusLabel('draft')}</option>
                </select>
                <button
                  type="button"
                  onClick={handleCreateDevice}
                  className="min-h-11 w-full shrink-0 rounded-lg border border-blue-400/40 bg-blue-500/20 px-3 py-2 text-xs font-semibold text-blue-100 hover:bg-blue-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 sm:w-auto"
                >
                  {t('smartDashboard.actions.createDevice')}
                </button>
              </div>
              <div className="mt-3 grid grid-cols-2 gap-2 sm:grid-cols-4">
                <div className="rounded-lg border border-white/10 bg-slate-800/60 p-2 text-xs text-slate-300">{t('smartDashboard.workflowModal.summary.total', { count: deviceSummary.total })}</div>
                <div className="rounded-lg border border-emerald-500/30 bg-emerald-500/10 p-2 text-xs text-emerald-200">{t('smartDashboard.workflowModal.summary.active', { count: deviceSummary.active })}</div>
                <div className="rounded-lg border border-amber-500/30 bg-amber-500/10 p-2 text-xs text-amber-200">{t('smartDashboard.workflowModal.summary.disabled', { count: deviceSummary.disabled })}</div>
                <div className="rounded-lg border border-slate-500/40 bg-slate-800/60 p-2 text-xs text-slate-300">{t('smartDashboard.workflowModal.summary.draft', { count: deviceSummary.draft })}</div>
              </div>
              <div className="mt-3 max-h-[45vh] space-y-2 overflow-auto pr-1 sm:max-h-[380px]">
                {filteredDevices.map((device) => (
                  <article
                    key={device.id}
                    role="button"
                    tabIndex={0}
                    className="cursor-pointer rounded-xl border border-white/10 bg-slate-800/50 p-3 transition-colors hover:bg-slate-800/70 focus:outline-none focus:ring-2 focus:ring-blue-500/50"
                    onClick={(e) => {
                      if ((e.target as HTMLElement).closest('button')) return;
                      handleSelectDevice(device);
                    }}
                    onKeyDown={(e) => {
                      if (e.key !== 'Enter' && e.key !== ' ') return;
                      if ((e.target as HTMLElement).closest('button')) return;
                      e.preventDefault();
                      handleSelectDevice(device);
                    }}
                  >
                    <div className="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between sm:gap-3">
                      <div className="min-w-0">
                        <p className="truncate text-sm font-semibold text-slate-100">{device.name}</p>
                        <p className="mt-1 text-[11px] font-mono text-slate-400">{device.id.slice(0, 8)} · {device.protocol}</p>
                      </div>
                      <span className={`rounded-md px-2 py-1 text-[10px] font-bold uppercase ${
                        device.status === 'active'
                          ? 'bg-emerald-500/20 text-emerald-200'
                          : device.status === 'disabled'
                            ? 'bg-amber-500/20 text-amber-200'
                            : 'bg-slate-700/70 text-slate-300'
                      }`}>
                        {deviceStatusLabel(device.status)}
                      </span>
                    </div>
                    <div className="mt-2 flex flex-col gap-2 sm:flex-row sm:flex-wrap sm:items-center sm:justify-between">
                      <p className="min-w-0 truncate text-[11px] text-slate-400">
                        {t('smartDashboard.workflowModal.lastTestAt', { value: device.last_test_at ? new Date(device.last_test_at).toLocaleString() : '-' })}
                      </p>
                      <div className="flex flex-wrap items-center gap-2">
                        <button
                          type="button"
                          onClick={(e) => {
                            e.stopPropagation();
                            openDeviceSetupModal(device.id);
                          }}
                          className="min-h-9 rounded-md border border-slate-400/30 bg-slate-700/60 px-2.5 py-1.5 text-[11px] font-semibold text-slate-100 hover:bg-slate-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-400"
                        >
                          {t('smartDashboard.workflowModal.setup')}
                        </button>
                        <button
                          type="button"
                          onClick={(e) => {
                            e.stopPropagation();
                            void handleTestDeviceConnection(device.id);
                          }}
                          disabled={testingDeviceId === device.id}
                          className="min-h-9 rounded-md border border-cyan-400/40 bg-cyan-500/20 px-2.5 py-1.5 text-[11px] font-semibold text-cyan-100 hover:bg-cyan-500/30 disabled:cursor-not-allowed disabled:opacity-60 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-500"
                        >
                          {testingDeviceId === device.id ? t('smartDashboard.workflowModal.testing') : t('smartDashboard.workflowModal.testConnection')}
                        </button>
                        {device.status !== 'draft' && (
                          <button
                            type="button"
                            onClick={(e) => {
                              e.stopPropagation();
                              void handleToggleDeviceStatusDirect(device.id);
                            }}
                            disabled={activatingDeviceId === device.id}
                            className="min-h-9 rounded-md border border-amber-400/40 bg-amber-500/20 px-2.5 py-1.5 text-[11px] font-semibold text-amber-100 hover:bg-amber-500/30 disabled:cursor-not-allowed disabled:opacity-60 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500"
                          >
                            {activatingDeviceId === device.id ? t('smartDashboard.workflowModal.processing') : device.status === 'active' ? t('smartDashboard.workflowModal.disable') : t('smartDashboard.workflowModal.enable')}
                          </button>
                        )}
                        <button
                          type="button"
                          onClick={(e) => {
                            e.stopPropagation();
                            requestDeleteDevice(device.id);
                          }}
                          disabled={deletingDeviceId === device.id}
                          className="min-h-9 rounded-md border border-rose-400/40 bg-rose-500/20 px-2.5 py-1.5 text-[11px] font-semibold text-rose-100 hover:bg-rose-500/30 disabled:cursor-not-allowed disabled:opacity-60 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-rose-500"
                        >
                          {deletingDeviceId === device.id ? t('smartDashboard.workflowModal.deleting') : t('common.delete')}
                        </button>
                        {selectedDeviceId === device.id && (
                          <span className="min-h-9 inline-flex items-center rounded-md border border-blue-400/40 bg-blue-500/20 px-2.5 py-1.5 text-[11px] font-semibold text-blue-100">
                            {t('smartDashboard.workflowModal.currentDevice')}
                          </span>
                        )}
                      </div>
                    </div>
                  </article>
                ))}
                {filteredDevices.length === 0 && (
                  <div className="rounded-xl border border-dashed border-slate-600 p-4 text-center text-xs text-slate-400">
                    {t('smartDashboard.workflowModal.emptyDevices')}
                  </div>
                )}
              </div>
            </section>
            <section className="flex min-h-0 flex-col rounded-2xl border border-white/10 bg-slate-900/60 p-3 sm:p-4">
              <div className="flex shrink-0 items-center justify-between gap-3">
                <div className="min-w-0">
                  <p className="text-[10px] uppercase tracking-wider text-slate-400 sm:text-xs">{t('smartDashboard.workflowModal.deviceSetupTitle')}</p>
                  <h4 className="truncate text-sm font-semibold text-slate-100">
                    {editingDeviceInModal
                      ? t('smartDashboard.workflowModal.editingDeviceTitle', { name: editingDeviceInModal.name })
                      : t('smartDashboard.workflowModal.selectDeviceFirst')}
                  </h4>
                </div>
                {editingDeviceInModal && (
                  <button
                    type="button"
                    onClick={() => setEditingDeviceInModal(null)}
                    className="min-h-9 rounded-md border border-slate-700 px-2.5 py-1.5 text-[11px] text-slate-200 hover:bg-slate-800"
                  >
                    {t('smartDashboard.workflowModal.closeSetup')}
                  </button>
                )}
              </div>
              {editingDeviceInModal ? (
                <div className="mt-3 space-y-3">
                  <div className="flex flex-wrap gap-2">
                    <button
                      type="button"
                      onClick={() => void handleTestDeviceConnection(editingDeviceInModal.id)}
                      disabled={testingDeviceId === editingDeviceInModal.id}
                      className="min-h-9 rounded-md border border-cyan-400/40 bg-cyan-500/20 px-3 py-1.5 text-[11px] font-semibold text-cyan-100 hover:bg-cyan-500/30 disabled:cursor-not-allowed disabled:opacity-60"
                    >
                      {testingDeviceId === editingDeviceInModal.id ? t('smartDashboard.workflowModal.testing') : t('smartDashboard.workflowModal.testConnection')}
                    </button>
                    <button
                      type="button"
                      onClick={() => void handleToggleDeviceStatusDirect(editingDeviceInModal.id)}
                      disabled={activatingDeviceId === editingDeviceInModal.id}
                      className="min-h-9 rounded-md border border-amber-400/40 bg-amber-500/20 px-3 py-1.5 text-[11px] font-semibold text-amber-100 hover:bg-amber-500/30 disabled:cursor-not-allowed disabled:opacity-60"
                    >
                      {activatingDeviceId === editingDeviceInModal.id
                        ? t('smartDashboard.workflowModal.processing')
                        : editingDeviceInModal.status === 'active'
                          ? t('smartDashboard.workflowModal.disableDevice')
                          : t('smartDashboard.workflowModal.enableDevice')}
                    </button>
                    <button
                      type="button"
                      onClick={() => requestDeleteDevice(editingDeviceInModal.id)}
                      disabled={deletingDeviceId === editingDeviceInModal.id}
                      className="min-h-9 rounded-md border border-rose-400/40 bg-rose-500/20 px-3 py-1.5 text-[11px] font-semibold text-rose-100 hover:bg-rose-500/30 disabled:cursor-not-allowed disabled:opacity-60"
                    >
                      {deletingDeviceId === editingDeviceInModal.id
                        ? t('smartDashboard.workflowModal.deleting')
                        : t('smartDashboard.workflowModal.deleteDevice')}
                    </button>
                  </div>
                  <div className="max-h-[45vh] overflow-auto rounded-xl border border-white/10 bg-slate-900/70 p-3 sm:max-h-[52vh]">
                    <DeviceForm
                      device={editingDeviceInModal}
                      onSubmit={handleSubmitDeviceSetup}
                      onCancel={() => setEditingDeviceInModal(null)}
                    />
                  </div>
                </div>
              ) : (
                <div className="mt-3 rounded-xl border border-dashed border-slate-600 bg-slate-900/60 p-4 text-xs text-slate-400">
                  {t('smartDashboard.workflowModal.setupHint')}
                </div>
              )}
            </section>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
