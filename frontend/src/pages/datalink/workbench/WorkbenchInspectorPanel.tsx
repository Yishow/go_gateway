import { useQueryClient } from '@tanstack/react-query';
import { useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { deviceKeys } from '../../../hooks/datalink/keys';
import {
  useDevicesQuery,
  useTestConnectionMutation,
} from '../../../hooks/datalink/useDevices';
import { useMappingsQuery } from '../../../hooks/datalink/useMappings';
import { usePointsQuery } from '../../../hooks/datalink/usePoints';
import { useTagsQuery } from '../../../hooks/datalink/useTags';
import { dbTargetAPI, modbusShareAPI } from '../../../services/datalink';
import {
  buildDeviceCapabilitySummary,
  buildDeviceConnectionSummary,
  getDeviceTestTimestampLabel,
  getWorkbenchDeviceStatusLabelKey,
  getWorkbenchProtocolLabelKey,
  parseDeviceConnectionConfig,
} from './workbenchDeviceFormModel';
import {
  buildAddressCanvasItems,
  buildSourceRuleCoverage,
  getDataTypeBitWidth,
} from './sourceCanvasModel';
import { computeOutputReadiness } from './workbenchOutputTypes';
import { useWorkbench } from './WorkbenchProvider';
import { WORKBENCH_STEP_META, type InspectorSelection } from './workbenchTypes';

function joinClasses(...classNames: Array<string | false | null | undefined>) {
  return classNames.filter(Boolean).join(' ');
}

function getDeviceStatusClasses(status: 'draft' | 'active' | 'disabled') {
  switch (status) {
    case 'active':
      return 'border-emerald-500/40 bg-emerald-500/10 text-emerald-300';
    case 'disabled':
      return 'border-rose-500/40 bg-rose-500/10 text-rose-300';
    case 'draft':
      return 'border-amber-500/40 bg-amber-500/10 text-amber-300';
  }
}

function getSelectionLabelKey(selection: InspectorSelection): string {
  switch (selection.kind) {
    case 'none':
      return 'workbench.inspector.noSelection';
    case 'device':
      return 'workbench.inspector.selectionKind.device';
    case 'rule':
      return 'workbench.inspector.selectionKind.rule';
    case 'span':
      return 'workbench.inspector.selectionKind.span';
    case 'tag':
      return 'workbench.inspector.selectionKind.tag';
    case 'outputCandidate':
      return 'workbench.inspector.selectionKind.outputCandidate';
  }
}

function getSelectionId(selection: InspectorSelection): string | null {
  switch (selection.kind) {
    case 'none':
      return null;
    case 'device':
      return selection.deviceId;
    case 'rule':
      return selection.ruleId;
    case 'span':
      return selection.spanAddress;
    case 'tag':
      return selection.tagId;
    case 'outputCandidate':
      return selection.tagId;
  }
}

function formatInspectorValue(value: unknown): string {
  if (value === null || value === undefined || value === '') {
    return '—';
  }

  return String(value);
}

function DeviceInspectorContent() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const testConnectionMutation = useTestConnectionMutation();
  const {
    openCloneDevicePanel,
    openCreateDevicePanel,
    openEditDevicePanel,
    recentDeviceTests,
    recordDeviceTest,
    selectedDeviceId,
  } = useWorkbench();
  const { data: devices = [] } = useDevicesQuery();
  const selectedDevice =
    devices.find((device) => device.id === selectedDeviceId) ?? null;

  if (!selectedDevice) {
    return (
      <div className="flex flex-1 flex-col justify-center gap-4 rounded-xl border border-dashed border-slate-700/60 bg-slate-950/30 p-6">
      <div className="space-y-2">
        <h3 className="text-sm font-semibold text-slate-200">
          {t('workbench.device.inspector.emptyTitle')}
        </h3>
        <p className="text-xs text-slate-400">
          {t('workbench.device.inspector.emptyDescription')}
        </p>
      </div>
        {devices.length > 0 ? (
          <button
            className="rounded-xl bg-cyan-500 px-3 py-2 text-sm font-semibold text-slate-950 transition hover:bg-cyan-400"
            onClick={openCreateDevicePanel}
            type="button"
          >
            {t('workbench.device.actions.create')}
          </button>
        ) : null}
      </div>
    );
  }

  const connectionConfig = parseDeviceConnectionConfig(selectedDevice.connection_config);
  const capabilitySummary = buildDeviceCapabilitySummary(
    selectedDevice.protocol,
    connectionConfig,
    t,
  );
  const connectionSummary = buildDeviceConnectionSummary(
    selectedDevice.protocol,
    connectionConfig,
  );
  const storedHistory = recentDeviceTests[selectedDevice.id] ?? [];
  const fallbackHistory =
    storedHistory.length === 0
      && (selectedDevice.last_test_success !== null || selectedDevice.last_test_at)
      ? [
          {
            id: `${selectedDevice.id}-latest`,
            testedAt: getDeviceTestTimestampLabel(selectedDevice.last_test_at, t),
            success: selectedDevice.last_test_success === true,
            message:
              selectedDevice.last_test_success === true
                ? t('workbench.device.card.testPassed')
                : selectedDevice.last_test_success === false
                  ? selectedDevice.last_test_error || t('workbench.device.card.testFailed')
                  : t('workbench.device.inspector.noTestYet'),
            latencyMs: null,
          },
        ]
      : [];
  const recentHistory = [...storedHistory, ...fallbackHistory].slice(0, 3);

  const handleTestConnection = async () => {
    if (testConnectionMutation.isPending) {
      return;
    }

    try {
      const result = await testConnectionMutation.mutateAsync(selectedDevice.id);
      await queryClient.invalidateQueries({ queryKey: deviceKeys.lists() });
      recordDeviceTest(selectedDevice.id, {
        testedAt: new Date().toISOString(),
        success: result.success,
        message: result.success
          ? t('workbench.device.messages.testSuccess')
          : result.error || t('workbench.device.messages.testFailed'),
        latencyMs: result.latency_ms ?? null,
      });
    } catch (error) {
      recordDeviceTest(selectedDevice.id, {
        testedAt: new Date().toISOString(),
        success: false,
        message:
          error instanceof Error
            ? error.message
            : t('workbench.device.messages.testFailed'),
        latencyMs: null,
      });
    }
  };

  return (
    <div className="flex h-full flex-col gap-4">
      <div className="space-y-3 rounded-xl border border-slate-700/40 bg-slate-950/40 p-4">
        <div className="space-y-2">
          <h3 className="text-lg font-semibold text-slate-50">{selectedDevice.name}</h3>
          <p className="text-xs text-slate-400">
            {selectedDevice.description || t('workbench.device.card.noDescription')}
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
          <span className="rounded-full border border-slate-700 bg-slate-950/70 px-2 py-0.5 text-[11px] text-slate-300">
            {t(getWorkbenchProtocolLabelKey(selectedDevice.protocol))}
          </span>
          <span
            className={joinClasses(
              'rounded-full border px-2 py-0.5 text-[11px] font-medium',
              getDeviceStatusClasses(selectedDevice.status),
            )}
          >
            {t(getWorkbenchDeviceStatusLabelKey(selectedDevice.status))}
          </span>
        </div>
      </div>

      <section className="space-y-3 rounded-xl border border-slate-700/40 bg-slate-950/40 p-4">
        <p className="text-[10px] font-bold uppercase tracking-[0.18em] text-slate-500">
          {t('workbench.device.inspector.capabilitySummary')}
        </p>
        <dl className="space-y-3 text-sm text-slate-200">
          {capabilitySummary.map((item) => (
            <div className="space-y-1" key={item.id}>
              <dt className="text-xs uppercase tracking-[0.14em] text-slate-500">
                {t(item.labelKey)}
              </dt>
              <dd className="font-medium text-slate-100">{item.value}</dd>
            </div>
          ))}
        </dl>
      </section>

      <section className="space-y-3 rounded-xl border border-slate-700/40 bg-slate-950/40 p-4">
        <p className="text-[10px] font-bold uppercase tracking-[0.18em] text-slate-500">
          {t('workbench.device.inspector.connectionSummary')}
        </p>
        <dl className="space-y-3 text-sm text-slate-200">
          {connectionSummary.map((item) => (
            <div className="flex items-center justify-between gap-4" key={item.labelKey}>
              <dt className="text-slate-400">{t(item.labelKey)}</dt>
              <dd className="text-right font-medium text-slate-100">{item.value}</dd>
            </div>
          ))}
        </dl>
      </section>

      <section className="space-y-3 rounded-xl border border-slate-700/40 bg-slate-950/40 p-4">
        <p className="text-[10px] font-bold uppercase tracking-[0.18em] text-slate-500">
          {t('workbench.device.inspector.recentTests')}
        </p>
        {recentHistory.length > 0 ? (
          <ol className="space-y-3">
            {recentHistory.map((entry) => (
              <li
                className="rounded-xl border border-slate-800 bg-slate-950/70 p-3"
                key={entry.id}
              >
                <p
                  className={joinClasses(
                    'text-sm font-medium',
                    entry.success ? 'text-emerald-300' : 'text-rose-300',
                  )}
                >
                  {entry.message}
                </p>
                <p className="mt-1 text-xs text-slate-500">
                  {getDeviceTestTimestampLabel(entry.testedAt, t)}
                </p>
              </li>
            ))}
          </ol>
        ) : (
          <p className="text-xs text-slate-400">{t('workbench.device.inspector.noTestYet')}</p>
        )}
      </section>

      <div className="mt-auto grid gap-2">
        <button
          className="rounded-xl border border-slate-700 bg-slate-950/60 px-3 py-2 text-sm font-medium text-slate-200 transition hover:border-slate-500 hover:text-slate-50"
          onClick={() => openEditDevicePanel(selectedDevice.id)}
          type="button"
        >
          {t('workbench.device.actions.edit')}
        </button>
        <button
          className="rounded-xl border border-slate-700 bg-slate-950/60 px-3 py-2 text-sm font-medium text-slate-200 transition hover:border-slate-500 hover:text-slate-50 disabled:cursor-not-allowed disabled:opacity-50"
          disabled={testConnectionMutation.isPending}
          onClick={() => void handleTestConnection()}
          type="button"
        >
          {testConnectionMutation.isPending
            ? t('workbench.device.actions.testing')
            : t('workbench.device.actions.testConnection')}
        </button>
        <button
          className="rounded-xl border border-slate-700 bg-slate-950/60 px-3 py-2 text-sm font-medium text-slate-200 transition hover:border-slate-500 hover:text-slate-50"
          onClick={() => openCloneDevicePanel(selectedDevice.id)}
          type="button"
        >
          {t('workbench.device.actions.clone')}
        </button>
      </div>
    </div>
  );
}

