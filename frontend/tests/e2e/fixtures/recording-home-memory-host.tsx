// Isolated memory view for inspector / agent-browser review: in-memory saved ranges only,
// no API, no device and no write. Query: ?scenario=mc|modbus|unconfirmed|conflict|w, ?light, ?en.
import { useState } from 'react';
import { createRoot } from 'react-dom/client';
import i18n from '../../../src/i18n/config';
import { RecordingHomeMemoryView } from '../../../src/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView';
import type { MemoryRuleInput } from '../../../src/features/datalink/workbench-v2/recording-home/recordingHomeMemoryMap';
import type { ProtocolType } from '../../../src/types/datalink';
import '../../../src/index.css';
import '../../../src/styles/semantic-modal.css';
import '../../../src/features/datalink/workbench-v2/styles/workbench-v2.css';
import '../../../src/features/datalink/workbench-v2/recording-home/recording-home.css';

const params = new URLSearchParams(window.location.search);
await i18n.changeLanguage(params.has('en') ? 'en' : 'zh-TW');
if (!params.has('light')) document.documentElement.classList.add('dark');
const span = (area: string, start: number, count: number, width: number) => Array.from({ length: count }, (_, index) => `${area}${start + index * width}`);
const scenarios: Record<string, { protocol: ProtocolType; inputs: MemoryRuleInput[] }> = {
  mc: { protocol: 'mc_3e', inputs: [
    { id: 'I', label: 'TEMP_', start_address: 'D0', count: 32, data_type: 'int16', skipped_addresses: [], proven: span('D', 0, 32, 1) },
    { id: 'F', label: 'FLOW_', start_address: 'D39', count: 2, data_type: 'float32', skipped_addresses: [], proven: span('D', 39, 2, 2) },
    { id: 'M', label: 'RUN_', start_address: 'M0', count: 8, data_type: 'bool', skipped_addresses: [], proven: span('M', 0, 8, 1) },
  ] },
  modbus: { protocol: 'modbus_tcp', inputs: [
    { id: 'A', label: 'PRESSURE_', start_address: '40001', count: 4, data_type: 'int16', skipped_addresses: ['40003'] },
    { id: 'B', label: 'ENERGY_', start_address: '40010', count: 3, data_type: 'float32', skipped_addresses: [] },
    { id: 'C', label: 'AMBIENT_', start_address: '30001', count: 2, data_type: 'uint16', skipped_addresses: [] },
  ] },
  // W is hexadecimal (16 per logical row); W9 float32 ends a 10-column rendered row on tablet.
  w: { protocol: 'mc_3e', inputs: [
    { id: 'W', label: 'LINK_', start_address: 'W9', count: 1, data_type: 'float32', skipped_addresses: [], proven: ['W9'] },
  ] },
  conflict: { protocol: 'modbus_tcp', inputs: [
    { id: 'A', label: 'LEFT_', start_address: '40001', count: 4, data_type: 'int16', skipped_addresses: [] },
    { id: 'B', label: 'RIGHT_', start_address: '40003', count: 1, data_type: 'float32', skipped_addresses: [] },
  ] },
  unconfirmed: { protocol: 'mc_3e', inputs: [
    { id: 'I', label: 'TEMP_', start_address: 'D0', count: 32, data_type: 'int16', skipped_addresses: [], proven: span('D', 0, 32, 1) },
    { id: 'S', label: 'LABEL_', start_address: 'D100', count: 2, data_type: 'string', skipped_addresses: [] },
  ] },
};
const scenario = scenarios[params.get('scenario') ?? 'mc'] ?? scenarios.mc;

function Root() {
  const [last, setLast] = useState('');
  return <main data-workbench-v2 data-testid="recording-home" className="gw-recording-home"><h1>{i18n.t('recording-home:title', { defaultValue: '記錄工作台' })}</h1>
    <div role="dialog" aria-labelledby="fixture-title" data-workbench-v2 className="gw-modal gw-home-source-modal" data-size="form" data-source-kind="points" style={{ color: 'var(--gw-text)', backgroundColor: 'var(--gw-raised)' }}>
      <header className="gw-modal-header"><div><h2 id="fixture-title">Line 3 PLC · {i18n.t('recording-home:source.memory.title')}</h2><p>{i18n.t('recording-home:source.sharedDescription')}</p></div>
        <button type="button">{i18n.t('recording-home:source.close')}</button></header>
      <RecordingHomeMemoryView protocol={scenario.protocol} inputs={scenario.inputs} status="ready" checking={false}
        onOpen={ruleId => setLast(`open:${ruleId}`)} onAdd={address => setLast(`add:${address}`)} onBack={() => setLast('back')} />
      <output data-testid="fixture-last-action" hidden>{last}</output>
    </div>
  </main>;
}
createRoot(document.getElementById('root')!).render(<Root />);
