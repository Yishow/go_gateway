import { useQueryClient } from '@tanstack/react-query';
import { studioV2WorkspaceKeys } from '../../../../hooks/datalink/keys';
import { useRecordingNavigation } from '../recording-editor/RecordingDraftNavigation';
import type { RuleQueueObservation } from '../../../../pages/datalink/workbench-v2/sourceQueueObservation';
import { useCallback, useRef, useState, type ReactNode, type MutableRefObject } from 'react';
import * as Dialog from '@radix-ui/react-dialog';
import { useTranslation } from 'react-i18next';
import { captureReturnFocus, restoreReturnFocus, type ReturnFocusTicket } from '../../../../utils/returnFocus';
import { semanticTokens } from '../tokens';
import '../../../../styles/semantic-modal.css';
/** `view: 'memory'` is a read-only projection of the range list; `startAddress` only prefills a new range;
 * `returnRuleId` / `returnNewRange` name the setting the memory view returns to. */
export interface HomeSourceTarget { workspaceId: string; deviceId: string; kind: 'connection' | 'points'; ruleId?: string; newRange?: boolean; view?: 'memory'; startAddress?: string; returnRuleId?: string; returnNewRange?: boolean; initialSelection?: boolean; epoch: number }
export type HomeSourceSelection = HomeSourceTarget | { workspaceId: string; kind: 'create'; deviceId?: never; ruleId?: never; newRange?: never; view?: never; startAddress?: never; returnRuleId?: never; returnNewRange?: never; initialSelection?: never; epoch: number };
export interface HomeSourceSelectOptions { view?: 'memory'; startAddress?: string; returnRuleId?: string; returnNewRange?: boolean }
/** Labels from the validated Home source projection of the same workspace; never a configuration owner. */
export interface HomeSourceDevice { id: string; name: string }
export type HomeRuleObserver = Pick<RuleQueueObservation, 'submitted' | 'saved' | 'failed'>;
export type HomeSourceLeaveGuard = (proceed: () => void) => void;
export interface HomeSourceControls {
  target: HomeSourceSelection; ruleObserverRef: MutableRefObject<HomeRuleObserver | undefined>; closeRef: MutableRefObject<(() => void) | undefined>;
  /** The current editor's leave protection (dirty confirmation, pending/unknown lock) for title-bar switches. */
  leaveRef: MutableRefObject<HomeSourceLeaveGuard | undefined>;
  /** Title-bar slot for controls owned by the current editor's draft. */
  headerSlot: HTMLElement | null;
  onCreated: (deviceId: string) => void; onClose: () => void; onDeviceName: (name: string) => void; onReady: () => void; onLocked: (locked: boolean) => void;
  onSelect: (ruleId?: string, newRange?: boolean, options?: HomeSourceSelectOptions) => void;
}
/** The home remains mounted. This host owns only the modal selection and DOM focus,
 * never a source mutation queue or an inferred persisted identity. */
