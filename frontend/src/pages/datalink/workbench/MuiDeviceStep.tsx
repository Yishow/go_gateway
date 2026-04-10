import { useQueryClient } from '@tanstack/react-query';
import { useEffect, useMemo, useRef, useState, type FormEvent } from 'react';
import { useTranslation } from 'react-i18next';
import Alert from '@mui/material/Alert';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import Chip from '@mui/material/Chip';
import CircularProgress from '@mui/material/CircularProgress';
import IconButton from '@mui/material/IconButton';
import InputAdornment from '@mui/material/InputAdornment';
import MenuItem from '@mui/material/MenuItem';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';
import AddIcon from '@mui/icons-material/Add';
import RefreshIcon from '@mui/icons-material/Refresh';
import SearchIcon from '@mui/icons-material/Search';
import {
  useCreateDeviceMutation,
  useDevicesQuery,
  useTestDraftConnectionMutation,
  useUpdateDeviceMutation,
} from '../../../hooks/datalink/useDevices';
import { deviceKeys } from '../../../hooks/datalink/keys';
import { workbenchExperimentTokens as tokens } from '../../../styles/workbench-experiment-tokens';
import type { ConnectionTestResult, ProtocolType } from '../../../types/datalink';
import { useWorkbench } from './WorkbenchProvider';
import { shouldClearSelectedDevice } from './deviceSelectionGuard';
import {
  buildDeviceDraftFromDevice,
  createEmptyDeviceDraft,
  getWorkbenchProtocolLabelKey,
  sanitizeDeviceConnectionConfig,
  WORKBENCH_PROTOCOLS,
  type DeviceDraft,
} from './workbenchDeviceFormModel';
import { MuiDeviceList } from './MuiDeviceList';
import { MuiDeviceEditor } from './MuiDeviceEditor';
import { MuiDeviceDetail } from './MuiDeviceDetail';
import { MuiDeviceDiagnosticsPanel } from './MuiDeviceDiagnosticsPanel';
import { SENTRY_SX } from './sentrySurfaceStyles';
import { MODBUS_DEFAULT_DATA_FORMAT } from './workbenchDeviceDataFormat';

type FieldErrorMap = Record<string, string>;

function validateFields(draft: DeviceDraft, t: (k: string) => string, full: boolean): FieldErrorMap {
  const e: FieldErrorMap = {};
  const req = (f: string, v: unknown, m: string) => { if (typeof v !== 'string' || !v.trim()) e[f] = m; };
  const reqN = (f: string, v: unknown, m: string) => { if (typeof v !== 'number' || Number.isNaN(v)) e[f] = m; };
  if (full && !draft.name.trim()) e.name = t('workbench.device.validation.nameRequired');
  switch (draft.protocol) {
    case 'modbus_tcp': case 'modbus_udp':
      req('host', draft.connectionConfig.host, t('workbench.device.validation.hostRequired'));
      reqN('port', draft.connectionConfig.port, t('workbench.device.validation.portRequired'));
      break;
    case 'modbus_rtu':
      req('serial_port', draft.connectionConfig.serial_port, t('workbench.device.validation.serialPortRequired'));
      break;
    case 'fatek_fbs': {
      const m = typeof draft.connectionConfig.mode === 'string' ? draft.connectionConfig.mode : 'tcp';
      if (m === 'serial') req('serial_port', draft.connectionConfig.serial_port, t('workbench.device.validation.serialPortRequired'));
      else { req('host', draft.connectionConfig.host, t('workbench.device.validation.hostRequired')); reqN('port', draft.connectionConfig.port, t('workbench.device.validation.portRequired')); }
      reqN('station_no', draft.connectionConfig.station_no, t('workbench.device.validation.stationRequired'));
      break;
    }
    case 'mc_3e':
      req('host', draft.connectionConfig.host, t('workbench.device.validation.hostRequired'));
      reqN('port', draft.connectionConfig.port, t('workbench.device.validation.portRequired'));
      break;
    case 'mqtt':
      req('broker_url', draft.connectionConfig.broker_url, t('workbench.device.validation.brokerUrlRequired'));
      req('client_id', draft.connectionConfig.client_id, t('workbench.device.validation.clientIdRequired'));
      if (!Array.isArray(draft.connectionConfig.topics) || draft.connectionConfig.topics.length === 0) e.topics = t('workbench.device.validation.topicsRequired');
      break;
  }
  return e;
}

