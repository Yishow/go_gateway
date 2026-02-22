import type { ModbusShareStatus } from '../../../types/datalink';

export interface SmartDashboardModbusPanelProps {
  goToLocalModbusWorkbench: () => void;
  loadModbusStatus: () => void;
  modbusStatus: ModbusShareStatus | null;
  modbusRegister: string;
  setModbusRegister: (value: string) => void;
  handleBindTagToModbus: () => void;
  handlePushCurrentValueToModbus: () => void;
  handleSyncModbusFromMappings: () => void;
}

/** Local Modbus Share 操作面板 */
export default function SmartDashboardModbusPanel({
  goToLocalModbusWorkbench,
  loadModbusStatus,
  modbusStatus,
  modbusRegister,
  setModbusRegister,
  handleBindTagToModbus,
  handlePushCurrentValueToModbus,
  handleSyncModbusFromMappings,
}: SmartDashboardModbusPanelProps) {
  return (
    <div className="space-y-3 border-t border-white/5 p-4">
      {/* 標題列 */}
      <div className="flex items-center justify-between">
        <p className="text-xs font-semibold tracking-wide text-slate-200">Local Modbus Share</p>
        <div className="flex items-center gap-2">
          <button
            type="button"
            onClick={goToLocalModbusWorkbench}
            className="rounded border border-blue-500/30 bg-blue-500/10 px-2 py-1 text-[10px] text-blue-200 hover:bg-blue-500/20 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
          >
            完整工作台
          </button>
          <button
            type="button"
            onClick={loadModbusStatus}
            className="rounded border border-slate-700 px-2 py-1 text-[10px] text-slate-300 hover:bg-slate-800/70 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
          >
            Refresh
          </button>
        </div>
      </div>

      {/* 狀態摘要 */}
      <div className="grid grid-cols-3 gap-1.5 rounded-lg border border-white/10 bg-slate-900/60 p-2.5 text-[11px]">
        <div>
          <p className="text-slate-500">狀態</p>
          <p className={`mt-0.5 font-semibold ${modbusStatus?.enabled ? 'text-emerald-300' : 'text-slate-400'}`}>
            {modbusStatus?.enabled ? 'Running' : 'Stopped'}
          </p>
        </div>
        <div>
          <p className="text-slate-500">Address</p>
          <p className="mt-0.5 font-mono text-slate-200">{modbusStatus?.address || '—'}</p>
        </div>
        <div>
          <p className="text-slate-500">Mappings</p>
          <p className="mt-0.5 font-semibold text-slate-100">{modbusStatus?.mapping_count ?? 0}</p>
        </div>
      </div>

      {/* Register 輸入 */}
      <label className="block text-[11px] text-slate-300">
        Register (Holding)
        <input
          value={modbusRegister}
          onChange={(e) => setModbusRegister(e.target.value)}
          className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs font-mono focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </label>

      {/* 動作按鈕 */}
      <div className="space-y-2">
        <div className="grid grid-cols-2 gap-2">
          <button
            type="button"
            onClick={handleBindTagToModbus}
            className="rounded-lg border border-indigo-500/40 bg-indigo-500/20 px-3 py-2 text-xs font-medium text-indigo-100 hover:bg-indigo-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500"
          >
            綁定 Tag 到 Register
          </button>
          <button
            type="button"
            onClick={handlePushCurrentValueToModbus}
            className="rounded-lg border border-emerald-500/40 bg-emerald-500/20 px-3 py-2 text-xs font-medium text-emerald-100 hover:bg-emerald-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
          >
            推送目前值
          </button>
        </div>
        <button
          type="button"
          onClick={handleSyncModbusFromMappings}
          className="w-full rounded-lg border border-blue-500/40 bg-blue-500/20 px-3 py-2 text-xs font-medium text-blue-100 hover:bg-blue-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
        >
          同步全部啟用映射
        </button>
      </div>
    </div>
  );
}
