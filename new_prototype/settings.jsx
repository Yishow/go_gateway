// 設定頁
// - 資料庫連接器池
// - 時序儲存策略
// - 排程預設
// - Local Modbus Share
// - 介面/系統

function SettingsPage({ state, setState }) {
  const settings = state.settings || {};
  const update = (patch) => setState((s) => ({ ...s, settings: { ...(s.settings || {}), ...patch } }));
  const updateNested = (key, patch) =>
    setState((s) => ({
      ...s,
      settings: { ...(s.settings || {}), [key]: { ...(s.settings?.[key] || {}), ...patch } },
    }));

  // Connector 池：以陣列管理
  const connectors = settings.connectors || [];
  const addConnector = () => {
    const id = `conn-${Date.now().toString(36)}`;
    update({
      connectors: [
        ...connectors,
        {
          id,
          name: `新連線 ${connectors.length + 1}`,
          kind: 'postgres',
          host: 'localhost',
          port: 5432,
          database: 'metrics',
          username: 'gw_writer',
          enabled: true,
          status: 'unknown',
          default_write_interval_seconds: 5,
        },
      ],
    });
  };
  const updateConnector = (id, patch) =>
    update({ connectors: connectors.map((c) => (c.id === id ? { ...c, ...patch } : c)) });
  const removeConnector = (id) => update({ connectors: connectors.filter((c) => c.id !== id) });

  const [testingId, setTestingId] = React.useState(null);
  const testConnector = (id) => {
    setTestingId(id);
    updateConnector(id, { status: 'testing' });
    setTimeout(() => {
      const ok = Math.random() > 0.15;
      updateConnector(id, {
        status: ok ? 'ready' : 'unreachable',
        last_check_at: new Date().toISOString(),
        last_check_error: ok ? '' : 'connection refused',
      });
      setTestingId(null);
    }, 900);
  };

  const ts = settings.timeseries || {};
  const sched = settings.scheduler || {};
  const modbusShare = settings.modbus_share || {};
  const general = settings.general || {};

  return (
    <div className="space-y-5 sweep-in">
      {/* 標題列 */}
      <div className="rounded-2xl border border-slate-700/40 bg-gradient-to-br from-slate-900/80 to-slate-900/40 p-5">
        <div className="flex items-center gap-3">
          <div className="grid place-items-center w-11 h-11 rounded-xl bg-blue-500/15 border border-blue-500/30 text-blue-300">
            <Icon name="sliders" className="w-5 h-5" />
          </div>
          <div>
            <h2 className="text-lg font-semibold text-slate-100">系統設定</h2>
            <p className="text-xs text-slate-400 mt-0.5">資料庫連接器池、時序儲存策略、排程器、Local Modbus Share、介面偏好</p>
          </div>
        </div>
      </div>

      {/* 1. 資料庫連接器池 */}
      <SectionCard
        title="資料庫連接器池"
        subtitle="可註冊多個目標資料庫，於映射時挑選"
        icon={<Icon name="db" className="w-5 h-5" />}
        aside={
          <Button size="sm" variant="secondary" onClick={addConnector} icon={<Icon name="plus" className="w-3.5 h-3.5" />}>
            新增連接器
          </Button>
        }
        contentClassName="p-0"
      >
        {connectors.length === 0 ? (
          <div className="px-5 py-8 text-center text-sm text-slate-500">
            尚未註冊任何連接器。點上方「新增連接器」開始。
          </div>
        ) : (
          <div className="divide-y divide-slate-800/60">
            {connectors.map((c) => {
              const kindMeta = DB_KINDS.find((k) => k.id === c.kind);
              const statusTone =
                c.status === 'ready' ? 'success' :
                c.status === 'testing' ? 'info' :
                c.status === 'unreachable' || c.status === 'auth_failed' ? 'error' : 'draft';
              const statusLabel =
                c.status === 'ready' ? '已就緒' :
                c.status === 'testing' ? '測試中…' :
                c.status === 'unreachable' ? '無法連線' :
                c.status === 'auth_failed' ? '認證失敗' : '未測試';
              return (
                <div key={c.id} className="px-5 py-4 hover:bg-slate-900/30">
                  <div className="flex items-start gap-4">
                    {/* 類型 */}
                    <div className="flex-shrink-0 grid place-items-center w-10 h-10 rounded-lg bg-slate-800/80 border border-slate-700 text-lg">
                      {kindMeta?.icon || '🗄️'}
                    </div>
                    {/* 表單 */}
                    <div className="flex-1 grid grid-cols-12 gap-3 min-w-0">
                      <Field label="連線名稱" className="col-span-12 sm:col-span-3">
                        <Input value={c.name} onChange={(e) => updateConnector(c.id, { name: e.target.value })} />
                      </Field>
                      <Field label="類型" className="col-span-6 sm:col-span-2">
                        <Select value={c.kind} onChange={(e) => updateConnector(c.id, { kind: e.target.value, status: 'unknown' })}>
                          {DB_KINDS.map((k) => <option key={k.id} value={k.id}>{k.name}</option>)}
                        </Select>
                      </Field>
                      <Field label="Host" className="col-span-6 sm:col-span-3">
                        <Input value={c.host} onChange={(e) => updateConnector(c.id, { host: e.target.value, status: 'unknown' })} className="font-mono" />
                      </Field>
                      <Field label="Port" className="col-span-4 sm:col-span-1">
                        <Input type="number" value={c.port} onChange={(e) => updateConnector(c.id, { port: +e.target.value, status: 'unknown' })} className="font-mono" />
                      </Field>
                      <Field label="Database" className="col-span-8 sm:col-span-3">
                        <Input value={c.database} onChange={(e) => updateConnector(c.id, { database: e.target.value })} className="font-mono" />
                      </Field>
                      <Field label="使用者" className="col-span-6 sm:col-span-3">
                        <Input value={c.username} onChange={(e) => updateConnector(c.id, { username: e.target.value })} className="font-mono" />
                      </Field>
                      <Field label="預設寫入間隔" className="col-span-3 sm:col-span-2" hint="秒">
                        <Input type="number" value={c.default_write_interval_seconds}
                          onChange={(e) => updateConnector(c.id, { default_write_interval_seconds: +e.target.value })}
                          className="font-mono" />
                      </Field>
                      <div className="col-span-3 sm:col-span-2 flex items-end gap-1.5">
                        <Toggle checked={c.enabled} onChange={(v) => updateConnector(c.id, { enabled: v })} label="啟用" size="sm" />
                      </div>
                      <div className="col-span-12 flex items-center gap-2 -mt-1">
                        <StatusChip tone={statusTone}>{statusLabel}</StatusChip>
                        {c.last_check_at && c.status === 'ready' && (
                          <span className="text-[10px] text-slate-500">上次檢查 {new Date(c.last_check_at).toLocaleTimeString()}</span>
                        )}
                        {c.last_check_error && (
                          <span className="text-[10px] text-red-300">{c.last_check_error}</span>
                        )}
                        <div className="ml-auto flex items-center gap-1">
                          <Button size="sm" variant="ghost" onClick={() => testConnector(c.id)}
                            disabled={testingId === c.id}
                            icon={<Icon name={testingId === c.id ? 'refresh' : 'play'} className={clsx('w-3.5 h-3.5', testingId === c.id && 'animate-spin')} />}>
                            {testingId === c.id ? '測試中' : '測試'}
                          </Button>
                          <Button size="sm" variant="ghost" onClick={() => removeConnector(c.id)} className="text-red-300 hover:text-red-200">
                            <Icon name="close" className="w-3.5 h-3.5" />
                          </Button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </SectionCard>

      {/* 2. 時序儲存策略 */}
      <div className="grid grid-cols-12 gap-5">
        <SectionCard
          className="col-span-12 lg:col-span-6"
          title="時序儲存策略"
          subtitle="決定資料庫的時間精度與分區方式"
          icon={<Icon name="spark" className="w-5 h-5" />}
        >
          <div className="grid grid-cols-2 gap-4">
            <Field label="時間精度">
              <Select value={ts.write_precision || 'millisecond'} onChange={(e) => updateNested('timeseries', { write_precision: e.target.value })}>
                <option value="second">秒 (second)</option>
                <option value="millisecond">毫秒 (millisecond)</option>
              </Select>
            </Field>
            <Field label="分區間隔">
              <Select value={ts.partition_interval || 'daily'} onChange={(e) => updateNested('timeseries', { partition_interval: e.target.value })}>
                <option value="daily">每日 (daily)</option>
                <option value="weekly">每週 (weekly)</option>
                <option value="monthly">每月 (monthly)</option>
              </Select>
            </Field>
            <Field label="批次寫入大小" hint="一次最多寫入的筆數">
              <Input type="number" value={ts.batch_size ?? 500}
                onChange={(e) => updateNested('timeseries', { batch_size: +e.target.value })} className="font-mono" />
            </Field>
            <Field label="保留天數" hint="超過天數的資料將被歸檔">
              <Input type="number" value={ts.retention_days ?? 90}
                onChange={(e) => updateNested('timeseries', { retention_days: +e.target.value })} className="font-mono" />
            </Field>
          </div>
        </SectionCard>

        {/* 3. 排程預設 */}
        <SectionCard
          className="col-span-12 lg:col-span-6"
          title="排程預設"
          subtitle="新建立的輪詢群組與點位的預設行為"
          icon={<Icon name="refresh" className="w-5 h-5" />}
        >
          <div className="grid grid-cols-2 gap-4">
            <Field label="預設輪詢間隔" hint="毫秒 (ms)">
              <Input type="number" value={sched.default_interval_ms ?? 1000}
                onChange={(e) => updateNested('scheduler', { default_interval_ms: +e.target.value })} className="font-mono" />
            </Field>
            <Field label="重試次數">
              <Input type="number" value={sched.default_retry_count ?? 3}
                onChange={(e) => updateNested('scheduler', { default_retry_count: +e.target.value })} className="font-mono" />
            </Field>
            <Field label="重試延遲" hint="毫秒 (ms)">
              <Input type="number" value={sched.default_retry_delay_ms ?? 500}
                onChange={(e) => updateNested('scheduler', { default_retry_delay_ms: +e.target.value })} className="font-mono" />
            </Field>
            <Field label="斷路器門檻" hint="連續錯誤次數">
              <Input type="number" value={sched.breaker_threshold ?? 10}
                onChange={(e) => updateNested('scheduler', { breaker_threshold: +e.target.value })} className="font-mono" />
            </Field>
            <div className="col-span-2 flex items-center gap-4 pt-1">
              <Toggle
                checked={sched.auto_start ?? true}
                onChange={(v) => updateNested('scheduler', { auto_start: v })}
                label="開機時自動啟動 collector"
              />
            </div>
          </div>
        </SectionCard>
      </div>

      {/* 4. Local Modbus Share */}
      <SectionCard
        title="Local Modbus Share"
        subtitle="把採集到的 Tag 重新發佈成 Modbus TCP，供下游 SCADA 訂閱"
        icon={<Icon name="cable" className="w-5 h-5" />}
        aside={<Toggle checked={modbusShare.enabled ?? false} onChange={(v) => updateNested('modbus_share', { enabled: v })} />}
      >
        {(modbusShare.enabled ?? false) ? (
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
            <Field label="綁定位址" hint="0.0.0.0 = 全部介面">
              <Input value={modbusShare.bind_address || '0.0.0.0'}
                onChange={(e) => updateNested('modbus_share', { bind_address: e.target.value })} className="font-mono" />
            </Field>
            <Field label="Port">
              <Input type="number" value={modbusShare.port ?? 5020}
                onChange={(e) => updateNested('modbus_share', { port: +e.target.value })} className="font-mono" />
            </Field>
            <Field label="Slave ID">
              <Input type="number" value={modbusShare.slave_id ?? 1}
                onChange={(e) => updateNested('modbus_share', { slave_id: +e.target.value })} className="font-mono" />
            </Field>
            <Field label="起始 Register">
              <Input type="number" value={modbusShare.base_register ?? 40001}
                onChange={(e) => updateNested('modbus_share', { base_register: +e.target.value })} className="font-mono" />
            </Field>
          </div>
        ) : (
          <p className="text-xs text-slate-500">啟用後可將所有 Tag 再次經由 Modbus TCP 提供給其他系統。</p>
        )}
      </SectionCard>

      {/* 5. 介面 / API / 診斷 */}
      <div className="grid grid-cols-12 gap-5">
        <SectionCard
          className="col-span-12 lg:col-span-4"
          title="介面"
          subtitle="主題與語言"
          icon={<Icon name="eye" className="w-5 h-5" />}
        >
          <div className="space-y-3">
            <Field label="主題">
              <Select value={general.theme || 'dark'} onChange={(e) => updateNested('general', { theme: e.target.value })}>
                <option value="dark">深色 (預設)</option>
                <option value="light">淺色</option>
                <option value="auto">跟隨系統</option>
              </Select>
            </Field>
            <Field label="語言">
              <Select value={general.locale || 'zh-TW'} onChange={(e) => updateNested('general', { locale: e.target.value })}>
                <option value="zh-TW">繁體中文</option>
                <option value="en">English</option>
              </Select>
            </Field>
            <Field label="位址顯示格式">
              <Select value={general.addr_format || 'modbus'} onChange={(e) => updateNested('general', { addr_format: e.target.value })}>
                <option value="modbus">Modbus 5 位 (40001)</option>
                <option value="hex">十六進位 (0x9C41)</option>
                <option value="raw">原始 (0)</option>
              </Select>
            </Field>
          </div>
        </SectionCard>

        <SectionCard
          className="col-span-12 lg:col-span-4"
          title="API 端點"
          subtitle="後端 Gateway 服務"
          icon={<Icon name="flow" className="w-5 h-5" />}
        >
          <div className="space-y-3">
            <Field label="Base URL" hint="不含尾斜線">
              <Input value={general.api_base || 'http://localhost:8080'}
                onChange={(e) => updateNested('general', { api_base: e.target.value })} className="font-mono" />
            </Field>
            <Field label="API 版本">
              <Select value={general.api_version || 'v1'} onChange={(e) => updateNested('general', { api_version: e.target.value })}>
                <option value="v1">v1</option>
                <option value="v2">v2 (beta)</option>
              </Select>
            </Field>
            <Field label="請求逾時" hint="秒">
              <Input type="number" value={general.timeout_seconds ?? 30}
                onChange={(e) => updateNested('general', { timeout_seconds: +e.target.value })} className="font-mono" />
            </Field>
          </div>
        </SectionCard>

        <SectionCard
          className="col-span-12 lg:col-span-4"
          title="診斷"
          subtitle="日誌與 SSE 串流"
          icon={<Icon name="alert" className="w-5 h-5" />}
        >
          <div className="space-y-3">
            <Field label="日誌等級">
              <Select value={general.log_level || 'info'} onChange={(e) => updateNested('general', { log_level: e.target.value })}>
                <option value="trace">trace</option>
                <option value="debug">debug</option>
                <option value="info">info (預設)</option>
                <option value="warn">warn</option>
                <option value="error">error</option>
              </Select>
            </Field>
            <Field label="SSE 心跳" hint="秒">
              <Input type="number" value={general.sse_heartbeat_seconds ?? 15}
                onChange={(e) => updateNested('general', { sse_heartbeat_seconds: +e.target.value })} className="font-mono" />
            </Field>
            <div className="pt-1 space-y-2">
              <Toggle checked={general.enable_debug_panel ?? false}
                onChange={(v) => updateNested('general', { enable_debug_panel: v })} label="顯示 Debug 面板" />
              <Toggle checked={general.enable_audit_log ?? true}
                onChange={(v) => updateNested('general', { enable_audit_log: v })} label="操作審計記錄" />
            </div>
          </div>
        </SectionCard>
      </div>

      {/* 底部儲存 */}
      <div className="sticky bottom-4 z-10 flex items-center justify-between rounded-xl border border-slate-700/60 bg-slate-900/95 backdrop-blur px-4 py-3 shadow-2xl shadow-black/40">
        <div className="text-[11px] text-slate-500 flex items-center gap-2">
          <Icon name="info" className="w-3.5 h-3.5" />
          設定會立即套用，並於下次重啟後生效。
        </div>
        <div className="flex items-center gap-2">
          <Button variant="ghost" size="sm">重設為預設</Button>
          <Button variant="success" size="md" icon={<Icon name="save" className="w-3.5 h-3.5" />}>
            儲存所有設定
          </Button>
        </div>
      </div>
    </div>
  );
}

window.SettingsPage = SettingsPage;
