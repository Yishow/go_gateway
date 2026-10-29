import type { StudioV2MySQLRecordingCapability } from '../../../../types/studioV2MySQLRecordingCapability';
import type { StudioV2SQLServerRecordingCapability } from '../../../../types/studioV2SQLServerRecordingCapability';
import type { WriteGroupDestinationDraft, WriteGroupStorageMode } from '../../../../types/studioV2WriteGroup';
import type { CircuitPurpose } from './circuitProfileBasic';

export interface RecordingCapabilitySnapshot<T> {
  data?: T;
  isLoading: boolean;
  isFetching: boolean;
  isError: boolean;
}
export interface CircuitProfileCapabilities {
  mysql?: RecordingCapabilitySnapshot<StudioV2MySQLRecordingCapability>;
  sqlserver?: RecordingCapabilitySnapshot<StudioV2SQLServerRecordingCapability>;
}
export interface CircuitProfileCapabilityContext {
  connectorKind?: string;
  destination?: Pick<WriteGroupDestinationDraft, 'connector_id' | 'connector_revision'>;
  capabilities?: CircuitProfileCapabilities;
}
export type CircuitLayoutBlockReason = 'layout_unsupported' | 'capability_unavailable' | 'storage_unsupported';
export type CircuitLayoutOptions = Record<CircuitPurpose, CircuitLayoutBlockReason | undefined>;

/** Layout/storage permission only; never substitutes for schema readiness, ownership or backend mutation CAS. */
export function circuitProfileLayoutBlockReason(
  context: CircuitProfileCapabilityContext, purpose: CircuitPurpose, storageMode: WriteGroupStorageMode,
): CircuitLayoutBlockReason | undefined {
  const { destination, capabilities } = context;
  if (!destination?.connector_id.trim() || !destination.connector_revision.trim()) return 'capability_unavailable';
  const kind = context.connectorKind?.toLowerCase();
  const postgres = kind === 'postgres' || kind === 'postgresql';
  const dynamic = kind === 'mysql' || kind === 'sqlserver';
  if (kind !== 'sqlite' && !postgres && !dynamic) return 'capability_unavailable';
  if (dynamic && purpose === 'wide') return 'layout_unsupported';
  const snapshot = kind === 'mysql' ? capabilities?.mysql : capabilities?.sqlserver;
  const capability = snapshot?.data;
  if (dynamic && (!snapshot || snapshot.isLoading || snapshot.isFetching || snapshot.isError || !capability ||
    capability.connector_id !== destination.connector_id || capability.connector_revision !== destination.connector_revision ||
    !capability.capability_revision.trim() || !capability.connected || capability.metadata_status !== 'verified' ||
    !capability.supports_managed_long || (kind === 'mysql' && !capabilities?.mysql?.data?.transactional_engine))) {
    return 'capability_unavailable';
  }
  if (storageMode === 'plain') return undefined;
  if (purpose === 'long' && storageMode === 'timescale' && postgres) return undefined;
  if (purpose === 'long' && storageMode === 'native_monthly' && dynamic && capability?.supports_native_partition) return undefined;
  return 'storage_unsupported';
}
