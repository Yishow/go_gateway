import React from 'react';
import { usePreviewStream, type StepResult } from '../../../hooks/usePreviewStream';
import './LivePreviewPanel.css';

/**
 * LivePreviewPanel 組件屬性
 */
interface LivePreviewPanelProps {
  /** Mapping ID */
  mappingId?: string;
  /** 是否自動開始預覽 */
  autoStart?: boolean;
  /** 自訂樣式類名 */
  className?: string;
}

/**
 * 即時預覽面板
 *
 * 顯示 Mapping 的即時 raw → transform → final 值
 */
export const LivePreviewPanel: React.FC<LivePreviewPanelProps> = ({
  mappingId = '',
  autoStart = false,
  className = '',
}) => {
  const {
    latestEvent,
    connectionState,
    error,
    connect,
    disconnect,
    isConnected,
  } = usePreviewStream({
    mappingId,
    autoConnect: autoStart && !!mappingId,
  });

  /**
   * 格式化值為字串顯示
   */
  const formatValue = (value: unknown): string => {
    if (value === null || value === undefined) {
      return '-';
    }
    if (typeof value === 'number') {
      return value.toLocaleString(undefined, { maximumFractionDigits: 4 });
    }
    return String(value);
  };

  /**
   * 取得連接狀態顯示文字
   */
  const getConnectionStatusText = (): string => {
    switch (connectionState) {
      case 'connecting':
        return '連線中...';
      case 'connected':
        return '已連線';
      case 'disconnected':
        return '未連線';
      case 'error':
        return '連線錯誤';
      default:
        return '未知';
    }
  };

  /**
   * 取得連接狀態顏色
   */
  const getConnectionStatusColor = (): string => {
    switch (connectionState) {
      case 'connected':
        return 'status-connected';
      case 'connecting':
        return 'status-connecting';
      case 'error':
        return 'status-error';
      default:
        return 'status-disconnected';
    }
  };

  /**
   * 渲染步驟結果
   */
  const renderStepResults = (steps: StepResult[]): React.ReactNode => {
    if (!steps || steps.length === 0) {
      return <div className="preview-no-steps">無轉換步驟</div>;
    }

    return (
      <div className="preview-steps">
        {steps.map((step, index) => (
          <div key={index} className="preview-step">
            <div className="step-header">
              <span className="step-index">步驟 {step.step_index + 1}</span>
              <span className="step-type">{step.step_type}</span>
            </div>
            <div className="step-values">
              <span className="step-input">{formatValue(step.input)}</span>
              <span className="step-arrow">→</span>
              <span className="step-output">{formatValue(step.output)}</span>
            </div>
            {step.error && <div className="step-error">{step.error}</div>}
          </div>
        ))}
      </div>
    );
  };

  return (
    <div className={`live-preview-panel ${className}`}>
      <div className="preview-header">
        <h3>即時預覽</h3>
        <div className={`connection-status ${getConnectionStatusColor()}`}>
          <span className="status-dot" />
          <span className="status-text">{getConnectionStatusText()}</span>
        </div>
      </div>

      {error && <div className="preview-error">{error}</div>}

      {!mappingId && (
        <div className="preview-placeholder">
          請選擇 Mapping 以開始預覽
        </div>
      )}

      {mappingId && !isConnected && (
        <div className="preview-actions">
          <button
            className="btn-start-preview"
            onClick={connect}
            disabled={connectionState === 'connecting'}
          >
            {connectionState === 'connecting' ? '連線中...' : '▶ 開始即時預覽'}
          </button>
        </div>
      )}

      {isConnected && (
        <div className="preview-actions">
          <button className="btn-stop-preview" onClick={disconnect}>
            ⏹ 停止預覽
          </button>
        </div>
      )}

      {latestEvent?.type === 'preview' && (
        <div className="preview-content">
          <div className="preview-value-row">
            <label>原始值 (Raw)</label>
            <span className="value raw-value">{formatValue(latestEvent.raw_value)}</span>
          </div>

          {renderStepResults(latestEvent.steps || [])}

          <div className="preview-value-row final">
            <label>最終值 (Final)</label>
            <span className="value final-value">{formatValue(latestEvent.final_value)}</span>
          </div>

          <div className="preview-meta">
            <span className="quality">品質: {latestEvent.quality}</span>
            <span className="timestamp">{latestEvent.timestamp}</span>
          </div>
        </div>
      )}

      {isConnected && latestEvent?.type === 'heartbeat' && (
        <div className="preview-waiting">
          等待資料中...
        </div>
      )}
    </div>
  );
};

export default LivePreviewPanel;
