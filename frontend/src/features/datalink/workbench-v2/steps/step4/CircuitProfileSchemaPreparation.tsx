import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import { groupProfileKind } from '../../state/circuitProfileBasic';
import { aggregateDestinationBlockers, type DestinationBlockerCode } from '../../state/circuitProfileBlockers';
import { CircuitProfileBlockers } from './CircuitProfileBlockers';
import { GroupSchemaPanel } from './writeGroup/GroupSchemaPanel';

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
}

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
      {preparedGroups.map((group) => {
        const kind = groupProfileKind(group) ?? 'instantaneous';
        return (
          <div
            key={`${group.id}:${group.revision}:${workspaceRevision}`}
            data-testid={`circuit-prepared-group-${kind}`}
            className="space-y-2 rounded border border-slate-800/80 bg-slate-950/30 p-3"
          >
            <p className="text-xs font-medium text-slate-300">
              {group.name} ({t(`step4.basic.circuit.${kind}_label`)})
            </p>
            <GroupSchemaPanel
              autoPreview
              group={group}
              workspaceId={workspaceId ?? group.workspace_id}
              workspaceRevision={workspaceRevision}
              dirty={dirtyGroupIds?.has(group.id) ?? false}
              destinationMatches={destinationMatches}
              readonly={readonly}
              onApplied={onApplied}
              onBlocker={report}
            />
          </div>
        );
      })}
    </section>
  );
};
