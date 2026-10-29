import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi, type Mock } from 'vitest';
import {
  CircuitProfileBatchEditor,
  type CircuitProfileBatchEditorProps,
} from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileBatchEditor';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { httpFailure, longProfileGroup } from '../../fixtures/circuitBatchPanel';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string) => key }),
}));

type Props = CircuitProfileBatchEditorProps;
const groups = (): WriteGroup[] => [
  longProfileGroup({ id: 'grp-a', circuitId: 'cir-a', circuitName: 'Feeder A', interval: 15, points: [0, 1] }),
  longProfileGroup({ id: 'grp-b', circuitId: 'cir-b', circuitName: 'Feeder B', interval: 60, points: [2, 3] }),
];
const savedUpdate = async ({ id, request }: { id: string; request: { group: object } }) =>
  ({ workspace_revision: `wrev-${id}`, group: { ...(request.group as WriteGroup), id, revision: `${id}-saved` } });

let updateWriteGroup: Mock<Props['updateWriteGroup']>;
let onSuccess: Mock<Props['onSuccess']>;
let onClose: Mock<Props['onClose']>;
let onBusyChange: Mock<NonNullable<Props['onBusyChange']>>;
const unhandled: unknown[] = [];
const trackUnhandled = (reason: unknown) => { unhandled.push(reason); };

/** Props this slice adds are passed loosely so the suite compiles before they exist. */
function editor(extra: Record<string, unknown> = {}, overrides: Partial<Props> = {}) {
  return (
    <CircuitProfileBatchEditor
      capabilityContext={{ connectorKind: 'sqlite', destination: { connector_id: 'conn-1', connector_revision: 'crev-1' } }}
      workspaceId="ws-1" workspaceRevision="wrev-1" connectorRevision="crev-1" savedLongGroups={groups()} readonly={false}
      onClose={onClose} onSuccess={onSuccess} updateWriteGroup={updateWriteGroup} onBusyChange={onBusyChange}
      {...overrides} {...extra}
    />
  );
}
function review() {
  fireEvent.click(screen.getByTestId('circuit-profile-batch-preview-btn'));
}
function confirm() {
  fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm'));
}

beforeEach(() => {
  vi.clearAllMocks();
  unhandled.length = 0;
  process.on('unhandledRejection', trackUnhandled);
  updateWriteGroup = vi.fn<Props['updateWriteGroup']>(savedUpdate);
  onSuccess = vi.fn<Props['onSuccess']>();
  onClose = vi.fn<Props['onClose']>();
  onBusyChange = vi.fn<NonNullable<Props['onBusyChange']>>();
});
afterEach(() => { process.off('unhandledRejection', trackUnhandled); });

