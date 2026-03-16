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
  useUpdateDeviceMutation,
} from '../../../hooks/datalink/useDevices';
import { deviceKeys } from '../../../hooks/datalink/keys';
import type {
  Device,
  DeviceStatus,
  ProtocolType,
} from '../../../types/datalink';
import { useWorkbench } from './WorkbenchProvider';
import {
  buildDeviceCapabilitySummary,
  buildDeviceDraftFromDevice,
  createDefaultDeviceConnectionConfig,
  createEmptyDeviceDraft,
  getWorkbenchDeviceStatusLabelKey,
  getWorkbenchProtocolLabelKey,
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

function getDeviceTestLabel(t: (key: string) => string, device: Device) {
  if (device.last_test_success === true) {
    return t('workbench.device.card.testPassed');
  }

  if (device.last_test_success === false) {
    return device.last_test_error || t('workbench.device.card.testFailed');
  }

  return t('workbench.device.card.notTested');
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

function validateDeviceDraft(
  draft: DeviceDraft,
  t: (key: string) => string,
): FieldErrorMap {
  const errors: FieldErrorMap = {};

  if (draft.name.trim() === '') {
    errors.name = t('workbench.device.validation.nameRequired');
  }

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

export function WorkbenchDeviceStep() {
  const { t } = useTranslation();
  const queryClient = useQueryClient();
  const {
    clearInspectorSelection,
    closeDevicePanel,
    devicePanelState,
    openCloneDevicePanel,
    openCreateDevicePanel,
    selectedDeviceId,
    setActiveStep,
    setInspectorSelection,
    setSelectedDeviceId,
  } = useWorkbench();
  const { data: devices = [], isLoading } = useDevicesQuery();
  const createDeviceMutation = useCreateDeviceMutation();
  const updateDeviceMutation = useUpdateDeviceMutation();

  const [searchQuery, setSearchQuery] = useState('');
  const [protocolFilter, setProtocolFilter] = useState<'all' | ProtocolType>('all');
  const [statusFilter, setStatusFilter] = useState<'all' | DeviceStatus>('all');
  const [draft, setDraft] = useState<DeviceDraft>(() => createEmptyDeviceDraft());
  const [fieldErrors, setFieldErrors] = useState<FieldErrorMap>({});
  const [notice, setNotice] = useState<DeviceNotice | null>(null);
  const pendingSelectedDeviceIdRef = useRef<string | null>(null);

  const selectedDevice = useMemo(
    () => devices.find((device) => device.id === selectedDeviceId) ?? null,
    [devices, selectedDeviceId],
  );

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

  useEffect(() => {
    if (!selectedDeviceId) {
      pendingSelectedDeviceIdRef.current = null;
      return;
    }

    const deviceExists = devices.some((device) => device.id === selectedDeviceId);
    if (deviceExists) {
      if (pendingSelectedDeviceIdRef.current === selectedDeviceId) {
        pendingSelectedDeviceIdRef.current = null;
      }
      return;
    }

    if (pendingSelectedDeviceIdRef.current === selectedDeviceId) {
      return;
    }

    if (selectedDeviceId && !deviceExists) {
      setSelectedDeviceId(null);
      clearInspectorSelection();
    }
  }, [clearInspectorSelection, devices, selectedDeviceId, setSelectedDeviceId]);

  useEffect(() => {
    if (!devicePanelState) {
      return;
    }

    if (devicePanelState.mode === 'create') {
      setDraft(createEmptyDeviceDraft());
      setFieldErrors({});
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
  }, [devicePanelState, devices]);

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
                placeholder="1"
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
                    value: 'binary',
                    label: t('workbench.device.connection.dataFormats.binary'),
                  },
                  {
                    value: 'ascii',
                    label: t('workbench.device.connection.dataFormats.ascii'),
                  },
                ]}
                value={connectionValueAsString('data_format') || 'binary'}
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

  return (
    <section className="relative space-y-5">
      <div className="flex flex-col gap-4 rounded-3xl border border-slate-800 bg-slate-950/40 p-5">
        <div className="space-y-2">
          <p className="text-xs font-semibold uppercase tracking-[0.24em] text-cyan-300">
            {t('workbench.device.eyebrow')}
          </p>
          <h2 className="text-2xl font-semibold text-slate-50">
            {t('workbench.device.title')}
          </h2>
          <p className="max-w-3xl text-sm text-slate-300">
            {t('workbench.device.description')}
          </p>
        </div>
        {devices.length > 0 ? (
          <div className="flex flex-wrap gap-3">
            <button
              className={primaryButtonClassName}
              onClick={openCreateDevicePanel}
              type="button"
            >
              {t('workbench.device.actions.create')}
            </button>
            <button
              className={ghostButtonClassName}
              onClick={() => void handleRefreshDevices()}
              type="button"
            >
              {t('workbench.device.actions.refresh')}
            </button>
            <button
              className={ghostButtonClassName}
              disabled={!selectedDevice}
              onClick={() => selectedDevice && openCloneDevicePanel(selectedDevice.id)}
              type="button"
            >
              {t('workbench.device.actions.clone')}
            </button>
            {selectedDevice ? (
              <button
                className={ghostButtonClassName}
                onClick={() => setActiveStep('source')}
                type="button"
              >
                {t('workbench.device.actions.continue')}
              </button>
            ) : null}
          </div>
        ) : null}
      </div>

      <div className="grid gap-4 rounded-3xl border border-slate-800 bg-slate-950/40 p-5 md:grid-cols-[minmax(0,1fr)_180px_180px]">
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
      </div>

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
        <div className="rounded-3xl border border-slate-800 bg-slate-950/40 p-8 text-sm text-slate-400">
          {t('workbench.device.loading')}
        </div>
      ) : null}

      {!isLoading && devices.length === 0 ? (
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

      {!isLoading && devices.length > 0 && filteredDevices.length === 0 ? (
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

      {!isLoading && filteredDevices.length > 0 ? (
        <div className="grid gap-4 xl:grid-cols-2">
          {filteredDevices.map((device) => {
            const isSelected = device.id === selectedDeviceId;
            const capabilitySummary = buildDeviceCapabilitySummary(
              device.protocol,
              parseDeviceConnectionConfig(device.connection_config),
              t,
            );

            return (
              <button
                key={device.id}
                aria-label={device.name}
                aria-pressed={isSelected}
                className={joinClasses(
                  'flex flex-col gap-4 rounded-3xl border p-5 text-left transition',
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
                <div className="flex items-start justify-between gap-3">
                  <div className="space-y-2">
                    <div className="flex flex-wrap items-center gap-2">
                      <span className="text-lg font-semibold text-slate-50">
                        {device.name}
                      </span>
                      <span className="rounded-full border border-slate-700 bg-slate-900/80 px-2 py-1 text-xs text-slate-300">
                        {t(getWorkbenchProtocolLabelKey(device.protocol))}
                      </span>
                    </div>
                    <p className="text-sm text-slate-400">
                      {device.description || t('workbench.device.card.noDescription')}
                    </p>
                  </div>
                  <span
                    className={joinClasses(
                      'rounded-full border px-2 py-1 text-xs font-medium',
                      getStatusClasses(device.status),
                    )}
                  >
                    {t(getWorkbenchDeviceStatusLabelKey(device.status))}
                  </span>
                </div>
                <dl className="grid gap-3 text-sm text-slate-300 sm:grid-cols-2">
                  <div className="space-y-1">
                    <dt className={labelClassName}>
                      {t('workbench.device.card.protocol')}
                    </dt>
                    <dd>{t(getWorkbenchProtocolLabelKey(device.protocol))}</dd>
                  </div>
                  <div className="space-y-1">
                    <dt className={labelClassName}>
                      {t('workbench.device.card.lastTest')}
                    </dt>
                    <dd>{getDeviceTestLabel(t, device)}</dd>
                  </div>
                </dl>
                <div className="grid gap-3 sm:grid-cols-2">
                  {capabilitySummary.map((item) => (
                    <div
                      className="rounded-2xl border border-slate-800 bg-slate-950/70 px-3 py-3"
                      key={`${device.id}-${item.id}`}
                    >
                      <p className="text-[10px] uppercase tracking-[0.18em] text-slate-500">
                        {t(item.labelKey)}
                      </p>
                      <p className="mt-1 text-sm font-medium text-slate-100">
                        {item.value}
                      </p>
                    </div>
                  ))}
                </div>
              </button>
            );
          })}
        </div>
      ) : null}

      {devicePanelState ? (
        <div className="absolute inset-0 z-20 flex justify-end bg-slate-950/70 p-4 backdrop-blur-sm">
          <div
            aria-modal="true"
            className="flex h-full w-full max-w-2xl flex-col overflow-hidden rounded-3xl border border-slate-800 bg-slate-950 shadow-2xl shadow-slate-950/60"
            role="dialog"
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
                  </div>

                  {renderConnectionFields()}
                </section>
              </div>

              <div className="flex flex-wrap justify-end gap-3 border-t border-slate-800 px-6 py-4">
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
        </div>
      ) : null}
    </section>
  );
}
