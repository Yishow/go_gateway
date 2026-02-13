import { useState, useMemo, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import type { Point, CreatePointRequest, ProtocolType } from '../../types/datalink';
import { DeviceTreeNav } from '../../components/datalink/DeviceTreeNav';
import { MemoryGrid } from '../../components/datalink/MemoryGrid';
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
import { Search, Bell, Settings, Box, Cpu, Sparkles, Keyboard, Upload, Download, Undo2, Redo2 } from 'lucide-react';

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

export default function SmartDashboard() {
  const { t } = useTranslation();
  const [selectedDeviceId, setSelectedDeviceId] = useState<string | null>(null);
  const [selectedAddresses, setSelectedAddresses] = useState<string[]>([]);
  const [selectedPoint, setSelectedPoint] = useState<Point | null>(null);
  const [isTreeCollapsed, setIsTreeCollapsed] = useState(false);
  const [panelType, setPanelType] = useState<'batch' | 'detail' | 'shortcuts' | null>(null);
  const [importDialogOpen, setImportDialogOpen] = useState(false);
  const [exportDialogOpen, setExportDialogOpen] = useState(false);

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

  const shortcuts = useSmartDashboardShortcuts({
    onBatchCreate: handleBatchCreate,
    onClosePanel: handleClosePanel,
    onToggleSidebar: handleToggleSidebar,
    onUndo: handleUndo,
    onRedo: handleRedo,
    onImport: handleImportShortcut,
    onExport: handleExportShortcut,
  });

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
                  <MemoryGrid
                    deviceId={selectedDevice.id}
                    protocol={selectedDevice.protocol}
                    centerAddress={getGridCenterAddress(selectedDevice.protocol)}
                    range={300}
                    existingPoints={allPoints}
                    selectedAddresses={selectedAddresses}
                    onSelect={setSelectedAddresses}
                    onCellClick={handleCellClick}
                  />
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
                    className="flex-1 flex items-center justify-center gap-2 px-3 py-2 text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 rounded-lg transition-colors border border-slate-700/50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
                  >
                    <Upload className="w-4 h-4" />
                    <span>{t('smartDashboard.import')}</span>
                  </button>
                  <button
                    onClick={() => setExportDialogOpen(true)}
                    disabled={allPoints.length === 0}
                    className="flex-1 flex items-center justify-center gap-2 px-3 py-2 text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 rounded-lg transition-colors border border-slate-700/50 disabled:opacity-50 disabled:cursor-not-allowed focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
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
                  className="flex-1 flex items-center justify-center gap-2 px-3 py-2 text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 rounded-lg transition-colors border border-slate-700/50 disabled:opacity-30 disabled:cursor-not-allowed focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
                  title={history.getUndoAction()?.description || t('smartDashboard.noUndo')}
                >
                  <Undo2 className="w-4 h-4" />
                  <span>{t('smartDashboard.undo')}</span>
                </button>
                <button
                  onClick={handleRedo}
                  disabled={!history.canRedo}
                  className="flex-1 flex items-center justify-center gap-2 px-3 py-2 text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 rounded-lg transition-colors border border-slate-700/50 disabled:opacity-30 disabled:cursor-not-allowed focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
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
                className="w-full flex items-center justify-center gap-2 px-4 py-2 text-xs text-slate-300 hover:text-white hover:bg-slate-800/50 rounded-lg transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
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
                className="w-full px-3 py-2 text-xs font-medium rounded-lg border border-blue-500/40 bg-blue-500/20 text-blue-200 hover:bg-blue-500/30 disabled:opacity-50 disabled:cursor-not-allowed focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
              >
                {validatePipelineMutation.isPending
                  ? t('smartDashboard.validating')
                  : t('smartDashboard.validateFlow')}
              </button>
              <button
                type="button"
                onClick={handleActivateFlow}
                disabled={!canActivate}
                className="w-full px-3 py-2 text-xs font-medium rounded-lg border border-emerald-500/40 bg-emerald-500/20 text-emerald-200 hover:bg-emerald-500/30 disabled:opacity-50 disabled:cursor-not-allowed focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500"
              >
                {t('smartDashboard.activateFlow')}
              </button>
              {hasError && (
                <button
                  type="button"
                  onClick={handleRecoverFlow}
                  className="w-full px-3 py-2 text-xs font-medium rounded-lg border border-amber-500/40 bg-amber-500/20 text-amber-200 hover:bg-amber-500/30 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-500"
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
