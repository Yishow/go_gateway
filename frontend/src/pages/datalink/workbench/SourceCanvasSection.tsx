import { ArrowDown, ArrowUp, ChevronDown, ChevronRight, CircleHelp } from 'lucide-react';
import { useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import {
  applyTemplateToPlanner,
  createTemplateFromPlanner,
  getDefaultNamingPrefixForProtocol,
  normalizeNamingPrefix,
  parsePlannerCountInput,
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
import { useDevicesQuery } from '../../../hooks/datalink/useDevices';
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
import { SourceCanvasStatusState } from './SourceCanvasStatusState';
import {
  clampLatticeColumns,
  LATTICE_COLUMNS_DEFAULT,
  LATTICE_COLUMNS_MAX,
  LATTICE_COLUMNS_MIN,
  readSourceCanvasLatticeColumns,
  writeSourceCanvasLatticeColumns,
} from './sourceCanvasLatticeColumns';
import { AddressLedger } from './AddressLedger';
import { SourceRuleLayerPanel } from './SourceRuleLayerPanel';
import { SourceStepRuleSummary } from './SourceStepRuleSummary';
import { SourceTriagePanel } from './SourceTriagePanel';
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
import { buildSourcePlanningDatabaseAdvisory } from './sourceStepRuleSummaryModel';

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

/**
 * 規則建立器／進階設定表單列：左欄位標籤、右輸入控制項（適用窄側欄）。
 */
const SOURCE_PLANNER_INLINE_LABEL_ROW_CLASS =
  'grid grid-cols-[minmax(0,38%)_minmax(0,1fr)] items-center gap-2 text-sm uppercase tracking-[0.16em] text-slate-400';

/**
 * 規則卡片內嵌編輯表單列：同 {@link SOURCE_PLANNER_INLINE_LABEL_ROW_CLASS}，字級較小以配合卡片密度。
 */
const SOURCE_PLANNER_INLINE_LABEL_ROW_CLASS_COMPACT =
  'grid grid-cols-[minmax(0,38%)_minmax(0,1fr)] items-center gap-2 text-[11px] uppercase tracking-[0.16em] text-slate-400';

/**
 * 規則建立器控制項外觀：與畫布工具列「每列格數」輸入框對齊之 padding、字級與圓角，使高度一致。
 */
const SOURCE_PLANNER_FIELD_CONTROL_CLASS =
  'min-w-0 w-full rounded-md border border-slate-700 bg-slate-950 px-2 py-1.5 text-xs text-slate-100';

/**
 * 位址／數值欄位用等寬字體；其餘同 {@link SOURCE_PLANNER_FIELD_CONTROL_CLASS}。
 */
const SOURCE_PLANNER_FIELD_CONTROL_MONO_CLASS = `${SOURCE_PLANNER_FIELD_CONTROL_CLASS} font-mono`;

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

type SourceCanvasSectionProps = {
  deskMode: 'inspect' | 'build' | 'triage';
};

export function SourceCanvasSection({ deskMode }: SourceCanvasSectionProps) {
  const { t } = useTranslation();
  const {
    activeOutputTarget,
    crossStepContext,
    selectedDeviceId,
    setActiveStep,
    setFocusedRuleId,
    setInspectorSelection,
    setSelectedDeviceId,
    setSourceStepNotice,
    setSourceStepInspectorBanner,
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
  const selectedDeviceProtocol = selectedDevice?.protocol;

  const [viewMode, setViewMode] = useState<SourceViewMode>('plan');
  const [valueFormat, setValueFormat] = useState<SourceValueFormat>('decimal');
  const [showAudit, setShowAudit] = useState(false);
  const [freezeLive, setFreezeLive] = useState(false);
  const [frozenLiveValues, setFrozenLiveValues] = useState<
    Readonly<Record<string, { raw_value?: unknown; timestamp?: string }>> | null
  >(null);
  const [startAddress, setStartAddress] = useState('40001');
  const [dataType, setDataType] = useState<DataType>('int16');
  const [count, setCount] = useState('4');
  const [namingPrefix, setNamingPrefix] = useState('');
  const [targetDataType, setTargetDataType] = useState<DataType | ''>('');
  const [scaleMultiplier, setScaleMultiplier] = useState<string>('');
  const [scaleOffset, setScaleOffset] = useState<string>('');
  const [plannerFieldTouched, setPlannerFieldTouched] = useState({
    namingPrefix: false,
    targetDataType: false,
    scaleMultiplier: false,
    scaleOffset: false,
  });
  const [hoverPreviewRuleId, setHoverPreviewRuleId] = useState<string | null>(null);
  const [dataFormat, setDataFormat] = useState<string>('');
  /** 規劃規則層側欄：規則建立器與已新增規則清單之分頁。 */
  const [ruleLayerTab, setRuleLayerTab] = useState<'planner' | 'rules'>('planner');
  /** 右欄「覆蓋總覽」區塊是否展開（收合時僅保留標題列；預設收合）。 */
  const [coverageOverviewOpen, setCoverageOverviewOpen] = useState(false);
  /** 已收合的規則卡片 id（預設空＝全部展開；使用者可收合以精簡側欄）。 */
  const [ruleCardCollapsedIds, setRuleCardCollapsedIds] = useState(() => new Set<string>());
  const [editingRuleId, setEditingRuleId] = useState<string | null>(null);
  const [editDraft, setEditDraft] = useState<{
    startAddress: string;
    count: string;
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
  /** 主工具列「更多」選單：收合凍結／稽核／模板等次要操作。 */
  const [sourceToolbarMoreOpen, setSourceToolbarMoreOpen] = useState(false);
  const sourceToolbarMoreRef = useRef<HTMLDivElement>(null);
  /** 來源畫布每列顯示幾個位址格（例 10 → 40001～40010 同一列）；自本機 storage 還原。 */
  const [canvasLatticeColumns, setCanvasLatticeColumns] = useState(LATTICE_COLUMNS_DEFAULT);
  const [templateNotice, setTemplateNotice] = useState<string | null>(null);
  const [templateRecoveryWarning, setTemplateRecoveryWarning] = useState<string | null>(null);
  const [templateWarning, setTemplateWarning] = useState<string | null>(null);
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
      setTemplateRecoveryWarning(null);
      setTemplates([]);
      return;
    }

    const hasStaleTemplate = restored.some(isTemplateStale);
    const nextTemplates = hasStaleTemplate
      ? sortTemplates(upgradeTemplates(restored))
      : sortTemplates(restored);

    if (hasStaleTemplate) {
      saveSourceTemplates(nextTemplates);
      setTemplateRecoveryWarning(
        t('workbench.source.templates.recoveredLegacy', {
          defaultValue:
            'Legacy templates were auto-upgraded. Review the range before reusing it.',
        }),
      );
    } else {
      setTemplateRecoveryWarning(null);
    }
    setTemplates(nextTemplates);
    // eslint-disable-next-line react-hooks/exhaustive-deps -- legacy template recovery only needs mount-time evaluation
  }, []);
  useEffect(() => {
    setTemplateWarning(null);
    setAppliedTemplate(null);
  }, [selectedDeviceId]);
  /** 在捲動區頂/底邊界攔截 wheel，避免位移傳到外層造成整塊版面跟動。 */
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
  /** 「更多」選單開啟時：點擊外區或 Escape 關閉，避免遮擋下層互動。 */
  useEffect(() => {
    if (!sourceToolbarMoreOpen) {
      return undefined;
    }

    const onPointerDown = (event: MouseEvent) => {
      if (
        sourceToolbarMoreRef.current &&
        !sourceToolbarMoreRef.current.contains(event.target as Node)
      ) {
        setSourceToolbarMoreOpen(false);
      }
    };
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setSourceToolbarMoreOpen(false);
      }
    };

    document.addEventListener('mousedown', onPointerDown);
    document.addEventListener('keydown', onKeyDown);
    return () => {
      document.removeEventListener('mousedown', onPointerDown);
      document.removeEventListener('keydown', onKeyDown);
    };
  }, [sourceToolbarMoreOpen]);

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
    setPlannerFieldTouched({
      namingPrefix: false,
      targetDataType: false,
      scaleMultiplier: false,
      scaleOffset: false,
    });
    setHoverPreviewRuleId(null);
  }, [selectedDeviceId]);

  /**
   * 切換設備（或設備資料載入完成）時，依協定重設規則建立器的命名前綴預設值。
   * 僅依 `deviceId` 與 `protocol` 觸發，避免 React Query 重取導致同一台設備下使用者已編輯的前綴被洗回。
   */
  useEffect(() => {
    if (!selectedDeviceId || !selectedDeviceProtocol) {
      return;
    }

    if (activeOutputTarget === 'database') {
      const advisory = buildSourcePlanningDatabaseAdvisory(selectedDeviceProtocol, dataType);
      if (!plannerFieldTouched.namingPrefix) {
        setNamingPrefix(advisory.namingPrefix);
      }
      if (!plannerFieldTouched.targetDataType) {
        setTargetDataType(advisory.targetDataType);
      }
      if (!plannerFieldTouched.scaleMultiplier) {
        setScaleMultiplier(advisory.scaleMultiplier);
      }
      if (!plannerFieldTouched.scaleOffset) {
        setScaleOffset(advisory.scaleOffset);
      }
      return;
    }

    if (!plannerFieldTouched.namingPrefix) {
      setNamingPrefix(getDefaultNamingPrefixForProtocol(selectedDeviceProtocol));
    }
    if (!plannerFieldTouched.targetDataType) {
      setTargetDataType('');
    }
    if (!plannerFieldTouched.scaleMultiplier) {
      setScaleMultiplier('');
    }
    if (!plannerFieldTouched.scaleOffset) {
      setScaleOffset('');
    }
  }, [
    activeOutputTarget,
    dataType,
    plannerFieldTouched.namingPrefix,
    plannerFieldTouched.scaleMultiplier,
    plannerFieldTouched.scaleOffset,
    plannerFieldTouched.targetDataType,
    selectedDeviceId,
    selectedDeviceProtocol,
  ]);

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

  useEffect(() => {
    setCanvasLatticeColumns(readSourceCanvasLatticeColumns());
  }, []);

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

  const runtimePointsHealthy = runtimeStatusQuery.data?.collectors[0]?.points_healthy ?? 0;
  const runtimePointsStale = runtimeStatusQuery.data?.collectors[0]?.points_stale ?? 0;
  const activeRuleId = useMemo(() => {
    const preferredIds = [crossStepContext.focusedRuleId, sourcePlanningState.selectedRuleId];
    for (const candidateId of preferredIds) {
      if (candidateId && rules.some((rule) => rule.id === candidateId)) {
        return candidateId;
      }
    }
    return rules[0]?.id ?? null;
  }, [
    crossStepContext.focusedRuleId,
    rules,
    sourcePlanningState.selectedRuleId,
  ]);
  const previewRuleId = hoverPreviewRuleId ?? activeRuleId;

  useEffect(() => {
    setSourceStepInspectorBanner(null);
    return () => {
      setSourceStepInspectorBanner(null);
    };
  }, [setSourceStepInspectorBanner]);

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
      /** `expand` 解析失敗時回傳空陣列；`[].every(...)` 會誤判為安全，須拒絕。 */
      if (occupied.length !== span) {
        return false;
      }
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

    const parsedCount = parsePlannerCountInput(count);
    if (parsedCount === null) {
      return;
    }

    const ruleId = getNextRuleId(rules);
    const nextRule: SourceRule = {
      id: ruleId,
      deviceId: selectedDevice.id,
      startAddress: startAddress.trim(),
      count: parsedCount,
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
    setHoverPreviewRuleId(null);
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
    setRuleCardCollapsedIds((prev) => {
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

  const handleSkipSelection = () => {
    setSourcePlanningState((currentState) => ({
      ...currentState,
      selectedAddress: null,
    }));
    setInspectorSelection({ kind: 'none' });
  };

  /**
   * 更新每列格數並寫入本機；空字或無效時還原為預設 16。
   */
  const handleCanvasLatticeColumnsChange = (raw: string) => {
    if (raw === '') {
      setCanvasLatticeColumns(LATTICE_COLUMNS_DEFAULT);
      writeSourceCanvasLatticeColumns(LATTICE_COLUMNS_DEFAULT);
      return;
    }
    const parsed = parseInt(raw, 10);
    if (!Number.isFinite(parsed)) {
      return;
    }
    const next = clampLatticeColumns(parsed);
    setCanvasLatticeColumns(next);
    writeSourceCanvasLatticeColumns(next);
  };

  const handleStartRuleEdit = (ruleId: string) => {
    const rule = rules.find((candidate) => candidate.id === ruleId);
    if (!rule) {
      return;
    }

    setRuleCardCollapsedIds((prev) => {
      const next = new Set(prev);
      next.delete(ruleId);
      return next;
    });

    setEditingRuleId(ruleId);
    setEditDraft({
      startAddress: rule.startAddress,
      count: String(rule.count),
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

    const parsedEditCount = parsePlannerCountInput(editDraft.count);
    if (parsedEditCount === null) {
      return;
    }

    const nextStartAddress = editDraft.startAddress.trim();
    const geometryChanged =
      rule.startAddress !== nextStartAddress ||
      rule.count !== parsedEditCount ||
      rule.dataType !== editDraft.dataType;
    const nextSkippedAddresses = geometryChanged ? [] : editDraft.skippedAddresses;

    if (rule.persisted) {
      await updateSourceRuleMutation.mutateAsync({
        id: ruleId,
        data: {
          start_address: nextStartAddress,
          count: parsedEditCount,
          data_type: editDraft.dataType,
          skipped_addresses: nextSkippedAddresses,
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
                  startAddress: nextStartAddress,
                  count: parsedEditCount,
                  dataType: editDraft.dataType,
                  skippedAddresses: nextSkippedAddresses,
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

    const parsedTemplateCount = parsePlannerCountInput(count);
    if (parsedTemplateCount === null) {
      setTemplateNotice(t('workbench.source.planner.countInvalid'));
      return;
    }

    const nextTemplate = createTemplateFromPlanner({
      templateName: trimmedName,
      draft: {
        startAddress,
        count: parsedTemplateCount,
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
    setCount(String(plannerDraft.count));
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

  /**
   * 批次「套用規則」：將畫布上無衝突、可安全建立的規劃根位址送交後端。
   * - 草稿規則（`persisted === false`）：呼叫建立來源規則 API（後端一併建立 point／連結）。
   * - 已持久化規則：改以更新 API 帶入合併後的 `skipped_addresses`，由後端 reconcile 並補齊缺漏 point
   *   （例如 point 曾被刪除、或僅有規則無連結之資料狀態）。
   *
   * @returns Promise<void>
   */
  const handleBatchCreate = async () => {
    if (!selectedDevice || safePointDefinitions.length === 0) {
      return;
    }

    const safeAddressSet = new Set(safePointDefinitions.map((definition) => definition.address));
    /** 與 `plannedPointDefinitions`（僅 `locked === false`）對齊，避免對保護規則送 API。 */
    const eligibleRules = rules
      .filter((rule) => rule.enabled && !rule.locked)
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
      setSourceStepNotice(t('workbench.source.planner.batchApplyNoEligibleRules'));
      return;
    }

    const results = await Promise.allSettled(
      eligibleRules.map(({ rule, skippedAddresses }) =>
        rule.persisted
          ? updateSourceRuleMutation.mutateAsync({
              id: rule.id,
              data: { skipped_addresses: skippedAddresses },
            })
          : // 建立請求不帶 `id`，由後端產生 UUID，避免與既有 `source_rules.id`（如舊版 `rule-1`）UNIQUE 衝突。
            createSourceRuleMutation.mutateAsync({
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

    const succeededDraftRuleIds = eligibleRules
      .filter((_, index) => results[index]?.status === 'fulfilled')
      .filter(({ rule }) => !rule.persisted)
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
          currentState.selectedRuleId &&
          succeededDraftRuleIds.includes(currentState.selectedRuleId)
            ? null
            : currentState.selectedRuleId,
        selectedAddress: null,
      }));
    }

    const successCount = results.filter((result) => result.status === 'fulfilled').length;
    const failureCount = results.filter((result) => result.status === 'rejected').length;

    setBatchCreateSummary({ successCount, failureCount });

    if (failureCount > 0) {
      let rejectReason: unknown;
      for (const result of results) {
        if (result.status === 'rejected') {
          rejectReason = result.reason;
          break;
        }
      }
      setSourceStepNotice(
        getErrorMessage(rejectReason, t('workbench.source.planner.batchApplyErrorFallback')),
      );
    } else {
      setSourceStepNotice(null);
    }
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
          <p className="text-sm font-semibold uppercase tracking-[0.18em] text-cyan-300">
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

  if (selectedDeviceId && sourceRulesQuery.isLoading) return <SourceCanvasStatusState state="loading" />;
  if (sourceRulesQuery.isError)
    return (
      <SourceCanvasStatusState
        error={sourceRulesQuery.error}
        onRetry={() => {
          void sourceRulesQuery.refetch();
        }}
        state="error"
      />
    );

  const workspaceModeTestId =
    deskMode === 'build' ? 'source-build-workspace' : 'source-inspect-workspace';
  const unmanagedIssueCount = items.filter(
    (item) => item.status === 'unmanaged' && item.mergeOffset === 0,
  ).length;
  const triageNeedsAttention = conflictQueue.length > 0 || unmanagedIssueCount > 0;
  const workspaceSummaryStrip = (
    <div
      className="grid gap-4 xl:grid-cols-[minmax(0,1.35fr)_minmax(18rem,0.9fr)]"
      data-testid="source-workspace-summary-strip"
    >
      <SourceStepRuleSummary
        activeOutputTarget={activeOutputTarget}
        activeRuleId={activeRuleId}
        applyDisabled={
          safePointDefinitions.length === 0 ||
          createPointMutation.isPending ||
          createSourceRuleMutation.isPending ||
          updateSourceRuleMutation.isPending
        }
        batchCreateSummary={batchCreateSummary}
        items={items}
        onApplyPlanning={() => {
          void handleBatchCreate();
        }}
        onHoverRule={setHoverPreviewRuleId}
        onSelectRule={handleSelectRule}
        previewRuleId={previewRuleId}
        protocol={selectedDevice.protocol}
        rules={rules}
      />
      <div className="min-w-0 space-y-3.5">
        <section className="space-y-2" data-testid="source-runtime-collection-panel">
          <div className="flex min-w-0 items-center justify-between gap-2">
            <p className="m-0 text-[10px] font-semibold uppercase tracking-[0.2em] text-cyan-400/95">
              {t('workbench.runtime.summary.title')}
            </p>
            <button
              type="button"
              data-testid="source-runtime-collection-hint"
              className="shrink-0 rounded-lg p-1 text-slate-500 transition-colors hover:bg-slate-800/90 hover:text-slate-200 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-cyan-400/80"
              aria-label={t('workbench.source.collection.hintHint')}
              title={t('workbench.source.collection.hint')}
            >
              <CircleHelp className="h-3.5 w-3.5" aria-hidden />
            </button>
          </div>
          <div className="grid grid-cols-2 gap-2">
            <article className="min-w-0 rounded-xl border border-emerald-500/20 bg-emerald-500/[0.08] px-2.5 py-2 shadow-sm shadow-black/20">
              <p className="text-[9px] font-semibold uppercase leading-tight tracking-[0.14em] text-emerald-200/75">
                {t('workbench.runtime.summary.pointsHealthy')}
              </p>
              <p className="mt-1 text-lg font-semibold tabular-nums tracking-tight text-emerald-100">
                {runtimePointsHealthy}
              </p>
            </article>
            <article className="min-w-0 rounded-xl border border-amber-500/20 bg-amber-500/[0.07] px-2.5 py-2 shadow-sm shadow-black/20">
              <p className="text-[9px] font-semibold uppercase leading-tight tracking-[0.14em] text-amber-200/75">
                {t('workbench.runtime.summary.pointsStale')}
              </p>
              <p className="mt-1 text-lg font-semibold tabular-nums tracking-tight text-amber-100">
                {runtimePointsStale}
              </p>
            </article>
          </div>
        </section>
        <section aria-label={t('workbench.source.summary.title')}>
          <p className="mb-2 text-[10px] font-semibold uppercase tracking-[0.18em] text-slate-400">
            {t('workbench.source.summary.title')}
          </p>
          <div className="grid min-w-0 cursor-help grid-cols-3 gap-1.5" title={t('workbench.source.summary.panelHint')}>
            <div
              className="flex min-w-0 flex-col rounded-lg border border-slate-700/60 bg-slate-900/50 px-1.5 py-2 shadow-sm shadow-black/15"
              title={t('workbench.source.summary.readyToCreateHint')}
            >
              <p className="truncate text-[8px] font-semibold uppercase tracking-[0.12em] text-slate-500">
                {t('workbench.source.summary.readyToCreate')}
              </p>
              <p
                className="mt-1 text-base font-semibold tabular-nums leading-none text-emerald-300/95"
                data-testid="source-workspace-ready-count"
              >
                <span data-testid="source-summary-ready-count">{readyToCreateCount}</span>
              </p>
            </div>
            <div
              className="flex min-w-0 flex-col rounded-lg border border-slate-700/60 bg-slate-900/50 px-1.5 py-2 shadow-sm shadow-black/15"
              title={t('workbench.source.summary.inConflictHint')}
            >
              <p className="truncate text-[8px] font-semibold uppercase tracking-[0.12em] text-slate-500">
                {t('workbench.source.summary.inConflict')}
              </p>
              <p
                className="mt-1 text-base font-semibold tabular-nums leading-none text-rose-300/95"
                data-testid="source-workspace-conflict-count"
              >
                <span data-testid="source-summary-conflict-count">{conflictCount}</span>
              </p>
            </div>
            <div
              className="flex min-w-0 flex-col rounded-lg border border-slate-700/60 bg-slate-900/50 px-1.5 py-2 shadow-sm shadow-black/15"
              title={t('workbench.source.summary.protectedHint')}
            >
              <p className="truncate text-[8px] font-semibold uppercase tracking-[0.12em] text-slate-500">
                {t('workbench.source.summary.protected')}
              </p>
              <p
                className="mt-1 text-base font-semibold tabular-nums leading-none text-amber-300/95"
                data-testid="source-workspace-protected-count"
              >
                <span data-testid="source-summary-protected-count">
                  {protectedPointDefinitions.length}
                </span>
              </p>
            </div>
          </div>
        </section>
      </div>
    </div>
  );
  const workspaceDiagnosticsStrip = (
    <section
      className="rounded-xl border border-slate-800/60 bg-slate-950/20 px-4 py-3"
      data-testid="source-workspace-diagnostics-strip"
    >
      <div className="flex flex-col gap-1 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <p className="text-[10px] font-semibold uppercase tracking-[0.18em] text-slate-500">
            {t('workbench.source.diagnostics.eyebrow')}
          </p>
          <p className="mt-1 text-sm text-slate-300">
            {t('workbench.source.diagnostics.description')}
          </p>
          <p className="mt-2 text-xs text-cyan-200">
            {t(`workbench.source.deskModeGuidance.${deskMode}`)}
          </p>
        </div>
        <p className="text-xs font-semibold uppercase tracking-[0.18em] text-cyan-300">
          {t(`workbench.source.deskMode.${deskMode}`)}
        </p>
      </div>
    </section>
  );
  const triageRecoveryPanel = (
    <section
      className="space-y-3 rounded-xl border border-rose-500/20 bg-rose-950/10 p-4"
      data-testid="source-triage-recovery-panel"
    >
      <div className="space-y-1">
        <p className="text-[10px] font-semibold uppercase tracking-[0.18em] text-rose-300">
          {t('workbench.source.triage.eyebrow')}
        </p>
        <h3 className="text-sm font-semibold text-slate-100">
          {triageNeedsAttention
            ? t('workbench.source.triage.title')
            : t('workbench.source.triage.clearTitle')}
        </h3>
        <p className="text-sm text-slate-300">
          {triageNeedsAttention
            ? t('workbench.source.triage.description', { count: conflictQueue.length + unmanagedIssueCount })
            : t('workbench.source.triage.clearDescription')}
        </p>
      </div>
      {triageNeedsAttention ? (
        <>
          {unmanagedIssueCount > 0 ? (
            <div className="rounded-lg border border-amber-500/20 bg-amber-500/[0.08] px-3 py-2 text-sm text-amber-100">
              {t('workbench.source.triage.unmanagedSummary', { count: unmanagedIssueCount })}
            </div>
          ) : null}
          {deskMode !== 'triage' && conflictQueue.length > 0 ? (
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
                                setRuleLayerTab('rules');
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
        </>
      ) : (
        <div
          className="rounded-lg border border-emerald-500/20 bg-emerald-500/[0.08] px-3 py-3 text-sm text-emerald-100"
          data-testid="source-triage-clear-state"
        >
          {t('workbench.source.triage.clearState')}
          <div className="mt-3 flex flex-wrap gap-2">
            <button
              type="button"
              className="rounded-lg border border-emerald-400/30 px-3 py-1.5 text-xs font-medium text-emerald-100"
              onClick={() => setActiveStep('tag')}
            >
              {t('workbench.source.handoff.toTag')}
            </button>
          </div>
        </div>
      )}
    </section>
  );

  const workspaceGrid = (
    <div
      className={[
        'grid min-h-0 flex-1 gap-4 overflow-hidden',
        deskMode === 'build'
          ? 'xl:grid-cols-[minmax(20rem,0.92fr)_minmax(0,1.08fr)]'
          : 'xl:grid-cols-[300px_minmax(0,1fr)]',
      ].join(' ')}
    >
      <SourceRuleLayerPanel
        panelTestId={deskMode === 'build' ? 'source-build-rule-panel' : 'source-rule-layer-panel'}
      >
          <aside className="flex min-h-0 flex-col gap-4 overflow-hidden">
          <section
            className="flex min-h-0 flex-1 flex-col gap-3 overflow-hidden rounded-2xl border border-slate-800/70 bg-slate-950/25 p-4"
            data-emphasis="primary"
            data-testid="source-rule-layer"
          >
            <div className="shrink-0 space-y-1">
              <p className="text-sm font-semibold uppercase tracking-[0.18em] text-cyan-300">
                {t('workbench.source.ruleLayer.eyebrow')}
              </p>
              <div className="flex flex-wrap items-center gap-1.5">
                <h3 className="m-0 min-w-0 text-sm font-semibold leading-snug text-slate-100">
                  {t('workbench.source.ruleLayer.title')}
                </h3>
                <button
                  type="button"
                  data-testid="source-rule-layer-description-hint"
                  className="shrink-0 rounded-md p-0.5 text-slate-400 transition-colors hover:bg-slate-800/80 hover:text-slate-200 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-cyan-400"
                  aria-label={t('workbench.source.ruleLayer.descriptionHint')}
                  title={t('workbench.source.ruleLayer.description')}
                >
                  <CircleHelp className="h-4 w-4" aria-hidden />
                </button>
              </div>
            </div>

            <div className="flex min-h-0 flex-1 flex-col gap-3 overflow-hidden">
              <div
                className="flex shrink-0 items-stretch gap-1 rounded-xl border border-slate-800/70 bg-slate-950/70 p-1"
                role="tablist"
                aria-label={t('workbench.source.ruleLayer.tabListAria')}
              >
                <button
                  type="button"
                  role="tab"
                  id="source-rule-layer-tab-planner"
                  data-testid="source-rule-layer-tab-planner"
                  aria-selected={ruleLayerTab === 'planner'}
                  aria-controls="source-rule-layer-panel-planner"
                  onClick={() => setRuleLayerTab('planner')}
                  className={
                    ruleLayerTab === 'planner'
                      ? 'flex min-h-10 min-w-0 flex-1 items-center justify-center gap-1 rounded-lg bg-cyan-500 px-2 py-2 text-center text-sm font-semibold text-slate-950'
                      : 'flex min-h-10 min-w-0 flex-1 items-center justify-center gap-1 rounded-lg border border-transparent px-2 py-2 text-center text-sm font-medium text-slate-500 transition hover:bg-slate-800/45 hover:text-slate-300'
                  }
                >
                  <span className="min-w-0 truncate">{t('workbench.source.planner.title')}</span>
                </button>
                <button
                  type="button"
                  role="tab"
                  id="source-rule-layer-tab-rules"
                  data-testid="source-rule-layer-tab-rules"
                  aria-selected={ruleLayerTab === 'rules'}
                  aria-controls="source-rule-layer-panel-rules"
                  aria-label={
                    rules.length > 0
                      ? `${t('workbench.source.ruleLayer.tabRules')} (${rules.length})`
                      : t('workbench.source.ruleLayer.tabRules')
                  }
                  onClick={() => setRuleLayerTab('rules')}
                  className={
                    ruleLayerTab === 'rules'
                      ? 'flex min-h-10 min-w-0 flex-1 items-center justify-center gap-1 rounded-lg bg-cyan-500 px-2 py-2 text-center text-sm font-semibold text-slate-950'
                      : 'flex min-h-10 min-w-0 flex-1 items-center justify-center gap-1 rounded-lg border border-transparent px-2 py-2 text-center text-sm font-medium text-slate-500 transition hover:bg-slate-800/45 hover:text-slate-300'
                  }
                >
                  <span className="min-w-0 truncate">{t('workbench.source.ruleLayer.tabRules')}</span>
                  {rules.length > 0 ? (
                    <span
                      className={
                        ruleLayerTab === 'rules'
                          ? 'shrink-0 rounded-md bg-slate-950/20 px-1.5 py-0.5 text-[10px] font-semibold tabular-nums text-slate-950'
                          : 'shrink-0 rounded-md bg-slate-800/55 px-1.5 py-0.5 text-[10px] font-semibold tabular-nums text-slate-400'
                      }
                      aria-hidden
                    >
                      {rules.length}
                    </span>
                  ) : null}
                </button>
              </div>

              <div className="relative min-h-0 flex-1 overflow-hidden">
              <div
                role="tabpanel"
                id="source-rule-layer-panel-planner"
                aria-labelledby="source-rule-layer-tab-planner"
                hidden={ruleLayerTab !== 'planner'}
                className="absolute inset-0 overflow-y-auto overscroll-contain scrollbar-auto-hide"
              >
                  <div className="space-y-3 rounded-xl border border-slate-800/70 bg-slate-950/70 p-3">
              <div className="flex flex-col gap-3">
                <label className={SOURCE_PLANNER_INLINE_LABEL_ROW_CLASS}>
                  <span className="min-w-0 break-words leading-snug">
                    {t('workbench.source.planner.startAddress')}
                  </span>
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
                    className={SOURCE_PLANNER_FIELD_CONTROL_MONO_CLASS}
                  />
                </label>
                <label className={SOURCE_PLANNER_INLINE_LABEL_ROW_CLASS}>
                  <span className="min-w-0 break-words leading-snug">{t('workbench.source.planner.count')}</span>
                  <input
                    aria-label={t('workbench.source.planner.count')}
                    onChange={(event) => {
                      clearAppliedTemplate();
                      setCount(event.target.value);
                    }}
                    type="number"
                    value={count}
                    className={SOURCE_PLANNER_FIELD_CONTROL_MONO_CLASS}
                  />
                </label>
                <label className={SOURCE_PLANNER_INLINE_LABEL_ROW_CLASS}>
                  <span className="min-w-0 break-words leading-snug">{t('workbench.source.planner.dataType')}</span>
                  <select
                    aria-label={t('workbench.source.planner.dataType')}
                    value={dataType}
                    onChange={(event) => {
                      clearAppliedTemplate();
                      setDataType(event.target.value as DataType);
                    }}
                    className={SOURCE_PLANNER_FIELD_CONTROL_CLASS}
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
                <label className={SOURCE_PLANNER_INLINE_LABEL_ROW_CLASS}>
                  <span className="min-w-0 break-words leading-snug">
                    {t('workbench.source.planner.namingPrefix')}
                  </span>
                  <input
                    aria-label={t('workbench.source.planner.namingPrefix')}
                    value={namingPrefix}
                    onChange={(event) => {
                      clearAppliedTemplate();
                      setPlannerFieldTouched((current) => ({ ...current, namingPrefix: true }));
                      setNamingPrefix(event.target.value);
                    }}
                    className={SOURCE_PLANNER_FIELD_CONTROL_CLASS}
                  />
                </label>
              </div>

              {/* 目標型態／字節序、偏移／倍率（各列左標籤、右控制項） */}
              <div className="flex flex-col gap-3">
                  <label className={SOURCE_PLANNER_INLINE_LABEL_ROW_CLASS}>
                    <span className="min-w-0 break-words leading-snug">
                      {t('workbench.source.planner.targetDataType')}
                    </span>
                    <select
                      aria-label={t('workbench.source.planner.targetDataType')}
                      value={targetDataType}
                      onChange={(event) => {
                        clearAppliedTemplate();
                        setPlannerFieldTouched((current) => ({
                          ...current,
                          targetDataType: true,
                        }));
                        setTargetDataType(event.target.value as DataType | '');
                      }}
                      className={SOURCE_PLANNER_FIELD_CONTROL_CLASS}
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
                  <label className={SOURCE_PLANNER_INLINE_LABEL_ROW_CLASS}>
                    <span className="min-w-0 break-words leading-snug">
                      {t('workbench.source.planner.dataFormat')}
                    </span>
                    <select
                      aria-label={t('workbench.source.planner.dataFormat')}
                      value={dataFormat}
                      onChange={(event) => {
                        clearAppliedTemplate();
                        setDataFormat(event.target.value);
                      }}
                      className={SOURCE_PLANNER_FIELD_CONTROL_CLASS}
                    >
                      <option value="">{t('workbench.source.planner.dataFormatDefault')}</option>
                      <option value="CDAB">CDAB</option>
                      <option value="ABCD">ABCD</option>
                      <option value="BADC">BADC</option>
                      <option value="DCBA">DCBA</option>
                    </select>
                  </label>
                  <label className={SOURCE_PLANNER_INLINE_LABEL_ROW_CLASS}>
                    <span className="min-w-0 break-words leading-snug">
                      {t('workbench.source.planner.scaleOffset')}
                    </span>
                    <input
                      aria-label={t('workbench.source.planner.scaleOffset')}
                      type="number"
                      step="any"
                      placeholder="0.0"
                      value={scaleOffset}
                      onChange={(event) => {
                        setPlannerFieldTouched((current) => ({ ...current, scaleOffset: true }));
                        setScaleOffset(event.target.value);
                      }}
                      className={SOURCE_PLANNER_FIELD_CONTROL_MONO_CLASS}
                    />
                  </label>
                  <label className={SOURCE_PLANNER_INLINE_LABEL_ROW_CLASS}>
                    <span className="min-w-0 break-words leading-snug">
                      {t('workbench.source.planner.scaleMultiplier')}
                    </span>
                    <input
                      aria-label={t('workbench.source.planner.scaleMultiplier')}
                      type="number"
                      step="any"
                      placeholder="1.0"
                      value={scaleMultiplier}
                      onChange={(event) => {
                        setPlannerFieldTouched((current) => ({
                          ...current,
                          scaleMultiplier: true,
                        }));
                        setScaleMultiplier(event.target.value);
                      }}
                      className={SOURCE_PLANNER_FIELD_CONTROL_MONO_CLASS}
                    />
                  </label>
              </div>

              {activeOutputTarget === 'database' ? (
                <div
                  className="rounded-lg border border-violet-500/20 bg-violet-500/[0.08] px-3 py-2 text-xs text-violet-100"
                  data-testid="source-planner-database-defaults"
                >
                  {t('workbench.source.databaseHints.formDefaults')}
                </div>
              ) : null}

              <button
                type="button"
                onClick={handleApplyPlan}
                disabled={parsePlannerCountInput(count) === null}
                className="w-full rounded-lg bg-cyan-500 px-3 py-2 text-sm font-medium text-slate-950 disabled:cursor-not-allowed disabled:bg-slate-600 disabled:text-slate-300"
              >
                {t('workbench.source.planner.addRule')}
              </button>
                  </div>
              </div>

              <div
                role="tabpanel"
                id="source-rule-layer-panel-rules"
                aria-labelledby="source-rule-layer-tab-rules"
                hidden={ruleLayerTab !== 'rules'}
                className="absolute inset-0 overflow-y-auto overscroll-contain pr-0.5 scrollbar-auto-hide"
              >
            {rules.length === 0 ? (
              <p className="rounded-xl border border-dashed border-slate-700/60 bg-slate-900/30 px-3 py-4 text-xs text-slate-400">
                {t('workbench.source.ruleLayer.empty')}
              </p>
            ) : (
              <div className="space-y-2">
                {rules.map((rule, ruleIndex) => {
                  const coverage = buildSourceRuleCoverage(rule, selectedDevice.protocol);
                  const isSelected = rule.id === selectedRuleId;
                  const isEditing = editingRuleId === rule.id;
                  const isRuleCardCollapsed = ruleCardCollapsedIds.has(rule.id);
                  const canMoveRuleUp = ruleIndex > 0;
                  const canMoveRuleDown = ruleIndex < rules.length - 1;
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
                      <div className="flex gap-2 sm:gap-3">
                        <div className="flex shrink-0 flex-col items-center gap-1">
                          <button
                            type="button"
                            className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg border border-slate-800/80 bg-slate-900/90 text-slate-400 shadow-sm shadow-black/25 transition hover:border-slate-600/80 hover:bg-slate-800/90 hover:text-slate-100 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-cyan-400"
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
                          {isRuleCardCollapsed ? null : (
                            <>
                              <button
                                type="button"
                                className="flex h-8 w-9 shrink-0 items-center justify-center rounded-lg border border-slate-700/70 bg-slate-900/80 text-slate-300 transition hover:border-slate-500/60 hover:bg-slate-800/80 hover:text-slate-100 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-cyan-400 disabled:cursor-not-allowed disabled:opacity-40"
                                aria-label={t('workbench.source.ruleLayer.moveUp')}
                                title={t('workbench.source.ruleLayer.moveUp')}
                                disabled={!canMoveRuleUp}
                                data-testid={`source-rule-${rule.id}-move-up`}
                                onClick={(event) => {
                                  event.stopPropagation();
                                  handleMoveRule(rule.id, -1);
                                }}
                              >
                                <ArrowUp className="h-4 w-4" aria-hidden />
                              </button>
                              <button
                                type="button"
                                className="flex h-8 w-9 shrink-0 items-center justify-center rounded-lg border border-slate-700/70 bg-slate-900/80 text-slate-300 transition hover:border-slate-500/60 hover:bg-slate-800/80 hover:text-slate-100 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-cyan-400 disabled:cursor-not-allowed disabled:opacity-40"
                                aria-label={t('workbench.source.ruleLayer.moveDown')}
                                title={t('workbench.source.ruleLayer.moveDown')}
                                disabled={!canMoveRuleDown}
                                data-testid={`source-rule-${rule.id}-move-down`}
                                onClick={(event) => {
                                  event.stopPropagation();
                                  handleMoveRule(rule.id, 1);
                                }}
                              >
                                <ArrowDown className="h-4 w-4" aria-hidden />
                              </button>
                            </>
                          )}
                        </div>
                        <button
                          type="button"
                          data-testid={`source-rule-${rule.id}`}
                          onClick={() => handleSelectRule(rule.id)}
                          className={[
                            'min-w-0 flex-1 rounded-xl border p-2.5 text-left shadow-inner shadow-black/15 transition',
                            'focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-cyan-400',
                            isSelected
                              ? 'border-cyan-500/40 bg-slate-950/50 ring-1 ring-cyan-400/25'
                              : 'border-slate-800/60 bg-slate-950/40 hover:border-slate-600/50 hover:bg-slate-950/65',
                          ].join(' ')}
                        >
                          <div className="flex min-w-0 flex-col gap-2">
                            <div className="flex min-w-0 flex-wrap items-center gap-2">
                              <span className="truncate font-mono text-sm font-semibold tracking-tight text-slate-50">
                                {rule.id}
                              </span>
                              <span
                                className={[
                                  'shrink-0 rounded-md px-2 py-0.5 text-[10px] font-semibold uppercase tracking-[0.14em]',
                                  rule.persisted
                                    ? 'border border-emerald-500/35 bg-emerald-500/[0.12] text-emerald-100'
                                    : 'border border-slate-600/50 bg-slate-800/60 text-slate-200',
                                ].join(' ')}
                              >
                                {rule.persisted
                                  ? t('workbench.source.ruleLayer.persistedBadge')
                                  : t('workbench.source.ruleLayer.draftBadge')}
                              </span>
                            </div>
                            <div className="flex min-w-0 items-center gap-2 rounded-lg bg-slate-900/80 px-2 py-1.5 font-mono text-[11px] tabular-nums text-cyan-200/90 ring-1 ring-slate-700/50">
                              <span className="min-w-0 truncate">{coverage.startAddress}</span>
                              <span className="shrink-0 text-slate-500" aria-hidden>
                                →
                              </span>
                              <span className="min-w-0 truncate">{coverage.endAddress}</span>
                            </div>
                            <p className="text-[11px] leading-relaxed text-slate-500">
                              {t('workbench.source.ruleLayer.meta', {
                                bitWidth: coverage.bitWidth,
                                count: rule.count,
                                cells: coverage.cellCount,
                              })}
                            </p>
                          </div>
                        </button>
                      </div>

                      {!isRuleCardCollapsed && isEditing && editDraft ? (
                        <div
                          className="space-y-2 rounded-lg border border-cyan-500/20 bg-slate-950/60 p-2 text-[11px] text-slate-300"
                          data-testid="rule-inline-edit-form"
                        >
                          <div className="flex flex-col gap-2">
                            <label className={SOURCE_PLANNER_INLINE_LABEL_ROW_CLASS_COMPACT}>
                              <span className="min-w-0 break-words leading-snug">
                                {t('workbench.source.planner.startAddress')}
                              </span>
                              <input
                                aria-label={t('workbench.source.planner.startAddress')}
                                value={editDraft.startAddress}
                                onChange={(event) =>
                                  setEditDraft((draft) =>
                                    draft ? { ...draft, startAddress: event.target.value } : draft,
                                  )
                                }
                                className={SOURCE_PLANNER_FIELD_CONTROL_MONO_CLASS}
                              />
                            </label>
                            <label className={SOURCE_PLANNER_INLINE_LABEL_ROW_CLASS_COMPACT}>
                              <span className="min-w-0 break-words leading-snug">
                                {t('workbench.source.planner.count')}
                              </span>
                              <input
                                aria-label={t('workbench.source.planner.count')}
                                type="number"
                                value={editDraft.count}
                                onChange={(event) =>
                                  setEditDraft((draft) =>
                                    draft ? { ...draft, count: event.target.value } : draft,
                                  )
                                }
                                className={SOURCE_PLANNER_FIELD_CONTROL_MONO_CLASS}
                              />
                            </label>
                            <label className={SOURCE_PLANNER_INLINE_LABEL_ROW_CLASS_COMPACT}>
                              <span className="min-w-0 break-words leading-snug">
                                {t('workbench.source.planner.dataType')}
                              </span>
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
                                className={SOURCE_PLANNER_FIELD_CONTROL_CLASS}
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
                          <div className="flex flex-col gap-2">
                              <label className={SOURCE_PLANNER_INLINE_LABEL_ROW_CLASS_COMPACT}>
                                <span className="min-w-0 break-words leading-snug">
                                  {t('workbench.source.planner.targetDataType')}
                                </span>
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
                                  className={SOURCE_PLANNER_FIELD_CONTROL_CLASS}
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
                              <label className={SOURCE_PLANNER_INLINE_LABEL_ROW_CLASS_COMPACT}>
                                <span className="min-w-0 break-words leading-snug">
                                  {t('workbench.source.planner.dataFormat')}
                                </span>
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
                                  className={SOURCE_PLANNER_FIELD_CONTROL_CLASS}
                                >
                                  <option value="">{t('workbench.source.planner.dataFormatDefault')}</option>
                                  <option value="CDAB">CDAB</option>
                                  <option value="ABCD">ABCD</option>
                                  <option value="BADC">BADC</option>
                                  <option value="DCBA">DCBA</option>
                                </select>
                              </label>
                              <label className={SOURCE_PLANNER_INLINE_LABEL_ROW_CLASS_COMPACT}>
                                <span className="min-w-0 break-words leading-snug">
                                  {t('workbench.source.planner.scaleOffset')}
                                </span>
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
                                  className={SOURCE_PLANNER_FIELD_CONTROL_MONO_CLASS}
                                />
                              </label>
                              <label className={SOURCE_PLANNER_INLINE_LABEL_ROW_CLASS_COMPACT}>
                                <span className="min-w-0 break-words leading-snug">
                                  {t('workbench.source.planner.scaleMultiplier')}
                                </span>
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
                                  className={SOURCE_PLANNER_FIELD_CONTROL_MONO_CLASS}
                                />
                              </label>
                          </div>
                          <div className="flex gap-2">
                            <button
                              type="button"
                              onClick={() => handleSaveRuleEdit(rule.id)}
                              disabled={parsePlannerCountInput(editDraft.count) === null}
                              className="rounded-lg bg-cyan-500 px-2 py-1 text-[11px] font-medium text-slate-950 disabled:cursor-not-allowed disabled:bg-slate-600 disabled:text-slate-300"
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
                          className="min-w-0 truncate rounded-lg border border-cyan-500/30 bg-cyan-500/5 px-2 py-1 text-center text-cyan-200"
                          onClick={() => handleStartRuleEdit(rule.id)}
                          type="button"
                        >
                          {t('workbench.source.ruleLayer.editStart')}
                        </button>
                        <button
                          className="min-w-0 truncate rounded-lg border border-rose-500/30 bg-rose-500/10 px-2 py-1 text-center text-rose-100"
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
              </div>
            </div>
          </section>
        </aside>
        </SourceRuleLayerPanel>

        <div
          className="flex min-h-0 min-w-0 flex-1 flex-col overflow-hidden"
          data-testid={
            deskMode === 'build'
              ? 'source-build-canvas-panel'
              : deskMode === 'triage'
                ? 'source-triage-canvas-panel'
                : 'source-inspect-canvas-panel'
          }
        >
        <section
          className="flex min-h-0 min-w-0 flex-1 flex-col gap-4 overflow-hidden rounded-2xl border border-slate-800 bg-slate-950/40 p-4"
          data-emphasis="supporting"
          data-testid="source-canvas-workspace"
        >
          <div
            className="flex shrink-0 flex-col gap-3 xl:flex-row xl:items-center xl:justify-between"
            data-testid="source-primary-toolbar"
          >
            <div
              className="flex flex-1 flex-col gap-2 xl:max-w-2xl"
              data-testid="source-diagnostic-toolbar"
            >
              <div>
                <p className="text-[10px] font-semibold uppercase tracking-[0.18em] text-slate-500">
                  {t('workbench.source.diagnostics.eyebrow')}
                </p>
                <p className="mt-1 text-sm text-slate-300">
                  {t('workbench.source.diagnostics.description')}
                </p>
              </div>
              <div className="flex flex-wrap items-center gap-2">
                {(['plan', 'live', 'link'] as const).map((mode) => (
                  <button
                    key={mode}
                    type="button"
                    onClick={() => setViewMode(mode)}
                    className={
                      viewMode === mode
                        ? 'rounded-lg border border-cyan-500/40 bg-cyan-500/12 px-3 py-1.5 text-xs font-semibold text-cyan-100'
                        : 'rounded-lg border border-slate-800 bg-slate-900 px-3 py-1.5 text-xs text-slate-300'
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
                  className="rounded-lg border border-slate-800 bg-slate-900 px-3 py-1.5 text-xs text-slate-200"
                >
                  <option value="decimal">{t('workbench.source.formats.decimal')}</option>
                  <option value="hex">{t('workbench.source.formats.hex')}</option>
                  <option value="binary">{t('workbench.source.formats.binary')}</option>
                  <option value="float">{t('workbench.source.formats.float')}</option>
                </select>
              </div>
            </div>
            <div className="flex shrink-0 flex-wrap items-center gap-2 xl:justify-end">
              <label
                className="inline-flex max-w-[9.5rem] shrink-0 flex-col gap-0.5 sm:max-w-none sm:inline-flex sm:flex-row sm:items-center sm:gap-1.5"
                title={t('workbench.source.canvas.latticeColumnsHint')}
              >
                <span className="whitespace-nowrap text-[10px] font-medium uppercase tracking-wide text-slate-500">
                  {t('workbench.source.canvas.latticeColumnsLabel')}
                </span>
                <input
                  type="number"
                  min={LATTICE_COLUMNS_MIN}
                  max={LATTICE_COLUMNS_MAX}
                  step={1}
                  inputMode="numeric"
                  aria-label={t('workbench.source.canvas.latticeColumnsLabel')}
                  data-testid="source-canvas-lattice-columns"
                  value={canvasLatticeColumns}
                  onChange={(event) => handleCanvasLatticeColumnsChange(event.target.value)}
                  onBlur={() => {
                    const next = clampLatticeColumns(canvasLatticeColumns);
                    if (next !== canvasLatticeColumns) {
                      setCanvasLatticeColumns(next);
                      writeSourceCanvasLatticeColumns(next);
                    }
                  }}
                  className="w-full min-w-[4.25rem] rounded-md border border-slate-700 bg-slate-950 px-2 py-1.5 font-mono text-xs text-slate-100 sm:w-[4.5rem]"
                />
              </label>
              <div ref={sourceToolbarMoreRef} className="relative">
                <button
                  type="button"
                  data-testid="source-toolbar-more-trigger"
                  aria-expanded={sourceToolbarMoreOpen}
                  aria-haspopup="menu"
                  aria-controls="source-toolbar-more-menu"
                  aria-label={t('workbench.source.toolbar.moreMenuAria')}
                  title={t('workbench.source.toolbar.moreMenuAria')}
                  className="inline-flex items-center rounded-lg border border-slate-800 bg-slate-900/70 px-3 py-2 text-xs font-medium text-slate-200 transition hover:border-slate-600 hover:bg-slate-800/80"
                  onClick={() => setSourceToolbarMoreOpen((open) => !open)}
                >
                  {t('workbench.source.toolbar.moreMenu')}
                  <ChevronDown
                    className={`ml-1 h-3.5 w-3.5 shrink-0 text-slate-400 transition-transform ${sourceToolbarMoreOpen ? 'rotate-180' : ''}`}
                    aria-hidden
                  />
                </button>
                {sourceToolbarMoreOpen ? (
                  <div
                    id="source-toolbar-more-menu"
                    role="menu"
                    aria-label={t('workbench.source.toolbar.moreMenuAria')}
                    data-testid="source-toolbar-more-menu"
                    className="absolute right-0 z-50 mt-1 min-w-[13rem] rounded-xl border border-slate-700 bg-slate-900 py-1 shadow-xl shadow-black/50"
                  >
                    <button
                      type="button"
                      role="menuitem"
                      className="flex w-full items-center px-3 py-2 text-left text-xs text-slate-200 transition-colors hover:bg-slate-800 focus-visible:bg-slate-800 focus-visible:outline-none"
                      onClick={() => {
                        handleToggleFreezeLive();
                        setSourceToolbarMoreOpen(false);
                      }}
                      title={
                        freezeLive
                          ? t('workbench.source.toolbar.unfreezeLiveHint')
                          : t('workbench.source.toolbar.freezeLiveHint')
                      }
                    >
                      {freezeLive
                        ? t('workbench.source.toolbar.unfreezeLive')
                        : t('workbench.source.toolbar.freezeLive')}
                    </button>
                    <button
                      type="button"
                      role="menuitem"
                      className="flex w-full items-center px-3 py-2 text-left text-xs text-slate-200 transition-colors hover:bg-slate-800 focus-visible:bg-slate-800 focus-visible:outline-none"
                      onClick={() => {
                        handleCaptureSnapshot();
                        setSourceToolbarMoreOpen(false);
                      }}
                      title={t('workbench.source.toolbar.snapshotCompareHint')}
                    >
                      {t('workbench.source.toolbar.snapshotCompare')}
                    </button>
                    <button
                      type="button"
                      role="menuitem"
                      className="flex w-full items-center px-3 py-2 text-left text-xs text-slate-200 transition-colors hover:bg-slate-800 focus-visible:bg-slate-800 focus-visible:outline-none"
                      onClick={() => {
                        setShowAudit((currentValue) => !currentValue);
                        setSourceToolbarMoreOpen(false);
                      }}
                      title={
                        showAudit
                          ? t('workbench.source.toolbar.hideAuditHint')
                          : t('workbench.source.toolbar.showAuditHint')
                      }
                    >
                      {showAudit
                        ? t('workbench.source.toolbar.hideAudit')
                        : t('workbench.source.toolbar.showAudit')}
                    </button>
                    <button
                      type="button"
                      role="menuitem"
                      className="flex w-full items-center px-3 py-2 text-left text-xs text-slate-200 transition-colors hover:bg-slate-800 focus-visible:bg-slate-800 focus-visible:outline-none"
                      onClick={() => {
                        handleOpenSaveTemplate();
                        setSourceToolbarMoreOpen(false);
                      }}
                      title={t('workbench.source.toolbar.saveTemplateHint')}
                    >
                      {t('workbench.source.toolbar.saveTemplate')}
                    </button>
                    <button
                      type="button"
                      role="menuitem"
                      disabled={templates.length === 0}
                      className="flex w-full items-center px-3 py-2 text-left text-xs text-slate-200 transition-colors hover:bg-slate-800 focus-visible:bg-slate-800 focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-45 disabled:hover:bg-slate-900"
                      onClick={() => {
                        handleOpenLoadTemplate();
                        setSourceToolbarMoreOpen(false);
                      }}
                      title={t('workbench.source.toolbar.loadTemplateHint')}
                    >
                      {t('workbench.source.toolbar.loadTemplate')}
                    </button>
                  </div>
                ) : null}
              </div>
              <button
                type="button"
                onClick={() => void handleCreateSelectedPoint()}
                disabled={!selectedPointDefinition || createPointMutation.isPending}
                title={t('workbench.source.actions.createSelectedPointsHint')}
                className="rounded-lg border border-slate-700/80 bg-slate-900/70 px-3 py-2 text-xs font-medium text-slate-200 transition hover:border-slate-600 hover:bg-slate-800/80 disabled:cursor-not-allowed disabled:border-slate-800 disabled:bg-slate-950/50 disabled:text-slate-600"
              >
                {t('workbench.source.actions.createSelectedPoints')}
              </button>
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
                latticeColumns={canvasLatticeColumns}
                onSelectAddress={handleSelectAddress}
                selectedAddress={selectedAddress}
                valueFormat={valueFormat}
                viewMode={viewMode}
              />
            </div>

            <div
              ref={sourceWorkspaceSecondaryScrollRef}
              className={[
                'min-h-0 touch-pan-y space-y-4 overflow-y-auto overscroll-contain pb-4 scrollbar-none',
                deskMode === 'build' ? 'max-h-[min(30vh,16rem)]' : 'max-h-[min(38vh,22rem)]',
              ].join(' ')}
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
                title={t('workbench.source.actions.createSelectedPointsHint')}
                className="rounded-lg border border-cyan-400/35 bg-cyan-500/10 px-3 py-1.5 text-xs font-medium text-cyan-100 transition hover:bg-cyan-500/15 disabled:cursor-not-allowed disabled:border-slate-800 disabled:bg-slate-950/40 disabled:text-slate-500"
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
            <div className="flex flex-wrap items-start justify-between gap-3">
              <button
                type="button"
                data-testid="source-coverage-overview-toggle"
                aria-expanded={coverageOverviewOpen}
                aria-label={
                  coverageOverviewOpen
                    ? t('workbench.source.coverage.collapseSection')
                    : t('workbench.source.coverage.expandSection')
                }
                onClick={() => setCoverageOverviewOpen((open) => !open)}
                className="flex min-w-0 flex-1 items-start gap-2 rounded-lg py-0.5 text-left transition-colors hover:bg-slate-800/35 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-cyan-400"
              >
                <span className="mt-0.5 shrink-0 text-slate-400" aria-hidden>
                  {coverageOverviewOpen ? (
                    <ChevronDown className="h-4 w-4" />
                  ) : (
                    <ChevronRight className="h-4 w-4" />
                  )}
                </span>
                <div className="min-w-0 space-y-1">
                  <p className="text-sm font-semibold uppercase tracking-[0.18em] text-cyan-300">
                    {t('workbench.source.coverage.eyebrow')}
                  </p>
                  <h3 className="text-sm font-semibold text-slate-100">
                    {t('workbench.source.coverage.title')}
                  </h3>
                </div>
              </button>
              {coverageSegments.length > 0 ? (
                <p className="shrink-0 text-xs text-slate-400">
                  {t('workbench.source.coverage.cells', {
                    count: coverageSegments.reduce(
                      (total, segment) => total + segment.cellCount,
                      0,
                    ),
                  })}
                </p>
              ) : null}
            </div>
            {coverageOverviewOpen ? (
              coverageSegments.length > 0 ? (
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
                </div>
              ) : (
                <p className="text-xs text-slate-400">
                  {t('workbench.source.coverage.empty')}
                </p>
              )
            ) : null}
          </section>

          {isSaveTemplateOpen ? (
            <div className="space-y-3 rounded-xl border border-slate-800 bg-slate-900/50 p-4">
              <label className={SOURCE_PLANNER_INLINE_LABEL_ROW_CLASS}>
                <span className="min-w-0 break-words leading-snug">{t('workbench.source.templates.name')}</span>
                <input
                  aria-label={t('workbench.source.templates.name')}
                  value={templateName}
                  onChange={(event) => setTemplateName(event.target.value)}
                  className={SOURCE_PLANNER_FIELD_CONTROL_CLASS}
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

          {deskMode !== 'triage' && conflictQueue.length > 0 ? (
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
                                setRuleLayerTab('rules');
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

          {templateRecoveryWarning || templateWarning ? (
            <p className="text-sm text-amber-200" data-testid="source-template-warning">
              {templateWarning ?? templateRecoveryWarning}
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
    </div>
  );

  const primaryWorkspace =
    deskMode === 'triage' ? (
      <SourceTriagePanel recoveryPanel={triageRecoveryPanel}>{workspaceGrid}</SourceTriagePanel>
    ) : (
      <div className="flex min-h-0 flex-1 flex-col overflow-hidden" data-testid={workspaceModeTestId}>
        {workspaceGrid}
      </div>
    );

  return (
    <section className="flex h-full min-h-0 flex-col gap-6 overflow-hidden">
      {deskMode === 'inspect' ? (
        <>
          {workspaceSummaryStrip}
          {workspaceDiagnosticsStrip}
          {primaryWorkspace}
        </>
      ) : (
        primaryWorkspace
      )}
    </section>
  );
}
