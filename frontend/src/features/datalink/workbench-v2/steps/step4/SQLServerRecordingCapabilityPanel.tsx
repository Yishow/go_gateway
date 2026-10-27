import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { StudioV2SQLServerRecordingCapability } from '../../../../../types/studioV2SQLServerRecordingCapability';

export interface SQLServerRecordingCapabilityPanelProps {
  capability?: StudioV2SQLServerRecordingCapability;
  isLoading: boolean;
  isError: boolean;
  onRetry: () => void;
}

export const SQLServerRecordingCapabilityPanel: React.FC<SQLServerRecordingCapabilityPanelProps> = ({
  capability, isLoading, isError, onRetry,
}) => {
  const { t } = useTranslation('workbench-v2');
  const serverVerified = Boolean(capability?.connected && capability.metadata_status === 'verified');
  const managedLongAvailable = Boolean(serverVerified && capability?.supports_managed_long);
  const nativeMonthlyAvailable = Boolean(managedLongAvailable && capability?.supports_native_partition);
  const operationStatus = (available: boolean) => t(`step4.basic.circuit.sqlserver_capability.${available ? 'available' : 'unavailable'}`);
  const metadataStatus = capability?.metadata_status
    ? t(`step4.basic.circuit.sqlserver_capability.metadata_status.${capability.metadata_status}`, {
        defaultValue: t('step4.basic.circuit.sqlserver_capability.metadata_status.unknown'),
      })
    : t('step4.basic.circuit.sqlserver_capability.metadata_status.unknown');

  return (
    <section
      data-testid="circuit-profile-sqlserver-capability"
      aria-label={t('step4.basic.circuit.sqlserver_capability.title')}
      className="space-y-2 rounded border border-slate-800/80 bg-slate-950/30 p-3 text-xs text-slate-300"
    >
      <h4 className="font-semibold text-slate-200">{t('step4.basic.circuit.sqlserver_capability.title')}</h4>
      <dl className="grid gap-x-4 gap-y-1 sm:grid-cols-2">
        <dt className="text-slate-500">{t('step4.basic.circuit.sqlserver_capability.database')}</dt>
        <dd data-testid="circuit-profile-sqlserver-database" className="text-slate-100">
          {capability?.database || t('step4.basic.circuit.sqlserver_capability.unknown')}
        </dd>
        <dt className="text-slate-500">{t('step4.basic.circuit.sqlserver_capability.schema')}</dt>
        <dd data-testid="circuit-profile-sqlserver-schema" className="text-slate-100">
          {capability?.schema || t('step4.basic.circuit.sqlserver_capability.unknown')}
        </dd>
        <dt className="text-slate-500">{t('step4.basic.circuit.sqlserver_capability.server_version')}</dt>
        <dd data-testid="circuit-profile-sqlserver-server-version" className="text-slate-100">
          {capability?.server_version || t('step4.basic.circuit.sqlserver_capability.unknown')}
        </dd>
        <dt className="text-slate-500">{t('step4.basic.circuit.sqlserver_capability.server_edition')}</dt>
        <dd data-testid="circuit-profile-sqlserver-server-edition" className="text-slate-100">
          {capability?.server_edition || t('step4.basic.circuit.sqlserver_capability.unknown')}
        </dd>
        <dt className="text-slate-500">{t('step4.basic.circuit.sqlserver_capability.compatibility_level')}</dt>
        <dd data-testid="circuit-profile-sqlserver-compatibility-level" className="text-slate-100">
          {capability?.compatibility_level ?? t('step4.basic.circuit.sqlserver_capability.unknown')}
        </dd>
        <dt className="text-slate-500">{t('step4.basic.circuit.sqlserver_capability.connection')}</dt>
        <dd data-testid="circuit-profile-sqlserver-connection" className="text-slate-100">
          {operationStatus(Boolean(capability?.connected))}
        </dd>
        <dt className="text-slate-500">{t('step4.basic.circuit.sqlserver_capability.metadata')}</dt>
        <dd data-testid="circuit-profile-sqlserver-metadata" className="text-slate-100">{metadataStatus}</dd>
        <dt className="text-slate-500">{t('step4.basic.circuit.sqlserver_capability.schema_view_definition')}</dt>
        <dd data-testid="circuit-profile-sqlserver-schema-view-definition" className="text-slate-100">
          {operationStatus(Boolean(serverVerified && capability?.permissions.schema_view_definition))}
        </dd>
      </dl>
      <ul className="grid gap-x-4 gap-y-1 sm:grid-cols-2" aria-label={t('step4.basic.circuit.sqlserver_capability.operations')}>
        <li data-testid="circuit-profile-sqlserver-capability-managed-long">
          {t('step4.basic.circuit.sqlserver_capability.managed_long')}: {operationStatus(managedLongAvailable)}
        </li>
        <li data-testid="circuit-profile-sqlserver-capability-native-partitions">
          {t('step4.basic.circuit.sqlserver_capability.native_partitions')}: {operationStatus(nativeMonthlyAvailable)}
        </li>
        <li data-testid="circuit-profile-sqlserver-capability-test-write">
          {t('step4.basic.circuit.sqlserver_capability.test_write')}: {operationStatus(Boolean(serverVerified && capability?.supports_group_test_write))}
        </li>
        <li data-testid="circuit-profile-sqlserver-capability-missing-repair">
          {t('step4.basic.circuit.sqlserver_capability.missing_repair')}: {operationStatus(Boolean(serverVerified && capability?.supports_missing_repair))}
        </li>
      </ul>
      {capability && capability.reasons.length > 0 && (
        <div data-testid="circuit-profile-sqlserver-capability-reasons" className="space-y-1 text-slate-400">
          <p>{t('step4.basic.circuit.sqlserver_capability.reasons_title')}</p>
          <ul className="list-disc space-y-0.5 pl-5">
            {capability.reasons.map((reason) => (
              <li key={reason}>{t(`step4.basic.circuit.sqlserver_capability.reason.${reason}`, {
                defaultValue: t('step4.basic.circuit.sqlserver_capability.reason.generic'),
              })}</li>
            ))}
          </ul>
        </div>
      )}
      {isLoading && <p role="status">{t('step4.basic.circuit.sqlserver_capability.loading')}</p>}
      {isError && (
        <div role="status" data-testid="circuit-profile-sqlserver-capability-error" className="flex flex-wrap items-center gap-2 text-amber-300">
          <span>{t('step4.basic.circuit.sqlserver_capability.load_error')}</span>
          <button type="button" onClick={onRetry} className="rounded border border-amber-800 px-2 py-1">
            {t('step4.basic.circuit.sqlserver_capability.retry')}
          </button>
        </div>
      )}
      {!capability && !isLoading && !isError && (
        <p role="status" className="text-slate-400">{t('step4.basic.circuit.sqlserver_capability.not_available')}</p>
      )}
    </section>
  );
};
