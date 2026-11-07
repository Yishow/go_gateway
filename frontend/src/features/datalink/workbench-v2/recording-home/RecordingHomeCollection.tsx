import { watchHomeReturnFocus } from './recordingHomeReturnFocus';
import { eligibleFocusTarget } from '@/utils/returnFocus';
import { RecordingHomeRecordRow } from './RecordingHomeRecordRow';
import type { recordingHomeResultSummary } from './recordingHomeResultSummary';
import { useLayoutEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroup, WriteGroupStatus } from '../../../../types/studioV2WriteGroup';
import type { RecordingReadonly } from '../state/recordingWorkflowView';
import type { RecordingHomeGuidance } from './recordingHomeGuidance';
import { recordingCollection, initialCollectionState, type HomeDeviceSummary, type RecordingCollectionState } from './recordingHomeCollections';
import type { HomePresentation } from './recordingHomePresentation';
import { ChevronLeft, ChevronRight, Search, X } from 'lucide-react';

export interface RecordingHomeCollectionProps {
  groups: readonly RecordingReadonly<WriteGroup>[]; devices: readonly HomeDeviceSummary[]; state: RecordingCollectionState;
  guidance: RecordingHomeGuidance; pendingDevices?: readonly string[]; summaries?: Readonly<Record<string, ReturnType<typeof recordingHomeResultSummary>>>; observedIds: readonly string[]; filterObservedIds?: readonly string[]; unverifiedObserved: boolean; onRecheck: () => void; sourcesAvailable: boolean; restored: HomePresentation;
  /** False while the first-load reading state hides the collection; restore waits so focus can land. */
  presentationReady?: boolean;
  /** Keep the current return-focus ticket only while initial home reads can change layout. */
  returnLayoutPending?: boolean;
  /** Observed presentation can restore focus without asserting current source freshness. */
  observedSourcesAvailable?: boolean;
  onDelete?: (group: RecordingReadonly<WriteGroup>) => void;
  onChange: (state: RecordingCollectionState) => void; onOpen: (id: string) => void; onOpenDevice: (id: string, kind: 'connection' | 'points') => void;
  onRemember: (scrollTop: number, focusId: string | null, page: number) => void;
}
export function RecordingHomeCollection({ groups, devices, state, guidance, summaries, pendingDevices = [], observedIds, filterObservedIds = observedIds, unverifiedObserved, onRecheck, sourcesAvailable, observedSourcesAvailable = sourcesAvailable, presentationReady = true, returnLayoutPending = false, restored, onChange, onOpen, onOpenDevice, onRemember, onDelete }: RecordingHomeCollectionProps) {
  const { t } = useTranslation('recording-home');
  // Devices are managed in the device strip above; this panel always lists recordings.
  const current = recordingCollection(groups, devices, state, new Set(filterObservedIds));
  const root = useRef<HTMLElement>(null), panel = useRef<HTMLDivElement>(null), heading = useRef<HTMLHeadingElement>(null), didRestore = useRef(false);
  const returnFocus = useRef<HTMLElement | null>(null), focusReveal = useRef<ReturnType<typeof watchHomeReturnFocus> | null>(null);
  useLayoutEffect(() => {
    if (didRestore.current || !observedSourcesAvailable || !presentationReady) return;
    didRestore.current = true;
    if (panel.current) {
      panel.current.scrollTop = Math.min(restored.scrollTop, Math.max(0, panel.current.scrollHeight - panel.current.clientHeight));
      const home = root.current?.closest('[data-testid="recording-home"]') ?? root.current;
      const target = [...(home?.querySelectorAll<HTMLButtonElement>('[data-home-focus]') ?? [])].find(node => node.dataset.homeFocus === restored.focusId);
      if (restored.focusId) {
        const focus = eligibleFocusTarget(target) ? target : eligibleFocusTarget(heading.current) ? heading.current : null;
        if (focus) { focus.focus({ preventScroll: true }); returnFocus.current = focus; }
      }
    }
  }, [observedSourcesAvailable, presentationReady, restored]);
  useLayoutEffect(() => {
    const target = returnFocus.current;
    if (!target || !observedSourcesAvailable || !presentationReady) return;
    const ticket = watchHomeReturnFocus(target); focusReveal.current = ticket;
    return () => ticket.cancel();
  }, [observedSourcesAvailable, presentationReady]);
  useLayoutEffect(() => {
    if (!returnLayoutPending) focusReveal.current?.finish();
  }, [returnLayoutPending, observedSourcesAvailable, presentationReady]);
  const patch = (update: Partial<RecordingCollectionState>) => onChange({ ...state, ...update, page: 1 });
  const remember = (id: string | null) => onRemember(panel.current?.scrollTop ?? 0, id, current.page);
  const destinations = [...new Map(groups.map(group => [group.destination.connector_id, group.destination.database || group.destination.connector_id])).entries()];
  const attention = guidance.attentionItems.filter(item => item.reason !== 'evidence_unavailable');
  return <section ref={root} aria-label={t('management')} className="gw-management">
    {pendingDevices.length > 0 && <aside className="gw-summary"><p>{t('result.devicePending')}</p>{pendingDevices.slice(0, 4).map(id => <button type="button" key={id} onClick={() => onOpenDevice(id, 'connection')}>{devices.find(device => device.id === id)?.name ?? id}</button>)}</aside>}
    <p className="gw-observation-note" data-testid="recording-observation-boundary">{t('observationBoundary')}</p>
    {!!observedIds.length && <aside className="gw-summary" aria-label={t('observedAttentionTitle')}>
      <h2 className="text-lg font-semibold">{t('observedAttentionTitle')}</h2>
      <ul>{attention.slice(0, 4).map(item => <li key={item.id}><button type="button" data-home-focus={`attention-record-${item.id}`} onClick={() => { remember(`attention-record-${item.id}`); onOpen(item.id); }}>{item.name}</button></li>)}</ul>
      {unverifiedObserved && <p>{t('observationNeedsRefresh')}</p>}
      <button type="button" onClick={() => onChange({ ...initialCollectionState(), pageSize: state.pageSize, observedOnly: true })}>{t('viewObserved')}</button>
    </aside>}
    <div className="gw-records-panel">
    <div className="gw-section-head">
      <div className="gw-section-title"><h2 ref={heading} tabIndex={-1}>{t('saved')}</h2><span className="gw-count">{groups.length}</span></div>
    </div>
    <form role="search" aria-label={t('filters')} className="gw-filters" onSubmit={event => event.preventDefault()}>
      <label className="gw-filter-search"><span className="gw-field-label">{t('keyword')}</span><span className="gw-search-field"><Search aria-hidden="true" size={16} />
        <input type="search" value={state.keyword} placeholder={t('keywordPlaceholder')} onChange={event => patch({ keyword: event.target.value })} /></span></label>
      <label><span className="gw-field-label">{t('deviceFilter')}</span><select value={state.deviceId} onChange={event => patch({ deviceId: event.target.value })}>
        <option value="">{t('allDevices')}</option>{devices.map(device => <option key={device.id} value={device.id}>{device.name}</option>)}
      </select></label>
      <label><span className="gw-field-label">{t('savedStatus')}</span><select value={state.status} onChange={event => patch({ status: event.target.value as WriteGroupStatus | '' })}>
        <option value="">{t('allStatuses')}</option>{(['draft', 'ready', 'running', 'disabled', 'deleted'] as const).map(status => <option key={status} value={status}>{t(`savedStatuses.${status}`)}</option>)}
      </select></label>
      <label><span className="gw-field-label">{t('destinationFilter')}</span><select value={state.destination} onChange={event => patch({ destination: event.target.value })}>
        <option value="">{t('allDestinations')}</option>{destinations.map(([id, name]) => <option key={id} value={id}>{name}</option>)}
      </select></label>
      <div className="gw-filter-options"><label className="gw-checkbox"><input type="checkbox" checked={state.observedOnly} onChange={event => patch({ observedOnly: event.target.checked })} />{t('observedFilter')}</label>
      <button type="button" className="gw-ghost" onClick={() => onChange({ ...initialCollectionState(), pageSize: state.pageSize })}><X aria-hidden="true" size={15} />{t('clearFilters')}</button></div>
    </form>
    {state.observedOnly && <div role="status"><p>{t(unverifiedObserved ? 'observationNeedsRefresh' : 'observedFilterBoundary')}</p>
      <button type="button" onClick={onRecheck}>{t('recheckVisible')}</button></div>}
    <div ref={panel} role="region" aria-label={t('tableRegion')} tabIndex={0} className="gw-table-scroll" data-testid="home-collection-scroll"
      onScroll={() => {
        // Restoring scrollTop dispatches a later scroll event. Keep the actual restored focus
        // so a subsequent Back/Forward visit does not replace its identity with null.
        const home = root.current?.closest('[data-testid="recording-home"]') ?? root.current;
        const focused = home?.querySelector<HTMLElement>('[data-home-focus]:focus');
        remember(focused?.dataset.homeFocus ?? null);
      }}>
      {current.empty ? <p role="status">{t(!sourcesAvailable ? 'unknown' : current.empty === 'empty' ? 'emptyCollection' : 'noResults')}</p> :
        <table data-testid="home-management-table" data-mode="recordings">
          <caption className="sr-only">{t('saved')}</caption>
          <thead><tr><th scope="col">{t('name')}</th><th scope="col">{t('recordSources')}</th><th scope="col">{t('destinationFilter')}</th><th scope="col">{t('recordWrites')}</th><th scope="col">{t('actions')}</th></tr></thead>
          <tbody>{current.rows.map(group => <RecordingHomeRecordRow key={group.id} group={group} devices={devices} summary={summaries?.[group.id]}
            attention={guidance.attentionItems.some(item => item.id === group.id && item.reason !== 'evidence_unavailable')}
            evidence={guidance.recordingItems.find(item => item.id === group.id)} onDelete={onDelete} onRecheck={onRecheck}
            onOpen={() => { remember(`record-${group.id}`); onOpen(group.id); }} />)}</tbody>
        </table>}
    </div>
    <footer className="gw-pagination"><div className="gw-pagination-summary"><span role="status" data-testid="home-page-status">{t('pageOf', { page: current.page, pages: current.pageCount, count: current.filteredCount })}</span>
      <label>{t('sort')}<select value={state.descending ? 'desc' : 'asc'} onChange={event => patch({ descending: event.target.value === 'desc' })}>
        <option value="asc">{t('ascending')}</option><option value="desc">{t('descending')}</option>
      </select></label>
      </div><div className="gw-pagination-controls"><label>{t('pageSize')}<select value={state.pageSize} onChange={event => patch({ pageSize: event.target.value === '50' ? 50 : 25 })}>
        <option value="25">25</option><option value="50">50</option>
      </select></label>

      <button type="button" className="gw-page-button" disabled={current.page <= 1} onClick={() => onChange({ ...state, page: current.page - 1 })}><ChevronLeft aria-hidden="true" size={16} />{t('previous')}</button>
      <button type="button" className="gw-page-button" disabled={current.page >= current.pageCount} onClick={() => onChange({ ...state, page: current.page + 1 })}>{t('next')}<ChevronRight aria-hidden="true" size={16} /></button>
    </div></footer>
    </div>
  </section>;
}
