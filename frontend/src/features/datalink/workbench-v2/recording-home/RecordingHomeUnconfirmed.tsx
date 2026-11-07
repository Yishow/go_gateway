import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { CircleHelp, TriangleAlert } from 'lucide-react';
import { recordingHomeNotRequested, type RecordingHomeAttentionItem, type RecordingHomeGuidance } from './recordingHomeGuidance';
import { RecordingStatusBadge } from './RecordingHomeBadges';
import { RecordingHomeRowActions } from './RecordingHomeRowActions';
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
  const [expanded, setExpanded] = useState<boolean | null>(null);
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
        <button type="button" aria-label={t('unconfirmed.inspect', { name: item.name })} onClick={() => onInspect(item)}>{t('unconfirmed.inspectShort')}</button>
        <button type="button" aria-label={t('unconfirmed.logs', { name: item.name })} onClick={() => onLogs(item)}>{t('logs')}</button>
        {group.status === 'deleted' && onDelete && (!matches || canDelete) && <RecordingHomeRowActions
          label={t('rowEvidence.actionsFor', { name: item.name })} hint={!matches ? t('rowEvidence.reloadHint') : undefined}
          items={[{ id: 'delete', label: t('rowEvidence.delete'), disabled: !canDelete, onSelect: () => onDelete(group) }]} />}
      </div>
    </li>;
  };
  return <details className="gw-home-unconfirmed" role="region" aria-label={t('unconfirmed.title')} data-testid="recording-home-unconfirmed"
    open={expanded ?? items.some(item => item.evidence.group.status !== 'deleted')} onToggle={event => setExpanded(event.currentTarget.open)}>
    <summary><h2>{t('deviceStatus.attention', { count: items.length })}</h2></summary>
    <ul className="gw-home-unconfirmed-list">{items.slice(0, SHOWN).map(row)}</ul>
    {items.length > SHOWN && <details className="gw-home-unconfirmed-more"><summary>{t('unconfirmed.more', { count: items.length - SHOWN })}</summary>
      <ul className="gw-home-unconfirmed-list">{items.slice(SHOWN).map(row)}</ul></details>}
  </details>;
}
