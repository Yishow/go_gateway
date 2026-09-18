import { savedRecordingScope, savedRecordingPlanScopeFields } from '../../fixtures/recordingScope';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { createInstance } from 'i18next';
import { I18nextProvider } from 'react-i18next';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { RecordingPlanSetupSection } from '@/features/datalink/workbench-v2/steps/step4/RecordingPlanSetupSection';
import { studioV2DatalinkApi } from '@/services/studioV2Workspace';
import en from '@/i18n/locales/en/workbench-v2.json';
import zh from '@/i18n/locales/zh-TW/workbench-v2.json';
import { INITIAL_STATE } from '@/features/datalink/workbench-v2/state/useWorkbenchV2State';

vi.mock('@/services/studioV2Workspace', () => ({
  studioV2DatalinkApi: { get: vi.fn(), post: vi.fn() },
}));

const cases = [
  { language: 'en', schema: 'Automatic table creation is currently unavailable.',
    write: 'Recording test writes are currently unavailable.', wait: 'Wait until this operation is available.',
    unknown: 'The test-write result is unconfirmed. Check the database state before trying again.' },
  { language: 'zh-TW', schema: '自動建表目前尚未開放。',
    write: '記錄方案試寫目前尚未開放。', wait: '請等待此操作開放後再使用。',
    unknown: '試寫結果尚未確認；請先查明資料庫狀態，避免重複寫入。' },
] as const;

async function mountRecording(language: string, managedSchema = true, testWrites = true) {
  const i18n = createInstance();
  await i18n.init({
    lng: language, fallbackLng: false, defaultNS: 'workbench-v2',
    resources: { en: { 'workbench-v2': en }, 'zh-TW': { 'workbench-v2': zh } },
    interpolation: { escapeValue: false }, react: { useSuspense: false },
  });
  vi.mocked(studioV2DatalinkApi.get).mockImplementation(async (url) => ({
    data: { success: true, data: url?.endsWith('/capabilities') ? [{
      kind: 'sqlite', supported: true, supports_managed_schema: managedSchema,
      supports_test_writes: testWrites, supports_transactions: true, supports_receipts: true,
      supported_modes: ['managed_recording', 'custom_table'],
    }] : [{ ...savedRecordingPlanScopeFields, id: 'plan-1', name: 'Fixture recording plan', revision: 'rev-1' }] },
  }) as never);
  const client = new QueryClient({ defaultOptions: {
    queries: { retry: false }, mutations: { retry: 1, retryDelay: 0 },
  } });
  const view = render(
    <QueryClientProvider client={client}>
      <I18nextProvider i18n={i18n}>
        <RecordingPlanSetupSection scope={savedRecordingScope} connector={{ ...INITIAL_STATE.db.connector, kind: 'sqlite', workspace_id: 'workspace-a', connector_id: 'saved-sqlite', identity_revision: 'identity-1', setup_revision: 'setup-1', persisted: true, save_state: 'saved' }} />
      </I18nextProvider>
    </QueryClientProvider>,
  );
  await screen.findByTestId('preview-schema-btn');
  return view;
}

const preview = {
  token: 'token-1', operation_id: 'op-1', plan_id: 'plan-1', plan_revision: 'rev-1',
  workspace_id: 'workspace-a', workspace_revision: 'setup-1', connector_id: 'saved-sqlite', table_prefix: 'gw_record_',
  statements: ['CREATE TABLE sample (id TEXT)'], tables: [{ name: 'sample', action: 'create', columns: ['id'] }],
  digest: 'a'.repeat(64), expires_at: '2026-09-15T12:00:00Z', created_at: '2026-09-15T11:00:00Z',
};

describe.each(cases)('recording API to $language operator messages', (copy) => {
  beforeEach(() => vi.resetAllMocks());

  it.each(['schema-apply', 'test-write'] as const)('keeps %s 501 safe through the real service, hook and translation', async (operation) => {
    const code = operation === 'schema-apply'
      ? 'RECORDING_SCHEMA_NOT_IMPLEMENTED' : 'RECORDING_TEST_WRITE_NOT_IMPLEMENTED';
    const requestId = `req-recording-${operation}-full-0123456789`;
    vi.mocked(studioV2DatalinkApi.post).mockImplementation(async (url) => {
      if (url?.endsWith('/schema-preview')) return { data: { success: true, data: preview } } as never;
      throw { response: { status: 501, data: { success: false, error: {
        code, message: 'raw database diagnostic must stay hidden', retryable: false,
        action: 'wait_for_supported_operation', request_id: requestId,
      } } } };
    });
    await mountRecording(copy.language);
    if (operation === 'schema-apply') {
      fireEvent.click(screen.getByTestId('preview-schema-btn'));
      await screen.findByTestId('apply-schema-btn');
    }
    const buttonId = operation === 'schema-apply' ? 'apply-schema-btn' : 'test-write-btn';
    await waitFor(() => expect(screen.getByTestId(buttonId)).toBeEnabled());
    fireEvent.click(screen.getByTestId(buttonId));
    const alert = await screen.findByTestId('recording-plan-action-error');
    expect(alert).toHaveTextContent(operation === 'schema-apply' ? copy.schema : copy.write);
    expect(alert).toHaveTextContent(copy.wait);
    expect(alert).toHaveTextContent(requestId);
    expect(alert).not.toHaveTextContent(code);
    expect(alert).not.toHaveTextContent('wait_for_supported_operation');
    expect(alert).not.toHaveTextContent('raw database diagnostic');
    expect(screen.queryByTestId('test-write-result')).not.toBeInTheDocument();
    if (operation === 'schema-apply') expect(screen.getByText(preview.statements[0])).toBeInTheDocument();
    expect(vi.mocked(studioV2DatalinkApi.post).mock.calls.filter(([url]) => url?.endsWith(`/${operation}`))).toHaveLength(1);
  });

  it('retains an operation identity from an unfamiliar nominal 200 without a success card', async () => {
    vi.mocked(studioV2DatalinkApi.post).mockResolvedValue({ data: { success: true,
      data: { status: 'unknown', operation_id: 'op-unconfirmed-1' },
    } } as never);
    await mountRecording(copy.language);
    await waitFor(() => expect(screen.getByTestId('test-write-btn')).toBeEnabled());
    fireEvent.click(screen.getByTestId('test-write-btn'));
    const unknown = await screen.findByTestId('recording-plan-test-unconfirmed');
    expect(unknown).toHaveTextContent(copy.unknown);
    expect(unknown).toHaveTextContent('op-unconfirmed-1');
    expect(screen.queryByTestId('test-write-result')).not.toBeInTheDocument();
    expect(studioV2DatalinkApi.post).toHaveBeenCalledTimes(1);
  });
});
