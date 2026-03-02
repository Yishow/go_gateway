import { useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import { gatewayAdapter } from '../../features/gateway/gatewayAdapter';

type RouteDraft = {
  id: string;
  method: 'GET' | 'POST' | 'PUT' | 'DELETE';
  path: string;
  target: string;
  enabled: boolean;
};

function createDefaultRoutes(): RouteDraft[] {
  return [
    {
      id: 'route-1',
      method: 'GET',
      path: '/api/v1/metrics',
      target: 'http://collector.internal:9000/metrics',
      enabled: true,
    },
    {
      id: 'route-2',
      method: 'POST',
      path: '/api/v1/write',
      target: 'http://collector.internal:9000/write',
      enabled: true,
    },
  ];
}

export default function GatewayExpertWorkbenchPage() {
  const [protocol, setProtocol] = useState('modbus-tcp');
  const [routes, setRoutes] = useState<RouteDraft[]>(createDefaultRoutes());
  const [plugins, setPlugins] = useState({
    requestId: true,
    cors: true,
    rateLimit: false,
  });

  const builderConfig = useMemo(
    () => ({
      timeoutMs: 5000,
      routes: routes.map((route) => ({
        method: route.method,
        path: route.path,
        target: route.target,
        enabled: route.enabled,
      })),
      plugins,
    }),
    [plugins, routes],
  );

  const [rawManifest, setRawManifest] = useState(() => JSON.stringify(builderConfig, null, 2));

  const payloadResult = useMemo(() => {
    try {
      const payload = gatewayAdapter.expertToPayload({
        protocol,
        configJson: rawManifest,
      });

      return {
        payload,
        error: null as string | null,
      };
    } catch (error) {
      return {
        payload: null,
        error: error instanceof Error ? error.message : '未知錯誤',
      };
    }
  }, [protocol, rawManifest]);

  function updateRoute(id: string, key: keyof Omit<RouteDraft, 'id'>, value: string | boolean) {
    setRoutes((prev) =>
      prev.map((route) => (route.id === id ? { ...route, [key]: value } : route)),
    );
  }

  function addRoute() {
    const index = routes.length + 1;
    setRoutes((prev) => [
      ...prev,
      {
        id: `route-${Date.now()}`,
        method: 'GET',
        path: `/api/v1/custom-${index}`,
        target: 'http://collector.internal:9000/custom',
        enabled: true,
      },
    ]);
  }

  function removeRoute(id: string) {
    setRoutes((prev) => prev.filter((route) => route.id !== id));
  }

  function syncBuilderToManifest() {
    setRawManifest(JSON.stringify(builderConfig, null, 2));
  }

  return (
    <main className="mx-auto max-w-7xl px-6 py-10 text-slate-100">
      <header className="rounded-2xl border border-slate-700/80 bg-gradient-to-r from-slate-950 to-slate-900 p-6 shadow-xl">
        <p className="text-xs uppercase tracking-[0.18em] text-cyan-300">Gateway C 方案</p>
        <h1 className="mt-2 text-2xl font-bold">Expert Workbench</h1>
        <p className="mt-2 text-sm text-slate-300">
          高密度路由編輯 + 插件鏈設定 + Raw Manifest 編輯器。此頁以 expertToPayload 產生 API 請求預覽。
        </p>
      </header>

      <section className="mt-6 grid gap-6 lg:grid-cols-[1.25fr_1fr]">
        <div className="space-y-6">
          <article className="rounded-xl border border-slate-700/70 bg-slate-900/70 p-5 shadow-lg backdrop-blur-sm">
            <div className="mb-4 flex items-center justify-between">
              <h2 className="text-lg font-semibold">Advanced Route Table</h2>
              <button
                type="button"
                onClick={addRoute}
                className="rounded-md border border-cyan-500/40 bg-cyan-500/10 px-3 py-1.5 text-xs font-medium text-cyan-200 transition hover:bg-cyan-500/20"
              >
                + 新增路由
              </button>
            </div>

            <div className="overflow-x-auto" data-testid="expert-route-table">
              <table className="min-w-full text-sm">
                <thead>
                  <tr className="border-b border-slate-700 text-slate-400">
                    <th className="px-2 py-2 text-left">Method</th>
                    <th className="px-2 py-2 text-left">Path</th>
                    <th className="px-2 py-2 text-left">Target</th>
                    <th className="px-2 py-2 text-left">Enabled</th>
                    <th className="px-2 py-2 text-right">操作</th>
                  </tr>
                </thead>
                <tbody>
                  {routes.map((route) => (
                    <tr key={route.id} className="border-b border-slate-800/70">
                      <td className="px-2 py-2">
                        <select
                          value={route.method}
                          onChange={(e) => updateRoute(route.id, 'method', e.target.value)}
                          className="rounded-md border border-slate-700 bg-slate-950 px-2 py-1"
                        >
                          <option value="GET">GET</option>
                          <option value="POST">POST</option>
                          <option value="PUT">PUT</option>
                          <option value="DELETE">DELETE</option>
                        </select>
                      </td>
                      <td className="px-2 py-2">
                        <input
                          value={route.path}
                          onChange={(e) => updateRoute(route.id, 'path', e.target.value)}
                          className="w-full rounded-md border border-slate-700 bg-slate-950 px-2 py-1"
                        />
                      </td>
                      <td className="px-2 py-2">
                        <input
                          value={route.target}
                          onChange={(e) => updateRoute(route.id, 'target', e.target.value)}
                          className="w-full rounded-md border border-slate-700 bg-slate-950 px-2 py-1"
                        />
                      </td>
                      <td className="px-2 py-2">
                        <input
                          type="checkbox"
                          checked={route.enabled}
                          onChange={(e) => updateRoute(route.id, 'enabled', e.target.checked)}
                          className="h-4 w-4"
                        />
                      </td>
                      <td className="px-2 py-2 text-right">
                        <button
                          type="button"
                          onClick={() => removeRoute(route.id)}
                          className="text-xs text-red-300 transition hover:text-red-200"
                        >
                          刪除
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </article>

          <article className="rounded-xl border border-slate-700/70 bg-slate-900/70 p-5 shadow-lg backdrop-blur-sm">
            <h2 className="text-lg font-semibold">Plugin Chain</h2>
            <p className="mt-1 text-xs text-slate-400">先提供常用開關，後續可延伸為拖曳式鏈路。</p>

            <div className="mt-4 grid gap-3 sm:grid-cols-3" data-testid="expert-plugin-chain">
              {[
                ['requestId', 'Request ID'],
                ['cors', 'CORS'],
                ['rateLimit', 'Rate Limit'],
              ].map(([key, label]) => {
                const pluginKey = key as keyof typeof plugins;
                return (
                  <label
                    key={key}
                    className="flex items-center justify-between rounded-lg border border-slate-700 bg-slate-950 px-3 py-2"
                  >
                    <span className="text-sm text-slate-200">{label}</span>
                    <input
                      type="checkbox"
                      checked={plugins[pluginKey]}
                      onChange={(e) => setPlugins((prev) => ({ ...prev, [pluginKey]: e.target.checked }))}
                    />
                  </label>
                );
              })}
            </div>
          </article>
        </div>

        <div className="space-y-6">
          <article className="rounded-xl border border-slate-700/70 bg-slate-900/70 p-5 shadow-lg backdrop-blur-sm">
            <div className="mb-3 flex items-center justify-between gap-3">
              <h2 className="text-lg font-semibold">Raw Manifest Editor</h2>
              <button
                type="button"
                onClick={syncBuilderToManifest}
                className="rounded-md border border-slate-600 px-3 py-1 text-xs text-slate-200 hover:border-cyan-400"
              >
                由 Builder 覆寫
              </button>
            </div>

            <label className="mb-2 block text-xs text-slate-400">Protocol</label>
            <select
              data-testid="expert-protocol-select"
              value={protocol}
              onChange={(e) => setProtocol(e.target.value)}
              className="mb-4 w-full rounded-md border border-slate-700 bg-slate-950 px-3 py-2 text-sm"
            >
              <option value="modbus-tcp">Modbus TCP</option>
              <option value="modbus-rtu">Modbus RTU</option>
              <option value="fatek-tcp">Fatek TCP</option>
              <option value="mc-tcp">Mitsubishi MC</option>
            </select>

            <textarea
              data-testid="expert-raw-manifest"
              value={rawManifest}
              onChange={(e) => setRawManifest(e.target.value)}
              className="h-72 w-full rounded-lg border border-slate-700 bg-[#0b1220] p-3 font-mono text-xs text-cyan-100"
            />

            {payloadResult.error ? (
              <p data-testid="expert-manifest-error" className="mt-3 rounded-md border border-red-500/40 bg-red-500/10 px-3 py-2 text-xs text-red-200">
                解析失敗：{payloadResult.error}
              </p>
            ) : null}
          </article>

          <article className="rounded-xl border border-slate-700/70 bg-slate-950 p-5 shadow-xl">
            <h2 className="text-lg font-semibold">Payload 預覽（Expert）</h2>
            <p className="mt-1 text-xs text-slate-400">送出前可確認是否符合既有 API 契約。</p>
            <pre
              data-testid="expert-payload-preview"
              className="mt-4 max-h-72 overflow-auto rounded-lg border border-slate-800 bg-black p-3 text-xs text-emerald-300"
            >
              {JSON.stringify(payloadResult.payload ?? { error: payloadResult.error }, null, 2)}
            </pre>
          </article>
        </div>
      </section>

      <Link to="/gateway/entry" className="mt-8 inline-block text-sm text-blue-300 hover:text-blue-200">
        ← 回雙入口選擇
      </Link>
    </main>
  );
}
