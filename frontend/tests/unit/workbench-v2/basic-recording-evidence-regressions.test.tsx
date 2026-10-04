import { createInstance } from 'i18next';
import { I18nextProvider } from 'react-i18next';
import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { BasicRecordingEvidence } from '@/features/datalink/workbench-v2/steps/step4/BasicRecordingEvidence';
import en from '@/i18n/locales/en/workbench-v2.json';
import { savedGroup } from '../../fixtures/writeGroupState';

describe('Basic recording first-bucket evidence', () => {
  it('stops claiming the first complete bucket is pending after a current receipt confirms SQL', async () => {
    const i18n = createInstance();
    await i18n.init({ lng: 'en', resources: { en: { 'workbench-v2': en } }, interpolation: { escapeValue: false } });
    const group = savedGroup({ applied_revision: 'applied-current' });
    const initial = <BasicRecordingEvidence intervalSeconds={60} deliveryError={false} />;
    const { rerender } = render(<I18nextProvider i18n={i18n}>{initial}</I18nextProvider>);
    expect(screen.getByText('The first complete UTC snapshot is still pending until its 60-second bucket closes.')).toBeInTheDocument();

    rerender(
      <I18nextProvider i18n={i18n}>
        <BasicRecordingEvidence intervalSeconds={60} deliveryError={false} group={group} committedEffect={{
          group_revision: 'applied-current', connector_revision: 'connector-current',
          record_id: 'formal-record', effect_key: 'formal-effect', payload_digest: 'formal-digest',
          committed_at: '2026-10-04T01:06:01Z',
        }} />
      </I18nextProvider>,
    );
    expect(screen.getByText('SQL commit confirmed')).toBeInTheDocument();
    expect(screen.getByText('formal-record')).toBeInTheDocument();
    expect(screen.queryByText(/The first complete UTC snapshot is still pending/)).not.toBeInTheDocument();
  });
});
