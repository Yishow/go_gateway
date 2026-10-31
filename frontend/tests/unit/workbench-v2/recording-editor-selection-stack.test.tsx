import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { useState } from 'react';
import { RecordingPointSelection } from '@/features/datalink/workbench-v2/recording-editor/RecordingPointSelection';
import type { EditorMember } from '@/features/datalink/workbench-v2/state/writeGroup/draft';
import type { GroupCandidate } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string, values?: { count?: number }) => values?.count !== undefined ? `${key}:${values.count}` : key }) }));
const rows: GroupCandidate[] = [1, 2].map(n => ({ key: `P${n}|T${n}`, device_id: `D${n}`, point_id: `P${n}`, tag_id: `T${n}`, tag_key: `tag${n}`, label: `Point ${n}`, device_name: `PLC ${n}`, address: '40001', target_type: 'float32' }));
function Harness({ port, disabled = false }: { disabled?: boolean; port?: import('@/pages/datalink/workbench-v2/deviceSourceSubmission').DeviceSourcePort } = {}) {
  const [members, setMembers] = useState<EditorMember[]>([]);
  return <><input aria-label="Recording name" defaultValue="Unsaved name" />
    <output data-testid="members">{JSON.stringify(members)}</output>
    <RecordingPointSelection deviceSource={port} scope="S" workspaceId="W" candidates={rows} members={members} onChange={setMembers} disabled={disabled}
      devices={[{ id: 'D1', name: 'PLC 1' }, { id: 'D2', name: 'PLC 2' }, { id: 'EMPTY', name: 'Empty PLC' }]} />
  </>;
}
describe('single recording picker modal stack', () => {
  it('keeps pending distinct from recording membership and Cancel preserves parent draft', () => {
    render(<Harness />); fireEvent.click(screen.getByTestId('recording-points-open'));
    expect(screen.getByRole('button', { name: 'picker.add:0' })).toBeDisabled();
    fireEvent.click(screen.getByRole('checkbox', { name: 'Point 1' }));
    expect(screen.getByTestId('members')).toHaveTextContent('[]');
    fireEvent.change(screen.getByLabelText('picker.device'), { target: { value: 'D2' } });
    expect(screen.queryByRole('checkbox', { name: 'Point 1' })).not.toBeInTheDocument();
    expect(screen.getByText('picker.hidden:1')).toBeVisible();
    fireEvent.click(screen.getByRole('button', { name: 'picker.cancel' }));
    expect(screen.getByLabelText('Recording name')).toHaveValue('Unsaved name');
    expect(screen.getByTestId('members')).toHaveTextContent('[]');
  });
  it('enters empty-device range buffer in one dialog and dirty cancel restores picker context', () => {
    render(<Harness />); fireEvent.click(screen.getByTestId('recording-points-open'));
    fireEvent.change(screen.getByLabelText('picker.device'), { target: { value: 'EMPTY' } });
    fireEvent.click(screen.getByRole('button', { name: 'picker.newSource' }));
    expect(screen.getAllByRole('dialog')).toHaveLength(1);
    expect(screen.getByTestId('recording-range-device')).toHaveTextContent('Empty PLC');
    fireEvent.change(screen.getByLabelText('range.count'), { target: { value: '1x' } });
    fireEvent.click(screen.getByRole('button', { name: 'source.back' }));
    expect(screen.getByRole('alertdialog')).toBeVisible();
    fireEvent.click(screen.getByRole('button', { name: 'device.keep' }));
    expect(screen.getByLabelText('range.count')).toHaveValue('1x');
    fireEvent.click(screen.getByRole('button', { name: 'source.back' }));
    fireEvent.click(screen.getByRole('button', { name: 'device.discard' }));
    expect(screen.getByLabelText('picker.device')).toHaveValue('EMPTY');
    expect(screen.getByTestId('members')).toHaveTextContent('[]');
  });
  it('all devices opens explicit choice without selecting the first device', () => {
    render(<Harness />); fireEvent.click(screen.getByTestId('recording-points-open'));
    fireEvent.click(screen.getByRole('button', { name: 'picker.newSource' }));
    expect(screen.getByLabelText('source.chooseDevice')).toHaveValue('');
    expect(screen.getByRole('button', { name: 'source.continue' })).toBeDisabled();
    expect(screen.queryByTestId('recording-range-device')).not.toBeInTheDocument();
  });
  it('saves device through the existing port in the same dialog then cancels range without removing the source', async () => {
    const submit = vi.fn(async (request: import('@/pages/datalink/workbench-v2/deviceSourceSubmission').DeviceSourceRequest) => ({ request, status: 'confirmed' as const,
      device: { id: 'PERSISTED', name: request.buffer.name, description: '', protocol: request.buffer.protocol, status: 'draft' as const, connection_config: '{}', last_test_at: null, last_test_success: null, last_test_error: '', created_at: '', updated_at: '' } }));
    const port = { submit, inspect: vi.fn(), observations: () => [] };
    render(<Harness port={port} />); fireEvent.click(screen.getByTestId('recording-points-open'));
    fireEvent.click(screen.getByRole('button', { name: 'picker.newSource' }));
    fireEvent.click(screen.getByRole('button', { name: 'source.newDevice' }));
    fireEvent.change(await screen.findByTestId('recording-source-name'), { target: { value: 'Saved PLC' } });
    fireEvent.change(screen.getByTestId('input-host'), { target: { value: '192.0.2.50' } });
    fireEvent.click(screen.getByTestId('recording-source-save'));
    await waitFor(() => expect(screen.getByTestId('recording-range-device')).toHaveTextContent('Saved PLC'));
    expect(screen.getAllByRole('dialog')).toHaveLength(1);
    fireEvent.click(screen.getByRole('button', { name: 'source.back' }));
    expect(screen.getByRole('option', { name: 'Saved PLC' })).toHaveValue('PERSISTED');
    fireEvent.click(screen.getByRole('button', { name: 'picker.cancel' }));
    expect(submit).toHaveBeenCalledTimes(1);
    expect(screen.getByTestId('members')).toHaveTextContent('[]');
  });

  it('keeps unknown submitted device and newer raw input across picker close/reopen without a replacement create', async () => {
    const submit = vi.fn(async (request: import('@/pages/datalink/workbench-v2/deviceSourceSubmission').DeviceSourceRequest) => ({ request, status: 'unknown' as const }));
    render(<Harness port={{ submit, inspect: vi.fn(), observations: () => [] }} />);
    const openDevice = async () => {
      fireEvent.click(screen.getByTestId('recording-points-open'));
      fireEvent.click(screen.getByRole('button', { name: 'picker.newSource' }));
      fireEvent.click(screen.getByRole('button', { name: 'source.newDevice' }));
      await screen.findByTestId('recording-source-name');
    };
    await openDevice();
    fireEvent.change(screen.getByTestId('recording-source-name'), { target: { value: 'Uncertain PLC' } });
    fireEvent.change(screen.getByTestId('input-host'), { target: { value: '192.0.2.50' } });
    fireEvent.click(screen.getByTestId('recording-source-save'));
    await screen.findByTestId('recording-source-unknown');
    fireEvent.change(screen.getByTestId('input-port'), { target: { value: 'invalid 5x' } });
    fireEvent.click(screen.getByTestId('recording-source-close'));
    fireEvent.click(screen.getByRole('button', { name: 'picker.cancel' }));
    await openDevice();
    expect(screen.getByTestId('input-port')).toHaveValue('invalid 5x');
    expect(screen.getByLabelText('device.protocol')).toBeDisabled();
    expect(screen.getByTestId('recording-source-save')).toBeDisabled();
    expect(submit).toHaveBeenCalledTimes(1);
  });

  it('late lease preparation cannot pull the operator out of the restored picker', async () => {
    let resolve!: (value: boolean) => void;
    const prepare = vi.fn(() => new Promise<boolean>(done => { resolve = done; }));
    render(<Harness port={{ prepare, submit: vi.fn(), inspect: vi.fn(), observations: () => [] }} />);
    fireEvent.click(screen.getByTestId('recording-points-open'));
    fireEvent.click(screen.getByRole('button', { name: 'picker.newSource' }));
    fireEvent.click(screen.getByRole('button', { name: 'source.newDevice' }));
    fireEvent.click(screen.getByRole('button', { name: 'source.back' }));
    await import('@testing-library/react').then(({ act }) => act(async () => resolve(true)));
    expect(screen.getByLabelText('picker.search')).toBeVisible();
    expect(screen.queryByTestId('recording-source-name')).not.toBeInTheDocument();
  });

  it('preserves hidden pending across an empty-device range round trip and adds only on explicit Add', () => {
    render(<Harness />); fireEvent.click(screen.getByTestId('recording-points-open'));
    fireEvent.click(screen.getByRole('checkbox', { name: 'Point 1' }));
    fireEvent.click(screen.getByRole('checkbox', { name: 'Point 2' }));
    fireEvent.change(screen.getByLabelText('picker.device'), { target: { value: 'EMPTY' } });
    fireEvent.click(screen.getByRole('button', { name: 'picker.newSource' }));
    fireEvent.click(screen.getByRole('button', { name: 'source.back' }));
    expect(screen.getByText('picker.hidden:2')).toBeVisible();
    expect(screen.getByTestId('members')).toHaveTextContent('[]');
    fireEvent.click(screen.getByRole('button', { name: 'picker.add:2' }));
    expect(JSON.parse(screen.getByTestId('members').textContent!)).toHaveLength(2);
    expect(screen.getByLabelText('Recording name')).toHaveValue('Unsaved name');
  });

  it('disables visible source Save immediately when the live recording permission is suspended', async () => {
    const port = { submit: vi.fn(), inspect: vi.fn(), observations: () => [] };
    const view = render(<Harness port={port} />);
    fireEvent.click(screen.getByTestId('recording-points-open'));
    fireEvent.click(screen.getByRole('button', { name: 'picker.newSource' }));
    fireEvent.click(screen.getByRole('button', { name: 'source.newDevice' }));
    fireEvent.change(await screen.findByTestId('recording-source-name'), { target: { value: 'Unsent' } });
    fireEvent.change(screen.getByTestId('input-host'), { target: { value: '192.0.2.50' } });
    view.rerender(<Harness port={port} disabled />);
    expect(screen.getByTestId('recording-source-save')).toBeDisabled();
    expect(port.submit).not.toHaveBeenCalled();
  });

});

