import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useDevicesQuery } from '../../../hooks/datalink/useDevices';
import { useMappingsQuery } from '../../../hooks/datalink/useMappings';
import { usePointsQuery } from '../../../hooks/datalink/usePoints';
import { useTagsQuery } from '../../../hooks/datalink/useTags';
import { dbTargetAPI, modbusShareAPI } from '../../../services/datalink';
import type {
  DatabaseTargetMapping,
  Device,
  ModbusShareMapping,
  ModbusShareStatus,
} from '../../../types/datalink';
import { DatabaseTargetBoard } from './DatabaseTargetBoard';
import { useWorkbench } from './WorkbenchProvider';
import type { WorkbenchOutputCandidate } from './workbenchOutputTypes';

type MappingConflict = {
  register: number;
  mappings: ModbusShareMapping[];
};

type OutputCandidate = WorkbenchOutputCandidate & {
  register: number | null;
  databasePath: string | null;
};

function getSelectedDevice(devices: Device[], selectedDeviceId: string | null) {
  return devices.find((device) => device.id === selectedDeviceId) ?? null;
}

function getErrorMessage(error: unknown, fallback: string): string {
  return error instanceof Error ? error.message : fallback;
}

export function LocalModbusBoard() {
  const { t } = useTranslation();
  const tRef = useRef(t);
  tRef.current = t;
  const {
    activeOutputTarget,
    selectedDeviceId,
    setActiveOutputTarget,
    setSelectedDeviceId,
  } = useWorkbench();
  const { data: devices = [] } = useDevicesQuery();
  const { data: tags = [] } = useTagsQuery();
  const { data: mappings = [] } = useMappingsQuery();
  const { data: points = [] } = usePointsQuery(
    selectedDeviceId ? { device_id: selectedDeviceId } : undefined,
  );

  const [status, setStatus] = useState<ModbusShareStatus | null>(null);
  const [shareMappings, setShareMappings] = useState<ModbusShareMapping[]>([]);
  const [dbTargetMappings, setDBTargetMappings] = useState<DatabaseTargetMapping[]>([]);
  const [selectedTagId, setSelectedTagId] = useState('');
  const [registerInput, setRegisterInput] = useState('0');
  const [serverPortInput, setServerPortInput] = useState('5020');
  const [message, setMessage] = useState('');
  const [isBusy, setIsBusy] = useState(true);

  const selectedDevice = getSelectedDevice(devices, selectedDeviceId);
  const pointById = useMemo(
    () => new Map(points.map((point) => [point.id, point])),
    [points],
  );
  const tagById = useMemo(
    () => new Map(tags.map((tag) => [tag.id, tag])),
    [tags],
  );

  const loadData = useCallback(async () => {
    setIsBusy(true);

    try {
      const [nextStatus, nextMappings, nextDBMappings] = await Promise.all([
        modbusShareAPI.status(),
        modbusShareAPI.listMappings(),
        dbTargetAPI.listMappings(),
      ]);

      setStatus(nextStatus);
      setShareMappings(nextMappings);
      setDBTargetMappings(nextDBMappings);
      setMessage('');
    } catch (error) {
      setMessage(
        getErrorMessage(
          error,
          tRef.current('workbench.output.results.statusFallback'),
        ),
      );
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

  const candidates = useMemo(() => {
    const pointIds = new Set(points.map((point) => point.id));
    const shareMappingByTagId = new Map(
      shareMappings.map((mapping) => [mapping.tag_id, mapping]),
    );
    const dbMappingByTagId = new Map(
      dbTargetMappings.map((mapping) => [mapping.tag_id, mapping]),
    );
    const nextCandidates: OutputCandidate[] = [];

    mappings
      .filter((mapping) => pointIds.has(mapping.point_id))
      .forEach((mapping) => {
        const point = pointById.get(mapping.point_id);
        const tag = tagById.get(mapping.tag_id);

        if (!point || !tag) {
          return;
        }

        const shareMapping = shareMappingByTagId.get(tag.id);
        const dbMapping = dbMappingByTagId.get(tag.id);

        nextCandidates.push({
          tagId: tag.id,
          tagKey: tag.key,
          pointName: point.name,
          pointAddress: point.address,
          dataType: tag.data_type,
          register: shareMapping?.register ?? null,
          databasePath: dbMapping
            ? `${dbMapping.table_schema}.${dbMapping.table_name}.${dbMapping.column_name}`
            : null,
          lastValue: point.last_value,
        });
      });

    return nextCandidates;
  }, [dbTargetMappings, mappings, pointById, points, shareMappings, tagById]);

  const candidateKey = candidates.map((candidate) => candidate.tagId).join('|');
  const selectedCandidate =
    candidates.find((candidate) => candidate.tagId === selectedTagId) ?? null;
  const selectedCandidateRegister = selectedCandidate?.register ?? null;

  useEffect(() => {
    if (candidates.length === 0) {
      setSelectedTagId('');
      setRegisterInput('0');
      return;
    }

    setSelectedTagId((previous) =>
      candidates.some((candidate) => candidate.tagId === previous)
        ? previous
        : candidates[0].tagId,
    );
  }, [candidateKey, candidates]);

  useEffect(() => {
    if (!selectedTagId) {
      return;
    }

    setRegisterInput(
      selectedCandidateRegister !== null ? String(selectedCandidateRegister) : '0',
    );
  }, [selectedCandidateRegister, selectedTagId]);

  const conflicts = useMemo<MappingConflict[]>(() => {
    const grouped = new Map<number, ModbusShareMapping[]>();

    shareMappings.forEach((mapping) => {
      const list = grouped.get(mapping.register) ?? [];
      list.push(mapping);
      grouped.set(mapping.register, list);
    });

    return Array.from(grouped.entries())
      .filter(([, items]) => items.length > 1)
      .map(([register, items]) => ({ register, mappings: items }));
  }, [shareMappings]);

  const canSync = Boolean(status?.enabled) && conflicts.length === 0;

  const handleStartServer = useCallback(async () => {
    const port = Number(serverPortInput);

    if (!Number.isInteger(port) || port <= 0 || port > 65535) {
      setMessage(t('workbench.output.results.invalidPort'));
      return;
    }

    setIsBusy(true);

    try {
      await modbusShareAPI.start(port);
      await loadData();
      setMessage(t('workbench.output.results.serverStarted', { port }));
    } catch (error) {
      setMessage(
        getErrorMessage(error, t('workbench.output.results.startFallback')),
      );
    } finally {
      setIsBusy(false);
    }
  }, [loadData, serverPortInput, t]);

  const handleStopServer = useCallback(async () => {
    setIsBusy(true);

    try {
      await modbusShareAPI.stop();
      await loadData();
      setMessage(t('workbench.output.results.serverStopped'));
    } catch (error) {
      setMessage(
        getErrorMessage(error, t('workbench.output.results.stopFallback')),
      );
    } finally {
      setIsBusy(false);
    }
  }, [loadData, t]);

  const handleBind = useCallback(async () => {
    if (!selectedCandidate) {
      setMessage(t('workbench.output.results.noTagSelected'));
      return;
    }

    const register = Number(registerInput);
    if (!Number.isInteger(register) || register < 0 || register > 65535) {
      setMessage(t('workbench.output.results.invalidRegister'));
      return;
    }

    setIsBusy(true);

    try {
      await modbusShareAPI.upsertMapping(selectedCandidate.tagId, register);
      await loadData();
      setMessage(
        t('workbench.output.results.mappingUpdated', {
          key: selectedCandidate.tagKey,
          register,
        }),
      );
    } catch (error) {
      setMessage(
        getErrorMessage(error, t('workbench.output.results.mappingFallback')),
      );
    } finally {
      setIsBusy(false);
    }
  }, [loadData, registerInput, selectedCandidate, t]);

  const handleDelete = useCallback(async () => {
    if (!selectedCandidate) {
      setMessage(t('workbench.output.results.noTagSelected'));
      return;
    }

    setIsBusy(true);

    try {
      await modbusShareAPI.deleteMapping(selectedCandidate.tagId);
      await loadData();
      setMessage(
        t('workbench.output.results.mappingDeleted', {
          key: selectedCandidate.tagKey,
        }),
      );
    } catch (error) {
      setMessage(
        getErrorMessage(error, t('workbench.output.results.deleteFallback')),
      );
    } finally {
      setIsBusy(false);
    }
  }, [loadData, selectedCandidate, t]);

  const handleSync = useCallback(async () => {
    if (!canSync) {
      setMessage(t('workbench.output.results.syncBlocked'));
      return;
    }

    setIsBusy(true);

    try {
      const result = await modbusShareAPI.sync();
      await loadData();
      setMessage(
        t('workbench.output.results.syncSummary', {
          updated: result.updated,
          skipped: result.skipped,
          errors: result.errors.length,
        }),
      );
    } catch (error) {
      setMessage(getErrorMessage(error, t('workbench.output.results.syncFallback')));
    } finally {
      setIsBusy(false);
    }
  }, [canSync, loadData, t]);

  const handlePushValue = useCallback(async () => {
    if (!selectedCandidate) {
      setMessage(t('workbench.output.results.noTagSelected'));
      return;
    }

    if (selectedCandidate.lastValue === null || selectedCandidate.lastValue === undefined) {
      setMessage(t('workbench.output.results.noPointValue'));
      return;
    }

    setIsBusy(true);

    try {
      await modbusShareAPI.writeTagValue(selectedCandidate.tagId, selectedCandidate.lastValue);
      setMessage(
        t('workbench.output.results.valuePushed', {
          key: selectedCandidate.tagKey,
        }),
      );
    } catch (error) {
      setMessage(
        getErrorMessage(error, t('workbench.output.results.writeFallback')),
      );
    } finally {
      setIsBusy(false);
    }
  }, [selectedCandidate, t]);

  if (!selectedDevice) {
    return (
      <section className="space-y-6 rounded-2xl border border-dashed border-slate-700 bg-slate-950/40 p-6">
        <div className="space-y-2">
          <p className="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-300">
            {t('workbench.output.empty.eyebrow')}
          </p>
          <h2 className="text-2xl font-semibold text-slate-50">
            {t('workbench.output.empty.deviceTitle')}
          </h2>
          <p className="max-w-2xl text-sm text-slate-300">
            {t('workbench.output.empty.deviceDescription')}
          </p>
        </div>

        <div className="grid gap-3 sm:grid-cols-2">
          {devices.map((device) => (
            <button
              key={device.id}
              type="button"
              onClick={() => setSelectedDeviceId(device.id)}
              aria-label={device.name}
              className="rounded-2xl border border-slate-800 bg-slate-900/80 p-4 text-left transition hover:border-cyan-500/40 hover:bg-slate-900"
            >
              <p className="text-sm font-semibold text-slate-50">{device.name}</p>
              <p className="mt-1 text-xs text-slate-400">{device.protocol}</p>
            </button>
          ))}
        </div>
      </section>
    );
  }

  if (candidates.length === 0) {
    return (
      <section className="space-y-6 rounded-2xl border border-dashed border-slate-700 bg-slate-950/40 p-6">
        <div className="space-y-2">
          <p className="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-300">
            {t('workbench.output.empty.eyebrow')}
          </p>
          <h2 className="text-2xl font-semibold text-slate-50">
            {t('workbench.output.empty.title')}
          </h2>
          <p className="max-w-2xl text-sm text-slate-300">
            {t('workbench.output.empty.description')}
          </p>
        </div>
      </section>
    );
  }

  return (
    <div className="space-y-6">
      <section className="space-y-4 rounded-2xl border border-slate-800 bg-slate-950/40 p-5">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
          <div className="space-y-2">
            <p className="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-300">
              {t('workbench.output.selection.eyebrow')}
            </p>
            <div>
              <h2 className="text-2xl font-semibold text-slate-50">
                {t('workbench.output.selection.title')}
              </h2>
              <p className="mt-1 max-w-2xl text-sm text-slate-300">
                {t('workbench.output.selection.description')}
              </p>
            </div>
          </div>

          <div className="flex flex-wrap gap-2">
            <button
              type="button"
              aria-pressed={activeOutputTarget === 'modbus'}
              onClick={() => setActiveOutputTarget('modbus')}
              className={
                activeOutputTarget === 'modbus'
                  ? 'rounded-lg bg-cyan-500 px-3 py-2 text-sm font-semibold text-slate-950'
                  : 'rounded-lg border border-slate-800 px-3 py-2 text-sm text-slate-300'
              }
            >
              {t('workbench.output.targetSwitcher.modbus')}
            </button>
            <button
              type="button"
              aria-pressed={activeOutputTarget === 'database'}
              onClick={() => setActiveOutputTarget('database')}
              className={
                activeOutputTarget === 'database'
                  ? 'rounded-lg bg-cyan-500 px-3 py-2 text-sm font-semibold text-slate-950'
                  : 'rounded-lg border border-slate-800 px-3 py-2 text-sm text-slate-300'
              }
            >
              {t('workbench.output.targetSwitcher.database')}
            </button>
          </div>
        </div>

        <div className="grid gap-3">
          {candidates.map((candidate) => {
            const active = selectedTagId === candidate.tagId;

            return (
              <button
                key={candidate.tagId}
                type="button"
                onClick={() => setSelectedTagId(candidate.tagId)}
                aria-pressed={active}
                data-testid={`output-candidate-${candidate.tagId}`}
                className={`grid gap-3 rounded-2xl border p-4 text-left transition xl:grid-cols-[minmax(0,1.1fr)_minmax(0,0.9fr)] ${
                  active
                    ? 'border-cyan-500/40 bg-cyan-500/5'
                    : 'border-slate-800 bg-slate-900/70'
                }`}
              >
                <span className="space-y-2">
                  <span className="block text-sm font-semibold text-slate-50">
                    {candidate.tagKey}
                  </span>
                  <span className="block text-xs text-slate-400">
                    {t('workbench.output.selection.candidateMeta', {
                      point: candidate.pointName,
                      address: candidate.pointAddress,
                      dataType: candidate.dataType,
                    })}
                  </span>
                  <span className="block text-xs text-slate-500">
                    {t('workbench.output.selection.latestValue', {
                      value:
                        candidate.lastValue === null || candidate.lastValue === undefined
                          ? '—'
                          : String(candidate.lastValue),
                    })}
                  </span>
                </span>

                <span className="flex flex-wrap gap-2">
                  <span
                    className={`rounded-full px-3 py-1 text-xs font-medium ${
                      candidate.register !== null
                        ? 'bg-emerald-500/10 text-emerald-200'
                        : 'bg-slate-800 text-slate-300'
                    }`}
                    data-testid={`output-modbus-status-${candidate.tagId}`}
                  >
                    {candidate.register !== null
                      ? `HR${candidate.register}`
                      : t('workbench.output.selection.unmapped')}
                  </span>
                  <span
                    className={`rounded-full px-3 py-1 text-xs font-medium ${
                      candidate.databasePath
                        ? 'bg-violet-500/10 text-violet-200'
                        : 'bg-slate-800 text-slate-300'
                    }`}
                    data-testid={`output-db-status-${candidate.tagId}`}
                  >
                    {candidate.databasePath
                      ? candidate.databasePath
                      : t('workbench.output.selection.databaseUnmapped')}
                  </span>
                </span>
              </button>
            );
          })}
        </div>
      </section>

      {activeOutputTarget === 'modbus' ? (
        <section className="grid gap-6 xl:grid-cols-[minmax(0,1.05fr)_minmax(320px,0.95fr)]">
          <div className="space-y-4 rounded-2xl border border-slate-800 bg-slate-950/40 p-5">
            <div className="space-y-2">
              <p className="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-300">
                {t('workbench.output.modbusStudio.eyebrow')}
              </p>
              <h3 className="text-xl font-semibold text-slate-50">
                {t('workbench.output.modbusStudio.title')}
              </h3>
            </div>

            <div className="grid gap-3 sm:grid-cols-2">
              {candidates.map((candidate) => (
                <article
                  key={`register-map-${candidate.tagId}`}
                  className="rounded-2xl border border-slate-800 bg-slate-900/70 p-4"
                >
                  <p className="text-sm font-semibold text-slate-50">{candidate.tagKey}</p>
                  <p className="mt-1 text-xs text-slate-400">
                    {candidate.pointAddress} · {candidate.pointName}
                  </p>
                  <p className="mt-3 text-sm text-slate-200">
                    {candidate.register !== null
                      ? t('workbench.output.selection.mapped', {
                          register: candidate.register,
                        })
                      : t('workbench.output.selection.unmapped')}
                  </p>
                </article>
              ))}
            </div>
          </div>

          <aside className="space-y-6 rounded-2xl border border-slate-800 bg-slate-950/40 p-5">
            <div className="grid gap-3 sm:grid-cols-3 xl:grid-cols-1">
              <article className="rounded-xl border border-slate-800 bg-slate-900/70 p-4">
                <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
                  {t('workbench.output.status.server')}
                </p>
                <p className="mt-2 text-lg font-semibold text-slate-50">
                  {status?.enabled
                    ? t('workbench.output.status.serverRunning', { address: status.address })
                    : t('workbench.output.status.serverStopped')}
                </p>
              </article>
              <article className="rounded-xl border border-slate-800 bg-slate-900/70 p-4">
                <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
                  {t('workbench.output.status.mappingCount')}
                </p>
                <p className="mt-2 text-lg font-semibold text-slate-50">
                  {status?.mapping_count ?? shareMappings.length}
                </p>
              </article>
              <article className="rounded-xl border border-slate-800 bg-slate-900/70 p-4">
                <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
                  {t('workbench.output.status.conflicts')}
                </p>
                <p className="mt-2 text-lg font-semibold text-slate-50">
                  {conflicts.length}
                </p>
              </article>
            </div>

            <div className="space-y-3 rounded-2xl border border-slate-800 bg-slate-900/70 p-4">
              <div className="grid gap-3">
                <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                  <span>{t('workbench.output.server.port')}</span>
                  <input
                    aria-label={t('workbench.output.server.port')}
                    value={serverPortInput}
                    onChange={(event) => setServerPortInput(event.target.value)}
                    className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
                    disabled={isBusy || Boolean(status?.enabled)}
                  />
                </label>
              </div>

              <div className="flex flex-wrap gap-2">
                <button
                  type="button"
                  disabled={isBusy || Boolean(status?.enabled)}
                  onClick={() => void handleStartServer()}
                  className="rounded-lg border border-emerald-400/30 bg-emerald-500/20 px-3 py-2 text-sm font-medium text-emerald-100 disabled:opacity-50"
                >
                  {t('workbench.output.actions.startServer')}
                </button>
                <button
                  type="button"
                  disabled={isBusy || !status?.enabled}
                  onClick={() => void handleStopServer()}
                  className="rounded-lg border border-rose-400/30 bg-rose-500/20 px-3 py-2 text-sm font-medium text-rose-100 disabled:opacity-50"
                >
                  {t('workbench.output.actions.stopServer')}
                </button>
                <button
                  type="button"
                  disabled={isBusy || !canSync}
                  onClick={() => void handleSync()}
                  className="rounded-lg border border-cyan-400/30 bg-cyan-500/20 px-3 py-2 text-sm font-medium text-cyan-100 disabled:opacity-50"
                >
                  {t('workbench.output.actions.sync')}
                </button>
              </div>
            </div>

            <div className="space-y-3 rounded-2xl border border-slate-800 bg-slate-900/70 p-4">
              <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                <span>{t('workbench.output.mapping.tag')}</span>
                <select
                  aria-label={t('workbench.output.mapping.tag')}
                  value={selectedTagId}
                  onChange={(event) => setSelectedTagId(event.target.value)}
                  className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
                >
                  {candidates.map((candidate) => (
                    <option key={candidate.tagId} value={candidate.tagId}>
                      {candidate.tagKey}
                    </option>
                  ))}
                </select>
              </label>

              <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
                <span>{t('workbench.output.mapping.register')}</span>
                <input
                  aria-label={t('workbench.output.mapping.register')}
                  value={registerInput}
                  onChange={(event) => setRegisterInput(event.target.value)}
                  className="w-full rounded-lg border border-slate-800 bg-slate-950 px-3 py-2 text-sm text-slate-100"
                />
              </label>

              <div className="flex flex-wrap gap-2">
                <button
                  type="button"
                  onClick={() => void handleBind()}
                  disabled={isBusy || !selectedCandidate}
                  className="rounded-lg bg-cyan-500 px-3 py-2 text-sm font-semibold text-slate-950 disabled:opacity-50"
                >
                  {t('workbench.output.actions.bind')}
                </button>
                <button
                  type="button"
                  onClick={() => void handleDelete()}
                  disabled={isBusy || !selectedCandidate || selectedCandidate.register === null}
                  className="rounded-lg border border-slate-700 px-3 py-2 text-sm text-slate-200 disabled:opacity-50"
                >
                  {t('workbench.output.actions.delete')}
                </button>
                <button
                  type="button"
                  onClick={() => void handlePushValue()}
                  disabled={isBusy || !selectedCandidate}
                  className="rounded-lg border border-slate-700 px-3 py-2 text-sm text-slate-200 disabled:opacity-50"
                >
                  {t('workbench.output.actions.pushValue')}
                </button>
              </div>
            </div>

            {conflicts.length > 0 ? (
              <div
                role="status"
                aria-live="polite"
                className="space-y-2 rounded-2xl border border-amber-500/30 bg-amber-500/10 p-4 text-sm text-amber-100"
              >
                <p>{t('workbench.output.conflicts.summary')}</p>
                <ul className="space-y-1">
                  {conflicts.map((conflict) => (
                    <li key={conflict.register}>
                      {t('workbench.output.conflicts.item', {
                        register: conflict.register,
                        count: conflict.mappings.length,
                      })}
                    </li>
                  ))}
                </ul>
              </div>
            ) : null}

            {message ? (
              <p
                role="status"
                aria-live="polite"
                className="rounded-xl border border-slate-800 bg-slate-900/70 px-3 py-2 text-sm text-slate-200"
              >
                {message}
              </p>
            ) : null}
          </aside>
        </section>
      ) : null}

      {activeOutputTarget === 'database' ? (
        <DatabaseTargetBoard candidates={candidates} />
      ) : null}
    </div>
  );
}
