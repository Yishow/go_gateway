import { act, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import TestPage from '../../../src/pages/TestPage';
import { PROFILE_STORAGE_KEY, type Profile } from '../../../src/types/profile';

const { mockUseCardMinimizeState, mockConnect, mockDisconnect, mockRead, mockWrite } = vi.hoisted(() => ({
  mockUseCardMinimizeState: {
    isMinimized: vi.fn<(cardType: string) => boolean>(() => false),
    minimizeCardsBatch: vi.fn(),
  },
  mockConnect: vi.fn(),
  mockDisconnect: vi.fn(),
  mockRead: vi.fn(),
  mockWrite: vi.fn(),
}));

vi.mock('../../../src/hooks/useCardMinimize', () => ({
  useCardMinimize: () => mockUseCardMinimizeState,
}));

vi.mock('../../../src/services/api', () => ({
  useTestAPI: () => ({
    connect: mockConnect,
    disconnect: mockDisconnect,
    getStatus: vi.fn(),
    read: mockRead,
    write: mockWrite,
    batch: vi.fn(),
    startMonitor: vi.fn(),
    stopMonitor: vi.fn(),
  }),
}));

vi.mock('../../../src/contexts/ToastContext', () => ({
  useToast: () => ({
    showError: vi.fn(),
    showWarning: vi.fn(),
    showSuccess: vi.fn(),
  }),
}));

vi.mock('../../../src/components/ProtocolSelector', () => ({
  default: () => <div data-testid="protocol-selector-mock">protocol-selector</div>,
}));

vi.mock('../../../src/components/ProfileSelector', () => ({
  default: () => <div data-testid="profile-selector-mock">profile-selector</div>,
}));

vi.mock('../../../src/components/TestOperations', () => ({
  default: () => <div data-testid="test-operations-mock">test-operations</div>,
}));

vi.mock('../../../src/components/DebugPanel', () => ({
  default: () => <div data-testid="debug-panel-mock">debug-panel</div>,
}));

vi.mock('../../../src/components/MonitorControl', () => ({
  default: () => <div data-testid="monitor-control-mock">monitor-control</div>,
}));

vi.mock('../../../src/components/MinimizedCardsBar', () => ({
  default: () => <div data-testid="minimized-cards-bar-mock">minimized-cards-bar</div>,
}));

vi.mock('../../../src/components/CardMinimizeButton', () => ({
  default: () => null,
}));

vi.mock('../../../src/components/DeviceScanner', () => ({
  default: () => <div data-testid="device-scanner-mock">device-scanner</div>,
}));

function seedSerialProfile(): void {
  const serialProfile: Profile = {
    id: 'profile-serial',
    name: 'RTU Serial Profile',
    protocol: 'modbus_rtu',
    connectionMode: 'serial',
    config: {
      serial: {
        port: 'COM1',
        baudRate: 9600,
        dataBits: 8,
        stopBits: 1,
        parity: 'N',
      },
    },
    createdAt: 1,
    updatedAt: 1,
  };

  localStorage.setItem(PROFILE_STORAGE_KEY, JSON.stringify([serialProfile]));
  localStorage.setItem(`${PROFILE_STORAGE_KEY}_current`, serialProfile.id);
}

describe('TestPage RTU disconnect behavior', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    vi.clearAllMocks();
    localStorage.clear();
    seedSerialProfile();
    mockUseCardMinimizeState.isMinimized.mockImplementation(() => false);
    mockConnect.mockResolvedValue({ connection_id: 'conn-serial-1', status: 'connected' });
    mockDisconnect.mockImplementation(() => new Promise<void>(() => {}));
    mockRead.mockResolvedValue({ values: [true], count: 1 });
  });

  afterEach(() => {
    vi.useRealTimers();
    localStorage.clear();
  });

  it('stops polling immediately when disconnect is requested', async () => {
    render(<TestPage />);

    await act(async () => {
      await Promise.resolve();
    });

    await act(async () => {
      vi.advanceTimersByTime(100);
    });

    await act(async () => {
      fireEvent.click(screen.getByRole('button', { name: '建立連線' }));
      await Promise.resolve();
    });

    await act(async () => {
      fireEvent.click(screen.getByRole('button', { name: '開始' }));
      await Promise.resolve();
    });

    expect(mockRead).toHaveBeenCalledTimes(1);

    await act(async () => {
      fireEvent.click(screen.getByRole('button', { name: '斷開連線' }));
      await Promise.resolve();
    });

    expect(mockDisconnect).toHaveBeenCalledWith('conn-serial-1');

    await act(async () => {
      vi.advanceTimersByTime(1200);
    });

    expect(mockRead).toHaveBeenCalledTimes(1);
  });
});
