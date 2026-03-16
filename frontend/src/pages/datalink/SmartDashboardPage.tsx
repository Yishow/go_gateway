import { useState, useMemo, useCallback, useEffect, useRef } from "react";
import { useTranslation } from "react-i18next";
import type {
  Device,
  Point,
  Tag,
  CreatePointRequest,
  CreateDeviceRequest,
  UpdateDeviceRequest,
  ProtocolType,
  DataType,
} from "../../types/datalink";
import { tagAPI } from "../../services/datalink";
import type { PlannedAllocation } from "../../components/datalink/MemoryGrid";
import { ImportDialog, ExportDialog } from "../../components/datalink/ImportExportDialog";
import SmartDashboardHeader from "./smart-dashboard/SmartDashboardHeader";
import SmartDashboardIntentNotices from "./smart-dashboard/SmartDashboardIntentNotices";
import SmartDashboardControlBar from "./smart-dashboard/SmartDashboardControlBar";
import type { SidebarTab } from "./smart-dashboard/SmartDashboardSidebar";
import SmartDashboardWorkspaceSection from "./smart-dashboard/SmartDashboardWorkspaceSection";
import { useSmartDashboardTagLinking } from "./smart-dashboard/useSmartDashboardTagLinking";
import { useSmartDashboardModbusActions } from "./smart-dashboard/useSmartDashboardModbusActions";
import { useSmartDashboardSidebarPanelMotion } from "./smart-dashboard/useSmartDashboardSidebarPanelMotion";
import SmartDashboardWorkflowModal from "./smart-dashboard/SmartDashboardWorkflowModal";
import SmartDashboardOverlays from "./smart-dashboard/SmartDashboardOverlays";
import SmartDashboardPanels from "./smart-dashboard/SmartDashboardPanels";
import SmartDashboardGridOverlaysSection from "./smart-dashboard/SmartDashboardGridOverlaysSection";
import { useSmartDashboardGridOverlays } from "./smart-dashboard/useSmartDashboardGridOverlays";
import { useSmartDashboardCommitFlow } from "./smart-dashboard/useSmartDashboardCommitFlow";
import { useSmartDashboardWorkspaceState } from "./smart-dashboard/useSmartDashboardWorkspaceState";
import {
  useDevicesQuery,
  useDeleteDeviceMutation,
  useToggleDeviceStatusMutation,
  useUpdateDeviceMutation,
  useTestConnectionMutation,
} from "../../hooks/datalink/useDevices";
import { usePollingGroupsQuery } from "../../hooks/datalink/usePollingGroups";
import { usePointsQuery, useCreatePointMutation, useDeletePointMutation } from "../../hooks/datalink/usePoints";
import { useToast } from "../../contexts/ToastContext";
import {
  useMappingsQuery,
  useValidatePipelineMutation,
  useCreateMappingMutation,
  useUpdateMappingMutation,
} from "../../hooks/datalink/useMappings";
import { useQueryClient } from "@tanstack/react-query";
import { useTagsQuery, useCreateTagMutation, useUpdateTagMutation } from "../../hooks/datalink/useTags";
import { tagKeys, mappingKeys } from "../../hooks/datalink/keys";
import { useSmartDashboardShortcuts } from "../../hooks/useKeyboardShortcuts";
import { usePointHistory } from "../../hooks/useHistory";
import { useFlowLifecycle, type FlowSegment, type FlowStatus } from "../../features/flow/stateMachine";
import {
  SOURCE_TEMPLATE_SCHEMA_VERSION,
  isTemplateStale,
  loadSourceTemplates,
  saveSourceTemplates,
  upgradeTemplates,
} from "../../features/datalink/sourceTemplateStorage";
import { buildBatchNamePreview } from "../../features/datalink/batchNaming";
import {
  applyTemplateToPlanner,
  createTemplateFromPlanner,
  normalizeNamingPrefix,
  upsertTemplateRecord,
} from "../../features/datalink/sourcePlannerContract";
import { findNearestValidContiguousSpan } from "../../features/datalink/allocationStrategy";
import { getAffectedMappingCountForTag } from "../../features/datalink/tagEditImpact";
import { buildGlobalTagGuardrail } from "../../features/datalink/globalTagGuardrails";
import {
  buildWorkbenchRedirect,
  isDashboardModalIntent,
  isDashboardSectionIntent,
  isLegacyDecommissionRoute,
} from "../../features/datalink/legacyRoutes";
import {
  MOTION_TOKENS,
  resolveIntentMotionClass,
  resolveScrollBehavior,
} from "../../features/datalink/motionGuidance";
import { getSpanByDataType, validateTypedOccupancyPlan } from "../../features/datalink/typedOccupancy";
import { runStructuralValidation } from "../../features/datalink/validationFlow";
import { addressParser } from "../../utils/addressParser";
import { useNavigate, useSearchParams } from "react-router-dom";

const FLOW_SEGMENTS: FlowSegment[] = ["source", "grid", "tag", "sink"];
type DashboardTab = "overview" | "devices" | "settings";

const STATUS_STYLE: Record<FlowStatus, string> = {
  draft: "bg-slate-700/70 text-slate-200 border-slate-600",
  validated: "bg-blue-500/20 text-blue-300 border-blue-500/40",
  active: "bg-emerald-500/20 text-emerald-300 border-emerald-500/40",
  error: "bg-red-500/20 text-red-300 border-red-500/40",
};

const RESET_SEGMENT_DIAGNOSTIC = {
  latestValue: "-",
  quality: "unknown" as const,
  timestamp: "-",
  error: "",
};

interface SourceTemplate {
  id: string;
  name: string;
  dataType: DataType;
  count: number;
  startAddress: string;
  updatedAt: string;
  lastUsedAt: string;
  version: number;
}

