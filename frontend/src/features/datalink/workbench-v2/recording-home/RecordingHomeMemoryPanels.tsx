import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { describeMemoryRule, locateMemoryOffset, type MemoryArea, type MemoryCell, type MemoryRuleInput } from './recordingHomeMemoryMap';
import { memoryValueKey, type MemoryValues } from './recordingHomeMemoryValues';
import { memoryTypeFamily } from './RecordingHomeMemoryCells';

interface RangeProps { area: MemoryArea; inputs: readonly MemoryRuleInput[]; onJump: (address: string) => void }
function RangeFacts({ area, rule }: { area: MemoryArea; rule: MemoryRuleInput }) {
  const { t } = useTranslation('recording-home');
  const span = describeMemoryRule(area.protocol, rule);
  return <dl className="gw-memory-range-facts">
    <div><dt>{t('source.memory.presentation.type')}</dt><dd>{rule.data_type}</dd></div>
    <div><dt>{t('source.memory.presentation.count')}</dt><dd>{rule.count}</dd></div>
    <div><dt>{t('source.memory.presentation.start')}</dt><dd>{span.start}</dd></div>
    <div><dt>{t('source.memory.presentation.end')}</dt><dd>{span.end ?? t('source.memory.presentation.uncertain')}</dd></div>
    <div><dt>{t('source.memory.presentation.footprint')}</dt><dd>{span.cells === undefined ? t('source.memory.presentation.uncertain')
      : t(`source.memory.presentation.${area.unit}Count`, { count: span.cells })}</dd></div>
  </dl>;
}

/** Navigation is a projection of the same saved rules, never a second source of ranges. */
export function MemoryRangeNavigation({ area, inputs, onJump, cells, selected }: RangeProps & { cells: readonly MemoryCell[]; selected?: MemoryCell }) {
  const { t } = useTranslation('recording-home');
  const rules = inputs.filter(rule => area.ruleIds.includes(rule.id));
  const selectedRule = selected?.state === 'used' ? selected.owners[0]?.ruleId : undefined;
  const windowRules = new Set(cells.flatMap(cell => cell.owners.map(owner => owner.ruleId)));
  const [expanded, setExpanded] = useState(() => typeof matchMedia === 'undefined' || !matchMedia('(max-width: 639px)').matches);
  return <aside className="gw-memory-range-nav" data-testid="home-memory-range-nav" aria-label={t('source.memory.presentation.rangeNavigation')}>
    <details open={expanded} onToggle={event => setExpanded(event.currentTarget.open)}><summary>{t('source.memory.presentation.savedRanges', { count: rules.length })}</summary>
      {!rules.length && <p>{t('source.memory.presentation.noRanges')}</p>}
      <ul>{rules.map(rule => <li key={rule.id} data-in-window={windowRules.has(rule.id) || undefined}>
        <button type="button" data-testid={`home-memory-range-jump-${rule.id}`} aria-current={selectedRule === rule.id ? true : undefined} onClick={() => onJump(rule.start_address)}>
          <strong>{rule.label}</strong><span>{describeMemoryRule(area.protocol, rule).start} · {rule.data_type}</span>
          {selectedRule === rule.id && <span>{t('source.memory.presentation.selectedRange')}</span>}
        </button>{windowRules.has(rule.id) && <p className="gw-memory-window-membership">{t('source.memory.presentation.inWindow')}</p>}<RangeFacts area={area} rule={rule} />
      </li>)}</ul>
    </details>
  </aside>;
}

/** Native disclosure preserves keyboard operation and does not alter saved range state. */
export function MemoryRangeChapters({ area, inputs, onJump }: RangeProps) {
  const { t } = useTranslation('recording-home');
  const rules = inputs.filter(rule => area.ruleIds.includes(rule.id));
  return <section className="gw-memory-chapters" data-testid="home-memory-chapters" aria-label={t('source.memory.presentation.rangeChapters')}>
    <h4>{t('source.memory.presentation.savedRanges', { count: rules.length })}</h4>
    {!rules.length && <p>{t('source.memory.presentation.noRanges')}</p>}
    <div className="gw-memory-chapter-list">{rules.map(rule => {
      const span = describeMemoryRule(area.protocol, rule);
      return <details key={rule.id} data-testid={`home-memory-chapter-${rule.id}`}>
        <summary><span className="gw-memory-chapter-summary"><strong>{rule.label}</strong><span>{span.start} – {span.end ?? t('source.memory.presentation.uncertain')}</span><span>{rule.data_type}</span>
          <span>{t('source.memory.presentation.valueCount', { count: rule.count })} · {span.cells === undefined ? t('source.memory.presentation.uncertain') : t(`source.memory.presentation.${area.unit}Count`, { count: span.cells })}</span></span></summary>
        <div className="gw-memory-chapter-body"><RangeFacts area={area} rule={rule} />
          <button type="button" data-testid={`home-memory-range-jump-${rule.id}`} onClick={() => onJump(rule.start_address)}>{t('source.memory.presentation.jumpRange')}</button>
        </div>
      </details>;
    })}</div>
  </section>;
}

