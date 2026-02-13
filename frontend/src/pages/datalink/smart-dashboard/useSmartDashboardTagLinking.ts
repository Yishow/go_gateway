import { useCallback, useEffect, useState } from 'react';
import type { DataType, Mapping, Point, Tag, UpdateTagRequest } from '../../../types/datalink';
import {
  buildGlobalTagEditDraft,
  hasGlobalTagEditChanges,
  toTagUpdateRequest,
} from '../../../features/datalink/tagEditImpact';

type TagLinkMode = 'existing' | 'create';

interface UseSmartDashboardTagLinkingParams {
  activePointForLink: Point | null;
  selectedMapping?: Mapping;
  tags: Tag[];
  linkedTag?: Tag;
  linkedTagAffectedMappingsCount: number;
  parsePipeline: () => unknown[];
  tagEditGuardrailWarning: string;
  createTag: (data: { key: string; display_name: string; data_type: DataType }) => Promise<Tag>;
  createMapping: (data: { point_id: string; tag_id: string; enabled: boolean }) => Promise<unknown>;
  updateMapping: (data: { id: string; data: { tag_id: string; enabled: boolean; transform_pipeline: unknown[] } }) => Promise<unknown>;
  updateTag: (data: { id: string; data: UpdateTagRequest }) => Promise<unknown>;
}