/**
 * Sentry Incident Desk — Device command center.
 *
 * Layout: three-zone (incident feed | diagnostics desk | system intel).
 * Diagnostics is the center narrative; task form (create/edit/clone) replaces
 * the center+right zones when a command is active.
 */
export function MuiDeviceStep() {
  const { t } = useTranslation();
  const qc = useQueryClient();
  const {
    clearInspectorSelection, closeDevicePanel, devicePanelState, openCreateDevicePanel,
    openEditDevicePanel, openCloneDevicePanel,
    selectedDeviceId, setActiveStep, setInspectorSelection, setSelectedDeviceId,
    recentDeviceTests, recordDeviceTest,
  } = useWorkbench();

  const { data: devices = [], isLoading, isSuccess, isFetching } = useDevicesQuery();
  const createMut = useCreateDeviceMutation();
  const updateMut = useUpdateDeviceMutation();
  const testMut = useTestDraftConnectionMutation();

  const [search, setSearch] = useState('');
  const [protocolFilter, setProtocolFilter] = useState<'all' | ProtocolType>('all');
  const [draft, setDraft] = useState<DeviceDraft>(() => createEmptyDeviceDraft());
  const [testResult, setTestResult] = useState<ConnectionTestResult | null>(null);
  const [fieldErrors, setFieldErrors] = useState<FieldErrorMap>({});
  const [notice, setNotice] = useState<{ tone: 'success' | 'error' | 'info'; message: string } | null>(null);
  const pendingRef = useRef<string | null>(null);
  const initPanelRef = useRef<string | null>(null);

  const filtered = useMemo(() => {
    const q = search.trim().toLowerCase();
    return devices.filter((d) => {
      const matchesSearch = q === '' || d.name.toLowerCase().includes(q) || d.description.toLowerCase().includes(q);
      const matchesProtocol = protocolFilter === 'all' || d.protocol === protocolFilter;
      return matchesSearch && matchesProtocol;
    });
  }, [devices, protocolFilter, search]);

  const selectedDevice = useMemo(() => devices.find((d) => d.id === selectedDeviceId) ?? null, [devices, selectedDeviceId]);

  useEffect(() => {
    if (!selectedDeviceId) { pendingRef.current = null; return; }
    if (shouldClearSelectedDevice({ selectedDeviceId, devices, isSuccess, isFetching, pendingSelectedDeviceId: pendingRef.current })) {
      setSelectedDeviceId(null); clearInspectorSelection();
    } else if (pendingRef.current === selectedDeviceId) { pendingRef.current = null; }
  }, [clearInspectorSelection, devices, isFetching, isSuccess, selectedDeviceId, setSelectedDeviceId]);

  useEffect(() => {
    if (!devicePanelState) { initPanelRef.current = null; return; }
    const key = devicePanelState.mode === 'create' ? 'create' : devicePanelState.mode === 'edit' ? `edit:${devicePanelState.deviceId}` : `clone:${devicePanelState.sourceDeviceId}`;
    if (initPanelRef.current === key) return;
    if (devicePanelState.mode === 'create') { setDraft(createEmptyDeviceDraft()); setFieldErrors({}); initPanelRef.current = key; return; }
    const src = devices.find((d) => d.id === (devicePanelState.mode === 'edit' ? devicePanelState.deviceId : devicePanelState.sourceDeviceId));
    if (!src) return;
    const next = buildDeviceDraftFromDevice(src);
    setDraft(devicePanelState.mode === 'clone' ? { ...next, name: '' } : next);
    setFieldErrors({}); initPanelRef.current = key;
  }, [devicePanelState, devices]);

  useEffect(() => { setTestResult(null); }, [draft]);

  const closePanel = () => { closeDevicePanel(); setFieldErrors({}); setTestResult(null); };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const errs = validateFields(draft, t, true);
    if (devicePanelState?.mode === 'clone') {
      const src = devices.find((d) => d.id === devicePanelState.sourceDeviceId);
      if (src && draft.name.trim() === src.name.trim()) errs.name = t('workbench.device.validation.cloneNameDistinct');
    }
    setFieldErrors(errs);
    if (Object.keys(errs).length > 0) return;
    try {
      if (devicePanelState?.mode === 'create' || devicePanelState?.mode === 'clone') {
        const created = await createMut.mutateAsync({ name: draft.name.trim(), description: draft.description.trim() || undefined, protocol: draft.protocol, connection_config: sanitizeDeviceConnectionConfig(draft.connectionConfig) });
        pendingRef.current = created.id; setSelectedDeviceId(created.id);
        setNotice({ tone: 'success', message: t('workbench.device.messages.createSuccess') });
      }
      if (devicePanelState?.mode === 'edit') {
        await updateMut.mutateAsync({ id: devicePanelState.deviceId, data: { name: draft.name.trim(), description: draft.description.trim(), connection_config: sanitizeDeviceConnectionConfig(draft.connectionConfig) } });
        setSelectedDeviceId(devicePanelState.deviceId);
        setNotice({ tone: 'success', message: t('workbench.device.messages.updateSuccess') });
      }
      closePanel();
    } catch (err) {
      setNotice({ tone: 'error', message: err instanceof Error ? err.message : t('workbench.device.messages.saveFailed') });
    }
  };

  const handleTest = async () => {
    if (testMut.isPending) return;
    const errs = validateFields(draft, t, false);
    setFieldErrors(errs);
    if (Object.keys(errs).length > 0) return;
    try {
      const connectionConfig = sanitizeDeviceConnectionConfig(draft.connectionConfig);
      if (
        (draft.protocol === 'modbus_tcp' || draft.protocol === 'modbus_udp' || draft.protocol === 'modbus_rtu') &&
        connectionConfig.data_format === MODBUS_DEFAULT_DATA_FORMAT
      ) {
        delete connectionConfig.data_format;
      }
      const r = await testMut.mutateAsync({ protocol: draft.protocol, connection_config: connectionConfig });
      setTestResult(r);
      if (selectedDeviceId) recordDeviceTest(selectedDeviceId, { testedAt: new Date().toISOString(), success: r.success, message: r.error || (r.success ? 'OK' : 'Failed'), latencyMs: r.latency_ms, phaseDetails: { connect: r.connect, probe: r.probe }, canActivate: r.can_activate });
      setNotice({ tone: r.success ? 'success' : 'error', message: r.success ? t('workbench.device.messages.draftTestSuccess') : r.error || t('workbench.device.messages.draftTestFailed') });
    } catch (err) {
      setTestResult(null);
      setNotice({ tone: 'error', message: err instanceof Error ? err.message : t('workbench.device.messages.draftTestFailed') });
    }
  };

  const modeLabel = devicePanelState ? devicePanelState.mode.toUpperCase() : null;
  const hasDevices = !isLoading && devices.length > 0;
  const isTaskMode = Boolean(devicePanelState);

  return (
    <Box component="section" sx={{ display: 'flex', flexDirection: 'column', height: '100%', minHeight: 0 }}>
      {/* ── Command Strip ── */}
      <Box data-testid="device-primary-toolbar" sx={{ ...SENTRY_SX.commandStrip, display: 'flex', alignItems: 'center', gap: 1.5, px: 2, py: 1, mb: 1.5 }}>
        {modeLabel ? <Chip label={modeLabel} size="small" color="warning" sx={SENTRY_SX.modeChip} /> : null}
        <Typography sx={SENTRY_SX.sectionLabel}>{t('workbench.device.eyebrow')}</Typography>
        <Box sx={{ flex: 1 }} />
        {hasDevices ? (
          <TextField size="small" placeholder={t('workbench.device.search.placeholder')} value={search} onChange={(e) => setSearch(e.target.value)}
            slotProps={{ input: { startAdornment: <InputAdornment position="start"><SearchIcon sx={{ fontSize: 16, color: tokens.text.muted }} /></InputAdornment> }, htmlInput: { 'aria-label': t('workbench.device.search.label') } }}
            sx={{ maxWidth: 200, '& .MuiOutlinedInput-root': { height: 30, fontSize: '12px' } }} />
        ) : null}
        {hasDevices ? (
          <TextField
            size="small"
            select
            label={t('workbench.device.filters.protocol')}
            value={protocolFilter}
            onChange={(e) => setProtocolFilter(e.target.value as 'all' | ProtocolType)}
            sx={{ minWidth: 168, '& .MuiOutlinedInput-root': { height: 36, fontSize: '12px' } }}
          >
            <MenuItem value="all">{t('workbench.device.filters.allProtocols')}</MenuItem>
            {WORKBENCH_PROTOCOLS.map((protocol) => (
              <MenuItem key={protocol} value={protocol}>
                {t(getWorkbenchProtocolLabelKey(protocol))}
              </MenuItem>
            ))}
          </TextField>
        ) : null}
        <IconButton
          size="small"
          aria-label={t('workbench.device.actions.refresh')}
          onClick={() => { void qc.invalidateQueries({ queryKey: deviceKeys.lists() }); setNotice({ tone: 'info', message: t('workbench.device.messages.refreshed') }); }}
          sx={SENTRY_SX.insetIconBtn}
        >
          <RefreshIcon sx={{ fontSize: 16 }} />
        </IconButton>
        {devices.length > 0 ? <Button size="small" startIcon={<AddIcon />} onClick={openCreateDevicePanel} sx={SENTRY_SX.insetBtn}>{t('workbench.device.actions.create')}</Button> : null}
      </Box>

      {notice ? <Alert severity={notice.tone === 'success' ? 'success' : notice.tone === 'error' ? 'error' : 'info'} onClose={() => setNotice(null)} sx={{ mb: 1 }}>{notice.message}</Alert> : null}
      {notice ? (
        <Box
          role="status"
          aria-live="polite"
          sx={{
            position: 'absolute',
            width: 1,
            height: 1,
            padding: 0,
            margin: -1,
            overflow: 'hidden',
            clip: 'rect(0, 0, 0, 0)',
            whiteSpace: 'nowrap',
            border: 0,
          }}
        >
          {notice.message}
        </Box>
      ) : null}

      {isLoading ? (
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 2, py: 8 }}>
          <CircularProgress size={24} /><Typography variant="body2" sx={{ color: tokens.text.muted }}>{t('workbench.device.loading')}</Typography>
        </Box>
      ) : null}

      {!isLoading && devices.length === 0 && !isTaskMode ? (
        <Box sx={{ ...SENTRY_SX.frostedPanel, p: 6, textAlign: 'center', mx: 'auto', maxWidth: 520 }}>
          <Typography sx={SENTRY_SX.sectionLabel}>{t('workbench.device.empty.eyebrow')}</Typography>
          <Typography variant="h6" sx={{ mt: 1.5, fontWeight: 700 }}>{t('workbench.device.empty.title')}</Typography>
          <Typography variant="body2" sx={{ color: tokens.text.secondary, mt: 1 }}>{t('workbench.device.empty.description')}</Typography>
          <Button startIcon={<AddIcon />} onClick={openCreateDevicePanel} sx={{ ...SENTRY_SX.insetBtn, mt: 3 }}>{t('workbench.device.actions.create')}</Button>
        </Box>
      ) : null}

      {/* ── Three-Zone Command Center ── */}
      {hasDevices || (isTaskMode && !isLoading) ? (
        <Box sx={{ flex: 1, minHeight: 0, display: 'grid', gap: 1.5, gridTemplateColumns: isTaskMode ? { xs: '1fr', lg: '200px 1fr' } : { xs: '1fr', lg: '200px 1fr 260px' }, overflow: 'hidden' }}>
          {/* Zone 1: Incident Feed */}
          <Box sx={{ overflow: 'auto', opacity: isTaskMode ? 0.6 : 1, transition: 'opacity 0.2s' }}>
            <MuiDeviceList devices={filtered} selectedDeviceId={selectedDeviceId} onSelect={(d) => { setSelectedDeviceId(d.id); setInspectorSelection({ kind: 'device', deviceId: d.id }); setNotice(null); }} />
          </Box>

          {/* Zone 2: Center workspace */}
          {isTaskMode && devicePanelState ? (
            <Box sx={{ overflow: 'auto' }}>
              <MuiDeviceEditor panelState={devicePanelState} draft={draft} setDraft={setDraft} fieldErrors={fieldErrors} setFieldErrors={setFieldErrors} draftTestResult={testResult} isSaving={createMut.isPending || updateMut.isPending} isTesting={testMut.isPending} onClose={closePanel} onSubmit={handleSubmit} onTestConnection={() => void handleTest()} />
            </Box>
          ) : (
            <Box sx={{ overflow: 'auto', display: 'flex', flexDirection: 'column', gap: 1.5 }}>
              <MuiDeviceDiagnosticsPanel device={selectedDevice} recentTests={selectedDeviceId ? recentDeviceTests[selectedDeviceId] ?? [] : []} onContinue={() => setActiveStep('source')} onEdit={selectedDevice ? () => openEditDevicePanel(selectedDevice.id) : undefined} onClone={selectedDevice ? () => openCloneDevicePanel(selectedDevice.id) : undefined} />
            </Box>
          )}

          {/* Zone 3: System Intel (hidden in task mode) */}
          {!isTaskMode && selectedDevice ? (
            <Box sx={{ overflow: 'auto' }}>
              <MuiDeviceDetail device={selectedDevice} />
            </Box>
          ) : null}
        </Box>
      ) : null}
    </Box>
  );
}
