import React, { useState } from 'react';
import { mappingAPI } from '../../../../services/datalink';
import type { WizardFormData } from '../MappingWizard';
import './Steps.css';

/**
 * CompleteStep 組件屬性
 */
interface CompleteStepProps {
  /** 表單資料 */
  formData: WizardFormData;
  /** 完成回呼 */
  onComplete: () => void;
}

/**
 * 完成步驟
 */
export const CompleteStep: React.FC<CompleteStepProps> = ({
  formData,
  onComplete,
}) => {
  const [saving, setSaving] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSave = async () => {
    try {
      setSaving(true);
      setError(null);

      // 建立映射
      await mappingAPI.create({
        point_id: formData.pointId,
        tag_id: formData.tagId,
        transform_pipeline: formData.transformPipeline,
        enabled: true,
      });

      setSuccess(true);
      onComplete();
    } catch (err) {
      setError(err instanceof Error ? err.message : '建立映射失敗');
    } finally {
      setSaving(false);
    }
  };

  if (success) {
    return (
      <div className="step-complete success">
        <div className="success-icon">🎉</div>
        <h3>映射建立成功！</h3>
        <p>資料將開始自動從點位流入標籤。</p>
      </div>
    );
  }

  return (
    <div className="step-complete">
      <div className="complete-icon">🚀</div>
      <h3>準備完成</h3>
      <p>所有設定已完成，點擊下方按鈕建立並啟用映射。</p>

      {error && <div className="step-error">{error}</div>}

      <div className="complete-summary">
        <h4>最終確認</h4>
        <dl>
          <dt>來源點位</dt>
          <dd>{formData.pointId}</dd>
          <dt>目標標籤</dt>
          <dd>{formData.tagId}</dd>
          <dt>轉換步驟</dt>
          <dd>
            {formData.transformPipeline.length === 0
              ? '無（直接傳遞）'
              : `${formData.transformPipeline.length} 個步驟`}
          </dd>
        </dl>
      </div>

      <button
        className="btn-activate"
        onClick={handleSave}
        disabled={saving}
      >
        {saving ? '建立中...' : '🚀 建立並啟用映射'}
      </button>
    </div>
  );
};

export default CompleteStep;
