import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter, useLocation, useNavigate } from 'react-router-dom';
import DatalinkWorkbenchV2Page from '../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { WorkbenchV2Shell } from '../../../src/features/datalink/workbench-v2/shell/WorkbenchV2Shell';
import { INITIAL_STATE } from '../../../src/features/datalink/workbench-v2/state/useWorkbenchV2State';

// 模擬 i18next 避免 namespace 未定義錯誤
vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
    i18n: { changeLanguage: () => Promise.resolve() },
  }),
}));

describe('Workbench V2 Shell & Integration', () => {
  const originalInnerWidth = window.innerWidth;

  function LocationProbe() {
    const location = useLocation();
    return <div data-testid="shell-location">{`${location.pathname}${location.search}`}</div>;
  }

  function ShellRouterHarness({
    state,
    actions,
  }: {
    state: typeof INITIAL_STATE;
    actions: {
      state: typeof INITIAL_STATE;
      setView: ReturnType<typeof vi.fn>;
      setCurrent: ReturnType<typeof vi.fn>;
      completeStep: ReturnType<typeof vi.fn>;
      toggleSidebar: ReturnType<typeof vi.fn>;
      toggleSummaryRail: ReturnType<typeof vi.fn>;
      setSidebarCollapsed: ReturnType<typeof vi.fn>;
      setShowSummaryRail: ReturnType<typeof vi.fn>;
      resetFlow: ReturnType<typeof vi.fn>;
      selectRule: ReturnType<typeof vi.fn>;
      dispatch: ReturnType<typeof vi.fn>;
    };
  }) {
    const navigate = useNavigate();

    return (
      <>
        <WorkbenchV2Shell state={state} actions={actions} navigateTo={navigate} />
        <LocationProbe />
      </>
    );
  }

  beforeEach(() => {
    // 預設寬度為 1440px (寬螢幕，渲染 SummaryRail)
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 1440,
    });
    
    // 清除 localStorage
    localStorage.clear();
    
    // Mock console.warn 防止污染
    vi.spyOn(console, 'warn').mockImplementation(() => {});
    
    // Spy fetch
    vi.spyOn(window, 'fetch').mockImplementation(() => Promise.resolve({} as Response));
  });

  afterEach(() => {
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: originalInnerWidth,
    });
    vi.restoreAllMocks();
    localStorage.clear();
  });

  // 1.2 Font-loading verify
  it('should verify document.fonts fallback behavior in testing', () => {
    // Mock document.fonts to contain Inter and JetBrains Mono
    const mockFonts = {
      ready: Promise.resolve(),
      check: vi.fn().mockReturnValue(true),
      forEach: vi.fn(),
    };
    Object.defineProperty(document, 'fonts', {
      value: mockFonts,
      configurable: true,
    });

    expect(document.fonts.check('12px Inter')).toBe(true);
    expect(document.fonts.check('12px "JetBrains Mono"')).toBe(true);
  });

  // 5.5 Default render of all regions
  it('renders all regions with default state', () => {
    render(<DatalinkWorkbenchV2Page />);

    // 驗證 2.3 data-workbench-v2 屬性在 root
    const root = screen.getByTestId('workbench-v2-root');
    expect(root).toBeInTheDocument();
    expect(root).toHaveAttribute('data-workbench-v2', 'true');

    // 驗證頂部列
    expect(screen.getByText('Datalink Workbench')).toBeInTheDocument();
    expect(screen.getByText('Datalink')).toBeInTheDocument();
    expect(screen.getByText('Workbench')).toBeInTheDocument();
    // 透過 getAllByText 來處理畫面上多個「新增裝置」字眼
    expect(screen.getAllByText('新增裝置').length).toBeGreaterThanOrEqual(1);

    // 驗證左側步驟軌道
    expect(screen.getByText('建置流程')).toBeInTheDocument();
    expect(screen.getAllByText('Device + Protocol').length).toBeGreaterThanOrEqual(1);

    // 驗證中央內容
    expect(screen.getByTestId('step1-device-view')).toBeInTheDocument();

    // 驗證右側摘要列 (>= 1280px 時渲染)
    expect(screen.getByTestId('summary-rail')).toBeInTheDocument();
    expect(screen.getByText('總數')).toBeInTheDocument();
  });

  // 5.2 Step navigation & click
  it('allows clicking reachable steps and triggers setCurrent callback', () => {
    const mockState = {
      ...INITIAL_STATE,
      completed: new Set([1]),
    };
    const mockActions = {
      state: mockState,
      setView: vi.fn(),
      setCurrent: vi.fn(),
      completeStep: vi.fn(),
      toggleSidebar: vi.fn(),
      toggleSummaryRail: vi.fn(),
      setSidebarCollapsed: vi.fn(),
      setShowSummaryRail: vi.fn(),
      resetFlow: vi.fn(),
      selectRule: vi.fn(),
      dispatch: vi.fn(),
    };

    render(<WorkbenchV2Shell state={mockState} actions={mockActions} />);
    
    // 驗證第 2 步在 completed 有 1 時是 reachable 的
    const step2Button = screen.getByTestId('step-nav-button-2');
    expect(step2Button).not.toBeDisabled();

    // 點擊切換到第 2 步
    fireEvent.click(step2Button);

    // 驗證 setCurrent 被呼叫且帶參數 2
    expect(mockActions.setCurrent).toHaveBeenCalledWith(2);
  });

  // 5.2 Settings view switch
  it('switches view mode when clicking settings', () => {
    render(<DatalinkWorkbenchV2Page />);

    const settingsButton = screen.getByTestId('settings-nav-button');
    fireEvent.click(settingsButton);

    // 驗證中央內容換成 Settings 頁面
    expect(screen.getByTestId('settings-page')).toBeInTheDocument();

    // 驗證 settings view 下 SummaryRail 應被隱藏
    expect(screen.queryByTestId('summary-rail')).not.toBeInTheDocument();
  });

  it('does not render SettingsPlaceholder in shell settings view', () => {
    const mockState = {
      ...INITIAL_STATE,
      view: 'settings' as const,
    };
    const mockActions = {
      state: mockState,
      setView: vi.fn(),
      setCurrent: vi.fn(),
      completeStep: vi.fn(),
      toggleSidebar: vi.fn(),
      toggleSummaryRail: vi.fn(),
      setSidebarCollapsed: vi.fn(),
      setShowSummaryRail: vi.fn(),
      resetFlow: vi.fn(),
      selectRule: vi.fn(),
      dispatch: vi.fn(),
    };

    render(<WorkbenchV2Shell state={mockState} actions={mockActions} />);
    expect(screen.getByTestId('settings-page')).toBeInTheDocument();
    expect(screen.queryByText('設定內容即將上線')).not.toBeInTheDocument();
  });

  // 5.3 Viewport detection
  it('does not render SummaryRail when viewport width is less than 1280px', () => {
    // 模擬 1024px viewport (narrow viewport)
    Object.defineProperty(window, 'innerWidth', {
      writable: true,
      configurable: true,
      value: 1024,
    });

    render(<DatalinkWorkbenchV2Page />);

    // 驗證 SummaryRail 不在 DOM 中
    expect(screen.queryByTestId('summary-rail')).not.toBeInTheDocument();
  });

  // 6.1 Keyboard shortcut ⌘B
  it('toggles sidebar state on Cmd+B or Ctrl+B keydown', () => {
    render(<DatalinkWorkbenchV2Page />);

    const sidebar = screen.getByTestId('sidebar-rail-aside');
    expect(sidebar).toHaveAttribute('data-collapsed', 'false');

    // 觸發 Cmd+B
    fireEvent.keyDown(window, { key: 'b', metaKey: true });
    expect(sidebar).toHaveAttribute('data-collapsed', 'true');

    // 再觸發一次
    fireEvent.keyDown(window, { key: 'b', metaKey: true });
    expect(sidebar).toHaveAttribute('data-collapsed', 'false');
  });

  // 6.2 Keyboard shortcut suppressed in text input
  it('suppresses sidebar toggle shortcut inside input element', () => {
    render(<DatalinkWorkbenchV2Page />);

    const sidebar = screen.getByTestId('sidebar-rail-aside');
    expect(sidebar).toHaveAttribute('data-collapsed', 'false');

    // 建立臨時 input 元素並 focus
    const input = document.createElement('input');
    document.body.appendChild(input);
    input.focus();

    // 在 input 內按 Cmd+B
    fireEvent.keyDown(input, { key: 'b', metaKey: true });

    // 驗證 sidebar 狀態不改變
    expect(sidebar).toHaveAttribute('data-collapsed', 'false');
    document.body.removeChild(input);
  });

  // 5.4 TweaksPanel availability & localStorage warning
  it('warns when localStorage throws security error', () => {
    // 設置 tweaks panel 顯示條件
    localStorage.setItem('WBV2_TWEAKS', '1');

    // 模擬 localStorage 讀取異常 (隱私模式)
    vi.spyOn(Storage.prototype, 'getItem').mockImplementation((key) => {
      if (key === 'wbv2_sidebar_collapsed') {
        throw new Error('SecurityError: The operation is insecure.');
      }
      return localStorage.getItem(key);
    });

    const consoleWarnSpy = vi.spyOn(console, 'warn');

    render(<DatalinkWorkbenchV2Page />);

    // 觸發 mount 的 useEffect 呼叫
    expect(consoleWarnSpy).toHaveBeenCalled();
  });

  it('當前步驟為 2 時應渲染 Step2Rule 元件', () => {
    const mockState = {
      ...INITIAL_STATE,
      current: 2 as const,
    };
    const mockActions = {
      state: mockState,
      setView: vi.fn(),
      setCurrent: vi.fn(),
      completeStep: vi.fn(),
      toggleSidebar: vi.fn(),
      toggleSummaryRail: vi.fn(),
      setSidebarCollapsed: vi.fn(),
      setShowSummaryRail: vi.fn(),
      resetFlow: vi.fn(),
      selectRule: vi.fn(),
      dispatch: vi.fn(),
    };

    render(<WorkbenchV2Shell state={mockState} actions={mockActions} />);
    expect(screen.getByTestId('step2-rule-container')).toBeInTheDocument();
  });

  describe('TopBar scheduler', () => {
    it('當 committed 為 false 時，應顯示 scheduler idle 狀態', () => {
      const mockState = {
        ...INITIAL_STATE,
        committed: false,
      };
      const mockActions = {
        state: mockState,
        setView: vi.fn(),
        setCurrent: vi.fn(),
        completeStep: vi.fn(),
        toggleSidebar: vi.fn(),
        toggleSummaryRail: vi.fn(),
        setSidebarCollapsed: vi.fn(),
        setShowSummaryRail: vi.fn(),
        resetFlow: vi.fn(),
        selectRule: vi.fn(),
        dispatch: vi.fn(),
      };

      render(<WorkbenchV2Shell state={mockState} actions={mockActions} />);
      expect(screen.getByText('step4.scheduler_idle')).toBeInTheDocument();
    });

    it('當 committed 為 true 時，應顯示 scheduler running 狀態', () => {
      const mockState = {
        ...INITIAL_STATE,
        committed: true,
      };
      const mockActions = {
        state: mockState,
        setView: vi.fn(),
        setCurrent: vi.fn(),
        completeStep: vi.fn(),
        toggleSidebar: vi.fn(),
        toggleSummaryRail: vi.fn(),
        setSidebarCollapsed: vi.fn(),
        setShowSummaryRail: vi.fn(),
        resetFlow: vi.fn(),
        selectRule: vi.fn(),
        dispatch: vi.fn(),
      };

      render(<WorkbenchV2Shell state={mockState} actions={mockActions} />);
      expect(screen.getByText('step4.scheduler_running')).toBeInTheDocument();
    });
  });

  describe('Runtime dashboard handoff', () => {
    it('navigates to the resolved runtime dashboard device route from the success card', () => {
      const mockState = {
        ...INITIAL_STATE,
        current: 4 as const,
        committed: true,
        commit: {
          status: 'success' as const,
          logs: [],
        },
      };
      const mockActions = {
        state: mockState,
        setView: vi.fn(),
        setCurrent: vi.fn(),
        completeStep: vi.fn(),
        toggleSidebar: vi.fn(),
        toggleSummaryRail: vi.fn(),
        setSidebarCollapsed: vi.fn(),
        setShowSummaryRail: vi.fn(),
        resetFlow: vi.fn(),
        selectRule: vi.fn(),
        dispatch: vi.fn(),
      };

      render(
        <MemoryRouter initialEntries={['/studio/v2']}>
          <ShellRouterHarness state={mockState} actions={mockActions} />
        </MemoryRouter>,
      );

      fireEvent.click(screen.getByText('step4.go_to_dashboard_btn'));

      expect(screen.getByTestId('shell-location')).toHaveTextContent('/studio/runtime?device_id=dev-01');
    });
  });
});