/** One strip segment per physical cell, with an exact index slider for keyboard and touch. */
export function MemoryAddressMap({ cells, selected, onSelect }: {
  cells: readonly MemoryCell[]; selected?: number; onSelect: (offset: number) => void;
}) {
  const { t } = useTranslation('recording-home');
  const index = Math.max(0, cells.findIndex(cell => cell.offset === selected));
  const current = cells[index];
  return <section className="gw-memory-address-map" data-testid="home-memory-map" aria-label={t('source.memory.presentation.addressMap')}>
    <div className="gw-memory-map-heading"><h4>{t('source.memory.presentation.addressMap')}</h4>
      <span>{current ? t('source.memory.presentation.selectedAddress', { address: current.address }) : t('source.memory.presentation.noSelection')}</span></div>
    <div className="gw-memory-map-strip" aria-hidden="true">{cells.map(cell => <span key={cell.offset} data-state={cell.state}
      data-selected={cell.offset === current?.offset || undefined} data-family={cell.state === 'used' ? memoryTypeFamily(cell.owners[0].dataType) : undefined} />)}</div>
    <input type="range" min={0} max={Math.max(0, cells.length - 1)} step={1} value={index} disabled={!cells.length}
      data-testid="home-memory-map-select" aria-label={t('source.memory.presentation.selectAddress')} aria-valuetext={current?.address}
      onChange={event => { const cell = cells[Number(event.target.value)]; if (cell) onSelect(cell.offset); }} />
    <div className="gw-memory-map-ends"><span>{cells[0]?.address}</span><span>{cells.at(-1)?.address}</span></div>
  </section>;
}

/** Inspect the selected physical slot, resolving a continuation's reading at its value start. */
interface InspectionProps {
  area: MemoryArea; cell?: MemoryCell; values?: MemoryValues; onChoose: (cell: MemoryCell) => void;
  onSelect: (offset: number) => void; cells: readonly MemoryCell[];
}
function inspection({ area, cell, values, cells }: InspectionProps) {
  const owner = cell?.state === 'used' ? cell.owners[0] : undefined;
  const start = owner && locateMemoryOffset(area, owner.valueAddress);
  const live = owner && !owner.partial && start !== undefined ? values?.byCell.get(memoryValueKey(area.key, start, owner.dataType)) : undefined;
  const index = cells.findIndex(item => item.offset === cell?.offset);
  const readingState = live?.error ? 'failed' : live?.value !== undefined ? live.stale ? 'stale' : 'available'
    : values?.status === 'loading' ? 'waiting' : 'missing';
  return { owner, start, live, index, readingState };
}
export function MemoryCellDetail(props: InspectionProps) {
  const { area, cell, values, cells, onSelect, onChoose } = props;
  const { t } = useTranslation('recording-home');
  const { owner, start, live, index, readingState } = inspection(props);
  return <aside className="gw-memory-detail gw-memory-inspection" data-memory-inspection data-testid="home-memory-detail" aria-label={t('source.memory.presentation.cellDetail')}>
    <header><h4>{t('source.memory.presentation.cellDetail')}</h4><div className="gw-memory-detail-navigation">
      <button type="button" data-testid="home-memory-detail-previous" aria-label={t('source.memory.presentation.previousAddress')} disabled={index <= 0}
        onClick={() => { if (index > 0) onSelect(cells[index - 1].offset); }}>‹</button>
      <button type="button" data-testid="home-memory-detail-next" aria-label={t('source.memory.presentation.nextAddress')} disabled={index < 0 || index >= cells.length - 1}
        onClick={() => { if (index >= 0 && index < cells.length - 1) onSelect(cells[index + 1].offset); }}>›</button>
    </div></header>
    {!cell ? <p>{t('source.memory.presentation.noSelection')}</p> : <>
      <p className="gw-memory-detail-address">{cell.address}</p>
      <dl className="gw-memory-detail-facts" data-detail-group="type">
        <div><dt>{t('source.memory.presentation.type')}</dt><dd>{owner?.dataType ?? t('source.memory.presentation.notApplicable')}</dd></div>
        <div><dt>{t('source.memory.presentation.valueStart')}</dt><dd>{owner?.valueAddress ?? t('source.memory.presentation.notApplicable')}</dd></div>
        <div><dt>{t('source.memory.presentation.width')}</dt><dd>{owner ? owner.partial ? t('source.memory.presentation.uncertain')
          : t(`source.memory.presentation.${area.unit}Count`, { count: owner.width }) : t('source.memory.presentation.notApplicable')}</dd></div>
      </dl>
      <dl className="gw-memory-detail-facts" data-detail-group="reading">
        <div><dt>{t('source.memory.presentation.currentReading')}</dt><dd className="gw-memory-detail-reading" data-value-key={owner && !owner.partial && start !== undefined ? memoryValueKey(area.key, start, owner.dataType) : undefined}>{live?.error ? '!' : live?.value ?? (values?.status === 'loading' ? '…' : t('source.memory.presentation.notApplicable'))}</dd></div>
        <div><dt>{t('source.memory.presentation.status')}</dt><dd>{t(`source.memory.presentation.state.${cell.state}`)}
          {owner?.skipped && <span> · {t('source.memory.presentation.skipped')}</span>}</dd></div>
        <div><dt>{t('source.memory.presentation.freshness')}</dt><dd>{t(`source.memory.presentation.reading.${readingState}`)}
          {live?.readAt && <time dateTime={live.readAt}> · {live.readAt}</time>}</dd></div>
      </dl>
      <div data-detail-group="actions">{cell.state !== 'unconfirmed' && <button type="button" data-testid="home-memory-detail-action" onClick={() => onChoose(cell)}>
        {t(`source.memory.presentation.action.${cell.state}`)}</button>}</div>
    </>}
  </aside>;
}

