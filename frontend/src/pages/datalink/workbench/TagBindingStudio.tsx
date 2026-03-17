import { useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useDevicesQuery } from '../../../hooks/datalink/useDevices';
import {
  useCreateMappingMutation,
  useMappingsQuery,
} from '../../../hooks/datalink/useMappings';
import { usePointsQuery } from '../../../hooks/datalink/usePoints';
import {
  useCreateTagMutation,
  useTagsQuery,
} from '../../../hooks/datalink/useTags';
import type { Device } from '../../../types/datalink';
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
  createdCount: number;
  linkedCount: number;
  skippedCount: number;
  failureCount: number;
  failures: TagBindingFailure[];
};

type TagBindingFlowMode = 'create' | 'existing';
type TagBindingStatusFilter = 'all' | 'unbound' | 'partial' | 'bound';

function getSelectedDevice(devices: Device[], selectedDeviceId: string | null) {
  return devices.find((device) => device.id === selectedDeviceId) ?? null;
}

function getErrorMessage(error: unknown, fallback: string): string {
  return error instanceof Error ? error.message : fallback;
}

function isNonNull<T>(value: T | null): value is T {
  return value !== null;
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
  const { selectedDeviceId, setActiveStep, setInspectorSelection, setSelectedDeviceId, sourcePlanningState } = useWorkbench();
  const { data: devices = [] } = useDevicesQuery();
  const { data: tags = [] } = useTagsQuery();
  const { data: mappings = [] } = useMappingsQuery();
  const { data: points = [] } = usePointsQuery(
    selectedDeviceId ? { device_id: selectedDeviceId } : undefined,
  );
  const createTagMutation = useCreateTagMutation();
  const createMappingMutation = useCreateMappingMutation();

  const selectedDevice = getSelectedDevice(devices, selectedDeviceId);
  const mappingByPointId = useMemo(
    () => new Map(mappings.map((mapping) => [mapping.point_id, mapping])),
    [mappings],
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
  const selectionHintKey = batchSummary
    && batchSummary.failureCount === 0
    && (batchSummary.createdCount > 0 || batchSummary.linkedCount > 0)
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

  const handleBatchBind = async () => {
    if (!selectedDevice || blockedSelectionCount > 0) {
      return;
    }

    if (flowMode === 'create') {
      if (createRequests.length === 0) {
        return;
      }

      const results = await Promise.all(
        createRequests.map(async (request) => {
          try {
            const createdTag = await createTagMutation.mutateAsync(request.tagRequest);

            try {
              await createMappingMutation.mutateAsync({
                point_id: request.pointId,
                tag_id: createdTag.id,
                enabled: true,
              });

              return null;
            } catch (error) {
              return {
                pointId: request.pointId,
                tagKey: request.tagKey,
                stage: 'mapping' as const,
                error: getErrorMessage(error, t('workbench.tag.results.mappingFallback')),
              };
            }
          } catch (error) {
            return {
              pointId: request.pointId,
              tagKey: request.tagKey,
              stage: 'tag' as const,
              error: getErrorMessage(error, t('workbench.tag.results.tagFallback')),
            };
          }
        }),
      );

      const failures = results.filter(isNonNull);

      setBatchSummary({
        createdCount: createRequests.length - failures.length,
        linkedCount: 0,
        skippedCount: batchDiffPreview.skipped.length,
        failureCount: failures.length,
        failures,
      });
      return;
    }

    if (existingRequests.length === 0) {
      return;
    }

    const results = await Promise.all(
      existingRequests.map(async (request) => {
        try {
          await createMappingMutation.mutateAsync({
            point_id: request.pointId,
            tag_id: request.tagId,
            enabled: true,
          });
          return null;
        } catch (error) {
          return {
            pointId: request.pointId,
            tagKey: request.tagKey,
            stage: 'mapping' as const,
            error: getErrorMessage(error, t('workbench.tag.results.mappingFallback')),
          };
        }
      }),
    );

    const existingFailures = results.filter(isNonNull);

    setBatchSummary({
      createdCount: 0,
      linkedCount: existingRequests.length - existingFailures.length,
      skippedCount: batchDiffPreview.skipped.length,
      failureCount: existingFailures.length,
      failures: existingFailures,
    });
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
            const selectedExistingTag =
              candidate.existingTagOptions.find(
                (option) => option.id === existingTagSelections[candidate.pointId],
              ) ?? null;

            return (
              <article
                key={candidate.pointId}
                data-layout="row"
                data-testid={`tag-candidate-${candidate.pointId}`}
                onClick={() =>
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
                    })
                }
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

                      {flowMode === 'create' ? (
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

                    {flowMode === 'existing' ? (
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
          <div className="space-y-2">
            <p className="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-300">
              {t('workbench.tag.template.eyebrow')}
            </p>
            <div>
              <h3 className="text-xl font-semibold text-slate-50">
                {t('workbench.tag.template.title')}
              </h3>
              <p className="mt-1 text-sm text-slate-300">
                {t('workbench.tag.template.description')}
              </p>
            </div>
          </div>

          <div className="grid gap-3">
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
          <button
            type="button"
            onClick={() => void handleBatchBind()}
            disabled={
              readyCount === 0
              || blockedSelectionCount > 0
              || createTagMutation.isPending
              || createMappingMutation.isPending
            }
            className="w-full rounded-xl bg-cyan-500 px-4 py-3 text-sm font-semibold text-slate-950 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {t('workbench.tag.actions.bind')}
          </button>
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
                  {flowMode === 'existing'
                    ? t('workbench.tag.results.linked')
                    : t('workbench.tag.results.created')}
                </dt>
                <dd
                  className="mt-1 text-lg font-semibold text-emerald-100"
                  data-testid={flowMode === 'existing' ? 'result-linked-count' : 'result-created-count'}
                >
                  {flowMode === 'existing' ? batchSummary.linkedCount : batchSummary.createdCount}
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
