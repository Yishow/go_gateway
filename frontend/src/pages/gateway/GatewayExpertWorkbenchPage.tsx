import { Link } from 'react-router-dom';

export default function GatewayExpertWorkbenchPage() {
  return (
    <main className="mx-auto max-w-4xl px-6 py-10 text-slate-100">
      <h1 className="text-2xl font-bold">Expert Workbench（W1 骨架）</h1>
      <p className="mt-2 text-sm text-slate-300">
        這是 W1 的最小可展示版：先提供路由與容器，後續再接 AdvancedRouteTable / PluginChain / RawManifestEditor。
      </p>

      <div className="mt-6 rounded-xl border border-dashed border-slate-600 p-5 text-sm text-slate-300">
        TODO（W3）: AdvancedRouteTable / PluginChainVisualizer / RawManifestEditor
      </div>

      <Link to="/gateway/entry" className="mt-6 inline-block text-sm text-blue-300 hover:text-blue-200">
        ← 回雙入口選擇
      </Link>
    </main>
  );
}
