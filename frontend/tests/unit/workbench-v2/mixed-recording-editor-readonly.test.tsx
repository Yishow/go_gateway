import { expect, it, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { GroupEditor } from '@/features/datalink/workbench-v2/steps/step4/writeGroup/GroupEditor';
import { savedGroup, savedGroupState } from '../../fixtures/writeGroupState';
const calls=vi.hoisted(()=>({probe:vi.fn(()=>{throw new Error('legacy editor probe reached for mixed recording');})}));
vi.mock('react-i18next',()=>({useTranslation:()=>({t:(key:string)=>key})}));
vi.mock('@/hooks/datalink/useStudioV2WriteGroups',()=>({useCreateWriteGroupMutation:calls.probe,useUpdateWriteGroupMutation:calls.probe,useWriteGroupReadinessQuery:calls.probe}));
vi.mock('@/hooks/datalink/useStudioV2WorkspaceDatabase',()=>({useStudioV2MySQLRecordingCapabilityQuery:calls.probe,useStudioV2SQLServerRecordingCapabilityQuery:calls.probe}));
it.each(['sqlite','mysql'] as const)('blocks mixed legacy editor before %s probes or mutation hooks',kind=>{
 calls.probe.mockClear();
 const base=savedGroup();const group={...base,record_id:'11111111-1111-4111-8111-111111111111',recording_profile:{kind:'mixed' as const,result_shape:'member-calculation-v1' as const},row_policy:{allowed_lateness_seconds:0}};
 const original=structuredClone(group);const close=vi.fn();const saved=vi.fn();
 expect(()=>render(<GroupEditor group={group} workspaceId={group.workspace_id} workspaceRevision="workspace-rev" connector={{...savedGroupState().db.connector,kind}} candidates={[]} excluded={[]} readonly={false} onClose={close} onSaved={saved} onReload={vi.fn()} fetchLatest={vi.fn()}/>)).not.toThrow();
 expect(calls.probe).not.toHaveBeenCalled();expect(screen.getByTestId('mixed-group-editor-readonly')).toBeInTheDocument();
 expect(screen.queryByTestId('group-save')).not.toBeInTheDocument();
 fireEvent.click(screen.getByRole('button'));expect(close).toHaveBeenCalledOnce();expect(saved).not.toHaveBeenCalled();
 expect(group).toEqual(original);
});
