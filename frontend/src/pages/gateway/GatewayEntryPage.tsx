import { type MouseEvent } from 'react';
import { Link } from 'react-router-dom';
import { Zap, Wrench, ArrowRight, Activity, Terminal, ShieldCheck, Database } from 'lucide-react';
import {
  confirmGatewayExpertDowngradeToQuick,
  useGatewayDraftStore,
} from '../../features/gateway/gatewayDraftStore';

export default function GatewayEntryPage() {
  const expertQuickCompatibility = useGatewayDraftStore((state) => state.expertQuickCompatibility);
  const expertUnsupportedKeys = useGatewayDraftStore((state) => state.expertUnsupportedKeys);

  function handleQuickEntryClick(event: MouseEvent<HTMLAnchorElement>) {
    if (expertQuickCompatibility !== 'unsupported') return;

    const detail = expertUnsupportedKeys.length > 0 ? `\n不支援欄位：${expertUnsupportedKeys.join(', ')}` : '';
    const confirmed = window.confirm(
      `【降級警告】目前 Expert 設定包含 Quick 不支援欄位。\n切換後將移除這些設定且無法自動復原。\n是否仍要切換到 Quick？${detail}`,
    );
    if (!confirmed) {
      event.preventDefault();
      return;
    }

    confirmGatewayExpertDowngradeToQuick();
  }

  return (
    <main className="relative mx-auto max-w-5xl px-6 py-16 text-slate-100 min-h-[calc(100vh-4rem)] flex flex-col justify-center">
      <div className="absolute inset-0 -z-10 h-full w-full bg-[radial-gradient(#1e293b_1px,transparent_1px)] [background-size:16px_16px] opacity-20"></div>

      <div className="text-center mb-16 relative z-10">
        <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-500/10 border border-blue-500/20 text-blue-400 text-sm font-medium mb-6">
          <Activity className="w-4 h-4" />
          <span>Go-Gateway Initialization</span>
        </div>
        <h1 className="text-4xl md:text-5xl font-extrabold tracking-tight mb-6 text-white">
          選擇工作模式
        </h1>
        <p className="text-lg text-slate-400 max-w-2xl mx-auto">
          依據您的部署情境與技術需求，選擇適當的初始化流程。系統將自動最佳化相應的資源分配與通訊參數。
        </p>
      </div>

      <section className="grid gap-8 md:grid-cols-2 relative z-10">
        <Link
          to="/gateway/quick-setup"
          data-testid="gateway-entry-quick-link"
          onClick={handleQuickEntryClick}
          className="group relative flex flex-col overflow-hidden rounded-2xl border border-slate-700/60 bg-slate-900/50 p-8 transition-all duration-300 hover:-translate-y-1 hover:border-blue-500/50 hover:bg-slate-800/80 hover:shadow-[0_8px_32px_-10px_rgba(59,130,246,0.3)] backdrop-blur-sm"
        >
          <div className="absolute inset-0 bg-gradient-to-br from-blue-500/5 via-transparent to-transparent opacity-0 transition-opacity duration-500 group-hover:opacity-100" />
          
          <div className="relative z-10 flex-1">
            <div className="flex items-start justify-between mb-6">
              <div className="inline-flex h-14 w-14 items-center justify-center rounded-xl bg-gradient-to-br from-blue-500/20 to-blue-600/10 text-blue-400 ring-1 ring-blue-500/30 shadow-inner">
                <Zap className="h-7 w-7" />
              </div>
              <span className="inline-flex items-center rounded-full bg-blue-500/10 px-2.5 py-1 text-xs font-semibold text-blue-400 ring-1 ring-inset ring-blue-500/20">
                推薦 (Recommended)
              </span>
            </div>
            
            <h2 className="mb-3 text-2xl font-bold text-white tracking-wide">快速設定</h2>
            <p className="text-slate-400 leading-relaxed mb-6 text-sm">
              標準設備代理場景。透過智慧型精靈引導，自動套用最佳實踐配置，在數分鐘內建立穩定的工業通訊連線。
            </p>

            <div className="space-y-3 mb-8">
              <div className="flex items-center text-xs text-slate-300 gap-2">
                <ShieldCheck className="w-4 h-4 text-emerald-400" /> 內建防呆與參數驗證
              </div>
              <div className="flex items-center text-xs text-slate-300 gap-2">
                <Database className="w-4 h-4 text-blue-400" /> 自動生成對應路由
              </div>
              {expertQuickCompatibility === 'unsupported' ? (
                <div data-testid="gateway-entry-downgrade-warning" className="flex items-center text-xs text-amber-300 gap-2">
                  <ShieldCheck className="w-4 h-4 text-amber-400" /> 【降級警告】切回 Quick 將移除不支援欄位
                </div>
              ) : null}
            </div>
          </div>
          
          <div className="relative z-10 mt-auto pt-6 border-t border-slate-700/50 flex items-center justify-between text-sm font-semibold text-slate-300 group-hover:text-blue-400 transition-colors">
            <span>開始快速部署</span>
            <ArrowRight className="h-5 w-5 transition-transform duration-300 group-hover:translate-x-1.5" />
          </div>
        </Link>

        <Link
          to="/gateway/expert-workbench"
          data-testid="gateway-entry-expert-link"
          className="group relative flex flex-col overflow-hidden rounded-2xl border border-slate-700/60 bg-slate-900/50 p-8 transition-all duration-300 hover:-translate-y-1 hover:border-emerald-500/50 hover:bg-slate-800/80 hover:shadow-[0_8px_32px_-10px_rgba(16,185,129,0.2)] backdrop-blur-sm"
        >
          <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/5 via-transparent to-transparent opacity-0 transition-opacity duration-500 group-hover:opacity-100" />
          
          <div className="relative z-10 flex-1">
            <div className="flex items-start justify-between mb-6">
              <div className="inline-flex h-14 w-14 items-center justify-center rounded-xl bg-gradient-to-br from-emerald-500/10 to-teal-500/5 text-emerald-400 ring-1 ring-emerald-500/30 shadow-inner">
                <Wrench className="h-7 w-7" />
              </div>
              <span className="inline-flex items-center rounded-full bg-slate-800 px-2.5 py-1 text-xs font-semibold text-slate-300 ring-1 ring-inset ring-slate-600">
                進階 (Advanced)
              </span>
            </div>
            
            <h2 className="mb-3 text-2xl font-bold text-white tracking-wide">專家工作台</h2>
            <p className="text-slate-400 leading-relaxed mb-6 text-sm">
              專為複雜拓撲架構打造。支援高密度多點路由配置、通訊逾時重試策略微調，提供 100% 的底層引擎控制權限。
            </p>

            <div className="space-y-3 mb-8">
              <div className="flex items-center text-xs text-slate-300 gap-2">
                <Terminal className="w-4 h-4 text-emerald-400" /> 原始 Payload 編輯器
              </div>
              <div className="flex items-center text-xs text-slate-300 gap-2">
                <Database className="w-4 h-4 text-purple-400" /> 跨協定橋接與多重路由
              </div>
            </div>
          </div>
          
          <div className="relative z-10 mt-auto pt-6 border-t border-slate-700/50 flex items-center justify-between text-sm font-semibold text-slate-300 group-hover:text-emerald-400 transition-colors">
            <span>進入專家模式</span>
            <ArrowRight className="h-5 w-5 transition-transform duration-300 group-hover:translate-x-1.5" />
          </div>
        </Link>
      </section>
    </main>
  );
}