function SourceInspectorContent() {
  const { t } = useTranslation();
  const {
    inspectorSelection,
    selectedDeviceId,
    setSourcePlanningState,
    sourcePlanningState,
  } = useWorkbench();
  const { data: devices = [] } = useDevicesQuery();
  const { data: points = [] } = usePointsQuery(
    selectedDeviceId ? { device_id: selectedDeviceId } : undefined,
  );
  const { data: mappings = [] } = useMappingsQuery();
  const { data: tags = [] } = useTagsQuery();
  const selectedDevice =
    devices.find((device) => device.id === selectedDeviceId) ?? null;

  const items = useMemo(() => {
    if (!selectedDevice) {
      return [];
    }

    return buildAddressCanvasItems({
      points,
      rules: sourcePlanningState.rules,
      mappings,
      tags,
      protocol: selectedDevice.protocol,
    });
  }, [mappings, points, selectedDevice, sourcePlanningState.rules, tags]);

  if (!selectedDevice) {
    return null;
  }

  if (inspectorSelection.kind === 'rule') {
    const rule =
      sourcePlanningState.rules.find((item) => item.id === inspectorSelection.ruleId)
      ?? null;
    if (!rule) {
      return null;
    }

    const coverage = buildSourceRuleCoverage(rule, selectedDevice.protocol);

    return (
      <div
        data-testid="source-rule-inspector"
        className="space-y-4 rounded-xl border border-slate-700/40 bg-slate-950/40 p-4"
      >
        <div className="space-y-1">
          <p className="text-[10px] font-bold uppercase tracking-[0.18em] text-slate-500">
            {t('workbench.source.inspector.rule.heading')}
          </p>
          <h3 className="text-sm font-semibold text-slate-100">{rule.namingPrefix}</h3>
        </div>

        <dl className="space-y-3 text-sm">
          <div className="space-y-1">
            <dt className="text-xs uppercase tracking-[0.14em] text-slate-500">
              {t('workbench.source.inspector.rule.coverage')}
            </dt>
            <dd className="font-medium text-slate-100" data-testid="source-rule-coverage">
              {coverage.startAddress} → {coverage.endAddress}
            </dd>
          </div>
          <div className="grid gap-3 sm:grid-cols-2">
            <div className="space-y-1">
              <dt className="text-xs uppercase tracking-[0.14em] text-slate-500">
                {t('workbench.source.inspector.rule.count')}
              </dt>
              <dd className="font-medium text-slate-100">
                {rule.count}
              </dd>
            </div>
            <div className="space-y-1">
              <dt className="text-xs uppercase tracking-[0.14em] text-slate-500">
                {t('workbench.source.inspector.rule.bitWidth')}
              </dt>
              <dd className="font-medium text-slate-100">
                {coverage.bitWidth}
              </dd>
            </div>
          </div>
          <div className="space-y-1">
            <dt className="text-xs uppercase tracking-[0.14em] text-slate-500">
              {t('workbench.source.inspector.rule.origin')}
            </dt>
            <dd className="font-medium text-slate-100">
              {rule.origin === 'template'
                ? `${t('workbench.source.inspector.rule.originTemplate')} · ${rule.templateName ?? '—'}`
                : t('workbench.source.inspector.rule.originManual')}
            </dd>
          </div>
        </dl>

        <div className="flex flex-wrap gap-2">
          <button
            type="button"
            onClick={() =>
              setSourcePlanningState((currentState) => ({
                ...currentState,
                rules: currentState.rules.map((item) =>
                  item.id === rule.id ? { ...item, enabled: !item.enabled } : item,
                ),
              }))
            }
            className="rounded-lg border border-slate-700 px-3 py-1.5 text-xs text-slate-200"
          >
            {rule.enabled
              ? t('workbench.source.ruleLayer.disable')
              : t('workbench.source.ruleLayer.enable')}
          </button>
          <button
            type="button"
            onClick={() =>
              setSourcePlanningState((currentState) => ({
                ...currentState,
                rules: currentState.rules.map((item) =>
                  item.id === rule.id ? { ...item, locked: !item.locked } : item,
                ),
              }))
            }
            className="rounded-lg border border-slate-700 px-3 py-1.5 text-xs text-slate-200"
          >
            {rule.locked
              ? t('workbench.source.ruleLayer.unlock')
              : t('workbench.source.ruleLayer.lock')}
          </button>
        </div>
      </div>
    );
  }

  if (inspectorSelection.kind === 'span') {
    const item =
      items.find((candidate) => candidate.address === inspectorSelection.spanAddress) ?? null;
    const rule =
      sourcePlanningState.rules.find((candidate) => candidate.id === inspectorSelection.ruleId)
      ?? null;
    const bitWidth = item?.point
      ? getDataTypeBitWidth(item.point.data_type)
      : rule
        ? getDataTypeBitWidth(rule.dataType)
        : 0;

    return (
      <div
        data-testid="source-span-inspector"
        className="space-y-4 rounded-xl border border-slate-700/40 bg-slate-950/40 p-4"
      >
        <div className="space-y-1">
          <p className="text-[10px] font-bold uppercase tracking-[0.18em] text-slate-500">
            {t('workbench.source.inspector.span.heading')}
          </p>
          <h3 className="text-sm font-semibold text-slate-100">
            {inspectorSelection.spanAddress}
          </h3>
        </div>

        <dl className="space-y-3 text-sm">
          <div className="space-y-1">
            <dt className="text-xs uppercase tracking-[0.14em] text-slate-500">
              {t('workbench.source.inspector.span.pointName')}
            </dt>
            <dd className="font-medium text-slate-100">
              {item?.point?.name ?? '—'}
            </dd>
          </div>
          <div className="grid gap-3 sm:grid-cols-2">
            <div className="space-y-1">
              <dt className="text-xs uppercase tracking-[0.14em] text-slate-500">
                {t('workbench.source.inspector.span.bitWidth')}
              </dt>
              <dd className="font-medium text-slate-100">{bitWidth || '—'}</dd>
            </div>
            <div className="space-y-1">
              <dt className="text-xs uppercase tracking-[0.14em] text-slate-500">
                {t('workbench.source.inspector.span.rawValue')}
              </dt>
              <dd className="font-medium text-slate-100">
                {formatInspectorValue(item?.liveValue ?? item?.point?.last_value)}
              </dd>
            </div>
          </div>
          <div className="space-y-1">
            <dt className="text-xs uppercase tracking-[0.14em] text-slate-500">
              {t('workbench.source.inspector.span.linkState')}
            </dt>
            <dd
              className="font-medium text-slate-100"
              data-testid="source-span-link-state"
            >
              {item?.linkLabelKey ? t(item.linkLabelKey) : '—'}
            </dd>
          </div>
        </dl>
      </div>
    );
  }

  return null;
}

