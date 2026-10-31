import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { WriteGroupSection } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/WriteGroupSection';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import type { DeviceSourceOutcome, DeviceSourcePort, DeviceSourceRequest } from '@/pages/datalink/workbench-v2/deviceSourceSubmission';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { savedGroupState, TABLE_COLUMNS } from '../../fixtures/writeGroupState';

vi.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key: string, options?: Record<string, unknown>) => (options && typeof options === 'object' && 'count' in options ? `${key}:${options.count}` : key) }),
}));

let metadata: { status: string; columns: typeof TABLE_COLUMNS } = { status: 'exists', columns: TABLE_COLUMNS };
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({
  useStep4TargetColumns: () => ({ columns: metadata.columns, status: metadata.status, refetch: vi.fn() }),
}));

vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({
  studioV2WorkspaceWriteGroupsAPI: {
    list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(),
    testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), get: vi.fn(),
  },
}));

const api = vi.mocked(studioV2WorkspaceWriteGroupsAPI);

function list(groups: WriteGroup[], revision = 'wrev-1') {
  return { workspace_id: 'ws-1', workspace_revision: revision, groups };
}

function renderSection(props: Partial<React.ComponentProps<typeof WriteGroupSection>> & { deviceSource?: DeviceSourcePort } = {}) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
  const utils = render(
    <QueryClientProvider client={client}>
      <WriteGroupSection state={savedGroupState()} workspaceId="ws-1" readonly={false} {...props} />
    </QueryClientProvider>,
  );
  return { client, ...utils };
}

beforeEach(() => {
  vi.clearAllMocks();
  metadata = { status: 'exists', columns: TABLE_COLUMNS };
  api.list.mockResolvedValue(list([]));
  api.delivery.mockResolvedValue({
    group_id: 'group-1', intake: { state: 'not_running' },
    stages: { collecting: 0, queued: 0, retrying: 0, blocked: 0, quarantined: 0, unknown: 0, sql_committed: 0, skipped: 0 },
    last_sql_committed_at: null, oldest_pending_seconds: 0, no_data_buckets: 0, skipped_buckets: 0, backlog: [],
    quota: { configured: false, state: 'unconfigured', scope: '', used_bytes: 0, max_bytes: 0, intake_refused: false },
  });
});


function sourcePort() {
  return { submit: vi.fn<(request: DeviceSourceRequest) => Promise<DeviceSourceOutcome>>(), inspect: vi.fn<DeviceSourcePort['inspect']>(), observations: (): readonly DeviceSourceOutcome[] => [] };
}
async function openSource(port: DeviceSourcePort) {
  renderSection({ deviceSource: port });
  fireEvent.click(await screen.findByTestId('group-create'));
  fireEvent.change(screen.getByTestId('group-name'), { target: { value: 'Unsaved recording name' } });
  fireEvent.click(screen.getByTestId('group-member-include-pt-0'));
  fireEvent.click(screen.getByTestId('recording-source-open'));
}
function enterDevice() {
  fireEvent.change(screen.getByTestId('recording-source-name'), { target: { value: 'New PLC' } });
  fireEvent.change(screen.getByTestId('input-host'), { target: { value: '192.0.2.10' } });
}
function confirmed(request: DeviceSourceRequest): DeviceSourceOutcome {
  return { request, status: 'confirmed', device: { id: 'persisted-server-device', name: request.buffer.name, description: '', protocol: request.buffer.protocol,
    connection_config: JSON.stringify(request.buffer.config), status: 'draft', last_test_at: null, last_test_success: null, last_test_error: '', created_at: '', updated_at: '' } };
}

