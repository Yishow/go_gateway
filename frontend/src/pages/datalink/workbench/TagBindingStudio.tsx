import { useEffect, useMemo, useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { useDevicesQuery } from '../../../hooks/datalink/useDevices';
import { mappingKeys } from '../../../hooks/datalink/keys';
import {
  useCreateMappingMutation,
  useDeleteMappingMutation,
  useMappingsQuery,
} from '../../../hooks/datalink/useMappings';
import { useDeletePointMutation, usePointsQuery } from '../../../hooks/datalink/usePoints';
import {
  useCreateTagMutation,
  useDeleteTagMutation,
  useTagsQuery,
} from '../../../hooks/datalink/useTags';
import { tagAPI } from '../../../services/datalink';
import { DATALINK_DATA_TYPES, type DataType, type Device, type Tag } from '../../../types/datalink';
import { useWorkbench } from './WorkbenchProvider';
import { countEligibleSpans } from './sourceCanvasModel';
import {
  buildBatchDiffPreview,
  buildTagBindingCandidates,
  buildTagBindingRequests,
  type TagBindingCandidate,
  type TagBindingStrategy,
} from './tagBindingModel';

type TagBindingFailure = {
  pointId: string;
  tagKey: string;
  stage: 'tag' | 'mapping' | 'point';
  error: string;
};

type TagBindingBatchSummary = {
  mode: 'bind' | 'unbind' | 'deletePoints';
  createdCount: number;
  linkedCount: number;
  unboundCount: number;
  deletedPointCount: number;
  skippedCount: number;
  failureCount: number;
  failures: TagBindingFailure[];
};

type TagBindingFlowMode = 'create' | 'existing';
type TagBindingStatusFilter = 'all' | 'unbound' | 'partial' | 'bound';
type TagLibraryFeedback = {
  tone: 'success' | 'error';
  message: string;
};

const tagDataTypeOptions: DataType[] = [...DATALINK_DATA_TYPES];

function getSelectedDevice(devices: Device[], selectedDeviceId: string | null) {
  return devices.find((device) => device.id === selectedDeviceId) ?? null;
}

function getErrorMessage(error: unknown, fallback: string): string {
  return error instanceof Error ? error.message : fallback;
}

/**
 * 將候選點位即時值格式化為顯示字串；空值統一為 em dash。
 *
 * @param value 原始或轉換後數值
 * @returns 顯示用字串
 */
function formatCandidateValue(value: unknown) {
  if (value === null || value === undefined || value === '') {
    return '—';
  }

  return String(value);
}

/**
 * 候選列狀態徽章：柔和對比、內高光，與整體深色介面一致。
 *
 * @param status 綁定狀態
 * @returns Tailwind 類別字串
 */
function getStatusToneClass(status: TagBindingCandidate['bindingStatus']) {
  switch (status) {
    case 'bound':
      return 'border-emerald-400/20 bg-emerald-500/[0.14] text-emerald-100 shadow-[inset_0_1px_0_0_rgba(255,255,255,0.06)]';
    case 'partial':
      return 'border-amber-400/20 bg-amber-500/[0.14] text-amber-100 shadow-[inset_0_1px_0_0_rgba(255,255,255,0.06)]';
    case 'unbound':
      return 'border-slate-500/25 bg-slate-950/55 text-slate-300 shadow-[inset_0_1px_0_0_rgba(255,255,255,0.04)]';
  }
}

export function TagBindingStudio() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const { selectedDeviceId, setActiveStep, setFocusedTagIds, setInspectorSelection, setSelectedDeviceId, sourcePlanningState } = useWorkbench();
  const { data: devices = [] } = useDevicesQuery();
  const tagsQuery = useTagsQuery();
  const tagsData = tagsQuery.data;
  const tags = useMemo(() => tagsData ?? [], [tagsData]);
  const { data: mappings = [] } = useMappingsQuery();
  const { data: points = [] } = usePointsQuery(
    selectedDeviceId ? { device_id: selectedDeviceId } : undefined,
  );
  const createTagMutation = useCreateTagMutation();
  const deleteTagMutation = useDeleteTagMutation();
  const createMappingMutation = useCreateMappingMutation();
  const deleteMappingMutation = useDeleteMappingMutation();
  const deletePointMutation = useDeletePointMutation();

  const selectedDevice = getSelectedDevice(devices, selectedDeviceId);
  const mappingByPointId = useMemo(
    () => new Map(mappings.map((mapping) => [mapping.point_id, mapping])),
    [mappings],
  );
  const tagById = useMemo(
    () => new Map(tags.map((tag) => [tag.id, tag])),
    [tags],
  );
  const [prefix, setPrefix] = useState('TAG');
  const [strategy, setStrategy] = useState<TagBindingStrategy>('address');
  const [flowMode, setFlowMode] = useState<TagBindingFlowMode>('create');
  const [statusFilter, setStatusFilter] = useState<TagBindingStatusFilter>('all');
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedPointIds, setSelectedPointIds] = useState<string[]>([]);
  const [existingTagSelections, setExistingTagSelections] = useState<Record<string, string>>(
    {},
  );
  const [batchSummary, setBatchSummary] = useState<TagBindingBatchSummary | null>(null);
  const [templateExpanded, setTemplateExpanded] = useState(false);
  const [tagLibraryQuery, setTagLibraryQuery] = useState('');
  const [tagDraftKey, setTagDraftKey] = useState('');
  const [tagDraftDisplayName, setTagDraftDisplayName] = useState('');
  const [tagDraftDataType, setTagDraftDataType] = useState<DataType>('int16');
  const [tagLibraryFeedback, setTagLibraryFeedback] = useState<TagLibraryFeedback | null>(null);
  const [isBatchUnbinding, setIsBatchUnbinding] = useState(false);
  const [isBatchDeletingPoints, setIsBatchDeletingPoints] = useState(false);

  const pointIdsKey = useMemo(() => points.map((point) => point.id).join('|'), [points]);

  const candidates = useMemo(
    () =>
      buildTagBindingCandidates({
        points,
        tags,
        mappings,
        template: {
          prefix,
          strategy,
        },
      }),
    [mappings, points, prefix, strategy, tags],
  );

  useEffect(() => {
    setSelectedPointIds([]);
    setBatchSummary(null);
  }, [pointIdsKey, selectedDeviceId]);

  useEffect(() => {
    setExistingTagSelections((currentState) => {
      let didChange = false;
      const nextState = { ...currentState };
      for (const candidate of candidates) {
        if (!nextState[candidate.pointId] && candidate.existingTagOptions[0]) {
          nextState[candidate.pointId] = candidate.existingTagOptions[0].id;
          didChange = true;
        }
      }
      return didChange ? nextState : currentState;
    });
  }, [candidates]);

  const filteredCandidates = useMemo(() => {
    const normalizedQuery = searchQuery.trim().toLowerCase();
    return candidates.filter((candidate) => {
      const matchesSearch = !normalizedQuery
        || candidate.pointName.toLowerCase().includes(normalizedQuery)
        || candidate.pointAddress.toLowerCase().includes(normalizedQuery)
        || candidate.previewKey.toLowerCase().includes(normalizedQuery);
      const matchesStatus =
        statusFilter === 'all' || candidate.bindingStatus === statusFilter;
      return matchesSearch && matchesStatus;
    });
  }, [candidates, searchQuery, statusFilter]);
  const mappingCountByTagId = useMemo(() => {
    const counts = new Map<string, number>();
    for (const mapping of mappings) {
      counts.set(mapping.tag_id, (counts.get(mapping.tag_id) ?? 0) + 1);
    }
    return counts;
  }, [mappings]);
  const filteredTagLibrary = useMemo(() => {
    const normalizedQuery = tagLibraryQuery.trim().toLowerCase();
    const sortedTags = [...tags].sort((left, right) => left.key.localeCompare(right.key));

    if (normalizedQuery === '') {
      return sortedTags;
    }

    return sortedTags.filter((tag) =>
      `${tag.key} ${tag.display_name} ${tag.data_type}`.toLowerCase().includes(normalizedQuery),
    );
  }, [tagLibraryQuery, tags]);
  const linkedTagCount = useMemo(
    () => tags.filter((tag) => (mappingCountByTagId.get(tag.id) ?? 0) > 0).length,
    [mappingCountByTagId, tags],
  );
  const unusedTagCount = tags.length - linkedTagCount;

  const selectedCandidates = useMemo(
    () =>
      candidates.filter((candidate) => selectedPointIds.includes(candidate.pointId)),
    [candidates, selectedPointIds],
  );

  const createRequests = useMemo(
    () => buildTagBindingRequests(candidates, selectedPointIds),
    [candidates, selectedPointIds],
  );

  const existingRequests = useMemo(
    () =>
      candidates
        .filter(
          (candidate) =>
            selectedPointIds.includes(candidate.pointId)
            && !candidate.alreadyLinked
            && Boolean(existingTagSelections[candidate.pointId]),
        )
        .map((candidate) => ({
          pointId: candidate.pointId,
          tagId: existingTagSelections[candidate.pointId]!,
          tagKey:
            candidate.existingTagOptions.find(
              (option) => option.id === existingTagSelections[candidate.pointId],
            )?.key ?? existingTagSelections[candidate.pointId]!,
        })),
    [candidates, existingTagSelections, selectedPointIds],
  );

  const batchDiffPreview = useMemo(
    () =>
      buildBatchDiffPreview({
        candidates,
        selectedPointIds,
        flowMode,
        existingTagSelections,
      }),
    [candidates, selectedPointIds, flowMode, existingTagSelections],
  );

  const blockedSelectionCount = selectedCandidates.filter((candidate) =>
    flowMode === 'create'
      ? candidate.conflict || candidate.alreadyLinked
      : candidate.alreadyLinked || !existingTagSelections[candidate.pointId],
  ).length;

  const readyCount = flowMode === 'create'
    ? createRequests.length
    : existingRequests.length;
  const generatedReviewCount = useMemo(
    () => candidates.filter((candidate) => candidate.alreadyLinked).length,
    [candidates],
  );
  const needsReviewCount = useMemo(
    () => candidates.filter((candidate) => candidate.bindingStatus !== 'bound').length,
    [candidates],
  );
  const selectedExceptionCount = useMemo(
    () => selectedCandidates.filter((candidate) => !candidate.alreadyLinked).length,
    [selectedCandidates],
  );
  const selectedBoundItems = useMemo(
    () =>
      selectedCandidates
        .map((candidate) => ({
          candidate,
          mapping: mappingByPointId.get(candidate.pointId) ?? null,
          boundTag: (() => {
            const mapping = mappingByPointId.get(candidate.pointId);
            return mapping ? tagById.get(mapping.tag_id) ?? null : null;
          })(),
        }))
        .filter(
          (item): item is {
            candidate: TagBindingCandidate;
            mapping: NonNullable<ReturnType<typeof mappingByPointId.get>>;
            boundTag: Tag | null;
          } => Boolean(item.mapping),
        ),
    [mappingByPointId, selectedCandidates, tagById],
  );
  const isAnyUnbindPending = deleteMappingMutation.isPending || isBatchUnbinding;
  const isPointDeleteBusy = deletePointMutation.isPending || isBatchDeletingPoints;
  const isTagBoardBusy = isAnyUnbindPending || isPointDeleteBusy;
  const selectionHintKey = batchSummary
    && batchSummary.failureCount === 0
    && (
      batchSummary.createdCount > 0
      || batchSummary.linkedCount > 0
      || batchSummary.unboundCount > 0
      || batchSummary.deletedPointCount > 0
    )
    ? 'workbench.tag.board.selectionHint.done'
    : selectedCandidates.length === 0
      ? 'workbench.tag.board.selectionHint.none'
      : blockedSelectionCount > 0
        ? 'workbench.tag.board.selectionHint.conflicts'
        : 'workbench.tag.board.selectionHint.ready';
  const flowModeHintKey = flowMode === 'create'
    ? 'workbench.tag.board.flowModeHint.create'
    : 'workbench.tag.board.flowModeHint.existing';

  const handleTogglePoint = (pointId: string) => {
    setSelectedPointIds((previous) =>
      previous.includes(pointId)
        ? previous.filter((id) => id !== pointId)
        : [...previous, pointId],
    );
    setBatchSummary(null);
  };

  const handleSelectAll = () => {
    setSelectedPointIds(filteredCandidates.map((candidate) => candidate.pointId));
    setBatchSummary(null);
  };

  const handleSelectBindable = () => {
    setSelectedPointIds(
      filteredCandidates
        .filter((candidate) =>
          flowMode === 'create'
            ? !candidate.conflict && !candidate.alreadyLinked
            : !candidate.alreadyLinked && candidate.existingTagOptions.length > 0,
        )
        .map((candidate) => candidate.pointId),
    );
    setBatchSummary(null);
  };

  const handleClearSelection = () => {
    setSelectedPointIds([]);
    setBatchSummary(null);
  };

  const handleExistingTagSelection = (pointId: string, tagId: string) => {
    setExistingTagSelections((currentState) => ({
      ...currentState,
      [pointId]: tagId,
    }));
    setBatchSummary(null);
  };

  const handleCreateStandaloneTag = async () => {
    if (!tagDraftKey.trim()) {
      return;
    }

    try {
      const createdTag = await createTagMutation.mutateAsync({
        key: tagDraftKey.trim(),
        display_name: tagDraftDisplayName.trim() || undefined,
        data_type: tagDraftDataType,
      });
      setTagDraftKey('');
      setTagDraftDisplayName('');
      setFocusedTagIds([createdTag.id]);
      setTagLibraryFeedback({
        tone: 'success',
        message: t('workbench.tag.master.createSuccess', { key: createdTag.key }),
      });
    } catch (error) {
      setTagLibraryFeedback({
        tone: 'error',
        message: getErrorMessage(error, t('workbench.tag.master.createFailed')),
      });
    }
  };

  const handleDeleteStandaloneTag = async (tag: Tag) => {
    if (!window.confirm(t('workbench.tag.master.deleteConfirm', { key: tag.key }))) {
      return;
    }

    try {
      await deleteTagMutation.mutateAsync(tag.id);
      setTagLibraryFeedback({
        tone: 'success',
        message: t('workbench.tag.master.deleteSuccess', { key: tag.key }),
      });
    } catch (error) {
      setTagLibraryFeedback({
        tone: 'error',
        message: getErrorMessage(error, t('workbench.tag.master.deleteFailed')),
      });
    }
  };

  const handleBatchBind = async () => {
    if (!selectedDevice || blockedSelectionCount > 0) {
      return;
    }

    if (flowMode === 'create') {
      if (createRequests.length === 0) {
        return;
      }

      const boundTagIds: string[] = [];
      const failures: TagBindingFailure[] = [];

      try {
        // 使用 batch API 一次建立所有 tags，避免 SQLite lock
        const batchResult = await tagAPI.batchCreate(
          createRequests.map((request) => request.tagRequest),
        );

        // 取得所有新建立的 tags（需要 tag ID 建立 mapping）
        const createdTagIds = batchResult.created;
        const tagErrors = batchResult.errors;

        // 記錄 tag 建立失敗
        for (const tagError of tagErrors) {
          const request = createRequests.find((r) => r.tagKey === tagError.key);
          if (request) {
            failures.push({
              pointId: request.pointId,
              tagKey: request.tagKey,
              stage: 'tag',
              error: tagError.error,
            });
          }
        }

        // 依序建立 mappings（逐一以避免 SQLite lock）
        if (createdTagIds.length > 0) {
          // 重新查詢以取得剛建立的 tags
          const { data: freshTags = [] } = await tagsQuery.refetch();

          for (const request of createRequests) {
            // 跳過 tag 建立失敗的
            if (tagErrors.some((e) => e.key === request.tagKey)) {
              continue;
            }

            const matchedTag = freshTags.find(
              (t) => t.key.toLowerCase() === request.tagKey.toLowerCase(),
            );
            if (!matchedTag) {
              failures.push({
                pointId: request.pointId,
                tagKey: request.tagKey,
                stage: 'tag',
                error: t('workbench.tag.results.tagFallback'),
              });
              continue;
            }

            try {
              await createMappingMutation.mutateAsync({
                point_id: request.pointId,
                tag_id: matchedTag.id,
                enabled: true,
              });
              boundTagIds.push(matchedTag.id);
            } catch (error) {
              failures.push({
                pointId: request.pointId,
                tagKey: request.tagKey,
                stage: 'mapping',
                error: getErrorMessage(error, t('workbench.tag.results.mappingFallback')),
              });
            }
          }
        }
      } catch (error) {
        // batch API 整體失敗
        for (const request of createRequests) {
          failures.push({
            pointId: request.pointId,
            tagKey: request.tagKey,
            stage: 'tag',
            error: getErrorMessage(error, t('workbench.tag.results.tagFallback')),
          });
        }
      }

      if (boundTagIds.length > 0) {
        setFocusedTagIds(boundTagIds);
      }

      setBatchSummary({
        mode: 'bind',
        createdCount: createRequests.length - failures.length,
        linkedCount: 0,
        unboundCount: 0,
        deletedPointCount: 0,
        skippedCount: batchDiffPreview.skipped.length,
        failureCount: failures.length,
        failures,
      });
      return;
    }

    if (existingRequests.length === 0) {
      return;
    }

    const linkedTagIds: string[] = [];
    const existingFailures: TagBindingFailure[] = [];

    // 依序建立 mappings 以避免 SQLite lock
    for (const request of existingRequests) {
      try {
        await createMappingMutation.mutateAsync({
          point_id: request.pointId,
          tag_id: request.tagId,
          enabled: true,
        });
        linkedTagIds.push(request.tagId);
      } catch (error) {
        existingFailures.push({
          pointId: request.pointId,
          tagKey: request.tagKey,
          stage: 'mapping',
          error: getErrorMessage(error, t('workbench.tag.results.mappingFallback')),
        });
      }
    }

    if (linkedTagIds.length > 0) {
      setFocusedTagIds(linkedTagIds);
    }

    setBatchSummary({
      mode: 'bind',
      createdCount: 0,
      linkedCount: existingRequests.length - existingFailures.length,
      unboundCount: 0,
      deletedPointCount: 0,
      skippedCount: batchDiffPreview.skipped.length,
      failureCount: existingFailures.length,
      failures: existingFailures,
    });
  };

  const handleUnbind = async (pointId: string) => {
    const mapping = mappingByPointId.get(pointId);
    if (!mapping) {
      return;
    }

    if (!window.confirm(t('workbench.tag.actions.unbindConfirm'))) {
      return;
    }

    const point = points.find((p) => p.id === pointId);

    try {
      await deleteMappingMutation.mutateAsync(mapping.id);
      setBatchSummary({
        mode: 'unbind',
        createdCount: 0,
        linkedCount: 0,
        unboundCount: 1,
        deletedPointCount: 0,
        skippedCount: 0,
        failureCount: 0,
        failures: [],
      });
      setSelectedPointIds((previous) => previous.filter((id) => id !== pointId));
    } catch {
      setBatchSummary({
        mode: 'unbind',
        createdCount: 0,
        linkedCount: 0,
        unboundCount: 0,
        deletedPointCount: 0,
        skippedCount: 0,
        failureCount: 1,
        failures: [
          {
            pointId,
            tagKey: point?.name ?? pointId,
            stage: 'mapping',
            error: t('workbench.tag.results.unbindFailed'),
          },
        ],
      });
    }
  };

  const handleBatchUnbind = async () => {
    if (selectedBoundItems.length === 0 || isBatchUnbinding || isPointDeleteBusy) {
      return;
    }

    if (
      !window.confirm(
        t('workbench.tag.actions.unbindSelectedConfirm', {
          count: selectedBoundItems.length,
        }),
      )
    ) {
      return;
    }

    const failures: TagBindingFailure[] = [];
    const succeededPointIds: string[] = [];
    setIsBatchUnbinding(true);

    try {
      for (const item of selectedBoundItems) {
        try {
          await deleteMappingMutation.mutateAsync(item.mapping.id);
          succeededPointIds.push(item.candidate.pointId);
        } catch {
          failures.push({
            pointId: item.candidate.pointId,
            tagKey: item.boundTag?.key ?? item.candidate.previewKey,
            stage: 'mapping',
            error: t('workbench.tag.results.unbindFailed'),
          });
        }
      }

      setBatchSummary({
        mode: 'unbind',
        createdCount: 0,
        linkedCount: 0,
        unboundCount: succeededPointIds.length,
        deletedPointCount: 0,
        skippedCount: 0,
        failureCount: failures.length,
        failures,
      });
      setSelectedPointIds((previous) => previous.filter((id) => !succeededPointIds.includes(id)));
    } finally {
      setIsBatchUnbinding(false);
    }
  };

  /**
   * 批次刪除目前勾選的 point；資料庫會 CASCADE 移除對應 mapping，並同步刷新快取。
   */
  const handleDeleteSelectedPoints = async () => {
    if (selectedPointIds.length === 0 || isBatchDeletingPoints) {
      return;
    }

    if (
      !window.confirm(
        t('workbench.tag.actions.deleteSelectedConfirm', { count: selectedPointIds.length }),
      )
    ) {
      return;
    }

    const failures: TagBindingFailure[] = [];
    const succeededPointIds: string[] = [];
    setIsBatchDeletingPoints(true);

    try {
      for (const pointId of selectedPointIds) {
        const point = points.find((p) => p.id === pointId);
        try {
          await deletePointMutation.mutateAsync(pointId);
          succeededPointIds.push(pointId);
        } catch (error) {
          failures.push({
            pointId,
            tagKey: point?.name ?? pointId,
            stage: 'point',
            error: getErrorMessage(error, t('workbench.tag.results.deletePointFailed')),
          });
        }
      }

      await queryClient.invalidateQueries({ queryKey: mappingKeys.lists() });

      setBatchSummary({
        mode: 'deletePoints',
        createdCount: 0,
        linkedCount: 0,
        unboundCount: 0,
        deletedPointCount: succeededPointIds.length,
        skippedCount: 0,
        failureCount: failures.length,
        failures,
      });
      setSelectedPointIds((previous) => previous.filter((id) => !succeededPointIds.includes(id)));
    } finally {
      setIsBatchDeletingPoints(false);
    }
  };

  if (!selectedDevice) {
    return (
      <section className="min-h-0 flex-1 space-y-6 overflow-y-auto overscroll-contain rounded-2xl border border-dashed border-slate-700 bg-slate-950/40 p-6">
        <div className="space-y-2">
          <p className="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-300">
            {t('workbench.tag.empty.eyebrow')}
          </p>
          <h2 className="text-2xl font-semibold text-slate-50">
            {t('workbench.tag.empty.deviceTitle')}
          </h2>
          <p className="max-w-2xl text-sm text-slate-300">
            {t('workbench.tag.empty.deviceDescription')}
          </p>
        </div>

        <div className="grid gap-3 sm:grid-cols-2">
          {devices.map((device) => (
            <button
              key={device.id}
              type="button"
              onClick={() => setSelectedDeviceId(device.id)}
              aria-label={device.name}
              className="rounded-2xl border border-slate-800 bg-slate-900/80 p-4 text-left transition hover:border-cyan-500/40 hover:bg-slate-900"
            >
              <p className="text-sm font-semibold text-slate-50">{device.name}</p>
              <p className="mt-1 text-xs text-slate-400">{device.protocol}</p>
            </button>
          ))}
        </div>
      </section>
    );
  }

  if (points.length === 0) {
    const eligibleSpanCount = selectedDevice
      ? countEligibleSpans({
          rules: sourcePlanningState.rules,
          points: [],
          protocol: selectedDevice.protocol,
        })
      : 0;

    return (
      <section className="min-h-0 flex-1 space-y-6 overflow-y-auto overscroll-contain rounded-2xl border border-dashed border-slate-700 bg-slate-950/40 p-6">
        <div className="space-y-2">
          <p className="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-300">
            {t('workbench.tag.empty.eyebrow')}
          </p>
          <h2 className="text-2xl font-semibold text-slate-50">
            {t('workbench.tag.empty.title')}
          </h2>
          <p className="max-w-2xl text-sm text-slate-300">
            {t('workbench.tag.empty.description')}
          </p>
        </div>

        <div className="flex items-center gap-3 rounded-xl border border-slate-800 bg-slate-900/70 px-4 py-3">
          <span
            className="text-sm text-slate-300"
            data-testid="tag-empty-eligible-spans"
            data-count={eligibleSpanCount}
          >
            {eligibleSpanCount > 0
              ? t('workbench.tag.empty.eligibleSpans', { count: eligibleSpanCount })
              : t('workbench.tag.empty.eligibleSpansNone')}
          </span>
        </div>

        <button
          type="button"
          onClick={() => setActiveStep('source')}
          data-testid="tag-empty-goto-source"
          className="rounded-xl border border-cyan-500/30 bg-cyan-500/10 px-4 py-3 text-sm font-semibold text-cyan-200 transition hover:bg-cyan-500/20"
        >
          {t('workbench.tag.empty.goToSource')}
        </button>
      </section>
    );
  }

  return (
    <section className="grid h-full min-h-0 gap-6 xl:grid-rows-1 xl:grid-cols-[minmax(0,1.45fr)_minmax(340px,0.95fr)]">
      <div
        className="relative flex h-full min-h-0 flex-col overflow-hidden rounded-2xl border border-slate-800/70 bg-gradient-to-b from-slate-900/95 to-slate-950 shadow-[0_1px_0_0_rgba(255,255,255,0.05)_inset]"
        data-testid="tag-board-surface"
      >
        <div
          aria-hidden="true"
          className="pointer-events-none absolute inset-x-0 top-0 h-px bg-gradient-to-r from-transparent via-cyan-500/35 to-transparent"
        />
        <div
          aria-hidden="true"
          className="pointer-events-none absolute -left-20 -top-28 h-56 w-56 rounded-full bg-cyan-500/[0.07] blur-3xl"
        />
        <div className="relative flex min-h-0 flex-1 flex-col divide-y divide-slate-800/70">
          <header className="shrink-0 space-y-1 px-4 pb-3 pt-4">
            <p className="text-[10px] font-semibold uppercase tracking-[0.2em] text-cyan-400/85">
              {t('workbench.tag.selection.eyebrow')}
            </p>
            <h2 className="text-lg font-semibold tracking-tight text-slate-50">
              {t('workbench.tag.selection.title')}
            </h2>
            <p className="line-clamp-2 text-xs leading-relaxed text-slate-500">
              {t('workbench.tag.selection.description')}
            </p>
          </header>

          <div className="flex shrink-0 flex-col" data-testid="tag-review-summary">
            <div className="space-y-1 px-4 py-2.5">
              <p className="text-[10px] font-semibold uppercase tracking-[0.14em] text-cyan-400/80">
                {t('workbench.tag.review.eyebrow')}
              </p>
              <h3 className="text-base font-semibold text-slate-50">
                {t('workbench.tag.review.title')}
              </h3>
              <p className="line-clamp-2 text-xs text-slate-500">
                {t('workbench.tag.review.description')}
              </p>
            </div>
            <dl className="flex min-h-[4.25rem] items-stretch border-t border-slate-800/60 bg-slate-950/40">
              <div className="flex min-w-0 flex-1 flex-col justify-center bg-emerald-500/[0.06] px-3 py-2 text-center sm:px-4">
                <dt className="text-[9px] font-medium uppercase tracking-[0.12em] text-emerald-400/90">
                  {t('workbench.tag.review.metrics.generated')}
                </dt>
                <dd className="mt-0.5">
                  <span
                    className="font-mono text-xl font-semibold tabular-nums text-emerald-100"
                    data-testid="tag-review-generated"
                    title={t('workbench.tag.review.metrics.generatedFootnote')}
                  >
                    {generatedReviewCount}
                  </span>
                </dd>
              </div>
              <div
                aria-hidden="true"
                className="w-px shrink-0 bg-gradient-to-b from-transparent via-slate-600/50 to-transparent"
              />
              <div className="flex min-w-0 flex-1 flex-col justify-center px-3 py-2 text-center sm:px-4">
                <dt className="text-[9px] font-medium uppercase tracking-[0.12em] text-amber-400/85">
                  {t('workbench.tag.review.metrics.needsReview')}
                </dt>
                <dd className="mt-0.5">
                  <span
                    className="font-mono text-xl font-semibold tabular-nums text-amber-100"
                    data-testid="tag-review-needs-review"
                  >
                    {needsReviewCount}
                  </span>
                </dd>
              </div>
              <div
                aria-hidden="true"
                className="w-px shrink-0 bg-gradient-to-b from-transparent via-slate-600/50 to-transparent"
              />
              <div className="flex min-w-0 flex-1 flex-col justify-center px-3 py-2 text-center sm:px-4">
                <dt className="text-[9px] font-medium uppercase tracking-[0.12em] text-slate-500">
                  {t('workbench.tag.review.metrics.selectedExceptions')}
                </dt>
                <dd className="mt-0.5">
                  <span
                    className="font-mono text-xl font-semibold tabular-nums text-slate-100"
                    data-testid="tag-review-selected-exceptions"
                  >
                    {selectedExceptionCount}
                  </span>
                </dd>
              </div>
            </dl>
          </div>

          <div className="flex shrink-0 flex-col gap-2 bg-slate-950/20 px-4 py-3 sm:flex-row sm:items-center">
            <input
              aria-label={t('workbench.tag.board.search')}
              value={searchQuery}
              onChange={(event) => setSearchQuery(event.target.value)}
              placeholder={t('workbench.tag.board.search')}
              className="min-w-0 flex-1 rounded-xl border-0 bg-slate-900/50 px-3.5 py-2.5 text-sm text-slate-100 shadow-inner shadow-black/20 placeholder:text-slate-500 ring-1 ring-slate-700/40 transition-[box-shadow,ring-color] focus:outline-none focus:ring-2 focus:ring-cyan-400/35"
            />
            <label className="sr-only" htmlFor="tag-board-status-filter">
              {t('workbench.tag.board.statusFilter')}
            </label>
            <select
              id="tag-board-status-filter"
              aria-label={t('workbench.tag.board.statusFilter')}
              value={statusFilter}
              onChange={(event) =>
                setStatusFilter(event.target.value as TagBindingStatusFilter)
              }
              className="h-[42px] w-full shrink-0 rounded-xl border-0 bg-slate-900/50 px-3.5 text-sm text-slate-100 shadow-inner shadow-black/15 ring-1 ring-slate-700/40 sm:w-[200px] focus:outline-none focus:ring-2 focus:ring-cyan-400/35"
            >
              <option value="all">{t('workbench.tag.board.statusFilters.all')}</option>
              <option value="unbound">
                {t('workbench.tag.board.statusFilters.unbound')}
              </option>
              <option value="partial">
                {t('workbench.tag.board.statusFilters.partial')}
              </option>
              <option value="bound">{t('workbench.tag.board.statusFilters.bound')}</option>
            </select>
          </div>

          <div className="shrink-0 space-y-2 px-4 py-3">
            <div className="flex flex-wrap items-center gap-2">
              <button
                type="button"
                onClick={handleSelectAll}
                className="rounded-xl bg-slate-800/70 px-3.5 py-2 text-sm text-slate-200 shadow-sm shadow-black/20 ring-1 ring-slate-700/45 transition hover:bg-slate-800 hover:ring-slate-600/55"
              >
                {t('workbench.tag.actions.selectAll')}
              </button>
              <button
                type="button"
                onClick={handleSelectBindable}
                className="rounded-xl bg-slate-800/70 px-3.5 py-2 text-sm text-slate-200 shadow-sm shadow-black/20 ring-1 ring-slate-700/45 transition hover:bg-slate-800 hover:ring-slate-600/55"
              >
                {t('workbench.tag.actions.selectBindable')}
              </button>
              <button
                type="button"
                onClick={handleClearSelection}
                className="rounded-xl bg-slate-800/70 px-3.5 py-2 text-sm text-slate-200 shadow-sm shadow-black/20 ring-1 ring-slate-700/45 transition hover:bg-slate-800 hover:ring-slate-600/55"
              >
                {t('workbench.tag.actions.clearSelection')}
              </button>
              <button
                type="button"
                data-testid="tag-delete-selected-points"
                onClick={() => void handleDeleteSelectedPoints()}
                disabled={selectedPointIds.length === 0 || isTagBoardBusy}
                className="rounded-xl bg-rose-500/[0.12] px-3.5 py-2 text-sm font-medium text-rose-100 ring-1 ring-rose-400/25 transition hover:bg-rose-500/20 hover:ring-rose-400/35 disabled:cursor-not-allowed disabled:opacity-50"
              >
                {t('workbench.tag.actions.deleteSelected')}
              </button>
            </div>
            <p className="text-xs leading-relaxed text-slate-500">{t(selectionHintKey)}</p>
          </div>

          <div
            className="min-h-0 flex-1 space-y-2 overflow-y-auto overscroll-contain p-2 sm:p-3 sm:pt-2"
            data-testid="tag-candidate-board"
            data-layout="two-line-board"
          >
          {filteredCandidates.map((candidate) => {
            const selected = selectedPointIds.includes(candidate.pointId);
            const boundMapping = mappingByPointId.get(candidate.pointId);
            const boundTag = boundMapping ? tagById.get(boundMapping.tag_id) : null;
            const selectedExistingTag =
              candidate.existingTagOptions.find(
                (option) => option.id === existingTagSelections[candidate.pointId],
              ) ?? null;

            return (
              <article
                key={candidate.pointId}
                data-layout="two-line"
                data-testid={`tag-candidate-${candidate.pointId}`}
                onClick={() => {
                  if (boundMapping) {
                    setFocusedTagIds([boundMapping.tag_id]);
                  }
                  setInspectorSelection({
                    kind: 'tag',
                    tagId:
                      boundMapping?.tag_id
                      ?? selectedExistingTag?.key
                      ?? candidate.previewKey,
                    pointId: candidate.pointId,
                    pointName: candidate.pointName,
                    pointAddress: candidate.pointAddress,
                    rawValue: candidate.rawValue,
                    transformedValue: candidate.transformedValue,
                    bitWidth: candidate.bitWidth,
                    cellSpan: candidate.cellSpan,
                    bindingStatus: candidate.bindingStatus,
                    conflictReason: candidate.conflictReason,
                    alreadyLinked: candidate.alreadyLinked,
                    existingTagLabel: selectedExistingTag?.displayName ?? null,
                    });
                }}
                className={[
                  'group/tag-card relative grid cursor-pointer grid-cols-[auto_minmax(0,1fr)] gap-x-3 gap-y-0 rounded-xl px-3 py-2.5 ring-1 ring-slate-800/70 transition-[background-color,box-shadow,ring-color] duration-200',
                  selected
                    ? 'bg-cyan-500/[0.07] ring-cyan-400/30 shadow-[0_4px_14px_-6px_rgba(0,0,0,0.55),inset_3px_0_0_0_rgba(34,211,238,0.65)]'
                    : candidate.alreadyLinked
                      ? 'bg-emerald-950/[0.22] ring-emerald-500/15 shadow-[0_4px_14px_-6px_rgba(0,0,0,0.5),inset_3px_0_0_0_rgba(52,211,153,0.55)] hover:bg-emerald-950/[0.32] hover:ring-emerald-400/25'
                      : 'bg-slate-900/40 shadow-[0_4px_14px_-6px_rgba(0,0,0,0.45),inset_3px_0_0_0_rgba(71,85,105,0.55)] hover:bg-slate-900/65 hover:ring-slate-600/50',
                ].join(' ')}
              >
                <div className="row-span-2 flex min-h-9 min-w-9 shrink-0 items-start justify-center pt-0.5">
                  <input
                    type="checkbox"
                    checked={selected}
                    onClick={(event) => event.stopPropagation()}
                    onChange={() => handleTogglePoint(candidate.pointId)}
                    aria-label={candidate.pointName}
                    className="h-4 w-4 rounded border-slate-600 bg-slate-950/80 text-cyan-400 shadow-inner shadow-black/30 ring-offset-2 ring-offset-slate-900 transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400/50"
                  />
                </div>

                <div className="flex min-w-0 gap-2.5">
                  <div className="min-w-0 flex-1 space-y-1">
                    <p className="truncate text-[0.9375rem] font-semibold leading-tight tracking-tight text-slate-50">
                      {candidate.pointName}
                    </p>
                    <p className="line-clamp-2 text-[11px] leading-snug text-slate-400 sm:line-clamp-1">
                      {t('workbench.tag.selection.pointMeta', {
                        address: candidate.pointAddress,
                        dataType: candidate.dataType,
                      })}
                      <span className="mx-1 text-slate-600" aria-hidden>
                        ·
                      </span>
                      {t('workbench.tag.board.span', {
                        cells: candidate.cellSpan,
                        bitWidth: candidate.bitWidth,
                      })}
                    </p>
                  </div>
                  <span
                    className={`shrink-0 self-start rounded-md border px-2 py-1 text-[10px] font-semibold tracking-wide ${getStatusToneClass(
                      candidate.bindingStatus,
                    )}`}
                    data-testid={`tag-status-${candidate.pointId}`}
                  >
                    {t(`workbench.tag.board.status.${candidate.bindingStatus}`)}
                  </span>
                </div>

                <div className="col-start-2 mt-2 min-w-0">
                  <div className="rounded-lg bg-slate-950/35 p-2 ring-1 ring-slate-800/60 backdrop-blur-[2px]">
                    <div className="flex flex-col gap-2.5 sm:flex-row sm:items-stretch sm:gap-3">
                      <div className="flex min-w-0 flex-1 flex-col gap-2.5 text-[11px] leading-snug sm:flex-row sm:gap-0 sm:divide-x sm:divide-slate-700/50">
                        <span className="min-w-0 sm:pr-3">
                          <span className="block text-[9px] font-semibold uppercase tracking-[0.12em] text-slate-500">
                            {t('workbench.tag.board.rawValue')}
                          </span>
                          <span
                            className="mt-0.5 block break-all font-mono text-[13px] tabular-nums text-slate-100"
                            data-testid={`tag-raw-${candidate.pointId}`}
                          >
                            {formatCandidateValue(candidate.rawValue)}
                          </span>
                        </span>
                        <span className="min-w-0 sm:pl-3">
                          <span className="block text-[9px] font-semibold uppercase tracking-[0.12em] text-slate-500">
                            {t('workbench.tag.board.transformedValue')}
                          </span>
                          <span className="mt-0.5 block break-all font-mono text-[13px] tabular-nums text-slate-100">
                            {formatCandidateValue(candidate.transformedValue)}
                          </span>
                        </span>
                      </div>
                      <div className="min-w-0 shrink-0 sm:w-[min(100%,13.5rem)]">
                        {candidate.alreadyLinked && boundTag ? (
                          <div className="flex items-center justify-end gap-2 sm:justify-start">
                            <span
                              data-testid={`tag-bound-label-${candidate.pointId}`}
                              className="min-w-0 truncate text-left text-[11px] font-medium text-emerald-200/95"
                            >
                              {t('workbench.tag.board.boundTag', { tag: boundTag.key })}
                            </span>
                            <button
                              type="button"
                              data-testid={`tag-unbind-${candidate.pointId}`}
                              onClick={(event) => {
                                event.stopPropagation();
                                void handleUnbind(candidate.pointId);
                              }}
                              disabled={isTagBoardBusy}
                              className="min-h-8 shrink-0 rounded-lg bg-rose-500/10 px-2.5 py-1 text-[11px] font-semibold text-rose-100 ring-1 ring-rose-400/20 transition hover:bg-rose-500/20 hover:ring-rose-400/35 disabled:opacity-50"
                            >
                              {t('workbench.tag.actions.unbind')}
                            </button>
                          </div>
                        ) : flowMode === 'create' ? (
                          <div
                            data-testid={`tag-preview-${candidate.pointId}`}
                            data-conflict={candidate.conflict ? 'true' : 'false'}
                            className={`rounded-lg px-2.5 py-1.5 text-xs font-semibold leading-snug ring-1 ${
                              candidate.conflict
                                ? 'bg-rose-500/15 text-rose-50 ring-rose-400/25'
                                : 'bg-gradient-to-br from-slate-900/90 to-slate-950 font-mono text-cyan-100 ring-cyan-500/15 shadow-inner shadow-black/20'
                            }`}
                          >
                            <span className="break-all">{candidate.previewKey}</span>
                          </div>
                        ) : null}

                        {flowMode === 'existing' && !candidate.alreadyLinked ? (
                          <select
                            aria-label={t('workbench.tag.board.existingTag')}
                            data-testid={`existing-tag-select-${candidate.pointId}`}
                            value={existingTagSelections[candidate.pointId] ?? ''}
                            onClick={(event) => event.stopPropagation()}
                            onChange={(event) =>
                              handleExistingTagSelection(
                                candidate.pointId,
                                event.target.value,
                              )
                            }
                            className="min-h-9 w-full rounded-lg border-0 bg-slate-900/70 px-2.5 py-1.5 text-xs text-slate-100 shadow-inner shadow-black/15 ring-1 ring-slate-700/45 focus:outline-none focus:ring-2 focus:ring-cyan-400/35"
                          >
                            <option value="">
                              {t('workbench.tag.board.noExistingTag')}
                            </option>
                            {candidate.existingTagOptions.map((option) => (
                              <option key={option.id} value={option.id}>
                                {option.displayName}
                              </option>
                            ))}
                          </select>
                        ) : null}
                      </div>
                    </div>
                  </div>
                </div>
              </article>
            );
          })}
          </div>
        </div>
      </div>

      <aside className="flex h-full min-h-0 flex-col gap-6 rounded-2xl border border-slate-800 bg-slate-950/40 p-5">
        <div
          className="shrink-0 space-y-4 rounded-2xl border border-slate-800 bg-slate-900/70 p-4"
          data-testid="tag-exception-tools"
        >
          <div className="space-y-1">
            <p className="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-300">
              {t('workbench.tag.exception.eyebrow')}
            </p>
            <h3 className="text-xl font-semibold text-slate-50">
              {t('workbench.tag.exception.title')}
            </h3>
            <p className="text-sm text-slate-300">
              {t('workbench.tag.exception.description')}
            </p>
          </div>

          <div className="flex flex-wrap gap-2">
            <button
              type="button"
              onClick={() => setFlowMode('create')}
              className={
                flowMode === 'create'
                  ? 'rounded-lg bg-cyan-500 px-3 py-2 text-sm font-semibold text-slate-950'
                  : 'rounded-lg border border-slate-800 px-3 py-2 text-sm text-slate-300'
              }
            >
              {t('workbench.tag.board.flow.create')}
            </button>
            <button
              type="button"
              onClick={() => setFlowMode('existing')}
              className={
                flowMode === 'existing'
                  ? 'rounded-lg bg-cyan-500 px-3 py-2 text-sm font-semibold text-slate-950'
                  : 'rounded-lg border border-slate-800 px-3 py-2 text-sm text-slate-300'
              }
            >
              {t('workbench.tag.board.flow.existing')}
            </button>
          </div>

          <p className="text-sm text-slate-300">{t(flowModeHintKey)}</p>
        </div>

        <div className="shrink-0 space-y-4">
          <button
            type="button"
            onClick={() => setTemplateExpanded((v) => !v)}
            className="flex w-full items-center justify-between text-left"
            data-testid="tag-template-toggle"
          >
            <div className="space-y-1">
              <p className="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-300">
                {t('workbench.tag.template.eyebrow')}
              </p>
              <h3 className="text-xl font-semibold text-slate-50">
                {t('workbench.tag.template.toggle')}
              </h3>
            </div>
            <span className="text-sm text-slate-400">{templateExpanded ? '▲' : '▼'}</span>
          </button>

          {templateExpanded ? (
            <div className="grid gap-3">
              <p className="text-sm text-slate-300">
                {t('workbench.tag.template.description')}
              </p>
              <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                <span>{t('workbench.tag.template.prefix')}</span>
                <input
                  aria-label={t('workbench.tag.template.prefix')}
                  value={prefix}
                  onChange={(event) => setPrefix(event.target.value)}
                  className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
                />
              </label>

              <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                <span>{t('workbench.tag.template.strategy')}</span>
                <select
                  aria-label={t('workbench.tag.template.strategy')}
                  value={strategy}
                  onChange={(event) =>
                    setStrategy(event.target.value as TagBindingStrategy)
                  }
                  className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
                >
                  <option value="address">
                    {t('workbench.tag.template.strategies.address')}
                  </option>
                  <option value="pointName">
                    {t('workbench.tag.template.strategies.pointName')}
                  </option>
                </select>
              </label>
            </div>
          ) : null}
        </div>

        <div
          className="relative flex min-h-0 flex-1 flex-col overflow-hidden rounded-2xl border border-slate-800/70 bg-gradient-to-b from-slate-900/95 to-slate-950 shadow-[0_1px_0_0_rgba(255,255,255,0.05)_inset]"
          data-testid="tag-master-surface"
        >
          <div
            aria-hidden="true"
            className="pointer-events-none absolute inset-x-0 top-0 h-px bg-gradient-to-r from-transparent via-cyan-500/35 to-transparent"
          />
          <div
            aria-hidden="true"
            className="pointer-events-none absolute -right-20 -top-28 h-56 w-56 rounded-full bg-cyan-500/[0.07] blur-3xl"
          />
          <div className="relative flex min-h-0 flex-1 flex-col divide-y divide-slate-800/70">
            <header className="shrink-0 space-y-1 px-4 pb-3 pt-4">
              <p className="text-[10px] font-semibold uppercase tracking-[0.2em] text-cyan-400/85">
                {t('workbench.tag.master.eyebrow')}
              </p>
              <h3 className="text-lg font-semibold tracking-tight text-slate-50">
                {t('workbench.tag.master.title')}
              </h3>
              <p className="line-clamp-2 text-xs leading-relaxed text-slate-500">
                {t('workbench.tag.master.description')}
              </p>
            </header>

            <dl className="flex min-h-[4.5rem] shrink-0 items-stretch bg-slate-950/40">
              <div className="flex min-w-0 flex-1 flex-col justify-center px-3 py-2.5 text-center sm:px-4">
                <dt className="text-[9px] font-medium uppercase tracking-[0.12em] text-slate-500">
                  {t('workbench.tag.master.metrics.total')}
                </dt>
                <dd className="mt-0.5">
                  <span
                    className="font-mono text-xl font-semibold tabular-nums text-slate-50"
                    data-testid="tag-master-total"
                    title={t('workbench.tag.master.metrics.totalFootnote')}
                  >
                    {tags.length}
                  </span>
                </dd>
              </div>
              <div
                aria-hidden="true"
                className="w-px shrink-0 bg-gradient-to-b from-transparent via-slate-600/50 to-transparent"
              />
              <div className="flex min-w-0 flex-1 flex-col justify-center bg-emerald-500/[0.06] px-3 py-2.5 text-center sm:px-4">
                <dt className="text-[9px] font-medium uppercase tracking-[0.12em] text-emerald-400/90">
                  {t('workbench.tag.master.metrics.linked')}
                </dt>
                <dd className="mt-0.5">
                  <span
                    className="font-mono text-xl font-semibold tabular-nums text-emerald-100"
                    data-testid="tag-master-linked"
                    title={t('workbench.tag.master.metrics.linkedFootnote')}
                  >
                    {linkedTagCount}
                  </span>
                </dd>
              </div>
              <div
                aria-hidden="true"
                className="w-px shrink-0 bg-gradient-to-b from-transparent via-slate-600/50 to-transparent"
              />
              <div className="flex min-w-0 flex-1 flex-col justify-center px-3 py-2.5 text-center sm:px-4">
                <dt className="text-[9px] font-medium uppercase tracking-[0.12em] text-violet-300/75">
                  {t('workbench.tag.master.metrics.unused')}
                </dt>
                <dd className="mt-0.5">
                  <span
                    className="font-mono text-xl font-semibold tabular-nums text-slate-100"
                    data-testid="tag-master-unused"
                    title={t('workbench.tag.master.metrics.unusedFootnote')}
                  >
                    {unusedTagCount}
                  </span>
                </dd>
              </div>
            </dl>

            <div className="shrink-0 bg-slate-950/25 px-4 py-3">
              <input
                aria-label={t('workbench.tag.master.search')}
                value={tagLibraryQuery}
                onChange={(event) => setTagLibraryQuery(event.target.value)}
                placeholder={t('workbench.tag.master.searchPlaceholder')}
                className="w-full rounded-lg border-0 bg-slate-900/70 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-600 ring-1 ring-slate-700/60 transition focus:outline-none focus:ring-2 focus:ring-cyan-500/30"
              />
            </div>

            <div className="shrink-0 space-y-2 bg-gradient-to-r from-cyan-950/25 via-slate-950/40 to-slate-950/50 px-4 py-3">
              <p className="text-[9px] font-semibold uppercase tracking-[0.14em] text-cyan-300/85">
                {t('workbench.tag.master.quickCreate')}
              </p>
              <div className="flex flex-col gap-2 sm:flex-row sm:items-end">
                <label className="min-w-0 flex-1 space-y-1">
                  <span className="sr-only">{t('workbench.tag.master.fields.key')}</span>
                  <input
                    aria-label={t('workbench.tag.master.fields.key')}
                    value={tagDraftKey}
                    onChange={(event) => setTagDraftKey(event.target.value)}
                    placeholder={t('workbench.tag.master.fields.key')}
                    className="w-full rounded-lg border-0 bg-slate-900/80 px-2.5 py-2 text-sm text-slate-100 placeholder:text-slate-600 ring-1 ring-slate-700/50 focus:outline-none focus:ring-2 focus:ring-cyan-500/25"
                  />
                </label>
                <label className="min-w-0 flex-1 space-y-1">
                  <span className="sr-only">{t('workbench.tag.master.fields.displayName')}</span>
                  <input
                    aria-label={t('workbench.tag.master.fields.displayName')}
                    value={tagDraftDisplayName}
                    onChange={(event) => setTagDraftDisplayName(event.target.value)}
                    placeholder={t('workbench.tag.master.fields.displayName')}
                    className="w-full rounded-lg border-0 bg-slate-900/80 px-2.5 py-2 text-sm text-slate-100 placeholder:text-slate-600 ring-1 ring-slate-700/50 focus:outline-none focus:ring-2 focus:ring-cyan-500/25"
                  />
                </label>
              </div>
              <div className="flex gap-2">
                <label className="min-w-0 flex-1">
                  <span className="sr-only">{t('workbench.tag.master.fields.dataType')}</span>
                  <select
                    aria-label={t('workbench.tag.master.fields.dataType')}
                    value={tagDraftDataType}
                    onChange={(event) => setTagDraftDataType(event.target.value as DataType)}
                    className="h-[38px] w-full rounded-lg border-0 bg-slate-900/80 px-2.5 text-sm text-slate-100 ring-1 ring-slate-700/50 focus:outline-none focus:ring-2 focus:ring-cyan-500/25"
                  >
                    {tagDataTypeOptions.map((dataType) => (
                      <option key={dataType} value={dataType}>
                        {dataType}
                      </option>
                    ))}
                  </select>
                </label>
                <button
                  type="button"
                  data-testid="tag-master-create"
                  onClick={() => void handleCreateStandaloneTag()}
                  disabled={!tagDraftKey.trim() || createTagMutation.isPending}
                  className="h-[38px] shrink-0 rounded-lg bg-gradient-to-r from-cyan-500 to-cyan-400 px-4 text-sm font-semibold text-slate-950 shadow-sm transition hover:from-cyan-400 hover:to-cyan-300 disabled:cursor-not-allowed disabled:from-slate-700 disabled:to-slate-700 disabled:text-slate-500 disabled:shadow-none"
                >
                  {t('workbench.tag.master.createAction')}
                </button>
              </div>
            </div>

            {tagLibraryFeedback ? (
              <p
                role="status"
                className={`shrink-0 px-4 py-2.5 text-sm ${
                  tagLibraryFeedback.tone === 'success'
                    ? 'bg-emerald-500/10 text-emerald-100'
                    : 'bg-rose-500/10 text-rose-100'
                }`}
              >
                {tagLibraryFeedback.message}
              </p>
            ) : null}

            <div
              className="min-h-0 flex-1 divide-y divide-slate-800/60 overflow-y-auto overscroll-contain"
              data-testid="tag-master-list"
            >
              {filteredTagLibrary.length > 0 ? (
                filteredTagLibrary.map((tag) => {
                  const bindingCount = mappingCountByTagId.get(tag.id) ?? 0;
                  const canDelete = bindingCount === 0;
                  const isLinked = bindingCount > 0;
                  const displayNameTrimmed = tag.display_name?.trim() ?? '';
                  const showSecondaryLabel =
                    displayNameTrimmed.length > 0 &&
                    displayNameTrimmed !== tag.key;

                  return (
                    <article
                      key={tag.id}
                      data-testid={`tag-master-row-${tag.id}`}
                      title={[tag.key, showSecondaryLabel ? displayNameTrimmed : '']
                        .filter(Boolean)
                        .join(' · ')}
                      className={[
                        'group/tag-master-row flex min-h-11 flex-nowrap items-center gap-2.5 px-4 py-2 transition-colors duration-150',
                        'hover:bg-slate-800/35',
                        isLinked ? 'border-l-[3px] border-l-emerald-500/50' : 'border-l-[3px] border-l-transparent',
                      ].join(' ')}
                    >
                      <div className="flex min-w-0 flex-1 items-center gap-2">
                        <p className="shrink-0 truncate font-mono text-sm font-medium tabular-nums text-slate-100">
                          {tag.key}
                        </p>
                        {showSecondaryLabel ? (
                          <>
                            <span
                              className="shrink-0 text-slate-600"
                              aria-hidden
                            >
                              ·
                            </span>
                            <p className="min-w-0 truncate text-xs text-slate-500">
                              {displayNameTrimmed}
                            </p>
                          </>
                        ) : null}
                      </div>

                      <div
                        aria-hidden
                        className="hidden h-5 w-px shrink-0 bg-gradient-to-b from-transparent via-slate-600/45 to-transparent sm:block"
                      />

                      <div className="flex shrink-0 flex-nowrap items-center gap-1.5">
                        <span className="rounded-md bg-slate-800/90 px-1.5 py-0.5 font-mono text-[10px] font-medium text-slate-400">
                          {tag.data_type}
                        </span>
                        <span
                          className={[
                            'whitespace-nowrap rounded-md px-1.5 py-0.5 text-[10px] font-medium tabular-nums',
                            isLinked
                              ? 'bg-emerald-500/15 text-emerald-200/90 ring-1 ring-emerald-500/20'
                              : 'bg-slate-800/60 text-slate-500',
                          ].join(' ')}
                        >
                          {t('workbench.tag.master.bindingCount', { count: bindingCount })}
                        </span>
                        {!canDelete ? (
                          <span className="whitespace-nowrap rounded-md border border-amber-500/25 bg-amber-500/10 px-1.5 py-0.5 text-[10px] font-medium text-amber-200/90">
                            {t('workbench.tag.master.inUse')}
                          </span>
                        ) : null}
                      </div>

                      <button
                        type="button"
                        data-testid={`tag-master-delete-${tag.id}`}
                        onClick={() => void handleDeleteStandaloneTag(tag)}
                        disabled={!canDelete || deleteTagMutation.isPending}
                        className="shrink-0 rounded-md px-2 py-1 text-[11px] font-medium text-rose-300/90 opacity-90 transition hover:bg-rose-500/15 hover:text-rose-200 hover:opacity-100 disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:bg-transparent sm:opacity-70 sm:group-hover/tag-master-row:opacity-100"
                      >
                        {t('workbench.tag.master.deleteAction')}
                      </button>
                    </article>
                  );
                })
              ) : (
                <div className="px-4 py-10 text-center text-sm text-slate-500">
                  {t('workbench.tag.master.empty')}
                </div>
              )}
            </div>
          </div>
        </div>

        <dl className="grid shrink-0 gap-2 sm:grid-cols-3 xl:grid-cols-3">
          <div className="rounded-lg border border-slate-800 bg-slate-900/70 p-3">
            <dt className="text-xs uppercase tracking-[0.18em] text-slate-400">
              {t('workbench.tag.metrics.selected')}
            </dt>
            <dd className="mt-1 text-xl font-semibold text-slate-50">
              {selectedCandidates.length}
            </dd>
          </div>
          <div className="rounded-lg border border-slate-800 bg-slate-900/70 p-3">
            <dt className="text-xs uppercase tracking-[0.18em] text-slate-400">
              {t('workbench.tag.metrics.ready')}
            </dt>
            <dd className="mt-1 text-xl font-semibold text-emerald-200">
              {readyCount}
            </dd>
          </div>
          <div className="rounded-lg border border-slate-800 bg-slate-900/70 p-3">
            <dt className="text-xs uppercase tracking-[0.18em] text-slate-400">
              {t('workbench.tag.metrics.blocked')}
            </dt>
            <dd className="mt-1 text-xl font-semibold text-rose-200">
              {blockedSelectionCount}
            </dd>
          </div>
        </dl>

        {selectedCandidates.length > 0 ? (
          <div
            data-testid="batch-diff-preview"
            className="shrink-0 space-y-3 rounded-2xl border border-slate-800 bg-slate-900/70 p-4"
          >
            <p className="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-300">
              {t('workbench.tag.diff.title')}
            </p>

            <dl className="grid grid-cols-2 gap-2">
              {flowMode === 'create' ? (
                <div className="rounded-xl border border-emerald-500/30 bg-emerald-500/10 p-3">
                  <dt className="text-[11px] uppercase tracking-[0.16em] text-emerald-300">
                    {t('workbench.tag.diff.toCreate')}
                  </dt>
                  <dd
                    className="mt-1 text-lg font-semibold text-emerald-100"
                    data-testid="diff-to-create-count"
                  >
                    {batchDiffPreview.toCreate.length}
                  </dd>
                </div>
              ) : (
                <div className="rounded-xl border border-cyan-500/30 bg-cyan-500/10 p-3">
                  <dt className="text-[11px] uppercase tracking-[0.16em] text-cyan-300">
                    {t('workbench.tag.diff.toBind')}
                  </dt>
                  <dd
                    className="mt-1 text-lg font-semibold text-cyan-100"
                    data-testid="diff-to-bind-count"
                  >
                    {batchDiffPreview.toBind.length}
                  </dd>
                </div>
              )}
              <div className="rounded-xl border border-slate-700 bg-slate-950/60 p-3">
                <dt className="text-[11px] uppercase tracking-[0.16em] text-slate-400">
                  {t('workbench.tag.diff.skipped')}
                </dt>
                <dd
                  className="mt-1 text-lg font-semibold text-slate-300"
                  data-testid="diff-skipped-count"
                >
                  {batchDiffPreview.skipped.length}
                </dd>
              </div>
            </dl>

            {(flowMode === 'create'
              ? batchDiffPreview.toCreate
              : batchDiffPreview.toBind
            ).length > 0 ? (
              <ul className="max-h-40 space-y-1 overflow-y-auto text-xs">
                {(flowMode === 'create'
                  ? batchDiffPreview.toCreate
                  : batchDiffPreview.toBind
                ).map((item) => (
                  <li
                    key={item.pointId}
                    className="flex items-center justify-between rounded-lg bg-slate-950/50 px-2 py-1.5"
                  >
                    <span className="truncate text-slate-300">{item.pointName}</span>
                    <span className="ml-2 shrink-0 font-mono text-cyan-200">{item.tagKey}</span>
                  </li>
                ))}
              </ul>
            ) : null}

            {batchDiffPreview.skipped.length > 0 ? (
              <div className="space-y-1">
                <p className="text-[11px] uppercase tracking-[0.16em] text-slate-500">
                  {t('workbench.tag.diff.skippedItems')}
                </p>
                <ul className="max-h-24 space-y-1 overflow-y-auto text-xs">
                  {batchDiffPreview.skipped.map((item) => (
                    <li
                      key={item.pointId}
                      className="flex items-center justify-between rounded-lg bg-amber-500/5 px-2 py-1.5"
                    >
                      <span className="truncate text-slate-400">{item.pointName}</span>
                      <span className="ml-2 shrink-0 text-amber-300">
                        {t(`workbench.tag.diff.skipReason.${item.reason}`)}
                      </span>
                    </li>
                  ))}
                </ul>
              </div>
            ) : null}
          </div>
        ) : null}

        {blockedSelectionCount > 0 ? (
          <p
            role="status"
            aria-live="polite"
            className="shrink-0 rounded-xl border border-rose-500/30 bg-rose-500/10 px-3 py-2 text-sm text-rose-100"
          >
            {t('workbench.tag.conflicts.summary')}
          </p>
        ) : null}

        {selectedCandidates.length > 0 ? (
          <div className="shrink-0 space-y-2">
            {selectedBoundItems.length > 0 ? (
                <button
                  type="button"
                  data-testid="tag-batch-unbind"
                  onClick={() => void handleBatchUnbind()}
                  disabled={isTagBoardBusy}
                  className="w-full rounded-xl border border-rose-500/30 bg-rose-500/10 px-4 py-3 text-sm font-semibold text-rose-100 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  {t('workbench.tag.actions.unbindSelected')}
                </button>
            ) : null}
            <button
              type="button"
              onClick={() => void handleBatchBind()}
                disabled={
                  readyCount === 0
                  || blockedSelectionCount > 0
                  || isTagBoardBusy
                  || createTagMutation.isPending
                  || createMappingMutation.isPending
                }
              className="w-full rounded-xl bg-cyan-500 px-4 py-3 text-sm font-semibold text-slate-950 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {t('workbench.tag.actions.bind')}
            </button>
          </div>
        ) : null}

        {batchSummary ? (
          <div
            role="status"
            aria-live="polite"
            className="shrink-0 space-y-3 rounded-2xl border border-slate-800 bg-slate-900/70 p-4"
          >
            <p className="text-sm font-semibold text-slate-50">
              {batchSummary.failureCount > 0
                ? t('workbench.tag.results.partialFailure')
                : t('workbench.tag.results.success')}
            </p>

            <dl className="grid grid-cols-2 gap-2 text-xs">
              <div className="rounded-lg border border-emerald-500/30 bg-emerald-500/10 p-2">
                <dt className="text-emerald-300">
                  {batchSummary.mode === 'deletePoints'
                    ? t('workbench.tag.results.deletedPoints')
                    : batchSummary.mode === 'unbind'
                      ? t('workbench.tag.results.unbound')
                      : flowMode === 'existing'
                        ? t('workbench.tag.results.linked')
                        : t('workbench.tag.results.created')}
                </dt>
                <dd
                  className="mt-1 text-lg font-semibold text-emerald-100"
                  data-testid={
                    batchSummary.mode === 'deletePoints'
                      ? 'result-deleted-points-count'
                      : flowMode === 'existing'
                        ? 'result-linked-count'
                        : 'result-created-count'
                  }
                >
                  {batchSummary.mode === 'deletePoints'
                    ? batchSummary.deletedPointCount
                    : batchSummary.mode === 'unbind'
                      ? batchSummary.unboundCount
                      : flowMode === 'existing'
                        ? batchSummary.linkedCount
                        : batchSummary.createdCount}
                </dd>
              </div>
              <div className="rounded-lg border border-rose-500/30 bg-rose-500/10 p-2">
                <dt className="text-rose-300">{t('workbench.tag.results.failedLabel')}</dt>
                <dd
                  className="mt-1 text-lg font-semibold text-rose-100"
                  data-testid="result-failed-count"
                >
                  {batchSummary.failureCount}
                </dd>
              </div>
              {batchSummary.skippedCount > 0 ? (
                <div className="rounded-lg border border-slate-700 bg-slate-950/60 p-2">
                  <dt className="text-slate-400">{t('workbench.tag.results.skippedLabel')}</dt>
                  <dd
                    className="mt-1 text-lg font-semibold text-slate-300"
                    data-testid="result-skipped-count"
                  >
                    {batchSummary.skippedCount}
                  </dd>
                </div>
              ) : null}
            </dl>

            {batchSummary.failures.length > 0 ? (
              <ul className="space-y-2 text-sm text-rose-100">
                {batchSummary.failures.map((failure) => (
                  <li
                    key={`${failure.pointId}-${failure.stage}`}
                    className="rounded-xl bg-rose-500/10 px-3 py-2"
                  >
                    <span className="font-medium">{failure.tagKey}</span>
                    <span className="ml-2">{failure.error}</span>
                  </li>
                ))}
              </ul>
            ) : null}
          </div>
        ) : null}
      </aside>
    </section>
  );
}
