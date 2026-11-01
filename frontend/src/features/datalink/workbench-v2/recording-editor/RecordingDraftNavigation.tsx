import { createContext, useCallback, useContext, useEffect, useLayoutEffect, useRef, useState, type ReactNode } from 'react';
import * as Dialog from '@radix-ui/react-dialog';
import { useTranslation } from 'react-i18next';
import '../../../../styles/semantic-modal.css';

export interface RecordingDraftGuard {
  /** Page intent only; never a server revision, receipt or mutation queue. */
  version: string;
  dirty: boolean;
  pending?: boolean;
  /** Reuse this already-open source modal for a page-leave confirmation. */
  host?: boolean;
  discard?: () => void;
}
type GuardReader = () => RecordingDraftGuard;
type Pending = { proceed: () => void; cancel?: () => void; signature: string; focus: HTMLElement | null };
interface Navigation {
  register: (key: symbol, read: GuardReader) => () => void;
  changed: () => void;
  blocked: () => boolean;
  request: (proceed: () => void, cancel?: () => void) => void;
  confirming: boolean;
  confirmation?: ReactNode;
  keep: () => void;
}
const NavigationContext = createContext<Navigation | null>(null);
const directNavigation: Navigation = { register: () => () => undefined, changed: () => undefined, blocked: () => false, request: proceed => proceed(), confirming: false, keep: () => undefined };
export function useRecordingNavigation() { return useContext(NavigationContext) ?? directNavigation; }

/** One transient leave decision shared by editor Close, page navigation and the router.
 * Keeping never rebuilds the form. Discard touches only registered unsent intent. */
export function RecordingDraftNavigation({ children }: { children: ReactNode }) {
  const { t } = useTranslation('recording-editor');
  const guards = useRef(new Map<symbol, { id: number; read: GuardReader }>());
  const serial = useRef(0);
  const [pending, setPending] = useState<Pending>();
  const [, refresh] = useState(0);
  const changed = useCallback(() => refresh(value => value + 1), []);
  const register = useCallback((key: symbol, read: GuardReader) => {
    guards.current.set(key, { id: ++serial.current, read });
    return () => { guards.current.delete(key); };
  }, []);
  const active = () => [...guards.current.values()].map(({ read }) => read()).filter(value => value.dirty || value.pending);
  const signature = () => JSON.stringify([...guards.current.values()].map(({ id, read }) => { const value = read(); return [id, value.version, value.dirty, value.pending, value.host]; }));
  const blocked = useCallback(() => [...guards.current.values()].some(({ read }) => { const value = read(); return value.dirty || value.pending; }), []);
  const request = useCallback((proceed: () => void, cancel?: () => void) => {
    if (!blocked()) { proceed(); return; }
    setPending({ proceed, cancel, signature: signature(), focus: document.activeElement instanceof HTMLElement ? document.activeElement : null });
  // Readers are live refs; callbacks never borrow a stale dirty value.
  }, [blocked]);
  const keep = () => { pending?.cancel?.(); setPending(undefined); const target = pending?.focus; requestAnimationFrame(() => { if (target?.isConnected) target.focus(); }); };
  const currentSignature = signature();
  useEffect(() => {
    if (pending && pending.signature !== currentSignature) { pending.cancel?.(); setPending(undefined); }
  }, [pending, currentSignature]);
  useEffect(() => {
    const unload = (event: BeforeUnloadEvent) => {
      if (blocked()) { event.preventDefault(); event.returnValue = ''; }
    };
    window.addEventListener('beforeunload', unload);
    return () => window.removeEventListener('beforeunload', unload);
  }, [blocked]);
  const focus = useRef<HTMLElement | null>(null);
  if (pending) focus.current = pending.focus;
  const hosted = [...guards.current.values()].some(({ read }) => read().host);
  const keepButton = useRef<HTMLButtonElement>(null);
  useLayoutEffect(() => { if (pending && hosted) keepButton.current?.focus(); }, [pending, hosted]);
  const confirmation = pending ? <>
    <header className="gw-modal-header"><Dialog.Title>{t('navigation.title')}</Dialog.Title>
      <Dialog.Description>{t('navigation.description')}</Dialog.Description></header>
    <div className="gw-modal-body space-y-3"><p>{t('navigation.savedSources')}</p>
      {active().some(value => value.pending) && <p role="status">{t('navigation.pending')}</p>}
    </div>
    <footer className="gw-modal-footer">
      <button ref={keepButton} type="button" className="gw-secondary" onClick={keep}>{t('navigation.keep')}</button>
      <button type="button" className="gw-primary" onClick={() => {
        if (signature() !== pending.signature) { keep(); return; }
        const leave = pending.proceed;
        for (const value of active()) value.discard?.();
        setPending(undefined); leave();
      }}>{t('navigation.discard')}</button>
    </footer>
  </> : undefined;
  return <NavigationContext.Provider value={{ register, changed, blocked, request, confirming: Boolean(pending), confirmation: hosted ? confirmation : undefined, keep }}>
    {children}
    <Dialog.Root open={Boolean(pending) && !hosted} onOpenChange={open => { if (!open) keep(); }}>
      <Dialog.Portal><Dialog.Overlay className="gw-modal-overlay" />
        <Dialog.Content role="alertdialog" className="gw-modal" data-size="confirm"
          onPointerDownOutside={event => event.preventDefault()}
          onCloseAutoFocus={event => { event.preventDefault(); if (focus.current?.isConnected) focus.current.focus(); }}>
          {confirmation}
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  </NavigationContext.Provider>;
}
export function RecordingDraftBoundary({ children }: { children: ReactNode }) {
  return useContext(NavigationContext) ? children : <RecordingDraftNavigation>{children}</RecordingDraftNavigation>;
}
export function useRecordingDraftGuard(value: RecordingDraftGuard) {
  const navigation = useRecordingNavigation();
  const current = useRef(value); current.current = value;
  const key = useRef(Symbol('recording-draft'));
  const { register, changed } = navigation;
  useLayoutEffect(() => register(key.current, () => current.current), [register]);
  useLayoutEffect(changed, [changed, value.version, value.dirty, value.pending, value.host]);
  return navigation;
}

// Denial-only projection of a failed page read. It cannot grant readiness or replace any source owner.
const RecordingReadUnavailable = createContext(false);
export function RecordingDraftReadState({ unavailable, children }: { unavailable: boolean; children: ReactNode }) {
  return <RecordingReadUnavailable.Provider value={unavailable}>{children}</RecordingReadUnavailable.Provider>;
}
export function useRecordingReadUnavailable() { return useContext(RecordingReadUnavailable); }
