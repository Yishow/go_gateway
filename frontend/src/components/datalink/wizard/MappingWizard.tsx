import React, { useState, useCallback } from 'react';
import './MappingWizard.css';
import StepIndicator from './StepIndicator';
import LivePreviewPanel from './LivePreviewPanel';
import DeviceStep from './steps/DeviceStep';
import PointStep from './steps/PointStep';
import TagStep from './steps/TagStep';
import TransformStep from './steps/TransformStep';
import PreviewStep from './steps/PreviewStep';
import CompleteStep from './steps/CompleteStep';
import { useToast } from '../../../contexts/ToastContext';
import type { TransformStep as TransformStepType } from '../../../types/datalink';

/**
 * Wizard 步驟定義
 */
export interface WizardStep {
  id: string;
  title: string;
  description: string;
  icon: string;
}

/**
 * Wizard 表單資料
 */
export interface WizardFormData {
  deviceId: string;
  pointId: string;
  tagId: string;
  transformPipeline: TransformStepType[];
  confirmed: boolean;
}

/**
 * Wizard 步驟列表
 */
const WIZARD_STEPS: WizardStep[] = [
  { id: 'device', title: '設備選擇', description: '選擇或新增設備', icon: '📱' },
  { id: 'point', title: '點位配置', description: '配置資料點位', icon: '📍' },
  { id: 'tag', title: '標籤選擇', description: '選擇目標標籤', icon: '🏷️' },
  { id: 'transform', title: '轉換設定', description: '配置轉換管線', icon: '⚙️' },
  { id: 'preview', title: '預覽確認', description: '驗證轉換結果', icon: '👁️' },
  { id: 'complete', title: '完成', description: '啟用映射', icon: '✅' },
];

/**
 * 初始表單資料
 */
const INITIAL_FORM_DATA: WizardFormData = {
  deviceId: '',
  pointId: '',
  tagId: '',
  transformPipeline: [],
  confirmed: false,
};

/**
 * MappingWizard 組件屬性
 */
interface MappingWizardProps {
  /** 完成時回呼 */
  onComplete?: (data: WizardFormData) => void;
  /** 取消時回呼 */
  onCancel?: () => void;
}

/**
 * Mapping 建立精靈
 *
 * 6 步驟導引式工作流程：設備 → 點位 → 標籤 → 轉換 → 預覽 → 完成
 */
