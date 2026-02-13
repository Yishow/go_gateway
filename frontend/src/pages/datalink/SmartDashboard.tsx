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
import { useMappingsQuery, useValidatePipelineMutation } from '../../hooks/datalink/useMappings';
import { useTagsQuery } from '../../hooks/datalink/useTags';
import { useSmartDashboardShortcuts } from '../../hooks/useKeyboardShortcuts';
import { usePointHistory } from '../../hooks/useHistory';
import { useFlowLifecycle, type FlowSegment, type FlowStatus } from '../../features/flow/stateMachine';
import { addressParser } from '../../utils/addressParser';
import { Search, Bell, Settings, Box, Cpu, Sparkles, Keyboard, Upload, Download, Undo2, Redo2, Save, FolderOpen, WandSparkles, Filter } from 'lucide-react';

const FLOW_SEGMENTS: FlowSegment[] = ['source', 'grid', 'tag', 'sink'];

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

const SOURCE_TEMPLATE_STORAGE_KEY = 'pipeline-studio-source-templates-v1';

const SPAN_BY_TYPE: Record<DataType, number> = {
  bool: 1,
  int16: 1,
  uint16: 1,
  int32: 2,
  uint32: 2,
  float32: 2,
  int64: 4,
  uint64: 4,
  float64: 4,
  string: 1,
};

interface SourceTemplate {
  id: string;
  name: string;
  dataType: DataType;
  count: number;
  startAddress: string;
  updatedAt: string;
}

