import { useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { Server, Network, Shield, Code, ArrowLeft, Terminal, Save } from 'lucide-react';
import { gatewayAdapter, type QuickDraft } from '../../features/gateway/gatewayAdapter';

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

  const handleChange = (key: keyof QuickDraft, value: string | number | boolean) => {
    setDraft((prev) => ({ ...prev, [key]: value }));
  };

  return (
    <main className="mx-auto max-w-6xl px-6 py-8 text-slate-100 min-h-screen flex flex-col">
      <div className="mb-8">
        <Link 
          to="/gateway/entry" 
          className="inline-flex items-center text-sm font-medium text-slate-400 hover:text-blue-400 transition-colors mb-4"
        >
          <ArrowLeft className="mr-2 h-4 w-4" />
          返回模式選擇
        </Link>
        <h1 className="text-3xl font-extrabold tracking-tight text-white mb-2">
          快速設定 <span className="text-blue-400 font-light ml-2">Quick Setup</span>
        </h1>
        <p className="text-slate-400 max-w-2xl">
          透過引導式表單快速建立連線，我們將自動套用最佳實踐與合理的預設值。
        </p>
      </div>
      
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 flex-1 items-start">
        {/* Form Area */}
        <div className="lg:col-span-7 space-y-6">
          {/* Section 1: Connection Intent */}
          <section className="relative overflow-hidden rounded-2xl border border-slate-700/50 bg-slate-900/40 p-6 shadow-lg backdrop-blur-sm transition-all hover:border-slate-600/50">
            <div className="flex items-center mb-6 border-b border-slate-800 pb-4">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-500/10 text-blue-400 mr-4 ring-1 ring-blue-500/20">
                <Server className="h-5 w-5" />
              </div>
              <h2 className="text-lg font-semibold text-slate-200">1. 設備連線 (Intent Form)</h2>
            </div>
            
            <div className="space-y-5">
              <div>
                <label className="block text-sm font-medium text-slate-300 mb-2">通訊協定</label>
                <div className="relative">
                  <select
                    data-testid="protocol-select"
                    value={draft.protocol}
                    onChange={(e) => handleChange('protocol', e.target.value)}
                    className="w-full appearance-none rounded-xl border border-slate-700 bg-slate-800/80 px-4 py-3 text-sm text-white focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 focus:outline-none transition-shadow"
                  >
                    <option value="modbus-tcp">Modbus TCP</option>
                    <option value="modbus-rtu">Modbus RTU</option>
                    <option value="fatek-tcp">Fatek TCP</option>
                    <option value="mc-tcp">Mitsubishi MC</option>
                  </select>
                  <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-4 text-slate-400">
                    <svg className="h-4 w-4 fill-current" viewBox="0 0 20 20"><path d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" /></svg>
                  </div>
                </div>
              </div>
              
              <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
                <div>
                  <label className="block text-sm font-medium text-slate-300 mb-2">主機位置 (Host)</label>
                  <input
                    type="text"
                    data-testid="host-input"
                    value={draft.host || ''}
                    onChange={(e) => handleChange('host', e.target.value)}
                    className="w-full rounded-xl border border-slate-700 bg-slate-800/80 px-4 py-3 text-sm text-white placeholder-slate-500 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 focus:outline-none transition-shadow"
                    placeholder="192.168.1.100"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-300 mb-2">通訊埠 (Port)</label>
                  <input
                    type="number"
                    data-testid="port-input"
                    value={draft.port || ''}
                    onChange={(e) => handleChange('port', parseInt(e.target.value, 10))}
                    className="w-full rounded-xl border border-slate-700 bg-slate-800/80 px-4 py-3 text-sm text-white focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 focus:outline-none transition-shadow"
                  />
                </div>
              </div>

              {draft.protocol.includes('modbus') && (
                <div className="animate-in fade-in slide-in-from-top-2 duration-300">
                  <label className="block text-sm font-medium text-slate-300 mb-2">站號 (Unit ID)</label>
                  <input
                    type="number"
                    data-testid="unit-id-input"
                    value={draft.unitID || ''}
                    onChange={(e) => handleChange('unitID', parseInt(e.target.value, 10))}
                    className="w-full rounded-xl border border-slate-700 bg-slate-800/80 px-4 py-3 text-sm text-white focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 focus:outline-none transition-shadow"
                  />
                </div>
              )}
              {draft.protocol.includes('fatek') && (
                <div className="animate-in fade-in slide-in-from-top-2 duration-300">
                  <label className="block text-sm font-medium text-slate-300 mb-2">站號 (Station)</label>
                  <input
                    type="number"
                    data-testid="station-input"
                    value={draft.station || ''}
                    onChange={(e) => handleChange('station', parseInt(e.target.value, 10))}
                    className="w-full rounded-xl border border-slate-700 bg-slate-800/80 px-4 py-3 text-sm text-white focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 focus:outline-none transition-shadow"
                  />
                </div>
              )}
            </div>
          </section>

          {/* Section 2: Route Builder */}
          <section className="relative overflow-hidden rounded-2xl border border-slate-700/50 bg-slate-900/40 p-6 shadow-lg backdrop-blur-sm transition-all hover:border-slate-600/50">
            <div className="flex items-center mb-6 border-b border-slate-800 pb-4">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-emerald-500/10 text-emerald-400 mr-4 ring-1 ring-emerald-500/20">
                <Network className="h-5 w-5" />
              </div>
              <h2 className="text-lg font-semibold text-slate-200">2. 路由設定 (Route Builder)</h2>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-slate-300 mb-2">API 路徑</label>
              <div className="relative flex items-center">
                <span className="absolute left-4 text-slate-500 font-mono text-sm">POST</span>
                <input
                  type="text"
                  data-testid="route-input"
                  value={draft.route || ''}
                  onChange={(e) => handleChange('route', e.target.value)}
                  className="w-full rounded-xl border border-slate-700 bg-slate-800/80 py-3 pl-16 pr-4 text-sm text-white placeholder-slate-500 focus:border-emerald-500 focus:ring-2 focus:ring-emerald-500/20 focus:outline-none transition-shadow font-mono"
                  placeholder="/api/v1/data"
                />
              </div>
              <p className="mt-2 text-xs text-slate-500">此路徑將作為 Gateway 轉發資料的進入點。</p>
            </div>
          </section>

          {/* Section 3: Security */}
          <section className="relative overflow-hidden rounded-2xl border border-slate-700/50 bg-slate-900/40 p-6 shadow-lg backdrop-blur-sm transition-all hover:border-slate-600/50">
            <div className="flex items-center mb-6 border-b border-slate-800 pb-4">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-purple-500/10 text-purple-400 mr-4 ring-1 ring-purple-500/20">
                <Shield className="h-5 w-5" />
              </div>
              <h2 className="text-lg font-semibold text-slate-200">3. 安全設定 (Auth)</h2>
            </div>
            
            <label className="group relative flex items-start gap-4 cursor-pointer rounded-xl border border-transparent p-4 transition-colors hover:bg-slate-800/50">
              <div className="flex h-6 items-center">
                <input
                  type="checkbox"
                  data-testid="auth-checkbox"
                  checked={draft.auth || false}
                  onChange={(e) => handleChange('auth', e.target.checked)}
                  className="peer sr-only"
                />
                <div className="h-5 w-5 rounded border border-slate-600 bg-slate-900 ring-offset-slate-900 transition-all peer-checked:border-blue-500 peer-checked:bg-blue-500 peer-focus-visible:ring-2 peer-focus-visible:ring-blue-500 peer-focus-visible:ring-offset-2 flex items-center justify-center">
                  <svg className="h-4 w-4 text-white opacity-0 transition-opacity peer-checked:opacity-100" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={3}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
                  </svg>
                </div>
              </div>
              <div className="flex flex-col">
                <span className="text-sm font-medium text-slate-200 group-hover:text-white transition-colors">啟用基本身分驗證 (Basic Auth)</span>
                <span className="text-xs text-slate-500 mt-1">要求 API 呼叫端提供帳號密碼，提升介面安全性。</span>
              </div>
            </label>
          </section>
        </div>

        {/* Payload Preview Area */}
        <div className="lg:col-span-5 lg:sticky lg:top-8 lg:h-[calc(100vh-8rem)]">
          <section className="flex flex-col h-full rounded-2xl border border-slate-700/80 bg-[#0d1117] shadow-2xl overflow-hidden ring-1 ring-white/5">
            {/* Terminal Header */}
            <div className="flex items-center justify-between px-4 py-3 bg-slate-800/40 border-b border-slate-700/80">
              <div className="flex items-center gap-2">
                <div className="h-3 w-3 rounded-full bg-red-500/80"></div>
                <div className="h-3 w-3 rounded-full bg-yellow-500/80"></div>
                <div className="h-3 w-3 rounded-full bg-green-500/80"></div>
              </div>
              <div className="flex items-center gap-2 text-xs font-mono text-slate-400">
                <Terminal className="h-3 w-3" />
                <span>Payload 預覽 (供驗證)</span>
              </div>
              <div className="w-10"></div> {/* Spacer for balance */}
            </div>
            
            {/* Code Content */}
            <div className="flex-1 overflow-auto p-5 relative min-h-[400px]">
              <div className="absolute top-5 right-5 text-slate-600/50 pointer-events-none">
                <Code className="h-24 w-24 opacity-10" />
              </div>
              <div data-testid="payload-preview" className="hidden">
                {JSON.stringify(payload, null, 2)}
              </div>
              <pre 
                className="font-mono text-[13px] leading-relaxed relative z-10"
              >
                {JSON.stringify(payload, null, 2).split('\n').map((line, i) => {
                  let colorClass = "text-slate-300";
                  if (line.includes('": "')) colorClass = "text-emerald-300"; // strings
                  else if (line.match(/": \d+/)) colorClass = "text-orange-300"; // numbers
                  else if (line.includes('": true') || line.includes('": false')) colorClass = "text-purple-300"; // booleans
                  else if (line.includes('":')) colorClass = "text-blue-300"; // keys
                  
                  return (
                    <div key={i} className="hover:bg-white/5 px-2 -mx-2 rounded transition-colors flex">
                      <span className="select-none text-slate-600 mr-4 text-right inline-block w-4 shrink-0">{i + 1}</span>
                      <span className={`${colorClass} whitespace-pre-wrap break-all`}>{line}</span>
                    </div>
                  );
                })}
              </pre>
            </div>
            
            {/* Action Footer */}
            <div className="p-4 border-t border-slate-800 bg-slate-900/50 backdrop-blur shrink-0">
              <p className="text-xs text-slate-500 mb-3 text-center">
                此預覽即時反映表單設定，將由 Gateway Adapter 轉譯
              </p>
              <button 
                data-testid="save-btn"
                className="group relative flex w-full items-center justify-center gap-2 rounded-xl bg-gradient-to-r from-blue-600 to-indigo-600 px-6 py-3 font-medium text-white shadow-[0_0_20px_-5px_rgba(59,130,246,0.5)] transition-all hover:scale-[1.02] hover:shadow-[0_0_25px_-5px_rgba(59,130,246,0.6)] focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-slate-900 active:scale-100"
                onClick={() => alert('目前為預覽模式，尚未發送至後端')}
              >
                <Save className="h-5 w-5 transition-transform group-hover:-translate-y-0.5" />
                <span>儲存並套用設定</span>
              </button>
            </div>
          </section>
        </div>
      </div>
    </main>
  );
}
