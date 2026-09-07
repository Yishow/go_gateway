import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { ChevronDown, ChevronUp, ExternalLink, ShieldCheck, AlertTriangle } from 'lucide-react';
import type {
  StudioV2RuntimeContextDevice,
  StudioV2RuntimeSetupContext,
  StudioV2RuntimeSetupMapping,
  StudioV2RuntimeSetupSourceRule,
} from '../../../types/studioV2RuntimeContext';
import {
  buildStudioV2FocusTarget,
  formatReadinessStatus,
  readinessStepLabel,
  readinessActionFor,
  fixInStudioLabel,
  compactIdentifier,
  handleStudioTargetClick,
  RuntimeSetupFact,
} from './RuntimeSetupHelpers';

interface RuntimeSetupContextPanelProps {
  setupContext: StudioV2RuntimeSetupContext | null;
  selectedDevice: StudioV2RuntimeContextDevice | null;
  navigateTo?: (target: string) => void;
}

export function RuntimeSetupContextPanel({
  setupContext,
  selectedDevice,
  navigateTo,
}: RuntimeSetupContextPanelProps) {
  const { t } = useTranslation('runtime-dashboard');
  const [isExpanded, setIsExpanded] = useState(false);

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
  const rules = setupContext.source_rules.slice(0, 4);
  const mappings = setupContext.mappings.slice(0, 6);
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
          <button
            type="button"
            onClick={() => setIsExpanded((prev) => !prev)}
            className="inline-flex items-center gap-1.5 rounded-xl border border-slate-700 bg-slate-800/60 px-3 py-1.5 text-xs font-medium text-slate-200 transition hover:bg-slate-700"
          >
            {isExpanded ? (
              <>
                <span>收合配置詳情</span>
                <ChevronUp className="h-3.5 w-3.5" />
              </>
            ) : (
              <>
                <span>展開配置詳情</span>
                <ChevronDown className="h-3.5 w-3.5" />
              </>
            )}
          </button>
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

      <div className={`mt-5 grid gap-4 xl:grid-cols-[0.95fr_1.05fr] ${isExpanded ? 'block' : 'block'}`}>
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
          <RuntimeSetupRules rules={rules} total={setupContext.source_rules.length} />
          <RuntimeSetupMappings mappings={mappings} total={setupContext.mappings.length} />
          <RuntimeSetupDatabase
            databaseConfig={databaseConfig}
            databaseTargetsCount={databaseTargets.length}
          />
        </div>
      </div>
    </section>
  );
}

function RuntimeSetupRules({ rules, total }: { rules: StudioV2RuntimeSetupSourceRule[]; total: number }) {
  const { t } = useTranslation('runtime-dashboard');
  return (
    <div className="rounded-2xl border border-slate-800/90 bg-slate-950/70 p-4">
      <h3 className="text-xs font-semibold uppercase tracking-wider text-slate-400">
        {t('setup.rules.title', 'Configured source rules')}
      </h3>
      {rules.length > 0 ? (
        <ul className="mt-2.5 space-y-1.5">
          {rules.map((rule) => (
            <li key={rule.id} className="text-xs text-slate-300">
              <span className="font-mono text-cyan-200">{rule.start_address}</span>
              <span className="mx-2 text-slate-600">/</span>
              <span>{`${rule.count} ${t('setup.rules.pointUnit', 'points')}`}</span>
              <span className="mx-2 text-slate-600">/</span>
              <span>{rule.target_data_type ?? rule.data_type}</span>
            </li>
          ))}
        </ul>
      ) : (
        <p className="mt-2 text-xs text-slate-400">{t('setup.rules.empty', 'No persisted source rules.')}</p>
      )}
      {total > rules.length && (
        <p className="mt-2 text-xs text-slate-500">{`+${total - rules.length} ${t('setup.moreSuffix', 'more')}`}</p>
      )}
    </div>
  );
}

function RuntimeSetupMappings({ mappings, total }: { mappings: StudioV2RuntimeSetupMapping[]; total: number }) {
  const { t } = useTranslation('runtime-dashboard');
  return (
    <div className="rounded-2xl border border-slate-800/90 bg-slate-950/70 p-4">
      <h3 className="text-xs font-semibold uppercase tracking-wider text-slate-400">
        {t('setup.mappings.title', 'Configured mappings')}
      </h3>
      {mappings.length > 0 ? (
        <ul className="mt-2.5 grid gap-2 sm:grid-cols-2">
          {mappings.map((mapping) => {
            const title = mapping.display_name || mapping.tag_key || t('setup.mappings.unbound', 'unbound tag');
            const subtitle = mapping.tag_key || (mapping.tag_id ? compactIdentifier(mapping.tag_id) : '');
            return (
              <li key={`${mapping.rule_id}-${mapping.point_id}`} className="rounded-lg bg-slate-900/60 p-2 text-xs text-slate-300">
                <div className="font-mono text-[11px] text-cyan-300">{mapping.address}</div>
                <div className="font-medium text-slate-100">{title}</div>
                {subtitle && <div className="text-[11px] text-slate-400">{subtitle}{mapping.unit ? ` / ${mapping.unit}` : ''}</div>}
              </li>
            );
          })}
        </ul>
      ) : (
        <p className="mt-2 text-xs text-slate-400">{t('setup.mappings.empty', 'No persisted mappings.')}</p>
      )}
      {total > mappings.length && (
        <p className="mt-2 text-xs text-slate-500">{`+${total - mappings.length} ${t('setup.moreSuffix', 'more')}`}</p>
      )}
    </div>
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
          <RuntimeSetupFact label={t('setup.database.connector', 'Connector')} value={databaseConfig.name} />
          <RuntimeSetupFact label={t('setup.database.database', 'Database')} value={databaseConfig.database} />
          <RuntimeSetupFact label={t('setup.database.table', 'Table')} value={databaseConfig.table} />
          <RuntimeSetupFact label={t('setup.database.targets', 'Targets')} value={databaseTargetsCount ? String(databaseTargetsCount) : t('setup.database.noTargets', 'No database targets')} />
        </dl>
      ) : (
        <p className="mt-2 text-xs text-slate-400">{t('setup.database.empty', 'No persisted database configuration.')}</p>
      )}
    </div>
  );
}