/** Compact map selection, not the full right-hand inspector. Shares canonical reading resolution. */
export function MemoryCellSummary(props: InspectionProps) {
  const { area, cell, values, cells, onSelect, onChoose } = props;
  const { t } = useTranslation('recording-home');
  const { owner, start, live, index, readingState } = inspection(props);
  return <section className="gw-memory-map-summary gw-memory-inspection" data-memory-inspection data-testid="home-memory-summary" aria-label={t('source.memory.presentation.cellDetail')}>
    {!cell ? <p>{t('source.memory.presentation.noSelection')}</p> : <>
      <span className="gw-memory-summary-field"><span>{t('source.memory.presentation.cellDetail')}</span><strong className="gw-memory-detail-address">{cell.address}</strong></span>
      <span className="gw-memory-summary-field"><span>{t('source.memory.presentation.type')}</span><strong>{owner?.dataType ?? t('source.memory.presentation.notApplicable')}</strong></span>
      {owner && <span className="gw-memory-summary-field"><span>{t('source.memory.presentation.valueStart')}</span><strong>{owner.valueAddress}</strong>
        <span>{owner.partial ? t('source.memory.presentation.uncertain') : t(`source.memory.presentation.${area.unit}Count`, { count: owner.width })}</span></span>}
      <span className="gw-memory-summary-field"><span>{t('source.memory.presentation.currentReading')}</span>
        <strong className="gw-memory-detail-reading" data-value-key={owner && !owner.partial && start !== undefined ? memoryValueKey(area.key, start, owner.dataType) : undefined}>
          {live?.error ? '!' : live?.value ?? (values?.status === 'loading' ? '…' : t('source.memory.presentation.notApplicable'))}</strong>
        <span>{t(`source.memory.presentation.reading.${readingState}`)}</span>{live?.readAt && <time dateTime={live.readAt}>{live.readAt}</time>}</span>
      <span className="gw-memory-summary-field"><span>{t('source.memory.presentation.status')}</span><strong>{t(`source.memory.presentation.state.${cell.state}`)}</strong>
        {owner?.skipped && <span>{t('source.memory.presentation.skipped')}</span>}</span>
      <div className="gw-memory-detail-navigation">
        <button type="button" data-testid="home-memory-detail-previous" aria-label={t('source.memory.presentation.previousAddress')} disabled={index <= 0}
          onClick={() => { if (index > 0) onSelect(cells[index - 1].offset); }}>‹</button>
        <button type="button" data-testid="home-memory-detail-next" aria-label={t('source.memory.presentation.nextAddress')} disabled={index < 0 || index >= cells.length - 1}
          onClick={() => { if (index >= 0 && index < cells.length - 1) onSelect(cells[index + 1].offset); }}>›</button>
      </div>
      {cell.state !== 'unconfirmed' && <button type="button" data-testid="home-memory-detail-action" onClick={() => onChoose(cell)}>{t(`source.memory.presentation.action.${cell.state}`)}</button>}
    </>}
  </section>;
}
