import { useQueryClient } from '@tanstack/react-query';
import {
  useEffect,
  useMemo,
  useRef,
  useState,
  type InputHTMLAttributes,
  type FormEvent,
} from 'react';
import { useTranslation } from 'react-i18next';
import {
  useCreateDeviceMutation,
  useDevicesQuery,
  useTestDraftConnectionMutation,
  useUpdateDeviceMutation,
} from '../../../hooks/datalink/useDevices';
import { deviceKeys } from '../../../hooks/datalink/keys';
import type {
  ConnectionTestResult,
  ConnectionTestStageResult,
  Device,
  DeviceStatus,
  ProtocolType,
} from '../../../types/datalink';
import { Spinner } from '../../../components/ui/spinner';
import { useWorkbench } from './WorkbenchProvider';
import { shouldClearSelectedDevice } from './deviceSelectionGuard';
import {
  buildDeviceCapabilitySummary,
  buildDeviceEndpointSummary,
  buildDeviceDraftFromDevice,
  createDefaultDeviceConnectionConfig,
  createEmptyDeviceDraft,
  getWorkbenchDeviceStatusLabelKey,
  getWorkbenchProtocolLabelKey,
  normalizeMc3eDataFormatValue,
  parseDeviceConnectionConfig,
  sanitizeDeviceConnectionConfig,
  WORKBENCH_DEVICE_STATUSES,
  WORKBENCH_PROTOCOLS,
  type DeviceConnectionConfig,
  type DeviceDraft,
} from './workbenchDeviceFormModel';

const inputClassName =
  'w-full rounded-xl border border-slate-700 bg-slate-950/70 px-3 py-2 text-sm text-slate-100 placeholder:text-slate-500 focus:border-cyan-400 focus:outline-none focus:ring-2 focus:ring-cyan-400/30';
const textareaClassName = `${inputClassName} min-h-24 resize-y`;
const labelClassName =
  'text-xs font-semibold uppercase tracking-[0.18em] text-slate-400';
const ghostButtonClassName =
  'inline-flex items-center justify-center rounded-xl border border-slate-700 bg-slate-950/60 px-3 py-2 text-sm font-medium text-slate-200 transition hover:border-slate-500 hover:text-slate-50 disabled:cursor-not-allowed disabled:opacity-50';
const primaryButtonClassName =
  'inline-flex items-center justify-center rounded-xl bg-cyan-500 px-3 py-2 text-sm font-semibold text-slate-950 transition hover:bg-cyan-400 disabled:cursor-not-allowed disabled:opacity-50';

type DeviceNotice = {
  tone: 'success' | 'error' | 'info';
  message: string;
};

type FieldErrorMap = Record<string, string>;

type SelectOption = {
  value: string;
  label: string;
};

function joinClasses(...classNames: Array<string | false | null | undefined>) {
  return classNames.filter(Boolean).join(' ');
}

function TextField({
  id,
  label,
  value,
  onChange,
  placeholder,
  error,
  inputMode,
  disabled,
  type = 'text',
}: {
  id: string;
  label: string;
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  error?: string;
  inputMode?: InputHTMLAttributes<HTMLInputElement>['inputMode'];
  disabled?: boolean;
  type?: string;
}) {
  return (
    <label className="space-y-2" htmlFor={id}>
      <span className={labelClassName}>{label}</span>
      <input
        id={id}
        className={joinClasses(
          inputClassName,
          error ? 'border-rose-500 focus:border-rose-400 focus:ring-rose-400/30' : undefined,
        )}
        disabled={disabled}
        inputMode={inputMode}
        onChange={(event) => onChange(event.target.value)}
        placeholder={placeholder}
        type={type}
        value={value}
      />
      {error ? <p className="text-xs text-rose-300">{error}</p> : null}
    </label>
  );
}

function TextAreaField({
  id,
  label,
  value,
  onChange,
  placeholder,
  error,
}: {
  id: string;
  label: string;
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  error?: string;
}) {
  return (
    <label className="space-y-2" htmlFor={id}>
      <span className={labelClassName}>{label}</span>
      <textarea
        id={id}
        className={joinClasses(
          textareaClassName,
          error ? 'border-rose-500 focus:border-rose-400 focus:ring-rose-400/30' : undefined,
        )}
        onChange={(event) => onChange(event.target.value)}
        placeholder={placeholder}
        value={value}
      />
      {error ? <p className="text-xs text-rose-300">{error}</p> : null}
    </label>
  );
}

function SelectField({
  id,
  label,
  value,
  onChange,
  options,
  error,
  disabled,
}: {
  id: string;
  label: string;
  value: string;
  onChange: (value: string) => void;
  options: SelectOption[];
  error?: string;
  disabled?: boolean;
}) {
  return (
    <label className="space-y-2" htmlFor={id}>
      <span className={labelClassName}>{label}</span>
      <select
        id={id}
        className={joinClasses(
          inputClassName,
          error ? 'border-rose-500 focus:border-rose-400 focus:ring-rose-400/30' : undefined,
        )}
        disabled={disabled}
        onChange={(event) => onChange(event.target.value)}
        value={value}
      >
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
      {error ? <p className="text-xs text-rose-300">{error}</p> : null}
    </label>
  );
}

function CheckboxField({
  id,
  label,
  checked,
  onChange,
}: {
  id: string;
  label: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
}) {
  return (
    <label
      className="flex items-center gap-3 rounded-2xl border border-slate-800 bg-slate-950/40 px-4 py-3 text-sm text-slate-200"
      htmlFor={id}
    >
      <input
        id={id}
        checked={checked}
        className="h-4 w-4 rounded border-slate-600 bg-slate-900 text-cyan-400 focus:ring-cyan-400/30"
        onChange={(event) => onChange(event.target.checked)}
        type="checkbox"
      />
      <span>{label}</span>
    </label>
  );
}

function getStatusClasses(status: DeviceStatus) {
  switch (status) {
    case 'active':
      return 'border-emerald-500/40 bg-emerald-500/10 text-emerald-200';
    case 'disabled':
      return 'border-rose-500/40 bg-rose-500/10 text-rose-200';
    case 'draft':
      return 'border-amber-500/40 bg-amber-500/10 text-amber-200';
  }
}

function getNoticeClasses(tone: DeviceNotice['tone']) {
  switch (tone) {
    case 'success':
      return 'border-emerald-500/40 bg-emerald-500/10 text-emerald-100';
    case 'error':
      return 'border-rose-500/40 bg-rose-500/10 text-rose-100';
    case 'info':
      return 'border-cyan-500/40 bg-cyan-500/10 text-cyan-100';
  }
}

