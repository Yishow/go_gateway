import { useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import {
  normalizeNamingPrefix,
  SOURCE_PLANNER_ALLOWED_DATA_TYPES,
} from '../../../features/datalink/sourcePlannerContract';
import { useDevicesQuery } from '../../../hooks/datalink/useDevices';
import {
  useCreatePointMutation,
  usePointsQuery,
} from '../../../hooks/datalink/usePoints';
import type { DataType, Device } from '../../../types/datalink';
import { AddressCanvas } from './AddressCanvas';
import { AddressLedger } from './AddressLedger';
import { useWorkbench } from './WorkbenchProvider';
import {
  buildAddressCanvasItems,
  buildPlannedPointAddresses,
  getConflictingPlannedPointAddresses,
  type SourceViewMode,
} from './sourceCanvasModel';

function getSelectedDevice(devices: Device[], selectedDeviceId: string | null) {
  return devices.find((device) => device.id === selectedDeviceId) ?? null;
}

export function SourceCanvasSection() {
  const { t } = useTranslation();
  const { selectedDeviceId, setSelectedDeviceId } = useWorkbench();
  const { data: devices = [] } = useDevicesQuery();
  const { data: points = [] } = usePointsQuery(
    selectedDeviceId ? { device_id: selectedDeviceId } : undefined,
  );
  const createPointMutation = useCreatePointMutation();

  const selectedDevice = getSelectedDevice(devices, selectedDeviceId);
  const [viewMode, setViewMode] = useState<SourceViewMode>('grid');
  const [startAddress, setStartAddress] = useState('40001');
  const [dataType, setDataType] = useState<DataType>('int16');
  const [count, setCount] = useState(4);
  const [namingPrefix, setNamingPrefix] = useState('SRC');
  const [plannedPointAddresses, setPlannedPointAddresses] = useState<string[]>([]);
  const [batchCreateSummary, setBatchCreateSummary] = useState<{
    successCount: number;
    failureCount: number;
  } | null>(null);

  const items = useMemo(() => {
    if (!selectedDevice) {
      return [];
    }

    return buildAddressCanvasItems({
      points,
      plannedPointAddresses,
      plannedDataType: dataType,
      protocol: selectedDevice.protocol,
    });
  }, [dataType, plannedPointAddresses, points, selectedDevice]);

  const conflictingPointAddresses = useMemo(() => {
    if (!selectedDevice) {
      return [];
    }

    return getConflictingPlannedPointAddresses({
      points,
      plannedPointAddresses,
      plannedDataType: dataType,
      protocol: selectedDevice.protocol,
    });
  }, [dataType, plannedPointAddresses, points, selectedDevice]);

  const handleApplyPlan = () => {
    if (!selectedDevice) {
      return;
    }

    setPlannedPointAddresses(
      buildPlannedPointAddresses({
        startAddress,
        count,
        dataType,
        protocol: selectedDevice.protocol,
      }),
    );
    setBatchCreateSummary(null);
  };

  const handleBatchCreate = async () => {
    if (
      !selectedDevice ||
      plannedPointAddresses.length === 0 ||
      conflictingPointAddresses.length > 0
    ) {
      return;
    }

    const normalizedPrefix = normalizeNamingPrefix(namingPrefix);

    const results = await Promise.allSettled(
      plannedPointAddresses.map((address) =>
        createPointMutation.mutateAsync({
          device_id: selectedDevice.id,
          address,
          data_type: dataType,
          name: `${normalizedPrefix}_${address}`,
        }),
      ),
    );

    setBatchCreateSummary({
      successCount: results.filter((result) => result.status === 'fulfilled').length,
      failureCount: results.filter((result) => result.status === 'rejected').length,
    });
  };

  if (!selectedDevice) {
    return (
      <section className="space-y-6 rounded-2xl border border-dashed border-slate-700 bg-slate-950/40 p-6">
        <div className="space-y-2">
          <p className="text-xs font-semibold uppercase tracking-[0.2em] text-cyan-300">
            {t('workbench.source.empty.eyebrow')}
          </p>
          <h2 className="text-2xl font-semibold text-slate-50">
            {t('workbench.source.empty.title')}
          </h2>
          <p className="max-w-2xl text-sm text-slate-300">
            {t('workbench.source.empty.description')}
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

  return (
    <section className="space-y-6">
      <div className="rounded-2xl border border-slate-800 bg-slate-950/40 p-4">
        <div className="flex flex-col gap-4 xl:flex-row xl:items-end xl:justify-between">
          <div className="flex flex-wrap gap-2">
            <button
              type="button"
              onClick={() => setViewMode('grid')}
              className={`rounded-lg px-3 py-2 text-sm ${
                viewMode === 'grid'
                  ? 'bg-cyan-500 text-slate-950'
                  : 'bg-slate-900 text-slate-300'
              }`}
            >
              {t('workbench.source.view.grid')}
            </button>
            <button
              type="button"
              onClick={() => setViewMode('table')}
              className={`rounded-lg px-3 py-2 text-sm ${
                viewMode === 'table'
                  ? 'bg-cyan-500 text-slate-950'
                  : 'bg-slate-900 text-slate-300'
              }`}
            >
              {t('workbench.source.view.table')}
            </button>
          </div>

          <div className="grid gap-3 md:grid-cols-4">
            <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
              <span>{t('workbench.source.planner.startAddress')}</span>
              <input
                aria-label={t('workbench.source.planner.startAddress')}
                value={startAddress}
                onChange={(event) => setStartAddress(event.target.value)}
                className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
              />
            </label>
            <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
              <span>{t('workbench.source.planner.dataType')}</span>
              <select
                aria-label={t('workbench.source.planner.dataType')}
                value={dataType}
                onChange={(event) => setDataType(event.target.value as DataType)}
                className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
              >
                {SOURCE_PLANNER_ALLOWED_DATA_TYPES.map((option) => (
                  <option key={option} value={option}>
                    {option}
                  </option>
                ))}
              </select>
            </label>
            <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
              <span>{t('workbench.source.planner.count')}</span>
              <input
                aria-label={t('workbench.source.planner.count')}
                type="number"
                min={1}
                value={count}
                onChange={(event) => setCount(Number(event.target.value) || 0)}
                className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
              />
            </label>
            <label className="space-y-1 text-xs uppercase tracking-[0.16em] text-slate-400">
              <span>{t('workbench.source.planner.namingPrefix')}</span>
              <input
                aria-label={t('workbench.source.planner.namingPrefix')}
                value={namingPrefix}
                onChange={(event) => setNamingPrefix(event.target.value)}
                className="w-full rounded-lg border border-slate-800 bg-slate-900 px-3 py-2 text-sm text-slate-100"
              />
            </label>
            <div className="flex items-end gap-2">
              <button
                type="button"
                onClick={handleApplyPlan}
                className="rounded-lg bg-cyan-500 px-3 py-2 text-sm font-medium text-slate-950"
              >
                {t('workbench.source.planner.apply')}
              </button>
            </div>
          </div>
        </div>

        <div className="mt-4 flex flex-wrap gap-2">
          <button
            type="button"
            onClick={() => void handleBatchCreate()}
            disabled={
              plannedPointAddresses.length === 0 ||
              conflictingPointAddresses.length > 0 ||
              createPointMutation.isPending
            }
            className="rounded-lg border border-slate-800 px-3 py-2 text-sm text-slate-300 disabled:opacity-50"
          >
            {t('workbench.source.planner.batchCreate')}
          </button>
          <button
            type="button"
            disabled
            className="rounded-lg border border-slate-800 px-3 py-2 text-sm text-slate-300 disabled:opacity-50"
          >
            {t('workbench.source.planner.import')}
          </button>
          <button
            type="button"
            disabled
            className="rounded-lg border border-slate-800 px-3 py-2 text-sm text-slate-300 disabled:opacity-50"
          >
            {t('workbench.source.planner.export')}
          </button>
        </div>

        {conflictingPointAddresses.length > 0 ? (
          <p className="mt-4 text-sm text-rose-300">
            {t('workbench.source.planner.conflictHint')}
          </p>
        ) : null}

        {batchCreateSummary ? (
          <p className="mt-4 text-sm text-slate-300">
            {t('workbench.source.planner.batchSummary', batchCreateSummary)}
          </p>
        ) : null}
      </div>

      {viewMode === 'grid' ? <AddressCanvas items={items} /> : <AddressLedger items={items} />}
    </section>
  );
}
