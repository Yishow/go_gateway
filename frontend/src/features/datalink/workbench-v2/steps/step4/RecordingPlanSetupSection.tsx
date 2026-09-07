import * as React from 'react';
import { useState } from 'react';
import {
  useApplySchemaMutation,
  useCreateRecordingPlanMutation,
  usePreviewSchemaMutation,
  useStudioV2WorkspaceRecordingPlansQuery,
  useTestWritePlanMutation,
} from '@/hooks/datalink/useStudioV2WorkspaceRecordingPlans';
import { SafeQueryBoundary } from '@/utils/SafeQueryBoundary';
import type { RecordingPlan, SchemaPreviewToken, StreamMode, TestWriteResult } from '@/types/recordingPlan';

export interface RecordingPlanSetupSectionProps {
  deviceId?: string;
  disabled?: boolean;
}

const RecordingPlanSetupSectionContent: React.FC<RecordingPlanSetupSectionProps> = ({
  deviceId,
  disabled = false,
}) => {
  const { data: plans = [], isLoading, refetch } = useStudioV2WorkspaceRecordingPlansQuery(true, deviceId);
  const createPlanMutation = useCreateRecordingPlanMutation();
  const previewSchemaMutation = usePreviewSchemaMutation();
  const applySchemaMutation = useApplySchemaMutation();
  const testWriteMutation = useTestWritePlanMutation();

  const [planName, setPlanName] = useState('標準生產設備記錄方案');
  const [streamMode, setStreamMode] = useState<StreamMode>('window_summary');
  const [retentionDays, setRetentionDays] = useState<number>(365);
  const [previewToken, setPreviewToken] = useState<SchemaPreviewToken | null>(null);
  const [testResult, setTestResult] = useState<TestWriteResult | null>(null);

  const activePlan = plans[0] as RecordingPlan | undefined;

  const handleCreateDefaultPlan = async () => {
    if (!deviceId) return;
    const newPlan: Partial<RecordingPlan> = {
      name: planName,
      timezone: 'Asia/Taipei',
      members: [
        {
          member_id: 'member-primary',
          measurement_id: 'meas-primary',
          equipment_id: deviceId,
          name: '主設備測量項目',
        },
      ],
      retention: {
        raw_days: 30,
        summary_days: retentionDays,
        events_days: 90,
        correction_horizon_hours: 24,
      },
      streams: [
        {
          stream_id: 'stream-primary',
          measurement_id: 'meas-primary',
          mode: streamMode,
          raw_policy: 'every_sample',
        },
      ],
      destinations: [
        {
          destination_id: 'dest-default-db',
          connector_id: 'default-sqlite',
          table_prefix: 'gw_record_',
        },
      ],
    };
    await createPlanMutation.mutateAsync(newPlan);
    refetch();
  };

  const handlePreviewSchema = async () => {
    if (!activePlan) return;
    const token = await previewSchemaMutation.mutateAsync({
      plan_id: activePlan.id,
      connector_id: 'default-sqlite',
      dialect: 'sqlite',
      table_prefix: 'gw_record_',
    });
    setPreviewToken(token);
  };

  const handleApplySchema = async () => {
    if (!previewToken) return;
    await applySchemaMutation.mutateAsync(previewToken.token);
    setPreviewToken(null);
    refetch();
  };

  const handleTestWrite = async () => {
    if (!activePlan) return;
    const res = await testWriteMutation.mutateAsync({
      plan_id: activePlan.id,
      table_prefix: 'gw_record_',
    });
    setTestResult(res);
  };

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
            方案已就緒: {activePlan.name} (Rev: {activePlan.revision})
          </span>
        )}
      </div>

      {!activePlan ? (
        <div className="bg-slate-900/60 border border-slate-700/50 rounded-lg p-4 space-y-3">
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div>
              <label className="text-[11px] text-slate-400 block mb-1">方案名稱</label>
              <input
                type="text"
                className="w-full bg-slate-950 border border-slate-700 rounded px-2.5 py-1.5 text-xs text-slate-200"
                value={planName}
                onChange={(e) => setPlanName(e.target.value)}
                disabled={disabled}
              />
            </div>
            <div>
              <label className="text-[11px] text-slate-400 block mb-1">記錄用途主線</label>
              <select
                className="w-full bg-slate-950 border border-slate-700 rounded px-2.5 py-1.5 text-xs text-slate-200"
                value={streamMode}
                onChange={(e) => setStreamMode(e.target.value as StreamMode)}
                disabled={disabled}
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
                disabled={disabled}
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
            disabled={disabled || createPlanMutation.isPending}
            data-testid="create-plan-btn"
          >
            {createPlanMutation.isPending ? '建立中…' : '建立記錄方案'}
          </button>
        </div>
      ) : (
        <div className="space-y-3">
          <div className="flex flex-wrap items-center gap-3">
            <button
              type="button"
              className="bg-slate-700 hover:bg-slate-600 disabled:opacity-50 text-slate-200 text-xs px-3.5 py-2 rounded-lg font-medium transition-colors"
              onClick={handlePreviewSchema}
              disabled={disabled || previewSchemaMutation.isPending}
              data-testid="preview-schema-btn"
            >
              {previewSchemaMutation.isPending ? '產生中…' : '預覽資料庫 DDL'}
            </button>

            {previewToken && (
              <button
                type="button"
                className="bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white text-xs px-3.5 py-2 rounded-lg font-medium transition-colors"
                onClick={handleApplySchema}
                disabled={disabled || applySchemaMutation.isPending}
                data-testid="apply-schema-btn"
              >
                {applySchemaMutation.isPending ? '建立中…' : '確認自動建表 (Managed DDL)'}
              </button>
            )}

            <button
              type="button"
              className="bg-violet-600 hover:bg-violet-500 disabled:opacity-50 text-white text-xs px-3.5 py-2 rounded-lg font-medium transition-colors"
              onClick={handleTestWrite}
              disabled={disabled || testWriteMutation.isPending}
              data-testid="test-write-btn"
            >
              {testWriteMutation.isPending ? '試寫中…' : '執行真實試寫與回讀驗證'}
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
              className={`p-3 rounded-lg border text-xs ${
                testResult.status === 'success'
                  ? 'bg-emerald-500/10 border-emerald-500/30 text-emerald-300'
                  : 'bg-rose-500/10 border-rose-500/30 text-rose-300'
              }`}
              data-testid="test-write-result"
            >
              <div className="font-semibold">
                試寫狀態: {testResult.status === 'success' ? '成功 (已驗證回讀且清理乾淨)' : '失敗'}
              </div>
              <div className="text-[11px] mt-0.5 opacity-80">
                寫入目標表: {testResult.table} · 測試記錄編號: {testResult.record_id} · 觀測時間: {testResult.observed_at}
              </div>
            </div>
          )}
        </div>
      )}
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
