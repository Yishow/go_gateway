import { beforeEach, describe, expect, it, vi } from 'vitest';
import { parseStudioV2WriteGroup, parseWriteGroupSuccessorCandidate } from '@/utils/studioV2WriteGroupDefinitionJson';
import { parseRecordingProfile, recordingProfileIntervalSeconds } from '@/utils/studioV2RecordingProfile';
import { studioV2WorkspaceWriteGroupsAPI } from '@/services/studioV2WorkspaceWriteGroups';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import type { WriteGroupCreateRequest } from '@/types/studioV2WriteGroup';
vi.mock('@/services/studioV2Workspace',()=>({studioV2DatalinkApi:{get:vi.fn(),post:vi.fn(),put:vi.fn(),delete:vi.fn()}}));
function mixed() {return {
 id:'group',workspace_id:'workspace',revision:'revision',applied_revision:'',name:'mixed',status:'draft',
 record_id:'11111111-1111-4111-8111-111111111111',
 layout_version:'managed_wide_profile_v1',recording_profile:{kind:'mixed',result_shape:'member-calculation-v1'},
 members:[{device_id:'device',point_id:'point',tag_id:'tag',source_revision:'source',mapping_revision:'mapping',target_column:'temperature',required:true,value_role:'gauge',source_value_type:'float64',result_value_type:'float64',calculation:'mean',cadence:{interval:15,unit:'second',alignment:'utc_epoch'}}],
 destination:{connector_id:'connector',connector_revision:'connector-rev',database:'db',table_schema:'',table_name:'mixed_values',storage_strategy:'managed'},
 row_policy:{allowed_lateness_seconds:0,incomplete_policy:'partial',record_key_column:'record_id',provenance_column:'provenance',due_at_column:'due_at'},write_policy:{mode:'append'},migration:{},created_at:'2026-10-31T00:00:00Z',updated_at:'2026-10-31T00:00:00Z',
};}
describe('canonical mixed recording contract',()=>{
 beforeEach(()=>vi.clearAllMocks());
 it('retains member intent, identity and frozen types without a group interval',()=>{
  const input=mixed();const result=parseStudioV2WriteGroup(input);
  expect(result).toEqual(input);
  expect(result?.row_policy).not.toHaveProperty('interval_seconds');
  expect(parseRecordingProfile(input.recording_profile)).toEqual(input.recording_profile);
  if(result?.recording_profile) expect(recordingProfileIntervalSeconds(result.recording_profile)).toBeNull();
 });
 it.each([5,15,60,600])('preserves cadence %s rather than applying legacy bounds',interval=>{
  const input=mixed();input.members[0].cadence.interval=interval;
  expect(parseStudioV2WriteGroup(input)?.members[0]).toEqual(input.members[0]);
 });
 it.each(['record_id','calculation','cadence','source_value_type','result_value_type','source_revision','mapping_revision'])('rejects absent %s',key=>{
  const input=mixed();if(key==='record_id')Reflect.deleteProperty(input,key);else Reflect.deleteProperty(input.members[0],key);
  expect(parseStudioV2WriteGroup(input)).toBeNull();
 });
 it.each([0,60,null])('rejects mixed group interval %s',interval_seconds=>{
  const input=mixed();Object.assign(input.row_policy,{interval_seconds});expect(parseStudioV2WriteGroup(input)).toBeNull();
 });
 it.each([{interval:0},{interval:1.5},{interval:9223372037},{unit:'hour'},{alignment:'activation'}])('rejects invalid cadence %j',patch=>{
  const input=mixed();Object.assign(input.members[0].cadence,patch);expect(parseStudioV2WriteGroup(input)).toBeNull();
 });
 it('clones create without dropping or aliasing canonical member facts',async()=>{
  const input=mixed();vi.mocked(studioV2DatalinkApi.post).mockResolvedValue({data:{success:true,data:{workspace_revision:'new',group:input}}});
  const request={workspace_id:'workspace',expected_workspace_revision:'old',expected_connector_revision:'connector-rev',group:input} as unknown as WriteGroupCreateRequest;
  await studioV2WorkspaceWriteGroupsAPI.create(request);
  const sent=vi.mocked(studioV2DatalinkApi.post).mock.calls[0][1] as WriteGroupCreateRequest;
  expect(sent.group.record_id).toBe(input.record_id);
  expect(sent.group.members[0]).toEqual(input.members[0]);
  expect(sent.group.members[0].cadence).not.toBe(input.members[0].cadence);
  expect(sent.group.row_policy).not.toHaveProperty('interval_seconds');
  expect(sent.group.row_policy).toHaveProperty('due_at_column','due_at');
 });
 it.each(['latest','mean','delta'])('round-trips explicit %s without changing source meaning',calculation=>{
  const input=mixed();input.members[0].calculation=calculation;
  if(calculation==='delta')Object.assign(input.members[0],{measurement_id:'measurement',value_role:'counter_raw'});
  expect(parseStudioV2WriteGroup(input)?.members[0]).toEqual(input.members[0]);
 });
 it('retains scaled mixed latest advanced readback without claiming runtime support',()=>{
  const input=mixed();input.layout_version='managed_long_v1';Object.assign(input,{storage_mode:'plain'});
  Object.assign(input.row_policy,{bucket_start_column:'bucket_start'});Reflect.deleteProperty(input.row_policy,'due_at_column');
  const proof={measurement_id:'measurement',definition_revision:'definition',source_binding_revision:'binding',series_epoch:'epoch',original_data_type:'uint64',scale:'0.1',unit:'kWh',semantic_kind:'counter',counter_policy:{allow_negative:false}};
  Object.assign(input.members[0],{calculation:'latest',measurement_id:'measurement',value_role:'counter_raw',source_value_type:'uint64',result_value_type:'decimal',exact_recording_source:proof});
  expect(parseStudioV2WriteGroup(input)?.members[0]).toEqual(input.members[0]);
 });
 it.each(['future','',null,{toString:()=> 'latest'}])('rejects unknown method %j',calculation=>{
  const input=mixed();Object.assign(input.members[0],{calculation});expect(parseStudioV2WriteGroup(input)).toBeNull();
 });
 it.each(['00000000-0000-0000-0000-000000000000','wrong','11111111-1111-4111-8111-11111111111A'])('rejects noncanonical record identity %s',record_id=>{
  const input=mixed();input.record_id=record_id;expect(parseStudioV2WriteGroup(input)).toBeNull();
 });
 it('keeps minute cadence intact and fails unsafe values',()=>{
  const input=mixed();input.members[0].cadence={interval:10,unit:'minute',alignment:'utc_epoch'};
  expect(parseStudioV2WriteGroup(input)?.members[0].cadence).toEqual(input.members[0].cadence);
  input.members[0].cadence.interval=153722868;expect(parseStudioV2WriteGroup(input)).toBeNull();
 });
 it('does not relabel legacy profiles with mixed result shape',()=>{
  expect(parseRecordingProfile({kind:'instantaneous',interval:15,unit:'second',result_shape:'member-calculation-v1'})).toBeNull();
  expect(parseRecordingProfile({kind:'mixed',result_shape:'unknown'})).toBeNull();
  expect(parseRecordingProfile({kind:'mixed',result_shape:'member-calculation-v1',sample_reduction:'mean'})).toBeNull();
 });
 it('keeps independent records on one source independent during update',async()=>{
  const first=mixed(),second=mixed();second.record_id='22222222-2222-4222-8222-222222222222';second.members[0].calculation='latest';
  vi.mocked(studioV2DatalinkApi.put).mockResolvedValue({data:{success:true,data:{workspace_revision:'next',group:second}}});
  await studioV2WorkspaceWriteGroupsAPI.update('group',{workspace_id:'workspace',expected_workspace_revision:'old',expected_connector_revision:'connector-rev',expected_group_revision:'revision',group:second} as unknown as Parameters<typeof studioV2WorkspaceWriteGroupsAPI.update>[1]);
  const sent=vi.mocked(studioV2DatalinkApi.put).mock.calls[0][1] as WriteGroupCreateRequest;
  expect(sent.group.record_id).toBe(second.record_id);expect(sent.group.members[0].calculation).toBe('latest');
  expect(first.members[0].calculation).toBe('mean');expect(sent.group.members[0].cadence).not.toBe(second.members[0].cadence);
 });

 it('does not authorize a legacy successor route for mixed records',()=>{
  const input=mixed();Object.assign(input,{id:'',revision:'',applied_revision:'',created_at:'0001-01-01T00:00:00Z',updated_at:'0001-01-01T00:00:00Z'});
  input.destination.table_name='';expect(parseWriteGroupSuccessorCandidate(input)).toBeNull();
 });

});
