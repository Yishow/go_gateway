import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { ModbusShareSettings } from '../state/types';
import type { ModbusShareStatus } from '../../../../types/modbusShare';
import { modbusShareAPI } from '../../../../services/modbusShare';
import { getSafeErrorMessage, type SafeErrorMessage } from '../../../../utils/typedErrors';

/**
 * ModbusShareSection 元件屬性
 */
interface ModbusShareSectionProps {
  settings: ModbusShareSettings;
  onChange: (patch: Partial<ModbusShareSettings>) => void;
  pending?: boolean;
  persistedBaseRegister?: number | null;
  shareStatus?: ModbusShareStatus | null;
}

/**
 * Local Modbus Share 設定區塊元件
 * 落地需求：「Local Modbus Share configuration」
 */
export function ModbusShareSection({
  settings,
  onChange,
  pending = false,
  persistedBaseRegister = null,
  shareStatus = null,
}: ModbusShareSectionProps) {
  const { t } = useTranslation('workbench-v2');
  const interactionLockRef = React.useRef(false);
  const [runtimeStatus, setRuntimeStatus] = React.useState(shareStatus);
  const [lifecyclePending, setLifecyclePending] = React.useState(false);
  const [lifecycleError, setLifecycleError] = React.useState<SafeErrorMessage | null>(null);

  const { enabled, bind_address, port, slave_id, capacity_registers } = settings;
  const expectedSettingsRevision = settings.expected_settings_revision
    || settings.settings_revision
    || runtimeStatus?.settings_revision;

  React.useEffect(() => {
    setRuntimeStatus(shareStatus);
  }, [shareStatus]);

  const handleLifecycle = async (action: 'start' | 'stop') => {
    if (lifecyclePending || pending) return;
    setLifecyclePending(true);
    setLifecycleError(null);
    try {
      const nextStatus = action === 'start'
        ? await modbusShareAPI.start({ port, expected_settings_revision: expectedSettingsRevision })
        : await modbusShareAPI.stop(expectedSettingsRevision);
      setRuntimeStatus(nextStatus);
    } catch (error) {
      setLifecycleError(getSafeErrorMessage(error, (key, options) => t(key, options)));
    } finally {
      setLifecyclePending(false);
    }
  };

  const handleToggle = () => {
    if (pending || interactionLockRef.current) {
      return;
    }
    interactionLockRef.current = true;
    onChange({ enabled: !enabled });
    queueMicrotask(() => {
      interactionLockRef.current = false;
    });
  };

  return (
    <div className="space-y-6 bg-gray-900/10 border border-gray-800 rounded-2xl p-6 backdrop-blur-sm">
      {/* 標題與側欄開關 */}
      <div className="flex justify-between items-start">
        <div>
          <h3 className="text-lg font-semibold text-white">
            {t('settings.modbus_share_title', 'Local Modbus Share')}
          </h3>
          <p className="text-xs text-gray-500 mt-1">
            {t('settings.modbus_share_subtitle', '將採集資料映射並轉發至本地虛擬 Modbus Server 暫存器，供上層 SCADA/HMI 直接抓取。')}
          </p>
        </div>
        
        {/* 總開關 */}
        <div className="flex items-center gap-2">
          <span className="text-xs text-gray-400 select-none">
            {enabled ? t('status.active', '啟用') : t('status.disabled', '未啟用')}
          </span>
          <button
            type="button"
            aria-label={enabled ? t('status.active', '啟用') : t('status.disabled', '未啟用')}
            aria-pressed={enabled}
            aria-busy={pending}
            disabled={pending}
            onClick={handleToggle}
            className={`
              relative inline-flex h-6 w-11 items-center rounded-full transition-colors duration-200 focus:outline-none cursor-pointer disabled:cursor-not-allowed disabled:opacity-60
              ${enabled ? 'bg-blue-600' : 'bg-gray-800'}
            `}
          >
            <span
              className={`
                inline-block h-4 w-4 transform rounded-full bg-white transition-transform duration-200
                ${enabled ? 'translate-x-6' : 'translate-x-1'}
              `}
            />
          </button>
        </div>
      </div>
      <ShareRuntimeStatus
        status={runtimeStatus}
        lifecyclePending={lifecyclePending}
        lifecycleError={lifecycleError}
        onStart={() => void handleLifecycle('start')}
        onStop={() => void handleLifecycle('stop')}
      />

      {/* 依據開關渲染內容 */}
      {enabled ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4 p-4 rounded-xl border border-gray-800/60 bg-gray-950/20 sweep-in">
          {/* 綁定位址 */}
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-gray-400">
              {t('settings.share_bind_address', '監聽本機位址')}
            </label>
            <input
              type="text"
              value={bind_address}
              onChange={(e) => onChange({ bind_address: e.target.value })}
              className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all font-mono"
              placeholder="0.0.0.0"
            />
          </div>

          {/* 通訊埠 */}
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-gray-400">
              {t('settings.share_port', '監聽通訊埠 (Port)')}
            </label>
            <input
              type="number"
              min={1}
              max={65535}
              value={port}
              onChange={(e) => onChange({ port: Math.max(1, parseInt(e.target.value) || 1) })}
              className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all font-mono"
              placeholder="5020"
            />
          </div>

          {/* Slave ID */}
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-gray-400">
              {t('settings.share_slave_id', '本機站號 (Slave ID)')}
            </label>
            <input
              type="number"
              min={1}
              max={247}
              value={slave_id}
              onChange={(e) => onChange({ slave_id: Math.max(1, parseInt(e.target.value) || 1) })}
              className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all font-mono"
              placeholder="1"
            />
          </div>

          {/* 起始暫存器位址由已保存 source-rule candidate 決定，不在 Settings 形成 browser-only truth。 */}
          <div className="space-y-1.5" data-testid="share-base-register-value">
            <label className="text-xs font-medium text-gray-400">
              {t('settings.share_base_register', '起始暫存器位址')}
            </label>
            <div className="w-full bg-gray-950/60 border border-gray-800 rounded-lg px-3 py-2 text-gray-300 text-sm font-mono">
              {persistedBaseRegister ?? t('settings.share_base_register_auto', '自動由已保存候選配置')}
            </div>
          </div>

          {/* 暫存器容量 */}
          <div className="space-y-1.5">
            <label className="text-xs font-medium text-gray-400">
              {t('settings.share_capacity_registers', '暫存器容量')}
            </label>
            <input
              type="number"
              min={1}
              max={32768}
              value={capacity_registers}
              onChange={(e) => onChange({ capacity_registers: Math.max(1, parseInt(e.target.value) || 1) })}
              className="w-full bg-gray-950 border border-gray-800 rounded-lg px-3 py-2 text-white text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-all font-mono"
              placeholder="32768"
            />
          </div>
        </div>
      ) : (
        <div className="flex items-center justify-center p-6 border border-dashed border-gray-800 rounded-xl bg-gray-950/10 text-gray-500 text-xs">
          <span>{t('settings.modbus_share_disabled_msg', 'Modbus 轉發服務已停用。開啟後點位分配將會與 Step 2 的 Modbus Share 暫存器配置連動。')}</span>
        </div>
      )}
    </div>
  );
}