export function RecordingHomeSourceHost({ children, renderSource }: {
  children: (open: (workspaceId: string, deviceId: string, kind: 'connection' | 'points', devices?: readonly HomeSourceDevice[]) => void, workspaceChanged: (id?: string) => void, create: (workspaceId: string) => void) => ReactNode;
  renderSource: (controls: HomeSourceControls) => ReactNode;
}) {
  const { t } = useTranslation('recording-home');
  const navigation = useRecordingNavigation();
  const context = useQueryClient();
  const [deviceName, setDeviceName] = useState('');
  const [initialSourceLoading, setInitialSourceLoading] = useState(false);
  const [locked, setLocked] = useState(false);
  const [devices, setDevices] = useState<readonly HomeSourceDevice[]>([]);
  const [headerSlot, setHeaderSlot] = useState<HTMLElement | null>(null);
  const [target, setTarget] = useState<HomeSourceSelection>();
  const ruleObserverRef = useRef<HomeRuleObserver | undefined>(undefined);
  const leaveRef = useRef<HomeSourceLeaveGuard | undefined>(undefined);
  const serial = useRef(0), owner = useRef<object>({}), workspace = useRef<string | undefined>(undefined);
  const ticket = useRef<ReturnFocusTicket | undefined>(undefined), closeRef = useRef<(() => void) | undefined>(undefined);
  /** Every new epoch drops the previous editor's owners before it can act again. */
  const reset = () => { closeRef.current = undefined; ruleObserverRef.current = undefined; leaveRef.current = undefined; setLocked(false); };
  const open = useCallback((workspaceId: string, deviceId: string, kind: 'connection' | 'points', listed?: readonly HomeSourceDevice[]) => {
    if (workspace.current !== workspaceId) return;
    const choices = kind === 'points' && listed?.some(item => item.id === deviceId) ? listed : [];
    setDevices(choices); setDeviceName(choices.find(item => item.id === deviceId)?.name ?? ''); setInitialSourceLoading(true);
    owner.current = {}; ticket.current = captureReturnFocus(owner.current, document.activeElement instanceof HTMLElement ? document.activeElement : null);
    reset(); setTarget({ workspaceId, deviceId, kind, initialSelection: true, epoch: ++serial.current });
  }, []);
  const create = useCallback((workspaceId: string) => {
    if (workspace.current !== workspaceId) return;
    setDevices([]); setDeviceName(''); setInitialSourceLoading(false); owner.current = {}; ticket.current = captureReturnFocus(owner.current, document.activeElement instanceof HTMLElement ? document.activeElement : null);
    reset(); setTarget({ workspaceId, kind: 'create', epoch: ++serial.current });
  }, []);
  const workspaceChanged = useCallback((id?: string) => {
    if (id && workspace.current !== id) { workspace.current = id; owner.current = {}; reset(); setDevices([]); setTarget(undefined); }
  }, []);
  const close = () => { if (closeRef.current) closeRef.current(); else setTarget(undefined); };
  const current = (captured: HomeSourceSelection) => captured.kind !== 'create' && workspace.current === captured.workspaceId && captured.epoch === serial.current;
  const select = (captured: HomeSourceSelection, options?: HomeSourceSelectOptions, ruleId?: string, newRange?: boolean) => {
    if (captured.kind === 'create' || !current(captured)) return;
    // React may invoke a state updater twice in Strict Mode; allocate the epoch only once.
    const epoch = ++serial.current;
    reset(); setInitialSourceLoading(true);
    setTarget(value => value?.epoch === captured.epoch ? { ...captured, ruleId, newRange, view: options?.view, startAddress: options?.startAddress, returnRuleId: options?.returnRuleId, returnNewRange: options?.returnNewRange, initialSelection: false, epoch } : value);
  };
  /** Another device of the same validated Home projection; its own scoped locator read decides what is editable. */
  const switchDevice = (captured: HomeSourceSelection, deviceId: string) => {
    if (captured.kind !== 'points' || !current(captured) || captured.deviceId === deviceId) return;
    const choice = devices.find(item => item.id === deviceId);
    if (!choice) return;
    const epoch = ++serial.current;
    reset(); setDeviceName(choice.name); setInitialSourceLoading(true);
    setTarget(value => value?.epoch === captured.epoch ? { workspaceId: captured.workspaceId, deviceId, kind: 'points', initialSelection: true, epoch } : value);
  };
  /** Title-bar switches go through the current editor's own leave protection. */
  const guarded = (proceed: () => void) => { if (locked) return; if (leaveRef.current) leaveRef.current(proceed); else proceed(); };
  const captured = target;
  const points = target?.kind === 'points';
  const settingsOpen = points && !initialSourceLoading && !target.view;
  // Even a single device is shown as the selected one; a just-created device is named by its own scoped read.
  const deviceChoices = !points ? [] : devices.length ? devices : deviceName ? [{ id: target.deviceId, name: deviceName }] : [];
  return <>{children(open, workspaceChanged, create)}<Dialog.Root open={Boolean(target)} onOpenChange={value => { if (!value) close(); }}>
    <Dialog.Portal><Dialog.Overlay className="gw-modal-overlay" /><Dialog.Content data-workbench-v2 className="gw-modal gw-home-source-modal" data-source-kind={target?.kind} data-source-view={target?.view} data-initial-loading={initialSourceLoading} data-size={navigation.confirming ? 'confirm' : 'form'} role={navigation.confirming ? 'alertdialog' : 'dialog'}
      style={{ color: semanticTokens.colors.text, backgroundColor: semanticTokens.colors.raised }}
      onPointerDownOutside={event => event.preventDefault()}
      // While the scoped read is pending, the frame's own Close keeps the initial focus (not a device choice).
      onOpenAutoFocus={event => { const loadingClose = (event.currentTarget as HTMLElement | null)?.querySelector<HTMLElement>('.gw-home-source-loading ~ .gw-modal-footer button'); if (loadingClose) { event.preventDefault(); loadingClose.focus(); } }}
      onEscapeKeyDown={event => { event.preventDefault(); if (!event.isComposing) { if (navigation.confirming) navigation.keep(); else close(); } }}
      onCloseAutoFocus={event => { event.preventDefault(); restoreReturnFocus(ticket.current, owner.current, document.querySelector<HTMLElement>('[data-testid="recording-home"] h1')); }}>
      <>{navigation.confirmation ?? <header className="gw-modal-header"><div><Dialog.Title>{deviceName && `${deviceName} · `}{t(target?.kind === 'create' ? 'createDevice' : target?.kind === 'connection' ? 'source.connectionTitle' : target?.view === 'memory' ? 'source.memory.title' : 'source.pointsTitle')}</Dialog.Title>
        <Dialog.Description>{t('source.sharedDescription')}</Dialog.Description></div>
        {points && <span ref={setHeaderSlot} className="gw-home-source-slot" />}
        {settingsOpen && <button type="button" data-testid="home-range-memory" disabled={locked} onClick={() => guarded(() => select(target, { view: 'memory', returnRuleId: target.ruleId, returnNewRange: Boolean(target.newRange) }))}>{t('source.memory.open')}</button>}
        {(initialSourceLoading || target?.kind !== 'connection' && !target?.ruleId && !target?.newRange && !target?.view) && <button type="button" data-testid="home-source-close" onClick={close}>{t('source.close')}</button>}
        {points && !target.view && deviceChoices.length > 0 && <div role="group" aria-label={t('source.devices')} className="gw-home-source-devices">{deviceChoices.map(item =>
          <button type="button" key={item.id} data-testid={`home-source-device-${item.id}`} aria-pressed={item.id === target.deviceId} disabled={locked}
            onClick={() => { if (item.id !== target.deviceId) guarded(() => switchDevice(target, item.id)); }}>
            <svg aria-hidden="true" viewBox="0 0 16 16" width="16" height="16"><rect x="2" y="3" width="12" height="8" rx="1.5" fill="none" stroke="currentColor" strokeWidth="1.5" /><path d="M5 13.5h6M8 11v2.5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" /></svg>
            <span>{item.name}</span></button>)}</div>}
      </header>}<div style={{ display: navigation.confirming ? 'none' : 'contents' }}>
      {captured && renderSource({ target: captured, closeRef, ruleObserverRef, leaveRef, headerSlot: navigation.confirming ? null : headerSlot, onCreated: deviceId => {
          if (!deviceId || workspace.current !== captured.workspaceId || captured.kind !== 'create' || captured.epoch !== serial.current) return;
          void context.invalidateQueries({ queryKey: [...studioV2WorkspaceKeys.all, 'recording-home-sources'] });
          // Allocate outside the updater: Strict Mode may invoke it twice.
          const epoch = ++serial.current;
          reset(); setInitialSourceLoading(true);
          setTarget(value => value?.epoch === captured.epoch ? { workspaceId: captured.workspaceId, deviceId, kind: 'points', initialSelection: true, epoch } : value);
        }, onDeviceName: name => { if (workspace.current === captured.workspaceId && captured.epoch === serial.current) setDeviceName(name); },
        onReady: () => { if (workspace.current === captured.workspaceId && captured.epoch === serial.current) setInitialSourceLoading(false); },
        onLocked: value => { if (workspace.current === captured.workspaceId && captured.epoch === serial.current) setLocked(value); },
        onClose: () => { if (workspace.current === captured.workspaceId) setTarget(value => value?.epoch === captured.epoch ? undefined : value); },
        onSelect: (ruleId, newRange, options) => select(captured, options, ruleId, newRange) })}
    </div></></Dialog.Content></Dialog.Portal>
  </Dialog.Root></>;
}
