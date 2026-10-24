import { createRoot } from 'react-dom/client';
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import '../../../src/index.css';
import zh from '../../../src/i18n/locales/zh-TW/workbench-v2.json';
import { MappingTable } from '../../../src/features/datalink/workbench-v2/steps/step3/MappingTable';
import type { Device, Mapping, Point } from '../../../src/features/datalink/workbench-v2/state/types';

await i18n.use(initReactI18next).init({
  lng: 'zh-TW',
  resources: { 'zh-TW': { 'workbench-v2': zh } },
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
}]));

createRoot(document.getElementById('root')!).render(
  <main style={{ padding: 24, color: '#e2e8f0', background: '#0f172a' }}>
    <header style={{ height: 180 }}>
      <h1>命名與換算</h1>
      <p>48 個點位 · 表格捲動驗證（模擬資料）</p>
    </header>
    <MappingTable points={points} mappings={mappings} selectedIdx={null}
      setSelectedIdx={() => undefined} devices={devices} rawValues={{}}
      connectionByDevice={{}} dispatch={() => undefined} />
    <footer style={{ height: 500 }}>表格外的頁面內容</footer>
  </main>,
);
