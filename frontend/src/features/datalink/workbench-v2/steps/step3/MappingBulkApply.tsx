import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { Device, Mapping, Point } from '../../state/types';
import type { WorkbenchV2Action } from '../../state/useWorkbenchV2State';

type Field = 'scale' | 'offset' | 'target_type' | 'unit';
interface Props {
  points: Point[];
  mappings: Record<string, Mapping>;
  devices: Device[];
  source?: Point | null;
  targets: Set<string>;
  setTargets: (ids: Set<string>) => void;
  dispatch: React.Dispatch<WorkbenchV2Action>;
}

export function MappingBulkApply({ points, mappings, devices, source, targets, setTargets, dispatch }: Props) {
  const { t } = useTranslation('workbench-v2');
  const [preview, setPreview] = useState<{ fields: Field[]; signature: string } | null>(null);
  const targetPoints = points.filter(p => targets.has(p.id) && p.id !== source?.id && mappings[p.id]?.enabled);
  const sourceMapping = source ? mappings[source.id] : undefined;
  // Any edit, source change, membership change or reload invalidates consent.
  const fingerprint = (m?: Mapping) => m && [m.point_id, m.persisted_point_id, m.tag_id, m.tag_key, m.display_name, m.enabled, m.scale, m.offset, m.unit, m.target_type];
  const signature = JSON.stringify([source?.id, fingerprint(sourceMapping), targetPoints.map(p => [p.id, p.device_id, p.address, fingerprint(mappings[p.id])])]);
  const validPreview = preview?.signature === signature ? preview : null;
  const label = (p: Point) => `${devices.find(d => d.id === p.device_id)?.name ?? p.device_id} · ${p.address} · ${mappings[p.id]?.display_name}`;
  const startPreview = (fields: Field[]) => setPreview({ fields, signature });
  const tag = sourceMapping?.tag_key ?? '';
  return <div className="space-y-3 border-t border-slate-800 p-3 text-xs">
    <p>{t('step3.bulk.explanation')}</p>
    <button type="button" onClick={() => { setTargets(new Set(points.filter(p => mappings[p.id]?.enabled && p.id !== source?.id).map(p => p.id))); setPreview(null); }} data-testid="bulk-select-all">{t('step3.bulk.selectAll')}</button>
    {' · '}<button type="button" onClick={() => { setTargets(new Set()); setPreview(null); }}>{t('step3.bulk.clear')}</button>
    {source && sourceMapping ? <>
      <p data-testid="bulk-source">{t('step3.table.bulkApplySource', { tag: label(source) })}</p>
      <div className="flex flex-wrap gap-3">
        <button type="button" data-testid="btn-bulk-apply-all" disabled={!targetPoints.length} onClick={() => startPreview(['scale', 'offset', 'target_type', 'unit'])}>{t('step3.table.bulkApplyAll', { tag })}</button>
        {(['scale', 'offset', 'target_type', 'unit'] as Field[]).map(field => <button type="button" key={field} disabled={!targetPoints.length} data-testid={`btn-bulk-apply-${field === 'target_type' ? 'target-type' : field}`} onClick={() => startPreview([field])}>{t(`step3.bulk.${field}`)}</button>)}
      </div>
    </> : <p data-testid="btn-bulk-apply-all-disabled">{t('step3.table.bulkApplyPlaceholder')}</p>}
    {validPreview && source && sourceMapping && <section aria-label={t('step3.bulk.preview')} className="space-y-2 rounded border border-amber-600 p-3" data-testid="bulk-preview">
      <h4>{t('step3.bulk.preview')} · {targetPoints.length}</h4>
      <p>{t('step3.bulk.warning')}</p>
      <ul>{targetPoints.map(p => <li key={p.id}>{label(p)}: {validPreview.fields.map(field => `${field}: ${String(mappings[p.id][field])} → ${String(sourceMapping[field])}`).join('; ')}</li>)}</ul>
      <button type="button" data-testid="bulk-confirm" onClick={() => { dispatch({ type: 'bulkApplyTransform', fromPointId: source.id, toPointIds: targetPoints.map(p => p.id), fields: validPreview.fields }); setPreview(null); setTargets(new Set()); }}>{t('step3.bulk.confirm')}</button>
      {' · '}<button type="button" onClick={() => setPreview(null)}>{t('step3.bulk.cancel')}</button>
    </section>}
  </div>;
}
