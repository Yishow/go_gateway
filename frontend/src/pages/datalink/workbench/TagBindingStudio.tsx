import { useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useDevicesQuery } from '../../../hooks/datalink/useDevices';
import { useMappingsQuery, useCreateMappingMutation } from '../../../hooks/datalink/useMappings';
import { usePointsQuery } from '../../../hooks/datalink/usePoints';
import { useCreateTagMutation, useTagsQuery } from '../../../hooks/datalink/useTags';
import type { Device } from '../../../types/datalink';
import { useWorkbench } from './WorkbenchProvider';
import {
  buildTagBindingCandidates,
  buildTagBindingRequests,
  type TagBindingStrategy,
} from './tagBindingModel';

type TagBindingFailure = {
  pointId: string;
  tagKey: string;
  stage: 'tag' | 'mapping';
  error: string;
};

type TagBindingBatchSummary = {
  successCount: number;
  failureCount: number;
  failures: TagBindingFailure[];
};

function getSelectedDevice(devices: Device[], selectedDeviceId: string | null) {
  return devices.find((device) => device.id === selectedDeviceId) ?? null;
}

function getErrorMessage(error: unknown, fallback: string): string {
  return error instanceof Error ? error.message : fallback;
}

export function TagBindingStudio() {
  const { t } = useTranslation();
  const { selectedDeviceId, setSelectedDeviceId } = useWorkbench();
  const { data: devices = [] } = useDevicesQuery();
  const { data: tags = [] } = useTagsQuery();
  const { data: mappings = [] } = useMappingsQuery();
  const { data: points = [] } = usePointsQuery(
    selectedDeviceId ? { device_id: selectedDeviceId } : undefined,
  );
  const createTagMutation = useCreateTagMutation();
  const createMappingMutation = useCreateMappingMutation();

  const selectedDevice = getSelectedDevice(devices, selectedDeviceId);
  const [prefix, setPrefix] = useState('TAG');
  const [strategy, setStrategy] = useState<TagBindingStrategy>('address');
  const [selectedPointIds, setSelectedPointIds] = useState<string[]>([]);
  const [batchSummary, setBatchSummary] = useState<TagBindingBatchSummary | null>(null);
  const pointIdsKey = useMemo(
    () => points.map((point) => point.id).join('|'),
    [points],
  );
  const allPointIds = useMemo(
    () => (pointIdsKey ? pointIdsKey.split('|') : []),
    [pointIdsKey],
  );

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
    setSelectedPointIds(allPointIds);
    setBatchSummary(null);
  }, [allPointIds, selectedDeviceId]);

  const selectedCandidates = useMemo(
    () =>
      candidates.filter((candidate) => selectedPointIds.includes(candidate.pointId)),
    [candidates, selectedPointIds],
  );

  const bindableRequests = useMemo(
    () => buildTagBindingRequests(candidates, selectedPointIds),
    [candidates, selectedPointIds],
  );

  const blockedSelectionCount = selectedCandidates.filter(
    (candidate) => candidate.conflict || candidate.alreadyLinked,
  ).length;

  const handleTogglePoint = (pointId: string) => {
    setSelectedPointIds((previous) =>
      previous.includes(pointId)
        ? previous.filter((id) => id !== pointId)
        : [...previous, pointId],
    );
    setBatchSummary(null);
  };

  const handleSelectAll = () => {
    setSelectedPointIds(candidates.map((candidate) => candidate.pointId));
    setBatchSummary(null);
  };

  const handleSelectBindable = () => {
    setSelectedPointIds(
      candidates
        .filter((candidate) => !candidate.conflict && !candidate.alreadyLinked)
        .map((candidate) => candidate.pointId),
    );
    setBatchSummary(null);
  };

  const handleClearSelection = () => {
    setSelectedPointIds([]);
    setBatchSummary(null);
  };

  const handleBatchBind = async () => {
    if (
      !selectedDevice ||
      bindableRequests.length === 0 ||
      blockedSelectionCount > 0
    ) {
      return;
    }

    const results = await Promise.all(
      bindableRequests.map(async (request) => {
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

    const failures = results.filter(
      (result): result is TagBindingFailure => result !== null,
    );

    setBatchSummary({
      successCount: bindableRequests.length - failures.length,
      failureCount: failures.length,
      failures,
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
      </section>
    );
  }

  return (
    <section className="grid gap-6 xl:grid-cols-[minmax(0,1.35fr)_minmax(320px,0.95fr)]">
      <div className="space-y-6 rounded-2xl border border-slate-800 bg-slate-950/40 p-5">
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
        </div>

        <div className="grid gap-3">
          {candidates.map((candidate) => {
            const selected = selectedPointIds.includes(candidate.pointId);

            return (
              <label
                key={candidate.pointId}
                className={`grid gap-3 rounded-2xl border p-4 transition md:grid-cols-[auto_minmax(0,1fr)_minmax(220px,0.8fr)] ${
                  selected
                    ? 'border-cyan-500/40 bg-cyan-500/5'
                    : 'border-slate-800 bg-slate-900/70'
                }`}
              >
                <span className="pt-1">
                  <input
                    type="checkbox"
                    checked={selected}
                    onChange={() => handleTogglePoint(candidate.pointId)}
                    aria-label={candidate.pointName}
                    className="h-4 w-4 rounded border-slate-700 bg-slate-950 text-cyan-400"
                  />
                </span>

                <span className="space-y-2">
                  <span className="block text-sm font-semibold text-slate-50">
                    {candidate.pointName}
                  </span>
                  <span className="block text-xs text-slate-400">
                    {t('workbench.tag.selection.pointMeta', {
                      address: candidate.pointAddress,
                      dataType: candidate.dataType,
                    })}
                  </span>
                  <span className="flex flex-wrap gap-2">
                    {candidate.alreadyLinked ? (
                      <span className="rounded-full bg-amber-500/10 px-2 py-1 text-xs text-amber-200">
                        {t('workbench.tag.badges.alreadyLinked')}
                      </span>
                    ) : null}
                    {candidate.conflictReason === 'existing-key' ? (
                      <span className="rounded-full bg-rose-500/10 px-2 py-1 text-xs text-rose-200">
                        {t('workbench.tag.badges.existingKey')}
                      </span>
                    ) : null}
                    {candidate.conflictReason === 'duplicate-preview' ? (
                      <span className="rounded-full bg-rose-500/10 px-2 py-1 text-xs text-rose-200">
                        {t('workbench.tag.badges.duplicatePreview')}
                      </span>
                    ) : null}
                  </span>
                </span>

                <span
                  data-testid={`tag-preview-${candidate.pointId}`}
                  data-conflict={candidate.conflict ? 'true' : 'false'}
                  className={`rounded-xl border px-3 py-2 text-sm font-medium ${
                    candidate.conflict
                      ? 'border-rose-500/40 bg-rose-500/10 text-rose-100'
                      : 'border-slate-800 bg-slate-950/80 text-cyan-100'
                  }`}
                >
                  {candidate.previewKey}
                </span>
              </label>
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

        <dl className="grid gap-3 sm:grid-cols-3 xl:grid-cols-1">
          <div className="rounded-xl border border-slate-800 bg-slate-900/70 p-4">
            <dt className="text-xs uppercase tracking-[0.18em] text-slate-400">
              {t('workbench.tag.metrics.selected')}
            </dt>
            <dd className="mt-2 text-2xl font-semibold text-slate-50">
              {selectedCandidates.length}
            </dd>
          </div>
          <div className="rounded-xl border border-slate-800 bg-slate-900/70 p-4">
            <dt className="text-xs uppercase tracking-[0.18em] text-slate-400">
              {t('workbench.tag.metrics.ready')}
            </dt>
            <dd className="mt-2 text-2xl font-semibold text-emerald-200">
              {bindableRequests.length}
            </dd>
          </div>
          <div className="rounded-xl border border-slate-800 bg-slate-900/70 p-4">
            <dt className="text-xs uppercase tracking-[0.18em] text-slate-400">
              {t('workbench.tag.metrics.blocked')}
            </dt>
            <dd className="mt-2 text-2xl font-semibold text-rose-200">
              {blockedSelectionCount}
            </dd>
          </div>
        </dl>

        {blockedSelectionCount > 0 ? (
          <p
            role="status"
            aria-live="polite"
            className="rounded-xl border border-rose-500/30 bg-rose-500/10 px-3 py-2 text-sm text-rose-100"
          >
            {t('workbench.tag.conflicts.summary')}
          </p>
        ) : null}

        <button
          type="button"
          onClick={() => void handleBatchBind()}
          disabled={
            bindableRequests.length === 0 ||
            blockedSelectionCount > 0 ||
            createTagMutation.isPending ||
            createMappingMutation.isPending
          }
          className="w-full rounded-xl bg-cyan-500 px-4 py-3 text-sm font-semibold text-slate-950 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {t('workbench.tag.actions.bind')}
        </button>

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
            <p className="text-sm text-slate-300">
              {t('workbench.tag.results.summary', batchSummary)}
            </p>

            {batchSummary.failures.length > 0 ? (
              <ul className="space-y-2 text-sm text-rose-100">
                {batchSummary.failures.map((failure) => (
                  <li key={`${failure.pointId}-${failure.stage}`} className="rounded-xl bg-rose-500/10 px-3 py-2">
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
