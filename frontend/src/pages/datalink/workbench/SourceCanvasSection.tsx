import { ChevronDown, ChevronRight, CircleHelp } from 'lucide-react';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import {
  applyTemplateToPlanner,
  createTemplateFromPlanner,
  normalizeNamingPrefix,
  SOURCE_PLANNER_DATA_TYPE_GROUPS,
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
import { useDevicesQuery, useToggleDeviceStatusMutation } from '../../../hooks/datalink/useDevices';
import { useMappingsQuery } from '../../../hooks/datalink/useMappings';
import {
  useCreatePointMutation,
  useDeletePointMutation,
  usePointsQuery,
} from '../../../hooks/datalink/usePoints';
import { useRuntimeStream } from '../../../hooks/datalink/useRuntimeStream';
import {
  useCreateSourceRuleMutation,
  useDeleteSourceRuleMutation,
  useDisableSourceRuleMutation,
  useEnableSourceRuleMutation,
  useSourceRulesQuery,
  useUpdateSourceRuleMutation,
} from '../../../hooks/datalink/useSourceRules';
import { useTagsQuery } from '../../../hooks/datalink/useTags';
import { runtimeAPI } from '../../../services/datalink';
import type {
  DataType,
  Device,
  ProtocolType,
  SourceRuleRecord,
} from '../../../types/datalink';
import {
  addressParser,
  getDefaultPlannerStartAddress,
} from '../../../utils/addressParser';
import { AddressCanvas } from './AddressCanvas';
import { AddressLedger } from './AddressLedger';
import { useWorkbench } from './WorkbenchProvider';
import { parseDeviceConnectionConfig } from './workbenchDeviceFormModel';
import {
  buildAddressCanvasItems,
  buildConflictQueue,
  buildCoverageOverviewSegments,
  buildPlannedPointAddresses,
  buildSourceRuleCoverage,
  getDataTypeCellSpan,
  type SourceRule,
  type SourceValueFormat,
  type SourceViewMode,
} from './sourceCanvasModel';

function getSelectedDevice(devices: Device[], selectedDeviceId: string | null) {
  return devices.find((device) => device.id === selectedDeviceId) ?? null;
}

function getErrorMessage(error: unknown, fallback: string): string {
  return error instanceof Error ? error.message : fallback;
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
    case 'unmanaged':
      return 'bg-amber-500/80';
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

function mapPersistedRuleToPlannerRule(rule: SourceRuleRecord): SourceRule {
  return {
    id: rule.id,
    deviceId: rule.device_id,
    startAddress: rule.start_address,
    count: rule.count,
    dataType: rule.data_type,
    namingPrefix: rule.naming_prefix,
    enabled: rule.enabled,
    locked: rule.locked,
    origin: rule.origin,
    templateName: rule.template_name,
    skippedAddresses: rule.skipped_addresses ?? [],
    persisted: true,
    updatedAt: rule.updated_at,
    targetDataType: rule.target_data_type,
    scaleMultiplier: rule.scale_multiplier,
    scaleOffset: rule.scale_offset,
    dataFormat: rule.data_format,
  };
}

function areRulesEqual(left: ReadonlyArray<SourceRule>, right: ReadonlyArray<SourceRule>) {
  if (left.length !== right.length) {
    return false;
  }

  return left.every((rule, index) => {
    const candidate = right[index];
    return (
      rule.id === candidate.id &&
      rule.deviceId === candidate.deviceId &&
      rule.startAddress === candidate.startAddress &&
      rule.count === candidate.count &&
      rule.dataType === candidate.dataType &&
      rule.namingPrefix === candidate.namingPrefix &&
      rule.enabled === candidate.enabled &&
      rule.locked === candidate.locked &&
      rule.origin === candidate.origin &&
      rule.templateName === candidate.templateName &&
      rule.persisted === candidate.persisted &&
      rule.updatedAt === candidate.updatedAt &&
      rule.skippedAddresses.join(',') === candidate.skippedAddresses.join(',') &&
      rule.targetDataType === candidate.targetDataType &&
      rule.scaleMultiplier === candidate.scaleMultiplier &&
      rule.scaleOffset === candidate.scaleOffset &&
      rule.dataFormat === candidate.dataFormat
    );
  });
}

/**
 * 合併遠端持久化規則與本地草稿，避免相同 `id` 造成 React key 重複或狀態分裂。
 * 優先採用持久化列；草稿僅在該 `id` 尚未存在於持久化集合時才納入；兩側各自再依 id 去重。
 *
 * @param persistedRules - 後端同步的規則（可含重複 id，僅保留第一筆）
 * @param draftRules - 僅 `persisted === false` 的裝置內草稿
 * @returns 合併後、id 唯一的規則陣列（先持久化、後草稿）
 */
function mergePersistedRulesWithDrafts(
  persistedRules: ReadonlyArray<SourceRule>,
  draftRules: ReadonlyArray<SourceRule>,
): SourceRule[] {
  const uniquePersisted = [...new Map(persistedRules.map((rule) => [rule.id, rule])).values()];
  const persistedIds = new Set(uniquePersisted.map((rule) => rule.id));
  const draftsWithoutPersistedId = draftRules.filter((rule) => !persistedIds.has(rule.id));
  const uniqueDrafts = [
    ...new Map(draftsWithoutPersistedId.map((rule) => [rule.id, rule])).values(),
  ];
  return [...uniquePersisted, ...uniqueDrafts];
}

/**
 * 依列表順序去重相同 `id`，若同 id 多筆則優先保留已持久化（`persisted`）的那一筆。
 *
 * @param rules - 單一裝置下的規則列
 * @returns id 唯一且順序與首次出現一致（內容可能替換為持久化版本）
 */
function dedupeDeviceRulesPreferPersisted(rules: ReadonlyArray<SourceRule>): SourceRule[] {
  const indexById = new Map<string, number>();
  const out: SourceRule[] = [];

  for (const rule of rules) {
    const existingIndex = indexById.get(rule.id);
    if (existingIndex === undefined) {
      indexById.set(rule.id, out.length);
      out.push(rule);
      continue;
    }

    const previous = out[existingIndex]!;
    const next =
      rule.persisted && !previous.persisted ? rule : previous.persisted && !rule.persisted ? previous : rule;
    out[existingIndex] = next;
  }

  return out;
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
      if (rule.skippedAddresses?.includes(address)) return;
      if (!definitions.has(address)) {
        definitions.set(address, buildPointDefinition(rule, address));
      }
    });
  }

  return [...definitions.values()];
}

