import './parameterHelp.css';
import { useCallback, useEffect, useId, useLayoutEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { ParameterHelpKey } from './parameterHelpContent';

/** A nonmodal top-layer note inside the existing focus scope. No business state or commands. */
export function ParameterHelp({ help, compact = false }: { help: ParameterHelpKey; compact?: boolean }) {
  const { t } = useTranslation('parameter-help');
  const id = useId(), trigger = useRef<HTMLButtonElement>(null), panel = useRef<HTMLDivElement>(null);
  const [open, setOpen] = useState(false);
  const [position, setPosition] = useState({ left: 8, top: 8, width: 352, maxHeight: 300 });
  const close = useCallback((restore: boolean) => { panel.current?.hidePopover?.(); setOpen(false); if (restore) trigger.current?.focus({ preventScroll: true }); }, []);
  useLayoutEffect(() => {
    if (!open || !trigger.current || !panel.current) return;
    const note = panel.current;
    note.showPopover?.();
    const place = () => {
      const rect = trigger.current!.getBoundingClientRect(), view = window.visualViewport;
      const width = view?.width ?? window.innerWidth, height = view?.height ?? window.innerHeight;
      const x = view?.offsetLeft ?? 0, y = view?.offsetTop ?? 0;
      const panelWidth = Math.max(0, Math.min(352, width - 16)), panelHeight = Math.min(note.scrollHeight || 240, height - 16);
      const below = rect.bottom + 8, above = rect.top - panelHeight - 8;
      setPosition({ width: panelWidth, left: Math.max(x + 8, Math.min(rect.left, x + width - panelWidth - 8)),
        top: Math.max(y + 8, Math.min(below + panelHeight <= y + height - 8 ? below : above, y + height - panelHeight - 8)), maxHeight: Math.max(80, height - 16) });
    };
    place(); note.focus({ preventScroll: true });
    window.addEventListener('resize', place); window.addEventListener('scroll', place, true);
    window.visualViewport?.addEventListener('resize', place); window.visualViewport?.addEventListener('scroll', place);
    return () => { window.removeEventListener('resize', place); window.removeEventListener('scroll', place, true);
      window.visualViewport?.removeEventListener('resize', place); window.visualViewport?.removeEventListener('scroll', place); };
  }, [open]);
  useEffect(() => {
    if (!open) return;
    // Capture before the existing Dialog's document listener; Escape only dismisses this note.
    const key = (event: KeyboardEvent) => { if (event.key === 'Escape' && !event.isComposing) { event.preventDefault(); event.stopPropagation(); close(true); } };
    const outside = (event: PointerEvent) => { if (event.target instanceof Node && !panel.current?.contains(event.target) && !trigger.current?.contains(event.target)) close(false); };
    window.addEventListener('keydown', key, true); window.addEventListener('pointerdown', outside, true);
    return () => { window.removeEventListener('keydown', key, true); window.removeEventListener('pointerdown', outside, true); };
  }, [open, close]);
  return <>
    <button ref={trigger} type="button" aria-label={t('open', { title: t(`${help}.title`) })} aria-expanded={open} aria-controls={id}
      data-parameter-help={help} style={compact ? undefined : { width: 44, height: 44, minWidth: 44, padding: 0, borderRadius: '50%' }} className="parameter-help-trigger inline-flex size-6 shrink-0 items-center justify-center rounded-full border border-[var(--gw-control)] bg-[var(--gw-surface)] text-sm font-semibold text-[var(--gw-text)] hover:bg-[var(--gw-active)] focus-visible:outline-solid focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--gw-focus)]"
      onClick={event => { event.preventDefault(); event.stopPropagation(); if (open) close(true); else setOpen(true); }}>ⓘ</button>
    <div ref={panel} id={id} role="note" aria-labelledby={`${id}-title`} tabIndex={-1} popover="manual" hidden={!open}
      style={{ display: open ? 'block' : 'none', position: 'fixed', inset: 'auto', margin: 0, left: position.left, top: position.top, width: position.width, maxHeight: position.maxHeight, zIndex: 100 }}
      className="overflow-y-auto rounded-xl border border-[var(--gw-control)] bg-[var(--gw-raised)] p-4 text-sm leading-relaxed text-[var(--gw-text)] shadow-lg outline-none">
      <h3 id={`${id}-title`} className="mb-2 font-semibold">{t(`${help}.title`)}</h3>
      <p>{t(`${help}.what`)}</p><p className="mt-2">{t(`${help}.how`)}</p><p className="mt-2 text-[var(--gw-muted)]">{t('example')}: {t(`${help}.example`)}</p>
      <button type="button" className="mt-3 rounded border border-[var(--gw-control)] px-3 py-1 text-[var(--gw-text)] hover:bg-[var(--gw-active)] focus-visible:outline-solid focus-visible:outline-2 focus-visible:outline-[var(--gw-focus)]" onClick={() => close(true)}>{t('close')}</button>
    </div>
  </>;
}
