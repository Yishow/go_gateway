import { useCallback, useEffect, useState } from 'react';
import { modbusShareAPI } from '../../../services/datalink';
import type { ModbusShareStatus, Point, Tag } from '../../../types/datalink';
import { designSystem } from '../../../styles/designSystem';

interface UseSmartDashboardModbusActionsParams {
  linkedTag?: Tag;
  selectedPoint: Point | null;
  showError: (message: string) => void;
  showInfo: (message: string) => void;
  showSuccess: (message: string) => void;
}

export function useSmartDashboardModbusActions({
  linkedTag,
  selectedPoint,
  showError,
  showInfo,
  showSuccess,
}: UseSmartDashboardModbusActionsParams) {
  const [modbusStatus, setModbusStatus] = useState<ModbusShareStatus | null>(null);
  const [modbusRegister, setModbusRegister] = useState('0');

  const loadModbusStatus = useCallback(async () => {
    try {
      const status = await modbusShareAPI.status();
      setModbusStatus(status);
    } catch (error) {
      const message = error instanceof Error ? error.message : '讀取 Modbus 分享狀態失敗';
      showError(message);
    }
  }, [showError]);

  useEffect(() => {
    void loadModbusStatus();
  }, [loadModbusStatus]);

  const handleBindTagToModbus = useCallback(async () => {
    if (!linkedTag?.id) {
      showError('請先選取已連結 Tag 的點位');
      return;
    }
    const register = Number(modbusRegister);
    if (!Number.isInteger(register) || register < 0 || register > 65535) {
      showError('Register 必須為 0-65535 的整數');
      return;
    }
    try {
      await modbusShareAPI.upsertMapping(linkedTag.id, register);
      await loadModbusStatus();
      showSuccess(`已綁定 ${linkedTag.key} -> HR${register}`);
    } catch (error) {
      const message = error instanceof Error ? error.message : '綁定失敗';
      showError(message);
    }
  }, [linkedTag?.id, linkedTag?.key, loadModbusStatus, modbusRegister, showError, showSuccess]);

  const handlePushCurrentValueToModbus = useCallback(async () => {
    if (!linkedTag?.id) {
      showError('請先選取已連結 Tag 的點位');
      return;
    }
    if (selectedPoint?.last_value === undefined || selectedPoint?.last_value === null) {
      showError('目前點位沒有可推送的值');
      return;
    }
    try {
      await modbusShareAPI.writeTagValue(linkedTag.id, selectedPoint.last_value);
      showSuccess(`已推送當前值到 Tag ${linkedTag.key} 的 Modbus 映射`);
    } catch (error) {
      const message = error instanceof Error ? error.message : '推送失敗';
      showError(message);
    }
  }, [linkedTag?.id, linkedTag?.key, selectedPoint?.last_value, showError, showSuccess]);

  const handleSyncModbusFromMappings = useCallback(async () => {
    try {
      const result = await modbusShareAPI.sync();
      await loadModbusStatus();
      const errorHint = result.errors.length > 0 ? `, errors=${result.errors.length}` : '';
      showInfo(`${designSystem.microcopy.feedback.success.validated}：updated=${result.updated}, skipped=${result.skipped}${errorHint}`);
    } catch (error) {
      const message = error instanceof Error ? error.message : designSystem.microcopy.feedback.error.validationFailed;
      showError(message);
    }
  }, [loadModbusStatus, showError, showInfo]);

  return {
    modbusStatus,
    modbusRegister,
    setModbusRegister,
    loadModbusStatus,
    handleBindTagToModbus,
    handlePushCurrentValueToModbus,
    handleSyncModbusFromMappings,
  };
}
