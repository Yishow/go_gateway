// Step 4: 儲存到資料庫 - DatabaseTargetMapping
// 連接器設定 + 表/欄位映射 + 提交 (Commit)

function Step4Database({ state, setState, onCommit }) {
  const { points = [], mappings = {}, db = {} } = state;
  const enabledPoints = points.filter((p) => !p.skipped);

  // 預設 connector
  const connector = state.db?.connector || {
    name: 'TimeSeries Prod',
    kind: 'postgres',
    host: 'tsdb.internal',
    port: 5432,
    database: 'gateway_metrics',
    username: 'gw_writer',
    status: 'ready',
    schema: 'public',
    table: 'sensor_readings',
    write_mode: 'insert',
    write_interval_seconds: 5,
    timestamp_column: 'ts',
  };

  const updateConn = (patch) =>
    setState((s) => ({ ...s, db: { ...s.db, connector: { ...connector, ...patch } } }));

  // 表欄位映射
  const tableInfo = (SAMPLE_DB_TABLES[connector.kind] || SAMPLE_DB_TABLES.postgres)[0];
  const columns = tableInfo?.columns || [];
  const columnNames = columns.filter((c) => !c.primary_key).map((c) => c.name);

  // 初始化每個 mapping 的 column 對應
  React.useEffect(() => {
    if (!enabledPoints.length) return;
    setState((s) => {
      const existing = s.db?.targets || {};
      const next = {};
      const used = new Set();
      enabledPoints.forEach((p, i) => {
        const m = s.mappings?.[p.id];
        const tagShort = (m?.tag_key || '').split('.').pop() || `col_${i + 1}`;
        // 嘗試精確匹配尾段且未被使用；找不到則照索引指派以避免衝突
        const exactMatch = columnNames.find(
          (c) => !used.has(c) && (c === tagShort || c.endsWith('_' + tagShort) || c.startsWith(tagShort + '_')),
        );
        const indexFallback = columnNames[i % columnNames.length];
        const auto = exactMatch || (used.has(indexFallback) ? columnNames.find((c) => !used.has(c)) : indexFallback) || indexFallback;
        used.add(auto);
        next[p.id] = existing[p.id] || {
          tag_id: `tag.${m?.tag_key}`,
          column_name: auto,
          enabled: true,
        };
      });
      return { ...s, db: { ...s.db, connector, targets: next } };
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [connector.kind, enabledPoints.length]);

  const targets = state.db?.targets || {};
  const updateTarget = (pid, patch) =>
    setState((s) => ({ ...s, db: { ...s.db, targets: { ...s.db.targets, [pid]: { ...s.db.targets[pid], ...patch } } } }));

  // 偵測衝突 - 兩個 mapping 寫入同一欄位
  const colUsage = {};
  Object.entries(targets).forEach(([pid, t]) => {
    if (!t?.enabled || !t.column_name) return;
    colUsage[t.column_name] = (colUsage[t.column_name] || []).concat(pid);
  });

  const dbKind = DB_KINDS.find((k) => k.id === connector.kind);

  // Commit 狀態
  const [committing, setCommitting] = React.useState(false);
  const [committed, setCommitted] = React.useState(state.committed || false);

  const startCommit = () => {
    setCommitting(true);
    setState((s) => ({ ...s, commit: { status: 'running', logs: [], started_at: new Date() } }));
    const logs = [
      { t: 0, label: `POST /devices × ${(state.devices || []).length}`, detail: (state.devices || []).map((d) => `${d.name} (${d.protocol})`).join(' · '), status: 'ok' },
      { t: 1, label: `POST /devices/:id/activate × ${(state.devices || []).length}`, detail: 'draft → active', status: 'ok' },
      { t: 2, label: `POST /source-rules × ${(state.rules || []).filter((r) => r.enabled).length}`, detail: (state.rules || []).map((r) => `${r.name}: ${r.start_address}×${r.count}`).join(' · ') || '—', status: 'ok' },
      { t: 3, label: `POST /points × ${enabledPoints.length}`, detail: 'bulk create', status: 'ok' },
      { t: 4, label: 'POST /polling-groups', detail: '快速輪詢 1s, enabled', status: 'ok' },
      { t: 5, label: `POST /tags × ${enabledPoints.length}`, detail: 'register tag keys', status: 'ok' },
      { t: 6, label: `POST /mappings × ${enabledPoints.length}`, detail: 'point ↔ tag, scale pipeline', status: 'ok' },
      { t: 7, label: 'POST /db-connectors/:id/test', detail: `${connector.kind} ${connector.host}:${connector.port}`, status: 'ok' },
      { t: 8, label: `POST /db-targets × ${Object.values(targets).filter((x) => x.enabled).length}`, detail: `→ ${connector.schema}.${connector.table}`, status: 'ok' },
      { t: 9, label: 'POST /scheduler/start', detail: 'collectors started', status: 'ok' },
    ];
    let i = 0;
    const tick = () => {
      if (i >= logs.length) {
        setCommitting(false);
        setCommitted(true);
        setState((s) => ({ ...s, committed: true, commit: { status: 'success', logs, finished_at: new Date() } }));
        return;
      }
      setState((s) => ({ ...s, commit: { ...s.commit, logs: logs.slice(0, i + 1) } }));
      i++;
      setTimeout(tick, 280);
    };
    setTimeout(tick, 200);
  };

  const enabledTargetCount = Object.values(targets).filter((t) => t.enabled).length;
  const hasConflict = Object.values(colUsage).some((arr) => arr.length > 1);

  return (
    <div className="grid grid-cols-12 gap-5 sweep-in">
      {/* 上：Connector */}
      <SectionCard
        className="col-span-12"
        title="資料庫連接器 (Connector)"
        subtitle="寫入時序值的目標資料庫"
        icon={<Icon name="db" className="w-5 h-5" />}
        aside={<StatusChip tone="success">已連線 · {connector.host}</StatusChip>}
      >
        <div className="grid grid-cols-12 gap-4 items-end">
          <div className="col-span-12 md:col-span-3">
            <span className="label">資料庫類型</span>
            <div className="grid grid-cols-4 gap-1.5">
              {DB_KINDS.map((k) => {
                const active = connector.kind === k.id;
                return (
                  <button
                    key={k.id}
                    onClick={() => updateConn({ kind: k.id })}
                    className={clsx(
                      'flex flex-col items-center justify-center rounded-lg border px-2 py-2.5 transition-all',
                      active ? 'border-blue-500/60 bg-blue-500/10' : 'border-slate-700/60 bg-slate-900/40 hover:border-slate-600',
                    )}
                  >
                    <span className="text-base leading-none">{k.icon}</span>
                    <span className={clsx('text-[10px] mt-1 font-medium', active ? 'text-blue-200' : 'text-slate-400')}>{k.name}</span>
                  </button>
                );
              })}
            </div>
          </div>
          <Field label="連線名稱" className="col-span-6 md:col-span-2">
            <Input value={connector.name} onChange={(e) => updateConn({ name: e.target.value })} />
          </Field>
          <Field label="Host" className="col-span-6 md:col-span-2">
            <Input value={connector.host} onChange={(e) => updateConn({ host: e.target.value })} className="font-mono" />
          </Field>
          <Field label="Port" className="col-span-3 md:col-span-1">
            <Input type="number" value={connector.port} onChange={(e) => updateConn({ port: +e.target.value })} className="font-mono" />
          </Field>
          <Field label="Database" className="col-span-9 md:col-span-2">
            <Input value={connector.database} onChange={(e) => updateConn({ database: e.target.value })} className="font-mono" />
          </Field>
          <Field label="Schema" className="col-span-6 md:col-span-1">
            <Input value={connector.schema} onChange={(e) => updateConn({ schema: e.target.value })} className="font-mono" />
          </Field>
          <Field label="Table" className="col-span-6 md:col-span-1">
            <Input value={connector.table} onChange={(e) => updateConn({ table: e.target.value })} className="font-mono" />
          </Field>
        </div>

        <div className="mt-4 flex flex-wrap items-center gap-3 border-t border-slate-800 pt-3 text-xs">
          <span className="text-slate-500">寫入策略：</span>
          <label className="inline-flex items-center gap-2 cursor-pointer">
            <input type="radio" checked={connector.write_mode === 'insert'} onChange={() => updateConn({ write_mode: 'insert' })} className="accent-blue-500" />
            <span className="text-slate-300">INSERT (時序追加)</span>
          </label>
          <label className="inline-flex items-center gap-2 cursor-pointer">
            <input type="radio" checked={connector.write_mode === 'upsert'} onChange={() => updateConn({ write_mode: 'upsert' })} className="accent-blue-500" />
            <span className="text-slate-300">UPSERT (依時間戳)</span>
          </label>
          <span className="ml-auto text-slate-500">寫入間隔</span>
          <input
            type="number"
            value={connector.write_interval_seconds}
            onChange={(e) => updateConn({ write_interval_seconds: +e.target.value })}
            className="w-16 bg-slate-950/60 border border-slate-700/60 rounded px-2 py-1 font-mono text-slate-200 text-center"
          />
          <span className="text-slate-500">秒</span>
        </div>
      </SectionCard>

      {/* 下左：欄位映射 */}
      <SectionCard
        className="col-span-12 xl:col-span-7"
        title={`Tag → 資料表欄位`}
        subtitle={`寫入 ${connector.schema}.${connector.table}`}
        icon={<Icon name="table" className="w-5 h-5" />}
        aside={hasConflict ? <StatusChip tone="error">欄位衝突</StatusChip> : <StatusChip tone="success">{enabledTargetCount} 對應</StatusChip>}
        contentClassName="p-0"
      >
        <div className="overflow-x-auto">
          <table className="w-full text-xs">
            <thead className="bg-slate-900/70 text-slate-500 border-b border-slate-800">
              <tr>
                <th className="px-3 py-2 text-left font-medium">Tag</th>
                <th className="px-3 py-2 text-left font-medium">點位</th>
                <th className="px-3 py-2 text-center font-medium w-8"></th>
                <th className="px-3 py-2 text-left font-medium">資料表欄位</th>
                <th className="px-3 py-2 text-left font-medium">欄位型態</th>
                <th className="px-3 py-2 text-right font-medium">啟用</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/60">
              {enabledPoints.map((p) => {
                const m = mappings[p.id];
                const t = targets[p.id] || {};
                const colInfo = columns.find((c) => c.name === t.column_name);
                const conflict = colUsage[t.column_name]?.length > 1;
                if (!m) return null;
                return (
                  <tr key={p.id} className="hover:bg-slate-800/30">
                    <td className="px-3 py-1.5">
                      <div className="font-mono text-blue-200 text-[11px]">{m.tag_key}</div>
                      <div className="text-[10px] text-slate-500">{m.display_name}</div>
                    </td>
                    <td className="px-3 py-1.5 font-mono text-slate-400 text-[11px]">{p.name}<br /><span className="text-slate-600">@{p.address}</span></td>
                    <td className="px-3 py-1.5 text-center text-slate-600">→</td>
                    <td className="px-3 py-1.5">
                      <div className="flex items-center gap-1.5">
                        <select
                          value={t.column_name || ''}
                          onChange={(e) => updateTarget(p.id, { column_name: e.target.value })}
                          className={clsx(
                            'bg-slate-950/60 border rounded px-2 py-1 text-[11px] font-mono focus:outline-none focus:border-blue-500/60',
                            conflict ? 'border-red-500/50 text-red-200' : 'border-slate-700/50 text-slate-200',
                          )}
                        >
                          {columnNames.map((c) => <option key={c} value={c}>{c}</option>)}
                        </select>
                        {conflict && <Icon name="alert" className="w-3.5 h-3.5 text-red-400" />}
                      </div>
                    </td>
                    <td className="px-3 py-1.5 font-mono text-[11px] text-slate-500">
                      {colInfo?.data_type || '—'}
                    </td>
                    <td className="px-3 py-1.5 text-right">
                      <Toggle checked={t.enabled ?? true} onChange={(v) => updateTarget(p.id, { enabled: v })} size="sm" />
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
        {hasConflict && (
          <div className="border-t border-red-500/30 bg-red-500/[0.06] px-4 py-2 text-[11px] text-red-200 flex items-start gap-1.5">
            <Icon name="alert" className="w-3.5 h-3.5 mt-0.5 flex-shrink-0" />
            偵測到多個 Tag 寫入同一資料表欄位。請調整以避免覆寫。
          </div>
        )}
      </SectionCard>

      {/* 下右：提交摘要 + Commit */}
      <SectionCard
        className="col-span-12 xl:col-span-5"
        title="提交摘要 (Commit)"
        subtitle="一次原子寫入所有設定"
        icon={<Icon name="save" className="w-5 h-5" />}
      >
        {!committed && !committing && (
          <>
            <ul className="space-y-2 text-xs">
              {[
                { label: '新增設備', value: `${(state.devices || []).length} 個`, sub: (state.devices || []).map((d) => `${d.name} (${d.protocol})`).join(' · ') },
                { label: '接入規則', value: `${(state.rules || []).filter((r) => r.enabled).length} 條`, sub: (state.rules || []).map((r) => `${r.name} (${r.count})`).join(' · ') },
                { label: '點位', value: `${enabledPoints.length} 個`, sub: '加入快速輪詢 (1s)' },
                { label: '標籤映射', value: `${Object.values(mappings).filter((m) => m.enabled).length} 個 Tag`, sub: '含 scale/offset/cast 管線' },
                { label: '資料庫寫入', value: `${enabledTargetCount} 欄位`, sub: `${dbKind?.name} → ${connector.schema}.${connector.table}` },
              ].map((row) => (
                <li key={row.label} className="flex items-start justify-between gap-3 rounded-lg border border-slate-800 bg-slate-950/30 px-3 py-2">
                  <div>
                    <div className="text-[10px] uppercase tracking-wider text-slate-500">{row.label}</div>
                    <div className="font-mono text-slate-100 text-sm mt-0.5">{row.value}</div>
                  </div>
                  <div className="text-[10px] text-slate-500 text-right max-w-[55%]">{row.sub}</div>
                </li>
              ))}
            </ul>
            <Button
              onClick={startCommit}
              disabled={hasConflict || enabledTargetCount === 0}
              variant="success"
              size="lg"
              className="w-full mt-4"
              icon={<Icon name="save" className="w-4 h-4" />}
            >
              提交並啟動排程器
            </Button>
            <p className="mt-2 text-[11px] text-slate-500 text-center">
              <Icon name="info" className="inline w-3 h-3 mr-1" />
              將執行 10 個 API 呼叫並啟動 collector。
            </p>
          </>
        )}

        {(committing || committed) && (
          <div>
            <div className="space-y-1.5 font-mono text-[11px]">
              {(state.commit?.logs || []).map((log, i) => (
                <div key={i} className="flex items-start gap-2 sweep-in">
                  <Icon name="check" className="w-3.5 h-3.5 mt-0.5 text-emerald-400 flex-shrink-0" />
                  <div className="flex-1 min-w-0">
                    <div className="text-slate-200">{log.label}</div>
                    <div className="text-slate-500 text-[10px]">{log.detail}</div>
                  </div>
                  <span className="text-emerald-400/80 text-[10px]">200</span>
                </div>
              ))}
              {committing && (
                <div className="flex items-center gap-2 text-slate-500">
                  <span className="w-1.5 h-1.5 rounded-full bg-blue-400 pulse-dot" />
                  正在執行下一個指令…
                </div>
              )}
            </div>

            {committed && (
              <div className="mt-4 rounded-xl border border-emerald-500/30 bg-emerald-500/[0.08] p-4 text-center">
                <div className="grid place-items-center w-10 h-10 mx-auto rounded-full bg-emerald-500/20 border border-emerald-500/40">
                  <Icon name="check" className="w-5 h-5 text-emerald-300" />
                </div>
                <div className="mt-2 text-sm font-semibold text-emerald-200">設定已套用 · 開始收集資料</div>
                <div className="text-[11px] text-emerald-300/70 mt-1">
                  Scheduler 已啟動 · 第一筆資料預計在 ~{connector.write_interval_seconds}s 後寫入
                </div>
                <Button variant="secondary" size="sm" className="mt-3" onClick={onCommit} icon={<Icon name="eye" className="w-3.5 h-3.5" />}>
                  前往 Runtime Dashboard
                </Button>
              </div>
            )}
          </div>
        )}
      </SectionCard>
    </div>
  );
}

window.Step4Database = Step4Database;
