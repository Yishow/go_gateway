import { render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { ShareOutputSummary } from '../../../src/features/datalink/workbench-v2/steps/step4/ShareOutputSummary';
import type { ModbusShareStatus } from '../../../src/types/modbusShare';
import { hydratedShareStatus, mcShareState } from './step4-share-helpers';
import { sourceRuleCandidateAPI } from '../../../src/services/sourceRuleCandidates';

vi.mock('../../../src/services/sourceRuleCandidates', () => ({
  sourceRuleCandidateAPI: { get: vi.fn() },
}));

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (_key: string, options?: { defaultValue?: string }) => options?.defaultValue ?? _key,
  }),
}));

describe('Local Modbus review surface', () => {
  it('shows backend register geometry for a multi-register allocation', () => {
    const status: ModbusShareStatus = {
      ...hydratedShareStatus,
      capacity_registers: 64,
      canonical_desired_mappings: [{
        workspace_id: 'workspace-1',
        source_rule_id: 'mc-rule',
        source_rule_revision: 'mc-rule-revision-1',
        tag_id: 'tag-i64',
        tag_key: 'mc.i64',
        display_name: 'I64 value',
        data_type: 'int64',
        share_start_register: 40010,
        zero_based_register: 9,
        span_registers: 4,
        stride_registers: 4,
        capacity_registers: 64,
      }],
    };

    render(<ShareOutputSummary state={mcShareState()} shareStatus={status} />);

    const row = screen.getByTestId('step4-share-output-row-tag-i64');
    expect(row).toHaveTextContent('mc.i64');
    expect(row).toHaveTextContent('int64');
    expect(row).toHaveTextContent('40010');
    expect(row).toHaveTextContent('40013');
    expect(row).toHaveTextContent('9');
    expect(row).toHaveTextContent('4');
    expect(row).toHaveTextContent('64');
  });

  it('keeps a failed allocation visible with safe diagnostic and recovery action', () => {
    render(<ShareOutputSummary state={mcShareState()} shareStatus={{
      ...hydratedShareStatus,
      candidate_set_status: 'blocked',
      candidate_statuses: ['blocked_conflict'],
      outcome: 'invalidated_unknown',
      invalidated_count: 1,
      error: {
        code: 'modbus_share_range_collision',
        message: 'overlapping ranges were rejected',
        retryable: true,
        request_id: 'request-review-1',
        action: 'Retry reconcile after resolving the overlap',
      },
    }} />);

    expect(screen.getByTestId('step4-share-output-failed')).toBeInTheDocument();
    expect(screen.getByTestId('step4-share-output-diagnostic')).toHaveTextContent('modbus_share_range_collision');
    expect(screen.getByTestId('step4-share-output-diagnostic')).toHaveTextContent('Retry reconcile after resolving the overlap');
    expect(screen.getByTestId('step4-share-output-diagnostic')).toHaveTextContent('blocked_conflict');
    expect(screen.getByTestId('step4-share-output-diagnostic')).toHaveTextContent('invalidated');
  });

  it('reads scoped candidate truth and surfaces collision reason plus ownership proof', async () => {
    vi.mocked(sourceRuleCandidateAPI.get).mockResolvedValue({
      source_rule_id: 'mc-rule',
      workspace_id: 'workspace-1',
      workspace_revision: 'workspace-revision-1',
      revision_id: 'mc-rule-revision-1',
      tags: { status: 'ready', candidates: [] },
      database_outputs: { status: 'ready', candidates: [] },
      local_modbus_outputs: {
        status: 'blocked',
        reason: 'range collision',
        candidates: [{
          id: 'candidate-i64',
          identity: {
            source_rule_id: 'mc-rule',
            candidate_type: 'local_modbus_outputs',
            candidate_kind: 'local_modbus',
            derived_from_rule_address: 'D0',
          },
          proposed_signature: 'signature-i64',
          address: 'D0',
          point_id: 'point-i64',
          tag_id: 'tag-i64',
          tag_key: 'mc.i64',
          display_name: 'I64 value',
          data_type: 'int64',
          register_count: 4,
          status: 'blocked_conflict',
          blocking_reason: 'local Modbus range overlaps another candidate',
        }],
      },
    });
    const mapping = {
      ...hydratedShareStatus.canonical_desired_mappings![0],
      tag_id: 'tag-i64',
      source_rule_revision: 'mc-rule-revision-1',
      data_type: 'int64' as const,
      span_registers: 4,
      stride_registers: 4,
      ownership_proof: {
        verified: true,
        workspace_id: 'workspace-1',
        source_rule_id: 'mc-rule',
        source_rule_revision: 'mc-rule-revision-1',
        basis: 'persisted source-rule mapping',
      },
    };
    const status = {
      ...hydratedShareStatus,
      canonical_desired_mappings: [mapping],
      canonical_plan: {
        workspace_id: 'workspace-1',
        workspace_revision: 'workspace-revision-1',
        settings_revision: 'settings-revision-1',
        signature: 'server-plan-signature',
        desired_mappings: [mapping],
      },
    };

    render(<ShareOutputSummary state={mcShareState()} shareStatus={status} />);

    await waitFor(() => expect(sourceRuleCandidateAPI.get).toHaveBeenCalledWith('mc-rule', {
      workspace_id: 'workspace-1',
      expected_workspace_revision: 'workspace-revision-1',
      revision_id: 'mc-rule-revision-1',
    }));
    await waitFor(() => expect(screen.getByTestId('step4-share-output-row-tag-i64')).toHaveTextContent('blocked_conflict'));
    expect(screen.getByTestId('step4-share-output-row-tag-i64')).toHaveTextContent('range overlaps another candidate');
    expect(screen.getByTestId('step4-share-output-row-tag-i64')).toHaveTextContent('verified');
  });

  it('keeps dirty_unknown status fail-closed and shows backend recovery action', () => {
    render(<ShareOutputSummary state={mcShareState()} shareStatus={{
      ...hydratedShareStatus,
      dirty_state: 'dirty_unknown',
      recovery: {
        code: 'modbus_share_dirty_unknown',
        retryable: true,
        action: 'run recovery reconcile',
        request_id: 'request-dirty-1',
      },
    }} />);

    expect(screen.getByTestId('step4-share-output-failed')).toBeInTheDocument();
    expect(screen.getByTestId('step4-share-output-diagnostic')).toHaveTextContent('dirty_unknown');
    expect(screen.getByTestId('step4-share-output-diagnostic')).toHaveTextContent('run recovery reconcile');
  });
});
