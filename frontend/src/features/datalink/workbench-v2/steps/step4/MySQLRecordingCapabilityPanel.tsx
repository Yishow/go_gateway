import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { StudioV2MySQLRecordingCapability } from '../../../../../types/studioV2MySQLRecordingCapability';

export interface MySQLRecordingCapabilityPanelProps {
  capability?: StudioV2MySQLRecordingCapability;
  isLoading: boolean;
  isError: boolean;
  onRetry: () => void;
}

function serverFamilyLabel(family: string | undefined, t: (key: string) => string): string {
  if (family === 'mysql') return t('step4.basic.circuit.mysql_capability.family.mysql');
  if (family === 'mariadb') return t('step4.basic.circuit.mysql_capability.family.mariadb');
  return t('step4.basic.circuit.mysql_capability.family.unknown');
}

export const MySQLRecordingCapabilityPanel: React.FC<MySQLRecordingCapabilityPanelProps> = ({
  capability, isLoading, isError, onRetry,
}) => {
  const { t } = useTranslation('workbench-v2');
  const serverVerified = Boolean(capability?.connected && capability.metadata_status === 'verified');
  const managedLongAvailable = Boolean(serverVerified && capability?.transactional_engine && capability.supports_managed_long);
  const nativeMonthlyAvailable = Boolean(managedLongAvailable && capability?.supports_native_partition);
  const operationStatus = (available: boolean) => t(`step4.basic.circuit.mysql_capability.${available ? 'available' : 'unavailable'}`);
  const metadataStatus = capability?.metadata_status
    ? t(`step4.basic.circuit.mysql_capability.metadata_status.${capability.metadata_status}`, {
        defaultValue: t('step4.basic.circuit.mysql_capability.metadata_status.unknown'),
      })
    : t('step4.basic.circuit.mysql_capability.metadata_status.unknown');

  return (
    <section
      data-testid="circuit-profile-mysql-capability"
      aria-label={t('step4.basic.circuit.mysql_capability.title')}
      className="space-y-2 rounded border border-slate-800/80 bg-slate-950/30 p-3 text-xs text-slate-300"
    >
      <h4 className="font-semibold text-slate-200">{t('step4.basic.circuit.mysql_capability.title')}</h4>
      <dl className="grid gap-x-4 gap-y-1 sm:grid-cols-2">
        <dt className="text-slate-500">{t('step4.basic.circuit.mysql_capability.server_family')}</dt>
        <dd data-testid="circuit-profile-mysql-server-family" className="text-slate-100">
          {serverFamilyLabel(capability?.server_family, t)}
        </dd>
        <dt className="text-slate-500">{t('step4.basic.circuit.mysql_capability.server_version')}</dt>
        <dd data-testid="circuit-profile-mysql-server-version" className="text-slate-100">
          {capability?.server_version || t('step4.basic.circuit.mysql_capability.unknown')}
        </dd>
        <dt className="text-slate-500">{t('step4.basic.circuit.mysql_capability.connection')}</dt>
        <dd data-testid="circuit-profile-mysql-connection" className="text-slate-100">
          {operationStatus(Boolean(capability?.connected))}
        </dd>
        <dt className="text-slate-500">{t('step4.basic.circuit.mysql_capability.metadata')}</dt>
        <dd data-testid="circuit-profile-mysql-metadata" className="text-slate-100">{metadataStatus}</dd>
        <dt className="text-slate-500">{t('step4.basic.circuit.mysql_capability.transactional_engine')}</dt>
        <dd data-testid="circuit-profile-mysql-transactional-engine" className="text-slate-100">
          {operationStatus(Boolean(serverVerified && capability?.transactional_engine))}
        </dd>
      </dl>
      <ul className="grid gap-x-4 gap-y-1 sm:grid-cols-2" aria-label={t('step4.basic.circuit.mysql_capability.operations')}>
        <li data-testid="circuit-profile-mysql-capability-managed-long">
          {t('step4.basic.circuit.mysql_capability.managed_long')}: {operationStatus(managedLongAvailable)}
        </li>
        <li data-testid="circuit-profile-mysql-capability-native-partitions">
          {t('step4.basic.circuit.mysql_capability.native_partitions')}: {operationStatus(nativeMonthlyAvailable)}
        </li>
        <li data-testid="circuit-profile-mysql-capability-test-write">
          {t('step4.basic.circuit.mysql_capability.test_write')}: {operationStatus(Boolean(serverVerified && capability?.supports_group_test_write))}
        </li>
        <li data-testid="circuit-profile-mysql-capability-missing-repair">
          {t('step4.basic.circuit.mysql_capability.missing_repair')}: {operationStatus(Boolean(serverVerified && capability?.supports_missing_repair))}
        </li>
      </ul>
      {capability && capability.reasons.length > 0 && (
        <div data-testid="circuit-profile-mysql-capability-reasons" className="space-y-1 text-slate-400">
          <p>{t('step4.basic.circuit.mysql_capability.reasons_title')}</p>
          <ul className="list-disc space-y-0.5 pl-5">
            {capability.reasons.map((reason) => (
              <li key={reason}>{t(`step4.basic.circuit.mysql_capability.reason.${reason}`, {
                defaultValue: t('step4.basic.circuit.mysql_capability.reason.generic'),
              })}</li>
            ))}
          </ul>
        </div>
      )}
      {isLoading && <p role="status">{t('step4.basic.circuit.mysql_capability.loading')}</p>}
      {isError && (
        <div role="status" data-testid="circuit-profile-mysql-capability-error" className="flex flex-wrap items-center gap-2 text-amber-300">
          <span>{t('step4.basic.circuit.mysql_capability.load_error')}</span>
          <button type="button" onClick={onRetry} className="rounded border border-amber-800 px-2 py-1">
            {t('step4.basic.circuit.mysql_capability.retry')}
          </button>
        </div>
      )}
      {!capability && !isLoading && !isError && (
        <p role="status" className="text-slate-400">{t('step4.basic.circuit.mysql_capability.not_available')}</p>
      )}
    </section>
  );
};