function getTagBindingClasses(status: 'bound' | 'unbound' | 'partial') {
  switch (status) {
    case 'bound':
      return 'border-emerald-500/30 bg-emerald-500/10 text-emerald-200';
    case 'partial':
      return 'border-amber-500/30 bg-amber-500/10 text-amber-200';
    case 'unbound':
      return 'border-slate-700 bg-slate-900/80 text-slate-300';
  }
}

function TagInspectorContent() {
  const { t } = useTranslation();
  const { inspectorSelection, selectedDeviceId } = useWorkbench();
  const { data: points = [] } = usePointsQuery(
    selectedDeviceId ? { device_id: selectedDeviceId } : undefined,
  );
  const { data: tags = [] } = useTagsQuery();
  const { data: mappings = [] } = useMappingsQuery();

  if (inspectorSelection.kind !== 'tag') {
    return null;
  }

  const point = inspectorSelection.pointId
    ? points.find((candidate) => candidate.id === inspectorSelection.pointId) ?? null
    : null;
  const mapping = inspectorSelection.pointId
    ? mappings.find((candidate) => candidate.point_id === inspectorSelection.pointId) ?? null
    : mappings.find((candidate) => candidate.tag_id === inspectorSelection.tagId) ?? null;
  const boundTag = mapping
    ? tags.find((candidate) => candidate.id === mapping.tag_id) ?? null
    : tags.find((candidate) => candidate.id === inspectorSelection.tagId) ?? null;
  const bindingStatus =
    inspectorSelection.bindingStatus ?? (mapping ? 'bound' : 'unbound');
  const displayTagKey =
    boundTag?.key
    ?? inspectorSelection.existingTagLabel
    ?? inspectorSelection.tagId;

  return (
    <div
      data-testid="tag-inspector-panel"
      className="space-y-4 rounded-xl border border-slate-700/40 bg-slate-950/40 p-4"
    >
      <div className="space-y-1">
        <p className="text-[10px] font-bold uppercase tracking-[0.18em] text-slate-500">
          {t('workbench.tag.inspector.heading')}
        </p>
        <h3
          className="text-sm font-semibold text-slate-100"
          data-testid="tag-inspector-tag-key"
        >
          {displayTagKey}
        </h3>
      </div>

      <span
        className={`inline-flex rounded-full border px-2 py-1 text-xs font-medium ${getTagBindingClasses(
          bindingStatus,
        )}`}
        data-testid="tag-inspector-binding-status"
      >
        {t(`workbench.tag.board.status.${bindingStatus}`)}
      </span>

      <dl className="space-y-3 text-sm">
        <div className="space-y-1">
          <dt className="text-xs uppercase tracking-[0.14em] text-slate-500">
            {t('workbench.tag.inspector.pointName')}
          </dt>
          <dd className="font-medium text-slate-100">
            {inspectorSelection.pointName ?? point?.name ?? '—'}
          </dd>
        </div>
        <div className="space-y-1">
          <dt className="text-xs uppercase tracking-[0.14em] text-slate-500">
            {t('workbench.tag.inspector.sourceAddress')}
          </dt>
          <dd
            className="font-medium text-slate-100"
            data-testid="tag-inspector-point-address"
          >
            {inspectorSelection.pointAddress ?? point?.address ?? '—'}
          </dd>
        </div>
        <div className="grid gap-3 sm:grid-cols-2">
          <div className="space-y-1">
            <dt className="text-xs uppercase tracking-[0.14em] text-slate-500">
              {t('workbench.tag.board.rawValue')}
            </dt>
            <dd className="font-medium text-slate-100">
              {formatInspectorValue(inspectorSelection.rawValue ?? point?.last_value)}
            </dd>
          </div>
          <div className="space-y-1">
            <dt className="text-xs uppercase tracking-[0.14em] text-slate-500">
              {t('workbench.tag.board.transformedValue')}
            </dt>
            <dd className="font-medium text-slate-100">
              {formatInspectorValue(inspectorSelection.transformedValue)}
            </dd>
          </div>
        </div>
        <div className="grid gap-3 sm:grid-cols-2">
          <div className="space-y-1">
            <dt className="text-xs uppercase tracking-[0.14em] text-slate-500">
              {t('workbench.tag.inspector.span')}
            </dt>
            <dd className="font-medium text-slate-100">
              {inspectorSelection.cellSpan ?? '—'}
            </dd>
          </div>
          <div className="space-y-1">
            <dt className="text-xs uppercase tracking-[0.14em] text-slate-500">
              {t('workbench.tag.inspector.bitWidth')}
            </dt>
            <dd className="font-medium text-slate-100">
              {inspectorSelection.bitWidth ?? '—'}
            </dd>
          </div>
        </div>
      </dl>

      {inspectorSelection.conflictReason ? (
        <div className="rounded-lg border border-rose-500/20 bg-rose-500/10 px-3 py-2 text-xs text-rose-200">
          {t(`workbench.tag.inspector.conflict.${inspectorSelection.conflictReason}`)}
        </div>
      ) : null}

      {inspectorSelection.alreadyLinked ? (
        <div className="rounded-lg border border-amber-500/20 bg-amber-500/10 px-3 py-2 text-xs text-amber-200">
          {t('workbench.tag.inspector.alreadyLinked')}
        </div>
      ) : null}
    </div>
  );
}

