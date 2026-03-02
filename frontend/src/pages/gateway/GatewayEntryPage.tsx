import { Link } from 'react-router-dom';
import { Zap, Wrench, ArrowRight } from 'lucide-react';

export default function GatewayEntryPage() {
  return (
    <main className="mx-auto max-w-4xl px-6 py-16 text-slate-100 min-h-[calc(100vh-4rem)] flex flex-col justify-center">
      <div className="text-center mb-12">
        <h1 className="text-4xl font-extrabold tracking-tight mb-4">
          <span className="bg-gradient-to-r from-blue-400 to-indigo-400 bg-clip-text text-transparent">
            選擇工作模式
          </span>
        </h1>
        <p className="text-lg text-slate-400 max-w-xl mx-auto">
          依據您的使用情境，選擇快速引導的入門流程，或是完整控制的專家模式。
        </p>
      </div>

      <section className="grid gap-6 md:grid-cols-2">
        {/* Quick Setup Card */}
        <Link
          to="/gateway/quick-setup"
          className="group relative flex flex-col justify-between overflow-hidden rounded-2xl border border-slate-700/50 bg-slate-800/40 p-8 transition-all hover:-translate-y-1 hover:border-blue-500/50 hover:bg-slate-800/60 hover:shadow-[0_0_30px_-5px_rgba(59,130,246,0.3)]"
        >
          <div className="absolute inset-0 bg-gradient-to-br from-blue-500/10 to-transparent opacity-0 transition-opacity group-hover:opacity-100" />
          
          <div>
            <div className="mb-6 inline-flex h-12 w-12 items-center justify-center rounded-xl bg-blue-500/20 text-blue-400 ring-1 ring-blue-500/30">
              <Zap className="h-6 w-6" />
            </div>
            <h2 className="mb-2 text-2xl font-bold text-slate-100 flex items-center gap-2">
              快速設定
              <span className="inline-flex items-center rounded-full bg-blue-500/10 px-2 py-0.5 text-xs font-medium text-blue-400 ring-1 ring-inset ring-blue-500/20">
                推薦
              </span>
            </h2>
            <p className="text-slate-400 leading-relaxed">
              適合常見的設備代理場景。透過引導式表單與智能預設值，在幾分鐘內完成您的第一個設備連線。
            </p>
          </div>
          
          <div className="mt-8 flex items-center text-sm font-medium text-blue-400">
            開始快速設定
            <ArrowRight className="ml-1 h-4 w-4 transition-transform group-hover:translate-x-1" />
          </div>
        </Link>

        {/* Expert Workbench Card */}
        <Link
          to="/gateway/expert-workbench"
          className="group relative flex flex-col justify-between overflow-hidden rounded-2xl border border-slate-700/50 bg-slate-800/40 p-8 transition-all hover:-translate-y-1 hover:border-emerald-500/50 hover:bg-slate-800/60 hover:shadow-[0_0_30px_-5px_rgba(16,185,129,0.2)]"
        >
          <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/5 to-transparent opacity-0 transition-opacity group-hover:opacity-100" />
          
          <div>
            <div className="mb-6 inline-flex h-12 w-12 items-center justify-center rounded-xl bg-slate-700/50 text-slate-300 ring-1 ring-slate-600">
              <Wrench className="h-6 w-6" />
            </div>
            <h2 className="mb-2 text-2xl font-bold text-slate-100">專家工作台</h2>
            <p className="text-slate-400 leading-relaxed">
              專為進階需求打造。支援高密度路由配置、複雜通訊參數微調，提供毫無限制的底層控制能力。
            </p>
          </div>
          
          <div className="mt-8 flex items-center text-sm font-medium text-slate-300 group-hover:text-emerald-400 transition-colors">
            進入專家模式
            <ArrowRight className="ml-1 h-4 w-4 transition-transform group-hover:translate-x-1" />
          </div>
        </Link>
      </section>
    </main>
  );
}
