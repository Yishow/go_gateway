import { expect, it } from 'vitest';
import { parseRuntimeRuleEvidence } from '@/services/runtimeRuleEvidenceParser';
import { ruleEvidenceFixture, evidenceScope } from './runtimeRuleEvidenceFixture';
import type { RuleEvidenceView, SourceResult } from '@/types/runtimeRuleEvidence';
const start='2026-10-27T09:59:45Z', end='2026-10-27T10:00:00Z';
function mixed(): RuleEvidenceView {
 const view=ruleEvidenceFixture();
 for(const f of [view.groups[0].saved!,view.groups[0].effective!]) {
  Object.assign(f,{recording_profile:null,interval_seconds:null,reduction:'member_calculation',mixed_recording:{record_id:'11111111-1111-4111-8111-111111111111',result_shape:'member-calculation-v1'}});
  Object.assign(f.members[0],{calculation:{method:'mean',interval:15,unit:'second',alignment:'utc_epoch',required:true,source_revision:'source',mapping_revision:'mapping',measurement_id:null,result_type:'float64'}});
 }
 const receipt={group_revision:'A',connector_revision:'7',record_id:'record',effect_key:'effect',payload_digest:'digest',committed_at:end,provenance:'runtime' as const};
 const result={state:'available',window_start:start,window_end:end,delta_known:null,coverage:'unknown',receipt,calculation:{result_shape:'member-calculation-v1',method:'mean',interval_seconds:15,alignment:'utc_epoch',due_at:end,next_due:null,outcome_digest:'a'.repeat(64),outcome_kind:'row',status:'ok',reason:null,result_present:true,row_emitted:true,raw_present:true,source_type:'float64',result_type:'float64',aligned_start:start,aligned_end:end,actual_start:start,actual_end:end,rounded:false}} as SourceResult;
 view.groups[0].delivery={counter_replay_quarantines:null,oldest_pending_scope:'group-all-revisions',oldest_pending_seconds:null,backlog:[],test_write:null,current:{intake:null,binding:{workspace_id:'W',workspace_revision:'W1',entity_revision:'A'},totals_scope:'group',shared:false,stages:{collecting:0,queued:0,retrying:0,blocked:0,quarantined:0,unknown:0,sql_committed:1,skipped:0},last_commit:receipt,sources:[{device_id:'D',point_id:'P',tag_id:'T',latest:result,last_committed:structuredClone(result)}],source_state:'available',next_due_at:null,next_retry_at:null,attempt_duration_ms:null,quota:null,storage:null}};
 return view;
}
function latest(v:RuleEvidenceView) {return v.groups[0].delivery!.current!.sources[0].latest!;}
it('preserves producer-backed selected results and mixed facts without a common interval',()=>{
 const v=mixed();const result=parseRuntimeRuleEvidence(v,evidenceScope);expect(result.ok).toBe(true);
 if(!result.ok)throw new Error(result.reason);
 expect(result.value.groups[0].state).toBe('effective');expect(result.value.groups[0].effective?.interval_seconds).toBeNull();
 expect(latest(result.value)).toMatchObject({calculation:{result_present:true,raw_present:true,method:'mean'}});
 expect(result.value).toEqual(v);
});
it('accepts wide not-due receipt with no actual window and no selected result',()=>{
 const v=mixed(),r=latest(v);Object.assign(r,{state:'not_due',window_start:null,window_end:null,coverage:null});
 Object.assign(r.calculation!,{status:'not_due',reason:'not_due',result_present:false,raw_present:false,due_at:'2026-10-27T09:59:55Z',next_due:'2026-10-27T10:00:00Z',aligned_start:null,aligned_end:null,actual_start:null,actual_end:null,rounded:null});
 v.groups[0].delivery!.current!.sources[0].last_committed=structuredClone(r);
 expect(parseRuntimeRuleEvidence(v,evidenceScope).ok).toBe(true);
});
it.each(['method','interval_seconds','source_type','result_type'])('rejects result intent mismatch %s',key=>{
 const v=mixed();Object.assign(latest(v).calculation!,{[key]:key==='interval_seconds'?60:key==='method'?'latest':'uint64'});
 expect(parseRuntimeRuleEvidence(v,evidenceScope)).toMatchObject({ok:false,reason:'stale-reference'});
});
it('rejects a fictitious common interval or missing member intent',()=>{
 const v=mixed();v.groups[0].effective!.interval_seconds=15;expect(parseRuntimeRuleEvidence(v,evidenceScope).ok).toBe(false);
 const missing=mixed();Reflect.deleteProperty(missing.groups[0].effective!.members[0],'calculation');expect(parseRuntimeRuleEvidence(missing,evidenceScope).ok).toBe(false);
});
it('rejects unknown result values or unsafe timestamps rather than promoting readiness',()=>{
 for(const patch of [{value:12},{method:'unknown'},{due_at:'2300-01-01T00:00:00Z'},{actual_end:start},{result_present:true,status:'missing'}]) {
  const v=mixed();Object.assign(latest(v).calculation!,patch);expect(parseRuntimeRuleEvidence(v,evidenceScope).ok).toBe(false);
 }
});
it.each(['legacy-profile','wrong-reduction','nil-record','minute-overflow','missing-source-type','delta-without-measurement'])('rejects contradictory mixed facts %s',mode=>{
 const v=mixed(),f=v.groups[0].effective!,m=f.members[0];
 if(mode==='legacy-profile')f.recording_profile={kind:'instantaneous',interval:15,unit:'second'};
 if(mode==='wrong-reduction')f.reduction='mean';
 if(mode==='nil-record')f.mixed_recording!.record_id='00000000-0000-0000-0000-000000000000';
 if(mode==='minute-overflow'){m.calculation!.unit='minute';m.calculation!.interval=153722868;}
 if(mode==='missing-source-type')m.source_value_type=null;
 if(mode==='delta-without-measurement')m.calculation!.method='delta';
 expect(parseRuntimeRuleEvidence(v,evidenceScope)).toMatchObject({ok:false,reason:'unavailable'});
});
it('validates delta-known against selected result, not raw sample presence',()=>{
 const v=mixed(),f=v.groups[0].effective!,r=latest(v);
 f.members[0].calculation!.method='delta';f.members[0].calculation!.measurement_id='measurement';
 r.calculation!.method='delta';r.calculation!.rounded=null;r.delta_known=true;
 v.groups[0].delivery!.current!.sources[0].last_committed=structuredClone(r);
 expect(parseRuntimeRuleEvidence(v,evidenceScope).ok).toBe(true);
 r.delta_known=false;expect(parseRuntimeRuleEvidence(v,evidenceScope).ok).toBe(false);
});
it('keeps unavailable mixed evidence unknown rather than deriving a result',()=>{
 const v=mixed();const c=v.groups[0].delivery!.current!;
 c.stages.sql_committed=0;c.last_commit=null;c.sources[0].last_committed=null;
 c.sources[0].latest={state:'unavailable',window_start:null,window_end:null,delta_known:null,coverage:null,receipt:null};
 const result=parseRuntimeRuleEvidence(v,evidenceScope);expect(result.ok).toBe(true);
 if(result.ok)expect(latest(result.value)).not.toHaveProperty('calculation');
});
it('rejects foreign member identity and mismatched receipt scope',()=>{
 const foreign=mixed();foreign.groups[0].delivery!.current!.sources[0].tag_id='foreign';
 expect(parseRuntimeRuleEvidence(foreign,evidenceScope)).toMatchObject({ok:false,reason:'foreign-scope'});
 const stale=mixed();latest(stale).receipt!.connector_revision='foreign';
 expect(parseRuntimeRuleEvidence(stale,evidenceScope)).toMatchObject({ok:false,reason:'stale-reference'});
});
it('rejects receipt on skipped result and any arbitrary value-bearing field',()=>{
 const skipped=mixed();Object.assign(latest(skipped).calculation!,{outcome_kind:'skipped',row_emitted:false});
 expect(parseRuntimeRuleEvidence(skipped,evidenceScope).ok).toBe(false);
 for(const key of ['value','raw_value','delta','samples']) {
  const v=mixed();Object.assign(latest(v).calculation!,{[key]:'secret'});expect(parseRuntimeRuleEvidence(v,evidenceScope).ok).toBe(false);
 }
});
it('retains nominal slow dispatch and truncated actual coverage with rounding evidence',()=>{
 const v=mixed(),f=v.groups[0].effective!,r=latest(v);
 f.members[0].calculation!.interval=60;
 Object.assign(r,{state:'partial',coverage:'partial',window_start:'2026-10-27T09:59:00Z',window_end:'2026-10-27T09:59:23Z'});
 Object.assign(r.calculation!,{interval_seconds:60,aligned_start:'2026-10-27T09:59:00Z',actual_start:'2026-10-27T09:59:00.000000000Z',actual_end:'2026-10-27T09:59:23Z',rounded:true});
 v.groups[0].delivery!.current!.sources[0].last_committed=structuredClone(r);
 const result=parseRuntimeRuleEvidence(v,evidenceScope);expect(result.ok).toBe(true);
 if(result.ok)expect(latest(result.value).calculation).toMatchObject({due_at:end,actual_end:'2026-10-27T09:59:23Z',rounded:true});
});
it('leaves legacy field absence unchanged',()=>{
 const result=parseRuntimeRuleEvidence(ruleEvidenceFixture(),evidenceScope);expect(result.ok).toBe(true);
 if(!result.ok)throw new Error(result.reason);
 expect(result.value.groups[0].effective).not.toHaveProperty('mixed_recording');
 expect(result.value.groups[0].effective!.members[0]).not.toHaveProperty('calculation');
});
it.each(['aligned-start','aligned-end','next-equal','next-before','long'])('rejects contradictory not-due evidence %s',mode=>{
 const v=mixed(),r=latest(v);Object.assign(r,{state:'not_due',window_start:null,window_end:null,coverage:null});
 Object.assign(r.calculation!,{status:'not_due',reason:'not_due',result_present:false,raw_present:false,due_at:'2026-10-27T09:59:55Z',next_due:'2026-10-27T10:00:00Z',aligned_start:null,aligned_end:null,actual_start:null,actual_end:null,rounded:null});
 if(mode==='aligned-start')r.calculation!.aligned_start=start;
 if(mode==='aligned-end')r.calculation!.aligned_end=end;
 if(mode==='next-equal')r.calculation!.next_due=r.calculation!.due_at;
 if(mode==='next-before')r.calculation!.next_due=start;
 if(mode==='long')v.groups[0].effective!.layout='managed_long_v1';
 v.groups[0].delivery!.current!.sources[0].last_committed=structuredClone(r);
 expect(parseRuntimeRuleEvidence(v,evidenceScope)).toMatchObject({ok:false,reason:'unavailable'});
});

