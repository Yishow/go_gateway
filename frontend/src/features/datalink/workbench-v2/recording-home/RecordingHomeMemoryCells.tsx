import type { CSSProperties, KeyboardEvent } from 'react';
import { useTranslation } from 'react-i18next';
import type { MemoryCell, MemoryOwner } from './recordingHomeMemoryMap';
import { memoryValueKey, type MemoryLiveValue } from './recordingHomeMemoryValues';

/** Family is presentation only: the exact persisted type remains visible and in every cell label. */
export function memoryTypeFamily(type: string) {
  return type.startsWith('uint') ? 'unsigned' : type.startsWith('int') ? 'signed'
    : type.startsWith('float') ? 'float' : type === 'bool' ? 'bool' : 'text';
}
function sameValue(cell: MemoryCell | undefined, owner: MemoryOwner) {
  const next = cell?.owners[0];
  return cell?.state === 'used' && !next?.partial && next?.ruleId === owner.ruleId && next.valueAddress === owner.valueAddress;
}

/** Each address remains a separate focus target; a visual label spans only contiguous physical slots. */
export function RecordingHomeMemoryCells({ cells, columns, areaKey, active, pending, describe, reading, onFocus, onHover, onKeyDown, onChoose }: {
  cells: MemoryCell[]; columns: number; areaKey: string; active?: number; pending: boolean;
  describe: (cell: MemoryCell) => string; reading: (cell: MemoryCell) => MemoryLiveValue | undefined;
  onFocus: (offset: number) => void; onHover: (offset: number) => void; onKeyDown: (event: KeyboardEvent, cell: MemoryCell) => void; onChoose: (cell: MemoryCell) => void;
}) {
  const { t } = useTranslation('recording-home');
  return cells.map((cell, index) => {
    const owner = cell.owners[0], value = cell.state === 'used' && owner.width > 1 && !owner.partial ? owner : undefined;
    const column = index % columns;
    const joinedBefore = Boolean(value && column > 0 && sameValue(cells[index - 1], value));
    const joinNext = Boolean(value && value.part < value.width && column < columns - 1 && sameValue(cells[index + 1], value));
    const continuesBelow = Boolean(value && value.part < value.width && (column === columns - 1 || index === cells.length - 1));
    let span = 1;
    if (value && !joinedBefore) {
      while (span < value.width - value.part + 1 && column + span < columns && sameValue(cells[index + span], value)) span++;
    }
    // Focus may scroll another cell under a stationary pointer and emit mouseenter.
    // Only actual mouse movement should replace the keyboard-selected detail.
    const props = { 'data-testid': `home-memory-cell-${cell.address}`, 'data-state': cell.state, 'data-offset': cell.offset, 'data-selected': cell.offset === active || undefined, 'aria-label': describe(cell),
      tabIndex: cell.offset === active ? 0 : -1, onFocus: () => onFocus(cell.offset), onMouseMove: () => onHover(cell.offset), onKeyDown: (event: KeyboardEvent) => onKeyDown(event, cell),
      'data-part': owner ? (owner.part === 1 ? 'first' : 'rest') : undefined,
      'data-family': cell.state === 'used' ? memoryTypeFamily(owner.dataType) : undefined,
      'data-value-key': cell.state === 'used' && owner.part === 1 ? memoryValueKey(areaKey, cell.offset, owner.dataType) : undefined,
      'data-join-next': joinNext || undefined, 'data-join-before': joinedBefore || undefined, 'data-multi': value ? true : undefined, title: describe(cell) };
    const content = <><span className="gw-memory-address">{cell.address}</span>
      {cell.state === 'unconfirmed' && <span className="gw-memory-value" aria-hidden="true">?</span>}
      {cell.state === 'conflict' && <><span className="gw-memory-value" aria-hidden="true">×{cell.owners.length}</span><span className="gw-memory-type" aria-hidden="true">{t('source.memory.state.conflict')}</span></>}
      {cell.state === 'used' && !joinedBefore && <span className="gw-memory-shared" data-shared-span={span} style={{ '--gw-memory-span': span } as CSSProperties} aria-hidden="true">
        <span className="gw-memory-type"><span>{owner.dataType}</span>{value && <span> {value.part}/{value.width}</span>}{owner.partial && <span> ?</span>}{owner.skipped && <span> {t('source.memory.cell.skippedShort')}</span>}</span>
        {owner.part === 1 ? <MemoryReading live={reading(cell)} pending={pending} /> : <span className="gw-memory-value gw-memory-continuation">{t('source.memory.cell.continued', { start: owner.valueAddress })}</span>}
      </span>}
      {/* Confirmed owner.part identifies the canonical start even outside this page. */}
      <span className="gw-memory-focus-information" aria-hidden="true"
        data-value-key={cell.state === 'used' && !owner.partial ? memoryValueKey(areaKey, cell.offset - owner.part + 1, owner.dataType) : undefined}>{describe(cell)}</span>
      {continuesBelow && <span className="gw-memory-cue" aria-hidden="true">{t('source.memory.cell.continues')}</span>}
    </>;
    return cell.state === 'unconfirmed' ? <div key={cell.offset} {...props}>{content}</div>
      : <button key={cell.offset} type="button" {...props} onClick={() => onChoose(cell)}>{content}</button>;
  });
}

function MemoryReading({ live, pending }: { live?: MemoryLiveValue; pending: boolean }) {
  const state = live?.error ? 'error' : live?.value !== undefined ? (live.stale ? 'stale' : 'value') : pending ? 'pending' : 'none';
  return <span className="gw-memory-value gw-memory-reading" data-reading={state}>{state === 'value' || state === 'stale' ? live!.value : state === 'error' ? '!' : state === 'pending' ? '…' : '—'}</span>;
}