export function SourceCanvasSection() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const {
    selectedDeviceId,
    setFocusedRuleId,
    setInspectorSelection,
    setSelectedDeviceId,
    sourcePlanningState,
    setSourcePlanningState,
    setSourcePlannerStartAddress,
  } = useWorkbench();
  const { data: devices = [] } = useDevicesQuery();
  const { data: points = [] } = usePointsQuery(
    selectedDeviceId ? { device_id: selectedDeviceId } : undefined,
  );
  const { data: mappings = [] } = useMappingsQuery();
  const { data: tags = [] } = useTagsQuery();
  const createPointMutation = useCreatePointMutation();
  const deletePointMutation = useDeletePointMutation();
  const sourceRulesQuery = useSourceRulesQuery(
    selectedDeviceId ? { device_id: selectedDeviceId } : undefined,
  );
  const createSourceRuleMutation = useCreateSourceRuleMutation();
  const updateSourceRuleMutation = useUpdateSourceRuleMutation();
  const deleteSourceRuleMutation = useDeleteSourceRuleMutation();
  const enableSourceRuleMutation = useEnableSourceRuleMutation();
  const disableSourceRuleMutation = useDisableSourceRuleMutation();
  const toggleDeviceCollectionMutation = useToggleDeviceStatusMutation();
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
  const [targetDataType, setTargetDataType] = useState<DataType | ''>('');
  const [scaleMultiplier, setScaleMultiplier] = useState<string>('');
  const [scaleOffset, setScaleOffset] = useState<string>('');
  const [dataFormat, setDataFormat] = useState<string>('');
  const [jumpAddress, setJumpAddress] = useState('');
  /** 規則建立器表單區是否展開（收合時僅顯示標題列）。 */
  const [plannerSectionOpen, setPlannerSectionOpen] = useState(true);
  /** 已收合的規則卡片 id（Set 內表示該卡操作區／編輯區隱藏）。 */
  const [collapsedRuleCardIds, setCollapsedRuleCardIds] = useState(() => new Set<string>());
  const [editingRuleId, setEditingRuleId] = useState<string | null>(null);
  const [editDraft, setEditDraft] = useState<{
    startAddress: string;
    count: number;
    dataType: DataType;
    skippedAddresses: string[];
    targetDataType: DataType | '';
    scaleMultiplier: string;
    scaleOffset: string;
    dataFormat: string;
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
  /** 來源步驟表面訊息（規則 API 錯誤、設備收集開關結果等）。 */
  const [sourceStepNotice, setSourceStepNotice] = useState<string | null>(null);
  const [appliedTemplate, setAppliedTemplate] = useState<{
    id: string;
    name: string;
  } | null>(null);
  /** 來源步驟右側「覆蓋總覽／工具列」捲動區，用於攔截 wheel 避免傳到外層造成整塊版面微幅跟動。 */
  const sourceWorkspaceSecondaryScrollRef = useRef<HTMLDivElement>(null);
  const currentDeviceRules = useMemo(() => {
    if (!selectedDeviceId) {
      return [];
    }

    const forDevice = sourcePlanningState.rules.filter(
      (rule) => rule.deviceId === selectedDeviceId,
    );
    return dedupeDeviceRulesPreferPersisted(forDevice);
  }, [selectedDeviceId, sourcePlanningState.rules]);
  const rules = currentDeviceRules;
  const selectedRuleId = sourcePlanningState.selectedRuleId;
  const selectedAddress = sourcePlanningState.selectedAddress;
  const rememberedStartAddress = selectedDeviceId
    ? sourcePlanningState.plannerStartAddressByDeviceId[selectedDeviceId] ?? null
    : null;
  const currentCapability = useMemo(
    () => buildTemplateCapabilitySnapshot(selectedDevice),
    [selectedDevice],
  );
  const persistedRules = useMemo(
    () => (sourceRulesQuery.data ?? []).map(mapPersistedRuleToPlannerRule),
    [sourceRulesQuery.data],
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
    setSourceStepNotice(null);
  }, [selectedDeviceId]);

  /**
   * 在捲動區已抵頂或抵底時阻斷 wheel 的預設行為，避免捲動鏈把剩餘位移傳給祖先節點（外層 flex 版面會「微微跟著動」）。
   *
   * 須使用 `{ passive: false }` 才能於邊界呼叫 `preventDefault`；僅在頂／底過捲時觸發，不影響區塊內正常捲動。
   */
  useEffect(() => {
    if (!selectedDeviceId) {
      return undefined;
    }

    const el = sourceWorkspaceSecondaryScrollRef.current;
    if (!el) {
      return undefined;
    }

    const onWheel = (event: WheelEvent) => {
      const { scrollTop, scrollHeight, clientHeight } = el;
      const { deltaY } = event;
      const edgeSlack = 2;
      const atTop = scrollTop <= edgeSlack;
      const atBottom = scrollTop + clientHeight >= scrollHeight - edgeSlack;

      if ((atTop && deltaY < 0) || (atBottom && deltaY > 0)) {
        event.preventDefault();
      }
    };

    el.addEventListener('wheel', onWheel, { passive: false });
    return () => {
      el.removeEventListener('wheel', onWheel);
    };
  }, [selectedDeviceId]);

  useEffect(() => {
    if (!selectedDevice) {
      return;
    }

    const nextStartAddress =
      rememberedStartAddress && rememberedStartAddress.trim().length > 0
        ? rememberedStartAddress
        : getDefaultPlannerStartAddress(selectedDevice.protocol);
    setStartAddress(nextStartAddress);
  }, [rememberedStartAddress, selectedDevice]);

  useEffect(() => {
    if (!selectedDeviceId || !sourceRulesQuery.isSuccess) {
      return;
    }

    setSourcePlanningState((currentState) => {
      const deviceRules = currentState.rules.filter((rule) => rule.deviceId === selectedDeviceId);
      const otherRules = currentState.rules.filter((rule) => rule.deviceId !== selectedDeviceId);
      const draftRules = deviceRules.filter((rule) => !rule.persisted);
      const nextRules = mergePersistedRulesWithDrafts(persistedRules, draftRules);
      if (areRulesEqual(deviceRules, nextRules)) {
        return currentState;
      }

      const nextSelectedRuleId = nextRules.some(
        (rule) => rule.id === currentState.selectedRuleId,
      )
        ? currentState.selectedRuleId
        : null;

      return {
        ...currentState,
        rules: [...otherRules, ...nextRules],
        selectedRuleId: nextSelectedRuleId,
        selectedAddress: nextSelectedRuleId ? currentState.selectedAddress : null,
      };
    });
  }, [persistedRules, selectedDeviceId, setSourcePlanningState, sourceRulesQuery.isSuccess]);

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
  const conflictQueue = useMemo(() => buildConflictQueue(items), [items]);

  const safePointDefinitions = useMemo(() => {
    if (!selectedDevice) return [];
    const itemMap = new Map(items.map((item) => [item.address, item]));

    return plannedPointDefinitions.filter((definition) => {
      const span = getDataTypeCellSpan(definition.dataType);
      const occupied = addressParser.expand(
        definition.address,
        span,
        selectedDevice.protocol,
      );
      return occupied.every((addr) => {
        const cell = itemMap.get(addr);
        return cell && cell.status === 'planned';
      });
    });
  }, [items, plannedPointDefinitions, selectedDevice]);

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
      deviceId: selectedDevice.id,
      startAddress: startAddress.trim(),
      count,
      dataType,
      namingPrefix,
      enabled: true,
      locked: false,
      origin: appliedTemplate ? 'template' : 'manual',
      templateName: appliedTemplate?.name,
      skippedAddresses: [],
      persisted: false,
      targetDataType: targetDataType || undefined,
      scaleMultiplier: scaleMultiplier ? Number(scaleMultiplier) : undefined,
      scaleOffset: scaleOffset ? Number(scaleOffset) : undefined,
      dataFormat: dataFormat || undefined,
    };

    setSourcePlanningState((currentState) => ({
      ...currentState,
      rules: [
        ...currentState.rules.filter((rule) => rule.deviceId !== selectedDevice.id),
        ...currentState.rules.filter((rule) => rule.deviceId === selectedDevice.id),
        nextRule,
      ],
      selectedRuleId: ruleId,
      selectedAddress: null,
    }));
    setBatchCreateSummary(null);
    setFocusedRuleId(ruleId);
    setInspectorSelection({ kind: 'rule', ruleId });
  };

  const handleDeleteRule = async (ruleId: string) => {
    const rule = rules.find((candidate) => candidate.id === ruleId);
    if (!rule) {
      return;
    }

    if (rule.persisted) {
      await deleteSourceRuleMutation.mutateAsync(ruleId);
      setSourcePlanningState((currentState) => ({
        ...currentState,
        rules: [
          ...currentState.rules.filter((candidate) => candidate.deviceId !== selectedDeviceId),
          ...currentState.rules.filter(
            (candidate) => candidate.deviceId === selectedDeviceId && candidate.id !== ruleId,
          ),
        ],
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
      return;
    }

    if (selectedDevice) {
      const plannedAddresses = new Set(
        buildPlannedPointAddresses({
          startAddress: rule.startAddress,
          count: rule.count,
          dataType: rule.dataType,
          protocol: selectedDevice.protocol,
        }).filter((addr) => !rule.skippedAddresses?.includes(addr)),
      );

      const orphanedPoints = points.filter((point) => plannedAddresses.has(point.address));
      for (const point of orphanedPoints) {
        deletePointMutation.mutate(point.id);
      }
    }

    setSourcePlanningState((currentState) => ({
      ...currentState,
      rules: [
        ...currentState.rules.filter((rule) => rule.deviceId !== selectedDeviceId),
        ...currentState.rules.filter(
          (candidate) => candidate.deviceId === selectedDeviceId && candidate.id !== ruleId,
        ),
      ],
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

  /**
   * 切換左欄單一規則卡片的展開／收合（不影響選取狀態）。
   */
  const toggleRuleCardCollapsed = (ruleId: string) => {
    setCollapsedRuleCardIds((prev) => {
      const next = new Set(prev);
      if (next.has(ruleId)) {
        next.delete(ruleId);
      } else {
        next.add(ruleId);
      }
      return next;
    });
  };

  const handleSelectAddress = (address: string) => {
    const clickedItem = items.find((item) => item.address === address);

    // Logical cell guard: snap to root address for multi-word values
    let resolvedAddress = address;
    if (clickedItem && clickedItem.mergeSpan > 1 && clickedItem.mergeOffset > 0 && selectedDevice) {
      resolvedAddress = addressParser.offset(
        address,
        -clickedItem.mergeOffset,
        selectedDevice.protocol,
      );
    }

    const resolvedItem = resolvedAddress !== address
      ? items.find((item) => item.address === resolvedAddress)
      : clickedItem;

    setSourcePlanningState((currentState) => ({
      ...currentState,
      selectedAddress: resolvedAddress,
      selectedRuleId: resolvedItem?.primaryRuleId ?? currentState.selectedRuleId,
    }));
    if (resolvedItem?.primaryRuleId) {
      setFocusedRuleId(resolvedItem.primaryRuleId);
    }
    setInspectorSelection({
      kind: 'span',
      spanAddress: resolvedAddress,
      ruleId: resolvedItem?.primaryRuleId ?? undefined,
    });
  };

  /**
   * 切換單一來源規則的 `enabled`（後端已持久化規則走 enable/disable API；草稿僅改本地狀態）。
   * 此開關不控制設備輪詢，僅影響該規則是否納管。
   */
  const handleToggleRuleEnabled = async (ruleId: string) => {
    const rule = rules.find((candidate) => candidate.id === ruleId);
    if (!rule) {
      return;
    }

    if (rule.persisted) {
      try {
        if (rule.enabled) {
          await disableSourceRuleMutation.mutateAsync(ruleId);
        } else {
          await enableSourceRuleMutation.mutateAsync(ruleId);
        }
        setSourceStepNotice(null);
      } catch (error) {
        setSourceStepNotice(
          getErrorMessage(error, t('workbench.source.collection.toggleRuleError')),
        );
      }
      return;
    }

    setSourcePlanningState((currentState) => ({
      ...currentState,
      rules: [
        ...currentState.rules.filter((candidate) => candidate.deviceId !== selectedDeviceId),
        ...currentState.rules
          .filter((candidate) => candidate.deviceId === selectedDeviceId)
          .map((candidate) =>
            candidate.id === ruleId ? { ...candidate, enabled: !candidate.enabled } : candidate,
          ),
      ],
    }));
  };

  /**
   * 依目前選擇的設備切換「是否由 Runtime 輪詢收集」（啟用設備 / 停用設備 API）。
   */
  const handleDeviceCollectionToggle = async () => {
    if (!selectedDevice) {
      return;
    }

    const wasActive = selectedDevice.status === 'active';

    try {
      await toggleDeviceCollectionMutation.mutateAsync({
        id: selectedDevice.id,
        currentStatus: selectedDevice.status,
      });
      await queryClient.invalidateQueries({ queryKey: ['runtime-status', selectedDeviceId] });
      setSourceStepNotice(
        wasActive ? t('workbench.source.collection.stopped') : t('workbench.source.collection.started'),
      );
    } catch (error) {
      setSourceStepNotice(getErrorMessage(error, t('workbench.source.collection.error')));
    }
  };

  const isRuleEnableMutating =
    disableSourceRuleMutation.isPending || enableSourceRuleMutation.isPending;

  const handleSkipConflictSpan = (ruleId: string, conflictCellAddress: string) => {
    if (!selectedDevice) return;

    setSourcePlanningState((currentState) => ({
      ...currentState,
      rules: [
        ...currentState.rules.filter((rule) => rule.deviceId !== selectedDeviceId),
        ...currentState.rules
          .filter((rule) => rule.deviceId === selectedDeviceId)
          .map((rule) => {
        if (rule.id !== ruleId) return rule;

        const plannedAddresses = buildPlannedPointAddresses({
          startAddress: rule.startAddress,
          count: rule.count,
          dataType: rule.dataType,
          protocol: selectedDevice.protocol,
        });

        const spanRoot = plannedAddresses.find((pointAddr) => {
          const occupied = addressParser.expand(
            pointAddr,
            getDataTypeCellSpan(rule.dataType),
            selectedDevice.protocol,
          );
          return occupied.includes(conflictCellAddress);
        });

        if (!spanRoot) return rule;

        return {
          ...rule,
          skippedAddresses: [...(rule.skippedAddresses ?? []), spanRoot],
        };
          }),
      ],
    }));
    setBatchCreateSummary(null);
  };

  const handleToggleRuleLocked = (ruleId: string) => {
    setSourcePlanningState((currentState) => ({
      ...currentState,
      rules: [
        ...currentState.rules.filter((rule) => rule.deviceId !== selectedDeviceId),
        ...currentState.rules
          .filter((rule) => rule.deviceId === selectedDeviceId)
          .map((rule) => (rule.id === ruleId ? { ...rule, locked: !rule.locked } : rule)),
      ],
    }));
  };

  const handleMoveRule = (ruleId: string, direction: -1 | 1) => {
    setSourcePlanningState((currentState) => ({
      ...currentState,
      rules: [
        ...currentState.rules.filter((rule) => rule.deviceId !== selectedDeviceId),
        ...moveRule(
          currentState.rules.filter((rule) => rule.deviceId === selectedDeviceId),
          ruleId,
          direction,
        ),
      ],
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
      skippedAddresses: [...rule.skippedAddresses],
      targetDataType: rule.targetDataType ?? '',
      scaleMultiplier: rule.scaleMultiplier?.toString() ?? '',
      scaleOffset: rule.scaleOffset?.toString() ?? '',
      dataFormat: rule.dataFormat ?? '',
    });
  };

  const handleSaveRuleEdit = async (ruleId: string) => {
    if (!editDraft) {
      return;
    }

    const rule = rules.find((candidate) => candidate.id === ruleId);
    if (!rule) {
      return;
    }

    if (rule.persisted) {
      await updateSourceRuleMutation.mutateAsync({
        id: ruleId,
        data: {
          start_address: editDraft.startAddress.trim(),
          count: editDraft.count,
          data_type: editDraft.dataType,
          skipped_addresses: editDraft.skippedAddresses,
          target_data_type: editDraft.targetDataType || null,
          scale_multiplier: editDraft.scaleMultiplier ? Number(editDraft.scaleMultiplier) : null,
          scale_offset: editDraft.scaleOffset ? Number(editDraft.scaleOffset) : null,
          data_format:
            editDraft.dataFormat.trim() === ''
              ? null
              : editDraft.dataFormat.trim().toUpperCase(),
        },
      });
      setEditingRuleId(null);
      setEditDraft(null);
      setBatchCreateSummary(null);
      setFocusedRuleId(ruleId);
      setInspectorSelection({ kind: 'rule', ruleId });
      return;
    }

    setSourcePlanningState((currentState) => ({
      ...currentState,
      rules: [
        ...currentState.rules.filter((rule) => rule.deviceId !== selectedDeviceId),
        ...currentState.rules
          .filter((rule) => rule.deviceId === selectedDeviceId)
          .map((rule) =>
            rule.id === ruleId
              ? {
                  ...rule,
                  startAddress: editDraft.startAddress.trim(),
                  count: editDraft.count,
                  dataType: editDraft.dataType,
                  skippedAddresses: editDraft.skippedAddresses,
                  targetDataType: editDraft.targetDataType || undefined,
                  scaleMultiplier: editDraft.scaleMultiplier ? Number(editDraft.scaleMultiplier) : undefined,
                  scaleOffset: editDraft.scaleOffset ? Number(editDraft.scaleOffset) : undefined,
                  dataFormat: editDraft.dataFormat || undefined,
                }
              : rule,
          ),
      ],
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
    if (selectedDeviceId) {
      setSourcePlannerStartAddress(selectedDeviceId, plannerDraft.startAddress);
    }
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
    if (!selectedDevice || safePointDefinitions.length === 0) {
      return;
    }

    const safeAddressSet = new Set(safePointDefinitions.map((definition) => definition.address));
    const eligibleRules = rules
      .filter((rule) => !rule.persisted && rule.enabled)
      .map((rule) => {
        const ruleDefinitions = buildRulePointDefinitions({
          rules: [rule],
          protocol: selectedDevice.protocol,
        });
        const conflictingRoots = ruleDefinitions
          .filter((definition) => !safeAddressSet.has(definition.address))
          .map((definition) => definition.address);
        const skippedAddresses = [...new Set([...(rule.skippedAddresses ?? []), ...conflictingRoots])];
        const safeDefinitionCount = ruleDefinitions.length - conflictingRoots.length;

        return {
          rule,
          skippedAddresses,
          safeDefinitionCount,
        };
      })
      .filter((candidate) => candidate.safeDefinitionCount > 0);

    if (eligibleRules.length === 0) {
      return;
    }

    const results = await Promise.allSettled(
      eligibleRules.map(({ rule, skippedAddresses }) =>
        createSourceRuleMutation.mutateAsync({
          id: rule.id,
          device_id: selectedDevice.id,
          start_address: rule.startAddress,
          count: rule.count,
          data_type: rule.dataType,
          naming_prefix: rule.namingPrefix,
          enabled: rule.enabled,
          locked: rule.locked,
          origin: rule.origin,
          template_name: rule.templateName,
          skipped_addresses: skippedAddresses,
          target_data_type: rule.targetDataType,
          scale_multiplier: rule.scaleMultiplier,
          scale_offset: rule.scaleOffset,
          data_format: rule.dataFormat,
        }),
      ),
    );

    const succeededRuleIds = eligibleRules
      .filter((_, index) => results[index]?.status === 'fulfilled')
      .map(({ rule }) => rule.id);

    if (succeededRuleIds.length > 0) {
      setSourcePlanningState((currentState) => ({
        ...currentState,
        rules: [
          ...currentState.rules.filter((rule) => rule.deviceId !== selectedDeviceId),
          ...currentState.rules.filter(
            (rule) =>
              rule.deviceId === selectedDeviceId &&
              (rule.persisted || !succeededRuleIds.includes(rule.id)),
          ),
        ],
        selectedRuleId:
          currentState.selectedRuleId && succeededRuleIds.includes(currentState.selectedRuleId)
            ? null
            : currentState.selectedRuleId,
        selectedAddress: null,
      }));
    }

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
      <section className="min-h-0 flex-1 space-y-6 overflow-y-auto overscroll-contain rounded-2xl border border-dashed border-slate-700 bg-slate-950/40 p-6">
        <div className="space-y-2">
          <p className="text-xs font-semibold uppercase tracking-[0.18em] text-cyan-300">
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
    <section className="flex h-full min-h-0 flex-col gap-6 overflow-hidden">
      <div
        className="shrink-0 rounded-2xl border border-slate-800 bg-slate-950/30 p-4"
        data-testid="source-runtime-collection-panel"
      >
        <div className="flex flex-col gap-4 lg:flex-row lg:items-stretch lg:justify-between lg:gap-6">
          <div className="min-w-0 flex-1 space-y-3">
            <p className="text-xs font-semibold uppercase tracking-[0.18em] text-cyan-300">
              {t('workbench.runtime.summary.title')}
            </p>
            <div className="flex flex-wrap items-end gap-x-4 gap-y-2">
              <div className="min-w-0 space-y-1">
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
              <div className="flex min-w-0 flex-wrap items-center gap-2 border-slate-800/70 lg:border-l lg:pl-4">
                <span className="text-xs font-medium text-slate-200">{selectedDevice.name}</span>
                <span className="rounded-full border border-slate-700/80 px-2 py-0.5 text-[10px] uppercase tracking-[0.12em] text-slate-400">
                  {selectedDevice.status}
                </span>
                <button
                  type="button"
                  data-testid="source-device-collection-toggle"
                  disabled={toggleDeviceCollectionMutation.isPending}
                  onClick={() => void handleDeviceCollectionToggle()}
                  className={
                    selectedDevice.status === 'active'
                      ? 'rounded-lg border border-rose-500/40 bg-rose-500/10 px-3 py-2 text-xs font-semibold text-rose-100 disabled:cursor-not-allowed disabled:opacity-50'
                      : 'rounded-lg bg-cyan-500 px-3 py-2 text-xs font-semibold text-slate-950 disabled:cursor-not-allowed disabled:opacity-50'
                  }
                >
                  {selectedDevice.status === 'active'
                    ? t('workbench.source.collection.stop')
                    : t('workbench.source.collection.start')}
                </button>
              </div>
            </div>
            <p
              className="text-[11px] leading-relaxed text-slate-500"
              title={t('workbench.source.collection.hint')}
            >
              {t('workbench.source.collection.hintShort')}
            </p>
            {sourceStepNotice ? (
              <p className="text-xs text-amber-200" data-testid="source-step-notice">
                {sourceStepNotice}
              </p>
            ) : null}
          </div>

          <div className="flex shrink-0 flex-row gap-3 sm:gap-4">
            <article className="min-w-[140px] flex-1 rounded-xl border border-slate-800/70 bg-slate-900/60 p-3 sm:min-w-[160px]">
              <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
                {t('workbench.runtime.summary.pointsHealthy')}
              </p>
              <p className="mt-2 text-lg font-semibold text-slate-50">
                {runtimeStatusQuery.data?.collectors[0]?.points_healthy ?? 0}
              </p>
            </article>
            <article className="min-w-[140px] flex-1 rounded-xl border border-slate-800/70 bg-slate-900/60 p-3 sm:min-w-[160px]">
              <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
                {t('workbench.runtime.summary.pointsStale')}
              </p>
              <p className="mt-2 text-lg font-semibold text-slate-50">
                {runtimeStatusQuery.data?.collectors[0]?.points_stale ?? 0}
              </p>
            </article>
          </div>
        </div>
      </div>

      <div className="grid min-h-0 flex-1 gap-4 overflow-hidden xl:grid-cols-[300px_minmax(0,1fr)]">
        <aside className="flex min-h-0 flex-col gap-4 overflow-hidden">
          <section
            className="flex min-h-0 flex-1 flex-col gap-3 overflow-hidden rounded-2xl border border-slate-800/70 bg-slate-950/25 p-4"
            data-emphasis="supporting"
            data-testid="source-rule-layer"
          >
            <div className="shrink-0 space-y-1">
              <p className="text-xs font-semibold uppercase tracking-[0.18em] text-cyan-300">
                {t('workbench.source.ruleLayer.eyebrow')}
              </p>
              <div className="flex items-start gap-2">
                <h3 className="min-w-0 flex-1 text-sm font-semibold leading-snug text-slate-100">
                  {t('workbench.source.ruleLayer.title')}
                </h3>
                <button
                  type="button"
                  data-testid="source-rule-layer-description-hint"
                  className="mt-0.5 shrink-0 rounded-md p-0.5 text-slate-400 transition-colors hover:bg-slate-800/80 hover:text-slate-200 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-cyan-400"
                  aria-label={t('workbench.source.ruleLayer.descriptionHint')}
                  title={t('workbench.source.ruleLayer.description')}
                >
                  <CircleHelp className="h-4 w-4" aria-hidden />
                </button>
              </div>
            </div>

            <div className="shrink-0 space-y-3 rounded-xl border border-slate-800/70 bg-slate-950/70 p-3">
              <button
                type="button"
                data-testid="source-planner-section-toggle"
                aria-expanded={plannerSectionOpen}
                aria-label={
                  plannerSectionOpen
                    ? t('workbench.source.planner.collapseSection')
                    : t('workbench.source.planner.expandSection')
                }
                onClick={() => setPlannerSectionOpen((open) => !open)}
                className="flex w-full items-start gap-2 rounded-lg text-left transition-colors hover:bg-slate-800/40 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-cyan-400"
              >
                <span className="mt-0.5 shrink-0 text-slate-400" aria-hidden>
                  {plannerSectionOpen ? (
                    <ChevronDown className="h-4 w-4" />
                  ) : (
                    <ChevronRight className="h-4 w-4" />
                  )}
                </span>
                <span className="min-w-0 flex-1 space-y-1">
                  <span className="block text-xs font-semibold uppercase tracking-[0.18em] text-cyan-300">
                    {t('workbench.source.planner.title')}
                  </span>
                  {plannerSectionOpen ? (
                    <span className="block text-xs text-slate-400">
                      {t('workbench.source.planner.helper')}
                    </span>
                  ) : (
                    <span className="block text-[11px] text-slate-500">
                      {t('workbench.source.planner.collapsedHint')}
                    </span>
                  )}
                </span>
              </button>

              {plannerSectionOpen ? (
              <>
              <div className="grid gap-3 sm:grid-cols-2">
                <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                  <span>{t('workbench.source.planner.startAddress')}</span>
                  <input
                    aria-label={t('workbench.source.planner.startAddress')}
                    value={startAddress}
                    onChange={(event) => {
                      clearAppliedTemplate();
                      const nextValue = event.target.value;
                      setStartAddress(nextValue);
                      if (selectedDeviceId) {
                        setSourcePlannerStartAddress(selectedDeviceId, nextValue);
                      }
                    }}
                    className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
                  />
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
                    {SOURCE_PLANNER_DATA_TYPE_GROUPS.map((group) => (
                      <optgroup key={group.labelKey} label={t(group.labelKey)}>
                        {group.types.map((entry) => (
                          <option
                            key={entry.value}
                            value={entry.value}
                            disabled={!entry.supported}
                          >
                            {entry.value}
                            {!entry.supported && entry.disabledReasonKey
                              ? ` — ${t(entry.disabledReasonKey)}`
                              : ''}
                          </option>
                        ))}
                      </optgroup>
                    ))}
                  </select>
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

              {/* 進階設定：目標型態／字節序、偏移／倍率（2×2 網格） */}
              <div className="space-y-3">
                <p className="text-xs font-semibold uppercase tracking-[0.18em] text-cyan-300">
                  {t('workbench.source.planner.advanced.heading')}
                </p>
                <div className="grid gap-3 sm:grid-cols-2">
                  <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                    <span>{t('workbench.source.planner.targetDataType')}</span>
                    <select
                      aria-label={t('workbench.source.planner.targetDataType')}
                      value={targetDataType}
                      onChange={(event) => {
                        clearAppliedTemplate();
                        setTargetDataType(event.target.value as DataType | '');
                      }}
                      className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
                    >
                      <option value="">{t('workbench.source.planner.sameAsReadType')}</option>
                      {SOURCE_PLANNER_DATA_TYPE_GROUPS.flatMap((group) =>
                        group.types
                          .filter((entry) => entry.supported)
                          .map((entry) => (
                            <option key={entry.value} value={entry.value}>
                              {entry.value}
                            </option>
                          )),
                      )}
                    </select>
                  </label>
                  <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                    <span>{t('workbench.source.planner.dataFormat')}</span>
                    <select
                      aria-label={t('workbench.source.planner.dataFormat')}
                      value={dataFormat}
                      onChange={(event) => {
                        clearAppliedTemplate();
                        setDataFormat(event.target.value);
                      }}
                      className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
                    >
                      <option value="">{t('workbench.source.planner.dataFormatDefault')}</option>
                      <option value="CDAB">CDAB</option>
                      <option value="ABCD">ABCD</option>
                      <option value="BADC">BADC</option>
                      <option value="DCBA">DCBA</option>
                    </select>
                  </label>
                  <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                    <span>{t('workbench.source.planner.scaleOffset')}</span>
                    <input
                      aria-label={t('workbench.source.planner.scaleOffset')}
                      type="number"
                      step="any"
                      placeholder="0.0"
                      value={scaleOffset}
                      onChange={(event) => setScaleOffset(event.target.value)}
                      className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
                    />
                  </label>
                  <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                    <span>{t('workbench.source.planner.scaleMultiplier')}</span>
                    <input
                      aria-label={t('workbench.source.planner.scaleMultiplier')}
                      type="number"
                      step="any"
                      placeholder="1.0"
                      value={scaleMultiplier}
                      onChange={(event) => setScaleMultiplier(event.target.value)}
                      className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
                    />
                  </label>
                </div>
              </div>

              <button
                type="button"
                onClick={handleApplyPlan}
                className="w-full rounded-lg bg-cyan-500 px-3 py-2 text-sm font-medium text-slate-950"
              >
                {t('workbench.source.planner.addRule')}
              </button>
              </>
              ) : null}
            </div>

            <div className="min-h-0 flex-1 overflow-y-auto overscroll-contain pr-0.5 scrollbar-auto-hide">
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
                  const isRuleCardCollapsed = collapsedRuleCardIds.has(rule.id);
                  return (
                    <article
                      className={[
                        'space-y-2 rounded-xl border p-3 transition',
                        isSelected
                          ? 'border-cyan-500/40 bg-cyan-500/10'
                          : 'border-slate-800/70 bg-slate-900/50',
                      ].join(' ')}
                      data-testid={`source-rule-${rule.id}-card`}
                      key={rule.id}
                    >
                      <div className="flex items-start gap-1">
                        <button
                          type="button"
                          className="mt-0.5 shrink-0 rounded-md p-1 text-slate-400 transition-colors hover:bg-slate-800/80 hover:text-slate-200 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-cyan-400"
                          aria-expanded={!isRuleCardCollapsed}
                          aria-label={
                            isRuleCardCollapsed
                              ? t('workbench.source.ruleLayer.expandCard')
                              : t('workbench.source.ruleLayer.collapseCard')
                          }
                          data-testid={`source-rule-${rule.id}-collapse`}
                          onClick={(event) => {
                            event.stopPropagation();
                            toggleRuleCardCollapsed(rule.id);
                          }}
                        >
                          {isRuleCardCollapsed ? (
                            <ChevronRight className="h-4 w-4" aria-hidden />
                          ) : (
                            <ChevronDown className="h-4 w-4" aria-hidden />
                          )}
                        </button>
                        <button
                          className="min-w-0 flex-1 space-y-2 text-left"
                          data-testid={`source-rule-${rule.id}`}
                          onClick={() => handleSelectRule(rule.id)}
                          type="button"
                        >
                          <div className="flex items-center justify-between gap-3">
                            <div className="flex items-center gap-2">
                              <span className="text-sm font-semibold text-slate-100">
                                {rule.id}
                              </span>
                              <span
                                className={[
                                  'rounded-full px-2 py-0.5 text-[10px] font-bold uppercase tracking-[0.12em]',
                                  rule.persisted
                                    ? 'border border-emerald-500/30 bg-emerald-500/10 text-emerald-200'
                                    : 'border border-slate-700/70 bg-slate-900/70 text-slate-300',
                                ].join(' ')}
                              >
                                {rule.persisted
                                  ? t('workbench.source.ruleLayer.persistedBadge')
                                  : t('workbench.source.ruleLayer.draftBadge')}
                              </span>
                            </div>
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
                      </div>

                      {!isRuleCardCollapsed && isEditing && editDraft ? (
                        <div
                          className="space-y-2 rounded-lg border border-cyan-500/20 bg-slate-950/60 p-2 text-[11px] text-slate-300"
                          data-testid="rule-inline-edit-form"
                        >
                          <div className="grid gap-2 sm:grid-cols-2">
                            <label className="space-y-1 text-[11px] uppercase tracking-[0.16em] text-slate-400">
                              <span>{t('workbench.source.planner.startAddress')}</span>
                              <input
                                aria-label={t('workbench.source.planner.startAddress')}
                                value={editDraft.startAddress}
                                onChange={(event) =>
                                  setEditDraft((draft) =>
                                    draft ? { ...draft, startAddress: event.target.value } : draft,
                                  )
                                }
                                className="w-full rounded-lg border border-slate-800 bg-slate-900 px-2 py-1 text-[11px] text-slate-100"
                              />
                            </label>
                            <label className="space-y-1 text-[11px] uppercase tracking-[0.16em] text-slate-400">
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
                                className="w-full rounded-lg border border-slate-800 bg-slate-900 px-2 py-1 text-[11px] text-slate-100"
                              />
                            </label>
                            <label className="space-y-1 text-[11px] uppercase tracking-[0.16em] text-slate-400 sm:col-span-2">
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
                                className="w-full rounded-lg border border-slate-800 bg-slate-900 px-2 py-1 text-[11px] text-slate-100"
                              >
                                {SOURCE_PLANNER_DATA_TYPE_GROUPS.map((group) => (
                                  <optgroup key={group.labelKey} label={t(group.labelKey)}>
                                    {group.types.map((entry) => (
                                      <option
                                        key={entry.value}
                                        value={entry.value}
                                        disabled={!entry.supported}
                                      >
                                        {entry.value}
                                        {!entry.supported && entry.disabledReasonKey
                                          ? ` — ${t(entry.disabledReasonKey)}`
                                          : ''}
                                      </option>
                                    ))}
                                  </optgroup>
                                ))}
                              </select>
                            </label>
                          </div>
                          {/* 進階設定：與規則建立器同樣為常駐雙欄；字級與同卡按鈕列 text-[11px] 對齊 */}
                          <div className="space-y-2">
                            <p className="text-[11px] font-semibold uppercase tracking-[0.18em] text-cyan-300">
                              {t('workbench.source.planner.advanced.heading')}
                            </p>
                            <div className="grid gap-2 sm:grid-cols-2">
                              <label className="space-y-1 text-[11px] uppercase tracking-[0.16em] text-slate-400">
                                <span>{t('workbench.source.planner.targetDataType')}</span>
                                <select
                                  aria-label={t('workbench.source.planner.targetDataType')}
                                  value={editDraft.targetDataType}
                                  onChange={(event) =>
                                    setEditDraft((draft) =>
                                      draft
                                        ? { ...draft, targetDataType: event.target.value as DataType | '' }
                                        : draft,
                                    )
                                  }
                                  className="w-full rounded-lg border border-slate-800 bg-slate-900 px-2 py-1 text-[11px] text-slate-100"
                                >
                                  <option value="">{t('workbench.source.planner.sameAsReadType')}</option>
                                  {SOURCE_PLANNER_DATA_TYPE_GROUPS.flatMap((group) =>
                                    group.types
                                      .filter((entry) => entry.supported)
                                      .map((entry) => (
                                        <option key={entry.value} value={entry.value}>
                                          {entry.value}
                                        </option>
                                      )),
                                  )}
                                </select>
                              </label>
                              <label className="space-y-1 text-[11px] uppercase tracking-[0.16em] text-slate-400">
                                <span>{t('workbench.source.planner.dataFormat')}</span>
                                <select
                                  aria-label={t('workbench.source.planner.dataFormat')}
                                  value={editDraft.dataFormat}
                                  onChange={(event) =>
                                    setEditDraft((draft) =>
                                      draft
                                        ? { ...draft, dataFormat: event.target.value }
                                        : draft,
                                    )
                                  }
                                  className="w-full rounded-lg border border-slate-800 bg-slate-900 px-2 py-1 text-[11px] text-slate-100"
                                >
                                  <option value="">{t('workbench.source.planner.dataFormatDefault')}</option>
                                  <option value="CDAB">CDAB</option>
                                  <option value="ABCD">ABCD</option>
                                  <option value="BADC">BADC</option>
                                  <option value="DCBA">DCBA</option>
                                </select>
                              </label>
                              <label className="space-y-1 text-[11px] uppercase tracking-[0.16em] text-slate-400">
                                <span>{t('workbench.source.planner.scaleOffset')}</span>
                                <input
                                  aria-label={t('workbench.source.planner.scaleOffset')}
                                  type="number"
                                  step="any"
                                  placeholder="0.0"
                                  value={editDraft.scaleOffset}
                                  onChange={(event) =>
                                    setEditDraft((draft) =>
                                      draft
                                        ? { ...draft, scaleOffset: event.target.value }
                                        : draft,
                                    )
                                  }
                                  className="w-full rounded-lg border border-slate-800 bg-slate-900 px-2 py-1 text-[11px] text-slate-100"
                                />
                              </label>
                              <label className="space-y-1 text-[11px] uppercase tracking-[0.16em] text-slate-400">
                                <span>{t('workbench.source.planner.scaleMultiplier')}</span>
                                <input
                                  aria-label={t('workbench.source.planner.scaleMultiplier')}
                                  type="number"
                                  step="any"
                                  placeholder="1.0"
                                  value={editDraft.scaleMultiplier}
                                  onChange={(event) =>
                                    setEditDraft((draft) =>
                                      draft
                                        ? { ...draft, scaleMultiplier: event.target.value }
                                        : draft,
                                    )
                                  }
                                  className="w-full rounded-lg border border-slate-800 bg-slate-900 px-2 py-1 text-[11px] text-slate-100"
                                />
                              </label>
                            </div>
                          </div>
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

                      {!isRuleCardCollapsed ? (
                      <div className="grid grid-cols-2 gap-2 text-[11px]">
                        <button
                          className="rounded-lg border border-slate-700/70 px-2 py-1 text-slate-300 disabled:cursor-not-allowed disabled:opacity-50"
                          disabled={isRuleEnableMutating}
                          onClick={() => void handleToggleRuleEnabled(rule.id)}
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
                            ? t('workbench.source.ruleLayer.unprotectPlan')
                            : t('workbench.source.ruleLayer.protectPlan')}
                        </button>
                        {rule.locked ? (
                          <p
                            className="col-span-2 text-[10px] text-amber-200/80"
                            data-testid={`rule-protect-hint-${rule.id}`}
                          >
                            {t('workbench.source.ruleLayer.protectHint')}
                          </p>
                        ) : null}
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
                      ) : null}
                    </article>
                  );
                })}
              </div>
            )}
            </div>
          </section>
        </aside>

        <section
          className="flex min-h-0 min-w-0 flex-1 flex-col gap-4 overflow-hidden rounded-2xl border border-slate-800 bg-slate-950/40 p-4"
          data-emphasis="primary"
          data-testid="source-canvas-workspace"
        >
          <div
            className="flex shrink-0 flex-col gap-3 xl:flex-row xl:items-start xl:justify-between"
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

          <div className="flex min-h-0 min-w-0 flex-1 flex-col gap-4 overflow-hidden">
            <div
              aria-label={t('workbench.source.canvas.scrollRegion')}
              className="min-h-0 flex-1 overflow-y-auto overscroll-contain rounded-xl border border-slate-800/50 bg-slate-950/30 px-2 py-3 scrollbar-auto-hide"
              data-testid="source-memory-scroll-region"
            >
              <AddressCanvas
                items={items}
                onSelectAddress={handleSelectAddress}
                selectedAddress={selectedAddress}
                valueFormat={valueFormat}
                viewMode={viewMode}
              />
            </div>

            <div
              ref={sourceWorkspaceSecondaryScrollRef}
              className="min-h-0 max-h-[min(46vh,26rem)] touch-pan-y space-y-4 overflow-y-auto overscroll-contain pb-4 scrollbar-none"
            >
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
                <p className="text-xs font-semibold uppercase tracking-[0.18em] text-cyan-300">
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

              <div className="flex flex-wrap items-center gap-2">
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
                    safePointDefinitions.length === 0 ||
                    createPointMutation.isPending
                  }
                  className="rounded-lg bg-cyan-500 px-3 py-2 text-sm font-semibold text-slate-950 disabled:cursor-not-allowed disabled:bg-slate-800 disabled:text-slate-500"
                >
                  {t('workbench.source.actions.createRulePoints')}
                  {conflictCount > 0 && safePointDefinitions.length > 0 ? (
                    <span className="ml-1 text-xs font-normal opacity-80">
                      ({safePointDefinitions.length})
                    </span>
                  ) : null}
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
              <label className="flex items-end gap-1 text-xs uppercase tracking-[0.16em] text-slate-500">
                <input
                  aria-label={t('workbench.source.toolbar.jumpToAddress')}
                  value={jumpAddress}
                  onChange={(event) => setJumpAddress(event.target.value)}
                  placeholder={t('workbench.source.toolbar.jumpToAddress')}
                  className="w-28 rounded-lg border border-slate-800 bg-slate-900 px-2 py-2 text-sm text-slate-100"
                />
                <button
                  type="button"
                  onClick={handleJumpToAddress}
                  className="rounded-lg border border-slate-800 bg-slate-900/60 px-3 py-2 text-xs text-slate-300"
                >
                  {t('workbench.source.toolbar.jump')}
                </button>
              </label>
            </div>
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

          {conflictQueue.length > 0 ? (
            <section
              className="space-y-2 rounded-xl border border-rose-500/30 bg-rose-900/10 p-3"
              data-testid="source-conflict-queue"
            >
              <p className="text-xs font-semibold uppercase tracking-[0.18em] text-rose-300">
                {t('workbench.source.conflictQueue.title')}
              </p>
              <p className="text-xs text-rose-200">
                {t('workbench.source.conflictQueue.step3Blocked')}
              </p>
              <div className="space-y-2">
                {conflictQueue.map((conflict) => {
                  const targetRuleId = conflict.ruleIds[conflict.ruleIds.length - 1];
                  return (
                    <article
                      key={conflict.id}
                      className="flex flex-wrap items-center gap-2 rounded-lg border border-rose-500/20 bg-rose-950/40 px-3 py-2"
                      data-testid={`conflict-item-${conflict.address}`}
                    >
                      <span className="font-mono text-xs text-rose-100">
                        {conflict.address}
                      </span>
                      <span className="text-xs text-rose-200">
                        {t(conflict.reasonKey)}
                      </span>
                      <span className="text-xs text-rose-300">
                        {conflict.ruleIds.join(', ')}
                      </span>
                      <div className="ml-auto flex gap-1">
                        {targetRuleId ? (
                          <>
                            <button
                              type="button"
                              className="rounded-lg border border-rose-500/30 px-2 py-1 text-[11px] text-rose-100"
                              onClick={() => {
                                handleSelectRule(targetRuleId);
                                handleStartRuleEdit(targetRuleId);
                              }}
                            >
                              {t('workbench.source.conflictQueue.editRule')}
                            </button>
                            <button
                              type="button"
                              className="rounded-lg border border-rose-500/30 px-2 py-1 text-[11px] text-rose-100"
                              onClick={() => handleSkipConflictSpan(targetRuleId, conflict.conflictCellAddress)}
                            >
                              {t('workbench.source.conflictQueue.skipSpan')}
                            </button>
                          </>
                        ) : null}
                      </div>
                    </article>
                  );
                })}
              </div>
            </section>
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
            </div>
          </div>
        </section>
      </div>
    </section>
  );
}
