import { useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { designSystem } from '../../../styles/designSystem';
import { useDevicesQuery } from '../../../hooks/datalink/useDevices';
import { useSourceRulesQuery } from '../../../hooks/datalink/useSourceRules';
import { useRuntimeStream } from '../../../hooks/datalink/useRuntimeStream';
import { DeviceStep } from './DeviceStep';
import { LaunchStep } from './LaunchStep';
import { OutputStep } from './OutputStep';
import { PipelinePreview } from './PipelinePreview';
import { PointsStep } from './PointsStep';
import { QUICKFLOW_STEPS, SAMPLE_DEVICE, SAMPLE_RULE, buildQuickFlowViewModel } from './quickflowModel';
import { StepRail } from './StepRail';
import type { Device } from '../../../types/datalink';
import type { QuickFlowStepId } from './quickflowTypes';

const gradientShell = 'linear-gradient(180deg, rgba(17,19,27,0.98), rgba(13,13,13,0.98))';
const themeTokens = designSystem.tokens;
const panelStyle = {
  borderRadius: 24,
  border: '1px solid rgba(122, 143, 216, 0.16)',
  background: 'rgba(18, 22, 34, 0.92)',
} as const;

export default function QuickFlowPage() {
  const { t } = useTranslation();
  const deviceQuery = useDevicesQuery();
  const [activeStep, setActiveStep] = useState<QuickFlowStepId>('device');
  const [selectedDeviceId, setSelectedDeviceId] = useState<string>('');

  useEffect(() => {
    if (!selectedDeviceId && deviceQuery.data?.[0]?.id) {
      setSelectedDeviceId(deviceQuery.data[0].id);
    }
  }, [deviceQuery.data, selectedDeviceId]);

  const liveDeviceId = selectedDeviceId || deviceQuery.data?.[0]?.id || null;
  const sourceRuleQuery = useSourceRulesQuery(liveDeviceId ? { device_id: liveDeviceId } : undefined);
  const runtime = useRuntimeStream({ deviceId: liveDeviceId, pointIds: ['point-40021', 'point-40022', 'point-40030'] });

  const devices = deviceQuery.data?.length ? deviceQuery.data : [SAMPLE_DEVICE];
  const selectedDevice = devices.find((device) => device.id === selectedDeviceId) ?? devices[0] ?? SAMPLE_DEVICE;
  const sourceRule = sourceRuleQuery.data?.[0] ?? SAMPLE_RULE;
  const model = useMemo(
    () =>
      buildQuickFlowViewModel({
        device: selectedDevice,
        sourceRule,
        runtimeState: runtime.connectionState,
        liveValues: runtime.liveValues,
      }),
    [runtime.connectionState, runtime.liveValues, selectedDevice, sourceRule],
  );

  if (deviceQuery.error) {
    return (
      <main style={{ minHeight: '100vh', background: gradientShell, padding: 24, color: themeTokens.colors.textPrimary }}>
        <section style={{ ...panelStyle, maxWidth: 720, margin: '0 auto', padding: 28 }}>
          <h1 style={{ margin: 0, fontSize: 30 }}>{t('quickflow.page.title')}</h1>
          <p style={{ marginTop: 14, color: themeTokens.colors.textSecondary }}>{t('quickflow.errors.devices')}</p>
          <button
            type="button"
            onClick={() => deviceQuery.refetch?.()}
            style={{ marginTop: 18, border: 0, borderRadius: 999, padding: '12px 18px', background: 'linear-gradient(135deg, #5d7cff, #9b5cf6)', color: '#fff', cursor: 'pointer' }}
          >
            {t('quickflow.actions.refreshDevices')}
          </button>
        </section>
      </main>
    );
  }

  return (
    <main
      style={{
        minHeight: '100vh',
        background: gradientShell,
        color: themeTokens.colors.textPrimary,
        fontFamily: `${themeTokens.font.family.sans}, "Noto Sans TC", sans-serif`,
        padding: 24,
      }}
    >
      <div style={{ width: 'min(1520px, calc(100vw - 48px))', margin: '0 auto', borderRadius: 28, padding: 24, border: '1px solid rgba(122,143,216,0.14)', background: 'rgba(8, 10, 16, 0.82)', boxShadow: themeTokens.shadow.glow.blue }}>
        <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: 18, marginBottom: 20 }}>
          <div>
            <span style={{ display: 'inline-flex', padding: '6px 12px', borderRadius: 999, fontSize: 12, letterSpacing: '0.08em', textTransform: 'uppercase', background: themeTokens.colors.primaryLight, color: '#dfe7ff', border: '1px solid rgba(122,143,216,0.18)' }}>
              {t('quickflow.page.badge')}
            </span>
            <h1 style={{ margin: '12px 0 0', fontSize: 34 }}>{t('quickflow.page.title')}</h1>
            <p style={{ marginTop: 10, maxWidth: 760, color: themeTokens.colors.textSecondary }}>{t('quickflow.page.description')}</p>
          </div>
          <div style={{ ...panelStyle, minWidth: 320, padding: 18 }}>
            <div style={{ color: themeTokens.colors.textSecondary, fontSize: 12 }}>{t('quickflow.page.currentDevice')}</div>
            <strong style={{ display: 'block', marginTop: 10, fontSize: 21 }}>{selectedDevice.name}</strong>
            <div style={{ marginTop: 8, color: '#dbe3ff' }}>{model.endpoint}</div>
            <div style={{ marginTop: 6, color: '#7fe4b3' }}>{model.unitLabel}</div>
          </div>
        </header>

        <section style={{ ...panelStyle, padding: 22, background: 'linear-gradient(135deg, rgba(23,27,42,0.96), rgba(12,14,22,0.96))' }}>
          <div style={{ display: 'grid', gap: 18, gridTemplateColumns: '1.6fr 1fr', alignItems: 'center' }}>
            <div>
              <div style={{ color: themeTokens.colors.textSecondary, fontSize: 12 }}>{t('quickflow.hero.kicker')}</div>
              <h2 style={{ margin: '8px 0 0', fontSize: 28 }}>{t('quickflow.hero.title')}</h2>
              <p style={{ marginTop: 10, color: '#9ba6c9' }}>{t('quickflow.hero.description')}</p>
            </div>
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, minmax(0, 1fr))', gap: 12 }}>
              {model.metrics.map((metric) => (
                <div key={metric.labelKey} style={{ padding: 16, borderRadius: 18, background: 'rgba(9,13,22,0.92)', border: '1px solid rgba(122,143,216,0.12)' }}>
                  <div style={{ color: '#9ba6c9', fontSize: 12 }}>{t(metric.labelKey)}</div>
                  <strong style={{ display: 'block', marginTop: 8, fontSize: 22 }}>{metric.value}</strong>
                  <div style={{ marginTop: 6, color: '#8b96b7', fontSize: 12 }}>{metric.detail}</div>
                </div>
              ))}
            </div>
          </div>
          <div style={{ marginTop: 18 }}>
            <StepRail steps={QUICKFLOW_STEPS} activeStep={activeStep} onStepChange={setActiveStep} />
          </div>
        </section>

        <div style={{ display: 'grid', gridTemplateColumns: '7fr 3fr', gap: 18, alignItems: 'start', marginTop: 18 }}>
          <section style={{ ...panelStyle, padding: 22 }}>
            {activeStep === 'device' ? <DeviceStep devices={devices as Device[]} metrics={model.metrics} selectedDeviceId={selectedDevice.id} model={model} onSelectDevice={setSelectedDeviceId} /> : null}
            {activeStep === 'points' ? <PointsStep model={model} /> : null}
            {activeStep === 'output' ? <OutputStep model={model} /> : null}
            {activeStep === 'launch' ? <LaunchStep model={model} /> : null}
            {sourceRuleQuery.error ? (
              <div style={{ marginTop: 18, padding: 16, borderRadius: 18, background: 'rgba(255,191,71,0.12)', color: '#ffe0a2' }}>
                {t('quickflow.errors.rules')}
              </div>
            ) : null}
          </section>
          <PipelinePreview runtimeState={model.runtimeState} nodes={model.pipelineNodes} blocker={model.blocker} />
        </div>
      </div>
    </main>
  );
}
