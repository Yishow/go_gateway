import { useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import {
  applyTemplateToPlanner,
  createTemplateFromPlanner,
  normalizeNamingPrefix,
  SOURCE_PLANNER_ALLOWED_DATA_TYPES,
  upsertTemplateRecord,
} from '../../../features/datalink/sourcePlannerContract';
import {
  isTemplateStale,
  loadSourceTemplates,
  saveSourceTemplates,
  upgradeTemplates,
  type SourceTemplateCapabilitySnapshot,
  type SourceTemplateRecord,
} from '../../../features/datalink/sourceTemplateStorage';
import { useDevicesQuery } from '../../../hooks/datalink/useDevices';
import { useMappingsQuery } from '../../../hooks/datalink/useMappings';
import {
  useCreatePointMutation,
  usePointsQuery,
} from '../../../hooks/datalink/usePoints';
import { useRuntimeStream } from '../../../hooks/datalink/useRuntimeStream';
import { useTagsQuery } from '../../../hooks/datalink/useTags';
import { runtimeAPI } from '../../../services/datalink';
import type { DataType, Device, ProtocolType } from '../../../types/datalink';
import { addressParser } from '../../../utils/addressParser';
import { AddressCanvas } from './AddressCanvas';
import { AddressLedger } from './AddressLedger';
import { useWorkbench } from './WorkbenchProvider';
import { parseDeviceConnectionConfig } from './workbenchDeviceFormModel';
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

function sortTemplates(templates: ReadonlyArray<SourceTemplateRecord>) {
  return [...templates].sort((left, right) =>
    right.lastUsedAt.localeCompare(left.lastUsedAt),
  );
}

function buildTemplateCapabilitySnapshot(
  device: Device | null,
): SourceTemplateCapabilitySnapshot | null {
  if (!device) {
    return null;
  }

  const connectionConfig = parseDeviceConnectionConfig(device.connection_config);
  const addressBase = (() => {
    switch (device.protocol) {
      case 'modbus_tcp':
      case 'modbus_udp':
      case 'modbus_rtu':
        return 'modbus-register';
      case 'mqtt':
        return 'topic-based';
      case 'fatek_fbs':
      case 'mc_3e':
        return 'protocol-native';
    }
  })();

  const wordOrder = (() => {
    switch (device.protocol) {
      case 'mc_3e':
        return typeof connectionConfig.data_format === 'string'
          ? connectionConfig.data_format
          : 'protocol-default';
      case 'mqtt':
        return 'not-applicable';
      case 'modbus_tcp':
      case 'modbus_udp':
      case 'modbus_rtu':
      case 'fatek_fbs':
        return 'protocol-default';
    }
  })();

  return {
    protocol: device.protocol,
    addressBase,
    wordOrder,
  };
}

function buildTemplateWarning(
  template: SourceTemplateRecord,
  currentCapability: SourceTemplateCapabilitySnapshot | null,
  t: (key: string, options?: Record<string, unknown>) => string,
) {
  if (!template.capabilitySnapshot || !currentCapability) {
    return null;
  }

  const mismatches: string[] = [];
  if (template.capabilitySnapshot.addressBase !== currentCapability.addressBase) {
    mismatches.push(t('workbench.source.templates.warning.addressBase'));
  }
  if (template.capabilitySnapshot.wordOrder !== currentCapability.wordOrder) {
    mismatches.push(t('workbench.source.templates.warning.wordOrder'));
  }

  if (mismatches.length === 0) {
    return null;
  }

  return t('workbench.source.templates.warning.message', {
    name: template.name,
    mismatches: mismatches.join(', '),
  });
}

type PointDefinition = {
  address: string;
  dataType: DataType;
  name: string;
};

function buildPointDefinition(rule: SourceRule, address: string): PointDefinition {
  const prefix = normalizeNamingPrefix(rule.namingPrefix);

  return {
    address,
    dataType: rule.dataType,
    name: `${prefix}_${address}`,
  };
}

function buildRulePointDefinitions(input: {
  rules: ReadonlyArray<SourceRule>;
  protocol: ProtocolType;
  locked?: boolean;
}) {
  const definitions = new Map<string, PointDefinition>();

  for (const rule of input.rules.filter(
    (candidate) =>
      candidate.enabled &&
      (input.locked === undefined || candidate.locked === input.locked),
  )) {
    const plannedAddresses = buildPlannedPointAddresses({
      startAddress: rule.startAddress,
      count: rule.count,
      dataType: rule.dataType,
      protocol: input.protocol,
    });

    plannedAddresses.forEach((address) => {
      if (!definitions.has(address)) {
        definitions.set(address, buildPointDefinition(rule, address));
      }
    });
  }

  return [...definitions.values()];
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
  const [showUtilityTools, setShowUtilityTools] = useState(false);
  const [freezeLive, setFreezeLive] = useState(false);
  const [frozenLiveValues, setFrozenLiveValues] = useState<
    Readonly<Record<string, { raw_value?: unknown; timestamp?: string }>> | null
  >(null);
  const [startAddress, setStartAddress] = useState('40001');
  const [dataType, setDataType] = useState<DataType>('int16');
  const [count, setCount] = useState(4);
  const [namingPrefix, setNamingPrefix] = useState('SRC');
  const [jumpAddress, setJumpAddress] = useState('');
  const [editingRuleId, setEditingRuleId] = useState<string | null>(null);
  const [editDraft, setEditDraft] = useState<{
    startAddress: string;
    count: number;
    dataType: DataType;
  } | null>(null);
  const [batchCreateSummary, setBatchCreateSummary] = useState<{
    successCount: number;
    failureCount: number;
  } | null>(null);
  const [templates, setTemplates] = useState<SourceTemplateRecord[]>([]);
  const [templateName, setTemplateName] = useState('');
  const [isSaveTemplateOpen, setIsSaveTemplateOpen] = useState(false);
  const [isLoadTemplateOpen, setIsLoadTemplateOpen] = useState(false);
  const [templateNotice, setTemplateNotice] = useState<string | null>(null);
  const [templateWarning, setTemplateWarning] = useState<string | null>(null);
  const [appliedTemplate, setAppliedTemplate] = useState<{
    id: string;
    name: string;
  } | null>(null);
  const rules = sourcePlanningState.rules;
  const selectedRuleId = sourcePlanningState.selectedRuleId;
  const selectedAddress = sourcePlanningState.selectedAddress;
  const currentCapability = useMemo(
    () => buildTemplateCapabilitySnapshot(selectedDevice),
    [selectedDevice],
  );

  useEffect(() => {
    const restored = loadSourceTemplates();
    if (restored.length === 0) {
      setTemplates([]);
      return;
    }

    const nextTemplates = restored.some(isTemplateStale)
      ? sortTemplates(upgradeTemplates(restored))
      : sortTemplates(restored);

    if (restored.some(isTemplateStale)) {
      saveSourceTemplates(nextTemplates);
    }

    setTemplates(nextTemplates);
  }, []);

  useEffect(() => {
    setTemplateWarning(null);
    setAppliedTemplate(null);
  }, [selectedDeviceId]);

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

    return buildRulePointDefinitions({
      rules,
      protocol: selectedDevice.protocol,
      locked: false,
    });
  }, [rules, selectedDevice]);

  const protectedPointDefinitions = useMemo(() => {
    if (!selectedDevice) {
      return [];
    }

    return buildRulePointDefinitions({
      rules,
      protocol: selectedDevice.protocol,
      locked: true,
    });
  }, [rules, selectedDevice]);

  const protectedAddressSet = useMemo(
    () => new Set(protectedPointDefinitions.map((definition) => definition.address)),
    [protectedPointDefinitions],
  );
  const readyToCreateCount = useMemo(
    () =>
      items.filter(
        (item) =>
          item.status === 'planned' &&
          item.mergeOffset === 0 &&
          !protectedAddressSet.has(item.address),
      ).length,
    [items, protectedAddressSet],
  );
  const conflictCount = useMemo(
    () => items.filter((item) => item.status === 'conflict' && item.mergeOffset === 0).length,
    [items],
  );
  const hasConflicts = items.some((item) => item.status === 'conflict');
  const selectedPointDefinition = useMemo(() => {
    if (!selectedAddress || !selectedDevice) {
      return null;
    }

    const selectedItem = items.find((item) => item.address === selectedAddress);
    if (!selectedItem || selectedItem.status !== 'planned' || !selectedItem.primaryRuleId) {
      return null;
    }

    const rule = rules.find(
      (candidate) => candidate.id === selectedItem.primaryRuleId && candidate.enabled,
    );
    if (!rule) {
      return null;
    }

    const rootAddress =
      selectedItem.mergeOffset > 0
        ? addressParser.offset(
            selectedItem.address,
            -selectedItem.mergeOffset,
            selectedDevice.protocol,
          )
        : selectedItem.address;

    return buildPointDefinition(rule, rootAddress);
  }, [items, rules, selectedAddress, selectedDevice]);

  const clearAppliedTemplate = () => {
    setAppliedTemplate(null);
    setTemplateWarning(null);
  };

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
      origin: appliedTemplate ? 'template' : 'manual',
      templateName: appliedTemplate?.name,
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

  const handleDeleteRule = (ruleId: string) => {
    setSourcePlanningState((currentState) => ({
      ...currentState,
      rules: currentState.rules.filter((rule) => rule.id !== ruleId),
      selectedRuleId:
        currentState.selectedRuleId === ruleId ? null : currentState.selectedRuleId,
      selectedAddress:
        currentState.selectedRuleId === ruleId ? null : currentState.selectedAddress,
    }));

    if (selectedRuleId === ruleId) {
      setFocusedRuleId(null);
      setInspectorSelection({ kind: 'none' });
    }

    setBatchCreateSummary(null);
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

  const handleSkipSelection = () => {
    setSourcePlanningState((currentState) => ({
      ...currentState,
      selectedAddress: null,
    }));
    setInspectorSelection({ kind: 'none' });
  };

  const handleStartRuleEdit = (ruleId: string) => {
    const rule = rules.find((candidate) => candidate.id === ruleId);
    if (!rule) {
      return;
    }

    setEditingRuleId(ruleId);
    setEditDraft({
      startAddress: rule.startAddress,
      count: rule.count,
      dataType: rule.dataType,
    });
  };

  const handleSaveRuleEdit = (ruleId: string) => {
    if (!editDraft) {
      return;
    }

    setSourcePlanningState((currentState) => ({
      ...currentState,
      rules: currentState.rules.map((rule) =>
        rule.id === ruleId
          ? {
              ...rule,
              startAddress: editDraft.startAddress.trim(),
              count: editDraft.count,
              dataType: editDraft.dataType,
            }
          : rule,
      ),
      selectedAddress: null,
    }));
    setEditingRuleId(null);
    setEditDraft(null);
    setBatchCreateSummary(null);
    setFocusedRuleId(ruleId);
    setInspectorSelection({ kind: 'rule', ruleId });
  };

  const handleCancelRuleEdit = () => {
    setEditingRuleId(null);
    setEditDraft(null);
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

  const handleOpenSaveTemplate = () => {
    setTemplateName(appliedTemplate?.name ?? '');
    setShowUtilityTools(true);
    setIsSaveTemplateOpen(true);
    setIsLoadTemplateOpen(false);
    setTemplateNotice(null);
  };

  const handleSaveTemplate = () => {
    const trimmedName = templateName.trim();
    if (!trimmedName) {
      setTemplateNotice(t('workbench.source.templates.nameRequired'));
      return;
    }

    const nextTemplate = createTemplateFromPlanner({
      templateName: trimmedName,
      draft: {
        startAddress,
        count,
        dataType,
      },
      preferredViewMode: viewMode,
      capabilitySnapshot: currentCapability ?? undefined,
    });
    const nextTemplates = sortTemplates(upsertTemplateRecord(templates, nextTemplate));

    saveSourceTemplates(nextTemplates);
    setTemplates(nextTemplates);
    setTemplateName('');
    setIsSaveTemplateOpen(false);
    setTemplateNotice(
      t('workbench.source.templates.saved', {
        name: nextTemplate.name,
      }),
    );
  };

  const handleOpenLoadTemplate = () => {
    setShowUtilityTools(true);
    setTemplates(sortTemplates(loadSourceTemplates()));
    setIsLoadTemplateOpen((currentValue) => !currentValue);
    setIsSaveTemplateOpen(false);
    setTemplateNotice(null);
  };

  const handleApplyTemplate = (template: SourceTemplateRecord) => {
    const plannerDraft = applyTemplateToPlanner(template);
    const now = new Date().toISOString();
    const nextTemplate = {
      ...template,
      lastUsedAt: now,
    };
    const nextTemplates = sortTemplates(upsertTemplateRecord(templates, nextTemplate));

    setStartAddress(plannerDraft.startAddress);
    setCount(plannerDraft.count);
    setDataType(plannerDraft.dataType);
    setViewMode(template.preferredViewMode ?? 'plan');
    setAppliedTemplate({
      id: template.id,
      name: template.name,
    });
    setTemplateWarning(buildTemplateWarning(template, currentCapability, t));
    setTemplateNotice(
      t('workbench.source.templates.loaded', {
        name: template.name,
      }),
    );
    setTemplates(nextTemplates);
    saveSourceTemplates(nextTemplates);
    setIsLoadTemplateOpen(false);
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

  const handleCreateSelectedPoint = async () => {
    if (!selectedDevice || !selectedPointDefinition) {
      return;
    }

    try {
      await createPointMutation.mutateAsync({
        device_id: selectedDevice.id,
        address: selectedPointDefinition.address,
        data_type: selectedPointDefinition.dataType,
        name: selectedPointDefinition.name,
      });
      setBatchCreateSummary({ successCount: 1, failureCount: 0 });
    } catch {
      setBatchCreateSummary({ successCount: 0, failureCount: 1 });
    }
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
      <div className="grid gap-3 rounded-2xl border border-slate-800 bg-slate-950/30 p-4 xl:grid-cols-[minmax(0,1fr)_180px_180px]">
        <div className="space-y-1">
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

        <article className="rounded-xl border border-slate-800/70 bg-slate-900/60 p-3">
          <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
            {t('workbench.runtime.summary.pointsHealthy')}
          </p>
          <p className="mt-2 text-lg font-semibold text-slate-50">
            {runtimeStatusQuery.data?.collectors[0]?.points_healthy ?? 0}
          </p>
        </article>

        <article className="rounded-xl border border-slate-800/70 bg-slate-900/60 p-3">
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
          <section
            className="space-y-3 rounded-2xl border border-slate-800/70 bg-slate-950/25 p-4"
            data-emphasis="supporting"
            data-testid="source-rule-layer"
          >
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

            <div className="space-y-3 rounded-xl border border-slate-800/70 bg-slate-950/70 p-3">
              <div className="space-y-1">
                <p className="text-[10px] font-bold uppercase tracking-[0.18em] text-cyan-300">
                  {t('workbench.source.planner.title')}
                </p>
                <p className="text-xs text-slate-400">
                  {t('workbench.source.planner.helper')}
                </p>
              </div>

              <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-1">
                <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                  <span>{t('workbench.source.planner.startAddress')}</span>
                  <input
                    aria-label={t('workbench.source.planner.startAddress')}
                    value={startAddress}
                    onChange={(event) => {
                      clearAppliedTemplate();
                      setStartAddress(event.target.value);
                    }}
                    className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
                  />
                </label>
                <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                  <span>{t('workbench.source.planner.dataType')}</span>
                  <select
                    aria-label={t('workbench.source.planner.dataType')}
                    value={dataType}
                    onChange={(event) => {
                      clearAppliedTemplate();
                      setDataType(event.target.value as DataType);
                    }}
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
                    min={1}
                    onChange={(event) => {
                      clearAppliedTemplate();
                      setCount(Number(event.target.value) || 0);
                    }}
                    type="number"
                    value={count}
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
              </div>

              <button
                type="button"
                onClick={handleApplyPlan}
                className="w-full rounded-lg bg-cyan-500 px-3 py-2 text-sm font-medium text-slate-950"
              >
                {t('workbench.source.planner.addRule')}
              </button>
            </div>

            {rules.length === 0 ? (
              <p className="rounded-xl border border-dashed border-slate-700/60 bg-slate-900/30 px-3 py-4 text-xs text-slate-400">
                {t('workbench.source.ruleLayer.empty')}
              </p>
            ) : (
              <div className="space-y-2">
                {rules.map((rule) => {
                  const coverage = buildSourceRuleCoverage(rule, selectedDevice.protocol);
                  const isSelected = rule.id === selectedRuleId;
                  const isEditing = editingRuleId === rule.id;
                  return (
                    <article
                      className={[
                        'space-y-2 rounded-xl border p-3 transition',
                        isSelected
                          ? 'border-cyan-500/40 bg-cyan-500/10'
                          : 'border-slate-800/70 bg-slate-900/50',
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

                      {isEditing && editDraft ? (
                        <div
                          className="space-y-2 rounded-lg border border-cyan-500/20 bg-slate-950/60 p-2"
                          data-testid="rule-inline-edit-form"
                        >
                          <label className="block space-y-1 text-[11px] uppercase tracking-[0.14em] text-slate-400">
                            <span>{t('workbench.source.planner.startAddress')}</span>
                            <input
                              aria-label={t('workbench.source.planner.startAddress')}
                              value={editDraft.startAddress}
                              onChange={(event) =>
                                setEditDraft((draft) =>
                                  draft ? { ...draft, startAddress: event.target.value } : draft,
                                )
                              }
                              className="w-full rounded-lg border border-slate-800 bg-slate-900 px-2 py-1 text-xs text-slate-100"
                            />
                          </label>
                          <label className="block space-y-1 text-[11px] uppercase tracking-[0.14em] text-slate-400">
                            <span>{t('workbench.source.planner.count')}</span>
                            <input
                              aria-label={t('workbench.source.planner.count')}
                              min={1}
                              type="number"
                              value={editDraft.count}
                              onChange={(event) =>
                                setEditDraft((draft) =>
                                  draft
                                    ? { ...draft, count: Number(event.target.value) || 0 }
                                    : draft,
                                )
                              }
                              className="w-full rounded-lg border border-slate-800 bg-slate-900 px-2 py-1 text-xs text-slate-100"
                            />
                          </label>
                          <label className="block space-y-1 text-[11px] uppercase tracking-[0.14em] text-slate-400">
                            <span>{t('workbench.source.planner.dataType')}</span>
                            <select
                              aria-label={t('workbench.source.planner.dataType')}
                              value={editDraft.dataType}
                              onChange={(event) =>
                                setEditDraft((draft) =>
                                  draft
                                    ? { ...draft, dataType: event.target.value as DataType }
                                    : draft,
                                )
                              }
                              className="w-full rounded-lg border border-slate-800 bg-slate-900 px-2 py-1 text-xs text-slate-100"
                            >
                              {SOURCE_PLANNER_ALLOWED_DATA_TYPES.map((option) => (
                                <option key={option} value={option}>
                                  {option}
                                </option>
                              ))}
                            </select>
                          </label>
                          <div className="flex gap-2">
                            <button
                              type="button"
                              onClick={() => handleSaveRuleEdit(rule.id)}
                              className="rounded-lg bg-cyan-500 px-2 py-1 text-[11px] font-medium text-slate-950"
                            >
                              {t('workbench.source.ruleLayer.editSave')}
                            </button>
                            <button
                              type="button"
                              onClick={handleCancelRuleEdit}
                              className="rounded-lg border border-slate-700/70 px-2 py-1 text-[11px] text-slate-300"
                            >
                              {t('workbench.source.ruleLayer.editCancel')}
                            </button>
                          </div>
                        </div>
                      ) : null}

                      <div className="grid grid-cols-2 gap-2 text-[11px]">
                        <button
                          className="rounded-lg border border-slate-700/70 px-2 py-1 text-slate-300"
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
                          className="rounded-lg border border-cyan-500/30 bg-cyan-500/5 px-2 py-1 text-cyan-200"
                          onClick={() => handleStartRuleEdit(rule.id)}
                          type="button"
                        >
                          {t('workbench.source.ruleLayer.editStart')}
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
                        <button
                          className="col-span-2 rounded-lg border border-rose-500/30 bg-rose-500/10 px-2 py-1 text-rose-100"
                          onClick={() => handleDeleteRule(rule.id)}
                          type="button"
                        >
                          {t('workbench.source.ruleLayer.delete')}
                        </button>
                      </div>
                    </article>
                  );
                })}
              </div>
            )}
          </section>
        </aside>

        <section
          className="space-y-4 rounded-2xl border border-slate-800 bg-slate-950/40 p-4"
          data-emphasis="primary"
          data-testid="source-canvas-workspace"
        >
          <div
            className="flex flex-col gap-3 xl:flex-row xl:items-start xl:justify-between"
            data-testid="source-primary-toolbar"
          >
            <div className="flex flex-wrap items-center gap-2">
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
            </div>
          </div>

          <AddressCanvas
            items={items}
            onSelectAddress={handleSelectAddress}
            selectedAddress={selectedAddress}
            valueFormat={valueFormat}
            viewMode={viewMode}
          />

          {selectedPointDefinition ? (
            <div
              className="flex flex-wrap items-center gap-2 rounded-xl border border-cyan-500/30 bg-cyan-500/5 px-3 py-2"
              data-testid="source-selection-toolbar"
            >
              <span className="text-xs font-medium text-cyan-200">
                {selectedPointDefinition.address}
              </span>
              <button
                type="button"
                onClick={() => void handleCreateSelectedPoint()}
                disabled={createPointMutation.isPending}
                className="rounded-lg bg-cyan-500 px-3 py-1.5 text-xs font-semibold text-slate-950 disabled:cursor-not-allowed disabled:bg-slate-800 disabled:text-slate-500"
              >
                {t('workbench.source.selectionToolbar.createPoints')}
              </button>
              <button
                type="button"
                onClick={() => handleSelectRule(
                  items.find((item) => item.address === selectedAddress)?.primaryRuleId ?? '',
                )}
                className="rounded-lg border border-slate-700/70 px-3 py-1.5 text-xs text-slate-300"
              >
                {t('workbench.source.selectionToolbar.addToRule')}
              </button>
              <button
                type="button"
                onClick={handleSkipSelection}
                className="rounded-lg border border-slate-700/70 px-3 py-1.5 text-xs text-slate-300"
              >
                {t('workbench.source.selectionToolbar.skip')}
              </button>
            </div>
          ) : null}

          <section
            className="space-y-2 rounded-xl border border-slate-800/60 bg-slate-950/20 p-3"
            data-testid="source-coverage-overview"
          >
            <div className="flex items-center justify-between gap-3">
              <div className="space-y-1">
                <p className="text-[10px] font-bold uppercase tracking-[0.2em] text-cyan-300">
                  {t('workbench.source.coverage.eyebrow')}
                </p>
                <h3 className="text-sm font-semibold text-slate-100">
                  {t('workbench.source.coverage.title')}
                </h3>
              </div>
              {coverageSegments.length > 0 ? (
                <p className="text-xs text-slate-400">
                  {t('workbench.source.coverage.cells', {
                    count: coverageSegments.reduce(
                      (total, segment) => total + segment.cellCount,
                      0,
                    ),
                  })}
                </p>
              ) : null}
            </div>
            {coverageSegments.length > 0 ? (
              <div className="space-y-3">
                <div className="flex gap-1">
                  {coverageSegments.map((segment) => (
                    <div
                      key={segment.id}
                      className={`h-2 rounded-full ${getCoverageSegmentClass(segment.status)}`}
                      style={{ flex: Math.max(segment.cellCount, 1) }}
                      title={`${segment.startAddress} → ${segment.endAddress}`}
                    />
                  ))}
                </div>
                <div className="flex flex-wrap gap-2">
                  {coverageSegments.map((segment) => (
                    <button
                      key={`${segment.id}-jump`}
                      className="rounded-lg border border-slate-800 bg-slate-900/60 px-3 py-2 text-left text-[11px] text-slate-300"
                      onClick={() => handleSelectAddress(segment.startAddress)}
                      type="button"
                    >
                      <p className="font-medium text-slate-100">
                        {segment.startAddress} → {segment.endAddress}
                      </p>
                      <p className="mt-1">
                        {t(`workbench.source.coverage.status.${segment.status}`)}
                      </p>
                </button>
              ))}
            </div>

            <section
              aria-label={t('workbench.source.summary.title')}
              className="grid gap-3 rounded-2xl border border-slate-800/70 bg-slate-900/60 p-3 xl:min-w-[420px]"
              data-testid="source-step-summary"
            >
              <div className="grid gap-2 sm:grid-cols-3">
                <article className="rounded-xl border border-slate-800/70 bg-slate-950/60 px-3 py-2">
                  <p className="text-[10px] font-bold uppercase tracking-[0.18em] text-slate-500">
                    {t('workbench.source.summary.readyToCreate')}
                  </p>
                  <p
                    className="mt-2 text-lg font-semibold text-emerald-200"
                    data-testid="source-summary-ready-count"
                  >
                    {readyToCreateCount}
                  </p>
                </article>
                <article className="rounded-xl border border-slate-800/70 bg-slate-950/60 px-3 py-2">
                  <p className="text-[10px] font-bold uppercase tracking-[0.18em] text-slate-500">
                    {t('workbench.source.summary.inConflict')}
                  </p>
                  <p
                    className="mt-2 text-lg font-semibold text-rose-200"
                    data-testid="source-summary-conflict-count"
                  >
                    {conflictCount}
                  </p>
                </article>
                <article className="rounded-xl border border-slate-800/70 bg-slate-950/60 px-3 py-2">
                  <p className="text-[10px] font-bold uppercase tracking-[0.18em] text-slate-500">
                    {t('workbench.source.summary.protected')}
                  </p>
                  <p
                    className="mt-2 text-lg font-semibold text-amber-200"
                    data-testid="source-summary-protected-count"
                  >
                    {protectedPointDefinitions.length}
                  </p>
                </article>
              </div>

              <div className="flex flex-wrap gap-2">
                <button
                  type="button"
                  onClick={() => void handleCreateSelectedPoint()}
                  disabled={!selectedPointDefinition || createPointMutation.isPending}
                  className="rounded-lg border border-cyan-500/40 bg-cyan-500/10 px-3 py-2 text-sm font-medium text-cyan-100 disabled:cursor-not-allowed disabled:border-slate-700 disabled:bg-slate-950/40 disabled:text-slate-500"
                >
                  {t('workbench.source.actions.createSelectedPoints')}
                </button>
                <button
                  type="button"
                  onClick={() => void handleBatchCreate()}
                  disabled={
                    plannedPointDefinitions.length === 0 ||
                    hasConflicts ||
                    createPointMutation.isPending
                  }
                  className="rounded-lg bg-cyan-500 px-3 py-2 text-sm font-semibold text-slate-950 disabled:cursor-not-allowed disabled:bg-slate-800 disabled:text-slate-500"
                >
                  {t('workbench.source.actions.createRulePoints')}
                </button>
              </div>
            </section>
          </div>
            ) : (
              <p className="text-xs text-slate-400">
                {t('workbench.source.coverage.empty')}
              </p>
            )}
          </section>

          <div
            className="space-y-3 rounded-xl border border-slate-800/60 bg-slate-950/20 px-3 py-3"
            data-testid="source-secondary-controls"
          >
            <div className="flex flex-wrap items-center gap-2">
              <button
                type="button"
                onClick={() => setShowUtilityTools((currentValue) => !currentValue)}
                className="rounded-lg border border-slate-800 bg-slate-900/60 px-3 py-2 text-xs text-slate-300"
              >
                {showUtilityTools
                  ? t('workbench.source.toolbar.hideTools')
                  : t('workbench.source.toolbar.moreTools')}
              </button>
            </div>

            {showUtilityTools ? (
              <>
                <div className="flex flex-wrap items-center gap-2">
                  <button
                    type="button"
                    onClick={handleToggleFreezeLive}
                    className="rounded-lg border border-slate-800 bg-slate-900/60 px-3 py-2 text-xs text-slate-300"
                  >
                    {freezeLive
                      ? t('workbench.source.toolbar.unfreezeLive')
                      : t('workbench.source.toolbar.freezeLive')}
                  </button>
                  <button
                    type="button"
                    onClick={handleCaptureSnapshot}
                    className="rounded-lg border border-slate-800 bg-slate-900/60 px-3 py-2 text-xs text-slate-300"
                  >
                    {t('workbench.source.toolbar.snapshotCompare')}
                  </button>
                  <button
                    type="button"
                    onClick={() => setShowAudit((currentValue) => !currentValue)}
                    className="rounded-lg border border-slate-800 bg-slate-900/60 px-3 py-2 text-xs text-slate-300"
                  >
                    {showAudit
                      ? t('workbench.source.toolbar.hideAudit')
                      : t('workbench.source.toolbar.showAudit')}
                  </button>
                  <button
                    type="button"
                    onClick={handleOpenSaveTemplate}
                    className="rounded-lg border border-slate-800 bg-slate-900/60 px-3 py-2 text-xs text-slate-300 disabled:opacity-50"
                  >
                    {t('workbench.source.toolbar.saveTemplate')}
                  </button>
                  <button
                    type="button"
                    onClick={handleOpenLoadTemplate}
                    disabled={templates.length === 0}
                    className="rounded-lg border border-slate-800 bg-slate-900/60 px-3 py-2 text-xs text-slate-300 disabled:opacity-50"
                  >
                    {t('workbench.source.toolbar.loadTemplate')}
                  </button>
                </div>

                <div className="flex flex-wrap items-end gap-2">
                  <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-500">
                    <span>{t('workbench.source.toolbar.jumpToAddress')}</span>
                    <input
                      aria-label={t('workbench.source.toolbar.jumpToAddress')}
                      value={jumpAddress}
                      onChange={(event) => setJumpAddress(event.target.value)}
                      className="w-36 rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
                    />
                  </label>
                  <button
                    type="button"
                    onClick={handleJumpToAddress}
                    className="rounded-lg border border-slate-800 bg-slate-900/60 px-3 py-2 text-xs text-slate-300"
                  >
                    {t('workbench.source.toolbar.jump')}
                  </button>
                </div>
              </>
            ) : null}
          </div>

          {isSaveTemplateOpen ? (
            <div className="space-y-3 rounded-xl border border-slate-800 bg-slate-900/50 p-4">
              <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                <span>{t('workbench.source.templates.name')}</span>
                <input
                  aria-label={t('workbench.source.templates.name')}
                  value={templateName}
                  onChange={(event) => setTemplateName(event.target.value)}
                  className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
                />
              </label>
              <div className="flex flex-wrap gap-2">
                <button
                  type="button"
                  onClick={handleSaveTemplate}
                  className="rounded-lg bg-cyan-500 px-3 py-2 text-sm font-medium text-slate-950"
                >
                  {t('workbench.source.templates.confirmSave')}
                </button>
                <button
                  type="button"
                  onClick={() => {
                    setIsSaveTemplateOpen(false);
                    setTemplateName('');
                  }}
                  className="rounded-lg border border-slate-800 px-3 py-2 text-sm text-slate-300"
                >
                  {t('workbench.source.templates.cancel')}
                </button>
              </div>
            </div>
          ) : null}

          {isLoadTemplateOpen ? (
            <div className="space-y-3 rounded-xl border border-slate-800 bg-slate-900/50 p-4">
              <div className="space-y-1">
                <p className="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">
                  {t('workbench.source.templates.library')}
                </p>
                <p className="text-xs text-slate-500">
                  {t('workbench.source.templates.description')}
                </p>
              </div>
              {templates.length > 0 ? (
                <div className="space-y-2">
                  {templates.map((template) => (
                    <button
                      key={template.id}
                      type="button"
                      aria-label={template.name}
                      onClick={() => handleApplyTemplate(template)}
                      className="block w-full rounded-xl border border-slate-800 bg-slate-950/60 px-3 py-3 text-left transition hover:border-cyan-500/40"
                    >
                      <p className="text-sm font-semibold text-slate-100">
                        {template.name}
                      </p>
                      <p className="mt-1 text-xs text-slate-400">
                        {template.startAddress} · {template.dataType} · {template.count}
                      </p>
                    </button>
                  ))}
                </div>
              ) : (
                <p className="text-xs text-slate-400">
                  {t('workbench.source.templates.empty')}
                </p>
              )}
            </div>
          ) : null}

          {hasConflicts ? (
            <p className="text-sm text-rose-300">
              {t('workbench.source.planner.conflictHint')}
            </p>
          ) : null}

          {templateWarning ? (
            <p className="text-sm text-amber-200" data-testid="source-template-warning">
              {templateWarning}
            </p>
          ) : null}

          {templateNotice ? (
            <p className="text-sm text-slate-300">{templateNotice}</p>
          ) : null}

          {batchCreateSummary ? (
            <p className="text-sm text-slate-300">
              {t('workbench.source.planner.batchSummary', batchCreateSummary)}
            </p>
          ) : null}

          {showAudit ? (
            <div className="rounded-xl border border-slate-800/60 bg-slate-950/20 p-3">
              <AddressLedger
                items={items}
                valueFormat={valueFormat}
                viewMode={viewMode}
              />
            </div>
          ) : null}
        </section>
      </div>
    </section>
  );
}