it('keeps ended fast members not-due without fabricating future NextDue',()=>{
 const v=mixed(),r=latest(v);Object.assign(r,{state:'not_due',window_start:null,window_end:null,coverage:null});
 Object.assign(r.calculation!,{status:'not_due',reason:'not_due',result_present:false,raw_present:false,next_due:null,aligned_start:null,aligned_end:null,actual_start:null,actual_end:null,rounded:null});
 v.groups[0].delivery!.current!.sources[0].last_committed=structuredClone(r);
 expect(parseRuntimeRuleEvidence(v,evidenceScope).ok).toBe(true);
});
it.each(['latest','delta'] as const)('rejects rounding evidence for non-mean %s',method=>{
 const v=mixed(),r=latest(v),m=v.groups[0].effective!.members[0];
 m.calculation!.method=method;r.calculation!.method=method;r.calculation!.rounded=false;
 if(method==='delta'){m.calculation!.measurement_id='measurement';r.delta_known=true;}
 v.groups[0].delivery!.current!.sources[0].last_committed=structuredClone(r);
 expect(parseRuntimeRuleEvidence(v,evidenceScope).ok).toBe(false);
});
it.each([start,end])('rejects non-future NextDue even for due rows %s',next=>{
 const v=mixed();latest(v).calculation!.next_due=next;
 expect(parseRuntimeRuleEvidence(v,evidenceScope).ok).toBe(false);
});

