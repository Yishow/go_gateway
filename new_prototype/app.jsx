// 主 App - Datalink Workbench MVP
// 4 步驟流程：新增裝置 → 接入規則 → 點位映射 → 儲存資料庫  + 設定頁

const TWEAK_DEFAULTS = /*EDITMODE-BEGIN*/{
  "showSummaryRail": true,
  "sidebarCollapsed": false,
  "accent": "blue"
}/*EDITMODE-END*/;

const STEPS = [
  { id: 1, key: 'device', title: '新增裝置', subtitle: 'Device + Protocol', icon: 'device' },
  { id: 2, key: 'rule', title: '接入規則', subtitle: 'Source Rules · 多條範圍', icon: 'rule' },
  { id: 3, key: 'mapping', title: '點位映射', subtitle: 'Point → Tag', icon: 'map' },
  { id: 4, key: 'database', title: '儲存資料庫', subtitle: 'DB Target · Commit', icon: 'db' },
];

// ──────────────────────────────────────────────────────────────────
// 左側導航 (步驟 + 設定)
// ──────────────────────────────────────────────────────────────────
function StepRailV2({ view, current, completed, collapsed, onJump, onSwitchView }) {
  return (
    <nav className="flex flex-col gap-1">
      {!collapsed && (
        <div className="px-1 mb-1 text-[10px] font-semibold uppercase tracking-[0.18em] text-slate-500">建置流程</div>
      )}
      {STEPS.map((s) => {
        const isCurrent = view === 'flow' && current === s.id;
        const isDone = completed.has(s.id);
        const isReachable = isDone || isCurrent || completed.has(s.id - 1) || s.id === 1;
        return (
          <button
            key={s.id}
            onClick={() => isReachable && (onSwitchView('flow'), onJump(s.id))}
            disabled={!isReachable}
            title={collapsed ? s.title : undefined}
            className={clsx(
              'relative flex items-center gap-3 rounded-xl border text-left transition-all w-full',
              collapsed ? 'p-1.5 justify-center' : 'px-3 py-2.5',
              isCurrent
                ? 'border-blue-500/40 bg-blue-500/10 shadow-lg shadow-blue-500/5'
                : isDone
                  ? 'border-emerald-500/20 bg-emerald-500/[0.04] hover:bg-emerald-500/[0.08]'
                  : !isReachable
                    ? 'border-slate-800 bg-slate-900/30 opacity-60 cursor-not-allowed'
                    : 'border-slate-800 bg-slate-900/30 hover:bg-slate-800/50',
            )}
          >
            <span
              className={clsx(
                'grid place-items-center w-9 h-9 rounded-lg border flex-shrink-0',
                isCurrent
                  ? 'border-blue-500/50 bg-blue-500/20 text-blue-200'
                  : isDone
                    ? 'border-emerald-500/40 bg-emerald-500/15 text-emerald-300'
                    : 'border-slate-700 bg-slate-800 text-slate-500',
              )}
            >
              {isDone && !isCurrent ? <Icon name="check" className="w-4 h-4" /> : <Icon name={s.icon} className="w-4 h-4" />}
            </span>
            {!collapsed && (
              <>
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-1.5">
                    <span className="text-[10px] font-mono text-slate-500">0{s.id}</span>
                    <span className={clsx('text-sm font-semibold truncate', isCurrent ? 'text-blue-100' : 'text-slate-200')}>{s.title}</span>
                  </div>
                  <div className="text-[11px] text-slate-500 truncate">{s.subtitle}</div>
                </div>
                {isCurrent && <span className="w-1.5 h-1.5 rounded-full bg-blue-400 pulse-dot" />}
              </>
            )}
          </button>
        );
      })}

      <div className={clsx('my-2 border-t border-slate-800', collapsed ? 'mx-1' : 'mx-1')} />

      {!collapsed && (
        <div className="px-1 mb-1 text-[10px] font-semibold uppercase tracking-[0.18em] text-slate-500">系統</div>
      )}

      <button
        onClick={() => onSwitchView('settings')}
        title={collapsed ? '設定' : undefined}
        className={clsx(
          'relative flex items-center gap-3 rounded-xl border text-left transition-all w-full',
          collapsed ? 'p-1.5 justify-center' : 'px-3 py-2.5',
          view === 'settings'
            ? 'border-blue-500/40 bg-blue-500/10 shadow-lg shadow-blue-500/5'
            : 'border-slate-800 bg-slate-900/30 hover:bg-slate-800/50',
        )}
      >
        <span
          className={clsx(
            'grid place-items-center w-9 h-9 rounded-lg border flex-shrink-0',
            view === 'settings'
              ? 'border-blue-500/50 bg-blue-500/20 text-blue-200'
              : 'border-slate-700 bg-slate-800 text-slate-500',
          )}
        >
          <Icon name="sliders" className="w-4 h-4" />
        </span>
        {!collapsed && (
          <div className="min-w-0 flex-1">
            <div className="text-sm font-semibold truncate text-slate-200">設定</div>
            <div className="text-[11px] text-slate-500 truncate">資料庫、排程、介面…</div>
          </div>
        )}
        {view === 'settings' && !collapsed && <span className="w-1.5 h-1.5 rounded-full bg-blue-400 pulse-dot" />}
      </button>
    </nav>
  );
}

