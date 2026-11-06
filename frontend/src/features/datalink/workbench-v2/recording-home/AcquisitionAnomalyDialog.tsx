import * as Dialog from '@radix-ui/react-dialog';
import { useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { TFunction } from 'i18next';
import '@/styles/semantic-modal.css';
import { captureReturnFocus, restoreReturnFocus, type ReturnFocusTicket } from '@/utils/returnFocus';
import { useDeviceAcquisitionAnomaliesQuery } from '../../../../hooks/datalink/useAcquisitionAnomalies';
import { AcquisitionAnomalyError, acquisitionAnomalyErrorKey } from '../../../../services/acquisitionAnomalies';
import {
  ACQUISITION_ANOMALY_KINDS, ACQUISITION_ANOMALY_STATES,
  type AcquisitionAnomalyFilter, type AcquisitionAnomalyItem, type AcquisitionAnomalyKind, type AcquisitionAnomalyState,
} from '../../../../types/acquisitionAnomaly';
import './acquisitionAnomalies.css';

/** Bounded page size, following the home 25/50 paging rule (design condition). */
export const ACQUISITION_ANOMALY_PAGE_SIZE = 25;

function formatDuration(ms: number, t: TFunction): string {
  if (ms < 1000) return t('acquisitionAnomalies.duration.underSecond');
  const total = Math.floor(ms / 1000), hours = Math.floor(total / 3600), minutes = Math.floor((total % 3600) / 60), seconds = total % 60;
  return [hours ? t('acquisitionAnomalies.duration.hours', { count: hours }) : '', minutes ? t('acquisitionAnomalies.duration.minutes', { count: minutes }) : '',
    seconds ? t('acquisitionAnomalies.duration.seconds', { count: seconds }) : ''].filter(Boolean).join(' ');
}

function EventRow({ item, format, t, now }: { item: AcquisitionAnomalyItem; format: Intl.DateTimeFormat; t: TFunction; now: number }) {
  const location = item.point_address ?? (item.point_id === null ? t('acquisitionAnomalies.deviceLevel') : t('acquisitionAnomalies.pointRemoved'));
  const value = item.last_value === null ? t('acquisitionAnomalies.noValue') : item.last_value === '' ? t('acquisitionAnomalies.emptyValue') : item.last_value;
  const duration = item.duration_ms ?? Math.max(0, now - Date.parse(item.started_at));
  const label = (column: string) => t(`acquisitionAnomalies.columns.${column}`);
  return <tr data-anomaly-state={item.state}>
    <td data-label={label('kind')}><div><strong>{t(`acquisitionAnomalies.kinds.${item.kind}`)}</strong>
      {item.detail_codes.length > 0 && <p className="gw-anomaly-reason">{item.detail_codes.map(code => t(`acquisitionAnomalies.reasons.${code}`)).join(t('acquisitionAnomalies.reasonSeparator'))}</p>}</div></td>
    <td data-label={label('state')}><span className="gw-anomaly-state" data-anomaly-state={item.state}>{t(`acquisitionAnomalies.states.${item.state}`)}</span></td>
    <td data-label={label('startedAt')}><time dateTime={item.started_at}>{format.format(new Date(item.started_at))}</time></td>
    <td data-label={label('endedAt')}>{item.ended_at === null ? t('acquisitionAnomalies.ongoing') : <time dateTime={item.ended_at}>{format.format(new Date(item.ended_at))}</time>}</td>
    <td data-label={label('duration')}>{formatDuration(duration, t)}</td>
    <td data-label={label('address')} title={item.point_name ?? undefined}><span className="gw-anomaly-address">{location}</span></td>
    <td data-label={label('lastValue')} className="gw-anomaly-value">{value}</td>
  </tr>;
}

/** Classified form modal listing one device's acquisition anomalies; read failures never render as an empty list. */
export function AcquisitionAnomalyDialog({ device, onClose }: { device: { id: string; name: string }; onClose: () => void }) {
  const { t, i18n } = useTranslation('recording-home');
  const [kind, setKind] = useState<AcquisitionAnomalyKind | ''>('');
  const [state, setState] = useState<AcquisitionAnomalyState | ''>('');
  const [page, setPage] = useState(1);
  const [now, setNow] = useState(Date.now);
  useEffect(() => {
    const timer = window.setInterval(() => setNow(Date.now()), 1000);
    return () => window.clearInterval(timer);
  }, []);
  // Close first, return focus to the trigger, then let the owner unmount the dialog.
  const [open, setOpen] = useState(true);
  const focusOwner = useRef({}), returnFocus = useRef<ReturnFocusTicket | null>(null);
  const filter: AcquisitionAnomalyFilter = { ...(kind ? { kind } : {}), ...(state ? { state } : {}), limit: ACQUISITION_ANOMALY_PAGE_SIZE, offset: (page - 1) * ACQUISITION_ANOMALY_PAGE_SIZE };
  const query = useDeviceAcquisitionAnomaliesQuery(device.id, filter);
  const format = new Intl.DateTimeFormat(i18n.language, { dateStyle: 'short', timeStyle: 'medium' });
  const pages = query.data ? Math.max(1, Math.ceil(query.data.total / ACQUISITION_ANOMALY_PAGE_SIZE)) : 1;
  const filtered = kind !== '' || state !== '';
  const error = query.error instanceof AcquisitionAnomalyError ? query.error : new AcquisitionAnomalyError('unknown', true);
  return <Dialog.Root open={open} onOpenChange={next => { if (!next) setOpen(false); }}>
    <Dialog.Portal>
      <Dialog.Overlay className="gw-modal-overlay" />
      <Dialog.Content className="gw-modal gw-anomaly-modal" data-size="form" data-testid="acquisition-anomaly-dialog"
        onOpenAutoFocus={() => { returnFocus.current = captureReturnFocus(focusOwner.current, document.activeElement instanceof HTMLElement ? document.activeElement : null); }}
        onCloseAutoFocus={event => { event.preventDefault(); restoreReturnFocus(returnFocus.current, focusOwner.current); onClose(); }}>
        <header className="gw-modal-header">
          <Dialog.Title>{t('acquisitionAnomalies.title', { device: device.name })}</Dialog.Title>
          <Dialog.Description className="gw-muted">{t('acquisitionAnomalies.description')}</Dialog.Description>
          <form role="search" aria-label={t('acquisitionAnomalies.filters')} className="gw-anomaly-filters" onSubmit={event => event.preventDefault()}>
            <label>{t('acquisitionAnomalies.kindFilter')}<select value={kind} onChange={event => { setKind(ACQUISITION_ANOMALY_KINDS.find(item => item === event.target.value) ?? ''); setPage(1); }}>
              <option value="">{t('acquisitionAnomalies.allKinds')}</option>
              {ACQUISITION_ANOMALY_KINDS.map(item => <option key={item} value={item}>{t(`acquisitionAnomalies.kinds.${item}`)}</option>)}
            </select></label>
            <label>{t('acquisitionAnomalies.stateFilter')}<select value={state} onChange={event => { setState(ACQUISITION_ANOMALY_STATES.find(item => item === event.target.value) ?? ''); setPage(1); }}>
              <option value="">{t('acquisitionAnomalies.allStates')}</option>
              {ACQUISITION_ANOMALY_STATES.map(item => <option key={item} value={item}>{t(`acquisitionAnomalies.states.${item}`)}</option>)}
            </select></label>
          </form>
        </header>
        <div className="gw-modal-body" aria-busy={query.isFetching}>
          {query.isPending && <p role="status">{t('acquisitionAnomalies.loading')}</p>}
          {query.isError && <div role="alert" className="gw-warning gw-anomaly-error">
            <p>{t(acquisitionAnomalyErrorKey(error))}</p>
            <button type="button" disabled={query.isFetching} onClick={() => { void query.refetch(); }}>{t('acquisitionAnomalies.retry')}</button>
          </div>}
          {query.isSuccess && (query.data.items.length === 0 && page > 1
            ? <div role="status"><p>{t('acquisitionAnomalies.pageEmpty')}</p><button type="button" onClick={() => setPage(1)}>{t('acquisitionAnomalies.firstPage')}</button></div>
            : query.data.items.length === 0
            ? <p role="status">{t(filtered ? 'acquisitionAnomalies.emptyFiltered' : 'acquisitionAnomalies.empty')}</p>
            : <div className="gw-anomaly-table-scroll" role="region" tabIndex={0} aria-label={t('acquisitionAnomalies.tableRegion')}>
              <table className="gw-anomaly-table">
                <caption className="sr-only">{t('acquisitionAnomalies.title', { device: device.name })}</caption>
                <thead><tr>{(['kind', 'state', 'startedAt', 'endedAt', 'duration', 'address', 'lastValue'] as const).map(column =>
                  <th key={column} scope="col">{t(`acquisitionAnomalies.columns.${column}`)}</th>)}</tr></thead>
                <tbody>{query.data.items.map(item => <EventRow key={item.id} item={item} format={format} t={t} now={now} />)}</tbody>
              </table>
            </div>)}
        </div>
        <footer className="gw-modal-footer">
          {query.isSuccess && query.data.total > ACQUISITION_ANOMALY_PAGE_SIZE ? <div>
            <span role="status">{t('acquisitionAnomalies.pageStatus', { page, pages, count: query.data.total })}</span>
            <button type="button" disabled={page <= 1} onClick={() => setPage(page - 1)}>{t('previous')}</button>
            <button type="button" disabled={page >= pages} onClick={() => setPage(page + 1)}>{t('next')}</button>
          </div> : <span />}
          <Dialog.Close asChild><button type="button">{t('acquisitionAnomalies.close')}</button></Dialog.Close>
        </footer>
      </Dialog.Content>
    </Dialog.Portal>
  </Dialog.Root>;
}
