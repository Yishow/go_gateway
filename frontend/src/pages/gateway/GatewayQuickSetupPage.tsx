import { useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import {
  ArrowLeft,
  CheckCircle2,
  ChevronRight,
  ClipboardCheck,
  Network,
  Save,
  Server,
  Shield,
  Terminal,
  Zap,
} from 'lucide-react';
import { gatewayAdapter, type QuickDraft } from '../../features/gateway/gatewayAdapter';
import {
  setGatewayQuickStep,
  updateGatewayQuickDraft,
  updateGatewayQuickWizardDraft,
  useGatewayDraftStore,
} from '../../features/gateway/gatewayDraftStore';
import { useTestAPI } from '../../services/api';

const STEPS = [
  'Step1 設備連線',
  'Step2 路由/來源規劃',
  'Step3 Tag/安全',
  'Step4 驗證',
  'Step5 提交',
] as const;

const isPositiveInteger = (value: number | undefined) =>
  typeof value === 'number' && Number.isInteger(value) && value > 0;

export default function GatewayQuickSetupPage() {
  const currentStep = useGatewayDraftStore((state) => state.quickStep);
  const draft = useGatewayDraftStore((state) => state.quickDraft);
  const wizardDraft = useGatewayDraftStore((state) => state.quickWizardDraft);
  const [stepErrors, setStepErrors] = useState<Partial<Record<number, string[]>>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitMessage, setSubmitMessage] = useState<string | null>(null);
  const [submitError, setSubmitError] = useState<string | null>(null);
  const payload = useMemo(() => gatewayAdapter.quickToPayload(draft), [draft]);
  const { connect } = useTestAPI();

  const validateStep = (step: number): string[] => {
    const errors: string[] = [];

    if (step === 1) {
      if (!draft.protocol?.trim()) errors.push('請選擇通訊協定。');
      if (!draft.host?.trim()) errors.push('請輸入主機位置。');
      if (!isPositiveInteger(draft.port)) errors.push('通訊埠需為正整數。');
      if (draft.protocol.includes('modbus') && !isPositiveInteger(draft.unitID)) {
        errors.push('Modbus 站號需為正整數。');
      }
      if (draft.protocol.includes('fatek') && !isPositiveInteger(draft.station)) {
        errors.push('Fatek 站號需為正整數。');
      }
    }

    if (step === 2) {
      if (!draft.route?.trim()) errors.push('請輸入 API 路徑。');
      if (draft.route && !draft.route.startsWith('/')) errors.push('API 路徑需以 "/" 開頭。');
      if (!wizardDraft.sourcePlan.trim()) errors.push('請輸入來源規劃。');
    }

    if (step === 3) {
      if (!wizardDraft.tagName.trim()) errors.push('請輸入 Tag 名稱。');
      if (!wizardDraft.securityReviewed) errors.push('請確認已完成安全檢查。');
    }

    if (step === 4 && !wizardDraft.validationConfirmed) {
      errors.push('請勾選「已完成驗證」才能進入提交。');
    }

    if (step === 5 && !wizardDraft.submitConfirmed) {
      errors.push('請勾選提交確認。');
    }

    return errors;
  };

  const validateAllSteps = (): Partial<Record<number, string[]>> => {
    const nextErrors: Partial<Record<number, string[]>> = {};
    [1, 2, 3, 4, 5].forEach((step) => {
      const errors = validateStep(step);
      if (errors.length > 0) nextErrors[step] = errors;
    });
    return nextErrors;
  };

  const allRequiredStepsValid = useMemo(
    () => [1, 2, 3, 4].every((step) => validateStep(step).length === 0),
    [draft, wizardDraft]
  );

  const handleQuickDraftChange = (key: keyof QuickDraft, value: string | number | boolean) => {
    updateGatewayQuickDraft({ [key]: value } as Partial<QuickDraft>);
    setStepErrors((prev) => ({ ...prev, [currentStep]: [] }));
  };

  const handleWizardDraftChange = (key: keyof typeof wizardDraft, value: string | boolean) => {
    updateGatewayQuickWizardDraft({ [key]: value } as Partial<typeof wizardDraft>);
    setStepErrors((prev) => ({ ...prev, [currentStep]: [] }));
  };

  const handleNext = () => {
    const errors = validateStep(currentStep);
    if (errors.length > 0) {
      setStepErrors((prev) => ({ ...prev, [currentStep]: errors }));
      return;
    }
    setStepErrors((prev) => ({ ...prev, [currentStep]: [] }));
    setGatewayQuickStep(currentStep + 1);
  };

  const handleBack = () => {
    setGatewayQuickStep(currentStep - 1);
  };

  const handleSubmit = async () => {
    const nextErrors = validateAllSteps();
    if (Object.keys(nextErrors).length > 0) {
      setStepErrors(nextErrors);
      const firstInvalidStep = [1, 2, 3, 4, 5].find((step) => nextErrors[step]?.length);
      if (firstInvalidStep) setGatewayQuickStep(firstInvalidStep);
      return;
    }

    setIsSubmitting(true);
    setSubmitError(null);
    setSubmitMessage(null);

    try {
      const response = await connect(payload.protocol, payload.config);
      setSubmitMessage(`連線測試成功：${response.connection_id} (${response.status})`);
    } catch (error) {
      const message = error instanceof Error ? error.message : '提交失敗';
      setSubmitError(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen bg-[#0B1120] text-slate-300 font-sans selection:bg-blue-500/30 relative">
      <div className="absolute inset-0 bg-[radial-gradient(#1e293b_1px,transparent_1px)] [background-size:16px_16px] opacity-10 pointer-events-none"></div>
      <div className="absolute top-0 inset-x-0 h-px bg-gradient-to-r from-transparent via-blue-500/50 to-transparent"></div>

      <main className="mx-auto max-w-[1400px] px-6 py-10 flex flex-col min-h-screen relative z-10">
        <header className="mb-8 flex flex-col md:flex-row md:items-end justify-between gap-6">
          <div>
            <Link
              to="/gateway/entry"
              className="inline-flex items-center text-xs font-semibold tracking-wider text-slate-400 hover:text-white transition-colors mb-6 uppercase bg-slate-800/50 px-3 py-1.5 rounded-md border border-slate-700/50"
            >
              <ArrowLeft className="mr-2 h-3.5 w-3.5" />
              返回模式選擇
            </Link>
            <div className="flex items-center gap-4">
              <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-blue-500/20 to-blue-600/10 border border-blue-500/20 shadow-inner">
                <Zap className="h-6 w-6 text-blue-400" />
              </div>
              <div>
                <h1 className="text-3xl font-extrabold tracking-tight text-white mb-1 flex items-center">
                  快速設定配置
                  <span className="ml-3 inline-flex items-center rounded-full bg-blue-500/10 px-2.5 py-0.5 text-xs font-medium text-blue-400 border border-blue-500/20">
                    Quick Setup
                  </span>
                </h1>
                <p className="text-sm text-slate-400">五步驟狀態機：連線、規劃、Tag/安全、驗證與提交。</p>
              </div>
            </div>
          </div>

          <div className="flex items-center gap-3 text-sm text-slate-400 bg-slate-900/50 px-4 py-2 rounded-lg border border-slate-800 shadow-sm">
            <span className="flex items-center gap-1.5">
              <CheckCircle2 className="w-4 h-4 text-emerald-500" /> 系統就緒
            </span>
            <span className="w-px h-4 bg-slate-700"></span>
            <span className="font-mono text-xs">
              {STEPS[currentStep - 1]}
            </span>
          </div>
        </header>

        <div className="mb-6 grid grid-cols-1 gap-2 md:grid-cols-5">
          {STEPS.map((step, index) => {
            const stepNumber = index + 1;
            const isActive = stepNumber === currentStep;
            return (
              <div
                key={step}
                className={`rounded-md border px-3 py-2 text-xs font-semibold ${
                  isActive
                    ? 'border-blue-500 bg-blue-500/15 text-blue-300'
                    : 'border-slate-700 bg-slate-900/60 text-slate-400'
                }`}
              >
                {step}
              </div>
            );
          })}
        </div>

        <div className="grid grid-cols-1 xl:grid-cols-12 gap-8 flex-1 items-start">
          <div className="xl:col-span-7 space-y-6">
            <section className="rounded-xl border border-slate-800 bg-[#111827]/80 p-7 shadow-lg backdrop-blur-md">
              {currentStep === 1 ? (
                <div className="space-y-6">
                  <div className="flex items-center gap-3">
                    <Server className="h-5 w-5 text-blue-400" />
                    <h2 className="text-base font-bold text-white">Step1 設備連線</h2>
                  </div>
                  <div className="grid gap-6">
                    <div>
                      <label className="block text-xs font-semibold tracking-wider text-slate-400 uppercase mb-2">通訊協定</label>
                      <div className="relative">
                        <select
                          data-testid="protocol-select"
                          value={draft.protocol}
                          onChange={(e) => handleQuickDraftChange('protocol', e.target.value)}
                          className="w-full appearance-none rounded-lg border border-slate-700 bg-slate-900 px-4 py-3 text-sm text-white focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
                        >
                          <option value="modbus-tcp">Modbus TCP</option>
                          <option value="modbus-rtu">Modbus RTU</option>
                          <option value="fatek-tcp">Fatek TCP</option>
                          <option value="mc-tcp">Mitsubishi MC</option>
                        </select>
                        <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-4 text-slate-500">
                          <ChevronRight className="h-4 w-4 rotate-90" />
                        </div>
                      </div>
                    </div>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <div>
                        <label className="block text-xs font-semibold tracking-wider text-slate-400 uppercase mb-2">主機位置 (Host)</label>
                        <input
                          type="text"
                          data-testid="host-input"
                          value={draft.host || ''}
                          onChange={(e) => handleQuickDraftChange('host', e.target.value)}
                          className="w-full rounded-lg border border-slate-700 bg-slate-900 px-4 py-3 text-sm text-white placeholder-slate-600 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none font-mono"
                          placeholder="192.168.1.100"
                        />
                      </div>
                      <div>
                        <label className="block text-xs font-semibold tracking-wider text-slate-400 uppercase mb-2">通訊埠 (Port)</label>
                        <input
                          type="number"
                          data-testid="port-input"
                          value={draft.port || ''}
                          onChange={(e) => handleQuickDraftChange('port', parseInt(e.target.value, 10))}
                          className="w-full rounded-lg border border-slate-700 bg-slate-900 px-4 py-3 text-sm text-white focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none font-mono"
                        />
                      </div>
                    </div>
                    {draft.protocol.includes('modbus') && (
                      <div>
                        <label className="block text-xs font-semibold tracking-wider text-slate-400 uppercase mb-2">站號 (Unit ID)</label>
                        <input
                          type="number"
                          data-testid="unit-id-input"
                          value={draft.unitID || ''}
                          onChange={(e) => handleQuickDraftChange('unitID', parseInt(e.target.value, 10))}
                          className="w-full max-w-[200px] rounded-lg border border-slate-700 bg-slate-900 px-4 py-3 text-sm text-white focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none font-mono"
                        />
                      </div>
                    )}
                    {draft.protocol.includes('fatek') && (
                      <div>
                        <label className="block text-xs font-semibold tracking-wider text-slate-400 uppercase mb-2">站號 (Station)</label>
                        <input
                          type="number"
                          data-testid="station-input"
                          value={draft.station || ''}
                          onChange={(e) => handleQuickDraftChange('station', parseInt(e.target.value, 10))}
                          className="w-full max-w-[200px] rounded-lg border border-slate-700 bg-slate-900 px-4 py-3 text-sm text-white focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none font-mono"
                        />
                      </div>
                    )}
                  </div>
                </div>
              ) : null}

              {currentStep === 2 ? (
                <div className="space-y-6">
                  <div className="flex items-center gap-3">
                    <Network className="h-5 w-5 text-emerald-400" />
                    <h2 className="text-base font-bold text-white">Step2 路由/來源規劃</h2>
                  </div>
                  <div>
                    <label className="block text-xs font-semibold tracking-wider text-slate-400 uppercase mb-2">API 路徑</label>
                    <input
                      type="text"
                      data-testid="route-input"
                      value={draft.route || ''}
                      onChange={(e) => handleQuickDraftChange('route', e.target.value)}
                      className="w-full rounded-lg border border-slate-700 bg-slate-900 py-3 px-4 text-sm text-white placeholder-slate-600 focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 focus:outline-none font-mono"
                      placeholder="/api/v1/data"
                    />
                  </div>
                  <div>
                    <label className="block text-xs font-semibold tracking-wider text-slate-400 uppercase mb-2">來源規劃 (Source Plan)</label>
                    <input
                      type="text"
                      data-testid="source-input"
                      value={wizardDraft.sourcePlan}
                      onChange={(e) => handleWizardDraftChange('sourcePlan', e.target.value)}
                      className="w-full rounded-lg border border-slate-700 bg-slate-900 py-3 px-4 text-sm text-white placeholder-slate-600 focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 focus:outline-none"
                      placeholder="line-a-source"
                    />
                  </div>
                </div>
              ) : null}

              {currentStep === 3 ? (
                <div className="space-y-6">
                  <div className="flex items-center gap-3">
                    <Shield className="h-5 w-5 text-purple-400" />
                    <h2 className="text-base font-bold text-white">Step3 Tag/安全</h2>
                  </div>
                  <div>
                    <label className="block text-xs font-semibold tracking-wider text-slate-400 uppercase mb-2">Tag 名稱</label>
                    <input
                      type="text"
                      data-testid="tag-input"
                      value={wizardDraft.tagName}
                      onChange={(e) => handleWizardDraftChange('tagName', e.target.value)}
                      className="w-full rounded-lg border border-slate-700 bg-slate-900 py-3 px-4 text-sm text-white placeholder-slate-600 focus:border-purple-500 focus:ring-1 focus:ring-purple-500 focus:outline-none"
                      placeholder="temperature"
                    />
                  </div>
                  <label className="flex items-center gap-3 rounded-lg border border-slate-700/60 bg-slate-900/50 p-4">
                    <input
                      type="checkbox"
                      data-testid="auth-checkbox"
                      checked={draft.auth || false}
                      onChange={(e) => handleQuickDraftChange('auth', e.target.checked)}
                      className="h-4 w-4 rounded border-slate-600 bg-slate-800 text-purple-500 focus:ring-purple-500"
                    />
                    <span className="text-sm text-slate-300">啟用基本身分驗證 (Basic Auth)</span>
                  </label>
                  <label className="flex items-center gap-3 rounded-lg border border-slate-700/60 bg-slate-900/50 p-4">
                    <input
                      type="checkbox"
                      data-testid="security-review-checkbox"
                      checked={wizardDraft.securityReviewed}
                      onChange={(e) => handleWizardDraftChange('securityReviewed', e.target.checked)}
                      className="h-4 w-4 rounded border-slate-600 bg-slate-800 text-purple-500 focus:ring-purple-500"
                    />
                    <span className="text-sm text-slate-300">我已確認 Tag 命名與安全策略。</span>
                  </label>
                </div>
              ) : null}

              {currentStep === 4 ? (
                <div className="space-y-6">
                  <div className="flex items-center gap-3">
                    <ClipboardCheck className="h-5 w-5 text-cyan-400" />
                    <h2 className="text-base font-bold text-white">Step4 驗證</h2>
                  </div>
                  <div className="rounded-lg border border-slate-700 bg-slate-900/60 p-4 text-sm">
                    <p>設備：{draft.protocol} @ {draft.host}:{draft.port}</p>
                    <p>路徑：{draft.route}</p>
                    <p>來源：{wizardDraft.sourcePlan}</p>
                    <p>Tag：{wizardDraft.tagName || '未設定'}</p>
                  </div>
                  <label className="flex items-center gap-3 rounded-lg border border-slate-700/60 bg-slate-900/50 p-4">
                    <input
                      type="checkbox"
                      data-testid="validation-confirm-checkbox"
                      checked={wizardDraft.validationConfirmed}
                      onChange={(e) => handleWizardDraftChange('validationConfirmed', e.target.checked)}
                      className="h-4 w-4 rounded border-slate-600 bg-slate-800 text-cyan-500 focus:ring-cyan-500"
                    />
                    <span className="text-sm text-slate-300">已完成欄位驗證與參數核對。</span>
                  </label>
                </div>
              ) : null}

              {currentStep === 5 ? (
                <div className="space-y-6">
                  <div className="flex items-center gap-3">
                    <Save className="h-5 w-5 text-amber-400" />
                    <h2 className="text-base font-bold text-white">Step5 提交</h2>
                  </div>
                  <label className="flex items-center gap-3 rounded-lg border border-slate-700/60 bg-slate-900/50 p-4">
                    <input
                      type="checkbox"
                      data-testid="submit-confirm-checkbox"
                      checked={wizardDraft.submitConfirmed}
                      onChange={(e) => handleWizardDraftChange('submitConfirmed', e.target.checked)}
                      className="h-4 w-4 rounded border-slate-600 bg-slate-800 text-amber-500 focus:ring-amber-500"
                    />
                    <span className="text-sm text-slate-300">確認以上設定可送出至 Gateway 測試流程。</span>
                  </label>
                </div>
              ) : null}

              {stepErrors[currentStep]?.length ? (
                <div
                  data-testid="step-error-message"
                  className="mt-6 rounded-md border border-red-500/30 bg-red-500/10 px-3 py-2 text-xs text-red-200 space-y-1"
                >
                  {stepErrors[currentStep]?.map((message) => (
                    <p key={message}>{message}</p>
                  ))}
                </div>
              ) : null}

              <div className="mt-8 flex items-center justify-between gap-3">
                <button
                  type="button"
                  data-testid="step-back-btn"
                  onClick={handleBack}
                  disabled={currentStep === 1}
                  className="rounded-md border border-slate-700 bg-slate-900 px-4 py-2 text-sm text-slate-200 disabled:cursor-not-allowed disabled:opacity-50"
                >
                  上一步
                </button>
                {currentStep < 5 ? (
                  <button
                    type="button"
                    data-testid="step-next-btn"
                    onClick={handleNext}
                    className="rounded-md border border-blue-500 bg-blue-600 px-4 py-2 text-sm font-semibold text-white hover:bg-blue-500"
                  >
                    下一步
                  </button>
                ) : (
                  <span className="text-xs text-slate-400">已到最後一步，請於右側提交。</span>
                )}
              </div>
            </section>
          </div>

          <div className="xl:col-span-5 xl:sticky xl:top-8 h-[calc(100vh-8rem)] min-h-[600px] flex flex-col">
            <section className="flex flex-col h-full rounded-xl border border-slate-700 bg-[#090D14] shadow-2xl overflow-hidden ring-1 ring-white/5 relative">
              <div className="flex items-center justify-between px-4 py-3 bg-[#111827] border-b border-slate-800">
                <div className="flex items-center gap-2">
                  <div className="h-3 w-3 rounded-full bg-[#FF5F56] border border-black/20"></div>
                  <div className="h-3 w-3 rounded-full bg-[#FFBD2E] border border-black/20"></div>
                  <div className="h-3 w-3 rounded-full bg-[#27C93F] border border-black/20"></div>
                </div>
                <div className="flex items-center gap-2 text-xs font-mono font-medium text-slate-400 px-3 py-1 bg-slate-900 rounded-md border border-slate-800">
                  <Terminal className="h-3.5 w-3.5 text-blue-400" />
                  <span>Payload 預覽 (供驗證)</span>
                </div>
                <div className="w-10"></div>
              </div>
              <div className="flex-1 overflow-auto p-5 font-mono text-xs leading-relaxed">
                <div data-testid="payload-preview" className="hidden">
                  {JSON.stringify(payload, null, 2)}
                </div>
                <pre className="whitespace-pre-wrap break-all text-slate-200">
                  {JSON.stringify(payload, null, 2)}
                </pre>
              </div>
              <div className="p-5 border-t border-slate-800 bg-[#111827] shrink-0">
                <button
                  data-testid="save-btn"
                  disabled={isSubmitting || currentStep !== 5 || !allRequiredStepsValid || !wizardDraft.submitConfirmed}
                  className="group relative flex w-full items-center justify-center gap-2 rounded-lg bg-blue-600 px-6 py-3.5 font-semibold text-white transition-all hover:bg-blue-500 disabled:cursor-not-allowed disabled:opacity-60"
                  onClick={handleSubmit}
                >
                  <Save className="h-5 w-5" />
                  <span>{isSubmitting ? '提交中...' : '部署設定 (Deploy)'}</span>
                </button>
                {submitMessage ? (
                  <p className="mt-3 rounded-md border border-emerald-500/30 bg-emerald-500/10 px-3 py-2 text-xs text-emerald-200">
                    {submitMessage}
                  </p>
                ) : null}
                {submitError ? (
                  <p className="mt-3 rounded-md border border-red-500/30 bg-red-500/10 px-3 py-2 text-xs text-red-200">
                    提交失敗：{submitError}
                  </p>
                ) : null}
                <p className="text-[11px] text-slate-500 mt-4 text-center font-mono">
                  系統將依據上述 Payload 結構發送至 /api/v1/test/connect
                </p>
              </div>
            </section>
          </div>
        </div>
      </main>
    </div>
  );
}
