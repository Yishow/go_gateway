import { useMemo } from 'react';
import type { Rule, Point, ShareLayout, Mapping, Device } from './types';
import {
  deriveAllPoints,
  computeShareLayout,
  detectAddressConflicts,
  derivePoints,
} from './sourceRule';

import type { ProtocolType } from '../../../../types/datalink';

const EMPTY_PROTOCOL_MAP: Record<string, ProtocolType> = Object.freeze({});

/**
 * 依設備清單建立設備 ID 至通訊協議對照表 (附帶 caching)
 *
 * @param devices 設備清單
 * @returns 鍵為 device.id，值為 protocol 的對照表
 */
export function useDeviceProtocolMap(devices: Device[]): Record<string, ProtocolType> {
  return useMemo(() => {
    const map: Record<string, ProtocolType> = {};
    devices.forEach((d) => {
      map[d.id] = d.protocol;
    });
    return map;
  }, [devices]);
}

/**
 * 衍生所有規則所對應的點位列表 (附帶 caching)
 *
 * @param rules 規則陣列
 * @param fallbackDeviceId 備用設備 ID
 * @param deviceProtocolMap 設備 ID 與協議對照表 (選填)
 * @returns 點位陣列
 */
export function useAllPoints(
  rules: Rule[],
  fallbackDeviceId: string,
  deviceProtocolMap: Record<string, ProtocolType> = EMPTY_PROTOCOL_MAP,
): Point[] {
  return useMemo(() => {
    return deriveAllPoints(rules, fallbackDeviceId, deviceProtocolMap);
  }, [rules, fallbackDeviceId, deviceProtocolMap]);
}

/**
 * 計算全域 Modbus Share 的位址分配 (附帶 caching)
 *
 * @param rules 規則陣列
 * @param baseRegister 記憶體起始暫存器號
 * @returns 鍵為 ruleId，值為 ShareLayout 的對照表
 */
export function useShareLayout(
  rules: Rule[],
  baseRegister: number
): Record<string, ShareLayout | null> {
  return useMemo(() => {
    return computeShareLayout(rules, baseRegister);
  }, [rules, baseRegister]);
}

/**
 * 檢測所有點位中是否有重複/衝突的 Modbus 地址 (附帶 caching)
 *
 * @param allPoints 所有點位陣列
 * @returns 衝突位址 Set
 */
export function useConflictAddrs(allPoints: Point[]): Set<string> {
  return useMemo(() => {
    return detectAddressConflicts(allPoints);
  }, [allPoints]);
}

/**
 * 根據單一規則衍生其下的所有點位資訊 (附帶 caching)
 *
 * @param rule 規則資料
 * @param deviceId 所屬設備 ID
 * @param skippedSet 已略過的位址集合 (Set)
 * @param protocol 設備通訊協議 (選填，預設 modbus_tcp)
 * @returns 點位陣列
 */
export function useRulePoints(
  rule: Rule | null | undefined,
  deviceId: string,
  skippedSet: Set<string>,
  protocol: ProtocolType = 'modbus_tcp',
): Point[] {
  return useMemo(() => {
    if (!rule) {
      return [];
    }
    return derivePoints(rule, deviceId, skippedSet, protocol);
  }, [rule, deviceId, skippedSet, protocol]);
}

/**
 * 點位映射的驗證 Hook，確認是否可以繼續下一步 (附帶 caching)
 *
 * @param mappings 映射對照表
 * @returns 驗證狀態物件
 */
export function useMappingValidation(
  mappings: Record<string, Mapping>
): { canContinue: boolean; emptyTagCount: number; enabledCount: number; totalCount: number } {
  return useMemo(() => {
    const keys = Object.keys(mappings);
    const totalCount = keys.length;
    let enabledCount = 0;
    let emptyTagCount = 0;

    keys.forEach((id) => {
      const m = mappings[id];
      if (m.enabled) {
        enabledCount++;
        if (!m.tag_key || m.tag_key.trim() === '') {
          emptyTagCount++;
        }
      }
    });

    const canContinue = enabledCount > 0 && emptyTagCount === 0;

    return {
      canContinue,
      emptyTagCount,
      enabledCount,
      totalCount,
    };
  }, [mappings]);
}

/**
 * 取得當前選取的點位與對應映射及模擬值 (附帶 caching)
 *
 * @param selectedIdx 選取的列索引
 * @param points 所有點位陣列
 * @param mappings 映射對照表
 * @returns 選定映射資訊物件或 null
 */
export function useSelectedMapping(
  selectedIdx: number | null,
  points: Point[],
  mappings: Record<string, Mapping>
): { point: Point; mapping: Mapping } | null {
  return useMemo(() => {
    if (selectedIdx === null || selectedIdx < 0 || selectedIdx >= points.length) {
      return null;
    }
    const point = points[selectedIdx];
    const mapping = mappings[point.id];
    if (!point || !mapping) {
      return null;
    }
    return {
      point,
      mapping,
    };
  }, [selectedIdx, points, mappings]);
}
