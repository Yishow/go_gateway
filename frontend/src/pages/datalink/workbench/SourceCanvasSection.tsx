import { useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import {
  normalizeNamingPrefix,
  SOURCE_PLANNER_ALLOWED_DATA_TYPES,
} from '../../../features/datalink/sourcePlannerContract';
import { useDevicesQuery } from '../../../hooks/datalink/useDevices';
import { useMappingsQuery } from '../../../hooks/datalink/useMappings';
import {
  useCreatePointMutation,
  usePointsQuery,
} from '../../../hooks/datalink/usePoints';
import { useRuntimeStream } from '../../../hooks/datalink/useRuntimeStream';
import { useTagsQuery } from '../../../hooks/datalink/useTags';
import { runtimeAPI } from '../../../services/datalink';
import type { DataType, Device } from '../../../types/datalink';
import { AddressCanvas } from './AddressCanvas';
import { AddressLedger } from './AddressLedger';
import { useWorkbench } from './WorkbenchProvider';
import {
  buildAddressCanvasItems,
  buildCoverageOverviewSegments,
  buildPlannedPointAddresses,
  buildSourceRuleCoverage,
  type SourceRule,
  type SourceValueFormat,
  type SourceViewMode,
} from './sourceCanvasModel';

function getSelectedDevice(devices: Device[], selectedDeviceId: string | null) {
  return devices.find((device) => device.id === selectedDeviceId) ?? null;
}

function getNextRuleId(rules: ReadonlyArray<SourceRule>) {
  const maxIndex = rules.reduce((currentMax, rule) => {
    const suffix = Number.parseInt(rule.id.replace('rule-', ''), 10);
    return Number.isNaN(suffix) ? currentMax : Math.max(currentMax, suffix);
  }, 0);

  return `rule-${maxIndex + 1}`;
}

function moveRule(rules: ReadonlyArray<SourceRule>, ruleId: string, direction: -1 | 1) {
  const fromIndex = rules.findIndex((rule) => rule.id === ruleId);
  const toIndex = fromIndex + direction;

  if (fromIndex === -1 || toIndex < 0 || toIndex >= rules.length) {
    return [...rules];
  }

  const nextRules = [...rules];
  const [movedRule] = nextRules.splice(fromIndex, 1);
  nextRules.splice(toIndex, 0, movedRule);
  return nextRules;
}

function getCoverageSegmentClass(status: ReturnType<typeof buildCoverageOverviewSegments>[number]['status']) {
  switch (status) {
    case 'conflict':
      return 'bg-rose-500/80';
    case 'gap':
      return 'bg-slate-700';
    case 'planned':
      return 'bg-sky-500/80';
    case 'used':
      return 'bg-emerald-500/80';
  }
}

export function SourceCanvasSection() {
  const { t } = useTranslation();
  const {
    selectedDeviceId,
    setFocusedRuleId,
    setInspectorSelection,
    setSelectedDeviceId,
    sourcePlanningState,
    setSourcePlanningState,
  } = useWorkbench();
  const { data: devices = [] } = useDevicesQuery();
  const { data: points = [] } = usePointsQuery(
    selectedDeviceId ? { device_id: selectedDeviceId } : undefined,
  );
  const { data: mappings = [] } = useMappingsQuery();
  const { data: tags = [] } = useTagsQuery();
  const createPointMutation = useCreatePointMutation();
  const runtimeStatusQuery = useQuery({
    queryKey: ['runtime-status', selectedDeviceId],
    queryFn: () => runtimeAPI.getStatus(selectedDeviceId ?? undefined),
    enabled: Boolean(selectedDeviceId),
    staleTime: 10_000,
  });
  const runtimeStream = useRuntimeStream({
    deviceId: selectedDeviceId,
    pointIds: points.map((point) => point.id),
  });
  const selectedDevice = getSelectedDevice(devices, selectedDeviceId);

  const [viewMode, setViewMode] = useState<SourceViewMode>('plan');
  const [valueFormat, setValueFormat] = useState<SourceValueFormat>('decimal');
  const [showAudit, setShowAudit] = useState(false);
  const [freezeLive, setFreezeLive] = useState(false);
  const [frozenLiveValues, setFrozenLiveValues] = useState<
    Readonly<Record<string, { raw_value?: unknown; timestamp?: string }>> | null
  >(null);
  const [startAddress, setStartAddress] = useState('40001');
  const [dataType, setDataType] = useState<DataType>('int16');
  const [count, setCount] = useState(4);
  const [namingPrefix, setNamingPrefix] = useState('SRC');
  const [jumpAddress, setJumpAddress] = useState('');
  const [batchCreateSummary, setBatchCreateSummary] = useState<{
    successCount: number;
    failureCount: number;
  } | null>(null);
  const rules = sourcePlanningState.rules;
  const selectedRuleId = sourcePlanningState.selectedRuleId;
  const selectedAddress = sourcePlanningState.selectedAddress;

  const effectiveLiveValues = freezeLive && frozenLiveValues
    ? frozenLiveValues
    : runtimeStream.liveValues;

  const items = useMemo(() => {
    if (!selectedDevice) {
      return [];
    }

    return buildAddressCanvasItems({
      points,
      rules,
      mappings,
      tags,
      protocol: selectedDevice.protocol,
      liveValues: effectiveLiveValues,
    });
  }, [effectiveLiveValues, mappings, points, rules, selectedDevice, tags]);

  const coverageSegments = useMemo(
    () => buildCoverageOverviewSegments(items),
    [items],
  );

  const plannedPointDefinitions = useMemo(() => {
    if (!selectedDevice) {
      return [];
    }

    const definitions = new Map<
      string,
      { address: string; dataType: DataType; name: string }
    >();

    for (const rule of rules.filter((candidate) => candidate.enabled)) {
      const prefix = normalizeNamingPrefix(rule.namingPrefix);
      const plannedAddresses = buildPlannedPointAddresses({
        startAddress: rule.startAddress,
        count: rule.count,
        dataType: rule.dataType,
        protocol: selectedDevice.protocol,
      });

      plannedAddresses.forEach((address) => {
        if (!definitions.has(address)) {
          definitions.set(address, {
            address,
            dataType: rule.dataType,
            name: `${prefix}_${address}`,
          });
        }
      });
    }

    return [...definitions.values()];
  }, [rules, selectedDevice]);

  const hasConflicts = items.some((item) => item.status === 'conflict');

  const handleApplyPlan = () => {
    if (!selectedDevice) {
      return;
    }

    const ruleId = getNextRuleId(rules);
    const nextRule: SourceRule = {
      id: ruleId,
      startAddress: startAddress.trim(),
      count,
      dataType,
      namingPrefix,
      enabled: true,
      locked: false,
      origin: 'manual',
    };

    setSourcePlanningState((currentState) => ({
      ...currentState,
      rules: [...currentState.rules, nextRule],
      selectedRuleId: ruleId,
      selectedAddress: null,
    }));
    setBatchCreateSummary(null);
    setFocusedRuleId(ruleId);
    setInspectorSelection({ kind: 'rule', ruleId });
  };

  const handleSelectRule = (ruleId: string) => {
    setSourcePlanningState((currentState) => ({
      ...currentState,
      selectedRuleId: ruleId,
    }));
    setFocusedRuleId(ruleId);
    setInspectorSelection({ kind: 'rule', ruleId });
  };

  const handleSelectAddress = (address: string) => {
    const selectedItem = items.find((item) => item.address === address);
    setSourcePlanningState((currentState) => ({
      ...currentState,
      selectedAddress: address,
      selectedRuleId: selectedItem?.primaryRuleId ?? currentState.selectedRuleId,
    }));
    if (selectedItem?.primaryRuleId) {
      setFocusedRuleId(selectedItem.primaryRuleId);
    }
    setInspectorSelection({
      kind: 'span',
      spanAddress: address,
      ruleId: selectedItem?.primaryRuleId ?? undefined,
    });
  };

  const handleToggleRuleEnabled = (ruleId: string) => {
    setSourcePlanningState((currentState) => ({
      ...currentState,
      rules: currentState.rules.map((rule) =>
        rule.id === ruleId ? { ...rule, enabled: !rule.enabled } : rule,
      ),
    }));
  };

  const handleToggleRuleLocked = (ruleId: string) => {
    setSourcePlanningState((currentState) => ({
      ...currentState,
      rules: currentState.rules.map((rule) =>
        rule.id === ruleId ? { ...rule, locked: !rule.locked } : rule,
      ),
    }));
  };

  const handleMoveRule = (ruleId: string, direction: -1 | 1) => {
    setSourcePlanningState((currentState) => ({
      ...currentState,
      rules: moveRule(currentState.rules, ruleId, direction),
    }));
  };

  const handleJumpToAddress = () => {
    const normalizedAddress = jumpAddress.trim();
    if (!normalizedAddress) {
      return;
    }

    const targetItem = items.find((item) => item.address === normalizedAddress);
    if (!targetItem) {
      return;
    }

    handleSelectAddress(targetItem.address);
  };

  const handleToggleFreezeLive = () => {
    if (freezeLive) {
      setFreezeLive(false);
      setFrozenLiveValues(null);
      return;
    }

    setFreezeLive(true);
    setFrozenLiveValues(runtimeStream.liveValues);
  };

  const handleCaptureSnapshot = () => {
    setFreezeLive(true);
    setFrozenLiveValues(runtimeStream.liveValues);
  };

  const handleBatchCreate = async () => {
    if (
      !selectedDevice ||
      plannedPointDefinitions.length === 0 ||
      hasConflicts
    ) {
      return;
    }

    const results = await Promise.allSettled(
      plannedPointDefinitions.map((definition) =>
        createPointMutation.mutateAsync({
          device_id: selectedDevice.id,
          address: definition.address,
          data_type: definition.dataType,
          name: definition.name,
        }),
      ),
    );

    setBatchCreateSummary({
      successCount: results.filter((result) => result.status === 'fulfilled').length,
      failureCount: results.filter((result) => result.status === 'rejected').length,
    });
  };

  if (!selectedDevice) {
    return (
      <section className="space-y-6 rounded-2xl border border-dashed border-slate-700 bg-slate-950/40 p-6">
        <div className="space-y-2">
          <p className="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-300">
            {t('workbench.source.empty.eyebrow')}
          </p>
          <h2 className="text-2xl font-semibold text-slate-50">
            {t('workbench.source.empty.title')}
          </h2>
          <p className="max-w-2xl text-sm text-slate-300">
            {t('workbench.source.empty.description')}
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

  return (
    <section className="space-y-6">
      <div className="grid gap-4 rounded-2xl border border-slate-800 bg-slate-950/40 p-4 md:grid-cols-4">
        <div className="space-y-1 md:col-span-2">
          <p className="text-xs font-semibold uppercase tracking-[0.18em] text-cyan-300">
            {t('workbench.runtime.summary.title')}
          </p>
          <p
            className="text-sm font-medium text-slate-50"
            data-testid="workbench-runtime-status"
          >
            {runtimeStatusQuery.data?.collectors[0]?.status === 'running'
              ? t('workbench.runtime.summary.running')
              : runtimeStatusQuery.data?.collectors[0]?.status === 'warning'
                ? t('workbench.runtime.summary.warning')
                : runtimeStatusQuery.data?.collectors[0]?.status === 'error'
                  ? t('workbench.runtime.summary.error')
                  : t('workbench.runtime.summary.idle')}
          </p>
          <p className="text-xs text-slate-400">
            {t('workbench.runtime.summary.uptime', {
              seconds: runtimeStatusQuery.data?.uptime_seconds ?? 0,
            })}
          </p>
        </div>

        <article className="rounded-xl border border-slate-800 bg-slate-900/70 p-4">
          <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
            {t('workbench.runtime.summary.pointsHealthy')}
          </p>
          <p className="mt-2 text-lg font-semibold text-slate-50">
            {runtimeStatusQuery.data?.collectors[0]?.points_healthy ?? 0}
          </p>
        </article>

        <article className="rounded-xl border border-slate-800 bg-slate-900/70 p-4">
          <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
            {t('workbench.runtime.summary.pointsStale')}
          </p>
          <p className="mt-2 text-lg font-semibold text-slate-50">
            {runtimeStatusQuery.data?.collectors[0]?.points_stale ?? 0}
          </p>
        </article>
      </div>

      <div className="grid gap-4 xl:grid-cols-[300px_minmax(0,1fr)]">
        <aside className="space-y-4">
          <section className="space-y-3 rounded-2xl border border-slate-800 bg-slate-950/40 p-4">
            <div className="space-y-1">
              <p className="text-[10px] font-bold uppercase tracking-[0.2em] text-cyan-300">
                {t('workbench.source.ruleLayer.eyebrow')}
              </p>
              <h3 className="text-sm font-semibold text-slate-100">
                {t('workbench.source.ruleLayer.title')}
              </h3>
              <p className="text-xs text-slate-400">
                {t('workbench.source.ruleLayer.description')}
              </p>
            </div>

            {rules.length === 0 ? (
              <p className="rounded-xl border border-dashed border-slate-700/60 bg-slate-900/40 px-3 py-4 text-xs text-slate-400">
                {t('workbench.source.ruleLayer.empty')}
              </p>
            ) : (
              <div className="space-y-3">
                {rules.map((rule) => {
                  const coverage = buildSourceRuleCoverage(rule, selectedDevice.protocol);
                  const isSelected = rule.id === selectedRuleId;
                  return (
                    <article
                      className={[
                        'space-y-3 rounded-xl border p-3 transition',
                        isSelected
                          ? 'border-cyan-500/40 bg-cyan-500/10'
                          : 'border-slate-800 bg-slate-900/70',
                      ].join(' ')}
                      data-testid={`source-rule-${rule.id}`}
                      key={rule.id}
                    >
                      <button
                        className="block w-full space-y-2 text-left"
                        onClick={() => handleSelectRule(rule.id)}
                        type="button"
                      >
                        <div className="flex items-center justify-between gap-3">
                          <span className="text-sm font-semibold text-slate-100">
                            {rule.id}
                          </span>
                          <span className="text-xs text-slate-400">
                            {coverage.startAddress} → {coverage.endAddress}
                          </span>
                        </div>
                        <p className="text-xs text-slate-400">
                          {t('workbench.source.ruleLayer.meta', {
                            bitWidth: coverage.bitWidth,
                            count: rule.count,
                            cells: coverage.cellCount,
                          })}
                        </p>
                      </button>

                      <div className="grid grid-cols-2 gap-2 text-xs">
                        <button
                          className="rounded-lg border border-slate-700/70 px-2 py-1 text-slate-200"
                          onClick={() => handleToggleRuleEnabled(rule.id)}
                          type="button"
                        >
                          {rule.enabled
                            ? t('workbench.source.ruleLayer.disable')
                            : t('workbench.source.ruleLayer.enable')}
                        </button>
                        <button
                          className="rounded-lg border border-slate-700/70 px-2 py-1 text-slate-200"
                          onClick={() => handleToggleRuleLocked(rule.id)}
                          type="button"
                        >
                          {rule.locked
                            ? t('workbench.source.ruleLayer.unlock')
                            : t('workbench.source.ruleLayer.lock')}
                        </button>
                        <button
                          className="rounded-lg border border-slate-700/70 px-2 py-1 text-slate-200"
                          onClick={() => handleMoveRule(rule.id, -1)}
                          type="button"
                        >
                          {t('workbench.source.ruleLayer.moveUp')}
                        </button>
                        <button
                          className="rounded-lg border border-slate-700/70 px-2 py-1 text-slate-200"
                          onClick={() => handleMoveRule(rule.id, 1)}
                          type="button"
                        >
                          {t('workbench.source.ruleLayer.moveDown')}
                        </button>
                      </div>
                    </article>
                  );
                })}
              </div>
            )}
          </section>
        </aside>

        <div className="space-y-4">
          <section className="space-y-4 rounded-2xl border border-slate-800 bg-slate-950/40 p-4">
            <div className="flex flex-col gap-3 xl:flex-row xl:items-center xl:justify-between">
              <div className="flex flex-wrap gap-2">
                {(['plan', 'live', 'link'] as const).map((mode) => (
                  <button
                    key={mode}
                    type="button"
                    onClick={() => setViewMode(mode)}
                    className={
                      viewMode === mode
                        ? 'rounded-lg bg-cyan-500 px-3 py-2 text-sm font-semibold text-slate-950'
                        : 'rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-300'
                    }
                  >
                    {t(`workbench.source.view.${mode}`)}
                  </button>
                ))}
              </div>

              <div className="flex flex-wrap items-center gap-2">
                <select
                  aria-label={t('workbench.source.toolbar.valueFormat')}
                  value={valueFormat}
                  onChange={(event) =>
                    setValueFormat(event.target.value as SourceValueFormat)
                  }
                  className="rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-200"
                >
                  <option value="decimal">{t('workbench.source.formats.decimal')}</option>
                  <option value="hex">{t('workbench.source.formats.hex')}</option>
                  <option value="binary">{t('workbench.source.formats.binary')}</option>
                  <option value="float">{t('workbench.source.formats.float')}</option>
                </select>
                <button
                  type="button"
                  onClick={handleToggleFreezeLive}
                  className="rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-300"
                >
                  {freezeLive
                    ? t('workbench.source.toolbar.unfreezeLive')
                    : t('workbench.source.toolbar.freezeLive')}
                </button>
                <button
                  type="button"
                  onClick={handleCaptureSnapshot}
                  className="rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-300"
                >
                  {t('workbench.source.toolbar.snapshotCompare')}
                </button>
                <button
                  type="button"
                  onClick={() => setShowAudit((currentValue) => !currentValue)}
                  className="rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-300"
                >
                  {showAudit
                    ? t('workbench.source.toolbar.hideAudit')
                    : t('workbench.source.toolbar.showAudit')}
                </button>
              </div>
            </div>

            <div className="grid gap-3 md:grid-cols-5">
              <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                <span>{t('workbench.source.planner.startAddress')}</span>
                <input
                  aria-label={t('workbench.source.planner.startAddress')}
                  value={startAddress}
                  onChange={(event) => setStartAddress(event.target.value)}
                  className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
                />
              </label>
              <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                <span>{t('workbench.source.planner.dataType')}</span>
                <select
                  aria-label={t('workbench.source.planner.dataType')}
                  value={dataType}
                  onChange={(event) => setDataType(event.target.value as DataType)}
                  className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
                >
                  {SOURCE_PLANNER_ALLOWED_DATA_TYPES.map((option) => (
                    <option key={option} value={option}>
                      {option}
                    </option>
                  ))}
                </select>
              </label>
              <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                <span>{t('workbench.source.planner.count')}</span>
                <input
                  aria-label={t('workbench.source.planner.count')}
                  type="number"
                  min={1}
                  value={count}
                  onChange={(event) => setCount(Number(event.target.value) || 0)}
                  className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
                />
              </label>
              <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                <span>{t('workbench.source.planner.namingPrefix')}</span>
                <input
                  aria-label={t('workbench.source.planner.namingPrefix')}
                  value={namingPrefix}
                  onChange={(event) => setNamingPrefix(event.target.value)}
                  className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
                />
              </label>
              <div className="flex items-end gap-2">
                <button
                  type="button"
                  onClick={handleApplyPlan}
                  className="w-full rounded-lg bg-cyan-500 px-3 py-2 text-sm font-medium text-slate-950"
                >
                  {t('workbench.source.planner.apply')}
                </button>
              </div>
            </div>

            <div className="grid gap-3 md:grid-cols-[minmax(0,1fr)_auto]">
              <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                <span>{t('workbench.source.toolbar.jumpToAddress')}</span>
                <input
                  aria-label={t('workbench.source.toolbar.jumpToAddress')}
                  value={jumpAddress}
                  onChange={(event) => setJumpAddress(event.target.value)}
                  className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
                />
              </label>
              <div className="flex items-end gap-2">
                <button
                  type="button"
                  onClick={handleJumpToAddress}
                  className="rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-300"
                >
                  {t('workbench.source.toolbar.jump')}
                </button>
              </div>
            </div>

            <div className="flex flex-wrap gap-2">
              <button
                type="button"
                onClick={() => void handleBatchCreate()}
                disabled={
                  plannedPointDefinitions.length === 0 ||
                  hasConflicts ||
                  createPointMutation.isPending
                }
                className="rounded-lg border border-slate-800 px-3 py-2 text-sm text-slate-300 disabled:opacity-50"
              >
                {t('workbench.source.planner.batchCreate')}
              </button>
              <button
                type="button"
                disabled
                className="rounded-lg border border-slate-800 px-3 py-2 text-sm text-slate-300 disabled:opacity-50"
              >
                {t('workbench.source.toolbar.saveTemplate')}
              </button>
              <button
                type="button"
                disabled
                className="rounded-lg border border-slate-800 px-3 py-2 text-sm text-slate-300 disabled:opacity-50"
              >
                {t('workbench.source.toolbar.loadTemplate')}
              </button>
            </div>

            {hasConflicts ? (
              <p className="text-sm text-rose-300">
                {t('workbench.source.planner.conflictHint')}
              </p>
            ) : null}

            {batchCreateSummary ? (
              <p className="text-sm text-slate-300">
                {t('workbench.source.planner.batchSummary', batchCreateSummary)}
              </p>
            ) : null}
          </section>

          <section
            className="space-y-3 rounded-2xl border border-slate-800 bg-slate-950/40 p-4"
            data-testid="source-coverage-overview"
          >
            <div className="space-y-1">
              <p className="text-[10px] font-bold uppercase tracking-[0.2em] text-cyan-300">
                {t('workbench.source.coverage.eyebrow')}
              </p>
              <h3 className="text-sm font-semibold text-slate-100">
                {t('workbench.source.coverage.title')}
              </h3>
            </div>
            {coverageSegments.length > 0 ? (
              <div className="space-y-3">
                <div className="flex gap-1">
                  {coverageSegments.map((segment) => (
                    <div
                      key={segment.id}
                      className={`h-3 rounded-full ${getCoverageSegmentClass(segment.status)}`}
                      style={{ flex: Math.max(segment.cellCount, 1) }}
                      title={`${segment.startAddress} → ${segment.endAddress}`}
                    />
                  ))}
                </div>
                <div className="grid gap-2 md:grid-cols-2 xl:grid-cols-4">
                  {coverageSegments.map((segment) => (
                    <button
                      key={`${segment.id}-jump`}
                      className="rounded-lg border border-slate-800 bg-slate-900/70 px-3 py-2 text-left text-xs text-slate-300"
                      onClick={() => handleSelectAddress(segment.startAddress)}
                      type="button"
                    >
                      <p className="font-medium text-slate-100">
                        {segment.startAddress} → {segment.endAddress}
                      </p>
                      <p className="mt-1">
                        {t(`workbench.source.coverage.status.${segment.status}`)} ·{' '}
                        {t('workbench.source.coverage.cells', {
                          count: segment.cellCount,
                        })}
                      </p>
                    </button>
                  ))}
                </div>
              </div>
            ) : (
              <p className="text-xs text-slate-400">
                {t('workbench.source.coverage.empty')}
              </p>
            )}
          </section>

          <AddressCanvas
            items={items}
            onSelectAddress={handleSelectAddress}
            selectedAddress={selectedAddress}
            valueFormat={valueFormat}
            viewMode={viewMode}
          />

          {showAudit ? (
            <AddressLedger
              items={items}
              valueFormat={valueFormat}
              viewMode={viewMode}
            />
          ) : null}
        </div>
      </div>
    </section>
  );
}
