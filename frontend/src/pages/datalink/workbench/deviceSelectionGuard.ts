type DeviceIdentity = {
  id: string;
};

type ShouldClearSelectedDeviceInput = {
  selectedDeviceId: string | null;
  devices: ReadonlyArray<DeviceIdentity>;
  isSuccess: boolean;
  isFetching: boolean;
  pendingSelectedDeviceId: string | null;
};

export function shouldClearSelectedDevice({
  selectedDeviceId,
  devices,
  isSuccess,
  isFetching,
  pendingSelectedDeviceId,
}: ShouldClearSelectedDeviceInput): boolean {
  if (!selectedDeviceId) {
    return false;
  }

  if (pendingSelectedDeviceId === selectedDeviceId) {
    return false;
  }

  if (!isSuccess || isFetching) {
    return false;
  }

  return !devices.some((device) => device.id === selectedDeviceId);
}
