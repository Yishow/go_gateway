import { useEffect, useLayoutEffect, useMemo, useRef, useState, type KeyboardEvent } from 'react';
import { useQueries } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import type { ProtocolType, SourceRuleRecord } from '../../../../types/datalink';
import { sourceRuleCandidateAPI } from '../../../../services/sourceRuleCandidates';
import { buildPersistedMCRuleState } from '../state/mcPersistedPoints';
import './recordingHomeMemory.css';
import './recordingHomeMemoryLayouts.css';
import { MemoryRangeNavigation, MemoryRangeChapters, MemoryAddressMap, MemoryCellDetail, MemoryCellSummary } from './RecordingHomeMemoryPanels';
import { RecordingHomeMemoryCells, memoryTypeFamily } from './RecordingHomeMemoryCells';
import { buildMemoryMap, initialMemoryWindow, locateMemoryOffset, memoryAddress, memoryRowStart, memoryRuleInputs, memoryWindow,
  type MemoryArea, type MemoryCell, type MemoryRuleInput } from './recordingHomeMemoryMap';
import { memoryValueKey, type MemoryLiveValue, type MemoryValues } from './recordingHomeMemoryValues';

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

const MEMORY_LAYOUTS = ['grid', 'navigator', 'map', 'chapters', 'detail'] as const;
export type HomeMemoryLayout = typeof MEMORY_LAYOUTS[number];
const AREA_NAMES: Record<string, string> = { HR: '4x', IR: '3x', CS: '0x', IS: '1x' };
export type HomeMemoryBlock = 'recovery' | 'local' | 'stale';
export function RecordingHomeMemoryView({ protocol, inputs, status, checking, blocked, listScopeMissing, values, onOpen, onAdd, onBack, onClose, onRetry }: {
  protocol?: ProtocolType; inputs: MemoryRuleInput[]; status: 'loading' | 'error' | 'ready'; checking: boolean; blocked?: HomeMemoryBlock; listScopeMissing?: boolean;
  /** The gateway's stored device readings; cells show these instead of point names. */
  values?: MemoryValues;
  onOpen: (ruleId: string) => void; onAdd: (address: string) => void; onBack: () => void; onClose: () => void; onRetry?: () => void;
}) {
  const { t } = useTranslation('recording-home');
  const key = JSON.stringify([protocol, inputs, listScopeMissing]);
  const map = useMemo(() => protocol ? buildMemoryMap(protocol, inputs, { listScopeMissing }) : undefined, [key]); // eslint-disable-line react-hooks/exhaustive-deps
  const retry = onRetry && <button type="button" onClick={onRetry}>{t('source.memory.retry')}</button>;
  const [layout, setLayout] = useState<HomeMemoryLayout>('grid');
  const [areaKey, setAreaKey] = useState<string>();
  const [starts, setStarts] = useState<Record<string, number>>({});
  const [focus, setFocus] = useState<number>();
  // Pointer inspection must not move the keyboard tab stop or trigger a cell action.
  const [inspected, setInspected] = useState<number>();
  const select = (offset: number) => { setInspected(undefined); setFocus(offset); };
  const [owners, setOwners] = useState<MemoryCell>();
  const [locate, setLocate] = useState(''), [locateError, setLocateError] = useState(false);
  const viewport = useRef<HTMLDivElement>(null);
  const heading = useRef<HTMLHeadingElement>(null), grid = useRef<HTMLDivElement>(null), tabs = useRef<HTMLDivElement>(null);
  const [tabEdges, setTabEdges] = useState({ before: false, after: false });
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
    const read = () => { const tracks = getComputedStyle(node).gridTemplateColumns.split(' ').filter(track => track && track !== 'none').length; setVisualColumns(tracks > 0 ? tracks : undefined); };
    read();
    const observer = typeof ResizeObserver === 'undefined' ? undefined : new ResizeObserver(read);
    observer?.observe(node);
    return () => observer?.disconnect();
  }, [status, checking, area?.key, map, layout]);
  // Whatever selected it (key, pointer or touch), the selected tab is scrolled fully into the tab row.
  const selectedTab = area?.key;
  useLayoutEffect(() => {
    if (selectedTab) document.getElementById(`home-memory-tab-${selectedTab}`)?.scrollIntoView?.({ block: 'nearest', inline: 'nearest' });
  }, [selectedTab, checking]);
  // Area tabs stay on one row; mark which side still has tabs so the edge can show it.
  useLayoutEffect(() => {
    const node = tabs.current;
    if (!node) return;
    const read = () => setTabEdges(current => {
      const next = { before: node.scrollLeft > 1, after: node.scrollLeft + node.clientWidth < node.scrollWidth - 1 };
      return current.before === next.before && current.after === next.after ? current : next;
    });
    read(); node.addEventListener('scroll', read, { passive: true });
    const observer = typeof ResizeObserver === 'undefined' ? undefined : new ResizeObserver(read);
    observer?.observe(node);
    return () => { node.removeEventListener('scroll', read); observer?.disconnect(); };
  }, [status, checking, map, layout]);
  // One subscription covers visible grid and detail readings; duplicate displays share a key.
  const reportVisible = values?.setVisible;
  useEffect(() => {
    const node = viewport.current;
    if (!node || !reportVisible) return undefined;
    const targets = [...node.querySelectorAll<HTMLElement>('[data-value-key]')];
    // Loading/error/checking replaces the cells: do not keep their old subscription.
    if (!targets.length) { reportVisible([]); return undefined; }
    if (typeof IntersectionObserver === 'undefined') { reportVisible([...new Set(targets.map(cell => cell.dataset.valueKey!))]); return undefined; }
    let observing = true;
    const shown = new Map<Element, string>();
    const observer = new IntersectionObserver(entries => {
      if (!observing) return;
      entries.forEach(entry => {
        const key = (entry.target as HTMLElement).dataset.valueKey;
        if (key && entry.isIntersecting) shown.set(entry.target, key); else shown.delete(entry.target);
      });
      reportVisible([...new Set(shown.values())]);
    });
    targets.forEach(cell => observer.observe(cell));
    // Do not clear here: populated hover/layout changes must not briefly unsubscribe.
    return () => { observing = false; observer.disconnect(); };
  }, [reportVisible, area?.key, starts, map, status, checking, blocked, layout, focus, inspected]);
  const body = () => {
    if (status === 'loading') return <p role="status">{t('loading')}</p>;
    if (status === 'error') return <div role="alert" className="gw-memory-notice"><p>{t('source.memory.readFailed')}{retry}</p></div>;
    // Missing proof is not an error while its scoped read is still pending.
    if (checking) return <>{blocked && <div role="alert" className="gw-memory-notice"><p>{t(`source.memory.blocked.${blocked}`)}</p></div>}
      <p role="status">{t('source.memory.checking')}</p></>;
    if (!map?.supported || !area) return <p role="alert">{t('source.memory.unsupported')}</p>;
    const window = initialMemoryWindow(area), start = starts[area.key] ?? window.start;
    const cells = memoryWindow(area, start, window.length, Boolean(blocked) || checking);
    const selectArea = (key: string) => { if (key !== area.key) setFocus(undefined); setAreaKey(key); setInspected(undefined); setOwners(undefined); setLocateError(false); };
    const issues = [...map.global, ...area.issues.filter(item => !map.global.includes(item))];
    const move = (next: number) => { setInspected(undefined); setStarts(current => ({ ...current, [area.key]: Math.max(area.base, Math.min(next, memoryRowStart(area, area.max))) })); setOwners(undefined); };
    const go = () => { const offset = locateMemoryOffset(area, locate); setLocateError(offset === undefined); if (offset !== undefined) { move(memoryRowStart(area, offset)); select(offset); } };
    const choose = (cell: MemoryCell) => {
      if (cell.state === 'empty') onAdd(cell.address);
      else if (cell.state === 'used') onOpen(cell.owners[0].ruleId);
      else if (cell.state === 'conflict') setOwners(cell);
    };
    const reading = (cell: MemoryCell): MemoryLiveValue | undefined => {
      const owner = cell.owners[0];
      if (cell.state !== 'used' || !owner || owner.partial || owner.part !== 1) return undefined;
      // A value's first register is its own address.
      return values?.byCell.get(memoryValueKey(area.key, cell.offset, owner.dataType));
    };
    const describe = (cell: MemoryCell) => {
      if (cell.state === 'empty') return t('source.memory.cell.empty', { address: cell.address });
      if (cell.state === 'unconfirmed') return t('source.memory.cell.unconfirmed', { address: cell.address });
      if (cell.state === 'conflict') return t('source.memory.cell.conflict', { address: cell.address, count: cell.owners.length, labels: cell.owners.map(owner => owner.label).join(', ') });
      const owner = cell.owners[0];
      const valueStart = locateMemoryOffset(area, owner.valueAddress);
      const live = !owner.partial && valueStart !== undefined ? values?.byCell.get(memoryValueKey(area.key, valueStart, owner.dataType)) : undefined;
      const freshness = live?.error ? 'failed' : live?.value !== undefined ? live.stale ? 'stale' : 'available'
        : values?.status === 'loading' ? 'waiting' : 'missing';
      return [t('source.memory.cell.used', { address: cell.address, label: owner.label, type: owner.dataType, unit: t(`source.memory.unitName.${area.unit}`), part: owner.part, width: owner.width }),
        owner.part > 1 && t('source.memory.cell.startsAt', { start: owner.valueAddress }), owner.skipped && t('source.memory.cell.skipped'), owner.partial && t('source.memory.cell.partial'),
        live && (live.error ? t('source.memory.value.cellError') : live.value !== undefined ? t('source.memory.value.cell', { value: live.value }) : ''),
        !owner.partial && `, ${t('source.memory.presentation.reading.' + freshness)}`,
        live?.readAt && t('source.memory.value.readAt', { time: clock(live.readAt) })].filter(Boolean).join('');
    };
    // The flat grid wraps at the rendered track count (e.g. hex W shows 10 per row on tablet),
    // so keyboard rows and wrap cues use it; the logical radix count is only a fallback.
    const columns = visualColumns ?? area.columns;
    const keyDown = (event: KeyboardEvent, cell: MemoryCell) => {
      const delta = { ArrowLeft: -1, ArrowRight: 1, ArrowUp: -columns, ArrowDown: columns }[event.key];
      if (!delta) return;
      event.preventDefault();
      const target = cells.find(item => item.offset === cell.offset + delta);
      if (target) { select(target.offset); grid.current?.querySelector<HTMLElement>(`[data-offset="${target.offset}"]`)?.focus(); }
    };
    const active = cells.some(cell => cell.offset === focus) ? focus : cells[0]?.offset;
    const selected = cells.find(cell => cell.offset === inspected) ?? cells.find(cell => cell.offset === active);
    const jump = (address: string) => { const offset = locateMemoryOffset(area, address); if (offset !== undefined) { move(memoryRowStart(area, offset)); select(offset); } };
    const types = [...new Set(cells.flatMap(cell => cell.owners.map(owner => owner.dataType)))];
    const messages = [...issues.map(issue => ({ key: `${issue.ruleId}-${issue.issue}`, text: t(`source.memory.issue.${issue.issue}`, { label: issue.label }), ruleId: issue.ruleId || undefined, retry: issue.issue === 'list-scope' })),
      ...blocked ? [{ key: blocked, text: t(`source.memory.blocked.${blocked}`), ruleId: undefined, retry: false }] : []];
    return <>
      <div ref={tabs} role="tablist" aria-label={t('source.memory.areas')} className="gw-memory-tabs" data-more-before={tabEdges.before} data-more-after={tabEdges.after}>
        {map.areas.map((item, index) => <button key={item.key} id={`home-memory-tab-${item.key}`} type="button" role="tab" aria-selected={item.key === area.key} aria-controls="home-memory-panel"
          tabIndex={item.key === area.key ? 0 : -1} data-testid={`home-memory-area-${item.key}`}
          onClick={event => { selectArea(item.key); event.currentTarget.scrollIntoView?.({ block: 'nearest', inline: 'nearest' }); }}
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
        <div className="gw-memory-toolbar">
          <form onSubmit={event => { event.preventDefault(); go(); }}><label><span className="gw-memory-locate-label">{t('source.memory.locate')}</span><input value={locate} placeholder={memoryAddress(area, area.base)}
            aria-invalid={locateError} onChange={event => { setLocate(event.target.value); setLocateError(false); }} /></label><button type="submit">{t('source.memory.go')}</button></form>
        <div className="gw-memory-nav" role="group" aria-label={t('source.memory.navigation')}>
          <button type="button" disabled={start <= area.base} onClick={() => move(start - window.length)}><span aria-hidden="true">‹ </span>{t('source.memory.previous')}</button>
          <p aria-live="polite">{t('source.memory.showing', { from: cells[0]?.address, to: cells.at(-1)?.address, unit: unit(area) })}</p>
          <button type="button" disabled={(cells.at(-1)?.offset ?? area.max) >= area.max} onClick={() => move(start + window.length)}>{t('source.memory.next')}<span aria-hidden="true"> ›</span></button>
        </div>
        </div>
        {locateError && <p role="alert">{t('source.memory.locateInvalid', { area: name(area) })}</p>}
        {area.ruleIds.length > 0 && layout !== 'navigator' && layout !== 'chapters' && <p className="gw-memory-jumps"><span>{t('source.memory.jumpTo')} <span className="gw-muted">{area.ruleIds.length}</span></span>{area.ruleIds.map(id => {
          const rule = inputs.find(item => item.id === id), offset = rule && locateMemoryOffset(area, rule.start_address);
          if (offset === undefined) return null;
          const text = `${memoryAddress(area, offset)} · ${rule!.label}`;
          return <button key={id} type="button" aria-label={text} title={text} onClick={() => { move(memoryRowStart(area, offset)); setFocus(offset); }}>{text}</button>;
        })}</p>}
        <div className="gw-memory-meta">
        {types.length > 0 && <ul className="gw-memory-type-legend" aria-label={t('source.memory.presentation.types')}>
          {types.map(type => <li key={type} data-family={memoryTypeFamily(type)}>{type}</li>)}
        </ul>}
          <ul className="gw-memory-legend" aria-label={t('source.memory.legend')}>{(['empty', 'used', 'conflict', 'unconfirmed'] as const).map(state =>
            <li key={state}><span className="gw-memory-swatch" data-state={state} aria-hidden="true" />{t(`source.memory.state.${state}`)}</li>)}</ul>
          {values && <p className="gw-memory-live" data-status={values.status} role={values.status === 'error' ? 'alert' : 'status'} data-testid="home-memory-values-status">
            <span className="gw-memory-live-dot" aria-hidden="true" />{(values.status === 'live' || values.status === 'polling') && values.updatedAt ? t(`source.memory.value.${values.status}Updated`, { time: clock(values.updatedAt) }) : t(`source.memory.value.${values.status}`)}</p>}</div>
        <div className="gw-memory-workspace" data-layout={layout}>
          {layout === 'navigator' && <MemoryRangeNavigation area={area} inputs={inputs} onJump={jump} cells={cells} selected={cells.find(cell => cell.offset === focus)} />}
          {layout === 'chapters' && <MemoryRangeChapters area={area} inputs={inputs} onJump={jump} />}
          <div className="gw-memory-main">
        {layout === 'map' && <MemoryAddressMap cells={cells} selected={selected?.offset} onSelect={select} />}
        {owners && <section className="gw-memory-owners" data-testid="home-memory-owners" aria-labelledby="home-memory-owners-title">
          <h4 id="home-memory-owners-title" tabIndex={-1} ref={focusOwnerHeading}>{t('source.memory.ownersTitle', { address: owners.address })}</h4>
          {owners.owners.map(owner => <button key={owner.ruleId} type="button" onClick={() => onOpen(owner.ruleId)}>{t('source.memory.owner', { label: owner.label, type: owner.dataType, start: owner.valueAddress })}</button>)}
          <button type="button" onClick={() => { setOwners(undefined); grid.current?.querySelector<HTMLElement>(`[data-offset="${owners.offset}"]`)?.focus(); }}>{t('source.memory.ownersCancel')}</button>
        </section>}
        <div ref={grid} role="group" aria-label={t('source.memory.gridLabel', { area: name(area), from: cells[0]?.address, to: cells.at(-1)?.address })} className="gw-memory-grid">
          <RecordingHomeMemoryCells cells={cells} columns={columns} areaKey={area.key} active={active} expanded={focus} pending={values?.status === 'loading'}
            describe={describe} reading={reading} onFocus={select} onHover={setInspected} onKeyDown={keyDown} onChoose={choose} />
        </div>
        {layout === 'map' && <MemoryCellSummary area={area} cell={selected} values={values} cells={cells} onSelect={select} onChoose={choose} />}
          </div>
          {layout === 'detail' && <MemoryCellDetail area={area} cell={selected} values={values} cells={cells} onSelect={select} onChoose={choose} />}
        </div>
      </div>
    </>;
  };
  return <><div ref={viewport} className="gw-modal-body" data-testid="home-memory-view" data-layout={layout}><section className="gw-home-memory" aria-labelledby="home-memory-heading" aria-describedby="home-memory-description">
    <h3 id="home-memory-heading" ref={heading} tabIndex={-1}>{t('source.memory.heading')}</h3>
    <p id="home-memory-description" className="gw-muted gw-memory-intro">{t('source.memory.intro')}</p>
    <p className="gw-muted gw-memory-intro-short">{t('source.memory.introShort')}</p>
    <div className="gw-memory-layouts" role="group" aria-label={t('source.memory.layouts.label')}>
      {MEMORY_LAYOUTS.map(item => <button key={item} type="button" data-testid={`home-memory-layout-${item}`} aria-pressed={layout === item}
        onClick={() => setLayout(item)}>{t(`source.memory.layouts.${item}`)}</button>)}
    </div>
    {body()}
  </section></div><footer className="gw-modal-footer"><button type="button" data-testid="home-memory-back" onClick={onBack}>{t('source.backToSettings')}</button>
    <button type="button" data-testid="home-memory-close" onClick={onClose}>{t('source.close')}</button></footer></>;
}

// Stable ref focuses a newly opened chooser, never a later reading or pointer update.
function focusOwnerHeading(node: HTMLHeadingElement | null) { node?.focus(); }

function clock(at: string | number): string {
  const date = new Date(at);
  return Number.isNaN(date.getTime()) ? String(at) : date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false });
}
