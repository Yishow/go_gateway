import { useEffect, useLayoutEffect, useMemo, useRef, useState, type KeyboardEvent } from 'react';
import { useQueries } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import type { ProtocolType, SourceRuleRecord } from '../../../../types/datalink';
import { sourceRuleCandidateAPI } from '../../../../services/sourceRuleCandidates';
import { buildPersistedMCRuleState } from '../state/mcPersistedPoints';
import './recordingHomeMemory.css';
import { buildMemoryMap, initialMemoryWindow, locateMemoryOffset, memoryAddress, memoryRowStart, memoryRuleInputs, memoryWindow,
  type MemoryArea, type MemoryCell, type MemoryRuleInput } from './recordingHomeMemoryMap';

export interface HomeRangeRead { rules: SourceRuleRecord[]; workspaceRevision?: string }
/** MC proof is read for exactly the list's workspace revision and rule revisions;
 * a new list read (or revision) uses new keys, so older proof never carries over. */
export function useHomeRangeInputs(read: HomeRangeRead | undefined, readAt: number, protocol: ProtocolType | undefined, workspaceId: string) {
  const revision = read?.workspaceRevision;
  const mc = protocol === 'mc_3e' && revision ? read.rules.filter(rule => rule.revision_id) : [];
  const queries = useQueries({ queries: mc.map(rule => ({
    queryKey: ['recording-home-mc-proof', workspaceId, revision, readAt, rule.id, rule.revision_id],
    queryFn: () => sourceRuleCandidateAPI.getTags(rule.id, { workspace_id: workspaceId, expected_workspace_revision: revision!, revision_id: rule.revision_id! }),
    retry: false, refetchOnWindowFocus: false, gcTime: 0, staleTime: Infinity,
  })) });
  const evidence = new Map(mc.map((rule, index) => [rule.id, queries[index]?.data ? buildPersistedMCRuleState(rule,
    { workspace_id: workspaceId, workspace_revision: revision!, revision_id: rule.revision_id! }, queries[index].data) : undefined]));
  const inputs = read && protocol ? memoryRuleInputs(read.rules, protocol, { id: workspaceId, revision }, id => evidence.get(id)) : [];
  return { inputs, checking: queries.some(query => query.isPending), listScopeMissing: protocol === 'mc_3e' && Boolean(read) && !revision };
}

