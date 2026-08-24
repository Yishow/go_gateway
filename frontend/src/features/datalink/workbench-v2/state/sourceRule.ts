import type { Rule, Point, ShareLayout, PointFunctionType } from './types';
import { addressParser } from '../../../../utils/addressParser';
import type { ProtocolType } from '../../../../types/datalink';

/** 可阻擋規則衍生或 Step 2 繼續的穩定原因碼。 */
export type RuleReadinessReason = 'unknown_device' | 'deleted_device' | 'invalid_address';
/** 以規則、設備與位址 identity 攜帶的 readiness 問題。 */
export interface RuleReadinessIssue {
  rule: Rule;
  ruleId: string;
  deviceId: string;
  startAddress: string;
  reason: RuleReadinessReason;
}
/**
 * 依資料型別取得其在暫存器中所佔用的寬度 (Stride / Register Count)
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
 * 根據位址字串與協議推斷功能碼或暫存器型態標籤 (通用函式)
 * 
 * @param addr 位址字串 (例如 "40001", "30001", "D0", "X0")
 * @param protocol 通訊協議 (預設 "modbus_tcp")
 * @returns 'coil' | 'discrete_input' | 'input_register' | 'holding_register' 或區域標籤
 */
function getPointFunctionType(addr: string, protocol: ProtocolType = 'modbus_tcp'): PointFunctionType {
  return addressParser.getAreaInfo(addr, protocol).typeLabel;
}

/**
 * 根據位址字串與協議推斷功能碼或暫存器型態標籤 (相容別名)
 */
export const fnFromAddr = getPointFunctionType;
/**
 * 格式化位址，將數值或字串正規化為標準字串
 * 
 * @param n 位址數值或字串
 * @returns 格式化後的位址字串
 */
export function formatAddr(n: number | string): string {
  return String(n);
}
/**
 * 判斷規則是否具備可衍生點位所需的設備 ownership 與合法位址。
 *
 * @param rule 待檢查規則
 * @param deviceProtocolMap 當前工作區的設備協議對照表
 * @returns 缺失原因；null 代表規則可依其設備協議衍生
 */
export function getRuleReadinessReason(
  rule: Rule,
  deviceProtocolMap: Record<string, ProtocolType>,
): RuleReadinessReason | null {
  const deviceId = rule.device_id?.trim() ?? '';
  if (!deviceId) {
    return 'unknown_device';
  }

  const protocol = deviceProtocolMap[deviceId];
  if (!protocol) {
    return rule.persisted ? 'deleted_device' : 'unknown_device';
  }

  return addressParser.validate(rule.start_address, protocol).valid
    ? null
    : 'invalid_address';
}
/**
 * 集中取得所有啟用規則的 readiness 問題，供 Step 2 的衍生與 UI gate 共用。
 */
export function getRuleReadinessIssues(
  rules: Rule[],
  deviceProtocolMap: Record<string, ProtocolType>,
): RuleReadinessIssue[] {
  return rules.flatMap((rule) => {
    if (!rule.enabled) {
      return [];
    }

    const reason = getRuleReadinessReason(rule, deviceProtocolMap);
    return reason
      ? [{
        rule,
        ruleId: rule.id,
        deviceId: rule.device_id,
        startAddress: rule.start_address,
        reason,
      }]
      : [];
  });
}
/**
 * 判斷 Step 2 是否可交接至映射流程。
 * 每個啟用規則都必須擁有當前設備、合法位址與至少一個未略過點位。
 */
export function isStep2Ready(
  rules: Rule[],
  deviceProtocolMap: Record<string, ProtocolType>,
): boolean {
  const enabledRules = rules.filter((rule) => rule.enabled);
  if (enabledRules.length === 0) {
    return false;
  }

  const points = deriveAllPoints(rules, deviceProtocolMap);
  const enabledPointRuleIds = new Set(
    points
      .filter((point) => point.enabled && !point.skipped)
      .map((point) => point.rule_id),
  );

  return enabledRules.every((rule) => (
    getRuleReadinessReason(rule, deviceProtocolMap) === null &&
    enabledPointRuleIds.has(rule.id)
  ));
}

/**
 * 根據規則與協議衍生出其底下的所有點位資訊
 * 
 * @param rule 規則資料
 * @param deviceId 所屬設備 ID
 * @param skippedSet 已略過的位址集合 (Set)
 * @param protocol 設備通訊協議 (預設 'modbus_tcp')
 * @returns 衍生出來的點位陣列
 */
export function derivePoints(
  rule: Rule,
  deviceId: string,
  skippedSet: Set<string>,
  protocol: ProtocolType = 'modbus_tcp',
): Point[] {
  const points: Point[] = [];
  const startAddr = rule.start_address?.trim() || '';
  if (!startAddr || !addressParser.validate(startAddr, protocol).valid) {
    return points;
  }

  const stride = dataTypeWidth(rule.data_type);
  try {
    const fnType = fnFromAddr(startAddr, protocol);

    for (let i = 0; i < rule.count; i++) {
      const addrStr = addressParser.offset(startAddr, i * stride, protocol);
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
  } catch {
    return [];
  }

  return points;
}

/**
 * 根據多個規則，衍生出整個工作區的全部點位資訊
 * 
 * @param rules 規則列表
 * @param deviceProtocolMap 設備 ID 至 Protocol 對照表
 * @returns 合併後的點位陣列
 */
export function deriveAllPoints(
  rules: Rule[],
  deviceProtocolMap: Record<string, ProtocolType>,
): Point[] {
  if (rules.length > 0 && Object.keys(deviceProtocolMap).length === 0) {
    return [];
  }

  const allPoints: Point[] = [];
  for (const rule of rules) {
    const devId = rule.device_id;
    if (!devId) {
      continue;
    }

    const readinessReason = getRuleReadinessReason(rule, deviceProtocolMap);
    if (readinessReason) {
      continue;
    }

    const skippedSet = new Set(rule.skipped_addresses || []);
    const protocol = deviceProtocolMap[devId];
    const points = derivePoints(rule, devId, skippedSet, protocol);
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
