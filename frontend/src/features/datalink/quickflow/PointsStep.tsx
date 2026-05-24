import { useTranslation } from 'react-i18next';
import type { QuickFlowViewModel } from './quickflowTypes';

interface PointsStepProps {
  model: QuickFlowViewModel;
}

export function PointsStep({ model }: PointsStepProps) {
  const { t } = useTranslation();

  return (
    <section style={{ display: 'grid', gap: 18 }}>
      <div>
        <h2 style={{ margin: 0, fontSize: 28 }}>{t('quickflow.sections.points.title')}</h2>
        <p style={{ marginTop: 8, color: '#9ba6c9' }}>{t('quickflow.sections.points.description')}</p>
      </div>
      <div style={{ padding: 20, borderRadius: 22, background: 'rgba(10, 14, 22, 0.92)', border: '1px solid rgba(122,143,216,0.12)' }}>
        <strong style={{ display: 'block', fontSize: 18 }}>{t('quickflow.sections.points.summaryTitle')}</strong>
        <div style={{ marginTop: 8, color: '#dfe6ff' }}>
          Holding Register {model.sourceRule.start_address}-{Number(model.sourceRule.start_address) + model.sourceRule.count - 1}
        </div>
        <div style={{ marginTop: 6, color: '#9ba6c9' }}>{t('quickflow.sections.points.summaryDetail', { prefix: model.sourceRule.naming_prefix })}</div>
      </div>
      <div style={{ display: 'grid', gap: 12 }}>
        {model.points.map((point) => (
          <article
            key={point.id}
            style={{
              padding: '16px 18px',
              borderRadius: 20,
              background: 'rgba(10, 14, 22, 0.9)',
              border: '1px solid rgba(122,143,216,0.12)',
            }}
          >
            <div style={{ display: 'flex', justifyContent: 'space-between', gap: 16 }}>
              <div>
                <strong style={{ display: 'block', fontSize: 17 }}>{point.label}</strong>
                <div style={{ marginTop: 6, color: '#9ba6c9' }}>Holding Register {point.address}</div>
              </div>
              <div style={{ fontSize: 24, fontWeight: 700, color: point.status === 'ok' ? '#7fe4b3' : '#ffd37f' }}>{point.displayValue}</div>
            </div>
          </article>
        ))}
      </div>
    </section>
  );
}
