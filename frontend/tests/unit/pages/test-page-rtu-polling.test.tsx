import { act, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import TestPage from '../../../src/pages/TestPage';
import { PROFILE_STORAGE_KEY, type Profile } from '../../../src/types/profile';

const { mockUseCardMinimizeState } = vi.hoisted(() => ({
  mockUseCardMinimizeState: {
    isMinimized: vi.fn<(cardType: string) => boolean>(() => false),
    minimizeCardsBatch: vi.fn(),
  },
}));

vi.mock('../../../src/hooks/useCardMinimize', () => ({
  useCardMinimize: () => mockUseCardMinimizeState,
}));

vi.mock('../../../src/components/ProtocolSelector', () => ({
  default: ({
    selectedProtocol,
    connectionMode,
  }: {
    selectedProtocol: string;
    connectionMode: string;
  }) => (
    <div data-testid="protocol-selector-mock">
      <p>{`protocol:${selectedProtocol}`}</p>
      <p>{`mode:${connectionMode}`}</p>
    </div>
  ),
}));

vi.mock('../../../src/components/ProfileSelector', () => ({
  default: () => <div data-testid="profile-selector-mock">profile-selector</div>,
}));

vi.mock('../../../src/components/ConfigForm', () => ({
  default: ({
    config,
    mode,
    protocol,
    connectionId,
    onConfigChange,
    onConnectionChange,
  }: {
    config: Record<string, string | number | string[] | undefined>;
    mode: string;
    protocol: string;
    connectionId: string | null;
    onConfigChange: (config: Record<string, string | number | string[] | undefined>) => void;
    onConnectionChange: (connectionId: string | null) => void;
  }) => (
    <div data-testid="config-form-mock">
      <p>{`config-protocol:${protocol}`}</p>
      <p>{`config-mode:${mode}`}</p>
      <p>{`config-connection:${connectionId ?? 'null'}`}</p>
      <button type="button" onClick={() => onConnectionChange('conn-serial-1')}>
        connect-success
      </button>
      <button
        type="button"
        onClick={() =>
          onConfigChange({
            ...config,
            baudRate: 19200,
          })
        }
      >
        change-config
      </button>
    </div>
  ),
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

vi.mock('../../../src/components/RTUPollingCard', () => ({
  default: ({ connectionId }: { connectionId: string | null }) => (
    <div data-testid="rtu-polling-card-mock">{`rtu-connection:${connectionId ?? 'null'}`}</div>
  ),
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

describe('TestPage RTU polling connection persistence', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    vi.clearAllMocks();
    localStorage.clear();
    seedSerialProfile();
    mockUseCardMinimizeState.isMinimized.mockImplementation(() => false);
  });

  afterEach(() => {
    vi.useRealTimers();
    localStorage.clear();
  });

  it('preserves the RTU polling connection after profile auto-save refreshes currentProfile', async () => {
    render(<TestPage />);

    await act(async () => {
      await Promise.resolve();
    });

    expect(screen.getByText('config-protocol:modbus_rtu')).toBeInTheDocument();
    expect(screen.getByText('config-mode:serial')).toBeInTheDocument();

    await act(async () => {
      vi.advanceTimersByTime(100);
    });

    fireEvent.click(screen.getByRole('button', { name: 'connect-success' }));

    expect(screen.getByTestId('rtu-polling-card-mock')).toHaveTextContent('rtu-connection:conn-serial-1');

    fireEvent.click(screen.getByRole('button', { name: 'change-config' }));

    await act(async () => {
      vi.advanceTimersByTime(1000);
    });

    await act(async () => {
      await Promise.resolve();
    });

    expect(screen.getByTestId('rtu-polling-card-mock')).toHaveTextContent('rtu-connection:conn-serial-1');
  });
});
