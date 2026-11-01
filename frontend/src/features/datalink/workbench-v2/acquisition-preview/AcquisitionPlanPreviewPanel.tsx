import { SectionCard } from '../components/SectionCard';
import { Button } from '../components/Button';
import { useTranslation } from 'react-i18next';
import type { AcquisitionPreviewScope } from '../../../../types/acquisitionPlanPreview';
import { useAcquisitionPlanPreview } from './useAcquisitionPlanPreview';
import { AcquisitionPreviewGroup } from './AcquisitionPreviewGroup';
export interface AcquisitionPlanPreviewPanelProps { scope?: AcquisitionPreviewScope; refreshWorkspaceHead?: boolean; unavailableReason?: 'unsaved' | 'scope' }

/** A read-only source subpanel. It has no mutation callbacks or save/activation gate. */
export function AcquisitionPlanPreviewPanel({ scope, refreshWorkspaceHead = false, unavailableReason = 'scope' }: AcquisitionPlanPreviewPanelProps) {
  const { t } = useTranslation('acquisition-preview');
  const { result, loading, refresh } = useAcquisitionPlanPreview(scope, refreshWorkspaceHead);
  const preview = result?.ok ? result.value : undefined;
  return <div data-testid="acquisition-plan-preview" className="min-w-0">
    <SectionCard title={t('title')} contentClassName="space-y-3 text-xs leading-relaxed text-slate-300"
      aside={scope && <Button type="button" variant="secondary" size="sm" disabled={loading} onClick={refresh}>{t('refresh')}</Button>}>
    <p className="text-slate-400">{t('readOnlyNote')}</p>
    {!scope ? <p role="status">{t(unavailableReason === 'unsaved' ? 'unsaved' : 'scopeUnavailable')}</p>
      : loading ? <p role="status">{t('loading')}</p>
      : result && !result.ok ? <p role="status">{t(result.reason === 'stale-reference' ? 'stale' : result.reason === 'foreign-scope' ? 'foreign' : 'unavailable')}</p>
      : preview?.status === 'unavailable' ? <div role="status"><p>{t('unavailable')}</p><details><summary>{t('details')}</summary><code className="break-all">{preview.reason}</code></details></div>
      : preview?.status === 'available' ? <>
        <p>{t('estimateNote')}</p>
        {preview.groups.map(group => <AcquisitionPreviewGroup key={group.polling_group_id} group={group} />)}
        {preview.groups.length === 0 && <p>{t('emptyInstalled')}</p>}
        <details><summary>{t('lossTitle')}</summary>
          {preview.output_loss.status === 'partial' ? <><p>{t('lossPartial')}</p>
            <ul>{preview.output_loss.retained_detail.map(row => <li key={row.point_id}>{t('lossRow', { point: row.point_id, count: row.losses })}</li>)}</ul>
          </> : <p>{t('lossUnavailable')}</p>}
        </details>
      </> : null}
    </SectionCard>
  </div>;
}
