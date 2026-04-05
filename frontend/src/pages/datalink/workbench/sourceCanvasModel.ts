import type {
  DataType,
  Mapping,
  Point,
  ProtocolType,
  Tag,
} from '../../../types/datalink';
import { addressParser } from '../../../utils/addressParser';

export type SourceViewMode = 'plan' | 'live' | 'link';

export type SourceValueFormat = 'decimal' | 'hex' | 'binary' | 'float';

export type SourceRule = {
  id: string;
  deviceId?: string;
  startAddress: string;
  count: number;
  dataType: DataType;
  namingPrefix: string;
  enabled: boolean;
  locked: boolean;
  origin: 'manual' | 'template';
  templateName?: string;
  skippedAddresses: string[];
  persisted?: boolean;
  updatedAt?: string;
  /** 目標資料型態（可空）。預設與 dataType 相同。 */
  targetDataType?: DataType;
  /** 縮放倍率（可空）。 */
  scaleMultiplier?: number;
  /** 偏移量（可空）。 */
  scaleOffset?: number;
  /** 字節序格式（可空）。 */
  dataFormat?: string;
};

export type AddressCanvasStatus = 'gap' | 'planned' | 'used' | 'unmanaged' | 'conflict';

export type AddressLinkState = 'needsPoint' | 'unbound' | 'draft' | 'ready' | 'blocked';

export type AddressCanvasItem = {
  address: string;
  status: AddressCanvasStatus;
  point?: Point;
  tagDisplayName?: string | null;
  tagKey?: string | null;
  ruleIds: ReadonlyArray<string>;
  primaryRuleId: string | null;
  mergeSpan: number;
  mergeOffset: number;
  ruleMergeOffset: number;
  linkState: AddressLinkState | null;
  linkLabelKey: string | null;
  liveValue: unknown;
  liveTimestamp: string | null;
};

export type CoverageOverviewSegment = {
  id: string;
  status: AddressCanvasStatus;
  startAddress: string;
  endAddress: string;
  cellCount: number;
};

type LiveValueSnapshot = {
  raw_value?: unknown;
  timestamp?: string;
};

type BuildAddressCanvasItemsInput =
  | {
      points: Point[];
      rules: ReadonlyArray<SourceRule>;
      mappings?: Mapping[];
      tags?: Tag[];
      protocol: ProtocolType;
      liveValues?: Readonly<Record<string, LiveValueSnapshot>>;
    }
  | {
      points: Point[];
      plannedPointAddresses: string[];
      plannedDataType: DataType;
      mappings?: Mapping[];
      tags?: Tag[];
      protocol: ProtocolType;
      liveValues?: Readonly<Record<string, LiveValueSnapshot>>;
    };

type PointOccupancy = {
  point: Point;
  mergeSpan: number;
  mergeOffset: number;
};

type RuleOccupancy = {
  ruleIds: string[];
  primaryRuleId: string;
  mergeSpan: number;
  mergeOffset: number;
};

const DATA_TYPE_CELL_SPAN: Record<DataType, number> = {
  bool: 1,
  int16: 1,
  uint16: 1,
  int32: 2,
  uint32: 2,
  float32: 2,
  int64: 4,
  uint64: 4,
  float64: 4,
  string: 1,
};

export function getDataTypeCellSpan(dataType: DataType): number {
  return DATA_TYPE_CELL_SPAN[dataType] ?? 1;
}

export function getDataTypeBitWidth(dataType: DataType): number {
  return getDataTypeCellSpan(dataType) * 16;
}

export function buildPlannedPointAddresses(input: {
  startAddress: string;
  count: number;
  dataType: DataType;
  protocol: ProtocolType;
}) {
  const span = getDataTypeCellSpan(input.dataType);
  const logicalCount = Math.max(0, input.count);

  return Array.from({ length: logicalCount }, (_, index) =>
    addressParser.offset(input.startAddress, index * span, input.protocol),
  );
}

export function buildSourceRuleCoverage(rule: SourceRule, protocol: ProtocolType) {
  const occupiedAddresses = expandRuleOccupiedAddresses(rule, protocol);
  const startAddress = occupiedAddresses[0] ?? rule.startAddress;
  const endAddress = occupiedAddresses.at(-1) ?? rule.startAddress;

  return {
    startAddress,
    endAddress,
    cellCount: occupiedAddresses.length,
    bitWidth: getDataTypeBitWidth(rule.dataType),
  };
}