export default function SmartDashboard() {
  const { t } = useTranslation();
  const [selectedDeviceId, setSelectedDeviceId] = useState<string | null>(null);
  const [selectedAddresses, setSelectedAddresses] = useState<string[]>([]);
  const [selectedPoint, setSelectedPoint] = useState<Point | null>(null);
  const [isTreeCollapsed, setIsTreeCollapsed] = useState(false);
  const [panelType, setPanelType] = useState<'batch' | 'detail' | 'shortcuts' | null>(null);
  const [importDialogOpen, setImportDialogOpen] = useState(false);
  const [exportDialogOpen, setExportDialogOpen] = useState(false);
  const [planDataType, setPlanDataType] = useState<DataType>('int16');
  const [planCount, setPlanCount] = useState(5);
  const [planStartAddress, setPlanStartAddress] = useState('40001');
  const [templateName, setTemplateName] = useState('');
  const [showConflictsOnly, setShowConflictsOnly] = useState(false);
  const [guideStage, setGuideStage] = useState<'idle' | 'grid'>('idle');
  const [sourceTemplates, setSourceTemplates] = useState<SourceTemplate[]>(() => {
    try {
      const raw = localStorage.getItem(SOURCE_TEMPLATE_STORAGE_KEY);
      if (!raw) return [];
      const parsed = JSON.parse(raw);
      return Array.isArray(parsed) ? parsed : [];
    } catch {
      return [];
    }
  });
  const searchInputRef = useRef<HTMLInputElement | null>(null);
  const gridSectionRef = useRef<HTMLElement | null>(null);

  const { data: devices = [] } = useDevicesQuery();
  const { data: pollingGroups = [] } = usePollingGroupsQuery();
  const { data: allPoints = [] } = usePointsQuery({ device_id: selectedDeviceId || undefined });
  const { data: mappings = [] } = useMappingsQuery();
  const { data: tags = [] } = useTagsQuery();
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

  const cellSpan = SPAN_BY_TYPE[planDataType];
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

  useEffect(() => {
    setSelectedAddresses([]);
    setSelectedPoint(null);
    setPanelType((prev) => (prev === 'detail' ? null : prev));
  }, [selectedDeviceId]);

  useEffect(() => {
    localStorage.setItem(SOURCE_TEMPLATE_STORAGE_KEY, JSON.stringify(sourceTemplates));
  }, [sourceTemplates]);

  useEffect(() => {
    if (!selectedDevice) return;
    setPlanStartAddress(selectedDevice.protocol.startsWith('modbus') ? '40001' : 'D0');
  }, [selectedDevice]);

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

    const next: SourceTemplate = {
      id: normalized.toLowerCase().replace(/\s+/g, '-'),
      name: normalized,
      dataType: planDataType,
      count: planCount,
      startAddress: planStartAddress,
      updatedAt: new Date().toISOString(),
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
  }, []);

  const handleDeleteTemplate = useCallback((templateId: string) => {
    setSourceTemplates((prev) => prev.filter((item) => item.id !== templateId));
  }, []);

  const handleAutoAllocate = useCallback(() => {
    if (!selectedDevice || !planStartAddress || totalPlannedCells <= 0) return;
    const used = new Set(allPoints.map((point) => point.address));
    const candidates = addressParser.expand(planStartAddress, 800, selectedDevice.protocol);

    const nextStart = candidates.find((candidateStart) => {
      const span = addressParser.expand(candidateStart, totalPlannedCells, selectedDevice.protocol);
      if (span.length !== totalPlannedCells) return false;
      return span.every((address) => !used.has(address));
    });

    if (nextStart) {
      setPlanStartAddress(nextStart);
    }
  }, [allPoints, planStartAddress, selectedDevice, totalPlannedCells]);

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
  const selectedMapping = useMemo(
    () => mappings.find((mapping) => selectedPoint && mapping.point_id === selectedPoint.id) || null,
    [mappings, selectedPoint]
  );
  const linkedTag = useMemo(
    () => tags.find((tag) => tag.id === selectedMapping?.tag_id) || null,
    [selectedMapping?.tag_id, tags]
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
      return;
    }
    if (!selectedMapping) {
      markError('tag', t('smartDashboard.flowErrors.missingMapping'));
      return;
    }
    try {
      const pipeline = parsePipeline();
      const result = await validatePipelineMutation.mutateAsync(pipeline);
      if (result.valid) {
        markValidated();
        return;
      }
      markError('grid', result.error || t('smartDashboard.flowErrors.validationFailed'));
    } catch (error) {
      const message = error instanceof Error ? error.message : t('smartDashboard.flowErrors.validationFailed');
      markError('grid', message);
    }
  }, [canValidate, markError, markValidated, parsePipeline, selectedMapping, t, validatePipelineMutation]);

  const handleActivateFlow = useCallback(() => {
    if (!canActivate) {
      markError('sink', t('smartDashboard.flowErrors.notValidated'));
      return;
    }
    if (!selectedMapping?.enabled) {
      markError('sink', t('smartDashboard.flowErrors.mappingDisabled'));
      return;
    }
    markActive();
  }, [canActivate, markActive, markError, selectedMapping?.enabled, t]);

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
    onActivateFlow: handleActivateFlow,
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
                        <span className="px-2 py-0.5 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-[10px] font-bold text-emerald-400 uppercase tracking-wider">
                          {t('smartDashboard.active')}
                        </span>
                      </div>
                      <div className="flex items-center gap-2 text-xs text-slate-400">
                        <span className="font-mono opacity-70">ID: {selectedDevice.id.slice(0, 8)}</span>
                        <span>•</span>
                        <span>{selectedDevice.protocol}</span>
                      </div>
                    </div>
                  </div>

                  <div className="flex gap-3">
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
                    <div className="mt-3 flex flex-wrap items-center gap-2">
                      <span className="text-xs text-slate-400">衝突格數: {planConflictCount}</span>
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
                onClick={handleActivateFlow}
                disabled={!canActivate}
                aria-keyshortcuts="Control+Shift+Enter"
                className="w-full px-3 py-2 text-xs font-medium rounded-lg border border-emerald-500/40 bg-emerald-500/20 text-emerald-100 hover:bg-emerald-500/30 disabled:opacity-50 disabled:cursor-not-allowed focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900"
              >
                {t('smartDashboard.activateFlow')}
              </button>
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