describe('real recording editor device source consumer', () => {
  it('keeps raw local fields unsent and cancels without changing the recording or existing members', async () => {
    const port = sourcePort(); await openSource(port); enterDevice();
    fireEvent.change(screen.getByTestId('input-port'), { target: { value: '1x' } });
    expect(screen.getByTestId('input-port')).toHaveValue('1x');
    expect(screen.getByTestId('recording-source-save')).toBeDisabled();
    expect(port.submit).not.toHaveBeenCalled();
    fireEvent.click(screen.getByTestId('recording-source-close'));
    expect(screen.getAllByRole('alertdialog')).toHaveLength(1);
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    fireEvent.click(screen.getByTestId('recording-source-discard'));
    expect(screen.queryByRole('alertdialog')).not.toBeInTheDocument();
    expect(screen.getByTestId('group-name')).toHaveValue('Unsaved recording name');
    expect(screen.getByTestId('group-member-include-pt-0')).toBeChecked();
    expect(port.submit).not.toHaveBeenCalled();
    expect(api.create).not.toHaveBeenCalled();
  });
  it('submits exactly once and retains a separately confirmed device when recording closes', async () => {
    const port = sourcePort(); port.submit.mockImplementation(async request => confirmed(request));
    await openSource(port); enterDevice();
    fireEvent.click(screen.getByTestId('recording-source-save'));
    expect(await screen.findByTestId('recording-source-confirmed')).toBeVisible();
    expect(port.submit).toHaveBeenCalledTimes(1);
    expect(port.submit.mock.calls[0][0]).toMatchObject({ workspaceId: 'ws-1', generation: 2, buffer: { name: 'New PLC', config: { host: '192.0.2.10', port: 502 } } });
    fireEvent.click(screen.getByTestId('recording-source-close'));
    expect(screen.getByTestId('group-name')).toHaveValue('Unsaved recording name');
    fireEvent.click(screen.getByTestId('group-editor-close'));
    expect(api.remove).not.toHaveBeenCalled();
    expect(api.create).not.toHaveBeenCalled();
    expect(port.submit).toHaveBeenCalledTimes(1);
  });
  it('settles the submitted generation without clearing newer input or resubmitting a create', async () => {
    const port = sourcePort(); let resolve!: (result: DeviceSourceOutcome) => void;
    port.submit.mockImplementation(() => new Promise(done => { resolve = done; }));
    await openSource(port); enterDevice();
    fireEvent.click(screen.getByTestId('recording-source-save'));
    fireEvent.change(screen.getByTestId('recording-source-name'), { target: { value: 'Newer unsent name' } });
    const request = port.submit.mock.calls[0][0]; resolve(confirmed(request));
    expect(await screen.findByTestId('recording-source-newer')).toBeVisible();
    expect(screen.getByTestId('recording-source-name')).toHaveValue('Newer unsent name');
    expect(screen.getByRole('dialog')).toBeVisible();
    expect(port.submit).toHaveBeenCalledTimes(1);
    fireEvent.click(screen.getByTestId('recording-source-save'));
    expect(port.submit.mock.calls[1][0].deviceId).toBe('persisted-server-device');
  });
  it('retains unknown operation and input, using inspect instead of replacement create', async () => {
    const port = sourcePort(); port.submit.mockImplementation(async request => ({ request, status: 'unknown' }));
    await openSource(port); enterDevice(); fireEvent.click(screen.getByTestId('recording-source-save'));
    expect(await screen.findByTestId('recording-source-unknown')).toBeVisible();
    expect(screen.getByTestId('recording-source-save')).toBeDisabled();
    fireEvent.click(screen.getByTestId('recording-source-inspect'));
    expect(port.inspect).toHaveBeenCalledWith(port.submit.mock.calls[0][0].requestId);
    expect(port.submit).toHaveBeenCalledTimes(1);
    expect(screen.getByTestId('recording-source-name')).toHaveValue('New PLC');
  });
});

describe('source retained-operation recovery consumer', () => {
  it('resumes an exact unknown request after editor remount without creating a replacement', async () => {
    const port = sourcePort();
    const request: DeviceSourceRequest = { workspaceId: 'ws-1', sessionId: 'older-editor', requestId: 'exact-request', generation: 4,
      buffer: { name: 'Unresolved PLC', description: '', protocol: 'modbus_tcp', config: { host: '192.0.2.20', port: 502, slave_id: 1, timeout: 5 } } };
    port.observations = () => [{ request, status: 'unknown' }];
    port.inspect.mockResolvedValue({ request, status: 'unknown' });
    renderSection({ deviceSource: port });
    fireEvent.click(await screen.findByTestId('group-create'));
    expect(screen.getByTestId('recording-source-open')).toBeDisabled();
    fireEvent.click(screen.getByTestId('recording-source-recover-exact-request'));
    expect(screen.getByTestId('recording-source-name')).toHaveValue('Unresolved PLC');
    fireEvent.click(screen.getByTestId('recording-source-inspect'));
    await waitFor(() => expect(port.inspect).toHaveBeenCalledWith('exact-request'));
    expect(port.submit).not.toHaveBeenCalled();
    expect(screen.getByTestId('recording-source-save')).toBeDisabled();
  });
  it('keeps the unknown state when inspect fails instead of reporting cancellation', async () => {
    const port = sourcePort(); port.submit.mockImplementation(async request => ({ request, status: 'unknown' }));
    port.inspect.mockRejectedValue(new Error('unavailable'));
    await openSource(port); enterDevice(); fireEvent.click(screen.getByTestId('recording-source-save'));
    await screen.findByTestId('recording-source-unknown');
    fireEvent.click(screen.getByTestId('recording-source-inspect'));
    await waitFor(() => expect(screen.getByTestId('recording-source-inspect')).toBeEnabled());
    expect(screen.getByTestId('recording-source-unknown')).toBeVisible();
    expect(port.submit).toHaveBeenCalledTimes(1);
  });
});

