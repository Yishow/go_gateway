import { RuntimeSettingsLink } from './RuntimeSettingsLink';
import type { RuntimeSourceSelection } from '../../../types/studioV2RuntimeContext';
import type { EvidenceScope } from '../../../types/runtimeRuleEvidence';
import { safeRuntimeIdentity } from '../../../utils/runtimeOperatorSafety';
import { RuntimeSourceRuleList } from './RuntimeSourceRuleList';
import { useTranslation } from 'react-i18next';
import { ExternalLink, ShieldCheck, AlertTriangle } from 'lucide-react';
import type {
  StudioV2RuntimeContextDevice,
  StudioV2RuntimeSetupContext,
} from '../../../types/studioV2RuntimeContext';
import {
  buildStudioV2FocusTarget,
  formatReadinessStatus,
  readinessStepLabel,
  readinessActionFor,
  fixInStudioLabel,
  handleStudioTargetClick,
  RuntimeSetupFact,
} from './RuntimeSetupHelpers';

interface RuntimeSetupContextPanelProps {
  setupContext: StudioV2RuntimeSetupContext | null;
  selectedDevice: StudioV2RuntimeContextDevice | null;
  navigateTo?: (target: string) => void;
  evidenceScope?: EvidenceScope | null;
  sourceSelection?: RuntimeSourceSelection | null;
  onSourceSelection?: (selection: RuntimeSourceSelection | null) => void;
}

