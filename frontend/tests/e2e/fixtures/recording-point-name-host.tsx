import { useMemo, useState } from 'react';
import { createRoot } from 'react-dom/client';
import i18n from '../../../src/i18n/config';
import { RecordingPointSelection } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingPointSelection';
import { RecordingDraftNavigation } from '../../../src/features/datalink/workbench-v2/recording-editor/RecordingDraftNavigation';
import type { GroupCandidate } from '../../../src/features/datalink/workbench-v2/state/writeGroup/candidates';
import type { EditorMember } from '../../../src/features/datalink/workbench-v2/state/writeGroup/draft';
import type { RecordingPointNamePort, PointNameSession } from '../../../src/pages/datalink/workbench-v2/recordingPointNameSubmission';
import '../../../src/index.css';
import '../../../src/features/datalink/workbench-v2/styles/workbench-v2.css';

await i18n.changeLanguage('zh-TW');
document.documentElement.classList.add('dark');
const candidates: GroupCandidate[] = Array.from({ length: 48 }, (_, index) => {
  const mc = index < 32, number = mc ? index : index - 32;
  return { key: `P${index}|T${index}`, device_id: mc ? 'mc' : 'modbus', point_id: `P${index}`, tag_id: `T${index}`,
    tag_key: `tag${index}`, label: mc ? `TAG_${number} (D${number})` : `MOD_${number} (${40001 + number})`,
    device_name: mc ? 'MC' : 'Modbus', address: mc ? `D${number}` : `${40001 + number}`, target_type: 'int16' };
});
const initialMembers: EditorMember[] = candidates.slice(0, 32).map(row => ({ key: row.key, device_id: row.device_id, point_id: row.point_id,
  tag_id: row.tag_id, entity_key: '', target_column: '', required: false }));

function Host() {
  const [members, setMembers] = useState(initialMembers);
  const [saves, setSaves] = useState<{ point: string; name: string }[]>([]);
  const name = useMemo<RecordingPointNamePort>(() => {
    const sessions = new Map<string, PointNameSession>();
    const names = new Map(candidates.map(row => [row.point_id, row.label]));
    return {
      prepare: async identity => {
        await new Promise(resolve => setTimeout(resolve, 300));
        const session = { id: identity.point_id, identity, displayName: names.get(identity.point_id)! };
        sessions.set(session.id, session);
        return { session, status: 'ready' };
      },
      submit: async (id, value) => {
        await new Promise(resolve => setTimeout(resolve, 80));
        const session = sessions.get(id)!;
        names.set(id, value); setSaves(current => [...current, { point: id, name: value }]);
        return { session, status: 'confirmed', currentName: value };
      },
      inspect: async () => undefined,
      release: () => undefined,
    };
  }, []);
  return <RecordingDraftNavigation><main data-workbench-v2="true" className="min-h-screen bg-slate-950 p-4 text-slate-100">
    <h1>測試資料：點位名稱自動儲存</h1><output data-testid="fixture-saves">{JSON.stringify(saves)}</output>
    <RecordingPointSelection scope="test" workspaceId="test" candidates={candidates} members={members} onChange={setMembers} disabled={false}
      devices={[{ id: 'mc', name: 'MC' }, { id: 'modbus', name: 'Modbus' }]}
      deviceSource={{ name, submit: async () => { throw new Error('Unexpected source write'); }, inspect: async () => undefined, observations: () => [] }} />
  </main></RecordingDraftNavigation>;
}
createRoot(document.getElementById('root')!).render(<Host />);
