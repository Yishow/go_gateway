// Isolated memory view for inspector / agent-browser review: in-memory saved ranges only,
// Mock sample readings only: no API, no device, no HSL connection and no write.
// Query: ?scenario=reference|mc|modbus|unconfirmed|conflict|w, ?light, ?en, ?controls.
// Optional blocked/read/values query parameters select the same mock-only controls below.
import { useCallback, useMemo, useState } from 'react';
import { createRoot } from 'react-dom/client';
import i18n from '../../../src/i18n/config';
import { RecordingHomeMemoryView, type HomeMemoryBlock } from '../../../src/features/datalink/workbench-v2/recording-home/RecordingHomeMemoryView';
import type { MemoryRuleInput } from '../../../src/features/datalink/workbench-v2/recording-home/recordingHomeMemoryMap';
import { buildMemoryMap, memoryWindow } from '../../../src/features/datalink/workbench-v2/recording-home/recordingHomeMemoryMap';
import { memoryValueKey, type MemoryValues } from '../../../src/features/datalink/workbench-v2/recording-home/recordingHomeMemoryValues';
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
  reference: { protocol: 'modbus_tcp', inputs: [
    { id: 'I16', label: 'SAMPLE_INT16_', start_address: '40001', count: 16, data_type: 'int16', skipped_addresses: [] },
    { id: 'F32', label: 'SAMPLE_FLOAT32_', start_address: '40017', count: 8, data_type: 'float32', skipped_addresses: [] },
    { id: 'F64', label: 'SAMPLE_FLOAT64_', start_address: '40041', count: 1, data_type: 'float64', skipped_addresses: [] },
    { id: 'U16', label: 'SAMPLE_UINT16_', start_address: '40045', count: 1, data_type: 'uint16', skipped_addresses: [] },
    { id: 'I32', label: 'SAMPLE_INT32_', start_address: '40047', count: 1, data_type: 'int32', skipped_addresses: [] },
    { id: 'REMOTE', label: 'SAMPLE_REMOTE_', start_address: '40201', count: 1, data_type: 'int16', skipped_addresses: [] },
    { id: 'IR', label: 'SAMPLE_INPUT_', start_address: '30201', count: 1, data_type: 'uint16', skipped_addresses: [] },
    { id: 'BOOL', label: 'SAMPLE_BOOL_', start_address: '00001', count: 2, data_type: 'bool', skipped_addresses: [] },
  ] },
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

// Explicit deterministic SAMPLE values, supplied through the same MemoryValues contract as production.
// “live” exercises presentation only; this fixture does not prove gateway streaming or device reads.
const sampledAt = '2026-11-04T12:00:00.000Z';
const sampleValues: MemoryValues = { status: 'live', updatedAt: Date.parse(sampledAt), byCell: new Map(
  buildMemoryMap(scenario.protocol, scenario.inputs).areas.flatMap(area =>
    // Include the distant test range as well as the initial 100-word window.
    [0, 100, 200].flatMap(offset => memoryWindow(area, area.base + offset, 100, false)).flatMap(cell => {
      const owner = cell.owners[0];
      if (cell.state !== 'used' || owner?.part !== 1) return [];
      const sample = owner.dataType === 'float32' ? '12.5' : owner.dataType === 'float64' ? '123456.125'
        : owner.dataType === 'uint16' ? '65535' : owner.dataType === 'int32' ? '-123456'
          : owner.dataType === 'bool' ? '1' : cell.offset === area.base ? '0' : '-42';
      return [[memoryValueKey(area.key, cell.offset, owner.dataType), { value: sample, readAt: sampledAt }] as const];
    })),
) };

const blockModes = ['none', 'recovery', 'local', 'stale'] as const;
const readModes = ['ready', 'loading', 'error', 'checking'] as const;
const valueModes = ['live', 'updated', 'stale', 'cell-error', 'error', 'loading', 'polling'] as const;
type ReadMode = typeof readModes[number];
type ValueMode = typeof valueModes[number];
function queryMode<T extends string>(key: string, choices: readonly T[], fallback: T): T {
  return queryChoice(params.get(key), choices, fallback);
}
function mockValues(mode: ValueMode): MemoryValues {
  const status = mode === 'loading' || mode === 'error' || mode === 'polling' ? mode : 'live';
  return { status, updatedAt: mode === 'updated' ? Date.parse(sampledAt) + 1000 : sampleValues.updatedAt,
    byCell: new Map([...sampleValues.byCell].map(([key, sample]) => [key,
      mode === 'loading' ? {} : mode === 'cell-error' ? { error: 'MOCK_READ_FAILURE', readAt: sampledAt }
        : mode === 'updated' ? { value: key.endsWith(':bool') ? '0' : '37', readAt: '2026-11-04T12:00:01.000Z' }
          : { ...sample, stale: mode === 'stale' || mode === 'error' }])) };
}

