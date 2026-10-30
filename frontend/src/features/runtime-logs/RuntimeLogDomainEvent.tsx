import type { MouseEvent } from 'react';
import { Link, useInRouterContext } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import type { RuntimeLogRecord } from '../../types/runtimeLogs';
import { projectRuntimeLogDomain } from '../../services/runtimeLogDomainProjection';
import { buildStudioEntityLocatorTarget } from '../../services/studioEntityLocator';
import { handleStudioTargetClick } from '../datalink/runtime-dashboard/RuntimeSetupHelpers';

const focus = 'rounded focus-visible:outline focus-visible:outline-2 focus-visible:outline-blue-600';
type Fact = [label: string, value: string | null];
function Facts({ items, unavailable }: { items: Fact[]; unavailable: string }) {
  return <dl className="grid gap-x-6 gap-y-2 sm:grid-cols-2">
    {items.map(([label, value]) => <div key={label} className="min-w-0">
      <dt className="text-slate-600">{label}</dt>
      <dd className="break-words font-mono text-slate-900">{value ?? unavailable}</dd>
    </div>)}
  </dl>;
}

/** Historical safe presentation only. Links navigate; they never execute an operation. */
export default function RuntimeLogDomainEvent({ record, navigateTo }: {
  record: RuntimeLogRecord; navigateTo?: (target: string) => void;
}) {
  const { t } = useTranslation('runtime-logs');
  const inRouter = useInRouterContext();
  const view = projectRuntimeLogDomain(record);
  if (view.kind === 'legacy') return null;
  const unavailable = t('domain.unavailable');
  const summary = view.summary;
  const summaryFacts: Fact[] = [[t('domain.firstSeen'), summary?.firstSeenAt ?? null],
    [t('domain.lastSeen'), summary?.lastSeenAt ?? null]];
  if (view.kind === 'unavailable') return <tr role="row" className="block align-top font-sans sm:table-row">
    <td role="cell" className="block whitespace-nowrap p-3 text-slate-600 sm:table-cell">{view.timestamp ? <time dateTime={view.timestamp}>{view.timestamp}</time> : unavailable}</td>
    <td role="cell" className="block p-3 sm:table-cell">{unavailable}</td><td role="cell" className="block p-3 sm:table-cell">{t('domain.source')}</td>
    <td role="cell" className="block min-w-0 max-w-3xl space-y-3 break-words p-3 sm:table-cell">
      <h2 className="font-semibold text-slate-900">{t('domain.unavailableProblem')}</h2>
      <p className="text-slate-700">{t('domain.unavailableImpact')}</p>
      <p className="text-slate-600">{t('domain.historical')}</p>
      <Facts items={[[t('domain.occurrences'), summary?.occurrences ?? null]]} unavailable={unavailable} />
      <details><summary className={`${focus} cursor-pointer text-blue-800`}>{t('domain.details')}</summary>
        <div className="mt-3"><Facts items={summaryFacts} unavailable={unavailable} /></div>
      </details>
    </td>
  </tr>;

  const codeKey = view.code.slice('domain.'.length);
  const target = view.action !== 'none' ? buildStudioEntityLocatorTarget(view.locator) : null;
  const linkLabel = t(view.locator.step === 4 ? 'domain.openGroup' : 'domain.openCurrent');
  const click = target ? handleStudioTargetClick(navigateTo, target) : undefined;
  const onClick = click ? (event: MouseEvent<HTMLAnchorElement>) => {
    if (event.defaultPrevented || event.button !== 0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;
    click(event);
  } : undefined;
  const identityFacts: Fact[] = [[t('domain.workspaceId'), view.binding.workspace_id]];
  const optionalIDs = [
    ['deviceId', view.locator.device_id], ['ruleId', view.locator.rule_id], ['mappingId', view.locator.mapping_id],
    ['groupId', view.locator.group_id], ['pointId', view.pointID], ['tagId', view.tagID], ['connectorId', view.connectorID],
  ] as const;
  for (const [label, value] of optionalIDs) if (value !== null) identityFacts.push([t(`domain.${label}`), value]);
  identityFacts.push([t('domain.entityRevision'), view.binding.entity_revision],
    [t('domain.workspaceRevision'), view.binding.workspace_revision]);
  if (view.locator.group_revision !== null) identityFacts.push([t('domain.groupRevision'), view.locator.group_revision]);
  if (view.connectorRevision !== null) identityFacts.push([t('domain.connectorRevision'), view.connectorRevision]);
  const detailFacts: Fact[] = [...identityFacts, ...summaryFacts, [t('domain.referenceId'), view.referenceID],
    [t('domain.eventCode'), view.code], [t('domain.evidenceKind'), t(`domain.evidenceKinds.${view.evidenceKind}`)]];
  const affected = view.locator.group_id ?? view.locator.mapping_id ?? view.locator.rule_id ?? view.pointID ?? view.locator.device_id;
  const receipt = view.receipt;
  const quota = view.quota;
  const linkClass = `${focus} inline-flex border border-blue-300 bg-blue-50 px-3 py-2 text-blue-900`;
  return <tr role="row" className="block align-top font-sans sm:table-row">
    <td role="cell" className="block whitespace-nowrap p-3 text-slate-600 sm:table-cell"><time dateTime={view.timestamp}>{view.timestamp}</time></td>
    <td role="cell" className="block p-3 sm:table-cell">{t(`levels.${view.level}`)}</td><td role="cell" className="block p-3 sm:table-cell">{t('domain.source')}</td>
    <td role="cell" className="block min-w-0 max-w-3xl space-y-3 break-words p-3 sm:table-cell">
      <h2 className="font-semibold text-slate-900">{t(`domain.problems.${codeKey}`)}</h2>
      <p className="text-slate-700">{t(`domain.impacts.${codeKey}`)}</p>
      <p className="text-slate-600">{t('domain.historical')}</p>
      {view.locator.group_id !== null && <p className="text-slate-600">{t('domain.groupScope')}</p>}
      <Facts items={[[t('domain.affected'), affected], [t('domain.stage'), t(`stages.${view.stage}`)],
        [t('domain.occurrences'), summary?.occurrences ?? null], [t('domain.reason'), t(`domain.reasons.${view.reason}`)],
        [t('domain.action'), t(`domain.actions.${view.action}`)]]} unavailable={unavailable} />
      {target && <div className="space-y-2">
        {inRouter && !navigateTo ? <Link to={target} className={linkClass}>{linkLabel}</Link>
          : <a href={target} onClick={onClick} className={linkClass}>{linkLabel}</a>}
        <p className="text-slate-600">{t(view.locator.step === 4 ? 'domain.groupSettingsHelp' : 'domain.currentSettingsHelp')}</p>
        <p className="text-slate-600">{t('domain.navigationOnly')}</p>
      </div>}
      <details><summary className={`${focus} cursor-pointer text-blue-800`}>{t('domain.details')}</summary>
        <div className="mt-3 space-y-4 rounded-lg border border-slate-200 bg-slate-50 p-3">
          <Facts items={detailFacts} unavailable={unavailable} />
          {receipt && <section aria-label={t('domain.receipt')} className="space-y-2">
            <h3 className="font-medium text-slate-900">{t('domain.receipt')}</h3>
            <Facts items={[[t('domain.recordId'), receipt.record_id], [t('domain.effectKey'), receipt.effect_key],
              [t('domain.payloadDigest'), receipt.payload_digest], [t('domain.committedAt'), receipt.committed_at],
              [t('domain.receiptGroupRevision'), receipt.group_revision], [t('domain.receiptConnectorRevision'), receipt.connector_revision],
              [t('domain.receiptProvenance'), t('domain.runtimeProvenance')]]} unavailable={unavailable} />
          </section>}
          {view.intakeState !== null && <Facts items={[[t('domain.intakeState'), t(`domain.intakeStates.${view.intakeState}`)]]} unavailable={unavailable} />}
          {view.pendingSamples !== null && <Facts items={[[t('domain.pendingSamples'), view.pendingSamples]]} unavailable={unavailable} />}
          {quota && <section aria-label={t('domain.quota')} className="space-y-2">
            <h3 className="font-medium text-slate-900">{t('domain.quota')}</h3>
            <Facts items={[[t('domain.quotaScope'), t(`domain.quotaScopes.${quota.scope}`)],
              [t('domain.quotaState'), t(`domain.quotaStates.${quota.state}`)],
              [t('domain.quotaUsed'), String(quota.used_bytes)], [t('domain.quotaMax'), String(quota.max_bytes)],
              [t('domain.intakeRefused'), t(quota.intake_refused ? 'domain.yes' : 'domain.no')]]} unavailable={unavailable} />
          </section>}
        </div>
      </details>
    </td>
  </tr>;
}
