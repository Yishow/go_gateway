import { useTranslation } from 'react-i18next';
import type { QuickFlowStepDefinition, QuickFlowStepId } from './quickflowTypes';

interface StepRailProps {
  steps: QuickFlowStepDefinition[];
  activeStep: QuickFlowStepId;
  onStepChange: (step: QuickFlowStepId) => void;
}

const railStyle = {
  display: 'grid',
  gridTemplateColumns: 'repeat(4, minmax(0, 1fr))',
  gap: 12,
} as const;

export function StepRail({ steps, activeStep, onStepChange }: StepRailProps) {
  const { t } = useTranslation();

  return (
    <div style={railStyle}>
      {steps.map((step, index) => {
        const active = step.id === activeStep;
        return (
          <button
            key={step.id}
            type="button"
            onClick={() => onStepChange(step.id)}
            aria-pressed={active}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: 12,
              padding: '14px 16px',
              borderRadius: 18,
              border: active ? '1px solid rgba(153, 120, 255, 0.46)' : '1px solid rgba(122, 143, 216, 0.14)',
              background: active
                ? 'linear-gradient(135deg, rgba(93,124,255,0.2), rgba(155,92,246,0.18))'
                : 'rgba(11, 16, 26, 0.9)',
              color: '#f5f7ff',
              cursor: 'pointer',
              textAlign: 'left',
            }}
          >
            <span
              style={{
                display: 'grid',
                placeItems: 'center',
                width: 34,
                height: 34,
                borderRadius: 12,
                background: active ? 'linear-gradient(135deg, #5d7cff, #9b5cf6)' : 'rgba(255,255,255,0.08)',
                fontWeight: 700,
              }}
            >
              {index + 1}
            </span>
            <span>
              <strong style={{ display: 'block', fontSize: 15 }}>{t(step.labelKey)}</strong>
              <span style={{ color: '#9ba6c9', fontSize: 12 }}>{t(step.descriptionKey)}</span>
            </span>
          </button>
        );
      })}
    </div>
  );
}
