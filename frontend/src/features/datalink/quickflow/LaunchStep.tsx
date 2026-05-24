import { useTranslation } from 'react-i18next';
import type { QuickFlowViewModel } from './quickflowTypes';

interface LaunchStepProps {
  model: QuickFlowViewModel;
}

export function LaunchStep({ model }: LaunchStepProps) {
  const { t } = useTranslation();

  return (
    <section style={{ display: 'grid', gap: 18 }}>
      <div>
        <h2 style={{ margin: 0, fontSize: 28 }}>{t('quickflow.sections.launch.title')}</h2>
        <p style={{ marginTop: 8, color: '#9ba6c9' }}>{t('quickflow.sections.launch.description')}</p>
      </div>
      <div style={{ display: 'grid', gap: 12 }}>
        {model.checks.map((check) => (
          <article
            key={check.id}
            style={{
              padding: '16px 18px',
              borderRadius: 20,
              background: 'rgba(10, 14, 22, 0.92)',
              border: '1px solid rgba(122,143,216,0.12)',
            }}
          >
            <strong style={{ display: 'block', fontSize: 17 }}>{check.title}</strong>
            <div style={{ marginTop: 8, color: '#9ba6c9' }}>{check.detail}</div>
          </article>
        ))}
      </div>
      <div style={{ padding: 18, borderRadius: 20, background: 'rgba(255, 92, 124, 0.08)', border: '1px solid rgba(255,92,124,0.18)' }}>
        <strong style={{ display: 'block', fontSize: 17 }}>{t('quickflow.sections.launch.blockerTitle')}</strong>
        <div style={{ marginTop: 8, color: '#ffcad5' }}>{model.blocker}</div>
      </div>
    </section>
  );
}