export function buildAddressCanvasItems(input: BuildAddressCanvasItemsInput): AddressCanvasItem[] {
  const pointOccupancy = buildPointOccupancyMap(input.points, input.protocol);
  const ruleOccupancy =
    'rules' in input
      ? buildRuleOccupancyMap(input.rules, input.protocol)
      : buildLegacyRuleOccupancyMap(
          input.plannedPointAddresses,
          input.plannedDataType,
          input.protocol,
        );
  const mappings = input.mappings ?? [];
  const tags = input.tags ?? [];
  const primaryTagByPointId = new Map<string, Tag>();
  for (const mapping of mappings) {
    if (primaryTagByPointId.has(mapping.point_id)) {
      continue;
    }
    const tag = tags.find((candidate) => candidate.id === mapping.tag_id);
    if (tag) {
      primaryTagByPointId.set(mapping.point_id, tag);
    }
  }
  const occupiedAddresses = new Set<string>([
    ...pointOccupancy.keys(),
    ...ruleOccupancy.keys(),
  ]);

  if (occupiedAddresses.size === 0) {
    return [];
  }

  const sortedOccupiedAddresses = [...occupiedAddresses].sort((left, right) =>
    sortAddresses(left, right, input.protocol),
  );
  const minAddress = sortedOccupiedAddresses[0]!;
  const maxAddress = sortedOccupiedAddresses[sortedOccupiedAddresses.length - 1]!;
  const cellCount = getSequentialCellCount(minAddress, maxAddress, input.protocol);

  const rulesList = 'rules' in input ? input.rules : [];
  const ruleById = new Map(rulesList.map((r) => [r.id, r] as const));

  return Array.from({ length: cellCount }, (_, offset) => {
    const address = addressParser.offset(minAddress, offset, input.protocol);
    const pointMeta = pointOccupancy.get(address);
    const ruleMeta = ruleOccupancy.get(address);
    const primaryTag = pointMeta?.point ? primaryTagByPointId.get(pointMeta.point.id) : undefined;
    const hasRuleConflict = (ruleMeta?.ruleIds.length ?? 0) > 1;

    let status: AddressCanvasStatus = 'gap';
    if (hasRuleConflict) {
      status = 'conflict';
    } else if (pointMeta && ruleMeta) {
      // 若 point 的合併特徵與 rule 一致，表示 point 是由該 rule 建立的 → used
      // 若不一致（如既有 uint16 與 float32 continuation 重疊），則為真正的衝突
      const mergeMatch =
        pointMeta.mergeSpan === ruleMeta.mergeSpan &&
        pointMeta.mergeOffset === ruleMeta.mergeOffset;
      status = mergeMatch ? 'used' : 'conflict';
    } else if (pointMeta) {
      status = 'unmanaged';
    } else if (ruleMeta) {
      status = 'planned';
    }

    const link = resolveLinkState({
      point: pointMeta?.point,
      hasRule: Boolean(ruleMeta),
      mappings,
      tags,
      hasConflict: status === 'conflict',
    });
    const liveSnapshot = pointMeta?.point ? input.liveValues?.[pointMeta.point.id] : undefined;

    let liveValue: unknown = liveSnapshot?.raw_value ?? pointMeta?.point?.last_value;
    const primaryId = ruleMeta?.primaryRuleId ?? null;
    if (primaryId) {
      const ruleForCell = ruleById.get(primaryId);
      liveValue = applyRuleLinearDisplayScale(liveValue, ruleForCell);
    }

    return {
      address,
      status,
      point: pointMeta?.point,
      tagDisplayName: primaryTag?.display_name ?? null,
      tagKey: primaryTag?.key ?? null,
      ruleIds: ruleMeta?.ruleIds ?? [],
      primaryRuleId: ruleMeta?.primaryRuleId ?? null,
      mergeSpan: pointMeta?.mergeSpan ?? ruleMeta?.mergeSpan ?? 1,
      mergeOffset: pointMeta?.mergeOffset ?? ruleMeta?.mergeOffset ?? 0,
      ruleMergeOffset: ruleMeta?.mergeOffset ?? 0,
      linkState: link.state,
      linkLabelKey: link.labelKey,
      liveValue,
      liveTimestamp:
        liveSnapshot?.timestamp ?? pointMeta?.point?.last_read_at ?? null,
    };
  });
}

