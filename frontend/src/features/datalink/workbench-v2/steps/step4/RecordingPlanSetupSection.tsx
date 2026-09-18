import * as React from 'react';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useSchemaOperationQuery, useStudioV2ConnectorCapabilitiesQuery, useStudioV2WorkspaceRecordingPlansQuery } from '@/hooks/datalink/useStudioV2WorkspaceRecordingPlans';
import { SafeQueryBoundary } from '@/utils/SafeQueryBoundary';
import type { StreamMode } from '@/types/recordingPlan';
import { recordingPlanMatchesScope, type RecordingPlanScope } from './recordingPlanMembership';
import { useRecordingPlanActions } from './useRecordingPlanActions';
import type { DbConnector } from '../../state/types';

export interface RecordingPlanSetupSectionProps {
  disabled?: boolean;
  connector?: DbConnector;
  scope?: RecordingPlanScope;
}

interface UnconfirmedNotice {
  operationId?: string;
}

/** Lets the operator query the durable operation instead of retrying blindly. */
const UnconfirmedOperationNotice: React.FC<{ notice: UnconfirmedNotice; testid: string; message: string }> = ({
  notice, testid, message,
}) => {
  const { t } = useTranslation('workbench-v2');
  const { data: checked, isError, isFetching, refetch } = useSchemaOperationQuery(notice.operationId);
  return (
    <div className="rounded-lg border border-amber-700 bg-amber-950/40 p-3 text-xs text-amber-300" data-testid={testid}>
      {message}
      {notice.operationId && (
        <span className="ml-1">({t('step4.recording_operation_id')}: {notice.operationId})</span>
      )}
      {notice.operationId && (
        <div className="mt-2 flex flex-wrap items-center gap-2">
          <button
            type="button"
            data-testid={`${testid}-check`}
            className="underline disabled:opacity-50"
            disabled={isFetching}
            onClick={() => void refetch()}
          >
            {t('step4.recording_operation_check')}
          </button>
          {isError && <span>{t('step4.recording_operation_check_failed')}</span>}
          {checked && (
            <span data-testid={`${testid}-result`}>
              {t('step4.recording_operation_status')}: {checked.status}
              {checked.next_action ? ` · ${checked.next_action}` : ''}
            </span>
          )}
        </div>
      )}
    </div>
  );
};