// ──────────────────────────────────────────────────────────────────
// 右側即時設定摘要
// ──────────────────────────────────────────────────────────────────
function SummaryRail({ state }) {
  const { devices = [], rules = [], points = [], mappings = {}, db = {} } = state;
  const enabledPoints = points.filter((p) => !p.skipped);
  const connector = db.connector;

  const Row = ({ label, value, mono, dim }) =>
    <div className="flex items-baseline justify-between gap-3 py-1">
      <span className="text-[10px] uppercase tracking-wider text-slate-500 flex-shrink-0">{label}</span>
      <span className={clsx('text-xs text-right truncate', mono && 'font-mono', dim ? 'text-slate-500' : 'text-slate-200')}>{value}</span>
    </div>;

  const totalRanges = rules.filter((r) => r.enabled).length;
  const totalPlanned = rules.filter((r) => r.enabled).reduce((sum, r) => sum + r.count - (r.skipped_addresses?.length || 0), 0);
  const testedDevs = devices.filter((d) => d.test?.status === 'success').length;

  return (
    <div className="space-y-3">
      <div className="rounded-xl border border-slate-700/60 bg-slate-900/40 p-3">
        <div className="flex items-center gap-2 mb-2">
          <Icon name="device" className="w-3.5 h-3.5 text-slate-400" />
          <span className="text-[11px] font-semibold uppercase tracking-wider text-slate-400">設備</span>
        </div>
        <Row label="總數" value={`${devices.length} 個`} mono />
        <Row label="已測試" value={`${testedDevs} / ${devices.length}`} mono dim={testedDevs !== devices.length} />
        <div className="mt-1.5 space-y-0.5">
          {devices.slice(0, 4).map((d, i) => {
            const colors = ['bg-blue-500', 'bg-emerald-500', 'bg-amber-500', 'bg-fuchsia-500'];
            return (
              <div key={d.id} className="flex items-center gap-1.5 text-[10px] font-mono text-slate-500">
                <span className={clsx('w-1.5 h-1.5 rounded-full flex-shrink-0', colors[i % colors.length])} />
                <span className="truncate text-slate-400">{d.name}</span>
                <span className="ml-auto text-slate-500">{d.protocol}</span>
              </div>
            );
          })}
          {devices.length > 4 && <div className="text-[10px] text-slate-600 pl-3">+ {devices.length - 4} 更多</div>}
        </div>
      </div>

      <div className="rounded-xl border border-slate-700/60 bg-slate-900/40 p-3">
        <div className="flex items-center gap-2 mb-2">
          <Icon name="rule" className="w-3.5 h-3.5 text-slate-400" />
          <span className="text-[11px] font-semibold uppercase tracking-wider text-slate-400">接入規則</span>
        </div>
        <Row label="條數" value={`${totalRanges} 條`} mono />
        <Row label="點位" value={`${enabledPoints.length || totalPlanned} 個`} mono />
      </div>

      <div className="rounded-xl border border-slate-700/60 bg-slate-900/40 p-3">
        <div className="flex items-center gap-2 mb-2">
          <Icon name="map" className="w-3.5 h-3.5 text-slate-400" />
          <span className="text-[11px] font-semibold uppercase tracking-wider text-slate-400">映射</span>
        </div>
        <Row label="Tags" value={Object.values(mappings).filter((m) => m.enabled).length || '—'} mono />
        <Row label="個別轉換" value={Object.values(mappings).filter((m) => m.scale !== 1 || m.offset !== 0).length} mono />
      </div>

      <div className="rounded-xl border border-slate-700/60 bg-slate-900/40 p-3">
        <div className="flex items-center gap-2 mb-2">
          <Icon name="db" className="w-3.5 h-3.5 text-slate-400" />
          <span className="text-[11px] font-semibold uppercase tracking-wider text-slate-400">資料庫寫入</span>
        </div>
        {connector ? (
          <>
            <Row label="類型" value={connector.kind} mono />
            <Row label="表" value={`${connector.schema}.${connector.table}`} mono />
            <Row label="間隔" value={`${connector.write_interval_seconds}s`} mono />
            <Row label="目標" value={`${Object.values(db.targets || {}).filter((t) => t.enabled).length} 欄位`} mono />
          </>
        ) : (
          <div className="text-[11px] text-slate-500 italic">尚未設定</div>
        )}
      </div>

      {state.committed && (
        <div className="rounded-xl border border-emerald-500/30 bg-emerald-500/[0.08] p-3 sweep-in">
          <div className="flex items-center gap-2 text-emerald-200 text-xs font-semibold">
            <Icon name="check" className="w-3.5 h-3.5" />
            已部署到 Runtime
          </div>
          <div className="text-[10px] text-emerald-300/70 mt-1">所有設定已成功寫入。</div>
        </div>
      )}
    </div>
  );
}

