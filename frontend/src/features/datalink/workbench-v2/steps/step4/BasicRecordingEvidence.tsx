import * as React from 'react';
import { useTranslation } from 'react-i18next';
import type { WriteGroup } from '../../../../../types/studioV2WriteGroup';
import type {
  WriteGroupCommittedEffect,
  WriteGroupDeliveryStages,
} from '../../../../../types/studioV2WriteGroupDelivery';

interface BasicRecordingEvidenceProps {
  intervalSeconds: number;
  deliveryError: boolean;
  stages?: WriteGroupDeliveryStages;
  committedEffect?: WriteGroupCommittedEffect;
  group?: WriteGroup;
}

export const BasicRecordingEvidence: React.FC<BasicRecordingEvidenceProps> = ({
  intervalSeconds,
  deliveryError,
  stages,
  committedEffect,
  group,
}) => {
  const { t } = useTranslation('workbench-v2');

  return (
    <section className="space-y-2 rounded border border-[var(--gw-border)] bg-[var(--gw-surface)] p-3" aria-label={t('step4.basic.evidence_title')} data-testid="basic-recording-evidence">
      <h4 className="text-xs font-semibold text-[var(--gw-text)]">{t('step4.basic.evidence_title')}</h4>
      {!committedEffect && <p className="text-xs text-[var(--gw-muted)]">{t('step4.basic.first_bucket', { interval: intervalSeconds })}</p>}
      {deliveryError && <p role="alert" className="text-xs text-[var(--gw-warning-fg)]">{t('step4.basic.delivery_failed')}</p>}
      <dl className="grid grid-cols-3 gap-2 text-xs">
        <div><dt className="text-[var(--gw-muted)]">{t('step4.basic.collecting')}</dt><dd className="text-[var(--gw-text)]">{stages ? stages.collecting : t('step4.basic.unconfirmed')}</dd></div>
        <div><dt className="text-[var(--gw-muted)]">{t('step4.basic.local')}</dt><dd className="text-[var(--gw-text)]">{stages ? stages.queued + stages.retrying : t('step4.basic.unconfirmed')}</dd></div>
        <div><dt className="text-[var(--gw-muted)]">{t('step4.basic.committed')}</dt><dd className="text-[var(--gw-text)]">{committedEffect ? t('step4.basic.committed_verified') : t('step4.basic.unconfirmed')}</dd></div>
      </dl>
      {committedEffect && group && (
        <details><summary className="cursor-pointer text-xs text-[var(--gw-muted)]">{t('step4.basic.diagnostics')}</summary>
        <dl className="grid gap-1 text-xs text-[var(--gw-muted)] sm:grid-cols-3" data-testid="basic-recording-committed-effect">
          <div><dt className="text-[var(--gw-muted)]">{t('step4.basic.effect_group')}</dt><dd className="break-all text-[var(--gw-text)]">{group.id} · {group.applied_revision}</dd></div>
          <div><dt className="text-[var(--gw-muted)]">{t('step4.basic.effect_record')}</dt><dd className="break-all text-[var(--gw-text)]">{committedEffect.record_id}</dd></div>
          <div><dt className="text-[var(--gw-muted)]">{t('step4.basic.effect_key')}</dt><dd className="break-all text-[var(--gw-text)]">{committedEffect.effect_key}</dd></div>
        </dl>
        </details>
      )}
    </section>
  );
};
