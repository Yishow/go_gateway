import { useId } from 'react';
import { useTranslation } from 'react-i18next';
import { CircleHelp, TriangleAlert } from 'lucide-react';
import { recordingHomeNotRequested, type RecordingHomeAttentionItem, type RecordingHomeGuidance } from './recordingHomeGuidance';
import { RecordingStatusBadge } from './RecordingHomeBadges';

const SHOWN = 4;
/** Lines that say why this item is not confirmed. Unknown is never rendered as healthy or resolved. */
function gapKeys(item: RecordingHomeAttentionItem): string[] {
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
export function RecordingHomeUnconfirmed({ guidance, onInspect, onLogs }: {
  guidance: RecordingHomeGuidance;
  onInspect: (item: RecordingHomeAttentionItem) => void;
  onLogs: (item: RecordingHomeAttentionItem) => void;
}) {
  const { t } = useTranslation('recording-home');
  const headingId = useId();
  // Only items whose details this page requested (or with known unresolved work); the page boundary
  // note already covers items not read. Nothing here requests more detail.
  const items = guidance.availability.state === 'available' ? guidance.attentionItems.filter(item => item.reason === 'operation_recovery' ||
    item.reason === 'evidence_unavailable' && !recordingHomeNotRequested(item)) : [];
  if (!items.length) return null;
  const row = (item: RecordingHomeAttentionItem) => {
    const unresolved = item.reason === 'operation_recovery';
    const Icon = unresolved ? TriangleAlert : CircleHelp;
    return <li key={item.id} data-unconfirmed-kind={unresolved ? 'unresolved' : 'unknown'}>
      <Icon aria-hidden="true" className="gw-home-unconfirmed-icon" size={18} />
      <div className="gw-home-unconfirmed-body">
        <p className="gw-home-unconfirmed-name"><strong>{item.name}</strong><RecordingStatusBadge status={item.evidence.group.status} /></p>
        <ul className="gw-home-unconfirmed-gaps">{gapKeys(item).map(key => <li key={key}>{t(key)}</li>)}</ul>
        {item.evidence.group.status === 'deleted' && <p className="gw-home-unconfirmed-note">{t('unconfirmed.removed')}</p>}
      </div>
      <div className="gw-home-unconfirmed-actions">
        <button type="button" aria-label={t('unconfirmed.inspect', { name: item.name })} onClick={() => onInspect(item)}>{t('unconfirmed.inspectShort')}</button>
        <button type="button" aria-label={t('unconfirmed.logs', { name: item.name })} onClick={() => onLogs(item)}>{t('logs')}</button>
      </div>
    </li>;
  };
  return <section className="gw-home-unconfirmed" aria-labelledby={headingId} data-testid="recording-home-unconfirmed">
    <h2 id={headingId}>{t('unconfirmed.title')}</h2>
    <ul className="gw-home-unconfirmed-list">{items.slice(0, SHOWN).map(row)}</ul>
    {items.length > SHOWN && <details className="gw-home-unconfirmed-more"><summary>{t('unconfirmed.more', { count: items.length - SHOWN })}</summary>
      <ul className="gw-home-unconfirmed-list">{items.slice(SHOWN).map(row)}</ul></details>}
  </section>;
}
