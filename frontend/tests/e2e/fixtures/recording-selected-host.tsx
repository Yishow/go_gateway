import { useState } from 'react';
import { createRoot } from 'react-dom/client';
import i18n from '../../../src/i18n/config';
import { SimpleRecordingShell } from '../../../src/features/datalink/workbench-v2/shell/SimpleRecordingShell';
import { RecordingEditorLayout } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingEditorLayout';
import { RecordingDraftNavigation } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import { RecordingPointSelection } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingPointSelection';
import { RecordingSelectedMembers } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingSelectedMembers';
import { RecordingPointOrder } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingPointOrder';
import { RecordingMemberCalculation } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingMemberCalculationPanel';
import { GroupIdentityFields } from '../../../src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupIdentityFields';
import { GroupPolicyFields } from '../../../src/features/datalink/workbench-v2/steps/step4/writeGroup/GroupPolicyFields';
import { RecordingMixedEditor } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingMixedEditor';
import { newDraft } from '../../../src/features/datalink/workbench-v2/state/writeGroup/draft';
import type { MemberCalculationContext } from '../../../src/features/datalink/workbench-v2/recording-editor/recordingMemberCalculation';
import type { PointOrderEntry, PointOrderObservation, RecordingPointOrderPort } from '../../../src/features/datalink/workbench-v2/recording-editor/useRecordingPointOrder';
import type { PointFormatScope, PointFormatSnapshot } from '../../../src/types/pointFormat';
import type { GroupCandidate } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
import type { EditorMember } from '../../../src/features/datalink/workbench-v2/state/writeGroup/draft';
import type { RecordingPointNamePort, PointNameSession } from '../../../src/pages/datalink/workbench-v2/recordingPointNameSubmission';
import '../../../src/index.css';
import '../../../src/features/datalink/workbench-v2/styles/workbench-v2.css';
import '@fontsource/inter/400.css';
import '@fontsource/inter/500.css';
import '@fontsource/inter/600.css';
import '@fontsource/jetbrains-mono/400.css';

const parameters = new URLSearchParams(window.location.search);
await i18n.changeLanguage(parameters.get('lang') === 'en' ? 'en' : 'zh-TW'); document.documentElement.classList.add(parameters.get('theme') === 'light' ? 'light' : 'dark');
const samePointOnDevices = parameters.has('samepoint');
const candidates: GroupCandidate[] = Array.from({ length: 32 }, (_, i) => ({ key: `P${i}|T${i}`, device_id: 'D', point_id: `P${i}`, tag_id: `T${i}`,
  tag_key: `tag${i}`, label: samePointOnDevices && i < 2 ? 'Power (D0)' : `TAG_${i} (D${i})`,
  device_name: samePointOnDevices && i < 2 ? `MC ${i === 0 ? 'A' : 'B'}` : parameters.has('longdevice') && i === 0 ? 'MC Building North' : 'MC',
  address: samePointOnDevices && i === 1 ? 'D0' : `D${i}`, target_type: i === 0 ? 'float32' : 'int16',
  ...(samePointOnDevices && i === 1 ? { device_id: 'D2' } : {}) }));
const mixed = parameters.has('mixed');
const initial: EditorMember[] = candidates.map((candidate, i) => ({ key: candidate.key, device_id: candidate.device_id, point_id: candidate.point_id,
  tag_id: candidate.tag_id, entity_key: '', target_column: `v_${i}`, required: true, ...(mixed ? { calculation: 'latest' as const, cadence: { interval: 15, unit: 'second' as const, alignment: 'utc_epoch' as const } } : {}) }));
const sessions = new Map<string, PointNameSession>();
const showEntityKey = parameters.has('entity');
const calculationContext: MemberCalculationContext = { workspaceId: 'W', definitions: [], layout: 'managed_wide_profile_v1', dialect: 'sqlite', capabilities: {
  version: 1, result_shape: 'member-calculation-v1', editor_supported: true, draft_supported: true, activation_supported: false, reason: 'record_readiness_required',
  methods: ['latest', 'mean', 'delta'], cadence: { alignment: 'utc_epoch', units: ['second', 'minute'], minimum: 1, max_seconds: 600, max_minutes: 10, window_bounds_checked: true },
  schema_dialects: ['sqlite'], schema_layouts: ['managed_wide_profile_v1'], unsupported_combinations: [] } };
