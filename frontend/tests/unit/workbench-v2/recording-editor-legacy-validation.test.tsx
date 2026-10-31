import { isKnownRuleValidationRejection } from '@/pages/datalink/workbench-v2/sourceQueueObservation';
import { RuleTabRail } from '@/features/datalink/workbench-v2/steps/step2/RuleTabRail';
import { deviceFixture, resetDeviceAutosaveMocks, studioV2WorkspaceAPI, studioV2WorkspaceDevicesAPI, studioV2RulesAPI, studioV2MappingsAPI } from './device-autosave-page.testHarness';
import { TABLE_COLUMNS } from '../../fixtures/writeGroupState';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, expect, it, vi } from 'vitest';
import { useStudioV2AutosaveState } from '@/pages/datalink/workbench-v2/useStudioV2AutosaveState';
import { RecordingPointSelection } from '@/features/datalink/workbench-v2/recording-editor/RecordingPointSelection';
import { buildGroupCandidates, type GroupCandidate } from '@/features/datalink/workbench-v2/state/writeGroup/candidates';
import { sourceRuleCandidateAPI } from '@/services/sourceRuleCandidates';
import { deviceAPI, pointAPI, tagAPI } from '@/services/datalink';
import type { SourceRuleRecord, StudioV2WorkspaceMappingRecord } from '@/types/datalink';
vi.mock('react-i18next', () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock('@/services/sourceRuleCandidates', () => ({ sourceRuleCandidateAPI: { get: vi.fn() } }));
vi.mock('@/services/datalink', async original => { const actual = await original<typeof import('@/services/datalink')>(); return { ...actual, deviceAPI: { ...actual.deviceAPI, checkReadiness: vi.fn(), testConnection: vi.fn(), activate: vi.fn() }, pointAPI: { ...actual.pointAPI, get: vi.fn() }, tagAPI: { ...actual.tagAPI, get: vi.fn() } }; });
vi.mock('@/features/datalink/workbench-v2/steps/step4/useStep4TargetColumns', () => ({ useStep4TargetColumns: () => ({ columns: TABLE_COLUMNS, status: 'exists', refetch: vi.fn() }) }));
vi.mock('@/services/studioV2WorkspaceWriteGroups', () => ({ studioV2WorkspaceWriteGroupsAPI: { list: vi.fn(), readiness: vi.fn(), delivery: vi.fn(), create: vi.fn(), update: vi.fn(), apply: vi.fn(), disable: vi.fn(), remove: vi.fn(), testWritePreview: vi.fn(), testWrite: vi.fn(), testWriteOperation: vi.fn(), get: vi.fn() } }));
vi.mock('@/services/studioV2LocatorEntry', () => ({ readStudioLocatorEditorEntry: vi.fn() }));
vi.mock('@/features/datalink/workbench-v2/recording-home/recordingHomeSources', () => ({ readRecordingHomeSources: async (scope: object) => ({ state: 'available', scope, value: { candidates: [], excluded: [], devices: [device] } }) }));
vi.mock('@/features/datalink/workbench-v2/recording-home/useRecordingHomeFacts', () => ({ useRecordingHomeFacts: () => ({ owners: {}, runtimeGroups: {}, histories: {}, pendingDevices: [], retry: vi.fn() }) }));
const device = deviceFixture({ id: 'dev-1', name: 'PLC', last_test_success: true });
let diagnosed = new Set<string>();
let sourceDevice = device;
let rule: SourceRuleRecord | undefined;
let mapping: StudioV2WorkspaceMappingRecord | undefined;
function Harness({ extra = [] }: { extra?: GroupCandidate[] }) {
 const owner = useStudioV2AutosaveState(true, undefined, false, undefined, 'W');
 if (!owner.workspaceHydrated) return <p>Loading</p>;
 return <><output data-testid="source-write-state">{owner.deviceSource.range?.sourceWriteState?.('dev-1', 'POINT') ?? 'settled'}</output><RuleTabRail rules={owner.state.rules} devices={owner.state.devices} selectedRuleId={owner.state.rules[0]?.id ?? null} dispatch={owner.actions.dispatch} />{owner.state.rules.map(item => <div key={item.id}><input aria-label="Legacy address" value={item.start_address} onChange={event => owner.actions.dispatch({type:'updateRule',ruleId:item.id,patch:{start_address:event.target.value}})} /><output data-testid="legacy-save-state">{item.save_state}</output><output data-testid="legacy-enabled">{String(item.enabled)}</output></div>)}<input aria-label="Recording draft" defaultValue="Keep my recording" /><RecordingPointSelection scope="S" workspaceId="W" candidates={[...extra, ...buildGroupCandidates(owner.state).candidates]}
  members={[]} onChange={vi.fn()} devices={owner.state.devices} disabled={false} deviceSource={owner.deviceSource} /></>;
}
beforeEach(() => {
 resetDeviceAutosaveMocks(); vi.clearAllMocks(); diagnosed = new Set(['dev-1']); sourceDevice = device; rule = undefined; mapping = undefined;
 studioV2WorkspaceAPI.get.mockResolvedValue({ id: 'W', kind: 'single', status: 'ready', ordered_device_ids: ['dev-1'], created_at: '', updated_at: '', modbus_share: { workspace_id: 'W', workspace_revision: 'SOURCE-REV' } });
 studioV2WorkspaceDevicesAPI.list.mockResolvedValue([device]);
 vi.mocked(deviceAPI.checkReadiness).mockImplementation(async id => ({ device_id: id, status: diagnosed.has(id) ? 'ready' : 'error', connect_status: diagnosed.has(id) ? 'success' : 'unknown', probe_status: diagnosed.has(id) ? 'success' : 'unknown', activation_allowed: diagnosed.has(id), planning_allowed: diagnosed.has(id), apply_allowed: diagnosed.has(id), blocking_reasons: [], checks: [] }));
 vi.mocked(deviceAPI.testConnection).mockImplementation(async id => { diagnosed.add(id); return { success: true, error: '', latency_ms: 1, can_activate: true }; });
 studioV2RulesAPI.list.mockImplementation(async () => rule ? [rule] : []);
 studioV2MappingsAPI.list.mockImplementation(async () => mapping ? [mapping] : []);
 studioV2RulesAPI.create.mockImplementation(async request => {
  rule = { ...request, naming_prefix: request.naming_prefix.toUpperCase().replace(/ /g, '-'), start_address: request.start_address.trim().toUpperCase(), id: request.id!, workspace_id: 'W', revision_id: 'R1', created_at: '', updated_at: '' } as SourceRuleRecord;
  return { ...rule, runtime_apply_status: 'not_running' };
 });
 studioV2MappingsAPI.create.mockImplementation(async request => {
  mapping = { ...request, id: 'MAP', workspace_id: 'W', device_id: sourceDevice.id, point_id: 'POINT', tag_id: 'TAG', created_at: '', updated_at: '' };
  return { ...mapping, runtime_apply_status: 'not_running' };
 });
 studioV2RulesAPI.update.mockImplementation(async (id, request) => { rule = { ...rule!, ...request, id, revision_id: 'R2' }; return { ...rule, runtime_apply_status: 'not_running' }; });
 studioV2MappingsAPI.update.mockImplementation(async (id, request) => { mapping = { ...mapping!, ...request, id }; return { ...mapping, runtime_apply_status: 'not_running' }; });
 vi.mocked(sourceRuleCandidateAPI.get).mockImplementation(async () => ({ source_rule_id: rule!.id, workspace_id: 'W', workspace_revision: 'SOURCE-REV', revision_id: rule!.revision_id!,
  tags: { status: 'ready', candidates: [{ id: 'CANDIDATE', identity: { source_rule_id: rule!.id, candidate_type: 'tag', candidate_kind: 'point', derived_from_rule_address: '40009' }, proposed_signature: 'SIG',
    address: '40009', point_id: 'POINT', tag_key: `${rule!.naming_prefix}_40009`, display_name: `${rule!.naming_prefix}_40009`, data_type: 'int16', status: 'draft' }] },
  database_outputs: { status: 'ready', candidates: [] }, local_modbus_outputs: { status: 'deferred', candidates: [] } }));
 vi.mocked(pointAPI.get).mockImplementation(async () => ({ id: 'POINT', device_id: sourceDevice.id, address: '40009', data_type: 'int16', enabled: rule?.enabled ?? true } as Awaited<ReturnType<typeof pointAPI.get>>));
 vi.mocked(tagAPI.get).mockImplementation(async () => ({ id: 'TAG', key: mapping!.tag_key, data_type: mapping!.target_type, display_name: mapping!.display_name, ...(mapping!.unit ? { unit: mapping!.unit } : {}), status: 'draft' } as Awaited<ReturnType<typeof tagAPI.get>>));
});
async function openRange() {
 const client = new QueryClient({ defaultOptions: { queries: { retry: false }, mutations: { retry: false } } });
 render(<QueryClientProvider client={client}><Harness /></QueryClientProvider>);
 fireEvent.click(await screen.findByTestId('recording-points-open'));
 fireEvent.change(screen.getByLabelText('picker.device'), { target: { value: 'dev-1' } });
 fireEvent.click(screen.getByRole('button', { name: 'picker.newSource' }));
 fireEvent.change(await screen.findByLabelText('range.name'), { target: { value: 'Sensor' } });
 fireEvent.change(screen.getByLabelText('range.address'), { target: { value: '40009' } });
 fireEvent.change(screen.getByLabelText('range.count'), { target: { value: '1' } });
 await waitFor(() => expect(screen.getByRole('button', { name: 'range.save' })).toBeEnabled());
}

async function confirmedSource() {
 await openRange(); fireEvent.click(screen.getByRole('button',{name:'range.save'}));
 await screen.findByLabelText('picker.search'); fireEvent.click(screen.getByRole('button',{name:'picker.cancel'}));
}
it('known HTTP400 validation permits corrected address and one later disable PUT without reload',async()=>{
 await confirmedSource();
 studioV2RulesAPI.update.mockRejectedValueOnce({response:{status:400,data:{success:false,error:{code:'validation',message:'source range overlap'}}}});
 fireEvent.change(screen.getByLabelText('Legacy address'),{target:{value:'40001'}});
 await waitFor(()=>expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(1));
 await waitFor(()=>expect(screen.getByTestId('legacy-save-state')).toHaveTextContent('save-error'));
 expect(screen.getByLabelText('Legacy address')).toHaveValue('40001');
 expect(rule?.start_address).toBe('40009'); expect(rule?.enabled).toBe(true);
 expect(screen.getByTestId('source-write-state')).toHaveTextContent('settled');
 fireEvent.change(screen.getByLabelText('Legacy address'),{target:{value:'40009'}});
 await waitFor(()=>expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(2));
 await waitFor(()=>expect(screen.getByTestId('legacy-save-state')).toHaveTextContent('saved'));
 fireEvent.click(screen.getByTestId(`rule-toggle-enabled-${rule!.id}`));
 await waitFor(()=>expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(3));
 await waitFor(()=>expect(screen.getByTestId('legacy-save-state')).toHaveTextContent('saved'));
 expect(studioV2RulesAPI.update.mock.calls[2][1]).toMatchObject({start_address:'40009',enabled:false});
 expect(rule?.enabled).toBe(false); expect(studioV2RulesAPI.create).toHaveBeenCalledTimes(1);
});

it.each([503,400])('ambiguous or non-validation HTTP%s never clears the old range into a resend',async status=>{
 await confirmedSource();
 studioV2RulesAPI.update.mockRejectedValueOnce({response:{status,data:{success:false,error:{code:status===503?'validation':'unknown_failure'}}}});
 fireEvent.change(screen.getByLabelText('Legacy address'),{target:{value:'40001'}});
 await waitFor(()=>expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(1));
 await waitFor(()=>expect(screen.getByTestId('legacy-save-state')).not.toHaveTextContent('saving'));
 fireEvent.change(screen.getByLabelText('Legacy address'),{target:{value:'40009'}});
 fireEvent.click(screen.getByTestId(`rule-toggle-enabled-${rule!.id}`));
 fireEvent.click(screen.getByTestId('recording-points-open'));
 fireEvent.click(screen.getByTestId('recording-range-resume'));
 fireEvent.click(screen.getByTestId('recording-range-inspect'));
 await waitFor(()=>expect(screen.getByTestId('recording-range-unknown')).toBeVisible());
 expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(1);
 expect(rule?.enabled).toBe(true);
});
it('a late known validation rejection sends only the newer corrected body, never replays the bad address',async()=>{
 await confirmedSource();
 let reject!: (reason: unknown)=>void;
 studioV2RulesAPI.update.mockImplementationOnce(()=>new Promise((_resolve,no)=>{reject=no;}));
 fireEvent.change(screen.getByLabelText('Legacy address'),{target:{value:'40001'}});
 await waitFor(()=>expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(1));
 fireEvent.change(screen.getByLabelText('Legacy address'),{target:{value:'40009'}});
 reject({response:{status:400,data:{error:{code:'validation'}}}});
 await waitFor(()=>expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(2));
 await waitFor(()=>expect(screen.getByTestId('legacy-save-state')).toHaveTextContent('saved'));
 expect(studioV2RulesAPI.update.mock.calls.map(call=>call[1].start_address)).toEqual(['40001','40009']);
});
it('a foreign legacy acknowledgement cannot clear the source uncertainty using the old receipt',async()=>{
 await confirmedSource();
 studioV2RulesAPI.update.mockImplementationOnce(async(_id,request)=>({...rule!,...request,id:'FOREIGN',revision_id:'R2',runtime_apply_status:'not_running'}));
 fireEvent.click(screen.getByTestId(`rule-toggle-enabled-${rule!.id}`));
 await waitFor(()=>expect(screen.getByTestId('legacy-save-state')).toHaveTextContent('save-error'));
 fireEvent.click(screen.getByTestId(`rule-toggle-enabled-${rule!.id}`));
 expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(1);
 expect(rule?.enabled).toBe(true);
});

it('independent: client-level validation code without server typed validation must not unlock uncertain legacy write',async()=>{
 await confirmedSource();
 studioV2RulesAPI.update.mockRejectedValueOnce({code:'validation',response:{status:400,data:{message:'unclassified error'}}});
 fireEvent.change(screen.getByLabelText('Legacy address'),{target:{value:'40001'}});
 await waitFor(()=>expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(1));
 await waitFor(()=>expect(screen.getByTestId('legacy-save-state')).not.toHaveTextContent('saving'));
 fireEvent.change(screen.getByLabelText('Legacy address'),{target:{value:'40009'}});
 await new Promise(resolve=>setTimeout(resolve,20));
 expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(1);
});

it.each([
 {error:{code:'validation'},response:{status:400,data:{message:'unknown'}}},
 {response:{status:400,data:{code:'validation'}}},
 {response:{status:400,data:{error:{error:{code:'validation'}}}}},
 {response:{status:400,data:{error:['validation']}}},
 {response:{status:503,data:{error:{code:'validation'}}}},
 {response:{status:400,data:{error:{code:'validation'+'x'.repeat(1024)}}}},
])('does not use alternative or malformed metadata as server zero-write proof (%#)', error => {
 expect(isKnownRuleValidationRejection(error)).toBe(false);
});

it('a later legacy edit on a confirmed range accepts canonical normalization without reviving the old receipt', async () => {
 await confirmedSource();
 studioV2RulesAPI.update.mockImplementationOnce(async (id, request) => {
  rule = { ...rule!, ...request, id, start_address: '40009', enabled: false, revision_id: 'R2' };
  return { ...rule, runtime_apply_status: 'not_running' };
 });
 fireEvent.change(screen.getByLabelText('Legacy address'), { target: { value: ' 40009 ' } });
 await waitFor(() => expect(screen.getByTestId('legacy-save-state')).toHaveTextContent('saved'));
 expect(screen.getByLabelText('Legacy address')).toHaveValue('40009');
 expect(screen.getByTestId('legacy-enabled')).toHaveTextContent('false');
 expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(1);
 fireEvent.change(screen.getByLabelText('Legacy address'), { target: { value: '40017' } });
 await waitFor(() => expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(2));
});

it.each([false, true])('a lost legacy PUT preserves attempted input, blocks old candidates, and reads current state without proving success (server wrote=%s)', async wrote => {
 await confirmedSource();
 const originalId = rule!.id;
 expect(screen.getByTestId('source-write-state')).toHaveTextContent('settled');
 studioV2RulesAPI.update.mockImplementationOnce(async (id, request) => {
  if (wrote) rule = { ...rule!, ...request, id, revision_id: 'R2' };
  throw new Error('lost response');
 });
 fireEvent.change(screen.getByLabelText('Legacy address'), { target: { value: '40711' } });
 await waitFor(() => expect(screen.getByTestId('legacy-save-state')).not.toHaveTextContent('saving'));
 expect(screen.getByTestId('source-write-state')).toHaveTextContent('unknown');
 fireEvent.click(screen.getByTestId('recording-points-open'));
 expect(screen.getByRole('checkbox')).toBeDisabled();
 expect(screen.getByRole('button', { name: 'picker.add' })).toBeDisabled();
 fireEvent.click(screen.getByTestId('recording-range-resume'));
 expect(screen.getByLabelText('range.address')).toHaveValue('40711');
 const reads = studioV2RulesAPI.list.mock.calls.length;
 fireEvent.click(screen.getByTestId('recording-range-inspect'));
 await waitFor(() => expect(studioV2RulesAPI.list).toHaveBeenCalledTimes(reads + 1));
 expect(screen.getByTestId('recording-range-current-state')).toHaveTextContent(wrote ? '40711' : '40009');
 expect(screen.getByTestId('recording-range-unknown')).toBeVisible();
 expect(screen.getByLabelText('range.address')).toHaveValue('40711');
 expect(screen.getByRole('button', { name: 'range.save' })).toBeDisabled();
 const review = screen.getByRole('link', { name: 'range.reviewSource' });
 expect(review).toHaveAttribute('target', '_blank'); expect(review.getAttribute('href')).toContain(`rule_id=${originalId}`);
 expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(1); expect(studioV2RulesAPI.create).toHaveBeenCalledTimes(1);
});

it('a failed current-state GET does not claim a successful inspection or release source uncertainty', async () => {
 await confirmedSource();
 studioV2RulesAPI.update.mockRejectedValueOnce(new Error('lost response'));
 fireEvent.change(screen.getByLabelText('Legacy address'), { target: { value: '40711' } });
 await waitFor(() => expect(screen.getByTestId('source-write-state')).toHaveTextContent('unknown'));
 fireEvent.click(screen.getByTestId('recording-points-open')); fireEvent.click(screen.getByTestId('recording-range-resume'));
 studioV2RulesAPI.list.mockRejectedValueOnce(new Error('GET failed'));
 fireEvent.click(screen.getByTestId('recording-range-inspect'));
 await screen.findByTestId('recording-range-inspection-unavailable');
 expect(screen.queryByTestId('recording-range-current-state')).not.toBeInTheDocument();
 expect(screen.getByLabelText('range.address')).toHaveValue('40711');
 expect(studioV2RulesAPI.update).toHaveBeenCalledTimes(1);
});
