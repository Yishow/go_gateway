import type { MouseEvent } from 'react';
import { useTranslation } from 'react-i18next';
import type {
  StudioV2RuntimeContextDevice,
  StudioV2RuntimeSetupContext,
  StudioV2RuntimeSetupMapping,
  StudioV2RuntimeSetupSourceRule,
} from '../../../types/studioV2RuntimeContext';
import type { StudioV2WorkspaceReadinessIssue } from '../../../types/studioV2WorkspaceReadiness';

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

  if (!setupContext) {
    return null;
  }

  const readiness = setupContext.readiness_summary;
  const visibleIssues = readiness?.issues.slice(0, 4) ?? [];
  const primaryIssue = readiness?.issues.find((issue) => issue.severity === 'blocking')
    ?? readiness?.issues[0]
    ?? null;
  const studioReturnTarget = buildStudioV2FocusTarget(primaryIssue);
  const rules = setupContext.source_rules.slice(0, 4);
  const mappings = setupContext.mappings.slice(0, 6);
  const databaseConfig = setupContext.database_config;
  const databaseTargets = setupContext.database_targets;
  const projectionIsStale =
    selectedDevice?.projection_alignment && selectedDevice.projection_alignment !== 'aligned';

  return (
    <section
      className="rounded-3xl border border-slate-800 bg-slate-900/70 p-6"
      data-testid="runtime-dashboard-setup-context"
    >
      <div className="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
        <div>
          <h2 className="text-lg font-semibold text-slate-50">
            {t('setup.title', 'Workspace readiness')}
          </h2>
          <p className="mt-1 text-sm text-slate-400">
            {t(
              'setup.description',
              'Saved Studio V2 setup conditions used by this runtime view.',
            )}
          </p>
        </div>
        <a
          href={studioReturnTarget}
          data-testid="runtime-dashboard-setup-return-link"
          onClick={handleStudioTargetClick(navigateTo, studioReturnTarget)}
          className="inline-flex w-fit rounded-xl border border-slate-700 px-3 py-2 text-sm font-medium text-slate-100 transition hover:border-cyan-400 hover:text-cyan-100"
        >
          {t('setup.returnToStudio', 'Return to Studio V2')}
        </a>
      </div>

      <div className="mt-5 grid gap-4 xl:grid-cols-[0.95fr_1.05fr]">
        <div className="space-y-4">
          <div className="rounded-2xl border border-slate-800 bg-slate-950/60 p-4">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <h3 className="text-sm font-semibold text-slate-100">
                {t('setup.readiness.title', 'Readiness result')}
              </h3>
              <span className={readiness?.ready ? 'text-sm text-emerald-300' : 'text-sm text-amber-300'}>
                {formatReadinessStatus(readiness?.ready, readiness?.blocking_count ?? 0, t)}
              </span>
            </div>
            {visibleIssues.length > 0 ? (
              <ul className="mt-3 space-y-3">
                {visibleIssues.map((issue) => (
                  <li key={`${issue.code}-${issue.scope}`} className="text-sm">
                    <div className="flex flex-wrap gap-2 text-slate-100">
                      <span className="font-semibold">{readinessStepLabel(issue)}</span>
                      <span className="font-mono text-xs text-amber-300">{issue.code}</span>
                    </div>
                    <p className="mt-1 text-slate-300">{issue.message}</p>
                    <div className="mt-2 flex flex-wrap items-center justify-between gap-2">
                      <p className="text-xs text-cyan-200">{readinessActionFor(issue, t)}</p>
                      <a
                        href={buildStudioV2FocusTarget(issue)}
                        data-testid={`runtime-dashboard-setup-fix-${issue.code}`}
                        onClick={handleStudioTargetClick(navigateTo, buildStudioV2FocusTarget(issue))}
                        className="rounded-lg border border-cyan-500/40 bg-cyan-500/10 px-2 py-1 text-xs font-semibold text-cyan-100 transition hover:bg-cyan-500/20"
                      >
                        {fixInStudioLabel(t, readinessStepLabel(issue))}
                      </a>
                    </div>
                  </li>
                ))}
              </ul>
            ) : (
              <p className="mt-3 text-sm text-slate-400">
                {t('setup.readiness.empty', 'No readiness blockers reported.')}
              </p>
            )}
          </div>

          {projectionIsStale ? (
            <div className="rounded-2xl border border-amber-500/30 bg-amber-500/10 p-4">
              <h3 className="text-sm font-semibold text-amber-100">
                {t('setup.projectionStale', 'Projection stale')}
              </h3>
              <dl className="mt-3 grid gap-2 text-sm text-amber-50/90 sm:grid-cols-2">
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

function RuntimeSetupRules({
  rules,
  total,
}: {
  rules: StudioV2RuntimeSetupSourceRule[];
  total: number;
}) {
  const { t } = useTranslation('runtime-dashboard');

  return (
    <div className="rounded-2xl border border-slate-800 bg-slate-950/60 p-4">
      <h3 className="text-sm font-semibold text-slate-100">
        {t('setup.rules.title', 'Configured source rules')}
      </h3>
      {rules.length > 0 ? (
        <ul className="mt-3 space-y-2">
          {rules.map((rule) => (
            <li key={rule.id} className="text-sm text-slate-300">
              <span className="font-mono text-cyan-200">{rule.start_address}</span>
              <span className="mx-2 text-slate-600">/</span>
              <span>{`${rule.count} ${t('setup.rules.pointUnit', 'points')}`}</span>
              <span className="mx-2 text-slate-600">/</span>
              <span>{rule.target_data_type ?? rule.data_type}</span>
            </li>
          ))}
        </ul>
      ) : (
        <p className="mt-3 text-sm text-slate-400">
          {t('setup.rules.empty', 'No persisted source rules.')}
        </p>
      )}
      {total > rules.length ? (
        <p className="mt-2 text-xs text-slate-500">
          {`+${total - rules.length} ${t('setup.moreSuffix', 'more')}`}
        </p>
      ) : null}
    </div>
  );
}

function RuntimeSetupMappings({
  mappings,
  total,
}: {
  mappings: StudioV2RuntimeSetupMapping[];
  total: number;
}) {
  const { t } = useTranslation('runtime-dashboard');

  return (
    <div className="rounded-2xl border border-slate-800 bg-slate-950/60 p-4">
      <h3 className="text-sm font-semibold text-slate-100">
        {t('setup.mappings.title', 'Configured mappings')}
      </h3>
      {mappings.length > 0 ? (
        <ul className="mt-3 grid gap-2 sm:grid-cols-2">
          {mappings.map((mapping) => {
            const title = mapping.display_name || mapping.tag_key || t('setup.mappings.unbound', 'unbound tag');
            const subtitle = mapping.tag_key || (mapping.tag_id ? compactIdentifier(mapping.tag_id) : '');
            return (
              <li key={`${mapping.rule_id}-${mapping.point_id}`} className="text-sm text-slate-300">
                <div className="font-mono text-xs text-cyan-200">{mapping.address}</div>
                <div className="font-semibold text-slate-100">{title}</div>
                {subtitle ? (
                  <div className="text-xs text-slate-400">
                    {subtitle}
                    {mapping.unit ? ` / ${mapping.unit}` : ''}
                  </div>
                ) : null}
              </li>
            );
          })}
        </ul>
      ) : (
        <p className="mt-3 text-sm text-slate-400">
          {t('setup.mappings.empty', 'No persisted mappings.')}
        </p>
      )}
      {total > mappings.length ? (
        <p className="mt-2 text-xs text-slate-500">
          {`+${total - mappings.length} ${t('setup.moreSuffix', 'more')}`}
        </p>
      ) : null}
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
    <div className="rounded-2xl border border-slate-800 bg-slate-950/60 p-4">
      <h3 className="text-sm font-semibold text-slate-100">
        {t('setup.database.title', 'Database setup')}
      </h3>
      {databaseConfig ? (
        <dl className="mt-3 grid gap-2 text-sm text-slate-300 sm:grid-cols-2">
          <RuntimeSetupFact label={t('setup.database.connector', 'Connector')} value={databaseConfig.name} />
          <RuntimeSetupFact label={t('setup.database.database', 'Database')} value={databaseConfig.database} />
          <RuntimeSetupFact label={t('setup.database.table', 'Table')} value={databaseConfig.table} />
          <RuntimeSetupFact
            label={t('setup.database.interval', 'Interval')}
            value={`${databaseConfig.write_interval_seconds}s`}
          />
          <RuntimeSetupFact label={t('setup.database.status', 'Status')} value={databaseConfig.status ?? '-'} />
          <RuntimeSetupFact
            label={t('setup.database.targets', 'Targets')}
            value={databaseTargetsCount > 0 ? String(databaseTargetsCount) : t('setup.database.noTargets', 'No persisted database targets')}
          />
        </dl>
      ) : (
        <p className="mt-3 text-sm text-slate-400">
          {t('setup.database.empty', 'No persisted database config.')}
        </p>
      )}
    </div>
  );
}

function RuntimeSetupFact({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <dt className="text-xs uppercase tracking-[0.18em] text-slate-500">{label}</dt>
      <dd className="mt-1 break-words font-medium text-slate-100">{value}</dd>
    </div>
  );
}

function compactIdentifier(value: string): string {
  return value.length > 18 ? `${value.slice(0, 8)}...${value.slice(-4)}` : value;
}

function buildStudioV2FocusTarget(issue: StudioV2WorkspaceReadinessIssue | null): string {
  if (!issue) {
    return '/studio/v2';
  }

  const params = new URLSearchParams({
    step: String(readinessIssueOwnerStepNumber(issue)),
    focus: 'readiness',
    issue: issue.code,
  });

  return `/studio/v2?${params.toString()}`;
}

function readinessIssueOwnerStepNumber(issue: StudioV2WorkspaceReadinessIssue): 1 | 2 | 3 | 4 {
  switch (issue.code) {
    case 'tag-missing':
    case 'mapping-missing':
      return 3;
    case 'database-target-missing':
      return 4;
    case 'device-connect-required':
    case 'device-probe-required':
      return 1;
    default:
      return readinessStepNumber(issue.step);
  }
}

function readinessStepLabel(issue: StudioV2WorkspaceReadinessIssue): string {
  return `Step ${readinessIssueOwnerStepNumber(issue)}`;
}

function fixInStudioLabel(
  t: (key: string, options?: Record<string, unknown>) => string,
  stepLabel: string,
): string {
  const fallback = `Fix in ${stepLabel}`;
  const translated = t('setup.fixInStudio', {
    step: stepLabel,
    defaultValue: fallback,
  });
  return translated === 'setup.fixInStudio' ? fallback : translated;
}

function handleStudioTargetClick(
  navigateTo: ((target: string) => void) | undefined,
  target: string,
): ((event: MouseEvent<HTMLAnchorElement>) => void) | undefined {
  if (!navigateTo) {
    return undefined;
  }

  return (event) => {
    event.preventDefault();
    navigateTo(target);
  };
}

function readinessStepNumber(step: StudioV2WorkspaceReadinessIssue['step']): 1 | 2 | 3 | 4 {
  switch (step) {
    case 'Step 2':
      return 2;
    case 'Step 3':
      return 3;
    case 'Step 4':
      return 4;
    case 'Step 1':
    default:
      return 1;
  }
}

function formatReadinessStatus(
  ready: boolean | undefined,
  blockingCount: number,
  t: (key: string, fallback: string) => string,
): string {
  if (ready) {
    return t('setup.readiness.ready', 'Ready');
  }
  return `${blockingCount} ${t('setup.readiness.blockerUnit', 'blockers')}`;
}

function readinessActionFor(
  issue: StudioV2WorkspaceReadinessIssue,
  t: (key: string, fallback: string) => string,
): string {
  switch (issue.code) {
    case 'point-missing':
      return t('setup.actions.pointMissing', 'Fix in Step 2: re-save this source rule to rebuild the missing derived point.');
    case 'database-target-missing':
      return t('setup.actions.databaseTargetMissing', 'Fix in Step 4: create a database target for this derived point.');
    case 'tag-missing':
      return t('setup.actions.tagMissing', 'Fix in Step 3: save the tag generated for this derived point.');
    case 'mapping-missing':
      return t('setup.actions.mappingMissing', 'Fix in Step 3: save the point-to-tag mapping.');
    case 'device-connect-required':
    case 'device-probe-required':
      return t('setup.actions.deviceProbeRequired', 'Fix in Step 1: pass connection/probe readiness before activation.');
    default:
      return t('setup.actions.default', 'Open Studio V2 and resolve the owning step shown above.');
  }
}
