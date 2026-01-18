import React from 'react';
import type { WizardStep } from './MappingWizard';
import './StepIndicator.css';

/**
 * StepIndicator 組件屬性
 */
interface StepIndicatorProps {
  /** 步驟列表 */
  steps: WizardStep[];
  /** 當前步驟索引 */
  currentStep: number;
  /** 點擊步驟回呼 */
  onStepClick?: (stepIndex: number) => void;
}

/**
 * 步驟指示器
 *
 * 頂部水平顯示所有步驟，點擊可跳回已完成的步驟
 */
export const StepIndicator: React.FC<StepIndicatorProps> = ({
  steps,
  currentStep,
  onStepClick,
}) => {
  return (
    <div className="step-indicator">
      {steps.map((step, index) => {
        const isCompleted = index < currentStep;
        const isCurrent = index === currentStep;
        const isClickable = index <= currentStep;

        return (
          <React.Fragment key={step.id}>
            <div
              className={`step-item ${isCompleted ? 'completed' : ''} ${isCurrent ? 'current' : ''} ${isClickable ? 'clickable' : ''}`}
              onClick={() => isClickable && onStepClick?.(index)}
              role={isClickable ? 'button' : undefined}
              tabIndex={isClickable ? 0 : undefined}
            >
              <div className="step-circle">
                {isCompleted ? (
                  <span className="check-mark">✓</span>
                ) : (
                  <span className="step-number">{index + 1}</span>
                )}
              </div>
              <div className="step-label">
                <span className="step-title">{step.title}</span>
              </div>
            </div>
            
            {index < steps.length - 1 && (
              <div className={`step-connector ${isCompleted ? 'completed' : ''}`} />
            )}
          </React.Fragment>
        );
      })}
    </div>
  );
};

export default StepIndicator;