export function RuntimeSetupContextPanel({
  setupContext,
  selectedDevice,
  navigateTo,
  evidenceScope = null,
  sourceSelection = null,
  onSourceSelection,
}: RuntimeSetupContextPanelProps) {
  const { t } = useTranslation('runtime-dashboard');
  const scope = evidenceScope && evidenceScope.device_id === selectedDevice?.device_id &&
    Object.values(evidenceScope).every((id) => Boolean(id) && safeRuntimeIdentity(id) === id) ? evidenceScope : null;
  const scopeKey = scope ? JSON.stringify([scope.workspace_id, scope.workspace_revision, scope.device_id]) : '';

  if (!setupContext) {
    return null;
  }

  const readiness = setupContext.readiness_summary;
  const visibleIssues = readiness?.issues.slice(0, 4) ?? [];
  const primaryIssue =
    readiness?.issues.find((issue) => issue.severity === 'blocking') ??
    readiness?.issues[0] ??
    null;
  const studioReturnTarget = buildStudioV2FocusTarget(primaryIssue);
  const databaseConfig = setupContext.database_config;
  const databaseTargets = setupContext.database_targets;
  const projectionIsStale =
    selectedDevice?.projection_alignment && selectedDevice.projection_alignment !== 'aligned';
  const hasBlockers = (readiness?.blocking_count ?? 0) > 0;

  return (
    <section
      className={`rounded-3xl border transition-all duration-300 ${
        hasBlockers
          ? 'border-amber-500/40 bg-amber-500/5'
          : 'border-slate-800/80 bg-slate-900/60'
      } p-5 shadow-lg`}
      data-testid="runtime-dashboard-setup-context"
    >
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <div
            className={`flex h-10 w-10 items-center justify-center rounded-xl border ${
              hasBlockers
                ? 'border-amber-500/30 bg-amber-500/10 text-amber-300'
                : 'border-cyan-500/30 bg-cyan-500/10 text-cyan-300'
            }`}
          >
            {hasBlockers ? <AlertTriangle className="h-5 w-5" /> : <ShieldCheck className="h-5 w-5" />}
          </div>
          <div>
            <div className="flex items-center gap-2.5">
          {scope && <RuntimeSettingsLink locator={{ step: 1, workspace_id: scope.workspace_id, device_id: scope.device_id, rule_id: null, mapping_id: null, group_id: null, group_revision: null }} navigateTo={navigateTo}>{t('settingsLink.device')}</RuntimeSettingsLink>}
              <h2 className="text-base font-semibold text-slate-100">
                {t('setup.title', 'Workspace readiness')}
              </h2>
              <span
                className={`inline-flex items-center rounded-md px-2 py-0.5 text-xs font-semibold ${
                  readiness?.ready
                    ? 'border border-emerald-500/30 bg-emerald-500/10 text-emerald-300'
                    : 'border border-amber-500/30 bg-amber-500/10 text-amber-300'
                }`}
              >
                {formatReadinessStatus(readiness?.ready, readiness?.blocking_count ?? 0, t)}
              </span>
            </div>
            <p className="mt-0.5 text-xs text-slate-400">
              {t(
                'setup.description',
                'Saved Studio V2 setup conditions used by this runtime view.',
              )}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2.5">
          <a
            href={studioReturnTarget}
            data-testid="runtime-dashboard-setup-return-link"
            onClick={handleStudioTargetClick(navigateTo, studioReturnTarget)}
            className="inline-flex items-center gap-1.5 rounded-xl border border-cyan-500/40 bg-cyan-500/10 px-3 py-1.5 text-xs font-medium text-cyan-200 transition hover:bg-cyan-500/20"
          >
            <span>{t('setup.returnToStudio', 'Return to Studio V2')}</span>
            <ExternalLink className="h-3.5 w-3.5" />
          </a>
        </div>
      </div>

      <div className="mt-5 grid gap-4 xl:grid-cols-[0.95fr_1.05fr]">
        <div className="space-y-4">
          <div className="rounded-2xl border border-slate-800/90 bg-slate-950/70 p-4">
            <div className="flex flex-wrap items-center justify-between gap-2 border-b border-slate-800 pb-2.5">
              <h3 className="text-xs font-semibold uppercase tracking-wider text-slate-400">
                {t('setup.readiness.title', 'Readiness result')}
              </h3>
              <span className={readiness?.ready ? 'text-xs font-medium text-emerald-300' : 'text-xs font-medium text-amber-300'}>
                {formatReadinessStatus(readiness?.ready, readiness?.blocking_count ?? 0, t)}
              </span>
            </div>
            {visibleIssues.length > 0 ? (
              <ul className="mt-3 space-y-3">
                {visibleIssues.map((issue) => (
                  <li key={`${issue.code}-${issue.scope}`} className="rounded-xl border border-amber-500/20 bg-amber-500/5 p-3 text-sm">
                    <div className="flex flex-wrap items-center justify-between gap-2 text-slate-100">
                      <span className="font-semibold text-amber-200">{readinessStepLabel(issue)}</span>
                      <span className="rounded bg-slate-900 px-2 py-0.5 font-mono text-xs text-amber-300">{issue.code}</span>
                    </div>
                    <p className="mt-1.5 text-xs text-slate-300 leading-relaxed">{issue.message}</p>
                    <div className="mt-2.5 flex flex-wrap items-center justify-between gap-2 border-t border-amber-500/10 pt-2">
                      <p className="text-xs text-cyan-200">{readinessActionFor(issue, t)}</p>
                      <a
                        href={buildStudioV2FocusTarget(issue)}
                        data-testid={`runtime-dashboard-setup-fix-${issue.code}`}
                        onClick={handleStudioTargetClick(navigateTo, buildStudioV2FocusTarget(issue))}
                        className="inline-flex items-center gap-1 rounded-lg border border-cyan-500/40 bg-cyan-500/10 px-2.5 py-1 text-xs font-semibold text-cyan-100 transition hover:bg-cyan-500/20"
                      >
                        {fixInStudioLabel(t, readinessStepLabel(issue))}
                      </a>
                    </div>
                  </li>
                ))}
              </ul>
            ) : (
              <p className="mt-3 text-xs text-slate-400">
                {t('setup.readiness.empty', 'No readiness blockers reported.')}
              </p>
            )}
          </div>

          {projectionIsStale ? (
            <div className="rounded-2xl border border-amber-500/30 bg-amber-500/10 p-4">
              <h3 className="text-xs font-semibold uppercase tracking-wider text-amber-100">
                {t('setup.projectionStale', 'Projection stale')}
              </h3>
              <dl className="mt-3 grid gap-2 text-xs text-amber-50/90 sm:grid-cols-2">
                <RuntimeSetupFact
                  label={t('setup.runtimeProjection', 'Runtime projection')}
                  value={selectedDevice?.runtime_projection_version ?? '-'}
                />
                <RuntimeSetupFact
                  label={t('setup.workspaceProjection', 'Workspace projection')}
                  value={selectedDevice?.workspace_projection_version ?? '-'}
                />
              </dl>
              {selectedDevice?.projection_message ? (
                <p className="mt-2 text-xs text-amber-100/80">{selectedDevice.projection_message}</p>
              ) : null}
            </div>
          ) : null}
        </div>

        <div className="space-y-4">
          {scope ? <RuntimeSourceRuleList key={scopeKey} scope={scope} navigateTo={navigateTo} initialSelection={sourceSelection} onSelection={onSourceSelection} /> : <p role="status" className="text-xs text-amber-100">{t('sourceDetail.scopeUnavailable')}</p>}
          <RuntimeSetupDatabase
            databaseConfig={databaseConfig}
            databaseTargetsCount={databaseTargets.length}
          />
        </div>
      </div>
    </section>
  );
}

function RuntimeSetupDatabase({
  databaseConfig,
  databaseTargetsCount,
}: {
  databaseConfig: StudioV2RuntimeSetupContext['database_config'];
  databaseTargetsCount: number;
}) {
  const { t } = useTranslation('runtime-dashboard');
  return (
    <div className="rounded-2xl border border-slate-800/90 bg-slate-950/70 p-4">
      <h3 className="text-xs font-semibold uppercase tracking-wider text-slate-400">
        {t('setup.database.title', 'Database setup')}
      </h3>
      {databaseConfig ? (
        <dl className="mt-2.5 grid gap-2 text-xs text-slate-300 sm:grid-cols-2">
          <RuntimeSetupFact label={t('setup.database.connector', 'Connector')} value={`${databaseConfig.name}${databaseConfig.kind ? ` (${databaseConfig.kind.toUpperCase()})` : ''}`} />
          <RuntimeSetupFact label={t('setup.database.table', 'Table')} value={databaseConfig.table} />
          <RuntimeSetupFact label={t('setup.database.targets', 'Targets')} value={databaseTargetsCount ? String(databaseTargetsCount) : t('setup.database.noTargets', 'No database targets')} />
        </dl>
      ) : (
        <p className="mt-2 text-xs text-slate-400">{t('setup.database.empty', 'No persisted database configuration.')}</p>
      )}
    </div>
  );
}

