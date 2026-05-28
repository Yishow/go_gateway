// Step 2: 接入規則 (SourceRule)
// 支援多條規則 - 上方規則 Tab 列、中間規則編輯器 + 視覺化、下方合併點位表

function Step2Rule({ state, setState, onContinue }) {
  const { rules = [], devices = [] } = state;
  const [selectedRuleId, setSelectedRuleId] = React.useState(() => rules[0]?.id);

  // 裝置色盤 (與 Step 1 對齊)
  const deviceColors = ['blue', 'emerald', 'amber', 'fuchsia', 'cyan', 'rose'];
  const deviceColorFor = (did) => deviceColors[devices.findIndex((d) => d.id === did) % deviceColors.length];
  const deviceById = (did) => devices.find((d) => d.id === did);

  // 確保選中的規則仍存在
  React.useEffect(() => {
    if (!rules.find((r) => r.id === selectedRuleId)) {
      setSelectedRuleId(rules[0]?.id);
    }
  }, [rules, selectedRuleId]);

  const selectedRule = rules.find((r) => r.id === selectedRuleId) || rules[0];

  const updateRule = (patch) => {
    setState((s) => ({
      ...s,
      rules: s.rules.map((r) => (r.id === selectedRuleId ? { ...r, ...patch } : r)),
    }));
  };

  const addRule = () => {
    const newRule = makeDefaultRule({
      device_id: selectedRule?.device_id || devices[0]?.id || null,
      name: `規則 ${rules.length + 1}`,
      start_address: '40001',
      count: 4,
      naming_prefix: `BLOCK${rules.length + 1}_`,
      scale_multiplier: 1,
      scale_offset: 0,
    });
    setState((s) => ({ ...s, rules: [...s.rules, newRule] }));
    setSelectedRuleId(newRule.id);
  };

  const removeRule = (id) => {
    if (rules.length <= 1) return;
    setState((s) => ({ ...s, rules: s.rules.filter((r) => r.id !== id) }));
  };

  const renameRule = (id, name) => {
    setState((s) => ({ ...s, rules: s.rules.map((r) => (r.id === id ? { ...r, name } : r)) }));
  };

  // 衍生所有點位 (使用每條規則自身的 device_id)
  const allPoints = deriveAllPoints(rules, devices[0]?.id || 'dev-tmp');
  const totalEnabled = allPoints.filter((p) => !p.skipped).length;

  // 偵測位址衝突
  const addrUsage = {};
  allPoints.forEach((p) => {
    if (p.skipped) return;
    addrUsage[p.address] = (addrUsage[p.address] || []).concat(p.rule_id);
  });
  const conflictAddrs = new Set(Object.entries(addrUsage).filter(([, arr]) => arr.length > 1).map(([k]) => k));

  // Modbus Share 對外發布位址計算
  const shareBaseRegister = state.settings?.modbus_share?.base_register ?? 40001;
  const shareLayout = computeShareLayout(rules, shareBaseRegister);
  const shareEnabledGlobally = state.settings?.modbus_share?.enabled ?? false;

  // 取得單一點位的 share 位址
  const shareAddrFor = (point) => {
    const layout = shareLayout[point.rule_id];
    if (!layout) return null;
    // 找到該點位在它規則內、未跳過的 index
    const r = rules.find((x) => x.id === point.rule_id);
    if (!r) return null;
    const skippedSet = new Set(r.skipped_addresses || []);
    let idx = 0;
    const startN = parseInt(String(r.start_address).replace(/[^0-9]/g, ''), 10);
    const width = dataTypeWidth(r.data_type);
    for (let i = 0; i < r.count; i++) {
      const addrStr = formatAddr(startN + i * width);
      if (skippedSet.has(addrStr)) continue;
      if (addrStr === point.address) {
        return layout.start + idx * layout.stride;
      }
      idx++;
    }
    return null;
  };

  // 網格批次選擇狀態
  const [gridSelection, setGridSelection] = React.useState(new Set());
  const [lastClickedIdx, setLastClickedIdx] = React.useState(null);

  // 重置選擇當切換規則
  React.useEffect(() => {
    setGridSelection(new Set());
    setLastClickedIdx(null);
  }, [selectedRuleId]);

  // 同步 points 到 state
  React.useEffect(() => {
    setState((s) => ({ ...s, points: allPoints }));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [JSON.stringify(rules)]);

  // 當前規則的本地計算
  const ruleStartN = selectedRule ? parseInt(String(selectedRule.start_address).replace(/[^0-9]/g, ''), 10) || 40001 : 0;
  const ruleWidth = selectedRule ? dataTypeWidth(selectedRule.data_type) : 1;
  const ruleEndN = ruleStartN + ((selectedRule?.count || 1) - 1) * ruleWidth;
  const fn = selectedRule ? fnFromAddr(selectedRule.start_address) : 'holding_register';
  const fnLabel = {
    coil: { label: 'Coil (R/W)', code: 'FC 01/05', tone: 'bg-amber-500/10 text-amber-200 border-amber-500/30' },
    discrete_input: { label: 'Discrete Input (RO)', code: 'FC 02', tone: 'bg-slate-700/40 text-slate-200 border-slate-500/30' },
    input_register: { label: 'Input Register (RO)', code: 'FC 04', tone: 'bg-cyan-500/10 text-cyan-200 border-cyan-500/30' },
    holding_register: { label: 'Holding Register', code: 'FC 03/06/16', tone: 'bg-blue-500/10 text-blue-200 border-blue-500/30' },
  }[fn];

  // 點位 - 該規則的 (使用該規則自身的 device_id)
  const rulePoints = selectedRule ? derivePoints(selectedRule, selectedRule.device_id || devices[0]?.id || 'dev-tmp', new Set(selectedRule.skipped_addresses || [])) : [];
  const ruleEnabledCount = rulePoints.filter((p) => !p.skipped).length;

  const toggleSkip = (addr) => {
    const set = new Set(selectedRule.skipped_addresses || []);
    if (set.has(addr)) set.delete(addr);
    else set.add(addr);
    updateRule({ skipped_addresses: [...set] });
  };

  // 規則顏色 (圖示用) - 給每條規則不同顏色，方便在合併網格區分
  const ruleColors = ['blue', 'emerald', 'amber', 'fuchsia', 'cyan', 'rose'];
  const colorFor = (rid) => ruleColors[rules.findIndex((r) => r.id === rid) % ruleColors.length];
  const colorClasses = {
    blue: { bg: 'bg-blue-500/15', border: 'border-blue-500/40', text: 'text-blue-100', solid: 'bg-blue-500' },
    emerald: { bg: 'bg-emerald-500/15', border: 'border-emerald-500/40', text: 'text-emerald-100', solid: 'bg-emerald-500' },
    amber: { bg: 'bg-amber-500/15', border: 'border-amber-500/40', text: 'text-amber-100', solid: 'bg-amber-500' },
    fuchsia: { bg: 'bg-fuchsia-500/15', border: 'border-fuchsia-500/40', text: 'text-fuchsia-100', solid: 'bg-fuchsia-500' },
    cyan: { bg: 'bg-cyan-500/15', border: 'border-cyan-500/40', text: 'text-cyan-100', solid: 'bg-cyan-500' },
    rose: { bg: 'bg-rose-500/15', border: 'border-rose-500/40', text: 'text-rose-100', solid: 'bg-rose-500' },
  };

  return (
    <div className="space-y-5 sweep-in">
      {/* ── 規則 Tab 列 ─────────────────────────────────── */}
      <SectionCard
        title="接入規則"
        subtitle="可同時定義多條範圍規則，每條展開為一組點位"
        icon={<Icon name="rule" className="w-5 h-5" />}
        aside={
          <div className="flex items-center gap-2">
            <StatusChip tone="info">{rules.length} 條規則 · 共 {totalEnabled} 點位</StatusChip>
            {conflictAddrs.size > 0 && <StatusChip tone="error">{conflictAddrs.size} 位址衝突</StatusChip>}
          </div>
        }
        contentClassName="p-3"
      >
        <div className="flex items-stretch gap-2 overflow-x-auto pb-1">
          {rules.map((r) => {
            const cc = colorClasses[colorFor(r.id)];
            const isActive = r.id === selectedRuleId;
            const rPoints = derivePoints(r, r.device_id || devices[0]?.id || 'dev-tmp', new Set(r.skipped_addresses || []));
            const rEnabled = rPoints.filter((p) => !p.skipped).length;
            return (
              <div
                key={r.id}
                onClick={() => setSelectedRuleId(r.id)}
                className={clsx(
                  'group relative flex-shrink-0 cursor-pointer rounded-lg border px-3 py-2 transition-all min-w-[200px]',
                  isActive ? `${cc.border} ${cc.bg} shadow-md` : 'border-slate-700/60 bg-slate-900/40 hover:border-slate-600 hover:bg-slate-800/60',
                )}
              >
                <div className="flex items-center gap-2">
                  <span className={clsx('w-2 h-2 rounded-full flex-shrink-0', cc.solid)} />
                  <input
                    value={r.name}
                    onChange={(e) => renameRule(r.id, e.target.value)}
                    onClick={(e) => e.stopPropagation()}
                    className={clsx(
                      'bg-transparent border-none outline-none text-sm font-semibold w-full',
                      isActive ? 'text-slate-100' : 'text-slate-300',
                    )}
                  />
                  <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                    <Toggle
                      size="sm"
                      checked={r.enabled}
                      onChange={(v) => {
                        setState((s) => ({ ...s, rules: s.rules.map((x) => (x.id === r.id ? { ...x, enabled: v } : x)) }));
                      }}
                    />
                    {rules.length > 1 && (
                      <button
                        onClick={(e) => { e.stopPropagation(); removeRule(r.id); }}
                        className="text-slate-500 hover:text-red-400 p-0.5"
                        title="刪除規則"
                      >
                        <Icon name="close" className="w-3.5 h-3.5" />
                      </button>
                    )}
                  </div>
                </div>
                <div className="mt-1 flex items-center justify-between text-[10px] font-mono">
                  <span className="text-slate-500">{r.start_address} · {r.data_type}</span>
                  <span className={isActive ? cc.text : 'text-slate-400'}>{rEnabled}/{r.count} pts</span>
                </div>
                {/* 所屬裝置 */}
                {devices.length > 1 && (() => {
                  const dev = deviceById(r.device_id);
                  const dcc = colorClasses[deviceColorFor(r.device_id)];
                  return (
                    <div className="mt-1 flex items-center gap-1 text-[10px]">
                      <span className={clsx('w-1.5 h-1.5 rounded-full', dcc?.solid || 'bg-slate-500')} />
                      <span className="text-slate-400 truncate">{dev?.name || '未指派裝置'}</span>
                    </div>
                  );
                })()}
              </div>
            );
          })}
          <button
            onClick={addRule}
            className="flex-shrink-0 flex items-center gap-1.5 px-3 py-2 rounded-lg border border-dashed border-slate-700 text-slate-400 hover:border-blue-500/40 hover:text-blue-300 hover:bg-blue-500/[0.04] transition-all text-sm"
          >
            <Icon name="plus" className="w-4 h-4" />
            新增規則
          </button>
        </div>
      </SectionCard>

      {/* ── 編輯器 + 視覺化 ────────────────────────────────── */}
      {selectedRule && (() => {
        const cc = colorClasses[colorFor(selectedRule.id)];
        return (
          <div className="grid grid-cols-12 gap-5">
            {/* 左：規則表單 */}
            <SectionCard
              className="col-span-12 lg:col-span-5"
              title={`編輯：${selectedRule.name}`}
              subtitle="位址範圍、命名、線性轉換"
              icon={<Icon name="rule" className="w-5 h-5" />}
              aside={<StatusChip tone="info">{ruleEnabledCount} / {selectedRule.count} 點位</StatusChip>}
            >
              <div className="space-y-4">
                {/* 所屬裝置 */}
                <Field label="所屬裝置" required hint={devices.length === 1 ? '只有一個設備，自動指派' : '此規則展開的點位將歸屬於這個設備'}>
                  <Select
                    value={selectedRule.device_id || ''}
                    onChange={(e) => updateRule({ device_id: e.target.value })}
                    disabled={devices.length <= 1}
                  >
                    <option value="">— 選擇設備 —</option>
                    {devices.map((d) => (
                      <option key={d.id} value={d.id}>
                        {d.name} ({d.protocol})
                      </option>
                    ))}
                  </Select>
                </Field>

                <div className="grid grid-cols-2 gap-3">
                  <Field label="起始位址" required hint="Modbus 4xxxx = Holding Register">
                    <Input
                      value={selectedRule.start_address}
                      onChange={(e) => updateRule({ start_address: e.target.value, skipped_addresses: [] })}
                      className="font-mono"
                    />
                  </Field>
                  <Field label="點位數量" required>
                    <Input
                      type="number" min="1" max="64"
                      value={selectedRule.count}
                      onChange={(e) => updateRule({ count: Math.max(1, Math.min(64, +e.target.value || 1)), skipped_addresses: [] })}
                      className="font-mono"
                    />
                  </Field>
                  <Field label="資料型態" required hint={`寬度 ${ruleWidth} register`}>
                    <Select
                      value={selectedRule.data_type}
                      onChange={(e) => updateRule({ data_type: e.target.value, skipped_addresses: [] })}
                    >
                      {DATA_TYPES.map((t) => <option key={t} value={t}>{t}</option>)}
                    </Select>
                  </Field>
                  <Field label="命名前綴" hint={`生成 ${selectedRule.naming_prefix}01, …`}>
                    <Input
                      value={selectedRule.naming_prefix}
                      onChange={(e) => updateRule({ naming_prefix: e.target.value })}
                      className="font-mono"
                    />
                  </Field>
                </div>

                <div className="rounded-xl border border-slate-700/60 bg-slate-950/40 p-4">
                  <div className="text-[11px] uppercase tracking-wider text-slate-500 mb-2">範圍摘要</div>
                  <div className="flex items-center gap-3 font-mono text-sm">
                    <span className="text-slate-200">{formatAddr(ruleStartN)}</span>
                    <span className="text-slate-600">→</span>
                    <span className="text-slate-200">{formatAddr(ruleEndN + ruleWidth - 1)}</span>
                    <span className="ml-auto text-[11px] text-slate-500">共 {selectedRule.count * ruleWidth} register</span>
                  </div>
                  <div className="mt-3 flex items-center gap-2">
                    <span className={clsx('chip', fnLabel.tone)}>{fnLabel.label}</span>
                    <span className="text-[11px] font-mono text-slate-500">{fnLabel.code}</span>
                  </div>
                </div>

                <details className="group rounded-xl border border-slate-700/60 bg-slate-950/40 p-4 open:bg-slate-900/40" open>
                  <summary className="flex items-center justify-between cursor-pointer list-none">
                    <span className="text-xs font-semibold uppercase tracking-wider text-slate-400">線性轉換 / 字節序 (規則預設值)</span>
                    <Icon name="chevron" className="w-4 h-4 text-slate-500 transition-transform group-open:rotate-90" />
                  </summary>
                  <p className="text-[10px] text-slate-500 mt-2 mb-3">
                    <Icon name="info" className="inline w-3 h-3 mr-0.5" />
                    這裡設的值會作為下一步「點位映射」的初始值，仍可在映射表中個別調整。
                  </p>
                  <div className="grid grid-cols-2 gap-3">
                    <Field label="Multiplier" hint="raw × multiplier + offset">
                      <Input type="number" step="0.01" value={selectedRule.scale_multiplier}
                        onChange={(e) => updateRule({ scale_multiplier: +e.target.value })} className="font-mono" />
                    </Field>
                    <Field label="Offset">
                      <Input type="number" step="0.01" value={selectedRule.scale_offset}
                        onChange={(e) => updateRule({ scale_offset: +e.target.value })} className="font-mono" />
                    </Field>
                    <Field label="Byte order" className="col-span-2">
                      <Select value={selectedRule.data_format || ''}
                        onChange={(e) => updateRule({ data_format: e.target.value })}>
                        <option value="">(預設 ABCD - Big Endian)</option>
                        <option value="ABCD">ABCD - Big Endian</option>
                        <option value="BADC">BADC - Big Endian byte-swap</option>
                        <option value="CDAB">CDAB - Little Endian word-swap</option>
                        <option value="DCBA">DCBA - Little Endian</option>
                      </Select>
                    </Field>
                  </div>
                </details>

                {/* Local Modbus Share 對外發布設定 */}
                <details className="group rounded-xl border border-slate-700/60 bg-slate-950/40 p-4 open:bg-slate-900/40">
                  <summary className="flex items-center justify-between cursor-pointer list-none">
                    <span className="text-xs font-semibold uppercase tracking-wider text-slate-400 flex items-center gap-2">
                      Local Modbus Share 位址
                      {!shareEnabledGlobally && <span className="chip border-slate-600/40 bg-slate-700/20 text-slate-400">全域未啟用</span>}
                      {shareEnabledGlobally && selectedRule.share_enabled && shareLayout[selectedRule.id] && (
                        <span className={clsx('chip border-emerald-500/30 bg-emerald-500/10 text-emerald-200')}>
                          {formatAddr(shareLayout[selectedRule.id].start)} … {formatAddr(shareLayout[selectedRule.id].end - 1)}
                        </span>
                      )}
                    </span>
                    <Icon name="chevron" className="w-4 h-4 text-slate-500 transition-transform group-open:rotate-90" />
                  </summary>
                  <p className="text-[10px] text-slate-500 mt-2 mb-3">
                    <Icon name="info" className="inline w-3 h-3 mr-0.5" />
                    這條規則展開的所有點位會依序對應到下列 register。可手動指定起點或留空讓系統自動接續配置。設定→Local Modbus Share 必須開啟。
                  </p>
                  <div className="grid grid-cols-2 gap-3">
                    <Field label="Share 起始 Register" hint={shareLayout[selectedRule.id]?.auto ? '自動配置' : '手動指定'}>
                      <Input
                        type="text"
                        placeholder="(自動)"
                        value={selectedRule.share_start_register ?? ''}
                        onChange={(e) => {
                          const v = e.target.value;
                          updateRule({ share_start_register: v === '' ? null : +v });
                        }}
                        className="font-mono"
                      />
                    </Field>
                    <Field label="Stride (跨距)" hint="留空 = 依資料型態寬度">
                      <Input
                        type="text"
                        placeholder={`(${dataTypeWidth(selectedRule.data_type)})`}
                        value={selectedRule.share_stride ?? ''}
                        onChange={(e) => {
                          const v = e.target.value;
                          updateRule({ share_stride: v === '' ? null : +v });
                        }}
                        className="font-mono"
                      />
                    </Field>
                  </div>
                  <div className="mt-3 flex items-center justify-between">
                    <Toggle
                      checked={selectedRule.share_enabled ?? true}
                      onChange={(v) => updateRule({ share_enabled: v })}
                      label="透過 Modbus Share 對外發布此規則的點位"
                    />
                    {selectedRule.share_start_register != null && (
                      <button
                        onClick={() => updateRule({ share_start_register: null, share_stride: null })}
                        className="text-[11px] text-slate-400 hover:text-blue-300 flex items-center gap-1"
                      >
                        <Icon name="refresh" className="w-3 h-3" /> 改回自動分配
                      </button>
                    )}
                  </div>
                </details>

                <div className="flex items-center gap-2 pt-1">
                  <Toggle checked={selectedRule.enabled} onChange={(v) => updateRule({ enabled: v })}
                    label="啟用此規則 (建立後立即輪詢)" />
                </div>
              </div>
            </SectionCard>

            {/* 右：當前規則的點位網格 - 多功能 */}
            <SectionCard
              className="col-span-12 lg:col-span-7"
              title="當前規則點位網格"
              subtitle="點擊切換略過 · Shift+點擊：範圍切換 · Ctrl/⌘+點擊：加入選取"
              icon={<Icon name="grid" className="w-5 h-5" />}
              aside={
                <div className="flex items-center gap-1.5">
                  <span className={clsx('chip', cc.border, cc.bg, cc.text)}>啟用 {ruleEnabledCount}</span>
                  {(selectedRule.skipped_addresses || []).length > 0 && (
                    <span className="chip border-amber-500/30 bg-amber-500/10 text-amber-200">略過 {(selectedRule.skipped_addresses || []).length}</span>
                  )}
                  {gridSelection.size > 0 && (
                    <span className="chip border-blue-500/40 bg-blue-500/15 text-blue-200">已選 {gridSelection.size}</span>
                  )}
                </div>
              }
            >
              {/* 批次操作工具列 */}
              <div className="mb-3 flex flex-wrap items-center gap-1.5 rounded-lg border border-slate-700/50 bg-slate-950/40 px-2 py-1.5">
                <Button
                  size="sm" variant="ghost"
                  onClick={() => {
                    // 全部啟用
                    updateRule({ skipped_addresses: [] });
                    setGridSelection(new Set());
                  }}
                  icon={<Icon name="check" className="w-3 h-3" />}
                >全部啟用</Button>
                <Button
                  size="sm" variant="ghost"
                  onClick={() => {
                    const startN = parseInt(String(selectedRule.start_address).replace(/[^0-9]/g, ''), 10);
                    const w = dataTypeWidth(selectedRule.data_type);
                    const all = [];
                    for (let i = 0; i < selectedRule.count; i++) all.push(formatAddr(startN + i * w));
                    updateRule({ skipped_addresses: all });
                  }}
                  icon={<Icon name="close" className="w-3 h-3" />}
                >全部略過</Button>
                <Button
                  size="sm" variant="ghost"
                  onClick={() => {
                    const startN = parseInt(String(selectedRule.start_address).replace(/[^0-9]/g, ''), 10);
                    const w = dataTypeWidth(selectedRule.data_type);
                    const skippedNow = new Set(selectedRule.skipped_addresses || []);
                    const inverted = [];
                    for (let i = 0; i < selectedRule.count; i++) {
                      const a = formatAddr(startN + i * w);
                      if (!skippedNow.has(a)) inverted.push(a);
                    }
                    updateRule({ skipped_addresses: inverted });
                  }}
                  icon={<Icon name="refresh" className="w-3 h-3" />}
                >反轉啟用</Button>
                <div className="h-4 w-px bg-slate-700 mx-1" />
                <Button
                  size="sm" variant="ghost"
                  disabled={gridSelection.size === 0}
                  onClick={() => {
                    const next = new Set(selectedRule.skipped_addresses || []);
                    gridSelection.forEach((a) => next.add(a));
                    updateRule({ skipped_addresses: [...next] });
                    setGridSelection(new Set());
                  }}
                >略過選取 ({gridSelection.size})</Button>
                <Button
                  size="sm" variant="ghost"
                  disabled={gridSelection.size === 0}
                  onClick={() => {
                    const next = new Set(selectedRule.skipped_addresses || []);
                    gridSelection.forEach((a) => next.delete(a));
                    updateRule({ skipped_addresses: [...next] });
                    setGridSelection(new Set());
                  }}
                >啟用選取</Button>
                {gridSelection.size > 0 && (
                  <Button size="sm" variant="ghost" onClick={() => setGridSelection(new Set())}>
                    清除選取
                  </Button>
                )}
                <div className="ml-auto text-[10px] font-mono text-slate-500 pr-1">
                  total {selectedRule.count} · stride {ruleWidth}
                </div>
              </div>

              <div className="rounded-xl border border-slate-700/60 bg-slate-950/60 p-4 bg-grid">
                <div className="mb-3 flex items-center justify-between">
                  <div className="font-mono text-[11px] text-slate-500">
                    {formatAddr(ruleStartN)} … {formatAddr(ruleEndN + ruleWidth - 1)}
                  </div>
                  <div className="font-mono text-[11px] text-slate-500">
                    <span className={cc.text}>{ruleEnabledCount}</span> active · <span className="text-amber-300">{(selectedRule.skipped_addresses || []).length}</span> skipped
                    {shareEnabledGlobally && selectedRule.share_enabled && shareLayout[selectedRule.id] && (
                      <> · <span className="text-emerald-300">share@{shareLayout[selectedRule.id].start}</span></>
                    )}
                  </div>
                </div>
                <div className="grid grid-cols-4 sm:grid-cols-6 md:grid-cols-8 gap-1.5">
                  {rulePoints.map((p, i) => {
                    const isSelected = gridSelection.has(p.address);
                    const isConflict = conflictAddrs.has(p.address);
                    const shareAddr = shareEnabledGlobally && selectedRule.share_enabled && shareLayout[selectedRule.id]
                      ? shareAddrFor(p) : null;
                    const onClick = (e) => {
                      // Shift+click：對 lastClickedIdx → i 的範圍批次切換 skip
                      if (e.shiftKey && lastClickedIdx != null) {
                        const [lo, hi] = [Math.min(lastClickedIdx, i), Math.max(lastClickedIdx, i)];
                        const set = new Set(selectedRule.skipped_addresses || []);
                        const isCurrentSkipped = (selectedRule.skipped_addresses || []).includes(p.address);
                        for (let j = lo; j <= hi; j++) {
                          const addr = rulePoints[j].address;
                          if (isCurrentSkipped) set.delete(addr);
                          else set.add(addr);
                        }
                        updateRule({ skipped_addresses: [...set] });
                      } else if (e.ctrlKey || e.metaKey) {
                        // 加入/移出選取
                        const next = new Set(gridSelection);
                        if (next.has(p.address)) next.delete(p.address);
                        else next.add(p.address);
                        setGridSelection(next);
                      } else {
                        // 一般點擊：toggle skip
                        toggleSkip(p.address);
                      }
                      setLastClickedIdx(i);
                    };
                    return (
                      <button
                        key={p.id}
                        onClick={onClick}
                        className={clsx(
                          'mem-cell border relative',
                          p.skipped
                            ? 'bg-amber-500/10 border-amber-500/30 text-amber-200'
                            : `${cc.bg} ${cc.border} ${cc.text} hover:brightness-125`,
                          isConflict && 'ring-2 ring-red-500/50 ring-offset-1 ring-offset-slate-950',
                          isSelected && 'ring-2 ring-blue-400 ring-offset-1 ring-offset-slate-950',
                        )}
                        title={[
                          `${p.name} @ ${p.address}`,
                          shareAddr != null ? `Share → ${formatAddr(shareAddr)}` : null,
                          isConflict ? '⚠ 與其他規則衝突' : null,
                        ].filter(Boolean).join('\n')}
                      >
                        <span className="addr">{p.address}</span>
                        <span className="meta">{p.name}</span>
                        {shareAddr != null && !p.skipped && (
                          <span className="absolute -bottom-1 right-0.5 text-[8px] font-mono text-emerald-300/90 bg-slate-950/80 px-1 rounded">
                            →{shareAddr}
                          </span>
                        )}
                        {isConflict && (
                          <span className="absolute -top-1 -right-1 grid place-items-center w-3.5 h-3.5 rounded-full bg-red-500 text-white text-[8px] font-bold">!</span>
                        )}
                      </button>
                    );
                  })}
                </div>
                <div className="mt-3 flex flex-wrap gap-x-4 gap-y-1 text-[10px] font-mono text-slate-500">
                  <span>· width {ruleWidth} | start {formatAddr(ruleStartN)}</span>
                  <span>· {selectedRule.naming_prefix}NN ({selectedRule.naming_prefix}01–{selectedRule.naming_prefix}{String(selectedRule.count).padStart(2, '0')})</span>
                  {conflictAddrs.size > 0 && (
                    <span className="text-red-300">· {[...conflictAddrs].filter((a) => rulePoints.some((p) => p.address === a)).length} 個位址與其他規則衝突</span>
                  )}
                </div>
              </div>
            </SectionCard>
          </div>
        );
      })()}

      {/* ── 合併點位表 ──────────────────────────────────── */}
      <SectionCard
        title={`將要建立的點位 (${totalEnabled})`}
        subtitle="所有啟用規則合併後的結果"
        icon={<Icon name="table" className="w-5 h-5" />}
        aside={
          <div className="flex items-center gap-2 text-[11px] text-slate-500">
            {shareEnabledGlobally && <StatusChip tone="info" dot={false}>Modbus Share 開啟</StatusChip>}
            <span>POST /source-rules × {rules.filter((r) => r.enabled).length} · POST /points × {totalEnabled}</span>
          </div>
        }
        contentClassName="p-0"
      >
        <div className="max-h-[280px] overflow-y-auto">
          <table className="w-full text-xs">
            <thead className="bg-slate-900/70 text-slate-500 sticky top-0">
              <tr>
                <th className="px-3 py-2 text-left font-medium">所屬設備</th>
                <th className="px-3 py-2 text-left font-medium">所屬規則</th>
                <th className="px-3 py-2 text-left font-medium">點位名稱</th>
                <th className="px-3 py-2 text-left font-medium">位址</th>
                <th className="px-3 py-2 text-left font-medium">型態</th>
                <th className="px-3 py-2 text-left font-medium">功能碼</th>
                <th className="px-3 py-2 text-left font-medium">
                  Share 位址
                  {!shareEnabledGlobally && <span className="ml-1 text-slate-600">(off)</span>}
                </th>
                <th className="px-3 py-2 text-right font-medium">狀態</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/60 font-mono">
              {allPoints.map((p) => {
                const cc = colorClasses[colorFor(p.rule_id)];
                const dcc = colorClasses[deviceColorFor(p.device_id)];
                const dev = deviceById(p.device_id);
                const conflict = conflictAddrs.has(p.address);
                const r = rules.find((x) => x.id === p.rule_id);
                const shareAddr = !p.skipped && r?.share_enabled && shareEnabledGlobally ? shareAddrFor(p) : null;
                return (
                  <tr key={p.id} className={clsx(p.skipped && 'opacity-40')}>
                    <td className="px-3 py-1.5">
                      <span className="inline-flex items-center gap-1.5">
                        <span className={clsx('w-1.5 h-1.5 rounded-full', dcc?.solid || 'bg-slate-500')} />
                        <span className="text-slate-300 text-[11px]">{dev?.name || '—'}</span>
                      </span>
                    </td>
                    <td className="px-3 py-1.5">
                      <span className="inline-flex items-center gap-1.5">
                        <span className={clsx('w-1.5 h-1.5 rounded-full', cc.solid)} />
                        <span className="text-slate-300 text-[11px]">{p.rule_name}</span>
                      </span>
                    </td>
                    <td className="px-3 py-1.5 text-slate-200">{p.name}</td>
                    <td className={clsx('px-3 py-1.5', conflict ? 'text-red-300' : 'text-blue-200')}>
                      {p.address}{conflict && <span className="ml-1 text-red-400">⚠</span>}
                    </td>
                    <td className="px-3 py-1.5 text-slate-400">{p.data_type}</td>
                    <td className="px-3 py-1.5 text-slate-500">{p.function}</td>
                    <td className="px-3 py-1.5">
                      {shareAddr != null ? (
                        <span className="text-emerald-300">{formatAddr(shareAddr)}</span>
                      ) : (
                        <span className="text-slate-600">—</span>
                      )}
                    </td>
                    <td className="px-3 py-1.5 text-right">
                      {p.skipped ? (
                        <StatusChip tone="warning" dot={false}>跳過</StatusChip>
                      ) : conflict ? (
                        <StatusChip tone="error" dot={false}>衝突</StatusChip>
                      ) : (
                        <StatusChip tone="success" dot={false}>建立</StatusChip>
                      )}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
        <div className="border-t border-slate-800 px-4 py-2.5 flex items-center justify-between">
          <div className="text-[11px] text-slate-500">
            <Icon name="info" className="inline w-3.5 h-3.5 mr-1" />
            指派到輪詢群組 <span className="font-mono text-slate-300">{DEFAULTS.pollingGroup.name}</span>
            {!shareEnabledGlobally && (
              <span className="ml-3">· 啟用 Share 請至 <span className="text-blue-300">設定 → Local Modbus Share</span></span>
            )}
          </div>
          <Button onClick={onContinue} disabled={totalEnabled === 0 || conflictAddrs.size > 0} icon={<Icon name="arrow" className="w-4 h-4" />}>
            繼續到映射
          </Button>
        </div>
      </SectionCard>
    </div>
  );
}

window.Step2Rule = Step2Rule;
