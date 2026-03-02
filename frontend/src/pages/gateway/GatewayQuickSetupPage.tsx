import { useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { Server, Network, Shield, Code, ArrowLeft, Terminal, Save, CheckCircle2, ChevronRight, Zap } from 'lucide-react';
import { gatewayAdapter, type QuickDraft } from '../../features/gateway/gatewayAdapter';
import { useTestAPI } from '../../services/api';

export default function GatewayQuickSetupPage() {
  const [draft, setDraft] = useState<QuickDraft>({
    protocol: 'modbus-tcp',
    host: '192.168.1.100',
    port: 502,
    unitID: 1,
    route: '/api/v1/data',
    auth: false,
  });

  const payload = useMemo(() => gatewayAdapter.quickToPayload(draft), [draft]);
  const { connect } = useTestAPI();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [submitMessage, setSubmitMessage] = useState<string | null>(null);
  const [submitError, setSubmitError] = useState<string | null>(null);

  const handleChange = (key: keyof QuickDraft, value: string | number | boolean) => {
    setDraft((prev) => ({ ...prev, [key]: value }));
  };

  const handleSubmit = async () => {
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
        {/* Header Section */}
        <header className="mb-10 flex flex-col md:flex-row md:items-end justify-between gap-6">
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
                <p className="text-sm text-slate-400">
                  配置邊緣設備連線參數，自動產生符合 Gateway 規範的 Payload 結構。
                </p>
              </div>
            </div>
          </div>
          
          <div className="flex items-center gap-3 text-sm text-slate-400 bg-slate-900/50 px-4 py-2 rounded-lg border border-slate-800 shadow-sm">
            <span className="flex items-center gap-1.5"><CheckCircle2 className="w-4 h-4 text-emerald-500" /> 系統就緒</span>
            <span className="w-px h-4 bg-slate-700"></span>
            <span className="font-mono text-xs">v1.4.2</span>
          </div>
        </header>
        
        <div className="grid grid-cols-1 xl:grid-cols-12 gap-8 flex-1 items-start">
          {/* Form Area */}
          <div className="xl:col-span-7 space-y-6">
            
            {/* Section 1: Connection Intent */}
            <section className="relative overflow-hidden rounded-xl border border-slate-800 bg-[#111827]/80 p-7 shadow-lg backdrop-blur-md">
              <div className="absolute top-0 left-0 w-1 h-full bg-blue-500/50"></div>
              
              <div className="flex items-center justify-between mb-8">
                <div className="flex items-center">
                  <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-blue-500/10 text-blue-400 mr-3 border border-blue-500/20">
                    <Server className="h-4 w-4" />
                  </div>
                  <div>
                    <h2 className="text-base font-bold text-white tracking-wide">1. 設備連線 (Intent Form)</h2>
                    <p className="text-xs text-slate-400 mt-0.5">定義目標設備的物理或網路通訊位置</p>
                  </div>
                </div>
              </div>
              
              <div className="grid gap-6">
                <div className="group">
                  <label className="block text-xs font-semibold tracking-wider text-slate-400 uppercase mb-2 group-focus-within:text-blue-400 transition-colors">通訊協定</label>
                  <div className="relative">
                    <select
                      data-testid="protocol-select"
                      value={draft.protocol}
                      onChange={(e) => handleChange('protocol', e.target.value)}
                      className="w-full appearance-none rounded-lg border border-slate-700 bg-slate-900 px-4 py-3 text-sm text-white focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none transition-all shadow-inner"
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
                  <div className="group">
                    <label className="block text-xs font-semibold tracking-wider text-slate-400 uppercase mb-2 group-focus-within:text-blue-400 transition-colors">主機位置 (Host)</label>
                    <input
                      type="text"
                      data-testid="host-input"
                      value={draft.host || ''}
                      onChange={(e) => handleChange('host', e.target.value)}
                      className="w-full rounded-lg border border-slate-700 bg-slate-900 px-4 py-3 text-sm text-white placeholder-slate-600 focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none transition-all shadow-inner font-mono"
                      placeholder="192.168.1.100"
                    />
                  </div>
                  <div className="group">
                    <label className="block text-xs font-semibold tracking-wider text-slate-400 uppercase mb-2 group-focus-within:text-blue-400 transition-colors">通訊埠 (Port)</label>
                    <input
                      type="number"
                      data-testid="port-input"
                      value={draft.port || ''}
                      onChange={(e) => handleChange('port', parseInt(e.target.value, 10))}
                      className="w-full rounded-lg border border-slate-700 bg-slate-900 px-4 py-3 text-sm text-white focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none transition-all shadow-inner font-mono"
                    />
                  </div>
                </div>

                {draft.protocol.includes('modbus') && (
                  <div className="animate-in fade-in slide-in-from-top-2 duration-300 group">
                    <label className="block text-xs font-semibold tracking-wider text-slate-400 uppercase mb-2 group-focus-within:text-blue-400 transition-colors">站號 (Unit ID)</label>
                    <input
                      type="number"
                      data-testid="unit-id-input"
                      value={draft.unitID || ''}
                      onChange={(e) => handleChange('unitID', parseInt(e.target.value, 10))}
                      className="w-full rounded-lg border border-slate-700 bg-slate-900 px-4 py-3 text-sm text-white focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none transition-all shadow-inner font-mono max-w-[200px]"
                    />
                  </div>
                )}
                {draft.protocol.includes('fatek') && (
                  <div className="animate-in fade-in slide-in-from-top-2 duration-300 group">
                    <label className="block text-xs font-semibold tracking-wider text-slate-400 uppercase mb-2 group-focus-within:text-blue-400 transition-colors">站號 (Station)</label>
                    <input
                      type="number"
                      data-testid="station-input"
                      value={draft.station || ''}
                      onChange={(e) => handleChange('station', parseInt(e.target.value, 10))}
                      className="w-full rounded-lg border border-slate-700 bg-slate-900 px-4 py-3 text-sm text-white focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none transition-all shadow-inner font-mono max-w-[200px]"
                    />
                  </div>
                )}
              </div>
            </section>

            {/* Section 2: Route Builder */}
            <section className="relative overflow-hidden rounded-xl border border-slate-800 bg-[#111827]/80 p-7 shadow-lg backdrop-blur-md">
              <div className="absolute top-0 left-0 w-1 h-full bg-emerald-500/50"></div>
              
              <div className="flex items-center justify-between mb-8">
                <div className="flex items-center">
                  <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-emerald-500/10 text-emerald-400 mr-3 border border-emerald-500/20">
                    <Network className="h-4 w-4" />
                  </div>
                  <div>
                    <h2 className="text-base font-bold text-white tracking-wide">2. 路由設定 (Route Builder)</h2>
                    <p className="text-xs text-slate-400 mt-0.5">定義資料在 Gateway 內的映射端點</p>
                  </div>
                </div>
              </div>
              
              <div className="group">
                <label className="block text-xs font-semibold tracking-wider text-slate-400 uppercase mb-2 group-focus-within:text-emerald-400 transition-colors">API 路徑</label>
                <div className="relative flex items-center">
                  <span className="absolute left-4 px-2 py-1 bg-slate-800 text-slate-300 font-mono text-[10px] font-bold tracking-wider rounded border border-slate-700">POST</span>
                  <input
                    type="text"
                    data-testid="route-input"
                    value={draft.route || ''}
                    onChange={(e) => handleChange('route', e.target.value)}
                    className="w-full rounded-lg border border-slate-700 bg-slate-900 py-3 pl-[4.5rem] pr-4 text-sm text-white placeholder-slate-600 focus:border-emerald-500 focus:ring-1 focus:ring-emerald-500 focus:outline-none transition-all shadow-inner font-mono"
                    placeholder="/api/v1/data"
                  />
                </div>
              </div>
            </section>

            {/* Section 3: Security */}
            <section className="relative overflow-hidden rounded-xl border border-slate-800 bg-[#111827]/80 p-7 shadow-lg backdrop-blur-md">
              <div className="absolute top-0 left-0 w-1 h-full bg-purple-500/50"></div>
              
              <div className="flex items-center justify-between mb-6">
                <div className="flex items-center">
                  <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-purple-500/10 text-purple-400 mr-3 border border-purple-500/20">
                    <Shield className="h-4 w-4" />
                  </div>
                  <div>
                    <h2 className="text-base font-bold text-white tracking-wide">3. 安全設定 (Auth)</h2>
                    <p className="text-xs text-slate-400 mt-0.5">端點存取控制與防護</p>
                  </div>
                </div>
              </div>
              
              <label className="group relative flex items-center gap-4 cursor-pointer rounded-lg border border-slate-700/50 bg-slate-900/50 p-4 transition-all hover:bg-slate-800/80 hover:border-slate-600">
                <div className="relative flex items-center justify-center">
                  <input
                    type="checkbox"
                    data-testid="auth-checkbox"
                    checked={draft.auth || false}
                    onChange={(e) => handleChange('auth', e.target.checked)}
                    className="peer sr-only"
                  />
                  <div className="h-5 w-5 rounded border border-slate-600 bg-slate-800 transition-all peer-checked:border-purple-500 peer-checked:bg-purple-500 shadow-inner"></div>
                  <svg className="absolute h-3.5 w-3.5 text-white opacity-0 transition-opacity peer-checked:opacity-100 pointer-events-none" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={3}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
                  </svg>
                </div>
                <div className="flex flex-col">
                  <span className="text-sm font-semibold text-slate-200 group-hover:text-white transition-colors">啟用基本身分驗證 (Basic Auth)</span>
                  <span className="text-xs text-slate-500 mt-1">要求 API 呼叫端提供標準帳號密碼，防止未授權存取。</span>
                </div>
              </label>
            </section>
          </div>

          {/* Payload Preview Area */}
          <div className="xl:col-span-5 xl:sticky xl:top-8 h-[calc(100vh-8rem)] min-h-[600px] flex flex-col">
            <section className="flex flex-col h-full rounded-xl border border-slate-700 bg-[#090D14] shadow-2xl overflow-hidden ring-1 ring-white/5 relative">
              {/* Terminal Header */}
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
                <div className="w-10"></div> {/* Spacer for balance */}
              </div>
              
              {/* Code Content */}
              <div className="flex-1 overflow-auto bg-[#090D14] p-0 relative group">
                {/* Line numbers background */}
                <div className="absolute left-0 top-0 bottom-0 w-12 bg-[#0D131F] border-r border-slate-800/80 pointer-events-none z-0"></div>
                
                <div className="absolute top-6 right-6 text-slate-700/30 pointer-events-none transition-opacity duration-500 group-hover:opacity-10">
                  <Code className="h-32 w-32" />
                </div>
                
                {/* Invisible element for tests to easily grab the raw payload */}
                <div data-testid="payload-preview" className="hidden">
                  {JSON.stringify(payload, null, 2)}
                </div>
                
                <div className="p-5 relative z-10 font-mono text-[13px] leading-relaxed">
                  {JSON.stringify(payload, null, 2).split('\n').map((line, i) => {
                    if (line.match(/"[^"]+":/)) line = line.replace(/"([^"]+)":/, (_m, p1) => `<span class="text-[#82AAFF]">"${p1}"</span>:`); // keys (blue)
                    if (line.match(/: "[^"]+"/)) line = line.replace(/: "([^"]+)"/, (_m, p1) => `: <span class="text-[#C3E88D]">"${p1}"</span>`); // string values (green)
                    if (line.match(/: \d+/)) line = line.replace(/: (\d+)/, (_m, p1) => `: <span class="text-[#F78C6C]">${p1}</span>`); // numbers (orange)
                    if (line.match(/: (true|false)/)) line = line.replace(/: (true|false)/, (_m, p1) => `: <span class="text-[#C792EA]">${p1}</span>`); // booleans (purple)
                    
                    return (
                      <div key={i} className="flex hover:bg-white/5 transition-colors -mx-5 px-5 group/line">
                        <span className="select-none text-slate-600/70 mr-6 text-right inline-block w-6 shrink-0 group-hover/line:text-slate-400 transition-colors">{i + 1}</span>
                        <span className="whitespace-pre-wrap break-all" dangerouslySetInnerHTML={{ __html: line }}></span>
                      </div>
                    );
                  })}
                </div>
              </div>
              
              {/* Action Footer */}
              <div className="p-5 border-t border-slate-800 bg-[#111827] shrink-0">
                <button
                  data-testid="save-btn"
                  disabled={isSubmitting}
                  className="group relative flex w-full items-center justify-center gap-2 rounded-lg bg-blue-600 px-6 py-3.5 font-semibold text-white transition-all hover:bg-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-[#111827] active:scale-[0.98] shadow-[0_0_15px_-3px_rgba(37,99,235,0.4)] hover:shadow-[0_0_20px_-3px_rgba(59,130,246,0.6)] border border-blue-500/50 disabled:cursor-not-allowed disabled:opacity-60"
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
