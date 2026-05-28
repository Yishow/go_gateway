// Step 3: 點位映射到 Tag
// 表格 + 即時預覽轉換管線

function Step3Mapping({ state, setState, onContinue }) {
  const { points = [], rules = [] } = state;
  const enabledPoints = points.filter((p) => !p.skipped);

  const [selectedIdx, setSelectedIdx] = React.useState(0);
  const selected = enabledPoints[selectedIdx] || enabledPoints[0];

  // 初始化每個點位的 mapping (從點位身上的 rule 預設值帶入)
  React.useEffect(() => {
    setState((s) => {
      const existing = s.mappings || {};
      const next = {};
      enabledPoints.forEach((p, i) => {
        const semantic = POINT_SEMANTIC[i % POINT_SEMANTIC.length];
        next[p.id] = existing[p.id] || {
          point_id: p.id,
          tag_key: semantic.tag_key,
          display_name: semantic.display,
          unit: semantic.unit,
          target_type: 'float64',
          scale: p._rule_scale ?? 1,
          offset: p._rule_offset ?? 0,
          enabled: true,
        };
      });
      return { ...s, mappings: next };
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [points.length]);

  const mappings = state.mappings || {};
  const updateMapping = (pid, patch) =>
    setState((s) => ({ ...s, mappings: { ...s.mappings, [pid]: { ...s.mappings[pid], ...patch } } }));

  // 模擬 raw 值 (顯示用)
  const sampleRaw = (i) => {
    const seeds = [243, 251, 1024, 985, 67, 542, 18, 1450];
    return seeds[i % seeds.length];
  };

  const computePreview = (raw, m) => {
    if (m?.target_type === 'bool') return raw !== 0 ? 'true' : 'false';
    const v = (raw * (m?.scale ?? 1) + (m?.offset ?? 0));
    const isFloat = String(m?.target_type || '').startsWith('float');
    return isFloat ? v.toFixed(2) : Math.round(v).toString();
  };

  const allHasTag = enabledPoints.every((p) => mappings[p.id]?.tag_key?.trim());

  return (
    <div className="grid grid-cols-12 gap-5 sweep-in">
      {/* 左：映射表 */}
      <SectionCard
        className="col-span-12 lg:col-span-8"
        title="點位 → Tag 映射"
        subtitle="每個點位對應一個 Tag (時序資料庫的命名鍵)"
        icon={<Icon name="map" className="w-5 h-5" />}
        aside={<StatusChip tone={allHasTag ? 'success' : 'warning'}>
          {Object.values(mappings).filter((m) => m.enabled).length} / {enabledPoints.length} 已啟用
        </StatusChip>}
        contentClassName="p-0"
      >
        <div className="overflow-x-auto">
          <table className="w-full text-xs">
            <thead className="bg-slate-900/70 text-slate-500 border-b border-slate-800">
              <tr>
                <th className="px-2 py-2.5 text-left font-medium w-1"></th>
                <th className="px-3 py-2.5 text-left font-medium">點位</th>
                <th className="px-3 py-2.5 text-left font-medium">位址</th>
                <th className="px-3 py-2.5 text-left font-medium w-[22%]">Tag Key</th>
                <th className="px-3 py-2.5 text-left font-medium">顯示名稱</th>
                <th className="px-3 py-2.5 text-left font-medium">單位</th>
                <th className="px-3 py-2.5 text-left font-medium">目標型態</th>
                <th className="px-3 py-2.5 text-right font-medium" title="raw × scale + offset">× Scale</th>
                <th className="px-3 py-2.5 text-right font-medium">+ Offset</th>
                <th className="px-3 py-2.5 text-right font-medium">啟用</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/60">
              {enabledPoints.map((p, i) => {
                const m = mappings[p.id] || {};
                const isSelected = i === selectedIdx;
                // 規則顏色 (對齊 Step 2)
                const ruleColors = ['blue', 'emerald', 'amber', 'fuchsia', 'cyan', 'rose'];
                const ruleIdx = rules.findIndex((r) => r.id === p.rule_id);
                const colorName = ruleColors[(ruleIdx >= 0 ? ruleIdx : 0) % ruleColors.length];
                const dotClass = {
                  blue: 'bg-blue-500', emerald: 'bg-emerald-500', amber: 'bg-amber-500',
                  fuchsia: 'bg-fuchsia-500', cyan: 'bg-cyan-500', rose: 'bg-rose-500',
                }[colorName];
                return (
                  <tr
                    key={p.id}
                    onClick={() => setSelectedIdx(i)}
                    className={clsx(
                      'cursor-pointer transition-colors',
                      isSelected ? 'bg-blue-500/[0.06] outline outline-1 -outline-offset-1 outline-blue-500/30' : 'hover:bg-slate-800/40',
                    )}
                  >
                    <td className="px-2 py-1.5">
                      <span className={clsx('block w-1.5 h-1.5 rounded-full', dotClass)} title={p.rule_name} />
                    </td>
                    <td className="px-3 py-1.5 font-mono text-slate-200 whitespace-nowrap">{p.name}</td>
                    <td className="px-3 py-1.5 font-mono text-blue-200">{p.address}</td>
                    <td className="px-3 py-1.5">
                      <input
                        value={m.tag_key || ''}
                        onChange={(e) => updateMapping(p.id, { tag_key: e.target.value })}
                        onClick={(e) => e.stopPropagation()}
                        className="w-full bg-transparent border border-slate-700/50 rounded px-2 py-1 font-mono text-[11px] text-slate-100 focus:border-blue-500/60 focus:outline-none"
                      />
                    </td>
                    <td className="px-3 py-1.5">
                      <input
                        value={m.display_name || ''}
                        onChange={(e) => updateMapping(p.id, { display_name: e.target.value })}
                        onClick={(e) => e.stopPropagation()}
                        className="w-full bg-transparent border border-slate-700/50 rounded px-2 py-1 text-[11px] text-slate-200 focus:border-blue-500/60 focus:outline-none"
                      />
                    </td>
                    <td className="px-3 py-1.5">
                      <input
                        value={m.unit || ''}
                        onChange={(e) => updateMapping(p.id, { unit: e.target.value })}
                        onClick={(e) => e.stopPropagation()}
                        className="w-14 bg-transparent border border-slate-700/50 rounded px-2 py-1 text-[11px] text-slate-300 focus:border-blue-500/60 focus:outline-none font-mono"
                      />
                    </td>
                    <td className="px-3 py-1.5">
                      <select
                        value={m.target_type || 'float64'}
                        onChange={(e) => updateMapping(p.id, { target_type: e.target.value })}
                        onClick={(e) => e.stopPropagation()}
                        className="bg-slate-950/60 border border-slate-700/50 rounded px-1.5 py-1 text-[11px] font-mono text-slate-200 focus:border-blue-500/60 focus:outline-none"
                      >
                        {DATA_TYPES.map((t) => <option key={t} value={t}>{t}</option>)}
                      </select>
                    </td>
                    <td className="px-3 py-1.5 text-right">
                      <input
                        type="number" step="0.01"
                        value={m.scale ?? 1}
                        onChange={(e) => updateMapping(p.id, { scale: +e.target.value })}
                        onClick={(e) => e.stopPropagation()}
                        className="w-16 bg-slate-950/40 border border-slate-700/50 rounded px-2 py-1 font-mono text-[11px] text-blue-200 text-right focus:border-blue-500/60 focus:outline-none"
                      />
                    </td>
                    <td className="px-3 py-1.5 text-right">
                      <input
                        type="number" step="0.01"
                        value={m.offset ?? 0}
                        onChange={(e) => updateMapping(p.id, { offset: +e.target.value })}
                        onClick={(e) => e.stopPropagation()}
                        className="w-16 bg-slate-950/40 border border-slate-700/50 rounded px-2 py-1 font-mono text-[11px] text-blue-200 text-right focus:border-blue-500/60 focus:outline-none"
                      />
                    </td>
                    <td className="px-3 py-1.5 text-right">
                      <div onClick={(e) => e.stopPropagation()} className="inline-block">
                        <Toggle checked={m.enabled ?? true} onChange={(v) => updateMapping(p.id, { enabled: v })} size="sm" />
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
        <div className="border-t border-slate-800 px-4 py-2.5 text-[11px] text-slate-500 flex items-center justify-between">
          <span>
            <Icon name="info" className="inline w-3.5 h-3.5 mr-1" />
            點選列以在右側看轉換管線預覽
          </span>
          <button
            onClick={() => {
              // 批次套用：將選中列的 scale / offset 套用到全部
              const m = mappings[selected?.id] || {};
              setState((s) => {
                const next = { ...s.mappings };
                Object.keys(next).forEach((k) => {
                  next[k] = { ...next[k], scale: m.scale, offset: m.offset, target_type: m.target_type };
                });
                return { ...s, mappings: next };
              });
            }}
            className="flex items-center gap-1 text-slate-400 hover:text-slate-200 transition-colors"
          >
            <Icon name="bolt" className="w-3.5 h-3.5" /> 將選中列的 Scale / Offset / 型態 套用到全部
          </button>
        </div>
      </SectionCard>

      {/* 右：轉換管線預覽 */}
      <SectionCard
        className="col-span-12 lg:col-span-4"
        title="轉換管線預覽"
        subtitle={selected ? `${selected.name} (${selected.address})` : '請選擇一個點位'}
        icon={<Icon name="flow" className="w-5 h-5" />}
      >
        {selected && mappings[selected.id] && (() => {
          const m = mappings[selected.id];
          const raw = sampleRaw(selectedIdx);
          const scaled = raw * (m.scale ?? 1) + (m.offset ?? 0);
          const final = computePreview(raw, m);
          return (
            <div className="space-y-4">
              {/* 轉換步驟 */}
              <ol className="relative space-y-2 before:absolute before:left-3.5 before:top-2 before:bottom-2 before:w-px before:bg-slate-700">
                <li className="relative flex gap-3 items-start">
                  <span className="grid place-items-center w-7 h-7 rounded-full bg-slate-800 border border-slate-700 text-[10px] font-bold text-slate-300 z-10">1</span>
                  <div className="flex-1 rounded-lg border border-slate-700/60 bg-slate-950/40 p-2.5">
                    <div className="text-[10px] uppercase tracking-wider text-slate-500">decode (raw)</div>
                    <div className="font-mono text-lg text-slate-100 mt-0.5">{raw}</div>
                    <div className="text-[10px] text-slate-500 mt-0.5">register @ {selected.address}, type={selected.data_type}</div>
                  </div>
                </li>
                <li className="relative flex gap-3 items-start">
                  <span className="grid place-items-center w-7 h-7 rounded-full bg-slate-800 border border-slate-700 text-[10px] font-bold text-slate-300 z-10">2</span>
                  <div className="flex-1 rounded-lg border border-slate-700/60 bg-slate-950/40 p-2.5">
                    <div className="text-[10px] uppercase tracking-wider text-slate-500">scale (linear)</div>
                    <div className="font-mono text-sm text-slate-200 mt-0.5">
                      <span className="text-slate-500">{raw}</span>
                      <span className="mx-1 text-slate-600">×</span>
                      <span className="text-blue-300">{m.scale}</span>
                      <span className="mx-1 text-slate-600">+</span>
                      <span className="text-blue-300">{m.offset}</span>
                      <span className="mx-1 text-slate-600">=</span>
                      <span className="text-emerald-300 font-semibold">{scaled.toFixed(3)}</span>
                    </div>
                  </div>
                </li>
                <li className="relative flex gap-3 items-start">
                  <span className="grid place-items-center w-7 h-7 rounded-full bg-slate-800 border border-slate-700 text-[10px] font-bold text-slate-300 z-10">3</span>
                  <div className="flex-1 rounded-lg border border-slate-700/60 bg-slate-950/40 p-2.5">
                    <div className="text-[10px] uppercase tracking-wider text-slate-500">cast to {m.target_type}</div>
                    <div className="font-mono text-sm text-slate-200 mt-0.5">cast({scaled.toFixed(3)} → {m.target_type})</div>
                  </div>
                </li>
                <li className="relative flex gap-3 items-start">
                  <span className="grid place-items-center w-7 h-7 rounded-full bg-emerald-500/20 border border-emerald-500/40 text-[10px] font-bold text-emerald-300 z-10">✓</span>
                  <div className="flex-1 rounded-lg border border-emerald-500/30 bg-emerald-500/[0.06] p-2.5">
                    <div className="text-[10px] uppercase tracking-wider text-emerald-400/80">final → Tag</div>
                    <div className="flex items-end justify-between mt-0.5">
                      <span className="font-mono text-2xl font-semibold text-emerald-200">{final}</span>
                      <span className="font-mono text-xs text-emerald-300/80">{m.unit}</span>
                    </div>
                    <div className="font-mono text-[10px] text-emerald-300/60 mt-1">→ tag.{m.tag_key}</div>
                  </div>
                </li>
              </ol>

              {/* 範例 payload */}
              <div className="rounded-lg border border-slate-700/60 bg-slate-950/60 p-2.5">
                <div className="text-[10px] uppercase tracking-wider text-slate-500 mb-1">POST /mappings</div>
                <pre className="text-[10px] font-mono leading-relaxed text-slate-400 overflow-x-auto">
{`{
  "point_id": "${selected.id.slice(0, 14)}…",
  "tag_id": "tag.${m.tag_key}",
  "transform_pipeline": [
    { "type": "decode" },
    { "type": "scale",
      "params": { "scale": ${m.scale}, "offset": ${m.offset} }},
    { "type": "cast",
      "params": { "to": "${m.target_type}" }}
  ],
  "enabled": ${m.enabled}
}`}
                </pre>
              </div>
            </div>
          );
        })()}

        <div className="mt-4 flex justify-end">
          <Button onClick={onContinue} disabled={!allHasTag} icon={<Icon name="arrow" className="w-4 h-4" />}>
            設定資料庫寫入
          </Button>
        </div>
      </SectionCard>
    </div>
  );
}

window.Step3Mapping = Step3Mapping;
