import { describe, expect, it, vi } from 'vitest';
import { fireEvent, render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import DatalinkWorkbenchV2Page from '../../../src/pages/datalink/workbench-v2/DatalinkWorkbenchV2Page';
import { Button, SectionCard } from '../../../src/features/datalink/workbench-v2/components';
import { StepRail } from '../../../src/features/datalink/workbench-v2/shell/StepRail';

const mockState = vi.hoisted(() => ({
  workspaceQuery: {
    data: undefined,
    isLoading: false,
    isSuccess: false,
    isError: false,
  },
}));

vi.mock('../../../src/hooks/datalink/useStudioV2Workspace', () => ({
  useStudioV2WorkspaceQuery: () => mockState.workspaceQuery,
}));

vi.mock('../../../src/hooks/datalink/useStudioV2WorkspaceAuditHistory', () => ({
  useStudioV2WorkspaceAuditHistoryQuery: () => ({
    data: { entries: [] },
    isError: false,
  }),
}));

vi.mock('../../../src/pages/datalink/workbench-v2/useStudioV2AutosaveState', () => ({
  useStudioV2AutosaveState: () => ({
    workspaceHydrated: false,
    devicesQuery: { isLoading: false, isError: false },
    rulesQuery: { isLoading: false, isError: false },
    mappingsQuery: { isLoading: false, isError: false },
    databaseConfigQuery: { isLoading: false, isError: false },
    databaseTargetsQuery: { isLoading: false, isError: false },
    state: {},
    actions: {},
  }),
}));

vi.mock('../../../src/hooks/datalink/useStudioV2WorkspaceActivation', () => ({
  useActivateStudioV2WorkspaceMutation: () => ({
    mutateAsync: vi.fn(),
  }),
}));

function renderWorkbenchPage() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
    },
  });

  return render(
    <QueryClientProvider client={queryClient}>
      <DatalinkWorkbenchV2Page />
    </QueryClientProvider>,
  );
}

describe('Workbench V2 redesign contract', () => {
  it('scopes the bootstrap loading surface to the workbench visual system', () => {
    mockState.workspaceQuery = {
      data: undefined,
      isLoading: true,
      isSuccess: false,
      isError: false,
    };

    renderWorkbenchPage();

    const loading = screen.getByTestId('workbench-v2-bootstrap-loading');
    expect(loading).toHaveAttribute('data-workbench-v2', 'true');
    expect(loading).toHaveAttribute('role', 'status');
    expect(loading).toHaveClass('wbv2-bootstrap-state');
    expect(screen.getByTestId('workbench-v2-bootstrap-shell')).toBeInTheDocument();
  });

  it('scopes the bootstrap error surface to the workbench visual system', () => {
    mockState.workspaceQuery = {
      data: undefined,
      isLoading: false,
      isSuccess: false,
      isError: true,
    };

    renderWorkbenchPage();

    const error = screen.getByTestId('workbench-v2-bootstrap-error');
    expect(error).toHaveAttribute('data-workbench-v2', 'true');
    expect(error).toHaveAttribute('role', 'alert');
    expect(error).toHaveClass('wbv2-bootstrap-state');
    expect(screen.getByTestId('workbench-v2-bootstrap-shell')).toBeInTheDocument();
  });

  it('adds pressed and visible focus feedback to shared buttons', () => {
    render(<Button>Commit</Button>);

    const button = screen.getByRole('button', { name: 'Commit' });
    expect(button).toHaveClass('active:scale-[0.98]');
    expect(button).toHaveClass('focus-visible:ring-offset-2');
  });

  it('uses the refined workbench surface class for section cards', () => {
    render(<SectionCard title="Connector">Content</SectionCard>);

    const section = screen.getByRole('region', { name: 'Connector' });
    expect(section).toHaveClass('wbv2-surface');
  });

  it('adds visible keyboard focus treatment to step rail buttons', () => {
    render(
      <StepRail
        view="flow"
        current={1}
        completed={new Set()}
        step2Ready={true}
        collapsed={false}
        onJump={vi.fn()}
        onSwitchView={vi.fn()}
      />,
    );

    expect(screen.getByTestId('step-nav-button-1')).toHaveClass('focus-visible:ring-2');
    expect(screen.getByTestId('settings-nav-button')).toHaveClass('focus-visible:ring-2');
  });

  it('blocks stale completed Step 3 and Step 4 when Step 2 is not ready', () => {
    const onJump = vi.fn();

    render(
      <StepRail
        view="flow"
        current={1}
        completed={new Set([1, 2, 3, 4])}
        step2Ready={false}
        collapsed={false}
        onJump={onJump}
        onSwitchView={vi.fn()}
      />,
    );

    expect(screen.getByTestId('step-nav-button-2')).toBeEnabled();
    expect(screen.getByTestId('step-nav-button-3')).toBeDisabled();
    expect(screen.getByTestId('step-nav-button-4')).toBeDisabled();
    fireEvent.click(screen.getByTestId('step-nav-button-3'));
    expect(onJump).not.toHaveBeenCalled();
  });
});