const names = new Map(candidates.map(candidate => [candidate.point_id, candidate.label]));
const name: RecordingPointNamePort = {
  prepare: async identity => {
    const session = { id: identity.point_id, identity, displayName: names.get(identity.point_id)! };
    sessions.set(session.id, session); return { session, status: 'ready' };
  },
  submit: async (id, value) => { const session = sessions.get(id)!; names.set(id, value); return { session, status: 'confirmed', currentName: value }; },
  inspect: async () => undefined, release: () => undefined,
};
const orderEntries = new Map<string, PointOrderObservation>(), orderListeners = new Set<() => void>();
let observations: readonly PointOrderObservation[] = [];
const publish = (scope: PointFormatScope, entry: PointOrderEntry) => {
  orderEntries.set(scope.point_id, { scope, entry }); observations = [...orderEntries.values()];
  orderListeners.forEach(listener => listener());
};
function snapshot(scope: PointFormatScope): PointFormatSnapshot {
  const applicable = scope.point_id === 'P0';
  return { ...scope, rule_id: 'R', workspace_revision: 'source-1', workspace_revision_origin: 'persisted', point_revision: 'P1', rule_revision: 'R1', connection_revision: 'C1', editable: true,
    resolution: { resolved: true, format_mode: applicable ? 'override' : 'legacy', configured_format: applicable ? 'ABCD' : '', effective_format: applicable ? 'ABCD' : '',
      effective_source: applicable ? 'override' : 'native', effective_policy: applicable ? 'override_scalar' : 'legacy_mc3e', native_compatibility: false,
      applicable, format_version: 1, workspace_revision: 'source-1', point_revision: 'P1', rule_revision: 'R1', connection_revision: 'C1' } };
}
const pointOrder: RecordingPointOrderPort = {
  available: () => true, subscribe: listener => { orderListeners.add(listener); return () => { orderListeners.delete(listener); }; },
  observations: () => observations, get: scope => orderEntries.get(scope.point_id)?.entry,
  load: async scope => { if (!orderEntries.has(scope.point_id)) publish(scope, { status: 'ready', snapshot: snapshot(scope) }); },
  choose: (scope, intent) => { const entry = orderEntries.get(scope.point_id)?.entry; if (entry) publish(scope, { ...entry, intent }); },
  save: async () => undefined, inspect: async () => undefined, recover: async () => undefined, restart: async () => undefined,
};
function Host() {
  const [members, setMembers] = useState(initial), [rows, setRows] = useState(candidates);
  const [draft, setDraft] = useState(() => ({ ...newDraft({ connector_id: 'C', connector_revision: 'CR', table_schema: '', table_name: '' }), name: 'Fixture MC' }));
  const [selectedScope, setSelectedScope] = useState<{ shown: number; total: number }>();
  return <main data-workbench-v2 className="bg-canvas min-h-screen font-sans text-[var(--gw-text)] antialiased"><RecordingDraftNavigation><SimpleRecordingShell mode="located" workspaceId="W" onBack={() => undefined} onAdvanced={() => undefined}>
    <div className="recording-editor-simple"><RecordingEditorLayout presentation="simple" memberCount={members.length} selectedScope={selectedScope} selectedScopeInHeading perPoint={false}
      identity={<GroupIdentityFields section="name" draft={draft} readonly={false} locked={false} onChange={patch => setDraft(current => ({ ...current, ...patch }))} onStorageStrategy={() => undefined} />}
      sources={<RecordingPointSelection scope="W/G" workspaceId="W" candidates={rows} members={members} onChange={setMembers} disabled={false}
        devices={samePointOnDevices ? [{ id: 'D', name: 'MC A' }, { id: 'D2', name: 'MC B' }] : [{ id: 'D', name: 'MC' }]} deviceSource={{ name, pointOrder, submit: async () => { throw new Error('Source creation is unavailable in this fixture'); }, inspect: async () => undefined, observations: () => [] }}
        onDiscovery={observation => setRows([...observation.value.candidates])}
        renderMembers={actions => <RecordingSelectedMembers {...actions} presentation="simple" managed scopeInHeading onScopeChange={setSelectedScope} candidates={rows} members={members} columns={initial.map((member, index) => ({ name: member.target_column, type: index === 0 ? 'REAL' : 'INTEGER', nullable: true, primary_key: false }))} metadataStatus="exists"
          suggestions={{ assignments: {}, unmatched: [] }} issues={[]} showEntityKey={showEntityKey} disabled={false} onChange={setMembers}
          renderSource={(member, inline) => <RecordingPointOrder simple inline={inline} port={pointOrder} workspaceId="W" deviceId={member.device_id} pointId={member.point_id} disabled={false} />}
          renderCalculation={mixed ? (member, tableCells) => <RecordingMemberCalculation compact tableCells={tableCells} member={member} context={calculationContext} disabled={false} onChange={next => setMembers(current => current.map(value => value.key === member.key ? next : value))} /> : undefined}
          pointOrderSort={{ workspaceId: 'W', port: pointOrder }} labels={{ remove: member => `${parameters.get('lang') === 'en' ? 'Remove' : '移除'} ${member.point_id}`, unavailable: '來源不可用', undo: '復原' }} />} />}
      method={<div style={{ minHeight: 480 }}>記錄方式內容</div>} destination={<><GroupIdentityFields section="destination" draft={draft} readonly={false} locked={false} onChange={patch => setDraft(current => ({ ...current, ...patch }))} onStorageStrategy={storage_strategy => setDraft(current => ({ ...current, storage_strategy }))} />{mixed && <RecordingMixedEditor section="destination" draft={{ ...draft, recording_profile: { kind: 'mixed', result_shape: 'member-calculation-v1' }, layout_version: 'managed_wide_profile_v1' }} saved={false} disabled={false} capabilities={calculationContext.capabilities} dialect="sqlite" onChange={setDraft} />}<div style={{ minHeight: 480 }}>目的地內容</div></>}
      policy={<GroupPolicyFields expanded draft={draft} readonly={false} profile={false} onChange={patch => setDraft(current => ({ ...current, ...patch }))} onDedupe={dedupe_capability => setDraft(current => ({ ...current, dedupe_capability }))} />}
      summary={<><strong>這組記錄</strong><p>1 台設備／32 點</p></>} />
    </div>
  </SimpleRecordingShell></RecordingDraftNavigation></main>;
}
createRoot(document.getElementById('root')!).render(<Host />);
