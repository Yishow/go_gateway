import { Button } from '../../../../components/ui/button';
import { useTranslation } from 'react-i18next';
import { Link, useInRouterContext } from 'react-router-dom';
import type { WriteGroup, WriteGroupListResponse } from '../../../../types/studioV2WriteGroup';
import { buildStudioEntityLocatorTarget } from '../../../../services/studioEntityLocator';
import { configuredDestinationIdentity, configuredDestinationReason, type ConfiguredDestinationContext } from './configuredRecordingDestinations';

// Existing semantic tokens are defined for light and dark in index.css; no new theme authority.
const SECONDARY_BUTTON = 'border border-[var(--gw-control)] bg-[var(--gw-surface)] text-[var(--gw-text)] hover:bg-[var(--gw-active)] disabled:text-[var(--gw-muted)] disabled:border-[var(--gw-border)] disabled:hover:bg-[var(--gw-surface)] aria-disabled:text-[var(--gw-muted)] aria-disabled:border-[var(--gw-border)] aria-disabled:hover:bg-[var(--gw-surface)] aria-disabled:opacity-100 disabled:opacity-100 focus-visible:outline-solid focus-visible:outline-2 focus-visible:outline-[var(--gw-focus)] focus-visible:outline-offset-2 focus-visible:ring-0 focus-visible:ring-offset-0';

export interface ConfiguredDestinationObservation {
  data?: WriteGroupListResponse;
  loading: boolean;
  failed: boolean;
  retry: () => void;
  /** Existing scoped editor owner; locator routes use the original locator link instead. */
  open?: (group: WriteGroup) => void;
}
export function ConfiguredDestinationChooser({ observation, context, workspaceRevision, disabled, selected, onSelect, onReset }: {
  observation: ConfiguredDestinationObservation; context: ConfiguredDestinationContext; workspaceRevision: string;
  disabled: boolean; selected?: WriteGroup; onSelect: (group: WriteGroup) => void; onReset: () => void;
}) {
  const { t } = useTranslation('recording-editor');
  const inRouter = useInRouterContext();
  const data = observation.data;
  const status = observation.loading ? 'loading' : observation.failed ? 'error' : !data ? 'unavailable' :
    data.workspace_id !== context.workspaceId || data.workspace_revision !== workspaceRevision ? 'stale' : undefined;
  const groups = data?.groups ?? [];
  return <section className="recording-configured-destinations" data-testid="configured-destinations" aria-labelledby="configured-destinations-title">
    <h4 id="configured-destinations-title">{t('destinations.title')}</h4>
    <p className="recording-configured-destinations__note">{t('destinations.configuredOnly')}</p>
    {status && <div role="status"><p>{t(`destinations.${status}`)}</p><Button type="button" variant="secondary" size="sm" className={SECONDARY_BUTTON} disabled={observation.loading} onClick={observation.retry}>{t('destinations.retry')}</Button></div>}
    {!status && groups.length === 0 && <p role="status">{t('destinations.empty')}</p>}
    {selected && <div role="status"><p>{t('destinations.intentOnly')}</p><Button type="button" variant="secondary" size="sm" className={SECONDARY_BUTTON} disabled={disabled} onClick={onReset}>{t('destinations.reset')}</Button></div>}
    {groups.length > 0 && <ul>
      {groups.map(group => {
        const reason = status ?? configuredDestinationReason(group, context);
        const target = buildStudioEntityLocatorTarget({ step: 4, workspace_id: group.workspace_id, group_id: group.id,
          group_revision: group.revision, device_id: null, rule_id: null, mapping_id: null });
        return <li key={`${configuredDestinationIdentity(group)}:${group.id}`} data-testid={`configured-destination-${group.id}`}>
          <div className="recording-configured-destinations__identity">
            <h5 className="recording-configured-destinations__name">{group.name || t('destinations.unnamed')}</h5>
            <p><span className="text-[var(--gw-muted)]">{t('destinations.table')} </span><span className="[overflow-wrap:anywhere]">{group.destination.table_name || t('destinations.unassigned')}</span></p>
            <p className="recording-configured-destinations__meta">{t('destinations.database')} {group.destination.database} · {t('destinations.schema')} {group.destination.table_schema || '—'} · {group.destination.storage_strategy} · {group.layout_version || 'managed_wide_v1'} · {group.storage_mode || 'plain'}</p>
            <p className="recording-configured-destinations__meta">{t('destinations.connector')} <span className="font-mono" title={group.destination.connector_id}>{group.destination.connector_id.slice(0, 8)}</span> · {t('destinations.revision')} <span className="font-mono" title={group.revision}>{group.revision.slice(0, 8)}</span></p>
          </div>
          <div className="recording-configured-destinations__actions">
            {reason && <p>{t(`destinations.reason.${reason}`, { defaultValue: t(`destinations.${reason}`) })}</p>}
            <Button type="button" variant="secondary" size="sm" className={SECONDARY_BUTTON} disabled={disabled || Boolean(reason)} onClick={() => {
              if (!disabled && !status && !configuredDestinationReason(group, context)) onSelect(group);
            }}>{t('destinations.select')}</Button>
            {group.status !== 'deleted' && target && !status && (observation.open && group.workspace_id === context.workspaceId ?
              <Button type="button" variant="secondary" size="sm" className={SECONDARY_BUTTON} disabled={disabled} onClick={() => observation.open?.(group)}>{t('destinations.editOriginal')}</Button> :
              <Button asChild variant="secondary" size="sm" className={SECONDARY_BUTTON}>
                {inRouter ? <Link to={target} aria-disabled={disabled} tabIndex={disabled ? -1 : undefined} onClick={event => { if (disabled) event.preventDefault(); }}>{t('destinations.editOriginal')}</Link> :
                  <a href={target} aria-disabled={disabled} tabIndex={disabled ? -1 : undefined} onClick={event => { if (disabled) event.preventDefault(); }}>{t('destinations.editOriginal')}</a>}
              </Button>)}
          </div>
        </li>;
      })}
    </ul>}
  </section>;
}
