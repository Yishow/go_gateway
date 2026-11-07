import { useEffect, useId, useLayoutEffect, useRef, useState } from 'react';
import { MoreHorizontal } from 'lucide-react';
import './recordingHomeRows.css';

export interface RecordingHomeRowAction {
  id: string;
  label: string;
  disabled?: boolean;
  onSelect: () => void;
}

/** Shared keyboard-operable actions menu for recording and device surfaces. */
export function RecordingHomeRowActions({ label, items, hint }: { label: string; items: readonly RecordingHomeRowAction[]; hint?: string }) {
  const [open, setOpen] = useState(false);
  const trigger = useRef<HTMLButtonElement>(null), root = useRef<HTMLDivElement>(null), menu = useRef<HTMLDivElement>(null);
  const focusEdge = useRef<'first' | 'last'>('first');
  const menuId = useId();
  const close = (restore: boolean) => { setOpen(false); if (restore) trigger.current?.focus(); };
  useLayoutEffect(() => {
    if (open) {
      const buttons = [...(menu.current?.querySelectorAll<HTMLButtonElement>('[role=menuitem]:not(:disabled)') ?? [])];
      (focusEdge.current === 'last' ? buttons.at(-1) : buttons[0])?.focus();
    }
  }, [open]);
  useEffect(() => {
    if (!open) return;
    const outside = (event: PointerEvent) => { if (!root.current?.contains(event.target as Node)) close(false); };
    document.addEventListener('pointerdown', outside);
    return () => document.removeEventListener('pointerdown', outside);
  }, [open]);
  const navigate = (event: React.KeyboardEvent) => {
    if (event.key === 'Escape') { event.preventDefault(); close(true); return; }
    if (!['ArrowDown', 'ArrowUp', 'Home', 'End'].includes(event.key)) return;
    event.preventDefault();
    const buttons = [...(menu.current?.querySelectorAll<HTMLButtonElement>('[role=menuitem]:not(:disabled)') ?? [])];
    if (!buttons.length) return;
    const current = buttons.indexOf(document.activeElement as HTMLButtonElement);
    const next = event.key === 'Home' ? 0 : event.key === 'End' ? buttons.length - 1 :
      event.key === 'ArrowDown' ? (current + 1) % buttons.length : (current <= 0 ? buttons.length - 1 : current - 1);
    buttons[next].focus();
  };
  return <div ref={root} className="gw-record-row-menu">
    <button ref={trigger} type="button" className="gw-row-icon-button" aria-label={label} aria-haspopup="menu" aria-expanded={open} aria-controls={menuId}
      onClick={() => { focusEdge.current = 'first'; setOpen(value => !value); }} onKeyDown={event => {
        if (open && event.key === 'Escape') { event.preventDefault(); close(true); }
        else if (!open && ['ArrowDown', 'ArrowUp'].includes(event.key)) { event.preventDefault(); focusEdge.current = event.key === 'ArrowUp' ? 'last' : 'first'; setOpen(true); }
      }}>
      <MoreHorizontal aria-hidden="true" size={18} />
    </button>
    {open && <div ref={menu} id={menuId} role="menu" className="gw-record-row-menu-list" onKeyDown={navigate}>
      {items.map(item => <button key={item.id} type="button" role="menuitem" disabled={item.disabled} onClick={() => { close(true); item.onSelect(); }}>{item.label}</button>)}
      {hint && <p className="gw-record-row-menu-hint">{hint}</p>}
    </div>}
  </div>;
}