export default function SmartDashboard() {
  const { t } = useTranslation();
  const { showSuccess, showError, showInfo } = useToast();
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();
  const [selectedDeviceId, setSelectedDeviceId] = useState<string | null>(null);
  const [selectedAddresses, setSelectedAddresses] = useState<string[]>([]);
  const [selectedPoint, setSelectedPoint] = useState<Point | null>(null);
  const [deviceSearchQuery, setDeviceSearchQuery] = useState("");
  const [deviceStatusFilter, setDeviceStatusFilter] = useState<"all" | "active" | "disabled" | "draft">("all");
  const [isSwitchingDevice, setIsSwitchingDevice] = useState(false);
  const [pendingSwitchDeviceId, setPendingSwitchDeviceId] = useState<string | null>(null);
  const [activatingDeviceId, setActivatingDeviceId] = useState<string | null>(null);
  const [showSwitchConfirmDialog, setShowSwitchConfirmDialog] = useState(false);
  const [panelType, setPanelType] = useState<"batch" | "detail" | "shortcuts" | null>(null);
  /** 「套用到網格」帶入的批量命名模板，關閉批量面板時清除 */
  const [batchInitialTemplate, setBatchInitialTemplate] = useState<string | null>(null);
  const [importDialogOpen, setImportDialogOpen] = useState(false);
  const [exportDialogOpen, setExportDialogOpen] = useState(false);
  const [planDataType, setPlanDataType] = useState<DataType>("int16");
  const [planCount, setPlanCount] = useState(5);
  const [batchNamePrefix, setBatchNamePrefix] = useState("SRC");
  const [planStartAddress, setPlanStartAddress] = useState("");
  /** Grid 顯示視窗起始位址；與 planStartAddress 同步，滾輪/按鈕可偏移 100 格 */
  const [gridViewStartAddress, setGridViewStartAddress] = useState("");
  const [allocationMessage, setAllocationMessage] = useState("");
  const [templateName, setTemplateName] = useState("");
  const [showConflictsOnly, setShowConflictsOnly] = useState(false);
  const [guideStage, setGuideStage] = useState<"idle" | "grid" | "commit">("idle");
  const [reducedMotion, setReducedMotion] = useState(false);
  const [sourceTemplates, setSourceTemplates] = useState<SourceTemplate[]>(() => loadSourceTemplates());
  const legacyRoute = searchParams.get("legacy");
  const rawSectionIntent = searchParams.get("section");
  const rawModalIntent = searchParams.get("modal");
  const createDeviceIntent = searchParams.get("createDevice");
  const [activeTab, setActiveTab] = useState<DashboardTab>("overview");
  const [sidebarTab, setSidebarTab] = useState<SidebarTab>("plan");
  const [isCreateDeviceModalOpen, setIsCreateDeviceModalOpen] = useState(false);
  const [justCreatedDeviceId, setJustCreatedDeviceId] = useState<string | null>(null);
  const [editingDeviceInModal, setEditingDeviceInModal] = useState<Device | null>(null);
  const [testingDeviceId, setTestingDeviceId] = useState<string | null>(null);
  const [deletingDeviceId, setDeletingDeviceId] = useState<string | null>(null);
  const [deleteConfirmDevice, setDeleteConfirmDevice] = useState<Device | null>(null);
  const searchInputRef = useRef<HTMLInputElement | null>(null);
  const gridSectionRef = useRef<HTMLElement | null>(null);
  const guideStageTimeoutRef = useRef<number | null>(null);

  const { data: devices = [] } = useDevicesQuery();
  const deleteDeviceMutation = useDeleteDeviceMutation();
  const toggleDeviceStatusMutation = useToggleDeviceStatusMutation();
  const updateDeviceMutation = useUpdateDeviceMutation();
  const testConnectionMutation = useTestConnectionMutation();
  const { data: pollingGroups = [] } = usePollingGroupsQuery();
  const { data: allPoints = [] } = usePointsQuery({ device_id: selectedDeviceId || undefined });
  const {
    gridPopoverOpen,
    setGridPopoverOpen,
    gridPopoverPoint,
    gridPopoverAddress,
    gridPopoverPosition,
    gridPopoverContentRef,
    gridContextMenuRef,
    gridContextMenu,
    setGridContextMenu,
    handleCellClick,
    handleCellContextMenu,
  } = useSmartDashboardGridOverlays({
    allPoints,
    selectedAddresses,
  });
  const { data: mappings = [] } = useMappingsQuery();
  const queryClient = useQueryClient();
  const refetchMappings = useCallback(
    () => queryClient.refetchQueries({ queryKey: mappingKeys.lists() }),
    [queryClient]
  );
  const { data: tags = [] } = useTagsQuery();
  const createTagMutation = useCreateTagMutation();
  const updateTagMutation = useUpdateTagMutation();
  const createMappingMutation = useCreateMappingMutation();
  const updateMappingMutation = useUpdateMappingMutation();
  const createPointMutation = useCreatePointMutation();
  const deletePointMutation = useDeletePointMutation();
  const history = usePointHistory({ maxHistory: 30 });
  const flow = useFlowLifecycle();
  const {
    state: flowState,
    canValidate,
    canActivate,
    hasError,
    setSource,
    setTag,
    setSink,
    setDiagnostics,
    markValidated,
    markActive,
    markError,
    resetDraft,
  } = flow;
  const validatePipelineMutation = useValidatePipelineMutation();

  const selectedDevice = useMemo(
    () => devices.find((d) => d.id === selectedDeviceId) || null,
    [devices, selectedDeviceId],
  );
  const filteredDevices = useMemo(() => {
    const keyword = deviceSearchQuery.trim().toLowerCase();
    return devices.filter((device) => {
      const byStatus = deviceStatusFilter === "all" || device.status === deviceStatusFilter;
      if (!byStatus) return false;
      if (!keyword) return true;
      return (
        device.name.toLowerCase().includes(keyword) ||
        device.protocol.toLowerCase().includes(keyword) ||
        device.id.toLowerCase().includes(keyword)
      );
    });
  }, [deviceSearchQuery, deviceStatusFilter, devices]);
  const deviceSummary = useMemo(
    () => ({
      total: devices.length,
      active: devices.filter((device) => device.status === "active").length,
      disabled: devices.filter((device) => device.status === "disabled").length,
      draft: devices.filter((device) => device.status === "draft").length,
    }),
    [devices],
  );

  const legacyRouteLabel = useMemo(() => {
    if (legacyRoute === "points") return t("nav.points");
    if (legacyRoute === "mappings") return t("nav.mappings");
    if (legacyRoute === "wizard") return t("nav.mappingWizard");
    return "";
  }, [legacyRoute, t]);
  const sectionIntent = useMemo(
    () => (isDashboardSectionIntent(rawSectionIntent) ? rawSectionIntent : null),
    [rawSectionIntent],
  );
  const modalIntent = useMemo(() => (isDashboardModalIntent(rawModalIntent) ? rawModalIntent : null), [rawModalIntent]);
  const sectionIntentLabel = useMemo(() => {
    if (sectionIntent === "devices") return t("nav.devices");
    if (sectionIntent === "settings") return t("nav.settings");
    return "";
  }, [sectionIntent, t]);
  const modalIntentLabel = useMemo(() => {
    if (!modalIntent) return "";
    if (modalIntent === "devices") return t("nav.devices");
    if (modalIntent === "settings") return t("nav.settings");
    if (modalIntent === "points") return t("nav.points");
    if (modalIntent === "mappings") return t("nav.mappings");
    if (modalIntent === "wizard") return t("nav.mappingWizard");
    if (modalIntent === "polling-groups") return t("nav.pollingGroups");
    if (modalIntent === "tags") return t("nav.tags");
    return modalIntent;
  }, [modalIntent, t]);

  const cellSpan = getSpanByDataType(planDataType);
  const typedPlanValidation = useMemo(
    () => validateTypedOccupancyPlan(planDataType, planCount),
    [planCount, planDataType],
  );
  const totalPlannedCells = typedPlanValidation.totalCells;

  const plannedAllocations = useMemo<PlannedAllocation[]>(() => {
    if (!selectedDevice || !planStartAddress || planCount <= 0) return [];
    const expanded = addressParser.expand(planStartAddress, totalPlannedCells, selectedDevice.protocol);
    if (expanded.length < totalPlannedCells) return [];

    const prefix = (batchNamePrefix.trim() || "SRC").toUpperCase();
    return Array.from({ length: planCount }).map((_, index) => {
      const start = index * cellSpan;
      const addresses = expanded.slice(start, start + cellSpan);
      return {
        id: `plan-${index}`,
        dataType: planDataType,
        addresses,
        label: `${prefix}-${String(index + 1).padStart(3, "0")}`,
      };
    });
  }, [batchNamePrefix, cellSpan, planCount, planDataType, planStartAddress, selectedDevice, totalPlannedCells]);

  const planAddresses = useMemo(
    () => plannedAllocations.flatMap((allocation) => allocation.addresses),
    [plannedAllocations],
  );

  const planConflictCount = useMemo(() => {
    const occupiedSet = new Set(allPoints.map((point) => point.address));
    return planAddresses.filter((address) => occupiedSet.has(address)).length;
  }, [allPoints, planAddresses]);
  const linkedAddresses = useMemo(() => {
    const pointAddressById = new Map(allPoints.map((point) => [point.id, point.address]));
    return mappings
      .map((mapping) => pointAddressById.get(mapping.point_id))
      .filter((address): address is string => Boolean(address));
  }, [allPoints, mappings]);
  const blockedAddresses = useMemo(() => {
    const blocked = new Set<string>();
    allPoints.forEach((point) => blocked.add(point.address));
    linkedAddresses.forEach((address) => blocked.add(address));
    return blocked;
  }, [allPoints, linkedAddresses]);
  const namePreview = useMemo(
    () =>
      buildBatchNamePreview(
        batchNamePrefix,
        planCount,
        allPoints.map((point) => point.name),
      ),
    [allPoints, batchNamePrefix, planCount],
  );
  const nameConflictCount = useMemo(() => namePreview.filter((item) => item.conflict).length, [namePreview]);

  const staleTemplateCount = useMemo(
    () => sourceTemplates.filter((template) => isTemplateStale(template)).length,
    [sourceTemplates],
  );
  const hasUnsavedChanges = useMemo(
    () => selectedAddresses.length > 0 || planAddresses.length > 0 || panelType === "batch" || panelType === "detail",
    [panelType, planAddresses.length, selectedAddresses.length],
  );

  useEffect(() => {
    setSelectedAddresses([]);
    setSelectedPoint(null);
    setPanelType((prev) => (prev === "detail" ? null : prev));
  }, [selectedDeviceId]);

  useEffect(() => {
    saveSourceTemplates(sourceTemplates);
  }, [sourceTemplates]);

  useEffect(() => {
    if (panelType !== "batch") setBatchInitialTemplate(null);
  }, [panelType]);

  useEffect(() => {
    const mediaQuery = window.matchMedia("(prefers-reduced-motion: reduce)");
    const apply = () => setReducedMotion(mediaQuery.matches);
    apply();
    mediaQuery.addEventListener("change", apply);
    return () => mediaQuery.removeEventListener("change", apply);
  }, []);
  useEffect(() => {
    return () => {
      if (guideStageTimeoutRef.current) {
        window.clearTimeout(guideStageTimeoutRef.current);
        guideStageTimeoutRef.current = null;
      }
    };
  }, []);

  /** 起始位址預設為空白，由使用者手動輸入或透過 Modbus 區域下拉選單填寫；不依 selectedDevice 自動填入 */

  const handleBatchCreate = useCallback(() => {
    if (selectedDeviceId) setPanelType("batch");
  }, [selectedDeviceId]);

  const handleClosePanel = useCallback(() => {
    setPanelType(null);
    setSelectedPoint(null);
  }, []);

  const applyDeviceSwitch = useCallback(
    async (nextDeviceId: string) => {
      setIsSwitchingDevice(true);
      setPendingSwitchDeviceId(nextDeviceId);
      try {
        const target = devices.find((device) => device.id === nextDeviceId);
        if (target && target.status !== "active" && target.status !== "disabled") {
          showError(`設備 ${target.name} 目前不可切換，請先完成啟用流程。`);
          showInfo("若設備離線，請先測試連線並確認來源協議設定。");
          setEditingDeviceInModal(target);
          const next = new URLSearchParams(searchParams);
          next.set("modal", "devices");
          setSearchParams(next, { replace: true });
          setPendingSwitchDeviceId(null);
          return;
        }
        setSelectedDeviceId(nextDeviceId);
        const next = new URLSearchParams(searchParams);
        if (next.get("modal") === "devices") {
          next.delete("modal");
          setSearchParams(next, { replace: true });
        }
        setPendingSwitchDeviceId(null);
      } catch (error) {
        const message = error instanceof Error ? error.message : "切換設備失敗";
        showError(message);
      } finally {
        setIsSwitchingDevice(false);
      }
    },
    [devices, searchParams, setSearchParams, showError, showInfo],
  );
  const requestDeviceSwitch = useCallback(
    (nextDeviceId: string) => {
      if (nextDeviceId === selectedDeviceId) {
        const next = new URLSearchParams(searchParams);
        if (next.get("modal") === "devices") {
          next.delete("modal");
          setSearchParams(next, { replace: true });
        }
        return;
      }
      if (hasUnsavedChanges && selectedDeviceId) {
        setPendingSwitchDeviceId(nextDeviceId);
        setShowSwitchConfirmDialog(true);
        return;
      }
      void applyDeviceSwitch(nextDeviceId);
    },
    [applyDeviceSwitch, hasUnsavedChanges, searchParams, selectedDeviceId, setSearchParams],
  );
  const handleToggleSidebar = useCallback(() => {
    setActiveTab("devices");
    const next = new URLSearchParams(searchParams);
    next.set("modal", "devices");
    setSearchParams(next, { replace: true });
  }, [searchParams, setSearchParams]);

  const handleImportPoints = useCallback(
    async (points: CreatePointRequest[]) => {
      const createdIds: string[] = [];
      for (const point of points) {
        const result = await createPointMutation.mutateAsync(point);
        if (result?.id) createdIds.push(result.id);
      }
      history.push({
        type: "import",
        description: t("smartDashboard.importedPoints", { count: points.length }),
        data: { pointIds: createdIds },
      });
    },
    [createPointMutation, history, t],
  );

  const handleUndo = useCallback(() => {
    history.undo();
  }, [history]);

  const handleRedo = useCallback(() => {
    history.redo();
  }, [history]);

  const handleImportShortcut = useCallback(() => {
    if (selectedDeviceId) setImportDialogOpen(true);
  }, [selectedDeviceId]);

  const handleExportShortcut = useCallback(() => {
    if (selectedDeviceId && allPoints.length > 0) setExportDialogOpen(true);
  }, [selectedDeviceId, allPoints.length]);

  const handleSearchShortcut = useCallback(() => {
    searchInputRef.current?.focus();
  }, []);

  const handleSaveTemplate = useCallback(() => {
    const normalized = templateName.trim();
    if (!normalized || !selectedDevice || !typedPlanValidation.valid) return;
    const next = createTemplateFromPlanner({
      templateName: normalized,
      draft: {
        dataType: planDataType,
        count: planCount,
        startAddress: planStartAddress,
      },
    });

    setSourceTemplates((prev) => upsertTemplateRecord(prev, next));
    setTemplateName("");
  }, [planCount, planDataType, planStartAddress, selectedDevice, templateName, typedPlanValidation.valid]);

  const handleLoadTemplate = useCallback((template: SourceTemplate) => {
    const applied = applyTemplateToPlanner(template);
    setPlanDataType(applied.dataType);
    setPlanCount(applied.count);
    setPlanStartAddress(applied.startAddress);
    const now = new Date().toISOString();
    setSourceTemplates((prev) =>
      prev.map((item) =>
        item.id === template.id
          ? {
              ...item,
              lastUsedAt: now,
              version: Math.max(item.version, SOURCE_TEMPLATE_SCHEMA_VERSION),
            }
          : item,
      ),
    );
  }, []);

  const handleUpgradeTemplates = useCallback(() => {
    setSourceTemplates((prev) => upgradeTemplates(prev));
  }, []);

  const handleDeleteTemplate = useCallback((templateId: string) => {
    setSourceTemplates((prev) => prev.filter((item) => item.id !== templateId));
  }, []);

  const handleAutoAllocate = useCallback(() => {
    if (!selectedDevice || !planStartAddress || totalPlannedCells <= 0) return;
    const result = findNearestValidContiguousSpan({
      startAddress: planStartAddress,
      spanSize: totalPlannedCells,
      protocol: selectedDevice.protocol,
      blockedAddresses,
    });

    if (result.startAddress) {
      setPlanStartAddress(result.startAddress);
      setAllocationMessage(`已自動配置到 ${result.startAddress}`);
      return;
    }

    if (result.reason === "invalid_start") {
      setAllocationMessage("起始位址格式無效，請先修正後再試一次。");
      return;
    }

    setAllocationMessage("在目前範圍內找不到可用連續區段，請調整起始位址或降低來源數量。");
  }, [blockedAddresses, planStartAddress, selectedDevice, totalPlannedCells]);

  const scheduleGuideStageReset = useCallback(() => {
    if (guideStageTimeoutRef.current) window.clearTimeout(guideStageTimeoutRef.current);
    guideStageTimeoutRef.current = window.setTimeout(() => {
      setGuideStage("idle");
      guideStageTimeoutRef.current = null;
    }, 1200);
  }, []);

  const handleApplyPlan = useCallback(() => {
    if (!planAddresses.length) return;
    setSelectedAddresses(planAddresses);
    setGuideStage("grid");
    gridSectionRef.current?.scrollIntoView({ behavior: resolveScrollBehavior(reducedMotion), block: "start" });
    gridSectionRef.current?.focus();
    scheduleGuideStageReset();
    setPanelType("batch");
    setBatchInitialTemplate(`${(batchNamePrefix.trim() || "SRC").toUpperCase()}-{index03}`);
  }, [planAddresses, batchNamePrefix, reducedMotion, scheduleGuideStageReset]);

  /** 批量建立點位成功後：依點位命名自動建立 Tag 並建立 mapping，再切到 Tag 分頁 */
  const handleBatchCreatedWithPoints = useCallback(
    async (points: Point[]) => {
      if (!points.length) {
        setSidebarTab("tag");
        return;
      }
      let tagCount = 0;
      let linkCount = 0;
      try {
        for (const point of points) {
          const name = (point.name || "").trim();
          if (!name) continue;
          let tagId: string | null = null;
          try {
            const tag = await createTagMutation.mutateAsync({
              key: name,
              display_name: name,
              data_type: point.data_type,
            });
            tagId = tag.id;
            tagCount += 1;
          } catch {
            await queryClient.invalidateQueries({ queryKey: tagKeys.lists() });
            const list = await queryClient.fetchQuery({
              queryKey: tagKeys.list(),
              queryFn: () => tagAPI.list(),
            });
            const existing = (list as Tag[]).find((t) => t.key === name);
            if (existing) tagId = existing.id;
          }
          if (tagId) {
            await createMappingMutation.mutateAsync({
              point_id: point.id,
              tag_id: tagId,
              enabled: true,
            });
            linkCount += 1;
          }
        }
        setSidebarTab("tag");
        if (tagCount > 0 || linkCount > 0) {
          showSuccess(`已建立 ${points.length} 個點位，並依命名建立 ${tagCount} 個 Tag、連結 ${linkCount} 筆映射。`);
        }
      } catch (err) {
        showError("依命名建立 Tag 時發生錯誤: " + (err instanceof Error ? err.message : "未知錯誤"));
        setSidebarTab("tag");
      }
    },
    [
      createTagMutation,
      createMappingMutation,
      queryClient,
      setSidebarTab,
      showSuccess,
      showError,
    ]
  );

  const getGridCenterAddress = useCallback((protocol: ProtocolType) => (protocol.startsWith("modbus") ? "40001" : "D0"), []);

  const primarySelectedAddress = selectedAddresses[0] || "";
  const selectedSourceAddress = primarySelectedAddress || selectedPoint?.address || "";
  const selectedPointFromGrid = useMemo(
    () => allPoints.find((point) => point.address === selectedSourceAddress) || null,
    [allPoints, selectedSourceAddress],
  );
  const activePointForLink = selectedPoint ?? selectedPointFromGrid;

  /** Popover 開啟且為空格位時，以點擊格位為 Tag 面板的選取來源，否則用側欄/網格選取 */
  const effectiveSelectedSourceAddress =
    gridPopoverOpen && !gridPopoverPoint ? gridPopoverAddress : selectedSourceAddress;
  const effectiveActivePointForLink =
    gridPopoverOpen && !gridPopoverPoint
      ? allPoints.find((p) => p.address === gridPopoverAddress) ?? null
      : activePointForLink;

  const selectedMapping = useMemo(
    () =>
      mappings.find(
        (mapping) => effectiveActivePointForLink && mapping.point_id === effectiveActivePointForLink.id
      ) || null,
    [effectiveActivePointForLink, mappings],
  );
  const linkedTag = useMemo(
    () => tags.find((tag) => tag.id === selectedMapping?.tag_id) || null,
    [selectedMapping?.tag_id, tags],
  );
  const linkedTagAffectedMappingsCount = useMemo(() => {
    if (!linkedTag?.id) return 0;
    return getAffectedMappingCountForTag(mappings, linkedTag.id);
  }, [linkedTag?.id, mappings]);
  const tagEditGuardrail = useMemo(
    () => buildGlobalTagGuardrail(linkedTagAffectedMappingsCount),
    [linkedTagAffectedMappingsCount],
  );
  const parsePipeline = useCallback(() => {
    if (!selectedMapping?.transform_pipeline) return [];
    try {
      const parsed = JSON.parse(selectedMapping.transform_pipeline);
      return Array.isArray(parsed) ? parsed : [];
    } catch {
      return [];
    }
  }, [selectedMapping?.transform_pipeline]);
  const {
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
  } = useSmartDashboardTagLinking({
    activePointForLink: effectiveActivePointForLink,
    selectedMapping: selectedMapping || undefined,
    tags,
    linkedTag: linkedTag || undefined,
    linkedTagAffectedMappingsCount,
    parsePipeline,
    tagEditGuardrailWarning: tagEditGuardrail.warningMessage,
    createTag: createTagMutation.mutateAsync,
    createMapping: createMappingMutation.mutateAsync,
    updateMapping: updateMappingMutation.mutateAsync,
    updateTag: updateTagMutation.mutateAsync,
    refetchMappings,
  });
  const {
    modbusStatus,
    modbusRegister,
    setModbusRegister,
    loadModbusStatus,
    handleBindTagToModbus,
    handlePushCurrentValueToModbus,
    handleSyncModbusFromMappings,
  } = useSmartDashboardModbusActions({
    linkedTag: linkedTag || undefined,
    selectedPoint,
    showError,
    showInfo,
    showSuccess,
  });

  /** Tag 面板 props，供側邊欄與網格 Popover 共用（Popover 開啟時用點擊格位為來源） */
  const tagPanelProps = useMemo(
    () => ({
      selectedSourceAddress: effectiveSelectedSourceAddress,
      activePointForLink: effectiveActivePointForLink ?? null,
      linkedTag: linkedTag ?? null,
      tagLinkMode,
      setTagLinkMode,
      selectedTagIdForLink,
      setSelectedTagIdForLink,
      tags,
      handleLinkTagToSelectedAddress,
      updateMappingPending: updateMappingMutation.isPending,
      createMappingPending: createMappingMutation.isPending,
      newTagKey,
      setNewTagKey,
      newTagDisplayName,
      setNewTagDisplayName,
      handleCreateTagAndLink,
      createTagPending: createTagMutation.isPending,
      tagLinkActionMessage,
      linkedTagAffectedMappingsCount,
      tagEditDisplayName,
      setTagEditDisplayName,
      tagEditUnit,
      setTagEditUnit,
      tagEditDescription,
      setTagEditDescription,
      handleSaveLinkedTagEdit,
      updateTagPending: updateTagMutation.isPending,
      pendingTagEdit,
      handleConfirmTagEdit,
      clearPendingTagEdit,
      tagEditMessage,
    }),
    [
      effectiveSelectedSourceAddress,
      effectiveActivePointForLink,
      linkedTag,
      tagLinkMode,
      setTagLinkMode,
      selectedTagIdForLink,
      setSelectedTagIdForLink,
      tags,
      handleLinkTagToSelectedAddress,
      updateMappingMutation.isPending,
      createMappingMutation.isPending,
      newTagKey,
      setNewTagKey,
      newTagDisplayName,
      setNewTagDisplayName,
      handleCreateTagAndLink,
      createTagMutation.isPending,
      tagLinkActionMessage,
      linkedTagAffectedMappingsCount,
      tagEditDisplayName,
      setTagEditDisplayName,
      tagEditUnit,
      setTagEditUnit,
      tagEditDescription,
      setTagEditDescription,
      handleSaveLinkedTagEdit,
      updateTagMutation.isPending,
      pendingTagEdit,
      handleConfirmTagEdit,
      clearPendingTagEdit,
      tagEditMessage,
    ]
  );

  const segmentFeedback = useMemo(() => {
    const sourceOk = Boolean(selectedDeviceId && selectedSourceAddress);
    const gridOk = planAddresses.length > 0 && planConflictCount === 0;
    const tagOk = Boolean(selectedMapping && linkedTag);
    const sinkOk = canActivate;

    return [
      {
        id: "source",
        label: t("smartDashboard.commit.segment.source.label"),
        ok: sourceOk,
        message: sourceOk
          ? t("smartDashboard.commit.segment.source.ready")
          : t("smartDashboard.commit.segment.source.missing"),
      },
      {
        id: "grid",
        label: t("smartDashboard.commit.segment.grid.label"),
        ok: gridOk,
        message: gridOk
          ? t("smartDashboard.commit.segment.grid.ready")
          : t("smartDashboard.commit.segment.grid.missing", { count: planConflictCount }),
      },
      {
        id: "tag",
        label: t("smartDashboard.commit.segment.tag.label"),
        ok: tagOk,
        message: tagOk
          ? t("smartDashboard.commit.segment.tag.ready", { key: linkedTag?.key ?? "" })
          : t("smartDashboard.commit.segment.tag.missing"),
      },
      {
        id: "sink",
        label: t("smartDashboard.commit.segment.sink.label"),
        ok: sinkOk,
        message: sinkOk
          ? t("smartDashboard.commit.segment.sink.ready")
          : t("smartDashboard.commit.segment.sink.missing"),
      },
    ] as const;
  }, [
    canActivate,
    linkedTag,
    planAddresses.length,
    planConflictCount,
    selectedDeviceId,
    selectedMapping,
    selectedSourceAddress,
    t,
  ]);
  const {
    commitQueueItems,
    commitQueueSummary,
    commitImpactSummary,
    preCommitLoadEstimate,
    motionQAGate,
    commitActionMessage,
    setCommitActionMessage,
    isCommitRunning,
    commitChunkResults,
    commitAuditPayload,
    hasFailedChunk,
    canRollback,
    handleCommitFlow,
    handleRetryFailedCommits,
    handleRollbackCommitRun,
  } = useSmartDashboardCommitFlow({
    allPoints,
    linkedAddresses,
    plannedAllocations,
    pendingGlobalTagEdit: Boolean(pendingTagEdit),
    pollingGroups,
    canActivate,
    mappingEnabled: Boolean(selectedMapping?.enabled),
    markActive,
    markError,
    setGuideStage,
    scheduleGuideStageReset,
    t,
  });

  const dismissLegacyNotice = useCallback(() => {
    if (!isLegacyDecommissionRoute(legacyRoute)) return;
    const next = new URLSearchParams(searchParams);
    next.delete("legacy");
    setSearchParams(next, { replace: true });
  }, [legacyRoute, searchParams, setSearchParams]);
  const dismissSectionIntentNotice = useCallback(() => {
    if (!sectionIntent) return;
    const next = new URLSearchParams(searchParams);
    next.delete("section");
    setSearchParams(next, { replace: true });
  }, [searchParams, sectionIntent, setSearchParams]);
  const closeWorkflowModal = useCallback(() => {
    if (!modalIntent) return;
    const next = new URLSearchParams(searchParams);
    next.delete("modal");
    setSearchParams(next, { replace: true });
  }, [modalIntent, searchParams, setSearchParams]);
  const goToLocalModbusWorkbench = useCallback(() => {
    const section = activeTab;
    navigate(
      buildWorkbenchRedirect({
        step: "output",
        target: "modbus",
        section,
      }),
    );
  }, [activeTab, navigate]);
  const openDeviceSetupModal = useCallback(
    (deviceId?: string | null) => {
      if (!deviceId) return;
      const target = devices.find((device) => device.id === deviceId) || null;
      if (!target) return;
      setEditingDeviceInModal(target);
      setActiveTab("devices");
      const next = new URLSearchParams(searchParams);
      next.set("modal", "devices");
      setSearchParams(next, { replace: true });
    },
    [devices, searchParams, setSearchParams],
  );
  const handleSubmitDeviceSetup = useCallback(
    async (data: CreateDeviceRequest | UpdateDeviceRequest) => {
      if (!editingDeviceInModal) return;
      await updateDeviceMutation.mutateAsync({
        id: editingDeviceInModal.id,
        data: data as UpdateDeviceRequest,
      });
      showSuccess(`已更新設備「${editingDeviceInModal.name}」設定`);
    },
    [editingDeviceInModal, showSuccess, updateDeviceMutation],
  );
  const handleTestDeviceConnection = useCallback(
    async (deviceId: string) => {
      setTestingDeviceId(deviceId);
      try {
        const result = await testConnectionMutation.mutateAsync(deviceId);
        if (result.success) {
          showSuccess(`連線成功，延遲 ${result.latency_ms} ms`);
          return;
        }
        showError(`連線失敗：${result.error || "未知錯誤"}`);
      } catch (error) {
        const message = error instanceof Error ? error.message : "測試連線失敗";
        showError(`連線失敗：${message}`);
      } finally {
        setTestingDeviceId(null);
      }
    },
    [showError, showSuccess, testConnectionMutation],
  );
  const handleChooseDevice = useCallback(() => {
    setActiveTab("devices");
    const next = new URLSearchParams(searchParams);
    next.set("modal", "devices");
    setSearchParams(next, { replace: true });
  }, [searchParams, setSearchParams]);
  const handleToggleDeviceStatusDirect = useCallback(
    async (deviceId: string, closeModalOnSuccess = false) => {
      const target = devices.find((device) => device.id === deviceId);
      if (!target) return;
      const willActivate = target.status !== "active";
      setActivatingDeviceId(deviceId);
      try {
        await toggleDeviceStatusMutation.mutateAsync({
          id: deviceId,
          currentStatus: target.status,
        });
        showSuccess(`已${willActivate ? "啟用" : "停用"}設備「${target.name}」`);
        if (willActivate) {
          setSelectedDeviceId(deviceId);
        } else if (selectedDeviceId === deviceId) {
          setSelectedAddresses([]);
        }
        if (closeModalOnSuccess) {
          const next = new URLSearchParams(searchParams);
          if (next.get("modal") === "devices") {
            next.delete("modal");
            setSearchParams(next, { replace: true });
          }
        }
      } catch (error) {
        const message = error instanceof Error ? error.message : `${willActivate ? "啟用" : "停用"}失敗`;
        showError(`設備${willActivate ? "啟用" : "停用"}失敗：${message}`);
      } finally {
        setActivatingDeviceId(null);
      }
    },
    [devices, searchParams, selectedDeviceId, setSearchParams, showError, showSuccess, toggleDeviceStatusMutation],
  );
  /** 點選設備卡片時：切換到該設備（draft 先啟用再切換）、關閉 modal 並保持連接 */
  const handleSelectDevice = useCallback(
    (device: Device) => {
      if (device.id === selectedDeviceId) {
        const next = new URLSearchParams(searchParams);
        if (next.get("modal") === "devices") {
          next.delete("modal");
          setSearchParams(next, { replace: true });
        }
        return;
      }
      if (device.status === "draft") {
        void handleToggleDeviceStatusDirect(device.id, true);
        return;
      }
      requestDeviceSwitch(device.id);
    },
    [
      handleToggleDeviceStatusDirect,
      requestDeviceSwitch,
      searchParams,
      selectedDeviceId,
      setSearchParams,
    ],
  );
  const requestDeleteDevice = useCallback(
    (deviceId: string) => {
      const target = devices.find((device) => device.id === deviceId) || null;
      setDeleteConfirmDevice(target);
    },
    [devices],
  );
  const handleConfirmDeleteDevice = useCallback(async () => {
    if (!deleteConfirmDevice) return;
    setDeletingDeviceId(deleteConfirmDevice.id);
    try {
      await deleteDeviceMutation.mutateAsync(deleteConfirmDevice.id);
      showSuccess(`已刪除設備「${deleteConfirmDevice.name}」`);
      if (selectedDeviceId === deleteConfirmDevice.id) {
        setSelectedDeviceId(null);
        setSelectedAddresses([]);
      }
      if (editingDeviceInModal?.id === deleteConfirmDevice.id) {
        setEditingDeviceInModal(null);
      }
      setDeleteConfirmDevice(null);
    } catch (error) {
      const message = error instanceof Error ? error.message : "刪除設備失敗";
      showError(`刪除設備失敗：${message}`);
    } finally {
      setDeletingDeviceId(null);
    }
  }, [deleteConfirmDevice, deleteDeviceMutation, editingDeviceInModal, selectedDeviceId, showError, showSuccess]);
  const handleCreateDevice = useCallback(() => {
    setIsCreateDeviceModalOpen(true);
    setActiveTab("devices");
    setJustCreatedDeviceId(null);
  }, []);
  const closeCreateDeviceModal = useCallback(() => {
    setIsCreateDeviceModalOpen(false);
    if (createDeviceIntent === "1") {
      const next = new URLSearchParams(searchParams);
      next.delete("createDevice");
      setSearchParams(next, { replace: true });
    }
  }, [createDeviceIntent, searchParams, setSearchParams]);

  useEffect(() => {
    if (!sectionIntent) return;
    setActiveTab(sectionIntent);
  }, [sectionIntent]);
  useEffect(() => {
    if (!modalIntent) return;
    if (modalIntent === "devices") setActiveTab("devices");
    if (modalIntent === "settings") setActiveTab("settings");
    if (modalIntent === "wizard") setIsCreateDeviceModalOpen(true);
  }, [modalIntent]);
  useEffect(() => {
    if (modalIntent !== "devices") {
      setEditingDeviceInModal(null);
    }
  }, [modalIntent]);
  useEffect(() => {
    if (createDeviceIntent === "1") {
      setActiveTab("devices");
      setIsCreateDeviceModalOpen(true);
    }
  }, [createDeviceIntent]);
  // 僅追蹤 id 變化，避免 setEditingDeviceInModal 觸發自身重複執行
  useEffect(() => {
    const editingId = editingDeviceInModal?.id;
    if (!editingId) return;
    const refreshed = devices.find((device) => device.id === editingId);
    if (!refreshed) {
      setEditingDeviceInModal(null);
      return;
    }
    setEditingDeviceInModal(refreshed);
  }, [devices, editingDeviceInModal?.id]);

  useEffect(() => {
    if (selectedDeviceId || devices.length === 0) return;
    setActiveTab("devices");
  }, [devices.length, selectedDeviceId]);
  const confirmSwitchToCreatedDevice = useCallback(
    (switchNow: boolean) => {
      if (switchNow && justCreatedDeviceId) {
        void applyDeviceSwitch(justCreatedDeviceId);
      }
      setJustCreatedDeviceId(null);
    },
    [applyDeviceSwitch, justCreatedDeviceId],
  );

  useEffect(() => {
    setSource(selectedDeviceId || "", selectedSourceAddress, selectedPoint?.id || "");
  }, [selectedDeviceId, selectedPoint?.id, selectedSourceAddress, setSource]);

  /** 僅在設備變更時初始化 grid 視窗起始位址；修改起始位址不聯動 grid，由「套用到網格」套用 */
  useEffect(() => {
    if (!selectedDevice) return;
    setGridViewStartAddress(getGridCenterAddress(selectedDevice.protocol));
  }, [selectedDevice, getGridCenterAddress]);

  const handleGridViewShift = useCallback(
    (delta: number) => {
      if (!selectedDevice) return;
      const current = gridViewStartAddress || planStartAddress || getGridCenterAddress(selectedDevice.protocol);
      const next = addressParser.offset(current, delta, selectedDevice.protocol);
      setGridViewStartAddress(next);
    },
    [gridViewStartAddress, planStartAddress, selectedDevice, getGridCenterAddress]
  );

  useEffect(() => {
    setTag(linkedTag?.id || "");
  }, [linkedTag?.id, setTag]);

  useEffect(() => {
    setSink(selectedMapping ? "timeseries" : "");
  }, [selectedMapping, setSink]);

  useEffect(() => {
    const sourceQuality = selectedPoint?.last_error ? "bad" : selectedDevice ? "good" : "unknown";
    const sourceValue =
      selectedPoint?.last_value !== undefined
        ? String(selectedPoint.last_value)
        : selectedDevice
          ? selectedDevice.name
          : "-";
    const sourceTime = selectedPoint?.last_read_at || "-";
    const sourceError = selectedPoint?.last_error || "";
    const gridQuality = selectedSourceAddress ? "good" : "unknown";
    const tagQuality = linkedTag ? (selectedMapping?.enabled ? "good" : "warning") : "unknown";
    const sinkQuality = selectedPoint?.last_error ? "bad" : selectedMapping?.enabled ? "good" : "unknown";

    setDiagnostics({
      source: {
        latestValue: sourceValue,
        quality: sourceQuality,
        timestamp: sourceTime,
        error: sourceError,
      },
      grid: {
        latestValue: selectedSourceAddress || "-",
        quality: gridQuality,
        timestamp: sourceTime,
        error: "",
      },
      tag: {
        latestValue: linkedTag ? `${linkedTag.key}` : "-",
        quality: tagQuality,
        timestamp: linkedTag?.updated_at || "-",
        error: "",
      },
      sink: {
        latestValue: flowState.sinkTarget || "-",
        quality: sinkQuality,
        timestamp: selectedMapping?.updated_at || "-",
        error: selectedPoint?.last_error || "",
      },
    });
  }, [
    flowState.sinkTarget,
    linkedTag,
    selectedMapping?.enabled,
    selectedMapping?.updated_at,
    selectedDevice,
    selectedSourceAddress,
    selectedPoint?.last_error,
    selectedPoint?.last_read_at,
    selectedPoint?.last_value,
    setDiagnostics,
  ]);

  const handleValidateFlow = useCallback(async () => {
    const structural = runStructuralValidation({
      selectedDeviceId,
      selectedSourceAddress,
      planAddressesCount: planAddresses.length,
      planConflictCount,
      hasSelectedMapping: Boolean(selectedMapping),
    });
    if (!structural.ok) {
      const structuralError = structural.error || t("smartDashboard.commit.validation.structureErrorDefault");
      if (structuralError.includes("來源")) {
        markError("source", t("smartDashboard.flowErrors.missingSource"));
      } else if (structuralError.includes("Tag")) {
        markError("tag", t("smartDashboard.flowErrors.missingMapping"));
      } else {
        markError("grid", structuralError);
      }
      setCommitActionMessage(t("smartDashboard.commit.validation.structureFailed", { error: structuralError }));
      return;
    }
    setCommitActionMessage(t("smartDashboard.commit.validation.structurePassed"));
    try {
      const pipeline = parsePipeline();
      const result = await validatePipelineMutation.mutateAsync(pipeline);
      if (result.valid) {
        markValidated();
        setCommitActionMessage(t("smartDashboard.commit.validation.validationPassed"));
        return;
      }
      markError("grid", result.error || t("smartDashboard.flowErrors.validationFailed"));
      setCommitActionMessage(
        t("smartDashboard.commit.validation.executeFailed", {
          error: result.error || t("smartDashboard.commit.validation.executeFailedDefault"),
        })
      );
    } catch (error) {
      const message = error instanceof Error ? error.message : t("smartDashboard.flowErrors.validationFailed");
      markError("grid", message);
      setCommitActionMessage(t("smartDashboard.commit.validation.executeFailed", { error: message }));
    }
  }, [
    markError,
    markValidated,
    parsePipeline,
    planAddresses.length,
    planConflictCount,
    selectedDeviceId,
    selectedMapping,
    selectedSourceAddress,
    setCommitActionMessage,
    t,
    validatePipelineMutation,
  ]);

  const handleRecoverFlow = useCallback(() => {
    resetDraft();
    setDiagnostics({
      source: { ...RESET_SEGMENT_DIAGNOSTIC },
      grid: { ...RESET_SEGMENT_DIAGNOSTIC },
      tag: { ...RESET_SEGMENT_DIAGNOSTIC },
      sink: { ...RESET_SEGMENT_DIAGNOSTIC },
    });
  }, [resetDraft, setDiagnostics]);

  const shortcuts = useSmartDashboardShortcuts({
    onBatchCreate: handleBatchCreate,
    onClosePanel: handleClosePanel,
    onToggleSidebar: handleToggleSidebar,
    onUndo: handleUndo,
    onRedo: handleRedo,
    onSearch: handleSearchShortcut,
    onValidateFlow: handleValidateFlow,
    onActivateFlow: handleCommitFlow,
    onRecoverFlow: hasError ? handleRecoverFlow : undefined,
    onImport: handleImportShortcut,
    onExport: handleExportShortcut,
  });
  const sidebarPanelMotion = useSmartDashboardSidebarPanelMotion({ guideStage, reducedMotion });
  const workspaceSectionState = useSmartDashboardWorkspaceState({
    workspace: {
      selectedDevice,
      planDataType,
      setPlanDataType,
      planCount,
      setPlanCount,
      totalPlannedCells,
      planStartAddress,
      setPlanStartAddress,
      handleAutoAllocate,
      handleApplyPlan,
      planConflictCount,
      typedPlanValidation,
      showConflictsOnly,
      flowSegments: FLOW_SEGMENTS,
      flowState,
      statusStyle: STATUS_STYLE,
      hasError,
      t,
      gridSectionRef,
      resolveIntentMotionClass,
      guideStage,
      reducedMotion,
      motionTokens: MOTION_TOKENS,
      modbusStatus,
      allPoints,
      linkedAddresses,
      selectedAddresses,
      plannedAllocations,
      setSelectedAddresses,
      handleCellClick,
      onCellContextMenu: handleCellContextMenu,
      getGridCenterAddress,
      gridViewStartAddress,
      onGridViewShift: handleGridViewShift,
      handleChooseDevice,
      handleCreateDevice,
    },
    sidebar: {
      sidebarTab,
      setSidebarTab,
      hasTagSelection: Boolean(selectedSourceAddress),
      planningTabProps: {
        planDataType,
        planCount,
        totalPlannedCells,
        planConflictCount,
        typedPlanValidation,
        allocationMessage,
        templateName,
        setTemplateName,
        handleSaveTemplate,
        sourceTemplates,
        handleLoadTemplate,
        handleDeleteTemplate,
        staleTemplateCount,
        sourceTemplateSchemaVersion: SOURCE_TEMPLATE_SCHEMA_VERSION,
        handleUpgradeTemplates,
        batchNamePrefix,
        setBatchNamePrefix,
        normalizeNamingPrefix,
        namePreview,
        nameConflictCount,
        showConflictsOnly,
        setShowConflictsOnly,
        t,
      },
      tagPanelProps,
      modbusPanelProps: {
        goToLocalModbusWorkbench,
        loadModbusStatus,
        modbusStatus,
        modbusRegister,
        setModbusRegister,
        handleBindTagToModbus,
        handlePushCurrentValueToModbus,
        handleSyncModbusFromMappings,
      },
      commitPanelProps: {
        commitQueueSummary,
        commitImpactSummary,
        preCommitLoadEstimate,
        motionQAGate,
        commitQueueItems,
        onValidateFlow: handleValidateFlow,
        canValidate,
        validating: validatePipelineMutation.isPending,
        onCommitFlow: handleCommitFlow,
        canCommit: canActivate,
        isCommitRunning,
        onRetryFailed: handleRetryFailedCommits,
        hasFailedChunk,
        onRollback: handleRollbackCommitRun,
        canRollback,
        segmentFeedback,
        commitActionMessage,
        commitAuditPayload,
        commitChunkResults,
        hasError,
        onRecoverFlow: handleRecoverFlow,
        t,
      },
      selectedDevice,
      onOpenWorkbench: goToLocalModbusWorkbench,
      onOpenImport: () => setImportDialogOpen(true),
      onOpenExport: () => setExportDialogOpen(true),
      canExport: allPoints.length > 0,
      onUndo: handleUndo,
      onRedo: handleRedo,
      canUndo: history.canUndo,
      canRedo: history.canRedo,
      undoDescription: history.getUndoAction()?.description || t("smartDashboard.noUndo"),
      redoDescription: history.getRedoAction()?.description || t("smartDashboard.noRedo"),
      onOpenShortcuts: () => setPanelType("shortcuts"),
      t,
    },
    sidebarPanelMotion,
  });

  return (
    <div className="flex flex-col min-h-[calc(100vh-11rem)] bg-gradient-to-br from-[#0B0F19] via-[#111827] to-[#0F172A] text-slate-100 font-sans rounded-2xl overflow-hidden selection:bg-blue-500/30">
      <SmartDashboardHeader
        searchInputRef={searchInputRef}
        productSubtitle={t("smartDashboard.productSubtitle")}
        searchLabel={t("smartDashboard.searchLabel")}
        searchPlaceholder={t("smartDashboard.searchPlaceholder")}
        notificationsLabel={t("smartDashboard.notifications")}
        preferencesLabel={t("smartDashboard.preferences")}
        profileMenuLabel={t("smartDashboard.profileMenu")}
        onOpenLocalModbusWorkbench={goToLocalModbusWorkbench}
      />

      <SmartDashboardIntentNotices
        showLegacyNotice={Boolean(legacyRoute && legacyRouteLabel)}
        legacyMigrationTitle={t("smartDashboard.legacyMigration.title")}
        legacyMigrationDescription={t("smartDashboard.legacyMigration.description", { route: legacyRouteLabel })}
        legacyMigrationDismissLabel={t("smartDashboard.legacyMigration.dismiss")}
        dismissLegacyNotice={dismissLegacyNotice}
        showSectionIntentNotice={Boolean(sectionIntent && sectionIntentLabel)}
        sectionIntentLabel={sectionIntentLabel}
        dismissSectionIntentNotice={dismissSectionIntentNotice}
        showModalIntentNotice={Boolean(modalIntent && modalIntentLabel)}
        modalIntentLabel={modalIntentLabel}
        closeWorkflowModal={closeWorkflowModal}
        t={t}
      />
      <SmartDashboardControlBar onChooseDevice={handleChooseDevice} onCreateDevice={handleCreateDevice} t={t} />
      <SmartDashboardWorkspaceSection {...workspaceSectionState} />
      <SmartDashboardWorkflowModal
        modalIntent={modalIntent}
        modalIntentLabel={modalIntentLabel}
        closeWorkflowModal={closeWorkflowModal}
        deviceSearchQuery={deviceSearchQuery}
        setDeviceSearchQuery={setDeviceSearchQuery}
        deviceStatusFilter={deviceStatusFilter}
        setDeviceStatusFilter={setDeviceStatusFilter}
        handleCreateDevice={handleCreateDevice}
        deviceSummary={deviceSummary}
        filteredDevices={filteredDevices}
        openDeviceSetupModal={openDeviceSetupModal}
        handleTestDeviceConnection={handleTestDeviceConnection}
        testingDeviceId={testingDeviceId}
        handleToggleDeviceStatusDirect={handleToggleDeviceStatusDirect}
        activatingDeviceId={activatingDeviceId}
        requestDeleteDevice={requestDeleteDevice}
        deletingDeviceId={deletingDeviceId}
        requestDeviceSwitch={requestDeviceSwitch}
        handleSelectDevice={handleSelectDevice}
        isSwitchingDevice={isSwitchingDevice}
        selectedDeviceId={selectedDeviceId}
        editingDeviceInModal={editingDeviceInModal}
        setEditingDeviceInModal={setEditingDeviceInModal}
        handleSubmitDeviceSetup={handleSubmitDeviceSetup}
        setPanelType={setPanelType}
        setActiveTab={setActiveTab}
        goToLocalModbusWorkbench={goToLocalModbusWorkbench}
        t={t}
      />

      <SmartDashboardOverlays
        showSwitchConfirmDialog={showSwitchConfirmDialog}
        setShowSwitchConfirmDialog={setShowSwitchConfirmDialog}
        pendingSwitchDeviceId={pendingSwitchDeviceId}
        applyDeviceSwitch={applyDeviceSwitch}
        setSelectedAddresses={setSelectedAddresses}
        setPendingSwitchDeviceId={setPendingSwitchDeviceId}
        deleteConfirmDevice={deleteConfirmDevice}
        setDeleteConfirmDevice={setDeleteConfirmDevice}
        deletingDeviceId={deletingDeviceId}
        handleConfirmDeleteDevice={handleConfirmDeleteDevice}
        isCreateDeviceModalOpen={isCreateDeviceModalOpen}
        closeCreateDeviceModal={closeCreateDeviceModal}
        setJustCreatedDeviceId={setJustCreatedDeviceId}
        justCreatedDeviceId={justCreatedDeviceId}
        confirmSwitchToCreatedDevice={confirmSwitchToCreatedDevice}
        t={t}
      />
      <SmartDashboardPanels
        panelType={panelType}
        setPanelType={setPanelType}
        selectedDevice={selectedDevice}
        selectedPoint={selectedPoint}
        selectedAddresses={selectedAddresses}
        setSelectedAddresses={setSelectedAddresses}
        pollingGroups={pollingGroups}
        shortcuts={shortcuts}
        batchCreateTitle={t("smartDashboard.batchCreate")}
        shortcutsTitle={t("smartDashboard.shortcuts")}
        pointDetailTitle={t("smartDashboard.pointDetail")}
        shortcutsHint={t("smartDashboard.shortcutsHint")}
        initialBatchTemplate={batchInitialTemplate ?? undefined}
        initialBatchDataType={batchInitialTemplate != null ? planDataType : undefined}
        onBatchClose={() => setBatchInitialTemplate(null)}
        onBatchCreated={handleBatchCreatedWithPoints}
      />

      <SmartDashboardGridOverlaysSection
        gridPopoverOpen={gridPopoverOpen}
        setGridPopoverOpen={setGridPopoverOpen}
        gridPopoverPosition={gridPopoverPosition}
        gridPopoverAddress={gridPopoverAddress}
        gridPopoverPoint={gridPopoverPoint}
        gridPopoverContentRef={gridPopoverContentRef}
        tagPanelProps={tagPanelProps}
        queryClient={queryClient}
        deletePoint={deletePointMutation.mutateAsync}
        selectedPoint={selectedPoint}
        setSelectedPoint={setSelectedPoint}
        showSuccess={showSuccess}
        showError={showError}
        gridContextMenu={gridContextMenu}
        gridContextMenuRef={gridContextMenuRef}
        setGridContextMenu={setGridContextMenu}
        t={t}
      />

      {selectedDevice && (
        <>
          <ImportDialog
            open={importDialogOpen}
            onOpenChange={setImportDialogOpen}
            deviceId={selectedDevice.id}
            deviceName={selectedDevice.name}
            onImport={handleImportPoints}
          />
          <ExportDialog
            open={exportDialogOpen}
            onOpenChange={setExportDialogOpen}
            points={allPoints}
            deviceId={selectedDevice.id}
            deviceName={selectedDevice.name}
          />
        </>
      )}
    </div>
  );
}