export function buildCoverageOverviewSegments(
  items: ReadonlyArray<AddressCanvasItem>,
): CoverageOverviewSegment[] {
  const segments: CoverageOverviewSegment[] = [];

  for (const item of items) {
    const lastSegment = segments.at(-1);
    if (!lastSegment || lastSegment.status !== item.status) {
      segments.push({
        id: `${item.status}-${item.address}`,
        status: item.status,
        startAddress: item.address,
        endAddress: item.address,
        cellCount: 1,
      });
      continue;
    }

    lastSegment.endAddress = item.address;
    lastSegment.cellCount += 1;
  }

  return segments;
}

/**
 * 將即時讀值轉成有限數字，供畫布縮放顯示；無法解析時回傳 `null`。
 *
 * @param value - point.last_value、SSE `raw_value` 或 JSON 字串
 * @returns 可參與線性縮放的數字，否則 `null`
 */
function parseLiveNumericForCanvas(value: unknown): number | null {
  if (value === null || value === undefined || value === '') {
    return null;
  }
  if (typeof value === 'number' && Number.isFinite(value)) {
    return value;
  }
  if (typeof value === 'string') {
    const trimmed = value.trim();
    if (trimmed === '') {
      return null;
    }
    const direct = Number(trimmed);
    if (Number.isFinite(direct)) {
      return direct;
    }
    try {
      const parsed: unknown = JSON.parse(trimmed);
      if (typeof parsed === 'number' && Number.isFinite(parsed)) {
        return parsed;
      }
    } catch {
      return null;
    }
  }
  return null;
}

/**
 * 套用來源規則的線性縮放於畫布顯示（與後端 mapping `scale` 步驟一致；point 仍存原始採集值）。
 *
 * @param value - 顯示前的即時值
 * @param rule - 該格所屬之主要規則（可空）
 * @returns 縮放後數字，或無法／無需縮放時回傳原值
 */
function applyRuleLinearDisplayScale(value: unknown, rule: SourceRule | undefined): unknown {
  if (!rule || (rule.scaleMultiplier === undefined && rule.scaleOffset === undefined)) {
    return value;
  }
  const n = parseLiveNumericForCanvas(value);
  if (n === null) {
    return value;
  }
  const scale = rule.scaleMultiplier ?? 1;
  const offset = rule.scaleOffset ?? 0;
  return n * scale + offset;
}

/**
 * 格式化來源值以供畫布顯示；處理多種後端回傳格式。
 *
 * 支援格式：
 * - number：直接格式化
 * - boolean：顯示 true/false
 * - string：嘗試解析 JSON（後端可能傳 `{"value":...,"raw_bytes":...}`）
 * - array：取第一個元素再格式化（一格一個邏輯值；MC／部分驅動可能回傳單元素或多餘欄位陣列）
 * - object：嘗試提取 `.value` 或轉成 JSON
 *
 * @param value - 來源值（來自 SSE `raw_value` 或 `point.last_value`）
 * @param format - 數值格式化模式
 * @returns 格式化後的字串，供格位覆蓋層顯示
 */
export function formatSourceValue(value: unknown, format: SourceValueFormat): string {
  if (value === null || value === undefined || value === '') {
    return '—';
  }

  if (typeof value === 'number') {
    if (!Number.isFinite(value)) {
      return 'NaN';
    }
    switch (format) {
      case 'hex':
        return `0x${Math.trunc(value).toString(16).toUpperCase()}`;
      case 'binary':
        return `0b${Math.trunc(value).toString(2)}`;
      case 'float':
        return value.toFixed(3).replace(/\.?0+$/, '');
      case 'decimal':
        return String(value);
    }
  }

  if (typeof value === 'boolean') {
    return value ? 'true' : 'false';
  }

  if (Array.isArray(value)) {
    if (value.length === 0) {
      return '—';
    }
    return formatSourceValue(value[0], format);
  }

  if (typeof value === 'object') {
    const record = value as Record<string, unknown>;
    if ('value' in record) {
      return formatSourceValue(record.value, format);
    }
    try {
      return JSON.stringify(value);
    } catch {
      return String(value);
    }
  }

  if (typeof value === 'string') {
    const trimmed = value.trim();
    if (trimmed === '') {
      return '—';
    }

    const direct = Number(trimmed);
    if (Number.isFinite(direct)) {
      return formatSourceValue(direct, format);
    }

    if (trimmed.startsWith('{') || trimmed.startsWith('[')) {
      try {
        const parsed: unknown = JSON.parse(trimmed);
        return formatSourceValue(parsed, format);
      } catch {
        // 無法解析，回傳原始字串
      }
    }

    return trimmed;
  }

  return String(value);
}

