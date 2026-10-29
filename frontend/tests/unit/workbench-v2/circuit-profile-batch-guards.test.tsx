import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi, type Mock } from 'vitest';
import {
  CircuitProfileBatchEditor,
  type CircuitProfileBatchEditorProps,
} from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileBatchEditor';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { savedGroup } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string, options?: Record<string, unknown>) =>
      options && 'count' in options ? `${key}:${options.count}` : key,
  }),
}));

function makeGroup(circuitId: string, circuitName: string, intervalSeconds: number, id: string): WriteGroup {
  return savedGroup({
    id,
    revision: `rev-${id}-1`,
    name: `${circuitName} (Instantaneous)`,
    status: 'ready',
    layout_version: 'managed_long_v1',
    recording_profile: { kind: 'instantaneous', interval: intervalSeconds, unit: 'second' },
    storage_mode: 'plain',
    destination: { ...savedGroup().destination, storage_strategy: 'managed', table_name: 'gw_long' },
    members: [
      {
        device_id: 'dev-1', point_id: `pt-${id}-1`, tag_id: `tag-${id}-1`,
        source_revision: 's1', mapping_revision: 'm1', target_column: '', required: false,
        circuit_id: circuitId, circuit_name: circuitName,
      },
    ],
    row_policy: { interval_seconds: intervalSeconds, allowed_lateness_seconds: 0, incomplete_policy: 'partial' },
    write_policy: { mode: 'append', dedupe_capability: 'receipt', recreate_missing: true },
  });
}

