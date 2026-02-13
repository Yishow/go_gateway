import React, { useState, useEffect } from 'react';
import { deviceAPI } from '../../../../services/datalink';
import type { Device } from '../../../../types/datalink';
import './Steps.css';

/**
 * DeviceStep 組件屬性
 */
interface DeviceStepProps {
  /** 已選擇的設備 ID */
  selectedDeviceId: string;
  /** 選擇回呼 */
  onSelect: (deviceId: string) => void;
  /** 錯誤訊息 */
  error?: string;
}

/**
 * 設備選擇步驟
 */
export const DeviceStep: React.FC<DeviceStepProps> = ({
  selectedDeviceId,
  onSelect,
  error,
}) => {
  const [devices, setDevices] = useState<Device[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    loadDevices();
  }, []);

  const loadDevices = async () => {
    try {
      setLoading(true);
      const data = await deviceAPI.list();
      setDevices(data);
    } catch (err) {
      console.error('Failed to load devices:', err);
    } finally {
      setLoading(false);
    }
  };

  const filteredDevices = devices.filter((device) =>
    device.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  if (loading) {
    return <div className="step-loading">載入中...</div>;
  }

  return (
    <div className="step-device">
      <div className="step-search">
        <input
          type="text"
          placeholder="搜尋設備..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="search-input"
        />
      </div>

      {error && <div className="step-error">{error}</div>}

      <div className="device-list">
        {filteredDevices.length === 0 ? (
          <div className="no-data">
            {searchQuery ? '沒有符合的設備' : '沒有設備，請先新增設備'}
          </div>
        ) : (
          filteredDevices.map((device) => (
            <button
              type="button"
              key={device.id}
              className={`device-card ${selectedDeviceId === device.id ? 'selected' : ''}`}
              onClick={() => onSelect(device.id)}
              aria-pressed={selectedDeviceId === device.id}
              aria-label={`選擇設備 ${device.name}`}
            >
              <div className="device-icon">📱</div>
              <div className="device-info">
                <h4>{device.name}</h4>
                <p>{device.protocol}</p>
              </div>
              <div className={`device-status status-${device.status}`}>
                {device.status}
              </div>
              {selectedDeviceId === device.id && (
                <div className="selected-check">✓</div>
              )}
            </button>
          ))
        )}
      </div>
    </div>
  );
};

export default DeviceStep;
