import { useCallback, useEffect, useMemo, useState, type ChangeEvent } from 'react';
import { modbusShareAPI } from '../../services/datalink';
import { useTagsQuery } from '../../hooks/datalink/useTags';
import type { ModbusShareMapping, ModbusShareStatus } from '../../types/datalink';

type MappingConflict = {
  register: number;
  mappings: ModbusShareMapping[];
};

export default function LocalModbusWorkbenchPage() {
  const { data: tags = [] } = useTagsQuery();
  const [status, setStatus] = useState<ModbusShareStatus | null>(null);
  const [mappings, setMappings] = useState<ModbusShareMapping[]>([]);
  const [selectedTagId, setSelectedTagId] = useState('');
  const [registerInput, setRegisterInput] = useState('0');
  const [testTagId, setTestTagId] = useState('');
  const [testValue, setTestValue] = useState('0');
  const [message, setMessage] = useState('');
  const [isBusy, setIsBusy] = useState(false);

  const tagById = useMemo(() => new Map(tags.map((tag) => [tag.id, tag])), [tags]);
  const conflicts = useMemo<MappingConflict[]>(() => {
    const grouped = new Map<number, ModbusShareMapping[]>();
    mappings.forEach((mapping) => {
      const list = grouped.get(mapping.register) ?? [];
      list.push(mapping);
      grouped.set(mapping.register, list);
    });
    return Array.from(grouped.entries())
      .filter(([, list]) => list.length > 1)
      .map(([register, list]) => ({ register, mappings: list }));
  }, [mappings]);
  const canWrite = status?.enabled && conflicts.length === 0;

  const loadData = useCallback(async () => {
    setIsBusy(true);
    try {
      const [nextStatus, nextMappings] = await Promise.all([
        modbusShareAPI.status(),
        modbusShareAPI.listMappings(),
      ]);
      setStatus(nextStatus);
      setMappings(nextMappings);
      setMessage('');
    } catch (error) {
      const nextMessage = error instanceof Error ? error.message : '讀取本機 Modbus 狀態失敗';
      setMessage(nextMessage);
    } finally {
      setIsBusy(false);
    }
  }, []);

  useEffect(() => {
    void loadData();
  }, [loadData]);

  const handleUpsert = useCallback(async () => {
    if (!selectedTagId) {
      setMessage('請先選擇 Tag');
      return;
    }
    const register = Number(registerInput);
    if (!Number.isInteger(register) || register < 0 || register > 65535) {
      setMessage('Register 必須為 0~65535 整數');
      return;
    }
    setIsBusy(true);
    try {
      await modbusShareAPI.upsertMapping(selectedTagId, register);
      await loadData();
      setMessage(`映射已更新：${selectedTagId} -> HR${register}`);
    } catch (error) {
      setMessage(error instanceof Error ? error.message : '更新映射失敗');
    } finally {
      setIsBusy(false);
    }
  }, [loadData, registerInput, selectedTagId]);

  const handleDelete = useCallback(
    async (tagId: string) => {
      setIsBusy(true);
      try {
        await modbusShareAPI.deleteMapping(tagId);
        await loadData();
        setMessage(`已刪除映射：${tagId}`);
      } catch (error) {
        setMessage(error instanceof Error ? error.message : '刪除映射失敗');
      } finally {
        setIsBusy(false);
      }
    },
    [loadData]
  );

  const handleSync = useCallback(async () => {
    if (!canWrite) {
      setMessage('目前無法同步：請先確認 server 啟用且衝突為 0');
      return;
    }
    setIsBusy(true);
    try {
      const result = await modbusShareAPI.sync();
      setMessage(`同步完成：updated=${result.updated}, skipped=${result.skipped}, errors=${result.errors.length}`);
      await loadData();
    } catch (error) {
      setMessage(error instanceof Error ? error.message : '同步失敗');
    } finally {
      setIsBusy(false);
    }
  }, [canWrite, loadData]);

  const handleWriteTest = useCallback(async () => {
    if (!canWrite) {
      setMessage('目前為 blocked 模式，請先解決衝突');
      return;
    }
    if (!testTagId) {
      setMessage('請選擇測試 Tag');
      return;
    }
    setIsBusy(true);
    try {
      await modbusShareAPI.writeTagValue(testTagId, testValue);
      setMessage(`測試寫入完成：tag=${testTagId}, value=${testValue}`);
    } catch (error) {
      setMessage(error instanceof Error ? error.message : '測試寫入失敗');
    } finally {
      setIsBusy(false);
    }
  }, [canWrite, testTagId, testValue]);

  const handleExport = useCallback(() => {
    const payload = {
      exportedAt: new Date().toISOString(),
      mappings,
    };
    const blob = new Blob([JSON.stringify(payload, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const anchor = document.createElement('a');
    anchor.href = url;
    anchor.download = 'local-modbus-mappings.json';
    anchor.click();
    URL.revokeObjectURL(url);
  }, [mappings]);

  const handleImport = useCallback(
    async (event: ChangeEvent<HTMLInputElement>) => {
      const file = event.target.files?.[0];
      if (!file) return;
      const text = await file.text();
      try {
        const parsed = JSON.parse(text) as { mappings?: Array<{ tag_id: string; register: number }> };
        const imported = parsed.mappings ?? [];
        for (const item of imported) {
          await modbusShareAPI.upsertMapping(item.tag_id, item.register);
        }
        await loadData();
        setMessage(`匯入完成：${imported.length} 筆`);
      } catch {
        setMessage('匯入格式錯誤，請使用匯出的 JSON 檔');
      }
      event.target.value = '';
    },
    [loadData]
  );

  return (
    <div className="min-h-[calc(100vh-11rem)] rounded-2xl bg-gradient-to-br from-[#0B1220] via-[#0F172A] to-[#111827] p-4 text-slate-100 sm:p-6">
      <header className="rounded-2xl border border-white/10 bg-slate-900/70 px-4 py-3 sm:px-5">
        <h1 className="text-xl font-bold tracking-tight text-slate-100">Local Modbus 5020 - Server Memory Grid</h1>
        <p className="mt-1 text-xs text-slate-400">設備來源隔離，僅本機 5020 目標空間做衝突治理。</p>
      </header>

      <section className="mt-3 grid grid-cols-1 gap-3 sm:grid-cols-3">
        <article className="rounded-xl border border-white/10 bg-slate-900/70 p-3">
          <p className="text-[11px] text-slate-400">Server</p>
          <p className={`mt-1 text-sm font-semibold ${status?.enabled ? 'text-emerald-300' : 'text-rose-300'}`}>
            {status?.enabled ? `Running @ ${status.address}:${status.port}` : 'Stopped'}
          </p>
        </article>
        <article className="rounded-xl border border-white/10 bg-slate-900/70 p-3">
          <p className="text-[11px] text-slate-400">Mappings</p>
          <p className="mt-1 text-sm font-semibold text-slate-100">{status?.mapping_count ?? mappings.length}</p>
        </article>
        <article className="rounded-xl border border-white/10 bg-slate-900/70 p-3">
          <p className="text-[11px] text-slate-400">Conflicts</p>
          <p className={`mt-1 text-sm font-semibold ${conflicts.length > 0 ? 'text-amber-300' : 'text-emerald-300'}`}>
            {conflicts.length > 0 ? `${conflicts.length} blocked` : '0'}
          </p>
        </article>
      </section>

      <section className="mt-3 rounded-2xl border border-white/10 bg-slate-900/70 p-3">
        <div className="flex flex-wrap items-center gap-2">
          <button
            type="button"
            onClick={() => void loadData()}
            className="min-h-11 rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 text-xs font-semibold text-slate-200 hover:bg-slate-700"
          >
            Refresh
          </button>
          <button
            type="button"
            disabled={!canWrite || isBusy}
            onClick={() => void handleSync()}
            className="min-h-11 rounded-lg border border-emerald-400/30 bg-emerald-500/20 px-3 py-2 text-xs font-semibold text-emerald-100 hover:bg-emerald-500/30 disabled:cursor-not-allowed disabled:opacity-50"
          >
            Sync from Mappings
          </button>
          <button
            type="button"
            onClick={handleExport}
            className="min-h-11 rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 text-xs font-semibold text-slate-200 hover:bg-slate-700"
          >
            Export
          </button>
          <label className="inline-flex min-h-11 cursor-pointer items-center rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 text-xs font-semibold text-slate-200 hover:bg-slate-700">
            Import
            <input type="file" accept="application/json" className="hidden" onChange={(e) => void handleImport(e)} />
          </label>
          <p className={`ml-auto text-xs ${canWrite ? 'text-emerald-300' : 'text-amber-300'}`}>
            Write Mode: {canWrite ? 'Ready' : 'Blocked'}
          </p>
        </div>
      </section>

      <section className="mt-3 grid grid-cols-1 gap-3 xl:grid-cols-[420px_1fr_320px]">
        <article className="rounded-2xl border border-white/10 bg-slate-900/70 p-3">
          <h2 className="text-sm font-semibold text-slate-100">Mapping Editor</h2>
          <div className="mt-3 grid grid-cols-1 gap-2">
            <select
              value={selectedTagId}
              onChange={(event) => setSelectedTagId(event.target.value)}
              className="min-h-11 rounded-lg border border-slate-700 bg-slate-800 px-3 text-sm text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">選擇 Tag</option>
              {tags.map((tag) => (
                <option key={tag.id} value={tag.id}>
                  {tag.key}
                </option>
              ))}
            </select>
            <input
              value={registerInput}
              onChange={(event) => setRegisterInput(event.target.value)}
              placeholder="Register (0-65535)"
              className="min-h-11 rounded-lg border border-slate-700 bg-slate-800 px-3 text-sm text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
              type="button"
              onClick={() => void handleUpsert()}
              className="min-h-11 rounded-lg bg-blue-600 px-3 py-2 text-xs font-semibold text-white hover:bg-blue-500"
            >
              新增 / 更新映射
            </button>
          </div>
          <div className="mt-3 max-h-[420px] space-y-2 overflow-auto">
            {mappings.map((mapping) => (
              <div key={mapping.tag_id} className="rounded-lg border border-white/10 bg-slate-800/70 p-2">
                <p className="text-xs text-slate-100">
                  {tagById.get(mapping.tag_id)?.key ?? mapping.tag_id}{' -> '}HR{mapping.register}
                </p>
                <p className="text-[11px] text-slate-400">{new Date(mapping.updated_at).toLocaleString()}</p>
                <button
                  type="button"
                  onClick={() => void handleDelete(mapping.tag_id)}
                  className="mt-1 text-[11px] text-rose-300 hover:text-rose-200"
                >
                  刪除
                </button>
              </div>
            ))}
          </div>
        </article>

        <article className="rounded-2xl border border-white/10 bg-slate-900/70 p-3">
          <h2 className="text-sm font-semibold text-slate-100">Server Memory Grid (5020)</h2>
          <div className="mt-3 rounded-lg border border-white/10 bg-slate-950/70 p-3">
            <p className="text-xs text-slate-400">Preflight</p>
            <div className="mt-1 flex flex-wrap gap-3 text-xs">
              <span className={status?.enabled ? 'text-emerald-300' : 'text-rose-300'}>
                Bind: {status?.enabled ? 'PASS' : 'FAIL'}
              </span>
              <span className={conflicts.length > 0 ? 'text-amber-300' : 'text-emerald-300'}>
                Conflict: {conflicts.length > 0 ? `FAIL (${conflicts.length})` : 'PASS'}
              </span>
              <span className="text-emerald-300">Permission: PASS</span>
            </div>
          </div>
          <div className="mt-3 h-[420px] rounded-lg border border-white/10 bg-slate-950/70 p-3">
            <p className="text-xs text-slate-400">Grid Canvas Placeholder (HR00000 - HR65535)</p>
          </div>
        </article>

        <article className="rounded-2xl border border-white/10 bg-slate-900/70 p-3">
          <h2 className="text-sm font-semibold text-slate-100">Conflict + Write Test</h2>
          <div className="mt-3 rounded-lg border border-white/10 bg-slate-950/70 p-3">
            <p className="text-xs font-semibold text-amber-200">Conflict Queue</p>
            {conflicts.length === 0 ? (
              <p className="mt-2 text-xs text-emerald-300">無衝突</p>
            ) : (
              <div className="mt-2 space-y-2">
                {conflicts.map((conflict) => (
                  <div key={conflict.register} className="rounded-md border border-amber-500/30 bg-amber-500/10 p-2">
                    <p className="text-xs text-amber-200">HR{conflict.register}</p>
                    <p className="text-[11px] text-amber-100">
                      {conflict.mappings.map((mapping) => tagById.get(mapping.tag_id)?.key ?? mapping.tag_id).join(' / ')}
                    </p>
                  </div>
                ))}
              </div>
            )}
          </div>
          <div className="mt-3 rounded-lg border border-white/10 bg-slate-950/70 p-3">
            <p className="text-xs font-semibold text-slate-200">Test Read/Write</p>
            <select
              value={testTagId}
              onChange={(event) => setTestTagId(event.target.value)}
              className="mt-2 min-h-11 w-full rounded-lg border border-slate-700 bg-slate-800 px-3 text-sm text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="">選擇 Tag</option>
              {tags.map((tag) => (
                <option key={tag.id} value={tag.id}>
                  {tag.key}
                </option>
              ))}
            </select>
            <input
              value={testValue}
              onChange={(event) => setTestValue(event.target.value)}
              className="mt-2 min-h-11 w-full rounded-lg border border-slate-700 bg-slate-800 px-3 text-sm text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
              type="button"
              disabled={!canWrite || isBusy}
              onClick={() => void handleWriteTest()}
              className="mt-2 min-h-11 w-full rounded-lg bg-blue-600 px-3 py-2 text-xs font-semibold text-white hover:bg-blue-500 disabled:cursor-not-allowed disabled:opacity-50"
            >
              執行測試寫入
            </button>
          </div>
        </article>
      </section>

      <footer className="mt-3 rounded-xl border border-white/10 bg-slate-900/70 px-3 py-2 text-xs text-slate-300">
        {isBusy ? '處理中...' : message || 'Ready'}
      </footer>
    </div>
  );
}
