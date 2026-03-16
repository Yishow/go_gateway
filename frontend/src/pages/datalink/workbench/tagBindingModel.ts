import { normalizeNamingPrefix } from '../../../features/datalink/sourcePlannerContract';
import type { CreateTagRequest, Mapping, Point, Tag } from '../../../types/datalink';
import { getDataTypeBitWidth, getDataTypeCellSpan } from './sourceCanvasModel';

export type TagBindingStrategy = 'address' | 'pointName';

export interface TagBindingTemplate {
  prefix: string;
  strategy: TagBindingStrategy;
}

export interface TagBindingCandidate {
  pointId: string;
  pointName: string;
  pointAddress: string;
  dataType: Point['data_type'];
  bitWidth: number;
  cellSpan: number;
  rawValue: unknown;
  transformedValue: unknown;
  bindingStatus: 'bound' | 'unbound' | 'partial';
  existingTagOptions: Array<{
    id: string;
    key: string;
    displayName: string;
  }>;
  previewKey: string;
  conflict: boolean;
  conflictReason: 'existing-key' | 'duplicate-preview' | null;
  alreadyLinked: boolean;
}

export interface TagBindingRequest {
  pointId: string;
  tagKey: string;
  tagRequest: CreateTagRequest;
}

export type TagBindingFlowMode = 'create' | 'existing';

export interface TagBindingDiffEntry {
  pointId: string;
  pointName: string;
  tagKey: string;
}

export interface TagBindingSkippedEntry {
  pointId: string;
  pointName: string;
  reason: 'conflict' | 'already-linked' | 'no-tag-selected';
}

export interface TagBindingBatchDiffPreview {
  toCreate: TagBindingDiffEntry[];
  toBind: TagBindingDiffEntry[];
  skipped: TagBindingSkippedEntry[];
}

function normalizeTagKeySegment(value: string): string {
  const compact = value.trim().toUpperCase();
  const sanitized = compact
    .replace(/[^A-Z0-9_-]+/g, '_')
    .replace(/_+/g, '_')
    .replace(/^_|_$/g, '');

  return sanitized || 'VALUE';
}

export function buildTagKey(point: Point, template: TagBindingTemplate): string {
  const baseValue =
    template.strategy === 'pointName' ? point.name : point.address;

  return `${normalizeNamingPrefix(template.prefix)}_${normalizeTagKeySegment(baseValue)}`;
}

export function buildTagBindingCandidates(input: {
  points: Point[];
  tags: Tag[];
  mappings: Mapping[];
  template: TagBindingTemplate;
}): TagBindingCandidate[] {
  const previewKeys = input.points.map((point) => buildTagKey(point, input.template));
  const previewKeyCounts = previewKeys.reduce<Map<string, number>>((acc, previewKey) => {
    const normalizedKey = previewKey.toLowerCase();
    acc.set(normalizedKey, (acc.get(normalizedKey) ?? 0) + 1);
    return acc;
  }, new Map());
  const existingKeys = new Set(input.tags.map((tag) => tag.key.trim().toLowerCase()));
  const linkedMappingsByPointId = input.mappings.reduce<Map<string, Mapping>>((acc, mapping) => {
    if (!acc.has(mapping.point_id)) {
      acc.set(mapping.point_id, mapping);
    }
    return acc;
  }, new Map());

  return input.points.map((point, index) => {
    const previewKey = previewKeys[index];
    const normalizedKey = previewKey.toLowerCase();
    const duplicatePreview = (previewKeyCounts.get(normalizedKey) ?? 0) > 1;
    const existingKey = existingKeys.has(normalizedKey);
    const alreadyLinked = linkedMappingsByPointId.has(point.id);
    const bindingStatus = alreadyLinked
      ? 'bound'
      : duplicatePreview || existingKey
        ? 'partial'
        : 'unbound';

    return {
      pointId: point.id,
      pointName: point.name,
      pointAddress: point.address,
      dataType: point.data_type,
      bitWidth: getDataTypeBitWidth(point.data_type),
      cellSpan: getDataTypeCellSpan(point.data_type),
      rawValue: point.last_value,
      transformedValue: point.last_value,
      bindingStatus,
      existingTagOptions: input.tags
        .filter((tag) => tag.data_type === point.data_type)
        .map((tag) => ({
          id: tag.id,
          key: tag.key,
          displayName: tag.display_name || tag.key,
        })),
      previewKey,
      conflict: duplicatePreview || existingKey,
      conflictReason: existingKey
        ? 'existing-key'
        : duplicatePreview
          ? 'duplicate-preview'
          : null,
      alreadyLinked,
    };
  });
}

export function buildTagBindingRequests(
  candidates: TagBindingCandidate[],
  selectedPointIds: string[],
): TagBindingRequest[] {
  const selected = new Set(selectedPointIds);

  return candidates
    .filter(
      (candidate) =>
        selected.has(candidate.pointId) &&
        !candidate.conflict &&
        !candidate.alreadyLinked,
    )
    .map((candidate) => ({
      pointId: candidate.pointId,
      tagKey: candidate.previewKey,
      tagRequest: {
        key: candidate.previewKey,
        data_type: candidate.dataType,
        display_name: candidate.pointName,
      },
    }));
}

export function buildBatchDiffPreview(input: {
  candidates: TagBindingCandidate[];
  selectedPointIds: string[];
  flowMode: TagBindingFlowMode;
  existingTagSelections: Record<string, string>;
}): TagBindingBatchDiffPreview {
  const selected = new Set(input.selectedPointIds);

  return input.candidates.reduce<TagBindingBatchDiffPreview>(
    (acc, candidate) => {
      if (!selected.has(candidate.pointId)) {
        return acc;
      }

      if (input.flowMode === 'create') {
        if (candidate.alreadyLinked) {
          acc.skipped.push({
            pointId: candidate.pointId,
            pointName: candidate.pointName,
            reason: 'already-linked',
          });
          return acc;
        }

        if (candidate.conflict) {
          acc.skipped.push({
            pointId: candidate.pointId,
            pointName: candidate.pointName,
            reason: 'conflict',
          });
          return acc;
        }

        acc.toCreate.push({
          pointId: candidate.pointId,
          pointName: candidate.pointName,
          tagKey: candidate.previewKey,
        });
        return acc;
      }

      if (candidate.alreadyLinked) {
        acc.skipped.push({
          pointId: candidate.pointId,
          pointName: candidate.pointName,
          reason: 'already-linked',
        });
        return acc;
      }

      const selectedTagId = input.existingTagSelections[candidate.pointId];
      if (!selectedTagId) {
        acc.skipped.push({
          pointId: candidate.pointId,
          pointName: candidate.pointName,
          reason: 'no-tag-selected',
        });
        return acc;
      }

      const selectedTagKey =
        candidate.existingTagOptions.find((option) => option.id === selectedTagId)?.key
        ?? selectedTagId;

      acc.toBind.push({
        pointId: candidate.pointId,
        pointName: candidate.pointName,
        tagKey: selectedTagKey,
      });
      return acc;
    },
    {
      toCreate: [],
      toBind: [],
      skipped: [],
    },
  );
}
