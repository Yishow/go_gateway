// Step 1: 新增裝置 - 支援多裝置 (tab 列管理)

function Step1Device({ state, setState, onContinue }) {
  const devices = state.devices || [];
  const [selectedDeviceId, setSelectedDeviceId] = React.useState(() => devices[0]?.id);
  const [testing, setTesting] = React.useState(false);

  React.useEffect(() => {
    if (!devices.find((d) => d.id === selectedDeviceId)) {
      setSelectedDeviceId(devices[0]?.id);
    }
  }, [devices, selectedDeviceId]);

  const selectedDevice = devices.find((d) => d.id === selectedDeviceId) || devices[0];

  const updateDevice = (patch) =>
    setState((s) => ({
      ...s,
      devices: s.devices.map((d) => (d.id === selectedDeviceId ? { ...d, ...patch } : d)),
    }));

  const updateConfig = (patch) =>
    setState((s) => ({
      ...s,
      devices: s.devices.map((d) => (d.id === selectedDeviceId ? { ...d, config: { ...d.config, ...patch } } : d)),
    }));

  const renameDevice = (id, name) =>
    setState((s) => ({ ...s, devices: s.devices.map((d) => (d.id === id ? { ...d, name } : d)) }));

  const addDevice = () => {
    const newDev = makeDefaultDevice({
      name: `設備 ${devices.length + 1}`,
      description: '',
      protocol: 'modbus_tcp',
      config: { host: '192.168.1.10' + devices.length, port: 502, slave_id: 1, timeout: 5 },
    });
    setState((s) => ({ ...s, devices: [...s.devices, newDev] }));
    setSelectedDeviceId(newDev.id);
  };

  const removeDevice = (id) => {
    if (devices.length <= 1) return;
    if (!confirm('刪除此設備將同時移除所屬規則。是否繼續?')) return;
    setState((s) => ({
      ...s,
      devices: s.devices.filter((d) => d.id !== id),
      rules: s.rules.filter((r) => r.device_id !== id),
    }));
  };

  // 裝置顏色 (與 Step 2 規則色盤共用)
  const deviceColors = ['blue', 'emerald', 'amber', 'fuchsia', 'cyan', 'rose'];
  const colorFor = (did) => deviceColors[devices.findIndex((d) => d.id === did) % deviceColors.length];
  const cs = {
    blue: { bg: 'bg-blue-500/15', border: 'border-blue-500/40', text: 'text-blue-100', solid: 'bg-blue-500' },
    emerald: { bg: 'bg-emerald-500/15', border: 'border-emerald-500/40', text: 'text-emerald-100', solid: 'bg-emerald-500' },
    amber: { bg: 'bg-amber-500/15', border: 'border-amber-500/40', text: 'text-amber-100', solid: 'bg-amber-500' },
    fuchsia: { bg: 'bg-fuchsia-500/15', border: 'border-fuchsia-500/40', text: 'text-fuchsia-100', solid: 'bg-fuchsia-500' },
    cyan: { bg: 'bg-cyan-500/15', border: 'border-cyan-500/40', text: 'text-cyan-100', solid: 'bg-cyan-500' },
    rose: { bg: 'bg-rose-500/15', border: 'border-rose-500/40', text: 'text-rose-100', solid: 'bg-rose-500' },
  };

  if (!selectedDevice) {
    return (
      <div className="text-center py-20 text-slate-500">
        尚無設備。<button onClick={addDevice} className="text-blue-300 hover:underline">新增第一個設備</button>
      </div>
    );
  }

  const device = selectedDevice;

  const protocolMeta = PROTOCOLS.find((p) => p.id === device.protocol) || PROTOCOLS[0];
  const isSerial = device.protocol === 'modbus_rtu';
  const isMqtt = device.protocol === 'mqtt';

  const runTest = () => {
    setTesting(true);
    updateDevice({ test: { status: 'running', stages: {} } });
    const stages = isMqtt
      ? ['resolve', 'connect', 'subscribe']
      : isSerial
        ? ['open_port', 'handshake', 'probe']
        : ['resolve', 'tcp', 'modbus_probe'];

    let idx = 0;
    const stageResults = {};
    const tick = () => {
      if (idx >= stages.length) {
        updateDevice({
          status: 'tested',
          test: {
            status: 'success',
            latency_ms: 42 + Math.floor(Math.random() * 18),
            stages: stageResults,
            tested_at: new Date(),
          },
        });
        setTesting(false);
        return;
      }
      const stage = stages[idx];
      stageResults[stage] = { status: 'success', latency_ms: 10 + Math.floor(Math.random() * 18) };
      updateDevice({ test: { status: 'running', stages: { ...stageResults }, progress: (idx + 1) / stages.length } });
      idx++;
      setTimeout(tick, 380);
    };
    setTimeout(tick, 220);
  };

  const test = device.test || {};
  const allTested = devices.every((d) => d.test?.status === 'success');
  const canContinue = allTested && devices.every((d) => d.name && d.protocol);
  const testedCount = devices.filter((d) => d.test?.status === 'success').length;

  return (
    <div className="space-y-5 sweep-in">
      {/* ── 設備 Tab 列 ─────────────────────────────────── */}
      <SectionCard
        title="設備列表"
        subtitle="可同時設定多個設備，每個設備有各自的協議與連線參數"
        icon={<Icon name="device" className="w-5 h-5" />}
        aside={
          <StatusChip tone={testedCount === devices.length ? 'success' : 'info'}>
            {devices.length} 個設備 · {testedCount} 已測試
          </StatusChip>
        }
        contentClassName="p-3"
      >
        <div className="flex items-stretch gap-2 overflow-x-auto pb-1">
          {devices.map((d) => {
            const cc = cs[colorFor(d.id)];
            const isActive = d.id === selectedDeviceId;
            const tStatus = d.test?.status === 'success' ? 'success' : d.test?.status === 'running' ? 'info' : 'draft';
            const tLabel = d.test?.status === 'success' ? `✓ ${d.test.latency_ms}ms` : d.test?.status === 'running' ? '測試中…' : '未測試';
            return (
              <div
                key={d.id}
                onClick={() => setSelectedDeviceId(d.id)}
                className={clsx(
                  'group relative flex-shrink-0 cursor-pointer rounded-lg border px-3 py-2 transition-all min-w-[220px]',
                  isActive ? `${cc.border} ${cc.bg} shadow-md` : 'border-slate-700/60 bg-slate-900/40 hover:border-slate-600 hover:bg-slate-800/60',
                )}
              >
                <div className="flex items-center gap-2">
                  <span className={clsx('w-2 h-2 rounded-full flex-shrink-0', cc.solid)} />
                  <input
                    value={d.name}
                    onChange={(e) => renameDevice(d.id, e.target.value)}
                    onClick={(e) => e.stopPropagation()}
                    className={clsx(
                      'bg-transparent border-none outline-none text-sm font-semibold w-full',
                      isActive ? 'text-slate-100' : 'text-slate-300',
                    )}
                  />
                  {devices.length > 1 && (
                    <button
                      onClick={(e) => { e.stopPropagation(); removeDevice(d.id); }}
                      className="opacity-0 group-hover:opacity-100 transition-opacity text-slate-500 hover:text-red-400 p-0.5"
                      title="刪除設備"
                    >
                      <Icon name="close" className="w-3.5 h-3.5" />
                    </button>
                  )}
                </div>
                <div className="mt-1 flex items-center justify-between text-[10px] font-mono">
                  <span className="text-slate-500">{d.protocol} · {d.config?.host || '—'}</span>
                  <StatusChip tone={tStatus} dot={false}>{tLabel}</StatusChip>
                </div>
              </div>
            );
          })}
          <button
            onClick={addDevice}
            className="flex-shrink-0 flex items-center gap-1.5 px-3 py-2 rounded-lg border border-dashed border-slate-700 text-slate-400 hover:border-blue-500/40 hover:text-blue-300 hover:bg-blue-500/[0.04] transition-all text-sm"
          >
            <Icon name="plus" className="w-4 h-4" />
            新增設備
          </button>
        </div>
      </SectionCard>

      {/* ── 編輯選中設備 ──────────────────────────────────── */}
      <div className="grid grid-cols-12 gap-5">
        {/* 左：基本資料 */}
        <SectionCard
          className="col-span-12 xl:col-span-7"
          title={`編輯：${device.name}`}
          subtitle="名稱、描述、通訊協議"
          icon={<Icon name="device" className="w-5 h-5" />}
        >
          <div className="grid grid-cols-2 gap-4">
            <Field label="設備名稱" required className="col-span-2 sm:col-span-1">
              <Input value={device.name} onChange={(e) => updateDevice({ name: e.target.value })} placeholder="例如：PLC-生產線-01" />
            </Field>
            <Field label="描述" className="col-span-2 sm:col-span-1">
              <Input value={device.description || ''} onChange={(e) => updateDevice({ description: e.target.value })} placeholder="選填說明" />
            </Field>

            <div className="col-span-2">
              <span className="label">通訊協議</span>
              <div className="grid grid-cols-2 md:grid-cols-3 gap-2">
                {PROTOCOLS.map((p) => {
                  const active = device.protocol === p.id;
                  return (
                    <button
                      key={p.id}
                      type="button"
                      onClick={() => updateDevice({ protocol: p.id, test: null, status: 'draft' })}
                      className={clsx(
                        'group relative rounded-lg border px-3 py-2.5 text-left transition-all',
                        active ? 'border-blue-500/60 bg-blue-500/10 shadow-md shadow-blue-500/10' : 'border-slate-700/60 bg-slate-900/40 hover:border-slate-600 hover:bg-slate-800/60',
                      )}
                    >
                      <div className="flex items-center justify-between">
                        <span className={clsx('text-sm font-semibold', active ? 'text-blue-200' : 'text-slate-200')}>{p.name}</span>
                        {active && <Icon name="check" className="w-4 h-4 text-blue-300" />}
                      </div>
                      <div className="mt-0.5 text-[11px] text-slate-500">{p.desc}</div>
                    </button>
                  );
                })}
              </div>
            </div>

            <div className="col-span-2 mt-2 rounded-xl border border-slate-700/60 bg-slate-950/40 p-4">
              <div className="mb-3 flex items-center justify-between">
                <div className="text-xs font-semibold uppercase tracking-wider text-slate-400">
                  連線參數 · {protocolMeta.name}
                </div>
                <span className="text-[11px] text-slate-500">connection_config (JSON)</span>
              </div>
              {!isSerial && !isMqtt && (
                <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
                  <Field label="Host" required>
                    <Input value={device.config.host || ''} onChange={(e) => updateConfig({ host: e.target.value })} placeholder="192.168.1.100" className="font-mono" />
                  </Field>
                  <Field label="Port" required>
                    <Input type="number" value={device.config.port || ''} onChange={(e) => updateConfig({ port: +e.target.value })} className="font-mono" />
                  </Field>
                  <Field label="Slave ID">
                    <Input type="number" value={device.config.slave_id || ''} onChange={(e) => updateConfig({ slave_id: +e.target.value })} className="font-mono" />
                  </Field>
                  <Field label="Timeout (s)">
                    <Input type="number" value={device.config.timeout || ''} onChange={(e) => updateConfig({ timeout: +e.target.value })} className="font-mono" />
                  </Field>
                </div>
              )}
              {isSerial && (
                <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
                  <Field label="Port"><Input value={device.config.port || '/dev/ttyUSB0'} onChange={(e) => updateConfig({ port: e.target.value })} className="font-mono" /></Field>
                  <Field label="Baud Rate"><Input type="number" value={device.config.baud || 9600} onChange={(e) => updateConfig({ baud: +e.target.value })} className="font-mono" /></Field>
                  <Field label="Parity"><Select value={device.config.parity || 'N'} onChange={(e) => updateConfig({ parity: e.target.value })}><option>N</option><option>E</option><option>O</option></Select></Field>
                  <Field label="Slave ID"><Input type="number" value={device.config.slave_id || 1} onChange={(e) => updateConfig({ slave_id: +e.target.value })} className="font-mono" /></Field>
                </div>
              )}
              {isMqtt && (
                <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                  <Field label="Broker"><Input value={device.config.broker || 'mqtts://broker.local:8883'} onChange={(e) => updateConfig({ broker: e.target.value })} className="font-mono" /></Field>
                  <Field label="Username"><Input value={device.config.username || ''} onChange={(e) => updateConfig({ username: e.target.value })} className="font-mono" /></Field>
                  <Field label="Client ID"><Input value={device.config.client_id || 'gw-01'} onChange={(e) => updateConfig({ client_id: e.target.value })} className="font-mono" /></Field>
                </div>
              )}
            </div>
          </div>
        </SectionCard>

        {/* 右：測試連線 */}
        <SectionCard
          className="col-span-12 xl:col-span-5"
          title="連線測試"
          subtitle="啟用前必須通過 readiness 檢查"
          icon={<Icon name="cable" className="w-5 h-5" />}
          aside={
            <StatusChip tone={test.status === 'success' ? 'success' : test.status === 'running' ? 'info' : 'draft'}>
              {test.status === 'success' ? '已通過' : test.status === 'running' ? '測試中…' : '尚未測試'}
            </StatusChip>
          }
        >
          <div className="space-y-4">
            <div className="rounded-lg border border-slate-700/60 bg-slate-950/60 p-3">
              <div className="text-[11px] uppercase tracking-wider text-slate-500 mb-2">送出 payload</div>
              <pre className="text-[11px] font-mono leading-relaxed text-slate-300 overflow-x-auto">
{`POST /api/v1/datalink/devices/test
{
  "protocol": "${device.protocol}",
  "connection_config": ${JSON.stringify(device.config, null, 2).replace(/\n/g, '\n  ')}
}`}
              </pre>
            </div>

            <div className="space-y-1.5">
              {(test.status ? Object.keys(test.stages || {}) : ['resolve', 'tcp', 'modbus_probe']).map((stage, i) => {
                const sr = (test.stages || {})[stage];
                return (
                  <div key={stage} className="flex items-center justify-between rounded-md border border-slate-800 bg-slate-900/40 px-3 py-2">
                    <div className="flex items-center gap-2.5">
                      <span className={clsx('grid place-items-center w-5 h-5 rounded-full text-[10px] font-bold',
                        sr ? 'bg-emerald-500/15 text-emerald-300 border border-emerald-500/30' : 'bg-slate-800 text-slate-500 border border-slate-700',
                      )}>
                        {sr ? <Icon name="check" className="w-3 h-3" /> : i + 1}
                      </span>
                      <span className="text-xs font-mono text-slate-300">{stage}</span>
                    </div>
                    {sr ? <span className="text-[11px] font-mono text-emerald-300">{sr.latency_ms}ms</span>
                      : test.status === 'running' ? <span className="text-[11px] text-slate-500">等候…</span>
                        : <span className="text-[11px] text-slate-600">—</span>}
                  </div>
                );
              })}
            </div>

            <div className="flex items-center gap-2">
              <Button onClick={runTest} disabled={testing} variant="secondary"
                icon={<Icon name={testing ? 'refresh' : 'play'} className={clsx('w-4 h-4', testing && 'animate-spin')} />}>
                {testing ? '測試中…' : '執行測試'}
              </Button>
              {test.status === 'success' && (
                <span className="text-xs text-emerald-300 flex items-center gap-1.5">
                  <Icon name="check" className="w-3.5 h-3.5" />
                  延遲 {test.latency_ms}ms · 已就緒
                </span>
              )}
            </div>

            {test.status === 'success' && (
              <div className="rounded-lg border border-emerald-500/20 bg-emerald-500/[0.06] p-3">
                <div className="flex items-center gap-2 text-xs text-emerald-200 font-medium mb-1.5">
                  <Icon name="check" className="w-4 h-4" /> Readiness 檢查通過
                </div>
                <ul className="text-[11px] text-emerald-100/80 space-y-0.5 pl-1">
                  <li>· Connection Configuration ✓</li>
                  <li>· Protocol Probe ✓ ({test.latency_ms}ms)</li>
                  <li>· 可進行下一步：設定接入規則</li>
                </ul>
              </div>
            )}
          </div>
        </SectionCard>
      </div>

      {/* 底部狀態 + 繼續 */}
      <div className="flex items-center justify-between rounded-xl border border-slate-700/60 bg-slate-900/40 px-4 py-3">
        <div className="flex items-center gap-3 text-xs">
          <span className="text-slate-400">{devices.length} 個設備，{testedCount} 已通過測試</span>
          {!allTested && <StatusChip tone="warning">尚有 {devices.length - testedCount} 個設備未通過測試</StatusChip>}
        </div>
        <Button onClick={onContinue} disabled={!canContinue} icon={<Icon name="arrow" className="w-4 h-4" />}>
          全部建立並繼續
        </Button>
      </div>
    </div>
  );
}

window.Step1Device = Step1Device;
