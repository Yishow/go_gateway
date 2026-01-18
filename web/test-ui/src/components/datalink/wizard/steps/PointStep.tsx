import React, { useState, useEffect } from 'react';
import { pointAPI } from '../../../../services/datalink';
import type { Point } from '../../../../types/datalink';
import './Steps.css';

/**
 * PointStep 組件屬性
 */
interface PointStepProps {
  /** 設備 ID */
  deviceId: string;
  /** 已選擇的點位 ID */
  selectedPointId: string;
  /** 選擇回呼 */
  onSelect: (pointId: string) => void;
  /** 錯誤訊息 */
  error?: string;
}

/**
 * 點位配置步驟
 */
export const PointStep: React.FC<PointStepProps> = ({
  deviceId,
  selectedPointId,
  onSelect,
  error,
}) => {
  const [points, setPoints] = useState<Point[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (deviceId) {
      loadPoints();
    }
  }, [deviceId]);

  const loadPoints = async () => {
    try {
      setLoading(true);
      const data = await pointAPI.list({ device_id: deviceId });
      setPoints(data);
    } catch (err) {
      console.error('Failed to load points:', err);
    } finally {
      setLoading(false);
    }
  };

  if (!deviceId) {
    return <div className="step-warning">請先選擇設備</div>;
  }

  if (loading) {
    return <div className="step-loading">載入中...</div>;
  }

  return (
    <div className="step-point">
      {error && <div className="step-error">{error}</div>}

      <div className="point-list">
        {points.length === 0 ? (
          <div className="no-data">
            此設備沒有點位，請新增點位
          </div>
        ) : (
          points.map((point) => (
            <div
              key={point.id}
              className={`point-card ${selectedPointId === point.id ? 'selected' : ''}`}
              onClick={() => onSelect(point.id)}
            >
              <div className="point-icon">📍</div>
              <div className="point-info">
                <h4>{point.name}</h4>
                <p>地址: {point.address} | 類型: {point.data_type}</p>
              </div>
              {selectedPointId === point.id && (
                <div className="selected-check">✓</div>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default PointStep;