function buildPointOccupancyMap(points: Point[], protocol: ProtocolType) {
  const occupancy = new Map<string, PointOccupancy>();

  for (const point of points) {
    const occupiedAddresses = expandOccupiedAddresses(
      point.address,
      point.data_type,
      protocol,
    );

    occupiedAddresses.forEach((address, index) => {
      occupancy.set(address, {
        point,
        mergeSpan: occupiedAddresses.length,
        mergeOffset: index,
      });
    });
  }

  return occupancy;
}

function buildRuleOccupancyMap(rules: ReadonlyArray<SourceRule>, protocol: ProtocolType) {
  const occupancy = new Map<string, RuleOccupancy>();

  for (const rule of rules.filter((candidate) => candidate.enabled)) {
    const pointAddresses = buildPlannedPointAddresses({
      startAddress: rule.startAddress,
      count: rule.count,
      dataType: rule.dataType,
      protocol,
    });

    for (const pointAddress of pointAddresses) {
      if (rule.skippedAddresses?.includes(pointAddress)) continue;
      const occupiedAddresses = expandOccupiedAddresses(
        pointAddress,
        rule.dataType,
        protocol,
      );

      occupiedAddresses.forEach((address, index) => {
        const current = occupancy.get(address);
        if (current) {
          occupancy.set(address, {
            ...current,
            ruleIds: [...current.ruleIds, rule.id],
          });
          return;
        }

        occupancy.set(address, {
          ruleIds: [rule.id],
          primaryRuleId: rule.id,
          mergeSpan: occupiedAddresses.length,
          mergeOffset: index,
        });
      });
    }
  }

  return occupancy;
}

function buildLegacyRuleOccupancyMap(
  plannedPointAddresses: ReadonlyArray<string>,
  plannedDataType: DataType,
  protocol: ProtocolType,
) {
  const occupancy = new Map<string, RuleOccupancy>();

  plannedPointAddresses.forEach((pointAddress, ruleIndex) => {
    const occupiedAddresses = expandOccupiedAddresses(
      pointAddress,
      plannedDataType,
      protocol,
    );
    const ruleId = `legacy-rule-${ruleIndex + 1}`;

    occupiedAddresses.forEach((address, index) => {
      occupancy.set(address, {
        ruleIds: [ruleId],
        primaryRuleId: ruleId,
        mergeSpan: occupiedAddresses.length,
        mergeOffset: index,
      });
    });
  });

  return occupancy;
}

function expandRuleOccupiedAddresses(rule: SourceRule, protocol: ProtocolType) {
  return buildPlannedPointAddresses({
    startAddress: rule.startAddress,
    count: rule.count,
    dataType: rule.dataType,
    protocol,
  }).flatMap((pointAddress) =>
    expandOccupiedAddresses(pointAddress, rule.dataType, protocol),
  );
}

function expandOccupiedAddresses(
  startAddress: string,
  dataType: DataType,
  protocol: ProtocolType,
) {
  return addressParser.expand(startAddress, getDataTypeCellSpan(dataType), protocol);
}

function getSequentialCellCount(
  startAddress: string,
  endAddress: string,
  protocol: ProtocolType,
) {
  try {
    return (
      addressParser.parse(endAddress, protocol).startNumber -
        addressParser.parse(startAddress, protocol).startNumber +
      1
    );
  } catch {
    return 1;
  }
}

function sortAddresses(left: string, right: string, protocol: ProtocolType) {
  try {
    return (
      addressParser.parse(left, protocol).startNumber -
      addressParser.parse(right, protocol).startNumber
    );
  } catch {
    return left.localeCompare(right);
  }
}

export type ConflictQueueItem = {
  id: string;
  address: string;
  conflictCellAddress: string;
  ruleIds: ReadonlyArray<string>;
  reason: 'rule-overlap' | 'point-overlap';
  reasonKey: string;
};

