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
export interface HomeSourceTarget { workspaceId: string; deviceId: string; kind: 'connection' | 'points'; ruleId?: string; newRange?: boolean; initialSelection?: boolean; epoch: number }
export type HomeSourceSelection = HomeSourceTarget | { workspaceId: string; kind: 'create'; deviceId?: never; ruleId?: never; newRange?: never; initialSelection?: never; epoch: number };
export type HomeRuleObserver = Pick<RuleQueueObservation, 'submitted' | 'saved' | 'failed'>;
export interface HomeSourceControls {
  target: HomeSourceSelection; ruleObserverRef: MutableRefObject<HomeRuleObserver | undefined>; closeRef: MutableRefObject<(() => void) | undefined>;
  onCreated: (deviceId: string) => void; onClose: () => void; onDeviceName: (name: string) => void; onSelect: (ruleId?: string, newRange?: boolean) => void;
}
/** The home remains mounted. This host owns only the modal selection and DOM focus,
 * never a source mutation queue or an inferred persisted identity. */
export function RecordingHomeSourceHost({ children, renderSource }: {
  children: (open: (workspaceId: string, deviceId: string, kind: 'connection' | 'points') => void, workspaceChanged: (id?: string) => void, create: (workspaceId: string) => void) => ReactNode;
  renderSource: (controls: HomeSourceControls) => ReactNode;
}) {
  const { t } = useTranslation('recording-home');
  const navigation = useRecordingNavigation();
  const context = useQueryClient();
  const [deviceName, setDeviceName] = useState('');
  const [target, setTarget] = useState<HomeSourceSelection>();
  const ruleObserverRef = useRef<HomeRuleObserver | undefined>(undefined);
  const serial = useRef(0), owner = useRef<object>({}), workspace = useRef<string | undefined>(undefined);
  const ticket = useRef<ReturnFocusTicket | undefined>(undefined), closeRef = useRef<(() => void) | undefined>(undefined);
  const open = useCallback((workspaceId: string, deviceId: string, kind: 'connection' | 'points') => {
    if (workspace.current !== workspaceId) return;
    setDeviceName(''); owner.current = {}; ticket.current = captureReturnFocus(owner.current, document.activeElement instanceof HTMLElement ? document.activeElement : null);
    closeRef.current = undefined; ruleObserverRef.current = undefined; setTarget({ workspaceId, deviceId, kind, initialSelection: true, epoch: ++serial.current });
  }, []);
  const create = useCallback((workspaceId: string) => {
    if (workspace.current !== workspaceId) return;
    setDeviceName(''); owner.current = {}; ticket.current = captureReturnFocus(owner.current, document.activeElement instanceof HTMLElement ? document.activeElement : null);
    closeRef.current = undefined; ruleObserverRef.current = undefined; setTarget({ workspaceId, kind: 'create', epoch: ++serial.current });
  }, []);
  const workspaceChanged = useCallback((id?: string) => {
    if (id && workspace.current !== id) { workspace.current = id; owner.current = {}; closeRef.current = undefined; ruleObserverRef.current = undefined; setTarget(undefined); }
  }, []);
  const close = () => { if (closeRef.current) closeRef.current(); else setTarget(undefined); };
  const captured = target;
  return <>{children(open, workspaceChanged, create)}<Dialog.Root open={Boolean(target)} onOpenChange={value => { if (!value) close(); }}>
    <Dialog.Portal><Dialog.Overlay className="gw-modal-overlay" /><Dialog.Content data-workbench-v2 className="gw-modal gw-home-source-modal" data-size={navigation.confirming ? 'confirm' : 'form'} role={navigation.confirming ? 'alertdialog' : 'dialog'}
      style={{ color: semanticTokens.colors.text, backgroundColor: semanticTokens.colors.raised }}
      onPointerDownOutside={event => event.preventDefault()}
      onEscapeKeyDown={event => { event.preventDefault(); if (!event.isComposing) { if (navigation.confirming) navigation.keep(); else close(); } }}
      onCloseAutoFocus={event => { event.preventDefault(); restoreReturnFocus(ticket.current, owner.current, document.querySelector<HTMLElement>('[data-testid="recording-home"] h1')); }}>
      <>{navigation.confirmation ?? <header className="gw-modal-header"><div><Dialog.Title>{deviceName && `${deviceName} · `}{t(target?.kind === 'create' ? 'createDevice' : target?.kind === 'connection' ? 'source.connectionTitle' : 'source.pointsTitle')}</Dialog.Title>
        <Dialog.Description>{t('source.sharedDescription')}</Dialog.Description></div>
        {target?.kind !== 'connection' && !target?.ruleId && <button type="button" data-testid="home-source-close" onClick={close}>{t('source.close')}</button>}</header>}<div style={{ display: navigation.confirming ? 'none' : 'contents' }}>
      {captured && renderSource({ target: captured, closeRef, ruleObserverRef, onCreated: deviceId => {
          if (!deviceId || workspace.current !== captured.workspaceId || captured.kind !== 'create' || captured.epoch !== serial.current) return;
          void context.invalidateQueries({ queryKey: [...studioV2WorkspaceKeys.all, 'recording-home-sources'] });
          closeRef.current = undefined; ruleObserverRef.current = undefined;
          setTarget(current => current?.epoch === captured.epoch ? { workspaceId: captured.workspaceId, deviceId, kind: 'points', initialSelection: true, epoch: ++serial.current } : current);
        }, onDeviceName: name => { if (workspace.current === captured.workspaceId) setDeviceName(name); }, onClose: () => { if (workspace.current === captured.workspaceId) setTarget(current => current?.epoch === captured.epoch ? undefined : current); },
        onSelect: (ruleId, newRange) => { if (captured.kind !== 'create' && workspace.current === captured.workspaceId) { closeRef.current = undefined; ruleObserverRef.current = undefined;
          setTarget(current => current?.epoch === captured.epoch ? { ...captured, ruleId, newRange, initialSelection: false, epoch: ++serial.current } : current); } } })}
    </div></></Dialog.Content></Dialog.Portal>
  </Dialog.Root></>;
}
