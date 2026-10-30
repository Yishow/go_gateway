import { useId, useLayoutEffect, useRef } from 'react';
import { cancelPicker, filterPicker, pendingAddition, pickerResults, sourceKey, togglePending } from './recordingPickerState';
import './recordingPointPicker.css';
import type { Device } from '../state/types';
import type { PickerProjection, PickerState, PointFilter, SourceIdentity } from './recordingPickerState';
export interface PickerLabels {
  title: string; search: string; device: string; pointType: string; allDevices: string; allTypes: string;
  types: Record<Exclude<PointFilter, 'all'>, string>; unknownType: string;
  point: string; address: string; alreadyAdded: string; unavailable: string; cancel: string;
  loading: string; failed: string; retry: string; empty: string; noMatches: string; clearFilters: string;
  newSource: string; invalidPending: string; previous: string; next: string; pageSize: string;
  matching: (count: number, added: number) => string;
  pending: (count: number) => string; hidden: (count: number) => string; add: (count: number) => string;
  page: (current: number, total: number) => string;
  removePending: (identity: SourceIdentity) => string;
}
export interface RecordingPointPickerProps {
  state: PickerState; projection: PickerProjection; devices: readonly Pick<Device, 'id' | 'name'>[];
  status: 'loading' | 'failed' | 'loaded'; labels: PickerLabels;
  onChange: (state: PickerState) => void;
  onCancel: (reason: 'cancel' | 'escape') => void;
  onAdd: (identities: readonly SourceIdentity[]) => void;
  onSource: (deviceId: string | undefined) => void; onRetry: () => void;
}
/** Content only. The shared shell owns modal sizing, close/dirty guards and focus trap/return. */
export function RecordingPointPicker({ state, projection, devices, status, labels, onChange, onCancel, onAdd, onSource, onRetry }: RecordingPointPickerProps) {
  const id = useId();
  const root = useRef<HTMLElement>(null), body = useRef<HTMLDivElement>(null);
  const restoredScope = useRef<string | undefined>(undefined);
  useLayoutEffect(() => {
    if (status !== 'loaded' || state.scope !== projection.scope) return;
    // Controlled result navigation must move the DOM, not only the remembered position.
    if (body.current) body.current.scrollTop = state.scrollTop;
    if (restoredScope.current === state.scope) return;
    restoredScope.current = state.scope;
    const target = [...(root.current?.querySelectorAll<HTMLElement>('[data-picker-focus]') ?? [])]
      .find(element => element.dataset.pickerFocus === state.focusKey);
    target?.focus({ preventScroll: true });
  }, [status, state.scope, state.scrollTop, state.focusKey, projection.scope]);
  const result = pickerResults(state, projection);
  const active = status === 'loaded' && state.scope === projection.scope;
  const members = new Set(projection.members.map(sourceKey));
  const pending = new Set(state.pending.map(sourceKey));
  const narrowed = Boolean(state.query.trim() || state.pointType !== 'all');
  const cancel = (reason: 'cancel' | 'escape') => { onChange(cancelPicker(state)); onCancel(reason); };
  const patch = (value: Parameters<typeof filterPicker>[1]) => onChange(filterPicker(state, value));
  return <section ref={root} className="recording-point-picker" aria-label={labels.title} onFocusCapture={event => {
    const key = (event.target as HTMLElement).dataset.pickerFocus;
    if (key && key !== state.focusKey) onChange({ ...state, focusKey: key });
  }} onKeyDown={event => {
    if (event.key === 'Escape') { event.preventDefault(); event.stopPropagation(); cancel('escape'); }
  }}>
    <div className="recording-point-picker__toolbar">
      <label className="recording-point-picker__search" htmlFor={`${id}-search`}>{labels.search}
        <input data-picker-focus="search" id={`${id}-search`} type="search" value={state.query} onChange={event => patch({ query: event.target.value })} />
      </label>
      <label htmlFor={`${id}-device`}>{labels.device}
        <select data-picker-focus="device" id={`${id}-device`} value={state.deviceId} onChange={event => patch({ deviceId: event.target.value })}>
          <option value="">{labels.allDevices}</option>
          {devices.map(device => <option key={device.id} value={device.id}>{device.name}</option>)}
        </select>
      </label>
      <label htmlFor={`${id}-type`}>{labels.pointType}
        <select data-picker-focus="type" id={`${id}-type`} value={state.pointType} onChange={event => {
          const type = event.target.value;
          if (type === 'all' || type === 'gauge' || type === 'counter' || type === 'state') patch({ pointType: type });
        }}>
          <option value="all">{labels.allTypes}</option>
          {(['gauge', 'counter', 'state'] as const).map(type => <option key={type} value={type}>{labels.types[type]}</option>)}
        </select>
      </label>
    </div>
    <div ref={body} className="recording-point-picker__body" onScroll={event => onChange({ ...state, scrollTop: event.currentTarget.scrollTop })}>
      {status === 'loading' && <p role="status">{labels.loading}</p>}
      {status === 'failed' && <div role="alert"><p>{labels.failed}</p><button type="button" onClick={onRetry}>{labels.retry}</button></div>}
      {status === 'loaded' && !active && <p role="alert">{labels.unavailable}</p>}
      {active && <>
        <p aria-live="polite">{labels.matching(result.total, result.added)}</p>
        {result.total === 0 ? <div>
          <p>{narrowed ? labels.noMatches : labels.empty}</p>
          {narrowed && <button type="button" onClick={() => patch({ query: '', pointType: 'all' })}>{labels.clearFilters}</button>}
        </div> : <table>
          <thead><tr><th scope="col">{labels.point}</th><th scope="col">{labels.device}</th><th scope="col">{labels.address}</th><th scope="col">{labels.pointType}</th></tr></thead>
          <tbody>{result.rows.map(row => {
            const key = sourceKey(row), added = members.has(key), evidence = projection.evidence[key];
            const semantic = evidence?.semantic;
            const meaning = semantic === 'signed_counter' ? 'counter' : semantic;
            const typeLabel = meaning === 'gauge' || meaning === 'counter' || meaning === 'state' ? labels.types[meaning] : labels.unknownType;
            const unavailable = evidence?.availability !== 'eligible';
            return <tr key={key}>
              <td><label><input data-picker-focus={key} type="checkbox" aria-label={row.label} checked={added || pending.has(key)} disabled={added || unavailable && !pending.has(key)}
                onChange={() => onChange(togglePending(state, row, projection))} />{row.label}</label>
                {added ? <small>{labels.alreadyAdded}</small> : unavailable && <small>{evidence?.reason || labels.unavailable}</small>}</td>
              <td>{row.device_name}</td><td>{row.address}</td><td>{typeLabel}</td>
            </tr>;
          })}</tbody>
        </table>}
        {result.pages > 1 && <nav aria-label={labels.page(result.page + 1, result.pages)}>
          <button type="button" disabled={result.page === 0} onClick={() => onChange({ ...state, page: result.page - 1, scrollTop: 0 })}>{labels.previous}</button>
          <span>{labels.page(result.page + 1, result.pages)}</span>
          <button type="button" disabled={result.page === result.pages - 1} onClick={() => onChange({ ...state, page: result.page + 1, scrollTop: 0 })}>{labels.next}</button>
          <label>{labels.pageSize}<select value={state.pageSize} onChange={event => patch({ pageSize: event.target.value === '50' ? 50 : 25 })}><option value="25">25</option><option value="50">50</option></select></label>
        </nav>}
        <button className="recording-point-picker__source" data-picker-focus="source" type="button" onClick={() => onSource(state.deviceId || undefined)}>{labels.newSource}</button>
      </>}
    </div>
    <footer role="contentinfo" className="recording-point-picker__footer">
      <div aria-live="polite"><strong>{labels.pending(state.pending.length)}</strong>
        {result.hiddenPending > 0 && <small>{labels.hidden(result.hiddenPending)}</small>}
        {result.invalidPending.length > 0 && <small role="alert">{labels.invalidPending}</small>}
        {result.invalidPending.map(identity => <button
          type="button" key={sourceKey(identity)} onClick={() => onChange({ ...state,
            pending: state.pending.filter(item => sourceKey(item) !== sourceKey(identity)),
          })}>{labels.removePending(identity)}</button>)}
      </div>
      <button type="button" onClick={() => cancel('cancel')}>{labels.cancel}</button>
      <button type="button" disabled={!active || !result.canAdd} onClick={() => {
        const identities = pendingAddition(state, projection);
        if (active && identities) onAdd(identities);
      }}>{labels.add(state.pending.length)}</button>
    </footer>
  </section>;
}
