import { useTranslation } from 'react-i18next';
import { Cpu, Server } from 'lucide-react';
import type { StudioV2RuntimeContextDevice } from '../../../types/studioV2RuntimeContext';

interface FocusedDeviceHeaderProps {
  selectedDeviceId: string | null;
  selectedDevice: StudioV2RuntimeContextDevice | null;
  devices: StudioV2RuntimeContextDevice[];
  onSelectDevice: (deviceId: string) => void;
}

export function FocusedDeviceHeader({
  selectedDeviceId,
  selectedDevice,
  devices,
  onSelectDevice,
}: FocusedDeviceHeaderProps) {
  const { t } = useTranslation('runtime-dashboard');
  const isRunning = selectedDevice?.running ?? false;
  const isUnavailable = selectedDevice?.availability_status === 'unavailable';

  return (
    <section
      className="rounded-3xl border border-slate-800/80 bg-slate-900/80 p-6 shadow-2xl backdrop-blur-md"
      data-testid="runtime-dashboard-header"
    >
      <div className="flex flex-col gap-5 lg:flex-row lg:items-center lg:justify-between">
        <div className="flex items-start gap-4">
          <div className="flex h-14 w-14 shrink-0 items-center justify-center rounded-2xl border border-cyan-500/30 bg-gradient-to-br from-cyan-500/20 to-slate-900 text-cyan-300 shadow-lg shadow-cyan-950/50">
            <Cpu className="h-7 w-7" />
          </div>
          <div className="space-y-1">
            <div className="flex flex-wrap items-center gap-2">
              <p className="text-[11px] font-semibold uppercase tracking-[0.24em] text-cyan-400">
                {t('header.eyebrow', 'Runtime Dashboard')}
              </p>
              {selectedDevice && (
                <span
                  className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-semibold ${
                    isRunning
                      ? 'border border-emerald-500/30 bg-emerald-500/10 text-emerald-300'
                      : 'border border-amber-500/30 bg-amber-500/10 text-amber-300'
                  }`}
                >
                  <span
                    className={`h-1.5 w-1.5 rounded-full ${
                      isRunning ? 'bg-emerald-400 animate-pulse' : 'bg-amber-400'
                    }`}
                  />
                  {isRunning ? '運轉中' : '已停止'}
                </span>
              )}
            </div>

            <h1 className="text-2xl font-bold tracking-tight text-slate-50 sm:text-3xl">
              {selectedDevice?.name ?? t('header.noDevice', 'Choose a device')}
            </h1>

            <p className="text-xs text-slate-400">
              {selectedDevice ? (
                isUnavailable ? (
                  <span className="text-amber-300 font-medium">
                    {t('header.deviceMetaUnavailable', 'Unavailable: runtime device context is unavailable.')}
                  </span>
                ) : (
                  <span className="inline-flex items-center gap-2">
                    <span className="rounded bg-slate-800 px-1.5 py-0.5 font-mono text-[11px] text-slate-300">
                      {t('header.deviceMeta', 'Protocol {{protocol}}', {
                        protocol: selectedDevice.protocol,
                      })}
                    </span>
                    <span className="text-slate-500">•</span>
                    <span>即時遙測串流已啟動</span>
                  </span>
                )
              ) : (
                t(
                  'header.deviceMetaEmpty',
                  'Use the committed device context to inspect live runtime health.',
                )
              )}
            </p>
          </div>
        </div>

        {/* 設備切換按鈕區 */}
        {devices.length > 0 && (
          <div className="flex flex-col gap-1.5">
            <span className="text-[11px] font-semibold uppercase tracking-wider text-slate-500">
              切換監控設備
            </span>
            <div className="flex flex-wrap gap-2" data-testid="runtime-dashboard-device-switcher">
              {devices.map((device) => {
                const isSelected = device.device_id === selectedDeviceId;
                const unavailable = device.availability_status === 'unavailable';
                return (
                  <button
                    key={device.device_id}
                    type="button"
                    aria-label={device.name}
                    onClick={() => onSelectDevice(device.device_id)}
                    className={`flex items-center gap-2 rounded-xl border px-3.5 py-2 text-xs font-medium transition-all ${
                      isSelected
                        ? 'border-cyan-400 bg-cyan-500/15 text-cyan-100 shadow-md shadow-cyan-950/40'
                        : 'border-slate-800 bg-slate-950/60 text-slate-300 hover:border-slate-700 hover:bg-slate-900/80 hover:text-slate-100'
                    }`}
                  >
                    <Server className={`h-3.5 w-3.5 ${isSelected ? 'text-cyan-400' : 'text-slate-500'}`} />
                    <span>{device.name}</span>
                    {unavailable && (
                      <div className="flex flex-col text-left">
                        <span className="rounded bg-amber-500/20 px-1 py-0.2 text-[10px] font-semibold text-amber-300">
                          {t('header.unavailable', 'Unavailable')}
                        </span>
                        {device.availability_reason && (
                          <span className="text-[10px] text-slate-400">
                            {t('header.unavailableAction', 'Review the device setup in Studio V2 and retry.')}
                          </span>
                        )}
                      </div>
                    )}
                  </button>
                );
              })}
            </div>
          </div>
        )}
      </div>

      {selectedDevice && !isUnavailable && (
        <div className="mt-5 flex flex-wrap items-center justify-between gap-3 rounded-2xl border border-cyan-500/20 bg-gradient-to-r from-cyan-500/10 via-slate-900/60 to-emerald-500/10 px-4 py-3 text-xs">
          <div className="flex items-center gap-2.5 text-slate-200">
            <span className="relative flex h-2 w-2">
              <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-400 opacity-75" />
              <span className="relative inline-flex h-2 w-2 rounded-full bg-emerald-500" />
            </span>
            <span className="font-semibold text-emerald-300">資料採集與儲存進行中</span>
            <span className="text-slate-500">•</span>
            <span className="text-slate-300">
              閘道正持續從 PLC 讀取數據，經由 Datalink 映射並即時寫入資料庫持久化儲存。
            </span>
          </div>
          <div className="flex items-center gap-2 text-[11px] font-mono text-cyan-300">
            <span>通訊協定: {selectedDevice.protocol}</span>
          </div>
        </div>
      )}
    </section>
  );
}
