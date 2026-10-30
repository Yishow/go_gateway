import { useLayoutEffect, useRef } from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroup, WriteGroupStatus } from '../../../../types/studioV2WriteGroup';
import type { RecordingReadonly } from '../state/recordingWorkflowView';
import type { RecordingHomeGuidance } from './recordingHomeGuidance';
import { recordingCollection, deviceCollection, initialCollectionState, type HomeDeviceSummary, type RecordingCollectionState } from './recordingHomeCollections';
import type { HomePresentation } from './recordingHomePresentation';

export interface RecordingHomeCollectionProps {
  groups: readonly RecordingReadonly<WriteGroup>[]; devices: readonly HomeDeviceSummary[]; state: RecordingCollectionState;
  guidance: RecordingHomeGuidance; observedIds: readonly string[]; filterObservedIds?: readonly string[]; unverifiedObserved: boolean; onRecheck: () => void; sourcesAvailable: boolean; restored: HomePresentation;
  onChange: (state: RecordingCollectionState) => void; onOpen: (id: string) => void; onOpenDevice: (id: string) => void;
  onRemember: (scrollTop: number, focusId: string | null, page: number) => void;
}
export function RecordingHomeCollection({ groups, devices, state, guidance, observedIds, filterObservedIds = observedIds, unverifiedObserved, onRecheck, sourcesAvailable, restored, onChange, onOpen, onOpenDevice, onRemember }: RecordingHomeCollectionProps) {
  const { t } = useTranslation('recording-home');
  const records = recordingCollection(groups, devices, state, new Set(filterObservedIds)), devicePage = deviceCollection(devices, state);
  const current = state.mode === 'devices' ? devicePage : records;
  const root = useRef<HTMLElement>(null), panel = useRef<HTMLDivElement>(null), heading = useRef<HTMLHeadingElement>(null), didRestore = useRef(false);
  useLayoutEffect(() => {
    if (didRestore.current || !sourcesAvailable) return;
    didRestore.current = true;
    if (panel.current) {
      panel.current.scrollTop = Math.min(restored.scrollTop, Math.max(0, panel.current.scrollHeight - panel.current.clientHeight));
      const target = [...(root.current?.querySelectorAll<HTMLButtonElement>('[data-home-focus]') ?? [])].find(node => node.dataset.homeFocus === restored.focusId);
      if (restored.focusId) (target ?? heading.current)?.focus({ preventScroll: true });
    }
  }, [sourcesAvailable, restored]);
  const patch = (update: Partial<RecordingCollectionState>) => onChange({ ...state, ...update, page: 1 });
  const remember = (id: string | null) => onRemember(panel.current?.scrollTop ?? 0, id, current.page);
  const destinations = [...new Map(groups.map(group => [group.destination.connector_id, group.destination.database || group.destination.connector_id])).entries()];
  const attention = guidance.attentionItems.filter(item => item.reason !== 'evidence_unavailable');
  return <section ref={root} aria-label={t('management')} className="space-y-4">
    <p data-testid="recording-observation-boundary">{t('observationBoundary')}</p>
    {!!observedIds.length && <aside aria-label={t('observedAttentionTitle')}>
      <h2 className="text-lg font-semibold">{t('observedAttentionTitle')}</h2>
      <ul>{attention.slice(0, 4).map(item => <li key={item.id}><button type="button" onClick={() => { remember(`record-${item.id}`); onOpen(item.id); }}>{item.name}</button></li>)}</ul>
      {unverifiedObserved && <p>{t('observationNeedsRefresh')}</p>}
      <button type="button" onClick={() => onChange({ ...initialCollectionState(), pageSize: state.pageSize, observedOnly: true })}>{t('viewObserved')}</button>
    </aside>}
    <aside aria-label={t('deviceSummary')}>
      <h2 className="text-lg font-semibold">{t('deviceSummary')}</h2>
      <ul data-testid="home-device-summary">{devicePage.summary.map(device => <li key={device.id}><button type="button" data-home-focus={`summary-device-${device.id}`} onClick={() => { remember(`summary-device-${device.id}`); onOpenDevice(device.id); }}>{device.name}</button></li>)}</ul>
      <button type="button" onClick={() => patch({ mode: state.mode === 'devices' ? 'recordings' : 'devices' })}>{t(state.mode === 'devices' ? 'showRecordings' : 'manageDevices')}</button>
    </aside>
    <h2 ref={heading} tabIndex={-1} className="text-lg font-semibold">{t(state.mode === 'devices' ? 'manageDevices' : 'saved')}</h2>
    <div className="flex flex-wrap gap-3">
      <label>{t('keyword')}<input value={state.keyword} onChange={event => patch({ keyword: event.target.value })} /></label>
      {state.mode === 'recordings' && <>
        <label><input type="checkbox" checked={state.observedOnly} onChange={event => patch({ observedOnly: event.target.checked })} />{t('observedFilter')}</label>
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
      <label>{t('sort')}<select value={state.descending ? 'desc' : 'asc'} onChange={event => patch({ descending: event.target.value === 'desc' })}>
        <option value="asc">{t('ascending')}</option><option value="desc">{t('descending')}</option>
      </select></label>
      <label>{t('pageSize')}<select value={state.pageSize} onChange={event => patch({ pageSize: event.target.value === '50' ? 50 : 25 })}>
        <option value="25">25</option><option value="50">50</option>
      </select></label>
      <button type="button" onClick={() => onChange({ ...initialCollectionState(), mode: state.mode, pageSize: state.pageSize })}>{t('clearFilters')}</button>
    </div>
    {state.observedOnly && <div role="status"><p>{t(unverifiedObserved ? 'observationNeedsRefresh' : 'observedFilterBoundary')}</p>
      <button type="button" onClick={onRecheck}>{t('recheckVisible')}</button></div>}
    <div ref={panel} className="max-h-[60dvh] overflow-auto" data-testid="home-collection-scroll"
      onScroll={() => remember(null)}>
      {current.empty ? <p role="status">{t(!sourcesAvailable ? 'unknown' : current.empty === 'empty' ? 'emptyCollection' : 'noResults')}</p> :
        <table className="w-full text-sm" data-testid="home-management-table">
          <thead><tr><th scope="col">{t('name')}</th><th scope="col">{t(state.mode === 'devices' ? 'protocol' : 'savedStatus')}</th><th scope="col">{t(state.mode === 'devices' ? 'availability' : 'table')}</th></tr></thead>
          <tbody>{state.mode === 'devices' ? devicePage.rows.map(device => <tr key={device.id}>
            <td><button type="button" data-home-focus={`device-${device.id}`} onClick={() => { remember(`device-${device.id}`); onOpenDevice(device.id); }}>{device.name}</button></td>
            <td>{device.protocol}</td><td>{t(`availabilityStates.${device.availability}`, { defaultValue: t('unknown') })}</td>
          </tr>) : records.rows.map(group => <tr key={group.id}>
            <td><button type="button" data-home-focus={`record-${group.id}`} onClick={() => { remember(`record-${group.id}`); onOpen(group.id); }}>{group.name}</button></td>
            <td>{t(`savedStatuses.${group.status}`)}{guidance.attentionItems.some(item => item.id === group.id && item.reason !== 'evidence_unavailable') ? ` · ${t('observedAttentionShort')}` : ` · ${t('statusUnknown')}`}</td>
            <td>{[group.destination.table_schema, group.destination.table_name].filter(Boolean).join('.')}</td>
          </tr>)}</tbody>
        </table>}
    </div>
    <footer className="flex items-center gap-4">
      <button type="button" disabled={current.page <= 1} onClick={() => onChange({ ...state, page: current.page - 1 })}>{t('previous')}</button>
      <span role="status" data-testid="home-page-status">{t('pageOf', { page: current.page, pages: current.pageCount, count: current.filteredCount })}</span>
      <button type="button" disabled={current.page >= current.pageCount} onClick={() => onChange({ ...state, page: current.page + 1 })}>{t('next')}</button>
    </footer>
  </section>;
}
