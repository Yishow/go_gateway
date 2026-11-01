import type { StudioV2WorkspaceDeviceRecord } from '../../../src/services/studioV2WorkspaceDevices';
import { useEffect, useRef } from 'react';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { useRecordingDeviceBuffer } from '../../../src/features/datalink/workbench-v2/recording-editor/useRecordingDeviceBuffer';
import { RecordingSourceDeviceContent } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingSourceDeviceContent';
import type { DeviceSourceOutcome, DeviceSourcePort } from '../../../src/pages/datalink/workbench-v2/deviceSourceSubmission';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
const tcp = { protocol: 'fatek_fbs' as const, config: { mode: 'tcp', host: '192.0.2.34', port: 500, station_no: 3, timeout: 5 } };
function setup() {
  const recentConnections = { read: vi.fn((_scope: string, protocol?: string, mode?: string) =>
    (!protocol || protocol === 'fatek_fbs' && mode !== 'serial') ? structuredClone(tcp) : undefined), clear: vi.fn() };
  const port: DeviceSourcePort = { recentConnections, observations: () => [], submit: vi.fn(), inspect: vi.fn() };
  const slot = { active: true, closeRef: { current: undefined }, onConfirmChange: vi.fn(), onReturn: vi.fn(), onSaved: vi.fn() };
  function Form() {
    const controller = useRecordingDeviceBuffer('workspace-1', port);
    const start = useRef(controller.start);
    useEffect(() => { void start.current(); }, []); // A newly opened source form.
    return <RecordingSourceDeviceContent controller={controller} slot={slot} />;
  }
  return { ...render(<Form />), port, recentConnections, slot };
}
describe('recent connection reuse controls use the existing source guard', () => {
  it('shows reuse and endpoint reminder, and clears name-only edits without discarding the name', async () => {
    const current = setup(); await screen.findByTestId('recording-source-reuse');
    expect(screen.getByText('recentConnection.verifyEndpoint')).toBeVisible();
    fireEvent.change(screen.getByTestId('recording-source-name'), { target: { value: 'Second PLC' } });
    fireEvent.click(screen.getByTestId('recording-source-clear-recent'));
    expect(screen.getByTestId('recording-source-name')).toHaveValue('Second PLC');
    expect(screen.queryByTestId('recording-source-reuse')).not.toBeInTheDocument();
    expect(screen.getByTestId('input-host')).toHaveValue('192.168.1.100');
    expect(current.recentConnections.clear).toHaveBeenCalledWith('workspace-1', 'fatek_fbs');
    expect(current.port.submit).not.toHaveBeenCalled();
  });
  it('keeps cache and edited raw input exactly unchanged when clear confirmation is cancelled', async () => {
    const current = setup(); await screen.findByTestId('recording-source-reuse');
    fireEvent.change(screen.getByTestId('input-port'), { target: { value: ' raw invalid ' } });
    fireEvent.click(screen.getByTestId('recording-source-clear-recent'));
    await waitFor(() => expect(current.slot.onConfirmChange).toHaveBeenLastCalledWith(true));
    fireEvent.click(screen.getByText('device.keep'));
    expect(screen.getByTestId('input-port')).toHaveValue(' raw invalid ');
    expect(screen.getByTestId('recording-source-reuse')).toBeVisible();
    expect(current.recentConnections.clear).not.toHaveBeenCalled();
    fireEvent.click(screen.getByTestId('recording-source-clear-recent'));
    fireEvent.click(screen.getByTestId('recording-source-discard'));
    expect(screen.getByTestId('input-port')).toHaveValue('500');
    expect(current.recentConnections.clear).toHaveBeenCalledOnce();
    expect(current.port.submit).not.toHaveBeenCalled();
  });
  it('switches FATEK subtype using only matching cache, stripping inactive TCP fields', async () => {
    const current = setup(); await screen.findByTestId('recording-source-reuse');
    fireEvent.change(screen.getByTestId('select-fatek-mode'), { target: { value: 'serial' } });
    expect(screen.getByTestId('input-rtu-port')).toHaveValue('');
    expect(screen.queryByTestId('recording-source-reuse')).not.toBeInTheDocument();
    expect(current.recentConnections.read).toHaveBeenLastCalledWith('workspace-1', 'fatek_fbs', 'serial');
    fireEvent.change(screen.getByTestId('select-fatek-mode'), { target: { value: 'tcp' } });
    fireEvent.click(screen.getByTestId('recording-source-discard'));
    expect(screen.getByTestId('input-host')).toHaveValue('192.0.2.34');
    expect(screen.getByTestId('recording-source-reuse')).toBeVisible();
    expect(current.port.submit).not.toHaveBeenCalled();
  });
});

