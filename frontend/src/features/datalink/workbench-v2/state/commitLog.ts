import type { WorkbenchV2State, CommitLog } from './types';

/**
 * 產生 10 步提交 Log 序列
 * 落地設計決策：「Commit log 序列：純函式 + reducer 串聯」
 * 
 * @param state 當前的工作台狀態
 * @returns 10 個 CommitLog 的陣列
 */
export function buildCommitLogSequence(state: WorkbenchV2State): CommitLog[] {
  const devices = state.devices || [];
  const enabledRules = (state.rules || []).filter(r => r.enabled);
  const enabledPoints = (state.points || []).filter(p => p.enabled);
  const connector = state.db.connector;
  
  // 計算啟用的 db targets 數量
  const enabledTargetCount = Object.values(state.db.targets || {}).filter(t => t.enabled).length;

  // 裝置連線詳細字串
  const devicesDetail = devices.map(d => `${d.name} (${d.protocol})`).join(' · ') || '無裝置';

  // 啟用規則詳細字串
  const rulesDetail = enabledRules.map(r => r.name).join(' · ') || '無規則';

  // 資料庫連線詳細字串 (區分 SQLite 是否有 schema)
  const dbTestDetail = connector.kind === 'sqlite'
    ? `sqlite ${connector.database}`
    : `${connector.kind} ${connector.host}:${connector.port}`;

  const dbTargetDetail = connector.kind === 'sqlite' || !connector.schema
    ? `→ ${connector.table}`
    : `→ ${connector.schema}.${connector.table}`;

  return [
    {
      label: `POST /devices × ${devices.length}`,
      detail: devicesDetail,
      status: 'pending'
    },
    {
      label: `POST /devices/:id/activate × ${devices.length}`,
      detail: 'draft → active',
      status: 'pending'
    },
    {
      label: `POST /source-rules × ${enabledRules.length}`,
      detail: rulesDetail,
      status: 'pending'
    },
    {
      label: `POST /points × ${enabledPoints.length}`,
      detail: 'bulk create',
      status: 'pending'
    },
    {
      label: 'POST /polling-groups',
      detail: '快速輪詢 1s, enabled',
      status: 'pending'
    },
    {
      label: `POST /tags × ${enabledPoints.length}`,
      detail: 'register tag keys',
      status: 'pending'
    },
    {
      label: `POST /mappings × ${enabledPoints.length}`,
      detail: 'point ↔ tag, scale pipeline',
      status: 'pending'
    },
    {
      label: 'POST /db-connectors/:id/test',
      detail: dbTestDetail,
      status: 'pending'
    },
    {
      label: `POST /db-targets × ${enabledTargetCount}`,
      detail: dbTargetDetail,
      status: 'pending'
    },
    {
      label: 'POST /scheduler/start',
      detail: 'collectors started',
      status: 'pending'
    }
  ];
}