it.each(['missing-as-available','not-due-label-on-due','misaligned-next','wrong-window-duration','result-without-raw','wrong-nominal-end','skipped-next-window'])('rejects contradictory producer contract %s', mode => {
 const v=mixed(), r=latest(v),c=r.calculation!;
 if(mode==='result-without-raw')c.raw_present=false;
 if(mode==='missing-as-available')Object.assign(c,{result_present:false,status:'missing',reason:'no-samples'});
 if(mode==='not-due-label-on-due')r.state='not_due';
 if(mode==='misaligned-next')c.next_due='2026-10-27T10:00:01Z';
 if(mode==='wrong-window-duration')c.aligned_start='2026-10-27T09:00:00Z';
 if(mode==='wrong-nominal-end'){c.aligned_start='2026-10-27T09:59:30Z';c.aligned_end=start;c.actual_start=c.aligned_start;c.actual_end=start;r.window_start=c.actual_start;r.window_end=c.actual_end;}
 if(mode==='skipped-next-window')c.next_due='2026-10-27T10:00:30Z';
 expect(parseRuntimeRuleEvidence(v,evidenceScope).ok).toBe(false);
});

it.each(['boundary','skipped','unaligned'])('rejects not-due next window inconsistent with member cadence %s',mode=>{
 const v=mixed(),r=latest(v);Object.assign(r,{state:'not_due',window_start:null,window_end:null,coverage:null});
 Object.assign(r.calculation!,{status:'not_due',reason:'not_due',result_present:false,raw_present:false,due_at:'2026-10-27T09:59:55Z',next_due:end,aligned_start:null,aligned_end:null,actual_start:null,actual_end:null,rounded:null});
 if(mode==='boundary')Object.assign(r.calculation!,{due_at:end,next_due:'2026-10-27T10:00:15Z'});
 if(mode==='skipped')r.calculation!.next_due='2026-10-27T10:00:15Z';
 if(mode==='unaligned')r.calculation!.next_due='2026-10-27T10:00:01Z';
 expect(parseRuntimeRuleEvidence(v,evidenceScope).ok).toBe(false);
});
it('preserves partial first-window nanoseconds without shifting the nominal due',()=>{
 const v=mixed(),r=latest(v);
 Object.assign(r,{state:'partial',coverage:'partial',window_start:'2026-10-27T09:59:50.000000001Z'});
 Object.assign(r.calculation!,{actual_start:r.window_start,next_due:'2026-10-27T10:00:15Z'});
 v.groups[0].delivery!.current!.sources[0].last_committed=structuredClone(r);
 expect(parseRuntimeRuleEvidence(v,evidenceScope).ok).toBe(true);
 r.calculation!.aligned_start='2026-10-27T09:59:45.000000001Z';
 expect(parseRuntimeRuleEvidence(v,evidenceScope).ok).toBe(false);
});
it('preserves unknown first-baseline delta with a raw sample but no selected delta',()=>{
 const v=mixed(),r=latest(v),c=r.calculation!;
 const member=v.groups[0].effective!.members[0].calculation!;member.method='delta';member.measurement_id='counter';
 Object.assign(r,{state:'unknown',delta_known:false,coverage:null,receipt:null});
 Object.assign(c,{method:'delta',rounded:null,status:'missing',reason:'baseline-required',result_present:false,raw_present:true,row_emitted:false,outcome_kind:'skipped',next_due:'2026-10-27T10:00:15Z'});
 v.groups[0].delivery!.current!.sources[0].last_committed=null;
 expect(parseRuntimeRuleEvidence(v,evidenceScope).ok).toBe(true);
});
it('uses mathematical epoch floor for negative-time not-due windows',()=>{
 const v=mixed(),r=latest(v);
 for(const f of [v.groups[0].saved!,v.groups[0].effective!])f.effective_at='1969-12-31T23:59:00Z';Object.assign(r,{state:'not_due',window_start:null,window_end:null,coverage:null});
 Object.assign(r.calculation!,{status:'not_due',reason:'not_due',result_present:false,raw_present:false,due_at:'1969-12-31T23:59:55Z',next_due:'1970-01-01T00:00:00Z',aligned_start:null,aligned_end:null,actual_start:null,actual_end:null,rounded:null});
 v.groups[0].delivery!.current!.sources[0].last_committed=structuredClone(r);
 expect(parseRuntimeRuleEvidence(v,evidenceScope).ok).toBe(true);
 r.calculation!.next_due='1970-01-01T00:00:15Z';
 expect(parseRuntimeRuleEvidence(v,evidenceScope).ok).toBe(false);
});
it('rejects status ok without the selected usable result',()=>{
 const v=mixed(),r=latest(v);r.state='unknown';r.calculation!.result_present=false;
 expect(parseRuntimeRuleEvidence(v,evidenceScope).ok).toBe(false);
});
it('rejects a not-due floor outside the signed UnixNano window domain',()=>{
 const v=mixed(),r=latest(v);Object.assign(r,{state:'not_due',window_start:null,window_end:null,coverage:null});
 Object.assign(r.calculation!,{status:'not_due',reason:'not_due',result_present:false,raw_present:false,due_at:'1677-09-21T00:12:43.145224192Z',next_due:'1677-09-21T00:12:45Z',aligned_start:null,aligned_end:null,actual_start:null,actual_end:null,rounded:null});
 expect(parseRuntimeRuleEvidence(v,evidenceScope).ok).toBe(false);
});
it('rejects invented partial member status while preserving partial coverage separately',()=>{
 const v=mixed(),r=latest(v);
 Object.assign(r,{state:'partial',coverage:'partial'});
 Object.assign(r.calculation!,{status:'partial',result_present:false});
 expect(parseRuntimeRuleEvidence(v,evidenceScope).ok).toBe(false);
});