describe('pending new-source subtype changes', () => {
  it('uses only matching subtype defaults after dirty confirmation without releasing the pending save', async () => {
    const current = setup(); await screen.findByTestId('recording-source-reuse');
    vi.mocked(current.port.submit).mockImplementation(() => new Promise(() => {}));
    fireEvent.change(screen.getByTestId('recording-source-name'), { target: { value: 'Pending PLC' } });
    fireEvent.click(screen.getByTestId('recording-source-save'));
    await screen.findByText('device.pending');
    fireEvent.change(screen.getByTestId('select-fatek-mode'), { target: { value: 'serial' } });
    expect(screen.getByTestId('recording-source-discard')).toBeVisible();
    fireEvent.click(screen.getByTestId('recording-source-discard'));
    expect(screen.getByTestId('input-rtu-port')).toHaveValue('');
    expect(screen.queryByTestId('recording-source-reuse')).not.toBeInTheDocument();
    expect(current.recentConnections.read).toHaveBeenLastCalledWith('workspace-1', 'fatek_fbs', 'serial');
    expect(screen.getByText('device.pending')).toBeVisible();
    expect(screen.getByTestId('recording-source-save')).toBeDisabled();
    expect(current.port.submit).toHaveBeenCalledOnce();
  });
});


describe('an acknowledgement does not decide an already open connection guard', () => {
  it.each([{ action: 'clear', active: true }, { action: 'clear', active: false }, { action: 'subtype', active: true }, { action: 'subtype', active: false }] as const)('retains pending $action through ACK (active=$active) and changes only the local buffer', async ({ action, active }) => {
    const current = setup(); await screen.findByTestId('recording-source-reuse');
    let acknowledge!: (outcome: DeviceSourceOutcome) => void;
    vi.mocked(current.port.submit).mockImplementation(() => new Promise(resolve => { acknowledge = resolve; }));
    fireEvent.change(screen.getByTestId('recording-source-name'), { target: { value: 'Saved PLC' } });
    fireEvent.change(screen.getByTestId('input-host'), { target: { value: '192.0.2.40' } });
    fireEvent.click(screen.getByTestId('recording-source-save'));
    await screen.findByText('device.pending');
    if (action === 'clear') fireEvent.click(screen.getByTestId('recording-source-clear-recent'));
    else fireEvent.change(screen.getByTestId('select-fatek-mode'), { target: { value: 'serial' } });
    expect(screen.getByTestId('recording-source-discard')).toBeVisible();
    const request = vi.mocked(current.port.submit).mock.calls[0][0];
    const saved = { id: 'saved-plc', name: request.buffer.name, protocol: request.buffer.protocol,
      connection_config: JSON.stringify(request.buffer.config) } as StudioV2WorkspaceDeviceRecord;
    current.slot.active = active; // Hidden modal recovery also must preserve an already chosen guard intent.
    await act(async () => { acknowledge({ request, status: 'confirmed', device: saved }); });
    expect(current.slot.onSaved).not.toHaveBeenCalled();
    expect(screen.getByTestId('recording-source-discard')).toBeVisible();
    fireEvent.click(screen.getByTestId('recording-source-discard'));
    if (action === 'clear') {
      expect(current.recentConnections.clear).toHaveBeenCalledWith('workspace-1', 'fatek_fbs');
      expect(screen.getByTestId('input-host')).toHaveValue('192.168.1.100');
    } else {
      expect(screen.getByTestId('select-fatek-mode')).toHaveValue('serial');
      expect(screen.getByTestId('input-rtu-port')).toHaveValue('');
    }
    expect(screen.getByTestId('recording-source-name')).toHaveValue('Saved PLC');
    expect(screen.getByTestId('recording-source-save')).toHaveTextContent('device.saveChanges');
    expect(JSON.parse(saved.connection_config).host).toBe('192.0.2.40');
    expect(current.port.submit).toHaveBeenCalledOnce();
    expect(current.slot.onSaved).not.toHaveBeenCalled();
  });
});

describe('cancelling a clear after acknowledgement', () => {
  it('leaves the connection buffer and recent cache unchanged', async () => {
    const current = setup(); await screen.findByTestId('recording-source-reuse');
    let acknowledge!: (outcome: DeviceSourceOutcome) => void;
    vi.mocked(current.port.submit).mockImplementation(() => new Promise(resolve => { acknowledge = resolve; }));
    fireEvent.change(screen.getByTestId('recording-source-name'), { target: { value: 'Keep PLC' } });
    fireEvent.change(screen.getByTestId('input-host'), { target: { value: '192.0.2.41' } });
    fireEvent.click(screen.getByTestId('recording-source-save'));
    await screen.findByText('device.pending');
    fireEvent.click(screen.getByTestId('recording-source-clear-recent'));
    const request = vi.mocked(current.port.submit).mock.calls[0][0];
    const saved = { id: 'saved-plc', name: request.buffer.name, protocol: request.buffer.protocol,
      connection_config: JSON.stringify(request.buffer.config) } as StudioV2WorkspaceDeviceRecord;
    await act(async () => { acknowledge({ request, status: 'confirmed', device: saved }); });
    expect(current.slot.onSaved).not.toHaveBeenCalled();
    fireEvent.click(screen.getByText('device.keep'));
    expect(screen.getByTestId('input-host')).toHaveValue('192.0.2.41');
    expect(screen.getByTestId('recording-source-name')).toHaveValue('Keep PLC');
    expect(current.recentConnections.clear).not.toHaveBeenCalled();
    expect(current.port.submit).toHaveBeenCalledOnce();
  });
});
