import { useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useDevicesQuery } from '../../../hooks/datalink/useDevices';
import {
  useCreateMappingMutation,
  useDeleteMappingMutation,
  useMappingsQuery,
} from '../../../hooks/datalink/useMappings';
import { usePointsQuery } from '../../../hooks/datalink/usePoints';
import {
  useCreateTagMutation,
  useDeleteTagMutation,
  useTagsQuery,
} from '../../../hooks/datalink/useTags';
import { tagAPI } from '../../../services/datalink';
import type { DataType, Device, Tag } from '../../../types/datalink';
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
  stage: 'tag' | 'mapping';
  error: string;
};

type TagBindingBatchSummary = {
  mode: 'bind' | 'unbind';
  createdCount: number;
  linkedCount: number;
  unboundCount: number;
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

const tagDataTypeOptions: DataType[] = [
  'bool',
  'int16',
  'int32',
  'int64',
  'uint16',
  'uint32',
  'uint64',
  'float32',
  'float64',
  'string',
];

function getSelectedDevice(devices: Device[], selectedDeviceId: string | null) {
  return devices.find((device) => device.id === selectedDeviceId) ?? null;
}

function getErrorMessage(error: unknown, fallback: string): string {
  return error instanceof Error ? error.message : fallback;
}

function formatCandidateValue(value: unknown) {
  if (value === null || value === undefined || value === '') {
    return '—';
  }

  return String(value);
}

function getStatusToneClass(status: TagBindingCandidate['bindingStatus']) {
  switch (status) {
    case 'bound':
      return 'border-emerald-500/40 bg-emerald-500/10 text-emerald-200';
    case 'partial':
      return 'border-amber-500/40 bg-amber-500/10 text-amber-200';
    case 'unbound':
      return 'border-slate-700 bg-slate-950/70 text-slate-300';
  }
}

export function TagBindingStudio() {
  const { t } = useTranslation();
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
  const selectionHintKey = batchSummary
    && batchSummary.failureCount === 0
    && (
      batchSummary.createdCount > 0
      || batchSummary.linkedCount > 0
      || batchSummary.unboundCount > 0
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
    if (selectedBoundItems.length === 0 || isBatchUnbinding) {
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
        skippedCount: 0,
        failureCount: failures.length,
        failures,
      });
      setSelectedPointIds((previous) => previous.filter((id) => !succeededPointIds.includes(id)));
    } finally {
      setIsBatchUnbinding(false);
    }
  };

  if (!selectedDevice) {
    return (
      <section className="space-y-6 rounded-2xl border border-dashed border-slate-700 bg-slate-950/40 p-6">
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
      <section className="space-y-6 rounded-2xl border border-dashed border-slate-700 bg-slate-950/40 p-6">
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
    <section className="grid gap-6 xl:grid-cols-[minmax(0,1.45fr)_minmax(340px,0.95fr)]">
      <div className="space-y-6 rounded-2xl border border-slate-800 bg-slate-950/40 p-5">
        <div className="space-y-4">
          <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
            <div className="space-y-2">
              <p className="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-300">
                {t('workbench.tag.selection.eyebrow')}
              </p>
              <div>
                <h2 className="text-2xl font-semibold text-slate-50">
                  {t('workbench.tag.selection.title')}
                </h2>
                <p className="mt-1 max-w-2xl text-sm text-slate-300">
                  {t('workbench.tag.selection.description')}
                </p>
              </div>
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
          </div>

          <p className="text-sm text-slate-300">{t(flowModeHintKey)}</p>

          <div className="grid gap-3 lg:grid-cols-[minmax(0,1fr)_220px]">
            <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
              <span>{t('workbench.tag.board.search')}</span>
              <input
                aria-label={t('workbench.tag.board.search')}
                value={searchQuery}
                onChange={(event) => setSearchQuery(event.target.value)}
                className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
              />
            </label>
            <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
              <span>{t('workbench.tag.board.statusFilter')}</span>
              <select
                aria-label={t('workbench.tag.board.statusFilter')}
                value={statusFilter}
                onChange={(event) =>
                  setStatusFilter(event.target.value as TagBindingStatusFilter)
                }
                className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
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
            </label>
          </div>

          <div className="flex flex-wrap items-center gap-2">
            <button
              type="button"
              onClick={handleSelectAll}
              className="rounded-lg border border-slate-800 px-3 py-2 text-sm text-slate-300"
            >
              {t('workbench.tag.actions.selectAll')}
            </button>
            <button
              type="button"
              onClick={handleSelectBindable}
              className="rounded-lg border border-slate-800 px-3 py-2 text-sm text-slate-300"
            >
              {t('workbench.tag.actions.selectBindable')}
            </button>
            <button
              type="button"
              onClick={handleClearSelection}
              className="rounded-lg border border-slate-800 px-3 py-2 text-sm text-slate-300"
            >
              {t('workbench.tag.actions.clearSelection')}
            </button>
          </div>

          <p className="text-sm text-slate-300">{t(selectionHintKey)}</p>
        </div>

        <div className="space-y-2" data-testid="tag-candidate-board" data-layout="row-board">
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
                data-layout="row"
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
                className={`rounded-xl border px-4 py-3 transition ${
                  selected
                    ? 'border-cyan-500/40 bg-cyan-500/5'
                    : 'border-slate-800 bg-slate-900/70'
                }`}
              >
                <div className="grid gap-3 lg:grid-cols-[auto_minmax(0,1.2fr)_minmax(96px,0.45fr)_minmax(96px,0.45fr)_minmax(220px,0.9fr)] lg:items-center">
                  <span className="pt-1">
                    <input
                      type="checkbox"
                      checked={selected}
                      onClick={(event) => event.stopPropagation()}
                      onChange={() => handleTogglePoint(candidate.pointId)}
                      aria-label={candidate.pointName}
                      className="h-4 w-4 rounded border-slate-700 bg-slate-950 text-cyan-400"
                    />
                  </span>

                  <div className="min-w-0 space-y-1">
                    <p className="text-sm font-semibold text-slate-50">
                      {candidate.pointName}
                    </p>
                    <p className="text-xs text-slate-400">
                      {t('workbench.tag.selection.pointMeta', {
                        address: candidate.pointAddress,
                        dataType: candidate.dataType,
                      })}
                    </p>
                    <p className="text-xs text-slate-500">
                      {t('workbench.tag.board.span', {
                        cells: candidate.cellSpan,
                        bitWidth: candidate.bitWidth,
                      })}
                    </p>
                  </div>

                  <div className="space-y-1">
                    <p className="text-[11px] uppercase tracking-[0.16em] text-slate-500">
                      {t('workbench.tag.board.rawValue')}
                    </p>
                    <p
                      className="text-sm font-medium text-slate-100"
                      data-testid={`tag-raw-${candidate.pointId}`}
                    >
                      {formatCandidateValue(candidate.rawValue)}
                    </p>
                  </div>
                  <div className="space-y-1">
                    <p className="text-[11px] uppercase tracking-[0.16em] text-slate-500">
                      {t('workbench.tag.board.transformedValue')}
                    </p>
                    <p className="text-sm font-medium text-slate-100">
                      {formatCandidateValue(candidate.transformedValue)}
                    </p>
                  </div>

                  <div className="min-w-0 space-y-2">
                    <div className="flex flex-wrap items-center justify-between gap-2">
                      <span
                        className={`inline-flex rounded-full border px-2 py-1 text-xs font-medium ${getStatusToneClass(
                          candidate.bindingStatus,
                        )}`}
                        data-testid={`tag-status-${candidate.pointId}`}
                      >
                        {t(`workbench.tag.board.status.${candidate.bindingStatus}`)}
                      </span>

                      {candidate.alreadyLinked && boundTag ? (
                        <div className="flex items-center gap-2">
                          <span
                            data-testid={`tag-bound-label-${candidate.pointId}`}
                            className="rounded-lg border border-emerald-500/30 bg-emerald-500/10 px-3 py-1.5 text-xs font-medium text-emerald-200"
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
                            disabled={isAnyUnbindPending}
                            className="rounded-lg border border-rose-500/30 bg-rose-500/10 px-2.5 py-1.5 text-xs font-medium text-rose-200 transition hover:bg-rose-500/20 disabled:opacity-50"
                          >
                            {t('workbench.tag.actions.unbind')}
                          </button>
                        </div>
                      ) : flowMode === 'create' ? (
                        <div
                          data-testid={`tag-preview-${candidate.pointId}`}
                          data-conflict={candidate.conflict ? 'true' : 'false'}
                          className={`rounded-lg border px-3 py-2 text-sm font-medium ${
                            candidate.conflict
                              ? 'border-rose-500/40 bg-rose-500/10 text-rose-100'
                              : 'border-slate-800 bg-slate-950/80 text-cyan-100'
                          }`}
                        >
                          {candidate.previewKey}
                        </div>
                      ) : null}
                    </div>

                    {flowMode === 'existing' && !candidate.alreadyLinked ? (
                      <label className="block space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                        <span>{t('workbench.tag.board.existingTag')}</span>
                        <select
                          data-testid={`existing-tag-select-${candidate.pointId}`}
                          value={existingTagSelections[candidate.pointId] ?? ''}
                          onClick={(event) => event.stopPropagation()}
                          onChange={(event) =>
                            handleExistingTagSelection(
                              candidate.pointId,
                              event.target.value,
                            )
                          }
                          className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
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
                      </label>
                    ) : null}
                  </div>
                </div>
              </article>
            );
          })}
        </div>
      </div>

      <aside className="space-y-6 rounded-2xl border border-slate-800 bg-slate-950/40 p-5">
        <div className="space-y-4">
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
          className="space-y-4 rounded-2xl border border-slate-800 bg-slate-900/70 p-4"
          data-testid="tag-master-surface"
        >
          <div className="space-y-1">
            <p className="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-300">
              {t('workbench.tag.master.eyebrow')}
            </p>
            <h3 className="text-xl font-semibold text-slate-50">
              {t('workbench.tag.master.title')}
            </h3>
            <p className="text-sm text-slate-300">
              {t('workbench.tag.master.description')}
            </p>
          </div>

          <dl className="grid grid-cols-3 gap-2">
            <div className="rounded-lg border border-slate-800 bg-slate-950/70 p-3">
              <dt className="text-[11px] uppercase tracking-[0.16em] text-slate-500">
                {t('workbench.tag.master.metrics.total')}
              </dt>
              <dd className="mt-1 text-lg font-semibold text-slate-50" data-testid="tag-master-total">
                {tags.length}
              </dd>
            </div>
            <div className="rounded-lg border border-emerald-500/20 bg-emerald-500/5 p-3">
              <dt className="text-[11px] uppercase tracking-[0.16em] text-emerald-300">
                {t('workbench.tag.master.metrics.linked')}
              </dt>
              <dd className="mt-1 text-lg font-semibold text-emerald-100" data-testid="tag-master-linked">
                {linkedTagCount}
              </dd>
            </div>
            <div className="rounded-lg border border-slate-800 bg-slate-950/70 p-3">
              <dt className="text-[11px] uppercase tracking-[0.16em] text-slate-400">
                {t('workbench.tag.master.metrics.unused')}
              </dt>
              <dd className="mt-1 text-lg font-semibold text-slate-100" data-testid="tag-master-unused">
                {unusedTagCount}
              </dd>
            </div>
          </dl>

          <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
            <span>{t('workbench.tag.master.search')}</span>
            <input
              aria-label={t('workbench.tag.master.search')}
              value={tagLibraryQuery}
              onChange={(event) => setTagLibraryQuery(event.target.value)}
              className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
            />
          </label>

          <div className="grid gap-3 rounded-xl border border-slate-800 bg-slate-950/60 p-3">
            <p className="text-xs font-semibold uppercase tracking-[0.16em] text-slate-400">
              {t('workbench.tag.master.quickCreate')}
            </p>
            <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
              <span>{t('workbench.tag.master.fields.key')}</span>
              <input
                aria-label={t('workbench.tag.master.fields.key')}
                value={tagDraftKey}
                onChange={(event) => setTagDraftKey(event.target.value)}
                className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
              />
            </label>
            <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
              <span>{t('workbench.tag.master.fields.displayName')}</span>
              <input
                aria-label={t('workbench.tag.master.fields.displayName')}
                value={tagDraftDisplayName}
                onChange={(event) => setTagDraftDisplayName(event.target.value)}
                className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
              />
            </label>
            <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
              <span>{t('workbench.tag.master.fields.dataType')}</span>
              <select
                aria-label={t('workbench.tag.master.fields.dataType')}
                value={tagDraftDataType}
                onChange={(event) => setTagDraftDataType(event.target.value as DataType)}
                className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
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
              className="rounded-xl border border-cyan-500/30 bg-cyan-500/10 px-4 py-3 text-sm font-semibold text-cyan-100 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {t('workbench.tag.master.createAction')}
            </button>
          </div>

          {tagLibraryFeedback ? (
            <p
              role="status"
              className={`rounded-xl px-3 py-2 text-sm ${
                tagLibraryFeedback.tone === 'success'
                  ? 'border border-emerald-500/30 bg-emerald-500/10 text-emerald-100'
                  : 'border border-rose-500/30 bg-rose-500/10 text-rose-100'
              }`}
            >
              {tagLibraryFeedback.message}
            </p>
          ) : null}

          <div className="space-y-2" data-testid="tag-master-list">
            {filteredTagLibrary.length > 0 ? (
              filteredTagLibrary.map((tag) => {
                const bindingCount = mappingCountByTagId.get(tag.id) ?? 0;
                const canDelete = bindingCount === 0;

                return (
                  <article
                    key={tag.id}
                    data-testid={`tag-master-row-${tag.id}`}
                    className="rounded-xl border border-slate-800 bg-slate-950/60 px-3 py-3"
                  >
                    <div className="flex items-start justify-between gap-3">
                      <div className="min-w-0 space-y-1">
                        <p className="truncate text-sm font-semibold text-slate-50">{tag.key}</p>
                        <p className="truncate text-xs text-slate-400">
                          {tag.display_name || t('workbench.tag.master.displayNameFallback')}
                        </p>
                      </div>
                      <button
                        type="button"
                        data-testid={`tag-master-delete-${tag.id}`}
                        onClick={() => void handleDeleteStandaloneTag(tag)}
                        disabled={!canDelete || deleteTagMutation.isPending}
                        className="rounded-lg border border-rose-500/30 bg-rose-500/10 px-2.5 py-1.5 text-xs font-medium text-rose-100 disabled:cursor-not-allowed disabled:border-slate-700 disabled:bg-slate-950/40 disabled:text-slate-500"
                      >
                        {t('workbench.tag.master.deleteAction')}
                      </button>
                    </div>
                    <div className="mt-3 flex flex-wrap items-center gap-2 text-xs">
                      <span className="rounded-full border border-slate-700 bg-slate-900 px-2 py-1 text-slate-300">
                        {tag.data_type}
                      </span>
                      <span className="rounded-full border border-slate-700 bg-slate-900 px-2 py-1 text-slate-300">
                        {t('workbench.tag.master.bindingCount', { count: bindingCount })}
                      </span>
                      {!canDelete ? (
                        <span className="rounded-full border border-amber-500/30 bg-amber-500/10 px-2 py-1 text-amber-200">
                          {t('workbench.tag.master.inUse')}
                        </span>
                      ) : null}
                    </div>
                  </article>
                );
              })
            ) : (
              <div className="rounded-xl border border-dashed border-slate-700 bg-slate-950/40 px-3 py-4 text-sm text-slate-400">
                {t('workbench.tag.master.empty')}
              </div>
            )}
          </div>
        </div>

        <dl className="grid gap-2 sm:grid-cols-3 xl:grid-cols-3">
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
            className="space-y-3 rounded-2xl border border-slate-800 bg-slate-900/70 p-4"
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
            className="rounded-xl border border-rose-500/30 bg-rose-500/10 px-3 py-2 text-sm text-rose-100"
          >
            {t('workbench.tag.conflicts.summary')}
          </p>
        ) : null}

        {selectedCandidates.length > 0 ? (
          <div className="space-y-2">
            {selectedBoundItems.length > 0 ? (
                <button
                  type="button"
                  data-testid="tag-batch-unbind"
                  onClick={() => void handleBatchUnbind()}
                  disabled={isAnyUnbindPending}
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
                  || isAnyUnbindPending
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
            className="space-y-3 rounded-2xl border border-slate-800 bg-slate-900/70 p-4"
          >
            <p className="text-sm font-semibold text-slate-50">
              {batchSummary.failureCount > 0
                ? t('workbench.tag.results.partialFailure')
                : t('workbench.tag.results.success')}
            </p>

            <dl className="grid grid-cols-2 gap-2 text-xs">
              <div className="rounded-lg border border-emerald-500/30 bg-emerald-500/10 p-2">
                <dt className="text-emerald-300">
                  {batchSummary.mode === 'unbind'
                    ? t('workbench.tag.results.unbound')
                    : flowMode === 'existing'
                      ? t('workbench.tag.results.linked')
                      : t('workbench.tag.results.created')}
                </dt>
                <dd
                  className="mt-1 text-lg font-semibold text-emerald-100"
                  data-testid={flowMode === 'existing' ? 'result-linked-count' : 'result-created-count'}
                >
                  {batchSummary.mode === 'unbind'
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
