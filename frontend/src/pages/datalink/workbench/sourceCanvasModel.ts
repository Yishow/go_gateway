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
  startAddress: string;
  count: number;
  dataType: DataType;
  namingPrefix: string;
  enabled: boolean;
  locked: boolean;
  origin: 'manual' | 'template';
  templateName?: string;
  skippedAddresses: string[];
};

export type AddressCanvasStatus = 'gap' | 'planned' | 'used' | 'conflict';

export type AddressLinkState = 'needsPoint' | 'unbound' | 'draft' | 'ready' | 'blocked';

export type AddressCanvasItem = {
  address: string;
  status: AddressCanvasStatus;
  point?: Point;
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

  return Array.from({ length: cellCount }, (_, offset) => {
    const address = addressParser.offset(minAddress, offset, input.protocol);
    const pointMeta = pointOccupancy.get(address);
    const ruleMeta = ruleOccupancy.get(address);
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
      status = 'used';
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

    return {
      address,
      status,
      point: pointMeta?.point,
      ruleIds: ruleMeta?.ruleIds ?? [],
      primaryRuleId: ruleMeta?.primaryRuleId ?? null,
      mergeSpan: pointMeta?.mergeSpan ?? ruleMeta?.mergeSpan ?? 1,
      mergeOffset: pointMeta?.mergeOffset ?? ruleMeta?.mergeOffset ?? 0,
      ruleMergeOffset: ruleMeta?.mergeOffset ?? 0,
      linkState: link.state,
      linkLabelKey: link.labelKey,
      liveValue: liveSnapshot?.raw_value ?? pointMeta?.point?.last_value,
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

export function formatSourceValue(value: unknown, format: SourceValueFormat): string {
  if (value === null || value === undefined || value === '') {
    return '—';
  }

  if (typeof value !== 'number' || !Number.isFinite(value)) {
    return String(value);
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
