import { useCallback, useEffect, useId, useRef, useState } from 'react';
import { createPortal } from 'react-dom';
import { Info } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import type { CSSProperties, FocusEvent } from 'react';
import type { Mapping } from '../../state/types';
import type { WorkbenchV2Action } from '../../state/useWorkbenchV2State';
import { MappingSaveFailure } from './MappingSaveFailure';

interface Props {
  pointId: string;
  mapping: Mapping;
  dispatch: React.Dispatch<WorkbenchV2Action>;
}

/** Compact save feedback; the portal keeps details outside the scrolling table. */
export function MappingSaveInfo({ pointId, mapping, dispatch }: Props) {
  const { t } = useTranslation('workbench-v2');
  const id = useId();
  const trigger = useRef<HTMLButtonElement>(null);
  const panel = useRef<HTMLDivElement>(null);
  const closeTimer = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);
  const pinned = useRef(false);
  const [open, setOpen] = useState(false);
  const [position, setPosition] = useState<CSSProperties>({});
  const state = mapping.save_state;
  const active = !!state && state !== 'idle';
  const label = t(`step3.saveStates.${state}`);
  const color = state === 'save-error' || state === 'draft-invalid'
    ? 'text-rose-400' : state === 'saving' ? 'text-amber-400' : 'text-emerald-400';

  const cancelClose = useCallback(() => { clearTimeout(closeTimer.current); }, []);
  const close = useCallback(() => { cancelClose(); pinned.current = false; setOpen(false); }, [cancelClose]);
  function reveal() {
    cancelClose();
    const rect = trigger.current?.getBoundingClientRect();
    if (!rect) return;
    const width = Math.min(320, window.innerWidth - 16);
    const above = rect.bottom + 200 > window.innerHeight;
    setPosition({ width, left: Math.max(8, Math.min(rect.right - width, window.innerWidth - width - 8)),
      ...(above ? { bottom: window.innerHeight - rect.top + 4 } : { top: rect.bottom + 4 }),
      maxHeight: Math.max(80, Math.min(280, above ? rect.top - 12 : window.innerHeight - rect.bottom - 12)),
    });
    setOpen(true);
  }
  function leave() {
    if (!pinned.current) closeTimer.current = setTimeout(() => setOpen(false), 120);
  }
  function blur(event: FocusEvent) {
    const next = event.relatedTarget as Node | null;
    if (!trigger.current?.contains(next) && !panel.current?.contains(next)) close();
  }

  useEffect(() => () => clearTimeout(closeTimer.current), []);
  useEffect(() => {
    if (!open || !active) return;
    const outside = (event: PointerEvent) => {
      const target = event.target as Node;
      if (!trigger.current?.contains(target) && !panel.current?.contains(target)) close();
    };
    const escape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') { event.stopPropagation(); close(); }
    };
    const scroll = (event: Event) => { if (!panel.current?.contains(event.target as Node)) close(); };
    document.addEventListener('pointerdown', outside);
    document.addEventListener('keydown', escape);
    window.addEventListener('scroll', scroll, true);
    window.addEventListener('resize', close);
    return () => {
      document.removeEventListener('pointerdown', outside);
      document.removeEventListener('keydown', escape);
      window.removeEventListener('scroll', scroll, true);
      window.removeEventListener('resize', close);
    };
  }, [open, active, close]);

  if (!active) return null;
  return <>
    <button type="button" ref={trigger} aria-label={label} aria-haspopup="dialog"
      aria-expanded={open} aria-controls={open ? id : undefined}
      data-testid={`mapping-save-info-${pointId}`}
      className={`shrink-0 rounded p-1 focus-visible:outline focus-visible:outline-2 focus-visible:outline-blue-400 ${color}`}
      onMouseEnter={reveal} onMouseLeave={leave}
      onFocus={event => { event.stopPropagation(); reveal(); }} onBlur={blur}
      onClick={event => { event.stopPropagation(); if (pinned.current) close(); else { pinned.current = true; reveal(); } }}>
      <Info size={16} aria-hidden="true" />
    </button>
    {open && createPortal(<div id={id} ref={panel} role="dialog" aria-label={label}
      data-testid={`mapping-save-state-${pointId}`} style={position}
      className={`fixed z-50 overflow-auto break-words rounded-md border border-slate-700 bg-slate-950 p-3 text-xs shadow-lg ${color}`}
      onMouseEnter={cancelClose} onMouseLeave={leave} onBlur={blur}
      onClick={event => event.stopPropagation()} onFocus={event => event.stopPropagation()}>
      {label}
      {state === 'save-error' && <MappingSaveFailure pointId={pointId} mapping={mapping} dispatch={dispatch} />}
    </div>, document.body)}
  </>;
}
