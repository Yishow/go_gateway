import { expect, it, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
import { basicMemberChange } from '@/features/datalink/workbench-v2/state/basicRecordingMembers';
import { buildGroupCandidates } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { groupProfileKind, defaultCircuitProfileLocalDraft } from '@/features/datalink/workbench-v2/state/circuitProfileBasic';
import { groupProfileInterval } from '@/features/datalink/workbench-v2/state/circuitProfileBatchOutcome';
import { groupToDraft, draftToRequestGroup, draftIsDirty, validateDraft, UnsupportedRecordingEditorError } from '@/features/datalink/workbench-v2/state/writeGroup/draft';
import { CircuitProfileStatusList } from '@/features/datalink/workbench-v2/steps/step4/CircuitProfileStatusList';
vi.mock('react-i18next',()=>({useTranslation:()=>({t:(key:string)=>key})}));
vi.mock('@/features/datalink/workbench-v2/steps/step4/CircuitProfileStatus',()=>({CircuitProfileStatus:()=> <div data-testid="legacy-profile-status"/>}));
function mixed(): WriteGroup {
 const base=savedGroup();
 return {...base,record_id:'11111111-1111-4111-8111-111111111111',layout_version:'managed_wide_profile_v1',recording_profile:{kind:'mixed',result_shape:'member-calculation-v1'},destination:{...base.destination,storage_strategy:'managed'},row_policy:{allowed_lateness_seconds:0,incomplete_policy:'partial',due_at_column:'due_at'},members:base.members.map(m=>({...m,calculation:'mean',cadence:{interval:15,unit:'second',alignment:'utc_epoch'},source_value_type:'float64',result_value_type:'float64',value_role:'gauge'}))};
}
it('does not use mixed members or cadence as a basic recording update',()=>{
 const group=mixed();const original=structuredClone(group);
 expect(basicMemberChange(group,'dev-1',buildGroupCandidates(savedGroupState()).candidates)).toBeUndefined();
 expect(group).toEqual(original);
});
it('does not classify mixed as a legacy circuit profile or borrow a group interval',()=>{
 const group=mixed();expect(groupProfileKind(group)).toBeUndefined();
 expect(groupProfileInterval(group,'instantaneous')).toBeUndefined();
 expect(groupProfileInterval(group,'cumulative')).toBeUndefined();
 const local=defaultCircuitProfileLocalDraft([], [group]);
 expect(local.instantaneous_enabled).toBe(false);expect(local.cumulative_enabled).toBe(false);
 expect(local.selected_keys).toEqual([]);
});
it('retains mixed facts in the read-only draft and refuses legacy mutation conversion',()=>{
 const group=mixed();const original=structuredClone(group);const draft=groupToDraft(group);
 expect(draft.record_id).toBe(group.record_id);expect(draft).not.toHaveProperty('interval_seconds');expect(draft.due_at_column).toBe('due_at');
 expect(draft.members[0]).toMatchObject(group.members[0]);expect(draft.members[0].cadence).not.toBe(group.members[0].cadence);
 expect(draftIsDirty(draft,group)).toBe(false);
 expect(validateDraft(draft,{destinationSaved:true})).toContainEqual({code:'profile-editor-unsupported'});
 expect(()=>draftToRequestGroup(draft,group.workspace_id)).toThrow(UnsupportedRecordingEditorError);
 expect(group).toEqual(original);
});
it('shows an explicit advanced handoff without a legacy status projection',()=>{
 const group=mixed();const open=vi.fn();render(<CircuitProfileStatusList groups={[group]} candidates={[]} onOpenAdvanced={open}/>);
 expect(screen.getByTestId('mixed-recording-readonly')).toBeInTheDocument();
 expect(screen.queryByTestId('legacy-profile-status')).not.toBeInTheDocument();
 fireEvent.click(screen.getByRole('button'));expect(open).toHaveBeenCalledWith(group.id);
});
it('does not let changing a local profile label erase the mixed record authority',()=>{
 const draft=groupToDraft(mixed());
 draft.recording_profile={kind:'instantaneous',interval:15,unit:'second'};
 expect(()=>draftToRequestGroup(draft,'workspace')).toThrow(UnsupportedRecordingEditorError);
 expect(validateDraft(draft,{destinationSaved:true})).toContainEqual({code:'profile-editor-unsupported'});
});
