import { watchHomeReturnFocus } from './recordingHomeReturnFocus';
import { eligibleFocusTarget } from '@/utils/returnFocus';
import { RecordingHomeRecordRow } from './RecordingHomeRecordRow';
import type { recordingHomeResultSummary } from './recordingHomeResultSummary';
import { useLayoutEffect, useRef, useState } from 'react';
import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroup, WriteGroupStatus } from '../../../../types/studioV2WriteGroup';
import type { RecordingReadonly } from '../state/recordingWorkflowView';
import type { RecordingHomeGuidance } from './recordingHomeGuidance';
import { recordingCollection, deviceCollection, initialCollectionState, type HomeDeviceSummary, type RecordingCollectionState } from './recordingHomeCollections';
import type { HomePresentation } from './recordingHomePresentation';
import { AcquisitionAnomalyButton, type AnomalySummaryStatus } from './AcquisitionAnomalyButton';
import { AcquisitionAnomalyDialog } from './AcquisitionAnomalyDialog';
import { Cable, Cpu, Rows3 } from 'lucide-react';
import { DeviceAvailabilityBadge } from './RecordingHomeBadges';
import { RecordingHomeStatusDistribution } from './RecordingHomeStatusDistribution';
import { DeviceConnectionBadge } from './RecordingHomeDeviceStrip';
import type { DeviceConnectionObservation } from '../../../../types/deviceConnection';

