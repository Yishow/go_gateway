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
import { LocalModbusRegisterMapCanvas } from './LocalModbusRegisterMapCanvas';
import { SourceRuleDatabaseTargetBoard } from './SourceRuleDatabaseTargetBoard';
import { SourceRuleLocalModbusReviewSurface } from './SourceRuleLocalModbusReviewSurface';
import { useWorkbench } from './WorkbenchProvider';
import type { WorkbenchOutputCandidate } from './workbenchOutputTypes';
import type { AutoMapStrategy, DryRunResult } from './workbenchOutputTypes';
import {
  fromModbusDisplayRegister,
  modbusDisplayRegisterMax,
  modbusDisplayRegisterMin,
  toModbusDisplayRegister,
} from './workbenchOutputTypes';

type MappingConflict = {
  register: number;
  mappings: ModbusShareMapping[];
};

type OutputCandidate = WorkbenchOutputCandidate & {
  register: number | null;
  databasePath: string | null;
};

function formatModbusRegister(register: number): string {
  return `HR${toModbusDisplayRegister(register)}`;
}

function getSelectedDevice(devices: Device[], selectedDeviceId: string | null) {
  return devices.find((device) => device.id === selectedDeviceId) ?? null;
}

function getErrorMessage(error: unknown, fallback: string): string {
  return error instanceof Error ? error.message : fallback;
}

function dataTypeWidth(dataType: string): number {
  switch (dataType) {
    case 'int32':
    case 'uint32':
    case 'float32':
      return 2;
    case 'int64':
    case 'uint64':
    case 'float64':
      return 4;
    default:
      return 1;
  }
}

function getRegisterSlots(register: number, dataType: string): number[] {
  const width = dataTypeWidth(dataType);
  return Array.from({ length: width }, (_, index) => register + index);
}

function buildRegisterUsage(
  mappings: ModbusShareMapping[],
): Map<number, ModbusShareMapping[]> {
  const usage = new Map<number, ModbusShareMapping[]>();

  mappings.forEach((mapping) => {
    getRegisterSlots(mapping.register, mapping.data_type).forEach((slot) => {
      const list = usage.get(slot) ?? [];
      list.push(mapping);
      usage.set(slot, list);
    });
  });

  return usage;
}

function computeAutoMap(
  strategy: AutoMapStrategy,
  candidates: OutputCandidate[],
  existingMappings: ModbusShareMapping[],
): Array<{ tagId: string; register: number }> {
  const occupied = new Set<number>(buildRegisterUsage(existingMappings).keys());

  const unmapped = candidates.filter((c) => c.register === null);
  const result: Array<{ tagId: string; register: number }> = [];

  if (strategy === 'sequential') {
    let cursor = occupied.size > 0 ? Math.max(...occupied) + 1 : 0;
    for (const candidate of unmapped) {
      const width = dataTypeWidth(candidate.dataType);
      while (Array.from({ length: width }, (_, i) => cursor + i).some((r) => occupied.has(r))) {
        cursor++;
      }
      result.push({ tagId: candidate.tagId, register: cursor });
      for (let i = 0; i < width; i++) {
        occupied.add(cursor + i);
      }
      cursor += width;
    }
  } else if (strategy === 'gapAware') {
    let cursor = 0;
    for (const candidate of unmapped) {
      const width = dataTypeWidth(candidate.dataType);
      while (Array.from({ length: width }, (_, i) => cursor + i).some((r) => occupied.has(r))) {
        cursor++;
      }
      result.push({ tagId: candidate.tagId, register: cursor });
      for (let i = 0; i < width; i++) {
        occupied.add(cursor + i);
      }
      cursor += width;
    }
  } else {
    let cursor = 0;
    for (const candidate of unmapped) {
      const width = dataTypeWidth(candidate.dataType);
      const alignment = width;
      cursor = Math.ceil(cursor / alignment) * alignment;
      while (Array.from({ length: width }, (_, i) => cursor + i).some((r) => occupied.has(r))) {
        cursor += alignment;
      }
      result.push({ tagId: candidate.tagId, register: cursor });
      for (let i = 0; i < width; i++) {
        occupied.add(cursor + i);
      }
      cursor += width;
    }
  }

  return result;
}

