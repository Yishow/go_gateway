import { useTranslation } from 'react-i18next';
import type { Device } from '../../../types/datalink';
import type { QuickFlowMetric, QuickFlowViewModel } from './quickflowTypes';

interface DeviceStepProps {
  devices: Device[];
  metrics: QuickFlowMetric[];
  selectedDeviceId: string;
  model: QuickFlowViewModel;
  onSelectDevice: (deviceId: string) => void;
}

function MetricCard({ metric }: { metric: QuickFlowMetric }) {
  const { t } = useTranslation();
  return (
    <div style={{ padding: 16, borderRadius: 18, background: 'rgba(9, 13, 22, 0.92)', border: '1px solid rgba(122,143,216,0.12)' }}>
      <div style={{ color: '#9ba6c9', fontSize: 12 }}>{t(metric.labelKey)}</div>
      <strong style={{ display: 'block', marginTop: 8, fontSize: 24 }}>{metric.value}</strong>
      <div style={{ marginTop: 6, color: '#8b96b7', fontSize: 12 }}>{metric.detail}</div>
    </div>
  );
}

export function DeviceStep({ devices, metrics, selectedDeviceId, model, onSelectDevice }: DeviceStepProps) {
  const { t } = useTranslation();

  return (
    <section style={{ display: 'grid', gap: 18 }}>
      <div>
        <h2 style={{ margin: 0, fontSize: 28 }}>{t('quickflow.sections.device.title')}</h2>
        <p style={{ marginTop: 8, color: '#9ba6c9' }}>{t('quickflow.sections.device.description')}</p>
      </div>
      <div style={{ display: 'grid', gap: 14 }}>
        {devices.map((device) => {
          const active = device.id === selectedDeviceId;
          return (
            <button
              key={device.id}
              type="button"
              onClick={() => onSelectDevice(device.id)}
              style={{
                padding: '16px 18px',
                borderRadius: 20,
                border: active ? '1px solid rgba(153,120,255,0.4)' : '1px solid rgba(122,143,216,0.12)',
                background: active ? 'linear-gradient(135deg, rgba(93,124,255,0.18), rgba(155,92,246,0.18))' : 'rgba(10, 14, 22, 0.9)',
                color: '#f5f7ff',
                textAlign: 'left',
                cursor: 'pointer',
              }}
            >
              <strong style={{ display: 'block', fontSize: 18 }}>{device.name}</strong>
              <div style={{ marginTop: 8, color: '#9ba6c9', fontSize: 14 }}>{model.endpoint}</div>
              <div style={{ marginTop: 6, color: '#7fe4b3', fontSize: 13 }}>{model.unitLabel}</div>
            </button>
          );
        })}
      </div>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(4, minmax(0, 1fr))', gap: 12 }}>
        {metrics.map((metric) => <MetricCard key={metric.labelKey} metric={metric} />)}
      </div>
    </section>
  );
}
