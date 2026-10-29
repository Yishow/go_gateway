import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi, type Mock } from 'vitest';
import {
  CircuitProfileBatchEditor,
  type CircuitProfileBatchEditorProps,
} from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileBatchEditor';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { httpFailure, longProfileGroup } from '../../fixtures/circuitBatchPanel';

vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));

type Props = CircuitProfileBatchEditorProps;
const instGroup = (interval = 15, revision = 'rev-inst-1'): WriteGroup =>
  ({ ...longProfileGroup({ id: 'inst-a', circuitId: 'cir-a', circuitName: 'Feeder A', interval, points: [0, 1] }), revision });
const cumGroup = (interval = 10, revision = 'rev-cum-1'): WriteGroup =>
  ({ ...longProfileGroup({ id: 'cum-a', circuitId: 'cir-a', circuitName: 'Feeder A', interval, points: [0, 1], kind: 'cumulative' }), revision });

let updateWriteGroup: Mock<Props['updateWriteGroup']>;
let onSuccess: Mock<Props['onSuccess']>;
let onClose: Mock<Props['onClose']>;

function editor(groups: WriteGroup[], workspaceRevision = 'wrev-1', extra: Record<string, unknown> = {}) {
  return (
    <CircuitProfileBatchEditor
      capabilityContext={{ connectorKind: 'sqlite', destination: { connector_id: 'conn-1', connector_revision: 'crev-1' } }}
      workspaceId="ws-1" workspaceRevision={workspaceRevision} connectorRevision="crev-1" savedLongGroups={groups} readonly={false}
      onClose={onClose} onSuccess={onSuccess} updateWriteGroup={updateWriteGroup} {...extra}
    />
  );
}
function reviewBoth(firstTime = true) {
  if (firstTime) fireEvent.click(screen.getByLabelText('step4.basic.circuit.batch_cum_interval'));
  fireEvent.change(screen.getByTestId('circuit-profile-batch-interval-inst'), { target: { value: '30' } });
  fireEvent.change(screen.getByTestId('circuit-profile-batch-interval-cum'), { target: { value: '20' } });
  fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
}
const saved = (revision: string) => async ({ id, request }: { id: string; request: { group: object } }) =>
  ({ workspace_revision: revision, group: { ...(request.group as WriteGroup), id, revision: `${id}-saved` } });

beforeEach(() => {
  vi.clearAllMocks();
  updateWriteGroup = vi.fn<Props['updateWriteGroup']>();
  onSuccess = vi.fn<Props['onSuccess']>();
  onClose = vi.fn<Props['onClose']>();
});

describe('Dual-profile batch shows which writer was saved and reviews only what is left', () => {
  it('reports the saved instantaneous writer and the rejected cumulative writer separately, then reviews only the cumulative change', async () => {
    updateWriteGroup
      .mockImplementationOnce(saved('wrev-2'))
      .mockRejectedValueOnce(httpFailure(409, 'WRITE_GROUP_CIRCUIT_PROFILE_CONFLICT'));
    const view = render(editor([instGroup(), cumGroup()]));
    reviewBoth();
    fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm'));

    await screen.findByTestId('circuit-profile-batch-status-cir-a-cumulative');
    expect(screen.getByTestId('circuit-profile-batch-status-cir-a-instantaneous')).toHaveTextContent('step4.basic.circuit.batch_status_saved');
    expect(screen.getByTestId('circuit-profile-batch-status-cir-a-cumulative')).toHaveTextContent('step4.basic.circuit.batch_status_failed');
    expect(screen.getByTestId('circuit-profile-batch-status-cir-a')).toHaveTextContent('step4.basic.circuit.batch_status_partial');
    expect(onSuccess).toHaveBeenCalledWith([expect.objectContaining({ id: 'inst-a' })], 'wrev-2');
    expect(screen.getByRole('alert')).toHaveTextContent('step4.basic.errors.circuit_profile_conflict');
    expect(updateWriteGroup).toHaveBeenCalledTimes(2);

    // The parent reloaded: the instantaneous writer is now 30 s with a new revision. A fresh review starts from that.
    fireEvent.click(screen.getByTestId('circuit-profile-batch-cancel'));
    view.rerender(editor([instGroup(30, 'rev-inst-2'), cumGroup()], 'wrev-2'));
    reviewBoth(false);
    expect(screen.getByTestId('circuit-profile-batch-preview')).toHaveTextContent('step4.basic.circuit.instantaneous_label: 30s');
    expect(screen.getByTestId('circuit-profile-batch-preview')).toHaveTextContent('10m → 20m');
    expect(screen.queryByTestId('circuit-profile-batch-status-cir-a-instantaneous')).not.toBeInTheDocument();

    updateWriteGroup.mockImplementationOnce(saved('wrev-3'));
    fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm'));
    await waitFor(() => expect(onClose).toHaveBeenCalledTimes(1));
    expect(updateWriteGroup).toHaveBeenCalledTimes(3);
    expect(updateWriteGroup.mock.calls[2][0]).toEqual(expect.objectContaining({ id: 'cum-a' }));
    expect(updateWriteGroup.mock.calls[2][0].request.expected_workspace_revision).toBe('wrev-2');
    expect(screen.getByTestId('circuit-profile-batch-status-cir-a-instantaneous')).toHaveTextContent('step4.basic.circuit.batch_status_unchanged');
  });

  it('keeps a lost cumulative reply unknown, blocks a silent re-review and never resends the writer', async () => {
    updateWriteGroup
      .mockImplementationOnce(saved('wrev-2'))
      .mockRejectedValueOnce(new Error('socket hang up'));
    render(editor([instGroup(), cumGroup()], 'wrev-1', { onVerify: vi.fn() }));
    reviewBoth();
    fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm'));

    await screen.findByTestId('circuit-profile-batch-unknown-note');
    expect(screen.getByTestId('circuit-profile-batch-status-cir-a-instantaneous')).toHaveTextContent('step4.basic.circuit.batch_status_saved');
    expect(screen.getByTestId('circuit-profile-batch-status-cir-a-cumulative')).toHaveTextContent('step4.basic.circuit.batch_status_unknown');
    expect(screen.getByTestId('circuit-profile-batch-cancel')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-batch-confirm')).toBeDisabled();
    expect(updateWriteGroup).toHaveBeenCalledTimes(2);
    expect(onClose).not.toHaveBeenCalled();
  });
});

describe('Closing after an unknown outcome leaves nothing to resend', () => {
  it('opens the next editor without the old plan and sends again only after a new explicit review', async () => {
    updateWriteGroup.mockRejectedValueOnce(new Error('socket hang up'));
    const first = render(editor([instGroup(), cumGroup()]));
    reviewBoth();
    fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm'));
    await screen.findByTestId('circuit-profile-batch-unknown-note');
    first.unmount();

    render(editor([instGroup(), cumGroup()]));
    expect(screen.queryByTestId('circuit-profile-batch-preview')).not.toBeInTheDocument();
    expect(screen.queryByTestId('circuit-profile-batch-unknown-note')).not.toBeInTheDocument();
    await new Promise((resolve) => setTimeout(resolve, 0));
    expect(updateWriteGroup).toHaveBeenCalledTimes(1);

    updateWriteGroup.mockImplementation(saved('wrev-2'));
    fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
    expect(screen.queryByTestId('circuit-profile-batch-status-cir-a-instantaneous')).not.toBeInTheDocument();
    expect(updateWriteGroup).toHaveBeenCalledTimes(1);
    fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm'));
    await waitFor(() => expect(updateWriteGroup).toHaveBeenCalledTimes(2));
  });
});
