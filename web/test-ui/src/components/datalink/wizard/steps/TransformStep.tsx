import React from 'react';
import TransformBuilder from '../../TransformBuilder';
import type { TransformStep as TransformStepType } from '../../../../types/datalink';
import './Steps.css';

/**
 * TransformStep 組件屬性
 */
interface TransformStepProps {
  /** 當前管線 */
  pipeline: TransformStepType[];
  /** 變更回呼 */
  onChange: (pipeline: TransformStepType[]) => void;
}

/**
 * 轉換設定步驟
 */
export const TransformStep: React.FC<TransformStepProps> = ({
  pipeline,
  onChange,
}) => {
  return (
    <div className="step-transform">
      <div className="step-info-box">
        <p>
          配置轉換管線，將原始值轉換為最終值。
          如果不需要轉換，可以直接跳過此步驟。
        </p>
      </div>

      <TransformBuilder
        steps={pipeline}
        onChange={onChange}
      />
    </div>
  );
};

export default TransformStep;
