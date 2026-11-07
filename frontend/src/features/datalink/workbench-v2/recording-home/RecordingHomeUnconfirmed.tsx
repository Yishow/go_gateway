import { useTranslation } from 'react-i18next';
import { CircleHelp, ScrollText, Search, Trash2, TriangleAlert } from 'lucide-react';
import { recordingHomeNotRequested, type RecordingHomeAttentionItem, type RecordingHomeGuidance } from './recordingHomeGuidance';
import { RecordingStatusBadge } from './RecordingHomeBadges';
import { recordingHomeLifecycleMatches } from './recordingHomeLifecycleMatches';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { RecordingReadonly } from '../state/recordingWorkflowView';

const SHOWN = 4;
/** Lines that say why an item needs attention or remains unconfirmed. */
function gapKeys(item: RecordingHomeAttentionItem): string[] {
  if (item.reason === 'delivery_attention' || item.reason === 'source_attention') return [`reason.${item.reason}`];
  if (item.reason === 'operation_recovery') {
    const owner = item.action.target === 'recording' ? item.action.intent?.operationOwner : undefined;
    return [owner ? `unconfirmed.unresolved.${owner}` : 'unconfirmed.unresolved.unknown'];
  }
  const gaps = item.unconfirmed;
  const keys = gaps ? [...(gaps.lifecycle ? ['unconfirmed.read.lifecycle'] : []), ...(gaps.delivery ? ['unconfirmed.read.delivery'] : []),
    ...gaps.owners.map(owner => `unconfirmed.owner.${owner}`)] : [];
  return keys.length ? keys : ['unconfirmed.current'];
}
/** Names every unresolved or unconfirmed item with its saved status, gaps and read-only check entries. */
export function RecordingHomeUnconfirmed({ guidance, onInspect, onLogs, onDelete }: {
  guidance: RecordingHomeGuidance;
  onInspect: (item: RecordingHomeAttentionItem) => void;
  onLogs: (item: RecordingHomeAttentionItem) => void;
  onDelete?: (group: RecordingReadonly<WriteGroup>) => void;
}) {
  const { t } = useTranslation('recording-home');
  // Only items whose details this page requested (or with known unresolved work); the page boundary
  // note already covers items not read. Nothing here requests more detail.
  const items = guidance.availability.state === 'available' ? guidance.attentionItems.filter(item => item.reason !== 'evidence_unavailable' ||
    !recordingHomeNotRequested(item)) : [];
  if (!items.length) return null;
  const row = (item: RecordingHomeAttentionItem) => {
    const attention = item.reason !== 'evidence_unavailable';
    const Icon = attention ? TriangleAlert : CircleHelp;
    const group = item.evidence.group;
    const matches = recordingHomeLifecycleMatches(group, item);
    const lifecycle = item.evidence.lifecycle;
    const canDelete = matches && lifecycle.state === 'available' && lifecycle.value.actions.some(action => action.action === 'remove' && action.allowed);
    return <li key={item.id} data-unconfirmed-kind={attention ? 'unresolved' : 'unknown'}>
      <Icon aria-hidden="true" className="gw-home-unconfirmed-icon" size={18} />
      <div className="gw-home-unconfirmed-body">
        <p className="gw-home-unconfirmed-name"><strong>{item.name}</strong><RecordingStatusBadge status={item.evidence.group.status} /></p>
        <ul className="gw-home-unconfirmed-gaps">{gapKeys(item).map(key => <li key={key}>{t(key)}</li>)}</ul>
        {item.evidence.group.status === 'deleted' && <p className="gw-home-unconfirmed-note">{t('unconfirmed.removed')}</p>}
      </div>
      <div className="gw-home-unconfirmed-actions">
        <button type="button" aria-label={t('unconfirmed.inspect', { name: item.name })} onClick={() => onInspect(item)}><Search aria-hidden="true" size={15} /><span>{t('unconfirmed.inspectShort')}</span></button>
        <button type="button" aria-label={t('unconfirmed.logs', { name: item.name })} onClick={() => onLogs(item)}><ScrollText aria-hidden="true" size={15} /><span>{t('logs')}</span></button>
        {group.status === 'deleted' && onDelete && (!matches || canDelete) && <button type="button" className="gw-row-icon-button gw-danger-ghost"
          aria-label={t('rowEvidence.deleteNamed', { name: item.name })} disabled={!canDelete} title={!matches ? t('rowEvidence.reloadHint') : t('rowEvidence.delete')}
          onClick={event => { event.currentTarget.focus(); onDelete(group); }}><Trash2 aria-hidden="true" size={16} /></button>}
      </div>
      {!matches && group.status === 'deleted' && onDelete && <p className="gw-home-unconfirmed-hint">{t('rowEvidence.reloadHint')}</p>}
    </li>;
  };
  return <section className="gw-home-unconfirmed" aria-label={t('unconfirmed.title')} data-testid="recording-home-unconfirmed">
    <header className="gw-home-unconfirmed-header"><TriangleAlert aria-hidden="true" size={18} />
      <h2>{t('deviceStatus.attention', { count: items.length })}</h2><p>{t('unconfirmed.hint')}</p></header>
    <ul className="gw-home-unconfirmed-list">{items.slice(0, SHOWN).map(row)}</ul>
    {items.length > SHOWN && <details className="gw-home-unconfirmed-more"><summary>{t('unconfirmed.more', { count: items.length - SHOWN })}</summary>
      <ul className="gw-home-unconfirmed-list">{items.slice(SHOWN).map(row)}</ul></details>}
  </section>;
}