it('keeps a same-owner submitted unknown device reachable for inspection when point discovery fails on reopen', async()=>{
 const discover=vi.fn().mockResolvedValue({candidates:[],devices:[]});
 const submit=vi.fn(async(request: import('@/pages/datalink/workbench-v2/deviceSourceSubmission').DeviceSourceRequest)=>({request,status:'unknown' as const}));
 const inspect = vi.fn(async () => ({ request: submit.mock.calls[0][0], status: 'unknown' as const }));
 render(<RecordingPointSelection scope="S" workspaceId="W" candidates={[]} members={[]} devices={[]} disabled={false} onChange={vi.fn()} discover={discover} deviceSource={{submit,inspect,observations:()=>[]}}/>);
 fireEvent.click(screen.getByTestId('recording-points-open'));
 fireEvent.click(await screen.findByRole('button',{name:'picker.newSource'}));
 fireEvent.click(screen.getByRole('button',{name:'source.newDevice'}));
 fireEvent.change(await screen.findByTestId('recording-source-name'),{target:{value:'Unknown PLC'}});
 fireEvent.change(screen.getByTestId('input-host'),{target:{value:'192.0.2.5'}});
 fireEvent.click(screen.getByTestId('recording-source-save'));
 await screen.findByTestId('recording-source-unknown');
 fireEvent.click(screen.getByTestId('recording-source-close'));
 fireEvent.click(screen.getByRole('button',{name:'picker.cancel'}));
 discover.mockRejectedValue(new Error('point service unavailable'));
 fireEvent.click(screen.getByTestId('recording-points-open'));
 await screen.findByText('picker.failed');
 // Existing operation should remain reachable even when unrelated point discovery fails.
 expect(screen.getByRole('button',{name:'picker.add:0'})).toBeDisabled();
 fireEvent.click(screen.getByTestId('recording-source-resume'));
 fireEvent.click(screen.getByTestId('recording-source-inspect'));
 expect(inspect).toHaveBeenCalledWith(submit.mock.calls[0][0].requestId);
 expect(screen.getAllByRole('dialog')).toHaveLength(1);
 expect(submit).toHaveBeenCalledTimes(1);
});

