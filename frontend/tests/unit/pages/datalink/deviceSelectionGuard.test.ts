import { describe, expect, it } from 'vitest';
import { shouldClearSelectedDevice } from '../../../../src/pages/datalink/workbench/deviceSelectionGuard';

describe('shouldClearSelectedDevice', () => {
  it('keeps the current selection during a transient refetch gap', () => {
    expect(
      shouldClearSelectedDevice({
        selectedDeviceId: 'device-1',
        devices: [],
        isSuccess: false,
        isFetching: true,
        pendingSelectedDeviceId: null,
      }),
    ).toBe(false);
  });

  it('keeps the current selection while a newly created device is still pending', () => {
    expect(
      shouldClearSelectedDevice({
        selectedDeviceId: 'device-1',
        devices: [],
        isSuccess: true,
        isFetching: false,
        pendingSelectedDeviceId: 'device-1',
      }),
    ).toBe(false);
  });

  it('clears the selection after a successful load confirms the device is gone', () => {
    expect(
      shouldClearSelectedDevice({
        selectedDeviceId: 'device-1',
        devices: [{ id: 'device-2' }],
        isSuccess: true,
        isFetching: false,
        pendingSelectedDeviceId: null,
      }),
    ).toBe(true);
  });
});