const AREA_NAMES: Record<string, string> = { HR: '4x', IR: '3x', CS: '0x', IS: '1x' };
export type HomeMemoryBlock = 'recovery' | 'local' | 'stale';
export function RecordingHomeMemoryView({ protocol, inputs, status, checking, blocked, listScopeMissing, onOpen, onAdd, onBack, onRetry }: {
  protocol?: ProtocolType; inputs: MemoryRuleInput[]; status: 'loading' | 'error' | 'ready'; checking: boolean; blocked?: HomeMemoryBlock; listScopeMissing?: boolean;
  onOpen: (ruleId: string) => void; onAdd: (address: string) => void; onBack: () => void; onRetry?: () => void;
}) {
  const { t } = useTranslation('recording-home');
  const key = JSON.stringify([protocol, inputs, listScopeMissing]);
  const map = useMemo(() => protocol ? buildMemoryMap(protocol, inputs, { listScopeMissing }) : undefined, [key]); // eslint-disable-line react-hooks/exhaustive-deps
  const retry = onRetry && <button type="button" onClick={onRetry}>{t('source.memory.retry')}</button>;
  const [areaKey, setAreaKey] = useState<string>();
  const [starts, setStarts] = useState<Record<string, number>>({});
  const [focus, setFocus] = useState<number>();
  const [owners, setOwners] = useState<MemoryCell>();
  const [locate, setLocate] = useState(''), [locateError, setLocateError] = useState(false);
  const heading = useRef<HTMLHeadingElement>(null), grid = useRef<HTMLDivElement>(null);
  const [visualColumns, setVisualColumns] = useState<number>();
  useEffect(() => { heading.current?.focus(); }, []);
  const area = map?.areas.find(item => item.key === areaKey) ?? map?.areas.find(item => item.ruleIds.length) ?? map?.areas[0];
  const name = (item: MemoryArea) => AREA_NAMES[item.key] ?? item.key;
  const unit = (item: MemoryArea, count?: number) => t(`source.memory.${item.unit}`, { count });
  const label = (rule: string) => inputs.find(item => item.id === rule)?.label ?? rule;
  // The rendered column count (desktop/tablet/phone) decides where a multi-register value wraps.
  useLayoutEffect(() => {
    const node = grid.current;
    if (!node) return;
    const read = () => { const tracks = getComputedStyle(node).gridTemplateColumns.split(' ').filter(track => track && track !== 'none').length; setVisualColumns(tracks > 1 ? tracks : undefined); };
    read();
    const observer = typeof ResizeObserver === 'undefined' ? undefined : new ResizeObserver(read);
    observer?.observe(node);
    return () => observer?.disconnect();
  }, [status, area?.key, map]);
  const body = () => {
    if (status === 'loading') return <p role="status">{t('loading')}</p>;
    if (status === 'error') return <div role="alert" className="gw-memory-notice"><p>{t('source.memory.readFailed')}{retry}</p></div>;
    if (!map?.supported || !area) return <p role="alert">{t('source.memory.unsupported')}</p>;
    const window = initialMemoryWindow(area), start = starts[area.key] ?? window.start;
    const cells = memoryWindow(area, start, window.length, Boolean(blocked) || checking);
    const selectArea = (key: string) => { setAreaKey(key); setOwners(undefined); setLocateError(false); };
    const issues = [...map.global, ...area.issues.filter(item => !map.global.includes(item))];
    const move = (next: number) => { setStarts(current => ({ ...current, [area.key]: Math.max(area.base, Math.min(next, memoryRowStart(area, area.max))) })); setOwners(undefined); };
    const go = () => { const offset = locateMemoryOffset(area, locate); setLocateError(offset === undefined); if (offset !== undefined) { move(memoryRowStart(area, offset)); setFocus(offset); } };
    const choose = (cell: MemoryCell) => {
      if (cell.state === 'empty') onAdd(cell.address);
      else if (cell.state === 'used') onOpen(cell.owners[0].ruleId);
      else if (cell.state === 'conflict') setOwners(cell);
    };
    const describe = (cell: MemoryCell) => {
      if (cell.state === 'empty') return t('source.memory.cell.empty', { address: cell.address });
      if (cell.state === 'unconfirmed') return t('source.memory.cell.unconfirmed', { address: cell.address });
      if (cell.state === 'conflict') return t('source.memory.cell.conflict', { address: cell.address, count: cell.owners.length, labels: cell.owners.map(owner => owner.label).join(', ') });
      const owner = cell.owners[0];
      return [t('source.memory.cell.used', { address: cell.address, label: owner.label, type: owner.dataType, unit: t(`source.memory.unitName.${area.unit}`), part: owner.part, width: owner.width }),
        owner.part > 1 && t('source.memory.cell.startsAt', { start: owner.valueAddress }), owner.skipped && t('source.memory.cell.skipped'), owner.partial && t('source.memory.cell.partial')].filter(Boolean).join('');
    };
    // The flat grid wraps at the rendered track count (e.g. hex W shows 10 per row on tablet),
    // so keyboard rows and wrap cues use it; the logical radix count is only a fallback.
    const columns = visualColumns ?? area.columns;
    const keyDown = (event: KeyboardEvent, cell: MemoryCell) => {
      const delta = { ArrowLeft: -1, ArrowRight: 1, ArrowUp: -columns, ArrowDown: columns }[event.key];
      if (!delta) return;
      event.preventDefault();
      const target = cells.find(item => item.offset === cell.offset + delta);
      if (target) { setFocus(target.offset); grid.current?.querySelector<HTMLElement>(`[data-offset="${target.offset}"]`)?.focus(); }
    };
    const active = cells.some(cell => cell.offset === focus) ? focus : cells[0]?.offset;
    const messages = [...issues.map(issue => ({ key: `${issue.ruleId}-${issue.issue}`, text: t(`source.memory.issue.${issue.issue}`, { label: issue.label }), ruleId: issue.ruleId || undefined, retry: issue.issue === 'list-scope' })),
      ...blocked ? [{ key: blocked, text: t(`source.memory.blocked.${blocked}`), ruleId: undefined, retry: false }] : []];
    return <>
      <div role="tablist" aria-label={t('source.memory.areas')} className="gw-memory-tabs">
        {map.areas.map((item, index) => <button key={item.key} id={`home-memory-tab-${item.key}`} type="button" role="tab" aria-selected={item.key === area.key} aria-controls="home-memory-panel"
          tabIndex={item.key === area.key ? 0 : -1} data-testid={`home-memory-area-${item.key}`} onClick={() => selectArea(item.key)}
          onKeyDown={event => {
            const count = map.areas.length, next = { ArrowRight: index + 1, ArrowLeft: index - 1, Home: 0, End: count - 1 }[event.key];
            if (next === undefined) return;
            event.preventDefault();
            const target = map.areas[(next + count) % count];
            selectArea(target.key); document.getElementById(`home-memory-tab-${target.key}`)?.focus();
          }}>
          {name(item)} <span>· {unit(item)} · {t('source.memory.rangeCount', { count: item.ruleIds.length })}</span></button>)}
      </div>
      <div role="tabpanel" id="home-memory-panel" aria-labelledby={`home-memory-tab-${area.key}`} className="gw-memory-panel">
        {messages.length > 0 && <div role="alert" className="gw-memory-notice">{messages.map(message => <p key={message.key}>{message.text}
          {message.ruleId && <button type="button" onClick={() => onOpen(message.ruleId!)}>{t('source.memory.openRange', { label: label(message.ruleId) })}</button>}
          {message.retry && retry}</p>)}</div>}
        {checking && <p role="status">{t('source.memory.checking')}</p>}
        <div className="gw-memory-toolbar">
          <button type="button" disabled={start <= area.base} onClick={() => move(start - window.length)}>{t('source.memory.previous')}</button>
          <form onSubmit={event => { event.preventDefault(); go(); }}><label>{t('source.memory.locate')}<input value={locate} placeholder={memoryAddress(area, area.base)}
            aria-invalid={locateError} onChange={event => { setLocate(event.target.value); setLocateError(false); }} /></label><button type="submit">{t('source.memory.go')}</button></form>
          <button type="button" disabled={(cells.at(-1)?.offset ?? area.max) >= area.max} onClick={() => move(start + window.length)}>{t('source.memory.next')}</button>
        </div>
        {locateError && <p role="alert">{t('source.memory.locateInvalid', { area: name(area) })}</p>}
        {area.ruleIds.length > 0 && <p className="gw-memory-jumps"><span>{t('source.memory.jumpTo')}</span>{area.ruleIds.map(id => {
          const rule = inputs.find(item => item.id === id), offset = rule && locateMemoryOffset(area, rule.start_address);
          return offset === undefined ? null : <button key={id} type="button" onClick={() => { move(memoryRowStart(area, offset)); setFocus(offset); }}>{rule!.label} · {memoryAddress(area, offset)}</button>;
        })}</p>}
        <div className="gw-memory-meta"><p className="gw-muted" aria-live="polite">{t('source.memory.showing', { from: cells[0]?.address, to: cells.at(-1)?.address, unit: unit(area) })}</p>
          <ul className="gw-memory-legend" aria-label={t('source.memory.legend')}>{(['empty', 'used', 'conflict', 'unconfirmed'] as const).map(state =>
            <li key={state}><span className="gw-memory-swatch" data-state={state} aria-hidden="true" />{t(`source.memory.state.${state}`)}</li>)}</ul></div>
        {owners && <section className="gw-memory-owners" data-testid="home-memory-owners" aria-labelledby="home-memory-owners-title">
          <h4 id="home-memory-owners-title" tabIndex={-1} ref={node => node?.focus()}>{t('source.memory.ownersTitle', { address: owners.address })}</h4>
          {owners.owners.map(owner => <button key={owner.ruleId} type="button" onClick={() => onOpen(owner.ruleId)}>{t('source.memory.owner', { label: owner.label, type: owner.dataType, start: owner.valueAddress })}</button>)}
          <button type="button" onClick={() => { setOwners(undefined); grid.current?.querySelector<HTMLElement>(`[data-offset="${owners.offset}"]`)?.focus(); }}>{t('source.memory.ownersCancel')}</button>
        </section>}
        <div ref={grid} role="group" aria-label={t('source.memory.gridLabel', { area: name(area), from: cells[0]?.address, to: cells.at(-1)?.address })} className="gw-memory-grid"
          style={{ ['--gw-memory-cols' as string]: area.columns, ['--gw-memory-cols-md' as string]: Math.min(area.columns, 10), ['--gw-memory-cols-sm' as string]: area.columns === 10 ? 5 : 4 }}>
          {cells.map((cell, index) => {
            const owner = cell.owners[0], value = cell.state === 'used' && owner.width > 1 && !owner.partial ? owner : undefined, column = index % columns;
            const next = cells[index + 1]?.owners[0], joinNext = Boolean(value && value.part < value.width && column < columns - 1 && cells[index + 1]?.state === 'used' && next?.ruleId === value.ruleId && next.valueAddress === value.valueAddress);
            const continuesBelow = Boolean(value && value.part < value.width && column === columns - 1), continued = Boolean(value && value.part > 1 && column === 0);
            const props = { 'data-testid': `home-memory-cell-${cell.address}`, 'data-state': cell.state, 'data-offset': cell.offset, 'aria-label': describe(cell),
              tabIndex: cell.offset === active ? 0 : -1, onFocus: () => setFocus(cell.offset), onKeyDown: (event: KeyboardEvent) => keyDown(event, cell),
              'data-part': owner ? (owner.part === 1 ? 'first' : 'rest') : undefined, 'data-join-next': joinNext || undefined, 'data-multi': value ? true : undefined, title: describe(cell) };
            const content = <><span className="gw-memory-address">{cell.address}</span>
              {cell.state === 'conflict' && <><span className="gw-memory-value" aria-hidden="true">×{cell.owners.length}</span><span className="gw-memory-type" aria-hidden="true">{t('source.memory.state.conflict')}</span></>}
              {cell.state === 'used' && <><span className="gw-memory-value" aria-hidden="true">{owner.part === 1 ? owner.label : continued ? t('source.memory.cell.continued', { start: owner.valueAddress }) : `↳ ${owner.valueAddress}`}</span>
                <span className="gw-memory-type" aria-hidden="true"><span>{owner.dataType}</span>{value && <span> {value.part}/{value.width}</span>}{owner.partial && <span> ?</span>}{owner.skipped && <span> {t('source.memory.cell.skippedShort')}</span>}</span>
                {continuesBelow && <span className="gw-memory-cue" aria-hidden="true">{t('source.memory.cell.continues')}</span>}</>}</>;
            return cell.state === 'unconfirmed' ? <div key={cell.offset} {...props}>{content}</div>
              : <button key={cell.offset} type="button" {...props} onClick={() => choose(cell)}>{content}</button>;
          })}
        </div>
      </div>
    </>;
  };
  return <><div className="gw-modal-body" data-testid="home-memory-view"><section className="gw-home-memory" aria-labelledby="home-memory-heading">
    <h3 id="home-memory-heading" ref={heading} tabIndex={-1}>{t('source.memory.heading')}</h3>
    <p className="gw-muted gw-memory-intro">{t('source.memory.intro')}</p>
    <p className="gw-muted gw-memory-intro-short">{t('source.memory.introShort')}</p>
    {body()}
  </section></div><footer className="gw-modal-footer"><button type="button" data-testid="home-memory-back" onClick={onBack}>{t('source.backToList')}</button></footer></>;
}
