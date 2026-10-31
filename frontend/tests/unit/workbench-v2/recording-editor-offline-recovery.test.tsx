import type { StudioV2WorkspaceDeviceRecord } from '@/services/studioV2WorkspaceDevices';
import type { SourceRuleRecord } from '@/types/datalink';
import {act,renderHook} from '@testing-library/react';
import {it,expect,vi} from 'vitest';
import {useRecordingRangeBuffer} from '@/features/datalink/workbench-v2/recording-editor/useRecordingRangeBuffer';
import type {RangeSourceOutcome,RangeSourceRequest} from '@/pages/datalink/workbench-v2/rangeSourceSubmission';
it('retains the current saved-inactive B receipt on resume when older saved-inactive A also exists',async()=>{
 const observations:RangeSourceOutcome[]=[];
 const port={prepare:vi.fn(async(id:string)=>({device:{id,name:id,protocol:'modbus_tcp',connection_config:'{}'} as StudioV2WorkspaceDeviceRecord,activationAllowed:false})),
 submit:vi.fn(async(request:RangeSourceRequest)=>{const result={request,status:'saved-inactive' as const,rule:{id:`R-${request.deviceId}`} as SourceRuleRecord,savedPoints:[{device_id:request.deviceId,point_id:`P-${request.deviceId}`,tag_id:`T-${request.deviceId}`}]};observations.push(result);return result;}),
 prepareSaved:vi.fn(),inspect:vi.fn(),observations:()=>observations};
 const h=renderHook(()=>useRecordingRangeBuffer('W',port));
 for(const id of ['A','B']){
  await act(async()=>{await h.result.current.start({id,name:id});});
  act(()=>h.result.current.change({name:`Range ${id}`,address:'40001',count:'1',dataType:'int16',scale:'1',offset:'0'}));
  await act(async()=>{await h.result.current.submit();});
 }
 expect(h.result.current.target?.id).toBe('B');
 act(()=>{h.result.current.resume();});
 expect(h.result.current.target?.id).toBe('B');
 expect(h.result.current.outcome?.request.deviceId).toBe('B');
 expect(port.submit).toHaveBeenCalledTimes(2);
 act(()=>{expect(h.result.current.resume(observations[0].request.requestId)).toBe(true);});
 expect(h.result.current.target?.id).toBe('A');
 act(()=>h.result.current.change({...h.result.current.buffer,count:'illegal retained'}));
 act(()=>{expect(h.result.current.resume(observations[1].request.requestId)).toBe(false);});
 expect(h.result.current.buffer.count).toBe('illegal retained');
 expect(h.result.current.target?.id).toBe('A');
 act(()=>{expect(h.result.current.resume(observations[0].request.requestId)).toBe(true);});
 expect(h.result.current.buffer.count).toBe('illegal retained');
});

import {render,screen,fireEvent} from '@testing-library/react';
import {RecordingPointSelection} from '@/features/datalink/workbench-v2/recording-editor/RecordingPointSelection';
vi.mock('react-i18next',()=>({useTranslation:()=>({t:(key:string)=>key})}));
it.each([['A','B'],['A','A']])('native consumer: explicit receipts remain reachable for devices %s and %s',async(first,second)=>{
 const observations:RangeSourceOutcome[]=[];
 const range={prepare:vi.fn(async(id:string)=>({device:{id,name:id,protocol:'modbus_tcp',connection_config:'{}'} as StudioV2WorkspaceDeviceRecord,activationAllowed:false})),
 submit:vi.fn(async(request:RangeSourceRequest)=>{const result={request,status:'saved-inactive' as const,rule:{id:`R-${request.deviceId}`} as SourceRuleRecord,savedPoints:[{device_id:request.deviceId,point_id:`P-${request.deviceId}`,tag_id:`T-${request.deviceId}`}]};observations.push(result);return result;}),prepareSaved:vi.fn(),inspect:vi.fn(),observations:()=>observations};
 render(<RecordingPointSelection scope="S" workspaceId="W" candidates={[]} members={[]} onChange={vi.fn()} disabled={false} devices={[{id:'A',name:'Device A'},{id:'B',name:'Device B'}]} deviceSource={{submit:vi.fn(),inspect:vi.fn(),observations:()=>[],range}}/>);
 fireEvent.click(screen.getByTestId('recording-points-open'));
 for(const id of [first,second]){
  fireEvent.change(screen.getByLabelText('picker.device'),{target:{value:id}});
  fireEvent.click(screen.getByRole('button',{name:'picker.newSource'}));
  await act(async()=>Promise.resolve());
  fireEvent.change(screen.getByLabelText('range.name'),{target:{value:`Range ${id}-${observations.length}`}});
  fireEvent.change(screen.getByLabelText('range.address'),{target:{value:'40001'}});
  fireEvent.change(screen.getByLabelText('range.count'),{target:{value:'1'}});
  fireEvent.click(screen.getByRole('button',{name:'range.save'}));
  await screen.findByTestId('recording-range-saved-inactive');
  fireEvent.click(screen.getByTestId('recording-range-back'));
 }
 fireEvent.click(screen.getByRole('button',{name:'picker.cancel'}));
 fireEvent.click(screen.getByTestId('recording-points-open'));
 fireEvent.click(screen.getByTestId('recording-range-resume'));
 expect(screen.getByLabelText('range.name')).toHaveValue(`Range ${second}-1`);
 expect(range.submit).toHaveBeenCalledTimes(2);
 fireEvent.click(screen.getByTestId('recording-range-back'));
 fireEvent.click(screen.getByTestId(`recording-range-recover-${observations[0].request.requestId}`));
 expect(screen.getByLabelText('range.name')).toHaveValue(`Range ${first}-0`);
 fireEvent.click(screen.getByTestId('recording-range-back'));
 fireEvent.click(screen.getByTestId(`recording-range-recover-${observations[1].request.requestId}`));
 expect(screen.getByLabelText('range.name')).toHaveValue(`Range ${second}-1`);
 expect(range.submit).toHaveBeenCalledTimes(2);
 expect(range.prepareSaved).not.toHaveBeenCalled();
});
