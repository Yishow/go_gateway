import type { Mapping, Tag, UpdateTagRequest } from '../../types/datalink';

export interface GlobalTagEditDraft {
  display_name: string;
  unit: string;
  description: string;
}

export function getAffectedMappingsForTag(mappings: Mapping[], tagId: string): Mapping[] {
  return mappings.filter((mapping) => mapping.tag_id === tagId);
}

export function getAffectedMappingCountForTag(mappings: Mapping[], tagId: string): number {
  return getAffectedMappingsForTag(mappings, tagId).length;
}

export function buildGlobalTagEditDraft(
  input: Pick<GlobalTagEditDraft, 'display_name' | 'unit' | 'description'>
): GlobalTagEditDraft {
  return {
    display_name: input.display_name.trim(),
    unit: input.unit.trim(),
    description: input.description.trim(),
  };
}

export function hasGlobalTagEditChanges(tag: Tag, draft: GlobalTagEditDraft): boolean {
  return (
    draft.display_name !== (tag.display_name || '') ||
    draft.unit !== (tag.unit || '') ||
    draft.description !== (tag.description || '')
  );
}

export function toTagUpdateRequest(draft: GlobalTagEditDraft): UpdateTagRequest {
  return {
    display_name: draft.display_name,
    unit: draft.unit,
    description: draft.description,
  };
}