function OutputInspectorContent() {
  const { t } = useTranslation();
  const { inspectorSelection, selectedDeviceId } = useWorkbench();
  const { data: tags = [] } = useTagsQuery();
  const { data: mappings = [] } = useMappingsQuery();
  const { data: points = [] } = usePointsQuery(
    selectedDeviceId ? { device_id: selectedDeviceId } : undefined,
  );
  const [modbusMappings, setModbusMappings] = useState<Array<{ tag_id: string; register: number }>>(
    [],
  );
  const [databaseMappings, setDatabaseMappings] = useState<
    Array<{ tag_id: string; table_schema: string; table_name: string; column_name: string }>
  >([]);
  const [loadError, setLoadError] = useState<string | null>(null);

  useEffect(() => {
    if (inspectorSelection.kind !== 'outputCandidate') {
      return;
    }

    let cancelled = false;
    void Promise.all([modbusShareAPI.listMappings(), dbTargetAPI.listMappings()])
      .then(([nextModbusMappings, nextDatabaseMappings]) => {
        if (cancelled) {
          return;
        }
        setModbusMappings(nextModbusMappings);
        setDatabaseMappings(nextDatabaseMappings);
        setLoadError(null);
      })
      .catch((error: unknown) => {
        if (cancelled) {
          return;
        }
        setModbusMappings([]);
        setDatabaseMappings([]);
        setLoadError(
          error instanceof Error
            ? error.message
            : t('workbench.output.inspector.loadFailed'),
        );
      });

    return () => {
      cancelled = true;
    };
  }, [inspectorSelection, t]);

  if (inspectorSelection.kind !== 'outputCandidate') {
    return null;
  }

  const tag = tags.find((item) => item.id === inspectorSelection.tagId) ?? null;
  const mapping = mappings.find((item) => item.tag_id === inspectorSelection.tagId) ?? null;
  const point = mapping ? points.find((item) => item.id === mapping.point_id) ?? null : null;
  const modbusMapping =
    modbusMappings.find((item) => item.tag_id === inspectorSelection.tagId) ?? null;
  const databaseMapping =
    databaseMappings.find((item) => item.tag_id === inspectorSelection.tagId) ?? null;
  const readiness = computeOutputReadiness(Boolean(modbusMapping), Boolean(databaseMapping));

  return (
    <div
      data-testid="inspector-trace-panel"
      className="space-y-4 rounded-xl border border-slate-700/40 bg-slate-950/40 p-4"
    >
      <div className="space-y-1">
        <p className="text-[10px] font-bold uppercase tracking-[0.18em] text-slate-500">
          {t('workbench.output.inspector.heading')}
        </p>
        <p
          className={`inline-flex rounded-full border px-2 py-1 text-xs font-medium ${
            readiness === 'ready'
              ? 'border-emerald-500/30 bg-emerald-500/10 text-emerald-200'
              : readiness === 'partial'
                ? 'border-amber-500/30 bg-amber-500/10 text-amber-200'
                : 'border-slate-700 bg-slate-900/80 text-slate-300'
          }`}
          data-testid="trace-readiness"
        >
          {t(`workbench.output.inspector.readiness.${readiness}`)}
        </p>
      </div>

      {loadError ? (
        <p
          role="status"
          aria-live="polite"
          className="rounded-lg border border-rose-500/20 bg-rose-500/10 px-3 py-2 text-xs text-rose-200"
        >
          {loadError}
        </p>
      ) : null}

      <dl className="space-y-3 text-sm">
        <div className="space-y-1">
          <dt className="text-xs uppercase tracking-[0.14em] text-slate-500">
            {t('workbench.output.inspector.sourceAddress')}
          </dt>
          <dd className="font-medium text-slate-100" data-testid="trace-source-address">
            {point?.address ?? '—'}
          </dd>
        </div>
        <div className="space-y-1">
          <dt className="text-xs uppercase tracking-[0.14em] text-slate-500">
            {t('workbench.output.inspector.tagKey')}
          </dt>
          <dd className="font-medium text-slate-100" data-testid="trace-tag-key">
            {tag?.key ?? inspectorSelection.tagId}
          </dd>
        </div>
        <div className="space-y-1">
          <dt className="text-xs uppercase tracking-[0.14em] text-slate-500">
            {t('workbench.output.inspector.outputModbus')}
          </dt>
          <dd className="font-medium text-slate-100" data-testid="trace-output-modbus">
            {modbusMapping ? `HR${modbusMapping.register}` : t('workbench.output.selection.unmapped')}
          </dd>
        </div>
        <div className="space-y-1">
          <dt className="text-xs uppercase tracking-[0.14em] text-slate-500">
            {t('workbench.output.inspector.outputDatabase')}
          </dt>
          <dd className="font-medium text-slate-100" data-testid="trace-output-database">
            {databaseMapping
              ? `${databaseMapping.table_schema}.${databaseMapping.table_name}.${databaseMapping.column_name}`
              : t('workbench.output.selection.databaseUnmapped')}
          </dd>
        </div>
      </dl>
    </div>
  );
}

