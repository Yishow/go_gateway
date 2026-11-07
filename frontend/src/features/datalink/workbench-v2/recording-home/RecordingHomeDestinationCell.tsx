import { useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { ArrowRight, Copy, Database, Info } from 'lucide-react';
import type { WriteGroup } from '@/types/studioV2WriteGroup';
import type { RecordingReadonly } from '../state/recordingWorkflowView';
import './recordingHomeRows.css';

export function RecordingHomeDestinationCell({ destination }: { destination: RecordingReadonly<WriteGroup>['destination'] }) {
  const { t } = useTranslation('recording-home');
  const [copyStatus, setCopyStatus] = useState<'copied' | 'copyFailed' | null>(null);
  const tableRef = useRef<HTMLSpanElement>(null);
  const identity = [destination.table_schema, destination.table_name].filter(Boolean).join('.');
  const split = Math.ceil(identity.length / 2);
  const copy = async () => {
    try {
      await navigator.clipboard.writeText(identity);
      setCopyStatus('copied');
    } catch {
      const selection = window.getSelection();
      if (selection && tableRef.current) {
        const range = document.createRange();
        range.selectNodeContents(tableRef.current);
        selection.removeAllRanges();
        selection.addRange(range);
      }
      setCopyStatus('copyFailed');
    }
  };
  return <>
    <span className="gw-home-flow" data-flow="destination"><ArrowRight aria-hidden="true" size={14} className="gw-home-flow-arrow" /><Database aria-hidden="true" size={14} /><strong className="gw-destination-name" title={destination.table_name}>{destination.storage_strategy === 'managed' ? t('managedTable') : destination.table_name || t('unknown')}</strong></span>
    <div className="gw-destination-identity-line">
      <span ref={tableRef} className="gw-table-identity" title={identity}><span>{identity.slice(0, split)}</span><span>{identity.slice(split)}</span></span>
      {identity && <button type="button" className="gw-row-icon-button" aria-label={t('rowEvidence.copyTable')} onClick={copy}><Copy aria-hidden="true" size={14} /></button>}
      {destination.connector_id && <span className="gw-connector-hint" title={`${t('product.connector')}: ${destination.connector_id}`} aria-label={`${t('product.connector')}: ${destination.connector_id}`} tabIndex={0}><Info aria-hidden="true" size={14} /></span>}
    </div>
    <p className="gw-destination-name">{(destination.database || '').split(/[\\/]/).filter(Boolean).at(-1) || t('unknown')}</p>
    {copyStatus && <p role="status">{t(`rowEvidence.${copyStatus}`)}</p>}
  </>;
}