export interface RecordingHomeCollectionProps {
  groups: readonly RecordingReadonly<WriteGroup>[]; devices: readonly HomeDeviceSummary[]; state: RecordingCollectionState;
  guidance: RecordingHomeGuidance; pendingDevices?: readonly string[]; summaries?: Readonly<Record<string, ReturnType<typeof recordingHomeResultSummary>>>; observedIds: readonly string[]; filterObservedIds?: readonly string[]; unverifiedObserved: boolean; onRecheck: () => void; sourcesAvailable: boolean; restored: HomePresentation;
  /** False while the first-load reading state hides the collection; restore waits so focus can land. */
  presentationReady?: boolean;
  /** Keep the current return-focus ticket only while initial home reads can change layout. */
  returnLayoutPending?: boolean;
  /** Observed presentation can restore focus without asserting current source freshness. */
  observedSourcesAvailable?: boolean;
  /** Active anomaly counts read by the route; absent means the status is unknown. */
  anomalySummary?: AnomalySummaryStatus;
  /** Called after the anomaly modal closes so the counts can be re-read. */
  onAnomaliesClosed?: () => void;
  connectionStates?: ReadonlyMap<string, DeviceConnectionObservation>;
  renderDeviceActions?: (device: HomeDeviceSummary) => ReactNode;
  onDelete?: (group: RecordingReadonly<WriteGroup>) => void;
  onChange: (state: RecordingCollectionState) => void; onOpen: (id: string) => void; onOpenDevice: (id: string, kind: 'connection' | 'points') => void;
  onRemember: (scrollTop: number, focusId: string | null, page: number) => void;
}
export function RecordingHomeCollection({ groups, devices, state, guidance, summaries, pendingDevices = [], observedIds, filterObservedIds = observedIds, unverifiedObserved, onRecheck, sourcesAvailable, observedSourcesAvailable = sourcesAvailable, presentationReady = true, returnLayoutPending = false, restored, anomalySummary, onAnomaliesClosed, connectionStates, renderDeviceActions, onChange, onOpen, onOpenDevice, onRemember, onDelete }: RecordingHomeCollectionProps) {
  const { t } = useTranslation('recording-home');
  const records = recordingCollection(groups, devices, state, new Set(filterObservedIds)), devicePage = deviceCollection(devices, state);
  const current = state.mode === 'devices' ? devicePage : records;
  const [anomalyDevice, setAnomalyDevice] = useState<{ id: string; name: string } | null>(null);
  const root = useRef<HTMLElement>(null), panel = useRef<HTMLDivElement>(null), heading = useRef<HTMLHeadingElement>(null), didRestore = useRef(false);
  const returnFocus = useRef<HTMLElement | null>(null), focusReveal = useRef<ReturnType<typeof watchHomeReturnFocus> | null>(null);
  useLayoutEffect(() => {
    if (didRestore.current || !observedSourcesAvailable || !presentationReady) return;
    didRestore.current = true;
    if (panel.current) {
      panel.current.scrollTop = Math.min(restored.scrollTop, Math.max(0, panel.current.scrollHeight - panel.current.clientHeight));
      const target = [...(root.current?.querySelectorAll<HTMLButtonElement>('[data-home-focus]') ?? [])].find(node => node.dataset.homeFocus === restored.focusId);
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
    <nav className="gw-home-object-modes" aria-label={t('objectModes')}>
      <button type="button" aria-pressed={state.mode === 'recordings'} onClick={() => patch({ mode: 'recordings' })}><Rows3 aria-hidden="true" size={16} />{t('recordingsTab')}</button>
      <button type="button" aria-pressed={state.mode === 'devices'} onClick={() => patch({ mode: 'devices' })}><Cpu aria-hidden="true" size={16} />{t('devicesTab')}</button>
    </nav>
    <p className="gw-object-purpose">{t(state.mode === 'devices' ? 'devicesPurpose' : 'recordingsPurpose')}</p>
    {state.mode === 'recordings' && <RecordingHomeStatusDistribution groups={groups} />}
    <h2 ref={heading} tabIndex={-1} className="text-lg font-semibold">{t(state.mode === 'devices' ? 'manageDevices' : 'saved')}</h2>
    <form role="search" aria-label={t('filters')} className="gw-filters" onSubmit={event => event.preventDefault()}>
      <label className="gw-filter-search">{t(state.mode === 'devices' ? 'deviceKeyword' : 'keyword')}<input type="search" value={state.keyword} onChange={event => patch({ keyword: event.target.value })} /></label>
      {state.mode === 'recordings' && <>

        <label>{t('deviceFilter')}<select value={state.deviceId} onChange={event => patch({ deviceId: event.target.value })}>
          <option value="">{t('allDevices')}</option>{devices.map(device => <option key={device.id} value={device.id}>{device.name}</option>)}
        </select></label>
        <label>{t('savedStatus')}<select value={state.status} onChange={event => patch({ status: event.target.value as WriteGroupStatus | '' })}>
          <option value="">{t('allStatuses')}</option>{(['draft', 'ready', 'running', 'disabled', 'deleted'] as const).map(status => <option key={status} value={status}>{t(`savedStatuses.${status}`)}</option>)}
        </select></label>
        <label>{t('destinationFilter')}<select value={state.destination} onChange={event => patch({ destination: event.target.value })}>
          <option value="">{t('allDestinations')}</option>{destinations.map(([id, name]) => <option key={id} value={id}>{name}</option>)}
        </select></label>
      </>}
      <div className="gw-filter-options">{state.mode === 'recordings' && <label className="gw-checkbox"><input type="checkbox" checked={state.observedOnly} onChange={event => patch({ observedOnly: event.target.checked })} />{t('observedFilter')}</label>}
      <button type="button" onClick={() => onChange({ ...initialCollectionState(), mode: state.mode, pageSize: state.pageSize })}>{t('clearFilters')}</button></div>
    </form>
    {state.observedOnly && <div role="status"><p>{t(unverifiedObserved ? 'observationNeedsRefresh' : 'observedFilterBoundary')}</p>
      <button type="button" onClick={onRecheck}>{t('recheckVisible')}</button></div>}
    <div ref={panel} role="region" aria-label={t('tableRegion')} tabIndex={0} className="gw-table-scroll" data-testid="home-collection-scroll"
      onScroll={() => {
        // Restoring scrollTop dispatches a later scroll event. Keep the actual restored focus
        // so a subsequent Back/Forward visit does not replace its identity with null.
        const focused = root.current?.querySelector<HTMLElement>('[data-home-focus]:focus');
        remember(focused?.dataset.homeFocus ?? null);
      }}>
      {current.empty ? <p role="status">{t(!sourcesAvailable ? 'unknown' : current.empty === 'empty' ? 'emptyCollection' : 'noResults')}</p> :
        <table data-testid="home-management-table" data-mode={state.mode}>
          <caption className="sr-only">{t(state.mode === 'devices' ? 'manageDevices' : 'saved')}</caption>
          <thead><tr><th scope="col">{t('name')}</th><th scope="col">{t(state.mode === 'devices' ? 'protocol' : 'savedStatus')}</th>
            {state.mode === 'devices' ? <><th scope="col">{t('availability')}</th><th scope="col">{t('deviceStatus.column')}</th></> : <><th scope="col">{t('recordSources')}</th><th scope="col">{t('destinationFilter')}</th><th scope="col">{t('recordWrites')}</th></>}<th scope="col">{t('actions')}</th></tr></thead>
          <tbody>{state.mode === 'devices' ? devicePage.rows.map(device => <tr key={device.id}>
            <td><strong>{device.name}</strong></td>
            <td><span className="gw-home-chip" data-protocol={device.protocol}><Cable aria-hidden="true" size={14} />{t(`protocolNames.${device.protocol}`, { defaultValue: device.protocol })}</span></td><td><DeviceAvailabilityBadge availability={device.availability} /></td><td><DeviceConnectionBadge state={connectionStates?.get(device.id)?.state ?? 'unavailable'} /></td><td className="gw-row-actions">
              <button type="button" data-home-focus={`device-${device.id}`} onClick={() => { remember(`device-${device.id}`); onOpenDevice(device.id, 'connection'); }}>{t('deviceConnection')}</button>
              <button type="button" onClick={() => { remember(`device-${device.id}`); onOpenDevice(device.id, 'points'); }}>{t('devicePoints')}</button>
              <AcquisitionAnomalyButton deviceId={device.id} status={anomalySummary} onOpen={() => setAnomalyDevice({ id: device.id, name: device.name })} />{renderDeviceActions?.(device)}</td>
          </tr>) : records.rows.map(group => <RecordingHomeRecordRow key={group.id} group={group} devices={devices} summary={summaries?.[group.id]}
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

      <button type="button" disabled={current.page <= 1} onClick={() => onChange({ ...state, page: current.page - 1 })}>{t('previous')}</button>
      <button type="button" disabled={current.page >= current.pageCount} onClick={() => onChange({ ...state, page: current.page + 1 })}>{t('next')}</button>
    </div></footer>
    {anomalyDevice && <AcquisitionAnomalyDialog device={anomalyDevice} onClose={() => { setAnomalyDevice(null); onAnomaliesClosed?.(); }} />}
  </section>;
}