export const MappingWizard: React.FC<MappingWizardProps> = ({
  onComplete,
  onCancel,
}) => {
  const [currentStep, setCurrentStep] = useState(0);
  const [formData, setFormData] = useState<WizardFormData>(INITIAL_FORM_DATA);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const { showSuccess, showError, showInfo } = useToast();

  /**
   * 更新表單資料
   */
  const updateFormData = useCallback((updates: Partial<WizardFormData>) => {
    setFormData((prev) => ({ ...prev, ...updates }));
    // 清除相關錯誤
    const keys = Object.keys(updates);
    setErrors((prev) => {
      const newErrors = { ...prev };
      keys.forEach((key) => delete newErrors[key]);
      return newErrors;
    });
  }, []);

  /**
   * 驗證當前步驟
   */
  const validateStep = useCallback((stepIndex: number): boolean => {
    const newErrors: Record<string, string> = {};

    switch (stepIndex) {
      case 0: // 設備選擇
        if (!formData.deviceId) {
          newErrors.deviceId = '請選擇一個設備';
        }
        break;
      case 1: // 點位配置
        if (!formData.pointId) {
          newErrors.pointId = '請選擇或新增一個點位';
        }
        break;
      case 2: // 標籤選擇
        if (!formData.tagId) {
          newErrors.tagId = '請選擇或新增一個標籤';
        }
        break;
      case 3: // 轉換設定
        // 轉換管線可以為空（pass-through）
        break;
      case 4: // 預覽確認
        if (!formData.confirmed) {
          newErrors.confirmed = '請確認預覽結果';
        }
        break;
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  }, [formData]);

  /**
   * 前往下一步
   */
  const goNext = useCallback(() => {
    if (validateStep(currentStep)) {
      if (currentStep < WIZARD_STEPS.length - 1) {
        setCurrentStep((prev) => prev + 1);
      }
    }
  }, [currentStep, validateStep]);

  /**
   * 返回上一步
   */
  const goPrev = useCallback(() => {
    if (currentStep > 0) {
      setCurrentStep((prev) => prev - 1);
    }
  }, [currentStep]);

  /**
   * 跳轉至指定步驟（僅允許跳回已完成的步驟）
   */
  const goToStep = useCallback((stepIndex: number) => {
    if (stepIndex <= currentStep) {
      setCurrentStep(stepIndex);
    }
  }, [currentStep]);

  /**
   * 完成精靈
   */
  const handleComplete = useCallback(() => {
    if (validateStep(currentStep)) {
      onComplete?.(formData);
    }
  }, [currentStep, formData, onComplete, validateStep]);

  /**
   * 儲存草稿
   */
  const saveDraft = useCallback(() => {
    const draft = {
      formData,
      currentStep,
      savedAt: new Date().toISOString(),
    };
    localStorage.setItem('mapping-wizard-draft', JSON.stringify(draft));
    showSuccess('草稿已儲存');
  }, [formData, currentStep, showSuccess]);

  /**
   * 載入草稿
   */
  const loadDraft = useCallback(() => {
    const draftJson = localStorage.getItem('mapping-wizard-draft');
    if (draftJson) {
      try {
        const draft = JSON.parse(draftJson);
        setFormData(draft.formData);
        setCurrentStep(draft.currentStep);
        showInfo(`已載入草稿 (儲存於 ${new Date(draft.savedAt).toLocaleString()})`);
      } catch {
        showError('載入草稿失敗');
      }
    } else {
      showInfo('沒有找到草稿');
    }
  }, [showInfo, showError]);

  /**
   * 渲染當前步驟內容
   */
  const renderStepContent = () => {
    switch (currentStep) {
      case 0:
        return (
          <DeviceStep
            selectedDeviceId={formData.deviceId}
            onSelect={(deviceId) => updateFormData({ deviceId })}
            error={errors.deviceId}
          />
        );
      case 1:
        return (
          <PointStep
            deviceId={formData.deviceId}
            selectedPointId={formData.pointId}
            onSelect={(pointId: string) => updateFormData({ pointId })}
            error={errors.pointId}
          />
        );
      case 2:
        return (
          <TagStep
            selectedTagId={formData.tagId}
            onSelect={(tagId: string) => updateFormData({ tagId })}
            error={errors.tagId}
          />
        );
      case 3:
        return (
          <TransformStep
            pipeline={formData.transformPipeline}
            onChange={(transformPipeline: TransformStepType[]) => updateFormData({ transformPipeline })}
          />
        );
      case 4:
        return (
          <PreviewStep
            formData={formData}
            onConfirm={() => updateFormData({ confirmed: true })}
            confirmed={formData.confirmed}
            error={errors.confirmed}
          />
        );
      case 5:
        return (
          <CompleteStep
            formData={formData}
            onComplete={handleComplete}
          />
        );
      default:
        return null;
    }
  };

  const isFirstStep = currentStep === 0;
  const isLastStep = currentStep === WIZARD_STEPS.length - 1;

  return (
    <div className="mapping-wizard">
      <div className="wizard-header">
        <h1>建立新映射</h1>
        <div className="wizard-actions">
          <button className="btn-secondary" onClick={loadDraft}>
            載入草稿
          </button>
          <button className="btn-secondary" onClick={saveDraft}>
            儲存草稿
          </button>
          {onCancel && (
            <button className="btn-cancel" onClick={onCancel}>
              取消
            </button>
          )}
        </div>
      </div>

      <StepIndicator
        steps={WIZARD_STEPS}
        currentStep={currentStep}
        onStepClick={goToStep}
      />

      <div className="wizard-body">
        <div className="wizard-content">
          <div className="step-header">
            <span className="step-icon">{WIZARD_STEPS[currentStep].icon}</span>
            <div>
              <h2>{WIZARD_STEPS[currentStep].title}</h2>
              <p>{WIZARD_STEPS[currentStep].description}</p>
            </div>
          </div>

          <div className="step-content">
            {renderStepContent()}
          </div>

          <div className="wizard-navigation">
            <button
              className="btn-prev"
              onClick={goPrev}
              disabled={isFirstStep}
            >
              ← 上一步
            </button>
            
            {!isLastStep ? (
              <button
                className="btn-next"
                onClick={goNext}
              >
                下一步 →
              </button>
            ) : (
              <button
                className="btn-complete"
                onClick={handleComplete}
              >
                ✓ 完成並啟用
              </button>
            )}
          </div>
        </div>

        <div className="wizard-sidebar">
          <LivePreviewPanel
            mappingId={formData.pointId && formData.tagId ? `${formData.pointId}-${formData.tagId}` : ''}
            autoStart={currentStep >= 4}
          />
        </div>
      </div>
    </div>
  );
};

export default MappingWizard;