it('does not resurrect a confirmed device omitted by a later successful authoritative discovery', async()=>{
 const discover=vi.fn().mockResolvedValue({candidates:[],devices:[]});
 const submit=vi.fn(async(request: import('@/pages/datalink/workbench-v2/deviceSourceSubmission').DeviceSourceRequest)=>({request,status:'confirmed' as const,device:{id:'SAVED',name:'Saved PLC',protocol:request.buffer.protocol,description:'',status:'draft' as const,connection_config:'{}',created_at:'',updated_at:'',last_test_at:null,last_test_success:null,last_test_error:''}}));
 render(<RecordingPointSelection scope="S" workspaceId="W" candidates={[]} members={[]} devices={[]} disabled={false} onChange={vi.fn()} discover={discover} deviceSource={{submit,inspect:vi.fn(),observations:()=>[]}}/>);
 fireEvent.click(screen.getByTestId('recording-points-open'));
 fireEvent.click(await screen.findByRole('button',{name:'picker.newSource'}));
 fireEvent.click(screen.getByRole('button',{name:'source.newDevice'}));
 fireEvent.change(await screen.findByTestId('recording-source-name'),{target:{value:'Saved PLC'}});
 fireEvent.change(screen.getByTestId('input-host'),{target:{value:'192.0.2.5'}});
 fireEvent.click(screen.getByTestId('recording-source-save'));
 await screen.findByTestId('recording-range-device');
 fireEvent.click(screen.getByTestId('recording-range-back'));
 expect(screen.getByRole('option',{name:'Saved PLC'})).toBeInTheDocument();
 fireEvent.click(screen.getByRole('button',{name:'picker.cancel'}));
 // An authoritative successful later collection says this device no longer exists.
 fireEvent.click(screen.getByTestId('recording-points-open'));
 await screen.findByRole('button',{name:'picker.newSource'});
 expect(discover).toHaveBeenCalledTimes(2);
 expect(screen.queryByRole('option',{name:'Saved PLC'})).not.toBeInTheDocument();
});

