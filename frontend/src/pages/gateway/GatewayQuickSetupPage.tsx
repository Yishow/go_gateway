import { Link } from 'react-router-dom';

export default function GatewayQuickSetupPage() {
  return (
    <main className="mx-auto max-w-4xl px-6 py-10 text-slate-100">
      <h1 className="text-2xl font-bold">Quick Setup（W1 骨架）</h1>
      <p className="mt-2 text-sm text-slate-300">
        這是 W1 的最小可展示版：先提供路由、頁面骨架與模式切換入口，暫不提交任何 API。
      </p>

      <div className="mt-6 rounded-xl border border-dashed border-slate-600 p-5 text-sm text-slate-300">
        TODO（W2）: IntentForm / SimpleRouteBuilder / BasicAuthToggle
      </div>

      <Link to="/gateway/entry" className="mt-6 inline-block text-sm text-blue-300 hover:text-blue-200">
        ← 回雙入口選擇
      </Link>
    </main>
  );
}
