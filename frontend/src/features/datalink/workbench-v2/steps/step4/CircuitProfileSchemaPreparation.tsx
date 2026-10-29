import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import { groupProfileKind } from '../../state/circuitProfileBasic';
import { aggregateDestinationBlockers, type DestinationBlockerCode } from '../../state/circuitProfileBlockers';
import { CircuitProfileBlockers } from './CircuitProfileBlockers';
import { useWriteGroupReadinessQuery } from '../../../../../hooks/datalink/useStudioV2WriteGroups';
import { GroupSchemaPanel } from './writeGroup/GroupSchemaPanel';
import { GroupTestWritePanel } from './writeGroup/GroupTestWritePanel';

export interface CircuitProfileSchemaPreparationProps {
  workspaceId?: string;
  workspaceRevision: string;
  preparedGroups: WriteGroup[];
  dirtyGroupIds?: ReadonlySet<string>;
  disabled: boolean;
  readonly: boolean;
  isPreparing: boolean;
  destinationMatches: boolean;
  onPrepare: () => Promise<void> | void;
  onApplied: () => void;
  /** Where the operator fixes a destination cause: the destination settings, or one exact saved group in Advanced. */
  onRepairDestination?: () => void;
  onOpenAdvanced?: (groupId: string) => void;
  /** The shared test-write surface is shown for a saved MySQL group; the backend flag alone opens the action. */
  isMySQL?: boolean;
  mysqlTestWriteAvailable?: boolean;
  isSQLServer?: boolean;
  sqlserverTestWriteAvailable?: boolean;
}

interface PreparedGroupProps {
  group: WriteGroup;
  workspaceId: string;
  workspaceRevision: string;
  dirty: boolean;
  readonly: boolean;
  destinationMatches: boolean;
  isMySQL: boolean;
  mysqlTestWriteAvailable: boolean;
  isSQLServer?: boolean;
  sqlserverTestWriteAvailable?: boolean;
  onApplied: () => void;
  onBlocker: (groupId: string, code: DestinationBlockerCode | undefined) => void;
}

interface PreparedGroupTestWriteProps {
  group: WriteGroup;
  dirty: boolean;
  readonly: boolean;
  destinationMatches: boolean;
  testWriteAvailable: boolean;
}

const PreparedGroupTestWrite: React.FC<PreparedGroupTestWriteProps> = ({
  group, dirty, readonly, destinationMatches, testWriteAvailable,
}) => {
  const { t } = useTranslation('workbench-v2');
  const readiness = useWriteGroupReadinessQuery(group.id, !dirty);
  const readinessIsCurrent = readiness.data?.group_revision === group.revision && !readiness.isFetching && !readiness.isError;
  const schemaReady = Boolean(readinessIsCurrent && readiness.data?.schema_ready);
  const canTestWrite = !readonly && !dirty && destinationMatches && group.status !== 'deleted' &&
    testWriteAvailable && schemaReady;
  const blockedReason = readonly ? t('step4.group.test_write.blocked_readonly')
    : dirty ? t('step4.group.test_write.blocked_dirty')
      : !destinationMatches ? t('step4.group.test_write.blocked_destination')
        : readiness.isError ? t('step4.group.test_write.blocked_readiness_unknown')
          : !readinessIsCurrent || readiness.isLoading || readiness.isFetching ? t('step4.group.test_write.blocked_readiness_loading')
            : !schemaReady ? t('step4.group.test_write.blocked_schema_not_ready')
              : !testWriteAvailable ? t('step4.group.test_write.blocked_capability')
                : group.status === 'deleted' ? t('step4.group.test_write.blocked_deleted') : undefined;

  return <GroupTestWritePanel
    workspaceId={group.workspace_id}
    groupId={group.id}
    // Capabilities gate new writes; querying an existing operation remains valid
    // while they refresh or an unrelated setup save moves the workspace revision.
    scopeKey={`${group.id}|${group.revision}`}
    enabled={canTestWrite}
    disabledReason={blockedReason}
  />;
};