function getDraftStageStatusClasses(
  status: ConnectionTestStageResult['status'] | undefined,
) {
  switch (status) {
    case 'success':
      return 'border-emerald-500/40 bg-emerald-500/10 text-emerald-200';
    case 'failed':
      return 'border-rose-500/40 bg-rose-500/10 text-rose-200';
    case 'skipped':
    default:
      return 'border-slate-700 bg-slate-900/70 text-slate-300';
  }
}

function getDraftStageMessage(stage: ConnectionTestStageResult | undefined) {
  if (!stage) {
    return '—';
  }

  return stage.error || stage.message || '—';
}

function getDeviceHealthLabel(t: (key: string) => string, device: Device) {
  if (device.last_test_success === true) {
    return t('workbench.device.card.testPassed');
  }

  if (device.last_test_success === false) {
    return t('workbench.device.card.testFailed');
  }

  return t('workbench.device.card.notTested');
}

function getDeviceHealthClasses(lastTestSuccess: boolean | null) {
  if (lastTestSuccess === true) {
    return 'border-emerald-500/40 bg-emerald-500/10 text-emerald-200';
  }

  if (lastTestSuccess === false) {
    return 'border-rose-500/40 bg-rose-500/10 text-rose-200';
  }

  return 'border-slate-700 bg-slate-900/70 text-slate-300';
}

function requireText(
  errors: FieldErrorMap,
  field: string,
  value: DeviceConnectionConfig[string],
  message: string,
) {
  if (typeof value !== 'string' || value.trim() === '') {
    errors[field] = message;
  }
}

function requireNumber(
  errors: FieldErrorMap,
  field: string,
  value: DeviceConnectionConfig[string],
  message: string,
) {
  if (typeof value !== 'number' || Number.isNaN(value)) {
    errors[field] = message;
  }
}

function validateDeviceConnectionFields(
  draft: DeviceDraft,
  t: (key: string) => string,
): FieldErrorMap {
  const errors: FieldErrorMap = {};

  switch (draft.protocol) {
    case 'modbus_tcp':
    case 'modbus_udp':
      requireText(
        errors,
        'host',
        draft.connectionConfig.host,
        t('workbench.device.validation.hostRequired'),
      );
      requireNumber(
        errors,
        'port',
        draft.connectionConfig.port,
        t('workbench.device.validation.portRequired'),
      );
      break;
    case 'modbus_rtu':
      requireText(
        errors,
        'serial_port',
        draft.connectionConfig.serial_port,
        t('workbench.device.validation.serialPortRequired'),
      );
      break;
    case 'fatek_fbs': {
      const mode =
        typeof draft.connectionConfig.mode === 'string'
          ? draft.connectionConfig.mode
          : 'tcp';
      if (mode === 'serial') {
        requireText(
          errors,
          'serial_port',
          draft.connectionConfig.serial_port,
          t('workbench.device.validation.serialPortRequired'),
        );
      } else {
        requireText(
          errors,
          'host',
          draft.connectionConfig.host,
          t('workbench.device.validation.hostRequired'),
        );
        requireNumber(
          errors,
          'port',
          draft.connectionConfig.port,
          t('workbench.device.validation.portRequired'),
        );
      }
      requireNumber(
        errors,
        'station_no',
        draft.connectionConfig.station_no,
        t('workbench.device.validation.stationRequired'),
      );
      break;
    }
    case 'mc_3e':
      requireText(
        errors,
        'host',
        draft.connectionConfig.host,
        t('workbench.device.validation.hostRequired'),
      );
      requireNumber(
        errors,
        'port',
        draft.connectionConfig.port,
        t('workbench.device.validation.portRequired'),
      );
      break;
    case 'mqtt':
      requireText(
        errors,
        'broker_url',
        draft.connectionConfig.broker_url,
        t('workbench.device.validation.brokerUrlRequired'),
      );
      requireText(
        errors,
        'client_id',
        draft.connectionConfig.client_id,
        t('workbench.device.validation.clientIdRequired'),
      );
      if (
        !Array.isArray(draft.connectionConfig.topics) ||
        draft.connectionConfig.topics.length === 0
      ) {
        errors.topics = t('workbench.device.validation.topicsRequired');
      }
      break;
  }

  return errors;
}

function validateDeviceDraft(
  draft: DeviceDraft,
  t: (key: string) => string,
): FieldErrorMap {
  const errors = validateDeviceConnectionFields(draft, t);

  if (draft.name.trim() === '') {
    errors.name = t('workbench.device.validation.nameRequired');
  }

  return errors;
}