describe('Batch save waits for the parent to reload the server authority', () => {
  it('keeps the lock and the editor open until the awaited parent refresh settles', async () => {
    let finishRefresh!: () => void;
    onSuccess.mockImplementation(() => new Promise<void>((resolve) => { finishRefresh = resolve; }));
    render(editor());
    review();
    confirm();
    await waitFor(() => expect(onSuccess).toHaveBeenCalledTimes(1));

    expect(onClose).not.toHaveBeenCalled();
    expect(onBusyChange).not.toHaveBeenCalledWith(false);
    expect(screen.getByTestId('circuit-profile-batch-close')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-batch-confirm')).toBeDisabled();

    await act(async () => finishRefresh());
    expect(onClose).toHaveBeenCalledTimes(1);
    expect(onBusyChange).toHaveBeenLastCalledWith(false);
    expect(updateWriteGroup).toHaveBeenCalledTimes(2);
  });

  it('shows a refresh failure after a committed save without calling it a failed mutation or resending', async () => {
    onSuccess.mockRejectedValue(new Error('GET /write-groups failed: dsn=postgres://user:secret@db/prod'));
    const view = render(editor());
    review();
    confirm();

    expect(await screen.findByTestId('circuit-profile-batch-refresh-error')).toHaveTextContent('step4.basic.circuit.batch_refresh_failed');
    await act(async () => { await new Promise((resolve) => setTimeout(resolve, 0)); });
    expect(unhandled).toEqual([]);
    expect(onClose).not.toHaveBeenCalled();
    expect(screen.getByTestId('circuit-profile-batch-status-cir-a-instantaneous')).toHaveTextContent('step4.basic.circuit.batch_status_saved');
    expect(screen.getByTestId('circuit-profile-batch-status-cir-b-instantaneous')).toHaveTextContent('step4.basic.circuit.batch_status_saved');
    expect(view.container.textContent).not.toMatch(/secret|postgres:/);
    expect(onBusyChange).toHaveBeenLastCalledWith(false);

    fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm'));
    expect(updateWriteGroup).toHaveBeenCalledTimes(2);
  });
});

describe('An unconfirmed save keeps the reviewed intent and offers a way to check', () => {
  it('shows unknown, keeps the old to new review, does not resend and exposes the check entry', async () => {
    updateWriteGroup.mockRejectedValue(new Error('response lost after server commit'));
    const onVerify = vi.fn();
    render(editor({ onVerify }));
    review();
    confirm();

    expect(await screen.findByTestId('circuit-profile-batch-status-cir-a-instantaneous')).toHaveTextContent('step4.basic.circuit.batch_status_unknown');
    expect(screen.getByTestId('circuit-profile-batch-preview')).toHaveTextContent('15s → 30s');
    expect(screen.getByTestId('circuit-profile-batch-unknown-note')).toBeInTheDocument();
    expect(screen.getByTestId('circuit-profile-batch-confirm')).toBeDisabled();
    expect(screen.getByTestId('circuit-profile-batch-status-cir-b-instantaneous')).toHaveTextContent('step4.basic.circuit.batch_status_not_attempted');
    expect(updateWriteGroup).toHaveBeenCalledTimes(1);
    expect(onClose).not.toHaveBeenCalled();

    fireEvent.click(screen.getByTestId('circuit-profile-batch-verify'));
    await waitFor(() => expect(onVerify).toHaveBeenCalledTimes(1));
    expect(updateWriteGroup).toHaveBeenCalledTimes(1);
  });

  it('reads an unknown profile back as saved only when the reloaded group proves the reviewed interval', async () => {
    updateWriteGroup.mockRejectedValue(new Error('response lost after server commit'));
    const view = render(editor({ onVerify: vi.fn() }));
    review();
    confirm();
    await screen.findByTestId('circuit-profile-batch-status-cir-a-instantaneous');

    const committed = longProfileGroup({ id: 'grp-a', circuitId: 'cir-a', circuitName: 'Feeder A', interval: 30, points: [0, 1] });
    view.rerender(editor({ onVerify: vi.fn() }, { savedLongGroups: [{ ...committed, revision: 'rev-grp-a-2' }, groups()[1]] }));
    expect(screen.getByTestId('circuit-profile-batch-status-cir-a-instantaneous')).toHaveTextContent('step4.basic.circuit.batch_status_saved');
    expect(screen.getByTestId('circuit-profile-batch-status-cir-b-instantaneous')).toHaveTextContent('step4.basic.circuit.batch_status_not_attempted');
  });
});

describe('Operator copy never carries raw backend text', () => {
  it('maps a typed 409 to its code copy and hides SQL, DSN and exception text', async () => {
    updateWriteGroup.mockRejectedValue(httpFailure(409, 'WRITE_GROUP_CIRCUIT_PROFILE_CONFLICT', 'pq: duplicate key postgres://user:secret@db/prod'));
    const view = render(editor());
    review();
    confirm();

    expect(await screen.findByTestId('circuit-profile-batch-status-cir-a-instantaneous')).toHaveTextContent('step4.basic.circuit.batch_status_failed');
    expect(view.container).toHaveTextContent('step4.basic.errors.circuit_profile_conflict');
    expect(view.container.textContent).not.toMatch(/secret|postgres:|pq:|duplicate key/);
  });

  it('does not render raw text for an unknown outcome either', async () => {
    updateWriteGroup.mockRejectedValue(new Error('read tcp 10.0.0.5:5432: i/o timeout, password=hunter2'));
    const view = render(editor());
    review();
    confirm();

    expect(await screen.findByTestId('circuit-profile-batch-status-cir-a-instantaneous')).toHaveTextContent('step4.basic.circuit.batch_status_unknown');
    expect(view.container.textContent).not.toMatch(/hunter2|10\.0\.0\.5|i\/o timeout/);
  });
});

describe('Preview and in-flight request are bound to the scope that reviewed them', () => {
  it('invalidates the reviewed preview when the scope identity changes before confirmation', () => {
    const view = render(editor({ scopeKey: 'ws-1/dev-1/cir-a/crev-1' }));
    review();
    expect(screen.getByTestId('circuit-profile-batch-confirm')).toBeEnabled();

    view.rerender(editor({ scopeKey: 'ws-1/dev-1/cir-a/crev-2' }));
    expect(screen.getByTestId('circuit-profile-batch-confirm')).toBeDisabled();
    expect(screen.getByRole('alert')).toHaveTextContent('step4.basic.reason.stale_intent');
    fireEvent.click(screen.getByTestId('circuit-profile-batch-confirm'));
    expect(updateWriteGroup).not.toHaveBeenCalled();
  });

  it('still releases its own lock and reports to the parent when unmounted mid-request, without closing a newer editor', async () => {
    let finishFirst!: (value: Awaited<ReturnType<typeof savedUpdate>>) => void;
    updateWriteGroup.mockImplementationOnce(() => new Promise((resolve) => { finishFirst = resolve as typeof finishFirst; }));
    const view = render(editor());
    review();
    confirm();
    await waitFor(() => expect(updateWriteGroup).toHaveBeenCalledTimes(1));
    view.unmount();

    await act(async () => finishFirst(await savedUpdate({ id: 'grp-a', request: { group: groups()[0] } })));
    await waitFor(() => expect(onSuccess).toHaveBeenCalledTimes(1));
    expect(onClose).not.toHaveBeenCalled();
    expect(onBusyChange.mock.calls.filter(([busy]) => busy === false)).toHaveLength(1);
  });
});
