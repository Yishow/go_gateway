import type { Point, Mapping, DbTarget } from './types';

function matchesColumn(col: string, token: string): boolean {
  if (!token) return false;
  return col === token || col.endsWith(`_${token}`) || col.startsWith(`${token}_`);
}

/**
 * 自動欄位匹配演算法
 * 
 * 優先級：
 * 1. existingTargets 已存在的配對 (保留不變)
 * 2. 精確匹配 / 後綴 / 前綴匹配 (exact match, endsWith('_tag'), startsWith('tag_'))
 * 3. 索引回退 (index fallback, 若未使用)
 * 4. 第一個未使用的欄位 (first unused)
 * 5. 若皆已使用，回退到索引對應欄位 (重複使用，會產生 conflict)
 * 
 * @param enabledPoints 啟用中的點位清單
 * @param mappings 點位與 Tag 的映射字典
 * @param columnNames 可用的資料庫欄位名稱（不含主鍵/時間戳）
 * @param existingTargets 目前已存在的匹配關係
 * @returns 新的 DbTarget 字典，Key 為 point_id
 */
export function autoAssignTargets(
  enabledPoints: Point[],
  mappings: Record<string, Mapping>,
  columnNames: string[],
  existingTargets: Record<string, DbTarget>
): Record<string, DbTarget> {
  const used = new Set<string>();
  const out: Record<string, DbTarget> = {};

  // 第一階段：先處理已存在的 targets，將其佔用的欄位加入 used
  for (const p of enabledPoints) {
    const m = mappings[p.id];
    if (!m) continue;

    if (existingTargets[p.id]) {
      out[p.id] = { ...existingTargets[p.id] };
      used.add(existingTargets[p.id].column_name);
    }
  }

  // 第二階段：處理未分配的點位
  for (let i = 0; i < enabledPoints.length; i++) {
    const p = enabledPoints[i];
    const m = mappings[p.id];
    if (!m) continue;

    // 如果在第一階段已經分配好，就跳過
    if (out[p.id]) {
      continue;
    }

    const tagShort = m.tag_key.split('.').pop() ?? `col_${i + 1}`;
    // 支援非 Modbus 點位帶有 r 前綴 (如 rd0) 時與資料庫欄位 (如 sensor_d0、d0) 之容錯匹配
    const tagShortAlt = tagShort.startsWith('r') && tagShort.length > 1 ? tagShort.slice(1) : '';
    
    // 1. 精確/前後綴匹配
    const exact = columnNames.find(c =>
      !used.has(c) && (matchesColumn(c, tagShort) || matchesColumn(c, tagShortAlt))
    );

    let chosenColumn = exact;

    if (!chosenColumn && columnNames.length > 0) {
      // 2. 索引回退 (index fallback)
      const fallback = columnNames[i % columnNames.length];
      if (!used.has(fallback)) {
        chosenColumn = fallback;
      } else {
        // 3. 第一個未使用的欄位
        const firstUnused = columnNames.find(c => !used.has(c));
        chosenColumn = firstUnused ?? fallback;
      }
    }

    const auto = chosenColumn ?? `col_${i + 1}`;
    used.add(auto);
    out[p.id] = {
      tag_id: `tag.${m.tag_key}`,
      column_name: auto,
      enabled: true
    };
  }

  return out;
}