function computeDryRun(
  candidates: OutputCandidate[],
  shareMappings: ModbusShareMapping[],
): DryRunResult[] {
  const registerUsage = buildRegisterUsage(shareMappings);

  return candidates.map((c) => {
    const mapping = shareMappings.find((m) => m.tag_id === c.tagId);
    if (!mapping) {
      return { tagId: c.tagId, tagKey: c.tagKey, register: -1, valid: false, reason: 'unmapped' };
    }
    const hasConflict = getRegisterSlots(mapping.register, mapping.data_type).some((slot) => {
      const users = registerUsage.get(slot) ?? [];
      return users.length > 1;
    });
    if (hasConflict) {
      return { tagId: c.tagId, tagKey: c.tagKey, register: mapping.register, valid: false, reason: 'conflict' };
    }
    return { tagId: c.tagId, tagKey: c.tagKey, register: mapping.register, valid: true };
  });
}

export function LocalModbusBoard() {
  const { t } = useTranslation();
  const tRef = useRef(t);
  tRef.current = t;
  const {
    activeOutputTarget,
    clearOutputSelectionState,
    crossStepContext,
    outputSelectionState,
    setInspectorSelection,
    setOutputSelectionState,
    setSelectedOutputTagId,
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
  const [registerInput, setRegisterInput] = useState(String(modbusDisplayRegisterMin));
  const [serverPortInput, setServerPortInput] = useState('5020');
  const [message, setMessage] = useState('');
  const [isBusy, setIsBusy] = useState(true);
  const [dryRunResults, setDryRunResults] = useState<DryRunResult[] | null>(null);

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
      setDryRunResults(null);
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
  const modbusSelectedTagId = outputSelectionState.modbus;
  const databaseSelectedTagId = outputSelectionState.database;
  const selectedTagId = outputSelectionState[activeOutputTarget];
  const selectedCandidate =
    candidates.find((candidate) => candidate.tagId === selectedTagId) ?? null;
  const selectedModbusCandidate =
    candidates.find((candidate) => candidate.tagId === modbusSelectedTagId) ?? null;
  const selectedModbusCandidateRegister = selectedModbusCandidate?.register ?? null;
  const requestedDisplayRegister = Number(registerInput);
  const requestedRegister = Number.isInteger(requestedDisplayRegister) &&
    requestedDisplayRegister >= modbusDisplayRegisterMin &&
    requestedDisplayRegister <= modbusDisplayRegisterMax
    ? fromModbusDisplayRegister(requestedDisplayRegister)
    : null;
  const highestMappedRegister = useMemo(
    () => shareMappings.reduce((max, mapping) => Math.max(max, mapping.register), 0),
    [shareMappings],
  );
  const viewportAnchorRegister = requestedRegister ?? selectedModbusCandidateRegister ?? highestMappedRegister;

  const { focusedTagIds } = crossStepContext;

  useEffect(() => {
    if (candidates.length === 0) {
      clearOutputSelectionState();
      setRegisterInput(String(modbusDisplayRegisterMin));
      return;
    }

    const fallbackSelectedTagId = focusedTagIds.find((id) =>
      candidates.some((candidate) => candidate.tagId === id),
    ) ?? candidates[0].tagId;
    const hasFocusedCandidate = focusedTagIds.some((id) =>
      candidates.some((candidate) => candidate.tagId === id),
    );

    setOutputSelectionState((currentSelection) => {
      let didChange = false;
      const nextSelection = { ...currentSelection };

      (['modbus', 'database'] as const).forEach((target) => {
        if (hasFocusedCandidate) {
          if (nextSelection[target] === fallbackSelectedTagId) {
            return;
          }
          nextSelection[target] = fallbackSelectedTagId;
          didChange = true;
          return;
        }

        if (candidates.some((candidate) => candidate.tagId === currentSelection[target])) {
          return;
        }
        if (nextSelection[target] === fallbackSelectedTagId) {
          return;
        }
        nextSelection[target] = fallbackSelectedTagId;
        didChange = true;
      });

      return didChange ? nextSelection : currentSelection;
    });
  }, [candidateKey, candidates, clearOutputSelectionState, focusedTagIds, setOutputSelectionState]);

  useEffect(() => {
    if (!modbusSelectedTagId) {
      return;
    }

    setRegisterInput(
      selectedModbusCandidateRegister !== null
        ? String(toModbusDisplayRegister(selectedModbusCandidateRegister))
        : String(modbusDisplayRegisterMin),
    );
  }, [modbusSelectedTagId, selectedModbusCandidateRegister]);

  const conflicts = useMemo<MappingConflict[]>(() => {
    return Array.from(buildRegisterUsage(shareMappings).entries())
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
    if (!selectedModbusCandidate) {
      setMessage(t('workbench.output.results.noTagSelected'));
      return;
    }

    const displayRegister = Number(registerInput);
    if (
      !Number.isInteger(displayRegister) ||
      displayRegister < modbusDisplayRegisterMin ||
      displayRegister > modbusDisplayRegisterMax
    ) {
      setMessage(t('workbench.output.results.invalidRegister'));
      return;
    }
    const register = fromModbusDisplayRegister(displayRegister);

    setIsBusy(true);

    try {
      await modbusShareAPI.upsertMapping(selectedModbusCandidate.tagId, register);
      await loadData();
      setMessage(
        t('workbench.output.results.mappingUpdated', {
          key: selectedModbusCandidate.tagKey,
          register: toModbusDisplayRegister(register),
        }),
      );
    } catch (error) {
      setMessage(
        getErrorMessage(error, t('workbench.output.results.mappingFallback')),
      );
    } finally {
      setIsBusy(false);
    }
  }, [loadData, registerInput, selectedModbusCandidate, t]);

  const handleDelete = useCallback(async () => {
    if (!selectedModbusCandidate) {
      setMessage(t('workbench.output.results.noTagSelected'));
      return;
    }

    setIsBusy(true);

    try {
      await modbusShareAPI.deleteMapping(selectedModbusCandidate.tagId);
      await loadData();
      setMessage(
        t('workbench.output.results.mappingDeleted', {
          key: selectedModbusCandidate.tagKey,
        }),
      );
    } catch (error) {
      setMessage(
        getErrorMessage(error, t('workbench.output.results.deleteFallback')),
      );
    } finally {
      setIsBusy(false);
    }
  }, [loadData, selectedModbusCandidate, t]);

  const handleSlotClick = useCallback(async (slotIndex: number, occupant: WorkbenchOutputCandidate | null) => {
    if (occupant) {
      setIsBusy(true);
      try {
        await modbusShareAPI.deleteMapping(occupant.tagId);
        await loadData();
        setMessage(
          t('workbench.output.results.slotUnbound', {
            key: occupant.tagKey,
            register: toModbusDisplayRegister(slotIndex),
          }),
        );
      } catch (error) {
        setMessage(
          getErrorMessage(error, t('workbench.output.results.deleteFallback')),
        );
      } finally {
        setIsBusy(false);
      }
      return;
    }

    if (!modbusSelectedTagId) {
      setMessage(t('workbench.output.results.noTagSelected'));
      return;
    }

    setIsBusy(true);
    try {
      await modbusShareAPI.upsertMapping(modbusSelectedTagId, slotIndex);
      await loadData();
      const boundTag = candidates.find((c) => c.tagId === modbusSelectedTagId);
      setMessage(
        t('workbench.output.results.slotBound', {
          key: boundTag?.tagKey ?? modbusSelectedTagId,
          register: toModbusDisplayRegister(slotIndex),
        }),
      );
    } catch (error) {
      setMessage(
        getErrorMessage(error, t('workbench.output.results.mappingFallback')),
      );
    } finally {
      setIsBusy(false);
    }
  }, [candidates, loadData, modbusSelectedTagId, t]);

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
    if (!selectedModbusCandidate) {
      setMessage(t('workbench.output.results.noTagSelected'));
      return;
    }

    if (selectedModbusCandidate.lastValue === null || selectedModbusCandidate.lastValue === undefined) {
      setMessage(t('workbench.output.results.noPointValue'));
      return;
    }

    setIsBusy(true);

    try {
      await modbusShareAPI.writeTagValue(
        selectedModbusCandidate.tagId,
        selectedModbusCandidate.lastValue,
      );
      setMessage(
        t('workbench.output.results.valuePushed', {
          key: selectedModbusCandidate.tagKey,
        }),
      );
    } catch (error) {
      setMessage(
        getErrorMessage(error, t('workbench.output.results.writeFallback')),
      );
    } finally {
      setIsBusy(false);
    }
  }, [selectedModbusCandidate, t]);

  const handleAutoMap = useCallback(
    async (strategy: AutoMapStrategy) => {
      const assignments = computeAutoMap(strategy, candidates, shareMappings);
      if (assignments.length === 0) {
        setMessage(t('workbench.output.modbusStudio.autoMap.noChanges'));
        return;
      }

      setIsBusy(true);
      try {
        await Promise.all(
          assignments.map((assignment) =>
            modbusShareAPI.upsertMapping(assignment.tagId, assignment.register),
          ),
        );
        await loadData();
        setMessage(
          t('workbench.output.modbusStudio.autoMap.applied', {
            count: assignments.length,
          }),
        );
      } catch (error) {
        setMessage(
          getErrorMessage(error, t('workbench.output.modbusStudio.autoMap.failed')),
        );
      } finally {
        setIsBusy(false);
      }
    },
    [candidates, loadData, shareMappings, t],
  );

  const handleDryRun = useCallback(() => {
    setDryRunResults(computeDryRun(candidates, shareMappings));
  }, [candidates, shareMappings]);

  if (!selectedDevice) {
    return (
      <section className="min-h-0 flex-1 space-y-6 overflow-y-auto overscroll-contain rounded-2xl border border-dashed border-slate-700 bg-slate-950/40 p-6">
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
      <section className="min-h-0 flex-1 space-y-6 overflow-y-auto overscroll-contain rounded-2xl border border-dashed border-slate-700 bg-slate-950/40 p-6">
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
    <div className="h-full min-h-0 space-y-6 overflow-y-auto overscroll-contain">
      <section
        className="space-y-3 rounded-2xl border border-slate-800 bg-slate-950/40 p-4"
        data-testid="output-primary-anchor"
        tabIndex={-1}
      >
        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex items-center gap-2">
            <span className="text-xs font-semibold uppercase tracking-[0.16em] text-cyan-300">
              {t('workbench.output.surface.activeTag')}
            </span>
            {selectedCandidate ? (
              <span className="rounded-full bg-cyan-500/10 px-3 py-1 text-xs font-medium text-cyan-200" data-testid="active-tag-badge">
                {selectedCandidate.tagKey}
              </span>
            ) : (
              <span className="text-xs text-slate-500">{t('workbench.output.surface.noneActive')}</span>
            )}
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

        <div className="flex flex-wrap gap-1.5" data-testid="output-tag-chips">
          {candidates.map((c) => {
            const isActive = c.tagId === selectedTagId;
            const modbusLabel = c.register !== null ? ` ${formatModbusRegister(c.register)}` : '';
            const dbLabel = c.databasePath ? ` DB` : '';
            const statusSuffix = activeOutputTarget === 'modbus' ? modbusLabel : dbLabel;
            return (
              <button
                key={c.tagId}
                type="button"
                data-testid={`output-candidate-${c.tagId}`}
                aria-pressed={isActive}
                onClick={() => {
                  setSelectedOutputTagId(activeOutputTarget, c.tagId);
                  setInspectorSelection({
                    kind: 'outputCandidate',
                    tagId: c.tagId,
                    target: activeOutputTarget,
                  });
                }}
                className={`rounded-lg border px-2.5 py-1.5 text-xs transition ${
                  isActive
                    ? 'border-cyan-400 bg-cyan-500/20 font-semibold text-cyan-100 ring-1 ring-cyan-400/50'
                    : 'border-slate-700 bg-slate-900/60 text-slate-300 hover:border-cyan-500/30'
                }`}
              >
                {c.tagKey}{statusSuffix ? <span className="ml-1 text-[10px] text-slate-400">{statusSuffix}</span> : null}
              </button>
            );
          })}
        </div>

        {selectedCandidate && (
          <p className="text-xs text-slate-400" data-testid="surface-hint">
            {activeOutputTarget === 'modbus'
              ? t('workbench.output.surface.clickSlotHint', { key: selectedCandidate.tagKey })
              : t('workbench.output.surface.clickColumnHint', { key: selectedCandidate.tagKey })}
          </p>
        )}
      </section>

      {activeOutputTarget === 'modbus' ? (
        <>
          <SourceRuleLocalModbusReviewSurface
            selectedTagId={selectedTagId}
            status={status}
            conflictCount={conflicts.length}
          />

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

              <LocalModbusRegisterMapCanvas
                candidates={candidates}
                shareMappings={shareMappings}
                conflicts={conflicts}
                selectedTagId={selectedTagId}
                viewportAnchorRegister={viewportAnchorRegister}
                onSlotClick={(slot, occupant) => void handleSlotClick(slot, occupant)}
                onAutoMap={(strategy) => void handleAutoMap(strategy)}
                onDryRun={handleDryRun}
                dryRunResults={dryRunResults}
                t={t}
              />
            </div>

            <aside
              className="space-y-6 rounded-2xl border border-slate-800 bg-slate-950/40 p-5"
              data-testid="modbus-secondary-panels"
              data-emphasis="supporting"
            >
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
                    data-testid="output-modbus-sync"
                    disabled={isBusy || !canSync}
                    onClick={() => void handleSync()}
                    className="rounded-lg border border-cyan-400/30 bg-cyan-500/20 px-3 py-2 text-sm font-medium text-cyan-100 disabled:opacity-50"
                  >
                    {t('workbench.output.actions.sync')}
                  </button>
                </div>
              </div>

              <div className="space-y-3 rounded-2xl border border-slate-800 bg-slate-900/70 p-4">
                <div className="rounded-xl border border-slate-800 bg-slate-950/70 p-3">
                  <p className="text-xs uppercase tracking-[0.18em] text-slate-400">
                    {t('workbench.output.mapping.selectedTag')}
                  </p>
                  <p className="mt-2 text-sm font-semibold text-slate-50">
                    {selectedModbusCandidate?.tagKey ?? '—'}
                  </p>
                </div>

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
                    disabled={isBusy || !selectedModbusCandidate}
                    className="rounded-lg bg-cyan-500 px-3 py-2 text-sm font-semibold text-slate-950 disabled:opacity-50"
                  >
                    {t('workbench.output.actions.bind')}
                  </button>
                  <button
                    type="button"
                    onClick={() => void handleDelete()}
                    disabled={isBusy || !selectedModbusCandidate || selectedModbusCandidate.register === null}
                    className="rounded-lg border border-slate-700 px-3 py-2 text-sm text-slate-200 disabled:opacity-50"
                  >
                    {t('workbench.output.actions.delete')}
                  </button>
                  <button
                    type="button"
                    onClick={() => void handlePushValue()}
                    disabled={isBusy || !selectedModbusCandidate}
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
                          register: toModbusDisplayRegister(conflict.register),
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
        </>
      ) : null}

      {activeOutputTarget === 'database' ? (
        <SourceRuleDatabaseTargetBoard
          candidates={candidates}
          selectedTagId={databaseSelectedTagId}
        />
      ) : null}
    </div>
  );
}