it('a pre-confirmation read that arrives late cannot remove the just-saved device, but a later-started read can', async () => {
 const empty = { candidates: [], devices: [] };
 let resolveOld!: (value: typeof empty) => void;
 const discover = vi.fn().mockResolvedValueOnce(empty).mockImplementationOnce(() => new Promise(done => { resolveOld = done; })).mockResolvedValue(empty);
 const submit = vi.fn(async (request: import('@/pages/datalink/workbench-v2/deviceSourceSubmission').DeviceSourceRequest) => ({ request, status: 'unknown' as const }));
 const inspect = vi.fn(async () => ({ request: submit.mock.calls[0][0], status: 'confirmed' as const, device: { id: 'SAVED', name: 'Saved PLC', protocol: 'modbus_tcp' as const,
   description: '', status: 'draft' as const, connection_config: '{}', created_at: '', updated_at: '', last_test_at: null, last_test_success: null, last_test_error: '' } }));
 render(<RecordingPointSelection scope="S" workspaceId="W" candidates={[]} members={[]} devices={[]} disabled={false} onChange={vi.fn()} discover={discover} deviceSource={{ submit, inspect, observations: () => [] }} />);
 fireEvent.click(screen.getByTestId('recording-points-open'));
 fireEvent.click(await screen.findByRole('button', { name: 'picker.newSource' }));
 fireEvent.click(screen.getByRole('button', { name: 'source.newDevice' }));
 fireEvent.change(await screen.findByTestId('recording-source-name'), { target: { value: 'Saved PLC' } });
 fireEvent.change(screen.getByTestId('input-host'), { target: { value: '192.0.2.5' } });
 fireEvent.click(screen.getByTestId('recording-source-save'));
 await screen.findByTestId('recording-source-unknown');
 fireEvent.click(screen.getByTestId('recording-source-close'));
 fireEvent.click(screen.getByRole('button', { name: 'picker.cancel' }));
 fireEvent.click(screen.getByTestId('recording-points-open')); // Read 2 starts before confirmed save.
 expect(screen.getByText('picker.loading')).toBeVisible();
 fireEvent.click(screen.getByTestId('recording-source-resume'));
 fireEvent.click(screen.getByTestId('recording-source-inspect'));
 await screen.findByTestId('recording-range-device');
 await act(async () => resolveOld(empty)); // Its old empty snapshot must not remove this new ID.
 expect(screen.getByTestId('recording-source-confirmed')).toBeVisible();
 fireEvent.click(screen.getByTestId('recording-range-back'));
 expect(screen.getByRole('option', { name: 'Saved PLC' })).toBeInTheDocument();
 fireEvent.click(screen.getByRole('button', { name: 'picker.cancel' }));
 fireEvent.click(screen.getByTestId('recording-points-open')); // Read 3 starts after confirmed save.
 await screen.findByRole('button', { name: 'picker.newSource' });
 expect(screen.queryByRole('option', { name: 'Saved PLC' })).not.toBeInTheDocument();
 expect(submit).toHaveBeenCalledTimes(1);
});
