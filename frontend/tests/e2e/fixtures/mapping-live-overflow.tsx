import { useState } from 'react';
import { createRoot } from 'react-dom/client';
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import '../../../src/index.css';
import '../../../src/features/datalink/workbench-v2/styles/workbench-v2.css';
import zh from '../../../src/i18n/locales/zh-TW/workbench-v2.json';
import mappingErrors from '../../../src/i18n/locales/zh-TW/mapping-errors.json';
import { MappingTable } from '../../../src/features/datalink/workbench-v2/steps/step3/MappingTable';
import type { Device, Mapping, Point } from '../../../src/features/datalink/workbench-v2/state/types';

await i18n.use(initReactI18next).init({
  lng: 'zh-TW', resources: { 'zh-TW': { 'workbench-v2': zh, 'mapping-errors': mappingErrors } },
});
const devices: Device[] = [{ id: 'plc', name: '測試設備', description: '', protocol: 'modbus_tcp', config: {}, status: 'draft', test: null }];
const points: Point[] = Array.from({ length: 48 }, (_, index) => ({
  id: `point-${index}`, device_id: 'plc', rule_id: 'rule', rule_name: '測試規則',
  name: `TAG_${index}`, address: String(40001 + index), data_type: 'int16',
  function: 'holding_register', width: 1, enabled: true, skipped: false, _rule_scale: 1, _rule_offset: 0,
}));
const mappings: Record<string, Mapping> = Object.fromEntries(points.map(point => [point.id, {
  point_id: point.id, tag_key: point.name, display_name: point.name,
  unit: '', target_type: 'int16', scale: 1, offset: 0, enabled: true,
  save_state: point.id === 'point-1' ? 'save-error' : 'saved',
  ...(point.id === 'point-1' ? { save_error_detail: { code: 'workspace_mapping_save_failed', status: 500, requestId: 'test-request' } } : {}),
}]));
const rawValues = Object.fromEntries(points.map((point, index) => [point.id, index + 1]));

function Preview() {
  const [loaded, setLoaded] = useState(false);
  return <div data-workbench-v2="true" data-testid="overflow-page" className="bg-canvas min-h-screen text-slate-100">
    <div className="wbv2-shell min-h-screen flex flex-col">
      <header className="p-5"><h1>命名與換算</h1>
        <button type="button" data-testid="load-live-values" onClick={() => setLoaded(true)}>載入模擬即時值</button>
      </header>
      {/* Keep the real positioned workspace outside the scrollable table, as in Studio. */}
      <div className="wbv2-workspace flex-1 flex px-5 py-5">
        <main className="flex-1 min-w-0">
          <MappingTable points={points} mappings={mappings} selectedIdx={null}
            setSelectedIdx={() => undefined} devices={devices} rawValues={loaded ? rawValues : {}}
            connectionByDevice={{}} dispatch={() => undefined} />
          <footer style={{ height: 120 }}>頁面內容底部</footer>
        </main>
      </div>
    </div>
  </div>;
}
createRoot(document.getElementById('root')!).render(<Preview />);
