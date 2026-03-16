import { normalizeNamingPrefix } from '../../../features/datalink/sourcePlannerContract';
import type { CreateTagRequest, Mapping, Point, Tag } from '../../../types/datalink';

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
  const linkedPointIds = new Set(input.mappings.map((mapping) => mapping.point_id));

  return input.points.map((point, index) => {
    const previewKey = previewKeys[index];
    const normalizedKey = previewKey.toLowerCase();
    const duplicatePreview = (previewKeyCounts.get(normalizedKey) ?? 0) > 1;
    const existingKey = existingKeys.has(normalizedKey);

    return {
      pointId: point.id,
      pointName: point.name,
      pointAddress: point.address,
      dataType: point.data_type,
      previewKey,
      conflict: duplicatePreview || existingKey,
      conflictReason: existingKey
        ? 'existing-key'
        : duplicatePreview
          ? 'duplicate-preview'
          : null,
      alreadyLinked: linkedPointIds.has(point.id),
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
