import { useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
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
    <main className="mx-auto max-w-5xl px-6 py-10 text-slate-100">
      <h1 className="text-2xl font-bold mb-6">快速設定 (Quick Setup)</h1>
      <p className="mb-8 text-sm text-slate-300">
        透過簡化表單快速建立連線設定。系統會自動套用合理的預設值。
      </p>
      
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div className="space-y-6">
          {/* IntentForm */}
          <section className="rounded-xl border border-slate-700 bg-slate-800/50 p-5 shadow-sm">
            <h2 className="text-lg font-semibold mb-4 text-slate-200">1. 設備連線 (Intent Form)</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-400 mb-1">通訊協定</label>
                <select
                  data-testid="protocol-select"
                  value={draft.protocol}
                  onChange={(e) => handleChange('protocol', e.target.value)}
                  className="w-full rounded-md border border-slate-600 bg-slate-900 px-3 py-2 text-sm text-white focus:border-blue-500 focus:outline-none"
                >
                  <option value="modbus-tcp">Modbus TCP</option>
                  <option value="modbus-rtu">Modbus RTU</option>
                  <option value="fatek-tcp">Fatek TCP</option>
                  <option value="mc-tcp">Mitsubishi MC</option>
                </select>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-slate-400 mb-1">主機位置 (Host)</label>
                  <input
                    type="text"
                    data-testid="host-input"
                    value={draft.host || ''}
                    onChange={(e) => handleChange('host', e.target.value)}
                    className="w-full rounded-md border border-slate-600 bg-slate-900 px-3 py-2 text-sm text-white focus:border-blue-500 focus:outline-none"
                    placeholder="192.168.1.100"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-slate-400 mb-1">通訊埠 (Port)</label>
                  <input
                    type="number"
                    data-testid="port-input"
                    value={draft.port || ''}
                    onChange={(e) => handleChange('port', parseInt(e.target.value, 10))}
                    className="w-full rounded-md border border-slate-600 bg-slate-900 px-3 py-2 text-sm text-white focus:border-blue-500 focus:outline-none"
                  />
                </div>
              </div>
              {draft.protocol.includes('modbus') && (
                <div>
                  <label className="block text-sm font-medium text-slate-400 mb-1">站號 (Unit ID)</label>
                  <input
                    type="number"
                    data-testid="unit-id-input"
                    value={draft.unitID || ''}
                    onChange={(e) => handleChange('unitID', parseInt(e.target.value, 10))}
                    className="w-full rounded-md border border-slate-600 bg-slate-900 px-3 py-2 text-sm text-white focus:border-blue-500 focus:outline-none"
                  />
                </div>
              )}
              {draft.protocol.includes('fatek') && (
                <div>
                  <label className="block text-sm font-medium text-slate-400 mb-1">站號 (Station)</label>
                  <input
                    type="number"
                    data-testid="station-input"
                    value={draft.station || ''}
                    onChange={(e) => handleChange('station', parseInt(e.target.value, 10))}
                    className="w-full rounded-md border border-slate-600 bg-slate-900 px-3 py-2 text-sm text-white focus:border-blue-500 focus:outline-none"
                  />
                </div>
              )}
            </div>
          </section>

          {/* SimpleRouteBuilder */}
          <section className="rounded-xl border border-slate-700 bg-slate-800/50 p-5 shadow-sm">
            <h2 className="text-lg font-semibold mb-4 text-slate-200">2. 路由設定 (Route Builder)</h2>
            <div>
              <label className="block text-sm font-medium text-slate-400 mb-1">API 路徑</label>
              <input
                type="text"
                data-testid="route-input"
                value={draft.route || ''}
                onChange={(e) => handleChange('route', e.target.value)}
                className="w-full rounded-md border border-slate-600 bg-slate-900 px-3 py-2 text-sm text-white focus:border-blue-500 focus:outline-none"
                placeholder="/api/v1/data"
              />
            </div>
          </section>

          {/* BasicAuthToggle */}
          <section className="rounded-xl border border-slate-700 bg-slate-800/50 p-5 shadow-sm">
            <h2 className="text-lg font-semibold mb-4 text-slate-200">3. 安全設定 (Auth)</h2>
            <label className="flex items-center space-x-3 cursor-pointer">
              <input
                type="checkbox"
                data-testid="auth-checkbox"
                checked={draft.auth || false}
                onChange={(e) => handleChange('auth', e.target.checked)}
                className="form-checkbox h-5 w-5 text-blue-500 rounded border-slate-600 bg-slate-900 focus:ring-blue-500 focus:ring-offset-slate-900"
              />
              <span className="text-sm font-medium text-slate-300">啟用基本身分驗證 (Basic Auth)</span>
            </label>
          </section>
        </div>

        {/* Payload Preview */}
        <div className="space-y-6">
          <section className="rounded-xl border border-slate-700 bg-slate-900 p-5 h-full flex flex-col shadow-sm">
            <h2 className="text-lg font-semibold mb-4 text-slate-200">Payload 預覽 (供驗證)</h2>
            <p className="text-xs text-slate-400 mb-4">
              此區域展示 `gatewayAdapter` 即時轉換的 API Payload，目前不直接呼叫後端。
            </p>
            <pre 
              data-testid="payload-preview"
              className="flex-1 bg-black p-4 rounded-md overflow-auto text-sm text-green-400 border border-slate-800"
            >
              {JSON.stringify(payload, null, 2)}
            </pre>
            
            <div className="mt-4 flex gap-3">
              <button 
                data-testid="save-btn"
                className="flex-1 bg-blue-600 hover:bg-blue-500 text-white py-2 px-4 rounded-md font-medium transition-colors"
                onClick={() => alert('目前為預覽模式，尚未發送至後端')}
              >
                儲存設定 (模擬)
              </button>
            </div>
          </section>
        </div>
      </div>

      <Link to="/gateway/entry" className="mt-8 inline-block text-sm text-blue-400 hover:text-blue-300 transition-colors">
        ← 回雙入口選擇
      </Link>
    </main>
  );
}