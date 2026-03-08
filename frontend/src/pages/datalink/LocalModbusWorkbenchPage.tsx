import { useCallback, useEffect, useMemo, useState, type ChangeEvent } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { modbusShareAPI } from '../../services/datalink';
import { useTagsQuery } from '../../hooks/datalink/useTags';
import type { ModbusShareMapping, ModbusShareStatus } from '../../types/datalink';
import { designSystem } from '../../styles/designSystem';

type MappingConflict = {
  register: number;
  mappings: ModbusShareMapping[];
};

export default function LocalModbusWorkbenchPage() {
  const [searchParams] = useSearchParams();
  const { data: tags = [] } = useTagsQuery();
  const [status, setStatus] = useState<ModbusShareStatus | null>(null);
  const [mappings, setMappings] = useState<ModbusShareMapping[]>([]);
  const [selectedTagId, setSelectedTagId] = useState('');
  const [registerInput, setRegisterInput] = useState('0');
  const [serverPortInput, setServerPortInput] = useState('5020');
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
  const section = searchParams.get('section');
  const returnTarget = section ? `/datalink?section=${section}` : '/datalink';

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

  useEffect(() => {
    if (status?.port && Number.isInteger(status.port) && status.port > 0) {
      setServerPortInput(String(status.port));
    }
  }, [status?.port]);

  const handleStartServer = useCallback(async () => {
    const port = Number(serverPortInput);
    if (!Number.isInteger(port) || port <= 0 || port > 65535) {
      setMessage('Port 必須為 1~65535 的整數');
      return;
    }

    setIsBusy(true);
    try {
      await modbusShareAPI.start(port);
      await loadData();
      setMessage(`Server 已啟動並綁定 127.0.0.1:${port}`);
    } catch (error) {
      setMessage(error instanceof Error ? error.message : '啟動 server 失敗');
    } finally {
      setIsBusy(false);
    }
  }, [loadData, serverPortInput]);

  const handleStopServer = useCallback(async () => {
    setIsBusy(true);
    try {
      await modbusShareAPI.stop();
      await loadData();
      setMessage('Server 已停止');
    } catch (error) {
      setMessage(error instanceof Error ? error.message : '停止 server 失敗');
    } finally {
      setIsBusy(false);
    }
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
      setMessage(`${designSystem.microcopy.feedback.success.validated}：updated=${result.updated}, skipped=${result.skipped}, errors=${result.errors.length}`);
      await loadData();
    } catch (error) {
      setMessage(error instanceof Error ? error.message : designSystem.microcopy.feedback.error.validationFailed);
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
        setMessage(`${designSystem.microcopy.feedback.success.created}：${imported.length} 筆`);
      } catch {
        setMessage(designSystem.microcopy.feedback.error.createFailed);
      }
      event.target.value = '';
    },
    [loadData]
  );

  return (
    <div className="min-h-[calc(100vh-11rem)] rounded-2xl bg-gradient-to-br from-[#0B1220] via-[#0F172A] to-[#111827] p-4 text-slate-100 sm:p-6">
      {/* ── Header ── */}
      <header className="rounded-2xl border border-white/10 bg-slate-900/70 px-4 py-3 sm:px-5">
        <div className="flex flex-wrap items-center justify-between gap-2">
          <div>
            <h1 className="text-xl font-bold tracking-tight text-slate-100">Local Modbus 工作台</h1>
            <p className="mt-1 text-xs text-slate-400">將 Tag 映射到本地 Modbus Server，供外部設備讀取</p>
          </div>
          <Link
            to={returnTarget}
            className="min-h-11 rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 text-xs font-semibold text-slate-100 hover:bg-slate-700 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
          >
            返回 Dashboard
          </Link>
        </div>
      </header>

      {/* ── 系統狀態 ── */}
      <section className="mt-4 rounded-2xl border border-white/10 bg-slate-900/70 p-4">
        <div className="flex items-center justify-between mb-3">
          <h2 className="text-sm font-semibold text-slate-200">系統狀態</h2>
          <button
            type="button"
            onClick={() => void loadData()}
            disabled={isBusy}
            className="rounded-lg border border-slate-600 bg-slate-800 px-2.5 py-1 text-xs font-semibold text-slate-200 hover:bg-slate-700 disabled:cursor-not-allowed disabled:opacity-50"
          >
            重新整理
          </button>
        </div>
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
          <article className="rounded-xl border border-white/10 bg-slate-800/50 p-3">
            <p className="text-[11px] font-medium text-slate-400 mb-1">Server 狀態</p>
            <p className={`text-sm font-semibold ${status?.enabled ? 'text-emerald-300' : 'text-rose-300'}`}>
              {status?.enabled ? `運行中 @ ${status.address}` : '已停止'}
            </p>
            {status?.enabled && (
              <p className="mt-1 text-[11px] text-slate-400">Port: {status.port}</p>
            )}
          </article>
          <article className="rounded-xl border border-white/10 bg-slate-800/50 p-3">
            <p className="text-[11px] font-medium text-slate-400 mb-1">映射數量</p>
            <p className="text-sm font-semibold text-slate-100">{status?.mapping_count ?? mappings.length}</p>
            <p className="mt-1 text-[11px] text-slate-400">已設定映射</p>
          </article>
          <article className="rounded-xl border border-white/10 bg-slate-800/50 p-3">
            <p className="text-[11px] font-medium text-slate-400 mb-1">衝突狀態</p>
            <p className={`text-sm font-semibold ${conflicts.length > 0 ? 'text-amber-300' : 'text-emerald-300'}`}>
              {conflicts.length > 0 ? `${conflicts.length} 個衝突` : '無衝突'}
            </p>
            <p className={`mt-1 text-[11px] ${canWrite ? 'text-emerald-300' : 'text-amber-300'}`}>
              {canWrite ? '寫入模式：就緒' : '寫入模式：已阻擋'}
            </p>
          </article>
        </div>
        <div className="mt-3 flex flex-wrap items-center gap-2 pt-3 border-t border-white/5">
          <input
            value={serverPortInput}
            onChange={(event) => setServerPortInput(event.target.value)}
            placeholder="Server Port"
            name={designSystem.forms.name.modbusPort}
            autoComplete={designSystem.forms.autocomplete.port}
            inputMode={designSystem.forms.inputmode.numeric}
            className="min-h-11 w-32 rounded-lg border border-slate-700 bg-slate-800 px-3 text-sm text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
            disabled={isBusy || !!status?.enabled}
          />
          <button
            type="button"
            disabled={isBusy || !!status?.enabled}
            onClick={() => void handleStartServer()}
            className="min-h-11 rounded-lg border border-emerald-400/30 bg-emerald-500/20 px-3 py-2 text-xs font-semibold text-emerald-100 hover:bg-emerald-500/30 disabled:cursor-not-allowed disabled:opacity-50"
          >
            啟動 Server
          </button>
          <button
            type="button"
            disabled={isBusy || !status?.enabled}
            onClick={() => void handleStopServer()}
            className="min-h-11 rounded-lg border border-rose-400/30 bg-rose-500/20 px-3 py-2 text-xs font-semibold text-rose-100 hover:bg-rose-500/30 disabled:cursor-not-allowed disabled:opacity-50"
          >
            停止 Server
          </button>
          <button
            type="button"
            disabled={!canWrite || isBusy}
            onClick={() => void handleSync()}
            className="min-h-11 rounded-lg border border-blue-400/30 bg-blue-500/20 px-3 py-2 text-xs font-semibold text-blue-100 hover:bg-blue-500/30 disabled:cursor-not-allowed disabled:opacity-50"
          >
            同步映射
          </button>
          <button
            type="button"
            onClick={handleExport}
            className="min-h-11 rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 text-xs font-semibold text-slate-200 hover:bg-slate-700"
          >
            匯出
          </button>
          <label className="inline-flex min-h-11 cursor-pointer items-center rounded-lg border border-slate-600 bg-slate-800 px-3 py-2 text-xs font-semibold text-slate-200 hover:bg-slate-700">
            匯入
            <input type="file" accept="application/json" className="hidden" onChange={(e) => void handleImport(e)} />
          </label>
        </div>
      </section>

      {/* ── 主要工作區：Mapping 編輯 / 衝突治理 / 寫入測試 ── */}
      <section className="mt-4 grid grid-cols-1 gap-4 lg:grid-cols-2 xl:grid-cols-[1fr_400px]">
        {/* Mapping 編輯 */}
        <article className="rounded-2xl border border-white/10 bg-slate-900/70 p-4">
          <h2 className="text-sm font-semibold text-slate-200 mb-3">Mapping 編輯</h2>
          <div className="space-y-3">
            <div className="grid grid-cols-1 gap-2">
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
                name={designSystem.forms.name.modbusRegister}
                autoComplete={designSystem.forms.autocomplete.number}
                inputMode={designSystem.forms.inputmode.numeric}
                className="min-h-11 rounded-lg border border-slate-700 bg-slate-800 px-3 text-sm text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                type="button"
                onClick={() => void handleUpsert()}
                disabled={isBusy}
                className="min-h-11 rounded-lg bg-blue-600 px-3 py-2 text-xs font-semibold text-white hover:bg-blue-500 disabled:cursor-not-allowed disabled:opacity-50"
              >
                {isBusy ? designSystem.microcopy.loading.saving : '新增 / 更新映射'}
              </button>
            </div>
            <div className="max-h-[400px] space-y-2 overflow-auto rounded-lg border border-white/5 bg-slate-950/50 p-2">
              {mappings.length === 0 ? (
                <p className="text-center text-xs text-slate-400 py-4">尚無映射</p>
              ) : (
                mappings.map((mapping) => (
                  <div key={mapping.tag_id} className="rounded-lg border border-white/10 bg-slate-800/70 p-2">
                    <div className="flex items-start justify-between">
                      <div className="min-w-0 flex-1">
                        <p className="text-xs font-medium text-slate-100">
                          {tagById.get(mapping.tag_id)?.key ?? mapping.tag_id}
                        </p>
                        <p className="mt-0.5 text-[11px] text-slate-400">HR{mapping.register}</p>
                        <p className="mt-1 text-[10px] text-slate-500">{new Date(mapping.updated_at).toLocaleString()}</p>
                      </div>
                      <button
                        type="button"
                        onClick={() => void handleDelete(mapping.tag_id)}
                        className="ml-2 rounded border border-rose-500/30 bg-rose-500/10 px-2 py-1 text-[11px] text-rose-300 hover:bg-rose-500/20"
                      >
                        刪除
                      </button>
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>
        </article>

        {/* 衝突治理 + 寫入測試 */}
        <div className="space-y-4">
          {/* 衝突治理 */}
          <article className="rounded-2xl border border-white/10 bg-slate-900/70 p-4">
            <h2 className="text-sm font-semibold text-slate-200 mb-3">衝突治理</h2>
            {conflicts.length === 0 ? (
              <div className="rounded-lg border border-emerald-500/20 bg-emerald-500/5 p-3 text-center">
                <p className="text-xs text-emerald-300">無衝突</p>
                <p className="mt-1 text-[11px] text-emerald-200/80">所有映射皆可正常寫入</p>
              </div>
            ) : (
              <div className="space-y-2">
                {conflicts.map((conflict) => (
                  <div key={conflict.register} className="rounded-lg border border-amber-500/30 bg-amber-500/10 p-3">
                    <p className="text-xs font-semibold text-amber-200">HR{conflict.register}</p>
                    <p className="mt-1 text-[11px] text-amber-100">
                      {conflict.mappings.map((mapping) => tagById.get(mapping.tag_id)?.key ?? mapping.tag_id).join(' / ')}
                    </p>
                    <p className="mt-1 text-[10px] text-amber-200/70">請刪除多餘映射以解決衝突</p>
                  </div>
                ))}
              </div>
            )}
          </article>

          {/* 寫入測試 */}
          <article className="rounded-2xl border border-white/10 bg-slate-900/70 p-4">
            <h2 className="text-sm font-semibold text-slate-200 mb-3">寫入測試</h2>
            <div className="space-y-2">
              <select
                value={testTagId}
                onChange={(event) => setTestTagId(event.target.value)}
                className="w-full min-h-11 rounded-lg border border-slate-700 bg-slate-800 px-3 text-sm text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
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
                placeholder="測試數值"
                name="test-value"
                autoComplete={designSystem.forms.autocomplete.number}
                inputMode={designSystem.forms.inputmode.decimal}
                className="w-full min-h-11 rounded-lg border border-slate-700 bg-slate-800 px-3 text-sm text-slate-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                type="button"
                disabled={!canWrite || isBusy || !testTagId}
                onClick={() => void handleWriteTest()}
                className="w-full min-h-11 rounded-lg bg-blue-600 px-3 py-2 text-xs font-semibold text-white hover:bg-blue-500 disabled:cursor-not-allowed disabled:opacity-50"
              >
                {isBusy ? designSystem.microcopy.loading.saving : '執行測試寫入'}
              </button>
              {!canWrite && (
                <p className="text-[11px] text-amber-300">請先解決衝突或啟動 Server</p>
              )}
            </div>
          </article>
        </div>
      </section>

      {/* ── 狀態訊息 ── */}
      {message && (
        <footer 
          className="mt-4 rounded-xl border border-white/10 bg-slate-900/70 px-4 py-2 text-xs text-slate-300"
          role="status"
          aria-live="polite"
          aria-atomic="true"
        >
          {isBusy ? designSystem.microcopy.loading.default : message}
        </footer>
      )}
    </div>
  );
}
