import type {
  CreateDeviceRequest,
  UpdateDeviceRequest,
} from '../../../../types/datalink';
import type { StudioV2WorkspaceDeviceRecord } from '../../../../services/studioV2WorkspaceDevices';
import type { Device, ProtocolId } from './types';

function asNonEmptyString(value: unknown): string | null {
  return typeof value === 'string' && value.trim() !== '' ? value.trim() : null;
}

function asPositiveNumber(value: unknown): number | null {
  return typeof value === 'number' && Number.isFinite(value) && value > 0 ? value : null;
}

export function isStudioV2DeviceValid(device: Device): boolean {
  if (!asNonEmptyString(device.name)) {
    return false;
  }

  if (device.protocol === 'modbus_rtu') {
    return Boolean(
      asNonEmptyString(device.config.port) &&
      asPositiveNumber(device.config.baud) &&
      asNonEmptyString(device.config.parity) &&
      asPositiveNumber(device.config.slave_id),
    );
  }

  if (device.protocol === 'mqtt') {
    return Boolean(
      asNonEmptyString(device.config.broker) &&
      asNonEmptyString(device.config.client_id),
    );
  }

  return Boolean(
    asNonEmptyString(device.config.host) &&
    asPositiveNumber(device.config.port) &&
    asPositiveNumber(device.config.slave_id) &&
    asPositiveNumber(device.config.timeout),
  );
}

export function toStudioV2DeviceCreateRequest(device: Device): CreateDeviceRequest {
  return {
    id: device.id,
    name: device.name,
    description: device.description,
    protocol: device.protocol,
    connection_config: { ...device.config },
  };
}

export function toStudioV2DeviceUpdateRequest(device: Device): UpdateDeviceRequest {
  return {
    name: device.name,
    description: device.description,
    connection_config: { ...device.config },
  };
}

export function hydrateStudioV2Device(record: StudioV2WorkspaceDeviceRecord): Device {
  let config: Record<string, unknown> = {};

  try {
    const parsed = JSON.parse(record.connection_config);
    if (parsed && typeof parsed === 'object') {
      config = parsed as Record<string, unknown>;
    }
  } catch {
    config = {};
  }

  return {
    id: record.id,
    name: record.name,
    description: record.description ?? '',
    protocol: record.protocol as ProtocolId,
    config,
    status: record.status === 'active' ? 'active' : 'draft',
    test: null,
    persisted: true,
    save_state: 'saved',
    save_error: null,
    runtime_apply_status: null,
    runtime_apply_message: null,
    availability_status: record.availability_status ?? 'available',
    availability_reason: record.availability_reason ?? null,
    running: record.running ?? false,
  };
}
