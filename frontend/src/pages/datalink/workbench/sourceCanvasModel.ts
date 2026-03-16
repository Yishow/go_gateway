import type { DataType, Point, ProtocolType } from '../../../types/datalink';
import { addressParser } from '../../../utils/addressParser';

export type SourceViewMode = 'grid' | 'table';

export type AddressCanvasItem = {
  address: string;
  status: 'planned' | 'used' | 'conflict';
  point?: Point;
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

export function buildAddressCanvasItems(input: {
  points: Point[];
  plannedPointAddresses: string[];
  plannedDataType: DataType;
  protocol: ProtocolType;
}): AddressCanvasItem[] {
  const itemMap = new Map<string, AddressCanvasItem>();
  const usedAddressMap = buildUsedAddressMap(input.points, input.protocol);
  const conflictAddresses = new Set(
    getConflictingPlannedPointAddresses({
      points: input.points,
      plannedPointAddresses: input.plannedPointAddresses,
      plannedDataType: input.plannedDataType,
      protocol: input.protocol,
    }),
  );

  for (const [address, point] of usedAddressMap.entries()) {
    itemMap.set(address, {
      address,
      status: 'used',
      point,
    });
  }

  for (const pointAddress of input.plannedPointAddresses) {
    for (const address of expandOccupiedAddresses(
      pointAddress,
      input.plannedDataType,
      input.protocol,
    )) {
      itemMap.set(address, {
        address,
        status: conflictAddresses.has(pointAddress) && usedAddressMap.has(address)
          ? 'conflict'
          : 'planned',
        point: usedAddressMap.get(address),
      });
    }
  }

  return [...itemMap.values()].sort((left, right) =>
    sortAddresses(left.address, right.address, input.protocol),
  );
}

export function getConflictingPlannedPointAddresses(input: {
  points: Point[];
  plannedPointAddresses: string[];
  plannedDataType: DataType;
  protocol: ProtocolType;
}) {
  const usedAddressMap = buildUsedAddressMap(input.points, input.protocol);

  return input.plannedPointAddresses.filter((pointAddress) =>
    expandOccupiedAddresses(pointAddress, input.plannedDataType, input.protocol).some((address) =>
      usedAddressMap.has(address),
    ),
  );
}

function buildUsedAddressMap(points: Point[], protocol: ProtocolType) {
  const usedAddressMap = new Map<string, Point>();

  for (const point of points) {
    for (const address of expandOccupiedAddresses(point.address, point.data_type, protocol)) {
      usedAddressMap.set(address, point);
    }
  }

  return usedAddressMap;
}

function expandOccupiedAddresses(
  startAddress: string,
  dataType: DataType,
  protocol: ProtocolType,
) {
  return addressParser.expand(
    startAddress,
    getDataTypeCellSpan(dataType),
    protocol,
  );
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