export function useSmartDashboardTagLinking({
  activePointForLink,
  selectedMapping,
  tags,
  linkedTag,
  linkedTagAffectedMappingsCount,
  parsePipeline,
  tagEditGuardrailWarning,
  createTag,
  createMapping,
  updateMapping,
  updateTag,
}: UseSmartDashboardTagLinkingParams) {
  const [tagLinkMode, setTagLinkMode] = useState<TagLinkMode>('existing');
  const [selectedTagIdForLink, setSelectedTagIdForLink] = useState('');
  const [newTagKey, setNewTagKey] = useState('');
  const [newTagDisplayName, setNewTagDisplayName] = useState('');
  const [tagLinkActionMessage, setTagLinkActionMessage] = useState('');
  const [tagEditDisplayName, setTagEditDisplayName] = useState('');
  const [tagEditUnit, setTagEditUnit] = useState('');
  const [tagEditDescription, setTagEditDescription] = useState('');
  const [tagEditMessage, setTagEditMessage] = useState('');
  const [pendingTagEdit, setPendingTagEdit] = useState<{
    display_name: string;
    unit: string;
    description: string;
  } | null>(null);

  useEffect(() => {
    setSelectedTagIdForLink(selectedMapping?.tag_id || '');
  }, [selectedMapping?.tag_id]);

  useEffect(() => {
    setTagEditDisplayName(linkedTag?.display_name || '');
    setTagEditUnit(linkedTag?.unit || '');
    setTagEditDescription(linkedTag?.description || '');
    setTagEditMessage('');
    setPendingTagEdit(null);
  }, [linkedTag?.description, linkedTag?.display_name, linkedTag?.id, linkedTag?.unit]);

  const handleLinkTagToSelectedAddress = useCallback(async () => {
    if (!activePointForLink) {
      setTagLinkActionMessage('請先選取已建立點位的格位，再進行 Tag 連結。');
      return;
    }

    const targetTagId = selectedTagIdForLink.trim();
    if (!targetTagId) {
      setTagLinkActionMessage('請先選擇要連結的既有 Tag。');
      return;
    }

    try {
      if (selectedMapping) {
        await updateMapping({
          id: selectedMapping.id,
          data: {
            tag_id: targetTagId,
            enabled: true,
            transform_pipeline: parsePipeline(),
          },
        });
        setTagLinkActionMessage(`已更新 ${activePointForLink.address} 的 Tag 連結。`);
        return;
      }

      await createMapping({
        point_id: activePointForLink.id,
        tag_id: targetTagId,
        enabled: true,
      });
      setTagLinkActionMessage(`已建立 ${activePointForLink.address} 的 Tag 連結。`);
    } catch (error) {
      const message = error instanceof Error ? error.message : '連結失敗';
      setTagLinkActionMessage(message);
    }
  }, [
    activePointForLink,
    createMapping,
    parsePipeline,
    selectedMapping,
    selectedTagIdForLink,
    updateMapping,
  ]);

  const handleCreateTagAndLink = useCallback(async () => {
    if (!activePointForLink) {
      setTagLinkActionMessage('請先選取已建立點位的格位，再建立 Tag。');
      return;
    }

    const normalizedKey = newTagKey.trim();
    if (!normalizedKey) {
      setTagLinkActionMessage('Tag Key 不可為空。');
      return;
    }

    const existing = tags.find((tag) => tag.key.toLowerCase() === normalizedKey.toLowerCase());
    if (existing) {
      setTagLinkActionMessage(`Tag Key ${normalizedKey} 已存在，請改用既有 Tag 模式。`);
      return;
    }

    try {
      const createdTag = await createTag({
        key: normalizedKey,
        display_name: newTagDisplayName.trim() || normalizedKey,
        data_type: activePointForLink.data_type,
      });

      if (selectedMapping) {
        await updateMapping({
          id: selectedMapping.id,
          data: {
            tag_id: createdTag.id,
            enabled: true,
            transform_pipeline: parsePipeline(),
          },
        });
      } else {
        await createMapping({
          point_id: activePointForLink.id,
          tag_id: createdTag.id,
          enabled: true,
        });
      }

      setSelectedTagIdForLink(createdTag.id);
      setTagLinkActionMessage(`已建立 Tag ${createdTag.key} 並完成連結。`);
      setNewTagKey('');
      setNewTagDisplayName('');
      setTagLinkMode('existing');
    } catch (error) {
      const message = error instanceof Error ? error.message : '建立 Tag 失敗';
      setTagLinkActionMessage(message);
    }
  }, [
    activePointForLink,
    createMapping,
    createTag,
    newTagDisplayName,
    newTagKey,
    parsePipeline,
    selectedMapping,
    tags,
    updateMapping,
  ]);

  const handleSaveLinkedTagEdit = useCallback(async () => {
    if (!linkedTag?.id) {
      setTagEditMessage('目前沒有可編輯的已連結 Tag。');
      return;
    }

    const nextEdit = buildGlobalTagEditDraft({
      display_name: tagEditDisplayName,
      unit: tagEditUnit,
      description: tagEditDescription,
    });
    if (!hasGlobalTagEditChanges(linkedTag, nextEdit)) {
      setTagEditMessage('沒有變更，無需儲存。');
      return;
    }

    setPendingTagEdit(nextEdit);
    setTagEditMessage(tagEditGuardrailWarning);
  }, [
    linkedTag,
    tagEditDescription,
    tagEditDisplayName,
    tagEditGuardrailWarning,
    tagEditUnit,
  ]);

  const handleConfirmTagEdit = useCallback(async () => {
    if (!linkedTag?.id || !pendingTagEdit) return;

    try {
      await updateTag({
        id: linkedTag.id,
        data: toTagUpdateRequest(pendingTagEdit),
      });
      setPendingTagEdit(null);
      setTagEditMessage(`已更新全域 Tag，影響 ${linkedTagAffectedMappingsCount} 個映射。`);
    } catch (error) {
      const message = error instanceof Error ? error.message : '更新 Tag 失敗';
      setTagEditMessage(message);
    }
  }, [linkedTag?.id, linkedTagAffectedMappingsCount, pendingTagEdit, updateTag]);

  const clearPendingTagEdit = useCallback(() => {
    setPendingTagEdit(null);
  }, []);

  return {
    tagLinkMode,
    setTagLinkMode,
    selectedTagIdForLink,
    setSelectedTagIdForLink,
    newTagKey,
    setNewTagKey,
    newTagDisplayName,
    setNewTagDisplayName,
    tagLinkActionMessage,
    tagEditDisplayName,
    setTagEditDisplayName,
    tagEditUnit,
    setTagEditUnit,
    tagEditDescription,
    setTagEditDescription,
    tagEditMessage,
    pendingTagEdit,
    handleLinkTagToSelectedAddress,
    handleCreateTagAndLink,
    handleSaveLinkedTagEdit,
    handleConfirmTagEdit,
    clearPendingTagEdit,
  };
}
