import * as React from 'react';
import { useTranslation } from 'react-i18next';
import { BLOCKER_REPAIR, type DestinationBlocker } from '../../state/circuitProfileBlockers';

export interface CircuitProfileBlockersProps {
  blockers: DestinationBlocker[];
  /** Opens the destination settings. */
  onRepairDestination?: () => void;
  /** Opens one exact saved group in Advanced. */
  onOpenAdvanced?: (groupId: string) => void;
}

/** Each cause of a blocked preparation once, with the profiles it affects and where to fix it. */
export const CircuitProfileBlockers: React.FC<CircuitProfileBlockersProps> = ({ blockers, onRepairDestination, onOpenAdvanced }) => {
  const { t } = useTranslation('workbench-v2');
  if (blockers.length === 0) return null;
  return (
    <div className="space-y-2" data-testid="circuit-profile-blockers">
      {blockers.map((blocker, index) => (
        <section key={blocker.key} role="alert" data-testid={`circuit-blocker-${blocker.code}-${index}`}
          className="space-y-1.5 rounded border border-amber-700/60 bg-amber-950/20 p-3 text-xs text-amber-100">
          <p className="font-semibold">{t(`step4.basic.blockers.${blocker.code}.title`)}</p>
          <p>{t(`step4.basic.blockers.${blocker.code}.next`)}</p>
          <p className="text-amber-200/80">{t('step4.basic.blockers.affects')}: {blocker.groups.map((group) => group.name).join(' · ')}</p>
          {BLOCKER_REPAIR[blocker.code] === 'destination' ? (
            <button type="button" data-testid="circuit-blocker-repair-destination" onClick={() => onRepairDestination?.()}
              className="rounded border border-amber-500/60 px-2 py-1 hover:bg-amber-900/40">
              {t('step4.basic.blockers.repair_destination')}
            </button>
          ) : (
            <div className="flex flex-wrap gap-2">
              {blocker.groups.map((group) => (
                <button key={group.id} type="button" data-testid={`circuit-blocker-repair-group-${group.id}`} onClick={() => onOpenAdvanced?.(group.id)}
                  className="rounded border border-amber-500/60 px-2 py-1 hover:bg-amber-900/40">
                  {t('step4.basic.blockers.repair_group', { name: group.name })}
                </button>
              ))}
            </div>
          )}
        </section>
      ))}
    </div>
  );
};
