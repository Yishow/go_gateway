/** Reveal restored keyboard focus without changing the collection's independently restored scroll. */
function revealDocumentFocus(target: HTMLElement) {
  const view = target.ownerDocument.defaultView, rect = target.getBoundingClientRect();
  if (!view || rect.height <= 0) return;
  const viewport = view.visualViewport, top = (viewport?.offsetTop ?? 0) + 8;
  const bottom = (viewport?.offsetTop ?? 0) + (viewport?.height ?? view.innerHeight) - 8;
  const movement = rect.top < top ? rect.top - top : rect.bottom > bottom ? Math.min(rect.top - top, rect.bottom - bottom) : 0;
  if (movement) view.scrollBy({ top: movement, behavior: 'instant' });
}

/** Current-visit ticket: late layout may reveal the same focus, but user intent always cancels it. */
export function watchHomeReturnFocus(target: HTMLElement) {
  const document = target.ownerDocument, view = document.defaultView;
  let active = true, observer: IntersectionObserver | undefined;
  const intents = ['wheel', 'pointerdown', 'keydown', 'touchstart'] as const;
  const cancel = () => {
    if (!active) return;
    active = false; observer?.disconnect();
    intents.forEach(type => document.removeEventListener(type, cancel, true));
    document.removeEventListener('focusin', focusChanged, true);
    view?.removeEventListener('blur', cancel);
  };
  const focusChanged = () => { if (document.activeElement !== target) cancel(); };
  const reveal = () => {
    if (!active) return;
    if (!target.isConnected || document.activeElement !== target) { cancel(); return; }
    revealDocumentFocus(target);
  };
  intents.forEach(type => document.addEventListener(type, cancel, { capture: true, passive: true }));
  document.addEventListener('focusin', focusChanged, true); view?.addEventListener('blur', cancel);
  if (typeof IntersectionObserver !== 'undefined') {
    // Siblings can translate this target without resizing its ancestors. Observe visibility,
    // using the actual viewport edge so subpixel rounding of the 8px reveal inset is harmless.
    observer = new IntersectionObserver(reveal, { threshold: 1 });
    observer.observe(target);
  }
  reveal();
  return { cancel, finish: () => { reveal(); cancel(); } };
}