function ShareRuntimeStatus({
  status,
  lifecyclePending,
  lifecycleError,
  onStart,
  onStop,
}: {
  status: ModbusShareStatus | null | undefined;
  lifecyclePending: boolean;
  lifecycleError: SafeErrorMessage | null;
  onStart: () => void;
  onStop: () => void;
}) {
  const { t } = useTranslation('workbench-v2');
  if (!status) return null;
  const hydrationState = status.hydration_state;
  if (hydrationState && hydrationState !== 'ready') {
    const key = hydrationState === 'failed' ? 'failed' : 'pending';
    return (
      <p data-testid="share-runtime-status" className="text-xs text-slate-400">
        {t(`settings.share_runtime_${key}`, key)}
      </p>
    );
  }
  const configured = status.configured_enabled ?? status.enabled;
  const running = status.running ?? (status.lifecycle_state === 'running' || status.bind_state === 'pass');
  const failed = status.failed ?? (status.bind_state === 'fail' || status.lifecycle_state === 'failed');
  const key = failed ? 'failed' : !configured ? 'disabled' : running ? 'running' : 'pending';
  return (
    <div className="space-y-2" data-testid="share-runtime-controls">
      <p data-testid="share-runtime-status" className="text-xs text-slate-400">
        {t(`settings.share_runtime_${key}`, key)}
      </p>
      {configured && (
        <div className="flex flex-wrap items-center gap-2">
          <button
            type="button"
            onClick={onStart}
            disabled={lifecyclePending || running}
            className="rounded-lg border border-emerald-400/30 px-3 py-1.5 text-xs text-emerald-200 disabled:opacity-50"
          >
            {t('settings.share_runtime_start', 'Start')}
          </button>
          <button
            type="button"
            onClick={onStop}
            disabled={lifecyclePending || !running}
            className="rounded-lg border border-amber-400/30 px-3 py-1.5 text-xs text-amber-200 disabled:opacity-50"
          >
            {t('settings.share_runtime_stop', 'Stop')}
          </button>
        </div>
      )}
      {lifecycleError && (
        <div className="text-xs text-red-200" data-testid="share-runtime-operation-error">
          <p>{lifecycleError.message}</p>
          {lifecycleError.requestId && (
            <p className="mt-1">{t('errors.request_id', 'Request ID')}: {lifecycleError.requestId}</p>
          )}
        </div>
      )}
    </div>
  );
}