export function buildConflictQueue(
  items: ReadonlyArray<AddressCanvasItem>,
): ConflictQueueItem[] {
  const queue: ConflictQueueItem[] = [];
  const emitted = new Set<string>();

  for (let i = 0; i < items.length; i++) {
    const item = items[i];
    if (item.status !== 'conflict') continue;

    // Resolve to root of the rule's logical span using rule-perspective offset
    const rootIndex = i - item.ruleMergeOffset;
    const rootAddress = rootIndex >= 0 ? items[rootIndex].address : item.address;

    if (emitted.has(rootAddress)) continue;
    emitted.add(rootAddress);

    const isRuleOverlap = item.ruleIds.length > 1;

    queue.push({
      id: `conflict-${rootAddress}`,
      address: rootAddress,
      conflictCellAddress: item.address,
      ruleIds: item.ruleIds,
      reason: isRuleOverlap ? 'rule-overlap' : 'point-overlap',
      reasonKey: isRuleOverlap
        ? 'workbench.source.conflictQueue.ruleOverlap'
        : 'workbench.source.conflictQueue.pointOverlap',
    });
  }

  return queue;
}

/**
 * Count eligible (ready-to-create) logical spans that Step 2 reports.
 *
 * Mirrors the exact `readyToCreateCount` semantics from SourceCanvasSection:
 *   1. Build the full address canvas from rules + existing points.
 *   2. Keep only items with `status === 'planned'` (excludes conflicts,
 *      existing-point overlap, and gaps).
 *   3. Keep only logical-root cells (`mergeOffset === 0`), not continuation
 *      cells of wide data types.
 *   4. Exclude addresses belonging to locked (protected) rules, because
 *      Step 2 reports them separately.
 */
export function countEligibleSpans(input: {
  rules: ReadonlyArray<SourceRule>;
  points: Point[];
  protocol: ProtocolType;
}): number {
  const items = buildAddressCanvasItems({
    points: input.points,
    rules: input.rules,
    protocol: input.protocol,
  });

  const protectedAddresses = buildProtectedAddressSet(input.rules, input.protocol);

  return items.filter(
    (item) =>
      item.status === 'planned' &&
      item.mergeOffset === 0 &&
      !protectedAddresses.has(item.address),
  ).length;
}

/**
 * Collect all point-level addresses produced by locked (protected) rules.
 */
function buildProtectedAddressSet(
  rules: ReadonlyArray<SourceRule>,
  protocol: ProtocolType,
): Set<string> {
  const addresses = new Set<string>();

  for (const rule of rules.filter((r) => r.enabled && r.locked)) {
    const planned = buildPlannedPointAddresses({
      startAddress: rule.startAddress,
      count: rule.count,
      dataType: rule.dataType,
      protocol,
    });

    for (const address of planned) {
      if (!rule.skippedAddresses?.includes(address)) {
        addresses.add(address);
      }
    }
  }

  return addresses;
}

function resolveLinkState(input: {
  point: Point | undefined;
  hasRule: boolean;
  mappings: Mapping[];
  tags: Tag[];
  hasConflict: boolean;
}): { state: AddressLinkState | null; labelKey: string | null } {
  if (input.hasConflict) {
    return {
      state: 'blocked',
      labelKey: 'workbench.source.link.conflict',
    };
  }

  if (!input.point) {
    return input.hasRule
      ? {
          state: 'needsPoint',
          labelKey: 'workbench.source.link.needsPoint',
        }
      : {
          state: null,
          labelKey: null,
        };
  }

  const mapping = input.mappings.find((candidate) => candidate.point_id === input.point?.id);
  if (!mapping) {
    return {
      state: 'unbound',
      labelKey: 'workbench.source.link.unbound',
    };
  }

  const tag = input.tags.find((candidate) => candidate.id === mapping.tag_id);
  if (!tag) {
    return {
      state: 'blocked',
      labelKey: 'workbench.source.link.missingTag',
    };
  }

  switch (tag.status) {
    case 'active':
      return {
        state: 'ready',
        labelKey: 'workbench.source.link.ready',
      };
    case 'draft':
      return {
        state: 'draft',
        labelKey: 'workbench.source.link.draft',
      };
    case 'retired':
      return {
        state: 'blocked',
        labelKey: 'workbench.source.link.blocked',
      };
  }
}
