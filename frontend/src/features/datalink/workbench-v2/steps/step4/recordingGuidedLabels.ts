import { useTranslation } from 'react-i18next';
import type { RecordingRuleSummaryLabels } from './RecordingRuleSummary';

const extras = ['unsaved', 'unsavedDescription', 'noScope', 'retryEvidence', 'draftHeading', 'savedDetails', 'effectiveDetails',
  'effectiveCounterUnavailable', 'unavailableSource', 'destinationFallback', 'selectedDevice', 'selectedCircuit', 'noProfiles',
  'tablePending', 'savedReferences', 'deliveryDescription'] as const;
const detailFields = ['point', 'tag', 'column', 'sourceRevision', 'mappingRevision', 'required', 'maxAge', 'valueRole', 'resultType',
  'measurement', 'originalType', 'definitionRevision', 'sourceBindingRevision', 'epoch', 'scale', 'physicalUnit', 'semanticKind', 'allowNegative',
  'resetThreshold', 'maxRollover', 'rolloverModulus', 'maxRate', 'schemaRevision', 'schemaDigest', 'writeMode', 'recreateMissing',
  'operation', 'operationStatus', 'readiness', 'schemaEvidence', 'delivery', 'intake', 'collecting', 'queued', 'retrying', 'blocked',
  'quarantined', 'unknown', 'sqlCommitted', 'skipped', 'lastCommit', 'readback', 'cleanup', 'backlog', 'connector', 'connectorRevision',
  'entityKeyColumn', 'groupKeyColumns', 'uniqueKeyColumns', 'valueColumn', 'qualityColumn', 'provenanceColumn', 'recordKeyColumn',
  'bucketStartColumn', 'groupIdColumn', 'deviceIdColumn', 'dedupeCapability', 'requestedAt', 'writeOutcome', 'cleanupStatus', 'writtenAt',
  'true', 'false', 'none', 'append', 'latest', 'gauge', 'counter_raw', 'state', 'counter', 'signed_counter'] as const;
const states = ['pending', 'running', 'succeeded', 'partial', 'failed', 'unknown', 'active', 'retiring', 'blocked', 'not_running',
  'unavailable', 'written_verified', 'written_unverified', 'not_attempted', 'cleaned'] as const;
export type RecordingGuidedLabels = RecordingRuleSummaryLabels & Readonly<Record<typeof extras[number], string>> & {
  readonly states: Readonly<Record<typeof states[number], string>>;
  readonly detailFields: Readonly<Record<typeof detailFields[number], string>>;
};

/** Explicit key allowlists keep translations typed; no server reason is used as a locale key. */
export function recordingGuidedLabels(t: (key: string) => string): RecordingGuidedLabels {
  function record<Key extends string>(prefix: string, keys: readonly Key[]): Readonly<Record<Key, string>> {
    return Object.fromEntries(keys.map(key => [key, t(prefix ? `${prefix}.${key}` : key)])) as Record<Key, string>;
  }
  return {
    ...record('', ['heading', 'saved', 'effective', 'savedDescription', 'effectiveDescription', 'unavailable', 'emptySources',
      'normalizedInterval', 'groupScope', 'requestedNotEffective', 'details', 'detailsTitle', 'detailsDescription', 'detailsUnavailable', 'close']),
    ...record('', extras),
    fields: record('fields', ['sources', 'layout', 'storageMode', 'profile', 'reduction', 'interval', 'lateness', 'incomplete',
      'destination', 'ownership', 'workspace', 'device', 'circuit', 'group', 'savedRevision', 'requestedRevision', 'effectiveRevision', 'asOf', 'effectiveAt', 'endsAt']),
    profiles: record('profiles', ['instantaneous', 'cumulative']), units: record('units', ['second', 'minute']),
    layouts: record('layouts', ['managed_wide_v1', 'managed_long_v1', 'managed_wide_profile_v1', 'unavailable']),
    storageModes: record('storageModes', ['plain', 'timescale', 'native_monthly', 'unavailable']),
    reductions: record('reductions', ['latest', 'mean', 'cumulative_delta', 'member_calculation', 'unavailable']),
    incompletePolicies: record('incompletePolicies', ['skip_row', 'partial', 'unavailable']),
    sourceTypes: record('sourceTypes', ['bool', 'int16', 'uint16', 'int32', 'uint32', 'int64', 'uint64', 'float32', 'float64', 'string']),
    reasons: record('reasons', ['missing', 'loading', 'error', 'unknown', 'scope_mismatch', 'revision_mismatch',
      'effective_snapshot_unavailable', 'operation_unresolved', 'stale', 'invalid_evidence']),
    detailFields: record('detailFields', detailFields), states: record('states', states),
  };
}

export function useRecordingGuidedLabels(): RecordingGuidedLabels {
  const { t } = useTranslation('workbench-v2-recording-guided');
  return recordingGuidedLabels((key) => t(key));
}