const PreparedGroup: React.FC<PreparedGroupProps> = ({
  group, workspaceId, workspaceRevision, dirty, readonly, destinationMatches, isMySQL,
  mysqlTestWriteAvailable, isSQLServer, sqlserverTestWriteAvailable, onApplied, onBlocker,
}) => {
  const { t } = useTranslation('workbench-v2');
  const [readinessRefreshKey, setReadinessRefreshKey] = React.useState(0);
  const handleApplied = () => {
    setReadinessRefreshKey((current) => current + 1);
    onApplied();
  };

  return (
    <div data-testid={`circuit-prepared-group-${groupProfileKind(group) ?? 'instantaneous'}`} className="space-y-2 rounded border border-slate-800/80 bg-slate-950/30 p-3">
      <p className="text-xs font-medium text-slate-300">
        {group.name} ({t(`step4.basic.circuit.${groupProfileKind(group) ?? 'instantaneous'}_label`)})
      </p>
      <GroupSchemaPanel
        key={`${group.id}:${group.revision}:${workspaceRevision}`}
        autoPreview
        group={group}
        workspaceId={workspaceId}
        workspaceRevision={workspaceRevision}
        dirty={dirty}
        destinationMatches={destinationMatches}
        readonly={readonly}
        onApplied={handleApplied}
        onBlocker={onBlocker}
      />
      {isMySQL && <PreparedGroupTestWrite
        key={`mysql-${readinessRefreshKey}`}
        group={group}
        dirty={dirty}
        readonly={readonly}
        destinationMatches={destinationMatches}
        testWriteAvailable={mysqlTestWriteAvailable}
      />}
      {isSQLServer && <PreparedGroupTestWrite
        key={`sqlserver-${readinessRefreshKey}`}
        group={group}
        dirty={dirty}
        readonly={readonly}
        destinationMatches={destinationMatches}
        testWriteAvailable={Boolean(sqlserverTestWriteAvailable)}
      />}
    </div>
  );
};

export const CircuitProfileSchemaPreparation: React.FC<CircuitProfileSchemaPreparationProps> = ({
  workspaceId,
  workspaceRevision,
  preparedGroups,
  dirtyGroupIds,
  disabled,
  readonly,
  isPreparing,
  destinationMatches,
  onPrepare,
  onApplied,
  onRepairDestination,
  onOpenAdvanced,
  isMySQL = false,
  mysqlTestWriteAvailable = false,
  isSQLServer = false,
  sqlserverTestWriteAvailable = false,
}) => {
  const { t } = useTranslation('workbench-v2');
  const [reports, setReports] = React.useState<Record<string, DestinationBlockerCode | undefined>>({});
  const report = React.useCallback((groupId: string, code: DestinationBlockerCode | undefined) => {
    setReports((current) => (current[groupId] === code ? current : { ...current, [groupId]: code }));
  }, []);

  return (
    <section
      data-testid="circuit-profile-schema-section"
      aria-label={t('step4.basic.circuit.step_schema')}
      className="space-y-3"
    >
      <div className="flex flex-wrap items-center justify-between gap-2">
        <h4 className="text-xs font-semibold text-slate-200">{t('step4.basic.circuit.step_schema')}</h4>
        <button
          type="button"
          data-testid="circuit-profile-prepare"
          disabled={disabled || readonly || isPreparing}
          onClick={() => void onPrepare()}
          className="rounded border border-slate-600 px-3 py-2 text-xs text-slate-200 disabled:opacity-40"
        >
          {t('step4.basic.prepare_schema')}
        </button>
        <button
          type="button"
          data-testid="basic-recording-prepare"
          disabled={disabled || readonly || isPreparing}
          className="hidden"
          aria-hidden="true"
          tabIndex={-1}
        />
      </div>
      <CircuitProfileBlockers blockers={aggregateDestinationBlockers(preparedGroups, reports)} onRepairDestination={onRepairDestination} onOpenAdvanced={onOpenAdvanced} />
      {preparedGroups.map((group) => (
          <PreparedGroup
            key={`${group.id}:${group.revision}`}
            group={group}
            workspaceId={workspaceId ?? group.workspace_id}
            workspaceRevision={workspaceRevision}
            dirty={dirtyGroupIds?.has(group.id) ?? false}
            readonly={readonly}
            destinationMatches={destinationMatches}
            isMySQL={isMySQL}
            mysqlTestWriteAvailable={mysqlTestWriteAvailable}
            isSQLServer={isSQLServer}
            sqlserverTestWriteAvailable={sqlserverTestWriteAvailable}
            onApplied={onApplied}
            onBlocker={report}
          />
      ))}
    </section>
  );
};
