import { useCallback, useEffect, useMemo, useRef } from "react";
import type { ProtocolType } from "../../../types/datalink";

type ModbusArea = "0" | "1" | "3" | "4";

const MODBUS_AREA_DEFAULTS: Record<ModbusArea, string> = {
  "0": "00001",
  "1": "10001",
  "3": "30001",
  "4": "40001",
};

export interface UseSmartDashboardWorkspaceContentStateInput {
  planStartAddress: string;
  setPlanStartAddress: (value: string) => void;
  setPlanCount: (value: number) => void;
  selectedDeviceProtocol: ProtocolType;
  gridViewStartAddress: string;
  getGridCenterAddress: (protocol: ProtocolType) => string;
  onGridViewShift: (delta: number) => void;
}

export function useSmartDashboardWorkspaceContentState({
  planStartAddress,
  setPlanStartAddress,
  setPlanCount,
  selectedDeviceProtocol,
  gridViewStartAddress,
  getGridCenterAddress,
  onGridViewShift,
}: UseSmartDashboardWorkspaceContentStateInput) {
  const lastWheelShiftAt = useRef(0);
  const gridScrollRef = useRef<HTMLDivElement>(null);
  const wheelShiftMs = 300;
  const gridRange = 100;

  const modbusArea = useMemo<ModbusArea>(() => {
    if (planStartAddress.startsWith("0")) return "0";
    if (planStartAddress.startsWith("1")) return "1";
    if (planStartAddress.startsWith("3")) return "3";
    return "4";
  }, [planStartAddress]);

  const centerAddress = useMemo(
    () => gridViewStartAddress || planStartAddress || getGridCenterAddress(selectedDeviceProtocol),
    [getGridCenterAddress, gridViewStartAddress, planStartAddress, selectedDeviceProtocol],
  );

  const handleModbusAreaChange = useCallback(
    (prefix: ModbusArea) => {
      setPlanStartAddress(MODBUS_AREA_DEFAULTS[prefix] ?? MODBUS_AREA_DEFAULTS["4"]);
    },
    [setPlanStartAddress],
  );

  const handlePlanStartAddressChange = useCallback(
    (address: string) => {
      setPlanStartAddress(address.toUpperCase());
    },
    [setPlanStartAddress],
  );

  const handlePlanCountChange = useCallback(
    (rawValue: string) => {
      const parsed = Number(rawValue);
      const bounded = Number.isFinite(parsed) ? Math.min(200, Math.max(1, Math.floor(parsed))) : 1;
      setPlanCount(bounded);
    },
    [setPlanCount],
  );

  useEffect(() => {
    const element = gridScrollRef.current;
    if (!element) return;

    const onWheel = (event: WheelEvent) => {
      if (event.deltaY === 0) return;
      const now = Date.now();
      if (now - lastWheelShiftAt.current < wheelShiftMs) {
        event.preventDefault();
        return;
      }
      lastWheelShiftAt.current = now;
      event.preventDefault();
      onGridViewShift(event.deltaY > 0 ? 100 : -100);
    };

    element.addEventListener("wheel", onWheel, { passive: false });
    return () => element.removeEventListener("wheel", onWheel);
  }, [onGridViewShift]);

  return {
    gridScrollRef,
    gridRange,
    modbusArea,
    centerAddress,
    handleModbusAreaChange,
    handlePlanStartAddressChange,
    handlePlanCountChange,
  };
}
