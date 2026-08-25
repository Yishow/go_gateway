import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { ShareOutputSummary } from '../../../src/features/datalink/workbench-v2/steps/step4/ShareOutputSummary';
import { Step4Database } from '../../../src/features/datalink/workbench-v2/steps/step4/Step4Database';
import { modbusShareStatusFromBootstrap } from '../../../src/types/modbusShare';
import { hydratedShareStatus, mcShareState, shareContext } from './step4-share-helpers';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (_key: string, options?: { defaultValue?: string }) => options?.defaultValue ?? _key,
  }),
}));

describe('Step 4 Share output summary', () => {
  it('does not manufacture a disabled status when bootstrap omits backend status', () => {
    expect(modbusShareStatusFromBootstrap({
      hydration_state: 'pending',
      readiness: false,
    })).toBeNull();
  });

  it('hydrates Share truth from the workspace bootstrap response shape', () => {
    const status = modbusShareStatusFromBootstrap({
      hydration_state: 'ready',
      workspace_id: 'workspace-bootstrap',
      workspace_revision: 'workspace-revision-9',
      settings_revision: 'settings-revision-9',
      readiness: true,
      readiness_token: 'readiness-token-9',
      error: {
        code: 'modbus_share_hydration_required',
        message: 'backend detail must stay out of the UI',
        retryable: true,
        request_id: 'req-bootstrap-9',
        action: 'Retry workspace bootstrap',
      },
      status: { ...hydratedShareStatus, enabled: false, bind_state: 'disabled' },
      canonical_plan: shareContext.canonical_plan,
    });

    expect(status).toMatchObject({
      workspace_id: 'workspace-bootstrap',
      workspace_revision: 'workspace-revision-9',
      settings_revision: 'settings-revision-9',
      readiness_token: 'readiness-token-9',
      hydration_state: 'ready',
      readiness: true,
      enabled: false,
      bind_state: 'disabled',
      error: {
        code: 'modbus_share_hydration_required',
        request_id: 'req-bootstrap-9',
      },
      canonical_plan: { signature: shareContext.canonical_plan?.signature },
    });
  });

  it('retains the non-Modbus Step4 destination surface while Share is configured', () => {
    render(<Step4Database state={mcShareState()} dispatch={vi.fn()} />);

    expect(screen.getByTestId('step4-destination-overview')).toBeInTheDocument();
  });

  it('shows non-Modbus MC D0~D3 as 40001 ~ 40004 without DB targets', () => {
    render(<ShareOutputSummary state={mcShareState()} shareStatus={hydratedShareStatus} />);

    expect(screen.getByTestId('step4-share-output-summary')).toHaveTextContent('D0 ~ D3');
    expect(screen.getByTestId('step4-share-output-summary')).toHaveTextContent('40001 ~ 40004');
  });

  it('uses canonical backend geometry even when browser rule layout is different', () => {
    const state = mcShareState({
      rules: mcShareState().rules.map((rule) => ({ ...rule, share_start_register: 41000, share_stride: 99 })),
    });
    render(<ShareOutputSummary state={state} shareStatus={hydratedShareStatus} />);

    expect(screen.getByTestId('step4-share-output-summary')).toHaveTextContent('40001 ~ 40004');
    expect(screen.getByTestId('step4-share-output-summary')).not.toHaveTextContent('41000');
  });

  it('does not invent rows for an empty canonical candidate set', () => {
    render(<ShareOutputSummary state={mcShareState()} shareStatus={{ ...hydratedShareStatus, canonical_desired_mappings: [] }} />);

    expect(screen.queryByTestId('step4-share-output-summary')).not.toBeInTheDocument();
    expect(screen.queryByTestId('step4-share-output-row-mc-rule')).not.toBeInTheDocument();
  });

  it('does not show rows when backend candidate state is blocked by a collision', () => {
    render(<ShareOutputSummary state={mcShareState()} shareStatus={{
      ...hydratedShareStatus,
      candidate_set_status: 'blocked',
      candidate_statuses: ['blocked_conflict'],
      error: {
        code: 'modbus_share_range_collision',
        message: 'hidden',
        retryable: false,
        request_id: 'request-collision',
      },
    }} />);

    expect(screen.getByTestId('step4-share-output-failed')).toBeInTheDocument();
    expect(screen.queryByTestId('step4-share-output-summary')).not.toBeInTheDocument();
  });

  it('blocks the summary when backend Share status is unavailable instead of using defaults', () => {
    render(<ShareOutputSummary state={mcShareState()} />);

    expect(screen.getByTestId('step4-share-output-blocked')).toBeInTheDocument();
    expect(screen.queryByTestId('step4-share-output-summary')).not.toBeInTheDocument();
    expect(screen.getByTestId('step4-share-output-blocked')).not.toHaveTextContent('40001 ~ 40004');
  });

  it('blocks before hydration readiness and shows disabled backend truth', () => {
    const { rerender } = render(
      <ShareOutputSummary
        state={mcShareState()}
        shareStatus={{ ...hydratedShareStatus, hydration_state: 'pending', readiness: false }}
      />,
    );
    expect(screen.getByTestId('step4-share-output-blocked')).toBeInTheDocument();

    rerender(
      <ShareOutputSummary
        state={mcShareState()}
        shareStatus={{ ...hydratedShareStatus, enabled: false, bind_state: 'disabled' }}
      />,
    );
    expect(screen.getByTestId('step4-share-output-disabled')).toBeInTheDocument();
    expect(screen.queryByTestId('step4-share-output-summary')).not.toBeInTheDocument();
  });

  it('prioritizes persisted disabled truth even when hydration and canonical data are absent', () => {
    render(
      <ShareOutputSummary
        state={mcShareState()}
        shareStatus={{
          ...hydratedShareStatus,
          configured_enabled: false,
          hydration_state: 'pending',
          readiness: false,
          canonical_desired_mappings: undefined,
          canonical_plan: undefined,
          bind_state: 'disabled',
        }}
      />,
    );

    expect(screen.getByTestId('step4-share-output-disabled')).toBeInTheDocument();
    expect(screen.queryByTestId('step4-share-output-blocked')).not.toBeInTheDocument();
    expect(screen.queryByTestId('step4-share-output-summary')).not.toBeInTheDocument();
  });

  it('keeps hydration truth ahead of pending and failed states when Share is enabled', () => {
    const { rerender } = render(
      <ShareOutputSummary
        state={mcShareState()}
        shareStatus={{
          ...hydratedShareStatus,
          configured_enabled: true,
          hydration_state: 'pending',
          readiness: false,
          running: false,
          bind_state: 'stopped',
          canonical_desired_mappings: undefined,
          canonical_plan: undefined,
        }}
      />,
    );

    expect(screen.getByTestId('step4-share-output-blocked')).toBeInTheDocument();
    expect(screen.queryByTestId('step4-share-output-pending')).not.toBeInTheDocument();

    rerender(
      <ShareOutputSummary
        state={mcShareState()}
        shareStatus={{
          ...hydratedShareStatus,
          configured_enabled: true,
          hydration_state: 'ready',
          readiness: true,
          running: true,
          bind_state: 'fail',
          failed: true,
          canonical_desired_mappings: undefined,
          canonical_plan: undefined,
        }}
      />,
    );

    expect(screen.getByTestId('step4-share-output-blocked')).toBeInTheDocument();
    expect(screen.queryByTestId('step4-share-output-failed')).not.toBeInTheDocument();
  });

  it('shows failed binding after hydrated canonical truth is ready', () => {
    render(
      <ShareOutputSummary
        state={mcShareState()}
        shareStatus={{
          ...hydratedShareStatus,
          configured_enabled: true,
          running: false,
          bind_state: 'fail',
          failed: true,
        }}
      />,
    );

    expect(screen.getByTestId('step4-share-output-failed')).toBeInTheDocument();
    expect(screen.queryByTestId('step4-share-output-pending')).not.toBeInTheDocument();
  });

  it('shows pending stopped binding after hydrated canonical truth is ready', () => {
    render(
      <ShareOutputSummary
        state={mcShareState()}
        shareStatus={{
          ...hydratedShareStatus,
          configured_enabled: true,
          running: false,
          bind_state: 'stopped',
          failed: false,
        }}
      />,
    );

    expect(screen.getByTestId('step4-share-output-pending')).toBeInTheDocument();
    expect(screen.queryByTestId('step4-share-output-failed')).not.toBeInTheDocument();
  });

  it('shows backend failed status instead of presenting a healthy Share layout', () => {
    render(
      <ShareOutputSummary
        state={mcShareState()}
        shareStatus={{
          ...hydratedShareStatus,
          enabled: false,
          bind_state: 'fail',
          lifecycle_state: 'failed',
          failed: true,
        }}
      />,
    );

    expect(screen.getByTestId('step4-share-output-failed')).toBeInTheDocument();
    expect(screen.queryByTestId('step4-share-output-summary')).not.toBeInTheDocument();
  });
});