export function WorkbenchInspectorPanel() {
  const { t } = useTranslation();
  const { activeStep, inspectorSelection } = useWorkbench();
  const stepMeta = WORKBENCH_STEP_META[activeStep];
  const hasSelection = inspectorSelection.kind !== 'none';
  const selectionId = getSelectionId(inspectorSelection);

  return (
    <aside
      aria-label={t('workbench.inspector.ariaLabel')}
      className="flex flex-col gap-4 rounded-2xl border border-slate-800 bg-slate-900/80 p-4"
      data-testid="workbench-inspector-panel"
    >
      <div className="space-y-1">
        <p className="text-[10px] font-bold uppercase tracking-[0.22em] text-slate-500">
          {t('workbench.inspector.heading')}
        </p>
        <h3 className="text-sm font-semibold text-slate-200">
          {t(stepMeta.labelKey)}
        </h3>
      </div>

      {activeStep === 'device' ? (
        <DeviceInspectorContent />
      ) : activeStep === 'source' &&
        (inspectorSelection.kind === 'rule' || inspectorSelection.kind === 'span') ? (
        <SourceInspectorContent />
      ) : activeStep === 'tag' && inspectorSelection.kind === 'tag' ? (
        <TagInspectorContent />
      ) : activeStep === 'output' && inspectorSelection.kind === 'outputCandidate' ? (
        <OutputInspectorContent />
      ) : hasSelection ? (
        <div
          className="space-y-2 rounded-xl border border-slate-700/40 bg-slate-950/40 p-4"
          data-testid="inspector-selection-context"
        >
          <p className="text-[10px] font-bold uppercase tracking-[0.18em] text-slate-500">
            {t(getSelectionLabelKey(inspectorSelection))}
          </p>
          {selectionId ? (
            <p
              className="truncate text-xs font-medium text-slate-200"
              data-testid="inspector-selection-id"
            >
              {selectionId}
            </p>
          ) : null}
        </div>
      ) : (
        <div className="flex flex-1 items-center justify-center rounded-xl border border-dashed border-slate-700/60 bg-slate-950/30 p-6">
          <p className="text-center text-xs text-slate-500" data-testid="inspector-placeholder">
            {t('workbench.inspector.placeholder')}
          </p>
        </div>
      )}
    </aside>
  );
}