// ──────────────────────────────────────────────────────────────────
// App
// ──────────────────────────────────────────────────────────────────
function App() {
  const [tweaks, setTweak] = useTweaks(TWEAK_DEFAULTS);
  const [view, setView] = React.useState('flow'); // 'flow' | 'settings'
  const [current, setCurrent] = React.useState(1);
  const [completed, setCompleted] = React.useState(new Set());
  const sidebarCollapsed = tweaks.sidebarCollapsed;
  const setSidebarCollapsed = (v) => setTweak('sidebarCollapsed', v);

  const [state, setState] = React.useState(() => {
    const dev0 = makeDefaultDevice({
      name: DEFAULTS.device.name,
      description: DEFAULTS.device.description,
      protocol: DEFAULTS.device.protocol,
      config: { ...DEFAULTS.device.config },
    });
    return {
      devices: [dev0],
      rules: [
        makeDefaultRule({
          device_id: dev0.id,
          name: 'Holding Registers',
          start_address: '40001',
          count: 8,
          data_type: 'int16',
          naming_prefix: 'SENSOR_',
          scale_multiplier: 0.1,
        }),
      ],
      points: [],
      mappings: {},
      db: {},
      settings: {
        connectors: [
          {
            id: 'conn-prod',
            name: 'TimeSeries Prod',
            kind: 'postgres',
            host: 'tsdb.internal',
            port: 5432,
            database: 'gateway_metrics',
            username: 'gw_writer',
            enabled: true,
            status: 'ready',
            last_check_at: new Date().toISOString(),
            default_write_interval_seconds: 5,
          },
        ],
        timeseries: { write_precision: 'millisecond', partition_interval: 'daily', batch_size: 500, retention_days: 90 },
        scheduler: { default_interval_ms: 1000, default_retry_count: 3, default_retry_delay_ms: 500, breaker_threshold: 10, auto_start: true },
        modbus_share: { enabled: true, bind_address: '0.0.0.0', port: 5020, slave_id: 1, base_register: 40001 },
        general: { theme: 'dark', locale: 'zh-TW', addr_format: 'modbus', api_base: 'http://localhost:8080', api_version: 'v1', timeout_seconds: 30, log_level: 'info', sse_heartbeat_seconds: 15, enable_debug_panel: false, enable_audit_log: true },
      },
      committed: false,
    };
  });

  const goNext = () => {
    setCompleted((s) => new Set(s).add(current));
    setCurrent((c) => Math.min(STEPS.length, c + 1));
    document.getElementById('step-content')?.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const goBack = () => {
    setCurrent((c) => Math.max(1, c - 1));
    document.getElementById('step-content')?.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const stepMeta = STEPS.find((s) => s.id === current);

  // 鍵盤捷徑
  React.useEffect(() => {
    const onKey = (e) => {
      if ((e.ctrlKey || e.metaKey) && e.key === 'b') {
        e.preventDefault();
        setSidebarCollapsed(!sidebarCollapsed);
      }
    };
    window.addEventListener('keydown', onKey);
    return () => window.removeEventListener('keydown', onKey);
  }, [sidebarCollapsed]);

  return (
    <div className="min-h-screen flex flex-col">
      {/* Top Bar */}
      <header className="border-b border-slate-800/80 bg-slate-950/60 backdrop-blur-md sticky top-0 z-30">
        <div className="px-5 lg:px-7 py-3 flex items-center gap-4">
          <button
            onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
            className="grid place-items-center w-8 h-8 rounded-lg border border-slate-700/60 text-slate-400 hover:text-slate-100 hover:bg-slate-800/60 transition-colors"
            title={sidebarCollapsed ? '展開側邊欄 (⌘B)' : '收合側邊欄 (⌘B)'}
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.75" className="w-4 h-4">
              <rect x="3" y="4" width="18" height="16" rx="2" />
              <line x1="9" y1="4" x2="9" y2="20" />
              {!sidebarCollapsed && <line x1="6" y1="8" x2="6" y2="8.01" strokeWidth="2" />}
            </svg>
          </button>

          <div className="flex items-center gap-2.5">
            <div className="grid place-items-center w-9 h-9 rounded-lg bg-gradient-to-br from-blue-500/30 to-emerald-500/20 border border-blue-500/30">
              <Icon name="bolt" className="w-4 h-4 text-blue-300" />
            </div>
            <div>
              <div className="text-sm font-semibold tracking-tight text-slate-100 leading-tight">Datalink Workbench</div>
              <div className="text-[10px] text-slate-500 leading-tight">go_gateway · v1 接入點位設定流程</div>
            </div>
          </div>

          <nav className="hidden md:flex items-center gap-1.5 ml-6 text-xs">
            <span className="text-slate-500">Datalink</span>
            <Icon name="chevron" className="w-3 h-3 text-slate-600" />
            <span className="text-slate-500">Workbench</span>
            <Icon name="chevron" className="w-3 h-3 text-slate-600" />
            <span className="text-slate-200 font-medium">{view === 'settings' ? '設定' : stepMeta?.title}</span>
          </nav>

          <div className="ml-auto flex items-center gap-2">
            <Button variant="ghost" size="sm" icon={<Icon name="save" className="w-3.5 h-3.5" />}>
              儲存草稿
            </Button>
            <Button variant="secondary" size="sm">取消</Button>
            <div className="h-6 w-px bg-slate-700/60 mx-1" />
            <div className="flex items-center gap-1.5 text-xs text-slate-400">
              <span className="w-1.5 h-1.5 rounded-full bg-emerald-400 pulse-dot" />
              <span>scheduler {state.committed ? 'running' : 'idle'}</span>
            </div>
          </div>
        </div>
      </header>

      {/* Body - 使用 flex 以支援收合 */}
      <div className="flex-1 flex gap-5 px-5 lg:px-7 py-5">
        {/* Left rail - 動態寬度 */}
        <aside
          className="min-w-0 overflow-hidden transition-all duration-200"
          style={{ flex: `0 0 ${sidebarCollapsed ? 64 : 232}px`, width: sidebarCollapsed ? 64 : 232 }}
        >
          <div className="sticky top-[68px] space-y-4">
            <StepRailV2
              view={view}
              current={current}
              completed={completed}
              collapsed={sidebarCollapsed}
              onJump={(id) => setCurrent(id)}
              onSwitchView={setView}
            />

            {!sidebarCollapsed && view === 'flow' && (
              <div className="rounded-xl border border-slate-700/40 bg-slate-900/30 p-3">
                <div className="text-[10px] uppercase tracking-wider text-slate-500 mb-1.5">本流程</div>
                <p className="text-[11px] leading-relaxed text-slate-400">
                  由「設備 → 接入規則 → 點位/映射 → 資料庫」四個階段組成，全部成功後才會啟動 collector 並寫入。
                </p>
              </div>
            )}
          </div>
        </aside>

        {/* Center content */}
        <main id="step-content" className="flex-1 min-w-0">
          {view === 'flow' ? (
            <>
              <div className="mb-4 flex items-center justify-between">
                <div>
                  <div className="flex items-center gap-2">
                    <span className="text-[11px] font-mono text-slate-500">STEP 0{current} / 04</span>
                    <span className="h-3 w-px bg-slate-700" />
                    <span className="text-[11px] text-slate-500">{stepMeta?.subtitle}</span>
                  </div>
                  <h2 className="text-xl font-semibold text-slate-100 mt-0.5">{stepMeta?.title}</h2>
                </div>
                <div className="hidden sm:flex items-center gap-1.5">
                  {STEPS.map((s) => (
                    <span
                      key={s.id}
                      className={clsx(
                        'h-1 rounded-full transition-all',
                        s.id < current ? 'w-8 bg-emerald-500/80' : s.id === current ? 'w-12 bg-blue-500' : 'w-8 bg-slate-700',
                      )}
                    />
                  ))}
                </div>
              </div>

              {current === 1 && <Step1Device state={state} setState={setState} onContinue={goNext} />}
              {current === 2 && <Step2Rule state={state} setState={setState} onContinue={goNext} />}
              {current === 3 && <Step3Mapping state={state} setState={setState} onContinue={goNext} />}
              {current === 4 && <Step4Database state={state} setState={setState} onCommit={() => alert('Runtime Dashboard 跳轉 (示意)')} />}

              <div className="mt-6 flex items-center justify-between pt-4 border-t border-slate-800/80">
                <Button onClick={goBack} disabled={current === 1} variant="ghost" icon={<Icon name="chevron" className="w-4 h-4 rotate-180" />}>
                  上一步
                </Button>
                <div className="flex items-center gap-2 text-[11px] text-slate-500">
                  <span className="kbd">Ctrl</span><span>+</span><span className="kbd">B</span>
                  <span>收合側欄</span>
                </div>
                <div className="invisible"><Button>placeholder</Button></div>
              </div>
            </>
          ) : (
            <SettingsPage state={state} setState={setState} />
          )}
        </main>

        {/* Right summary rail */}
        {tweaks.showSummaryRail && view === 'flow' && (
          <aside className="hidden xl:block flex-shrink-0 w-[240px]">
            <div className="sticky top-[68px]">
              <div className="text-[10px] font-semibold uppercase tracking-wider text-slate-500 mb-2 px-1">即時設定摘要</div>
              <SummaryRail state={state} />
            </div>
          </aside>
        )}
      </div>

      {/* Tweaks panel */}
      <TweaksPanel title="Tweaks">
        <TweakSection title="版面">
          <TweakToggle label="收合側邊欄" value={tweaks.sidebarCollapsed} onChange={(v) => setTweak('sidebarCollapsed', v)} />
          <TweakToggle label="顯示右側摘要欄" value={tweaks.showSummaryRail} onChange={(v) => setTweak('showSummaryRail', v)} />
        </TweakSection>
        <TweakSection title="操作">
          <TweakButton label="重置流程" onClick={() => {
            setCurrent(1);
            setView('flow');
            setCompleted(new Set());
            setState((s) => {
              const dev0 = makeDefaultDevice({ name: DEFAULTS.device.name, protocol: DEFAULTS.device.protocol, config: { ...DEFAULTS.device.config } });
              return {
                ...s,
                devices: [dev0],
                rules: [makeDefaultRule({ device_id: dev0.id, name: 'Holding Registers', start_address: '40001', count: 8, naming_prefix: 'SENSOR_', scale_multiplier: 0.1 })],
                points: [],
                mappings: {},
                db: {},
                committed: false,
                commit: undefined,
              };
            });
          }} />
          <TweakButton label="加入第二台設備 (示範)" onClick={() => {
            setState((s) => {
              const newDev = makeDefaultDevice({ name: 'PLC-生產線-02', protocol: 'mc_3e', config: { host: '192.168.1.101', port: 5007, slave_id: 1, timeout: 5 } });
              return {
                ...s,
                devices: [...s.devices, newDev],
                rules: [
                  ...s.rules,
                  makeDefaultRule({
                    device_id: newDev.id,
                    name: 'Mitsubishi D-Reg',
                    start_address: '00100',
                    count: 4,
                    data_type: 'int16',
                    naming_prefix: 'MC_',
                    scale_multiplier: 1,
                  }),
                ],
              };
            });
            setCurrent(1);
            setView('flow');
          }} />
          <TweakButton label="加入第二條規則 (示範)" onClick={() => {
            setState((s) => ({
              ...s,
              rules: [
                ...s.rules,
                makeDefaultRule({
                  device_id: s.devices[0]?.id,
                  name: 'Coils (DI)',
                  start_address: '00001',
                  count: 4,
                  data_type: 'bool',
                  naming_prefix: 'STATUS_',
                  scale_multiplier: 1,
                  scale_offset: 0,
                }),
              ],
            }));
            setCurrent(2);
            setView('flow');
          }} />
          <TweakButton label="跳到 Step 4 (含示例資料)" onClick={() => {
            setCompleted(new Set([1, 2, 3]));
            setView('flow');
            setState((s) => {
              const points = deriveAllPoints(s.rules, s.devices[0]?.id);
              return {
                ...s,
                devices: s.devices.map((d) => ({ ...d, test: d.test || { status: 'success', latency_ms: 38, tested_at: new Date() } })),
                points,
              };
            });
            setCurrent(4);
          }} />
        </TweakSection>
      </TweaksPanel>
    </div>
  );
}

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(<App />);
