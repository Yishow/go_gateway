import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { Device } from '../../state/types';
import { getStagesForProtocol } from '../../state/protocols';
import { ReadinessStages } from './ReadinessStages';
import { Button } from '../../components/Button';
import { Icon } from '../../components/Icon';
import { SectionCard } from '../../components/SectionCard';

export interface ConnectionTestPanelProps {
  /** 當前所選的設備資料 */
  device: Device;
  /** 點按執行測試時的回調 */
  onRunTest: (deviceId: string) => void | Promise<void>;
}

/**
 * 設備連線測試與診斷面板 (ConnectionTestPanel)
 * 
 * 落地需求 **Send payload preview block**。
 * 展示預計傳送的設定 JSON Payload、連線與握手狀態的即時推進，以及測試成功時的綠色卡片通知。
 */
export const ConnectionTestPanel: React.FC<ConnectionTestPanelProps> = ({ device, onRunTest }) => {
  const { t } = useTranslation('workbench-v2');

  const { protocol, config, test } = device;

  // 驗證欄位是否已填寫完整以決定測試按鈕是否啟用
  const isConfigValid = (): boolean => {
    if (protocol === 'modbus_rtu') {
      return !!config.port && config.baud !== undefined && config.baud !== '';
    }
    if (protocol === 'mqtt') {
      return !!config.broker && !!config.client_id;
    }
    // TCP 系
    return !!config.host && config.port !== undefined && config.port !== '';
  };

  const isTesting = test?.status === 'running';
  const isSuccess = test?.status === 'success';
  const isBtnDisabled = !isConfigValid() || isTesting;

  // 取得當前顯示的 stages 序列
  // 如果 test 已經存在，應使用 test.stages 的狀態，否則使用預設 pending 序列
  const currentStages = React.useMemo(() => {
    const defaultStages = getStagesForProtocol(protocol);
    if (!test || !test.stages) return defaultStages;
    if ('connect' in test.stages || 'probe' in test.stages) {
      return [
        {
          id: 'connect',
          label: 'step1.stages.connect',
          group: 'connect' as const,
          status: test.stages.connect?.status ?? 'pending',
          latency_ms: test.stages.connect?.latency_ms,
          message: test.stages.connect?.message,
        },
        {
          id: 'probe',
          label: 'step1.stages.probe',
          group: 'probe' as const,
          status: test.stages.probe?.status ?? 'pending',
          latency_ms: test.stages.probe?.latency_ms,
          message: test.stages.probe?.message,
        },
      ];
    }

    return defaultStages.map((stage) => {
      const liveStage = test.stages[stage.id];
      return {
        ...stage,
        status: liveStage ? liveStage.status : 'pending',
        latency_ms: liveStage ? liveStage.latency_ms : undefined,
        message: liveStage ? liveStage.message : undefined,
      };
    });
  }, [protocol, test]);

  // Payload 預覽內容
  const payloadPreview = React.useMemo(() => {
    const payload = {
      protocol,
      ...config,
    };
    return JSON.stringify(payload, null, 2);
  }, [protocol, config]);

  return (
    <SectionCard
      title={t('step1.panels.diagnostic')}
      icon="device"
      className="flex h-full flex-col justify-between"
      data-testid="connection-test-panel"
    >
      <div className="flex-1 space-y-5">
        {/* Payload 預覽塊 */}
        <div className="space-y-2">
          <span className="text-xs font-semibold text-slate-400">
            {t('step1.panels.payload_preview')}
          </span>
          <pre
            className="max-h-[140px] overflow-auto rounded-lg border border-slate-800 bg-slate-950/80 p-3 font-mono text-[11px] text-slate-400 leading-relaxed scrollbar-thin"
            data-testid="payload-preview"
          >
            {payloadPreview}
          </pre>
          <span className="text-[10px] text-slate-500 block">
            {t('step1.panels.payload_hint')}
          </span>
        </div>

        {/* Readiness 診斷步驟 */}
        <div className="space-y-2">
          <span className="text-xs font-semibold text-slate-400 block">
            {t('step1.panels.stages_title')}
          </span>
          <div className="rounded-lg border border-slate-800 bg-slate-900/10 p-3.5">
            <ReadinessStages stages={currentStages} testStatus={test?.status} />
          </div>
        </div>

        {/* 測試成功後的綠色卡片通知 */}
        {isSuccess && (
          <div
            className="flex items-start gap-3 rounded-lg border border-emerald-500/25 bg-emerald-500/10 p-3 text-emerald-400"
            data-testid="success-readiness-card"
          >
            <div className="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-emerald-500 text-slate-950 mt-0.5">
              <Icon name="check" className="h-3 w-3 stroke-[2.5]" />
            </div>
            <div className="space-y-1">
              <span className="text-xs font-semibold block text-slate-200">
                {t('step1.panels.test_success_title')}
              </span>
              <p className="text-[11px] text-slate-400 leading-normal">
                {t('step1.panels.test_success_desc', { latency: test.latency_ms ?? 0 })}
              </p>
            </div>
          </div>
        )}
      </div>

      {/* 底部動作按鈕 */}
      <div className="mt-5 border-t border-slate-800/80 pt-4 flex justify-end">
        <Button
          onClick={() => onRunTest(device.id)}
          disabled={isBtnDisabled}
          variant={isSuccess ? 'secondary' : 'primary'}
          className="w-full sm:w-auto"
          data-testid="run-test-button"
        >
          {isTesting ? (
            <span className="flex items-center gap-1.5 justify-center">
              <Icon name="refresh" className="h-3.5 w-3.5 animate-spin" />
              {t('step1.buttons.testing')}
            </span>
          ) : (
            <span className="flex items-center gap-1.5 justify-center">
              <Icon name="play" className="h-3.5 w-3.5" />
              {t('step1.buttons.run_test')}
            </span>
          )}
        </Button>
      </div>
    </SectionCard>
  );
};