export function WorkbenchDeviceStep() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const {
    clearInspectorSelection,
    closeDevicePanel,
    devicePanelState,
    openCreateDevicePanel,
    selectedDeviceId,
    setActiveStep,
    setInspectorSelection,
    setSelectedDeviceId,
  } = useWorkbench();
  const {
    data: devices = [],
    isLoading,
    isSuccess: devicesQuerySuccess,
    isFetching: devicesQueryFetching,
  } = useDevicesQuery();
  const createDeviceMutation = useCreateDeviceMutation();
  const testDraftConnectionMutation = useTestDraftConnectionMutation();
  const updateDeviceMutation = useUpdateDeviceMutation();

  const [searchQuery, setSearchQuery] = useState('');
  const [protocolFilter, setProtocolFilter] = useState<'all' | ProtocolType>('all');
  const [statusFilter, setStatusFilter] = useState<'all' | DeviceStatus>('all');
  const [draft, setDraft] = useState<DeviceDraft>(() => createEmptyDeviceDraft());
  const [draftTestResult, setDraftTestResult] = useState<ConnectionTestResult | null>(null);
  const [fieldErrors, setFieldErrors] = useState<FieldErrorMap>({});
  const [notice, setNotice] = useState<DeviceNotice | null>(null);
  const pendingSelectedDeviceIdRef = useRef<string | null>(null);
  const initializedPanelKeyRef = useRef<string | null>(null);

  const filteredDevices = useMemo(() => {
    return devices.filter((device) => {
      const matchesSearch =
        searchQuery.trim() === '' ||
        device.name.toLowerCase().includes(searchQuery.trim().toLowerCase()) ||
        device.description.toLowerCase().includes(searchQuery.trim().toLowerCase());
      const matchesProtocol =
        protocolFilter === 'all' || device.protocol === protocolFilter;
      const matchesStatus =
        statusFilter === 'all' || device.status === statusFilter;

      return matchesSearch && matchesProtocol && matchesStatus;
    });
  }, [devices, protocolFilter, searchQuery, statusFilter]);

  const selectedDevice = useMemo(
    () => devices.find((device) => device.id === selectedDeviceId) ?? null,
    [devices, selectedDeviceId],
  );

  const selectedDeviceConnectionConfig = useMemo(
    () =>
      selectedDevice
        ? parseDeviceConnectionConfig(selectedDevice.connection_config)
        : null,
    [selectedDevice],
  );

  const selectedDeviceEndpoint = useMemo(
    () =>
      selectedDevice && selectedDeviceConnectionConfig
        ? buildDeviceEndpointSummary(selectedDevice.protocol, selectedDeviceConnectionConfig)
        : null,
    [selectedDevice, selectedDeviceConnectionConfig],
  );

  const selectedDeviceCapabilitySummary = useMemo(
    () =>
      selectedDevice && selectedDeviceConnectionConfig
        ? buildDeviceCapabilitySummary(
            selectedDevice.protocol,
            selectedDeviceConnectionConfig,
            t,
          )
        : [],
    [selectedDevice, selectedDeviceConnectionConfig, t],
  );

  useEffect(() => {
    if (!selectedDeviceId) {
      pendingSelectedDeviceIdRef.current = null;
      return;
    }

    const shouldClear = shouldClearSelectedDevice({
      selectedDeviceId,
      devices,
      isSuccess: devicesQuerySuccess,
      isFetching: devicesQueryFetching,
      pendingSelectedDeviceId: pendingSelectedDeviceIdRef.current,
    });

    if (!shouldClear) {
      if (pendingSelectedDeviceIdRef.current === selectedDeviceId) {
        pendingSelectedDeviceIdRef.current = null;
      }
      return;
    }

    setSelectedDeviceId(null);
    clearInspectorSelection();
  }, [
    clearInspectorSelection,
    devices,
    devicesQueryFetching,
    devicesQuerySuccess,
    selectedDeviceId,
    setSelectedDeviceId,
  ]);

  useEffect(() => {
    if (!devicePanelState) {
      initializedPanelKeyRef.current = null;
      return;
    }

    const panelKey =
      devicePanelState.mode === 'create'
        ? 'create'
        : devicePanelState.mode === 'edit'
          ? `edit:${devicePanelState.deviceId}`
          : `clone:${devicePanelState.sourceDeviceId}`;

    if (initializedPanelKeyRef.current === panelKey) {
      return;
    }

    if (devicePanelState.mode === 'create') {
      setDraft(createEmptyDeviceDraft());
      setFieldErrors({});
      initializedPanelKeyRef.current = panelKey;
      return;
    }

    const sourceDevice =
      devicePanelState.mode === 'edit'
        ? devices.find((device) => device.id === devicePanelState.deviceId)
        : devices.find((device) => device.id === devicePanelState.sourceDeviceId);

    if (!sourceDevice) {
      return;
    }

    const nextDraft = buildDeviceDraftFromDevice(sourceDevice);
    setDraft(
      devicePanelState.mode === 'clone'
        ? {
            ...nextDraft,
            name: '',
          }
        : nextDraft,
    );
    setFieldErrors({});
    initializedPanelKeyRef.current = panelKey;
  }, [devicePanelState, devices]);

  useEffect(() => {
    setDraftTestResult(null);
  }, [draft]);

  const isSaving =
    createDeviceMutation.isPending || updateDeviceMutation.isPending;

  const protocolOptions: SelectOption[] = useMemo(
    () => [
      { value: 'all', label: t('workbench.device.filters.allProtocols') },
      ...WORKBENCH_PROTOCOLS.map((protocol) => ({
        value: protocol,
        label: t(getWorkbenchProtocolLabelKey(protocol)),
      })),
    ],
    [t],
  );

  const statusOptions: SelectOption[] = useMemo(
    () => [
      { value: 'all', label: t('workbench.device.filters.allStatuses') },
      ...WORKBENCH_DEVICE_STATUSES.map((status) => ({
        value: status,
        label: t(getWorkbenchDeviceStatusLabelKey(status)),
      })),
    ],
    [t],
  );

  const connectionValueAsString = (key: string) => {
    const value = draft.connectionConfig[key];
    if (typeof value === 'string') {
      return value;
    }

    if (typeof value === 'number') {
      return String(value);
    }

    if (Array.isArray(value)) {
      return value.join('\n');
    }

    return '';
  };

  const setConnectionValue = (
    key: string,
    value: DeviceConnectionConfig[string],
  ) => {
    setDraft((currentDraft) => ({
      ...currentDraft,
      connectionConfig: {
        ...currentDraft.connectionConfig,
        [key]: value,
      },
    }));
  };

  const setConnectionNumber = (key: string, rawValue: string) => {
    if (rawValue === '') {
      setConnectionValue(key, undefined);
      return;
    }

    if (/^\d+$/.test(rawValue)) {
      setConnectionValue(key, Number.parseInt(rawValue, 10));
    }
  };

  const setTopics = (rawValue: string) => {
    setConnectionValue(
      'topics',
      rawValue
        .split(/\n|,/)
        .map((topic) => topic.trim())
        .filter((topic) => topic.length > 0),
    );
  };

  const closePanel = () => {
    closeDevicePanel();
    setFieldErrors({});
    setDraftTestResult(null);
  };

  const handleRefreshDevices = async () => {
    await queryClient.invalidateQueries({ queryKey: deviceKeys.lists() });
    setNotice({
      tone: 'info',
      message: t('workbench.device.messages.refreshed'),
    });
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const nextErrors = validateDeviceDraft(draft, t);
    if (devicePanelState?.mode === 'clone') {
      const sourceDevice = devices.find(
        (device) => device.id === devicePanelState.sourceDeviceId,
      );
      if (sourceDevice && draft.name.trim() === sourceDevice.name.trim()) {
        nextErrors.name = t('workbench.device.validation.cloneNameDistinct');
      }
    }
    setFieldErrors(nextErrors);
    if (Object.keys(nextErrors).length > 0) {
      return;
    }

    try {
      if (devicePanelState?.mode === 'create' || devicePanelState?.mode === 'clone') {
        const createdDevice = await createDeviceMutation.mutateAsync({
          name: draft.name.trim(),
          description: draft.description.trim() || undefined,
          protocol: draft.protocol,
          connection_config: sanitizeDeviceConnectionConfig(draft.connectionConfig),
        });

        pendingSelectedDeviceIdRef.current = createdDevice.id;
        setSelectedDeviceId(createdDevice.id);
        setNotice({
          tone: 'success',
          message: t('workbench.device.messages.createSuccess'),
        });
      }

      if (devicePanelState?.mode === 'edit') {
        await updateDeviceMutation.mutateAsync({
          id: devicePanelState.deviceId,
          data: {
            name: draft.name.trim(),
            description: draft.description.trim(),
            connection_config: sanitizeDeviceConnectionConfig(draft.connectionConfig),
          },
        });

        setSelectedDeviceId(devicePanelState.deviceId);
        setNotice({
          tone: 'success',
          message: t('workbench.device.messages.updateSuccess'),
        });
      }

      closePanel();
    } catch (error) {
      setNotice({
        tone: 'error',
        message:
          error instanceof Error
            ? error.message
            : t('workbench.device.messages.saveFailed'),
      });
    }
  };

  const handleTestDraftConnection = async () => {
    if (testDraftConnectionMutation.isPending) {
      return;
    }

    const nextErrors = validateDeviceConnectionFields(draft, t);
    setFieldErrors(nextErrors);
    if (Object.keys(nextErrors).length > 0) {
      return;
    }

    try {
      const result = await testDraftConnectionMutation.mutateAsync({
        protocol: draft.protocol,
        connection_config: sanitizeDeviceConnectionConfig(draft.connectionConfig),
      });
      setDraftTestResult(result);
      setNotice({
        tone: result.success ? 'success' : 'error',
        message: result.success
          ? t('workbench.device.messages.draftTestSuccess')
          : result.error || t('workbench.device.messages.draftTestFailed'),
      });
    } catch (error) {
      setDraftTestResult(null);
      setNotice({
        tone: 'error',
        message:
          error instanceof Error
            ? error.message
            : t('workbench.device.messages.draftTestFailed'),
      });
    }
  };

  const renderConnectionFields = () => {
    switch (draft.protocol) {
      case 'modbus_tcp':
      case 'modbus_udp':
        return (
          <>
            <div className="grid gap-4 md:grid-cols-2">
              <TextField
                id="device-host"
                label={t('workbench.device.connection.host')}
                onChange={(value) => setConnectionValue('host', value)}
                placeholder="192.168.1.10"
                value={connectionValueAsString('host')}
                error={fieldErrors.host}
              />
              <TextField
                id="device-port"
                label={t('workbench.device.connection.port')}
                onChange={(value) => setConnectionNumber('port', value)}
                placeholder="502"
                value={connectionValueAsString('port')}
                error={fieldErrors.port}
                inputMode="numeric"
              />
            </div>
            <div className="grid gap-4 md:grid-cols-2">
              <TextField
                id="device-slave-id"
                label={t('workbench.device.connection.slaveId')}
                onChange={(value) => setConnectionNumber('slave_id', value)}
                placeholder="1"
                value={connectionValueAsString('slave_id')}
                inputMode="numeric"
              />
              <TextField
                id="device-timeout"
                label={t('workbench.device.connection.timeout')}
                onChange={(value) => setConnectionNumber('timeout', value)}
                placeholder="5"
                value={connectionValueAsString('timeout')}
                inputMode="numeric"
              />
            </div>
          </>
        );
      case 'modbus_rtu':
        return (
          <>
            <TextField
              id="device-serial-port"
              label={t('workbench.device.connection.serialPort')}
              onChange={(value) => setConnectionValue('serial_port', value)}
              placeholder="/dev/ttyUSB0"
              value={connectionValueAsString('serial_port')}
              error={fieldErrors.serial_port}
            />
            <div className="grid gap-4 md:grid-cols-2">
              <SelectField
                id="device-baud-rate"
                label={t('workbench.device.connection.baudRate')}
                onChange={(value) =>
                  setConnectionValue('baud_rate', Number.parseInt(value, 10))
                }
                options={[9600, 19200, 38400, 57600, 115200].map((baudRate) => ({
                  value: String(baudRate),
                  label: String(baudRate),
                }))}
                value={connectionValueAsString('baud_rate') || '9600'}
              />
              <SelectField
                id="device-data-bits"
                label={t('workbench.device.connection.dataBits')}
                onChange={(value) =>
                  setConnectionValue('data_bits', Number.parseInt(value, 10))
                }
                options={[
                  { value: '7', label: '7' },
                  { value: '8', label: '8' },
                ]}
                value={connectionValueAsString('data_bits') || '8'}
              />
              <SelectField
                id="device-stop-bits"
                label={t('workbench.device.connection.stopBits')}
                onChange={(value) =>
                  setConnectionValue('stop_bits', Number.parseInt(value, 10))
                }
                options={[
                  { value: '1', label: '1' },
                  { value: '2', label: '2' },
                ]}
                value={connectionValueAsString('stop_bits') || '1'}
              />
              <SelectField
                id="device-parity"
                label={t('workbench.device.connection.parity')}
                onChange={(value) => setConnectionValue('parity', value)}
                options={[
                  { value: 'none', label: t('device.parityNone') },
                  { value: 'even', label: t('device.parityEven') },
                  { value: 'odd', label: t('device.parityOdd') },
                ]}
                value={connectionValueAsString('parity') || 'none'}
              />
            </div>
            <div className="grid gap-4 md:grid-cols-2">
              <TextField
                id="device-rtu-slave-id"
                label={t('workbench.device.connection.slaveId')}
                onChange={(value) => setConnectionNumber('slave_id', value)}
                placeholder="1"
                value={connectionValueAsString('slave_id')}
                inputMode="numeric"
              />
              <TextField
                id="device-rtu-timeout"
                label={t('workbench.device.connection.timeout')}
                onChange={(value) => setConnectionNumber('timeout', value)}
                placeholder="5"
                value={connectionValueAsString('timeout')}
                inputMode="numeric"
              />
            </div>
          </>
        );
      case 'fatek_fbs': {
        const mode =
          typeof draft.connectionConfig.mode === 'string'
            ? draft.connectionConfig.mode
            : 'tcp';

        return (
          <>
            <SelectField
              id="device-fatek-mode"
              label={t('workbench.device.connection.mode')}
              onChange={(value) => {
                if (value === 'tcp') {
                  setDraft((currentDraft) => {
                    const {
                      serial_port: _serialPort,
                      baud_rate: _baudRate,
                      data_bits: _dataBits,
                      stop_bits: _stopBits,
                      parity: _parity,
                      ...rest
                    } = currentDraft.connectionConfig;
                    const base = createDefaultDeviceConnectionConfig('fatek_fbs');
                    return {
                      ...currentDraft,
                      connectionConfig: {
                        ...base,
                        ...rest,
                        mode: 'tcp',
                      },
                    };
                  });
                  return;
                }
                setConnectionValue('mode', value);
                if (value === 'serial') {
                  setDraft((currentDraft) => ({
                    ...currentDraft,
                    connectionConfig: {
                      ...currentDraft.connectionConfig,
                      mode: value,
                      serial_port:
                        typeof currentDraft.connectionConfig.serial_port === 'string'
                          ? currentDraft.connectionConfig.serial_port
                          : '',
                      baud_rate:
                        typeof currentDraft.connectionConfig.baud_rate === 'number'
                          ? currentDraft.connectionConfig.baud_rate
                          : 9600,
                      data_bits:
                        typeof currentDraft.connectionConfig.data_bits === 'number'
                          ? currentDraft.connectionConfig.data_bits
                          : 7,
                      stop_bits:
                        typeof currentDraft.connectionConfig.stop_bits === 'number'
                          ? currentDraft.connectionConfig.stop_bits
                          : 1,
                      parity:
                        typeof currentDraft.connectionConfig.parity === 'string'
                          ? currentDraft.connectionConfig.parity
                          : 'even',
                    },
                  }));
                }
              }}
              options={[
                { value: 'tcp', label: t('workbench.device.mode.tcp') },
                { value: 'serial', label: t('workbench.device.mode.serial') },
              ]}
              value={mode}
            />
            {mode === 'serial' ? (
              <>
                <div className="grid gap-4 md:grid-cols-2">
                  <TextField
                    id="device-fatek-serial-port"
                    label={t('workbench.device.connection.serialPort')}
                    onChange={(value) => setConnectionValue('serial_port', value)}
                    placeholder="/dev/ttyUSB0"
                    value={connectionValueAsString('serial_port')}
                    error={fieldErrors.serial_port}
                  />
                  <SelectField
                    id="device-fatek-baud-rate"
                    label={t('workbench.device.connection.baudRate')}
                    onChange={(value) =>
                      setConnectionValue('baud_rate', Number.parseInt(value, 10))
                    }
                    options={[9600, 19200, 38400, 57600, 115200].map((baudRate) => ({
                      value: String(baudRate),
                      label: String(baudRate),
                    }))}
                    value={connectionValueAsString('baud_rate') || '9600'}
                  />
                </div>
                <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                  <SelectField
                    id="device-fatek-data-bits"
                    label={t('workbench.device.connection.dataBits')}
                    onChange={(value) =>
                      setConnectionValue('data_bits', Number.parseInt(value, 10))
                    }
                    options={[
                      { value: '7', label: '7' },
                      { value: '8', label: '8' },
                    ]}
                    value={connectionValueAsString('data_bits') || '7'}
                  />
                  <SelectField
                    id="device-fatek-stop-bits"
                    label={t('workbench.device.connection.stopBits')}
                    onChange={(value) =>
                      setConnectionValue('stop_bits', Number.parseInt(value, 10))
                    }
                    options={[
                      { value: '1', label: '1' },
                      { value: '2', label: '2' },
                    ]}
                    value={connectionValueAsString('stop_bits') || '1'}
                  />
                  <SelectField
                    id="device-fatek-parity"
                    label={t('workbench.device.connection.parity')}
                    onChange={(value) => setConnectionValue('parity', value)}
                    options={[
                      { value: 'none', label: t('device.parityNone') },
                      { value: 'even', label: t('device.parityEven') },
                      { value: 'odd', label: t('device.parityOdd') },
                    ]}
                    value={connectionValueAsString('parity') || 'even'}
                  />
                </div>
              </>
            ) : (
              <div className="grid gap-4 md:grid-cols-2">
                <TextField
                  id="device-fatek-host"
                  label={t('workbench.device.connection.host')}
                  onChange={(value) => setConnectionValue('host', value)}
                  placeholder="192.168.1.10"
                  value={connectionValueAsString('host')}
                  error={fieldErrors.host}
                />
                <TextField
                  id="device-fatek-port"
                  label={t('workbench.device.connection.port')}
                  onChange={(value) => setConnectionNumber('port', value)}
                  placeholder="500"
                  value={connectionValueAsString('port')}
                  error={fieldErrors.port}
                  inputMode="numeric"
                />
              </div>
            )}
            <div className="grid gap-4 md:grid-cols-2">
              <TextField
                id="device-fatek-station"
                label={t('workbench.device.connection.stationNo')}
                onChange={(value) => setConnectionNumber('station_no', value)}
                placeholder="1"
                value={connectionValueAsString('station_no')}
                error={fieldErrors.station_no}
                inputMode="numeric"
              />
              <TextField
                id="device-fatek-timeout"
                label={t('workbench.device.connection.timeout')}
                onChange={(value) => setConnectionNumber('timeout', value)}
                placeholder="5"
                value={connectionValueAsString('timeout')}
                inputMode="numeric"
              />
            </div>
          </>
        );
      }
      case 'mc_3e':
        return (
          <>
            <div className="grid gap-4 md:grid-cols-2">
              <TextField
                id="device-mc-host"
                label={t('workbench.device.connection.host')}
                onChange={(value) => setConnectionValue('host', value)}
                placeholder="192.168.1.10"
                value={connectionValueAsString('host')}
                error={fieldErrors.host}
              />
              <TextField
                id="device-mc-port"
                label={t('workbench.device.connection.port')}
                onChange={(value) => setConnectionNumber('port', value)}
                placeholder="5000"
                value={connectionValueAsString('port')}
                error={fieldErrors.port}
                inputMode="numeric"
              />
            </div>
            <div className="grid gap-4 md:grid-cols-2">
              <TextField
                id="device-mc-network"
                label={t('workbench.device.connection.networkNo')}
                onChange={(value) => setConnectionNumber('network_no', value)}
                placeholder="0"
                value={connectionValueAsString('network_no')}
                inputMode="numeric"
              />
              <TextField
                id="device-mc-pc"
                label={t('workbench.device.connection.pcNo')}
                onChange={(value) => setConnectionNumber('pc_no', value)}
                placeholder="255"
                value={connectionValueAsString('pc_no')}
                inputMode="numeric"
              />
              <TextField
                id="device-mc-io"
                label={t('workbench.device.connection.ioNo')}
                onChange={(value) => setConnectionNumber('io_no', value)}
                placeholder="1023"
                value={connectionValueAsString('io_no')}
                inputMode="numeric"
              />
              <TextField
                id="device-mc-station"
                label={t('workbench.device.connection.stationNo')}
                onChange={(value) => setConnectionNumber('station_no', value)}
                placeholder="0"
                value={connectionValueAsString('station_no')}
                inputMode="numeric"
              />
            </div>
            <div className="grid gap-4 md:grid-cols-2">
              <TextField
                id="device-mc-timeout"
                label={t('workbench.device.connection.timeout')}
                onChange={(value) => setConnectionNumber('timeout', value)}
                placeholder="5"
                value={connectionValueAsString('timeout')}
                inputMode="numeric"
              />
              <SelectField
                id="device-mc-format"
                label={t('workbench.device.connection.dataFormat')}
                onChange={(value) => setConnectionValue('data_format', value)}
                options={[
                  {
                    value: 'ABCD',
                    label: t('workbench.device.connection.dataFormats.abcd'),
                  },
                  {
                    value: 'BADC',
                    label: t('workbench.device.connection.dataFormats.badc'),
                  },
                  {
                    value: 'CDAB',
                    label: t('workbench.device.connection.dataFormats.cdab'),
                  },
                  {
                    value: 'DCBA',
                    label: t('workbench.device.connection.dataFormats.dcba'),
                  },
                ]}
                value={normalizeMc3eDataFormatValue(
                  draft.connectionConfig.data_format,
                )}
              />
            </div>
          </>
        );
      case 'mqtt':
        return (
          <>
            <div className="grid gap-4 md:grid-cols-2">
              <TextField
                id="device-mqtt-broker"
                label={t('workbench.device.connection.brokerUrl')}
                onChange={(value) => setConnectionValue('broker_url', value)}
                placeholder="mqtt://broker.local:1883"
                value={connectionValueAsString('broker_url')}
                error={fieldErrors.broker_url}
              />
              <TextField
                id="device-mqtt-client-id"
                label={t('workbench.device.connection.clientId')}
                onChange={(value) => setConnectionValue('client_id', value)}
                placeholder="datalink-workbench"
                value={connectionValueAsString('client_id')}
                error={fieldErrors.client_id}
              />
            </div>
            <div className="grid gap-4 md:grid-cols-2">
              <TextField
                id="device-mqtt-username"
                label={t('workbench.device.connection.username')}
                onChange={(value) => setConnectionValue('username', value)}
                value={connectionValueAsString('username')}
              />
              <TextField
                id="device-mqtt-password"
                label={t('workbench.device.connection.password')}
                onChange={(value) => setConnectionValue('password', value)}
                type="password"
                value={connectionValueAsString('password')}
              />
            </div>
            <TextAreaField
              id="device-mqtt-topics"
              label={t('workbench.device.connection.topics')}
              onChange={setTopics}
              placeholder="factory/line1/temperature"
              value={connectionValueAsString('topics')}
              error={fieldErrors.topics}
            />
            <div className="grid gap-4 md:grid-cols-2">
              <SelectField
                id="device-mqtt-qos"
                label={t('workbench.device.connection.qos')}
                onChange={(value) => setConnectionValue('qos', Number.parseInt(value, 10))}
                options={[
                  { value: '0', label: '0' },
                  { value: '1', label: '1' },
                  { value: '2', label: '2' },
                ]}
                value={connectionValueAsString('qos') || '0'}
              />
              <CheckboxField
                id="device-mqtt-tls"
                checked={Boolean(draft.connectionConfig.use_tls)}
                label={t('workbench.device.connection.useTls')}
                onChange={(checked) => setConnectionValue('use_tls', checked)}
              />
            </div>
          </>
        );
    }
  };

  const renderInlineEditor = () => {
    if (!devicePanelState) return null;
    return (
      <div
        className="flex flex-col overflow-hidden rounded-3xl border border-slate-800 bg-slate-950/40"
        data-testid="device-inline-editor"
      >
        <div className="flex items-start justify-between gap-4 border-b border-slate-800 px-6 py-5">
          <div className="space-y-2">
            <p className="text-xs font-semibold uppercase tracking-[0.24em] text-cyan-300">
              {t('workbench.device.panel.eyebrow')}
            </p>
            <h3 className="text-2xl font-semibold text-slate-50">
              {devicePanelState.mode === 'create'
                ? t('workbench.device.panel.createTitle')
                : devicePanelState.mode === 'clone'
                  ? t('workbench.device.panel.cloneTitle')
                  : t('workbench.device.panel.editTitle')}
            </h3>
            <p className="text-sm text-slate-300">
              {devicePanelState.mode === 'create'
                ? t('workbench.device.panel.createDescription')
                : devicePanelState.mode === 'clone'
                  ? t('workbench.device.panel.cloneDescription')
                  : t('workbench.device.panel.editDescription')}
            </p>
          </div>
          <button
            className={ghostButtonClassName}
            onClick={closePanel}
            type="button"
          >
            {t('common.cancel')}
          </button>
        </div>

        <form className="flex min-h-0 flex-1 flex-col" onSubmit={handleSubmit}>
          <div className="min-h-0 flex-1 space-y-6 overflow-y-auto px-6 py-6">
            <div className="grid gap-4 md:grid-cols-2">
              <TextField
                id="device-name"
                label={t('workbench.device.fields.name')}
                onChange={(value) =>
                  setDraft((currentDraft) => ({
                    ...currentDraft,
                    name: value,
                  }))
                }
                placeholder="Mixer PLC"
                value={draft.name}
                error={fieldErrors.name}
              />
              <SelectField
                id="device-protocol"
                label={t('workbench.device.fields.protocol')}
                onChange={(value) => {
                  const nextProtocol = value as ProtocolType;
                  setDraft((currentDraft) => ({
                    ...currentDraft,
                    protocol: nextProtocol,
                    connectionConfig: createDefaultDeviceConnectionConfig(nextProtocol),
                  }));
                }}
                options={WORKBENCH_PROTOCOLS.map((protocol) => ({
                  value: protocol,
                  label: t(getWorkbenchProtocolLabelKey(protocol)),
                }))}
                value={draft.protocol}
                disabled={devicePanelState.mode !== 'create'}
              />
              <div className="md:col-span-2">
                <TextAreaField
                  id="device-description"
                  label={t('workbench.device.fields.description')}
                  onChange={(value) =>
                    setDraft((currentDraft) => ({
                      ...currentDraft,
                      description: value,
                    }))
                  }
                  placeholder={t('workbench.device.fields.descriptionPlaceholder')}
                  value={draft.description}
                />
              </div>
            </div>

            <section className="space-y-4 rounded-3xl border border-slate-800 bg-slate-900/50 p-5">
              <div className="space-y-2">
                <p className="text-xs font-semibold uppercase tracking-[0.24em] text-cyan-300">
                  {t('workbench.device.connection.eyebrow')}
                </p>
                <h4 className="text-xl font-semibold text-slate-50">
                  {t('workbench.device.connection.title')}
                </h4>
                <p className="text-sm text-slate-300">
                  {t('workbench.device.connection.description')}
                </p>
                <div className="rounded-2xl border border-cyan-500/20 bg-slate-950/50 px-4 py-3 text-xs text-slate-300">
                  {t('workbench.device.connection.backendHostHint')}
                </div>
              </div>

              {renderConnectionFields()}

              {draftTestResult ? (
                <div
                  className="space-y-3 rounded-2xl border border-slate-800 bg-slate-950/60 p-4"
                  data-testid="device-draft-test-result"
                >
                  <div className="flex flex-wrap items-center justify-between gap-2">
                    <p className="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">
                      {t('workbench.device.actions.testDraftConnection')}
                    </p>
                    <span
                      className={joinClasses(
                        'rounded-full border px-2 py-0.5 text-[11px] font-medium',
                        getDraftStageStatusClasses(
                          draftTestResult.success ? 'success' : 'failed',
                        ),
                      )}
                    >
                      {draftTestResult.success
                        ? t('workbench.device.inspector.phaseStatus.success')
                        : t('workbench.device.inspector.phaseStatus.failed')}
                    </span>
                  </div>
                  <div className="grid gap-3 md:grid-cols-2">
                    {([
                      ['connect', draftTestResult.connect],
                      ['probe', draftTestResult.probe],
                    ] as const).map(([stageKey, stage]) => (
                      <div
                        className="rounded-2xl border border-slate-800 bg-slate-900/80 p-3"
                        key={stageKey}
                      >
                        <div className="flex items-center justify-between gap-2">
                          <span className="text-xs font-semibold uppercase tracking-[0.18em] text-slate-500">
                            {t(`workbench.device.inspector.phases.${stageKey}`)}
                          </span>
                          <span
                            className={joinClasses(
                              'rounded-full border px-2 py-0.5 text-[11px] font-medium',
                              getDraftStageStatusClasses(stage?.status),
                            )}
                          >
                            {t(
                              `workbench.device.inspector.phaseStatus.${stage?.status ?? 'skipped'}`,
                            )}
                          </span>
                        </div>
                        <p className="mt-2 text-sm text-slate-200">
                          {getDraftStageMessage(stage)}
                        </p>
                      </div>
                    ))}
                  </div>
                  {draftTestResult.can_activate === false ? (
                    <p className="text-xs font-medium text-amber-300">
                      {t('workbench.device.inspector.activationBlocked')}
                    </p>
                  ) : null}
                </div>
              ) : null}
            </section>
          </div>

          <div className="flex flex-wrap justify-end gap-3 border-t border-slate-800 px-6 py-4">
            <button
              className={ghostButtonClassName}
              disabled={testDraftConnectionMutation.isPending}
              onClick={() => void handleTestDraftConnection()}
              type="button"
            >
              {testDraftConnectionMutation.isPending
                ? t('workbench.device.actions.testingDraft')
                : t('workbench.device.actions.testDraftConnection')}
            </button>
            <button
              className={ghostButtonClassName}
              onClick={closePanel}
              type="button"
            >
              {t('common.cancel')}
            </button>
            <button
              className={primaryButtonClassName}
              disabled={isSaving}
              type="submit"
            >
              {isSaving
                ? t('workbench.device.actions.saving')
                : t('workbench.device.actions.save')}
            </button>
          </div>
        </form>
      </div>
    );
  };

  return (
    <section className="flex h-full min-h-0 flex-col overflow-hidden">
      <div className="min-h-0 flex-1 space-y-4 overflow-y-auto overscroll-contain">
      {devices.length > 0 ? (
        <div
          className="grid gap-3 rounded-3xl border border-slate-800 bg-slate-950/40 p-4 xl:grid-cols-[minmax(0,1fr)_180px_180px_auto] xl:items-end"
          data-testid="device-primary-toolbar"
        >
          <TextField
            id="workbench-device-search"
            label={t('workbench.device.search.label')}
            onChange={setSearchQuery}
            placeholder={t('workbench.device.search.placeholder')}
            value={searchQuery}
          />
          <SelectField
            id="workbench-device-protocol-filter"
            label={t('workbench.device.filters.protocol')}
            onChange={(value) =>
              setProtocolFilter(value === 'all' ? 'all' : (value as ProtocolType))
            }
            options={protocolOptions}
            value={protocolFilter}
          />
          <SelectField
            id="workbench-device-status-filter"
            label={t('workbench.device.filters.status')}
            onChange={(value) =>
              setStatusFilter(value === 'all' ? 'all' : (value as DeviceStatus))
            }
            options={statusOptions}
            value={statusFilter}
          />
          <div className="flex flex-wrap items-end justify-end gap-2">
            <button
              className={ghostButtonClassName}
              onClick={() => void handleRefreshDevices()}
              type="button"
            >
              {t('workbench.device.actions.refresh')}
            </button>
            <button
              className={primaryButtonClassName}
              onClick={openCreateDevicePanel}
              type="button"
            >
              {t('workbench.device.actions.create')}
            </button>
          </div>
        </div>
      ) : null}

      {notice ? (
        <div
          aria-live="polite"
          className={joinClasses(
            'rounded-2xl border px-4 py-3 text-sm',
            getNoticeClasses(notice.tone),
          )}
          role="status"
        >
          {notice.message}
        </div>
      ) : null}

      {isLoading ? (
        <div className="flex items-center justify-center gap-3 rounded-3xl border border-slate-800 bg-slate-950/40 p-8">
          <Spinner size="lg" />
          <span className="text-sm text-slate-400">{t('workbench.device.loading')}</span>
        </div>
      ) : null}

      {!isLoading && devices.length === 0 && !devicePanelState ? (
        <div className="rounded-3xl border border-dashed border-slate-700 bg-slate-950/40 p-8 text-center">
          <p className="text-xs font-semibold uppercase tracking-[0.24em] text-cyan-300">
            {t('workbench.device.empty.eyebrow')}
          </p>
          <h3 className="mt-3 text-2xl font-semibold text-slate-50">
            {t('workbench.device.empty.title')}
          </h3>
          <p className="mx-auto mt-3 max-w-2xl text-sm text-slate-300">
            {t('workbench.device.empty.description')}
          </p>
          <button
            className={joinClasses(primaryButtonClassName, 'mt-6')}
            onClick={openCreateDevicePanel}
            type="button"
          >
            {t('workbench.device.actions.create')}
          </button>
        </div>
      ) : null}

      {!isLoading && devices.length > 0 && filteredDevices.length === 0 && !devicePanelState ? (
        <div className="rounded-3xl border border-slate-800 bg-slate-950/40 p-8 text-center text-sm text-slate-300">
          <p className="text-xs font-semibold uppercase tracking-[0.24em] text-cyan-300">
            {t('workbench.device.filteredEmpty.eyebrow')}
          </p>
          <h3 className="mt-3 text-2xl font-semibold text-slate-50">
            {t('workbench.device.filteredEmpty.title')}
          </h3>
          <p className="mt-3">{t('workbench.device.filteredEmpty.description')}</p>
        </div>
      ) : null}

      {!isLoading && filteredDevices.length === 0 && devicePanelState ? (
        <div className="mx-auto max-w-4xl">
          {renderInlineEditor()}
        </div>
      ) : null}

      {!isLoading && filteredDevices.length > 0 ? (
        <div className="grid gap-4 xl:grid-cols-[minmax(380px,0.78fr)_minmax(0,1.22fr)]">
          <div className="space-y-3">
            {filteredDevices.map((device) => {
              const isSelected = device.id === selectedDeviceId;
              const connectionConfig = parseDeviceConnectionConfig(device.connection_config);
              const capabilitySummary = buildDeviceCapabilitySummary(
                device.protocol,
                connectionConfig,
                t,
              );
              const endpoint = buildDeviceEndpointSummary(device.protocol, connectionConfig);
              const capabilityHints = capabilitySummary.filter(
                (item) => item.id === 'unit-id' || item.id === 'address-base',
              );

              return (
                <button
                  data-testid={`device-row-${device.id}`}
                  key={device.id}
                  aria-label={device.name}
                  aria-pressed={isSelected}
                  className={joinClasses(
                    'flex w-full flex-col gap-3 rounded-2xl border px-4 py-3 text-left transition',
                    isSelected
                      ? 'border-cyan-400 bg-cyan-500/10 shadow-lg shadow-cyan-950/20'
                      : 'border-slate-800 bg-slate-950/40 hover:border-slate-600 hover:bg-slate-900/70',
                  )}
                  onClick={() => {
                    setSelectedDeviceId(device.id);
                    setInspectorSelection({ kind: 'device', deviceId: device.id });
                    setNotice(null);
                  }}
                  type="button"
                >
                  <div className="min-w-0 space-y-1.5">
                    <div className="flex flex-wrap items-center gap-2">
                      <span className="truncate text-base font-semibold text-slate-50">
                        {device.name}
                      </span>
                      <span className="rounded-full border border-slate-700 bg-slate-900/80 px-2 py-1 text-xs text-slate-300">
                        {t(getWorkbenchProtocolLabelKey(device.protocol))}
                      </span>
                    </div>
                    <p
                      className="truncate text-sm text-slate-400"
                      data-testid={`device-endpoint-${device.id}`}
                    >
                      {endpoint}
                    </p>
                  </div>

                  <div
                    className="flex min-w-0 flex-wrap items-center gap-2 sm:justify-between"
                    data-testid={`device-capability-strip-${device.id}`}
                  >
                    <div className="flex min-w-0 flex-wrap items-center gap-2">
                      {capabilityHints.map((item) => (
                        <span
                          className="inline-flex shrink-0 items-center gap-x-2 gap-y-0.5 rounded-2xl border border-slate-700 bg-slate-950/80 px-3 py-1.5 text-xs text-slate-300"
                          key={`${device.id}-${item.id}`}
                        >
                          <span className="whitespace-nowrap uppercase tracking-[0.18em] text-slate-500">
                            {t(item.labelKey)}
                          </span>
                          <span className="font-medium whitespace-nowrap text-slate-100">
                            {item.value}
                          </span>
                        </span>
                      ))}
                    </div>
                    <div className="flex flex-wrap items-center gap-2 sm:justify-end">
                      <span
                        className={joinClasses(
                          'shrink-0 rounded-full border px-2 py-1 text-xs font-medium whitespace-nowrap',
                          getStatusClasses(device.status),
                        )}
                      >
                        {t(getWorkbenchDeviceStatusLabelKey(device.status))}
                      </span>
                      <span
                        className={joinClasses(
                          'shrink-0 rounded-full border px-2.5 py-1 text-[11px] font-medium whitespace-nowrap',
                          getDeviceHealthClasses(device.last_test_success),
                        )}
                        data-testid={`device-health-${device.id}`}
                      >
                        {getDeviceHealthLabel(t, device)}
                      </span>
                    </div>
                  </div>
                </button>
              );
            })}
          </div>

          {devicePanelState ? (
            renderInlineEditor()
          ) : (
            <div
              className="rounded-3xl border border-slate-800 bg-slate-950/40 p-6"
              data-testid="device-detail-panel"
            >
              {selectedDevice ? (
                <div className="space-y-6">
                  <div className="space-y-2">
                    <p className="text-xs font-semibold uppercase tracking-[0.24em] text-cyan-300">
                      {t('workbench.device.eyebrow')}
                    </p>
                    <div className="flex flex-wrap items-center gap-2">
                      <h3 className="text-2xl font-semibold text-slate-50">{selectedDevice.name}</h3>
                      <span className="rounded-full border border-slate-700 bg-slate-900/80 px-2 py-1 text-xs text-slate-300">
                        {t(getWorkbenchProtocolLabelKey(selectedDevice.protocol))}
                      </span>
                      <span
                        className={joinClasses(
                          'rounded-full border px-2 py-1 text-xs font-medium',
                          getStatusClasses(selectedDevice.status),
                        )}
                      >
                        {t(getWorkbenchDeviceStatusLabelKey(selectedDevice.status))}
                      </span>
                    </div>
                    {selectedDeviceEndpoint ? (
                      <p
                        className="text-sm text-slate-300"
                        data-testid="device-detail-endpoint"
                      >
                        {selectedDeviceEndpoint}
                      </p>
                    ) : null}
                  </div>

                  <div className="flex min-w-0 flex-col gap-4 xl:flex-row xl:items-stretch">
                    <div className="w-full min-w-0 space-y-3 rounded-2xl border border-slate-800 bg-slate-900/50 p-4 xl:flex-1">
                      <p className="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">
                        {t('workbench.device.inspector.capabilitySummary')}
                      </p>
                      <div className="grid grid-cols-[repeat(auto-fit,minmax(min(100%,220px),1fr))] gap-3">
                        {selectedDeviceCapabilitySummary.map((item) => (
                          <div
                            className="min-w-0 rounded-2xl border border-slate-800 bg-slate-950/50 px-4 py-3"
                            key={`selected-${item.id}`}
                          >
                            <p className="text-[11px] uppercase tracking-[0.18em] text-slate-500 break-words">
                              {t(item.labelKey)}
                            </p>
                            <p className="mt-2 break-words text-sm font-medium text-slate-100">
                              {item.value}
                            </p>
                          </div>
                        ))}
                      </div>
                    </div>

                    <div className="w-full shrink-0 space-y-3 rounded-2xl border border-slate-800 bg-slate-900/50 p-4 xl:w-64 xl:max-w-sm">
                      <p className="text-xs font-semibold uppercase tracking-[0.18em] text-slate-400">
                        {t('workbench.device.card.lastTest')}
                      </p>
                      <span
                        className={joinClasses(
                          'inline-flex rounded-full border px-2.5 py-1 text-xs font-medium',
                          getDeviceHealthClasses(selectedDevice.last_test_success),
                        )}
                      >
                        {getDeviceHealthLabel(t, selectedDevice)}
                      </span>
                      <button
                        className={joinClasses(primaryButtonClassName, 'w-full')}
                        onClick={() => setActiveStep('source')}
                        type="button"
                      >
                        {t('workbench.device.actions.continue')}
                      </button>
                    </div>
                  </div>
                </div>
              ) : (
                <div className="space-y-3">
                  <p className="text-xs font-semibold uppercase tracking-[0.24em] text-cyan-300">
                    {t('workbench.device.eyebrow')}
                  </p>
                  <h3 className="text-2xl font-semibold text-slate-50">
                    {t('workbench.device.inspector.emptyTitle')}
                  </h3>
                  <p className="max-w-2xl text-sm text-slate-300">
                    {t('workbench.device.inspector.emptyDescription')}
                  </p>
                </div>
              )}
            </div>
          )}
        </div>
      ) : null}
      </div>
    </section>
  );
}
