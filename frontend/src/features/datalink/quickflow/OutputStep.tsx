import { useTranslation } from 'react-i18next';
import type { QuickFlowViewModel } from './quickflowTypes';

interface OutputStepProps {
  model: QuickFlowViewModel;
}

export function OutputStep({ model }: OutputStepProps) {
  const { t } = useTranslation();

  return (
    <section style={{ display: 'grid', gap: 18 }}>
      <div>
        <h2 style={{ margin: 0, fontSize: 28 }}>{t('quickflow.sections.output.title')}</h2>
        <p style={{ marginTop: 8, color: '#9ba6c9' }}>{t('quickflow.sections.output.description')}</p>
      </div>
      <div style={{ display: 'grid', gap: 12 }}>
        {model.outputs.map((target) => (
          <article
            key={target.id}
            style={{
              padding: '16px 18px',
              borderRadius: 20,
              background: 'rgba(10, 14, 22, 0.92)',
              border: '1px solid rgba(122,143,216,0.12)',
            }}
          >
            <div style={{ display: 'flex', justifyContent: 'space-between', gap: 12 }}>
              <div>
                <strong style={{ display: 'block', fontSize: 18 }}>{target.title}</strong>
                <div style={{ marginTop: 6, color: '#9ba6c9' }}>{target.detail}</div>
              </div>
              <span
                style={{
                  alignSelf: 'flex-start',
                  padding: '6px 10px',
                  borderRadius: 999,
                  background: target.status === 'ok' ? 'rgba(61,214,140,0.12)' : 'rgba(255,92,124,0.12)',
                  color: target.status === 'ok' ? '#bafad7' : '#ffbac8',
                  fontSize: 12,
                }}
              >
                {target.status === 'ok' ? t('quickflow.badges.ready') : t('quickflow.badges.blocked')}
              </span>
            </div>
          </article>
        ))}
      </div>
    </section>
  );
}
