import { useQuery } from '@tanstack/react-query';
import { modbusShareAPI } from '../../services/datalink';

export const modbusShareStatusKey = ['datalink', 'modbus-share', 'status'] as const;

/** Reads backend Share lifecycle/hydration truth without manufacturing defaults. */
export function useModbusShareStatusQuery(enabled = true) {
  return useQuery({
    queryKey: modbusShareStatusKey,
    queryFn: () => modbusShareAPI.status(),
    enabled,
    retry: false,
    refetchOnWindowFocus: false,
  });
}
