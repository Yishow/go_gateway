import type { Rule, Point, ShareLayout } from './types';

/**
 * 依資料型別取得其在 Modbus 中所佔用的暫存器寬度 (Stride / Register Count)
 * 
 * @param type 資料型別
 * @returns 暫存器數量 (1, 2, 4 或 10)
 */
export function dataTypeWidth(type: string): number {
  switch (type) {
    case 'int32':
    case 'uint32':
    case 'float32':
      return 2;
    case 'int64':
    case 'uint64':
    case 'float64':
      return 4;
    case 'string':
      return 10;
    case 'bool':
    case 'int16':
    case 'uint16':
    default:
      return 1;
  }
}

/**
 * 根據位址字串推斷 Modbus Function Code
 * 
 * @param addr 位址字串 (例如 "40001", "30001", "00001")
 * @returns 'coil' | 'discrete_input' | 'input_register' | 'holding_register'
 */
export function fnFromAddr(addr: string): 'coil' | 'discrete_input' | 'input_register' | 'holding_register' {
  const cleanAddr = addr.trim();
  if (cleanAddr.startsWith('0')) {
    return 'coil';
  }
  if (cleanAddr.startsWith('1')) {
    return 'discrete_input';
  }
  if (cleanAddr.startsWith('3')) {
    return 'input_register';
  }
  return 'holding_register'; // 預設或以 4 開頭者
}

/**
 * 格式化位址，將數字轉換為字串
 * 
 * @param n 位址數值或字串
 * @returns 格式化後的位址字串
 */
export function formatAddr(n: number | string): string {
  return String(n);
}

/**
 * 根據規則衍生出其底下的所有點位資訊
 * 
 * @param rule 規則資料
 * @param deviceId 所屬設備 ID
 * @param skippedSet 已略過的位址集合 (Set)
 * @returns 衍生出來的點位陣列
 */
export function derivePoints(rule: Rule, deviceId: string, skippedSet: Set<string>): Point[] {
  const points: Point[] = [];
  const startNum = parseInt(rule.start_address.replace(/[^0-9]/g, ''), 10) || 40001;
  const stride = dataTypeWidth(rule.data_type);
  const fnType = fnFromAddr(rule.start_address);

  for (let i = 0; i < rule.count; i++) {
    const addrNum = startNum + i * stride;
    const addrStr = formatAddr(addrNum);
    const isSkipped = skippedSet.has(addrStr);

    points.push({
      id: `${rule.id}-p-${i}`,
      device_id: deviceId,
      rule_id: rule.id,
      rule_name: rule.name,
      name: `${rule.naming_prefix}${i}`,
      address: addrStr,
      data_type: rule.data_type,
      function: fnType,
      width: stride,
      enabled: rule.enabled && !isSkipped,
      skipped: isSkipped,
      _rule_scale: rule.scale_multiplier,
      _rule_offset: rule.scale_offset,
    });
  }

  return points;
}

/**
 * 根據多個規則，衍生出整個工作區的全部點位資訊
 * 
 * @param rules 規則列表
 * @param fallbackDeviceId 備用設備 ID (若 rule 的 device_id 為空時使用)
 * @returns 合併後的點位陣列
 */
export function deriveAllPoints(rules: Rule[], fallbackDeviceId: string): Point[] {
  const allPoints: Point[] = [];
  for (const rule of rules) {
    const devId = rule.device_id || fallbackDeviceId;
    const skippedSet = new Set(rule.skipped_addresses || []);
    const points = derivePoints(rule, devId, skippedSet);
    allPoints.push(...points);
  }
  return allPoints;
}

/**
 * 計算所有規則在 Modbus Share 記憶體區段中的位址佈局
 * 
 * 落地設計決策：「Modbus Share 位址布局：自動接續 + 手動覆寫」
 * 
 * @param rules 規則列表
 * @param baseRegister 記憶體起始暫存器號
 * @returns 鍵為 rule.id，值為 ShareLayout 或 null 的對照表
 */
export function computeShareLayout(rules: Rule[], baseRegister: number): Record<string, ShareLayout | null> {
  const layout: Record<string, ShareLayout | null> = {};
  let cursor = baseRegister;

  for (const rule of rules) {
    if (!rule.share_enabled) {
      layout[rule.id] = null;
      continue;
    }

    const stride = rule.share_stride ?? dataTypeWidth(rule.data_type);
    const enabledCount = rule.count - (rule.skipped_addresses?.length ?? 0);

    let start: number;
    let isAuto = true;

    if (rule.share_start_register != null) {
      start = rule.share_start_register;
      isAuto = false;
    } else {
      start = cursor;
      isAuto = true;
    }

    const end = start + enabledCount * stride;
    layout[rule.id] = {
      start,
      stride,
      end,
      auto: isAuto,
    };

    // cursor 取最大值，防止手動指定的起點落後時把後續自動分配往回拉
    cursor = Math.max(cursor, end);
  }

  return layout;
}

/**
 * 檢測多個點位是否存在跨規則的位址衝突
 * 
 * @param allPoints 全體點位陣列
 * @returns 存在衝突的位址集合 (Set)
 */
export function detectAddressConflicts(allPoints: Point[]): Set<string> {
  const usage: Record<string, string[]> = {};

  for (const p of allPoints) {
    if (p.skipped || !p.enabled) continue;
    (usage[p.address] ??= []).push(p.rule_id);
  }

  const conflictAddrs = new Set<string>();
  Object.entries(usage).forEach(([addr, ruleIds]) => {
    // 跨規則且有多個引用時判定為衝突
    if (ruleIds.length > 1) {
      conflictAddrs.add(addr);
    }
  });

  return conflictAddrs;
}
