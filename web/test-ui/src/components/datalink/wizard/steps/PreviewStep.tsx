import React from 'react';
import { LivePreviewPanel } from '../LivePreviewPanel';
import type { WizardFormData } from '../MappingWizard';
import './Steps.css';

/**
 * PreviewStep 組件屬性
 */
interface PreviewStepProps {
  /** 表單資料 */
  formData: WizardFormData;
  /** 確認回呼 */
  onConfirm: () => void;
  /** 是否已確認 */
  confirmed: boolean;
  /** 錯誤訊息 */
  error?: string;
}

/**
 * 預覽確認步驟
 */
export const PreviewStep: React.FC<PreviewStepProps> = ({
  formData,
  onConfirm,
  confirmed,
  error,
}) => {
  const mappingId = formData.pointId && formData.tagId 
    ? `${formData.pointId}-${formData.tagId}` 
    : '';

  return (
    <div className="step-preview">
      <div className="step-info-box">
        <p>
          開始即時預覽以驗證轉換結果。確認無誤後點擊「確認結果」按鈕。
        </p>
      </div>

      {error && <div className="step-error">{error}</div>}

      <div className="preview-summary">
        <h4>映射摘要</h4>
        <dl>
          <dt>設備 ID</dt>
          <dd>{formData.deviceId || '-'}</dd>
          <dt>點位 ID</dt>
          <dd>{formData.pointId || '-'}</dd>
          <dt>標籤 ID</dt>
          <dd>{formData.tagId || '-'}</dd>
          <dt>轉換步驟</dt>
          <dd>{formData.transformPipeline.length} 個</dd>
        </dl>
      </div>

      <div className="preview-panel-wrapper">
        <LivePreviewPanel
          mappingId={mappingId}
          autoStart={true}
        />
      </div>

      <div className="confirm-section">
        {confirmed ? (
          <div className="confirmed-badge">
            <span>✓</span> 已確認預覽結果
          </div>
        ) : (
          <button
            className="btn-confirm"
            onClick={onConfirm}
          >
            確認結果
          </button>
        )}
      </div>
    </div>
  );
};

export default PreviewStep;
