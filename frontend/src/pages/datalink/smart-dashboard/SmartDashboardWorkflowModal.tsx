import type { CreateDeviceRequest, Device, UpdateDeviceRequest } from '../../../types/datalink';
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
  isSwitchingDevice: boolean;
  selectedDeviceId: string | null;
  editingDeviceInModal: Device | null;
  setEditingDeviceInModal: (device: Device | null) => void;
  handleSubmitDeviceSetup: (data: CreateDeviceRequest | UpdateDeviceRequest) => Promise<void>;
  setPanelType: (panel: PanelType) => void;
  setActiveTab: (tab: DashboardTab) => void;
  goToLocalModbusWorkbench: () => void;
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
  requestDeviceSwitch,
  isSwitchingDevice,
  selectedDeviceId,
  editingDeviceInModal,
  setEditingDeviceInModal,
  handleSubmitDeviceSetup,
  setPanelType,
  setActiveTab,
  goToLocalModbusWorkbench,
}: SmartDashboardWorkflowModalProps) {
  if (!modalIntent) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/70 p-4">
      <div className="w-full max-w-5xl rounded-2xl border border-white/10 bg-gradient-to-br from-slate-900/95 via-slate-900 to-slate-950 p-4 shadow-2xl">
        <div className="mb-3 flex items-center justify-between">
          <div>
            <p className="text-xs uppercase tracking-wider text-indigo-200/80">Dashboard Modal</p>
            <h3 className="text-base font-semibold text-slate-100">{modalIntentLabel}</h3>
          </div>
          <button
            type="button"
            onClick={closeWorkflowModal}
            className="rounded-md border border-slate-700 px-2 py-1 text-xs text-slate-200 hover:bg-slate-800"
          >
            關閉
          </button>
        </div>
        {modalIntent === 'devices' ? (
          <div className="grid grid-cols-1 gap-4 xl:grid-cols-[minmax(0,1.15fr)_minmax(0,0.85fr)]">
            <section className="rounded-2xl border border-white/10 bg-slate-900/60 p-4">
              <div className="flex flex-wrap items-center gap-2">
                <input
                  value={deviceSearchQuery}
                  onChange={(e) => setDeviceSearchQuery(e.target.value)}
                  placeholder="搜尋設備名稱 / protocol / ID"
                  className="min-h-11 flex-1 rounded-lg border border-slate-700 bg-slate-800/80 px-3 py-2 text-sm text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
                <select
                  value={deviceStatusFilter}
                  onChange={(e) => setDeviceStatusFilter(e.target.value as DeviceStatusFilter)}
                  className="min-h-11 rounded-lg border border-slate-700 bg-slate-800/80 px-3 py-2 text-sm text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="all">全部狀態</option>
                  <option value="active">active</option>
                  <option value="disabled">disabled</option>
                  <option value="draft">draft</option>
                </select>
                <button
                  type="button"
                  onClick={handleCreateDevice}
                  className="min-h-11 rounded-lg border border-blue-400/40 bg-blue-500/20 px-3 py-2 text-xs font-semibold text-blue-100 hover:bg-blue-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
                >
                  新增設備
                </button>
              </div>
              <div className="mt-3 grid grid-cols-2 gap-2 md:grid-cols-4">
                <div className="rounded-lg border border-white/10 bg-slate-800/60 p-2 text-xs text-slate-300">總數: <span className="font-semibold text-slate-100">{deviceSummary.total}</span></div>
                <div className="rounded-lg border border-emerald-500/30 bg-emerald-500/10 p-2 text-xs text-emerald-200">active: <span className="font-semibold">{deviceSummary.active}</span></div>
                <div className="rounded-lg border border-amber-500/30 bg-amber-500/10 p-2 text-xs text-amber-200">disabled: <span className="font-semibold">{deviceSummary.disabled}</span></div>
                <div className="rounded-lg border border-slate-500/40 bg-slate-800/60 p-2 text-xs text-slate-300">draft: <span className="font-semibold text-slate-100">{deviceSummary.draft}</span></div>
              </div>
              <div className="mt-3 max-h-[380px] space-y-2 overflow-auto pr-1">
                {filteredDevices.map((device) => (
                  <article key={device.id} className="rounded-xl border border-white/10 bg-slate-800/50 p-3">
                    <div className="flex items-start justify-between gap-3">
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
                        {device.status}
                      </span>
                    </div>
                    <div className="mt-2 flex flex-wrap items-center justify-between gap-2">
                      <p className="text-[11px] text-slate-400">上次測試: {device.last_test_at ? new Date(device.last_test_at).toLocaleString() : '-'}</p>
                      <div className="flex flex-wrap items-center gap-2">
                        <button
                          type="button"
                          onClick={() => openDeviceSetupModal(device.id)}
                          className="min-h-9 rounded-md border border-slate-400/30 bg-slate-700/60 px-2.5 py-1.5 text-[11px] font-semibold text-slate-100 hover:bg-slate-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-400"
                        >
                          設定
                        </button>
                        <button
                          type="button"
                          onClick={() => void handleTestDeviceConnection(device.id)}
                          disabled={testingDeviceId === device.id}
                          className="min-h-9 rounded-md border border-cyan-400/40 bg-cyan-500/20 px-2.5 py-1.5 text-[11px] font-semibold text-cyan-100 hover:bg-cyan-500/30 disabled:cursor-not-allowed disabled:opacity-60 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-500"
                        >
                          {testingDeviceId === device.id ? '測試中...' : '測試連線'}
                        </button>
                        {device.status === 'draft' ? (
                          <button
                            type="button"
                            onClick={() => void handleToggleDeviceStatusDirect(device.id, true)}
                            disabled={activatingDeviceId === device.id}
                            className="min-h-9 rounded-md border border-amber-400/40 bg-amber-500/20 px-2.5 py-1.5 text-[11px] font-semibold text-amber-100 hover:bg-amber-500/30 disabled:cursor-not-allowed disabled:opacity-60 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500"
                          >
                            {activatingDeviceId === device.id ? '啟用中...' : '啟用並切換'}
                          </button>
                        ) : (
                          <button
                            type="button"
                            onClick={() => void handleToggleDeviceStatusDirect(device.id)}
                            disabled={activatingDeviceId === device.id}
                            className="min-h-9 rounded-md border border-amber-400/40 bg-amber-500/20 px-2.5 py-1.5 text-[11px] font-semibold text-amber-100 hover:bg-amber-500/30 disabled:cursor-not-allowed disabled:opacity-60 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500"
                          >
                            {activatingDeviceId === device.id ? '處理中...' : device.status === 'active' ? '停用' : '啟用'}
                          </button>
                        )}
                        <button
                          type="button"
                          onClick={() => requestDeleteDevice(device.id)}
                          disabled={deletingDeviceId === device.id}
                          className="min-h-9 rounded-md border border-rose-400/40 bg-rose-500/20 px-2.5 py-1.5 text-[11px] font-semibold text-rose-100 hover:bg-rose-500/30 disabled:cursor-not-allowed disabled:opacity-60 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-rose-500"
                        >
                          {deletingDeviceId === device.id ? '刪除中...' : '刪除'}
                        </button>
                        <button
                          type="button"
                          onClick={() => requestDeviceSwitch(device.id)}
                          disabled={isSwitchingDevice || selectedDeviceId === device.id}
                          className="min-h-9 rounded-md border border-blue-400/40 bg-blue-500/20 px-2.5 py-1.5 text-[11px] font-semibold text-blue-100 hover:bg-blue-500/30 disabled:cursor-not-allowed disabled:opacity-60 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
                        >
                          {isSwitchingDevice
                            ? '切換中...'
                            : selectedDeviceId === device.id
                              ? '目前設備'
                              : '切換'}
                        </button>
                      </div>
                    </div>
                  </article>
                ))}
                {filteredDevices.length === 0 && (
                  <div className="rounded-xl border border-dashed border-slate-600 p-4 text-center text-xs text-slate-400">
                    查無符合條件的設備
                  </div>
                )}
              </div>
            </section>
            <section className="rounded-2xl border border-white/10 bg-slate-900/60 p-4">
              <div className="flex items-center justify-between gap-3">
                <div>
                  <p className="text-xs uppercase tracking-wider text-slate-400">Device Setup</p>
                  <h4 className="text-sm font-semibold text-slate-100">
                    {editingDeviceInModal ? `設定：${editingDeviceInModal.name}` : '請先選擇設備'}
                  </h4>
                </div>
                {editingDeviceInModal && (
                  <button
                    type="button"
                    onClick={() => setEditingDeviceInModal(null)}
                    className="min-h-9 rounded-md border border-slate-700 px-2.5 py-1.5 text-[11px] text-slate-200 hover:bg-slate-800"
                  >
                    關閉設定
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
                      {testingDeviceId === editingDeviceInModal.id ? '測試中...' : '測試連線'}
                    </button>
                    <button
                      type="button"
                      onClick={() => void handleToggleDeviceStatusDirect(editingDeviceInModal.id)}
                      disabled={activatingDeviceId === editingDeviceInModal.id}
                      className="min-h-9 rounded-md border border-amber-400/40 bg-amber-500/20 px-3 py-1.5 text-[11px] font-semibold text-amber-100 hover:bg-amber-500/30 disabled:cursor-not-allowed disabled:opacity-60"
                    >
                      {activatingDeviceId === editingDeviceInModal.id
                        ? '處理中...'
                        : editingDeviceInModal.status === 'active'
                          ? '停用設備'
                          : '啟用設備'}
                    </button>
                    <button
                      type="button"
                      onClick={() => requestDeleteDevice(editingDeviceInModal.id)}
                      disabled={deletingDeviceId === editingDeviceInModal.id}
                      className="min-h-9 rounded-md border border-rose-400/40 bg-rose-500/20 px-3 py-1.5 text-[11px] font-semibold text-rose-100 hover:bg-rose-500/30 disabled:cursor-not-allowed disabled:opacity-60"
                    >
                      {deletingDeviceId === editingDeviceInModal.id
                        ? '刪除中...'
                        : '刪除設備'}
                    </button>
                  </div>
                  <div className="max-h-[52vh] overflow-auto rounded-xl border border-white/10 bg-slate-900/70 p-3">
                    <DeviceForm
                      device={editingDeviceInModal}
                      onSubmit={handleSubmitDeviceSetup}
                      onCancel={() => setEditingDeviceInModal(null)}
                    />
                  </div>
                </div>
              ) : (
                <div className="mt-3 rounded-xl border border-dashed border-slate-600 bg-slate-900/60 p-4 text-xs text-slate-400">
                  從左側設備卡片點擊「設定」，即可在此直接編輯來源協議、連線參數、重試策略並測試連線，不再跳轉到獨立設定頁。
                </div>
              )}
            </section>
          </div>
        ) : (
          <>
            <div className="rounded-xl border border-white/10 bg-slate-800/50 p-3 text-sm text-slate-200">
              <p>
                {modalIntent === 'settings' && '系統設定以 Dashboard 內嵌設定模式開啟。'}
                {modalIntent === 'points' && '點位流程已整合到 Source Planner + Batch 建立。'}
                {modalIntent === 'mappings' && '映射流程已整合到 Flow + Tag Linkage 區。'}
                {modalIntent === 'wizard' && '精靈流程以新增設備 modal 承載。'}
                {modalIntent === 'polling-groups' && '輪詢群組管理透過設定與點位流程整合。'}
                {modalIntent === 'tags' && 'Tag 管理與全域編輯整合在右側面板。'}
              </p>
            </div>
            <div className="mt-3 grid grid-cols-1 gap-2 sm:grid-cols-2">
              <button
                type="button"
                onClick={() => {
                  if (modalIntent === 'wizard') handleCreateDevice();
                  if (modalIntent === 'points') setPanelType('batch');
                  if (modalIntent === 'settings') setActiveTab('settings');
                  if (modalIntent === 'mappings' || modalIntent === 'tags') setActiveTab('overview');
                  closeWorkflowModal();
                }}
                className="min-h-11 rounded-lg border border-blue-400/40 bg-blue-500/20 px-3 py-2 text-xs font-semibold text-blue-100 hover:bg-blue-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
              >
                開啟對應流程
              </button>
              <button
                type="button"
                onClick={() => {
                  goToLocalModbusWorkbench();
                  closeWorkflowModal();
                }}
                className="min-h-11 rounded-lg border border-cyan-400/40 bg-cyan-500/20 px-3 py-2 text-xs font-semibold text-cyan-100 hover:bg-cyan-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-500"
              >
                前往 Server Memory Grid
              </button>
            </div>
          </>
        )}
      </div>
    </div>
  );
}