describe('CircuitProfileBatchEditor boundary guards and partial failures', () => {
  let updateWriteGroup: Mock<CircuitProfileBatchEditorProps['updateWriteGroup']>;
  let ensureCircuitProfile: Mock<NonNullable<CircuitProfileBatchEditorProps['ensureCircuitProfile']>>;
  let onSuccess: Mock<CircuitProfileBatchEditorProps['onSuccess']>;
  let onClose: Mock<CircuitProfileBatchEditorProps['onClose']>;
  let onOperationActiveChange: Mock<NonNullable<CircuitProfileBatchEditorProps['onOperationActiveChange']>>;
  let savedLongGroups: WriteGroup[];

  beforeEach(() => {
    vi.clearAllMocks();
    savedLongGroups = [
      makeGroup('cir-a', 'Feeder A', 15, 'grp-a'),
      makeGroup('cir-b', 'Feeder B', 60, 'grp-b'),
    ];
    updateWriteGroup = vi.fn<CircuitProfileBatchEditorProps['updateWriteGroup']>();
    ensureCircuitProfile = vi.fn<NonNullable<CircuitProfileBatchEditorProps['ensureCircuitProfile']>>();
    onSuccess = vi.fn<CircuitProfileBatchEditorProps['onSuccess']>();
    onClose = vi.fn<CircuitProfileBatchEditorProps['onClose']>();
    onOperationActiveChange = vi.fn<NonNullable<CircuitProfileBatchEditorProps['onOperationActiveChange']>>();
  });

  it('detects external workspaceRevision drift and blocks confirm with stale notice', () => {
    const { rerender } = render(
      <CircuitProfileBatchEditor
        capabilityContext={{ connectorKind: 'sqlite', destination: { connector_id: 'conn-1', connector_revision: 'crev-1' } }}
        workspaceId="ws-1"
        workspaceRevision="wrev-1"
        savedLongGroups={savedLongGroups}
        readonly={false}
        onClose={onClose}
        onSuccess={onSuccess}
        updateWriteGroup={updateWriteGroup}
        ensureCircuitProfile={ensureCircuitProfile}
        onOperationActiveChange={onOperationActiveChange}
      />,
    );

    fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
    expect(screen.getByTestId('circuit-profile-batch-confirm')).toBeEnabled();

    // External change occurs: workspaceRevision bumps to wrev-2
    rerender(
      <CircuitProfileBatchEditor
        capabilityContext={{ connectorKind: 'sqlite', destination: { connector_id: 'conn-1', connector_revision: 'crev-1' } }}
        workspaceId="ws-1"
        workspaceRevision="wrev-2"
        savedLongGroups={savedLongGroups}
        readonly={false}
        onClose={onClose}
        onSuccess={onSuccess}
        updateWriteGroup={updateWriteGroup}
        ensureCircuitProfile={ensureCircuitProfile}
        onOperationActiveChange={onOperationActiveChange}
      />,
    );

    expect(screen.getByRole('alert')).toHaveTextContent('step4.basic.reason.stale_intent');
    expect(screen.getByTestId('circuit-profile-batch-confirm')).toBeDisabled();
    fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm'));
    expect(updateWriteGroup).not.toHaveBeenCalled();
  });

  it('preserves partial success and updates parent revision when a subsequent update fails', async () => {
    updateWriteGroup
      .mockResolvedValueOnce({
        workspace_revision: 'wrev-2',
        group: { ...savedLongGroups[0], revision: 'rev-grp-a-2' },
      })
      .mockRejectedValueOnce(new Error('stale_group_revision: rev-grp-b-1'));

    render(
      <CircuitProfileBatchEditor
        capabilityContext={{ connectorKind: 'sqlite', destination: { connector_id: 'conn-1', connector_revision: 'crev-1' } }}
        workspaceId="ws-1"
        workspaceRevision="wrev-1"
        savedLongGroups={savedLongGroups}
        readonly={false}
        onClose={onClose}
        onSuccess={onSuccess}
        updateWriteGroup={updateWriteGroup}
        ensureCircuitProfile={ensureCircuitProfile}
        onOperationActiveChange={onOperationActiveChange}
      />,
    );

    fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
    fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm'));

    await waitFor(() => {
      // onSuccess should be called with partial results (grp-a succeeded with wrev-2)
      expect(onSuccess).toHaveBeenCalledWith(
        [expect.objectContaining({ id: 'grp-a' })],
        'wrev-2',
      );
    });

    // Error is rendered with safe message without raw backend error
    expect(screen.getByRole('alert')).toBeInTheDocument();
    expect(screen.queryByText(/stale_group_revision/)).not.toBeInTheDocument();
    // Modal stays open so user can inspect or retry
    expect(onClose).not.toHaveBeenCalled();
    expect(onOperationActiveChange).toHaveBeenCalledWith(false);
  });

  it('guards against duplicate confirm submissions while in-flight', async () => {
    let resolveFirst: (val: { workspace_revision: string; group: WriteGroup }) => void = () => {};
    updateWriteGroup.mockReturnValue(new Promise((resolve) => {
      resolveFirst = resolve;
    }));

    render(
      <CircuitProfileBatchEditor
        capabilityContext={{ connectorKind: 'sqlite', destination: { connector_id: 'conn-1', connector_revision: 'crev-1' } }}
        workspaceId="ws-1"
        workspaceRevision="wrev-1"
        savedLongGroups={savedLongGroups}
        readonly={false}
        onClose={onClose}
        onSuccess={onSuccess}
        updateWriteGroup={updateWriteGroup}
        ensureCircuitProfile={ensureCircuitProfile}
        onOperationActiveChange={onOperationActiveChange}
      />,
    );

    fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
    const confirmBtn = screen.getByTestId('circuit-profile-batch-confirm');
    fireEvent.click(confirmBtn);

    // Second click during in-flight save
    fireEvent.click(confirmBtn);

    expect(updateWriteGroup).toHaveBeenCalledTimes(1);
    expect(onOperationActiveChange).toHaveBeenCalledWith(true);

    // Complete first
    resolveFirst({
      workspace_revision: 'wrev-2',
      group: { ...savedLongGroups[0], revision: 'rev-grp-a-2' },
    });
  });

  it('disables preview and confirm when readonly is true', () => {
    render(
      <CircuitProfileBatchEditor
        capabilityContext={{ connectorKind: 'sqlite', destination: { connector_id: 'conn-1', connector_revision: 'crev-1' } }}
        workspaceId="ws-1"
        workspaceRevision="wrev-1"
        savedLongGroups={savedLongGroups}
        readonly={true}
        onClose={onClose}
        onSuccess={onSuccess}
        updateWriteGroup={updateWriteGroup}
        ensureCircuitProfile={ensureCircuitProfile}
      />,
    );

    expect(screen.getByTestId('circuit-profile-batch-preview-btn')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-batch-interval-inst')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-batch-interval-cum')).toBeDisabled();
  });
});
