import { useState, useMemo, useCallback, useEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import type { Point, CreatePointRequest, ProtocolType, DataType } from '../../types/datalink';
import { DeviceTreeNav } from '../../components/datalink/DeviceTreeNav';
import { MemoryGrid, type PlannedAllocation } from '../../components/datalink/MemoryGrid';
import { QuickActions } from '../../components/datalink/QuickActions';
import { SlidePanel } from '../../components/datalink/SlidePanel';
import { BatchPointCreator } from '../../components/datalink/BatchPointCreator';
import { PointDetailPanel } from '../../components/datalink/PointDetailPanel';
import { ImportDialog, ExportDialog } from '../../components/datalink/ImportExportDialog';
import { useDevicesQuery } from '../../hooks/datalink/useDevices';
import { usePollingGroupsQuery } from '../../hooks/datalink/usePollingGroups';
import { usePointsQuery, useCreatePointMutation } from '../../hooks/datalink/usePoints';
import {
  useMappingsQuery,
  useValidatePipelineMutation,
  useCreateMappingMutation,
  useUpdateMappingMutation,
} from '../../hooks/datalink/useMappings';
import { useTagsQuery, useCreateTagMutation, useUpdateTagMutation } from '../../hooks/datalink/useTags';
import { useSmartDashboardShortcuts } from '../../hooks/useKeyboardShortcuts';
import { usePointHistory } from '../../hooks/useHistory';
import { useFlowLifecycle, type FlowSegment, type FlowStatus } from '../../features/flow/stateMachine';
import {
  SOURCE_TEMPLATE_SCHEMA_VERSION,
  isTemplateStale,
  loadSourceTemplates,
  saveSourceTemplates,
  upgradeTemplates,
} from '../../features/datalink/sourceTemplateStorage';
import { buildBatchNamePreview } from '../../features/datalink/batchNaming';
import { findNearestValidContiguousSpan } from '../../features/datalink/allocationStrategy';
import {
  canExecuteCommit,
  executeCommitLifecycle,
  retryFailedLifecycle,
  rollbackCommitLifecycle,
  summarizeCommitImpact,
  type CommitQueueViewStatus,
  type QueueBaseStatus,
} from '../../features/datalink/commitLifecycle';
import { buildCommitAuditPayload, type CommitAuditPayload } from '../../features/datalink/commitAudit';
import {
  buildGlobalTagEditDraft,
  getAffectedMappingCountForTag,
  hasGlobalTagEditChanges,
  toTagUpdateRequest,
} from '../../features/datalink/tagEditImpact';
import { getSpanByDataType } from '../../features/datalink/typedOccupancy';
import { addressParser } from '../../utils/addressParser';
import { useSearchParams } from 'react-router-dom';
import { Search, Bell, Settings, Box, Cpu, Sparkles, Keyboard, Upload, Download, Undo2, Redo2, Save, FolderOpen, WandSparkles, Filter } from 'lucide-react';
import { modbusShareAPI } from '../../services/datalink';
import type { ModbusShareStatus } from '../../types/datalink';

const FLOW_SEGMENTS: FlowSegment[] = ['source', 'grid', 'tag', 'sink'];
const COMMIT_CHUNK_SIZE = 8;

const STATUS_STYLE: Record<FlowStatus, string> = {
  draft: 'bg-slate-700/70 text-slate-200 border-slate-600',
  validated: 'bg-blue-500/20 text-blue-300 border-blue-500/40',
  active: 'bg-emerald-500/20 text-emerald-300 border-emerald-500/40',
  error: 'bg-red-500/20 text-red-300 border-red-500/40',
};

const RESET_SEGMENT_DIAGNOSTIC = {
  latestValue: '-',
  quality: 'unknown' as const,
  timestamp: '-',
  error: '',
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
  const [searchParams, setSearchParams] = useSearchParams();
  const [selectedDeviceId, setSelectedDeviceId] = useState<string | null>(null);
  const [selectedAddresses, setSelectedAddresses] = useState<string[]>([]);
  const [selectedPoint, setSelectedPoint] = useState<Point | null>(null);
  const [isTreeCollapsed, setIsTreeCollapsed] = useState(false);
  const [panelType, setPanelType] = useState<'batch' | 'detail' | 'shortcuts' | null>(null);
  const [importDialogOpen, setImportDialogOpen] = useState(false);
  const [exportDialogOpen, setExportDialogOpen] = useState(false);
  const [planDataType, setPlanDataType] = useState<DataType>('int16');
  const [planCount, setPlanCount] = useState(5);
  const [batchNamePrefix, setBatchNamePrefix] = useState('SRC');
  const [planStartAddress, setPlanStartAddress] = useState('40001');
  const [allocationMessage, setAllocationMessage] = useState('');
  const [templateName, setTemplateName] = useState('');
  const [showConflictsOnly, setShowConflictsOnly] = useState(false);
  const [guideStage, setGuideStage] = useState<'idle' | 'grid'>('idle');
  const [sourceTemplates, setSourceTemplates] = useState<SourceTemplate[]>(() => loadSourceTemplates());
  const [modbusStatus, setModbusStatus] = useState<ModbusShareStatus | null>(null);
  const [modbusRegister, setModbusRegister] = useState('0');
  const [modbusActionMessage, setModbusActionMessage] = useState('');
  const [tagLinkMode, setTagLinkMode] = useState<'existing' | 'create'>('existing');
  const [selectedTagIdForLink, setSelectedTagIdForLink] = useState('');
  const [newTagKey, setNewTagKey] = useState('');
  const [newTagDisplayName, setNewTagDisplayName] = useState('');
  const [tagLinkActionMessage, setTagLinkActionMessage] = useState('');
  const [commitActionMessage, setCommitActionMessage] = useState('');
  const [commitQueueRunStatus, setCommitQueueRunStatus] = useState<Record<string, 'success' | 'failed'>>({});
  const [lastCommitSnapshot, setLastCommitSnapshot] = useState<Record<string, 'success' | 'failed'> | null>(null);
  const [isCommitRunning, setIsCommitRunning] = useState(false);
  const [commitChunkResults, setCommitChunkResults] = useState<
    Array<{ chunk: number; totalChunks: number; success: number; failed: number; status: 'success' | 'failed' }>
  >([]);
  const [commitAuditPayload, setCommitAuditPayload] = useState<CommitAuditPayload | null>(null);
  const [failedChunkRetryQueue, setFailedChunkRetryQueue] = useState<number[]>([]);
  const [tagEditDisplayName, setTagEditDisplayName] = useState('');
  const [tagEditUnit, setTagEditUnit] = useState('');
  const [tagEditDescription, setTagEditDescription] = useState('');
  const [tagEditMessage, setTagEditMessage] = useState('');
  const [pendingTagEdit, setPendingTagEdit] = useState<{
    display_name: string;
    unit: string;
    description: string;
  } | null>(null);
  const legacyRoute = searchParams.get('legacy');
  const searchInputRef = useRef<HTMLInputElement | null>(null);
  const gridSectionRef = useRef<HTMLElement | null>(null);

  const { data: devices = [] } = useDevicesQuery();
  const { data: pollingGroups = [] } = usePollingGroupsQuery();
  const { data: allPoints = [] } = usePointsQuery({ device_id: selectedDeviceId || undefined });
  const { data: mappings = [] } = useMappingsQuery();
  const { data: tags = [] } = useTagsQuery();
  const createTagMutation = useCreateTagMutation();
  const updateTagMutation = useUpdateTagMutation();
  const createMappingMutation = useCreateMappingMutation();
  const updateMappingMutation = useUpdateMappingMutation();
  const createPointMutation = useCreatePointMutation();
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
    [devices, selectedDeviceId]
  );

  const legacyRouteLabel = useMemo(() => {
    if (legacyRoute === 'points') return t('nav.points');
    if (legacyRoute === 'mappings') return t('nav.mappings');
    if (legacyRoute === 'wizard') return t('nav.mappingWizard');
    return '';
  }, [legacyRoute, t]);

  const cellSpan = getSpanByDataType(planDataType);
  const totalPlannedCells = planCount * cellSpan;

  const plannedAllocations = useMemo<PlannedAllocation[]>(() => {
    if (!selectedDevice || !planStartAddress || planCount <= 0) return [];
    const expanded = addressParser.expand(planStartAddress, totalPlannedCells, selectedDevice.protocol);
    if (expanded.length < totalPlannedCells) return [];

    return Array.from({ length: planCount }).map((_, index) => {
      const start = index * cellSpan;
      const addresses = expanded.slice(start, start + cellSpan);
      return {
        id: `plan-${index}`,
        dataType: planDataType,
        addresses,
        label: `S${index + 1}`,
      };
    });
  }, [cellSpan, planCount, planDataType, planStartAddress, selectedDevice, totalPlannedCells]);

  const planAddresses = useMemo(
    () => plannedAllocations.flatMap((allocation) => allocation.addresses),
    [plannedAllocations]
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
    () => buildBatchNamePreview(batchNamePrefix, planCount, allPoints.map((point) => point.name)),
    [allPoints, batchNamePrefix, planCount]
  );
  const nameConflictCount = useMemo(
    () => namePreview.filter((item) => item.conflict).length,
    [namePreview]
  );

  const staleTemplateCount = useMemo(
    () => sourceTemplates.filter((template) => isTemplateStale(template)).length,
    [sourceTemplates]
  );

  useEffect(() => {
    setSelectedAddresses([]);
    setSelectedPoint(null);
    setPanelType((prev) => (prev === 'detail' ? null : prev));
  }, [selectedDeviceId]);

  useEffect(() => {
    saveSourceTemplates(sourceTemplates);
  }, [sourceTemplates]);

  useEffect(() => {
    if (!selectedDevice) return;
    setPlanStartAddress(selectedDevice.protocol.startsWith('modbus') ? '40001' : 'D0');
  }, [selectedDevice]);

  const loadModbusStatus = useCallback(async () => {
    try {
      const status = await modbusShareAPI.status();
      setModbusStatus(status);
    } catch (error) {
      const message = error instanceof Error ? error.message : '讀取 Modbus 分享狀態失敗';
      setModbusActionMessage(message);
    }
  }, []);

  useEffect(() => {
    loadModbusStatus();
  }, [loadModbusStatus]);

  const handleBatchCreate = useCallback(() => {
    if (selectedDeviceId) setPanelType('batch');
  }, [selectedDeviceId]);

  const handleClosePanel = useCallback(() => {
    setPanelType(null);
    setSelectedPoint(null);
  }, []);

  const handleToggleSidebar = useCallback(() => {
    setIsTreeCollapsed((prev) => !prev);
  }, []);

  const handleImportPoints = useCallback(
    async (points: CreatePointRequest[]) => {
      const createdIds: string[] = [];
      for (const point of points) {
        const result = await createPointMutation.mutateAsync(point);
        if (result?.id) createdIds.push(result.id);
      }
      history.push({
        type: 'import',
        description: t('smartDashboard.importedPoints', { count: points.length }),
        data: { pointIds: createdIds },
      });
    },
    [createPointMutation, history, t]
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
    if (!normalized || !selectedDevice) return;
    const now = new Date().toISOString();

    const next: SourceTemplate = {
      id: normalized.toLowerCase().replace(/\s+/g, '-'),
      name: normalized,
      dataType: planDataType,
      count: planCount,
      startAddress: planStartAddress,
      updatedAt: now,
      lastUsedAt: now,
      version: SOURCE_TEMPLATE_SCHEMA_VERSION,
    };

    setSourceTemplates((prev) => {
      const existingIndex = prev.findIndex((item) => item.id === next.id);
      if (existingIndex === -1) return [next, ...prev].slice(0, 20);
      const copied = [...prev];
      copied[existingIndex] = next;
      return copied;
    });
    setTemplateName('');
  }, [planCount, planDataType, planStartAddress, selectedDevice, templateName]);

  const handleLoadTemplate = useCallback((template: SourceTemplate) => {
    setPlanDataType(template.dataType);
    setPlanCount(template.count);
    setPlanStartAddress(template.startAddress);
    const now = new Date().toISOString();
    setSourceTemplates((prev) =>
      prev.map((item) =>
        item.id === template.id
          ? {
              ...item,
              lastUsedAt: now,
              version: Math.max(item.version, SOURCE_TEMPLATE_SCHEMA_VERSION),
            }
          : item
      )
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

    if (result.reason === 'invalid_start') {
      setAllocationMessage('起始位址格式無效，請先修正後再試一次。');
      return;
    }

    setAllocationMessage('在目前範圍內找不到可用連續區段，請調整起始位址或降低來源數量。');
  }, [blockedAddresses, planStartAddress, selectedDevice, totalPlannedCells]);

  const handleApplyPlan = useCallback(() => {
    if (!planAddresses.length) return;
    setSelectedAddresses(planAddresses);
    setGuideStage('grid');
    gridSectionRef.current?.scrollIntoView({ behavior: 'smooth', block: 'start' });
    window.setTimeout(() => setGuideStage('idle'), 1200);
  }, [planAddresses]);

  const handleCellClick = (_addr: string, point?: Point) => {
    if (point) {
      setSelectedPoint(point);
      setPanelType('detail');
      return;
    }
    setSelectedPoint(null);
  };

  const getGridCenterAddress = (protocol: ProtocolType) =>
    protocol.startsWith('modbus') ? '40001' : 'D0';

  const primarySelectedAddress = selectedAddresses[0] || '';
  const selectedSourceAddress = primarySelectedAddress || selectedPoint?.address || '';
  const baseCommitQueueItems = useMemo(() => {
    const usedSet = new Set(allPoints.map((point) => point.address));
    const linkedSet = new Set(linkedAddresses);
    return plannedAllocations.map((allocation, index) => {
      const conflictCount = allocation.addresses.filter((address) => usedSet.has(address)).length;
      const linkedCount = allocation.addresses.filter((address) => linkedSet.has(address)).length;
      const status: QueueBaseStatus = conflictCount > 0 ? 'conflict' : linkedCount > 0 ? 'linked' : 'pending';
      return {
        id: allocation.id,
        label: allocation.label,
        type: allocation.dataType,
        addresses: allocation.addresses,
        status,
        order: index + 1,
      };
    });
  }, [allPoints, linkedAddresses, plannedAllocations]);
  const commitQueueItems = useMemo(() => {
    return baseCommitQueueItems.map((item) => {
      const runStatus = commitQueueRunStatus[item.id];
      const viewStatus: CommitQueueViewStatus =
        runStatus === 'success'
          ? 'committed'
          : runStatus === 'failed'
            ? 'failed'
            : item.status;
      return {
        ...item,
        viewStatus,
      };
    });
  }, [baseCommitQueueItems, commitQueueRunStatus]);
  useEffect(() => {
    setCommitQueueRunStatus((prev) => {
      const next: Record<string, 'success' | 'failed'> = {};
      baseCommitQueueItems.forEach((item) => {
        if (prev[item.id]) next[item.id] = prev[item.id];
      });
      return next;
    });
  }, [baseCommitQueueItems]);
  const commitQueueSummary = useMemo(() => {
    return commitQueueItems.reduce(
      (acc, item) => {
        acc.total += 1;
        if (item.viewStatus === 'conflict') acc.conflict += 1;
        if (item.viewStatus === 'linked') acc.linked += 1;
        if (item.viewStatus === 'pending') acc.pending += 1;
        if (item.viewStatus === 'failed') acc.failed += 1;
        if (item.viewStatus === 'committed') acc.committed += 1;
        return acc;
      },
      { total: 0, pending: 0, linked: 0, conflict: 0, failed: 0, committed: 0 }
    );
  }, [commitQueueItems]);
  const commitImpactSummary = useMemo(
    () => summarizeCommitImpact(commitQueueItems, Boolean(pendingTagEdit)),
    [commitQueueItems, pendingTagEdit]
  );
  const selectedPointFromGrid = useMemo(
    () => allPoints.find((point) => point.address === selectedSourceAddress) || null,
    [allPoints, selectedSourceAddress]
  );
  const activePointForLink = selectedPoint ?? selectedPointFromGrid;
  const selectedMapping = useMemo(
    () => mappings.find((mapping) => activePointForLink && mapping.point_id === activePointForLink.id) || null,
    [activePointForLink, mappings]
  );
  const linkedTag = useMemo(
    () => tags.find((tag) => tag.id === selectedMapping?.tag_id) || null,
    [selectedMapping?.tag_id, tags]
  );
  const linkedTagAffectedMappingsCount = useMemo(() => {
    if (!linkedTag?.id) return 0;
    return getAffectedMappingCountForTag(mappings, linkedTag.id);
  }, [linkedTag?.id, mappings]);
  const parsePipeline = useCallback(() => {
    if (!selectedMapping?.transform_pipeline) return [];
    try {
      const parsed = JSON.parse(selectedMapping.transform_pipeline);
      return Array.isArray(parsed) ? parsed : [];
    } catch {
      return [];
    }
  }, [selectedMapping?.transform_pipeline]);

  const segmentFeedback = useMemo(() => {
    const sourceOk = Boolean(selectedDeviceId && selectedSourceAddress);
    const gridOk = planAddresses.length > 0 && planConflictCount === 0;
    const tagOk = Boolean(selectedMapping && linkedTag);
    const sinkOk = canActivate;

    return [
      {
        id: 'source',
        label: 'Source',
        ok: sourceOk,
        message: sourceOk ? '來源已選定' : '請先選取來源設備與格位',
      },
      {
        id: 'grid',
        label: 'Grid',
        ok: gridOk,
        message: gridOk ? '格位驗證通過' : `仍有衝突 ${planConflictCount} 格或尚未規劃`,
      },
      {
        id: 'tag',
        label: 'Tag',
        ok: tagOk,
        message: tagOk ? `已連結 Tag ${linkedTag?.key}` : '請先完成 Tag 連結',
      },
      {
        id: 'sink',
        label: 'Sink',
        ok: sinkOk,
        message: sinkOk ? '可提交到 DB' : '需先完成 Validate',
      },
    ] as const;
  }, [canActivate, linkedTag?.key, planAddresses.length, planConflictCount, selectedDeviceId, selectedMapping, selectedSourceAddress]);

  useEffect(() => {
    setSelectedTagIdForLink(selectedMapping?.tag_id || '');
  }, [selectedMapping?.tag_id]);

  useEffect(() => {
    setTagEditDisplayName(linkedTag?.display_name || '');
    setTagEditUnit(linkedTag?.unit || '');
    setTagEditDescription(linkedTag?.description || '');
    setTagEditMessage('');
    setPendingTagEdit(null);
  }, [linkedTag?.description, linkedTag?.display_name, linkedTag?.id, linkedTag?.unit]);

  const handleLinkTagToSelectedAddress = useCallback(async () => {
    if (!activePointForLink) {
      setTagLinkActionMessage('請先選取已建立點位的格位，再進行 Tag 連結。');
      return;
    }

    const targetTagId = selectedTagIdForLink.trim();
    if (!targetTagId) {
      setTagLinkActionMessage('請先選擇要連結的既有 Tag。');
      return;
    }

    try {
      if (selectedMapping) {
        await updateMappingMutation.mutateAsync({
          id: selectedMapping.id,
          data: {
            tag_id: targetTagId,
            enabled: true,
            transform_pipeline: parsePipeline(),
          },
        });
        setTagLinkActionMessage(`已更新 ${activePointForLink.address} 的 Tag 連結。`);
        return;
      }

      await createMappingMutation.mutateAsync({
        point_id: activePointForLink.id,
        tag_id: targetTagId,
        enabled: true,
      });
      setTagLinkActionMessage(`已建立 ${activePointForLink.address} 的 Tag 連結。`);
    } catch (error) {
      const message = error instanceof Error ? error.message : '連結失敗';
      setTagLinkActionMessage(message);
    }
  }, [
    activePointForLink,
    createMappingMutation,
    parsePipeline,
    selectedMapping,
    selectedTagIdForLink,
    updateMappingMutation,
  ]);

  const handleCreateTagAndLink = useCallback(async () => {
    if (!activePointForLink) {
      setTagLinkActionMessage('請先選取已建立點位的格位，再建立 Tag。');
      return;
    }

    const normalizedKey = newTagKey.trim();
    if (!normalizedKey) {
      setTagLinkActionMessage('Tag Key 不可為空。');
      return;
    }

    const existing = tags.find((tag) => tag.key.toLowerCase() === normalizedKey.toLowerCase());
    if (existing) {
      setTagLinkActionMessage(`Tag Key ${normalizedKey} 已存在，請改用既有 Tag 模式。`);
      return;
    }

    try {
      const createdTag = await createTagMutation.mutateAsync({
        key: normalizedKey,
        display_name: newTagDisplayName.trim() || normalizedKey,
        data_type: activePointForLink.data_type,
      });

      if (selectedMapping) {
        await updateMappingMutation.mutateAsync({
          id: selectedMapping.id,
          data: {
            tag_id: createdTag.id,
            enabled: true,
            transform_pipeline: parsePipeline(),
          },
        });
      } else {
        await createMappingMutation.mutateAsync({
          point_id: activePointForLink.id,
          tag_id: createdTag.id,
          enabled: true,
        });
      }

      setSelectedTagIdForLink(createdTag.id);
      setTagLinkActionMessage(`已建立 Tag ${createdTag.key} 並完成連結。`);
      setNewTagKey('');
      setNewTagDisplayName('');
      setTagLinkMode('existing');
    } catch (error) {
      const message = error instanceof Error ? error.message : '建立 Tag 失敗';
      setTagLinkActionMessage(message);
    }
  }, [
    activePointForLink,
    createMappingMutation,
    createTagMutation,
    newTagDisplayName,
    newTagKey,
    parsePipeline,
    selectedMapping,
    tags,
    updateMappingMutation,
  ]);

  const handleSaveLinkedTagEdit = useCallback(async () => {
    if (!linkedTag?.id) {
      setTagEditMessage('目前沒有可編輯的已連結 Tag。');
      return;
    }

    const nextEdit = buildGlobalTagEditDraft({
      display_name: tagEditDisplayName,
      unit: tagEditUnit,
      description: tagEditDescription,
    });
    if (!hasGlobalTagEditChanges(linkedTag, nextEdit)) {
      setTagEditMessage('沒有變更，無需儲存。');
      return;
    }

    setPendingTagEdit(nextEdit);
    setTagEditMessage(`請確認差異後再儲存，預估影響 ${linkedTagAffectedMappingsCount} 個映射。`);
  }, [
    linkedTag,
    linkedTagAffectedMappingsCount,
    tagEditDescription,
    tagEditDisplayName,
    tagEditUnit,
  ]);

  const handleConfirmTagEdit = useCallback(async () => {
    if (!linkedTag?.id || !pendingTagEdit) return;

    try {
      await updateTagMutation.mutateAsync({
        id: linkedTag.id,
        data: toTagUpdateRequest(pendingTagEdit),
      });
      setPendingTagEdit(null);
      setTagEditMessage(`已更新全域 Tag，影響 ${linkedTagAffectedMappingsCount} 個映射。`);
    } catch (error) {
      const message = error instanceof Error ? error.message : '更新 Tag 失敗';
      setTagEditMessage(message);
    }
  }, [
    linkedTag?.id,
    linkedTagAffectedMappingsCount,
    pendingTagEdit,
    updateTagMutation,
  ]);

  const handleBindTagToModbus = useCallback(async () => {
    if (!linkedTag?.id) {
      setModbusActionMessage('請先選取已連結 Tag 的點位');
      return;
    }
    const register = Number(modbusRegister);
    if (!Number.isInteger(register) || register < 0 || register > 65535) {
      setModbusActionMessage('Register 必須為 0-65535 的整數');
      return;
    }
    try {
      await modbusShareAPI.upsertMapping(linkedTag.id, register);
      await loadModbusStatus();
      setModbusActionMessage(`已綁定 ${linkedTag.key} -> HR${register}`);
    } catch (error) {
      const message = error instanceof Error ? error.message : '綁定失敗';
      setModbusActionMessage(message);
    }
  }, [linkedTag?.id, linkedTag?.key, loadModbusStatus, modbusRegister]);

  const handlePushCurrentValueToModbus = useCallback(async () => {
    if (!linkedTag?.id) {
      setModbusActionMessage('請先選取已連結 Tag 的點位');
      return;
    }
    if (selectedPoint?.last_value === undefined || selectedPoint?.last_value === null) {
      setModbusActionMessage('目前點位沒有可推送的值');
      return;
    }
    try {
      await modbusShareAPI.writeTagValue(linkedTag.id, selectedPoint.last_value);
      setModbusActionMessage(`已推送當前值到 Tag ${linkedTag.key} 的 Modbus 映射`);
    } catch (error) {
      const message = error instanceof Error ? error.message : '推送失敗';
      setModbusActionMessage(message);
    }
  }, [linkedTag?.id, linkedTag?.key, selectedPoint?.last_value]);

  const handleSyncModbusFromMappings = useCallback(async () => {
    try {
      const result = await modbusShareAPI.sync();
      await loadModbusStatus();
      const errorHint = result.errors.length > 0 ? `, errors=${result.errors.length}` : '';
      setModbusActionMessage(`同步完成: updated=${result.updated}, skipped=${result.skipped}${errorHint}`);
    } catch (error) {
      const message = error instanceof Error ? error.message : '同步失敗';
      setModbusActionMessage(message);
    }
  }, [loadModbusStatus]);

  const dismissLegacyNotice = useCallback(() => {
    if (!legacyRoute) return;
    const next = new URLSearchParams(searchParams);
    next.delete('legacy');
    setSearchParams(next, { replace: true });
  }, [legacyRoute, searchParams, setSearchParams]);

  useEffect(() => {
    setSource(selectedDeviceId || '', selectedSourceAddress, selectedPoint?.id || '');
  }, [selectedDeviceId, selectedPoint?.id, selectedSourceAddress, setSource]);

  useEffect(() => {
    setTag(linkedTag?.id || '');
  }, [linkedTag?.id, setTag]);

  useEffect(() => {
    setSink(selectedMapping ? 'timeseries' : '');
  }, [selectedMapping, setSink]);

  useEffect(() => {
    const sourceQuality = selectedPoint?.last_error ? 'bad' : selectedDevice ? 'good' : 'unknown';
    const sourceValue = selectedPoint?.last_value === undefined ? '-' : String(selectedPoint.last_value);
    const sourceTime = selectedPoint?.last_read_at || '-';
    const sourceError = selectedPoint?.last_error || '';
    const gridQuality = selectedSourceAddress ? 'good' : 'unknown';
    const tagQuality = linkedTag ? (selectedMapping?.enabled ? 'good' : 'warning') : 'unknown';
    const sinkQuality = selectedPoint?.last_error ? 'bad' : selectedMapping?.enabled ? 'good' : 'unknown';

    setDiagnostics({
      source: {
        latestValue: sourceValue,
        quality: sourceQuality,
        timestamp: sourceTime,
        error: sourceError,
      },
      grid: {
        latestValue: selectedSourceAddress || '-',
        quality: gridQuality,
        timestamp: sourceTime,
        error: '',
      },
      tag: {
        latestValue: linkedTag ? `${linkedTag.key}` : '-',
        quality: tagQuality,
        timestamp: linkedTag?.updated_at || '-',
        error: '',
      },
      sink: {
        latestValue: flowState.sinkTarget || '-',
        quality: sinkQuality,
        timestamp: selectedMapping?.updated_at || '-',
        error: selectedPoint?.last_error || '',
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
    if (!canValidate) {
      markError('source', t('smartDashboard.flowErrors.missingSource'));
      setCommitActionMessage('Validate 失敗：來源段尚未完成。');
      return;
    }
    if (!selectedMapping) {
      markError('tag', t('smartDashboard.flowErrors.missingMapping'));
      setCommitActionMessage('Validate 失敗：Tag 段尚未完成。');
      return;
    }
    try {
      const pipeline = parsePipeline();
      const result = await validatePipelineMutation.mutateAsync(pipeline);
      if (result.valid) {
        markValidated();
        setCommitActionMessage('Validate 通過：可執行 Commit。');
        return;
      }
      markError('grid', result.error || t('smartDashboard.flowErrors.validationFailed'));
      setCommitActionMessage(`Validate 失敗：${result.error || 'Grid 驗證未通過'}`);
    } catch (error) {
      const message = error instanceof Error ? error.message : t('smartDashboard.flowErrors.validationFailed');
      markError('grid', message);
      setCommitActionMessage(`Validate 失敗：${message}`);
    }
  }, [canValidate, markError, markValidated, parsePipeline, selectedMapping, t, validatePipelineMutation]);

  const handleCommitFlow = useCallback(() => {
    const commitGate = canExecuteCommit({
      canActivate,
      mappingEnabled: Boolean(selectedMapping?.enabled),
    });
    if (!commitGate.ok) {
      if (commitGate.reason === 'not_validated') {
        markError('sink', t('smartDashboard.flowErrors.notValidated'));
        setCommitActionMessage('Commit 失敗：請先完成 Validate。');
        return;
      }
      markError('sink', t('smartDashboard.flowErrors.mappingDisabled'));
      setCommitActionMessage('Commit 失敗：Mapping 尚未啟用。');
      return;
    }
    setIsCommitRunning(true);
    setLastCommitSnapshot(commitQueueRunStatus);
    setCommitChunkResults([]);
    setFailedChunkRetryQueue([]);

    const chunks = Array.from(
      { length: Math.ceil(baseCommitQueueItems.length / COMMIT_CHUNK_SIZE) },
      (_, index) => baseCommitQueueItems.slice(index * COMMIT_CHUNK_SIZE, (index + 1) * COMMIT_CHUNK_SIZE)
    );
    let rollingStatus = { ...commitQueueRunStatus };
    const nextChunkResults: Array<{
      chunk: number;
      totalChunks: number;
      success: number;
      failed: number;
      status: 'success' | 'failed';
    }> = [];
    let successCount = 0;
    let failedCount = 0;

    chunks.forEach((chunkItems, index) => {
      const chunkResult = executeCommitLifecycle(chunkItems, rollingStatus);
      rollingStatus = chunkResult.nextStatus;
      successCount += chunkResult.successCount;
      failedCount += chunkResult.failedCount;
      nextChunkResults.push({
        chunk: index + 1,
        totalChunks: chunks.length,
        success: chunkResult.successCount,
        failed: chunkResult.failedCount,
        status: chunkResult.failedCount > 0 ? 'failed' : 'success',
      });
    });

    setCommitQueueRunStatus(rollingStatus);
    setCommitChunkResults(nextChunkResults);
    setFailedChunkRetryQueue(
      nextChunkResults
        .filter((chunkResult) => chunkResult.status === 'failed')
        .map((chunkResult) => chunkResult.chunk - 1)
    );
    const nextQueueItems = baseCommitQueueItems.map((item) => {
      const runStatus = rollingStatus[item.id];
      const viewStatus: CommitQueueViewStatus =
        runStatus === 'success'
          ? 'committed'
          : runStatus === 'failed'
            ? 'failed'
            : item.status;
      return {
        ...item,
        viewStatus,
      };
    });
    const nextImpact = summarizeCommitImpact(nextQueueItems, Boolean(pendingTagEdit));
    setCommitAuditPayload(
      buildCommitAuditPayload({
        queueItems: nextQueueItems,
        chunkResults: nextChunkResults,
        impact: nextImpact,
      })
    );
    setIsCommitRunning(false);

    if (failedCount > 0) {
      markError('sink', `Commit 部分失敗：${failedCount} 筆失敗，請執行 Retry 或 Rollback。`);
      setCommitActionMessage(`Commit 部分成功：成功 ${successCount}、失敗 ${failedCount}。`);
      return;
    }

    markActive();
    setCommitActionMessage(`Commit 成功：${successCount} 筆已提交並啟用流程。`);
  }, [baseCommitQueueItems, canActivate, commitQueueRunStatus, markActive, markError, pendingTagEdit, selectedMapping?.enabled, t]);

  const handleRetryFailedCommits = useCallback(() => {
    if (failedChunkRetryQueue.length === 0) {
      setCommitActionMessage('沒有可重試的失敗項目。');
      return;
    }

    const chunks = Array.from(
      { length: Math.ceil(baseCommitQueueItems.length / COMMIT_CHUNK_SIZE) },
      (_, index) => baseCommitQueueItems.slice(index * COMMIT_CHUNK_SIZE, (index + 1) * COMMIT_CHUNK_SIZE)
    );

    let rollingStatus = { ...commitQueueRunStatus };
    let recovered = 0;
    let remainingFailed = 0;
    const nextFailedChunkQueue: number[] = [];

    failedChunkRetryQueue.forEach((chunkIndex) => {
      const chunkItems = chunks[chunkIndex] || [];
      const chunkRetry = retryFailedLifecycle(chunkItems, rollingStatus);
      rollingStatus = chunkRetry.nextStatus;
      recovered += chunkRetry.recovered;
      remainingFailed += chunkRetry.remainingFailed;
      if (chunkRetry.remainingFailed > 0) nextFailedChunkQueue.push(chunkIndex);
    });

    setCommitQueueRunStatus(rollingStatus);
    setFailedChunkRetryQueue(nextFailedChunkQueue);

    if (remainingFailed > 0) {
      setCommitActionMessage(`Retry 完成：恢復 ${recovered} 筆，仍有 ${remainingFailed} 筆衝突。`);
      return;
    }

    setCommitActionMessage(`Retry 成功：已恢復 ${recovered} 筆失敗項目。`);
  }, [baseCommitQueueItems, commitQueueRunStatus, failedChunkRetryQueue]);

  const handleRollbackCommitRun = useCallback(() => {
    const rolledBack = rollbackCommitLifecycle(lastCommitSnapshot);
    if (!rolledBack) {
      setCommitActionMessage('目前沒有可回滾的提交快照。');
      return;
    }
    setCommitQueueRunStatus(rolledBack);
    setFailedChunkRetryQueue([]);
    setLastCommitSnapshot(null);
    setCommitActionMessage('已回滾到上次 Commit 前的佇列狀態。');
  }, [lastCommitSnapshot]);

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

  return (
    <div className="flex flex-col min-h-[calc(100vh-11rem)] bg-gradient-to-br from-[#0B0F19] via-[#111827] to-[#0F172A] text-slate-100 font-sans rounded-2xl overflow-hidden selection:bg-blue-500/30">
      <header className="px-4 py-3 sm:px-6 flex flex-wrap items-center justify-between gap-3 border-b border-white/5">
        <div className="flex items-center gap-3 min-w-0">
          <div className="flex items-center justify-center w-10 h-10 rounded-xl bg-blue-600/20 text-blue-400 ring-1 ring-blue-500/30">
            <Box className="w-6 h-6" />
          </div>
          <div>
            <h1 className="text-lg sm:text-xl font-bold tracking-tight bg-gradient-to-r from-white to-slate-400 bg-clip-text text-transparent">
              GoGateway
            </h1>
            <p className="text-[10px] text-slate-400 font-medium tracking-wider uppercase">
              {t('smartDashboard.productSubtitle')}
            </p>
          </div>
        </div>

        <div className="order-3 w-full sm:order-2 sm:w-auto sm:flex-1 sm:max-w-lg sm:mx-4 relative group">
          <label htmlFor="dashboard-search" className="sr-only">
            {t('smartDashboard.searchLabel')}
          </label>
          <div className="absolute inset-0 bg-blue-500/20 blur-xl opacity-0 group-focus-within:opacity-100 transition-opacity rounded-full" />
          <Search className="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 group-focus-within:text-blue-400 transition-colors" />
          <input
            ref={searchInputRef}
            id="dashboard-search"
            type="text"
            placeholder={t('smartDashboard.searchPlaceholder')}
            className="w-full bg-slate-900/50 backdrop-blur-md border border-slate-700/50 rounded-full pl-11 pr-4 py-2.5 text-sm text-slate-200 placeholder:text-slate-500 focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500/50 transition-all shadow-lg shadow-black/20"
          />
        </div>

        <div className="order-2 sm:order-3 flex items-center gap-3">
          <div className="flex items-center gap-2 p-1 bg-slate-800/50 rounded-full border border-slate-700/50">
            <button
              type="button"
              aria-label={t('smartDashboard.notifications')}
              className="p-2 text-slate-400 hover:text-white hover:bg-slate-700/50 rounded-full transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
            >
              <Bell className="w-5 h-5" />
            </button>
            <button
              type="button"
              aria-label={t('smartDashboard.preferences')}
              className="p-2 text-slate-400 hover:text-white hover:bg-slate-700/50 rounded-full transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
            >
              <Settings className="w-5 h-5" />
            </button>
          </div>
          <button
            type="button"
            aria-label={t('smartDashboard.profileMenu')}
            className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center text-sm font-bold shadow-lg shadow-blue-500/25 ring-2 ring-slate-900 hover:ring-offset-2 hover:ring-offset-slate-900 transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
          >
            Y
          </button>
        </div>
      </header>

      {legacyRoute && legacyRouteLabel && (
        <section className="mx-3 mt-3 sm:mx-4 rounded-2xl border border-amber-300/30 bg-amber-500/10 px-4 py-3 text-amber-100 shadow-lg shadow-amber-900/10 transition-all duration-300">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <div>
              <p className="text-xs uppercase tracking-wider text-amber-200/80">{t('smartDashboard.legacyMigration.title')}</p>
              <p className="text-sm">{t('smartDashboard.legacyMigration.description', { route: legacyRouteLabel })}</p>
            </div>
            <button
              type="button"
              onClick={dismissLegacyNotice}
              className="rounded-lg border border-amber-300/40 bg-amber-500/20 px-3 py-1.5 text-xs font-medium hover:bg-amber-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-400"
            >
              {t('smartDashboard.legacyMigration.dismiss')}
            </button>
          </div>
        </section>
      )}

      <div className="flex-1 p-3 sm:p-4 grid grid-cols-1 xl:grid-cols-[auto_1fr_300px] gap-4 min-h-0">
        <div className={`${isTreeCollapsed ? 'xl:w-20' : 'xl:w-[260px]'} min-h-[280px] xl:min-h-0 transition-all duration-300`}>
          <div className="h-full bg-slate-900/60 backdrop-blur-xl border border-white/5 rounded-3xl overflow-hidden shadow-2xl relative">
            <div className="absolute top-0 left-0 right-0 h-40 bg-gradient-to-b from-blue-500/5 to-transparent pointer-events-none" />
            <DeviceTreeNav
              devices={devices}
              selectedDeviceId={selectedDeviceId}
              onSelectDevice={setSelectedDeviceId}
              isCollapsed={isTreeCollapsed}
              onToggleCollapse={() => setIsTreeCollapsed(!isTreeCollapsed)}
              onReorder={() => {}}
            />
          </div>
        </div>

        <div className="min-w-0 flex flex-col">
          <div className="flex-1 bg-slate-900/40 backdrop-blur-md border border-white/5 rounded-3xl overflow-hidden shadow-2xl flex flex-col relative min-h-[420px]">
            {selectedDevice ? (
              <>
                <div className="min-h-16 px-4 sm:px-8 py-3 flex flex-wrap items-center justify-between gap-3 border-b border-white/5 bg-white/[0.02]">
                  <div className="flex items-center gap-4 min-w-0">
                    <div className="w-10 h-10 rounded-2xl bg-gradient-to-br from-emerald-500/20 to-teal-500/20 flex items-center justify-center border border-emerald-500/30">
                      <Cpu className="w-5 h-5 text-emerald-400" />
                    </div>
                    <div className="min-w-0">
                      <div className="flex items-center gap-2 flex-wrap">
                        <h2 className="text-lg font-bold text-white tracking-wide truncate">{selectedDevice.name}</h2>
                        <span
                          className={`flex items-center gap-1.5 px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider ${
                            selectedDevice.status === 'active'
                              ? 'bg-emerald-500/10 border border-emerald-500/20 text-emerald-400'
                              : 'bg-slate-500/20 border border-slate-500/30 text-slate-400'
                          }`}
                        >
                          <span
                            className={`w-1.5 h-1.5 rounded-full ${
                              selectedDevice.status === 'active' ? 'bg-emerald-400 animate-pulse' : 'bg-slate-400'
                            }`}
                            aria-hidden
                          />
                          {selectedDevice.status === 'active'
                            ? t('smartDashboard.connected')
                            : t('smartDashboard.disconnected')}
                        </span>
                      </div>
                      <div className="flex items-center gap-2 text-xs text-slate-400">
                        <span className="font-mono opacity-70">ID: {selectedDevice.id.slice(0, 8)}</span>
                        <span>•</span>
                        <span>{selectedDevice.protocol}</span>
                      </div>
                    </div>
                  </div>

                  <div className="flex flex-wrap items-center gap-3">
                    <div className="px-4 py-2 rounded-xl bg-slate-800/50 border border-white/5 flex flex-col items-end min-w-[100px]">
                      <span className="text-[10px] text-slate-500 uppercase font-bold tracking-wider">
                        {t('smartDashboard.lastSync')}
                      </span>
                      <span className="text-sm font-mono text-slate-200">
                        {selectedDevice.last_test_at
                          ? new Date(selectedDevice.last_test_at).toLocaleTimeString()
                          : '-'}
                      </span>
                    </div>
                    <div className="px-4 py-2 rounded-xl bg-slate-800/50 border border-white/5 flex flex-col items-end min-w-[110px]">
                      <span className="text-[10px] text-slate-500 uppercase font-bold tracking-wider">
                        {t('smartDashboard.cycleTime')}
                      </span>
                      <span className="text-sm font-mono text-blue-400">100 ms</span>
                    </div>
                  </div>
                </div>

                <div className="flex-1 overflow-auto p-4 sm:p-8 scrollbar-thin scrollbar-thumb-slate-700/50 scrollbar-track-transparent">
                  <section className="mb-6 rounded-2xl border border-white/10 bg-slate-900/60 p-4">
                    <div className="mb-3 flex flex-wrap items-center justify-between gap-2">
                      <h3 className="text-sm font-semibold text-slate-100">Source Planner</h3>
                      <span className="text-xs text-slate-400">
                        {planDataType} x {planCount} = {totalPlannedCells} cells
                      </span>
                    </div>
                    <div className="grid grid-cols-1 xl:grid-cols-[1fr_1fr_auto_auto] gap-3">
                      <label className="text-xs text-slate-300">
                        Data Type
                        <select
                          value={planDataType}
                          onChange={(e) => setPlanDataType(e.target.value as DataType)}
                          className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                        >
                          <option value="int16">int16 (1 cell)</option>
                          <option value="int32">int32 (2 cells)</option>
                          <option value="float32">float32 (2 cells)</option>
                          <option value="float64">float64 (4 cells)</option>
                        </select>
                      </label>
                      <label className="text-xs text-slate-300">
                        Source Count
                        <input
                          type="number"
                          min={1}
                          max={200}
                          value={planCount}
                          onChange={(e) => setPlanCount(Math.max(1, Number(e.target.value) || 1))}
                          className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                      </label>
                      <label className="text-xs text-slate-300">
                        Start Address
                        <input
                          value={planStartAddress}
                          onChange={(e) => setPlanStartAddress(e.target.value.toUpperCase())}
                          className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-blue-500"
                        />
                      </label>
                      <div className="flex items-end gap-2">
                        <button
                          type="button"
                          onClick={handleAutoAllocate}
                          className="inline-flex min-h-11 items-center gap-2 rounded-lg border border-sky-500/40 bg-sky-500/15 px-3 py-2 text-xs font-medium text-sky-100 hover:bg-sky-500/25 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sky-500"
                        >
                          <WandSparkles className="h-4 w-4" />
                          Auto
                        </button>
                        <button
                          type="button"
                          onClick={handleApplyPlan}
                          className="inline-flex min-h-11 items-center gap-2 rounded-lg border border-emerald-500/40 bg-emerald-500/20 px-3 py-2 text-xs font-medium text-emerald-100 hover:bg-emerald-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
                        >
                          套用到 Grid
                        </button>
                      </div>
                    </div>
                    <div className="mt-3 grid grid-cols-1 xl:grid-cols-[1fr_auto] gap-3">
                      <label className="text-xs text-slate-300">
                        儲存模板
                        <div className="mt-1 flex gap-2">
                          <input
                            value={templateName}
                            onChange={(e) => setTemplateName(e.target.value)}
                            placeholder="例如: line-a-float32-10"
                            className="w-full rounded-lg border border-slate-700 bg-slate-800/80 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                          />
                          <button
                            type="button"
                            onClick={handleSaveTemplate}
                            disabled={!templateName.trim()}
                            className="inline-flex min-h-11 items-center gap-2 rounded-lg border border-indigo-500/40 bg-indigo-500/20 px-3 py-2 text-xs font-medium text-indigo-100 hover:bg-indigo-500/30 disabled:cursor-not-allowed disabled:opacity-40 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500"
                          >
                            <Save className="h-4 w-4" />
                            儲存
                          </button>
                        </div>
                      </label>
                      <div className="flex items-end gap-2">
                        <button
                          type="button"
                          onClick={() => setShowConflictsOnly((prev) => !prev)}
                          className={`inline-flex min-h-11 items-center gap-2 rounded-lg border px-3 py-2 text-xs font-medium focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500 ${
                            showConflictsOnly
                              ? 'border-amber-500/40 bg-amber-500/20 text-amber-100'
                              : 'border-slate-700 bg-slate-800/70 text-slate-200'
                          }`}
                        >
                          <Filter className="h-4 w-4" />
                          只看衝突
                        </button>
                      </div>
                    </div>
                    <div className="mt-3 rounded-lg border border-white/10 bg-slate-800/40 p-3">
                      <div className="flex flex-wrap items-center justify-between gap-2">
                        <label className="text-xs text-slate-300">
                          批次命名前綴
                          <input
                            value={batchNamePrefix}
                            onChange={(e) => setBatchNamePrefix(e.target.value.toUpperCase())}
                            className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-blue-500"
                            placeholder="例如: LINEA"
                          />
                        </label>
                        <div className="text-xs text-slate-400">
                          預覽 {namePreview.length} 筆，名稱衝突 {nameConflictCount} 筆
                        </div>
                      </div>
                      <div className="mt-2 grid grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-2">
                        {namePreview.slice(0, 12).map((item) => (
                          <div
                            key={item.sequence}
                            className={`rounded-md border px-2 py-1 text-[11px] font-mono ${
                              item.conflict
                                ? 'border-rose-500/40 bg-rose-500/10 text-rose-200'
                                : 'border-slate-700 bg-slate-900/70 text-slate-200'
                            }`}
                          >
                            {item.name}
                          </div>
                        ))}
                      </div>
                    </div>
                    {staleTemplateCount > 0 && (
                      <div className="mt-3 rounded-lg border border-amber-500/30 bg-amber-500/10 px-3 py-2 text-xs text-amber-100">
                        <div className="flex flex-wrap items-center justify-between gap-2">
                          <p>
                            偵測到 {staleTemplateCount} 個舊版模板，建議升級到 v{SOURCE_TEMPLATE_SCHEMA_VERSION} 以確保流程一致性。
                          </p>
                          <button
                            type="button"
                            onClick={handleUpgradeTemplates}
                            className="rounded-md border border-amber-400/40 bg-amber-500/20 px-2 py-1 text-[11px] font-semibold hover:bg-amber-500/30"
                          >
                            升級模板
                          </button>
                        </div>
                      </div>
                    )}
                    <div className="mt-3 flex flex-wrap items-center gap-2">
                      <span className="text-xs text-slate-400">衝突格數: {planConflictCount}</span>
                      {allocationMessage && (
                        <span className="text-xs text-sky-200">{allocationMessage}</span>
                      )}
                      {sourceTemplates.slice(0, 6).map((template) => (
                        <div
                          key={template.id}
                          className="inline-flex items-center gap-1 rounded-full border border-white/10 bg-slate-800/60 pl-2 pr-1 py-1 text-[11px] text-slate-200"
                        >
                          <button
                            type="button"
                            onClick={() => handleLoadTemplate(template)}
                            className="inline-flex items-center gap-1 cursor-pointer hover:text-white"
                          >
                            <FolderOpen className="h-3.5 w-3.5" />
                            {template.name}
                          </button>
                          <span className="rounded bg-slate-700/60 px-1.5 py-0.5 text-[10px] text-slate-300">
                            v{template.version}
                          </span>
                          <span className="text-[10px] text-slate-400">
                            {new Date(template.lastUsedAt).toLocaleDateString()}
                          </span>
                          <button
                            type="button"
                            onClick={() => handleDeleteTemplate(template.id)}
                            className="rounded-full px-1 text-slate-400 hover:bg-slate-700 hover:text-white"
                            aria-label={`Delete template ${template.name}`}
                          >
                            ×
                          </button>
                        </div>
                      ))}
                    </div>
                  </section>
                  <section className="mb-6 rounded-2xl border border-white/10 bg-slate-900/50 p-4">
                    <div className="mb-3 flex items-center justify-between gap-2">
                      <h3 className="text-sm font-semibold text-slate-100">{t('smartDashboard.flowTitle')}</h3>
                      <span className={`text-xs px-2 py-1 rounded-lg border ${STATUS_STYLE[flowState.status]}`}>
                        {t(`smartDashboard.flowStatus.${flowState.status}`)}
                      </span>
                    </div>
                    {hasError && (
                      <p className="mb-3 rounded-lg border border-red-500/40 bg-red-500/10 px-3 py-2 text-xs text-red-200">
                        {t('smartDashboard.flowErrorHint')}
                      </p>
                    )}
                    <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-3">
                      {FLOW_SEGMENTS.map((segment) => {
                        const diag = flowState.diagnostics[segment];
                        return (
                          <article key={segment} className="rounded-xl border border-white/10 bg-slate-800/40 p-3">
                            <div className="flex items-center justify-between">
                              <p className="text-xs uppercase tracking-wider text-slate-400">
                                {t(`smartDashboard.flowSegments.${segment}.title`)}
                              </p>
                              <span
                                className={`h-2 w-2 rounded-full ${
                                  diag.quality === 'good'
                                    ? 'bg-emerald-400'
                                    : diag.quality === 'warning'
                                      ? 'bg-yellow-400'
                                      : diag.quality === 'bad'
                                        ? 'bg-red-400'
                                        : 'bg-slate-500'
                                }`}
                              />
                            </div>
                            <p className="mt-1 text-xs text-slate-500">{t(`smartDashboard.flowSegments.${segment}.subtitle`)}</p>
                            <p className="mt-2 text-sm text-slate-100 font-mono truncate">{diag.latestValue}</p>
                            <p className="mt-1 text-[11px] text-slate-400 truncate">{diag.timestamp}</p>
                            {diag.error && <p className="mt-1 text-[11px] text-red-300 truncate">{diag.error}</p>}
                          </article>
                        );
                      })}
                    </div>
                  </section>
                  <section
                    ref={gridSectionRef}
                    className={`rounded-2xl border border-white/10 bg-slate-900/50 transition-all duration-300 ${
                      guideStage === 'grid' ? 'ring-2 ring-sky-500/70 ring-offset-2 ring-offset-slate-900 motion-safe:animate-pulse' : ''
                    }`}
                  >
                    <MemoryGrid
                      deviceId={selectedDevice.id}
                      protocol={selectedDevice.protocol}
                      centerAddress={planStartAddress || getGridCenterAddress(selectedDevice.protocol)}
                      range={300}
                      existingPoints={allPoints}
                      linkedAddresses={linkedAddresses}
                      selectedAddresses={selectedAddresses}
                      plannedAllocations={plannedAllocations}
                      showConflictsOnly={showConflictsOnly}
                      onSelect={setSelectedAddresses}
                      onCellClick={handleCellClick}
                    />
                  </section>
                </div>
              </>
            ) : (
              <div className="flex-1 flex flex-col items-center justify-center text-slate-400 p-6 sm:p-8 relative overflow-hidden">
                <div className="absolute inset-0 bg-gradient-to-t from-blue-500/5 to-transparent opacity-50" />
                <div className="relative z-10 flex flex-col items-center text-center">
                  <div className="w-24 h-24 sm:w-32 sm:h-32 rounded-3xl bg-gradient-to-br from-blue-500/20 to-purple-500/20 flex items-center justify-center mb-6 sm:mb-8 border border-white/5 shadow-[0_0_50px_rgba(59,130,246,0.2)]">
                    <Sparkles className="w-10 h-10 sm:w-12 sm:h-12 text-blue-400" />
                  </div>
                  <h3 className="text-xl sm:text-2xl font-bold text-white mb-3">{t('smartDashboard.welcomeTitle')}</h3>
                  <p className="text-slate-300 max-w-md mb-6 leading-relaxed">{t('smartDashboard.welcomeDescription')}</p>
                  <p className="text-blue-300 opacity-80">{t('smartDashboard.welcomeHint')}</p>
                </div>
              </div>
            )}
          </div>
        </div>

        <div className="min-h-[300px] xl:min-h-0">
          <div className="h-full bg-slate-900/60 backdrop-blur-xl border border-white/5 rounded-3xl overflow-hidden shadow-2xl">
            <QuickActions
              device={selectedDevice}
              selectedCount={selectedAddresses.length}
              onBatchCreate={() => setPanelType('batch')}
              onQuickMapping={() => {}}
              onTestConnection={() => {}}
            />
            {selectedDevice && (
              <div className="px-4 py-2 border-t border-white/5">
                <div className="flex gap-2">
                  <button
                    onClick={() => setImportDialogOpen(true)}
                    aria-keyshortcuts="Control+I"
                    className="flex-1 flex items-center justify-center gap-2 px-3 py-2 text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 rounded-lg transition-colors border border-slate-700/50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900"
                  >
                    <Upload className="w-4 h-4" />
                    <span>{t('smartDashboard.import')}</span>
                  </button>
                  <button
                    onClick={() => setExportDialogOpen(true)}
                    disabled={allPoints.length === 0}
                    aria-keyshortcuts="Control+E"
                    className="flex-1 flex items-center justify-center gap-2 px-3 py-2 text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 rounded-lg transition-colors border border-slate-700/50 disabled:opacity-50 disabled:cursor-not-allowed focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900"
                  >
                    <Download className="w-4 h-4" />
                    <span>{t('smartDashboard.export')}</span>
                  </button>
                </div>
              </div>
            )}
            <div className="px-4 py-2 border-t border-white/5">
              <div className="flex gap-2">
                <button
                  onClick={handleUndo}
                  disabled={!history.canUndo}
                  aria-keyshortcuts="Control+Z"
                  className="flex-1 flex items-center justify-center gap-2 px-3 py-2 text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 rounded-lg transition-colors border border-slate-700/50 disabled:opacity-30 disabled:cursor-not-allowed focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900"
                  title={history.getUndoAction()?.description || t('smartDashboard.noUndo')}
                >
                  <Undo2 className="w-4 h-4" />
                  <span>{t('smartDashboard.undo')}</span>
                </button>
                <button
                  onClick={handleRedo}
                  disabled={!history.canRedo}
                  aria-keyshortcuts="Control+Y"
                  className="flex-1 flex items-center justify-center gap-2 px-3 py-2 text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 rounded-lg transition-colors border border-slate-700/50 disabled:opacity-30 disabled:cursor-not-allowed focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900"
                  title={history.getRedoAction()?.description || t('smartDashboard.noRedo')}
                >
                  <Redo2 className="w-4 h-4" />
                  <span>{t('smartDashboard.redo')}</span>
                </button>
              </div>
            </div>
            <div className="p-4 border-t border-white/5">
              <button
                onClick={() => setPanelType('shortcuts')}
                className="w-full flex items-center justify-center gap-2 px-4 py-2 text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 rounded-lg transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900"
              >
                <Keyboard className="w-4 h-4" />
                <span>{t('smartDashboard.shortcuts')}</span>
                <span className="ml-auto text-[10px] font-mono opacity-60">?</span>
              </button>
            </div>
            <div className="p-4 border-t border-white/5 space-y-2">
              <div className="flex items-center justify-between">
                <p className="text-xs font-semibold tracking-wide text-slate-200">Commit Queue</p>
                <span className="text-[10px] text-slate-400">Total {commitQueueSummary.total}</span>
              </div>
              <div className="grid grid-cols-3 gap-2 text-[10px]">
                <div className="rounded border border-slate-700 bg-slate-800/60 px-2 py-1 text-slate-200">
                  Pending {commitQueueSummary.pending}
                </div>
                <div className="rounded border border-indigo-500/30 bg-indigo-500/10 px-2 py-1 text-indigo-100">
                  Linked {commitQueueSummary.linked}
                </div>
                <div className="rounded border border-rose-500/30 bg-rose-500/10 px-2 py-1 text-rose-100">
                  Conflict {commitQueueSummary.conflict}
                </div>
              </div>
              <div className="grid grid-cols-2 gap-2 text-[10px]">
                <div className="rounded border border-emerald-500/30 bg-emerald-500/10 px-2 py-1 text-emerald-100">
                  Committed {commitQueueSummary.committed}
                </div>
                <div className="rounded border border-amber-500/30 bg-amber-500/10 px-2 py-1 text-amber-100">
                  Failed {commitQueueSummary.failed}
                </div>
              </div>
              <div className="rounded border border-cyan-500/30 bg-cyan-500/10 px-2 py-2 text-[10px] text-cyan-100 space-y-1">
                <p className="font-semibold tracking-wide">Commit Impact</p>
                <div className="grid grid-cols-3 gap-2">
                  <div>New Points {commitImpactSummary.newPoints}</div>
                  <div>Global Tag Updates {commitImpactSummary.globalTagUpdates}</div>
                  <div>Conflicts {commitImpactSummary.conflicts}</div>
                </div>
              </div>
              <div className="max-h-40 overflow-y-auto space-y-1.5 pr-1">
                {commitQueueItems.length === 0 ? (
                  <p className="rounded border border-slate-700 bg-slate-900/60 px-2 py-2 text-[11px] text-slate-400">
                    尚無待提交規劃
                  </p>
                ) : (
                  commitQueueItems.map((item) => (
                    <div
                      key={item.id}
                      className={`rounded border px-2 py-1.5 text-[11px] ${
                        item.viewStatus === 'conflict' || item.viewStatus === 'failed'
                          ? 'border-rose-500/30 bg-rose-500/10 text-rose-100'
                          : item.viewStatus === 'linked'
                            ? 'border-indigo-500/30 bg-indigo-500/10 text-indigo-100'
                            : item.viewStatus === 'committed'
                              ? 'border-emerald-500/30 bg-emerald-500/10 text-emerald-100'
                            : 'border-slate-700 bg-slate-900/60 text-slate-200'
                      }`}
                    >
                      <div className="flex items-center justify-between">
                        <span className="font-mono">#{item.order} {item.label}</span>
                        <span className="uppercase text-[10px]">{item.viewStatus}</span>
                      </div>
                      <div className="mt-0.5 text-[10px] opacity-80">
                        {item.type} · {item.addresses[0]}..{item.addresses[item.addresses.length - 1]}
                      </div>
                    </div>
                  ))
                )}
              </div>
            </div>
            <div className="p-4 border-t border-white/5 space-y-2">
              <button
                type="button"
                onClick={handleValidateFlow}
                disabled={!canValidate || validatePipelineMutation.isPending}
                aria-keyshortcuts="Control+Enter"
                className="w-full px-3 py-2 text-xs font-medium rounded-lg border border-blue-500/40 bg-blue-500/20 text-blue-100 hover:bg-blue-500/30 disabled:opacity-50 disabled:cursor-not-allowed focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900"
              >
                {validatePipelineMutation.isPending
                  ? t('smartDashboard.validating')
                  : t('smartDashboard.validateFlow')}
              </button>
              <button
                type="button"
                onClick={handleCommitFlow}
                disabled={!canActivate || isCommitRunning}
                aria-keyshortcuts="Control+Shift+Enter"
                className="w-full px-3 py-2 text-xs font-medium rounded-lg border border-emerald-500/40 bg-emerald-500/20 text-emerald-100 hover:bg-emerald-500/30 disabled:opacity-50 disabled:cursor-not-allowed focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900"
              >
                {isCommitRunning ? 'Commit 執行中...' : 'Commit 到 DB'}
              </button>
              <div className="grid grid-cols-2 gap-2">
                <button
                  type="button"
                  onClick={handleRetryFailedCommits}
                  disabled={failedChunkRetryQueue.length === 0}
                  className="rounded-lg border border-amber-500/40 bg-amber-500/20 px-3 py-2 text-xs font-medium text-amber-100 hover:bg-amber-500/30 disabled:opacity-50 disabled:cursor-not-allowed focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500"
                >
                  Retry 失敗 Chunk
                </button>
                <button
                  type="button"
                  onClick={handleRollbackCommitRun}
                  disabled={!lastCommitSnapshot}
                  className="rounded-lg border border-slate-600 bg-slate-700/50 px-3 py-2 text-xs font-medium text-slate-100 hover:bg-slate-700 disabled:opacity-50 disabled:cursor-not-allowed focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-400"
                >
                  Rollback
                </button>
              </div>
              <div className="rounded-lg border border-white/10 bg-slate-900/60 p-2 space-y-1">
                {segmentFeedback.map((segment) => (
                  <div key={segment.id} className="flex items-center justify-between text-[11px]">
                    <span className="text-slate-300">{segment.label}</span>
                    <span className={segment.ok ? 'text-emerald-300' : 'text-amber-300'}>
                      {segment.message}
                    </span>
                  </div>
                ))}
              </div>
              {commitActionMessage && (
                <p className="rounded border border-slate-700 bg-slate-900/70 px-2 py-1.5 text-[11px] text-slate-300">
                  {commitActionMessage}
                </p>
              )}
              {commitAuditPayload && (
                <a
                  href={commitAuditPayload.traceLink}
                  className="inline-flex text-[11px] text-cyan-300 underline decoration-cyan-400/40 underline-offset-2 hover:text-cyan-200"
                >
                  查看稽核追蹤
                </a>
              )}
              {commitChunkResults.length > 0 && (
                <div className="rounded-lg border border-white/10 bg-slate-900/60 p-2 space-y-1">
                  <p className="text-[11px] font-semibold text-slate-200">Chunk 結果</p>
                  {commitChunkResults.map((chunkResult) => (
                    <div key={chunkResult.chunk} className="flex items-center justify-between text-[10px]">
                      <span className="text-slate-300">
                        Chunk {chunkResult.chunk}/{chunkResult.totalChunks}
                      </span>
                      <span className={chunkResult.status === 'success' ? 'text-emerald-300' : 'text-amber-300'}>
                        success {chunkResult.success} / failed {chunkResult.failed}
                      </span>
                    </div>
                  ))}
                </div>
              )}
              {commitAuditPayload && (
                <div id="commit-audit-trace" className="rounded-lg border border-cyan-500/20 bg-slate-900/70 p-2 space-y-1">
                  <p className="text-[11px] font-semibold text-cyan-200">Audit Payload</p>
                  <p className="text-[10px] text-slate-300">
                    {commitAuditPayload.createdAt} · new {commitAuditPayload.summary.newPoints} · tag{' '}
                    {commitAuditPayload.summary.globalTagUpdates} · conflict {commitAuditPayload.summary.conflicts}
                  </p>
                  <pre className="max-h-28 overflow-auto rounded bg-slate-950/70 p-2 text-[10px] text-slate-300">
                    {JSON.stringify(commitAuditPayload, null, 2)}
                  </pre>
                </div>
              )}
              {hasError && (
                <button
                  type="button"
                  onClick={handleRecoverFlow}
                  aria-keyshortcuts="Alt+R"
                  className="w-full px-3 py-2 text-xs font-medium rounded-lg border border-amber-500/40 bg-amber-500/20 text-amber-100 hover:bg-amber-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900"
                >
                  {t('smartDashboard.recoverFlow')}
                </button>
              )}
            </div>
            <div className="p-4 border-t border-white/5 space-y-3">
              <div className="flex items-center justify-between">
                <p className="text-xs font-semibold tracking-wide text-slate-200">Tag Linkage</p>
                <span className="text-[10px] text-slate-400 font-mono">
                  {selectedSourceAddress || '-'}
                </span>
              </div>
              <div className="rounded-lg border border-white/10 bg-slate-900/60 p-3 text-[11px] text-slate-300 space-y-2">
                <p>選取格位: {selectedSourceAddress || '尚未選取'}</p>
                <p>點位: {activePointForLink?.name || '尚未建立點位'}</p>
                <p>目前 Tag: {linkedTag?.key || '未連結'}</p>
              </div>
              <div className="grid grid-cols-2 gap-2">
                <button
                  type="button"
                  onClick={() => setTagLinkMode('existing')}
                  className={`rounded-lg border px-3 py-2 text-xs font-medium focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 ${
                    tagLinkMode === 'existing'
                      ? 'border-blue-500/40 bg-blue-500/20 text-blue-100'
                      : 'border-slate-700 bg-slate-800/70 text-slate-300'
                  }`}
                >
                  選擇既有 Tag
                </button>
                <button
                  type="button"
                  onClick={() => setTagLinkMode('create')}
                  className={`rounded-lg border px-3 py-2 text-xs font-medium focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 ${
                    tagLinkMode === 'create'
                      ? 'border-indigo-500/40 bg-indigo-500/20 text-indigo-100'
                      : 'border-slate-700 bg-slate-800/70 text-slate-300'
                  }`}
                >
                  新建 Tag
                </button>
              </div>
              {tagLinkMode === 'existing' ? (
                <div className="space-y-2">
                  <label className="block text-[11px] text-slate-300">
                    已有 Tag
                    <select
                      value={selectedTagIdForLink}
                      onChange={(e) => setSelectedTagIdForLink(e.target.value)}
                      className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-blue-500"
                    >
                      <option value="">請選擇 Tag</option>
                      {tags.map((tag) => (
                        <option key={tag.id} value={tag.id}>
                          {tag.key} ({tag.data_type})
                        </option>
                      ))}
                    </select>
                  </label>
                  <button
                    type="button"
                    onClick={handleLinkTagToSelectedAddress}
                    disabled={!activePointForLink || !selectedTagIdForLink || updateMappingMutation.isPending || createMappingMutation.isPending}
                    className="w-full rounded-lg border border-blue-500/40 bg-blue-500/20 px-3 py-2 text-xs font-medium text-blue-100 hover:bg-blue-500/30 disabled:opacity-50 disabled:cursor-not-allowed focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
                  >
                    套用 Tag 連結
                  </button>
                </div>
              ) : (
                <div className="space-y-2">
                  <label className="block text-[11px] text-slate-300">
                    Tag Key
                    <input
                      value={newTagKey}
                      onChange={(e) => setNewTagKey(e.target.value)}
                      placeholder="例如: line_a_temp"
                      className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs font-mono focus:outline-none focus:ring-2 focus:ring-indigo-500"
                    />
                  </label>
                  <label className="block text-[11px] text-slate-300">
                    Display Name
                    <input
                      value={newTagDisplayName}
                      onChange={(e) => setNewTagDisplayName(e.target.value)}
                      placeholder="例如: Line A Temperature"
                      className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-indigo-500"
                    />
                  </label>
                  <button
                    type="button"
                    onClick={handleCreateTagAndLink}
                    disabled={!activePointForLink || !newTagKey.trim() || createTagMutation.isPending || createMappingMutation.isPending || updateMappingMutation.isPending}
                    className="w-full rounded-lg border border-indigo-500/40 bg-indigo-500/20 px-3 py-2 text-xs font-medium text-indigo-100 hover:bg-indigo-500/30 disabled:opacity-50 disabled:cursor-not-allowed focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500"
                  >
                    建立並連結 Tag
                  </button>
                </div>
              )}
              {tagLinkActionMessage && (
                <p className="rounded border border-slate-700 bg-slate-900/70 px-2 py-1.5 text-[11px] text-slate-300">
                  {tagLinkActionMessage}
                </p>
              )}
              <div className="rounded-lg border border-white/10 bg-slate-900/60 p-3 space-y-2">
                <div className="flex items-center justify-between">
                  <p className="text-[11px] font-semibold tracking-wide text-slate-200">全域 Tag 內嵌編輯</p>
                  <span className="text-[10px] text-slate-400">
                    影響映射: {linkedTagAffectedMappingsCount}
                  </span>
                </div>
                <label className="block text-[11px] text-slate-300">
                  Display Name
                  <input
                    value={tagEditDisplayName}
                    onChange={(e) => setTagEditDisplayName(e.target.value)}
                    disabled={!linkedTag}
                    className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-amber-500 disabled:opacity-50"
                  />
                </label>
                <label className="block text-[11px] text-slate-300">
                  Unit
                  <input
                    value={tagEditUnit}
                    onChange={(e) => setTagEditUnit(e.target.value)}
                    disabled={!linkedTag}
                    className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-amber-500 disabled:opacity-50"
                  />
                </label>
                <label className="block text-[11px] text-slate-300">
                  Description
                  <input
                    value={tagEditDescription}
                    onChange={(e) => setTagEditDescription(e.target.value)}
                    disabled={!linkedTag}
                    className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-amber-500 disabled:opacity-50"
                  />
                </label>
                <button
                  type="button"
                  onClick={handleSaveLinkedTagEdit}
                  disabled={!linkedTag || updateTagMutation.isPending}
                  className="w-full rounded-lg border border-amber-500/40 bg-amber-500/20 px-3 py-2 text-xs font-medium text-amber-100 hover:bg-amber-500/30 disabled:opacity-50 disabled:cursor-not-allowed focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500"
                >
                  儲存全域 Tag 變更
                </button>
                {pendingTagEdit && linkedTag && (
                  <div className="rounded-lg border border-amber-500/30 bg-amber-500/10 p-2.5 space-y-2">
                    <p className="text-[11px] text-amber-100">
                      變更預覽（第二次確認）: 將影響 {linkedTagAffectedMappingsCount} 個映射
                    </p>
                    <div className="space-y-1 text-[11px] text-slate-200">
                      <p>
                        display_name: <span className="text-slate-400">{linkedTag.display_name || '-'}</span> →{' '}
                        <span className="text-amber-100">{pendingTagEdit.display_name || '-'}</span>
                      </p>
                      <p>
                        unit: <span className="text-slate-400">{linkedTag.unit || '-'}</span> →{' '}
                        <span className="text-amber-100">{pendingTagEdit.unit || '-'}</span>
                      </p>
                      <p>
                        description: <span className="text-slate-400">{linkedTag.description || '-'}</span> →{' '}
                        <span className="text-amber-100">{pendingTagEdit.description || '-'}</span>
                      </p>
                    </div>
                    <div className="grid grid-cols-2 gap-2">
                      <button
                        type="button"
                        onClick={handleConfirmTagEdit}
                        disabled={updateTagMutation.isPending}
                        className="rounded-md border border-amber-400/40 bg-amber-500/20 px-2 py-1.5 text-[11px] font-medium text-amber-100 hover:bg-amber-500/30 disabled:opacity-50"
                      >
                        確認寫入
                      </button>
                      <button
                        type="button"
                        onClick={() => setPendingTagEdit(null)}
                        className="rounded-md border border-slate-700 bg-slate-800/70 px-2 py-1.5 text-[11px] text-slate-200 hover:bg-slate-700"
                      >
                        取消
                      </button>
                    </div>
                  </div>
                )}
                {tagEditMessage && (
                  <p className="rounded border border-slate-700 bg-slate-900/70 px-2 py-1.5 text-[11px] text-slate-300">
                    {tagEditMessage}
                  </p>
                )}
              </div>
            </div>
            <div className="p-4 border-t border-white/5 space-y-3">
              <div className="flex items-center justify-between">
                <p className="text-xs font-semibold tracking-wide text-slate-200">Local Modbus Share</p>
                <button
                  type="button"
                  onClick={loadModbusStatus}
                  className="rounded border border-slate-700 px-2 py-1 text-[10px] text-slate-300 hover:bg-slate-800/70 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
                >
                  Refresh
                </button>
              </div>
              <div className="rounded-lg border border-white/10 bg-slate-900/60 p-3 text-[11px] text-slate-300">
                <p>狀態: {modbusStatus?.enabled ? 'Running' : 'Stopped'}</p>
                <p>Address: {modbusStatus?.address || '-'} (Port 5020)</p>
                <p>Mappings: {modbusStatus?.mapping_count ?? 0}</p>
              </div>
              <label className="block text-[11px] text-slate-300">
                Register (Holding)
                <input
                  value={modbusRegister}
                  onChange={(e) => setModbusRegister(e.target.value)}
                  className="mt-1 w-full rounded-lg border border-slate-700 bg-slate-800/80 px-2.5 py-2 text-xs font-mono focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </label>
              <div className="grid grid-cols-1 gap-2">
                <button
                  type="button"
                  onClick={handleBindTagToModbus}
                  className="w-full rounded-lg border border-indigo-500/40 bg-indigo-500/20 px-3 py-2 text-xs font-medium text-indigo-100 hover:bg-indigo-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500"
                >
                  綁定目前 Tag 到 Register
                </button>
                <button
                  type="button"
                  onClick={handlePushCurrentValueToModbus}
                  className="w-full rounded-lg border border-emerald-500/40 bg-emerald-500/20 px-3 py-2 text-xs font-medium text-emerald-100 hover:bg-emerald-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
                >
                  推送目前值到 Modbus
                </button>
                <button
                  type="button"
                  onClick={handleSyncModbusFromMappings}
                  className="w-full rounded-lg border border-blue-500/40 bg-blue-500/20 px-3 py-2 text-xs font-medium text-blue-100 hover:bg-blue-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
                >
                  同步全部啟用映射
                </button>
              </div>
              {modbusActionMessage && (
                <p className="rounded border border-slate-700 bg-slate-900/70 px-2 py-1.5 text-[11px] text-slate-300">
                  {modbusActionMessage}
                </p>
              )}
            </div>
          </div>
        </div>
      </div>

      <SlidePanel
        isOpen={panelType !== null}
        title={
          panelType === 'batch'
            ? t('smartDashboard.batchCreate')
            : panelType === 'shortcuts'
              ? t('smartDashboard.shortcuts')
              : t('smartDashboard.pointDetail')
        }
        onClose={() => setPanelType(null)}
      >
        {panelType === 'batch' && selectedDevice && (
          <BatchPointCreator
            deviceId={selectedDevice.id}
            protocol={selectedDevice.protocol}
            preselectedAddresses={selectedAddresses}
            pollingGroups={pollingGroups}
            onCreated={() => {
              setPanelType(null);
              setSelectedAddresses([]);
            }}
            onCancel={() => setPanelType(null)}
          />
        )}

        {panelType === 'detail' && selectedPoint && (
          <PointDetailPanel
            point={selectedPoint}
            onUpdate={() => setPanelType(null)}
            onDelete={() => setPanelType(null)}
            onClose={() => setPanelType(null)}
          />
        )}

        {panelType === 'shortcuts' && (
          <div className="space-y-4 p-4">
            <p className="text-sm text-slate-500 dark:text-slate-400 mb-4">{t('smartDashboard.shortcutsHint')}</p>
            <div className="space-y-3">
              {shortcuts.map((shortcut, index) => (
                <div key={index} className="flex items-center justify-between py-2 px-3 bg-slate-100 dark:bg-slate-800/50 rounded-lg">
                  <span className="text-sm text-slate-700 dark:text-slate-300">{shortcut.description}</span>
                  <kbd className="px-2 py-1 text-xs font-mono bg-slate-200 dark:bg-slate-700 rounded border border-slate-300 dark:border-slate-600 text-slate-700 dark:text-slate-300">
                    {shortcut.ctrl && 'Ctrl+'}
                    {shortcut.alt && 'Alt+'}
                    {shortcut.shift && 'Shift+'}
                    {shortcut.key}
                  </kbd>
                </div>
              ))}
            </div>
          </div>
        )}
      </SlidePanel>

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