function Root() {
  const [last, setLast] = useState('');
  const [blocked, setBlocked] = useState<HomeMemoryBlock | 'none'>(() => queryMode('blocked', blockModes, 'none'));
  const [read, setRead] = useState<ReadMode>(() => queryMode('read', readModes, 'ready'));
  const [valueMode, setValueMode] = useState<ValueMode>(() => queryMode('values', valueModes, 'live'));
  const [screen, setScreen] = useState('memory'), [instance, setInstance] = useState(1);
  const [retryCount, setRetryCount] = useState(0), [visibleKeys, setVisibleKeys] = useState<readonly string[]>([]);
  const reportVisible = useCallback((keys: readonly string[]) => setVisibleKeys([...new Set(keys)]), []);
  const values = useMemo(() => ({ ...mockValues(valueMode), setVisible: reportVisible }), [valueMode, reportVisible]);
  const dismiss = (action: 'back' | 'close') => { setLast(action); setScreen(action); setVisibleKeys([]); };
  const reopen = () => { setInstance(current => current + 1); setScreen('memory'); };
  return <main data-workbench-v2 data-testid="recording-home" className="gw-recording-home"><h1>{i18n.t('recording-home:title', { defaultValue: '記錄工作台' })}</h1>
    {screen === 'memory' ? <div role="dialog" aria-labelledby="fixture-title" data-workbench-v2 className="gw-modal gw-home-source-modal" data-size="form" data-source-kind="points" data-source-view="memory" data-fixture-instance={instance} style={{ color: 'var(--gw-text)', backgroundColor: 'var(--gw-raised)' }}>
      <header className="gw-modal-header"><div><h2 id="fixture-title">Line 3 PLC · {i18n.t('recording-home:source.memory.title')}</h2><p>{i18n.t('recording-home:source.sharedDescription')}</p>
      {/* Keep fixture-only content in the header so the real body retains its bounded scrolling row. */}
      <p data-testid="fixture-sample-notice">MOCK SAMPLE readings · not connected to HSL or a device</p>
      {params.has('controls') && <details data-testid="fixture-controls"><summary>MOCK SAMPLE controls</summary>
        <div style={{ display: 'grid', gap: 4 }}>
          <label>Mock source blocking <select value={blocked} onChange={event => setBlocked(queryChoice(event.target.value, blockModes, 'none'))}>
            {blockModes.map(mode => <option key={mode}>{mode}</option>)}</select></label>
          <label>Mock range read <select value={read} onChange={event => setRead(queryChoice(event.target.value, readModes, 'ready'))}>
            {readModes.map(mode => <option key={mode}>{mode}</option>)}</select></label>
          <label>Mock value presentation <select value={valueMode} onChange={event => setValueMode(queryChoice(event.target.value, valueModes, 'live'))}>
            {valueModes.map(mode => <option key={mode}>{mode}</option>)}</select></label>
        </div>
      </details>}
      </div></header>
      <RecordingHomeMemoryView key={instance} values={values} protocol={scenario.protocol} inputs={scenario.inputs}
        status={read === 'checking' ? 'ready' : read} checking={read === 'checking'} blocked={blocked === 'none' ? undefined : blocked}
        onOpen={ruleId => setLast(`open:${ruleId}`)} onAdd={address => setLast(`add:${address}`)} onBack={() => dismiss('back')} onClose={() => dismiss('close')}
        onRetry={() => { setRetryCount(count => count + 1); setRead('ready'); }} />
    </div> : <section data-testid="fixture-dismissed"><p>{screen === 'back' ? 'Back to MOCK SAMPLE settings' : 'MOCK SAMPLE memory closed'}</p>
      <button type="button" data-testid="fixture-reopen" onClick={reopen}>Reopen MOCK SAMPLE memory</button></section>}
    <output data-testid="fixture-last-action" hidden>{last}</output>
    <output data-testid="fixture-retry-count" hidden>{retryCount}</output>
    <output data-testid="fixture-visible-keys" hidden>{JSON.stringify(visibleKeys)}</output>
  </main>;
}
function queryChoice<T extends string>(value: string | null, choices: readonly T[], fallback: T): T {
  return choices.find(choice => choice === value) ?? fallback;
}
createRoot(document.getElementById('root')!).render(<Root />);