const RecordingPlanSetupSectionContent: React.FC<RecordingPlanSetupSectionProps> = ({
  disabled = false, connector, scope,
}) => {
  const { data: plans = [], isLoading, isError, refetch } = useStudioV2WorkspaceRecordingPlansQuery(true);
  const savedTarget = connector?.persisted === true && connector.save_state === 'saved' &&
    Boolean(connector.connector_id && connector.identity_revision) && connector.workspace_id === scope?.workspaceId;
  const { data: capabilities = [] } = useStudioV2ConnectorCapabilitiesQuery(savedTarget, connector?.kind);
  const { t } = useTranslation('workbench-v2');
  const [selectedPlanId, setSelectedPlanId] = useState('');
  const selectionScope = JSON.stringify([scope?.fingerprint, connector?.connector_id, connector?.identity_revision]);
  React.useEffect(() => { setSelectedPlanId(''); }, [selectionScope]);
  const selected = selectedPlanId ? plans.find(plan => plan.id === selectedPlanId) : (plans.length === 1 ? plans[0] : undefined);
  const activePlan = selected && recordingPlanMatchesScope(selected, scope, connector) ? selected : undefined;
  const selectedCapability = capabilities.find(capability => capability.kind === connector?.kind);
  const managedSchemaSupported = selectedCapability?.supports_managed_schema === true;
  const testWritesSupported = selectedCapability?.supports_test_writes === true;
  const canPrepare = savedTarget && Boolean(scope?.members.length) && !isError;
  const {
    planName, setPlanName, streamMode, setStreamMode, retentionDays, setRetentionDays,
    previewToken, testResult, actionError, applyUnconfirmed, applyOperation, testUnconfirmed, busy,
    handleCreateDefaultPlan, handlePreviewSchema, handleApplySchema, handleTestWrite,
  } = useRecordingPlanActions({ activePlan, connector, scope, disabled: disabled || !canPrepare,
    managedSchemaSupported, testWritesSupported, refetch });

  if (isError) return <div role="alert" className="text-xs text-amber-300">
    {t('step4.recording_plans_load_failed')}
    <button type="button" data-testid="recording-plan-retry" className="ml-2 underline" onClick={() => void refetch()}>{t('step4.recording_retry')}</button>
  </div>;

  if (isLoading) {
    return <div className="text-xs text-slate-400">載入記錄方案中…</div>;
  }

  return (
    <div
      className="bg-slate-800/40 border border-slate-700/60 rounded-xl p-5 space-y-4"
      data-testid="recording-plan-setup-section"
    >
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-sm font-semibold text-slate-200">
            遙測記錄方案與自動建表 (Recording Plan & Managed Schema)
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            依資料用途自動生成資料庫結構與索引，免手動撰寫 SQL DDL 與欄位配對。
          </p>
        </div>
        {activePlan && (
          <span className="text-[11px] bg-sky-500/10 text-sky-400 border border-sky-500/20 px-2.5 py-1 rounded-full font-mono">
            {t('step4.recording_selected_plan')}: {activePlan.name} (Rev: {activePlan.revision})
          </span>
        )}
      </div>

      {actionError && (
        <div className="rounded-lg border border-rose-800 bg-rose-950/40 p-3 text-xs text-rose-300" role="alert" data-testid="recording-plan-action-error">
          <div>{actionError.message}</div>
          {actionError.action && <div className="mt-1">{actionError.action}</div>}
          {actionError.requestId && (
            <div className="mt-1">{t('errors.request_id')}: {actionError.requestId}</div>
          )}
        </div>
      )}

      {!savedTarget && <p role="status" className="text-xs text-amber-300">{t('step4.recording_save_destination')}</p>}

      {applyUnconfirmed && (
        <UnconfirmedOperationNotice
          notice={applyUnconfirmed}
          testid="recording-plan-operation-unconfirmed"
          message={t('step4.recording_operation_unconfirmed')}
        />
      )}

      {applyOperation && (
        <div className="rounded-lg border border-slate-700 bg-slate-900/60 p-3 text-xs text-slate-300" data-testid="recording-plan-operation-status">
          <div>
            {t('step4.recording_operation_status')}: {applyOperation.status}
            {applyOperation.executed_statements > 0 && (
              <span className="ml-1">
                ({t('step4.recording_operation_executed')}: {applyOperation.executed_statements})
              </span>
            )}
          </div>
          {applyOperation.next_action && <div className="mt-1 text-amber-300">{applyOperation.next_action}</div>}
        </div>
      )}

      {plans.length > 0 && <div className="space-y-2">
        <label className="text-xs text-slate-300" htmlFor="recording-plan-selection">{t('step4.recording_select_plan')}</label>
        <select id="recording-plan-selection" data-testid="recording-plan-selection" value={activePlan?.id ?? ''}
          disabled={disabled || (busy !== null && busy !== 'preview')}
          onChange={event => setSelectedPlanId(event.target.value)} className="ml-2 rounded border border-slate-700 bg-slate-950 p-2 text-xs">
          <option value="">{t('step4.recording_choose_plan')}</option>
          {plans.map(plan => <option key={plan.id} value={plan.id} disabled={!recordingPlanMatchesScope(plan, scope, connector)}>{plan.name}</option>)}
        </select>
        {!activePlan && <p role="status" className="text-xs text-amber-300">{t('step4.recording_plan_scope_required')}</p>}
      </div>}
      {!scope?.members.length && <p role="status" className="text-xs text-amber-300">{t('step4.recording_members_required')}</p>}
      {plans.length === 0 ? (
        <div className="bg-slate-900/60 border border-slate-700/50 rounded-lg p-4 space-y-3">
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div>
              <label className="text-[11px] text-slate-400 block mb-1">方案名稱</label>
              <input
                type="text"
                className="w-full bg-slate-950 border border-slate-700 rounded px-2.5 py-1.5 text-xs text-slate-200"
                value={planName}
                onChange={(e) => setPlanName(e.target.value)}
                disabled={disabled || busy !== null}
              />
            </div>
            <div>
              <label className="text-[11px] text-slate-400 block mb-1">記錄用途主線</label>
              <select
                className="w-full bg-slate-950 border border-slate-700 rounded px-2.5 py-1.5 text-xs text-slate-200"
                value={streamMode}
                onChange={(e) => setStreamMode(e.target.value as StreamMode)}
                disabled={disabled || busy !== null}
              >
                <option value="window_summary">時序統計曲線 (Window Summary · 推薦)</option>
                <option value="usage_interval">累積用量差分 (Usage Interval · 電表水量)</option>
                <option value="raw_history">全明細保存 (Raw History · 完整樣本)</option>
                <option value="state_changes">狀態變化與警報 (State Changes)</option>
                <option value="latest_only">僅留最新值 (Latest Only)</option>
              </select>
            </div>
            <div>
              <label className="text-[11px] text-slate-400 block mb-1">摘要保留期限</label>
              <select
                className="w-full bg-slate-950 border border-slate-700 rounded px-2.5 py-1.5 text-xs text-slate-200"
                value={retentionDays}
                onChange={(e) => setRetentionDays(Number(e.target.value))}
                disabled={disabled || busy !== null}
              >
                <option value={90}>90 天 (近三個月)</option>
                <option value={180}>180 天 (半年)</option>
                <option value={365}>365 天 (一年 · 推薦)</option>
                <option value={730}>730 天 (兩年)</option>
              </select>
            </div>
          </div>

          <button
            type="button"
            className="bg-sky-600 hover:bg-sky-500 disabled:opacity-50 text-white text-xs px-4 py-2 rounded-lg font-medium transition-colors"
            onClick={handleCreateDefaultPlan}
            disabled={disabled || !canPrepare || busy !== null || !planName.trim()}
            data-testid="create-plan-btn"
          >
            {busy === 'create' ? '建立中…' : '建立記錄方案'}
          </button>
        </div>
      ) : activePlan ? (
        <div className="space-y-3">
          {(!managedSchemaSupported || !testWritesSupported) && (
            <div className="rounded-lg border border-amber-700/70 bg-amber-950/30 p-3 text-xs text-amber-300" data-testid="recording-plan-capability-warning">
              {t('step4.recording_operations_unavailable')}
            </div>
          )}
          <div className="flex flex-wrap items-center gap-3">
            <button
              type="button"
              className="bg-slate-700 hover:bg-slate-600 disabled:opacity-50 text-slate-200 text-xs px-3.5 py-2 rounded-lg font-medium transition-colors"
              onClick={handlePreviewSchema}
              disabled={disabled || !canPrepare || busy !== null}
              data-testid="preview-schema-btn"
            >
              {busy === 'preview' ? '產生中…' : '預覽資料庫 DDL'}
            </button>

            {previewToken && (
              <button
                type="button"
                className="bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white text-xs px-3.5 py-2 rounded-lg font-medium transition-colors"
                onClick={handleApplySchema}
                disabled={disabled || !canPrepare || !managedSchemaSupported || busy !== null || Boolean(applyUnconfirmed)}
                data-testid="apply-schema-btn"
              >
                {busy === 'apply' ? '建立中…' : '確認自動建表 (Managed DDL)'}
              </button>
            )}

            <button
              type="button"
              className="bg-violet-600 hover:bg-violet-500 disabled:opacity-50 text-white text-xs px-3.5 py-2 rounded-lg font-medium transition-colors"
              onClick={handleTestWrite}
              disabled={disabled || !canPrepare || !testWritesSupported || busy !== null || Boolean(testUnconfirmed)}
              data-testid="test-write-btn"
            >
              {busy === 'test' ? '試寫中…' : '執行真實試寫與回讀驗證'}
            </button>
          </div>

          {previewToken && (
            <div className="bg-slate-950/80 border border-slate-800 rounded-lg p-3 text-xs font-mono space-y-1.5 max-h-40 overflow-y-auto">
              <div className="text-slate-400 font-sans text-[11px] mb-1">
                預覽產生的 DDL 敘述句 ({previewToken.statements.length} 項)：
              </div>
              {previewToken.statements.map((stmt, i) => (
                <div key={i} className="text-sky-300">
                  {stmt}
                </div>
              ))}
            </div>
          )}

          {testResult && (
            <div
              className={`p-3 rounded-lg border text-xs ${testResult.status === 'written_verified'
                ? 'bg-emerald-500/10 border-emerald-500/30 text-emerald-300'
                : testResult.status === 'unknown'
                  ? 'bg-slate-500/10 border-slate-500/30 text-slate-300'
                  : testResult.status === 'written_unverified'
                    ? 'bg-amber-500/10 border-amber-500/30 text-amber-300'
                    : 'bg-rose-500/10 border-rose-500/30 text-rose-300'}`}
              data-testid="test-write-result"
            >
              <div className="font-semibold">
                {testResult.status === 'written_verified'
                  ? t('step4.recording_written_verified')
                  : testResult.status === 'written_unverified'
                    ? t('step4.recording_written_unverified')
                    : testResult.status === 'unknown'
                      ? t('step4.recording_write_unknown')
                      : t('step4.recording_write_failed')}
              </div>
              <div className="text-[11px] mt-0.5 opacity-80">
                寫入目標表: {testResult.table} · 測試記錄編號: {testResult.record_id} · 觀測時間: {testResult.observed_at}
              </div>
            </div>
          )}

          {testUnconfirmed && (
            <UnconfirmedOperationNotice
              notice={testUnconfirmed}
              testid="recording-plan-test-unconfirmed"
              message={t('step4.recording_test_unconfirmed')}
            />
          )}
        </div>
      ) : null}
    </div>
  );
};

export const RecordingPlanSetupSection: React.FC<RecordingPlanSetupSectionProps> = (props) => {
  return (
    <SafeQueryBoundary>
      <RecordingPlanSetupSectionContent {...props} />
    </SafeQueryBoundary>
  );
};
