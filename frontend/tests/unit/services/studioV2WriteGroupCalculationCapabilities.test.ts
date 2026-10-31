import { expect, it } from 'vitest';
import { parseWriteGroupListData, parseWriteGroupMutationData } from '@/utils/studioV2WriteGroupJson';
import { savedGroup } from '../../fixtures/writeGroupState';
function capabilities() { return {
 version:1,result_shape:'member-calculation-v1',draft_supported:true,activation_supported:false,editor_supported:false,reason:'production_acceptance_pending',methods:['latest','mean','delta'],
 cadence:{alignment:'utc_epoch',units:['second','minute'],minimum:1,max_seconds:9223372036,max_minutes:153722867,window_bounds_checked:true},
 schema_dialects:['sqlite','postgres'],schema_layouts:['managed_wide_profile_v1','managed_long_v1'],
 unsupported_combinations:[
 {method:'latest',source_transform:'exact_counter_scale',phase:'draft',reason:'mixed_exact_source_draft_unsupported'},
 {method:'delta',source_transform:'exact_counter_scale',phase:'draft',reason:'mixed_exact_source_draft_unsupported'},
 {method:'latest',source_transform:'exact_counter_scale',phase:'runtime',reason:'mixed_exact_source_latest_unsupported'},
 ],
}; }
function list(cap:unknown=undefined) {return {workspace_id:'ws-1',workspace_revision:'wr',groups:[savedGroup()],...(cap===undefined?{}:{calculation_capabilities:cap})};}
it('retains detached backend preparation capabilities in both list and mutation, never promotes readiness',()=>{
 const original=capabilities();const result=parseWriteGroupListData(list(original));
 expect(result?.calculation_capabilities).toEqual(original);
 const mutation=parseWriteGroupMutationData({workspace_revision:'wr',group:savedGroup(),calculation_capabilities:original});
 expect(mutation?.calculation_capabilities).toEqual(original);
 expect(result?.calculation_capabilities?.draft_supported).toBe(true);
 expect(result?.calculation_capabilities?.activation_supported).toBe(false);
 expect(result?.calculation_capabilities?.editor_supported).toBe(false);
 original.methods[0]='changed';original.cadence.units[0]='changed';original.unsupported_combinations[0].reason='changed';
 original.schema_dialects[0]='mysql';original.schema_layouts[0]='managed_wide_v1';
 expect(result?.calculation_capabilities?.methods[0]).toBe('latest');
 expect(result?.calculation_capabilities?.cadence.units[0]).toBe('second');
 expect(result?.calculation_capabilities?.schema_dialects[0]).toBe('sqlite');
 expect(result?.calculation_capabilities?.schema_layouts[0]).toBe('managed_wide_profile_v1');
 expect(mutation?.calculation_capabilities?.unsupported_combinations[0].reason).toBe('mixed_exact_source_draft_unsupported');
});
it('does not invent capabilities for an older backend',()=>{
 expect(parseWriteGroupListData(list())).not.toHaveProperty('calculation_capabilities');
 expect(parseWriteGroupMutationData({workspace_revision:'wr',group:savedGroup()})).not.toHaveProperty('calculation_capabilities');
});
it.each([null,{},false,{...capabilities(),version:2},{...capabilities(),result_shape:'raw'},
 {...capabilities(),activation_supported:'true'},{...capabilities(),methods:['mean','unknown']},
 {...capabilities(),schema_dialects:['mysql']},{...capabilities(),schema_layouts:['managed_wide_v1']},
 {...capabilities(),unsupported_combinations:[{method:'latest',phase:'future',source_transform:'exact_counter_scale',reason:'blocked'}]},
 {...capabilities(),cadence:{...capabilities().cadence,alignment:'activation'}},
 {...capabilities(),cadence:{...capabilities().cadence,max_seconds:Number.MAX_SAFE_INTEGER+1}},
 {...capabilities(),cadence:{...capabilities().cadence,window_bounds_checked:false}},
])('fails closed for malformed present capabilities %#',cap=>{
 expect(parseWriteGroupListData(list(cap))).toBeNull();
 expect(parseWriteGroupMutationData({workspace_revision:'wr',group:savedGroup(),calculation_capabilities:cap})).toBeNull();
});

it.each(['methods','reason','cadence','unsupported_combinations','activation_supported','editor_supported'])('rejects a missing required %s rather than filling defaults',key=>{
 const cap=capabilities();Reflect.deleteProperty(cap,key);expect(parseWriteGroupListData(list(cap))).toBeNull();
});
it.each([
 {...capabilities(),methods:['mean','mean']},
 {...capabilities(),cadence:{...capabilities().cadence,units:['second','second']}},
 {...capabilities(),cadence:{...capabilities().cadence,minimum:0}},
 {...capabilities(),cadence:{...capabilities().cadence,max_minutes:153722868}},
 {...capabilities(),cadence:{...capabilities().cadence,max_seconds:1.5}},
 {...capabilities(),new_permission:true},
 {...capabilities(),unsupported_combinations:[capabilities().unsupported_combinations[0],capabilities().unsupported_combinations[0]]},
])('rejects ambiguous or unsafe capability material %#',cap=>{expect(parseWriteGroupListData(list(cap))).toBeNull();});
