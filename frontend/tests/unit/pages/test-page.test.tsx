import { act, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import type { Profile } from '../../../src/types/profile';
import TestPage from '../../../src/pages/TestPage';

const {
  mockUseProfilesState,
  mockUseCardMinimizeState,
} = vi.hoisted(() => ({
  mockUseProfilesState: {
    currentProfile: null as Profile | null,
    updateCurrentProfileConfig: vi.fn(),
    updateProfile: vi.fn(),
    profiles: [] as Profile[],
    currentProfileId: null as string | null,
  },
  mockUseCardMinimizeState: {
    isMinimized: vi.fn<(cardType: string) => boolean>(() => false),
    minimizeCardsBatch: vi.fn(),
  },
}));

vi.mock('../../../src/hooks/useProfiles', () => ({
  useProfiles: () => mockUseProfilesState,
}));

vi.mock('../../../src/hooks/useCardMinimize', () => ({
  useCardMinimize: () => mockUseCardMinimizeState,
}));

vi.mock('../../../src/components/ProtocolSelector', () => ({
  default: ({
    selectedProtocol,
    connectionMode,
    onProtocolChange,
    onConnectionModeChange,
  }: {
    selectedProtocol: string;
    connectionMode: string;
    onProtocolChange: (protocol: string) => void;
    onConnectionModeChange: (mode: string) => void;
  }) => (
    <div data-testid="protocol-selector-mock">
      <p>{`protocol:${selectedProtocol}`}</p>
      <p>{`mode:${connectionMode}`}</p>
      <button
        type="button"
        onClick={() => {
          onProtocolChange('modbus_rtu');
          onConnectionModeChange('serial');
        }}
      >
        switch-to-serial
      </button>
    </div>
  ),
}));

vi.mock('../../../src/components/ProfileSelector', () => ({
  default: ({
    onProfileChange,
  }: {
    onProfileChange: (profile: Profile) => void;
  }) => (
    <div data-testid="profile-selector-mock">
      <button
        type="button"
        onClick={() =>
          onProfileChange({
            id: 'profile-2',
            name: '備援 Profile',
            protocol: 'modbus_tcp',
            connectionMode: 'tcp',
            config: {
              tcp: {
                host: '10.0.0.20',
                port: 1502,
              },
            },
            createdAt: 2,
            updatedAt: 2,
          })
        }
      >
        switch-profile
      </button>
    </div>
  ),
}));

vi.mock('../../../src/components/ConfigForm', () => ({
  default: ({
    mode,
    protocol,
    onConnectionChange,
    onMinimize,
  }: {
    mode: string;
    protocol: string;
    onConnectionChange: (connectionId: string | null) => void;
    onMinimize: () => void;
  }) => (
    <div data-testid="config-form-mock">
      <p>{`config-protocol:${protocol}`}</p>
      <p>{`config-mode:${mode}`}</p>
      <button type="button" onClick={() => onConnectionChange('conn-12345678')}>
        connect-success
      </button>
      <button type="button" onClick={onMinimize}>
        minimize-config
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
  default: () => <div data-testid="rtu-polling-card-mock">rtu-polling-card</div>,
}));

/**
 * 建立頁面測試用的 Profile 狀態，確保 `TestPage` 初始化時能載入既有 TCP 設定。
 */
function createProfileState(): void {
  mockUseProfilesState.currentProfile = {
    id: 'profile-1',
    name: '預設 Profile',
    protocol: 'modbus_tcp',
    connectionMode: 'tcp',
    config: {
      tcp: {
        host: '192.168.0.10',
        port: 502,
      },
    },
    createdAt: 1,
    updatedAt: 1,
  };
  mockUseProfilesState.profiles = [mockUseProfilesState.currentProfile];
  mockUseProfilesState.currentProfileId = 'profile-1';
}

describe('TestPage', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    vi.clearAllMocks();
    createProfileState();
    mockUseCardMinimizeState.isMinimized.mockImplementation(() => false);
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  /**
   * 統一渲染 `TestPage`，讓各測試專注在頁面行為，不重複建立樣板。
   */
  function renderPage() {
    return render(<TestPage />);
  }

  it('initializes batch card minimization for non-config panels', () => {
    renderPage();

    expect(mockUseCardMinimizeState.minimizeCardsBatch).toHaveBeenCalledTimes(1);
    expect(mockUseCardMinimizeState.minimizeCardsBatch).toHaveBeenCalledWith(
      expect.arrayContaining([
        expect.objectContaining({ type: 'operations', title: '測試操作' }),
        expect.objectContaining({ type: 'scanner', title: '設備掃描' }),
        expect.objectContaining({ type: 'monitor', title: '即時監控' }),
        expect.objectContaining({ type: 'debug', title: '調試面板' }),
      ])
    );
  });

  it('collapses the config card after a successful connection and allows expanding it again', async () => {
    renderPage();

    expect(screen.getByTestId('config-form-mock')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'connect-success' }));

    await act(async () => {
      vi.advanceTimersByTime(600);
    });

    expect(screen.queryByTestId('config-form-mock')).not.toBeInTheDocument();
    expect(screen.getByText('點擊展開')).toBeInTheDocument();
    expect(screen.getByText('192.168.0.10')).toBeInTheDocument();
    expect(screen.getByText('502')).toBeInTheDocument();
    expect(screen.getByText('conn-123...')).toBeInTheDocument();

    fireEvent.click(screen.getByText('點擊展開'));

    expect(screen.getByTestId('config-form-mock')).toBeInTheDocument();
  });

  it('switches protocol family to serial mode while preserving mode synchronization', async () => {
    renderPage();

    expect(screen.getByText('config-mode:tcp')).toBeInTheDocument();

    await act(async () => {
      vi.advanceTimersByTime(150);
    });

    fireEvent.click(screen.getByRole('button', { name: 'switch-to-serial' }));

    await act(async () => {
      vi.advanceTimersByTime(100);
    });

    expect(screen.getByText('config-mode:serial')).toBeInTheDocument();
    expect(mockUseProfilesState.updateProfile).toHaveBeenCalledWith(
      'profile-1',
      expect.objectContaining({
        protocol: 'modbus_tcp',
        connectionMode: 'tcp',
        config: expect.objectContaining({
          tcp: {
            host: '192.168.0.10',
            port: 502,
          },
        }),
      })
    );
  });
});
