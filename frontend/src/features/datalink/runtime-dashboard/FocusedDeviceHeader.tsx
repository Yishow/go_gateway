import { useTranslation } from 'react-i18next';
import type { Device } from '../../../types/datalink';

interface FocusedDeviceHeaderProps {
  selectedDeviceId: string | null;
  selectedDevice: Device | null;
  devices: Device[];
  onSelectDevice: (deviceId: string) => void;
}

export function FocusedDeviceHeader({
  selectedDeviceId,
  selectedDevice,
  devices,
  onSelectDevice,
}: FocusedDeviceHeaderProps) {
  const { t } = useTranslation('runtime-dashboard');

  return (
    <section
      className="rounded-3xl border border-slate-800 bg-slate-900/70 p-6 shadow-[0_20px_60px_rgba(15,23,42,0.35)]"
      data-testid="runtime-dashboard-header"
    >
      <div className="flex flex-col gap-5 lg:flex-row lg:items-end lg:justify-between">
        <div className="space-y-2">
          <p className="text-xs font-semibold uppercase tracking-[0.24em] text-cyan-300">
            {t('header.eyebrow', 'Runtime Dashboard')}
          </p>
          <div className="space-y-1">
            <h1 className="text-3xl font-semibold text-slate-50">
              {selectedDevice?.name ?? t('header.noDevice', 'Choose a device')}
            </h1>
            <p className="text-sm text-slate-300">
              {selectedDevice
                ? t('header.deviceMeta', 'Protocol {{protocol}}', {
                    protocol: selectedDevice.protocol,
                  })
                : t(
                    'header.deviceMetaEmpty',
                    'Use the committed device context to inspect live runtime health.',
                  )}
            </p>
          </div>
        </div>

        <div className="flex flex-wrap gap-2" data-testid="runtime-dashboard-device-switcher">
          {devices.map((device) => {
            const isSelected = device.id === selectedDeviceId;
            return (
              <button
                key={device.id}
                type="button"
                onClick={() => onSelectDevice(device.id)}
                className={`rounded-xl border px-3 py-2 text-sm font-medium transition ${
                  isSelected
                    ? 'border-cyan-400/60 bg-cyan-500/10 text-cyan-100'
                    : 'border-slate-700 bg-slate-950/60 text-slate-200 hover:border-slate-500 hover:text-slate-50'
                }`}
              >
                {device.name}
              </button>
            );
          })}
        </div>
      </div>
    </section>
  );
}