describe('saved-but-unavailable source and next explicit device', () => {
  it('keeps saved-but-unavailable distinct and retries only the acknowledged device identity', async () => {
    const port = sourcePort(); port.submit.mockImplementation(async request => ({ ...confirmed(request), status: 'persisted-unavailable' }));
    await openSource(port); enterDevice(); fireEvent.click(screen.getByTestId('recording-source-save'));
    expect(await screen.findByTestId('recording-source-persisted-unavailable')).toBeVisible();
    expect(screen.queryByTestId('recording-source-confirmed')).not.toBeInTheDocument();
    fireEvent.change(screen.getByTestId('input-host'), { target: { value: '192.0.2.11' } });
    fireEvent.click(screen.getByTestId('recording-source-save'));
    expect(port.submit).toHaveBeenCalledTimes(2);
    expect(port.submit.mock.calls[1][0].deviceId).toBe('persisted-server-device');
  });
  it('opens a fresh unsent form only on an explicit next source action after confirmed close', async () => {
    const port = sourcePort(); port.submit.mockImplementation(async request => confirmed(request));
    await openSource(port); enterDevice(); fireEvent.click(screen.getByTestId('recording-source-save'));
    await screen.findByTestId('recording-source-confirmed');
    fireEvent.click(screen.getByTestId('recording-source-close'));
    fireEvent.click(screen.getByTestId('recording-source-open'));
    expect(screen.getByTestId('recording-source-name')).toHaveValue('');
    expect(port.submit).toHaveBeenCalledTimes(1);
    expect(screen.getByTestId('group-name')).toHaveValue('Unsaved recording name');
  });
});

describe('single source modal guard', () => {
  it('uses the same dirty guard for Escape and preserves raw fields when keeping edits', async () => {
    const port = sourcePort(); await openSource(port); enterDevice();
    fireEvent.change(screen.getByTestId('input-port'), { target: { value: ' 1x ' } });
    fireEvent.keyDown(screen.getByTestId('input-port'), { key: 'Escape' });
    expect(screen.getByRole('alertdialog')).toBeVisible();
    fireEvent.keyDown(screen.getByRole('alertdialog'), { key: 'Escape' });
    expect(screen.getByRole('dialog')).toBeVisible();
    expect(screen.getByTestId('input-port')).toHaveValue(' 1x ');
    expect(port.submit).not.toHaveBeenCalled();
  });
});

describe('independent unresolved protocol intent', () => {
 it('does not allow a pending-create protocol to diverge after unknown and become locked on recovery', async () => {
  const port=sourcePort(); port.submit.mockImplementation(async request=>({request,status:'unknown'}));
  await openSource(port); enterDevice(); fireEvent.click(screen.getByTestId('recording-source-save'));
  await screen.findByTestId('recording-source-unknown');
  const select=screen.getByLabelText('device.protocol');
  expect(select).toBeDisabled();
  fireEvent.change(screen.getByTestId('input-port'), { target: { value: 'new invalid 1x' } });
  const request = port.submit.mock.calls[0][0]; port.inspect.mockResolvedValue(confirmed(request));
  fireEvent.click(screen.getByTestId('recording-source-inspect'));
  await screen.findByTestId('recording-source-newer');
  expect(screen.getByLabelText('device.protocol')).toHaveValue('modbus_tcp');
  expect(screen.getByTestId('input-port')).toHaveValue('new invalid 1x');
  expect(port.submit).toHaveBeenCalledTimes(1);
 });
});
